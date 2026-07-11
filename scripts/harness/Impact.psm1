Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

function New-ImpactResult {
  param([bool]$Passed, [string]$Status, [string]$ErrorCode, [string]$OutcomePath, [string[]]$SharedSeams)
  [pscustomobject]@{ Passed=$Passed; Status=$Status; ErrorCode=$ErrorCode; OutcomePath=$OutcomePath; SharedSeams=$SharedSeams }
}

function Get-HarnessImpactCommand {
  param([Parameter(Mandatory)][string]$CommandId, [Parameter(Mandatory)][string]$RepositoryRoot)
  $pwsh = (Get-Command pwsh -CommandType Application -ErrorAction Stop).Source
  $fixture = Join-Path $RepositoryRoot 'scripts/tests/fixtures/harness/impact-command-probe.ps1'
  switch ($CommandId) {
    'impact-probe-one' { return [pscustomobject]@{ FilePath=$pwsh; Arguments=@('-NoProfile','-File',$fixture,'-Name','one'); Command='impact-probe-one'; Target='fake'; EvidenceClass='contract' } }
    'impact-probe-two' { return [pscustomobject]@{ FilePath=$pwsh; Arguments=@('-NoProfile','-File',$fixture,'-Name','two'); Command='impact-probe-two'; Target='fake'; EvidenceClass='contract' } }
    'governance-contracts' { return [pscustomobject]@{ FilePath=$pwsh; Arguments=@('-NoProfile','-File',(Join-Path $RepositoryRoot 'scripts/tests/governance-contracts.tests.ps1')); Command='governance-contracts'; Target='fake'; EvidenceClass='contract' } }
    'harness-aliases' { return [pscustomobject]@{ FilePath=$pwsh; Arguments=@('-NoProfile','-File',(Join-Path $RepositoryRoot 'scripts/tests/harness-aliases.tests.ps1')); Command='harness-aliases'; Target='fake'; EvidenceClass='contract' } }
    'context-compiler-tests' { return [pscustomobject]@{ FilePath=$pwsh; Arguments=@('-NoProfile','-File',(Join-Path $RepositoryRoot 'scripts/tests/context-compiler.tests.ps1')); Command='context-compiler-tests'; Target='fake'; EvidenceClass='contract' } }
    'orchestration-tests' { return [pscustomobject]@{ FilePath=$pwsh; Arguments=@('-NoProfile','-File',(Join-Path $RepositoryRoot 'scripts/tests/harness-orchestration.tests.ps1')); Command='orchestration-tests'; Target='fake'; EvidenceClass='contract' } }
    'harness-evals' { return [pscustomobject]@{ FilePath=$pwsh; Arguments=@('-NoProfile','-File',(Join-Path $RepositoryRoot 'scripts/tests/harness-eval.tests.ps1')); Command='harness-evals'; Target='fake'; EvidenceClass='contract' } }
    default { return $null }
  }
}

