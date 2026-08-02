CREATE TABLE outbound_emails (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    provider TEXT NOT NULL,
    idempotency_key TEXT NOT NULL UNIQUE,
    recipient_email TEXT NOT NULL,
    subject TEXT NOT NULL,
    status TEXT NOT NULL,
    provider_message_id TEXT,
    provider_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    error_message TEXT,
    attempts INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    sent_at TIMESTAMPTZ,
    CONSTRAINT outbound_emails_type_valid CHECK (type IN ('order_confirmation', 'payment_receipt')),
    CONSTRAINT outbound_emails_status_valid CHECK (status IN ('pending', 'sent', 'failed')),
    CONSTRAINT outbound_emails_provider_not_empty CHECK (provider <> ''),
    CONSTRAINT outbound_emails_recipient_email_not_empty CHECK (recipient_email <> ''),
    CONSTRAINT outbound_emails_subject_not_empty CHECK (subject <> ''),
    CONSTRAINT outbound_emails_attempts_nonnegative CHECK (attempts >= 0)
);

CREATE INDEX outbound_emails_order_id_idx ON outbound_emails (order_id);
CREATE INDEX outbound_emails_status_idx ON outbound_emails (status);
