-- 0089: additive columns on orders_marketplace_orders — IC-03 (M-02
-- sync-core-seam, F-01-core-ddl).
--
-- Every column here is NULL by default: honest-unknown until a later
-- milestone's producer (M-03 shipment-enrich, M-06 backfill+decomposition)
-- populates it. Zero ALTER that changes an existing column's type or
-- nullability — purely additive, IF NOT EXISTS on every ADD COLUMN so a
-- re-run is a no-op.
ALTER TABLE orders_marketplace_orders
    ADD COLUMN IF NOT EXISTS pack_id text;

-- Soft link to order_shipments (0088) — no FK, same reasoning as that table's
-- own soft link back to provider_order_id.
ALTER TABLE orders_marketplace_orders
    ADD COLUMN IF NOT EXISTS provider_shipment_id text;

-- Enum values come from the existing domain.DeriveOrderBucket
-- (orders/domain/order_bucket.go) — this migration does not invent a new
-- vocabulary, it only gives the derived value somewhere to land once the
-- derivation call site moves from read-path to ingest (M-03).
ALTER TABLE orders_marketplace_orders
    ADD COLUMN IF NOT EXISTS bucket text;

-- Watermark for incremental sync (IC-06) — verbatim from the provider, never
-- computed by us.
ALTER TABLE orders_marketplace_orders
    ADD COLUMN IF NOT EXISTS date_last_updated_ml timestamptz;

-- Buyer fiscal, typed only (ADR-03/R-6: raw billing_info payload is
-- forbidden). Exactly the fields the comprador_fiscal drawer renders today
-- per p5-prerequisites.md §2 (buyer_fiscal_reader.go DTO) — no more, no less:
-- uf_nome/state_name and fetched_at are deliberately excluded, they are not
-- rendered.
ALTER TABLE orders_marketplace_orders
    ADD COLUMN IF NOT EXISTS buyer_name text;
ALTER TABLE orders_marketplace_orders
    ADD COLUMN IF NOT EXISTS buyer_doc_type text;
ALTER TABLE orders_marketplace_orders
    ADD COLUMN IF NOT EXISTS buyer_doc_number text;
ALTER TABLE orders_marketplace_orders
    ADD COLUMN IF NOT EXISTS buyer_address_street text;
ALTER TABLE orders_marketplace_orders
    ADD COLUMN IF NOT EXISTS buyer_address_number text;
ALTER TABLE orders_marketplace_orders
    ADD COLUMN IF NOT EXISTS buyer_address_city text;
ALTER TABLE orders_marketplace_orders
    ADD COLUMN IF NOT EXISTS buyer_address_state_code text;
ALTER TABLE orders_marketplace_orders
    ADD COLUMN IF NOT EXISTS buyer_address_zip text;
ALTER TABLE orders_marketplace_orders
    ADD COLUMN IF NOT EXISTS buyer_address_country text;

-- Decomposition (IC-03 Persistence Expectations): the jsonb breakdown, plus
-- two columns extracted from it for sort/filter. Unknown decomposition
-- fields are absent from the JSON, never zero (IC-03); these two mirror
-- columns are written in the same transaction as the jsonb, never inferred
-- independently.
ALTER TABLE orders_marketplace_orders
    ADD COLUMN IF NOT EXISTS decomposition jsonb;
ALTER TABLE orders_marketplace_orders
    ADD COLUMN IF NOT EXISTS net_amount numeric(14,2);
ALTER TABLE orders_marketplace_orders
    ADD COLUMN IF NOT EXISTS margin_pct numeric(7,4);

CREATE INDEX IF NOT EXISTS idx_orders_marketplace_orders_bucket
    ON orders_marketplace_orders (tenant_id, bucket);
