# P5 Sol Decomposition Audit — r03 (VERBATIM)

> Provenance: GPT-5.6 Sol / model_reasoning_effort=medium / --sandbox read-only, codex exec OS-process.
> Prompt: `C:\Users\LEANDR~1.THE\AppData\Local\Temp\claude\C--Users-leandro-theodoro-Documents-marketplace-central\3b870e2f-2c67-4f47-9e53-6939f20380e4\scratchpad\prompt-p5-sol-r03.md` · Log: `C:\Users\LEANDR~1.THE\AppData\Local\Temp\claude\C--Users-leandro-theodoro-Documents-marketplace-central\3b870e2f-2c67-4f47-9e53-6939f20380e4\scratchpad\agent__p5-sol-r03.log`
> Input manifest: `planning-reviews/p5-input-r03.sha256` — top digest `8448f39eca1aeda724f0242bc33009275756a8dda989a160e86f68089a748142`
> Verdict: NEEDS-FOLD. Stdout below verbatim.

---
## Manifest verification

Digest match: **yes**.

- 47/47 manifested files exist and match their per-file SHA-256.
- Recomputed canonical top digest: `8448f39eca1aeda724f0242bc33009275756a8dda989a160e86f68089a748142`.

## Fold verification

| r02 finding ID | Status | Evidence | Note |
| --- | --- | --- | --- |
| P5-R2-01 | **NOT-CLOSED** | `M-01-erp-xlsx-identity/F-01-identity-semantics-fix/feature.md`, §Ownership: “tipos catalog em `packages/sdk-runtime/src/index.ts` (**additive-lock grant**)”; but `M-01-erp-xlsx-identity/milestone.md`, §Ownership: “OpenAPI seção `/erp/imports*`, `packages/sdk-runtime/src/erpImport.ts`” and “export 1-linha em `sdk-runtime/src/index.ts`”. | The feature and mission matrix grant catalog type edits, but the milestone executable write-set still describes `index.ts` only as a one-line barrel export and omits the catalog OpenAPI/SDK grant. |
| P5-R2-02 | **NOT-CLOSED** | `M-01-erp-xlsx-identity/F-01-identity-semantics-fix/feature.md`, §Negative Scenarios: “REFERENCIA vazia/whitespace ⇒ `ean: null`, `refforn: null`”; IC-01: “`refforn` … TGFPRO.REFFORN.” | Blank REFERENCIA cannot independently null a valid REFFORN. The negative scenario still contradicts the Brief, EARS, and IC-01. |
| P5-R2-03 | **NOT-CLOSED** | IC-03 §Operations: “`GET /market/aggregates?codprod=`”; `M-05-anuncios-sinais/F-01-listings-signals-api/feature.md`, §Inputs: “agregados/veredictos por `codprods`”. | The producing HTTP brief is corrected, but a consuming interface description retains the rejected identifier token. Reconciliation’s claimed zero-occurrence sweep is false. |
| P5-R2-04 | **CLOSED** | IC-03 §Operations: “**200 síncrono** … sem tabela de job, sem polling”; M-02 F-04: “**200** com sumário”; M-06 F-01: “POST síncrono retorna 200 … sem polling”. | Contract, producer, and consumer now agree. |
| P5-R2-05 | **NOT-CLOSED** | `mission.md`, ADR-01: “M-03 envelope = único write path”; Runtime Contract and ADR-08 instead qualify it as the only provider-target write path and allow local M-04 writes. | The old universal rule survives unqualified in a binding cross-cutting decision. |
| P5-R2-06 | **CLOSED** | M-07 F-02 §Expected Output: “`source`, `fetched_at`/freshness, `n_offers`/`n_sellers`, `match_status`”. | Simulador now propagates the complete IC-03 evidence set. |
| P5-R2-07 | **CLOSED** | IC-04 §Decomposition debits “`tarifa_full`”; M-07 F-01 output includes “`difal, tarifa_full, custo`” and propagates null in `full`. | Contract formula and executable decomposition agree. |
| P5-R2-08 | **CLOSED** | IC-04 §Operations and M-07 F-01 §Expected Output both specify “`PUT /pricing/difal/{uf}`”. | Exact route propagated. |
| P5-R2-09 | **CLOSED** | M-09 §Dependencies: “M-01, M-03, M-04, M-05, M-08. M-07 NÃO produz artefato”; identical producer set in `mission.md` DAG. | False M-07 dependency removed. |
| P5-R2-10 | **CLOSED** | R-04 §Caveats item 4: “Product-specific taxation is OUTSIDE … no per-SKU fiscal override surface exists in MIS-004.” | R-04 now matches IC-04 scope. |

## New findings

### P5-R3-01

- Audit check violated: regression sweep — IC-05 gear-menu propagation and no new product scope.
- Cited excerpt: IC-05 fixes “⚙(menu: Configurações → seção DIFAL read/drawer; Integrações; Catálogo; Estoque).” M-03 F-02 calls its menu “EXATO IC-05” but adds “`Vínculos` (entrada secundária).”
- Exact defect locus: `M-03-shell-retheme/F-02-header-nav-routes/feature.md`, §Expected Output, offending token `Vínculos` inside the gear menu.
- Severity: **BLOCKING**.
- Yes-if: the executable menu contains only the IC-05 entries; Vínculos remains outside global navigation and reachable through the already-approved Anúncios/importação path and registered route.

Standing spot-checks found no additional defect: DAG forcing edges and the M-07→M-08 intra-wave gate are present; M-09 dependencies match the mission DAG; active migration blocks remain disjoint at 0045–0069 with reserve 0070–0074; the fixture seam is `apps/server_core/internal/platform/migrate/runner_test.go`, currently asserting 41 migrations.

