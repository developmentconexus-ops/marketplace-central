# P3 Sol Counter-proposal — r01 (VERBATIM)

> Dispatch: gpt-5.6-sol / model_reasoning_effort=medium / codex exec OS-process / sandbox read-only.
> Prompt: scratchpad prompt-p3-sol-r01.md · log: agent__p3-sol-r01.log · exit 0 · 2026-07-17.
> Input: frozen manifest p3-input-r01.sha256 (top-digest c95678eb2819a8a4de3aaedfb540f68b87b5c90111259de279c5e77d3899b1d7). BLIND — Claude candidate denied.
> Conteúdo abaixo = stdout final do Sol, sem edição.

---
## Outcome

**MIS-004 — mvp-demo.** Deliver a local Docker demo for Monday, 2026-07-20, that imports the client’s ERP stock from `.xlsx`, preserves `CODPROD`/EAN/`REFFORN` identity honestly, links products to the connected seller account, and carries the user through stock → market evidence → “vale a pena vender?” → real margin simulation. It also delivers the approved rethemed shell, functional read-only Pedidos, minimal DIFAL visibility, and live Mercado Livre reads while making every provider write impossible in the demo configuration. Pre-listing price absence is a valid visible outcome, never fabricated evidence. (`REPLAN-BRIEF-2026-07-17.md` §§3–5; `p1-clarified-decisions-2026-07-17.md` P1a–P1c)

**MIS-005 — produto-completo.** Productize the remaining platform: production write enablement, multi-tenant authentication, webhook-driven convergence, scheduled market history, full Mercado radar, complete product detail, claims/returns, fiscal workflows, configuration, reconciliation, Full inventory, and any external price provider only after contractual and empirical homologation. MIS-005 receives milestone headlines now and must pass its own P4–P7 planning before execution. (`REPLAN-BRIEF-2026-07-17.md` §4; `p1-clarified-decisions-2026-07-17.md` Accepted assumptions)

## Architecture spine

### ADR-01: Dual ERP source behind one canonical import boundary

- **Decision:** Keep the existing Oracle/Sankhya reader as an alternate source, but make a versioned `.xlsx` import through the `product_links` import workflow the demo path. Both sources normalize into the same catalog-owned product representation. Required spreadsheet columns are `CODPROD`, `DESCRPROD`, `CUSTO`, and `ESTOQUE_FISICO`; optional columns are `ESTOQUE_RESERVADO`, `EAN`, `REFFORN`, `MARCA`, and `NCM`. Persist import protocol, file hash, source, and import time; do not depend on the workbook remaining mounted at runtime. (`p1-clarified-decisions-2026-07-17.md` P1a and Accepted assumptions; `repo-baseline-2026-07-17.md` ERP read path)
- **Prevents:** One worker embedding spreadsheet rules in the UI while another builds an incompatible Oracle-only domain.
- **Must preserve:** Missing optional values remain unknown; malformed required values reject the affected row with inspectable reasons. Oracle unavailability remains an error, not zero stock/cost. (`repo-baseline-2026-07-17.md` ERP read path; `HARNESS-PROFILE.md` §7)
- **Trade-off:** The demo adds an ingestion path that later needs lifecycle and retention hardening.
- **Validation impact:** Use the real client-format workbook or an operator-approved non-sensitive fixture; prove row counts, rejection reasons, import protocol, and identical canonical identity semantics across `.xlsx` and Oracle contract tests.

### ADR-02: Canonical identity belongs to `catalog`

