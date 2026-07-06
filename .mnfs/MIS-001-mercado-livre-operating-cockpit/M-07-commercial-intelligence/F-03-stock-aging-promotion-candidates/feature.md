# F-03-stock-aging-promotion-candidates

```yaml
id: F-03
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-07
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
```

## Brief

Identify stock aging and promotion candidates using stock quantity, last sale/exit, sales velocity, margin quality, and product eligibility.

## Inputs

- Internal stock.
- Sales history.
- Product metadata.
- Margin quality.

## Expected Output

- Candidate list with reason codes such as `aged_stock`, `high_stock_low_velocity`, `promotion_candidate`, `blocked_low_margin`.

## Constraints

- Do not recommend products ineligible for Mercado Livre.
- Low margin blocks aggressive promotion recommendation unless operator overrides in a later feature.

## Validation Expectations

- Tests cover aged stock, high stock with velocity, blocked low margin, and missing sales history.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: Create `spec.md`.
- Required files/evidence: F-03/validation.md.
- Blockers or open decisions: aging thresholds.
