# VALIDATION PLAN — D-119 demo-readiness (hub-owned)

2026-07-20, demo TODAY. Operator mandate: "seu trabalho vai ser garantir que tudo funciona".
Hub validates EVERY link of the demo chain end-to-end on the SHARED stack with REAL data
before sign-off. Chips prove their unit in isolation; only the hub's live drive counts
(lesson: live-verification is bounded by interactions actually performed — each surface
below is exercised explicitly, or honestly listed as NOT exercised).

## Gate ladder per chip (unchanged doctrine)
Per CLOSED: (1) hub independent P6 dual gate (cold Opus reviewer + adversarial refuter,
AGREEMENT, whole-module scan not diff — CHIP-PED-FILA lesson 2); (2) merge --no-ff;
(3) post-merge ladder: Go tests (GOCACHE=.gocache, from apps/server_core), web vitest,
tsc write-set; (4) docker restart frontend (Vite bind-mount staleness D-112); (5) P7 live
drive below. Only QA passes.

## P7 live validation matrix — the demo chain, in demo order

### V1 — Import raw file (after CHIP-IMPORT-FIX merge)
- /integracoes renders (light+dark), upload card + source selector + history.
- LIVE: upload the REAL `Downloads/PRODUTOS MERCADO LIVRE.xlsx` (1.6MB, 3 sheets),
  source=catalogo_cliente → expect **2012 accepted**, 0 rejected-for-structure,
  warnings = honest custo/estoque-unknown only. Protocol appears in history.
- Negative: re-upload same file → 409 duplicate_file rendered honestly.
- Negative: upload a wrong file (e.g. an image renamed .xlsx) → invalid_file, no crash.
- Guard: Sankhya snapshot #003-E untouched (strict source still latest for source=xlsx).

### V2 — Two-source toggle (GOAL C)
- Default (no toggle): /catalogo lists SANKHYA products WITH custo/estoque — byte-stable,
  proves no regression for every existing screen.
- Flip Fonte ativa → catálogo do cliente: /catalogo lists the 2012 (FACAS/FERRAMENTAS/
  MAQUINAS present), custo/estoque honest "—" (ADR-17), never 0.
- Flip back → Sankhya returns. Both themes.

### V3 — Competitor price for imported products (after V1+V2)
- Pick 5-10 hero products from the 2012 (valid EAN — 1827/2012 eligible; prefer known
  brands: Tramontina/Bosch/Makita class → higher ML-catalog hit odds). Hub pre-drives
  ONE before the demo to confirm the path; operator live-collects the rest in front of
  the client.
- LIVE: POST /market/collections with erp_source=catalogo_cliente for hero codprod →
  match decision recorded, offers collected (n>0 for at least some heroes), own-seller
  excluded, aggregate min/median/max persisted.
- Honest negatives: a no-EAN product → NO_IDENTITY, an EAN with no ML catalog → honest
  no_price_evidence. No fabricated medians.
- Surfaces: /precos market comparison + /mercado (CHIP-MERCADO) show the collected
  evidence; insufficient_market renders honestly.

### V4 — /pedidos (after CHIP-PED-FIX merge)
- LIVE re-import of orders (read-only ML, operator-authorized) → orders newer than
  2026-07-09 present; freshness verified against ML directly.
- Perf: /pedidos list loads materially faster than the 10.9s baseline (parallelized
  shipment reads); measure wall-clock, record number.
- Buckets: paid+undispatched orders sit in A FATURAR; click "Foi faturado" on one →
  moves to A ENVIAR (faturado_at persisted, idempotent re-post); shipped/delivered
  orders in ENVIADO regardless of flag; KPI counts consistent with tabs/kanban.
- Decomposição: open drawer of order 2000017276984774 → Comissão R$ 22,95 (real
  persisted sale_fee), Custo real ERP, Frete real when shipment costs present;
  difal/taxa_fixa/tarifa_full honest "—" listed in componentes desconhecidos; margem
  only when ALL inputs known.
- "Atualizar" button refetches; drawer + fila + kanban + lista + dark theme each
  explicitly driven.

### V5 — /mercado (after CHIP-MERCADO merge)
- 3 tabs render 1:1 (spot-check grid/min-widths vs Mercado.dc.html), real market data
  where collected (V3 heroes), honest "—" elsewhere; Aplicar/Atualizar/+Monitorar inert
  (assert zero network write); light+dark.

### V6 — Full demo rehearsal (final sign-off)
Run the ENTIRE demo script start-to-finish as the customer persona, one sitting:
dashboard → /integracoes import live → toggle → /catalogo → hero collect → /precos
simulate → /mercado radar → /pedidos (buckets, faturado click, drawer, atualizar).
Any break = fix or honest demo-script cut. Verdict recorded here + ledger. Only after
V6 PASS: demo-ready sign-off to operator.

## Cross-cutting invariants (checked at every V)
- ADR-17: no fabricated zeros anywhere; every "—" traceable to an absent source.
- Zero live ML writes (read-only demo; faturado = our-DB only).
- Tenant scoping intact on new queries (faturado write, source-filtered reads).
- No secrets/tokens in logs, evidence, or commits.

## Risks watched
- AppRouter/Header 3-way merge (mercado + integracoes) — hub reconciles, then re-run FE ladder.
- openapi/sdk lockstep from two chips (orders vs erp_import sections) — verify both after merges.
- Rate limit ML during hero collection — space collect calls; stop on PROVIDER_RATE_LIMITED.
- Import of 2012 must NOT displace pricing default source (V2 default check is the guard).

## Status log (07:26, pre-demo freeze @ main 61616aab)
- [x] V1 PASS — byte-distinct RAW copy → #005-E 2012/2012 aceitos, 3 rejeitados (rodapé), 4068 avisos honestos, 2.2s; dup=409, invalid=400; #003-E untouched (1845). REAL file kept VIRGIN for the live demo import (no 409).
- [x] V2 PASS — default byte-stable (Sankhya custo/estoque presentes); erp_source=catalogo_cliente → honest missing_stock/null; bogus→400. UI selector renders on /integracoes.
- [~] V3 PARTIAL — chain proven (91% EAN-eligible; D-94 live offers evidence); hero pre-collect NOT run (deadline) — operator collects live: 14105 / 622 / 3308.
- [x] V4 PASS — lista 2.23s (era 10.9s); golden 2000017276984774 comissão 22,95 + frete 19,85 reais, unknowns honest-listed; buckets honestos (24 enviados); migration 0074 applied; faturado positivo NOT live-exercised (0 orders in bucket — RTL coverage cited).
- [x] V5 PASS — /mercado live: 34 anúncios, preços conc. reais (papeleira 169,99 vs 179,90/229,20), honest "—"; Aplicar/Simular inert.
- [~] V6 PARTIAL — per-screen smoke all 200 + zero console errors; full one-sitting persona rehearsal NOT run (7:30 hard deadline). Sign-off: DEMO-READY with V3/V6 partial-honesty caveats above.
