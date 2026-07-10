$ErrorActionPreference = 'Stop'

$verifier = Join-Path $PSScriptRoot 'verify-baseline.ps1'
$fixtureRoot = Join-Path $env:TEMP ("verify-baseline-" + [guid]::NewGuid().ToString('N'))
$repositoryRoot = Join-Path $fixtureRoot 'repository'
$fixtureOwnershipMarker = Join-Path $fixtureRoot '.verify-baseline-fixture-owner'
$fixtureTemplateDirectory = Join-Path $fixtureRoot 'git-template'
$fixtureHooksDirectory = Join-Path $fixtureRoot 'hooks'
$fixtureRootCreated = $false
$fixtureOwnershipMarkerCreated = $false

function Write-Fixture {
    param(
        [string]$Name,
        [string[]]$InventoryRows,
        [string[]]$LedgerRows
    )

    $inventory = Join-Path $repositoryRoot "$Name-inventory.tsv"
    $ledger = Join-Path $repositoryRoot "$Name-ledger.tsv"
    @("status`tpath") + $InventoryRows | Set-Content -Path $inventory -Encoding utf8
    @("path`tstatus`tdisposition`tcohort`tcommit_sha`towner`treason`tblocker`tnext_validation") + $LedgerRows |
        Set-Content -Path $ledger -Encoding utf8
    return @{ Inventory = $inventory; Ledger = $ledger }
}

function Invoke-Git {
    param([string[]]$Arguments)

    & git -C $repositoryRoot @Arguments *> $null
    if ($LASTEXITCODE -ne 0) {
        throw "git $($Arguments -join ' ') failed with exit $LASTEXITCODE."
    }
}

function Assert-FixtureGitIsolation {
    $hooksPath = (& git -C $repositoryRoot config --local --get core.hooksPath)
    $expectedHooksPath = Join-Path $fixtureRoot 'hooks'
    if (
        [string]::IsNullOrWhiteSpace($hooksPath) -or
        [IO.Path]::GetFullPath($hooksPath) -ne [IO.Path]::GetFullPath($expectedHooksPath)
    ) {
        throw 'temporary fixture repository must set a local fixture-owned core.hooksPath.'
    }

    if (-not (Test-Path -LiteralPath $expectedHooksPath -PathType Container) -or @(Get-ChildItem -LiteralPath $expectedHooksPath).Count -ne 0) {
        throw 'temporary fixture repository hooks directory must exist and be empty.'
    }

    $commitSigning = (& git -C $repositoryRoot config --local --get commit.gpgSign)
    if ($commitSigning -ne 'false') {
        throw 'temporary fixture repository must disable local commit signing.'
    }
}

function Assert-VerifierExit {
    param(
        [string]$CaseName,
        [hashtable]$Fixture,
        [int]$ExpectedExit,
        [switch]$AllowRetainedState
    )

    $previousPreference = $ErrorActionPreference
    $ErrorActionPreference = 'Continue'
    $verifierArguments = @(
        '-NoProfile',
        '-ExecutionPolicy', 'Bypass',
        '-File', $verifier,
        '-Inventory', $Fixture.Inventory,
        '-Ledger', $Fixture.Ledger
    )
    if ($AllowRetainedState) {
        $verifierArguments += '-AllowRetainedState'
    }
    Push-Location $repositoryRoot
    try {
        & powershell @verifierArguments *> $null
        $actualExit = $LASTEXITCODE
    }
    finally {
        Pop-Location
        $ErrorActionPreference = $previousPreference
    }
    if ($actualExit -ne $ExpectedExit) {
        throw "$CaseName expected verifier exit $ExpectedExit but received $actualExit."
    }
}

