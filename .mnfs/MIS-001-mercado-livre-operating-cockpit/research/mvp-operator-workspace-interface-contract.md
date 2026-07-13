# MVP Operator Workspace Interface Contract

```yaml
id: IC-003
type: interface-contract
status: planned
owner: Mission Strategist
parent: MIS-001
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-0
lifecycle_scope: support
```

## Boundary

Object-centered browser workspaces composed from existing SDK/domain read operations.

## Why This Contract Exists

M-09, M-13, and M-14 must use the same product/listing/sale identities, quality
states, route names, attention behavior, and simulation semantics without merging
backend ownership or inventing client-side business rules.

## Resources Or Entities

- `ProductRef`: positive integer `internal_product_id` equal to Sankhya `CODPROD`.
- `ListingRef`: `installation_id`, nonblank `provider_item_id`, nullable
  `provider_variation_id`.
- `SaleRef`: `installation_id`, nonblank `provider_order_id`.
- `SourceFact`: `source`, nullable `value`, `quality`, nullable `observed_at`,
  and nullable `quality_reason`. Quality is server-owned; React never derives it
  from elapsed time, amount, or absence alone.
- `AttentionItem`: stable derived key, `kind`, `severity`, entity ref, `quality`,
  `observed_at`, and target client URL.
- `StockSimulation`: listing ref, current quantity, nullable recommended quantity,
  policy/evidence references, source timestamps, preview payload, `executed=false`.

## Operations

| Operation | Trigger | Input | Output | Notes |
| --- | --- | --- | --- | --- |
| ListAttention | Open `/` or change installation | nullable installation filter | attention items | Severity desc, then non-null `observed_at` before null, timestamp desc, stable key asc |
| ListProducts | Open `/products` | search/filter/page | products | Product name asc, then `internal_product_id` asc |
| OpenProduct | Navigate to `/products/:productId` | positive `productId` | Product 360 sections | Uses existing SDK reads; no client calculations |
| ListListings | Open `/listings` | installation, attention/filter/page | listings | Attention severity desc, title asc, listing identity asc |
| OpenListing | Navigate to listing detail | URL-safe composite listing ref | listing, link, stock, recent sales refs | Sales: provider update desc, then order ID asc; deep links to Product and Sales |
| ListSales | Open `/sales` | installation, quality/filter/page | orders and margin summary | Provider update desc, order id desc |
| OpenSale | Navigate to `/sales/:providerOrderId` | installation plus order id | order, items, source inputs, profit snapshots | Items: provider item asc, null variation last, line ID asc; inputs: source asc then key asc; snapshots: calculated time desc then snapshot ID asc |
| OpenOperations | Open `/operations` | optional installation | connection, auth state, probes, runs | Installations: needs-action first, provider name asc, installation ID asc; probes: non-null observed time before null, time desc, probe ID asc; runs: non-null start time before null, time desc, run ID asc |
| ReviewStockSimulation | Select listing recommendation | ListingRef | StockSimulation | `executed` is always `false` in MVP |

### Attention Derivation

Attention is a client read-model over server-owned states, not a new domain rule.
The browser may construct the stable key/target URL and merge lists, but it copies
server-owned health/risk/quality and never derives stock or margin amounts. Exact
presentation mapping is: unavailable credential/source → `integration/critical`;
link `conflict` → `link/critical`; link `unresolved` → `link/warning`; a server-owned
stock risk → `stock/<server severity>`; incomplete/conflicting margin →
`margin/warning`; other `stale|unknown` source facts → `source_data/warning`.
Duplicates with the same `kind` and entity use one stable key. Resolving the source
condition removes the item on the next read; there is no attention write endpoint.

## Fields

### Required Inputs

- IDs are trimmed and nonblank; `internal_product_id` is a positive integer.
- Installation context is explicit in URL/query state for listing, sale, and operations.

### Required Outputs

- Every external fact supplies `source` and `quality`.
- Known values from external sources supply `observed_at` in RFC3339 UTC.
- Unknown values are JSON `null`, never numeric/string zero placeholders.
- User-visible actions carry a stable target URL.

## Enums And Statuses

- `quality`: `current|stale|unknown|conflict`.
- `attention_kind`: `integration|link|stock|margin|source_data`.
- `severity`: `critical|warning|info`.
- `simulation_state`: `draft|reviewed`.
- Actor attribution shown for local manual records: `operator_supplied_unverified`.

`current` means a successful source observation exists within the domain-owned
freshness window. `stale` means a previously known persisted value exists but its
freshness window elapsed or the latest refresh failed. `unknown` means no usable
value or successful observation exists. `conflict` means two or more eligible facts
cannot be deterministically reconciled. Features use the owning domain's configured
freshness boundary; they may not invent a UI-only duration.

