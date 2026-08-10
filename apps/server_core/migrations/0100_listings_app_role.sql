-- Same reasoning as 0098: RLS FORCEd in 0099 is decorative until the
-- application role that cannot bypass it is granted the schema. mpc_app
-- already exists (0098); this only extends it to the listings schema.
GRANT USAGE ON SCHEMA listings TO mpc_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON listings.listings            TO mpc_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON listings.listing_variations  TO mpc_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON listings.source_observations TO mpc_app;
-- No GRANT on any other schema (0098:24-25).
