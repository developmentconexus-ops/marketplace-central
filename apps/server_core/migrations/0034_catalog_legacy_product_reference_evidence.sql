-- Legacy text references are retained as audit evidence. They are deliberately
-- not attached to CODPROD until an approved deterministic mapping source exists.
CREATE TABLE IF NOT EXISTS catalog_legacy_product_reference_evidence (
  tenant_id text NOT NULL,
  consumer text NOT NULL,
  record_id text NOT NULL,
  legacy_product_id text NOT NULL,
  internal_product_id bigint,
  resolution_status text NOT NULL CHECK (resolution_status IN ('mapped', 'not_found', 'identity_conflict')),
  created_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (tenant_id, consumer, record_id)
);

INSERT INTO catalog_legacy_product_reference_evidence (tenant_id, consumer, record_id, legacy_product_id, internal_product_id, resolution_status)
SELECT tenant_id, 'classification', classification_id || ':' || product_id, product_id, NULL, 'not_found'
FROM classification_products
ON CONFLICT (tenant_id, consumer, record_id) DO NOTHING;

INSERT INTO catalog_legacy_product_reference_evidence (tenant_id, consumer, record_id, legacy_product_id, internal_product_id, resolution_status)
SELECT tenant_id, 'pricing_simulation', simulation_id, product_id, NULL, 'not_found'
FROM pricing_simulations
ON CONFLICT (tenant_id, consumer, record_id) DO NOTHING;
