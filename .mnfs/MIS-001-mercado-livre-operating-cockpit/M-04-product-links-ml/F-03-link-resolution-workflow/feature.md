# F-03-link-resolution-workflow

```yaml
id: F-03
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-04
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
```

## Brief

Add API/SDK/UI workflow for operators to approve, reject, and resolve product link candidates with audit evidence.

## Inputs

- Link candidates.
- Product candidate evidence.
- Existing frontend package conventions.

## Expected Output

- Operators see listing evidence, internal product candidates, confidence reason, and conflict/unresolved states.
- Link audit records who/what changed state and why.

## Constraints

- UI uses SDK runtime only.
- Link approval does not apply stock writes.

## Interaction Model

- List candidates -> filter by state -> inspect listing/product evidence -> approve/reject/resolve -> refetch link state.
- Conflict and unresolved states remain visible until resolved or rejected.

## Validation Expectations

- API tests cover approve/reject/conflict.
- UI tests cover loading, error, empty, conflict, and approval states.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: Create `spec.md`.
- Required files/evidence: F-03/validation.md.
- Blockers or open decisions: None.