function Get-HarnessChangedPaths {
  param([string]$RepositoryRoot, [string]$BaseSha)
  $tracked = @(& git -C $RepositoryRoot diff --name-only $BaseSha 2>$null)
  if ($LASTEXITCODE -ne 0) { throw 'HIMPACT_BASE_UNAVAILABLE' }
  $untracked = @(& git -C $RepositoryRoot ls-files --others --exclude-standard 2>$null)
  if ($LASTEXITCODE -ne 0) { throw 'HIMPACT_CHANGED_PATHS_UNAVAILABLE' }
  @($tracked + $untracked | ForEach-Object { $_.Trim().Replace('\','/') } | Where-Object { $_ } | Sort-Object -Unique)
}

function Test-HarnessImpactPath {
  param([string]$Path, [string[]]$Allowed)
  @($Allowed | Where-Object { $Path -eq $_ -or $Path.StartsWith($_.TrimEnd('/','*') + '/', [StringComparison]::Ordinal) }).Count -gt 0
}

function Get-HarnessImpactSeams {
  param([string]$RepositoryRoot, [string[]]$ChangedPaths)
  $registry = Get-Content -Raw -LiteralPath (Join-Path $RepositoryRoot 'contracts/governance/shared-seams.json') | ConvertFrom-Json
  @($registry.seams | Where-Object { $seam = $_; @($ChangedPaths | Where-Object { $path = $_; @($seam.exclusive_paths | Where-Object { $path -eq $_ -or $path.StartsWith($_.TrimEnd('/') + '/', [StringComparison]::Ordinal) -or $_.StartsWith($path.TrimEnd('/') + '/', [StringComparison]::Ordinal) }).Count -gt 0 }).Count -gt 0 } | ForEach-Object id | Sort-Object -Unique)
}

function Invoke-HarnessImpactGate {
  [CmdletBinding()]
  param([Parameter(Mandatory)][string]$RepositoryRoot, [Parameter(Mandatory)][string]$ContextPath, [Parameter(Mandatory)][string]$RunId, [Parameter(Mandatory)][string]$RunDirectory)
  Import-Module (Join-Path $RepositoryRoot 'scripts/harness/Context.psm1') -Force -Global
  Import-Module (Join-Path $RepositoryRoot 'scripts/harness/Execution.psm1') -Force
  Import-Module (Join-Path $RepositoryRoot 'scripts/harness/Evidence.psm1') -Force
  Import-Module (Join-Path $RepositoryRoot 'scripts/harness/Environment.psm1') -Force
  $outcomePath = Join-Path $RunDirectory 'outcome.json'; $relativeOutcome = "scripts/.runs/$RunId/outcome.json"; New-Item -ItemType Directory -Path $RunDirectory -Force | Out-Null
  $head = (& git -C $RepositoryRoot rev-parse HEAD).Trim(); $branch = (& git -C $RepositoryRoot branch --show-current).Trim(); if (-not $branch) { $branch='detached' }
  $packResult = Test-HarnessContextPack -Path $ContextPath -RepositoryRoot $RepositoryRoot -RequireCurrentBase
  $baseSha = if ($packResult.Pack -and [string]$packResult.Pack.base_sha -match '^[0-9a-f]{40}$') { [string]$packResult.Pack.base_sha } else { '0000000000000000000000000000000000000000' }
  $risk = if ($packResult.Pack -and [string]$packResult.Pack.risk.level -in @('L0','L1','L2','L3')) { [string]$packResult.Pack.risk.level } else { 'L0' }
  $records = [Collections.Generic.List[object]]::new(); $changedPaths=@(); $seams=@(); $failure=''; $environment = New-HarnessChildEnvironment -RepositoryRoot $RepositoryRoot -LaneId 'unit'
  $environment['GIT_CONFIG_COUNT'] = '1'; $environment['GIT_CONFIG_KEY_0'] = 'safe.directory'; $environment['GIT_CONFIG_VALUE_0'] = [IO.Path]::GetFullPath($RepositoryRoot)
  if (-not $packResult.Passed) { $failure = [string]$packResult.ErrorCode } else {
    $changedPaths = Get-HarnessChangedPaths -RepositoryRoot $RepositoryRoot -BaseSha $baseSha
    $outside = @($changedPaths | Where-Object { -not (Test-HarnessImpactPath -Path $_ -Allowed @($packResult.Pack.paths.allowed)) })
    if ($outside.Count -gt 0) { $failure = 'HIMPACT_PATH_OUTSIDE_SCOPE' } else {
      $seams = Get-HarnessImpactSeams -RepositoryRoot $RepositoryRoot -ChangedPaths $changedPaths
      foreach ($selected in @($packResult.Pack.commands)) {
        $definition = Get-HarnessImpactCommand -CommandId ([string]$selected.command_id) -RepositoryRoot $RepositoryRoot
        if ($null -eq $definition -or $definition.Target -ne [string]$selected.target_label) { $failure='HIMPACT_COMMAND_UNREGISTERED'; break }
        $timer=[Diagnostics.Stopwatch]::StartNew()
        $process = Invoke-HarnessProcess -Request (New-HarnessProcessRequest -FilePath $definition.FilePath -ArgumentList $definition.Arguments -WorkingDirectory $RepositoryRoot -Environment $environment -TimeoutSeconds 120 -RedactionCandidates @($RepositoryRoot))
        $timer.Stop()
        $reason = if ($process.ExitCode -eq 0) { 'passed' } elseif ($process.TimedOut) { 'HEXEC_TIMEOUT' } else { [string]$process.ReasonCode }
        $records.Add([ordered]@{id=[string]$selected.id;command=[string]$definition.Command;stage='validation';target_label=$definition.Target;evidence_class=$definition.EvidenceClass;duration_ms=$timer.ElapsedMilliseconds;exit_code=$process.ExitCode;reason=$reason;artifact_paths=@($relativeOutcome)})
        if ($process.ExitCode -ne 0) { $failure='HIMPACT_COMMAND_FAILED'; break }
      }
    }
  }
  if ($records.Count -eq 0) { $records.Add([ordered]@{id='impact-preflight';command='impact-preflight';stage='preflight';target_label='fake';evidence_class='contract';duration_ms=0;exit_code=1;reason=if($failure){$failure}else{'HIMPACT_UNRECORDED'};artifact_paths=@($relativeOutcome)}) }
  $aggregate = if ($failure -match '^(CTX_|HIMPACT_(?:PATH|BASE|CHANGED|COMMAND_UNREGISTERED))') { 'blocked' } elseif ($failure) { 'failed' } else { 'passed' }
  $outcome = New-HarnessOutcome -RunId $RunId -BaseSha $baseSha -CommitSha $head -Branch $branch -Dirty (@(& git -C $RepositoryRoot status --porcelain).Count -gt 0) -Risk $risk -ChangedPaths $changedPaths -Tools @{ pwsh = 'registered' } -Commands @($records) -RemainingRisks $(if($failure){@($failure)}else{@()}) -AggregateClassification $aggregate
  Write-HarnessOutcome -Outcome $outcome -Path $outcomePath | Out-Null
  [pscustomobject]@{ Passed=($aggregate -eq 'passed'); Status=$aggregate; ErrorCode=$failure; OutcomePath=$relativeOutcome; SharedSeams=$seams }
}

Export-ModuleMember -Function Get-HarnessImpactCommand, Invoke-HarnessImpactGate
