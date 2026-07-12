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

The Docker live-Oracle runner needs a secure local credential handoff without
Compose-wide `.env` inheritance. The correction must accept only the three
explicit `MPC_SANKHYA_ORACLE_*` inputs from an ignored local `.env` or the
same exact caller-process keys, map them one way to the governed container
`MPC_ORACLE_*` names, and remain a canonical Linux/CGO/Instant Client execution
path that cannot invoke an application runtime.

## Requirements

- Build the existing `docker/dev/backend.Dockerfile` and run only
  `TestOracleLiveSmoke` in an isolated Docker container.
- Load an ignored local `.env` using a narrow parser that accepts exactly
  `MPC_SANKHYA_ORACLE_USERNAME`, `MPC_SANKHYA_ORACLE_PASSWORD`, and
  `MPC_SANKHYA_ORACLE_CONNECT_STRING`; reject every other `.env` key. Caller
  process values may be used only for those same exact three keys and take
  precedence over the corresponding local `.env` values. Reject generic
  `MPC_ORACLE_*` and every ambient alias for container handoff. Map only the
  resolved values to the pre-existing
  governed container names `MPC_ORACLE_USERNAME`, `MPC_ORACLE_PASSWORD`, and
  `MPC_ORACLE_CONNECT_STRING`; inject by key, never as Docker argument values
  or persisted content.
- Set `MPC_ORACLE_LIVE_TEST=1` and the image's `/opt/oracle/instantclient` as
  `MPC_ORACLE_LIB_DIR`; do not use Compose-wide `.env` inheritance, aliases,
  migrations, or an application entrypoint.
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

`scripts/run-live-oracle-docker.ps1` parses the ignored local `.env` narrowly:
it accepts precisely the three Sankhya credential assignments and fails closed
on any other key. The same exact caller-process keys override their local
counterparts. It creates a fresh child process environment containing resolved
values under the governed Oracle container keys, then calls Docker with key-only
`--env` switches. This one-way boundary mapping is not a generic input
fallback. It builds an ephemeral tagged image from the existing backend Dockerfile and uses `go test` with the exact package and `-run
^TestOracleLiveSmoke$`. A test seam returns the constructed invocation without
launching Docker so Pester fixtures never require credentials or Docker.

## Edge Cases

- Missing or whitespace credentials: stop before Docker build or run and list
  only missing key names.
- Missing Docker: stop before credentials are used for a container launch.
- Caller-process values take precedence over matching local `.env` values;
  generic names and ambient aliases are rejected and can never satisfy the
  credential preflight. Any non-whitelisted `.env` key stops before Docker is
  launched.
- A live test failure is surfaced as a failed execution, never rewritten as a
  success.

## Acceptance Criteria

### F05-AC01

The canonical runner invokes only the registered read-only Oracle test with the
existing Dockerfile, exactly five allowlisted runtime keys, and no Compose,
migration, or server command. Traces to milestone criterion ID: `M-06-C02`.
Proven by: `docker-live-oracle-runner-tests`.

### F05-AC02

The runner preflights Docker and resolved Sankhya credentials, and deterministic
fixtures prove the `.env` whitelist, caller-process precedence, generic and
ambient-alias rejection, and that values are absent from Docker arguments and
test output.
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
