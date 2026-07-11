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
  param([string]$Path)
  if ([string]::IsNullOrWhiteSpace($Path) -or [IO.Path]::IsPathRooted($Path) -or $Path -match '(^|[\\/])\.\.([\\/]|$)') { return $null }
  $normalized = $Path.Trim().Replace('\', '/').TrimEnd('/')
  if ($normalized.StartsWith('./', [StringComparison]::Ordinal)) { $normalized = $normalized.Substring(2) }
  $normalized = $normalized -replace '/\*\*$', '' -replace '/\*$', ''
  if ([string]::IsNullOrWhiteSpace($normalized) -or $normalized -match '[:*?\[]') { return $null }
  return $normalized
}

function Test-ContextPathIntersection {
  param([string]$Left, [string]$Right)
  $leftPath = $Left.TrimEnd('/'); $rightPath = $Right.TrimEnd('/')
  return $leftPath -eq $rightPath -or $leftPath.StartsWith($rightPath + '/', [StringComparison]::Ordinal) -or $rightPath.StartsWith($leftPath + '/', [StringComparison]::Ordinal)
}

function Resolve-ContextRepositoryPath {
  param([string]$RepositoryRoot, [string]$Path)
  Join-Path $RepositoryRoot ($Path.Replace('/', [IO.Path]::DirectorySeparatorChar))
}

function Get-ContextFileHash {
  param([string]$Path)
  ([Convert]::ToHexString([Security.Cryptography.SHA256]::HashData([IO.File]::ReadAllBytes($Path)))).ToLowerInvariant()
}

function Get-MarkdownSection {
  param([string]$Content, [string]$Heading)
  $match = [regex]::Match($Content, "(?ms)^## $([regex]::Escape($Heading))\s*\r?\n(?<body>.*?)(?=^## |\z)")
  if (-not $match.Success) { return '' }
  return $match.Groups['body'].Value.Trim()
}

function ConvertTo-CompactMarkdownText {
  param([string]$Text)
  $value = $Text -replace '(?m)^\s*[-*]\s+', '' -replace '`', '' -replace '\s+', ' '
  return $value.Trim()
}

function Get-ObservableDone {
  param([string]$FeatureContent)
  $section = Get-MarkdownSection $FeatureContent 'Expected Output'
  if ([string]::IsNullOrWhiteSpace($section)) { $section = Get-MarkdownSection $FeatureContent 'Validation Expectations' }
  $items = @([regex]::Matches($section, '(?m)^-\s+(?<item>.+(?:\r?\n\s{2,}.+)*)') | ForEach-Object { ConvertTo-CompactMarkdownText $_.Groups['item'].Value })
  if ($items.Count -eq 0 -and -not [string]::IsNullOrWhiteSpace($section)) { $items = @(ConvertTo-CompactMarkdownText $section) }
  return $items
}

function Get-PlanDeclaredPaths {
  param([string]$PlanContent)
  $paths = [Collections.Generic.List[string]]::new()
  foreach ($match in [regex]::Matches($PlanContent, '`(?<path>(?:AGENTS\.md|package\.json|scripts/[^` <]+|contracts/[^` <]+|\.mnfs/[^` <]+))`')) {
    $normalized = ConvertTo-ContextPath $match.Groups['path'].Value
    if ($null -ne $normalized) { $paths.Add($normalized) }
  }
  @($paths | Sort-Object -Unique)
}

function Get-ContextCommandCatalog {
  param([string]$PlanContent, [string]$BaseSha)
  $mapping = Get-MarkdownSection $PlanContent 'Verification Mapping'
  $definitions = @(
    [pscustomobject]@{ Id = 'governance-contracts'; Needle = 'governance-contracts.tests.ps1'; Command = 'pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/governance-contracts.tests.ps1' },
    [pscustomobject]@{ Id = 'governance-drift'; Needle = 'governance-drift.tests.ps1'; Command = 'pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/governance-drift.tests.ps1' },
    [pscustomobject]@{ Id = 'governance-current'; Needle = 'harness:governance'; Command = "npm run harness:governance -- -BaseSha $BaseSha" },
    [pscustomobject]@{ Id = 'context-compiler'; Needle = 'context-compiler.tests.ps1'; Command = 'pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/context-compiler.tests.ps1' },
    [pscustomobject]@{ Id = 'context-current'; Needle = 'current F-07 compile/validate'; Command = 'npm run harness:context:validate -- -ContextPath <context-pack> -RequireCurrentBase' }
  )
  @($definitions | Where-Object { $mapping.Contains($_.Needle, [StringComparison]::OrdinalIgnoreCase) } | ForEach-Object {
    [ordered]@{ id = $_.Id; command = $_.Command; target_label = 'fake'; evidence_type = 'assumed' }
  })
}

