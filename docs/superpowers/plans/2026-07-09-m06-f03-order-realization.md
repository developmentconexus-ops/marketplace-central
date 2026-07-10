# M-06 F-03 Order Realization Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Persist and expose canonical order realization so cancelled and unknown orders never receive fabricated contribution or margin.

**Architecture:** `profitability` owns the order-fact contract and canonical realization enum. The `profitability/adapters/orders` boundary maps orders-owned provider facts into `realized`, `not_realized`, or `unknown`; calculation consumes only profitability types, persists realization with every snapshot, and exposes the same contract through OpenAPI, SDK, and UI.

**Tech Stack:** Go 1.x, pgx v5/PostgreSQL 16, OpenAPI YAML, TypeScript, React/Vitest.

## Global Constraints

- `paid` maps to `realized`; `cancelled` and `canceled` map to `not_realized`; every other or blank status maps to `unknown`.
- Only `profitability/adapters/orders` may import `orders/application` or `orders/domain`.
- `not_realized` keeps known inputs visible but always sets contribution and margin to `nil`, quality to `not_realized`, and flag `order_cancelled`.
- `unknown` keeps known inputs visible but always sets contribution and margin to `nil`, quality to `incomplete`, and flag `order_state_unknown`.
- Unknown order state must never become `realized` by default or inference.
- Missing monetary values remain `nil`; no unknown value becomes zero.
- OpenAPI, SDK, PostgreSQL, and UI enum/field contracts change together.
- New migration is forward-safe: existing snapshots become `unknown`, the write default is removed, and the enum check is `NOT VALID`.
- Do not change refund/reversal accounting or infer realization from payments, timestamps, or money.

---

### Task 1: Profitability-Owned Order Facts And Adapter Policy

**Files:**
- Create: `apps/server_core/internal/modules/profitability/domain/order_fact.go`
- Modify: `apps/server_core/internal/modules/profitability/ports/order_reader.go`
- Modify: `apps/server_core/internal/modules/profitability/adapters/orders/order_reader.go`
- Create: `apps/server_core/internal/modules/profitability/adapters/orders/order_reader_test.go`
- Modify: `apps/server_core/internal/modules/profitability/application/service.go`
- Modify: `apps/server_core/internal/modules/profitability/application/service_test.go`

**Interfaces:**
- Produces `domain.OrderRealizationState` constants `realized`, `not_realized`, `unknown`.
- Produces `domain.OrderFact` and `domain.OrderItemFact` with only fields needed by input assembly and snapshot calculation.
- `ports.OrderReader.ListOrders(context.Context, string, int) ([]domain.OrderFact, error)`.

- [ ] **Step 1: Write failing adapter and boundary tests**

Add table-driven cases in `order_reader_test.go` for Mercado Livre `paid`, `cancelled`, `canceled`, blank, another status, and another provider. Assert the mapped profitability fact and all item/timestamp/link fields. Add a structural/import assertion that profitability domain, ports, and application do not expose orders-domain types.

```go
tests := []struct {
    provider, status string
    want domain.OrderRealizationState
}{
    {"mercado_livre", "paid", domain.OrderRealizationRealized},
    {"mercado_livre", "cancelled", domain.OrderRealizationNotRealized},
    {"mercado_livre", "canceled", domain.OrderRealizationNotRealized},
    {"mercado_livre", "", domain.OrderRealizationUnknown},
    {"mercado_livre", "confirmed", domain.OrderRealizationUnknown},
    {"other", "paid", domain.OrderRealizationUnknown},
}
```

- [ ] **Step 2: Run RED**

Run:

```powershell
cd apps/server_core
$env:GOCACHE="$PWD\.gocache"
go test ./internal/modules/profitability/adapters/orders ./internal/modules/profitability/application -run 'TestOrderReader|TestImportMarginInputs' -count=1
```

Expected: compile failure because profitability-owned order facts and realization constants do not exist.

- [ ] **Step 3: Implement the owned contract and mapping**

Create exact domain types:

