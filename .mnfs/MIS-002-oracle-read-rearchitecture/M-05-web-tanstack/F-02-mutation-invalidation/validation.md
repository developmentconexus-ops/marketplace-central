# F-02-mutation-invalidation — Validation

All commands run by the Milestone Orchestrator on a quiet machine, real exit codes captured (never piped through `tail`/`head`, which masks the runner's exit code).

## Commands + results

| Command | Exit | Result |
|---|---|---|
| `npx vitest run --root . packages/feature-product-links packages/feature-inventory packages/feature-orders packages/feature-products --testTimeout=30000` | **0** | **32 passed / 32** (4 files) |
| `npm run build` (in `apps/web`) | **0** | ✓ 1832 modules transformed, built in 1m 27s |
| `npm test` (repo root, canonical) | **0** | **11 passed / 11** (4 files) |

Build warnings are pre-existing `"use client"` directive notices from `@tanstack/react-query`, `react-router`, and `lucide-react` — bundler-level, unrelated to this feature.

## M-05-C04 — mutations invalidate the correct namespaces

Crosswalk assertions (spy on `queryClient.invalidateQueries`, asserting exact namespaces AND that unrelated namespaces are untouched):
- `approves a candidate and refreshes workflows` — `['linkage']` + `['catalog']`
- `rejects a listing identity` — `['linkage']` + `['catalog']`
- `manually resolves a workflow` — `['linkage']` + `['catalog']`
- `invalidates catalog and profitability after importing margin inputs` — `['catalog']` + `['profitability']`
- `applies manual action and shows result message` — `['inventory']`

Failure paths (a rejected mutation must invalidate nothing):
- `registers linkage as uncached and does not invalidate after a failed write`
- `does not invalidate after a failed margin-input import`
- `does not invalidate inventory after a failed stock action`

Linkage cache options: `registers linkage as uncached ...` asserts `staleTime: 0` and `gcTime: 0` (v5 naming) on the linkage workflow query.

## Key discipline
`grep -rnE 'queryKey: \[("|'\'')' packages/feature-*/src/ apps/web/src/` → **no matches**. All keys and invalidation targets resolve through `@marketplace-central/web-query`. Namespaces limited to the four IC-01 reserved roots; none invented.

## Seams
`git status --porcelain | grep -E '\.go$|sdk-runtime|openapi'` → **no matches**. No Go, OpenAPI, or `sdk-runtime` edits. `f02-packet.md` / `f02-fix.md` are git-excluded and uncommitted.

## Contract note — IC-01 "product edit → ['catalog']" is MOOT
IC-01's Invalidation Crosswalk carries a `product edit → ['catalog']` row. F-01 implemented the orchestrator's Option B decision and **deleted the legacy `ProductsPage`**, which held the only product-edit write in the app. No product-edit write now exists, so the row has no implementation target. It is satisfied vacuously, not by code. No page was invented to satisfy it. Flagged for QA as a contract/scope interaction rather than a coverage gap.

## Defects found and fixed during verification
Codex's implementation was contract-correct on the crosswalk, namespaces, cache options, and error handling, but two defects were found by orchestrator verification (Codex was killed by host teardown before it could self-verify or commit):

1. **Production defect — `mutationFn` passed by reference at 5 sites.** TanStack v5 invokes `mutationFn(variables, context)`, so `mutationFn: client.approveProductLinkCandidate` forwarded a second `{client, meta, mutationKey}` argument into the SDK method and unbound it from the client object. Proven by spy output:
   ```
   1st spy call:
     [ { "installation_id": "inst-1", "limit": 50 },
   +   { "client": QueryClient {}, "meta": undefined, "mutationKey": undefined } ]
   ```
   Fixed by wrapping each site in an arrow passing only the request, typed via `Parameters<Client["method"]>[0]` (no `any`).

2. **Test defect — error assertions raced the production retry.** Tests built their QueryClient from `createWebQueryClient()`, which carries the production default `retry: 1`, so `renders error state when workflow loading fails` waited out the retry delay and never saw the error text. Fixed with a test-local `createTestQueryClient()` that derives from the production factory and disables only `retry`. **Production defaults (`retry: 1`, `refetchOnWindowFocus: false`) are unchanged** — they are contract-fixed.

## Note on a discarded test run
An earlier run reported 8 failures, including F-01's CatalogPage tests. That run executed concurrently with a `npm run build` and leftover Codex processes; its timings (environment 538s, setup 171s, failures all "timed out in 5000ms") show machine contention, not defects. Re-run on a quiet machine it completed in 7.4s and CatalogPage passed. Discarded as invalid; the 5 genuine failures it masked are the two defects above.

## Changed paths
```
packages/web-query/src/index.ts                          (+11)  queryKeyNamespaces, linkageQueryKeys
packages/feature-product-links/src/ProductLinksPage.tsx         useQuery(staleTime 0/gcTime 0) + 3 mutations
packages/feature-product-links/src/ProductLinksPage.test.tsx
packages/feature-product-links/package.json               (+2)  @tanstack/react-query
packages/feature-orders/src/OrdersPage.tsx                      margin-input import mutation
packages/feature-orders/src/OrdersPage.test.tsx
packages/feature-inventory/src/StockSeguroPage.tsx              stock action mutation
packages/feature-inventory/src/StockSeguroPage.test.tsx
.mnfs/.../F-02-mutation-invalidation/{spec.md,plan.md,validation.md}
```
