#Requires -Version 7.0
$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$agentsPath = Join-Path $repoRoot 'AGENTS.md'
$docsIndexPath = Join-Path $repoRoot 'docs/README.md'
$rootReadmePath = Join-Path $repoRoot 'README.md'
$compatibilityPath = Join-Path $repoRoot 'docs/engineering/rebaseline/README.md'
$knowledgeRoutesPath = Join-Path $repoRoot 'contracts/governance/knowledge-routes.json'
$bootstrapBudgetBytes = 20KB

function Assert-True {
  param([bool]$Condition, [string]$Message)
  if (-not $Condition) { throw $Message }
}

function Test-ExactPaths {
  param([string[]]$Actual, [string[]]$Expected)
  return (($Actual | Sort-Object) -join '|') -eq (($Expected | Sort-Object) -join '|')
}

function Test-WithinBudget {
  param([long]$Bytes, [long]$Budget)
  return $Bytes -le $Budget
}

# Positive and negative controls for the two predicates that protect the
# bootstrap. A comparison that can only return true is not a control.
Assert-True (Test-ExactPaths -Actual @('AGENTS.md', 'docs/README.md') -Expected @('docs/README.md', 'AGENTS.md')) `
  'exact-path predicate rejected the intended two-file bootstrap'
Assert-True (-not (Test-ExactPaths -Actual @('AGENTS.md', 'docs/README.md', 'ARCHITECTURE.md') -Expected @('AGENTS.md', 'docs/README.md'))) `
  'exact-path predicate accepted an extra default-read authority'
Assert-True (Test-WithinBudget -Bytes $bootstrapBudgetBytes -Budget $bootstrapBudgetBytes) `
  'bootstrap budget predicate rejected its inclusive boundary'
Assert-True (-not (Test-WithinBudget -Bytes ($bootstrapBudgetBytes + 1) -Budget $bootstrapBudgetBytes)) `
  'bootstrap budget predicate accepted an oversized bootstrap'

foreach ($requiredPath in @($agentsPath, $docsIndexPath, $rootReadmePath, $compatibilityPath, $knowledgeRoutesPath)) {
  Assert-True (Test-Path -LiteralPath $requiredPath -PathType Leaf) "missing documentation authority file: $requiredPath"
}

$agents = Get-Content -Raw -LiteralPath $agentsPath
$docsIndex = Get-Content -Raw -LiteralPath $docsIndexPath
$rootReadme = Get-Content -Raw -LiteralPath $rootReadmePath
$compatibility = Get-Content -Raw -LiteralPath $compatibilityPath

$bootstrapBytes = (Get-Item -LiteralPath $agentsPath).Length + (Get-Item -LiteralPath $docsIndexPath).Length
Assert-True (Test-WithinBudget -Bytes $bootstrapBytes -Budget $bootstrapBudgetBytes) `
  "default bootstrap exceeds $bootstrapBudgetBytes bytes: actual=$bootstrapBytes"

Assert-True ($agents.Contains('DEFAULT READ: `AGENTS.md` → `docs/README.md`. STOP.')) `
  'AGENTS.md does not state the exact two-file default read'
Assert-True ($agents.Contains('sole authority for current program status')) `
  'AGENTS.md does not delegate status exclusively to docs/README.md'
Assert-True ($agents.Contains('one PR per D-stage')) `
  'AGENTS.md lost the one-PR-per-D-stage workflow'
Assert-True ($agents.Contains('BUILD NOW') -and $agents.Contains('SEAM NOW') -and $agents.Contains('PROVE FIRST') -and $agents.Contains('DEFER')) `
  'AGENTS.md lost staged-delivery classification'
Assert-True ($agents.Contains('`AI-DIALOG.md` is temporary')) `
  'AGENTS.md does not make AI-DIALOG temporary'
Assert-True ($agents.Contains('superseded legacy topology') -and $agents.Contains('Never silence current security')) `
  'AGENTS.md does not distinguish scoped retirement of legacy-only checks from weakening current safety controls'
Assert-True (-not $agents.Contains('docs/engineering/rebaseline/README.md')) `
  'AGENTS.md still routes the default session through the retired status path'

Assert-True ($docsIndex.Contains('<!-- program-status-authority -->')) `
  'docs/README.md lacks the unique program-status authority marker'
Assert-True ($docsIndex.Contains('D5 — API — OPEN / ACTIVE')) `
  'docs/README.md lost the current D5 stage'
