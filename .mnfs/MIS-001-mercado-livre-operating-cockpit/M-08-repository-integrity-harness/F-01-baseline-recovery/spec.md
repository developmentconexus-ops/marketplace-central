# F-01 Baseline Recovery — Specification

```yaml
id: F-01
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: M-08
created: 2026-07-10
updated: 2026-07-10
validation_level: QA-3
lifecycle_scope: feature
```

## Feature ID

F-01-baseline-recovery

## Problem

The current `main` checkout contains completed and user-owned M-03 through M-06
work that is absent from `HEAD`. It cannot safely seed a worktree until every
pre-F-01 dirty path is preserved, classified, and represented by an intentional
commit or an explicit retained-state record.

The first read-only snapshot found 88 unstaged tracked paths, zero staged paths,
and 224 untracked paths (312 total). This is the original inventory. The
`spec.md` and `plan.md` files created by this planning checkpoint are controlled
F-01 exclusions and must not be counted as original dirty paths.

## Requirements

- Capture a machine-readable inventory of every original dirty path before
  classification, without reading suspected secret contents.
- Maintain an ownership ledger in which every original path appears exactly
  once and is assigned to a reviewed commit cohort or an explicit
  `owner-needed` retained-state record.
- Preserve all user-owned state: no reset, revert, stash, clean, mass
  checkout, deletion of an unknown path, or staging of `.env` files.
- Verify each commit cohort with the scoped evidence required by the work it
  contains before it is accepted.
- Finish only at an intentionally committed clean baseline SHA, with a verifier
  proving the ledger is complete and `git status --short` empty.
- Keep raw command output in an ignored run-ID directory; commit only concise,
  target-labelled evidence and path metadata without secret values or buyer PII.

## Non-Goals

- Fixing product behavior, legacy defects, M-09 through M-12 work, or provider
  listing data.
- Creating a worktree before the clean baseline is accepted.
- Treating mocks, fixtures, or a compile result as live Oracle, provider, or
  browser proof.
- Guessing the owner or purpose of an unknown path.

## Design

F-01 is a serial repository-reconciliation workflow. It records the starting
inventory by pathname and Git status only, then produces a committed ownership
ledger and a deterministic baseline verifier. Cohorts are reviewed one at a
time, validated with their existing scoped commands, staged by explicit path,
and committed only after review. A path that cannot be safely attributed stays
unstaged and is recorded as `owner-needed`; that prevents a clean-baseline
claim and stops the affected cohort rather than discarding it.

The final verifier compares the original inventory against the ledger, confirms
each path has one disposition, reports retained records, and requires a clean
Git status only when the feature is ready to hand off for M-08-C01 QA. Git and
verifier evidence is labelled `fake`; it does not assert any live integration.

## Edge Cases

- A path disappears from the working tree: stop, compare it with the initial
  inventory, and obtain owner direction; never restore or delete it by force.
- A likely secret or buyer-PII path is found: record only its path and
  disposition, quarantine it from staging, and block the relevant cohort.
- A scoped validation command fails: leave the cohort unaccepted and record the
  exact command and blocker; do not fold it into another commit.
- A new path appears after the initial snapshot: classify it separately as
  F-01-created or concurrent user state; it must not silently alter original
  inventory completeness.
- An unrelated user change arrives while reconciliation is running: preserve it
  as a new retained-state entry and do not stage it without owner confirmation.

## Acceptance Criteria

1. Every original dirty path from the 312-entry snapshot is present exactly
   once in the ownership ledger with a commit cohort or retained-state
   disposition.
   - Traces to milestone criterion ID: `M-08-C01 Baseline Integrity`.
   - Proven by: target `fake` baseline inventory verifier against the ledger.
2. No original path is lost, silently reverted, or unattributed during the
   reconciliation.
   - Traces to milestone criterion ID: `M-08-C01 Baseline Integrity`.
   - Proven by: target `fake` baseline inventory verifier and reviewed Git
     history/path evidence.
3. The accepted baseline SHA has no staged, unstaged, or untracked paths, and
   each original path is reachable through its recorded commit or an accepted
   retained-state record.
   - Traces to milestone criterion ID: `M-08-C01 Baseline Integrity`.
   - Proven by: target `fake` `git status --short` plus the baseline inventory
     verifier.

## Handoff

- Current status: `spec_ready`
- Next owner: Feature Implementer
- Next action: Execute the split implementation plan in a fresh build session.
- Required files/evidence: initial inventory, ownership ledger, scoped cohort
  evidence, verifier result, commit list, and final Git status.
- Blockers or open decisions: unknown-path ownership, suspected secret/PII
  disposition, and any cohort whose required validation cannot run.
