[CmdletBinding()]
param(
  [ValidateSet('unit', 'integration', 'live', 'browser', 'provider-write', 'governance-validate', 'governance-drift', 'governance', 'context-compile', 'context-validate', 'cold')]
  [string]$Command = 'unit',
  [switch]$PreflightOnly,
  [string]$EnvFile,
  [string]$Target = 'oracle',
  [string]$DatabaseUrl,
  [string]$Provider,
  [string]$Actor,
  [string]$IdempotencyKey,
  [string]$BaseSha,
  [string]$FeaturePath,
  [string[]]$AllowedPath,
  [string]$ContextPath,
  [switch]$RequireCurrentBase,
  [switch]$Execute
  ,[string]$CandidateSha
)

$ErrorActionPreference = 'Stop'
$repoRoot = Split-Path -Parent $PSScriptRoot
$runRoot = Join-Path $PSScriptRoot '.runs'
$runId = [guid]::NewGuid().ToString('N')
$runDir = Join-Path $runRoot $runId

Import-Module (Join-Path $PSScriptRoot 'harness/Environment.psm1') -Force
Import-Module (Join-Path $PSScriptRoot 'harness/Execution.psm1') -Force
Import-Module (Join-Path $PSScriptRoot 'harness/Postgres.psm1') -Force
Import-Module (Join-Path $PSScriptRoot 'harness/Evidence.psm1') -Force

function Resolve-HarnessApplication {
  param([Parameter(Mandatory)][string]$Name)

  $application = Get-Command $Name -CommandType Application -ErrorAction SilentlyContinue | Select-Object -First 1
  if ($null -eq $application -or [string]::IsNullOrWhiteSpace([string]$application.Source)) {
    throw "HEXEC_FILE_NOT_FOUND tool=$Name"
  }
  return [IO.Path]::GetFullPath([string]$application.Source)
}

function Get-HarnessNpmInvocation {
  $npmPath = Resolve-HarnessApplication -Name 'npm'
  if ([Environment]::OSVersion.Platform -ne [PlatformID]::Win32NT) {
    return [pscustomobject]@{ FilePath = $npmPath; ArgumentPrefix = @() }
  }

  $nodePath = Resolve-HarnessApplication -Name 'node'
  $npmCliPath = Join-Path (Split-Path -Parent $npmPath) 'node_modules/npm/bin/npm-cli.js'
  if (-not (Test-Path -LiteralPath $npmCliPath -PathType Leaf)) { throw 'HEXEC_FILE_NOT_FOUND tool=npm' }
  return [pscustomobject]@{ FilePath = $nodePath; ArgumentPrefix = @([IO.Path]::GetFullPath($npmCliPath)) }
}

function Write-HarnessProcessResult {
  param([Parameter(Mandatory)][object]$Result)

  if (-not [string]::IsNullOrEmpty([string]$Result.Stdout)) { Write-Output ([string]$Result.Stdout).TrimEnd("`r", "`n") }
  if (-not [string]::IsNullOrEmpty([string]$Result.Stderr)) { Write-Output ([string]$Result.Stderr).TrimEnd("`r", "`n") }
}

function Write-Summary {
  param([string]$TargetType, [string]$Status)
  New-Item -ItemType Directory -Path $runDir -Force | Out-Null
  $summary = @("target=$TargetType", "status=$Status", "run_id=$runId")
  $summary | Set-Content -LiteralPath (Join-Path $runDir 'summary.txt') -Encoding utf8
  $summary | ForEach-Object { Write-Output $_ }
  Write-Output "run_dir=scripts/.runs/$runId"
}

