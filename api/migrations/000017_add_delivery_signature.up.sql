-- Who actually took the parcel, when, and where. The delivered fulfilment event
-- already records the moment an admin marked the order delivered, which is not
-- the same fact: handover happens in the world, the click happens later. The
-- client's delivered email shows the handover, so it is stored separately.
ALTER TABLE orders
    ADD COLUMN delivery_signed_for_by TEXT,
    ADD COLUMN delivery_signed_at TIMESTAMPTZ,
    ADD COLUMN delivery_location TEXT,
    ADD CONSTRAINT orders_delivery_signed_for_by_length CHECK (
        delivery_signed_for_by IS NULL OR char_length(delivery_signed_for_by) <= 120
    ),
    ADD CONSTRAINT orders_delivery_location_length CHECK (
        delivery_location IS NULL OR char_length(delivery_location) <= 240
    );
