-- Listings context: its own schema, its own writer, no foreign key leaving it.
-- Identity is the provider's listing id scoped by tenant + channel account;
-- Listings mints nothing. Every fact column is a state/value/reason triple
-- because "ML said nothing" and "ML said zero" are different facts (0097 is
-- the ratified precedent).
CREATE SCHEMA IF NOT EXISTS listings;

CREATE TABLE IF NOT EXISTS listings.listings (
    tenant_id           text        NOT NULL,
    channel             text        NOT NULL,
    account_external_id text        NOT NULL,
    listing_id          text        NOT NULL,
    version             integer     NOT NULL,
    title_state         text        NOT NULL,
    title_value         text        NULL,
    title_reason        text        NULL,
    status_state        text        NOT NULL,
    status_value        text        NULL,
    status_reason       text        NULL,
    listing_type_state  text        NOT NULL,
    listing_type_value  text        NULL,
    listing_type_reason text        NULL,
    price_state         text        NOT NULL,
    price_amount        text        NULL,
    price_currency      text        NULL,
    price_reason        text        NULL,
    available_qty_state  text       NOT NULL,
    available_qty_value  integer    NULL,
    available_qty_reason text       NULL,
    seller_sku_state    text        NOT NULL,
    seller_sku_value    text        NULL,
    seller_sku_reason   text        NULL,
    gtin_state          text        NOT NULL,
    gtin_value          text        NULL,
    gtin_reason         text        NULL,
    last_payload_hash   text        NOT NULL,
    created_at          timestamptz NOT NULL DEFAULT now(),
    updated_at          timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT listings_pkey PRIMARY KEY (tenant_id, channel, account_external_id, listing_id),
    CONSTRAINT listings_version_positive CHECK (version >= 1),
    CONSTRAINT listings_title_state CHECK (title_state IN ('known','estimated','unknown','not_applicable')),
    CONSTRAINT listings_status_state CHECK (status_state IN ('known','estimated','unknown','not_applicable')),
    CONSTRAINT listings_listing_type_state CHECK (listing_type_state IN ('known','estimated','unknown','not_applicable')),
    CONSTRAINT listings_price_state CHECK (price_state IN ('known','estimated','unknown','not_applicable')),
    CONSTRAINT listings_available_qty_state CHECK (available_qty_state IN ('known','estimated','unknown','not_applicable')),
    CONSTRAINT listings_seller_sku_state CHECK (seller_sku_state IN ('known','estimated','unknown','not_applicable')),
    CONSTRAINT listings_gtin_state CHECK (gtin_state IN ('known','estimated','unknown','not_applicable')),
    -- known carries a value; unknown carries a reason and no value (0097:26-31).
    CONSTRAINT listings_title_consistent CHECK (
        (title_state = 'known' AND title_value IS NOT NULL)
     OR (title_state = 'estimated' AND title_value IS NOT NULL AND title_reason IS NOT NULL)
     OR (title_state IN ('unknown','not_applicable') AND title_value IS NULL AND title_reason IS NOT NULL)),
    CONSTRAINT listings_price_consistent CHECK (
        (price_state = 'known' AND price_amount IS NOT NULL AND price_currency IS NOT NULL)
     OR (price_state = 'estimated' AND price_amount IS NOT NULL AND price_currency IS NOT NULL AND price_reason IS NOT NULL)
     OR (price_state IN ('unknown','not_applicable') AND price_amount IS NULL AND price_currency IS NULL AND price_reason IS NOT NULL))
);

CREATE TABLE IF NOT EXISTS listings.listing_variations (
    tenant_id           text    NOT NULL,
    channel             text    NOT NULL,
    account_external_id text    NOT NULL,
    listing_id          text    NOT NULL,
    variation_id        text    NOT NULL,
    price_state         text    NOT NULL,
    price_amount        text    NULL,
    price_currency      text    NULL,
    price_reason        text    NULL,
    available_qty_state  text   NOT NULL,
    available_qty_value  integer NULL,
    available_qty_reason text   NULL,
    seller_sku_state    text    NOT NULL,
    seller_sku_value    text    NULL,
    seller_sku_reason   text    NULL,
    gtin_state          text    NOT NULL,
    gtin_value          text    NULL,
    gtin_reason         text    NULL,
    CONSTRAINT listing_variations_pkey
        PRIMARY KEY (tenant_id, channel, account_external_id, listing_id, variation_id),
    CONSTRAINT listing_variations_listing_fkey
        FOREIGN KEY (tenant_id, channel, account_external_id, listing_id)
        REFERENCES listings.listings (tenant_id, channel, account_external_id, listing_id)
        ON DELETE CASCADE
);

-- One row per distinct payload the channel ever showed us: the real bytes,
-- kept for reconciliation and for the linking leg (§15.3).
CREATE TABLE IF NOT EXISTS listings.source_observations (
    tenant_id           text        NOT NULL,
    channel             text        NOT NULL,
    account_external_id text        NOT NULL,
    listing_id          text        NOT NULL,
    payload_hash        text        NOT NULL,
    payload             jsonb       NOT NULL,
    observed_at         timestamptz NOT NULL,
    recorded_at         timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT source_observations_pkey
        PRIMARY KEY (tenant_id, channel, account_external_id, listing_id, payload_hash),
    CONSTRAINT source_observations_listing_fkey
        FOREIGN KEY (tenant_id, channel, account_external_id, listing_id)
        REFERENCES listings.listings (tenant_id, channel, account_external_id, listing_id)
        ON DELETE CASCADE
);

ALTER TABLE listings.listings            ENABLE ROW LEVEL SECURITY;
ALTER TABLE listings.listing_variations  ENABLE ROW LEVEL SECURITY;
ALTER TABLE listings.source_observations ENABLE ROW LEVEL SECURITY;
ALTER TABLE listings.listings            FORCE ROW LEVEL SECURITY;
ALTER TABLE listings.listing_variations  FORCE ROW LEVEL SECURITY;
ALTER TABLE listings.source_observations FORCE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS tenant_isolation ON listings.listings;
CREATE POLICY tenant_isolation ON listings.listings
    USING (tenant_id = current_setting('app.tenant_id', true))
    WITH CHECK (tenant_id = current_setting('app.tenant_id', true));
DROP POLICY IF EXISTS tenant_isolation ON listings.listing_variations;
CREATE POLICY tenant_isolation ON listings.listing_variations
    USING (tenant_id = current_setting('app.tenant_id', true))
    WITH CHECK (tenant_id = current_setting('app.tenant_id', true));
DROP POLICY IF EXISTS tenant_isolation ON listings.source_observations;
CREATE POLICY tenant_isolation ON listings.source_observations
    USING (tenant_id = current_setting('app.tenant_id', true))
    WITH CHECK (tenant_id = current_setting('app.tenant_id', true));
