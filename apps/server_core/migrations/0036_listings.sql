-- No provider variation is normalized to the literal '-' at the application boundary.
CREATE TABLE IF NOT EXISTS listings (
    tenant_id TEXT NOT NULL,
    installation_id TEXT NOT NULL,
    provider TEXT NOT NULL,
    provider_listing_id TEXT NOT NULL,
    variation_id TEXT NOT NULL,
    title TEXT NOT NULL,
    listing_type_code TEXT,
    status TEXT NOT NULL,
    price_amount NUMERIC,
    price_currency TEXT,
    published_quantity INT,
    sync_state TEXT NOT NULL,
    sync_error JSONB,
    quality_score INT,
    sales_30d INT,
    fetched_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, installation_id, provider_listing_id, variation_id),
    CONSTRAINT listings_status_check
        CHECK (status IN ('active', 'paused', 'closed', 'unknown')),
    CONSTRAINT listings_sync_state_check
        CHECK (sync_state IN ('synced', 'error', 'stale', 'queued', 'syncing', 'paused_sync')),
    CONSTRAINT listings_quality_score_check
        CHECK (quality_score IS NULL OR quality_score BETWEEN 0 AND 100),
    CONSTRAINT listings_price_pair_check
        CHECK ((price_amount IS NULL AND price_currency IS NULL)
            OR (price_amount IS NOT NULL AND price_currency IS NOT NULL))
);

CREATE INDEX IF NOT EXISTS idx_listings_f02_title_key
    ON listings (tenant_id, installation_id, title, provider_listing_id, variation_id);

CREATE TABLE IF NOT EXISTS listing_sync_events (
    tenant_id TEXT NOT NULL,
    installation_id TEXT NOT NULL,
    provider_listing_id TEXT NOT NULL,
    variation_id TEXT NOT NULL,
    event_id TEXT NOT NULL,
    at TIMESTAMPTZ NOT NULL,
    kind TEXT NOT NULL,
    message_pt TEXT NOT NULL,
    PRIMARY KEY (tenant_id, installation_id, provider_listing_id, variation_id, event_id),
    CONSTRAINT listing_sync_events_kind_check
        CHECK (kind IN ('synced', 'sync_error', 'closed', 'paused', 'refreshed'))
);

CREATE INDEX IF NOT EXISTS idx_listing_sync_events_timeline
    ON listing_sync_events (
        tenant_id, installation_id, provider_listing_id, variation_id, at DESC, event_id DESC
    );
