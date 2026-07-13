$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$modulePath = Join-Path $repoRoot 'scripts/harness/Context.psm1'
$featurePath = '.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/F-10-pragmatic-harness-cutover'
$allowedPaths = @(
  'contracts/governance/**',
  'scripts/harness/**',
  'scripts/tests/**',
  'scripts/harness.ps1',
  'package.json',
  'AGENTS.md',
  '.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/F-10-pragmatic-harness-cutover/**'
)
$fixtureRoot = Join-Path ([IO.Path]::GetTempPath()) ("mpc-context-fixtures-{0}" -f [guid]::NewGuid().ToString('N'))

function Assert-Equal {
  param([object]$Actual, [object]$Expected, [string]$Case)
  if ($Actual -ne $Expected) { throw "$Case expected '$Expected', got '$Actual'" }
}

function Copy-Pack {
  param([object]$Pack)
  $Pack | ConvertTo-Json -Depth 100 | ConvertFrom-Json -Depth 100
}

function Write-Pack {
  param([object]$Pack, [string]$Name)
  $path = Join-Path $fixtureRoot "$Name.json"
  $Pack | ConvertTo-Json -Depth 100 | Set-Content -LiteralPath $path -Encoding utf8
  return $path
}

function Assert-PackFailure {
  param([object]$Pack, [string]$Name, [string]$Code, [switch]$RequireCurrentBase)
  $path = Write-Pack -Pack $Pack -Name $Name
  $result = Test-HarnessContextPack -Path $path -RepositoryRoot $repoRoot -RequireCurrentBase:$RequireCurrentBase
  Assert-Equal $result.Passed $false $Name
  Assert-Equal $result.ErrorCode $Code $Name
}

function Set-FixtureContent {
  param([string]$Path, [string]$Content)
  $directory = Split-Path -Parent $Path
  New-Item -ItemType Directory -Path $directory -Force | Out-Null
  Set-Content -LiteralPath $Path -Value $Content -Encoding utf8
}

