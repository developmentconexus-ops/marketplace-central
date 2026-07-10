# M-06 F-03 Task 3 Report - PostgreSQL, OpenAPI, And SDK Contract

## Status

`DONE_WITH_CONCERNS`

Implementation and deterministic contract verification are complete. Real PostgreSQL integration validation remains pending because `MC_DATABASE_URL` was absent in this session.

## Scope

- Owning module: `profitability`.
- Contract surfaces: PostgreSQL snapshot storage, OpenAPI schemas, and `@marketplace-central/sdk-runtime` types/fixtures.
- Data source: realization-aware `domain.ProfitSnapshot` values produced by Tasks 1-2.
- Side effects: replace/read profitability snapshot rows only; no provider write, UI work, refund behavior, or inference change.
- Verification path: PostgreSQL-gated repository integration test, profitability/composition Go suite, SDK Vitest suite, focused contract inspection.

## TDD RED Evidence

### PostgreSQL repository test

Command (from `apps/server_core`, with `GOCACHE=$PWD\.gocache`):

```text
go test ./internal/modules/profitability/adapters/postgres -run TestProfitSnapshotRealizationPersistence -count=1
```

Environment label and exact result:

```text
MC_DATABASE_URL_PRESENT=false
ok  	marketplace-central/apps/server_core/internal/modules/profitability/adapters/postgres	2.526s
```

This is skipped/compile evidence only. It is not PostgreSQL integration evidence and does not validate migration execution or repository round-trip behavior against a database.

### SDK contract test

The first sandboxed invocation failed before test collection because Vitest could not read the workspace config path (`Access is denied` / could not resolve `vitest.config.ts`). The same required command was rerun with approved filesystem access and produced the meaningful contract RED:

```text
npm run test --workspace @marketplace-central/sdk-runtime

src/index.test.ts (35 tests | 2 failed)
expected undefined to be 'realized'
expected undefined to be 'not_realized'
Test Files  1 failed (1)
Tests  2 failed | 33 passed (35)
```

The failures were expected: the SDK response fixtures did not yet carry required `realization_state` values.

## Minimal Implementation

- Added migration `0031` verbatim from the Task 3 brief.
- `Store.ReplaceSnapshots` now inserts the required `realization_state` value.
- `Store.ListSnapshots` now selects and scans `realization_state` in the matching ordinal position.
- OpenAPI now defines `OrderRealizationState = realized | not_realized | unknown`, requires `realization_state` on `ProfitabilityProfitSnapshot`, adds quality `not_realized`, and adds flags `order_cancelled` and `order_state_unknown`.
- SDK mirrors those exact strings and requires `realization_state` on `ProfitabilityProfitSnapshot`.
- SDK fixtures cover realized and not-realized snapshots; the cancelled fixture retains `null` contribution and margin.
- The PostgreSQL integration test covers realized, cancelled, and unknown rows, including realization, quality, flags, and nil contribution/margin readback.

## Migration Semantics

- Adding `realization_state text NOT NULL DEFAULT 'unknown'` backfills existing snapshot rows safely as `unknown` while establishing non-null storage.
- Dropping the default makes every subsequent repository insert provide realization explicitly.
- The `NOT VALID` check allows rollout without a separate historical validation scan while enforcing the exact three-value enum for new or updated rows.
- No monetary column or nullability behavior changes.

## Fresh GREEN Evidence

### Go

Command (from `apps/server_core`, with `GOCACHE=$PWD\.gocache`):

```text
go test ./internal/modules/profitability/... ./internal/composition -count=1
```

Result: exit `0`; adapters/orders, adapters/postgres, application, and transport passed; internalread/domain/ports/composition compiled with no test files. The postgres package result includes a skipped environment-gated integration test, not a live database pass.

Focused explicit skip evidence:

```text
=== RUN   TestProfitSnapshotRealizationPersistence
profit_snapshot_integration_test.go:17: MC_DATABASE_URL not set; real PostgreSQL profit snapshot realization validation pending
--- SKIP: TestProfitSnapshotRealizationPersistence (0.00s)
PASS
ok  	marketplace-central/apps/server_core/internal/modules/profitability/adapters/postgres	2.237s
```

### SDK

Command:

```text
npm run test --workspace @marketplace-central/sdk-runtime
```

Exact summary:

```text
✓ src/index.test.ts (35 tests) 15ms
Test Files  1 passed (1)
Tests  35 passed (35)
```

## Changed Paths

- `apps/server_core/migrations/0031_profitability_order_realization.sql`
- `apps/server_core/internal/modules/profitability/adapters/postgres/store.go`
- `apps/server_core/internal/modules/profitability/adapters/postgres/profit_snapshot_integration_test.go`
- `contracts/api/marketplace-central.openapi.yaml`
- `packages/sdk-runtime/src/index.ts`
- `packages/sdk-runtime/src/index.test.ts`
- `.superpowers/sdd/m06-f03-task-3-report.md`
- `.superpowers/sdd/m06-f03-correction-report.md`

## Formatting And Self-Review

- Ran `gofmt -w` on `store.go` and `profit_snapshot_integration_test.go`.
- Reviewed the repository persistence ordinals together: INSERT has 21 matching columns/placeholders/arguments; SELECT and Scan have 20 matching ordinals; `realization_state` is after `currency` in both relevant shapes.
- Reviewed the migration against the exact brief text.
- Reviewed OpenAPI and SDK enum spellings against domain strings; all use exactly `realized`, `not_realized`, `unknown`, `order_cancelled`, and `order_state_unknown`.
- Confirmed nullable monetary fields remain nullable and no zero/default conversion was introduced.
- Confirmed no Task 4 UI, refund/reversal, or realization inference changes were made.
- Existing dirty OpenAPI/SDK hunks were preserved; Task 3 edits were additive and localized.
- No files were staged or committed.

## Concern / Remaining Gate

No real database was available. The controller must apply migrations `0028` through `0031` to PostgreSQL 16 and rerun `TestProfitSnapshotRealizationPersistence` with `MC_DATABASE_URL` to establish actual migration and repository integration evidence.
