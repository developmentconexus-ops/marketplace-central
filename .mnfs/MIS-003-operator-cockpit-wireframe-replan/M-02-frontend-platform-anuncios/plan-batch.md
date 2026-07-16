## Plan summary

- **Execution order:** F-01 → F-02 → F-03, strictly sequential, one frontend writer.
- **Slice count:** F-01 = 3; F-02 = 6; F-03 = 5; total = 14.
- **Complex slices:** F-01/S1 and F-03/S5.
  - F-01/S1 requires bidirectional URL/context synchronization across four context states.
  - F-03/S5 is a 202/409 attach-and-poll state machine with terminal invalidation.
- **Standard slices:** all remaining slices use Luna high.
- **Scope:** only `apps/web`, `packages/web-query`, and `packages/ui`; `packages/sdk-runtime` is consumed read-only. No backend, migration, OpenAPI, SDK, or implementation-evidence edits.
- **Execution gate:** implementation must not begin until the blocking contract/scope findings below are adjudicated. In particular, M02-C03, the state-component feature order, dependency declarations, failure-copy literals, IC-02 post-M-01 drift, and `/orders` proxy/SPA collision are not safely decidable by an implementer.
- **W1 collision result:** CHIP-M02 remains the exclusive frontend writer. CHIP-M03 owns mutation backend/OpenAPI/SDK sections; CHIP-SAT owns dashboard/orders/sync/market backend and disjoint OpenAPI/SDK sections. No present source-file collision exists.
- **Additive contract-lock result:** none is required for the authorized frontend write-set. If the required internal dependency declarations are approved, `package-lock.json` is outside the named frontend seam and requires a hub-granted additive-only lock plus dependency-change authorization.
- **Current checkout:** `a49168e641ffd6f61932ca57c29b1d1bdcde2fb0`. M-01’s listing SDK operations are present at this SHA.

## F-01 slices

### 1. Installation context and URL selection state

- **Goal:** add `InstallationContext`/`useInstallation()` with the IC-05 shape `{installationId, setInstallationId, installations, status}`; one TanStack-managed installation request; known-ID selection; first-installation fallback plus URL rewrite; URL-only persistence; empty/error/retry states; selection changes without refetch.
- **Files touched:**
  - `apps/web/src/app/InstallationContext.tsx` — new.
  - `apps/web/src/app/InstallationContext.test.tsx` — new.
  - `apps/web/src/App.tsx` or `apps/web/src/app/AppRouter.tsx` — provider placement, one of these only after the router/provider ownership is fixed.
- **Failing test first:**
  - Render with `?installation=inst_ghost`; assert first fixture selected and URL rewritten.
  - Render with known ID; assert no rewrite.
  - Change selection; assert current pathname and unrelated query parameters remain and no reload occurs.
  - Navigate across three routes and switch once; assert exactly one `listIntegrationInstallations()` request.
  - Empty response produces `status="empty"`; rejection produces `status="error"` and retry starts one replacement request.
- **Done criteria:**
  - M02-C02 and M02-C03 pass for context-owned requests.
  - No `localStorage`, raw `fetch`, or `useEffect`-managed server fetch.
  - Selection changes never re-enter loading and never cause another installation fetch.
  - Query key and error-state ownership use the hub-ratified resolutions from Findings 1–3; the worker must not invent them.
- **Complexity:** **complex → Sol low**. Four-state async context plus two-way URL reconciliation can produce loops, silent URL/UI mismatch, or duplicate requests.
- **Estimated diff:** 220–290 lines.
- **Risks:**
  - React StrictMode remount behavior can expose duplicate request assertions unless the same QueryClient/key owns the request.
  - `setSearchParams` can erase `tab`, `filter.*`, and `q` unless it clones and updates the current parameters.
  - Provider placement outside `BrowserRouter` makes router hooks illegal; placement inside individual routes would recreate the fetch.
  - Blocked by the absent contracted installation query key and unavailable F-02 error/empty components.

### 2. Route map, query-preserving redirects, and future-workspace placeholders

- **Goal:** implement the IC-05 route table verbatim, all six query-preserving `<Navigate replace>` redirects, keep `/` on current Dashboard, retain `/classifications` and `/marketplaces` off-sidebar, and mount pt-BR placeholders for unfinished assigned routes. Create the `/anuncios` placeholder that F-03 replaces without another router edit.
- **Files touched:**
  - `apps/web/src/app/AppRouter.tsx`.
  - `apps/web/src/app/AppRouter.test.tsx`.
  - `apps/web/src/app/LegacyRedirect.tsx` — new, only if needed to keep redirect behavior isolated.
  - `apps/web/src/pages/WorkspacePlaceholder.tsx` — new.
  - `apps/web/src/pages/AnunciosPage.tsx` — new placeholder, later replaced by F-03.
- **Failing test first:**
  - Table-driven test for all six redirects, each retaining the complete search string and using replacement navigation.
  - Direct tests for `/anuncios`, `/catalogo/produtos/:productId`, `/protocolos/:protocolId`, `/classifications`, and `/marketplaces`.
  - Assert no route outside the IC-05 table was added.
  - Assert direct unfinished-route deep links keep the shell and display the contracted construction state.
