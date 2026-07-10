# Environment and Database Access

This page documents runtime environment keys and practical database access for local development.

## Runtime Sources

- `.env` is loaded by [start-server.ps1](../../.claude/start-server.ps1)
- Preview server definition is in [launch.json](../../.claude/launch.json)

## Relevant `.env` Keys

Do not store secrets in wiki. Keep values in local `.env` only.

Server/runtime:
- `API_PORT`
- `MPC_WEB_ORIGIN`
- `MC_DEFAULT_TENANT_ID`

Marketplace auth providers:
- `MPC_PROVIDER_MERCADOLIVRE_CLIENT_ID`
- `MPC_PROVIDER_MERCADOLIVRE_CLIENT_SECRET`
- `MPC_PROVIDER_MAGALU_CLIENT_ID`
- `MPC_PROVIDER_MAGALU_CLIENT_SECRET`
- `MPC_OAUTH_REDIRECT_URI`
- `MPC_OAUTH_HMAC_SECRET`
- `MPC_ENCRYPTION_KEY`

Databases:
- `MC_DATABASE_URL` (Marketplace Central DB)
- `MPC_ORACLE_USERNAME` (Oracle internal-read username)
- `MPC_ORACLE_PASSWORD` (Oracle internal-read password)
- `MPC_ORACLE_CONNECT_STRING` (Oracle service/connect string)
- `MPC_ORACLE_LIB_DIR` (optional Oracle client library directory for Windows/macOS runtime)

Legacy Oracle source keys still found in local operator setups:
- `SANKHYA_ORACLE_USER`
- `SANKHYA_ORACLE_PASSWORD`
- `SANKHYA_ORACLE_HOST`
- `SANKHYA_ORACLE_PORT`
- `SANKHYA_ORACLE_SERVICE_NAME`

When only the legacy `SANKHYA_ORACLE_*` keys exist, map them into the new MPC-owned Oracle runtime keys before running live validation.

Legacy VTEX keys:
- `VTEX_APP_KEY`
- `VTEX_APP_TOKEN`
- `VTEX_ACCOUNT`

These are legacy only after ADR-005. They are not required for current Marketplace Central startup, validation, or Mercado Livre operations. Do not add new VTEX-dependent workflows.

## Starting Local Services

Web:

```powershell
npm run dev --workspace=apps/web
```

Server (loads `.env` first):

```powershell
powershell -ExecutionPolicy Bypass -File .claude/start-server.ps1
```

Note:
- `launch.json` may show a generic preview port, but effective server bind comes from `API_PORT` in `.env`.

## Database Access

Use your local `psql` client and connection URLs from `.env`.

PowerShell example (loads `.env` into current shell and strips surrounding quotes when present):

```powershell
function Normalize-EnvValue([string]$value) {
  if ($null -eq $value) { return $null }
  $trimmed = $value.Trim()
  if ($trimmed.Length -ge 2) {
    if (($trimmed.StartsWith("'") -and $trimmed.EndsWith("'")) -or ($trimmed.StartsWith('"') -and $trimmed.EndsWith('"'))) {
      return $trimmed.Substring(1, $trimmed.Length - 2)
    }
  }
  return $trimmed
}

Get-Content .env |
  Where-Object { $_ -match '^[^#]' -and $_ -match '=' } |
  ForEach-Object {
    $parts = $_ -split '=', 2
    [System.Environment]::SetEnvironmentVariable($parts[0].Trim(), (Normalize-EnvValue $parts[1]), 'Process')
  }
```

Connect to Marketplace Central DB:

```powershell
psql "$env:MC_DATABASE_URL"
```

Quick verification queries:

```sql
SELECT tenant_id, account_id, marketplace_code, integration_installation_id
FROM marketplace_accounts
ORDER BY marketplace_code, account_id
LIMIT 20;
```

```sql
SELECT tenant_id, installation_id, provider_code, status
FROM integration_installations
ORDER BY provider_code, installation_id
LIMIT 20;
```

## Marketplace + Integrations Sanity Check

The intended relationship is one-way and nullable:

```sql
marketplace_accounts.integration_installation_id
  -> integration_installations(tenant_id, installation_id)
```

No reverse FK is expected from integrations tables back to marketplace tables.

## Oracle Live Validation Notes

The M-03 Oracle-first milestone was validated live on Windows with:

- Oracle Instant Client on disk and reachable in `PATH`
- a `cgo`-capable Go run using `godror`
- a user-scoped portable `gcc` toolchain
- normalized `.env` loading to avoid quoted-password failures

If the local environment still uses legacy Sankhya env names, map them before running Oracle validation:

```powershell
$env:MPC_ORACLE_USERNAME = $env:SANKHYA_ORACLE_USER
$env:MPC_ORACLE_PASSWORD = $env:SANKHYA_ORACLE_PASSWORD
$env:MPC_ORACLE_CONNECT_STRING = "$($env:SANKHYA_ORACLE_HOST):$($env:SANKHYA_ORACLE_PORT)/$($env:SANKHYA_ORACLE_SERVICE_NAME)"
```

Windows live-validation session setup:

```powershell
$env:PATH = 'C:\Users\leandro.theodoro\AppData\Local\codex-tools\winlibs-mcf-ucrt\bin;C:\oracle\instantclient_23_0;' + $env:PATH
$env:CC = 'C:\Users\leandro.theodoro\AppData\Local\codex-tools\winlibs-mcf-ucrt\bin\gcc.exe'
$env:CXX = 'C:\Users\leandro.theodoro\AppData\Local\codex-tools\winlibs-mcf-ucrt\bin\g++.exe'
$env:CGO_ENABLED = '1'
$env:MPC_ORACLE_LIB_DIR = 'C:/oracle/instantclient_23_0'
$env:MPC_ORACLE_LIVE_TEST = '1'
```

Focused live proof command:

```powershell
cd apps/server_core
$env:GOCACHE = (Join-Path (Get-Location) '.gocache')
go test ./internal/modules/internal_read/adapters/oracle -run TestOracleLiveSmoke -v
```

Expected proof scope:
- product lookup
- sellable stock
- current price
- cost as-of
- sales history
- tax inputs

## Architecture Notice

- `MS_DATABASE_URL` is no longer part of the target Marketplace Central architecture.
- `MS_TENANT_ID` is also legacy as part of that removed path.
- Internal product/stock/price/cost/tax/sales reads now run through Oracle adapters inside `apps/server_core`.
- Oracle live validation depends on a `cgo`-capable Go toolchain and Oracle client libraries for the `godror` driver.
- `MS_DATABASE_URL` and `MS_TENANT_ID` are not part of the target MPC runtime path and must not be reintroduced as integration shortcuts.
