# F-06-knowledge-authority-cutover

```yaml
id: F-06
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

Establish one current owner for architecture, machine governance, human
runbooks, and execution state; migrate only current unique ADRs and remove
`.brain` without a compatibility layer.

## Inputs

- `AGENTS.md`, `ARCHITECTURE.md`, wiki, `.mnfs`, `.brain`, current M04-M06
  validation results, and the approved repository-native harness design.

## Expected Output

- Current ADRs 004-006 under `docs/architecture/decisions/` with an index.
- Active startup/truth/update guidance points to architecture, wiki, `.mnfs`,
  and the reserved machine-governance boundary.
- Architecture module status reconciled from accepted MNFS evidence.
- Zero tracked `.brain` paths; superseded history remains in Git.

## Constraints

- Do not copy superseded ADRs 001-003, stale pulse/roadmap/session state, or the
  completed legacy plan.
- Historical evidence may retain truthful `.brain` path references.
- Do not renumber ADRs or rewrite M04-M06 verdicts.

## Negative Scenarios

- A current unique decision has no new owner: block deletion.
- Active guidance still requires `.brain`: fail the active-reference gate.
- Architecture claims M06 passed: fail; implementation exists but milestone
  remains blocked on required live evidence.

## Validation Expectations

- `git ls-files -- .brain` emits no path.
- Active-reference scan emits no `.brain` or `Nexus Brain` match.
- ADR index lists exactly current decisions 004, 005, and 006.
- Architecture distinguishes passed M04/M05 from implemented-but-blocked M06.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution.

## Handoff

- Current status: Briefed.
- Next owner: Feature Implementer.
- Next action: Create `spec.md` and `plan.md` for the atomic authority cutover.
- Required files/evidence: migration ledger, active-reference scan, tracked-file
  scan, and diff review.
- Blockers or open decisions: None; operator approved `.brain` removal.
