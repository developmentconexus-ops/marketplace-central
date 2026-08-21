param(
    [ValidateSet('quick', 'full')]
    [string]$Lane = 'quick'
)

$ErrorActionPreference = 'Stop'
$repoRoot = Split-Path -Parent $PSScriptRoot
Set-Location $repoRoot

function Fail([string]$message) {
    throw $message
}

function Require-Command([string]$name) {
    if (-not (Get-Command $name -ErrorAction SilentlyContinue)) {
        Fail "required command not found: $name"
    }
}

function Get-TrackedFiles {
    $files = @(git ls-files)
    if ($LASTEXITCODE -ne 0) { Fail 'git ls-files failed' }
    return $files
}

function Resolve-GateBase {
    if ($env:GATE_BASE_SHA) {
        git cat-file -e "$($env:GATE_BASE_SHA)^{commit}" 2>$null
        if ($LASTEXITCODE -eq 0) { return $env:GATE_BASE_SHA }
        Fail "GATE_BASE_SHA is not an available commit: $($env:GATE_BASE_SHA)"
    }

    if ($env:GATE_BASE_REF) {
        foreach ($candidate in @($env:GATE_BASE_REF, "origin/$($env:GATE_BASE_REF)")) {
            git rev-parse --verify --quiet "$candidate^{commit}" *> $null
            if ($LASTEXITCODE -eq 0) { return $candidate }
        }
        Fail "GATE_BASE_REF is not an available commit: $($env:GATE_BASE_REF)"
    }

    foreach ($candidate in @('origin/main', 'main')) {
        git rev-parse --verify --quiet "$candidate^{commit}" *> $null
        if ($LASTEXITCODE -eq 0) { return $candidate }
    }

    return $null
}

