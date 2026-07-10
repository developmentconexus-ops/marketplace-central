CREATE TABLE IF NOT EXISTS product_link_listing_snapshots (
  tenant_id text NOT NULL,
  installation_id text NOT NULL,
  provider_code text NOT NULL,
  provider_item_id text NOT NULL,
  provider_variation_id text NOT NULL DEFAULT '',
  provider_status text NOT NULL DEFAULT '',
  seller_sku text NOT NULL DEFAULT '',
  ean text NOT NULL DEFAULT '',
  title text NOT NULL DEFAULT '',
  available_quantity integer,
  source_updated_at timestamptz,
  fetched_at timestamptz NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (tenant_id, installation_id, provider_item_id, provider_variation_id)
);
