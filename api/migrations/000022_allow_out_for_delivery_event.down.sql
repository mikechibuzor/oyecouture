-- Narrowing this again would orphan any event already recorded at the stage, so
-- the rows are removed first. That is lossy, which is why it belongs only in a
-- down migration nobody should need to run.
DELETE FROM order_fulfillment_events WHERE stage = 'out_for_delivery';

ALTER TABLE order_fulfillment_events
    DROP CONSTRAINT IF EXISTS order_fulfillment_events_stage_valid;

ALTER TABLE order_fulfillment_events
    ADD CONSTRAINT order_fulfillment_events_stage_valid CHECK (
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
    );
