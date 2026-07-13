[CmdletBinding()]
param(
  [switch]$PreflightOnly
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$script:RepositoryRoot = Split-Path -Parent $PSScriptRoot
$script:ProfilePath = Join-Path $script:RepositoryRoot 'docker/live-oracle/profile.json'
$script:LocalEnvPath = Join-Path $script:RepositoryRoot '.env'
$script:CallerConnectionKeys = @(
  'MPC_SANKHYA_ORACLE_USERNAME',
  'MPC_SANKHYA_ORACLE_PASSWORD',
  'MPC_SANKHYA_ORACLE_HOST',
  'MPC_SANKHYA_ORACLE_PORT',
  'MPC_SANKHYA_ORACLE_CONNECT_STRING'
)
$script:IgnoredLocalEnvKeys = @('MPC_SANKHYA_ORACLE_SCHEMA')
$script:AllowedLocalEnvKeys = @($script:CallerConnectionKeys + $script:IgnoredLocalEnvKeys)
$script:ReservedLocalEnvPrefix = 'MPC_SANKHYA_ORACLE_'
$script:RejectedLocalEnvAliasPrefixes = @('MPC_ORACLE_', 'SANKHYA_ORACLE_', 'ORACLE_')
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
  $seenKeys = [Collections.Generic.HashSet[string]]::new([StringComparer]::Ordinal)
  if (-not (Test-Path -LiteralPath $EnvFilePath -PathType Leaf)) { return $values }

  $lineNumber = 0
  foreach ($line in Get-Content -LiteralPath $EnvFilePath) {
    $lineNumber += 1
    $trimmed = $line.Trim()
    if ([string]::IsNullOrWhiteSpace($trimmed) -or $trimmed.StartsWith('#')) { continue }
    if ($trimmed -notmatch '^(?<key>[A-Za-z_][A-Za-z0-9_]*)=') {
      throw "live Oracle .env invalid_line=$lineNumber"
    }

    $key = $Matches.key
    if (-not $key.StartsWith($script:ReservedLocalEnvPrefix, [StringComparison]::Ordinal)) {
      if (@($script:RejectedLocalEnvAliasPrefixes | Where-Object { $key.StartsWith($_, [StringComparison]::Ordinal) }).Count -gt 0) {
        throw "live Oracle .env unsupported_key=$key"
      }
      continue
    }
    if ($key -notin $script:AllowedLocalEnvKeys) { throw "live Oracle .env unsupported_key=$key" }
    if (-not $seenKeys.Add($key)) { throw "live Oracle .env duplicate_key=$key" }

    if ($key -in $script:IgnoredLocalEnvKeys) { continue }

    $value = $trimmed.Substring($key.Length + 1).Trim()
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
  foreach ($key in $script:CallerConnectionKeys) {
    $processValue = [Environment]::GetEnvironmentVariable($key, 'Process')
    $value = if (-not [string]::IsNullOrWhiteSpace($processValue)) { $processValue.Trim() } else { $localValues[$key] }
    if ([string]::IsNullOrWhiteSpace($value)) { $missing.Add($key); continue }
    $values[$key] = $value
  }
  if ($missing.Count -gt 0) { throw "live Oracle preflight missing_keys=$($missing -join ',')" }
  $values
}

function New-LiveOracleDockerPlan {
  param([string]$EnvFilePath = $script:LocalEnvPath)

  $profile = Get-LiveOracleDockerProfile
  $connection = Get-LiveOracleCredentialValues -EnvFilePath $EnvFilePath
  $containerEnvironment = [ordered]@{
    MPC_ORACLE_USERNAME = $connection.MPC_SANKHYA_ORACLE_USERNAME
    MPC_ORACLE_PASSWORD = $connection.MPC_SANKHYA_ORACLE_PASSWORD
    MPC_ORACLE_CONNECT_STRING = "$($connection.MPC_SANKHYA_ORACLE_HOST):$($connection.MPC_SANKHYA_ORACLE_PORT)/$($connection.MPC_SANKHYA_ORACLE_CONNECT_STRING)"
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
      '-run', [string]$profile.go_test_regex, '-count=1', '-v'
    )
  }
}

