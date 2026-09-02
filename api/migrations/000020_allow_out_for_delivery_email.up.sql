-- The out-for-delivery email type was added in Go without widening this
-- constraint, so enqueuing one violated the CHECK and rolled back the whole
-- fulfilment transition. Moving an order to out for delivery failed outright.
--
-- Any new email.Type* constant needs a line here too. TestEmailTypesAreAllowedByTheDatabase
-- exists to catch the next one before it reaches production.
ALTER TABLE outbound_emails
    DROP CONSTRAINT IF EXISTS outbound_emails_type_valid;

ALTER TABLE outbound_emails
    ADD CONSTRAINT outbound_emails_type_valid CHECK (
        type = ANY (ARRAY[
            'order_confirmation',
            'payment_receipt',
            'fulfillment_update',
            'order_dispatched',
            'order_out_for_delivery',
            'order_delivered'
        ])
    );