- **Decision:** `catalog` owns `ProductIdentity`, keyed by canonical internal SKU `CODPROD`. `REFERENCIA` maps to EAN/GTIN and `REFFORN` to manufacturer reference. Correct the existing Oracle reader contract that currently treats `REFERENCIA` as a generic manufacturer/reference value. `seller_sku` may resolve only to `CODPROD`, never to EAN or `REFFORN`. (`pricing-intelligence-implementation-handoff.md` §4; `repo-baseline-2026-07-17.md` ERP read path)
- **Prevents:** Cross-source identity corruption and false links caused by the baseline reader defect.
- **Must preserve:** GTIN checksum and uniqueness checks; at least two independent anchors for automatic acceptance; hard contradictions override EAN. Fuzzy/title matching ranks candidates but cannot auto-accept them. (`pricing-intelligence-implementation-handoff.md` §5)
- **Trade-off:** More imports remain in review instead of producing an apparently complete demo.
- **Validation impact:** Include known EAN-collision, duplicate-GTIN, kit, variant, and missing-EAN cases. Missing EAN permits `REVIEW`/`NO_CANDIDATE`, not title-only `ACCEPT`.

### ADR-03: Links and imports are workflows, not catalog identity

- **Decision:** `product_links` owns import runs, candidate generation, confidence/reasons, manual resolution, undo, batch preview, and protocol. It references catalog identities but does not redefine them. Linking is an MPC-local mutation and never alters Mercado Livre. (`API-MAP.md` “Vinculos e Importacao”; `HANDOFF.md` “Telas prontas”; `repo-baseline-2026-07-17.md` module inventory)
- **Prevents:** Catalog, listing, and UI workers each creating separate link/confidence models.
- **Must preserve:** Provider payloads end at adapters; persisted workflow states are tenant-scoped.
- **Trade-off:** Import and identity become separate concepts that must be joined explicitly.
- **Validation impact:** Demonstrate confirm/review/no-match paths, manual selection, undo, batch protocol, and a visible confirmation that nothing changed in ML.

### ADR-04: Price evidence lives in the existing `market` module

- **Decision:** Expand the contract-only `market` module rather than create a parallel pricing-intelligence module. It owns `MarketPriceSnapshot`, `CompetitiveSignal`, `ValidatedOffer`, and `MarketPriceAggregate`; `ProductIdentity` remains in `catalog`. `pricing` consumes a published market read contract and must not persist its own ambiguous `market_price`. (`repo-baseline-2026-07-17.md` Modules backend; `pricing-intelligence-implementation-handoff.md` §6)
- **Prevents:** Incompatible market-price meanings and duplicate snapshot tables.
- **Must preserve:** Source, match reasons, contradictions, observation/fetch/expiry times, request IDs, sample size, and currency/condition evidence; every query carries `tenant_id`.
- **Trade-off:** Simulator delivery depends on the market contract even when evidence is absent.
- **Validation impact:** Contract tests must distinguish our sale price, winner price, competitive target, and catalog aggregates; prove tenant isolation and ADR-17 failure preservation.

### ADR-05: One owner extends the Mercado Livre read adapter

- **Decision:** A single MIS-004 milestone owns all changes to `connectors/adapters/mercado_livre/capability_adapter.go` and the corresponding ports: own-item `sale_price`, `price_to_win`, catalog discovery/detail, optional catalog offers, shipments/costs/delays, and shipping options. Market and orders consume normalized MPC-native ports. (`repo-baseline-2026-07-17.md` Adapter Mercado Livre; `HARNESS-CORE.md` §3)
- **Prevents:** Concurrent edits to the shared adapter and leakage of ML DTOs into `market` or `orders`.
- **Must preserve:** Existing OAuth credential resolution; provider payloads die at the adapter; `/products/{id}/items` is default-off, fully paginated when enabled, emits telemetry, and returns an explicit unavailable/fallback state. No `/sites/MLB/search`, competitor-item dependency, scraping, proxy, or unhomologated provider. (`pricing-intelligence-implementation-handoff.md` §§2.2, 3, 8)
- **Trade-off:** Orders waits on the shared adapter owner for shipment-related reads.
- **Validation impact:** Real connected-account read evidence is required; mocks prove normalized port shape only.

### ADR-06: Identity state and commercial verdict remain distinct

