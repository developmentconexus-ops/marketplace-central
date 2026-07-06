# F-04-kit-bundle-opportunity-engine

```yaml
id: F-04
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

Suggest kit/bundle opportunities from sales history, complementary product groups, available stock, and margin quality.

## Inputs

- Sales history.
- Product groups.
- Stock availability.
- Margin quality.
- Operator-defined incompatible groups or eligible categories.

## Expected Output

- Bundle candidates with component products, stock feasibility, estimated contribution, evidence, and confidence.

## Constraints

- No automatic listing creation.
- Suggestions require operator review.
- Missing margin or stock feasibility blocks high-confidence candidate.

## Validation Expectations

- Tests cover valid bundle, insufficient component stock, missing margin, incompatible group, and low-confidence suggestion.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: Create `spec.md`.
- Required files/evidence: F-04/validation.md.
- Blockers or open decisions: operator-defined compatibility rules.
