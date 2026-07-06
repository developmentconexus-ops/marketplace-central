# F-02-price-recommendation-candidates

```yaml
id: F-02
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

Generate manual-review price recommendation candidates based on cost, current ML price, fees, stock, and margin bands.

## Inputs

- Linked listings.
- Profit snapshots.
- Current price data.
- Margin guardrail policy.

## Expected Output

- Recommendation rows with current price, target price/range, reason, evidence, and quality.

## Constraints

- No Mercado Livre price write.
- Price review risk must be visible when future write support is planned.

## Validation Expectations

- Tests cover increase, decrease, no-change, missing cost, and low-quality blocked recommendation.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: Create `spec.md`.
- Required files/evidence: F-02/validation.md.
- Blockers or open decisions: exact margin bands.
