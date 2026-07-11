Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

class HarnessProcessRequest {
  [string]$FilePath
  [string[]]$ArgumentList
  [string]$WorkingDirectory
  [System.Collections.Generic.Dictionary[string,string]]$Environment
  [int]$TimeoutSeconds
  [string[]]$RedactionCandidates
}

class HarnessProcessResult {
  [int]$ExitCode
  [string]$Stdout
  [string]$Stderr
  [string]$ReasonCode
  [bool]$TimedOut
}

function Protect-HarnessText {
  param(
    [AllowNull()][string]$Text,
    [string[]]$Candidates
  )

  if ($null -eq $Text) { return '' }
  $safe = $Text
  $forms = [System.Collections.Generic.HashSet[string]]::new([StringComparer]::Ordinal)
  foreach ($candidate in @($Candidates)) {
    if ([string]::IsNullOrWhiteSpace($candidate) -or $candidate.Length -lt 3) { continue }
    [void]$forms.Add($candidate)
    [void]$forms.Add([uri]::EscapeDataString($candidate))
  }
  foreach ($form in @($forms | Sort-Object Length -Descending)) {
    $safe = $safe.Replace($form, '[redacted]', [StringComparison]::Ordinal)
  }
  $safe = [regex]::Replace(
    $safe,
    '(?i)(\b[A-Z][A-Z0-9_]*(?:PASSWORD|SECRET|TOKEN|CREDENTIAL|CONNECT_STRING)\b\s*=\s*)[^,;\s]+',
    '$1[redacted]')
  $safe = [regex]::Replace($safe, '(?i)(\b(?:postgres(?:ql)?|mysql|oracle)://)[^\s,;]+', '$1[redacted]')
  return $safe
}

function New-HarnessProcessRequest {
  [CmdletBinding()]
  param(
    [Parameter(Mandatory)][string]$FilePath,
    [Parameter(Mandatory)][AllowEmptyCollection()][string[]]$ArgumentList,
    [Parameter(Mandatory)][string]$WorkingDirectory,
    [Parameter(Mandatory)][System.Collections.IDictionary]$Environment,
    [ValidateRange(1, 86400)][int]$TimeoutSeconds = 600,
    [string[]]$RedactionCandidates = @()
  )

  if (-not [IO.Path]::IsPathFullyQualified($FilePath) -or -not (Test-Path -LiteralPath $FilePath -PathType Leaf)) {
    throw 'HEXEC_FILE_NOT_FOUND'
  }
  if (-not [IO.Path]::IsPathFullyQualified($WorkingDirectory)) {
    throw 'HEXEC_WORKING_DIRECTORY_INVALID'
  }
  try { $canonicalWorkingDirectory = [IO.Path]::GetFullPath($WorkingDirectory) }
  catch { throw 'HEXEC_WORKING_DIRECTORY_INVALID' }
  if (-not (Test-Path -LiteralPath $canonicalWorkingDirectory -PathType Container)) { throw 'HEXEC_WORKING_DIRECTORY_INVALID' }

  $request = [HarnessProcessRequest]::new()
  $request.FilePath = [IO.Path]::GetFullPath($FilePath)
  $request.ArgumentList = @($ArgumentList | ForEach-Object { [string]$_ })
  $request.WorkingDirectory = $canonicalWorkingDirectory
  $request.Environment = [System.Collections.Generic.Dictionary[string,string]]::new([StringComparer]::OrdinalIgnoreCase)
  foreach ($key in $Environment.Keys) { $request.Environment[[string]$key] = [string]$Environment[$key] }
  $request.TimeoutSeconds = $TimeoutSeconds
  $request.RedactionCandidates = @($RedactionCandidates | Where-Object { -not [string]::IsNullOrWhiteSpace($_) })
  return $request
}

