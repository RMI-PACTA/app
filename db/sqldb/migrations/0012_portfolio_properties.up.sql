BEGIN;

ALTER TABLE incomplete_upload 
    ADD COLUMN properties JSONB NOT NULL DEFAULT '{}';
ALTER TABLE incomplete_upload 
    DROP COLUMN holdings_date;

ALTER TABLE portfolio 
    ADD COLUMN properties JSONB NOT NULL DEFAULT '{}';
ALTER TABLE portfolio
    DROP COLUMN holdings_date;

COMMIT;