# F-02-mutation-invalidation — Spec

## Goal
Every Oracle-backed write in `apps/web` invalidates exactly the TanStack Query namespaces IC-01's Invalidation Crosswalk assigns it, so a successful mutation makes dependent reads observably fresh without any manual refetch path.

## Behaviour

### Invalidation Crosswalk (IC-01)
| Write | Site | Invalidates |
|---|---|---|
| linkage confirm | `ProductLinksPage.approveProductLinkCandidate` | `['linkage']` + `['catalog']` |
| linkage reject | `ProductLinksPage.rejectProductLinkListing` | `['linkage']` + `['catalog']` |
| linkage manual resolve | `ProductLinksPage.manualResolveProductLink` | `['linkage']` + `['catalog']` |
| margin-input import | `OrdersPage.importProfitabilityMarginInputs` | `['catalog']` + `['profitability']` |
| stock action | `StockSeguroPage.applyInventoryManualStockAction` | `['inventory']` |

- Invalidation is `onSuccess`-only. A rejected mutation invalidates nothing and surfaces `actionError`.
- No optimistic updates (out of scope; rollback surface for no contract value).

### Linkage is never cached
The linkage workflow read uses `staleTime: 0` and `gcTime: 0` (TanStack v5 naming). Linkage is authorization-sensitive; it is never served from cache. The manual `useEffect` load and manual post-action reload are removed — `invalidateQueries` is the single refresh path.

### Key discipline
All queryKeys and invalidation targets come from `@marketplace-central/web-query`. Namespaces are limited to the four IC-01 reserved roots: `['catalog']`, `['inventory']`, `['linkage']`, `['profitability']`. Zero inline queryKey literals in feature packages.

## Out of scope
CatalogPage (F-01), Dashboard, MarketplaceSettings, IntegrationsHub, PricingSimulator, Classifications. No Go, OpenAPI, or `sdk-runtime` edits.
