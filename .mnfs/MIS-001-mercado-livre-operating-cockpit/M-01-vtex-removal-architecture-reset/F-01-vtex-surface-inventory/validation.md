# Feature Validation

```yaml
id: F-01
type: feature-validation
status: quick_validation_passed
owner: Feature Implementer
parent: F-01-vtex-surface-inventory
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-01-vtex-surface-inventory

## Summary

Quick validation passed for the inventory-only feature. VTEX references were inventoried across backend, contracts, SDK, frontend, docs, tests, env, and migrations; each group is classified as `remove`, `legacy-doc-retain`, or `migration-risk`. No production code was removed.

## Quick Validation Result

- Result: Pass
- Result owner: Feature Implementer
- Decision date: 2026-07-06
- Final feature state for handoff: quick_validation_passed

## Quick Validation State

- fixup_attempts: 0
- max_fixup_attempts: 1
- last_feature_validation_result: Pass

## Spec Adherence

- Spec satisfied: Yes
- Deviations: None.
- Reason: F-01 required inventory and classification only; generated artifacts satisfy the brief and no removal was performed.

## Changes Made

- File: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-01-vtex-removal-architecture-reset/F-01-vtex-surface-inventory/spec.md`
- Change: Added feature spec with VTEX inventory classification and M-01 acceptance mapping.

- File: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-01-vtex-removal-architecture-reset/F-01-vtex-surface-inventory/plan.md`
- Change: Added execution plan, verification commands, and `split_decision: single`.

- File: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-01-vtex-removal-architecture-reset/F-01-vtex-surface-inventory/validation.md`
- Change: Added quick validation evidence, categorized counts, risks, and handoff.

## Categorized Inventory

### Backend Routes, Adapters, Ports, Tests

- Classification: `remove`
- Files with backend/internal/test/migration VTEX matches: 34
- Match count in backend/internal/test/migration scope: 682
- Active route owner for F-02: `apps/server_core/internal/modules/connectors/transport/http_handler.go`
- Adapter/wiring owner for F-02: `apps/server_core/internal/composition/root.go`, `apps/server_core/internal/modules/connectors/adapters/vtex/**`
- Port/application/domain owner for F-02: `apps/server_core/internal/modules/connectors/ports/vtex_catalog.go`, connector orchestrator/executor/domain/repository surfaces.
- Test owner for F-02: connector route/orchestrator/executor/domain/router/root tests plus VTEX adapter integration tests.

Files:

- `apps/server_core/internal/composition/root.go`
- `apps/server_core/internal/modules/connectors/transport/http_handler.go`
- `apps/server_core/internal/modules/connectors/ports/vtex_catalog.go`
- `apps/server_core/internal/modules/connectors/ports/repository.go`
- `apps/server_core/internal/modules/connectors/application/orchestrator.go`
- `apps/server_core/internal/modules/connectors/application/executor.go`
- `apps/server_core/internal/modules/connectors/domain/batch.go`
- `apps/server_core/internal/modules/connectors/domain/errors.go`
- `apps/server_core/internal/modules/connectors/domain/mapping.go`
- `apps/server_core/internal/modules/connectors/domain/pipeline.go`
- `apps/server_core/internal/modules/connectors/adapters/postgres/repository.go`
- `apps/server_core/internal/modules/connectors/adapters/vtex/http/adapter.go`
- `apps/server_core/internal/modules/connectors/adapters/vtex/http/adapter_test.go`
- `apps/server_core/internal/modules/connectors/adapters/vtex/http/client.go`
- `apps/server_core/internal/modules/connectors/adapters/vtex/http/config.go`
- `apps/server_core/internal/modules/connectors/adapters/vtex/http/credentials.go`
- `apps/server_core/internal/modules/connectors/adapters/vtex/http/errors.go`
- `apps/server_core/internal/modules/connectors/adapters/vtex/http/integration_test.go`
- `apps/server_core/internal/modules/connectors/adapters/vtex/http/mapper.go`
- `apps/server_core/internal/modules/connectors/adapters/vtex/http/types.go`
- `apps/server_core/internal/modules/connectors/adapters/vtex/stub/adapter.go`
- `apps/server_core/internal/modules/integrations/adapters/feesync/marketplace_executor.go`
- `apps/server_core/internal/modules/integrations/adapters/feesync/marketplace_executor_test.go`
- `apps/server_core/tests/unit/router_registration_test.go`
- `apps/server_core/tests/unit/root_router_test.go`
- `apps/server_core/tests/unit/marketplaces_handler_test.go`
- `apps/server_core/tests/unit/connectors_handler_test.go`
- `apps/server_core/tests/unit/connectors_executor_test.go`
- `apps/server_core/tests/unit/connectors_orchestrator_test.go`
- `apps/server_core/tests/unit/connectors_domain_test.go`
- `apps/server_core/tests/integration/marketplaces_repository_test.go`
- `apps/server_core/tests/integration/vtex_adapter_test.go`
- `apps/server_core/tests/integration/phase1_smoke_test.go`

