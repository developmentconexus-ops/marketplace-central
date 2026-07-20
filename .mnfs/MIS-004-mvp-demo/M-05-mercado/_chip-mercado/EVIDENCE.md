# CHIP-MERCADO — Evidence Pack

**Chip:** CHIP-MERCADO · **Mission:** MIS-004-mvp-demo · **Milestone:** M-05 (mercado radar)
**Branch:** `claude/friendly-goldberg-1e3045` (worktree, UNPUSHED) · **Base-SHA:** `9e0beb41`
**Commit-SHA:** HEAD of `claude/friendly-goldberg-1e3045` (self-referential — the authoritative
40-hex sha is reported in the CLOSED hub event and via `git rev-parse HEAD`; single commit, unpushed)

## (a) Goal
Build the missing `/mercado` "radar" screen 1:1 to the ratified design
(`.mnfs/MIS-004-mvp-demo/design/handoff/Mercado.dc.html`), in the exact mold of the merged
`/anuncios` (CHIP-ANUN-IA), `/precos` (CHIP-SIM-PARAMS), `/pedidos` (CHIP-PED-FILA) parity chips.
Was a nav stub ("Mercado — em breve"); now a live three-tab radar.

## (b) Scope boundary — FE-only, honored
Write-set is 14 files, ALL under `apps/web/src/**`. No backend / Go / OpenAPI / sdk-runtime /
migration changes. Endpoints + SDK consumed as-is.

```
apps/web/src/pages/mercado/MercadoPage.tsx          (page: header, tabs, queries)
apps/web/src/pages/mercado/RepricingTable.tsx       (Reprecificação tab)
apps/web/src/pages/mercado/OportunidadesTable.tsx   (Oportunidades tab)
apps/web/src/pages/mercado/MonitoradosTab.tsx       (Monitorados tab — honest empty)
apps/web/src/pages/mercado/oportunidades.ts         (buildOppRows join helper)
apps/web/src/pages/mercado/mercadoFormatters.ts     (money / position / collectedAt — honest null→dash)
apps/web/src/pages/mercado/marketClient.ts          (standalone IC-03 client hook)
apps/web/src/pages/mercado/{mercadoFormatters,oportunidades,MercadoPage}.test.{ts,tsx}
apps/web/src/routes/mercado.tsx                     (route wrapper)
apps/web/src/app/AppRouter.tsx                       (+ /mercado route)
apps/web/src/app/Header.tsx                          (nav: Mercado stub → NavLink)
apps/web/src/app/Header.test.tsx                     (test updated for active nav)
```

## (c) Design → impl parity matrix (vs Mercado.dc.html)
| Design element | Impl | Match |
|---|---|---|
| Header: "Mercado" 22px/700 | `MercadoPage` h1 `text-[22px] font-bold` | ✓ (LIVE fontWeight 700) |
| "Categoria ▾" pill (inert) | disabled button, border-border pill | ✓ |
| "Margem mín: 18%" accent pill | accent-soft/accent-ink pill, `MARGIN_MIN_PCT` (static label only) | ✓ |
| "fonte: busca pública ML · coletado …" faint | derived from `listings.as_of` (honest "—" if unknown) | ✓ honest |
| "Atualizar agora" button (inert) | disabled | ✓ |
| 3 tabs, underline-active, default Reprecificação | role=tab, border-b-2 accent when active | ✓ |
| Reprice grid `minmax(170px,1.3fr) 84px 84px 84px 66px 76px 90px 110px minmax(150px,1fr) 120px`, min-w 1020 | exact `GRID_COLS` + `minWidth:1020` | ✓ |
| Opp grid `84px minmax(170px,1.3fr) 76px 90px 110px 120px 110px minmax(150px,1fr) 120px`, min-w 960 | exact | ✓ |
| Mon grid `96px minmax(180px,1.4fr) minmax(150px,1.1fr) 120px 130px 100px 96px`, min-w 920 | exact | ✓ |
| Opp footer "ordenado por …" | "ordenado por concorrência observada" (truthful — sort key IS n_sellers) | ✓ honest |
| Legend line under each table | present all three | ✓ |
| Monitorados filter chips + "+ Monitorar…" | rendered, inert | ✓ |

Sanctioned honesty deviations from the mock (ADR-17 overrides the demo mock, as in sibling chips):
design's "ordenado por potencial de margem" → "ordenado por concorrência observada" (a margin sort
would need a fabricated margin); Monitorados demo rows → honest empty state (no snapshot pipeline).

