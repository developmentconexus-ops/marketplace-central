CREATE TABLE IF NOT EXISTS orders_marketplace_orders (
    tenant_id text NOT NULL,
    installation_id text NOT NULL,
    provider_code text NOT NULL,
    provider_order_id text NOT NULL,
    provider_status text NOT NULL DEFAULT '',
    provider_status_detail text NOT NULL DEFAULT '',
    provider_created_at timestamptz NULL,
    provider_closed_at timestamptz NULL,
    provider_updated_at timestamptz NULL,
    fetched_at timestamptz NOT NULL,
    shipping_id text NOT NULL DEFAULT '',
    cancellation_detail text NOT NULL DEFAULT '',
    tags_json jsonb NOT NULL DEFAULT '[]'::jsonb,
    raw_provider_ref jsonb NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    PRIMARY KEY (tenant_id, installation_id, provider_order_id)
);

CREATE TABLE IF NOT EXISTS orders_marketplace_order_items (
    tenant_id text NOT NULL,
    installation_id text NOT NULL,
    provider_order_id text NOT NULL,
    line_no integer NOT NULL,
    provider_item_id text NOT NULL,
    provider_variation_id text NOT NULL DEFAULT '',
    seller_sku text NOT NULL DEFAULT '',
    title text NOT NULL DEFAULT '',
    quantity integer NOT NULL,
    unit_price numeric(14,2) NULL,
    sale_fee_amount numeric(14,2) NULL,
    link_quality text NOT NULL,
    internal_product_id integer NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    PRIMARY KEY (tenant_id, installation_id, provider_order_id, line_no)
);

CREATE TABLE IF NOT EXISTS orders_marketplace_order_payments (
    tenant_id text NOT NULL,
    installation_id text NOT NULL,
    provider_order_id text NOT NULL,
    line_no integer NOT NULL,
    provider_payment_id text NOT NULL,
    provider_status text NOT NULL DEFAULT '',
    transaction_amount numeric(14,2) NULL,
    total_paid_amount numeric(14,2) NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    PRIMARY KEY (tenant_id, installation_id, provider_order_id, line_no)
);

CREATE INDEX IF NOT EXISTS idx_orders_marketplace_orders_tenant_updated
    ON orders_marketplace_orders (tenant_id, installation_id, provider_updated_at DESC, fetched_at DESC);
