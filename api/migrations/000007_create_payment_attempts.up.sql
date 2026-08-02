CREATE TABLE payment_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    provider TEXT NOT NULL,
    reference TEXT NOT NULL UNIQUE,
    amount_kobo BIGINT NOT NULL,
    currency TEXT NOT NULL DEFAULT 'NGN',
    status TEXT NOT NULL,
    authorization_url TEXT,
    access_code TEXT,
    provider_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    verified_at TIMESTAMPTZ,
    paid_at TIMESTAMPTZ,
    CONSTRAINT payment_attempts_provider_not_empty CHECK (provider <> ''),
    CONSTRAINT payment_attempts_reference_not_empty CHECK (reference <> ''),
    CONSTRAINT payment_attempts_amount_positive CHECK (amount_kobo > 0),
    CONSTRAINT payment_attempts_currency_not_empty CHECK (currency <> ''),
    CONSTRAINT payment_attempts_status_not_empty CHECK (status <> '')
);

CREATE INDEX payment_attempts_order_id_idx ON payment_attempts (order_id);
CREATE INDEX payment_attempts_status_idx ON payment_attempts (status);