- **Decision:** The shared market/identity state vocabulary is `ACCEPT`, `REVIEW`, `REJECT`, `NO_CANDIDATE`, `NO_PRICE_EVIDENCE`, and `INSUFFICIENT_MARKET`. “Vale a pena vender?” is calculated only from an accepted identity plus sufficient explicit cost/fee/freight/tax/market inputs. Otherwise it displays the blocking evidence state, not a commercial verdict. Fewer than five valid sellers is `INSUFFICIENT_MARKET`. (`pricing-intelligence-implementation-handoff.md` §§1, 5, 7)
- **Prevents:** Treating a matched product as proof of price or presenting null `buy_box_winner` as zero.
- **Must preserve:** Source, collection age, `n_offers`, `n_sellers`, match confidence, and limitations appear on every price UI. `price_to_win` is an ML competitive target, not “lowest market price.”
- **Trade-off:** Some client products may show no numeric comparison during the demo.
- **Validation impact:** Browser QA must exercise all empty/uncertain states, including `buy_box_winner=null`, and confirm that no `R$0`, zero margin, or misleading green verdict appears.

### ADR-07: Snapshot validity follows ADR-17

- **Decision:** Successful observations append evidence; network errors, null provider fields, expiration, and insufficient samples record status without overwriting the latest valid snapshot with zero/default. On-demand preparation is the only MIS-004 collection trigger; daily scheduling belongs to MIS-005. (`pricing-intelligence-implementation-handoff.md` §§2.3, 6; `p1-clarified-decisions-2026-07-17.md` P1c and Accepted assumptions)
- **Prevents:** A transient provider failure corrupting the demo’s last known evidence.
- **Must preserve:** `observed_at`, `fetched_at`, expiry, raw status, source, and staleness remain distinguishable.
- **Trade-off:** The UI must explain stale-but-valid evidence separately from a current failed fetch.
- **Validation impact:** Negative persistence test: valid snapshot → failed/null collection → valid snapshot remains intact and visibly aged, with the failed attempt inspectable.

### ADR-08: Retheme-first frontend convergence

- **Decision:** The first frontend milestone owns the approved paper-and-green shell, typography, light/dark theme, nav, responsive table primitives, drawers, shared chips, `Layout`, and `AppRouter` route registration. Subsequent UI milestones own disjoint route subtrees and reuse those primitives. Mercado and Repasses remain visible but disabled as “em breve”; Vínculos stays outside global nav. (`p1-clarified-decisions-2026-07-17.md` P1b; `README.md` Design Tokens; `HANDOFF.md` Design decisions and Nav; `design-screens-2026-07-17.md` Shell/nav)
- **Prevents:** Every screen independently rebuilding shell/theme/nav and colliding in `Layout` or `AppRouter`.
- **Must preserve:** No inline table editing; canonical nav; final PT-BR copy; Instrument Sans/IBM Plex Mono; no fake counts from mocks. (`HANDOFF.md` Design decisions; `design-screens-2026-07-17.md` Dashboard and Shell/nav)
- **Trade-off:** Feature pages cannot close before the shared shell lands.
- **Validation impact:** Visual QA first validates the shell in both themes and narrow/table-overflow states; route milestones then validate only their owned pages plus shared-shell regression.

### ADR-09: Demo mutations use M-03 but provider execution is unreachable

- **Decision:** Execution starts only after M-03 merges. Import/link and any other user-triggered mutation use its preview/protocol envelope. Price/create/pause/faturar controls may show preview and protocol, but the local demo configuration has no executable Mercado Livre write path; zero live provider writes are permitted. (`p1-clarified-decisions-2026-07-17.md` P1b; `REPLAN-BRIEF-2026-07-17.md` §§4–5; `repo-baseline-2026-07-17.md` M-03 mutation envelope)
- **Prevents:** A demo click bypassing governance and modifying the connected seller account.
- **Must preserve:** Any future provider execution requires resolved linkage, explicit policy/source time, duplicate protection, and audit through the M-03 envelope. (`HARNESS-PROFILE.md` §7)
- **Trade-off:** “Aplicar” demonstrates the governed intent/protocol, not remote effect.
- **Validation impact:** Network/audit proof must show no `PUT`/`POST` reaches ML; direct calls outside `/mutations` fail; protocol and preview remain inspectable.

