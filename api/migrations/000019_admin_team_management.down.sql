-- Any deactivated admin can sign in again once the column is gone, so their
-- sessions are cleared first rather than leaving dormant access behind.
DELETE FROM admin_sessions
WHERE admin_user_id IN (SELECT id FROM admin_users WHERE active = FALSE);

ALTER TABLE admin_users
    DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS active;
