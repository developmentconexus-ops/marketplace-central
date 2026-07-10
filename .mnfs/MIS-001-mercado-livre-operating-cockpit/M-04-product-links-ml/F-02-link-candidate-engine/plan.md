# Feature Plan

```yaml
id: F-02
type: feature-plan
status: planned
owner: Feature Implementer
parent: F-02
created: 2026-07-08
updated: 2026-07-08
validation_level: QA-0
lifecycle_scope: feature
```

## Steps

1. Add `product_links` candidate domain types, persistence, and listing-snapshot read support.
2. Implement exact-first candidate generation over persisted snapshots using the `internal_read` product-linking boundary.
3. Add manual candidate-generation and candidate-listing HTTP surfaces and wire them into the root composition.
4. Update OpenAPI and `packages/sdk-runtime` for the new contract.
5. Run focused Go and SDK tests, then validate generation live against imported Mercado Livre snapshots plus Oracle-backed product lookup.

## Verification Commands

- Command: `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.gocache'; go test ./apps/server_core/internal/modules/product_links/... ./apps/server_core/internal/composition -count=1`
  - Satisfies criterion ID: M-04-F02-C01
  - Expected result: Pass. Candidate generation behavior and transport contract hold.
- Command: `npm test -- --run packages/sdk-runtime/src/index.test.ts`
  - Satisfies criterion ID: M-04-F02-C02
  - Expected result: Pass. SDK exposes generation/listing APIs.
- Command: `Invoke-RestMethod -Uri 'http://localhost:8080/product-links/link-candidates/generations' -Method Post -ContentType 'application/json' -Body '{"installation_id":"inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98","limit":20}'`
  - Satisfies criterion ID: M-04-F02-C03
  - Expected result: Pass. Live candidates are generated from Mercado Livre snapshots against Oracle lookup.

## Rollback/Risk Notes

- Risk: exact signals from `seller_sku` and `ean` disagree and the code accidentally prefers one instead of surfacing conflict.
- Recovery: keep exact-match arbitration centralized in one generation service and re-run live generation after fixes.

## Handoff

- Current status: `planned`
- Next owner: Feature Implementer
- Next action: keep manual-first flow and let F-03 consume persisted candidates
- Required files/evidence: spec, validation, migration, API/SDK contract
- Blockers or open decisions: scheduler deferred on purpose
