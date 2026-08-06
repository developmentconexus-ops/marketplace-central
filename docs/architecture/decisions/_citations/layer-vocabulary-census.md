# Layer-vocabulary census — the seven non-canonical directory names

Read-only investigation. Scope: `apps/server_core/internal/modules/**` for the seven names
(`readmodel`, `events`, `composition`, `background`, `registry`, `observability`, `integration`),
plus the two verification claims about `sourcekind` and `tenant_config`.

Canonical five (for contrast, not re-audited here): `domain`, `application`, `ports`, `adapters`,
`transport`.

---

## readmodel

| Occurrence | .go files | non-blank LOC | doc.go-only? |
|---|---|---|---|
| `catalog/readmodel/doc.go` | 1 | 1 | yes |
| `connectors/readmodel/doc.go` | 1 | 1 | yes |
| `marketplaces/readmodel/doc.go` | 1 | 1 | yes |
| `pricing/readmodel/doc.go` | 1 | 1 | yes |

Content, verbatim, all four files: a single line, `package readmodel`. No comment, no type, no
function.

What it does: nothing. Each is an empty package stub — the directory exists and compiles but
declares zero symbols.

Importers: `grep -rn "modules/<x>/readmodel"` across `apps/server_core` (excluding `_test.go`)
returned **zero matches** for all four occurrences — no importer, cross-module or in-module.

## events

| Occurrence | .go files | non-blank LOC | doc.go-only? |
|---|---|---|---|
| `catalog/events/doc.go` | 1 | 1 | yes |
| `connectors/events/doc.go` | 1 | 1 | yes |
| `marketplaces/events/doc.go` | 1 | 1 | yes |
| `pricing/events/doc.go` | 1 | 1 | yes |

Content, verbatim, all four files: a single line, `package events`.

What it does: nothing — same empty-stub shape as `readmodel`.

Importers: `grep -rn "modules/<x>/events"` (excluding `_test.go`) returned **zero matches** for
all four — no importer anywhere.

## composition (in-module occurrences)

| Occurrence | .go files | non-blank LOC | doc.go-only? |
|---|---|---|---|
| `listings/composition/` | 3 (incl. 2 test files) | 398 | no |
| `orders/composition/` | 2 (incl. 1 test file) | 179 | no |
| `product_links/composition/` | 2 (incl. 1 test file) | 128 | no |
| `sync/composition/` | 4 (incl. 1 test file) | 377 | no |

What it does, per occurrence:
- `listings/composition/scheduler.go` — builds one `sync/application.Scheduler` per active
  Mercado Livre installation and registers the listings backfill/sweep job on each, because a
  listings run is scoped to one ML installation while `sync/composition`'s pattern assumes a
  single tenant-wide scheduler (comment at `listings/composition/scheduler.go:1-13`).
- `orders/composition/scheduler.go` — same fan-out-per-installation pattern, wiring
  `orders/application.NewOrdersJob` onto a `sync/application.Scheduler` per ML installation so
  order-sync failures surface through `sync_state`/`/sync/health` instead of a silent parallel
  ticker (`orders/composition/scheduler.go:1-16`).
- `product_links/composition/refresher.go` — `LinkCandidateRefresher` regenerates product-link
  candidates for every installation of a tenant whenever product data changes
  (`product_links/composition/refresher.go:21-29`).
- `sync/composition/` — three files: `installation_scheduler.go` builds a single-installation
  `sync/application.Scheduler` with no job registered yet (extracted so other modules stop
  importing `sync/adapters/postgres` directly, which would violate `GOV_MODULE_LAYER`);
  `products_job.go` wires the ERP-sourced products sync job and defines its cursor type;
  `scheduler.go` is a package-doc file plus two exported sentinel constants,
  `InstallationScopeERP` and `InstallationScopeMarket`, used as `sync_state.installation_id`
  values for tenant-wide (non-per-installation) entities.

Cross-module importers (non-test, `file:line`):
- `sync/composition` — imported from **3** owning-module-external files:
  `internal/composition/root.go:119`, `internal/modules/listings/composition/scheduler.go:29`,
  `internal/modules/orders/composition/scheduler.go:28`.
- `listings/composition` — imported from **1**: `internal/composition/root.go:69`.
- `orders/composition` — imported from **1**: `internal/composition/root.go:95`.
- `product_links/composition` — imported from **1**: `internal/composition/root.go:110`.

In-module importers: none observed — each package's own scheduler/refresher file is the sole
non-test `.go` file in its directory; nothing else inside the same owning module imports it.

### Top-level `composition/` package vs. the in-module `composition/` dirs

