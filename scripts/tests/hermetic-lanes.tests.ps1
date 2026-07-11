$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$harness = Join-Path $repoRoot 'scripts/harness.ps1'

function Invoke-Harness {
  param([string[]]$Arguments)
  $output = & pwsh -NoProfile -ExecutionPolicy Bypass -File $harness @Arguments 2>&1
  [pscustomobject]@{ ExitCode = $LASTEXITCODE; Output = ($output -join "`n") }
}

function Assert-Contains([string]$Value, [string]$Needle, [string]$Name) {
  if ($Value -notmatch [regex]::Escape($Needle)) { throw "${Name}: expected [$Needle] in [$Value]" }
}

function Assert-NotContains([string]$Value, [string]$Needle, [string]$Name) {
  if ($Value -match [regex]::Escape($Needle)) { throw "${Name}: found forbidden [$Needle] in [$Value]" }
}

if (-not (Test-Path -LiteralPath $harness -PathType Leaf)) {
  throw 'harness script is missing'
}

$previousEnvironment = @{}
foreach ($name in @('MC_DATABASE_URL','MPC_TEST_DATABASE_URL','DATABASE_URL','PGHOST','MS_TENANT_ID','MPC_PROVIDER_FOO','MPC_ORACLE_PASSWORD','SANKHYA_ORACLE_USER','MPC_WEB_PROXY_TARGET','NGROK_AUTHTOKEN')) {
  $previousEnvironment[$name] = [Environment]::GetEnvironmentVariable($name)
}
try {
  $env:MC_DATABASE_URL = 'postgres://secret.example/db'
  $unitBlocked = Invoke-Harness @('-Command', 'unit', '-PreflightOnly')
  if ($unitBlocked.ExitCode -eq 0) { throw 'unit preflight must reject inherited live configuration' }

  $env:MPC_TEST_DATABASE_URL = 'postgres://ambient-secret.example/mpc_test_ambient'
  $env:DATABASE_URL = 'postgres://database-secret.example/db'
  $env:PGHOST = 'pg-secret.example'
  $env:MS_TENANT_ID = 'tenant-secret'
  $env:MPC_PROVIDER_FOO = 'provider-secret'
  $env:MPC_ORACLE_PASSWORD = 'oracle-secret'
  $env:SANKHYA_ORACLE_USER = 'legacy-secret'
  $env:MPC_WEB_PROXY_TARGET = 'https://proxy-secret.example'
  $env:NGROK_AUTHTOKEN = 'tunnel-secret'
  $unitPatterns = Invoke-Harness @('-Command', 'unit', '-PreflightOnly')
  if ($unitPatterns.ExitCode -eq 0) { throw 'unit preflight must reject provider/oracle/proxy/tunnel/database configuration' }
  foreach ($key in @('MPC_TEST_DATABASE_URL','DATABASE_URL','PGHOST','MS_TENANT_ID','MPC_PROVIDER_FOO','MPC_ORACLE_PASSWORD','SANKHYA_ORACLE_USER','MPC_WEB_PROXY_TARGET','NGROK_AUTHTOKEN')) {
    Assert-Contains $unitPatterns.Output $key "unit forbidden key $key"
  }
  foreach ($value in @('ambient-secret.example','database-secret.example','pg-secret.example','tenant-secret','provider-secret','oracle-secret','legacy-secret','proxy-secret.example','tunnel-secret')) {
    Assert-NotContains $unitPatterns.Output $value "unit forbidden value $value"
  }
} finally {
  foreach ($name in $previousEnvironment.Keys) { [Environment]::SetEnvironmentVariable($name, $previousEnvironment[$name], 'Process') }
}

$tempEnv = Join-Path ([IO.Path]::GetTempPath()) ('mpc-hermetic-env-' + [guid]::NewGuid().ToString('N') + '.env')
try {
  @(
    'MPC_ORACLE_USERNAME=alice',
    'MPC_ORACLE_PASSWORD=oracle-secret',
    'MPC_ORACLE_CONNECT_STRING=oracle.example/ORCL',
    'MPC_PROVIDER_MERCADOLIVRE_CLIENT_SECRET=provider-secret'
  ) | Set-Content -LiteralPath $tempEnv -Encoding utf8

  $live = Invoke-Harness @('-Command', 'live', '-PreflightOnly', '-EnvFile', $tempEnv)
  if ($live.ExitCode -ne 0) { throw "live preflight failed: $($live.Output)" }
  Assert-Contains $live.Output 'target=live-oracle' 'live target'
  Assert-Contains $live.Output 'key=MPC_ORACLE_USERNAME' 'live key'
  Assert-NotContains $live.Output 'oracle-secret' 'oracle secret redaction'
  Assert-NotContains $live.Output 'provider-secret' 'provider secret redaction'

  $missing = Invoke-Harness @('-Command', 'live', '-PreflightOnly', '-EnvFile', (Join-Path ([IO.Path]::GetTempPath()) 'missing-mpc.env'))
  if ($missing.ExitCode -eq 0) { throw 'live preflight must reject missing configuration' }
  Assert-Contains $missing.Output 'missing_keys=MPC_ORACLE_USERNAME,MPC_ORACLE_PASSWORD,MPC_ORACLE_CONNECT_STRING' 'missing key names'
  Assert-NotContains $missing.Output 'secret' 'missing output redaction'
} finally {
  Remove-Item -LiteralPath $tempEnv -Force -ErrorAction SilentlyContinue
}

