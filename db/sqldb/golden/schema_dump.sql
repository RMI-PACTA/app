--
-- PostgreSQL database dump
--

-- Dumped from database version 14.9 (Debian 14.9-1.pgdg120+1)
-- Dumped by pg_dump version 14.9 (Debian 14.9-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: analysis_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.analysis_type AS ENUM (
    'audit',
    'report'
);


ALTER TYPE public.analysis_type OWNER TO postgres;

--
-- Name: audit_log_action; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.audit_log_action AS ENUM (
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


ALTER TYPE public.audit_log_action OWNER TO postgres;

--
-- Name: audit_log_actor_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.audit_log_actor_type AS ENUM (
    'USER',
    'ADMIN',
    'SUPER_ADMIN',
    'SYSTEM'
);


ALTER TYPE public.audit_log_actor_type OWNER TO postgres;

--
-- Name: audit_log_target_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.audit_log_target_type AS ENUM (
    'USER',
    'PORTFOLIO',
    'PORTFOLIO_GROUP',
    'INITIATIVE',
    'PACTA_VERSION',
    'ANALYSIS',
    'INCOMPLETE_UPLOAD'
);


ALTER TYPE public.audit_log_target_type OWNER TO postgres;

--
-- Name: authn_mechanism; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.authn_mechanism AS ENUM (
    'EMAIL_AND_PASS'
);


ALTER TYPE public.authn_mechanism OWNER TO postgres;

--
-- Name: failure_code; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.failure_code AS ENUM (
    'UNKNOWN'
);


ALTER TYPE public.failure_code OWNER TO postgres;

--
-- Name: file_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.file_type AS ENUM (
    'csv',
    'yaml',
    'zip',
    'html'
);


ALTER TYPE public.file_type OWNER TO postgres;

--
-- Name: language; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.language AS ENUM (
    'en',
    'de',
    'fr',
    'es'
);


ALTER TYPE public.language OWNER TO postgres;

--
-- Name: track_applied_migration(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.track_applied_migration() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE _current_version integer;
BEGIN
    SELECT COALESCE(MAX(version),0) FROM schema_migrations_history INTO _current_version;
    IF new.dirty = 'f' AND new.version > _current_version THEN
        INSERT INTO schema_migrations_history(version) VALUES (new.version);
    END IF;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.track_applied_migration() OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: analysis; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.analysis (
    id text NOT NULL,
    analysis_type public.analysis_type NOT NULL,
    owner_id text NOT NULL,
    pacta_version_id text NOT NULL,
    portfolio_snapshot_id text NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    ran_at timestamp with time zone,
    completed_at timestamp with time zone,
    failure_code public.failure_code,
    failure_message text
);


ALTER TABLE public.analysis OWNER TO postgres;

--
-- Name: analysis_artifact; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.analysis_artifact (
    id text NOT NULL,
    analysis_id text NOT NULL,
    blob_id text NOT NULL,
    admin_debug_enabled boolean NOT NULL,
    shared_to_public boolean NOT NULL
);


ALTER TABLE public.analysis_artifact OWNER TO postgres;

--
-- Name: audit_log; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.audit_log (
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    actor_type public.audit_log_actor_type NOT NULL,
    actor_id text NOT NULL,
    actor_owner_id text NOT NULL,
    action public.audit_log_action NOT NULL,
    primary_target_type public.audit_log_target_type NOT NULL,
    primary_target_id text NOT NULL,
    primary_target_owner_id text NOT NULL,
    secondary_target_type public.audit_log_target_type,
    secondary_target_id text NOT NULL,
    secondary_target_owner_id text,
    id text NOT NULL
);


ALTER TABLE public.audit_log OWNER TO postgres;

--
-- Name: blob; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.blob (
    id text NOT NULL,
    blob_uri text NOT NULL,
    file_type public.file_type NOT NULL,
    file_name text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.blob OWNER TO postgres;

--
-- Name: incomplete_upload; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.incomplete_upload (
    id text NOT NULL,
    owner_id text NOT NULL,
    admin_debug_enabled boolean NOT NULL,
    blob_id text,
    name text NOT NULL,
    description text NOT NULL,
    holdings_date timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    ran_at timestamp with time zone,
    completed_at timestamp with time zone,
    failure_code public.failure_code,
    failure_message text
);


ALTER TABLE public.incomplete_upload OWNER TO postgres;

--
-- Name: initiative; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.initiative (
    id text NOT NULL,
    name text NOT NULL,
    affiliation text NOT NULL,
    public_description text NOT NULL,
    internal_description text NOT NULL,
    requires_invitation_to_join boolean NOT NULL,
    is_accepting_new_members boolean NOT NULL,
    is_accepting_new_portfolios boolean NOT NULL,
    pacta_version_id text NOT NULL,
    language public.language NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.initiative OWNER TO postgres;

--
-- Name: initiative_invitation; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.initiative_invitation (
    id text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    used_at timestamp with time zone,
    initiative_id text NOT NULL,
    used_by_user_id text
);


ALTER TABLE public.initiative_invitation OWNER TO postgres;

--
-- Name: initiative_user_relationship; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.initiative_user_relationship (
    user_id text NOT NULL,
    initiative_id text NOT NULL,
    manager boolean NOT NULL,
    member boolean NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.initiative_user_relationship OWNER TO postgres;

--
-- Name: owner; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.owner (
    id text NOT NULL,
    user_id text,
    initiative_id text,
    CONSTRAINT owner_is_always_well_defined CHECK ((num_nonnulls(user_id, initiative_id) = 1))
);


ALTER TABLE public.owner OWNER TO postgres;

--
-- Name: pacta_user; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pacta_user (
    id text NOT NULL,
    authn_mechanism public.authn_mechanism NOT NULL,
    authn_id text NOT NULL,
    entered_email text NOT NULL,
    canonical_email text NOT NULL,
    admin boolean NOT NULL,
    super_admin boolean NOT NULL,
    name text NOT NULL,
    preferred_language public.language,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.pacta_user OWNER TO postgres;

--
-- Name: pacta_version; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pacta_version (
    id text NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    digest text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    is_default boolean,
    CONSTRAINT is_default_is_true_or_null CHECK (is_default)
);


ALTER TABLE public.pacta_version OWNER TO postgres;

--
-- Name: portfolio; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.portfolio (
    id text NOT NULL,
    owner_id text NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    holdings_date timestamp with time zone,
    blob_id text NOT NULL,
    admin_debug_enabled boolean NOT NULL,
    number_of_rows integer
);


ALTER TABLE public.portfolio OWNER TO postgres;

--
-- Name: portfolio_group; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.portfolio_group (
    id text NOT NULL,
    owner_id text NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.portfolio_group OWNER TO postgres;

--
-- Name: portfolio_group_membership; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.portfolio_group_membership (
    portfolio_id text,
    portfolio_group_id text,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.portfolio_group_membership OWNER TO postgres;

--
-- Name: portfolio_initiative_membership; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.portfolio_initiative_membership (
    portfolio_id text NOT NULL,
    initiative_id text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    added_by_user_id text
);


ALTER TABLE public.portfolio_initiative_membership OWNER TO postgres;

--
-- Name: portfolio_snapshot; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.portfolio_snapshot (
    id text NOT NULL,
    portfolio_id text,
    portfolio_group_id text,
    initiative_id text,
    portfolio_ids text[],
    CONSTRAINT snapshot_is_well_formed CHECK ((num_nonnulls(portfolio_id, portfolio_group_id, initiative_id) = 1))
);


ALTER TABLE public.portfolio_snapshot OWNER TO postgres;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO postgres;

--
-- Name: schema_migrations_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migrations_history (
    id integer NOT NULL,
    version bigint NOT NULL,
    applied_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.schema_migrations_history OWNER TO postgres;

--
-- Name: schema_migrations_history_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.schema_migrations_history_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.schema_migrations_history_id_seq OWNER TO postgres;

--
-- Name: schema_migrations_history_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.schema_migrations_history_id_seq OWNED BY public.schema_migrations_history.id;


--
-- Name: schema_migrations_history id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schema_migrations_history ALTER COLUMN id SET DEFAULT nextval('public.schema_migrations_history_id_seq'::regclass);


--
-- Name: analysis_artifact analysis_artifact_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.analysis_artifact
    ADD CONSTRAINT analysis_artifact_pkey PRIMARY KEY (id);


--
-- Name: analysis analysis_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.analysis
    ADD CONSTRAINT analysis_pkey PRIMARY KEY (id);


--
-- Name: audit_log audit_log_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.audit_log
    ADD CONSTRAINT audit_log_pkey PRIMARY KEY (id);


--
-- Name: blob blob_blob_uri_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.blob
    ADD CONSTRAINT blob_blob_uri_key UNIQUE (blob_uri);


--
-- Name: blob blob_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.blob
    ADD CONSTRAINT blob_pkey PRIMARY KEY (id);


--
-- Name: incomplete_upload incomplete_upload_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.incomplete_upload
    ADD CONSTRAINT incomplete_upload_pkey PRIMARY KEY (id);


--
-- Name: initiative_invitation initiative_invitation_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.initiative_invitation
    ADD CONSTRAINT initiative_invitation_pkey PRIMARY KEY (id);


--
-- Name: initiative initiative_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.initiative
    ADD CONSTRAINT initiative_pkey PRIMARY KEY (id);


--
-- Name: initiative_user_relationship initiative_user_relationship_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.initiative_user_relationship
    ADD CONSTRAINT initiative_user_relationship_pkey PRIMARY KEY (user_id, initiative_id);


--
-- Name: pacta_version is_default_only_1_true; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pacta_version
    ADD CONSTRAINT is_default_only_1_true UNIQUE (is_default);


--
-- Name: owner owner_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.owner
    ADD CONSTRAINT owner_pkey PRIMARY KEY (id);


--
-- Name: pacta_user pacta_user_authn_mechanism_authn_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pacta_user
    ADD CONSTRAINT pacta_user_authn_mechanism_authn_id_key UNIQUE (authn_mechanism, authn_id);


--
-- Name: pacta_user pacta_user_canonical_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pacta_user
    ADD CONSTRAINT pacta_user_canonical_email_key UNIQUE (canonical_email);


--
-- Name: pacta_user pacta_user_entered_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pacta_user
    ADD CONSTRAINT pacta_user_entered_email_key UNIQUE (entered_email);


--
-- Name: pacta_user pacta_user_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pacta_user
    ADD CONSTRAINT pacta_user_pkey PRIMARY KEY (id);


--
-- Name: pacta_version pacta_version_digest_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pacta_version
    ADD CONSTRAINT pacta_version_digest_key UNIQUE (digest);


--
-- Name: pacta_version pacta_version_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pacta_version
    ADD CONSTRAINT pacta_version_pkey PRIMARY KEY (id);


--
-- Name: portfolio_group portfolio_group_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.portfolio_group
    ADD CONSTRAINT portfolio_group_pkey PRIMARY KEY (id);


--
-- Name: portfolio portfolio_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.portfolio
    ADD CONSTRAINT portfolio_pkey PRIMARY KEY (id);


--
-- Name: portfolio_snapshot portfolio_snapshot_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.portfolio_snapshot
    ADD CONSTRAINT portfolio_snapshot_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations_history schema_migrations_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schema_migrations_history
    ADD CONSTRAINT schema_migrations_history_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: owner_by_initiative_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX owner_by_initiative_id ON public.owner USING btree (initiative_id);


--
-- Name: owner_by_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX owner_by_user_id ON public.owner USING btree (user_id);


--
-- Name: schema_migrations track_applied_migrations; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER track_applied_migrations AFTER INSERT ON public.schema_migrations FOR EACH ROW EXECUTE FUNCTION public.track_applied_migration();


--
-- Name: analysis_artifact analysis_artifact_analysis_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.analysis_artifact
    ADD CONSTRAINT analysis_artifact_analysis_id_fkey FOREIGN KEY (analysis_id) REFERENCES public.analysis(id) ON DELETE RESTRICT;


--
-- Name: analysis_artifact analysis_artifact_blob_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.analysis_artifact
    ADD CONSTRAINT analysis_artifact_blob_id_fkey FOREIGN KEY (blob_id) REFERENCES public.blob(id) ON DELETE RESTRICT;


--
-- Name: analysis analysis_owner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.analysis
    ADD CONSTRAINT analysis_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public.owner(id) ON DELETE RESTRICT;


--
-- Name: analysis analysis_pacta_version_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.analysis
    ADD CONSTRAINT analysis_pacta_version_id_fkey FOREIGN KEY (pacta_version_id) REFERENCES public.pacta_version(id) ON DELETE RESTRICT;


--
-- Name: analysis analysis_portfolio_snapshot_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.analysis
    ADD CONSTRAINT analysis_portfolio_snapshot_id_fkey FOREIGN KEY (portfolio_snapshot_id) REFERENCES public.portfolio_snapshot(id) ON DELETE RESTRICT;


--
-- Name: incomplete_upload incomplete_upload_blob_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.incomplete_upload
    ADD CONSTRAINT incomplete_upload_blob_id_fkey FOREIGN KEY (blob_id) REFERENCES public.blob(id) ON DELETE RESTRICT;


--
-- Name: incomplete_upload incomplete_upload_owner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.incomplete_upload
    ADD CONSTRAINT incomplete_upload_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public.owner(id) ON DELETE RESTRICT;


--
-- Name: initiative_invitation initiative_invitation_initiative_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.initiative_invitation
    ADD CONSTRAINT initiative_invitation_initiative_id_fkey FOREIGN KEY (initiative_id) REFERENCES public.initiative(id) ON DELETE RESTRICT;


--
-- Name: initiative_invitation initiative_invitation_used_by_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.initiative_invitation
    ADD CONSTRAINT initiative_invitation_used_by_user_id_fkey FOREIGN KEY (used_by_user_id) REFERENCES public.pacta_user(id) ON DELETE RESTRICT;


--
-- Name: initiative initiative_pacta_version_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.initiative
    ADD CONSTRAINT initiative_pacta_version_id_fkey FOREIGN KEY (pacta_version_id) REFERENCES public.pacta_version(id) ON DELETE RESTRICT;


--
-- Name: owner owner_initiative_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.owner
    ADD CONSTRAINT owner_initiative_id_fkey FOREIGN KEY (initiative_id) REFERENCES public.initiative(id) ON DELETE RESTRICT;


--
-- Name: owner owner_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.owner
    ADD CONSTRAINT owner_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.pacta_user(id) ON DELETE RESTRICT;


--
-- Name: portfolio portfolio_blob_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.portfolio
    ADD CONSTRAINT portfolio_blob_id_fkey FOREIGN KEY (blob_id) REFERENCES public.blob(id) ON DELETE RESTRICT;


--
-- Name: portfolio_group_membership portfolio_group_membership_portfolio_group_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.portfolio_group_membership
    ADD CONSTRAINT portfolio_group_membership_portfolio_group_id_fkey FOREIGN KEY (portfolio_group_id) REFERENCES public.portfolio_group(id) ON DELETE RESTRICT;


--
-- Name: portfolio_group_membership portfolio_group_membership_portfolio_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.portfolio_group_membership
    ADD CONSTRAINT portfolio_group_membership_portfolio_id_fkey FOREIGN KEY (portfolio_id) REFERENCES public.portfolio(id) ON DELETE RESTRICT;


--
-- Name: portfolio_group portfolio_group_owner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.portfolio_group
    ADD CONSTRAINT portfolio_group_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public.owner(id) ON DELETE RESTRICT;


--
-- Name: portfolio_initiative_membership portfolio_initiative_membership_added_by_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.portfolio_initiative_membership
    ADD CONSTRAINT portfolio_initiative_membership_added_by_user_id_fkey FOREIGN KEY (added_by_user_id) REFERENCES public.pacta_user(id) ON DELETE RESTRICT;


--
-- Name: portfolio_initiative_membership portfolio_initiative_membership_initiative_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.portfolio_initiative_membership
    ADD CONSTRAINT portfolio_initiative_membership_initiative_id_fkey FOREIGN KEY (initiative_id) REFERENCES public.initiative(id) ON DELETE RESTRICT;


--
-- Name: portfolio_initiative_membership portfolio_initiative_membership_portfolio_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.portfolio_initiative_membership
    ADD CONSTRAINT portfolio_initiative_membership_portfolio_id_fkey FOREIGN KEY (portfolio_id) REFERENCES public.portfolio(id) ON DELETE RESTRICT;


--
-- Name: portfolio portfolio_owner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.portfolio
    ADD CONSTRAINT portfolio_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public.owner(id) ON DELETE RESTRICT;


--
-- Name: portfolio_snapshot portfolio_snapshot_initiative_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.portfolio_snapshot
    ADD CONSTRAINT portfolio_snapshot_initiative_id_fkey FOREIGN KEY (initiative_id) REFERENCES public.initiative(id) ON DELETE RESTRICT;


--
-- Name: portfolio_snapshot portfolio_snapshot_portfolio_group_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.portfolio_snapshot
    ADD CONSTRAINT portfolio_snapshot_portfolio_group_id_fkey FOREIGN KEY (portfolio_group_id) REFERENCES public.portfolio_group(id) ON DELETE RESTRICT;


--
-- Name: portfolio_snapshot portfolio_snapshot_portfolio_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.portfolio_snapshot
    ADD CONSTRAINT portfolio_snapshot_portfolio_id_fkey FOREIGN KEY (portfolio_id) REFERENCES public.portfolio(id) ON DELETE RESTRICT;


--
-- PostgreSQL database dump complete
--

