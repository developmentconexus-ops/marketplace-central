# F-06 Quantity Cost Semantics — Validation

```yaml
id: F-06
type: feature-validation
status: quick_validation_passed
owner: Feature Implementer
parent: M-06
created: 2026-07-12
updated: 2026-07-12
validation_level: QA-0
lifecycle_scope: feature
```

## Summary

Quick validation passed. The import path now records Oracle CUSSEMICM as
`per_unit` at the read boundary and persists its cost input as an
item-line total. The focused proof covers quantities 1, 2, and 7: costs are
10, 20, and 70 for unit cost 10, while sale fee remains 3 and ICMS/IPI/PIS/
COFINS remain 1/2/3/4 for every quantity. Unknown cost remains `nil` with
`missing` quality and `missing_cost`; no monetary zero is introduced.

## Quick Validation Result

- Result: Pass
- Result owner: Feature Implementer
- Decision date: 2026-07-12
- Final feature state for handoff: `quick_validation_passed`

## Quick Validation State

- fixup_attempts: 0
- max_fixup_attempts: 1
- last_feature_validation_result: Pass

## Spec Adherence

- Spec satisfied: Yes
- Deviations: None.
- Reason: CUSSEMICM is labelled `per_unit` in `CostAsOf`, extended only once
  at profitability item import, and receives durable source-reference
  evidence. Sale fee and each tax component retain their line-total amount and
  source-reference scope evidence.

## Changes Made

- File: `apps/server_core/internal/modules/internal_read/domain/internal_cost.go`
  - Change: Added the explicit `CostAmountScopePerUnit` metadata to `CostAsOf`.
- File: `apps/server_core/internal/modules/internal_read/adapters/oracle/reader.go`
  - Change: Marks every Oracle CUSSEMICM result, including missing results, as
    per-unit.
- File: `apps/server_core/internal/modules/profitability/application/service.go`
  - Change: Extends only non-nil CUSSEMICM by item quantity; preserves nil;
    writes source references for cost unit-to-line, sale-fee line-total, and
    tax line-total semantics.
- File: `apps/server_core/internal/modules/profitability/application/service_test.go`
  - Change: Added quantity 1/2/7 and nil-cost regression coverage.

## Commands Run

- Command: `gofmt -w apps/server_core/internal/modules/internal_read/domain/internal_cost.go apps/server_core/internal/modules/internal_read/adapters/oracle/reader.go apps/server_core/internal/modules/profitability/application/service.go apps/server_core/internal/modules/profitability/application/service_test.go`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: scoped Go files are formatted.
  - Actual: command exited 0 with no diagnostic output.
  - Artifact: formatted files listed in Changes Made.
  - Blocking condition: None.
- Command: `cd apps/server_core; $env:GOCACHE="$PWD\\.gocache"; go test ./internal/modules/profitability/application ./internal/modules/internal_read/... -count=1`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: profitability application and internal-read packages pass.
  - Actual: exit 0; profitability application plus internal-read fake, Oracle,
    application, domain, and ports packages all passed.
  - Artifact: command output in this Feature execution task.
  - Blocking condition: None.
- Command: `cd apps/server_core; $env:GOCACHE="$PWD\\.gocache"; go test ./internal/modules/profitability/application -run "TestImportMarginInputsExtendsOnlyUnitCostByQuantity|TestImportMarginInputsPreservesUnknownCost" -count=1 -v`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: quantities 1/2/7 extend only cost and nil cost remains nil with
    missing quality.
  - Actual: both named tests passed; package exit 0.
  - Artifact: command output in this Feature execution task, including both
    named `PASS` lines.
  - Blocking condition: None.
- Command: `git diff --check -- .mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-06-quantity-cost-semantics apps/server_core/internal/modules/profitability/application/service.go apps/server_core/internal/modules/profitability/application/service_test.go apps/server_core/internal/modules/internal_read/domain/internal_cost.go apps/server_core/internal/modules/internal_read/adapters/oracle/reader.go`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: no whitespace errors in the F-06 change set.
  - Actual: exit 0; no diff-check errors (only repository line-ending warnings).
  - Artifact: command output in this Feature execution task.
  - Blocking condition: None.

## Manual QA

- QA level: QA-0
- Flow or step: Read focused verbose test observables for quantity and unknown
  cost semantics.
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Expected: line cost changes with quantity; sale fee/tax do not; unknown cost
  is not zero.
- Actual: both focused tests passed; their assertions cover the stated values
  and source-reference scope evidence.
- Blocking condition: None.

## Evidence

- Artifact: `spec.md`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Blocking condition: None.
- Artifact: `plan.md`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Blocking condition: None.
- Artifact: `validation.md`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Blocking condition: None.

## Risks

- This feature proves deterministic application behavior only; Milestone QA
  remains responsible for any later integrated or live evidence.
- Existing unrelated worktree changes were preserved and are not part of this
  feature commit.

## Milestone Acceptance Review

- Decision: accepted.
- Reviewer: Milestone Orchestrator.
- Reviewed commit: `2284c1d3bfcfa359a66777baad6c339083973538`.
- Scope: accepted. The eight committed paths are the F-06 execution artifacts,
  CUSSEMICM amount-scope metadata, Oracle read metadata, and the profitability
  composition/test seam. No API, SDK, provider, Oracle-write, Candidate A,
  manual-adjustment, or principal-boundary path changed.
- Evidence: accepted. The Feature records `ran` evidence for the registered
  profitability/internal-read suite and focused quantity 1/2/7 plus
  nil-cost tests. The Orchestrator independently reran both commands on the
  reviewed commit with exit 0 on 2026-07-12.
- Constraint review: accepted. The focused test proves line cost `10/20/70`
  from unit cost `10`, unchanged sale fee `3`, unchanged tax components
  `1/2/3/4`, and nil cost remaining nil with `missing_cost`; no zero default
  is introduced.
- QA routing: no formal feature QA invoked. This is a bounded deterministic
  cost-composition correction without an auth, PII, secret, multi-role, or
  live-runtime change. Formal milestone review and proportional QA remain
  required after all correction seams integrate.
- Next owner: Milestone Orchestrator.

## Handoff

- Current status: `quick_validation_passed`
- Next owner: Milestone Orchestrator
- Next action: Review the signed semantics, scoped diff, and ran evidence;
  integrate this single feature commit before fixed-SHA milestone review/QA.
- Required files/evidence: `feature.md`, `spec.md`, `plan.md`, `validation.md`,
  and the two Go test commands above.
- Blockers or open decisions: None.
