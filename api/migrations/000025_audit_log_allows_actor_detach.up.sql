-- Removing a team member failed for anyone who had ever used the admin.
--
-- Two safeguards in migration 21 contradicted each other. admin_audit_log.actor_admin_id is
-- ON DELETE SET NULL, and that table's own header says why: deleting an admin detaches their
-- entries while actor_email keeps them attributable. But the append-only trigger refuses
-- every UPDATE, including the one the foreign key performs to do the detaching. So
-- DELETE /admin/team/{id} raised "admin_audit_log is append-only" and the Team page could not
-- remove anybody who had done anything.
--
-- The fix permits exactly that one change: actor_admin_id going from a value to NULL, with
-- every other column identical. Any other update, and every delete, is still refused.
--
-- What the table promises is unchanged. The evidence is actor_email, actor_role, method,
-- path, status_code and occurred_at, and none of those can be touched. Only the link to a
-- row that no longer exists can be dropped.
CREATE OR REPLACE FUNCTION admin_audit_log_is_append_only()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'UPDATE'
       AND OLD.actor_admin_id IS NOT NULL
       AND NEW.actor_admin_id IS NULL
       AND NEW.id          IS NOT DISTINCT FROM OLD.id
       AND NEW.actor_email IS NOT DISTINCT FROM OLD.actor_email
       AND NEW.actor_role  IS NOT DISTINCT FROM OLD.actor_role
       AND NEW.method      IS NOT DISTINCT FROM OLD.method
       AND NEW.path        IS NOT DISTINCT FROM OLD.path
       AND NEW.status_code IS NOT DISTINCT FROM OLD.status_code
       AND NEW.occurred_at IS NOT DISTINCT FROM OLD.occurred_at
    THEN
        RETURN NEW;
    END IF;

    RAISE EXCEPTION 'admin_audit_log is append-only';
END;
$$ LANGUAGE plpgsql;
