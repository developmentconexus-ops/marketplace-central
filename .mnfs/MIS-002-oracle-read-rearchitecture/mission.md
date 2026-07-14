# MIS-002-oracle-read-rearchitecture

```yaml
id: MIS-002
type: mission
status: planned
owner: Mission Strategist
parent: none
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-2
lifecycle_scope: mission
planning_phase: ready
```

## Objective

Oracle-backed reads in `apps/server_core` answer in O(1) Oracle queries per page (today 1+3N sequential round trips), protected by a server-side TTL cache honoring the existing `FreshnessPolicy`, and consumed by `apps/web` through TanStack Query — with no Oracle data mirror, preserving ADR-006 (direct read) and the "unknown never becomes zero" rule.

## Outcome

- `GET /catalog/products` returns a cursor-paginated page (50 items) built from exactly one Oracle query, in sub-second latency for 2–10 concurrent operators over a 30k–100k product catalog.
- Inventory stock-risk and profitability import paths use batch IN-list reads (chunked at 500) instead of per-item loops.
- Sankhya assisted-linkage operations cost 1–2 Oracle round trips (configuration validation cached), always fresh (never cached).
- Every Oracle-backed response carries `as_of`; operators can force-refresh past the cache.
- Latency per port method, slow queries (>1s), and pool stats are observable in logs.

## Scope

Blocks chosen at P1 gate: Core (observability baseline + batch/paginated ports + pool 12 + deadlines), Server cache L2 (FreshnessPolicy TTL + singleflight), Frontend TanStack Query. Runtime-dimensioning report injected 6 hardening decisions (see `research/oracle-runtime-dimensioning.md`).

## Domain Scope

- (2) Core entities: catalog product page facts (product+stock+price+cost), sellable stock batch, cost/tax batch, sales history (capped), Sankhya linkage candidates/descendants.
- (3) Lifecycle: cache entry fresh→stale (TTL per data class)→evicted-on-mutation.
- (6) Audit & history: slow-query log, per-method latency, pool stats; linkage reads never cached (audit freshness).
- (8) Search/filter: catalog search bounded (FETCH FIRST 50, debounced client-side).
- (9) Admin/config: pool/TTL/timeout env knobs (existing `MPC_ORACLE_*` plus new cache/deadline knobs).
- (10) Integration: Oracle/Sankhya via `internal_read` adapters only; web via `packages/sdk-runtime`.

## Non-Scope

- F5 Postgres read-model / Oracle mirror — conditional on post-M-04 metrics; requires ADR-006 revision first (operator decision).
- Resilience block (transient-error retry/backoff, graceful shutdown + `Close()`, `/readyz` Oracle probe) — operator declined this mission; acceptable for internal tool; blips surface as errors until retried by the operator.
- Redis / distributed cache — single-instance deployment today; in-memory cache suffices.
- Numeric performance SLO (p95 target) — operator declined; performance validated structurally (query count per page vs baseline).
- Link-candidates generation batching beyond 2k listings/run — known limit with explicit trigger (risk R7).

## Current State

- Oracle adapter (`apps/server_core/internal/modules/internal_read/adapters/oracle/`) on godror v0.51.0; uncommitted hardening refactor in working tree (structured ConnectionParams, CallTimeout, `Database` interface, validated pool envs) — MUST be preserved and integrated, not discarded.
- N+1 read paths: catalog 1+3N (`catalog/adapters/internalread/reader.go:67-111`), profitability per item×tax-source (`profitability/application/service.go:239-287`), inventory per snapshot (`inventory/application/stock_risk_service.go:29-84`), Sankhya per-candidate lines + per-call ValidateConfiguration (`sankhya_linkage_reader.go:104,160`).
- Pool default 4 sessions (`config.go:13`); no HTTP server timeouts (`cmd/server/main.go:35`); `FreshnessPolicy` plumbed but unread; no metrics; `Database.Stats()` never called.
- `apps/web`: React 19 + Vite, hand-written TS SDK (`packages/sdk-runtime`), no data-fetching cache library.
- Full session evidence: `research/oracle-runtime-dimensioning.md`.

## Clarified Decisions

