Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

function New-StateResult {
  param([bool]$Passed, [string]$ErrorCode = '', [hashtable]$Data = @{})
  $result = [ordered]@{ Passed=$Passed; ErrorCode=$(if($Passed){$null}else{$ErrorCode}) }
  foreach ($key in $Data.Keys) { $result[$key] = $Data[$key] }
  [pscustomobject]$result
}

function ConvertTo-HarnessStateId {
  param([string]$Value)
  if ([string]::IsNullOrWhiteSpace($Value) -or $Value -notmatch '^[A-Za-z0-9][A-Za-z0-9._-]*$') { throw 'HSTATE_ID_INVALID' }
  $Value
}

function Get-HarnessStatePath {
  param([string]$StateRoot, [string]$Kind, [string]$Id)
  Join-Path (Join-Path $StateRoot $Kind) ((ConvertTo-HarnessStateId $Id) + '.json')
}

function Write-HarnessStateNew {
  param([string]$Path, [hashtable]$Value)
  New-Item -ItemType Directory -Path (Split-Path -Parent $Path) -Force | Out-Null
  $json = $Value | ConvertTo-Json -Depth 30
  try {
    $stream = [IO.File]::Open($Path, [IO.FileMode]::CreateNew, [IO.FileAccess]::Write, [IO.FileShare]::None)
    try { $writer = [IO.StreamWriter]::new($stream, [Text.UTF8Encoding]::new($false)); try { $writer.Write($json) } finally { $writer.Dispose() } } finally { $stream.Dispose() }
    $true
  } catch [IO.IOException] { $false }
}

function New-HarnessLease {
  [CmdletBinding()]
  param([Parameter(Mandatory)][string]$StateRoot, [Parameter(Mandatory)][string]$CheckoutId, [Parameter(Mandatory)][string]$WriterId, [string[]]$SharedSeam = @())
  try {
    $checkout = ConvertTo-HarnessStateId $CheckoutId; $writer = ConvertTo-HarnessStateId $WriterId
    $seams = @($SharedSeam | ForEach-Object { ConvertTo-HarnessStateId $_ } | Sort-Object -Unique)
  } catch { return New-StateResult $false $_.Exception.Message }
  $record = [ordered]@{ checkout_id=$checkout; writer_id=$writer; shared_seams=$seams; acquired_at=[DateTime]::UtcNow.ToString('o'); disposition='active' }
  $paths = @((Get-HarnessStatePath $StateRoot 'leases' "checkout-$checkout")) + @($seams | ForEach-Object { Get-HarnessStatePath $StateRoot 'leases' "seam-$_" })
  $created = [Collections.Generic.List[string]]::new()
  foreach ($path in $paths) {
    if (-not (Write-HarnessStateNew -Path $path -Value $record)) {
      foreach ($owned in $created) { Remove-Item -LiteralPath $owned -Force -ErrorAction SilentlyContinue }
      return New-StateResult $false 'HSTATE_LEASE_CONFLICT' @{ ConflictPath=$path }
    }
    $created.Add($path)
  }
  New-StateResult $true '' @{ LeasePaths=@($created); CheckoutId=$checkout; WriterId=$writer }
}

function Get-HarnessLeaseRecovery {
  [CmdletBinding()]
  param([Parameter(Mandatory)][string]$StateRoot, [Parameter(Mandatory)][string]$CheckoutId)
  try { $path = Get-HarnessStatePath $StateRoot 'leases' "checkout-$CheckoutId" } catch { return New-StateResult $false $_.Exception.Message }
  if (-not (Test-Path -LiteralPath $path -PathType Leaf)) { return New-StateResult $false 'HSTATE_LEASE_MISSING' }
  $lease = Get-Content -Raw -LiteralPath $path | ConvertFrom-Json -Depth 30
  New-StateResult $true '' @{ Disposition='explicit-disposition-required'; CheckoutId=[string]$lease.checkout_id; WriterId=[string]$lease.writer_id; LeasePath=$path }
}

