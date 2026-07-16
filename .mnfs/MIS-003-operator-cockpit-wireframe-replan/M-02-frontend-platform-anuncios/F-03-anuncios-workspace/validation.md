# F-03 anuncios-workspace — validation evidence

Feature closed 2026-07-16 by CHIP-M02. Plan source: `../plan-batch.md` §F-03 slices (Sol-medium,
ledger row 1). Spec: `feature.md` + IC-02 (`../../research/listings-read-interface-contract.md`)
+ IC-05.

## Slices

| Slice | Commit | Worker | Review (ledger row) | Result |
| --- | --- | --- | --- | --- |
| F03-S1 URL codec + query adapters (anunciosQueryState/anunciosQueries/AnunciosPage skeleton) | 00efa5b8 | Luna high (row 20) | row 21 ACCEPT-WITH-CONDITIONS — 2 importants fixed (invalid_filter-gated retry; search replace:true) | 10 scoped / 96 full |
| F03-S2 AnunciosTable + ListingsSummary (honest unknowns, IC-02 labels) | c95c4836 | Luna high (row 22) | row 23 ACCEPT (StatCard divergence deliberate, see rulings) | 11 scoped / 107 full |
| F03-S3 cursor pagination + cross-page selection + disabled bulk | a36ec327 | Luna high (row 24, teardown-killed post-delivery) | row 25 ACCEPT-WITH-CONDITIONS — 3 coverage importants fixed | 13 scoped / 111 full |
| F03-S4 ListingDetailPanel drawer + ▸ técnico + ui closeLabel additive | 8aec4aa5 | Luna high (row 26) | row 27 ACCEPT-WITH-CONDITIONS — 2 importants + fixture suggestion fixed | 30 scoped / 123 full |
| F03-S5 ListingsRefreshControl (202 / 409-attach / poll / one-shot invalidate) | 6693babc | Sol low, complex (row 28) | row 29 ACCEPT | 7 scoped / 130 full |

## EARS / criteria coverage (component level)

- URL-as-truth round-trip: anunciosQueryState.test.ts — tab/q/filter.* encode+parse identical;
  unknown keys and invalid enums dropped; clearFilters removes all `filter.*`, preserves
  q/tab/installation (adjudication #11). Tabs map per feature.md (todos absent / ativos→
  `status=active` / pausados→`status=paused` / pendencia→`has_exception=true`).
- Honest nulls (ADR-17): AnunciosTable.test.tsx + ListingDetailPanel.test.tsx — null price/
  stock/sales/margin/quality → "—" never 0; cost null → margin hint byte-exact
  "sem custo no ERP → não simulado"; summary nullable counters → UnknownValue.
- IC-02 label maps: fixed sync_state pt-BR map exhaustive via `satisfies` (error→"com erro" red);
  link conflict → ConflictTag "divergente".
- M02-C09 refresh/409: ListingsRefreshControl.test.tsx — 202 captures operation_run_id; typed
  409 `refresh_in_progress` attaches via `error.details.operation_run_id` with no error render;
  polling exclusively TanStack `refetchInterval` 2s over `syncQueryKeys.runs` +
  `listIntegrationOperationRuns` (adjudication #12, no custom timers), stops at terminal
  (proven by exhausted-mock call-count); succeeded invalidates `queryKeyNamespaces.listings`
  exactly once (rerender-proof); failed/cancelled → honest failure text, no invalidation.
- invalid_filter retry: AnunciosPage.test.tsx — retry clears all `filter.*` only when
  `error.code === "invalid_filter"`; any other error → plain refetch keeping filters.
- Selection persistence: AnunciosSelection.test.tsx — opaque composite ids accumulate across
  cursor pages, survive back-navigation, header checkbox page-scoped + indeterminate, cleared
  on installation switch; count chip `N selecionado(s)`.
- Cursor pagination: first-page query key byte-identical to pre-cursor shape (cursor
  conditionally spread); tab/filter/search/installation change resets cursor + stack
  (asserted via lastOptions); Anterior/Próxima disabled states.
- Detail drawer: ListingDetailPanel.test.tsx — getListing once under listingsQueryKeys.detail,
  no page refetch; timeline in API order (non-monotonic fixture proves no client sort);
  `message_pt` primary, `message_provider` only behind `<details><summary>▸ técnico`;
  close "Fechar painel" (ui DetailPanel additive `closeLabel` prop, default unchanged, 8/8).
- Mutation affordances: bulk Pausar/Atualizar preço/Re-sync and drawer Corrigir/Simular/Pausar
  disabled, `title="disponível em breve"`, zero SDK calls on click.
- Zero non-TanStack fetches: pages use only useQuery/useMutation over the SDK client (reviewer-
  verified per slice); no raw fetch in F-03 files.

## Rulings applied

- SummaryCounter diverges from shared StatCard deliberately (wireframe-2a; StatCard.value cannot
  take UnknownValue JSX) — row 23.
- ui DetailPanel `closeLabel` additive prop (chip seam) so pt-BR "Fechar painel" doesn't leak
  English into the workspace; default "Close panel" keeps existing consumers byte-identical.
- Refresh status copy chip-pinned pt-BR (na fila / em andamento / concluído / falhou /
  cancelado) — submitted for ratification at CLOSED.
- Browser proof (feature.md Validation Expectations screenshot) is hub-owned: chip never boots
  the dev stack (profile seam); deferred to milestone QA live-drive.

## Suite state at feature close

`cd apps/web; npm test -- --configLoader native` → 22 files, 130/130 green (exit 0, chip
re-verified after each slice). packages/ui DetailPanel 8/8. Plain `npm test` config-load
failure = pre-existing Node 26/esbuild issue; occasional npm exit-1 with all-green output =
teardown flake (both field findings in CLOSED payload).
