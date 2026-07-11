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

$safeKeys = @('SystemRoot', 'WINDIR', 'ComSpec', 'PATH', 'PATHEXT', 'TEMP', 'TMP', 'GOCACHE', 'GOMODCACHE', 'GOPROXY', 'GOSUMDB')
$contaminatedKeys = @(
  'MC_DATABASE_URL', 'MS_DATABASE_URL', 'MPC_TEST_DATABASE_URL', 'MPC_PRODUCT_LINKS_POSTGRES_URL',
  'DATABASE_URL', 'PGHOST', 'MS_TENANT_ID', 'MPC_PROVIDER_FOO', 'MPC_ORACLE_PASSWORD',
  'SANKHYA_ORACLE_USER', 'MPC_OAUTH_REDIRECT_URI', 'MPC_WEB_PROXY_TARGET',
  'ME_CLIENT_ID', 'ME_CLIENT_SECRET', 'ME_REDIRECT_URI', 'HTTP_PROXY', 'HTTPS_PROXY',
  'NGROK_AUTHTOKEN', 'RUN_MIGRATIONS', 'MC_MIGRATIONS_DIR', 'GOCACHE', 'GOMODCACHE', 'GOPROXY', 'GOSUMDB',
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
  foreach ($key in $safeKeys | Where-Object { $_ -notin @('GOCACHE', 'GOMODCACHE', 'GOPROXY', 'GOSUMDB') }) {
    $parentValue = [Environment]::GetEnvironmentVariable($key, 'Process')
    if (-not [string]::IsNullOrWhiteSpace($parentValue)) {
      Assert-True ($child.ContainsKey($key) -and $child[$key] -eq $parentValue) "safe tool key missing or changed: $key"
    }
  }
  foreach ($key in $contaminatedKeys) {
    if ($key -in @('GOCACHE', 'GOMODCACHE', 'GOPROXY', 'GOSUMDB')) { continue }
    Assert-True (-not $child.ContainsKey($key)) "unit child inherited forbidden key $key"
  }
  $expectedCache = [IO.Path]::GetFullPath((Join-Path $repoRoot 'apps/server_core/.gocache'))
  Assert-True ([IO.Path]::GetFullPath([string]$child['GOCACHE']) -eq $expectedCache) 'unit child GOCACHE is not canonical and repository-local'
  $expectedModuleCache = [IO.Path]::GetFullPath((Join-Path $repoRoot 'apps/server_core/.gomodcache'))
  Assert-True ([IO.Path]::GetFullPath([string]$child['GOMODCACHE']) -eq $expectedModuleCache) 'unit child GOMODCACHE is not canonical and repository-local'
  Assert-True ($child['GOMODCACHE'] -ne $child['GOCACHE']) 'unit Go caches must use distinct paths'
  Assert-True ($child['GOPROXY'] -eq 'off' -and $child['GOSUMDB'] -eq 'off') 'unit Go dependency network is not disabled'

  $relativeRootError = ''
  try { New-HarnessChildEnvironment -RepositoryRoot '.' -LaneId 'unit' | Out-Null }
  catch { $relativeRootError = $_.Exception.Message }
  Assert-True ($relativeRootError -match 'HENV_REPOSITORY_ROOT_INVALID') 'relative repository root was accepted'

  $foreignLocation = [IO.Path]::GetTempPath()
  Set-Location -LiteralPath $foreignLocation
  try {
    $foreignChild = New-HarnessChildEnvironment -RepositoryRoot $repoRoot -LaneId 'unit'
    Assert-True ($foreignChild['GOCACHE'] -eq $expectedCache) 'absolute repository root depends on caller CWD'
    Assert-True ($foreignChild['GOMODCACHE'] -eq $expectedModuleCache) 'module cache depends on caller CWD'
  } finally { Set-Location -LiteralPath $originalLocation }

  foreach ($runtimeCase in @(
    @{ Values = @{ HARNESS_UNKNOWN_EXPLICIT = 'value' }; Label = 'unknown' },
    @{ Values = @{ MPC_PROVIDER_MERCADOLIVRE_CLIENT_ID = 'value' }; Label = 'lane-forbidden' }
  )) {
    $runtimeError = ''
    try { New-HarnessChildEnvironment -RepositoryRoot $repoRoot -LaneId 'unit' -RuntimeValues $runtimeCase.Values | Out-Null }
    catch { $runtimeError = $_.Exception.Message }
    Assert-True ($runtimeError -match 'HENV_RUNTIME_KEY_FORBIDDEN') "explicit $($runtimeCase.Label) runtime value was not rejected"
  }

  $precedenceKeys = @('MPC_ORACLE_USERNAME', 'MPC_ORACLE_CONNECT_STRING')
  $precedenceBefore = @{}
  try {
    foreach ($key in $precedenceKeys) { $precedenceBefore[$key] = [Environment]::GetEnvironmentVariable($key, 'Process') }
    [Environment]::SetEnvironmentVariable('MPC_ORACLE_USERNAME', 'process-user', 'Process')
    [Environment]::SetEnvironmentVariable('MPC_ORACLE_CONNECT_STRING', 'process-connect', 'Process')
    $processWins = New-HarnessChildEnvironment -RepositoryRoot $repoRoot -LaneId 'live-oracle' -EnvFile $envFile
    Assert-True ($processWins['MPC_ORACLE_USERNAME'] -eq 'process-user') 'process runtime value did not override EnvFile'
    Assert-True ($processWins['MPC_ORACLE_CONNECT_STRING'] -eq 'process-connect') 'process canonical value did not override EnvFile alias'
    Assert-True (-not $processWins.ContainsKey('MPC_PRODUCT_LINKS_POSTGRES_URL')) 'lane-forbidden process value entered live Oracle child'

    $explicitWins = New-HarnessChildEnvironment -RepositoryRoot $repoRoot -LaneId 'live-oracle' -EnvFile $envFile -RuntimeValues @{
      MPC_ORACLE_USERNAME = 'explicit-user'
      SANKHYA_ORACLE_CONNECT_STRING = 'explicit-legacy-connect'
    }
    Assert-True ($explicitWins['MPC_ORACLE_USERNAME'] -eq 'explicit-user') 'explicit runtime value did not override process value'
    Assert-True ($explicitWins['MPC_ORACLE_CONNECT_STRING'] -eq 'explicit-legacy-connect') 'explicit alias did not override process canonical value'
  }
  finally {
    foreach ($key in $precedenceBefore.Keys) { [Environment]::SetEnvironmentVariable($key, $precedenceBefore[$key], 'Process') }
  }

  $fromTolerantFile = New-HarnessChildEnvironment -RepositoryRoot $repoRoot -LaneId 'live-oracle' -EnvFile $envFile
  Assert-True ($fromTolerantFile['MPC_ORACLE_CONNECT_STRING'] -eq 'fixture-oracle-connect') 'EnvFile aliases did not normalize'
  Assert-True (-not $fromTolerantFile.ContainsKey('MPC_PRODUCT_LINKS_POSTGRES_URL')) 'EnvFile lane-forbidden key entered child'

  $canonicalWins = New-HarnessChildEnvironment -RepositoryRoot $repoRoot -LaneId 'live-oracle' -RuntimeValues @{
    MPC_ORACLE_CONNECT_STRING = 'canonical-connect'
    SANKHYA_ORACLE_CONNECT_STRING = 'legacy-connect'
    SANKHYA_ORACLE_HOST = 'legacy-host'
    SANKHYA_ORACLE_PORT = '1521'
    SANKHYA_ORACLE_SERVICE_NAME = 'legacy-service'
  }
  Assert-True ($canonicalWins['MPC_ORACLE_CONNECT_STRING'] -eq 'canonical-connect') 'canonical Oracle connect string lost precedence'

  $legacyConnect = New-HarnessChildEnvironment -RepositoryRoot $repoRoot -LaneId 'live-oracle' -RuntimeValues @{
    SANKHYA_ORACLE_CONNECT_STRING = 'legacy-connect'
    SANKHYA_ORACLE_HOST = 'legacy-host'
    SANKHYA_ORACLE_PORT = '1521'
    SANKHYA_ORACLE_SERVICE_NAME = 'legacy-service'
  }
  Assert-True ($legacyConnect['MPC_ORACLE_CONNECT_STRING'] -eq 'legacy-connect') 'direct legacy Oracle connect string lost precedence'
  Assert-True (-not $legacyConnect.ContainsKey('SANKHYA_ORACLE_CONNECT_STRING')) 'legacy Oracle alias leaked into child'

  $tupleConnect = New-HarnessChildEnvironment -RepositoryRoot $repoRoot -LaneId 'live-oracle' -RuntimeValues @{
    SANKHYA_ORACLE_HOST = 'tuple-host'
    SANKHYA_ORACLE_PORT = '1522'
    SANKHYA_ORACLE_SERVICE_NAME = 'tuple-service'
  }
  Assert-True ($tupleConnect['MPC_ORACLE_CONNECT_STRING'] -eq 'tuple-host:1522/tuple-service') 'complete Oracle tuple did not normalize'

  $incompleteTuple = New-HarnessChildEnvironment -RepositoryRoot $repoRoot -LaneId 'live-oracle' -RuntimeValues @{
    SANKHYA_ORACLE_HOST = 'host-only'
  }
  Assert-True (-not $incompleteTuple.ContainsKey('MPC_ORACLE_CONNECT_STRING')) 'incomplete Oracle tuple fabricated a connect string'

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
