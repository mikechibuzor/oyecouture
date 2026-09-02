-- Rows of the removed type would violate the narrower constraint, so they go
-- first. They are delivery notifications, not payment records.
DELETE FROM outbound_emails WHERE type = 'order_out_for_delivery';

ALTER TABLE outbound_emails
    DROP CONSTRAINT IF EXISTS outbound_emails_type_valid;

ALTER TABLE outbound_emails
    ADD CONSTRAINT outbound_emails_type_valid CHECK (
        type = ANY (ARRAY[
            'order_confirmation',
            'payment_receipt',
            'fulfillment_update',
            'order_dispatched',
            'order_delivered'
        ])
    );
