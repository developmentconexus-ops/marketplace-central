# F-04-orders-margin-ui

```yaml
id: F-04
type: feature-brief
status: quick_validation_passed
owner: Mission Strategist
parent: M-06
created: 2026-07-06
updated: 2026-07-09
validation_level: QA-0
lifecycle_scope: feature
```

## Brief

Build orders and margin UI so operators can inspect sale profitability, missing inputs, and manual adjustments.

## Inputs

- Orders API.
- Profitability API.
- SDK runtime.

## Expected Output

- Orders list with status, item, link quality, margin quality, and contribution.
- Order detail shows input breakdown and adjustment history.

## Constraints

- React does not calculate margin.
- PII is minimized.

## Interaction Model

- List orders -> filter missing margin inputs -> open order -> add freight/commission adjustment -> margin snapshot refreshes.

## Validation Expectations

- UI tests cover loading, error, empty, missing input, complete margin, negative margin, and manual adjustment states.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: quick_validation_passed.
- Next owner: Milestone Orchestrator.
- Next action: Validate M-06 milestone end-to-end over the now-verified backend + UI slices.
- Required files/evidence: F-04/validation.md.
- Blockers or open decisions: None.
