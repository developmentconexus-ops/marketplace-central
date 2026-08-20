#Requires -Version 7.0
$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$agentsPath = Join-Path $repoRoot 'AGENTS.md'
$docsIndexPath = Join-Path $repoRoot 'docs/README.md'
$rootReadmePath = Join-Path $repoRoot 'README.md'
$compatibilityRelativePath = 'docs/engineering/' + 'rebaseline/README.md'
$compatibilityPath = Join-Path $repoRoot $compatibilityRelativePath
$reconciliationPath = Join-Path $repoRoot 'docs/engineering/rebaseline/DECISION-RECONCILIATION-BASELINE.md'
$adrRegistryPath = Join-Path $repoRoot 'docs/architecture/decisions/README.md'
$knowledgeRoutesPath = Join-Path $repoRoot 'contracts/governance/knowledge-routes.json'
$dialogPath = Join-Path $repoRoot 'AI-DIALOG.md'
$bootstrapBudgetBytes = 20KB
$oldStatusRouter = 'docs/engineering/' + 'rebaseline/README.md'

function Assert-True {
  param([bool]$Condition, [string]$Message)
  if (-not $Condition) { throw $Message }
}

function Test-ExactPaths {
  param([string[]]$Actual, [string[]]$Expected)
  return ((@($Actual) | Sort-Object) -join '|') -eq ((@($Expected) | Sort-Object) -join '|')
}

function Test-WithinBudget {
  param([long]$Bytes, [long]$Budget)
  return $Bytes -le $Budget
}

function Test-NoForbiddenPaths {
  param([string[]]$TrackedPaths, [string[]]$ForbiddenPaths)
  foreach ($forbiddenPath in $ForbiddenPaths) {
    $hit = @($TrackedPaths | Where-Object {
      $_ -eq $forbiddenPath -or $_.StartsWith("$forbiddenPath/", [StringComparison]::Ordinal)
    })
    if ($hit.Count -gt 0) { return $false }
  }
  return $true
}

function Get-CurrentBranchName {
  foreach ($candidate in @($env:GITHUB_HEAD_REF, $env:GITHUB_REF_NAME)) {
    if (-not [string]::IsNullOrWhiteSpace($candidate)) { return $candidate.Trim() }
  }
  $branch = [string](& git -C $repoRoot branch --show-current 2>$null)
  return $branch.Trim()
}

function Resolve-DocumentationCandidateRef {
  if (-not (Test-Path -LiteralPath $dialogPath -PathType Leaf)) { return 'HEAD' }

  $branch = Get-CurrentBranchName
  Assert-True ($branch -like 'review/*') `
    "AI-DIALOG.md may exist only on an isolated review/* branch, current=$branch"

  $dialog = Get-Content -Raw -LiteralPath $dialogPath
  $candidateMatch = [regex]::Match($dialog, '(?m)^candidate head\s+([0-9a-f]{40})\s*$')
  Assert-True $candidateMatch.Success 'AI-DIALOG.md does not declare an exact candidate head SHA'
  $candidateRef = $candidateMatch.Groups[1].Value

  & git -C $repoRoot cat-file -e "$candidateRef^{commit}" 2>$null
  Assert-True ($LASTEXITCODE -eq 0) "AI-DIALOG.md candidate head is unavailable: $candidateRef"

  $reviewChanges = @(& git -C $repoRoot diff --name-only "$candidateRef..HEAD" 2>$null)
  Assert-True (Test-ExactPaths -Actual $reviewChanges -Expected @('AI-DIALOG.md')) `
    "review branch contains changes outside AI-DIALOG.md: $($reviewChanges -join ', ')"

  return $candidateRef
}

function Get-TrackedPathsAtRef {
  param([string]$Ref)
  $paths = @(& git -C $repoRoot ls-tree -r --name-only $Ref 2>$null)
  Assert-True ($LASTEXITCODE -eq 0) "cannot enumerate candidate tree: $Ref"
  Assert-True ($paths.Count -gt 0) "candidate tree is empty: $Ref"
  return $paths
}

# Positive and negative controls for the predicates that protect the bootstrap
# and merge candidate. A comparison that can only return true is not a control.
Assert-True (Test-ExactPaths -Actual @('AGENTS.md', 'docs/README.md') -Expected @('docs/README.md', 'AGENTS.md')) `
  'exact-path predicate rejected the intended two-file bootstrap'
Assert-True (-not (Test-ExactPaths -Actual @('AGENTS.md', 'docs/README.md', 'ARCHITECTURE.md') -Expected @('AGENTS.md', 'docs/README.md'))) `
  'exact-path predicate accepted an extra default-read authority'
