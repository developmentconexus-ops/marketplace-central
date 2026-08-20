#Requires -Version 7.0
<#
.SYNOPSIS
  Validates that a PR title is a Conventional Commit subject.

.DESCRIPTION
  With linear history and squash merges, the PR title
  becomes the commit subject, so validating the title is sufficient and
  validating every commit in the branch is not. The repository already writes
  Conventional Commits consistently; this makes the convention enforced instead
  of observed.

  The title arrives as a parameter, never interpolated into a shell command --
  a PR title is attacker-controlled text and the workflow passes it through an
  environment variable for exactly that reason.
#>
[CmdletBinding()]
param(
  [Parameter(Mandatory)][AllowEmptyString()][string]$Title
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

# The Conventional Commits v1.0.0 shape: type, optional (scope), optional ! for
# a breaking change, then ": " and a non-empty subject. The type list is the
# conventional set; 'style' is included for completeness even though Prettier
# and gofmt own formatting here.
$pattern = '^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)(\([a-z0-9][a-z0-9_./-]*\))?!?: \S'

if ($Title -match $pattern) {
  Write-Output "pr_title=conventional status=passed"
  exit 0
}
Write-Output "pr_title=nonconventional status=failed"
Write-Output "title_length=$($Title.Length)"
Write-Output 'expected: type(scope)?: subject -- e.g. "docs(api): author canonical Product OAD", types: feat fix docs style refactor perf test build ci chore revert'
exit 1
