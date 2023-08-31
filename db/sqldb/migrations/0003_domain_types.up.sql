BEGIN;

DROP TABLE user_account;

DROP TYPE auth_provider;

CREATE TYPE language AS ENUM (
    'en',
    'de',
    'fr',
    'es'
);

CREATE TYPE authn_mechanism AS ENUM (
    'EMAIL_AND_PASS'
);

CREATE TABLE pacta_user (
    id TEXT PRIMARY KEY,
    authn_mechanism authn_mechanism NOT NULL,
    authn_id TEXT NOT NULL,
    -- What the user actually entered during sign up.
    entered_email TEXT NOT NULL UNIQUE,
    -- The validated, cleaned, consistently cased email address - excludes plus aliases, for example.
    canonical_email TEXT NOT NULL UNIQUE,
    admin BOOLEAN NOT NULL,
    super_admin BOOLEAN NOT NULL,
    name TEXT NOT NULL,
    -- Null until the user explicitly choses something. Will be defaulted based on the domain accessed from.
    preferred_language language,
    created_at TIMESTAMPTZ NOT NULL,
    UNIQUE(authn_mechanism, authn_id)
);

CREATE TABLE pacta_version (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    digest TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL,
    -- NULL = is not default, TRUE = is default, enforced via the checks below.
    is_default BOOLEAN
);
ALTER TABLE pacta_version ADD CONSTRAINT is_default_is_true_or_null CHECK (is_default);
ALTER TABLE pacta_version ADD CONSTRAINT is_default_only_1_true UNIQUE (is_default);

