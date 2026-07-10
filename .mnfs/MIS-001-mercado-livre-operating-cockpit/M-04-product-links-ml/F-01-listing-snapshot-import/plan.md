# Feature Plan

```yaml
id: F-01
type: feature-plan
status: planned
owner: Feature Implementer
parent: F-01
created: 2026-07-08
updated: 2026-07-08
validation_level: QA-0
lifecycle_scope: feature
```

## Steps

1. Add a minimal `product_links` module with a listing-snapshot domain model, import service, and tenant-scoped Postgres upsert repository.
2. Normalize connector listing output into one linkable snapshot row per item or variation.
3. Add a manual HTTP import endpoint and wire it through `composition/root.go`.
4. Update OpenAPI and `packages/sdk-runtime` so the import trigger is part of the supported contract.
5. Run focused Go and SDK tests, then validate the import against the live Mercado Livre installation and confirm idempotent re-import behavior in Postgres.

## Verification Commands

- Command: `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.gocache'; go test ./apps/server_core/internal/modules/product_links/... ./apps/server_core/internal/composition -count=1`
  - Satisfies criterion ID: M-04-F01-C01
  - Expected result: Pass. Import normalization and transport contract hold.
- Command: `npm test -- --run packages/sdk-runtime/src/index.test.ts`
  - Satisfies criterion ID: M-04-F01-C02
  - Expected result: Pass. SDK exposes the manual import contract.
- Command: `Invoke-RestMethod -Uri 'http://localhost:8080/product-links/listing-snapshots/imports' -Method Post -ContentType 'application/json' -Body '{"installation_id":"inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98","limit":5}'`
  - Satisfies criterion ID: M-04-F01-C03
  - Expected result: Pass. Live Mercado Livre snapshots import into MPC.

## Rollback/Risk Notes

- Risk: provider rows duplicate when item/variation identity is not normalized consistently.
- Recovery: keep idempotency at the Postgres key level and re-import from the live provider surface after fixes.

## Handoff

- Current status: `planned`
- Next owner: Feature Implementer
- Next action: keep this slice read-only and prepare candidate generation on top of persisted snapshots
- Required files/evidence: spec, validation, migration, API/SDK contract
- Blockers or open decisions: none for the manual-first slice
