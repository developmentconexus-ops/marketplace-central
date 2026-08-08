$ErrorActionPreference = 'Stop'

# The gate's counting functions, tested against recorded tool output.
#
# The fixtures below are not invented shapes. Each block is the literal form the
# tool emits, including the ANSI escapes vitest wraps its summary in -- a bare
# `^ *Tests` match never fires against the real stream, which is why the plain
# text is produced by a function under test rather than assumed.
#
# The cases that matter are the vacuous ones. A suite that ran nothing and a
# suite that passed everything both exit 0; every assertion here that expects a
# zero is asserting that the difference is visible.

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

# --- go test -------------------------------------------------------------

$goHappy = @"
=== RUN   TestOne
--- PASS: TestOne (0.00s)
=== RUN   TestTwo
=== RUN   TestTwo/subcase
    --- PASS: TestTwo/subcase (0.00s)
--- PASS: TestTwo (0.00s)
PASS
ok  	example/pkg	0.012s
"@
$goMeasured = Measure-GateGoTest -Text $goHappy
Assert-True ($goMeasured.Run -eq 3) "go run count wrong: $($goMeasured.Run)"
Assert-True ($goMeasured.Passed -eq 3) "go pass count wrong: $($goMeasured.Passed)"
Assert-True ($goMeasured.Failed -eq 0) 'go fail count wrong'
Assert-True ($goMeasured.Skipped -eq 0) 'go skip count wrong'

# The whole reason this lane counts at all. `go test` prints this and exits 0.
$goEmpty = @"
testing: warning: no tests to run
PASS
ok  	example/pkg	0.001s [no tests to run]
"@
$goEmptyMeasured = Measure-GateGoTest -Text $goEmpty
Assert-True ($goEmptyMeasured.Run -eq 0) 'an empty package set did not measure as zero runs'
Assert-True ($goEmptyMeasured.Passed -eq 0) 'an empty package set did not measure as zero passes'

# HARNESS-PROFILE.md records this exact shape: every test skipped, package `ok`,
# exit 0, with the slice's whole reason to exist among the skips. Run>0 alone
# would accept it, which is why the gate's guard reads run==0 OR passed==0.
$goAllSkipped = @"
=== RUN   TestOne
--- SKIP: TestOne (0.00s)
=== RUN   TestTwo
--- SKIP: TestTwo (0.00s)
PASS
ok  	example/pkg	0.003s
"@
$goSkipMeasured = Measure-GateGoTest -Text $goAllSkipped
Assert-True ($goSkipMeasured.Run -eq 2) 'a fully skipped suite lost its run count'
Assert-True ($goSkipMeasured.Skipped -eq 2) 'a fully skipped suite lost its skip count'
Assert-True ($goSkipMeasured.Passed -eq 0) 'a fully skipped suite reported passes it did not have'

$goFailing = @"
=== RUN   TestOne
--- FAIL: TestOne (0.00s)
FAIL
"@
Assert-True ((Measure-GateGoTest -Text $goFailing).Failed -eq 1) 'a failure was not counted'

# --- vitest --------------------------------------------------------------

$esc = [char]27
$vitestHappy = "$esc[2m Test Files $esc[22m $esc[1m$esc[32m72 passed$esc[39m$esc[22m$esc[90m (72)$esc[39m`n$esc[2m      Tests $esc[22m $esc[1m$esc[32m605 passed$esc[39m$esc[22m$esc[90m (605)$esc[39m"
$vitestMeasured = Measure-GateVitest -Text $vitestHappy
Assert-True ($vitestMeasured.Files -eq 72) "vitest file count wrong: $($vitestMeasured.Files)"
Assert-True ($vitestMeasured.Tests -eq 605) "vitest test count wrong: $($vitestMeasured.Tests)"
Assert-True ($vitestMeasured.Passed -eq 605) "vitest pass count wrong: $($vitestMeasured.Passed)"
Assert-True ($vitestMeasured.Failed -eq 0) 'vitest reported failures on a clean run'

Assert-True ((Remove-GateAnsi -Text $vitestHappy) -match '(?m)^\s*Tests\s+605 passed') 'ANSI stripping did not expose the summary line'

$vitestMixed = "  Tests  2 failed | 603 passed (605)"
$vitestMixedMeasured = Measure-GateVitest -Text $vitestMixed
Assert-True ($vitestMixedMeasured.Failed -eq 2) "vitest mixed failure count wrong: $($vitestMixedMeasured.Failed)"
Assert-True ($vitestMixedMeasured.Passed -eq 603) "vitest mixed pass count wrong: $($vitestMixedMeasured.Passed)"

