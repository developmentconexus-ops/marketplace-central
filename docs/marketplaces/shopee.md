# Shopee - API Integration Reference

**Channel ID:** `shopee`
**Current auth strategy:** `api_key` placeholder in MPC, blocked
**Target auth strategy:** signed partner auth (`shopee_partner` or `signed_partner`)
**Rollout stage:** `blocked` / **Execution mode:** `blocked`
**Status: BLOCKED** - pending partner credentials and signed auth implementation

---

## Current Status

Shopee is visible in the integration provider catalog but cannot be connected yet. This is intentional.

Use the framework notes before implementing:

- [Marketplace Integration Framework](../../wiki/framework/marketplace-integration-framework.md)
- [Adding a Marketplace Provider](../../wiki/framework/adding-a-marketplace-provider.md)
- [Shopee Fit Analysis](../../wiki/framework/shopee-fit-analysis.md)

| Capability | Status | Reason |
|---|---|---|
| publish | blocked | Connector not implemented |
| priceSync | seeded | Static fee seeder exists; API sync not implemented |
| stockSync | blocked | Connector not implemented |
| orders | blocked | Connector not implemented |
| messages | blocked | Connector not implemented |
| questions | blocked | Connector not implemented |
| freightQuotes | blocked | Connector not implemented |
| webhooks | blocked | Webhook behavior not confirmed |
| sandbox | blocked | Partner credentials needed |

---

## Platform

Shopee uses Shopee Open Platform for third-party seller integrations. Access requires:

1. Applying as a developer or partner through Shopee Open Platform
2. Receiving a `partner_id` and `partner_key` from Shopee
3. Completing seller/shop authorization
4. Exchanging the returned code for access and refresh tokens
5. Signing subsequent API calls

Official reference used for the current framework decision:

```text
https://cdngarenanow-a.akamaihd.net/shopee/seller/seller_cms/c575929f948611337e1249564c2b8ff6/%5BTW%5D%5BOpen%20API%5DAPI%20v1_v2%E6%8E%88%E6%AC%8A%E6%96%B9%E6%B3%95%20%282020_09%29_newnew.pdf
```

## Auth Mechanism

Shopee Open Platform v2 is an interactive signed partner flow, not plain OAuth2.

Authorize:

```text
GET /api/v2/shop/auth_partner
  ?partner_id={partner_id}
  &redirect={redirect_url}
  &timestamp={unix_timestamp}
  &sign={signature}
```

Callback:

```text
{redirect_url}?code={code}&shop_id={shop_id}
```

Token exchange:

```text
POST /api/v2/auth/token/get
```

Refresh:

```text
POST /api/v2/auth/access_token/get
```

API requests use signed parameters. The common v2 signing base is:

```text
partner_id + api_path + timestamp + access_token + shop_id
```

Example shape:

```text
GET /api/v2/item/get_item_list
  ?partner_id={SHOPEE_PARTNER_ID}
  &timestamp={unix_timestamp}
  &access_token={SELLER_ACCESS_TOKEN}
  &shop_id={SELLER_SHOP_ID}
  &sign={hmac_signature}
```

The Brazil account may still require partner approval, sandbox/UAT configuration, and app-specific endpoint confirmation before implementation.

## Likely API Capabilities

These capabilities are candidates only. Add them to `DeclaredCapabilities` only when connector execution exists or the capability is intentionally documented as planned.

| Operation | Likely endpoint | Notes |
|---|---|---|
| Create product | `POST /api/v2/product/add_item` | With title, description, images, price, stock |
| Update stock | `POST /api/v2/product/update_stock` | By item/model identifiers |
| Update price | `POST /api/v2/product/update_price` | By item/model identifiers |
| Get orders | `GET /api/v2/order/get_order_list` | Order status filter |
| Order detail | `GET /api/v2/order/get_order_detail` | Full order and items |
| Ship order | `POST /api/v2/logistics/ship_order` | Mark as shipped / logistics flow |
| Messages | provider-specific | Confirm before implementing |
| Webhooks | portal/config driven | Confirm before implementing |

---

## Readiness Checklist

Complete all items before promoting Shopee from blocked to available:

- [ ] Apply for Shopee Open Platform Brazil developer/partner access
- [ ] Receive `partner_id` and `partner_key`
- [ ] Confirm production and sandbox/UAT base URLs
- [ ] Confirm callback URL registration requirements
- [ ] Confirm app permissions and seller tier limits
- [ ] Validate which capabilities are available to our seller tier
- [ ] Confirm webhook event types and registration mechanism
- [ ] Implement signed authorize URL generation in Go
- [ ] Implement token exchange and refresh in Go
- [ ] Test signatures with deterministic timestamp fixtures
- [ ] Complete seller shop authorization flow

---

## When Unblocked: Implementation Steps

Implement through the framework:

1. Add `shopee_partner` or `signed_partner` auth strategy if accepted.
2. Update `apps/server_core/internal/modules/integrations/adapters/shopee/auth_adapter.go`.
3. Read env vars:
   - `MPC_PROVIDER_SHOPEE_PARTNER_ID`
   - `MPC_PROVIDER_SHOPEE_PARTNER_KEY`
   - `MPC_PROVIDER_SHOPEE_BASE_URL`
4. Implement `StartAuthorize`, `ExchangeCallback`, and `Refresh`.
5. Store `shop_id` as provider account identity and encrypted credential extra.
6. Keep `fee_source=seed` until API fee sync is confirmed.
7. Update OpenAPI, SDK, frontend auth labels, and framework docs.
8. Run backend tests, frontend tests, and browser validation on `/marketplaces`.

---

## Current Decision

Shopee fits the Marketplace Central framework. It should remain blocked until the signed partner auth adapter is implemented and validated.
