# F-02 Slice 7 validation

## Scope and composition

Integration tests close M01-C04..C09 for the read-only listings spine. They drive real HTTP requests through `httpx.NewRouteClassMux`, the production `ReadHandler.Register`, `NewReadService`, the real Postgres listings repository over harness migrations, and the real integrations Postgres repository/service/adapter for installation existence. The sole external-data substitution is a deterministic in-process `ports.CostReader` fake at the Oracle boundary; it is explicit fake evidence, not live-Oracle evidence. Because `root.go` deliberately composes an unavailable marketplace-policy source, the tests use the same explicit pricing-policy contract as the listings read-service unit tests.

## Contract evidence

`TestListingsReadContractEndToEnd` covers the six-row real-cursor page walk; exact envelope keys and present-null fields; every filter key (including `exception=sync_error`); title/provider-ID/seller-SKU search; by-product counts, states, tie ordering and null-last grouping; variation detail identity; last-ten timeline ordering; malformed/absent 404s; the status plus `error.code` matrix; ADR-17 null-cost honesty; deterministic known below-margin behavior and summary counters; and cross-tenant isolation across list, q, filter, get, group, summary, and cursor continuation.

## Performance evidence

`TestListingsReadPerformance2000` seeds 2,000 tenant/installation-scoped listing, product-link, and snapshot rows, warms the endpoint, records 100 sequential `limit=50` calls, calculates nearest-rank p95, captures the keyset `EXPLAIN`, and traces the summary endpoint to prove exactly one listings conditional-aggregate query.

- p95: **3.2563 ms** (nearest-rank over 100 sequential `GET /listings?limit=50` calls, 5 warmups) — well under the 500 ms ceiling.
- keyset plan/index line: `Index Only Scan using idx_listings_f02_title_key on listings l` with `Index Cond: ((tenant_id = …) AND (installation_id = …) AND (ROW(title, provider_listing_id, variation_id) > ROW('Title 1000','MLBPERF1000','-')))`, wrapped by `Limit (rows=51)`. No `Seq Scan`. Planner stats primed with `ANALYZE listings` after the 2,000-row bulk load so the plan is deterministic.
- summary listings aggregate query count: **1** (exactly one tenant-scoped `count(*) FILTER` conditional-aggregate query, traced via pgx `QueryTracer` — D-20 one-query rule; the Oracle cost batch is a separate bounded port call, not a per-row query).
- environment/SHA and 100 samples: go=go1.26.4, os=windows/amd64, SHA=(local worktree); 100 samples all in the 1.5–4.1 ms band (min ~1.55 ms, max ~4.03 ms).
- exact lane commands: ephemeral `postgres:16-bookworm` → `CREATE DATABASE mpc_test_<32hex>` → `go run ./cmd/testdb migrate` (applied 36) → `go test -tags=integration -run TestListingsRead -v -count=1 ./tests/integration`.

## Lane result (2026-07-15, milestone-owner run)

`go test -tags=integration -run TestListingsRead -v` → **PASS** (`ok marketplace-central/apps/server_core/tests/integration 5.994s`):
- `TestListingsReadPerformance2000` PASS (0.77s) — p95 3.2563 ms, keyset Index Only Scan, summary query count 1.
- `TestListingsReadContractEndToEnd` PASS (2.06s) — all 8 subtests green: cursor walk + JSON null contract, all filter keys, q (title/provider-id/SKU), by-product tie/null-last, detail timeline + 404, error matrix (status + `error.code`), null-cost honesty + known below-margin + summary counters, tenant isolation across all read paths + cursor.

Real ephemeral Postgres, real migrations (0036), real registrar/handler/repository; the only substitution is the deterministic `ports.CostReader` fake at the Oracle boundary (explicit fake, not live-Oracle evidence). Completes **M01-C04..C09**.
