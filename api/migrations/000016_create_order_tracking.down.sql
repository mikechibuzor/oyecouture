DROP INDEX IF EXISTS outbound_emails_delivery_queue_idx;

DELETE FROM outbound_emails
WHERE type IN ('fulfillment_update', 'order_dispatched', 'order_delivered');

ALTER TABLE outbound_emails
    DROP CONSTRAINT outbound_emails_type_valid,
    DROP CONSTRAINT outbound_emails_status_valid,
    DROP COLUMN IF EXISTS fulfillment_event_id,
    DROP COLUMN IF EXISTS payload,
    DROP COLUMN IF EXISTS next_attempt_at,
    DROP COLUMN IF EXISTS locked_at,
    ADD CONSTRAINT outbound_emails_type_valid CHECK (
        type IN ('order_confirmation', 'payment_receipt')
    ),
    ADD CONSTRAINT outbound_emails_status_valid CHECK (
        status IN ('pending', 'sent', 'failed')
    );

DROP TABLE IF EXISTS order_fulfillment_events;
DROP INDEX IF EXISTS orders_tracking_token_unique_idx;

ALTER TABLE orders
    DROP CONSTRAINT IF EXISTS orders_fulfillment_stage_valid,
    DROP CONSTRAINT IF EXISTS orders_fulfillment_version_nonnegative,
    DROP COLUMN IF EXISTS fulfillment_stage,
    DROP COLUMN IF EXISTS fulfillment_version,
    DROP COLUMN IF EXISTS tracking_token,
    DROP COLUMN IF EXISTS estimated_delivery_date,
    DROP COLUMN IF EXISTS carrier_name,
    DROP COLUMN IF EXISTS carrier_tracking_number,
    DROP COLUMN IF EXISTS carrier_tracking_url;
