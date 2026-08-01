-- 0091: listing_variations — IC-07 (M-04 listings-ddl), new child table.
--
-- PK is the full parent tuple plus variation_id: (tenant_id, installation_id,
-- provider, provider_listing_id, variation_id). installation_id is required
-- in the tuple because the real PK of the parent `listings` table is
-- (tenant_id, installation_id, provider_listing_id, variation_id) —
-- 0036_listings.sql:2-31 — a child key without installation_id could not
-- address parent rows per installation.
--
-- No FK to listings: repo convention for this table shape (order_shipments,
-- 0088, keeps its own soft link to orders unenforced on purpose); backfill
-- and per-listing hydration can each populate parent/child independently.
--
-- No raw column: the parent listing's own `raw` (0090) already covers the
-- full /items payload; a variation has no independent raw fetch.
CREATE TABLE IF NOT EXISTS listing_variations (
    tenant_id           text        NOT NULL,
    installation_id     text        NOT NULL,
    provider            text        NOT NULL,
    provider_listing_id text        NOT NULL,
    variation_id        text        NOT NULL,
    price               numeric(14,2),
    available_quantity  int,
    sold_quantity       int,
    seller_sku          text,
    attributes          jsonb,
    last_seen_at        timestamptz,
    absent_since        timestamptz,
    created_at          timestamptz NOT NULL DEFAULT now(),
    updated_at          timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, installation_id, provider, provider_listing_id, variation_id)
);
