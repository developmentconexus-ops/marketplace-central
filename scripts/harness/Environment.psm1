Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$script:SafeToolKeys = @('SystemRoot', 'WINDIR', 'ComSpec', 'PATH', 'PATHEXT', 'TEMP', 'TMP')

function New-CaseInsensitiveStringDictionary {
  return [System.Collections.Generic.Dictionary[string,string]]::new([StringComparer]::OrdinalIgnoreCase)
}

function Read-HarnessEnvironmentFile {
  [CmdletBinding()]
  param([string]$Path)

  $values = New-CaseInsensitiveStringDictionary
  if ([string]::IsNullOrWhiteSpace($Path)) { return $values }
  if (-not (Test-Path -LiteralPath $Path -PathType Leaf)) { return $values }

  foreach ($line in Get-Content -LiteralPath $Path) {
    if ($line -match '^\s*#' -or $line -notmatch '^\s*([A-Za-z_][A-Za-z0-9_]*)\s*=') { continue }
    $key = $Matches[1]
    $value = ($line -split '=', 2)[1].Trim()
    if ($value.Length -ge 2 -and (
      ($value.StartsWith('"') -and $value.EndsWith('"')) -or
      ($value.StartsWith("'") -and $value.EndsWith("'")))) {
      $value = $value.Substring(1, $value.Length - 2)
    }
    $values[$key] = $value
  }
  return $values
}

function Get-HarnessRegistry {
  param(
    [Parameter(Mandatory)][string]$RepositoryRoot,
    [Parameter(Mandatory)][string]$Name
  )

  $path = Join-Path $RepositoryRoot "contracts/governance/$Name.json"
  if (-not (Test-Path -LiteralPath $path -PathType Leaf)) {
    throw "HENV_RUNTIME_CONTRACT_DRIFT registry_missing=$Name"
  }
  try { return Get-Content -Raw -LiteralPath $path | ConvertFrom-Json -Depth 100 }
  catch { throw "HENV_RUNTIME_CONTRACT_DRIFT registry_invalid=$Name" }
}

function Get-ExplicitRuntimeValues {
  param(
    [System.Collections.IDictionary]$RuntimeValues,
    [Parameter(Mandatory)][object]$Lane,
    [Parameter(Mandatory)][object[]]$RuntimeKeys
  )

  $values = New-CaseInsensitiveStringDictionary
  $knownKeys = @{}
  foreach ($entry in $RuntimeKeys) { $knownKeys[[string]$entry.key] = $entry }
  if ($null -ne $RuntimeValues) {
    foreach ($key in $RuntimeValues.Keys) {
      $name = [string]$key
      if (-not $knownKeys.ContainsKey($name) -or @($Lane.allowed_runtime_keys) -notcontains $name) {
        throw "HENV_RUNTIME_KEY_FORBIDDEN lane=$($Lane.id) key=$name"
      }
      $value = [string]$RuntimeValues[$key]
      if (-not [string]::IsNullOrWhiteSpace($value)) { $values[$name] = $value }
    }
  }
  return $values
}

function Get-ProcessRuntimeValues {
  param([Parameter(Mandatory)][object]$Lane)

  $values = New-CaseInsensitiveStringDictionary
  foreach ($allowedKey in @($Lane.allowed_runtime_keys)) {
    $name = [string]$allowedKey
    $value = [Environment]::GetEnvironmentVariable($name, 'Process')
    if (-not [string]::IsNullOrWhiteSpace($value)) { $values[$name] = $value }
  }
  return $values
}

