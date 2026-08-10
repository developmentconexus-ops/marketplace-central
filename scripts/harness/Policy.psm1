Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

function New-PolicyIssue {
  param([string]$ErrorCode, [string]$Id, [string]$Path = '')
  [pscustomobject]@{ ErrorCode = $ErrorCode; Id = $Id; Path = $Path }
}

# A registry entry's kind decides where its root must live. Absent means module,
# so the entries that predate ADR-023's amendment need no edit. Read through
# PSObject because Set-StrictMode -Version Latest throws on a missing property.
function Get-RegistryEntryKind {
  param([object]$Entry)
  $property = $Entry.PSObject.Properties['kind']
  if ($null -eq $property -or [string]::IsNullOrWhiteSpace([string]$property.Value)) { return 'module' }
  return [string]$property.Value
}

function Get-RegistryEntryRoot {
  param([object]$Entry)
  switch (Get-RegistryEntryKind $Entry) {
    'context' { "apps/server_core/internal/contexts/$($Entry.id)" }
    default { "apps/server_core/internal/modules/$($Entry.id)" }
  }
}

function New-PolicyResult {
  param(
    [Collections.IEnumerable]$Violations = @(),
    [Collections.IEnumerable]$BaselineExceptions = @(),
    [hashtable]$Documents = @{},
    # -1 means "this result is not from a tree walk" (validate) or "the walk
    # never finished". Only Test-GovernanceDrift sets a real count; a caller
    # ratcheting drift findings must refuse both -1 and 0, because a verdict
    # from a walk that saw no files is a verdict about nothing.
    [int]$FilesScanned = -1
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
    FilesScanned = $FilesScanned
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
  $registryNames = @('modules', 'runtime-config', 'execution-lanes', 'invariants', 'shared-seams', 'knowledge-routes')
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
  # Dependencies name modules, so a context id must not become resolvable here
  # just by existing in the same registry.
  $moduleIds = @($modules | Where-Object { (Get-RegistryEntryKind $_) -eq 'module' } | ForEach-Object { [string]$_.id })
  # Uniqueness is on the pair (kind, id): internal/modules/catalog and
  # internal/contexts/catalog coexist for the whole migration and both are
  # legitimately called catalog.
  $moduleKeys = @($modules | ForEach-Object { [pscustomobject]@{ id = "$(Get-RegistryEntryKind $_)/$([string]$_.id)" } })
  if (-not (Test-UniqueIds $moduleKeys) -or -not (Test-UniqueIds @($documents.modules.temporary_exceptions))) {
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
  $laneRuntimeKeys = @($lanes.allowed_runtime_keys | ForEach-Object { [string]$_ } | Sort-Object -Unique)
  $registryRuntimeKeys = @(
    $runtimeKeys |
      Where-Object { @($_.readers | Where-Object kind -eq 'registry').Count -gt 0 } |
      ForEach-Object { [string]$_.key } |
      Sort-Object -Unique
  )
  if (($laneRuntimeKeys -join '|') -ne ($registryRuntimeKeys -join '|')) {
    $issues.Add((New-PolicyIssue 'RCFG_REGISTRY_READER_COVERAGE' 'lane-runtime-keys' 'contracts/governance/runtime-config.json'))
  }
  foreach ($key in @($runtimeKeys | Where-Object { $_.key -in $laneRuntimeKeys })) {
    if (@($key.readers | Where-Object { $_.kind -eq 'registry' -and $_.status -eq 'approved' }).Count -ne 1) {
      $issues.Add((New-PolicyIssue 'RCFG_REGISTRY_READER_COVERAGE' ([string]$key.key) 'contracts/governance/runtime-config.json'))
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
  $extensions = @('.go', '.ps1', '.psm1', '.sh', '.ts', '.tsx', '.js', '.mjs', '.yml', '.yaml')
  @(Get-ChildItem -LiteralPath $RepositoryRoot -Recurse -File -Force | Where-Object {
    $relative = ConvertTo-NormalizedPath $RepositoryRoot $_.FullName
    # Derived or foreign trees are excluded at ANY depth. Anchoring these to the root let
    # whole checkouts of OTHER branches into the scan: .worktrees/<branch>/ and
    # .claude/worktrees/<branch>/ together contributed 264 of 266 apparent new violations
    # in the Onda 0 set diff, plus a 25-minute stall. A gate may only ever read the
    # checkout under test — anything else is a cross-branch measurement.
    $relative -notmatch '(?:^|/)(?:\.git|\.claude|node_modules|\.gocache|\.gomodcache|\.worktrees)/' -and
      $relative -notmatch '^(?:\.mnfs|scripts/\.runs|scripts/tests|contracts/governance)/' -and
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
    '(?i)(?<![A-Za-z0-9_.])(?:getenv|lookupenv)\(\s*["''](?<key>[A-Z][A-Z0-9_]*)["'']',
    '(?:process|import\.meta)\.env\.(?<key>[A-Z][A-Z0-9_]*)',
    'GetEnvironmentVariable\(\s*["''](?<key>[A-Z][A-Z0-9_]*)["'']',
    '\$env:(?<key>[A-Z][A-Z0-9_]*)',
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

function Test-RegistryDrivenEnvironmentReader {
  param([Parameter(Mandatory)][string]$Content)

  $tokens = $null
  $parseErrors = $null
  $ast = [System.Management.Automation.Language.Parser]::ParseInput($Content, [ref]$tokens, [ref]$parseErrors)
  if (@($parseErrors).Count -gt 0) { return $false }

  $functions = @($ast.FindAll({
    param($node)
    $node -is [System.Management.Automation.Language.FunctionDefinitionAst]
  }, $true))
  $runtimeFunctions = @($functions | Where-Object Name -eq 'Get-ProcessRuntimeValues')
  $childFunctions = @($functions | Where-Object Name -eq 'New-HarnessChildEnvironment')
  if ($runtimeFunctions.Count -ne 1 -or $childFunctions.Count -ne 1) { return $false }
  $runtimeFunction = $runtimeFunctions[0]
  $childFunction = $childFunctions[0]

  $dynamicLookups = @($ast.FindAll({
    param($node)
    $node -is [System.Management.Automation.Language.InvokeMemberExpressionAst] -and
      [string]$node.Member.Value -eq 'GetEnvironmentVariable' -and
      $node.Extent.Text -match '\(\s*\$[A-Za-z_][A-Za-z0-9_]*\s*,'
  }, $true))
  $runtimeLookups = @($dynamicLookups | Where-Object {
    $_.Extent.StartOffset -ge $runtimeFunction.Extent.StartOffset -and $_.Extent.EndOffset -le $runtimeFunction.Extent.EndOffset
  })
  $childLookups = @($dynamicLookups | Where-Object {
    $_.Extent.StartOffset -ge $childFunction.Extent.StartOffset -and $_.Extent.EndOffset -le $childFunction.Extent.EndOffset
  })
  if ($dynamicLookups.Count -ne 2 -or $runtimeLookups.Count -ne 1 -or $childLookups.Count -ne 1) { return $false }
  if ($runtimeLookups[0].Extent.Text -notmatch '^\[Environment\]::GetEnvironmentVariable\(\s*\$name\s*,\s*[''\"]Process[''\"]\s*\)$') { return $false }
  if ($childLookups[0].Extent.Text -notmatch '^\[Environment\]::GetEnvironmentVariable\(\s*\$key\s*,\s*[''\"]Process[''\"]\s*\)$') { return $false }

  $runtimeLoops = @($runtimeFunction.FindAll({
    param($node)
    $node -is [System.Management.Automation.Language.ForEachStatementAst] -and
      $node.Variable.VariablePath.UserPath -eq 'allowedKey' -and
      ($node.Condition.Extent.Text -replace '\s', '') -eq '@($Lane.allowed_runtime_keys)'
  }, $true))
  $safeToolLoops = @($childFunction.FindAll({
    param($node)
    $node -is [System.Management.Automation.Language.ForEachStatementAst] -and
      $node.Variable.VariablePath.UserPath -eq 'key' -and
      ($node.Condition.Extent.Text -replace '\s', '') -eq '$script:SafeToolKeys'
  }, $true))
  if ($runtimeLoops.Count -ne 1 -or $safeToolLoops.Count -ne 1) { return $false }
  if ($runtimeFunction.Extent.Text -notmatch '\$name\s*=\s*\[string\]\s*\$allowedKey') { return $false }
  if ($childFunction.Extent.Text -notmatch 'Get-HarnessRegistry\s+-RepositoryRoot\s+\$root\s+-Name\s+[''"]runtime-config[''"]') { return $false }
  return $true
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
  # Both directions below are about internal/modules only. A context declared in
  # the same registry must not be demanded as a module directory.
  $declaredModules = @($moduleRegistry.modules | Where-Object { (Get-RegistryEntryKind $_) -eq 'module' } | ForEach-Object { [string]$_.id } | Sort-Object)
  foreach ($id in @($actualModules | Where-Object { $_ -notin $declaredModules })) { $issues.Add((New-PolicyIssue 'GOV_MODULE_COVERAGE' $id "apps/server_core/internal/modules/$id")) }
  foreach ($id in @($declaredModules | Where-Object { $_ -notin $actualModules })) { $issues.Add((New-PolicyIssue 'GOV_MODULE_COVERAGE' $id "apps/server_core/internal/modules/$id")) }
  foreach ($module in @($moduleRegistry.modules)) {
    if ($module.root -ne (Get-RegistryEntryRoot $module) -or -not (Test-Path -LiteralPath (Resolve-RepositoryPath $RepositoryRoot ([string]$module.root)) -PathType Container)) {
      $issues.Add((New-PolicyIssue 'GOV_MODULE_COVERAGE' ([string]$module.id) ([string]$module.root)))
    }
  }

  # Walk the tree, not the registry. Every rule above iterates the registry, and
  # iterating the registry can never find what nobody registered: an unregistered
  # directory produces zero findings, which reads as compliance. The coverage
  # rule for internal/modules already walks the filesystem; internal/contexts had
  # no such rule at all.
  $contextsRoot = Resolve-RepositoryPath $RepositoryRoot 'apps/server_core/internal/contexts'
  if (Test-Path -LiteralPath $contextsRoot -PathType Container) {
    $registeredContexts = @($moduleRegistry.modules | Where-Object { (Get-RegistryEntryKind $_) -eq 'context' } | ForEach-Object { [string]$_.id })
    foreach ($directory in @(Get-ChildItem -LiteralPath $contextsRoot -Directory)) {
      if ($directory.Name -notin $registeredContexts) {
        $issues.Add((New-PolicyIssue 'GOV_CONTEXT_UNREGISTERED' $directory.Name "apps/server_core/internal/contexts/$($directory.Name)"))
      }
    }
  }

  # Modules only. This map is keyed by id alone, so a context entry sharing an id
  # would silently replace the module's — and take its declared dependencies with
  # it. That is why uniqueness is on (kind, id) and lookups filter by kind.
  $moduleById = @{}; foreach ($module in @($moduleRegistry.modules | Where-Object { (Get-RegistryEntryKind $_) -eq 'module' })) { $moduleById[[string]$module.id] = $module }
  foreach ($file in @($files | Where-Object { (ConvertTo-NormalizedPath $RepositoryRoot $_.FullName) -match '^apps/server_core/internal/modules/[^/]+/.+(?<!_test)\.go$' })) {
    $path = ConvertTo-NormalizedPath $RepositoryRoot $file.FullName
    if ($path -notmatch '^apps/server_core/internal/modules/(?<source>[^/]+)/') { continue }
    $source = $Matches.source
    # Undeclared module already surfaces as GOV_MODULE_COVERAGE above; skip
    # dependency checks instead of crashing on $moduleById[$source] under StrictMode.
    if (-not $moduleById.ContainsKey($source)) { continue }
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
    $compositionNeedle = switch (Get-RegistryEntryKind $module) {
      'context' { "/contexts/$($module.id)" }
      default { "/modules/$($module.id)" }
    }
    if ($composition -notmatch [regex]::Escape($compositionNeedle)) {
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

  # The `production-panic` check used to sit here. It is now forbidigo in
  # apps/server_core/.golangci.yml, counted by the lint-go ratchet. It left
  # because of what it cost, not because the rule was wrong: every sanctioned
  # panic needed a hand-written exception carrying a path, a symbol, a verbatim
  # source fingerprint compared by exact string equality, an occurrence count,
  # a paragraph of prose and a removal_owner -- and one of them had to be
  # declared TWICE because the same git blob materializes with either line
  # ending under core.autocrlf, so the fingerprint matched in one checkout and
  # not the other. Eight call sites, seven entries, ~90 lines of registry.
  # The ratchet expresses the same thing as `forbidigo: 8`.

  $runtime = $documents.'runtime-config'
  $reads = Get-EnvironmentReads $RepositoryRoot $files
  $registryReaderPaths = @(
    $runtime.keys.readers |
      Where-Object { $_.kind -eq 'registry' -and $_.status -eq 'approved' } |
      ForEach-Object { [string]$_.path } |
      Sort-Object -Unique
  )
  foreach ($file in $files) {
    $path = ConvertTo-NormalizedPath $RepositoryRoot $file.FullName
    $content = Get-Content -Raw -LiteralPath $file.FullName
    if ($content -match 'os\.(?:Getenv|LookupEnv)\(\s*(?!["''])[^)]+\)' -or $content -match 'GetEnvironmentVariable\(\s*\$[A-Za-z_]') {
      if ($path -notin $registryReaderPaths -or -not (Test-RegistryDrivenEnvironmentReader -Content $content)) {
        $issues.Add((New-PolicyIssue 'RCFG_DYNAMIC_READER_UNBOUNDED' 'dynamic-reader' $path))
      }
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
      if ($reader.kind -eq 'registry' -and ($reader.status -ne 'approved' -or -not (Test-RegistryDrivenEnvironmentReader -Content $content))) {
        $issues.Add((New-PolicyIssue 'RCFG_DYNAMIC_READER_UNBOUNDED' ([string]$key.key) ([string]$reader.path)))
      } elseif ($reader.kind -ne 'registry' -and $content -notmatch "(?<![A-Z0-9_])$([regex]::Escape([string]$key.key))(?![A-Z0-9_])") {
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
      $trackedChanged = @(& git -C $RepositoryRoot diff --name-only $BaseSha -- 2>$null | ForEach-Object { ([string]$_).Replace('\', '/') })
      $trackedStatus = $LASTEXITCODE
      $untrackedChanged = @(& git -C $RepositoryRoot ls-files --others --exclude-standard 2>$null | ForEach-Object { ([string]$_).Replace('\', '/') })
      $untrackedStatus = $LASTEXITCODE
      $changed = @($trackedChanged + $untrackedChanged | Sort-Object -Unique)
      if ($trackedStatus -ne 0 -or $untrackedStatus -ne 0) { $issues.Add((New-PolicyIssue 'GOV_SEMANTIC_DRIFT' 'base-sha-unavailable')) }
      else {
        $apiChanged = 'contracts/api/marketplace-central.openapi.yaml' -in $changed
        # ADR-016 guards the SDK's *types* -- the contract shape the frontend
        # consumes, which lives in the non-test sources under src/. Test files
        # exercise that shape but do not declare it, and package.json or a
        # vitest config is test plumbing; neither carries the drift the
        # same-commit rule exists to prevent.
        $sdkChanged = @($changed | Where-Object {
            $_ -like 'packages/sdk-runtime/src/*' -and
            $_ -notmatch '\.(test|spec)\.(ts|tsx|js|mjs)$'
          }).Count -gt 0
        if ($apiChanged -xor $sdkChanged) { $issues.Add((New-PolicyIssue 'GOV_API_SDK_SPLIT' 'api-sdk-atomicity')) }
      }
    }
  }
  return New-PolicyResult -Violations $issues -BaselineExceptions ($baselines | Sort-Object Id) -Documents $documents -FilesScanned (@($files).Count)
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
