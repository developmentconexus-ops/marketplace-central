# CHIP-SIM-MATRIZ — Slice Plan (cold Opus planner)

**Scope (BINDING):** `apps/web/src/pages/precos/**` ONLY. No SDK/OpenAPI edits (FROZEN). No new endpoints.
**Goal:** Make the design's PRODUCT MATRIX the main surface of `/precos`; row-click opens the EXISTING simular panel (DecompositionPanel + SolverPanel) wired to the selected row. Preserve Parâmetros + DIFAL drawers.

---

## 0. Facts confirmed by reading (do not re-prove)

- `PricingPage.tsx` current layout = `grid grid-cols-[280px_1fr]`: left `<aside>` product-list picker + right `<section>` (preço input, modalidade, `region-decomposicao`, `region-solver`, `region-comparacao`, `region-aplicar`, `region-cenarios`). Two drawers (`ParamsDrawer`, `DifalDrawer`) rendered at root.
- `selected` derives from `selectedId`: **null ⇒ defaults to `products[0]`** (line 117); explicit-but-absent id ⇒ `null` (honest empty). This default-to-first is load-bearing for the existing "renders shell regions on mount" test.
- `CatalogProductFact` (sdk `index.ts:202`): `internal_product_id:number`, `description`, `manufacturer_reference`, `current_price.amount:string|null` (**NOSSO PREÇO**, canonical dot-decimal e.g. `"89.00"`), `cost.amount:string|null` (**CUSTO**).
- `pricingDecompose(PricingCalcInput)` → `PricingDecomposeResponse.decomposition` with `margem_valor:string|null` (**retorno/un**), `margem_pct:string|null` (**chip %**), `custo:string|null`. `comissao_pct` MUST stay OMITTED (resolver chain). Row input = `{ preco: current_price.amount, modalidade, product_id }` — current_price is already dot-decimal, **no ptBr conversion** needed for rows.
- Market: `MarketAggregatesClient.listMarketAggregates(codprods:string[])` → `MarketPriceIntelAggregate[]` with `product_id:string`, `median:Money|null`, `status ∈ OK|INSUFFICIENT_MARKET|NO_PRICE_EVIDENCE`. **NO rank/position field** ⇒ OMIT "19º/23", keep median as PREÇO MERCADO. Standalone client seam is injectable (MarketComparison pattern: `client?` prop, defaults `createMarketPriceIntelClient({baseUrl: apiBaseUrl()})`).
- UI kit (`@marketplace-central/ui`): `MarginChip({ marginPct:number|null, thresholds?:{healthy,tight} })` — renders `—` (band `unknown`) when null, else `${marginPct}%` colored by bands (default 18/10; feed profile `limiar_verde_pct`/`limiar_amarelo_pct`). `UnknownValue({ hint? })` → `—` with title.
- `DecompositionPanel` + `SolverPanel` signatures unchanged — REUSE, do not rebuild.

---

## 1. Component decomposition (DECISION)

### New file: `PricingMatrix.tsx`
Presentational table + its own two data lanes. Props:

```ts
interface PricingMatrixProps {
  products: CatalogProductFact[];
  selectedId: number | null;          // for aria-pressed / row highlight (selected derives in page)
  onSelect: (id: number) => void;     // row click
  modalidade: string;                 // page modalidade → per-row decompose input
  profile: PricingCalcProfile;        // thresholds for MarginChip
  marketClient?: MarketAggregatesClient; // INJECTABLE (tests); default standalone IC-03 client
}
```

**Data lane A — market (ONE call, page-wide):**
`useQuery(["market","matrix-aggregates", codprods.join(",")], () => marketClient.listMarketAggregates(codprods))` where `codprods = products.map(p => String(p.internal_product_id))`. Result indexed into `Map<string, MarketPriceIntelAggregate>` by `product_id`. **Single fan-in call for all rows** (requirement: ALL codprods in one call).

