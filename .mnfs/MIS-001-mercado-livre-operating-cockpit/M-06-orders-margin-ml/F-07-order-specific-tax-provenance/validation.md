# F-07 order-specific tax provenance validation

```yaml
id: F-07
type: feature-validation
status: quick_validation_passed
owner: Feature Implementer
parent: F-07
created: 2026-07-12
updated: 2026-07-12
validation_level: QA-1
lifecycle_scope: feature
```

## Feature ID

F-07-order-specific-tax-provenance

## Summary

Quick validation passed. Oracle tax reads now require an exact positive
`NUNOTA`/`SEQUENCIA` identity, retain product and incidence as consistency
predicates, return the identity in provenance, and never use product/date
aggregation as an order-item fact. Because the current order fact has no
owner-verified Oracle source linkage, profitability deliberately supplies no
identity and records all four tax inputs as missing.

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
- Deviations: No live Oracle execution was attempted; the dispatch requested
  proportional focused proof and the feature does not claim real linkage.
- Reason: Deterministic tests prove the contract and missing-state boundary.

## Changes Made

- Internal-read tax domain: added verified Oracle document/line identity and
  returned provenance.
- Oracle/fake adapters: exact identity selection and explicit missing tax when
  identity or rows are absent; partial components remain nullable.
- Profitability port/adapter/service: carries the identity contract and passes
  an empty identity until upstream provides an owner-verified mapping.
- Tests: exact SQL predicates/arguments, no-identity behavior, fake provenance
  and partial tax, and resolved-product/no-guessed-provenance composition.

## Commands Run

- Command: `$env:GOCACHE=(Join-Path (Resolve-Path ..\..).Path '.gocache'); go test ./internal/modules/internal_read/adapters/oracle ./internal/modules/internal_read/adapters/fake`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: exact source selection and missing/partial tax tests pass.
  - Actual: both packages `ok` (cached on final rerun), exit code 0.
  - Artifact: this validation record and package test output in task execution.
  - Blocking condition: none.
- Command: `$env:GOCACHE=(Join-Path (Resolve-Path ..\..).Path '.gocache'); go test ./internal/modules/profitability/...`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: profitability boundary and snapshot tests pass.
  - Actual: all profitability packages passed, exit code 0.
  - Artifact: this validation record and package test output in task execution.
  - Blocking condition: none.
- Command: `git diff --check`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: no whitespace errors.
  - Actual: exit code 0; only line-ending notices for working-copy policy.
  - Artifact: this validation record.
  - Blocking condition: none.

## Manual QA

- QA level: QA-1
- Flow or step: inspect Oracle query and profitability call boundary.
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Expected: exact `NUNOTA`, `SEQUENCIA`, `CODPROD`, `CODINC` predicates; no
  date predicate; no inferred identity from resolved product/date.
- Actual: implementation and focused tests show those exact boundaries.
- Blocking condition: none.

## Evidence

- Artifact: `context.json`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Blocking condition: none; compiled and validated against accepted HEAD.
- Artifact: `validation.md`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Blocking condition: none.

## Risks

- Real resolved-order margin remains operationally incomplete until an
  owner-approved upstream process supplies exact Oracle `NUNOTA` and
  `SEQUENCIA` for each marketplace order item.
- This feature does not establish or approve that cross-system mapping and does
  not claim live Oracle resolved-order evidence.

## Handoff

- Current status: `quick_validation_passed`
- Next owner: Milestone Orchestrator
- Next action: Review and integrate the scoped commit; retain the exact Oracle
  source-linkage gap in milestone QA/live-evidence reporting.
- Required files/evidence: feature, spec, plan, context, validation, commit.
- Blockers or open decisions: No implementation blocker. Real linkage remains
  a named upstream owner decision/input before resolved-order tax can complete.
