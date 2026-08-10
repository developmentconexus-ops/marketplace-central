$ErrorActionPreference = 'Stop'

# Measure-GateGuards: the counter the `guards` lane (issue #3) uses to hold a
# targeted `go test -run` output against the inventory's expected test names.

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$module = Join-Path $repoRoot 'scripts/harness/Gate.psm1'
$gateScript = Join-Path $repoRoot 'scripts/gate.ps1'

function Assert-True([bool]$Condition, [string]$Message) {
  if (-not $Condition) { throw $Message }
}

if (-not (Test-Path -LiteralPath $module -PathType Leaf)) {
  throw 'RED: Gate.psm1 missing; the gate has no countable core'
}
if (-not (Test-Path -LiteralPath $gateScript -PathType Leaf)) {
  throw 'RED: scripts/gate.ps1 missing; the gate has no entry point'
}
Import-Module $module -Force

$happy = @"
=== RUN   TestAlpha
--- PASS: TestAlpha (0.00s)
=== RUN   TestBeta
--- PASS: TestBeta (0.01s)
PASS
ok      marketplace-central/apps/server_core/internal/arch     0.117s
"@
$m = Measure-GateGuards -Text $happy -Expected @('TestAlpha', 'TestBeta')
Assert-True ($m.Missing.Count -eq 0) "happy path reported missing: $($m.Missing -join ',')"
Assert-True ($m.Ran -eq 2) "happy path ran=$($m.Ran), want 2"

$m = Measure-GateGuards -Text $happy -Expected @('TestAlpha', 'TestBeta', 'TestGamma')
Assert-True ($m.Missing -contains 'TestGamma') 'a renamed/deleted inventory test must surface as missing'

$vacuous = @"
testing: warning: no tests to run
PASS
ok      marketplace-central/apps/server_core/internal/arch     0.021s [no tests to run]
"@
$m = Measure-GateGuards -Text $vacuous -Expected @('TestAlpha')
Assert-True ($m.Ran -eq 0) 'a no-tests-to-run output must count zero RUN lines'
Assert-True ($m.Missing -contains 'TestAlpha') 'a vacuous run satisfies no expectation'

$prefix = @"
=== RUN   TestAlphaBravo
--- PASS: TestAlphaBravo (0.00s)
PASS
"@
$m = Measure-GateGuards -Text $prefix -Expected @('TestAlpha')
Assert-True ($m.Missing -contains 'TestAlpha') 'TestAlphaBravo must not satisfy TestAlpha (word boundary)'


# --- Test-GateGuardInventory: the verdict logic, extracted from gate.ps1 ---
#
# Measure-GateGuards above proves the text-matching helper works. It never
# exercises the verdict logic that decides PASS/FAIL from that measurement --
# the missing-name failure, the zero-tests-ran failure, the missing-token
# failure, the empty-inventory failure. That logic used to live only inside
# Invoke-GateGuards in gate.ps1, unreachable from a self-test, so a regression
# that deleted any one of those checks would pass every committed test here.
# These cases call the extracted function directly, with synthetic inventories
# and injected runner script blocks -- no real `go test`, no real pwsh spawn --
# so the failure logic is exercised cheaply and persistently.

$goRunnerNoop = { param($Package, $Pattern) [pscustomobject]@{ ExitCode = 0; Text = '' } }
$pwshRunnerNoop = { param($FilePath) [pscustomobject]@{ ExitCode = 0; Text = '' } }

