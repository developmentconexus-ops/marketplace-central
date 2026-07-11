$ErrorActionPreference = 'Stop'

$repositoryRoot = (Resolve-Path (Join-Path $PSScriptRoot '../..')).Path
$governanceRoot = Join-Path $repositoryRoot 'contracts/governance'
$schemaRoot = Join-Path $governanceRoot 'schemas'

function Assert-True {
  param([bool]$Condition, [string]$Message)
  if (-not $Condition) { throw $Message }
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

  $contextSchema = Join-Path $schemaRoot 'context-pack.schema.json'
  Assert-True (Test-Path -LiteralPath $contextSchema -PathType Leaf) 'missing schema: contracts/governance/schemas/context-pack.schema.json'

  Assert-True ($modules.modules.Count -eq 11) 'modules registry must contain exactly 11 modules'
  Assert-Unique @($modules.modules.id) 'modules'
  Assert-Unique @($runtime.keys.key) 'runtime keys'
  Assert-Unique @($lanes.lanes.id) 'execution lanes'
  Assert-Unique @($invariants.invariants.id) 'invariants'
  Assert-Unique @($seams.seams.id) 'shared seams'

  $moduleIds = @($modules.modules.id)
  foreach ($module in $modules.modules) {
    foreach ($dependency in $module.dependencies) {
      Assert-True ($dependency -in $moduleIds) "module $($module.id) references unknown dependency $dependency"
    }
  }

  $laneIds = @($lanes.lanes.id)
  Assert-True ((($laneIds | Sort-Object) -join '|') -eq (('browser','integration','live-oracle','live-provider-read','provider-write','unit' | Sort-Object) -join '|')) 'execution lane set is not exact'
  foreach ($key in $runtime.keys) {
    foreach ($lane in $key.allowed_lanes) {
      Assert-True ($lane -in $laneIds) "runtime key $($key.key) references unknown lane $lane"
    }
    if ($key.ContainsKey('alias_for')) {
      Assert-True ($key.alias_for -in @($runtime.keys.key)) "runtime alias $($key.key) references unknown canonical key $($key.alias_for)"
    }
  }

  $unit = @($lanes.lanes | Where-Object id -eq 'unit')[0]
  Assert-True (-not $unit.inherit_parent) 'unit lane must not inherit its parent environment'
  Assert-True ($unit.network -eq 'disabled') 'unit lane network must be disabled'
  Assert-True ($unit.database -eq 'disabled') 'unit lane database must be disabled'
  Assert-True ($unit.target_label -eq 'fake') 'unit lane target must be fake'
  Assert-True ($unit.allowed_runtime_keys.Count -eq 0) 'unit lane must allow no application runtime keys'

  $providerWrite = @($lanes.lanes | Where-Object id -eq 'provider-write')[0]
  $requiredWriteGates = @('actor','idempotency','execute','resolved-link','policy','source-timestamp','before-after-audit')
  Assert-True ((($providerWrite.gates.id | Sort-Object) -join '|') -eq (($requiredWriteGates | Sort-Object) -join '|')) 'provider-write gates are incomplete'
  Assert-True (@($providerWrite.gates | Where-Object { -not $_.required }).Count -eq 0) 'every provider-write gate must be required'

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
    Assert-ExactProperties $exception @('id','rule_id','paths','reason','removal_owner') "invariant exception $($exception.id)"
    Assert-True ($exception.rule_id -in $invariantIds) "invariant exception $($exception.id) references unknown invariant"
    Assert-True (@($exception.paths | Where-Object { $_ -match '[*?\[]' }).Count -eq 0) "invariant exception $($exception.id) is not exact"
  }

  $contextFixture = @{
    schema_version = '1.0'
    context_id = 'ctx-f07-contract'
    feature_path = '.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/F-07-governance-context-compiler'
    base_sha = ('a' * 40)
    objective = 'Create executable governance contracts.'
    observable_done = @('Registries validate.')
    criteria = @(@{ id = 'F07-AC01'; milestone_id = 'M-08-C09'; proof_commands = @('pwsh governance-contracts.tests.ps1') })
    sources = @(@{ path = 'contracts/governance/README.md'; sha256 = ('b' * 64) })
    risk = @{ level = 'L2'; review_policy = 'exclusive-seam-review'; advisory_model = 'none' }
    paths = @{ allowed = @('contracts/governance'); forbidden = @('apps/server_core/internal/modules') }
    shared_seams = @('api-sdk')
    side_effects = @('repository-write')
    commands = @(@{ id = 'governance-contracts'; command = 'pwsh governance-contracts.tests.ps1'; target_label = 'fake'; evidence_type = 'ran' })
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
