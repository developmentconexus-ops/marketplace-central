# Local Runtime with Docker PostgreSQL Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Provide a repeatable Windows host workflow that runs Marketplace Central backend and frontend locally while Docker supplies only PostgreSQL.

**Architecture:** A PowerShell library owns safe `.env` parsing, prerequisite checks, PID/log state, PostgreSQL readiness, and lifecycle helpers. `scripts/dev-local.ps1` is a thin command dispatcher. The host Go server and Vite dev server inherit sanitized environment values; Docker is called only for `postgres`. A wiki runbook and a personal Codex skill make the workflow discoverable and repeatable.

**Tech Stack:** PowerShell 7, Go 1.25, Node/npm/Vite, Docker Compose (PostgreSQL only), PostgreSQL, Markdown, Codex skills.

## Global Constraints

- Work on Windows; do not use WSL.
- Preserve the heavily dirty worktree; never reset, revert, stage, or commit unrelated files.
- Do not print, persist, or copy secrets from `.env` into logs, wiki, test output, or skill files.
- Docker may start and inspect only the `postgres` service; backend and frontend must run on the Windows host.
- Never run `docker compose down`, remove containers, remove volumes, clear PostgreSQL data, restart Docker Desktop, or repair Docker from this workflow.
- `stop` may terminate only backend/frontend processes whose PIDs were written by this script.
- Use `GOCACHE=.gocache` for local Go commands.
- Keep `MPC_WEB_PROXY_TARGET=http://localhost:8080` for host Vite processes.
- Real provider and database validation must remain explicitly distinct from mocks and compile-only checks.

---

## File Structure

| File | Responsibility |
| --- | --- |
| `scripts/lib/dev-local-runtime.ps1` | Safe environment parsing, prerequisite checks, Docker-Postgres health, child-process/PID/log lifecycle, and local build/test helpers. |
| `scripts/dev-local.ps1` | Public `up`, `status`, `build`, `test`, and `stop` command dispatcher. |
| `scripts/tests/dev-local-runtime.tests.ps1` | Dependency-free PowerShell regression tests for parser, PID scoping, and no-secret status rendering. |
| `wiki/operations/local-runtime.md` | Operator runbook for daily start, update/build, diagnostics, Oracle live prerequisites, and scope limits. |
| `wiki/README.md` | Link to the local runtime runbook. |
| `package.json` | Convenience aliases for the five PowerShell commands. |
| `$CODEX_HOME/skills/local-mpc-runtime/SKILL.md` | Reusable Codex procedure for local runtime use. |
| `$CODEX_HOME/skills/local-mpc-runtime/agents/openai.yaml` | Generated skill UI metadata. |

### Task 1: Runtime library and deterministic regression tests

**Files:**
- Create: `scripts/lib/dev-local-runtime.ps1`
- Create: `scripts/tests/dev-local-runtime.tests.ps1`

**Interfaces:**
- Produces `Import-DevLocalEnvironment`, `ConvertTo-DevLocalEnvironmentValue`, `Get-DevLocalStateDirectory`, `Read-DevLocalPid`, `Write-DevLocalPid`, `Remove-DevLocalPid`, `Test-DevLocalPidOwnedByRuntime`, and `Get-DevLocalRedactedStatus`.
- Consumed by `scripts/dev-local.ps1` in Task 2.

- [ ] **Step 1: Write the failing parser and PID-scope tests**

Create `scripts/tests/dev-local-runtime.tests.ps1`:

