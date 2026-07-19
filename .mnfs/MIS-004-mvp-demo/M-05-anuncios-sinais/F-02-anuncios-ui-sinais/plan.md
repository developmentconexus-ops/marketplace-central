# F-02 anuncios-ui-sinais — SLICE PLAN (P2 batch, cold Opus)

```yaml
id: F-02
type: feature-plan
parent: M-05-anuncios-sinais
mission: MIS-004-mvp-demo
base_sha: 8b6c4b3093f9465cd3b91209b054af4fa702171a
planner: P2 BATCH PLANNER (cold Opus)
lane: contingency §12 (implementers = sonnet)
slices: 6
depends_on: F-01 (listings.ts SDK surface + OpenAPI market_signal fields must be merged first)
internal_dag: S1 -> S2 -> S3 ; S4 (after S1) ; S5 (after S2/S3) ; S6 retheme (after S2-S5)
```

## ⚠️ SDK-LISTINGS-M05 GRANT OVERRIDE (BINDING — supersedes EVERY `withListingSignals` / `ListingWithSignal` / `listings.ts` reference in the slice cards below)

Hub grant **SDK-LISTINGS-M05** VOIDED the standalone `listings.ts` + `withListingSignals` wrapper — they DO NOT EXIST. The F-01 signal surface is now INLINE on the base SDK types in `@marketplace-central/sdk-runtime` (index), and F-01 is committed on THIS branch:
- `ListingReadModel.market_signal?: ListingMarketSignal | null` and `.signal_status?: SignalStatus` — already on the return of `client.listListings` / `client.listListingsByProduct` / `client.getListing`.
- `ListingSummaryExceptions.sem_vinculo? / abaixo_custo? / sem_evidencia?: number | null`.
- `ListingException` union already includes `sem_vinculo | abaixo_custo | sem_evidencia`.
- New exported types: `SignalStatus`, `ListingMarketSignal`, `ListingSignalEvidence`.

THEREFORE in every slice below:
- **Do NOT call/import `withListingSignals`.** Use `const client = useClient()` DIRECTLY — items are already typed with the optional signal fields. (`app/ClientContext.tsx` stays FORBIDDEN + untouched; no wrapper is needed to reach the fields.)
- **Do NOT import** `ListingWithSignal` / `ListingSummaryExceptionsWithSignals` / anything from `./listings`. Import `ListingReadModel`, `ListingException`, `SignalStatus`, `ListingMarketSignal`, `ListingSummaryExceptions` from `@marketplace-central/sdk-runtime`.
- Signal fields are OPTIONAL (`?`) at the type level — render DEFENSIVELY: treat a missing `signal_status` as the honest absent state (NO_PRICE_EVIDENCE / "—"), never fabricate a number.
- **F02-S1** drops the "wrap client in withListingSignals" step — it reduces to the query-state `exceptionValues` + `exceptionLabels` additive extension ONLY.
- **F02-S3** types `anunciosSummaryQuery` via the inline `ListingSummaryExceptions` (3 optional counters), NOT `ListingSummaryExceptionsWithSignals`.

Everything else in each card (write_set files, failing-test-first, done_criteria, ADR-17 honest-state rules, W1 regression guard) STANDS.

---

## Scope recap (from feature.md + validation-contract C03 + IC-05)

Extend the LIVE W1 `AnunciosPage` additively: VS MERCADO column, clickable exception chips (URL filter),
agrupar-por-produto toggle (by-product endpoint), drawer evidence section, and retheme to M-03 token classes.
Zero regression on W1 behavior (deep-links C08/C09, filters/pagination/drawer/refresh). No fetch outside the SDK;
all signal data via the F-01 `withListingSignals(client)` surface. Honest negative states only — never a fabricated
number for an unlinked or evidence-less listing (ADR-17).

## Command reference (HARNESS-PROFILE §3)

