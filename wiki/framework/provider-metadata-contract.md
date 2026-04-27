# Provider Metadata Contract

## Purpose

`ProviderDefinition.Metadata` is the contract between backend provider plugins and the provider catalog UX.

The backend owns the truth. The frontend reads metadata through `sdk-runtime` and renders states/actions from that data.

## Source Of Truth

Go:

```text
apps/server_core/internal/modules/integrations/domain/provider_definition.go
apps/server_core/internal/modules/integrations/adapters/<provider>/auth_adapter.go
```

SDK:

```text
packages/sdk-runtime/src/index.ts
IntegrationProviderMetadata
```

Frontend:

```text
packages/feature-marketplaces/src/ProviderCatalogCard.tsx
packages/feature-marketplaces/src/ProviderCatalogPanel.tsx
```

## Stable Metadata Keys

| Key | Type | Required | Meaning |
|---|---|---:|---|
| `country` | string | yes | Primary country or market, currently `BR` |
| `rollout_stage` | string | yes | Product rollout stage: `v1`, `wave_2`, `blocked` |
| `execution_mode` | string | yes | Runtime availability: `available`, `planned`, `blocked` |
| `unavailable_reason` | string | when blocked | Human-readable reason shown by UI |
| `fee_source` | string | yes | Fee schedule source: `api_sync`, `seed`, `manual` |
| `baseline_commission_percent` | number | yes | Catalog-level baseline commission hint |
| `baseline_fixed_fee_amount` | number | yes | Catalog-level fixed fee hint |
| `credential_schema` | array | yes | Manual credential fields, empty for pure interactive providers if not needed |
| `docs_url` | string | recommended | Link to provider or internal docs |

## Credential Schema Shape

```json
[
  { "key": "client_id", "label": "Client ID", "secret": false },
  { "key": "client_secret", "label": "Client Secret", "secret": true }
]
```

Rules:

- `key` is the stable machine name sent to backend metadata/credentials.
- `label` is display text.
- `secret` controls masked input behavior.
- Do not include values in metadata. Metadata describes fields only.

## Execution Mode Semantics

| Value | Meaning | UI Behavior |
|---|---|---|
| `available` | Provider can be installed or configured now | CTA enabled based on install mode/status |
| `planned` | Provider is visible but not yet actionable | CTA disabled or future-state text |
| `blocked` | Provider cannot be connected due external/internal blocker | CTA disabled, show `unavailable_reason` |

`blocked` wins over any installation state in the provider card.

## Rollout Stage Semantics

| Value | Meaning |
|---|---|
| `v1` | Supported in first production version |
| `wave_2` | Planned for later rollout |
| `blocked` | Cannot progress until blocker is resolved |

Rollout stage is product planning. Execution mode is runtime availability. Do not use one as a substitute for the other.

## Fee Source Semantics

| Value | Meaning |
|---|---|
| `api_sync` | Fee schedules are fetched from provider API |
| `seed` | Fee schedules are deterministic curated rows |
| `manual` | Fees require manual configuration |

If `fee_source=api_sync`, there should be a connector syncer or a documented planned state.

## Declared Capabilities

`DeclaredCapabilities` is not inside metadata, but it is part of the same provider catalog contract.

Common capability keys:

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

The frontend shows capability labels. Backend services also gate runtime operations by these keys.

## Example: Available Interactive Provider

```go
Metadata: map[string]any{
    "country": "BR",
    "rollout_stage": "v1",
    "execution_mode": "available",
    "fee_source": "api_sync",
    "baseline_commission_percent": 0.16,
    "baseline_fixed_fee_amount": 0.0,
    "credential_schema": []map[string]any{
        {"key": "client_id", "label": "Client ID", "secret": false},
        {"key": "client_secret", "label": "Client Secret", "secret": true},
    },
}
```

## Example: Blocked Provider

```go
Metadata: map[string]any{
    "country": "BR",
    "rollout_stage": "blocked",
    "execution_mode": "blocked",
    "unavailable_reason": "Provider access is blocked until partner credentials are approved.",
    "fee_source": "seed",
    "baseline_commission_percent": 0.14,
    "baseline_fixed_fee_amount": 0.0,
    "credential_schema": []map[string]any{},
}
```

## Known Contract Cleanup

There is a vocabulary mismatch to clean up before a production freeze:

- Integration provider metadata uses `execution_mode=available|planned|blocked`.
- Some older marketplace definition types/docs use `live|blocked`.

The integration provider contract should be treated as canonical for the provider catalog UX.
