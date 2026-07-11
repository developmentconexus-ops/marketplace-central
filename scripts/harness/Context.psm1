Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$script:RepositoryRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)

function New-ContextResult {
  param([bool]$Passed, [string]$ErrorCode = '', [string]$Id = '', [string]$Path = '', [object]$Pack = $null)
  [pscustomobject]@{
    PSTypeName = 'MarketplaceCentral.Harness.ContextResult'
    Passed = $Passed
    Status = if ($Passed) { 'passed' } else { 'failed' }
    ErrorCode = if ($Passed) { $null } else { $ErrorCode }
    Id = $Id
    Path = $Path
    Pack = $Pack
  }
}

function ConvertTo-ContextPath {
  param([string]$Path, [switch]$AllowScope)
  if ([string]::IsNullOrWhiteSpace($Path) -or [IO.Path]::IsPathRooted($Path) -or $Path -match '(^|[\\/])\.\.([\\/]|$)') { return $null }
  $normalized = $Path.Trim().Replace('\', '/').TrimEnd('/')
  if ($normalized.StartsWith('./', [StringComparison]::Ordinal)) { $normalized = $normalized.Substring(2) }
  if ($AllowScope) { $normalized = $normalized -replace '/\*\*$', '' }
  if ([string]::IsNullOrWhiteSpace($normalized) -or $normalized -match '[:*?\[]') { return $null }
  return $normalized
}

function Resolve-ContextRepositoryPath {
  param([string]$RepositoryRoot, [string]$Path)
  Join-Path $RepositoryRoot ($Path.Replace('/', [IO.Path]::DirectorySeparatorChar))
}

function Get-SafeContextPath {
  param([string]$RepositoryRoot, [string]$Path)
  if ([string]::IsNullOrWhiteSpace($Path)) { return '' }
  $root = [IO.Path]::GetFullPath($RepositoryRoot).TrimEnd([IO.Path]::DirectorySeparatorChar, [IO.Path]::AltDirectorySeparatorChar)
  $full = [IO.Path]::GetFullPath($Path)
  if ($full -eq $root -or $full.StartsWith($root + [IO.Path]::DirectorySeparatorChar, [StringComparison]::OrdinalIgnoreCase)) {
    return [IO.Path]::GetRelativePath($root, $full).Replace('\', '/')
  }
  return [IO.Path]::GetFileName($full)
}

function Get-ContextFileHash {
  param([string]$Path)
  ([Convert]::ToHexString([Security.Cryptography.SHA256]::HashData([IO.File]::ReadAllBytes($Path)))).ToLowerInvariant()
}

function Get-MarkdownSection {
  param([string]$Content, [string]$Heading)
  $match = [regex]::Match($Content, "(?ms)^## $([regex]::Escape($Heading))\s*\r?\n(?<body>.*?)(?=^## |\z)")
  if (-not $match.Success) { return '' }
  $match.Groups['body'].Value.Trim()
}

function ConvertTo-CompactMarkdownText {
  param([string]$Text)
  (($Text -replace '(?m)^\s*[-*]\s+', '' -replace '`', '' -replace '\s+', ' ').Trim())
}

function Get-ObservableDone {
  param([string]$FeatureContent)
  $section = Get-MarkdownSection $FeatureContent 'Expected Output'
  if ([string]::IsNullOrWhiteSpace($section)) { $section = Get-MarkdownSection $FeatureContent 'Validation Expectations' }
  $items = @([regex]::Matches($section, '(?m)^-\s+(?<item>.+(?:\r?\n\s{2,}.+)*)') | ForEach-Object { ConvertTo-CompactMarkdownText $_.Groups['item'].Value })
  if ($items.Count -eq 0 -and -not [string]::IsNullOrWhiteSpace($section)) { $items = @(ConvertTo-CompactMarkdownText $section) }
  $items
}

function Test-PathIntersection {
  param([string]$Left, [string]$Right)
  $leftPath = $Left.TrimEnd('/'); $rightPath = $Right.TrimEnd('/')
  $leftPath -eq $rightPath -or $leftPath.StartsWith($rightPath + '/', [StringComparison]::Ordinal) -or $rightPath.StartsWith($leftPath + '/', [StringComparison]::Ordinal)
}

function Test-RequestedPathAllowed {
  param([string]$Requested, [string[]]$Declared)
  @($Declared | Where-Object { $Requested -eq $_ -or $Requested.StartsWith($_ + '/', [StringComparison]::Ordinal) }).Count -gt 0
}

function Get-MissionRootPath {
  param([string]$FeaturePath)
  $parts = @($FeaturePath.Split('/'))
  $missionIndex = -1
  for ($index = 0; $index -lt $parts.Count; $index++) {
    if ($parts[$index] -match '^MIS-[0-9]+') { $missionIndex = $index; break }
  }
  if ($missionIndex -lt 0) { return $null }
  ($parts[0..$missionIndex] -join '/')
}

function Import-FeatureWorkContract {
  param([string]$PlanContent, [string]$RepositoryRoot)
  $headings = @([regex]::Matches($PlanContent, '(?m)^## Machine Work Contract[ \t]*$'))
  $section = Get-MarkdownSection $PlanContent 'Machine Work Contract'
  $block = [regex]::Match($section, '(?ms)^```json[ \t]*\r?\n(?<json>.*?)\r?\n```[ \t]*$')
  if ($headings.Count -ne 1 -or -not $block.Success) { return New-ContextResult $false 'CTX_FEATURE_INVALID' 'machine-work-contract-count' }
  $temporaryPath = Join-Path ([IO.Path]::GetTempPath()) "feature-work-contract-$([guid]::NewGuid().ToString('N')).json"
  try {
    Set-Content -LiteralPath $temporaryPath -Value $block.Groups['json'].Value -Encoding utf8
    $schemaPath = Join-Path $RepositoryRoot 'contracts/governance/schemas/feature-work-contract.schema.json'
    if (-not (Test-Json -LiteralPath $temporaryPath -SchemaFile $schemaPath -ErrorAction Stop)) { throw 'schema rejected work contract' }
    $contract = Get-Content -Raw -LiteralPath $temporaryPath | ConvertFrom-Json -Depth 100
    return New-ContextResult $true -Pack $contract
  } catch {
    return New-ContextResult $false 'CTX_FEATURE_INVALID' 'machine-work-contract-schema'
  } finally {
    Remove-Item -LiteralPath $temporaryPath -Force -ErrorAction SilentlyContinue
  }
}

function Expand-ContextCommandTemplate {
  param([string]$Template, [string]$BaseSha, [string]$ContextPath, [string]$FeaturePath)
  $variables = @([regex]::Matches($Template, '\{(?<name>[a-z_]+)\}') | ForEach-Object { $_.Groups['name'].Value } | Sort-Object -Unique)
  if (@($variables | Where-Object { $_ -notin @('base_sha', 'candidate_sha', 'context_path', 'feature_path') }).Count -gt 0) { return $null }
  $expanded = $Template.Replace('{base_sha}', $BaseSha).Replace('{candidate_sha}', $BaseSha).Replace('{context_path}', $ContextPath).Replace('{feature_path}', $FeaturePath)
  if ($expanded -match '\{[^{}]+\}') { return $null }
  $expanded
}

function Get-ContextRisk {
  param([string[]]$AllowedPaths, [object[]]$SharedSeams, [object[]]$Lanes, [string[]]$SideEffects)
  $moduleScopes = @($AllowedPaths | ForEach-Object { if ($_ -match '^apps/server_core/internal/modules/(?<module>[^/]+)') { $Matches.module } } | Sort-Object -Unique)
  $live = @($Lanes | Where-Object { $_.network -eq 'live' -or $_.id -eq 'provider-write' }).Count -gt 0 -or @($SideEffects | Where-Object { $_ -in @('provider-read', 'provider-write', 'database-write') }).Count -gt 0
  $level = if ($live) { 'L3' } elseif ($SharedSeams.Count -gt 0 -or $moduleScopes.Count -gt 1) { 'L2' } elseif ($moduleScopes.Count -eq 1) { 'L1' } else { 'L0' }
  [ordered]@{
    level = $level
    review_policy = @{ L0 = 'self-review'; L1 = 'standard-review'; L2 = 'independent-review'; L3 = 'safety-review' }[$level]
    advisory_model = @{ L0 = 'fast-capable-model'; L1 = 'general-reasoning-model'; L2 = 'strong-reasoning-model'; L3 = 'strongest-safety-reviewed-model' }[$level]
  }
}

function Get-EstimatedContextTokens {
  param([object]$Pack)
  $withoutEstimate = [ordered]@{}
  foreach ($property in $Pack.PSObject.Properties) {
    if ($property.Name -ne 'estimated_input_tokens') { $withoutEstimate[$property.Name] = $property.Value }
  }
  $json = $withoutEstimate | ConvertTo-Json -Depth 100 -Compress
  [int][Math]::Ceiling([Text.Encoding]::UTF8.GetByteCount($json) / 4.0)
}

function Test-JsonEqual {
  param([object]$Left, [object]$Right)
  ($Left | ConvertTo-Json -Depth 100 -Compress) -ceq ($Right | ConvertTo-Json -Depth 100 -Compress)
}

function New-CanonicalContextPack {
  param([string]$RepositoryRoot, [string]$FeaturePath, [string]$BaseSha, [string[]]$AllowedPath, [string]$ContextPath)
  if ($BaseSha -notmatch '^[0-9a-f]{40}$') { return New-ContextResult $false 'CTX_BASE_SHA_INVALID' 'base-sha' }
  $normalizedFeature = ConvertTo-ContextPath $FeaturePath
  if ($null -eq $normalizedFeature) { return New-ContextResult $false 'CTX_FEATURE_INVALID' 'feature-path' $FeaturePath }
  $missionRoot = Get-MissionRootPath $normalizedFeature
  if ($null -eq $missionRoot) { return New-ContextResult $false 'CTX_FEATURE_INVALID' 'mission-root' $normalizedFeature }
  $featureDirectory = Resolve-ContextRepositoryPath $RepositoryRoot $normalizedFeature
  $milestoneDirectory = Split-Path -Parent $featureDirectory
  $planPath = Join-Path $featureDirectory 'plan.md'
  if (-not (Test-Path -LiteralPath $planPath -PathType Leaf)) { return New-ContextResult $false 'CTX_FEATURE_INVALID' 'feature-sources' $normalizedFeature }
  $planContent = Get-Content -Raw -LiteralPath $planPath
  $contractResult = Import-FeatureWorkContract $planContent $RepositoryRoot
  if (-not $contractResult.Passed) { return $contractResult }
  $contract = $contractResult.Pack

  $declaredAllowed = @($contract.allowed_paths | ForEach-Object { ConvertTo-ContextPath ([string]$_) -AllowScope } | Sort-Object -Unique)
  $forbidden = @($contract.forbidden_paths | ForEach-Object { ConvertTo-ContextPath ([string]$_) -AllowScope } | Sort-Object -Unique)
  if (@($declaredAllowed | Where-Object { $null -eq $_ }).Count -gt 0 -or @($forbidden | Where-Object { $null -eq $_ }).Count -gt 0) { return New-ContextResult $false 'CTX_FEATURE_INVALID' 'work-contract-path' }
  $requested = [Collections.Generic.List[string]]::new()
  foreach ($path in $AllowedPath) {
    $normalized = ConvertTo-ContextPath $path -AllowScope
    if ($null -eq $normalized -or -not (Test-RequestedPathAllowed $normalized $declaredAllowed)) { return New-ContextResult $false 'CTX_PATH_OUTSIDE_SCOPE' 'allowed-path' $path }
    if (@($forbidden | Where-Object { Test-PathIntersection $normalized $_ }).Count -gt 0) { return New-ContextResult $false 'CTX_PATH_SCOPE_CONFLICT' 'allowed-forbidden' $path }
    $requested.Add($normalized)
  }
  $requestedPaths = @($requested | Sort-Object -Unique)
  if ($requestedPaths.Count -eq 0) { return New-ContextResult $false 'CTX_PATH_OUTSIDE_SCOPE' 'allowed-path' }

  $featurePathFile = Join-Path $featureDirectory 'feature.md'
  $specPath = Join-Path $featureDirectory 'spec.md'
  $sourcePaths = @(
    $featurePathFile
    $specPath
    $planPath
    (Join-Path $milestoneDirectory 'milestone.md')
    (Join-Path $milestoneDirectory 'validation-contract.md')
    (Resolve-ContextRepositoryPath $RepositoryRoot "$missionRoot/mission.md")
    (Join-Path $RepositoryRoot 'AGENTS.md')
  )
  foreach ($extra in @($contract.required_sources)) {
    $normalizedExtra = ConvertTo-ContextPath ([string]$extra)
    if ($null -eq $normalizedExtra) { return New-ContextResult $false 'CTX_FEATURE_INVALID' 'required-source' ([string]$extra) }
    $sourcePaths += Resolve-ContextRepositoryPath $RepositoryRoot $normalizedExtra
  }
  $sourcePaths = @($sourcePaths | Select-Object -Unique)
  foreach ($sourcePath in $sourcePaths) {
    if (-not (Test-Path -LiteralPath $sourcePath -PathType Leaf)) { return New-ContextResult $false 'CTX_SOURCE_MISSING' 'source' (Get-SafeContextPath $RepositoryRoot $sourcePath) }
  }

  $featureContent = Get-Content -Raw -LiteralPath $featurePathFile
  $specContent = Get-Content -Raw -LiteralPath $specPath
  $validationContent = Get-Content -Raw -LiteralPath (Join-Path $milestoneDirectory 'validation-contract.md')
  $objective = ConvertTo-CompactMarkdownText (Get-MarkdownSection $featureContent 'Brief')
  $done = @(Get-ObservableDone $featureContent)
  if ([string]::IsNullOrWhiteSpace($objective) -or $done.Count -eq 0) { return New-ContextResult $false 'CTX_FEATURE_INVALID' 'mnfs-headings' $normalizedFeature }

  Import-Module (Join-Path $PSScriptRoot 'Policy.psm1') -Force
  $lanesDocument = Import-GovernanceDocument -Path (Join-Path $RepositoryRoot 'contracts/governance/execution-lanes.json') -SchemaPath (Join-Path $RepositoryRoot 'contracts/governance/schemas/execution-lanes.schema.json')
  $seamsDocument = Import-GovernanceDocument -Path (Join-Path $RepositoryRoot 'contracts/governance/shared-seams.json') -SchemaPath (Join-Path $RepositoryRoot 'contracts/governance/schemas/shared-seams.schema.json')
  if (-not $lanesDocument.Passed -or -not $seamsDocument.Passed) { return New-ContextResult $false 'CTX_FEATURE_INVALID' 'governance-contracts' }
  $laneById = @{}; foreach ($lane in @($lanesDocument.Document.lanes)) { $laneById[[string]$lane.id] = $lane }
  $commandIds = @($contract.commands.id)
  if (@($commandIds | Group-Object | Where-Object Count -gt 1).Count -gt 0) { return New-ContextResult $false 'CTX_PROOF_REFERENCE_INVALID' 'command-id' }
  $commands = [Collections.Generic.List[object]]::new()
  $usedLanes = [Collections.Generic.List[object]]::new()
  foreach ($command in @($contract.commands)) {
    if (-not $laneById.ContainsKey([string]$command.lane_id)) { return New-ContextResult $false 'CTX_PROOF_REFERENCE_INVALID' ([string]$command.id) }
    $expanded = Expand-ContextCommandTemplate ([string]$command.command_template) $BaseSha $ContextPath $normalizedFeature
    if ($null -eq $expanded) { return New-ContextResult $false 'CTX_FEATURE_INVALID' ([string]$command.id) }
    $lane = $laneById[[string]$command.lane_id]
    $usedLanes.Add($lane)
    $targetLabel = if ([string]$command.id -like 'cold-real-*') { 'cold-gate' } else { [string]$lane.target_label }
    $commands.Add([ordered]@{ id = [string]$command.id; command = $expanded; target_label = $targetLabel; evidence_type = 'assumed' })
  }
  $forbiddenEffects = @($contract.side_effects.forbidden)
  if ('external-network' -in $forbiddenEffects -and @($usedLanes | Where-Object network -eq 'live').Count -gt 0) { return New-ContextResult $false 'CTX_SIDE_EFFECT_CONFLICT' 'external-network' }
  if ('provider-write' -in $forbiddenEffects -and @($usedLanes | Where-Object id -eq 'provider-write').Count -gt 0) { return New-ContextResult $false 'CTX_SIDE_EFFECT_CONFLICT' 'provider-write' }
  $sideEffects = @($contract.side_effects.allowed)
  if ('none' -in $sideEffects -and $sideEffects.Count -gt 1) { return New-ContextResult $false 'CTX_SIDE_EFFECT_CONFLICT' 'none' }

  $criteria = [Collections.Generic.List[object]]::new()
  foreach ($criterion in @($contract.criteria)) {
    if ($specContent -notmatch "(?m)^### $([regex]::Escape([string]$criterion.id))(?:\s|$)") { return New-ContextResult $false 'CTX_CRITERION_PROOF_MISSING' ([string]$criterion.id) }
    if ($validationContent -notmatch "(?m)^ID:\s*$([regex]::Escape([string]$criterion.milestone_criterion_id))\s*$") { return New-ContextResult $false 'CTX_PROOF_REFERENCE_INVALID' ([string]$criterion.id) }
    if (@($criterion.command_ids).Count -eq 0) { return New-ContextResult $false 'CTX_CRITERION_PROOF_MISSING' ([string]$criterion.id) }
    foreach ($proof in @($criterion.command_ids)) { if ($proof -notin $commandIds) { return New-ContextResult $false 'CTX_PROOF_REFERENCE_INVALID' ([string]$criterion.id) } }
    $criteria.Add([ordered]@{ id = [string]$criterion.id; milestone_id = [string]$criterion.milestone_criterion_id; proof_commands = @($criterion.command_ids) })
  }
  $seams = @(Get-SharedSeams -Registry $seamsDocument.Document -Paths $requestedPaths | Select-Object -ExpandProperty id | Sort-Object -Unique)
  $risk = Get-ContextRisk $requestedPaths $seams @($usedLanes) $sideEffects
  $sources = @($sourcePaths | ForEach-Object { [ordered]@{ path = Get-SafeContextPath $RepositoryRoot $_; sha256 = Get-ContextFileHash $_ } })
  $stopConditions = @($contract.stop_conditions | ForEach-Object { "$($_.code): $($_.condition)" })
  $packData = [ordered]@{
    schema_version = '1.0'
    context_id = ([IO.Path]::GetFileName($normalizedFeature).ToLowerInvariant() -replace '_', '-')
    feature_path = $normalizedFeature
    base_sha = $BaseSha
    objective = $objective
    observable_done = $done
    criteria = @($criteria)
    sources = $sources
    risk = $risk
    paths = [ordered]@{ allowed = $requestedPaths; forbidden = $forbidden }
    shared_seams = $seams
    side_effects = $sideEffects
    commands = @($commands)
    stop_conditions = $stopConditions
    retry_budget = [int]$contract.retry_budget.max_correction_attempts
    handoff = [ordered]@{ target = 'Milestone Orchestrator'; reason = "return fields: $(@($contract.handoff_fields) -join ', ')" }
  }
  $pack = [pscustomobject]$packData
  $estimate = Get-EstimatedContextTokens $pack
  if ($estimate -gt 2000) { return New-ContextResult $false 'CTX_TOKEN_BUDGET_EXCEEDED' 'estimated-input-tokens' }
  $pack | Add-Member -NotePropertyName estimated_input_tokens -NotePropertyValue $estimate
  New-ContextResult $true -Id ([string]$pack.context_id) -Path $ContextPath -Pack $pack
}

function New-HarnessContextPack {
  [CmdletBinding()]
  param([Parameter(Mandatory)][string]$FeaturePath, [Parameter(Mandatory)][string]$BaseSha, [Parameter(Mandatory)][string[]]$AllowedPath, [Parameter(Mandatory)][string]$OutputPath)
  $safeOutputPath = Get-SafeContextPath $script:RepositoryRoot $OutputPath
  $result = New-CanonicalContextPack $script:RepositoryRoot $FeaturePath $BaseSha $AllowedPath $safeOutputPath
  if (-not $result.Passed) { return $result }
  $outputDirectory = Split-Path -Parent $OutputPath
  if (-not [string]::IsNullOrWhiteSpace($outputDirectory)) { New-Item -ItemType Directory -Path $outputDirectory -Force | Out-Null }
  $result.Pack | ConvertTo-Json -Depth 100 | Set-Content -LiteralPath $OutputPath -Encoding utf8
  $result.Path = $safeOutputPath
  $result
}

function Test-HarnessContextPack {
  [CmdletBinding()]
  param([Parameter(Mandatory)][string]$Path, [Parameter(Mandatory)][string]$RepositoryRoot, [switch]$RequireCurrentBase)
  $safePath = Get-SafeContextPath $RepositoryRoot $Path
  if (-not (Test-Path -LiteralPath $Path -PathType Leaf)) { return New-ContextResult $false 'CTX_SOURCE_MISSING' 'context-pack' $safePath }
  try { $pack = Get-Content -Raw -LiteralPath $Path | ConvertFrom-Json -Depth 100 } catch { return New-ContextResult $false 'CTX_SCHEMA_INVALID' 'context-pack' $safePath }
  if ([string]$pack.base_sha -notmatch '^[0-9a-f]{40}$') { return New-ContextResult $false 'CTX_BASE_SHA_INVALID' 'base-sha' }
  if ($RequireCurrentBase) {
    $currentBase = @(& git -C $RepositoryRoot rev-parse HEAD 2>$null)
    if ($LASTEXITCODE -ne 0 -or $currentBase.Count -ne 1 -or $currentBase[0].Trim() -ne [string]$pack.base_sha) { return New-ContextResult $false 'CTX_BASE_SHA_MISMATCH' 'base-sha' }
  }
  foreach ($property in @('allowed', 'forbidden')) {
    foreach ($candidate in @($pack.paths.$property)) {
      $normalized = ConvertTo-ContextPath ([string]$candidate)
      if ($null -eq $normalized -or $normalized -ne [string]$candidate) { return New-ContextResult $false 'CTX_PATH_OUTSIDE_SCOPE' $property ([string]$candidate) }
    }
  }
  foreach ($criterion in @($pack.criteria)) {
    if (@($criterion.proof_commands).Count -eq 0) { return New-ContextResult $false 'CTX_CRITERION_PROOF_MISSING' ([string]$criterion.id) }
    foreach ($proof in @($criterion.proof_commands)) { if ($proof -notin @($pack.commands.id)) { return New-ContextResult $false 'CTX_PROOF_REFERENCE_INVALID' ([string]$criterion.id) } }
  }
  if ([int]$pack.estimated_input_tokens -gt 2000) { return New-ContextResult $false 'CTX_TOKEN_BUDGET_EXCEEDED' 'estimated-input-tokens' }
  foreach ($source in @($pack.sources)) {
    $normalized = ConvertTo-ContextPath ([string]$source.path)
    if ($null -eq $normalized -or $normalized -ne [string]$source.path) { return New-ContextResult $false 'CTX_PATH_OUTSIDE_SCOPE' 'source' ([string]$source.path) }
    $sourcePath = Resolve-ContextRepositoryPath $RepositoryRoot $normalized
    if (-not (Test-Path -LiteralPath $sourcePath -PathType Leaf)) { return New-ContextResult $false 'CTX_SOURCE_MISSING' 'source' $normalized }
    if ((Get-ContextFileHash $sourcePath) -ne [string]$source.sha256) { return New-ContextResult $false 'CTX_SOURCE_HASH_MISMATCH' 'source' $normalized }
  }
  $schemaPath = Join-Path $RepositoryRoot 'contracts/governance/schemas/context-pack.schema.json'
  try { if (-not (Test-Json -LiteralPath $Path -SchemaFile $schemaPath -ErrorAction Stop)) { throw 'schema rejected pack' } } catch { return New-ContextResult $false 'CTX_SCHEMA_INVALID' 'context-pack' $safePath }
  $expected = New-CanonicalContextPack $RepositoryRoot ([string]$pack.feature_path) ([string]$pack.base_sha) @($pack.paths.allowed) $safePath
  if (-not $expected.Passed) { return $expected }
  $canonical = $expected.Pack
  if (-not (Test-JsonEqual $pack.sources $canonical.sources)) { return New-ContextResult $false 'CTX_SOURCE_HASH_MISMATCH' 'sources' }
  if (-not (Test-JsonEqual $pack.criteria $canonical.criteria)) { return New-ContextResult $false 'CTX_PROOF_REFERENCE_INVALID' 'criteria' }
  if (-not (Test-JsonEqual $pack.paths $canonical.paths)) { return New-ContextResult $false 'CTX_PATH_SCOPE_CONFLICT' 'paths' }
  if (-not (Test-JsonEqual $pack.shared_seams $canonical.shared_seams)) { return New-ContextResult $false 'CTX_SHARED_SEAM_UNDECLARED' 'shared-seams' }
  if (-not (Test-JsonEqual $pack.side_effects $canonical.side_effects)) { return New-ContextResult $false 'CTX_SIDE_EFFECT_CONFLICT' 'side-effects' }
  foreach ($field in @('context_id', 'feature_path', 'objective', 'observable_done', 'risk', 'stop_conditions', 'retry_budget', 'handoff')) {
    if (-not (Test-JsonEqual $pack.$field $canonical.$field)) { return New-ContextResult $false 'CTX_FEATURE_INVALID' $field }
  }
  if (-not (Test-JsonEqual $pack.commands $canonical.commands)) { return New-ContextResult $false 'CTX_TARGET_EVIDENCE_INFLATION' 'commands' }
  if ([int]$pack.estimated_input_tokens -ne [int]$canonical.estimated_input_tokens) { return New-ContextResult $false 'CTX_SCHEMA_INVALID' 'estimated-input-tokens' }
  New-ContextResult $true -Id ([string]$pack.context_id) -Path $safePath -Pack $pack
}

Export-ModuleMember -Function New-HarnessContextPack, Test-HarnessContextPack
