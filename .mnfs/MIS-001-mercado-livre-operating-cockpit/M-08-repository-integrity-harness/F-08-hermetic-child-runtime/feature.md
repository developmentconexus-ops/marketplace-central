# F-08-hermetic-child-runtime

```yaml
id: F-08
type: feature-brief
status: accepted
owner: Feature Implementer
parent: M-08
created: 2026-07-10
updated: 2026-07-11
validation_level: QA-3
lifecycle_scope: feature
```

## Mission

MIS-001 Mercado Livre Operating Cockpit.

## Milestone

M-08 Repository Integrity and Deterministic Harness.

## Brief

Replace F-02's incomplete inherited-environment denylist with fresh
lane-specific child environments and split the monolithic runner into focused,
typed PowerShell modules without changing public command aliases.

## Inputs

- F-02 blocked evidence, accepted F-07 execution/runtime contracts, current
  harness script/tests, Go/npm workspace commands, and Windows process model.

## Expected Output

- Thin stable `scripts/harness.ps1` dispatcher.
- Focused policy, environment, execution, and context modules under
  `scripts/harness/`; evidence and state remain owned by F-04 and F-05.
- Unit child environment built from a small safe OS/tool allowlist.
- Existing npm harness aliases remain compatible and are behavior-tested.

## Constraints

- Do not fix F-02 by adding missed keys to another denylist.
- Child execution must be CWD-independent and restore no ambient environment
  into the parent.
- Business/provider rules remain outside the harness.

## Negative Scenarios

- Any unapproved inherited variable reaches unit subprocess: fail fixture.
- Subprocess output contains a configured secret/candidate credential: fail
  redaction gate.
- Working directory differs from repository root: public command still resolves
  canonical paths.
- Unsupported lane or target: fail before process start.

## Validation Expectations

- Sentinel external variables, including arbitrary previously unknown names,
  are absent inside the unit child process.
- Required safe tool variables, repository-local `GOCACHE`/`GOMODCACHE`, and
  network-off Go settings are present only in the child; parent values remain
  unchanged.
- Existing lane tests and all root aliases exercise behavior, not name presence,
  and report expected target classes.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` record the completed feature execution
and focused fake/contract evidence.

## Handoff

- Current status: `accepted`.
- Next owner: F-03 Feature Implementer.
- Next action: Implement the explicit ephemeral PostgreSQL lifecycle on the
  accepted typed environment/process boundary.
- Required files/evidence: `spec.md`, `plan.md`, `validation.md`, commits
  `4beb9d7` through the Phase 4 evidence commit, focused suites, normal unit run,
  and final current-base context pack.
- Blockers or open decisions: No F-08 blocker. F-04 still owns cold dependency
  provisioning and rich evidence; M-08 remains in progress.