function Test-LiveOracleDockerAvailable {
  $docker = Get-Command docker -CommandType Application -ErrorAction SilentlyContinue | Select-Object -First 1
  if ($null -eq $docker) { throw 'live Oracle preflight docker_unavailable' }
  $preflight = New-LiveOracleDockerPreflightPlan -DockerPath ([string]$docker.Source)
  [void](Invoke-LiveOracleDockerProcess -StartInfo $preflight.StartInfo -Phase 'preflight')
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

  $standardOutput = $process.StandardOutput.ReadToEnd()
  [void]$process.StandardError.ReadToEnd()
  $process.WaitForExit()
  if ($process.ExitCode -ne 0) { throw "live Oracle Docker $Phase failed exit_code=$($process.ExitCode); output suppressed" }

  return $standardOutput
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
  param(
    [switch]$PreflightOnly,
    [string]$EnvFilePath = $script:LocalEnvPath
  )

  $plan = New-LiveOracleDockerPlan -EnvFilePath $EnvFilePath
  $dockerPath = Test-LiveOracleDockerAvailable
  if ($PreflightOnly) { return }

  [void](Invoke-LiveOracleDockerCommand -DockerPath $dockerPath -Arguments $plan.BuildArguments -Environment @{} -Phase 'build')
  $runOutput = Invoke-LiveOracleDockerCommand -DockerPath $dockerPath -Arguments $plan.RunArguments -Environment $plan.ContainerEnvironment -Phase 'run'
  if ($runOutput -notmatch '(?m)^MPC_C05_POSITIVE_CODPROD_OBSERVED=true\s*$') {
    throw 'live Oracle Docker run failed exit_code=1; output suppressed'
  }
}

function Write-LiveOracleDockerTelemetry {
  param(
    [Parameter(Mandatory)][ValidateSet('ready', 'passed', 'blocked')][string]$Status,
    [Parameter(Mandatory)][ValidateSet('preflight', 'complete', 'failed')][string]$Phase,
    [Parameter(Mandatory)][int]$ExitCode,
    [string]$Reason
  )

  Write-Output "status=$Status"
  Write-Output "phase=$Phase"
  Write-Output "exit_code=$ExitCode"
  if (-not [string]::IsNullOrWhiteSpace($Reason)) { Write-Output "reason=$Reason" }
}

function Get-LiveOracleDockerSafeReason {
  param([Parameter(Mandatory)][System.Management.Automation.ErrorRecord]$ErrorRecord)

  $message = $ErrorRecord.Exception.Message
  $safePatterns = @(
    '^live Oracle \.env invalid_line=\d+$',
    '^live Oracle \.env unsupported_key=[A-Za-z_][A-Za-z0-9_]*$',
    '^live Oracle preflight missing_keys=[A-Za-z0-9_,]+$',
    '^live Oracle preflight docker_unavailable$',
    '^live Oracle Docker (preflight|build|run) failed exit_code=\d+; output suppressed$'
  )
  if (@($safePatterns | Where-Object { $message -match $_ }).Count -gt 0) { return $message }
  return $null
}

function Invoke-LiveOracleDockerEntrypoint {
  param(
    [switch]$PreflightOnly,
    [switch]$EmitC05Evidence,
    [string]$EnvFilePath = $script:LocalEnvPath,
    [Parameter(Mandatory)][ref]$ExitCode
  )

  try {
    Invoke-LiveOracleDockerRunner -PreflightOnly:$PreflightOnly -EnvFilePath $EnvFilePath
    if ($PreflightOnly) {
      Write-LiveOracleDockerTelemetry -Status 'ready' -Phase 'preflight' -ExitCode 0
    } else {
      Write-LiveOracleDockerTelemetry -Status 'passed' -Phase 'complete' -ExitCode 0
      if ($EmitC05Evidence) {
        Write-Output "frozen_sha=$(git -C $script:RepositoryRoot rev-parse HEAD)"
        Write-Output 'source=oracle/sankhya'
        Write-Output "observed_at=$([DateTimeOffset]::UtcNow.ToString('o'))"
        Write-Output 'read_only=true'
        Write-Output 'positive_codprod_observed=true'
      }
    }
    $ExitCode.Value = 0
  } catch {
    $reason = Get-LiveOracleDockerSafeReason -ErrorRecord $_
    Write-LiveOracleDockerTelemetry -Status 'blocked' -Phase 'failed' -ExitCode 1 -Reason $reason
    $ExitCode.Value = 1
  }
}

if ($MyInvocation.InvocationName -ne '.') {
  $exitCode = 0
  Invoke-LiveOracleDockerEntrypoint -PreflightOnly:$PreflightOnly -EmitC05Evidence -ExitCode ([ref]$exitCode)
  exit $exitCode
}
