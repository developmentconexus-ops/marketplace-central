# M-06 F-03 Order Realization Design

## Decision

Profitability calculation must distinguish a realized sale from a cancelled or unknown order state. A cancelled order remains inspectable, but it is not treated as revenue-producing activity and cannot receive contribution or margin values.

## Ownership And Boundaries

- `orders` remains the source of persisted provider order facts.
- `profitability` owns the canonical calculation state:
  - `realized`
  - `not_realized`
  - `unknown`
- `profitability/ports.OrderReader` returns profitability-owned order facts rather than `orders/domain` types.
- `profitability/adapters/orders` is the only profitability package that imports orders application/domain types. It translates provider status into the canonical realization state.
- The calculation service never branches on Mercado Livre status strings.

## Initial Status Policy

For the Mercado Livre adapter:

- `paid` maps to `realized`.
- `cancelled` or `canceled` maps to `not_realized`.
- Any other or blank status maps to `unknown`.

Other providers map to `unknown` until their adapter defines an explicit policy. Unknown is not silently treated as realized.

## Calculation Semantics

### Realized

The existing contribution and margin rules apply. Missing inputs continue to produce `nil` contribution/margin and explicit quality flags.

### Not Realized

- Revenue, fee, cost, tax, freight, commission, and manual-adjustment inputs remain visible when known.
- `contribution_amount` and `margin_percent` are always `nil`.
- Snapshot quality is `not_realized`.
- Snapshot flags include `order_cancelled`.
- `negative_margin` is never emitted for a cancelled order.

### Unknown

- Known inputs remain visible.
- `contribution_amount` and `margin_percent` are `nil` because realization cannot be proven.
- Snapshot quality is `incomplete`.
- Snapshot flags include `order_state_unknown`.

## Persistence And Contracts

- `ProfitSnapshot` gains required `realization_state`.
- Profit snapshot quality adds `not_realized`.
- Profit snapshot flags add `order_cancelled` and `order_state_unknown`.
- A forward migration adds `realization_state` to existing snapshots as `unknown`, removes the write default, and installs a `NOT VALID` enum check so new writes are constrained without blocking on history.
- PostgreSQL repository, OpenAPI, SDK, and orders UI types remain aligned.

## Error Handling

- Snapshot calculation requires an order reader; missing configuration returns the existing structured profitability configuration error.
- Inputs whose order context cannot be found receive `unknown`, never `realized`.
- No provider state is guessed from payments, timestamps, or monetary amounts.

## Verification

- TDD unit coverage for realized, cancelled, unknown, incomplete, manual adjustment, and negative-margin cases.
- Adapter tests prove Mercado Livre status mapping and module-boundary isolation.
- PostgreSQL integration tests prove realization state, flags, and `nil` contribution/margin survive persistence and readback.
- OpenAPI/SDK tests prove the new enum and field contract.
- Live QA imports real Mercado Livre orders and confirms paid/cancelled semantics where the account exposes both.
- Full runtime QA must also prove one resolved product link with real Oracle cost/tax data before M-06 can pass; fake or unresolved-link evidence is insufficient.

## Non-Goals

- Refund/reversal accounting is not inferred in this correction.
- Unknown provider statuses are not forced into realized or cancelled.
- Historical invalid audit rows are not rewritten as part of this feature.
