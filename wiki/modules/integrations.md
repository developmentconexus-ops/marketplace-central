# Module: Integrations

Layer: technical operations
Path: `apps/server_core/internal/modules/integrations/`
Frontend: `packages/feature-integrations/`

## Main Question It Answers

"Is auth healthy? Are credentials valid? Are operations running?"

## What This Module Owns

| Entity | Purpose |
|--------|---------|
| `ProviderDefinition` | Provider metadata (family, auth strategy, capabilities) |
| `Installation` | Tenant-specific provider installation |
| `Credential` | Encrypted token/API key material and lifecycle metadata |
| `CapabilityState` | Per-capability operational state (enabled/degraded/disabled) |
| `OperationRun` | Audit trail for operations such as `fee_sync` |
| `AuthSession` + `OAuthState` | OAuth flow and refresh lifecycle state |

## Plugin Framework

- Provider packages register their own `ProviderDefinition` and auth factory.
- Connector packages register optional fee syncers.
- Composition root reads registries to build runtime dependencies.
- Providers appear in catalog (`/integrations/providers`) even when tenant has no installation yet.

## Database Focus

```sql
integration_installations
  tenant_id
  installation_id
  provider_code
  status

integration_credentials
  tenant_id
  installation_id
  credential_id
  is_active

integration_operation_runs
  tenant_id
  run_id
  installation_id
  operation_type
  status
```

Important behavior:
- No FK from `integration_installations` back to marketplace tables.
- This module is operationally independent from Marketplaces at schema level.
- Integration catalog availability is code-defined and seeded; tenant state lives in installations.

## HTTP Routes (Current)

```text
GET  /integrations/providers
GET  /integrations/installations
POST /integrations/installations
GET  /integrations/installations/:id/auth-status
POST /integrations/installations/:id/authorize
POST /integrations/installations/:id/disconnect
POST /integrations/installations/:id/sync-fees
```

## Fee Sync Flow

1. Client calls `POST /integrations/installations/:id/sync-fees`
2. `FeeSyncService.StartSync` loads installation + provider
3. Provider mapping resolves marketplace code
4. Adapter fetches fee data from external API
5. Marketplaces `FeeSchedule` rows are upserted
6. `OperationRun` records success/failure for audit

No cross-module SQL join is required.