- **Done criteria:**
  - M02-C01 passes at component level.
  - Redirect destinations and dynamic parameter names are byte-exact IC-05 values.
  - No route-path or redirect-map redesign.
  - F-03 can replace only `AnunciosPage.tsx`; it does not need to reopen `AppRouter.tsx`.
- **Complexity:** **standard → Luna high**.
- **Estimated diff:** 190–270 lines.
- **Risks:**
  - A string-only `<Navigate to="/destino">` drops the search string.
  - `/orders` cannot pass both browser deep-link redirect and current IC-05 proxy requirements until Finding 8 is resolved.
  - Whether M-04/M-05 paths mount legacy pages or placeholders is contradictory in the current milestone artifacts; the worker must not choose.
  - `EmptyState` does not exist until F-02, creating the feature-order defect in Finding 2.

### 3. Sidebar, top-bar pills, and installation-preserving navigation

- **Goal:** rebuild `Layout` with the eight-item IC-05 sidebar order, omit Mercado, keep configuration routes off-sidebar, add the functional `ML: <nome> ▾` selector and static empresa pill, and preserve `?installation=` while changing workspace routes.
- **Files touched:**
  - `apps/web/src/app/Layout.tsx`.
  - `apps/web/src/app/Layout.test.tsx` — new.
- **Failing test first:**
  - Assert the eight pt-BR labels and exact order.
  - Assert Mercado, Classifications, and Marketplaces are absent from main navigation.
  - With a ready installation context, select another account and assert the current pathname/unrelated query state remain.
  - Follow a sidebar link and assert `installation` is carried to the destination.
  - Empty context shows the contracted connection hint; ready context displays the installation-derived account name.
- **Done criteria:**
  - Layout and top bar satisfy IC-05 without changing route paths or adding product behavior.
  - M02-C02 navigation persistence is green.
  - All new shell copy is pt-BR.
- **Complexity:** **standard → Luna high**.
- **Estimated diff:** 200–280 lines.
- **Risks:**
  - Plain `NavLink to="/path"` drops the current query string.
  - Current icon labels and `DetailPanel` accessibility text contain English; only touched/new shell copy may be changed in this feature.
  - The exact installation field used as `<nome>` is not fixed by IC-05; contract owner should confirm `external_account_name` versus `display_name`.

## F-02 slices

### 1. Query stale-time registry, namespaces, and exact key builders

- **Goal:** extend the existing registry and namespaces without reshaping current keys; add IC-05 key builders verbatim.
- **Files touched:**
  - `packages/web-query/src/index.ts`.
  - `packages/web-query/src/index.test.ts` — new.
- **Failing test first:**
  - Snapshot `QUERY_STALE_TIME`.
  - Assert current `catalog`, `stock`, and `pricecost` values remain unchanged.
  - Assert the five new values and five new namespaces.
  - Assert exact output shapes for every listings, mutations, orders, and sync builder.
- **Done criteria:**
  - M02-C05 passes with `{catalog:300000, stock:45000, pricecost:120000, listings:45000, mutations:5000, orders:120000, sync:30000, market:300000}`.
  - Existing catalog/inventory/linkage/profitability keys are unchanged.
  - No installation-query key is invented here without contract ratification.
- **Complexity:** **standard → Luna high**.
- **Estimated diff:** 130–210 lines.
- **Risks:**
  - Renesting an existing key invalidates in-flight caches and violates IC-05 compatibility.
  - Object key stability must be deterministic; page filters should be normalized before reaching builders.

### 2. Exhaustive mutation invalidation crosswalk

- **Goal:** implement the single `invalidateAfterMutation(queryClient, type)` helper and exact namespace-prefix invalidations.
- **Files touched:**
  - `packages/web-query/src/invalidation.ts` — new.
  - `packages/web-query/src/invalidation.test.ts` — new.
  - `packages/web-query/src/index.ts` — export only.
- **Failing test first:**
  - One table row per IC-05 crosswalk case.
  - Assert both the exact namespace set and absence of extra invalidations.
  - Assert unknown input throws a typed error and performs zero invalidations.
  - Add the catalog-only product-enrichment case only after its discriminator is ratified.
- **Done criteria:**
  - M02-C04 passes.
  - Namespace-prefix invalidation uses `invalidateQueries({queryKey: namespace})`.
  - Envelope-write invalidation has no per-page duplicate implementation.
- **Complexity:** **standard → Luna high**.
- **Estimated diff:** 150–220 lines.
- **Risks:**
  - Parallel `Promise.all` invalidations can make call-order assertions brittle; tests should compare namespace sets unless ordering becomes contractual.
  - IC-05 does not give an exact product-enrichment discriminator; the worker must not invent one.

### 3. IC-03 failure-copy module

- **Goal:** provide one exhaustive failure-code-to-pt-BR mapping and the exact fallback `Falha desconhecida ({code})`.
- **Files touched:**
  - `packages/web-query/src/failureCopy.ts` — new.
  - `packages/web-query/src/failureCopy.test.ts` — new.
  - `packages/web-query/src/index.ts` — export only.
