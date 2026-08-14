# Negative fixture for issue #1 (I1-edge). Mandatory, and it is the deliverable.
#
# /orders serialises buyer name, CPF/CNPJ and billing address
# (internal/modules/orders/transport/http_handler.go:608-618) and the composed
# middleware chain has no identity check at all
# (internal/composition/root.go:994). The edge is therefore the only control on
# that route, and the boot assertion that was supposed to retire it is deferred
# along with authentication, which the rebaseline decides at D2 (ADR-035). So
# "the door is closed" is a claim
# about two config files nobody re-reads, unless something re-reads them. This is
# that something.
#
# It drives the REAL deploy/Caddyfile and the REAL docker/dev/oauth-edge.Caddyfile
# through real Caddy, against stub upstreams, and asserts:
#   - non-200 on /orders JSON without credentials, and with wrong credentials;
#   - 200 WITH credentials, from the backend stub — a positive control, because a
#     stack that is simply broken also returns non-200 and would pass a
#     deny-only assertion;
#   - the SPA hard-nav (Accept: text/html) still reaches the web stub;
#   - /healthz still reaches the backend — the control is scoped, not a blanket
#     outage;
#   - the ngrok edge publishes the OAuth callback and denies everything else.
#
# Run: npm run harness:edge   (pwsh scripts/harness.ps1 -Command edge)
# Requires: Docker. No ERP, no dev stack, no provider credentials, no secrets.

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$runId = [guid]::NewGuid().ToString('N').Substring(0, 12)
$network = "mpc-edge-$runId"
$containers = [System.Collections.Generic.List[string]]::new()
$assertions = 0

function Resolve-DockerApplication {
  $command = Get-Command 'docker' -CommandType Application -ErrorAction SilentlyContinue | Select-Object -First 1
  if ($null -eq $command) { throw 'HEDGE_FILE_NOT_FOUND tool=docker' }
  return [IO.Path]::GetFullPath($command.Source)
}

$docker = Resolve-DockerApplication

function Invoke-Docker {
  param([Parameter(Mandatory)][string[]]$Arguments, [switch]$IgnoreFailure)
  $output = & $docker @Arguments 2>&1
  if ($LASTEXITCODE -ne 0 -and -not $IgnoreFailure) {
    throw "HEDGE_DOCKER_FAILED args=$($Arguments -join ' ') exit_code=$LASTEXITCODE output=$($output -join ' ')"
  }
  return ($output | Out-String).Trim()
}

function New-StubContainer {
  param(
    [Parameter(Mandatory)][string]$Name,
    [Parameter(Mandatory)][string]$Alias,
    [Parameter(Mandatory)][string]$CaddyfilePath
  )
  $full = [IO.Path]::GetFullPath($CaddyfilePath)
  if (-not (Test-Path -LiteralPath $full -PathType Leaf)) { throw "HEDGE_FIXTURE_MISSING path=$full" }
  Invoke-Docker -Arguments @(
    'run', '-d', '--name', $Name, '--network', $network, '--network-alias', $Alias,
    '-v', "${full}:/etc/caddy/Caddyfile:ro", 'caddy:2-alpine'
  ) | Out-Null
  $containers.Add($Name)
}

function New-EdgeContainer {
  param(
    [Parameter(Mandatory)][string]$Name,
    [Parameter(Mandatory)][string]$CaddyfilePath,
    [string[]]$EnvironmentPairs = @()
  )
  $full = [IO.Path]::GetFullPath($CaddyfilePath)
  if (-not (Test-Path -LiteralPath $full -PathType Leaf)) { throw "HEDGE_FIXTURE_MISSING path=$full" }
  $arguments = [System.Collections.Generic.List[string]]::new()
  $arguments.AddRange([string[]]@('run', '-d', '--name', $Name, '--network', $network))
  foreach ($pair in $EnvironmentPairs) { $arguments.AddRange([string[]]@('-e', $pair)) }
  $arguments.AddRange([string[]]@('-p', '127.0.0.1::80', '-v', "${full}:/etc/caddy/Caddyfile:ro", 'caddy:2-alpine'))
  Invoke-Docker -Arguments $arguments.ToArray() | Out-Null
  $containers.Add($Name)

  $mapping = Invoke-Docker -Arguments @('port', $Name, '80')
  $port = ($mapping -split "`n" | Select-Object -First 1).Split(':')[-1].Trim()
  if ([string]::IsNullOrWhiteSpace($port) -or ($port -as [int]) -lt 1) {
    throw "HEDGE_PORT_UNRESOLVED container=$Name mapping=$mapping"
  }
  return [int]$port
}