```powershell
$ErrorActionPreference = 'Stop'
$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
. (Join-Path $repoRoot 'scripts/lib/dev-local-runtime.ps1')

function Assert-Equal([object]$Actual, [object]$Expected, [string]$Name) {
  if ($Actual -ne $Expected) { throw "$Name: expected [$Expected], got [$Actual]" }
}

Assert-Equal (ConvertTo-DevLocalEnvironmentValue '"quoted value"') 'quoted value' 'double quotes removed'
Assert-Equal (ConvertTo-DevLocalEnvironmentValue "'quoted value'") 'quoted value' 'single quotes removed'
Assert-Equal (ConvertTo-DevLocalEnvironmentValue 'abc=def') 'abc=def' 'equals preserved'

$state = Join-Path ([System.IO.Path]::GetTempPath()) ('mpc-dev-local-test-' + [guid]::NewGuid())
New-Item -ItemType Directory -Path $state | Out-Null
try {
  Write-DevLocalPid -StateDirectory $state -Name 'backend' -ProcessId $PID
  Assert-Equal (Read-DevLocalPid -StateDirectory $state -Name 'backend') $PID 'pid roundtrip'
  if (-not (Test-DevLocalPidOwnedByRuntime -StateDirectory $state -Name 'backend' -ProcessId $PID)) { throw 'owned PID was not recognized' }
  $status = Get-DevLocalRedactedStatus -StateDirectory $state
  if ($status -match 'MC_DATABASE_URL|postgres://') { throw 'status leaked configuration' }
} finally {
  Remove-Item -LiteralPath $state -Recurse -Force
}

Write-Output 'PASS dev-local-runtime tests'
```

- [ ] **Step 2: Run the test to verify RED**

Run:

```powershell
pwsh -NoProfile -File scripts/tests/dev-local-runtime.tests.ps1
```

Expected: FAIL because `scripts/lib/dev-local-runtime.ps1` does not exist.

- [ ] **Step 3: Implement the minimal safe library**

Create `scripts/lib/dev-local-runtime.ps1` with these exact safeguards:

```powershell
Set-StrictMode -Version Latest

function ConvertTo-DevLocalEnvironmentValue([string]$Value) {
  $trimmed = $Value.Trim()
  if ($trimmed.Length -ge 2 -and (($trimmed.StartsWith('"') -and $trimmed.EndsWith('"')) -or ($trimmed.StartsWith("'") -and $trimmed.EndsWith("'")))) {
    return $trimmed.Substring(1, $trimmed.Length - 2)
  }
  return $trimmed
}

function Import-DevLocalEnvironment([string]$RepositoryRoot) {
  $envPath = Join-Path $RepositoryRoot '.env'
  if (-not (Test-Path -LiteralPath $envPath -PathType Leaf)) { throw '.env is required for local runtime' }
  Get-Content -LiteralPath $envPath | ForEach-Object {
    if ($_ -notmatch '^\s*[^#\s][^=]*=') { return }
    $parts = $_ -split '=', 2
    [Environment]::SetEnvironmentVariable($parts[0].Trim().TrimStart([char]0xFEFF), (ConvertTo-DevLocalEnvironmentValue $parts[1]), 'Process')
  }
  if ([string]::IsNullOrWhiteSpace($env:MC_DATABASE_URL)) { throw 'MC_DATABASE_URL is required in .env' }
}

function Get-DevLocalStateDirectory([string]$RepositoryRoot) { Join-Path $RepositoryRoot '.tmp/dev-local' }
function Get-DevLocalPidPath([string]$StateDirectory, [string]$Name) { Join-Path $StateDirectory "$Name.pid" }
function Write-DevLocalPid([string]$StateDirectory, [string]$Name, [int]$ProcessId) {
  $process = Get-Process -Id $ProcessId -ErrorAction Stop
  $record = [pscustomobject]@{ process_id = $ProcessId; started_at_utc = $process.StartTime.ToUniversalTime().ToString('O') }
  New-Item -ItemType Directory -Force -Path $StateDirectory | Out-Null
  $record | ConvertTo-Json -Compress | Set-Content -LiteralPath (Get-DevLocalPidPath $StateDirectory $Name) -NoNewline
}
function Read-DevLocalProcessRecord([string]$StateDirectory, [string]$Name) { $path = Get-DevLocalPidPath $StateDirectory $Name; if (-not (Test-Path -LiteralPath $path)) { return $null }; return (Get-Content -LiteralPath $path -Raw | ConvertFrom-Json) }
function Read-DevLocalPid([string]$StateDirectory, [string]$Name) { $record = Read-DevLocalProcessRecord $StateDirectory $Name; if ($null -eq $record) { return $null }; return [int]$record.process_id }
function Remove-DevLocalPid([string]$StateDirectory, [string]$Name) { Remove-Item -LiteralPath (Get-DevLocalPidPath $StateDirectory $Name) -Force -ErrorAction SilentlyContinue }
function Test-DevLocalPidOwnedByRuntime([string]$StateDirectory, [string]$Name, [int]$ProcessId) {
  $record = Read-DevLocalProcessRecord $StateDirectory $Name
  if ($null -eq $record -or [int]$record.process_id -ne $ProcessId) { return $false }
  $process = Get-Process -Id $ProcessId -ErrorAction SilentlyContinue
  return $null -ne $process -and $process.StartTime.ToUniversalTime().ToString('O') -eq [string]$record.started_at_utc
}
function Get-DevLocalRedactedStatus([string]$StateDirectory) { return "state_directory=$StateDirectory backend_pid=$(Read-DevLocalPid $StateDirectory 'backend') frontend_pid=$(Read-DevLocalPid $StateDirectory 'frontend')" }
```

