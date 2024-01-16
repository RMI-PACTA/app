BEGIN;

ALTER TABLE portfolio
    ADD COLUMN holdings_date TIMESTAMPTZ;
ALTER TABLE portfolio 
    DROP COLUMN properties;

ALTER TABLE incomplete_upload 
    ADD COLUMN holdings_date TIMESTAMPTZ;
ALTER TABLE incomplete_upload
    DROP COLUMN properties;

COMMIT;