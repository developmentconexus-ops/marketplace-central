# F-01-baseline-recovery

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-08
created: 2026-07-10
updated: 2026-07-10
validation_level: QA-3
lifecycle_scope: feature
```

## Mission
MIS-001 Mercado Livre Operating Cockpit.

## Milestone
M-08 Repository Integrity and Deterministic Harness.

## Brief
Recover the current dirty checkout into an intentional, reviewable commit graph without losing or rewriting user-owned work.

## Inputs
- Current `git status`, tracked diff, untracked inventory, MNFS artifacts, and validation evidence.
- Existing commits and milestone ownership M-01 through M-06.

## Expected Output
- Path ownership ledger and superseded-artifact classification.
- Small intentional commits grouped by coherent completed work.
- Accepted clean baseline SHA with zero unstaged, staged, or untracked paths.

## Constraints
- No reset, revert, stash, clean, mass checkout, or deletion of unknown files.
- Do not fix product behavior while classifying baseline.
- Preserve unrelated user changes and secrets; `.env` is never staged.

## Negative Scenarios
- Unknown ownership: stop that path, record owner-needed; do not guess.
- Failing scoped validation: keep change unaccepted and record exact blocker.
- Secret/PII candidate: quarantine from staging and report path only.

## Validation Expectations
- Every original dirty path appears exactly once in ownership ledger.
- Accepted baseline reports `git status --short` empty and all commits have scoped test/evidence references.

## Execution Artifact Rules
`spec.md`, `plan.md`, and `validation.md` are created during feature execution.

## Handoff
- Current status: Briefed.
- Next owner: Feature Implementer.
- Next action: Create `spec.md` and `plan.md` from current git inventory.
- Required files/evidence: ownership ledger, commit list, validation.md.
- Blockers or open decisions: Unknown path ownership stops only affected commit group.

