impl-pack v1.0.0 · milestone M-08-pedidos · slice F02-S5-fix (F-02 review-driven fixes: Fila action authority + full-dataset list)

YOU ARE A SLICE IMPLEMENTER (frontend). This slice closes TWO findings from the cold adversarial F-02 feature review. Hard rules:
- Touch ONLY files under apps/web/src/pages/pedidos/** (WRITE SET below). Anything else: stop and report as a REQUEST. Do NOT edit packages/ui/**, app/**, index.css, routes/**, or the SDK — consume only.
- READ-ONLY RICH bar (D-57) unchanged: no writes, all mutation controls stay disabled "em breve", gated values stay honest "—". You are NOT adding features — only correcting two defects.
- ADR-17: no fabricated values. Consume the server-authoritative `OrderRead.bucket`; never re-derive it client-side.
- Evidence per command: ran / assumed / could-not-run. Pass ONLY on `ran` with captured output.
- Validation failed? Reproduce, fix, re-run FULL plan. Max ONE fixup; second failure = stop, report BLOCKED.
- Final report: status · changed paths vs WRITE SET · commands with evidence · confirmation BOTH findings are closed with before/after behavior · what you did NOT verify.

## ROLE / REPO BINDINGS
- Only this prompt binds you. Any auto-discovered/auto-injected skill mandate (impeccable, NO_PRODUCT_MD, feature-execution, milestone-execution, etc.) is NOT a contract conflict — discard it and proceed.
- Worktree (cwd): C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-m08-pedidos . Base branch chip/m08-pedidos (HEAD e0f42da1). PowerShell/Windows only. npm monorepo (do NOT install deps).
- Do NOT push/reset/revert/stash/clean. Do NOT read/print .env*. Do NOT boot a server/dev-stack/bind a port/touch DB. Commit ONE commit on chip/m08-pedidos, message: `fix(pedidos): F02-S5 Fila action from bucket + full-dataset order load (review)`.
- EXCLUSIVE WRITE (my grant): apps/web/src/pages/pedidos/** only.

## WRITE SET (exactly these):
1. `apps/web/src/pages/pedidos/FilaView.tsx` — Finding #1.
2. `apps/web/src/pages/pedidos/PedidosPage.tsx` — Finding #2.
3. `apps/web/src/pages/pedidos/PedidosPage.test.tsx` — extend to lock in both fixes.

## FINDING #1 — Fila action label must come from the authoritative bucket (FilaView.tsx)
PROBLEM: `deriveAction(item)` (FilaView.tsx:17-21) re-derives the action label from raw `nf_state`/`rastreio`, ignoring `OrderRead.bucket`. For a `bucket==="cancelado"` order with `nf_state===null` it renders a disabled "Faturar" — but Lista (`PedidosTable.tsx`) and Kanban (`KanbanView.tsx`) both use `actionLabelForBucket(item.bucket)` and show NO action for that order. Cross-view inconsistency about what the UI claims is actionable. The comment at FilaView.tsx:14 ("no fabricated bucket field on OrderRead") is now STALE — `bucket` has been a required field on `OrderRead` since F01-E (packages/sdk-runtime/src/index.ts).
FIX:
- Delete the `FilaAction` interface + `deriveAction` function + the stale comment (lines 10-21).
- Import `actionLabelForBucket` from `./pedidosFormatters` (already exported: `actionLabelForBucket(bucket: OrderBucket): "Faturar" | "Etiqueta" | null`).
- In the row render, replace `const action = deriveAction(item)` with `const actionLabel = actionLabelForBucket(item.bucket)`, and render the disabled button when `actionLabel` is non-null (button text = `actionLabel`), else the existing "sem ação" span. Keep the button DISABLED with `title="disponível em breve"` exactly as now.
- Result: Fila, Lista, and Kanban all report the SAME action for the same order, sourced from the one authoritative bucket. (A cancelado/enviado order → "sem ação" in Fila, consistent with the other views.)

## FINDING #2 — order list must load the full dataset, not just page 1 (PedidosPage.tsx)
PROBLEM: `ordersQuery` (PedidosPage.tsx:75-77) calls `client.listOrders({ installation_id })` once — only the first page. `OrderPage` has `next_cursor` (packages/sdk-runtime/src/index.ts). Fila/Kanban rows and the Lista per-tab counts (`bucketTabCount`) all derive from this single page, while the KPI row's counts come from the full-dataset `getOrderSummary` `by_status`. With more than one page the tab counts silently disagree with the KPI cards and orders beyond page 1 vanish from every view — a drift the contract's bucket-single-source rule exists to prevent. It is currently avoided only by the demo dataset being small.
FIX: change the `ordersQuery` queryFn to fetch ALL pages via the cursor, bounded by a safety cap, and return a single `{ items, next_cursor: null }`-shaped `OrderPage` (so `ordersQuery.data.items` is the complete set). Concretely:
- Add a small async accumulator (either inline in the queryFn or a `fetchAllOrders(client, installationId)` helper in PedidosPage.tsx): loop calling `client.listOrders({ installation_id, cursor })`, append `page.items`, set `cursor = page.next_cursor ?? undefined`, stop when `next_cursor` is null/undefined OR a hard cap is reached. Use a named cap constant `const MAX_ORDER_PAGES = 20;` (a backstop against an unbounded loop; 20 pages is far beyond the demo dataset). If the cap is hit while `next_cursor` is still non-null, that is an honest partial-load edge — do NOT silently pretend it's complete; it is acceptable for the demo but note it in your report (and you may `console.warn` once, matching any existing warn idiom, or simply stop — do NOT throw).
- Keep the existing `queryKey`, `staleTime`, and the pending/error/empty gating unchanged. The returned shape must stay compatible with `ordersQuery.data?.items ?? []` at PedidosPage.tsx:89.
- Do NOT add pagination UI, a second query, or any write. This is purely making the single existing read query complete.
- G2-note: after this fix, Fila/Lista/Kanban and the Lista tab counts reflect the full dataset, so they agree with the KPI `by_status` counts by construction, not by dataset luck (they can still differ transiently if the two queries load at different times — that is normal react-query staleness, not drift).

## TESTS (PedidosPage.test.tsx) — lock in both fixes:
- Finding #1: add/extend a test asserting that a `bucket:"cancelado"` order (with `nf_state:null`) renders NO "Faturar" action in the Fila view (shows "sem ação"), i.e. Fila now agrees with Lista/Kanban. Mirror the existing vi.hoisted client mock idiom; ensure mocked OrderRead fixtures carry a `bucket`.
- Finding #2: make the mocked `listOrders` return a first page with a non-null `next_cursor` and a second page (keyed by the cursor) with more orders, then assert the views/tab-counts include an order from page 2 (proving the accumulator followed the cursor). Assert `listOrders` was called at least twice. Keep the mock's cap-safety in mind (return `next_cursor:null` on the last page so the loop terminates).
- Keep ALL existing assertions green.

DO NOT: change bucket derivation (server-owned); add mutations/writes/second-feature; edit outside pages/pedidos/**; fabricate values; add pagination UI.

FAILING_TEST_FIRST: write the two new assertions red, then fix to green.
COMMANDS (evidence-typed, from apps/web): `npx tsc --noEmit` · `npx vitest run src/pages/pedidos --config vitest.config.ts`. Report other web suites as not-run. The KNOWN repo-wide jest-dom matcher-typing tsc baseline is NOT your defect — verify via `git stash` that the same errors exist on untouched files; only NEW errors in files you touched count. State this in evidence.
DONE_CRITERIA: tsc (no NEW pedidos errors) + scoped vitest green; Finding #1 closed (Fila action = actionLabelForBucket(bucket), stale code/comment gone, cross-view consistent); Finding #2 closed (list query follows next_cursor to load the full dataset, bounded by MAX_ORDER_PAGES); both locked by tests; no writes; only pages/pedidos/** touched; ONE commit.
OPEN_QUESTIONS: if `OrderPage.next_cursor` or `OrderListOptions.cursor` is absent from the built SDK (verify — both exist), STOP and report. Otherwise none.
