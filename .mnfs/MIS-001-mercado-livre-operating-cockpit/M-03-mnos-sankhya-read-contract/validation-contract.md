# Milestone Validation Contract

```yaml
id: M-03
type: milestone-validation-contract
status: planned
owner: Mission Strategist
parent: MIS-001
created: 2026-07-06
updated: 2026-07-07
validation_level: QA-0
lifecycle_scope: milestone
```

## Milestone ID

M-03

## QA Level

QA-0 planning contract.

## Required Outcome

Marketplace Central owns a real Oracle-backed internal-read boundary inside `apps/server_core` for product, stock, price, cost, tax, and sales facts, with explicit quality states, no ERP write path, and no false validation claims.

## Criteria

## Criterion: MPC Owns The Oracle Read Contract
ID: M-03-C01
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: focused Go tests for `internal_read/domain`, `ports`, and contract-level policy objects
- Expected: downstream modules depend only on MPC-owned domain/port types; Oracle table/query details do not leak outside the adapter boundary.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/validation-result.md`
Blocking failure: business modules import Oracle/driver/query details directly.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Sellable Stock And Margin Inputs Are Explicit And Auditable
ID: M-03-C02
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: contract tests plus Oracle adapter mapping tests
- Expected: sellable stock, cost basis, tax inputs, and source freshness semantics are explicit in MPC-owned contracts and trace to named Oracle evidence; missing facts remain quality states rather than zero/defaults.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/validation-result.md`
Blocking failure: stock/cost/tax semantics are hidden in ad hoc queries, undocumented constants, or silent defaults.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Real Oracle Adapter Works For Claimed Runtime Paths
ID: M-03-C03
Level: Milestone
Type: Integration
Required: Yes
Status: Pending
Evidence:
- Command: real-environment validation command or targeted harness against operator-approved Oracle access
- Expected: at least one end-to-end read for each claimed runtime surface succeeds against the real Oracle source, with captured evidence for product, stock, cost/tax, and source timestamps.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/validation-result.md`
Blocking failure: milestone claims live behavior without direct Oracle-backed evidence.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Security And Boundary Discipline Hold
ID: M-03-C04
Level: Milestone
Type: Security
Required: Yes
Status: Pending
Evidence:
- Command: secret-safety checks, write-path grep, and adapter-boundary inspection
- Expected: Oracle credentials/secrets never appear in logs/artifacts; no ERP write SQL exists; no ERP mirror tables are introduced outside MPC-owned state.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/validation-result.md`
Blocking failure: secret leakage, write-path introduction, or ERP mirror creep.
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- Legacy MNOS artifacts may be cited as reference evidence, but they are not sufficient alone.
- MPC-owned specs must name the Oracle evidence sources, contract semantics, and adapter ownership clearly.
- Fake/test seams may prove consumer behavior and deterministic edge cases, but they do not prove live Oracle behavior.
- M-03 may be marked passed only when real Oracle validation evidence exists for every runtime behavior claimed in the milestone result.

## Blocking Failures

- Any ERP write path inside `internal_read`.
- Any silent zero/default for missing stock, cost, tax, or source freshness.
- Any Oracle query logic leaking into downstream business modules.
- Any validation result that implies real integration proof without fresh Oracle-backed evidence.

## Retry Policy

- correction_attempts:
- max_correction_attempts: 2
- last_validation_result:

## Handoff

- Current status: planned.
- Next owner: QA Validator after Oracle-first execution.
- Next action: validate the rewritten contract, adapter implementation, and real-environment evidence.
- Required files/evidence: rewritten F-*/validation.md plus direct Oracle evidence artifacts.
- Blockers or open decisions: real Oracle validation harness and exact evidence capture flow must be defined during feature execution.
