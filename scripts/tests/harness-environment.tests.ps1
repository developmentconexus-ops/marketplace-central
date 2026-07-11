$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$module = Join-Path $repoRoot 'scripts/harness/Environment.psm1'

function Assert-True([bool]$Condition, [string]$Message) {
  if (-not $Condition) { throw $Message }
}

if (-not (Test-Path -LiteralPath $module -PathType Leaf)) {
  throw 'RED: Environment.psm1 missing; fresh child environment contract is not implemented'
}
Import-Module $module -Force

$safeKeys = @('SystemRoot', 'WINDIR', 'ComSpec', 'PATH', 'PATHEXT', 'TEMP', 'TMP', 'GOCACHE')
$contaminatedKeys = @(
  'MC_DATABASE_URL', 'MS_DATABASE_URL', 'MPC_TEST_DATABASE_URL', 'MPC_PRODUCT_LINKS_POSTGRES_URL',
  'DATABASE_URL', 'PGHOST', 'MS_TENANT_ID', 'MPC_PROVIDER_FOO', 'MPC_ORACLE_PASSWORD',
  'SANKHYA_ORACLE_USER', 'MPC_OAUTH_REDIRECT_URI', 'MPC_WEB_PROXY_TARGET',
  'ME_CLIENT_ID', 'ME_CLIENT_SECRET', 'ME_REDIRECT_URI', 'HTTP_PROXY', 'HTTPS_PROXY',
  'NGROK_AUTHTOKEN', 'RUN_MIGRATIONS', 'MC_MIGRATIONS_DIR', 'GOCACHE',
  ('HARNESS_UNKNOWN_' + [guid]::NewGuid().ToString('N'))
)
$before = @{}
$originalLocation = (Get-Location).Path
foreach ($key in $contaminatedKeys) {
  $before[$key] = [Environment]::GetEnvironmentVariable($key, 'Process')
  [Environment]::SetEnvironmentVariable($key, "fixture-$key", 'Process')
}

$envFile = Join-Path ([IO.Path]::GetTempPath()) ('harness-unit-' + [guid]::NewGuid().ToString('N') + '.env')
try {
  @(
    'MPC_PRODUCT_LINKS_POSTGRES_URL=fixture-db',
    'ME_CLIENT_SECRET=fixture-provider',
    'SANKHYA_ORACLE_USER=fixture-oracle-user',
    'SANKHYA_ORACLE_PASSWORD=fixture-oracle-password',
    'SANKHYA_ORACLE_CONNECT_STRING=fixture-oracle-connect',
    'SANKHYA_ORACLE_HOST=fixture-oracle-host',
    'SANKHYA_ORACLE_PORT=1521',
    'SANKHYA_ORACLE_SERVICE_NAME=fixture-oracle-service'
  ) | Set-Content -LiteralPath $envFile -Encoding utf8

  $child = New-HarnessChildEnvironment -RepositoryRoot $repoRoot -LaneId 'unit' -EnvFile $envFile
  $keys = @($child.Keys)
  Assert-True (@($keys | Where-Object { $_ -notin $safeKeys }).Count -eq 0) 'unit child contains a non-safe key'
  foreach ($key in $safeKeys | Where-Object { $_ -ne 'GOCACHE' }) {
    $parentValue = [Environment]::GetEnvironmentVariable($key, 'Process')
    if (-not [string]::IsNullOrWhiteSpace($parentValue)) {
      Assert-True ($child.ContainsKey($key) -and $child[$key] -eq $parentValue) "safe tool key missing or changed: $key"
    }
  }
  foreach ($key in $contaminatedKeys) {
    if ($key -eq 'GOCACHE') { continue }
    Assert-True (-not $child.ContainsKey($key)) "unit child inherited forbidden key $key"
  }
  $expectedCache = [IO.Path]::GetFullPath((Join-Path $repoRoot 'apps/server_core/.gocache'))
  Assert-True ([IO.Path]::GetFullPath([string]$child['GOCACHE']) -eq $expectedCache) 'unit child GOCACHE is not canonical and repository-local'

  foreach ($key in $contaminatedKeys) {
    Assert-True ([Environment]::GetEnvironmentVariable($key, 'Process') -eq "fixture-$key") "parent environment mutated for $key"
  }
  Assert-True ((Get-Location).Path -eq $originalLocation) 'parent location mutated'

  $unsupported = $null
  $marker = Join-Path ([IO.Path]::GetTempPath()) ('harness-unsupported-' + [guid]::NewGuid().ToString('N'))
  $runRoot = Join-Path $repoRoot 'scripts/.runs'
  $runCountBefore = @(Get-ChildItem -LiteralPath $runRoot -Directory -ErrorAction SilentlyContinue).Count
  try {
    New-HarnessChildEnvironment -RepositoryRoot $repoRoot -LaneId 'not-a-lane' | Out-Null
    Set-Content -LiteralPath $marker -Value 'started'
  } catch { $unsupported = $_.Exception.Message }
  Assert-True ($unsupported -match 'HENV_LANE_UNSUPPORTED') 'unsupported lane lacks stable reason code'
  Assert-True (-not (Test-Path -LiteralPath $marker)) 'unsupported lane reached process-start marker'
  Assert-True (@(Get-ChildItem -LiteralPath $runRoot -Directory -ErrorAction SilentlyContinue).Count -eq $runCountBefore) 'unsupported lane created run directory'
} finally {
  foreach ($key in $before.Keys) { [Environment]::SetEnvironmentVariable($key, $before[$key], 'Process') }
  Remove-Item -LiteralPath $envFile -Force -ErrorAction SilentlyContinue
  Set-Location -LiteralPath $originalLocation
}

Write-Output 'PASS harness environment tests'
