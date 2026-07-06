# Milestone Validation Contract

```yaml
id: M-01
type: milestone-validation-contract
status: planned
owner: Mission Strategist
parent: MIS-001
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: milestone
```

## Milestone ID

M-01

## QA Level

QA-0 planning contract.

## Required Outcome

VTEX has no active runtime, contract, SDK, or UI path in Marketplace Central target architecture.

## Criteria

## Criterion: VTEX Active Surfaces Removed
ID: M-01-C01
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: `rg -n "VTEX|vtex" apps packages contracts wiki docs`
- Expected: Matches are absent from active runtime/contract/SDK/UI paths or explicitly marked historical/legacy.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-01-vtex-removal-architecture-reset/validation-result.md`
Blocking failure: Active route, SDK method, adapter, frontend page, or target doc still exposes VTEX as current path.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Contract And SDK Aligned
ID: M-01-C02
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: frontend/sdk tests and OpenAPI inspection
- Expected: No `/connectors/vtex/*` routes or `publishToVTEX` SDK methods remain; tests pass.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-01-vtex-removal-architecture-reset/validation-result.md`
Blocking failure: SDK and OpenAPI disagree after removal.
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- F-01 inventory evidence must be linked before deletion claims.
- F-02 removal evidence includes tests and `rg`.
- F-03 doc evidence includes changed docs list.

## Blocking Failures

- Runtime still mounts VTEX routes.
- Frontend still navigates to VTEX publisher.
- SDK still exposes VTEX publish commands.

## Retry Policy

- correction_attempts:
- max_correction_attempts: 2
- last_validation_result:

## Handoff

- Current status: planned.
- Next owner: QA Validator after milestone execution.
- Next action: Validate removal evidence.
- Required files/evidence: F-*/validation.md.
- Blockers or open decisions: None.
