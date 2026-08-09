$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$harness = Join-Path $repoRoot 'scripts/harness.ps1'

function Invoke-Harness([string[]]$Arguments) {
  $output = & pwsh -NoProfile -ExecutionPolicy Bypass -File $harness @Arguments 2>&1
  [pscustomobject]@{ ExitCode = $LASTEXITCODE; Output = ($output -join "`n") }
}

function Assert-NotContains([string]$Value, [string]$Needle, [string]$Name) {
  if ($Value -match [regex]::Escape($Needle)) { throw "$Name leaked forbidden value" }
}

$unknownKey = 'HARNESS_UNKNOWN_' + [guid]::NewGuid().ToString('N')
$keys = @(
  'MC_DATABASE_URL', 'MPC_TEST_DATABASE_URL', 'MPC_UNKNOWN_DATABASE_TARGET',
  'MPC_PROVIDER_FOO', 'MPC_ORACLE_PASSWORD', 'SANKHYA_ORACLE_USER',
  'ME_CLIENT_ID', 'ME_CLIENT_SECRET', 'ME_REDIRECT_URI', 'HTTP_PROXY', 'GOCACHE', $unknownKey
)
$before = @{}
$envFile = Join-Path ([IO.Path]::GetTempPath()) ('harness-lane-' + [guid]::NewGuid().ToString('N') + '.env')
try {
  foreach ($key in $keys) {
    $before[$key] = [Environment]::GetEnvironmentVariable($key, 'Process')
    [Environment]::SetEnvironmentVariable($key, "fixture-$key", 'Process')
  }
  @(
    'MPC_UNKNOWN_DATABASE_TARGET=fixture-file-db',
    'ME_CLIENT_SECRET=fixture-file-provider',
    'SANKHYA_ORACLE_PASSWORD=fixture-file-oracle'
  ) | Set-Content -LiteralPath $envFile -Encoding utf8

  $unit = Invoke-Harness @('-Command', 'unit', '-PreflightOnly', '-EnvFile', $envFile)
  if ($unit.ExitCode -ne 0) { throw 'unit must ignore ambient and EnvFile runtime configuration' }
  if ($unit.Output -notmatch 'target=fake' -or $unit.Output -notmatch 'env=ignored') { throw 'unit classification missing' }
  foreach ($key in $keys) { Assert-NotContains $unit.Output "fixture-$key" "unit $key" }
  foreach ($value in @('fixture-file-db', 'fixture-file-provider', 'fixture-file-oracle')) { Assert-NotContains $unit.Output $value 'unit EnvFile' }
  $runMatch = [regex]::Match($unit.Output, 'run_dir=(scripts/\.runs/[a-f0-9]+)')
  if (-not $runMatch.Success) { throw 'unit did not report minimal run directory' }
  $summary = Get-Content -Raw -LiteralPath (Join-Path (Join-Path $repoRoot $runMatch.Groups[1].Value) 'summary.txt')
  foreach ($key in $keys) { Assert-NotContains $summary "fixture-$key" "unit summary $key" }
  foreach ($value in @('fixture-file-db', 'fixture-file-provider', 'fixture-file-oracle')) { Assert-NotContains $summary $value 'unit summary EnvFile' }

  foreach ($key in $keys) {
    if ([Environment]::GetEnvironmentVariable($key, 'Process') -ne "fixture-$key") { throw "parent environment mutated for $key" }
  }

  $integration = Invoke-Harness @('-Command', 'integration', '-PreflightOnly')
  if ($integration.ExitCode -ne 0 -or $integration.Output -notmatch 'target=ephemeral-postgres' -or $integration.Output -notmatch 'status=ready') { throw 'integration preflight must be contact-free and ready' }
  $browser = Invoke-Harness @('-Command', 'browser', '-PreflightOnly')
  if ($browser.ExitCode -ne 0 -or $browser.Output -notmatch 'target=browser') { throw 'browser preflight classification changed' }
  # The guard now validates actor/idempotency arguments before it ever reaches
  # the network-rejection message, so the structural signal to pin is the
  # blocked status plus a nonzero exit -- both still prove no network was
  # touched, whatever the first refusal happens to be.
  $provider = Invoke-Harness @('-Command', 'provider-write', '-Provider', 'mercado_livre')
  if ($provider.ExitCode -eq 0 -or $provider.Output -notmatch 'status=blocked') { throw 'provider write no-network guard changed' }
} finally {
  foreach ($key in $before.Keys) { [Environment]::SetEnvironmentVariable($key, $before[$key], 'Process') }
  Remove-Item -LiteralPath $envFile -Force -ErrorAction SilentlyContinue
}

$legacyKeys = @(
  'MPC_ORACLE_USERNAME', 'MPC_ORACLE_PASSWORD', 'MPC_ORACLE_CONNECT_STRING',
  'SANKHYA_ORACLE_USER', 'SANKHYA_ORACLE_PASSWORD', 'SANKHYA_ORACLE_CONNECT_STRING'
)
$legacyBefore = @{}
$legacyFile = Join-Path ([IO.Path]::GetTempPath()) ('harness-legacy-' + [guid]::NewGuid().ToString('N') + '.env')
try {
  foreach ($key in $legacyKeys) {
    $legacyBefore[$key] = [Environment]::GetEnvironmentVariable($key, 'Process')
    [Environment]::SetEnvironmentVariable($key, $null, 'Process')
  }
  @(
    'SANKHYA_ORACLE_USER=fixture-legacy-user',
    'SANKHYA_ORACLE_PASSWORD=fixture-legacy-password',
    'SANKHYA_ORACLE_CONNECT_STRING=fixture-legacy-connect'
  ) | Set-Content -LiteralPath $legacyFile -Encoding utf8
  # Alias normalization became silent (Environment.psm1 resolves alias_for
  # without per-key telemetry), so the old key= message pin has nothing to
  # match. The behavioral proof is stronger anyway: preflight requires the
  # canonical trio, so reaching status=ready from a file that carries ONLY the
  # legacy aliases is the normalization working.
  $legacy = Invoke-Harness @('-Command', 'live', '-Target', 'oracle', '-PreflightOnly', '-EnvFile', $legacyFile)
  if ($legacy.ExitCode -ne 0) { throw 'legacy EnvFile aliases must satisfy live Oracle preflight' }
  if ($legacy.Output -notmatch 'status=ready') { throw 'legacy aliases did not reach a ready preflight' }
  foreach ($value in @('fixture-legacy-user', 'fixture-legacy-password', 'fixture-legacy-connect')) {
    Assert-NotContains $legacy.Output $value 'legacy alias output'
  }
} finally {
  foreach ($key in $legacyBefore.Keys) { [Environment]::SetEnvironmentVariable($key, $legacyBefore[$key], 'Process') }
  Remove-Item -LiteralPath $legacyFile -Force -ErrorAction SilentlyContinue
}

Write-Output 'PASS hermetic lane tests'
