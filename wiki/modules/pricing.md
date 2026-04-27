# Module: Pricing Simulator

Layer: business intelligence
Path: `apps/server_core/internal/modules/pricing/`
Frontend: `packages/feature-simulator/`

## Main Question It Answers

"Can we sell this product in this marketplace under this policy and still make enough margin?"

The pricing module is the seller decision engine. It turns product cost, selling price, marketplace fees, fixed fees, freight, and margin thresholds into clear viability signals before the seller publishes, reprices, or changes policy.

## Product Goal

The simulator should help operators compare many products across many marketplace policies quickly enough to support real pricing decisions.

It should make these answers obvious:

- Which products are healthy, risky, or unviable by marketplace
- Why margin changed: cost, commission, fixed fee, freight, or selling price
- Which price source was used: current seller price, suggested price, or manual override
- Whether freight is reliable or missing because dimensions, Melhor Envio connection, or quote data are unavailable
- Which policies/channels need configuration before results can be trusted

## What This Module Owns

| Entity / Concept | Purpose |
|------------------|---------|
| `Simulation` | Persisted single simulation summary for audit/history |
| `BatchSimulationItem` | Product x policy comparison row returned to the simulator UI |
| Margin calculation | Canonical server-side calculation of margin amount, margin percent, and status |
| Price source selection | Chooses product `price_amount`, `suggested_price`, or explicit override |
| Freight resolution | Uses policy shipping mode and Melhor Envio quote state to decide freight amount/source |
| Fee resolution | Uses marketplace fee schedules unless a policy-level commission override is present |

## What This Module Does Not Own

- Product catalog truth. It reads product data through a `ProductProvider` port backed by `catalog`.
- Marketplace account/policy ownership. It reads policies through a `PolicyProvider` port backed by `marketplaces`.
- Provider auth, credentials, or fee-sync operations. Those belong to `integrations`.
- Direct marketplace API calls. Runtime marketplace IO belongs behind connector/integration ports.
- Frontend-only margin math. React renders results; Go owns the calculation.

## Current User Workflow

The `/simulator` page is a batch comparison workflow:

1. Load catalog products, classifications, taxonomy nodes, marketplace policies, and Melhor Envio status.
2. Operator filters products by search, taxonomy, or classification.
3. Operator selects products manually or via classification pill.
4. Operator enters origin/destination CEP.
5. Operator chooses `my_price` or `suggested_price`.
6. Operator can override a product x policy selling price after results appear.
7. UI calls `POST /pricing/simulations/batch`.
8. Results render inline as a product x marketplace/policy matrix with margin, cost, commission, freight, fixed fee, and status.
9. Operator can filter by health and export CSV.

## Calculation Contract

For each product x policy:

```text
selling_price =
  price_override[product_id::policy_id]
  OR suggested_price when price_source = "suggested_price" and product has suggested_price
  OR product.price_amount

commission_percent =
  policy.commission_override
  OR fee_schedule(marketplace_code, category_id, listing_type).commission_percent
  OR policy.commission_percent

commission_amount = selling_price * commission_percent

margin_amount =
  selling_price
  - product.cost_amount
  - commission_amount
  - policy.fixed_fee_amount
  - freight_amount

margin_percent = margin_amount / selling_price
```

All monetary values use `float64` in Go domain/application code and `numeric(14,2)` in Postgres until an ADR changes the money convention.

## Status Rules

Single simulation:

| Status | Meaning |
|--------|---------|
| `healthy` | Margin is valid and greater than or equal to the requested minimum |
| `warning` | Margin is non-negative but below the requested minimum |
| `critical` | Margin is negative, invalid, or impossible to calculate |

Batch simulation:

| Status | Meaning |
|--------|---------|
| `healthy` | Margin is greater than `20%` and freight is available |
| `warning` | Margin is between `10%` and `20%` and freight is available |
| `critical` | Margin is below `10%` or freight is unavailable |

Current batch thresholds are implementation constants. If they become tenant-configurable, that belongs in marketplace policies or a dedicated pricing rules table, not React.

## Freight Behavior

Policy `shipping_provider` controls freight source:

| Value | Behavior |
|-------|----------|
| `fixed` | Use `policy.default_shipping`, source `fixed` |
| `marketplace` | Use `policy.default_shipping`, source `marketplace` |
| `melhor_envio` | Quote Melhor Envio when connected and product dimensions/weight are present |

Freight source values in batch results:

| Source | Meaning |
|--------|---------|
| `melhor_envio` | Live quote was available and used |
| `fixed` | Fixed policy freight was used |
| `marketplace` | Marketplace-handled/default policy freight was used |
| `no_dimensions` | Product cannot be quoted because dimensions/weight are missing |
| `me_not_connected` | Policy requires Melhor Envio but account is not connected |
| `me_error` | Melhor Envio was connected but quote failed or returned no usable result |

When freight is unavailable, status must be `critical` even if price margin would otherwise look acceptable.

## Fee Schedule Behavior

Fee lookup precedence:

1. `policy.commission_override` when present
2. Marketplace fee schedule lookup by `marketplace_code` + product taxonomy category
3. `policy.commission_percent`

Product taxonomy node is currently used as the fee schedule category proxy. Empty category falls back to `default`.

Fee schedule lookup is intentionally through a pricing-owned port (`FeeScheduleLookup`) and adapter. Pricing must not import marketplaces domain types directly.