function New-HarnessProcessResult {
  param(
    [int]$ExitCode,
    [string]$Stdout,
    [string]$Stderr,
    [string]$ReasonCode,
    [bool]$TimedOut,
    [string[]]$Candidates
  )

  $result = [HarnessProcessResult]::new()
  $result.ExitCode = $ExitCode
  $result.Stdout = Protect-HarnessText -Text $Stdout -Candidates $Candidates
  $result.Stderr = Protect-HarnessText -Text $Stderr -Candidates $Candidates
  $result.ReasonCode = $ReasonCode
  $result.TimedOut = $TimedOut
  return $result
}

function Read-CompletedHarnessStream {
  param(
    [Parameter(Mandatory)][System.Threading.Tasks.Task[string]]$Task,
    [ValidateRange(0, 10000)][int]$WaitMilliseconds
  )

  if (-not $Task.IsCompleted -and $WaitMilliseconds -gt 0) {
    try { [void]$Task.Wait($WaitMilliseconds) } catch { return '' }
  }
  if (-not $Task.IsCompletedSuccessfully) { return '' }
  return [string]$Task.Result
}

function Invoke-HarnessProcess {
  [CmdletBinding()]
  param([Parameter(Mandatory)][HarnessProcessRequest]$Request)

  $startInfo = [Diagnostics.ProcessStartInfo]::new()
  $startInfo.FileName = $Request.FilePath
  $startInfo.WorkingDirectory = $Request.WorkingDirectory
  $startInfo.UseShellExecute = $false
  $startInfo.CreateNoWindow = $true
  $startInfo.RedirectStandardOutput = $true
  $startInfo.RedirectStandardError = $true
  $startInfo.Environment.Clear()
  foreach ($key in $Request.Environment.Keys) { $startInfo.Environment[$key] = $Request.Environment[$key] }
  foreach ($argument in $Request.ArgumentList) { [void]$startInfo.ArgumentList.Add($argument) }

  $process = [Diagnostics.Process]::new()
  $process.StartInfo = $startInfo
  try {
    try {
      if (-not $process.Start()) {
        return New-HarnessProcessResult -ExitCode -1 -Stdout '' -Stderr 'HEXEC_START_FAILED' -ReasonCode 'HEXEC_START_FAILED' -TimedOut $false -Candidates $Request.RedactionCandidates
      }
    } catch {
      $message = Protect-HarnessText -Text $_.Exception.Message -Candidates $Request.RedactionCandidates
      return New-HarnessProcessResult -ExitCode -1 -Stdout '' -Stderr "HEXEC_START_FAILED $message" -ReasonCode 'HEXEC_START_FAILED' -TimedOut $false -Candidates $Request.RedactionCandidates
    }

    $stdoutTask = $process.StandardOutput.ReadToEndAsync()
    $stderrTask = $process.StandardError.ReadToEndAsync()
    $completed = $process.WaitForExit($Request.TimeoutSeconds * 1000)
    if (-not $completed) {
      try { $process.Kill($true) } catch { }
      [void]$process.WaitForExit(2000)
    }
    $drainWait = if ($completed) { 5000 } else { 2000 }
    $stdout = Read-CompletedHarnessStream -Task $stdoutTask -WaitMilliseconds $drainWait
    $stderr = Read-CompletedHarnessStream -Task $stderrTask -WaitMilliseconds $drainWait
    if (-not $completed) {
      return New-HarnessProcessResult -ExitCode -1 -Stdout $stdout -Stderr $stderr -ReasonCode 'HEXEC_TIMEOUT' -TimedOut $true -Candidates $Request.RedactionCandidates
    }
    $reason = if ($process.ExitCode -eq 0) { '' } else { 'HEXEC_EXIT_NONZERO' }
    return New-HarnessProcessResult -ExitCode $process.ExitCode -Stdout $stdout -Stderr $stderr -ReasonCode $reason -TimedOut $false -Candidates $Request.RedactionCandidates
  } finally {
    $process.Dispose()
  }
}

Export-ModuleMember -Function New-HarnessProcessRequest, Invoke-HarnessProcess
