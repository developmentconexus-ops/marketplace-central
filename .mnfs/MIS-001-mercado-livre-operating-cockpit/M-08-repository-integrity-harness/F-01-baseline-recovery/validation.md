# F-01 Baseline Recovery — Quick Validation

```yaml
id: F-01
type: feature-validation
status: quick_validation_passed
owner: Feature Implementer
parent: M-08
created: 2026-07-10
updated: 2026-07-10
validation_level: QA-3
lifecycle_scope: feature
```

## Feature ID

F-01-baseline-recovery

## Summary

The path-only baseline harness and its deterministic fixtures passed quick
validation. C01 mechanical readiness is now proven: all 312 original paths
have an explicit reconciliation disposition (296 committed and 16 retained
local raw artifacts). No original path was staged, removed, or altered by
F-01 finalization.

## Quick Validation Result

- Result: Passed (C01 mechanical readiness)
- Result owner: Feature Implementer
- Decision date: 2026-07-10
- Final feature state for handoff: quick_validation_passed

## Quick Validation State

- fixup_attempts: 1 (PowerShell 5 compatibility correction in the new verifier)
- max_fixup_attempts: 1
- last_feature_validation_result: quick_validation_passed by C01 mechanical gate; M-08 QA verdict remains unclaimed

## Superseded Historical Evidence

- The pre-correction command
  `powershell -NoProfile -ExecutionPolicy Bypass -File scripts/verify-baseline.ps1 -Inventory .mnfs/.../baseline-inventory.tsv -Ledger .mnfs/.../ownership-ledger.md`
  ran **without** `-AllowRetainedState` and recorded exit 0 with
  `inventory=312 committed=0 retained=312`. It is preserved only as historical
  evidence.
- That result is superseded by `f56a9ee` and `a021455`; it must not be read as
  current verifier behavior. Current controlled fake-fixture evidence is the
  explicit default-retained rejection and dirty-worktree coverage recorded in
  the two corrections below.

## Correction — Default Baseline Cleanliness

- Original quality finding: `scripts/verify-baseline.ps1:109` accepted retained
  ledger state and skipped the Git dirty check unless `-RequireCleanStatus` was
  provided.
- Target: fake (deterministic TSV fixtures and the local dirty worktree).
- RED command: `powershell -NoProfile -ExecutionPolicy Bypass -File scripts/test-verify-baseline.ps1`
  - Actual: exit 1; the retained-state fixture expected exit 1 but received 0.
  - Artifact: ignored `runs/20260710T-baseline-harness-correction/red-retained-default.log`.
- GREEN command: `powershell -NoProfile -ExecutionPolicy Bypass -File scripts/test-verify-baseline.ps1`
  - Actual: exit 0; default mode rejects retained ledger state and the dirty
    worktree. The historical `-AllowRetainedState` diagnostic claim is
    superseded by the correction below.
  - Artifact: ignored `runs/20260710T-baseline-harness-correction/green-test-verify-baseline.log`.
- Historical blocker: 312 paths were retained pending attribution. The final
  ledger maps 296 to intentional commits and retains 16 named local raw artifacts.

## Correction — Diagnostic Opt-In Does Not Waive Git Cleanliness

- Status: current semantics; this correction supersedes the earlier historical
  claim that `-AllowRetainedState` could accept a dirty diagnostic fixture.
- Root cause: the verifier used the retained-state opt-in to skip both the
  retained-disposition rejection and the Git dirty-worktree check; the former
  is waivable, the latter is not.
- Test isolation: `scripts/test-verify-baseline.ps1` now creates its own
  temporary Git repository, commits its fixtures, and introduces/removes only
  its controlled `controlled-dirty.txt` path. It does not use the ambient
  primary worktree as a test condition.
- Target: fake (temporary Git repository and deterministic TSV fixtures).
- RED command: `powershell -NoProfile -ExecutionPolicy Bypass -File scripts/test-verify-baseline.ps1`
  - Actual: exit 1; `dirty retained diagnostic baseline` expected exit 1 but
    received exit 0 before the verifier change.
  - Artifact: ignored `runs/20260710T-baseline-harness-correction/red-allow-retained-dirty-temp-repo.log`.
