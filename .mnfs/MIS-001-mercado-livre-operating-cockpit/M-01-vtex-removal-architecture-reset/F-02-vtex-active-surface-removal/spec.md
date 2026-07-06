# Feature Spec

```yaml
id: F-02
type: feature-spec
status: executed
owner: Feature Implementer
parent: F-02-vtex-active-surface-removal
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-02-vtex-active-surface-removal

## Problem

F-01 proved that Marketplace Central still exposes VTEX as an active runtime, contract, SDK, UI, adapter, and test surface. This contradicts ADR-005 and blocks Mercado Livre-first cockpit work because new features could still depend on the legacy VTEX publication pipeline.

## Requirements

- Requirement: Remove active `/connectors/vtex/*` backend routes and root composition dependencies on VTEX credentials/adapters.
- Acceptance evidence: Go tests and `rg` show no active `/connectors/vtex` route registration, VTEX adapter import, `VTEX_APP_*` runtime dependency, or `VTEXCatalogPort` production dependency.

- Requirement: Remove VTEX-specific connector pipeline code and provider adapter code while preserving generic connector fee syncers and Melhor Envio auth routes.
- Acceptance evidence: Go tests pass for remaining modules; `/connectors/melhor-envio/*` remains registered.

- Requirement: Remove OpenAPI VTEX paths, schemas, operation ids, and active VTEX error enum entries.
- Acceptance evidence: targeted `rg` over OpenAPI finds no `/connectors/vtex`, `publishToVTEX`, `validateVTEXConnection`, `RetryVTEXBatch`, or `CONNECTORS_VTEX_*`.

- Requirement: Remove SDK VTEX types/methods/tests and frontend VTEX publisher pages/navigation.
- Acceptance evidence: frontend tests/build pass; `rg` finds no `publishToVTEX`, `VTEXPublishPage`, `/connectors/vtex`, or `VTEXProduct` in active app/package sources.

- Requirement: Replace VTEX-only test fixtures/placeholders in generic marketplace/classification/product tests with non-VTEX wording.
- Acceptance evidence: package tests pass and active test fixtures no longer present VTEX as a current marketplace target.

## Non-Goals

- Do not edit historical migrations in place.
- Do not drop database tables or columns.
- Do not remove Mercado Livre, Shopee, Magalu, Amazon, Leroy Merlin, MadeiraMadeira, Melhor Envio, or generic integrations behavior.
- Do not implement new Mercado Livre listing/stock/order capabilities in F-02.
- Do not rewrite historical docs; F-03 owns documentation truth alignment.

## Design

Remove VTEX active surfaces in one coordinated edit:

- Backend: keep `connectors` module for fee syncers and Melhor Envio auth; remove VTEX publication transport, orchestrator/executor/domain/repository port/adapter code, and root router VTEX credential dependency.
- Contract/SDK: delete VTEX OpenAPI operations and SDK client methods/types.
- Frontend: remove `feature-connectors` pages from app dependencies/router/nav and delete VTEX page package.
- Tests: delete VTEX-specific tests and update generic marketplace fixtures to Mercado Livre or neutral labels.

`apps/server_core/migrations/0005_connectors.sql` remains unchanged as migration history and is tracked as migration risk for a later forward migration decision.

## Edge Cases

- `connectors` still has non-VTEX adapters used by integration provider fee syncing; these must remain.
- Some files mention VTEX only in historical docs or ADRs; F-02 does not rewrite them unless they are active code/test surfaces.
- Root router must build without VTEX credentials after this feature.
- `.env` may still contain local legacy keys; F-03/docs can mark/remove docs, but F-02 removes runtime use.

## Acceptance Criteria

- Criterion: Backend no longer mounts VTEX routes, imports VTEX adapter packages, or requires VTEX credentials to start.
- Traces to milestone criterion ID: M-01-C01
- Proven by (verification command or QA step): `GOCACHE=.gocache go test ./...` in `apps/server_core`; targeted `rg -n "connectors/vtex|adapters/vtex|VTEX_APP_|VTEXCatalogPort" apps/server_core/internal apps/server_core/tests`.

- Criterion: OpenAPI and SDK no longer expose VTEX publish/validate/batch surfaces.
- Traces to milestone criterion ID: M-01-C02
- Proven by (verification command or QA step): targeted `rg -n "/connectors/vtex|publishToVTEX|validateVTEXConnection|RetryVTEXBatch|CONNECTORS_VTEX|VTEXProduct" contracts packages/sdk-runtime`.

- Criterion: Frontend no longer navigates to or renders VTEX publisher/batch pages.
- Traces to milestone criterion ID: M-01-C01
- Proven by (verification command or QA step): `npm test -- --run` and `npm run build`; targeted `rg -n "VTEXPublishPage|/connectors/vtex|feature-connectors|VTEX Publisher" apps packages`.

- Criterion: Generic provider catalog and integrations behavior are preserved.
- Traces to milestone criterion ID: M-01-C01
- Proven by (verification command or QA step): backend tests including integrations/provider tests and frontend tests continue passing; Melhor Envio connector route registration remains tested.

## Handoff

- Current status: `executed`
- Next owner: Milestone Orchestrator
- Next action: Use removal evidence for M-01 validation and F-03 alignment.
- Required files/evidence: F-01 inventory, F-02 plan, code diffs, validation commands.
- Blockers or open decisions: none for active-surface removal; migration cleanup remains future migration-risk.
