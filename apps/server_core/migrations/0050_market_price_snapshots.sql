CREATE TABLE IF NOT EXISTS market_price_snapshots (
    tenant_id TEXT NOT NULL,
    scope TEXT NOT NULL,
    ref_id TEXT NOT NULL,
    source TEXT NOT NULL,
    price_amount NUMERIC,
    price_currency TEXT,
    observed_at TIMESTAMPTZ NOT NULL,
    fetched_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL,
    failure_reason TEXT,
    request_id TEXT NOT NULL,
    PRIMARY KEY (tenant_id, scope, ref_id, source, fetched_at),
    CONSTRAINT market_price_snapshots_scope_check CHECK (scope IN ('listing', 'catalog_product')),
    CONSTRAINT market_price_snapshots_source_check CHECK (source IN ('ml_sale_price', 'ml_price_to_win', 'ml_catalog_offers')),
    CONSTRAINT market_price_snapshots_price_pair_check CHECK ((price_amount IS NULL AND price_currency IS NULL) OR (price_amount > 0 AND price_currency = 'BRL')),
    CONSTRAINT market_price_snapshots_status_check CHECK (status IN ('VALID', 'FAILED', 'EXPIRED')),
    CONSTRAINT market_price_snapshots_evidence_check CHECK ((status = 'FAILED' AND price_amount IS NULL AND failure_reason IS NOT NULL) OR (status IN ('VALID', 'EXPIRED') AND price_amount IS NOT NULL AND failure_reason IS NULL))
);

CREATE INDEX market_price_snapshots_latest_valid_idx
    ON market_price_snapshots (tenant_id, scope, ref_id, source, fetched_at DESC)
    WHERE status = 'VALID';
