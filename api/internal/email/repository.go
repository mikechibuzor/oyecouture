package email

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Reserve(context.Context, Event) (bool, error)
	MarkSent(context.Context, string, SendResult, []byte) error
	MarkFailed(context.Context, string, string, []byte) error
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresRepository stores outbound email events in PostgreSQL.
func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// Reserve creates or reopens one outbound email row for a retry-safe send.
func (repository *PostgresRepository) Reserve(
	requestContext context.Context,
	event Event,
) (bool, error) {
	const retryExisting = `
		UPDATE outbound_emails
		SET status = $2, recipient_email = $3, subject = $4,
			error_message = NULL, updated_at = NOW()
		WHERE idempotency_key = $1
			AND (
				status = $5
				OR (status = $2 AND updated_at < NOW() - INTERVAL '5 minutes')
			)
	`

	commandTag, err := repository.pool.Exec(
		requestContext, retryExisting, event.IdempotencyKey, StatusPending,
		event.RecipientEmail, event.Subject, StatusFailed,
	)
	if err != nil {
		return false, fmt.Errorf("reserve existing outbound email: %w", err)
	}
	if commandTag.RowsAffected() > 0 {
		return true, nil
	}

	const insert = `
		INSERT INTO outbound_emails (
			order_id, type, provider, idempotency_key, recipient_email, subject, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (idempotency_key) DO NOTHING
	`

	commandTag, err = repository.pool.Exec(
		requestContext, insert, event.OrderID, event.Type, event.Provider,
		event.IdempotencyKey, event.RecipientEmail, event.Subject, StatusPending,
	)
	if err != nil {
		return false, fmt.Errorf("reserve outbound email: %w", err)
	}

	return commandTag.RowsAffected() > 0, nil
}

// MarkSent records a successful provider delivery attempt.
func (repository *PostgresRepository) MarkSent(
	requestContext context.Context,
	idempotencyKey string,
	result SendResult,
	providerPayload []byte,
) error {
	const query = `
		UPDATE outbound_emails
		SET status = $2, provider_message_id = $3, provider_payload = $4::jsonb,
			error_message = NULL, attempts = attempts + 1, sent_at = NOW(), updated_at = NOW()
		WHERE idempotency_key = $1
	`

	if _, err := repository.pool.Exec(
		requestContext, query, idempotencyKey, StatusSent, result.MessageID,
		jsonPayload(providerPayload),
	); err != nil {
		return fmt.Errorf("mark outbound email sent: %w", err)
	}

	return nil
}

// MarkFailed records the provider or render error for later retry/reconciliation.
func (repository *PostgresRepository) MarkFailed(
	requestContext context.Context,
	idempotencyKey string,
	errorMessage string,
	providerPayload []byte,
) error {
	const query = `
		UPDATE outbound_emails
		SET status = $2, provider_payload = $3::jsonb, error_message = $4,
			attempts = attempts + 1, updated_at = NOW()
		WHERE idempotency_key = $1
	`

	if _, err := repository.pool.Exec(
		requestContext, query, idempotencyKey, StatusFailed,
		jsonPayload(providerPayload), errorMessage,
	); err != nil {
		return fmt.Errorf("mark outbound email failed: %w", err)
	}

	return nil
}

// jsonPayload stores provider responses as JSONB even when the provider returns text.
func jsonPayload(payload []byte) string {
	if len(payload) == 0 {
		return "{}"
	}
	if json.Valid(payload) {
		return string(payload)
	}

	wrappedPayload, err := json.Marshal(map[string]string{"raw": string(payload)})
	if err != nil {
		return "{}"
	}

	return string(wrappedPayload)
}
