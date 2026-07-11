Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$script:Pairs = @{
  'fake'='contract'; 'external-dependency-registry'='provisioning'; 'ephemeral-postgres'='integration'; 'dev-invariance'='integration';
  'live-oracle'='live'; 'live-provider'='live'; 'browser'='integration'; 'provider-write'='production-like'; 'cold-gate'='contract'
}

function Test-SafeEvidenceText {
  param([AllowNull()][string]$Value, [switch]$AllowImageIdentity)
  if ($null -eq $Value) { return $true }
  if ($Value -match '(?i)(?:\b(?:password|secret|token|credential|authorization|connect_string)\b|\b(?:postgres(?:ql)?|mysql|oracle)://|\b[a-z][a-z0-9+.-]*://|\\\\|(?:^|[\s])(?:[A-Za-z]:[\\/]|/)|\b[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}\b)') { return $false }
  if (-not $AllowImageIdentity -and $Value -match '(?i)\bsha256:\S+') { return $false }
  return $true
}

function ConvertTo-SafeArtifactPath {
  param([string]$Path)
  if ([string]::IsNullOrWhiteSpace($Path) -or $Path -match '^(?:[A-Za-z]:[\\/]|/|\\\\)' -or $Path -match '(^|[\\/])\.\.([\\/]|$)' -or $Path -match '[*?]' -or -not (Test-SafeEvidenceText $Path)) { throw 'EVIDENCE_UNSAFE_ARTIFACT_PATH' }
  return $Path.Replace('\\','/')
}

function New-HarnessCommandRecord {
  param([Parameter(Mandatory)][object]$Command)
  $target = [string]$Command.target_label; $klass = [string]$Command.evidence_class
  if (-not $script:Pairs.ContainsKey($target) -or $script:Pairs[$target] -ne $klass) { throw 'EVIDENCE_TARGET_CLASS_MISMATCH' }
  foreach ($field in @('id','command','stage','reason')) {
    if (-not (Test-SafeEvidenceText ([string]$Command.$field))) { throw 'EVIDENCE_REDACTION_FAILED' }
  }
  $artifacts = @($Command.artifact_paths | ForEach-Object { ConvertTo-SafeArtifactPath ([string]$_) })
  return [pscustomobject][ordered]@{
    id=[string]$Command.id; command=[string]$Command.command; stage=[string]$Command.stage; target_label=$target; evidence_class=$klass
    duration_ms=[int]$Command.duration_ms; exit_code=[int]$Command.exit_code; reason=[string]$Command.reason; artifact_paths=$artifacts
  }
}

function Test-HarnessImageIdentity {
  param([AllowEmptyString()][string]$Value)
  return [string]::IsNullOrEmpty($Value) -or $Value -match '^sha256:[0-9a-f]{64}$'
}