Assert-True ($docsIndex.Contains('Author and prove the canonical Product OpenAPI Description')) `
  'docs/README.md lost the exact next action'
Assert-True ($docsIndex.Contains('95 operations · 29 ordinary Permissions · Principal kinds H/A/S only')) `
  'docs/README.md lost the accepted Product surface'
Assert-True ($docsIndex.Contains('https://conexus.fun/marketplace-central/problems/product/{slug}')) `
  'docs/README.md lost the canonical Product Problem namespace'
Assert-True ($docsIndex.Contains('BLOCKED until D9 is accepted')) `
  'docs/README.md lost the implementation block'
Assert-True ($docsIndex.Contains('Do not recursively read every D-stage')) `
  'docs/README.md lost selective-read discipline'

$statusAuthorityFiles = @(& git -C $repoRoot grep -l --fixed-strings '<!-- program-status-authority -->' -- '*.md' 2>$null)
Assert-True (Test-ExactPaths -Actual $statusAuthorityFiles -Expected @('docs/README.md')) `
  "program-status marker is not unique: $($statusAuthorityFiles -join ', ')"

$knowledgeRoutes = Get-Content -Raw -LiteralPath $knowledgeRoutesPath | ConvertFrom-Json
$rootBootstrapRoutes = @($knowledgeRoutes.routes | Where-Object id -eq 'root-bootstrap')
Assert-True ($rootBootstrapRoutes.Count -eq 1) `
  "knowledge-routes must contain exactly one root-bootstrap route, found=$($rootBootstrapRoutes.Count)"
$rootBootstrapPaths = @($rootBootstrapRoutes[0].selectors | ForEach-Object path)
Assert-True (Test-ExactPaths -Actual $rootBootstrapPaths -Expected @('AGENTS.md', 'docs/README.md')) `
  "knowledge-routes root-bootstrap is not exact: $($rootBootstrapPaths -join ', ')"

Assert-True ($rootReadme.Contains('[`AGENTS.md`](AGENTS.md)') -and $rootReadme.Contains('[`docs/README.md`](docs/README.md)')) `
  'root README does not route to the two-file bootstrap'
Assert-True (-not $rootReadme.Contains('docs/engineering/rebaseline/README.md')) `
  'root README still routes through the retired status path'

Assert-True ((Get-Item -LiteralPath $compatibilityPath).Length -le 1200) `
  'compatibility pointer grew into another router'
Assert-True ($compatibility.Contains('Compatibility pointer only')) `
  'old rebaseline README is not explicitly a compatibility pointer'
Assert-True ($compatibility.Contains('docs/README.md')) `
  'compatibility pointer does not route to docs/README.md'
Assert-True (-not $compatibility.Contains('<!-- program-status-authority -->')) `
  'compatibility pointer claims program-status authority'

$forbiddenPaths = @(
  'AI-DIALOG.md',
  'docs/superpowers',
  'docs/engineering/rebaseline/cockpit.html',
  'scripts/tests/rebaseline-cockpit.test.mjs',
  'docs/engineering/standards/evidence-grounded-production-engineering-for-llm-agents.md'
)
foreach ($relativePath in $forbiddenPaths) {
  Assert-True (-not (Test-Path -LiteralPath (Join-Path $repoRoot $relativePath))) `
    "temporary/duplicate documentation surface remains active: $relativePath"
}

# Every relative Markdown document link in the index must resolve. Future
# artifacts such as contracts/api/product/openapi.yaml are deliberately code
# literals rather than links until they exist.
$docsDirectory = Split-Path -Parent $docsIndexPath
$linkMatches = [regex]::Matches($docsIndex, '\]\((?!https?://|#)([^)]+\.md)(?:#[^)]+)?\)')
Assert-True ($linkMatches.Count -gt 0) 'docs/README.md contains no routed Markdown documents'
foreach ($match in $linkMatches) {
  $target = $match.Groups[1].Value
  $resolvedTarget = [IO.Path]::GetFullPath((Join-Path $docsDirectory $target))
  Assert-True (Test-Path -LiteralPath $resolvedTarget -PathType Leaf) `
    "docs/README.md contains a dead document link: $target"
}

Write-Output "PASS documentation authority bootstrap_bytes=$bootstrapBytes links=$($linkMatches.Count) root_routes=$($rootBootstrapPaths.Count)"
