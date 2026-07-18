CREATE TABLE IF NOT EXISTS market_match_decisions (
    tenant_id TEXT NOT NULL,
    product_id TEXT NOT NULL,
    catalog_product_id TEXT,
    match_status TEXT NOT NULL,
    anchors TEXT[] NOT NULL,
    contradictions TEXT[] NOT NULL,
    confidence_band TEXT NOT NULL,
    decided_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (tenant_id, product_id, decided_at),
    CONSTRAINT market_match_decisions_match_status_check CHECK (match_status IN ('ACCEPT', 'REVIEW', 'REJECT', 'NO_CANDIDATE')),
    CONSTRAINT market_match_decisions_confidence_band_check CHECK (confidence_band IN ('ALTA', 'MEDIA', 'BAIXA'))
);

CREATE INDEX market_match_decisions_product_idx
    ON market_match_decisions (tenant_id, product_id, decided_at DESC);
