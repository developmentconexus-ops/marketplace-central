# Docker live-Oracle validation

`scripts/run-live-oracle-docker.ps1` is the canonical live Oracle validation
runner. It is a narrow, database-read-only lane; the Windows harness remains a
compatible fallback and is not configured or invoked here.

## Prerequisites

Docker must be available. The runner resolves these canonical values from the
ignored local `.env` file:

- `MPC_SANKHYA_ORACLE_USERNAME`
- `MPC_SANKHYA_ORACLE_PASSWORD`
- `MPC_SANKHYA_ORACLE_CONNECT_STRING`

The local `.env` parser accepts exactly those three assignments (blank lines
and comments are permitted). Any other assignment, duplicate key, or malformed
line stops the runner before Docker is invoked. It does not load Compose or
inherit a Compose-wide `.env`. Do not place credentials in shell history or a
command transcript.

For an approved secure handoff, a nonempty caller-process value with the same
exact name overrides its corresponding local `.env` value. No PowerShell
credential parameters are accepted. Generic `MPC_ORACLE_*` names and ambient
aliases are never read as caller input and cannot satisfy preflight. The runner
maps resolved values one way to the pre-existing governed container keys
`MPC_ORACLE_USERNAME`, `MPC_ORACLE_PASSWORD`, and
`MPC_ORACLE_CONNECT_STRING` used by the Oracle configuration. Those generic
names are a container boundary contract, never caller input aliases. Docker
receives key references rather than credential values as command arguments.

Before Docker preflight, build, or run, the runner clears the Docker child process
environment. It restores only a small Windows execution allowlist (`SystemRoot`,
`WINDIR`, `ComSpec`, `PATH`, `PATHEXT`, temporary-directory, and Docker user
configuration keys). The `docker version --format {{.Server.Version}}`
preflight has no runtime values at all. For the `docker run` child the runner
then adds exactly the five listed Oracle runtime keys. Ambient MPC, Oracle,
provider, database, and legacy alias values cannot enter the Docker child or the
container.

## Procedure

First perform the non-executing preflight:

```powershell
pwsh -NoProfile -File scripts/run-live-oracle-docker.ps1 -PreflightOnly
```

When it returns `status=ready`, run the target:

```powershell
pwsh -NoProfile -File scripts/run-live-oracle-docker.ps1
```

The runner builds `docker/dev/backend.Dockerfile`, mounts the checkout
read-only, and executes only:

```text
go test ./internal/modules/internal_read/adapters/oracle -run ^TestOracleLiveSmoke$ -count=1
```

Inside the container it injects exactly `MPC_ORACLE_USERNAME`,
`MPC_ORACLE_PASSWORD`, `MPC_ORACLE_CONNECT_STRING`, `MPC_ORACLE_LIVE_TEST=1`,
and `MPC_ORACLE_LIB_DIR=/opt/oracle/instantclient`. It starts no application
entrypoint, migration, Compose service, provider operation, or database write.

Docker command output is suppressed to avoid accidental secret disclosure. A
missing Docker installation, absent value, or live test failure must be
recorded as blocked/failed evidence; it is never a successful validation.
