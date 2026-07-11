$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$serverRoot = Join-Path $repoRoot 'apps/server_core'
$go = (Get-Command go -CommandType Application -ErrorAction Stop).Source
$environment = @{
  GOCACHE = Join-Path $repoRoot '.gocache'
  GOPROXY = 'off'
  GOSUMDB = 'off'
}

Push-Location -LiteralPath $serverRoot
try {
  foreach ($entry in $environment.GetEnumerator()) { [Environment]::SetEnvironmentVariable($entry.Key, $entry.Value, 'Process') }
  & $go test ./internal/platform/migrate ./internal/testsupport/postgres -count=1
  $exitCode = $LASTEXITCODE
} finally {
  Pop-Location
}

exit $exitCode