### ADR-10: Polling and explicit refresh, no webhooks in MIS-004

- **Decision:** Use authenticated GET reconciliation, on-demand pre-demo collection, and bounded UI polling/refetch. Do not register or depend on webhooks or a daily scheduler in MIS-004. (`p1-clarified-decisions-2026-07-17.md` P1a; `REPLAN-BRIEF-2026-07-17.md` §4)
- **Prevents:** The three-day demo depending on callback infrastructure and asynchronous event convergence.
- **Must preserve:** Every displayed read carries source time/freshness; polling failure produces an honest stale/error state.
- **Trade-off:** Changes may appear later than they would under webhooks.
- **Validation impact:** Live-read QA proves refresh and stale/error behavior; webhook delivery is explicitly not a pass criterion.

### ADR-11: Minimal DIFAL and calculation configuration live in `pricing`

- **Decision:** `pricing` owns a tenant-scoped calculation profile and the versioned 27-UF DIFAL reference/override table. Resolution follows global tenant parameters → product cost/weight/overrides → listing-specific live commission/freight. The simulator owns the DIFAL toggle; orders consumes a read-only computed chip/column. No scheduling, reminders, or paid state in MIS-004. (`HANDOFF.md` Design decisions; `design-screens-2026-07-17.md` Simulador and Configurações; `p1-clarified-decisions-2026-07-17.md` P1a)
- **Prevents:** Simulator and orders embedding conflicting UF tables or percentages.
- **Must preserve:** The mock’s hardcoded SP behavior is forbidden. Missing origin/destination/rate produces unknown/unavailable, not 0%. Commission and freight come from live ML reads; they are never hardcoded. (`design-screens-2026-07-17.md` Simulador; `API-MAP.md` Simulador)
- **Trade-off:** MIS-004 exposes only the calculation drawer and read-only fiscal indication, not complete configuration/fiscal operations.
- **Validation impact:** Prove one source of truth across simulator and orders, toggle on/off recalculation, destination-sensitive behavior, and unknown-rate handling.

### ADR-12: Local Docker is the single demo runtime

- **Decision:** The supported demo topology is `npm run docker:dev`, backend `:8080`, frontend `:5174`, local Postgres, `.xlsx` import, and the existing encrypted ML installation. No cloud deployment is required for MIS-004. (`p1-clarified-decisions-2026-07-17.md` P1b; `repo-baseline-2026-07-17.md` Tenancy/auth)
- **Prevents:** Workers validating against incompatible ports, databases, or remote environments.
- **Must preserve:** Secrets remain environment-provided; no credentials enter fixtures or evidence. The current fixed tenant is explicit demo debt, not implicit multi-tenant security.
- **Trade-off:** Client access depends on the operator’s machine and connected environment.
- **Validation impact:** Every milestone still follows L0–L4; final QA starts from a clean local stack and exercises the real workbook and real read-only provider seam. (`HARNESS-CORE.md` §5)

### ADR-13: Contract-first lock and hand-written SDK atomicity

- **Decision:** One contract milestone owns all MIS-004 OpenAPI sections and `packages/sdk-runtime`; it edits the hand-written SDK in the same commit as OpenAPI. Implementation milestones may not independently alter either surface; contract correction requires a lock request back to that owner. (`repo-baseline-2026-07-17.md` OpenAPI and sdk-runtime; `HARNESS-PROFILE.md` §§5, 7)
- **Prevents:** Multiple chips colliding in the shared spec/client and assuming nonexistent code generation.
- **Must preserve:** Path/schema/client parity and tenant-aware request shapes.
- **Trade-off:** Contract defects can serialize otherwise independent work.
- **Validation impact:** Contract/SDK/handler parity is checked at L0–L2; a spec-only or SDK-only commit fails the milestone.

