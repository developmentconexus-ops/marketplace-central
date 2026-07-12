# F-05-docker-live-oracle-runner

```yaml
id: F-05
type: feature-brief
status: accepted
owner: Feature Implementer
parent: M-06
created: 2026-07-12
updated: 2026-07-12
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-001 Mercado Livre Operating Cockpit.

## Milestone

M-06 Orders + Margin ML.

## Brief

Create the canonical Docker-based live-Oracle validation lane. It must use the
existing backend Docker image's Linux Go/CGO/Oracle Instant Client runtime,
inject only the explicit live-Oracle runtime values, and invoke only the
registered `TestOracleLiveSmoke` target. Windows remains an optional compatible
fallback, not the canonical runner.

## Inputs

- `docker/dev/backend.Dockerfile`
- `apps/server_core/internal/modules/internal_read/adapters/oracle/reader_live_test.go`
- `contracts/governance/execution-lanes.json`
- `contracts/governance/runtime-config.json`
- current Windows harness live preflight at `scripts/harness.ps1`

## Expected Output

- A narrow Docker command/profile and concise operator documentation.
- Explicit secret-safe runtime injection, with no Compose-wide `.env` use.
- A regression check proving the runner selects the registered read-only Go test
  and cannot start migrations or the application server.

## Constraints

- Docker is canonical; Windows is optional fallback only.
- Allow only the live-Oracle lane variables needed by the existing Oracle
  config/test, including `MPC_ORACLE_LIVE_TEST=1` and the in-container library
  directory.
- No Oracle, provider, or MPC database writes; no migrations; no server boot;
  no application Compose services; no broad `.env` use; no secrets in files,
  logs, tests, or validation evidence.
- Do not modify the existing harness control-plane seam, Oracle reader/domain
  behavior, Docker development stack, OpenAPI/SDK, or migrations.

## Validation Expectations

- Run focused deterministic validation for the runner/profile and documentation.
- Build/run the live target only after Docker availability and safe explicit
  injection are verified; otherwise record a truthful blocked live-run result.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution.

## Handoff

- Current status: accepted.
- Next owner: Milestone Orchestrator for fixed-SHA review and proportional QA.
- Next action: retain the Docker daemon blocker as pending operational evidence;
  when access is restored, rerun the Docker preflight and then the canonical
  runner with credentials supplied only through the calling process or secure
  parameters.
- Required files/evidence: `spec.md`, `plan.md`, `validation.md`, focused
  runner checks, targeted Oracle config result, and any safe live-target result.
- Blockers or open decisions: Docker daemon/API remains unavailable; the
  implementation is accepted, but the actual live lane remains unproven and
  cannot be represented as a pass.