**Data lane B — per-row margin (N calls, `useQueries`):**
`useQueries({ queries: products.map(p => ({ queryKey: ["pricing","matrix-decompose", p.internal_product_id, p.current_price.amount, modalidade], queryFn: () => client.pricingDecompose({ preco: p.current_price!.amount, modalidade, product_id: p.internal_product_id }), enabled: p.current_price.amount !== null })) })`. `client` = `useClient()` (central seam — already mocked via `ClientContext` in page tests). `comissao_pct` OMITTED.
- Row `i` margin = `results[i].data?.decomposition`; `margem_valor`/`margem_pct` null (e.g. SEM_CUSTO) or query disabled ⇒ MARGEM cell = `<UnknownValue hint="margem: M-07" />`.

**Client seam for tests:** market = injectable `marketClient` prop (MarketComparison pattern); decompose + catalog = `ClientContext` `useClient()` (mock the module, as `PricingPage.test.tsx` already does). Golden test mocks `../../app/ClientContext` and injects `marketClient`.

### Row → cell mapping
| Column | Source | Unknown/honest state |
|---|---|---|
| SKU | `manufacturer_reference ?? #id` | — |
| DESCRIÇÃO | `description ?? manufacturer_reference ?? #id` | — |
| CUSTO | `cost.amount` | `null ⇒ UnknownValue "—"` |
| NOSSO PREÇO | `current_price.amount` | `null ⇒ UnknownValue "—"` |
| PREÇO MERCADO | aggregate.median WHEN `status==="OK"` | else `—` (missing aggregate OR non-OK) — **rank OMITTED** |
| MARGEM | `margem_valor` (mono retorno/un) + `MarginChip(Number(margem_pct), thresholds)` | margem null / no decompose ⇒ `UnknownValue hint="margem: M-07"` |
| VEREDICTO | map `aggregate.status` | see below |