Add these concrete helpers in the same file. Every `throw` message names the
missing prerequisite or failed surface, never a secret value:

```powershell
function Assert-DevLocalPrerequisites([string]$RepositoryRoot) {
  foreach ($name in 'go','node','npm','docker') {
    if (-not (Get-Command $name -ErrorAction SilentlyContinue)) { throw "$name is required for local runtime" }
  }
  if (-not (Test-Path -LiteralPath (Join-Path $RepositoryRoot 'apps/server_core/go.mod'))) { throw 'apps/server_core/go.mod is required' }
  if (-not (Test-Path -LiteralPath (Join-Path $RepositoryRoot 'apps/web/package.json'))) { throw 'apps/web/package.json is required' }
}

function Start-DevLocalPostgres([string]$RepositoryRoot) {
  & docker compose up -d postgres
  if ($LASTEXITCODE -ne 0) { throw 'PostgreSQL Docker service did not start; inspect docker compose ps and PostgreSQL logs' }
}

function Wait-DevLocalPostgres([string]$RepositoryRoot) {
  $deadline = (Get-Date).AddMinutes(3)
  do {
    $health = (& docker inspect marketplace-central-postgres-1 --format '{{.State.Health.Status}}' 2>$null | Select-Object -Last 1).Trim()
    if ($health -eq 'healthy') { return }
    Start-Sleep -Seconds 2
  } while ((Get-Date) -lt $deadline)
  throw 'PostgreSQL did not become healthy in 3 minutes; inspect docker compose logs postgres. The script will not repair Docker or remove data.'
}

function Invoke-DevLocalMigrations([string]$RepositoryRoot) {
  $env:GOCACHE = Join-Path $RepositoryRoot '.gocache'
  & go run ./apps/server_core/cmd/migrate
  if ($LASTEXITCODE -ne 0) { throw 'Migrations failed; backend/frontend were not started' }
}

function Start-DevLocalChild([string]$RepositoryRoot, [string]$Name, [string]$FilePath, [string[]]$ArgumentList) {
  $state = Get-DevLocalStateDirectory $RepositoryRoot
  $stdoutLog = Join-Path $state "$Name.stdout.log"
  $stderrLog = Join-Path $state "$Name.stderr.log"
  $existing = Read-DevLocalPid $state $Name
  if ($existing -and (Get-Process -Id $existing -ErrorAction SilentlyContinue)) { throw "$Name is already running with PID $existing" }
  Remove-DevLocalPid $state $Name
  $process = Start-Process -FilePath $FilePath -ArgumentList $ArgumentList -WorkingDirectory $RepositoryRoot -WindowStyle Hidden -RedirectStandardOutput $stdoutLog -RedirectStandardError $stderrLog -PassThru
  Write-DevLocalPid $state $Name $process.Id
}

function Start-DevLocalBackend([string]$RepositoryRoot) {
  $command = "`$env:GOCACHE = Join-Path (Get-Location) '.gocache'; go run ./apps/server_core/cmd/server"
  Start-DevLocalChild $RepositoryRoot 'backend' 'pwsh' @('-NoProfile','-ExecutionPolicy','Bypass','-Command',$command)
}

function Start-DevLocalFrontend([string]$RepositoryRoot) {
  $env:MPC_WEB_PROXY_TARGET = 'http://localhost:8080'
  Start-DevLocalChild $RepositoryRoot 'frontend' 'npm.cmd' @('run','dev','--workspace','@marketplace-central/web','--','--host','127.0.0.1','--port','5174')
}

