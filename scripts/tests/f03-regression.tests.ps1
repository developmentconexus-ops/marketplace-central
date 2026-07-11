param([Parameter(Mandatory)][ValidatePattern('^[0-9a-f]{40}$')][string]$BaseSha)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$pwsh = [IO.Path]::GetFullPath((Get-Command pwsh -CommandType Application -ErrorAction Stop | Select-Object -First 1).Source)
Import-Module (Join-Path $repoRoot 'scripts/harness/Execution.psm1') -Force
$environment = [System.Collections.Generic.Dictionary[string,string]]::new([StringComparer]::OrdinalIgnoreCase)
foreach ($key in @('SystemRoot', 'WINDIR', 'ComSpec', 'PATH', 'PATHEXT', 'TEMP', 'TMP')) {
  $value = [Environment]::GetEnvironmentVariable($key, 'Process')
  if (-not [string]::IsNullOrWhiteSpace($value)) { $environment[$key] = $value }
}

function Invoke-RegressionScript([string[]]$Arguments) {
  $request = New-HarnessProcessRequest -FilePath $pwsh -ArgumentList (@('-NoProfile', '-ExecutionPolicy', 'Bypass') + $Arguments) -WorkingDirectory $repoRoot -Environment $environment -TimeoutSeconds 1200
  $result = Invoke-HarnessProcess -Request $request
  if (-not [string]::IsNullOrWhiteSpace($result.Stdout)) { [Console]::Out.WriteLine($result.Stdout.Trim()) }
  if (-not [string]::IsNullOrWhiteSpace($result.Stderr)) { [Console]::Error.WriteLine($result.Stderr.Trim()) }
  return $result.ExitCode
}
$scripts = @(
  'scripts/tests/postgres-contract.tests.ps1',
  'scripts/tests/postgres-lifecycle.tests.ps1',
  'scripts/tests/governance-contracts.tests.ps1'
)
foreach ($script in $scripts) {
  $exitCode = Invoke-RegressionScript @('-File', (Join-Path $repoRoot $script))
  if ($exitCode -ne 0) { throw "F-03 regression failed script=$script" }
}
$governanceExit = Invoke-RegressionScript @('-File', (Join-Path $repoRoot 'scripts/harness.ps1'), '-Command', 'governance-drift', '-BaseSha', $BaseSha)
if ($governanceExit -ne 0) { throw 'F-03 governance drift regression failed' }

Write-Output 'PASS F-03 regression gate target=fake'
