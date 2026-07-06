# F-02-stock-risk-engine

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-05
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
```

## Brief

Compare internal sellable stock and Mercado Livre announced stock for resolved links and classify risk.

## Inputs

- Product links.
- Internal stock reads.
- Provider stock snapshots.
- Stock policy.

## Expected Output

- Risk states: `healthy`, `oversell`, `undersell`, `stale`, `unresolved`, `conflict`, `ineligible`, `unsupported`.
- Recommended quantity and policy evidence for actionable rows.

## Constraints

- Do not call provider synchronously from the UI read endpoint.
- Do not produce actionable recommendation for blocked states.

## Inputs/Outputs

- Input: link id, listing snapshot, internal stock, provider stock, policy.
- Output: risk row with state, quantities, recommended quantity, source timestamps, blocking reason.

## Negative Scenarios

- Missing internal stock returns `stale` or blocked data-quality state.
- Conflict link returns `conflict` and no recommendation.

## Validation Expectations

- Unit tests cover all risk states and source timestamp handling.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: Create `spec.md`.
- Required files/evidence: F-02/validation.md.
- Blockers or open decisions: source freshness threshold default.
