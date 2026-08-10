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
  # Errors attributed to their file so the lane can ratchet per file: a new
  # error in a clean file is an unknown key, whatever another file shrank by.
  $byFile = @{}
  foreach ($match in @([regex]::Matches($Text, '(?m)^(.+?)\(\d+,\d+\):\s+error\s+TS\d+'))) {
    $errorPath = $match.Groups[1].Value.Trim()
    if ($byFile.ContainsKey($errorPath)) { $byFile[$errorPath] = $byFile[$errorPath] + 1 } else { $byFile[$errorPath] = 1 }
  }
  return [pscustomobject]@{
    Checked = $project.Count
    Paths   = $project
    Errors  = @([regex]::Matches($Text, '(?m):\s+error\s+TS\d+')).Count
    ByFile  = $byFile
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

    `Checked` is the effective file set, and it needs `--log-level debug` to
    exist. At that verbosity Prettier emits one `resolve config from '<path>'`
    line per file it actually reads, and nothing else in its output distinguishes
    a formatted tree from a glob that matched nothing. Counted distinct, because
    the same path resolving twice is one file. Without the debug flag this is 0
    and the caller must treat that as the failure it is.
  #>
  [CmdletBinding()]
  param([Parameter(Mandatory)][AllowEmptyString()][string]$Text)

  $plain = Remove-GateAnsi -Text $Text
  $summary = [regex]::Match($plain, '(?m)Code style issues found in (?<count>\d+) file')
  $listed = @([regex]::Matches($plain, '(?m)^\[warn\] (?<path>(?!Code style)\S.*)$') |
    ForEach-Object { $_.Groups['path'].Value.Trim() })
  $resolved = @{}
  foreach ($match in [regex]::Matches($plain, "(?m)resolve config from '(?<path>[^']+)'")) {
    $resolved[$match.Groups['path'].Value] = $true
  }
  return [pscustomobject]@{
    Clean       = ($plain -match 'All matched files use Prettier code style!')
    Unformatted = if ($summary.Success) { [int]$summary.Groups['count'].Value } else { $listed.Count }
    Paths       = $listed
    Checked     = $resolved.Count
  }
}

function Measure-GateEslintRules {
  <#
    Reads `eslint --print-config <file>` and returns the rules that are actually
    ON for that file.

    This is the ESLint half of the lint-go lane's enabled-linter assertion, and
    it exists because a findings count cannot tell "the code got better" from
    "the rule got turned off". Reproduced 2026-08-08 on PR #21: setting two of
    the six rules to "off" dropped the lane's total from 26 to 13, and the ratchet
    reported the drop as a shrink worth committing.

    ESLint's severity is 0/1/2 or "off"/"warn"/"error", and a configured rule is
    an array whose first element is the severity. Anything resolving to 0 is
    absent from the returned set, exactly as if it had never been named.
  #>
  [CmdletBinding()]
  param([Parameter(Mandatory)][AllowEmptyString()][string]$Json)

  if ([string]::IsNullOrWhiteSpace($Json)) { return @() }
  $document = $Json | ConvertFrom-Json
  if ($null -eq $document -or $null -eq $document.rules) { return @() }

  $enabled = @()
  foreach ($property in $document.rules.PSObject.Properties) {
    $value = $property.Value
    $severity = if ($value -is [array] -or $value -is [System.Collections.IList]) {
      if (@($value).Count -eq 0) { 0 } else { @($value)[0] }
    } else { $value }
    $isOff = ($null -eq $severity) -or ($severity -eq 0) -or ("$severity" -eq 'off')
    if (-not $isOff) { $enabled += $property.Name }
  }
  return @($enabled | Sort-Object)
}