- WEB_TYPECHECK: `cd apps/web && npx tsc --noEmit`
- WEB_TEST: `cd apps/web && npx vitest run`
- WEB_TEST_ROUTER: `cd apps/web && npx vitest run src/app/AppRouter.test.tsx`
- QA_LIVE (L2, P7 only): docker compose dev-stack up, live-drive `http://localhost:5174/anuncios` (light + dark).

## Ground-truth seam map (base 8b6c4b30)

- `pages/AnunciosPage.tsx:78` main; fetch `useQuery({queryKey: listingsQueryKeys.page(installationId, pageOptions),
  queryFn: () => client.listListings(pageOptions)})`; `summaryQuery = anunciosSummaryQuery(client, installationId)`;
  `exceptionLabels: Record<ListingException,...>` (exhaustive — new enum values force a compile update).
- `pages/AnunciosTable.tsx` columns (`renderMargin`, `renderLinkState`); `UnknownValue`/`ConflictTag` from `@marketplace-central/ui`.
- `pages/anunciosQueryState.ts`: `exceptionValues = ["sync_error","stale","unlinked","below_margin"]`; `setFilterParam`
  writes `filter.<key>`; `parse/apply/clearFilters/toListingListOptions`.
- `pages/anunciosQueries.ts`: `anunciosSummaryQuery` (uses `listingsQueryKeys.summary`); `AnunciosClient = Pick<Client,...>`.
- `pages/ListingDetailPanel.tsx`: drawer via `DetailPanel` / `DetailBody` sections.
- `routes/anuncios.tsx`: thin `<AnunciosPage/>` wrapper.
- `packages/web-query`: `FreshnessIndicator({asOf})`, `formatAsOf`, `QUERY_STALE_TIME.listings`, `listingsQueryKeys` (web-query is hub/M-03-owned — CONSUME, do not edit).
- `app/ClientContext.tsx` (FORBIDDEN): `useClient()` returns the full `Client` — so `withListingSignals(useClient())`
  wraps it in-page without touching ClientContext.
- M-03 tokens: `apps/web/src/index.css` (CONSUME token classes; do NOT edit).

---

## SLICE F02-S1 — SDK wiring + query-state exception extension (foundation)

- **id:** F02-S1
- **complexity:** standard
- **goal:** Route the page through `withListingSignals(useClient())` so items are typed `ListingWithSignal`, and extend
  URL query-state to accept the 3 new exception values additively (no W1 filter behavior changed).
- **failing_test_first:** extend `pages/anunciosQueryState.test.ts` (or NEW if absent): `?exception=sem_vinculo`
  parses/round-trips; invalid `?exception=xyz` is ignored and cleared (no crash); the existing 4 values still parse
  identically (regression).
- **validation_kind:** L0 (WEB_TYPECHECK) + L1 (WEB_TEST).
- **commands:** WEB_TYPECHECK ; WEB_TEST
- **expected_artifacts:** query-state test transcript; tsc clean with `ListingWithSignal` flowing into the page.
- **write_set:**
  - `pages/anunciosQueryState.ts` — extend `exceptionValues` to include `sem_vinculo|abaixo_custo|sem_evidencia`
    (drives parse/validate/clear). No change to `setFilterParam` mechanism.
  - `pages/AnunciosPage.tsx` — wrap client: `const client = withListingSignals(useClient())`; queryFn now yields
    `ListingWithSignal[]`; `exceptionLabels` extended to the new enum values (PT copy). No fetch added outside SDK.
- **done_criteria:** all query-state tests green incl. invalid-ignored; existing 4-value behavior unchanged; tsc clean.
- **open_questions:** none.

---

## SLICE F02-S2 — VS MERCADO column (honest per-state rendering)

- **id:** F02-S2
- **complexity:** complex
- **goal:** Add the VS MERCADO column to the table: position + delta vs `price_to_win` (mono), with distinct honest
  rendering per `signal_status` — never a fabricated number.
- **failing_test_first:** NEW/extended `pages/AnunciosTable.test.tsx`: `OK` ⇒ position + delta shown (mono);
  `SEM_VINCULO` ⇒ grey chip linking to `/vinculos` (no number); `NO_PRICE_EVIDENCE` ⇒ `UnknownValue` "—" + tooltip;
  `STALE` ⇒ value + `FreshnessIndicator` (amber). Renders all 4 states in one table without throwing.
