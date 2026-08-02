-- The ghost admin's single operator account. Listed in the data-model doc but
-- never actually created, so it is created here alongside sessions.
CREATE TABLE admin_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'owner',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT admin_users_email_not_empty CHECK (email <> ''),
    CONSTRAINT admin_users_password_hash_not_empty CHECK (password_hash <> ''),
    CONSTRAINT admin_users_role_valid CHECK (role IN ('owner', 'staff'))
);

-- Email is the login identifier, so it must be unique and case-insensitive.
CREATE UNIQUE INDEX admin_users_email_unique_idx ON admin_users (lower(email));

-- Server-side sessions: revocable and they survive a restart, unlike a
-- self-contained token. The cookie carries only an opaque random id.
CREATE TABLE admin_sessions (
    token TEXT PRIMARY KEY,
    admin_user_id UUID NOT NULL REFERENCES admin_users (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT admin_sessions_token_not_empty CHECK (token <> '')
);

CREATE INDEX admin_sessions_admin_user_id_idx ON admin_sessions (admin_user_id);
CREATE INDEX admin_sessions_expires_at_idx ON admin_sessions (expires_at);
