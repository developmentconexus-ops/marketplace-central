# F-03-ephemeral-postgres-migrations

```yaml
id: F-03
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
Run PostgreSQL tests against isolated run-scoped databases with canonical migrations, seeds, fixtures, and guaranteed cleanup.

## Inputs
- Migration runner, PostgreSQL repositories/tests, registry seeds, Compose runtime, and lane contract from F-02.

## Expected Output
- Generated `mpc_test_<run-id>` database/project namespace.
- CWD-independent migration source, idempotence check, FK-complete fixture builders, and teardown in `finally`.

## Constraints
- Guard must reject database names outside `mpc_test_*`.
- Tests cannot use persistent dev volume or fixed container name.
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
- Current status: Briefed.
- Next owner: Feature Implementer after F-02.
- Next action: Create `spec.md` and `plan.md` for ephemeral DB lifecycle.
- Required files/evidence: migration/idempotence/cleanup transcripts.
- Blockers or open decisions: Exact command names frozen by F-02.

