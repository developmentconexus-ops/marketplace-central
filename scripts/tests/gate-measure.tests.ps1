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

# --- archscan ------------------------------------------------------------

$archscanOutput = @'
internal/modules/integrations/domain/lifecycle.go:17: adapters/vendor-token-outside-adapters: shopee in identifier
internal/modules/orders/application/ingest_service.go:334: adapters/vendor-token-outside-adapters: mercado_livre in string literal
internal/contexts/catalog/internal/domain/product_test.go:158: facts/value-discarded: the bool from .Value() is discarded
archscan: scanned=909 findings=3
'@
$archMeasured = Measure-GateArchscan -Text $archscanOutput
Assert-True ($archMeasured.Scanned -eq 909) "archscan scanned count wrong: $($archMeasured.Scanned)"
Assert-True ($archMeasured.Total -eq 3) "archscan findings count wrong: $($archMeasured.Total)"
Assert-True ($archMeasured.ByRule['adapters/vendor-token-outside-adapters'] -eq 2) 'archscan by-rule count wrong'
Assert-True ($archMeasured.ByRule['facts/value-discarded'] -eq 1) 'archscan by-rule count wrong for facts'

# A stream without the summary line is a run that did not finish, and -1 is not
# 0: "scanned nothing" is a verdict, "nothing is known" is not.
$archDead = Measure-GateArchscan -Text 'internal/x.go:1: facts/value-discarded: the bool'
Assert-True ($archDead.Scanned -eq -1) 'a truncated archscan stream reported a scanned count'
Assert-True ((Measure-GateArchscan -Text '').Scanned -eq -1) 'an empty archscan stream reported a scanned count'

$archCommitted = ConvertTo-GateCountMap -Object $baselineDocument.archscan.by_rule
$archSum = 0
foreach ($value in $archCommitted.Values) { $archSum += [int]$value }
Assert-True ($archSum -eq [int]$baselineDocument.archscan.total) `
  "the archscan baseline's per-rule counts sum to $archSum but it states a total of $($baselineDocument.archscan.total)"

Assert-True ($gateSourceForLint -match '\$measurement\.Scanned -eq 0') `
  'the arch lane no longer fails when archscan parsed zero Go files'
Assert-True ($gateSourceForLint -match '\$measurement\.Scanned -lt 0') `
  'the arch lane no longer fails when archscan died before printing its count'

# --- boundary (ADR-023) ---------------------------------------------------
#
# The literal shape `go test -v` prints for the red-by-design detector. The
# `by target:` block deliberately repeats the count-tab-name shape of the
# `by origin layer:` block, because that is the trap: a regex over the whole
# stream reads every violation twice under a second key set.

$boundaryRed = @"
=== RUN   TestModuleBoundaryADR023
    module_boundary_arch_test.go:153: boundary_files=782
    module_boundary_arch_test.go:230: 5 violation(s)
        ADR-023 §2 module boundary violated

        by origin layer:
          3`tadapters
          1`tapplication
          1`t(module root, no layer)

        by target:
          4`tconnectors/domain
          1`tinternal_read/domain

        sites:
          internal/modules/catalog/adapters/x.go:10`tcatalog/adapters -> connectors/domain
--- FAIL: TestModuleBoundaryADR023 (0.10s)
FAIL
FAIL	marketplace-central/apps/server_core/internal/composition	0.427s
FAIL
"@
$boundaryMeasured = Measure-GateBoundary -Text $boundaryRed
Assert-True ($boundaryMeasured.Files -eq 782) "boundary files count wrong: $($boundaryMeasured.Files)"
Assert-True ($boundaryMeasured.Total -eq 5) "boundary violation count wrong: $($boundaryMeasured.Total)"
Assert-True ($boundaryMeasured.ByOrigin['adapters'] -eq 3) 'boundary by-origin count wrong'
Assert-True ($boundaryMeasured.ByOrigin['(module root, no layer)'] -eq 1) 'the unsanctioned module-root origin was dropped'
Assert-True (@($boundaryMeasured.ByOrigin.Keys).Count -eq 3) `
  "boundary read $(@($boundaryMeasured.ByOrigin.Keys).Count) origins; the by-target block leaked into the origin map"

