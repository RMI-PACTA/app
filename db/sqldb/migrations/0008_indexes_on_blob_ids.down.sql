BEGIN;

DROP INDEX portfolio_by_blob_id;
DROP INDEX incomplete_upload_by_blob_id;
DROP INDEX analysis_artifact_by_blob_id;

COMMIT;