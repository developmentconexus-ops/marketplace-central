$ErrorActionPreference = 'Stop'

$repositoryRoot = (Resolve-Path (Join-Path $PSScriptRoot '../..')).Path
$governanceRoot = Join-Path $repositoryRoot 'contracts/governance'
$schemaRoot = Join-Path $governanceRoot 'schemas'

function Assert-True {
  param([bool]$Condition, [string]$Message)
  if (-not $Condition) { throw $Message }
}

$retiredRuntimeKeys = @(
  ('MC_' + 'MIGRATIONS_DIR'),
  ('MPC_' + 'PRODUCT_LINKS_' + 'LIVE_TEST'),
  ('MPC_' + 'PRODUCT_LINKS_' + 'POSTGRES_URL'),
  ('MPC_' + 'PRODUCT_LINKS_' + 'INSTALLATION_ID')
)
# git grep, not rg: git is the one tool this repository cannot run without,
# while a missing rg either throws or -- suppressed -- returns an empty match
# list and turns this sweep into a vacuous pass. Tracked files are also the
# right universe here; scripts/.runs is untracked and needs no exclusion.
foreach ($retiredKey in $retiredRuntimeKeys) {
  $references = @(& git grep --fixed-strings -l -- $retiredKey `
      apps/server_core scripts docker contracts/governance docker-compose.yml ':!*.md' 2>$null)
  Assert-True ($references.Count -eq 0) "retired runtime key still referenced: $retiredKey"
}

function Assert-Unique {
  param([object[]]$Values, [string]$Label)
  $duplicates = @($Values | Group-Object | Where-Object Count -gt 1 | ForEach-Object Name)
  Assert-True ($duplicates.Count -eq 0) "$Label contains duplicate IDs: $($duplicates -join ', ')"
}

function Assert-ExactProperties {
  param([hashtable]$Value, [string[]]$Expected, [string]$Label)
  $actual = @($Value.Keys | Sort-Object)
  $wanted = @($Expected | Sort-Object)
  Assert-True (($actual -join '|') -eq ($wanted -join '|')) "$Label has unexpected shape: $($actual -join ', ')"
}

function Test-RuntimeExceptionReferences {
  param([hashtable]$Registry)
  $ruleIds = @($Registry.exception_rules)
  return @($Registry.temporary_exceptions | Where-Object { $_.rule_id -notin $ruleIds }).Count -eq 0
}

function Read-GovernanceDocument {
  param([string]$Name)
  $documentPath = Join-Path $governanceRoot "$Name.json"
  $schemaPath = Join-Path $schemaRoot "$Name.schema.json"
  Assert-True (Test-Path -LiteralPath $documentPath -PathType Leaf) "missing registry: contracts/governance/$Name.json"
  Assert-True (Test-Path -LiteralPath $schemaPath -PathType Leaf) "missing schema: contracts/governance/schemas/$Name.schema.json"
  Assert-True (Test-Json -LiteralPath $documentPath -SchemaFile $schemaPath -ErrorAction Stop) "schema rejected $Name.json"

  $document = Get-Content -Raw -LiteralPath $documentPath | ConvertFrom-Json -AsHashtable
  $document.unknown_contract_property = $true
  $invalidPath = Join-Path ([IO.Path]::GetTempPath()) "governance-$Name-$([guid]::NewGuid().ToString('N')).json"
  try {
    $document | ConvertTo-Json -Depth 30 | Set-Content -LiteralPath $invalidPath -Encoding utf8
    Assert-True (-not (Test-Json -LiteralPath $invalidPath -SchemaFile $schemaPath -ErrorAction SilentlyContinue)) "$Name schema accepted an unknown root property"
  }
  finally {
    Remove-Item -LiteralPath $invalidPath -Force -ErrorAction SilentlyContinue
  }

  $document.Remove('unknown_contract_property')
  return $document
}

try {
  $modules = Read-GovernanceDocument 'modules'
  $runtime = Read-GovernanceDocument 'runtime-config'
  $lanes = Read-GovernanceDocument 'execution-lanes'
  $invariants = Read-GovernanceDocument 'invariants'
  $seams = Read-GovernanceDocument 'shared-seams'
  $knowledgeRoutes = Read-GovernanceDocument 'knowledge-routes'

  $contextSchema = Join-Path $schemaRoot 'context-pack.schema.json'
  Assert-True (Test-Path -LiteralPath $contextSchema -PathType Leaf) 'missing schema: contracts/governance/schemas/context-pack.schema.json'
  $workContractSchema = Join-Path $schemaRoot 'feature-work-contract.schema.json'
  $currentWorkContract = @{
    schema_version='1.0'; feature_id='F-10'; required_sources=@('scripts/harness.ps1'); allowed_paths=@('scripts/harness.ps1'); forbidden_paths=@('apps/**')
    side_effects=@{allowed=@('repository-write');forbidden=@('provider-write')}; commands=@(@{id='impact-probe-one';command_id='impact-probe-one';lane_id='unit';expected_exit_code=0})
    criteria=@(@{id='F10-AC01';milestone_criterion_id='M-08-C17';command_ids=@('impact-probe-one')}); stop_conditions=@(@{code='impact-stop';condition='stop'}) ; retry_budget=@{max_correction_attempts=1}; handoff_fields=@('status')
  }
  $currentWorkContractPath = Join-Path ([IO.Path]::GetTempPath()) "feature-work-contract-current-$([guid]::NewGuid().ToString('N')).json"
  try { $currentWorkContract | ConvertTo-Json -Depth 20 | Set-Content -LiteralPath $currentWorkContractPath -Encoding utf8; Assert-True (Test-Json -LiteralPath $currentWorkContractPath -SchemaFile $workContractSchema -ErrorAction Stop) 'work contract schema rejected registered command ID' }
  finally { Remove-Item -LiteralPath $currentWorkContractPath -Force -ErrorAction SilentlyContinue }
  Assert-True (Test-Path -LiteralPath $workContractSchema -PathType Leaf) 'missing schema: contracts/governance/schemas/feature-work-contract.schema.json'
  $f10PlanPath = Join-Path $repositoryRoot '.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/F-10-pragmatic-harness-cutover/plan.md'
  $f10Plan = Get-Content -Raw -LiteralPath $f10PlanPath
  $workContractBlocks = @([regex]::Matches($f10Plan, '(?ms)^## Machine Work Contract\s*\r?\n\s*```json\s*\r?\n(?<json>.*?)\r?\n```\s*$'))
  Assert-True ($workContractBlocks.Count -eq 1) 'F-10 plan must contain exactly one machine work contract JSON block'
  $workContractPath = Join-Path ([IO.Path]::GetTempPath()) "feature-work-contract-$([guid]::NewGuid().ToString('N')).json"
  try {
    Set-Content -LiteralPath $workContractPath -Value $workContractBlocks[0].Groups['json'].Value -Encoding utf8
    Assert-True (Test-Json -LiteralPath $workContractPath -SchemaFile $workContractSchema -ErrorAction Stop) 'F-10 machine work contract failed schema validation'
    $invalidWorkContract = Get-Content -Raw -LiteralPath $workContractPath | ConvertFrom-Json -AsHashtable
    $invalidWorkContract.unknown_contract_property = $true
    $invalidWorkContract | ConvertTo-Json -Depth 30 | Set-Content -LiteralPath $workContractPath -Encoding utf8
    Assert-True (-not (Test-Json -LiteralPath $workContractPath -SchemaFile $workContractSchema -ErrorAction SilentlyContinue)) 'work contract schema accepted an unknown root property'
  }
  finally {
    Remove-Item -LiteralPath $workContractPath -Force -ErrorAction SilentlyContinue
  }

  # The registry key is the pair (kind, id): a module and a context may share an
  # id (catalog does). The honest count is the tree itself -- every directory
  # under internal/modules and internal/contexts has exactly one entry of its
  # kind, and no entry names a directory that does not exist. A hardcoded count
  # rotted silently from 11 to 21 before this assertion replaced it.
  $moduleEntries = @($modules.modules | Where-Object { -not $_.ContainsKey('kind') -or $_.kind -eq 'module' })
  $contextEntries = @($modules.modules | Where-Object { $_.ContainsKey('kind') -and $_.kind -eq 'context' })
  Assert-True (($moduleEntries.Count + $contextEntries.Count) -eq $modules.modules.Count) 'registry entry with an unknown kind'
  $moduleDirs = @(Get-ChildItem -Directory (Join-Path $repositoryRoot 'apps/server_core/internal/modules') | ForEach-Object Name | Sort-Object)
  $contextDirs = @(Get-ChildItem -Directory (Join-Path $repositoryRoot 'apps/server_core/internal/contexts') | ForEach-Object Name | Sort-Object)
  Assert-True ((@($moduleEntries.id | Sort-Object) -join '|') -eq ($moduleDirs -join '|')) `
    "module registry ids do not mirror internal/modules: registry=[$(@($moduleEntries.id | Sort-Object) -join ',')] tree=[$($moduleDirs -join ',')]"
  Assert-True ((@($contextEntries.id | Sort-Object) -join '|') -eq ($contextDirs -join '|')) `
    "context registry ids do not mirror internal/contexts: registry=[$(@($contextEntries.id | Sort-Object) -join ',')] tree=[$($contextDirs -join ',')]"
  Assert-Unique @($moduleEntries.id) 'modules'
  Assert-Unique @($contextEntries.id) 'contexts'
  Assert-Unique @($runtime.keys.key) 'runtime keys'
  Assert-Unique @($lanes.lanes.id) 'execution lanes'
  Assert-Unique @($invariants.invariants.id) 'invariants'
  Assert-Unique @($seams.seams.id) 'shared seams'
  Assert-Unique @($knowledgeRoutes.routes.id) 'knowledge routes'

  $moduleIds = @($modules.modules.id)
  foreach ($module in $modules.modules) {
    foreach ($dependency in $module.dependencies) {
      Assert-True ($dependency -in $moduleIds) "module $($module.id) references unknown dependency $dependency"
    }
  }

  $laneIds = @($lanes.lanes.id)
  Assert-True ((($laneIds | Sort-Object) -join '|') -eq (('browser','dev-invariance','integration','live-oracle','live-provider-read','provider-write','unit' | Sort-Object) -join '|')) 'execution lane set is not exact'
  foreach ($key in $runtime.keys) {
    foreach ($lane in $key.allowed_lanes) {
      Assert-True ($lane -in $laneIds) "runtime key $($key.key) references unknown lane $lane"
    }
    if ($key.ContainsKey('alias_for')) {
      Assert-True ($key.alias_for -in @($runtime.keys.key)) "runtime alias $($key.key) references unknown canonical key $($key.alias_for)"
    }
    foreach ($reader in $key.readers) {
      $expectedProperties = if ($reader.status -eq 'temporary_exception') { @('path','kind','status','exception_id') } else { @('path','kind','status') }
      Assert-ExactProperties $reader $expectedProperties "runtime reader $($key.key):$($reader.path)"
      # Only 'direct' is exceptional by definition -- production code reaching
      # past the typed config. A 'test' reader can be approved: the schema
      # permits it explicitly (its kind=registry special-case shows per-kind
      # status was considered), and the registry carries an approved test
      # reader merged through review. Contract outranks this test.
      if ($reader.kind -eq 'direct') {
        Assert-True ($reader.status -eq 'temporary_exception') "direct runtime reader is not exceptional: $($key.key):$($reader.path)"
      }
      if ($reader.status -eq 'temporary_exception') {
        $linkedException = @($runtime.temporary_exceptions | Where-Object id -eq $reader.exception_id)
        Assert-True ($linkedException.Count -eq 1) "runtime reader exception link is unresolved: $($key.key):$($reader.path)"
        Assert-True ($linkedException[0].key -eq $key.key -and $linkedException[0].path -eq $reader.path) "runtime reader exception link is not exact: $($key.key):$($reader.path)"
      }
    }
  }
  $referencedRuntimeExceptions = @($runtime.keys.readers | Where-Object status -eq 'temporary_exception' | ForEach-Object exception_id)
  foreach ($exception in $runtime.temporary_exceptions) {
    Assert-True (@($referencedRuntimeExceptions | Where-Object { $_ -eq $exception.id }).Count -eq 1) "runtime exception is not referenced exactly once: $($exception.id)"
  }

  $registryDrivenKeys = @(
    $lanes.lanes |
      ForEach-Object allowed_runtime_keys |
      Sort-Object -Unique
  )
  foreach ($keyName in $registryDrivenKeys) {
    $keyRecord = @($runtime.keys | Where-Object key -eq $keyName)[0]
    $reader = @($keyRecord.readers | Where-Object {
      $_.path -eq 'scripts/harness/Environment.psm1' -and $_.kind -eq 'registry' -and $_.status -eq 'approved'
    })
    Assert-True ($reader.Count -eq 1) "registry-driven environment reader missing: $keyName"
  }
  $declaredRegistryKeys = @(
    $runtime.keys |
      Where-Object { @($_.readers | Where-Object kind -eq 'registry').Count -gt 0 } |
      ForEach-Object key |
      Sort-Object -Unique
  )
  Assert-True (($declaredRegistryKeys -join '|') -eq ($registryDrivenKeys -join '|')) 'registry readers do not exactly cover all lane runtime keys'

  $invalidRegistryReader = $runtime | ConvertTo-Json -Depth 30 | ConvertFrom-Json -AsHashtable
  $invalidReader = @($invalidRegistryReader.keys | ForEach-Object readers | Where-Object kind -eq 'registry')[0]
  $invalidReader.status = 'temporary_exception'
  $invalidReader.exception_id = 'invalid-registry-reader'
  $invalidRegistryPath = Join-Path ([IO.Path]::GetTempPath()) "runtime-registry-reader-$([guid]::NewGuid().ToString('N')).json"
  try {
    $invalidRegistryReader | ConvertTo-Json -Depth 30 | Set-Content -LiteralPath $invalidRegistryPath -Encoding utf8
    Assert-True (-not (Test-Json -LiteralPath $invalidRegistryPath -SchemaFile (Join-Path $schemaRoot 'runtime-config.schema.json') -ErrorAction SilentlyContinue)) 'runtime schema accepted a non-approved registry reader'
  }
  finally {
    Remove-Item -LiteralPath $invalidRegistryPath -Force -ErrorAction SilentlyContinue
  }

  $expectedReaders = @(
    @{ key='SERVER_ADDR'; path='docker/dev/backend-entrypoint.sh'; kind='edge'; status='approved' },
    @{ key='API_PORT'; path='docker/dev/backend-entrypoint.sh'; kind='edge'; status='approved' },
    @{ key='MPC_ORACLE_LIB_DIR'; path='docker/dev/backend-entrypoint.sh'; kind='edge'; status='approved' },
    @{ key='MPC_ORACLE_USERNAME'; path='docker/dev/backend-entrypoint.sh'; kind='edge'; status='approved' },
    @{ key='MPC_ORACLE_PASSWORD'; path='docker/dev/backend-entrypoint.sh'; kind='edge'; status='approved' },
    @{ key='MPC_ORACLE_CONNECT_STRING'; path='docker/dev/backend-entrypoint.sh'; kind='edge'; status='approved' },
    @{ key='SANKHYA_ORACLE_USER'; path='docker/dev/backend-entrypoint.sh'; kind='edge'; status='approved' },
    @{ key='SANKHYA_ORACLE_PASSWORD'; path='docker/dev/backend-entrypoint.sh'; kind='edge'; status='approved' },
    @{ key='SANKHYA_ORACLE_HOST'; path='docker/dev/backend-entrypoint.sh'; kind='edge'; status='approved' },
    @{ key='SANKHYA_ORACLE_PORT'; path='docker/dev/backend-entrypoint.sh'; kind='edge'; status='approved' },
    @{ key='SANKHYA_ORACLE_SERVICE_NAME'; path='docker/dev/backend-entrypoint.sh'; kind='edge'; status='approved' },
    @{ key='MPC_OAUTH_REDIRECT_URI'; path='docker/dev/ngrok-entrypoint.sh'; kind='edge'; status='approved' },
    @{ key='MPC_TEST_DATABASE_URL'; path='apps/server_core/internal/testsupport/postgres/target.go'; kind='typed'; status='approved' },
    @{ key='MPC_ORACLE_LIVE_TEST'; path='apps/server_core/internal/modules/internal_read/adapters/oracle/reader_live_test.go'; kind='test'; status='temporary_exception'; owner='M-08/F-04' }
  )
  foreach ($expected in $expectedReaders) {
    $keyRecord = @($runtime.keys | Where-Object key -eq $expected.key)[0]
    $reader = @($keyRecord.readers | Where-Object { $_.path -eq $expected.path -and $_.kind -eq $expected.kind -and $_.status -eq $expected.status })
    Assert-True ($reader.Count -eq 1) "runtime reader baseline missing: $($expected.key):$($expected.path)"
    if ($expected.status -eq 'temporary_exception') {
      $exception = @($runtime.temporary_exceptions | Where-Object id -eq $reader[0].exception_id)
      Assert-True ($exception.Count -eq 1) "runtime reader exception missing: $($expected.key):$($expected.path)"
      Assert-True ($exception[0].key -eq $expected.key -and $exception[0].path -eq $expected.path) "runtime reader exception is not exact: $($expected.key):$($expected.path)"
      Assert-True ($exception[0].removal_owner -eq $expected.owner) "runtime reader exception owner mismatch: $($expected.key):$($expected.path)"
    }
  }

  Assert-True (Test-RuntimeExceptionReferences $runtime) 'runtime exception has a dangling rule_id'
  $runtimeWithDanglingRule = $runtime | ConvertTo-Json -Depth 30 | ConvertFrom-Json -AsHashtable
  $runtimeWithDanglingRule.temporary_exceptions[0].rule_id = 'missing-runtime-rule'
  Assert-True (-not (Test-RuntimeExceptionReferences $runtimeWithDanglingRule)) 'runtime reference validation accepted a dangling rule_id'

  $unit = @($lanes.lanes | Where-Object id -eq 'unit')[0]
  Assert-True (-not $unit.inherit_parent) 'unit lane must not inherit its parent environment'
  Assert-True ($unit.network -eq 'disabled') 'unit lane network must be disabled'
  Assert-True ($unit.database -eq 'disabled') 'unit lane database must be disabled'
  Assert-True ($unit.target_label -eq 'fake') 'unit lane target must be fake'
  Assert-True ($unit.allowed_runtime_keys.Count -eq 0) 'unit lane must allow no application runtime keys'

  $liveOracle = @($lanes.lanes | Where-Object id -eq 'live-oracle')[0]
  Assert-True ((($liveOracle.side_effects | Sort-Object) -join '|') -eq 'database-read|session-temporary-write') 'live-oracle lane must declare database-read plus the session-temporary-write envelope (EXPLAIN PLAN PLAN_TABLE) and nothing else'

  $providerWrite = @($lanes.lanes | Where-Object id -eq 'provider-write')[0]
  $requiredWriteGates = @('actor','idempotency','execute','resolved-link','policy','source-timestamp','before-after-audit')
  Assert-True ((($providerWrite.gates.id | Sort-Object) -join '|') -eq (($requiredWriteGates | Sort-Object) -join '|')) 'provider-write gates are incomplete'
  Assert-True (@($providerWrite.gates | Where-Object { -not $_.required }).Count -eq 0) 'every provider-write gate must be required'

  $devInvariance = @($lanes.lanes | Where-Object id -eq 'dev-invariance')[0]
  Assert-True ($devInvariance.database -eq 'ephemeral-postgres' -and $devInvariance.target_label -eq 'dev-invariance') 'dev invariance target is misclassified'
  Assert-True (($devInvariance.side_effects -join '|') -eq 'database-read') 'dev invariance must be read-only against dev'
  Assert-True (@($devInvariance.allowed_runtime_keys).Count -eq 0) 'dev invariance cannot inherit a database target'

  $invariantIds = @($invariants.invariants.id)
  foreach ($exception in @($modules.temporary_exceptions)) {
    Assert-ExactProperties $exception @('id','rule_id','source_module','target_module','target_layer','path','import_path','reason','removal_owner') "module exception $($exception.id)"
    Assert-True ($exception.rule_id -in $invariantIds) "module exception $($exception.id) references unknown invariant"
    Assert-True (($exception.path + $exception.import_path) -notmatch '[*?\[]') "module exception $($exception.id) is not exact"
  }
  foreach ($exception in @($runtime.temporary_exceptions)) {
    Assert-ExactProperties $exception @('id','rule_id','key','path','reason','removal_owner') "runtime exception $($exception.id)"
    Assert-True ($exception.key -in @($runtime.keys.key)) "runtime exception $($exception.id) references unknown key"
    Assert-True ($exception.path -notmatch '[*?\[]') "runtime exception $($exception.id) is not exact"
  }
  foreach ($exception in @($invariants.temporary_exceptions)) {
    Assert-ExactProperties $exception @('id','rule_id','paths','occurrences','reason','removal_owner') "invariant exception $($exception.id)"
    Assert-True ($exception.rule_id -in $invariantIds) "invariant exception $($exception.id) references unknown invariant"
    Assert-True (@($exception.paths | Where-Object { $_ -match '[*?\[]' }).Count -eq 0) "invariant exception $($exception.id) is not exact"
  }

  # The occurrence-coverage assertions stood here. They existed because
  # `production-panic` was the only rule whose exceptions carried per-site
  # occurrence records -- a path, a symbol, a verbatim source fingerprint and a
  # count -- and a duplicated or stale record silently widened the exception.
  # The rule is now forbidigo under the lint-go ratchet and no exception in this
  # registry carries occurrences any more, so there is nothing left to be exact
  # about. The assertion below proves that rather than assuming it.
  Assert-True (@($invariants.temporary_exceptions | Where-Object { @($_.occurrences).Count -gt 0 }).Count -eq 0) `
    'an invariant exception carries occurrences; the occurrence-exactness assertions were removed with production-panic and would not cover it'

  $contextFixture = @{
    schema_version = '1.0'
    context_id = 'ctx-f07-contract'
    feature_path = '.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/F-07-governance-context-compiler'
    base_sha = ('a' * 40)
    objective = 'Create executable governance contracts.'
    observable_done = @('Registries validate.')
    criteria = @(@{ id = 'F07-AC01'; milestone_id = 'M-08-C09'; proof_commands = @('pwsh governance-contracts.tests.ps1') })
    sources = @(@{ path = 'contracts/governance/README.md'; sha256 = ('b' * 64); selector = 'governance overview'; reason = 'fixture source'; estimated_tokens = 120 })
    risk = @{ level = 'L2'; review_policy = 'exclusive-seam-review'; advisory_model = 'none' }
    paths = @{ allowed = @('contracts/governance'); forbidden = @('apps/server_core/internal/modules') }
    shared_seams = @('api-sdk')
    side_effects = @('repository-write')
    commands = @(@{ id = 'governance-contracts'; command_id = 'governance-contracts'; target_label = 'fake'; evidence_type = 'ran' })
    stop_conditions = @('source-hash-mismatch')
    retry_budget = 1
    handoff = @{ target = 'Milestone Orchestrator'; reason = 'Phase complete.' }
    estimated_input_tokens = 500
  }
  $contextPath = Join-Path ([IO.Path]::GetTempPath()) "governance-context-$([guid]::NewGuid().ToString('N')).json"
  try {
    $contextFixture | ConvertTo-Json -Depth 30 | Set-Content -LiteralPath $contextPath -Encoding utf8
    Assert-True (Test-Json -LiteralPath $contextPath -SchemaFile $contextSchema -ErrorAction Stop) 'context fixture failed schema validation'
  }
  finally {
    Remove-Item -LiteralPath $contextPath -Force -ErrorAction SilentlyContinue
  }

  Write-Output 'PASS governance contract tests'
  exit 0
}
catch {
  Write-Output "FAIL governance contract tests: $($_.Exception.Message)"
  exit 1
}
