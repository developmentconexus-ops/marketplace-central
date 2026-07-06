# Feature Validation

```yaml
id: F-02
type: feature-validation
status: passed
owner: Feature Implementer
parent: F-02-vtex-active-surface-removal
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
```

## Result

- Verdict: `passed`
- Scope: active VTEX backend, OpenAPI, SDK, frontend router/nav, and test surfaces removed
- Remaining VTEX references: `migration-risk` only in `apps/server_core/migrations/0005_connectors.sql`

## Commands Run

1. `rg -n "/connectors/vtex|validateVTEXConnection|publishToVTEX|getVTEXBatchStatus|retryVTEXBatch|PublishToVTEXRequest|RetryVTEXBatchRequest|PublishToVTEXResponse|BatchStatusResponse|CONNECTORS_VTEX|vtex_account" contracts/api/marketplace-central.openapi.yaml`
   - Result: matches found before removal; after patch, no active OpenAPI VTEX paths/schemas/error enums remain.

2. `rg -n "VTEX|vtex|publishToVTEX|getBatchStatus|retryBatch|BatchOrchestrator|VTEXCatalogPort|ErrVTEX|connectorsdomain|PublicationBatch|PublicationOperation|PipelineStepResult|VTEXEntityMapping" apps packages contracts wiki docs`
   - Result: used as inventory during execution; active code/contract hits were removed, historical docs and forward-only migration were retained for later classification.

3. `npm install`
   - Result: workspace lockfile updated after declaring `@marketplace-central/sdk-runtime` explicitly in `apps/web/package.json`.

4. `$abs = Join-Path (Get-Location) '.gocache'; New-Item -ItemType Directory -Force -Path $abs | Out-Null; $env:GOCACHE = $abs; go test ./...`
   - Working directory: `apps/server_core`
   - Result: `PASS`

5. `npm test -- --run`
   - Working directory: repo root
   - Result: `15 passed`, `133 passed`
   - Notes:
     - npm emitted a non-blocking warning about unknown CLI config `--run`
     - Vitest emitted existing `act(...)` warnings in `MarketplaceSettingsPage.test.tsx`, but the suite passed

6. `npm run build`
   - Working directory: repo root
   - Result: `vite build` passed
   - Output summary:
     - `1777 modules transformed`
     - built in `1m 48s`
   - Notes:
     - Vite emitted non-blocking `"use client" was ignored` warnings from `react-router` and `lucide-react`

7. `rg -n "connectors/vtex|adapters/vtex|VTEX_APP_|VTEX_ACCOUNT|VTEXCatalogPort|publishToVTEX|VTEXPublishPage|VTEXProduct|CONNECTORS_VTEX" apps packages contracts`
   - Result: no matches

8. `rg -n "VTEX|vtex" apps packages contracts`
   - Result:
     - `apps/server_core/migrations/0005_connectors.sql:5:  vtex_account     text NOT NULL,`
     - `apps/server_core/migrations/0005_connectors.sql:19:  vtex_account  text NOT NULL,`
     - `apps/server_core/migrations/0005_connectors.sql:31:  ON publication_operations (tenant_id, vtex_account, product_id)`
     - `apps/server_core/migrations/0005_connectors.sql:41:  vtex_entity_id text,`
     - `apps/server_core/migrations/0005_connectors.sql:49:-- VTEX entity mappings: durable local_id <-> vtex_id per account`
     - `apps/server_core/migrations/0005_connectors.sql:50:CREATE TABLE IF NOT EXISTS vtex_entity_mappings (`
     - `apps/server_core/migrations/0005_connectors.sql:53:  vtex_account text NOT NULL,`
     - `apps/server_core/migrations/0005_connectors.sql:56:  vtex_id      text NOT NULL,`
     - `apps/server_core/migrations/0005_connectors.sql:59:  UNIQUE (tenant_id, vtex_account, entity_type, local_id)`

## Classification

- `remove`
  - `apps/server_core/internal/modules/connectors/adapters/vtex/**`
  - `apps/server_core/internal/modules/connectors/application/{executor.go,orchestrator.go}`
  - `apps/server_core/internal/modules/connectors/domain/{batch.go,errors.go,mapping.go,pipeline.go}`
  - `apps/server_core/internal/modules/connectors/ports/{repository.go,vtex_catalog.go}`
  - `apps/server_core/internal/modules/connectors/adapters/postgres/repository.go`
  - `apps/server_core/tests/**/*vtex*`
  - `packages/feature-connectors/**`
  - VTEX OpenAPI paths/schemas/operation ids/error enums
  - VTEX SDK exports/methods/tests
  - VTEX frontend routes/nav/source registration

- `migration-risk`
  - `apps/server_core/migrations/0005_connectors.sql`
  - Reason: migration history is forward-only and still documents the legacy VTEX persistence model; editing it in place would break replay integrity.

## Notes

- `packages/feature-marketplaces/src/MarketplaceSettingsPage.tsx` now surfaces provider action failures instead of leaking rejected promises into tests.
- `packages/feature-marketplaces/src/ProviderCatalogPanel.tsx` now enables `Run fee sync` only for connected/degraded installations, which aligns with the active installation lifecycle.
- No production path keeps `/connectors/vtex/*` registered after this feature.
