package payments

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

const initializedAttemptUniqueIndex = "payment_attempts_one_initialized_attempt_per_order_idx"

type Repository interface {
	OrderByPayToken(context.Context, string) (Order, error)
	ReusableInitializedAttempt(context.Context, string, string, int64, string) (InitializeResult, bool, error)
	CreateAttempt(context.Context, Attempt) error
	MarkAttemptInitialized(context.Context, string, InitializeResult, []byte) error
	MarkAttemptFailed(context.Context, string, []byte) error
	RecordVerification(context.Context, string, Verification) error
	ConfirmSuccessfulPayment(context.Context, string, Verification) (Order, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (repository *PostgresRepository) OrderByPayToken(
	requestContext context.Context,
	payToken string,
) (Order, error) {
	const query = `
		SELECT id, order_no, pay_token, status, kind, customer_name, email, phone,
			total_kobo, created_at
		FROM orders
		WHERE pay_token = $1
	`

	var order Order
	if err := repository.pool.QueryRow(requestContext, query, payToken).Scan(
		&order.ID, &order.OrderNo, &order.PayToken, &order.Status, &order.Kind,
		&order.CustomerName, &order.Email, &order.Phone, &order.TotalKobo, &order.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Order{}, ErrOrderNotFound
		}
		return Order{}, fmt.Errorf("query payment order: %w", err)
	}

	items, err := repository.itemsForOrder(requestContext, order.ID)
	if err != nil {
		return Order{}, err
	}
	order.Items = items

	return order, nil
}

func (repository *PostgresRepository) ReusableInitializedAttempt(
	requestContext context.Context,
	orderID string,
	provider string,
	amountKobo int64,
	currency string,
) (InitializeResult, bool, error) {
	const query = `
		SELECT authorization_url, access_code, reference
		FROM payment_attempts
		WHERE order_id = $1
			AND provider = $2
			AND amount_kobo = $3
			AND currency = $4
			AND status = $5
			AND authorization_url IS NOT NULL
			AND authorization_url <> ''
		ORDER BY updated_at DESC, created_at DESC
		LIMIT 1
	`

	var result InitializeResult
	if err := repository.pool.QueryRow(
		requestContext, query, orderID, provider, amountKobo, currency, StatusInitialized,
	).Scan(&result.AuthorizationURL, &result.AccessCode, &result.Reference); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return InitializeResult{}, false, nil
		}
		return InitializeResult{}, false, fmt.Errorf("query reusable payment attempt: %w", err)
	}

	return result, true, nil
}

func (repository *PostgresRepository) CreateAttempt(
	requestContext context.Context,
	attempt Attempt,
) error {
	const query = `
		INSERT INTO payment_attempts (
			order_id, provider, reference, amount_kobo, currency, status, provider_payload
		) VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb)
	`

	payload := jsonPayload(attempt.ProviderPayload)
	if _, err := repository.pool.Exec(
		requestContext, query, attempt.OrderID, attempt.Provider, attempt.Reference,
		attempt.AmountKobo, attempt.Currency, attempt.Status, payload,
	); err != nil {
		return fmt.Errorf("insert payment attempt: %w", err)
	}

	return nil
}

func (repository *PostgresRepository) MarkAttemptInitialized(
	requestContext context.Context,
	reference string,
	result InitializeResult,
	providerPayload []byte,
) error {
	const query = `
		UPDATE payment_attempts
		SET status = $2, authorization_url = $3, access_code = $4,
			provider_payload = $5::jsonb, updated_at = NOW()
		WHERE reference = $1
	`

	if _, err := repository.pool.Exec(
		requestContext, query, reference, StatusInitialized,
		result.AuthorizationURL, result.AccessCode, jsonPayload(providerPayload),
	); err != nil {
		if isInitializedAttemptConflict(err) {
			return ErrPaymentInProgress
		}
		return fmt.Errorf("mark payment initialized: %w", err)
	}

	return nil
}

func (repository *PostgresRepository) MarkAttemptFailed(
	requestContext context.Context,
	reference string,
	providerPayload []byte,
) error {
	const query = `
		UPDATE payment_attempts
		SET status = $2, provider_payload = $3::jsonb, updated_at = NOW()
		WHERE reference = $1
	`

	if _, err := repository.pool.Exec(
		requestContext, query, reference, StatusFailed, jsonPayload(providerPayload),
	); err != nil {
		return fmt.Errorf("mark payment failed: %w", err)
	}

	return nil
}

