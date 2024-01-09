BEGIN;

DROP TABLE owner_merges;
DROP TABLE user_merges;

-- There isn't a way to delete a value from an enum, so this is the workaround
-- https://stackoverflow.com/a/56777227/17909149

ALTER TABLE audit_log ALTER action TYPE TEXT;

DROP TYPE audit_log_action;
CREATE TYPE audit_log_action AS ENUM (
    'CREATE',
    'UPDATE',
    'DELETE',
    'ADD_TO',
    'REMOVE_FROM',
    'ENABLE_ADMIN_DEBUG',
    'DISABLE_ADMIN_DEBUG',
    'DOWNLOAD',
    'ENABLE_SHARING',
    'DISABLE_SHARING');

ALTER TABLE audit_log 
    ALTER action TYPE audit_log_action 
        USING audit_log_action::audit_log_action;

COMMIT;
