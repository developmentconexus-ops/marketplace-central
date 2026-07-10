# M-06 F-03 Task 2 Implementation Report

## Status

`DONE`

No commit, staging, reset, revert, or cleanup operation was performed.

## Scope

Implemented realization-aware profitability snapshot calculation only. Migration, PostgreSQL mapping/tests, OpenAPI, SDK, UI, and refund/reversal accounting were not changed.

## TDD RED Evidence

Command (from `apps/server_core` with `GOCACHE=$PWD\.gocache`):

```powershell
go test ./internal/modules/profitability/application -run 'TestCalculateSnapshots' -count=1
```

Exit: `1`.

Observed output:

```text
service_test.go:433:78: item.RealizationState undefined
service_test.go:444:41: undefined: profitabilitydomain.ProfitSnapshotNotRealized
service_test.go:450:53: undefined: profitabilitydomain.ProfitFlagOrderCancelled
service_test.go:461:53: undefined: profitabilitydomain.ProfitFlagOrderStateUnknown
FAIL marketplace-central/apps/server_core/internal/modules/profitability/application [build failed]
```

The failure was feature-related: tests referenced the required realization field, quality, and flags before production code defined them. There were no syntax, fixture, or environment failures.

## Implementation

- Added `ProfitSnapshot.RealizationState`.
- Added `not_realized`, `order_cancelled`, and `order_state_unknown` snapshot constants.
- Required the existing profitability-owned `OrderReader` for snapshot calculation.
- Read up to 1000 profitability-owned order facts and built an order-ID-to-realization-state map.
- Normalized absent, empty, or unsupported states to `unknown`; no monetary or timestamp inference is used.
- Passed realization state through item and order snapshot construction.
- Applied realization before missing-money and negative-margin calculation.
- Preserved all known revenue, fee, cost, tax, freight, commission, and manual-adjustment values for cancelled and unknown snapshots.
- Preserved nullable monetary inputs; no missing value is converted to zero.
- Preserved realized complete, incomplete, manual-adjustment, and negative-margin behavior.

## Focused GREEN Evidence

Command:

```powershell
go test ./internal/modules/profitability/application -run 'TestCalculateSnapshots' -count=1
```

Exit: `0`.

```text
ok marketplace-central/apps/server_core/internal/modules/profitability/application 1.807s
```

The test verifies explicit realized, cancelled, and unknown facts plus a missing order ID. Cancelled snapshots are `not_realized`, have nil contribution/margin, include `order_cancelled`, omit `negative_margin`, and retain known inputs. Explicit unknown and missing-order snapshots are `unknown`/`incomplete`, have nil contribution/margin, include `order_state_unknown`, and retain known inputs.

## Fresh GREEN / Regression Evidence

Required command (from `apps/server_core` with `GOCACHE=$PWD\.gocache`):

```powershell
go test ./internal/modules/profitability/... ./internal/composition -count=1
```

Final fresh exit: `0`.

```text
?  marketplace-central/apps/server_core/internal/modules/profitability/adapters/internalread [no test files]
ok marketplace-central/apps/server_core/internal/modules/profitability/adapters/orders 1.065s
ok marketplace-central/apps/server_core/internal/modules/profitability/adapters/postgres 2.202s
ok marketplace-central/apps/server_core/internal/modules/profitability/application 0.503s
?  marketplace-central/apps/server_core/internal/modules/profitability/domain [no test files]
?  marketplace-central/apps/server_core/internal/modules/profitability/ports [no test files]
ok marketplace-central/apps/server_core/internal/modules/profitability/transport 1.470s
?  marketplace-central/apps/server_core/internal/composition [no test files]
```

This is deterministic unit/contract regression evidence using repository test doubles where applicable. It is not live provider or database evidence.

## Formatting And Boundary Evidence

`gofmt` was run on all three changed Go files. A final `gofmt -d` check returned no output and exit `0`.

Boundary command:

```powershell
rg -n 'internal/modules/orders' internal/modules/profitability/domain internal/modules/profitability/ports internal/modules/profitability/application
```

Result: empty output, exit `1` (no matches). Profitability domain, ports, and application do not import the orders module.

## Exact Changed Paths

- `apps/server_core/internal/modules/profitability/domain/input.go`
- `apps/server_core/internal/modules/profitability/application/service.go`
- `apps/server_core/internal/modules/profitability/application/service_test.go`
- `.superpowers/sdd/m06-f03-task-2-report.md`
- `.superpowers/sdd/m06-f03-correction-report.md`

## Self-Review

- Confirmed the service reads `OrderReader.ListOrders` once with limit 1000 before snapshot construction.
- Confirmed missing order IDs and invalid/empty states normalize only to `unknown`, never to realized.
- Confirmed cancelled and unknown branches clear contribution/margin and return before missing-money and negative-margin logic.
- Confirmed cancelled flags cannot contain `negative_margin` and known values remain on snapshots.
- Confirmed F-02 adjustment application remains unchanged and its values remain visible before realization gating.
- Confirmed no orders-module import was introduced outside the existing adapter boundary.
- Confirmed no out-of-scope persistence, contract, SDK, UI, migration, or refund behavior was implemented.

## Concerns

None within Task 2. Persistence and public-contract propagation of `RealizationState` remain intentionally deferred to Tasks 3-4.