- Resolved:
  - Scope blocks: Core + Server cache L2 + TanStack Query (Resilience declined).
  - Quality: baseline only (maintainability, credential redaction, CGO pinning); no numeric perf SLO.
  - Catalog pagination: NEW cursor envelope `{items, next_cursor, page_size}` — breaking; web is sole consumer and migrates in the same milestone. OpenAPI + `packages/sdk-runtime` updated together.
  - Execution harness: existing MNFS repo flow (Portfolio→Milestone→workers) with hub-style acceptance: independent reviewer per feature before accept, `merge --no-ff` per milestone, post-merge test ladder.
  - Execution workers: `mpc-implementer` and `mpc-verifier` pinned to `gpt-5.6-luna` high reasoning (operator directive 2026-07-13); lean process, minimum ceremony.
- Accepted assumptions (each reversible):
  - Pool max sessions default 12 (env `MPC_ORACLE_POOL_MAX_SESSIONS`; DBA validates at rollout; env-only change to revert).
  - Entity-shaped ports (GetSellableStock etc.) remain for single-item lookups; only hot paths move to batch ports.
  - Cache is in-memory per instance + `golang.org/x/sync/singleflight`; Redis only if multi-instance later (cache is an adapter seam, swappable).
  - Cursor = opaque base64 of keyset key (CODPROD); documented in IC-01.
  - TTLs: catalog/taxonomy 5min, stock 45s, price/cost 2min, linkage 0 (never); env-tunable.
  - TanStack staleTime mirrors L2 TTLs.
  - Interactive route deadline 15s; batch route deadline 120s; batch Oracle semaphore 4.
  - Catalog scale assumed 30k–100k active products and LAN latency to Oracle — both MEASURED in M-01 before M-02 cutover (risks R1/R8).
  - Oracle read ports have NO tenant dimension (no port takes TenantID; ERP facts are installation-global, single-tenant deployment). L2 cache key therefore excludes tenant. Revisit + add tenant to key AND an isolation criterion if per-tenant Oracle scoping is ever introduced (reversible: key formula is one adapter function).
- Owner decisions still open: None - all P1 answers received 2026-07-13.
- Blocked items: None - planning proceeded to validation phase.

### Clarification Interview

| # | Taxon | Question | Proposed default | Operator answer |
| --- | --- | --- | --- | --- |
| 1 | persistence/reset | Capability blocks in scope | Core only | Core + Cache L2 + TanStack (Resilience out) |
| 2 | validation expectations | Quality bars | Baseline | Baseline only (no numeric p95) |
| 3 | UI convergence | Catalog pagination API shape | Cursor envelope (breaking) | Cursor envelope confirmed |
| 4 | build/runtime conventions | Execution harness style | MNFS + hub acceptance | MNFS + hub acceptance; workers gpt-5.6-luna high; minimum bureaucracy |

## Architecture Spine

### System Shape

Three boundaries: React SPA (`apps/web`) → Go API (`apps/server_core`) → Oracle ERP (read-only via `internal_read` adapters) and Postgres (marketplace-owned data, unchanged). The SPA never reads Oracle or Postgres directly. Oracle SQL and driver types exist only in `internal_read/adapters/oracle`.

### Runtime Topology

- Single `apps/server_core` process; one godror pool (12 sessions) opened at composition root.
- Cache L2 lives in-process, keyed per (port method, canonical params). No tenant key component: Oracle read ports carry no tenant dimension (ERP facts installation-global — see Accepted assumptions).
- `apps/web` served by Vite dev / static build; consumes API via `packages/sdk-runtime` wrapped by TanStack Query.

### Runtime Contract

- Oracle: source of truth for product/stock/price/cost/tax/sales facts; read-only; unavailability → typed `source_unavailable`, never fallback data, never zero.
- Postgres: sole store for marketplace-owned entities; untouched by this mission.
- Route classes: interactive (15s context deadline) vs batch (120s) — declared per route at transport registration; batch Oracle work passes through a 4-permit semaphore.
- Every Oracle-backed response includes `as_of` (RFC3339 UTC, the time facts were read from Oracle or served from cache).

### Cross-Cutting Decisions

