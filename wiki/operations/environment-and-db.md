# Environment and Database Access

This page documents runtime environment keys and practical database access for local development.

## Runtime Sources

- `.env` is loaded by [start-server.ps1](/C:/Users/leandro.theodoro.MN-NTB-LEANDROT/Documents/marketplace-central/.claude/start-server.ps1)
- Preview server definition is in [launch.json](/C:/Users/leandro.theodoro.MN-NTB-LEANDROT/Documents/marketplace-central/.claude/launch.json)

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
- `MS_DATABASE_URL` (MetalShopping DB, catalog source)
- `MS_TENANT_ID`

Legacy VTEX keys:
- `VTEX_APP_KEY`
- `VTEX_APP_TOKEN`
- `VTEX_ACCOUNT`

These are legacy only after ADR-005. Do not add new VTEX-dependent workflows.

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

PowerShell example (loads `.env` into current shell):

```powershell
Get-Content .env |
  Where-Object { $_ -match '^[^#]' -and $_ -match '=' } |
  ForEach-Object {
    $parts = $_ -split '=', 2
    [System.Environment]::SetEnvironmentVariable($parts[0].Trim(), $parts[1].Trim(), 'Process')
  }
```

Connect to Marketplace Central DB:

```powershell
psql "$env:MC_DATABASE_URL"
```

Connect to MetalShopping DB:

```powershell
psql "$env:MS_DATABASE_URL"
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