Assert-True (Test-WithinBudget -Bytes $bootstrapBudgetBytes -Budget $bootstrapBudgetBytes) `
  'bootstrap budget predicate rejected its inclusive boundary'
Assert-True (-not (Test-WithinBudget -Bytes ($bootstrapBudgetBytes + 1) -Budget $bootstrapBudgetBytes)) `
  'bootstrap budget predicate accepted an oversized bootstrap'
Assert-True (Test-NoForbiddenPaths -TrackedPaths @('AGENTS.md', 'docs/README.md') -ForbiddenPaths @('AI-DIALOG.md', 'docs/superpowers')) `
  'forbidden-path predicate rejected a clean candidate fixture'
Assert-True (-not (Test-NoForbiddenPaths -TrackedPaths @('AGENTS.md', 'AI-DIALOG.md') -ForbiddenPaths @('AI-DIALOG.md'))) `
  'forbidden-path predicate accepted AI-DIALOG.md in a candidate fixture'
Assert-True (-not (Test-NoForbiddenPaths -TrackedPaths @('docs/superpowers/specs/a.md') -ForbiddenPaths @('docs/superpowers'))) `
  'forbidden-path predicate accepted a nested forbidden directory fixture'

foreach ($requiredPath in @($agentsPath, $docsIndexPath, $rootReadmePath, $reconciliationPath, $adrRegistryPath, $knowledgeRoutesPath)) {
  Assert-True (Test-Path -LiteralPath $requiredPath -PathType Leaf) "missing documentation authority file: $requiredPath"
}

$agents = Get-Content -Raw -LiteralPath $agentsPath
$docsIndex = Get-Content -Raw -LiteralPath $docsIndexPath
$rootReadme = Get-Content -Raw -LiteralPath $rootReadmePath
$reconciliation = Get-Content -Raw -LiteralPath $reconciliationPath
$adrRegistry = Get-Content -Raw -LiteralPath $adrRegistryPath

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
Assert-True ($agents.Contains('subject population is zero') -and $agents.Contains('full replacement coverage')) `
  'AGENTS.md permits legacy-check retirement without attributable zero-population or replacement-coverage evidence'
Assert-True ($agents.Contains('Never silence current security')) `
  'AGENTS.md lost the protected current-safety control classes'
Assert-True ($agents.Contains('Dependency or lockfile changes require explicit declared scope')) `
  'AGENTS.md lost explicit dependency-change scope'
Assert-True ($agents.Contains('Do not alter accepted target architecture merely to preserve legacy code or tests')) `
  'AGENTS.md permits legacy implementation to distort accepted target architecture'
Assert-True (-not $agents.Contains($oldStatusRouter)) `
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
Assert-True ($docsIndex.Contains('Stage-local status and read-order prose accepted before this checkpoint is superseded by this index')) `
  'docs/README.md does not supersede historical stage-local routing prose'
Assert-True ($docsIndex.Contains('non-authoritative local availability copy') -and $docsIndex.Contains('developmentconexus-ops/conexus-methodology/METHOD.md') -and $docsIndex.Contains('v1.0.0')) `
  'docs/README.md lost canonical Method provenance for the local availability copy'
Assert-True ($docsIndex.Contains('| `README.md` | SUPPORTING |')) `
  'docs/README.md does not assign the root landing page a documentation disposition'

$statusAuthorityFiles = @(& git -C $repoRoot grep -l --fixed-strings '<!-- program-status-authority -->' -- '*.md' 2>$null)
Assert-True (Test-ExactPaths -Actual $statusAuthorityFiles -Expected @('docs/README.md')) `
  "program-status marker is not unique: $($statusAuthorityFiles -join ', ')"

$allowedHistoricalRouterReferences = @(
  'docs/architecture/decisions/035-architecture-rebaseline-governs-target-design.md',
  'docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md',
  'docs/engineering/rebaseline/D4-R1-PUBLICATION-INPUT.md',
  'docs/engineering/rebaseline/D5-B2-OPENAPI-WIRE-AUTHORITY-TOOLING.md'
)
$routerReferenceFiles = @(& git -C $repoRoot grep -l --fixed-strings $oldStatusRouter -- '*.md' '*.json' '*.ps1' '*.yml' '*.yaml' 2>$null)
Assert-True (Test-ExactPaths -Actual $routerReferenceFiles -Expected $allowedHistoricalRouterReferences) `
  "retired status-router references are not limited to explicitly superseded historical prose: $($routerReferenceFiles -join ', ')"
Assert-True ($adrRegistry.Contains('2026-08-20 authority-surface amendment') -and $adrRegistry.Contains('docs/README.md')) `
  'ADR registry does not explicitly amend ADR-035 status-routing consequences'

