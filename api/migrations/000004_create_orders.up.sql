-- Human-readable order references (OYE-1001, OYE-1002, ...). A sequence keeps
-- them unique without a lookup, and they are what the customer quotes on WhatsApp.
CREATE SEQUENCE order_number_seq START 1001;

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_no TEXT NOT NULL UNIQUE DEFAULT 'OYE-' || nextval('order_number_seq'),
    kind TEXT NOT NULL,
    status TEXT NOT NULL,
    customer_name TEXT NOT NULL,
    email TEXT NOT NULL,
    phone TEXT NOT NULL,
    -- Address and delivery are unknown for bespoke orders: they are agreed on
    -- WhatsApp alongside measurements, so they stay nullable until checkout.
    address TEXT,
    delivery_zone TEXT,
    delivery_fee_kobo BIGINT NOT NULL DEFAULT 0,
    subtotal_kobo BIGINT NOT NULL,
    total_kobo BIGINT NOT NULL,
    -- Opaque token behind the /pay/<token> link the owner shares once measurements
    -- are agreed. Generated up front so no second write is needed later.
    pay_token TEXT NOT NULL UNIQUE,
    payment_provider TEXT,
    payment_ref TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT orders_kind_valid CHECK (kind IN ('bespoke', 'ready-to-wear')),
    CONSTRAINT orders_status_valid CHECK (
        status IN (
            'pending_confirmation',
            'awaiting_payment',
            'paid',
            'in_progress',
            'shipped',
            'delivered',
            'cancelled',
            'refunded'
        )
    ),
    CONSTRAINT orders_customer_name_not_empty CHECK (customer_name <> ''),
    CONSTRAINT orders_email_not_empty CHECK (email <> ''),
    CONSTRAINT orders_phone_not_empty CHECK (phone <> ''),
    CONSTRAINT orders_money_nonnegative CHECK (
        delivery_fee_kobo >= 0 AND subtotal_kobo >= 0 AND total_kobo >= 0
    )
);

CREATE INDEX orders_created_at_idx ON orders (created_at DESC);
CREATE INDEX orders_status_idx ON orders (status);

CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    -- Products are never hard-deleted, but keep the order readable even if one is.
    product_id UUID REFERENCES products (id) ON DELETE SET NULL,
    -- Snapshots: what the customer actually bought, at the price actually charged.
    name_snapshot TEXT NOT NULL,
    price_kobo_snapshot BIGINT NOT NULL,
    size TEXT,
    quantity INTEGER NOT NULL,
    CONSTRAINT order_items_name_not_empty CHECK (name_snapshot <> ''),
    CONSTRAINT order_items_price_nonnegative CHECK (price_kobo_snapshot >= 0),
    CONSTRAINT order_items_quantity_valid CHECK (quantity BETWEEN 1 AND 10)
);

CREATE INDEX order_items_order_id_idx ON order_items (order_id);

-- Optional: submitted by customers who already know their measurements, so the
-- WhatsApp conversation can skip straight to fabric and delivery.
CREATE TABLE measurements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL UNIQUE REFERENCES orders (id) ON DELETE CASCADE,
    neck_cm NUMERIC(5, 1),
    chest_cm NUMERIC(5, 1),
    waist_cm NUMERIC(5, 1),
    hips_cm NUMERIC(5, 1),
    sleeve_cm NUMERIC(5, 1),
    inseam_cm NUMERIC(5, 1),
    length_cm NUMERIC(5, 1),
    height_cm NUMERIC(5, 1),
    weight_kg NUMERIC(5, 1),
    notes TEXT
);
