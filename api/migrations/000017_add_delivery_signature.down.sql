ALTER TABLE orders
    DROP CONSTRAINT IF EXISTS orders_delivery_location_length,
    DROP CONSTRAINT IF EXISTS orders_delivery_signed_for_by_length,
    DROP COLUMN IF EXISTS delivery_location,
    DROP COLUMN IF EXISTS delivery_signed_at,
    DROP COLUMN IF EXISTS delivery_signed_for_by;