function Wait-EdgeReady {
  param([Parameter(Mandatory)][int]$Port, [Parameter(Mandatory)][string]$Name)
  for ($attempt = 1; $attempt -le 40; $attempt++) {
    try {
      Invoke-WebRequest -Uri "http://127.0.0.1:$Port/" -Method Get -SkipHttpErrorCheck -MaximumRedirection 0 -TimeoutSec 5 | Out-Null
      return
    } catch {
      Start-Sleep -Milliseconds 500
    }
  }
  $logs = Invoke-Docker -Arguments @('logs', $Name) -IgnoreFailure
  throw "HEDGE_EDGE_NOT_READY container=$Name port=$Port logs=$logs"
}

function Get-EdgeResponse {
  param(
    [Parameter(Mandatory)][int]$Port,
    [Parameter(Mandatory)][string]$Path,
    [string]$Accept = 'application/json',
    [string]$User,
    [string]$Password
  )
  $headers = @{ Accept = $Accept }
  if (-not [string]::IsNullOrEmpty($User)) {
    $token = [Convert]::ToBase64String([Text.Encoding]::UTF8.GetBytes("${User}:${Password}"))
    $headers['Authorization'] = "Basic $token"
  }
  return Invoke-WebRequest -Uri "http://127.0.0.1:$Port$Path" -Method Get -Headers $headers `
    -SkipHttpErrorCheck -MaximumRedirection 0 -TimeoutSec 20
}

function Assert-Status {
  param(
    [Parameter(Mandatory)][object]$Response,
    [Parameter(Mandatory)][int]$Expected,
    [Parameter(Mandatory)][string]$Because
  )
  if ([int]$Response.StatusCode -ne $Expected) {
    throw "failure_token=test=$Because expected_status=$Expected actual_status=$([int]$Response.StatusCode)"
  }
  $script:assertions++
  Write-Output "ok $Because status=$Expected"
}

function Assert-Body {
  param(
    [Parameter(Mandatory)][object]$Response,
    [Parameter(Mandatory)][string]$Expected,
    [Parameter(Mandatory)][string]$Because
  )
  $body = [string]$Response.Content
  if ($body -notlike "*$Expected*") {
    throw "failure_token=test=$Because expected_body_substring=$Expected actual_body=$body"
  }
  $script:assertions++
  Write-Output "ok $Because body~=$Expected"
}

try {
  Invoke-Docker -Arguments @('network', 'create', $network) | Out-Null

  # A password generated per run: nothing to commit, nothing to leak, and no
  # value in this file that could ever be mistaken for a real credential.
  $passwordBytes = [byte[]]::new(24)
  [Security.Cryptography.RandomNumberGenerator]::Fill($passwordBytes)
  $password = [Convert]::ToHexString($passwordBytes).ToLowerInvariant()
  $user = 'fixture-ops'
  $hash = Invoke-Docker -Arguments @('run', '--rm', 'caddy:2-alpine', 'caddy', 'hash-password', '--plaintext', $password)
  if ($hash -notlike '$2a$*') { throw "HEDGE_HASH_UNEXPECTED value_shape=$($hash.Substring(0, [Math]::Min(4, $hash.Length)))" }

  New-StubContainer -Name "mpc-edge-backend-$runId" -Alias 'backend' -CaddyfilePath (Join-Path $repoRoot 'scripts/tests/fixtures/edge/backend-stub.Caddyfile')
  New-StubContainer -Name "mpc-edge-web-$runId" -Alias 'web' -CaddyfilePath (Join-Path $repoRoot 'scripts/tests/fixtures/edge/web-stub.Caddyfile')

  # ---------------------------------------------------------------------------
  # Part 1 — the production edge, deploy/Caddyfile, verbatim.
  # ---------------------------------------------------------------------------
  $prodPort = New-EdgeContainer -Name "mpc-edge-prod-$runId" -CaddyfilePath (Join-Path $repoRoot 'deploy/Caddyfile') -EnvironmentPairs @(
    'MPC_DOMAIN=:80',
    "MPC_ORDERS_BASIC_AUTH_USER=$user",
    "MPC_ORDERS_BASIC_AUTH_HASH=$hash"
  )
  Wait-EdgeReady -Port $prodPort -Name "mpc-edge-prod-$runId"

  # The claim of issue #1: non-200 without credentials.
  Assert-Status -Response (Get-EdgeResponse -Port $prodPort -Path '/orders') -Expected 401 -Because 'orders-json-no-credentials-is-denied'
  Assert-Status -Response (Get-EdgeResponse -Port $prodPort -Path '/orders/ML-123') -Expected 401 -Because 'orders-subpath-no-credentials-is-denied'
  Assert-Status -Response (Get-EdgeResponse -Port $prodPort -Path '/orders' -User $user -Password 'wrong-password') -Expected 401 -Because 'orders-json-wrong-credentials-is-denied'

  # Positive controls. A stack that is merely broken also returns non-200; these
  # are what separate "the door is closed" from "the building fell down".
  $authorised = Get-EdgeResponse -Port $prodPort -Path '/orders' -User $user -Password $password
  Assert-Status -Response $authorised -Expected 200 -Because 'orders-json-with-credentials-reaches-backend'
  Assert-Body -Response $authorised -Expected 'backend-stub-body' -Because 'orders-json-with-credentials-is-answered-by-the-backend'

  $hardNav = Get-EdgeResponse -Port $prodPort -Path '/orders' -Accept 'text/html'
  Assert-Status -Response $hardNav -Expected 200 -Because 'orders-hard-nav-still-serves-the-spa'
  Assert-Body -Response $hardNav -Expected 'web-stub-index-html' -Because 'orders-hard-nav-is-answered-by-the-spa'

  $health = Get-EdgeResponse -Port $prodPort -Path '/healthz'
  Assert-Status -Response $health -Expected 200 -Because 'healthz-is-unaffected'
  Assert-Body -Response $health -Expected 'backend-stub-healthz' -Because 'healthz-still-reaches-the-backend'

  # ---------------------------------------------------------------------------
  # Part 2 — the ngrok edge, docker/dev/oauth-edge.Caddyfile, verbatim.
  # ---------------------------------------------------------------------------
  $oauthPort = New-EdgeContainer -Name "mpc-edge-oauth-$runId" -CaddyfilePath (Join-Path $repoRoot 'docker/dev/oauth-edge.Caddyfile')
  Wait-EdgeReady -Port $oauthPort -Name "mpc-edge-oauth-$runId"

  $callback = Get-EdgeResponse -Port $oauthPort -Path '/integrations/auth/callback?code=c&state=s'
  Assert-Status -Response $callback -Expected 200 -Because 'tunnel-publishes-the-oauth-callback'
  Assert-Body -Response $callback -Expected 'backend-stub-body' -Because 'tunnel-callback-reaches-the-backend'

  Assert-Status -Response (Get-EdgeResponse -Port $oauthPort -Path '/orders') -Expected 403 -Because 'tunnel-does-not-publish-orders'
  Assert-Status -Response (Get-EdgeResponse -Port $oauthPort -Path '/orders/ML-123') -Expected 403 -Because 'tunnel-does-not-publish-orders-subpath'
  Assert-Status -Response (Get-EdgeResponse -Port $oauthPort -Path '/orders' -Accept 'text/html') -Expected 403 -Because 'tunnel-does-not-publish-orders-as-html-either'
  Assert-Status -Response (Get-EdgeResponse -Port $oauthPort -Path '/dashboard') -Expected 403 -Because 'tunnel-does-not-publish-the-api-surface'
  Assert-Status -Response (Get-EdgeResponse -Port $oauthPort -Path '/integrations/installations/abc/auth/start') -Expected 403 -Because 'tunnel-does-not-publish-the-rest-of-integrations'

  Assert-Status -Response (Get-EdgeResponse -Port $oauthPort -Path '/integrations' -Accept 'text/html') -Expected 200 -Because 'tunnel-answers-the-callback-landing-in-place'

  # Anti-vacuity: a run that asserted nothing is byte-identical to a green run
  # without it, unless the count itself is asserted.
  if ($assertions -lt 17) { throw "failure_token=test=assertion-count-too-low assertions=$assertions expected_min=17" }
  Write-Output "assertions=$assertions"
  Write-Output "PASS edge PII deny run_id=$runId prod_port=$prodPort oauth_port=$oauthPort assertions=$assertions"
} finally {
  foreach ($container in $containers) {
    Invoke-Docker -Arguments @('rm', '-f', $container) -IgnoreFailure | Out-Null
  }
  Invoke-Docker -Arguments @('network', 'rm', $network) -IgnoreFailure | Out-Null
}
