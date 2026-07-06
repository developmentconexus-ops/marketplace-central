# F-01-margin-guardrail-policy

```yaml
id: F-01
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

Define margin guardrail policies for marketplace recommendations, including minimum margin thresholds and blocked states.

## Inputs

- Profit snapshots.
- Stock policies.
- Operator-defined margin targets.

## Expected Output

- Policy model can classify recommendation eligibility by product/listing/group.

## Constraints

- No automatic price writes.
- Missing margin blocks high-confidence recommendation.

## Validation Expectations

- Tests cover above-threshold, below-threshold, missing cost, and missing fee cases.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: Create `spec.md`.
- Required files/evidence: F-01/validation.md.
- Blockers or open decisions: exact threshold values.
