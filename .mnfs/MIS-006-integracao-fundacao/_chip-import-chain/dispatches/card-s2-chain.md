SLICE CARD — S2 · the import detail screen reads the REAL chain from the server

## Why this slice exists (read before writing)

`GET /erp/imports/{id}/chain` landed in a previous chip with OpenAPI + SDK, and **no screen has ever called
it**. This slice is the first consumer, and the live drive of this screen is what discharges an operator
waiver. Two things are therefore non-negotiable:

- The decomposition counters come from `client.getErpImportChain(id)` and from nowhere else. A counter you
  compute on the client from `listErpImports` is the LOOK of the chain without its truth — that shortcut is
  the exact thing this chip exists to prevent. Do not derive, sum, or infer any of the three numbers.
- A counter that the payload does not carry renders `—`, never `0`. See the integrity section of the bindings.

## write_set (nothing else)

- `apps/web/src/pages/importacoes/useErpImportChain.ts`          (NEW)
- `apps/web/src/pages/importacoes/ImportChainPanel.tsx`           (NEW)
- `apps/web/src/pages/importacoes/ImportChainPanel.test.tsx`      (NEW)
- `apps/web/src/pages/importacoes/ImportacaoDetailPage.tsx`       (NEW)
- `apps/web/src/pages/importacoes/ImportacaoDetailPage.test.tsx`  (NEW)
- `apps/web/src/pages/importacoes/ImportacaoSection.tsx`
- `apps/web/src/pages/importacoes/ImportacaoSection.test.tsx`
- `apps/web/src/routes/importacoes.tsx`
- `apps/web/src/app/AppRouter.tsx`

## Contract facts (verified — do not re-derive)

- SDK method: `client.getErpImportChain(importId: string): Promise<ErpImportChain>`
  (`packages/sdk-runtime/src/index.ts:1901`).
- `ErpImportChain` (`packages/sdk-runtime/src/erpImport.ts:48`):
  `{ protocol: string; importados: number; vinculados: number; enfileirados: number; queue_read_at: string }`.
  All five are declared REQUIRED and non-nullable in `contracts/api/marketplace-central.openapi.yaml:8078`.
  **The runtime guard stays anyway** — a required field that arrives absent is version drift between server
  and client, and drift is exactly when honesty is worth the most. Do NOT change the OpenAPI or the SDK; the
  types are correct as declared and this slice must not touch `contracts/**` or `packages/sdk-runtime/**`.
- SDK errors are thrown as plain objects `{ status: number, error: string }` (see
  `packages/sdk-runtime/src/index.ts:1715`). For a missing import the server answers `404` with
  `error: "import_not_found"` (`contracts/api/marketplace-central.openapi.yaml:3277`).
- The app's shared QueryClient defaults to `retry: 1` (`packages/web-query/src/index.ts:82`).
- Existing page-local query-hook idiom to match: `apps/web/src/pages/vinculos/useErpImports.ts` (exported
  `*QueryKeys` object + `useQuery` with `QUERY_STALE_TIME.listings` from `@marketplace-central/web-query`).
  You may NOT edit that file; write your own module in `pages/importacoes/`.
- `apps/web/src/routes/importacoes.tsx` currently exports `ImportacoesRoute`; `AppRouter.tsx` registers
  `/importacoes` OUTSIDE the `<InstallationGatedRoutes />` block.

## What to build

**1. `useErpImportChain.ts`** — a page-local query hook over `client.getErpImportChain(importId)`, following
the `useErpImports.ts` idiom (own exported query-keys object, `useClient()` from `../../app/ClientContext`,
`enabled: Boolean(importId)`, `staleTime: QUERY_STALE_TIME.listings`).

One deliberate deviation from the app default, and it needs a short comment saying why: do NOT retry on a 4xx.
"This import does not exist" is a settled answer, not a flake, and the default single retry only delays the
honest error by a round trip. A 5xx may still use the default single retry.

**2. `ImportChainPanel.tsx`** — props: `{ importId: string }`. Renders a card (`rounded-card border
border-border bg-surface p-4`) titled `Cadeia da importação`, with a one-line subtitle naming what the reader
is looking at (`importados → vinculados → enfileirados`, read from the server).

- Pending → `<LoadingState />`.
- Error, or settled with no data → an ERROR state. Use `<ErrorState onRetry={...} detail={...} />` with
  `detail` = `"Importação não encontrada."` when the error is the 404/`import_not_found` case, and an honest
  generic sentence otherwise. Wrap it in an element carrying `data-testid="erp-import-chain-error"`.
  It must be IMPOSSIBLE for the failure path to render counters — an import whose chain could not be read
  must not look like an import that produced nothing.
- Success → the protocol from the payload, the three counters, and `queue_read_at` formatted with
  `formatDateTime`. Give the counter container `data-testid="erp-import-chain"` and each counter value
  `data-testid="erp-import-chain-importados" | "...-vinculados" | "...-enfileirados"`.
  Labels in pt-BR: `Produtos do import`, `Vinculados`, `Enfileirados`.