# A passing run still names its file count -- that is why the test logs it on
# every outcome -- and measures as zero violations, not as unknown.
$boundaryGreen = @"
=== RUN   TestModuleBoundaryADR023
    module_boundary_arch_test.go:153: boundary_files=782
--- PASS: TestModuleBoundaryADR023 (0.09s)
PASS
ok  	marketplace-central/apps/server_core/internal/composition	0.412s
"@
$boundaryGreenMeasured = Measure-GateBoundary -Text $boundaryGreen
Assert-True ($boundaryGreenMeasured.Files -eq 782) 'a passing boundary run lost its file count'
Assert-True ($boundaryGreenMeasured.Total -eq 0) "a passing boundary run measured $($boundaryGreenMeasured.Total) violations, not zero"

# A red run whose origin block went missing (or half-parsed) still reports a
# violation total. The measurer returns what it saw -- an empty origin map next
# to a nonzero total -- and the LANE must refuse the mismatch, because the
# ratchet compares origins and an all-zero map would sail under any baseline.
$boundaryNoOrigins = @"
=== RUN   TestModuleBoundaryADR023
    module_boundary_arch_test.go:153: boundary_files=782
    module_boundary_arch_test.go:228: 234 violation(s)
        ADR-023 §2 module boundary violated
--- FAIL: TestModuleBoundaryADR023 (0.10s)
FAIL
"@
$boundaryNoOriginsMeasured = Measure-GateBoundary -Text $boundaryNoOrigins
Assert-True ($boundaryNoOriginsMeasured.Total -eq 234) 'the originless fixture lost its violation total'
Assert-True (@($boundaryNoOriginsMeasured.ByOrigin.Keys).Count -eq 0) 'the originless fixture grew origins from nowhere'
Assert-True ($gateSourceForLint -match '\$originTotal -ne \$measurement\.Total') `
  'the boundary lane no longer rejects a violation total that the by-origin counts do not account for'

# A compile error prints neither the log line nor a verdict. Both facts default
# to -1: "nothing is known" must not read as "walked nothing" or "found nothing".
$boundaryDead = Measure-GateBoundary -Text "FAIL	marketplace-central/apps/server_core/internal/composition [build failed]"
Assert-True ($boundaryDead.Files -eq -1) 'a build-failed boundary run reported a file count'
Assert-True ($boundaryDead.Total -eq -1) 'a build-failed boundary run reported a violation count'

$boundaryCommitted = ConvertTo-GateCountMap -Object $baselineDocument.boundary.by_origin
$boundarySum = 0
foreach ($value in $boundaryCommitted.Values) { $boundarySum += [int]$value }
Assert-True ($boundarySum -eq [int]$baselineDocument.boundary.total) `
  "the boundary baseline's per-origin counts sum to $boundarySum but it states a total of $($baselineDocument.boundary.total)"

Assert-True ($gateSourceForLint -match '\$measurement\.Files -eq 0') `
  'the boundary lane no longer fails when the walk parsed zero Go files'
Assert-True ($gateSourceForLint -match '\$measurement\.Files -lt 0') `
  'the boundary lane no longer fails when the run died before the walk finished'
Assert-True ($gateSourceForLint -match 'boundary_files') `
  'the boundary lane no longer reads the boundary_files count at all'
$boundaryTestSource = Get-Content -Raw -LiteralPath (Join-Path $repoRoot 'apps/server_core/internal/composition/module_boundary_arch_test.go')
Assert-True ($boundaryTestSource -match [regex]::Escape('t.Logf("boundary_files=%d", filesParsed)')) `
  'the detector no longer logs its file count on every outcome; the lane cannot tell an empty walk from a clean tree'

# --- governance ------------------------------------------------------------
#
# No text fixture: the lane calls Test-GovernanceDrift directly and groups its
# structured violations. What is tested is the grouping and the lane's guards.

