# F-03 Validation

```yaml
id: M-06-F-03
type: feature-validation
status: quick_validation_passed
owner: Feature Implementer
parent: M-06
created: 2026-07-09
updated: 2026-07-09
```

## Scope

Validate persisted profit snapshot calculation only.

## Local Contract Validation

- Command:
  - `cd apps/server_core; $env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\.gocache'; go test ./internal/modules/profitability/... ./internal/composition -count=1`
- Result:
  - Passed.
  - profitability application and transport tests passed, including complete, incomplete, and negative-margin scenarios.
- Command:
  - `npm run test --workspace @marketplace-central/sdk-runtime`
- Result:
  - Passed with `35` tests.
  - SDK now covers calculate/list profit snapshot endpoints.

## Runtime Validation

- Environment:
  - Docker backend on `http://localhost:8080`
  - restart command: `docker compose restart backend`
  - backend boot evidence: `applied 1 migration(s)` for `0029_profitability_profit_snapshots.sql`
- Calculation evidence:
  - `POST /profitability/profit-snapshots/calculate`
  - body: `{"installation_id":"inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98","limit":20}`
  - result: `calculated_count=2`
- Persisted readback:
  - `GET /profitability/profit-snapshots?installation_id=inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98&limit=20`
  - returned:
    - item snapshot with `revenue_amount=19.9`, `sale_fee_amount=9.63`, `quality=incomplete`, flags `missing_link`, `missing_cost`, `missing_tax`
    - order snapshot with aggregated `revenue_amount=19.9`, `sale_fee_amount=9.63`, manual `freight_amount=12.5`, `quality=incomplete`, flags `missing_link`, `missing_commission`

## Live Provider Boundary

- Live validation covers:
  - calculation over a real Mercado Livre order previously imported in F-01
  - persisted margin inputs from F-02
  - persisted manual freight adjustment from F-02
  - profit snapshot persistence in local Postgres
- Live validation does not cover:
  - a fully complete live order with resolved link plus Oracle cost/tax inputs in this session
  - profitability UI

## Open Blockers

- None for F-03 backend calculation.

## Postgres Evidence

- Command:
  - `select count(*) as snapshot_count from profitability_profit_snapshots where installation_id = 'inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98';`
  - `select scope, provider_item_id, revenue_amount, sale_fee_amount, freight_amount, quality from profitability_profit_snapshots where installation_id = 'inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98' order by scope, provider_item_id;`
- Result:
  - `snapshot_count=2`
  - item snapshot row persisted with revenue/fee
  - order snapshot row persisted with aggregated revenue/fee and manual freight `12.50`

## Quality Honesty Evidence

- Unknown cost/tax remained absent from the numeric result and explicit in `flags`.
- Manual freight stopped being falsely treated as missing after recalculation.
- Order-level snapshot now aggregates item revenue/fee instead of reporting them as missing.