- **Each counter renders `<UnknownValue hint="…" />` when its value is not a finite number** (absent, null,
  `NaN`, a string — anything). Never `0`, never a blank cell. Write ONE small local helper for a counter and
  use it three times; do not repeat the guard inline three times.
- `enfileirados` is a queue depth AT AN INSTANT, not an import-history total — the number falls as the queue
  drains. Render `queue_read_at` next to the counters so a smaller number on the next visit reads as drainage
  instead of data loss. `formatDateTime` returning `null` is itself an unknown → `<UnknownValue />`, not a
  blank.

**3. `ImportacaoDetailPage.tsx`** — the detail screen for one import. Reads `importId` from the route with
`useParams` (react-router-dom, already a dependency). Page shell matches `ImportacoesPage.tsx` (same
`mx-auto flex max-w-5xl flex-col gap-[14px]` section + `<header>` with the `text-[22px] font-bold
tracking-tight text-ink` `<h1>`). Include a `<Link to="/importacoes">` back to the list. Renders
`<ImportChainPanel importId={importId} />`.

If the route somehow produces no `importId`, render an honest message — never call the endpoint with an empty
id and never render a chain shell with dashes as if a real import had unknown numbers.

**4. Route** — add `ImportacaoDetailRoute` to `routes/importacoes.tsx` and register
`<Route path="/importacoes/:importId" element={<ImportacaoDetailRoute />} />` in `AppRouter.tsx`, immediately
after the `/importacoes` route and OUTSIDE the installation gate for the same reason the list is: an ERP
import does not depend on a connected marketplace account.

**5. `ImportacaoSection.tsx`** — each row gets a `<Link to={`/importacoes/${item.import_id}`}>` labelled
`Ver cadeia`, placed next to the existing `Ver detalhes` button and styled to match it. This is the only
navigation into the detail screen, so without it the screen is unreachable.

This makes the component require a router context. `ImportacaoSection` is rendered by `ImportacoesPage`,
`ImportacaoDetailPage`'s siblings and `IntegracoesPage`, all of which mount under `BrowserRouter` in
`AppRouter` — so the app is fine. The TESTS are what change: `ImportacaoSection.test.tsx` and
`ImportacoesPage.test.tsx` (the latter is NOT in your write_set — so if it breaks, STOP and report it rather
than editing it).

## Tests (failing FIRST, then green)

`ImportChainPanel.test.tsx` — mock `../../app/ClientContext` the way `ImportacaoSection.test.tsx` does, with a
`getErpImportChain` spy. Cases:

1. **Consumption is real.** Resolve a chain payload whose numbers appear NOWHERE else (e.g. `importados: 137`,
   `vinculados: 42`, `enfileirados: 9`). Assert those three values render AND that `getErpImportChain` was
   called with the import id. Numbers that could have been derived from a list payload would not prove
   anything — these cannot.
2. **A missing counter is `—`, never `0`.** Resolve a payload with `vinculados` ABSENT (build it as a partial
   object and cast). Assert the `vinculados` testid has text content `—` and assert explicitly that it is NOT
   `0`. Assert the other two counters still render their real values — one unknown field must not blank the
   whole card.
3. **A null counter is `—`, never `0`.** Same as (2) with `enfileirados: null`.
4. **404 is honest.** Reject with `{ status: 404, error: "import_not_found" }`. Assert the error testid is in
   the document with the "não encontrada" wording, and assert `queryByTestId("erp-import-chain")` is `null` —
   the failure must not render as an empty/zero chain.
5. **A 5xx is honest too.** Reject with `{ status: 500, error: "internal_error" }` and assert an error state
   renders (different wording from the 404 case, and still no chain).

Use `new QueryClient({ defaultOptions: { queries: { retry: false } } })` in the tests so a rejection settles
immediately.

`ImportacaoDetailPage.test.tsx` — render the page inside `MemoryRouter` with `initialEntries` pointing at
`/importacoes/imp_1` and a matching `<Routes>/<Route path="/importacoes/:importId">`, and assert the chain
panel renders for that id (i.e. `getErpImportChain` called with `"imp_1"`).

`ImportacaoSection.test.tsx` — wrap the existing `renderSection()` helper in `MemoryRouter` so the new `Link`
mounts. **All five existing assertions must survive unchanged** — a test weakened during a change is a
regression wearing a refactor's clothes. Add ONE assertion: the row exposes a link to
`/importacoes/{import_id}`. In your report, state exactly what you changed in this file and why.

## G-questions to answer in your report

- G1: is a page-local query hook right for the WHOLE system here, or does this belong in
  `packages/web-query`? (Look at what already lives in each and answer from the code, not from taste.)
- G2: one to three lines on anything the card left open that you decided.
- G3: does anything here block a named upcoming seam?
