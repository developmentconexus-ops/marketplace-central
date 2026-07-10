$ErrorActionPreference = 'Stop'

$verifier = Join-Path $PSScriptRoot 'verify-baseline.ps1'
$fixtureRoot = Join-Path $env:TEMP ("verify-baseline-" + [guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Path $fixtureRoot | Out-Null

function Write-Fixture {
    param(
        [string]$Name,
        [string[]]$InventoryRows,
        [string[]]$LedgerRows
    )

    $inventory = Join-Path $fixtureRoot "$Name-inventory.tsv"
    $ledger = Join-Path $fixtureRoot "$Name-ledger.tsv"
    @("status`tpath") + $InventoryRows | Set-Content -Path $inventory -Encoding utf8
    @("path`tstatus`tdisposition`tcohort`tcommit_sha`towner`treason`tblocker`tnext_validation") + $LedgerRows |
        Set-Content -Path $ledger -Encoding utf8
    return @{ Inventory = $inventory; Ledger = $ledger }
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
    & powershell @verifierArguments *> $null
    $actualExit = $LASTEXITCODE
    $ErrorActionPreference = $previousPreference
    if ($actualExit -ne $ExpectedExit) {
        throw "$CaseName expected verifier exit $ExpectedExit but received $actualExit."
    }
}

try {
    $missing = Write-Fixture -Name 'missing' -InventoryRows @(" M`talpha.txt", "??`tbeta.txt") -LedgerRows @(
        "alpha.txt`t M`tcommitted`tM-03`tdeadbeef`towner`tverified`tnone`tgit status"
    )
    Assert-VerifierExit -CaseName 'missing path' -Fixture $missing -ExpectedExit 1

    $duplicate = Write-Fixture -Name 'duplicate' -InventoryRows @(" M`talpha.txt") -LedgerRows @(
        "alpha.txt`t M`tcommitted`tM-03`tdeadbeef`towner`tverified`tnone`tgit status",
        "alpha.txt`t M`tcommitted`tM-03`tdeadbeef`towner`tverified`tnone`tgit status"
    )
    Assert-VerifierExit -CaseName 'duplicate ledger path' -Fixture $duplicate -ExpectedExit 1

    $invalid = Write-Fixture -Name 'invalid' -InventoryRows @(" M`talpha.txt") -LedgerRows @(
        "alpha.txt`t M`tunknown-disposition`tM-03`tdeadbeef`towner`tverified`tnone`tgit status"
    )
    Assert-VerifierExit -CaseName 'invalid disposition' -Fixture $invalid -ExpectedExit 1

    $retained = Write-Fixture -Name 'retained' -InventoryRows @(" M`talpha.txt") -LedgerRows @(
        "alpha.txt`t M`tretained-owner-needed`tM-08`t`towner`tawaiting ownership`tscoped validation needed`tgit status"
    )
    Assert-VerifierExit -CaseName 'retained state requires explicit opt-in' -Fixture $retained -ExpectedExit 1
    Assert-VerifierExit -CaseName 'retained state diagnostic opt-in' -Fixture $retained -ExpectedExit 0 -AllowRetainedState

    $complete = Write-Fixture -Name 'complete' -InventoryRows @(" M`talpha.txt", "??`tbeta.txt") -LedgerRows @(
        "alpha.txt`t M`tcommitted`tM-03`tdeadbeef`towner`tverified`tnone`tgit status",
        "beta.txt`t??`tcommitted`tM-04	cafebabe`towner`tverified`tnone`tgit status"
    )
    Assert-VerifierExit -CaseName 'dirty worktree requires explicit opt-in' -Fixture $complete -ExpectedExit 1
    Assert-VerifierExit -CaseName 'complete ledger diagnostic opt-in' -Fixture $complete -ExpectedExit 0 -AllowRetainedState
    Write-Output 'PASS: invalid fixtures and default retained/dirty state were rejected; explicit diagnostic opt-in was accepted.'
}
finally {
    Remove-Item -LiteralPath $fixtureRoot -Recurse -Force -ErrorAction SilentlyContinue
}
