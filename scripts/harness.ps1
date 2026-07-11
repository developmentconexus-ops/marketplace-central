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

function Get-EnvFileMap {
  param([string]$Path)
  $result = @{}
  if (-not (Test-Path -LiteralPath $Path -PathType Leaf)) { return $result }
  foreach ($line in Get-Content -LiteralPath $Path) {
    if ($line -match '^\s*#' -or $line -notmatch '^\s*([A-Za-z_][A-Za-z0-9_]*)\s*=') { continue }
    $key = $Matches[1]
    $value = ($line -split '=', 2)[1].Trim()
    if ($value.Length -ge 2 -and (($value.StartsWith('"') -and $value.EndsWith('"')) -or ($value.StartsWith("'") -and $value.EndsWith("'")))) {
      $value = $value.Substring(1, $value.Length - 2)
    }
    $result[$key] = $value
  }
  return $result
}

function Get-ConfiguredValue {
  param([hashtable]$FileValues, [string]$Key)
  $processValue = [Environment]::GetEnvironmentVariable($Key)
  if (-not [string]::IsNullOrWhiteSpace($processValue)) { return $processValue }
  if ($FileValues.ContainsKey($Key)) { return [string]$FileValues[$Key] }
  $legacy = @{
    MPC_ORACLE_USERNAME = 'SANKHYA_ORACLE_USER'
    MPC_ORACLE_PASSWORD = 'SANKHYA_ORACLE_PASSWORD'
    MPC_ORACLE_CONNECT_STRING = 'SANKHYA_ORACLE_CONNECT_STRING'
  }
  if ($legacy.ContainsKey($Key)) {
    $legacyKey = $legacy[$Key]
    $legacyProcessValue = [Environment]::GetEnvironmentVariable($legacyKey)
    if (-not [string]::IsNullOrWhiteSpace($legacyProcessValue)) { return $legacyProcessValue }
    if ($FileValues.ContainsKey($legacyKey)) { return [string]$FileValues[$legacyKey] }
  }
  if ($Key -eq 'MPC_ORACLE_CONNECT_STRING') {
    $oracleHost = [Environment]::GetEnvironmentVariable('SANKHYA_ORACLE_HOST')
    $oraclePort = [Environment]::GetEnvironmentVariable('SANKHYA_ORACLE_PORT')
    $oracleService = [Environment]::GetEnvironmentVariable('SANKHYA_ORACLE_SERVICE_NAME')
    if ([string]::IsNullOrWhiteSpace($oracleHost)) { $oracleHost = [string]$FileValues['SANKHYA_ORACLE_HOST'] }
    if ([string]::IsNullOrWhiteSpace($oraclePort)) { $oraclePort = [string]$FileValues['SANKHYA_ORACLE_PORT'] }
    if ([string]::IsNullOrWhiteSpace($oracleService)) { $oracleService = [string]$FileValues['SANKHYA_ORACLE_SERVICE_NAME'] }
    if (-not [string]::IsNullOrWhiteSpace($oracleHost) -and -not [string]::IsNullOrWhiteSpace($oraclePort) -and -not [string]::IsNullOrWhiteSpace($oracleService)) {
      return "$oracleHost`:$oraclePort/$oracleService"
    }
  }
  return ''
}

function Get-ForbiddenUnitEnvironmentKeys {
  $patterns = @(
    '^(?:MC|MS|MPC|MPC_TEST)_DATABASE_URL$',
    '^(?:DATABASE_URL|POSTGRES(?:QL)?_URL|REDIS_URL|MYSQL_URL)$',
    '^PG(?:HOST|PORT|DATABASE|USER|PASSWORD|SERVICE|SSLMODE)$',
    '^(?:MC|MS|MPC)_TENANT_ID$',
    '^MPC_PROVIDER_.+',
    '^MPC_ORACLE_.+',
    '^SANKHYA_ORACLE_.+',
    '^MPC_OAUTH_.+',
    '^(?:HTTP|HTTPS|ALL|NO)_PROXY$',
    '^MPC_WEB_(?:ORIGIN|PROXY_TARGET)$',
    '^(?:MPC|SANKHYA|NGROK|TUNNEL|CLOUDFLARE)_.*(?:PROXY|TUNNEL|AUTHTOKEN|LIVE|PROD|PRODUCTION|TARGET).*',
    '^(?:RUN_MIGRATIONS|MC_MIGRATIONS_DIR)$',
    '^MPC_.*MIGRATIONS?.*'
  )
  $present = @()
  foreach ($entry in [Environment]::GetEnvironmentVariables().GetEnumerator()) {
    $name = [string]$entry.Key
    $value = [string]$entry.Value
    if ([string]::IsNullOrWhiteSpace($value)) { continue }
    if ($patterns | Where-Object { $name -match $_ }) { $present += $name }
  }
  return @($present | Sort-Object -Unique)
}

