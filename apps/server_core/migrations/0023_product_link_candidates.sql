CREATE TABLE IF NOT EXISTS product_link_candidates (
  tenant_id text NOT NULL,
  candidate_id text NOT NULL,
  installation_id text NOT NULL,
  provider_code text NOT NULL,
  provider_item_id text NOT NULL,
  provider_variation_id text NOT NULL DEFAULT '',
  internal_product_id integer,
  internal_product_name text NOT NULL DEFAULT '',
  internal_reference_code text NOT NULL DEFAULT '',
  state text NOT NULL,
  match_input text NOT NULL,
  match_value text NOT NULL DEFAULT '',
  source_snapshot_fetched_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (tenant_id, candidate_id)
);

CREATE INDEX IF NOT EXISTS product_link_candidates_installation_idx
  ON product_link_candidates (tenant_id, installation_id, updated_at DESC);
