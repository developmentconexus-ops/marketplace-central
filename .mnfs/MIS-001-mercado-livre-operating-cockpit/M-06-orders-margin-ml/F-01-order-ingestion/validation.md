# F-01 Validation

```yaml
id: M-06-F-01
type: feature-validation
status: quick_validation_passed
owner: Feature Implementer
parent: M-06
created: 2026-07-09
updated: 2026-07-09
```

## Scope

Validate order ingestion only.

## Local Contract Validation

- Command:
  - `cd apps/server_core; $env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\.gocache'; go test ./internal/modules/orders/... ./internal/modules/connectors/... ./internal/modules/integrations/transport ./internal/composition -count=1`
- Result:
  - Passed.
  - `orders` application and transport tests passed.
  - existing Mercado Livre connector tests still passed after extending order snapshot mapping with created/closed/updated timestamps, tags, item sale fee, and payment amount fields.
- Command:
  - `npm run test --workspace @marketplace-central/sdk-runtime`
- Result:
  - Passed with `30` tests.
  - SDK now covers `importMarketplaceOrders` and `listMarketplaceOrders`.

## Runtime Validation

- Environment:
  - Docker backend on `http://localhost:8080`
  - restart command: `docker compose restart backend`
  - backend boot evidence: `applied 1 migration(s)` for `0027_orders_marketplace_orders.sql`
- Probe evidence:
  - `GET /integrations/installations/inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98/probes/orders?limit=1`
  - Returned one real Mercado Livre order:
    - `provider_order_id=2000012659424976`
    - `provider_status=cancelled`
    - `provider_item_id=MLB4834373620`
    - `seller_sku=S-18`
    - `sale_fee_amount=9.63`
    - `payment_id=122193559766`
    - `shipping_id=45338454770`
- Import evidence:
  - `POST /orders/import` with body `{"installation_id":"inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98","limit":1}`
  - Returned `imported_count=1`, `skipped_count=0`
  - Returned normalized item with `link_quality=missing`, proving missing link does not block ingestion
- Persistence evidence:
  - `GET /orders?installation_id=inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98&limit=1`
  - Returned the persisted normalized order snapshot with tags, item sale fee, and payment totals

## Live Provider Boundary

- Live validation covers:
  - Mercado Livre order read through the existing runtime capability path
  - orders-module import and persistence against the real local Postgres database
  - repeated import idempotency at header/item/payment row-count level
- Live validation does not cover:
  - profitability calculations
  - cost/tax/freight enrichment
  - UI, because F-01 intentionally ships backend/API/SDK first

## Open Blockers

- None for F-01 initial ingestion slice.

## Idempotency Evidence

- Command:
  - re-run `POST /orders/import` with the same body
  - query Postgres:
    - `select count(*) as orders_count from orders_marketplace_orders where installation_id = 'inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98';`
    - `select count(*) as items_count from orders_marketplace_order_items where installation_id = 'inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98';`
    - `select count(*) as payments_count from orders_marketplace_order_payments where installation_id = 'inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98';`
- Result:
  - second import completed without duplicate persistence
  - `orders_count=1`
  - `items_count=1`
  - `payments_count=1`
