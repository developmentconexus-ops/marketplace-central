# F-02-oracle-adapter-implementation

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-03
created: 2026-07-06
updated: 2026-07-07
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-001 Mercado Livre Operating Cockpit.

## Milestone

M-03 Oracle internal-read contract.

## Brief

Implement the real Oracle adapter boundary for `internal_read`, including config loading, connection ownership, query organization, mapping, deterministic test seams, and explicit error/quality behavior.

## Inputs

- ADR-006.
- Rewritten F-01 contract.
- Existing Go module patterns.

## Expected Output

- A real Oracle adapter exists inside `internal_read/adapters/oracle`.
- Query/mapping responsibility is isolated inside the adapter boundary.
- A deterministic fake adapter remains available for unit and downstream tests.

## Constraints

- No ERP write path.
- No secret leakage in logs, errors, or artifacts.
- No ad hoc Oracle access from downstream modules.

## Inputs/Outputs

- Input: typed `Reader` requests from MPC modules.
- Output: product candidates, stock, price, cost, tax, and sales facts with quality states and source metadata.

## Negative Scenarios

- Oracle unavailable becomes an MPC-owned unavailable error.
- Unsupported query shape becomes an explicit contract error.
- Ambiguous product resolution becomes an explicit contract error or quality outcome.

## Validation Expectations

- Unit tests prove config secrecy, mapping behavior, and fake seam determinism.
- Real-environment validation proves the claimed Oracle runtime paths with actual source data.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: Rewrite `spec.md` around the real Oracle adapter.
- Required files/evidence: F-02/validation.md.
- Blockers or open decisions: Oracle driver/runtime choice and validation harness details must be finalized during execution.
