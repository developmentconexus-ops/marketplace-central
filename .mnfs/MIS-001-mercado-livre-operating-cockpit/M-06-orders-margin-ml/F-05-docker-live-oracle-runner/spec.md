# F-05 Docker live-Oracle runner specification

```yaml
id: F-05
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: F-05
created: 2026-07-12
updated: 2026-07-12
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-05-docker-live-oracle-runner

## Problem

The Docker live-Oracle runner still expects the obsolete generic
`MPC_ORACLE_*` credential namespace, while its provisioned caller process uses
the explicit `MPC_SANKHYA_ORACLE_*` namespace. The correction must remain a
canonical Linux/CGO/Instant Client execution path that cannot inherit Compose
configuration or invoke an application runtime.

## Requirements

- Build the existing `docker/dev/backend.Dockerfile` and run only
  `TestOracleLiveSmoke` in an isolated Docker container.
- Accept only `MPC_SANKHYA_ORACLE_USERNAME`,
  `MPC_SANKHYA_ORACLE_PASSWORD`, and `MPC_SANKHYA_ORACLE_CONNECT_STRING` from
  explicit parameters or the caller process. Reject generic `MPC_ORACLE_*` and
  ambient aliases; inject accepted values into Docker by key, never as Docker
  argument values or persisted content.
- Set `MPC_ORACLE_LIVE_TEST=1` and the image's `/opt/oracle/instantclient` as
  `MPC_ORACLE_LIB_DIR`; do not load `.env`, aliases, Compose, migrations, or an
  application entrypoint.
- Preflight Docker and nonblank canonical credentials before build/run. A
  missing prerequisite is a truthful non-success result.
- Add deterministic Pester coverage for command shape, exact environment
  allowlist, and secret-safe argument handling, plus operator instructions.

## Non-Goals

- Changing the Oracle test, Docker development stack, Compose, governance, or
  Windows harness.
- Connecting to Oracle when Docker or credentials are unavailable, or claiming
  that an unrun live test passed.
- Running migrations, app/server processes, providers, or database writes.

## Design

`scripts/run-live-oracle-docker.ps1` creates a fresh child process environment
containing only the explicit Sankhya credential keys, then calls Docker with key-only
`--env` switches. It builds an ephemeral tagged image from the existing backend
Dockerfile and uses `go test` with the exact package and `-run
^TestOracleLiveSmoke$`. A test seam returns the constructed invocation without
launching Docker so Pester fixtures never require credentials or Docker.

## Edge Cases

- Missing or whitespace credentials: stop before Docker build or run and list
  only missing key names.
- Missing Docker: stop before credentials are used for a container launch.
- Explicit parameter and process values: explicit parameters take precedence;
  generic names, ambient aliases, and `.env` are rejected or ignored as
  applicable and can never satisfy the credential preflight.
- A live test failure is surfaced as a failed execution, never rewritten as a
  success.

## Acceptance Criteria

### F05-AC01

The canonical runner invokes only the registered read-only Oracle test with the
existing Dockerfile, exactly five allowlisted runtime keys, and no Compose,
migration, or server command. Traces to milestone criterion ID: `M-06-C02`.
Proven by: `docker-live-oracle-runner-tests`.

### F05-AC02

The runner preflights Docker and explicit Sankhya credentials, and deterministic
fixtures prove that generic and ambient aliases cannot satisfy preflight and
that values are absent from Docker arguments and test output.
Traces to milestone criterion ID: `M-06-C02`. Proven by:
`docker-live-oracle-runner-tests`; a live execution is recorded separately only
when its prerequisites are actually available.

## Handoff

- Current status: `spec_ready`
- Next owner: Feature Implementer
- Next action: Write `plan.md`, compile and validate the scoped context pack,
  then implement.
- Required files/evidence: runner, Pester test, documentation, context pack,
  and validation record.
- Blockers or open decisions: None.
