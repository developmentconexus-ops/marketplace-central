$ErrorActionPreference = 'Stop'

$repositoryRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$policyPath = Join-Path $repositoryRoot 'scripts/harness/Policy.psm1'
if (-not (Test-Path -LiteralPath $policyPath -PathType Leaf)) {
  Write-Error 'Policy.psm1 is missing; RED governance drift fixture cannot run.'
  exit 1
}

Import-Module $policyPath -Force

function Assert-True {
  param([bool]$Condition, [string]$Message)
  if (-not $Condition) { throw $Message }
}

function Write-FixtureFile {
  param([string]$Root, [string]$Path, [string]$Content)
  $fullPath = Join-Path $Root ($Path -replace '/', [IO.Path]::DirectorySeparatorChar)
  $parent = Split-Path -Parent $fullPath
  if ($parent) { New-Item -ItemType Directory -Path $parent -Force | Out-Null }
  Set-Content -LiteralPath $fullPath -Value $Content -Encoding utf8NoBOM
}

function New-PositiveFixture {
  $root = Join-Path ([IO.Path]::GetTempPath()) ("mpc-governance-{0}" -f [guid]::NewGuid().ToString('N'))
  New-Item -ItemType Directory -Path $root -Force | Out-Null
  Copy-Item -LiteralPath (Join-Path $repositoryRoot 'contracts') -Destination $root -Recurse

  $modules = Get-Content -Raw (Join-Path $root 'contracts/governance/modules.json') | ConvertFrom-Json
  foreach ($module in $modules.modules) {
    Write-FixtureFile $root "$($module.root)/domain/doc.go" "package domain`n"
  }
  $composition = ($modules.modules | ForEach-Object { "import _ `"marketplace-central/apps/server_core/internal/modules/$($_.id)`"" }) -join "`n"
  Write-FixtureFile $root 'apps/server_core/internal/composition/root.go' "package composition`n$composition`n"

  $runtime = Get-Content -Raw (Join-Path $root 'contracts/governance/runtime-config.json') | ConvertFrom-Json
  $readerKeys = @{}
  foreach ($entry in $runtime.keys) {
    foreach ($reader in $entry.readers) {
      if (-not $readerKeys.ContainsKey($reader.path)) { $readerKeys[$reader.path] = [Collections.Generic.List[string]]::new() }
      $readerKeys[$reader.path].Add([string]$entry.key)
    }
  }
  foreach ($readerPath in $readerKeys.Keys) {
    $keys = @($readerKeys[$readerPath] | Sort-Object -Unique)
    if ($readerPath -like '*.go') {
      $reads = ($keys | ForEach-Object { "var _ = os.Getenv(`"$_`")" }) -join "`n"
      Write-FixtureFile $root $readerPath "package fixture`nimport `"os`"`n$reads`n"
    } elseif ($readerPath -like '*.ts' -or $readerPath -like '*.tsx') {
      $reads = ($keys | ForEach-Object { "const $_ = process.env.$_;" }) -join "`n"
      Write-FixtureFile $root $readerPath $reads
    } elseif ($readerPath -like '*.ps1') {
      $reads = ($keys | ForEach-Object { "[Environment]::GetEnvironmentVariable('$_') | Out-Null" }) -join "`n"
      Write-FixtureFile $root $readerPath $reads
    } else {
      $reads = ($keys | ForEach-Object { 'echo ${' + $_ + ':-}' }) -join "`n"
      Write-FixtureFile $root $readerPath $reads
    }
  }

  Write-FixtureFile $root 'apps/server_core/internal/modules/product_links/application/resolution_service.go' @'
package application
func fixture() { panic("unsupported product link transition input") }
'@
  Write-FixtureFile $root 'apps/server_core/migrations/0021_integration_operation_run_evidence.sql' '-- fixture'
  Write-FixtureFile $root 'apps/server_core/migrations/0021_integrations_provider_auth_strategy_shopee_partner.sql' '-- fixture'
  Write-FixtureFile $root 'contracts/api/marketplace-central.openapi.yaml' 'openapi: 3.1.0'
  Write-FixtureFile $root 'packages/sdk-runtime/src/index.ts' 'export {}'
  Write-FixtureFile $root 'apps/web/src/index.ts' 'export {}'
  return $root
}

function Assert-FailureCode {
  param([string]$Root, [string]$ExpectedCode, [string]$BaseSha = '')
  $result = Test-GovernanceDrift -RepositoryRoot $Root -BaseSha $BaseSha
  Assert-True (-not $result.Passed) "expected governance drift failure $ExpectedCode"
  Assert-True ($ExpectedCode -in @($result.Violations.ErrorCode)) "missing reason code $ExpectedCode; actual=$(@($result.Violations.ErrorCode) -join ',')"
}

$fixtures = [Collections.Generic.List[string]]::new()
try {
  $positive = New-PositiveFixture; $fixtures.Add($positive)
  $positiveResult = Test-GovernanceDrift -RepositoryRoot $positive
  Assert-True $positiveResult.Passed "positive fixture failed: $(@($positiveResult.Violations.ErrorCode) -join ',')"
  Assert-True ('production-panic-product-link-transition' -in @($positiveResult.BaselineExceptions.Id)) 'positive fixture did not report panic baseline exception'
  Assert-True ('migration-prefix-0021-duplicate' -in @($positiveResult.BaselineExceptions.Id)) 'positive fixture did not report migration baseline exception'

  $missingModule = New-PositiveFixture; $fixtures.Add($missingModule)
  Remove-Item -LiteralPath (Join-Path $missingModule 'apps/server_core/internal/modules/catalog') -Recurse -Force
  Assert-FailureCode $missingModule 'GOV_MODULE_COVERAGE'

  $undeclaredDependency = New-PositiveFixture; $fixtures.Add($undeclaredDependency)
  Write-FixtureFile $undeclaredDependency 'apps/server_core/internal/modules/catalog/application/usecase.go' 'package application; import _ "marketplace-central/apps/server_core/internal/modules/pricing/domain"'
  Assert-FailureCode $undeclaredDependency 'GOV_MODULE_DEPENDENCY'

  $applicationImport = New-PositiveFixture; $fixtures.Add($applicationImport)
  Write-FixtureFile $applicationImport 'apps/server_core/internal/modules/catalog/application/http.go' 'package application; import _ "net/http"'
  Assert-FailureCode $applicationImport 'GOV_APPLICATION_IMPORT'

  $undeclaredReader = New-PositiveFixture; $fixtures.Add($undeclaredReader)
  Write-FixtureFile $undeclaredReader 'apps/server_core/internal/platform/config/rogue.go' 'package config; import "os"; var _ = os.Getenv("MPC_ROGUE_SECRET")'
  Assert-FailureCode $undeclaredReader 'RCFG_UNDECLARED_READ'

  $dynamicReader = New-PositiveFixture; $fixtures.Add($dynamicReader)
  Write-FixtureFile $dynamicReader 'apps/server_core/internal/platform/config/dynamic.go' 'package config; import "os"; func read(key string) string { return os.Getenv(key) }'
  Assert-FailureCode $dynamicReader 'RCFG_DYNAMIC_READER_UNBOUNDED'

  $aliasCollision = New-PositiveFixture; $fixtures.Add($aliasCollision)
  $runtimePath = Join-Path $aliasCollision 'contracts/governance/runtime-config.json'
  $runtime = Get-Content -Raw $runtimePath | ConvertFrom-Json
  $duplicate = $runtime.keys | Where-Object key -eq 'SANKHYA_ORACLE_USER' | Select-Object -First 1
  $duplicate = $duplicate.PSObject.Copy(); $duplicate.alias_for = 'MPC_ORACLE_CONNECT_STRING'
  $runtime.keys += $duplicate
  $runtime | ConvertTo-Json -Depth 20 | Set-Content -LiteralPath $runtimePath -Encoding utf8NoBOM
  Assert-FailureCode $aliasCollision 'RCFG_ALIAS_COLLISION'

  $secretClassification = New-PositiveFixture; $fixtures.Add($secretClassification)
  $runtimePath = Join-Path $secretClassification 'contracts/governance/runtime-config.json'
  $runtime = Get-Content -Raw $runtimePath | ConvertFrom-Json
  ($runtime.keys | Where-Object key -eq 'MPC_ENCRYPTION_KEY').sensitivity = 'public'
  $runtime | ConvertTo-Json -Depth 20 | Set-Content -LiteralPath $runtimePath -Encoding utf8NoBOM
  Assert-FailureCode $secretClassification 'RCFG_SECRET_CLASS_MISMATCH'

  $panicFixture = New-PositiveFixture; $fixtures.Add($panicFixture)
  Write-FixtureFile $panicFixture 'apps/server_core/internal/modules/catalog/domain/panic.go' 'package domain; func bad() { panic("new panic") }'
  Assert-FailureCode $panicFixture 'GOV_PRODUCTION_PANIC'

  $migrationFixture = New-PositiveFixture; $fixtures.Add($migrationFixture)
  Write-FixtureFile $migrationFixture 'apps/server_core/migrations/0099_first.sql' '-- fixture'
  Write-FixtureFile $migrationFixture 'apps/server_core/migrations/0099_second.sql' '-- fixture'
  Assert-FailureCode $migrationFixture 'GOV_MIGRATION_PREFIX'

  $frontendFixture = New-PositiveFixture; $fixtures.Add($frontendFixture)
  Write-FixtureFile $frontendFixture 'apps/web/src/rogue.ts' 'export const load = () => fetch("/api")'
  Assert-FailureCode $frontendFixture 'GOV_FRONTEND_FETCH'

  $atomicFixture = New-PositiveFixture; $fixtures.Add($atomicFixture)
  & git -C $atomicFixture init --quiet
  & git -C $atomicFixture config core.autocrlf false
  & git -C $atomicFixture config user.email 'fixture@example.invalid'
  & git -C $atomicFixture config user.name 'Fixture'
  & git -C $atomicFixture add .
  & git -C $atomicFixture commit --quiet -m baseline
  $baseSha = (& git -C $atomicFixture rev-parse HEAD).Trim()
  Add-Content -LiteralPath (Join-Path $atomicFixture 'contracts/api/marketplace-central.openapi.yaml') -Value '# changed'
  Assert-FailureCode $atomicFixture 'GOV_API_SDK_SPLIT' $baseSha

  $invariants = Get-Content -Raw (Join-Path $positive 'contracts/governance/invariants.json') | ConvertFrom-Json
  Assert-True ('frontend-fetch' -in @((Get-ApplicableInvariants -Registry $invariants -Paths @('apps/web/src/rogue.ts')).id)) 'applicable invariant lookup missed frontend path'
  $seams = Get-Content -Raw (Join-Path $positive 'contracts/governance/shared-seams.json') | ConvertFrom-Json
  Assert-True ('api-sdk' -in @((Get-SharedSeams -Registry $seams -Paths @('packages/sdk-runtime/src/index.ts')).id)) 'shared seam lookup missed SDK path'

  Write-Output 'PASS governance drift tests'
} finally {
  foreach ($fixture in $fixtures) {
    if (Test-Path -LiteralPath $fixture) { Remove-Item -LiteralPath $fixture -Recurse -Force }
  }
}
