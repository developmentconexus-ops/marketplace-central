$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$serverRoot = Join-Path $repoRoot 'apps/server_core'
$go = (Get-Command go -CommandType Application -ErrorAction Stop | Select-Object -First 1).Source
$environment = @{
  GOCACHE = Join-Path $serverRoot '.gocache'
  GOMODCACHE = Join-Path $serverRoot '.gomodcache'
  GOPROXY = 'off'
  GOSUMDB = 'off'
}

$originalLocation = (Get-Location).Path
$originalEnvironment = @{}
foreach ($entry in [Environment]::GetEnvironmentVariables('Process').GetEnumerator()) {
  $originalEnvironment[[string]$entry.Key] = [string]$entry.Value
}
$exitCode = 1
try {
  Set-Location -LiteralPath $serverRoot
  foreach ($entry in $environment.GetEnumerator()) { [Environment]::SetEnvironmentVariable($entry.Key, $entry.Value, 'Process') }
  $output = & $go test ./internal/platform/migrate ./internal/testsupport/postgres -count=1 2>&1
  $output | ForEach-Object { Write-Output $_ }
  $exitCode = $LASTEXITCODE
  # `go test` exits 0 for a package with no test files too ("? ... [no test
  # files]"), so exit code alone cannot distinguish ran-and-passed from
  # compiled-and-ran-nothing. Both packages must report `ok`.
  $okCount = @($output | Where-Object { "$_" -match '^ok\s' }).Count
  if ($exitCode -eq 0 -and $okCount -ne 2) {
    Write-Output "FAIL postgres go tests: expected 2 'ok' package lines, saw $okCount"
    $exitCode = 1
  } elseif ($exitCode -eq 0) {
    Write-Output 'PASS postgres go tests packages=2'
  }
} finally {
  Set-Location -LiteralPath $originalLocation
  foreach ($key in @([Environment]::GetEnvironmentVariables('Process').Keys)) {
    if (-not $originalEnvironment.ContainsKey([string]$key)) {
      [Environment]::SetEnvironmentVariable([string]$key, $null, 'Process')
    }
  }
  foreach ($key in $originalEnvironment.Keys) {
    [Environment]::SetEnvironmentVariable($key, $originalEnvironment[$key], 'Process')
  }
}

exit $exitCode