$governanceMeasured = Measure-GateGovernance -Violations @(
  [pscustomobject]@{ ErrorCode = 'RCFG_UNAPPROVED_READER'; Id = 'MC_X'; Path = 'a.go' },
  [pscustomobject]@{ ErrorCode = 'RCFG_UNAPPROVED_READER'; Id = 'MC_Y'; Path = 'b.go' },
  [pscustomobject]@{ ErrorCode = 'GOV_MODULE_LAYER'; Id = 'catalog-sync-adapters'; Path = 'c.go' }
)
Assert-True ($governanceMeasured.Total -eq 3) "governance total wrong: $($governanceMeasured.Total)"
Assert-True ($governanceMeasured.ByCode['RCFG_UNAPPROVED_READER'] -eq 2) 'governance by-code count wrong'
Assert-True ((Measure-GateGovernance -Violations @()).Total -eq 0) 'an empty violation list did not measure as zero'

$governanceCommitted = ConvertTo-GateCountMap -Object $baselineDocument.'governance-drift'.by_code
$governanceSum = 0
foreach ($value in $governanceCommitted.Values) { $governanceSum += [int]$value }
Assert-True ($governanceSum -eq [int]$baselineDocument.'governance-drift'.total) `
  "the governance baseline's per-code counts sum to $governanceSum but it states a total of $($baselineDocument.'governance-drift'.total)"

Assert-True ($gateSourceForLint -match '\$drift\.FilesScanned -le 0') `
  'the governance lane no longer fails when the drift walk scanned no files'
Assert-True ($gateSourceForLint -match '-not \$validate\.Passed') `
  'the governance lane no longer blocks at zero on a validate failure'
Assert-True ($gateSourceForLint -match '\$registries -lt 6') `
  'the governance lane no longer asserts that all six registries loaded'
$policySource = Get-Content -Raw -LiteralPath (Join-Path $repoRoot 'scripts/harness/Policy.psm1')
Assert-True ($policySource -match '-FilesScanned \(@\(\$files\)\.Count\)') `
  'Test-GovernanceDrift no longer reports how many files its walk covered; the lane cannot tell an empty walk from a clean tree'

# --- eslint --------------------------------------------------------------

$eslintReport = @'
[
  { "filePath": "C:\\repo\\apps\\web\\src\\a.ts", "messages": [
      { "ruleId": "@typescript-eslint/no-floating-promises", "severity": 2 },
      { "ruleId": "react-hooks/exhaustive-deps", "severity": 2 } ] },
  { "filePath": "C:\\repo\\packages\\ui\\src\\b.tsx", "messages": [
      { "ruleId": "@typescript-eslint/no-floating-promises", "severity": 2 } ] },
  { "filePath": "C:\\repo\\packages\\ui\\src\\c.tsx", "messages": [] }
]
'@
$eslintMeasured = Measure-GateEslint -Json $eslintReport
Assert-True ($eslintMeasured.Files -eq 3) "eslint file count wrong: $($eslintMeasured.Files)"
Assert-True ($eslintMeasured.Total -eq 3) "eslint finding count wrong: $($eslintMeasured.Total)"
Assert-True ($eslintMeasured.ByRule['@typescript-eslint/no-floating-promises'] -eq 2) 'eslint per-rule count wrong'

# The file with no messages must still be reported as covered. It is the whole
# reason the lane reads `Paths` rather than only counting findings: a file that
# was linted and had nothing to say and a file that was never linted are the same
# zero in the findings total, and only the path list tells them apart.
Assert-True (@($eslintMeasured.Paths) -contains 'C:/repo/packages/ui/src/c.tsx') `
  'a file with no findings was dropped from the coverage list'

# A parse failure has no ruleId. Folding it into the findings total would let a
# file that could not be read at all shrink the ratchet the moment it is ignored.
$eslintFatal = @'
[ { "filePath": "C:\\repo\\x.ts", "messages": [ { "fatal": true, "message": "Parsing error", "severity": 2 } ] } ]
'@
$fatalMeasured = Measure-GateEslint -Json $eslintFatal
Assert-True ($fatalMeasured.Fatal -eq 1) 'a parse failure was not counted as fatal'
Assert-True ($fatalMeasured.Total -eq 0) 'a parse failure was counted as a rule finding'

Assert-True ((Measure-GateEslint -Json '').Files -eq -1) 'a missing eslint report measured as zero files linted'

# --- prettier ------------------------------------------------------------

