# F-02-powershell-real-read-lanes

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-14
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-3
lifecycle_scope: feature
```

## Mission

MIS-001 MVP operator journey.

## Milestone

M-14 Real Vertical MVP Validation.

## Brief

Reuse existing harness lanes and close only proven gaps needed for PowerShell-safe
PostgreSQL, Mercado Livre, Oracle, redaction, and no-write validation.

## Inputs

- Existing harness lanes, `contracts/governance/execution-lanes.json`, the M-08
  runtime inventory, accepted module/live tests, and IC-004.
- M-14 validation criteria and F-01 no-write/scenario selection.

## Inputs/Outputs

- Validation records the actual PowerShell command, cwd, reviewed SHA, exit code,
  evidence path, and allowed read-only side effects.
- Go commands use absolute resolved GOCACHE and existing `./internal/...` paths from
  `apps/server_core`.

## Negative Scenarios

- Missing package/cwd/env name: preflight fails before live access.
- Unix `VAR=value command` syntax: registry validation rejects it.
- Provider method beyond approved read lane: no-write preflight rejects it.
- Command output contains secret/PII pattern: redact and fail evidence acceptance.

## Expected Output

QA has executable Windows commands for the proportional ladder and cannot silently
translate a broken command or widen live scope.

## Constraints

- Do not add a second CI, dependency installer, cache purge, or cold clone.
- Local secret values stay outside repository and command artifacts.
- Owned paths: only existing harness/testability files proved necessary by the
  plan, the test-only Mercado Livre live-read seam if required, and this feature
  root. Governance environment declarations remain read-only.
- Forbidden paths: production application behavior except testability seams explicitly accepted.

## Criteria IDs

- M-14-C01 Real source provenance.
- M-14-C03 PowerShell command validity.
- M-14-C04 Persistence and idempotency.
- M-14-C05 Evidence security and honesty.
- M-14-C06 No provider mutation.

## Validation Expectations

- Commands use absolute GOCACHE, existing paths, correct cwd, PowerShell syntax,
  sanitized output, and read-only side effects.
- Every live command preflights safely before final QA.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: Briefed.
- Next owner: Feature Implementer after F-01.
- Next action: Create spec/plan and close only demonstrated validation-lane gaps.
- Required files/evidence: harness schemas, IC-004, M-14 contract, and `validation.md`.
- Blockers or open decisions: Runtime capability conflict stops the lane.
