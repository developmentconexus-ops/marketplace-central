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

$required = @(
    'AGENTS.md',
    'ARCHITECTURE.md',
    'README.md',
    'docs/README.md',
    'docs/architecture/decisions/README.md',
    'docs/engineering/standards/root-cause-global-maximum-method.md',
    'docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md',
    'docs/engineering/rebaseline/D1-DOMAINS-BOUNDARIES.md',
    'docs/engineering/rebaseline/D2-IDENTITY-TENANT-DATA-OWNERSHIP.md',
    'docs/engineering/rebaseline/D3-COMMUNICATION-EVENTS.md',
    'docs/engineering/rebaseline/D4-EXTERNAL-INTEGRATIONS.md',
    'docs/engineering/rebaseline/D4-R1-PUBLICATION-INPUT.md',
    'docs/engineering/rebaseline/D5-API.md',
    'docs/engineering/rebaseline/D5-B2-OPENAPI-WIRE-AUTHORITY-TOOLING.md',
    'scripts/check-pr-title.ps1',
    'scripts/gate.ps1',
    'package.json',
    'package-lock.json',
    '.github/workflows/ci.yml',
    '.github/workflows/pr-title.yml'
)

$missing = @($required | Where-Object { -not (Test-Path -LiteralPath (Join-Path $root $_) -PathType Leaf) })
if ($missing.Count -gt 0) {
    Fail ('Required authority/mechanics missing: ' + ($missing -join ', '))
}

$forbiddenPrefixes = @(
    'apps/',
    'packages/',
    'deploy/',
    'docker/',
    'design-system/',
    'evidence/',
    'contracts/events/',
    'contracts/gate/',
    'contracts/governance/',
    'docs/superpowers/',
    '.github/ISSUE_TEMPLATE/',
    '.mnfs/',
    'scripts/harness/'
)

$forbiddenExact = @(
    '.dockerignore',
    '.prettierignore',
    '.prettierrc.json',
    '.github/pull_request_template.md',
    'Makefile',
    'contracts/api/marketplace-central.openapi.yaml',
    'docker-compose.yml',
    'eslint.config.mjs',
    'go.work',
    'go.work.sum',
    'run-server.ps1',
    'tsconfig.base.json',
    'tsconfig.eslint.json',
    'vitest.config.ts',
    'docs/operations/live-oracle-docker.md',
    'docs/engineering/rebaseline/README.md',
    'scripts/dev-local.ps1',
    'scripts/harness.ps1',
    'scripts/release.sh',
    'scripts/run-live-oracle-docker.ps1',
    'scripts/test-verify-baseline.ps1',
    'scripts/verify-baseline.ps1',
    'skills-lock.json'
)

function Test-LegacyPath([string]$Path) {
    $normalized = $Path.Replace('\', '/')
    if ($forbiddenExact -contains $normalized) { return $true }
    if ($normalized -like '*/AI-DIALOG.md' -or $normalized -eq 'AI-DIALOG.md') { return $true }
    foreach ($prefix in $forbiddenPrefixes) {
        if ($normalized.StartsWith($prefix, [StringComparison]::Ordinal)) { return $true }
    }
    return $false
}

[string[]]$tracked = @(& git -C $root ls-files --cached)
if ($LASTEXITCODE -ne 0) { Fail 'Unable to enumerate the Git index.' }

$legacy = @($tracked | Where-Object { Test-LegacyPath $_ })
if ($legacy.Count -gt 0) {
    Fail ("Retired legacy paths are tracked:`n" + ($legacy -join "`n"))
}

$negativeFixture = @('apps/server_core/main.go')
$negativeHits = @($negativeFixture | Where-Object { Test-LegacyPath $_ })
if ($negativeHits.Count -ne 1) {
    Fail 'Zero-legacy guard negative control did not detect its synthetic fixture.'
}

$indexPath = Join-Path $root 'docs/README.md'
$agentPath = Join-Path $root 'AGENTS.md'
$index = Get-Content -LiteralPath $indexPath -Raw
$agent = Get-Content -LiteralPath $agentPath -Raw

$bootstrapBytes = (Get-Item -LiteralPath $agentPath).Length + (Get-Item -LiteralPath $indexPath).Length
if ($bootstrapBytes -gt 20480) {
    Fail "Two-file bootstrap exceeds 20 KiB: $bootstrapBytes bytes"
}

$indexDir = Split-Path -Parent $indexPath
$linkMatches = [regex]::Matches($index, '\[[^\]]+\]\(([^)]+)\)')
$relativeLinks = 0
foreach ($match in $linkMatches) {
    $href = $match.Groups[1].Value.Trim()
    if (-not $href -or $href.StartsWith('#') -or $href -match '^[A-Za-z][A-Za-z0-9+.-]*:') { continue }
    $target = ($href -split '#', 2)[0]
    $target = ($target -split '\?', 2)[0]
    if (-not $target) { continue }
    $relativeLinks++
    $resolved = [System.IO.Path]::GetFullPath((Join-Path $indexDir $target))
    if (-not (Test-Path -LiteralPath $resolved)) {
        Fail "docs/README.md contains a dead relative link: $href"
    }
}

$markers = @(
    'D5 — API — OPEN / ACTIVE',
    'Author and prove the canonical Product OpenAPI Description',
    'contracts/api/product/openapi.yaml',
    'BLOCKED until D9 is accepted',
    'Active runtime baseline'
)
foreach ($marker in $markers) {
    if (-not $index.Contains($marker)) { Fail "docs/README.md is missing authority marker: $marker" }
}
if (-not $agent.Contains('DEFAULT READ: `AGENTS.md` → `docs/README.md`. STOP.')) {
    Fail 'AGENTS.md no longer declares the two-file default bootstrap.'
}

$package = Get-Content -LiteralPath (Join-Path $root 'package.json') -Raw | ConvertFrom-Json
foreach ($property in @('workspaces', 'dependencies', 'devDependencies')) {
    if ($package.PSObject.Properties.Name -contains $property) {
        Fail "package.json contains retired property: $property"
    }
}

Write-Host "gate lane: $Lane"
Write-Host "required files: $($required.Count)"
Write-Host "tracked files inspected: $($tracked.Count)"
Write-Host "bootstrap_bytes: $bootstrapBytes"
Write-Host "docs_index_relative_links: $relativeLinks"
Write-Host 'legacy tracked paths: 0'
Write-Host 'negative controls: 1/1'
Write-Host 'gate: PASS'