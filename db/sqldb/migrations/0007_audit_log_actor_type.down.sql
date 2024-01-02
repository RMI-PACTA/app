BEGIN;

-- There isn't a way to delete a value from an enum, so this is the workaround
-- https://stackoverflow.com/a/56777227/17909149

ALTER TABLE audit_log ALTER actor_type TYPE TEXT;

DROP TYPE audit_log_actor_type;
CREATE TYPE audit_log_actor_type AS ENUM (
    'USER',
    'ADMIN',
    'SUPER_ADMIN',
    'SYSTEM');

ALTER TABLE audit_log
    ALTER actor_type TYPE actor_type 
        USING actor_type::actor_type;

COMMIT;
