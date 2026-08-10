#Requires -Version 7.0
<#
.SYNOPSIS
  The verification gate. One entry point, run identically on a developer machine
  and on a CI runner.

.DESCRIPTION
  Issue #2 of the repository audit of 2026-08-07: "the verifier is a set of
  scripts to remember, not one product with one entry point". The checks already
  existed and were good; nothing invoked them as a set, so which ones ran on any
  given change was a function of what someone remembered.

  This file is the set: gofmt, Prettier, `go build` + `go vet`, typecheck,
  governance (contract validate + semantic drift), archscan, the ADR-023
  module-boundary detector, golangci-lint, ESLint, `go test ./...`, vitest.
  `npm run gate` runs it locally; `.github/workflows/ci.yml`
  invokes the same lanes by name. Neither owns a copy of the logic -- a gate the
  operator cannot reproduce on their own machine is its own defect class, and a
  gate whose CI form has drifted from its local form is that class wearing a
  disguise.

  Every lane prints an attributable count and fails when that count is zero. A
  tool that executed nothing exits 0, so "green" and "did not run" are otherwise
  byte-identical. That equivalence is the finding, not a detail of it.

.PARAMETER Lane
  Which lane to run. 'all' runs every lane in cost order -- cheap lanes first, so
  a twenty-second formatting error does not first burn several minutes of
  compilation.

.PARAMETER ContinueOnFailure
  Run every remaining lane after one fails, instead of stopping. Useful locally
  when the intent is to see the whole picture before starting to fix it. The exit
  code is unaffected: any failed lane still exits non-zero.
