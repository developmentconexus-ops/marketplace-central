$ErrorActionPreference = 'Stop'
$repoRoot = Split-Path -Parent $PSScriptRoot
Set-Location $repoRoot

function Fail([string]$message) { throw $message }
function Require-Command([string]$name) {
    if (-not (Get-Command $name -ErrorAction SilentlyContinue)) { Fail "required command not found: $name" }
}
function Resolve-GateBase {
    if ($env:GATE_BASE_SHA -and $env:GATE_BASE_SHA -notmatch '^0+$') {
        git cat-file -e "$($env:GATE_BASE_SHA)^{commit}" 2>$null
        if ($LASTEXITCODE -eq 0) { return $env:GATE_BASE_SHA }
    }
    foreach ($candidate in @('origin/main', 'main')) {
        git rev-parse --verify --quiet "$candidate^{commit}" *> $null
        if ($LASTEXITCODE -eq 0) { return $candidate }
    }
    return $null
}
function Test-ChangedPathMatches([string[]]$patterns) {
    foreach ($file in $changedFiles) {
        $normalized = $file.Replace('\', '/')
        foreach ($pattern in $patterns) {
            if ($normalized -match $pattern) { return $true }
        }
    }
    return $false
}

Require-Command git
Require-Command node
Require-Command npm
Require-Command go

$requiredFiles = @(
    'README.md',
    'AGENTS.md',
    'ARCHITECTURE.md',
    'docs/index.md',
    'docs/roadmap.md',
    'docs/development/engineering-method.md',
    'docs/development/frontend-product-experience-planning-method.md',
    'docs/development/engineering-rules.md',
    'contracts/api/product/openapi.yaml',
    'contracts/api/product/redocly.yaml',
    'scripts/verify-product-oad.mjs',
    'scripts/gate.ps1',
    'package.json'
)
foreach ($path in $requiredFiles) {
    if (-not (Test-Path $path -PathType Leaf)) { Fail "required file missing: $path" }
}

$trackedFiles = @(git ls-files)
if ($LASTEXITCODE -ne 0) { Fail 'git ls-files failed' }

$openApiEntrypoints = @($trackedFiles | Where-Object { $_.Replace('\', '/') -match '(^|/)openapi[^/]*\.ya?ml$' })
if ($openApiEntrypoints.Count -ne 1 -or $openApiEntrypoints[0].Replace('\', '/') -ne 'contracts/api/product/openapi.yaml') {
    Fail "tracked OpenAPI entrypoint authority mismatch: $($openApiEntrypoints -join ', ')"
}

$workflowFiles = @($trackedFiles | Where-Object { $_.Replace('\', '/').StartsWith('.github/workflows/') })
foreach ($workflow in $workflowFiles) {
    if ((Get-Content -Raw $workflow).Contains('pull_request_target')) { Fail "unsafe pull_request_target trigger found in $workflow" }
}

$roadmap = Get-Content -Raw 'docs/roadmap.md'
$implementationBlocked = $roadmap -match '(?i)\| Implementation \| \*\*BLOCKED'

$base = Resolve-GateBase
$head = if ($env:GATE_HEAD_SHA) { $env:GATE_HEAD_SHA } else { 'HEAD' }
$changedFiles = @()
$diffRange = 'unavailable'
if ($base) {
    git cat-file -e "$head^{commit}" 2>$null
    if ($LASTEXITCODE -ne 0) { Fail "gate head is not an available commit: $head" }

    $diffRange = "$base...$head"
    $changedFiles = @(git diff --name-only $diffRange)
    if ($LASTEXITCODE -ne 0) { Fail "git diff failed for range $diffRange" }

    $diffCheck = git diff --check $diffRange 2>&1
    $conflictMarkers = @($diffCheck | Select-String -SimpleMatch 'leftover conflict marker')
    if ($conflictMarkers.Count -gt 0) {
        Fail "candidate diff contains unresolved merge-conflict markers:`n$($conflictMarkers -join "`n")"
    }
}

if ($implementationBlocked) {
    foreach ($file in $changedFiles) {
        $normalized = $file.Replace('\', '/')
        if ($normalized -match '^(apps|cmd|internal|server|backend|frontend|src|migrations)/') {
            Fail "implementation is blocked by roadmap but candidate changes implementation surface: $normalized"
        }
    }
}

$productProofPatterns = @(
    '^contracts/api/product/',
    '^scripts/gate\.ps1$',
    '^scripts/verify-product-oad(?:-[^/]+)?\.mjs$',
    '^scripts/fixtures/product-oad-[^/]+\.json$',
    '^scripts/verify-oad-source-reachability\.mjs$',
    '^scripts/verify-operational-read-contract\.mjs$',
    '^scripts/verify-performance-evidence-knowledge\.mjs$',
    '^scripts/verify-notification-oad\.mjs$',
    '^scripts/verify-authorization-request-oad\.mjs$',
    '^scripts/lib/publication-requirements-oad-proof\.mjs$',
    '^package\.json$',
    '^\.node-version$',
    '^\.github/workflows/ci\.yml$'
)
$productProofAffected = if (-not $base) { $true } else { Test-ChangedPathMatches $productProofPatterns }

if ($productProofAffected) {
    $productProof = & node 'scripts/verify-product-oad.mjs' 2>&1
    $productProofExit = $LASTEXITCODE
    $productProof | ForEach-Object { Write-Host $_ }
    if ($productProofExit -ne 0) { Fail 'Product OAD proof failed' }
    $productProofStatus = 'PASS'
} else {
    $productProofStatus = 'SKIPPED_NOT_AFFECTED'
}

Write-Host 'gate: PASS'
Write-Host "required_files: $($requiredFiles.Count)"
Write-Host "implementation_blocked: $implementationBlocked"
Write-Host "diff_range: $diffRange changed_files: $($changedFiles.Count)"
Write-Host "product_oad_proof: $productProofStatus"