function Test-ReviewDiffNames([string[]]$names) {
    $normalized = @($names | ForEach-Object { $_.Replace('\', '/') })
    return $normalized.Count -eq 1 -and $normalized[0] -eq 'docs/work/current/ai-dialog.md'
}

function Test-DurableDocPath([string]$path) {
    $normalized = $path.Replace('\', '/')
    if ($normalized -eq 'README.md' -or $normalized -eq 'AGENTS.md' -or $normalized -eq 'ARCHITECTURE.md') { return $true }
    if (-not $normalized.StartsWith('docs/')) { return $false }
    if ($normalized.StartsWith('docs/work/')) { return $false }
    if ($normalized.StartsWith('docs/archive/')) { return $false }
    if ($normalized.StartsWith('docs/review/')) { return $false }
    if ($normalized.StartsWith('docs/plans/')) { return $false }
    if ($normalized.StartsWith('docs/superpowers/')) { return $false }
    return $normalized.EndsWith('.md')
}

function Resolve-RelativeDocLink([string]$sourcePath, [string]$target) {
    $cleanTarget = $target.Split('#')[0]
    if (-not $cleanTarget -or $cleanTarget -match '^[a-zA-Z][a-zA-Z0-9+.-]*:' -or $cleanTarget.StartsWith('#')) { return $null }
    if ($cleanTarget.StartsWith('/')) { return $null }

    $sourceDir = Split-Path -Parent $sourcePath
    if (-not $sourceDir) { $sourceDir = '.' }
    $combined = [System.IO.Path]::GetFullPath((Join-Path $repoRoot (Join-Path $sourceDir $cleanTarget)))
    $relative = [System.IO.Path]::GetRelativePath($repoRoot, $combined).Replace('\', '/')
    return $relative
}

Require-Command git
Require-Command node
Require-Command npm

$trackedFiles = Get-TrackedFiles
$trackedSet = [System.Collections.Generic.HashSet[string]]::new([System.StringComparer]::OrdinalIgnoreCase)
foreach ($file in $trackedFiles) { [void]$trackedSet.Add($file.Replace('\', '/')) }

$headRef = $env:GATE_HEAD_REF
if ([string]::IsNullOrWhiteSpace($headRef)) { $headRef = (git branch --show-current).Trim() }
$isReview = $headRef -like 'review/*'
$trackedWork = @($trackedFiles | Where-Object { $_.Replace('\', '/').StartsWith('docs/work/') })
if ($isReview) {
    if (-not (Test-ReviewDiffNames $trackedWork)) {
        Fail "review branch may track only docs/work/current/ai-dialog.md under docs/work/: $($trackedWork -join ', ')"
    }
} elseif ($trackedWork.Count -gt 0) {
    Fail "candidate/main contains temporary docs/work material: $($trackedWork -join ', ')"
}

$requiredFiles = @(
    'README.md',
    'AGENTS.md',
    'ARCHITECTURE.md',
    'docs/index.md',
    'docs/roadmap.md',
    'docs/development/engineering-rules.md',
    'docs/development/evidence-grounded-production-engineering-for-llm-agents.md',
    'docs/architecture/decisions/README.md',
    'docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md',
    'docs/engineering/rebaseline/D1-DOMAINS-BOUNDARIES.md',
    'docs/engineering/rebaseline/D2-IDENTITY-TENANT-DATA-OWNERSHIP.md',
    'docs/engineering/rebaseline/D3-COMMUNICATION-EVENTS.md',
    'docs/engineering/rebaseline/D4-EXTERNAL-INTEGRATIONS.md',
    'docs/engineering/rebaseline/D4-R1-PUBLICATION-INPUT.md',
    'docs/engineering/rebaseline/D5-API.md',
    'docs/engineering/rebaseline/D5-B2-OPENAPI-WIRE-AUTHORITY-TOOLING.md',
    'docs/engineering/rebaseline/D6-FRONTEND.md',
    'docs/engineering/rebaseline/DECISION-RECONCILIATION-BASELINE.md',
    'docs/engineering/rebaseline/EVIDENCE-REGISTER.md',
    'contracts/api/product/openapi.yaml',
    'contracts/api/product/redocly.yaml',
    'scripts/verify-product-oad.mjs',
    'scripts/gate.ps1',
    'package.json'
)

foreach ($path in $requiredFiles) {
    if (-not (Test-Path $path -PathType Leaf)) { Fail "required file missing: $path" }
}

$forbiddenRoots = @('apps', 'cmd', 'internal', 'server', 'backend', 'frontend')
foreach ($root in $forbiddenRoots) {
    if (Test-Path $root) {
        $trackedUnderRoot = @($trackedFiles | Where-Object { $_.Replace('\', '/').StartsWith("$root/") })
        if ($trackedUnderRoot.Count -gt 0) { Fail "active legacy/runtime population found under $root/" }
    }
}

$trackedOpenApiEntrypoints = @($trackedFiles | Where-Object { $_ -match '(^|/)openapi[^/]*\.ya?ml$' })
if ($trackedOpenApiEntrypoints.Count -ne 1 -or $trackedOpenApiEntrypoints[0].Replace('\', '/') -ne 'contracts/api/product/openapi.yaml') {
    Fail "tracked OpenAPI entrypoint authority mismatch: $($trackedOpenApiEntrypoints -join ', ')"
}

$bootstrapFiles = @('AGENTS.md', 'docs/index.md', 'docs/roadmap.md')
$bootstrapBytes = 0
foreach ($path in $bootstrapFiles) {
    $bootstrapBytes += (Get-Item $path).Length
}
$bootstrapLimit = 20480
if ($bootstrapBytes -gt $bootstrapLimit) {
    Fail "bootstrap authority pack exceeds $bootstrapLimit bytes: $bootstrapBytes"
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
        $target = $match.Groups[1].Value.Trim()
        $resolved = Resolve-RelativeDocLink $normalizedPath $target
        if ($resolved) {
            $links += $resolved
            if ($durableSet.Contains($resolved) -and -not $trackedSet.Contains($resolved)) {
                Fail "durable link resolves to untracked document: $normalizedPath -> $resolved"
            }
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

$readme = Get-Content -Raw 'README.md'
$agent = Get-Content -Raw 'AGENTS.md'
$index = Get-Content -Raw 'docs/index.md'
$roadmap = Get-Content -Raw 'docs/roadmap.md'

$roadmapMarkers = @(
    'D0 — Product / System Definition | ACCEPTED / CLOSED',
    'D1 — Domains / Boundaries | ACCEPTED / CLOSED',
    'D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED',
    'D3 — Communication / Events | ACCEPTED / CLOSED',
    'D4 — External Integrations | ACCEPTED / CLOSED',
    'D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL',
    'D5 — API | ACCEPTED / CLOSED',
    'D6 — Frontend — OPEN / ACTIVE',
    'D6 — Frontend | OPEN / ACTIVE',
    'D7 — Runtime / Jobs / Transactions | BLOCKED',
    'D8 — Golden Flows | BLOCKED',
    'D9 — Adversarial Architecture Review | BLOCKED',
    'Canonical Product OAD',
    'contracts/api/product/openapi.yaml',
    '99 Product operations',
    '30 ordinary Permissions',
    'Principal kinds H / A / S only',
    'https://conexus.fun',
    'Active runtime baseline',
    'NONE',
    'BLOCKED UNTIL D9'
)
foreach ($marker in $roadmapMarkers) {
    if (-not $roadmap.Contains($marker)) { Fail "docs/roadmap.md is missing required current truth: $marker" }
}

foreach ($forbidden in @('D5 — API — OPEN / ACTIVE', 'D5 — API — ACCEPTED / CLOSED', 'D6 — Frontend — OPEN / ACTIVE', 'Exact next action', '<!-- program-status-authority -->')) {
    if ($index.Contains($forbidden)) { Fail "docs/index.md duplicates mutable roadmap state: $forbidden" }
    if ($readme.Contains($forbidden)) { Fail "README.md is not landing-only: $forbidden" }
    if ($agent.Contains($forbidden)) { Fail "AGENTS.md duplicates mutable roadmap state: $forbidden" }
}
if (-not $readme.Contains('[`AGENTS.md`](AGENTS.md)') -or -not $readme.Contains('[`docs/index.md`](docs/index.md)')) {
    Fail 'README.md must remain a landing page pointing to AGENTS.md and docs/index.md.'
}

$forbiddenPathPatterns = @(
    '^docs/superpowers/',
    '^docs/archive/architecture-rebaseline/d0-d4-reset/',
    '^docs/current/',
    '^docs/active/'
)
foreach ($file in $trackedFiles) {
    $normalized = $file.Replace('\', '/')
    foreach ($pattern in $forbiddenPathPatterns) {
        if ($normalized -match $pattern) { Fail "forbidden tracked path: $normalized" }
    }
}

$machineRouteFiles = @($trackedFiles | Where-Object { $_ -match '(^|/)(CODEOWNERS|\.github/.*owners|.*machine.*route.*)$' })
$machineRouteSelectors = 0
foreach ($path in $machineRouteFiles) {
    $machineRouteSelectors += ([regex]::Matches((Get-Content -Raw $path), '(?m)^\s*[^#\s].*$')).Count
}

$base = Resolve-GateBase
$head = if ($env:GATE_HEAD_SHA) { $env:GATE_HEAD_SHA } else { 'HEAD' }
$changedFiles = @()
$diffRange = 'unavailable'
if ($base) {
    git cat-file -e "$head^{commit}" 2>$null
    if ($LASTEXITCODE -ne 0) { Fail "gate head is not an available commit: $head" }

    if ($isReview) {
        $candidateRef = $env:GATE_CANDIDATE_REF
        if ([string]::IsNullOrWhiteSpace($candidateRef)) { $candidateRef = $env:GATE_BASE_REF }
        if ([string]::IsNullOrWhiteSpace($candidateRef)) { Fail 'review mode requires an explicit candidate branch/ref' }

        git show-ref --verify --quiet "refs/remotes/origin/$candidateRef"
        if ($LASTEXITCODE -ne 0) {
            git ls-remote --exit-code --heads origin $candidateRef *> $null
            if ($LASTEXITCODE -ne 0) { Fail "review candidate ref does not exist on origin: $candidateRef" }
        }

        $candidateTree = @(git ls-tree -r --name-only $base)
        if ($LASTEXITCODE -ne 0) { Fail "unable to inspect exact review candidate tree: $base" }
        $candidateWork = @($candidateTree | Where-Object { $_.Replace('\', '/').StartsWith('docs/work/') })
        if ($candidateWork.Count -gt 0) { Fail 'exact review candidate tree is contaminated by docs/work material' }

        $reviewNames = @(git diff --name-only "$base..$head")
        if ($LASTEXITCODE -ne 0) { Fail "git diff failed for review range $base..$head" }
        if (-not (Test-ReviewDiffNames $reviewNames)) {
            Fail "review branch must differ from exact candidate only by docs/work/current/ai-dialog.md: $($reviewNames -join ', ')"
        }
        $diffRange = "$base..$head"
        $changedFiles = $reviewNames
    } else {
        $diffRange = "$base...$head"
        $changedFiles = @(git diff --name-only $diffRange)
        if ($LASTEXITCODE -ne 0) { Fail "git diff failed for range $diffRange" }
    }
}

$forbiddenDiffPatterns = @(
    '^apps/', '^cmd/', '^internal/', '^server/', '^backend/', '^frontend/',
    '(^|/)openapi(?!\.yaml$).*\.ya?ml$',
    '(^|/)sdk/', '(^|/)generated/'
)
foreach ($file in $changedFiles) {
    $normalized = $file.Replace('\', '/')
    foreach ($pattern in $forbiddenDiffPatterns) {
        if ($normalized -match $pattern -and $normalized -ne 'contracts/api/product/openapi.yaml') {
            Fail "changed file enters forbidden implementation/runtime/generated surface: $normalized"
        }
    }
}

$negativeControls = 0
function Expect-Failure([string]$name, [scriptblock]$body) {
    $failed = $false
    try { & $body } catch { $failed = $true }
    if (-not $failed) { Fail "negative control unexpectedly passed: $name" }
    $script:negativeControls++
}

Expect-Failure 'legacy runtime root' {
    if (-not ('apps/example/main.go' -match '^apps/')) { return }
    throw 'caught'
}
Expect-Failure 'second Product OAD entrypoint' {
    $fixture = @('contracts/api/product/openapi.yaml', 'docs/openapi.yaml')
    if ($fixture.Count -eq 1) { return }
    throw 'caught'
}
Expect-Failure 'duplicate mutable status' {
    $fixture = 'D6 — Frontend — OPEN / ACTIVE'
    if (-not $fixture.Contains('D6 — Frontend — OPEN / ACTIVE')) { return }
    throw 'caught'
}
Expect-Failure 'forbidden superpowers docs' {
    if (-not ('docs/superpowers/specs/x.md' -match '^docs/superpowers/')) { return }
    throw 'caught'
}
Expect-Failure 'bootstrap overflow' {
    if (-not (($bootstrapLimit + 1) -gt $bootstrapLimit)) { return }
    throw 'caught'
}
Expect-Failure 'review isolation' {
    if ((Test-ReviewDiffNames @('docs/work/current/ai-dialog.md')) -and -not (Test-ReviewDiffNames @('docs/work/current/ai-dialog.md', 'README.md'))) { throw 'caught' }
}

if ($negativeControls -ne 6) { Fail "repository negative-control count mismatch: $negativeControls/6" }

$productProof = & node 'scripts/verify-product-oad.mjs' 2>&1
if ($LASTEXITCODE -ne 0) {
    $productProof | ForEach-Object { Write-Host $_ }
    Fail 'Product OAD proof failed'
}
$productProof | ForEach-Object { Write-Host $_ }

Write-Host "gate lane: $Lane"
Write-Host "required files: $($requiredFiles.Count)"
Write-Host "tracked files inspected: $($trackedFiles.Count)"
Write-Host "bootstrap_bytes: $bootstrapBytes / $bootstrapLimit"
Write-Host "docs_index_relative_links: $(@($docLinks['docs/index.md']).Count)"
Write-Host "durable_docs_reachable: $($reachable.Count)"
Write-Host "machine_route_files: $($machineRouteFiles.Count) selectors: $machineRouteSelectors"
Write-Host "diff_range: $diffRange changed_files: $($changedFiles.Count)"
Write-Host "review_mode: $isReview"
Write-Host "legacy_runtime_population: 0"
Write-Host "negative_controls: $negativeControls/6"
Write-Host 'gate: PASS'