function Get-ContextCriteria {
  param([string]$SpecContent, [string]$PlanContent, [string]$ValidationContent, [object[]]$Commands)
  $mapping = Get-MarkdownSection $PlanContent 'Verification Mapping'
  $criteria = [Collections.Generic.List[object]]::new()
  foreach ($match in [regex]::Matches($SpecContent, '(?ms)^### (?<id>F\d{2}-AC\d{2})[^\r\n]*\r?\n(?<body>.*?)(?=^### |^## |\z)')) {
    $id = $match.Groups['id'].Value
    $body = $match.Groups['body'].Value
    $milestoneMatch = [regex]::Match($body, '`(?<id>M-\d{2}-C\d{2})`')
    $proofLine = [regex]::Match($mapping, "(?m)^- $([regex]::Escape($id)):\s*(?<proof>.+)$")
    $proofIds = [Collections.Generic.List[string]]::new()
    if ($proofLine.Success) {
      foreach ($command in $Commands) {
        $needle = switch ($command.id) {
          'governance-contracts' { 'governance-contracts.tests.ps1' }
          'governance-drift' { 'governance-drift.tests.ps1' }
          'governance-current' { 'harness:governance' }
          'context-compiler' { 'context' }
          'context-current' { 'current F-07 compile/validate' }
        }
        if ($proofLine.Groups['proof'].Value.Contains($needle, [StringComparison]::OrdinalIgnoreCase)) { $proofIds.Add([string]$command.id) }
      }
    }
    $milestoneId = if ($milestoneMatch.Success) { $milestoneMatch.Groups['id'].Value } else { '' }
    if (-not [string]::IsNullOrWhiteSpace($milestoneId) -and $ValidationContent -notmatch "(?m)^ID:\s*$([regex]::Escape($milestoneId))\s*$") { $milestoneId = '' }
    $criteria.Add([ordered]@{ id = $id; milestone_id = $milestoneId; proof_commands = @($proofIds | Sort-Object -Unique) })
  }
  return @($criteria)
}

function Get-ContextRisk {
  param([string[]]$AllowedPaths, [object[]]$SharedSeams)
  $moduleScopes = @($AllowedPaths | Where-Object { $_ -match '^apps/server_core/internal/modules/(?<module>[^/]+)' } | ForEach-Object { if ($_ -match '^apps/server_core/internal/modules/(?<module>[^/]+)') { $Matches.module } } | Sort-Object -Unique)
  $level = if (@($AllowedPaths | Where-Object { $_ -match '(?:^|/)(?:external|live|provider-write)(?:/|$)' }).Count -gt 0) { 'L3' } elseif ($SharedSeams.Count -gt 0 -or $moduleScopes.Count -gt 1) { 'L2' } elseif ($moduleScopes.Count -eq 1) { 'L1' } else { 'L0' }
  $policy = @{ L0 = 'self-review'; L1 = 'standard-review'; L2 = 'independent-review'; L3 = 'safety-review' }[$level]
  $model = @{ L0 = 'fast-capable-model'; L1 = 'general-reasoning-model'; L2 = 'strong-reasoning-model'; L3 = 'strongest-safety-reviewed-model' }[$level]
  [ordered]@{ level = $level; review_policy = $policy; advisory_model = $model }
}

function Get-EstimatedContextTokens {
  param([object]$Pack)
  $withoutEstimate = [ordered]@{}
  foreach ($property in $Pack.PSObject.Properties) {
    if ($property.Name -ne 'estimated_input_tokens') { $withoutEstimate[$property.Name] = $property.Value }
  }
  [int][Math]::Ceiling(($withoutEstimate | ConvertTo-Json -Depth 100 -Compress).Length / 4.0)
}

