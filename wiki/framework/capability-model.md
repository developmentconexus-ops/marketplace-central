# Capability Model

## Purpose

Capabilities describe what a provider can do and whether each operation is available at catalog time and runtime.

Marketplace Central uses two related capability models:

1. Static business/catalog capability profile in `marketplaces`.
2. Runtime provider capability state in `integrations`.

Do not merge them. They answer different questions.

## Static Marketplace Capabilities

Path:

```text
apps/server_core/internal/modules/marketplaces/domain/marketplace_def.go
```

`MarketplaceDefinition.CapabilityProfile` is catalog-oriented. It describes what the marketplace channel is expected to support from a product/business perspective.

Statuses:

```text
supported
partial
planned
blocked
```

Used by:

- marketplace definitions
- marketplace setup UX
- planning and visibility

Not used for:

- auth health
- operation run state
- credential validity

## Runtime Integration Capabilities

Paths:

```text
apps/server_core/internal/modules/integrations/domain/capability_state.go
apps/server_core/internal/modules/integrations/application/capability_service.go
```

`CapabilityState` is installation-specific runtime state.

Statuses:

```text
enabled
degraded
disabled
requires_reauth
unsupported
```

Used by:

- runtime health
- operation results
- auth/reauth needs
- provider-specific failures

## Declared Capabilities

Path:

```text
apps/server_core/internal/modules/integrations/domain/provider_definition.go
```

`ProviderDefinition.DeclaredCapabilities` declares what the provider plugin says it can support.

Common keys:

```text
catalog_publish
pricing_fee_sync
inventory_sync
order_read
message_read
message_reply
shipment_tracking
webhook_receive
```

Runtime services may gate operations by these strings. For example, fee sync requires `pricing_fee_sync`.

## Fee Sync Capability Flow

```text
ProviderDefinition.DeclaredCapabilities includes pricing_fee_sync
  -> tenant installation exists
  -> installation status is connected or degraded
  -> FeeSyncService.StartSync creates operation run
  -> MarketplaceExecutor executes registered syncer or fallback seed
  -> capability state is updated
```

Important paths:

```text
apps/server_core/internal/modules/integrations/application/fee_sync_service.go
apps/server_core/internal/modules/integrations/adapters/feesync/marketplace_executor.go
apps/server_core/internal/modules/marketplaces/ports/fee_syncer.go
```

## Connector Capability Execution

`connectors` implements external API behavior. A connector should:

- call provider APIs
- sign requests
- map provider payloads
- validate provider responses
- return domain-facing results through ports

A connector should not:

- own tenant marketplace accounts
- own pricing policies
- own credential lifecycle
- perform React/UI logic

## Capability Availability Rules

Use this rule when deciding a provider state:

| Situation | Metadata / Capability State |
|---|---|
| Provider docs/credentials unavailable | `execution_mode=blocked`, static capabilities `blocked` |
| Provider planned but not implemented | `execution_mode=planned`, static capabilities `planned` |
| Provider install/auth works, no runtime capability yet | auth available; capability omitted or `planned` |
| Runtime capability implemented and tested | declare capability and register executor |
| Runtime capability failed for one tenant | update `CapabilityState` to `degraded` or `requires_reauth` |

## Shopee Example

Current Shopee state:

- business marketplace exists
- integration provider exists
- execution mode is blocked
- capabilities are visible but not connectable
- fee seeder exists for static schedules
- real signed auth client is not implemented yet

Future Shopee state should not be marked available until signed partner authorization, token exchange, and refresh are implemented and tested.
