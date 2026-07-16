# MIS-003-operator-cockpit-wireframe-replan

```yaml
id: MIS-003
type: mission
status: in_progress
owner: Mission Strategist
parent: none
created: 2026-07-14
updated: 2026-07-16
validation_level: QA-0
lifecycle_scope: mission
planning_phase: ready
```

## Objective

Rebuild the marketplace-central operator cockpit to the deck-2 wireframe (`research/wireframe-inventory.md`): a pt-BR, desktop-first, ERP-dense UI where an operator sees every anúncio, its sync/link/quality state, and executes provider writes exclusively through a previewed, audited mutation protocolo — on a marketplace-agnostic platform where Mercado Livre is one adapter.

Observable when done: an operator opens `/anuncios`, filters by "🔴 erro", opens the drawer for a failing listing, fixes the missing "Marca" attribute through the corrigir-atributo flow, approves the protocolo, and watches it apply — with every state transition, counter, and freshness stamp visible and honest.

## Outcome

- Canonical `listings` read model + API feeding Anúncios, Detalhe do produto, and Visão geral counters.
- One durable mutation envelope (protocolo) through which ALL provider writes route: price, stock, vínculo, pause/resync, attribute edit; listing-create reserved as contract-only.
- Wireframe-faithful workspaces: Anúncios (2a), Detalhe do produto (2b, minus Concorrência), Catálogo (1f), Vínculos (1g), Estoque (1h), Preços & Simulador (1i, minus @mercado), Visão geral (1e), Pedidos read (1j), Integrações & Sync central (1k).
- Market-data contract pinned (schema, port, honest-empty API); zero collector, zero market UI.
- MIS-001 M-13/M-14/M-07 formally superseded by this mission.

## Scope

Full-stack: new backend modules (`listings`, mutation envelope, `market` contract), OpenAPI + sdk-runtime same-commit, frontend platform seam + workspace rebuilds. Basic frontend, strictly wireframe-driven; a successor mission enhances UI (operator decision, P1b-1).

## Domain Scope

- (1) Actors & roles: single operator, actor label `operator_supplied_unverified`; no role tiers.
- (2) Core entities: produto (ERP, read), anúncio/listing (new canonical read model), vínculo, estoque (físico/reservado/disponível/seguro), pedido (read), política de preço, protocolo de mutação (new), conta/installation, market observation/reference (contract-only).
- (3) Lifecycle & states: listing sync_state (6 values, IC-02); mutation protocolo lifecycle (8 states, IC-03); link states (ADR-008, unchanged).
- (5) Classification & metadata: exception band filters (sync_error/stale/unlinked/below_margin), tabs, quality %, completude.
- (6) Audit & history: protocolo items = immutable audit with before/after; listing timeline; sync central 90d view.
- (8) Search/filter/reporting: cursor-paged lists, canonical filter grammar (IC-02), summary counters, CSV export deferred (Non-Scope).
- (9) Admin & config: installations/contas (existing), políticas de preço (existing), legacy config pages kept off main nav.
- (10) Integration & external: ML adapter (read + gated writes), Oracle/Sankhya reads (existing, governed), market collector port (contract-only).

## Non-Scope

- Pedidos faturar/etiqueta (ERP write seam) — operator declined; high governance cost, separate mission.
- Mercado UI surfaces (tela 2c, aba Concorrência 2b, colunas @mercado 2d) — deferred to successor mission behind gates G1–G7; fake seeded market data is negative-value MVP (operator, P3 review).
- Live market-data collector — same gate chain; G1 currently failed.
- Full editor wizard 1l (9 steps) / listing-create runtime — cut to corrigir-atributo mini-flow; create is contract-only (operator, P3 review).
- Monitorados (termos de busca), Perguntas ML — operator declined.
- Full multi-conta/multi-empresa — single conta functional, context designed for multi (P1b-4).
- CSV export buttons, notification bell behavior, Shopee/Amazon/Magalu connectors — wireframe decoration this mission; extend later without contract change.

## Current State

See `research/frontend-state.md` (R-02) and `research/backend-surface.md` (R-03). Key gaps: no listing entity, no mutation queue, no PriceWriter, no installation context provider, mixed EN/pt-BR copy, 5 direct-fetch pages. Foundation to consume: MIS-002 catalog cursor/search + oracle batch ports + L2 cache; M-09 canonical CODPROD + nullable facts; product_links workflows; orders read; `integration_operation_runs`; provider-write gate vocabulary in `contracts/governance/execution-lanes.json`.

