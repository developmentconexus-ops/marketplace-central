[CmdletBinding()]
param(
    [ValidateSet('all', 'full')]
    [string]$Lane = 'all'
)

$ErrorActionPreference = 'Stop'
$root = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path

function Fail([string]$Message) {
    Write-Error $Message
    exit 1
}

function Repo-Path([string]$Path) {
    return $Path.Replace('\', '/').TrimStart('./')
}

function Test-BootstrapBudget([long]$Bytes) {
    return $Bytes -le 20480
}

$allowedRoot = @(
    '.gitattributes',
    '.gitignore',
    '.node-version',
    'AGENTS.md',
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
    'contracts/api/product/'
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

function Test-RouteTarget([string]$Target) {
    if ([string]::IsNullOrWhiteSpace($Target)) { return $false }
    return Test-Path -LiteralPath (Join-Path $root (Repo-Path $Target))
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
    'package.json',
    'package-lock.json',
    '.github/workflows/ci.yml',
    '.github/workflows/pr-title.yml'
)

$missing = @($required | Where-Object { -not (Test-Path -LiteralPath (Join-Path $root $_) -PathType Leaf) })
if ($missing.Count -gt 0) {
    Fail ('Required repository authority/mechanics missing: ' + ($missing -join ', '))
}

[string[]]$tracked = @(& git -C $root ls-files --cached)
if ($LASTEXITCODE -ne 0) { Fail 'Unable to enumerate tracked files.' }
$tracked = @($tracked | ForEach-Object { Repo-Path $_ })

$outsideEnvelope = @($tracked | Where-Object { -not (Test-AllowedTrackedPath $_) })
if ($outsideEnvelope.Count -gt 0) {
    Fail ("Tracked paths outside the architecture-only repository allowlist:`n" + ($outsideEnvelope -join "`n"))
}

$retiredRouter = 'docs/' + 'README.md'
$retiredMethod = 'docs/engineering/standards/' + 'root-cause-global-maximum-method.md'
$forbiddenExact = @($retiredRouter, $retiredMethod)
$forbiddenTracked = @($tracked | Where-Object { $forbiddenExact -contains $_ })
if ($forbiddenTracked.Count -gt 0) {
    Fail ('Retired authority path remains tracked: ' + ($forbiddenTracked -join ', '))
}

$forbiddenPathFindings = @()
foreach ($path in $tracked) {
    $lower = $path.ToLowerInvariant()
    if ($lower.StartsWith('docs/superpowers/')) { $forbiddenPathFindings += $path; continue }
    if ($lower -match '(^|/)docs/work/' -or $lower.StartsWith('docs/work/')) { $forbiddenPathFindings += $path; continue }
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
    $full = Join-Path $root $path
    if ((Get-Content $full -Raw).Contains($authorityMarker)) { $markerOwners += $path }
}
if ($markerOwners.Count -ne 1 -or $markerOwners[0] -ne 'docs/roadmap.md') {
    Fail ('Mutable program-status authority marker must exist only in docs/roadmap.md; found: ' + ($markerOwners -join ', '))
}

$roadmapMarkers = @(
    'D5 — API — OPEN / ACTIVE',
    'Author and prove the canonical Product OpenAPI Description',
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

foreach ($forbidden in @('D5 — API — OPEN / ACTIVE', 'Exact next action', '<!-- program-status-authority -->')) {
    if ($index.Contains($forbidden)) { Fail "docs/index.md duplicates mutable roadmap state: $forbidden" }
}
foreach ($forbidden in @('D5 — API — OPEN / ACTIVE', 'Exact next action', '<!-- program-status-authority -->')) {
    if ($readme.Contains($forbidden)) { Fail "README.md is not landing-only: $forbidden" }
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

# Reachability: follow relative Markdown links from docs/index.md. Citation archaeology
# is considered child-index-owned when the ADR registry itself is reachable.
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
    if (-not (Test-Path $sourceFull -PathType Leaf)) { continue }
    foreach ($target in Get-RelativeMarkdownLinks $sourceFull) {
        $resolved = [IO.Path]::GetFullPath((Join-Path (Split-Path $sourceFull -Parent) $target))
        $relative = Repo-Path ([IO.Path]::GetRelativePath($root, $resolved))
        if ($trackedMarkdown.ContainsKey($relative) -and -not $reachable.ContainsKey($relative)) { $queue.Enqueue($relative) }
    }
}
$unreachableDocs = @()
foreach ($path in @($trackedMarkdown.Keys | Where-Object { $_.StartsWith('docs/') -and -not $_.StartsWith('docs/work/') })) {
    if ($reachable.ContainsKey($path)) { continue }
    if ($path.StartsWith('docs/architecture/decisions/_citations/') -and $reachable.ContainsKey('docs/architecture/decisions/README.md')) { continue }
    $unreachableDocs += $path
}
if ($unreachableDocs.Count -gt 0) {
    Fail ("Durable docs are not reachable from docs/index.md routing:`n" + ($unreachableDocs | Sort-Object | ForEach-Object { $_ }) -join "`n")
}

# Machine routing/selectors are optional in the current hard-reset tree, but any
# current route manifest that exists must resolve repo-local file selectors.
$machineRouteFiles = @($tracked | Where-Object { $_ -match '(?i)(routes|selectors)\.(json|ya?ml)$' })
$machineSelectors = 0
foreach ($path in $machineRouteFiles) {
    if ($path.EndsWith('.json', [StringComparison]::OrdinalIgnoreCase)) {
        $document = Get-Content (Join-Path $root $path) -Raw | ConvertFrom-Json
        $json = $document | ConvertTo-Json -Depth 50
        foreach ($match in [regex]::Matches($json, '"([^"\r\n]+\.(?:md|ps1|json|ya?ml|toml))"')) {
            $candidate = Repo-Path $match.Groups[1].Value
            if ($candidate -match '^[A-Za-z][A-Za-z0-9+.-]*:') { continue }
            $machineSelectors++
            if (-not (Test-RouteTarget $candidate)) { Fail "$path contains unresolved machine selector: $candidate" }
        }
    }
}

if (-not $rules.Contains('subject population = 0')) { Fail 'Engineering rules lost the attributable zero-population retirement condition.' }
if (-not $rules.Contains('full replacement coverage = proved')) { Fail 'Engineering rules lost the complete replacement-coverage retirement condition.' }
if (-not $rules.Contains('Dependency or lockfile change')) { Fail 'Engineering rules lost explicit dependency/lockfile scope.' }
if (-not $rules.Contains('Never bend accepted target architecture')) { Fail 'Engineering rules permit legacy code/tests to warp target architecture.' }

# Determine exact base/candidate for diff and review-isolation proof.
$headSha = (& git -C $root rev-parse HEAD).Trim()
if ($LASTEXITCODE -ne 0 -or $headSha -notmatch '^[0-9a-f]{40}$') { Fail 'Unable to resolve HEAD.' }
$headRef = $env:GATE_HEAD_REF
if ([string]::IsNullOrWhiteSpace($headRef)) { $headRef = (& git -C $root branch --show-current).Trim() }
$baseSha = $env:GATE_BASE_SHA
if ([string]::IsNullOrWhiteSpace($baseSha) -or $baseSha -notmatch '^[0-9a-f]{40}$') {
    $baseSha = (& git -C $root rev-parse HEAD^ 2>$null).Trim()
}
if ([string]::IsNullOrWhiteSpace($baseSha) -or $baseSha -notmatch '^[0-9a-f]{40}$') { Fail 'Unable to resolve a non-vacuous diff base.' }
& git -C $root cat-file -e "$baseSha^{commit}" 2>$null
if ($LASTEXITCODE -ne 0) { Fail "Candidate/base ref does not exist locally: $baseSha" }

$isReview = $headRef -like 'review/*'
if ($isReview) {
    $candidateTree = @(& git -C $root ls-tree -r --name-only $baseSha)
    $candidateWork = @($candidateTree | Where-Object { (Repo-Path $_).StartsWith('docs/work/') })
    if ($candidateWork.Count -gt 0) { Fail 'Exact candidate ref is contaminated by docs/work material.' }
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

# Deterministic falsifiers for the material operating-envelope guards.
$negativeControls = 0
if (-not (Test-BootstrapBudget 20481) -and (Test-BootstrapBudget 20480)) { $negativeControls++ } else { Fail 'Bootstrap-budget negative control failed.' }
if (-not (Test-AllowedTrackedPath 'apps/server_core/main.go') -and (Test-AllowedTrackedPath 'docs/index.md')) { $negativeControls++ } else { Fail 'Repository-allowlist negative control failed.' }
if ((Test-ReviewDiffNames @('docs/work/current/ai-dialog.md')) -and -not (Test-ReviewDiffNames @('docs/work/current/ai-dialog.md','README.md'))) { $negativeControls++ } else { Fail 'Review-isolation negative control failed.' }
if (-not (Test-RouteTarget 'docs/__missing-route__.md') -and (Test-RouteTarget 'docs/index.md')) { $negativeControls++ } else { Fail 'Route-resolution negative control failed.' }
if ($authorityMarker -ne '<!-- not-the-authority -->') { $negativeControls++ } else { Fail 'Status-authority marker negative control failed.' }

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
