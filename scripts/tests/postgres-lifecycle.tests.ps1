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

function Invoke-ProbeLifecycle(
  [string]$FailureOperations,
  [hashtable]$ProbeOptions = @{},
  [int]$ReadyMaxAttempts = 3,
  [int]$ReadyTimeoutMilliseconds = 10000
) {
  $runId = [guid]::NewGuid().ToString('N')
  $password = 'fixture-' + [guid]::NewGuid().ToString('N')
  $log = Join-Path ([IO.Path]::GetTempPath()) ("postgres-lifecycle-$runId.jsonl")
  $environment = [System.Collections.Generic.Dictionary[string,string]]::new($base, [StringComparer]::OrdinalIgnoreCase)
  $environment['HARNESS_POSTGRES_PROBE_LOG'] = $log
  if (-not [string]::IsNullOrWhiteSpace($FailureOperations)) { $environment['HARNESS_POSTGRES_PROBE_FAIL_OPERATIONS'] = $FailureOperations }
  foreach ($key in $ProbeOptions.Keys) { $environment[$key] = [string]$ProbeOptions[$key] }
  $spec = New-HarnessPostgresRunSpec -RepositoryRoot $repoRoot -RunId $runId -Password $password -DockerFilePath $node -DockerArgumentPrefix @($probe, 'docker')
  $result = Invoke-HarnessPostgresLifecycle -RunSpec $spec -BaseEnvironment $environment -GoFilePath $node -GoArgumentPrefix @($probe, 'go') -TimeoutSeconds 10 -ReadyMaxAttempts $ReadyMaxAttempts -ReadyRetryDelayMilliseconds 0 -ReadyTimeoutMilliseconds $ReadyTimeoutMilliseconds
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
  $expected = @('daemon', 'image', 'preflight-name', 'preflight-label', 'start', 'ownership', 'ready', 'port', 'create', 'migrate', 'migrate', 'tests', 'drop', 'drop-verify', 'remove', 'inventory-label', 'inventory-name')
  Assert-True (($operations -join ',') -ceq ($expected -join ',')) "lifecycle order mismatch: $($operations -join ',')"
  Assert-True ($happy.Result.ExitCode -eq 0) 'happy lifecycle returned nonzero'
  Assert-True ($happy.Result.MigrationsAppliedFirst -eq 32) 'first migration count is not exact canonical inventory'
  Assert-True ($happy.Result.MigrationsAppliedSecond -eq 0) 'second migration run is not idempotent'
  Assert-True (@($happy.Result.ResourceInventory).Count -eq 0) 'happy lifecycle reported leaked resources'
  $drop = @($happy.Calls | Where-Object operation -eq 'drop')[0]
  Assert-True (($drop.args -join ' ') -match 'DROP DATABASE' -and ($drop.args -join ' ') -match 'WITH \(FORCE\)') 'cleanup does not force-drop active connections'

  $readyEventually = Invoke-ProbeLifecycle '' @{ HARNESS_POSTGRES_PROBE_READY_FAILURES = '2' }
  $runs += $readyEventually
  Assert-True ($readyEventually.Result.ExitCode -eq 0) 'bounded readiness retry did not accept eventual success'
  Assert-True (@($readyEventually.Calls | Where-Object operation -eq 'ready').Count -eq 3) 'readiness retry count changed'

  $readyTimeout = Invoke-ProbeLifecycle '' @{ HARNESS_POSTGRES_PROBE_READY_FAILURES = '99' }
  $runs += $readyTimeout
  Assert-True ($readyTimeout.Result.PrimaryReasonCode -eq 'HPG_READY_TIMEOUT') 'readiness exhaustion lacks stable reason'
  Assert-True (@($readyTimeout.Calls | Where-Object operation -eq 'ready').Count -eq 3) 'readiness exceeded bounded attempts'

  $readyDeadline = Invoke-ProbeLifecycle '' @{ HARNESS_POSTGRES_PROBE_READY_FAILURES = '99' } 99 1
  $runs += $readyDeadline
  $deadlineAttempts = @($readyDeadline.Calls | Where-Object operation -eq 'ready').Count
  Assert-True ($readyDeadline.Result.PrimaryReasonCode -eq 'HPG_READY_TIMEOUT') 'readiness deadline lacks stable reason'
  Assert-True ($deadlineAttempts -ge 1 -and $deadlineAttempts -lt 99) 'readiness ignored its deadline bound'

  $hangingReadyWatch = [Diagnostics.Stopwatch]::StartNew()
  $hangingReady = Invoke-ProbeLifecycle '' @{ HARNESS_POSTGRES_PROBE_READY_HANG_MILLISECONDS = '10000' } 99 1
  $hangingReadyWatch.Stop()
  $runs += $hangingReady
  Assert-True ($hangingReady.Result.PrimaryReasonCode -eq 'HPG_READY_TIMEOUT') 'hanging readiness lacks stable timeout reason'
  Assert-True ($hangingReadyWatch.Elapsed.TotalSeconds -lt 4) 'hanging pg_isready exceeded bounded subprocess deadline'
  $hangingCall = @($hangingReady.Calls | Where-Object operation -eq 'ready')[0]
  Assert-True (@($hangingCall.args) -contains '--timeout') 'pg_isready lacks its own bounded timeout'

  foreach ($conflict in @('HARNESS_POSTGRES_PROBE_NAME_CONFLICT', 'HARNESS_POSTGRES_PROBE_LABEL_CONFLICT')) {
    $run = Invoke-ProbeLifecycle '' @{ $conflict = '1' }
    $runs += $run
    Assert-True ($run.Result.PrimaryReasonCode -eq 'HPG_RESOURCE_CONFLICT') "$conflict lacks stable conflict reason"
    Assert-True (-not (@($run.Calls.operation) -contains 'start')) "$conflict started a container"
    Assert-True (-not (@($run.Calls.operation) -contains 'remove')) "$conflict removed a resource it did not own"
  }

  foreach ($preflight in @(
    @{ Fail = 'daemon'; Reason = 'HPG_DOCKER_UNAVAILABLE' },
    @{ Fail = 'image'; Reason = 'HPG_IMAGE_MISSING' }
  )) {
    $run = Invoke-ProbeLifecycle $preflight.Fail
    $runs += $run
    Assert-True ($run.Result.PrimaryReasonCode -eq $preflight.Reason) "$($preflight.Fail) lacks $($preflight.Reason)"
    Assert-True (-not (@($run.Calls.operation) -contains 'start')) "$($preflight.Fail) failure started a container"
    Assert-True (-not (@($run.Calls.operation) -contains 'remove')) "$($preflight.Fail) failure removed an unowned resource"
  }

  $unconfirmed = Invoke-ProbeLifecycle 'ownership'
  $runs += $unconfirmed
  Assert-True ($unconfirmed.Result.PrimaryReasonCode -eq 'HPG_CONTAINER_START_FAILED') 'ownership confirmation failure lacks stable start reason'
  Assert-True (-not (@($unconfirmed.Calls.operation) -contains 'remove')) 'unconfirmed container ownership allowed removal'

  $startUnconfirmed = Invoke-ProbeLifecycle 'start'
  $runs += $startUnconfirmed
  Assert-True ($startUnconfirmed.Result.PrimaryReasonCode -eq 'HPG_CONTAINER_START_FAILED') 'start failure lacks stable reason'
  Assert-True (-not (@($startUnconfirmed.Calls.operation) -contains 'remove')) 'failed start allowed removal without ownership confirmation'

  $dropRemains = Invoke-ProbeLifecycle '' @{ HARNESS_POSTGRES_PROBE_DROP_REMAINS = '1' }
  $runs += $dropRemains
  Assert-True (@($dropRemains.Result.CleanupReasonCodes) -contains 'HPG_DATABASE_DROP_FAILED') 'remaining database after successful DROP was not detected'
  Assert-True (@($dropRemains.Calls.operation) -contains 'drop-verify') 'database absence was not verified through maintenance database'

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
    @{ Fail = 'drop-verify'; Reason = 'HPG_DATABASE_DROP_FAILED' },
    @{ Fail = 'remove'; Reason = 'HPG_CONTAINER_REMOVE_FAILED' }
  )) {
    $run = Invoke-ProbeLifecycle $case.Fail
    $runs += $run
    Assert-True ($run.Result.ExitCode -ne 0) "$($case.Fail) failure returned zero"
    Assert-True (($run.Result.PrimaryReasonCode -eq $case.Reason) -or (@($run.Result.CleanupReasonCodes) -contains $case.Reason)) "$($case.Fail) lacks $($case.Reason)"
    Assert-True (@($run.Calls.operation) -contains 'remove') "$($case.Fail) failure skipped owned container removal attempt"
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