#>
[CmdletBinding()]
param(
  [ValidateSet('all', 'full', 'gofmt', 'format', 'build', 'typecheck', 'governance', 'arch', 'boundary', 'lint-go', 'lint-web', 'test-go', 'test-web', 'census', 'ci-hygiene', 'selftest', 'guards', 'integration', 'edge')]
  [string]$Lane = 'all',
  [switch]$ContinueOnFailure
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$repositoryRoot = [IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..'))
$serverCore = Join-Path $repositoryRoot 'apps/server_core'
Import-Module (Join-Path $PSScriptRoot 'harness/Gate.psm1') -Force
# For Get-HarnessIntegrationTestPackages: the census lane reuses the harness's
# own integration-package discovery instead of keeping a second copy of the rule.
Import-Module (Join-Path $PSScriptRoot 'harness/Postgres.psm1') -Force

# HARNESS-PROFILE §2 binds the Go caches to an ABSOLUTE path. The relative
# `GOCACHE=.gocache` that AGENTS.md still prints was struck from the profile on
# 2026-07-28 after a chip copied it and got an 83-byte EXIT 1 with zero `=== RUN`
# -- a vacuous green handed over by the doctrine itself.
$env:GOCACHE = [IO.Path]::GetFullPath((Join-Path $serverCore '.gocache'))
$env:GOMODCACHE = [IO.Path]::GetFullPath((Join-Path $serverCore '.gomodcache'))

$logDirectory = Join-Path $repositoryRoot '.gate'
if (-not (Test-Path -LiteralPath $logDirectory)) {
  [void](New-Item -ItemType Directory -Path $logDirectory -Force)
}

function Invoke-GateTool {
  <#
    Runs a tool, streams its output to the console AND captures it for counting.
    Both halves matter: the operator watching a lane needs to see it work, and
    the lane needs the text to prove it did. `scripts/harness/Postgres.psm1`
    carried the opposite arrangement until 2026-08-08 -- the success path never
    read stdout, so a lane that compiled zero test packages and one that ran
    every integration test produced byte-identical evidence.
  #>
  param(
    [Parameter(Mandatory)][string]$Name,
    [Parameter(Mandatory)][string]$FilePath,
    [Parameter(Mandatory)][AllowEmptyCollection()][string[]]$ArgumentList,
    [string]$WorkingDirectory = $repositoryRoot,
    # For a tool run at a debug verbosity that the lane reads but the operator
    # should not have to. The log file is written either way, so the evidence is
    # identical; only the console differs.
    [switch]$Quiet
  )

  $logPath = Join-Path $logDirectory "$Name.log"
  Push-Location -LiteralPath $WorkingDirectory
  try {
    # 2>&1 folds stderr into the stream: go and tsc write real diagnostics there,
    # and a count taken over stdout alone would miss them.
    if ($Quiet) {
      & $FilePath @ArgumentList 2>&1 | Out-File -LiteralPath $logPath -Encoding utf8
    } else {
      & $FilePath @ArgumentList 2>&1 | Tee-Object -FilePath $logPath | Out-Host
    }
    $exitCode = $LASTEXITCODE
  } finally {
    Pop-Location
  }
  $text = if (Test-Path -LiteralPath $logPath) { Get-Content -LiteralPath $logPath -Raw } else { '' }
  return [pscustomobject]@{
    ExitCode = $exitCode
    Text     = if ($null -eq $text) { '' } else { $text }
    LogPath  = $logPath
  }
}

function New-GateVerdict {
  param(
    [Parameter(Mandatory)][string]$Lane,
    [Parameter(Mandatory)][bool]$Passed,
    [Parameter(Mandatory)][string]$Counts,
    [string]$Reason = ''
  )
  return [pscustomobject]@{
    Lane   = $Lane
    Passed = $Passed
    Counts = $Counts
    Reason = $Reason
  }
}

function Invoke-GateGofmt {
  # `git ls-files`, not a walk of the working tree. On Windows a tree walk reaches
  # `.gomodcache` -- 245 MB of vendored sources that are not this repository's to
  # format, and which once inflated a measured 404 to 2721.
  #
  # --others --exclude-standard, because `--cached` alone sees only tracked files.
  # Measured 2026-08-08: a deliberately misformatted new file placed in the tree
  # left this lane reporting `unformatted=0`. On a CI checkout every file is
  # tracked and the omission is invisible; locally it means the gate is blind to
  # exactly the file the developer is writing. --exclude-standard keeps .gitignore
  # in force, so the .gomodcache property above survives.
  #
  # Filtered to paths that still exist. `--cached` reports the index, and a file
  # deleted from the worktree but not yet staged is still in it. Handing that path
  # to gofmt yields `GetFileAttributesEx <path>: The system cannot find the file
  # specified.` on stderr, which 2>&1 folds into the output and Measure-GateGofmt
  # counts as an unformatted file -- a red lane blaming formatting for a perfectly
  # legitimate uncommitted deletion. Measured 2026-08-08: exit 2, one stderr line
  # per missing path.
  $listed = @(& git -C $repositoryRoot ls-files --cached --others --exclude-standard '*.go')
  $discovered = @($listed | Where-Object { Test-Path -LiteralPath (Join-Path $repositoryRoot $_) -PathType Leaf })
  Write-Host "listed=$($listed.Count) discovered=$($discovered.Count)"
  if ($discovered.Count -eq 0) {
    return New-GateVerdict -Lane 'gofmt' -Passed $false -Counts 'discovered=0' `
      -Reason 'no Go files reached the formatter; the filter is broken, not the tree'
  }

  # Batched because the argument list is a few hundred paths and Windows caps the
  # command line at 32 KB.
  #
  # Resolved against the repository root, not left relative: `git ls-files` prints
  # paths relative to that root, while `gofmt` resolves them against the process's
  # working directory. Those agree only when the gate is launched from the root,
  # which is what `npm run gate` and the CI step both happen to do -- so the
  # disagreement would surface as a pile of "file does not exist" lines counted as
  # violations the first time someone ran it from anywhere else.
  $output = New-Object System.Text.StringBuilder
  $gofmtExitFailures = 0
  for ($index = 0; $index -lt $discovered.Count; $index += 200) {
    $batch = @($discovered[$index..([Math]::Min($index + 199, $discovered.Count - 1))] |
      ForEach-Object { Join-Path $repositoryRoot $_ })
    $result = & gofmt -l @batch 2>&1
    if ($LASTEXITCODE -ne 0) { $gofmtExitFailures++ }
    foreach ($line in @($result)) { [void]$output.AppendLine([string]$line) }
  }
  $measurement = Measure-GateGofmt -Text $output.ToString()
  $counts = "run=$($discovered.Count) unformatted=$($measurement.Unformatted) tool_failures=$gofmtExitFailures"
  Write-Host $counts
  # Checked before the violation count, because a batch that errored makes the
  # output untrustworthy in both directions: `gofmt -l` exits 0 while listing
  # unformatted files, so a non-zero exit is never a finding about the code -- it
  # is the tool saying it did not do the job it was asked to do.
  if ($gofmtExitFailures -gt 0) {
    $measurement.Paths | ForEach-Object { Write-Host $_ }
    return New-GateVerdict -Lane 'gofmt' -Passed $false -Counts $counts `
      -Reason "gofmt exited non-zero on $gofmtExitFailures batch(es); its output above is not a formatting verdict"
  }
  if ($measurement.Unformatted -ne 0) {
    $measurement.Paths | ForEach-Object { Write-Host $_ }
    return New-GateVerdict -Lane 'gofmt' -Passed $false -Counts $counts `
      -Reason 'gofmt -l is non-empty. Run `gofmt -w` over the paths listed above.'
  }
  return New-GateVerdict -Lane 'gofmt' -Passed $true -Counts $counts
}

function Invoke-GateBuild {
  # HARNESS-PROFILE §11 defines the deterministic L0 lane as build AND vet. The
  # test lane is not a substitute: `go test` without an explicit -vet flag runs a
  # curated subset of vet's checks (`go help testflag`), and vet's checks over
  # non-test code are outside it entirely.
  #
  # `go list ./...` first: build and vet report only exit codes, so without an
  # independent package count this lane cannot tell "everything compiled" from
  # "the pattern matched nothing" (issue #3, mechanism (b)).
  $list = Invoke-GateTool -Name 'go-list' -FilePath 'go' -ArgumentList @('list', './...') -WorkingDirectory $serverCore -Quiet
  if ($list.ExitCode -ne 0) {
    return New-GateVerdict -Lane 'build' -Passed $false -Counts 'packages=unknown' -Reason "go list ./... exited $($list.ExitCode)"
  }
  # Count only lines shaped like an import path: Invoke-GateTool folds stderr in
  # (see its header), and `go list` announces an empty match with a warning on
  # stderr -- counting that line would make the zero-package floor unreachable.
  $packages = @($list.Text -split "`r?`n" | Where-Object { $_ -match '^\S+$' -and $_ -notmatch '^go:' }).Count
  if ($packages -eq 0) {
    return New-GateVerdict -Lane 'build' -Passed $false -Counts 'packages=0' `
      -Reason 'go list ./... resolved zero packages; the universe is empty, not clean'
  }
  $build = Invoke-GateTool -Name 'go-build' -FilePath 'go' -ArgumentList @('build', './...') -WorkingDirectory $serverCore
  if ($build.ExitCode -ne 0) {
    return New-GateVerdict -Lane 'build' -Passed $false -Counts "packages=$packages build=fail" -Reason "go build ./... exited $($build.ExitCode)"
  }
  $vet = Invoke-GateTool -Name 'go-vet' -FilePath 'go' -ArgumentList @('vet', './...') -WorkingDirectory $serverCore
  if ($vet.ExitCode -ne 0) {
    return New-GateVerdict -Lane 'build' -Passed $false -Counts "packages=$packages build=ok vet=fail" -Reason "go vet ./... exited $($vet.ExitCode)"
  }
  Write-Host "packages=$packages build=ok vet=ok"
  return New-GateVerdict -Lane 'build' -Passed $true -Counts "packages=$packages build=ok vet=ok"
}

# Pinned, not `@latest`. A gate whose tool version floats reports a different
# number on a different day, and the ratchet would attribute the change to
# whoever happened to push that morning. Bump this deliberately, in a PR that
# also re-measures the baseline.
$script:GolangciLintVersion = 'v2.12.2'

# The linters .golangci.yml asks for, plus `typecheck`, which golangci-lint
# always runs. Asserted against what the tool reports as enabled: a config file
# that failed to load leaves the tool running its own defaults and reporting a
# healthy small number over the wrong rules.
$script:GolangciLintExpected = @(
  'bodyclose', 'errcheck', 'errorlint', 'exhaustive', 'forbidigo', 'ineffassign',
  'noctx', 'rowserrcheck', 'sqlclosecheck', 'staticcheck', 'typecheck', 'unused'
) | Sort-Object

# The same assertion for the TypeScript side, and it arrived one review later
# than the Go one. Reported by CodeRabbit on PR #21 and reproduced 2026-08-08:
# setting `no-misused-promises` and `rules-of-hooks` to "off" in
# `eslint.config.mjs` left this lane PASSing with `total=13 baseline=26`, and it
# printed `shrink: ... -- commit the lower baseline in this PR`. Turning a rule
# off is reported as progress by any instrument that only counts findings, and
# the ratchet then locks the lower number in.
#
# GATE-TOPOLOGY §2.3a names these six. The list is the contract; the baseline's
# by_rule keys are checked against it too, so the two cannot drift apart.
$script:EslintExpected = @(
  '@typescript-eslint/await-thenable',
  '@typescript-eslint/no-floating-promises',
  '@typescript-eslint/no-misused-promises',
  '@typescript-eslint/no-unused-vars',
  'react-hooks/exhaustive-deps',
  'react-hooks/rules-of-hooks'
) | Sort-Object

function Invoke-GateGovernance {
  # Two halves of scripts/harness/Policy.psm1, called directly rather than
  # through scripts/harness.ps1 so the lane reads structured violations instead
  # of re-parsing its own text output.
  #
  # Validate -- the six registries against their schemas plus every
  # cross-reference -- blocks at zero: it passes clean today (measured
  # 2026-08-08) and a governance contract that stops parsing has no legitimate
  # transitional state. Drift -- the semantic walk of the tree against those
  # contracts -- is a ratchet at 58, because its findings are the legacy tree's
  # unregistered readers and undeclared dependencies, and paying them down is
  # behavior-review work that an instrument PR must not do.
  $baselinePath = Join-Path $repositoryRoot 'contracts/gate/baselines.json'
  if (-not (Test-Path -LiteralPath $baselinePath -PathType Leaf)) {
    return New-GateVerdict -Lane 'governance' -Passed $false -Counts 'baseline=missing' `
      -Reason "no ratchet baseline at $baselinePath."
  }
  $baselineDocument = Get-Content -LiteralPath $baselinePath -Raw | ConvertFrom-Json
  $baseline = ConvertTo-GateCountMap -Object $baselineDocument.'governance-drift'.by_code

  Import-Module (Join-Path $PSScriptRoot 'harness/Policy.psm1') -Force

  $validate = Test-GovernanceContracts -RepositoryRoot $repositoryRoot
  $registries = @($validate.Documents.Keys).Count
  if (-not $validate.Passed) {
    foreach ($violation in @($validate.Violations)) { Write-Host "  $($violation.ErrorCode) $($violation.Id) $($violation.Path)" }
    return New-GateVerdict -Lane 'governance' -Passed $false -Counts "registries=$registries validate=$(@($validate.Violations).Count)" `
      -Reason 'the governance contracts do not validate. This half of the lane blocks at zero: a registry that stops parsing has no transitional state.'
  }
  if ($registries -lt 6) {
    return New-GateVerdict -Lane 'governance' -Passed $false -Counts "registries=$registries" `
      -Reason 'fewer than the six governance registries loaded. A validate pass over a partial set is not a pass.'
  }

  # The base SHA only feeds the API<->SDK atomicity check, which is a statement
  # about a change and therefore needs a base to diff against. CI provides the
  # PR base as GATE_BASE_SHA; locally the merge base against origin/main is the
  # same fact. HEAD is the honest fallback -- an empty diff, so the atomicity
  # check is silent -- and the counts line says which one was used, because a
  # check that silently did not run is this gate's founding defect class.
  $baseSha = $env:GATE_BASE_SHA
  $baseSource = 'env'
  if ([string]::IsNullOrWhiteSpace($baseSha)) {
    $baseSha = (& git -C $repositoryRoot merge-base HEAD origin/main 2>$null)
    $baseSource = 'merge-base'
  }
  if ([string]::IsNullOrWhiteSpace($baseSha) -or $baseSha -notmatch '^[0-9a-f]{40}$') {
    $baseSha = (& git -C $repositoryRoot rev-parse HEAD)
    $baseSource = 'head'
  }

  $drift = Test-GovernanceDrift -RepositoryRoot $repositoryRoot -BaseSha $baseSha
  $measurement = Measure-GateGovernance -Violations @($drift.Violations)
  $counts = "registries=$registries files=$($drift.FilesScanned) total=$($measurement.Total) baseline=$($baselineDocument.'governance-drift'.total) baselined_exceptions=$(@($drift.BaselineExceptions).Count) base=$baseSource"
  Write-Host $counts
  foreach ($code in @($measurement.ByCode.Keys | Sort-Object)) {
    Write-Host ("  {0,-32} {1,4}  baseline {2,4}" -f $code, $measurement.ByCode[$code], $(if ($baseline.ContainsKey($code)) { $baseline[$code] } else { 'n/a' }))
  }

  if ($drift.FilesScanned -le 0) {
    return New-GateVerdict -Lane 'governance' -Passed $false -Counts $counts `
      -Reason "the drift walk scanned $($drift.FilesScanned) file(s). A verdict over the empty set is not a verdict."
  }

  $measuredMap = @{}
  foreach ($key in @($baseline.Keys)) { $measuredMap[$key] = 0 }
  foreach ($key in @($measurement.ByCode.Keys)) { $measuredMap[$key] = $measurement.ByCode[$key] }
  $comparison = Compare-GateRatchet -Measured $measuredMap -Baseline $baseline
  foreach ($line in @($comparison.Decreased)) { Write-Host "shrink: $line -- commit the lower baseline in this PR" }
  if (-not $comparison.Passed) {
    foreach ($line in @($comparison.Increased)) { Write-Host "increase: $line" }
    foreach ($line in @($comparison.Unknown)) { Write-Host "unknown error code, absent from the baseline: $line" }
    foreach ($violation in @($drift.Violations)) { Write-Host "  $($violation.ErrorCode) $($violation.Id) $($violation.Path)" }
    return New-GateVerdict -Lane 'governance' -Passed $false -Counts $counts `
      -Reason 'governance drift findings increased over the committed baseline, or an error code appeared that the baseline does not name.'
  }
  return New-GateVerdict -Lane 'governance' -Passed $true -Counts $counts
}

function Invoke-GateBoundary {
  # TestModuleBoundaryADR023, the detector for ADR-023 §2: a module is
  # importable by another module only at X/ports. It is red on main by design --
  # 234 cross-module imports in the legacy tree, owned by Onda 1 of
  # docs/superpowers/plans/2026-08-05-arquitetura-protocolo-de-modulo-plan.md --
  # so the test-go lane skips it and THIS lane runs it as a ratchet: the debt
  # may shrink and may not grow. When Onda 1 task 1.5 lands and the count
  # reaches zero, this lane keeps the test enforcing at zero; delete the skip in
  # Invoke-GateGoTest then, not this lane.
  $baselinePath = Join-Path $repositoryRoot 'contracts/gate/baselines.json'
  if (-not (Test-Path -LiteralPath $baselinePath -PathType Leaf)) {
    return New-GateVerdict -Lane 'boundary' -Passed $false -Counts 'baseline=missing' `
      -Reason "no ratchet baseline at $baselinePath."
  }
  $baselineDocument = Get-Content -LiteralPath $baselinePath -Raw | ConvertFrom-Json
  $baseline = ConvertTo-GateCountMap -Object $baselineDocument.boundary.by_origin

  # -v so the boundary_files log line prints on a passing run too; -count=1
  # because a cached verdict is not a verdict about this commit. -Quiet: the
  # violation list is 234 lines of legacy sites; the lane prints the per-layer
  # counts and the log file keeps the sites.
  $result = Invoke-GateTool -Name 'boundary' -FilePath 'go' -Quiet `
    -ArgumentList @('test', './internal/composition/', '-run', 'TestModuleBoundaryADR023', '-count=1', '-v') `
    -WorkingDirectory $serverCore
  $measurement = Measure-GateBoundary -Text $result.Text

  $counts = "files=$($measurement.Files) total=$($measurement.Total) baseline=$($baselineDocument.boundary.total)"
  Write-Host $counts
  foreach ($origin in @($measurement.ByOrigin.Keys | Sort-Object)) {
    Write-Host ("  {0,-28} {1,4}  baseline {2,4}" -f $origin, $measurement.ByOrigin[$origin], $(if ($baseline.ContainsKey($origin)) { $baseline[$origin] } else { 'n/a' }))
  }

  if ($measurement.Files -lt 0) {
    return New-GateVerdict -Lane 'boundary' -Passed $false -Counts $counts `
      -Reason "no boundary_files line in the test output (exit $($result.ExitCode)). The run died before the walk finished, or the count was deleted from the test; see $($result.LogPath)."
  }
  if ($measurement.Files -eq 0) {
    return New-GateVerdict -Lane 'boundary' -Passed $false -Counts $counts `
      -Reason 'the boundary walk parsed zero Go files. A verdict over the empty set is not a verdict.'
  }
  if ($measurement.Total -lt 0) {
    return New-GateVerdict -Lane 'boundary' -Passed $false -Counts $counts `
      -Reason "the test neither passed nor reported a violation count (exit $($result.ExitCode))."
  }
  # Exit 1 is a red test, the normal state under a ratchet over a
  # red-by-design detector. Exit 0 is the debt paid off. Anything else is the
  # toolchain failing, not a boundary verdict.
  if ($result.ExitCode -notin @(0, 1)) {
    return New-GateVerdict -Lane 'boundary' -Passed $false -Counts $counts `
      -Reason "go test exited $($result.ExitCode), which is neither a pass nor a red test."
  }

  # The ratchet compares ByOrigin, not Total. If the origin block goes missing
  # or half-parses while violations still exist, the measured map degrades to
  # zeros and the ratchet would pass on a measurement it never actually took.
  $originTotal = [int](@($measurement.ByOrigin.Values | Measure-Object -Sum).Sum)
  if ($originTotal -ne $measurement.Total) {
    return New-GateVerdict -Lane 'boundary' -Passed $false -Counts $counts `
      -Reason "the by-origin counts sum to $originTotal but the detector reported $($measurement.Total) violation(s); the origin block was missing or only partially parsed, so there is no measurement to ratchet."
  }

  $measuredMap = @{}
  foreach ($key in @($baseline.Keys)) { $measuredMap[$key] = 0 }
  foreach ($key in @($measurement.ByOrigin.Keys)) { $measuredMap[$key] = $measurement.ByOrigin[$key] }
  $comparison = Compare-GateRatchet -Measured $measuredMap -Baseline $baseline
  foreach ($line in @($comparison.Decreased)) { Write-Host "shrink: $line -- commit the lower baseline in this PR" }
  if (-not $comparison.Passed) {
    foreach ($line in @($comparison.Increased)) { Write-Host "increase: $line" }
    foreach ($line in @($comparison.Unknown)) { Write-Host "unknown origin layer, absent from the baseline: $line" }
    return New-GateVerdict -Lane 'boundary' -Passed $false -Counts $counts `
      -Reason "module-boundary violations increased over the committed baseline, or a new origin layer appeared. The sites are in $($result.LogPath)."
  }
  return New-GateVerdict -Lane 'boundary' -Passed $true -Counts $counts
}

function Invoke-GateArch {
  # The architecture detectors: rules the Go compiler cannot express. The
  # scanner had zero invokers from the day it was written until 2026-08-08, and
  # the absence was rational rather than forgotten -- substring matching and a
  # guessed exemption list produced 483 findings of which 1 was real, and an
  # instrument that cries 483 times for 1 defect costs more to read than it
  # returns. The scanner was fixed first (segment matching, declared
  # vendor-aware layers, a scanned-file count) and wired here after, at 6.
  $baselinePath = Join-Path $repositoryRoot 'contracts/gate/baselines.json'
  if (-not (Test-Path -LiteralPath $baselinePath -PathType Leaf)) {
    return New-GateVerdict -Lane 'arch' -Passed $false -Counts 'baseline=missing' `
      -Reason "no ratchet baseline at $baselinePath."
  }
  $baselineDocument = Get-Content -LiteralPath $baselinePath -Raw | ConvertFrom-Json
  $baseline = ConvertTo-GateCountMap -Object $baselineDocument.archscan.by_rule

  $result = Invoke-GateTool -Name 'archscan' -FilePath 'go' `
    -ArgumentList @('run', './internal/arch/cmd/archscan') `
    -WorkingDirectory $serverCore
  $measurement = Measure-GateArchscan -Text $result.Text

  $counts = "scanned=$($measurement.Scanned) total=$($measurement.Total) baseline=$($baselineDocument.archscan.total)"
  Write-Host $counts
  foreach ($rule in @($measurement.ByRule.Keys | Sort-Object)) {
    Write-Host ("  {0,-42} {1,4}  baseline {2,4}" -f $rule, $measurement.ByRule[$rule], $(if ($baseline.ContainsKey($rule)) { $baseline[$rule] } else { 'n/a' }))
  }

  if ($measurement.Scanned -lt 0) {
    return New-GateVerdict -Lane 'arch' -Passed $false -Counts $counts `
      -Reason "archscan printed no summary line (exit $($result.ExitCode)). A stream without the count is a run that did not finish."
  }
  if ($measurement.Scanned -eq 0) {
    return New-GateVerdict -Lane 'arch' -Passed $false -Counts $counts `
      -Reason 'archscan parsed zero Go files. A verdict over the empty set is not a verdict.'
  }
  # Exit 1 is findings, the normal state under a ratchet. Exit 2 is the scanner
  # failing to run, and anything else is unknown.
  if ($result.ExitCode -notin @(0, 1)) {
    return New-GateVerdict -Lane 'arch' -Passed $false -Counts $counts `
      -Reason "archscan exited $($result.ExitCode), which is neither clean nor findings."
  }

  $measuredMap = @{}
  foreach ($key in @($baseline.Keys)) { $measuredMap[$key] = 0 }
  foreach ($key in @($measurement.ByRule.Keys)) { $measuredMap[$key] = $measurement.ByRule[$key] }
  $comparison = Compare-GateRatchet -Measured $measuredMap -Baseline $baseline
  foreach ($line in @($comparison.Decreased)) { Write-Host "shrink: $line -- commit the lower baseline in this PR" }
  if (-not $comparison.Passed) {
    foreach ($line in @($comparison.Increased)) { Write-Host "increase: $line" }
    foreach ($line in @($comparison.Unknown)) { Write-Host "unknown rule, absent from the baseline: $line" }
    return New-GateVerdict -Lane 'arch' -Passed $false -Counts $counts `
      -Reason 'archscan findings increased over the committed baseline, or a rule appeared that the baseline does not name.'
  }
  return New-GateVerdict -Lane 'arch' -Passed $true -Counts $counts
}

function Invoke-GateLintGo {
  $baselinePath = Join-Path $repositoryRoot 'contracts/gate/baselines.json'
  if (-not (Test-Path -LiteralPath $baselinePath -PathType Leaf)) {
    return New-GateVerdict -Lane 'lint-go' -Passed $false -Counts 'baseline=missing' `
      -Reason "no ratchet baseline at $baselinePath. A lint gate with no baseline either blocks on everything or on nothing."
  }
  $baselineDocument = Get-Content -LiteralPath $baselinePath -Raw | ConvertFrom-Json
  $baseline = ConvertTo-GateCountMap -Object $baselineDocument.'golangci-lint'.by_linter

  # The analysed file set is a function of the build tags in force, so the count
  # is a function of the platform. Measured 2026-08-08: the Linux runner reports
  # errcheck=30 and this Windows host reports 28, because gcc is present there
  # and absent here, so `//go:build cgo` selects `oracle_cgo.go` and
  # `open_cgo.go` on one and their `!cgo` twins on the other.
  #
  # The baseline is therefore pinned to the enforcing platform. A run elsewhere
  # still blocks on an increase -- an increase is real everywhere -- but its
  # shrinks are not evidence of anything, and must not be printed as though they
  # were. Fewer findings because fewer files compiled is the same lie as a green
  # suite that ran nothing, one level up.
  $goos = (& go env GOOS 2>$null)
  $cgo = (& go env CGO_ENABLED 2>$null)
  $platform = "goos=$goos cgo=$cgo"
  $baselinePlatform = "goos=$($baselineDocument.'golangci-lint'.measured_on_goos) cgo=$($baselineDocument.'golangci-lint'.measured_on_cgo)"
  $platformMatches = ($platform -eq $baselinePlatform)
  Write-Host "$platform baseline_platform=$baselinePlatform comparable=$platformMatches"

  # Same universe guard as the build lane: the ratchet compares totals, and a
  # total taken over zero analyzed packages would read as the ratchet's best day
  # ever (issue #3, mechanism (b)).
  $list = Invoke-GateTool -Name 'golangci-go-list' -FilePath 'go' -ArgumentList @('list', './...') -WorkingDirectory $serverCore -Quiet
  if ($list.ExitCode -ne 0) {
    return New-GateVerdict -Lane 'lint-go' -Passed $false -Counts 'packages=unknown' -Reason "go list ./... exited $($list.ExitCode)"
  }
  # Count only lines shaped like an import path: Invoke-GateTool folds stderr in
  # (see its header), and `go list` announces an empty match with a warning on
  # stderr -- counting that line would make the zero-package floor unreachable.
  $packages = @($list.Text -split "`r?`n" | Where-Object { $_ -match '^\S+$' -and $_ -notmatch '^go:' }).Count
  if ($packages -eq 0) {
    return New-GateVerdict -Lane 'lint-go' -Passed $false -Counts 'packages=0' `
      -Reason 'go list ./... resolved zero packages; a lint verdict over the empty set is not a verdict'
  }

  $reportPath = Join-Path $logDirectory 'golangci.json'
  if (Test-Path -LiteralPath $reportPath) { Remove-Item -LiteralPath $reportPath -Force }
  # `go run <module>@<version>` rather than an installed binary or a CI-only
  # action: the same invocation has to work on a developer machine and on a
  # runner, and it leaves go.mod and go.sum untouched (verified 2026-08-08).
  $result = Invoke-GateTool -Name 'golangci-lint' -FilePath 'go' `
    -ArgumentList @(
      'run', "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$($script:GolangciLintVersion)",
      'run', "--output.json.path=$reportPath", './...') `
    -WorkingDirectory $serverCore

  if (-not (Test-Path -LiteralPath $reportPath -PathType Leaf)) {
    return New-GateVerdict -Lane 'lint-go' -Passed $false -Counts "packages=$packages report=missing" `
      -Reason "golangci-lint wrote no report (exit $($result.ExitCode)). The lane cannot report zero findings for a run that did not happen."
  }
  $measurement = Measure-GateGolangciLint -Json (Get-Content -LiteralPath $reportPath -Raw)

  # Exit 1 is what golangci-lint returns when it found issues, which is the
  # normal state under a ratchet. Anything else is the tool failing to run, and
  # that must not read as a clean lane.
  if ($result.ExitCode -notin @(0, 1)) {
    return New-GateVerdict -Lane 'lint-go' -Passed $false -Counts "packages=$packages total=$($measurement.Total)" `
      -Reason "golangci-lint exited $($result.ExitCode), which is neither clean nor findings."
  }

  $enabledDrift = @(Compare-Object -ReferenceObject $script:GolangciLintExpected -DifferenceObject $measurement.Enabled)
  if ($enabledDrift.Count -ne 0) {
    $detail = ($enabledDrift | ForEach-Object { "$($_.SideIndicator) $($_.InputObject)" }) -join ', '
    return New-GateVerdict -Lane 'lint-go' -Passed $false -Counts "packages=$packages enabled=$($measurement.Enabled.Count)" `
      -Reason "the enabled linter set is not the configured one ($detail). A count over the wrong rules is not a smaller count."
  }

  $measuredMap = @{}
  foreach ($key in @($baseline.Keys)) { $measuredMap[$key] = 0 }
  foreach ($key in @($measurement.ByLinter.Keys)) { $measuredMap[$key] = $measurement.ByLinter[$key] }
  $comparison = Compare-GateRatchet -Measured $measuredMap -Baseline $baseline

  $counts = "packages=$packages total=$($measurement.Total) baseline=$($baselineDocument.'golangci-lint'.total) enabled=$($measurement.Enabled.Count)"
  Write-Host $counts
  foreach ($line in @($measurement.ByLinter.Keys | Sort-Object)) {
    Write-Host ("  {0,-14} {1,4}  baseline {2,4}" -f $line, $measurement.ByLinter[$line], $(if ($baseline.ContainsKey($line)) { $baseline[$line] } else { 'n/a' }))
  }
  foreach ($line in @($comparison.Decreased)) {
    if ($platformMatches) {
      Write-Host "shrink: $line -- commit the lower baseline in this PR"
    } else {
      Write-Host "lower here, NOT a shrink: $line -- measured on $platform, baseline on $baselinePlatform. Different build tags select different files; do not commit this number."
    }
  }
  if (-not $comparison.Passed) {
    foreach ($line in @($comparison.Increased)) { Write-Host "increase: $line" }
    foreach ($line in @($comparison.Unknown)) { Write-Host "unknown linter, absent from the baseline: $line" }
    return New-GateVerdict -Lane 'lint-go' -Passed $false -Counts $counts `
      -Reason 'golangci-lint findings increased over the committed baseline, or a linter appeared that the baseline does not name.'
  }
  return New-GateVerdict -Lane 'lint-go' -Passed $true -Counts $counts
}

function Get-GateWebSources {
  <#
    The .ts/.tsx files this repository owns, by the same rule the gofmt lane uses:
    `git ls-files` rather than a tree walk, tracked plus untracked-not-ignored, so
    the file the developer is writing right now is in scope and `.gomodcache` is
    not.

    Generated output and tsc's `.d.ts` emit are dropped here because they are
    dropped in `eslint.config.mjs` too; the two lists have to agree or the
    coverage assertion below reports drift that is really disagreement between
    two ignore lists.
  #>
  $listed = @(& git -C $repositoryRoot ls-files --cached --others --exclude-standard '*.ts' '*.tsx')
  return @($listed |
    Where-Object { $_ -notmatch 'generated/' -and $_ -notmatch '\.d\.ts$' } |
    Where-Object { Test-Path -LiteralPath (Join-Path $repositoryRoot $_) -PathType Leaf } |
    Sort-Object)
}

function Invoke-GateFormat {
  # Prettier, blocking at zero and never ratcheted. GATE-TOPOLOGY §2.3a: formatting
  # admits no legitimate exception, so a ratchet would only record how long the
  # drift had been tolerated.
  #
  # Markdown is deliberately out of scope and `.prettierignore` carries the
  # measurement -- Prettier's markdown printer damages nested lists whose
  # continuation line sits inside a wrapped code span, in 77 of 236 files here.
  #
  # --log-level debug, and the reason is the whole anti-vacuity guard for this
  # lane. At its default verbosity Prettier names no file on a clean run, so its
  # output cannot tell a formatted tree from a glob that matched nothing: it
  # prints the same clean marker either way. At debug it emits one
  # `resolve config from '<path>'` line per file it actually processed, which is
  # the effective file set and the only number here that a widened
  # `.prettierignore` can move. The stream is ~3 lines per file, so it goes to
  # the log and not to the console.
  $result = Invoke-GateTool -Name 'prettier' -FilePath 'npx' -Quiet `
    -ArgumentList @('--no-install', 'prettier', '--check', '--log-level', 'debug', '**/*.{ts,tsx,mjs,json,yml,yaml}')
  $measurement = Measure-GatePrettier -Text $result.Text

  # Two counts, and they answer different questions. `candidates` comes from git
  # and says how many files of these extensions exist; `checked` comes from
  # Prettier and says how many it read. The first was the lane's only guard until
  # CodeRabbit reported on PR #21 that it does not constrain the second.
  # Reproduced 2026-08-08: appending an ignore pattern for every extension left
  # this lane PASSing with `candidates=342 unformatted=0 clean_marker=True` while
  # Prettier checked zero files. Measured on the committed tree, `checked` is 264
  # and drops to 0 the moment the ignore list swallows the glob.
  $candidates = @(& git -C $repositoryRoot ls-files --cached --others --exclude-standard `
      '*.ts' '*.tsx' '*.mjs' '*.json' '*.yml' '*.yaml').Count
  $counts = "candidates=$candidates checked=$($measurement.Checked) unformatted=$($measurement.Unformatted) clean_marker=$($measurement.Clean)"
  Write-Host $counts
  foreach ($path in @($measurement.Paths)) { Write-Host "  unformatted: $path" }

  if ($candidates -eq 0) {
    return New-GateVerdict -Lane 'format' -Passed $false -Counts $counts `
      -Reason 'git reports no file of any formatted extension. The lane would pass over an empty set.'
  }
  if ($measurement.Checked -eq 0) {
    return New-GateVerdict -Lane 'format' -Passed $false -Counts $counts `
      -Reason "prettier read no file at all, over $candidates candidate(s). Its clean marker is about the empty set. Check .prettierignore and the glob."
  }
  # A floor rather than equality: `candidates` deliberately overcounts, because
  # .prettierignore excludes the lockfile, generated output and recorded evidence
  # on purpose and each exclusion is argued in that file. What must never happen
  # again is the ignore list quietly eating most of the tree, so the floor is set
  # where an accident is loud and the standing exclusions are not.
  $floor = [math]::Ceiling($candidates * 0.5)
  if ($measurement.Checked -lt $floor) {
    return New-GateVerdict -Lane 'format' -Passed $false -Counts $counts `
      -Reason "prettier read $($measurement.Checked) of $candidates candidate file(s), below the floor of $floor. An ignore rule is excluding more than .prettierignore accounts for."
  }
  if ($result.ExitCode -ne 0 -or $measurement.Unformatted -gt 0) {
    return New-GateVerdict -Lane 'format' -Passed $false -Counts $counts `
      -Reason "prettier exited $($result.ExitCode) with $($measurement.Unformatted) unformatted file(s). Run: npx prettier --write `"**/*.{ts,tsx,json,yml,yaml}`""
  }
  if (-not $measurement.Clean) {
    return New-GateVerdict -Lane 'format' -Passed $false -Counts $counts `
      -Reason 'prettier exited 0 without printing its clean marker. Silence is not a verdict.'
  }
  return New-GateVerdict -Lane 'format' -Passed $true -Counts $counts
}

function Invoke-GateLintWeb {
  $baselinePath = Join-Path $repositoryRoot 'contracts/gate/baselines.json'
  if (-not (Test-Path -LiteralPath $baselinePath -PathType Leaf)) {
    return New-GateVerdict -Lane 'lint-web' -Passed $false -Counts 'baseline=missing' `
      -Reason "no ratchet baseline at $baselinePath."
  }
  $baselineDocument = Get-Content -LiteralPath $baselinePath -Raw | ConvertFrom-Json
  $baseline = ConvertTo-GateCountMap -Object $baselineDocument.eslint.by_rule

  # The baseline names the rules it counted. If that set and $script:EslintExpected
  # disagree, one of the two moved alone, and every number below is over a rule set
  # nobody declared.
  $baselineDrift = @(Compare-Object -ReferenceObject $script:EslintExpected -DifferenceObject @($baseline.Keys | Sort-Object))
  if ($baselineDrift.Count -ne 0) {
    $detail = ($baselineDrift | ForEach-Object { "$($_.SideIndicator) $($_.InputObject)" }) -join ', '
    return New-GateVerdict -Lane 'lint-web' -Passed $false -Counts "baseline_rules=$($baseline.Count)" `
      -Reason "the baseline's by_rule keys are not the configured rule set ($detail)."
  }

  # Rule set before findings, for the same reason coverage comes before findings
  # below: a rule that is off produces zero findings and looks like a rule that
  # found nothing. `--print-config` resolves the flat config the way ESLint itself
  # will for that file, so this reads the effective severity rather than the text
  # of eslint.config.mjs.
  $sources = Get-GateWebSources
  if ($sources.Count -eq 0) {
    return New-GateVerdict -Lane 'lint-web' -Passed $false -Counts 'sources=0' `
      -Reason 'git reports no .ts/.tsx source. The lane would lint an empty set.'
  }
  $probe = Invoke-GateTool -Name 'eslint-print-config' -FilePath 'npx' -Quiet `
    -ArgumentList @('--no-install', 'eslint', '--print-config', $sources[0])
  $enabled = Measure-GateEslintRules -Json $probe.Text
  $enabledDrift = @(Compare-Object -ReferenceObject $script:EslintExpected -DifferenceObject @($enabled))
  if ($probe.ExitCode -ne 0 -or $enabledDrift.Count -ne 0) {
    $detail = ($enabledDrift | ForEach-Object { "$($_.SideIndicator) $($_.InputObject)" }) -join ', '
    return New-GateVerdict -Lane 'lint-web' -Passed $false -Counts "enabled=$(@($enabled).Count) expected=$($script:EslintExpected.Count) probe_exit=$($probe.ExitCode)" `
      -Reason "the enabled rule set for $($sources[0]) is not the configured one ($detail). A rule that is off reports zero findings, and the ratchet would take the drop for a shrink."
  }

  $reportPath = Join-Path $logDirectory 'eslint.json'
  if (Test-Path -LiteralPath $reportPath) { Remove-Item -LiteralPath $reportPath -Force }
  $result = Invoke-GateTool -Name 'eslint' -FilePath 'npx' `
    -ArgumentList @('--no-install', 'eslint', '.', '--format', 'json', '--output-file', $reportPath)

  if (-not (Test-Path -LiteralPath $reportPath -PathType Leaf)) {
    return New-GateVerdict -Lane 'lint-web' -Passed $false -Counts 'report=missing' `
      -Reason "eslint wrote no report (exit $($result.ExitCode)). The lane cannot report zero findings for a run that did not happen."
  }
  $measurement = Measure-GateEslint -Json (Get-Content -LiteralPath $reportPath -Raw)

  # Coverage before findings, and this order is the lesson of the lane. The
  # type-aware rules -- no-floating-promises, no-misused-promises, await-thenable
  # -- need a file to belong to a tsconfig. `tsconfig.eslint.json` exists to cover
  # all of them, and a file that slips out of it is still reported by ESLint,
  # still contributes zero findings, and still looks exactly like a clean file.
  # So the tracked source list is the reference and every entry must appear.
  $expected = $sources
  $linted = @{}
  $prefix = ($repositoryRoot -replace '\\', '/').TrimEnd('/') + '/'
  foreach ($path in @($measurement.Paths)) { $linted[($path -replace [regex]::Escape($prefix), '')] = $true }
  $uncovered = @($expected | Where-Object { -not $linted.ContainsKey($_) })

  $counts = "linted=$($measurement.Files) sources=$($expected.Count) uncovered=$($uncovered.Count) enabled=$(@($enabled).Count) total=$($measurement.Total) baseline=$($baselineDocument.eslint.total) fatal=$($measurement.Fatal)"
  Write-Host $counts
  foreach ($rule in @($measurement.ByRule.Keys | Sort-Object)) {
    Write-Host ("  {0,-42} {1,4}  baseline {2,4}" -f $rule, $measurement.ByRule[$rule], $(if ($baseline.ContainsKey($rule)) { $baseline[$rule] } else { 'n/a' }))
  }

  if ($uncovered.Count -gt 0) {
    foreach ($path in @($uncovered | Select-Object -First 20)) { Write-Host "  uncovered: $path" }
    return New-GateVerdict -Lane 'lint-web' -Passed $false -Counts $counts `
      -Reason 'ESLint did not lint every tracked .ts/.tsx source. A file outside the lint project is reported as clean by rules that could not run on it.'
  }
  if ($measurement.Fatal -gt 0) {
    return New-GateVerdict -Lane 'lint-web' -Passed $false -Counts $counts `
      -Reason 'ESLint could not parse one or more files. A file that failed to parse was not linted, whatever the findings count says.'
  }
  # Exit 1 is findings, which is the normal state under a ratchet. Anything else
  # is the tool failing to run.
  if ($result.ExitCode -notin @(0, 1)) {
    return New-GateVerdict -Lane 'lint-web' -Passed $false -Counts $counts `
      -Reason "eslint exited $($result.ExitCode), which is neither clean nor findings."
  }

  $measuredMap = @{}
  foreach ($key in @($baseline.Keys)) { $measuredMap[$key] = 0 }
  foreach ($key in @($measurement.ByRule.Keys)) { $measuredMap[$key] = $measurement.ByRule[$key] }
  $comparison = Compare-GateRatchet -Measured $measuredMap -Baseline $baseline
  foreach ($line in @($comparison.Decreased)) { Write-Host "shrink: $line -- commit the lower baseline in this PR" }
  if (-not $comparison.Passed) {
    foreach ($line in @($comparison.Increased)) { Write-Host "increase: $line" }
    foreach ($line in @($comparison.Unknown)) { Write-Host "unknown rule, absent from the baseline: $line" }
    return New-GateVerdict -Lane 'lint-web' -Passed $false -Counts $counts `
      -Reason 'ESLint findings increased over the committed baseline, or a rule fired that the baseline does not name.'
  }
  return New-GateVerdict -Lane 'lint-web' -Passed $true -Counts $counts
}

function Invoke-GateGoTest {
  # ./... and not ./tests/unit/...: the latter is 14 test files, internal/ holds
  # 353. The integration lane is invisible here regardless -- it sits behind
  # `//go:build integration`, which ./... does not compile.
  #
  # -skip TestModuleBoundaryADR023: that test is red on main by design (234
  # cross-module import violations) and the `boundary` lane runs it as a
  # ratchet, so the debt blocks on increase without turning this lane red.
  # Paying it down is Onda 1 of
  # docs/superpowers/plans/2026-08-05-arquitetura-protocolo-de-modulo-plan.md,
  # tasks 1.1 through 1.6. DELETE THIS SKIP when task 1.5 lands.
  #
  # -count=1 defeats the test cache: a cached PASS is not evidence that this
  # commit's tests ran.
  $result = Invoke-GateTool -Name 'go-test' -FilePath 'go' `
    -ArgumentList @('test', './...', '-skip', 'TestModuleBoundaryADR023', '-count=1', '-v') `
    -WorkingDirectory $serverCore
  $measurement = Measure-GateGoTest -Text $result.Text
  $counts = "run=$($measurement.Run) pass=$($measurement.Passed) skip=$($measurement.Skipped) fail=$($measurement.Failed)"
  Write-Host $counts
  if ($result.ExitCode -ne 0) {
    return New-GateVerdict -Lane 'test-go' -Passed $false -Counts $counts -Reason "go test exited $($result.ExitCode)"
  }
  if ($measurement.Run -eq 0 -or $measurement.Passed -eq 0) {
    return New-GateVerdict -Lane 'test-go' -Passed $false -Counts $counts `
      -Reason 'the lane exited 0 having executed nothing, or having passed nothing. That is the defect this line exists to catch.'
  }
  return New-GateVerdict -Lane 'test-go' -Passed $true -Counts $counts
}

function Invoke-GateTypecheck {
  # --listFiles is what makes this lane countable. `tsc --noEmit` prints nothing
  # on success, so a project that resolved zero source files -- a broken
  # `include`, a renamed directory -- exits 0 with the same silence as a clean
  # compile.
  #
  # tsconfig.eslint.json and not apps/web/tsconfig.json: the web project covers
  # 194 of the tracked sources and leaves the package test files dark, which is
  # how 12 type errors sat green in files vitest executes every run. The lint
  # project already exists as the one project declared over everything, so the
  # typecheck lane checks the same set, and the census below holds it there.
  $result = Invoke-GateTool -Name 'tsc' -FilePath 'npx' `
    -ArgumentList @('--no-install', 'tsc', '-p', 'tsconfig.eslint.json', '--noEmit', '--listFiles')
  $measurement = Measure-GateTsc -Text $result.Text

  # Coverage before errors, same order and same reason as lint-web: a tracked
  # source the project does not load is reported clean by a compiler that never
  # saw it. The tracked list is the reference and every entry must be loaded.
  $sources = Get-GateWebSources
  $loaded = @{}
  $prefix = ($repositoryRoot -replace '\\', '/').TrimEnd('/') + '/'
  foreach ($path in @($measurement.Paths)) {
    $loaded[(($path -replace '\\', '/') -replace [regex]::Escape($prefix), '')] = $true
  }
  $uncovered = @($sources | Where-Object { -not $loaded.ContainsKey($_) })

  $counts = "checked=$($measurement.Checked) sources=$($sources.Count) uncovered=$($uncovered.Count) errors=$($measurement.Errors)"
  Write-Host $counts
  if ($measurement.Checked -eq 0) {
    return New-GateVerdict -Lane 'typecheck' -Passed $false -Counts $counts `
      -Reason 'tsc loaded no project files. The project resolved nothing, which is not the same as compiling cleanly.'
  }
  if ($uncovered.Count -gt 0) {
    foreach ($path in @($uncovered | Select-Object -First 20)) { Write-Host "  uncovered: $path" }
    return New-GateVerdict -Lane 'typecheck' -Passed $false -Counts $counts `
      -Reason 'tsc did not load every tracked .ts/.tsx source. A file outside the project is reported clean by a compiler that never saw it.'
  }
  # Exit 2 is type errors, the normal state under a ratchet; anything else above
  # 0 is the tool failing to run.
  if ($result.ExitCode -notin @(0, 2)) {
    return New-GateVerdict -Lane 'typecheck' -Passed $false -Counts $counts -Reason "tsc exited $($result.ExitCode), which is neither clean nor type errors."
  }

  $baselinePath = Join-Path $repositoryRoot 'contracts/gate/baselines.json'
  if (-not (Test-Path -LiteralPath $baselinePath -PathType Leaf)) {
    return New-GateVerdict -Lane 'typecheck' -Passed $false -Counts 'baseline=missing' `
      -Reason "no ratchet baseline at $baselinePath."
  }
  $baselineDocument = Get-Content -LiteralPath $baselinePath -Raw | ConvertFrom-Json
  $baseline = ConvertTo-GateCountMap -Object $baselineDocument.tsc.by_file
  $comparison = Compare-GateRatchet -Measured $measurement.ByFile -Baseline $baseline
  foreach ($line in @($comparison.Decreased)) { Write-Host "shrink: $line -- commit the lower baseline in this PR" }
  if (-not $comparison.Passed) {
    foreach ($line in @($comparison.Increased)) { Write-Host "increase: $line" }
    foreach ($line in @($comparison.Unknown)) { Write-Host "unknown file, absent from the baseline: $line" }
    return New-GateVerdict -Lane 'typecheck' -Passed $false -Counts $counts `
      -Reason 'type errors increased over the committed baseline, or appeared in a file the baseline does not name.'
  }
  return New-GateVerdict -Lane 'typecheck' -Passed $true -Counts $counts
}

function Invoke-GateTestWeb {
  # The root config is the single vitest entry point: glob discovery over every
  # workspace (GATE-TOPOLOGY §2.3), so a test file in a package no lane names
  # still runs. The census lane holds the count against the tree.
  $result = Invoke-GateTool -Name 'vitest' -FilePath 'npx' `
    -ArgumentList @('--no-install', 'vitest', 'run', '--config', 'vitest.config.ts')
  $measurement = Measure-GateVitest -Text $result.Text
  $counts = "files=$($measurement.Files) tests=$($measurement.Tests) pass=$($measurement.Passed) fail=$($measurement.Failed)"
  Write-Host $counts
  if ($result.ExitCode -ne 0) {
    return New-GateVerdict -Lane 'test-web' -Passed $false -Counts $counts -Reason "vitest exited $($result.ExitCode)"
  }
  # A positive assertion, not a search for the string "0". A run that reports "no
  # test files", crashes before the summary, or whose output format drifted all
  # look like success to a check that only looks for a zero.
  if ($measurement.Passed -le 0) {
    return New-GateVerdict -Lane 'test-web' -Passed $false -Counts $counts `
      -Reason 'no summary line of the form "Tests <n> passed" with n>0. The suite did not execute.'
  }
  return New-GateVerdict -Lane 'test-web' -Passed $true -Counts $counts
}

function Invoke-GateSelftest {
  # The gate's own instruments, proved on themselves. Discovery is the glob plus
  # one naming convention: `*.integration.tests.ps1` needs Docker and belongs to
  # the integration lane, everything else runs here. No file list to forget --
  # a new self-test is in the lane the moment the file exists.
  $files = @(Get-ChildItem -Path (Join-Path $repositoryRoot 'scripts/tests') -Filter '*.tests.ps1' |
      Where-Object { $_.Name -notlike '*.integration.tests.ps1' } |
      Sort-Object Name)
  if ($files.Count -eq 0) {
    return New-GateVerdict -Lane 'selftest' -Passed $false -Counts 'files=0' `
      -Reason 'the glob reached no self-test file. An empty lane is not a green lane.'
  }

  $passedFiles = @()
  $failedFiles = @()
  foreach ($file in $files) {
    $result = Invoke-GateTool -Name "selftest-$($file.BaseName -replace '\.tests$', '')" -FilePath 'pwsh' `
      -ArgumentList @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', $file.FullName)
    # Exit code alone is not evidence: a Pester file invoked directly runs its
    # Describe blocks without setting the exit code, so a failed assertion can
    # exit 0. Each file must therefore also present a positive pass token --
    # either the script convention (`PASS ...` on its own line) or Pester
    # summaries whose failure count is zero and whose pass count is not.
    # Pester colors its summary, so the ANSI escapes sit between the numbers and
    # the words. Match the plain text.
    $plainText = $result.Text -replace "`e\[[0-9;]*m", ''
    $pesterPassed = 0
    $pesterFailed = 0
    foreach ($match in @([regex]::Matches($plainText, 'Tests Passed:\s*(\d+),\s*Failed:\s*(\d+)'))) {
      $pesterPassed += [int]$match.Groups[1].Value
      $pesterFailed += [int]$match.Groups[2].Value
    }
    $scriptToken = $plainText -match '(?m)^PASS '
    $ok = ($result.ExitCode -eq 0) -and ($pesterFailed -eq 0) -and ($scriptToken -or $pesterPassed -gt 0)
    if ($ok) { $passedFiles += $file.Name } else { $failedFiles += "$($file.Name) exit=$($result.ExitCode) pester_failed=$pesterFailed token=$(if ($scriptToken) { 'script' } elseif ($pesterPassed -gt 0) { 'pester' } else { 'none' })" }
  }

  $counts = "files=$($files.Count) pass=$($passedFiles.Count) fail=$($failedFiles.Count)"
  Write-Host $counts
  if ($failedFiles.Count -gt 0) {
    foreach ($line in $failedFiles) { Write-Host "  failed: $line" }
    return New-GateVerdict -Lane 'selftest' -Passed $false -Counts $counts `
      -Reason 'a self-test file failed, exited cleanly without a pass token, or reported Pester failures. The instruments are the thing under test here.'
  }
  return New-GateVerdict -Lane 'selftest' -Passed $true -Counts $counts
}

function Invoke-GateGuards {
  # Issue #3 mechanism (a): the inventory of guards, each held to the fixture
  # that proves it can fail -- and the fixture must EXECUTE, not merely exist.
  #
  # The verdict logic itself lives in Test-GateGuardInventory
  # (scripts/harness/Gate.psm1), so scripts/tests/guards-lane.tests.ps1 can
  # exercise it against synthetic inventories and injected runners. This
  # function stays the thin wrapper: parse the inventory, build the REAL
  # runners on Invoke-GateTool, hand both to the extracted function, print its
  # counts and failures exactly as before.
  $inventoryPath = Join-Path $repositoryRoot 'contracts/gate/guards.json'
  if (-not (Test-Path -LiteralPath $inventoryPath -PathType Leaf)) {
    return New-GateVerdict -Lane 'guards' -Passed $false -Counts 'inventory=missing' `
      -Reason 'contracts/gate/guards.json does not exist; an unlisted guard population is unmeasurable'
  }
  $inventory = Get-Content -LiteralPath $inventoryPath -Raw | ConvertFrom-Json

  $goRunner = {
    param([string]$Package, [string]$Pattern)
    $toolName = 'guards-' + (($Package -replace '[^A-Za-z0-9]+', '-').Trim('-'))
    Invoke-GateTool -Name $toolName -FilePath 'go' `
      -ArgumentList @('test', $Package, '-run', $Pattern, '-count=1', '-v') -WorkingDirectory $serverCore
  }
  # Same invocation Invoke-GateSelftest uses (pwsh -NoProfile -ExecutionPolicy
  # Bypass -File <path>), through the same Invoke-GateTool wrapper, so the
  # evidence this lane collects is the same shape as the selftest lane's.
  $pwshRunner = {
    param([string]$FilePath)
    $leaf = [IO.Path]::GetFileName($FilePath)
    $toolName = 'guards-pwsh-' + (($leaf -replace '\.tests\.ps1$', '') -replace '[^A-Za-z0-9]+', '-').Trim('-')
    Invoke-GateTool -Name $toolName -FilePath 'pwsh' `
      -ArgumentList @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', $FilePath)
  }

  $result = Test-GateGuardInventory -Inventory $inventory -RepositoryRoot $repositoryRoot `
    -GoRunner $goRunner -PwshRunner $pwshRunner

  # The empty-inventory short circuit prints nothing and returns immediately,
  # matching the pre-extraction lane exactly: no counts line was ever printed
  # for that case, only the terse 'entries=0' verdict.
  if ($result.Entries -eq 0) {
    return New-GateVerdict -Lane 'guards' -Passed $false -Counts $result.Counts -Reason $result.Reason
  }

  Write-Host $result.Counts
  if ($result.Failures.Count -gt 0) {
    foreach ($line in $result.Failures) { Write-Host "  guard: $line" }
    return New-GateVerdict -Lane 'guards' -Passed $false -Counts $result.Counts -Reason $result.Reason
  }
  return New-GateVerdict -Lane 'guards' -Passed $true -Counts $result.Counts
}

function Invoke-GateCensus {
  <#
    Every test file in the tree, held against the union of what the gate's lanes
    actually execute (GATE-TOPOLOGY §2.3: kills the dark-test-file class rather
    than one instance at a time). Static -- nothing runs; discovery is compared,
    not execution. Exemptions live in contracts/gate/census-exempt.json with a
    reason each, and a stale exemption (file gone, or file now reachable) fails
    the lane the same way an uncovered file does.
  #>
  # TypeScript: the tree census against vitest's own discovery over the root
  # config. `vitest list` resolves the same globs `vitest run` executes, so a
  # test file the run cannot see is named here before it goes dark.
  $tsCensus = @(& git -C $repositoryRoot ls-files --cached --others --exclude-standard '*.test.ts' '*.test.tsx' '*.spec.ts' '*.spec.tsx' |
      Where-Object { Test-Path -LiteralPath (Join-Path $repositoryRoot $_) -PathType Leaf } | Sort-Object)
  $listResult = Invoke-GateTool -Name 'census-vitest-list' -FilePath 'npx' -Quiet `
    -ArgumentList @('--no-install', 'vitest', 'list', '--config', 'vitest.config.ts', '--filesOnly')
  $discovered = @{}
  foreach ($line in @($listResult.Text -split "`r?`n")) {
    # ANSI first: a colored `[project]` prefix would defeat the bracket strip
    # and every discovered path would read as uncovered.
    $clean = ($line -replace "`e\[[0-9;]*m", '' -replace '^\[[^\]]+\]\s+', '').Trim() -replace '\\', '/'
    if ($clean) { $discovered[$clean] = $true }
  }
  $uncoveredTs = @($tsCensus | Where-Object { -not $discovered.ContainsKey($_) })

  # Go: the census against what `go test` can compile -- the unit set (./...)
  # plus the hermetic integration set, both with CGO_ENABLED=0 exactly as the
  # lanes run them, so the census measures the lanes and not a hypothetical
  # toolchain. The integration package list comes from the harness's own
  # discovery function; two copies of that rule is how they come to disagree.
  $goCensus = @(& git -C $repositoryRoot ls-files --cached --others --exclude-standard '*_test.go' |
      Where-Object { Test-Path -LiteralPath (Join-Path $repositoryRoot $_) -PathType Leaf } | Sort-Object)
  $goListFormat = '{{range .TestGoFiles}}{{$.Dir}}/{{.}}{{"\n"}}{{end}}{{range .XTestGoFiles}}{{$.Dir}}/{{.}}{{"\n"}}{{end}}'
  $previousCgo = $env:CGO_ENABLED
  try {
    $env:CGO_ENABLED = '0'
    $unitList = Invoke-GateTool -Name 'census-go-list' -FilePath 'go' -Quiet `
      -ArgumentList @('list', '-f', $goListFormat, './...') -WorkingDirectory $serverCore
    $integrationPackages = @(Get-HarnessIntegrationTestPackages -RepositoryRoot $repositoryRoot)
    $integrationList = Invoke-GateTool -Name 'census-go-list-integration' -FilePath 'go' -Quiet `
      -ArgumentList (@('list', '-tags=integration', '-f', $goListFormat) + $integrationPackages) -WorkingDirectory $serverCore
  } finally {
    $env:CGO_ENABLED = $previousCgo
  }
  $reachable = @{}
  $serverPrefix = (($serverCore -replace '\\', '/').TrimEnd('/')) + '/'
  foreach ($line in @(($unitList.Text + "`n" + $integrationList.Text) -split "`r?`n")) {
    $clean = ($line.Trim() -replace '\\', '/')
    if (-not $clean) { continue }
    $relative = 'apps/server_core/' + ($clean -replace [regex]::Escape($serverPrefix), '')
    $reachable[$relative] = $true
  }

  $exemptDocument = Get-Content -LiteralPath (Join-Path $repositoryRoot 'contracts/gate/census-exempt.json') -Raw | ConvertFrom-Json
  $exempt = @{}
  $staleExemptions = @()
  foreach ($entry in @($exemptDocument.exempt)) {
    if (-not (Test-Path -LiteralPath (Join-Path $repositoryRoot $entry.path) -PathType Leaf)) {
      $staleExemptions += "$($entry.path) (file no longer exists)"
    } elseif ($reachable.ContainsKey([string]$entry.path)) {
      $staleExemptions += "$($entry.path) (now reachable by the gate; the exemption is stale)"
    } else {
      $exempt[[string]$entry.path] = $true
    }
  }
  $uncoveredGo = @($goCensus | Where-Object { -not $reachable.ContainsKey($_) -and -not $exempt.ContainsKey($_) })

  $counts = "ts_census=$($tsCensus.Count) ts_discovered=$($discovered.Count) go_census=$($goCensus.Count) go_reachable=$($reachable.Count) exempt=$($exempt.Count) uncovered=$($uncoveredTs.Count + $uncoveredGo.Count) stale_exemptions=$($staleExemptions.Count)"
  Write-Host $counts
  if ($tsCensus.Count -eq 0 -or $discovered.Count -eq 0 -or $goCensus.Count -eq 0 -or $reachable.Count -eq 0) {
    return New-GateVerdict -Lane 'census' -Passed $false -Counts $counts `
      -Reason 'a census or discovery side came back empty. A comparison over an empty set proves nothing.'
  }
  if ($uncoveredTs.Count -gt 0 -or $uncoveredGo.Count -gt 0 -or $staleExemptions.Count -gt 0) {
    foreach ($path in @($uncoveredTs)) { Write-Host "  dark ts: $path" }
    foreach ($path in @($uncoveredGo)) { Write-Host "  dark go: $path" }
    foreach ($line in @($staleExemptions)) { Write-Host "  stale exemption: $line" }
    return New-GateVerdict -Lane 'census' -Passed $false -Counts $counts `
      -Reason 'a test file exists that no gate lane executes, or an exemption no longer describes reality. Name it in a lane or in census-exempt.json with a reason.'
  }
  return New-GateVerdict -Lane 'census' -Passed $true -Counts $counts
}

function Invoke-GateCiHygiene {
  # GATE-TOPOLOGY §2.3/§3: `pull_request_target` runs workflow code with write
  # credentials against a PR's HEAD, which is the textbook CI privilege
  # escalation. No workflow here needs it, so the string itself is the finding.
  $files = @(& git -C $repositoryRoot ls-files --cached --others --exclude-standard '.github/**' |
      Where-Object { Test-Path -LiteralPath (Join-Path $repositoryRoot $_) -PathType Leaf })
  if ($files.Count -eq 0) {
    return New-GateVerdict -Lane 'ci-hygiene' -Passed $false -Counts 'scanned=0' `
      -Reason 'no file under .github/ was scanned. Zero scanned is not zero findings.'
  }
  $hits = @()
  foreach ($file in $files) {
    foreach ($match in @(Select-String -LiteralPath (Join-Path $repositoryRoot $file) -Pattern 'pull_request_target' -SimpleMatch)) {
      $hits += "$($file):$($match.LineNumber)"
    }
  }
  $counts = "scanned=$($files.Count) hits=$($hits.Count)"
  Write-Host $counts
  if ($hits.Count -gt 0) {
    foreach ($hit in $hits) { Write-Host "  pull_request_target at $hit" }
    return New-GateVerdict -Lane 'ci-hygiene' -Passed $false -Counts $counts `
      -Reason 'a workflow uses pull_request_target, which runs PR code with write credentials.'
  }
  return New-GateVerdict -Lane 'ci-hygiene' -Passed $true -Counts $counts
}

function Invoke-GateIntegration {
  # Wraps `harness.ps1 -Command integration`: ephemeral Postgres in Docker,
  # embedded migrations, no ERP and no dev stack (GATE-TOPOLOGY §2.3). The
  # harness prints tests_run= with -1 for "the test step was never reached",
  # which is a different fact from 0 -- both fail here, for different reasons
  # the count makes attributable.
  $result = Invoke-GateTool -Name 'integration' -FilePath 'pwsh' `
    -ArgumentList @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', (Join-Path $repositoryRoot 'scripts/harness.ps1'), '-Command', 'integration')
  $match = [regex]::Match($result.Text, 'tests_run=(-?\d+) tests_passed=(-?\d+) tests_skipped=(-?\d+) tests_failed=(-?\d+)')
  $run = if ($match.Success) { [int]$match.Groups[1].Value } else { -1 }
  $passed = if ($match.Success) { [int]$match.Groups[2].Value } else { -1 }
  $skipped = if ($match.Success) { [int]$match.Groups[3].Value } else { -1 }
  $failed = if ($match.Success) { [int]$match.Groups[4].Value } else { -1 }
  $counts = "run=$run pass=$passed skip=$skipped fail=$failed"
  Write-Host $counts
  if ($result.ExitCode -ne 0) {
    return New-GateVerdict -Lane 'integration' -Passed $false -Counts $counts -Reason "the integration harness exited $($result.ExitCode)"
  }
  if ($run -le 0 -or $passed -le 0) {
    return New-GateVerdict -Lane 'integration' -Passed $false -Counts $counts `
      -Reason 'the lane exited 0 having run nothing or passed nothing. Pulled-vs-green is byte-identical without this line.'
  }
  if ($failed -ne 0) {
    return New-GateVerdict -Lane 'integration' -Passed $false -Counts $counts -Reason 'integration tests failed.'
  }
  return New-GateVerdict -Lane 'integration' -Passed $true -Counts $counts
}

