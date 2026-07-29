CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    section TEXT NOT NULL,
    category TEXT NOT NULL,
    price_kobo BIGINT NOT NULL,
    description TEXT NOT NULL,
    tag TEXT NOT NULL DEFAULT '',
    has_sizes BOOLEAN NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT products_slug_not_empty CHECK (slug <> ''),
    CONSTRAINT products_name_not_empty CHECK (name <> ''),
    CONSTRAINT products_section_valid CHECK (
        section IN ('suits', 'trad', 'casuals-men', 'casuals-women')
    ),
    CONSTRAINT products_category_not_empty CHECK (category <> ''),
    CONSTRAINT products_price_kobo_nonnegative CHECK (price_kobo >= 0),
    CONSTRAINT products_sort_order_nonnegative CHECK (sort_order >= 0)
);

CREATE INDEX products_active_section_sort_order_idx
    ON products (active, section, sort_order, created_at);

CREATE TABLE product_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    r2_key TEXT NOT NULL UNIQUE,
    position SMALLINT NOT NULL,
    CONSTRAINT product_images_r2_key_not_empty CHECK (r2_key <> ''),
    CONSTRAINT product_images_position_valid CHECK (position BETWEEN 1 AND 3),
    CONSTRAINT product_images_product_position_unique UNIQUE (product_id, position)
);

CREATE INDEX product_images_product_id_idx ON product_images (product_id);

CREATE TABLE settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    CONSTRAINT settings_key_not_empty CHECK (key <> '')
);

CREATE TABLE delivery_zones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    fee_kobo BIGINT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    CONSTRAINT delivery_zones_name_not_empty CHECK (name <> ''),
    CONSTRAINT delivery_zones_fee_kobo_nonnegative CHECK (fee_kobo >= 0)
);
