-- The catalog schema's RLS was FORCEd in 0097 and was still decorative: the
-- application connects as the table owner, which is also a superuser with
-- rolbypassrls. Every policy was evaluated against a role that skips policies.
--
-- This creates the role the application is meant to use. It is NOSUPERUSER and
-- NOBYPASSRLS explicitly rather than by default, because the defect being fixed
-- is precisely an attribute nobody looked at.
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'mpc_app') THEN
        CREATE ROLE mpc_app NOLOGIN NOSUPERUSER NOBYPASSRLS NOCREATEDB NOCREATEROLE;
    END IF;
END
$$;

ALTER ROLE mpc_app NOSUPERUSER NOBYPASSRLS;

GRANT USAGE ON SCHEMA catalog TO mpc_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON catalog.products            TO mpc_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON catalog.product_identifiers TO mpc_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON catalog.source_product_keys TO mpc_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON catalog.source_observations TO mpc_app;

-- No GRANT on any other schema. A role that can read the legacy tables would
-- make the boundary this context exists to draw invisible again.
