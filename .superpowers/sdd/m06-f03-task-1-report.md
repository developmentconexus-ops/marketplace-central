# M-06 F-03 Task 1 Report - Profitability-Owned Order Facts

## Status

`DONE`

## Implementation

- Added the exact profitability-owned `OrderRealizationState`, `OrderLinkQuality`, `OrderFact`, and `OrderItemFact` contract.
- Changed `ports.OrderReader.ListOrders(context.Context, string, int)` to return `[]domain.OrderFact`.
- Kept the orders application/domain dependency exclusively in `profitability/adapters/orders` and translated every contract field there.
- Applied the required realization policy: Mercado Livre `paid` is `realized`; `cancelled` and `canceled` are `not_realized`; blank, unsupported, and non-Mercado Livre states are `unknown`.
- Translated all orders link qualities into profitability-owned link qualities.
- Changed profitability application helpers to consume only profitability-owned facts. Application code does not inspect provider status strings.
- Did not implement snapshot realization, persistence, OpenAPI, SDK, UI, or any manual-adjustment behavior change.

## RED Evidence

Command run from `apps/server_core` before any production edit:

```powershell
$env:GOCACHE="$PWD\.gocache"
go test ./internal/modules/profitability/adapters/orders ./internal/modules/profitability/application -run 'TestOrderReader|TestImportMarginInputs' -count=1
```

Exit code: `1`.

Relevant expected compile output:

```text
internal\modules\profitability\application\service_test.go:14:31: undefined: profitabilitydomain.OrderFact
internal\modules\profitability\application\service_test.go:75:124: undefined: profitabilitydomain.OrderItemFact
internal\modules\profitability\application\service_test.go:76:118: undefined: profitabilitydomain.OrderLinkResolved
internal\modules\profitability\adapters\orders\order_reader_test.go:55:32: undefined: profitabilitydomain.OrderFact
internal\modules\profitability\adapters\orders\order_reader_test.go:58:42: undefined: profitabilitydomain.OrderRealizationRealized
internal\modules\profitability\adapters\orders\order_reader_test.go:83:36: undefined: profitabilitydomain.OrderRealizationState
FAIL marketplace-central/apps/server_core/internal/modules/profitability/adapters/orders [build failed]
FAIL marketplace-central/apps/server_core/internal/modules/profitability/application [build failed]
```

Explanation: the tests correctly demanded the new profitability-owned order contract and policy constants before they existed. The failure was caused by the missing feature, not a typo or environment failure.

## GREEN Evidence

Fresh command run from `apps/server_core`:

```powershell
$env:GOCACHE="$PWD\.gocache"
go test ./internal/modules/profitability/adapters/orders ./internal/modules/profitability/application -count=1
```

Exit code: `0`.

```text
ok  marketplace-central/apps/server_core/internal/modules/profitability/adapters/orders 1.406s
ok  marketplace-central/apps/server_core/internal/modules/profitability/application 0.443s
```

This is unit/contract evidence using test doubles; no live provider, database, credential, queue, or network integration was exercised or claimed.

## Gofmt Evidence

Ran `gofmt -w` over all six scoped Go files. A subsequent `gofmt -d` over the same paths produced no output between `GOFMT_DIFF_BEGIN` and `GOFMT_DIFF_END`, confirming canonical formatting.

## Boundary Search Evidence

Exact command run from `apps/server_core`:

```powershell
rg -n "internal/modules/orders" internal/modules/profitability/domain internal/modules/profitability/ports internal/modules/profitability/application
```

Exit code: `1`; output was empty. For `rg`, this is the expected success condition meaning no forbidden orders-module imports exist in profitability domain, ports, or application. A supplemental search found orders imports only in `profitability/adapters/orders` and its test.

## Exact Changed Paths

- `apps/server_core/internal/modules/profitability/domain/order_fact.go`
- `apps/server_core/internal/modules/profitability/ports/order_reader.go`
- `apps/server_core/internal/modules/profitability/adapters/orders/order_reader.go`
- `apps/server_core/internal/modules/profitability/adapters/orders/order_reader_test.go`
- `apps/server_core/internal/modules/profitability/application/service.go`
- `apps/server_core/internal/modules/profitability/application/service_test.go`
- `.superpowers/sdd/m06-f03-task-1-report.md`
- `.superpowers/sdd/m06-f03-correction-report.md`

## Self-Review

- Contract names, field types, and constant values match the Task 1 brief exactly.
- Adapter mapping covers every `OrderFact` and `OrderItemFact` field and propagates lister errors without translation loss.
- Realization policy is provider-specific and remains outside profitability application/domain/ports.
- Unknown or missing orders link qualities conservatively map to profitability `missing`.
- No Task 2 realization filtering/calculation was added, and F-02 manual-adjustment code was not changed.
- Scoped status review showed only the six intended Go paths; the shared worktree's unrelated changes were preserved.
- No files were staged and no commit was created.

## Concerns

- No implementation concern found within Task 1 scope.
- Validation is intentionally limited to deterministic unit/contract behavior; runtime integration remains outside this task.