# 1. a go_tests entry whose name never appears in the runner's output.
$invNamedMissing = [pscustomobject]@{
  go_tests      = @([pscustomobject]@{ id = 'g1'; package = './internal/x'; test = 'TestFoo' })
  pwsh_files    = @()
  presence_only = @()
}
$goRunnerOtherName = { param($Package, $Pattern) [pscustomobject]@{ ExitCode = 0; Text = "=== RUN   TestOther`n--- PASS: TestOther (0.00s)`nPASS`n" } }
$rNamedMissing = Test-GateGuardInventory -Inventory $invNamedMissing -RepositoryRoot $repoRoot -GoRunner $goRunnerOtherName -PwshRunner $pwshRunnerNoop
Assert-True (-not $rNamedMissing.Passed) 'a named-missing go test did not fail the inventory'
Assert-True (@($rNamedMissing.Failures) -contains "no '--- PASS: TestFoo' in ./internal/x") `
  "the named-missing failure line was not reported: $($rNamedMissing.Failures -join '; ')"

# 2. a runner output with zero `=== RUN` lines.
$goRunnerZeroRun = { param($Package, $Pattern) [pscustomobject]@{ ExitCode = 0; Text = "testing: warning: no tests to run`nPASS`nok  `tpkg`t0.001s [no tests to run]`n" } }
$rZeroRun = Test-GateGuardInventory -Inventory $invNamedMissing -RepositoryRoot $repoRoot -GoRunner $goRunnerZeroRun -PwshRunner $pwshRunnerNoop
Assert-True (-not $rZeroRun.Passed) 'a zero-RUN-lines output did not fail the inventory'
Assert-True (@($rZeroRun.Failures) -contains 'package ./internal/x: zero tests ran; the inventory names tests that do not execute') `
  "the zero-tests-ran failure line was not reported: $($rZeroRun.Failures -join '; ')"

# 3. a Go runner exiting non-zero.
$goRunnerNonZero = { param($Package, $Pattern) [pscustomobject]@{ ExitCode = 1; Text = "=== RUN   TestFoo`n--- PASS: TestFoo (0.00s)`nPASS`n" } }
$rGoExit = Test-GateGuardInventory -Inventory $invNamedMissing -RepositoryRoot $repoRoot -GoRunner $goRunnerNonZero -PwshRunner $pwshRunnerNoop
Assert-True (-not $rGoExit.Passed) 'a non-zero go runner exit did not fail the inventory'
Assert-True (@($rGoExit.Failures) -contains 'package ./internal/x exited 1') `
  "the go-exit failure line was not reported: $($rGoExit.Failures -join '; ')"

# 4. a pwsh_files entry whose `guard_ran=<id>` token is absent from the output.
$invPwsh = [pscustomobject]@{
  go_tests      = @()
  pwsh_files    = @([pscustomobject]@{ id = 'p1'; file = 'scripts/tests/gate-measure.tests.ps1'; token = 'guard_ran=synthetic-token' })
  presence_only = @()
}
$pwshRunnerNoToken = { param($FilePath) [pscustomobject]@{ ExitCode = 0; Text = 'PASS something unrelated' } }
$rTokenMissing = Test-GateGuardInventory -Inventory $invPwsh -RepositoryRoot $repoRoot -GoRunner $goRunnerNoop -PwshRunner $pwshRunnerNoToken
Assert-True (-not $rTokenMissing.Passed) 'a missing guard_ran token did not fail the inventory'
Assert-True (@($rTokenMissing.Failures) -contains "no 'guard_ran=synthetic-token' in scripts/tests/gate-measure.tests.ps1") `
  "the missing-token failure line was not reported: $($rTokenMissing.Failures -join '; ')"

# 5. a pwsh runner exiting non-zero.
$pwshRunnerNonZero = { param($FilePath) [pscustomobject]@{ ExitCode = 1; Text = 'guard_ran=synthetic-token' } }
$rPwshExit = Test-GateGuardInventory -Inventory $invPwsh -RepositoryRoot $repoRoot -GoRunner $goRunnerNoop -PwshRunner $pwshRunnerNonZero
Assert-True (-not $rPwshExit.Passed) 'a non-zero pwsh runner exit did not fail the inventory'
Assert-True (@($rPwshExit.Failures) -contains 'scripts/tests/gate-measure.tests.ps1 exited 1') `
  "the pwsh-exit failure line was not reported: $($rPwshExit.Failures -join '; ')"

# 6. a pwsh_files entry pointing outside the selftest discovery glob.
$invOutsideGlob = [pscustomobject]@{
  go_tests      = @()
  pwsh_files    = @([pscustomobject]@{ id = 'p2'; file = 'scripts/harness/Gate.psm1'; token = 'guard_ran=whatever' })
  presence_only = @()
}
$pwshRunnerAny = { param($FilePath) [pscustomobject]@{ ExitCode = 0; Text = 'guard_ran=whatever' } }
$rOutsideGlob = Test-GateGuardInventory -Inventory $invOutsideGlob -RepositoryRoot $repoRoot -GoRunner $goRunnerNoop -PwshRunner $pwshRunnerAny
Assert-True (-not $rOutsideGlob.Passed) 'a pwsh_files entry outside the selftest glob did not fail the inventory'
Assert-True (@($rOutsideGlob.Failures) -contains "scripts/harness/Gate.psm1 sits outside the selftest lane's discovery glob; its execution is delegated to nothing") `
  "the delegation failure line was not reported: $($rOutsideGlob.Failures -join '; ')"

