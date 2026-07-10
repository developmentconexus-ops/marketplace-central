# F-03 Specification - Profit Snapshot Calculation

```yaml
id: M-06-F-03
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: M-06
created: 2026-07-09
updated: 2026-07-09
```

## Goal

Calculate persisted `ProfitSnapshot` rows per order and per item from the already assembled margin inputs, keeping every unknown value explicit and preserving operator trust in the quality result.

## Scope

This feature includes:

- `ProfitSnapshot` domain and persistence
- order-level and item-level calculation
- application of manual adjustments
- quality flags for missing revenue, fee, cost, freight, commission, and tax
- API/SDK surface to calculate and inspect profit snapshots

This feature does not include:

- profitability UI
- automatic provider freight recovery beyond current inputs
- speculative estimation that would hide missing data

## Calculation Rules

- Revenue uses the persisted `revenue` input amount.
- Provider fee uses the persisted `sale_fee` input amount.
- Cost uses `CUSSEMICM`-derived `cost` input only.
- Tax uses the sum of `tax_icms`, `tax_ipi`, `tax_pis`, and `tax_cofins` when present.
- Freight uses order-level/manual freight inputs and adjustments.
- Commission uses order-level/manual commission inputs and adjustments.
- Contribution value is:
  - `revenue - provider_fee - cost - tax_total - freight_total - commission_total + generic_adjustments`
- Margin percent is:
  - `contribution / revenue * 100`
  - only when revenue is known and non-zero

## Quality Rules

- Any missing required monetary component keeps the corresponding snapshot amount `nil` where appropriate or keeps the final margin/contribution `nil`.
- Missing values produce explicit flags; no missing component may become `0`.
- Required flags for this slice:
  - `missing_revenue`
  - `missing_sale_fee`
  - `missing_cost`
  - `missing_tax`
  - `missing_freight`
  - `missing_commission`
  - `missing_link`
  - `negative_margin`
- Snapshot quality summary:
  - `complete` when all required components are known
  - `incomplete` when any required component is missing
  - `negative_margin` when the calculation is complete and contribution is below zero

## Manual Adjustment Rules

- Freight adjustments affect order-level freight total.
- Commission adjustments affect order-level commission total.
- Cost adjustments affect item-level cost total when item-scoped, otherwise order-level adjustment total.
- Generic adjustments are additive and can be positive or negative.

## Persistence

Add first-class persistence for calculated snapshots:

- `profitability_profit_snapshots`

Snapshots are replaceable per installation/order/item scope on recalculation.

## API Surface

- `POST /profitability/profit-snapshots/calculate`
- `GET /profitability/profit-snapshots?installation_id=...`

## Validation Requirements

Tests must cover:

- complete margin
- missing cost
- missing link
- manual adjustment
- canceled order
- negative-margin order

Runtime evidence should prove:

- the real imported order can generate a persisted profit snapshot
- missing cost/tax/freight still leave an explicit incomplete quality
- manual freight adjustment changes the calculated totals on recalculation