function Assert-UnitEnvironment {
  $present = Get-ForbiddenUnitEnvironmentKeys
  if ($present.Count -gt 0) { throw "unit environment rejected; forbidden_keys=$($present -join ',')" }
}

function Protect-OutputLine {
  param([object]$Line)
  $safe = [string]$Line
  $safe = [regex]::Replace($safe, '(?i)(\b(?:MPC|SANKHYA|MC|MS)_[A-Z0-9_]*(?:PASSWORD|SECRET|TOKEN|CREDENTIAL|CONNECT_STRING)\b\s*=\s*)[^,;\s]+', '$1[redacted]')
  $safe = [regex]::Replace($safe, '(?i)(\b(?:postgres(?:ql)?|mysql|oracle)://)[^\s,;]+', '$1[redacted]')
  return $safe
}

function Invoke-CapturedProcess {
  param([string]$FilePath, [string[]]$Arguments)
  $captured = @(& $FilePath @Arguments 2>&1)
  $exitCode = $LASTEXITCODE
  $safe = @($captured | ForEach-Object { Protect-OutputLine $_ })
  return [pscustomobject]@{ ExitCode = $exitCode; Output = $safe }
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
  Assert-UnitEnvironment
  Write-Output 'target=fake'
  Write-Output 'env=ignored'
  Write-Output 'postgres=disabled'
  Write-Output 'oracle=disabled'
  Write-Output 'provider_network=disabled'
  Write-Output 'migrations=disabled'
  if ($PreflightOnly) { Write-Summary -TargetType 'fake' -Status 'ready'; return }

  Push-Location (Join-Path $repoRoot 'apps/server_core')
  try {
    $previousGoCache = [Environment]::GetEnvironmentVariable('GOCACHE', 'Process')
    [Environment]::SetEnvironmentVariable('GOCACHE', (Join-Path (Get-Location) '.gocache'), 'Process')
    $go = Invoke-CapturedProcess -FilePath 'go' -Arguments @('test', './tests/unit/...', '-count=1')
    [Environment]::SetEnvironmentVariable('GOCACHE', $previousGoCache, 'Process')
    $go.Output | ForEach-Object { Write-Output $_ }
    if ($go.ExitCode -ne 0) { throw 'unit Go tests failed' }
  }
  finally { Pop-Location }
  $web = Invoke-CapturedProcess -FilePath 'npm' -Arguments @('run', 'test', '--workspace', '@marketplace-central/web', '--', '--run')
  $web.Output | ForEach-Object { Write-Output $_ }
  if ($web.ExitCode -ne 0) { throw 'unit web tests failed' }
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
  if ($EnvFile) { $path = $EnvFile } else { $path = Join-Path $repoRoot '.env' }
  $values = Get-EnvFileMap -Path $path
  $allowlist = @{
    oracle = @('MPC_ORACLE_USERNAME', 'MPC_ORACLE_PASSWORD', 'MPC_ORACLE_CONNECT_STRING', 'MPC_ORACLE_LIB_DIR', 'MPC_ORACLE_LIVE_TEST', 'SANKHYA_ORACLE_HOST', 'SANKHYA_ORACLE_PORT', 'SANKHYA_ORACLE_SERVICE_NAME', 'SANKHYA_ORACLE_USER', 'SANKHYA_ORACLE_PASSWORD')
    provider = @('MPC_PROVIDER_MERCADOLIVRE_CLIENT_ID', 'MPC_PROVIDER_MERCADOLIVRE_CLIENT_SECRET', 'MPC_PROVIDER_MERCADOLIVRE_ACCESS_TOKEN', 'MPC_PROVIDER_MERCADOLIVRE_LIVE_TEST')
  }
  if (-not $allowlist.ContainsKey($Target)) { throw 'live target must be oracle or provider' }
  $required = if ($Target -eq 'oracle') { @('MPC_ORACLE_USERNAME', 'MPC_ORACLE_PASSWORD', 'MPC_ORACLE_CONNECT_STRING') } else { @('MPC_PROVIDER_MERCADOLIVRE_CLIENT_ID', 'MPC_PROVIDER_MERCADOLIVRE_CLIENT_SECRET') }
  $missing = @($required | Where-Object { [string]::IsNullOrWhiteSpace((Get-ConfiguredValue -FileValues $values -Key $_)) })
  Write-Output "target=live-$Target"
  foreach ($key in $allowlist[$Target]) {
    if (-not [string]::IsNullOrWhiteSpace((Get-ConfiguredValue -FileValues $values -Key $key))) { Write-Output "key=$key" }
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
  Write-Output (Protect-OutputLine $_.Exception.Message)
  exit 1
}