# A run that crashed before the summary, and a run that matched no test files.
# Neither may measure as zero-and-fine: -1 says the result was never stated,
# which is a different fact from a suite that ran and found nothing.
$vitestCrashed = "RUN  v3.0.0`nError: Cannot find module './setup'"
Assert-True ((Measure-GateVitest -Text $vitestCrashed).Passed -eq -1) 'a crashed vitest run reported a pass count'
Assert-True ((Measure-GateVitest -Text $vitestCrashed).Failed -eq -1) 'a crashed vitest run claimed zero failures'
$vitestNoFiles = "No test files found, exiting with code 0"
Assert-True ((Measure-GateVitest -Text $vitestNoFiles).Passed -eq -1) 'a no-files vitest run reported a pass count'
Assert-True ((Measure-GateVitest -Text $vitestNoFiles).Failed -eq -1) 'a no-files vitest run claimed zero failures'

# --- tsc -----------------------------------------------------------------

$tscHappy = @"
C:/repo/node_modules/typescript/lib/lib.es5.d.ts
C:/repo/node_modules/@types/react/index.d.ts
C:/repo/apps/web/src/main.tsx
C:/repo/apps/web/src/app/Client.ts
"@
$tscMeasured = Measure-GateTsc -Text $tscHappy
Assert-True ($tscMeasured.Checked -eq 2) "tsc project file count wrong: $($tscMeasured.Checked)"
Assert-True ($tscMeasured.Errors -eq 0) 'tsc reported errors on a clean run'

# node_modules only: a project whose `include` resolved nothing still loads the
# toolchain's own 300-odd declaration files, so counting every listed path would
# report a healthy number over zero source files.
$tscNoProjectFiles = @"
C:/repo/node_modules/typescript/lib/lib.es5.d.ts
C:/repo/node_modules/typescript/lib/lib.dom.d.ts
"@
Assert-True ((Measure-GateTsc -Text $tscNoProjectFiles).Checked -eq 0) 'a project that resolved no sources reported files checked'

$tscErrors = @"
C:/repo/apps/web/src/main.tsx
apps/web/src/components/MutationPreviewModal.tsx(210,9): error TS2741: Property 'onRetry' is missing.
"@
$tscErrorMeasured = Measure-GateTsc -Text $tscErrors
Assert-True ($tscErrorMeasured.Errors -eq 1) "tsc error count wrong: $($tscErrorMeasured.Errors)"
Assert-True ($tscErrorMeasured.Checked -eq 1) 'a tsc error line was counted as a project file'

# --- gofmt ---------------------------------------------------------------

$gofmtClean = ''
Assert-True ((Measure-GateGofmt -Text $gofmtClean).Unformatted -eq 0) 'a clean gofmt run reported violations'
$gofmtDirty = "internal/a/a.go`ninternal/b/b.go`n"
$gofmtMeasured = Measure-GateGofmt -Text $gofmtDirty
Assert-True ($gofmtMeasured.Unformatted -eq 2) "gofmt violation count wrong: $($gofmtMeasured.Unformatted)"
Assert-True ($gofmtMeasured.Paths -contains 'internal/b/b.go') 'gofmt did not report the offending path'

# --- the guard, read off the entry point ---------------------------------
#
# The measurers above can be right while the lane that consumes them is wrong.
# These assertions pin the two conditions that turn a count into a verdict, so
# that weakening either one -- dropping the `-eq 0` on passes, or accepting a
# tsc run over zero files -- fails here rather than silently on a green PR.

$gateSource = Get-Content -Raw -LiteralPath $gateScript
Assert-True ($gateSource -match '\$measurement\.Run -eq 0 -or \$measurement\.Passed -eq 0') `
  'the go test lane no longer fails on a suite that ran nothing or passed nothing'
Assert-True ($gateSource -match '\$measurement\.Checked -eq 0') `
  'the typecheck lane no longer fails on a project that resolved no sources'
Assert-True ($gateSource -match '\$measurement\.Passed -le 0') `
  'the vitest lane no longer fails on a run with no stated result'
Assert-True ($gateSource -match '\$discovered\.Count -eq 0') `
  'the gofmt lane no longer fails when the file filter reaches nothing'
# Measured 2026-08-08: linux/cgo=1 reports errcheck=30, windows/cgo=0 reports 28.
# Without this, a Windows run reports a shrink it did not earn and someone
# commits a baseline the enforcing platform cannot meet.
Assert-True ($gateSource -match '\$platformMatches') `
  'the lint lane no longer distinguishes a real shrink from a smaller analysed file set'
Assert-True ($gateSource -match 'NOT a shrink') `
  'the lint lane once again prints a cross-platform decrease as though it were progress'
Assert-True ($gateSource -match '\$gofmtExitFailures -gt 0') `
  'the gofmt lane no longer fails when gofmt itself exits non-zero'
Assert-True ($gateSource -match 'Test-Path -LiteralPath \(Join-Path \$repositoryRoot \$_\) -PathType Leaf') `
  'the gofmt lane no longer filters index paths deleted from the worktree'

# --- golangci-lint -------------------------------------------------------

