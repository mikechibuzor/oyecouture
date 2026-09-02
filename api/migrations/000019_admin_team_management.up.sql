-- Team management for the ghost admin.
--
-- active exists so access can be revoked without deleting the account. A
-- deleted admin takes their name off every audit trail that references them,
-- which is the opposite of what an audit trail is for; a deactivated one keeps
-- the history and simply cannot sign in.
--
-- Deletion already revokes sessions: admin_sessions.admin_user_id is ON DELETE
-- CASCADE. Deactivation does not, so the session lookup has to check the owning
-- account is still active, and the code deletes their sessions on revoke.
ALTER TABLE admin_users
    ADD COLUMN active BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
