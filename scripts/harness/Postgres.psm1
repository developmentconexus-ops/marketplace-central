Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

Import-Module (Join-Path $PSScriptRoot 'Execution.psm1') -Force

function Copy-HarnessEnvironment {
  param([Parameter(Mandatory)][System.Collections.IDictionary]$Source)
  $copy = [System.Collections.Generic.Dictionary[string,string]]::new([StringComparer]::OrdinalIgnoreCase)
  foreach ($key in $Source.Keys) { $copy[[string]$key] = [string]$Source[$key] }
  return $copy
}

function Resolve-HarnessPostgresDockerApplication {
  [CmdletBinding()]
  param([string]$Name = 'docker')

  $application = Get-Command $Name -CommandType Application -ErrorAction SilentlyContinue | Select-Object -First 1
  if ($null -eq $application -or [string]::IsNullOrWhiteSpace([string]$application.Source)) {
    throw 'HPG_DOCKER_MISSING'
  }
  return [IO.Path]::GetFullPath([string]$application.Source)
}

function New-HarnessPostgresRunSpec {
  [CmdletBinding()]
  param(
    [Parameter(Mandatory)][string]$RepositoryRoot,
    [Parameter(Mandatory)][AllowEmptyString()][string]$RunId,
    [Parameter(Mandatory)][string]$Password,
    [Parameter(Mandatory)][string]$DockerFilePath,
    [string[]]$DockerArgumentPrefix = @()
  )

  if ($RunId -cnotmatch '^[0-9a-f]{32}$') { throw 'HPG_RUN_ID_INVALID' }
  if (-not [IO.Path]::IsPathFullyQualified($RepositoryRoot)) { throw 'HPG_RUN_ID_INVALID repository_root' }
  $root = [IO.Path]::GetFullPath($RepositoryRoot)
  if (-not (Test-Path -LiteralPath $root -PathType Container)) { throw 'HPG_RUN_ID_INVALID repository_root' }
  if (-not [IO.Path]::IsPathFullyQualified($DockerFilePath) -or -not (Test-Path -LiteralPath $DockerFilePath -PathType Leaf)) {
    throw 'HPG_DOCKER_MISSING'
  }
  if ([string]::IsNullOrWhiteSpace($Password)) { throw 'HPG_RUN_ID_INVALID password' }

  return [pscustomobject]@{
    RunId = $RunId
    DatabaseName = "mpc_test_$RunId"
    ContainerName = "mpc-pg-$RunId"
    Label = "marketplace-central.harness.run=$RunId"
    RepositoryRoot = $root
    Password = $Password
    DockerFilePath = [IO.Path]::GetFullPath($DockerFilePath)
    DockerArgumentPrefix = @($DockerArgumentPrefix)
  }
}

function Invoke-HarnessPostgresProcess {
  param(
    [Parameter(Mandatory)][string]$FilePath,
    [Parameter(Mandatory)][AllowEmptyCollection()][string[]]$ArgumentPrefix,
    [Parameter(Mandatory)][string[]]$Arguments,
    [Parameter(Mandatory)][string]$WorkingDirectory,
    [Parameter(Mandatory)][System.Collections.IDictionary]$Environment,
    [Parameter(Mandatory)][int]$TimeoutSeconds,
    [string[]]$RedactionCandidates = @()
  )
  $request = New-HarnessProcessRequest -FilePath $FilePath -ArgumentList (@($ArgumentPrefix) + @($Arguments)) -WorkingDirectory $WorkingDirectory -Environment $Environment -TimeoutSeconds $TimeoutSeconds -RedactionCandidates (@($RedactionCandidates) + @($WorkingDirectory))
  return Invoke-HarnessProcess -Request $request
}

function Get-HarnessMigrationCount {
  param([string]$Output)
  if ($Output -match '(?m)^applied\s+(\d+)\s+migration\(s\)\s*$') { return [int]$Matches[1] }
  return -1
}

