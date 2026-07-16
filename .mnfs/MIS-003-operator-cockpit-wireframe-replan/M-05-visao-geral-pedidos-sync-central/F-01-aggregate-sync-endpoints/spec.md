# F-01-aggregate-sync-endpoints — spec

```yaml
id: F-01
type: feature-spec
status: implemented
owner: CHIP-SAT (wave W1, MIS-003)
parent: M-05
created: 2026-07-16
updated: 2026-07-16
base_anchor: a49168e641ffd6f61932ca57c29b1d1bdcde2fb0
branch: chip/sat-m05f01-m06f02
```

## What was built

Three backend-only read surfaces (zero frontend files touched; the feature-brief line about
vite proxy rows was superseded by the mission plan seam map — proxy rows are frontend seam,
owned by M-05 FE chips):

### 1. Dashboard summary — `GET /dashboard/summary?installation_id=`
- New thin `dashboard` module (`apps/server_core/internal/modules/dashboard/{application,transport}`)
  composing four module application services via their ports/domain layers only: orders
  counters, listings summary (IC-02 reuse), product_links pending/missing-GTIN summary,
  integrations last-sync projection. No cross-module SQL.
- ADR-17 honest degradation: each counter nullable; a failed sub-read → JSON `null` for that
  counter + source name appended to `degraded[]`; all-sources-failed → 200 all-null + full
  degraded list. Never 0, never 500.
- Unknown installation → 404 `installation_not_found`. Tenant scoping on every sub-read.
- `last_sync_at` map values nullable (`Record<string, string | null> | null`).

### 2. Orders read API — evolve `GET /orders` in place + new `GET /orders/{provider_order_id}`
- Existing orders module evolved additive-only: tenant-scoped keyset pagination
  (`created_at DESC, id` tie-break), filters `status`, `date_from`, `date_to`, `q`
  (BARE param names — transport 400s `filter.*` forms), `limit`, opaque `cursor`.
- Canonical read model only (`domain.OrderReadModel`); `raw_provider_ref` never serialized
  in read paths (provider payloads stay at adapters). NF state nullable, exact-evidence only.
- Malformed cursor → 400 `invalid_cursor`; malformed dates → 400 `invalid_filter`;
  unknown order → 404 `order_not_found`; cross-tenant read impossible (predicate, tested).
- **Lock-exception D-02 (hub ruling 2026-07-16):** legacy `listMarketplaceOrders` 200
  response retyped to the new `ListOrdersResponse` envelope (contract follows committed
  runtime — the old `MarketplaceOrder` shape with required `installation_id`/`fetched_at`
  was no longer emitted anywhere after R1). `MarketplaceOrder` schema kept with
  `deprecated: true`; legacy SDK method kept `@deprecated` (live consumer:
  `packages/feature-orders/src/OrdersPage.tsx`); `importMarketplaceOrders`/POST untouched.

### 3. Sync runs — `GET /sync/runs`
- Integrations module: keyset list over `integration_operation_runs`,
  `started_at DESC` newest first, filters `module`, `status`, fixed 90-day window
  (`started_at >= now() - 90d`, mission §Audit & history). Running run → `status=running`,
  `finished_at` null.

### Contracts + composition
- OpenAPI (`contracts/openapi.yaml`) + `sdk-runtime` same-commit pairs:
  C2 = `getDashboardSummary` + `listSyncRuns` (491b18fb); C1 = `listOrders` + `getOrder`
  + D-02 retype (ce12fd7f).
- Composition root wires dashboard + sync-runs registrations and swaps orders to
  `NewHandlerWithReader` (a9bd1511) with route smoke tests.
- Governance registration: `dashboard` module + `/sync` prefix in
  `contracts/governance/modules.json` (17cce1a9).

## Slice map (commit per reviewed-green slice)

| Slice | Commit | Content |
|---|---|---|
| O1+S1 | 1fc12534 | orders query/cursor grammar + sync-runs keyset repo |
| D2 | bee6df32 | product-link pending + missing-GTIN summary service |
| S2 | 5950b094 | sync-runs application list + GET /sync/runs transport |
| O2 | 56c531f2 | tenant-scoped keyset order reads (+O1 review conditions) |
| D1 | 1efc4e7b | today/7d order counters |
| D3 | 484b204b | latest sync-run projection per module |
| D4 | 59aa7e15 | honest summary composition (ADR-17) |
| D5 | 970f1583 | GET /dashboard/summary transport |
| O3 | fa447d65 | /orders evolved in place + installation-scoped detail |
| R1 | a9bd1511 | composition root wiring |
| C2 | 491b18fb | dashboard/sync OpenAPI + SDK |
| I1 | 07bc2122 | end-to-end integration test (`TestAggregateSyncReadContract`) |
| C1 | ce12fd7f | orders OpenAPI + SDK (D-02) |
| gov | 17cce1a9 | modules.json registration |

Every slice: failing test first, independent sonnet review (REVIEW-STANDARD), ACCEPT before
dependent slice. Review verdicts + dispatch rows: `CHIP-SAT-DISPATCH-LEDGER.md` rows 1–43b.