- **validation_kind:** L0 (WEB_TYPECHECK) + L1 (WEB_TEST).
- **commands:** WEB_TYPECHECK ; WEB_TEST
- **expected_artifacts:** table test transcript covering 4 states; no hardcoded hex (token classes only).
- **write_set:**
  - `pages/AnunciosTable.tsx` — new `renderMarketSignal(item)` column reading `item.market_signal`/`item.signal_status`;
    consume `FreshnessIndicator` (web-query) + `UnknownValue` (`@marketplace-central/ui`); `SEM_VINCULO` uses a router
    link to `/vinculos`.
- **done_criteria:** 4 states render distinctly; SEM_VINCULO/NO_PRICE_EVIDENCE never render a numeric market price;
  STALE shows age; W1 columns unchanged; L1 green.
- **open_questions:** none.

---

## SLICE F02-S3 — Exception chips (summary counters -> URL filter)

- **id:** F02-S3
- **complexity:** standard
- **goal:** Render exception chips with summary counts in the header; clicking a chip writes `?exception=<value>` to the
  URL (deep-linkable, F5-restorable) and filters the table.
- **failing_test_first:** extend `pages/AnunciosPage.test.tsx`: chip click ⇒ URL gains `?exception=sem_vinculo` and the
  list query re-issues with that filter; summary-query failure ⇒ chips hidden, table renders normally (isolated error,
  negative scenario from feature.md).
- **validation_kind:** L0 (WEB_TYPECHECK) + L1 (WEB_TEST).
- **commands:** WEB_TYPECHECK ; WEB_TEST
- **expected_artifacts:** page test transcript (click ⇒ URL + filtered query; summary-fail ⇒ chips hidden).
- **write_set:**
  - `pages/ListingsSummary.tsx` (or the summary header region of `AnunciosPage.tsx`) — render `exceptions.sem_vinculo`
    / `abaixo_custo` / `sem_evidencia` counters as clickable chips using `applyAnunciosQueryState` to set the filter.
  - `pages/anunciosQueries.ts` — extend `anunciosSummaryQuery` return typing to the extended
    `ListingSummaryExceptionsWithSignals` (via listings.ts). Query key `listingsQueryKeys.summary` reused (no new key).
- **done_criteria:** chip click deep-links + filters; F5 restores; summary error isolates (chips hidden, table ok); L1 green.
- **open_questions:** none.

---

## SLICE F02-S4 — Agrupar-por-produto toggle (by-product endpoint)

- **id:** F02-S4
- **complexity:** complex
- **goal:** Add a URL-driven "agrupar por produto" toggle that switches the list to `listListingsByProduct` (existing
  endpoint), rendering product-group rows while keeping per-listing signal columns inside each group.
- **failing_test_first:** extend `pages/AnunciosPage.test.tsx`: toggle on ⇒ URL flag set + `listListingsByProduct`
  query issued (key `listingsQueryKeys.byProduct`); groups render per CODPROD; a 1-listing group renders as a normal
  group (no special collapse — negative scenario). Toggle off restores flat list (W1) unchanged.
- **validation_kind:** L0 (WEB_TYPECHECK) + L1 (WEB_TEST).
- **commands:** WEB_TYPECHECK ; WEB_TEST
- **expected_artifacts:** page test transcript (toggle ⇒ by-product query + grouped render; off ⇒ flat W1).
- **write_set:**
  - `pages/anunciosQueryState.ts` — additive URL param for the group toggle (parse/apply/clear).
  - `pages/AnunciosPage.tsx` — conditional query: `withListingSignals(client).listListingsByProduct(...)` when grouped;
    group rendering reuses the VS MERCADO cell from S2.
  - `pages/AnunciosTable.tsx` — group-row rendering variant (per-listing signal columns preserved within group).
