package orders

import (
	"context"
	"fmt"
	"strings"
)

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
		filter.Limit, filter.Offset,
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
	`

	var total int
	if err := repository.pool.QueryRow(
		requestContext, query,
		string(filter.Status), string(filter.Kind), searchPattern,
	).Scan(&total); err != nil {
		return 0, fmt.Errorf("count orders: %w", err)
	}

	return total, nil
}

// escapeLikePattern neutralises the ILIKE wildcards so a customer searching for
// "100%" matches that literal text rather than every order.
func escapeLikePattern(value string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

	return replacer.Replace(value)
}
