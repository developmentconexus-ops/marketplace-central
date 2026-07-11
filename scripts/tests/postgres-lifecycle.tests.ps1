$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$module = Join-Path $repoRoot 'scripts/harness/Postgres.psm1'
$probe = Join-Path $repoRoot 'scripts/tests/fixtures/harness/postgres-docker-probe.mjs'

function Assert-True([bool]$Condition, [string]$Message) {
  if (-not $Condition) { throw $Message }
}

if (-not (Test-Path -LiteralPath $module -PathType Leaf)) {
  throw 'RED: Postgres.psm1 missing; cleanup lifecycle is not implemented'
}
Import-Module $module -Force

$node = [IO.Path]::GetFullPath((Get-Command node -CommandType Application -ErrorAction Stop).Source)
$base = [System.Collections.Generic.Dictionary[string,string]]::new([StringComparer]::OrdinalIgnoreCase)
foreach ($key in @('SystemRoot', 'WINDIR', 'ComSpec', 'PATH', 'PATHEXT', 'TEMP', 'TMP')) {
  $value = [Environment]::GetEnvironmentVariable($key, 'Process')
  if (-not [string]::IsNullOrWhiteSpace($value)) { $base[$key] = $value }
}

function Invoke-ProbeLifecycle([string]$FailureOperations) {
  $runId = [guid]::NewGuid().ToString('N')
  $password = 'fixture-' + [guid]::NewGuid().ToString('N')
  $log = Join-Path ([IO.Path]::GetTempPath()) ("postgres-lifecycle-$runId.jsonl")
  $environment = [System.Collections.Generic.Dictionary[string,string]]::new($base, [StringComparer]::OrdinalIgnoreCase)
  $environment['HARNESS_POSTGRES_PROBE_LOG'] = $log
  if (-not [string]::IsNullOrWhiteSpace($FailureOperations)) { $environment['HARNESS_POSTGRES_PROBE_FAIL_OPERATIONS'] = $FailureOperations }
  $spec = New-HarnessPostgresRunSpec -RepositoryRoot $repoRoot -RunId $runId -Password $password -DockerFilePath $node -DockerArgumentPrefix @($probe, 'docker')
  $result = Invoke-HarnessPostgresLifecycle -RunSpec $spec -BaseEnvironment $environment -GoFilePath $node -GoArgumentPrefix @($probe, 'go') -TimeoutSeconds 10
  $calls = if (Test-Path -LiteralPath $log) { @(Get-Content -LiteralPath $log | ForEach-Object { $_ | ConvertFrom-Json }) } else { @() }
  [pscustomobject]@{ Result = $result; Calls = $calls; Log = $log; Password = $password }
}

function Remove-ProbeRun([object]$Run) {
  Remove-Item -LiteralPath $Run.Log, "$($Run.Log).state.json" -Force -ErrorAction SilentlyContinue
}

$runs = @()
try {
  $happy = Invoke-ProbeLifecycle ''
  $runs += $happy
  $operations = @($happy.Calls.operation)
  $expected = @('start', 'ready', 'port', 'create', 'migrate', 'migrate', 'tests', 'drop', 'remove', 'inventory')
  Assert-True (($operations -join ',') -ceq ($expected -join ',')) "lifecycle order mismatch: $($operations -join ',')"
  Assert-True ($happy.Result.ExitCode -eq 0) 'happy lifecycle returned nonzero'
  Assert-True ($happy.Result.MigrationsAppliedFirst -eq 32) 'first migration count is not exact canonical inventory'
  Assert-True ($happy.Result.MigrationsAppliedSecond -eq 0) 'second migration run is not idempotent'
  Assert-True (@($happy.Result.ResourceInventory).Count -eq 0) 'happy lifecycle reported leaked resources'
  $drop = @($happy.Calls | Where-Object operation -eq 'drop')[0]
  Assert-True (($drop.args -join ' ') -match 'DROP DATABASE' -and ($drop.args -join ' ') -match 'WITH \(FORCE\)') 'cleanup does not force-drop active connections'

  $testsFail = Invoke-ProbeLifecycle 'tests'
  $runs += $testsFail
  Assert-True ($testsFail.Result.ExitCode -eq 17) 'child exit 17 was not preserved'
  Assert-True ($testsFail.Result.PrimaryReasonCode -eq 'HPG_TEST_FAILED') 'test failure reason changed'
  Assert-True (@($testsFail.Calls.operation) -contains 'drop' -and @($testsFail.Calls.operation) -contains 'remove') 'test failure skipped cleanup'

  foreach ($case in @(
    @{ Fail = 'ready'; Reason = 'HPG_READY_TIMEOUT' },
    @{ Fail = 'port'; Reason = 'HPG_PORT_UNAVAILABLE' },
    @{ Fail = 'create'; Reason = 'HPG_DATABASE_CREATE_FAILED' },
    @{ Fail = 'drop'; Reason = 'HPG_DATABASE_DROP_FAILED' },
    @{ Fail = 'remove'; Reason = 'HPG_CONTAINER_REMOVE_FAILED' }
  )) {
    $run = Invoke-ProbeLifecycle $case.Fail
    $runs += $run
    Assert-True ($run.Result.ExitCode -ne 0) "$($case.Fail) failure returned zero"
    Assert-True (($run.Result.PrimaryReasonCode -eq $case.Reason) -or (@($run.Result.CleanupReasonCodes) -contains $case.Reason)) "$($case.Fail) lacks $($case.Reason)"
    Assert-True (@($run.Calls.operation) -contains 'remove') "$($case.Fail) failure skipped container removal attempt"
    if ($case.Fail -eq 'remove') {
      Assert-True (@($run.Result.CleanupReasonCodes) -contains 'HPG_RESOURCE_LEAK') 'nonempty inventory lacks HPG_RESOURCE_LEAK'
      Assert-True (@($run.Result.ResourceInventory).Count -gt 0) 'remove failure hid leaked inventory'
    }
  }

  $combined = Invoke-ProbeLifecycle 'tests,drop'
  $runs += $combined
  Assert-True ($combined.Result.ExitCode -eq 17) 'cleanup failure masked primary child exit'
  Assert-True ($combined.Result.PrimaryReasonCode -eq 'HPG_TEST_FAILED') 'cleanup failure replaced primary reason'
  Assert-True (@($combined.Result.CleanupReasonCodes) -contains 'HPG_DATABASE_DROP_FAILED') 'cleanup failure was discarded'
  Assert-True (@($combined.Calls.operation) -contains 'remove') 'combined failure skipped remove'
} finally {
  foreach ($run in $runs) { Remove-ProbeRun $run }
}

Write-Output 'PASS postgres lifecycle tests target=fake'
