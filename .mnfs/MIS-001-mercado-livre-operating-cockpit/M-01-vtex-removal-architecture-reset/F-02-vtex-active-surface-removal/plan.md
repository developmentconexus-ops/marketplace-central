# Feature Plan

```yaml
id: F-02
type: feature-plan
status: executed
owner: Feature Implementer
parent: F-02-vtex-active-surface-removal
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
split_decision: single
split_reason: Active VTEX backend, OpenAPI, SDK, and frontend surfaces are tightly coupled; partial removal would leave broken contracts, so one coordinated removal with full validation is safer.
```

## Feature ID

F-02-vtex-active-surface-removal

## Steps

1. Remove VTEX publication routes from connector transport while preserving Melhor Envio auth/status routes.
2. Remove VTEX adapter/root-router wiring and the runtime requirement for `VTEX_APP_KEY`/`VTEX_APP_TOKEN`.
3. Delete VTEX pipeline code and VTEX-specific backend tests; update generic tests to non-VTEX fixtures.
4. Remove VTEX OpenAPI paths, schemas, operation ids, and error enum entries.
5. Remove SDK VTEX types/client methods/tests.
6. Remove frontend VTEX publisher routes/nav/pages/package dependency; update generic frontend fixtures/placeholders.
7. Run Go tests, frontend tests/build, and targeted `rg` checks; record evidence in `validation.md`.

## Files Expected To Change

- Path: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-01-vtex-removal-architecture-reset/F-02-vtex-active-surface-removal/spec.md`
- Reason: Feature spec and acceptance criteria.

- Path: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-01-vtex-removal-architecture-reset/F-02-vtex-active-surface-removal/plan.md`
- Reason: Feature plan and verification mapping.

- Path: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-01-vtex-removal-architecture-reset/F-02-vtex-active-surface-removal/validation.md`
- Reason: Quick validation evidence.

- Path: `apps/server_core/internal/composition/root.go`
- Reason: Remove VTEX adapter import/wiring and credential startup requirement.

- Path: `apps/server_core/internal/modules/connectors/transport/http_handler.go`
- Reason: Preserve Melhor Envio routes and remove VTEX publication/validation routes.

- Path: `apps/server_core/internal/modules/connectors/application/orchestrator.go`, `executor.go`
- Reason: Delete VTEX publication pipeline application code.

- Path: `apps/server_core/internal/modules/connectors/domain/batch.go`, `errors.go`, `mapping.go`, `pipeline.go`
- Reason: Delete VTEX publication pipeline domain code.

- Path: `apps/server_core/internal/modules/connectors/ports/repository.go`, `vtex_catalog.go`
- Reason: Delete VTEX publication repository/adapter ports.

- Path: `apps/server_core/internal/modules/connectors/adapters/postgres/repository.go`, `apps/server_core/internal/modules/connectors/adapters/vtex/**`
- Reason: Delete active VTEX persistence adapter and provider adapter code.

- Path: `apps/server_core/tests/**`
- Reason: Delete VTEX-specific tests and replace VTEX generic fixtures.

- Path: `contracts/api/marketplace-central.openapi.yaml`
- Reason: Remove active VTEX contract paths/schemas/errors.

- Path: `packages/sdk-runtime/src/index.ts`, `packages/sdk-runtime/src/index.test.ts`
- Reason: Remove VTEX SDK types/methods/tests.

- Path: `apps/web/src/app/AppRouter.tsx`, `apps/web/src/app/Layout.tsx`, `apps/web/src/index.css`, `apps/web/package.json`
- Reason: Remove VTEX page/nav/package dependency/source scan.

- Path: `packages/feature-connectors/**`
- Reason: Delete VTEX publisher UI package.

- Path: `packages/feature-*/src/**`
- Reason: Replace generic test/placeholders that use VTEX as an active example.

## Verification Commands

- Command: `GOCACHE=.gocache go test ./...` from `apps/server_core`
- Satisfies criterion ID: M-01-C01
- Expected result: Go tests pass after backend VTEX removal while generic integrations/connectors behavior remains valid.

- Command: `npm test -- --run`
- Satisfies criterion ID: M-01-C01, M-01-C02
- Expected result: Frontend and SDK tests pass after removing VTEX SDK/UI surfaces.

- Command: `npm run build`
- Satisfies criterion ID: M-01-C01, M-01-C02
- Expected result: Web build succeeds without feature-connectors dependency or VTEX routes.

- Command: `rg -n "connectors/vtex|adapters/vtex|VTEX_APP_|VTEX_ACCOUNT|VTEXCatalogPort|publishToVTEX|VTEXPublishPage|VTEXProduct|CONNECTORS_VTEX" apps packages contracts`
- Satisfies criterion ID: M-01-C01, M-01-C02
- Expected result: No active runtime/contract/SDK/UI matches remain.

- Command: `rg -n "VTEX|vtex" apps packages contracts wiki docs`
- Satisfies criterion ID: M-01-C01
- Expected result: Remaining matches are absent from active code/contract/SDK/UI or are historical docs for F-03.

## QA Steps

- Step: Inspect root router tests to confirm startup no longer depends on VTEX credentials and Melhor Envio route registration remains covered.
- Expected result: Test names/expectations reflect the new Mercado Livre-first reset.

- Step: Inspect OpenAPI/SDK surfaces after removal.
- Expected result: Contract and SDK both omit VTEX publish/validate/batch operations.

## Rollback/Risk Notes

- Rollback requires restoring the removed files from Git if F-02 fails validation.
- Migration history is not edited in F-02; existing databases may still contain legacy VTEX tables/columns until a future forward migration.
- Removing `feature-connectors` requires deleting the web dependency and Tailwind source reference together.
- Generic tests must retain coverage by using Mercado Livre/provider-neutral fixtures rather than deleting unrelated marketplace behavior.

## Handoff

- Current status: `executed`
- Next owner: Milestone Orchestrator
- Next action: Validate milestone criteria against F-02 evidence and remaining legacy classifications.
- Required files/evidence: changed paths, Go/frontend test output, targeted `rg`.
- Blockers or open decisions: none.
