# F-01 Baseline Recovery — Plan

```yaml
id: F-01
type: feature-plan
status: planned
owner: Feature Implementer
parent: M-08
created: 2026-07-10
updated: 2026-07-10
validation_level: QA-3
lifecycle_scope: feature
split_decision: split
split_reason: More than six ordered steps and unresolved path ownership require a fresh build session.
```

## Feature ID

F-01-baseline-recovery

## Execution Boundary

This plan is planning-only. It authorizes no baseline reconciliation until a
fresh Feature Implementer receives it in `build` mode. The current checkout
stays serial; no worktree may be created. Destructive Git commands and unknown
file deletion remain prohibited.

## Steps

1. Capture the original inventory as Git status metadata in an ignored
   run-ID directory, excluding only the known F-01 `spec.md` and `plan.md`
   planning artifacts. Record the command, SHA, branch, path count, and a
   content-free path/status list. Label this evidence `fake`.
2. Add a focused failing test or fixture for the baseline reconciler that
   rejects a missing original path, duplicate ledger row, or unrecognised
   disposition. Run it first and record the expected failure.
3. Create the minimal deterministic baseline verifier and concise ownership
   ledger format. It must consume the captured inventory, accept only explicit
   dispositions, and emit no file contents or environment values.
4. Run the verifier against a deliberately incomplete fixture and confirm it
   fails; then run it against the complete captured inventory and confirm it
   reports every original path exactly once.
5. Classify paths into reviewed M-03 through M-06/M-08 cohorts, generated or
   historical documentation, and `owner-needed` records. A suspected secret or
   PII path is path-only, quarantined, and blocks its cohort; no ownership is
   inferred.
6. For one cohort at a time, inspect only its directly relevant artifact and
   run its scoped validation. Stage explicit paths only, review the staged
   diff, and create an intentional commit with the evidence reference. Do not
   mix shared seams or unrelated user changes.
7. For an unresolved or failed cohort, preserve all files unstaged and append
   an explicit retained-state record with owner, reason, exact failed command
   and blocker, and required next validation. Do not claim the clean baseline
   while any retained state exists.
8. Re-run the full baseline verifier after each accepted cohort; verify the
   ledger’s original inventory has no missing or duplicate path and record the
   associated commit SHA.
9. At the candidate accepted SHA, run `git status --short` and the verifier.
   Produce M-08-C01 mechanical evidence only when the visible checkout is
   clean and every original path resolves to an intentional commit or an
   explicit retained local raw artifact accepted with `-AllowRetainedState`.
   Without that opt-in, retained state blocks the verifier; QA Validator alone
   decides the M-08-C01 verdict.

## Files Expected To Change

- Create: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/F-01-baseline-recovery/ownership-ledger.md`
  - Reason: committed, path-only disposition record for the original inventory.
- Create: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/F-01-baseline-recovery/baseline-inventory.tsv`
  - Reason: committed content-free original inventory consumed by the verifier.
- Create: `scripts/verify-baseline.ps1`
  - Reason: deterministic Windows-native ledger and Git-status verifier.
- Create: `scripts/test-verify-baseline.ps1`
  - Reason: fixture-based TDD coverage for missing, duplicate, and invalid
    inventory dispositions.
- Modify: `.gitignore`
  - Reason: ignore raw run-ID logs while retaining concise ledger/evidence.
- Modify: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/F-01-baseline-recovery/validation.md`
  - Reason: quick-validation evidence and handoff after the build session.
- Modify: existing dirty paths only after their ledger cohort, scoped evidence,
  and staging set are explicitly reviewed.

## Verification Commands

- Command: `powershell -ExecutionPolicy Bypass -File scripts/test-verify-baseline.ps1`
  - Target: `fake`
  - Satisfies: acceptance criteria 1 and 2 / `M-08-C01 Baseline Integrity`.
  - Expected result: the incomplete, duplicate, and invalid fixtures fail; the
    complete fixture passes without reading path contents or environment values.
- Command: `powershell -ExecutionPolicy Bypass -File scripts/verify-baseline.ps1 -Inventory <path> -Ledger <path>`
  - Target: `fake`
  - Satisfies: acceptance criteria 1, 2, and 3 / `M-08-C01 Baseline Integrity`.
  - Expected result: each original path occurs exactly once, every committed
    SHA resolves and contains the recorded path in its commit tree/diff, all
    dispositions are valid, and the command passes only at a clean accepted
    SHA. Explicit retained local raw artifacts require `-AllowRetainedState`.
- Command: `git status --short`
  - Target: `fake`
  - Satisfies: acceptance criterion 3 / `M-08-C01 Baseline Integrity`.
  - Expected result: no output at the candidate accepted baseline SHA.
- Command: `git log --oneline --decorate -n <cohort-count>`
  - Target: `fake`
  - Satisfies: acceptance criteria 1 and 2 / `M-08-C01 Baseline Integrity`.
  - Expected result: each accepted ledger cohort names its intentional commit.

## QA Steps

- Review the original inventory count and verify the two planning artifacts are
  the only controlled exclusions.
- Review each staged cohort before commit; confirm no `.env`, secret value, or
  buyer PII is staged and no unrelated path appears.
- Independently re-run the verifier and final Git status before forwarding
  M-08-C01 evidence to the Milestone Orchestrator.

## Rollback/Risk Notes

- There is no destructive rollback. A failed cohort remains preserved and
  unstaged with its evidence and owner-needed record.
- Commit only paths already verified in the ledger; correct a bad ledger entry
  with a new intentional commit rather than rewriting history.
- Raw logs stay in an ignored run-ID directory. Committed artifacts contain
  paths, command classifications, and SHA references only.

## Handoff

- Current status: `planned`
- Next owner: Feature Implementer
- Next action: Fresh `build` session executes steps 1–9 and produces
  `validation.md`.
- Required files/evidence: `spec.md`, this plan, original inventory, ledger,
  cohort validations, verifier output path, commit list, and final Git status.
- Blockers or open decisions: any unknown path owner, secret/PII candidate,
  or validation failure blocks the affected cohort and prevents baseline
  acceptance.
