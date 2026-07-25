-- 0081: remove the marketplaces this platform never integrated with, and the
-- tables no code reads.
--
-- The platform integrates with exactly one marketplace (Mercado Livre) and one
-- ERP (Sankhya). Magalu, Shopee, Amazon, Leroy Merlin and Madeira Madeira were
-- carried as seeded definitions with invented commission rates and adapters
-- whose only implemented method returned ErrNotImplemented: the product offered
-- accounts it could never read, write or authorize, and priced against rates
-- nobody quoted. Their Go adapters are deleted in the same change, so nothing
-- re-seeds these rows at boot.
--
-- The seeded Mercado Livre fee rows go too. They do not match how this platform
-- computes fees: the commission comes per listing from the ML Fees API
-- (sale_fee, per unit) and the fixed fee is 50% of the price for items up to
-- R$12,50 — not one flat 16%/22% per listing type with a zero fixed fee. Worse,
-- they contradicted pricing_tariff_defaults (13%/16%), which is what the
-- simulator actually reads. The table stays as the sink for api_sync rows.

DELETE FROM marketplace_fee_schedules;

-- Pricing policies hang off an account, not a marketplace code, so they follow
-- the accounts they belong to.
DELETE FROM marketplace_pricing_policies
WHERE account_id IN (
  SELECT account_id FROM marketplace_accounts WHERE marketplace_code <> 'mercado_livre'
);

DELETE FROM marketplace_accounts
WHERE marketplace_code <> 'mercado_livre';

DELETE FROM marketplace_definitions
WHERE marketplace_code <> 'mercado_livre';

DELETE FROM integration_provider_definitions
WHERE provider_code <> 'mercado_livre';

-- Orphan tables: created by earlier migrations, referenced by no production
-- code (only by the migration tests that assert their shape). They are dropped
-- rather than left empty so the schema stops advertising capabilities the app
-- does not have.
--   vtex_entity_mappings                        — VTEX was never integrated.
--   publication_batches / publication_operations — publication pipeline never built.
--   pipeline_step_results                        — same pipeline.
--   platform_migrations                          — superseded by schema_migrations.
--   catalog_legacy_product_reference_evidence    — legacy catalog reference model.
--   catalog_codprod_compatibility_candidates     — same.
-- Dropped child-first: pipeline_step_results references publication_operations,
-- which references publication_batches.
DROP TABLE IF EXISTS vtex_entity_mappings;
DROP TABLE IF EXISTS pipeline_step_results;
DROP TABLE IF EXISTS publication_operations;
DROP TABLE IF EXISTS publication_batches;
DROP TABLE IF EXISTS platform_migrations;
DROP TABLE IF EXISTS catalog_legacy_product_reference_evidence;
DROP TABLE IF EXISTS catalog_codprod_compatibility_candidates;