$prettierClean = @'
Checking formatting...
All matched files use Prettier code style!
'@
$prettierMeasured = Measure-GatePrettier -Text $prettierClean
Assert-True ($prettierMeasured.Clean) 'a clean prettier run did not set the clean marker'
Assert-True ($prettierMeasured.Unformatted -eq 0) 'a clean prettier run reported unformatted files'

$prettierDirty = @'
Checking formatting...
[warn] apps/web/src/App.tsx
[warn] packages/ui/src/Badge.tsx
[warn] Code style issues found in 2 files. Run Prettier with --write to fix.
'@
$dirtyMeasured = Measure-GatePrettier -Text $prettierDirty
Assert-True (-not $dirtyMeasured.Clean) 'an unformatted tree still set the clean marker'
Assert-True ($dirtyMeasured.Unformatted -eq 2) "prettier unformatted count wrong: $($dirtyMeasured.Unformatted)"
Assert-True (@($dirtyMeasured.Paths).Count -eq 2) 'the summary line was counted as an unformatted path'
Assert-True (@($dirtyMeasured.Paths) -contains 'apps/web/src/App.tsx') 'prettier did not name the unformatted file'

# Silence is the case this exists for. A run that died before printing anything
# has no warnings either, and the clean marker is what separates the two.
$prettierSilent = Measure-GatePrettier -Text ''
Assert-True (-not $prettierSilent.Clean) 'an empty prettier stream was read as a clean tree'

# The clean marker over the empty set, which is the case the marker cannot rule
# out on its own. Reported by CodeRabbit on PR #21 and reproduced 2026-08-08:
# with every extension in .prettierignore the lane read zero files and still
# passed. `Checked` comes from the debug stream and is the only number that moves.
Assert-True ($prettierMeasured.Checked -eq 0) `
  'a run without --log-level debug reported files as checked; the caller must be able to see that number is unavailable'

$prettierDebug = @'
[debug] resolve config from 'C:\repo\apps\web\src\App.tsx'
[debug] loaded cached config
[debug] applied config
[debug] resolve config from 'C:\repo\apps\web\src\App.tsx'
[debug] resolve config from 'C:\repo\packages\ui\src\Badge.tsx'
Checking formatting...
All matched files use Prettier code style!
'@
$debugMeasured = Measure-GatePrettier -Text $prettierDebug
Assert-True ($debugMeasured.Checked -eq 2) `
  "prettier checked-file count wrong: $($debugMeasured.Checked); the same path resolving twice is one file"
Assert-True ($debugMeasured.Clean) 'the debug stream lost the clean marker'

# --- the eslint rule set ---------------------------------------------------
#
# A rule that is off produces zero findings and is indistinguishable from a rule
# that found nothing. Reproduced 2026-08-08: two rules set to "off" dropped the
# lane's total from 26 to 13 and the ratchet printed the drop as a shrink worth
# committing. Severity arrives as a number, a string, or the head of an array.
$printConfig = @'
{"rules":{
  "@typescript-eslint/no-floating-promises":[2],
  "@typescript-eslint/no-misused-promises":"off",
  "@typescript-eslint/no-unused-vars":[0,{"argsIgnorePattern":"^_"}],
  "react-hooks/exhaustive-deps":"warn",
  "react-hooks/rules-of-hooks":2
}}
'@
$rulesMeasured = @(Measure-GateEslintRules -Json $printConfig)
Assert-True ($rulesMeasured -contains '@typescript-eslint/no-floating-promises') 'an enabled rule was dropped from the effective set'
Assert-True ($rulesMeasured -contains 'react-hooks/exhaustive-deps') 'a rule at "warn" was read as off'
Assert-True ($rulesMeasured -contains 'react-hooks/rules-of-hooks') 'a bare numeric severity was read as off'
Assert-True (-not ($rulesMeasured -contains '@typescript-eslint/no-misused-promises')) 'a rule set to "off" was reported as enabled'
Assert-True (-not ($rulesMeasured -contains '@typescript-eslint/no-unused-vars')) 'a rule with severity 0 in an options array was reported as enabled'
Assert-True (@(Measure-GateEslintRules -Json '').Count -eq 0) 'an empty --print-config stream was read as a rule set'

# --- the eslint baseline, and the two lanes' guards -----------------------

