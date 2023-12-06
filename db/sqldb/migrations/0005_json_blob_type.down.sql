BEGIN;

-- There isn't a way to delete a value from an enum, so this is the workaround
-- https://stackoverflow.com/a/56777227/17909149

ALTER TABLE blob ALTER file_type TYPE TEXT;

DROP TYPE file_type;
CREATE TYPE file_type AS ENUM (
    'csv',
    'yaml',
    'zip',
    'html');

ALTER TABLE blob 
    ALTER file_type TYPE file_type 
        USING file_type::file_type;

COMMIT;
