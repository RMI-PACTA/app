BEGIN;

CREATE EXTENSION IF NOT EXISTS pg_trgm;
SET pg_trgm.similarity_threshold = 0.1;

CREATE INDEX user_name_gin_index ON pacta_user USING gin (name gin_trgm_ops); 
CREATE INDEX user_canonical_email_gin_index ON pacta_user USING gin (canonical_email gin_trgm_ops); 

COMMIT;