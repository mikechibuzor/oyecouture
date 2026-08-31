ALTER TABLE orders
    ADD COLUMN fulfillment_stage TEXT,
    ADD COLUMN fulfillment_version INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN tracking_token TEXT,
    ADD COLUMN estimated_delivery_date DATE,
    ADD COLUMN carrier_name TEXT,
    ADD COLUMN carrier_tracking_number TEXT,
    ADD COLUMN carrier_tracking_url TEXT,
    ADD CONSTRAINT orders_fulfillment_stage_valid CHECK (
        fulfillment_stage IS NULL OR fulfillment_stage IN (
            'payment_confirmed',
            'measurements_approved',
            'preparing_order',
            'in_production',
            'quality_check',
            'ready_for_dispatch',
            'dispatched',
            'delivered'
        )
    ),
    ADD CONSTRAINT orders_fulfillment_version_nonnegative CHECK (fulfillment_version >= 0);

CREATE UNIQUE INDEX orders_tracking_token_unique_idx
    ON orders (tracking_token)
    WHERE tracking_token IS NOT NULL;

CREATE TABLE order_fulfillment_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    stage TEXT NOT NULL,
    public_note TEXT,
    created_by_admin_id UUID REFERENCES admin_users (id) ON DELETE SET NULL,
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT order_fulfillment_events_stage_valid CHECK (
        stage IN (
            'payment_confirmed',
            'measurements_approved',
            'preparing_order',
            'in_production',
            'quality_check',
            'ready_for_dispatch',
            'dispatched',
            'delivered'
        )
    ),
    CONSTRAINT order_fulfillment_events_note_length CHECK (
        public_note IS NULL OR char_length(public_note) <= 500
    )
);

CREATE INDEX order_fulfillment_events_order_id_idx
    ON order_fulfillment_events (order_id, occurred_at, id);

UPDATE orders
SET fulfillment_stage = CASE
    WHEN status = 'paid' THEN 'payment_confirmed'
    WHEN status = 'in_progress' AND kind = 'bespoke' THEN 'in_production'
    WHEN status = 'in_progress' THEN 'preparing_order'
    WHEN status = 'shipped' THEN 'dispatched'
    WHEN status = 'delivered' THEN 'delivered'
    ELSE fulfillment_stage
END
WHERE status IN ('paid', 'in_progress', 'shipped', 'delivered');

UPDATE orders
SET tracking_token = replace(gen_random_uuid()::text, '-', '') ||
    replace(gen_random_uuid()::text, '-', '')
WHERE fulfillment_stage IN ('dispatched', 'delivered')
    AND tracking_token IS NULL;

INSERT INTO order_fulfillment_events (order_id, stage, occurred_at)
SELECT id, fulfillment_stage, created_at
FROM orders
WHERE fulfillment_stage IS NOT NULL;

ALTER TABLE outbound_emails
    DROP CONSTRAINT outbound_emails_type_valid,
    DROP CONSTRAINT outbound_emails_status_valid,
    ADD COLUMN fulfillment_event_id UUID REFERENCES order_fulfillment_events (id) ON DELETE SET NULL,
    ADD COLUMN payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN locked_at TIMESTAMPTZ,
    ADD CONSTRAINT outbound_emails_type_valid CHECK (
        type IN (
            'order_confirmation',
            'payment_receipt',
            'fulfillment_update',
            'order_dispatched',
            'order_delivered'
        )
    ),
    ADD CONSTRAINT outbound_emails_status_valid CHECK (
        status IN ('pending', 'processing', 'sent', 'failed')
    );

CREATE INDEX outbound_emails_delivery_queue_idx
    ON outbound_emails (next_attempt_at, created_at)
    WHERE status IN ('pending', 'processing', 'failed');
