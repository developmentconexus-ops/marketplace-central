Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

# The counting half of the gate, kept apart from the process-spawning half so it
# can be tested against recorded tool output instead of against a live toolchain.
#
# Every function here answers the same question in a different dialect: how many
# units did this tool actually execute? The question exists because a tool that
# executed nothing exits 0. `go test` on an empty package set, `vitest` with no
# matching files, and `tsc` over a project that resolved no sources are all
# indistinguishable from success by exit code alone -- which is finding V2 of the
# audit of 2026-08-07, and the reason `scripts/harness/Postgres.psm1` grew the
# same guard for the integration lane.

function Measure-GateGoTest {
  <#
    Counts `go test -v` verdict lines.

    The `=== RUN` and `--- PASS:` anchors tolerate leading whitespace because Go
    indents subtest results. Counting subtests is deliberate: they are executed
    units, and a package whose only real assertions live in subtests would
    otherwise measure as a handful of parents.

    SKIP is counted separately and never folded into PASS. A suite where every
    test skipped prints `ok` and exits 0 -- HARNESS-PROFILE.md records exactly
    that shape as RUN 27 / PASS 1 / SKIP 26, with the slice's whole reason to
    exist among the skips.
  #>
  [CmdletBinding()]
  param([Parameter(Mandatory)][AllowEmptyString()][string]$Text)

  return [pscustomobject]@{
    Run     = @([regex]::Matches($Text, '(?m)^\s*=== RUN\s')).Count
    Passed  = @([regex]::Matches($Text, '(?m)^\s*--- PASS:')).Count
    Failed  = @([regex]::Matches($Text, '(?m)^\s*--- FAIL:')).Count
    Skipped = @([regex]::Matches($Text, '(?m)^\s*--- SKIP:')).Count
  }
}

function Remove-GateAnsi {
  [CmdletBinding()]
  param([Parameter(Mandatory)][AllowEmptyString()][string]$Text)

  return [regex]::Replace($Text, "`e\[[0-9;]*[A-Za-z]", '')
}

function Measure-GateVitest {
  <#
    Reads vitest's summary block. The counts are taken from the summary rather
    than by counting per-file lines, because a run that crashed after printing
    some file lines and before printing the summary must not measure as a pass.
    No summary means -1, which is not zero: zero is a suite that ran and found
    nothing, -1 is a suite whose result was never stated.
  #>
  [CmdletBinding()]
  param([Parameter(Mandatory)][AllowEmptyString()][string]$Text)

  $plain = Remove-GateAnsi -Text $Text
  $files = [regex]::Match($plain, '(?m)^\s*Test Files\s+.*?\((?<total>\d+)\)\s*$')
  $tests = [regex]::Match($plain, '(?m)^\s*Tests\s+.*?\((?<total>\d+)\)\s*$')
  $passed = [regex]::Match($plain, '(?m)^\s*Tests\s+(?:.*?\|\s*)?(?<passed>\d+)\s+passed')
  $failed = [regex]::Match($plain, '(?m)^\s*Tests\s+(?<failed>\d+)\s+failed')

  return [pscustomobject]@{
    Files   = if ($files.Success) { [int]$files.Groups['total'].Value } else { -1 }
    Tests   = if ($tests.Success) { [int]$tests.Groups['total'].Value } else { -1 }
    Passed  = if ($passed.Success) { [int]$passed.Groups['passed'].Value } else { -1 }
    # Three states, not two. A summary that names no failures is the fact "zero
    # failed"; no summary at all is "unknown", and reporting that as zero turns a
    # crashed run into an affirmative claim that nothing failed -- the one thing
    # AGENTS.md forbids an unknown operational fact from becoming.
    Failed  = if ($failed.Success) { [int]$failed.Groups['failed'].Value }
              elseif ($tests.Success) { 0 }
              else { -1 }
  }
}

function Measure-GateTsc {
  <#
    Counts the project files `tsc --listFiles` reported it loaded, excluding
    `node_modules` -- the 319 declaration files under there are the toolchain's
    own, and counting them would make a run over zero source files report 319.

    Error lines are shaped `path(line,col): error TS0000: message` and so never
    match the path-only pattern.
  #>
  [CmdletBinding()]
  param([Parameter(Mandatory)][AllowEmptyString()][string]$Text)

  $lines = @($Text -split "`r?`n" | ForEach-Object { $_.Trim() } | Where-Object { $_ })
  $project = @($lines | Where-Object {
      $_ -match '\.(ts|tsx|mts|cts|js|jsx|json)$' -and
      $_ -notmatch '[/\\]node_modules[/\\]' -and
      $_ -notmatch ':\s+error\s+TS\d+'
    })
  return [pscustomobject]@{
    Checked = $project.Count
    Errors  = @([regex]::Matches($Text, '(?m):\s+error\s+TS\d+')).Count
  }
}

