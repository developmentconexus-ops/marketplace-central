$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$modulePath = Join-Path $repoRoot 'scripts/harness/Context.psm1'
$featurePath = '.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/F-07-governance-context-compiler'
$allowedPaths = @('contracts/governance/**', 'scripts/harness/**', 'scripts/tests/**', 'scripts/harness.ps1', 'package.json', 'AGENTS.md')
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
  if ([int]$positive.estimated_input_tokens -gt 2000) { throw 'positive estimate exceeds 2000' }

  $invalidSha = New-HarnessContextPack -FeaturePath $featurePath -BaseSha 'abc' -AllowedPath $allowedPaths -OutputPath (Join-Path $fixtureRoot 'invalid-sha.json')
  Assert-Equal $invalidSha.ErrorCode 'CTX_BASE_SHA_INVALID' 'invalid SHA'

  $case = Copy-Pack $positive; $case.base_sha = '0000000000000000000000000000000000000000'
  Assert-PackFailure $case 'stale-head' 'CTX_BASE_SHA_MISMATCH' -RequireCurrentBase

  $case = Copy-Pack $positive; $case.sources[0].path = 'missing/source.md'
  Assert-PackFailure $case 'missing-source' 'CTX_SOURCE_MISSING'

  $case = Copy-Pack $positive; $case.sources[0].sha256 = ('0' * 64)
  Assert-PackFailure $case 'mutated-source' 'CTX_SOURCE_HASH_MISMATCH'

  $case = Copy-Pack $positive; $case.criteria[0].proof_commands = @()
  Assert-PackFailure $case 'missing-proof' 'CTX_CRITERION_PROOF_MISSING'

  $case = Copy-Pack $positive; $case.criteria[0].proof_commands = @('not-a-command')
  Assert-PackFailure $case 'dangling-proof' 'CTX_PROOF_REFERENCE_INVALID'

  $case = Copy-Pack $positive; $case.paths.forbidden = @($case.paths.allowed[0])
  Assert-PackFailure $case 'path-overlap' 'CTX_PATH_SCOPE_CONFLICT'

  $case = Copy-Pack $positive; $case.paths.allowed = @($case.paths.allowed + 'apps/server_core')
  Assert-PackFailure $case 'outside-scope' 'CTX_PATH_OUTSIDE_SCOPE'

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

  $case = Copy-Pack $positive; $case.stop_conditions = @($case.stop_conditions + ('x' * 8100)); $case.estimated_input_tokens = 2001
  Assert-PackFailure $case 'token-budget' 'CTX_TOKEN_BUDGET_EXCEEDED'

  Write-Output 'PASS context compiler tests'
} finally {
  Remove-Item -LiteralPath $fixtureRoot -Recurse -Force -ErrorAction SilentlyContinue
}