- **done_criteria:** toggle deep-links; grouped view shows CODPROD groups with per-listing signals; flat view identical
  to W1 when off; L1 green.
- **open_questions:** none.

---

## SLICE F02-S5 — Drawer evidence section

- **id:** F02-S5
- **complexity:** standard
- **goal:** Add an Evidência section to the listing drawer: source, fetched_at (via FreshnessIndicator), freshness,
  `match_status`, sample size `n_offers`/`n_sellers`, and a product link — small sample visible, never hidden (IC-03).
- **failing_test_first:** extend `pages/ListingDetailPanel.test.tsx`: OK item ⇒ evidence section shows source +
  FreshnessIndicator + n_offers/n_sellers + match_status + product link; NO_PRICE_EVIDENCE ⇒ named state, no fabricated
  value; SEM_VINCULO ⇒ no evidence section / link to /vinculos.
- **validation_kind:** L0 (WEB_TYPECHECK) + L1 (WEB_TEST).
- **commands:** WEB_TYPECHECK ; WEB_TEST
- **expected_artifacts:** drawer test transcript covering present/absent evidence.
- **write_set:**
  - `pages/ListingDetailPanel.tsx` — new `DetailBody` Evidência section reading `item.market_signal.evidence`;
    consume `FreshnessIndicator`; small `n_offers`/`n_sellers` shown verbatim (never suppressed).
- **done_criteria:** evidence fields render for OK/STALE; negative states show named copy not values; drawer W1 sections
  unchanged; L1 green.
- **open_questions:** none.

---

## SLICE F02-S6 — Retheme to M-03 tokens + router-test guard (APPTEST-M05)

- **id:** F02-S6
- **complexity:** standard
- **goal:** Replace hardcoded colors on the Anuncios* surface with M-03 token classes (`apps/web/src/index.css`),
  verify light+dark, and — only if the new markup invalidates the `/anuncios` router assertion — update that single
  case (APPTEST-M05).
- **failing_test_first:** (a) grep/lint guard: zero hardcoded hex / raw `slate-`/`blue-`/`red-`/`emerald-` literals in
  Anuncios* files (token classes only); (b) `AppRouter.test.tsx:119-123` "renders the anuncios workspace" (tab "Todos")
  stays green — extend/repair ONLY this case if S1-S5 markup changed the tab structure.
- **validation_kind:** L0 (WEB_TYPECHECK) + L1 (WEB_TEST + WEB_TEST_ROUTER).
- **commands:** WEB_TYPECHECK ; WEB_TEST ; WEB_TEST_ROUTER
- **expected_artifacts:** router test green; grep proof of no hardcoded color; (P7) light+dark screenshots.
- **write_set:**
  - `pages/AnunciosPage.tsx`, `pages/AnunciosTable.tsx`, `pages/ListingDetailPanel.tsx`, `pages/ListingsSummary.tsx`,
    `pages/ListingsRefreshControl.tsx` — swap hardcoded colors for M-03 token classes.
  - `app/AppRouter.test.tsx` — APPTEST-M05: the `/anuncios` case ONLY, iff invalidated.
- **done_criteria:** no hardcoded color; light+dark OK; router "Todos" assertion green; W1 deep-links C08/C09 pass; L1 green.
- **open_questions:** none.

## Regression guard (F-02, all slices)

Every slice is additive to the LIVE W1 page. W1 acceptance (C08/C09 deep-links, filters, ordering, pagination, drawer,
manual refresh) MUST remain green — treated as the standing regression suite per slice. `withListingSignals` re-uses the
same base client and query keys, so no duplicate fetch and no new invalidation seam.

---

## ⚠️ RESHAPE — HUB RULING D-56 (BINDING; supersedes F02-S2 column shape, F02-S5 drawer shape, F02-S6)