function Invoke-GateEdge {
  # Wraps `harness.ps1 -Command edge`: the issue #1 PII fixture driving the real
  # Caddyfiles through real Caddy against stub upstreams, Docker only. The
  # fixture refuses to pass under 17 assertions; the lane re-asserts the floor
  # so weakening the fixture's own guard is visible here too.
  $result = Invoke-GateTool -Name 'edge' -FilePath 'pwsh' `
    -ArgumentList @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', (Join-Path $repositoryRoot 'scripts/harness.ps1'), '-Command', 'edge')
  $match = [regex]::Match($result.Text, '(?m)^assertions=(\d+)')
  $assertions = if ($match.Success) { [int]$match.Groups[1].Value } else { -1 }
  $counts = "assertions=$assertions"
  Write-Host $counts
  if ($result.ExitCode -ne 0) {
    return New-GateVerdict -Lane 'edge' -Passed $false -Counts $counts -Reason "the edge harness exited $($result.ExitCode)"
  }
  if ($assertions -lt 17) {
    return New-GateVerdict -Lane 'edge' -Passed $false -Counts $counts `
      -Reason 'fewer than 17 edge PII assertions ran (-1: no assertions= line at all). The deny surface was not exercised.'
  }
  return New-GateVerdict -Lane 'edge' -Passed $true -Counts $counts
}

