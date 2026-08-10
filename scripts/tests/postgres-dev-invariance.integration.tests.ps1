Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
Import-Module (Join-Path $repoRoot 'scripts/harness/Environment.psm1') -Force
Import-Module (Join-Path $repoRoot 'scripts/harness/Execution.psm1') -Force

function Resolve-Application([string]$Name) {
  $command = Get-Command $Name -CommandType Application -ErrorAction Stop | Select-Object -First 1
  return [IO.Path]::GetFullPath($command.Source)
}

$docker = Resolve-Application 'docker'
$dockerCompose = Resolve-Application 'docker-compose'
$pwsh = Resolve-Application 'pwsh'
$environment = New-HarnessChildEnvironment -RepositoryRoot $repoRoot -LaneId 'dev-invariance'

function Invoke-Typed([string]$FilePath, [string[]]$Arguments, [int]$TimeoutSeconds = 1200) {
  $request = New-HarnessProcessRequest -FilePath $FilePath -ArgumentList $Arguments -WorkingDirectory $repoRoot -Environment $environment -TimeoutSeconds $TimeoutSeconds
  return Invoke-HarnessProcess -Request $request
}

function Get-DevDigest {
  $containerResult = Invoke-Typed $dockerCompose @('ps', '--status', 'running', '--quiet', 'postgres') 30
  if ($containerResult.ExitCode -ne 0) { throw 'dev observer unavailable: compose postgres lookup failed' }
  $containerID = $containerResult.Stdout.Trim()
  if ($containerID -notmatch '^[0-9a-f]{12,64}$') { throw 'dev observer unavailable: compose postgres is not uniquely running' }

  $ownership = Invoke-Typed $docker @('inspect', '--format', '{{ index .Config.Labels "com.docker.compose.project" }}|{{ index .Config.Labels "com.docker.compose.service" }}', $containerID) 30
  if ($ownership.ExitCode -ne 0 -or $ownership.Stdout.Trim() -cne 'marketplace-central|postgres') { throw 'dev observer target ownership mismatch' }

  function Invoke-ReadOnlyPsql([string]$Query) {
    Invoke-Typed $docker @(
      'exec', '--env', 'PGOPTIONS=-cdefault_transaction_read_only=on', $containerID,
      'psql', '--username', 'marketplace', '--dbname', 'marketplace_central',
      '--quiet', '--tuples-only', '--no-align', '--set', 'ON_ERROR_STOP=1', '--command', $Query
    ) 120
  }

  $readOnly = Invoke-ReadOnlyPsql 'SHOW default_transaction_read_only'
  if ($readOnly.ExitCode -ne 0 -or $readOnly.Stdout.Trim() -cne 'on') { throw 'dev observer unavailable: session is not read-only' }
  $tables = Invoke-ReadOnlyPsql "SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE' ORDER BY table_name"
  if ($tables.ExitCode -ne 0) { throw 'dev observer unavailable: table inventory query failed' }
  $tableNames = @($tables.Stdout -split '\r?\n' | Where-Object { -not [string]::IsNullOrWhiteSpace($_) })
  $lines = [System.Collections.Generic.List[string]]::new()
  foreach ($tableName in $tableNames) {
    if ($tableName -cnotmatch '^[a-z_][a-z0-9_]*$') { throw 'dev observer unavailable: unsafe table identifier' }
    $count = Invoke-ReadOnlyPsql "SELECT count(*) FROM public.`"$tableName`""
    if ($count.ExitCode -ne 0 -or $count.Stdout.Trim() -notmatch '^\d+$') { throw 'dev observer unavailable: table count query failed' }
    $lines.Add("$tableName|$($count.Stdout.Trim())")
  }
  if ($lines.Count -eq 0) { throw 'dev observer unavailable: no public tables found' }
  $bytes = [Text.Encoding]::UTF8.GetBytes(($lines -join "`n"))
  $digest = [Convert]::ToHexString([Security.Cryptography.SHA256]::HashData($bytes)).ToLowerInvariant()
  return [pscustomobject]@{ Digest = $digest; TableCount = $lines.Count }
}

$before = Get-DevDigest
$nested = Invoke-Typed $pwsh @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', (Join-Path $repoRoot 'scripts/harness.ps1'), '-Command', 'integration') 1800
$after = Get-DevDigest
if ($before.Digest -cne $after.Digest -or $before.TableCount -ne $after.TableCount) { throw 'dev-state-change: persistent dev table-count digest changed' }
if ($nested.ExitCode -ne 0) {
  $safeDiagnostics = @((@($nested.Stdout, $nested.Stderr) -join "`n") -split '\r?\n' | Where-Object {
    $_ -match '^(?:status=(?:blocked|failed)|child_diagnostic=(?:hpg=[A-Z0-9_]+|sqlstate=[0-9A-Z]{5}|package=marketplace-central/[A-Za-z0-9_./-]+|test=Test[A-Za-z0-9_]+|subtest=Test[A-Za-z0-9_/]+|at=[A-Za-z0-9_]+_test\.go:\d+|constraint=[a-z0-9_]+|HPG_CHILD_OUTPUT_REDACTED)|migrations_(?:first|second)=-?\d+|resource_count=\d+)$'
  })
  foreach ($line in $safeDiagnostics) { Write-Output "nested_diagnostic=$line" }
  throw 'nested ephemeral PostgreSQL integration failed'
}

Write-Output "PASS postgres dev invariance target=dev-invariance nested_target=ephemeral-postgres digest=$($before.Digest) table_count=$($before.TableCount)"
