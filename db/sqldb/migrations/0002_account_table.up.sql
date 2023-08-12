BEGIN;

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