package email

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Reserve(context.Context, Event) (bool, error)
	Enqueue(context.Context, Event) error
	ClaimNext(context.Context) (QueuedEvent, bool, error)
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

// Enqueue stores a fulfillment email for asynchronous, retryable delivery.
func (repository *PostgresRepository) Enqueue(
	requestContext context.Context,
	event Event,
) error {
	const query = `
		INSERT INTO outbound_emails (
			order_id, fulfillment_event_id, type, provider, idempotency_key,
			recipient_email, subject, status, payload, next_attempt_at
		) VALUES ($1, NULLIF($2, '')::uuid, $3, $4, $5, $6, $7, $8, $9::jsonb, NOW())
		ON CONFLICT (idempotency_key) DO NOTHING
	`

	payload := jsonPayload(event.Payload)
	if _, err := repository.pool.Exec(
		requestContext, query, event.OrderID, event.FulfillmentEventID,
		event.Type, event.Provider, event.IdempotencyKey, event.RecipientEmail,
		event.Subject, StatusPending, payload,
	); err != nil {
		return fmt.Errorf("enqueue outbound email: %w", err)
	}

	return nil
}

// ClaimNext atomically leases one due fulfillment email. Stale leases are
// reclaimed so a process crash cannot strand a customer notification.
func (repository *PostgresRepository) ClaimNext(
	requestContext context.Context,
) (QueuedEvent, bool, error) {
	const query = `
		WITH candidate AS (
			SELECT id
			FROM outbound_emails
			WHERE type IN ($1, $2, $3)
				AND next_attempt_at <= NOW()
				AND (
					status IN ($4, $5)
					OR (status = $6 AND locked_at < NOW() - INTERVAL '5 minutes')
				)
			ORDER BY created_at, id
			FOR UPDATE SKIP LOCKED
			LIMIT 1
		)
		UPDATE outbound_emails AS email
		SET status = $6, locked_at = NOW(), updated_at = NOW()
		FROM candidate
		WHERE email.id = candidate.id
		RETURNING email.idempotency_key, email.type, email.recipient_email,
			email.subject, email.payload
	`

	var event QueuedEvent
	err := repository.pool.QueryRow(
		requestContext, query,
		TypeFulfillmentUpdate, TypeOrderDispatched, TypeOrderDelivered,
		StatusPending, StatusFailed, StatusProcessing,
	).Scan(
		&event.IdempotencyKey, &event.Type, &event.RecipientEmail,
		&event.Subject, &event.Payload,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return QueuedEvent{}, false, nil
		}
		return QueuedEvent{}, false, fmt.Errorf("claim outbound email: %w", err)
	}

	return event, true, nil
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
			error_message = NULL, attempts = attempts + 1, sent_at = NOW(),
			locked_at = NULL, updated_at = NOW()
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
			attempts = attempts + 1, locked_at = NULL,
			next_attempt_at = NOW() + LEAST(60, POWER(2, LEAST(attempts, 6))::int) * INTERVAL '1 minute',
			updated_at = NOW()
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
