# Milestone Validation Contract

```yaml
id: M-03
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

M-03

## QA Level

QA-0 planning contract.

## Required Outcome

MPC can compute internal sellable stock and margin input quality using MNOS/Sankhya semantics without Sankhya writes.

## Criteria

## Criterion: Sellable Stock Contract
ID: M-03-C01
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: read contract unit tests
- Expected: `sellable_stock = SUM(ESTOQUE - RESERVADO)` where `CODEMP IN (1,2)` and `CODLOCAL=10101`; `CODLOCAL=10108` contributes 0 under default policy.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/validation-result.md`
Blocking failure: Reserved or showroom stock is announced as sellable.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Missing Inputs Are Quality Flags
ID: M-03-C02
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: read contract unit tests
- Expected: Missing product yields unresolved candidate; missing cost yields `missing_cost`; missing tax yields `missing_tax`; no missing numeric value is returned as 0.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/validation-result.md`
Blocking failure: Missing cost/tax/stock is silently converted to zero.
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- Tests cite MNOS source semantics in `spec.md`.
- Secret-safety tests cover Oracle/Sankhya connection error messages.
- Fake adapter tests may prove contract preservation and quality-flag behavior, but they do not prove live Sankhya/Oracle reads.
- M-03 may be marked fully passed only when either:
  - its scope is explicitly limited to contract/seam readiness in the milestone brief and validation result, or
  - real Sankhya/Oracle validation evidence is present for the claimed runtime behavior.

## Blocking Failures

- Any Sankhya write path.
- Any ERP table mirror beyond MPC-owned snapshots.
- Any default use of `CUSVARIAVEL` instead of `CUSSEMICM` for initial margin.

## Retry Policy

- correction_attempts:
- max_correction_attempts: 2
- last_validation_result:

## Handoff

- Current status: planned.
- Next owner: QA Validator after execution.
- Next action: Validate read contract semantics.
- Required files/evidence: F-*/validation.md.
- Blockers or open decisions: None.
