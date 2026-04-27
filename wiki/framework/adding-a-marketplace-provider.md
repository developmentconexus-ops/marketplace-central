# Adding a Marketplace Provider

## Purpose

This is the implementation runbook for adding a marketplace provider to Marketplace Central.

Use it as the checklist for LLM agents and human engineers. If a provider needs a special flow, document the special flow in `docs/marketplaces/<provider>.md` and keep this runbook as the common path.

## Required Decisions

Before coding, answer:

| Question | Example |
|---|---|
| Provider code? | `shopee` |
| Business marketplace code? | Usually same as provider code |
| Auth strategy? | `oauth2`, `lwa`, `api_key`, `token`, provider-specific |
| Install mode? | `interactive`, `manual`, `hybrid` |
| Rollout state? | `available`, `planned`, `blocked` |
| Capabilities? | `pricing_fee_sync`, `order_read`, `message_read` |
| Fee source? | `api_sync`, `seed`, `manual` |
| Needs connector execution now? | yes/no |

## Step 1: Add or Update Marketplace Business Definition

Path:

```text
apps/server_core/internal/modules/marketplaces/registry/<provider>.go
```

Implement or update `MarketplacePlugin`:

```go
type MarketplacePlugin interface {
    Code() string
    Definition() domain.MarketplaceDefinition
    SeedFees(ctx context.Context, pool *pgxpool.Pool) error
    NewConnector(credentials map[string]string) (MarketplaceConnector, error)
}
```

Use this layer for:

- display name
- business marketplace code
- static capability profile for marketplace catalog
- pricing fee source
- policy/account compatibility
- optional seed fees for providers without a dedicated syncer

Do not use this layer for:

- OAuth token exchange
- credential encryption
- external API request signing
- operation runs
- runtime auth health

## Step 2: Add or Update Integration Provider Adapter

Path:

```text
apps/server_core/internal/modules/integrations/adapters/<provider>/auth_adapter.go
```

Register the provider definition in `init()`:

```go
func init() {
    integrationsproviders.RegisterDefinition(domain.ProviderDefinition{
        ProviderCode: "<provider>",
        TenantID:     "system",
        Family:       domain.IntegrationFamilyMarketplace,
        DisplayName:  "<Display Name>",
        AuthStrategy: domain.AuthStrategyOAuth2,
        InstallMode:  domain.InstallModeInteractive,
        Metadata:     map[string]any{...},
        DeclaredCapabilities: []string{...},
        IsActive: true,
    })

    integrationsproviders.RegisterAuthFactory(func() application.MarketplaceAuthAdapter {
        return NewAdapter(Config{...})
    })
}
```

Implement `application.MarketplaceAuthAdapter`:

```go
ProviderCode() string
StartAuthorize(ctx, input)
ExchangeCallback(ctx, input)
VerifyAPIKey(ctx, input)
Refresh(ctx, input)
```

Manual providers still implement this interface. Unsupported methods return `domain.ErrNotSupported`.

## Step 3: Add Auth Strategy Only When Needed

Current auth strategies live in:

```text
apps/server_core/internal/modules/integrations/domain/lifecycle.go
```

Current values:

```text
oauth2
lwa
api_key
token
none
unknown
```

Add a new strategy only when existing values would misrepresent the provider contract. Shopee is an example candidate because its Open Platform v2 flow is signed partner auth, not plain OAuth2.

If adding a strategy, update:

- Go domain enum
- OpenAPI schema
- `packages/sdk-runtime/src/index.ts`
- frontend label rendering/tests
- framework docs

## Step 4: Add Optional Fee Syncer

Path:

```text
apps/server_core/internal/modules/connectors/adapters/<provider>/fee_sync.go
```

Implement:

```go
type FeeScheduleSyncer interface {
    MarketplaceCode() string
    Sync(ctx context.Context, repo FeeScheduleRepository) (int, error)
}
```

Register:

```go
func init() {
    integrationsproviders.RegisterFeeSyncerFactory(func() ports.FeeScheduleSyncer {
        return NewFeeSyncer()
    })
}
```

Rules:

- API-backed fee sync should call provider API and upsert schedules.
- Static fee sync may seed deterministic rows.
- `pricing_fee_sync` must be declared in `ProviderDefinition.DeclaredCapabilities` before runtime fee sync is allowed.
- Fee sync is enabled only when installation status is `connected` or `degraded`.

## Step 5: Activate Registration In Composition

Path:

```text
apps/server_core/internal/composition/root.go
```

Add side-effect imports:

```go
_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/<provider>"
_ "marketplace-central/apps/server_core/internal/modules/connectors/adapters/<provider>"
```

Side-effect imports are mandatory. A provider package that is not imported by composition will not register at runtime.

## Step 6: Update Provider Metadata

Required metadata keys:

```text
country
rollout_stage
execution_mode
fee_source
baseline_commission_percent
baseline_fixed_fee_amount
credential_schema
```

Required when blocked:

```text
unavailable_reason
```

Recommended:

```text
docs_url
```

See [Provider Metadata Contract](provider-metadata-contract.md).

## Step 7: Update Contracts and SDK When Shape Changes

Update these only when adding new enum values, routes, or stable metadata shape:

```text
contracts/api/marketplace-central.openapi.yaml
packages/sdk-runtime/src/index.ts
packages/sdk-runtime/src/index.test.ts
```

Frontend feature packages must keep using `sdk-runtime`; no direct backend calls.

## Step 8: Update Provider Docs

Create or update:

```text
docs/marketplaces/<provider>.md
```

Include:

- provider code
- auth strategy
- install mode
- env variables
- official docs links
- sandbox/production notes
- supported capabilities
- blocked/planned/available status
- API endpoints used by MPC
- validation evidence

## Step 9: Tests and Verification

Backend:

```powershell
$env:GOCACHE=".gocache"; go test ./internal/modules/integrations/... ./internal/modules/marketplaces/... ./internal/composition/...
```

Frontend:

```powershell
cmd /c npm exec --workspace @marketplace-central/web vitest run packages/sdk-runtime/src/index.test.ts packages/feature-marketplaces/src/MarketplaceSettingsPage.test.tsx
```

Browser validation when UX behavior changes:

```text
/marketplaces
```

Confirm:

- provider appears in grid
- auth label is correct
- blocked/planned/available state is correct
- create installation action works when available
- fee sync is disabled until connected/degraded
- manual credential panel appears only for manual/hybrid providers

## Done Gate

A provider addition is complete only when:

- provider definition self-registers
- auth factory self-registers or intentionally returns unsupported operations
- composition imports the package for registration
- metadata contract is complete
- marketplace business definition is present when it is a marketplace channel
- connector syncer exists for any declared executable capability
- OpenAPI/SDK/frontend are updated when contract changes
- impacted tests pass
- provider docs and wiki are updated
