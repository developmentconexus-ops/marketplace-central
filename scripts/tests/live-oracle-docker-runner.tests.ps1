$ErrorActionPreference = 'Stop'

if (-not (Get-Module -ListAvailable -Name Pester)) { throw 'Pester is required for live Oracle Docker runner tests' }
Import-Module Pester -ErrorAction Stop

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
. (Join-Path $repoRoot 'scripts/run-live-oracle-docker.ps1')

function Assert-RunnerCondition {
  param([bool]$Condition, [string]$Message)
  if (-not $Condition) { throw $Message }
}

function New-CredentialFixtureFile {
  param([string[]]$Lines)

  $path = [IO.Path]::GetTempFileName()
  Set-Content -LiteralPath $path -Value $Lines -Encoding utf8
  $path
}

function Invoke-WithoutSankhyaCallerConnectionValues {
  param([Parameter(Mandatory)][scriptblock]$Action)

  $original = [ordered]@{}
  try {
    foreach ($key in $script:CallerConnectionKeys) {
      $original[$key] = [Environment]::GetEnvironmentVariable($key, 'Process')
      [Environment]::SetEnvironmentVariable($key, $null, 'Process')
    }
    & $Action
  } finally {
    foreach ($key in $script:CallerConnectionKeys) {
      [Environment]::SetEnvironmentVariable($key, $original[$key], 'Process')
    }
  }
}