function Invoke-Unit {
  $goPath = Resolve-HarnessApplication -Name 'go'
  $npm = Get-HarnessNpmInvocation
  $childEnvironment = New-HarnessChildEnvironment -RepositoryRoot $repoRoot -LaneId 'unit'
  $goWorkingDirectory = [IO.Path]::GetFullPath((Join-Path $repoRoot 'apps/server_core'))

  Write-Output 'target=fake'
  Write-Output 'env=ignored'
  Write-Output 'postgres=disabled'
  Write-Output 'oracle=disabled'
  Write-Output 'provider_network=disabled'
  Write-Output 'migrations=disabled'
  if ($PreflightOnly) { Write-Summary -TargetType 'fake' -Status 'ready'; return }

  $goRequest = New-HarnessProcessRequest -FilePath $goPath -ArgumentList @('test', './tests/unit/...', '-count=1') -WorkingDirectory $goWorkingDirectory -Environment $childEnvironment -TimeoutSeconds 1200
  $go = Invoke-HarnessProcess -Request $goRequest
  Write-HarnessProcessResult $go
  if ($go.ExitCode -ne 0) { throw "unit Go tests failed reason=$($go.ReasonCode) exit_code=$($go.ExitCode)" }

  $webArguments = @($npm.ArgumentPrefix) + @('run', 'test', '--workspace', '@marketplace-central/web', '--', '--run')
  $webRequest = New-HarnessProcessRequest -FilePath $npm.FilePath -ArgumentList $webArguments -WorkingDirectory $repoRoot -Environment $childEnvironment -TimeoutSeconds 1200
  $web = Invoke-HarnessProcess -Request $webRequest
  Write-HarnessProcessResult $web
  if ($web.ExitCode -ne 0) { throw "unit web tests failed reason=$($web.ReasonCode) exit_code=$($web.ExitCode)" }
  Write-Summary -TargetType 'fake' -Status 'passed'
}

function Invoke-Integration {
  if (-not [string]::IsNullOrWhiteSpace($DatabaseUrl)) { throw 'HPG_EXTERNAL_TARGET_FORBIDDEN' }
  Write-Output 'target=ephemeral-postgres'
  Write-Output 'key=MPC_TEST_DATABASE_URL'
  Write-Output 'migrations=embedded'
  if ($PreflightOnly) { Write-Output 'status=ready'; return }

  $dockerPath = Resolve-HarnessPostgresDockerApplication
  $goPath = Resolve-HarnessApplication -Name 'go'
  $childEnvironment = New-HarnessChildEnvironment -RepositoryRoot $repoRoot -LaneId 'integration'
  $passwordBytes = [byte[]]::new(24)
  [Security.Cryptography.RandomNumberGenerator]::Fill($passwordBytes)
  $password = [Convert]::ToHexString($passwordBytes).ToLowerInvariant()
  $spec = New-HarnessPostgresRunSpec -RepositoryRoot $repoRoot -RunId $runId -Password $password -DockerFilePath $dockerPath
  $result = Invoke-HarnessPostgresLifecycle -RunSpec $spec -BaseEnvironment $childEnvironment -GoFilePath $goPath -TimeoutSeconds 1200
  Write-Output "migrations_first=$($result.MigrationsAppliedFirst)"
  Write-Output "migrations_second=$($result.MigrationsAppliedSecond)"
  Write-Output "resource_count=$(@($result.ResourceInventory).Count)"
  Write-Output "port=$($result.HostPort)"
  if ($result.ExitCode -ne 0) {
    foreach ($token in @($result.FailureDiagnosticTokens)) { Write-Output "child_diagnostic=$token" }
    $reasons = @($result.PrimaryReasonCode) + @($result.CleanupReasonCodes) | Where-Object { -not [string]::IsNullOrWhiteSpace($_) }
    throw "postgres lifecycle failed reasons=$($reasons -join ',') exit_code=$($result.ExitCode)"
  }
  Write-Summary -TargetType 'ephemeral-postgres' -Status 'passed'
}

function Invoke-Live {
  if ($EnvFile) {
    $path = if ([IO.Path]::IsPathFullyQualified($EnvFile)) { [IO.Path]::GetFullPath($EnvFile) } else { [IO.Path]::GetFullPath((Join-Path $repoRoot $EnvFile)) }
  } else {
    $path = Join-Path $repoRoot '.env'
  }
  if ($Target -notin @('oracle', 'provider')) { throw 'live target must be oracle or provider' }
  $laneId = if ($Target -eq 'oracle') { 'live-oracle' } else { 'live-provider-read' }
  $childEnvironment = New-HarnessChildEnvironment -RepositoryRoot $repoRoot -LaneId $laneId -EnvFile $path
  $required = if ($Target -eq 'oracle') { @('MPC_ORACLE_USERNAME', 'MPC_ORACLE_PASSWORD', 'MPC_ORACLE_CONNECT_STRING') } else { @('MPC_PROVIDER_MERCADOLIVRE_CLIENT_ID', 'MPC_PROVIDER_MERCADOLIVRE_CLIENT_SECRET') }
  $missing = @($required | Where-Object { -not $childEnvironment.ContainsKey($_) -or [string]::IsNullOrWhiteSpace($childEnvironment[$_]) })
  Write-Output "target=live-$Target"
  foreach ($key in @($childEnvironment.Keys | Sort-Object)) {
    if ($key -notin @('SystemRoot', 'WINDIR', 'ComSpec', 'PATH', 'PATHEXT', 'TEMP', 'TMP', 'GOCACHE')) { Write-Output "key=$key" }
  }
  if ($missing.Count -gt 0) { throw "live preflight missing_keys=$($missing -join ',')" }
  Write-Output 'provider_write=disabled'
  Write-Summary -TargetType "live-$Target" -Status 'ready'
}

