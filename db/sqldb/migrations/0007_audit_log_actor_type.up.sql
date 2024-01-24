BEGIN;

ALTER TYPE audit_log_actor_type ADD VALUE 'OWNER';
ALTER TYPE audit_log_actor_type ADD VALUE 'PUBLIC'; 

COMMIT;