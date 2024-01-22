BEGIN;

DROP INDEX user_name_gin_index;
DROP INDEX user_canonical_email_gin_index;

COMMIT;