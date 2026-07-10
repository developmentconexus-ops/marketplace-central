CREATE TABLE IF NOT EXISTS profitability_profit_snapshots (
    tenant_id text NOT NULL,
    installation_id text NOT NULL,
    provider_order_id text NOT NULL,
    provider_item_id text NOT NULL DEFAULT '',
    provider_variation_id text NOT NULL DEFAULT '',
    scope text NOT NULL,
    currency text NOT NULL,
    revenue_amount numeric(14,2) NULL,
    sale_fee_amount numeric(14,2) NULL,
    cost_amount numeric(14,2) NULL,
    tax_amount numeric(14,2) NULL,
    freight_amount numeric(14,2) NULL,
    commission_amount numeric(14,2) NULL,
    adjustment_amount numeric(14,2) NULL,
    contribution_amount numeric(14,2) NULL,
    margin_percent numeric(14,2) NULL,
    quality text NOT NULL,
    flags_json jsonb NOT NULL DEFAULT '[]'::jsonb,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_profitability_profit_snapshots_installation
    ON profitability_profit_snapshots (tenant_id, installation_id, provider_order_id, provider_item_id, updated_at DESC);
