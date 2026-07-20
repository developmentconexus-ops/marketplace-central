# Refactor Inventory — Frontend (MIS-006-integracao-fundacao)

Base: main @138aac3d. Target docs: SYSTEM-BLUEPRINT.md, SCENARIO-WALKTHROUGH.md, INTEGRATION-DATA-CONTRACT.md.
Rule of scope: verified real file:line by opening files (not grep-guessed).

## 1. Oportunidades (`apps/web/src/pages/mercado/oportunidades.ts`)

| path:line | current state | target (blueprint §2,§6,§7) | action | why |
|---|---|---|---|---|
| `oportunidades.ts:49-101` `buildOppRows()` | 100% client-side join of 3 fetched endpoint payloads (facts, aggregates, verdicts), no backend query, no persistence | `/mercado` Oportunidades = `products_mirror − product_links JOIN market_aggregates`, one SQL view read by API (blueprint §6 diagram + §7 table + literal SQL in SCENARIO-WALKTHROUGH.md:117-125) | REFACTOR (become a thin renderer over a new backend endpoint) | violates "tela NUNCA chama API externa / lê Postgres" golden rule; today does client-side cross-referencing of live endpoint calls instead |
| `oportunidades.ts:68` `if (!agg \|\| agg.status !== "OK" ...) continue` | only filter is aggregate status; never reads `product_links`/`listings` | target: `LEFT JOIN product_links pl ... WHERE pl.link_id IS NULL` — must exclude products we already sell | REFACTOR (drop client filter once backend view does the join) | F3 in SCENARIO-WALKTHROUGH.md flaw table — "Exclusão 'não vendemos' ausente" |
| `oportunidades.ts:70,88-93` cost/gap sort using `fact.cost.amount` | sorts by raw ML-median minus ERP-cost gap; rows w/ null cost silently sink via `-Infinity` | target A6: costo NULL must still surface, ranked by demand, labeled "margem indisponível" — two sort modes (margem / demanda) | REFACTOR | current single-mode sort hides NULL-cost rows at the bottom instead of offering demand-only ranking |
| `MercadoPage.tsx:21-23` / `MonitoradosTab.tsx:28-84` | "Monitorados" tab renders design chrome only; all controls `disabled`; explicit comment "not wired yet" | target: scope "monitorado" (coluna + regra elegibilidade + tela) is Fase-1 real feature (blueprint §7 row `/mercado Oportunidades`) | CREATE (new, not a repair) | grep `monitor` in backend = 0 (per audit); this is intentionally a stub with honest empty state, not a bug — needs real backend before FE can stop being a stub |

## 2. /importacoes (no dedicated route — folded into /vinculos + /integracoes)

| path:line | current state | target (blueprint §2 step 6, §7) | action | why |
|---|---|---|---|---|
| `apps/web/src/app/AppRouter.tsx:42-61` | no `/importacoes` route exists at all; closest legacy alias is `/integrations → /integracoes` (line 59) | blueprint §7: `/importacoes` is its own screen showing "cadeia protocolo→vínculos→coletas" | CREATE (or explicitly fold target into /integracoes — needs a scope decision) | route doesn't exist; current UI splits the same concern across two pages instead |
| `apps/web/src/pages/vinculos/ImportacaoSection.tsx:115-146` `ImportacaoSection()` | read-only history list of `erp_import_*` protocols: shows `accepted_count/rejected_count/warning_count/imported_at` + expandable issue rows (`ImportRowDetail`, `useErpImportDetail`) | target cadeia: "N importados → N vinculados → N coletas" — must show downstream link + collection counts per protocol, not just parse-stage counts | REFACTOR | component only reads `erp_import_products`/protocol table; never joins `product_links` or `sync_state`/collection queue, so the "vinculado"/"coletado" legs of the chain are entirely absent |
| `apps/web/src/pages/vinculos/ImportacaoSection.tsx:6-8` used from `VinculosPage.tsx:8,146` AND `IntegracoesPage.tsx:6,397` | same component instance rendered on two different pages (/vinculos and /integracoes) | blueprint implies /importacoes as one coherent screen | REFACTOR / collision note | shared-component seam: any change to `ImportacaoSection` affects both `/vinculos` and `/integracoes` pages simultaneously — flag before touching |

## 3. /integracoes (`apps/web/src/pages/integracoes/IntegracoesPage.tsx`)

