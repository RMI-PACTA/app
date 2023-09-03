BEGIN;

ALTER TABLE audit_log ALTER COLUMN secondary_target_type SET NOT NULL;
ALTER TABLE audit_log ALTER COLUMN actor_owner_id DROP NOT NULL;
ALTER TABLE ALTER COLUMN primary_target_owner_id DROP NOT NULL;
ALTER TABLE audit_log ADD CONSTRAINT audit_log_actor_owner_id_fkey FOREIGN KEY (actor_owner_id) REFERENCES owner(id) ON DELETE RESTRICT;
ALTER TABLE audit_log ADD CONSTRAINT audit_log_primary_target_owner_id_fkey FOREIGN KEY (primary_target_owner_id) REFERENCES owner(id) ON DELETE RESTRICT;
ALTER TABLE audit_log ADD CONSTRAINT audit_log_secondary_target_owner_id_fkey FOREIGN KEY (secondary_target_owner_id) REFERENCES owner(id) ON DELETE RESTRICT;
ALTER TABLE audit_log DROP COLUMN id;

COMMIT;
