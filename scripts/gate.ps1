[CmdletBinding()]
param(
    [ValidateSet('all', 'full')]
    [string]$Lane = 'all'
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$root = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path

function Fail([string]$Message) {
    Write-Error $Message
    exit 1
}

function Repo-Path([string]$Path) {
    $normalized = $Path.Replace('\', '/')
    while ($normalized.StartsWith('./', [StringComparison]::Ordinal)) {
        $normalized = $normalized.Substring(2)
    }
    return $normalized.TrimStart('/')
}

function Test-BootstrapBudget([long]$Bytes) {
    return $Bytes -le 20480
}

$allowedRoot = @(
    '.gitattributes',
    '.gitignore',
    '.node-version',
    'AGENTS.md',
    'ARCHITECTURE.md',
    'CLAUDE.md',
    'README.md',
    'package.json',
    'package-lock.json'
)

$allowedPrefixes = @(
    '.claude/',
    '.codex/',
    '.github/workflows/',
    'docs/',
    'scripts/',
    'contracts/api/product/',
    'qualification/'
)

function Test-AllowedTrackedPath([string]$Path) {
    $normalized = Repo-Path $Path
    if ($allowedRoot -contains $normalized) { return $true }
    foreach ($prefix in $allowedPrefixes) {
        if ($normalized.StartsWith($prefix, [StringComparison]::Ordinal)) { return $true }
    }
    return $false
}

function Test-ReviewDiffNames([string[]]$Names) {
    return $Names.Count -eq 1 -and (Repo-Path $Names[0]) -eq 'docs/work/current/ai-dialog.md'
}

function Get-RelativeMarkdownLinks([string]$FilePath) {
    $text = Get-Content -LiteralPath $FilePath -Raw
    $links = @()
    foreach ($match in [regex]::Matches($text, '\[[^\]]+\]\(([^)]+)\)')) {
        $href = $match.Groups[1].Value.Trim()
        if (-not $href -or $href.StartsWith('#') -or $href -match '^[A-Za-z][A-Za-z0-9+.-]*:') { continue }
        $target = ($href -split '#', 2)[0]
        $target = ($target -split '\?', 2)[0]
        if ($target) { $links += $target }
    }
    return @($links)
}

function Test-RepoTarget([string]$Target) {
    if ([string]::IsNullOrWhiteSpace($Target)) { return $false }
    return Test-Path -LiteralPath (Join-Path $root (Repo-Path $Target))
}

$required = @(
    'AGENTS.md',
    'README.md',
    'ARCHITECTURE.md',
    'docs/index.md',
    'docs/roadmap.md',
    'docs/development/engineering-rules.md',
    'docs/architecture/decisions/README.md',
    'docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md',
    'docs/engineering/rebaseline/D1-DOMAINS-BOUNDARIES.md',
    'docs/engineering/rebaseline/D2-IDENTITY-TENANT-DATA-OWNERSHIP.md',
    'docs/engineering/rebaseline/D3-COMMUNICATION-EVENTS.md',
    'docs/engineering/rebaseline/D4-EXTERNAL-INTEGRATIONS.md',
    'docs/engineering/rebaseline/D4-R1-PUBLICATION-INPUT.md',
    'docs/engineering/rebaseline/D5-API.md',
    'docs/engineering/rebaseline/D5-B2-OPENAPI-WIRE-AUTHORITY-TOOLING.md',
    'docs/engineering/rebaseline/EVIDENCE-REGISTER.md',
    'scripts/check-pr-title.ps1',
    'scripts/gate.ps1',
    'scripts/verify-product-oad.mjs',
    'package.json',
    'package-lock.json',
    '.github/workflows/ci.yml',
    '.github/workflows/pr-title.yml'
)

$missing = @($required | Where-Object { -not (Test-Path -LiteralPath (Join-Path $root $_) -PathType Leaf) })
if ($missing.Count -gt 0) {
    Fail ('Required repository authority/mechanics missing: ' + ($missing -join ', '))
}

$headSha = (& git -C $root rev-parse HEAD).Trim()
if ($LASTEXITCODE -ne 0 -or $headSha -notmatch '^[0-9a-f]{40}$') { Fail 'Unable to resolve HEAD.' }
$headRef = $env:GATE_HEAD_REF
if ([string]::IsNullOrWhiteSpace($headRef)) { $headRef = (& git -C $root branch --show-current).Trim() }
$isReview = $headRef -like 'review/*'

[string[]]$tracked = @(& git -C $root ls-files --cached)
if ($LASTEXITCODE -ne 0) { Fail 'Unable to enumerate tracked files.' }
$tracked = @($tracked | ForEach-Object { Repo-Path $_ })

$outsideEnvelope = @($tracked | Where-Object { -not (Test-AllowedTrackedPath $_) })
if ($outsideEnvelope.Count -gt 0) {
    Fail ("Tracked paths outside the architecture-only repository allowlist:`n" + ($outsideEnvelope -join "`n"))
}

$retiredRouter = 'docs/' + 'README.md'
$retiredMethod = 'docs/engineering/standards/' + 'root-cause-global-maximum-method.md'
$forbiddenTracked = @($tracked | Where-Object { $_ -in @($retiredRouter, $retiredMethod) })
if ($forbiddenTracked.Count -gt 0) {
    Fail ('Retired authority path remains tracked: ' + ($forbiddenTracked -join ', '))
}

$forbiddenPathFindings = @()
foreach ($path in $tracked) {
    $lower = $path.ToLowerInvariant()
    if ($lower.StartsWith('docs/work/')) {
        if ($isReview -and $lower -eq 'docs/work/current/ai-dialog.md') { continue }
        $forbiddenPathFindings += $path
        continue
    }
    if ($lower.StartsWith('docs/superpowers/')) { $forbiddenPathFindings += $path; continue }
    if ($lower -match '(^|/)(old|archive)(/|$)') { $forbiddenPathFindings += $path; continue }
    if ($lower -match '(^|/)[^/]*(handoff|dialogue|ai-dialog)[^/]*\.md$') { $forbiddenPathFindings += $path; continue }
}
if ($forbiddenPathFindings.Count -gt 0) {
    Fail ("Temporary/archive documentation contaminates the candidate:`n" + ($forbiddenPathFindings -join "`n"))
}

$agentPath = Join-Path $root 'AGENTS.md'
$indexPath = Join-Path $root 'docs/index.md'
$roadmapPath = Join-Path $root 'docs/roadmap.md'
$readmePath = Join-Path $root 'README.md'
$rulesPath = Join-Path $root 'docs/development/engineering-rules.md'

$bootstrapBytes = (Get-Item $agentPath).Length + (Get-Item $indexPath).Length + (Get-Item $roadmapPath).Length
if (-not (Test-BootstrapBudget $bootstrapBytes)) {
    Fail "Bootstrap exceeds 20 KiB: $bootstrapBytes bytes"
}

$agent = Get-Content $agentPath -Raw
$index = Get-Content $indexPath -Raw
$roadmap = Get-Content $roadmapPath -Raw
$readme = Get-Content $readmePath -Raw
$rules = Get-Content $rulesPath -Raw

if (-not $agent.Contains('AGENTS.md` → `docs/index.md` → `docs/roadmap.md`')) {
    Fail 'AGENTS.md does not declare the canonical fresh-actor route.'
}

$authorityMarker = '<!-- program-status-authority -->'
$markerOwners = @()
foreach ($path in @($tracked | Where-Object { $_.EndsWith('.md', [StringComparison]::OrdinalIgnoreCase) })) {
    if ((Get-Content (Join-Path $root $path) -Raw).Contains($authorityMarker)) { $markerOwners += $path }
}
if ($markerOwners.Count -ne 1 -or $markerOwners[0] -ne 'docs/roadmap.md') {
    Fail ('Mutable program-status authority marker must exist only in docs/roadmap.md; found: ' + ($markerOwners -join ', '))
}

$roadmapMarkers = @(
    'D0 — Product / System Definition | ACCEPTED / CLOSED',
    'D1 — Domains / Boundaries | ACCEPTED / CLOSED',
    'D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED',
    'D3 — Communication / Events | ACCEPTED / CLOSED',
    'D4 — External Integrations | ACCEPTED / CLOSED',
    'D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL',
    'D5 — API — ACCEPTED / CLOSED',
    'D5 — API | ACCEPTED / CLOSED',
    'D6 — Frontend — NEXT / NOT STARTED',
    'D6 — Frontend | NEXT / NOT STARTED',
    'D7 — Runtime / Jobs / Transactions | BLOCKED',
    'D8 — Golden Flows | BLOCKED',
    'D9 — Adversarial Architecture Review | BLOCKED',
    'Canonical Product OAD',
    'contracts/api/product/openapi.yaml',
    '95 Product operations',
    '29 ordinary Permissions',
    'Principal kinds H / A / S only',
    'https://conexus.fun',
    'Active runtime baseline',
    'NONE',
    'BLOCKED UNTIL D9'
)
foreach ($marker in $roadmapMarkers) {
    if (-not $roadmap.Contains($marker)) { Fail "docs/roadmap.md is missing required current truth: $marker" }
}

foreach ($forbidden in @('D5 — API — OPEN / ACTIVE', 'D5 — API — ACCEPTED / CLOSED', 'D6 — Frontend — NEXT / NOT STARTED', 'Exact next action', '<!-- program-status-authority -->')) {
    if ($index.Contains($forbidden)) { Fail "docs/index.md duplicates mutable roadmap state: $forbidden" }
    if ($readme.Contains($forbidden)) { Fail "README.md is not landing-only: $forbidden" }
    if ($agent.Contains($forbidden)) { Fail "AGENTS.md duplicates mutable roadmap state: $forbidden" }
}
if (-not $readme.Contains('[`AGENTS.md`](AGENTS.md)') -or -not $readme.Contains('[`docs/index.md`](docs/index.md)')) {
    Fail 'README.md must remain a landing page pointing to AGENTS.md and docs/index.md.'
}

$activeRoutingFiles = @(
    'AGENTS.md',
    'README.md',
    'docs/index.md',
    'docs/roadmap.md',
    'docs/development/engineering-rules.md',
    'docs/architecture/decisions/README.md'
)
foreach ($path in $activeRoutingFiles) {
    $text = Get-Content (Join-Path $root $path) -Raw
    if ($text.Contains($retiredRouter)) { Fail "$path still routes through retired $retiredRouter" }
    if ($text.Contains($retiredMethod)) { Fail "$path still routes through retired local Method copy" }
}

$indexRelativeLinks = 0
foreach ($target in Get-RelativeMarkdownLinks $indexPath) {
    $indexRelativeLinks++
    $resolved = [IO.Path]::GetFullPath((Join-Path (Split-Path $indexPath -Parent) $target))
    if (-not (Test-Path -LiteralPath $resolved)) { Fail "docs/index.md contains a dead relative link: $target" }
}

$durableWorkDependencies = @()
foreach ($path in @($tracked | Where-Object { $_.EndsWith('.md', [StringComparison]::OrdinalIgnoreCase) -and -not $_.StartsWith('docs/work/') })) {
    $full = Join-Path $root $path
    foreach ($target in Get-RelativeMarkdownLinks $full) {
        $resolved = [IO.Path]::GetFullPath((Join-Path (Split-Path $full -Parent) $target))
        $relative = Repo-Path ([IO.Path]::GetRelativePath($root, $resolved))
        if ($relative.StartsWith('docs/work/', [StringComparison]::Ordinal)) {
            $durableWorkDependencies += "$path -> $target"
        }
    }
}
if ($durableWorkDependencies.Count -gt 0) {
    Fail ("Durable docs depend on temporary docs/work material:`n" + ($durableWorkDependencies -join "`n"))
}

$trackedMarkdown = @{}
foreach ($path in @($tracked | Where-Object { $_.EndsWith('.md', [StringComparison]::OrdinalIgnoreCase) })) { $trackedMarkdown[$path] = $true }
$reachable = @{}
$queue = [Collections.Generic.Queue[string]]::new()
$queue.Enqueue('docs/index.md')
while ($queue.Count -gt 0) {
    $source = $queue.Dequeue()
    if ($reachable.ContainsKey($source)) { continue }
    $reachable[$source] = $true
    $sourceFull = Join-Path $root $source
    foreach ($target in Get-RelativeMarkdownLinks $sourceFull) {
        $resolved = [IO.Path]::GetFullPath((Join-Path (Split-Path $sourceFull -Parent) $target))
        $relative = Repo-Path ([IO.Path]::GetRelativePath($root, $resolved))
        if ($trackedMarkdown.ContainsKey($relative) -and -not $reachable.ContainsKey($relative)) { $queue.Enqueue($relative) }
    }
}
$unreachableDocs = @()
foreach ($path in @($trackedMarkdown.Keys | Where-Object { $_.StartsWith('docs/') -and -not $_.StartsWith('docs/work/') })) {
    if (-not $reachable.ContainsKey($path)) { $unreachableDocs += $path }
}
if ($unreachableDocs.Count -gt 0) {
    Fail ("Durable docs are not reachable from docs/index.md or a routed child document:`n" + (($unreachableDocs | Sort-Object) -join "`n"))
}

$machineRouteFiles = @($tracked | Where-Object { $_ -match '(?i)(routes|selectors)\.(json|ya?ml)$' })
$machineSelectors = 0
foreach ($path in $machineRouteFiles) {
    if (-not $path.EndsWith('.json', [StringComparison]::OrdinalIgnoreCase)) { continue }
    $document = Get-Content (Join-Path $root $path) -Raw | ConvertFrom-Json
    $json = $document | ConvertTo-Json -Depth 50
    foreach ($match in [regex]::Matches($json, '"([^"\r\n]+\.(?:md|ps1|json|ya?ml|toml))"')) {
        $candidate = Repo-Path $match.Groups[1].Value
        if ($candidate -match '^[A-Za-z][A-Za-z0-9+.-]*:') { continue }
        $machineSelectors++
        if (-not (Test-RepoTarget $candidate)) { Fail "$path contains unresolved machine selector: $candidate" }
    }
}

if (-not $rules.Contains('subject population = 0')) { Fail 'Engineering rules lost the attributable zero-population retirement condition.' }
if (-not $rules.Contains('full replacement coverage = proved')) { Fail 'Engineering rules lost the complete replacement-coverage retirement condition.' }
if (-not $rules.Contains('Dependency or lockfile change')) { Fail 'Engineering rules lost explicit dependency/lockfile scope.' }
if (-not $rules.Contains('Never bend accepted target architecture')) { Fail 'Engineering rules permit legacy code/tests to warp target architecture.' }

$baseSha = $env:GATE_BASE_SHA
if ([string]::IsNullOrWhiteSpace($baseSha) -or $baseSha -notmatch '^[0-9a-f]{40}$') {
    $originMain = (& git -C $root rev-parse --verify refs/remotes/origin/main 2>$null).Trim()
    if ($LASTEXITCODE -eq 0 -and $originMain -match '^[0-9a-f]{40}$') {
        $baseSha = (& git -C $root merge-base $headSha $originMain).Trim()
    } else {
        $baseSha = (& git -C $root rev-parse HEAD^ 2>$null).Trim()
    }
}
if ([string]::IsNullOrWhiteSpace($baseSha) -or $baseSha -notmatch '^[0-9a-f]{40}$') { Fail 'Unable to resolve a non-vacuous diff base.' }
& git -C $root cat-file -e "$baseSha^{commit}" 2>$null
if ($LASTEXITCODE -ne 0) { Fail "Base/candidate commit does not exist locally: $baseSha" }

if ($isReview) {
    $candidateRef = $env:GATE_CANDIDATE_REF
    if ([string]::IsNullOrWhiteSpace($candidateRef)) { $candidateRef = $env:GATE_BASE_REF }
    if ([string]::IsNullOrWhiteSpace($candidateRef)) { Fail 'Review mode requires an explicit candidate branch/ref.' }

    & git -C $root show-ref --verify --quiet "refs/remotes/origin/$candidateRef"
    if ($LASTEXITCODE -ne 0) {
        & git -C $root ls-remote --exit-code --heads origin $candidateRef *> $null
        if ($LASTEXITCODE -ne 0) { Fail "Review candidate ref does not exist on origin: $candidateRef" }
    }

    $candidateTree = @(& git -C $root ls-tree -r --name-only $baseSha)
    $candidateWork = @($candidateTree | Where-Object { (Repo-Path $_).StartsWith('docs/work/') })
    if ($candidateWork.Count -gt 0) { Fail 'Exact candidate tree is contaminated by docs/work material.' }

    $reviewNames = @(& git -C $root diff --name-only "$baseSha..$headSha")
    if (-not (Test-ReviewDiffNames $reviewNames)) {
        Fail ('Review branch must differ from exact candidate only by docs/work/current/ai-dialog.md; found: ' + ($reviewNames -join ', '))
    }
    $diffRange = "$baseSha..$headSha"
} else {
    $currentWork = @($tracked | Where-Object { $_.StartsWith('docs/work/') })
    if ($currentWork.Count -gt 0) { Fail 'Candidate/main contains docs/work material.' }
    $diffRange = "$baseSha...$headSha"
}

$changedFiles = @(& git -C $root diff --name-only $diffRange)
if ($LASTEXITCODE -ne 0) { Fail "Unable to enumerate diff range $diffRange" }
if ($changedFiles.Count -eq 0) { Fail "Diff check is vacuous: $diffRange contains zero changed files." }
$diffCheck = @(& git -C $root diff --check $diffRange 2>&1)
if ($LASTEXITCODE -ne 0 -or $diffCheck.Count -gt 0) {
    Fail ("git diff --check failed for $diffRange`n" + ($diffCheck -join "`n"))
}

$package = Get-Content (Join-Path $root 'package.json') -Raw | ConvertFrom-Json
foreach ($property in @('workspaces', 'dependencies', 'devDependencies')) {
    if ($package.PSObject.Properties.Name -contains $property) { Fail "package.json contains retired property: $property" }
}

if ($Lane -eq 'full') {
    & node (Join-Path $root 'scripts/verify-product-oad.mjs')
    if ($LASTEXITCODE -ne 0) { Fail 'Canonical Product OAD proof failed.' }
}

$negativeControls = 0
if (-not (Test-BootstrapBudget 20481) -and (Test-BootstrapBudget 20480)) { $negativeControls++ } else { Fail 'Bootstrap-budget negative control failed.' }
if (-not (Test-AllowedTrackedPath 'apps/server_core/main.go') -and (Test-AllowedTrackedPath 'docs/index.md')) { $negativeControls++ } else { Fail 'Repository-allowlist negative control failed.' }
if ((Test-ReviewDiffNames @('docs/work/current/ai-dialog.md')) -and -not (Test-ReviewDiffNames @('docs/work/current/ai-dialog.md','README.md'))) { $negativeControls++ } else { Fail 'Review-isolation negative control failed.' }
if (-not (Test-RepoTarget 'docs/__missing-route__.md') -and (Test-RepoTarget 'docs/index.md')) { $negativeControls++ } else { Fail 'Route-resolution negative control failed.' }
if ($authorityMarker -eq '<!-- program-status-authority -->' -and $authorityMarker -ne '<!-- not-the-authority -->') { $negativeControls++ } else { Fail 'Status-authority marker negative control failed.' }

Write-Host "gate lane: $Lane"
Write-Host "required files: $($required.Count)"
Write-Host "tracked files inspected: $($tracked.Count)"
Write-Host "bootstrap_bytes: $bootstrapBytes / 20480"
Write-Host "docs_index_relative_links: $indexRelativeLinks"
Write-Host "durable_docs_reachable: $($reachable.Count)"
Write-Host "machine_route_files: $($machineRouteFiles.Count) selectors: $machineSelectors"
Write-Host "diff_range: $diffRange changed_files: $($changedFiles.Count)"
Write-Host "review_mode: $isReview"
Write-Host 'legacy_runtime_population: 0'
Write-Host "negative_controls: $negativeControls/5"
Write-Host 'gate: PASS'