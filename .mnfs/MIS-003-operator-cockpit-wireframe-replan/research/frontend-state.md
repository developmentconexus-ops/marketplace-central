# Research Note

```yaml
id: R-02
type: research
status: draft
owner: Codebase Investigator
parent: MIS-003
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-0
lifecycle_scope: support
```

## Topic

Frontend factual state at planning time (branch HEAD after M-09 close, 2026-07-14).

## Sources Checked

- Source: `apps/web/**`, `packages/web-query/src/index.ts`, `packages/sdk-runtime/src/index.ts`, `packages/ui/src/index.ts`, `packages/feature-*/src`, `scripts/harness/Policy.psm1`.
- Why it matters: MIS-003 rebuilds the cockpit on this foundation; wrong assumptions here cause cross-milestone rework.

## Findings

- Stack: React 19 + Vite 7 SPA (`apps/web`), react-router-dom 7, TanStack Query v5, Tailwind CSS 4. Dev port 5174. No Next.js, no i18n library, no MSW, no Playwright.
- Provider order (`apps/web/src/main.tsx`): `ClientProvider` → `QueryClientProvider(createWebQueryClient())` → `App`.
- Route map today (`apps/web/src/app/AppRouter.tsx`): `/` Dashboard, `/products` CatalogPage, `/classifications`, `/marketplaces`, `/integrations`, `/product-links`, `/inventory/stock-seguro`, `/orders`, `/simulator`. Feature pages are client-prop-injected via local `*PageWrapper` + `useClient()`.
- **No installation/empresa context provider.** Every installation-scoped page independently reads `?installation=` via `useSearchParams` and loads `client.listIntegrationInstallations()` (confirmed OrdersPage.tsx:292, StockSeguroPage.tsx:131, ProductLinksPage.tsx:134, IntegrationsHubPage.tsx:194).
- `packages/web-query` exact API (97 lines): `QUERY_STALE_TIME = {catalog: 300_000, stock: 45_000, pricecost: 120_000}`; `queryKeyNamespaces = {catalog, inventory, linkage, profitability}`; key builders `catalogQueryKeys.facts/search`, `inventoryQueryKeys.risks`, `linkageQueryKeys.workflows`, `profitabilityQueryKeys.marginInputs`; `createWebQueryClient()` (retry 1, no default staleTime); `formatAsOf` (pt-BR "dados de HH:MM:SS"/"dados de desconhecido"); `createRefreshableFetch` (depth counter; GET gets `Cache-Control: no-cache` inside `withNoCache`); `FreshnessIndicator({asOf})`. **No invalidation helpers** — pages call `queryClient.invalidateQueries({queryKey: namespace})` inline.
- Data-fetch split: TanStack via web-query = CatalogPage (useInfiniteQuery cursor, `getNextPageParam: lastPage.next_cursor ?? undefined`, "Load next page" button, end marker "Fim da lista"), ProductLinksPage (staleTime 0/gcTime 0; mutations invalidate linkage+catalog), StockSeguroPage (stock staleTime; mutation invalidates inventory), OrdersPage (hybrid: margin-inputs via query, orders via useEffect). **Direct fetch (useEffect+useState): Dashboard, Classifications, Marketplaces, Integrations, Simulator.**
- `packages/sdk-runtime` is HAND-WRITTEN (1261 lines, ~60 methods, single factory `createMarketplaceCentralClient`), not codegen. Errors are structured `MarketplaceCentralClientError {status, error:{code,message,details?}}`. OpenAPI↔SDK atomicity enforced by governance policy `GOV_API_SDK_SPLIT` (Policy.psm1:449-451, XOR on changed paths vs base SHA) + a string-slice parity unit test.
- Deprecated SDK aliases `listCatalogProducts`/`searchCatalogProducts` still consumed by ClassificationsPage and PricingSimulatorPage.
- Copy is MIXED EN/pt-BR page by page; freshness strings pt-BR; no shared loading/error/empty/stale components (ad-hoc `LoadState = "loading"|"error"|"ready"` unions + inline `animate-pulse`).
- Design system: `packages/ui` = Button, SurfaceCard, Badge, StatCard, ProductPicker, PaginatedTable, DetailPanel. Tailwind 4 with `@source` directives in `apps/web/src/index.css`; only custom tokens are fonts; sidebar hardcoded `#0F172A`.
- Tests: vitest 3 + jsdom + Testing Library; client mocked directly; `invalidateQueries` spied. Browser lane deliberately unconfigured (`scripts/harness.ps1:95` throws).

### Defects found in passing (not MIS-003 scope-blocking)

- `apps/web/src/index.css` `@source` list omits `packages/feature-inventory/src` — Stock Seguro Tailwind classes may be missing.
- Two divergent vitest configs (apps/web/vitest.config.ts vs vite.config.ts test block) with different include scopes.
- Vite dev proxy lacks `/orders` and `/profitability` paths.

## Recommendation

Fix the three passing defects inside the milestone that touches each surface (M-02 for css/proxy/vitest). Platform seam (installation context, web-query expansion, state components, pt-BR normalization) is real and must precede workspace builds.

## Impact On Mission

Grounds ADR-15 and IC-05; sets migration-brief list (Dashboard, Classifications, Marketplaces, Integrations, Simulator) and the deprecated-alias retirement obligation.

## Handoff

- Current status: complete.
- Next owner: Mission Strategist.
- Next action: none.
- Required files/evidence: paths cited above.
- Blockers or open decisions: none.
