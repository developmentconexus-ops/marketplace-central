$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$serverRoot = Join-Path $repoRoot 'apps/server_core'
$go = (Get-Command go -CommandType Application -ErrorAction Stop).Source
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
  & $go test ./internal/platform/migrate ./internal/testsupport/postgres -count=1
  $exitCode = $LASTEXITCODE
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
