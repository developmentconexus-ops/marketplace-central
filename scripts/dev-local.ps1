param([ValidateSet('up','status','build','test','stop', ErrorMessage='unsupported command')][string]$Command = 'status')

$ErrorActionPreference = 'Stop'
$repositoryRoot = Split-Path -Parent $PSScriptRoot
. (Join-Path $repositoryRoot 'scripts/lib/dev-local-runtime.ps1')

switch ($Command) {
  'up' {
    Import-DevLocalEnvironment $repositoryRoot
    Assert-DevLocalPrerequisites $repositoryRoot
    Start-DevLocalPostgres $repositoryRoot
    Wait-DevLocalPostgres $repositoryRoot
    Invoke-DevLocalMigrations $repositoryRoot
    Start-DevLocalBackend $repositoryRoot
    Start-DevLocalFrontend $repositoryRoot
    Get-DevLocalRuntimeStatus $repositoryRoot
  }
  'status' { Get-DevLocalRuntimeStatus $repositoryRoot }
  'build' { Import-DevLocalEnvironment $repositoryRoot; Invoke-DevLocalBuild $repositoryRoot }
  'test' { Import-DevLocalEnvironment $repositoryRoot; Invoke-DevLocalTests $repositoryRoot }
  'stop' { Stop-DevLocalChildren $repositoryRoot }
}
