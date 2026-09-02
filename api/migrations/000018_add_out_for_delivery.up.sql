-- Adds the out-for-delivery stage and the rider details it carries.
--
-- It sits between dispatched and delivered but is deliberately optional: a DHL
-- parcel goes straight from dispatched to delivered, while a Lagos order is
-- handed to Oyecouture's own rider and gets this step. Forcing every order
-- through it would make international orders sit in a stage that never ends.
--
-- Rider fields are separate from the carrier ones already on the table. A
-- courier is a company with a tracking number; a rider is a person with a phone
-- and a time window, and the customer needs entirely different things from each.
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
            'out_for_delivery',
            'delivered'
        )
    );

ALTER TABLE orders
    ADD COLUMN rider_name TEXT,
    ADD COLUMN rider_contact TEXT,
    ADD COLUMN delivery_window TEXT,
    ADD CONSTRAINT orders_rider_name_length CHECK (
        rider_name IS NULL OR char_length(rider_name) <= 120
    ),
    ADD CONSTRAINT orders_rider_contact_length CHECK (
        rider_contact IS NULL OR char_length(rider_contact) <= 60
    ),
    ADD CONSTRAINT orders_delivery_window_length CHECK (
        delivery_window IS NULL OR char_length(delivery_window) <= 120
    );