## MIS-004 milestone split

Execution precondition: all three W1 chips—M-02 frontend platform, M-03 mutation envelope, and SAT contracts—must be merged, accepted, and used as the common branch point. They were not merged at baseline `cd74b401`; planning must not infer their eventual diff. (`repo-baseline-2026-07-17.md` Git / estado W1)

| ID / slug | Headline | Owned surfaces | Dependencies |
|---|---|---|---|
| **M004-01 / `mvp-contract-lock`** | Freeze the MVP boundary contracts before parallel implementation | All MIS-004 additions in `contracts/api/marketplace-central.openapi.yaml`; matching hand-written `packages/sdk-runtime`; file-format contract for `.xlsx`; shared enums and null/error semantics. OpenAPI sections: product-link imports/resolutions, catalog identity, market evidence/signals, pricing profile/simulations, orders projections, dashboard summary. No migration. | Accepted W1 contract state. Must close before API consumers close; it may run parallel with shell work. |
| **M004-02 / `retheme-shell`** | Land the approved shell and reusable visual primitives first on the FE track | `apps/web` shared theme/tokens, `Layout`, canonical nav, `AppRouter`, route placeholders, query/error/empty-state frame, shared tables/drawers/chips. FE seams: all shell-level files; no business page ownership. No migration. | M-02 W1 merge. First FE milestone; all page milestones depend on it. |
| **M004-03 / `xlsx-identity-links`** | Import client stock and produce governed product-link candidates | Modules: `catalog`, `internal_read`, `product_links`; `.xlsx` adapter and Oracle identity correction; import/candidate/resolution persistence. OpenAPI implementation only against M004-01 sections. FE route ownership: none. Migration block **0045–0049**. | M004-01 identity/import contract; M-03 envelope for local mutations. Supplies canonical products, cost, stock, and links to market, simulator, orders, and product UI. |
| **M004-04 / `official-market-evidence`** | Collect official read-only competitive signals and honest pre-listing evidence | Modules: exclusive `connectors` ML capability-adapter owner plus `market`; normalized ports; snapshots/signals/offers/aggregates. Owns all shared ML adapter extensions, including shipment/shipping reads needed later by orders/pricing. OpenAPI implementation against market sections. Migration block **0050–0054**. | M004-01 market contract and identity vocabulary; connected ML installation. Can execute alongside M004-03, but integrated acceptance uses its canonical identities. |
| **M004-05 / `simulator-difal-service`** | Produce real price↔margin calculations from ERP and ML evidence | Modules: `pricing` and published reads from `profitability`; calculation profile, DIFAL seed/overrides, decomposition, scenarios, verdict prerequisites. No connector edits. OpenAPI pricing/config sections. Migration block **0055–0059**. | M004-03 cost/stock identity; M004-04 market, fee, freight, and competitive-signal ports; M004-01 schemas. |
| **M004-06 / `orders-read-model`** | Deliver functional Pedidos data and honest profitability/DIFAL projections | Module: `orders`; order/shipment projection, SLA queue, list/Kanban state, masked buyer detail, timeline, tracking, decomposition, DIFAL read projection. No connector or pricing-table ownership. OpenAPI orders sections. Migration block **0060–0064**. | M004-03 product/link/cost data; M004-04 normalized order/shipment reads; M004-05 calculation/DIFAL read contract. |
| **M004-07 / `inventory-to-verdict-ui`** | Complete the demo’s ERP-stock-to-verdict frontend journey | FE route trees: Vínculos/Importação, Anúncios retheme, Produto Detalhe header/verdict/linked listings/stock. Route-local components only; no shell/AppRouter edits. No migration. | M004-02 shell and M004-01 SDK. May build concurrently with backend, but closes only against M004-03 and M004-04 real data. |
| **M004-08 / `simulator-ui`** | Deliver the approved simulator matrix, panel, parameters, and scenarios | FE seam: `/simulator` route tree only; bidirectional price/margin interaction, evidence states, DIFAL toggle, scenario UI, disabled provider-execution preview. No migration. | M004-02 shell and M004-01 SDK; closes against M004-05 and M-03 preview/protocol. |
| **M004-09 / `orders-ui`** | Deliver functional Fila, Lista, read-only Kanban, and order drawer | FE seam: `/orders` route tree only; KPIs, filters, SLA/DIFAL ordering, status tabs, read-only Kanban, masked detail drawer. Cancelados/Devoluções visibly deferred. No migration. | M004-02 shell and M004-01 SDK; closes against M004-06. |
| **M004-10 / `dashboard-cuttable`** | Build an ML-only dashboard from already available aggregates | Module: `dashboard`; FE route `/`; no new provider calls or snapshot model. No planned migration. | M004-03 through M004-06 data. Last priority and explicitly cuttable without breaking the central story. (`p1-clarified-decisions-2026-07-17.md` P1a) |
| **M004-11 / `demo-integration`** | Assemble and prove the complete local read-only demo | Exclusive composition-root/module-registration integration, Docker demo configuration, route convergence, real `.xlsx` import run, real provider-read evidence, provider-write denial, and final story rehearsal. No new domain scope. Migration block **0065–0069 reserved only for approved integration correction**, not self-assigned. | All non-cut milestones merged; M004-10 optional. Full dual gate and fresh QA live-drive remain mandatory. |

