CREATE TABLE IF NOT EXISTS market_validated_offers (
    tenant_id TEXT NOT NULL,
    catalog_product_id TEXT NOT NULL,
    seller_id TEXT NOT NULL,
    price_amount NUMERIC,
    price_currency TEXT,
    condition TEXT NOT NULL,
    match_status TEXT NOT NULL,
    fetched_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (tenant_id, catalog_product_id, seller_id, fetched_at),
    CONSTRAINT market_validated_offers_price_pair_check CHECK ((price_amount IS NULL AND price_currency IS NULL) OR (price_amount > 0 AND price_currency = 'BRL')),
    CONSTRAINT market_validated_offers_match_status_check CHECK (match_status IN ('ACCEPT', 'REVIEW', 'REJECT', 'NO_CANDIDATE'))
);

CREATE INDEX market_validated_offers_catalog_idx
    ON market_validated_offers (tenant_id, catalog_product_id, fetched_at DESC);
