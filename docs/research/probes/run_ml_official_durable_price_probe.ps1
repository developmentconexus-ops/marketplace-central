$ErrorActionPreference = 'Stop'
$repo = 'C:\Users\leandro.theodoro\Documents\marketplace-central'
$packageDir = Join-Path $repo 'apps\server_core\internal\modules\connectors\adapters\mercado_livre'
$probeSource = Join-Path $repo 'docs\research\probes\ml_official_durable_price_probe_test.go'
$probeTarget = Join-Path $packageDir 'ml_official_durable_price_probe_test.go'
$overlayPath = Join-Path $env:TEMP 'ml_official_durable_price_overlay.json'
$inputPath = Join-Path $repo 'docs\research\evidence\2026-07-13-ml-official-durable-price-input.json'
$outputPath = if ($env:ML_DURABLE_OUTPUT) { $env:ML_DURABLE_OUTPUT } else { 'C:\tmp\ml_official_durable_price_result.json' }

. (Join-Path $repo 'scripts\lib\dev-local-runtime.ps1')
Import-DevLocalEnvironment $repo
$env:MC_DATABASE_URL = 'postgres://marketplace:marketplace@127.0.0.1:5435/marketplace_central?sslmode=disable'
$env:GOCACHE = Join-Path $repo '.gocache'
$env:ML_DURABLE_INPUT = $inputPath
$env:ML_DURABLE_OUTPUT = $outputPath

$overlay = [ordered]@{ Replace = [ordered]@{ $probeTarget = $probeSource } }
$overlay | ConvertTo-Json -Depth 4 | Set-Content -LiteralPath $overlayPath -Encoding UTF8
$postgresWasRunning = -not [string]::IsNullOrWhiteSpace((& docker compose ps --status running -q postgres))

try {
  if (-not $postgresWasRunning) {
    & docker compose up -d postgres
    if ($LASTEXITCODE -ne 0) { throw 'failed to start postgres' }
  }
  Push-Location (Join-Path $repo 'apps\server_core')
  try {
    & go test -overlay $overlayPath -run TestMLOfficialDurablePriceProbe -count=1 -v ./internal/modules/connectors/adapters/mercado_livre
    $testExitCode = $LASTEXITCODE
  } finally {
    Pop-Location
  }
  exit $testExitCode
} finally {
  if (-not $postgresWasRunning) {
    & docker compose stop postgres | Out-Host
  }
  Remove-Item -LiteralPath $overlayPath -Force -ErrorAction SilentlyContinue
}
