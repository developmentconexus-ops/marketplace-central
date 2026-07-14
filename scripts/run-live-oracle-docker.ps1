[CmdletBinding()]
param(
  [switch]$PreflightOnly,
  [switch]$PrepareImage,
  [switch]$EmitBaseline
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
  'TEMP', 'TMP', 'USERPROFILE', 'HOMEDRIVE', 'HOMEPATH',
  'APPDATA', 'LOCALAPPDATA', 'ProgramFiles', 'ProgramData'
)
$script:PhaseTimeouts = [ordered]@{
  preflight = [TimeSpan]::FromSeconds(30)
  build = [TimeSpan]::FromMinutes(10)
  run = [TimeSpan]::FromMinutes(3)
  inspect = [TimeSpan]::FromSeconds(30)
  promote = [TimeSpan]::FromSeconds(30)
}
$script:RuntimeFingerprintCache = $null

function Get-LiveOracleDockerProfile {
  Get-Content -Raw -LiteralPath $script:ProfilePath | ConvertFrom-Json -Depth 10
}

function Get-LiveOracleRuntimeFingerprint {
  param([Parameter(Mandatory)][object]$Profile)

  if (-not [string]::IsNullOrWhiteSpace($script:RuntimeFingerprintCache)) {
    return $script:RuntimeFingerprintCache
  }

  $inputs = [Collections.Generic.List[string]]::new()
  @(
    [string]$Profile.dockerfile,
    'go.work',
    'go.work.sum',
    'apps/server_core/go.mod',
    'apps/server_core/go.sum'
  ) | ForEach-Object { $inputs.Add($_) }
  Get-ChildItem -LiteralPath (Join-Path $script:RepositoryRoot 'apps/server_core') -Recurse -File -Filter '*.go' |
    Where-Object { $_.FullName -notmatch '[\\/]\.(gocache|gomodcache)[\\/]' } |
    ForEach-Object { [IO.Path]::GetRelativePath($script:RepositoryRoot, $_.FullName).Replace('\', '/') } |
    Sort-Object -Unique |
    ForEach-Object { $inputs.Add($_) }
  $manifest = @($inputs | ForEach-Object {
    $path = Join-Path $script:RepositoryRoot $_
    if (-not (Test-Path -LiteralPath $path -PathType Leaf)) { throw "live Oracle runtime input_missing=$_" }
    "$_=$((Get-FileHash -Algorithm SHA256 -LiteralPath $path).Hash.ToLowerInvariant())"
  }) -join "`n"
  $bytes = [Text.Encoding]::UTF8.GetBytes($manifest)
  $sha = [Security.Cryptography.SHA256]::Create()
  try {
    $script:RuntimeFingerprintCache = -join @($sha.ComputeHash($bytes) | ForEach-Object { $_.ToString('x2') })
    return $script:RuntimeFingerprintCache
  } finally { $sha.Dispose() }
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
  $fingerprint = Get-LiveOracleRuntimeFingerprint -Profile $profile
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
    RuntimeFingerprint = $fingerprint
    ContainerEnvironment = $containerEnvironment
    InspectArguments = @('image', 'inspect', '--format', '{{ index .Config.Labels "mpc.live-oracle.runtime-fingerprint" }}', [string]$profile.image_tag)
    BuildArguments = @('build', '--progress=plain', '--file', (Join-Path $script:RepositoryRoot ([string]$profile.dockerfile)), '--label', "mpc.live-oracle.runtime-fingerprint=$fingerprint", '--tag', [string]$profile.build_tag, $script:RepositoryRoot)
    PromoteArguments = @('tag', [string]$profile.build_tag, [string]$profile.image_tag)
    RunArguments = @(
      'run', '--rm', '--read-only', '--cap-drop', 'ALL', '--security-opt', 'no-new-privileges',
      '--tmpfs', '/tmp:rw,noexec,nosuid,size=64m'
    ) + @($script:ContainerKeys | ForEach-Object { @('--env', $_) }) + @(
      [string]$profile.image_tag, [string]$profile.test_binary,
      '-test.run', [string]$profile.go_test_regex, '-test.count=1', '-test.v'
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
    [Parameter(Mandatory)][string]$Phase,
    [TimeSpan]$Timeout = $script:PhaseTimeouts[$Phase]
  )

  if ($Timeout -le [TimeSpan]::Zero) { throw "live Oracle Docker $Phase invalid_timeout" }

  $process = [Diagnostics.Process]::new()
  $process.StartInfo = $StartInfo
  [void]$process.Start()

  # Both redirected pipes must be drained concurrently. Waiting on one stream
  # first can deadlock when Docker fills the other OS pipe buffer.
  $standardOutputTask = $process.StandardOutput.ReadToEndAsync()
  $standardErrorTask = $process.StandardError.ReadToEndAsync()
  $completed = $process.WaitForExit([Math]::Max(1, [int][Math]::Ceiling($Timeout.TotalMilliseconds)))
  if (-not $completed) {
    try { $process.Kill($true) } catch { try { $process.Kill() } catch { } }
    try { $process.WaitForExit() } catch { }
    try { [void]$standardOutputTask.GetAwaiter().GetResult() } catch { }
    try { [void]$standardErrorTask.GetAwaiter().GetResult() } catch { }
    throw "live Oracle Docker $Phase timed_out timeout_seconds=$([int][Math]::Ceiling($Timeout.TotalSeconds)); output suppressed"
  }

  $standardOutput = $standardOutputTask.GetAwaiter().GetResult()
  [void]$standardErrorTask.GetAwaiter().GetResult()
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
    [switch]$PrepareImage,
    [switch]$EmitBaseline,
    [string]$EnvFilePath = $script:LocalEnvPath
  )

  $plan = New-LiveOracleDockerPlan -EnvFilePath $EnvFilePath
  $dockerPath = Test-LiveOracleDockerAvailable
  if ($PreflightOnly) { return }

  if ($PrepareImage) {
    [void](Invoke-LiveOracleDockerCommand -DockerPath $dockerPath -Arguments $plan.BuildArguments -Environment @{} -Phase 'build')
    [void](Invoke-LiveOracleDockerCommand -DockerPath $dockerPath -Arguments $plan.PromoteArguments -Environment @{} -Phase 'promote')
    return
  }

  try {
    $observedFingerprint = (Invoke-LiveOracleDockerCommand -DockerPath $dockerPath -Arguments $plan.InspectArguments -Environment @{} -Phase 'inspect').Trim()
  } catch {
    throw 'live Oracle runtime image unavailable; run -PrepareImage'
  }
  if ($observedFingerprint -ne $plan.RuntimeFingerprint) {
    throw 'live Oracle runtime image stale; run -PrepareImage'
  }
  $runOutput = Invoke-LiveOracleDockerCommand -DockerPath $dockerPath -Arguments $plan.RunArguments -Environment $plan.ContainerEnvironment -Phase 'run'
  if ($EmitBaseline) {
    $baselineLines = Get-LiveOracleBaselineLines -RunOutput $runOutput
    if ($null -eq $baselineLines) {
      throw 'live Oracle Docker run failed exit_code=1; output suppressed'
    }
    return $baselineLines
  }
  if ($runOutput -notmatch '(?m)^MPC_C05_POSITIVE_CODPROD_OBSERVED=true\s*$') {
    throw 'live Oracle Docker run failed exit_code=1; output suppressed'
  }
}

# Get-LiveOracleBaselineLines extracts ONLY the three whitelisted baseline
# markers from suppressed run output. Each must match its strict shape exactly
# (integer or INDEX|FULL_SCAN); every other line of container output is
# discarded so no SQL text, bind values, DSN, credential, or driver message can
# surface. Returns the ordered sanitized lines, or $null if any is missing.
function Get-LiveOracleBaselineLines {
  param([Parameter(Mandatory)][AllowEmptyString()][string]$RunOutput)

  $patterns = [ordered]@{
    'MPC_BASELINE_TGFPRO_ACTIVE_COUNT' = '^MPC_BASELINE_TGFPRO_ACTIVE_COUNT=[0-9]+$'
    'MPC_BASELINE_RTT_MS'              = '^MPC_BASELINE_RTT_MS=[0-9]+$'
    'MPC_BASELINE_PAGE_PLAN'           = '^MPC_BASELINE_PAGE_PLAN=(INDEX|FULL_SCAN)$'
  }
  $lines = $RunOutput -split "`r?`n"
  $result = [System.Collections.Generic.List[string]]::new()
  foreach ($key in $patterns.Keys) {
    $match = $null
    foreach ($line in $lines) {
      $trimmed = $line.Trim()
      if ($trimmed -match $patterns[$key]) { $match = $trimmed; break }
    }
    if ($null -eq $match) { return $null }
    $result.Add($match)
  }
  return $result.ToArray()
}

function Write-LiveOracleDockerTelemetry {
  param(
    [Parameter(Mandatory)][ValidateSet('ready', 'passed', 'blocked')][string]$Status,
    [Parameter(Mandatory)][ValidateSet('preflight', 'prepared', 'complete', 'failed')][string]$Phase,
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
    '^live Oracle runtime image (unavailable|stale); run -PrepareImage$',
    '^live Oracle Docker (preflight|build|inspect|promote|run) failed exit_code=\d+; output suppressed$',
    '^live Oracle Docker (preflight|build|inspect|promote|run) timed_out timeout_seconds=\d+; output suppressed$'
  )
  if (@($safePatterns | Where-Object { $message -match $_ }).Count -gt 0) { return $message }
  return $null
}

function Invoke-LiveOracleDockerEntrypoint {
  param(
    [switch]$PreflightOnly,
    [switch]$PrepareImage,
    [switch]$EmitC05Evidence,
    [switch]$EmitBaseline,
    [string]$EnvFilePath = $script:LocalEnvPath,
    [Parameter(Mandatory)][ref]$ExitCode
  )

  try {
    $runnerResult = Invoke-LiveOracleDockerRunner -PreflightOnly:$PreflightOnly -PrepareImage:$PrepareImage -EmitBaseline:$EmitBaseline -EnvFilePath $EnvFilePath
    if ($PreflightOnly) {
      Write-LiveOracleDockerTelemetry -Status 'ready' -Phase 'preflight' -ExitCode 0
    } elseif ($PrepareImage) {
      Write-LiveOracleDockerTelemetry -Status 'ready' -Phase 'prepared' -ExitCode 0
    } else {
      Write-LiveOracleDockerTelemetry -Status 'passed' -Phase 'complete' -ExitCode 0
      if ($EmitBaseline -and $null -ne $runnerResult) {
        foreach ($line in $runnerResult) { Write-Output $line }
      }
      if ($EmitC05Evidence -and -not $EmitBaseline) {
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
  Invoke-LiveOracleDockerEntrypoint -PreflightOnly:$PreflightOnly -PrepareImage:$PrepareImage -EmitC05Evidence -EmitBaseline:$EmitBaseline -ExitCode ([ref]$exitCode)
  exit $exitCode
}
