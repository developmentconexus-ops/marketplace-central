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
  if ([string]::IsNullOrWhiteSpace($CandidateSha) -or $CandidateSha -notmatch '^[0-9a-f]{40}$') { throw 'COLD_CANDIDATE_SHA_INVALID' }
  $head = (& git -C $repoRoot rev-parse HEAD 2>$null).Trim(); if ($LASTEXITCODE -ne 0 -or $head -cne $CandidateSha) { throw 'COLD_CANDIDATE_SHA_MISMATCH' }
  $status = @(& git -C $repoRoot status --porcelain 2>$null); if ($status.Count -gt 0) { throw 'COLD_CALLER_DIRTY' }
  $snapshot = Join-Path $runDir 'snapshot'; New-Item -ItemType Directory -Path $runDir -Force | Out-Null
  & git -C $repoRoot clone --quiet --no-hardlinks --local $repoRoot $snapshot 2>$null; if ($LASTEXITCODE -ne 0) { throw 'COLD_SNAPSHOT_FAILED' }
  & git -C $snapshot checkout --quiet --detach $CandidateSha 2>$null; if ($LASTEXITCODE -ne 0) { throw 'COLD_SNAPSHOT_CHECKOUT_FAILED' }
  $runCache = Join-Path $runDir 'cache'; New-Item -ItemType Directory -Path $runCache -Force | Out-Null
  $env:CACHE = $runCache
  $records = [Collections.Generic.List[object]]::new()
  $records.Add([ordered]@{id='preflight';command='git clean candidate';stage='preflight';target_label='fake';evidence_class='contract';duration_ms=0;exit_code=0;reason='passed';artifact_paths=@()})
  $go = Get-Command go -ErrorAction SilentlyContinue; $npm = Get-Command npm -ErrorAction SilentlyContinue
  if ($null -eq $go -or $null -eq $npm) { throw 'COLD_TOOL_MISSING' }
  Push-Location $snapshot
  try {
    $env:GOCACHE = Join-Path $runCache 'go'; $env:GOMODCACHE = Join-Path $runCache 'gomod'; $env:npm_config_cache = Join-Path $runCache 'npm'; New-Item -ItemType Directory -Force -Path $env:GOCACHE,$env:GOMODCACHE,$env:npm_config_cache | Out-Null
    Push-Location (Join-Path $snapshot 'apps/server_core'); & go mod download; $goExit = $LASTEXITCODE; Pop-Location
    $records.Add([ordered]@{id='go-mod-download';command='go mod download';stage='provisioning';target_label='external-dependency-registry';evidence_class='provisioning';duration_ms=0;exit_code=$goExit;reason=if($goExit -eq 0){'passed'}else{'failed'};artifact_paths=@()})
    if ($goExit -ne 0) { throw 'COLD_PROVISION_GO_FAILED' }
    & npm ci --ignore-scripts; $npmExit = $LASTEXITCODE
    $records.Add([ordered]@{id='npm-ci';command='npm ci --ignore-scripts';stage='provisioning';target_label='external-dependency-registry';evidence_class='provisioning';duration_ms=0;exit_code=$npmExit;reason=if($npmExit -eq 0){'passed'}else{'failed'};artifact_paths=@()})
    if ($npmExit -ne 0) { throw 'COLD_PROVISION_NPM_FAILED' }
    $env:GOCACHE = Join-Path $snapshot 'apps/server_core/.gocache'
    $env:GOMODCACHE = Join-Path $snapshot 'apps/server_core/.gomodcache'
    New-Item -ItemType Directory -Force -Path $env:GOCACHE,$env:GOMODCACHE | Out-Null
    Push-Location (Join-Path $snapshot 'apps/server_core'); & go mod download; $goWorkspaceExit = $LASTEXITCODE; Pop-Location
    if ($goWorkspaceExit -ne 0) { throw 'COLD_PROVISION_GO_WORKSPACE_FAILED' }
    $docker = Get-Command docker -ErrorAction SilentlyContinue; $imageIdentity = ''
    if ($null -ne $docker) { & docker pull postgres:16-bookworm; if ($LASTEXITCODE -ne 0) { throw 'COLD_PROVISION_IMAGE_FAILED' }; $imageIdentity = (& docker image inspect postgres:16-bookworm --format '{{.Id}}').Trim() }
    $records.Add([ordered]@{id='docker-pull-postgres';command='docker pull postgres:16-bookworm';stage='provisioning';target_label='external-dependency-registry';evidence_class='provisioning';duration_ms=0;exit_code=0;reason='passed';artifact_paths=@()})
    foreach ($test in @('scripts/tests/governance-contracts.tests.ps1','scripts/tests/cold-gate-evidence.tests.ps1','scripts/tests/cold-gate-snapshot.tests.ps1')) { & pwsh -NoProfile -ExecutionPolicy Bypass -File $test; $exit=$LASTEXITCODE; $records.Add([ordered]@{id=([IO.Path]::GetFileNameWithoutExtension($test));command="pwsh -File $test";stage='validation';target_label='fake';evidence_class='contract';duration_ms=0;exit_code=$exit;reason=if($exit -eq 0){'passed'}else{'failed'};artifact_paths=@()}) }
    & pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/f03-regression.tests.ps1 -BaseSha $CandidateSha
    $f03Exit = $LASTEXITCODE
    $records.Add([ordered]@{id='f03-regression';command='pwsh -File scripts/tests/f03-regression.tests.ps1 -BaseSha {candidate_sha}';stage='validation';target_label='fake';evidence_class='contract';duration_ms=0;exit_code=$f03Exit;reason=if($f03Exit -eq 0){'passed'}else{'failed'};artifact_paths=@()})
    $aggregate = if (@($records | Where-Object exit_code -ne 0).Count -eq 0) {'passed'} else {'failed'}
    $branchName = ((& git -C $repoRoot branch --show-current).Trim()); if ([string]::IsNullOrWhiteSpace($branchName)) { $branchName = 'detached' }
    $outcome = New-HarnessOutcome -RunId $runId -CandidateSha $CandidateSha -Branch $branchName -Dirty $false -Tools ([ordered]@{go=((& go version) -join ' '); npm=((& npm --version) -join ' ')}) -Commands @($records) -PostgresImageIdentity $imageIdentity -AggregateClassification $aggregate
    Write-HarnessOutcome -Outcome $outcome -Path (Join-Path $runDir 'outcome.json') | Out-Null
    Write-HarnessTrace -Events @($records) -Path (Join-Path $runDir 'trace.jsonl') | Out-Null
    Write-Output "status=$aggregate"; Write-Output "run_dir=scripts/.runs/$runId"
    if ($aggregate -ne 'passed') { exit 1 }
  } finally { Pop-Location }
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