The top-level package lives at `apps/server_core/internal/composition` (not under
`internal/modules/`) — package name `composition`, 11 files, 3331 total lines
(`root.go` 1020 lines, `market_adapters.go` 679, `pricing_adapters.go` 171,
`orders_adapters.go` 118, `orders_ingest_adapters.go` 89, plus four `_test.go` files and
`catalog_routes_test.go`). It is the **application root / composition-root**: `root.go` imports
adapters, application services, ports, and transport packages from essentially every domain
module (catalog, classifications, connectors, dashboard, erp_import, integrations, …), wires
concrete adapters to ports, registers HTTP routes, and starts schedulers — the single place the
whole dependency graph is assembled and (per grep) `http.ListenAndServe` is called.

Difference from the in-module `composition/` dirs: the top-level package is the **one and only**
composition root for the whole binary and owns cross-module wiring + server startup. The four
in-module `composition/` dirs (`listings`, `orders`, `product_links`, `sync`) are each scoped to
**one owning module** and do narrower, single-purpose wiring (one scheduler or one refresher);
they exist specifically so that other modules and the top-level root can obtain a ready-to-use
scheduler/refresher without importing that module's `adapters/postgres` package directly (stated
motivation in `sync/composition/installation_scheduler.go:12-18`). The top-level root then
imports all four in-module composition packages (anchors above) — it is a consumer of them, not
a duplicate of them.

## background

| Occurrence | .go files | non-blank LOC | doc.go-only? |
|---|---|---|---|
| `integrations/background/` | 6 (incl. 3 test files) | 382 | no |
| `mutations/background/` | 2 (incl. 1 test file) | 200 | no |

What it does, per occurrence:
- `integrations/background/fee_sync_scheduler.go` — periodic job that lists installations and
  syncs marketplace fee schedules (`integrations/background/fee_sync_scheduler.go:1-14`).
- `integrations/background/refresh_ticker.go` — periodic ticker that lists sessions expiring
  soon and triggers token refresh (`integrations/background/refresh_ticker.go:1-14`).
- `integrations/background/state_cleanup.go` — `StateCleanup` deletes expired OAuth state rows
  on an interval (`integrations/background/state_cleanup.go:1-15`).
- `mutations/background/poller.go` — `Passer`/`InstallationProvider`-driven background poller
  over a set of installations (`mutations/background/poller.go:1-14`).

Cross-module importers (non-test, `file:line`):
- `integrations/background` — imported at `internal/composition/root.go:43`.
- `mutations/background` — imported at `internal/composition/root.go:86`.

Both counts: **1** cross-module importer each (the top-level composition root), **0** in-module
importers observed.

## registry

| Occurrence | .go files | non-blank LOC | doc.go-only? |
|---|---|---|---|
| `marketplaces/registry/` | 4 (incl. 1 test file) | 180 | no |

What it does: `marketplaces/registry` is a plugin registry — `registry.go` holds a package-level
`[]MarketplacePlugin` slice and a `register()` function called from each plugin's `init()`
(`marketplaces/registry/registry.go:1-19`); `plugin.go` declares the `MarketplaceConnector`
interface and `ErrNotImplemented`; `mercado_livre.go` registers the Mercado Livre plugin via
`init()` (`marketplaces/registry/mercado_livre.go:9-13`).

Cross-module importers (non-test, `file:line`): **2**
- `internal/composition/root.go:78`
- `internal/modules/integrations/adapters/feesync/marketplace_executor.go:11`

In-module importer: **1** — `internal/modules/marketplaces/application/fee_schedule_service.go:9`.

## observability

| Occurrence | .go files | non-blank LOC | doc.go-only? |
|---|---|---|---|
| `internal_read/observability/` | 4 (incl. 1 test file) | 467 | no |

What it does: `config.go` defines a `Config` struct plus default thresholds (slow-query
threshold, pool-stats interval) (`internal_read/observability/config.go:1-13`); `pool_stats.go`
periodically samples a `PoolStatsSource` (backed by `internal_read/ports.PoolStats`) and logs
pool statistics (`internal_read/observability/pool_stats.go:1-13`); `timing.go` wraps
`internal_read` query execution to log slow queries, using `internal_read/domain` and
`internal_read/ports` (`internal_read/observability/timing.go:1-13`).

Cross-module importers (non-test, `file:line`): **1** — `internal/composition/root.go:53`.
In-module importers: **0** observed (the package is self-contained; nothing else under
`internal_read/` imports it in non-test code).

## integration

| Occurrence | .go files | non-blank LOC | doc.go-only? |
|---|---|---|---|
| `mutations/integration/` | 2 (both `_test.go`) | 249 | no (not a `doc.go`, but 100% test files) |

Confirmed: the directory contains **only** `_test.go` files — `error_matrix_test.go` and
`lifecycle_test.go`. Both declare `package integration_test` and both carry the build tag on
line 1:

```
//go:build integration
```
(verified at `mutations/integration/error_matrix_test.go:1` and
`mutations/integration/lifecycle_test.go:1`).

Importers: **0** — `grep -rln "modules/mutations/integration"` across non-test `.go` files
returned no matches; being `package integration_test` with an `integration` build tag, this
directory is not importable as a library package at all — it is a hermetic-lane test fixture.

---

## Verification claims

