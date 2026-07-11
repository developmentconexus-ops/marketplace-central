$ErrorActionPreference = 'Stop'
$root = (Resolve-Path (Join-Path $PSScriptRoot '../..')).Path
function Assert-True { param([bool]$Condition,[string]$Message) if (-not $Condition) { throw $Message } }
try {
  $harness = Get-Content -Raw (Join-Path $root 'scripts/harness.ps1')
  Assert-True ($harness -match "'f03-integration'[\s\S]+'ephemeral-postgres'[\s\S]+'integration'") 'F-03 real integration is not target/evidence classified'
  Assert-True ($harness -match "'-Command','integration'") 'cold gate does not invoke the F-03 integration lane'
  Assert-True ($harness -match 'context-current') 'cold gate does not current-validate its candidate context pack'
  Assert-True ($harness -match "'cold-gate' 'contract'") 'current context record is not classified as top-level cold-gate evidence'
  Assert-True ($harness -match '--pull=never') 'cold inventory does not assert F-03 pull-never contract'
  Assert-True ($harness -match 'duration_ms') 'cold records do not capture actual duration'
  Assert-True ($harness -match 'Compare-HarnessOutcomeProjection') 'cold outcome comparison is not automatic and recorded'
  Write-Output 'PASS cold gate integration contract target=ephemeral-postgres'
  exit 0
} catch { Write-Output "FAIL cold gate integration contract: $($_.Exception.Message)"; exit 1 }
