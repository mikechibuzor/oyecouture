package orders

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// ProductsBySlugs loads the trusted price and name for each requested slug.
// Inactive products are omitted, so an order can never be placed against one.
func (repository *PostgresRepository) ProductsBySlugs(
	requestContext context.Context,
	slugs []string,
) (map[string]CatalogueProduct, error) {
	const query = `
		SELECT id, slug, name, price_kobo, has_sizes
		FROM products
		WHERE active = TRUE AND slug = ANY($1)
	`

	rows, err := repository.pool.Query(requestContext, query, slugs)
	if err != nil {
		return nil, fmt.Errorf("query products: %w", err)
	}
	defer rows.Close()

	catalogue := make(map[string]CatalogueProduct)
	for rows.Next() {
		var product CatalogueProduct
		if err := rows.Scan(
			&product.ID, &product.Slug, &product.Name, &product.PriceKobo, &product.HasSizes,
		); err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}
		catalogue[product.Slug] = product
	}

	return catalogue, rows.Err()
}

// Create writes the order, its lines and any measurements in one transaction so
// a failure can never leave an order without the items it was priced from.
func (repository *PostgresRepository) Create(
	requestContext context.Context,
	order *Order,
	items []Item,
	measurements *Measurements,
	input CreateOrderInput,
) error {
	transaction, err := repository.pool.Begin(requestContext)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = transaction.Rollback(requestContext) }()

	const insertOrder = `
		INSERT INTO orders (
			kind, status, customer_name, email, phone,
			subtotal_kobo, total_kobo, pay_token
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, order_no, created_at
	`

	var orderID string
	if err := transaction.QueryRow(
		requestContext, insertOrder,
		order.Kind, order.Status, order.CustomerName, input.Email, input.Phone,
		order.SubtotalKobo, order.TotalKobo, order.PayToken,
	).Scan(&orderID, &order.OrderNo, &order.CreatedAt); err != nil {
		return fmt.Errorf("insert order: %w", err)
	}

	itemRows := make([][]any, 0, len(items))
	for _, item := range items {
		var size any
		if item.Size != "" {
			size = item.Size
		}
		itemRows = append(itemRows, []any{
			orderID, item.productID, item.Name, item.PriceKobo, size, item.Quantity,
		})
	}

	if _, err := transaction.CopyFrom(
		requestContext,
		pgx.Identifier{"order_items"},
		[]string{"order_id", "product_id", "name_snapshot", "price_kobo_snapshot", "size", "quantity"},
		pgx.CopyFromRows(itemRows),
	); err != nil {
		return fmt.Errorf("insert order items: %w", err)
	}

	if measurements != nil {
		const insertMeasurements = `
			INSERT INTO measurements (
				order_id, neck_cm, chest_cm, waist_cm, hips_cm,
				sleeve_cm, inseam_cm, length_cm, height_cm, weight_kg, notes
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`
		if _, err := transaction.Exec(
			requestContext, insertMeasurements, orderID,
			measurements.NeckCm, measurements.ChestCm, measurements.WaistCm, measurements.HipsCm,
			measurements.SleeveCm, measurements.InseamCm, measurements.LengthCm,
			measurements.HeightCm, measurements.WeightKg, measurements.Notes,
		); err != nil {
			return fmt.Errorf("insert measurements: %w", err)
		}
	}

	if err := transaction.Commit(requestContext); err != nil {
		return fmt.Errorf("commit order: %w", err)
	}

	return nil
}
