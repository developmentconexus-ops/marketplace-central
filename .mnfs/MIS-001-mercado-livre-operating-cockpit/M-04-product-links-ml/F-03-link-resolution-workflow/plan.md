# Feature Plan

```yaml
id: F-03
type: feature-plan
status: planned
owner: Feature Implementer
parent: F-03
created: 2026-07-08
updated: 2026-07-08
validation_level: QA-0
lifecycle_scope: feature
```

## Steps

1. Add backend link-resolution workflow types, persistence, migration, and tests for approve/reject/manual-resolve with audit evidence.
2. Extend transport, OpenAPI, and `packages/sdk-runtime` with list filters, resolution commands, and response models that expose current link state plus audit evidence.
3. Add a dedicated product-links UI route and package that lists candidates, shows evidence, and performs operator actions through the SDK.
4. Run focused backend, SDK, and UI tests, then record evidence in `validation.md`.

## Verification Commands

- Command: `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.gocache'; go test ./apps/server_core/internal/modules/product_links/... ./apps/server_core/internal/composition -count=1`
  - Satisfies criterion ID: M-04-C01
  - Expected result: Pass. Candidate generation plus resolution workflow stay green together.
- Command: `npm run test --workspace @marketplace-central/sdk-runtime -- src/index.test.ts`
  - Satisfies criterion ID: M-04-C01
  - Expected result: Pass. SDK reflects new product-links workflow routes and models.
- Command: `npm run test --workspace @marketplace-central/feature-product-links`
  - Satisfies criterion ID: M-04-C01
  - Expected result: Pass. UI proves operator workflow states and actions.

## Rollback/Risk Notes

- Risk: candidate generation truth and operator link truth drift if stored in the same structure.
- Recovery: keep generated candidates immutable as recommendation evidence and persist operator truth in dedicated link/audit records.
- Risk: UI gets forced into the integrations surface and couples unrelated operator flows.
- Recovery: keep product-links UI in its own feature package and route.

## Handoff

- Current status: `planned`
- Next owner: Feature Implementer
- Next action: implement backend workflow and audit persistence first
- Required files/evidence: migration, backend tests, API/SDK/UI parity, updated validation artifact
- Blockers or open decisions: none
