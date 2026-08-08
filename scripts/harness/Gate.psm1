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

Export-ModuleMember -Function Measure-GateGoTest, Measure-GateVitest, Measure-GateTsc, Measure-GateGofmt, Remove-GateAnsi
