$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$evalModule = Join-Path $repoRoot 'scripts/harness/Evals.psm1'
$corpusPath = Join-Path $repoRoot 'contracts/governance/harness-evals.json'
$outOfScopeFixture = Join-Path $repoRoot 'scripts/tests/fixtures/harness-eval/out-of-scope-write/attempt.json'

if (-not (Test-Path -LiteralPath $evalModule -PathType Leaf)) { throw 'RED: deterministic harness eval module is missing' }
if (-not (Test-Path -LiteralPath $corpusPath -PathType Leaf)) { throw 'RED: versioned harness eval corpus is missing' }
if (-not (Test-Path -LiteralPath $outOfScopeFixture -PathType Leaf)) { throw 'RED: out-of-scope-write dogfood fixture is missing' }

function Assert-True([bool]$Condition, [string]$Message) {
  if (-not $Condition) { throw $Message }
}

Import-Module $evalModule -Force
$corpus = Get-Content -Raw -LiteralPath $corpusPath | ConvertFrom-Json -Depth 30
$expected = @(
  'unknown-to-zero', 'openapi-without-sdk', 'forbidden-import', 'mock-as-live',
  'unsafe-provider-write', 'stale-context', 'out-of-scope-write', 'concurrent-seam',
  'risk-misroute', 'unrelated-read', 'resume-failure', 'recursive-fan-out', 'cold-surface-active'
)
Assert-True (@($corpus.cases).Count -eq $expected.Count) 'RED: corpus does not contain every required deterministic case'
foreach ($caseId in $expected) { Assert-True (@($corpus.cases | Where-Object case_id -eq $caseId).Count -eq 1) "RED: missing eval case $caseId" }
$outOfScopeCase = @($corpus.cases | Where-Object case_id -eq 'out-of-scope-write')[0]
$fixture = Get-Content -Raw -LiteralPath $outOfScopeFixture | ConvertFrom-Json -Depth 10
Assert-True ($outOfScopeCase.expected.reason -eq 'HEVAL_PATH_OUT_OF_SCOPE') 'RED: out-of-scope-write rejection reason drifted'
Assert-True ($fixture.attempted_path -eq 'apps/server_core/internal/modules/inventory/domain/stock_policy.go') 'RED: dogfood fixture must attempt the pinned out-of-scope path'

$suite = Invoke-HarnessEvalSuite -CorpusPath $corpusPath
Assert-True $suite.Passed 'deterministic eval suite did not match pinned expectations'
foreach ($result in @($suite.Results)) {
  Assert-True (($result.verdict -eq 'rejected') -and ($result.exit_code -eq 1) -and ($result.target -eq 'fake')) "eval result shape failed for $($result.case_id)"
  Assert-True ((($result.initial_context.required_read_count -gt 0) -and ($result.initial_context.unrelated_read_count -eq 0)) -or ($result.case_id -eq 'unrelated-read')) "context metrics failed for $($result.case_id)"
  Assert-True ($result.artifact_path -match '^scripts/\.runs/') "artifact path is not relative for $($result.case_id)"
}
$runPath = Join-Path $repoRoot ('scripts/.runs/harness-eval-tests/' + [guid]::NewGuid().ToString('N') + '/results.json')
$written = Write-HarnessEvalResults -Suite $suite -OutputPath $runPath
Assert-True (Test-Path -LiteralPath $written -PathType Leaf) 'result manifest was not written'
Write-Output 'PASS harness eval tests'
