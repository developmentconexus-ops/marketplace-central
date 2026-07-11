$ErrorActionPreference = 'Stop'
$root = (Resolve-Path (Join-Path $PSScriptRoot '../..')).Path
function Assert-True { param([bool]$Condition,[string]$Message) if (-not $Condition) { throw $Message } }
try {
  Assert-True (Test-Path (Join-Path $root 'scripts/harness.ps1')) 'harness entrypoint missing'
  $harness = Get-Content -Raw (Join-Path $root 'scripts/harness.ps1')
  Assert-True ($harness -match "'cold'") 'cold command missing'
  Assert-True ($harness -match 'CandidateSha') 'candidate SHA argument missing'
  Assert-True (Test-Path (Join-Path $root 'contracts/governance/schemas/harness-outcome.schema.json')) 'outcome schema missing'
  Assert-True ($harness -match 'snapshot' -and $harness -match 'GOMODCACHE') 'isolated snapshot/cache contract missing'
  Assert-True ($harness -match 'New-HarnessChildEnvironment.+cold-provision') 'cold provisioning does not use an isolated child environment'
  Assert-True ($harness -match 'COLD_DOCKER_MISSING') 'missing Docker does not have a stable cold block reason'
  Assert-True ($harness -match 'COLD_IMAGE_IDENTITY_UNRESOLVED') 'unresolved image identity does not block'
  $evidence = Get-Content -Raw (Join-Path $root 'scripts/harness/Evidence.psm1')
  Assert-True ($evidence -match 'Test-Json') 'cold outcome is not schema-validated before persistence'
  Assert-True ($harness -match 'go-mod-download[\s\S]*npm-ci[\s\S]*docker-pull-postgres') 'provisioning order is not frozen Go to npm to Docker'
  Assert-True ($harness -match 'governance-validate[\s\S]*governance-drift[\s\S]*f03-integration') 'required cold validation inventory is incomplete'
  Assert-True ($harness -match 'COLD_CALLER_DIRTY') 'dirty caller does not have stable blocked outcome reason'
  Assert-True ($harness -match '-Dirty \$callerDirty') 'dirty caller state is not recorded in the blocked outcome'
  Import-Module (Join-Path $root 'scripts/harness/Evidence.psm1') -Force
  $canonical = Resolve-HarnessCanonicalCheckoutRoot -SourceRoot $root -ExpectedRoot $root
  Assert-True ([IO.Path]::IsPathFullyQualified($canonical)) 'canonical source root is not absolute'
  $mismatch = $false
  try { Resolve-HarnessCanonicalCheckoutRoot -SourceRoot $root -ExpectedRoot (Split-Path -Parent $root) | Out-Null } catch { $mismatch = $_.Exception.Message -eq 'COLD_SOURCE_ROOT_MISMATCH' }
  Assert-True $mismatch 'ancestor/divergent source root did not fail closed'
  $wildcard = $false
  try { New-HarnessScopedGitCloneArguments -CanonicalSourceRoot ($root + '\*') -SnapshotPath (Join-Path $root 'scripts/.runs/snapshot') | Out-Null } catch { $wildcard = $_.Exception.Message -eq 'COLD_GIT_TRUST_SCOPE_INVALID' }
  Assert-True $wildcard 'wildcard safe.directory source was accepted'
  $suffix = $false
  try { New-HarnessScopedGitCloneArguments -CanonicalSourceRoot ($root + '\.git') -SnapshotPath (Join-Path $root 'scripts/.runs/snapshot') | Out-Null } catch { $suffix = $_.Exception.Message -eq 'COLD_GIT_TRUST_SCOPE_INVALID' }
  Assert-True $suffix 'suffixed .git trust source was accepted'
  $multiple = $false
  try { New-HarnessScopedGitCloneArguments -CanonicalSourceRoot ($root + ';safe.directory=' + $root) -SnapshotPath (Join-Path $root 'scripts/.runs/snapshot') | Out-Null } catch { $multiple = $_.Exception.Message -eq 'COLD_GIT_TRUST_SCOPE_INVALID' }
  Assert-True $multiple 'multiple safe.directory values were accepted'
  $cloneArgs = @(New-HarnessScopedGitCloneArguments -CanonicalSourceRoot $canonical -SnapshotPath (Join-Path $root 'scripts/.runs/snapshot'))
  Assert-True ((@($cloneArgs | Where-Object { $_ -eq '-c' }).Count) -eq 1) 'scoped clone must have exactly one config override'
  Assert-True ($cloneArgs -contains "safe.directory=$canonical\.git") 'clone does not scope safe.directory to canonical gitdir'
  Assert-True (-not ($cloneArgs -contains '--no-local')) 'unsafe --no-local fallback was introduced'
  Assert-True (@($cloneArgs | Where-Object { [string]$_ -match '(?i)safe\.directory=.*(?:\*|\?|;)' }).Count -eq 0) 'wildcard/multiple safe.directory value was projected'
  $scopeBefore = @(& git config --global --get-all safe.directory 2>$null; & git config --system --get-all safe.directory 2>$null; & git config --local --get-all safe.directory 2>$null)
  $scopeAfter = @(& git config --global --get-all safe.directory 2>$null; & git config --system --get-all safe.directory 2>$null; & git config --local --get-all safe.directory 2>$null)
  Assert-True ((($scopeBefore -join "`n") -ceq ($scopeAfter -join "`n"))) 'Git global/system/local config changed during scoped trust preparation'
  Import-Module (Join-Path $root 'scripts/harness/Execution.psm1') -Force
  $node = (Get-Command node -CommandType Application | Select-Object -First 1).Source
  $probe = Join-Path $root 'scripts/tests/fixtures/harness/child-probe.mjs'
  $spaceGitdir = 'C:\fixture root with spaces\.git'
  $launcherArgs = @($probe, 'inspect', '-c', "safe.directory=$spaceGitdir", 'clone', '--local', 'C:\fixture root with spaces', 'C:\run snapshot')
  $launcherEnv = @{}
  foreach ($name in @('SystemRoot', 'WINDIR', 'ComSpec', 'PATH', 'PATHEXT', 'TEMP', 'TMP')) {
    $value = [Environment]::GetEnvironmentVariable($name, 'Process')
    if (-not [string]::IsNullOrWhiteSpace($value)) { $launcherEnv[$name] = $value }
  }
  $launcherResult = Invoke-HarnessProcess -Request (New-HarnessProcessRequest -FilePath $node -ArgumentList $launcherArgs -WorkingDirectory $root -Environment $launcherEnv -TimeoutSeconds 10)
  Assert-True ($launcherResult.ExitCode -eq 0) 'launcher argv probe failed'
  $launcherObserved = $launcherResult.Stdout | ConvertFrom-Json
  $expectedLauncherArgs = @($launcherArgs | Select-Object -Skip 2)
  Assert-True (($launcherObserved.args -join "`n") -eq ($expectedLauncherArgs -join "`n")) 'launcher changed scoped clone argv boundaries'
  $rootOnlyTarget = Join-Path $env:TEMP ('mpc-root-only-' + [guid]::NewGuid().ToString('N'))
  try {
    & git -c "safe.directory=$canonical" clone --quiet --no-hardlinks --local $canonical $rootOnlyTarget 2>$null
    Assert-True ($LASTEXITCODE -eq 128) 'root-only safe.directory unexpectedly bypassed dubious ownership'
  } finally { Remove-Item -LiteralPath $rootOnlyTarget -Recurse -Force -ErrorAction SilentlyContinue }
  $fixtureRoot = Join-Path $env:TEMP ('mpc-scoped-git-' + [guid]::NewGuid().ToString('N'))
  New-Item -ItemType Directory -Path $fixtureRoot -Force | Out-Null
  try {
    & git init --quiet $fixtureRoot 2>$null
    Assert-True ($LASTEXITCODE -eq 0) 'fixture repository initialization failed'
    $fixtureCanonical = Resolve-HarnessCanonicalCheckoutRoot -SourceRoot $fixtureRoot -ExpectedRoot $fixtureRoot
    $fixtureSnapshot = Join-Path $fixtureRoot 'snapshot'
    $fixtureArgs = @(New-HarnessScopedGitCloneArguments -CanonicalSourceRoot $fixtureCanonical -SnapshotPath $fixtureSnapshot)
  Assert-True ($fixtureArgs -contains "safe.directory=$fixtureCanonical\.git") 'fixture trust did not target direct child gitdir'
    & git @fixtureArgs 2>$null
    Assert-True ($LASTEXITCODE -eq 0) 'exact scoped trust clone did not permit validated fixture source'
  } finally { Remove-Item -LiteralPath $fixtureRoot -Recurse -Force -ErrorAction SilentlyContinue }
  $gitfileRoot = Join-Path $env:TEMP ('mpc-gitfile-' + [guid]::NewGuid().ToString('N'))
  New-Item -ItemType Directory -Path $gitfileRoot -Force | Out-Null
  try {
    Set-Content -LiteralPath (Join-Path $gitfileRoot '.git') -Value 'gitdir: C:\external' -Encoding ascii
    $gitfileRejected = $false
    try { New-HarnessScopedGitCloneArguments -CanonicalSourceRoot $gitfileRoot -SnapshotPath (Join-Path $gitfileRoot 'snapshot') | Out-Null } catch { $gitfileRejected = $_.Exception.Message -eq 'GITDIR_UNTRUSTED' }
    Assert-True $gitfileRejected 'external gitfile trust was accepted'
  } finally { Remove-Item -LiteralPath $gitfileRoot -Recurse -Force -ErrorAction SilentlyContinue }
  $projectionCommands = @([ordered]@{ id='snapshot'; command='git clone detached snapshot'; stage='preflight'; target_label='fake'; evidence_class='contract'; duration_ms=1; exit_code=0; reason='passed'; artifact_paths=@() })
  $projection = Get-HarnessOutcomeProjection (New-HarnessOutcome -RunId 'scope-test' -CandidateSha ('a' * 40) -Branch 'main' -Dirty $false -Tools @{} -Commands $projectionCommands -AggregateClassification 'passed')
  Assert-True (($projection.commands[0].command -notmatch '(?i)([A-Z]:[\\/]|/|\\\\)') -and ($projection.commands[0].command -notmatch '(?i)safe\.directory')) 'absolute trust source leaked into projected command'
  Write-Output 'PASS cold gate snapshot contract'
  exit 0
} catch { Write-Output "FAIL cold gate snapshot contract: $($_.Exception.Message)"; exit 1 }