func (repository *PostgresRepository) RecordVerification(
	requestContext context.Context,
	expectedOrderID string,
	verification Verification,
) error {
	const query = `
		UPDATE payment_attempts
		SET status = $2, provider_payload = $3::jsonb,
			verified_at = NOW(), paid_at = $4, updated_at = NOW()
		WHERE reference = $1
			AND ($5 = '' OR order_id = NULLIF($5, '')::uuid)
	`

	paidAt := nullableTime(verification.PaidAt)
	commandTag, err := repository.pool.Exec(
		requestContext, query, verification.Reference, verification.Status,
		jsonPayload(verification.ProviderPayload), paidAt, expectedOrderID,
	)
	if err != nil {
		return fmt.Errorf("record payment verification: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return ErrReferenceMismatch
	}

	return nil
}

func (repository *PostgresRepository) ConfirmSuccessfulPayment(
	requestContext context.Context,
	expectedOrderID string,
	verification Verification,
) (Order, error) {
	transaction, err := repository.pool.Begin(requestContext)
	if err != nil {
		return Order{}, fmt.Errorf("begin payment confirmation: %w", err)
	}
	defer func() { _ = transaction.Rollback(requestContext) }()

	const attemptQuery = `
		SELECT pa.order_id, pa.provider, pa.amount_kobo, pa.currency,
			o.order_no, o.pay_token, o.status, o.kind, o.customer_name, o.email, o.phone,
			o.total_kobo, o.created_at
		FROM payment_attempts pa
		JOIN orders o ON o.id = pa.order_id
		WHERE pa.reference = $1
			AND ($2 = '' OR pa.order_id = NULLIF($2, '')::uuid)
		FOR UPDATE OF pa, o
	`

	var attempt Attempt
	var order Order
	if err := transaction.QueryRow(requestContext, attemptQuery, verification.Reference, expectedOrderID).Scan(
		&attempt.OrderID, &attempt.Provider, &attempt.AmountKobo, &attempt.Currency,
		&order.OrderNo, &order.PayToken, &order.Status, &order.Kind,
		&order.CustomerName, &order.Email, &order.Phone, &order.TotalKobo, &order.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Order{}, ErrReferenceNotFound
		}
		return Order{}, fmt.Errorf("query payment attempt: %w", err)
	}
	order.ID = attempt.OrderID

	if err := verificationMatchesOrder(attempt, order, verification); err != nil {
		return Order{}, err
	}

	paidAt := nullableTime(verification.PaidAt)
	const updateAttempt = `
		UPDATE payment_attempts
		SET status = $2, provider_payload = $3::jsonb,
			verified_at = NOW(), paid_at = $4, updated_at = NOW()
		WHERE reference = $1
	`
	if _, err := transaction.Exec(
		requestContext, updateAttempt, verification.Reference, StatusSuccess,
		jsonPayload(verification.ProviderPayload), paidAt,
	); err != nil {
		return Order{}, fmt.Errorf("mark payment successful: %w", err)
	}

	if order.Status == OrderStatusPaid {
		if err := transaction.Commit(requestContext); err != nil {
			return Order{}, fmt.Errorf("commit idempotent payment confirmation: %w", err)
		}
		order.Items, err = repository.itemsForOrder(requestContext, order.ID)
		return order, err
	}
	if order.Status != OrderStatusAwaitingPayment {
		return Order{}, ErrOrderNotPayable
	}

	const updateOrder = `
		UPDATE orders
		SET status = $2, payment_provider = $3, payment_ref = $4
		WHERE id = $1
	`
	if _, err := transaction.Exec(
		requestContext, updateOrder, order.ID, OrderStatusPaid, attempt.Provider, verification.Reference,
	); err != nil {
		return Order{}, fmt.Errorf("mark order paid: %w", err)
	}

	if err := transaction.Commit(requestContext); err != nil {
		return Order{}, fmt.Errorf("commit payment confirmation: %w", err)
	}

	order.Status = OrderStatusPaid
	order.Items, err = repository.itemsForOrder(requestContext, order.ID)
	if err != nil {
		return Order{}, err
	}

	return order, nil
}

// verificationMatchesOrder is the guard that stops an order being marked paid
// for the wrong money. The provider's amount must equal both the attempt we
// created and the order total, so neither a tampered attempt nor a repriced
// order can slip through, and the currency must be the one we charge in.
//
// It is a pure function so it can be tested without a database: it is the last
// check standing between a partial payment and a fulfilled order.
func verificationMatchesOrder(attempt Attempt, order Order, verification Verification) error {
	if attempt.AmountKobo != verification.AmountKobo || order.TotalKobo != verification.AmountKobo {
		return ErrPaymentAmountMismatch
	}
	if attempt.Currency != verification.Currency || verification.Currency != CurrencyNGN {
		return ErrPaymentCurrency
	}

	return nil
}

func (repository *PostgresRepository) itemsForOrder(
	requestContext context.Context,
	orderID string,
) ([]Item, error) {
	const query = `
		SELECT name_snapshot, COALESCE(size, ''), quantity, price_kobo_snapshot * quantity
		FROM order_items
		WHERE order_id = $1
		ORDER BY name_snapshot ASC, id ASC
	`

	rows, err := repository.pool.Query(requestContext, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("query payment order items: %w", err)
	}
	defer rows.Close()

	items := []Item{}
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.Name, &item.Size, &item.Quantity, &item.LineKobo); err != nil {
			return nil, fmt.Errorf("scan payment order item: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read payment order items: %w", err)
	}

	return items, nil
}

func jsonPayload(payload []byte) string {
	if len(payload) == 0 {
		return "{}"
	}

	return string(payload)
}

func nullableTime(value *time.Time) pgtype.Timestamptz {
	if value == nil {
		return pgtype.Timestamptz{}
	}

	return pgtype.Timestamptz{Time: *value, Valid: true}
}

func isInitializedAttemptConflict(err error) bool {
	var postgresError *pgconn.PgError
	return errors.As(err, &postgresError) &&
		postgresError.Code == "23505" &&
		postgresError.ConstraintName == initializedAttemptUniqueIndex
}
