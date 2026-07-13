-- Forward-only compatibility bridge. Legacy text remains immutable evidence;
-- active consumers read only rows whose identity was proved by exact positive
-- CODPROD equality against governed candidate evidence.
CREATE TABLE IF NOT EXISTS catalog_codprod_compatibility_candidates (
  candidate_evidence_id text PRIMARY KEY,
  tenant_id text NOT NULL,
  legacy_product_id text NOT NULL,
  internal_product_id bigint NOT NULL CHECK (internal_product_id > 0),
  source text NOT NULL CHECK (source = 'oracle/sankhya'),
  observed_at timestamptz NOT NULL,
  read_only boolean NOT NULL CHECK (read_only)
);

ALTER TABLE product_enrichments ADD COLUMN IF NOT EXISTS internal_product_id bigint;
ALTER TABLE product_enrichments ADD COLUMN IF NOT EXISTS resolution_status text NOT NULL DEFAULT 'not_found';
ALTER TABLE classification_products ADD COLUMN IF NOT EXISTS internal_product_id bigint;
ALTER TABLE classification_products ADD COLUMN IF NOT EXISTS resolution_status text NOT NULL DEFAULT 'not_found';
ALTER TABLE pricing_simulations ADD COLUMN IF NOT EXISTS internal_product_id bigint;
ALTER TABLE pricing_simulations ADD COLUMN IF NOT EXISTS resolution_status text NOT NULL DEFAULT 'not_found';

ALTER TABLE product_enrichments ADD CONSTRAINT product_enrichments_codprod_resolution_check CHECK (
  (resolution_status = 'mapped' AND internal_product_id > 0 AND product_id = internal_product_id::text)
  OR (resolution_status IN ('not_found', 'identity_conflict') AND internal_product_id IS NULL)
);
ALTER TABLE classification_products ADD CONSTRAINT classification_products_codprod_resolution_check CHECK (
  (resolution_status = 'mapped' AND internal_product_id > 0 AND product_id = internal_product_id::text)
  OR (resolution_status IN ('not_found', 'identity_conflict') AND internal_product_id IS NULL)
);
ALTER TABLE pricing_simulations ADD CONSTRAINT pricing_simulations_codprod_resolution_check CHECK (
  (resolution_status = 'mapped' AND internal_product_id > 0 AND product_id = internal_product_id::text)
  OR (resolution_status IN ('not_found', 'identity_conflict') AND internal_product_id IS NULL)
);

INSERT INTO catalog_legacy_product_reference_evidence
  (tenant_id, consumer, record_id, legacy_product_id, internal_product_id, resolution_status)
SELECT tenant_id, 'product_enrichment', product_id, product_id, NULL, 'not_found'
FROM product_enrichments
ON CONFLICT (tenant_id, consumer, record_id) DO NOTHING;

CREATE OR REPLACE FUNCTION reconcile_catalog_codprod_compatibility() RETURNS void
LANGUAGE plpgsql AS $$
BEGIN
  WITH resolutions AS (
    SELECT e.tenant_id, e.consumer, e.record_id,
           count(c.candidate_evidence_id) AS candidate_count,
           min(c.internal_product_id) AS matched_id
    FROM catalog_legacy_product_reference_evidence e
    LEFT JOIN catalog_codprod_compatibility_candidates c
      ON c.tenant_id = e.tenant_id
     AND c.legacy_product_id = e.legacy_product_id
     AND e.legacy_product_id = c.internal_product_id::text
     AND c.internal_product_id > 0
    GROUP BY e.tenant_id, e.consumer, e.record_id
  )
  UPDATE catalog_legacy_product_reference_evidence e
  SET internal_product_id = CASE WHEN r.candidate_count = 1 THEN r.matched_id ELSE NULL END,
      resolution_status = CASE
        WHEN r.candidate_count = 1 THEN 'mapped'
        WHEN r.candidate_count = 0 THEN 'not_found'
        ELSE 'identity_conflict'
      END
  FROM resolutions r
  WHERE e.tenant_id = r.tenant_id AND e.consumer = r.consumer AND e.record_id = r.record_id;

  UPDATE product_enrichments p
  SET internal_product_id = e.internal_product_id, resolution_status = e.resolution_status
  FROM catalog_legacy_product_reference_evidence e
  WHERE e.tenant_id = p.tenant_id AND e.consumer = 'product_enrichment'
    AND e.record_id = p.product_id;

  UPDATE classification_products p
  SET internal_product_id = e.internal_product_id, resolution_status = e.resolution_status
  FROM catalog_legacy_product_reference_evidence e
  WHERE e.tenant_id = p.tenant_id AND e.consumer = 'classification'
    AND e.record_id = p.classification_id || ':' || p.product_id;

  UPDATE pricing_simulations p
  SET internal_product_id = e.internal_product_id, resolution_status = e.resolution_status
  FROM catalog_legacy_product_reference_evidence e
  WHERE e.tenant_id = p.tenant_id AND e.consumer = 'pricing_simulation'
    AND e.record_id = p.simulation_id;
END;
$$;

SELECT reconcile_catalog_codprod_compatibility();
