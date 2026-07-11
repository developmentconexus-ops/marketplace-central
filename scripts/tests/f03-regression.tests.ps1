param(
  [Parameter(Mandatory)][ValidatePattern('^[0-9a-f]{40}$')][string]$BaseSha,
  [switch]$EnvironmentProbe,
  [switch]$ContextWorker
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$scriptPath = [IO.Path]::GetFullPath($PSCommandPath)
$pwsh = [IO.Path]::GetFullPath((Get-Command pwsh -CommandType Application -ErrorAction Stop | Select-Object -First 1).Source)
$harness = Join-Path $repoRoot 'scripts/harness.ps1'
Import-Module (Join-Path $repoRoot 'scripts/harness/Environment.psm1') -Force
Import-Module (Join-Path $repoRoot 'scripts/harness/Execution.psm1') -Force

$contaminatedKeys = @(
  'MC_DATABASE_URL',
  'MPC_TEST_DATABASE_URL',
  'MPC_ORACLE_PASSWORD',
  'MPC_PROVIDER_MERCADOLIVRE_CLIENT_SECRET'
)

if ($EnvironmentProbe) {
  $count = @($contaminatedKeys | Where-Object { -not [string]::IsNullOrWhiteSpace([Environment]::GetEnvironmentVariable($_, 'Process')) }).Count
  Write-Output "contamination_count=$count"
  return
}

if ($ContextWorker) {
  $featurePath = '.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/F-03-ephemeral-postgres-migrations'
  $plan = Get-Content -Raw -LiteralPath (Join-Path $repoRoot "$featurePath/plan.md")
  $contractJson = [regex]::Match($plan, '(?ms)^```json\s*\r?\n(?<json>.*?)\r?\n```').Groups['json'].Value
  $allowedPaths = @((ConvertFrom-Json $contractJson -Depth 100).allowed_paths)
  $compiledOutput = @(& $harness -Command context-compile -FeaturePath $featurePath -BaseSha $BaseSha -AllowedPath $allowedPaths)
  $compiledOutput | ForEach-Object { Write-Output $_ }
  $contextPath = [string](($compiledOutput | Where-Object { $_ -like 'artifact_path=*' } | Select-Object -Last 1) -replace '^artifact_path=', '')
  if ([string]::IsNullOrWhiteSpace($contextPath)) { throw 'F-03 current context compile did not return an artifact' }
  & $harness -Command context-validate -ContextPath $contextPath -RequireCurrentBase
  Write-Output 'PASS F-03 context current target=fake'
  return
}

function Invoke-RegressionChild([string]$Path, [string[]]$Arguments = @()) {
  $request = New-HarnessProcessRequest `
    -FilePath $pwsh `
    -ArgumentList (@('-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', [IO.Path]::GetFullPath($Path)) + $Arguments) `
    -WorkingDirectory $repoRoot `
    -Environment $script:childEnvironment `
    -TimeoutSeconds 1200 `
    -RedactionCandidates @($repoRoot)
  $result = Invoke-HarnessProcess -Request $request
  if (-not [string]::IsNullOrWhiteSpace($result.Stdout)) { [Console]::Out.WriteLine($result.Stdout.Trim()) }
  if (-not [string]::IsNullOrWhiteSpace($result.Stderr)) { [Console]::Error.WriteLine($result.Stderr.Trim()) }
  return $result
}

$originalValues = @{}
try {
  foreach ($key in $contaminatedKeys) {
    $originalValues[$key] = [Environment]::GetEnvironmentVariable($key, 'Process')
    [Environment]::SetEnvironmentVariable($key, "fixture-parent-$key", 'Process')
  }
  $script:childEnvironment = New-HarnessChildEnvironment -RepositoryRoot $repoRoot -LaneId 'unit'

  $probe = Invoke-RegressionChild $scriptPath @('-BaseSha', $BaseSha, '-EnvironmentProbe')
  if ($probe.ExitCode -ne 0 -or $probe.Stdout.Trim() -cne 'contamination_count=0') { throw 'F-03 regression child inherited parent runtime targets' }

  foreach ($script in @(
    'scripts/tests/postgres-contract.tests.ps1',
    'scripts/tests/postgres-lifecycle.tests.ps1',
    'scripts/tests/postgres-go.tests.ps1',
    'scripts/tests/harness-environment.tests.ps1',
    'scripts/tests/harness-execution.tests.ps1',
    'scripts/tests/hermetic-lanes.tests.ps1',
    'scripts/tests/harness-aliases.tests.ps1',
    'scripts/tests/governance-contracts.tests.ps1',
    'scripts/tests/context-compiler.tests.ps1'
  )) {
    $result = Invoke-RegressionChild (Join-Path $repoRoot $script)
    if ($result.ExitCode -ne 0) { throw "F-03 regression child failed path=$([IO.Path]::GetFileName($script)) exit_code=$($result.ExitCode)" }
  }

  $governance = Invoke-RegressionChild $harness @('-Command', 'governance', '-BaseSha', $BaseSha)
  if ($governance.ExitCode -ne 0) { throw "F-03 governance regression failed exit_code=$($governance.ExitCode)" }
  $context = Invoke-RegressionChild $scriptPath @('-BaseSha', $BaseSha, '-ContextWorker')
  if ($context.ExitCode -ne 0 -or $context.Stdout -notmatch '(?m)^PASS F-03 context current target=fake\r?$') { throw 'F-03 current context regression failed' }
} finally {
  foreach ($key in $contaminatedKeys) { [Environment]::SetEnvironmentVariable($key, $originalValues[$key], 'Process') }
}

Write-Output 'PASS F-03 regression gate target=fake contamination=blocked'
