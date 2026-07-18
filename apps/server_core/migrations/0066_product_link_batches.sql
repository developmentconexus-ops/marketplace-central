CREATE TABLE IF NOT EXISTS product_link_batches (
  tenant_id text NOT NULL,
  batch_id text NOT NULL,
  installation_id text NOT NULL,
  actor_type text,
  actor_id text,
  actor_name text,
  requested_count integer NOT NULL DEFAULT 0,
  applied_count integer NOT NULL DEFAULT 0,
  failed_count integer NOT NULL DEFAULT 0,
  status text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (tenant_id, batch_id)
);
