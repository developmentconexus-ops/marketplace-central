# F-02-margin-input-model

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-06
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
```

## Brief

Define margin input entities for revenue, sale fee, cost, tax, freight, commission, and manual adjustments.

## Inputs

- Order snapshots.
- IC-002 cost/tax contract.
- Operator decision: cost basis `CUSSEMICM`.

## Expected Output

- Margin input model with explicit source, timestamp, amount, currency, and quality state.
- Manual adjustment scope supports order-level and item-level entries.

## Constraints

- Do not use `CUSVARIAVEL` as initial cost basis.
- Do not treat missing freight/tax/fee as zero without quality flag.

## Validation Expectations

- Tests cover complete input, missing cost, missing freight, manual freight, manual commission, and adjustment audit.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: Create `spec.md`.
- Required files/evidence: F-02/validation.md.
- Blockers or open decisions: exact default manual adjustment categories can be refined in spec.
