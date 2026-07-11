$ErrorActionPreference = 'Stop'
$root = (Resolve-Path (Join-Path $PSScriptRoot '../..')).Path
function Assert-True { param([bool]$Condition,[string]$Message) if (-not $Condition) { throw $Message } }
try {
  $schema = Join-Path $root 'contracts/governance/schemas/harness-outcome.schema.json'
  Assert-True (Test-Path -LiteralPath $schema -PathType Leaf) 'missing harness outcome schema'
  $modulePath = Join-Path $root 'scripts/harness/Evidence.psm1'
  Assert-True (Test-Path -LiteralPath $modulePath -PathType Leaf) 'missing evidence module'
  Import-Module $modulePath -Force
  $run = Join-Path $root ('scripts/.runs/test-' + [guid]::NewGuid().ToString('N'))
  New-Item -ItemType Directory -Path $run -Force | Out-Null
  try {
    $commands = @(
      [ordered]@{ id='governance'; command='pwsh governance'; stage='validation'; target_label='fake'; evidence_class='contract'; duration_ms=1; exit_code=0; reason='passed'; artifact_paths=@('contracts/governance') },
      [ordered]@{ id='cold-provision'; command='go mod download'; stage='provisioning'; target_label='external-dependency-registry'; evidence_class='provisioning'; duration_ms=1; exit_code=0; reason='passed'; artifact_paths=@() }
    )
    $manifest = New-HarnessOutcome -RunId 'a1' -CandidateSha ('a' * 40) -Branch 'main' -Dirty $false -Tools ([ordered]@{ go='go version'; npm='npm --version' }) -Commands $commands -PostgresImageIdentity 'sha256:abc' -AggregateClassification 'passed'
    $path = Join-Path $run 'outcome.json'
    Write-HarnessOutcome -Outcome $manifest -Path $path
    Assert-True (Test-Json -LiteralPath $path -SchemaFile $schema -ErrorAction Stop) 'outcome failed schema validation'
    $json = Get-Content -Raw $path | ConvertFrom-Json -Depth 30
    Assert-True ($json.commands.Count -eq 2) 'command inventory not preserved'
    Assert-True ($json.acceptance_link -eq $null) 'acceptance link must be unset'
    $unsafe = $commands + @([ordered]@{ id='bad'; command='echo SECRET'; stage='validation'; target_label='fake'; evidence_class='contract'; duration_ms=1; exit_code=0; reason=''; artifact_paths=@('C:\secret') })
    $threw = $false; try { New-HarnessOutcome -RunId 'a2' -CandidateSha ('a' * 40) -Branch 'main' -Dirty $false -Tools @{} -Commands $unsafe -PostgresImageIdentity '' -AggregateClassification 'failed' | Out-Null } catch { $threw = $true }
    Assert-True $threw 'unsafe evidence was accepted'
    Write-Output 'PASS cold gate evidence contract'
  } finally { Remove-Item -LiteralPath $run -Recurse -Force -ErrorAction SilentlyContinue }
  exit 0
} catch { Write-Output "FAIL cold gate evidence contract: $($_.Exception.Message)"; exit 1 }