# Cost order. `needs:` in ci.yml expresses the same ordering across jobs; here it
# is the order of this array.
$lanes = [ordered]@{
  'gofmt'      = ${function:Invoke-GateGofmt}
  'format'     = ${function:Invoke-GateFormat}
  'build'      = ${function:Invoke-GateBuild}
  'typecheck'  = ${function:Invoke-GateTypecheck}
  'governance' = ${function:Invoke-GateGovernance}
  'arch'       = ${function:Invoke-GateArch}
  'boundary'   = ${function:Invoke-GateBoundary}
  'lint-go'    = ${function:Invoke-GateLintGo}
  'lint-web'   = ${function:Invoke-GateLintWeb}
  'test-go'    = ${function:Invoke-GateGoTest}
  'test-web'   = ${function:Invoke-GateTestWeb}
  'census'     = ${function:Invoke-GateCensus}
  'ci-hygiene' = ${function:Invoke-GateCiHygiene}
  'selftest'   = ${function:Invoke-GateSelftest}
  'guards'     = ${function:Invoke-GateGuards}
  'integration' = ${function:Invoke-GateIntegration}
  'edge'       = ${function:Invoke-GateEdge}
}

# 'all' is the everyday set; 'full' adds the lanes whose cost is minutes rather
# than seconds (GATE-TOPOLOGY §2.3 verify-full). Both are the same product --
# CI's verify-full job and a pre-merge local run say `full`, a push says `all`.
$fullOnly = @('selftest', 'guards', 'integration', 'edge')
$selected = switch ($Lane) {
  'all' { @($lanes.Keys | Where-Object { $_ -notin $fullOnly }) }
  'full' { @($lanes.Keys) }
  default { @($Lane) }
}
$verdicts = @()
foreach ($name in $selected) {
  Write-Host ''
  Write-Host "=== gate lane=$name ==="
  $verdict = & $lanes[$name]
  $verdicts += $verdict
  Write-Host "lane=$name verdict=$(if ($verdict.Passed) { 'PASS' } else { 'FAIL' }) $($verdict.Counts)"
  if (-not $verdict.Passed) {
    Write-Host "reason: $($verdict.Reason)"
    if (-not $ContinueOnFailure) { break }
  }
}

Write-Host ''
Write-Host '=== gate summary ==='
foreach ($verdict in $verdicts) {
  Write-Host ("{0,-10} {1,-4} {2}" -f $verdict.Lane, $(if ($verdict.Passed) { 'PASS' } else { 'FAIL' }), $verdict.Counts)
}
$ran = @($verdicts).Count
$failed = @($verdicts | Where-Object { -not $_.Passed }).Count
$skippedLanes = @($selected).Count - $ran
Write-Host "lanes_selected=$(@($selected).Count) lanes_ran=$ran lanes_failed=$failed lanes_not_reached=$skippedLanes"

if ($ran -eq 0) {
  Write-Host 'gate: FAIL -- no lane ran. The gate itself is the vacuous green here.'
  exit 1
}
if ($failed -ne 0) {
  Write-Host 'gate: FAIL'
  exit 1
}
Write-Host 'gate: PASS'
exit 0