Context: HUB ratified `DESIGN-REFERENCE.md` + `design/handoff/Anuncios.dc.html` (main `8144238`) as the 1:1 visual-QA
gate. My ESCALATION-DR surfaced that the ratified Anuncios screen has NO standalone vs-mercado column and NO
market-evidence drawer panel. RULING D-56 = **Option 1: design wins, F-01 backend reused via bounded RESHAPE (no
revert; keep all commits `0f8d36f6`/`b4358923`).** The 1:1 gate = fidelity to the design's VISUAL LANGUAGE, honestly
powered by F-01's real signal logic. `ListingReadModel` already exposes every needed field (`sync_state`,`sync_error`,
`published_quantity`,`quality_score`,`pending_issue`,`price`,`link`,`market_signal`,`signal_status`) — pure FE work.

Ratified Anuncios table cols (9): `☑ · MLB(mono) · TÍTULO · PRODUTO · PREÇO(valor + %chip vs mercado) · EST · SYNC(status
pill) · QUAL(%) · PENDÊNCIA`. Decision (ACKed to hub, proceed-unless-countered): REMOVE table cols Modalidade / Vendas-30d /
Margem (extra vs prototype = layout finding); RELOCATE Margem + Modalidade into the drawer (prototype drawer shows both);
Vendas-30d dropped from table (data stays in model). Tokens: `apps/web/src/index.css` `@theme` already maps `--color-*` →
the exact ratified hex (light+dark). Non-regression (findings if violated): valor R$ + chip % colorido, NUNCA dupla sublinha;
sem edição inline; sem gradientes/emojis/Inter/Roboto.

### SLICE F02-R1 — table restructure + price-cell vs-mercado chip (supersedes S2 column)

- **id:** F02-R1  · **complexity:** complex  · **supersedes:** F02-S2 standalone "Vs. mercado" column (`0f8d36f6`)
- **goal:** Reshape `AnunciosTable` to the ratified 9-column set. DROP the standalone "Vs. mercado" `<th>`/cell. FOLD
  `renderMarketSignal` into the **PREÇO cell** as `valor R$` + a compact colored **%-chip vs mercado** (999px pill,
  token colors). Remove cols Modalidade / Vendas-30d / Margem. Add/confirm cols EST(`published_quantity`),
  SYNC(`sync_state` status pill), QUAL(`quality_score` %), PENDÊNCIA(`pending_issue`). MLB(mono) + TÍTULO + PRODUTO(`link`).
- **honest states (ADR-17) in the price chip:** `OK` → % delta chip verde (below/at) / vermelho (above) per sign;
  `STALE` → âmbar % chip + FreshnessIndicator age; `NO_PRICE_EVIDENCE` → NO chip, price only (honest absent, never 0/%);
  `SEM_VINCULO` → NO number, link "/vinculos"; missing `signal_status` → defensive absent. NUNCA dupla sublinha na célula.
  EST/QUAL/PENDÊNCIA honest when their field is null (UnknownValue "—", never fabricated 0).
- **failing_test_first:** extend `AnunciosTable.test.tsx`: (a) header = exactly the 9 prototype cols, NO "Vs. mercado"/
  "Modalidade"/"Vendas 30d"/"Margem" `<th>`; (b) price cell renders value + %chip for OK/STALE, NO chip/number for
  NO_PRICE_EVIDENCE + SEM_VINCULO (queryByText null on % and price-number resp.); (c) STALE shows freshness marker;
  (d) EST/QUAL null → "—" not "0". Red first.
- **validation_kind:** L0 (WEB_TYPECHECK stash-diff) + L1 (WEB_TEST). **commands:** WEB_TYPECHECK ; WEB_TEST
- **write_set:** `pages/AnunciosTable.tsx` (+ `.test.tsx`) ONLY. Update `TABLE_COLUMN_COUNT`. Token classes only (retheme
  finalized in R3, but do NOT introduce NEW raw literals here).
- **done_criteria:** 9-col header matches prototype; price-cell chip honest 4-state + no double-underline; EST/QUAL/PEND
  honest-null; W1 selection/rows/deep-link regression green; L1 green.
- **open_questions:** none (column-drop decision ruled above).

### SLICE F02-R2 — drawer fold into prototype card idiom (supersedes S5 panel shape)

