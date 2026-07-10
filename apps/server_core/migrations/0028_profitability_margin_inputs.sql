CREATE TABLE IF NOT EXISTS profitability_margin_inputs (
    tenant_id text NOT NULL,
    installation_id text NOT NULL,
    provider_order_id text NOT NULL,
    provider_item_id text NOT NULL DEFAULT '',
    provider_variation_id text NOT NULL DEFAULT '',
    scope text NOT NULL,
    kind text NOT NULL,
    amount numeric(14,2) NULL,
    currency text NOT NULL,
    source_system text NOT NULL,
    source_reference text NOT NULL DEFAULT '',
    observed_at timestamptz NULL,
    quality text NOT NULL,
    quality_reason text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_profitability_margin_inputs_installation
    ON profitability_margin_inputs (tenant_id, installation_id, provider_order_id, provider_item_id, kind);

CREATE TABLE IF NOT EXISTS profitability_manual_adjustments (
    tenant_id text NOT NULL,
    adjustment_id text PRIMARY KEY,
    installation_id text NOT NULL,
    provider_order_id text NOT NULL,
    provider_item_id text NOT NULL DEFAULT '',
    provider_variation_id text NOT NULL DEFAULT '',
    scope text NOT NULL,
    category text NOT NULL,
    amount numeric(14,2) NOT NULL,
    currency text NOT NULL,
    reason text NOT NULL,
    actor_type text NOT NULL,
    actor_id text NOT NULL,
    actor_name text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL
);
