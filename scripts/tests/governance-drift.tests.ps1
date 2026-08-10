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
  # Written after the runtime-config reader loop below, not here: that loop also
  # owns apps/server_core/internal/composition/root.go, because four keys declare
  # it as a reader. Writing it here first meant the loop overwrote the imports and
  # every composition_required module came back as GOV_COMPOSITION_MISSING.

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
    if ($readerPath -eq 'scripts/harness/Environment.psm1') {
      Write-FixtureFile $root $readerPath @'
$script:SafeToolKeys = @('PATH')
function Get-HarnessRegistry { param($RepositoryRoot, $Name) @{} }
function Get-ProcessRuntimeValues {
  param($Lane)
  foreach ($allowedKey in @($Lane.allowed_runtime_keys)) {
    $name = [string]$allowedKey
    $value = [Environment]::GetEnvironmentVariable($name, 'Process')
  }
}
function New-HarnessChildEnvironment {
  $runtime = Get-HarnessRegistry -RepositoryRoot $root -Name 'runtime-config'
  foreach ($key in $script:SafeToolKeys) {
    $value = [Environment]::GetEnvironmentVariable($key, 'Process')
  }
}
'@
    } elseif ($readerPath -eq 'apps/server_core/internal/composition/root.go') {
      # This one file has two jobs: it is the composition root AND a declared
      # reader of four runtime keys. It has to satisfy both in one write.
      $reads = ($keys | ForEach-Object { "var _ = os.Getenv(`"$_`")" }) -join "`n"
      Write-FixtureFile $root $readerPath "package composition`nimport `"os`"`n$composition`n$reads`n"
    } elseif ($readerPath -like '*.go') {
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

  # The loop only writes the composition root while some key still declares it a
  # reader. If that ever stops being true the file must still exist, or the
  # fixture would go red for a reason that has nothing to do with the test.
  if (-not (Test-Path -LiteralPath (Join-Path $root ('apps/server_core/internal/composition/root.go' -replace '/', [IO.Path]::DirectorySeparatorChar)) -PathType Leaf)) {
    Write-FixtureFile $root 'apps/server_core/internal/composition/root.go' "package composition`n$composition`n"
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

function Test-ReviewFailureCode {
  param([Collections.Generic.List[string]]$Failures, [string]$Root, [string]$ExpectedCode, [string]$Case, [string]$BaseSha = '')
  $result = Test-GovernanceDrift -RepositoryRoot $Root -BaseSha $BaseSha
  if ($result.Passed) {
    $Failures.Add("$Case false negative: drift passed; expected $ExpectedCode")
  } elseif ($ExpectedCode -notin @($result.Violations.ErrorCode)) {
    $Failures.Add("$Case wrong code: expected $ExpectedCode; actual=$(@($result.Violations.ErrorCode) -join ',')")
  }
}

$fixtures = [Collections.Generic.List[string]]::new()
try {
  $reviewFailures = [Collections.Generic.List[string]]::new()
  $positive = New-PositiveFixture; $fixtures.Add($positive)
  $positiveResult = Test-GovernanceDrift -RepositoryRoot $positive
  Assert-True $positiveResult.Passed "positive fixture failed: $(@($positiveResult.Violations.ErrorCode) -join ',')"
  Assert-True ('production-panic-product-link-transition' -in @($positiveResult.BaselineExceptions.Id)) 'positive fixture did not report panic baseline exception'
  Assert-True ('migration-prefix-0021-duplicate' -in @($positiveResult.BaselineExceptions.Id)) 'positive fixture did not report migration baseline exception'

  $missingModule = New-PositiveFixture; $fixtures.Add($missingModule)
  Remove-Item -LiteralPath (Join-Path $missingModule 'apps/server_core/internal/modules/catalog') -Recurse -Force
  Assert-FailureCode $missingModule 'GOV_MODULE_COVERAGE'
  Write-Output 'guard_ran=gov-module-coverage'

  $undeclaredDependency = New-PositiveFixture; $fixtures.Add($undeclaredDependency)
  Write-FixtureFile $undeclaredDependency 'apps/server_core/internal/modules/catalog/application/usecase.go' 'package application; import _ "marketplace-central/apps/server_core/internal/modules/pricing/domain"'
  Assert-FailureCode $undeclaredDependency 'GOV_MODULE_DEPENDENCY'
  Write-Output 'guard_ran=gov-module-dependency'

  $applicationImport = New-PositiveFixture; $fixtures.Add($applicationImport)
  Write-FixtureFile $applicationImport 'apps/server_core/internal/modules/catalog/application/http.go' 'package application; import _ "net/http"'
  Assert-FailureCode $applicationImport 'GOV_APPLICATION_IMPORT'
  Write-Output 'guard_ran=gov-application-import'

  # GOV_COMPOSITION_MISSING: a composition_required module whose import line is
  # absent from the composition root. `classifications` is chosen over `catalog`
  # because the module-coverage fixture above already deletes catalog's tree, and
  # this fixture needs a module that stays otherwise intact.
  $compositionMissing = New-PositiveFixture; $fixtures.Add($compositionMissing)
  $compositionRootPath = Join-Path $compositionMissing 'apps/server_core/internal/composition/root.go'
  $compositionRootContent = Get-Content -Raw -LiteralPath $compositionRootPath
  $compositionRootContent = $compositionRootContent -replace `
    [regex]::Escape('import _ "marketplace-central/apps/server_core/internal/modules/classifications"' + "`n"), ''
  Assert-True ($compositionRootContent -notmatch '/modules/classifications"') 'fixture premise broke: the classifications import line was not removed'
  Set-Content -LiteralPath $compositionRootPath -Value $compositionRootContent -Encoding utf8NoBOM
  Assert-FailureCode $compositionMissing 'GOV_COMPOSITION_MISSING'
  Write-Output 'guard_ran=gov-composition-missing'

  $undeclaredReader = New-PositiveFixture; $fixtures.Add($undeclaredReader)
  Write-FixtureFile $undeclaredReader 'apps/server_core/internal/platform/config/rogue.go' 'package config; import "os"; var _ = os.Getenv("MPC_ROGUE_SECRET")'
  Assert-FailureCode $undeclaredReader 'RCFG_UNDECLARED_READ'

  $dynamicReader = New-PositiveFixture; $fixtures.Add($dynamicReader)
  Write-FixtureFile $dynamicReader 'apps/server_core/internal/platform/config/dynamic.go' 'package config; import "os"; func read(key string) string { return os.Getenv(key) }'
  Assert-FailureCode $dynamicReader 'RCFG_DYNAMIC_READER_UNBOUNDED'

  $staleRegistryReader = New-PositiveFixture; $fixtures.Add($staleRegistryReader)
  $registryReaderPath = Join-Path $staleRegistryReader 'scripts/harness/Environment.psm1'
  $registryReader = Get-Content -Raw -LiteralPath $registryReaderPath
  $registryReader = $registryReader.Replace("    `$value = [Environment]::GetEnvironmentVariable(`$name, 'Process')", "    `$value = ''")
  Set-Content -LiteralPath $registryReaderPath -Value $registryReader -Encoding utf8NoBOM
  Assert-FailureCode $staleRegistryReader 'RCFG_DYNAMIC_READER_UNBOUNDED'

  $extraRegistryReader = New-PositiveFixture; $fixtures.Add($extraRegistryReader)
  Add-Content -LiteralPath (Join-Path $extraRegistryReader 'scripts/harness/Environment.psm1') -Value "[Environment]::GetEnvironmentVariable(`$evil, 'Process') | Out-Null" -Encoding utf8NoBOM
  Assert-FailureCode $extraRegistryReader 'RCFG_DYNAMIC_READER_UNBOUNDED'

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
  Write-Output 'guard_ran=gov-production-panic'

  $migrationFixture = New-PositiveFixture; $fixtures.Add($migrationFixture)
  Write-FixtureFile $migrationFixture 'apps/server_core/migrations/0099_first.sql' '-- fixture'
  Write-FixtureFile $migrationFixture 'apps/server_core/migrations/0099_second.sql' '-- fixture'
  Assert-FailureCode $migrationFixture 'GOV_MIGRATION_PREFIX'
  Write-Output 'guard_ran=gov-migration-prefix'

  $frontendFixture = New-PositiveFixture; $fixtures.Add($frontendFixture)
  Write-FixtureFile $frontendFixture 'apps/web/src/rogue.ts' 'export const load = () => fetch("/api")'
  Assert-FailureCode $frontendFixture 'GOV_FRONTEND_FETCH'
  Write-Output 'guard_ran=gov-frontend-fetch'

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
  Write-Output 'guard_ran=gov-api-sdk-split'

  $powerShellEnvFixture = New-PositiveFixture; $fixtures.Add($powerShellEnvFixture)
  Write-FixtureFile $powerShellEnvFixture 'scripts/review-reader.ps1' '$value = $env:MPC_REVIEW_ROGUE_SECRET'
  Test-ReviewFailureCode $reviewFailures $powerShellEnvFixture 'RCFG_UNDECLARED_READ' 'PowerShell env literal'

  $injectedGetterFixture = New-PositiveFixture; $fixtures.Add($injectedGetterFixture)
  Write-FixtureFile $injectedGetterFixture 'apps/server_core/internal/platform/config/review_getter.go' 'package config; func read(getenv func(string) string) string { return getenv("MPC_REVIEW_GETTER_SECRET") }'
  Test-ReviewFailureCode $reviewFailures $injectedGetterFixture 'RCFG_UNDECLARED_READ' 'injected Go getter literal'

  $untrackedSdkFixture = New-PositiveFixture; $fixtures.Add($untrackedSdkFixture)
  & git -C $untrackedSdkFixture init --quiet
  & git -C $untrackedSdkFixture config core.autocrlf false
  & git -C $untrackedSdkFixture config user.email 'fixture@example.invalid'
  & git -C $untrackedSdkFixture config user.name 'Fixture'
  & git -C $untrackedSdkFixture add .
  & git -C $untrackedSdkFixture commit --quiet -m baseline
  $untrackedBaseSha = (& git -C $untrackedSdkFixture rev-parse HEAD).Trim()
  Write-FixtureFile $untrackedSdkFixture 'packages/sdk-runtime/src/review-bypass.ts' 'export const reviewBypass = true'
  Test-ReviewFailureCode $reviewFailures $untrackedSdkFixture 'GOV_API_SDK_SPLIT' 'untracked SDK-only change' $untrackedBaseSha

  # ADR-016 guards the SDK's types, not its tests or its plumbing: a change
  # that declares no contract shape must not demand an OpenAPI edit.
  $sdkTestOnlyFixture = New-PositiveFixture; $fixtures.Add($sdkTestOnlyFixture)
  & git -C $sdkTestOnlyFixture init --quiet
  & git -C $sdkTestOnlyFixture config core.autocrlf false
  & git -C $sdkTestOnlyFixture config user.email 'fixture@example.invalid'
  & git -C $sdkTestOnlyFixture config user.name 'Fixture'
  & git -C $sdkTestOnlyFixture add .
  & git -C $sdkTestOnlyFixture commit --quiet -m baseline
  $sdkTestOnlyBaseSha = (& git -C $sdkTestOnlyFixture rev-parse HEAD).Trim()
  Write-FixtureFile $sdkTestOnlyFixture 'packages/sdk-runtime/src/index.test.ts' 'export {}'
  Write-FixtureFile $sdkTestOnlyFixture 'packages/sdk-runtime/package.json' '{}'
  $sdkTestOnlyResult = Test-GovernanceDrift -RepositoryRoot $sdkTestOnlyFixture -BaseSha $sdkTestOnlyBaseSha
  Assert-True ('GOV_API_SDK_SPLIT' -notin @($sdkTestOnlyResult.Violations.ErrorCode)) 'test-only SDK change wrongly demanded an OpenAPI edit'
  # And the exclusion must not widen: a src file alongside the test still fires.
  Write-FixtureFile $sdkTestOnlyFixture 'packages/sdk-runtime/src/review-shape.ts' 'export const reviewShape = true'
  Test-ReviewFailureCode $reviewFailures $sdkTestOnlyFixture 'GOV_API_SDK_SPLIT' 'SDK src change alongside test' $sdkTestOnlyBaseSha

  $multilinePanicFixture = New-PositiveFixture; $fixtures.Add($multilinePanicFixture)
  Write-FixtureFile $multilinePanicFixture 'apps/server_core/internal/modules/catalog/domain/review_panic.go' @'
package domain
func bad() {
  panic(
    "review"
  )
}
'@
  Test-ReviewFailureCode $reviewFailures $multilinePanicFixture 'GOV_PRODUCTION_PANIC' 'multiline production panic'

  # A context under internal/contexts must be registrable, and its id may collide
  # with a module's — internal/modules/<id> and internal/contexts/<id> coexist for
  # the whole migration — so uniqueness is on the pair (kind, id). The id used
  # here is deliberately one the live registry already spends on a module, and
  # deliberately NOT the context the repository has actually registered: the
  # assertion has to prove the pair key, not re-read today's modules.json.
  $contextFixture = New-PositiveFixture; $fixtures.Add($contextFixture)
  $contextRegistryPath = Join-Path $contextFixture 'contracts/governance/modules.json'
  $contextRegistry = Get-Content -Raw $contextRegistryPath | ConvertFrom-Json
  Assert-True (@($contextRegistry.modules | Where-Object { $_.id -eq 'pricing' }).Count -eq 1) 'fixture premise broke: pricing is no longer exactly one registry entry'
  $contextRegistry.modules += [pscustomobject]@{
    id                   = 'pricing'
    kind                 = 'context'
    root                 = 'apps/server_core/internal/contexts/pricing'
    code_owner_path      = 'apps/server_core/internal/contexts/pricing'
    composition_required = $false
    openapi_prefixes     = @()
    dependencies         = @()
  }
  $contextRegistry | ConvertTo-Json -Depth 20 | Set-Content -LiteralPath $contextRegistryPath -Encoding utf8NoBOM
  Write-FixtureFile $contextFixture 'apps/server_core/internal/contexts/pricing/contracts/doc.go' "package contracts`n"
  $contextResult = Test-GovernanceDrift -RepositoryRoot $contextFixture
  Assert-True $contextResult.Passed "registered context rejected: $(@($contextResult.Violations.ErrorCode) -join ',')"

  # The rule that actually closes the hole. Walking the registry can never find
  # what nobody registered, so the finding has to come from walking the tree.
  $orphanContextFixture = New-PositiveFixture; $fixtures.Add($orphanContextFixture)
  Write-FixtureFile $orphanContextFixture 'apps/server_core/internal/contexts/orphan/contracts/doc.go' "package contracts`n"
  $orphanResult = Test-GovernanceDrift -RepositoryRoot $orphanContextFixture
  Assert-True (-not $orphanResult.Passed) 'unregistered context directory produced no finding'
  Assert-True ('GOV_CONTEXT_UNREGISTERED' -in @($orphanResult.Violations.ErrorCode)) "missing GOV_CONTEXT_UNREGISTERED; actual=$(@($orphanResult.Violations.ErrorCode) -join ',')"
  Assert-True ('orphan' -in @($orphanResult.Violations.Id)) "GOV_CONTEXT_UNREGISTERED did not name the directory; actual=$(@($orphanResult.Violations.Id) -join ',')"
  Write-Output 'guard_ran=gov-context-unregistered'

  # GOV_MODULE_LAYER: a cross-module import that lands on adapters/ is coupling
  # with an extra step, not translation (ADR-023 §2). Policy.psm1:378-382. The file
  # sits in domain/, not application/, so GOV_APPLICATION_IMPORT stays out of the
  # blast radius and the assertion isolates the code under test. The source module
  # must be DECLARED (Policy.psm1:370 skips undeclared sources); orders is, at
  # contracts/governance/modules.json:125.
  $layerFixture = New-PositiveFixture; $fixtures.Add($layerFixture)
  Write-FixtureFile $layerFixture 'apps/server_core/internal/modules/orders/domain/uses_catalog_adapters.go' "package domain`n`nimport _ `"marketplace-central/apps/server_core/internal/modules/catalog/adapters/postgres`"`n"
  $layerResult = Test-GovernanceDrift -RepositoryRoot $layerFixture
  Assert-True (-not $layerResult.Passed) 'cross-module adapters import produced no finding'
  Assert-True ('GOV_MODULE_LAYER' -in @($layerResult.Violations.ErrorCode)) "expected GOV_MODULE_LAYER, got: $(@($layerResult.Violations.ErrorCode) -join ',')"
  Write-Output 'guard_ran=gov-module-layer'

  # GOV_POSTGRES_DRIVER: database/sql under adapters/postgres bypasses pgdb.
  # Policy.psm1:412-417.
  $driverFixture = New-PositiveFixture; $fixtures.Add($driverFixture)
  Write-FixtureFile $driverFixture 'apps/server_core/internal/modules/orders/adapters/postgres/uses_database_sql.go' "package postgres`n`nimport _ `"database/sql`"`n"
  $driverResult = Test-GovernanceDrift -RepositoryRoot $driverFixture
  Assert-True (-not $driverResult.Passed) 'database/sql under adapters/postgres produced no finding'
  Assert-True ('GOV_POSTGRES_DRIVER' -in @($driverResult.Violations.ErrorCode)) "expected GOV_POSTGRES_DRIVER, got: $(@($driverResult.Violations.ErrorCode) -join ',')"
  Write-Output 'guard_ran=gov-postgres-driver'

  Assert-True ($reviewFailures.Count -eq 0) ($reviewFailures -join '; ')

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