## HTTP Routes

```text
GET  /pricing/simulations
POST /pricing/simulations
POST /pricing/simulations/batch
```

`POST /pricing/simulations` is the legacy/single simulation path. It persists a compact summary to `pricing_simulations`.

`POST /pricing/simulations/batch` is the primary simulator UX path. It returns computed comparison rows and does not currently persist each row.

## API Shapes

Single simulation response:

```json
{
  "simulation_id": "sim_001",
  "tenant_id": "tenant_default",
  "product_id": "prod_001",
  "account_id": "acc_001",
  "margin_amount": 26.0,
  "margin_percent": 0.1733,
  "status": "warning"
}
```

Batch result row:

```json
{
  "product_id": "prod_001",
  "policy_id": "pol_001",
  "selling_price": 150.0,
  "cost_amount": 80.0,
  "commission_amount": 24.0,
  "freight_amount": 20.0,
  "fixed_fee_amount": 0.0,
  "margin_amount": 26.0,
  "margin_percent": 0.1733,
  "status": "warning",
  "freight_source": "fixed"
}
```

## Database Focus

```sql
pricing_simulations
  simulation_id
  tenant_id
  product_id
  account_id
  input_snapshot_json
  result_snapshot_json
  created_at

pricing_manual_overrides
  override_id
  tenant_id
  product_id
  account_id
  target_price_amount
  notes
```

Important behavior:

- `pricing_simulations` is tenant-owned and list queries must scope by `tenant_id`.
- Single simulation saves idempotently with `ON CONFLICT (simulation_id) DO NOTHING`.
- Batch simulation currently returns transient results and does not write batch rows.
- `pricing_manual_overrides` exists in schema but is not the current simulator's source of truth for inline overrides.

## Ports And Dependencies

| Port | Backing module / dependency | Purpose |
|------|-----------------------------|---------|
| `Repository` | Postgres adapter | Save/list single simulation summaries |
| `ProductProvider` | Catalog adapter | Load products, costs, prices, taxonomy, dimensions |
| `PolicyProvider` | Marketplaces adapter | Load marketplace pricing policies |
| `FeeScheduleLookup` | Marketplaces fee schedule adapter | Resolve marketplace/category commission rates |
| `FreightQuoter` | Connector/integration adapter | Quote Melhor Envio freight when policy requires it |

Dependency rule:

```text
pricing/application -> pricing/ports -> adapters
pricing must not import net/http, pgx, catalog domain, or marketplaces domain in application code
```

## Frontend Contract

`packages/feature-simulator` should remain a workflow surface, not a calculator.

Frontend responsibilities:

- Load product/policy/classification data through `sdk-runtime`
- Provide filtering, selection, CEP inputs, price source toggle, inline override UX, CSV export
- Render loading, error, empty, running, and result states
- Display server-returned margin breakdown and health state

Frontend must not:

- Recalculate margin or fee logic
- Fetch backend endpoints directly
- Treat local state/exported CSV as source of truth
- Hide freight or fee uncertainty when backend marks a result as critical

## Current Implementation Notes

- Batch simulation is the real v1 workflow.
- Single simulation remains useful for history/audit and simpler API tests.
- Batch results include commission/freight/fixed-fee breakdown, solving an earlier UX gap where margin was visible but not explainable.
- UI uses a matrix layout because operators compare products across marketplace policies repeatedly.
- The page currently depends on Melhor Envio connection status only for operator visibility; the backend remains authoritative for freight status.

## Known Gaps

- Batch results are not persisted as a named run, so there is no durable comparison history.
- Inline price overrides live in frontend request state; they are not saved to `pricing_manual_overrides`.
- Batch health thresholds (`10%`, `20%`) are hard-coded constants.
- Single simulation OpenAPI says cost can be resolved from catalog when omitted, but the current single simulation service expects the request value it receives.
- No dedicated pricing rules engine exists yet for taxes, campaign discounts, ads cost, payment fees, or marketplace-specific promotion mechanics.
- CSV export is frontend-generated from current visible rows, not a backend report artifact.

## Future Direction

For a production-grade SaaS v1+, the simulator should evolve toward:

- Saved simulation runs with run metadata, selected products, selected policies, and snapshots
- Persistent manual price overrides with audit trail and optional expiry
- Configurable margin thresholds per policy/channel/category
- Clear comparison between current price, suggested price, minimum viable price, and operator override
- Provider-specific fee accuracy through synced fee schedules and explicit stale-data warnings
- Export/report generation from backend for reproducibility
- Scenario modeling for campaigns, marketplace ads, tax rules, and freight strategy

## LLM Implementation Guardrails

When changing this module:

- Start from the formula and status rules above; do not invent frontend math.
- Keep new calculation inputs in Go request/domain/application types first, then expose via OpenAPI and `sdk-runtime`.
- If a new value affects margin, return it in `BatchSimulationItem` so the UI can explain the result.
- If a new dependency is needed, define a pricing port and add an adapter; do not import another module's domain into pricing application code.
- If behavior changes in `/pricing/simulations/batch`, update `contracts/api/marketplace-central.openapi.yaml`, `packages/sdk-runtime`, simulator tests, and this wiki page in the same task.
- Treat missing data honestly. Unknown freight/fees should produce explicit source/status, not optimistic margins.