function Invoke-DevLocalBuild([string]$RepositoryRoot) {
  $env:GOCACHE = Join-Path $RepositoryRoot '.gocache'
  & go build ./apps/server_core/cmd/server; if ($LASTEXITCODE -ne 0) { throw 'Go build failed' }
  & npm run build --workspace @marketplace-central/web; if ($LASTEXITCODE -ne 0) { throw 'Vite build failed' }
}

function Invoke-DevLocalTests([string]$RepositoryRoot) {
  $env:GOCACHE = Join-Path $RepositoryRoot '.gocache'
  & go test ./apps/server_core/internal/modules/orders/... ./apps/server_core/internal/modules/profitability/... -count=1; if ($LASTEXITCODE -ne 0) { throw 'Go runtime tests failed' }
  & npm run test --workspace @marketplace-central/web -- --run; if ($LASTEXITCODE -ne 0) { throw 'Web tests failed' }
}

function Stop-DevLocalChildren([string]$RepositoryRoot) {
  $state = Get-DevLocalStateDirectory $RepositoryRoot
  foreach ($name in 'backend','frontend') {
    $processId = Read-DevLocalPid $state $name
    if ($processId -and (Test-DevLocalPidOwnedByRuntime $state $name $processId)) { Stop-Process -Id $processId -Force -ErrorAction SilentlyContinue }
    Remove-DevLocalPid $state $name
  }
}
```

Logs contain process output only; never emit an environment dump. PIDs are
recorded with the child start time so `stop` rejects a reused Windows PID.

- [ ] **Step 4: Run the regression test to verify GREEN**

Run:

```powershell
pwsh -NoProfile -File scripts/tests/dev-local-runtime.tests.ps1
```

Expected: `PASS dev-local-runtime tests`.

- [ ] **Step 5: Parse-check the library**

Run:

```powershell
[scriptblock]::Create((Get-Content scripts/lib/dev-local-runtime.ps1 -Raw)) | Out-Null
```

Expected: exit code 0.

### Task 2: Public local runtime command dispatcher

**Files:**
- Create: `scripts/dev-local.ps1`
- Modify: `package.json`
- Test: `scripts/tests/dev-local-runtime.tests.ps1`

**Interfaces:**
- Consumes all Task 1 functions.
- Produces CLI contract: `pwsh -File scripts/dev-local.ps1 -Command up|status|build|test|stop`.

- [ ] **Step 1: Add a failing command-dispatch assertion**

Append to `scripts/tests/dev-local-runtime.tests.ps1`:

```powershell
$unknown = & pwsh -NoProfile -File (Join-Path $repoRoot 'scripts/dev-local.ps1') -Command unknown 2>&1
if ($LASTEXITCODE -eq 0 -or ($unknown -join "`n") -notmatch 'unsupported command') { throw 'unknown command must fail without starting a service' }
```

- [ ] **Step 2: Run it to verify RED**

Run:

```powershell
pwsh -NoProfile -File scripts/tests/dev-local-runtime.tests.ps1
```

Expected: FAIL because `scripts/dev-local.ps1` does not exist.

- [ ] **Step 3: Implement command behavior**

Create the dispatcher with an explicit ValidateSet:

```powershell
param([ValidateSet('up','status','build','test','stop')][string]$Command = 'status')
$ErrorActionPreference = 'Stop'
$repositoryRoot = Split-Path -Parent $PSScriptRoot
. (Join-Path $repositoryRoot 'scripts/lib/dev-local-runtime.ps1')

switch ($Command) {
  'up' {
    Import-DevLocalEnvironment $repositoryRoot
    Assert-DevLocalPrerequisites $repositoryRoot
    Start-DevLocalPostgres $repositoryRoot
    Wait-DevLocalPostgres $repositoryRoot
    Invoke-DevLocalMigrations $repositoryRoot
    Start-DevLocalBackend $repositoryRoot
    Start-DevLocalFrontend $repositoryRoot
    Get-DevLocalRuntimeStatus $repositoryRoot
  }
  'status' { Get-DevLocalRuntimeStatus $repositoryRoot }
  'build' { Import-DevLocalEnvironment $repositoryRoot; Invoke-DevLocalBuild $repositoryRoot }
  'test' { Import-DevLocalEnvironment $repositoryRoot; Invoke-DevLocalTests $repositoryRoot }
  'stop' { Stop-DevLocalChildren $repositoryRoot }
}
```

`Start-DevLocalPostgres` runs exactly `docker compose up -d postgres`; it does
not use `down`, `rm`, `restart`, or Docker Desktop commands. `Wait-DevLocalPostgres` times out with a clear message and does not attempt recovery. Backend starts with `go run ./apps/server_core/cmd/server`; frontend starts with `npm run dev --workspace @marketplace-central/web -- --host 127.0.0.1 --port 5174`, inheriting `MPC_WEB_PROXY_TARGET=http://localhost:8080`. Build runs `go build ./apps/server_core/cmd/server` and `npm run build --workspace @marketplace-central/web`.