### Proposed waves and parallelism

1. **Wave 0 — converge the base:** merge and validate W1. No MIS-004 chip branches from baseline `cd74b401`.
2. **Wave 1 — freeze seams:** M004-01 and M004-02 run in parallel: contract artifacts versus frontend shell.
3. **Wave 2 — independent data foundations:** M004-03 and M004-04 run in parallel after the contract lock. They own disjoint modules and migration blocks; M004-04 alone owns the shared ML adapter.
4. **Wave 3 — backend consumers plus route-local frontend:** M004-05 follows the imported identity and market read contracts. M004-06 may begin contract-first but cannot close before M004-03–05. M004-07, M004-08, and M004-09 may run concurrently because their route trees are disjoint; their close gates wait for their real backend dependencies.
5. **Wave 4 — optional aggregate:** M004-10 runs only if the central story is already green.
6. **Wave 5 — integration:** merge smallest/lowest-risk accepted chips first, then M004-11 owns composition wiring and the integrated clean-stack ladder. The shared test database is serialized whenever migration fingerprints differ. (`HARNESS-CORE.md` §§3, 5)

### Collision pre-assignments

| Seam | Pre-assignment |
|---|---|
| OpenAPI + SDK | Exclusive M004-01 contract lock. OpenAPI and hand-written SDK land together. Other milestones request corrections; they do not edit these surfaces independently. |
| Migration numbers | M004-03: 0045–0049; M004-04: 0050–0054; M004-05: 0055–0059; M004-06: 0060–0064; M004-11 contingency: 0065–0069. Every new migration also updates the runner’s hardcoded count. (`repo-baseline-2026-07-17.md` Migrations) |
| ML capability adapter | Exclusive M004-04 owner for `capability_adapter.go` and normalized ML read ports. Orders/pricing consume ports; no additive edits to this adapter from their chips. |
| DB tables | `catalog` identity and import-source records: M004-03; market snapshots/signals/offers/aggregates: M004-04; calculation/DIFAL tables: M004-05; order projections: M004-06. Cross-module reads use published interfaces. |
| FE shell | M004-02 exclusively owns `Layout`, nav, `AppRouter`, theme, and shared visual primitives. Route placeholders are registered there. |
| FE routes | M004-07 owns Vínculos/Anúncios/Produto; M004-08 owns Simulador; M004-09 owns Pedidos; M004-10 owns Dashboard. No route chip edits the shell seam. |
| Composition root | Module workers publish constructors; M004-11 exclusively wires shared composition-root registration after merges. Any unavoidable earlier registration is a named, additive-only, time-boxed lock. |
| Integration DB | Divergent migration fingerprints run serially; matching fingerprints may run concurrently. (`HARNESS-CORE.md` §3) |

