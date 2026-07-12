# F-06-quantity-cost-semantics

```yaml
id: F-06
type: feature-brief
status: briefed
owner: Milestone Orchestrator
parent: M-06
created: 2026-07-12
updated: 2026-07-12
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-001 Mercado Livre Operating Cockpit.

## Milestone

M-06 Orders + Margin ML.

## Brief

Correct the quantity-composition seam for profitability item snapshots. The
signed owner decision defines `GetCostAsOf` / `CUSSEMICM` as a per-unit cost;
an item-line cost is `CUSSEMICM × Mercado Livre item quantity`. Mercado Livre
`sale_fee` and each Oracle tax component are already line totals and must not
be quantity-multiplied.

## Inputs

- `M-06-orders-margin-ml/validation-contract.md` criterion M-06-C02.
- `M-06-orders-margin-ml/milestone-review.md` round-2 ★4 finding.
- `M-06-orders-margin-ml/_gate-evidence/round-2/gate-results.md` quantity
  evidence.
- `research/mnos-sankhya-read-interface-contract.md` `GetCostAsOf` contract.
- Signed M-06 owner decision dated 2026-07-12.

## Expected Output

- Explicit, durable scope metadata/evidence for unit CUSSEMICM versus line
  cost, sale fee, and tax component semantics.
- Correct item line-cost composition without changing unknown-value behavior.
- Tests for quantities 1, 2, and 7, plus nil/unknown cost preservation.

## Constraints

- Do not estimate tax or multiply line-sale fees or line-tax components.
- Do not change manual adjustments, trusted-principal authorization, API,
  OpenAPI, SDK runtime, provider reads/writes, Oracle writes, Candidate A, or
  unrelated M-06 seams.
- Preserve tenant scope and never convert unknown monetary values to zero.

## Validation Expectations

- Targeted profitability tests prove the signed line-composition rules and
  nil/unknown preservation.
- Record exact executed commands and observables as `ran` evidence in
  `validation.md`.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution,
not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: Create `spec.md`, then `plan.md`, then implement the bounded correction and record `validation.md`.
- Required files/evidence: `spec.md`, `plan.md`, `validation.md`, and targeted Go test output.
- Blockers or open decisions: None; the 2026-07-12 signed owner decision fixes all quantity semantics for this feature.