function Get-HarnessPostgresFailureTokens {
  param([AllowNull()][string]$Text)

  $tokens = [System.Collections.Generic.HashSet[string]]::new([StringComparer]::Ordinal)
  foreach ($match in [regex]::Matches([string]$Text, '\bHPG_[A-Z0-9_]{1,64}\b')) { [void]$tokens.Add("hpg=$($match.Value)") }
  foreach ($match in [regex]::Matches([string]$Text, '\bSQLSTATE\s+(?<code>[0-9A-Z]{5})\b')) { [void]$tokens.Add("sqlstate=$($match.Groups['code'].Value)") }
  foreach ($match in [regex]::Matches([string]$Text, '(?m)^FAIL\s+(?<package>marketplace-central/[A-Za-z0-9_./-]{1,240})(?:\s|$)')) { [void]$tokens.Add("package=$($match.Groups['package'].Value)") }
  foreach ($match in [regex]::Matches([string]$Text, '(?m)^--- FAIL: (?<test>Test[A-Za-z0-9_]{1,160})(?:\s|$)')) { [void]$tokens.Add("test=$($match.Groups['test'].Value)") }
  $safe = @($tokens | Sort-Object | Select-Object -First 24)
  if ($safe.Count -eq 0) { return @('HPG_CHILD_OUTPUT_REDACTED') }
  return $safe
}

