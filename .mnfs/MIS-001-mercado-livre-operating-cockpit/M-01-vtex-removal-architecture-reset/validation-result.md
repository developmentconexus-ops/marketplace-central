# Milestone Validation Result

```yaml
id: M-01
type: milestone-validation-result
status: passed
owner: QA Validator
parent: MIS-001
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: milestone
```

## Milestone

M-01-vtex-removal-architecture-reset

## Verdict

- Result: `passed`
- Blocking failures: none
- Summary: VTEX no longer has active runtime, contract, SDK, or frontend UI surfaces in Marketplace Central. Remaining VTEX references are limited to forward-only migrations and historical/legacy docs.

## Feature Evidence

- F-01 inventory: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-01-vtex-removal-architecture-reset/F-01-vtex-surface-inventory/validation.md`
- F-02 removal: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-01-vtex-removal-architecture-reset/F-02-vtex-active-surface-removal/validation.md`
- F-03 truth alignment: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-01-vtex-removal-architecture-reset/F-03-architecture-truth-alignment/validation.md`

## Criterion Review

### M-01-C01 — VTEX Active Surfaces Removed

- Status: `passed`
- Command:
  - `rg -n "VTEX|vtex" apps packages contracts wiki docs`
- Expected:
  - Matches are absent from active runtime/contract/SDK/UI paths or explicitly marked historical/legacy.
- Actual:
  - `apps/packages/contracts` active surfaces are clean except for forward-only migration residue in `apps/server_core/migrations/0005_connectors.sql`
  - `wiki/operations/environment-and-db.md` retains VTEX env keys as explicitly legacy-only
  - `docs/**` retains historical VTEX design/research/evidence material only
- Blocking failure observed:
  - `No`

### M-01-C02 — Contract And SDK Aligned

- Status: `passed`
- Commands:
  - `rg -n "connectors/vtex|adapters/vtex|publishToVTEX|getVTEXBatchStatus|retryVTEXBatch|VTEXPublishPage|VTEXProduct|CONNECTORS_VTEX|VTEX_APP_|VTEX_ACCOUNT" apps packages contracts`
  - `$abs = Join-Path (Get-Location) '.gocache'; New-Item -ItemType Directory -Force -Path $abs | Out-Null; $env:GOCACHE = $abs; go test ./...` from `apps/server_core`
  - `npm test -- --run`
  - `npm run build`
- Expected:
  - No `/connectors/vtex/*` routes or `publishToVTEX` SDK methods remain; tests pass.
- Actual:
  - Active route/SDK/OpenAPI search returned no matches
  - Backend test suite passed
  - Frontend/sdk tests passed (`15` files, `133` tests)
  - Web build passed (`vite build`, `1777 modules transformed`)
- Blocking failure observed:
  - `No`

## Remaining Residue Classification

- `migration-risk`
  - `apps/server_core/migrations/0005_connectors.sql`
  - Reason: forward-only migration history still contains legacy VTEX table/column names and must not be rewritten in place.

- `legacy-doc-retain`
  - `wiki/operations/environment-and-db.md`
  - `docs/vtex-integration-reference.md`
  - `docs/superpowers/**`
  - `.brain/decisions/005-mercado-livre-first-control-plane.md`
  - Reason: these files now document VTEX only as legacy/historical context, not as an active implementation target.

## Validation Notes

- During F-02 validation, `apps/web/package.json` needed an explicit `@marketplace-central/sdk-runtime` dependency and `package-lock.json` refresh so the web build could resolve the SDK directly. This is now fixed and validated.
- `packages/feature-marketplaces` received two small hardening fixes discovered by the validation suite:
  - provider action errors are now surfaced in the UI instead of leaking rejected promises
  - fee sync is now gated to `connected`/`degraded` installations
- Frontend test output still includes non-blocking React `act(...)` warnings in marketplace loading-state setup; suite passes and the warning remains outside M-01 scope.

## Handoff

- Milestone status: `ready for mission continuation`
- Next recommended action: start the Mercado Livre product-links and stock-reconciliation mission planning after mission owner acknowledges M-01 completion
- Open blockers: none for M-01
