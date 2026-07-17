# Module: Integrations

Layer: technical operations
Path: `apps/server_core/internal/modules/integrations/`
Frontend: removed — route renders `WorkspacePlaceholder` stub (M-02 COR-2, `packages/feature-integrations` deleted; rebuild pending)

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

## Provider Catalog Baseline

Current marketplace providers in the integrations catalog:

- `mercado_livre` - `oauth2`, interactive
- `magalu` - `oauth2`, interactive
- `amazon` - `lwa`, interactive
- `leroy_merlin` - `api_key`, manual
- `madeira_madeira` - `token`, manual
- `shopee` - `shopee_partner`, interactive, available (`execution_mode=available`, rollout `v1`)

Catalog metadata is the source for rollout and UX semantics (`rollout_stage`, `execution_mode`, `unavailable_reason`, baseline fee hints, credential schema).

## Amazon SP-API Auth Notes

Amazon identifiers are easy to mix up. Keep this mapping straight:

| Pattern | Meaning | Used as |
|---------|---------|---------|
| `amzn1.sellerapps.app.*` | SP-API application ID | Consent `application_id` parameter |
| `amzn1.application-oa2-client.*` | LWA client ID | Token exchange `client_id` |
| `amzn1.oa2-cs.*` | LWA client secret | Client secret for token exchange |
| `amzn1.sp.solution.*` | Solution identifier | Not the consent `application_id` |

Draft/testing consent flows should use `version=beta` in the authorize URL.

### Troubleshooting `INTEGRATIONS_AUTH_CONFIGURATION_INVALID`

- Confirm the consent URL uses the SP-API application ID from `amzn1.sellerapps.app.*`
- Confirm token exchange uses the LWA client ID from `amzn1.application-oa2-client.*`
- Confirm the client secret comes from `amzn1.oa2-cs.*`
- Do not substitute the solution identifier (`amzn1.sp.solution.*`) for consent setup
- Recheck the authorize URL parameters before retrying the install flow

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
