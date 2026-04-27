# Shopee API Best Practices

Last verified: 2026-04-27  
Scope: Shopee "API Best Practices" category (29 documents, including BR subcategory content).

## What we covered

Primary tracks documented by Shopee:

- Product lifecycle
  - product creation prep (`209`)
  - creating product (`211`)
  - creating/publishing global product (`213`, `215`)
  - variants (`219`)
  - base info (`221`)
  - stock and price (`223`)
  - guidelines overview (`217`)
- Fulfillment and post-order
  - first mile binding (`225`)
  - order management (`229`)
  - return/refund (`227`)
- Regional and specialist tracks
  - SIP best practices (`261`)
  - BR logistics/invoice/order/auth operational guides (`286`, `292`, `381`, `382`, `383`, `385`, `568`, `697`)
  - direct delivery, instant mart, xpress package-free, livestream, ads, AMS, video (`290`, `643`, `677`, `669`, `277`, `702`, `706`)
  - sensitive market/domain extensions (auto parts compatibility `378`)

## Best-practice contract for MPC

### 1) Product and catalog

- Always resolve category + attribute constraints before item creation.
- Validate required/optional attributes and input data types per category.
- Use staged flow:
  1. media upload
  2. category/attributes
  3. description
  4. price/stock
- Keep global product flow separate from shop product publication flow.

### 2) Variants

- Model tier structure explicitly; variant mutations may require structural migration.
- Treat variant delete/reorder as state transitions, not blind overwrites.
- Keep SKU identity mapping stable across tier updates.

### 3) Price and stock

- Implement dedicated price and stock update paths; do not conflate both updates.
- For global products, use global-level update paths where required.
- Use reconciliation jobs to detect divergence between intended state and Shopee state.

### 4) Orders and fulfillment

- Build around explicit order status and package fulfillment status machines.
- Support split and cancel split transitions safely.
- For logistics, support mode-specific handling (pickup/drop-off/self-deliver/express variations).
- Upload invoice docs and logistics transitions as retry-safe, idempotent operations.

### 5) Returns and refunds

- Encode return status flow separately from order status flow.
- Handle dispute and non-dispute branches distinctly.
- Persist return reason/solution metadata for auditability.

### 6) Push + pull consistency

- Use push to reduce latency, but reconcile critical states with read APIs.
- Deduplicate events by stable identifiers and replay windows.
- Keep dead-letter handling for malformed/out-of-order notifications.

### 7) Authorization and security

- Authorization and refresh must be fully automated and observable.
- Sensitive-data integrations must pass security prerequisites before production scope usage.
- For BR and logistics-specialized tracks, keep market-specific behavior behind capability flags.

## Recommended module ownership in MPC

- `integrations`
  - provider definition
  - install/auth/credential lifecycle
  - token refresh + auth-state transitions
- `connectors`
  - signed Shopee HTTP execution
  - request/response DTO mapping
  - capability handlers (catalog, order, return, logistics, push)
- domain modules (`catalog`, `orders`, `messaging`, `alerts`)
  - canonical business state and SLA policy

## Implementation checklist before enabling Shopee in production

1. Signed auth and refresh adapter tested with deterministic signature fixtures.
2. Idempotent handlers for product, stock/price, order, logistics, and return mutations.
3. Push ingestion pipeline with dedupe + replay safety.
4. Status mapping table for all V2 states used by active capabilities.
5. Market-specific flags documented (especially BR-only APIs and operational requirements).
6. Security controls aligned for sensitive-data access requests.
