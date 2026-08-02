package orders

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Detail reads the full admin view for one order.
func (repository *PostgresRepository) Detail(
	requestContext context.Context,
	orderNo string,
) (Detail, error) {
	const query = `
		SELECT
			o.id, o.order_no, o.kind, o.status, o.customer_name, o.email, o.phone,
			o.address, o.delivery_zone, o.delivery_fee_kobo,
			o.subtotal_kobo, o.total_kobo, o.pay_token,
			o.payment_provider, o.payment_ref, o.created_at,
			m.neck_cm::float8, m.chest_cm::float8, m.waist_cm::float8,
			m.hips_cm::float8, m.sleeve_cm::float8, m.inseam_cm::float8,
			m.length_cm::float8, m.height_cm::float8, m.weight_kg::float8,
			m.notes
		FROM orders o
		LEFT JOIN measurements m ON m.order_id = o.id
		WHERE o.order_no = $1
	`

	var detail Detail
	var address pgtype.Text
	var deliveryZone pgtype.Text
	var paymentProvider pgtype.Text
	var paymentRef pgtype.Text
	var measurementScalars measurementScan
	var measurementNotes pgtype.Text

	if err := repository.pool.QueryRow(requestContext, query, orderNo).Scan(
		&detail.orderID, &detail.OrderNo, &detail.Kind, &detail.Status,
		&detail.CustomerName, &detail.Email, &detail.Phone,
		&address, &deliveryZone, &detail.DeliveryFeeKobo,
		&detail.SubtotalKobo, &detail.TotalKobo, &detail.PayToken,
		&paymentProvider, &paymentRef, &detail.CreatedAt,
		&measurementScalars.NeckCm, &measurementScalars.ChestCm, &measurementScalars.WaistCm,
		&measurementScalars.HipsCm, &measurementScalars.SleeveCm,
		&measurementScalars.InseamCm, &measurementScalars.LengthCm,
		&measurementScalars.HeightCm, &measurementScalars.WeightKg, &measurementNotes,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Detail{}, ErrOrderNotFound
		}
		return Detail{}, fmt.Errorf("query order detail: %w", err)
	}

	detail.Address = textValue(address)
	detail.DeliveryZone = textValue(deliveryZone)
	detail.PaymentProvider = textValue(paymentProvider)
	detail.PaymentRef = textValue(paymentRef)
	detail.Measurements = measurementScalars.toMeasurements(textValue(measurementNotes))

	items, err := repository.itemsForOrder(requestContext, detail.orderID)
	if err != nil {
		return Detail{}, err
	}
	detail.Items = items

	return detail, nil
}

