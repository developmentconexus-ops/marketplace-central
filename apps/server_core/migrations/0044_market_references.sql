-- IC-04 contract-only append-only market references; no production seed data.
CREATE TABLE IF NOT EXISTS market_references (
    tenant_id TEXT NOT NULL,
    product_id TEXT NOT NULL,
    catalog_product_id TEXT,
    match_state TEXT NOT NULL,
    match_method TEXT,
    catalog_stats JSONB,
    evidence_state TEXT NOT NULL,
    captured_at TIMESTAMPTZ NOT NULL,
    source TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, product_id, captured_at),
    CONSTRAINT market_references_match_state_check CHECK (match_state IN ('accept', 'review', 'reject', 'no_candidate')),
    CONSTRAINT market_references_evidence_state_check CHECK (evidence_state IN ('observed', 'insufficient_market', 'no_price_evidence')),
    CONSTRAINT market_references_source_check CHECK (source IN ('official_api', 'vendor', 'manual'))
);