try {
    New-Item -ItemType Directory -Path $fixtureRoot | Out-Null
    $fixtureRootCreated = $true
    New-Item -ItemType File -Path $fixtureOwnershipMarker -Value 'test-verify-baseline fixture owner' | Out-Null
    $fixtureOwnershipMarkerCreated = Test-Path -LiteralPath $fixtureOwnershipMarker -PathType Leaf
    if (-not $fixtureOwnershipMarkerCreated) {
        throw 'temporary fixture ownership marker was not created.'
    }

    New-Item -ItemType Directory -Path $repositoryRoot | Out-Null
    New-Item -ItemType Directory -Path $fixtureTemplateDirectory | Out-Null
    New-Item -ItemType Directory -Path $fixtureHooksDirectory | Out-Null
    Invoke-Git -Arguments @('init', "--template=$fixtureTemplateDirectory")
    Invoke-Git -Arguments @('config', 'user.email', 'verify-baseline@example.test')
    Invoke-Git -Arguments @('config', 'user.name', 'Verify Baseline Test')
    Invoke-Git -Arguments @('config', 'core.hooksPath', $fixtureHooksDirectory)
    Invoke-Git -Arguments @('config', 'commit.gpgSign', 'false')
    Assert-FixtureGitIsolation

    $missing = Write-Fixture -Name 'missing' -InventoryRows @(" M`talpha.txt", "??`tbeta.txt") -LedgerRows @(
        "alpha.txt`t M`tcommitted`tM-03`tdeadbeef`towner`tverified`tnone`tgit status"
    )
    $duplicate = Write-Fixture -Name 'duplicate' -InventoryRows @(" M`talpha.txt") -LedgerRows @(
        "alpha.txt`t M`tcommitted`tM-03`tdeadbeef`towner`tverified`tnone`tgit status",
        "alpha.txt`t M`tcommitted`tM-03`tdeadbeef`towner`tverified`tnone`tgit status"
    )
    $invalid = Write-Fixture -Name 'invalid' -InventoryRows @(" M`talpha.txt") -LedgerRows @(
        "alpha.txt`t M`tunknown-disposition`tM-03`tdeadbeef`towner`tverified`tnone`tgit status"
    )
    $complete = Write-Fixture -Name 'complete' -InventoryRows @(" M`talpha.txt") -LedgerRows @(
        "alpha.txt`t M`tcommitted`tM-03`tdeadbeef`towner`tverified`tnone`tgit status"
    )
    $retained = Write-Fixture -Name 'retained' -InventoryRows @(" M`talpha.txt") -LedgerRows @(
        "alpha.txt`t M`tretained-owner-needed`tM-08`t`towner`tawaiting ownership`tscoped validation needed`tgit status"
    )
    Set-Content -Path (Join-Path $repositoryRoot 'alpha.txt') -Value 'fixture tracked path' -Encoding utf8
    Invoke-Git -Arguments @('add', '--all')
    Invoke-Git -Arguments @('commit', '-m', 'test baseline fixtures')

    $fixtureCommitSha = (& git -C $repositoryRoot rev-parse HEAD).Trim()
    if ([string]::IsNullOrWhiteSpace($fixtureCommitSha)) {
        throw 'fixture commit SHA was not resolved.'
    }

    $complete = Write-Fixture -Name 'complete' -InventoryRows @(" M`talpha.txt") -LedgerRows @(
        "alpha.txt`t M`tcommitted`tM-03`t$fixtureCommitSha`towner`tverified`tnone`tgit status"
    )
    $invalidSha = Write-Fixture -Name 'invalid-sha' -InventoryRows @(" M`talpha.txt") -LedgerRows @(
        "alpha.txt`t M`tcommitted`tM-03`tdeadbeef`towner`tverified`tnone`tgit status"
    )
    $validMissingPath = Write-Fixture -Name 'valid-missing-path' -InventoryRows @("??`tmissing-target.txt") -LedgerRows @(
        "missing-target.txt`t??`tcommitted`tM-03`t$fixtureCommitSha`towner`tverified`tnone`tgit status"
    )
    $validIncludedPath = Write-Fixture -Name 'valid-included-path' -InventoryRows @(" M`talpha.txt") -LedgerRows @(
        "alpha.txt`t M`tcommitted`tM-03`t$fixtureCommitSha`towner`tverified`tnone`tgit status"
    )
    Invoke-Git -Arguments @('add', '--all')
    Invoke-Git -Arguments @('commit', '-m', 'add commit provenance fixtures')

    Assert-VerifierExit -CaseName 'missing path' -Fixture $missing -ExpectedExit 1
    Assert-VerifierExit -CaseName 'duplicate ledger path' -Fixture $duplicate -ExpectedExit 1
    Assert-VerifierExit -CaseName 'invalid disposition' -Fixture $invalid -ExpectedExit 1
    Assert-VerifierExit -CaseName 'retained state requires explicit opt-in' -Fixture $retained -ExpectedExit 1
    Assert-VerifierExit -CaseName 'clean committed baseline' -Fixture $complete -ExpectedExit 0
    Assert-VerifierExit -CaseName 'invalid committed SHA' -Fixture $invalidSha -ExpectedExit 1
    Assert-VerifierExit -CaseName 'valid SHA missing original path' -Fixture $validMissingPath -ExpectedExit 1
    Assert-VerifierExit -CaseName 'valid SHA including original path' -Fixture $validIncludedPath -ExpectedExit 0

    $dirtyPath = Join-Path $repositoryRoot 'controlled-dirty.txt'
    Set-Content -Path $dirtyPath -Value 'controlled dirty state' -Encoding utf8
    Assert-VerifierExit -CaseName 'dirty committed baseline' -Fixture $complete -ExpectedExit 1
    Assert-VerifierExit -CaseName 'dirty retained diagnostic baseline' -Fixture $retained -ExpectedExit 1 -AllowRetainedState

    Remove-Item -LiteralPath $dirtyPath -Force
    Assert-VerifierExit -CaseName 'clean retained diagnostic baseline' -Fixture $retained -ExpectedExit 0 -AllowRetainedState
    Write-Output 'PASS: temporary repository rejects dirty baselines in every mode and accepts only explicit clean retained diagnostics.'
}
finally {
    if (
        $fixtureRootCreated -and
        $fixtureOwnershipMarkerCreated -and
        (Test-Path -LiteralPath $fixtureOwnershipMarker -PathType Leaf)
    ) {
        Remove-Item -LiteralPath $fixtureRoot -Recurse -Force -ErrorAction SilentlyContinue
    }
}
