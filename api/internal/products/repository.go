package products

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
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
		queryConditions = append(queryConditions, fmt.Sprintf(`(
    product.name ILIKE %[1]s ESCAPE '\'
    OR product.category ILIKE %[1]s ESCAPE '\'
    OR product.description ILIKE %[1]s ESCAPE '\'
)`, searchParameter))
	}

	query := productSelectQuery + "\nWHERE " + strings.Join(queryConditions, "\n    AND ")
	return query + productQuerySuffix, queryArguments
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
