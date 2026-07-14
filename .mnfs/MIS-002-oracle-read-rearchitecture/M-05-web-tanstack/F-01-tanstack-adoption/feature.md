# F-01-tanstack-adoption

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-05
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-2
lifecycle_scope: feature
```

## Mission

MIS-002. IC-01: queryKey namespaces, staleTime table, envelope shape.

## Milestone

M-05. Requires M-02 SDK method (`listCatalogProductFacts`).

## Brief

Adopt TanStack Query v5 in `apps/web`: QueryClientProvider at root with defaults (retry 1, refetchOnWindowFocus false); migrate Oracle-backed views (catalog list, stock/inventory view, profitability view) to hooks; catalog uses `useInfiniteQuery` over the IC-01 cursor envelope; catalog search uses a debounced (300ms) `useQuery` on `GET /catalog/products/search` with queryKey `['catalog','search',q]`; every migrated view shows an `as_of` freshness indicator ("dados de HH:mm:ss") + refresh button (refetch with no-cache header).

## Inputs

- `apps/web` React 19, hand-written fetch SDK at `packages/sdk-runtime/src/index.ts` (M-02 adds typed list method; verify shape before writing hooks).
- IC-01: queryKeys `['catalog', 'facts', {cursor}]` style namespaces, staleTime catalog 300_000 / stock 45_000 / pricecost 120_000 ms, envelope fields.
- Verify-at-install: `@tanstack/react-query` v5 (`npm view @tanstack/react-query version`; v5 API — `gcTime` not `cacheTime`, object-syntax `useQuery`).

## Expected Output

- Dependency added intentionally (one package.json change + lockfile).
- `apps/web/src/queries/` (or repo-conventional location) with hooks per data class; QueryClient wiring; as_of indicator component; refresh wired to `queryClient.refetchQueries` + no-cache request option in sdk call.
- Migrated views use hooks only; no direct fetch remains in them.
- One intentional commit.

## Constraints

- SDK stays the transport (hooks call sdk-runtime methods; no parallel fetch layer).
- staleTime values from IC-01, env/config-tunable at build level is NOT required (hardcoded constants module acceptable).
- Do not migrate non-Oracle views (scope creep).
- No mutation invalidation here (F-02).

## Interaction Model

- View mount within staleTime → render from cache instantly, no spinner, indicator shows cached as_of.
- Scroll end (catalog) → fetchNextPage with `next_cursor`; `next_cursor=null` → "fim da lista", no further calls.
- Refresh click → no-cache request → indicator updates; button disabled while inflight.
- Server 503 → error state with retry button (1 auto-retry only).

## Negative Scenarios

- While server returns 503 `source_unavailable`, when a view loads, the system shall show error state (not blank/zeroed data) with retry.
- While `next_cursor` is null, when user hits list end, the system shall issue no further requests.

## Validation Expectations

- Component tests with mocked sdk: infinite chain, staleTime remount (1 call), as_of render + refresh header assertion.
- `npm run build` green; TanStack version verification note in validation.md.

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` at execution time.

## Handoff

- Current status: briefed
- Next owner: Feature Implementer (mpc-implementer, gpt-5.6-luna high)
- Next action: spec → plan → implement → evidence
- Required files/evidence: `validation.md`
- Blockers or open decisions: None
