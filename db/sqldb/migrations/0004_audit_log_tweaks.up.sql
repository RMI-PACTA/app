BEGIN;

ALTER TABLE audit_log ADD COLUMN id TEXT PRIMARY KEY;
ALTER TABLE audit_log DROP CONSTRAINT audit_log_actor_owner_id_fkey;
ALTER TABLE audit_log DROP CONSTRAINT audit_log_primary_target_owner_id_fkey;
ALTER TABLE audit_log DROP CONSTRAINT audit_log_secondary_target_owner_id_fkey;
ALTER TABLE audit_log ALTER COLUMN actor_owner_id SET NOT NULL;
ALTER TABLE audit_log ALTER COLUMN primary_target_owner_id SET NOT NULL;
ALTER TABLE audit_log ALTER COLUMN secondary_target_type DROP NOT NULL;

COMMIT;