## (d) Data sources (existing endpoints — no shapes invented)
- **Reprecificação** → `client.listListings({installation_id})` → `ListingReadModel`; MEU PREÇO=`price`,
  MENOR CONC.=`market_signal.min_valid`, MEDIANA=`market_signal.median`, POSIÇÃO=`market_signal.position`,
  VENDAS 30D=`sales_30d`. Single query, no join.
- **Oportunidades** → `client.listCatalogProductFacts()` × `marketClient.listMarketAggregates` +
  `listMarketVerdicts` (IC-03). Only products with an `OK`-status aggregate surface (`buildOppRows`).
  SKU=`internal_product_id` (ERP codprod, the house identifier), PRODUTO=`description`, CUSTO=`cost.amount`,
  MEDIANA ML=`aggregate.median`, **CONCORRENTES=`aggregate.n_sellers`** (distinct competing sellers,
  deduped-by-seller per the IC-03 contract — NOT the raw `n_offers`, which would overstate competition).
  Rows ordered by `n_sellers` desc (the same value shown in CONCORRENTES → the footer claim is literally
  true), tie-break SKU. `n_offers` is not carried on the row at all.
- **Monitorados** → NO backing endpoint exists (watchlists / variação-7d need a daily-snapshot store not
  built; design API-REF flags it ⚠️). → honest empty state + inert chrome.

