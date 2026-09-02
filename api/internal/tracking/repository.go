package tracking

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/uzor-co/oyecouture/api/internal/email"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// Admin loads the operational tracking state and audited events for one order.
func (repository *PostgresRepository) Admin(
	requestContext context.Context,
	orderNo string,
) (AdminState, error) {
	const query = `
		SELECT id, order_no, kind, fulfillment_stage, fulfillment_version,
			estimated_delivery_date, COALESCE(carrier_name, ''),
			COALESCE(carrier_tracking_number, ''), COALESCE(carrier_tracking_url, ''),
			COALESCE(tracking_token, ''), customer_name, email,
			subtotal_kobo, delivery_fee_kobo, total_kobo
		FROM orders
		WHERE order_no = $1
	`

	var state AdminState
	var stage pgtype.Text
	var estimatedDeliveryDate pgtype.Date
	err := repository.pool.QueryRow(requestContext, query, orderNo).Scan(
		&state.orderID, &state.OrderNo, &state.Kind, &stage, &state.Version,
		&estimatedDeliveryDate, &state.Shipment.CarrierName,
		&state.Shipment.CarrierTrackingNumber, &state.Shipment.CarrierTrackingURL,
		&state.trackingToken, &state.customerName, &state.customerEmail,
		&state.subtotalKobo, &state.deliveryFeeKobo, &state.totalKobo,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return AdminState{}, ErrNotFound
		}
		return AdminState{}, fmt.Errorf("query admin tracking: %w", err)
	}
	if stage.Valid {
		state.CurrentStage = Stage(stage.String)
	}
	if estimatedDeliveryDate.Valid {
		date := estimatedDeliveryDate.Time
		state.EstimatedDeliveryDate = &date
	}

	events, err := repository.events(requestContext, state.orderID, true)
	if err != nil {
		return AdminState{}, err
	}
	state.Events = events

	// Needed by the dispatch and delivery emails, which reproduce the order.
	items, err := repository.items(requestContext, state.orderID)
	if err != nil {
		return AdminState{}, err
	}
	state.items = items

	return state, nil
}

// Public resolves only a valid bearer token and returns a minimized order view.
func (repository *PostgresRepository) Public(
	requestContext context.Context,
	trackingToken string,
) (PublicState, error) {
	const query = `
		SELECT id, order_no, kind, fulfillment_stage, estimated_delivery_date,
			COALESCE(carrier_name, ''), COALESCE(carrier_tracking_number, ''),
			COALESCE(carrier_tracking_url, '')
		FROM orders
		WHERE tracking_token = $1
			AND fulfillment_stage IS NOT NULL
	`

	var state PublicState
	var orderID string
	var estimatedDeliveryDate pgtype.Date
	if err := repository.pool.QueryRow(requestContext, query, trackingToken).Scan(
		&orderID, &state.OrderNo, &state.Kind, &state.CurrentStage,
		&estimatedDeliveryDate, &state.Shipment.CarrierName,
		&state.Shipment.CarrierTrackingNumber, &state.Shipment.CarrierTrackingURL,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return PublicState{}, ErrNotFound
		}
		return PublicState{}, fmt.Errorf("query public tracking: %w", err)
	}
	if estimatedDeliveryDate.Valid {
		date := estimatedDeliveryDate.Time
		state.EstimatedDeliveryDate = &date
	}

	items, err := repository.items(requestContext, orderID)
	if err != nil {
		return PublicState{}, err
	}
	state.Items = items
	state.Events, err = repository.events(requestContext, orderID, false)
	if err != nil {
		return PublicState{}, err
	}

	return state, nil
}