$lintReport = @'
{
  "Issues": [
    { "FromLinter": "errcheck",  "Pos": { "Filename": "internal/a/a.go" } },
    { "FromLinter": "errcheck",  "Pos": { "Filename": "internal/b/b.go" } },
    { "FromLinter": "bodyclose", "Pos": { "Filename": "internal/c/c.go" } }
  ],
  "Report": {
    "Linters": [
      { "Name": "errcheck",  "Enabled": true },
      { "Name": "bodyclose", "Enabled": true },
      { "Name": "gocyclo" }
    ]
  }
}
'@
$lintMeasured = Measure-GateGolangciLint -Json $lintReport
Assert-True ($lintMeasured.Total -eq 3) "golangci total wrong: $($lintMeasured.Total)"
Assert-True ($lintMeasured.ByLinter['errcheck'] -eq 2) 'golangci per-linter count wrong'
Assert-True (@($lintMeasured.Enabled) -join ',' -eq 'bodyclose,errcheck') "golangci enabled set wrong: $(@($lintMeasured.Enabled) -join ',')"

# A clean run still writes a report, and `Issues` is null rather than an empty
# array. Reading that as -1 would make a genuinely clean lane look like a lane
# that never ran.
$lintClean = '{ "Issues": null, "Report": { "Linters": [ { "Name": "errcheck", "Enabled": true } ] } }'
Assert-True ((Measure-GateGolangciLint -Json $lintClean).Total -eq 0) 'a clean golangci report did not measure as zero issues'

# No report at all is a different fact: the tool did not run to completion, and
# the lane must not report its findings as zero.
Assert-True ((Measure-GateGolangciLint -Json '').Total -eq -1) 'a missing golangci report measured as zero findings'

# --- the ratchet ---------------------------------------------------------

$baseline = @{ errcheck = 28; bodyclose = 8; unused = 7 }

$flat = Compare-GateRatchet -Measured @{ errcheck = 28; bodyclose = 8; unused = 7 } -Baseline $baseline
Assert-True ($flat.Passed) 'an unchanged count did not pass the ratchet'
Assert-True ($flat.Increased.Count -eq 0 -and $flat.Decreased.Count -eq 0) 'an unchanged count reported movement'

$grew = Compare-GateRatchet -Measured @{ errcheck = 29; bodyclose = 8; unused = 7 } -Baseline $baseline
Assert-True (-not $grew.Passed) 'the ratchet passed an increase'
Assert-True ($grew.Increased[0] -match 'errcheck=29') 'the ratchet did not name the linter that grew'

$shrank = Compare-GateRatchet -Measured @{ errcheck = 20; bodyclose = 8; unused = 7 } -Baseline $baseline
Assert-True ($shrank.Passed) 'the ratchet blocked a decrease'
Assert-True ($shrank.Decreased[0] -match 'errcheck=20') 'a decrease was not reported for committing'

# A linter that appears in the run and not in the baseline. Treating it as a
# baseline of zero would be stricter and sounds safer, but it hides the drift
# behind a finding about code; the drift is the thing to report.
$unknown = Compare-GateRatchet -Measured @{ errcheck = 28; bodyclose = 8; unused = 7; noctx = 1 } -Baseline $baseline
Assert-True (-not $unknown.Passed) 'the ratchet passed a linter the baseline does not name'
Assert-True ($unknown.Unknown[0] -match 'noctx=1') 'the unnamed linter was not reported'

# The opposite direction is a shrink to zero, which is where this is supposed to
# go -- not an error.
$gone = Compare-GateRatchet -Measured @{ errcheck = 28; bodyclose = 8 } -Baseline $baseline
Assert-True ($gone.Passed) 'a linter reaching zero findings failed the ratchet'
Assert-True ($gone.Decreased[0] -match 'unused=0') 'a linter reaching zero was not reported as a shrink'

# --- the committed baseline is the one the lane reads --------------------
#
# A baseline whose per-linter numbers do not add up to its own stated total is
# how a ratchet quietly gains headroom.

$baselinePath = Join-Path $repoRoot 'contracts/gate/baselines.json'
Assert-True (Test-Path -LiteralPath $baselinePath -PathType Leaf) 'RED: no ratchet baseline committed'
$baselineDocument = Get-Content -LiteralPath $baselinePath -Raw | ConvertFrom-Json
$committed = ConvertTo-GateCountMap -Object $baselineDocument.'golangci-lint'.by_linter
$sum = 0
foreach ($value in $committed.Values) { $sum += [int]$value }
Assert-True ($sum -eq [int]$baselineDocument.'golangci-lint'.total) `
  "the golangci-lint baseline's per-linter counts sum to $sum but it states a total of $($baselineDocument.'golangci-lint'.total)"

$gateSourceForLint = Get-Content -Raw -LiteralPath $gateScript
foreach ($linter in @($committed.Keys)) {
  Assert-True ($gateSourceForLint -match [regex]::Escape("'$linter'")) `
    "the baseline names $linter but the lane's expected enabled set does not"
}
Assert-True ($gateSourceForLint -match '\$enabledDrift\.Count -ne 0') `
  'the lint lane no longer fails when the enabled linter set drifts from the configured one'

Write-Output 'PASS gate measurement tests'
