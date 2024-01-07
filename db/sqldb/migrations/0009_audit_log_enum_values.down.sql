BEGIN;

-- There isn't a way to delete a value from an enum, so this is the workaround
-- https://stackoverflow.com/a/56777227/17909149

ALTER TABLE audit_log 
    ALTER primary_target_type TYPE TEXT,
    ALTER secondary_target_type TYPE TEXT,
    ALTER action TYPE TEXT;

DROP TYPE audit_log_target_type;
CREATE TYPE audit_log_target_type AS ENUM (
    'USER',
    'PORTFOLIO',
    'PORTFOLIO_GROUP',
    'INITIATIVE',
    'PACTA_VERSION',
    'ANALYSIS',
    'INCOMPLETE_UPLOAD');

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
    ALTER primary_target_type TYPE audit_log_target_type USING primary_target_type::audit_log_target_type,
    ALTER secondary_target_type TYPE audit_log_target_type USING secondary_target_type::audit_log_target_type,
    ALTER action TYPE audit_log_action USING action::audit_log_action;

COMMIT;