function Invoke-HarnessPostgresLifecycle {
  [CmdletBinding()]
  param(
    [Parameter(Mandatory)][object]$RunSpec,
    [Parameter(Mandatory)][System.Collections.IDictionary]$BaseEnvironment,
    [Parameter(Mandatory)][string]$GoFilePath,
    [string[]]$GoArgumentPrefix = @(),
    [ValidateRange(1, 3600)][int]$TimeoutSeconds = 600,
    [ValidateRange(1, 600)][int]$ReadyMaxAttempts = 60,
    [ValidateRange(0, 60000)][int]$ReadyRetryDelayMilliseconds = 1000,
    [ValidateRange(1, 600000)][int]$ReadyTimeoutMilliseconds = 60000,
    [string[]]$TestArguments = @('test', '-tags=integration', './tests/integration', './internal/modules/orders/adapters/postgres', './internal/modules/profitability/adapters/postgres', './internal/modules/product_links/application', '-count=1'),
    [switch]$HoldConnectionDuringCleanupTest
  )

  $dockerEnvironment = Copy-HarnessEnvironment $BaseEnvironment
  $dockerEnvironment['POSTGRES_PASSWORD'] = [string]$RunSpec.Password
  $goEnvironment = Copy-HarnessEnvironment $BaseEnvironment
  $cleanupCodes = [System.Collections.Generic.List[string]]::new()
  $inventory = @()
  $state = @{ PrimaryCode = ''; PrimaryExit = 0 }
  $databaseCreated = $false
  $containerStartAttempted = $false
  $containerOwned = $false
  $firstCount = -1
  $secondCount = -1
  $hostPort = 0
  $targetURL = ''
  $failureDiagnosticTokens = @()
  $heldConnectionConfirmed = $false

  function Invoke-Docker {
    param(
      [Parameter(Mandatory)][string[]]$Arguments,
      [ValidateRange(1, 3600)][int]$ProcessTimeoutSeconds = $TimeoutSeconds
    )
    return Invoke-HarnessPostgresProcess -FilePath $RunSpec.DockerFilePath -ArgumentPrefix @($RunSpec.DockerArgumentPrefix) -Arguments $Arguments -WorkingDirectory $RunSpec.RepositoryRoot -Environment $dockerEnvironment -TimeoutSeconds $ProcessTimeoutSeconds -RedactionCandidates @($RunSpec.Password, $targetURL, $RunSpec.RepositoryRoot)
  }
  function Invoke-Go([string[]]$Arguments) {
    return Invoke-HarnessPostgresProcess -FilePath $GoFilePath -ArgumentPrefix @($GoArgumentPrefix) -Arguments $Arguments -WorkingDirectory (Join-Path $RunSpec.RepositoryRoot 'apps/server_core') -Environment $goEnvironment -TimeoutSeconds $TimeoutSeconds -RedactionCandidates @($RunSpec.Password, $targetURL, $RunSpec.RepositoryRoot)
  }
  function Set-Primary([string]$Code, [int]$ExitCode) {
    if ([string]::IsNullOrWhiteSpace([string]$state.PrimaryCode)) {
      $state.PrimaryCode = $Code
      $state.PrimaryExit = if ($ExitCode -ne 0) { $ExitCode } else { 1 }
    }
  }
  function Add-CleanupCode([string]$Code) {
    if (-not $cleanupCodes.Contains($Code)) { [void]$cleanupCodes.Add($Code) }
  }
  function Get-ResourceNames([string]$Filter) {
    $query = Invoke-Docker @('ps', '--all', '--filter', $Filter, '--format', '{{.Names}}')
    if ($query.ExitCode -ne 0) { return [pscustomobject]@{ Passed = $false; Names = @() } }
    $names = @($query.Stdout -split '\r?\n' | Where-Object { -not [string]::IsNullOrWhiteSpace($_) })
    return [pscustomobject]@{ Passed = $true; Names = $names }
  }

  try {
    do {
    $daemon = Invoke-Docker @('version', '--format', '{{.Server.Version}}')
    if ($daemon.ExitCode -ne 0) { Set-Primary 'HPG_DOCKER_UNAVAILABLE' $daemon.ExitCode; break }

    $image = Invoke-Docker @('image', 'inspect', 'postgres:16-bookworm', '--format', '{{.Id}}')
    if ($image.ExitCode -ne 0) { Set-Primary 'HPG_IMAGE_MISSING' $image.ExitCode; break }

    $nameBefore = Get-ResourceNames "name=^/$([regex]::Escape($RunSpec.ContainerName))$"
    $labelBefore = Get-ResourceNames "label=$($RunSpec.Label)"
    if (-not $nameBefore.Passed -or -not $labelBefore.Passed) { Set-Primary 'HPG_DOCKER_UNAVAILABLE' 1; break }
    if (@($nameBefore.Names).Count -gt 0 -or @($labelBefore.Names).Count -gt 0) { Set-Primary 'HPG_RESOURCE_CONFLICT' 1; break }

    $containerStartAttempted = $true
    $start = Invoke-Docker @(
      'run', '--detach', '--rm', '--pull=never',
      '--name', $RunSpec.ContainerName,
      '--label', $RunSpec.Label,
      '--publish', '127.0.0.1::5432',
      '--tmpfs', '/var/lib/postgresql/data:rw,noexec,nosuid,size=512m',
      '--env', 'POSTGRES_PASSWORD',
      '--env', 'POSTGRES_DB=postgres',
      'postgres:16-bookworm'
    )
    if ($start.ExitCode -ne 0) { Set-Primary 'HPG_CONTAINER_START_FAILED' $start.ExitCode; break }

    $ownership = Invoke-Docker @('inspect', '--format', '{{.Name}}|{{ index .Config.Labels "marketplace-central.harness.run" }}', $RunSpec.ContainerName)
    $expectedOwnership = "/$($RunSpec.ContainerName)|$($RunSpec.RunId)"
    if ($ownership.ExitCode -ne 0 -or $ownership.Stdout.Trim() -cne $expectedOwnership) {
      Set-Primary 'HPG_CONTAINER_START_FAILED' $ownership.ExitCode
      break
    }
    $containerOwned = $true

    $ready = $null
    $readyWatch = [Diagnostics.Stopwatch]::StartNew()
    for ($attempt = 1; $attempt -le $ReadyMaxAttempts; $attempt++) {
      $remainingMilliseconds = [Math]::Max(1, $ReadyTimeoutMilliseconds - $readyWatch.ElapsedMilliseconds)
      $readyProcessTimeoutSeconds = [Math]::Min($TimeoutSeconds, [Math]::Max(1, [Math]::Ceiling($remainingMilliseconds / 1000.0)))
      $ready = Invoke-Docker -ProcessTimeoutSeconds $readyProcessTimeoutSeconds -Arguments @('exec', $RunSpec.ContainerName, 'pg_isready', '--username', 'postgres', '--dbname', 'postgres', '--timeout', [string]$readyProcessTimeoutSeconds)
      if ($ready.ExitCode -eq 0) { break }
      if ($attempt -ge $ReadyMaxAttempts -or $readyWatch.ElapsedMilliseconds -ge $ReadyTimeoutMilliseconds) { break }
      if ($ReadyRetryDelayMilliseconds -gt 0) { Start-Sleep -Milliseconds $ReadyRetryDelayMilliseconds }
    }
    $readyWatch.Stop()
    if ($null -eq $ready -or $ready.ExitCode -ne 0) { Set-Primary 'HPG_READY_TIMEOUT' $(if ($null -eq $ready) { 1 } else { $ready.ExitCode }); break }

    $port = Invoke-Docker @('port', $RunSpec.ContainerName, '5432/tcp')
    if ($port.ExitCode -ne 0 -or $port.Stdout -notmatch '(?m)^127\.0\.0\.1:(\d+)\s*$') { Set-Primary 'HPG_PORT_UNAVAILABLE' $port.ExitCode; break }
    $hostPort = [int]$Matches[1]
    if ($hostPort -lt 1 -or $hostPort -gt 65535) { Set-Primary 'HPG_PORT_UNAVAILABLE' 1; break }
    $escapedPassword = [uri]::EscapeDataString([string]$RunSpec.Password)
    $targetURL = "postgresql://postgres:$escapedPassword@127.0.0.1:$hostPort/$($RunSpec.DatabaseName)?sslmode=disable"
    $goEnvironment['MPC_TEST_DATABASE_URL'] = $targetURL

    $create = Invoke-Docker @('exec', $RunSpec.ContainerName, 'psql', '--username', 'postgres', '--dbname', 'postgres', '--set', 'ON_ERROR_STOP=1', '--command', "CREATE DATABASE $($RunSpec.DatabaseName)")
    if ($create.ExitCode -ne 0) { Set-Primary 'HPG_DATABASE_CREATE_FAILED' $create.ExitCode; break }
    $databaseCreated = $true

    $migrationArgs = @('run', './cmd/testdb', 'migrate')
    $migration1 = Invoke-Go $migrationArgs
    if ($migration1.ExitCode -ne 0) { Set-Primary 'HPG_MIGRATION_FAILED' $migration1.ExitCode; break }
    $firstCount = Get-HarnessMigrationCount $migration1.Stdout
    $migration2 = Invoke-Go $migrationArgs
    if ($migration2.ExitCode -ne 0) { Set-Primary 'HPG_MIGRATION_FAILED' $migration2.ExitCode; break }
    $secondCount = Get-HarnessMigrationCount $migration2.Stdout
    if ($firstCount -ne 32 -or $secondCount -ne 0) { Set-Primary 'HPG_MIGRATION_NOT_IDEMPOTENT' 1; break }

    if ($HoldConnectionDuringCleanupTest) {
      $heldConnection = Invoke-Docker @('exec', '--detach', $RunSpec.ContainerName, 'psql', '--username', 'postgres', '--dbname', $RunSpec.DatabaseName, '--command', 'SELECT pg_sleep(300)')
      if ($heldConnection.ExitCode -ne 0) { Set-Primary 'HPG_TEST_FAILED' $heldConnection.ExitCode; break }
      for ($heldAttempt = 1; $heldAttempt -le 20; $heldAttempt++) {
        $heldCheck = Invoke-Docker @('exec', $RunSpec.ContainerName, 'psql', '--username', 'postgres', '--dbname', 'postgres', '--quiet', '--tuples-only', '--no-align', '--set', 'ON_ERROR_STOP=1', '--command', "SELECT count(*) FROM pg_stat_activity WHERE datname = '$($RunSpec.DatabaseName)' AND state = 'active' AND query LIKE 'SELECT pg_sleep(300)%'")
        if ($heldCheck.ExitCode -eq 0 -and $heldCheck.Stdout.Trim() -match '^[1-9][0-9]*$') { $heldConnectionConfirmed = $true; break }
        Start-Sleep -Milliseconds 100
      }
      if (-not $heldConnectionConfirmed) { Set-Primary 'HPG_TEST_FAILED' 1; break }
    }
    $tests = Invoke-Go @($TestArguments)
    if ($tests.ExitCode -ne 0) {
      $failureDiagnosticTokens = @(Get-HarnessPostgresFailureTokens (@($tests.Stdout, $tests.Stderr) -join "`n"))
      Set-Primary 'HPG_TEST_FAILED' $tests.ExitCode
      break
    }
    } while ($false)
  } finally {
    if ($databaseCreated -and $containerOwned) {
      $drop = Invoke-Docker @('exec', $RunSpec.ContainerName, 'psql', '--username', 'postgres', '--dbname', 'postgres', '--set', 'ON_ERROR_STOP=1', '--command', "DROP DATABASE $($RunSpec.DatabaseName) WITH (FORCE)")
      if ($drop.ExitCode -ne 0) { Add-CleanupCode 'HPG_DATABASE_DROP_FAILED' }
      $dropVerify = Invoke-Docker @('exec', $RunSpec.ContainerName, 'psql', '--username', 'postgres', '--dbname', 'postgres', '--tuples-only', '--no-align', '--set', 'ON_ERROR_STOP=1', '--command', "SELECT datname FROM pg_database WHERE datname = '$($RunSpec.DatabaseName)'")
      if ($dropVerify.ExitCode -ne 0 -or -not [string]::IsNullOrWhiteSpace($dropVerify.Stdout)) { Add-CleanupCode 'HPG_DATABASE_DROP_FAILED' }
    }
    if ($containerOwned) {
      $remove = Invoke-Docker @('rm', '--force', $RunSpec.ContainerName)
      if ($remove.ExitCode -ne 0) { Add-CleanupCode 'HPG_CONTAINER_REMOVE_FAILED' }
    }
    if ($containerStartAttempted) {
      $labelAfter = Get-ResourceNames "label=$($RunSpec.Label)"
      $nameAfter = Get-ResourceNames "name=^/$([regex]::Escape($RunSpec.ContainerName))$"
      if (-not $labelAfter.Passed -or -not $nameAfter.Passed) {
        Add-CleanupCode 'HPG_RESOURCE_LEAK'
      } else {
        $inventory = @(@($labelAfter.Names) + @($nameAfter.Names) | Sort-Object -Unique)
        if ($inventory.Count -gt 0) { Add-CleanupCode 'HPG_RESOURCE_LEAK' }
      }
    }
  }

  if ([string]::IsNullOrWhiteSpace([string]$state.PrimaryCode) -and $cleanupCodes.Count -gt 0) {
    $state.PrimaryExit = 1
  }
  return [pscustomobject]@{
    ExitCode = $state.PrimaryExit
    PrimaryReasonCode = $state.PrimaryCode
    CleanupReasonCodes = @($cleanupCodes)
    MigrationsAppliedFirst = $firstCount
    MigrationsAppliedSecond = $secondCount
    ResourceInventory = @($inventory)
    RunId = $RunSpec.RunId
    ContainerName = $RunSpec.ContainerName
    DatabaseName = $RunSpec.DatabaseName
    HostPort = $hostPort
    FailureDiagnosticTokens = @($failureDiagnosticTokens)
    HeldConnectionConfirmed = $heldConnectionConfirmed
  }
}

Export-ModuleMember -Function Resolve-HarnessPostgresDockerApplication, New-HarnessPostgresRunSpec, Invoke-HarnessPostgresLifecycle
