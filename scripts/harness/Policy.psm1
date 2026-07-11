Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

function New-PolicyIssue {
  param([string]$ErrorCode, [string]$Id, [string]$Path = '')
  [pscustomobject]@{ ErrorCode = $ErrorCode; Id = $Id; Path = $Path }
}

function New-PolicyResult {
  param(
    [Collections.IEnumerable]$Violations = @(),
    [Collections.IEnumerable]$BaselineExceptions = @(),
    [hashtable]$Documents = @{}
  )
  $violationList = @($Violations)
  [pscustomobject]@{
    PSTypeName = 'MarketplaceCentral.Governance.PolicyResult'
    Passed = $violationList.Count -eq 0
    Status = if ($violationList.Count -eq 0) { 'passed' } else { 'failed' }
    ErrorCode = if ($violationList.Count -eq 0) { $null } else { [string]$violationList[0].ErrorCode }
    Violations = $violationList
    BaselineExceptions = @($BaselineExceptions)
    Documents = $Documents
  }
}

function ConvertTo-NormalizedPath {
  param([string]$RepositoryRoot, [string]$Path)
  $root = [IO.Path]::GetFullPath($RepositoryRoot).TrimEnd([IO.Path]::DirectorySeparatorChar, [IO.Path]::AltDirectorySeparatorChar)
  $full = [IO.Path]::GetFullPath($Path)
  if (-not $full.StartsWith($root + [IO.Path]::DirectorySeparatorChar, [StringComparison]::OrdinalIgnoreCase) -and $full -ne $root) { return '' }
  return $full.Substring($root.Length).TrimStart([IO.Path]::DirectorySeparatorChar, [IO.Path]::AltDirectorySeparatorChar).Replace('\', '/')
}

function Resolve-RepositoryPath {
  param([string]$RepositoryRoot, [string]$RelativePath)
  Join-Path $RepositoryRoot ($RelativePath.Replace('/', [IO.Path]::DirectorySeparatorChar))
}

function Import-GovernanceDocument {
  [CmdletBinding()]
  param([Parameter(Mandatory)][string]$Path, [Parameter(Mandatory)][string]$SchemaPath)

  if (-not (Test-Path -LiteralPath $Path -PathType Leaf) -or -not (Test-Path -LiteralPath $SchemaPath -PathType Leaf)) {
    return [pscustomobject]@{ Passed = $false; ErrorCode = 'GOV_DOCUMENT_MISSING'; Document = $null; Path = $Path }
  }
  try {
    $valid = Test-Json -LiteralPath $Path -SchemaFile $SchemaPath -ErrorAction Stop
    if (-not $valid) { throw 'schema rejected document' }
    $document = Get-Content -Raw -LiteralPath $Path | ConvertFrom-Json -Depth 100
    return [pscustomobject]@{ Passed = $true; ErrorCode = $null; Document = $document; Path = $Path }
  } catch {
    return [pscustomobject]@{ Passed = $false; ErrorCode = 'GOV_SCHEMA_INVALID'; Document = $null; Path = $Path }
  }
}

function Test-UniqueIds {
  param([object[]]$Items, [string]$Property = 'id')
  $values = @($Items | ForEach-Object { [string]$_.$Property })
  return @($values | Group-Object | Where-Object Count -gt 1).Count -eq 0
}

function Test-GovernanceContracts {
  [CmdletBinding()]
  param([Parameter(Mandatory)][string]$RepositoryRoot)

  $documents = @{}
  $issues = [Collections.Generic.List[object]]::new()
  $registryNames = @('modules', 'runtime-config', 'execution-lanes', 'invariants', 'shared-seams')
  foreach ($name in $registryNames) {
    $path = Resolve-RepositoryPath $RepositoryRoot "contracts/governance/$name.json"
    $schemaPath = Resolve-RepositoryPath $RepositoryRoot "contracts/governance/schemas/$name.schema.json"
    $loaded = Import-GovernanceDocument -Path $path -SchemaPath $schemaPath
    if (-not $loaded.Passed) {
      $issues.Add((New-PolicyIssue $loaded.ErrorCode $name "contracts/governance/$name.json"))
    } else {
      $documents[$name] = $loaded.Document
    }
  }
  $contextSchema = Resolve-RepositoryPath $RepositoryRoot 'contracts/governance/schemas/context-pack.schema.json'
  if (-not (Test-Path -LiteralPath $contextSchema -PathType Leaf)) {
    $issues.Add((New-PolicyIssue 'GOV_DOCUMENT_MISSING' 'context-pack-schema' 'contracts/governance/schemas/context-pack.schema.json'))
  }
  if ($issues.Count -gt 0) { return New-PolicyResult -Violations $issues -Documents $documents }

  $modules = @($documents.modules.modules)
  $moduleIds = @($modules.id)
  if (-not (Test-UniqueIds $modules) -or -not (Test-UniqueIds @($documents.modules.temporary_exceptions))) {
    $issues.Add((New-PolicyIssue 'GOV_REFERENCE_INVALID' 'modules-duplicate-id' 'contracts/governance/modules.json'))
  }
  foreach ($module in $modules) {
    foreach ($dependency in @($module.dependencies)) {
      if ($dependency -notin $moduleIds -or $dependency -eq $module.id) {
        $issues.Add((New-PolicyIssue 'GOV_REFERENCE_INVALID' "module-dependency-$($module.id)-$dependency" 'contracts/governance/modules.json'))
      }
    }
  }
  foreach ($exception in @($documents.modules.temporary_exceptions)) {
    if ($exception.source_module -notin $moduleIds -or $exception.target_module -notin $moduleIds -or $exception.rule_id -ne 'module-target-layer' -or [string]$exception.path -match '[*?\[]') {
      $issues.Add((New-PolicyIssue 'GOV_REFERENCE_INVALID' ([string]$exception.id) 'contracts/governance/modules.json'))
    }
  }

  $runtimeKeys = @($documents.'runtime-config'.keys)
  $runtimeExceptions = @($documents.'runtime-config'.temporary_exceptions)
  $keyGroups = @($runtimeKeys | Group-Object key)
  foreach ($group in @($keyGroups | Where-Object Count -gt 1)) {
    $targets = @($group.Group.alias_for | ForEach-Object { [string]$_ } | Sort-Object -Unique)
    $code = if (@($group.Group | Where-Object lifecycle -eq 'legacy_alias').Count -gt 0 -or $targets.Count -gt 1) { 'RCFG_ALIAS_COLLISION' } else { 'GOV_REFERENCE_INVALID' }
    $issues.Add((New-PolicyIssue $code ([string]$group.Name) 'contracts/governance/runtime-config.json'))
  }
  if (-not (Test-UniqueIds $runtimeExceptions)) {
    $issues.Add((New-PolicyIssue 'GOV_REFERENCE_INVALID' 'runtime-exception-duplicate-id' 'contracts/governance/runtime-config.json'))
  }
  $runtimeKeyNames = @($runtimeKeys.key)
  $runtimeExceptionIds = @($runtimeExceptions.id)
  foreach ($key in $runtimeKeys) {
    if ($key.lifecycle -eq 'legacy_alias' -and ([string]::IsNullOrWhiteSpace([string]$key.alias_for) -or $key.alias_for -notin $runtimeKeyNames -or $key.alias_for -eq $key.key)) {
      $issues.Add((New-PolicyIssue 'RCFG_ALIAS_UNDECLARED' ([string]$key.key) 'contracts/governance/runtime-config.json'))
    }
    if ([string]$key.key -match '(?:PASSWORD|SECRET|TOKEN|DATABASE_URL|POSTGRES_URL|ENCRYPTION_KEY)$' -and $key.sensitivity -ne 'secret') {
      $issues.Add((New-PolicyIssue 'RCFG_SECRET_CLASS_MISMATCH' ([string]$key.key) 'contracts/governance/runtime-config.json'))
    }
    foreach ($reader in @($key.readers)) {
      if ($reader.status -eq 'temporary_exception' -and ([string]::IsNullOrWhiteSpace([string]$reader.exception_id) -or $reader.exception_id -notin $runtimeExceptionIds)) {
        $issues.Add((New-PolicyIssue 'GOV_REFERENCE_INVALID' ([string]$key.key) ([string]$reader.path)))
      }
    }
  }
  foreach ($exception in $runtimeExceptions) {
    $key = $runtimeKeys | Where-Object key -eq $exception.key | Select-Object -First 1
    $reader = @($key.readers | Where-Object { $_.path -eq $exception.path -and $_.exception_id -eq $exception.id })
    if ($null -eq $key -or $reader.Count -ne 1 -or [string]$exception.path -match '[*?\[]') {
      $issues.Add((New-PolicyIssue 'GOV_REFERENCE_INVALID' ([string]$exception.id) 'contracts/governance/runtime-config.json'))
    }
  }

  $lanes = @($documents.'execution-lanes'.lanes)
  if (-not (Test-UniqueIds $lanes)) { $issues.Add((New-PolicyIssue 'GOV_REFERENCE_INVALID' 'lane-duplicate-id' 'contracts/governance/execution-lanes.json')) }
  foreach ($lane in $lanes) {
    foreach ($keyName in @($lane.allowed_runtime_keys)) {
      $key = $runtimeKeys | Where-Object key -eq $keyName | Select-Object -First 1
      if ($null -eq $key -or $lane.id -notin @($key.allowed_lanes)) {
        $issues.Add((New-PolicyIssue 'RCFG_LANE_VIOLATION' "$($lane.id)-$keyName" 'contracts/governance/execution-lanes.json'))
      }
    }
  }

  $invariants = @($documents.invariants.invariants)
  if (-not (Test-UniqueIds $invariants) -or -not (Test-UniqueIds @($documents.invariants.temporary_exceptions))) {
    $issues.Add((New-PolicyIssue 'GOV_REFERENCE_INVALID' 'invariant-duplicate-id' 'contracts/governance/invariants.json'))
  }
  foreach ($exception in @($documents.invariants.temporary_exceptions)) {
    $rule = $invariants | Where-Object id -eq $exception.rule_id | Select-Object -First 1
    if ($null -eq $rule -or $rule.exception_mode -eq 'none' -or @($exception.paths | Where-Object { [string]$_ -match '[*?\[]' }).Count -gt 0) {
      $issues.Add((New-PolicyIssue 'GOV_REFERENCE_INVALID' ([string]$exception.id) 'contracts/governance/invariants.json'))
    }
  }
  if (-not (Test-UniqueIds @($documents.'shared-seams'.seams))) {
    $issues.Add((New-PolicyIssue 'GOV_REFERENCE_INVALID' 'seam-duplicate-id' 'contracts/governance/shared-seams.json'))
  }
  foreach ($seam in @($documents.'shared-seams'.seams)) {
    if (@($seam.exclusive_paths | Where-Object { [string]$_ -match '[*?\[]' }).Count -gt 0) {
      $issues.Add((New-PolicyIssue 'GOV_REFERENCE_INVALID' ([string]$seam.id) 'contracts/governance/shared-seams.json'))
    }
  }
  return New-PolicyResult -Violations $issues -Documents $documents
}

function Add-BaselineException {
  param([Collections.Generic.List[object]]$List, [hashtable]$Seen, [string]$Id, [string]$Path)
  if (-not $Seen.ContainsKey($Id)) {
    $Seen[$Id] = $true
    $List.Add([pscustomobject]@{ Id = $Id; Path = $Path })
  }
}

function Get-SourceFiles {
  param([string]$RepositoryRoot)
  $extensions = @('.go', '.ps1', '.sh', '.ts', '.tsx', '.js', '.mjs', '.yml', '.yaml')
  @(Get-ChildItem -LiteralPath $RepositoryRoot -Recurse -File -Force | Where-Object {
    $relative = ConvertTo-NormalizedPath $RepositoryRoot $_.FullName
    $relative -notmatch '^(?:\.git|\.mnfs|node_modules|scripts/\.runs|scripts/tests|contracts/governance)/' -and
      ($_.Extension -in $extensions -or $relative -eq 'docker-compose.yml')
  })
}

function Get-EnvironmentReads {
  param([string]$RepositoryRoot, [object[]]$Files)
  $reads = [Collections.Generic.List[object]]::new()
  $ignored = @('PATH', 'TEMP', 'TMP', 'GOCACHE', 'DEV')
  $ownedPattern = '^(?:MC_|MS_|MPC_|MPC_TEST_|SANKHYA_|ME_|NGROK_|VITE_|SERVER_ADDR$|API_PORT$|RUN_MIGRATIONS$)'
  $patterns = @(
    'os\.(?:Getenv|LookupEnv)\(\s*["''](?<key>[A-Z][A-Z0-9_]*)["'']',
    '(?:process|import\.meta)\.env\.(?<key>[A-Z][A-Z0-9_]*)',
    'GetEnvironmentVariable\(\s*["''](?<key>[A-Z][A-Z0-9_]*)["'']',
    '\$\{(?<key>[A-Z][A-Z0-9_]*)'
  )
  foreach ($file in $Files) {
    $relative = ConvertTo-NormalizedPath $RepositoryRoot $file.FullName
    $content = Get-Content -Raw -LiteralPath $file.FullName
    foreach ($pattern in $patterns) {
      foreach ($match in [regex]::Matches($content, $pattern)) {
        $key = [string]$match.Groups['key'].Value
        if ($key -notin $ignored -and $key -match $ownedPattern) {
          $reads.Add([pscustomobject]@{ Key = $key; Path = $relative })
        }
      }
    }
  }
  @($reads | Sort-Object Key, Path -Unique)
}

function Test-GovernanceDrift {
  [CmdletBinding()]
  param([Parameter(Mandatory)][string]$RepositoryRoot, [string]$BaseSha = '')

  $contract = Test-GovernanceContracts -RepositoryRoot $RepositoryRoot
  if (-not $contract.Passed) { return $contract }
  $documents = $contract.Documents
  $issues = [Collections.Generic.List[object]]::new()
  $baselines = [Collections.Generic.List[object]]::new()
  $baselineSeen = @{}
  $files = Get-SourceFiles $RepositoryRoot

  $moduleRegistry = $documents.modules
  $moduleRoot = Resolve-RepositoryPath $RepositoryRoot 'apps/server_core/internal/modules'
  $actualModules = if (Test-Path -LiteralPath $moduleRoot -PathType Container) { @(Get-ChildItem -LiteralPath $moduleRoot -Directory | Select-Object -ExpandProperty Name | Sort-Object) } else { @() }
  $declaredModules = @($moduleRegistry.modules.id | Sort-Object)
  foreach ($id in @($actualModules | Where-Object { $_ -notin $declaredModules })) { $issues.Add((New-PolicyIssue 'GOV_MODULE_COVERAGE' $id "apps/server_core/internal/modules/$id")) }
  foreach ($id in @($declaredModules | Where-Object { $_ -notin $actualModules })) { $issues.Add((New-PolicyIssue 'GOV_MODULE_COVERAGE' $id "apps/server_core/internal/modules/$id")) }
  foreach ($module in @($moduleRegistry.modules)) {
    if ($module.root -ne "apps/server_core/internal/modules/$($module.id)" -or -not (Test-Path -LiteralPath (Resolve-RepositoryPath $RepositoryRoot ([string]$module.root)) -PathType Container)) {
      $issues.Add((New-PolicyIssue 'GOV_MODULE_COVERAGE' ([string]$module.id) ([string]$module.root)))
    }
  }

  $moduleById = @{}; foreach ($module in @($moduleRegistry.modules)) { $moduleById[[string]$module.id] = $module }
  foreach ($file in @($files | Where-Object { (ConvertTo-NormalizedPath $RepositoryRoot $_.FullName) -match '^apps/server_core/internal/modules/[^/]+/.+(?<!_test)\.go$' })) {
    $path = ConvertTo-NormalizedPath $RepositoryRoot $file.FullName
    if ($path -notmatch '^apps/server_core/internal/modules/(?<source>[^/]+)/') { continue }
    $source = $Matches.source
    $content = Get-Content -Raw -LiteralPath $file.FullName
    foreach ($match in [regex]::Matches($content, '["'']marketplace-central/apps/server_core/internal/modules/(?<target>[a-z_]+)/(?<layer>[a-z_]+)[^"'']*["'']')) {
      $target = [string]$match.Groups['target'].Value; $layer = [string]$match.Groups['layer'].Value
      if ($target -eq $source) { continue }
      if ($target -notin @($moduleById[$source].dependencies)) {
        $issues.Add((New-PolicyIssue 'GOV_MODULE_DEPENDENCY' "$source-$target" $path))
      }
      if ($layer -in @('adapters', 'transport', 'registry')) {
        $importPath = $match.Value.Trim('"', "'")
        $exception = @($moduleRegistry.temporary_exceptions | Where-Object { $_.source_module -eq $source -and $_.target_module -eq $target -and $_.target_layer -eq $layer -and $_.path -eq $path -and $_.import_path -eq $importPath })
        if ($exception.Count -eq 1) { Add-BaselineException $baselines $baselineSeen ([string]$exception[0].id) $path }
        else { $issues.Add((New-PolicyIssue 'GOV_MODULE_LAYER' "$source-$target-$layer" $path)) }
      }
    }
  }

  $compositionPath = Resolve-RepositoryPath $RepositoryRoot 'apps/server_core/internal/composition/root.go'
  $composition = if (Test-Path -LiteralPath $compositionPath -PathType Leaf) { Get-Content -Raw -LiteralPath $compositionPath } else { '' }
  foreach ($module in @($moduleRegistry.modules | Where-Object composition_required)) {
    if ($composition -notmatch [regex]::Escape("/modules/$($module.id)")) {
      $issues.Add((New-PolicyIssue 'GOV_COMPOSITION_MISSING' ([string]$module.id) 'apps/server_core/internal/composition/root.go'))
    }
  }

  foreach ($file in @($files | Where-Object { (ConvertTo-NormalizedPath $RepositoryRoot $_.FullName) -match '^apps/server_core/internal/modules/[^/]+/application/.+(?<!_test)\.go$' })) {
    $path = ConvertTo-NormalizedPath $RepositoryRoot $file.FullName
    $content = Get-Content -Raw -LiteralPath $file.FullName
    $sourceModule = if ($path -match '^apps/server_core/internal/modules/(?<source>[^/]+)/') { $Matches.source } else { '' }
    $forbidden = $content -match '["''](?:net/http|github\.com/jackc/pgx[^"'']*)["'']'
    foreach ($match in [regex]::Matches($content, '["''][^"'']*/modules/(?<target>[^/]+)/(?:adapters|transport|registry)(?:/[^"'']*)?["'']')) {
      if ($match.Groups['target'].Value -ne $sourceModule) { $forbidden = $true }
    }
    if ($forbidden) {
      $issues.Add((New-PolicyIssue 'GOV_APPLICATION_IMPORT' ([IO.Path]::GetFileNameWithoutExtension($path)) $path))
    }
  }

  foreach ($file in @($files | Where-Object { (ConvertTo-NormalizedPath $RepositoryRoot $_.FullName) -match '^apps/server_core/internal/modules/[^/]+/adapters/postgres/.+\.go$' })) {
    $path = ConvertTo-NormalizedPath $RepositoryRoot $file.FullName
    if ((Get-Content -Raw -LiteralPath $file.FullName) -match '["'']database/sql["'']') {
      $exception = @($documents.invariants.temporary_exceptions | Where-Object { $_.rule_id -eq 'postgres-driver' -and $path -in @($_.paths) })
      if ($exception.Count -eq 1) { Add-BaselineException $baselines $baselineSeen ([string]$exception[0].id) $path } else { $issues.Add((New-PolicyIssue 'GOV_POSTGRES_DRIVER' 'database-sql' $path)) }
    }
  }

  foreach ($file in @($files | Where-Object { (ConvertTo-NormalizedPath $RepositoryRoot $_.FullName) -match '^apps/server_core/(?!.*_test\.go$).+\.go$' })) {
    $path = ConvertTo-NormalizedPath $RepositoryRoot $file.FullName
    $content = Get-Content -Raw -LiteralPath $file.FullName
    $panicMatches = @([regex]::Matches($content, 'panic\s*\([^\r\n]+\)'))
    if ($panicMatches.Count -eq 0) { continue }
    $exceptions = @($documents.invariants.temporary_exceptions | Where-Object { $_.rule_id -eq 'production-panic' -and $path -in @($_.paths) })
    $covered = 0
    foreach ($exception in $exceptions) {
      foreach ($occurrence in @($exception.occurrences | Where-Object path -eq $path)) {
        $actualCount = @($panicMatches | Where-Object { $_.Value -eq $occurrence.fingerprint }).Count
        if ($actualCount -eq [int]$occurrence.count) { $covered += $actualCount; Add-BaselineException $baselines $baselineSeen ([string]$exception.id) $path }
      }
    }
    if ($covered -ne $panicMatches.Count) { $issues.Add((New-PolicyIssue 'GOV_PRODUCTION_PANIC' 'panic' $path)) }
  }

  $runtime = $documents.'runtime-config'
  $reads = Get-EnvironmentReads $RepositoryRoot $files
  foreach ($file in $files) {
    $path = ConvertTo-NormalizedPath $RepositoryRoot $file.FullName
    if ($path -eq 'scripts/harness.ps1') { continue }
    $content = Get-Content -Raw -LiteralPath $file.FullName
    if ($content -match 'os\.(?:Getenv|LookupEnv)\(\s*(?!["''])[^)]+\)' -or $content -match 'GetEnvironmentVariable\(\s*\$[A-Za-z_]') {
      $issues.Add((New-PolicyIssue 'RCFG_DYNAMIC_READER_UNBOUNDED' 'dynamic-reader' $path))
    }
  }
  $keysByName = @{}; foreach ($key in @($runtime.keys)) { if (-not $keysByName.ContainsKey([string]$key.key)) { $keysByName[[string]$key.key] = $key } }
  $observedReaderPairs = @{}; foreach ($read in $reads) { $observedReaderPairs["$($read.Key)|$($read.Path)"] = $true }
  foreach ($read in $reads) {
    if (-not $keysByName.ContainsKey([string]$read.Key)) {
      $issues.Add((New-PolicyIssue 'RCFG_UNDECLARED_READ' ([string]$read.Key) ([string]$read.Path))); continue
    }
    $key = $keysByName[[string]$read.Key]
    $reader = @($key.readers | Where-Object path -eq $read.Path)
    if ($reader.Count -ne 1) { $issues.Add((New-PolicyIssue 'RCFG_UNAPPROVED_READER' ([string]$read.Key) ([string]$read.Path))); continue }
    if ($reader[0].status -eq 'temporary_exception') { Add-BaselineException $baselines $baselineSeen ([string]$reader[0].exception_id) ([string]$read.Path) }
  }
  foreach ($key in @($runtime.keys)) {
    foreach ($reader in @($key.readers)) {
      $fullPath = Resolve-RepositoryPath $RepositoryRoot ([string]$reader.path)
      if (-not (Test-Path -LiteralPath $fullPath -PathType Leaf)) {
        $issues.Add((New-PolicyIssue 'RCFG_READER_MISSING' ([string]$key.key) ([string]$reader.path))); continue
      }
      $content = Get-Content -Raw -LiteralPath $fullPath
      if ($content -notmatch "(?<![A-Z0-9_])$([regex]::Escape([string]$key.key))(?![A-Z0-9_])") {
        $issues.Add((New-PolicyIssue 'RCFG_READER_MISSING' ([string]$key.key) ([string]$reader.path)))
      } elseif ($reader.status -eq 'temporary_exception' -and -not $observedReaderPairs.ContainsKey("$($key.key)|$($reader.path)")) {
        $issues.Add((New-PolicyIssue 'RCFG_READER_MISSING' ([string]$key.key) ([string]$reader.path)))
      } elseif ($reader.status -eq 'temporary_exception') {
        Add-BaselineException $baselines $baselineSeen ([string]$reader.exception_id) ([string]$reader.path)
      }
    }
  }

  $migrationRoot = Resolve-RepositoryPath $RepositoryRoot 'apps/server_core/migrations'
  if (Test-Path -LiteralPath $migrationRoot -PathType Container) {
    foreach ($group in @(Get-ChildItem -LiteralPath $migrationRoot -File -Filter '*.sql' | Where-Object Name -match '^(?<prefix>[0-9]+)_' | Group-Object { if ($_.Name -match '^(?<prefix>[0-9]+)_') { $Matches.prefix } } | Where-Object Count -gt 1)) {
      $paths = @($group.Group | ForEach-Object { ConvertTo-NormalizedPath $RepositoryRoot $_.FullName } | Sort-Object)
      $exception = @($documents.invariants.temporary_exceptions | Where-Object { $_.rule_id -eq 'migration-prefix' -and (@($_.paths | Sort-Object) -join '|') -eq ($paths -join '|') })
      if ($exception.Count -eq 1) { Add-BaselineException $baselines $baselineSeen ([string]$exception[0].id) ($paths -join ',') }
      else { $issues.Add((New-PolicyIssue 'GOV_MIGRATION_PREFIX' ([string]$group.Name) ($paths -join ','))) }
    }
  }

  foreach ($file in @($files | Where-Object { (ConvertTo-NormalizedPath $RepositoryRoot $_.FullName) -match '^(?:apps/web/src|packages)/(?!sdk-runtime/).+\.(?:ts|tsx|js|mjs)$' })) {
    $path = ConvertTo-NormalizedPath $RepositoryRoot $file.FullName
    if ((Get-Content -Raw -LiteralPath $file.FullName) -match '(?<![A-Za-z0-9_])fetch\s*\(') {
      $exception = @($documents.invariants.temporary_exceptions | Where-Object { $_.rule_id -eq 'frontend-fetch' -and $path -in @($_.paths) })
      if ($exception.Count -eq 1) { Add-BaselineException $baselines $baselineSeen ([string]$exception[0].id) $path } else { $issues.Add((New-PolicyIssue 'GOV_FRONTEND_FETCH' 'direct-fetch' $path)) }
    }
  }

  if (-not [string]::IsNullOrWhiteSpace($BaseSha)) {
    if ($BaseSha -notmatch '^[0-9a-f]{40}$') {
      $issues.Add((New-PolicyIssue 'GOV_SEMANTIC_DRIFT' 'base-sha-invalid'))
    } else {
      $changed = @(& git -C $RepositoryRoot diff --name-only $BaseSha -- 2>$null | ForEach-Object { ([string]$_).Replace('\', '/') })
      if ($LASTEXITCODE -ne 0) { $issues.Add((New-PolicyIssue 'GOV_SEMANTIC_DRIFT' 'base-sha-unavailable')) }
      else {
        $apiChanged = 'contracts/api/marketplace-central.openapi.yaml' -in $changed
        $sdkChanged = @($changed | Where-Object { $_ -eq 'packages/sdk-runtime' -or $_ -like 'packages/sdk-runtime/*' }).Count -gt 0
        if ($apiChanged -xor $sdkChanged) { $issues.Add((New-PolicyIssue 'GOV_API_SDK_SPLIT' 'api-sdk-atomicity')) }
      }
    }
  }
  return New-PolicyResult -Violations $issues -BaselineExceptions ($baselines | Sort-Object Id) -Documents $documents
}

function Test-PathIntersection {
  param([string]$Left, [string]$Right)
  $leftPath = $Left.TrimEnd('/'); $rightPath = $Right.TrimEnd('/')
  return $leftPath -eq $rightPath -or $leftPath.StartsWith($rightPath + '/', [StringComparison]::Ordinal) -or $rightPath.StartsWith($leftPath + '/', [StringComparison]::Ordinal)
}

function Get-ApplicableInvariants {
  [CmdletBinding()]
  param([Parameter(Mandatory)][object]$Registry, [Parameter(Mandatory)][string[]]$Paths)
  @($Registry.invariants | Where-Object {
    $invariant = $_
    @($Paths | Where-Object { $path = $_; @($invariant.scope_paths | Where-Object { (Test-PathIntersection $path ([string]$_)) }).Count -gt 0 }).Count -gt 0
  })
}

function Get-SharedSeams {
  [CmdletBinding()]
  param([Parameter(Mandatory)][object]$Registry, [Parameter(Mandatory)][string[]]$Paths)
  @($Registry.seams | Where-Object {
    $seam = $_
    @($Paths | Where-Object { $path = $_; @($seam.exclusive_paths | Where-Object { (Test-PathIntersection $path ([string]$_)) }).Count -gt 0 }).Count -gt 0
  })
}

Export-ModuleMember -Function Import-GovernanceDocument, Test-GovernanceContracts, Test-GovernanceDrift, Get-ApplicableInvariants, Get-SharedSeams