function Measure-GateGofmt {
  <#
    `gofmt -l` prints one path per unformatted file and nothing at all when the
    tree is clean, so the count of output lines IS the violation count. The
    number of files fed to it is measured separately by the caller; without it a
    broken file filter reports the same clean zero as a formatted tree.
  #>
  [CmdletBinding()]
  param([Parameter(Mandatory)][AllowEmptyString()][string]$Text)

  $paths = @($Text -split "`r?`n" | ForEach-Object { $_.Trim() } | Where-Object { $_ })
  return [pscustomobject]@{
    Unformatted = $paths.Count
    Paths       = $paths
  }
}

function Measure-GateGolangciLint {
  <#
    Reads golangci-lint's JSON report. The counts come from the issue list rather
    than from the human-readable summary, because the summary is formatting and
    the list is the fact.

    `Enabled` is returned so the caller can assert the enabled linter set against
    the one the config asks for. A config file that failed to load leaves
    golangci-lint running its own defaults and reporting a healthy small number
    over the wrong rules -- green, for a reason unrelated to the code.
  #>
  [CmdletBinding()]
  param([Parameter(Mandatory)][AllowEmptyString()][string]$Json)

  if ([string]::IsNullOrWhiteSpace($Json)) {
    return [pscustomobject]@{ Total = -1; ByLinter = @{}; Enabled = @() }
  }
  $report = $Json | ConvertFrom-Json
  $issues = @()
  if ($report.PSObject.Properties.Name -contains 'Issues' -and $null -ne $report.Issues) {
    $issues = @($report.Issues)
  }
  $byLinter = @{}
  foreach ($issue in $issues) {
    $linter = [string]$issue.FromLinter
    if (-not $byLinter.ContainsKey($linter)) { $byLinter[$linter] = 0 }
    $byLinter[$linter] = $byLinter[$linter] + 1
  }
  $enabled = @()
  if ($report.PSObject.Properties.Name -contains 'Report' -and $null -ne $report.Report -and
      $report.Report.PSObject.Properties.Name -contains 'Linters') {
    $enabled = @($report.Report.Linters |
      Where-Object { $_.PSObject.Properties.Name -contains 'Enabled' -and $_.Enabled } |
      ForEach-Object { [string]$_.Name })
  }
  return [pscustomobject]@{
    Total    = $issues.Count
    ByLinter = $byLinter
    Enabled  = @($enabled | Sort-Object)
  }
}

function Measure-GateEslint {
  <#
    Reads ESLint's JSON report.

    `Paths` is the point of this function, not `Total`. ESLint reports one entry
    per file it linted, including files with nothing to say, so the report states
    what was covered as well as what was found -- and coverage is the failure mode
    here. A type-aware rule over a file that belongs to no tsconfig cannot run;
    ESLint reports the file, the rule stays silent, and the lane reads a clean
    zero it did not earn. The caller compares `Paths` against the tracked source
    list and fails on a file that is missing.

    `Fatal` is counted apart from the rules. A parse failure is a message with no
    `ruleId`, so folding it into the findings total would let a file that could
    not be read at all shrink the ratchet the moment someone adds it to an
    ignore list.
  #>
  [CmdletBinding()]
  param([Parameter(Mandatory)][AllowEmptyString()][string]$Json)

  if ([string]::IsNullOrWhiteSpace($Json)) {
    return [pscustomobject]@{ Files = -1; Total = -1; Fatal = -1; ByRule = @{}; Paths = @() }
  }
  $report = @($Json | ConvertFrom-Json)
  $byRule = @{}
  $fatal = 0
  $paths = @()
  foreach ($file in $report) {
    $paths += ([string]$file.filePath -replace '\\', '/')
    foreach ($message in @($file.messages)) {
      if ($message.PSObject.Properties.Name -contains 'fatal' -and $message.fatal) { $fatal++; continue }
      $rule = if ($message.ruleId) { [string]$message.ruleId } else { '(no-rule)' }
      if (-not $byRule.ContainsKey($rule)) { $byRule[$rule] = 0 }
      $byRule[$rule] = $byRule[$rule] + 1
    }
  }
  return [pscustomobject]@{
    Files  = $report.Count
    Total  = @($byRule.Values | Measure-Object -Sum).Sum
    Fatal  = $fatal
    ByRule = $byRule
    Paths  = $paths
  }
}