function New-HarnessCheckpoint {
  [CmdletBinding()]
  param([Parameter(Mandatory)][string]$StateRoot, [Parameter(Mandatory)][string]$CheckpointId, [Parameter(Mandatory)][string]$FeaturePath, [Parameter(Mandatory)][string]$FeatureState, [Parameter(Mandatory)][string]$BaseSha, [Parameter(Mandatory)][string]$CurrentSha, [string[]]$ChangedPath=@(), [string[]]$EvidencePath=@(), [string]$Blocker='', [Parameter(Mandatory)][string]$NextAction, [string]$TaskCorrelationId='')
  try { $id=ConvertTo-HarnessStateId $CheckpointId } catch { return New-StateResult $false $_.Exception.Message }
  if ($BaseSha -notmatch '^[0-9a-f]{40}$' -or $CurrentSha -notmatch '^[0-9a-f]{40}$' -or [string]::IsNullOrWhiteSpace($NextAction)) { return New-StateResult $false 'HSTATE_CHECKPOINT_INVALID' }
  $path = Get-HarnessStatePath $StateRoot 'checkpoints' $id
  $record = [ordered]@{ checkpoint_id=$id; feature_path=$FeaturePath; feature_state=$FeatureState; base_sha=$BaseSha; current_sha=$CurrentSha; changed_paths=@($ChangedPath | Sort-Object -Unique); evidence_paths=@($EvidencePath | Sort-Object -Unique); blocker=$Blocker; next_action=$NextAction; task_correlation_id=$TaskCorrelationId; recorded_at=[DateTime]::UtcNow.ToString('o') }
  if (-not (Write-HarnessStateNew -Path $path -Value $record)) { return New-StateResult $false 'HSTATE_CHECKPOINT_EXISTS' }
  New-StateResult $true '' @{ CheckpointId=$id; Path=$path }
}

function Get-HarnessResume {
  [CmdletBinding()]
  param([Parameter(Mandatory)][string]$StateRoot, [Parameter(Mandatory)][string]$CheckpointId)
  try { $path = Get-HarnessStatePath $StateRoot 'checkpoints' $CheckpointId } catch { return New-StateResult $false $_.Exception.Message }
  if (-not (Test-Path -LiteralPath $path -PathType Leaf)) { return New-StateResult $false 'HSTATE_CHECKPOINT_MISSING' }
  $record = Get-Content -Raw -LiteralPath $path | ConvertFrom-Json -Depth 30
  New-StateResult $true '' @{ CheckpointId=[string]$record.checkpoint_id; NextAction=[string]$record.next_action; FeaturePath=[string]$record.feature_path; FeatureState=[string]$record.feature_state; BaseSha=[string]$record.base_sha; CurrentSha=[string]$record.current_sha; TaskCorrelationId=[string]$record.task_correlation_id; EvidencePaths=@($record.evidence_paths); ChangedPaths=@($record.changed_paths) }
}

function Get-HarnessRiskPolicy {
  [CmdletBinding()]
  param([Parameter(Mandatory)][ValidateSet('L0','L1','L2','L3')][string]$RiskLevel)
  $policies = @{
    L0=@{review='self-review';command_ids=@('governance-contracts');real_qa='not-required'}
    L1=@{review='combined-final-review';command_ids=@('context-compiler-tests','governance-contracts');real_qa='not-required'}
    L2=@{review='fixed-commit-review';command_ids=@('context-compiler-tests','governance-contracts');real_qa='contract-selected'}
    L3=@{review='independent-safety-and-quality';command_ids=@('context-compiler-tests','governance-contracts');real_qa='named-real-target-required'}
  }
  [pscustomobject]$policies[$RiskLevel]
}

function Test-HarnessCompactHandoff {
  [CmdletBinding()]
  param([Parameter(Mandatory)][hashtable]$Value)
  $required=@('status','checkpoint_id','commit','changed_paths','evidence','review','blockers','next')
  if (@($required | Where-Object { -not $Value.ContainsKey($_) -or $null -eq $Value[$_] }).Count -gt 0) { return New-StateResult $false 'HSTATE_HANDOFF_INVALID' }
  if ([string]$Value.checkpoint_id -notmatch '^[A-Za-z0-9][A-Za-z0-9._-]*$' -or [string]$Value.commit -notmatch '^[0-9a-f]{40}$') { return New-StateResult $false 'HSTATE_HANDOFF_INVALID' }
  New-StateResult $true '' @{ CheckpointId=[string]$Value.checkpoint_id }
}

Export-ModuleMember -Function New-HarnessLease, Get-HarnessLeaseRecovery, New-HarnessCheckpoint, Get-HarnessResume, Get-HarnessRiskPolicy, Test-HarnessCompactHandoff
