$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$module = Join-Path $repoRoot 'scripts/harness/Execution.psm1'
$probe = Join-Path $repoRoot 'scripts/tests/fixtures/harness/child-probe.mjs'

function Assert-True([bool]$Condition, [string]$Message) {
  if (-not $Condition) { throw $Message }
}

if (-not (Test-Path -LiteralPath $module -PathType Leaf)) {
  throw 'RED: Execution.psm1 missing; typed shell-free process contract is not implemented'
}
Import-Module $module -Force

$node = (Get-Command node -CommandType Application -ErrorAction Stop).Source
$originalLocation = (Get-Location).Path
$originalCache = [Environment]::GetEnvironmentVariable('GOCACHE', 'Process')
$baseEnvironment = [System.Collections.Generic.Dictionary[string,string]]::new([StringComparer]::OrdinalIgnoreCase)
foreach ($name in @('SystemRoot', 'WINDIR', 'ComSpec', 'PATH', 'PATHEXT', 'TEMP', 'TMP')) {
  $value = [Environment]::GetEnvironmentVariable($name, 'Process')
  if (-not [string]::IsNullOrWhiteSpace($value)) { $baseEnvironment[$name] = $value }
}

$literalArguments = @('space value', 'quote"value', 'amp&value', 'literal$()', 'semi;value')
$inspectRequest = New-HarnessProcessRequest -FilePath $node -ArgumentList (@($probe, 'inspect') + $literalArguments) -WorkingDirectory $repoRoot -Environment $baseEnvironment -TimeoutSeconds 10
$inspect = Invoke-HarnessProcess -Request $inspectRequest
Assert-True ($inspect.ExitCode -eq 0) 'argument probe failed'
$observed = $inspect.Stdout | ConvertFrom-Json
Assert-True (($observed.args -join "`n") -eq ($literalArguments -join "`n")) 'arguments did not round-trip literally'
Assert-True ([IO.Path]::GetFullPath($observed.cwd) -eq [IO.Path]::GetFullPath($repoRoot)) 'child CWD was not explicit'
$observedKeys = @($observed.env.PSObject.Properties.Name | Sort-Object)
$expectedKeys = @($baseEnvironment.Keys | Sort-Object)
Assert-True (($observedKeys -join "`n") -eq ($expectedKeys -join "`n")) 'child process environment was not exact'

$relativeWorkingDirectoryError = ''
try {
  New-HarnessProcessRequest -FilePath $node -ArgumentList @($probe, 'inspect') -WorkingDirectory '.' -Environment $baseEnvironment -TimeoutSeconds 10 | Out-Null
} catch { $relativeWorkingDirectoryError = $_.Exception.Message }
Assert-True ($relativeWorkingDirectoryError -match 'HEXEC_WORKING_DIRECTORY_INVALID') 'relative working directory was accepted'

Set-Location -LiteralPath ([IO.Path]::GetTempPath())
try {
  $foreignRequest = New-HarnessProcessRequest -FilePath $node -ArgumentList @($probe, 'inspect') -WorkingDirectory $repoRoot -Environment $baseEnvironment -TimeoutSeconds 10
  $foreignResult = Invoke-HarnessProcess -Request $foreignRequest
  $foreignObserved = $foreignResult.Stdout | ConvertFrom-Json
  Assert-True ([IO.Path]::GetFullPath($foreignObserved.cwd) -eq [IO.Path]::GetFullPath($repoRoot)) 'absolute working directory depends on caller CWD'
} finally { Set-Location -LiteralPath $originalLocation }

$streamRequest = New-HarnessProcessRequest -FilePath $node -ArgumentList @($probe, 'streams', '2200000') -WorkingDirectory $repoRoot -Environment $baseEnvironment -TimeoutSeconds 20
$streams = Invoke-HarnessProcess -Request $streamRequest
Assert-True ($streams.ExitCode -eq 0 -and $streams.Stdout.Length -eq 2200000 -and $streams.Stderr.Length -eq 2200000) 'concurrent large streams were truncated or deadlocked'

