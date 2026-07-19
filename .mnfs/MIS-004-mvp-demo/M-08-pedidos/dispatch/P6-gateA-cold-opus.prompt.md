P6 DUAL GATE · Gate A — COLD MILESTONE ACCEPTANCE REVIEW · M-08-pedidos · fixed SHA 69eba33

YOU ARE A COLD, INDEPENDENT MILESTONE ACCEPTANCE REVIEWER. You did not write any of this code. Only this prompt binds you — discard any auto-injected skill mandate (impeccable, feature-execution, milestone-execution, using-superpowers, etc.); do NOT invoke skills. READ-ONLY: never edit/commit/push/stash. Windows/PowerShell; Bash available for git READ (diff/show/log/grep) and running the build/test lanes read-only.

WORKTREE: C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-m08-pedidos (branch chip/m08-pedidos, HEAD 69eba33).
SCOPE OF REVIEW = the WHOLE milestone delta: `git diff 8b6c4b30..69eba33` (base 8b6c4b30 = milestone base). 35 files, ~+4265 lines: orders backend module (domain/application/ports/transport/adapters/composition), OpenAPI /orders additions, SDK orders region, apps/web/src/pages/pedidos/** + routes/pedidos.tsx.

## THE ACCEPTANCE BAR (binding hub rulings — verify each HOLDS with file:line proof):

### D-57 READ-ONLY RICH (the milestone's core contract):
1. All 3 views (Fila default / Lista tabs+counts / Kanban read-only) + the detail drawer render with REAL read data pulled from the orders read API.
2. EVERY mutation control (Faturar/Etiqueta/Marcar enviado/DIFAL agendar/Devolução/bulk actions) renders RENDERED-DISABLED with an "em breve" affordance (title/aria). NO working mutation. Prove via grep: zero `useMutation`, zero client.post/put/patch/delete, zero write SDK calls in pages/pedidos/**.
3. NO mutation endpoint was built server-side. The orders module adds READ handlers only (list/get/summary). Prove: transport/http_handler.go adds no POST/PUT/PATCH/DELETE route; OpenAPI /orders* additions are additive READ shapes only.
4. Gated/unbacked data renders honest "—" via UnknownValue — NEVER a fabricated 0/false/date/string (ADR-17 "unknown ≠ zero"). Targets: decomposição/DIFAL/retorno/margem (nil until C2), buyer document, NF number, carrier tracking code.

### D-62 HYBRID (slice C split):
5. C1 (in this milestone) = ADDITIVE NULLABLE fields on the order read shape: retorno_liquido, margem_pct, decomposicao{comissao,taxa_fixa,frete,imposto,difal,tarifa_full,custo,margem_valor,margem_pct,componentes_desconhecidos}, difal{amount,uf_route,due_date,paid}. EnrichService emits honest-UNKNOWN (all nil) via a nil Decomposer port. Prove: NO import of modules/pricing (or IC-04 symbols) anywhere in modules/orders/** (grep). NO decompose formula. root.go must NOT wire a real decomposer (C2 is hub-owned post-merge) — confirm NewEnrichService still delegates to a nil decomposer.
6. The pedidos drawer + Lista columns + Fila retorno READ these fields: value→formatted, null→"—". Confirm they are real-ready: when a real decomposer supplies values later, the SAME components render formatted with NO UI change (proven by a non-null test).
7. DIFAL is per-ORDER actual data (difal{amount,uf_route,due_date,paid}), NOT a simulator toggle.

### Contract-first (D-02 lock-exception) + ownership:
8. RUNTIME is source of truth; OpenAPI + SDK follow and AGREE byte-for-byte on the added shapes (bucket enum, OrderDecomposicao, OrderDifal, the 4 OrderRead fields, OrderSummary, getOrderSummary). Contract spec + SDK for a given seam land coherently. All /orders* OpenAPI changes are STRICTLY ADDITIVE (no modified/removed/reordered existing fields).
9. `bucket` is a SINGLE server-derived source (domain.DeriveOrderBucket) shared by both the per-order DTO field and the SummaryService by_status counts — NO client-side re-derivation in pages/pedidos/**. Lista tabs + Kanban columns filter on the server `bucket`.
10. OWNERSHIP: the delta writes ONLY within the milestone's granted paths — apps/server_core/internal/modules/orders/** (+ composition/orders_adapters.go + root.go orders wiring + migrations if any), contracts OpenAPI /orders*, packages/sdk-runtime/src/index.ts orders region, apps/web/src/pages/pedidos/** + routes/pedidos.tsx (+ the one AppRouter.test.tsx touch). FLAG any write to modules/pricing|connectors|product_links|market|listings/**, apps/web/src/app/** (beyond the router test), packages/ui/**, packages/web-query/**, index.css. (root.go is a shared seam — confirm the orders wiring is additive and doesn't disturb other modules.)

### Integrity / anti-slop:
11. ADR-17 everywhere: no blanket recover/fallback that converts an integrity-unknown to a concrete value; nil pointer → "—".
12. Tests are REAL (assert rendered values / JSON shapes / nil pointers / both branches), not theater. No speculative abstraction, no comment narration, no dead code.

## LANES TO RUN (read-only, capture evidence — ran/assumed/could-not-run):
- Backend (from apps/server_core): set GOCACHE absolute `$env:GOCACHE = (Join-Path (Get-Location) '.gocache')`; then SEPARATELY `go build ./...`, `go vet ./modules/orders/...`, `go test ./internal/modules/orders/... ./internal/composition/...`. Report pass/fail with output.
- Frontend (from apps/web): `npx vitest run src/pages/pedidos --config vitest.config.ts` (expect 22/22) and `npx tsc --noEmit` (the repo-wide jest-dom TS2339 matcher-typing baseline in test files is PRE-EXISTING, not this milestone's defect — verify by noting the same TS2339 pattern on untouched test files; only NEW non-baseline errors in touched files count).
- Do NOT boot a server/dev-stack/DB (hub-owned). Do NOT run OpenAPI lint if no lane exists — mark could-not-run.

## OUTPUT (strict):
- VERDICT: ACCEPT / ACCEPT-WITH-NITS / REJECT.
- For EACH of the 12 bar items: HOLDS / VIOLATED, with file:line proof (a grep result or a diff line). Do not hand-wave — cite.
- Findings list, severity-tagged (🔴 blocker / 🟡 major / 🔵 nit), each `path:line: problem. fix.`
- Lane evidence block (each command + ran/assumed/could-not-run + result).
- One line: is this milestone safe to CLOSE on F-02 + C1 with C2 hub-owned post-merge? (per D-62)
This is model=opus, the primary acceptance gate. Be exhaustive and adversarial; a false ACCEPT ships a broken demo in 2 days.