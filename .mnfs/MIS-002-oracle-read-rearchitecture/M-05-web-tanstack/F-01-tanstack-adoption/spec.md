# F-01-tanstack-adoption

```yaml
id: F-01
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: M-05
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-2
lifecycle_scope: feature
```

## Feature ID

F-01-tanstack-adoption

## Problem

The Oracle-backed catalog, stock-risk, and profitability views use manual effect-driven loading, so route changes and remounts lose the L1 cache and operators cannot see or force-refresh the source freshness represented by the response envelope.

## Requirements

- Add TanStack Query v5 with one root QueryClientProvider configured with `retry: 1` and `refetchOnWindowFocus: false`.
- Replace the legacy products page with a catalog facts page using the IC-01 envelope, `useInfiniteQuery`, debounced search, reserved queryKey namespaces, contract stale times, and null/quality-preserving display.
- Migrate stock-risk and profitability reads to `useQuery` without changing the F-02 mutation seams.
- Render every migrated view's `as_of` as local `dados de HH:mm:ss` and provide a disabled-while-inflight refresh that issues `Cache-Control: no-cache` through the apps/web-owned fetch wrapper.
- Keep the SDK as the only transport and do not edit `packages/sdk-runtime/src/**`, Go, or OpenAPI files.

## Acceptance evidence

Component tests with mocked SDK clients prove cursor append/null-stop, staleTime remount behavior, local `as_of` rendering, and no-cache refresh header/update behavior. `npm run build` passes in apps/web, and the installed TanStack version is recorded from `npm view`.

## Non-Goals

- Do not migrate Dashboard, MarketplaceSettings, IntegrationsHub, PricingSimulator, Classifications, or existing catalog lookups outside the assigned views.
- Do not add mutation invalidation; F-02 owns mutation call sites.
- Do not edit SDK transport/types, server code, OpenAPI, or create a parallel fetch layer.
- Do not retain the legacy enrichment ProductsPage or its tests/export/route wrapper/nav entry.

## Design

Create the lean `@marketplace-central/web-query` workspace package as the single source of reserved queryKey builders, IC-01 stale-time constants, QueryClient factory/defaults, as-of formatting, and a mutable no-cache fetch controller. `ClientContext` owns the controller and passes its fetch implementation to the SDK client; migrated feature pages receive the client through existing AppRouter wrappers and use the shared query client/hooks. The catalog package owns catalog query hooks/page, inventory owns the stock-risk query, and orders owns the profitability query while retaining existing order/installation/manual-write behavior around those reads.

Catalog facts use `['catalog','facts',{params}]`, `useInfiniteQuery` with an undefined initial page parameter and `next_cursor ?? undefined`, plus a 300ms debounced `['catalog','search',q]` query disabled for empty input. Inventory uses `['inventory',{installation_id,filters}]` with 45s staleTime; profitability uses `['profitability',{installation_id}]` with 120s staleTime. Refresh wraps query refetch in the client no-cache controller.

## Edge Cases

- `next_cursor: null` is terminal and must not cause another request.
- Null numeric facts render unknown/`—` and their quality flags; no zero/default substitution is allowed.
- Empty catalog search is disabled; search results have no pagination continuation.
- No installation selected leaves installation-dependent reads idle until selection is available.
- SDK/server error responses render explicit error states and retry controls rather than blank or zeroed data.
- Concurrent refreshes are prevented by disabling the control while the no-cache refetch burst is in flight.

## Acceptance Criteria

- Criterion: Infinite catalog pagination over cursor envelope. Traces to milestone criterion ID: M-05-C01. Proven by: `npm test` component tests in apps/web with mocked SDK, asserting appended pages, verbatim `next_cursor`, null termination, and no page-1 refetch.
- Criterion: staleTime prevents redundant refetch. Traces to milestone criterion ID: M-05-C02. Proven by: component tests with a fetch/client spy asserting one call for remount within the configured stale window and a refetch after the catalog stale window.
- Criterion: as_of indicator and manual refresh. Traces to milestone criterion ID: M-05-C03. Proven by: component test rendering an envelope, clicking refresh, asserting `dados de HH:mm:ss`, `Cache-Control: no-cache`, and the updated indicator.
- Criterion: Build health and no bypassing fetches. Traces to milestone criterion ID: M-05-C05. Proven by: `npm run build` in apps/web and scoped search showing migrated views call SDK methods through query hooks with no direct fetch.

## Handoff

- Current status: `spec_ready`
- Next owner: Feature Implementer
- Next action: Write `plan.md`, implement the scoped feature, and capture validation evidence.
- Required files/evidence: feature brief, IC-01, M-05 milestone/validation contract, exact command output, and changed-path list.
- Blockers or open decisions: None.
