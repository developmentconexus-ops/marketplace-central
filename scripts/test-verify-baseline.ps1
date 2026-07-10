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
        [int]$ExpectedExit
    )

    $previousPreference = $ErrorActionPreference
    $ErrorActionPreference = 'Continue'
    & powershell -NoProfile -ExecutionPolicy Bypass -File $verifier -Inventory $Fixture.Inventory -Ledger $Fixture.Ledger *> $null
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

    $complete = Write-Fixture -Name 'complete' -InventoryRows @(" M`talpha.txt", "??`tbeta.txt") -LedgerRows @(
        "alpha.txt`t M`tcommitted`tM-03`tdeadbeef`towner`tverified`tnone`tgit status",
        "beta.txt`t??`tcommitted`tM-04	cafebabe`towner`tverified`tnone`tgit status"
    )
    Assert-VerifierExit -CaseName 'complete ledger' -Fixture $complete -ExpectedExit 0
    Write-Output 'PASS: missing, duplicate, and invalid fixtures were rejected; complete fixture was accepted.'
}
finally {
    Remove-Item -LiteralPath $fixtureRoot -Recurse -Force -ErrorAction SilentlyContinue
}
