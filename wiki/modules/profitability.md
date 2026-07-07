# Module: Profitability

Layer: business intelligence
Path: `apps/server_core/internal/modules/profitability/` (planned)

## Main Question It Answers

"How much margin did this Mercado Livre sale actually generate, and how trustworthy is that number?"

## What This Module Owns

| Entity | Purpose |
|---|---|
| `ProfitSnapshot` | Per-order and per-item revenue/cost/margin result |
| `CostInput` | Internal cost, provider fee, freight, tax, and manual adjustment inputs |
| `MarginQuality` | Completeness flags: complete, estimated, missing cost, missing freight, missing link |

## Rules

- Profitability is calculated server-side, never in React.
- Unknown values are explicit data-quality states, not zero defaults.
- Cost must use the internal cost provider and the sale date when available.
- Manual freight/fee/tax adjustments require audit notes.
- `internal_read` missing cost or tax inputs must remain nil with `missing_cost` / `missing_tax` flags so margin quality is incomplete, not silently zeroed.

## Initial Mercado Livre Scope

- Revenue from Mercado Livre order/payment values.
- Provider sale fee from order item or listing fee endpoint when simulating.
- Internal product cost from Sankhya/MetalShopping as-of sale date.
- Freight and extra costs manual until provider data is proven reliable.
