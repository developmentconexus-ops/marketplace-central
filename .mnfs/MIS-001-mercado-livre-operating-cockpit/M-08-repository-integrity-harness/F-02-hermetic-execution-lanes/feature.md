# F-02-hermetic-execution-lanes

```yaml
id: F-02
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
Create explicit unit, integration, live, browser, and provider-write lanes with environment allowlists and no implicit side effects.

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
- Current status: Briefed.
- Next owner: Feature Implementer after F-01.
- Next action: Create `spec.md` and `plan.md` with exact command names.
- Required files/evidence: lane tests and environment-redaction evidence.
- Blockers or open decisions: F-01 baseline SHA.

