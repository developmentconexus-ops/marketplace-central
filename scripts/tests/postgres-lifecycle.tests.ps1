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

$node = [IO.Path]::GetFullPath((Get-Command node -CommandType Application -ErrorAction Stop | Select-Object -First 1).Source)
$base = [System.Collections.Generic.Dictionary[string,string]]::new([StringComparer]::OrdinalIgnoreCase)
foreach ($key in @('SystemRoot', 'WINDIR', 'ComSpec', 'PATH', 'PATHEXT', 'TEMP', 'TMP')) {
  $value = [Environment]::GetEnvironmentVariable($key, 'Process')
  if (-not [string]::IsNullOrWhiteSpace($value)) { $base[$key] = $value }
}

function Invoke-ProbeLifecycle(
  [string]$FailureOperations,
  [hashtable]$ProbeOptions = @{},
  [int]$ReadyMaxAttempts = 3,
  [int]$ReadyTimeoutMilliseconds = 10000,
  [bool]$HoldConnection = $false,
  [string[]]$TestArguments = @()
) {
  $runId = [guid]::NewGuid().ToString('N')
  $password = 'fixture-' + [guid]::NewGuid().ToString('N')
  $log = Join-Path ([IO.Path]::GetTempPath()) ("postgres-lifecycle-$runId.jsonl")
  $environment = [System.Collections.Generic.Dictionary[string,string]]::new($base, [StringComparer]::OrdinalIgnoreCase)
  $environment['HARNESS_POSTGRES_PROBE_LOG'] = $log
  if (-not [string]::IsNullOrWhiteSpace($FailureOperations)) { $environment['HARNESS_POSTGRES_PROBE_FAIL_OPERATIONS'] = $FailureOperations }
  foreach ($key in $ProbeOptions.Keys) { $environment[$key] = [string]$ProbeOptions[$key] }
  $spec = New-HarnessPostgresRunSpec -RepositoryRoot $repoRoot -RunId $runId -Password $password -DockerFilePath $node -DockerArgumentPrefix @($probe, 'docker')
  # 60s and not 10s: this budget bounds each fake-docker subprocess, and every
  # one of them is a node start. Under whole-suite load a node start has been
  # measured past the old 10s, and the harness then reported the reason for the
  # step it was on -- HPG_DOCKER_UNAVAILABLE where the case under test expected
  # HPG_CONTAINER_START_FAILED. The bounds these cases actually assert are their
  # own (readiness deadline, subprocess deadline); this one only has to be large
  # enough that a slow machine is not mistaken for a broken harness.
  $result = Invoke-HarnessPostgresLifecycle -RunSpec $spec -BaseEnvironment $environment -GoFilePath $node -GoArgumentPrefix @($probe, 'go') -TimeoutSeconds 60 -ReadyMaxAttempts $ReadyMaxAttempts -ReadyRetryDelayMilliseconds 0 -ReadyTimeoutMilliseconds $ReadyTimeoutMilliseconds -CreateMaxAttempts 3 -CreateRetryDelayMilliseconds 0 -HoldConnectionDuringCleanupTest:$HoldConnection -TestArguments $TestArguments
  $calls = if (Test-Path -LiteralPath $log) { @(Get-Content -LiteralPath $log | ForEach-Object { $_ | ConvertFrom-Json }) } else { @() }
  [pscustomobject]@{ Result = $result; Calls = $calls; Log = $log; Password = $password; ExpectedMigrationCount = $spec.ExpectedMigrationCount }
}

function Remove-ProbeRun([object]$Run) {
  Remove-Item -LiteralPath $Run.Log, "$($Run.Log).state.json" -Force -ErrorAction SilentlyContinue
}

