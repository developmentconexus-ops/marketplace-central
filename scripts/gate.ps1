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

  This file is the set. `npm run gate` runs it locally; `.github/workflows/ci.yml`
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
  [ValidateSet('all', 'gofmt', 'build', 'test-go', 'typecheck', 'test-web')]
  [string]$Lane = 'all',
  [switch]$ContinueOnFailure
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$repositoryRoot = [IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..'))
$serverCore = Join-Path $repositoryRoot 'apps/server_core'
Import-Module (Join-Path $PSScriptRoot 'harness/Gate.psm1') -Force

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
    [string]$WorkingDirectory = $repositoryRoot
  )

  $logPath = Join-Path $logDirectory "$Name.log"
  Push-Location -LiteralPath $WorkingDirectory
  try {
    # 2>&1 folds stderr into the stream: go and tsc write real diagnostics there,
    # and a count taken over stdout alone would miss them.
    & $FilePath @ArgumentList 2>&1 | Tee-Object -FilePath $logPath | Out-Host
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
  $discovered = @(& git -C $repositoryRoot ls-files --cached --others --exclude-standard '*.go')
  Write-Host "discovered=$($discovered.Count)"
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
  for ($index = 0; $index -lt $discovered.Count; $index += 200) {
    $batch = @($discovered[$index..([Math]::Min($index + 199, $discovered.Count - 1))] |
      ForEach-Object { Join-Path $repositoryRoot $_ })
    $result = & gofmt -l @batch 2>&1
    foreach ($line in @($result)) { [void]$output.AppendLine([string]$line) }
  }
  $measurement = Measure-GateGofmt -Text $output.ToString()
  $counts = "run=$($discovered.Count) unformatted=$($measurement.Unformatted)"
  Write-Host $counts
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
  $build = Invoke-GateTool -Name 'go-build' -FilePath 'go' -ArgumentList @('build', './...') -WorkingDirectory $serverCore
  if ($build.ExitCode -ne 0) {
    return New-GateVerdict -Lane 'build' -Passed $false -Counts 'build=fail' -Reason "go build ./... exited $($build.ExitCode)"
  }
  $vet = Invoke-GateTool -Name 'go-vet' -FilePath 'go' -ArgumentList @('vet', './...') -WorkingDirectory $serverCore
  if ($vet.ExitCode -ne 0) {
    return New-GateVerdict -Lane 'build' -Passed $false -Counts 'build=ok vet=fail' -Reason "go vet ./... exited $($vet.ExitCode)"
  }
  Write-Host 'build=ok vet=ok'
  return New-GateVerdict -Lane 'build' -Passed $true -Counts 'build=ok vet=ok'
}

function Invoke-GateGoTest {
  # ./... and not ./tests/unit/...: the latter is 14 test files, internal/ holds
  # 353. The integration lane is invisible here regardless -- it sits behind
  # `//go:build integration`, which ./... does not compile.
  #
  # -skip TestModuleBoundaryADR023: that test is red on main by design (234
  # cross-module import violations) and already has an owner -- Onda 1 of
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
  $result = Invoke-GateTool -Name 'tsc' -FilePath 'npx' `
    -ArgumentList @('--no-install', 'tsc', '-p', 'apps/web/tsconfig.json', '--noEmit', '--listFiles')
  $measurement = Measure-GateTsc -Text $result.Text
  $counts = "checked=$($measurement.Checked) errors=$($measurement.Errors)"
  Write-Host $counts
  if ($result.ExitCode -ne 0) {
    return New-GateVerdict -Lane 'typecheck' -Passed $false -Counts $counts -Reason "tsc exited $($result.ExitCode)"
  }
  if ($measurement.Checked -eq 0) {
    return New-GateVerdict -Lane 'typecheck' -Passed $false -Counts $counts `
      -Reason 'tsc loaded no project files. The project resolved nothing, which is not the same as compiling cleanly.'
  }
  return New-GateVerdict -Lane 'typecheck' -Passed $true -Counts $counts
}

function Invoke-GateTestWeb {
  $result = Invoke-GateTool -Name 'vitest' -FilePath 'npm' `
    -ArgumentList @('run', 'test', '--workspace', '@marketplace-central/web', '--', '--run')
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

# Cost order. `needs:` in ci.yml expresses the same ordering across jobs; here it
# is the order of this array.
$lanes = [ordered]@{
  'gofmt'     = ${function:Invoke-GateGofmt}
  'build'     = ${function:Invoke-GateBuild}
  'typecheck' = ${function:Invoke-GateTypecheck}
  'test-go'   = ${function:Invoke-GateGoTest}
  'test-web'  = ${function:Invoke-GateTestWeb}
}

$selected = if ($Lane -eq 'all') { @($lanes.Keys) } else { @($Lane) }
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
