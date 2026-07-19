# F-01 produto-detalhe — Batch Implementation Plan (P2)

```yaml
id: F-01
type: feature-plan
parent: M-06
base_sha: 89de2fef824fd61093ea2c4b340dbf1c5a759778
branch: chip-m06-produto
planner: cold Opus subagent (contingency lane §12, codex quota outage)
created: 2026-07-19
```

## Codebase-map findings (brief-vs-base deltas — all handled honest ADR-17)

- **No per-domain SDK client files.** All catalog/links/listings reads on the barrel
  `@marketplace-central/sdk-runtime` `index.ts`; market on standalone `market.ts`. Feature brief's
  `catalog`/`productLinks.ts`/`listings.ts` client-file assumption is wrong.
- **Market client NOT in `useClient()`** — replicate `pages/precos/MarketComparison.tsx:19-63`
  local-client pattern (`apiBaseUrl()` + `useMemo(createMarketPriceIntelClient)`), test-injectable via
  `client?` prop. In-scope (owned page dir), no `app/**` edit.
- **`MarketVerdict.verdict_label` ALWAYS null** (M-07-owned; OpenAPI:7561, market.ts:123). Margin
  categorical (saudavel/viavel/apertado/nao_vale) unimplemented at base ⇒ rendered `UnknownValue`
  hint "veredicto de margem — M-07". Honest verdict = price_evidence_status + faixa + evidência +
  custo + blocking_state. → ESCALATION FINDING-1 to hub.
- **No físico/reservado/disponível split** — only single `stock_quantity` CanonicalNumericSourceFact
  (index.ts:199). Contract C04 anticipates this: físico = stock_quantity; reservado + disponível =
  UnknownValue "—". → ESCALATION FINDING-2 to hub.

## Cross-cutting contracts (bind every slice)

- Route param `productId` = catalog `internal_product_id` as string. `getCatalogProduct(Number(productId))`
  (number sig index.ts:2045); market/listings use raw string. Parse once in page root;
  `!Number.isFinite(numericId) || productId===""` ⇒ not-found ErrorState (== 404 branch).
- Barrel client from `useClient()` (`../app/ClientContext`) — import, never create. `getJson` throws
  `{status,error}` on non-2xx (index.ts:1674); `error.status===404` = not-found.
- Market client local (mirror MarketComparison), throws `{status,error}` on non-2xx.
- `installationId` from `useInstallation()` (`../app/InstallationContext`); produto route under
  `<InstallationProvider>` (AppRouter.tsx:37,43). Only Anúncios tab needs it.
- Query keys = inline literal arrays in owned files (no web-query factory edit — `packages/**`
  FORBIDDEN); Anúncios reuses exported read-only `listingsQueryKeys.byProduct` so invalidation prefix
  bites.
- Every null numeric/identity ⇒ `<UnknownValue/>` (renders `—`, optional `hint`), never `""`/`0`.
- Design: paper+green className tokens (`text-ink/muted/faint`, `bg-surface`, `border-border`,
  `text-accent`, `var(--font-mono)`), reuse from AnunciosPage/MarketComparison. NO hardcoded hex.
- ui exports confirmed: EmptyState, ErrorState, LoadingState, UnknownValue, SurfaceCard,
  StatCard{label,value,sub?,className?}, Badge, Button, FreshnessIndicator{asOf}. No shared Tabs —
  hand-roll role="tablist"/"tab" (mirror AnunciosPage.tsx:198-215).
- Test harness: `vi.mock` ClientContext + InstallationContext, QueryClientProvider(retry:false) +
  MemoryRouter; widgets take `productId` (+ market `client?`) as props.
- Collection response shape: `{status:"COMPLETED"|"PARTIAL", "decisões":[…], contagens, causas}`
  (accented key `"decisões"`).

## Slice cards (serialize S1→S7; ProdutoPage.tsx is the hot composition file)

### S1 — Scaffold: route swap + local market client + tab-state codec + page shell
- validation_kind: unit-vitest + build
- commands: `npm run test --workspace apps/web -- src/pages/produto/ProdutoPage.test.tsx src/pages/produto/productQueryState.test.ts` · `npm run test --workspace apps/web -- src/app/AppRouter.test.tsx` · `npm run build --workspace apps/web`
- write_set: `routes/produto.tsx` (swap → `<ProdutoPage/>`), `pages/produto/ProdutoPage.tsx` (new shell: param parse, not-found guard, tablist, 3 placeholder panels), `pages/produto/productQueryState.ts` (new: `ProdutoTab`, `parseProdutoTab`, `applyProdutoTab`), `pages/produto/marketClient.ts` (new: `apiBaseUrl()` + `useProdutoMarketClient(injected?)`), `pages/produto/ProdutoPage.test.tsx`, `pages/produto/productQueryState.test.ts`, `app/AppRouter.test.tsx` (GRANTED — lines 123-125 assertion only → assert new page e.g. `findByRole("tab",{name:"Veredicto"})`).
- failing_test_first: parseProdutoTab default veredicto for missing/invalid, round-trips anuncios/estoque; ProdutoPage renders 3 tabs, deep-link `?tab=estoque` selects Estoque, tab click updates `?tab=`, invalid/empty productId ⇒ not-found ErrorState + /catalogo link; AppRouter updated assertion passes.
- open_questions: []

