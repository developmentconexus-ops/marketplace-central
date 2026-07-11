$ErrorActionPreference = 'Stop'
$root = (Resolve-Path (Join-Path $PSScriptRoot '../..')).Path
function Assert-True { param([bool]$Condition,[string]$Message) if (-not $Condition) { throw $Message } }
try {
  Assert-True (Test-Path (Join-Path $root 'scripts/harness.ps1')) 'harness entrypoint missing'
  $harness = Get-Content -Raw (Join-Path $root 'scripts/harness.ps1')
  Assert-True ($harness -match "'cold'") 'cold command missing'
  Assert-True ($harness -match 'CandidateSha') 'candidate SHA argument missing'
  Assert-True (Test-Path (Join-Path $root 'contracts/governance/schemas/harness-outcome.schema.json')) 'outcome schema missing'
  Assert-True ($harness -match 'snapshot' -and $harness -match 'GOMODCACHE') 'isolated snapshot/cache contract missing'
  Write-Output 'PASS cold gate snapshot contract'
  exit 0
} catch { Write-Output "FAIL cold gate snapshot contract: $($_.Exception.Message)"; exit 1 }