## Clarified Decisions

- Resolved: see interview table + P1a/P1c menus below.
- Accepted assumptions:
  - Listing read-model staleTime class = 45s (mirrors stock class; listing rows carry published stock). Reversible: one constant in web-query + docs.
  - Protocolo item cap 2000, chunk 20, preview TTL 15 min. Reversible: config values in IC-03, no schema impact.
  - Sidebar omits "Mercado" until the successor mission (deferred UI should not present dead nav). Reversible: one Layout row.
  - `/classifications` and `/marketplaces` stay at current routes off the main nav. Reversible: redirect rows later.
  - `sales_30d` stays null this mission (no provider metrics source planned). Reversible: nullable field already in schema.
  - Dev proxy additions and Tailwind `@source`/vitest-config fixes land in M-02 (defects found in R-02). Reversible: config-only.
  - `retried_from` nullable field on MutationProtocol (IC-03) links a retry clone to its terminal source protocol. Reversible: nullable column, no consumer depends on it yet.
  - Category-attributes read mounts under `/listings` (IC-02 row granted to M-06 F-01) with L2 cache class `category_meta` (24h TTL). Reversible: one route row + one cache-class constant.
  - Sibling staleTimes pinned in ICs: `mutations` 5s, `sync` 30s, `market` 300s; protocolo polling 2s while `applying` (IC-03). Reversible: constants in web-query/IC tables, no schema impact.
- Owner decisions still open: None — all P1/P3 gates answered 2026-07-14.
- Blocked items: None — market collection is out of scope, not blocked-in-scope.

### Clarification Interview

| # | Taxon | Question | Proposed default | Operator answer |
| --- | --- | --- | --- | --- |
| 1 | build/runtime conventions | Backend scope: full-stack vs contracts+mocks? | Full-stack | Full-stack; platform must be marketplace-agnostic (canonical domain language, ML = one adapter, future adapters pluggable); YAGNI, no redundancy; basic frontend now, wireframe-faithful; successor mission enhances UI |
| 2 | lifecycle/transitions | Editor scope + SKU note interpretation | Editor in scope, SKU obrigatório | Confirmed: HUB listing writes always carry SELLER_SKU linked to produto ERP (later narrowed to mini-flow at P3 review) |
| 3 | UI convergence | Deck-2 unified as target? | Deck 2 | Confirmed; 1a–1d are source material |
| 4 | actor model | Multi-conta/empresa delivery | Single conta functional, context designed | Confirmed |
| 5 | persistence/reset | Market data acquisition | Contracts + mocked collector | Confirmed (then hardened at P3 review: contracts only, no mocked UI, test doubles in `_test` only) |
| 6 | validation expectations | MVP over-engineering check | Dispatch independent review | Done — architecture-analyst review 2026-07-14; findings accepted (defer Mercado UI, cut wizard, 6 milestones, 4 spine amendments) |

## Architecture Spine

### System Shape

Three surfaces: React SPA (`apps/web` + `packages/*`), Go modular monolith (`apps/server_core`, domain/application/ports/adapters/transport per module), PostgreSQL + Oracle-read-only + ML API behind adapters. Browser talks ONLY to sdk-runtime → server_core HTTP. Provider payloads never cross out of adapters.

### Runtime Topology

- Implementation roots: `apps/server_core` (Go), `apps/web` + `packages/{sdk-runtime,web-query,ui,feature-*}` (TS). Artifact root: `.mnfs/MIS-003-operator-cockpit-wireframe-replan/`.
- New backend modules: `listings` (read model + ingestion), mutation envelope (module home fixed in ADR-13), `market` (contract-only).
- Stores: PostgreSQL (all new tables tenant-scoped); Oracle read-only via existing internal_read; no new stores.
- Background: ONE new in-process poller (mutation applier). No external queue/bus/scheduler infra.

### Runtime Contract

- Listing facts: `listings` module owns the read model; connectors ML adapter is sole ingestion source; product_links stays sole owner of link truth (joined at read).
- Provider writes: mutation envelope is the ONLY write path (UI → `/mutations` → poller → capability adapter). Direct capability calls from transport are forbidden.
- Frontend seams: AppRouter/Layout/InstallationContext/web-query written by M-02 once; later milestones add only their assigned route rows (IC-05).
- Unknown facts stay null end-to-end; UI renders explicit unknown copy (IC-05). Never zero/default.
- OpenAPI + sdk-runtime change in the same commit (existing governance `GOV_API_SDK_SPLIT`).

