# Milestone Validation Contract

```yaml
id: M-04
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

M-04

## QA Level

QA-0 planning contract.

## Required Outcome

Product Links can resolve existing Mercado Livre listing/variation identities to internal products with audit and blocked conflict states.

## Criteria

## Criterion: Exact Match Priority
ID: M-04-C01
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: product link service tests
- Expected: Exact EAN or `seller_sku` candidate outranks title heuristic; multiple exact matches create conflict, not auto-resolution.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-04-product-links-ml/validation-result.md`
Blocking failure: Ambiguous match becomes resolved without operator approval.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Link State Blocks Writes
ID: M-04-C02
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: product-links workflow tests, workflow API responses, and persisted state inspection
- Expected: `unresolved`, `conflict`, and `rejected` link states remain explicit and distinguishable from `resolved`, so downstream stock-action milestones can block unsafe proposals/apply paths instead of guessing.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-04-product-links-ml/validation-result.md`
Blocking failure: Product-link workflow collapses unresolved/conflict/rejected into an implicitly usable resolved-like state.
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- API/OpenAPI/SDK changes are validated together.
- UI tests prove conflict and empty states.

## Blocking Failures

- Provider listing id stored as integer instead of string.
- Link changes lack audit.

## Retry Policy

- correction_attempts:
- max_correction_attempts: 2
- last_validation_result:

## Handoff

- Current status: planned.
- Next owner: QA Validator after execution.
- Next action: Validate link state behavior.
- Required files/evidence: F-*/validation.md.
- Blockers or open decisions: None.
