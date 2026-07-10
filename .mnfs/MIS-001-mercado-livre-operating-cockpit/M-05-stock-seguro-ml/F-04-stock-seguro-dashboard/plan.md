# Feature Plan

```yaml
id: F-04
type: feature-plan
status: planned
owner: Feature Implementer
parent: F-04
created: 2026-07-09
updated: 2026-07-09
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-04-stock-seguro-dashboard

## Implementation Steps

1. Add inventory backend contracts and persistence.
   - Create inventory ports for reading listing snapshots/workflows, reading installations, and storing stock actions.
   - Add a Postgres stock action repository and any required migration.
   - Add focused tests for repository/idempotency behavior.

2. Add inventory application services and transport.
   - Implement risk list orchestration from product links + internal read + inventory policy/risk classifier.
   - Implement manual apply orchestration that resolves installation/provider context and delegates to the Mercado Livre stock writer.
   - Add HTTP handlers and tests for request validation, success, and structured failures.

3. Update contract surfaces.
   - Extend OpenAPI for inventory risk/action endpoints.
   - Extend `packages/sdk-runtime` types, client methods, and tests.

4. Add the Stock Seguro UI.
   - Create `packages/feature-inventory` with filterable cockpit, row detail, and explicit confirmatory manual action flow.
   - Add UI tests for loading, error, empty, healthy, oversell, undersell, stale, unresolved, conflict, and action results.
   - Register the route and sidebar entry in `apps/web`.

5. Validate end to end.
   - Run focused Go tests for inventory/application/transport/repository paths.
   - Run web/package tests for router, SDK, and feature UI.
   - Use the built-in browser to validate `/inventory/stock-seguro` on desktop and mobile.
   - If operator approval is provided for a live stock write, record whether the action used a real provider target or stayed local-only.

## Validation Commands

- `cd apps/server_core; $env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/modules/inventory/... -count=1`
- `cd apps/server_core; $env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/composition ./internal/modules/product_links/... ./internal/modules/inventory/... -count=1`
- `npm run test --workspace @marketplace-central/sdk-runtime`
- `npm run test --workspace @marketplace-central/feature-inventory`
- `npm run test --workspace @marketplace-central/web`

## Risks

- Oracle/internal read dependency can make the risk list unavailable; handlers must fail clearly rather than fabricate zero stock.
- Existing product link snapshots may not cover every listing identity; unsupported rows still need deterministic rendering.
- Live provider writes must remain explicit and operator-approved during QA to avoid accidental production-like side effects.

## Handoff

- Current status: `planned`
- Next owner: Feature Implementer
- Next action: start backend inventory vertical, then contract, then UI
- Required files/evidence: implementation diff and `validation.md`
- Blockers or open decisions: migration naming and live-write QA target will be finalized during implementation