## MIS-005 milestone headlines

| Milestone | Headline |
|---|---|
| **M005-01 Production identity and tenancy** | Replace the fixed demo tenant/no-auth posture with authenticated tenant resolution, tenant administration, authorization, and isolation evidence. (`repo-baseline-2026-07-17.md` Tenancy/auth) |
| **M005-02 Production mutation execution** | Enable governed ML writes through M-03 only, with resolved linkage, policy/source time, idempotency, audit, reconciliation, and operator-controlled rollout. |
| **M005-03 Webhook convergence** | Add `orders_v2`, shipments, claims, items, payments, item competition, and Full-stock events, with GET reconciliation and duplicate/out-of-order protection. (`API-MAP.md` General rules) |
| **M005-04 Scheduled market history** | Start daily snapshots, retained evidence, 7/90-day history, alerts, and collection observability; no retroactive history claims. |
| **M005-05 Mercado radar** | Deliver Reprecificação, Oportunidades, Monitorados, official benchmarks, item-to-win signals, and alert workflows without depending on forbidden public-search behavior. |
| **M005-06 External provider qualification** | Run contract review → frozen canary-five → batch-50. Integrate a provider only after OEM/storage/display rights and matching gates pass; otherwise close as no-provider. (`pricing-intelligence-implementation-handoff.md` §§9–10) |
| **M005-07 Complete product workspace** | Add Concorrência, Pedidos, Histórico, auditoria, Dados, visits, coverage history, movements, and Full inventory/operations. |
| **M005-08 Complete orders lifecycle** | Add cancellations, claims, returns, reputation incentives, reverse logistics, dispute actions, refund state, NF-e/faturamento, and bulk labels. |
| **M005-09 Fiscal and configuration center** | Complete 27-UF DIFAL lifecycle, exceptions, scheduling, reminders, paid state, notifications, and the approved global → product → listing configuration hierarchy. |
| **M005-10 Repasses and reconciliation** | Add scheduled asynchronous Mercado Pago release-report ingestion, billing reconciliation, release calendar, retention detail, and ERP reconciliation. (`API-MAP.md` Repasses) |
| **M005-11 Inventory and fulfillment depth** | Add Full inventory, fulfillment operations/events, coverage alerts, stock movements, and production-grade ERP convergence. |
| **M005-12 Analytics and operational hardening** | Finish dashboard analytics, production deployment/runtime, observability, recovery, retention, security hardening, and full live-integration validation. |

## Top risks