**Claim 1 — `sourcekind` has no layer dirs at all.**
**TRUE.** `apps/server_core/internal/modules/sourcekind/` contains exactly two files at the
module root: `sourcekind.go` and `sourcekind_test.go`. No subdirectories exist (confirmed via
directory listing) — no `domain/`, `application/`, `ports/`, `adapters/`, `transport/`, nor any
of the seven names under investigation.

**Claim 2 — `tenant_config` has ONLY `transport/` and no domain/application/ports/adapters.**
**TRUE for the "no domain/application/ports/adapters" half; the module is not "only
`transport/`"** — it also has non-layer files sitting directly at the module root.
`apps/server_core/internal/modules/tenant_config/` contains: `active_source.go`,
`active_source_test.go`, `context.go`, `context_test.go`, `repository.go`, `repository_test.go`
at the module root, plus one subdirectory `transport/` (`http_handler.go`,
`http_handler_test.go`). No `domain/`, `application/`, `ports/`, or `adapters/` directories
exist under `tenant_config/`.

What the root-level code does: `active_source.go`/`context.go` define the `Config`,
`ActiveSource`, `SourceKind`, `SellableAssortment` types and the `ActiveSourceLookup` interface
(domain-shaped code with no directory to hold it); `repository.go` defines `Repository`, a
Postgres-backed implementation of `ActiveSourceLookup` that both reads and writes the
`active_source` table directly via `pgxpool.Pool.QueryRow`/`Exec` (adapter-shaped code sitting at
module root, not under `adapters/`).

SQL location claim — **the exact quoted line does not match**. `tenant_config/repository.go:40`
is:

```go
	FROM active_source
```

— one line inside a multi-line SQL string that starts at line 37 (`r.pool.QueryRow(ctx, \`` )
and ends at line 43. So the plan's claimed anchor `tenant_config/repository.go:40` **is a real
line inside a raw SQL string** (`FROM active_source`), confirming the substance of the claim —
SQL text lives directly in `repository.go`, not behind an `adapters/postgres` package — but the
single line at :40 in isolation is only the `FROM` clause fragment, not the full statement. The
three SQL statements in the file span `repository.go:37-43` (SELECT), `74-82` (INSERT ... ON
CONFLICT), and `92-98` (UPDATE). `Repository.Set` writes via `Exec` and `Repository.Get` reads
via `QueryRow`/`Scan`, both directly against `*pgxpool.Pool` with no adapter-layer indirection.

---

## Judgment inputs

| Name | Strongest fact FOR keeping | Strongest fact AGAINST |
|---|---|---|
| `readmodel` | None of the 4 occurrences has any content or importer to lose — keeping costs nothing but also delivers nothing. | All 4 occurrences are a single `package readmodel` line with zero types/functions and **zero importers anywhere** (grep across the whole tree) — pure dead stub, 4 directories for nothing. |
| `events` | Same neutral cost-to-keep as `readmodel`. | Identical shape to `readmodel`: 4× single-line `package events`, **zero importers**. |
| `composition` (in-module) | Cross-module wiring is real and load-bearing: `sync/composition` alone has 3 cross-module importers and exists specifically to let other modules avoid importing `adapters/postgres` directly (a stated anti-`GOV_MODULE_LAYER`-violation measure). | The name collides with the top-level `internal/composition` package that already owns "composition root" — two things named `composition` at different scopes/purposes in the same codebase is exactly the vocabulary ambiguity the census was commissioned to resolve. |
| `background` | Both occurrences are substantive, non-trivial scheduled-job code (fee sync, token refresh, OAuth state cleanup, installation polling) genuinely wired in from the root (`root.go:43`, `root.go:86`). | Only 1 cross-module importer each (the root itself) and 0 in-module importers — the directory boundary buys no real encapsulation beyond "things called from a timer," which `application` could also hold. |
| `registry` | Has the richest import graph of the seven: 2 cross-module + 1 in-module importer, and a real behavioral contract (plugin `init()`-registration pattern) that doesn't map cleanly onto any of the five canonical layers. | Only exists in one module (`marketplaces`) — a single-occurrence pattern is cheap to fold into `application` or `adapters` without inventing a sixth vocabulary word. |
| `observability` | Self-contained, real code (467 LOC) with a single clear cross-cutting purpose (pool stats + slow-query logging) that doesn't fit `domain/application/ports/adapters/transport` naturally. | Single occurrence, single importer (`root.go:53`), zero in-module importers — nothing forces it to be its own layer name rather than living under `adapters` (it wraps `ports.PoolStats`/`internal_read.domain` already). |
| `integration` | Zero risk to keep: it is a hermetic hidden-behind-a-build-tag test fixture (`//go:build integration`), not a shippable package, so it competes with nothing in the "five canonical layers" vocabulary. | It isn't a layer at all — 100% `_test.go`, `package integration_test`, unimportable — so counting it alongside real layer names conflates "test lane" with "architecture layer." |