## Verdict

NEEDS-FOLD

## Paths checked

Manifest and planning reviews:

- `.mnfs/MIS-004-mvp-demo/planning-reviews/p5-input-r03.sha256`
- `.mnfs/MIS-004-mvp-demo/planning-reviews/p5-sol-decomposition-audit-r02.md`
- `.mnfs/MIS-004-mvp-demo/planning-reviews/p5-reconciliation-r02.md`
- `.mnfs/MIS-004-mvp-demo/planning-reviews/p5-decomposition-passes-r01.md`
- `.mnfs/MIS-004-mvp-demo/planning-reviews/p5-sol-decomposition-audit-r01.md`
- `.mnfs/MIS-004-mvp-demo/planning-reviews/p5-reconciliation-r01.md`
- `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-reconciliation-r01.md`

Mission and research:

- `.mnfs/MIS-004-mvp-demo/mission.md`
- `.mnfs/MIS-004-mvp-demo/research/identity-matching-interface-contract.md`
- `.mnfs/MIS-004-mvp-demo/research/erp-xlsx-import-interface-contract.md`
- `.mnfs/MIS-004-mvp-demo/research/market-evidence-read-interface-contract.md`
- `.mnfs/MIS-004-mvp-demo/research/pricing-difal-interface-contract.md`
- `.mnfs/MIS-004-mvp-demo/research/fe-shell-seams-interface-contract.md`
- `.mnfs/MIS-004-mvp-demo/research/ml-read-ports-interface-contract.md`
- `.mnfs/MIS-004-mvp-demo/research/difal-interna-rates-2026.md`
- `.mnfs/MIS-004-mvp-demo/research/design-screens-2026-07-17.md`
- `.mnfs/MIS-004-mvp-demo/research/p1-clarified-decisions-2026-07-17.md`
- `.mnfs/MIS-004-mvp-demo/research/repo-baseline-2026-07-17.md`
- `.mnfs/MIS-004-mvp-demo/research/w1-merge-addendum-2026-07-17.md`

Milestones and features:

- `.mnfs/MIS-004-mvp-demo/M-01-erp-xlsx-identity/milestone.md`
- `.mnfs/MIS-004-mvp-demo/M-01-erp-xlsx-identity/F-01-identity-semantics-fix/feature.md`
- `.mnfs/MIS-004-mvp-demo/M-01-erp-xlsx-identity/F-02-erp-import-module/feature.md`
- `.mnfs/MIS-004-mvp-demo/M-01-erp-xlsx-identity/F-03-reader-adapter-selection/feature.md`
- `.mnfs/MIS-004-mvp-demo/M-02-price-intel-core/milestone.md`
- `.mnfs/MIS-004-mvp-demo/M-02-price-intel-core/F-01-ml-adapter-read-ports/feature.md`
- `.mnfs/MIS-004-mvp-demo/M-02-price-intel-core/F-02-market-persistence/feature.md`
- `.mnfs/MIS-004-mvp-demo/M-02-price-intel-core/F-03-identity-resolver/feature.md`
- `.mnfs/MIS-004-mvp-demo/M-02-price-intel-core/F-04-collect-verdict-api/feature.md`
- `.mnfs/MIS-004-mvp-demo/M-03-shell-retheme/milestone.md`
- `.mnfs/MIS-004-mvp-demo/M-03-shell-retheme/F-01-theme-tokens-fonts/feature.md`
- `.mnfs/MIS-004-mvp-demo/M-03-shell-retheme/F-02-header-nav-routes/feature.md`
- `.mnfs/MIS-004-mvp-demo/M-03-shell-retheme/F-03-shared-primitives/feature.md`
- `.mnfs/MIS-004-mvp-demo/M-04-vinculos-import-ui/milestone.md`
- `.mnfs/MIS-004-mvp-demo/M-04-vinculos-import-ui/F-01-product-links-api-gaps/feature.md`
- `.mnfs/MIS-004-mvp-demo/M-04-vinculos-import-ui/F-02-vinculos-screen/feature.md`
- `.mnfs/MIS-004-mvp-demo/M-05-anuncios-sinais/milestone.md`
- `.mnfs/MIS-004-mvp-demo/M-05-anuncios-sinais/F-01-listings-signals-api/feature.md`
- `.mnfs/MIS-004-mvp-demo/M-05-anuncios-sinais/F-02-anuncios-ui-sinais/feature.md`
- `.mnfs/MIS-004-mvp-demo/M-07-simulador/milestone.md`
- `.mnfs/MIS-004-mvp-demo/M-07-simulador/F-01-pricing-calc-difal-service/feature.md`
- `.mnfs/MIS-004-mvp-demo/M-07-simulador/F-02-simulador-ui/feature.md`
- `.mnfs/MIS-004-mvp-demo/M-08-pedidos/milestone.md`
- `.mnfs/MIS-004-mvp-demo/M-08-pedidos/F-01-orders-projection-api/feature.md`
- `.mnfs/MIS-004-mvp-demo/M-08-pedidos/F-02-pedidos-ui/feature.md`
- `.mnfs/MIS-004-mvp-demo/M-06-produto-detalhe/milestone.md`
- `.mnfs/MIS-004-mvp-demo/M-06-produto-detalhe/F-01-produto-detalhe-page/feature.md`
- `.mnfs/MIS-004-mvp-demo/M-09-dashboard-demo/milestone.md`
- `.mnfs/MIS-004-mvp-demo/M-09-dashboard-demo/F-01-dashboard-mpc/feature.md`

Permitted live-code verification:

- `apps/server_core/internal/platform/migrate/runner_test.go`
- `apps/server_core/migrations/` — filename inventory only.