// UpdateStatus persists one service-approved lifecycle change.
func (repository *PostgresRepository) UpdateStatus(
	requestContext context.Context,
	orderNo string,
	status Status,
) error {
	const query = `UPDATE orders SET status = $2 WHERE order_no = $1`

	commandTag, err := repository.pool.Exec(requestContext, query, orderNo, string(status))
	if err != nil {
		return fmt.Errorf("update order status: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return ErrOrderNotFound
	}

	return nil
}

// itemsForOrder reads the priced item snapshots for one order.
func (repository *PostgresRepository) itemsForOrder(
	requestContext context.Context,
	orderID string,
) ([]Item, error) {
	const query = `
		SELECT
			COALESCE(p.slug, ''),
			oi.name_snapshot,
			COALESCE(oi.size, ''),
			oi.quantity,
			oi.price_kobo_snapshot,
			oi.price_kobo_snapshot * oi.quantity
		FROM order_items oi
		LEFT JOIN products p ON p.id = oi.product_id
		WHERE oi.order_id = $1
		ORDER BY oi.name_snapshot ASC, oi.id ASC
	`

	rows, err := repository.pool.Query(requestContext, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("query order items: %w", err)
	}
	defer rows.Close()

	items := []Item{}
	for rows.Next() {
		var item Item
		if err := rows.Scan(
			&item.Slug, &item.Name, &item.Size, &item.Quantity,
			&item.PriceKobo, &item.LineKobo,
		); err != nil {
			return nil, fmt.Errorf("scan order item: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read order items: %w", err)
	}

	return items, nil
}

type measurementScan struct {
	NeckCm   pgtype.Float8
	ChestCm  pgtype.Float8
	WaistCm  pgtype.Float8
	HipsCm   pgtype.Float8
	SleeveCm pgtype.Float8
	InseamCm pgtype.Float8
	LengthCm pgtype.Float8
	HeightCm pgtype.Float8
	WeightKg pgtype.Float8
}

// toMeasurements drops a missing measurement row while preserving submitted zeros.
func (scan measurementScan) toMeasurements(notes string) *Measurements {
	if !scan.hasValues() && notes == "" {
		return nil
	}

	return &Measurements{
		NeckCm:   floatPointer(scan.NeckCm),
		ChestCm:  floatPointer(scan.ChestCm),
		WaistCm:  floatPointer(scan.WaistCm),
		HipsCm:   floatPointer(scan.HipsCm),
		SleeveCm: floatPointer(scan.SleeveCm),
		InseamCm: floatPointer(scan.InseamCm),
		LengthCm: floatPointer(scan.LengthCm),
		HeightCm: floatPointer(scan.HeightCm),
		WeightKg: floatPointer(scan.WeightKg),
		Notes:    notes,
	}
}

// hasValues tells a fully absent row apart from an all-empty but noted one.
func (scan measurementScan) hasValues() bool {
	for _, value := range []pgtype.Float8{
		scan.NeckCm, scan.ChestCm, scan.WaistCm, scan.HipsCm, scan.SleeveCm,
		scan.InseamCm, scan.LengthCm, scan.HeightCm, scan.WeightKg,
	} {
		if value.Valid {
			return true
		}
	}

	return false
}

// floatPointer converts a nullable pgx value into the JSON shape the web uses.
func floatPointer(value pgtype.Float8) *float64 {
	if !value.Valid {
		return nil
	}

	return &value.Float64
}

// textValue converts nullable text columns to empty strings for omitted JSON fields.
func textValue(value pgtype.Text) string {
	if !value.Valid {
		return ""
	}

	return value.String
}

// List reads one page of orders for the admin, newest first. The total is taken
// in the same query with a window function, so paging costs one round trip.
func (repository *PostgresRepository) List(
	requestContext context.Context,
	filter ListFilter,
) (Page, error) {
	const query = `
		SELECT
			o.order_no, o.kind, o.status, o.customer_name, o.email, o.phone,
			(
				SELECT COALESCE(SUM(quantity), 0)
				FROM order_items
				WHERE order_id = o.id
			) AS item_count,
			o.total_kobo, o.created_at,
			COUNT(*) OVER () AS total_count
		FROM orders o
		WHERE ($1 = '' OR o.status = $1)
		  AND ($2 = '' OR o.kind = $2)
		  AND (
				$3 = ''
				OR o.order_no ILIKE $3 ESCAPE '\'
				OR o.customer_name ILIKE $3 ESCAPE '\'
				OR o.email ILIKE $3 ESCAPE '\'
				OR o.phone ILIKE $3 ESCAPE '\'
		  )
		  AND ($6::timestamptz IS NULL OR o.created_at >= $6)
		  AND ($7::timestamptz IS NULL OR o.created_at < $7)
		ORDER BY o.created_at DESC, o.order_no DESC
		LIMIT $4 OFFSET $5
	`

	searchPattern := ""
	if filter.Search != "" {
		searchPattern = "%" + escapeLikePattern(filter.Search) + "%"
	}

	rows, err := repository.pool.Query(
		requestContext, query,
		string(filter.Status), string(filter.Kind), searchPattern,
		filter.Limit, filter.Offset, filter.From, filter.To,
	)
	if err != nil {
		return Page{}, fmt.Errorf("query orders: %w", err)
	}
	defer rows.Close()

	page := Page{Orders: make([]Summary, 0, filter.Limit), Limit: filter.Limit, Offset: filter.Offset}
	for rows.Next() {
		var summary Summary
		var totalCount int
		if err := rows.Scan(
			&summary.OrderNo, &summary.Kind, &summary.Status, &summary.CustomerName,
			&summary.Email, &summary.Phone, &summary.ItemCount,
			&summary.TotalKobo, &summary.CreatedAt, &totalCount,
		); err != nil {
			return Page{}, fmt.Errorf("scan order: %w", err)
		}

		page.Orders = append(page.Orders, summary)
		page.Total = totalCount
	}
	if err := rows.Err(); err != nil {
		return Page{}, fmt.Errorf("read orders: %w", err)
	}

	// COUNT(*) OVER () only sees returned rows, so a page past the end would
	// report a total of zero and the UI could not tell "no matches" from "you
	// have paged too far". Only that one case pays for a second query.
	if len(page.Orders) == 0 && filter.Offset > 0 {
		total, err := repository.countOrders(requestContext, filter, searchPattern)
		if err != nil {
			return Page{}, err
		}
		page.Total = total
	}

	return page, nil
}

// countOrders totals the orders matching a filter, ignoring paging.
func (repository *PostgresRepository) countOrders(
	requestContext context.Context,
	filter ListFilter,
	searchPattern string,
) (int, error) {
	const query = `
		SELECT COUNT(*)
		FROM orders o
		WHERE ($1 = '' OR o.status = $1)
		  AND ($2 = '' OR o.kind = $2)
		  AND (
				$3 = ''
				OR o.order_no ILIKE $3 ESCAPE '\'
				OR o.customer_name ILIKE $3 ESCAPE '\'
				OR o.email ILIKE $3 ESCAPE '\'
				OR o.phone ILIKE $3 ESCAPE '\'
		  )
		  AND ($4::timestamptz IS NULL OR o.created_at >= $4)
		  AND ($5::timestamptz IS NULL OR o.created_at < $5)
	`

	var total int
	if err := repository.pool.QueryRow(
		requestContext, query,
		string(filter.Status), string(filter.Kind), searchPattern,
		filter.From, filter.To,
	).Scan(&total); err != nil {
		return 0, fmt.Errorf("count orders: %w", err)
	}

	return total, nil
}

// Metrics counts the orders by state in a single pass. A nil `since` covers
// all time.
func (repository *PostgresRepository) Metrics(
	requestContext context.Context,
	since *time.Time,
) (Metrics, error) {
	const query = `
		SELECT
			COUNT(*),
			COUNT(*) FILTER (WHERE status = $1),
			COUNT(*) FILTER (WHERE status = $2),
			COUNT(*) FILTER (WHERE status = ANY($3)),
			COALESCE(SUM(total_kobo) FILTER (WHERE status = ANY($3)), 0)
		FROM orders
		WHERE $4::timestamptz IS NULL OR created_at >= $4
	`

	paid := make([]string, 0, len(PaidStatuses))
	for _, status := range PaidStatuses {
		paid = append(paid, string(status))
	}

	var metrics Metrics
	if err := repository.pool.QueryRow(
		requestContext, query,
		string(StatusPendingConfirmation), string(StatusAwaitingPayment), paid, since,
	).Scan(
		&metrics.Total, &metrics.PendingConfirmation,
		&metrics.AwaitingPayment, &metrics.Paid, &metrics.RevenueKobo,
	); err != nil {
		return Metrics{}, fmt.Errorf("query order metrics: %w", err)
	}

	return metrics, nil
}

// escapeLikePattern neutralises the ILIKE wildcards so a customer searching for
// "100%" matches that literal text rather than every order.
func escapeLikePattern(value string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

	return replacer.Replace(value)
}
