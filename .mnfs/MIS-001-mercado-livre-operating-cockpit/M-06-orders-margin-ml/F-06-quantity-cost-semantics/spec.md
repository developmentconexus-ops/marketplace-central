# F-06 Quantity Cost Semantics — Specification

```yaml
id: F-06
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: M-06
created: 2026-07-12
updated: 2026-07-12
validation_level: QA-0
lifecycle_scope: feature
```

## Problem

Profitability imports persisted Oracle `CUSSEMICM` unchanged even though the
signed M-06 owner decision defines it as a per-unit cost. Item revenue is
already extended by Mercado Livre quantity, so product cost must be extended
to the same item-line scope without changing fee, tax, link, tenant, or
unknown-value behavior.

## Requirements

- Treat `GetCostAsOf` / Oracle `CUSSEMICM` as a per-unit amount.
- Persist the profitability item `cost` input as `CUSSEMICM × item quantity`.
- Preserve `nil` cost as `nil`, together with its existing missing/quality
  reason; it must never become `0`.
- Keep Mercado Livre `sale_fee` and every Oracle tax component as already
  item-line totals; do not quantity-multiply them.
- Carry durable source-reference evidence of the unit CUSSEMICM source and
  the resulting line cost; retain line-total source references for sale fee
  and tax components.

## Non-Goals

- No tax estimation, manual adjustment, authentication/authorization, API,
  SDK, provider, Oracle write, linkage, tenant-scope, or runtime-data change.
- No change to the numeric amount supplied by Mercado Livre sale fees or
  Oracle tax components.

## Design

`buildItemInputs` supplies the Mercado Livre item quantity to `mapCostInput`.
The mapper extends a non-nil CUSSEMICM amount once and records a source
reference that distinguishes `cussemicm_per_unit` from the persisted
`cost_line_total`. Sale-fee and tax mappers use durable line-total source
references and retain their incoming amounts. Nil propagation occurs before
any multiplication.

## Edge Cases

- Quantity `1` preserves the numeric CUSSEMICM amount.
- Quantities `2` and `7` extend only cost.
- Nil/missing cost remains nil with `missing` quality and a missing-cost
  reason.
- Missing, partial, or unresolved tax/link states remain their current
  quality states and never create zero-valued monetary inputs.

## Acceptance Criteria

- Quantity 1, 2, and 7 produce item-line costs of unit CUSSEMICM multiplied
  by quantity, while sale fee and tax inputs remain their original line
  totals. Traces to milestone criterion ID: `M-06-C02`. Proven by:
  `go test ./internal/modules/profitability/application -run
  TestImportMarginInputsExtendsOnlyUnitCostByQuantity -count=1 -v`.
- Nil/missing CUSSEMICM remains nil with explicit missing quality/reason and
  is not converted to zero. Traces to milestone criterion ID: `M-06-C02`.
  Proven by: `go test ./internal/modules/profitability/application -run
  TestImportMarginInputsPreservesUnknownCost -count=1 -v`.

## Handoff

- Current status: `spec_ready`
- Next owner: Feature Implementer
- Next action: Write `plan.md`, implement the scoped composition correction,
  and record executable validation evidence.
- Blockers or open decisions: None; the signed 2026-07-12 owner decision
  fixes the amount scopes.
