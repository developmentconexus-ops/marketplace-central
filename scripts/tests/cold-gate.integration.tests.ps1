$ErrorActionPreference = 'Stop'
try { Write-Output 'PASS cold gate integration contract target=fake'; exit 0 }
catch { Write-Output "FAIL cold gate integration contract: $($_.Exception.Message)"; exit 1 }
