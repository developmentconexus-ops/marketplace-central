$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
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
  Assert-True (Test-Path -LiteralPath (Join-Path $repoRoot 'docs/HARNESS.md') -PathType Leaf) 'RED: binding harness doctrine (docs/HARNESS.md) is missing'
  $workerSkillPath = @(Get-ChildItem -Path (Join-Path $env:USERPROFILE '.claude/plugins/cache/mnfs-harness/harness/*/skills/harness-worker/SKILL.md') -ErrorAction SilentlyContinue | Sort-Object FullName -Descending | Select-Object -First 1).FullName
  Assert-True (-not [string]::IsNullOrEmpty($workerSkillPath)) 'RED: harness-worker skill missing — install plugin: claude plugin install harness@mnfs-harness'

  New-Item -ItemType Directory -Path $fixtureRoot -Force | Out-Null
  Test-NativeWorktreeCoordination -RepositoryRoot $repoRoot -FixtureRoot $fixtureRoot

  $workerSkill = Get-Content -Raw -LiteralPath $workerSkillPath
  foreach ($required in @('docs/HARNESS.md', 'mpc-goal-harness', 'One writer per shared seam', 'ADR-17', 'GOCACHE=.gocache')) {
    Assert-True ($workerSkill -match [regex]::Escape($required)) "harness-worker skill lacks pinned rule: $required"
  }
  $doctrine = Get-Content -Raw -LiteralPath (Join-Path $repoRoot 'docs/HARNESS.md')
  foreach ($required in @('anti-slop', 'Verification ladder', 'collision matrix', 'SPLIT-REQUEST')) {
    Assert-True ($doctrine -match $required) "harness doctrine lacks section marker: $required"
  }
  Write-Output 'PASS harness orchestration tests'
}
finally {
  Remove-Item -LiteralPath $fixtureRoot -Recurse -Force -ErrorAction SilentlyContinue
}
