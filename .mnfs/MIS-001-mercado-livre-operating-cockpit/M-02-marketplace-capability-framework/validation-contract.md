# Milestone Validation Contract

```yaml
id: M-02
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

M-02

## QA Level

QA-0 planning contract.

## Required Outcome

Marketplace capability contracts exist and Mercado Livre implements the first adapter spine without leaking business policy into connector code.

## Criteria

## Criterion: Capability Ports Exist
ID: M-02-C01
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: Go tests and source inspection
- Expected: Listing, stock, and order capabilities are represented as small ports with provider ids as strings and no Mercado Livre endpoint payloads in business modules.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-02-marketplace-capability-framework/validation-result.md`
Blocking failure: Business service imports Mercado Livre adapter package.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Unsupported Capability Is Explicit
ID: M-02-C02
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: capability service tests
- Expected: Unsupported provider operation returns `CONNECTORS_PROVIDER_UNSUPPORTED_SHAPE` or capability `unsupported`, not nil success.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-02-marketplace-capability-framework/validation-result.md`
Blocking failure: Unsupported capability silently succeeds or panics.
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- Adapter tests use documented Mercado Livre item/order/variation shapes.
- Boundary checks prove no provider-specific business imports.
- Adapter/documentation tests may prove contract and mapping readiness, but they do not prove live Mercado Livre runtime integration.
- M-02 may be marked fully passed only when either:
  - its scope is explicitly limited to capability-contract and adapter-spine readiness in the milestone brief and validation result, or
  - real Mercado Livre environment validation evidence is present for the claimed runtime behavior.

## Blocking Failures

- Provider payload structs cross into domain/application modules.
- Capability status cannot represent unsupported/degraded.

## Retry Policy

- correction_attempts:
- max_correction_attempts: 2
- last_validation_result:

## Handoff

- Current status: planned.
- Next owner: QA Validator after execution.
- Next action: Validate capability seams.
- Required files/evidence: F-*/validation.md.
- Blockers or open decisions: None.
