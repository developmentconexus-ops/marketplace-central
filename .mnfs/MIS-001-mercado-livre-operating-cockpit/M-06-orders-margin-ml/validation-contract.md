# Milestone Validation Contract

```yaml
id: M-06
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

M-06

## QA Level

QA-0 planning contract.

## Required Outcome

Orders + Margin provides idempotent order state and honest profitability quality.

## Criteria

## Criterion: Idempotent Order Ingestion
ID: M-06-C01
Level: Milestone
Type: Reliability
Required: Yes
Status: Pending
Evidence:
- Command: orders service/repository tests
- Expected: Reprocessing the same provider order id updates the same order and items without duplicates.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/validation-result.md`
Blocking failure: Duplicate provider order rows or duplicate item rows are created.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Margin Quality Honesty
ID: M-06-C02
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: profitability service tests
- Expected: Missing link, cost, freight, sale fee, or tax produces a specific quality flag; unknown values are not converted to 0.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/validation-result.md`
Blocking failure: Margin calculation treats missing value as zero without quality flag.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Manual Adjustment Audit
ID: M-06-C03
Level: Milestone
Type: Observability
Required: Yes
Status: Pending
Evidence:
- Command: manual adjustment API/repository tests
- Expected: Freight/commission/manual cost adjustment stores amount, type, order/item scope, operator, timestamp, and note.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/validation-result.md`
Blocking failure: Manual margin input can be changed without audit note.
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- Mercado Livre order fields used by adapter must cite official docs in feature specs.
- UI tests cover missing quality flags and completed margin.

## Blocking Failures

- Order ingestion blocked by missing product link.
- Unknown cost/freight/fee/tax treated as zero.
- Buyer PII exposed beyond necessary operational fields.

## Retry Policy

- correction_attempts: 1
- max_correction_attempts: 2
- last_validation_result: Fail (round 2, reviewed SHA 5548ae406cb26d0703c111236d703281bb227d3e)

Correction authority is recorded append-only in
`corrections/correction-task.md`. The owner-designated baseline is round-2
Fail; F-06 is attempt 1 of 2. Earlier unnumbered fixes remain historical
context and do not consume attempts.

## Handoff

- Current status: planned.
- Next owner: QA Validator after execution.
- Next action: Validate Orders + Margin.
- Required files/evidence: F-*/validation.md.
- Blockers or open decisions: None for planning.
