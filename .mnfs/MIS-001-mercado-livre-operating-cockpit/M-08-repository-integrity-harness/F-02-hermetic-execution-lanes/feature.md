# F-02-hermetic-execution-lanes

```yaml
id: F-02
type: feature-brief
status: superseded
owner: Mission Strategist
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
Historical denylist-based lane spike. Its blocked validation remains evidence;
the accepted allowlisted child runtime is owned by F-08.

## Inputs
- Current scripts, Compose files, backend entrypoint, Go tests, npm workspaces, and `.env` consumers.

## Expected Output
- Versioned Windows-first command surface delegated by root npm aliases.
- Unit lane consumes no `.env`; live lanes consume only named keys; provider writes use a separate actor-required command.

## Constraints
- No library/framework migration.
- No provider-specific business rule in harness.
- Do not mask current inventory or contract failures.

## Negative Scenarios
- Unit lane detects DB/network/live variable: fail before tests.
- Live configuration incomplete: fail with missing key names only.
- Provider-write command lacks actor/idempotency: reject before network.

## Validation Expectations
- Unit preflight proves zero external targets and exits 0 with live `.env` present but ignored.
- Live preflight lists target type and key names without values.

## Execution Artifact Rules
`spec.md`, `plan.md`, and `validation.md` are created during feature execution.

## Handoff
- Current status: Superseded, not accepted; blocked validation is unchanged.
- Next owner: None; F-08 owns the replacement architecture.
- Next action: Do not reopen the denylist implementation.
- Required files/evidence: Existing F-02 spec, plan, and validation as history.
- Blockers or open decisions: None.