### Cross-Cutting Decisions

| Decision | Status | Prevents | Must preserve | Validation impact |
| --- | --- | --- | --- | --- |
| ADR-12 Canonical `listings` module (read-only) | decided | ML-shaped transport; listing logic leaking into product_links | canonical fields, composite listing_id, nullable unknowns; ingestion via connectors capability only; freshness = `fetched_at` + manual refresh (review 9c) | M-01 criteria |
| ADR-13 Mutation envelope = protocolo table + in-process poller | decided | N ad-hoc write paths; queue-infra over-build | 7 provider-write gates structural; per-item immutable audit (review 9b: per-item rows, cap 2000); failure-code taxonomy (9d); restart-resume; StockActionService folded in | M-03 criteria |
| ADR-14 Market contract-only behind CollectorPort | decided | fabricated market facts; scraping drift | 6-signal separation; honest-empty; no production adapter; G1–G7 sequence for successor | M-06 criteria |
| ADR-15 Frontend platform seam | decided | per-page context/invalidation reinvention | IC-05 route map, redirects, query contract, 6-state vocabulary, pt-BR-only new copy | M-02 criteria |
| ADR-16 SKU invariant | decided | mislinked live listings | `listing_edit`/`listing_create` require resolved link; SELLER_SKU = CODPROD | M-03/M-06 negative criteria |
| ADR-17 Unknown semantics end-to-end | decided | zero-for-unknown regressions on new surfaces | all new fact fields nullable from first migration | all milestones |

Accepted trade-offs: ADR-12 — narrow column overlap with product_links snapshots (consolidation is a later cheap migration). ADR-13 — in-process poller limits horizontal scale; acceptable single-tenant, table is the contract for later infra swap. ADR-14 — Mercado screens absent from MVP despite being wireframe headline; honesty over demo-value. ADR-15 — legacy pages (Classifications, Marketplaces) keep direct fetch until rebuilt; recorded as migration briefs, not silent debt. ADR-16 — none. ADR-17 — none.

## Shared Contracts

| Contract | Boundary | Path | Why it exists |
| --- | --- | --- | --- |
| IC-02 Listings read | listings module ↔ SDK ↔ UI ↔ IC-03 selection | `research/listings-read-interface-contract.md` | 6 consumers of listing rows; filter grammar shared with bulk |
| IC-03 Mutation envelope | all write UIs ↔ envelope ↔ provider adapters | `research/mutation-envelope-interface-contract.md` | 5 write surfaces, one lifecycle, 7 gates |
| IC-04 Market data | future collectors ↔ market module ↔ future UI | `research/market-data-interface-contract.md` | pin semantics before any consumer exists |
| IC-05 Frontend platform | AppRouter/context/web-query ↔ all pages | `research/frontend-platform-interface-contract.md` | routes, keys, invalidation crosswalk, state copy |
| IC-01 (MIS-002) catalog read + dormant product-edit row | catalog ↔ UI | `.mnfs/MIS-002-oracle-read-rearchitecture/research/catalog-read-interface-contract.md` | dormant row ACTIVATED by M-04 enrichment edit |

## Milestone Strategy

| ID | Name | System change | Why this order | Path |
| --- | --- | --- | --- | --- |
| M-01 | listings-read-spine | New `listings` module: schema, ML ingestion, GET /listings + by-product + summary + refresh, OpenAPI+SDK | Everything downstream consumes listing rows; filter grammar must exist before M-02/M-03 | `M-01-listings-read-spine/milestone.md` |
| M-02 | frontend-platform-anuncios | Platform seam (routes, redirects, context, web-query expansion, state components) + Anúncios workspace 2a read | Seam needs a real consumer to be proven; Anúncios is the flagship read | `M-02-frontend-platform-anuncios/milestone.md` |
| M-03 | mutation-envelope-writes | Protocolo tables + poller + price/stock/link/pause/resync writes + bulk + preview UI + protocolo detail page | Depends on M-01 selection grammar; own dual gate (provider-write surface) | `M-03-mutation-envelope-writes/milestone.md` |
| M-04 | read-workspaces-catalogo-produto-estoque-vinculos | Detalhe do produto 2b, Catálogo 1f, IC-01 activation, Estoque 1h, Vínculos 1g rebuilds | Read screens over existing modules + M-02 seam; Estoque/Vínculos write actions consume M-03 | `M-04-read-workspaces/milestone.md` |
| M-05 | visao-geral-pedidos-sync-central | Summary aggregates endpoint, Visão geral 1e, Pedidos read 1j, Integrações & sync central 1k | Aggregates need M-01 counters + M-03 protocolo list; closes the read loop | `M-05-visao-geral-pedidos-sync-central/milestone.md` |
| M-06 | corrigir-atributo-market-contracts | ML category-attribute read + corrigir-atributo flow via envelope (`listing_edit`) + `market` module contract-only | Needs M-03 envelope + M-04 product detail as host surface | `M-06-corrigir-atributo-market-contracts/milestone.md` |