- GREEN command: `powershell -NoProfile -ExecutionPolicy Bypass -File scripts/test-verify-baseline.ps1`
  - Actual: exit 0; clean committed passes, dirty committed fails, dirty
    retained with `-AllowRetainedState` fails, and clean retained with
    `-AllowRetainedState` passes.
  - Artifact: ignored `runs/20260710T-baseline-harness-correction/green-allow-retained-dirty-temp-repo.log`.
- Historical limitation: this fake-only correction did not establish a clean
  candidate baseline. C01 now provides that separate mechanical gate; it does
  not claim functional validation for M-06.

## Correction — Isolated Fixture Git and Guarded Cleanup

- Status: current controlled fixture evidence.
- Root cause: the temporary fixture repository inherited ambient Git template,
  hook, and commit-signing configuration; cleanup also recursively targeted a
  generated path without proving fixture ownership first.
- RED command: `powershell -NoProfile -ExecutionPolicy Bypass -File scripts/test-verify-baseline.ps1`
  - Target: fake (fixture-owned temporary Git repository and deterministic TSV fixtures).
  - Actual: exit 1; the new local configuration assertion failed with
    `temporary fixture repository must set a local fixture-owned core.hooksPath.`
  - Artifact: ignored
    `runs/20260710T-baseline-harness-correction/red-fixture-git-isolation.log`.
- GREEN command: `powershell -NoProfile -ExecutionPolicy Bypass -File scripts/test-verify-baseline.ps1`
  - Target: fake (fixture-owned temporary Git repository and deterministic TSV fixtures).
  - Actual: exit 0; the fixture root is created before use, receives an
    ownership marker, initializes Git with its own empty template, configures
    local `core.hooksPath` to its empty fixture directory, and sets local
    `commit.gpgSign=false`. Recursive cleanup runs only while that created root
    and marker still exist.
  - Artifact: ignored
    `runs/20260710T-baseline-harness-correction/green-fixture-git-isolation.log`.
- Historical blocker: the prior ledger retained all 312 paths. The final ledger
  resolves every row to a commit or an explicit retained local raw artifact.

## Spec Adherence

- Spec satisfied: Passed for C01 mechanical readiness
- Deviations: 16 explicit local raw artifacts remain intentionally uncommitted.
- Reason: The ledger maps each inventory path exactly once: 296 to intentional
  reconciliation commits and 16 to their source-to-local M-06 evidence owner.

## Changes Made

- File: `.gitignore`
  - Change: ignores F-01 raw run artifacts.
- File: `scripts/test-verify-baseline.ps1`
  - Change: red/green fixture coverage for missing, duplicate, and invalid ledger rows.
- File: `scripts/verify-baseline.ps1`
  - Change: deterministic path/status, disposition, and clean-status verifier.
- File: `baseline-inventory.tsv`
  - Change: content-free original 312-path snapshot.
- File: `ownership-ledger.md`
  - Change: exact 312-row reconciliation map: 296 committed and 16 explicit
    retained-owner-needed local raw artifacts.

## Commands Run

- Command: `git status --porcelain=v1 --untracked-files=all`
  - Target: fake
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: 88 tracked plus 224 original untracked paths after excluding F-01 planning artifacts.
  - Actual: 88 tracked plus 226 untracked paths; `spec.md` and `plan.md` are the two controlled exclusions, leaving 312 original paths.
  - Artifact: ignored `runs/20260710T000000-baseline-capture/git-status-pre-f01.txt`, SHA-256 `DDEAEA7431D227AAABAAD8B59971624114EFC8366585B3AEA9B761E74AAB5290`.
  - Blocking condition: none.
- Command: `powershell -NoProfile -ExecutionPolicy Bypass -File scripts/test-verify-baseline.ps1`
  - Target: fake
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: missing, duplicate, and invalid fixtures reject; complete fixture accepts.
  - Actual: exit 0 after the required RED run (verifier absent) and GREEN implementation.
  - Artifact: ignored `runs/20260710T000000-baseline-capture/test-verify-baseline.log`.
  - Blocking condition: none.
