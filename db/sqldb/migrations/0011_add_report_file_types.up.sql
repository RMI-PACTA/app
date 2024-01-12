BEGIN;

ALTER TYPE file_type ADD VALUE 'txt';
ALTER TYPE file_type ADD VALUE 'css';
ALTER TYPE file_type ADD VALUE 'js';
ALTER TYPE file_type ADD VALUE 'ttf';
ALTER TYPE file_type ADD VALUE ''; -- Unknown file types

COMMIT;