### Quality Rendering Rules

- `current` has a trustworthy known value (including numeric zero) and non-null
  `observed_at`.
- `stale` retains its last trustworthy value and timestamp and requires a
  `quality_reason` explaining expired freshness or failed refresh.
- `unknown` has `value=null` and a required `quality_reason`.
- `conflict` has canonical `value=null`, keeps candidate evidence inspectable, and
  blocks dependent actions.
- Only the owning application/domain assigns quality. Browser code renders the
  received value, quality, and reason without recomputing them.

## Error Cases

- Missing or malformed identity: reject navigation/action and show `invalid_identity`.
- Entity absent from current installation: show `not_found` and preserve list context.
- Source unavailable: keep last persisted facts, mark `stale` or `unknown`, and show
  `source_unavailable`; unknown includes null + nonblank reason; do not zero values.
- Ambiguous product/listing relationship: show `conflict`, block stock simulation review.
- Stock recommendation absent: show `simulation_unavailable`; do not create a payload.

## Persistence Expectations

- Workspaces read domain-owned persisted records; they do not create duplicate UI tables.
- URL path and query state persist selected installation, filters, and entity identity.
- Reloading a deep link restores the same persisted entity or returns `not_found`.
- Simulation review may persist an MPC-local record, but no provider execution record.

## Canonical Examples

Success attention item:

```json
{
  "key": "stock:inst-1:MLB123:-",
  "kind": "stock",
  "severity": "warning",
  "quality": "current",
  "observed_at": "2026-07-13T18:00:00Z",
  "target_url": "/listings/inst-1~MLB123~-?attention=stock"
}
```

Unknown margin fact:

```json
{
  "source": "sankhya_oracle",
  "value": null,
  "quality": "unknown",
  "observed_at": null,
  "quality_reason": "tax fact was not available from the source"
}
```

Stock preview:

```json
{
  "installation_id": "inst-1",
  "provider_item_id": "MLB123",
  "provider_variation_id": null,
  "current_quantity": 8,
  "recommended_quantity": 6,
  "simulation_state": "reviewed",
  "executed": false
}
```

Rejected navigation:

```json
{
  "error": {"code": "invalid_identity", "message": "productId must be a positive integer"}
}
```

## Error Matrix

| Case | Status | Code | Notes |
| --- | --- | --- | --- |
| Malformed product/listing/order identity | 400 | `invalid_identity` | Client route shows a recoverable error and list link |
| Entity not found for installation | 404 | `not_found` | No cross-installation fallback |
| Oracle or Mercado Livre read unavailable | 503 | `source_unavailable` | Prior value → stale; no prior value → unknown with null + reason |
| Product/listing identity ambiguous | 409 | `identity_conflict` | Stock simulation review is blocked |
| No recommendation or required evidence | 422 | `simulation_unavailable` | No preview payload is created |
| Manual local adjustment missing required label/reason | 400 | `invalid_request` | No adjustment record is created |

### IC-002 Boundary Mapping

IC-002 remains authoritative inside the Sankhya read boundary. Its transport codes
map at the IC-003 workspace boundary without changing the original HTTP status:

| IC-002 code | IC-003 code | Workspace quality |
| --- | --- | --- |
| `SANKHYA_READ_UNAVAILABLE` | `source_unavailable` | persisted known value → `stale`; otherwise `unknown` |
| `SANKHYA_PRODUCT_NOT_FOUND` | `not_found` | `unknown` |
| `SANKHYA_PRODUCT_AMBIGUOUS` | `identity_conflict` | `conflict` |
| `SANKHYA_QUERY_UNSUPPORTED` | not exposed publicly | Stop and amend IC-002/IC-003 plus OpenAPI/SDK before feature use |

IC-002 `quality_flags` are source diagnostics, not the IC-003 workspace `quality`.
`complete` maps to `current`; any `missing_*` flag maps the affected fact to
`unknown`; `ambiguous_product` maps to `conflict`; `stale_source` with a prior
trustworthy value maps to `stale`, and without one maps to `unknown`.
The original flags remain available as evidence and are never renamed in storage.

The same preservation rule applies to Mercado Livre/integration error details:
Operations may show the upstream code and next action, while object workspaces map
not-found, transient/unavailable, identity-conflict, and unsupported-preview cases
to the closest IC-003 code above. Features may not invent a third alias.

## Non-API Validation Outcomes

