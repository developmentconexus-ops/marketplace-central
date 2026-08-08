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
  golangci-lint, ESLint, `go test ./...`, vitest. `npm run gate` runs it locally; `.github/workflows/ci.yml`
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
  [ValidateSet('all', 'gofmt', 'format', 'build', 'lint-go', 'lint-web', 'test-go', 'typecheck', 'test-web')]
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
  'bodyclose', 'errcheck', 'errorlint', 'exhaustive', 'ineffassign',
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
    return New-GateVerdict -Lane 'lint-go' -Passed $false -Counts 'report=missing' `
      -Reason "golangci-lint wrote no report (exit $($result.ExitCode)). The lane cannot report zero findings for a run that did not happen."
  }
  $measurement = Measure-GateGolangciLint -Json (Get-Content -LiteralPath $reportPath -Raw)

  # Exit 1 is what golangci-lint returns when it found issues, which is the
  # normal state under a ratchet. Anything else is the tool failing to run, and
  # that must not read as a clean lane.
  if ($result.ExitCode -notin @(0, 1)) {
    return New-GateVerdict -Lane 'lint-go' -Passed $false -Counts "total=$($measurement.Total)" `
      -Reason "golangci-lint exited $($result.ExitCode), which is neither clean nor findings."
  }

  $enabledDrift = @(Compare-Object -ReferenceObject $script:GolangciLintExpected -DifferenceObject $measurement.Enabled)
  if ($enabledDrift.Count -ne 0) {
    $detail = ($enabledDrift | ForEach-Object { "$($_.SideIndicator) $($_.InputObject)" }) -join ', '
    return New-GateVerdict -Lane 'lint-go' -Passed $false -Counts "enabled=$($measurement.Enabled.Count)" `
      -Reason "the enabled linter set is not the configured one ($detail). A count over the wrong rules is not a smaller count."
  }

  $measuredMap = @{}
  foreach ($key in @($baseline.Keys)) { $measuredMap[$key] = 0 }
  foreach ($key in @($measurement.ByLinter.Keys)) { $measuredMap[$key] = $measurement.ByLinter[$key] }
  $comparison = Compare-GateRatchet -Measured $measuredMap -Baseline $baseline

  $counts = "total=$($measurement.Total) baseline=$($baselineDocument.'golangci-lint'.total) enabled=$($measurement.Enabled.Count)"
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
  'format'    = ${function:Invoke-GateFormat}
  'build'     = ${function:Invoke-GateBuild}
  'typecheck' = ${function:Invoke-GateTypecheck}
  'lint-go'   = ${function:Invoke-GateLintGo}
  'lint-web'  = ${function:Invoke-GateLintWeb}
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