Add aliases to root `package.json`:

```json
"local:up": "pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/dev-local.ps1 -Command up",
"local:status": "pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/dev-local.ps1 -Command status",
"local:build": "pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/dev-local.ps1 -Command build",
"local:test": "pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/dev-local.ps1 -Command test",
"local:stop": "pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/dev-local.ps1 -Command stop"
```

- [ ] **Step 4: Run command regression and non-mutating status**

Run:

```powershell
pwsh -NoProfile -File scripts/tests/dev-local-runtime.tests.ps1
pwsh -NoProfile -File scripts/dev-local.ps1 -Command status
```

Expected: tests PASS; status reports state/PIDs/health without displaying any `.env` value and does not change Docker state.

- [ ] **Step 5: Validate script safety by search**

Run:

```powershell
rg -n 'compose (down|rm)|volume rm|desktop restart|docker system|Remove-Item.*postgres' scripts/dev-local.ps1 scripts/lib/dev-local-runtime.ps1
```

Expected: no matches.

### Task 3: Operator runbook and project command discovery

**Files:**
- Create: `wiki/operations/local-runtime.md`
- Modify: `wiki/README.md`
- Modify: `README.md`

**Interfaces:**
- Consumes Task 2 public command names.
- Produces a safe operator path for daily runtime and update validation.

- [ ] **Step 1: Write the failing documentation contract check**

Run:

```powershell
$required = 'local:up','local:status','local:build','local:test','local:stop','MPC_WEB_PROXY_TARGET','MC_DATABASE_URL','Docker only runs PostgreSQL'
$text = Get-Content wiki/operations/local-runtime.md -Raw
$missing = $required | Where-Object { $text -notmatch [regex]::Escape($_) }
if ($missing) { throw "missing runbook content: $($missing -join ', ')" }
```

Expected: FAIL because the runbook does not yet exist.

- [ ] **Step 2: Write the runbook and index links**

Write these exact operator sections:

- prerequisites without secret values;
- initial local dependency installation;
- daily `npm run local:up`, status, log locations, and browser URLs;
- safe `npm run local:stop` semantics;
- update workflow: `git pull` is operator-owned, then `npm ci`, `go mod download`, `npm run local:build`, `npm run local:test`, and `npm run local:up`;
- Oracle live prerequisites as an opt-in real-integration step;
- Docker limitations: PostgreSQL only, never invoke `down` or repair Docker from the runtime script;
- failure diagnostics with health and log commands;
- explicit statement that mocks do not prove live providers.

Add `[Local Runtime (host backend/frontend + Docker PostgreSQL)](operations/local-runtime.md)` under Operations in `wiki/README.md` and this concise pointer in root `README.md`:

```markdown
For Windows local runtime, use `npm run local:up`; Docker is used only for PostgreSQL. See [the local runtime runbook](wiki/operations/local-runtime.md).
```

- [ ] **Step 3: Run the documentation contract check**

Run the exact Step 1 command again.

Expected: exit code 0.

- [ ] **Step 4: Proofread and diff-check only scoped docs**

Run:

```powershell
git diff --check -- wiki/operations/local-runtime.md wiki/README.md README.md package.json
```

Expected: no whitespace errors.

### Task 4: Personal `local-mpc-runtime` skill

**Files:**
- Create: `$CODEX_HOME/skills/local-mpc-runtime/SKILL.md`
- Create: `$CODEX_HOME/skills/local-mpc-runtime/agents/openai.yaml`