function New-GenericFeatureFixture {
  param(
    [string]$Name,
    [string]$Objective = 'Compile a generic bounded context.',
    [string]$ValidationContent = "# Validation`n`n## Criterion: Generic`nID: M-99-C01`n"
  )
  $missionRoot = Join-Path $repoRoot "scripts/.runs/$Name/MIS-999-generic"
  $milestoneRoot = Join-Path $missionRoot 'M-99-generic'
  $featureRoot = Join-Path $milestoneRoot 'F-42-generic-context'
  $featureRelative = [IO.Path]::GetRelativePath($repoRoot, $featureRoot).Replace('\', '/')
  $notePath = Join-Path $featureRoot 'notes.md'
  $noteRelative = [IO.Path]::GetRelativePath($repoRoot, $notePath).Replace('\', '/')
  Set-FixtureContent (Join-Path $missionRoot 'mission.md') "# Generic Mission`n"
  Set-FixtureContent (Join-Path $milestoneRoot 'milestone.md') "# Generic Milestone`n"
  Set-FixtureContent (Join-Path $milestoneRoot 'validation-contract.md') $ValidationContent
  Set-FixtureContent $notePath "generic source`n"
  Set-FixtureContent (Join-Path $featureRoot 'feature.md') "# F-42`n`n## Brief`n`n$Objective`n`n## Expected Output`n`n- A generic pack validates.`n"
  Set-FixtureContent (Join-Path $featureRoot 'spec.md') "# Generic Spec`n`n## Acceptance Criteria`n`n### F42-AC01 — Generic proof`n`n- Traces to milestone criterion ID: ``M-99-C01``.`n"
  $contract = [ordered]@{
    schema_version = '1.0'
    feature_id = 'F-42'
    required_sources = @($noteRelative)
    allowed_paths = @("$featureRelative/**")
    forbidden_paths = @('apps/**')
    side_effects = [ordered]@{ allowed = @(); forbidden = @('database-mutation', 'external-network', 'provider-write') }
    commands = @([ordered]@{ id = 'generic-proof'; command_id = 'impact-probe-one'; lane_id = 'unit'; expected_exit_code = 0 })
    criteria = @([ordered]@{ id = 'F42-AC01'; milestone_criterion_id = 'M-99-C01'; command_ids = @('generic-proof') })
    stop_conditions = @([ordered]@{ code = 'generic-stop'; condition = 'Stop on generic contract drift.' })
    retry_budget = [ordered]@{ max_correction_attempts = 1 }
    handoff_fields = @('status', 'evidence', 'next')
  }
  $plan = "# Generic Plan`n`n## Machine Work Contract`n`n``````json`n$($contract | ConvertTo-Json -Depth 20)`n```````n"
  Set-FixtureContent (Join-Path $featureRoot 'plan.md') $plan
  return $featureRelative
}

New-Item -ItemType Directory -Path $fixtureRoot -Force | Out-Null
try {
  if (-not (Test-Path -LiteralPath $modulePath -PathType Leaf)) {
    throw 'Context.psm1 functions do not exist'
  }
  Import-Module $modulePath -Force

  $baseSha = (& git -C $repoRoot rev-parse HEAD).Trim()
  $positivePath = Join-Path $fixtureRoot 'positive.json'
  $compiled = New-HarnessContextPack -FeaturePath $featurePath -BaseSha $baseSha -AllowedPath $allowedPaths -OutputPath $positivePath
  Assert-Equal $compiled.Passed $true 'positive compile'
  $validated = Test-HarnessContextPack -Path $positivePath -RepositoryRoot $repoRoot -RequireCurrentBase
  Assert-Equal $validated.Passed $true 'positive validation'
  $positive = Get-Content -Raw -LiteralPath $positivePath | ConvertFrom-Json -Depth 100
  Assert-Equal $positive.risk.level 'L2' 'positive risk'
  Assert-Equal (@($positive.commands | Where-Object target_label -ne 'fake').Count) 0 'positive targets'
  if ([int]$positive.estimated_input_tokens -gt 2000) {
    Assert-Equal $positive.risk.level 'L2' 'overflow risk level'
    Assert-Equal (@($positive.sources | Where-Object { [string]::IsNullOrWhiteSpace([string]$_.overflow_reason) }).Count) 0 'overflow reasons'
  }

  $invalidSha = New-HarnessContextPack -FeaturePath $featurePath -BaseSha 'abc' -AllowedPath $allowedPaths -OutputPath (Join-Path $fixtureRoot 'invalid-sha.json')
  Assert-Equal $invalidSha.ErrorCode 'CTX_BASE_SHA_INVALID' 'invalid SHA'

  $case = Copy-Pack $positive; $case.base_sha = '0000000000000000000000000000000000000000'
  Assert-PackFailure $case 'stale-head' 'CTX_BASE_SHA_MISMATCH' -RequireCurrentBase

  $case = Copy-Pack $positive; $case.sources[0].path = 'missing/source.md'
  Assert-PackFailure $case 'missing-source' 'CTX_SOURCE_MISSING'

  $case = Copy-Pack $positive; $case.sources[0].sha256 = ('0' * 64)
  Assert-PackFailure $case 'mutated-source' 'CTX_SOURCE_HASH_MISMATCH'

  $case = Copy-Pack $positive; $case.objective = 'tampered objective'
  Assert-PackFailure $case 'tampered-objective' 'CTX_FEATURE_INVALID'

  $case = Copy-Pack $positive; $case.risk.level = 'L0'
  Assert-PackFailure $case 'tampered-risk' 'CTX_TOKEN_BUDGET_EXCEEDED'

  $case = Copy-Pack $positive; $case.criteria[0].proof_commands = @()
  Assert-PackFailure $case 'missing-proof' 'CTX_CRITERION_PROOF_MISSING'

  $case = Copy-Pack $positive; $case.criteria[0].proof_commands = @('not-a-command')
  Assert-PackFailure $case 'dangling-proof' 'CTX_PROOF_REFERENCE_INVALID'

  $case = Copy-Pack $positive; $case.paths.forbidden = @($case.paths.allowed[0])
  Assert-PackFailure $case 'path-overlap' 'CTX_PATH_SCOPE_CONFLICT'

  $case = Copy-Pack $positive; $case.paths.allowed = @($case.paths.allowed + 'apps/server_core')
  Assert-PackFailure $case 'outside-scope' 'CTX_PATH_OUTSIDE_SCOPE'

  $ancestor = New-HarnessContextPack -FeaturePath $featurePath -BaseSha $baseSha -AllowedPath @('.mnfs/**') -OutputPath (Join-Path $fixtureRoot 'ancestor.json')
  Assert-Equal $ancestor.ErrorCode 'CTX_PATH_OUTSIDE_SCOPE' 'ancestor scope escalation'

  $case = Copy-Pack $positive; $case.paths.allowed = @($case.paths.allowed + 'C:/rooted/path')
  Assert-PackFailure $case 'rooted-path' 'CTX_PATH_OUTSIDE_SCOPE'

  $case = Copy-Pack $positive; $case.paths.allowed = @($case.paths.allowed + '../parent')
  Assert-PackFailure $case 'parent-path' 'CTX_PATH_OUTSIDE_SCOPE'

  $case = Copy-Pack $positive; $case.shared_seams = @($case.shared_seams | Where-Object { $_ -ne 'dependency-graph' })
  Assert-PackFailure $case 'undeclared-seam' 'CTX_SHARED_SEAM_UNDECLARED'

  $case = Copy-Pack $positive; $case.side_effects = @('provider-write')
  Assert-PackFailure $case 'side-effect-conflict' 'CTX_SIDE_EFFECT_CONFLICT'

  $case = Copy-Pack $positive; $case.commands[0].target_label = 'live-provider'; $case.commands[0].evidence_type = 'ran'
  Assert-PackFailure $case 'target-inflation' 'CTX_TARGET_EVIDENCE_INFLATION'

  $case = Copy-Pack $positive; $case.risk.level = 'L0'; $case.stop_conditions = @($case.stop_conditions + ('x' * 8100)); $case.estimated_input_tokens = 2001
  Assert-PackFailure $case 'token-budget' 'CTX_TOKEN_BUDGET_EXCEEDED'

  $genericFeature = New-GenericFeatureFixture -Name "generic-$([guid]::NewGuid().ToString('N'))"
  $genericPath = Join-Path $fixtureRoot 'generic.json'
  $generic = New-HarnessContextPack -FeaturePath $genericFeature -BaseSha $baseSha -AllowedPath @("$genericFeature/**") -OutputPath $genericPath
  Assert-Equal $generic.Passed $true 'generic fixture compile'
  $genericPack = Get-Content -Raw -LiteralPath $genericPath | ConvertFrom-Json -Depth 100
  Assert-Equal $genericPack.criteria[0].id 'F42-AC01' 'generic criterion'
  Assert-Equal $genericPack.commands[0].id 'generic-proof' 'generic command'
  Assert-Equal $genericPack.commands[0].command_id 'impact-probe-one' 'generic command registry ID'

  $headingFeature = New-GenericFeatureFixture -Name "heading-$([guid]::NewGuid().ToString('N'))" -ValidationContent "# Validation`n`n### M-99-C01 — Generic criterion`n"
  $heading = New-HarnessContextPack -FeaturePath $headingFeature -BaseSha $baseSha -AllowedPath @("$headingFeature/**") -OutputPath (Join-Path $fixtureRoot 'heading.json')
  Assert-Equal $heading.Passed $true 'heading criterion compile'

  $inlineFeature = New-GenericFeatureFixture -Name "inline-$([guid]::NewGuid().ToString('N'))" -ValidationContent "# Validation`n`nThe criterion M-99-C01 appears only in inline prose.`n"
  $inline = New-HarnessContextPack -FeaturePath $inlineFeature -BaseSha $baseSha -AllowedPath @("$inlineFeature/**") -OutputPath (Join-Path $fixtureRoot 'inline.json')
  Assert-Equal $inline.ErrorCode 'CTX_PROOF_REFERENCE_INVALID' 'inline criterion rejected'

  $utf8Feature = New-GenericFeatureFixture -Name "utf8-$([guid]::NewGuid().ToString('N'))" -Objective ([string]::new([char]0x00E9, 8100))
  $utf8 = New-HarnessContextPack -FeaturePath $utf8Feature -BaseSha $baseSha -AllowedPath @("$utf8Feature/**") -OutputPath (Join-Path $fixtureRoot 'utf8.json')
  Assert-Equal $utf8.ErrorCode 'CTX_TOKEN_BUDGET_EXCEEDED' 'UTF-8 multibyte estimate'

  $missingRelative = 'scripts/.runs/missing-context-pack.json'
  $dispatch = @(& pwsh -NoProfile -ExecutionPolicy Bypass -File (Join-Path $repoRoot 'scripts/harness.ps1') -Command context-validate -ContextPath $missingRelative 2>&1)
  Assert-Equal $LASTEXITCODE 1 'missing relative dispatcher exit'
  Assert-Equal (@($dispatch | Where-Object { $_ -eq "path=$missingRelative" }).Count) 1 'missing relative dispatcher path'
  if (($dispatch -join "`n") -match '(?i)Users[/\\]|leandro') { throw 'dispatcher exposed an absolute user path' }

  Write-Output 'PASS context compiler tests'
} finally {
  Get-ChildItem -LiteralPath (Join-Path $repoRoot 'scripts/.runs') -Directory -ErrorAction SilentlyContinue |
    Where-Object Name -match '^(?:generic|heading|inline|utf8)-' |
    ForEach-Object { Remove-Item -LiteralPath $_.FullName -Recurse -Force -ErrorAction SilentlyContinue }
  Remove-Item -LiteralPath $fixtureRoot -Recurse -Force -ErrorAction SilentlyContinue
}
