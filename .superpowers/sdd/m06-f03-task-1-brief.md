# M-06 F-03 Task 1 Brief - Profitability-Owned Order Facts

Source plan: `docs/superpowers/plans/2026-07-09-m06-f03-order-realization.md`, Task 1.

## Goal

Remove orders-domain types from profitability domain/ports/application and translate provider order facts into a profitability-owned contract at `profitability/adapters/orders`.

## Files

- Create `apps/server_core/internal/modules/profitability/domain/order_fact.go`
- Modify `apps/server_core/internal/modules/profitability/ports/order_reader.go`
- Modify `apps/server_core/internal/modules/profitability/adapters/orders/order_reader.go`
- Create `apps/server_core/internal/modules/profitability/adapters/orders/order_reader_test.go`
- Modify `apps/server_core/internal/modules/profitability/application/service.go`
- Modify `apps/server_core/internal/modules/profitability/application/service_test.go`
- Append `.superpowers/sdd/m06-f03-correction-report.md`

## Exact Contract

```go
type OrderRealizationState string

const (
    OrderRealizationRealized    OrderRealizationState = "realized"
    OrderRealizationNotRealized OrderRealizationState = "not_realized"
    OrderRealizationUnknown     OrderRealizationState = "unknown"
)

type OrderLinkQuality string

const (
    OrderLinkResolved   OrderLinkQuality = "resolved"
    OrderLinkRejected   OrderLinkQuality = "rejected"
    OrderLinkConflict   OrderLinkQuality = "conflict"
    OrderLinkUnresolved OrderLinkQuality = "unresolved"
    OrderLinkMissing    OrderLinkQuality = "missing"
)

type OrderFact struct {
    InstallationID      string
    ProviderOrderID     string
    RealizationState    OrderRealizationState
    ProviderCreatedAt   *time.Time
    ProviderClosedAt    *time.Time
    ProviderUpdatedAt   *time.Time
    FetchedAt            time.Time
    Items                []OrderItemFact
}

type OrderItemFact struct {
    ProviderItemID      string
    ProviderVariationID string
    Quantity            int
    UnitPrice           *float64
    SaleFeeAmount       *float64
    LinkQuality         OrderLinkQuality
    InternalProductID   *int
}
```

`ports.OrderReader.ListOrders(context.Context, string, int)` returns `[]domain.OrderFact`.

## Status Policy

- Mercado Livre `paid` -> `realized`.
- Mercado Livre `cancelled` or `canceled` -> `not_realized`.
- Blank, any other Mercado Livre state, or any other provider -> `unknown`.
- Application code must not inspect provider status strings.

## TDD

1. Add adapter table tests for all status-policy cases plus complete field mapping.
2. Update application tests to demand profitability-owned facts before production edits.
3. Run RED and record the expected compile failure.
4. Implement the contract, adapter mapping, port, and application signature changes.
5. Run:

```powershell
cd apps/server_core
$env:GOCACHE="$PWD\.gocache"
go test ./internal/modules/profitability/adapters/orders ./internal/modules/profitability/application -count=1
rg -n "internal/modules/orders" internal/modules/profitability/domain internal/modules/profitability/ports internal/modules/profitability/application
```

Expected: tests pass and `rg` finds no imports.

## Limits

- Do not implement snapshot realization behavior, persistence, OpenAPI, SDK, or UI; those are later tasks.
- Do not change manual-adjustment behavior from F-02.
- Do not commit in the shared worktree.
- Preserve unrelated changes and report exact paths and tests.
