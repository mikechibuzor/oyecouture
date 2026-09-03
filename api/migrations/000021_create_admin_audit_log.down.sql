-- The triggers block DELETE, so they have to go before the table can be
-- dropped by anything that removes rows first.
DROP TRIGGER IF EXISTS admin_audit_log_no_delete ON admin_audit_log;
DROP TRIGGER IF EXISTS admin_audit_log_no_update ON admin_audit_log;
DROP FUNCTION IF EXISTS admin_audit_log_is_append_only();
DROP TABLE IF EXISTS admin_audit_log;
