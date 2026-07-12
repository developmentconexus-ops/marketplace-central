$ErrorActionPreference = 'Stop'

if (-not (Get-Module -ListAvailable -Name Pester)) { throw 'Pester is required for live Oracle Docker runner tests' }
Import-Module Pester -ErrorAction Stop

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
. (Join-Path $repoRoot 'scripts/run-live-oracle-docker.ps1')

function Assert-RunnerCondition {
  param([bool]$Condition, [string]$Message)
  if (-not $Condition) { throw $Message }
}

Describe 'Docker live Oracle runner' {
  It 'uses the existing backend image and exactly the read-only smoke test' {
    $plan = New-LiveOracleDockerPlan -Username 'fixture-user' -Password 'fixture-password' -ConnectString 'fixture-connect'

    Assert-RunnerCondition ($plan.BuildArguments -contains (Join-Path $repoRoot 'docker/dev/backend.Dockerfile')) 'backend Dockerfile missing'
    Assert-RunnerCondition ($plan.RunArguments -contains 'go') 'Go command missing'
    Assert-RunnerCondition ($plan.RunArguments -contains 'test') 'Go test command missing'
    Assert-RunnerCondition ($plan.RunArguments -contains './internal/modules/internal_read/adapters/oracle') 'Oracle package missing'
    Assert-RunnerCondition ($plan.RunArguments -contains '^TestOracleLiveSmoke$') 'smoke test regex missing'
    Assert-RunnerCondition (($plan.RunArguments -join ' ') -notmatch 'docker\s+compose|\.env|migrations?|(^|\s)(air|serve|start)(\s|$)|entrypoint') 'forbidden command fragment present'
    Assert-RunnerCondition (@($plan.RunArguments | Where-Object { $_ -eq '--mount' }).Count -eq 1) 'read-only workspace mount missing'
    Assert-RunnerCondition (($plan.RunArguments -join ' ') -match 'readonly') 'workspace mount is not read-only'
  }

  It 'maps the explicit Sankhya caller keys to the governed container keys by Docker key reference' {
    $plan = New-LiveOracleDockerPlan -Username 'fixture-user' -Password 'fixture-password' -ConnectString 'fixture-connect'
    $expected = @('MPC_ORACLE_USERNAME', 'MPC_ORACLE_PASSWORD', 'MPC_ORACLE_CONNECT_STRING', 'MPC_ORACLE_LIVE_TEST', 'MPC_ORACLE_LIB_DIR')

    Assert-RunnerCondition ((@($plan.ContainerEnvironment.Keys) -join ',') -eq ($expected -join ',')) 'container environment key set changed'
    $envReferences = @($plan.RunArguments | Where-Object { $_ -like 'MPC_*' })
    Assert-RunnerCondition (($envReferences -join ',') -eq ($expected -join ',')) 'Docker environment forwarding changed'
    Assert-RunnerCondition ($plan.ContainerEnvironment['MPC_ORACLE_USERNAME'] -eq 'fixture-user') 'Sankhya username was not mapped to the governed container key'
    Assert-RunnerCondition ($plan.ContainerEnvironment['MPC_ORACLE_PASSWORD'] -eq 'fixture-password') 'Sankhya password was not mapped to the governed container key'
    Assert-RunnerCondition ($plan.ContainerEnvironment['MPC_ORACLE_CONNECT_STRING'] -eq 'fixture-connect') 'Sankhya connect string was not mapped to the governed container key'
    Assert-RunnerCondition (($plan.RunArguments -join ' ') -notmatch 'fixture-user|fixture-password|fixture-connect') 'credential value entered Docker arguments'
  }

  It 'rejects incomplete explicit Sankhya credentials before a Docker invocation plan exists' {
    $errorMessage = ''
    try { New-LiveOracleDockerPlan -Username 'fixture-user' -Password '' -ConnectString 'fixture-connect' | Out-Null } catch { $errorMessage = $_.Exception.Message }
    Assert-RunnerCondition ($errorMessage -like '*missing_keys=MPC_SANKHYA_ORACLE_PASSWORD*') 'missing explicit Sankhya credential was not rejected'
  }

  It 'rejects generic and ambient credential aliases when explicit Sankhya keys are absent' {
    $rejectedKeys = @(
      'MPC_SANKHYA_ORACLE_USERNAME', 'MPC_SANKHYA_ORACLE_PASSWORD', 'MPC_SANKHYA_ORACLE_CONNECT_STRING',
      'MPC_ORACLE_USERNAME', 'MPC_ORACLE_PASSWORD', 'MPC_ORACLE_CONNECT_STRING',
      'SANKHYA_ORACLE_USER', 'ORACLE_USERNAME'
    )
    $original = @{}
    try {
      foreach ($key in $rejectedKeys) {
        $original[$key] = [Environment]::GetEnvironmentVariable($key, 'Process')
        [Environment]::SetEnvironmentVariable($key, $null, 'Process')
      }
      [Environment]::SetEnvironmentVariable('MPC_ORACLE_USERNAME', 'generic-user', 'Process')
      [Environment]::SetEnvironmentVariable('MPC_ORACLE_PASSWORD', 'generic-password', 'Process')
      [Environment]::SetEnvironmentVariable('MPC_ORACLE_CONNECT_STRING', 'generic-connect', 'Process')
      [Environment]::SetEnvironmentVariable('SANKHYA_ORACLE_USER', 'alias-user', 'Process')
      [Environment]::SetEnvironmentVariable('ORACLE_USERNAME', 'ambient-user', 'Process')
      $errorMessage = ''
      try { Get-LiveOracleCredentialValues | Out-Null } catch { $errorMessage = $_.Exception.Message }
      Assert-RunnerCondition ($errorMessage -like '*MPC_SANKHYA_ORACLE_USERNAME*') 'generic credentials were accepted as an explicit Sankhya username'
      Assert-RunnerCondition ($errorMessage -like '*MPC_SANKHYA_ORACLE_PASSWORD*') 'generic credentials were accepted as an explicit Sankhya password'
      Assert-RunnerCondition ($errorMessage -like '*MPC_SANKHYA_ORACLE_CONNECT_STRING*') 'generic credentials were accepted as an explicit Sankhya connect string'
    } finally {
      foreach ($key in $rejectedKeys) {
        [Environment]::SetEnvironmentVariable($key, $original[$key], 'Process')
      }
    }
  }

  It 'strips ambient configuration from the Docker child and keeps exactly five governed runtime keys' {
    $ambient = [ordered]@{
      MPC_UNRELATED_SETTING = 'ambient-mpc'
      MPC_DATABASE_URL = 'ambient-database'
      ORACLE_HOME = 'ambient-oracle'
      MELI_ACCESS_TOKEN = 'ambient-provider'
      SANKHYA_ORACLE_USER = 'ambient-alias'
    }
    $original = @{}
    try {
      foreach ($key in $ambient.Keys) {
        $original[$key] = [Environment]::GetEnvironmentVariable($key, 'Process')
        [Environment]::SetEnvironmentVariable($key, [string]$ambient[$key], 'Process')
      }
      $plan = New-LiveOracleDockerPlan -Username 'fixture-user' -Password 'fixture-password' -ConnectString 'fixture-connect'
      $startInfo = New-LiveOracleDockerProcessStartInfo -DockerPath 'C:\fixture\docker.exe' -Arguments $plan.RunArguments -Environment $plan.ContainerEnvironment
      $expectedRuntime = @('MPC_ORACLE_USERNAME', 'MPC_ORACLE_PASSWORD', 'MPC_ORACLE_CONNECT_STRING', 'MPC_ORACLE_LIVE_TEST', 'MPC_ORACLE_LIB_DIR')
      $expectedKeys = @($script:DockerExecutionEnvironmentKeys | Where-Object { -not [string]::IsNullOrWhiteSpace([Environment]::GetEnvironmentVariable($_, 'Process')) }) + $expectedRuntime
      $actualKeys = @($startInfo.Environment.Keys)
      $actualRuntime = @($actualKeys | Where-Object { $_ -like 'MPC_*' })

      Assert-RunnerCondition ((@($actualKeys | Sort-Object) -join ',') -eq (@($expectedKeys | Sort-Object) -join ',')) 'Docker child environment was not an execution-key plus runtime-key allowlist'
      Assert-RunnerCondition (($actualRuntime -join ',') -eq ($expectedRuntime -join ',')) 'Docker child runtime environment was not exactly the five canonical keys'
      foreach ($key in $ambient.Keys) {
        Assert-RunnerCondition (-not ($actualKeys -contains $key)) "ambient key entered Docker child: $key"
      }
      Assert-RunnerCondition ($startInfo.Environment['MPC_ORACLE_USERNAME'] -eq 'fixture-user') 'Docker child did not receive the governed username key'
    } finally {
      foreach ($key in $ambient.Keys) {
        [Environment]::SetEnvironmentVariable($key, $original[$key], 'Process')
      }
    }
  }

  It 'preflights Docker through an execution-only child with no ambient or runtime keys' {
    $ambient = [ordered]@{
      MPC_UNRELATED_SETTING = 'ambient-mpc'
      MPC_SANKHYA_ORACLE_USERNAME = 'ambient-runtime'
      MPC_DATABASE_URL = 'ambient-database'
      ORACLE_HOME = 'ambient-oracle'
      MELI_ACCESS_TOKEN = 'ambient-provider'
      SANKHYA_ORACLE_USER = 'ambient-alias'
    }
    $original = @{}
    try {
      foreach ($key in $ambient.Keys) {
        $original[$key] = [Environment]::GetEnvironmentVariable($key, 'Process')
        [Environment]::SetEnvironmentVariable($key, [string]$ambient[$key], 'Process')
      }
      $preflight = New-LiveOracleDockerPreflightPlan -DockerPath 'C:\fixture\docker.exe'
      $expectedKeys = @($script:DockerExecutionEnvironmentKeys | Where-Object { -not [string]::IsNullOrWhiteSpace([Environment]::GetEnvironmentVariable($_, 'Process')) })
      $actualKeys = @($preflight.StartInfo.Environment.Keys)

      Assert-RunnerCondition (($preflight.Arguments -join ',') -eq 'version,--format,{{.Server.Version}}') 'Docker preflight arguments changed'
      Assert-RunnerCondition ((@($actualKeys | Sort-Object) -join ',') -eq (@($expectedKeys | Sort-Object) -join ',')) 'Docker preflight child was not an execution-key-only allowlist'
      Assert-RunnerCondition (@($actualKeys | Where-Object { $_ -like 'MPC_*' }).Count -eq 0) 'Docker preflight child included a runtime key'
      foreach ($key in $ambient.Keys) {
        Assert-RunnerCondition (-not ($actualKeys -contains $key)) "ambient key entered Docker preflight child: $key"
      }
    } finally {
      foreach ($key in $ambient.Keys) {
        [Environment]::SetEnvironmentVariable($key, $original[$key], 'Process')
      }
    }
  }
}
