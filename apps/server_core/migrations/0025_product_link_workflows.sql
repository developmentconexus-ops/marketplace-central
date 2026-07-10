CREATE TABLE IF NOT EXISTS product_links (
  tenant_id text NOT NULL,
  installation_id text NOT NULL,
  provider_code text NOT NULL,
  provider_item_id text NOT NULL,
  provider_variation_id text NOT NULL DEFAULT '',
  state text NOT NULL,
  source_candidate_id text NOT NULL DEFAULT '',
  internal_product_id integer,
  internal_product_name text NOT NULL DEFAULT '',
  internal_reference_code text NOT NULL DEFAULT '',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (tenant_id, installation_id, provider_item_id, provider_variation_id)
);

CREATE INDEX IF NOT EXISTS product_links_installation_idx
  ON product_links (tenant_id, installation_id, updated_at DESC);

CREATE TABLE IF NOT EXISTS product_link_audit_entries (
  tenant_id text NOT NULL,
  audit_id text NOT NULL,
  installation_id text NOT NULL,
  provider_code text NOT NULL,
  provider_item_id text NOT NULL,
  provider_variation_id text NOT NULL DEFAULT '',
  action text NOT NULL,
  reason text NOT NULL DEFAULT '',
  source_candidate_id text NOT NULL DEFAULT '',
  actor_type text NOT NULL DEFAULT '',
  actor_id text NOT NULL DEFAULT '',
  actor_name text NOT NULL DEFAULT '',
  previous_state text NOT NULL,
  next_state text NOT NULL,
  previous_internal_product_id integer,
  next_internal_product_id integer,
  created_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (tenant_id, audit_id)
);

CREATE INDEX IF NOT EXISTS product_link_audit_entries_installation_idx
  ON product_link_audit_entries (tenant_id, installation_id, created_at DESC);