### OpenAPI

- Classification: `remove`
- Files: 1
- Match count: 32
- F-02 owner: `contracts/api/marketplace-central.openapi.yaml`
- Active surfaces: `/connectors/vtex/validate-connection`, `/connectors/vtex/publish`, `/connectors/vtex/publish/batch/{batch_id}`, `/connectors/vtex/publish/batch/{batch_id}/retry`, `PublishToVTEXRequest`, `RetryVTEXBatchRequest`, `PublishToVTEXResponse`, `CONNECTORS_VTEX_*`.

### SDK Types, Methods, Tests

- Classification: `remove`
- Files: 2
- Match count: 30
- F-02 owner: `packages/sdk-runtime/src/index.ts`, `packages/sdk-runtime/src/index.test.ts`
- Active surfaces: `VTEXProduct`, `PublishBatchRequest`, `PublishBatchResponse`, `BatchStatus` fields tied to `vtex_account`, `publishToVTEX`, `getBatchStatus`, `retryBatch`, `/connectors/vtex/*`.

### Frontend Pages, Navigation, Tests

- Classification: `remove`
- Files with package/app frontend VTEX matches excluding SDK/OpenAPI: 18
- Match count in contracts/packages/apps/web scope including SDK/OpenAPI: 245
- F-02 route/nav owners: `apps/web/src/app/AppRouter.tsx`, `apps/web/src/app/Layout.tsx`
- F-02 page owners: `packages/feature-connectors/src/VTEXPublishPage.tsx`, `packages/feature-connectors/src/BatchDetailPage.tsx`, tests, and export.
- Mixed fixture owners: marketplace icon/account/integrations/classifications/products tests and placeholders should replace VTEX examples with provider-neutral or Mercado Livre-safe examples.

Files:

- `apps/web/src/app/AppRouter.tsx`
- `apps/web/src/app/Layout.tsx`
- `packages/feature-connectors/src/VTEXPublishPage.tsx`
- `packages/feature-connectors/src/VTEXPublishPage.test.tsx`
- `packages/feature-connectors/src/BatchDetailPage.tsx`
- `packages/feature-connectors/src/BatchDetailPage.test.tsx`
- `packages/feature-connectors/src/index.ts`
- `packages/feature-marketplaces/src/AccountPanel.tsx`
- `packages/feature-marketplaces/src/AccountPanel.test.tsx`
- `packages/feature-marketplaces/src/AccountCard.test.tsx`
- `packages/feature-marketplaces/src/components/MarketplaceIcon.tsx`
- `packages/feature-marketplaces/src/components/MarketplaceIcon.test.tsx`
- `packages/feature-integrations/src/IntegrationsHubPage.test.tsx`
- `packages/feature-classifications/src/ClassificationsPage.tsx`
- `packages/feature-classifications/src/ClassificationsPage.test.tsx`
- `packages/feature-products/src/ProductsPage.tsx`
- `packages/feature-products/src/ProductsPage.test.tsx`

### Docs, Wiki, Brain, Historical Planning

- Classification: `legacy-doc-retain`
- Files: 38
- Match count: 1,333
- F-03 owner: legacy-mark or update historical docs so no current planning path instructs future VTEX work.
- Retain as architecture truth: `ARCHITECTURE.md`, `.brain/decisions/005-mercado-livre-first-control-plane.md`, `.mnfs/MIS-001-mercado-livre-operating-cockpit/**`.
- Update/legacy-mark likely needed: `docs/vtex-integration-reference.md`, old `docs/superpowers/plans/*vtex*`, old `docs/superpowers/specs/*vtex*`, `IMPLEMENTATION_PLAN.md`, `.brain/roadmap.json`, `.brain/system-pulse.md`, `wiki/operations/environment-and-db.md`.

### Environment

- Classification: `remove`
- Files: `.env`
- Evidence: `.env` contains the key names `VTEX_APP_KEY`, `VTEX_APP_TOKEN`, and `VTEX_ACCOUNT`.
- Secret handling: values were observed locally but intentionally redacted from this artifact.
- F-02/F-03 owner: remove runtime dependency and update env docs; rotate/revoke external credentials if these values are live.

### Migrations

- Classification: `migration-risk`
- Files: `apps/server_core/migrations/0005_connectors.sql`
- Match count: 10
- Risk: creates `vtex_account`, `vtex_entity_id`, and `vtex_entity_mappings`; forward-only migration history should not be edited blindly.
- F-02/Future owner: remove active dependencies first, then decide whether a new forward migration should drop/quarantine legacy VTEX tables/columns after data state is understood.

## Commands Run

