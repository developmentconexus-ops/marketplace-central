# M-06 F-03 Task 2 Brief - Realization-Aware Snapshot Calculation

Source plan: `docs/superpowers/plans/2026-07-09-m06-f03-order-realization.md`, Task 2.

## Goal

Consume the profitability-owned order realization contract from Task 1 so cancelled and unknown orders retain known inputs but never receive contribution or margin.

## Files

- Modify `apps/server_core/internal/modules/profitability/domain/input.go`
- Modify `apps/server_core/internal/modules/profitability/application/service.go`
- Modify `apps/server_core/internal/modules/profitability/application/service_test.go`
- Append `.superpowers/sdd/m06-f03-correction-report.md`

## Required Interfaces

- Consume `ports.OrderReader.ListOrders` and `domain.OrderRealizationState` from Task 1.
- Add `ProfitSnapshot.RealizationState`.
- Add quality `ProfitSnapshotNotRealized = "not_realized"`.
- Add flags `ProfitFlagOrderCancelled = "order_cancelled"` and `ProfitFlagOrderStateUnknown = "order_state_unknown"`.

## Exact Semantics

- Realized: retain all existing complete/incomplete/manual/negative-margin behavior.
- Not realized: keep known revenue, fee, cost, tax, freight, commission, and manual-adjustment inputs visible; contribution and margin are always `nil`; quality is `not_realized`; flags include `order_cancelled`; `negative_margin` is absent.
- Unknown: keep known inputs visible; contribution and margin are always `nil`; quality is `incomplete`; flags include `order_state_unknown`.
- Missing order IDs map to `unknown`; never infer realization from payments, timestamps, or monetary values.
- Missing monetary facts remain `nil`; no unknown becomes zero.

## TDD

1. Add explicit realized, cancelled, and unknown order facts to snapshot tests.
2. For cancelled, assert realization `not_realized`, quality `not_realized`, nil contribution/margin, `order_cancelled`, and absence of `negative_margin`.
3. For unknown, assert realization `unknown`, quality `incomplete`, nil contribution/margin, and `order_state_unknown`.
4. Assert known input amounts remain present in cancelled and unknown snapshots.
5. Run RED before production edits:

```powershell
cd apps/server_core
$env:GOCACHE="$PWD\.gocache"
go test ./internal/modules/profitability/application -run 'TestCalculateSnapshots' -count=1
```

6. Require `s.orders` in `CalculateSnapshots`, read up to 1000 order facts, build `map[string]OrderRealizationState`, and pass it through snapshot construction. Missing IDs use `unknown`.
7. Handle realization before missing-money/negative-margin calculation.
8. Run `gofmt` on changed Go files.
9. Run GREEN/regression:

```powershell
go test ./internal/modules/profitability/... ./internal/composition -count=1
```

## Limits

- Do not implement migration, PostgreSQL mapping/tests, OpenAPI, SDK, or UI; those are Tasks 3-4.
- Do not change refund/reversal accounting.
- Do not change F-02 manual-adjustment invariants.
- Do not import orders application/domain outside `profitability/adapters/orders`.
- Do not commit, stage, reset, revert, or clean in the shared worktree.
- Report deterministic unit/contract evidence accurately; it is not live provider/database evidence.

