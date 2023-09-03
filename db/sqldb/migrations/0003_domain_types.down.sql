BEGIN;

DROP TABLE audit_log;

DROP TYPE audit_log_target_type;

DROP TYPE audit_log_actor_type;

DROP TYPE audit_log_action;

DROP TABLE pacta_version;

DROP TABLE analysis_artifact;
	
DROP TABLE analysis;

DROP TYPE analysis_type;

DROP TABLE portfolio_snapshot;

DROP TABLE portfolio_initiative_membership;
	
DROP TABLE portfolio_group_membership;

DROP TABLE portfolio_group;

DROP TABLE portfolio;

DROP TABLE incomplete_upload;

DROP TYPE failure_code;

DROP TABLE owner;

DROP TABLE blob;

DROP TYPE file_type;

DROP TABLE initiative_user_relationship;

DROP TABLE initiative_invitation;

DROP TABLE initiative;

DROP TABLE pacta_user;

DROP TYPE authn_mechanism;

DROP TYPE language;

CREATE TYPE auth_provider AS ENUM ('GOOGLE', 'FACEBOOK', 'EMAIL_AND_PASS');

CREATE TABLE user_account (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  auth_provider_type auth_provider NOT NULL,
  auth_provider_id TEXT NOT NULL
);
CREATE INDEX account_auth_provider_id_idx ON user_account (auth_provider_id);

COMMIT;