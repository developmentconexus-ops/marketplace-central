# Module: Marketplaces

Layer: business configuration
Path: `apps/server_core/internal/modules/marketplaces/`
Frontend: removed — route renders `WorkspacePlaceholder` stub (M-02 COR-2, `packages/feature-marketplaces` deleted; rebuild pending)

## Main Question It Answers

"Which channels do we sell on, and under what pricing rules?"

## What This Module Owns

| Entity | Purpose |
|--------|---------|
| `MarketplaceDefinition` | Code-driven marketplace catalog (code, display name, auth mode, capabilities, credential schema) |
| `Account` | Tenant channel configuration and linkage to a marketplace |
| `Policy` | Pricing and SLA settings used by pricing/simulation flows |
| `FeeSchedule` | Marketplace commission data, usually synced from integration providers |

## Database Focus

```sql
marketplace_accounts
  tenant_id
  account_id
  marketplace_code
  integration_installation_id -- nullable FK to integration_installations

marketplace_pricing_policies
  tenant_id
  policy_id
  account_id -- FK to marketplace_accounts(tenant_id, account_id)
```

Important behavior:
- `integration_installation_id` is optional.
- Account can exist without an active integration installation.
- Pricing can still run with `Policy` + `FeeSchedule` even when integration is disconnected.

## HTTP Routes (Current)

```text
GET  /marketplaces/accounts
POST /marketplaces/accounts
GET  /marketplaces/policies
POST /marketplaces/policies
GET  /marketplaces/definitions
GET  /marketplaces/fee-schedules?marketplace_code=
POST /admin/fee-schedules/seed
```

## Runtime Notes

- Definitions are plugin-based in `modules/marketplaces/registry/*` and auto-registered via `init()`.
- Startup sync seeds definitions and optional stub fee rows.
- Frontend marketplace forms are definition-driven from API payloads.
- Marketplace catalog UX in `feature-marketplaces` is provider-first and reads operational connection state from integrations (`/integrations/providers` + `/integrations/installations`), while marketplace accounts/policies stay focused on pricing setup.