Dependencies: M-01 → {M-02, M-03}; M-02 → {M-04, M-05}; M-03 → {M-04 (write actions), M-05 (protocolo list), M-06}; M-04 → M-05 (AppRouter/nav seam — one milestone writer at a time); M-04 → M-06 (host surface). M-01 ∥ nothing; M-02 ∥ M-03 after M-01.

## Parallel Execution Plan

Replanned 2026-07-16 at the M-01 boundary (M-01 `passed`). Feature-grain refinement of the
milestone DAG above — where this plan and the milestone-grain line conflict, this plan wins
(core §0). Binding for the hub board and every chip prompt.

### Waves

| Wave | Chips (concurrent) | Contents | Starts when |
| --- | --- | --- | --- |
| W1 | CHIP-M02 ∥ CHIP-M03 ∥ CHIP-SAT | M-02 (all) · M-03 (all; F-04 FE bits gated, see below) · SAT = M-05 F-01 + M-06 F-02 | now (M-01 passed) |
| W2 | CHIP-M04 | M-04 (all) | M-02 CLOSED + M-03 CLOSED (merged) |
| W3 | CHIP-M05 ∥ CHIP-M06 | M-05 F-02/F-03/F-04 · M-06 F-01 | M-04 CLOSED; M-05 FE also needs SAT's M-05 F-01 merged; M-06 F-01 needs SAT's M-06 F-02 merged + M-03 |

### Feature-grain DAG refinements (deltas vs milestone-grain line)

- **M-05 F-01 depends only on M-01** (aggregate/orders/sync-runs endpoints, zero frontend).
  Pulled forward into W1 SAT chip. The M-02 dependency applies only to M-05 F-02..F-04.
- **M-06 F-02 depends only on M-01** (market contract module, no adapter, no UI; feature.md
  already marks it parallel-eligible). Pulled forward into W1 SAT chip. M-03/M-04 dependencies
  apply only to M-06 F-01.
- **M-03 F-04 ← M-02 F-03** (undeclared edge, now declared): preview/confirm modal mounts in
  the Anúncios workspace M-02 F-03 builds. M-03 F-01..F-03 (backend) proceed regardless; F-04
  FE work starts only after the hub confirms M-02 F-03 merged and triggers CHIP-M03's rebase.
- **M-06 internal order corrected**: F-02 runs FIRST (W1, inside SAT), F-01 later (W3, rebases
  on merged F-02). The one-writer concern the old F-01→F-02 order guarded is handled by
  disjoint OpenAPI sections + the additive composition-root lock below.

### Shared-seam pre-assignments (W1 concurrency contract)

- **OpenAPI (`contracts/api/marketplace-central.openapi.yaml`) + `packages/sdk-runtime`** —
  disjoint sections, additive-only: CHIP-M03 = mutation/protocolo paths + their schemas;
  CHIP-SAT(M-05 F-01) = dashboard-summary, orders, sync-runs paths; CHIP-SAT(M-06 F-02) =
  market + category-attribute paths. No chip touches another's section; shared preamble/info
  blocks are hub-owned.
- **Migration blocks** (existing max = 0037): CHIP-M03 = **0038–0042**; CHIP-SAT(M-06 F-02) =
  **0043–0045**. M-05 F-01 gets no block (reads existing tables); if its planner finds a new
  table is required → `REQUEST` to hub for a block, never self-assign.
- **Additive contract-locks pre-granted** (core §5 mechanism; released at each CLOSED, diffs
  called out in the event): (1) server composition root — module registration lines only:
  CHIP-M03 registers mutations module, CHIP-SAT registers dashboard/orders/sync + market
  modules; (2) `connectors` PriceWriter/StockWriter wiring — CHIP-M03 F-02 only.
- **Frontend platform seam (AppRouter/nav/web-query/state components)** — CHIP-M02 exclusive
  in W1. CHIP-M04 exclusive in W2. W3: M-05 FE vs M-06 F-01 route rows are disjoint but the
  files are shared — hub serializes merges (rebase-then-merge, one at a time).