Describe 'Docker live Oracle runner' {
  It 'uses the existing backend image and exactly the read-only smoke test' {
    $fixture = New-CredentialFixtureFile -Lines @('MPC_SANKHYA_ORACLE_USERNAME=u', 'MPC_SANKHYA_ORACLE_PASSWORD=p', 'MPC_SANKHYA_ORACLE_HOST=fixture-host', 'MPC_SANKHYA_ORACLE_PORT=1521', 'MPC_SANKHYA_ORACLE_CONNECT_STRING=fixture-service')
    try {
      $plan = Invoke-WithoutSankhyaCallerConnectionValues { New-LiveOracleDockerPlan -EnvFilePath $fixture }

      Assert-RunnerCondition ($plan.BuildArguments -contains (Join-Path $repoRoot 'docker/dev/backend.Dockerfile')) 'backend Dockerfile missing'
      Assert-RunnerCondition ($plan.RunArguments -contains 'go') 'Go command missing'
      Assert-RunnerCondition ($plan.RunArguments -contains 'test') 'Go test command missing'
      Assert-RunnerCondition ($plan.RunArguments -contains './internal/modules/internal_read/adapters/oracle') 'Oracle package missing'
      Assert-RunnerCondition ($plan.RunArguments -contains '^TestOracleLiveSmoke$') 'smoke test regex missing'
      Assert-RunnerCondition (($plan.RunArguments -join ' ') -notmatch 'docker\s+compose|\.env|migrations?|(^|\s)(air|serve|start)(\s|$)|entrypoint') 'forbidden command fragment present'
      Assert-RunnerCondition (@($plan.RunArguments | Where-Object { $_ -eq '--mount' }).Count -eq 1) 'read-only workspace mount missing'
      Assert-RunnerCondition (($plan.RunArguments -join ' ') -match 'readonly') 'workspace mount is not read-only'
    } finally { Remove-Item -LiteralPath $fixture -Force }
  }

  It 'constructs the governed connect string from the five Sankhya service-name inputs' {
    $fixture = New-CredentialFixtureFile -Lines @('MPC_SANKHYA_ORACLE_USERNAME=u', 'MPC_SANKHYA_ORACLE_PASSWORD=p', 'MPC_SANKHYA_ORACLE_HOST=fixture-host', 'MPC_SANKHYA_ORACLE_PORT=1521', 'MPC_SANKHYA_ORACLE_CONNECT_STRING=fixture-service')
    try {
      $plan = Invoke-WithoutSankhyaCallerConnectionValues { New-LiveOracleDockerPlan -EnvFilePath $fixture }
      $expected = @('MPC_ORACLE_USERNAME', 'MPC_ORACLE_PASSWORD', 'MPC_ORACLE_CONNECT_STRING', 'MPC_ORACLE_LIVE_TEST', 'MPC_ORACLE_LIB_DIR')

      Assert-RunnerCondition ((@($plan.ContainerEnvironment.Keys) -join ',') -eq ($expected -join ',')) 'container environment key set changed'
      Assert-RunnerCondition ((@($plan.RunArguments | Where-Object { $_ -like 'MPC_*' }) -join ',') -eq ($expected -join ',')) 'Docker environment forwarding changed'
      Assert-RunnerCondition ($plan.ContainerEnvironment['MPC_ORACLE_USERNAME'] -eq 'u') 'local username was not mapped to the governed container key'
      Assert-RunnerCondition ($plan.ContainerEnvironment['MPC_ORACLE_PASSWORD'] -eq 'p') 'local password was not mapped to the governed container key'
      Assert-RunnerCondition ($plan.ContainerEnvironment['MPC_ORACLE_CONNECT_STRING'] -eq 'fixture-host:1521/fixture-service') 'service-name inputs were not constructed as the governed connect string'
      Assert-RunnerCondition (($plan.RunArguments -join ' ') -notmatch '(^|\s)(u|p|fixture-host|1521|fixture-service)(\s|$)') 'connection value entered Docker arguments'
    } finally { Remove-Item -LiteralPath $fixture -Force }
  }

  It 'permits but does not consume or forward the unrelated local schema' {
    $fixture = New-CredentialFixtureFile -Lines @('MPC_SANKHYA_ORACLE_USERNAME=u', 'MPC_SANKHYA_ORACLE_PASSWORD=p', 'MPC_SANKHYA_ORACLE_HOST=fixture-host', 'MPC_SANKHYA_ORACLE_PORT=1521', 'MPC_SANKHYA_ORACLE_CONNECT_STRING=fixture-service', 'MPC_SANKHYA_ORACLE_SCHEMA=fixture-schema')
    try {
      $localValues = Invoke-WithoutSankhyaCallerConnectionValues { Get-LiveOracleLocalEnvValues -EnvFilePath $fixture }
      $plan = Invoke-WithoutSankhyaCallerConnectionValues { New-LiveOracleDockerPlan -EnvFilePath $fixture }
      Assert-RunnerCondition (-not $localValues.Contains('MPC_SANKHYA_ORACLE_SCHEMA')) 'local schema was consumed'
      Assert-RunnerCondition (-not (@($plan.ContainerEnvironment.Keys) -contains 'MPC_SANKHYA_ORACLE_SCHEMA')) 'local schema was forwarded'
      Assert-RunnerCondition (($plan.ContainerEnvironment['MPC_ORACLE_CONNECT_STRING']) -eq 'fixture-host:1521/fixture-service') 'local schema influenced routing'
    } finally { Remove-Item -LiteralPath $fixture -Force }
  }

  It 'ignores unrelated non-reserved local assignments without consuming or forwarding them' {
    $fixture = New-CredentialFixtureFile -Lines @('MPC_SANKHYA_ORACLE_USERNAME=u', 'MPC_SANKHYA_ORACLE_PASSWORD=p', 'MPC_SANKHYA_ORACLE_HOST=fixture-host', 'MPC_SANKHYA_ORACLE_PORT=1521', 'MPC_SANKHYA_ORACLE_CONNECT_STRING=fixture-service', 'MC_DATABASE_URL=unrelated-fixture-value')
    try {
      $localValues = Invoke-WithoutSankhyaCallerConnectionValues { Get-LiveOracleLocalEnvValues -EnvFilePath $fixture }
      $plan = Invoke-WithoutSankhyaCallerConnectionValues { New-LiveOracleDockerPlan -EnvFilePath $fixture }
      Assert-RunnerCondition (-not $localValues.Contains('MC_DATABASE_URL')) 'unrelated local assignment was consumed'
      Assert-RunnerCondition (-not (@($plan.ContainerEnvironment.Keys) -contains 'MC_DATABASE_URL')) 'unrelated local assignment was forwarded'
      Assert-RunnerCondition (($plan.ContainerEnvironment['MPC_ORACLE_CONNECT_STRING']) -eq 'fixture-host:1521/fixture-service') 'unrelated local assignment influenced routing'
    } finally { Remove-Item -LiteralPath $fixture -Force }
  }

  It 'rejects an unknown reserved local key before requesting Docker from the runner' {
    $fixture = New-CredentialFixtureFile -Lines @('MPC_SANKHYA_ORACLE_USERNAME=u', 'MPC_SANKHYA_ORACLE_PASSWORD=p', 'MPC_SANKHYA_ORACLE_HOST=fixture-host', 'MPC_SANKHYA_ORACLE_PORT=1521', 'MPC_SANKHYA_ORACLE_CONNECT_STRING=fixture-service', 'MPC_SANKHYA_ORACLE_TNS_ADMIN=x')
    try {
      Mock -CommandName Test-LiveOracleDockerAvailable -MockWith { throw 'Docker preflight must not run for an invalid local file' }
      $errorMessage = ''
      try { Invoke-WithoutSankhyaCallerConnectionValues { Invoke-LiveOracleDockerRunner -EnvFilePath $fixture } | Out-Null } catch { $errorMessage = $_.Exception.Message }

      Assert-RunnerCondition ($errorMessage -eq 'live Oracle .env unsupported_key=MPC_SANKHYA_ORACLE_TNS_ADMIN') 'runner did not reject the unknown reserved local key first'
      Assert-MockCalled -CommandName Test-LiveOracleDockerAvailable -Times 0 -Exactly
    } finally { Remove-Item -LiteralPath $fixture -Force }
  }

  It 'rejects generic local aliases before requesting Docker from the runner' {
    $fixture = New-CredentialFixtureFile -Lines @('MPC_SANKHYA_ORACLE_USERNAME=u', 'MPC_SANKHYA_ORACLE_PASSWORD=p', 'MPC_SANKHYA_ORACLE_HOST=fixture-host', 'MPC_SANKHYA_ORACLE_PORT=1521', 'MPC_SANKHYA_ORACLE_CONNECT_STRING=fixture-service', 'MPC_ORACLE_USERNAME=x')
    try {
      Mock -CommandName Test-LiveOracleDockerAvailable -MockWith { throw 'Docker preflight must not run for an invalid local file' }
      $errorMessage = ''
      try { Invoke-WithoutSankhyaCallerConnectionValues { Invoke-LiveOracleDockerRunner -EnvFilePath $fixture } | Out-Null } catch { $errorMessage = $_.Exception.Message }

      Assert-RunnerCondition ($errorMessage -eq 'live Oracle .env unsupported_key=MPC_ORACLE_USERNAME') 'runner did not reject the generic local alias first'
      Assert-MockCalled -CommandName Test-LiveOracleDockerAvailable -Times 0 -Exactly
    } finally { Remove-Item -LiteralPath $fixture -Force }
  }

  It 'gives only the same exact caller-process keys precedence over the local file' {
    $fixture = New-CredentialFixtureFile -Lines @('MPC_SANKHYA_ORACLE_USERNAME=u', 'MPC_SANKHYA_ORACLE_PASSWORD=p', 'MPC_SANKHYA_ORACLE_HOST=fixture-host', 'MPC_SANKHYA_ORACLE_PORT=1521', 'MPC_SANKHYA_ORACLE_CONNECT_STRING=fixture-service')
    $original = [ordered]@{}
    try {
      $callerValues = [ordered]@{
        MPC_SANKHYA_ORACLE_USERNAME = 'caller-user'
        MPC_SANKHYA_ORACLE_PASSWORD = 'caller-password'
        MPC_SANKHYA_ORACLE_HOST = 'caller-host'
        MPC_SANKHYA_ORACLE_PORT = '1522'
        MPC_SANKHYA_ORACLE_CONNECT_STRING = 'caller-service'
      }
      foreach ($key in $script:CallerConnectionKeys) {
        $original[$key] = [Environment]::GetEnvironmentVariable($key, 'Process')
        [Environment]::SetEnvironmentVariable($key, $callerValues[$key], 'Process')
      }
      $plan = New-LiveOracleDockerPlan -EnvFilePath $fixture
      Assert-RunnerCondition ($plan.ContainerEnvironment['MPC_ORACLE_USERNAME'] -eq 'caller-user') 'caller username did not override local value'
      Assert-RunnerCondition ($plan.ContainerEnvironment['MPC_ORACLE_PASSWORD'] -eq 'caller-password') 'caller password did not override local value'
      Assert-RunnerCondition ($plan.ContainerEnvironment['MPC_ORACLE_CONNECT_STRING'] -eq 'caller-host:1522/caller-service') 'caller service inputs did not override local values'
    } finally {
      foreach ($key in $script:CallerConnectionKeys) { [Environment]::SetEnvironmentVariable($key, $original[$key], 'Process') }
      Remove-Item -LiteralPath $fixture -Force
    }
  }

  It 'rejects generic and ambient credential aliases when the exact Sankhya keys are absent' {
    $fixture = New-CredentialFixtureFile -Lines @()
    $rejectedKeys = @('MPC_ORACLE_USERNAME', 'MPC_ORACLE_PASSWORD', 'MPC_ORACLE_CONNECT_STRING', 'SANKHYA_ORACLE_USER', 'ORACLE_USERNAME')
    $original = [ordered]@{}
    try {
      foreach ($key in $rejectedKeys) {
        $original[$key] = [Environment]::GetEnvironmentVariable($key, 'Process')
        [Environment]::SetEnvironmentVariable($key, 'x', 'Process')
      }
      $errorMessage = ''
      try { Invoke-WithoutSankhyaCallerConnectionValues { Get-LiveOracleCredentialValues -EnvFilePath $fixture } | Out-Null } catch { $errorMessage = $_.Exception.Message }
      foreach ($key in $script:CallerConnectionKeys) {
        Assert-RunnerCondition ($errorMessage -like "*${key}*") "ambient alias satisfied $key"
      }
    } finally {
      foreach ($key in $rejectedKeys) { [Environment]::SetEnvironmentVariable($key, $original[$key], 'Process') }
      Remove-Item -LiteralPath $fixture -Force
    }
  }

  It 'strips ambient configuration from the Docker child and keeps exactly five governed runtime keys' {
    $fixture = New-CredentialFixtureFile -Lines @('MPC_SANKHYA_ORACLE_USERNAME=u', 'MPC_SANKHYA_ORACLE_PASSWORD=p', 'MPC_SANKHYA_ORACLE_HOST=fixture-host', 'MPC_SANKHYA_ORACLE_PORT=1521', 'MPC_SANKHYA_ORACLE_CONNECT_STRING=fixture-service')
    $ambient = [ordered]@{ MPC_UNRELATED_SETTING = 'x'; MPC_DATABASE_URL = 'x'; ORACLE_HOME = 'x'; MELI_ACCESS_TOKEN = 'x'; SANKHYA_ORACLE_USER = 'x' }
    $original = [ordered]@{}
    try {
      foreach ($key in $ambient.Keys) {
        $original[$key] = [Environment]::GetEnvironmentVariable($key, 'Process')
        [Environment]::SetEnvironmentVariable($key, [string]$ambient[$key], 'Process')
      }
      $plan = Invoke-WithoutSankhyaCallerConnectionValues { New-LiveOracleDockerPlan -EnvFilePath $fixture }
      $startInfo = New-LiveOracleDockerProcessStartInfo -DockerPath 'C:\fixture\docker.exe' -Arguments $plan.RunArguments -Environment $plan.ContainerEnvironment
      $expectedRuntime = @('MPC_ORACLE_USERNAME', 'MPC_ORACLE_PASSWORD', 'MPC_ORACLE_CONNECT_STRING', 'MPC_ORACLE_LIVE_TEST', 'MPC_ORACLE_LIB_DIR')
      $expectedKeys = @($script:DockerExecutionEnvironmentKeys | Where-Object { -not [string]::IsNullOrWhiteSpace([Environment]::GetEnvironmentVariable($_, 'Process')) }) + $expectedRuntime
      $actualKeys = @($startInfo.Environment.Keys)

      Assert-RunnerCondition ((@($actualKeys | Sort-Object) -join ',') -eq (@($expectedKeys | Sort-Object) -join ',')) 'Docker child environment was not an execution-key plus runtime-key allowlist'
      Assert-RunnerCondition ((@($actualKeys | Where-Object { $_ -like 'MPC_*' }) -join ',') -eq ($expectedRuntime -join ',')) 'Docker child runtime environment was not exactly the five canonical keys'
      foreach ($key in $ambient.Keys) { Assert-RunnerCondition (-not ($actualKeys -contains $key)) "ambient key entered Docker child: $key" }
    } finally {
      foreach ($key in $ambient.Keys) { [Environment]::SetEnvironmentVariable($key, $original[$key], 'Process') }
      Remove-Item -LiteralPath $fixture -Force
    }
  }

  It 'preflights Docker through an execution-only child with no ambient or runtime keys' {
    $ambient = [ordered]@{ MPC_UNRELATED_SETTING = 'x'; MPC_SANKHYA_ORACLE_USERNAME = 'x'; MPC_DATABASE_URL = 'x'; ORACLE_HOME = 'x'; MELI_ACCESS_TOKEN = 'x'; SANKHYA_ORACLE_USER = 'x' }
    $original = [ordered]@{}
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
      foreach ($key in $ambient.Keys) { Assert-RunnerCondition (-not ($actualKeys -contains $key)) "ambient key entered Docker preflight child: $key" }
    } finally {
      foreach ($key in $ambient.Keys) { [Environment]::SetEnvironmentVariable($key, $original[$key], 'Process') }
    }
  }
}
