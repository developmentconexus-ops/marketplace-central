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

Write-Output 'guard_ran=guards-lane-counters'
Write-Output 'PASS guards-lane measurements'
