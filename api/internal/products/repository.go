package products

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const productSelectQuery = `
SELECT
    product.id::text,
    product.slug,
    product.name,
    product.section,
    product.category,
    product.price_kobo,
    product.description,
    product.tag,
    product.has_sizes,
    product.sort_order,
    product.created_at,
    COALESCE(
        ARRAY_AGG(product_image.id::text ORDER BY product_image.position)
            FILTER (WHERE product_image.id IS NOT NULL),
        ARRAY[]::text[]
    ) AS image_ids,
    COALESCE(
        ARRAY_AGG(product_image.r2_key ORDER BY product_image.position)
            FILTER (WHERE product_image.id IS NOT NULL),
        ARRAY[]::text[]
    ) AS image_keys,
    COALESCE(
        ARRAY_AGG(product_image.position::integer ORDER BY product_image.position)
            FILTER (WHERE product_image.id IS NOT NULL),
        ARRAY[]::integer[]
    ) AS image_positions
FROM products AS product
LEFT JOIN product_images AS product_image ON product_image.product_id = product.id`

const productQuerySuffix = `
GROUP BY product.id
ORDER BY product.sort_order, product.created_at, product.id`

const adminProductSelectQuery = `
SELECT
    product.id::text,
    product.slug,
    product.name,
    product.section,
    product.category,
    product.price_kobo,
    product.description,
    product.tag,
    product.has_sizes,
    product.active,
    product.sort_order,
    product.created_at,
    COALESCE(
        ARRAY_AGG(product_image.id::text ORDER BY product_image.position)
            FILTER (WHERE product_image.id IS NOT NULL),
        ARRAY[]::text[]
    ) AS image_ids,
    COALESCE(
        ARRAY_AGG(product_image.r2_key ORDER BY product_image.position)
            FILTER (WHERE product_image.id IS NOT NULL),
        ARRAY[]::text[]
    ) AS image_keys,
    COALESCE(
        ARRAY_AGG(product_image.position::integer ORDER BY product_image.position)
            FILTER (WHERE product_image.id IS NOT NULL),
        ARRAY[]::integer[]
    ) AS image_positions
FROM products AS product
LEFT JOIN product_images AS product_image ON product_image.product_id = product.id`

const adminProductQuerySuffix = `
GROUP BY product.id
ORDER BY product.active DESC, product.sort_order, product.created_at, product.id
LIMIT $%d OFFSET $%d`

// Repository defines product persistence operations used by the service layer.
type Repository interface {
	List(context.Context, ListFilter) ([]Product, error)
	GetBySlug(context.Context, string) (Product, error)
}

// PostgresRepository reads storefront products from PostgreSQL.
type PostgresRepository struct {
	databasePool *pgxpool.Pool
}

// NewPostgresRepository creates a PostgreSQL-backed product repository.
func NewPostgresRepository(databasePool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{databasePool: databasePool}
}

// List returns active products matching the validated section and search filters.
func (repository *PostgresRepository) List(
	requestContext context.Context,
	filter ListFilter,
) ([]Product, error) {
	query, queryArguments := buildListQuery(filter)

	products, err := repository.queryProducts(requestContext, query, queryArguments...)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}

	return products, nil
}

// GetBySlug returns one active product for its exact canonical slug.
func (repository *PostgresRepository) GetBySlug(
	requestContext context.Context,
	productSlug string,
) (Product, error) {
	query := productSelectQuery + `
WHERE product.active = TRUE AND product.slug = $1` + productQuerySuffix

	products, err := repository.queryProducts(requestContext, query, productSlug)
	if err != nil {
		return Product{}, fmt.Errorf("get product by slug: %w", err)
	}
	if len(products) == 0 {
		return Product{}, ErrProductNotFound
	}

	return products[0], nil
}

