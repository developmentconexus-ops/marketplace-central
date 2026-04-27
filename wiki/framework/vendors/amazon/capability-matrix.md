# Amazon SP-API Capability Matrix

Last verified: 2026-04-27

| Capability | Amazon area | Main references | MPC owner | Notes |
|---|---|---|---|---|
| `auth_installation` | LWA OAuth + website/appstore authorization | `authorizing-selling-partner-api-applications`, `website-authorization-workflow` | `integrations` | Includes callback state validation and token lifecycle |
| `request_signing` | AWS SigV4 | `connecting-to-the-selling-partner-api` | `connectors` | Required on all SP-API calls |
| `restricted_data_access` | Tokens API (RDT) | `authorization-with-the-restricted-data-token` | `integrations` + `connectors` | Isolate PII-capable paths |
| `catalog_listings` | Catalog Items + Listings APIs | `catalog-items-api`, `manage-product-listings-guide` | `connectors` + `catalog` | Product type and listing issue handling required |
| `stock_price_management` | Listings/Pricing/Feeds | `product-pricing-api`, `feeds-api-best-practices` | `connectors` + `catalog`/`pricing` | Throttle-aware mutation orchestration |
| `orders_sync` | Orders APIs | `orders-api`, `orders-api-rate-limits` | `connectors` + `orders` | Strong pagination and status mapping needed |
| `logistics_operations` | Fulfillment + shipping APIs | `merchant-fulfillment-api`, `easy-ship-api`, `fulfillment-outbound-api` | `connectors` + `orders` | Capability flags by fulfillment mode/market |
| `returns_refunds` | Returns and post-order finance docs | `report-type-values-returns`, `invoices-api` | `connectors` + `orders` | Market and role dependent |
| `notifications_ingestion` | Notifications API (SQS/EventBridge) | `notifications-api`, `set-up-notifications-with-amazon-sqs` | `connectors` + `messaging`/`orders` | Dedupe and replay-safe consumers required |
| `reports_analytics` | Reports/Data Kiosk | `reports-api`, `data-kiosk-workflow-guide` | `connectors` + analytics consumers | Batch/offline extraction use cases |
| `finance_payments` | Finances/Invoices/transactions | `list-transactions-by-order`, `invoices-api` | `connectors` + `orders`/finance consumers | Strong auditability required |
| `vendor_mode` | Vendor APIs | `register-as-a-private-vendor`, `vendor-df-workflow-guide` | `connectors` | Separate from seller mode; phase-gated |

## Runtime rollout note (current)

- Amazon provider is operational in live execution mode.
- Production usage remains capability-gated and should be enabled incrementally:
  1. auth + signed read calls
  2. listings/stock/price
  3. orders/fulfillment
  4. notifications and restricted-data tracks

