-- Saved pricing scenarios (IC-04): named simulation inputs a tenant can store
-- and re-run. CRUD is tenant-scoped; listing is newest-first (see index). The
-- scenario inputs are held as jsonb (opaque to the schema); the engine
-- re-decomposes from them on demand.
CREATE TABLE IF NOT EXISTS pricing_scenarios (
    id         text NOT NULL,
    tenant_id  text NOT NULL,
    name       text NOT NULL,
    payload    jsonb NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, id)
);

CREATE INDEX pricing_scenarios_tenant_created_idx
    ON pricing_scenarios (tenant_id, created_at DESC);
