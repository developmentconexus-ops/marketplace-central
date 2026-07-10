# M-06 F-03 Task 3 Review Package

## Baseline

- Branch: `main`; HEAD remains `8dba7db`.
- No commit or staging operation was performed.
- The profitability module is untracked and OpenAPI/SDK already had unrelated dirty hunks before Task 3, so no clean conventional task-only Git diff exists.

## Task 3 Review Surface

Read these task-specific surfaces once:

1. Full migration: `apps/server_core/migrations/0031_profitability_order_realization.sql`
2. `Store.ReplaceSnapshots` and `Store.ListSnapshots` in `apps/server_core/internal/modules/profitability/adapters/postgres/store.go`
3. Full test: `apps/server_core/internal/modules/profitability/adapters/postgres/profit_snapshot_integration_test.go`
4. `OrderRealizationState`, `ProfitSnapshotQuality`, `ProfitSnapshotFlag`, and `ProfitabilityProfitSnapshot` schemas around lines 2666-2705 of `contracts/api/marketplace-central.openapi.yaml`
5. SDK realization/quality/flag/snapshot types around lines 380-405 of `packages/sdk-runtime/src/index.ts`
6. SDK profitability fixtures/assertions around lines 1100-1190 of `packages/sdk-runtime/src/index.test.ts`

Inspect outside these surfaces only for a concrete named ordinal/contract risk.

## Evidence

- Requirements: `.superpowers/sdd/m06-f03-task-3-brief.md`
- Implementer report: `.superpowers/sdd/m06-f03-task-3-report.md`
- Aggregate evidence: `.superpowers/sdd/m06-f03-correction-report.md`

## Validation Level

- Go profitability/composition regression passed, but the new PostgreSQL integration test skipped because `MC_DATABASE_URL` was absent.
- SDK contract RED produced two missing-realization failures; GREEN passed 35/35.
- Real migration/repository validation remains explicitly pending for the controller's PostgreSQL 16 gate and must not be treated as passed here.
