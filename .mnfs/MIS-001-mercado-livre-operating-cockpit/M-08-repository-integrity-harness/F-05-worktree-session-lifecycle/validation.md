# F-05 — Goal Orchestration Control Plane Validation

```yaml
id: F-05
type: feature-validation
status: quick_validation_passed
owner: Feature Implementer
parent: F-05
created: 2026-07-11
updated: 2026-07-11
validation_level: QA-3
lifecycle_scope: feature
```

## Feature ID

F-05-worktree-session-lifecycle

## Summary

The deterministic control-plane contract passed in the normal checkout. All
executed evidence is `fake`/contract evidence; F-09 and independent QA retain
the operator-observed milestone dogfood and any real-target decision.

## Quick Validation Result

- Result: Pass
- Result owner: Feature Implementer
- Decision date: 2026-07-11
- Final feature state for handoff: `quick_validation_passed`

## Quick Validation State

- fixup_attempts: 1 (impact child Git safe-directory environment only)
- max_fixup_attempts: 1
- last_feature_validation_result: Pass

## Spec Adherence

- Spec satisfied: Yes
- Deviations: One consolidated L2 review correction added F05-AC07 and native
  worktree proof before feature acceptance review.
- Reason: Routes, selector metadata, ignored state, leases, checkpoints,
  resume, risk policies, registered commands, and progressive skill all remain
  inside the declared control-plane paths.

## Changes Made

- `contracts/governance/knowledge-routes.json` and schema: canonical compact selectors.
- `scripts/harness/Context.psm1`: route resolution, selector/hash metadata, and L2/L3 overflow accounting.
- `scripts/harness/State.psm1`: ignored non-destructive leases, checkpoints, resume, and handoff validation.
- `scripts/harness/Impact.psm1`: registered deterministic tests and process-local Git safe-directory config.
- `AGENTS.md`, M-08 guide, and repo skill: compact bootstrap with progressive detail.
- Focused tests: deterministic route, state, risk, governance, and impact coverage.
- Consolidated L2 correction: a native detached-worktree fixture starts from the
  named SHA and proves the primary checkout is unchanged.

## Commands Run

- Command: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/harness-orchestration.tests.ps1`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: routes, conflicting leases, checkpoint resume, risk policies,
    skill bounds, handoff, and native detached-worktree coordination pass.
  - Actual: `Preparing worktree (detached HEAD 4265d850)` then
    `PASS harness orchestration tests`.
  - Artifact: terminal output, 2026-07-11; fixture-owned native worktree is
    removed by `git worktree remove` after assertion.
  - Blocking condition: None.
- Command: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/context-compiler.tests.ps1`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: current packs validate; stale/tampered/oversized packs fail closed.
  - Actual: `PASS context compiler tests`.
  - Artifact: terminal output, 2026-07-11.
  - Blocking condition: None.
- Command: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/governance-contracts.tests.ps1`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: route and existing governance schemas reject invalid documents.
  - Actual: `PASS governance contract tests`.
  - Artifact: terminal output, 2026-07-11.
  - Blocking condition: None.
- Command: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/harness.ps1 -Command governance -BaseSha 9ee2f0821599b31d115fe03d854a0490b7a60994`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: governance contracts and drift pass against the accepted base.
  - Actual: `status=passed`; only pre-existing baseline exceptions reported.
  - Artifact: `contracts/governance`.
  - Blocking condition: None.
- Command: `New-HarnessContextPack` plus `scripts/harness.ps1 -Command context-validate -ContextPath scripts/.runs/f05-context-pack.json -RequireCurrentBase`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: F-05 pack is hash-current and names every selector.
  - Actual: passed at 2,663 estimated tokens, L2 with per-source `overflow_reason`.
  - Artifact: `scripts/.runs/f05-context-pack.json` (ignored local state).
  - Blocking condition: None.
- Command: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/harness.ps1 -Command impact -ContextPath scripts/.runs/f05-context-pack.json`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: current checkout executes only registered fake commands.
  - Actual: `status=passed`, seam `harness-control`.
  - Artifact: `scripts/.runs/b12daa02720c433bbddcb361672564a0/outcome.json` (ignored local state).
  - Blocking condition: None.

## Manual QA

- QA level: QA-3
- Flow or step: operator-observed milestone task/subagent dogfood and fresh-session continuation.
- Status: Not run
- Evidence type: assumed
- Owner: F-09 / independent QA Validator
- Expected: supported native controls are observed and resume remains artifact-only.
- Actual: deliberately deferred by the F-05/F-09 sequence; no live capability claim is made here.
- Blocking condition: not a load-bearing F-05 quick-validation criterion.

## Evidence

- Artifact: `scripts/.runs/b12daa02720c433bbddcb361672564a0/outcome.json`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Blocking condition: None.
- Artifact: F-05 `spec.md` and `plan.md`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Blocking condition: None.

## Risks

- The 2,000-token target is exceeded for this L2 harness seam; every selected
  source records a reason. F-09 must collect efficiency metrics from dogfood.
- Native task/subagent controls are documented as operator-observed capability,
  not an API contract. Independent QA still owns milestone verdicts.
- The C06 fixture proves Git coordination only; it makes no isolation, clean
  machine, dependency, Docker, or external-runtime claim.

## Handoff

- Current status: `quick_validation_passed`
- Next owner: Milestone Orchestrator
- Next action: perform one proportional fixed-commit review and accept/reject F-05; then dispatch F-09 only after acceptance.
- Required files/evidence: feature brief, spec, plan, changed paths, validation evidence, current-checkout impact outcome.
- Blockers or open decisions: None.

## Milestone acceptance review

- Decision: Accepted.
- Reviewer: Milestone Orchestrator.
- Review basis: fixed-commit L2 review found the missing M-08-C06 proof. The
  sole consolidated correction added a native detached-worktree fixture; its
  privileged local execution passed and verified the primary checkout remains
  unchanged.
- Scope: all F-05 criteria map to `ran` deterministic evidence. Native control
  surfaces remain operator-observed capability, not a persisted API. No live
  target or milestone verdict is claimed.
- Next owner: Milestone Orchestrator awaiting Portfolio authorization for F-09.