**VEREDICTO map (price-evidence ONLY — categorical MARGIN verdict NOT invented, that's M-07):**
`OK→"OK"`, `NO_PRICE_EVIDENCE→"SEM_EVIDENCIA"`, `INSUFFICIENT_MARKET→"MERCADO_INSUFICIENTE"`, **missing aggregate → "SEM_EVIDENCIA"** (honest closest).

**"novo" tag** = product has no ML listing. Deferred to Slice B (per-row `listListingsByProduct` via `useQueries`, gated on `installationId`, "novo" when `groups[0].listings` empty). Kept OUT of Slice A so the golden test stays focused. Fan-out flagged in Risks.

**Design tokens (DESIGN-REFERENCE):** wrapper `overflow-x:auto` + inner `min-width:900px`; header row `--surface2`, 11px, letter-spacing .04em, uppercase; body rows 12.5px; margin chips already colored by `MarginChip`. `data-testid="pricing-matrix"`, each row `data-testid="matrix-row-<id>"`.

---

## 2. PricingPage restructure (DECISION)

Replace the `grid grid-cols-[280px_1fr]` (aside picker + section) with **matrix main + conditional 380px side panel**:

```
<div className="flex gap-4 p-4">
  <div className="flex min-w-0 flex-1 flex-col gap-4">   // MAIN
    <header/>                     // h1 + params-trigger + difal-trigger  (UNCHANGED)
    {profileQuery.isError ? alert : null}
    <PricingMatrix               // NEW main surface
      products={products}
      selectedId={selectedId}
      onSelect={setSelectedId}
      modalidade={modalidade}
      profile={profile}
    />
  </div>
  {selected !== null && (        // SIDE PANEL — the PRESERVED existing section, 380px
    <aside className="w-[380px] shrink-0 flex flex-col gap-4" aria-label={`Simular · ${productLabel(selected)}`}>
      {productMissing ? notice : null}
      <preço input + modalidade block/>   // UNCHANGED
      <div data-testid="region-decomposicao"> <DecompositionPanel .../> </div>
      <div data-testid="region-solver"> <SolverPanel .../> </div>
      <div data-testid="region-comparacao"> <MarketComparison .../> </div>
      <div data-testid="region-aplicar"> <ApplyPriceAction .../> </div>
      <div data-testid="region-cenarios"> <ScenariosPanel .../> </div>
    </aside>
  )}
  <ParamsDrawer/> <DifalDrawer/>  // UNCHANGED at root
</div>
```

**State: nothing moves.** `selectedId`, `modalidade`, `precoInput`, `paramsOpen`, `difalOpen`, all queries/mutations stay in `PricingPage`. Matrix owns only its two internal query lanes. Row click = `onSelect(id) => setSelectedId(id)`.

**Why regions survive on mount:** `selected` still default-to-first (selectedId null ⇒ products[0]), so `selected !== null` is true on mount ⇒ side panel + all `region-*` testids render ⇒ existing "renders shell regions" test passes. The `productMissing` empty-state path is unchanged.

**Test-safety edit (Slice B):** add `vi.mock("./PricingMatrix", () => ({ PricingMatrix: () => <div data-testid="pricing-matrix-stub" /> }))` to `PricingPage.test.tsx` — same sanctioned pattern already used for `MarketComparison`/`SolverPanel`/`ApplyPriceAction`/`ScenariosPanel` (stub child panels to avoid their side effects: matrix's standalone market call + N decompose calls). All existing region/drawer/scenario/pt-BR assertions keep passing because they target the preserved side panel, not the matrix.

---

## 3. Slice cards (2 slices)

### Slice A — golden test + matrix component (DAY 1, test-first)
- **Goal:** `PricingMatrix.tsx` renders columns SKU·DESCRIÇÃO·CUSTO·NOSSO PREÇO·PREÇO MERCADO·MARGEM·VEREDICTO for a multi-product list, with honest unknowns. Rank omitted.
- **write_set:** `apps/web/src/pages/precos/PricingMatrix.tsx` (new), `apps/web/src/pages/precos/PricingMatrix.test.tsx` (new).
- **failing-test-first (EXEMPLO-IO golden, 3 products in one render):**
  - (a) custo `120` preço `189`, aggregate median `169`/status OK/n_sellers 8, decompose→`margem_valor` + `margem_pct` non-null ⇒ row shows retorno + MarginChip `%` + PREÇO MERCADO `169` + VEREDICTO `OK`.
  - (b) same custo/preço, aggregate status `NO_PRICE_EVIDENCE` ⇒ PREÇO MERCADO `—`, VEREDICTO `SEM_EVIDENCIA`, MARGEM still shown (decompose from custo/preço).
  - (c) `cost.amount=null`, decompose returns `margem_valor=null`/`margem_pct=null` ⇒ CUSTO `—`, MARGEM `UnknownValue` `—`.
  - Assert market queried **once** with all three codprods: `listMarketAggregates(["...","...","..."])`; assert rank string `/º/` **absent**.
  - Mock `../../app/ClientContext` `useClient` → `{ pricingDecompose: vi.fn((req)=> per-product response) }`; inject `marketClient={{ listMarketAggregates: vi.fn() }}`. Render helper = `QueryClientProvider` (+ `MemoryRouter` only if needed).
- **done:** matrix renders all 7 columns; unknowns honest (ADR-17); rank omitted; single market call; per-row decompose keyed & `comissao_pct` absent; MarginChip fed profile thresholds. New test GREEN.
- **validate:** `pnpm --filter @marketplace-central/web test -- PricingMatrix` ; `pnpm --filter @marketplace-central/web exec tsc --noEmit` (baseline = main's 10 type-only, cite don't re-prove).

### Slice B — wire matrix into PricingPage + "novo" tag
- **Goal:** matrix becomes the main surface; row-click opens the preserved 380px simular panel; "novo" tag on rows with no listing; existing tests stay green.
- **write_set:** `apps/web/src/pages/precos/PricingPage.tsx` (restructure §2), `apps/web/src/pages/precos/PricingPage.test.tsx` (add PricingMatrix stub + 1 wiring assertion), `apps/web/src/pages/precos/PricingMatrix.tsx` (add "novo" listing lane).
- **failing-test-first:** in `PricingPage.test.tsx` add: matrix stub present on mount; region-* + drawers still present; (new) clicking a stub-exposed row seam calls `setSelectedId` → side panel reflects product. In `PricingMatrix.test.tsx` add a "novo" case (product whose `listListingsByProduct` → empty listings ⇒ `novo` tag; product with a listing ⇒ no tag). Gate listing lane on injected client so test controls it.
- **done:** matrix main + conditional 380px panel; row click drives selection; "novo" honest (no listing) — never fabricated; all pre-existing PricingPage assertions (deep-link params, save profile, scenario reload valid/missing, pt-BR normalize, DIFAL table) GREEN.
- **validate:** `pnpm --filter @marketplace-central/web test -- precos` ; `tsc --noEmit`.

> Kept to 2 slices (bounded chip). If Slice B's "novo" listing fan-out proves heavy in QA, split "novo" into a follow-up SPLIT-REQUEST to the hub rather than expanding this chip.

---

## 4. Verification map (criterion → assertion)

| Requirement | Verified by |
|---|---|
| P0 matrix is main surface; 7 columns | Slice A golden test: column headers + 3 rows rendered; Slice B: `pricing-matrix` occupies main col |
| Row-click → existing simular panel (380px), reuse not rebuild | Slice B: click row seam → `selectedId` set → `region-decomposicao`/`region-solver` reflect product; DecompositionPanel/SolverPanel imports unchanged |
| Parâmetros + DIFAL drawers preserved | existing `params-trigger`/`difal-trigger`/`params-drawer`/`difal-drawer` tests unchanged & green (Slice B) |
| Rank OMITTED, market price kept | Slice A: assert `/\dº\/\d/` absent; PREÇO MERCADO = median when OK |
| Same endpoints; market in ONE call | Slice A: `listMarketAggregates` called once with all codprods |
| Per-row margin via pricingDecompose, comissao_pct omitted | Slice A: decompose mock asserted `not.toHaveProperty("comissao_pct")` |
| P1 VEREDICTO = price-evidence map only | Slice A (a/b) + INSUFFICIENT_MARKET case → `OK/SEM_EVIDENCIA/MERCADO_INSUFICIENTE` |
| Categorical MARGIN verdict NOT invented → UnknownValue hint "margem: M-07" | Slice A (c): MARGEM `UnknownValue` with `title="margem: M-07"` |
| ADR-17 unknown ≠ 0 | (a/b/c): CUSTO `—`, PREÇO MERCADO `—`, MARGEM `—` — never `0` |
| EXEMPLO-IO (a)(b)(c) | Slice A golden test = literal three cases |
| "novo" when no listing | Slice B: listing-empty ⇒ tag; listing-present ⇒ no tag |

---

## 5. Risks & mitigations

1. **Per-row decompose fan-out (N `pricingDecompose`).** Bounded by page `limit:50`; demo catalog ~10 products. React-query parallelizes + caches by key; `enabled` gated on `current_price!=null`. Mitigation: keep the page's existing 50 cap; note as perf follow-up, do NOT add batching (no batch endpoint; endpoints FROZEN).
2. **"novo" listing lane adds a 2nd N fan-out** (`listListingsByProduct` per row). Isolated to Slice B, gated on `installationId`. If QA flags cost, SPLIT-REQUEST "novo" to a follow-up chip — matrix core (Slice A) is independently shippable.
3. **Breaking existing PricingPage tests.** Avoided by (a) preserving the entire `<section>` as the side panel with all `region-*` testids, (b) keeping `selected` default-to-first so regions render on mount, (c) adding a `PricingMatrix` `vi.mock` stub in `PricingPage.test.tsx` (same pattern as the 4 existing child stubs) so the matrix's market/decompose side effects don't hit the page test.
4. **Standalone market client hits network in tests if un-stubbed.** Golden test injects `marketClient`; page test stubs the whole matrix. No live IC-03 dependency in unit tests.
5. **`current_price.amount` is dot-decimal already** — do NOT run `ptBrMoneyToDot` on row values (that's for operator free-text only). Confirmed from SDK fixture (`"89.00"`).
