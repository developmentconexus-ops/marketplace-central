$ErrorActionPreference = 'Stop'
$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
. (Join-Path $repoRoot 'scripts/lib/dev-local-runtime.ps1')

function Assert-Equal([object]$Actual, [object]$Expected, [string]$Name) {
  if ($Actual -ne $Expected) { throw "${Name}: expected [$Expected], got [$Actual]" }
}

Assert-Equal (ConvertTo-DevLocalEnvironmentValue '"quoted value"') 'quoted value' 'double quotes removed'
Assert-Equal (ConvertTo-DevLocalEnvironmentValue "'quoted value'") 'quoted value' 'single quotes removed'
Assert-Equal (ConvertTo-DevLocalEnvironmentValue 'abc=def') 'abc=def' 'equals preserved'

$dockerPathTestRoot = Join-Path ([System.IO.Path]::GetTempPath()) ('mpc-dev-local-docker-path-test-' + [guid]::NewGuid())
$originalProcessPath = $env:PATH
$originalProgramFiles = $env:ProgramFiles
$originalUserPath = [Environment]::GetEnvironmentVariable('PATH', 'User')
$originalMachinePath = [Environment]::GetEnvironmentVariable('PATH', 'Machine')
try {
  $dockerDesktopBin = Join-Path $dockerPathTestRoot 'Docker\Docker\resources\bin'
  New-Item -ItemType Directory -Path $dockerDesktopBin -Force | Out-Null
  $env:ProgramFiles = $dockerPathTestRoot
  $env:PATH = 'C:\runtime-test-path'
  $script:dockerDiscoveryChecks = 0
  function Get-Command {
    param([string]$Name)
    if ($Name -eq 'docker') {
      $script:dockerDiscoveryChecks++
      if (($env:PATH -split [IO.Path]::PathSeparator) -contains $dockerDesktopBin) { return [pscustomobject]@{ Name = 'docker' } }
      return $null
    }
    return Microsoft.PowerShell.Core\Get-Command @PSBoundParameters
  }

  Ensure-DevLocalDockerPath

  if (($env:PATH -split [IO.Path]::PathSeparator)[0] -ne $dockerDesktopBin) { throw 'Docker Desktop bin was not prepended to the process PATH' }
  Assert-Equal $script:dockerDiscoveryChecks 2 'Docker discovery is retried after PATH correction'
  Assert-Equal ([Environment]::GetEnvironmentVariable('PATH', 'User')) $originalUserPath 'user PATH was not changed'
  Assert-Equal ([Environment]::GetEnvironmentVariable('PATH', 'Machine')) $originalMachinePath 'machine PATH was not changed'
} finally {
  Remove-Item Function:Get-Command -ErrorAction SilentlyContinue
  $env:PATH = $originalProcessPath
  $env:ProgramFiles = $originalProgramFiles
  Remove-Item -LiteralPath $dockerPathTestRoot -Recurse -Force -ErrorAction SilentlyContinue
}

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

$childRoot = Join-Path ([System.IO.Path]::GetTempPath()) ('mpc-dev-local-child-test-' + [guid]::NewGuid())
$childState = Get-DevLocalStateDirectory $childRoot
try {
  New-Item -ItemType Directory -Path $childRoot | Out-Null
  Start-DevLocalChild -RepositoryRoot $childRoot -Name 'backend' -FilePath 'pwsh' -ArgumentList @('-NoProfile', '-Command', 'Start-Sleep -Seconds 30')
  if (-not (Test-Path -LiteralPath $childState -PathType Container)) { throw 'child startup did not create the state directory' }
} finally {
  $childPid = Read-DevLocalPid -StateDirectory $childState -Name 'backend'
  if ($childPid) { Stop-Process -Id $childPid -Force -ErrorAction SilentlyContinue }
  Remove-Item -LiteralPath $childRoot -Recurse -Force -ErrorAction SilentlyContinue
}

$staleRoot = Join-Path ([System.IO.Path]::GetTempPath()) ('mpc-dev-local-stale-test-' + [guid]::NewGuid())
$staleState = Get-DevLocalStateDirectory $staleRoot
try {
  New-Item -ItemType Directory -Path $staleState -Force | Out-Null
  $staleRecord = [pscustomobject]@{ process_id = $PID; started_at_utc = '2000-01-01T00:00:00.0000000Z' }
  $staleRecord | ConvertTo-Json -Compress | Set-Content -LiteralPath (Get-DevLocalPidPath $staleState 'backend') -NoNewline
  Start-DevLocalChild -RepositoryRoot $staleRoot -Name 'backend' -FilePath 'pwsh' -ArgumentList @('-NoProfile', '-Command', 'Start-Sleep -Seconds 30')
  $startedPid = Read-DevLocalPid -StateDirectory $staleState -Name 'backend'
  if (-not $startedPid -or $startedPid -eq $PID) { throw 'stale PID record blocked child startup' }
} finally {
  $childPid = Read-DevLocalPid -StateDirectory $staleState -Name 'backend'
  if ($childPid -and $childPid -ne $PID) { Stop-Process -Id $childPid -Force -ErrorAction SilentlyContinue }
  Remove-Item -LiteralPath $staleRoot -Recurse -Force -ErrorAction SilentlyContinue
}

$unknown = & pwsh -NoProfile -File (Join-Path $repoRoot 'scripts/dev-local.ps1') -Command unknown 2>&1
if ($LASTEXITCODE -eq 0 -or ($unknown -join "`n") -notmatch 'unsupported command') { throw 'unknown command must fail without starting a service' }

$script:dockerInvocationCount = 0
function docker {
  param([Parameter(ValueFromRemainingArguments = $true)][string[]]$Arguments)
  $script:dockerInvocationCount++
  if (($Arguments -join ' ') -ne "inspect marketplace-central-postgres-1 --format {{.State.Health.Status}}") { throw 'runtime status issued a mutating Docker command' }
  'healthy'
}

function Invoke-WebRequest {
  [CmdletBinding()]
  param([string]$Uri, [int]$TimeoutSec, [switch]$UseBasicParsing)
  [pscustomobject]@{ StatusCode = 200 }
}

$runtimeStatus = Get-DevLocalRuntimeStatus -RepositoryRoot $repoRoot
if ($runtimeStatus -notmatch 'backend_health=healthy' -or $runtimeStatus -notmatch 'frontend_health=healthy' -or $runtimeStatus -notmatch 'postgres_health=healthy') { throw 'runtime status did not report every health surface' }
if ($runtimeStatus -match 'MC_DATABASE_URL|postgres://') { throw 'runtime status leaked configuration' }
Assert-Equal $script:dockerInvocationCount 1 'runtime status uses one Docker inspection'

$script:waitClockReads = 0
function docker {
  param([Parameter(ValueFromRemainingArguments = $true)][string[]]$Arguments)
}

function Get-Date {
  [CmdletBinding()]
  param()
  $script:waitClockReads++
  if ($script:waitClockReads -eq 1) { return [datetime]'2026-07-10T00:00:00Z' }
  return [datetime]'2026-07-10T00:03:01Z'
}

function Start-Sleep {
  param([int]$Seconds)
}

$postgresTimeout = $null
try {
  Wait-DevLocalPostgres -RepositoryRoot $repoRoot
} catch {
  $postgresTimeout = $_.Exception.Message
}
if ($postgresTimeout -notmatch 'PostgreSQL did not become healthy in 3 minutes') { throw 'empty Docker inspection did not reach the PostgreSQL timeout' }

Write-Output 'PASS dev-local-runtime tests'
