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

$previousDatabaseUrl = $env:MC_DATABASE_URL
try {
  $env:MC_DATABASE_URL = 'postgres://secret.example/db'
  $unitBlocked = Invoke-Harness @('-Command', 'unit', '-PreflightOnly')
  if ($unitBlocked.ExitCode -eq 0) { throw 'unit preflight must reject inherited live configuration' }
} finally {
  $env:MC_DATABASE_URL = $previousDatabaseUrl
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
  Assert-Contains $missing.Output 'missing_keys=' 'missing key names'
  Assert-NotContains $missing.Output 'secret' 'missing output redaction'
} finally {
  Remove-Item -LiteralPath $tempEnv -Force -ErrorAction SilentlyContinue
}

$provider = Invoke-Harness @('-Command', 'provider-write', '-Provider', 'mercado_livre')
if ($provider.ExitCode -eq 0) { throw 'provider-write must reject missing actor/idempotency' }
Assert-Contains $provider.Output 'actor' 'provider actor guard'
Assert-Contains $provider.Output 'idempotency' 'provider idempotency guard'
Assert-NotContains $provider.Output 'http' 'provider write before network'

$package = Get-Content -Raw (Join-Path $repoRoot 'package.json') | ConvertFrom-Json
foreach ($name in @('harness:unit', 'harness:integration', 'harness:live', 'harness:browser', 'harness:provider-write')) {
  if (-not $package.scripts.PSObject.Properties.Name.Contains($name)) { throw "missing npm alias $name" }
}

Write-Output 'PASS hermetic lane tests'
