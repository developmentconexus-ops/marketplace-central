# F-01 Baseline Recovery — Quick Validation

```yaml
id: F-01
type: feature-validation
status: blocked
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
validation. The feature remains blocked: all 312 original paths are retained
as `retained-owner-needed` because this fresh build session cannot safely
infer ownership or run a cohort's scoped validation. No original path was
staged, removed, or altered.

## Quick Validation Result

- Result: Blocked
- Result owner: Feature Implementer
- Decision date: 2026-07-10
- Final feature state for handoff: blocked

## Quick Validation State

- fixup_attempts: 1 (PowerShell 5 compatibility correction in the new verifier)
- max_fixup_attempts: 1
- last_feature_validation_result: blocked by retained original state, not by a verifier defect

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
- Remaining blocker: 312 original paths remain retained pending confirmed owner
  attribution and scoped validation; feature status remains `blocked`.

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
- Remaining blocker: this fake-only correction does not establish a clean
  candidate baseline or resolve ownership/scoped validation for the 312
  retained original paths. F-01 remains `blocked`.

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
- Remaining blocker: 312 original paths remain `retained-owner-needed` pending
  confirmed ownership and scoped validation; F-01 remains `blocked`.

## Spec Adherence

- Spec satisfied: Blocked
- Deviations: No original cohort was committed.
- Reason: Ownership and scoped validation were not independently confirmed for
  any original path. The ledger preserves the exact 312-path inventory with an
  explicit retained-state record instead of guessing.

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
  - Change: one retained-owner-needed row for every original path.

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

- Artifact: `baseline-inventory.tsv`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Blocking condition: none.
- Artifact: `ownership-ledger.md`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Blocking condition: all original paths retained pending ownership.
- Artifact: ignored raw run directory
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Blocking condition: none.

## Risks

- The baseline is not clean and must not seed a worktree.
- `retained-owner-needed` rows must be resolved by a confirmed owner; committing
  or deleting them without scoped validation would violate F-01 constraints.

## Handoff

- Current status: `blocked`
- Next owner: Milestone Orchestrator
- Next action: assign ownership and validation commands to retained cohorts,
  then dispatch a fresh scoped reconciler for each confirmed cohort.
- Required files/evidence: feature brief, spec, plan, baseline inventory,
  ledger, raw run artifacts, and each cohort's future scoped validation.
- Blockers or open decisions: 312 retained original paths; no original path has
  a confirmed owner, completed scoped validation, or intentional commit SHA.
