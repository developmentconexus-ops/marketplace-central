# F-03-profit-snapshot-calculation

```yaml
id: F-03
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

Calculate profit snapshots per order and item using revenue, Mercado Livre sale fee, `CUSSEMICM`, taxes, freight, and manual adjustments with quality flags.

## Inputs

- Margin input model.
- Product links.
- Sankhya cost/tax reads.

## Expected Output

- `ProfitSnapshot` with gross revenue, provider fee, cost, tax, freight, adjustments, contribution value, margin percent, and quality.

## Constraints

- Server-side calculation only.
- Unknown values remain null/quality flagged.

## Validation Expectations

- Unit tests cover complete margin, missing cost, missing link, manual adjustment, canceled order, and negative-margin order.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: quick_validation_passed.
- Next owner: Milestone Orchestrator.
- Next action: Proceed to F-04 orders and margin UI over the validated snapshots and quality flags.
- Required files/evidence: F-03/validation.md.
- Blockers or open decisions: None.