- **`GET /orders` contract-satisfiability directive** (binding on M-05 F-01 planner): OpenAPI
  already defines `listMarketplaceOrders` (limit-based) on the same path. Evolve the existing
  operation in place, additive-only (cursor params added alongside, existing params/response
  fields preserved). Never author a duplicate path/operationId. If cursor semantics cannot be
  added without breaking the existing contract → `ESCALATION`, not a workaround.

### Bookkeeping

- CHIP-SAT closes at **feature** grain: it reports `CLOSED` for M-05 F-01 and M-06 F-02 (own
  evidence per feature); milestones M-05/M-06 stay `in_progress`/open until their W3 chips
  close, then dual gate + QA at milestone grain per core §6.
- Each chip pins the per-milestone governance base anchor (profile §2) at dispatch; drift is
  measured against that 40-hex SHA, not `main`-at-merge-time.
- Known-failure allowlist (profile §2) applies: chips cite, never re-prove, listed failures.

## Quality Attributes

| Attribute | Target (concrete) | Owner (ADR/seam) | Validation criterion |
| --- | --- | --- | --- |
| Q1 Performance | `GET /listings` p95 < 500ms at 2k seeded rows, limit 50; summary endpoint single query; UI first data paint < 2s dev-local | ADR-12 / IC-02 | M-01 VC criterion with measured evidence (API targets) + M02-C11 (UI first data paint) |
| Q2 Security | tenant_id scoping on every new table/query; no token/PII in logs/protocolo payloads/evidence; STRIDE: tampering → envelope state machine rejects out-of-order transitions (409); repudiation → actor label + immutable items (label is `operator_supplied_unverified` per ADR-009 — it records the claim durably, it does not authenticate it); spoofing/elevation-of-privilege on the write path → explicitly declined this mission, see Non-Functional Scope; info-disclosure → `message_provider` sanitized (no tokens) | ADR-13 / IC-03 | M-03 VC negative criteria per gate |
| Q3 Reliability | poller resumes approved/applying after restart; idempotent re-apply proven; failed preserved never auto-mutated | ADR-13 | M-03 VC restart + idempotency criteria |
| Q4 Observability | every protocolo item carries failure code + pt + provider message; sync central lists runs + protocolos with erro traduzido; listing timeline last 10 events | IC-03 / M-05 | M-03/M-05 VC render/API criteria |
| Q5 Usability | pt-BR-only new copy per IC-05 glossary; 6-state vocabulary components used on every rebuilt page; deep links survive reload | ADR-15 / IC-05 | M-02 VC reload/copy criteria; per-workspace criteria |
| Q6 Maintainability | module layering per governance `modules.json`; OpenAPI+SDK same-commit; TanStack-only server state in rebuilt pages; no provider payload beyond adapters | existing governance + ADR-15 | governance lane green in every milestone VC |

## Non-Functional Scope

| Declined attribute | Reason |
| --- | --- |
| Q7 Compatibility/portability | Chrome-latest desktop only this mission (operator, P1c); wireframe is desktop-first denso |
| Caller authentication on the write path (`POST /mutations` + preview/approve) | Declined with reason: single trusted operator, no role tiers (Domain Scope, ADR-009), dev-local same-origin deployment with no external network exposure (IC-02 Transport: same-origin only; cross-origin requires new ADR). `actor` stays client-supplied `operator_supplied_unverified`; server-side identity verification (session/service auth binding write execution to a verified caller) is deferred to a successor mission before any internet-exposed deployment. Until then the deployment boundary (localhost/private network) is the authentication perimeter. |

## Validation Strategy

- Per feature: `<feature-root>/validation.md` (quick proof: tests run, endpoints exercised, screens rendered).
- Per milestone: `<milestone-root>/validation-contract.md` with stable criteria IDs (`M0n-C0m`); QA rollup at `<milestone-root>/validation-result.md`; dual gate = Codex `mpc-verifier` (gpt-5.6-sol, medium) + independent Claude review at fixed SHA; only QA passes.
- Mission: `validation-contract.md` → `validation-result.md` at mission root.
- Lanes: unit + integration (ephemeral-postgres) prove contracts; live-provider-read lane only where a criterion explicitly says so; provider-write live execution only under the governed provider-write lane with production-like evidence — M-03 validates lifecycle against a stub capability adapter in integration lane; any LIVE ML write needs explicit operator authorization at execution time.
- Mocks never claim live integration (ADR-010 carried forward).
- Browser evidence: manual/scripted walkthrough with screenshots to evidence paths (browser lane remains unconfigured; do not invent a Playwright stack this mission).

