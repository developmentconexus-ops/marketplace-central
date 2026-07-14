# F-02-mutation-invalidation — Plan

1. **Extend `packages/web-query`** — add `queryKeyNamespaces` (the four IC-01 reserved roots) and `linkageQueryKeys.workflows(installation_id)`. No new namespaces.
2. **ProductLinksPage** — migrate the workflow read to `useQuery` keyed by `linkageQueryKeys.workflows(...)` with `staleTime: 0` / `gcTime: 0`; wrap the three linkage writes in `useMutation` with `onSuccess` → invalidate `['linkage']` + `['catalog']`; delete the manual `useEffect` load and the manual post-action reload.
3. **OrdersPage** — wrap `importProfitabilityMarginInputs` in `useMutation`; `onSuccess` → invalidate `['catalog']` + `['profitability']`.
4. **StockSeguroPage** — wrap `applyInventoryManualStockAction` in `useMutation`; `onSuccess` → invalidate `['inventory']`; drop the manual post-action refetch F-01 left in place.
5. **Tests** — per-write invalidation assertions (exact namespaces, unrelated namespaces untouched), a failure path asserting zero invalidation, and linkage cache options.
6. **Verify** — feature suites, `npm run build`, canonical `npm test`, inline-queryKey grep; one intentional commit.

## Execution note
Steps 1–4 were implemented by Codex (`gpt-5.6-luna`, high reasoning) under `-s workspace-write`. Two successive host teardowns killed the Codex process mid-run before it could commit or self-report, and the second teardown killed it before it applied the correction in step 7. The orchestrator applied step 7 directly and ran all verification. See `validation.md` and the milestone checkpoint.

7. **Correction (applied by orchestrator)** — fix two defects found by orchestrator verification:
   - `mutationFn` passed by reference at 5 sites, leaking TanStack v5's second `context` argument into SDK calls;
   - tests building their QueryClient from `createWebQueryClient()` (production `retry: 1`), so error-state assertions raced the retry delay.
