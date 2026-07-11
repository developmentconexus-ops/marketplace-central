$ErrorActionPreference = 'Stop'
$root = (Resolve-Path (Join-Path $PSScriptRoot '../..')).Path
function Assert-True { param([bool]$Condition,[string]$Message) if (-not $Condition) { throw $Message } }
try {
  $schema = Join-Path $root 'contracts/governance/schemas/harness-outcome.schema.json'
  Assert-True (Test-Path -LiteralPath $schema -PathType Leaf) 'missing harness outcome schema'
  $modulePath = Join-Path $root 'scripts/harness/Evidence.psm1'
  Assert-True (Test-Path -LiteralPath $modulePath -PathType Leaf) 'missing evidence module'
  Import-Module $modulePath -Force
  $processException = [OperationCanceledException]::new('cancelled https://example.test/token C:\\secret token=super-secret')
  $simulatedRecord = $null
  try { throw $processException } catch { $simulatedRecord = [ordered]@{ id='go-mod-download'; exit_code=1; reason=(Get-HarnessProcessFailureReason -CommandId 'go-mod-download' -Exception $_.Exception) } }
  Assert-True ($simulatedRecord.reason -eq 'COLD_PROCESS_AWAIT_EXCEPTION_GO_MOD_DOWNLOAD') 'simulated process-await cancellation was not recorded'
  $safeReason = Get-HarnessProcessFailureReason -CommandId 'go-mod-download' -Exception $processException
  Assert-True ($safeReason -eq 'COLD_PROCESS_AWAIT_EXCEPTION_GO_MOD_DOWNLOAD') 'go-mod-download await exceptions need a stable safe reason'
  Assert-True (Test-SafeEvidenceText $safeReason) 'classified process exception reason must be safe evidence text'
  Assert-True ($safeReason -notmatch '(?i)(https?://|[A-Za-z]:[\\/]|token|secret|cancelled)') 'process exception text leaked into classified reason'
  $unknownReason = Get-HarnessProcessFailureReason -CommandId 'unregistered-process' -Exception ([Exception]::new('sensitive https://example.test/value'))
  Assert-True ($unknownReason -eq 'COLD_PROCESS_AWAIT_EXCEPTION_UNKNOWN') 'unknown process-await exceptions must remain fail-closed'
  $harnessText = Get-Content -Raw (Join-Path $root 'scripts/harness.ps1')
  Assert-True ($harnessText -match 'Invoke-HarnessProcess[\s\S]*Get-HarnessProcessFailureReason[\s\S]*throw') 'cold process await exception path is not classified before outer catch'
  $run = Join-Path $root ('scripts/.runs/test-' + [guid]::NewGuid().ToString('N'))
  New-Item -ItemType Directory -Path $run -Force | Out-Null
  try {
    $commands = @(
      [ordered]@{ id='governance'; command='pwsh governance'; stage='validation'; target_label='fake'; evidence_class='contract'; duration_ms=1; exit_code=0; reason='passed'; artifact_paths=@('contracts/governance') },
      [ordered]@{ id='cold-provision'; command='go mod download'; stage='provisioning'; target_label='external-dependency-registry'; evidence_class='provisioning'; duration_ms=1; exit_code=0; reason='passed'; artifact_paths=@() }
    )
    $manifest = New-HarnessOutcome -RunId 'a1' -CandidateSha ('a' * 40) -Branch 'main' -Dirty $false -Tools ([ordered]@{ go='go version'; npm='npm --version' }) -Commands $commands -PostgresImageIdentity ('sha256:' + ('a' * 64)) -AggregateClassification 'passed'
    $path = Join-Path $run 'outcome.json'
    Write-HarnessOutcome -Outcome $manifest -Path $path
    Assert-True ((Get-Content -Raw $path | ConvertFrom-Json -Depth 30).schema_version -eq '1.0') 'outcome failed schema validation'
    $json = Get-Content -Raw $path | ConvertFrom-Json -Depth 30
    Assert-True ($json.commands.Count -eq 2) 'command inventory not preserved'
    Assert-True ($json.acceptance_link -eq 'unaccepted') 'acceptance link must use the stable unaccepted representation'
    Assert-True (@($json.tools.PSObject.Properties.Name | Sort-Object) -join '|' -eq 'go|npm') 'tool identity map was not preserved'
    Assert-True (Test-Json -LiteralPath $path -SchemaFile $schema -ErrorAction Stop) 'outcome does not validate against committed schema'
    Assert-True ($json.commands[0].id -eq 'governance') 'stable command ID was not preserved'
    $projection = Get-HarnessOutcomeProjection -Outcome $manifest
    Assert-True ($projection.commands[0].id -eq 'governance') 'stable projection did not preserve command order'
    $trace = Join-Path $run 'trace.jsonl'
    Write-HarnessTrace -Events $commands -Path $trace
    $traceLines = @(Get-Content -LiteralPath $trace)
    Assert-True ($traceLines.Count -eq 2) 'trace did not preserve all records'
    foreach ($line in $traceLines) { Assert-True (Test-SafeEvidenceText $line) 'trace emitted unsafe text' }
    $unsafe = $commands + @([ordered]@{ id='bad'; command='echo SECRET'; stage='validation'; target_label='fake'; evidence_class='contract'; duration_ms=1; exit_code=0; reason=''; artifact_paths=@('C:\secret') })
    $threw = $false; try { New-HarnessOutcome -RunId 'a2' -CandidateSha ('a' * 40) -Branch 'main' -Dirty $false -Tools @{} -Commands $unsafe -PostgresImageIdentity '' -AggregateClassification 'failed' | Out-Null } catch { $threw = $true }
    Assert-True $threw 'unsafe evidence was accepted'
    foreach ($unsafeValue in @('https://example.test/token', '\\server\share\secret', '/absolute/path', 'buyer@example.test', 'token=secret-value', 'sha256:postgres-image')) {
      $rejected = $false
      try { New-HarnessOutcome -RunId 'a3' -CandidateSha ('a' * 40) -Branch $unsafeValue -Dirty $false -Tools @{} -Commands $commands -PostgresImageIdentity '' -AggregateClassification 'blocked' | Out-Null } catch { $rejected = $true }
      Assert-True $rejected "unsafe evidence was accepted: $unsafeValue"
    }
    $invalid = Get-Content -Raw $path | ConvertFrom-Json -AsHashtable
    $invalid.acceptance_link = 'https://example.test/acceptance'
    $invalidPath = Join-Path $run 'invalid.json'
    $invalid | ConvertTo-Json -Depth 30 | Set-Content -LiteralPath $invalidPath -Encoding utf8
    Assert-True (-not (Test-Json -LiteralPath $invalidPath -SchemaFile $schema -ErrorAction SilentlyContinue)) 'schema accepted URL acceptance link'
    Write-Output 'PASS cold gate evidence contract'
  } finally { Remove-Item -LiteralPath $run -Recurse -Force -ErrorAction SilentlyContinue }
  exit 0
} catch { Write-Output "FAIL cold gate evidence contract: $($_.Exception.Message)"; exit 1 }
