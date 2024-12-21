BEGIN;

-- There isn't a way to delete a value from an enum, so this is the workaround
-- https://stackoverflow.com/a/56777227/17909149

ALTER TABLE analysis ALTER analysis_type TYPE TEXT;

DROP TYPE analysis_type;
CREATE TYPE analysis_type AS ENUM (
    'audit',
    'report'
);

ALTER TABLE analysis
    ALTER analysis_type TYPE analysis_type
        USING analysis_type::analysis_type;

COMMIT;