function Invoke-Browser {
  Write-Output 'target=browser'
  Write-Output 'provider_network=disabled'
  if (-not $PreflightOnly) { throw 'browser runner is not configured; use browser automation as an explicit lane' }
  Write-Summary -TargetType 'browser' -Status 'ready'
}

function Invoke-ProviderWrite {
  Write-Output 'target=live-provider'
  if ([string]::IsNullOrWhiteSpace($Provider)) { throw 'provider is required' }
  $missing = @()
  if ([string]::IsNullOrWhiteSpace($Actor)) { $missing += 'actor' }
  if ([string]::IsNullOrWhiteSpace($IdempotencyKey)) { $missing += 'idempotency_key' }
  if ($missing.Count -gt 0) { throw "provider write missing=$($missing -join ','); rejected before network" }
  if (-not $Execute) { throw 'explicit -Execute is required before network' }
  throw 'provider write adapter is intentionally outside F-02; no network was invoked'
}

function Write-GovernanceResult {
  param([object]$Result)
  if ($Result.Passed) {
    Write-Output 'status=passed'
  } else {
    Write-Output 'status=failed'
    foreach ($violation in @($Result.Violations)) {
      Write-Output "error_code=$($violation.ErrorCode)"
      Write-Output "id=$($violation.Id)"
      if (-not [string]::IsNullOrWhiteSpace([string]$violation.Path)) { Write-Output "path=$($violation.Path)" }
    }
  }
  foreach ($exception in @($Result.BaselineExceptions)) {
    Write-Output "baseline_exception=$($exception.Id)"
  }
  Write-Output 'artifact_path=contracts/governance'
}

function Invoke-Governance {
  param([ValidateSet('validate', 'drift', 'all')][string]$Mode)
  Import-Module (Join-Path $PSScriptRoot 'harness/Policy.psm1') -Force
  if ($Mode -eq 'validate') {
    $result = Test-GovernanceContracts -RepositoryRoot $repoRoot
  } else {
    if ([string]::IsNullOrWhiteSpace($BaseSha) -or $BaseSha -notmatch '^[0-9a-f]{40}$') {
      Write-Output 'status=failed'
      Write-Output 'error_code=GOV_SEMANTIC_DRIFT'
      Write-Output 'id=base-sha-invalid'
      Write-Output 'artifact_path=contracts/governance'
      exit 1
    }
    if ($Mode -eq 'all') {
      $contracts = Test-GovernanceContracts -RepositoryRoot $repoRoot
      if (-not $contracts.Passed) { Write-GovernanceResult $contracts; exit 1 }
    }
    $result = Test-GovernanceDrift -RepositoryRoot $repoRoot -BaseSha $BaseSha
  }
  Write-GovernanceResult $result
  if (-not $result.Passed) { exit 1 }
}

function Write-ContextResult {
  param([object]$Result, [string]$ArtifactPath)
  Write-Output "status=$($Result.Status)"
  if (-not $Result.Passed) {
    Write-Output "error_code=$($Result.ErrorCode)"
    if (-not [string]::IsNullOrWhiteSpace([string]$Result.Id)) { Write-Output "id=$($Result.Id)" }
    if (-not [string]::IsNullOrWhiteSpace([string]$Result.Path)) { Write-Output "path=$($Result.Path)" }
  }
  Write-Output "artifact_path=$ArtifactPath"
}

