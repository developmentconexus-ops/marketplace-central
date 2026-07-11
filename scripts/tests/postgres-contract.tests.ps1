$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$module = Join-Path $repoRoot 'scripts/harness/Postgres.psm1'
$harness = Join-Path $repoRoot 'scripts/harness.ps1'
$probe = Join-Path $repoRoot 'scripts/tests/fixtures/harness/postgres-docker-probe.mjs'

function Assert-True([bool]$Condition, [string]$Message) {
  if (-not $Condition) { throw $Message }
}

if (-not (Test-Path -LiteralPath $module -PathType Leaf)) {
  throw 'RED: Postgres.psm1 missing; ephemeral PostgreSQL contract is not implemented'
}
Import-Module $module -Force
$moduleSource = Get-Content -Raw -LiteralPath $module
Assert-True ($moduleSource -match '\[AllowEmptyCollection\(\)\]\[string\[\]\]\$ArgumentPrefix') 'real executable path cannot bind an empty argument prefix'

$node = [IO.Path]::GetFullPath((Get-Command node -CommandType Application -ErrorAction Stop).Source)
$runId = '0123456789abcdef0123456789abcdef'
$password = 'fixture-' + [guid]::NewGuid().ToString('N')
$log = Join-Path ([IO.Path]::GetTempPath()) ('postgres-contract-' + [guid]::NewGuid().ToString('N') + '.jsonl')
$environment = [System.Collections.Generic.Dictionary[string,string]]::new([StringComparer]::OrdinalIgnoreCase)
foreach ($key in @('SystemRoot', 'WINDIR', 'ComSpec', 'PATH', 'PATHEXT', 'TEMP', 'TMP')) {
  $value = [Environment]::GetEnvironmentVariable($key, 'Process')
  if (-not [string]::IsNullOrWhiteSpace($value)) { $environment[$key] = $value }
}
$environment['HARNESS_POSTGRES_PROBE_LOG'] = $log

try {
  $spec = New-HarnessPostgresRunSpec -RepositoryRoot $repoRoot -RunId $runId -Password $password -DockerFilePath $node -DockerArgumentPrefix @($probe, 'docker')
  Assert-True ($spec.RunId -ceq $runId) 'run ID changed'
  Assert-True ($spec.DatabaseName -ceq "mpc_test_$runId") 'database name is not tied to run ID'
  Assert-True ($spec.ContainerName -ceq "mpc-pg-$runId") 'container name is not tied to run ID'
  Assert-True ($spec.Label -ceq "marketplace-central.harness.run=$runId") 'run label is not tied to run ID'
  Assert-True ([IO.Path]::IsPathFullyQualified($spec.RepositoryRoot)) 'repository root is not absolute'

  $missingDocker = ''
  try { Resolve-HarnessPostgresDockerApplication -Name ('missing-docker-' + [guid]::NewGuid().ToString('N')) | Out-Null } catch { $missingDocker = $_.Exception.Message }
  Assert-True ($missingDocker -match 'HPG_DOCKER_MISSING') 'missing Docker executable lacks stable reason'

  foreach ($invalid in @('', 'ABCDEF0123456789abcdef0123456789', '0123456789abcdef', '../0123456789abcdef0123456789abcdef')) {
    $before = Test-Path -LiteralPath $log
    $message = ''
    try {
      New-HarnessPostgresRunSpec -RepositoryRoot $repoRoot -RunId $invalid -Password $password -DockerFilePath $node -DockerArgumentPrefix @($probe, 'docker') | Out-Null
    } catch { $message = $_.Exception.Message }
    Assert-True ($message -match 'HPG_RUN_ID_INVALID') "invalid run ID did not fail closed: $invalid"
    Assert-True ((Test-Path -LiteralPath $log) -eq $before) 'invalid run ID caused a process side effect'
  }

  $runsBefore = @(Get-ChildItem (Join-Path $repoRoot 'scripts/.runs') -Directory -ErrorAction SilentlyContinue).Name
  $external = & pwsh -NoProfile -ExecutionPolicy Bypass -File $harness -Command integration -DatabaseUrl 'postgresql://fixture@127.0.0.1:49152/mpc_test_0123456789abcdef0123456789abcdef' 2>&1
  $externalExit = $LASTEXITCODE
  $runsAfter = @(Get-ChildItem (Join-Path $repoRoot 'scripts/.runs') -Directory -ErrorAction SilentlyContinue).Name
  Assert-True ($externalExit -ne 0) 'caller database URL was accepted'
  Assert-True (($external -join "`n") -match 'HPG_EXTERNAL_TARGET_FORBIDDEN') 'caller URL lacks stable rejection reason'
  Assert-True (($runsBefore -join "`n") -ceq ($runsAfter -join "`n")) 'caller URL rejection created a run directory'
  Assert-True (-not (Test-Path -LiteralPath $log)) 'caller URL rejection invoked Docker'

  $result = Invoke-HarnessPostgresLifecycle -RunSpec $spec -BaseEnvironment $environment -GoFilePath $node -GoArgumentPrefix @($probe, 'go') -TimeoutSeconds 10
  Assert-True ($result.ExitCode -eq 0) 'happy fake lifecycle failed'
  $calls = @(Get-Content -LiteralPath $log | ForEach-Object { $_ | ConvertFrom-Json })
  $start = @($calls | Where-Object operation -eq 'start')[0]
  $startArgs = @($start.args)
  Assert-True ($startArgs -contains '--pull=never') 'Docker run may pull an image'
  Assert-True (($startArgs -join ' ') -match '--tmpfs /var/lib/postgresql/data:') 'Docker run lacks bounded tmpfs'
  Assert-True (($startArgs -join ' ') -match '127\.0\.0\.1::5432') 'Docker run lacks dynamic loopback publication'
  Assert-True (-not ($startArgs -contains '--volume') -and -not ($startArgs -contains '-v')) 'Docker run uses a volume'
  Assert-True ($start.hasPostgresPassword -and -not $start.secretInArgs) 'password was not environment-only'
  $serialized = $calls | ConvertTo-Json -Depth 8
  Assert-True ($serialized -notmatch [regex]::Escape($password)) 'password leaked into probe transcript'
  Assert-True (@($calls | Where-Object { -not [IO.Path]::IsPathFullyQualified([string]$_.cwd) }).Count -eq 0) 'child CWD was not explicit'
} finally {
  Remove-Item -LiteralPath $log, "$log.state.json" -Force -ErrorAction SilentlyContinue
}

Write-Output 'PASS postgres contract tests target=fake'