function New-HarnessContextPack {
  [CmdletBinding()]
  param(
    [Parameter(Mandatory)][string]$FeaturePath,
    [Parameter(Mandatory)][string]$BaseSha,
    [Parameter(Mandatory)][string[]]$AllowedPath,
    [Parameter(Mandatory)][string]$OutputPath
  )

  if ($BaseSha -notmatch '^[0-9a-f]{40}$') { return New-ContextResult $false 'CTX_BASE_SHA_INVALID' 'base-sha' }
  $normalizedFeaturePath = ConvertTo-ContextPath $FeaturePath
  if ($null -eq $normalizedFeaturePath) { return New-ContextResult $false 'CTX_FEATURE_INVALID' 'feature-path' $FeaturePath }
  $normalizedAllowed = [Collections.Generic.List[string]]::new()
  foreach ($path in $AllowedPath) {
    $normalized = ConvertTo-ContextPath $path
    if ($null -eq $normalized) { return New-ContextResult $false 'CTX_PATH_OUTSIDE_SCOPE' 'allowed-path' $path }
    $normalizedAllowed.Add($normalized)
  }
  $normalizedAllowed = @($normalizedAllowed | Sort-Object -Unique)

  $featureDirectory = Resolve-ContextRepositoryPath $script:RepositoryRoot $normalizedFeaturePath
  $milestoneDirectory = Split-Path -Parent $featureDirectory
  $sourceFiles = @(
    Join-Path $featureDirectory 'feature.md'
    Join-Path $featureDirectory 'spec.md'
    Join-Path $featureDirectory 'plan.md'
    Join-Path $milestoneDirectory 'validation-contract.md'
  )
  if (@($sourceFiles | Where-Object { -not (Test-Path -LiteralPath $_ -PathType Leaf) }).Count -gt 0) { return New-ContextResult $false 'CTX_FEATURE_INVALID' 'feature-sources' $normalizedFeaturePath }

  $featureContent = Get-Content -Raw -LiteralPath $sourceFiles[0]
  $specContent = Get-Content -Raw -LiteralPath $sourceFiles[1]
  $planContent = Get-Content -Raw -LiteralPath $sourceFiles[2]
  $validationContent = Get-Content -Raw -LiteralPath $sourceFiles[3]
  $objective = ConvertTo-CompactMarkdownText (Get-MarkdownSection $featureContent 'Brief')
  $done = @(Get-ObservableDone $featureContent)
  $commands = @(Get-ContextCommandCatalog $planContent $BaseSha)
  $criteria = @(Get-ContextCriteria $specContent $planContent $validationContent $commands)
  if ([string]::IsNullOrWhiteSpace($objective) -or $done.Count -eq 0 -or $criteria.Count -eq 0) { return New-ContextResult $false 'CTX_FEATURE_INVALID' 'mnfs-headings' $normalizedFeaturePath }
  if (@($criteria | Where-Object { [string]::IsNullOrWhiteSpace($_.milestone_id) -or @($_.proof_commands).Count -eq 0 }).Count -gt 0) { return New-ContextResult $false 'CTX_CRITERION_PROOF_MISSING' 'criterion-proof' $normalizedFeaturePath }

  Import-Module (Join-Path $PSScriptRoot 'Policy.psm1') -Force
  $governance = Test-GovernanceContracts -RepositoryRoot $script:RepositoryRoot
  if (-not $governance.Passed) { return New-ContextResult $false 'CTX_FEATURE_INVALID' 'governance-contracts' 'contracts/governance' }
  $seams = @(Get-SharedSeams -Registry $governance.Documents.'shared-seams' -Paths $normalizedAllowed | Select-Object -ExpandProperty id | Sort-Object -Unique)
  $risk = Get-ContextRisk $normalizedAllowed $seams
  $retryMatch = [regex]::Match($validationContent, '(?m)^- max_correction_attempts:\s*(?<count>[0-3])\s*$')
  $retryBudget = if ($retryMatch.Success) { [int]$retryMatch.Groups['count'].Value } else { 0 }
  $sources = @($sourceFiles | ForEach-Object {
    $relative = [IO.Path]::GetRelativePath($script:RepositoryRoot, $_).Replace('\', '/')
    [ordered]@{ path = $relative; sha256 = Get-ContextFileHash $_ }
  })

  $packData = [ordered]@{
    schema_version = '1.0'
    context_id = ([IO.Path]::GetFileName($normalizedFeaturePath).ToLowerInvariant() -replace '_', '-')
    feature_path = $normalizedFeaturePath
    base_sha = $BaseSha
    objective = $objective
    observable_done = $done
    criteria = $criteria
    sources = $sources
    risk = $risk
    paths = [ordered]@{ allowed = $normalizedAllowed; forbidden = @('apps/server_core', 'contracts/api', 'packages/sdk-runtime') }
    shared_seams = $seams
    side_effects = @('repository-write')
    commands = $commands
    stop_conditions = @('base or declared source changes before dispatch', 'requested path exceeds declared feature scope', 'undeclared shared seam or contract contradiction appears')
    retry_budget = $retryBudget
    handoff = [ordered]@{ target = 'Milestone Orchestrator'; reason = 'review criterion-complete feature evidence' }
  }
  $pack = [pscustomobject]$packData
  $estimate = Get-EstimatedContextTokens $pack
  if ($estimate -gt 2000) { return New-ContextResult $false 'CTX_TOKEN_BUDGET_EXCEEDED' 'estimated-input-tokens' }
  $pack | Add-Member -NotePropertyName estimated_input_tokens -NotePropertyValue $estimate
  $outputDirectory = Split-Path -Parent $OutputPath
  if (-not [string]::IsNullOrWhiteSpace($outputDirectory)) { New-Item -ItemType Directory -Path $outputDirectory -Force | Out-Null }
  $pack | ConvertTo-Json -Depth 100 | Set-Content -LiteralPath $OutputPath -Encoding utf8
  return New-ContextResult $true -Id ([string]$pack.context_id) -Path $OutputPath -Pack $pack
}

function Test-HarnessContextPack {
  [CmdletBinding()]
  param(
    [Parameter(Mandatory)][string]$Path,
    [Parameter(Mandatory)][string]$RepositoryRoot,
    [switch]$RequireCurrentBase
  )

  if (-not (Test-Path -LiteralPath $Path -PathType Leaf)) { return New-ContextResult $false 'CTX_SOURCE_MISSING' 'context-pack' $Path }
  try { $pack = Get-Content -Raw -LiteralPath $Path | ConvertFrom-Json -Depth 100 } catch { return New-ContextResult $false 'CTX_SCHEMA_INVALID' 'context-pack' $Path }
  if ([string]$pack.base_sha -notmatch '^[0-9a-f]{40}$') { return New-ContextResult $false 'CTX_BASE_SHA_INVALID' 'base-sha' }

  foreach ($property in @('allowed', 'forbidden')) {
    foreach ($candidate in @($pack.paths.$property)) {
      $normalized = ConvertTo-ContextPath ([string]$candidate)
      if ($null -eq $normalized -or $normalized -ne [string]$candidate) { return New-ContextResult $false 'CTX_PATH_OUTSIDE_SCOPE' $property ([string]$candidate) }
    }
  }
  foreach ($source in @($pack.sources)) {
    $normalized = ConvertTo-ContextPath ([string]$source.path)
    if ($null -eq $normalized -or $normalized -ne [string]$source.path) { return New-ContextResult $false 'CTX_PATH_OUTSIDE_SCOPE' 'source' ([string]$source.path) }
  }
  foreach ($criterion in @($pack.criteria)) {
    if (@($criterion.proof_commands).Count -eq 0) { return New-ContextResult $false 'CTX_CRITERION_PROOF_MISSING' ([string]$criterion.id) }
  }
  if ([int]$pack.estimated_input_tokens -gt 2000) { return New-ContextResult $false 'CTX_TOKEN_BUDGET_EXCEEDED' 'estimated-input-tokens' }

  if ($RequireCurrentBase) {
    $currentBase = @(& git -C $RepositoryRoot rev-parse HEAD 2>$null)
    if ($LASTEXITCODE -ne 0 -or $currentBase.Count -ne 1 -or $currentBase[0].Trim() -ne [string]$pack.base_sha) { return New-ContextResult $false 'CTX_BASE_SHA_MISMATCH' 'base-sha' }
  }
  foreach ($source in @($pack.sources)) {
    $sourcePath = Resolve-ContextRepositoryPath $RepositoryRoot ([string]$source.path)
    if (-not (Test-Path -LiteralPath $sourcePath -PathType Leaf)) { return New-ContextResult $false 'CTX_SOURCE_MISSING' 'source' ([string]$source.path) }
    if ((Get-ContextFileHash $sourcePath) -ne [string]$source.sha256) { return New-ContextResult $false 'CTX_SOURCE_HASH_MISMATCH' 'source' ([string]$source.path) }
  }

  $planSource = @($pack.sources | Where-Object path -like '*/plan.md' | Select-Object -First 1)
  if ($planSource.Count -ne 1) { return New-ContextResult $false 'CTX_FEATURE_INVALID' 'plan-source' }
  $planContent = Get-Content -Raw -LiteralPath (Resolve-ContextRepositoryPath $RepositoryRoot ([string]$planSource[0].path))
  $declaredPaths = @(Get-PlanDeclaredPaths $planContent)
  foreach ($allowed in @($pack.paths.allowed)) {
    if (@($declaredPaths | Where-Object { Test-ContextPathIntersection ([string]$allowed) ([string]$_) }).Count -eq 0) { return New-ContextResult $false 'CTX_PATH_OUTSIDE_SCOPE' 'allowed-path' ([string]$allowed) }
  }
  foreach ($allowed in @($pack.paths.allowed)) {
    if (@($pack.paths.forbidden | Where-Object { Test-ContextPathIntersection ([string]$allowed) ([string]$_) }).Count -gt 0) { return New-ContextResult $false 'CTX_PATH_SCOPE_CONFLICT' 'allowed-forbidden' ([string]$allowed) }
  }

  Import-Module (Join-Path $PSScriptRoot 'Policy.psm1') -Force
  $seamDocument = Import-GovernanceDocument -Path (Join-Path $RepositoryRoot 'contracts/governance/shared-seams.json') -SchemaPath (Join-Path $RepositoryRoot 'contracts/governance/schemas/shared-seams.schema.json')
  if (-not $seamDocument.Passed) { return New-ContextResult $false 'CTX_FEATURE_INVALID' 'shared-seams' }
  $expectedSeams = @(Get-SharedSeams -Registry $seamDocument.Document -Paths @($pack.paths.allowed) | Select-Object -ExpandProperty id | Sort-Object -Unique)
  $actualSeams = @($pack.shared_seams | Sort-Object -Unique)
  if (($expectedSeams -join '|') -ne ($actualSeams -join '|')) { return New-ContextResult $false 'CTX_SHARED_SEAM_UNDECLARED' 'shared-seams' }

  $commandIds = @($pack.commands.id)
  foreach ($criterion in @($pack.criteria)) {
    foreach ($proof in @($criterion.proof_commands)) {
      if ($proof -notin $commandIds) { return New-ContextResult $false 'CTX_PROOF_REFERENCE_INVALID' ([string]$criterion.id) }
    }
  }
  if ('provider-write' -in @($pack.side_effects) -and @($pack.commands | Where-Object target_label -eq 'provider-write').Count -eq 0) { return New-ContextResult $false 'CTX_SIDE_EFFECT_CONFLICT' 'provider-write' }
  if ('database-write' -in @($pack.side_effects) -and @($pack.commands | Where-Object target_label -eq 'ephemeral-postgres').Count -eq 0) { return New-ContextResult $false 'CTX_SIDE_EFFECT_CONFLICT' 'database-write' }
  if ('provider-read' -in @($pack.side_effects) -and @($pack.commands | Where-Object target_label -eq 'live-provider').Count -eq 0) { return New-ContextResult $false 'CTX_SIDE_EFFECT_CONFLICT' 'provider-read' }
  if ('browser' -in @($pack.side_effects) -and @($pack.commands | Where-Object target_label -eq 'browser').Count -eq 0) { return New-ContextResult $false 'CTX_SIDE_EFFECT_CONFLICT' 'browser' }
  foreach ($command in @($pack.commands)) {
    if ($command.target_label -ne 'fake' -and [string]$command.command -notmatch '(?i)harness\.ps1.+(?:integration|live|browser|provider-write)|harness:(?:integration|live|browser|provider-write)') { return New-ContextResult $false 'CTX_TARGET_EVIDENCE_INFLATION' ([string]$command.id) }
  }

  $schemaPath = Join-Path $RepositoryRoot 'contracts/governance/schemas/context-pack.schema.json'
  try { if (-not (Test-Json -LiteralPath $Path -SchemaFile $schemaPath -ErrorAction Stop)) { throw 'schema rejected pack' } } catch { return New-ContextResult $false 'CTX_SCHEMA_INVALID' 'context-pack' $Path }
  $actualEstimate = Get-EstimatedContextTokens $pack
  if ($actualEstimate -gt 2000) { return New-ContextResult $false 'CTX_TOKEN_BUDGET_EXCEEDED' 'estimated-input-tokens' }
  if ($actualEstimate -ne [int]$pack.estimated_input_tokens) { return New-ContextResult $false 'CTX_SCHEMA_INVALID' 'estimated-input-tokens' }
  return New-ContextResult $true -Id ([string]$pack.context_id) -Path $Path -Pack $pack
}

Export-ModuleMember -Function New-HarnessContextPack, Test-HarnessContextPack
