# F-03-ephemeral-postgres-migrations

```yaml
id: F-03
type: feature-brief
status: accepted
owner: Milestone Orchestrator
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
Run PostgreSQL tests against isolated run-scoped databases with canonical migrations, seeds, fixtures, and guaranteed cleanup.

## Inputs
- Migration runner, PostgreSQL repositories/tests, registry seeds, Compose runtime, and lane contract from F-02.

## Expected Output
- Generated `mpc_test_<run-id>` database plus run-labelled disposable Docker
  container, loopback random port, and `tmpfs` data directory.
- CWD-independent migration source, idempotence check, FK-complete fixture builders, and teardown in `finally`.

## Constraints
- Guard must reject database names outside `mpc_test_*`.
- Tests cannot use persistent dev volume, Compose dev project, caller database
  URL, fixed container name, or fixed host port.
- No production schema weakening to accommodate fixtures.

## Negative Scenarios
- Dev/live URL supplied: reject before migration.
- Test failure: cleanup still executes and records exit code.
- Missing seed/FK: fixture fails explicitly; do not bypass constraint.

## Validation Expectations
- Fresh DB applies all migrations; second run applies zero.
- Integration test creates then removes test resources; persistent dev row counts stay unchanged.

## Execution Artifact Rules
`spec.md`, `plan.md`, and `validation.md` are created during feature execution.

## Handoff

- Current status: `accepted`.
- Next owner: F-04 Feature Implementer.
- Next action: Prove cold provisioning, cold execution, and the remaining
  environment-dependent evidence without weakening the accepted F-03 lanes.
- Required files/evidence: `spec.md`, `plan.md`, `validation.md`, commits through
  `7ea94672`, real run summaries/transcripts, digest/count inventory, and zero
  resource inventory.
- Blockers or open decisions: None for F-03. M-08 remains in progress; F-04
  owns cold provisioning and Oracle/provider evidence.
