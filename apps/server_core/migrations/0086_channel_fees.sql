-- 0086: channel_fees — IC-01 shared fee ledger (M-02 sync-core-seam, F-01-core-ddl).
--
-- One current-value row per (tenant, provider, installation, subject, layer,
-- fee_kind); the natural unique key below is the upsert target. No history in
-- this migration — additive later if the mission needs it (IC-01 Compatibility
-- Rules). No producer writes into this table yet: M-05 (layer 2, listings) and
-- M-06 (layer 3, orders) land in later milestones against this shape.
--
-- subject_id is an opaque, provider-verbatim string. The format is documented
-- here for readers, never enforced by the DB (IC-01 Fields):
--   listing    -> <provider_listing_id>                  (e.g. MLB123)
--   order      -> <provider_order_id>
--   order_line -> <provider_order_id>:<provider_item_id>
--   category   -> <category_id>                           (MIS-008)
--
-- currency is NOT NULL exactly when value_type='amount' — a percent fee has no
-- currency of its own (channel_fees_currency_when_amount_check).
--
-- origem includes api_shipping_options as an additive reservation (IC-01): no
-- producer in this mission writes it yet — listing freight stays
-- honest-unknown until then.
--
-- No FK: subject_id/installation_id reference provider-side entities, not rows
-- this schema owns. No raw column: detail is the only free-form payload, and
-- it is the source's own decomposition (IC-01 Canonical Examples), not a full
-- API dump.
CREATE TABLE IF NOT EXISTS channel_fees (
    id               bigserial      PRIMARY KEY,
    tenant_id        text           NOT NULL,
    provider         text           NOT NULL,
    installation_id  text           NOT NULL,
    layer            smallint       NOT NULL,
    fee_kind         text           NOT NULL,
    subject_type     text           NOT NULL,
    subject_id       text           NOT NULL,
    value_type       text           NOT NULL,
    value            numeric(14,4)  NOT NULL,
    currency         char(3),
    detail           jsonb,
    origem           text           NOT NULL,
    coletado_em      timestamptz    NOT NULL,
    source_time      timestamptz,
    created_at       timestamptz    NOT NULL DEFAULT now(),
    updated_at       timestamptz    NOT NULL DEFAULT now(),
    CONSTRAINT channel_fees_layer_check
        CHECK (layer IN (1, 2, 3)),
    CONSTRAINT channel_fees_fee_kind_check
        CHECK (fee_kind IN ('commission', 'freight')),
    CONSTRAINT channel_fees_subject_type_check
        CHECK (subject_type IN ('category', 'listing', 'order', 'order_line')),
    CONSTRAINT channel_fees_value_type_check
        CHECK (value_type IN ('percent', 'amount')),
    CONSTRAINT channel_fees_currency_when_amount_check
        CHECK ((value_type = 'amount' AND currency IS NOT NULL)
            OR (value_type = 'percent' AND currency IS NULL)),
    CONSTRAINT channel_fees_origem_check
        CHECK (origem IN ('api_listing_prices', 'api_shipping_options', 'api_order', 'api_shipment', 'config')),
    CONSTRAINT channel_fees_natural_key_key UNIQUE (tenant_id, provider, installation_id, subject_type, subject_id, layer, fee_kind)
);