function Resolve-HarnessRuntimeValues {
  param(
    [Parameter(Mandatory)][object]$Lane,
    [Parameter(Mandatory)][object[]]$RuntimeKeys,
    [Parameter(Mandatory)][System.Collections.Generic.Dictionary[string,string]]$Configured
  )

  $resolved = New-CaseInsensitiveStringDictionary
  $keysByName = @{}
  foreach ($entry in $RuntimeKeys) { $keysByName[[string]$entry.key] = $entry }

  foreach ($allowedKey in @($Lane.allowed_runtime_keys)) {
    $name = [string]$allowedKey
    if (-not $keysByName.ContainsKey($name)) {
      throw "HENV_RUNTIME_CONTRACT_DRIFT lane=$($Lane.id) key=$name"
    }
    $entry = $keysByName[$name]
    if (@($entry.allowed_lanes) -notcontains [string]$Lane.id) {
      throw "HENV_RUNTIME_CONTRACT_DRIFT lane=$($Lane.id) key=$name"
    }
  }

  foreach ($name in @($Lane.allowed_runtime_keys)) {
    $entry = $keysByName[[string]$name]
    if ([string]$entry.lifecycle -eq 'legacy_alias') { continue }
    if ($Configured.ContainsKey([string]$name)) {
      $value = $Configured[[string]$name]
      if (-not [string]::IsNullOrWhiteSpace($value)) { $resolved[[string]$name] = $value }
    }
  }

  $tupleAliases = @('SANKHYA_ORACLE_HOST', 'SANKHYA_ORACLE_PORT', 'SANKHYA_ORACLE_SERVICE_NAME')
  foreach ($entry in $RuntimeKeys | Where-Object { [string]$_.lifecycle -eq 'legacy_alias' }) {
    $alias = [string]$entry.key
    $canonical = [string]$entry.alias_for
    if ($tupleAliases -contains $alias) { continue }
    if (@($Lane.allowed_runtime_keys) -notcontains $alias) { continue }
    if ($resolved.ContainsKey($canonical) -or -not $Configured.ContainsKey($alias)) { continue }
    $value = $Configured[$alias]
    if (-not [string]::IsNullOrWhiteSpace($value)) { $resolved[$canonical] = $value }
  }

  if (-not $resolved.ContainsKey('MPC_ORACLE_CONNECT_STRING')) {
    $parts = @('SANKHYA_ORACLE_HOST', 'SANKHYA_ORACLE_PORT', 'SANKHYA_ORACLE_SERVICE_NAME')
    if (@($parts | Where-Object { -not $Configured.ContainsKey($_) -or [string]::IsNullOrWhiteSpace($Configured[$_]) }).Count -eq 0) {
      $resolved['MPC_ORACLE_CONNECT_STRING'] = "$($Configured[$parts[0]]):$($Configured[$parts[1]])/$($Configured[$parts[2]])"
    }
  }

  return $resolved
}

function New-HarnessChildEnvironment {
  [CmdletBinding()]
  param(
    [Parameter(Mandatory)][string]$RepositoryRoot,
    [Parameter(Mandatory)][string]$LaneId,
    [string]$EnvFile,
    [System.Collections.IDictionary]$RuntimeValues
  )

  if (-not [IO.Path]::IsPathFullyQualified($RepositoryRoot)) {
    throw 'HENV_REPOSITORY_ROOT_INVALID'
  }
  try { $root = [IO.Path]::GetFullPath($RepositoryRoot) }
  catch { throw 'HENV_REPOSITORY_ROOT_INVALID' }
  if (-not (Test-Path -LiteralPath $root -PathType Container)) { throw 'HENV_REPOSITORY_ROOT_INVALID' }
  $lanes = Get-HarnessRegistry -RepositoryRoot $root -Name 'execution-lanes'
  $runtime = Get-HarnessRegistry -RepositoryRoot $root -Name 'runtime-config'
  $lane = @($lanes.lanes | Where-Object { [string]$_.id -eq $LaneId })
  if ($lane.Count -ne 1) { throw "HENV_LANE_UNSUPPORTED lane=$LaneId" }
  $lane = $lane[0]
  if ([bool]$lane.inherit_parent) { throw "HENV_PARENT_INHERITANCE_FORBIDDEN lane=$LaneId" }

  $child = New-CaseInsensitiveStringDictionary
  foreach ($key in $script:SafeToolKeys) {
    $value = [Environment]::GetEnvironmentVariable($key, 'Process')
    if (-not [string]::IsNullOrWhiteSpace($value)) { $child[$key] = $value }
  }
  if (-not $child.ContainsKey('PATH')) { throw 'HENV_TOOL_KEY_MISSING key=PATH' }
  if ($LaneId -eq 'unit') {
    $child['GOCACHE'] = [IO.Path]::GetFullPath((Join-Path $root 'apps/server_core/.gocache'))
    $child['GOMODCACHE'] = [IO.Path]::GetFullPath((Join-Path $root 'apps/server_core/.gomodcache'))
    $child['GOPROXY'] = 'off'
    $child['GOSUMDB'] = 'off'
  }

  $explicitValues = Get-ExplicitRuntimeValues -RuntimeValues $RuntimeValues -Lane $lane -RuntimeKeys @($runtime.keys)
  if ($LaneId -eq 'unit') { return $child }

  $sources = @(
    (Read-HarnessEnvironmentFile -Path $EnvFile),
    (Get-ProcessRuntimeValues -Lane $lane),
    $explicitValues
  )
  foreach ($configured in $sources) {
    $resolved = Resolve-HarnessRuntimeValues -Lane $lane -RuntimeKeys @($runtime.keys) -Configured $configured
    foreach ($key in $resolved.Keys) {
      if (@($lane.allowed_runtime_keys) -notcontains $key) {
        throw "HENV_RUNTIME_KEY_FORBIDDEN lane=$LaneId key=$key"
      }
      $child[$key] = $resolved[$key]
    }
  }
  return $child
}

Export-ModuleMember -Function New-HarnessChildEnvironment, Read-HarnessEnvironmentFile
