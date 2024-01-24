BEGIN;

-- Creates indexes on blob_id columns for faster lookups when performing ownership lookups.
CREATE INDEX analysis_artifact_by_blob_id ON analysis_artifact (blob_id);
CREATE INDEX incomplete_upload_by_blob_id ON incomplete_upload (blob_id);
CREATE INDEX portfolio_by_blob_id ON portfolio (blob_id);

COMMIT;