# Feature Validation

```yaml
id: F-01
type: feature-validation
status: passed
owner: Feature Implementer
parent: F-01
created: 2026-07-08
updated: 2026-07-08
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-01-listing-snapshot-import

## Summary

`product_links` now imports Mercado Livre listing snapshots through the live installation-aware read surface, persists them in a dedicated tenant-scoped table, and supports idempotent manual re-import.

## Current Validation State

- Result: Passed for local contract behavior and live Mercado Livre read/import behavior
- Result owner: Feature Implementer
- Decision date: 2026-07-08
- Final feature state for handoff: ready_for_candidate_generation

## Evidence

- Command: `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.gocache'; go test ./apps/server_core/internal/modules/product_links/... ./apps/server_core/internal/composition -count=1`
  - Result: Pass
- Command: `npm test -- --run packages/sdk-runtime/src/index.test.ts`
  - Result: Pass
- Command: `Invoke-RestMethod -Uri 'http://localhost:8080/product-links/listing-snapshots/imports' -Method Post -ContentType 'application/json' -Body '{"installation_id":"inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98","limit":5}'`
  - Result: Pass
  - Observed: `imported_count = 5` with live Mercado Livre item data including `provider_item_id`, `ean`, `available_quantity`, `provider_status`, and `fetched_at`
- Command: `docker compose exec -T postgres psql -U marketplace -d marketplace_central -c "SELECT COUNT(*), MAX(updated_at) FROM product_link_listing_snapshots WHERE tenant_id = 'tenant_default' AND installation_id = 'inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98';"`
  - Result: Pass
  - Observed: first import persisted `5` rows; second import kept `COUNT(*) = 5` while `MAX(updated_at)` advanced, proving upsert-style idempotent re-import

## Observed

- Item-level listings without variations persist one row with empty `provider_variation_id`.
- Listings with variations normalize into variation-scoped rows instead of duplicating item-parent state.
- The manual import contract is now available in OpenAPI and `packages/sdk-runtime`.
- Live validation is read/import only; no Mercado Livre write operation was executed.

## Scope Declaration

- contract_validated: Yes
- integration_validated: Yes, for live Mercado Livre read/import into MPC Postgres
- provider_write_validated: No
- blocked_for_real_validation: No for this feature slice

## Handoff

- Current status: `passed`
- Next owner: Feature Implementer
- Next action: build F-02 candidate generation on top of persisted snapshot truth
- Required files/evidence: candidate rules, confidence states, unresolved/conflict semantics
- Blockers or open decisions: none for this import slice
