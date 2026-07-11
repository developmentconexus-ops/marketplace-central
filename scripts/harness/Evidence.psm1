Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$script:Pairs = @{
  'fake'='contract'; 'external-dependency-registry'='provisioning'; 'ephemeral-postgres'='integration'; 'dev-invariance'='integration';
  'live-oracle'='live'; 'live-provider'='live'; 'browser'='integration'; 'provider-write'='production-like'; 'cold-gate'='contract'
}
function Test-SafeEvidenceText {
  param([AllowNull()][string]$Value)
  if ($null -eq $Value) { return $true }
  if ($Value -match '(?i)(password|secret|token|credential|connect_string|authorization)\s*[:=]') { return $false }
  if ($Value -match '(?i)\b(?:postgres(?:ql)?|mysql|oracle)://') { return $false }
  if ($Value -match '(?i)\b[A-Z]:[\\/]|^/') { return $false }
  return $true
}
function ConvertTo-SafeArtifactPath {
  param([string]$Path)
  if ([string]::IsNullOrWhiteSpace($Path) -or $Path -match '^(?:[A-Za-z]:[\\/]|/)' -or $Path -match '(^|[\\/])\.\.([\\/]|$)' -or $Path -match '[*?]') { throw 'EVIDENCE_UNSAFE_ARTIFACT_PATH' }
  return $Path.Replace('\\','/')
}
function New-HarnessOutcome {
  [CmdletBinding()]
  param([Parameter(Mandatory)][string]$RunId,[Parameter(Mandatory)][string]$CandidateSha,[Parameter(Mandatory)][string]$Branch,[Parameter(Mandatory)][bool]$Dirty,[Parameter(Mandatory)][object]$Tools,[Parameter(Mandatory)][object[]]$Commands,[string]$PostgresImageIdentity='', [ValidateSet('passed','failed','blocked')][string]$AggregateClassification='passed')
  if ($CandidateSha -notmatch '^[0-9a-f]{40}$') { throw 'EVIDENCE_CANDIDATE_SHA_INVALID' }
  if (-not (Test-SafeEvidenceText $Branch)) { throw 'EVIDENCE_REDACTION_FAILED' }
  $safeCommands = [Collections.Generic.List[object]]::new()
  foreach ($command in @($Commands)) {
    $target = [string]$command.target_label; $klass = [string]$command.evidence_class
    if (-not $script:Pairs.ContainsKey($target) -or $script:Pairs[$target] -ne $klass) { throw "EVIDENCE_TARGET_CLASS_MISMATCH target=$target class=$klass" }
    foreach ($field in @('id','command','stage','reason')) { if (-not (Test-SafeEvidenceText ([string]$command.$field))) { throw 'EVIDENCE_REDACTION_FAILED' } }
    $artifacts = @($command.artifact_paths | ForEach-Object { ConvertTo-SafeArtifactPath ([string]$_) })
    $safeCommands.Add([ordered]@{ id=[string]$command.id; command=[string]$command.command; stage=[string]$command.stage; target_label=$target; evidence_class=$klass; duration_ms=[int]$command.duration_ms; exit_code=[int]$command.exit_code; reason=[string]$command.reason; artifact_paths=$artifacts })
  }
  $toolMap = [ordered]@{}
  foreach ($property in $Tools.PSObject.Properties) { if (-not (Test-SafeEvidenceText ([string]$property.Value))) { throw 'EVIDENCE_REDACTION_FAILED' }; $toolMap[[string]$property.Name] = [string]$property.Value }
  [pscustomobject][ordered]@{ schema_version='1.0'; run_id=$RunId; candidate_sha=$CandidateSha; branch=$Branch; dirty=$Dirty; acceptance_link=$null; tools=$toolMap; commands=@($safeCommands); postgres_image_identity=$PostgresImageIdentity; aggregate_classification=$AggregateClassification }
}
function Write-HarnessOutcome { param([Parameter(Mandatory)][object]$Outcome,[Parameter(Mandatory)][string]$Path) $dir=Split-Path -Parent $Path; if($dir){New-Item -ItemType Directory -Path $dir -Force|Out-Null}; $Outcome|ConvertTo-Json -Depth 40|Set-Content -LiteralPath $Path -Encoding utf8; return $Path }
function Write-HarnessTrace { param([Parameter(Mandatory)][object[]]$Events,[Parameter(Mandatory)][string]$Path) $dir=Split-Path -Parent $Path; if($dir){New-Item -ItemType Directory -Path $dir -Force|Out-Null}; foreach($event in $Events){$event|ConvertTo-Json -Compress -Depth 20|Add-Content -LiteralPath $Path -Encoding utf8}; return $Path }
Export-ModuleMember -Function New-HarnessOutcome, Write-HarnessOutcome, Write-HarnessTrace
