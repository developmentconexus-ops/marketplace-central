# F-04 Orders Margin UI Spec

## Goal

Deliver an operator-facing orders workspace that uses the persisted `orders` and `profitability` APIs to inspect order profitability, surface missing inputs, and apply manual adjustments without placing business math in React.

## Scope

- New web route for Mercado Livre order operations.
- Orders list backed by `GET /orders`.
- Profitability detail backed by persisted `profitability` endpoints.
- Manual adjustment flow that persists an operator adjustment and refreshes snapshots.
- UI states for loading, error, empty, incomplete margin, complete margin, negative margin, and manual adjustments.

## Route and Ownership

- Route: `/orders`
- Frontend package: `packages/feature-orders`
- Owning UI concept: orders is the operator entrypoint; profitability is rendered as order intelligence inside the same workflow.

## User Experience

The page is a two-pane operational workspace:

- Left pane lists orders for the selected installation.
- Top controls allow import orders, import margin inputs, and recalculate snapshots.
- Filters allow focusing on margin quality (`all`, `incomplete`, `negative_margin`, `complete`).
- Each order row shows status, item count, link quality, margin quality, contribution, and missing-input signals.
- Right pane shows the selected order with:
  - order profitability summary
  - item-level profitability rows
  - input breakdown grouped by order/item scope
  - manual adjustment history
  - manual adjustment form

## Data Rules

- React never calculates margin, contribution, or quality.
- React may aggregate presentation-only labels, counts, grouping, and severity ordering.
- Missing provider/internal facts remain visible as missing or incomplete, never coerced to zero.
- The detail view must minimize PII and avoid introducing buyer-identifying fields.

## Interaction Contract

1. Operator selects an installation.
2. Operator imports orders if needed.
3. Operator imports profitability inputs and recalculates snapshots.
4. Operator filters incomplete or negative-margin orders.
5. Operator opens an order, inspects inputs/flags, and appends an adjustment.
6. The page refreshes adjustments and snapshots after the write.

## Rendering Rules

- Order-level snapshot drives the list `margin quality` and `contribution` when present.
- If no order-level snapshot exists yet, the row renders `Not calculated`.
- Link quality is derived from order items using worst-state precedence:
  `conflict` > `unresolved` > `rejected` > `missing` > `resolved`
- Item detail uses item-level snapshots when present.
- Manual adjustments render newest first and show scope, category, amount, reason, actor, and timestamp.

## Validation Target

- Browser validation runs against the local app with the real backend already populated by M-06 live evidence.
- UI tests cover loading, error, empty, incomplete margin, complete margin, negative margin, and manual adjustment rendering/refresh behavior.
