[CmdletBinding()]
param(
    [Parameter(Mandatory)]
    [string]$Inventory,
    [Parameter(Mandatory)]
    [string]$Ledger,
    [switch]$RequireCleanStatus
)

$ErrorActionPreference = 'Stop'
$validDispositions = @(
    'committed',
    'retained-owner-needed',
    'retained-secret-pii',
    'retained-validation-failed'
)

function Read-Tsv {
    param([string]$Path)

    if (-not (Test-Path -LiteralPath $Path -PathType Leaf)) {
        throw "Required input does not exist: $Path"
    }
    return @(Import-Csv -LiteralPath $Path -Delimiter "`t")
}

function Read-Ledger {
    param([string]$Path)

    if ([IO.Path]::GetExtension($Path) -ne '.md') {
        return Read-Tsv -Path $Path
    }

    $lines = @(Get-Content -LiteralPath $Path)
    $start = [array]::IndexOf($lines, '<!-- baseline-ledger:start -->')
    $end = [array]::IndexOf($lines, '<!-- baseline-ledger:end -->')
    if ($start -lt 0 -or $end -le $start + 1) {
        throw 'Ledger Markdown must contain a non-empty baseline-ledger block.'
    }
    return @($lines[($start + 1)..($end - 1)] | ConvertFrom-Csv -Delimiter "`t")
}

$inventoryRows = Read-Tsv -Path $Inventory
$ledgerRows = Read-Ledger -Path $Ledger
$errors = [Collections.Generic.List[string]]::new()

if ($inventoryRows.Count -eq 0) {
    $errors.Add('Inventory has no rows.')
}

$inventoryByPath = [Collections.Generic.Dictionary[string, object]]::new([StringComparer]::Ordinal)
foreach ($row in $inventoryRows) {
    if ([string]::IsNullOrWhiteSpace($row.path) -or [string]::IsNullOrWhiteSpace($row.status)) {
        $errors.Add('Inventory contains a row without path or status.')
        continue
    }
    if ($inventoryByPath.ContainsKey($row.path)) {
        $errors.Add("Inventory duplicates path: $($row.path)")
    }
    else {
        $inventoryByPath.Add($row.path, $row)
    }
}

$ledgerByPath = [Collections.Generic.Dictionary[string, Collections.Generic.List[object]]]::new([StringComparer]::Ordinal)
foreach ($row in $ledgerRows) {
    if ([string]::IsNullOrWhiteSpace($row.path)) {
        $errors.Add('Ledger contains a row without path.')
        continue
    }
    if ($validDispositions -notcontains $row.disposition) {
        $errors.Add("Ledger has invalid disposition for path: $($row.path)")
    }
    foreach ($field in @('status', 'cohort', 'owner', 'reason', 'blocker', 'next_validation')) {
        if ([string]::IsNullOrWhiteSpace($row.$field)) {
            $errors.Add("Ledger path is missing ${field}: $($row.path)")
        }
    }
    if ($row.disposition -eq 'committed' -and [string]::IsNullOrWhiteSpace($row.commit_sha)) {
        $errors.Add("Committed ledger path is missing commit_sha: $($row.path)")
    }
    if (-not $ledgerByPath.ContainsKey($row.path)) {
        $ledgerByPath[$row.path] = [Collections.Generic.List[object]]::new()
    }
    $ledgerByPath[$row.path].Add($row)
}

foreach ($path in $inventoryByPath.Keys) {
    if (-not $ledgerByPath.ContainsKey($path)) {
        $errors.Add("Ledger is missing inventory path: $path")
        continue
    }
    $rows = $ledgerByPath[$path]
    if ($rows.Count -ne 1) {
        $errors.Add("Ledger path occurs $($rows.Count) times: $path")
        continue
    }
    if ($rows[0].status -ne $inventoryByPath[$path].status) {
        $errors.Add("Ledger status differs from inventory for path: $path")
    }
}

foreach ($path in $ledgerByPath.Keys) {
    if (-not $inventoryByPath.ContainsKey($path)) {
        $errors.Add("Ledger contains a path absent from inventory: $path")
    }
}

$retained = @($ledgerRows | Where-Object { $_.disposition -ne 'committed' })
if ($RequireCleanStatus) {
    if ($retained.Count -gt 0) {
        $errors.Add("Clean baseline is blocked by $($retained.Count) retained-state record(s).")
    }
    $status = @(git status --porcelain=v1 --untracked-files=all)
    if ($LASTEXITCODE -ne 0) {
        $errors.Add('git status failed while checking the candidate baseline.')
    }
    elseif ($status.Count -gt 0) {
        $errors.Add("Candidate baseline has $($status.Count) dirty path(s).")
    }
}

if ($errors.Count -gt 0) {
    $errors | ForEach-Object { Write-Error $_ }
    exit 1
}

$committed = @($ledgerRows | Where-Object { $_.disposition -eq 'committed' }).Count
Write-Output "PASS: inventory=$($inventoryRows.Count) committed=$committed retained=$($retained.Count)"
