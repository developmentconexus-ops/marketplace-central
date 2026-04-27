# Shopee Capability Matrix

Last verified: 2026-04-27

| Capability | Shopee area | Main reference docs | MPC owner | Notes |
|---|---|---|---|---|
| `auth_installation` | Partner authorization + token lifecycle | `20`, `16`, `385` | `integrations` | Signed URL auth, code exchange, access/refresh token rotation |
| `catalog_create_update` | Product/global product/variants/base info | `209`, `211`, `213`, `215`, `219`, `221`, `217` | `connectors` + `catalog` | Category/attribute-first flow; global and shop flows differ |
| `stock_price_management` | Price/stock updates | `223` | `connectors` + `catalog`/`pricing` | Keep independent price and stock update paths |
| `orders_sync` | Order list/detail/status/splits | `229`, `383` | `connectors` + `orders` | Explicit order + package fulfillment state mapping required |
| `returns_refunds` | Return/refund lifecycle and disputes | `227` | `connectors` + `orders` | Separate return state machine from order state machine |
| `logistics_operations` | First mile, shipping, tracking, invoice, quotation | `225`, `286`, `292`, `382`, `568`, `697`, `290`, `677` | `connectors` + `orders` | Strong BR specialization; capability flags by market are required |
| `push_ingestion` | Push subscriptions and event ingestion | `18` | `connectors` + `messaging`/`orders` | Push must be deduplicated and reconciled with pull APIs |
| `sandbox_validation` | Test account and sandbox procedures | `644`, `643` | `integrations` + `connectors` | Use pre-production evidence gates before rollout |
| `advanced_channels` | Ads, livestream, AMS, video, auto parts | `277`, `669`, `702`, `706`, `378` | `connectors` | Optional/phase-based capabilities, no implicit enablement |
| `sensitive_data_access` | Sensitive scope eligibility and controls | `718` | `integrations` + security operations | Requires additional approval/security evidence |

## Notes on IDs

- ID references map to `https://open.shopee.com/developer-guide/<ID>`.
- Full source map: `wiki/framework/vendors/shopee/sources.md`.