| path:line | current state | target (blueprint §3 onboarding saga, §7) | action | why |
|---|---|---|---|---|
| `IntegracoesPage.tsx:385-400` `IntegracoesPage()` | assembles 4 cards: `ActiveSourceCard`, `UploadCard`, `ProviderConnectCard`, `ImportacaoSection` — this is the real, functioning screen (route wired, live ML OAuth start per d782a89f) | target: same screen, plus onboarding saga steps 2-6 (backfill anúncios/pedidos progress, candidatos vínculo, coleta mercado, "pronto") per blueprint §3 mermaid | REFACTOR (extend, base already exists and is not a stub) | KEEP the 4 existing cards; screen currently stops after "Conectar" (OAuth start) — nothing shows post-connect progress |
| `IntegracoesPage.tsx:280-319` `ActiveSourceCard()` + note at 270-274 | `useActiveErpSource()` (from `@marketplace-central/web-query`) — comment explicitly says "The choice persists in **localStorage**" | target §2: "Config fonte ativa POR TENANT, no banco (`MC_ERP_SOURCE` morre)" | REFACTOR | localStorage-per-browser violates "por tenant, no banco" — same operator on a different machine loses/duplicates the toggle; this is the literal target callout ("MC_ERP_SOURCE morre") |
| `IntegracoesPage.tsx:328-383` `ProviderConnectCard()` / `connect()` | live: lists installations, reuses a `pending_connection` one or creates new, calls `startIntegrationAuthorization`, redirects to `auth_url` (real OAuth start, not a stub) | target §3 step 1 "Conectar OAuth ML" — matches; steps 2-6 (backfill/candidatos/coleta/pronto) not present | KEEP as-is for step 1; CREATE steps 2-6 | button is genuinely wired (matches memory note d782a89f); the gap is everything AFTER connect |
| `IntegracoesPage.tsx:1,7-8` imports `ErpImportSourceInput` from sdk-runtime, reuses `ImportacaoSection`/`useErpImportDetail` from `../vinculos/*` | cross-page import from `pages/vinculos` into `pages/integracoes` | n/a — architecture note | REFACTOR (consider promoting to a shared module if /importacoes becomes its own screen) | current relative import `../vinculos/ImportacaoSection` couples the two page folders; fine short-term, but blocks cleanly splitting /importacoes out later |

## 4. Vínculo/linkage UI (`apps/web/src/pages/vinculos/*`)

| path:line | current state | target (blueprint §2 step 4, SCENARIO-WALKTHROUGH FASE 2, F5) | action | why |
|---|---|---|---|---|
| `VinculosPage.tsx:50-146` `VinculosPage()` | real screen: KPI row (Pendentes/Alta confiança/Resolvidos hoje), tabs Fila/Resolvidos, calls `listProductLinkCandidates` / `listProductLinkWorkflows` via SDK (`vinculosQueryKeys.ts`) | matches target's "candidatos vínculo" concept; target adds AUTO-APPROVED exact-EAN links flowing straight to Resolvidos without operator action | KEEP screen shell; REFACTOR data assumption | screen already assumes a manual approval queue; once backend auto-approves exact-EAN (F5, not yet wired — `resolution_service.go:129` per audit, backend-side), the Fila tab's population will shrink and Resolvidos needs an "auto" provenance badge |
| `QueueRow.tsx:79-81,155` `gtinEqual = candidate.match_input === "ean" && ...`, `pill(bandLabels[candidate.confidence_band], ...)` | already renders EAN-match confidence band (`ALTA`/etc.) and an explicit "GTIN unknown → —" honesty rule — real, non-stub display logic | matches target's confidence-band display need | KEEP | correctly implements ADR-17 honesty pattern already; no fabricated data found here |
| `VinculosPage.tsx:20-35` `isResolved`/`isResolvedToday` client-side filters over `listProductLinkWorkflows` | KPI counts computed client-side from a full workflow list fetch, not from a backend aggregate | minor — blueprint doesn't mandate server-side KPI aggregation explicitly, but golden rule favors DB-computed counts | REFACTOR (low priority) | fine at current scale (34-2000 rows) but same class of client-side aggregation the audit flagged for Oportunidades |

## 5. SDK client surface (`packages/sdk-runtime/src/index.ts`)

| path:line | current state | target | action | why |
|---|---|---|---|---|
| `index.ts:1849` `createErpImport(file, source, fileName?)` | POST upload, source-typed (`catalogo_cliente` \| `xlsx`) — real, used by `useErpImportUpload.ts` | matches blueprint §2 step 1 (upload) | KEEP | functioning |
| `index.ts:1855` `listErpImports()` → `GET /erp/imports` | protocol history list only; no linked/collected counts in payload (confirmed via `ErpImportSummary` fields used in `ImportacaoSection.tsx:92-109`: accepted/rejected/warning/imported_at only) | target needs a richer shape carrying vínculos+coletas counts per protocol for the /importacoes chain | REFACTOR (extend response type + client method) | client-side type genuinely lacks the fields the target screen needs; not just a display gap |
| `index.ts:1864` `startIntegrationAuthorization(installationId)` | real OAuth-start call, used live in `ProviderConnectCard` | matches target step 1 | KEEP | — |
| `index.ts:1904,1908` `listProductLinkCandidates` / `listProductLinkWorkflows` | real, paginated (`limit` param, default 20) | matches target "candidatos vínculo" | KEEP | — |
| No client method found for: onboarding/backfill progress (`sync_state`), source-active-per-tenant (DB), monitored watchlist | absent | target needs new endpoints: onboarding progress, DB-backed active-source config, monitored-scope CRUD | CREATE | no grep hits for any of these method names in sdk-runtime — confirms these are net-new, not misnamed existing calls |

## Cross-cutting collision surfaces (routes/nav/shared components — touch with care)

- `apps/web/src/app/AppRouter.tsx:42-61` — owns all top-level routes; adding `/importacoes` (if scope decides to split it out) is a route-table edit here plus a nav-link edit in `Header.tsx`.
- `apps/web/src/pages/vinculos/ImportacaoSection.tsx` — shared between `/vinculos` (VinculosPage.tsx:146) and `/integracoes` (IntegracoesPage.tsx:397). Any refactor of this component changes both pages.
- `@marketplace-central/web-query`'s `useActiveErpSource` — currently localStorage-backed; migrating to DB-backed per-tenant config touches this shared package, not just `IntegracoesPage.tsx`.
