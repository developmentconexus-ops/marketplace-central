Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$script:ReasonByInvariant = @{
  'unknown-to-zero' = 'HEVAL_UNKNOWN_TO_ZERO'
  'openapi-without-sdk' = 'HEVAL_API_SDK_DRIFT'
  'forbidden-import' = 'HEVAL_FORBIDDEN_IMPORT'
  'mock-as-live' = 'HEVAL_EVIDENCE_INFLATION'
  'unsafe-provider-write' = 'HEVAL_PROVIDER_WRITE_UNSAFE'
  'stale-context' = 'HEVAL_CONTEXT_STALE'
  'out-of-scope-write' = 'HEVAL_PATH_OUT_OF_SCOPE'
  'concurrent-seam' = 'HEVAL_SEAM_CONFLICT'
  'risk-misroute' = 'HEVAL_RISK_MISROUTE'
  'unrelated-read' = 'HEVAL_CONTEXT_UNRELATED_READ'
  'resume-failure' = 'HEVAL_RESUME_FAILURE'
  'recursive-fan-out' = 'HEVAL_RECURSIVE_FAN_OUT'
  'cold-surface-active' = 'HEVAL_COLD_SURFACE_ACTIVE'
}

function Test-HarnessEvalCase {
  param([Parameter(Mandatory)][object]$Case)
  $required = @('case_id', 'invariant', 'expected', 'initial_context', 'artifact_path')
  if (@($required | Where-Object { -not $Case.PSObject.Properties[$_] }).Count -gt 0) { throw 'HEVAL_CASE_INVALID' }
  if (-not $script:ReasonByInvariant.ContainsKey([string]$Case.invariant)) { throw 'HEVAL_CASE_UNKNOWN' }
  if ([string]$Case.artifact_path -notmatch '^scripts/\.runs/[A-Za-z0-9._/-]+\.json$') { throw 'HEVAL_ARTIFACT_PATH_INVALID' }
  $context = $Case.initial_context
  foreach ($field in @('total_estimated_tokens', 'required_read_count', 'route_miss_count', 'unrelated_read_count', 'risk_level')) {
    if (-not $context.PSObject.Properties[$field]) { throw 'HEVAL_CONTEXT_INVALID' }
  }
  if ([int]$context.total_estimated_tokens -gt 2000 -and [string]$context.risk_level -notin @('L2', 'L3')) { throw 'HEVAL_CONTEXT_BUDGET_INVALID' }
  if ([int]$context.total_estimated_tokens -gt 2000 -and @($context.overflow_reasons).Count -eq 0) { throw 'HEVAL_CONTEXT_OVERFLOW_UNJUSTIFIED' }
  if ([int]$context.unrelated_read_count -gt 0 -and [string]$Case.invariant -ne 'unrelated-read') { throw 'HEVAL_CONTEXT_UNRELATED_UNDECLARED' }
  if ([int]$context.route_miss_count -gt 0 -and -not [bool]$Case.observed.route_miss_resolved -and [string]$Case.invariant -ne 'unrelated-read') { throw 'HEVAL_ROUTE_MISS_UNRESOLVED' }
}

function Test-HarnessEvalViolation {
  param([Parameter(Mandatory)][object]$Case)
  switch ([string]$Case.invariant) {
    'unknown-to-zero' { return [bool]$Case.observed.unknown_defaulted }
    'openapi-without-sdk' { return -not [bool]$Case.observed.sdk_updated }
    'forbidden-import' { return -not [bool]$Case.observed.import_allowed }
    'mock-as-live' { return [bool]$Case.observed.evidence_promoted }
    'unsafe-provider-write' { return -not [bool]$Case.observed.safe_write_proof }
    'stale-context' { return -not [bool]$Case.observed.hash_current }
    'out-of-scope-write' { return -not [bool]$Case.observed.changed_path_allowed }
    'concurrent-seam' { return [bool]$Case.observed.lease_conflict }
    'risk-misroute' { return -not [bool]$Case.observed.route_correct }
    'unrelated-read' { return [int]$Case.initial_context.unrelated_read_count -gt 0 }
    'resume-failure' { return -not [bool]$Case.observed.next_action_reconstructed }
    'recursive-fan-out' { return -not [bool]$Case.observed.depth_one }
    'cold-surface-active' { return [bool]$Case.observed.cold_surface_found }
  }
  throw 'HEVAL_CASE_UNKNOWN'
}

function Invoke-HarnessEvalSuite {
  [CmdletBinding()]
  param([Parameter(Mandatory)][string]$CorpusPath)
  $corpus = Get-Content -Raw -LiteralPath $CorpusPath | ConvertFrom-Json -Depth 30
  if ([string]$corpus.schema_version -ne '1.0' -or [string]::IsNullOrWhiteSpace([string]$corpus.registry_id)) { throw 'HEVAL_CORPUS_INVALID' }
  $cases = @($corpus.cases)
  if ($cases.Count -eq 0 -or @($cases.case_id | Group-Object | Where-Object Count -gt 1).Count -gt 0) { throw 'HEVAL_CORPUS_INVALID' }
  $results = [Collections.Generic.List[object]]::new()
  foreach ($case in $cases) {
    $timer = [Diagnostics.Stopwatch]::StartNew()
    Test-HarnessEvalCase -Case $case
    $timer.Stop()
    $violation = Test-HarnessEvalViolation -Case $case
    $reason = if ($violation) { $script:ReasonByInvariant[[string]$case.invariant] } else { 'HEVAL_NO_VIOLATION' }
    $verdict = if ($violation) { 'rejected' } else { 'accepted' }
    $exitCode = if ($violation) { 1 } else { 0 }
    $expected = $case.expected
    $matches = [string]$expected.verdict -eq $verdict -and [string]$expected.reason -eq $reason -and [int]$expected.exit_code -eq $exitCode -and [string]$expected.target -eq 'fake'
    $results.Add([ordered]@{
      case_id = [string]$case.case_id
      verdict = $verdict
      reason = $reason
      duration_ms = [int]$timer.ElapsedMilliseconds
      initial_context = [ordered]@{
        total_estimated_tokens = [int]$case.initial_context.total_estimated_tokens
        required_read_count = [int]$case.initial_context.required_read_count
        route_miss_count = [int]$case.initial_context.route_miss_count
        unrelated_read_count = [int]$case.initial_context.unrelated_read_count
      }
      target = 'fake'
      evidence_class = 'contract'
      exit_code = $exitCode
      artifact_path = [string]$case.artifact_path
      expected_match = $matches
    })
  }
  [pscustomobject]@{ Passed=(@($results | Where-Object { -not $_.expected_match }).Count -eq 0); Results=@($results) }
}

function Write-HarnessEvalResults {
  [CmdletBinding()]
  param([Parameter(Mandatory)][object]$Suite, [Parameter(Mandatory)][string]$OutputPath)
  $fullPath = [IO.Path]::GetFullPath($OutputPath)
  if ([IO.Path]::GetExtension($fullPath) -ne '.json') { throw 'HEVAL_OUTPUT_PATH_INVALID' }
  New-Item -ItemType Directory -Path (Split-Path -Parent $fullPath) -Force | Out-Null
  [ordered]@{ schema_version='1.0'; target='fake'; evidence_class='contract'; passed=[bool]$Suite.Passed; results=@($Suite.Results) } | ConvertTo-Json -Depth 30 | Set-Content -LiteralPath $fullPath -Encoding utf8
  $fullPath
}

Export-ModuleMember -Function Invoke-HarnessEvalSuite, Write-HarnessEvalResults
