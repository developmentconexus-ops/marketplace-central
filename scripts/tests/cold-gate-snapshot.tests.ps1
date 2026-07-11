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
  Assert-True ($harness -match 'New-HarnessChildEnvironment.+cold-provision') 'cold provisioning does not use an isolated child environment'
  Assert-True ($harness -match 'COLD_DOCKER_MISSING') 'missing Docker does not have a stable cold block reason'
  Assert-True ($harness -match 'COLD_IMAGE_IDENTITY_UNRESOLVED') 'unresolved image identity does not block'
  $evidence = Get-Content -Raw (Join-Path $root 'scripts/harness/Evidence.psm1')
  Assert-True ($evidence -match 'Test-Json') 'cold outcome is not schema-validated before persistence'
  Assert-True ($harness -match 'go-mod-download[\s\S]*npm-ci[\s\S]*docker-pull-postgres') 'provisioning order is not frozen Go to npm to Docker'
  Assert-True ($harness -match 'governance-validate[\s\S]*governance-drift[\s\S]*f03-integration') 'required cold validation inventory is incomplete'
  Assert-True ($harness -match 'COLD_CALLER_DIRTY') 'dirty caller does not have stable blocked outcome reason'
  Assert-True ($harness -match '-Dirty \$callerDirty') 'dirty caller state is not recorded in the blocked outcome'
  Write-Output 'PASS cold gate snapshot contract'
  exit 0
} catch { Write-Output "FAIL cold gate snapshot contract: $($_.Exception.Message)"; exit 1 }
