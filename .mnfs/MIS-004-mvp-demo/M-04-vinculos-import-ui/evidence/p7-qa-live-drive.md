# P7 QA live-drive evidence — M-04 vinculos-import-ui (SCOPED, hub ruling option 2)

Stack: hub-booted, backend :8080 + frontend :5174, shared dev PG, `MPC_PROVIDER_WRITES_ENABLED` unset (ML dispatcher OFF).
installation_id = `inst-mercado_livre-d373dc64-577f-4950-b11b-b90244f30cb2` (Mercado Livre / METALNOBREACABAMENTOS).
Data context (hub-ratified): synthetic-only ERP catalog (products 2002-2004), zero EAN/seller_sku overlap with connected ML listings ⇒ every candidate legitimately NO_CANDIDATE (ADR-17 honest output, NOT a defect). Candidates regenerated (POST link-candidates/generations) to populate enrichment fields — sanctioned by hub.

## PART A — UI render live-drive (browser pane at :5174/vinculos) — CAPTURED ✅

Source: `mcp__Claude_Browser__read_page` accessibility tree (DOM-grounded).

- **Shell + nav render**: banner "Marketplace Central", primary nav (Visão geral / Anúncios / Simulador / Pedidos + "em breve" pills), theme toggle, installation combobox (Mercado Livre selected), settings menu (DIFAL/Integrações/Catálogo/Estoque).
- **/vinculos page renders**: heading "Vínculos" + subtitle; **KPI Resumo** region with 3 tiles (Pendentes / Alta confiança / Resolvidos hoje); **tablist "Filtros de vínculos"** with tabs **Fila** + **Resolvidos** (role=tab); queue **region "Fila de candidatos"** with table columns Selecionar / Produto / Melhor candidato / Sinais / Confiança / Ações.
- **NO_CANDIDATE honesty live-proven (ADR-17)** — every row (34 candidates rendered) shows:
  - "Sem candidato" pill (honest, not a fabricated match)
  - honest message `Nenhuma correspondência encontrada para MLB<id>.`
  - anchor-signal list, all UNAVAILABLE with honest detail: `seller_sku: seller_sku sem correspondência` / `ean: ean sem correspondência` / `marca: marca inexistente no lado provider` / `refforn: refforn inexistente no lado provider`
  - Confiança column = **"sem confiança sem candidato"** — honest UnknownValue, **NOT** a fake "0%" / "NaN" / "0.00%". (Confirms ADR-17 no-silent-zero on the confidence surface; the double-scale "9500%" P6 defect is structurally absent — chip renders honest text, not a scaled numeral. Numeric-band happy-path chip covered by unit fixtures + P6 gate per hub ruling.)
  - per-row actions **Abrir** (drawer) + **Aprovar** present.
- Drawer open (clicked row "Abrir") registered; DOM re-read confirmed queue still mounted. Full drawer-body capture + manual-resolve UI drive INTERRUPTED by stack-down (see below).

## PART B — API contract sweep — NOT RUN (blocked: stack went down)

Prepared sweep (`scratchpad/p7-sweep.ps1`) covers: manual-resolve→leaves-queue→undo audit chain (to synthetic product 2002), ALREADY_RESOLVED semantics, ALREADY_UNDONE/SUPERSEDED 409, NOT_FOUND 404, PRODUCT_NOT_FOUND 404, INSUFFICIENT_ANCHORS 422, batch-preview dry-run (zero-persist re-GET proof), batch-apply failed[] itemization, empty-batch 422, C04 import history + #002-E rejection report. **Could not execute** — backend :8080 became unreachable before the sweep ran (HTTP:000 / curl exit 7).

## INTERRUPTION — hub stack went down mid-QA (2026-07-18)

`docker ps` mid-sweep: NO `marketplace-central` containers running. Ports :8080 + :5174 both HTTP:000. A different project's stack (`metaldocs-*`, api :8081, postgres :5433) now occupies docker — machine/docker contention reclaimed the marketplace-central stack. Chip cannot restart (hub seam / chip never boots servers/compose). **REQUEST sent to hub to restore the marketplace-central stack; PART B sweep resumes once up.**

VERDICT: **HELD** — PART A (UI render + NO_CANDIDATE honesty) PASS; PART B (contract/mechanics) pending stack restoration. No milestone verdict until PART B completes on the restored stack + independent review.
