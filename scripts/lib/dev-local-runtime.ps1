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
function Read-DevLocalProcessRecord([string]$StateDirectory, [string]$Name) { $path = Get-DevLocalPidPath $StateDirectory $Name; if (-not (Test-Path -LiteralPath $path)) { return $null }; return (Get-Content -LiteralPath $path -Raw | ConvertFrom-Json -DateKind String) }
function Read-DevLocalPid([string]$StateDirectory, [string]$Name) { $record = Read-DevLocalProcessRecord $StateDirectory $Name; if ($null -eq $record) { return $null }; return [int]$record.process_id }
function Remove-DevLocalPid([string]$StateDirectory, [string]$Name) { Remove-Item -LiteralPath (Get-DevLocalPidPath $StateDirectory $Name) -Force -ErrorAction SilentlyContinue }
function Test-DevLocalPidOwnedByRuntime([string]$StateDirectory, [string]$Name, [int]$ProcessId) {
  $record = Read-DevLocalProcessRecord $StateDirectory $Name
  if ($null -eq $record -or [int]$record.process_id -ne $ProcessId) { return $false }
  $process = Get-Process -Id $ProcessId -ErrorAction SilentlyContinue
  return $null -ne $process -and $process.StartTime.ToUniversalTime().ToString('O') -eq [string]$record.started_at_utc
}
function Get-DevLocalRedactedStatus([string]$StateDirectory) { return "state_directory=$StateDirectory backend_pid=$(Read-DevLocalPid $StateDirectory 'backend') frontend_pid=$(Read-DevLocalPid $StateDirectory 'frontend')" }

function Get-DevLocalHttpHealth([string]$Uri) {
  try {
    $response = Invoke-WebRequest -Uri $Uri -TimeoutSec 2 -UseBasicParsing -ErrorAction Stop
    if ($response.StatusCode -ge 200 -and $response.StatusCode -lt 300) { return 'healthy' }
    return "http_$($response.StatusCode)"
  } catch {
    return 'unavailable'
  }
}

function Ensure-DevLocalDockerPath() {
  if (Get-Command docker -ErrorAction SilentlyContinue) { return }
  if (-not $IsWindows) { return }

  $dockerDesktopBin = Join-Path $env:ProgramFiles 'Docker\Docker\resources\bin'
  if (-not (Test-Path -LiteralPath $dockerDesktopBin -PathType Container)) { return }

  $env:PATH = "$dockerDesktopBin$([IO.Path]::PathSeparator)$env:PATH"
  [void](Get-Command docker -ErrorAction SilentlyContinue)
}

function Get-DevLocalPostgresHealth() {
  try {
    $health = [string]((& docker inspect marketplace-central-postgres-1 --format '{{.State.Health.Status}}' 2>$null | Select-Object -Last 1))
    if ([string]::IsNullOrWhiteSpace($health)) { return 'unavailable' }
    return $health.Trim()
  } catch {
    return 'unavailable'
  }
}

function Get-DevLocalRuntimeStatus([string]$RepositoryRoot) {
  Ensure-DevLocalDockerPath
  $state = Get-DevLocalStateDirectory $RepositoryRoot
  $redactedState = Get-DevLocalRedactedStatus $state
  $backendHealth = Get-DevLocalHttpHealth 'http://127.0.0.1:8080/healthz'
  $frontendHealth = Get-DevLocalHttpHealth 'http://127.0.0.1:5174'
  $postgresHealth = Get-DevLocalPostgresHealth
  return "$redactedState backend_health=$backendHealth frontend_health=$frontendHealth postgres_health=$postgresHealth"
}

function Assert-DevLocalPrerequisites([string]$RepositoryRoot) {
  Ensure-DevLocalDockerPath
  foreach ($name in 'go','node','npm','docker') {
    if (-not (Get-Command $name -ErrorAction SilentlyContinue)) { throw "$name is required for local runtime" }
  }
  if (-not (Test-Path -LiteralPath (Join-Path $RepositoryRoot 'apps/server_core/go.mod'))) { throw 'apps/server_core/go.mod is required' }
  if (-not (Test-Path -LiteralPath (Join-Path $RepositoryRoot 'apps/web/package.json'))) { throw 'apps/web/package.json is required' }
}

function Start-DevLocalPostgres([string]$RepositoryRoot) {
  Ensure-DevLocalDockerPath
  & docker compose up -d postgres
  if ($LASTEXITCODE -ne 0) { throw 'PostgreSQL Docker service did not start; inspect docker compose ps and PostgreSQL logs' }
}

function Wait-DevLocalPostgres([string]$RepositoryRoot) {
  Ensure-DevLocalDockerPath
  $deadline = (Get-Date).AddMinutes(3)
  do {
    $inspection = @(& docker inspect marketplace-central-postgres-1 --format '{{.State.Health.Status}}' 2>$null)
    $health = if ($inspection.Count -eq 0) { '' } else { [string]$inspection[$inspection.Count - 1] }
    $health = $health.Trim()
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
  New-Item -ItemType Directory -Force -Path $state | Out-Null
  $stdoutLog = Join-Path $state "$Name.stdout.log"
  $stderrLog = Join-Path $state "$Name.stderr.log"
  $existing = Read-DevLocalPid $state $Name
  if ($existing -and (Test-DevLocalPidOwnedByRuntime $state $Name $existing)) { throw "$Name is already running with PID $existing" }
  Remove-DevLocalPid $state $Name
  $startArgs = @{
    FilePath               = $FilePath
    ArgumentList           = $ArgumentList
    WorkingDirectory       = $RepositoryRoot
    RedirectStandardOutput = $stdoutLog
    RedirectStandardError  = $stderrLog
    PassThru               = $true
  }
  # -WindowStyle exists only on Windows PowerShell editions; linux pwsh rejects the parameter outright.
  if ($IsWindows) { $startArgs.WindowStyle = 'Hidden' }
  $process = Start-Process @startArgs
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