function Measure-GateArchscan {
  <#
    Reads archscan's output: one `file:line: rule: detail` per finding and a
    final `archscan: scanned=N findings=M` line.

    `Scanned` is the anti-vacuity number and it defaults to -1, not 0: a stream
    with no summary line means the tool died before finishing, and the caller
    must not read that as "scanned nothing" -- it is "nothing is known". The
    tool itself exits 2 when scanned is 0; the count is re-read here so the
    lane's verdict does not rest on an exit code alone.
  #>
  [CmdletBinding()]
  param([Parameter(Mandatory)][AllowEmptyString()][string]$Text)

  $plain = Remove-GateAnsi -Text $Text
  $summary = [regex]::Match($plain, '(?m)^archscan: scanned=(?<scanned>\d+) findings=(?<findings>\d+)\s*$')
  $byRule = @{}
  foreach ($match in [regex]::Matches($plain, '(?m)^\S+?:\d+: (?<rule>[a-z]+/[a-z-]+): ')) {
    $rule = $match.Groups['rule'].Value
    if (-not $byRule.ContainsKey($rule)) { $byRule[$rule] = 0 }
    $byRule[$rule] = $byRule[$rule] + 1
  }
  return [pscustomobject]@{
    Scanned  = if ($summary.Success) { [int]$summary.Groups['scanned'].Value } else { -1 }
    Total    = if ($summary.Success) { [int]$summary.Groups['findings'].Value } else { -1 }
    ByRule   = $byRule
  }
}

function Measure-GateBoundary {
  <#
    Reads `go test -v` output of TestModuleBoundaryADR023, the ADR-023 §2
    module-boundary detector that is red on main by design.

    Three facts, each with its own default when the stream does not state it:

      - Files: the `boundary_files=N` line the test logs on every outcome.
        Defaults to -1 -- a stream without it means the run died before the walk
        finished, or the log line was deleted, and neither may read as "walked
        an empty tree".
      - Total: the `N violation(s)` failure line, or 0 when the test PASSed.
        Defaults to -1 when neither appears: a compile error prints neither.
      - ByOrigin: the `by origin layer:` block only. The `by target:` block has
        the same count-tab-name shape, so the block is cut out before its lines
        are read -- a regex over the whole stream would double-count every
        violation under a second key set.
  #>
  [CmdletBinding()]
  param([Parameter(Mandatory)][AllowEmptyString()][string]$Text)

  $plain = Remove-GateAnsi -Text $Text
  $files = [regex]::Match($plain, '(?m)boundary_files=(?<n>\d+)\s*$')
  $violations = [regex]::Match($plain, '(?m)^\s*\S+_test\.go:\d+: (?<n>\d+) violation\(s\)')
  $passed = $plain -match '(?m)^--- PASS: TestModuleBoundaryADR023'

  $byOrigin = @{}
  $block = [regex]::Match($plain, 'by origin layer:\r?\n(?<block>(?:\s*\d+\t[^\r\n]+\r?\n)+)')
  if ($block.Success) {
    foreach ($match in [regex]::Matches($block.Groups['block'].Value, '(?m)^\s*(?<count>\d+)\t(?<origin>[^\r\n]+?)\s*$')) {
      $byOrigin[$match.Groups['origin'].Value] = [int]$match.Groups['count'].Value
    }
  }
  return [pscustomobject]@{
    Files    = if ($files.Success) { [int]$files.Groups['n'].Value } else { -1 }
    Total    = if ($violations.Success) { [int]$violations.Groups['n'].Value } elseif ($passed) { 0 } else { -1 }
    ByOrigin = $byOrigin
  }
}

