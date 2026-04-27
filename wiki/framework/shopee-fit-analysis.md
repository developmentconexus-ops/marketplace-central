# Shopee Fit Analysis

## Question

Can Shopee fit the Marketplace Central marketplace integration framework?

Yes. Shopee fits the framework well, but it should not be modeled as plain OAuth2.

## Current MPC State

Shopee currently exists in three places:

| Layer | Path | Current Role |
|---|---|---|
| marketplace business definition | `apps/server_core/internal/modules/marketplaces/registry/shopee.go` | visible but blocked marketplace channel |
| integration provider | `apps/server_core/internal/modules/integrations/adapters/shopee/auth_adapter.go` | manual API-key placeholder, blocked metadata |
| connector fee seed | `apps/server_core/internal/modules/connectors/adapters/shopee/fee_seed.go` | deterministic fee schedule seeder |

Current provider metadata says:

```text
rollout_stage=blocked
execution_mode=blocked
auth_strategy=api_key
install_mode=manual
```

That was a safe placeholder while Shopee partner access and official flow details were unconfirmed.

## Shopee Open Platform Fit

Shopee Open Platform v2 uses a signed partner authorization model:

```text
GET /api/v2/shop/auth_partner
POST /api/v2/auth/token/get
POST /api/v2/auth/access_token/get
```

The authorize URL uses:

```text
partner_id
redirect
timestamp
sign
```

The callback returns:

```text
code
shop_id
```

The token endpoint returns:

```text
access_token
refresh_token
expire_in
```

API calls are signed using partner credentials and request parameters.

## Framework Mapping

| Shopee Concept | MPC Framework Location |
|---|---|
| Partner ID / Partner Key | provider adapter env config, never frontend |
| Seller shop authorization URL | `StartAuthorize` in Shopee auth adapter |
| Callback `code` + `shop_id` | `ExchangeCallback` in Shopee auth adapter |
| Access token / refresh token | encrypted integration credential |
| Shop ID | provider account ID / credential extra |
| Token refresh | `Refresh` in Shopee auth adapter |
| Signed API requests | connector adapter helper, not React |
| Fee schedules | existing connector fee seed, future API sync if supported |
| Orders/messages/catalog | future connector capabilities |

## Recommended Auth Strategy

Do not use `oauth2` for Shopee unless Shopee exposes a standard OAuth2 flow for the specific Brazil app.

Recommended options:

1. `shopee_partner`
   - Most explicit.
   - Best if only Shopee uses this pattern.

2. `signed_partner`
   - More reusable.
   - Better if other marketplaces use the same partner_id/partner_key signed auth pattern.

Either choice requires updates to:

```text
apps/server_core/internal/modules/integrations/domain/lifecycle.go
contracts/api/marketplace-central.openapi.yaml
packages/sdk-runtime/src/index.ts
packages/feature-marketplaces
wiki/framework/auth-strategies.md
```

## Recommended Implementation Shape

Update:

```text
apps/server_core/internal/modules/integrations/adapters/shopee/auth_adapter.go
```

Add config:

```text
MPC_PROVIDER_SHOPEE_PARTNER_ID
MPC_PROVIDER_SHOPEE_PARTNER_KEY
MPC_PROVIDER_SHOPEE_BASE_URL
```

Adapter behavior:

```text
StartAuthorize:
  validate partner config
  create timestamp
  sign partner_id + path + timestamp
  build /api/v2/shop/auth_partner URL

ExchangeCallback:
  require code
  require shop_id from callback/provider account input
  sign token endpoint request
  POST /api/v2/auth/token/get
  return access_token, refresh_token, shop_id

Refresh:
  require refresh_token
  require stored shop_id
  sign refresh endpoint request
  POST /api/v2/auth/access_token/get
  rotate credential
```

Credential payload should include:

```text
access_token
refresh_token
provider_account_id = shop_id
shop_id
```

## Rollout Metadata After Implementation

When signed auth is implemented but operational APIs are still limited:

```text
rollout_stage=v1
execution_mode=available
auth_strategy=shopee_partner or signed_partner
install_mode=interactive
fee_source=seed
```

Capabilities should only declare what is true:

```text
pricing_fee_sync
```

Add `order_read`, `message_read`, `catalog_publish`, or `inventory_sync` only when connector execution exists or the capability is explicitly planned/blocked in docs.

## Why Shopee Fits

Shopee needs:

- provider-level auth
- tenant-specific installation
- encrypted credentials
- refresh lifecycle
- signed connector calls
- capability-specific rollout

Those are exactly the responsibilities of `integrations` plus `connectors`.

It should not require a new module, a frontend special case, or a marketplace-specific database design.

## Current Blockers

- Need confirmed Shopee Brazil partner credentials.
- Need official app environment details: production vs sandbox/UAT base URL.
- Need callback shape validation in our backend route.
- Need signed request tests using deterministic timestamps.
- Need frontend label support for the chosen auth strategy.

## Decision

Shopee should be promoted from blocked placeholder to interactive provider only after a first-class signed partner auth adapter exists and passes backend tests.

Until then, keeping Shopee visible but blocked is the correct professional SaaS behavior.
