# F-01 Plan - Order Ingestion

```yaml
id: M-06-F-01
type: feature-plan
status: drafted
owner: Feature Implementer
parent: M-06
created: 2026-07-09
updated: 2026-07-09
```

## Slice Decision

Implement the smallest durable `orders` slice that already respects the final architecture:

- new business module
- manual ingestion entrypoint
- normalized persistence
- link-quality projection
- idempotent upsert

Do not start with UI or profitability.

## Work Breakdown

1. Create `orders` domain and ports

- define order, item, payment, link-quality, and import result types
- define provider-source and repository ports
- define link reader port with only the listing identity and quality data needed

2. Create Postgres schema and repository

- add migration `0027_orders_marketplace_orders.sql`
- create order, item, and payment tables
- implement transactional upsert repository
- enforce tenant scope and deterministic replacement/update of item/payment snapshots

3. Create adapters

- source adapter delegating to `integrations` provider order read path
- product-links adapter mapping existing workflow/candidate truth into orders link-quality states

4. Create application service

- validate input
- import provider snapshots
- normalize order/items/payments
- enrich items with link quality
- persist atomically
- return import summary

5. Wire composition and transport

- register module in `composition/root.go`
- add minimal HTTP endpoint for import
- add minimal HTTP endpoint for list/read if needed for verification
- keep API shape business-owned and Mercado Livre-neutral

6. Add contract and SDK only for the new orders endpoints

- update OpenAPI
- update `packages/sdk-runtime`
- avoid frontend work unless needed for validation

7. Validate

- unit tests for normalization/idempotency/link-quality behavior
- repository tests for upsert semantics
- transport tests for method/input handling
- real runtime import against the active Mercado Livre installation when possible

## Initial File Targets

- `apps/server_core/internal/modules/orders/domain/*`
- `apps/server_core/internal/modules/orders/application/*`
- `apps/server_core/internal/modules/orders/ports/*`
- `apps/server_core/internal/modules/orders/adapters/postgres/*`
- `apps/server_core/internal/modules/orders/adapters/integrations/*`
- `apps/server_core/internal/modules/orders/adapters/productlinks/*`
- `apps/server_core/internal/modules/orders/transport/http_handler.go`
- `apps/server_core/internal/composition/root.go`
- `apps/server_core/migrations/0027_orders_marketplace_orders.sql`
- `contracts/api/marketplace-central.openapi.yaml`
- `packages/sdk-runtime/src/index.ts`

## Verification Commands

Backend targeted:

```powershell
cd apps/server_core
$env:GOCACHE="$PWD\\.gocache"
go test ./internal/modules/orders/... ./internal/composition -count=1
```

SDK targeted:

```powershell
npm run test --workspace @marketplace-central/sdk-runtime
```

Runtime target:

```powershell
Invoke-RestMethod -Method Post -Uri "http://localhost:8080/orders/import" -ContentType "application/json" -Body '{"installation_id":"inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98","limit":1}'
```

## Risks To Watch

- duplicate child rows on re-import if item/payment replacement is not transactional
- accidentally coupling `orders` directly to `product_links` persistence
- over-storing provider payloads or buyer data
- letting profitability concerns leak into this first ingestion slice
