# Integration Extensibility

## Goal

Adding a provider should be predictable, with provider-owned packages and minimal central wiring.

## Catalog Model

- `ProviderDefinition` is the catalog source of truth (what can be installed).
- `Installation` is the tenant-specific connection (what is installed).
- Capability keys define what each provider can do (`pricing_fee_sync`, `order_read`, `message_read`, and similar).

This supports a general integration catalog, not only marketplace channels.

## Current Plugin Framework

| Area | Current Behavior |
|------|------------------|
| Provider definitions | Registered by provider packages via `init()` hooks |
| Auth adapters | Registered by provider packages via auth factories |
| Fee syncers | Registered by connector packages via syncer factories |
| Composition root | Consumes provider/auth/sync registries, no per-provider wiring lists |
| Frontend provider selector | Data-driven from `/integrations/providers` |

## Adding a New Provider Now

Typical new files:

```text
apps/server_core/internal/modules/integrations/adapters/<provider>/auth_adapter.go
apps/server_core/internal/modules/connectors/adapters/<provider>/fee_sync.go   (optional)
```

Expected registration ownership:

- Provider package registers:
  - definition metadata
  - auth adapter factory (reads provider env keys)
- Connector package registers:
  - optional fee syncer

After that, provider definition seeding and runtime adapter assembly happen automatically through registries.