$legacyEnvironment = @{}
foreach ($name in @('SANKHYA_ORACLE_USER','SANKHYA_ORACLE_PASSWORD','SANKHYA_ORACLE_CONNECT_STRING')) {
  $legacyEnvironment[$name] = [Environment]::GetEnvironmentVariable($name)
}
try {
  $env:SANKHYA_ORACLE_USER = 'legacy-user'
  $env:SANKHYA_ORACLE_PASSWORD = 'legacy-password'
  $env:SANKHYA_ORACLE_CONNECT_STRING = 'legacy-connect'
  $legacy = Invoke-Harness @('-Command', 'live', '-Target', 'oracle', '-PreflightOnly')
  if ($legacy.ExitCode -ne 0) { throw "legacy process aliases should satisfy oracle preflight: $($legacy.Output)" }
  foreach ($key in @('MPC_ORACLE_USERNAME','MPC_ORACLE_PASSWORD','MPC_ORACLE_CONNECT_STRING')) {
    Assert-Contains $legacy.Output "key=$key" "legacy alias $key"
  }
  foreach ($value in @('legacy-user','legacy-password','legacy-connect')) {
    Assert-NotContains $legacy.Output $value "legacy value $value"
  }
} finally {
  foreach ($name in $legacyEnvironment.Keys) { [Environment]::SetEnvironmentVariable($name, $legacyEnvironment[$name], 'Process') }
}

$integrationEnvironment = @{}
foreach ($name in @('MPC_TEST_DATABASE_URL','MPC_ORACLE_PASSWORD','MPC_PROVIDER_FOO','MPC_OAUTH_HMAC_SECRET','HTTP_PROXY','MPC_WEB_PROXY_TARGET')) {
  $integrationEnvironment[$name] = [Environment]::GetEnvironmentVariable($name)
}
try {
  $env:MPC_TEST_DATABASE_URL = 'postgres://ambient-secret.example/mpc_test_ambient'
  $env:MPC_ORACLE_PASSWORD = 'oracle-secret'
  $env:MPC_PROVIDER_FOO = 'provider-secret'
  $env:MPC_OAUTH_HMAC_SECRET = 'oauth-secret'
  $env:HTTP_PROXY = 'https://proxy-secret.example'
  $env:MPC_WEB_PROXY_TARGET = 'https://tunnel-secret.example'
  $integrationAmbient = Invoke-Harness @('-Command', 'integration', '-PreflightOnly')
  if ($integrationAmbient.ExitCode -eq 0) { throw 'integration must fail closed without an explicit F-03 database target' }
  Assert-Contains $integrationAmbient.Output 'MPC_TEST_DATABASE_URL' 'integration explicit target requirement'
  foreach ($value in @('ambient-secret.example','oracle-secret','provider-secret','oauth-secret','proxy-secret.example','tunnel-secret.example')) {
    Assert-NotContains $integrationAmbient.Output $value "integration value $value"
  }
  $integrationExplicit = Invoke-Harness @('-Command', 'integration', '-PreflightOnly', '-DatabaseUrl', 'postgres://explicit-secret.example/mpc_test_f02')
  if ($integrationExplicit.ExitCode -eq 0) { throw 'integration must remain blocked until F-03 supplies lifecycle' }
  Assert-Contains $integrationExplicit.Output 'target=ephemeral-postgres' 'integration target type'
  Assert-Contains $integrationExplicit.Output 'F-03' 'integration lifecycle blocker'
  Assert-NotContains $integrationExplicit.Output 'explicit-secret.example' 'integration URL redaction'
} finally {
  foreach ($name in $integrationEnvironment.Keys) { [Environment]::SetEnvironmentVariable($name, $integrationEnvironment[$name], 'Process') }
}

$provider = Invoke-Harness @('-Command', 'provider-write', '-Provider', 'mercado_livre')
if ($provider.ExitCode -eq 0) { throw 'provider-write must reject missing actor/idempotency' }
Assert-Contains $provider.Output 'actor' 'provider actor guard'
Assert-Contains $provider.Output 'idempotency' 'provider idempotency guard'
Assert-NotContains $provider.Output 'http' 'provider write before network'

$providerReady = Invoke-Harness @('-Command', 'provider-write', '-Provider', 'mercado_livre', '-Actor', 'operator-1', '-IdempotencyKey', 'idem-1')
if ($providerReady.ExitCode -eq 0) { throw 'provider-write must remain blocked without explicit execute' }
Assert-Contains $providerReady.Output 'explicit -Execute' 'provider explicit execute guard'
Assert-NotContains $providerReady.Output 'http' 'provider ready no network'

$providerExecute = Invoke-Harness @('-Command', 'provider-write', '-Provider', 'mercado_livre', '-Actor', 'operator-1', '-IdempotencyKey', 'idem-1', '-Execute')
if ($providerExecute.ExitCode -eq 0) { throw 'provider-write execute must be rejected in F-02' }
Assert-Contains $providerExecute.Output 'outside F-02' 'provider execute scope blocker'
Assert-NotContains $providerExecute.Output 'http' 'provider execute no network'

$package = Get-Content -Raw (Join-Path $repoRoot 'package.json') | ConvertFrom-Json
foreach ($name in @('harness:unit', 'harness:integration', 'harness:live', 'harness:browser', 'harness:provider-write')) {
  if (-not $package.scripts.PSObject.Properties.Name.Contains($name)) { throw "missing npm alias $name" }
}
$npmUnit = & npm run harness:unit -- -PreflightOnly 2>&1
if ($LASTEXITCODE -ne 0) { throw "npm harness:unit alias failed: $($npmUnit -join "`n")" }
Assert-Contains ($npmUnit -join "`n") 'target=fake' 'npm unit alias behavior'

Write-Output 'PASS hermetic lane tests'
