# F-01-stock-policy-model

```yaml
id: F-01
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

Define `StockPolicy` for sellable scope, default buffer, source freshness, and product ineligibility rules.

## Inputs

- IC-002 stock contract.
- Operator decisions: `CODEMP IN (1,2)`, `CODLOCAL=10101`, exclude showroom `10108`, buffer default 1.

## Expected Output

- Domain value objects for stock scope, buffer, freshness, and eligibility.
- Policy can later add group/weight/size/margin blocks.

## Constraints

- No auto-write policy in this mission.
- Product exclusion rules must be explicit and visible.

## Validation Expectations

- Unit tests prove default formula and excluded showroom stock.
- Policy test proves buffer cannot produce negative recommended stock.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: Create `spec.md`.
- Required files/evidence: F-01/validation.md.
- Blockers or open decisions: exact product group/weight/size exclusion examples.
