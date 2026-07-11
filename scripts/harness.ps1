[CmdletBinding()]
param(
  [ValidateSet('unit', 'integration', 'live', 'browser', 'provider-write', 'governance-validate', 'governance-drift', 'governance', 'context-compile', 'context-validate', 'impact')]
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
Import-Module (Join-Path $PSScriptRoot 'harness/Postgres.psm1') -Force
Import-Module (Join-Path $PSScriptRoot 'harness/Evidence.psm1') -Force
Import-Module (Join-Path $PSScriptRoot 'harness/Impact.psm1') -Force

function Resolve-HarnessApplication {
  param([Parameter(Mandatory)][string]$Name)
  $application = Get-Command $Name -CommandType Application -ErrorAction SilentlyContinue | Select-Object -First 1
  if ($null -eq $application -or [string]::IsNullOrWhiteSpace([string]$application.Source)) { throw "HEXEC_FILE_NOT_FOUND tool=$Name" }
  [IO.Path]::GetFullPath([string]$application.Source)
}

function Get-HarnessNpmInvocation {
  $npmPath = Resolve-HarnessApplication -Name 'npm'
  if ([Environment]::OSVersion.Platform -ne [PlatformID]::Win32NT) { return [pscustomobject]@{ FilePath = $npmPath; ArgumentPrefix = @() } }
  $nodePath = Resolve-HarnessApplication -Name 'node'
  $npmCliPath = Join-Path (Split-Path -Parent $npmPath) 'node_modules/npm/bin/npm-cli.js'
  if (-not (Test-Path -LiteralPath $npmCliPath -PathType Leaf)) { throw 'HEXEC_FILE_NOT_FOUND tool=npm' }
  [pscustomobject]@{ FilePath = $nodePath; ArgumentPrefix = @([IO.Path]::GetFullPath($npmCliPath)) }
}

function Write-HarnessProcessResult {
  param([Parameter(Mandatory)][object]$Result)
  if (-not [string]::IsNullOrEmpty([string]$Result.Stdout)) { Write-Output ([string]$Result.Stdout).TrimEnd("`r", "`n") }
  if (-not [string]::IsNullOrEmpty([string]$Result.Stderr)) { Write-Output ([string]$Result.Stderr).TrimEnd("`r", "`n") }
}

function Write-Summary {
  param([string]$TargetType, [string]$Status)
  New-Item -ItemType Directory -Path $runDir -Force | Out-Null
  @("target=$TargetType", "status=$Status", "run_id=$runId") | Set-Content -LiteralPath (Join-Path $runDir 'summary.txt') -Encoding utf8
  Write-Output "target=$TargetType"; Write-Output "status=$Status"; Write-Output "run_id=$runId"; Write-Output "run_dir=scripts/.runs/$runId"
}

function Invoke-Unit {
  $goPath = Resolve-HarnessApplication -Name 'go'; $npm = Get-HarnessNpmInvocation
  $environment = New-HarnessChildEnvironment -RepositoryRoot $repoRoot -LaneId 'unit'
  Write-Output 'target=fake'; Write-Output 'env=ignored'; Write-Output 'postgres=disabled'; Write-Output 'oracle=disabled'; Write-Output 'provider_network=disabled'; Write-Output 'migrations=disabled'
  if ($PreflightOnly) { Write-Summary -TargetType 'fake' -Status 'ready'; return }
  $go = Invoke-HarnessProcess -Request (New-HarnessProcessRequest -FilePath $goPath -ArgumentList @('test', './tests/unit/...', '-count=1') -WorkingDirectory ([IO.Path]::GetFullPath((Join-Path $repoRoot 'apps/server_core'))) -Environment $environment -TimeoutSeconds 1200)
  Write-HarnessProcessResult $go; if ($go.ExitCode -ne 0) { throw "unit Go tests failed reason=$($go.ReasonCode) exit_code=$($go.ExitCode)" }
  $web = Invoke-HarnessProcess -Request (New-HarnessProcessRequest -FilePath $npm.FilePath -ArgumentList (@($npm.ArgumentPrefix) + @('run', 'test', '--workspace', '@marketplace-central/web', '--', '--run')) -WorkingDirectory $repoRoot -Environment $environment -TimeoutSeconds 1200)
  Write-HarnessProcessResult $web; if ($web.ExitCode -ne 0) { throw "unit web tests failed reason=$($web.ReasonCode) exit_code=$($web.ExitCode)" }
  Write-Summary -TargetType 'fake' -Status 'passed'
}

function Invoke-Integration {
  if (-not [string]::IsNullOrWhiteSpace($DatabaseUrl)) { throw 'HPG_EXTERNAL_TARGET_FORBIDDEN' }
  Write-Output 'target=ephemeral-postgres'; Write-Output 'key=MPC_TEST_DATABASE_URL'; Write-Output 'migrations=embedded'
  if ($PreflightOnly) { Write-Output 'status=ready'; return }
  $passwordBytes = [byte[]]::new(24); [Security.Cryptography.RandomNumberGenerator]::Fill($passwordBytes)
  $spec = New-HarnessPostgresRunSpec -RepositoryRoot $repoRoot -RunId $runId -Password ([Convert]::ToHexString($passwordBytes).ToLowerInvariant()) -DockerFilePath (Resolve-HarnessPostgresDockerApplication)
  $result = Invoke-HarnessPostgresLifecycle -RunSpec $spec -BaseEnvironment (New-HarnessChildEnvironment -RepositoryRoot $repoRoot -LaneId 'integration') -GoFilePath (Resolve-HarnessApplication -Name 'go') -TimeoutSeconds 1200
  Write-Output "migrations_first=$($result.MigrationsAppliedFirst)"; Write-Output "migrations_second=$($result.MigrationsAppliedSecond)"; Write-Output "resource_count=$(@($result.ResourceInventory).Count)"; Write-Output "port=$($result.HostPort)"
  if ($result.ExitCode -ne 0) { throw "postgres lifecycle failed reasons=$((@($result.PrimaryReasonCode) + @($result.CleanupReasonCodes) | Where-Object { $_ } ) -join ',') exit_code=$($result.ExitCode)" }
  Write-Summary -TargetType 'ephemeral-postgres' -Status 'passed'
}

function Invoke-Live {
  $path = if ($EnvFile) { if ([IO.Path]::IsPathFullyQualified($EnvFile)) { [IO.Path]::GetFullPath($EnvFile) } else { [IO.Path]::GetFullPath((Join-Path $repoRoot $EnvFile)) } } else { Join-Path $repoRoot '.env' }
  if ($Target -notin @('oracle', 'provider')) { throw 'live target must be oracle or provider' }
  $laneId = if ($Target -eq 'oracle') { 'live-oracle' } else { 'live-provider-read' }
  $environment = New-HarnessChildEnvironment -RepositoryRoot $repoRoot -LaneId $laneId -EnvFile $path
  $required = if ($Target -eq 'oracle') { @('MPC_ORACLE_USERNAME', 'MPC_ORACLE_PASSWORD', 'MPC_ORACLE_CONNECT_STRING') } else { @('MPC_PROVIDER_MERCADOLIVRE_CLIENT_ID', 'MPC_PROVIDER_MERCADOLIVRE_CLIENT_SECRET') }
  $missing = @($required | Where-Object { -not $environment.ContainsKey($_) -or [string]::IsNullOrWhiteSpace($environment[$_]) })
  Write-Output "target=live-$Target"; if ($missing.Count -gt 0) { throw "live preflight missing_keys=$($missing -join ',')" }; Write-Output 'provider_write=disabled'; Write-Summary -TargetType "live-$Target" -Status 'ready'
}

function Invoke-Browser { Write-Output 'target=browser'; Write-Output 'provider_network=disabled'; if (-not $PreflightOnly) { throw 'browser runner is not configured; use browser automation as an explicit lane' }; Write-Summary -TargetType 'browser' -Status 'ready' }
function Invoke-ProviderWrite { Write-Output 'target=live-provider'; if ([string]::IsNullOrWhiteSpace($Provider)) { throw 'provider is required' }; if ([string]::IsNullOrWhiteSpace($Actor) -or [string]::IsNullOrWhiteSpace($IdempotencyKey)) { throw 'provider write requires actor and idempotency_key' }; if (-not $Execute) { throw 'explicit -Execute is required before network' }; throw 'provider write adapter is intentionally outside F-02; no network was invoked' }

function Write-GovernanceResult {
  param([object]$Result)
  if ($Result.Passed) { Write-Output 'status=passed' } else { Write-Output 'status=failed'; foreach ($violation in @($Result.Violations)) { Write-Output "error_code=$($violation.ErrorCode)"; Write-Output "id=$($violation.Id)"; if ($violation.Path) { Write-Output "path=$($violation.Path)" } } }
  foreach ($exception in @($Result.BaselineExceptions)) { Write-Output "baseline_exception=$($exception.Id)" }; Write-Output 'artifact_path=contracts/governance'
}

function Invoke-Governance {
  param([ValidateSet('validate', 'drift', 'all')][string]$Mode)
  Import-Module (Join-Path $PSScriptRoot 'harness/Policy.psm1') -Force
  if ($Mode -eq 'validate') { $result = Test-GovernanceContracts -RepositoryRoot $repoRoot } else {
    if ([string]::IsNullOrWhiteSpace($BaseSha) -or $BaseSha -notmatch '^[0-9a-f]{40}$') { Write-Output 'status=failed'; Write-Output 'error_code=GOV_SEMANTIC_DRIFT'; Write-Output 'id=base-sha-invalid'; exit 1 }
    if ($Mode -eq 'all') { $contracts = Test-GovernanceContracts -RepositoryRoot $repoRoot; if (-not $contracts.Passed) { Write-GovernanceResult $contracts; exit 1 } }
    $result = Test-GovernanceDrift -RepositoryRoot $repoRoot -BaseSha $BaseSha
  }
  Write-GovernanceResult $result; if (-not $result.Passed) { exit 1 }
}

function Invoke-Context {
  param([ValidateSet('compile', 'validate')][string]$Mode)
  Import-Module (Join-Path $PSScriptRoot 'harness/Context.psm1') -Force
  if ($Mode -eq 'compile') { $path = "scripts/.runs/$runId/context-pack.json"; $result = if ([string]::IsNullOrWhiteSpace($FeaturePath) -or @($AllowedPath).Count -eq 0) { [pscustomobject]@{ Passed=$false; Status='failed'; ErrorCode='CTX_FEATURE_INVALID'; Id='compile-input'; Path='' } } else { New-HarnessContextPack -FeaturePath $FeaturePath -BaseSha $BaseSha -AllowedPath $AllowedPath -OutputPath (Join-Path $repoRoot $path) } }
  else {
    $path = if ($ContextPath) { $ContextPath } else { 'context-pack.json' }
    if ([string]::IsNullOrWhiteSpace($ContextPath)) { $result = [pscustomobject]@{ Passed=$false; Status='failed'; ErrorCode='CTX_SOURCE_MISSING'; Id='context-pack'; Path='' } }
    else { $resolvedContextPath = if ([IO.Path]::IsPathRooted($ContextPath)) { $ContextPath } else { Join-Path $repoRoot $ContextPath }; $result = Test-HarnessContextPack -Path $resolvedContextPath -RepositoryRoot $repoRoot -RequireCurrentBase:$RequireCurrentBase }
  }
  $displayPath = if ($result.Path) { [string]$result.Path } else { $path }
  Write-Output "status=$($result.Status)"; if (-not $result.Passed) { Write-Output "error_code=$($result.ErrorCode)"; Write-Output "id=$($result.Id)"; Write-Output "path=$displayPath" }; Write-Output "artifact_path=$path"; if (-not $result.Passed) { exit 1 }
}

function Invoke-Impact {
  if ([string]::IsNullOrWhiteSpace($ContextPath)) { throw 'HIMPACT_CONTEXT_REQUIRED' }
  $result = Invoke-HarnessImpactGate -RepositoryRoot $repoRoot -ContextPath $ContextPath -RunId $runId -RunDirectory $runDir
  Write-Output "status=$($result.Status)"; Write-Output "outcome_path=$($result.OutcomePath)"; foreach ($seam in @($result.SharedSeams)) { Write-Output "shared_seam=$seam" }; if (-not $result.Passed) { Write-Output "error_code=$($result.ErrorCode)"; exit 1 }
}

try {
  switch ($Command) {
    'unit' { Invoke-Unit }; 'integration' { Invoke-Integration }; 'live' { Invoke-Live }; 'browser' { Invoke-Browser }; 'provider-write' { Invoke-ProviderWrite }
    'governance-validate' { Invoke-Governance -Mode validate }; 'governance-drift' { Invoke-Governance -Mode drift }; 'governance' { Invoke-Governance -Mode all }
    'context-compile' { Invoke-Context -Mode compile }; 'context-validate' { Invoke-Context -Mode validate }; 'impact' { Invoke-Impact }
  }
} catch { Write-Output 'status=blocked'; Write-Output $_.Exception.Message; exit 1 }
