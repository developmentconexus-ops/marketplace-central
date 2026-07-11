param([Parameter(Mandatory)][ValidateSet('one','two')][string]$Name)
Write-Output "impact-probe=$Name"