| Decision | Status | Prevents | Must preserve | Validation impact |
| --- | --- | --- | --- | --- |
| ADR-01 Use-case-shaped paginated/batch ports in `internal_read` (1 Oracle query per page/batch, keyset cursor) | decided | N+1 re-emerging via caller composition | nil+quality-flags (never zero); SQL only in adapter | query-count criteria per page |
| ADR-02 Catalog API cursor envelope `{items, next_cursor, page_size}` + `as_of` | decided | two consumers inventing pagination | OpenAPI + sdk-runtime updated in same commit | contract tests + web renders page |
| ADR-03 In-memory TTL cache + singleflight honoring `FreshnessPolicy.MaxAge`; linkage never cached; evict-on-mutation | decided | per-consumer ad hoc caches; stale audit reads | MaxAge=0 bypasses cache (force refresh) | cache-hit log criterion; linkage-uncached criterion |
| ADR-04 Pool 12; http.Server Read/Write/ReadHeader timeouts; per-route-class context deadlines (15s/120s); batch semaphore 4 | decided | pool exhaustion; unbounded request duration; batch starving interactive | godror CallTimeout stays | timeout + concurrency criteria |
| ADR-05 IN-list chunking at 500; `ImportMarginInputs.Limit` ceiling 200; `GetSalesHistory` row cap 5000 | decided | ORA-01795; unbounded imports/result sets | explicit 422 when Limit ceiling exceeded; explicit `truncated=true` marker when sales cap hit — never SILENT truncation | M-02-C03 (limit range), M-02-C06 (search bound), M-03-C02 (sales truncation), M-03-C06 (import ceiling) |
| ADR-06 Execution: MNFS + hub-style acceptance; workers gpt-5.6-luna high; EXPLAIN PLAN gate before M-02 cutover | decided | slow view plans shipping to prod | live-Oracle lane read-only (`scripts/run-live-oracle-docker.ps1`) | M-01-C04 precondition for M-02 |

Accepted trade-offs: ADR-02 breaks the current catalog response shape (single consumer, migrated same milestone); ADR-03 cold cache per instance and staleness ≤TTL (visible via `as_of`); ADR-04 higher Oracle session pressure (DBA validates); ADR-05 caps can reject oversized legitimate requests (explicit 422 with limit stated).

## Shared Contracts

| Contract | Boundary | Path | Why it exists |
| --- | --- | --- | --- |
| IC-01 Catalog paginated read + cache semantics | API + cache + batch rules | `research/catalog-read-interface-contract.md` | envelope shape, cursor, `as_of`, TTL classes, route classes, chunking, error matrix — multiple workers touch this seam |

## Milestone Strategy

| ID | Name | System change | Why this order | Path |
| --- | --- | --- | --- | --- |
| M-01 | foundation-observability | uncommitted adapter refactor integrated; pool 12; server+transport timeouts; latency/slow-query/pool-stats logging; baseline evidence (catalog COUNT, network RTT, EXPLAIN PLAN of page query) | measure before changing; M-02 cutover gated on plan evidence | `M-01-foundation-observability/` |
| M-02 | catalog-batch-cutover | `ListCatalogProductFacts` paginated port + bounded search variant (1 query/page); API envelope + OpenAPI + SDK + web adjust | biggest win; defines envelope all later work consumes | `M-02-catalog-batch-cutover/` |
| M-03 | batch-inventory-profitability-sankhya | stock batch port; cost+tax batch (chunk 500, batch route class, semaphore); SalesHistory cap 5000 (truncated marker); ImportMarginInputs ceiling 200 (422); Sankhya candidates+lines 1 JOIN; ValidateConfiguration startup-once | reuses M-02 port pattern; kills remaining N+1 | `M-03-batch-inventory-profitability-sankhya/` |
| M-04 | server-cache | TTL+singleflight cache on hot ports per ADR-03; `as_of` from cache; evict-on-mutation | needs final port shapes from M-02/M-03 | `M-04-server-cache/` |
| M-05 | web-tanstack | TanStack Query adoption; products/inventory first; mutation invalidation; force-refresh control | needs M-02 envelope; last so staleTimes mirror real TTLs | `M-05-web-tanstack/` |

Execution is parallelized in 3 waves per `execution-plan.md`: W1 = M-01; W2 = M-02 ∥ M-03 (disjoint module seams; `composition/root.go` declared conflict point, M-03 rebases); W3 = M-04 ∥ M-05 (server-only vs web-only). Gates and QA contracts unchanged; merge order fixed in the plan.

## Quality Attributes

| Attribute | Target (concrete) | Owner (ADR/seam) | Validation criterion |
| --- | --- | --- | --- |
| Q6 Maintainability | domain/application/ports/adapters/transport boundaries intact; `GOCACHE=.gocache go test ./...` green per milestone | ADR-01 | MIS-02-C03 |
| Q2 Security (credentials) | no Oracle credential/DSN/raw driver error in any log line or API error body (existing `safeOracleCause` discipline extended to new code) | adapter seam | MIS-02-C04 |
| Q7 Compatibility | godror v0.51.0 + CGO + Instant Client pinned; live lane runs via governed Docker runner | ADR-06 | M-01-C04 evidence produced by that lane |