- Command: `rg -l -i --hidden -g '!/.git/**' -g '!node_modules/**' -g '!dist/**' -g '!coverage/**' "vtex" apps/server_core/internal apps/server_core/tests apps/server_core/migrations`
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Expected: Lists backend/internal/test/migration VTEX files.
- Actual: 34 files listed.
- Artifact: this `validation.md` categorized inventory.
- Blocking condition: none.

- Command: `rg -l -i --hidden -g '!/.git/**' -g '!node_modules/**' -g '!dist/**' -g '!coverage/**' "vtex" contracts packages apps/web/src`
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Expected: Lists OpenAPI, SDK, frontend page/nav/test VTEX files.
- Actual: 20 files listed.
- Artifact: this `validation.md` categorized inventory.
- Blocking condition: none.

- Command: `rg -l -i --hidden -g '!/.git/**' -g '!node_modules/**' -g '!dist/**' -g '!coverage/**' "vtex" wiki docs ARCHITECTURE.md IMPLEMENTATION_PLAN.md README.md AGENTS.md .brain`
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Expected: Lists docs/wiki/brain/historical files containing VTEX references.
- Actual: 38 files listed.
- Artifact: this `validation.md` categorized inventory.
- Blocking condition: none.

- Command: `Get-ChildItem -Force -File -Name '.env*' | ForEach-Object { Select-String -Path $_ -Pattern 'VTEX|vtex|MPC_.*VTEX' -CaseSensitive:$false }`
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Expected: Identifies VTEX env key names.
- Actual: `.env` contains `VTEX_APP_KEY`, `VTEX_APP_TOKEN`, and `VTEX_ACCOUNT`; values redacted.
- Artifact: this `validation.md` environment section.
- Blocking condition: none.

- Command: `rg -n "connectors/vtex|publishToVTEX|VTEXPublishPage|VTEXCatalogPort|adapters/vtex|VTEX_APP_|VTEX_ACCOUNT|vtex_entity_mappings|vtex_account" apps packages contracts wiki docs ARCHITECTURE.md IMPLEMENTATION_PLAN.md .brain`
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Expected: Pinpoints active route, SDK, adapter, env, and migration identifiers.
- Actual: Identified active backend routes, OpenAPI paths/schemas, SDK methods, frontend routes/nav/pages, adapter imports, env keys, and migration objects. Output was large; summarized by category above.
- Artifact: this `validation.md` categorized inventory.
- Blocking condition: none.

- Command: `rg --files apps/server_core/internal/modules/connectors/adapters/vtex packages/feature-connectors | sort`
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Expected: Lists VTEX adapter files and feature connector package files.
- Actual: Listed VTEX HTTP/stub adapter files and `packages/feature-connectors` page/test/export files.
- Artifact: this `validation.md` backend/frontend sections.
- Blocking condition: none.

- Command: `git status --short`
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Expected: F-01 should not delete production code.
- Actual: `.mnfs/` is untracked; no production code changes were made in F-01.
- Artifact: this `validation.md` changes made section.
- Blocking condition: none.

## Manual QA

- QA level: QA-0
- Flow or step: Checked F-01 brief categories against validation inventory.
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Expected: backend, contracts, SDK, frontend, docs, tests, env, and migrations are all named.
- Actual: All requested categories are present.
- Blocking condition: none.

- QA level: QA-0
- Flow or step: Checked artifacts for copied VTEX secret values.
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Expected: secret values are not present in `spec.md`, `plan.md`, or `validation.md`.
- Actual: artifacts mention only env key names and redaction.
- Blocking condition: none.

## Evidence

- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-01-vtex-removal-architecture-reset/F-01-vtex-surface-inventory/spec.md`
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Blocking condition: none.

- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-01-vtex-removal-architecture-reset/F-01-vtex-surface-inventory/plan.md`
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Blocking condition: none.

- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-01-vtex-removal-architecture-reset/F-01-vtex-surface-inventory/validation.md`
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Blocking condition: none.

## Risks

- `.env` contains real-looking VTEX credentials; do not publish or paste values, and rotate/revoke externally if live.
- Removing active VTEX code may expose generic connector abstractions that were actually VTEX-shaped; F-02 should remove provider-specific naming without weakening future capability-port design.
- `0005_connectors.sql` is migration risk; forward-only migration policy means removal of persisted VTEX columns/tables needs a new migration decision.
- Historical docs are noisy; F-03 should retain history only with explicit legacy labels and ensure current architecture docs do not instruct future VTEX work.

## Handoff

- Current status: `quick_validation_passed`
- Next owner: Milestone Orchestrator
- Next action: Review F-01 inventory and dispatch F-02 VTEX active surface removal.
- Required files/evidence: `feature.md`, `spec.md`, `plan.md`, `validation.md`.
- Blockers or open decisions: none for F-01; F-02 must decide exact migration handling for `0005_connectors.sql` and whether integration catalog keeps a legacy/unavailable VTEX provider entry.
