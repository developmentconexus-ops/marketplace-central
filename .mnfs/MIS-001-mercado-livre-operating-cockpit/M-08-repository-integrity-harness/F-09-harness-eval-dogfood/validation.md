# F-09 Harness Eval and Dogfood — Validation

```yaml
id: F-09
type: feature-validation
status: quick_validation_passed
owner: Feature Implementer
parent: M-08
created: 2026-07-11
updated: 2026-07-11
validation_level: QA-3
lifecycle_scope: feature
```

## Feature ID

F-09-harness-eval-dogfood

## Summary

The deterministic eval corpus, independently added dogfood slice, current
context pack, canonical checkpoint/resume, and registered impact gate passed.
All executed evidence is `fake`/`contract`; it proves no real target and makes
no provider, Oracle, browser, or product claim.

## Quick Validation Result

- Result: Pass
- Result owner: Feature Implementer
- Decision date: 2026-07-11
- Final feature state for handoff: quick_validation_passed

Feature quick validation is not feature acceptance or a milestone verdict.
Milestone Orchestrator performs the fixed-commit review; independent QA alone
can write `validation-result.md` or pass M-08.

## Evidence Honesty

- All listed commands and artifacts have evidence type `ran`.
- Test and impact targets are `fake`; evidence class is `contract`.
- The native task/subagent dispatch, compact result, and fresh continuation are
  operator-observed capability evidence, not a persisted task API.

## Quick Validation State

- fixup_attempts: 1
- max_fixup_attempts: 1
- last_feature_validation_result: passed
- Consolidated correction: the grader now evaluates each case's declared
  observed condition before comparing its pinned verdict/reason. It also aligns
  the dogfood reason with `HEVAL_PATH_OUT_OF_SCOPE`; the full focused plan was
  rerun afterwards.

## Spec Adherence

- Spec satisfied: Yes
- Deviations: None.
- Reason: corpus is closed and deterministic; manifest data never becomes
  executable argv; dogfood worker changed only its three pinned paths.

## Changes Made

- `contracts/governance/harness-evals.json` — versioned 13-case negative
  corpus with pinned verdict/reason, context metrics, target, exit, and
  relative artifact path.
- `scripts/harness/Evals.psm1` — pure closed grader and result-manifest writer.
- `scripts/harness/Impact.psm1` — registered `harness-evals` structured argv.
- `scripts/tests/harness-eval.tests.ps1` — corpus/result-shape proof.
- `scripts/tests/fixtures/harness-eval/out-of-scope-write/attempt.json` — fresh
  bounded-worker dogfood fixture.
- `F-09/{spec,plan,validation}.md` — execution truth.

## Commands Run

- Command: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/harness-eval.tests.ps1`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: 13 pinned rejections, complete fake/contract result manifest.
  - Actual: `PASS harness eval tests` after the fresh worker added only the
    out-of-scope row, test assertion, and fixture.
  - Artifact: `scripts/.runs/harness-eval-tests/` (ignored result manifests).
  - Blocking condition: None.
- Command: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/harness-orchestration.tests.ps1`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: native detached worktree, lease/checkpoint/resume, and capability
    boundary pass without changing primary checkout.
  - Actual: `PASS harness orchestration tests`; run with the approved local
    worktree permission because the sandbox otherwise denies `.git/worktrees`.
  - Artifact: terminal evidence; temporary worktree removed by its fixture.
  - Blocking condition: None.
- Command: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/context-compiler.tests.ps1`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: current context remains route-selected and fail-closed.
  - Actual: `PASS context compiler tests`.
  - Artifact: `scripts/.runs/b4060a73d0354372a104921154fcb204/context-pack.json`.
  - Blocking condition: None.
- Command: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/governance-contracts.tests.ps1`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: governance contracts remain valid.
  - Actual: `PASS governance contract tests`.
  - Artifact: terminal evidence.
  - Blocking condition: None.
- Command: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/harness-aliases.tests.ps1`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: no active cold alias/lane remains.
  - Actual: `PASS harness alias tests`.
  - Artifact: terminal evidence.
  - Blocking condition: None.
- Command: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/harness.ps1 -Command context-validate -ContextPath scripts/.runs/b4060a73d0354372a104921154fcb204/context-pack.json -RequireCurrentBase`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: source hashes and accepted base are current.
  - Actual: `status=passed`.
  - Artifact: `scripts/.runs/b4060a73d0354372a104921154fcb204/context-pack.json`.
  - Blocking condition: None.
- Command: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/harness.ps1 -Command impact -ContextPath scripts/.runs/b4060a73d0354372a104921154fcb204/context-pack.json`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: exactly registered unit command IDs run against the normal
    checkout.
  - Actual: `status=passed`, shared seam `harness-control`.
  - Artifact: `scripts/.runs/9005f5c4492e48f49e4749eaa86f5ed9/outcome.json`.
  - Blocking condition: None.

## Manual QA

- QA level: QA-3
- Flow or step: fresh bounded dogfood dispatch after runner creation.
  - Status: Pass
  - Evidence type: ran
  - Owner: Milestone Orchestrator / fresh worker
  - Expected: worker changes only the manifest row, eval test, and
    out-of-scope fixture; deterministic grader decides the result.
  - Actual: worker returned only the three pinned paths and the grader passed
    with `HEVAL_PATH_OUT_OF_SCOPE`.
  - Blocking condition: None.
- Flow or step: artifact-only continuation.
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: checkpoint, Git/base, pack, and MNFS artifacts reconstruct the
    exact next action without portfolio transcript.
  - Actual: `F09-final-b7b884c` resumed as `quick_validation_passed` with next
    action `Milestone Orchestrator reviews the fixed F-09 commit`; task ID is
    correlation metadata only.
  - Blocking condition: None.

## Evidence

- Artifact: `scripts/.runs/b4060a73d0354372a104921154fcb204/context-pack.json`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Detail: 7,005 estimated tokens across 14 required bootstrap/selector
    sources. It is L2 and every source carries an `overflow_reason`; required
    reads are explicit, unrelated module sources are 0, and route misses are 0.
  - Blocking condition: None.
- Artifact: `scripts/.runs/9005f5c4492e48f49e4749eaa86f5ed9/outcome.json`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Detail: normal-checkout registered impact outcome, fake/contract target.
  - Blocking condition: None.
- Artifact: `scripts/.runs/state/checkpoints/F09-final-b7b884c.json`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Detail: canonical checkpoint and artifact-only resume proof.
  - Blocking condition: None.

## Risks

- The 7,005-token L2 pack exceeds the 2,000 target because its required
  bootstrap/validation and harness selectors are source-justified. Future work
  should reduce selectors only with a named route change, never by omitting a
  required safety source.
- This feature proves fake/contract behavior only. Independent milestone QA
  must retain real-target distinctions for any criterion that requires them.

## Handoff

- Current status: `quick_validation_passed`
- Next owner: Milestone Orchestrator
- Next action: Commit this fixed feature, perform one proportional L2 review,
  then accept/reject it and dispatch independent milestone QA on the fixed SHA.
- Required files/evidence: feature brief, spec, plan, validation, pack,
  impact outcome, checkpoint/resume record, and fresh-worker handoff.
- Blockers or open decisions: None.