## Risks

| id | risk | likelihood | impact | mitigation | trigger | owner |
| --- | --- | --- | --- | --- | --- | --- |
| RK-01 | ML rate limits during bulk apply (unmeasured, G6 pending) | M | M | chunk 20 + sequential per installation + `provider_rate_limited` retryable code; start conservative | 429s in protocolo failures | M-03 implementer |
| RK-02 | Listing ingestion volume/pagination surprises (1.3k listings) | M | M | full-page pull via existing capability; operation-run failure honest; refresh 409 guard | refresh runs failing | M-01 implementer |
| RK-03 | Live provider writes cannot be safely exercised in validation | H | M | validate lifecycle on stub adapter in integration lane; live write = explicit operator-authorized provider-write lane run, optional for milestone pass | QA gate | QA Validator |
| RK-04 | Seam contention: M-02 files touched by later milestones | M | H | IC-05 route-row ownership; one writer per seam; milestone briefs name allowed paths | overlapping diffs | Milestone Orchestrator |
| RK-05 | Legacy direct-fetch pages regress under new router | L | M | redirects tested; untouched pages keep routes; migration briefs explicit | route smoke failures | M-02 implementer |
| RK-06 | G1 OAuth defect (StartAuthorize clobbers connected) trips refresh/ingestion flows | M | M | MIS-003 touches only read + gated writes on an already-connected installation; defect fix stays successor-mission scope unless it blocks ingestion — then correction-scoped | ingestion auth failures | Milestone Orchestrator |
| RK-07 | Scope creep back toward Mercado UI ("just render the tab") | M | M | Non-Scope explicit; IC-04 has no UI consumer; readiness review checks | PR adding market UI | Mission Strategist |

## Research Links

- R-01 `research/wireframe-inventory.md` — screen truth.
- R-02 `research/frontend-state.md` — frontend facts + passing defects.
- R-03 `research/backend-surface.md` — gap matrix.
- R-04 `research/market-intelligence-digest.md` — market gates + still-binding MIS-001 ADRs.

## Migration Briefs

Pages remaining on direct fetch after MIS-003, each to be migrated to TanStack-via-web-query when next touched (brief = one paragraph here, not a milestone):

- **Classifications** (`packages/feature-classifications`): swap `listCatalogProducts` deprecated alias → `listCatalogProductFacts`; wrap list in `catalogQueryKeys`-style keys under `catalog` namespace; mutations invalidate `catalog`.
- **Marketplaces** (`packages/feature-marketplaces`): reads → new `marketplaces` namespace keys; staleTime `pricecost`.
- **Simulator** (rebuilt in M-04 as `/precos`): migrates then; deprecated aliases retired there.
- **Dashboard** (replaced by Visão geral in M-05): delete, don't migrate.
- **Integrations** (rebuilt in M-05): migrates then under `sync` namespace.
Deprecated SDK aliases `listCatalogProducts`/`searchCatalogProducts` are removed in the same milestone that removes their last consumer (M-04/M-05 — whichever lands second).

## Supersession

MIS-001 milestones M-13, M-14, M-07 are superseded by MIS-003 (operator decision 2026-07-14). Their briefs are source material; carried-forward invariants are recorded in R-04. M-06 stays failed-preserved (ADR-011); M-09 passed (frozen SHA 2eabecbc); M-10/11/12 stay deferred. On apply of this mission, each superseded milestone.md receives a `status: superseded_by: MIS-003` frontmatter note — no other MIS-001 artifact is edited.

## Handoff

- Current status: planning artifacts authored; readiness review pending (P7).
- Current owner: Mission Strategist.
- Next owner: Milestone Orchestrator (M-01) after readiness `planned`.
- Next action: run P7 readiness crew; then operator starts M-01 via the goal-harness handoff.
- Required artifact paths: `mission.md`, `validation-contract.md`, `architecture-map.md`, `research/*.md` (R-01..04, IC-02..05), `M-0n-*/milestone.md` + `validation-contract.md` + `F-nn-*/feature.md` (six milestones).
- Required evidence paths: `<feature-root>/validation.md`; `<milestone-root>/validation-result.md`; `<mission-root>/validation-result.md`.
- Blocked decisions: None — all gates answered.
