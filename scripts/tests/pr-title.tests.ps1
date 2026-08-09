$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$checker = Join-Path $repoRoot 'scripts/check-pr-title.ps1'

function Invoke-TitleCheck([string]$Title) {
  & pwsh -NoProfile -ExecutionPolicy Bypass -File $checker -Title $Title *> $null
  return $LASTEXITCODE
}

# Titles this repository has actually merged, plus the shapes the convention
# names explicitly (breaking-change !, multi-segment scope).
$accepted = @(
  'feat(gate): typecheck every tracked source, census-guarded, per-file ratchet',
  'fix(catalogingest): reject empty GTIN',
  'docs(audit): record the verifier topology',
  'ci: add golangci-lint as a shrink-only ratchet',
  'chore(deps): bump vitest to 3.2.4',
  'feat!: drop the v1 pricing endpoint',
  'refactor(sdk-runtime): split the error contract module',
  'test(web): one vitest entry point over every workspace'
)
foreach ($title in $accepted) {
  if ((Invoke-TitleCheck $title) -ne 0) { throw "conventional title rejected: $title" }
}

# GATE-TOPOLOGY §2.3a names 'wip' as the negative fixture; the rest are the
# near-misses that make the regex load-bearing -- capitalized type, missing
# space, empty scope, empty subject, bare sentence.
$rejected = @(
  'wip',
  'WIP: finish the gate',
  'Fix bug',
  'feat:missing space after colon',
  'feat(): empty scope',
  'feat(gate):',
  'feature(gate): unknown type',
  'update stuff',
  ''
)
foreach ($title in $rejected) {
  if ((Invoke-TitleCheck $title) -eq 0) { throw "nonconventional title accepted: '$title'" }
}

Write-Output "PASS pr-title checker accepted=$($accepted.Count) rejected=$($rejected.Count)"