function Invoke-Context {
  param([ValidateSet('compile', 'validate')][string]$Mode)
  Import-Module (Join-Path $PSScriptRoot 'harness/Context.psm1') -Force
  if ($Mode -eq 'compile') {
    $artifactPath = "scripts/.runs/$runId/context-pack.json"
    if ([string]::IsNullOrWhiteSpace($FeaturePath) -or @($AllowedPath).Count -eq 0) {
      $result = [pscustomobject]@{ Passed = $false; Status = 'failed'; ErrorCode = 'CTX_FEATURE_INVALID'; Id = 'compile-input'; Path = '' }
    } else {
      $result = New-HarnessContextPack -FeaturePath $FeaturePath -BaseSha $BaseSha -AllowedPath $AllowedPath -OutputPath (Join-Path $repoRoot $artifactPath)
    }
  } else {
    $artifactPath = 'context-pack.json'
    if ([string]::IsNullOrWhiteSpace($ContextPath)) {
      $result = [pscustomobject]@{ Passed = $false; Status = 'failed'; ErrorCode = 'CTX_SOURCE_MISSING'; Id = 'context-pack'; Path = '' }
    } else {
      $resolvedContextPath = if ([IO.Path]::IsPathRooted($ContextPath)) { $ContextPath } else { Join-Path $repoRoot $ContextPath }
      $artifactPath = if ([IO.Path]::GetFullPath($resolvedContextPath).StartsWith([IO.Path]::GetFullPath($repoRoot), [StringComparison]::OrdinalIgnoreCase)) { [IO.Path]::GetRelativePath($repoRoot, $resolvedContextPath).Replace('\', '/') } else { 'context-pack.json' }
      $result = Test-HarnessContextPack -Path $resolvedContextPath -RepositoryRoot $repoRoot -RequireCurrentBase:$RequireCurrentBase
    }
  }
  Write-ContextResult $result $artifactPath
  if (-not $result.Passed) { exit 1 }
}

function Invoke-Cold {
  New-Item -ItemType Directory -Path $runDir -Force | Out-Null
  $records = [Collections.Generic.List[object]]::new(); $imageIdentity = ''; $blocked = $false; $callerDirty = $false
  $head = ''; $safeCandidate = if ($CandidateSha -match '^[0-9a-f]{40}$') { $CandidateSha } else { '0000000000000000000000000000000000000000' }
  $branchName = ((& git -C $repoRoot branch --show-current).Trim()); if ([string]::IsNullOrWhiteSpace($branchName)) { $branchName = 'detached' }
  function Add-ColdRecord([string]$Id,[string]$Command,[string]$Stage,[string]$Target,[string]$Class,[int]$Exit,[string]$Reason,[long]$Duration = 0) {
    $records.Add([ordered]@{id=$Id;command=$Command;stage=$Stage;target_label=$Target;evidence_class=$Class;duration_ms=$Duration;exit_code=$Exit;reason=$Reason;artifact_paths=@()})
  }
  function Invoke-ColdProcess([string]$Id,[string]$CommandText,[string]$Stage,[string]$Target,[string]$Class,[string]$File,[string[]]$Arguments,[string]$Working,[System.Collections.IDictionary]$Environment) {
    $timer=[Diagnostics.Stopwatch]::StartNew()
    try {
      $result=Invoke-HarnessProcess -Request (New-HarnessProcessRequest -FilePath $File -ArgumentList $Arguments -WorkingDirectory $Working -Environment $Environment -TimeoutSeconds 1200 -RedactionCandidates @($repoRoot))
      $timer.Stop()
      Add-ColdRecord $Id $CommandText $Stage $Target $Class $result.ExitCode $(if($result.ExitCode -eq 0){'passed'}elseif($result.TimedOut){'HEXEC_TIMEOUT'}else{$result.ReasonCode}) $timer.ElapsedMilliseconds
      return $result
    } catch {
      $timer.Stop()
      Add-ColdRecord $Id $CommandText $Stage $Target $Class 1 (Get-HarnessProcessFailureReason -CommandId $Id -Exception $_.Exception) $timer.ElapsedMilliseconds
      throw
    }
  }
  try {
    if ([string]::IsNullOrWhiteSpace($CandidateSha) -or $CandidateSha -notmatch '^[0-9a-f]{40}$') { Add-ColdRecord 'preflight' 'git clean candidate' 'preflight' 'fake' 'contract' 1 'COLD_CANDIDATE_SHA_INVALID'; $blocked=$true; return }
    try {
      $expectedCheckoutRoot = (Get-Item -LiteralPath (Join-Path $PSScriptRoot '..') -Force).FullName
      $canonicalSource = Resolve-HarnessCanonicalCheckoutRoot -SourceRoot $repoRoot -ExpectedRoot $expectedCheckoutRoot
    } catch {
      $sourceReason = [string]$_.Exception.Message
      if ($sourceReason -notmatch '^COLD_SOURCE_(?:ROOT_INVALID|ROOT_MISSING|REPARSE_UNSAFE|ROOT_MISMATCH)$') { $sourceReason = 'COLD_SOURCE_ROOT_UNSAFE' }
      Add-ColdRecord 'preflight' 'git canonical checkout source' 'preflight' 'fake' 'contract' 1 $sourceReason; $blocked=$true; return
    }
    $headExit = 0; $head = (& git -C $canonicalSource rev-parse HEAD 2>$null).Trim(); $headExit = $LASTEXITCODE
    if ($headExit -ne 0 -or $head -cne $CandidateSha) { Add-ColdRecord 'preflight' 'git clean candidate' 'preflight' 'fake' 'contract' 1 'COLD_CANDIDATE_SHA_MISMATCH'; $blocked=$true; return }
    $status = @(& git -C $repoRoot status --porcelain 2>$null); if ($status.Count -gt 0) { $callerDirty = $true; Add-ColdRecord 'preflight' 'git clean candidate' 'preflight' 'fake' 'contract' 1 'COLD_CALLER_DIRTY'; $blocked=$true; return }
    Add-ColdRecord 'preflight' 'git clean candidate' 'preflight' 'fake' 'contract' 0 'passed'
    $snapshot = Join-Path $runDir 'snapshot'
    try {
      $gitPath = Resolve-HarnessApplication 'git'
      $cloneArguments = New-HarnessScopedGitCloneArguments -CanonicalSourceRoot $canonicalSource -SnapshotPath $snapshot
      $cloneResult = Invoke-HarnessProcess -Request (New-HarnessProcessRequest -FilePath $gitPath -ArgumentList $cloneArguments -WorkingDirectory $canonicalSource -Environment @{} -TimeoutSeconds 120 -RedactionCandidates @($repoRoot))
    } catch {
      $snapshotReason = [string]$_.Exception.Message
      if ($snapshotReason -notmatch '^(?:GITDIR_UNTRUSTED|COLD_GIT_TRUST_SCOPE_INVALID|COLD_SNAPSHOT_PATH_INVALID)$') { $snapshotReason = 'COLD_SNAPSHOT_FAILED' }
      Add-ColdRecord 'snapshot' 'git clone detached snapshot' 'preflight' 'fake' 'contract' 1 $snapshotReason; $blocked=$true; return
    }
    if ($cloneResult.ExitCode -ne 0) { Add-ColdRecord 'snapshot' 'git clone detached snapshot' 'preflight' 'fake' 'contract' 1 'COLD_SNAPSHOT_FAILED'; $blocked=$true; return }
    & git -C $snapshot checkout --quiet --detach $CandidateSha 2>$null; $snapshotHead=(& git -C $snapshot rev-parse HEAD 2>$null).Trim(); $snapshotStatus=@(& git -C $snapshot status --porcelain 2>$null)
    if ($LASTEXITCODE -ne 0 -or $snapshotHead -cne $CandidateSha -or $snapshotStatus.Count -gt 0) { Add-ColdRecord 'snapshot' 'git clone detached snapshot' 'preflight' 'fake' 'contract' 1 'COLD_SNAPSHOT_INVARIANT_FAILED'; $blocked=$true; return }
    Add-ColdRecord 'snapshot' 'git clone detached snapshot' 'preflight' 'fake' 'contract' 0 'passed'
    $runCache = Join-Path $runDir 'cache'; New-Item -ItemType Directory -Force -Path $runCache | Out-Null
    $coldEnvironment = New-HarnessChildEnvironment -RepositoryRoot $snapshot -LaneId 'cold-provision'
    $coldEnvironment['GOCACHE']=Join-Path $runCache 'go-build'; $coldEnvironment['GOMODCACHE']=Join-Path $runCache 'go-mod'; $coldEnvironment['npm_config_cache']=Join-Path $runCache 'npm'; $coldEnvironment['npm_config_userconfig']=Join-Path $runCache 'npmrc'; $coldEnvironment['npm_config_globalconfig']=Join-Path $runCache 'npm-globalrc'
    New-Item -ItemType Directory -Force -Path $coldEnvironment['GOCACHE'],$coldEnvironment['GOMODCACHE'],$coldEnvironment['npm_config_cache'] | Out-Null
    $goPath=Resolve-HarnessApplication 'go'; $npm=Get-HarnessNpmInvocation
    $goResult=Invoke-ColdProcess 'go-mod-download' 'go mod download' 'provisioning' 'external-dependency-registry' 'provisioning' $goPath @('mod','download') (Join-Path $snapshot 'apps/server_core') $coldEnvironment
    if($goResult.ExitCode -ne 0){$blocked=$true; return}
    $npmResult=Invoke-ColdProcess 'npm-ci' 'npm ci --ignore-scripts' 'provisioning' 'external-dependency-registry' 'provisioning' $npm.FilePath (@($npm.ArgumentPrefix)+@('ci','--ignore-scripts')) $snapshot $coldEnvironment
    if($npmResult.ExitCode -ne 0){$blocked=$true; return}
    try {$dockerPath=Resolve-HarnessApplication 'docker'} catch { Add-ColdRecord 'docker-pull-postgres' 'docker pull postgres:16-bookworm' 'provisioning' 'external-dependency-registry' 'provisioning' 1 'COLD_DOCKER_MISSING'; $blocked=$true; return }
    $dockerResult=Invoke-ColdProcess 'docker-pull-postgres' 'docker pull postgres:16-bookworm' 'provisioning' 'external-dependency-registry' 'provisioning' $dockerPath @('pull','postgres:16-bookworm') $snapshot $coldEnvironment
    if($dockerResult.ExitCode -ne 0){$blocked=$true; return}
    $imageResult=Invoke-HarnessProcess -Request (New-HarnessProcessRequest -FilePath $dockerPath -ArgumentList @('image','inspect','postgres:16-bookworm','--format','{{.Id}}') -WorkingDirectory $snapshot -Environment $coldEnvironment -TimeoutSeconds 120 -RedactionCandidates @($repoRoot)); $imageIdentity=$imageResult.Stdout.Trim()
    if($imageResult.ExitCode -ne 0 -or $imageIdentity -notmatch '^sha256:[0-9a-f]{64}$'){ Add-ColdRecord 'docker-image-identity' 'docker image inspect postgres:16-bookworm' 'provisioning' 'external-dependency-registry' 'provisioning' 1 'COLD_IMAGE_IDENTITY_UNRESOLVED'; $blocked=$true; return }
    Add-ColdRecord 'docker-image-identity' 'docker image inspect postgres:16-bookworm' 'provisioning' 'external-dependency-registry' 'provisioning' 0 'passed'
    $pwshPath=Resolve-HarnessApplication 'pwsh'
    foreach($entry in @(
      @{id='governance-validate'; args=@('-Command','governance-validate'); text='harness governance validate'; target='fake'; class='contract'},
      @{id='governance-drift'; args=@('-Command','governance-drift','-BaseSha',$CandidateSha); text='harness governance drift'; target='fake'; class='contract'},
      @{id='unit'; args=@('-Command','unit'); text='harness unit'; target='fake'; class='contract'}
    )) { [void](Invoke-ColdProcess $entry.id $entry.text 'validation' $entry.target $entry.class $pwshPath (@('-NoProfile','-ExecutionPolicy','Bypass','-File',(Join-Path $snapshot 'scripts/harness.ps1'))+$entry.args) $snapshot $coldEnvironment) }
    [void](Invoke-ColdProcess 'root-test-alias' 'npm run test:run' 'validation' 'fake' 'contract' $npm.FilePath (@($npm.ArgumentPrefix)+@('run','test:run')) $snapshot $coldEnvironment)
    [void](Invoke-ColdProcess 'root-build-alias' 'npm run build' 'validation' 'fake' 'contract' $npm.FilePath (@($npm.ArgumentPrefix)+@('run','build')) $snapshot $coldEnvironment)
    $contextTimer=[Diagnostics.Stopwatch]::StartNew(); $featurePath='.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/F-04-cold-gate-evidence'; $planPath=Join-Path $snapshot "$featurePath/plan.md"
    $planText=Get-Content -Raw -LiteralPath $planPath; $contractJson=[regex]::Match($planText,'(?ms)^```json\s*\r?\n(?<json>.*?)\r?\n```').Groups['json'].Value; $contract=$contractJson|ConvertFrom-Json -Depth 100
    Import-Module (Join-Path $snapshot 'scripts/harness/Context.psm1') -Force
    $contextPath=Join-Path $runDir 'context-pack.json'; $compiled=New-HarnessContextPack -FeaturePath $featurePath -BaseSha $CandidateSha -AllowedPath @($contract.allowed_paths) -OutputPath $contextPath
    $current=if($compiled.Passed){Test-HarnessContextPack -Path $contextPath -RepositoryRoot $snapshot -RequireCurrentBase}else{$compiled}; $contextTimer.Stop()
    Add-ColdRecord 'context-current' 'compile and current-validate F-04 context pack' 'validation' 'cold-gate' 'contract' $(if($current.Passed){0}else{1}) $(if($current.Passed){'passed'}else{[string]$current.ErrorCode}) $contextTimer.ElapsedMilliseconds
    # F-03 uses --pull=never inside its accepted ephemeral PostgreSQL lifecycle.
    [void](Invoke-ColdProcess 'f03-integration' 'harness integration --pull=never' 'validation' 'ephemeral-postgres' 'integration' $pwshPath @('-NoProfile','-ExecutionPolicy','Bypass','-File',(Join-Path $snapshot 'scripts/harness.ps1'),'-Command','integration') $snapshot $coldEnvironment)
  } catch { Add-ColdRecord 'cold-gate' 'cold gate orchestration' 'cleanup' 'cold-gate' 'contract' 1 'COLD_UNEXPECTED_BLOCK'; $blocked=$true }
  finally {
    if($records.Count -eq 0){Add-ColdRecord 'preflight' 'git clean candidate' 'preflight' 'fake' 'contract' 1 'COLD_PREFLIGHT_UNRECORDED'}
    $aggregate=if($blocked){'blocked'}elseif(@($records|Where-Object exit_code -ne 0).Count -gt 0){'failed'}else{'passed'}
    $tools=[ordered]@{}; try{$tools.go=((& go version) -join ' ');$tools.npm=((& npm --version) -join ' ');$tools.docker=((& docker --version) -join ' ')}catch{$tools.harness='tool-version-unavailable'}
    $outcome=New-HarnessOutcome -RunId $runId -CandidateSha $safeCandidate -Branch $branchName -Dirty $callerDirty -Tools $tools -Commands @($records) -PostgresImageIdentity $imageIdentity -AggregateClassification $aggregate
    Write-HarnessOutcome -Outcome $outcome -Path (Join-Path $runDir 'outcome.json')|Out-Null; Write-HarnessTrace -Events @($records) -Path (Join-Path $runDir 'trace.jsonl')|Out-Null
    $prior=Get-ChildItem -LiteralPath $runRoot -Filter outcome.json -Recurse -ErrorAction SilentlyContinue|Where-Object{$_.FullName -ne (Join-Path $runDir 'outcome.json')}|ForEach-Object{try{Get-Content -Raw $_.FullName|ConvertFrom-Json}catch{$null}}|Where-Object{$null -ne $_ -and $_.candidate_sha -eq $safeCandidate}|Select-Object -Last 1
    if($null -ne $prior){[ordered]@{candidate_sha=$safeCandidate;projection_equal=(Compare-HarnessOutcomeProjection $prior $outcome)}|ConvertTo-Json|Set-Content -LiteralPath (Join-Path $runDir 'comparison.json') -Encoding utf8}
    Write-Output "status=$aggregate"; Write-Output "run_dir=scripts/.runs/$runId"; if($aggregate -ne 'passed'){exit 1}
  }
}

try {
  switch ($Command) {
    'unit' { Invoke-Unit }
    'integration' { Invoke-Integration }
    'live' { Invoke-Live }
    'browser' { Invoke-Browser }
    'provider-write' { Invoke-ProviderWrite }
    'governance-validate' { Invoke-Governance -Mode validate }
    'governance-drift' { Invoke-Governance -Mode drift }
    'governance' { Invoke-Governance -Mode all }
    'context-compile' { Invoke-Context -Mode compile }
    'context-validate' { Invoke-Context -Mode validate }
    'cold' { Invoke-Cold }
  }
} catch {
  Write-Output "status=blocked"
  Write-Output $_.Exception.Message
  exit 1
}
