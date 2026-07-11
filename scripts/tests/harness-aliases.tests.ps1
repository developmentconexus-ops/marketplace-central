$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$harness = Join-Path $repoRoot 'scripts/harness.ps1'
$foreign = Join-Path ([IO.Path]::GetTempPath()) ('harness-cwd-' + [guid]::NewGuid().ToString('N'))

function Invoke-FromForeign([string]$Alias, [string[]]$Arguments) {
  Push-Location $foreign
  try {
    $output = & npm --prefix $repoRoot run $Alias -- @Arguments 2>&1
    [pscustomobject]@{ ExitCode = $LASTEXITCODE; Output = ($output -join "`n") }
  } finally { Pop-Location }
}

function Assert-Result($Result, [int]$ExitCode, [string]$Needle, [string]$Name) {
  if ($Result.ExitCode -ne $ExitCode) { throw "$Name exit mismatch" }
  if ($Result.Output -notmatch [regex]::Escape($Needle)) { throw "$Name behavior mismatch" }
}

New-Item -ItemType Directory -Path $foreign -Force | Out-Null
$unknownKey = 'HARNESS_UNKNOWN_' + [guid]::NewGuid().ToString('N')
$contaminatedKeys = @($unknownKey, 'MC_DATABASE_URL')
$before = @{}
try {
  foreach ($key in $contaminatedKeys) { $before[$key] = [Environment]::GetEnvironmentVariable($key, 'Process') }
  [Environment]::SetEnvironmentVariable($unknownKey, 'fixture-contamination', 'Process')
  [Environment]::SetEnvironmentVariable('MC_DATABASE_URL', 'fixture-contamination', 'Process')
  Assert-Result (Invoke-FromForeign 'harness:unit' @('-PreflightOnly')) 0 'target=fake' 'unit alias'
  Assert-Result (Invoke-FromForeign 'harness:integration' @('-PreflightOnly')) 1 'F-03' 'integration alias'
  Assert-Result (Invoke-FromForeign 'harness:live' @('-PreflightOnly', '-EnvFile', (Join-Path $foreign 'missing.env'))) 1 'target=live-oracle' 'live alias'
  Assert-Result (Invoke-FromForeign 'harness:browser' @('-PreflightOnly')) 0 'target=browser' 'browser alias'
  Assert-Result (Invoke-FromForeign 'harness:provider-write' @()) 1 'provider is required' 'provider-write alias'

  $sha = (git -C $repoRoot rev-parse HEAD).Trim()
  Assert-Result (Invoke-FromForeign 'harness:governance' @('-BaseSha', $sha)) 0 'status=passed' 'governance alias'
  Assert-Result (Invoke-FromForeign 'harness:context:compile' @()) 1 'CTX_FEATURE_INVALID' 'context compile alias'
  Assert-Result (Invoke-FromForeign 'harness:context:validate' @()) 1 'CTX_SOURCE_MISSING' 'context validate alias'

  Push-Location $foreign
  try {
    $direct = & pwsh -NoProfile -ExecutionPolicy Bypass -File $harness -Command unit -PreflightOnly 2>&1
    if ($LASTEXITCODE -ne 0 -or ($direct -join "`n") -notmatch 'target=fake') { throw 'absolute dispatcher invocation failed from foreign CWD' }
  } finally { Pop-Location }
} finally {
  foreach ($key in $before.Keys) { [Environment]::SetEnvironmentVariable($key, $before[$key], 'Process') }
  Remove-Item -LiteralPath $foreign -Recurse -Force -ErrorAction SilentlyContinue
}

Write-Output 'PASS harness alias tests'