// ListAdmin returns products for the owner, including archived products.
func (repository *PostgresRepository) ListAdmin(
	requestContext context.Context,
	filter AdminListFilter,
) (AdminProductPage, error) {
	query, countQuery, queryArguments := buildAdminListQuery(filter)

	products, err := repository.queryAdminProducts(requestContext, query, queryArguments...)
	if err != nil {
		return AdminProductPage{}, fmt.Errorf("list admin products: %w", err)
	}

	total, err := repository.countAdminProducts(requestContext, countQuery, queryArguments[:len(queryArguments)-2]...)
	if err != nil {
		return AdminProductPage{}, fmt.Errorf("count admin products: %w", err)
	}

	return AdminProductPage{
		Products: products,
		Total:    total,
		Limit:    filter.Limit,
		Offset:   filter.Offset,
	}, nil
}

// GetAdminBySlug returns one product for editing, regardless of active state.
func (repository *PostgresRepository) GetAdminBySlug(
	requestContext context.Context,
	productSlug string,
) (AdminProduct, error) {
	query := adminProductSelectQuery + `
WHERE product.slug = $1
GROUP BY product.id`

	products, err := repository.queryAdminProducts(requestContext, query, productSlug)
	if err != nil {
		return AdminProduct{}, fmt.Errorf("get admin product by slug: %w", err)
	}
	if len(products) == 0 {
		return AdminProduct{}, ErrProductNotFound
	}

	return products[0], nil
}

