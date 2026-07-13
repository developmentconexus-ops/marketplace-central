# F-08 Assisted Sankhya Linkage Discovery Validation

```yaml
id: F-08
type: feature-validation
status: quick_validation_passed
owner: Feature Implementer
parent: F-08
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-1
lifecycle_scope: feature
```

## Feature ID

F-08-assisted-sankhya-linkage-discovery

## Summary

Quick validation passed for the bounded documentation Feature. Repository facts
were reconciled; the fail-closed admin and implementation contracts cover all
acceptance criteria. No live Oracle discovery ran because the registered work
contract ID has no executable repository definition. Runtime field/TOP facts
remain unknown and production activation remains gated.

## Quick Validation Result

- Result: Pass
- Result owner: Feature Implementer
- Decision date: 2026-07-13
- Final feature state for handoff: `quick_validation_passed`

## Quick Validation State

- fixup_attempts: 0
- max_fixup_attempts: 1
- last_feature_validation_result: Pass

## Spec Adherence

- Spec satisfied: Yes
- Deviations: No bounded live SELECT was executed.
- Reason: `live-oracle-linkage-discovery` has no executable definition in
  repository truth; substituting ad hoc SQL or the generic Docker smoke test
  would violate the governed-command and evidence contract. The Feature brief
  explicitly permits unknown runtime facts with a fail-closed deployable spec.

## Changes Made

- `feature.md`: added compiler-required brief and observable outputs.
- `spec.md`, `plan.md`, `context.json`: acceptance, execution, and validated
  context contracts.
- `discovery.md`: classified repository facts, runtime facts, inferences, exact
  predicates, and unknowns.
- `sankhya-admin-spec.md`: deployable configured header-field contract and
  explicit no-required-line-field decision.
- `implementation-contract.md`: next-worker boundaries, ledger, idempotency,
  audit, unknowns, and exact migration/API/SDK seams.

## Commands Run

- Command: import `scripts/harness/Context.psm1`; run
  `New-HarnessContextPack` for the F-08 path/base/allowed path, then
  `Test-HarnessContextPack -RequireCurrentBase`.
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: compile and validate minimal context at accepted SHA.
  - Actual: both returned `Passed=True`, context ID
    `f-08-assisted-sankhya-linkage-discovery`, risk `L3`.
  - Artifact: `context.json`
  - Blocking condition: none.
- Command ID: `live-oracle-linkage-discovery`; definition lookup with
  `rg -l --fixed-strings live-oracle-linkage-discovery scripts contracts apps`.
  - Status: Not run
  - Evidence type: could-not-run
  - Owner: Feature Implementer
  - Expected: a governed bounded SELECT command with values-safe output.
  - Actual: executable definition count `0`; no Docker or Oracle command ran.
  - Artifact: `discovery.md`; this validation limitation.
  - Blocking condition: runtime activation remains blocked; documentation
    Feature completion is not blocked because unknowns are explicitly allowed.
- Command ID: `go-tax-provenance-focused`; from `apps/server_core`, set
  repository-root `GOCACHE=.gocache` and run
  `go test ./internal/modules/internal_read/adapters/oracle ./internal/modules/internal_read/adapters/fake`.
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: exact-identity and missing-tax contract tests pass.
  - Actual: both packages `ok` (cached), exit code 0.
  - Artifact: this validation record; F-07 deterministic contract.
  - Blocking condition: none.
- Command: bounded sensitive-literal scan over F-08 Markdown/JSON.
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: no credential/secret/PII value assignments.
  - Actual: `sensitive_literal_hits=0`.
  - Artifact: this validation record.
  - Blocking condition: none.
- Command ID: `git-diff-check`; `git diff --cached --check` after staging the
  owned Feature directory, plus scoped staged-path inspection.
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: no whitespace errors and only the allowed path staged.
  - Actual: exit code 0; staged paths were limited to the F-08 directory.
  - Artifact: this validation record.
  - Blocking condition: none.

## Manual QA

- QA level: QA-1
- Flow or step: trace candidate → explicit confirmation → append-only exact
  313 line mapping → exact `TGFVAR` descendants → profitability identity.
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Expected: candidates never prove linkage; partial/one-to-many lineage stays
  explicit; unknown/conflicting lineage never becomes zero or resolved.
- Actual: all three artifacts preserve the boundary; no line custom field is
  required for assisted-only proof.
- Blocking condition: none for Feature handoff; activation gates remain.

## Evidence

- Artifact: `discovery.md`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Blocking condition: none; runtime facts explicitly unknown.
- Artifact: `sankhya-admin-spec.md`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Blocking condition: activation requires admin-supplied field and uniqueness evidence.
- Artifact: `implementation-contract.md`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Blocking condition: none for specification; next worker requires approved scope.
- Artifact: `context.json`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Blocking condition: none; validated against accepted current HEAD.

## Risks

- Exact deployed field metadata, uniqueness, effective TOP behavior, and live
  `TGFVAR` cardinality remain unknown until a registered bounded live command
  runs. Confirmation must remain disabled.
- Cancellation/devolution policy remains an owner/admin input; mappings and
  audit events must never be deleted or reused meanwhile.
- This record is Feature quick-validation evidence, not milestone acceptance
  or QA pass.

## Handoff

- Current status: `quick_validation_passed`
- Next owner: Milestone Orchestrator
- Next action: review/integrate this documentation commit; separately authorize
  an executable bounded linkage discovery command before runtime activation,
  then scope the implementation Feature against `implementation-contract.md`.
- Required files/evidence: spec, plan, context, discovery, admin spec,
  implementation contract, validation, commit
- Blockers or open decisions: no blocker to accepting this Feature; runtime
  activation is blocked on field/uniqueness/TOP/lineage evidence and sanctioned
  admin ownership.
