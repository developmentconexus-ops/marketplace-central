SLICE CARD — S3 · `protocol` is guarded like every other field on the card

## Why this slice exists (read before writing)

A cold adversarial reviewer read the chain panel and found one honesty asymmetry
(`REVIEW-ADVERSARIAL.md`, attack 1, ressalva (a)):

> `protocol` at `:39` is rendered raw — a `200` with a missing `protocol` yields a blank span,
> not `—`, while its three sibling counters are guarded.

That is the ADR-17 rule applied unevenly INSIDE ONE CARD. The three counters route through
`renderCounter` and render `<UnknownValue />`; `queue_read_at` routes through `formatDateTime` and
falls back to `<UnknownValue />`. `protocol` alone renders whatever arrived, so an absent or empty
protocol shows as blank space — which reads to the operator as "this import has no protocol"
rather than "we did not receive one". A blank is a quieter lie than a zero, not a smaller one.

The OpenAPI declares `protocol` REQUIRED (`contracts/api/marketplace-central.openapi.yaml:8078`),
same as the counters. That is exactly why the counters are guarded anyway: a required field that
arrives absent is version drift between server and client, and drift is when honesty pays. Apply
the same reasoning to `protocol`. Do NOT change the contract or the SDK.

## write_set (nothing else)

- `apps/web/src/pages/importacoes/ImportChainPanel.tsx`
- `apps/web/src/pages/importacoes/ImportChainPanel.test.tsx`

## Contract facts (verified — do not re-derive)

- `ErpImportChain.protocol: string`, declared required and non-nullable
  (`packages/sdk-runtime/src/erpImport.ts:48-54`).
- `UnknownValue` comes from `@marketplace-central/ui` (optional `hint?: string`, renders `—` in
  `text-faint`) and is ALREADY imported in this file — reuse it, do not hand-roll a dash.
- The existing counter guard is `renderCounter` at `ImportChainPanel.tsx:15-17`:
  `typeof value === "number" && Number.isFinite(value) ? value : <UnknownValue hint={hint} />`.
- The current raw render is `ImportChainPanel.tsx:39`:
  `Protocolo <span className="font-mono font-medium text-ink">{chainQuery.data.protocol}</span>`.

## What to build

Route `protocol` through an unknown-guard in the same spirit as `renderCounter`, so the whole card
tells the truth the same way.

A protocol is UNKNOWN when it is not a non-empty string: absent, `null`, `undefined`, not a string,
or a string that is empty/whitespace-only. A whitespace-only protocol renders as blank space on
screen, which is indistinguishable from no protocol at all — treat it as unknown, do not trim-and-
render. Anything else renders verbatim (do NOT normalise, uppercase, or reformat a real protocol —
`#001-E` is an operator-facing identifier and must survive byte-for-byte).

Keep the `font-mono font-medium text-ink` styling for a KNOWN protocol. When unknown, render
`<UnknownValue hint="…" />` with a pt-BR hint consistent with the three existing hints
(`"produtos do import desconhecidos"` etc. at `:45,51,57`) — do not wrap `UnknownValue` in the
`text-ink` span, since it carries its own `text-faint`.

Write ONE small local helper next to `renderCounter` and use it once; do not inline the guard, and
do not generalise `renderCounter` into a polymorphic any-value renderer — a second consumer does
not exist. Two small named helpers that each say what they mean beat one clever one.

## Tests (failing FIRST, then green)

Add to `ImportChainPanel.test.tsx`, matching the existing mock idiom in that file
(`getErpImportChain` spy + `QueryClient` with `retry: false`):

1. **A known protocol renders verbatim.** Payload with `protocol: "#001-E"` → assert `#001-E` is on
   screen. (The existing test at `:29-44` already resolves a `protocol`, so make this assertion
   specific enough not to duplicate it — assert the exact string survives.)
2. **An absent protocol renders `—`, not blank.** Build the payload without `protocol` (partial
   object + cast, as tests at `:46-61` already do). Assert the em dash renders AND assert the three
   counters still render their real values — one unknown field must not blank the rest of the card.
3. **An empty/whitespace protocol renders `—`.** `protocol: "   "`. This is the case a naive
   `?? "—"` fallback would miss, so it is the one that proves the guard is real.

Do NOT weaken or rewrite any of the five existing tests. If one of them breaks, that is a finding —
report it, do not adapt the assertion.

## Verification (run each ONCE, from the worktree root)

The worktree now has a REAL `node_modules` (`npm ci`, lockfile-faithful). There is no chip vitest
config any more — use the stock one.

- Typecheck: `npx --no-install tsc --noEmit -p apps/web/tsconfig.json`
  15 PRE-EXISTING errors are expected and are NOT yours (`pages/anunciosQueries.ts`,
  `pages/anunciosQueryState.test.ts`, `pages/AnunciosTable.test.tsx`,
  `pages/ListingsRefreshControl.test.tsx`, `pages/mutations/*`, `pages/produto/ProdutoPage*.test.tsx`,
  `pages/vinculos/QueueRow.tsx`, `pages/vinculos/VinculoDrawer.tsx`). Only errors in the two files
  you touched are yours. Do not fix the others.
- Unit tests, this file: `npx --no-install vitest run src/pages/importacoes/ImportChainPanel.test.tsx`
  run from `apps/web`.
- Full suite baseline to preserve: **65 files / 518 tests passing**. Your slice should end at 65
  files and 521 tests. If any previously-passing file goes red, STOP and report BLOCKED.

If a command genuinely cannot run in your sandbox, record it as `could-not-run (sandbox)` and say so
— do not spend your one fixup on tooling, and do not report Pass on a command you did not run.

## G-questions to answer in your report

- G1: is a second local helper right for the WHOLE system, or should `renderCounter` generalise?
  Answer from the code (how many consumers exist NOW), not from taste.
- G2: one to three lines on anything the card left open that you decided.
- G3: does anything here block a named upcoming seam?