function New-HarnessOutcome {
  [CmdletBinding()]
  param(
    [Parameter(Mandatory)][string]$RunId,[Parameter(Mandatory)][string]$CandidateSha,[Parameter(Mandatory)][string]$Branch,
    [Parameter(Mandatory)][bool]$Dirty,[Parameter(Mandatory)][object]$Tools,[Parameter(Mandatory)][object[]]$Commands,
    [AllowEmptyString()][string]$PostgresImageIdentity='',[ValidateSet('passed','failed','blocked')][string]$AggregateClassification='passed',
    [AllowNull()][object]$AcceptanceLink=$null
  )
  if ($CandidateSha -notmatch '^[0-9a-f]{40}$') { throw 'EVIDENCE_CANDIDATE_SHA_INVALID' }
  if ($RunId -notmatch '^[a-z0-9][a-z0-9-]*$' -or -not (Test-SafeEvidenceText $RunId)) { throw 'EVIDENCE_RUN_ID_INVALID' }
  if (-not (Test-SafeEvidenceText $Branch) -or -not (Test-HarnessImageIdentity $PostgresImageIdentity)) { throw 'EVIDENCE_REDACTION_FAILED' }
  if ($null -eq $AcceptanceLink) { $AcceptanceLink = 'unaccepted' }
  if ([string]$AcceptanceLink -notmatch '^(?:unaccepted|F-[0-9]{2}@[0-9a-f]{40})$' -or -not (Test-SafeEvidenceText ([string]$AcceptanceLink))) { throw 'EVIDENCE_ACCEPTANCE_LINK_INVALID' }
  $safeCommands = [Collections.Generic.List[object]]::new()
  foreach ($command in @($Commands)) { $safeCommands.Add((New-HarnessCommandRecord $command)) }
  $ids = @($safeCommands | ForEach-Object id); if (@($ids | Select-Object -Unique).Count -ne $ids.Count) { throw 'EVIDENCE_COMMAND_ID_DUPLICATE' }
  $toolMap = [ordered]@{}
  $toolEntries = if ($Tools -is [System.Collections.IDictionary]) { @($Tools.Keys | ForEach-Object { [pscustomobject]@{ Name = [string]$_; Value = [string]$Tools[$_] } }) } else { @($Tools.PSObject.Properties) }
  foreach ($property in $toolEntries) {
    if (-not (Test-SafeEvidenceText ([string]$property.Name)) -or -not (Test-SafeEvidenceText ([string]$property.Value))) { throw 'EVIDENCE_REDACTION_FAILED' }
    $toolMap[[string]$property.Name] = [string]$property.Value
  }
  return [pscustomobject][ordered]@{ schema_version='1.0'; run_id=$RunId; candidate_sha=$CandidateSha; branch=$Branch; dirty=$Dirty; acceptance_link=$AcceptanceLink; tools=$toolMap; commands=@($safeCommands); postgres_image_identity=$PostgresImageIdentity; aggregate_classification=$AggregateClassification }
}

function Get-HarnessOutcomeProjection {
  param([Parameter(Mandatory)][object]$Outcome)
  return [pscustomobject][ordered]@{candidate_sha=$Outcome.candidate_sha; commands=@($Outcome.commands | ForEach-Object {[ordered]@{id=$_.id;command=$_.command;target_label=$_.target_label;evidence_class=$_.evidence_class;exit_code=$_.exit_code;reason=$_.reason}});postgres_image_identity=$Outcome.postgres_image_identity;aggregate_classification=$Outcome.aggregate_classification}
}

function Compare-HarnessOutcomeProjection {
  param([Parameter(Mandatory)][object]$Left,[Parameter(Mandatory)][object]$Right)
  return ((Get-HarnessOutcomeProjection $Left | ConvertTo-Json -Depth 30 -Compress) -ceq (Get-HarnessOutcomeProjection $Right | ConvertTo-Json -Depth 30 -Compress))
}

function Write-HarnessOutcome {
  param([Parameter(Mandatory)][object]$Outcome,[Parameter(Mandatory)][string]$Path)
  $dir=Split-Path -Parent $Path; if($dir){New-Item -ItemType Directory -Path $dir -Force|Out-Null}
  $Outcome | ConvertTo-Json -Depth 40 | Set-Content -LiteralPath $Path -Encoding utf8
  $schema = Join-Path (Split-Path -Parent (Split-Path -Parent $PSScriptRoot)) 'contracts/governance/schemas/harness-outcome.schema.json'
  if (-not (Test-Json -LiteralPath $Path -SchemaFile $schema -ErrorAction SilentlyContinue)) { Remove-Item -LiteralPath $Path -Force -ErrorAction SilentlyContinue; throw 'EVIDENCE_SCHEMA_INVALID' }
  return $Path
}

function Write-HarnessTrace {
  param([Parameter(Mandatory)][object[]]$Events,[Parameter(Mandatory)][string]$Path)
  $dir=Split-Path -Parent $Path; if($dir){New-Item -ItemType Directory -Path $dir -Force|Out-Null}
  $safe = @($Events | ForEach-Object { New-HarnessCommandRecord $_ })
  foreach($event in $safe){$line=$event|ConvertTo-Json -Compress -Depth 20; if(-not (Test-SafeEvidenceText $line)){throw 'EVIDENCE_REDACTION_FAILED'}; $line|Add-Content -LiteralPath $Path -Encoding utf8}
  return $Path
}

Export-ModuleMember -Function Test-SafeEvidenceText, New-HarnessOutcome, Get-HarnessOutcomeProjection, Compare-HarnessOutcomeProjection, Write-HarnessOutcome, Write-HarnessTrace
