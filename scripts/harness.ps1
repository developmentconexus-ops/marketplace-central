[CmdletBinding()]
param(
  [ValidateSet('unit', 'integration', 'live', 'browser', 'provider-write', 'governance-validate', 'governance-drift', 'governance', 'context-compile', 'context-validate')]
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
)

$ErrorActionPreference = 'Stop'
$repoRoot = Split-Path -Parent $PSScriptRoot
$runRoot = Join-Path $PSScriptRoot '.runs'
$runId = [guid]::NewGuid().ToString('N')
$runDir = Join-Path $runRoot $runId

Import-Module (Join-Path $PSScriptRoot 'harness/Environment.psm1') -Force
Import-Module (Join-Path $PSScriptRoot 'harness/Execution.psm1') -Force

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
  $url = $DatabaseUrl
  if ([string]::IsNullOrWhiteSpace($url)) { throw 'integration requires explicit -DatabaseUrl (MPC_TEST_DATABASE_URL); ambient configuration is ignored until F-03' }
  try { $uri = [Uri]$url } catch { throw 'integration database target is invalid' }
  if ($uri.Scheme -notin @('postgres', 'postgresql') -or $uri.AbsolutePath -notmatch '/mpc_test_[A-Za-z0-9_-]+$') {
    throw 'integration database must be an mpc_test_* PostgreSQL target'
  }
  Write-Output 'target=ephemeral-postgres'
  Write-Output 'key=MPC_TEST_DATABASE_URL'
  Write-Output 'migrations=delegated'
  throw 'integration blocked until F-03 supplies explicit ephemeral-postgres lifecycle; no database was contacted'
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
  }
} catch {
  Write-Output "status=blocked"
  Write-Output $_.Exception.Message
  exit 1
}
