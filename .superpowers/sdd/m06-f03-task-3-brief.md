# M-06 F-03 Task 3 Brief - PostgreSQL, OpenAPI, And SDK Contract

Source plan: `docs/superpowers/plans/2026-07-09-m06-f03-order-realization.md`, Task 3.

## Goal

Persist required order realization on every profitability snapshot and expose the exact same enum/field contract through OpenAPI and `packages/sdk-runtime`.

## Files

- Create `apps/server_core/migrations/0031_profitability_order_realization.sql`
- Modify `apps/server_core/internal/modules/profitability/adapters/postgres/store.go`
- Create `apps/server_core/internal/modules/profitability/adapters/postgres/profit_snapshot_integration_test.go`
- Modify `contracts/api/marketplace-central.openapi.yaml`
- Modify `packages/sdk-runtime/src/index.ts`
- Modify `packages/sdk-runtime/src/index.test.ts`
- Append `.superpowers/sdd/m06-f03-correction-report.md`

## Exact Contract

- Persist required `realization_state` on every `profitability_profit_snapshots` row.
- API/SDK type: `OrderRealizationState = "realized" | "not_realized" | "unknown"`.
- `ProfitabilityProfitSnapshot` requires `realization_state`.
- Profit snapshot quality includes `not_realized`.
- Profit snapshot flags include `order_cancelled` and `order_state_unknown`.
- OpenAPI, SDK, PostgreSQL, and domain strings must match exactly.
- Missing monetary values remain nullable; no unknown-to-zero/default.

## TDD

1. Add a PostgreSQL integration test gated by `MC_DATABASE_URL`. Persist realized, cancelled, and unknown snapshots, list them back, and assert realization, quality, flags, and nil contribution/margin.
2. Update SDK fixtures/tests to require `realization_state` and preserve `not_realized` plus `order_cancelled`.
3. Run RED before production edits:

```powershell
cd apps/server_core
$env:GOCACHE="$PWD\.gocache"
go test ./internal/modules/profitability/adapters/postgres -run TestProfitSnapshotRealizationPersistence -count=1
cd ../../..
npm run test --workspace @marketplace-central/sdk-runtime
```

If `MC_DATABASE_URL` is absent, the repository test may skip; record that accurately and ensure compile/contract RED is demonstrated without claiming real PostgreSQL evidence.

## Exact Migration

```sql
ALTER TABLE profitability_profit_snapshots
    ADD COLUMN realization_state text NOT NULL DEFAULT 'unknown';

ALTER TABLE profitability_profit_snapshots
    ALTER COLUMN realization_state DROP DEFAULT;

ALTER TABLE profitability_profit_snapshots
    ADD CONSTRAINT profitability_profit_snapshots_realization_state_valid
        CHECK (realization_state IN ('realized', 'not_realized', 'unknown')) NOT VALID;
```

## Implementation

- Include `realization_state` in `Store.ReplaceSnapshots` INSERT.
- Include `realization_state` in `ListSnapshots` SELECT and Scan.
- Add OpenAPI `OrderRealizationState`; require `realization_state` in `ProfitabilityProfitSnapshot`; add quality/flag enum values.
- Mirror exact types/strings and fixtures in SDK.
- Run `gofmt` on changed Go files.

## GREEN

```powershell
cd apps/server_core
$env:GOCACHE="$PWD\.gocache"
go test ./internal/modules/profitability/... ./internal/composition -count=1
cd ../../..
npm run test --workspace @marketplace-central/sdk-runtime
```

Controller will separately apply migrations `0028` through `0031` to real PostgreSQL 16 and rerun with `MC_DATABASE_URL` during the real gate.

## Limits

- Do not implement UI; that is Task 4.
- Do not change refund/reversal behavior or realization inference.
- Do not commit, stage, reset, revert, or clean.
- Fake/skip/compile evidence is not real PostgreSQL validation; label it accurately.
- Preserve every unrelated dirty-worktree change.