- **Failing test first:**
  - Iterate the 12 IC-03 codes:
    `provider_validation`, `provider_rate_limited`, `provider_unavailable`, `provider_auth`, `listing_paused_remote`, `link_unresolved`, `policy_missing`, `sku_invariant_violation`, `stale_source`, `conflict_remote_changed`, `type_not_enabled`, `internal`.
  - Assert each returns its ratified non-empty pt-BR literal.
  - Assert unknown input returns the byte-exact fallback with the original code interpolated.
- **Done criteria:**
  - M02-C07 passes.
  - No English fallback or mapping by provider message.
  - Copy literals come from a contract-owner amendment; implementation does not author product copy.
- **Complexity:** **standard → Luna high**.
- **Estimated diff:** 100–170 lines.
- **Risks:**
  - The 12 fixed strings are absent from IC-05/IC-03 today, so this slice is blocked.
  - Broad `string` typing can hide missing enum cases; the known-code mapping must be exhaustive.

### 4. Loading, error, and empty state components

- **Goal:** add the first half of the shared state vocabulary with fixed copy, retry behavior, detail/hint slots, and accessible semantics.
- **Files touched:**
  - `packages/ui/src/LoadingState.tsx` — new.
  - `packages/ui/src/ErrorState.tsx` — new.
  - `packages/ui/src/EmptyState.tsx` — new.
  - `packages/ui/src/LoadStates.test.tsx` — new.
  - `packages/ui/src/index.ts`.
- **Failing test first:**
  - `LoadingState` renders `Carregando…`.
  - `ErrorState` without a detail renders only `Erro ao carregar.` and button `Tentar novamente`; clicking invokes `onRetry` once.
  - `ErrorState` with detail appends it without replacing the fixed prefix.
  - `EmptyState` renders `Nenhum registro encontrado.` and an optional hint.
- **Done criteria:**
  - Relevant loading/error/empty rows of M02-C06 pass byte-exactly.
  - No English strings or speculative general state framework.
- **Complexity:** **standard → Luna high**.
- **Estimated diff:** 150–220 lines.
- **Risks:**
  - Error details must not accidentally render raw provider English outside the technical disclosure.
  - Creating these only in F-02 leaves F-01’s required EmptyState/ErrorState unsatisfied at F-01 acceptance.

### 5. Unknown, conflict, and freshness state components

- **Goal:** add `UnknownValue`, `ConflictTag`, and the contracted `packages/ui` `FreshnessIndicator`; reuse the existing `formatAsOf` behavior rather than duplicate timestamp formatting.
- **Files touched:**
  - `packages/ui/src/UnknownValue.tsx` — new.
  - `packages/ui/src/ConflictTag.tsx` — new.
  - `packages/ui/src/FreshnessIndicator.tsx` — new.
  - `packages/ui/src/FactStates.test.tsx` — new.
  - `packages/ui/src/index.ts`.
  - `packages/web-query/src/index.ts` — compatibility re-export or existing-export adjustment only after the package dependency direction is approved.
- **Failing test first:**
  - `UnknownValue` always renders `—`; a non-empty hint is accessible as a tooltip; absent hint creates no empty tooltip.
  - `ConflictTag` renders `divergente` with amber styling and optional detail.
  - `FreshnessIndicator({asOf})` renders `dados de HH:MM:SS` through `formatAsOf`.
  - Existing `web-query` consumers retain a compatible import or are covered by typecheck.
- **Done criteria:**
  - Remaining rows of M02-C06 pass.
  - Existing signature `{asOf: string | null | undefined}` is preserved.
  - ADR-17 rendering is centralized through `UnknownValue`.
- **Complexity:** **standard → Luna high**.
- **Estimated diff:** 170–250 lines.
- **Risks:**
  - `FreshnessIndicator` currently exists only in `web-query`; moving/re-exporting it without an approved internal dependency risks a package cycle or undeclared import.
  - Current aria label is English (`Data freshness`) and must become pt-BR on the rebuilt component without silently changing untouched legacy copy expectations.

### 6. CSS discovery, proxy contract, and canonical Vitest scope

- **Goal:** add the missing Tailwind source, required non-colliding proxy prefixes, resolve the `/orders` document/API collision as ratified, and make the invoked Vitest config discover app, `web-query`, and `ui` tests.
- **Files touched:**
  - `apps/web/src/index.css`.
  - `apps/web/vite.config.ts`.
  - `apps/web/src/app/viteProxy.test.ts`.
  - `apps/web/vitest.config.ts`.
- **Failing test first:**
  - Assert Tailwind scans `packages/feature-inventory/src`.
  - Assert the ratified IC-05 proxy rows.
  - Assert HTML navigation to `/orders` is not swallowed by the API proxy while API requests still reach the backend.
  - Add a discovery smoke proving tests under `packages/web-query/src` and `packages/ui/src` are included by the config used by `npm test`.
- **Done criteria:**
  - The proxy/CSS passing defects are fixed.
  - `npm test` runs all new F-01/F-02/F-03, `web-query`, and `ui` tests through the actual `apps/web/vitest.config.ts`.
  - No Playwright/MSW/new dependency is introduced.