function Measure-GateGovernance {
  <#
    Groups a PolicyResult's violations by error code, for the ratchet. The lane
    calls Test-GovernanceDrift directly and hands the structured Violations list
    here, so unlike the other measurers there is no text to parse -- but the
    grouping still lives in this module, because the tests that prove the lane
    counts correctly have to be able to call it with fixtures.
  #>
  [CmdletBinding()]
  param([AllowEmptyCollection()][object[]]$Violations = @())

  $byCode = @{}
  foreach ($violation in @($Violations)) {
    $code = [string]$violation.ErrorCode
    if (-not $byCode.ContainsKey($code)) { $byCode[$code] = 0 }
    $byCode[$code] = $byCode[$code] + 1
  }
  return [pscustomobject]@{
    Total  = @($Violations).Count
    ByCode = $byCode
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

function Measure-GateGuards {
  <#
    Holds a targeted `go test -run -v` output against the inventory's expected
    test names. `\b` after the escaped name: `--- PASS: TestA` must not satisfy
    an expectation of `TestAB`, and vice-versa Go suffixes the line with the
    duration so a bare $ anchor would never match.
  #>
  [CmdletBinding()]
  param(
    [Parameter(Mandatory)][AllowEmptyString()][string]$Text,
    [Parameter(Mandatory)][AllowEmptyCollection()][string[]]$Expected
  )
  $passed = [Collections.Generic.List[string]]::new()
  $missing = [Collections.Generic.List[string]]::new()
  foreach ($name in $Expected) {
    if ([regex]::IsMatch($Text, "(?m)^\s*--- PASS: $([regex]::Escape($name))\b")) { [void]$passed.Add($name) } else { [void]$missing.Add($name) }
  }
  return [pscustomobject]@{
    Ran     = @([regex]::Matches($Text, '(?m)^\s*=== RUN\s')).Count
    Passed  = @($passed)
    Missing = @($missing)
    Failed  = @([regex]::Matches($Text, '(?m)^\s*--- FAIL:')).Count
  }
}

function Test-GateGuardInventory {
  <#
    The guards lane's verdict logic (issue #3 mechanism (a)), extracted from
    `Invoke-GateGuards` in scripts/gate.ps1 so it can be exercised from a
    self-test against synthetic inventories and injected runners, instead of
    only from a live `go test` / pwsh spawn against the committed
    contracts/gate/guards.json.

    Takes the parsed inventory object, the repository root the pwsh_files and
    presence_only paths resolve against, and two injected runner script
    blocks -- one per execution path:

      GoRunner:   param($Package, $Pattern) -> [pscustomobject]@{ ExitCode; Text }
      PwshRunner: param($FilePath)          -> [pscustomobject]@{ ExitCode; Text }

    The real caller builds these on Invoke-GateTool exactly as the lane did
    before extraction; a test builds them as bare script blocks returning
    canned output, so the verdict logic runs with no real process spawned.

    Returns a result object carrying the failure lines and the counts this
    lane prints: Entries (go+pwsh+presence total), GoExecuted, PwshEntries,
    PwshExecuted, Presence, plus Passed/Reason/Counts for the verdict itself.
    Counts is pre-formatted so the caller does not have to reconstruct the
    lane's two distinct shapes (the 'entries=0' short-circuit versus the full
    five-field line) -- that reconstruction is exactly the kind of drift this
    extraction exists to prevent.
  #>
  [CmdletBinding()]
  param(
    [Parameter(Mandatory)]$Inventory,
    [Parameter(Mandatory)][string]$RepositoryRoot,
    [Parameter(Mandatory)][scriptblock]$GoRunner,
    [Parameter(Mandatory)][scriptblock]$PwshRunner
  )

  $goEntries = @($Inventory.go_tests); $pwshEntries = @($Inventory.pwsh_files); $presenceEntries = @($Inventory.presence_only)
  $total = $goEntries.Count + $pwshEntries.Count + $presenceEntries.Count
  if ($total -eq 0) {
    return [pscustomobject]@{
      Passed       = $false
      Reason       = 'an empty inventory is not a green inventory'
      Counts       = 'entries=0'
      Entries      = 0
      GoExecuted   = 0
      PwshEntries  = 0
      PwshExecuted = 0
      Presence     = 0
      Failures     = @()
    }
  }

  $failures = [Collections.Generic.List[string]]::new()
  $goExecuted = 0
  foreach ($group in @($goEntries | Group-Object -Property package)) {
    $names = @($group.Group | ForEach-Object { [string]$_.test })
    $pattern = '^(' + (@($names | ForEach-Object { [regex]::Escape($_) }) -join '|') + ')$'
    $result = & $GoRunner $group.Name $pattern
    $measurement = Measure-GateGuards -Text $result.Text -Expected $names
    if ($result.ExitCode -ne 0) { [void]$failures.Add("package $($group.Name) exited $($result.ExitCode)") }
    if ($measurement.Ran -eq 0) { [void]$failures.Add("package $($group.Name): zero tests ran; the inventory names tests that do not execute") }
    foreach ($name in $measurement.Missing) { [void]$failures.Add("no '--- PASS: $name' in $($group.Name)") }
    $goExecuted += @($measurement.Passed).Count
  }

  # pwsh_files entries are EXECUTED here, not merely inspected for a string.
  # Grouped by file first: several ids can share one fixture file (governance-drift
  # alone carries 11), and running that file once per id would run it 11 times for
  # no additional evidence. Each distinct file runs once; every entry pointing at
  # it is checked against that one output.
  $pwshExecuted = 0
  foreach ($group in @($pwshEntries | Group-Object -Property file)) {
    $file = Join-Path $RepositoryRoot ([string]$group.Name)
    if (-not (Test-Path -LiteralPath $file -PathType Leaf)) {
      [void]$failures.Add("$($group.Name) is missing")
      continue
    }
    $leaf = [IO.Path]::GetFileName($file)
    $inSelftestGlob = ($leaf -like '*.tests.ps1') -and ($leaf -notlike '*.integration.tests.ps1') -and
      ([IO.Path]::GetFullPath([IO.Path]::GetDirectoryName($file)) -eq [IO.Path]::GetFullPath((Join-Path $RepositoryRoot 'scripts/tests')))
    if (-not $inSelftestGlob) {
      [void]$failures.Add("$($group.Name) sits outside the selftest lane's discovery glob; its execution is delegated to nothing")
    }

    $result = & $PwshRunner $file
    $pwshExecuted++
    if ($result.ExitCode -ne 0) { [void]$failures.Add("$($group.Name) exited $($result.ExitCode)") }
    foreach ($entry in $group.Group) {
      $token = [string]$entry.token
      if (-not ([regex]::IsMatch($result.Text, "(?m)^\s*$([regex]::Escape($token))\s*$"))) {
        [void]$failures.Add("no '$token' in $($group.Name)")
      }
    }
  }

  foreach ($entry in $presenceEntries) {
    $file = Join-Path $RepositoryRoot ([string]$entry.file)
    if (-not (Test-Path -LiteralPath $file -PathType Leaf)) { [void]$failures.Add("$($entry.file) is missing"); continue }
    if (-not (Select-String -LiteralPath $file -Pattern ([regex]::Escape([string]$entry.anchor)) -Quiet)) {
      [void]$failures.Add("$($entry.file): anchor '$($entry.anchor)' not found")
    }
  }

  $counts = "entries=$total go_executed=$goExecuted pwsh_entries=$($pwshEntries.Count) pwsh_executed=$pwshExecuted presence=$($presenceEntries.Count) failures=$($failures.Count)"
  $passed = ($failures.Count -eq 0)
  return [pscustomobject]@{
    Passed       = $passed
    Reason       = if ($passed) { '' } else { 'a registered guard has no executing fixture. A guard that cannot demonstrate a catch is indistinguishable from an absent one.' }
    Counts       = $counts
    Entries      = $total
    GoExecuted   = $goExecuted
    PwshEntries  = $pwshEntries.Count
    PwshExecuted = $pwshExecuted
    Presence     = $presenceEntries.Count
    Failures     = @($failures)
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
  Remove-GateAnsi, Measure-GateGolangciLint, Measure-GateEslint, Measure-GateEslintRules, Measure-GatePrettier, Measure-GateArchscan, `
  Measure-GateBoundary, Measure-GateGovernance, Compare-GateRatchet, ConvertTo-GateCountMap, Measure-GateGuards, Test-GateGuardInventory