function Measure-GatePrettier {
  <#
    Reads `prettier --check` output.

    Prettier prints one `[warn] <path>` line per file whose formatting differs and
    a final `[warn] Code style issues found in N files`. On a clean run it prints
    `All matched files use Prettier code style!` and names no file, which is the
    reason the caller must count the candidate set itself: a glob that matched
    nothing and a tree that is perfectly formatted produce the same two lines.

    `Clean` is the affirmative marker, not the absence of warnings. A run that
    died before printing anything has no warnings either.
  #>
  [CmdletBinding()]
  param([Parameter(Mandatory)][AllowEmptyString()][string]$Text)

  $plain = Remove-GateAnsi -Text $Text
  $summary = [regex]::Match($plain, '(?m)Code style issues found in (?<count>\d+) file')
  $listed = @([regex]::Matches($plain, '(?m)^\[warn\] (?<path>(?!Code style)\S.*)$') |
    ForEach-Object { $_.Groups['path'].Value.Trim() })
  return [pscustomobject]@{
    Clean       = ($plain -match 'All matched files use Prettier code style!')
    Unformatted = if ($summary.Success) { [int]$summary.Groups['count'].Value } else { $listed.Count }
    Paths       = $listed
  }
}

function Compare-GateRatchet {
  <#
    Shrink-only comparison of measured counts against a committed baseline.

    Three verdicts, and the third is the one that makes this an instrument rather
    than a formality:

      - measured <= baseline for every key: PASS, and any shrink is reported so
        the lower number can be committed in the same PR that earned it.
      - measured > baseline for any key: FAIL. This is the ratchet.
      - a key present in the measurement and absent from the baseline: FAIL. An
        unknown key means the enabled set moved without the baseline moving with
        it. Treating it as a baseline of zero would be strictly stricter and
        sounds safer, but it hides the drift behind a finding about code; the
        drift is the thing to report.

    A key present in the baseline and absent from the measurement is a shrink to
    zero, not an error -- that is the direction this is supposed to move.
  #>
  [CmdletBinding()]
  param(
    [Parameter(Mandatory)][hashtable]$Measured,
    [Parameter(Mandatory)][hashtable]$Baseline
  )

  $increased = @()
  $decreased = @()
  $unknown = @()
  foreach ($key in @($Measured.Keys | Sort-Object)) {
    if (-not $Baseline.ContainsKey($key)) {
      $unknown += "$key=$($Measured[$key])"
      continue
    }
    $delta = [int]$Measured[$key] - [int]$Baseline[$key]
    if ($delta -gt 0) { $increased += "$key=$($Measured[$key]) baseline=$($Baseline[$key]) +$delta" }
    elseif ($delta -lt 0) { $decreased += "$key=$($Measured[$key]) baseline=$($Baseline[$key]) $delta" }
  }
  foreach ($key in @($Baseline.Keys | Sort-Object)) {
    if ($Measured.ContainsKey($key)) { continue }
    if ([int]$Baseline[$key] -gt 0) { $decreased += "$key=0 baseline=$($Baseline[$key]) -$($Baseline[$key])" }
  }
  return [pscustomobject]@{
    Passed    = ($increased.Count -eq 0 -and $unknown.Count -eq 0)
    Increased = $increased
    Decreased = $decreased
    Unknown   = $unknown
  }
}

function ConvertTo-GateCountMap {
  <#
    PSCustomObject to hashtable. ConvertFrom-Json yields the former and
    Compare-GateRatchet takes the latter; doing this inline at both call sites is
    how the two ends come to disagree about a missing key.
  #>
  [CmdletBinding()]
  param([Parameter(Mandatory)][AllowNull()]$Object)

  $map = @{}
  if ($null -eq $Object) { return $map }
  foreach ($property in $Object.PSObject.Properties) { $map[$property.Name] = [int]$property.Value }
  return $map
}

Export-ModuleMember -Function Measure-GateGoTest, Measure-GateVitest, Measure-GateTsc, Measure-GateGofmt, `
  Remove-GateAnsi, Measure-GateGolangciLint, Measure-GateEslint, Measure-GatePrettier, `
  Compare-GateRatchet, ConvertTo-GateCountMap