### S2 — Header widget (identity + custo + completude)
- validation_kind: unit-vitest
- commands: `npm run test --workspace apps/web -- src/pages/produto/ProdutoHeader.test.tsx`
- write_set: `pages/produto/ProdutoHeader.tsx` (new), `.test.tsx` (new), `ProdutoPage.tsx` (compose header above tablist).
- failing_test_first: null EAN/ref/ncm/cost fixture ⇒ `—` per null, "sem EAN — matching limitado a REVIEW" chip, completude chip per quality_flag; full fixture ⇒ real values + custo observed_at freshness. No literal `0`/`""` where source null.
- open_questions: []

### S3 — Veredicto box (verdict + synchronous collect + refetch, no polling)
- validation_kind: unit-vitest
- commands: `npm run test --workspace apps/web -- src/pages/produto/VeredictoBox.test.tsx`
- write_set: `pages/produto/VeredictoBox.tsx` (new), `.test.tsx` (new), `ProdutoPage.tsx` (mount veredicto panel, wire injected market client + costFact).
- detail: headline price_evidence_status; INSUFFICIENT_MARKET ⇒ exact `"{n_sellers} vendedores — mínimo 5"`; faixa min_valid/median/n_sellers; custo; evidência inputs_used rows; blocking_state incl SEM_CUSTO; margin label always UnknownValue hint M-07. Button coletar ⇒ useMutation collectMarketPriceIntel → onSuccess invalidate `["market","verdict",productId]` + `["listings","by-product"]`; disabled while isPending; label "coletando…"; NO polling/collection query; error ⇒ inline ErrorState/toast.
- failing_test_first: (1) INSUFFICIENT_MARKET n_sellers:2 ⇒ "2 vendedores — mínimo 5"; (2) OK ⇒ faixa + inputs_used rows + margin `—`+M-07 hint; (3) NO_PRICE_EVIDENCE 200 not 404 ⇒ headline, no fabricated price; (4) collect happy ⇒ disabled in-flight, both invalidations, no interval/2nd query; (5) collect error 504 ⇒ error surface, button re-enabled.
- open_questions: []

### S4 — Aba Anúncios vinculados
- validation_kind: unit-vitest
- commands: `npm run test --workspace apps/web -- src/pages/produto/AnunciosVinculadosTab.test.tsx`
- write_set: `pages/produto/AnunciosVinculadosTab.tsx` (new), `.test.tsx` (new), `ProdutoPage.tsx` (mount anuncios panel).
- detail: listListingsByProduct({installation_id, product_id}) via `listingsQueryKeys.byProduct(installationId,{product_id})`; per-listing match_status, position rank/total, price_to_win, delta_pct + FreshnessIndicator(evidence.fetched_at); null ⇒ UnknownValue; no listings ⇒ EmptyState + /vinculos link (not blank table).
- failing_test_first: group w/ 2 listings (one full signal, one null/SEM_VINCULO) ⇒ rank/price/delta/freshness for first, `—` for null; empty groups ⇒ EmptyState /vinculos; query invoked with product_id.
- open_questions: []

### S5 — Aba Estoque
- validation_kind: unit-vitest
- commands: `npm run test --workspace apps/web -- src/pages/produto/EstoqueTab.test.tsx`
- write_set: `pages/produto/EstoqueTab.tsx` (new), `.test.tsx` (new), `ProdutoPage.tsx` (mount estoque panel).
- detail: físico = stock_quantity.value + observed_at freshness; reservado = UnknownValue; disponível = DESCONHECIDO UnknownValue; value null OR missing_stock flag ⇒ "importar planilha" state (no 0); Concorrência/Pedidos/Histórico ABSENT.
- failing_test_first: value:42 ⇒ físico 42 + freshness, reservado + disponível `—`; value:null/missing_stock ⇒ importar planilha, no 0; assert 3 legacy tabs absent.
- open_questions: []

