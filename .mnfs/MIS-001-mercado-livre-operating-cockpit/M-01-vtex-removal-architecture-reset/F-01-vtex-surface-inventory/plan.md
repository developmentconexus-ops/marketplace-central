# Feature Plan

```yaml
id: F-01
type: feature-plan
status: executed
owner: Feature Implementer
parent: F-01-vtex-surface-inventory
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
split_decision: single
split_reason: Inventory-only feature changes only three MNFS artifacts and uses read-only repository searches.
```

## Feature ID

F-01-vtex-surface-inventory

## Steps

1. Read required project context, mission artifacts, M-01 validation contract, and F-01 brief.
2. Run broad `rg` inventory over backend, contracts, packages, frontend app, docs/wiki/brain, migrations, and env files.
3. Run targeted `rg` inspection for active route, SDK, frontend, adapter, env, and migration identifiers.
4. Classify findings as `remove`, `legacy-doc-retain`, or `migration-risk`.
5. Write `spec.md` with the classification inventory and acceptance criteria mapped to M-01 criteria.
6. Write `validation.md` with exact commands, counts, redacted env evidence, changed paths, and F-02/F-03 handoff.

## Files Expected To Change

- Path: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-01-vtex-removal-architecture-reset/F-01-vtex-surface-inventory/spec.md`
- Reason: Record inventory requirements, classification, and acceptance criteria.

- Path: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-01-vtex-removal-architecture-reset/F-01-vtex-surface-inventory/plan.md`
- Reason: Record execution steps, verification mapping, and split decision.

- Path: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-01-vtex-removal-architecture-reset/F-01-vtex-surface-inventory/validation.md`
- Reason: Record commands run, results, evidence, risks, and handoff.

## Verification Commands

- Command: `rg -l -i --hidden -g '!/.git/**' -g '!node_modules/**' -g '!dist/**' -g '!coverage/**' "vtex" apps/server_core/internal apps/server_core/tests apps/server_core/migrations`
- Satisfies criterion ID: M-01-C01
- Expected result: Lists backend/internal/test/migration files containing VTEX references for classification.

- Command: `rg -l -i --hidden -g '!/.git/**' -g '!node_modules/**' -g '!dist/**' -g '!coverage/**' "vtex" contracts packages apps/web/src`
- Satisfies criterion ID: M-01-C01, M-01-C02
- Expected result: Lists OpenAPI, SDK, frontend page/nav/test files containing VTEX references for classification.

- Command: `rg -l -i --hidden -g '!/.git/**' -g '!node_modules/**' -g '!dist/**' -g '!coverage/**' "vtex" wiki docs ARCHITECTURE.md IMPLEMENTATION_PLAN.md README.md AGENTS.md .brain`
- Satisfies criterion ID: M-01-C01
- Expected result: Lists docs/wiki/brain/historical files containing VTEX references for `legacy-doc-retain` or update classification.

- Command: `Get-ChildItem -Force -File -Name '.env*' | ForEach-Object { Select-String -Path $_ -Pattern 'VTEX|vtex|MPC_.*VTEX' -CaseSensitive:$false }`
- Satisfies criterion ID: M-01-C01
- Expected result: Identifies VTEX env key names while validation redacts values.

- Command: `rg -n "connectors/vtex|publishToVTEX|VTEXPublishPage|VTEXCatalogPort|adapters/vtex|VTEX_APP_|VTEX_ACCOUNT|vtex_entity_mappings|vtex_account" apps packages contracts wiki docs ARCHITECTURE.md IMPLEMENTATION_PLAN.md .brain`
- Satisfies criterion ID: M-01-C01, M-01-C02
- Expected result: Pinpoints active route, SDK, adapter, env, and migration identifiers that must have F-02/F-03 owners.

- Command: `git status --short`
- Satisfies criterion ID: M-01-C01
- Expected result: Only MNFS artifacts are untracked/changed; no production deletion in F-01.

## QA Steps

- Step: Manually compare validation categories against F-01 brief categories: backend, contracts, SDK, frontend, docs, tests, env, migrations.
- Expected result: Every requested category is present in `validation.md`.

- Step: Manually check that secret values from `.env` are not copied into `spec.md`, `plan.md`, or `validation.md`.
- Expected result: Artifacts mention only key names and redaction.

## Rollback/Risk Notes

- Rollback is deleting the three F-01 generated artifacts only; no production code is touched.
- `.env` contains VTEX credential values; do not paste them into artifacts or final handoff.
- Migration `0005_connectors.sql` is a forward-only migration and should be treated as migration risk in F-02, not casually edited.
- Historical docs are noisy; F-03 should legacy-mark or archive rather than delete all history blindly.

## Handoff

- Current status: `executed`
- Next owner: Milestone Orchestrator
- Next action: Use inventory evidence to drive F-02 removal and F-03 truth alignment.
- Required files/evidence: `spec.md`, `plan.md`, `validation.md`, exact command results.
- Blockers or open decisions: none for F-01.