- **Complexity:** **standard → Luna high**, unless the hub selects custom Vite middleware for `/orders`; that variant must be replanned as complex.
- **Estimated diff:** 100–180 lines.
- **Risks:**
  - Vite proxies every request whose path matches a configured prefix before SPA transformation; `/orders` therefore collides with the required legacy browser route.
  - `vite.config.ts` contains a broader test block, but the package script explicitly invokes `vitest.config.ts`; changing only the former is test theater.
  - `/dashboard` and `/sync` remain M-05-owned and must not be added by this slice.

## F-03 slices

### 1. Anúncios URL codec and TanStack query adapters

- **Goal:** make URL query state the single source for `tab`, `filter.*`, and `q`; translate it to IC-02 `ListingListOptions`; add TanStack query options/hooks using only the IC-05 builders and SDK client.
- **Files touched:**
  - `apps/web/src/pages/anunciosQueryState.ts` — new.
  - `apps/web/src/pages/anunciosQueryState.test.ts` — new.
  - `apps/web/src/pages/anunciosQueries.ts` — new.
  - `apps/web/src/pages/AnunciosPage.tsx` — replaces F-01 placeholder.
  - `apps/web/src/app/AppRouter.test.tsx` — update the `/anuncios` route expectation without changing router source.
- **Failing test first:**
  - Round-trip canonical deep link `installation=inst_1&tab=pendencia&filter.exception=sync_error`.
  - Assert tab mappings: Todos = no status filter; Ativos = `active`; Pausados = `paused`; Com pendência = `has_exception=true`.
  - Assert specific exception filters remain chips rather than tabs.
  - Assert installation changes clear component-only selection while preserving valid route query state.
  - Assert listing and summary query functions call SDK methods, not `fetch`.
- **Done criteria:**
  - URL state restores identically after remount/reload.
  - M02-C08’s state portion and M02-C10’s server-state architecture are satisfied.
  - Unknown filter keys are not silently sent to the API.
- **Complexity:** **standard → Luna high**.
- **Estimated diff:** 220–290 lines.
- **Risks:**
  - Encoding booleans or enums differently from IC-02 causes server `invalid_filter`.
  - A component-state mirror of URL state will drift on back/forward navigation.
  - The response behavior for server-side `invalid_filter` lacks the offending-key detail needed by the brief.

### 2. Summary strip and listing table with honest unknowns

- **Goal:** render the exception counters and the wireframe-derived listing columns using independent summary/page queries and all shared state components.
- **Files touched:**
  - `apps/web/src/pages/AnunciosPage.tsx`.
  - `apps/web/src/pages/AnunciosTable.tsx` — new.
  - `apps/web/src/pages/AnunciosTable.test.tsx` — new.
  - `apps/web/src/pages/ListingsSummary.tsx` — new.
  - `apps/web/src/pages/ListingsSummary.test.tsx` — new.
- **Failing test first:**
  - Table renders title/MLB, modality, link state, price, published stock, sales 30d, margin state, sync label, and system-master badges.
  - Null price/stock/sales/quality/cost-driven margin render `UnknownValue`, never zero.
  - Null cost gives margin hint `sem custo no ERP → não simulado`.
  - `sync_state="error"` renders `com erro`.
  - Summary failure renders inline ErrorState while a successful table remains visible, and vice versa.
  - Filtered empty response renders fixed EmptyState plus a clear-filter hint.
- **Done criteria:**
  - M02-C09 unknown rendering passes.
  - M02-C08 visible tab/chip/row state is enabled.
  - M02-C11 is enabled by parallel page/summary requests after installation resolution; final proof remains browser timing.
- **Complexity:** **standard → Luna high**.
- **Estimated diff:** 240–300 lines.
- **Risks:**
  - Current SDK uses `below_margin_worst_case` and nullable summary counters, not the older IC-02 `below_margin` fields.
  - The wireframe says grouped-by-product while the feature brief emphasizes `listListings`; no worker may choose flat versus grouped until Finding 9 is resolved.
  - Formatting unknown decimal facts before null checks can turn absence into `R$ 0,00`.

### 3. Cursor pagination and cross-page selection

- **Goal:** add cursor navigation and a selection set keyed by opaque composite `listing_id`; preserve selected IDs across pages and expose disabled bulk actions with the exact tooltip.
- **Files touched:**
  - `apps/web/src/pages/AnunciosPage.tsx`.
  - `apps/web/src/pages/AnunciosTable.tsx`.
  - `apps/web/src/pages/AnunciosSelection.test.tsx` — new.
- **Failing test first:**
  - Select rows on page one, advance by `next_cursor`, select another row, and assert accumulated count.
  - Navigate back without losing selections.
  - Change installation and assert selection clears.
  - Render Pausar, Atualizar preço, and Re-sync controls disabled; assert tooltip `disponível em breve`.
  - Assert disabled controls have no click handlers and make no SDK calls.
- **Done criteria:**
  - Cross-page selection behavior passes.
  - Mutation buttons meet the read-only constraint and contain no dead handlers.
  - Composite IDs remain opaque strings and are never parsed.
