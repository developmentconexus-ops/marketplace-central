-- 0087: divergences — IC-02 shared divergence ledger (M-02 sync-core-seam,
-- F-01-core-ddl).
--
-- Detection always happens at ingest (ADR-10): a daily stock producer (M-05)
-- and a 5-minute order producer (M-06) each Evaluate() into the SAME table
-- instead of inventing incompatible shapes (append-events vs. one-open-row).
-- No producer writes into this table yet.
--
-- Row "open" = resolved_at IS NULL; there is no separate status column, the
-- NULL is the state (IC-02 Enums And Statuses). Convergence sets resolved_at
-- and the row stays as history — never a physical delete. A later flap opens
-- a brand new row (new detected_at); the partial unique index below permits
-- exactly that, one open row per (entity, kind) at a time.
--
-- entity_id formats mirror IC-01's subject_id (IC-02 Fields):
--   estoque: <provider_listing_id>            (no variation)
--            <provider_listing_id>:<variation_id>  (with variation)
--   tarifa:  <provider_order_id>:<provider_item_id>
--
-- No FK: entity_id references provider-side/mirror entities, not a row this
-- schema owns.
--
-- R-5: both observed_at timestamps are always dated (NOT NULL) so staleness
-- is distinguishable from real disagreement.
CREATE TABLE IF NOT EXISTS divergences (
    id                     bigserial      PRIMARY KEY,
    tenant_id              text           NOT NULL,
    provider               text           NOT NULL,
    entity_type            text           NOT NULL,
    entity_id              text           NOT NULL,
    kind                   text           NOT NULL,
    expected_value         numeric(14,4)  NOT NULL,
    observed_value         numeric(14,4)  NOT NULL,
    expected_source        text           NOT NULL,
    observed_source        text           NOT NULL,
    expected_observed_at   timestamptz    NOT NULL,
    observed_observed_at   timestamptz    NOT NULL,
    detected_at            timestamptz    NOT NULL,
    last_evaluated_at      timestamptz    NOT NULL,
    resolved_at            timestamptz,
    created_at             timestamptz    NOT NULL DEFAULT now(),
    updated_at             timestamptz    NOT NULL DEFAULT now(),
    CONSTRAINT divergences_entity_type_check
        CHECK (entity_type IN ('listing', 'order_line')),
    CONSTRAINT divergences_kind_check
        CHECK (kind IN ('estoque', 'tarifa'))
);

-- The load-bearing constraint (IC-02 Persistence Expectations): at most one
-- OPEN row per (tenant, provider, entity_type, entity_id, kind). Postgres has
-- no partial UNIQUE table constraint, so this is expressed as a partial
-- unique index instead.
CREATE UNIQUE INDEX IF NOT EXISTS divergences_one_open_row
    ON divergences (tenant_id, provider, entity_type, entity_id, kind)
    WHERE resolved_at IS NULL;

-- Cheap badge count: how many open divergences a tenant has, without scanning
-- resolved history.
CREATE INDEX IF NOT EXISTS idx_divergences_open
    ON divergences (tenant_id, resolved_at)
    WHERE resolved_at IS NULL;