$knowledgeRoutes = Get-Content -Raw -LiteralPath $knowledgeRoutesPath | ConvertFrom-Json
$rootBootstrapRoutes = @($knowledgeRoutes.routes | Where-Object id -eq 'root-bootstrap')
Assert-True ($rootBootstrapRoutes.Count -eq 1) `
  "knowledge-routes must contain exactly one root-bootstrap route, found=$($rootBootstrapRoutes.Count)"
$rootBootstrapPaths = @($rootBootstrapRoutes[0].selectors | ForEach-Object path)
Assert-True (Test-ExactPaths -Actual $rootBootstrapPaths -Expected @('AGENTS.md', 'docs/README.md')) `
  "knowledge-routes root-bootstrap is not exact: $($rootBootstrapPaths -join ', ')"
$allRouteSelectors = @($knowledgeRoutes.routes | ForEach-Object selectors)
foreach ($selector in $allRouteSelectors) {
  Assert-True (Test-Path -LiteralPath (Join-Path $repoRoot $selector.path) -PathType Leaf) `
    "knowledge-routes selector path does not exist: $($selector.path)"
}
Assert-True (@($allRouteSelectors.path | Where-Object { $_ -eq 'docs/engineering/standards/root-cause-global-maximum-method.md' }).Count -eq 1) `
  'knowledge-routes does not route tooling to the local Method availability copy'

Assert-True ($rootReadme.Contains('Landing page only')) `
  'root README is not explicitly a pure landing-page pointer'
Assert-True ($rootReadme.Contains('[`AGENTS.md`](AGENTS.md)') -and $rootReadme.Contains('[`docs/README.md`](docs/README.md)')) `
  'root README does not route to the two-file bootstrap'
Assert-True (-not $rootReadme.Contains($oldStatusRouter)) `
  'root README still routes through the retired status path'
foreach ($forbiddenStatusPhrase in @('Current posture', 'Architecture Rebaseline D0→D9', 'blocked until D9', 'exact next action', 'D5 — API')) {
  Assert-True (-not $rootReadme.Contains($forbiddenStatusPhrase)) `
    "root README carries independent program-status prose: $forbiddenStatusPhrase"
}

if (Test-Path -LiteralPath $compatibilityPath -PathType Leaf) {
  $compatibility = Get-Content -Raw -LiteralPath $compatibilityPath
  Assert-True ((Get-Item -LiteralPath $compatibilityPath).Length -le 1200) `
    'compatibility pointer grew into another router'
  Assert-True ($compatibility.Contains('Compatibility pointer only')) `
    'old rebaseline README is not explicitly a compatibility pointer'
  Assert-True ($compatibility.Contains('docs/README.md')) `
    'compatibility pointer does not route to docs/README.md'
  Assert-True (-not $compatibility.Contains('<!-- program-status-authority -->')) `
    'compatibility pointer claims program-status authority'
} else {
  Assert-True ($routerReferenceFiles.Count -eq 0) `
    "compatibility pointer was retired while tracked references remain: $($routerReferenceFiles -join ', ')"
}

Assert-True ($reconciliation.Contains('on-demand routing map')) `
  'Decision Reconciliation Baseline is not explicitly on-demand'
Assert-True ($reconciliation.Contains('Do **not** read this file during the default fresh-session bootstrap')) `
  'Decision Reconciliation Baseline still behaves as default bootstrap'
Assert-True (-not $reconciliation.Contains('Always read:')) `
  'Decision Reconciliation Baseline still mandates the old recursive read shape'
Assert-True ($reconciliation.Contains('Exact program status and the next action live only in `docs/README.md`')) `
  'Decision Reconciliation Baseline claims status ownership outside docs/README.md'

$forbiddenCandidatePaths = @(
  'AI-DIALOG.md',
  'docs/superpowers',
  'docs/engineering/rebaseline/cockpit.html',
  'scripts/tests/rebaseline-cockpit.test.mjs',
  'docs/engineering/standards/evidence-grounded-production-engineering-for-llm-agents.md'
)
$candidateRef = Resolve-DocumentationCandidateRef
$candidatePaths = Get-TrackedPathsAtRef $candidateRef
Assert-True (Test-NoForbiddenPaths -TrackedPaths $candidatePaths -ForbiddenPaths $forbiddenCandidatePaths) `
  "temporary/duplicate documentation surface remains in merge candidate $candidateRef"

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

Write-Output "PASS documentation authority bootstrap_bytes=$bootstrapBytes links=$($linkMatches.Count) root_routes=$($rootBootstrapPaths.Count) candidate_ref=$candidateRef"