## Non-Functional Scope

| Declined attribute | Reason |
| --- | --- |
| Q1 Performance (numeric SLO) | operator declined; validated structurally by Oracle query count per page vs baseline |
| Q3 Reliability (retry/shutdown/readyz) | operator declined Resilience block; internal tool, operator retries on blip |
| Q4 Observability as quality TARGET | instrumentation ships as M-01 tooling; not a mission-level bar beyond its criteria |
| Q5 Usability/a11y | no UI redesign in scope; only data-fetching layer changes |
| Cross-tenant cache isolation | declined with reason: Oracle read ports carry no tenant input; ERP facts are installation-global (single-tenant deployment). Recorded under Accepted assumptions with explicit revisit trigger |

## Validation Strategy

QA-2. Feature quick-validation by implementer; milestone verdicts by QA Validator against milestone contracts; mission rollup against `validation-contract.md`. Live-Oracle evidence only through the governed read-only lane (`scripts/run-live-oracle-docker.ps1`); mocks prove contract behavior. Hub-style acceptance: independent reviewer per feature before milestone accept; `merge --no-ff`; post-merge ladder (`go build`, `GOCACHE=.gocache go test ./...`, web `tsc`/build when touched).

## Risks

| id | risk | likelihood | impact | mitigation | trigger | owner |
| --- | --- | --- | --- | --- | --- | --- |
| R1 | Oracle view plans (VW_PRECO_TABELA, VW_FAT_VENDA_ITEM) make the 1-query page slow (keyset ORDER BY full-scans view) | M | H | EXPLAIN PLAN evidence in M-01 BEFORE cutover; fallback = JOIN base TGF* tables directly | M-01-C04 fails plan review | Milestone M-01/M-02 |
| R2 | IN-list >1000 → ORA-01795 | M | M | chunk at 500 in adapter contract (ADR-05) | first large import | M-03 |
| R3 | 15s deadline kills legitimate batch imports | H | M | route classes 15s/120s (ADR-04) | batch route registered as interactive | M-02/M-03 |
| R4 | batch import starves pool, interactive latency spikes | M | M | 4-permit batch semaphore (ADR-04) | import during business hours | M-03 |
| R5 | catalog text search LIKE full-scans per keystroke | M | M | FETCH FIRST 50 + client debounce; DBA index review noted | search latency >1s in slow-query log | M-02/M-05 |
| R6 | operator unaware data is cached/stale | H | M | `as_of` in every response + force-refresh (MaxAge=0) | operator confusion report | M-04/M-05 |
| R7 | link-candidates generation O(listings) degrades >2k listings | L | M | known limit; explicit trigger to batch by key IN-list | generation run >2k listings | future |
| R8 | network to Oracle is WAN not LAN (RTT 20–50ms) | L | H | M-01 measures RTT; if WAN, raise TTLs, batch becomes even more critical | M-01 baseline evidence | M-01 |
| R9 | mutation leaves L1/L2 cache stale | M | M | evict-on-mutation contract (IC-01) + TanStack queryKey invalidation | stale read after confirm/apply | M-04/M-05 |
| R10 | real catalog count ≫ assumption | L | M | M-01 COUNT via live lane feeds baseline | count >150k | M-01 |

## Handoff

- Current status: planned — readiness Ready (3 rounds, see `readiness-review.md`)
- Current owner: Mission Strategist
- Next owner: Milestone Orchestrator (M-01)
- Next action: hub-orchestrated parallel execution per `execution-plan.md` — W0 commit working tree, W1 M-01 chip, W2 M-02∥M-03 chips, W3 M-04∥M-05 chips (workers mpc-implementer/mpc-verifier, gpt-5.6-luna high; terminal callbacks to hub session)
- Required artifact paths: `mission.md`, `validation-contract.md`, `research/catalog-read-interface-contract.md`, `research/oracle-runtime-dimensioning.md`, `M-0{1..5}-*/milestone.md`, `M-0{1..5}-*/validation-contract.md`, feature briefs under each milestone
- Required evidence paths: `<feature-root>/validation.md`; `<milestone-root>/validation-result.md`; `<mission-root>/validation-result.md`
- Blocked decisions: None - all operator answers received
