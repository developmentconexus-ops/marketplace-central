$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$contextModule = Join-Path $repoRoot 'scripts/harness/Context.psm1'
$stateModule = Join-Path $repoRoot 'scripts/harness/State.psm1'
$skillPath = Join-Path $repoRoot '.agents/skills/mpc-goal-harness/SKILL.md'
$fixtureRoot = Join-Path ([IO.Path]::GetTempPath()) ("mpc-orchestration-$([guid]::NewGuid().ToString('N'))")

function Assert-True([bool]$Condition, [string]$Message) {
  if (-not $Condition) { throw $Message }
}

function Test-NativeWorktreeCoordination {
  param([string]$RepositoryRoot, [string]$FixtureRoot)
  $baseSha = (& git -C $RepositoryRoot rev-parse HEAD).Trim()
  $primaryBefore = @(& git -C $RepositoryRoot status --porcelain)
  $worktreePath = Join-Path $FixtureRoot 'native-worktree'
  $added = $false
  try {
    @(& git -C $RepositoryRoot worktree add --detach $worktreePath $baseSha) | Out-Null
    Assert-True ($LASTEXITCODE -eq 0) 'native worktree creation failed'
    $added = $true
    Assert-True ((& git -C $worktreePath rev-parse HEAD).Trim() -eq $baseSha) 'worktree did not start from named SHA'
    $taskBranch = [string](@(& git -C $worktreePath branch --show-current) -join '')
    Assert-True ([string]::IsNullOrWhiteSpace($taskBranch.Trim())) 'worktree is not isolated detached task state'
    Assert-True ((@(& git -C $RepositoryRoot status --porcelain) -join "`n") -eq ($primaryBefore -join "`n")) 'native worktree changed primary checkout'
  }
  finally {
    if ($added) { @(& git -C $RepositoryRoot worktree remove --force $worktreePath) | Out-Null }
  }
}

try {
  Assert-True (Test-Path -LiteralPath $contextModule -PathType Leaf) 'RED: Context.psm1 is missing'
  Assert-True (Test-Path -LiteralPath $stateModule -PathType Leaf) 'RED: State.psm1 is missing'
  Assert-True (Test-Path -LiteralPath $skillPath -PathType Leaf) 'RED: mpc-goal-harness skill is missing'
  Import-Module $contextModule -Force
  Import-Module $stateModule -Force

  $route = Resolve-HarnessKnowledgeRoute -RepositoryRoot $repoRoot -RouteId 'harness-control-plane'
  Assert-True ($route.Passed) 'harness control-plane route did not resolve'
  Assert-True (@($route.Pack.selectors).Count -ge 3) 'route selectors are incomplete'
  Assert-True (@($route.Pack.selectors | Where-Object { $_.path -eq 'scripts/harness/State.psm1' }).Count -eq 1) 'state selector is missing'
  $unknownRoute = Resolve-HarnessKnowledgeRoute -RepositoryRoot $repoRoot -RouteId 'unknown-route'
  Assert-True (-not $unknownRoute.Passed -and $unknownRoute.ErrorCode -eq 'CTX_KNOWLEDGE_ROUTE_UNKNOWN') 'unknown route did not fail closed'

  New-Item -ItemType Directory -Path $fixtureRoot -Force | Out-Null
  Test-NativeWorktreeCoordination -RepositoryRoot $repoRoot -FixtureRoot $fixtureRoot
  $first = New-HarnessLease -StateRoot $fixtureRoot -CheckoutId 'primary' -WriterId 'writer-a' -SharedSeam @('harness-control')
  Assert-True ($first.Passed) 'first writer did not acquire lease'
  $second = New-HarnessLease -StateRoot $fixtureRoot -CheckoutId 'primary' -WriterId 'writer-b' -SharedSeam @('harness-control')
  Assert-True (-not $second.Passed -and $second.ErrorCode -eq 'HSTATE_LEASE_CONFLICT') 'competing writer was not rejected'
  $recovery = Get-HarnessLeaseRecovery -StateRoot $fixtureRoot -CheckoutId 'primary'
  Assert-True ($recovery.Passed -and $recovery.Disposition -eq 'explicit-disposition-required') 'recovery changed an existing lease'

  $checkpoint = New-HarnessCheckpoint -StateRoot $fixtureRoot -CheckpointId 'F05-state-fixture' -FeaturePath 'feature/path' -FeatureState 'implementing' -BaseSha ('1' * 40) -CurrentSha ('2' * 40) -ChangedPath @('scripts/harness/State.psm1') -EvidencePath @('scripts/tests/harness-orchestration.tests.ps1') -Blocker '' -NextAction 'run focused validation' -TaskCorrelationId 'optional-task-id'
  Assert-True ($checkpoint.Passed) 'checkpoint was not persisted'
  $resume = Get-HarnessResume -StateRoot $fixtureRoot -CheckpointId 'F05-state-fixture'
  Assert-True ($resume.Passed -and $resume.NextAction -eq 'run focused validation') 'resume did not reconstruct next action'
  Assert-True ($resume.CheckpointId -eq 'F05-state-fixture') 'checkpoint ID is not canonical'
  Assert-True ($resume.TaskCorrelationId -eq 'optional-task-id') 'task ID correlation was not preserved as optional metadata'

  foreach ($expected in @{
    L0 = 'self-review'; L1 = 'combined-final-review'; L2 = 'fixed-commit-review'; L3 = 'independent-safety-and-quality'
  }.GetEnumerator()) {
    $policy = Get-HarnessRiskPolicy -RiskLevel $expected.Key
    Assert-True ($policy.review -eq $expected.Value) "risk policy mismatch for $($expected.Key)"
    Assert-True (@($policy.command_ids).Count -gt 0) "risk policy has no registered command IDs for $($expected.Key)"
  }

  $skill = Get-Content -Raw -LiteralPath $skillPath
  foreach ($required in @('Portfolio -> Milestone packet', 'Milestone -> Feature packet', 'Feature Plan', 'Feature Execution', 'portfolio_task_id', 'milestone_task_id', 'context_files', 'knowledge_routes', 'constraints', 'operator-observed', 'fresh-session fallback', 'registered command IDs')) {
    Assert-True ($skill -match [regex]::Escape($required)) "skill lacks capability boundary: $required"
  }

  $handoff = Test-HarnessCompactHandoff -Value @{
    status = 'quick_validation_passed'; checkpoint_id = 'F05-state-fixture'; commit = ('3' * 40); changed_paths = @('scripts/harness/State.psm1'); evidence = @('scripts/tests/harness-orchestration.tests.ps1'); review = 'pending orchestrator review'; blockers = ''; next = 'Milestone Orchestrator review'
  }
  Assert-True ($handoff.Passed) 'compact handoff was rejected'
  Write-Output 'PASS harness orchestration tests'
}
finally {
  Remove-Item -LiteralPath $fixtureRoot -Recurse -Force -ErrorAction SilentlyContinue
}