- Command: `powershell -NoProfile -ExecutionPolicy Bypass -File scripts/verify-baseline.ps1 -Inventory <incomplete> -Ledger <incomplete>`
  - Target: fake
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: reject a ledger missing an inventory path.
  - Actual: exit 1 as expected.
  - Artifact: ignored `runs/20260710T000000-baseline-capture/verifier-incomplete.log`.
  - Blocking condition: none.
- Command: `[Historical — superseded] powershell -NoProfile -ExecutionPolicy Bypass -File scripts/verify-baseline.ps1 -Inventory .mnfs/.../baseline-inventory.tsv -Ledger .mnfs/.../ownership-ledger.md`
  - Target: fake
  - Status: Superseded
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: historical pre-correction check of exactly one valid ledger row
    for each original path.
  - Actual: historical pre-correction exit 0; `inventory=312 committed=0 retained=312`.
    This command omitted `-AllowRetainedState`; `f56a9ee` and `a021455`
    supersede it, so it is not current behavior.
  - Artifact: ignored `runs/20260710T000000-baseline-capture/verifier-complete-inventory.log`.
  - Blocking condition: retained records intentionally prevent clean-baseline readiness.
- Command: `git status --short` plus `verify-baseline.ps1 -RequireCleanStatus`
  - Target: fake
  - Status: Not run
  - Evidence type: could-not-run
  - Owner: Milestone Orchestrator
  - Expected: empty status and no retained state only at a candidate accepted SHA.
  - Actual: no candidate accepted SHA exists; 312 retained records remain.
  - Artifact: none.
  - Blocking condition: owner attribution and scoped validation are required before an original cohort can be committed.

## Manual QA

- QA level: QA-3
- Flow or step: Review staged F-01 paths before commit.
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Expected: only controlled F-01 artifacts and verifier files are staged.
- Actual: eight explicit paths staged; `git diff --cached --check` exited 0 and the staged path list contains no original dirty path or raw run artifact.
- Blocking condition: none.

## Evidence

## C01 Final Mechanical Gate

- Command: `powershell -NoProfile -ExecutionPolicy Bypass -File scripts/verify-baseline.ps1 -Inventory .mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/F-01-baseline-recovery/baseline-inventory.tsv -Ledger .mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/F-01-baseline-recovery/ownership-ledger.md -AllowRetainedState`
  - Target: fake (mechanical inventory/ledger reconciliation against the local Git worktree).
  - Actual: exit 0; `PASS: inventory=312 committed=296 retained=16`.
  - Artifact: ignored `runs/20260710T183900-c01-ledger-finalization/verify-baseline-allow-retained.log`.
- Command: `git status --short`
  - Target: fake (local Git worktree cleanliness check).
  - Actual: exit 0; visible status count 0.
  - Artifact: ignored `runs/20260710T183900-c01-ledger-finalization/git-status-short.log`.
- Decision: C01 readiness is proven mechanically, but an M-08 QA verdict is not claimed. M-06 remains functionally blocked; its scoped functional validation is not implied by this baseline reconciliation.

- Artifact: `baseline-inventory.tsv`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Blocking condition: none.
- Artifact: `ownership-ledger.md`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Blocking condition: none for C01 mechanical reconciliation; M-06 functional
    validation remains separately blocked.
- Artifact: ignored raw run directory
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Blocking condition: none.

## Risks

- C01 proves a clean Git worktree only with `-AllowRetainedState`; the 16 raw
  local artifacts are intentionally excluded from Git and must not be staged.
- M-06 remains functionally blocked; F-01 reconciliation does not replace its
  scoped validation.

## Handoff

- Current status: `quick_validation_passed`
- Next owner: final F-01 SPEC+QUALITY reviewers
- Next action: review the reconciliation ledger and C01 mechanical evidence;
  do not infer an M-08 QA verdict or M-06 functional completion.
- Required files/evidence: feature brief, spec, plan, baseline inventory,
  ledger, raw run artifacts, and each cohort's future scoped validation.
- Blockers or open decisions: M-06 remains functionally blocked; the 16 named
  local raw artifacts remain owned by its local-runtime evidence flow.