```go
type OrderRealizationState string

const (
    OrderRealizationRealized    OrderRealizationState = "realized"
    OrderRealizationNotRealized OrderRealizationState = "not_realized"
    OrderRealizationUnknown     OrderRealizationState = "unknown"
)

type OrderFact struct {
    InstallationID string
    ProviderOrderID string
    RealizationState OrderRealizationState
    ProviderCreatedAt *time.Time
    ProviderClosedAt *time.Time
    ProviderUpdatedAt *time.Time
    FetchedAt time.Time
    Items []OrderItemFact
}

type OrderLinkQuality string

const (
    OrderLinkResolved   OrderLinkQuality = "resolved"
    OrderLinkRejected   OrderLinkQuality = "rejected"
    OrderLinkConflict   OrderLinkQuality = "conflict"
    OrderLinkUnresolved OrderLinkQuality = "unresolved"
    OrderLinkMissing    OrderLinkQuality = "missing"
)

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

Map all fields in the adapter, including an explicit one-to-one translation of orders link quality into `OrderLinkQuality`; keep provider status strings out of application code. Replace application helper signatures and test stubs with the owned types.

- [ ] **Step 4: Run GREEN and boundary search**

```powershell
go test ./internal/modules/profitability/adapters/orders ./internal/modules/profitability/application -count=1
rg -n "internal/modules/orders" internal/modules/profitability/domain internal/modules/profitability/ports internal/modules/profitability/application
```

Expected: tests pass; `rg` returns no matches.

- [ ] **Step 5: Record handoff**

Append RED/GREEN output and changed paths to `.superpowers/sdd/m06-f03-correction-report.md`. Do not commit in the shared worktree.

---

### Task 2: Realization-Aware Snapshot Calculation

**Files:**
- Modify: `apps/server_core/internal/modules/profitability/domain/input.go`
- Modify: `apps/server_core/internal/modules/profitability/application/service.go`
- Modify: `apps/server_core/internal/modules/profitability/application/service_test.go`

**Interfaces:**
- Consumes `ports.OrderReader.ListOrders` and `domain.OrderRealizationState` from Task 1.
- Produces `ProfitSnapshot.RealizationState`.
- Adds quality `ProfitSnapshotNotRealized` and flags `ProfitFlagOrderCancelled`, `ProfitFlagOrderStateUnknown`.

- [ ] **Step 1: Write failing calculation tests**

Add explicit realized, cancelled, and unknown order facts to the snapshot test. Assert:

```go
// cancelled
snapshot.RealizationState == domain.OrderRealizationNotRealized
snapshot.Quality == domain.ProfitSnapshotNotRealized
snapshot.ContributionAmount == nil
snapshot.MarginPercent == nil
contains(snapshot.Flags, domain.ProfitFlagOrderCancelled)
!contains(snapshot.Flags, domain.ProfitFlagNegativeMargin)

// unknown
snapshot.RealizationState == domain.OrderRealizationUnknown
snapshot.Quality == domain.ProfitSnapshotIncomplete
snapshot.ContributionAmount == nil
snapshot.MarginPercent == nil
contains(snapshot.Flags, domain.ProfitFlagOrderStateUnknown)
```

Known revenue/fee/cost/tax/freight/commission must remain present in both scenarios.

- [ ] **Step 2: Run RED**

```powershell
go test ./internal/modules/profitability/application -run 'TestCalculateSnapshots' -count=1
```

Expected: compile/assertion failure for missing realization state, quality, and flags.

- [ ] **Step 3: Thread order context through calculation**

Require `s.orders` in `CalculateSnapshots`, read up to 1000 order facts, build `map[string]OrderRealizationState`, and pass it to `buildProfitSnapshots`. Every accumulator receives the mapped state; missing order IDs use `unknown`.

In `finalizeSnapshot`, handle realization before missing-money/negative-margin calculation:

```go
switch s.RealizationState {
case profitabilitydomain.OrderRealizationNotRealized:
    s.ContributionAmount = nil
    s.MarginPercent = nil
    flags[profitabilitydomain.ProfitFlagOrderCancelled] = struct{}{}
    delete(flags, profitabilitydomain.ProfitFlagNegativeMargin)
    s.Quality = profitabilitydomain.ProfitSnapshotNotRealized
    s.Flags = flattenFlags(flags)
    return
case profitabilitydomain.OrderRealizationUnknown:
    s.ContributionAmount = nil
    s.MarginPercent = nil
    flags[profitabilitydomain.ProfitFlagOrderStateUnknown] = struct{}{}
    s.Quality = profitabilitydomain.ProfitSnapshotIncomplete
    s.Flags = flattenFlags(flags)
    return
}
```

- [ ] **Step 4: Run GREEN and full profitability regression**

```powershell
go test ./internal/modules/profitability/... ./internal/composition -count=1
```

Expected: all packages pass; existing complete/incomplete/manual/negative tests remain valid with explicit realized facts.

- [ ] **Step 5: Record handoff**

Append evidence to `.superpowers/sdd/m06-f03-correction-report.md`. Do not commit.

---

### Task 3: PostgreSQL, OpenAPI, And SDK Contract

**Files:**
- Create: `apps/server_core/migrations/0031_profitability_order_realization.sql`
- Modify: `apps/server_core/internal/modules/profitability/adapters/postgres/store.go`
- Create: `apps/server_core/internal/modules/profitability/adapters/postgres/profit_snapshot_integration_test.go`
- Modify: `contracts/api/marketplace-central.openapi.yaml`
- Modify: `packages/sdk-runtime/src/index.ts`
- Modify: `packages/sdk-runtime/src/index.test.ts`

**Interfaces:**
- Persists required `realization_state` on every `profitability_profit_snapshots` row.
- API/SDK type `OrderRealizationState = "realized" | "not_realized" | "unknown"`.

- [ ] **Step 1: Write failing repository and SDK tests**

The PostgreSQL test is gated by `MC_DATABASE_URL`. Persist realized, cancelled, and unknown snapshots; list them back and assert realization, quality, flags, and nil contribution/margin. The SDK test fixture must include `realization_state` and assert `not_realized` plus `order_cancelled` decode unchanged.

- [ ] **Step 2: Run RED**

```powershell
cd apps/server_core
$env:GOCACHE="$PWD\.gocache"
go test ./internal/modules/profitability/adapters/postgres -run TestProfitSnapshotRealizationPersistence -count=1
cd ../../..
npm run test --workspace @marketplace-central/sdk-runtime
```

Expected: Go compile/SQL-shape failure and TypeScript test/type failure because realization is absent.

- [ ] **Step 3: Implement migration and repository mapping**

Migration content:

```sql
ALTER TABLE profitability_profit_snapshots
    ADD COLUMN realization_state text NOT NULL DEFAULT 'unknown';

