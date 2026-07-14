# M-05-web-tanstack

```yaml
id: M-05
type: milestone
status: planned
owner: Mission Strategist
parent: MIS-002
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-2
lifecycle_scope: milestone
```

## Mission

MIS-002 (`../mission.md`), IC-01 queryKey namespaces + staleTime table.

## Outcome

L1 frontend cache live: `apps/web` adopts TanStack Query v5 over the sdk-runtime client; Oracle-backed views use `useQuery`/`useInfiniteQuery` with staleTime mirroring server TTLs; `as_of` rendered as "dados de HH:mm:ss" with manual refresh (refetch + no-cache); mutations invalidate the right queryKey namespaces.

## Why This Milestone Exists

Without L1, every route change refetches; operators see spinners on data fetched seconds ago (journeys J1/J5). TanStack was an explicit user request. Mutation invalidation prevents the stale-after-write confusion (J7).

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | tanstack-adoption | add @tanstack/react-query v5 (verify-at-install), QueryClient defaults, query hooks for catalog page (useInfiniteQuery over IC-01 cursor), stock/profitability views; as_of freshness indicator + refresh button |
| F-02 | mutation-invalidation | wrap existing mutations (linkage confirm, product edits) with invalidateQueries on IC-01 namespaces; linkage candidate fetches use staleTime 0 + gcTime 0 |

## Dependencies

M-02 (envelope + SDK method) hard requirement; M-04 runs in parallel (staleTime/TTL coupling is contract-fixed in IC-01, not runtime-coupled).

F-02 does NOT wait for F-01 completion: queryKey namespaces and the invalidation crosswalk are contract-fixed (IC-01). Handshake = F-01's first step lands the QueryClientProvider wiring + hooks module skeleton; F-02 then builds mutation invalidation against the fixed namespaces in parallel.

## Feature Parallelization

| Lane | Content | Starts |
| --- | --- | --- |
| A | F-01: QueryClient wiring (step 1), then view migration + as_of/refresh | immediately |
| B | F-02: mutation invalidation wrappers + linkage staleTime 0/gcTime 0 | after lane A step 1 |

Seam ownership: F-01 owns provider/hooks/views; F-02 owns mutation call sites only. Crosswalk labels come verbatim from IC-01 — inventing a namespace is a reject.

## Risks

R2 (cursor misuse — useInfiniteQuery `getNextPageParam` reads `next_cursor` verbatim), R7 (visible staleness — as_of indicator + refresh), verify-at-install: TanStack v5 API surface.

## Done Means

`npm run build` green in `apps/web`; catalog view scrolls through cursor pages without refetching page 1; route-change back within staleTime → no network call (devtools/network evidence); refresh button issues no-cache request and updates as_of; linkage mutation invalidates catalog+linkage namespaces (refetch observed); no remaining direct fetch calls in migrated views.

## Handoff

- Current status: planned
- Next owner: Milestone Orchestrator
- Next action: dispatch F-01 (mpc-implementer, gpt-5.6-luna high)
- Required files/evidence: feature `validation.md` files; `validation-result.md`
- Blockers or open decisions: None

## Correction Handoff

Not applicable during initial planning.