$exitRequest = New-HarnessProcessRequest -FilePath $node -ArgumentList @($probe, 'exit', '17') -WorkingDirectory $repoRoot -Environment $baseEnvironment -TimeoutSeconds 10
$failed = Invoke-HarnessProcess -Request $exitRequest
Assert-True ($failed.ExitCode -eq 17 -and $failed.ReasonCode -eq 'HEXEC_EXIT_NONZERO') 'child exit 17 was not preserved'

$candidate = 'fixture:@/' + [guid]::NewGuid().ToString('N')
$secretEnvironment = [System.Collections.Generic.Dictionary[string,string]]::new($baseEnvironment, [StringComparer]::OrdinalIgnoreCase)
$secretEnvironment['HARNESS_TEST_CANDIDATE'] = $candidate
$redactRequest = New-HarnessProcessRequest -FilePath $node -ArgumentList @($probe, 'redact') -WorkingDirectory $repoRoot -Environment $secretEnvironment -TimeoutSeconds 10 -RedactionCandidates @($candidate)
$redacted = Invoke-HarnessProcess -Request $redactRequest
Assert-True (($redacted.Stdout + $redacted.Stderr) -notmatch [regex]::Escape($candidate)) 'credential escaped redaction'
Assert-True (($redacted.Stdout + $redacted.Stderr) -notmatch [regex]::Escape([uri]::EscapeDataString($candidate))) 'URI-encoded credential escaped redaction'
Assert-True (($redacted.Stdout + $redacted.Stderr) -match '\[redacted\]') 'redaction marker missing'

$implicitCandidate = 'implicit:@/' + [guid]::NewGuid().ToString('N')
$implicitEnvironment = [System.Collections.Generic.Dictionary[string,string]]::new($baseEnvironment, [StringComparer]::OrdinalIgnoreCase)
$implicitEnvironment['HARNESS_TEST_CANDIDATE'] = $implicitCandidate
$implicitRequest = New-HarnessProcessRequest -FilePath $node -ArgumentList @($probe, 'redact') -WorkingDirectory $repoRoot -Environment $implicitEnvironment -TimeoutSeconds 10
$implicitRedacted = Invoke-HarnessProcess -Request $implicitRequest
Assert-True ($implicitRedacted.Stdout -notmatch [regex]::Escape($implicitCandidate)) 'standalone stdout leaked an implicit environment candidate'
Assert-True ($implicitRedacted.Stderr -notmatch [regex]::Escape($implicitCandidate)) 'standalone stderr leaked an implicit environment candidate'
Assert-True (($implicitRedacted.Stdout + $implicitRedacted.Stderr) -match '\[redacted\]') 'implicit environment redaction marker missing'

$marker = Join-Path ([IO.Path]::GetTempPath()) ('harness-tree-' + [guid]::NewGuid().ToString('N'))
$timeoutRequest = New-HarnessProcessRequest -FilePath $node -ArgumentList @($probe, 'tree', $marker) -WorkingDirectory $repoRoot -Environment $baseEnvironment -TimeoutSeconds 1
$timedOut = Invoke-HarnessProcess -Request $timeoutRequest
Assert-True ($timedOut.ReasonCode -eq 'HEXEC_TIMEOUT') 'timeout lacks stable reason code'
Start-Sleep -Seconds 3
Assert-True (-not (Test-Path -LiteralPath $marker)) 'timeout did not kill process tree'

$executionSource = Get-Content -Raw -LiteralPath $module
Assert-True ($executionSource -notmatch '\.WaitForExit\(\s*\)') 'timeout path contains unbounded WaitForExit'
Assert-True ($executionSource -notmatch 'GetAwaiter\(\)\.GetResult\(\)') 'stream drain contains an unbounded task wait'

Assert-True ((Get-Location).Path -eq $originalLocation) 'parent location mutated'
Assert-True ([Environment]::GetEnvironmentVariable('GOCACHE', 'Process') -eq $originalCache) 'parent environment mutated'
Remove-Item -LiteralPath $marker -Force -ErrorAction SilentlyContinue

Write-Output 'PASS harness execution tests'