CREATE TABLE initiative (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    affiliation TEXT NOT NULL,
    public_description TEXT NOT NULL,
    internal_description TEXT NOT NULL,
    requires_invitation_to_join BOOLEAN NOT NULL,
    is_accepting_new_members BOOLEAN NOT NULL,
    is_accepting_new_portfolios BOOLEAN NOT NULL,
    pacta_version_id TEXT NOT NULL REFERENCES pacta_version (id) ON DELETE RESTRICT,
    language language NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE initiative_invitation (
    id TEXT PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    initiative_id TEXT NOT NULL REFERENCES initiative (id) ON DELETE RESTRICT,
    used_by_user_id TEXT REFERENCES pacta_user (id) ON DELETE RESTRICT
);

CREATE TABLE initiative_user_relationship (
    user_id TEXT NOT NULL,
    initiative_id TEXT NOT NULL,
    manager BOOLEAN NOT NULL,
    member BOOLEAN NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY(user_id, initiative_id)
);

CREATE TYPE file_type AS ENUM (
    'csv',
    'yaml',
    'zip',
    'html'
);

CREATE TABLE blob (
    id TEXT PRIMARY KEY,
    blob_uri TEXT NOT NULL UNIQUE,
    file_type file_type NOT NULL,
    file_name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

-- We use the owner abstraction since there are (at minimum) two types of entities that can own data:
-- users and initiatives. As we expect that complexity may increase in the future (with organizations)
-- I think this indirection is worthwhile.
CREATE TABLE owner (
    id TEXT PRIMARY KEY,
    user_id TEXT REFERENCES pacta_user (id) ON DELETE RESTRICT,
    initiative_id TEXT REFERENCES initiative (id) ON DELETE RESTRICT
);
ALTER TABLE owner ADD CONSTRAINT owner_is_always_well_defined CHECK (NUM_NONNULLS(user_id, initiative_id) = 1);
CREATE INDEX owner_by_user_id ON owner USING btree (user_id);
CREATE INDEX owner_by_initiative_id ON owner USING btree (initiative_id);

CREATE TYPE failure_code AS ENUM ('UNKNOWN');

CREATE TABLE incomplete_upload (
    id TEXT PRIMARY KEY,
    owner_id TEXT NOT NULL REFERENCES owner (id) ON DELETE RESTRICT,
    admin_debug_enabled BOOLEAN NOT NULL,
    blob_id TEXT REFERENCES blob (id) ON DELETE RESTRICT, -- Will be deleted when the upload completes (succeeds/fails)
    name TEXT NOT NULL,
    description TEXT NOT NULL,    
    holdings_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL,
    ran_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    failure_code failure_code,
    failure_message TEXT
);

CREATE TABLE portfolio (
    id TEXT PRIMARY KEY,
    owner_id TEXT NOT NULL REFERENCES owner (id) ON DELETE RESTRICT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    holdings_date TIMESTAMPTZ,
    blob_id TEXT NOT NULL REFERENCES blob (id) ON DELETE RESTRICT,
    admin_debug_enabled BOOLEAN NOT NULL,
    -- These are up for debate, but their basic goal is to help the user identify and differentiate between portfolios.
    number_of_rows INT8
);

CREATE TABLE portfolio_group (
    id TEXT PRIMARY KEY,
    owner_id TEXT NOT NULL REFERENCES owner (id) ON DELETE RESTRICT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL 
);

CREATE TABLE portfolio_group_membership (
    portfolio_id TEXT REFERENCES portfolio (id) ON DELETE RESTRICT,
    portfolio_group_id TEXT REFERENCES portfolio_group (id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL
);
    
CREATE TABLE portfolio_initiative_membership (
    portfolio_id TEXT NOT NULL REFERENCES portfolio (id) ON DELETE RESTRICT,
    initiative_id TEXT NOT NULL REFERENCES initiative (id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL,
    added_by_user_id TEXT REFERENCES pacta_user (id) ON DELETE RESTRICT
);

-- Stores the portfolios that were members to this snapshot at the time that
-- something was run, since portfolios can be deleted etc. Since we store the
-- PACTA version used alongside artifacts, and we treat portfolios as immutable,
-- this allows us to know if a report or audit is "stale" i.e. includes a
-- different set of portfolios than would be used for a current run.
CREATE TABLE portfolio_snapshot (
    id TEXT PRIMARY KEY,
    portfolio_id TEXT REFERENCES portfolio (id) ON DELETE RESTRICT,
    portfolio_group_id TEXT REFERENCES portfolio_group (id) ON DELETE RESTRICT,
    initiative_id TEXT REFERENCES initiative (id) ON DELETE RESTRICT,
    portfolio_ids TEXT[] -- Might include portfolios that have been deleted.
);
ALTER TABLE portfolio_snapshot ADD CONSTRAINT snapshot_is_well_formed CHECK (NUM_NONNULLS(portfolio_id, portfolio_group_id, initiative_id) = 1);

CREATE TYPE analysis_type AS ENUM (
    'audit',
    'report'
);

CREATE TABLE analysis (
    id TEXT PRIMARY KEY,
    analysis_type analysis_type NOT NULL,
    owner_id TEXT NOT NULL REFERENCES owner (id) ON DELETE RESTRICT,
    pacta_version_id TEXT NOT NULL REFERENCES pacta_version (id) ON DELETE RESTRICT,
    portfolio_snapshot_id TEXT NOT NULL REFERENCES portfolio_snapshot (id) ON DELETE RESTRICT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    ran_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    failure_code failure_code,
    failure_message TEXT
);
    
CREATE TABLE analysis_artifact (
    id TEXT PRIMARY KEY,
    analysis_id TEXT NOT NULL REFERENCES analysis (id) ON DELETE RESTRICT,
    blob_id TEXT NOT NULL REFERENCES blob (id) ON DELETE RESTRICT,
    admin_debug_enabled BOOLEAN NOT NULL,
    shared_to_public BOOLEAN NOT NULL
);

CREATE TYPE audit_log_action AS ENUM (
    'CREATE',
    'UPDATE',
    'DELETE',
    'ADD_TO',
    'REMOVE_FROM',
    'ENABLE_ADMIN_DEBUG',
    'DISABLE_ADMIN_DEBUG',
    'DOWNLOAD',
    'ENABLE_SHARING',
    'DISABLE_SHARING'
);

CREATE TYPE audit_log_actor_type AS ENUM (
    'USER',
    'ADMIN',
    'SUPER_ADMIN',
    'SYSTEM'
);

CREATE TYPE audit_log_target_type AS ENUM (
    'USER',
    'PORTFOLIO',
    'PORTFOLIO_GROUP',
    'INITIATIVE',
    'PACTA_VERSION',
    'ANALYSIS',
    'INCOMPLETE_UPLOAD'
);

CREATE TABLE audit_log (
    time TIMESTAMPTZ NOT NULL,
    actor_type audit_log_actor_type NOT NULL,
    actor_id TEXT NOT NULL,
    actor_owner_id TEXT REFERENCES owner (id) ON DELETE RESTRICT,
    action audit_log_action NOT NULL,
    primary_target_type audit_log_target_type NOT NULL,
    primary_target_id TEXT NOT NULL,
    primary_target_owner_id TEXT REFERENCES owner (id) ON DELETE RESTRICT,
    secondary_target_type audit_log_target_type NOT NULL,
    secondary_target_id TEXT NOT NULL,
    secondary_target_owner_id TEXT REFERENCES owner (id) ON DELETE RESTRICT
);

COMMIT;