- **Complexity:** **standard → Luna high**.
- **Estimated diff:** 170–240 lines.
- **Risks:**
  - Selecting by row index or provider listing ID causes collisions across installations/variations.
  - Cursor state should not be mistaken for mutation selection grammar; no IC-03 request is built in M-02.
  - Native disabled buttons do not emit hover events consistently; tooltip implementation must remain accessible without adding a dependency.

### 4. Listing detail drawer and translated/technical separation

- **Goal:** load `getListing` on row activation and render timeline, translated error, technical raw detail, margin state, freshness, and disabled future actions.
- **Files touched:**
  - `apps/web/src/pages/ListingDetailPanel.tsx` — new.
  - `apps/web/src/pages/ListingDetailPanel.test.tsx` — new.
  - `apps/web/src/pages/AnunciosPage.tsx`.
- **Failing test first:**
  - Opening a row invokes `getListing(listing_id)` under `listingsQueryKeys.detail`.
  - Timeline remains newest-first as returned; no client re-sorting or provider-state translation.
  - `message_pt` is primary and `message_provider` appears only under `▸ técnico`.
  - Null facts use UnknownValue.
  - Corrigir/Simular/Pausar controls are disabled with `disponível em breve` and no handlers.
- **Done criteria:**
  - Detail drawer matches the read-only portion of wireframe 2a.
  - All new visible copy is pt-BR; provider English remains behind the technical disclosure.
- **Complexity:** **standard → Luna high**.
- **Estimated diff:** 190–270 lines.
- **Risks:**
  - Existing shared `DetailPanel` has English `aria-label="Close panel"`; using it unchanged would leak English into a rebuilt surface.
  - The current SDK contains worst-case margin facts, not a direct margin percentage; the drawer must not invent one.

### 5. Refresh, 409 attachment, polling, and terminal invalidation

- **Goal:** drive `refreshListings`; display queued/running progress; attach to the active operation on `409 refresh_in_progress`; poll the available operation-run read until terminal; invalidate the listings namespace exactly once on completion.
- **Files touched:**
  - `apps/web/src/pages/ListingsRefreshControl.tsx` — new.
  - `apps/web/src/pages/ListingsRefreshControl.test.tsx` — new.
  - `apps/web/src/pages/AnunciosPage.tsx`.
  - `apps/web/src/pages/anunciosQueries.ts`.
- **Failing test first:**
  - 202 response captures `operation_run_id`, starts progress, and polls.
  - Typed 409 extracts `error.details.operation_run_id`, attaches to that run, and shows no error toast.
  - Queued/running continue polling; succeeded/failed/cancelled stop it.
  - Success invalidates `queryKeyNamespaces.listings`, causing page and summary refetch.
  - Terminal failure renders honest translated failure state and does not fabricate successful freshness.
  - Changing installation cancels/abandons the old observed run.
- **Done criteria:**
  - M02-C09’s 409 behavior passes.
  - No raw fetch or custom timer/effect polling; use TanStack `refetchInterval`.
  - No mutation endpoint is called.
  - Browser validation can observe progress and updated page/summary state.
- **Complexity:** **complex → Sol low**. This is a multi-entry async state machine: new 202 run, attached 409 run, polling, terminal stop, installation change, and one-time invalidation.
- **Estimated diff:** 230–300 lines.
- **Risks:**
  - SDK has no direct operation-run-by-ID method; polling must use the available installation operation list and find the target ID unless the contract owner provides another read.
  - IC-05 has no explicit listing-refresh operation query key; misuse of an invented key or unrelated namespace is forbidden.
  - Repeated terminal renders can repeatedly invalidate unless terminal transition is guarded by query state.

## Write-sets

### F-01 write-DAG

| Slice | Exact write-set | Requires | Validation criteria |
| --- | --- | --- | --- |
| F01-S1 | `apps/web/src/app/InstallationContext.tsx`; `apps/web/src/app/InstallationContext.test.tsx`; exactly one provider-placement edit in `apps/web/src/App.tsx` or `apps/web/src/app/AppRouter.tsx` | Ratified installation query key; ratified F-01/F-02 state-component sequencing; direct UI dependency authorization | M02-C02, M02-C03, contributes M02-C10 |
| F01-S2 | `apps/web/src/app/AppRouter.tsx`; `apps/web/src/app/AppRouter.test.tsx`; optional new `apps/web/src/app/LegacyRedirect.tsx`; new `apps/web/src/pages/WorkspacePlaceholder.tsx`; new placeholder `apps/web/src/pages/AnunciosPage.tsx` | F01-S1 reviewed-green; route-placeholder conflict and `/orders` collision resolved | M02-C01; route prerequisite for M02-C08 |
| F01-S3 | `apps/web/src/app/Layout.tsx`; new `apps/web/src/app/Layout.test.tsx` | F01-S2 reviewed-green | M02-C02; sidebar portion of M02-C08 |

Ordering is strict S1 → S2 → S3 because provider placement, routes, and Layout share the router/context seam. F-01 is not acceptable while its required EmptyState/ErrorState behavior depends on unimplemented F-02 outputs.