// CreateAdmin inserts one owner-managed product.
func (repository *PostgresRepository) CreateAdmin(
	requestContext context.Context,
	draft AdminProductDraft,
) (AdminProduct, error) {
	_, err := repository.databasePool.Exec(
		requestContext,
		`INSERT INTO products (
			slug, name, section, category, price_kobo, description, tag,
			has_sizes, active, sort_order
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		draft.Slug,
		draft.Name,
		draft.Section,
		draft.Category,
		draft.PriceKobo,
		draft.Description,
		draft.Tag,
		draft.HasSizes,
		draft.Active,
		draft.SortOrder,
	)
	if err != nil {
		if isUniqueConstraintViolation(err, "products_slug_key") {
			return AdminProduct{}, ErrDuplicateProductSlug
		}
		return AdminProduct{}, fmt.Errorf("create admin product: %w", err)
	}

	return repository.GetAdminBySlug(requestContext, draft.Slug)
}

// UpdateAdmin changes owner-editable product metadata.
func (repository *PostgresRepository) UpdateAdmin(
	requestContext context.Context,
	currentSlug string,
	draft AdminProductDraft,
) (AdminProduct, error) {
	commandTag, err := repository.databasePool.Exec(
		requestContext,
		`UPDATE products
		SET slug = $2,
			name = $3,
			section = $4,
			category = $5,
			price_kobo = $6,
			description = $7,
			tag = $8,
			has_sizes = $9,
			active = $10,
			sort_order = $11
		WHERE slug = $1`,
		currentSlug,
		draft.Slug,
		draft.Name,
		draft.Section,
		draft.Category,
		draft.PriceKobo,
		draft.Description,
		draft.Tag,
		draft.HasSizes,
		draft.Active,
		draft.SortOrder,
	)
	if err != nil {
		if isUniqueConstraintViolation(err, "products_slug_key") {
			return AdminProduct{}, ErrDuplicateProductSlug
		}
		return AdminProduct{}, fmt.Errorf("update admin product: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return AdminProduct{}, ErrProductNotFound
	}

	return repository.GetAdminBySlug(requestContext, draft.Slug)
}

// UpdateActive changes one product's storefront availability.
func (repository *PostgresRepository) UpdateActive(
	requestContext context.Context,
	productSlug string,
	active bool,
) error {
	commandTag, err := repository.databasePool.Exec(
		requestContext,
		`UPDATE products SET active = $2 WHERE slug = $1`,
		productSlug,
		active,
	)
	if err != nil {
		return fmt.Errorf("update product active: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return ErrProductNotFound
	}

	return nil
}

// AddImage appends one product image at the next gallery position.
func (repository *PostgresRepository) AddImage(
	requestContext context.Context,
	productSlug string,
	imageKey string,
) (AdminProduct, error) {
	const query = `
		INSERT INTO product_images (product_id, r2_key, position)
		SELECT
			product.id,
			$2,
			COALESCE(
				(SELECT MAX(position) + 1 FROM product_images WHERE product_id = product.id),
				1
			)
		FROM products AS product
		WHERE product.slug = $1
		RETURNING id
	`

	var imageID string
	if err := repository.databasePool.QueryRow(
		requestContext,
		query,
		productSlug,
		imageKey,
	).Scan(&imageID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return AdminProduct{}, ErrProductNotFound
		}
		if isCheckConstraintViolation(err, "product_images_position_valid") {
			return AdminProduct{}, ErrImageLimitReached
		}

		return AdminProduct{}, fmt.Errorf("add product image: %w", err)
	}

	return repository.GetAdminBySlug(requestContext, productSlug)
}

// DeleteImage removes one image from a product and compacts the remaining order.
func (repository *PostgresRepository) DeleteImage(
	requestContext context.Context,
	productSlug string,
	imageID string,
) (AdminProduct, error) {
	transaction, err := repository.databasePool.Begin(requestContext)
	if err != nil {
		return AdminProduct{}, fmt.Errorf("begin delete image: %w", err)
	}
	defer func() {
		_ = transaction.Rollback(requestContext)
	}()

	productID, err := productIDBySlug(requestContext, transaction, productSlug)
	if err != nil {
		return AdminProduct{}, err
	}

	commandTag, err := transaction.Exec(
		requestContext,
		`DELETE FROM product_images WHERE id = $1 AND product_id = $2`,
		imageID,
		productID,
	)
	if err != nil {
		return AdminProduct{}, fmt.Errorf("delete product image: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return AdminProduct{}, ErrImageNotFound
	}

	if err := compactImagePositions(requestContext, transaction, productID); err != nil {
		return AdminProduct{}, err
	}

	if err := transaction.Commit(requestContext); err != nil {
		return AdminProduct{}, fmt.Errorf("commit delete image: %w", err)
	}

	return repository.GetAdminBySlug(requestContext, productSlug)
}

// ReorderImages replaces product image positions with the requested complete order.
func (repository *PostgresRepository) ReorderImages(
	requestContext context.Context,
	productSlug string,
	imageIDs []string,
) (AdminProduct, error) {
	transaction, err := repository.databasePool.Begin(requestContext)
	if err != nil {
		return AdminProduct{}, fmt.Errorf("begin reorder images: %w", err)
	}
	defer func() {
		_ = transaction.Rollback(requestContext)
	}()

	productID, err := productIDBySlug(requestContext, transaction, productSlug)
	if err != nil {
		return AdminProduct{}, err
	}

	imageKeysByID, err := productImageKeysByID(requestContext, transaction, productID)
	if err != nil {
		return AdminProduct{}, err
	}
	if len(imageKeysByID) != len(imageIDs) {
		return AdminProduct{}, ErrInvalidImageReorder
	}

	orderedImageKeys := make([]string, 0, len(imageIDs))
	for _, imageID := range imageIDs {
		imageKey, exists := imageKeysByID[imageID]
		if !exists {
			return AdminProduct{}, ErrInvalidImageReorder
		}
		orderedImageKeys = append(orderedImageKeys, imageKey)
	}

	if _, err := transaction.Exec(
		requestContext,
		`DELETE FROM product_images WHERE product_id = $1`,
		productID,
	); err != nil {
		return AdminProduct{}, fmt.Errorf("clear product images for reorder: %w", err)
	}

	for imageIndex, imageID := range imageIDs {
		if _, err := transaction.Exec(
			requestContext,
			`INSERT INTO product_images (id, product_id, r2_key, position)
			VALUES ($1, $2, $3, $4)`,
			imageID,
			productID,
			orderedImageKeys[imageIndex],
			imageIndex+1,
		); err != nil {
			return AdminProduct{}, fmt.Errorf("insert reordered product image: %w", err)
		}
	}

	if err := transaction.Commit(requestContext); err != nil {
		return AdminProduct{}, fmt.Errorf("commit reorder images: %w", err)
	}

	return repository.GetAdminBySlug(requestContext, productSlug)
}

// isUniqueConstraintViolation identifies a named PostgreSQL unique constraint failure.
func isUniqueConstraintViolation(err error, constraintName string) bool {
	var postgresError *pgconn.PgError
	if !errors.As(err, &postgresError) {
		return false
	}

	return postgresError.Code == "23505" && postgresError.ConstraintName == constraintName
}

// isCheckConstraintViolation identifies a named PostgreSQL check constraint failure.
func isCheckConstraintViolation(err error, constraintName string) bool {
	var postgresError *pgconn.PgError
	if !errors.As(err, &postgresError) {
		return false
	}

	return postgresError.Code == "23514" && postgresError.ConstraintName == constraintName
}

type productTransaction interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

// productIDBySlug resolves a product slug inside an image transaction.
func productIDBySlug(
	requestContext context.Context,
	transaction productTransaction,
	productSlug string,
) (string, error) {
	var productID string
	if err := transaction.QueryRow(
		requestContext,
		`SELECT id::text FROM products WHERE slug = $1`,
		productSlug,
	).Scan(&productID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrProductNotFound
		}

		return "", fmt.Errorf("query product id: %w", err)
	}

	return productID, nil
}

// compactImagePositions keeps gallery positions contiguous after deletion.
func compactImagePositions(
	requestContext context.Context,
	transaction productTransaction,
	productID string,
) error {
	const query = `
		WITH ordered_images AS (
			SELECT
				id,
				ROW_NUMBER() OVER (ORDER BY position, id)::smallint AS position
			FROM product_images
			WHERE product_id = $1
		)
		UPDATE product_images
		SET position = ordered_images.position
		FROM ordered_images
		WHERE product_images.id = ordered_images.id
	`

	if _, err := transaction.Exec(requestContext, query, productID); err != nil {
		return fmt.Errorf("compact product image positions: %w", err)
	}

	return nil
}

// productImageKeysByID reads the current image set for exact-order validation.
func productImageKeysByID(
	requestContext context.Context,
	transaction productTransaction,
	productID string,
) (map[string]string, error) {
	rows, err := transaction.Query(
		requestContext,
		`SELECT id::text, r2_key FROM product_images WHERE product_id = $1`,
		productID,
	)
	if err != nil {
		return nil, fmt.Errorf("query product image keys: %w", err)
	}
	defer rows.Close()

	imageKeysByID := make(map[string]string)
	for rows.Next() {
		var imageID string
		var imageKey string
		if err := rows.Scan(&imageID, &imageKey); err != nil {
			return nil, fmt.Errorf("scan product image key: %w", err)
		}

		imageKeysByID[imageID] = imageKey
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate product image keys: %w", err)
	}

	return imageKeysByID, nil
}

// buildListQuery creates parameterized SQL for the validated optional filters.
func buildListQuery(filter ListFilter) (string, []any) {
	queryConditions := []string{"product.active = TRUE"}
	queryArguments := make([]any, 0, 2)

	// Add each optional filter with a positional parameter in deterministic order.
	if filter.Section != "" {
		queryArguments = append(queryArguments, filter.Section)
		queryConditions = append(
			queryConditions,
			fmt.Sprintf("product.section = $%d", len(queryArguments)),
		)
	}
	if filter.Query != "" {
		queryArguments = append(queryArguments, buildSearchPattern(filter.Query))
		searchParameter := fmt.Sprintf("$%d", len(queryArguments))
		// Section is searched as well as category because the storefront's own
		// Top Searches offer both kinds of term side by side: Kaftans and
		// Blazers are categories, while Suits and Casuals are sections. Without
		// section here the two section chips returned nothing at all, which is
		// what the client reported after launch.
		queryConditions = append(queryConditions, fmt.Sprintf(`(
    product.name ILIKE %[1]s ESCAPE '\'
    OR product.category ILIKE %[1]s ESCAPE '\'
    OR product.section ILIKE %[1]s ESCAPE '\'
    OR product.description ILIKE %[1]s ESCAPE '\'
)`, searchParameter))
	}

	query := productSelectQuery + "\nWHERE " + strings.Join(queryConditions, "\n    AND ")
	return query + productQuerySuffix, queryArguments
}

// buildAdminListQuery creates safe SQL for owner-facing catalogue management.
func buildAdminListQuery(filter AdminListFilter) (string, string, []any) {
	queryConditions := []string{"TRUE"}
	queryArguments := make([]any, 0, 5)

	if filter.Section != "" {
		queryArguments = append(queryArguments, filter.Section)
		queryConditions = append(
			queryConditions,
			fmt.Sprintf("product.section = $%d", len(queryArguments)),
		)
	}
	if filter.Active != nil {
		queryArguments = append(queryArguments, *filter.Active)
		queryConditions = append(
			queryConditions,
			fmt.Sprintf("product.active = $%d", len(queryArguments)),
		)
	}
	if filter.Query != "" {
		queryArguments = append(queryArguments, buildSearchPattern(filter.Query))
		searchParameter := fmt.Sprintf("$%d", len(queryArguments))
		// Matches the storefront search, so the owner looking for "casuals" in
		// the admin gets the same results a customer would.
		queryConditions = append(queryConditions, fmt.Sprintf(`(
    product.name ILIKE %[1]s ESCAPE '\'
    OR product.slug ILIKE %[1]s ESCAPE '\'
    OR product.category ILIKE %[1]s ESCAPE '\'
    OR product.section ILIKE %[1]s ESCAPE '\'
    OR product.description ILIKE %[1]s ESCAPE '\'
)`, searchParameter))
	}

	whereClause := "WHERE " + strings.Join(queryConditions, "\n    AND ")
	countQuery := "SELECT COUNT(*) FROM products AS product\n" + whereClause

	queryArguments = append(queryArguments, filter.Limit, filter.Offset)
	limitParameter := len(queryArguments) - 1
	offsetParameter := len(queryArguments)
	query := adminProductSelectQuery + "\n" + whereClause +
		fmt.Sprintf(adminProductQuerySuffix, limitParameter, offsetParameter)

	return query, countQuery, queryArguments
}

// buildSearchPattern escapes PostgreSQL wildcard characters before adding contains matching.
func buildSearchPattern(searchQuery string) string {
	escapedQuery := strings.NewReplacer(
		`\`, `\\`,
		`%`, `\%`,
		`_`, `\_`,
	).Replace(searchQuery)

	return "%" + escapedQuery + "%"
}

// queryProducts executes one product query and reconstructs its ordered image arrays.
func (repository *PostgresRepository) queryProducts(
	requestContext context.Context,
	query string,
	queryArguments ...any,
) ([]Product, error) {
	// Execute the aggregate query once to avoid per-product image lookups.
	productRows, err := repository.databasePool.Query(requestContext, query, queryArguments...)
	if err != nil {
		return nil, fmt.Errorf("query products: %w", err)
	}
	defer productRows.Close()

	// Scan each product row and rebuild its strongly typed ordered image collection.
	queryProducts := make([]Product, 0)
	for productRows.Next() {
		product, err := scanProduct(productRows)
		if err != nil {
			return nil, err
		}

		queryProducts = append(queryProducts, product)
	}
	if err := productRows.Err(); err != nil {
		return nil, fmt.Errorf("iterate products: %w", err)
	}

	return queryProducts, nil
}

// queryAdminProducts executes one owner-facing product query.
func (repository *PostgresRepository) queryAdminProducts(
	requestContext context.Context,
	query string,
	queryArguments ...any,
) ([]AdminProduct, error) {
	productRows, err := repository.databasePool.Query(requestContext, query, queryArguments...)
	if err != nil {
		return nil, fmt.Errorf("query admin products: %w", err)
	}
	defer productRows.Close()

	adminProducts := make([]AdminProduct, 0)
	for productRows.Next() {
		product, err := scanAdminProduct(productRows)
		if err != nil {
			return nil, err
		}

		adminProducts = append(adminProducts, product)
	}
	if err := productRows.Err(); err != nil {
		return nil, fmt.Errorf("iterate admin products: %w", err)
	}

	return adminProducts, nil
}

// countAdminProducts returns the total rows matching the current admin filters.
func (repository *PostgresRepository) countAdminProducts(
	requestContext context.Context,
	query string,
	queryArguments ...any,
) (int, error) {
	var total int
	if err := repository.databasePool.QueryRow(requestContext, query, queryArguments...).Scan(&total); err != nil {
		return 0, fmt.Errorf("count products: %w", err)
	}

	return total, nil
}

// scanProduct maps one aggregate database row into a product with ordered images.
func scanProduct(productRow pgx.Row) (Product, error) {
	var product Product
	var productSection string
	var imageIDs []string
	var imageKeys []string
	var imagePositions []int32

	if err := productRow.Scan(
		&product.ID,
		&product.Slug,
		&product.Name,
		&productSection,
		&product.Category,
		&product.PriceKobo,
		&product.Description,
		&product.Tag,
		&product.HasSizes,
		&product.SortOrder,
		&product.CreatedAt,
		&imageIDs,
		&imageKeys,
		&imagePositions,
	); err != nil {
		return Product{}, fmt.Errorf("scan product: %w", err)
	}

	// Reject inconsistent aggregate arrays instead of silently mismatching image data.
	if len(imageIDs) != len(imageKeys) || len(imageIDs) != len(imagePositions) {
		return Product{}, fmt.Errorf("scan product %q: inconsistent image data", product.Slug)
	}

	product.Section = Section(productSection)
	product.Images = make([]Image, 0, len(imageIDs))
	for imageIndex := range imageIDs {
		product.Images = append(product.Images, Image{
			ID:       imageIDs[imageIndex],
			R2Key:    imageKeys[imageIndex],
			Position: imagePositions[imageIndex],
		})
	}

	return product, nil
}

// scanAdminProduct maps one aggregate database row into the owner-facing model.
func scanAdminProduct(productRow pgx.Row) (AdminProduct, error) {
	var product AdminProduct
	var productSection string
	var imageIDs []string
	var imageKeys []string
	var imagePositions []int32

	if err := productRow.Scan(
		&product.ID,
		&product.Slug,
		&product.Name,
		&productSection,
		&product.Category,
		&product.PriceKobo,
		&product.Description,
		&product.Tag,
		&product.HasSizes,
		&product.Active,
		&product.SortOrder,
		&product.CreatedAt,
		&imageIDs,
		&imageKeys,
		&imagePositions,
	); err != nil {
		return AdminProduct{}, fmt.Errorf("scan admin product: %w", err)
	}

	if len(imageIDs) != len(imageKeys) || len(imageIDs) != len(imagePositions) {
		return AdminProduct{}, fmt.Errorf("scan admin product %q: inconsistent image data", product.Slug)
	}

	product.Section = Section(productSection)
	product.Images = make([]Image, 0, len(imageIDs))
	for imageIndex := range imageIDs {
		product.Images = append(product.Images, Image{
			ID:       imageIDs[imageIndex],
			R2Key:    imageKeys[imageIndex],
			Position: imagePositions[imageIndex],
		})
	}

	return product, nil
}