| Risk | Likelihood | Impact | Mitigation |
|---|---:|---:|---|
| Three-day schedule leaves almost no recovery time | High | Critical | Freeze contracts first; dispatch disjoint module/route chips; keep Dashboard strictly cuttable; rehearse the central `.xlsx` → verdict → simulator → Pedidos story before optional polish. |
| `buy_box_winner` and its range remain null, producing `NO_PRICE_EVIDENCE` during the client demo | High | High | Preflight the actual workbook on demand; present the honest state deliberately; prioritize active-listing `sale_price`/`price_to_win` evidence; never substitute zero or an unsupported “market price.” The frozen probe observed null in 22/22 catalog products. (`pricing-intelligence-implementation-handoff.md` §2.3) |
| W1 chips are unmerged or merge with surfaces different from baseline assumptions | High | Critical | Make post-W1 accepted main the mandatory branch point; rerun collision/contract checks after merge; do not reconstruct M-02/M-03 from planning notes. (`repo-baseline-2026-07-17.md` Git / estado W1) |
| Client workbook violates the proposed column/data contract | Medium | High | Obtain a pre-demo dry run; provide row-level validation and rejection evidence; require the four mandatory columns; never coerce missing cost/stock to zero. |
| Existing ML app permissions or live account data do not support a required read | Medium | High | Run read-only capability preflight early against the real installation; expose unavailable/stale states; keep forbidden endpoints and fallback scraping out. |
| Shared OpenAPI/SDK or ML adapter changes serialize several tracks | Medium | High | Exclusive M004-01 contract lock and M004-04 adapter ownership; use published ports; reject opportunistic cross-chip edits. |
| Identity collision creates a convincing but wrong product/price match | Medium | Critical | Deterministic two-anchor gate, hard negatives, manual review, collision fixtures, and no title-only auto-accept. |
| DIFAL seed is mistaken for tax advice or authoritative live fiscal data | Medium | High | Label source/version and demo scope; make destination and toggle explicit; surface unknown values; defer payment lifecycle and ongoing fiscal maintenance to MIS-005. |
| Demo controls accidentally execute a provider mutation | Low | Critical | Provider-write execution absent/disabled in demo config; all mutation UI goes through M-03 preview/protocol; validate network/audit logs for zero provider writes. |
| Parallel migrations break shared integration runs | Medium | Medium | Preallocated number blocks, exclusive table owners, runner-count update, and serialized tests for divergent migration fingerprints. |
| Dashboard consumes time while the core story remains unstable | Medium | Medium | Do not start M004-10 until all mandatory journeys are integrated and green; cut it without reopening scope. |

## Deliberate exclusions

- **No live Mercado Livre writes:** no remote price, stock, pause, create-listing, faturamento, or other mutation. The demo shows governed preview/protocol only.
- **No webhooks or daily scheduler:** polling, explicit refresh, and on-demand preparation are sufficient for the local read-only demonstration.
- **No scraping, browser automation, public-search fallback, proxy/stealth infrastructure, or the closed NO-GO providers.** They are unsupported legally or empirically. (`pricing-intelligence-implementation-handoff.md` §3)
- **No unhomologated external price provider:** DataForSEO/Precifica or any successor remains MIS-005 work behind contract, canary, and batch gates.
- **No unconditional `/products/{id}/items`:** it remains default-off and can be enabled only with full pagination, telemetry, and explicit fallback.
- **No promise of automatic pre-listing market price:** valid results include review and absence-of-evidence states.
- **No complete Mercado radar, monitored history, repricing, opportunities, or alerts:** those require scheduler/history and belong to MIS-005.
- **No complete DIFAL workflow or Configurações application:** exclude scheduling, reminders, paid state, exception editing UI, and general integrations/notification settings. MIS-004 keeps only the seed, simulator toggle, and orders indication.
- **No claims, returns, cancellation workflow, disputes, reputation actions, reverse logistics, or production NF-e:** Pedidos covers current operational reads, SLA, timeline, tracking, and read-only profitability.
- **No full Produto Detalhe:** Concorrência, Pedidos, Histórico, auditoria, Dados, visits, and Full-depth views remain MIS-005.
- **No Repasses/reconciliation:** asynchronous CSV release-report ingestion is not necessary for the demo story.
- **No real auth/multi-tenancy initiative:** preserve explicit tenant predicates and secret handling, but defer the greenfield authentication/tenant-resolution surface to MIS-005.
- **No cloud deployment or production SLO work:** MIS-004 is validated on the fixed local Docker topology.
- **No copied `.dc.html` production code or placeholder business data:** the design files are high-fidelity references; counts, multichannel Dashboard content, hardcoded DIFAL-SP behavior, and other mock defects are excluded. (`README.md` About the Design Files; `design-screens-2026-07-17.md` Dashboard, Simulador, Shell/nav)