**Pagination / completeness (CORRECTIVE #1 — hub refuter finding).** Both `/listings` and
`/catalog/products` are cursor-paginated (`ListingPage.next_cursor` / `CatalogProductFactPage.next_cursor`,
`sdk-runtime/src/index.ts:232,364`). The radar must scan the WHOLE catalog — a partial first page
silently under-reports counts and drops rows past the page, i.e. a dishonest count (ADR-17). Decision:
- **Cursor-walk both endpoints to exhaustion** (`walkAll`, `MercadoPage.tsx`) under a bounded ceiling
  `PAGE_SIZE=100 × MAX_PAGES=20 ⇒ ≤2000 items/tab`. `walkAll` stops on `next_cursor == null`; returns
  `truncated=true` iff the ceiling is hit with pages still remaining.
- **Honest count:** Reprecificação count = backed `getListingsSummary().total` (falls back to walked
  length if the summary read fails — never a partial page presented as complete); Oportunidades count =
  rendered rows.
- **No silent cap:** when a walk hits the ceiling, an honest `TruncationNote` renders ("Varredura
  limitada aos primeiros 2.000 itens — refine…"), in BOTH the populated and the empty branch of each tab
  (ADR-17).
- **Bounded query fan-out:** the codprod list is chunked to the `/market` `MaxReadIDs=200` cap
  (`chunk`, `oportunidades.ts`); per-chunk aggregates/verdicts are concatenated in fact order (via
  `Promise.all().flatMap`, order guaranteed by the JS spec), preserving `buildOppRows`' positional
  fact↔verdict alignment across chunks.

## (e) ADR-17 + D-57 honesty
- **Honest "—" (UnknownValue)** for every unbacked value: Reprice MARGEM ATUAL / SE IGUALAR MENOR /
  SUGESTÃO; Opp **MARGEM EST.** (operational margin is M-07-owned — a naive `(median−cost)/median`
  omits ML commission, freight, DIFAL, so it is NEVER shown as a value/pill/sort), VENDAS LÍDER 30D
  (no snapshot), VEREDICTO recommendation label (`verdict_label` always null pre-M-07 per market.ts),
  plus POSIÇÃO/VENDAS-30D when the field is null. No fabricated medians/margins/counts/timestamps.
  Tab counts are LIVE (not the design's 1.284/112/7).
- **D-57 inert:** Aplicar, Simular, Atualizar agora, Criar anúncio, + Monitorar…, Categoria ▾, filter
  chips are all `disabled`. Zero live Mercado Livre write; no `collectMarketPriceIntel`/PUT/POST anywhere;
  only mutation is `setTab` + read-only `refetch()`.
- **insufficient_market / no_price_evidence:** aggregate `status !== "OK"` rows are excluded from
  Oportunidades (no crash, honest empty); verdict evidence-state rendered as an honest note.

## (f) Verification
- **vitest** (apps/web): `src/pages/mercado` + touched `src/app/Header.test.tsx` + `src/app/AppRouter.test.tsx`
  → **48/48 PASS** (mercadoFormatters 4, oportunidades 5, MercadoPage 13, Header 7, AppRouter 19).
  Covers: EXEMPLO-IO 90008 golden (buildOppRows join + OportunidadesTable render), honest "—" rendering,
  CONCORRENTES=n_sellers (asserts 5 shown / 6-offer absent), sort-by-competition (discriminating fixture),
  footer-copy lock, tab switching, inert-action asserts (no `collectMarketPriceIntel`), theme tokens,
  and (CORRECTIVE #1) the cursor-walk + truncation + chunk paths below.
- **Regression guards proven load-bearing (revert → must-fail):**
  - `not.toHaveProperty("nOffers")` → RED when `nOffers` reintroduced to OppRow (offer-mislabel guard).
  - `Criar anúncio` `toBeDisabled()` → RED when `disabled` removed (D-57 inert guard).
  - sort test `["99999","90008"]` → RED when sort degraded to SKU-only (n_sellers-sort guard; fixture is
    discriminating — higher-seller item has the alphabetically-later SKU + fewer offers).
  - footer `/ordenado por concorrência observada/` → RED when copy changed (sort-claim-truthful guard).
  - **(CORRECTIVE #1) listings cursor-walk** `toHaveBeenCalledTimes(2)` + page-2 rows present → RED when
    `walkAll` reverted to a single-page fetch (proved: 2 walk tests both FAIL on `page < 1`).
  - **(CORRECTIVE #1) catalog cursor-walk + full-codprod batch** `listMarketAggregates(["90008","90009"])`
    → RED when the catalog walk stops at page 1 (page-2 codprod missing from the batch).
  - **(CORRECTIVE #1) chunk order/boundary** `chunk([a..e],2)===[[a,b],[c,d],[e]]` → RED when `chunk`
    degraded to a single batch (`return [arr]`).
  - **(CORRECTIVE #1) truncation surfaced in the EMPTY reprice branch** ("Varredura limitada…" present with
    zero rows) → RED when the empty branch reverts to a bare `EmptyState` (the exact asymmetry the hub
    refuter caught; proved FAIL on revert).
  - **(CORRECTIVE #1) summary-failure fallback** (count = walked length, no crash) → RED when
    `getListingsSummary(...).catch(()=>null)` reverts to an unguarded call.
  - **(CORRECTIVE #1, re-gate round 2) Oportunidades truncation banner** (populated + empty branch) →
    both RED when `oppTruncated` is forced `false` (endless-catalog-cursor fixtures drive the opp walk
    to the ceiling; a 2nd chip re-gate refuter caught this branch as unguarded — reprice was covered but
    the symmetric opp plumbing was not).
- **Accepted residual (non-blocking):** the `>200`-codprod multi-chunk path is proven at the `chunk()`
  unit level + single-chunk page integration; the page-level `Promise.all().flatMap` concatenation order
  is guaranteed by the JS spec (Promise.all preserves input order). No 201-item page fixture was added.
- **tsc --noEmit** (apps/web): **0 errors in write-set** (`grep -iE "mercado|Header.tsx|AppRouter.tsx"`
  of tsc output = empty). Remaining tsc errors are the repo's known pre-existing type-only baseline in
  files this chip never touched — cited per profile, not re-proved.
- **LIVE render** — throwaway vite `:5241` (strictPort 127.0.0.1, worktree-bound; NOT 5174/8080), driven
  via DOM + computed-style (screenshots infra-dead F-ENV-10), torn down after:

LIVE-VERIFIED: /mercado renders at :5247 against the live API (installation
inst-mercado_livre-d373dc64…, METALNOBREACABAMENTOS). Reprecificação → honest count **34** (from
getListingsSummary.total), 34 real listings (MLB4735324525 R$ 69,90 MEU PREÇO/MENOR/MEDIANA from
market_signal; MLB4735378521 MENOR R$ 965,00 / MEDIANA R$ 1.681,11), M-07 columns honest "—"; no
truncation banner (34 < 2000 = full walk). Oportunidades → **CORRECTIVE #1 proven live: 7 real
opportunities now surface** (90034 CONCORRENTES 96, 90031/38, 90010/35, 90033/27, 90032/20, 90008/11,
90009/5) — under the OLD single-page (limit 50) fetch these codprods sat past page 1 and the tab showed
**0** (the silent-truncation bug); the catalog cursor-walk (live: `/catalog/products?limit=100`
returns `next_cursor` ⇒ multi-page) now reaches them. Sorted by CONCORRENTES desc (96→5, footer
"ordenado por concorrência observada" literally true); MARGEM EST./VENDAS LÍDER 30D/VEREDICTO all honest
"—" + "mercado observado"; no truncation banner (this catalog ≤2000 ⇒ full walk, no false cap claim).
Monitorados → "Nenhum monitoramento ativo." + "+ Monitorar…" disabled. Nav "Mercado" active link, no
"em breve". Pagination change touches no tokens (light↔dark unchanged from prior round). All action
controls disabled = zero live ML write.

## (g) Gate — chip-side P6 DUAL (cold Opus + adversarial sonnet refuter), reading the WHOLE module
Four rounds; each refuter finding was a REAL defect, remediated, then re-gated to a clean agreement:

| Round | Cold (Opus) | Refuter (sonnet) | Outcome |
|---|---|---|---|
| 1 | **FAIL** — MARGEM EST. was a client-fabricated, decision-colored (green/amber/red vs 18%), sort-driving margin (ADR-17 NON-NEGOTIABLE; margin is M-07-owned) | NO-REFUTATION (missed it) | Disagree → remediated: MARGEM EST. → honest "—"; removed `estimateMarginPct`/`marginTone`/`MARGIN_PILL_CLASS`/`formatPercent`/`OppRow.marginPct`; sort no longer margin-based |
| 2 | PASS | **REFUTED** — CONCORRENTES rendered `n_offers` (raw offers) under a "competitors" header, overstating competition; `n_sellers` was fetched but discarded | Disagree → remediated: CONCORRENTES renders `n_sellers` |
| 3 | PASS | **REFUTED** — footer "ordenado por concorrência observada" contradicted an `n_offers` sort (file's own vocabulary defines concorrência = n_sellers) | Disagree → remediated: sort key → `n_sellers` desc; `nOffers` removed from OppRow; footer now literally true; footer-copy test added |
| 4 | **PASS** (6/6, both prior fixes confirmed) | **NO-REFUTATION** (no new blocking defect beyond the 3 fixed + 1 adjudicated) | **AGREEMENT** |

Adjudicated non-blocking (raised round 4, dismissed on merits by both): SKU column shows
`String(internal_product_id)` (ERP codprod, the house identifier), not the nullable `reference` — a real
backed value, honest, not an ADR-17 issue.

### CORRECTIVE #1 — hub gate caught a real blocker the chip's own 4 rounds missed

The **hub** P6 dual gate DISAGREED (cold PASS, refuter REFUTED) on one blocking defect the chip's own
gate never exercised: **silent page-1 truncation.** `MercadoPage` fetched exactly one 50-item page of the
cursor-paginated `/listings` and `/catalog/products`, so counts under-reported and rows past the page were
silently dropped. It evaded all 4 chip rounds because the live-verify DB happened to have 34 listings / 0
opps — the cap was never exercised (verification is bounded by what you exercise; the exact blind spot).

Remediation (this round): cursor-walk both endpoints under a bounded ceiling, honest backed count +
summary-failure fallback, honest `TruncationNote` on ceiling-hit (both branches, both tabs), `/market`
codprod batch chunked to 200 with order-preserving concatenation. See (d)/(f).

Chip-side **re-gate** on the corrected module (cold `harness:gate-reviewer` + adversarial refuter,
read-only, whole module) then DISAGREED again — both found real complementary gaps: refuter REFUTED (the
reprice **empty-state** branch dropped the banner — asymmetric with the opp branch, a live instance of the
same silent-cap class); cold FAIL (the `TruncationNote` ceiling path had **zero test coverage**). Both
remediated: empty-branch banner made symmetric, and 3 truncation/fallback tests added (populated-ceiling,
empty-ceiling, summary-failure) — each proven load-bearing (revert → RED). Live re-verified: the walk now
surfaces **7 real opportunities that were previously hidden** (see LIVE-VERIFIED).

Chip-side re-gate round 2 then DISAGREED once more — cold PASS, refuter REFUTED: the **Oportunidades**
truncation banner (both branches) had zero test coverage (reprice was covered, the symmetric opp
plumbing was not — a revert would go undetected). Remediated: 2 opp-tab ceiling tests added
(populated + empty), proven load-bearing (force `oppTruncated=false` → both RED).

Chip-side dual gate = **AGREEMENT** on the corrected artifact after all remediations above (final
read-only cold + refuter pass over the whole module, every re-gate finding fixed and must-fail
proven). Verified across both rounds: ADR-17 honest "—" on every unbacked field, honest counts + surfaced
truncation (never silent), CONCORRENTES = n_sellers (not n_offers), sort n_sellers desc with a
literally-true footer, chunk order-preserving fact↔verdict alignment, D-57 inert (no live ML write /
`collectMarketPriceIntel`), FE-only scope, and discriminating (non-theater) regression tests.

P6-DUAL-GATE: AGREEMENT

(Chip-side gate marker — records the chip's own dual-gate result. Milestone-level promotion,
merge, post-merge ladder and P7 browser QA remain HUB-OWNED.)

## (h) Bindings honored
No push, no merge, no server boot / :8080 bind / .env load, no dep ritual (npm ci = profile env-prep),
no reset/revert/stash/clean of unknown state (guard reverts were Edit→Edit, never `git checkout`).
Merge + post-merge ladder + P7 browser QA are HUB-OWNED.
