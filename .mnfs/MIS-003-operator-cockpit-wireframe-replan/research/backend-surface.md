# Research Note

```yaml
id: R-03
type: research
status: draft
owner: Codebase Investigator
parent: MIS-003
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-0
lifecycle_scope: support
```

## Topic

Backend factual state and gap matrix vs the wireframe (planning time, post-M-09, MIS-002 merged at QA-PASS SHA `314b1ef3`).

## Sources Checked

- Source: `apps/server_core/internal/**`, `contracts/api/marketplace-central.openapi.yaml` (version 2026-07-13), `contracts/governance/*`, `docs/architecture/decisions/004..007`, migrations.
- Why it matters: milestone split must build on what exists; inventing parallel machinery is the failure mode.

## Findings

- Go modular monolith `apps/server_core`, strict domain/application/ports/adapters/transport per module. Modules: catalog, classifications, marketplaces, pricing, product_links, internal_read, inventory, orders, profitability, connectors, integrations. Composition root `internal/composition/root.go`.
- ML lives in two adapter seams: `connectors/adapters/mercado_livre/capability_adapter.go` (ProbeAccount, ListListings, FeeQuotes, StockReads, Orders, and an implemented-but-never-live-run `UpdateAvailableQuantity` StockWriter with idempotency-key requirement) and `integrations/adapters/mercadolivre/auth_adapter.go` (OAuth). Sankhya/Oracle: `internal_read/adapters/oracle/` (godror; ADR-006/007).
- OpenAPI surface: ~50 operations across /catalog, /classifications, /integrations (installations, auth, probes, operations, fee-sync), /product-links (snapshots import, candidates, workflows, resolutions approve/reject/manual), /inventory (stock-risks, manual-apply), /orders (import, list, sankhya-linkage), /profitability, /marketplaces (definitions, fee-schedules, accounts, policies), /pricing (simulations, batch), /connectors/melhor-envio. MIS-002 added catalog cursor page/search; oracle batch reads are internal Go ports, not HTTP paths.

### Gap matrix vs wireframe

| Capability | State | Evidence |
| --- | --- | --- |
| Listing entity/read model | **Absent** — only thin import snapshots in product_links (`0022`: no price, no modalidade, no quality, no pendência, no sync state) | migration 0022; connectors `ListingSnapshot` |
| Listing sync scheduler | Absent — pull import operator-triggered only | `POST /product-links/listing-snapshots/imports` |
| Outbound mutation queue/protocolo | **Absent** — zero matches for outbox/queue; background jobs only in integrations (credential refresh, state cleanup, fee-sync scheduler) | grep apps/server_core |
| Price publish to ML | Absent — no PriceWriter capability | connectors module |
| Stock write to ML | Partial — StockWriter adapter implemented, policy-gated, never live-exercised | capability_adapter.go:226; ARCHITECTURE.md:182 |
| Envelope lifecycle precedent | **Exists** — `inventory_stock_actions`: proposed→approved→applied/failed/skipped/blocked, idempotency key, before/after audit | `0026`, stock_action_service.go |
| Product↔listing linkage | Exists fully (workflows, candidates, resolutions, audit) | product_links module |
| Orders read | Exists (import, persistence `0027`, list, assisted Sankhya linkage `0033`) | orders module |
| Sync events read | Partial — `integration_operation_runs` per installation | `0016/0017` |
| Competitor/market data | Absent in code; research corpus only | grep + docs/research |
| Fee schedules | Exists (list/seed/sync + scheduler) | marketplaces module |
| Price policies + simulator | Exists (policies, single + batch simulation) | pricing module |
| Dashboard aggregates (exception counters) | Absent as endpoint | — |

- Multi-tenancy: tenant-ready schema (all business tables carry `tenant_id`), single-tenant in practice via `cfg.DefaultTenantID` (`MC_DEFAULT_TENANT_ID`, default `tenant_default`) injected at composition. Oracle read ports deliberately carry no tenant dimension. ML account entities: `integration_installations`, `integration_credentials` (encrypted), `integration_oauth_states`, `integration_auth_sessions`, `integration_capability_states`, `integration_operation_runs`.
- Governance: `contracts/governance/execution-lanes.json` defines lanes unit/integration/dev-invariance/live-oracle (read-only + session-temporary-write)/live-provider-read/browser/**provider-write** (gates: actor, idempotency, execute, resolved-link, policy, source-timestamp, before-after-audit; evidence class production-like). `modules.json`, `runtime-config.json`, `invariants.json`, `shared-seams.json`; validated by `scripts/tests/governance-contracts.tests.ps1`.
- IC-01 dormant row (MIS-002 catalog-read-interface-contract.md:102,108): "product edit (DORMANT) | catalog | ['catalog']" — any reintroduced product-edit surface MUST invalidate server class `catalog` and client key `['catalog']`.
- L2 freshness cache: `internal_read/adapters/cache/cache.go` — catalog 300s / stock 45s / pricecost 120s (env-tunable), LRU + singleflight, `InvalidateClass` generation counter, `Cache-Control: no-cache` bypass. Decorates CatalogPageReader + cost/tax BatchReader + StockBatchReader.
- Build/test: harness lanes via `scripts/harness.ps1` (npm `harness:*`); GOCACHE absolute (`apps/server_core/.gocache` in harness, root `.gocache` dev-local); live Oracle governed runner `scripts/run-live-oracle-docker.ps1` + `docker/live-oracle/profile.json`.

## Recommendation

Build the listings read model as a NEW module (product_links stays a linkage-workflow domain); implement the mutation envelope as one `mutation_protocols` table + in-process poller (no outbox/bus), folding `inventory_stock_actions` semantics in; reuse connectors `ListListings` as sole listing ingestion path.

## Impact On Mission

Grounds ADR-12/13; fixes milestone dependency order (listings spine first); confirms provider-write gate vocabulary already exists and must be implemented, not invented.

## Handoff

- Current status: complete.
- Next owner: Mission Strategist.
- Next action: none.
- Required files/evidence: paths cited above.
- Blockers or open decisions: none.
