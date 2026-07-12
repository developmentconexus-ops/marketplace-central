[CmdletBinding()]
param(
  [switch]$PreflightOnly
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$script:RepositoryRoot = Split-Path -Parent $PSScriptRoot
$script:ProfilePath = Join-Path $script:RepositoryRoot 'docker/live-oracle/profile.json'
$script:LocalEnvPath = Join-Path $script:RepositoryRoot '.env'
$script:CallerCredentialKeys = @('MPC_SANKHYA_ORACLE_USERNAME', 'MPC_SANKHYA_ORACLE_PASSWORD', 'MPC_SANKHYA_ORACLE_CONNECT_STRING')
$script:ContainerCredentialKeys = @('MPC_ORACLE_USERNAME', 'MPC_ORACLE_PASSWORD', 'MPC_ORACLE_CONNECT_STRING')
$script:ContainerKeys = @($script:ContainerCredentialKeys + @('MPC_ORACLE_LIVE_TEST', 'MPC_ORACLE_LIB_DIR'))
$script:DockerExecutionEnvironmentKeys = @(
  'SystemRoot', 'WINDIR', 'ComSpec', 'PATH', 'PATHEXT',
  'TEMP', 'TMP', 'USERPROFILE', 'APPDATA', 'LOCALAPPDATA'
)

function Get-LiveOracleDockerProfile {
  Get-Content -Raw -LiteralPath $script:ProfilePath | ConvertFrom-Json -Depth 10
}

function Get-LiveOracleLocalEnvValues {
  param([string]$EnvFilePath = $script:LocalEnvPath)

  $values = [ordered]@{}
  if (-not (Test-Path -LiteralPath $EnvFilePath -PathType Leaf)) { return $values }

  $lineNumber = 0
  foreach ($line in Get-Content -LiteralPath $EnvFilePath) {
    $lineNumber += 1
    $trimmed = $line.Trim()
    if ([string]::IsNullOrWhiteSpace($trimmed) -or $trimmed.StartsWith('#')) { continue }
    if ($trimmed -notmatch '^(?<key>[A-Za-z_][A-Za-z0-9_]*)=(?<value>.*)$') {
      throw "live Oracle .env invalid_line=$lineNumber"
    }

    $key = $Matches.key
    if ($key -notin $script:CallerCredentialKeys) { throw "live Oracle .env unsupported_key=$key" }
    if ($values.Contains($key)) { throw "live Oracle .env duplicate_key=$key" }

    $value = $Matches.value.Trim()
    if ($value.Length -ge 2 -and (($value.StartsWith('"') -and $value.EndsWith('"')) -or ($value.StartsWith("'") -and $value.EndsWith("'")))) {
      $value = $value.Substring(1, $value.Length - 2)
    }
    $values[$key] = $value
  }
  $values
}

function Get-LiveOracleCredentialValues {
  param([string]$EnvFilePath = $script:LocalEnvPath)

  $localValues = Get-LiveOracleLocalEnvValues -EnvFilePath $EnvFilePath
  $values = [ordered]@{}
  $missing = [Collections.Generic.List[string]]::new()
  foreach ($key in $script:CallerCredentialKeys) {
    $processValue = [Environment]::GetEnvironmentVariable($key, 'Process')
    $value = if (-not [string]::IsNullOrWhiteSpace($processValue)) { $processValue } else { $localValues[$key] }
    if ([string]::IsNullOrWhiteSpace($value)) { $missing.Add($key); continue }
    $values[$key] = $value
  }
  if ($missing.Count -gt 0) { throw "live Oracle preflight missing_keys=$($missing -join ',')" }
  $values
}

function New-LiveOracleDockerPlan {
  param([string]$EnvFilePath = $script:LocalEnvPath)

  $profile = Get-LiveOracleDockerProfile
  $credentials = Get-LiveOracleCredentialValues -EnvFilePath $EnvFilePath
  $containerEnvironment = [ordered]@{
    MPC_ORACLE_USERNAME = $credentials.MPC_SANKHYA_ORACLE_USERNAME
    MPC_ORACLE_PASSWORD = $credentials.MPC_SANKHYA_ORACLE_PASSWORD
    MPC_ORACLE_CONNECT_STRING = $credentials.MPC_SANKHYA_ORACLE_CONNECT_STRING
    MPC_ORACLE_LIVE_TEST = '1'
    MPC_ORACLE_LIB_DIR = [string]$profile.oracle_lib_dir
  }
  [pscustomobject]@{
    Profile = $profile
    ContainerEnvironment = $containerEnvironment
    BuildArguments = @('build', '--file', (Join-Path $script:RepositoryRoot ([string]$profile.dockerfile)), '--tag', [string]$profile.image_tag, $script:RepositoryRoot)
    RunArguments = @(
      'run', '--rm', '--mount', "type=bind,source=$script:RepositoryRoot,target=$($profile.workspace),readonly",
      '--workdir', "$($profile.workspace)/apps/server_core"
    ) + @($script:ContainerKeys | ForEach-Object { @('--env', $_) }) + @(
      [string]$profile.image_tag, 'go', 'test', [string]$profile.go_package,
      '-run', [string]$profile.go_test_regex, '-count=1'
    )
  }
}

function Test-LiveOracleDockerAvailable {
  $docker = Get-Command docker -CommandType Application -ErrorAction SilentlyContinue | Select-Object -First 1
  if ($null -eq $docker) { throw 'live Oracle preflight docker_unavailable' }
  $preflight = New-LiveOracleDockerPreflightPlan -DockerPath ([string]$docker.Source)
  Invoke-LiveOracleDockerProcess -StartInfo $preflight.StartInfo -Phase 'preflight'
  [string]$docker.Source
}

function New-LiveOracleDockerProcessStartInfo {
  param(
    [Parameter(Mandatory)][string]$DockerPath,
    [Parameter(Mandatory)][string[]]$Arguments,
    [Parameter(Mandatory)][System.Collections.IDictionary]$Environment
  )

  $startInfo = [Diagnostics.ProcessStartInfo]::new()
  $startInfo.FileName = $DockerPath
  $startInfo.UseShellExecute = $false
  $startInfo.RedirectStandardOutput = $true
  $startInfo.RedirectStandardError = $true
  foreach ($argument in $Arguments) { [void]$startInfo.ArgumentList.Add($argument) }

  # Docker must not inherit ambient application, provider, database, or legacy
  # configuration. Retain only the OS support values Docker needs to execute,
  # then add the explicitly allowlisted runtime values for the run invocation.
  $startInfo.Environment.Clear()
  foreach ($key in $script:DockerExecutionEnvironmentKeys) {
    $value = [Environment]::GetEnvironmentVariable($key, 'Process')
    if (-not [string]::IsNullOrWhiteSpace($value)) { $startInfo.Environment[$key] = $value }
  }
  foreach ($key in $Environment.Keys) { $startInfo.Environment[$key] = [string]$Environment[$key] }
  $startInfo
}

function New-LiveOracleDockerPreflightPlan {
  param([Parameter(Mandatory)][string]$DockerPath)

  $arguments = @('version', '--format', '{{.Server.Version}}')
  [pscustomobject]@{
    Arguments = $arguments
    StartInfo = New-LiveOracleDockerProcessStartInfo -DockerPath $DockerPath -Arguments $arguments -Environment ([ordered]@{})
  }
}

function Invoke-LiveOracleDockerProcess {
  param(
    [Parameter(Mandatory)][Diagnostics.ProcessStartInfo]$StartInfo,
    [Parameter(Mandatory)][string]$Phase
  )

  $process = [Diagnostics.Process]::new()
  $process.StartInfo = $StartInfo
  [void]$process.Start()
  [void]$process.StandardOutput.ReadToEnd()
  [void]$process.StandardError.ReadToEnd()
  $process.WaitForExit()
  if ($process.ExitCode -ne 0) { throw "live Oracle Docker $Phase failed exit_code=$($process.ExitCode); output suppressed" }
}

function Invoke-LiveOracleDockerCommand {
  param(
    [Parameter(Mandatory)][string]$DockerPath,
    [Parameter(Mandatory)][string[]]$Arguments,
    [Parameter(Mandatory)][System.Collections.IDictionary]$Environment,
    [Parameter(Mandatory)][string]$Phase
  )

  $startInfo = New-LiveOracleDockerProcessStartInfo -DockerPath $DockerPath -Arguments $Arguments -Environment $Environment
  Invoke-LiveOracleDockerProcess -StartInfo $startInfo -Phase $Phase
}

function Invoke-LiveOracleDockerRunner {
  param([switch]$PreflightOnly)

  $dockerPath = Test-LiveOracleDockerAvailable
  $plan = New-LiveOracleDockerPlan
  if ($PreflightOnly) { Write-Output 'status=ready'; Write-Output 'target=live-oracle'; Write-Output 'runner=docker'; return }

  Invoke-LiveOracleDockerCommand -DockerPath $dockerPath -Arguments $plan.BuildArguments -Environment @{} -Phase 'build'
  Invoke-LiveOracleDockerCommand -DockerPath $dockerPath -Arguments $plan.RunArguments -Environment $plan.ContainerEnvironment -Phase 'run'
  Write-Output 'status=passed'; Write-Output 'target=live-oracle'; Write-Output 'runner=docker'
}

if ($MyInvocation.InvocationName -ne '.') {
  try {
    Invoke-LiveOracleDockerRunner -PreflightOnly:$PreflightOnly
  } catch {
    Write-Output 'status=blocked'
    Write-Output $_.Exception.Message
    exit 1
  }
}
