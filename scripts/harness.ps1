[CmdletBinding()]
param(
  [ValidateSet('unit', 'integration', 'live', 'browser', 'provider-write')]
  [string]$Command = 'unit',
  [switch]$PreflightOnly,
  [string]$EnvFile,
  [string]$Target = 'oracle',
  [string]$DatabaseUrl,
  [string]$Provider,
  [string]$Actor,
  [string]$IdempotencyKey,
  [switch]$Execute
)

$ErrorActionPreference = 'Stop'
$repoRoot = Split-Path -Parent $PSScriptRoot
$runRoot = Join-Path $PSScriptRoot '.runs'
$runId = [guid]::NewGuid().ToString('N')
$runDir = Join-Path $runRoot $runId
New-Item -ItemType Directory -Path $runDir -Force | Out-Null

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
    if ($FileValues.ContainsKey($legacyKey)) { return [string]$FileValues[$legacyKey] }
  }
  if ($Key -eq 'MPC_ORACLE_CONNECT_STRING' -and $FileValues.ContainsKey('SANKHYA_ORACLE_HOST') -and $FileValues.ContainsKey('SANKHYA_ORACLE_PORT') -and $FileValues.ContainsKey('SANKHYA_ORACLE_SERVICE_NAME')) {
    return "$($FileValues['SANKHYA_ORACLE_HOST']):$($FileValues['SANKHYA_ORACLE_PORT'])/$($FileValues['SANKHYA_ORACLE_SERVICE_NAME'])"
  }
  return ''
}

function Assert-UnitEnvironment {
  $forbidden = @(
    'MC_DATABASE_URL', 'MS_DATABASE_URL', 'MS_TENANT_ID', 'MC_MIGRATIONS_DIR', 'RUN_MIGRATIONS',
    'MPC_ORACLE_USERNAME', 'MPC_ORACLE_PASSWORD', 'MPC_ORACLE_CONNECT_STRING', 'MPC_ORACLE_LIB_DIR', 'MPC_ORACLE_LIVE_TEST',
    'SANKHYA_ORACLE_HOST', 'SANKHYA_ORACLE_PORT', 'SANKHYA_ORACLE_SERVICE_NAME', 'SANKHYA_ORACLE_USER', 'SANKHYA_ORACLE_PASSWORD',
    'MPC_PROVIDER_MERCADOLIVRE_CLIENT_ID', 'MPC_PROVIDER_MERCADOLIVRE_CLIENT_SECRET', 'MPC_PROVIDER_MERCADOLIVRE_ACCESS_TOKEN', 'MPC_PROVIDER_MERCADOLIVRE_LIVE_TEST',
    'MPC_OAUTH_HMAC_SECRET', 'MPC_OAUTH_REDIRECT_URI', 'MPC_WEB_ORIGIN', 'HTTP_PROXY', 'HTTPS_PROXY', 'ALL_PROXY', 'NO_PROXY'
  )
  $present = @($forbidden | Where-Object { -not [string]::IsNullOrWhiteSpace([Environment]::GetEnvironmentVariable($_)) })
  if ($present.Count -gt 0) { throw "unit environment rejected; forbidden_keys=$($present -join ',')" }
}

function Write-Summary {
  param([string]$TargetType, [string]$Status)
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
  try { & go test ./tests/unit/... -count=1; if ($LASTEXITCODE -ne 0) { throw 'unit Go tests failed' } }
  finally { Pop-Location }
  & npm run test --workspace @marketplace-central/web -- --run
  if ($LASTEXITCODE -ne 0) { throw 'unit web tests failed' }
  Write-Summary -TargetType 'fake' -Status 'passed'
}

function Invoke-Integration {
  $url = if ($DatabaseUrl) { $DatabaseUrl } else { [Environment]::GetEnvironmentVariable('MPC_TEST_DATABASE_URL') }
  if ([string]::IsNullOrWhiteSpace($url)) { throw 'integration requires MPC_TEST_DATABASE_URL (ephemeral-postgres)' }
  try { $uri = [Uri]$url } catch { throw 'integration database target is invalid' }
  if ($uri.Scheme -notin @('postgres', 'postgresql') -or $uri.AbsolutePath -notmatch '/mpc_test_[A-Za-z0-9_-]+$') {
    throw 'integration database must be an mpc_test_* PostgreSQL target'
  }
  Write-Output 'target=ephemeral-postgres'
  Write-Output 'migrations=delegated'
  if ($PreflightOnly) { Write-Summary -TargetType 'ephemeral-postgres' -Status 'ready'; return }
  $old = [Environment]::GetEnvironmentVariable('MC_DATABASE_URL')
  try {
    [Environment]::SetEnvironmentVariable('MC_DATABASE_URL', $url, 'Process')
    Push-Location (Join-Path $repoRoot 'apps/server_core')
    try { & go test ./tests/integration/... -count=1; if ($LASTEXITCODE -ne 0) { throw 'integration tests failed' } }
    finally { Pop-Location }
  } finally { [Environment]::SetEnvironmentVariable('MC_DATABASE_URL', $old, 'Process') }
  Write-Summary -TargetType 'ephemeral-postgres' -Status 'passed'
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

try {
  switch ($Command) {
    'unit' { Invoke-Unit }
    'integration' { Invoke-Integration }
    'live' { Invoke-Live }
    'browser' { Invoke-Browser }
    'provider-write' { Invoke-ProviderWrite }
  }
} catch {
  Write-Output "status=blocked"
  Write-Output (($_.Exception.Message -replace '(?i)(secret|password|token|credential)[^,; ]*', '$1 [redacted]'))
  exit 1
}
