# F-02-read-adapter-seam

```yaml
id: F-02
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

Create an application-facing read seam for internal product/stock/cost facts, with fake/test adapter first and a path for real Oracle read-only implementation.

## Inputs

- IC-002.
- Existing Go module patterns.

## Expected Output

- Ports and fake adapter support Stock Seguro and Product Links tests without live Oracle.
- Real adapter design preserves read-only credentials and SELECT-only behavior.

## Constraints

- No production Oracle write path.
- No secrets in logs/errors.

## Inputs/Outputs

- Input: `codprod`, EAN/SKU/title, sale date, stock policy scope.
- Output: product candidates, stock, cost, quality flags.

## Negative Scenarios

- Oracle unavailable maps to `SANKHYA_READ_UNAVAILABLE`.
- Ambiguous product maps to `SANKHYA_PRODUCT_AMBIGUOUS`.

## Validation Expectations

- Unit tests prove fake adapter returns exact stock formula and quality flags.
- Secret-safety test proves missing env errors list names only, not values.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: Create `spec.md`.
- Required files/evidence: F-02/validation.md.
- Blockers or open decisions: Live Oracle validation may require operator environment setup.