`externally_blocked` is an M-14 orchestration/checkpoint status owned by IC-004.
It is not an HTTP/API code, browser quality, or domain state and must not enter
production API models.

## Database Shape

- No new shared workspace table is mandated.
- Existing domain tables retain their keys and `tenant_id`.
- A persisted simulation, if required, belongs to `inventory` and stores
  `executed=false`; it cannot reuse the provider-action applied state.

## Seed Data

- Deterministic UI fixtures use one installation `inst-mvp-ml`, products `1001` and
  `1002`, listing `MLB-MVP-1`, and order `ORDER-MVP-1`.
- Product 1001 is resolved/current; product 1002 is conflict/unknown.
- Reset recreates only fixture-owned rows and never touches real integration rows.

## Timestamp And ID Semantics

- IDs use domain types above; composite client path segments use URL-safe `~` joining
  in order `installation~item~variation`, with `-` for null variation.
- Timestamps are RFC3339 UTC strings; display localization is browser-only.
- Provider update time and MPC fetch time remain separate fields.

## Compatibility Rules

- Existing server routes and SDK methods remain authoritative for M-13.
- Any missing backend field requires IC-003 amendment and atomic OpenAPI/SDK change.
- Legacy client routes redirect while preserving installation/filter query state.

## Route Namespace

- Server route prefixes remain the existing OpenAPI prefixes: `/catalog`,
  `/integrations`, `/product-links`, `/inventory`, `/orders`, `/profitability`.
- Reserved client routes: `/`, `/products`, `/products/:productId`, `/listings`,
  `/listings/:listingRef`, `/sales`, `/sales/:providerOrderId`, `/operations`.
- Secondary batch simulator remains `/simulator`.
- Legacy redirects: `/product-links` and `/inventory/stock-seguro` → `/listings`;
  `/orders` → `/sales`; `/marketplaces` and `/integrations` → `/operations`.
- No M-13 feature may mount another top-level route.

## Transport And Integration

- No operator auth cookie or session exists in MVP (`sameSite`, `secure`, and
  `httpOnly` are not applicable because no cookie is created).
- Browser uses `createMarketplaceCentralClient` with `VITE_API_BASE_URL`; SDK fetch
  credentials remain omitted.
- Local browser origin is `http://localhost:5174`; local API origin is
  `http://localhost:8080`. `ClientContext.tsx` and the existing development proxy
  own that connection. M-13 must not widen network exposure or introduce credentialed
  browser requests; production origin hardening belongs with post-MVP auth/runtime work.
- Mercado Livre OAuth returns to `MPC_WEB_ORIGIN=http://localhost:5174`; this is an
  integration credential flow, not an operator session.
- Dev proxy ownership remains `apps/web/vite.config.ts`; no broad proxy prefix may be
  added without updating its exact-prefix test.
- Existing Mercado Livre OAuth tokens remain server-side integration credentials.
- Existing `go.mod` and root npm lockfile own versions; M-09/M-13/M-14 add no runtime
  dependency unless a separately reviewed route gap proves it necessary.

## Shared Mutable Seams

- `apps/web/src/app/AppRouter.tsx`, `Layout.tsx`, shared installation context, and
  additions to shared `packages/ui`: M-13/F-01 exclusively owns them until its
  commit is accepted. F-01 wires all reserved routes and legacy redirects to stable
  workspace outlets; M-13/F-02 through F-05 supply outlet components and are
  router/layout/context/`packages/ui` read-only. A proved gap returns to the
  Orchestrator for a serialized F-01 follow-up or newly scoped seam owner.
- `contracts/api/marketplace-central.openapi.yaml` plus
  `packages/sdk-runtime/src/index.ts`: one atomic `api-sdk` seam.
- `apps/web/vite.config.ts` plus its exact-prefix test: one `web-proxy` seam.
- Migration numeric sequence: one writer; M-09/F-03 owns any required forward migration.
- Architecture/MNFS contracts: Milestone Orchestrator integrates changes serially.

## Must Preserve

- Domain ownership, tenant scoping, source timestamps, null/unknown semantics,
  provider payload confinement, and OpenAPI/SDK parity.
- UI language must be consistently Portuguese for operator copy.
- Attention is a view over domain state, not an independent source of truth.

## Must Not Decide In Feature Execution

- Production authentication or authorization.
- Provider-write execution semantics.
- New top-level routes, alternate identity encoding, or new quality enums.
- Client-side stock, tax, fee, cost, or margin calculations.

## Validation Impact

M-09 proves identity; M-13 proves routes, states, deep links, and simulation labels;
M-14 proves the real vertical journey and evidence provenance.
