CREATE TABLE IF NOT EXISTS market_competitive_signals (
    tenant_id TEXT NOT NULL,
    listing_id TEXT NOT NULL,
    our_price_amount NUMERIC NOT NULL,
    our_price_currency TEXT NOT NULL,
    winner_price_amount NUMERIC,
    winner_price_currency TEXT,
    target_price_amount NUMERIC,
    target_price_currency TEXT,
    position_rank INTEGER,
    position_total INTEGER,
    fetched_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (tenant_id, listing_id, fetched_at),
    CONSTRAINT market_competitive_signals_our_price_check CHECK (our_price_amount > 0 AND our_price_currency = 'BRL'),
    CONSTRAINT market_competitive_signals_winner_price_pair_check CHECK ((winner_price_amount IS NULL AND winner_price_currency IS NULL) OR (winner_price_amount > 0 AND winner_price_currency = 'BRL')),
    CONSTRAINT market_competitive_signals_target_price_pair_check CHECK ((target_price_amount IS NULL AND target_price_currency IS NULL) OR (target_price_amount > 0 AND target_price_currency = 'BRL')),
    CONSTRAINT market_competitive_signals_position_pair_check CHECK ((position_rank IS NULL AND position_total IS NULL) OR (position_rank >= 1 AND position_total >= 1 AND position_rank <= position_total))
);

CREATE INDEX market_competitive_signals_listing_idx
    ON market_competitive_signals (tenant_id, listing_id, fetched_at DESC);

CREATE TABLE IF NOT EXISTS market_aggregates (
    tenant_id TEXT NOT NULL,
    product_id TEXT NOT NULL,
    median_amount NUMERIC,
    median_currency TEXT,
    min_valid_amount NUMERIC,
    min_valid_currency TEXT,
    n_offers INTEGER NOT NULL,
    n_sellers INTEGER NOT NULL,
    source TEXT NOT NULL,
    fetched_at TIMESTAMPTZ NOT NULL,
    computed_at TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL,
    PRIMARY KEY (tenant_id, product_id, source, computed_at),
    CONSTRAINT market_aggregates_median_pair_check CHECK ((median_amount IS NULL AND median_currency IS NULL) OR (median_amount > 0 AND median_currency = 'BRL')),
    CONSTRAINT market_aggregates_min_valid_pair_check CHECK ((min_valid_amount IS NULL AND min_valid_currency IS NULL) OR (min_valid_amount > 0 AND min_valid_currency = 'BRL')),
    CONSTRAINT market_aggregates_counts_check CHECK (n_offers >= 0 AND n_sellers >= 0),
    CONSTRAINT market_aggregates_source_check CHECK (source IN ('ml_sale_price', 'ml_price_to_win', 'ml_catalog_offers')),
    CONSTRAINT market_aggregates_status_check CHECK (status IN ('OK', 'INSUFFICIENT_MARKET', 'NO_PRICE_EVIDENCE'))
);

CREATE INDEX market_aggregates_product_idx
    ON market_aggregates (tenant_id, product_id, computed_at DESC);