ALTER TABLE profitability_profit_snapshots
    ALTER COLUMN realization_state DROP DEFAULT;

ALTER TABLE profitability_profit_snapshots
    ADD CONSTRAINT profitability_profit_snapshots_realization_state_valid
        CHECK (realization_state IN ('realized', 'not_realized', 'unknown')) NOT VALID;
```

Include `realization_state` in `Store.ReplaceSnapshots` INSERT and `ListSnapshots` SELECT/Scan.

- [ ] **Step 4: Align OpenAPI and SDK**

Add `OrderRealizationState`, require `realization_state` in `ProfitabilityProfitSnapshot`, add `not_realized` quality, and add flags `order_cancelled` and `order_state_unknown`. Mirror exact strings in SDK types and fixtures.

- [ ] **Step 5: Run GREEN**

```powershell
cd apps/server_core
$env:GOCACHE="$PWD\.gocache"
go test ./internal/modules/profitability/... ./internal/composition -count=1
cd ../../..
npm run test --workspace @marketplace-central/sdk-runtime
```

Expected: all pass. Controller will separately apply `0028`, `0029`, `0030`, and `0031` to PostgreSQL 16 and rerun with `MC_DATABASE_URL`.

- [ ] **Step 6: Record handoff**

Append evidence and changed paths to `.superpowers/sdd/m06-f03-correction-report.md`. Do not commit.

---

### Task 4: Orders UI Semantics And Regression

**Files:**
- Modify: `packages/feature-orders/src/OrdersPage.tsx`
- Modify: `packages/feature-orders/src/OrdersPage.test.tsx`
- Modify if required by type integration: `apps/web/src/app/ClientContext.tsx`

**Interfaces:**
- Consumes SDK realization state, quality, and flags from Task 3.
- Keeps the F-02 manual-adjustment actor required in `OrdersClient`.

- [ ] **Step 1: Write failing UI tests**

Add cancelled and unknown fixtures. Assert cancelled orders render `Not realized`, `Order cancelled`, and an em dash for contribution/margin; unknown orders render `Incomplete` plus `Order state unknown`. Assert no negative-margin label for cancelled orders. Update the local `OrdersClient` actor field from optional to required.

- [ ] **Step 2: Run RED**

```powershell
npm run test --workspace @marketplace-central/feature-orders -- OrdersPage.test.tsx
```

Expected: failing labels/filters before UI support.

- [ ] **Step 3: Implement semantic presentation**

Add `not_realized` to `QualityFilter`, filter options, summary, `qualityLabel`, and `qualityTone`. Render realization flags as operational state, not as “Missing inputs”; use `Data quality` for missing flags and `Order not realized` for cancellation. Do not calculate money in React.

- [ ] **Step 4: Run GREEN, web regression, and build**

```powershell
npm run test --workspace @marketplace-central/feature-orders -- OrdersPage.test.tsx
npm run test --workspace @marketplace-central/web -- AppRouter.test.tsx ClientContext.test.tsx viteProxy.test.ts
npm run build --workspace @marketplace-central/web
```

Expected: all tests and build pass.

- [ ] **Step 5: Record handoff**

Append evidence and changed paths to `.superpowers/sdd/m06-f03-correction-report.md`. Do not commit.

---

## Final Gate

After all four task reviews are clean, the controller must:

1. Apply migrations through `0031` to PostgreSQL 16 and run full Go/SDK/web checks with fresh exit codes.
2. Drive `/orders` in the built-in browser at desktop and mobile widths.
3. Import real Mercado Livre orders and verify live paid/cancelled behavior when available.
4. Resolve or use one real product link and prove Oracle `CUSSEMICM` plus tax inputs reach a complete realized snapshot.
5. Run the full independent M-06 milestone gate; no spot-check can produce a pass verdict.
