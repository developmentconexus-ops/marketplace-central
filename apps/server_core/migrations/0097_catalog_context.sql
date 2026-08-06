-- Catalog context: its own schema, its own writer, no foreign key leaving it.
--
-- Every key leads with tenant_id and RLS is FORCEd, so a query that forgets the
-- tenant returns nothing instead of returning somebody else's catalogue.
CREATE SCHEMA IF NOT EXISTS catalog;

CREATE TABLE IF NOT EXISTS catalog.products (
    tenant_id           text        NOT NULL,
    product_id          text        NOT NULL,
    version             integer     NOT NULL,
    -- The knowledge state is a column because "the source said nothing" and
    -- "the source said it has no name" are different facts. A nullable
    -- description alone cannot tell them apart.
    description_state   text        NOT NULL,
    description_value   text        NULL,
    description_reason  text        NULL,
    last_payload_hash   text        NOT NULL,
    created_at          timestamptz NOT NULL DEFAULT now(),
    updated_at          timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT products_pkey PRIMARY KEY (tenant_id, product_id),
    CONSTRAINT products_version_positive CHECK (version >= 1),
    CONSTRAINT products_state_known CHECK (
        description_state IN ('known', 'estimated', 'unknown', 'not_applicable')
    ),
    -- known carries a value; unknown carries a reason and no value.
    CONSTRAINT products_state_consistent CHECK (
        (description_state = 'known'      AND description_value IS NOT NULL)
     OR (description_state = 'estimated'  AND description_value IS NOT NULL AND description_reason IS NOT NULL)
     OR (description_state IN ('unknown', 'not_applicable')
                                          AND description_value IS NULL     AND description_reason IS NOT NULL)
    )
);

CREATE TABLE IF NOT EXISTS catalog.product_identifiers (
    tenant_id  text NOT NULL,
    product_id text NOT NULL,
    kind       text NOT NULL,
    value      text NOT NULL,
    CONSTRAINT product_identifiers_pkey PRIMARY KEY (tenant_id, product_id, kind, value),
    CONSTRAINT product_identifiers_product_fkey
        FOREIGN KEY (tenant_id, product_id)
        REFERENCES catalog.products (tenant_id, product_id) ON DELETE CASCADE
);

-- Deliberately NOT unique on (tenant_id, kind, value): two products sharing a
-- bad EAN is real master data, and a unique index would turn it into an insert
-- failure at 03:00 instead of a conflict a human can look at.
CREATE INDEX IF NOT EXISTS product_identifiers_lookup
    ON catalog.product_identifiers (tenant_id, kind, value);

CREATE TABLE IF NOT EXISTS catalog.source_product_keys (
    tenant_id       text NOT NULL,
    source_system   text NOT NULL,
    source_instance text NOT NULL,
    object_kind     text NOT NULL,
    external_key    text NOT NULL,
    product_id      text NOT NULL,
    CONSTRAINT source_product_keys_pkey
        PRIMARY KEY (tenant_id, source_system, source_instance, object_kind, external_key),
    CONSTRAINT source_product_keys_product_fkey
        FOREIGN KEY (tenant_id, product_id)
        REFERENCES catalog.products (tenant_id, product_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS catalog.source_observations (
    tenant_id       text        NOT NULL,
    product_id      text        NOT NULL,
    payload_hash    text        NOT NULL,
    source_system   text        NOT NULL,
    object_kind     text        NOT NULL,
    external_key    text        NOT NULL,
    observed_at     timestamptz NOT NULL,
    recorded_at     timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT source_observations_pkey PRIMARY KEY (tenant_id, product_id, payload_hash),
    CONSTRAINT source_observations_product_fkey
        FOREIGN KEY (tenant_id, product_id)
        REFERENCES catalog.products (tenant_id, product_id) ON DELETE CASCADE
);

ALTER TABLE catalog.products             ENABLE ROW LEVEL SECURITY;
ALTER TABLE catalog.product_identifiers  ENABLE ROW LEVEL SECURITY;
ALTER TABLE catalog.source_product_keys  ENABLE ROW LEVEL SECURITY;
ALTER TABLE catalog.source_observations  ENABLE ROW LEVEL SECURITY;

-- FORCE, because without it the table owner silently bypasses every policy and
-- the whole mechanism is decorative in exactly the environment we test in.
ALTER TABLE catalog.products             FORCE ROW LEVEL SECURITY;
ALTER TABLE catalog.product_identifiers  FORCE ROW LEVEL SECURITY;
ALTER TABLE catalog.source_product_keys  FORCE ROW LEVEL SECURITY;
ALTER TABLE catalog.source_observations  FORCE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS tenant_isolation ON catalog.products;
CREATE POLICY tenant_isolation ON catalog.products
    USING (tenant_id = current_setting('app.tenant_id', true))
    WITH CHECK (tenant_id = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS tenant_isolation ON catalog.product_identifiers;
CREATE POLICY tenant_isolation ON catalog.product_identifiers
    USING (tenant_id = current_setting('app.tenant_id', true))
    WITH CHECK (tenant_id = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS tenant_isolation ON catalog.source_product_keys;
CREATE POLICY tenant_isolation ON catalog.source_product_keys
    USING (tenant_id = current_setting('app.tenant_id', true))
    WITH CHECK (tenant_id = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS tenant_isolation ON catalog.source_observations;
CREATE POLICY tenant_isolation ON catalog.source_observations
    USING (tenant_id = current_setting('app.tenant_id', true))
    WITH CHECK (tenant_id = current_setting('app.tenant_id', true));
