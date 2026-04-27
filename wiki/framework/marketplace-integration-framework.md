# Marketplace Integration Framework

## Purpose

The marketplace framework lets Marketplace Central add providers as plugins instead of one-off code paths.

The target shape is:

```text
provider package registers catalog/auth
connector package registers capability execution
composition imports packages for registration only
frontend renders provider state from SDK data
```

This keeps Marketplace Central ready for many marketplaces without turning the UI, composition root, or business modules into provider-specific control flow.

## Module Ownership

| Module | Owns | Does Not Own |
|---|---|---|
| `marketplaces` | business account setup, pricing policies, fee schedules, static marketplace definitions | OAuth, provider credentials, external API calls |
| `integrations` | provider catalog, installations, credentials, auth sessions, capability state, operation runs | pricing policy meaning, seller business rules |
| `connectors` | external API adapters and capability execution | tenant business state, marketplace account ownership |
| `packages/sdk-runtime` | typed backend client methods | business logic |
| `packages/feature-marketplaces` | provider catalog UX and marketplace setup UX | direct fetch calls, commission calculations |

## Runtime Flow

```text
web
  -> sdk-runtime
  -> server_core transport
  -> application services
  -> postgres

integrations provider catalog:
  provider package init()
  -> providers.RegisterDefinition(...)
  -> composition/root.go side-effect import
  -> ProviderService.SeedProviderDefinitions(...)
  -> GET /integrations/providers

auth lifecycle:
  create installation
  -> start authorize or submit credentials
  -> rotate encrypted credential
  -> upsert auth session
  -> mark installation connected/healthy

fee sync:
  installation + provider definition
  -> FeeSyncService.StartSync
  -> MarketplaceExecutor
  -> connector FeeScheduleSyncer
  -> marketplace fee schedules
  -> integration operation run + capability state
```

## Core Code Paths

| Concern | Path / Symbol |
|---|---|
| Composition root | `apps/server_core/internal/composition/root.go` |
| Provider registry | `apps/server_core/internal/modules/integrations/adapters/providers/registry.go` |
| Provider definition | `apps/server_core/internal/modules/integrations/domain/provider_definition.go` |
| Auth lifecycle | `apps/server_core/internal/modules/integrations/application/auth_flow_service.go` |
| Installation lifecycle enums | `apps/server_core/internal/modules/integrations/domain/lifecycle.go` |
| Marketplace plugin interface | `apps/server_core/internal/modules/marketplaces/registry/plugin.go` |
| Marketplace registry | `apps/server_core/internal/modules/marketplaces/registry/registry.go` |
| Fee sync executor | `apps/server_core/internal/modules/integrations/adapters/feesync/marketplace_executor.go` |
| SDK contract | `packages/sdk-runtime/src/index.ts` |
| Provider UX | `packages/feature-marketplaces/src/MarketplaceSettingsPage.tsx` |
| Provider panel | `packages/feature-marketplaces/src/ProviderCatalogPanel.tsx` |
| Provider card | `packages/feature-marketplaces/src/ProviderCatalogCard.tsx` |

## Registration Model

Provider packages self-register in `init()`:

```go
func init() {
    providers.RegisterDefinition(definition)
    providers.RegisterAuthFactory(factory)
}
```

Connector packages register optional capability executors:

```go
func init() {
    providers.RegisterFeeSyncerFactory(func() ports.FeeScheduleSyncer {
        return NewFeeSyncer()
    })
}
```

Composition activates packages with side-effect imports:

```go
_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/<provider>"
_ "marketplace-central/apps/server_core/internal/modules/connectors/adapters/<provider>"
```

Composition must not manually build provider-specific auth adapters or syncers except through the registry.

## Database Ownership

Global/system catalog:

```text
integration_provider_definitions
```

Tenant-owned integration runtime:

```text
integration_installations
integration_credentials
integration_auth_sessions
integration_oauth_states
integration_capability_states
integration_operation_runs
```

Tenant-owned marketplace business configuration:

```text
marketplace_accounts
marketplace_pricing_policies
marketplace_fee_schedules
```

Tenant-owned queries must scope by `tenant_id`. Provider definitions are system-owned reference rows and are intentionally global.

## Design Guardrails

- Do not put marketplace API calls in React.
- Do not put provider auth in `marketplaces`.
- Do not put pricing policy behavior in `integrations`.
- Do not add direct `fetch()` calls in frontend feature packages.
- Do not create provider-specific branching in composition beyond side-effect imports.
- Do not declare a capability as available unless there is an execution path or an intentional blocked/planned state.
- Do not mark docs-only framework changes complete without proofreading and committing the docs.
