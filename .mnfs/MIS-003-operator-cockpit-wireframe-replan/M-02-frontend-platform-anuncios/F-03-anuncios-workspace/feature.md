# F-03-anuncios-workspace

```yaml
id: F-03
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-02
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-003. Binding contracts: IC-05 (route, keys, states, glossary), IC-02 (endpoints, filter grammar, item shape), R-01 screen 2a inventory. ADR-15/17.

## Milestone

M-02 frontend-platform-anuncios. Depends on F-01 + F-02.

## Brief

Build `/anuncios` — Anúncios unificado (wireframe 2a), read-only. Tabs (Todos / Ativos / Pausados / Com pendência — wireframe 2a exact set), filter bar mapping 1:1 to IC-02 filter grammar + `q`, cursor-paginated table with wireframe columns (anúncio/título+MLB, modalidade tag, vínculo, preço, estoque, vendas 30d, margem, sync_state tag, ⚑ badges), summary counters strip from `getListingsSummary`, listing detail drawer from `getListing` (timeline included), manual refresh button driving `refreshListings` + polling operation run. Bulk-selection checkboxes render and accumulate selection state; action buttons render disabled with tooltip "disponível em breve" (M-03 wires them).

EARS:
- While the operator sets tab/filter/search, when state changes, the URL query shall encode it (`tab`, `filter.*`, `q`) and a reload shall restore the identical view.
- While a page of listings loads, when a fact is null, the cell shall render `UnknownValue` with the IC-05 hint (e.g. cost null → margem "—" hint "sem custo no ERP → não simulado"), never zero.
- While the operator clicks "Atualizar", when refresh is accepted (202), the UI shall show sync progress and refetch page + summary on run completion; when 409 `refresh_in_progress`, the UI shall attach to the running operation instead of erroring.
- While a row's sync_state is error, when rendered, the tag shall show "com erro" (red) per IC-02 label map.
- While selection spans pages, when the operator paginates, the selection count chip shall persist selected ids.

## Inputs

- IC-02 operations + canonical examples, IC-05 keys (`listingsQueryKeys.*`) + state components, R-01 screen 2a element inventory (tabs, columns, badges, counters), F-01 context/shell, F-02 components.

## Expected Output

- `/anuncios` page using exclusively `useQuery` + `listingsQueryKeys` (no direct fetch).
- Tab→filter mapping: Todos = no filter; Ativos = `filter.status=active`; Pausados = `filter.status=paused`; Com pendência = `filter.has_exception=true` (IC-02 grammar — no wildcards). Specific exception kinds (e.g. `filter.exception=below_margin`) are filter-bar chips, not tabs.
- Detail drawer with timeline + "▸ técnico" collapsible raw section.
- Component tests: URL round-trip (state→URL→state), null-render, 409-attach, selection persistence.

## Constraints

- Read-only: no mutation endpoints called; disabled buttons carry no handlers beyond tooltip.
- All copy pt-BR per glossary; provider raw strings only inside "▸ técnico".
- Columns/tabs from wireframe 2a only — no invented screens (mission rule).
- Uses only IC-05 seam files; does not modify router/context/web-query (one writer per seam — those are F-01/F-02 outputs).

## Negative Scenarios

- API 400 `invalid_filter` (stale URL param) → ErrorState with retry that clears offending filter.
- Empty result under filter → EmptyState "Nenhum registro encontrado." + hint to clear filters.
- Summary endpoint error while table ok → counters strip shows ErrorState inline; table unaffected (independent queries).

## Interaction Model

Query state lives in URL (single source); component state only for drawer-open + selection set. Selection keyed by composite listing_id, survives pagination, cleared on installation switch. Freshness: `FreshnessIndicator` from query `dataUpdatedAt`; staleTime listings 45s; manual refresh = `refreshListings` then invalidate `listings` namespace on run completion.

## Validation Expectations

- Vitest output: URL round-trip, null-render ("—" not "0"), 409-attach, selection tests green.
- Browser proof: `/anuncios?installation=X&tab=pendencia&filter.exception=sync_error` screenshot after F5 showing same tab+filter+rows.
- Network proof (devtools or test): zero non-TanStack fetches from the page.

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` created during feature execution.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (after F-02 accepted).
- Next action: compile context pack; read IC-02/IC-05 + R-01 §2a only.
- Required files/evidence: `validation.md` in this folder.
- Blockers or open decisions: none.
