# Lane: frontend

Method: `C:\Users\leandro.theodoro\Documents\MetalDocs\docs\engineering\repo-audit-playbook.md` (cited by path, not vendored).

## The goal, in the operator's words

> I have run something really deep into MetalDocs and I am changing the way I code there to move to
> something more professional towards issues, PRs, PR review, CodeRabbit mechanical full validation
> and so much more. For that I had to identify every error in my code, my platform, to improve it and
> to create this full validation. I want to run it here as well so we move on the same path, this way
> it gets so much harder to send bad PRs.

Calibration, chosen explicitly by the operator: **solid professional level.** Not Google-tier.
That means: clear dependencies, clear modules, clear consumable surfaces, no hand-maintained
redundancy, following rules that have existed for decades, optimised for maintainability and future
scaling.

The operative sentence is the last one — **"this way it gets so much harder to send bad PRs."** The
program's success condition is a *mechanism* that makes a bad change hard to land, not a cleaner
codebase. A remediation that fixes code and leaves the judging to discipline has failed this goal
even if every finding is closed.

## Findings

| ID | class | finding | evidence | scale |
|---|---|---|---|---|
| F1 | gap | No automated mechanism connects `contracts/api/marketplace-central.openapi.yaml` to `packages/sdk-runtime`'s hand-written types. The only cross-reference is a code **comment**, not a compiler/codegen link. | `packages/sdk-runtime/src/index.ts:1747` (comment citation only); `grep -rln "openapi.yaml" apps/web/src packages` → only comment hits; no `zod`/`io-ts`/`ajv`/`joi`/`yup` in any FE `package.json` | 1 spec file, 1 hand-written SDK, 0 agreement checks |
| F2 | gap/cicd | `tsc --noEmit` is not wired into any npm script, `harness.ps1` command, or CI workflow. Run directly, it fails with 12 real errors, 3 in shipped (non-test) code. `npm run harness:unit` (the wired FE gate) only runs vitest, whose esbuild transform does not type-check. | `cd apps/web && npx --no-install tsc --noEmit -p tsconfig.json` → 12 `error TS...`; `.github/workflows/` has only `release-images.yml` (no test/lint/typecheck step); `scripts/harness.ps1:56-66` `Invoke-Unit` calls only `go test` + `npm run test --workspace web` | 12 tsc errors repo-wide (3 production, 9 test-double drift), 0 wired |
| F3 | hazard | Two live, routed components call `<ErrorState detail="..." />` without the required `onRetry` prop — a real `TS2741` compile error proving a runtime `TypeError: onRetry is not a function` the moment a user clicks "Tentar novamente" after either screen errors. | `apps/web/src/pages/mutations/MutationPreviewModal.tsx:210`; `apps/web/src/pages/mutations/MutationResultSummary.tsx:22`; `ErrorStateProps.onRetry` is required at `packages/ui/src/ErrorState.tsx:4` | 2 call sites, both reachable from `/protocolos/:protocolId` |
| F4 | hazard | Zero React error boundaries anywhere in the app. Any uncaught render exception (including F3) crashes the whole SPA to a blank/broken screen with no recovery UI. | `grep -rn "ErrorBoundary" apps/web/src packages` → 0 matches | 0 boundaries across ~35k LOC |
| F5 | duplication | Four independent implementations of the same env/base-URL resolution logic (`apiBaseUrl`/`resolveApiBaseUrl`); three are near byte-identical hand copies. None of the three market-client copies pass the app's token-refresh `fetchImpl` that `useClient()` wires in. | `apps/web/src/app/ClientContext.tsx:9` (canonical, tested); `apps/web/src/pages/precos/MarketComparison.tsx:21`; `apps/web/src/pages/mercado/marketClient.ts:12`; `apps/web/src/pages/produto/marketClient.ts:11` | 4 files; 3 untested as standalone functions |
| F6 | duplication/drift | Money/percent formatting reimplemented outside the canonical `packages/ui/src/money.ts`, whose own header comment documents this exact defect as already fixed once. One reimplementation lives **inside the same package** as the canonical formatter. | `packages/ui/src/money.ts:1-6` (comment: "One money formatter for every screen... Before this existed each screen chose..."); duplicates at `apps/web/src/pages/mercado/mercadoFormatters.ts:16`, `apps/web/src/pages/pedidos/pedidosFormatters.ts:71,106`, `packages/ui/src/ProductPicker.tsx:44` | 3 files / 4 functions reimplemented |
| F7 | gap | Shared `ErrorState` + typed `ApiErrorCode`/`isApiError`/`hasCode` exist, but most call sites don't use them: of 41 `<ErrorState>` usages, 18 pass no `detail` (generic fallback), 20 pass a hardcoded per-screen string, only ~8 derive `detail` from the actual typed error. A rate-limit, an auth failure, and a network error render the same sentence on ~80% of screens. | `grep -rn "<ErrorState" apps/web/src packages --include=*.tsx` → 41; `detail="` → 20; no `detail` prop → 18; default fallback at `packages/ui/src/ErrorState.tsx:11` ("Erro ao carregar.") | 33/41 (~80%) of error surfaces discard the backend's discriminated error code |
| F8 | gap | Two routed, live screens (`/classifications`, `/estoque`'s `ClassificationsPage`) have never had their test files executed by any wired command — `vitest.config.ts`'s `include` name-pins one file (matches ledger B-7) and omits `feature-classifications`/`feature-inventory` entirely. `ClassificationsPage.tsx` is also the only non-test FE file using `catch (err: any)` (7×) and the only routed page doing data-fetching with no react-query at all. | `apps/web/vitest.config.ts:11-17`; live run `npm run test --workspace @marketplace-central/web -- --run` → 72 files/605 tests, none from `feature-classifications`/`feature-inventory`; `apps/web/src/app/AppRouter.tsx:2-3,73,83` (both routed); `packages/feature-classifications/src/ClassificationsPage.tsx:103,180,208,228,267,303,326` (7× `catch (err: any)`) | 2 packages/screens untested; 7 `any` sites concentrated in 1 file |
| F9 | idiom/duplication | Three coexisting, incompatible data-fetching patterns: (a) direct react-query render off `query.data/isLoading/isError` — majority, 42 files; (b) react-query **plus** a hand-mirrored parallel `useState`/`useEffect` copy of the same data (`StockSeguroPage.tsx`); (c) no react-query at all — raw `useState`+`useEffect`+`Promise.all` (`ClassificationsPage.tsx`, and dead `DashboardPage.tsx`). | `grep -rl "@tanstack/react-query"` → 42 files; `packages/feature-inventory/src/StockSeguroPage.tsx:182-219`; `packages/feature-classifications/src/ClassificationsPage.tsx:64-108` | 3 mechanisms, 2 outlier files |
| F10 | idiom | A real defect class (ledger D-11: boolean-ternary chain over react-query state doesn't partition all reachable states) was fixed once, excellently, with an exhaustive-by-construction `query.status` switch — but applied to exactly one file. 9 other files still use the vulnerable `isLoading ? … : isError ? …` pattern that produced the original bug. | Fix: `apps/web/src/pages/integracoes/SyncHealthCard.tsx:127-190`; vulnerable pattern still present: `grep -rln "isLoading ? " apps/web/src packages --include=*.tsx` minus SyncHealthCard → 9 files | 1 fixed / 9 unmigrated |
| F11 | dead code | `apps/web/src/pages/DashboardPage.tsx` (160 lines) — a complete, compiling, second Dashboard with zero importers anywhere; the live route uses `pages/dashboard/DashboardPage.tsx` instead. Still calls two live SDK methods and reimplements its own loading state machine and its own client-side margin average. | `apps/web/src/routes/dashboard.tsx:1` imports only `../pages/dashboard/DashboardPage`; exact-path grep for `pages/DashboardPage` importers → 0 | 1 file, 160 lines, 2 live SDK calls kept warm for nothing |
| F12 | idiom | The npm workspace glob (`apps/web` + `packages/*`) resolves 3 packages with **zero code**: `feature-connectors` (nothing tracked), `feature-orders` (nothing tracked), `feature-simulator` (package.json only, no src). `apps/landing` is a single `index.html`, no build wiring. `design-system/` is markdown documentation only. | `git ls-files packages/feature-connectors packages/feature-orders` → empty; `git ls-files packages/feature-simulator` → `package.json` only; `git ls-files apps/landing` → `index.html` only | 3 empty package shells + 1 non-functional app shell in the workspace inventory |
| F13 | hazard/idiom | The `sdk-runtime` module count established as fact 5 (2,595 lines / 172 interfaces) is the `index.ts` file alone. Four more hand-written, hand-maintained modules sit beside it, outside that count, each documented as "deliberately independent" of `index.ts` (own type unions, own client factory, own `apiBaseUrl` duplication per F5). | `packages/sdk-runtime/src/{activeSource,dashboard,erpImport,market}.ts` = 42+20+99+193 = 354 non-test lines; `market.ts:174` comment: "D-F4-o keeps this file independent of ./index.ts" | +354 lines / +2 client factories not covered by established fact 5's headline number |
| F14 | idiom (minor) | One instance of genuine client-side money aggregation outside the backend: DIFAL totals for the currently-loaded order page are summed in the browser rather than read from a backend-computed total. | `apps/web/src/pages/pedidos/PedidosPage.tsx:185` `difalPending.reduce((sum, order) => sum + (order.difal.amount ?? 0), 0)` | 1 site (plus the same pattern, unreachable, in dead F11) |
| F15 | idiom (positive-leaning, sized) | Quick accessibility pass: 136 `aria-*` attributes, 59 `role=` usages, 0 `<img>` tags (icon-based UI, no missing-alt risk), 0 non-semantic `div`/`span` `onClick` handlers. Not exhaustive — no automated axe run, no contrast/keyboard-trap check performed. | `grep -rno "aria-[a-z]*=" … \| wc -l` → 136; `role=` → 59; `<div[^>]*onClick` → 0 | Reasonable baseline for ~24k LOC; unverified beyond grep |

## The five heaviest, with detail

**1. F3 + F4 — a confirmed, reachable runtime crash with no safety net (hazard).**
`MutationPreviewModal.tsx:210` and `MutationResultSummary.tsx:22` render
`<ErrorState detail="..." />` without the `onRetry` prop that `ErrorStateProps` requires
(`packages/ui/src/ErrorState.tsx:3-6`). This is not a hypothetical: `tsc --noEmit` reports it as
`TS2741`. At runtime, React does not type-check JSX — the component mounts, and only breaks the
instant the user clicks "Tentar novamente" (`ErrorState.tsx:12-14` calls `onClick={onRetry}` on
`undefined`). Because `grep -rn "ErrorBoundary"` returns nothing anywhere in the app, that
`TypeError` is not contained to the widget — it propagates to React's default top-level unhandled
error behavior and takes the whole SPA down. Both call sites are reachable from
`/protocolos/:protocolId` (`AppRouter.tsx:72`), a route that exists specifically to track a mutation
that is, by definition, likely to fail sometimes. Neither `npm run harness:unit` nor any CI workflow
would have caught this — see F2.

**2. F2 — the type-checker is configured correctly and never runs (gap → mechanism failure).**
`tsconfig.base.json:6` sets `"strict": true`. That part of "professional level" is already done.
But nothing in `scripts/harness.ps1`, `package.json`, or `.github/workflows/` (which has exactly one
workflow, for release images, no test/lint/typecheck step at all) ever invokes `tsc --noEmit`.
Running it directly surfaces 12 errors, 3 of them the F3 production bug, the rest test-double drift
(mocks no longer matching updated SDK interfaces — itself a smaller hand-sync signal). This is the
purest instance of the audit's own operative sentence: the tool that would have caught F3 exists,
is correctly configured, and is disconnected from every path a PR can take to landing. A discipline
rule ("remember to run tsc") is not a mechanism; per the brief's own framing this must be recorded as
exactly that kind of gap.

**3. F1 — the contract seam, measured from the side that consumes it.**
The established facts (5-7) describe the seam's shape: 2,595 hand-written lines, 172 interfaces,
4 independent hand-transcriptions from domain to DTO to OpenAPI to SDK, and a governance rule that
only checks same-commit change, never agreement. From the frontend's side the picture is not better:
there is no codegen (`oapi-codegen`/`openapi-typescript` approved 2026-08-07, not installed — fact
7), no runtime schema validation library in any FE `package.json`, and the only place the OpenAPI
spec is even mentioned in FE source is a single code comment
(`packages/sdk-runtime/src/index.ts:1747`) that cites it as documentation, not as an input to
anything that runs. **Direct answer to the brief's question: there is no mechanism.** If the Go
backend renames or retypes a response field without a correct, manual, simultaneous edit to
`sdk-runtime`, nothing fails at build time (TypeScript trusts the hand-written type unconditionally)
and nothing fails at runtime unless a dereferenced field happens to be `undefined` somewhere that
throws. The frontend "learns" of a contract change only when a human notices wrong data on a screen.

**4. F6 — the money-formatter was unified once, in writing, and the unification didn't hold.**
`packages/ui/src/money.ts:1-6` is unusually explicit about the defect it fixes: it names the exact
symptom ("the same value rendered as 'R$ 53,90' on one screen and 'R$ 53.9' on the next") and states
"One money formatter for every screen." That comment is contradicted by measurement: three more
`formatMoney`/`formatCurrency`/`formatPercent` functions exist outside it
(`mercadoFormatters.ts:16`, `pedidosFormatters.ts:71,106`), and the most telling one is inside
`packages/ui` itself — `ProductPicker.tsx:44` defines its own `formatCurrency` using
`toLocaleString` directly, one file away from the canonical `Intl.NumberFormat`-based
implementation it never imports. This is `drift`, not just `duplication`: a fix that a comment
declares complete, disproven by grep, inside the very module meant to be the single source.

**5. F8 — the one file combining every risk this lane measures.**
`packages/feature-classifications/src/ClassificationsPage.tsx`, routed live at `/classifications`
(`AppRouter.tsx:73`), is simultaneously: the only non-test FE source file using `catch (err: any)`
(7 occurrences, `err?.error?.message` string-fished instead of the typed `isApiError`/`hasCode`
contract used elsewhere); the only routed screen doing data-fetching with no react-query at all
(raw `useState`/`useEffect`/`Promise.all`, no caching, no invalidation, no retry); and — confirmed
by actually running `npm run harness:unit` and reading its output (72 files, 605 tests) — a screen
whose own test file has never executed under any wired command, because `vitest.config.ts`'s
`include` list omits its package entirely. This is exactly the kind of file the operator's stated
goal ("harder to send bad PRs") needs to catch, and today nothing does.

## Real FE inventory (what counts as the frontend, and how established)

The frontend is **not** `apps/web/src` alone. Established via `git ls-files` per directory
(counted with `cat | wc -l` to avoid an `xargs`/`wc -l` batching artifact that silently truncates
totals on file lists over ~150 entries — a instrument bug worth flagging on its own: naive
`xargs wc -l | tail -1` undercounted by roughly a third on the first attempt):

| Path | Role | Non-test+test LOC (ts/tsx) | Notes |
|---|---|---|---|
| `apps/web/src` | App shell, routing, 9 page groups | 24,299 | The "app" in the narrow sense |
| `packages/sdk-runtime/src` | Hand-written API client + types | 5,990 | `index.ts` (2,595, established fact 5) + `index.test.ts` (2,402) + 4 more hand-written modules (F13) |
| `packages/ui/src` | Shared component library | 1,854 | `ErrorState`, `DataTable`, `PaginatedTable`, `money.ts`, etc. |
| `packages/web-query/src` | Query-key/invalidation/failure-copy helpers | 714 | Thin, react-query-adjacent utilities |
| `packages/feature-classifications/src` | Routed feature package (`/classifications`) | 862 | See F8 |
| `packages/feature-inventory/src` | Routed feature package (`/estoque`) | 888 | Untested test file, see F8 |
| `packages/feature-products/src` | Routed feature package (`/catalogo`) | 632 | Correctly wired (exact-file-pinned per ledger B-7) |
| `packages/feature-connectors` | Workspace member | 0 | Nothing tracked at all, not even `package.json` |
| `packages/feature-orders` | Workspace member | 0 | Nothing tracked |
| `packages/feature-simulator` | Workspace member | 0 code | `package.json` only |
| `apps/landing` | Separate app | 0 | Single `index.html`, no build wiring, no `package.json` |
| `design-system/` | Design documentation | 0 code | 5 markdown files, cited by comments elsewhere, not imported by any build |

Sum of code-bearing directories ≈ 35,239 lines, which reconciles with PHASE-0's repo-wide established
measurement of 35,423 (`git ls-files '*.ts' '*.tsx'`) — the small remainder is `apps/web`'s
top-level config/entry files. **Conclusion: the established repo-wide LOC total already captures the
real frontend** (nothing large is hiding outside `git ls-files '*.ts' '*.tsx'`); what was missing
from a naive read was the *shape* — three of nine workspace packages are empty shells (F12), and the
contract surface (`sdk-runtime`) is wider than its headline file by 354 lines across four more
independently-documented "deliberately separate" modules (F13).

## The contract seam from the FE side

- **SDK types consumed:** 74 non-test files import from `@marketplace-central/sdk-runtime`
  (`grep -rl "@marketplace-central/sdk-runtime" apps/web/src packages --include=*.ts --include=*.tsx
  | grep -v .test. | wc -l`). `index.ts` exports 220 top-level `interface`/`type` declarations by
  direct grep count (established fact 5 cites "172 interfaces" — the difference is plausibly
  interfaces-only vs interfaces+type-aliases; not reconciled, flagged `unverified` rather than
  contradicting the established fact).
- **SDK types re-declared locally instead of imported:** 125 local `interface`/`type` declarations
  exist outside `sdk-runtime` across `apps/web/src` + `packages/{feature-*,ui,web-query}`. Spot
  review of the largest concentrations (`tariffBadge.tsx`, `useVinculosQueue.ts`, `MercadoPage.tsx`,
  `anunciosQueryState.ts`, `StockSeguroPage.tsx`) shows these are legitimate FE-local UI state
  (tab unions, filter state, component props) — **not** shadow copies of backend DTOs. This is a
  negative finding worth stating plainly: unlike the SDK's own 4-copy problem (established facts
  5-6), the app layer is not independently re-deriving backend shapes.
- **Request bodies built by hand instead of through the SDK:** none found. Every HTTP call in the
  app funnels through the SDK's client factories; a repo-wide grep for `fetch(` outside
  `sdk-runtime/src` returned 0 hits (careful to distinguish from `fetchImpl(`, which is the SDK's
  internal choke point and appears only inside `sdk-runtime`).
- **How the frontend learns a backend contract changed: there is no mechanism.** See F1/heaviest #3
  above. The governance rule that exists (`GOV_API_SDK_SPLIT`, established fact 6) is same-commit
  only. No codegen, no schema validation, no contract test comparing SDK types to the OpenAPI spec or
  to a live/fixture backend response was found. The one test named for this purpose
  (`packages/sdk-runtime/src/errorContract.golden.test.ts`) mocks its own `fetchImpl` with a
  hand-written response envelope — it pins the SDK's *own* error-parsing logic against itself, not
  against the OpenAPI spec or the Go backend.

## What is actually fine

- **`sdk-runtime` type-checks clean in isolation.** `cd packages/sdk-runtime && npx --no-install tsc
  --noEmit` → no output, exit 0. The SDK's own internal consistency is solid; the gap is external
  agreement (F1), not internal quality.
- **No `fetch()` bypass anywhere.** Every HTTP call goes through the SDK's `fetchImpl` choke point
  (see above). This is a real, load-bearing invariant worth protecting as-is.
- **`any` is nearly absent from real code.** `: any` / `as any`: 0 hits in `apps/web/src` non-test
  files; the only non-test hit across all of `packages/` is the 7 sites in
  `ClassificationsPage.tsx` already captured under F8. `@ts-ignore`/`@ts-expect-error` appear only
  in test files, and where inspected (`errorContract.golden.test.ts`) are used correctly to pin a
  type-level contract, not to silence a real error.
- **Shared component library is real and widely adopted, not aspirational.** `ErrorState` has 41
  call sites, `PaginatedTable`/`DataTable` have 5+, `money.ts`'s `formatMoney`/`formatMoneyOr` are
  the majority pattern despite F6's exceptions. The unification effort exists and mostly works; the
  gap is depth of adoption, not absence of the effort.
- **Honest-unknown discipline (ADR-17) shows up consistently in FE code, not just backend.** Every
  money/percent formatter reviewed (`money.ts`, `mercadoFormatters.ts`, `pedidosFormatters.ts`)
  explicitly returns `null`/em-dash rather than fabricating `0` for a missing value, each with a
  comment citing the rule.
- **`tsconfig.base.json` has `"strict": true`.** The type system is configured to a professional
  standard; the failure mode (F2) is that nothing runs it, not that it's configured loosely.
- **D-11's actual fix (`SyncHealthCard.tsx:127-190`) is a genuinely good pattern**, worth using as
  the template when F10 is addressed: switching the discriminant to react-query's closed
  `QueryStatus` union makes the render exhaustive by TypeScript construction rather than by manual
  ternary bookkeeping, and the comment block cites the exact live-drive evidence that motivated it.
- **No non-semantic click handlers, no missing alt text** (F15) — a reasonable accessibility
  baseline for the codebase's size, though not exhaustively verified.

## Unverified / needs judgment

- Whether `oapi-codegen`/`openapi-typescript` (approved 2026-08-07 per established fact 7) would
  close the F1 gap without further hand-work — not attempted; out of scope for a read-only lane.
- The 172-vs-220 interface-count discrepancy against established fact 5 — not reconciled.
- Whether the 9 test-file `tsc` errors (mock shape drift, part of F2's 12) indicate stale test
  doubles that happen to still pass at runtime (because vitest's esbuild transform never
  type-checks) or an actual coverage gap — observed as compile errors only, not chased into runtime
  behavior.
- Whether other IC-03/market call sites beyond the 3 named in F5 also bypass the app's token-refresh
  wrapper — only those 3 were checked directly; not an exhaustive sweep of all `fetchImpl` call
  sites.
- Accessibility beyond the grep-level pass in F15: no contrast check, no keyboard-trap testing, no
  automated axe/lighthouse run performed.

## Commands run

```
git ls-files '*.go' | grep -v '_test.go$' | xargs wc -l | tail -1   (established, not re-run)
git ls-files apps/landing ; git ls-files design-system
git ls-files packages/feature-connectors packages/feature-orders packages/feature-simulator
git ls-files packages/sdk-runtime | xargs -I{} wc -l {} | sort -rn
for p in sdk-runtime ui web-query feature-classifications feature-inventory feature-products; do
  git ls-files "packages/$p/src/*.ts" "packages/$p/src/*.tsx" -z | xargs -0 cat | wc -l
done
git ls-files 'apps/web/src/**/*.ts' 'apps/web/src/**/*.tsx' -z | xargs -0 cat | wc -l
cat apps/web/package.json apps/web/vitest.config.ts apps/web/tsconfig.json tsconfig.base.json package.json
grep -n "feature-classifications|feature-inventory|..." scripts/harness.ps1  (0 hits)
grep -n "unit|vitest" scripts/harness.ps1 -A 15                              (Invoke-Unit body)
npm run test --workspace @marketplace-central/web -- --run                   (live run: 72 files / 605 tests)
grep -n "feature-classifications|feature-inventory|feature-products|feature-simulator" apps/web/src (AppRouter wiring)
cat apps/web/src/app/AppRouter.tsx
grep -rl "@marketplace-central/sdk-runtime" apps/web/src packages --include=*.ts --include=*.tsx | grep -v .test. | wc -l
grep -c "^export interface|^export type" packages/sdk-runtime/src/index.ts
grep -rEc "^(export )?(interface|type) [A-Z]" apps/web/src packages/{feature-*,web-query,ui}
grep -rn "fetch(" apps/web/src packages --include=*.ts --include=*.tsx | grep -v .test.   (0 outside SDK)
cat apps/web/src/pages/mercado/marketClient.ts apps/web/src/pages/produto/marketClient.ts
grep -n "apiBaseUrl|VITE_API_BASE_URL|baseUrl" apps/web/src/pages/precos/MarketComparison.tsx
grep -rn "VITE_API_BASE_URL" apps/web/src packages --include=*.ts --include=*.tsx | grep -v .test.
cat apps/web/src/app/ClientContext.tsx
Grep ": any\b|as any\b" / "@ts-ignore|@ts-expect-error" across apps/web/src and packages
cat packages/feature-classifications/src/ClassificationsPage.tsx (relevant sections)
grep -rn "^export function format|^function format" ... | grep -iE "money|currency|price|percent"
sed -n across packages/ui/src/money.ts, mercadoFormatters.ts, pedidosFormatters.ts, ProductPicker.tsx
find .github -type f ; grep -n "tsc |test|lint" .github/workflows/release-images.yml
cd apps/web && npx --no-install tsc --noEmit -p tsconfig.json               (12 errors, listed above)
cd packages/sdk-runtime && npx --no-install tsc --noEmit                     (0 errors)
grep -n "onRetry|detail=" around MutationPreviewModal.tsx:210, MutationResultSummary.tsx:1-30
grep -rn "ErrorBoundary" apps/web/src packages                               (0 hits)
grep -n "isLoading ?|SyncHealthCard.tsx status pattern" apps/web/src/pages/integracoes/SyncHealthCard.tsx
grep -rln "isLoading ? " apps/web/src packages --include=*.tsx | grep -v .test.
grep -rn "aria-[a-z]*=|role=|<div[^>]*onClick|<img " apps/web/src packages --include=*.tsx
grep -rn "\* 100|/ 100|reduce((s" apps/web/src packages --include=*.ts --include=*.tsx | grep -v test
```

`git status --porcelain --untracked-files=all` at end of this lane's work shows only this report
file as new/modified.