$runs = @()
try {
  foreach ($inventoryState in @('missing', 'empty')) {
    $invalidRoot = Join-Path ([IO.Path]::GetTempPath()) ("postgres-inventory-$inventoryState-$([guid]::NewGuid().ToString('N'))")
    try {
      New-Item -ItemType Directory -Path $invalidRoot -Force | Out-Null
      if ($inventoryState -eq 'empty') { New-Item -ItemType Directory -Path (Join-Path $invalidRoot 'apps/server_core/migrations') -Force | Out-Null }
      $threw = $false
      try {
        New-HarnessPostgresRunSpec -RepositoryRoot $invalidRoot -RunId ([guid]::NewGuid().ToString('N')) -Password 'fixture-invalid-inventory' -DockerFilePath $node | Out-Null
      } catch {
        $threw = $true
        Assert-True ($_.Exception.Message -eq 'HPG_MIGRATION_INVENTORY_INVALID') "$inventoryState inventory lacks stable failure reason"
      }
      Assert-True $threw "$inventoryState migration inventory did not fail closed"
    } finally {
      Remove-Item -LiteralPath $invalidRoot -Recurse -Force -ErrorAction SilentlyContinue
    }
  }

  $happy = Invoke-ProbeLifecycle ''
  $runs += $happy
  $operations = @($happy.Calls.operation)
  $expected = @('daemon', 'image', 'preflight-name', 'preflight-label', 'start', 'ownership', 'ready', 'port', 'create-verify', 'create', 'migrate', 'migrate', 'tests', 'drop', 'drop-verify', 'remove', 'inventory-label', 'inventory-name')
  Assert-True (($operations -join ',') -ceq ($expected -join ',')) "lifecycle order mismatch: $($operations -join ',')"
  Assert-True ($happy.Result.ExitCode -eq 0) 'happy lifecycle returned nonzero'
  Assert-True ($happy.ExpectedMigrationCount -gt 0) 'canonical migration inventory is empty'
  Assert-True ($happy.Result.MigrationsAppliedFirst -eq $happy.ExpectedMigrationCount) 'first migration count is not exact canonical inventory'
  Assert-True ($happy.Result.MigrationsAppliedSecond -eq 0) 'second migration run is not idempotent'
  Assert-True (@($happy.Result.ResourceInventory).Count -eq 0) 'happy lifecycle reported leaked resources'
  Assert-True ($happy.Result.TestsRun -eq 2 -and $happy.Result.TestsPassed -eq 2 -and $happy.Result.TestsSkipped -eq 0 -and $happy.Result.TestsFailed -eq 0) "happy lifecycle did not count the tests it ran: run=$($happy.Result.TestsRun) passed=$($happy.Result.TestsPassed) skipped=$($happy.Result.TestsSkipped) failed=$($happy.Result.TestsFailed)"
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

  # The attempt count comes from the harness, not from the fixture's log. Each
  # attempt's subprocess budget is derived from the deadline that is left and
  # floored at one second, so with a 1ms deadline the fixture has one second to
  # start node and append its row -- and a child killed at that budget appends
  # nothing. Counting child rows therefore read 0 attempts on a run that made
  # one, which is how this assertion failed under whole-suite load on a tree
  # that behaved correctly. `Result.ReadyAttempts` is incremented by the loop
  # itself and cannot be lost that way. The log is still checked below for what
  # only the child can attest: the arguments it was called with.
  $readyDeadline = Invoke-ProbeLifecycle '' @{ HARNESS_POSTGRES_PROBE_READY_FAILURES = '99' } 99 1
  $runs += $readyDeadline
  $deadlineAttempts = [int]$readyDeadline.Result.ReadyAttempts
  Assert-True ($readyDeadline.Result.PrimaryReasonCode -eq 'HPG_READY_TIMEOUT') 'readiness deadline lacks stable reason'
  Assert-True ($deadlineAttempts -ge 1 -and $deadlineAttempts -lt 99) "readiness ignored its deadline bound: attempts=$deadlineAttempts"

  $hangingReadyWatch = [Diagnostics.Stopwatch]::StartNew()
  # 1500ms, because the assertion below reads the arguments the child recorded:
  # the child appends its row before it hangs, but it still has to start, and a
  # 1ms deadline gives it only the one-second floor to do that under load.
  $hangingReady = Invoke-ProbeLifecycle '' @{ HARNESS_POSTGRES_PROBE_READY_HANG_MILLISECONDS = '60000' } 99 1500
  $hangingReadyWatch.Stop()
  $runs += $hangingReady
  Assert-True ($hangingReady.Result.PrimaryReasonCode -eq 'HPG_READY_TIMEOUT') 'hanging readiness lacks stable timeout reason'
  # The fixture hangs pg_isready for 60s; an unbounded harness therefore takes
  # >= 60s for even one ready attempt. The gap between the bound and the hang is
  # what keeps this measuring the harness instead of the machine: at a 10s hang
  # and an 8s bound, a loaded machine spending a few seconds on node starts and
  # teardown failed this assertion on a harness that bounded the call correctly.
  Assert-True ($hangingReadyWatch.Elapsed.TotalSeconds -lt 30) 'hanging pg_isready exceeded bounded subprocess deadline'
  $hangingCalls = @($hangingReady.Calls | Where-Object operation -eq 'ready')
  # Separated from the argument check on purpose: an empty log and a call made
  # without --timeout used to fail with the same message, and only one of them
  # is the harness losing its bound.
  Assert-True ($hangingCalls.Count -ge 1) "hanging readiness recorded no call to inspect: attempts=$($hangingReady.Result.ReadyAttempts)"
  Assert-True (@($hangingCalls[0].args) -contains '--timeout') 'pg_isready lacks its own bounded timeout'

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

  # `go test` exits 0 over a package set with no tests in it. Every other step of
  # this lifecycle succeeds here -- container, migrations, cleanup -- so before
  # the count existed this run was indistinguishable from the happy one above.
  $vacuous = Invoke-ProbeLifecycle '' @{ HARNESS_POSTGRES_PROBE_TEST_COUNT = '0' }
  $runs += $vacuous
  Assert-True ($vacuous.Result.ExitCode -ne 0) 'a run that executed zero tests reported success'
  Assert-True ($vacuous.Result.PrimaryReasonCode -eq 'HPG_TEST_VACUOUS') "zero-test run lacks stable reason: $($vacuous.Result.PrimaryReasonCode)"
  Assert-True (@($vacuous.Result.FailureDiagnosticTokens) -contains 'tests_run=0') 'zero-test run is not attributable from its tokens'
  Assert-True (@($vacuous.Calls.operation) -contains 'drop' -and @($vacuous.Calls.operation) -contains 'remove') 'vacuous run skipped cleanup'
  Write-Output 'guard_ran=integration-vacuous-run'

  # The profile's own signature: every test ran and every test skipped, package
  # exits 0. Counting `=== RUN` alone accepts this, which is why the guard reads
  # `passed -eq 0` as well.
  $allSkipped = Invoke-ProbeLifecycle '' @{ HARNESS_POSTGRES_PROBE_TEST_COUNT = '3'; HARNESS_POSTGRES_PROBE_TEST_SKIPPED = '3' }
  $runs += $allSkipped
  Assert-True ($allSkipped.Result.TestsRun -eq 3 -and $allSkipped.Result.TestsSkipped -eq 3 -and $allSkipped.Result.TestsPassed -eq 0) "skips were not counted: run=$($allSkipped.Result.TestsRun) passed=$($allSkipped.Result.TestsPassed) skipped=$($allSkipped.Result.TestsSkipped)"
  Assert-True ($allSkipped.Result.ExitCode -ne 0) 'a fully skipped suite reported success'
  Assert-True ($allSkipped.Result.PrimaryReasonCode -eq 'HPG_TEST_VACUOUS') "fully skipped run lacks stable reason: $($allSkipped.Result.PrimaryReasonCode)"

  # A partially skipped suite is not vacuous -- something asserted.
  $partiallySkipped = Invoke-ProbeLifecycle '' @{ HARNESS_POSTGRES_PROBE_TEST_COUNT = '3'; HARNESS_POSTGRES_PROBE_TEST_SKIPPED = '2' }
  $runs += $partiallySkipped
  Assert-True ($partiallySkipped.Result.ExitCode -eq 0) 'a suite with one real pass was rejected as vacuous'
  Assert-True ($partiallySkipped.Result.TestsSkipped -eq 2 -and $partiallySkipped.Result.TestsPassed -eq 1) "partial skip counts wrong: passed=$($partiallySkipped.Result.TestsPassed) skipped=$($partiallySkipped.Result.TestsSkipped)"

  # Caller-supplied arguments take the same -v, or every custom-argument run
  # would report zero tests and be rejected as vacuous.
  $customArguments = Invoke-ProbeLifecycle '' @{} 3 10000 $false @('test', '-tags=integration', './tests/integration', '-count=1')
  $runs += $customArguments
  Assert-True ($customArguments.Result.ExitCode -eq 0) 'custom test arguments were rejected as vacuous'
  Assert-True ($customArguments.Result.TestsRun -eq 2) "custom test arguments produced no countable output: run=$($customArguments.Result.TestsRun)"
  $customCall = @($customArguments.Calls | Where-Object operation -eq 'tests')[0]
  Assert-True (@($customCall.args) -contains '-v') 'custom test arguments were not forced verbose'

  $testsFail = Invoke-ProbeLifecycle 'tests'
  $runs += $testsFail
  Assert-True ($testsFail.Result.ExitCode -eq 17) 'child exit 17 was not preserved'
  Assert-True ($testsFail.Result.PrimaryReasonCode -eq 'HPG_TEST_FAILED') 'test failure reason changed'
  $failureTokens = @($testsFail.Result.FailureDiagnosticTokens)
  Assert-True ($failureTokens -contains 'hpg=HPG_TEST_FAILED_SENTINEL') 'safe child failure token was hidden'
  Assert-True (($failureTokens -join "`n").Length -le 8192) 'child failure tokens exceeded bound'
  Assert-True (($failureTokens -join "`n") -notmatch '(?i)(?:private|person@|customer|[a-z]:\\)') 'arbitrary child output entered structured diagnostics'
  Assert-True (($failureTokens -join "`n") -notmatch [regex]::Escape($testsFail.Password)) 'password leaked in child failure tokens'
  Assert-True (($failureTokens -join "`n") -notmatch [regex]::Escape($repoRoot)) 'repository root leaked in child failure tokens'
  Assert-True (($failureTokens -join "`n") -notmatch 'postgres(?:ql)?://[^\s]+') 'database target leaked in child failure tokens'
  Assert-True (@($testsFail.Calls.operation) -contains 'drop' -and @($testsFail.Calls.operation) -contains 'remove') 'test failure skipped cleanup'

  # A failing lane has to say WHERE, or the only move left is to re-run it until it
  # goes green. These three tokens are the difference between "TestX failed" and a
  # diagnosis; they are source and schema identifiers, so the redaction guarantees
  # asserted for every other failure still have to hold for them.
  $richFailure = Invoke-ProbeLifecycle 'tests' @{ HARNESS_POSTGRES_PROBE_RICH_TEST_FAILURE = '1' }
  $runs += $richFailure
  $richTokens = @($richFailure.Result.FailureDiagnosticTokens)
  Assert-True ($richTokens -contains 'subtest=TestProbeContract/seeded_rows_survive_the_walk') "failing subtest was not named: $($richTokens -join ',')"
  Assert-True ($richTokens -contains 'at=probe_contract_test.go:290') "failing assertion site was not named: $($richTokens -join ',')"
  Assert-True ($richTokens -contains 'constraint=listings_pkey') "violated constraint was not named: $($richTokens -join ',')"
  Assert-True ($richTokens -contains 'test=TestProbeContract' -and $richTokens -contains 'sqlstate=23505') 'the pre-existing token shapes stopped being emitted'
  # A must-fail test logs the constraint it provoked and then PASSES. Swept
  # rather than attributed, its file:line and its constraint are emitted beside
  # the real failure and read exactly like it.
  Assert-True ($richTokens -notcontains 'at=probe_mustfail_test.go:41') "a passing test's assertion site was reported as a failure site: $($richTokens -join ',')"
  Assert-True ($richTokens -notcontains 'constraint=probe_only_one_open_row') "a passing test's constraint was reported as violated: $($richTokens -join ',')"
  # Interleaved parallel output: the passing test's line prints between the
  # failing test's line and the failing verdict. Position says both belong to
  # the failure; identity says only one does.
  Assert-True ($richTokens -contains 'at=probe_parallel_test.go:44') "the failing parallel test's site was lost to interleaving: $($richTokens -join ',')"
  Assert-True ($richTokens -notcontains 'at=probe_parallel_test.go:77') "an interleaved passing test's site was attributed to the failure: $($richTokens -join ',')"
  Assert-True ($richTokens -notcontains 'constraint=probe_tolerated_key') "an interleaved passing test's constraint was attributed to the failure: $($richTokens -join ',')"

  # The count cap alone does not bound the output: 32 maximum-shape subtest
  # tokens are roughly 28 KiB, four times the contract asserted below.
  $longFailure = Invoke-ProbeLifecycle 'tests' @{ HARNESS_POSTGRES_PROBE_LONG_TOKEN_FAILURE = '1' }
  $runs += $longFailure
  $longTokens = @($longFailure.Result.FailureDiagnosticTokens)
  Assert-True ($longTokens.Count -ge 1) 'maximum-shape tokens produced no diagnostics at all'
  Assert-True ((($longTokens -join ' ').Length) -le 8192) "aggregate diagnostics exceeded 8192 bytes: $(($longTokens -join ' ').Length)"
  Assert-True (@($longTokens | Where-Object { $_ -like 'subtest=*' -and $_.Length -lt 200 }).Count -eq 0) 'a long token was truncated instead of dropped'
  Assert-True (($richTokens -join "`n") -notmatch '(?i)(?:private|person@|customer|[a-z]:\\)') 'arbitrary child output entered structured diagnostics through the new tokens'
  Assert-True (($richTokens -join "`n") -notmatch [regex]::Escape($richFailure.Password)) 'password leaked through the new tokens'
  Assert-True (($richTokens -join "`n") -notmatch [regex]::Escape($repoRoot)) 'repository root leaked through the new tokens'
  Assert-True (($richTokens -join "`n") -notmatch 'postgres(?:ql)?://[^\s]+') 'database target leaked through the new tokens'

  $arbitraryFailure = Invoke-ProbeLifecycle 'tests' @{ HARNESS_POSTGRES_PROBE_ARBITRARY_TEST_FAILURE = '1' }
  $runs += $arbitraryFailure
  Assert-True ((@($arbitraryFailure.Result.FailureDiagnosticTokens) -join ',') -ceq 'HPG_CHILD_OUTPUT_REDACTED') 'arbitrary child output did not collapse to safe fallback'

  $heldConnection = Invoke-ProbeLifecycle 'tests' @{} 3 10000 $true
  $runs += $heldConnection
  Assert-True ($heldConnection.Result.HeldConnectionConfirmed) 'held connection was not confirmed through pg_stat_activity'
  Assert-True (@($heldConnection.Calls.operation) -contains 'held-check') 'held connection confirmation query did not run'

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