### F-02 write-DAG

| Slice | Exact write-set | Requires | Validation criteria |
| --- | --- | --- | --- |
| F02-S1 | `packages/web-query/src/index.ts`; new `packages/web-query/src/index.test.ts` | F-01 reviewed-green | M02-C05 |
| F02-S2 | new `packages/web-query/src/invalidation.ts`; new `packages/web-query/src/invalidation.test.ts`; export in `packages/web-query/src/index.ts` | F02-S1 reviewed-green; enrichment discriminator ratified | M02-C04 |
| F02-S3 | new `packages/web-query/src/failureCopy.ts`; new `packages/web-query/src/failureCopy.test.ts`; export in `packages/web-query/src/index.ts` | F02-S2 reviewed-green; 12 copy literals ratified | M02-C07 |
| F02-S4 | new `packages/ui/src/{LoadingState,ErrorState,EmptyState}.tsx`; new `packages/ui/src/LoadStates.test.tsx`; `packages/ui/src/index.ts` | F02-S3 reviewed-green | M02-C06 |
| F02-S5 | new `packages/ui/src/{UnknownValue,ConflictTag,FreshnessIndicator}.tsx`; new `packages/ui/src/FactStates.test.tsx`; `packages/ui/src/index.ts`; compatibility-only edit to `packages/web-query/src/index.ts` | F02-S4 reviewed-green; dependency direction authorized | M02-C06 |
| F02-S6 | `apps/web/src/index.css`; `apps/web/vite.config.ts`; `apps/web/src/app/viteProxy.test.ts`; `apps/web/vitest.config.ts` | F02-S5 reviewed-green; `/orders` routing decision ratified | Enables all F-02 tests and M02-C10 |

S1–S6 are sequential because S2/S3 share the web-query export surface, S4/S5 share the UI export surface, and S6 makes those package tests visible to the invoked test command.

### F-03 write-DAG

| Slice | Exact write-set | Requires | Validation criteria |
| --- | --- | --- | --- |
| F03-S1 | new `apps/web/src/pages/anunciosQueryState.ts`; new test beside it; new `apps/web/src/pages/anunciosQueries.ts`; replace `apps/web/src/pages/AnunciosPage.tsx`; update `apps/web/src/app/AppRouter.test.tsx` only | F-02 reviewed-green; query-key and invalid-filter findings resolved | M02-C08 state restoration; M02-C10 |
| F03-S2 | `apps/web/src/pages/AnunciosPage.tsx`; new `AnunciosTable.tsx`/test; new `ListingsSummary.tsx`/test | F03-S1 reviewed-green; IC-02 current-shape and grouping rulings | M02-C08, M02-C09, enables M02-C11 |
| F03-S3 | `apps/web/src/pages/AnunciosPage.tsx`; `apps/web/src/pages/AnunciosTable.tsx`; new `AnunciosSelection.test.tsx` | F03-S2 reviewed-green | M02-C08 selection/deep-link workspace behavior |
| F03-S4 | new `apps/web/src/pages/ListingDetailPanel.tsx`/test; `apps/web/src/pages/AnunciosPage.tsx` | F03-S3 reviewed-green | M02-C09 unknown semantics |
| F03-S5 | new `apps/web/src/pages/ListingsRefreshControl.tsx`/test; `apps/web/src/pages/AnunciosPage.tsx`; `apps/web/src/pages/anunciosQueries.ts` | F03-S4 reviewed-green; refresh-operation query key ratified | M02-C09; contributes M02-C10 and M02-C11 |

After F03-S5 is reviewed-green and committed, emit the milestone’s required `COMMITTED` event identifying F-03 so the hub can rebase/unblock CHIP-M03 F-04. No M-03 file is touched by M-02.

## Contract satisfiability

### IC-05 versus current frontend

| IC-05 element | Current state at `a49168e…` | Satisfiability |
| --- | --- | --- |
| Route map | Only `/`, `/products`, `/classifications`, `/marketplaces`, `/integrations`, `/product-links`, `/inventory/stock-seguro`, `/orders`, `/simulator` exist | Additive F-01 work required |
| Redirects | No `<Navigate>` redirects | Additive F-01 work required |
| Sidebar | Nine English/legacy items; wrong order | Replacement F-01 work required |
| InstallationContext | Absent | Additive F-01 work required |
| Single installation fetch | Inventory, Integrations, Orders, Product Links, and Marketplaces independently call `listIntegrationInstallations` | Not satisfiable under current M-02 scope if these pages mount at new routes |
| `QUERY_STALE_TIME` | Exactly `{catalog:300000, stock:45000, pricecost:120000}` | Existing values match and can be extended |
| Existing namespaces | `catalog`, `inventory`, `linkage`, `profitability` | Must be preserved; five IC-05 namespaces are absent |
| Key builders | Only catalog/inventory/linkage/profitability builders exist | All IC-05 builders are absent |
| Installation query key | No contracted namespace/builder exists | Blocking omission |
| `invalidateAfterMutation` | Absent; legacy pages invalidate inline | Additive F-02 work, but enrichment discriminator is missing |
| `failureCopy` | Absent | Additive F-02 work, but fixed strings are missing |
| Loading/Error/Empty/Unknown/Conflict | Absent from `packages/ui` | Additive F-02 work |
| `FreshnessIndicator` | Exists in `packages/web-query/src/index.ts`, signature `{asOf: string \| null \| undefined}`, calls `formatAsOf`; not in `packages/ui`; aria label is English | Contracted package/signature only partially satisfied |
| `formatAsOf` | Exists; formats `dados de HH:MM:SS`; unknown/invalid gives `dados de desconhecido` | Reusable, must not be duplicated |
| Tailwind source | `feature-inventory` absent | F-02 fix required |
| Proxy rows | `/listings`, `/mutations`, `/market`, `/orders`, `/profitability` absent | F-02 required; `/orders` collides with SPA redirect |
| Vitest | `apps/web` script invokes `vitest.config.ts`, which includes app tests plus only `CatalogPage.test.tsx`; broader `vite.config.ts` test block is not used | New web-query/UI tests would be skipped unless F02-S6 fixes discovery |

