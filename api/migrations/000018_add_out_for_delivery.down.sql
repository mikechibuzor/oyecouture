-- Any order sitting in out_for_delivery is moved back to dispatched first,
-- otherwise restoring the narrower constraint fails against live rows.
UPDATE orders
SET fulfillment_stage = 'dispatched'
WHERE fulfillment_stage = 'out_for_delivery';

ALTER TABLE orders
    DROP CONSTRAINT IF EXISTS orders_delivery_window_length,
    DROP CONSTRAINT IF EXISTS orders_rider_contact_length,
    DROP CONSTRAINT IF EXISTS orders_rider_name_length,
    DROP COLUMN IF EXISTS delivery_window,
    DROP COLUMN IF EXISTS rider_contact,
    DROP COLUMN IF EXISTS rider_name;

ALTER TABLE orders
    DROP CONSTRAINT IF EXISTS orders_fulfillment_stage_valid;

ALTER TABLE orders
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
    );
