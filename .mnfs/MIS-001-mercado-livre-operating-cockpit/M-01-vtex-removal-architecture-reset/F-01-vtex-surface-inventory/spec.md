# Feature Spec

```yaml
id: F-01
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: F-01-vtex-surface-inventory
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-01-vtex-surface-inventory

## Problem

ADR-005 makes VTEX legacy and forbids new VTEX work, but active runtime, OpenAPI, SDK, UI, tests, docs, environment keys, and migration surfaces still reference VTEX. M-01 needs a precise inventory before F-02 removes active surfaces, so deletion does not accidentally break generic marketplace foundations or lose migration risk context.

## Requirements

- Requirement: Inventory active VTEX references across backend routes, adapters, ports, application/domain code, tests, OpenAPI, SDK runtime, frontend pages/navigation/tests, docs/wiki, environment files, and migrations.
- Acceptance evidence: `validation.md` records exact `rg`/inspection commands, categorized counts, and file-level findings.

- Requirement: Classify each VTEX reference group as `remove`, `legacy-doc-retain`, or `migration-risk`.
- Acceptance evidence: this spec lists every discovered VTEX file group by category and planned owner.

- Requirement: Do not remove or alter production code in F-01.
- Acceptance evidence: changed paths are limited to F-01 MNFS artifacts.

- Requirement: Do not classify Mercado Livre or generic marketplace references as VTEX.
- Acceptance evidence: inventory commands search explicit `VTEX|vtex` terms and classification separates provider-generic marketplace examples from active VTEX surfaces.

## Non-Goals

- Remove VTEX routes, adapter code, SDK methods, frontend pages, tests, docs, env keys, or database objects.
- Redesign connector/module boundaries.
- Add Mercado Livre capability ports or business modules.
- Run full test suites; this is an inventory-only feature.
- Copy secret values from `.env` into feature artifacts.

## Design

F-01 creates a deletion map for F-02 and a documentation alignment map for F-03. The classification is file-level except where a file contains both generic marketplace behavior and VTEX-only fixtures/placeholders; those mixed files are marked with the relevant local action.

Classification meanings:

- `remove`: active VTEX runtime, contract, SDK, UI, adapter, provider-specific test, env dependency, or fixture text that should disappear or be replaced with Mercado Livre/generic examples in F-02.
- `legacy-doc-retain`: historical or ADR/research references that should remain only when explicitly marked legacy/historical and not presented as target architecture.
- `migration-risk`: persisted database shape, historical migrations, or existing local secret material that needs a deliberate transition plan rather than blind deletion.

## Inventory Classification

### Backend Routes, Adapters, Ports, Application, Domain

Classification: `remove`

- `apps/server_core/internal/composition/root.go`: imports and wires `connectors/adapters/vtex/http`, reads VTEX env credentials, creates VTEX adapter, and injects it into `NewBatchOrchestrator`.
- `apps/server_core/internal/modules/connectors/transport/http_handler.go`: exposes `/connectors/vtex/publish`, `/connectors/vtex/publish/batch/*`, and `/connectors/vtex/validate-connection`; request/response fields use `vtex_account`; error codes include VTEX semantics.
- `apps/server_core/internal/modules/connectors/ports/vtex_catalog.go`: provider-specific `VTEXCatalogPort` and VTEX parameter/data structs.
- `apps/server_core/internal/modules/connectors/application/orchestrator.go`: orchestrates VTEX catalog pipeline and validation through `VTEXCatalogPort`.
- `apps/server_core/internal/modules/connectors/application/executor.go`: executes VTEX pipeline steps and records VTEX entity ids.
- `apps/server_core/internal/modules/connectors/domain/batch.go`, `mapping.go`, `pipeline.go`, `errors.go`: VTEX account/entity/error domain names.
- `apps/server_core/internal/modules/connectors/ports/repository.go`: repository methods and mappings tied to VTEX account/entity names.
- `apps/server_core/internal/modules/connectors/adapters/postgres/repository.go`: SQL for VTEX publication operations and `vtex_entity_mappings`.
- `apps/server_core/internal/modules/connectors/adapters/vtex/http/*`: full HTTP adapter, client, credentials, mapper, errors, config, tests, and integration test.
- `apps/server_core/internal/modules/connectors/adapters/vtex/stub/adapter.go`: VTEX stub adapter.

Mixed/generic note:

- `apps/server_core/internal/modules/integrations/adapters/feesync/marketplace_executor.go` and test contain `vtex` only as a disabled marketplace executor branch/fixture. Classification: `remove` if VTEX provider is removed from active provider catalog in F-02; otherwise convert to legacy/unavailable provider fixture only if the integration catalog intentionally preserves provider history.

### Backend Tests

Classification: `remove`

- `apps/server_core/tests/unit/router_registration_test.go`: expects VTEX routes and stubs `VTEXCatalogPort`.
- `apps/server_core/tests/unit/root_router_test.go`: requires VTEX env credentials for root router construction.
- `apps/server_core/tests/unit/connectors_handler_test.go`: tests `/connectors/vtex/*` handler behavior.
- `apps/server_core/tests/unit/connectors_executor_test.go`: tests VTEX publication pipeline and VTEX mapping/error paths.
- `apps/server_core/tests/unit/connectors_orchestrator_test.go`: tests VTEX batch orchestration.
- `apps/server_core/tests/unit/connectors_domain_test.go`: tests VTEX domain mappings/errors.
- `apps/server_core/tests/integration/vtex_adapter_test.go`: live/adapter tests using VTEX env.
- `apps/server_core/internal/modules/connectors/adapters/vtex/http/adapter_test.go` and `integration_test.go`: VTEX adapter tests.

Mixed/generic fixture updates:

- `apps/server_core/tests/unit/marketplaces_handler_test.go` and `apps/server_core/tests/integration/marketplaces_repository_test.go`: use `vtex` as marketplace account/policy fixture. Classification: `remove` VTEX fixture strings or replace with non-VTEX provider fixtures during F-02.
- `apps/server_core/tests/integration/phase1_smoke_test.go`: contains VTEX smoke reference. Classification: `remove` or replace with generic marketplace smoke fixture.

### OpenAPI Contract

Classification: `remove`

- `contracts/api/marketplace-central.openapi.yaml`: remove `/connectors/vtex/validate-connection`, `/connectors/vtex/publish`, `/connectors/vtex/publish/batch/{batch_id}`, `/connectors/vtex/publish/batch/{batch_id}/retry`, `PublishToVTEXRequest`, `RetryVTEXBatchRequest`, `PublishToVTEXResponse`, `validateVTEXConnection`, `publishToVTEX`, `getVTEXBatchStatus`, `retryVTEXBatch`, and `CONNECTORS_VTEX_*` active error enum entries.

### SDK Runtime

Classification: `remove`

- `packages/sdk-runtime/src/index.ts`: remove `VTEXProduct`, VTEX batch request/response fields where provider-specific, `publishToVTEX`, VTEX batch status route client methods, and `/connectors/vtex/*` URL strings.
- `packages/sdk-runtime/src/index.test.ts`: remove VTEX publish/batch tests and replace generic marketplace fixture values where they are not intentionally testing VTEX.

### Frontend Pages, Navigation, and Tests

Classification: `remove`

- `apps/web/src/app/AppRouter.tsx`: remove imports and routes for `VTEXPublishPage` and `BatchDetailPage` at `/connectors/vtex` and `/connectors/vtex/batch/:batchId`.
- `apps/web/src/app/Layout.tsx`: remove nav item `{ to: "/connectors/vtex", label: "VTEX Publisher" }`.
- `packages/feature-connectors/src/VTEXPublishPage.tsx`: remove VTEX publisher page.
- `packages/feature-connectors/src/VTEXPublishPage.test.tsx`: remove VTEX publisher tests.
- `packages/feature-connectors/src/BatchDetailPage.tsx`: remove VTEX batch detail page or replace only if a provider-neutral batch page is explicitly retained.
- `packages/feature-connectors/src/BatchDetailPage.test.tsx`: remove VTEX batch detail tests.
- `packages/feature-connectors/src/index.ts`: stop exporting VTEX pages.

Mixed/generic fixture updates:

- `packages/feature-marketplaces/src/components/MarketplaceIcon.tsx` and test: remove VTEX brand color if provider catalog no longer exposes active VTEX; retain only if explicitly marked legacy/unavailable.
- `packages/feature-marketplaces/src/AccountPanel.tsx`: placeholder `"My VTEX Store"` should be changed to a non-VTEX example.
- `packages/feature-marketplaces/src/AccountPanel.test.tsx`, `AccountCard.test.tsx`: replace VTEX fixtures.
- `packages/feature-integrations/src/IntegrationsHubPage.test.tsx`: replace or legacy-mark VTEX provider/installations fixtures.
- `packages/feature-classifications/src/ClassificationsPage.tsx`, `ClassificationsPage.test.tsx`, `packages/feature-products/src/ProductsPage.tsx`, `ProductsPage.test.tsx`: replace `"VTEX Ready"` placeholder/test fixture with provider-neutral wording.

### Docs, Wiki, Brain, and Historical Planning

Classification: `legacy-doc-retain`

- `ARCHITECTURE.md`: already marks VTEX legacy and not target architecture; retain.
- `.brain/decisions/005-mercado-livre-first-control-plane.md`: ADR-005 explains why VTEX is removed; retain.
- `.brain/system-pulse.md`, `.brain/roadmap.json`, `.brain/session-log.md`: retain/update only to reflect M-01 progress and avoid future VTEX work language.
- `docs/research/2026-07-06-mercado-livre-operating-cockpit.md`: retain as pivot research.
- `wiki/operations/environment-and-db.md`: already marks VTEX env keys legacy; retain/update after env removal.
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/**`: retain as current mission evidence.

Classification: `legacy-doc-retain` with explicit legacy/historical marking required in F-03

- `docs/vtex-integration-reference.md`
- `docs/superpowers/specs/*vtex*`, `docs/superpowers/plans/*vtex*`
- historical `docs/superpowers/*` files that mention VTEX as prior plan or scaffold
- `IMPLEMENTATION_PLAN.md`: historical truth only; ensure it is not used as current planning.

### Environment

Classification: `remove`

- `.env` contains `VTEX_APP_KEY`, `VTEX_APP_TOKEN`, and `VTEX_ACCOUNT`. F-02 should remove the runtime dependency and rotate/revoke externally if these were real credentials. Values must not be copied into artifacts.
- `apps/server_core/internal/modules/connectors/adapters/vtex/http/credentials.go` reads `VTEX_APP_KEY` and `VTEX_APP_TOKEN`; remove with adapter.
- Tests setting `VTEX_APP_KEY`, `VTEX_APP_TOKEN`, and `VTEX_ACCOUNT` should be deleted or rewritten.

### Migrations and Persisted Data

Classification: `migration-risk`

- `apps/server_core/migrations/0005_connectors.sql`: creates VTEX-specific columns and `vtex_entity_mappings`. Because migrations are forward-only, F-02 must not edit history casually. The safe path is to stop active code from depending on these structures and create a later forward migration/drop/rename only if current database state allows it.
- `apps/server_core/internal/modules/connectors/adapters/postgres/repository.go` and repository tests depend on the migration shape; remove with active VTEX pipeline or quarantine behind legacy data access if historical rows must remain readable.

## Edge Cases

- Some tests use `vtex` only as a marketplace fixture, not as VTEX adapter behavior. These still conflict with ADR-005 if they present VTEX as an active provider, but they should be replaced carefully to avoid deleting generic marketplace coverage.
- Historical docs are allowed to mention VTEX only when clearly marked as legacy/historical.
- `.env` has real-looking VTEX credentials; validation evidence must redact values and should recommend external rotation/revocation.
- Forward-only migrations mean removing code is simpler than removing existing database objects.

## Acceptance Criteria

- Criterion: Inventory covers backend, contracts, SDK, frontend, docs, tests, env, and migrations with category counts and file-level classifications.
- Traces to milestone criterion ID: M-01-C01
- Proven by: `validation.md` command `rg -l -i --hidden ... "vtex" apps/server_core/internal apps/server_core/tests apps/server_core/migrations`; `rg -l -i --hidden ... "vtex" contracts packages apps/web/src`; `rg -l -i --hidden ... "vtex" wiki docs ARCHITECTURE.md IMPLEMENTATION_PLAN.md README.md AGENTS.md .brain`; `.env*` redacted inspection.

- Criterion: Active VTEX route, SDK method, adapter, and frontend navigation owners are identified for F-02 removal.
- Traces to milestone criterion ID: M-01-C01
- Proven by: `validation.md` targeted command for `connectors/vtex|publishToVTEX|VTEXPublishPage|VTEXCatalogPort|adapters/vtex|VTEX_APP_|VTEX_ACCOUNT|vtex_entity_mappings|vtex_account`.

- Criterion: OpenAPI and SDK VTEX surfaces are explicitly listed for F-02 alignment.
- Traces to milestone criterion ID: M-01-C02
- Proven by: `validation.md` OpenAPI and SDK findings for `/connectors/vtex/*`, `publishToVTEX`, and related request/response schemas.

- Criterion: No production deletion occurs in F-01.
- Traces to milestone criterion ID: M-01-C01
- Proven by: `git status --short` and changed path list limited to F-01 artifacts.

## Handoff

- Current status: `spec_ready`
- Next owner: Feature Implementer
- Next action: Write `plan.md`, complete validation evidence, and hand off to Milestone Orchestrator.
- Required files/evidence: `feature.md`, `spec.md`, `plan.md`, `validation.md`, M-01 validation contract.
- Blockers or open decisions: none for F-01; F-02 must decide forward migration handling and provider catalog treatment for legacy VTEX.