### SDK/runtime operations at this SHA

`packages/sdk-runtime` is hand-written at this SHA, despite some planning text describing generated SDK output.

| Required operation | Available method and return type |
| --- | --- |
| Listings page | `listListings(options: ListingListOptions): Promise<ListingPage>` |
| Grouped listings | `listListingsByProduct(options: ListingListOptions): Promise<ListingGroupPage>` |
| Listing detail | `getListing(id: string): Promise<ListingDetail>` |
| Summary | `getListingsSummary(installationId: string): Promise<ListingSummary>` |
| Refresh | `refreshListings(req: RefreshListingsRequest): Promise<RefreshListingsAccepted>` |
| Installations | `listIntegrationInstallations(): Promise<ListResponse<IntegrationInstallation>>` |
| Pollable operation list | `listIntegrationOperationRuns(installationId: string): Promise<ListResponse<IntegrationOperationRun>>` |

Also exported:

- `ListingsRefreshConflictError` with `status: 409`, code `refresh_in_progress`, and `details.operation_run_id`.
- `IntegrationOperationRun.status` = `queued | running | succeeded | failed | cancelled`.
- `ListingListOptions` contains the IC-02 filter fields, `q`, `limit`, and `cursor`.

### IC-02 drift after M-01

Current OpenAPI and SDK are mutually aligned but differ from the research IC-02 document:

- SDK `ListingReadModel` has `below_margin_worst_case` and `icms_worst_case_by_uf`; IC-02 still says `below_margin`.
- SDK summary has nullable `below_margin_worst_case` and `margin_unknown`; IC-02 still says `exceptions.below_margin`.
- SDK `ListingStatus` includes post-contract provider statuses beyond IC-02’s four-value enum.
- The OpenAPI requires the newer worst-case fields, so the repository truth order places OpenAPI/SDK above stale `.mnfs` research.
- F-03 cannot safely implement the old names while consuming the current SDK. The contract owner must ratify the post-M-01 shape in IC-02 or explicitly direct M-02 to treat current OpenAPI/SDK as the amended read contract.

### Sibling W1 claims and locks

- CHIP-M03: mutation/protocol backend, OpenAPI mutation sections, SDK mutation rows. Its F-04 frontend work is explicitly gated on M-02 F-03 merge/rebase.
- CHIP-SAT: dashboard/orders/sync/market backend plus disjoint OpenAPI/SDK sections.
- Neither sibling currently claims `apps/web`, `packages/web-query`, or `packages/ui`.
- M-02 must not edit OpenAPI, SDK, backend, migrations, composition root, `/dashboard` proxy, or `/sync` proxy.
- **Additive contract-locks required by the authorized plan: none.**
- **Conditional lock request:** if hub approves adding direct internal dependencies in `apps/web/package.json` and/or `packages/ui/package.json`, the consequent root `package-lock.json` change needs an additive-only lock because it lies outside the named FE seam.

## Findings

### Blocking planning defects

1. **Installation query-key omission.** IC-05 requires all new server state through its builders but defines no installation-list namespace/key. F-01 cannot invent this key.

2. **F-01/F-02 feature-order contradiction.** F-01 requires `EmptyState` and `ErrorState`, while F-02 owns their creation and cannot start until F-01 is accepted. Temporary local duplicates would violate the anti-slop/dead-code rules. Resolve by moving the minimum shared states to F-01, prelanding them before F-01, or changing the acceptance order explicitly.

3. **M02-C03 is incompatible with current scope and legacy mounting.** Five legacy feature packages independently fetch installations. M-02 is restricted to `apps/web`, `packages/web-query`, and `packages/ui`, so it cannot make those pages consume `useInstallation()`. This also conflicts with the F-01 statement that legacy pages render unmodified at new paths. Hub must choose:
   - placeholders at M-04/M-05 paths until rebuild, or
   - expand M-02’s write-set to the affected feature packages and replan the diff.

4. **Placeholder versus legacy-page contradiction.** The milestone says future M-04/M-05 routes are construction stubs; F-01 says legacy pages remain reachable at new paths. The latter also exposes existing English UI on rebuilt routes, conflicting with the milestone’s pt-BR blocking rule.

