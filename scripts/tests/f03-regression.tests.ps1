param([Parameter(Mandatory)][ValidatePattern('^[0-9a-f]{40}$')][string]$BaseSha)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$pwsh = [IO.Path]::GetFullPath((Get-Command pwsh -CommandType Application -ErrorAction Stop | Select-Object -First 1).Source)

function Invoke-RegressionScript([string]$Path, [string[]]$Arguments = @()) {
  & $pwsh -NoProfile -ExecutionPolicy Bypass -File $Path @Arguments
  if ($LASTEXITCODE -ne 0) { throw "F-03 regression child failed path=$([IO.Path]::GetFileName($Path)) exit_code=$LASTEXITCODE" }
}
$scripts = @(
  'scripts/tests/postgres-contract.tests.ps1',
  'scripts/tests/postgres-lifecycle.tests.ps1',
  'scripts/tests/postgres-go.tests.ps1',
  'scripts/tests/harness-environment.tests.ps1',
  'scripts/tests/harness-execution.tests.ps1',
  'scripts/tests/hermetic-lanes.tests.ps1',
  'scripts/tests/harness-aliases.tests.ps1',
  'scripts/tests/governance-contracts.tests.ps1',
  'scripts/tests/context-compiler.tests.ps1'
)
foreach ($script in $scripts) {
  Invoke-RegressionScript (Join-Path $repoRoot $script)
}

$harness = Join-Path $repoRoot 'scripts/harness.ps1'
Invoke-RegressionScript $harness @('-Command', 'governance', '-BaseSha', $BaseSha)

$featurePath = '.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/F-03-ephemeral-postgres-migrations'
$plan = Get-Content -Raw -LiteralPath (Join-Path $repoRoot "$featurePath/plan.md")
$contractJson = [regex]::Match($plan, '(?ms)^```json\s*\r?\n(?<json>.*?)\r?\n```').Groups['json'].Value
$allowedPaths = @((ConvertFrom-Json $contractJson -Depth 100).allowed_paths)
$compiledOutput = @(& $harness -Command context-compile -FeaturePath $featurePath -BaseSha $BaseSha -AllowedPath $allowedPaths)
$contextPath = [string](($compiledOutput | Where-Object { $_ -like 'artifact_path=*' } | Select-Object -Last 1) -replace '^artifact_path=', '')
if ([string]::IsNullOrWhiteSpace($contextPath)) { throw 'F-03 current context compile did not return an artifact' }
& $harness -Command context-validate -ContextPath $contextPath -RequireCurrentBase

Write-Output 'PASS F-03 regression gate target=fake'
