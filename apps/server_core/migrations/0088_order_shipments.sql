-- 0088: order_shipments — IC-03 shipment persistence (M-02 sync-core-seam,
-- F-01-core-ddl).
--
-- Typed shipment facts only. No `raw` column, by design (IC-03, P7 r01 B-7,
-- ADR-03 amended): a shipment payload carries PII of delivery (receiver name,
-- street address, ZIP) — the same class of fact `cmd/mlprobe/main.go` flags as
-- PII. `raw jsonb` stays scoped to `listings` only; the typed columns below are
-- the COMPLETE persistence of shipment for this schema.
--
-- status/substatus are the provider's own vocabulary, verbatim — no CHECK,
-- because that vocabulary is owned by Mercado Livre, not by this repo, and can
-- change without our say.
--
-- No FK: provider_order_id is a soft link to orders_marketplace_orders, kept
-- unenforced on purpose (order_shipments and orders can each independently be
-- populated first during backfill).
--
-- sla_limit_at comes from GET /shipments/{id}'s lead_time.estimated_delivery_limit.date, NOT
-- from the separate /shipments/{id}/sla sub-resource (deliberately never called — honest gap,
-- see shipment_ingest_reader.go's GetShipmentDetail doc comment); cost fields from /costs
-- (header X-Costs-New); destination fields are what the screen already shows today.
CREATE TABLE IF NOT EXISTS order_shipments (
    tenant_id            text           NOT NULL,
    provider             text           NOT NULL,
    provider_shipment_id text           NOT NULL,
    provider_order_id    text           NOT NULL,
    status               text           NOT NULL,
    substatus            text,
    logistic_type        text,
    tracking_number      text,
    tracking_method      text,
    sla_status           text,
    sla_limit_at         timestamptz,
    cost_gross           numeric(14,2),
    cost_seller          numeric(14,2),
    currency             char(3),
    receiver_name        text,
    dest_city            text,
    dest_state           text,
    dest_zip             text,
    source_time          timestamptz,
    fetched_at           timestamptz    NOT NULL,
    created_at           timestamptz    NOT NULL DEFAULT now(),
    updated_at           timestamptz    NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, provider, provider_shipment_id)
);

CREATE INDEX IF NOT EXISTS idx_order_shipments_provider_order_id
    ON order_shipments (provider_order_id);