$eslintCommitted = ConvertTo-GateCountMap -Object $baselineDocument.eslint.by_rule
$eslintSum = 0
foreach ($value in $eslintCommitted.Values) { $eslintSum += [int]$value }
Assert-True ($eslintSum -eq [int]$baselineDocument.eslint.total) `
  "the eslint baseline's per-rule counts sum to $eslintSum but it states a total of $($baselineDocument.eslint.total)"

# Every rule the baseline names must be a rule the config turns on. A baseline
# entry for a rule nobody enabled is a number that can never move, and it makes
# the total look larger than the debt actually is.
$eslintConfig = Get-Content -Raw -LiteralPath (Join-Path $repoRoot 'eslint.config.mjs')
foreach ($rule in @($eslintCommitted.Keys)) {
  Assert-True ($eslintConfig -match [regex]::Escape("`"$rule`"")) `
    "the eslint baseline names $rule but eslint.config.mjs does not enable it"
}

# The type-aware rules are the reason `tsconfig.eslint.json` exists. Without a
# project they cannot run, and ESLint still reports the file as linted.
Assert-True (Test-Path -LiteralPath (Join-Path $repoRoot 'tsconfig.eslint.json') -PathType Leaf) `
  'RED: tsconfig.eslint.json missing; the type-aware rules have no project to resolve against'
Assert-True ($eslintConfig -match 'project:\s*\["\./tsconfig\.eslint\.json"\]') `
  'eslint.config.mjs no longer points the parser at the lint project; the type-aware rules go silent without failing'

Assert-True ($gateSourceForLint -match '\$uncovered\.Count -gt 0') `
  'the lint-web lane no longer fails when a tracked source was not linted'
Assert-True ($gateSourceForLint -match '\$measurement\.Fatal -gt 0') `
  'the lint-web lane no longer fails on a file ESLint could not parse'
Assert-True ($gateSourceForLint -match '\$candidates -eq 0') `
  'the format lane no longer fails when the formatted-file set is empty'
Assert-True ($gateSourceForLint -match '-not \$measurement\.Clean') `
  'the format lane once again reads silence as a clean tree'

# The two guards CodeRabbit's review bought, pinned so they cannot be dropped
# quietly. Each names a lane that passed over nothing before it existed.
Assert-True ($gateSourceForLint -match '\$measurement\.Checked -eq 0') `
  'the format lane no longer fails when prettier read no file at all'
Assert-True ($gateSourceForLint -match "--log-level', 'debug'") `
  'the format lane dropped --log-level debug; without it prettier names no file and Checked is always 0'
#
# Anchored on the reason text, not on `$enabledDrift`: the lint-go lane uses that
# same variable name, so a pattern naming it matches whether or not the lint-web
# guard still exists. The first draft of this assertion did exactly that.
Assert-True ($gateSourceForLint -match 'the ratchet would take the drop for a shrink') `
  'the lint-web lane no longer asserts its enabled rule set; a rule turned off would read as a shrink'
Assert-True ($gateSourceForLint -match 'Measure-GateEslintRules') `
  'the lint-web lane no longer measures the effective rule set at all'
Assert-True ($gateSourceForLint -match '--print-config') `
  'the lint-web lane no longer resolves the effective config; the text of eslint.config.mjs is not what ESLint runs'

# The configured rule set and the baseline's keys are one declaration in two
# files. A rule added to the config without a baseline entry, or the reverse,
# means every count below is over a set nobody declared.
$expectedRules = @([regex]::Matches($gateSourceForLint, "(?m)^\s*'(?<rule>[@a-z-]+/[a-z-]+)',?\s*$") |
  ForEach-Object { $_.Groups['rule'].Value } | Sort-Object -Unique)
Assert-True ($expectedRules.Count -eq @($eslintCommitted.Keys).Count) `
  "gate.ps1 declares $($expectedRules.Count) eslint rules and the baseline names $(@($eslintCommitted.Keys).Count)"
foreach ($rule in @($eslintCommitted.Keys)) {
  Assert-True ($expectedRules -contains $rule) "the baseline names $rule and gate.ps1's expected rule set does not"
}

Write-Output 'PASS gate measurement tests'
