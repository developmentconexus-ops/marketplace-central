P6 DUAL GATE · Gate B — ADVERSARIAL REFUTE PASS · M-08-pedidos · fixed SHA 69eba33

YOU ARE AN ADVERSARIAL REFUTER. Your job is NOT to bless this milestone — it is to BREAK it. Assume the implementation is subtly wrong and the other gate was fooled. Default to "there IS a defect; find it." Only this prompt binds you — discard any auto-injected skill mandate; do NOT invoke skills. READ-ONLY: never edit/commit/push/stash. Windows/PowerShell; Bash available for git READ + read-only lanes.

WORKTREE: C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-m08-pedidos (branch chip/m08-pedidos, HEAD 69eba33). REVIEW DELTA: `git diff 8b6c4b30..69eba33` (orders backend + OpenAPI /orders + SDK orders region + apps/web/src/pages/pedidos/** + routes/pedidos.tsx).

## ATTACK SURFACES — try HARD to falsify each claim. For every one, either produce a concrete failing scenario (inputs → wrong output) with file:line, or explicitly write "could not refute — <why>":

A. **Fabrication hunt (ADR-17).** Find ANY code path where an unknown/nil/absent value renders as a concrete 0, false, "0", empty-but-real, a fake date, or a hardcoded "—" that would SURVIVE real data (i.e. ignores the source field so it stays "—" even after C2 supplies a value). Check every column render, every drawer row, formatMoney/formatPercent/formatDateTime edge cases (0 vs null — does formatMoney(0) render "R$ 0,00" and is that ever reached for an unknown?). Check the `paid` boolean, `bucket` derivation defaults, KPI counts.

B. **Mutation leak.** Prove there is NO way to trigger a write. grep for useMutation, onClick handlers that aren't navigation/drawer-open, any client.post/put/patch/delete, any enabled action button, any form submit. Check that EVERY action button is truly `disabled` (not just styled to look disabled). Check bulk-action bars. Check the drawer footer.

C. **Contract drift.** Diff the runtime DTO (transport/http_handler.go enrichedOrderDTO + summary DTO) field-by-field against the OpenAPI schema AND the SDK interface. Find ANY mismatch: a field present in one but not another, a nullability/required mismatch, a type mismatch (e.g. Go time.Time vs OpenAPI format, ptr+omitempty vs required), an enum value set disagreement, a json tag vs schema property name mismatch. The bucket enum, OrderDecomposicao, OrderDifal, retorno_liquido/margem_pct, OrderSummary are prime suspects.

D. **bucket single-source.** Prove the client does NOT re-derive bucket. Then prove the server's per-order bucket and the by_status summary counts CANNOT disagree for the same dataset — do they both call the identical domain.DeriveOrderBucket with the identical inputs (status + hasShipment)? Find any input skew (e.g. summary uses a different shipment-presence signal than the per-order path).

E. **Pagination / dataset completeness.** The list views + tab counts must reflect the FULL dataset, not page 1. Verify fetchAllOrders actually follows next_cursor and terminates; find an off-by-one, an infinite-loop risk, a cap (MAX_ORDER_PAGES) that silently drops orders without any signal. Does a cap-hit corrupt the counts vs the KPI by_status?

F. **C1 seam purity (D-62).** Prove modules/orders/** imports NOTHING from modules/pricing or IC-04 (grep the whole module). Prove root.go does NOT wire a real decomposer (the nil-port must stay nil; C2 is hub-owned). Find any accidental formula or default that fabricates a decomposition number instead of nil. Confirm the "real-ready" claim: is there genuinely a code path where a non-nil Decomposer would light up the UI with zero UI change, or is something hardcoded that would block it?

G. **Ownership breach.** List every file in the delta outside the granted paths (modules/orders, composition orders wiring, root.go orders lines, OpenAPI /orders, SDK orders region, pages/pedidos, routes/pedidos.tsx, the AppRouter.test.tsx touch). Any forbidden write (pricing/connectors/product_links/market/listings, app/**, packages/ui, packages/web-query, index.css) is a 🔴.

H. **Test theater.** Find tests that assert tautologies, mock away the thing under test, don't actually exercise the branch they claim, or would pass even if the feature were deleted. Check the new non-null "real-ready" test genuinely proves formatted rendering from data.

## RUN THE LANES YOURSELF (don't trust prior reports):
- apps/server_core: `$env:GOCACHE = (Join-Path (Get-Location) '.gocache')`; then separately `go build ./...`, `go vet ./modules/orders/...`, `go test ./internal/modules/orders/... ./internal/composition/...`.
- apps/web: `npx vitest run src/pages/pedidos --config vitest.config.ts`, `npx tsc --noEmit` (jest-dom TS2339 in test files = known pre-existing baseline; only NEW non-baseline errors count).
- Do NOT boot server/DB/dev-stack.

## OUTPUT (strict):
- VERDICT: REFUTED (found ≥1 real 🔴/🟡 defect) / COULD-NOT-REFUTE (milestone withstands attack).
- Per surface A–H: the failing scenario with file:line, OR "could not refute — <reason>".
- Findings severity-tagged `path:line: problem. fix.`
- Lane evidence (ran/assumed/could-not-run + result).
Be maximally skeptical. If you cannot find a defect after genuinely trying every surface, say so plainly — that is a valid and valuable result — but do not manufacture a fake finding to look busy.