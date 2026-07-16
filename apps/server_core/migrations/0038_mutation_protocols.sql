-- source_as_of remains unknown in draft and is required by domain rules from previewed onward.
CREATE TABLE IF NOT EXISTS mutation_protocols (
    tenant_id TEXT NOT NULL,
    protocol_id TEXT NOT NULL,
    installation_id TEXT NOT NULL,
    type TEXT NOT NULL,
    state TEXT NOT NULL,
    actor TEXT NOT NULL,
    intent JSONB NOT NULL,
    selection JSONB NOT NULL,
    totals JSONB NOT NULL,
    source_as_of TIMESTAMPTZ,
    retried_from TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    previewed_at TIMESTAMPTZ,
    approved_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    PRIMARY KEY (tenant_id, protocol_id),
    CONSTRAINT mutation_protocols_type_check
        CHECK (type IN ('price_update', 'stock_correct', 'link_apply', 'listing_pause', 'listing_resync', 'listing_edit', 'listing_create')),
    CONSTRAINT mutation_protocols_state_check
        CHECK (state IN ('draft', 'previewed', 'approved', 'applying', 'applied', 'partially_failed', 'failed_preserved', 'cancelled'))
);

CREATE INDEX IF NOT EXISTS mutation_protocols_claim_idx
    ON mutation_protocols (tenant_id, installation_id, created_at)
    WHERE state IN ('approved', 'applying');

CREATE TABLE IF NOT EXISTS mutation_items (
    tenant_id TEXT NOT NULL,
    protocol_id TEXT NOT NULL,
    seq INT NOT NULL,
    item_id TEXT NOT NULL,
    listing_id TEXT NOT NULL,
    idempotency_key TEXT NOT NULL,
    before JSONB,
    after JSONB NOT NULL,
    state TEXT NOT NULL,
    failure JSONB,
    applied_at TIMESTAMPTZ,
    PRIMARY KEY (tenant_id, protocol_id, seq),
    CONSTRAINT mutation_items_idempotency_key_unique UNIQUE (tenant_id, idempotency_key),
    CONSTRAINT mutation_items_state_check
        CHECK (state IN ('previewed', 'approved', 'applying', 'applied', 'failed', 'skipped'))
);