# 7. a presence_only entry whose anchor is absent.
$invPresenceMissing = [pscustomobject]@{
  go_tests      = @()
  pwsh_files    = @()
  presence_only = @([pscustomobject]@{ id = 'pr1'; file = 'scripts/gate.ps1'; anchor = 'ZZZ_NOT_PRESENT_ANCHOR_XYZ' })
}
$rAnchorMissing = Test-GateGuardInventory -Inventory $invPresenceMissing -RepositoryRoot $repoRoot -GoRunner $goRunnerNoop -PwshRunner $pwshRunnerNoop
Assert-True (-not $rAnchorMissing.Passed) 'a missing presence anchor did not fail the inventory'
Assert-True (@($rAnchorMissing.Failures) -contains "scripts/gate.ps1: anchor 'ZZZ_NOT_PRESENT_ANCHOR_XYZ' not found") `
  "the missing-anchor failure line was not reported: $($rAnchorMissing.Failures -join '; ')"

# 8. an empty inventory.
$invEmpty = [pscustomobject]@{ go_tests = @(); pwsh_files = @(); presence_only = @() }
$rEmpty = Test-GateGuardInventory -Inventory $invEmpty -RepositoryRoot $repoRoot -GoRunner $goRunnerNoop -PwshRunner $pwshRunnerNoop
Assert-True (-not $rEmpty.Passed) 'an empty inventory did not fail'
Assert-True ($rEmpty.Entries -eq 0) "an empty inventory reported entries=$($rEmpty.Entries), want 0"
Assert-True ($rEmpty.Reason -match 'empty inventory') "the empty-inventory reason was not reported: $($rEmpty.Reason)"

# 9. a fully satisfied inventory -- zero failures, so the assertions above are
# not trivially satisfied by everything failing.
$invHappy = [pscustomobject]@{
  go_tests      = @([pscustomobject]@{ id = 'g1'; package = './internal/x'; test = 'TestFoo' })
  pwsh_files    = @([pscustomobject]@{ id = 'p1'; file = 'scripts/tests/gate-measure.tests.ps1'; token = 'guard_ran=synthetic-token' })
  presence_only = @([pscustomobject]@{ id = 'pr1'; file = 'scripts/gate.ps1'; anchor = 'Invoke-GateGuards' })
}
$goRunnerHappy = { param($Package, $Pattern) [pscustomobject]@{ ExitCode = 0; Text = "=== RUN   TestFoo`n--- PASS: TestFoo (0.00s)`nPASS`n" } }
$pwshRunnerHappy = { param($FilePath) [pscustomobject]@{ ExitCode = 0; Text = 'guard_ran=synthetic-token' } }
$rHappy = Test-GateGuardInventory -Inventory $invHappy -RepositoryRoot $repoRoot -GoRunner $goRunnerHappy -PwshRunner $pwshRunnerHappy
Assert-True ($rHappy.Passed) "a fully satisfied inventory unexpectedly failed: $($rHappy.Failures -join '; ')"
Assert-True (@($rHappy.Failures).Count -eq 0) 'a fully satisfied inventory reported failures'

Write-Output 'guard_ran=guards-lane-counters'
Write-Output 'PASS guards-lane measurements'
