BEGIN;

ALTER TYPE audit_log_action ADD VALUE 'TRANSFER_OWNERSHIP'; 

CREATE TABLE user_merges (
    from_user_id TEXT NOT NULL,
    to_user_id TEXT NOT NULL,
    actor_user_id TEXT NOT NULL,
    merged_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE owner_merges (
    from_owner_id TEXT NOT NULL,
    to_owner_id TEXT NOT NULL,
    actor_user_id TEXT NOT NULL,
    merged_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

COMMIT;