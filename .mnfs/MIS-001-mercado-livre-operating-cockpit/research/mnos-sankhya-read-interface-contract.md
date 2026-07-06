# Interface Contract

```yaml
id: IC-002
type: interface-contract
status: planned
owner: Mission Strategist
parent: MIS-001
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: support
```

## Boundary

MPC application modules to MNOS/Sankhya read edge.

## Why This Contract Exists

MPC needs internal product, stock, price, cost, tax, and sales facts. MNOS already mapped the Sankhya Oracle model. This contract prevents MPC feature workers from inventing ad hoc SQL or copying the entire MNOS system.

## Resources Or Entities

- `InternalProduct`: `CODPROD`, description, EAN/reference, group, brand, physical attributes.
- `InternalStockBalance`: stock bucket or aggregated sellable stock.
- `InternalPrice`: current/as-of price.
- `InternalCost`: cost reference from `CUSSEMICM` as-of.
- `InternalSalesHistory`: signed item sales history.
- `InternalTaxInput`: item tax values when needed for margin quality.

## Operations

| Operation | Trigger | Input | Output | Notes |
| --- | --- | --- | --- | --- |
| FindProductsForLinking | Listing import/link candidate generation | EAN, seller SKU, title tokens | product candidates sorted by exactness | Exact EAN/reference beats title heuristics |
| GetSellableStock | Stock reconciliation | `codprod`, policy company/location scope | sellable stock quantity and source buckets | Mission default: `SUM(ESTOQUE - RESERVADO)` where `CODEMP IN (1,2)` and `CODLOCAL=10101` |
| GetCurrentPrice | Price/margin display | `codprod`, optional table/local | current price | Default table `CODTAB=0`, local `CODLOCAL=10101` from MNOS price view |
| GetCostAsOf | Order margin | `codprod`, `codemp`, sale date | `CUSSEMICM` cost | Pick as-of `DTATUAL <= sale_date`, max date |
| GetSalesHistory | Commercial intelligence | product/group/date window | signed item sales aggregates | Use `VW_FAT_VENDA_ITEM` for item/product dimensions |
| GetTaxInputs | Margin quality | order/product/date reference | tax values or missing quality state | Use `VW_IMPOSTO_ITEM` where item-level tax is needed |

## Fields

### Required Inputs

- `codprod`: integer, positive.
- `codemp`: integer when cost or company-specific stock is requested.
- `codlocal`: integer when location-specific stock/price is requested.
- `sale_date`: date for as-of cost and tax context.
- `ean`: nullable string.
- `seller_sku`: nullable string.

### Required Outputs

- `codprod`: integer.
- `produto`: string.
- `ean`: nullable string from product reference/barcode source.
- `codgrupo_prod`: nullable integer.
- `grupo_produto`: nullable string.
- `codmarca`: nullable integer.
- `marca`: nullable string.
- `sellable_stock`: decimal/float, can be negative if stock-reserved is negative.
- `stock_source_scope`: object containing company/location filters.
- `cost_cussemicm`: nullable decimal/float.
- `quality_flags`: list of explicit missing/stale/conflict flags.
- `source_fetched_at`: timestamp.

## Enums And Statuses

- Data quality: `complete`, `missing_product`, `missing_stock`, `missing_cost`, `missing_tax`, `ambiguous_product`, `stale_source`.
- Stock source scope: `revenda`, `showroom_excluded`, `custom_policy`.

## Error Cases

- Oracle/MNOS read unavailable.
- Product not found.
- Multiple exact product matches for a linking key.
- Cost as-of missing.
- Unsupported query shape.

## Persistence Expectations

MPC does not mirror ERP rows wholesale. MPC may persist snapshots used for audit, risk read models, and reproducible action evidence. Snapshot rows must store source timestamps/fetch timestamps and the policy used.

## Canonical Examples

Sellable stock success:

```json
{
  "codprod": 42664,
  "sellable_stock": 3,
  "stock_source_scope": {
    "codemp": [1, 2],
    "codlocal": [10101],
    "formula": "SUM(ESTOQUE - RESERVADO)"
  },
  "quality_flags": [],
  "source_fetched_at": "2026-07-06T15:00:00Z"
}
```

Cost missing:

```json
{
  "codprod": 42664,
  "codemp": 1,
  "sale_date": "2026-07-06",
  "cost_cussemicm": null,
  "quality_flags": ["missing_cost"]
}
```

## Error Matrix

| Case | Status | Code | Notes |
| --- | --- | --- | --- |
| Sankhya read unavailable | 503 | `SANKHYA_READ_UNAVAILABLE` | Blocks refresh, does not erase previous snapshots |
| Product not found | 404 | `SANKHYA_PRODUCT_NOT_FOUND` | Link candidate stays unresolved |
| Multiple exact products | 409 | `SANKHYA_PRODUCT_AMBIGUOUS` | Operator resolution required |
| Cost as-of missing | 200 | quality flag `missing_cost` | Margin is incomplete, not failed ingestion |
| Unsupported query shape | 400 | `SANKHYA_QUERY_UNSUPPORTED` | Contract must be extended before feature use |

## Database Shape

- No ERP mirror tables are introduced by this contract.
- MPC tables may store:
  - `tenant_id`
  - `codprod`
  - source values used for decisions
  - `source_fetched_at`
  - `policy_id`
  - data quality flags

## Seed Data

- One product with EAN and positive sellable stock.
- One product with reserved stock reducing sellable stock.
- One product in showroom-only stock (`CODLOCAL=10108`) expected to produce zero sellable stock under default policy.
- One product missing cost for margin-quality validation.

## Timestamp And ID Semantics

- Sankhya ids remain numeric in internal fields.
- Provider ids remain strings and are not mixed with `CODPROD`.
- Timestamps returned to frontend use RFC3339 UTC.
- Oracle date-only values are converted to date strings where time is not meaningful.

## Compatibility Rules

- MNOS source semantics must be cited in feature specs before adding new query operations.
- New read operations extend this contract; they do not bypass it with module-local SQL.
- Read-only Sankhya access is invariant.

## Route Namespace

- No direct public Sankhya routes are introduced.
- Business routes expose derived MPC views under `/product-links`, `/inventory`, `/orders`, and `/profitability`.

## Transport And Integration

- Sankhya credentials are environment/secret-managed and never logged.
- Read edge must enforce SELECT-only/read-only behavior.
- Feature execution must not add Sankhya write paths.

## Must Preserve

- `SUM(ESTOQUE - RESERVADO)` stock default for Stock Seguro.
- `CODEMP IN (1,2)` and `CODLOCAL=10101` default sellable scope.
- `CUSSEMICM` cost as the initial margin cost basis.
- Missing values become quality flags.

## Must Not Decide In Feature Execution

- Do not replace `CUSSEMICM` with `CUSVARIAVEL` for initial margin.
- Do not include showroom stock in default sellable stock.
- Do not copy complete ERP tables into MPC.

## Validation Impact

- Contract tests must prove default stock excludes `CODLOCAL=10108`.
- Margin tests must prove missing cost yields `missing_cost`, not zero margin.