- **id:** F02-R2  · **complexity:** complex  · **supersedes:** F02-S5 evidence-panel shape (`b4358923`) — KEEP the evidence data, re-skin only
- **goal:** Re-skin `ListingDetailPanel` to the ratified drawer language: sticky ~300px, header (MLB mono + short title + ✕),
  foto placeholder + status pill, sync-error box (when `sync_error`), a **2×2 bordered-card grid** [Preço · Est. publicado ·
  Margem est. · Qualidade], action row, **LINHA DO TEMPO** timeline, "Abrir edição completa →" link. Then render the F-01
  **market evidence** (ofertas/vendedores/posição/preço-p-ganhar) as ADDITIVE bordered cards in the SAME grid/token idiom
  (not a foreign panel). Honest states preserved exactly (SEM_VINCULO/NO_PRICE_EVIDENCE/STALE/OK, unknown≠zero).
- **failing_test_first:** extend `ListingDetailPanel.test.tsx`: 2×2 grid cards present (Preço/Est/Margem/Qualidade); evidence
  cards render for OK/STALE with real numbers; SEM_VINCULO → no market numbers + link; NO_PRICE_EVIDENCE → "—" not 0;
  timeline block present. Red first.
- **validation_kind:** L0 + L1. **commands:** WEB_TYPECHECK ; WEB_TEST
- **write_set:** `pages/ListingDetailPanel.tsx` (+ `.test.tsx`) ONLY. (Margem + Modalidade relocated here from the table per R1.)
- **done_criteria:** drawer matches prototype shape (cards/grid/timeline); evidence retained in the design idiom; honest
  states intact; no inline edit; L1 green.
- **open_questions:** none.

### SLICE F02-R3 — retheme to ratified tokens + APPTEST-M05 (was F02-S6)

- **id:** F02-R3  · **complexity:** standard  · **replaces:** F02-S6
- **goal:** Swap ALL raw tailwind color literals (`slate-`,`blue-`,`red-`,`emerald-`,`amber-`,`gray-`, any hex) across the
  Anuncios* surface for the ratified token classes wired in `index.css @theme` (`text-ink`,`text-muted`,`text-faint`,
  `bg-surface`,`bg-surface-2`,`border-border`,`text-accent`,`bg-accent-soft`,`text-accent-ink`,`text-warn`,`bg-warn-soft`,
  `text-amber`,`bg-amber-soft`,`text-info`,`bg-info-soft`). Fonts: body already Instrument Sans; monospace cells → IBM Plex
  Mono utility. Resolve light+dark (tokens do both). NO behavior change.
- **failing_test_first:** grep guard = ZERO raw color literals in the 5 Anuncios* write_set files; full WEB_TEST stays green
  (behavior identical); APPTEST-M05 `/anuncios` case asserts AnunciosPage mounts (extend ONLY that case iff R1 markup changed it).
- **validation_kind:** L0 + L1 (+ WEB_TEST_ROUTER). **commands:** WEB_TYPECHECK ; WEB_TEST ; WEB_TEST_ROUTER
- **write_set:** `pages/AnunciosPage.tsx`, `pages/AnunciosTable.tsx`, `pages/ListingDetailPanel.tsx`, `pages/ListingsSummary.tsx`,
  `pages/ListingsRefreshControl.tsx` (token swap) + `app/AppRouter.test.tsx` (`/anuncios` case ONLY, APPTEST-M05).
- **done_criteria:** zero raw color literals; light+dark OK; router assertion green; full suite green; ready for P5→P7 1:1 QA.
- **open_questions:** none.

**RESHAPE DAG:** R1 → R2 (R2 relocates Margem/Modalidade dropped by R1; both touch different files so may overlap, but R2's
drawer content depends on R1's column-drop decision) → R3 (retheme touches all, runs last, after R1+R2 committed).
Each slice: failing-test-first, independent review (implementer ≠ reviewer), fix→delta-re-review. Then S5-review is MOOT
(superseded by R2); P5 ladder → P6 dual gate → P7 QA 1:1 vs `Anuncios.dc.html`.
