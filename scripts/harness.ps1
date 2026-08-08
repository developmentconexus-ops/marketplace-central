[CmdletBinding()]
param(
  [ValidateSet('unit', 'integration', 'live', 'browser', 'edge', 'provider-write', 'governance-validate', 'governance-drift', 'governance', 'pg-session-up', 'pg-session-down')]
  [string]$Command = 'unit',
  [switch]$PreflightOnly,
  [string]$EnvFile,
  [string]$Target = 'oracle',
  [string]$DatabaseUrl,
  [string]$Provider,
  [string]$Actor,
  [string]$IdempotencyKey,
  [string]$BaseSha,
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
  $dockerPath = Resolve-HarnessPostgresDockerApplication
  $session = Get-HarnessPostgresSession -RepositoryRoot $repoRoot -DockerFilePath $dockerPath
  if ($null -ne $session) {
    Write-Output 'container=session-reuse'
    $spec = New-HarnessPostgresRunSpec -RepositoryRoot $repoRoot -RunId $runId -Password ([string]$session.Password) -DockerFilePath $dockerPath
    $result = Invoke-HarnessPostgresLifecycle -RunSpec $spec -Session $session -BaseEnvironment (New-HarnessChildEnvironment -RepositoryRoot $repoRoot -LaneId 'integration') -GoFilePath (Resolve-HarnessApplication -Name 'go') -TimeoutSeconds 1200
  } else {
    Write-Output 'container=ephemeral'
    $passwordBytes = [byte[]]::new(24); [Security.Cryptography.RandomNumberGenerator]::Fill($passwordBytes)
    $spec = New-HarnessPostgresRunSpec -RepositoryRoot $repoRoot -RunId $runId -Password ([Convert]::ToHexString($passwordBytes).ToLowerInvariant()) -DockerFilePath $dockerPath
    $result = Invoke-HarnessPostgresLifecycle -RunSpec $spec -BaseEnvironment (New-HarnessChildEnvironment -RepositoryRoot $repoRoot -LaneId 'integration') -GoFilePath (Resolve-HarnessApplication -Name 'go') -TimeoutSeconds 1200
  }
  Write-Output "migrations_first=$($result.MigrationsAppliedFirst)"; Write-Output "migrations_second=$($result.MigrationsAppliedSecond)"; Write-Output "resource_count=$(@($result.ResourceInventory).Count)"; Write-Output "port=$($result.HostPort)"
  # -1 means the test step was never reached (the lane died in setup), which is a
  # different fact from 0 and must not print as one.
  Write-Output "tests_run=$($result.TestsRun) tests_passed=$($result.TestsPassed) tests_failed=$($result.TestsFailed)"
  foreach ($token in @($result.FailureDiagnosticTokens)) { Write-Output "failure_token=$token" }
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

function Invoke-Edge {
  # Issue #1 negative fixture. Drives deploy/Caddyfile and
  # docker/dev/oauth-edge.Caddyfile through real Caddy against stub upstreams.
  # Docker only — no ERP, no dev stack, no provider credentials, no secrets,
  # which is why this lane is CI-able (GATE-TOPOLOGY.md §5).
  Write-Output 'target=edge-caddy'; Write-Output 'postgres=disabled'; Write-Output 'oracle=disabled'; Write-Output 'provider_network=disabled'
  if ($PreflightOnly) { Write-Summary -TargetType 'edge-caddy' -Status 'ready'; return }
  $pwshPath = Resolve-HarnessApplication -Name 'pwsh'
  $script = [IO.Path]::GetFullPath((Join-Path $PSScriptRoot 'tests/edge-pii-deny.integration.tests.ps1'))
  if (-not (Test-Path -LiteralPath $script -PathType Leaf)) { throw "HEXEC_FILE_NOT_FOUND script=$script" }
  $environment = New-HarnessChildEnvironment -RepositoryRoot $repoRoot -LaneId 'integration'
  $result = Invoke-HarnessProcess -Request (New-HarnessProcessRequest -FilePath $pwshPath -ArgumentList @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', $script) -WorkingDirectory $repoRoot -Environment $environment -TimeoutSeconds 900)
  Write-HarnessProcessResult $result
  if ($result.ExitCode -ne 0) { throw "edge fixture failed reason=$($result.ReasonCode) exit_code=$($result.ExitCode)" }
  Write-Summary -TargetType 'edge-caddy' -Status 'passed'
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

function Invoke-PgSessionUp {
  $dockerPath = Resolve-HarnessPostgresDockerApplication
  $session = Start-HarnessPostgresSession -RepositoryRoot $repoRoot -DockerFilePath $dockerPath
  Write-Output "container=$($session.ContainerName)"; Write-Output "port=$($session.HostPort)"
  Write-Summary -TargetType 'pg-session' -Status 'ready'
}

function Invoke-PgSessionDown {
  $dockerPath = Resolve-HarnessPostgresDockerApplication
  Stop-HarnessPostgresSession -RepositoryRoot $repoRoot -DockerFilePath $dockerPath
  Write-Summary -TargetType 'pg-session' -Status 'stopped'
}

try {
  switch ($Command) {
    'unit' { Invoke-Unit }; 'integration' { Invoke-Integration }; 'live' { Invoke-Live }; 'browser' { Invoke-Browser }; 'edge' { Invoke-Edge }; 'provider-write' { Invoke-ProviderWrite }
    'governance-validate' { Invoke-Governance -Mode validate }; 'governance-drift' { Invoke-Governance -Mode drift }; 'governance' { Invoke-Governance -Mode all }
    'pg-session-up' { Invoke-PgSessionUp }; 'pg-session-down' { Invoke-PgSessionDown }
  }
} catch { Write-Output 'status=blocked'; Write-Output $_.Exception.Message; exit 1 }
