$ErrorActionPreference = 'Stop'
$repoRoot = Split-Path -Parent $PSScriptRoot
Set-Location $repoRoot

function Fail([string]$message) { throw $message }
function Require-Command([string]$name) {
    if (-not (Get-Command $name -ErrorAction SilentlyContinue)) { Fail "required command not found: $name" }
}
function Test-DurableDocPath([string]$path) {
    $normalized = $path.Replace('\', '/')
    if ($normalized -eq 'README.md' -or $normalized -eq 'AGENTS.md' -or $normalized -eq 'ARCHITECTURE.md') { return $true }
    if (-not $normalized.StartsWith('docs/')) { return $false }
    foreach ($prefix in @('docs/work/', 'docs/archive/', 'docs/review/', 'docs/plans/', 'docs/superpowers/')) {
        if ($normalized.StartsWith($prefix)) { return $false }
    }
    return $normalized.EndsWith('.md')
}
function Resolve-RelativeDocLink([string]$sourcePath, [string]$target) {
    $cleanTarget = $target.Split('#')[0]
    if (-not $cleanTarget -or $cleanTarget.StartsWith('#') -or $cleanTarget.StartsWith('/') -or $cleanTarget -match '^[a-zA-Z][a-zA-Z0-9+.-]*:') { return $null }
    $sourceDir = Split-Path -Parent $sourcePath
    if (-not $sourceDir) { $sourceDir = '.' }
    $combined = [System.IO.Path]::GetFullPath((Join-Path $repoRoot (Join-Path $sourceDir $cleanTarget)))
    return [System.IO.Path]::GetRelativePath($repoRoot, $combined).Replace('\', '/')
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

Require-Command git
Require-Command node
Require-Command npm
Require-Command go

$trackedFiles = @(git ls-files)
if ($LASTEXITCODE -ne 0) { Fail 'git ls-files failed' }
$trackedSet = [System.Collections.Generic.HashSet[string]]::new([System.StringComparer]::OrdinalIgnoreCase)
foreach ($file in $trackedFiles) { [void]$trackedSet.Add($file.Replace('\', '/')) }

$headRef = $env:GATE_HEAD_REF
if ([string]::IsNullOrWhiteSpace($headRef)) { $headRef = (git branch --show-current).Trim() }
$isReview = $headRef -like 'review/*'
$trackedWork = @($trackedFiles | Where-Object { $_.Replace('\', '/').StartsWith('docs/work/') })
if ($isReview) {
    if ($trackedWork.Count -ne 1 -or $trackedWork[0].Replace('\', '/') -ne 'docs/work/current/ai-dialog.md') {
        Fail "review branch may track only docs/work/current/ai-dialog.md under docs/work/: $($trackedWork -join ', ')"
    }
} elseif ($trackedWork.Count -gt 0) {
    Fail "candidate/main contains temporary docs/work material: $($trackedWork -join ', ')"
}

$trackedAiDialog = @($trackedFiles | Where-Object { $_.Replace('\', '/') -match '(?i)(^|/)ai-dialog(?:\.md)?$' })
if (-not $isReview -and $trackedAiDialog.Count -gt 0) {
    Fail "candidate/main contains review transport: $($trackedAiDialog -join ', ')"
}

$requiredFiles = @(
    'README.md', 'AGENTS.md', 'ARCHITECTURE.md', 'docs/index.md', 'docs/roadmap.md',
    'docs/development/engineering-rules.md', 'contracts/api/product/openapi.yaml',
    'contracts/api/product/redocly.yaml', 'scripts/verify-product-oad.mjs', 'scripts/gate.ps1', 'package.json'
)
foreach ($path in $requiredFiles) {
    if (-not (Test-Path $path -PathType Leaf)) { Fail "required file missing: $path" }
}

$readme = Get-Content -Raw 'README.md'
$agent = Get-Content -Raw 'AGENTS.md'
$index = Get-Content -Raw 'docs/index.md'
$roadmap = Get-Content -Raw 'docs/roadmap.md'
$engineeringRules = Get-Content -Raw 'docs/development/engineering-rules.md'

$methodologyRepository = 'developmentconexus-ops/conexus-methodology'
$methodologyPin = '9c7210d1504bef01c0d134a6c3ae8627deebb535'
foreach ($token in @($methodologyRepository, $methodologyPin, 'ROUTER.md')) {
    if (-not $agent.Contains($token)) { Fail "AGENTS.md is missing canonical methodology routing token: $token" }
}
if ($agent.Contains("$methodologyRepository/blob/main")) { Fail 'AGENTS.md must not consume floating methodology main' }
foreach ($token in @($methodologyRepository, $methodologyPin)) {
    if (-not $engineeringRules.Contains($token)) { Fail "engineering-rules.md is missing canonical methodology token: $token" }
}

$bootstrapFiles = @('AGENTS.md', 'docs/index.md', 'docs/roadmap.md')
$bootstrapBytes = 0
foreach ($path in $bootstrapFiles) { $bootstrapBytes += (Get-Item $path).Length }
$bootstrapLimit = 20480
if ($bootstrapBytes -gt $bootstrapLimit) { Fail "bootstrap authority pack exceeds $bootstrapLimit bytes: $bootstrapBytes" }

if (-not $roadmap.Contains('<!-- program-status-authority -->')) { Fail 'docs/roadmap.md must own mutable program status' }
if (-not $roadmap.Contains('contracts/api/product/openapi.yaml')) { Fail 'roadmap lost canonical Product OAD route' }
foreach ($surface in @(@{name='README.md';text=$readme}, @{name='AGENTS.md';text=$agent}, @{name='docs/index.md';text=$index})) {
    if ($surface.text.Contains('<!-- program-status-authority -->') -or $surface.text.Contains('Exact next action')) {
        Fail "$($surface.name) duplicates mutable roadmap authority"
    }
}
if (-not $readme.Contains('[`AGENTS.md`](AGENTS.md)') -or -not $readme.Contains('[`docs/index.md`](docs/index.md)')) {
    Fail 'README.md must remain a landing page pointing to AGENTS.md and docs/index.md'
}

$openApiEntrypoints = @($trackedFiles | Where-Object { $_.Replace('\', '/') -match '(^|/)openapi[^/]*\.ya?ml$' })
if ($openApiEntrypoints.Count -ne 1 -or $openApiEntrypoints[0].Replace('\', '/') -ne 'contracts/api/product/openapi.yaml') {
    Fail "tracked OpenAPI entrypoint authority mismatch: $($openApiEntrypoints -join ', ')"
}

foreach ($pattern in @('^docs/superpowers/', '^docs/current/', '^docs/active/')) {
    foreach ($file in $trackedFiles) {
        if ($file.Replace('\', '/') -match $pattern) { Fail "forbidden tracked path: $file" }
    }
}

$durableDocs = @($trackedFiles | Where-Object { Test-DurableDocPath $_ })
$durableSet = [System.Collections.Generic.HashSet[string]]::new([System.StringComparer]::OrdinalIgnoreCase)
foreach ($path in $durableDocs) { [void]$durableSet.Add($path.Replace('\', '/')) }
$docLinks = @{}
foreach ($path in $durableDocs) {
    $normalizedPath = $path.Replace('\', '/')
    $text = Get-Content -Raw $normalizedPath
    $links = @()
    foreach ($match in [regex]::Matches($text, '\[[^\]]*\]\(([^)]+)\)')) {
        $resolved = Resolve-RelativeDocLink $normalizedPath $match.Groups[1].Value.Trim()
        if (-not $resolved) { continue }
        $links += $resolved
        if ($resolved.EndsWith('.md') -and -not $trackedSet.Contains($resolved)) {
            Fail "broken relative document link: $normalizedPath -> $resolved"
        }
    }
    $docLinks[$normalizedPath] = $links
}
$reachable = [System.Collections.Generic.HashSet[string]]::new([System.StringComparer]::OrdinalIgnoreCase)
$queue = [System.Collections.Generic.Queue[string]]::new()
foreach ($root in @('README.md', 'AGENTS.md', 'docs/index.md')) {
    if ($durableSet.Contains($root)) { [void]$reachable.Add($root); $queue.Enqueue($root) }
}
while ($queue.Count -gt 0) {
    $current = $queue.Dequeue()
    foreach ($link in $docLinks[$current]) {
        if ($durableSet.Contains($link) -and $reachable.Add($link)) { $queue.Enqueue($link) }
    }
}
foreach ($path in $durableDocs) {
    $normalizedPath = $path.Replace('\', '/')
    if (-not $reachable.Contains($normalizedPath)) { Fail "durable document is not reachable from bootstrap routers: $normalizedPath" }
}

$implementationBlocked = $roadmap -match '(?i)Implementation[^\n]*BLOCKED UNTIL D9'
if ($implementationBlocked) {
    foreach ($root in @('apps', 'cmd', 'internal', 'server', 'backend', 'frontend')) {
        $population = @($trackedFiles | Where-Object { $_.Replace('\', '/').StartsWith("$root/") })
        if ($population.Count -gt 0) { Fail "implementation is blocked by roadmap but tracked runtime exists under $root/" }
    }
}

$workflowFiles = @($trackedFiles | Where-Object { $_.Replace('\', '/').StartsWith('.github/workflows/') })
foreach ($workflow in $workflowFiles) {
    if ((Get-Content -Raw $workflow).Contains('pull_request_target')) { Fail "unsafe pull_request_target trigger found in $workflow" }
}

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
}
if ($implementationBlocked) {
    foreach ($file in $changedFiles) {
        $normalized = $file.Replace('\', '/')
        if ($normalized -match '^(apps|cmd|internal|server|backend|frontend)/') {
            Fail "implementation is blocked by roadmap but candidate changes runtime surface: $normalized"
        }
    }
}

$productProof = & node 'scripts/verify-product-oad.mjs' 2>&1
$productProofExit = $LASTEXITCODE
$productProof | ForEach-Object { Write-Host $_ }
if ($productProofExit -ne 0) { Fail 'Product OAD proof failed' }

Write-Host 'gate: PASS'
Write-Host "tracked_files: $($trackedFiles.Count)"
Write-Host "bootstrap_bytes: $bootstrapBytes/$bootstrapLimit"
Write-Host "durable_docs_reachable: $($reachable.Count)"
Write-Host "methodology_pin: $methodologyRepository@$methodologyPin"
Write-Host "implementation_blocked: $implementationBlocked"
Write-Host "diff_range: $diffRange changed_files: $($changedFiles.Count)"
Write-Host 'product_oad_proof: PASS'