// Transition locks the order, verifies optimistic state, and atomically writes
// the stage, immutable event, legacy projection, and optional email job.
func (repository *PostgresRepository) Transition(
	requestContext context.Context,
	command TransitionCommand,
) error {
	transaction, err := repository.pool.Begin(requestContext)
	if err != nil {
		return fmt.Errorf("begin fulfillment transition: %w", err)
	}
	defer func() { _ = transaction.Rollback(requestContext) }()

	const lockQuery = `
		SELECT id, fulfillment_stage, fulfillment_version
		FROM orders
		WHERE order_no = $1
		FOR UPDATE
	`
	var orderID string
	var currentStage pgtype.Text
	var currentVersion int
	if err := transaction.QueryRow(requestContext, lockQuery, command.OrderNo).Scan(
		&orderID, &currentStage, &currentVersion,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("lock order tracking: %w", err)
	}
	if !currentStage.Valid {
		return ErrUnavailable
	}
	if Stage(currentStage.String) != command.CurrentStage || currentVersion != command.ExpectedVersion {
		return ErrVersionConflict
	}

	const insertEvent = `
		INSERT INTO order_fulfillment_events (
			order_id, stage, public_note, created_by_admin_id, occurred_at
		) VALUES ($1, $2, NULLIF($3, ''), NULLIF($4, '')::uuid, $5)
		RETURNING id
	`
	var eventID string
	if err := transaction.QueryRow(
		requestContext, insertEvent, orderID, string(command.NextStage), command.PublicNote,
		command.ActorAdminID, command.OccurredAt,
	).Scan(&eventID); err != nil {
		return fmt.Errorf("insert fulfillment event: %w", err)
	}

	const updateOrder = `
		UPDATE orders
		SET fulfillment_stage = $2,
			fulfillment_version = fulfillment_version + 1,
			status = $3,
			tracking_token = COALESCE(NULLIF($4, ''), tracking_token),
			estimated_delivery_date = COALESCE($5, estimated_delivery_date),
			carrier_name = CASE WHEN $2 = 'dispatched' THEN NULLIF($6, '') ELSE carrier_name END,
			carrier_tracking_number = CASE WHEN $2 = 'dispatched' THEN NULLIF($7, '') ELSE carrier_tracking_number END,
			carrier_tracking_url = CASE WHEN $2 = 'dispatched' THEN NULLIF($8, '') ELSE carrier_tracking_url END
		WHERE id = $1
	`
	if _, err := transaction.Exec(
		requestContext, updateOrder, orderID, string(command.NextStage), command.LegacyOrderStatus,
		command.TrackingToken, command.EstimatedDeliveryDate, command.CarrierName,
		command.CarrierTrackingNumber, command.CarrierTrackingURL,
	); err != nil {
		return fmt.Errorf("update order fulfillment: %w", err)
	}

	if command.Notification != nil {
		command.Notification.FulfillmentEventID = eventID
		if err := enqueueEmail(requestContext, transaction, *command.Notification); err != nil {
			return err
		}
	}

	if err := transaction.Commit(requestContext); err != nil {
		return fmt.Errorf("commit fulfillment transition: %w", err)
	}

	return nil
}

func (repository *PostgresRepository) events(
	requestContext context.Context,
	orderID string,
	includeActor bool,
) ([]Event, error) {
	const query = `
		SELECT event.id, event.stage, COALESCE(event.public_note, ''),
			CASE WHEN $2 THEN COALESCE(admin.email, '') ELSE '' END,
			event.occurred_at
		FROM order_fulfillment_events event
		LEFT JOIN admin_users admin ON admin.id = event.created_by_admin_id
		WHERE event.order_id = $1
		ORDER BY event.occurred_at, event.id
	`

	rows, err := repository.pool.Query(requestContext, query, orderID, includeActor)
	if err != nil {
		return nil, fmt.Errorf("query fulfillment events: %w", err)
	}
	defer rows.Close()

	events := []Event{}
	for rows.Next() {
		var event Event
		if err := rows.Scan(
			&event.ID, &event.Stage, &event.PublicNote, &event.ActorEmail, &event.OccurredAt,
		); err != nil {
			return nil, fmt.Errorf("scan fulfillment event: %w", err)
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read fulfillment events: %w", err)
	}

	return events, nil
}

func (repository *PostgresRepository) items(
	requestContext context.Context,
	orderID string,
) ([]Item, error) {
	const query = `
		SELECT name_snapshot, COALESCE(size, ''), quantity,
			price_kobo_snapshot * quantity
		FROM order_items
		WHERE order_id = $1
		ORDER BY id
	`
	rows, err := repository.pool.Query(requestContext, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("query tracking items: %w", err)
	}
	defer rows.Close()

	items := []Item{}
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.Name, &item.Size, &item.Quantity, &item.LineKobo); err != nil {
			return nil, fmt.Errorf("scan tracking item: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read tracking items: %w", err)
	}

	return items, nil
}

func enqueueEmail(
	requestContext context.Context,
	transaction pgx.Tx,
	event email.Event,
) error {
	const query = `
		INSERT INTO outbound_emails (
			order_id, fulfillment_event_id, type, provider, idempotency_key,
			recipient_email, subject, status, payload, next_attempt_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9::jsonb, NOW())
		ON CONFLICT (idempotency_key) DO NOTHING
	`
	if _, err := transaction.Exec(
		requestContext, query, event.OrderID, event.FulfillmentEventID, event.Type,
		event.Provider, event.IdempotencyKey, event.RecipientEmail, event.Subject,
		email.StatusPending, string(event.Payload),
	); err != nil {
		return fmt.Errorf("enqueue fulfillment email: %w", err)
	}

	return nil
}
