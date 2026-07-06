# F-03-data-quality-rules

```yaml
id: F-03
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-03
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-001 Mercado Livre Operating Cockpit.

## Milestone

M-03 MNOS/Sankhya read contract.

## Brief

Define reusable data-quality flags for internal inputs so missing/ambiguous/stale facts block or degrade downstream workflows visibly.

## Inputs

- IC-002 quality enum.
- Product Links, Inventory, Profitability planned states.

## Expected Output

- Shared value objects/enums for `missing_product`, `missing_stock`, `missing_cost`, `missing_tax`, `ambiguous_product`, `stale_source`.
- Tests prove missing values do not become zero.

## Constraints

- Do not add UI copy here beyond enum naming.
- Do not decide product exclusion policy here.

## Validation Expectations

- Unit tests cover every quality enum and at least one blocked stock case and one incomplete margin case.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: Create `spec.md`.
- Required files/evidence: F-03/validation.md.
- Blockers or open decisions: None.
