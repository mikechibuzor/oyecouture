-- An immutable record of every change made through the admin.
--
-- The point of this table is that the shop owner can verify what was done to
-- their store rather than take anyone's word for it, the developer included.
-- Two things follow from that.
--
-- First, actor_email is snapshotted alongside actor_admin_id. Deleting an admin
-- sets the id to NULL, and without the snapshot their entries would become
-- anonymous - the trail would erase itself exactly when someone had reason to
-- erase it.
--
-- Second, rows cannot be changed or removed once written. That is enforced by a
-- trigger below rather than by convention, because "we do not edit it" and "it
-- cannot be edited" are different promises.
CREATE TABLE admin_audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_admin_id UUID REFERENCES admin_users (id) ON DELETE SET NULL,
    actor_email TEXT NOT NULL,
    actor_role TEXT NOT NULL,
    method TEXT NOT NULL,
    path TEXT NOT NULL,
    -- The HTTP status the request ended with. Refused attempts are recorded
    -- too: an audit log that only shows successes hides the more interesting
    -- half, such as someone repeatedly trying to reach a route they cannot.
    status_code INTEGER NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT admin_audit_log_actor_email_not_empty CHECK (actor_email <> ''),
    CONSTRAINT admin_audit_log_method_not_empty CHECK (method <> ''),
    CONSTRAINT admin_audit_log_path_not_empty CHECK (path <> '')
);

-- The log is read newest-first and nothing else, so one index serves it.
CREATE INDEX admin_audit_log_occurred_at_idx ON admin_audit_log (occurred_at DESC);

CREATE OR REPLACE FUNCTION admin_audit_log_is_append_only()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'admin_audit_log is append-only';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER admin_audit_log_no_update
    BEFORE UPDATE ON admin_audit_log
    FOR EACH ROW EXECUTE FUNCTION admin_audit_log_is_append_only();

CREATE TRIGGER admin_audit_log_no_delete
    BEFORE DELETE ON admin_audit_log
    FOR EACH ROW EXECUTE FUNCTION admin_audit_log_is_append_only();