**Interfaces:**
- Consumes `scripts/dev-local.ps1` and `wiki/operations/local-runtime.md`.
- Produces an automatically discoverable Codex skill for phrases such as
  “subir MPC local”, “rodar sem Docker”, “validar runtime local”, and
  “atualizar build local”.

- [ ] **Step 1: Initialize the skill with the supplied generator**

Resolve the target once before initialization; default to the discovered
personal skill root when `CODEX_HOME` is unset:

```powershell
$skillRoot = if ($env:CODEX_HOME) { Join-Path $env:CODEX_HOME 'skills' } else { 'C:\Users\leandro.theodoro\.codex\skills' }
```

Run from the Skill Creator directory:

```powershell
python scripts/init_skill.py local-mpc-runtime --path $skillRoot --interface display_name='Local MPC Runtime' --interface short_description='Run Marketplace Central locally with Docker PostgreSQL only.' --interface default_prompt='Use the local Marketplace Central runtime workflow.'
```

Expected: creates `$skillRoot/local-mpc-runtime` with `SKILL.md` and `agents/openai.yaml`.

- [ ] **Step 2: Write the concise skill instructions**

The skill must direct Codex to:

1. read `wiki/operations/local-runtime.md` first;
2. use `npm run local:status` before changing state;
3. use `npm run local:up` for host backend/frontend and Docker PostgreSQL only;
4. use `npm run local:build` after code changes and `npm run local:test` for validation;
5. use `npm run local:stop` only for script-owned host processes;
6. never call Docker cleanup/restart/volume removal through this skill;
7. keep secret values out of output;
8. label mock, PostgreSQL, Oracle, Mercado Livre, and browser evidence separately.

Use the following frontmatter:

```yaml
---
name: local-mpc-runtime
description: Run, build, test, update, or diagnose Marketplace Central locally on Windows while Docker supplies PostgreSQL only. Use when asked to start MPC without Dockerized backend/frontend, check local runtime health, run local builds/tests, or perform a safe local update.
---
```

- [ ] **Step 3: Validate the new skill**

Run:

```powershell
python scripts/quick_validate.py (Join-Path $skillRoot 'local-mpc-runtime')
```

Expected: validation passes with valid frontmatter and folder name.

- [ ] **Step 4: Verify generated metadata matches the skill**

Run:

```powershell
Get-Content (Join-Path $skillRoot 'local-mpc-runtime/agents/openai.yaml')
```

Expected: display name, short description, and default prompt describe the local MPC runtime without credentials or Docker cleanup.

### Task 5: End-to-end local runtime verification and handoff

**Files:**
- Modify: `wiki/operations/local-runtime.md` only if observed command behavior differs from the documented contract.

**Interfaces:**
- Consumes Tasks 1–4.
- Produces evidence that backend/frontend host startup is separate from Docker PostgreSQL.

- [ ] **Step 1: Check prerequisite/status behavior before start**

Run:

```powershell
npm run local:status
```

Expected: reports PostgreSQL and host process state without secret values or destructive Docker action.

- [ ] **Step 2: Build and test locally**

Run:

```powershell
npm run local:build
npm run local:test
```

Expected: Go and Vite build successfully; test output identifies skipped real-PostgreSQL tests if `MC_DATABASE_URL` is unavailable rather than claiming live evidence.

- [ ] **Step 3: Start and verify the host services**

Run:

```powershell
npm run local:up
Invoke-WebRequest -UseBasicParsing http://127.0.0.1:8080/healthz
Invoke-WebRequest -UseBasicParsing http://127.0.0.1:5174
```

Expected: backend returns HTTP 200; Vite returns HTTP 200; `docker compose ps` shows only PostgreSQL from this runtime command, not backend/frontend.

- [ ] **Step 4: Verify safe stop semantics**

Run:

```powershell
npm run local:stop
npm run local:status
```

Expected: script-owned backend/frontend PIDs are gone; PostgreSQL remains running; no Docker container or volume removal appears in output.

- [ ] **Step 5: Final scoped verification**

Run:

```powershell
pwsh -NoProfile -File scripts/tests/dev-local-runtime.tests.ps1
git diff --check -- scripts wiki/operations/local-runtime.md wiki/README.md README.md package.json
```

Expected: test PASS and no whitespace errors.