5. **Dependency authorization is required.**
   - `apps/web` does not declare `@marketplace-central/ui`, although F-01/F-03 must consume shared UI states.
   - `packages/ui` does not declare `@marketplace-central/web-query`, although the contracted UI `FreshnessIndicator` must use existing `formatAsOf`.
   - Undeclared workspace imports or relative cross-package imports are not acceptable substitutes.
   - Any manifest/lock change requires a hub `REQUEST`; no slice assumes approval.

6. **Failure-copy literals are absent.** IC-05 says 12 fixed pt-BR strings exist but supplies only the fallback literal. IC-03 supplies codes and meanings, not byte-exact UI strings. Contract owner must add the 12 mappings.

7. **Product-enrichment discriminator is absent.** The invalidation crosswalk names “product enrichment edit” but does not define the exact helper input literal. M02-C04 writes `product-enrichment` descriptively. A typed exhaustive helper cannot choose among `product-enrichment`, `product_enrichment`, or a separate API without contract clarification.

8. **`/orders` proxy conflicts with `/orders` SPA redirect.** Vite 7 proxies every request whose path matches the prefix before SPA transformation. Current tests deliberately assert `/orders` is absent to avoid route shadowing. IC-05 requires both the proxy and direct legacy browser redirect. The hub must ratify either:
   - distinct API/client paths,
   - a content-negotiating custom Vite middleware design, or
   - a different redirect validation/runtime arrangement.
   A simple proxy row cannot satisfy M02-C01.

9. **IC-02 is stale against accepted M-01 OpenAPI/SDK.** The worst-case margin field and summary shapes must be reconciled before F03-S2.

10. **Flat versus grouped 2a table is unresolved.** R-01 calls 2a grouped by product; F-03 centers `listListings`, while `listListingsByProduct` is available. The worker must not select a presentation/data operation without a brief amendment.

11. **`invalid_filter` lacks offending-key details.** F-03 requires retry to clear the offending filter, but IC-02 specifies only code/message. Either define error details, permit clearing every submitted `filter.*`, or constrain the behavior to client-detected invalid enum values.

12. **Refresh polling key is unspecified.** The SDK supplies `listIntegrationOperationRuns`, but IC-05 gives no listing-refresh operation-run key. Contract owner must explicitly allow `syncQueryKeys.runs(installationId, filters)` for this purpose or add a listings refresh-run key.

### Non-blocking implementation risks

- First data paint has an installation-fetch dependency before listing and summary requests can begin; once selected, listing and summary must start in parallel.
- `DetailPanel` currently contains English accessibility copy and cannot be reused unchanged on a rebuilt pt-BR surface.
- Summary margin counters are nullable in the current SDK; null must render UnknownValue, not zero.
- No direct operation-run detail method exists; polling the installation’s operation list may be less efficient but is available.
- M-03 F-04 will later edit the Anúncios workspace. M-02 should keep bulk-action mounting localized and emit the F-03 `COMMITTED` event promptly, but it must not invent M-03’s modal interface.

### Test-infrastructure facts

- Vitest 3, jsdom, Testing Library, global test APIs, and jest-dom are already installed.
- No MSW, Playwright, or browser automation dependency exists; none may be added.
- Existing tests mock client methods directly and wrap query consumers in `QueryClientProvider`.
- Query tests use `createWebQueryClient()` or a retry-disabled `QueryClient`.
- Existing invalidation tests spy on `queryClient.invalidateQueries`.
- `apps/web/package.json` runs `vitest --config vitest.config.ts`; this is the canonical invoked config today.
- `apps/web/vitest.config.ts` does not discover `packages/web-query` or `packages/ui`; F02-S6 must fix that before those tests count as evidence.
- `vite.config.ts` has a broader, divergent test block; workers must not assume it is used by `npm test`.
- Browser evidence remains a manual/fresh-persona walkthrough with screenshots. `harness:browser` is intentionally unconfigured and must not be represented as passing automation.
- M02-C11 requires a real dev-local browser timing trace with the 2k fixture; a mocked component test cannot satisfy it.

### Worker false-alarm signatures and ladder facts

- `HPG_MIGRATION_FAILED` with `migrations_first=-1` means the hermetic Go module cache was not warmed; it is not a frontend or SQL defect.
- First-create `3D000 database does not exist` may be the PostgreSQL initialization race absorbed by the integration-lane retry.
- `TestPhase1SmokeFlow` is the profile’s known pre-existing L1 allowlisted failure; cite it rather than re-proving non-linkage.
- Governance must run in a clean detached worktree using the milestone’s accepted full 40-hex BaseSha. Short SHA produces `GOV_SEMANTIC_DRIFT id=base-sha-invalid`; scanning a main checkout containing worktrees can false-fail.
- Windows Go commands require absolute `GOCACHE`; this frontend milestone still runs the full platform ladder because it changes platform routing/configuration.
- Chips must request the hub-owned Docker dev stack; they must not start bare servers, bind `:8080`/`:5174`, or load `.env` contents.
