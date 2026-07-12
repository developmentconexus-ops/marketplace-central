# Docker live-Oracle validation

`scripts/run-live-oracle-docker.ps1` is the canonical live Oracle validation
runner. It is a narrow, database-read-only lane; the Windows harness remains a
compatible fallback and is not configured or invoked here.

## Prerequisites

Docker must be available and these canonical values must be nonempty in the
calling process, or supplied as PowerShell parameters:

- `MPC_ORACLE_USERNAME`
- `MPC_ORACLE_PASSWORD`
- `MPC_ORACLE_CONNECT_STRING`

No `.env` file, Compose service, or `SANKHYA_*` alias is read. Do not place
credentials in shell history, a command transcript, or a file. Supplying the
values via the caller process is preferred; the runner forwards them to Docker
by environment-key reference rather than as Docker command arguments.

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
