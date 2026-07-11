$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$contextModule = Join-Path $repoRoot 'scripts/harness/Context.psm1'
$impactModule = Join-Path $repoRoot 'scripts/harness/Impact.psm1'
$fixtureRoot = Join-Path $repoRoot ("scripts/.runs/impact-fixture-{0}" -f [guid]::NewGuid().ToString('N'))

function Assert-Equal([object]$Actual, [object]$Expected, [string]$Label) { if ($Actual -ne $Expected) { throw "$Label expected '$Expected', got '$Actual'" } }
function Write-File([string]$Path, [string]$Content) { New-Item -ItemType Directory -Path (Split-Path -Parent $Path) -Force | Out-Null; Set-Content -LiteralPath $Path -Value $Content -Encoding utf8 }
function New-Feature([string]$Name, [string[]]$AllowedPaths, [string]$SecondCommand = 'impact-probe-two') {
  $mission = Join-Path $fixtureRoot "MIS-999-$Name"; $milestone = Join-Path $mission 'M-99-impact'; $feature = Join-Path $milestone 'F-10-impact'
  $relative = [IO.Path]::GetRelativePath($repoRoot, $feature).Replace('\','/')
  Write-File (Join-Path $mission 'mission.md') '# Mission'
  Write-File (Join-Path $milestone 'milestone.md') '# Milestone'
  Write-File (Join-Path $milestone 'validation-contract.md') "# Validation`n`nID: M-08-C15`n"
  Write-File (Join-Path $feature 'feature.md') "# F-10`n`n## Brief`n`nRun a registered impact fixture.`n`n## Expected Output`n`n- Two registered commands run."
  Write-File (Join-Path $feature 'spec.md') "# Spec`n`n### F10-AC01`n`nRegistered fixture proof."
  $contract = [ordered]@{ schema_version='1.0'; feature_id='F-10'; required_sources=@(); allowed_paths=$AllowedPaths; forbidden_paths=@('apps/**'); side_effects=@{allowed=@('repository-write');forbidden=@('external-network','provider-write')}; commands=@(@{id='first';command_id='impact-probe-one';lane_id='unit';expected_exit_code=0},@{id='second';command_id=$SecondCommand;lane_id='unit';expected_exit_code=0}); criteria=@(@{id='F10-AC01';milestone_criterion_id='M-08-C15';command_ids=@('first','second')}); stop_conditions=@(@{code='fixture-stop';condition='Stop on invalid input.'}); retry_budget=@{max_correction_attempts=1}; handoff_fields=@('status') }
  Write-File (Join-Path $feature 'plan.md') ("# Plan`n`n## Machine Work Contract`n`n``````json`n" + ($contract | ConvertTo-Json -Depth 20) + "`n``````")
  $relative
}

Import-Module $contextModule -Force
Import-Module $impactModule -Force
$base = (& git -C $repoRoot rev-parse HEAD).Trim()
$changed = @((& git -C $repoRoot diff --name-only $base) + (& git -C $repoRoot ls-files --others --exclude-standard) | ForEach-Object { $_.Trim().Replace('\','/') } | Where-Object { $_ } | Sort-Object -Unique)
$allowed = @($changed + 'AGENTS.md' + 'scripts/.runs/**' | Sort-Object -Unique)
try {
  $feature = New-Feature -Name 'positive' -AllowedPaths $allowed
  $packPath = Join-Path $fixtureRoot 'positive.json'
  $pack = New-HarnessContextPack -FeaturePath $feature -BaseSha $base -AllowedPath $allowed -OutputPath $packPath
  Assert-Equal $pack.Passed $true 'positive pack'
  $result = Invoke-HarnessImpactGate -RepositoryRoot $repoRoot -ContextPath $packPath -RunId 'impact-positive' -RunDirectory (Join-Path $fixtureRoot 'run-positive')
  Assert-Equal $result.Passed $true 'positive gate'
  $outcome = Get-Content -Raw -LiteralPath (Join-Path $fixtureRoot 'run-positive/outcome.json') | ConvertFrom-Json
  Assert-Equal (($outcome.commands.id -join '|')) 'first|second' 'ordered registered commands'
  Assert-Equal $outcome.aggregate_classification 'passed' 'positive aggregate'
  Assert-Equal $outcome.base_sha $base 'outcome base SHA'
  Assert-Equal $outcome.commit_sha $base 'outcome commit SHA'
  if (@($outcome.changed_paths | Where-Object { $_ -match '^(?:[A-Za-z]:|/|\\)' }).Count -ne 0) { throw 'outcome persisted an absolute path' }

  $unknownFeature = New-Feature -Name 'unknown' -AllowedPaths $allowed -SecondCommand 'not-registered'
  $unknownPack = New-HarnessContextPack -FeaturePath $unknownFeature -BaseSha $base -AllowedPath $allowed -OutputPath (Join-Path $fixtureRoot 'unknown.json')
  Assert-Equal $unknownPack.Passed $true 'unknown pack compiles'
  $unknown = Invoke-HarnessImpactGate -RepositoryRoot $repoRoot -ContextPath (Join-Path $fixtureRoot 'unknown.json') -RunId 'impact-unknown' -RunDirectory (Join-Path $fixtureRoot 'run-unknown')
  Assert-Equal $unknown.Passed $false 'unknown command rejects'; Assert-Equal $unknown.ErrorCode 'HIMPACT_COMMAND_UNREGISTERED' 'unknown command reason'

  $stale = Get-Content -Raw -LiteralPath $packPath | ConvertFrom-Json; $stale.base_sha = ('0' * 40); $stale | ConvertTo-Json -Depth 100 | Set-Content -LiteralPath (Join-Path $fixtureRoot 'stale.json') -Encoding utf8
  $staleResult = Invoke-HarnessImpactGate -RepositoryRoot $repoRoot -ContextPath (Join-Path $fixtureRoot 'stale.json') -RunId 'impact-stale' -RunDirectory (Join-Path $fixtureRoot 'run-stale')
  Assert-Equal $staleResult.Passed $false 'stale pack rejects'; Assert-Equal $staleResult.ErrorCode 'CTX_BASE_SHA_MISMATCH' 'stale pack reason'

  $unknownTarget = Get-Content -Raw -LiteralPath $packPath | ConvertFrom-Json; $unknownTarget.commands[0].target_label = 'unknown-target'; $unknownTarget | ConvertTo-Json -Depth 100 | Set-Content -LiteralPath (Join-Path $fixtureRoot 'unknown-target.json') -Encoding utf8
  $unknownTargetResult = Invoke-HarnessImpactGate -RepositoryRoot $repoRoot -ContextPath (Join-Path $fixtureRoot 'unknown-target.json') -RunId 'impact-unknown-target' -RunDirectory (Join-Path $fixtureRoot 'run-unknown-target')
  Assert-Equal $unknownTargetResult.Passed $false 'unknown target rejects'; Assert-Equal $unknownTargetResult.ErrorCode 'CTX_SCHEMA_INVALID' 'unknown target reason'

  $narrowFeature = New-Feature -Name 'narrow' -AllowedPaths @('AGENTS.md','scripts/.runs/**')
  $narrowPack = New-HarnessContextPack -FeaturePath $narrowFeature -BaseSha $base -AllowedPath @('AGENTS.md','scripts/.runs/**') -OutputPath (Join-Path $fixtureRoot 'narrow.json')
  Assert-Equal $narrowPack.Passed $true 'narrow pack compiles'
  $narrow = Invoke-HarnessImpactGate -RepositoryRoot $repoRoot -ContextPath (Join-Path $fixtureRoot 'narrow.json') -RunId 'impact-narrow' -RunDirectory (Join-Path $fixtureRoot 'run-narrow')
  Assert-Equal $narrow.Passed $false 'out-of-scope rejects'; Assert-Equal $narrow.ErrorCode 'HIMPACT_PATH_OUTSIDE_SCOPE' 'out-of-scope reason'

  Write-Output 'PASS harness impact tests'
} finally { Remove-Item -LiteralPath $fixtureRoot -Recurse -Force -ErrorAction SilentlyContinue }
