# F-02 Plan - Margin Input Model

```yaml
id: M-06-F-02
type: feature-plan
status: drafted
owner: Feature Implementer
parent: M-06
created: 2026-07-09
updated: 2026-07-09
```

## Slice Decision

Build the profitability truth foundation before any final math:

- profitability domain types
- order-item input assembly from `orders` + `internal_read` + `product_links`
- append-only manual adjustments
- inspectable API/SDK contracts

## Work Breakdown

1. Create `profitability` domain and ports

- define component kinds, quality states, input rows, and manual adjustment rows
- define readers for orders, links, and internal facts
- define store ports for assembled inputs and manual adjustments

2. Create Postgres schema and repositories

- add migration `0028_profitability_margin_inputs.sql`
- persist assembled margin input rows
- persist append-only manual adjustments

3. Create module adapters

- orders adapter reading persisted order snapshots
- product-links adapter exposing resolved link truth as needed
- internal-read adapter mapping cost/tax contracts into profitability inputs

4. Create application services

- import/assemble margin inputs from one installation
- manual adjustment create/list services
- keep margin calculation out of scope

5. Wire transport and contracts

- add endpoints for import/list inputs and create/list adjustments
- update OpenAPI and `sdk-runtime` together

6. Validate

- Go tests for complete/missing/manual scenarios
- SDK tests for request paths
- runtime validation against the real imported order from F-01 plus manual adjustment persistence

## Initial File Targets

- `apps/server_core/internal/modules/profitability/domain/*`
- `apps/server_core/internal/modules/profitability/application/*`
- `apps/server_core/internal/modules/profitability/ports/*`
- `apps/server_core/internal/modules/profitability/adapters/postgres/*`
- `apps/server_core/internal/modules/profitability/adapters/orders/*`
- `apps/server_core/internal/modules/profitability/adapters/productlinks/*`
- `apps/server_core/internal/modules/profitability/adapters/internalread/*`
- `apps/server_core/internal/modules/profitability/transport/http_handler.go`
- `apps/server_core/internal/composition/root.go`
- `apps/server_core/migrations/0028_profitability_margin_inputs.sql`
- `contracts/api/marketplace-central.openapi.yaml`
- `packages/sdk-runtime/src/index.ts`

## Verification Commands

Backend:

```powershell
cd apps/server_core
$env:GOCACHE="$PWD\\.gocache"
go test ./internal/modules/profitability/... ./internal/composition -count=1
```

SDK:

```powershell
npm run test --workspace @marketplace-central/sdk-runtime
```

Runtime:

```powershell
Invoke-RestMethod -Method Post -Uri "http://localhost:8080/profitability/margin-inputs/import" -ContentType "application/json" -Body '{"installation_id":"inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98","limit":1}'
```
