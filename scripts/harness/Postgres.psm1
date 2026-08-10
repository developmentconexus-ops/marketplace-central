Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

Import-Module (Join-Path $PSScriptRoot 'Execution.psm1')

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

function Get-HarnessCanonicalMigrationCount {
  param([Parameter(Mandatory)][string]$RepositoryRoot)

  $migrationDirectory = Join-Path $RepositoryRoot 'apps/server_core/migrations'
  if (-not (Test-Path -LiteralPath $migrationDirectory -PathType Container)) {
    throw 'HPG_MIGRATION_INVENTORY_INVALID'
  }
  try {
    $migrations = @(Get-ChildItem -LiteralPath $migrationDirectory -File -ErrorAction Stop | Where-Object {
      $_.Name.EndsWith('.sql', [StringComparison]::Ordinal)
    })
  } catch {
    throw 'HPG_MIGRATION_INVENTORY_INVALID'
  }
  if ($migrations.Count -lt 1) { throw 'HPG_MIGRATION_INVENTORY_INVALID' }
  return [int]$migrations.Count
}

function Get-HarnessIntegrationTestPackages {
  [CmdletBinding()]
  param([Parameter(Mandatory)][string]$RepositoryRoot)

  $serverRoot = Join-Path $RepositoryRoot 'apps/server_core'
  $packages = [System.Collections.Generic.SortedSet[string]]::new([StringComparer]::Ordinal)
  [void]$packages.Add('./tests/integration')
  # The whole module, not just internal/modules: migrations/icms_matrix_test.go
  # carried the integration tag for a month while this walk could not see it --
  # a test in the tree that no lane ran. Dot-directories are excluded because
  # .gocache/.gomodcache live under the module root and recursing into them
  # turns a directory walk into a cache walk.
  $walkRoots = @(Get-ChildItem -LiteralPath $serverRoot -Directory -ErrorAction SilentlyContinue |
      Where-Object { $_.Name -notlike '.*' })
  foreach ($walkRoot in $walkRoots) {
    $testFiles = @(Get-ChildItem -LiteralPath $walkRoot.FullName -Recurse -File -Filter '*_test.go' -ErrorAction SilentlyContinue)
    foreach ($file in $testFiles) {
      $head = @(Get-Content -LiteralPath $file.FullName -TotalCount 5 -ErrorAction SilentlyContinue)
      if (@($head | Where-Object { $_ -match '^//go:build\b.*\bintegration\b' }).Count -gt 0) {
        $relative = [IO.Path]::GetRelativePath($serverRoot, $file.DirectoryName).Replace('\', '/')
        [void]$packages.Add("./$relative")
      }
    }
  }
  return @($packages)
}

function Get-HarnessPostgresSessionStatePath {
  param([Parameter(Mandatory)][string]$RepositoryRoot)
  return Join-Path $RepositoryRoot 'scripts/.runs/pg-session.json'
}

function Get-HarnessPostgresSessionContainerName {
  param([Parameter(Mandatory)][string]$RepositoryRoot)
  # Per-checkout name: hub and chip worktrees each own their session container,
  # so one checkout starting a session never removes another checkout's container.
  $canonical = [IO.Path]::GetFullPath($RepositoryRoot).TrimEnd('\', '/').ToLowerInvariant()
  $digest = [Security.Cryptography.SHA256]::HashData([Text.Encoding]::UTF8.GetBytes($canonical))
  $suffix = [Convert]::ToHexString($digest).Substring(0, 8).ToLowerInvariant()
  return "mpc-pg-session-$suffix"
}

function Get-HarnessPostgresSession {
  [CmdletBinding()]
  param(
    [Parameter(Mandatory)][string]$RepositoryRoot,
    [Parameter(Mandatory)][string]$DockerFilePath,
    [string[]]$DockerArgumentPrefix = @()
  )

  $statePath = Get-HarnessPostgresSessionStatePath -RepositoryRoot $RepositoryRoot
  if (-not (Test-Path -LiteralPath $statePath -PathType Leaf)) { return $null }
  try { $state = Get-Content -Raw -LiteralPath $statePath | ConvertFrom-Json } catch { return $null }
  if ($null -eq $state -or
      [string]::IsNullOrWhiteSpace([string]$state.ContainerName) -or
      $state.ContainerName -cne (Get-HarnessPostgresSessionContainerName -RepositoryRoot $RepositoryRoot) -or
      [string]::IsNullOrWhiteSpace([string]$state.Password) -or
      [int]$state.HostPort -lt 1 -or [int]$state.HostPort -gt 65535) {
    return $null
  }
  $environment = [System.Collections.Generic.Dictionary[string,string]]::new([StringComparer]::OrdinalIgnoreCase)
  $request = New-HarnessProcessRequest -FilePath $DockerFilePath -ArgumentList (@($DockerArgumentPrefix) + @('inspect', '--format', '{{.State.Running}}', [string]$state.ContainerName)) -WorkingDirectory $RepositoryRoot -Environment $environment -TimeoutSeconds 60 -RedactionCandidates @([string]$state.Password)
  $inspect = Invoke-HarnessProcess -Request $request
  if ($inspect.ExitCode -ne 0 -or $inspect.Stdout.Trim() -cne 'true') { return $null }
  return [pscustomobject]@{
    ContainerName = [string]$state.ContainerName
    HostPort = [int]$state.HostPort
    Password = [string]$state.Password
  }
}

function Start-HarnessPostgresSession {
  [CmdletBinding()]
  param(
    [Parameter(Mandatory)][string]$RepositoryRoot,
    [Parameter(Mandatory)][string]$DockerFilePath,
    [string[]]$DockerArgumentPrefix = @(),
    [ValidateRange(1, 600)][int]$ReadyMaxAttempts = 60,
    [ValidateRange(0, 60000)][int]$ReadyRetryDelayMilliseconds = 1000
  )

  $existing = Get-HarnessPostgresSession -RepositoryRoot $RepositoryRoot -DockerFilePath $DockerFilePath -DockerArgumentPrefix $DockerArgumentPrefix
  if ($null -ne $existing) { return $existing }

  $containerName = Get-HarnessPostgresSessionContainerName -RepositoryRoot $RepositoryRoot
  $passwordBytes = [byte[]]::new(24)
  [Security.Cryptography.RandomNumberGenerator]::Fill($passwordBytes)
  $password = [Convert]::ToHexString($passwordBytes).ToLowerInvariant()
  $environment = [System.Collections.Generic.Dictionary[string,string]]::new([StringComparer]::OrdinalIgnoreCase)
  $environment['POSTGRES_PASSWORD'] = $password
  function Invoke-SessionDocker([string[]]$Arguments) {
    $request = New-HarnessProcessRequest -FilePath $DockerFilePath -ArgumentList (@($DockerArgumentPrefix) + @($Arguments)) -WorkingDirectory $RepositoryRoot -Environment $environment -TimeoutSeconds 120 -RedactionCandidates @($password)
    return Invoke-HarnessProcess -Request $request
  }

  # Stale state or dead container: clear both before starting fresh.
  [void](Invoke-SessionDocker @('rm', '--force', $containerName))
  $statePath = Get-HarnessPostgresSessionStatePath -RepositoryRoot $RepositoryRoot
  Remove-Item -LiteralPath $statePath -Force -ErrorAction SilentlyContinue

  $start = Invoke-SessionDocker @(
    'run', '--detach', '--rm', '--pull=never',
    '--name', $containerName,
    '--label', 'marketplace-central.harness.session=true',
    '--publish', '127.0.0.1::5432',
    '--tmpfs', '/var/lib/postgresql/data:rw,noexec,nosuid,size=512m',
    '--env', 'POSTGRES_PASSWORD',
    '--env', 'POSTGRES_DB=postgres',
    'postgres:16-bookworm'
  )
  if ($start.ExitCode -ne 0) { throw 'HPG_CONTAINER_START_FAILED' }

  $ready = $null
  for ($attempt = 1; $attempt -le $ReadyMaxAttempts; $attempt++) {
    $ready = Invoke-SessionDocker @('exec', $containerName, 'pg_isready', '--username', 'postgres', '--dbname', 'postgres', '--timeout', '5')
    if ($ready.ExitCode -eq 0) { break }
    if ($attempt -lt $ReadyMaxAttempts -and $ReadyRetryDelayMilliseconds -gt 0) { Start-Sleep -Milliseconds $ReadyRetryDelayMilliseconds }
  }
  if ($null -eq $ready -or $ready.ExitCode -ne 0) {
    [void](Invoke-SessionDocker @('rm', '--force', $containerName))
    throw 'HPG_READY_TIMEOUT'
  }

  $port = Invoke-SessionDocker @('port', $containerName, '5432/tcp')
  if ($port.ExitCode -ne 0 -or $port.Stdout -notmatch '(?m)^127\.0\.0\.1:(\d+)\s*$') {
    [void](Invoke-SessionDocker @('rm', '--force', $containerName))
    throw 'HPG_PORT_UNAVAILABLE'
  }
  $hostPort = [int]$Matches[1]

  $stateDirectory = Split-Path -Parent $statePath
  New-Item -ItemType Directory -Path $stateDirectory -Force | Out-Null
  @{ ContainerName = $containerName; HostPort = $hostPort; Password = $password } | ConvertTo-Json | Set-Content -LiteralPath $statePath -Encoding utf8
  return [pscustomobject]@{ ContainerName = $containerName; HostPort = $hostPort; Password = $password }
}

function Stop-HarnessPostgresSession {
  [CmdletBinding()]
  param(
    [Parameter(Mandatory)][string]$RepositoryRoot,
    [Parameter(Mandatory)][string]$DockerFilePath,
    [string[]]$DockerArgumentPrefix = @()
  )

  $environment = [System.Collections.Generic.Dictionary[string,string]]::new([StringComparer]::OrdinalIgnoreCase)
  $request = New-HarnessProcessRequest -FilePath $DockerFilePath -ArgumentList (@($DockerArgumentPrefix) + @('rm', '--force', (Get-HarnessPostgresSessionContainerName -RepositoryRoot $RepositoryRoot))) -WorkingDirectory $RepositoryRoot -Environment $environment -TimeoutSeconds 120 -RedactionCandidates @()
  [void](Invoke-HarnessProcess -Request $request)
  $statePath = Get-HarnessPostgresSessionStatePath -RepositoryRoot $RepositoryRoot
  Remove-Item -LiteralPath $statePath -Force -ErrorAction SilentlyContinue
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
  $expectedMigrationCount = Get-HarnessCanonicalMigrationCount -RepositoryRoot $root
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
    ExpectedMigrationCount = $expectedMigrationCount
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
  # Everything above names WHAT failed and nothing about WHERE. A run of this lane
  # failed on 2026-08-10 with `test=TestListingsReadContractEndToEnd` plus two
  # sqlstates and passed on re-run with the same tree; the tokens could not say
  # which assertion, which subtest, or which constraint, so the only available
  # move was to re-run it -- which is how a red lane becomes a habit of ignoring
  # red. The three shapes below are source and schema identifiers, never row data:
  # a Go subtest name (spaces are already underscores by the time go test prints
  # it), a _test.go file and line, and a Postgres constraint name.
  foreach ($match in [regex]::Matches([string]$Text, '(?m)^\s+--- FAIL: (?<test>Test[A-Za-z0-9_]{1,160}(?:/[A-Za-z0-9_]{1,120}){1,6})(?:\s|$)')) { [void]$tokens.Add("subtest=$($match.Groups['test'].Value)") }
  #
  # `at=` and `constraint=` are attributed, not swept. A test's own log lines sit
  # between its `=== RUN` and its verdict line, so a file:line or a constraint
  # name only belongs to a failure when the verdict that closes it is FAIL. Swept
  # instead, a real run of this lane emitted 22 `at=` tokens and 3 constraints
  # from tests that PASSED -- must-fail tests log the constraint they provoked --
  # and the sort below would have spent the cap on them before reaching the one
  # failing line. Interleaved parallel output can still misattribute a line; the
  # tokens are source identifiers either way, never row data.
  $pending = [System.Collections.Generic.List[string]]::new()
  foreach ($line in ([string]$Text -split '\r?\n')) {
    if ($line -match '^\s*--- (?<verdict>PASS|FAIL|SKIP):') {
      if ($Matches['verdict'] -ceq 'FAIL') { foreach ($token in $pending) { [void]$tokens.Add($token) } }
      $pending.Clear()
      continue
    }
    if ($line -match '^\s*(?<file>[A-Za-z0-9_]{1,80}_test\.go):(?<line>\d{1,6}):') { $pending.Add("at=$($Matches['file']):$($Matches['line'])") }
    if ($line -match 'constraint "(?<constraint>[a-z0-9_]{1,63})"') { $pending.Add("constraint=$($Matches['constraint'])") }
  }
  $safe = @($tokens | Sort-Object | Select-Object -First 32)
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
    [ValidateRange(1, 60)][int]$CreateMaxAttempts = 15,
    [ValidateRange(0, 60000)][int]$CreateRetryDelayMilliseconds = 1000,
    [string[]]$TestArguments = @(),
    [switch]$HoldConnectionDuringCleanupTest,
    [object]$Session = $null
  )

  if (-not $RunSpec.PSObject.Properties['ExpectedMigrationCount'] -or $RunSpec.ExpectedMigrationCount -isnot [int] -or $RunSpec.ExpectedMigrationCount -lt 1) {
    throw 'HPG_MIGRATION_INVENTORY_INVALID'
  }
  $expectedMigrationCount = [int]$RunSpec.ExpectedMigrationCount
  if (@($TestArguments).Count -eq 0) {
    $TestArguments = @('test', '-tags=integration') + @(Get-HarnessIntegrationTestPackages -RepositoryRoot $RunSpec.RepositoryRoot) + @('-count=1')
  }
  # -v is not decoration, and it is applied to CALLER-SUPPLIED arguments too.
  # Without it `go test` prints one `ok <pkg>` line per package and nothing else,
  # and this function used to discard even that: on success `$tests.Stdout` was
  # never read, never returned, never written to the run directory. A lane that
  # compiled zero test packages and one that ran every integration test produced
  # byte-identical output -- `status=passed`. The per-test lines are what make
  # the counts below a measurement, so a caller that passed its own arguments
  # without -v would otherwise trip the vacuity guard on a run that did work.
  if (-not (@($TestArguments) -contains '-v')) { $TestArguments = @($TestArguments) + @('-v') }
  $dockerEnvironment = Copy-HarnessEnvironment $BaseEnvironment
  $dockerEnvironment['POSTGRES_PASSWORD'] = [string]$RunSpec.Password
  $goEnvironment = Copy-HarnessEnvironment $BaseEnvironment
  # The hermetic lane compiles the same file set on every platform. Without
  # this, a linux runner with gcc has cgo on by default and `integration && cgo`
  # files -- the live Oracle tests -- enter a lane that has no Oracle, while a
  # bare Windows host silently excludes them. Live-cgo tests belong to the
  # live-oracle lane alone.
  $goEnvironment['CGO_ENABLED'] = '0'
  $cleanupCodes = [System.Collections.Generic.List[string]]::new()
  $inventory = @()
  $state = @{ PrimaryCode = ''; PrimaryExit = 0 }
  $databaseCreated = $false
  $containerStartAttempted = $false
  $containerOwned = $false
  $execContainer = if ($null -ne $Session) { [string]$Session.ContainerName } else { [string]$RunSpec.ContainerName }
  $firstCount = -1
  $secondCount = -1
  $hostPort = 0
  $targetURL = ''
  $failureDiagnosticTokens = @()
  # Counted here rather than inferred from what the child logged: a readiness
  # attempt whose subprocess is killed at its own budget writes nothing, so a
  # caller counting child rows reads 0 attempts on a run that made one. That
  # miscount is what made the deadline assertion in postgres-lifecycle.tests.ps1
  # fail under load on a tree that behaved correctly.
  $readyAttempts = 0
  $heldConnectionConfirmed = $false
  $testsRun = -1
  $testsPassed = -1
  $testsSkipped = -1
  $testsFailed = -1

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

    if ($null -ne $Session) {
      # Session mode: reuse the long-lived container started by Start-HarnessPostgresSession.
      # Verify it is still running, then skip straight to per-run database creation.
      $sessionAlive = Invoke-Docker @('inspect', '--format', '{{.State.Running}}', $execContainer)
      if ($sessionAlive.ExitCode -ne 0 -or $sessionAlive.Stdout.Trim() -cne 'true') { Set-Primary 'HPG_CONTAINER_START_FAILED' 1; break }
      $hostPort = [int]$Session.HostPort
      if ($hostPort -lt 1 -or $hostPort -gt 65535) { Set-Primary 'HPG_PORT_UNAVAILABLE' 1; break }
    } else {
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
      $readyAttempts++
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
    }
    $escapedPassword = [uri]::EscapeDataString([string]$RunSpec.Password)
    $targetURL = "postgresql://postgres:$escapedPassword@127.0.0.1:$hostPort/$($RunSpec.DatabaseName)?sslmode=disable"
    $goEnvironment['MPC_TEST_DATABASE_URL'] = $targetURL

    # pg_isready can pass during the image's first-boot init phase (server restarts once
    # after init), so the first CREATE DATABASE may hit a restarting server. Retry until
    # the database provably exists instead of trusting readiness.
    $create = $null
    for ($createAttempt = 1; $createAttempt -le $CreateMaxAttempts; $createAttempt++) {
      $exists = Invoke-Docker @('exec', $execContainer, 'psql', '--username', 'postgres', '--dbname', 'postgres', '--quiet', '--tuples-only', '--no-align', '--set', 'ON_ERROR_STOP=1', '--command', "SELECT 1 FROM pg_database WHERE datname = '$($RunSpec.DatabaseName)'")
      if ($exists.ExitCode -eq 0 -and $exists.Stdout.Trim() -eq '1') { $create = $exists; break }
      $create = Invoke-Docker @('exec', $execContainer, 'psql', '--username', 'postgres', '--dbname', 'postgres', '--set', 'ON_ERROR_STOP=1', '--command', "CREATE DATABASE $($RunSpec.DatabaseName)")
      if ($create.ExitCode -eq 0) { break }
      if ($createAttempt -lt $CreateMaxAttempts -and $CreateRetryDelayMilliseconds -gt 0) { Start-Sleep -Milliseconds $CreateRetryDelayMilliseconds }
    }
    if ($null -eq $create -or $create.ExitCode -ne 0) { Set-Primary 'HPG_DATABASE_CREATE_FAILED' $(if ($null -eq $create) { 1 } else { $create.ExitCode }); break }
    $databaseCreated = $true

    $migrationArgs = @('run', './cmd/testdb', 'migrate')
    $migration1 = Invoke-Go $migrationArgs
    if ($migration1.ExitCode -ne 0) { Set-Primary 'HPG_MIGRATION_FAILED' $migration1.ExitCode; break }
    $firstCount = Get-HarnessMigrationCount $migration1.Stdout
    $migration2 = Invoke-Go $migrationArgs
    if ($migration2.ExitCode -ne 0) { Set-Primary 'HPG_MIGRATION_FAILED' $migration2.ExitCode; break }
    $secondCount = Get-HarnessMigrationCount $migration2.Stdout
    if ($firstCount -ne $expectedMigrationCount -or $secondCount -ne 0) { Set-Primary 'HPG_MIGRATION_NOT_IDEMPOTENT' 1; break }

    if ($HoldConnectionDuringCleanupTest) {
      $heldConnection = Invoke-Docker @('exec', '--detach', $execContainer, 'psql', '--username', 'postgres', '--dbname', $RunSpec.DatabaseName, '--command', 'SELECT pg_sleep(300)')
      if ($heldConnection.ExitCode -ne 0) { Set-Primary 'HPG_TEST_FAILED' $heldConnection.ExitCode; break }
      for ($heldAttempt = 1; $heldAttempt -le 20; $heldAttempt++) {
        $heldCheck = Invoke-Docker @('exec', $execContainer, 'psql', '--username', 'postgres', '--dbname', 'postgres', '--quiet', '--tuples-only', '--no-align', '--set', 'ON_ERROR_STOP=1', '--command', "SELECT count(*) FROM pg_stat_activity WHERE datname = '$($RunSpec.DatabaseName)' AND state = 'active' AND query LIKE 'SELECT pg_sleep(300)%'")
        if ($heldCheck.ExitCode -eq 0 -and $heldCheck.Stdout.Trim() -match '^[1-9][0-9]*$') { $heldConnectionConfirmed = $true; break }
        Start-Sleep -Milliseconds 100
      }
      if (-not $heldConnectionConfirmed) { Set-Primary 'HPG_TEST_FAILED' 1; break }
    }
    $tests = Invoke-Go @($TestArguments)
    $testOutput = @($tests.Stdout, $tests.Stderr) -join "`n"
    # Subtest lines are indented, top-level ones are not, so both anchors allow
    # leading whitespace. Subtests are counted as units of their own: a run that
    # executed one parent and forty children did forty-one things.
    $testsRun = @([regex]::Matches($testOutput, '(?m)^\s*=== RUN\s')).Count
    $testsPassed = @([regex]::Matches($testOutput, '(?m)^\s*--- PASS:')).Count
    $testsSkipped = @([regex]::Matches($testOutput, '(?m)^\s*--- SKIP:')).Count
    $testsFailed = @([regex]::Matches($testOutput, '(?m)^\s*--- FAIL:')).Count
    if ($tests.ExitCode -ne 0) {
      $failureDiagnosticTokens = @(Get-HarnessPostgresFailureTokens $testOutput)
      Set-Primary 'HPG_TEST_FAILED' $tests.ExitCode
      break
    }
    # `go test` exits 0 over a package set that contains no tests, over a set
    # that is empty, and over a set where every test called t.Skip. All three are
    # the same green as a full run. This is the only place that can tell them
    # apart, because it is the only place that holds the output.
    #
    # The skipped case is the one the profile already names: HARNESS-PROFILE.md
    # records `RUN 27 / PASS 1 / SKIP 26` from a lane whose environment was never
    # loaded, so every DB test hit SkipWithoutTarget and the package printed `ok`
    # with exit 0 -- with the slice's whole reason to exist among the skips. That
    # is why the guard is `passed -eq 0` and not only `run -eq 0`: a suite that
    # ran and asserted nothing has proved exactly as much as one that did not run.
    if ($testsRun -eq 0 -or $testsPassed -eq 0) {
      $failureDiagnosticTokens = @("tests_run=$testsRun", "tests_passed=$testsPassed", "tests_skipped=$testsSkipped")
      Set-Primary 'HPG_TEST_VACUOUS' 1
      break
    }
    } while ($false)
  } finally {
    if ($databaseCreated -and ($containerOwned -or $null -ne $Session)) {
      $drop = Invoke-Docker @('exec', $execContainer, 'psql', '--username', 'postgres', '--dbname', 'postgres', '--set', 'ON_ERROR_STOP=1', '--command', "DROP DATABASE $($RunSpec.DatabaseName) WITH (FORCE)")
      if ($drop.ExitCode -ne 0) { Add-CleanupCode 'HPG_DATABASE_DROP_FAILED' }
      $dropVerify = Invoke-Docker @('exec', $execContainer, 'psql', '--username', 'postgres', '--dbname', 'postgres', '--tuples-only', '--no-align', '--set', 'ON_ERROR_STOP=1', '--command', "SELECT datname FROM pg_database WHERE datname = '$($RunSpec.DatabaseName)'")
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
    ReadyAttempts = $readyAttempts
    HeldConnectionConfirmed = $heldConnectionConfirmed
    TestsRun = $testsRun
    TestsPassed = $testsPassed
    TestsSkipped = $testsSkipped
    TestsFailed = $testsFailed
  }
}

Export-ModuleMember -Function Resolve-HarnessPostgresDockerApplication, New-HarnessPostgresRunSpec, Invoke-HarnessPostgresLifecycle, Get-HarnessIntegrationTestPackages, Get-HarnessPostgresSession, Start-HarnessPostgresSession, Stop-HarnessPostgresSession
