# M-06 F-02 Correction Report

## Status

`DONE_WITH_CONCERNS`

## Assigned Failure

Manual adjustment actor metadata and append-only audit history are not fully enforced or proven.

## Backend RED/GREEN Evidence

- RED/GREEN scope: the previous worker correction handoff, represented by
  `m06-f02-correction-brief.md` and the profitability application and PostgreSQL
  test files, requires focused failing tests before the smallest backend change.
- The current workspace contains focused application coverage for actor identity,
  valid scope/category combinations, append-only same-timestamp adjustments, and
  PostgreSQL append-only readback/constraints.
- This contract-only stage did not execute Go tests and no prior worker command
  output was available in the handoff files. Those backend results are therefore
  recorded as prior-worker evidence, not re-observed command evidence.

## Controller PostgreSQL Evidence

- Applied `0028_profitability_margin_inputs.sql` to a disposable PostgreSQL 16 instance: `CREATE TABLE`, `CREATE INDEX`, `CREATE TABLE`.
- Applied `0030_profitability_manual_adjustment_invariants.sql`: `ALTER TABLE`.
- Command:
  - `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\.gocache'; $env:MC_DATABASE_URL='postgres://mpc:mpc-test@localhost:5436/marketplace_central?sslmode=disable'; go test -json -count=1 -timeout=5m ./internal/modules/profitability/...`
- Result:
  - Exit code `0` on 2026-07-09.
  - `TestManualAdjustmentsAppendOnlyReadbackAndConstraints` passed against PostgreSQL.
  - Freight and commission events at the same timestamp persisted as two distinct IDs with signed amounts and complete actor/reason/time readback.
  - Database constraints rejected empty actor type/id, reason, currency, invalid scope, and invalid category.

## Contract Evidence

- OpenAPI now requires `actor` on `CreateProfitabilityManualAdjustmentRequest`.
- The SDK request type requires `actor: ProfitabilityActor` while continuing to
  post the original request object unchanged.
- The SDK test compares the complete JSON request body, including
  `actor_type` and `actor_id`.

## Validation

- `npm run test --workspace @marketplace-central/sdk-runtime`
- Result: passed: 1 test file and 35 tests passed (`src/index.test.ts`, 28ms;
  total duration 5.08s).

## Changed Paths

- `apps/server_core/internal/modules/profitability/application/service.go`
- `apps/server_core/internal/modules/profitability/application/service_test.go`
- `apps/server_core/internal/modules/profitability/adapters/postgres/manual_adjustment_integration_test.go`
- `apps/server_core/migrations/0030_profitability_manual_adjustment_invariants.sql`
- `contracts/api/marketplace-central.openapi.yaml`
- `packages/sdk-runtime/src/index.ts`
- `packages/sdk-runtime/src/index.test.ts`
- `.superpowers/sdd/m06-f02-correction-report.md`

## Remaining Concern

- None within the assigned F-02 correction. Live application/runtime validation remains part of the later full M-06 gate.

## Fix Review Cycle

- Status: `DONE`; real PostgreSQL revalidation completed.
- Assigned findings addressed:
  - every `0030` check constraint now uses `NOT VALID`, preserving enforcement
    for new or updated rows without validating existing history during migration;
  - PostgreSQL constraint cases now include order scope with a provider item,
    order scope with a provider variation, and item scope without a provider item.
- Smallest-fix paths:
  - `apps/server_core/migrations/0030_profitability_manual_adjustment_invariants.sql`
  - `apps/server_core/internal/modules/profitability/adapters/postgres/manual_adjustment_integration_test.go`
  - `.superpowers/sdd/m06-f02-correction-report.md`
- Formatting: `gofmt -w apps/server_core/internal/modules/profitability/adapters/postgres/manual_adjustment_integration_test.go` exited `0`.
- Local validation command:
  - `$env:MC_DATABASE_URL=''; $env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\.gocache'; go test ./internal/modules/profitability/adapters/postgres -run TestManualAdjustmentsAppendOnlyReadbackAndConstraints -count=1 -v`
- Local validation result: exit code `0`; package compiled and the gated test
  reported `SKIP` because `MC_DATABASE_URL` was not set. This is not real-database
  evidence.
- Real-DB evidence:
  - created isolated database `f02_migration_test` in PostgreSQL 16;
  - applied `0028`, inserted one historical row with blank reason and actor fields (`INSERT 0 1`), then applied corrected `0030` successfully (`ALTER TABLE`);
  - readback confirmed the historical row remained and all seven constraints were installed with `convalidated = false`;
  - re-ran `go test -json -count=1 -timeout=5m ./internal/modules/profitability/...` against that database; exit code `0`;
  - the PostgreSQL test proved rejection of order scope with item/variation identity and item scope without provider item identity.
- Original review findings resolved in code: `Yes`; runtime revalidation: `Complete`.
- Handoff target: controller/QA Validator for the fresh PostgreSQL proof and
  review verdict refresh.
