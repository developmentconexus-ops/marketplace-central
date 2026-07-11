Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
Import-Module (Join-Path $repoRoot 'scripts/harness/Environment.psm1') -Force
Import-Module (Join-Path $repoRoot 'scripts/harness/Postgres.psm1') -Force

function Resolve-Application([string]$Name) {
  $command = Get-Command $Name -CommandType Application -ErrorAction Stop | Select-Object -First 1
  return [IO.Path]::GetFullPath($command.Source)
}

$docker = Resolve-HarnessPostgresDockerApplication
$go = Resolve-Application 'go'
$node = Resolve-Application 'node'
$probe = Join-Path $repoRoot 'scripts/tests/fixtures/harness/postgres-docker-probe.mjs'
$environment = New-HarnessChildEnvironment -RepositoryRoot $repoRoot -LaneId 'integration'
$environment['HARNESS_REAL_GO'] = $go
$passwordBytes = [byte[]]::new(24)
[Security.Cryptography.RandomNumberGenerator]::Fill($passwordBytes)
$password = [Convert]::ToHexString($passwordBytes).ToLowerInvariant()
$runId = [guid]::NewGuid().ToString('N')
$spec = New-HarnessPostgresRunSpec -RepositoryRoot $repoRoot -RunId $runId -Password $password -DockerFilePath $docker
$result = Invoke-HarnessPostgresLifecycle -RunSpec $spec -BaseEnvironment $environment -GoFilePath $node -GoArgumentPrefix @($probe, 'real-go') -HoldConnectionDuringCleanupTest -TimeoutSeconds 1200

if ($result.ExitCode -ne 17 -or $result.PrimaryReasonCode -ne 'HPG_TEST_FAILED') { throw 'real forced child exit was not preserved' }
if (@($result.CleanupReasonCodes).Count -ne 0) { throw "real forced cleanup failed reasons=$($result.CleanupReasonCodes -join ',')" }
if (@($result.ResourceInventory).Count -ne 0) { throw 'real forced cleanup leaked labelled resources' }
if ($result.MigrationsAppliedFirst -ne 32 -or $result.MigrationsAppliedSecond -ne 0) { throw 'real forced cleanup migration counts changed' }
if ($result.HostPort -lt 1) { throw 'real forced cleanup lacked dynamic loopback port' }
if (-not $result.HeldConnectionConfirmed) { throw 'real forced cleanup did not confirm active held connection' }

Write-Output "PASS postgres Docker failure cleanup target=ephemeral-postgres run_id=$runId port=$($result.HostPort) resource_count=0 child_exit=17 held_connection=confirmed-and-force-dropped"
