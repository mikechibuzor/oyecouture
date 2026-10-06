-- Back to refusing every update and delete, which also restores the bug this fixed.
CREATE OR REPLACE FUNCTION admin_audit_log_is_append_only()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'admin_audit_log is append-only';
END;
$$ LANGUAGE plpgsql;
