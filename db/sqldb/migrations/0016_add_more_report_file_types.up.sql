BEGIN;

ALTER TYPE file_type ADD VALUE 'xlsx';
ALTER TYPE file_type ADD VALUE 'rds';
ALTER TYPE file_type ADD VALUE 'css.map';

COMMIT;
