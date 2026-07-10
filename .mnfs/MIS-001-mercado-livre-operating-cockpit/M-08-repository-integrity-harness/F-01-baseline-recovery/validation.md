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
- Command: `powershell -NoProfile -ExecutionPolicy Bypass -File scripts/verify-baseline.ps1 -Inventory .mnfs/.../baseline-inventory.tsv -Ledger .mnfs/.../ownership-ledger.md`
  - Target: fake
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: exactly one valid ledger row for each original path.
  - Actual: exit 0; `inventory=312 committed=0 retained=312`.
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