### S6 — Tab deep-link + F5 restore hardening (integration)
- validation_kind: unit-vitest
- commands: `npm run test --workspace apps/web -- src/pages/produto/ProdutoPage.deeplink.test.tsx`
- write_set: `pages/produto/ProdutoPage.deeplink.test.tsx` (new); contingent `ProdutoPage.tsx`/`productQueryState.ts`.
- failing_test_first: deep-link `?tab=anuncios` ⇒ Anúncios panel, others not; remount same URL ⇒ still Anúncios; `?tab=bogus` ⇒ Veredicto; tab switch updates location.search.
- open_questions: []

### S7 — Partial-failure isolation + light/dark computed-style QA
- validation_kind: unit-vitest + manual-computed-style
- commands: `npm run test --workspace apps/web -- src/pages/produto/ProdutoPage.partialFailure.test.tsx` · `npm run build --workspace apps/web`
- write_set: `pages/produto/ProdutoPage.partialFailure.test.tsx` (new); contingent widget files S2-S5 for isolated isError branch.
- failing_test_first: market rejects while catalog resolves ⇒ Veredicto ErrorState+retry, Header renders; catalog rejects non-404 ⇒ Header ErrorState+retry, tablist present; retry ⇒ refetch.
- open_questions: []

## Write-DAG

| File | Slices | Role |
|---|---|---|
| routes/produto.tsx | S1 | one-line swap |
| app/AppRouter.test.tsx | S1 (GRANTED, 1 assertion) | line 123-125 |
| pages/produto/productQueryState.ts | S1 create, S6 contingent | tab codec |
| pages/produto/marketClient.ts | S1 | local market client hook |
| **pages/produto/ProdutoPage.tsx** | **S1,S2,S3,S4,S5,S6(c),S7(c)** | HOT composition root |
| pages/produto/ProdutoHeader.tsx | S2 create, S7(c) | widget |
| pages/produto/VeredictoBox.tsx | S3 create, S7(c) | widget |
| pages/produto/AnunciosVinculadosTab.tsx | S4 create, S7(c) | widget |
| pages/produto/EstoqueTab.tsx | S5 create, S7(c) | widget |

Serialization: ProdutoPage.tsx is the only multi-slice file → run strictly S1→S7 (sync workers). Widget
bodies S2⊥S4⊥S5 independent; only ProdutoPage compose edits serialize. S3 invalidate prefix
`["listings","by-product"]` must match S4's `listingsQueryKeys.byProduct` key by prefix.

## Verification map (M-06-C01..C05)

- **C01** header honest nulls → `ProdutoHeader.test.tsx` (S2) + computed-style null fixture: every gap = `—`, none `0`.
- **C02** veredicto sync collect+refetch no polling → `VeredictoBox.test.tsx` (S3): disabled in-flight, both invalidations, no interval; INSUFFICIENT_MARKET exact copy; error branch. Fake timers assert no polling.
- **C03** anúncios tab + empty → `AnunciosVinculadosTab.test.tsx` (S4): signal columns; empty ⇒ /vinculos CTA.
- **C04** estoque físico + reservado/disponível `—` → `EstoqueTab.test.tsx` (S5): físico+freshness, two `—`, missing_stock ⇒ importar planilha, 3 legacy tabs absent.
- **C05** ownership clean + light/dark → S6/S7 tests + `npm run build`; `git diff --name-only 89de2fef..HEAD` ⊆ {pages/produto/**, routes/produto.tsx, app/AppRouter.test.tsx}; computed-style light+dark token resolve + a11y tablist roles; zero backend/OpenAPI/migration/sdk/governance. tsc baseline = 10 type-only (D-97), new error is ours.

## Seam-closure checklist

- Route swap routes/produto.tsx → `<ProdutoPage/>` — [owned S1]. AppRouter.tsx:43 already points ProdutoRoute, NOT edited.
- AppRouter.test.tsx assertion (123-125) — [granted S1], single assertion.
- Local market-client wiring — [owned S1]. Rule-of-three deferral: apiBaseUrl 2nd/3rd copy; shared hoist needs app/**||packages/** (FORBIDDEN) → local accepted, deferral flagged for future app-owned chip.
- Query-key additions — inline literals owned files; Anúncios reuses read-only listingsQueryKeys.byProduct. No web-query edit — [owned S1/S3/S4].
- No SDK/OpenAPI/migration/backend — composes existing reads; data gaps → UnknownValue; real missing endpoint → ESCALATION, never built. [n/a / escalation path]
- No governance/modules.json — no new package/module. [n/a]
- installationId — existing useInstallation(), route under InstallationProvider. [owned S4 read-only]
- Light+dark — className tokens, getComputedStyle both themes + a11y-tree (rasterizer broken). [owned S7/QA]
- Env — npm ci at worktree root before lanes (env prep). No server/:8080/:5174/.env — hub seam. Commit per green slice.
