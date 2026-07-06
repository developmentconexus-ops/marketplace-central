# Milestone Validation Contract

```yaml
id: M-07
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

M-07

## QA Level

QA-0 planning contract.

## Required Outcome

Commercial recommendations are evidence-backed, quality-aware, and manual-review only.

## Criteria

## Criterion: Recommendation Evidence
ID: M-07-C01
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: commercial intelligence service tests
- Expected: Every recommendation includes source metrics, quality state, reason code, and no automatic provider write.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-07-commercial-intelligence/validation-result.md`
Blocking failure: Recommendation lacks evidence or triggers provider write.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Low Quality Blocks Strong Claims
ID: M-07-C02
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: recommendation quality tests
- Expected: Missing link, missing cost, incomplete margin, or stale stock prevents high-confidence price/kit recommendation.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-07-commercial-intelligence/validation-result.md`
Blocking failure: Low-quality data produces high-confidence recommendation.
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- Recommendation tests include positive and blocked cases.
- UI shows evidence, not just suggestion text.

## Blocking Failures

- Automatic price/stock/listing write.
- Recommendation using missing cost as zero.
- Kit suggestion without stock/margin evidence.

## Retry Policy

- correction_attempts:
- max_correction_attempts: 2
- last_validation_result:

## Handoff

- Current status: planned.
- Next owner: QA Validator after execution.
- Next action: Validate recommendations after implementation.
- Required files/evidence: F-*/validation.md.
- Blockers or open decisions: Business threshold values.
