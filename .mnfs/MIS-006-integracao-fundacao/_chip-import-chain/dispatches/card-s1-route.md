SLICE CARD — S1 · `/importacoes` gets its own route and the import history moves onto it

## write_set (nothing else)

- `apps/web/src/pages/importacoes/ImportacaoSection.tsx`      (already moved here by `git mv`; needs an import fix)
- `apps/web/src/pages/importacoes/ImportacaoSection.test.tsx`  (already moved here by `git mv`)
- `apps/web/src/pages/importacoes/ImportacoesPage.tsx`         (NEW)
- `apps/web/src/pages/importacoes/ImportacoesPage.test.tsx`    (NEW)
- `apps/web/src/routes/importacoes.tsx`                        (NEW)
- `apps/web/src/app/AppRouter.tsx`
- `apps/web/src/app/Header.tsx`
- `apps/web/src/pages/integracoes/IntegracoesPage.tsx`         (ONE import line — see below)
- `apps/web/src/pages/vinculos/VinculosPage.tsx`               (EXACTLY TWO LINES — see below)

## State of the repo right now (verified — do not re-derive)

- `ImportacaoSection.tsx` and its test were already `git mv`'d from `pages/vinculos/` to
  `pages/importacoes/`. The component is otherwise untouched. Its import on line 5 still reads
  `from "./useErpImports"` and is therefore BROKEN — the hook module did NOT move.
- `apps/web/src/pages/vinculos/useErpImports.ts` STAYS WHERE IT IS. Do not move it, do not edit it, do not
  delete it. It is already imported from `pages/integracoes/` today, so the cross-directory import is an
  established pattern here, not a new one. From `pages/importacoes/` the correct specifier is
  `"../vinculos/useErpImports"`.
- `ImportacaoSection` is currently rendered in TWO places:
  - `apps/web/src/pages/vinculos/VinculosPage.tsx` — import on line 8, `<ImportacaoSection />` on line 159.
  - `apps/web/src/pages/integracoes/IntegracoesPage.tsx` — import on line 6, `<ImportacaoSection />` on line 449.
- `apps/web/src/app/AppRouter.tsx` — `/integracoes` is registered on line 66, OUTSIDE the
  `<InstallationGatedRoutes />` block that starts on line 75. `/vinculos` (line 79) is INSIDE that gate.
- `apps/web/src/app/Header.tsx` — the settings dropdown (`role="group"`, `aria-label="Menu de configurações"`)
  holds `<Link to="/integracoes">Integrações</Link>` at lines 172-177, then `/catalogo`, then `/estoque`.
- `apps/web/src/routes/integracoes.tsx` is the pattern for a route module: it imports the page and exports a
  `XxxRoute()` component that renders it.

## What to build

**1. Fix the moved component's import.** In `pages/importacoes/ImportacaoSection.tsx`, change the
`./useErpImports` specifier to `../vinculos/useErpImports`. Nothing else in that file changes in this slice.

**2. `ImportacoesPage.tsx`** — the new screen. It is a page shell that renders `<ImportacaoSection />`.
Follow the shell used by `IntegracoesPage` (read `apps/web/src/pages/integracoes/IntegracoesPage.tsx` at the
`export function IntegracoesPage()` line for the exact markup): an `aria-labelledby`'d `<section>` with
`className="mx-auto flex max-w-5xl flex-col gap-[14px]"`, a `<header>` with an `<h1 id=...>` styled
`text-[22px] font-bold tracking-tight text-ink` and a `<p className="mt-1 text-sm text-muted">` subtitle.
Title: `Importações`. Subtitle: one sentence saying this is the ERP import history by protocol, read-only.

**3. `routes/importacoes.tsx`** — mirror `routes/integracoes.tsx`: export `ImportacoesRoute()` rendering
`<ImportacoesPage />`.

**4. Register the route in `AppRouter.tsx`** — `<Route path="/importacoes" element={<ImportacoesRoute />} />`
placed with `/integracoes`, i.e. **OUTSIDE** the `<Route element={<InstallationGatedRoutes />}>` block, plus
its import alongside the other route imports.

This placement is DECIDED, not yours to re-open — but it needs one honest comment. The existing comment on
lines 62-65 already states the rule ("Setup and ERP-side screens must render with no marketplace account
connected … the catalog, stock and import screens read the ERP mirror, which exists before any marketplace
does"). The import screen contradicted that comment only because it lived inside `/vinculos`, which is gated.
Extend the existing comment (do not duplicate it) so it names `/importacoes` as one of the screens it covers.

**5. Nav entry in `Header.tsx`** — add `<Link to="/importacoes">Importações</Link>` inside the same settings
dropdown, immediately after the `Integrações` link, with the identical `className`. This link is the answer to
"where did the import history go" for an operator who used to find it on `/vinculos`.

**6. `VinculosPage.tsx` — EXACTLY TWO LINES, both deletions.** Delete the `import { ImportacaoSection } from
"./ImportacaoSection";` line (line 8) and the `<ImportacaoSection />` line (line 159). Touch NOTHING else in
that file — not a rename, not a reflow, not a comment. Another chip is editing `/vinculos` in parallel and any
third hunk in this file is a collision that fails the chip.

**7. `IntegracoesPage.tsx` — ONE LINE.** Change the import specifier on line 6 from
`"../vinculos/ImportacaoSection"` to `"../importacoes/ImportacaoSection"`. The `<ImportacaoSection />` render
on line 449 STAYS — `/integracoes` is the upload screen and the history there is the receipt for the upload the
operator just performed. Removing it would be a regression, not cleanup.

## Tests (failing FIRST, then green)

- `ImportacoesPage.test.tsx` (NEW): render `<ImportacoesPage />` inside a `QueryClientProvider` (use
  `new QueryClient({ defaultOptions: { queries: { retry: false } } })`) with `../../app/ClientContext` mocked
  the same way `ImportacaoSection.test.tsx` mocks it (`vi.mock`, returning `listErpImports`/`getErpImport`
  spies). Assert: the `Importações` heading renders, and an import row from a mocked `listErpImports` payload
  renders on the page. This is what proves the screen actually mounts the history rather than just a title.
- `ImportacaoSection.test.tsx`: it must stay green in its new home. Its five existing assertions must survive
  UNCHANGED — a test "adapted" into weaker coverage during a move is a regression wearing a refactor's clothes.
  Change it ONLY if the move genuinely breaks a mechanism (e.g. a mock specifier path), and say exactly what
  you changed and why in your report.

Do not add a test that asserts on `AppRouter` route wiring in this slice; `apps/web/src/app/AppRouter.test.tsx`
is NOT in your write_set.

## G-questions to answer in your report

- G1: does `/importacoes` sitting outside the installation gate hold for the WHOLE system — i.e. is there any
  read on this screen that needs a connected marketplace account? (Read what `useErpImportsList` /
  `getErpImport` actually call before answering.)
- G2: one to three lines on anything you decided that the card left open.
- G3: does anything here block a named upcoming seam? A detail route `/importacoes/:importId` lands in the NEXT
  slice — leave room for it, but do NOT build it here.
