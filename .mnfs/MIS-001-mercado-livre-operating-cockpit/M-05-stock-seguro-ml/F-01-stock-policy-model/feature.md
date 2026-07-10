# F-01-stock-policy-model

```yaml
id: F-01
type: feature-brief
status: quick_validation_passed
owner: Feature Implementer
parent: M-05
created: 2026-07-06
updated: 2026-07-09
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

- Current status: quick_validation_passed.
- Next owner: Milestone Orchestrator.
- Next action: Accept F-01 into M-05 sequencing and start F-02 risk engine on top of the inventory-owned policy model.
- Required files/evidence: `spec.md`, `plan.md`, `implementation-report.md`, and `validation.md`.
- Blockers or open decisions: exact product group/weight/size/margin exclusion examples remain deferred to later owner refinement.
