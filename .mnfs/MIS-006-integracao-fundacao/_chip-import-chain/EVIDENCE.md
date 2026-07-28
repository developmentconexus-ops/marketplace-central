# EVIDENCE — CHIP-IMPORT-CHAIN (MIS-006 · M-06 · wave 2)

Branch `chip/import-chain`. Base floor `5441fe18f64171ef61cb03b51b5bf66e2922e4eb`.
**Last code commit `aceff011`** — everything after it is evidence only, so the reviewable diff ends
there. Scope: `apps/web` only — screens `/importacoes` (new) and `/integracoes` (verification only).

**Status: P6 gate PENDING (hub's). This document claims no verdict** — see Close. One cold review ran
and returned REPROVA on the close-out; a single pass is not a dual gate.

Companion artifacts, all on disk: `REVIEW-ADVERSARIAL.md` (cold reviewer, verbatim),
`DISPATCH-LEDGER.md` (every dispatch + artifact bytes), `HUB-RULINGS.md` (R1–R5),
`dispatches/` (packs, bindings, cards, prompts, logs, `.last.md`, `.result.json`).

## Commits

| Commit | Author | Subject |
|---|---|---|
| `4b76a287` | dispatch #1 (codex gpt-5.6-luna/high) | `feat(importacoes): add dedicated history route` |
| `1bffdcfd` | dispatch #2 (codex gpt-5.6-luna/high) | `feat(web): wire ERP import chain detail` |
| `67e4a3d5` | **chip session (glue, disclosed)** | `test(importacoes): repair the suites the chain link turned red` |
| `aceff011` | dispatch #3 (codex gpt-5.6-luna/high) | `fix(web): guard import chain protocol` |

`67e4a3d5` is the only code this orchestrating session wrote: ~16 lines across three test files. See I9.

## Verification ladder

| Level | Type | Result |
|---|---|---|
| L0 typecheck | `ran` | `npx --no-install tsc --noEmit -p apps/web/tsconfig.json` → **15 errors**, the exact pre-existing set, **NONE in any file this chip touched**, none outside `apps/web/src/`. Identical across consecutive runs. |
| L1 unit | `ran` | `npx --no-install vitest run` (stock `apps/web/vitest.config.ts`) → **65 files / 521 tests, all passing**. Base was 62 files / 511 tests. |
| L2 live drive | `ran` (**by the HUB, not by this chip**) | Executed 2026-07-28, report committed to `main` at `7099d9f`: `.mnfs/MIS-006-integracao-fundacao/_chip-import-chain/L2-hub-live-drive.md`. Cited by path below per R-14, never recopied. **The waiver's named risk is paid: the wiring stands up.** |
| L3/L4 | n/a | No Go, no migrations, no contract change in this chip. |

Environment: the lanes ran against a real `node_modules` installed by `npm ci` at the worktree root
(lockfile-faithful, 157 packages, exit 0), per `docs/HARNESS-PROFILE.md` §3 and hub ruling R1.
`node_modules/@marketplace-central/*` junctions into `sharp-pike-3387c1`, so both lanes provably
compile THIS worktree's `packages/`.

**Correction of record.** An earlier chip-side measurement used a junction to the main checkout plus a
chip-local vitest config (the pack's instruction, now superseded). The hub ruled that evidence
NON-DISCRIMINATING and it was right: the 15 errors and their composition (QueueRow×2 +
VinculoDrawer×1 + 12 baseline) exist in `main` as well — `INCOMPARABLE` is present in main's
`packages/sdk-runtime/src/index.ts` while main's `QueueRow`/`VinculoDrawer` do not handle it — so
reading the WRONG tree yields the identical number and the identical composition. The observable
passed in both worlds and therefore proved neither. Re-measured after `npm ci`: tsc 15 (same),
vitest 65/518 (same, pre-S3). Two independently built dependency trees landing on identical numbers
is the comparison that does discriminate.

## Who measured what

Every criterion below says which. The distinction is not bookkeeping: **this chip could not run L2 by
construction** (chips never boot a server or bind `:8080`), and the hub's drive touched things no unit
test in this diff can reach.

- **Chip's own lanes (tsc, vitest):** static shape and rendering logic against a mocked `useClient`.
  They prove the component tells the truth about whatever payload it is handed. They cannot prove a
  payload ever arrives.
- **Hub's L2 drive (`7099d9f`, `L2-hub-live-drive.md`):** the running app against a real database,
  with counters cross-checked by **direct SQL, never against the screen that was under test**. It
  proves the request is wired, reaches the server, and returns numbers that match the database.

Neither substitutes for the other, and the two halves of I5 land on opposite sides of that line —
see I5.

## Criteria

### I1 — `/importacoes` route exists; gate placement decided AND justified — PASS (code + live)

`AppRouter.tsx:68-69` registers `/importacoes` and `/importacoes/:importId` **OUTSIDE** the
`<Route element={<InstallationGatedRoutes />}>` block (`:78`).

**The gate was WRONG for this screen, and this is a declared correction, not a silent change.** The
comment at `AppRouter.tsx:63-66` already stated the rule before this chip existed:

> Setup and ERP-side screens must render with no marketplace account connected: /integracoes is
> where the account is connected, and the catalog, stock and import (/importacoes) screens read the
> ERP mirror, which exists before any marketplace does.

The import history contradicted that comment only because it lived inside `/vinculos`, which IS
gated — an accident of location, not a decision. The chip extended the existing comment to name
`/importacoes` rather than duplicating it.

Justification verified against what the screens actually read, not against intent: `useErpImportsList`
/ `useErpImportDetail` (`pages/vinculos/useErpImports.ts:14-33`) call `client.listErpImports()` /
`client.getErpImport(id)`, and `useErpImportChain.ts:20` calls `client.getErpImportChain(importId)` —
**none takes an `installationId`, and no component under `pages/importacoes/` imports
`InstallationContext`.** The cold reviewer independently confirmed this (attack 5).

**Live (hub, `L2-hub-live-drive.md` §"I1 — a rota e o gate de instalação"):** the gated routes' links
carry `?installation=inst-mercado_livre-…`; `/importacoes`, `/integracoes`, `/catalogo` and `/estoque`
are clean `href`s. That is an observable the code alone could not give — the gate's presence or
absence is visible in the emitted markup, not merely in the JSX nesting. No unit test in this diff
asserts router placement.

### I2 — section promoted without losing function; BOTH render sites declared — PASS

Per hub ruling R3, which amended this criterion after the pack missed a site. `ImportacaoSection` was
rendered in **two** places at base, and each has an explicit destination:

| Site at base | Destination | Why |
|---|---|---|
| `pages/vinculos/VinculosPage.tsx:8` (import) + `:159` (render) | **REMOVED** | `/vinculos` is a gated screen; the import history does not belong behind a marketplace gate. |
| `pages/integracoes/IntegracoesPage.tsx:6` (import) + `:449` (render) | **KEPT** (import path updated to `../importacoes/ImportacaoSection`) | `/integracoes` is the upload screen; the history there is the receipt for the upload the operator just performed. Removing it would be a silent regression. |
| — | **NEW**: `/importacoes` owns the screen and the chain | Dedicated home + nav entry. |

Zero surface regression: the operator who used to find the history on `/vinculos` now finds it via
the settings dropdown (`Header.tsx:178-183`, `<Link to="/importacoes">Importações</Link>`) and still
sees it on `/integracoes`.

**Live (hub, §"I2 — os DOIS sítios de render"):** on `/vinculos` the "Importação" section is GONE; on
`/integracoes` it is PRESENT, showing `#001-E` with Aceitos/Rejeitados/Avisos and both `Ver detalhes`
and `Ver cadeia`. Both halves of the R3 decision were checked on screen, not inferred from the diff —
which matters because the site the pack originally missed is exactly the one a diff-only reading
loses again.

Component and test moved by `git mv` (history preserved). **All 14 base assertions survive
byte-identical; exactly one added.** Verified by string:
`diff <(git show 5441fe18:…/vinculos/ImportacaoSection.test.tsx | grep "expect(") <(git show HEAD:…/importacoes/ImportacaoSection.test.tsx | grep "expect(")`
→ single line added, the `Ver cadeia` href assertion. No assertion weakened.

### I3 — `VinculosPage.tsx` restricted to two lines — PASS

`git diff 5441fe18f64171ef61cb03b51b5bf66e2922e4eb HEAD -- apps/web/src/pages/vinculos/VinculosPage.tsx`
→ exactly two hunks, both pure deletions: the `import { ImportacaoSection }` line and the
`<ImportacaoSection />` line. No third hunk, no rename, no reflow, no comment. `--numstat` reads
`0 2` — **zero added lines**, so the incoming `VinculosDesign.golden.test.tsx` from CHIP-VINC-NEUTRO
(which renders `<VinculosPage/>` and sweeps `container.querySelectorAll("[class]")` for off-theme
colour classes) cannot go red on anything this chip introduced: a diff that adds no line adds no
class. The hub ran the equivalent sweep live and got zero, and is merging in the order that keeps a
failure from pointing at the wrong chip. No collision with
CHIP-VINC-NEUTRO. The cold reviewer re-derived this line-by-line against the base copy and reached
the same conclusion (attack 4).

### I4 — the chain is really consumed — PASS

`ImportChainPanel.tsx` binds the three counters ONLY to `chainQuery.data.importados|vinculados|enfileirados`,
fed by `useErpImportChain` → `client.getErpImportChain(importId)`
(`packages/sdk-runtime/src/index.ts:1901` → `GET /erp/imports/{id}/chain`).

Proven by absence as well as presence:
`grep -rn "listErpImports\|useErpImportsList\|accepted_count\|rejected_count" ImportChainPanel.tsx useErpImportChain.ts ImportacaoDetailPage.tsx`
→ **no match**. The counters cannot be derived from the list payload; there is no `reduce`, no
`filter().length`, no sum anywhere in the chain path. The reviewer's repo-wide grep found no other
call site (attack 2).

The test that proves binding rather than mere fetching uses values that appear nowhere else
(`importados: 137, vinculados: 42, enfileirados: 9`), and tests 2/3 use different values again
(138/8, 139/41), so a hardcoded panel is refuted by the suite as a whole.

**Live (hub, §"Os números na tela" + §"Conferência contra o DB"):** on
`/importacoes/eac3ac9e-…` the screen rendered `Protocolo #001-E · 55 · 0 · 55 · Fila lida em:
28/07/2026, 13:32`, with `GET …/chain → 200 OK` observed in the browser's network tab. That network
observation is what closes I4 against the class of defect the unit tests structurally cannot see: a
panel that recomputed from `listErpImports` would render numbers without that request. Each counter
was then cross-checked by **direct SQL** — `count(*) erp_import_products` = 55, `product_links`
state=resolved = 0, `jsonb_array_length(cursor->'pending')` = 55 — never against the screen itself.

The data was not fabricated: the hub declined to INSERT rows (that would manufacture the very fact
under test) and instead uploaded the `example-erp.xlsx` fixture through `POST /erp/imports`,
`201` in 0.24s. The chain read therefore traversed the production write path end to end.

### I5 — ADR-17: absent counter renders `—`, NEVER `0` — PASS (with a scope note the hub should read)

Guard: `renderCounter` (`ImportChainPanel.tsx:15-17`) —
`typeof value === "number" && Number.isFinite(value) ? value : <UnknownValue hint={hint} />`.
Adversarially enumerated by the reviewer: `undefined`, `null`, `NaN`, `±Infinity`, `"0"`, `"42"` all
render `—`; a real `0` renders `0` (correct — the payload carried it); `-1` renders `-1` (visible
rather than laundered).

Tests assert BOTH halves — the em dash AND explicitly `not.toHaveTextContent("0")`
(`ImportChainPanel.test.tsx`, absent-field and null-field cases). Asserting the dash alone would pass
against a blank cell, which is why the not-zero assertion is the one that matters.

The failure path cannot fabricate a chain: `ImportChainPanel.tsx:32` short-circuits on
`isError || !data` and returns only the error div; the chain block is the `else`. There is no state
in which counters and an error coexist. `ImportacaoDetailPage.tsx:20-26` likewise renders an honest
message for a missing `importId` rather than a dash-filled shell.

`protocol` was rendered raw at first — a `200` without it would have shown blank space while its
three siblings showed `—`. Found by the cold reviewer, fixed in `aceff011` via `renderProtocol`
(unknown = not a non-empty string, so empty/whitespace-only also renders `—`; a real protocol renders
byte-for-byte with no output-side normalisation). `queue_read_at` was already guarded through
`formatDateTime(...) ?? <UnknownValue />`.

**Scope note, and it is the one place where the two evidence halves do not overlap.** The OpenAPI
declares all five fields `required` and non-nullable
(`contracts/api/marketplace-central.openapi.yaml:8078`) and the server always emits integers, so
against a CONFORMING server the `—` path cannot fire. The guard defends against version drift between
server and client — precisely when honesty is worth the most — which means **the absent-field half of
I5 is provable only by unit test, by construction. The hub's live route cannot exercise it and does
not claim to.**

What the live drive DID exercise is the discriminant's opposite side, which is the easier one to get
wrong: `vinculados` was a **known** `0`, and the screen rendered `0`, not `—`
(`L2-hub-live-drive.md` §"Os números na tela"). A `—` there would have been exactly as much a lie as
a fabricated zero in the absent case — it would tell the operator "we do not know" when the server
knew and said zero. ADR-17 is a two-sided rule and both sides now have evidence, from different
sources: absent→`—` from this chip's tests, known-zero→`0` from the hub's drive.

### I6 — a failed chain read is honest — PASS

`ImportChainPanel.tsx:32-35` renders `data-testid="erp-import-chain-error"` and nothing else.
404/`import_not_found` → "Importação não encontrada."; anything else → "Não foi possível carregar a
cadeia da importação." The 404 discrimination keys on the shape the SDK actually throws
(`packages/sdk-runtime/src/index.ts:1715` throws `{ status, error }`), so the branch is real in
production and not test-only. Tests assert `queryByTestId("erp-import-chain")` is `null` in both the
404 and 5xx cases — an import whose chain could not be read must not look like an import that
produced nothing.

`useErpImportChain.ts:23-28` deliberately deviates from the app default `retry: 1`: a 4xx is a
settled answer and is not retried; a 5xx keeps the single retry. Evidence the predicate works: the
404 test settles immediately with no timeout, while the 5xx test needs ~1s for the retry.

**Live (hub, §"Id inexistente"):** `/importacoes/00000000-0000-0000-0000-000000000000` with
`GET …/chain → 404` rendered "Erro ao carregar. Importação não encontrada." plus "Tentar novamente",
and **zero counters on screen** — no chain of zeros. This is the criterion where a mocked rejection
and a real 404 could most plausibly diverge (the SDK's thrown shape has to survive the real transport
for the branch to pick the right message), and it did not diverge.

### I7 — F-02 VERIFIED, not rebuilt — PASS

Obligation was to verify and declare. **Nothing was rebuilt; not one line of the active-source
feature was written by this chip.** Verification by string, on the path that produces the behavior:

- `packages/web-query/src/activeSource.ts` — `useActiveSourceQuery` (`:26`), `useSetActiveSourceMutation`
  (`:40`), exported at `packages/web-query/src/index.ts:174-175`.
- `apps/web/src/pages/integracoes/IntegracoesPage.tsx:3` imports both; `:297-298` consumes them.
- Landed in commit **`8048f98b`** — `feat(web): make "Fonte ativa" the tenant's real source instead of
  a browser preference` (`git log -1 -- packages/web-query/src/activeSource.ts`).
- The only surviving `localStorage` string in `IntegracoesPage.tsx` is at `:276`, inside a comment
  explaining the browser-preference bug that was removed — not a code path.

**Contract drift REPORT (I7).** `milestone.md` cites `GET/PUT /tenants/{tenant_id}/active-source`,
which is NOT the landed endpoint. The landed path is `/config/active-source`
(`contracts/api/marketplace-central.openapi.yaml:3290`). Stale references:
`M-06-telas-sdk/milestone.md:137`, `:188`, `:298` (`:298` already flags it) and
`M-06-telas-sdk/validation-contract.md:220`. Also `M-02-mirror-port-active-source/milestone.md:172`.
Not corrected here — mission docs are the hub's.

**Live (hub, §"I7 / F-02"): this is the test F-02 always promised and no suite in this repo could
provide.** Clicking "Planilha Sankhya" issued `PUT 200`; the subsequent `GET` returned
`xlsx`/`upload_snapshot`; the DB recorded `set_at = 2026-07-28 16:34:00.754941+00`. And the global
invalidation is real, not card-local: `/catalogo` stopped serving the Sankhya mirror (10529 rows) and
began serving the xlsx snapshot — `REF-1001 Example Product 1001`, `— (missing_price)`, cost
`R$ 13,34`. The whole app changed source, which is the actual claim F-02 makes and which no unit test
that mocks the client can make. State was restored to `sankhya` afterwards.

Note the honest-unknown behavior visible in that verification for free: the xlsx snapshot's missing
price rendered `— (missing_price)`, not `R$ 0,00`.

### I8 — F-03 resolved: SATISFIED, no contract change — PASS

`getErpImportChain` already existed and is now consumed; **no OpenAPI or SDK change was needed or
made.** The screen's only other need is the row link's identifier, and `ErpImportSummary.import_id`
already carries it (`packages/sdk-runtime/src/erpImport.ts:32`, used at `ImportacaoSection.tsx:148`).

`ErpImportChain` (`erpImport.ts:48-54`), the client method (`index.ts:1901`) and the OpenAPI schema
(`:8078`) are pre-existing CHIP-ANCHORS-2 work, consumed unmodified. Verified by reconciliation:
zero files under `contracts/**` or `packages/sdk-runtime/**` changed. The grant to touch
`/erp/imports*` + the erp-import SDK module was therefore never exercised.

### I9 — no collateral damage — PASS

`git diff --name-only 5441fe18f64171ef61cb03b51b5bf66e2922e4eb HEAD`, code paths:

```
apps/web/src/app/AppRouter.tsx
apps/web/src/app/Header.tsx
apps/web/src/pages/importacoes/{ImportChainPanel,ImportacaoDetailPage,ImportacaoSection,ImportacoesPage}.tsx (+ .test.tsx)
apps/web/src/pages/importacoes/useErpImportChain.ts
apps/web/src/pages/integracoes/IntegracoesPage.tsx (+ .test.tsx)
apps/web/src/pages/vinculos/VinculosPage.tsx
apps/web/src/routes/importacoes.tsx
```

- **Zero** `apps/server_core/**`, **zero** `migrations/**`, **zero** `contracts/**`, **zero**
  `packages/sdk-runtime/**` (grep of the reconciliation → `ZERO`).
- From `pages/vinculos/`: only `ImportacaoSection.tsx` + `.test.tsx` (deleted by the move) and
  `VinculosPage.tsx` — exactly what I9 permits. `QueueRow.tsx`, `VinculoDrawer.tsx`, `QueueTab.tsx`,
  `ResolvidosTab.tsx`, `useVinculosResolved.ts` untouched.
- tsc: 15 before, 15 after, none in this chip's files.
- vitest: 62 files / 511 tests before → **65 files / 521 tests** after, all green.
- `apps/web/vitest.chip.config.ts`: deleted, and proven never to have entered git —
  `git log --all -- apps/web/vitest.chip.config.ts` empty, `git ls-files` empty. The reviewer flagged
  that if it were tracked this would be a hard FAIL; it never was.
- Both `node_modules` junctions removed before the install; the main checkout was verified intact
  immediately after unlinking (99 / 2 entries).

**Disclosure — chip-authored glue, commit `67e4a3d5`.** S2's new `Ver cadeia` `<Link>` made
`ImportacaoSection` require a router context, turning `ImportacoesPage.test.tsx` and
`IntegracoesPage.test.tsx` red. The app was never affected — both mount under `BrowserRouter`; only
the tests mounted the pages bare. I wrapped both render helpers in `MemoryRouter` (no assertion
changed — reviewer confirmed `IntegracoesPage.test.tsx` holds 31 `expect(` and 11 `it(` in both base
and HEAD) and gave the 5xx assertion the ~1s the hook's own retry actually takes. ~16 lines, three
test files, no production behavior. Judged inside the core §4.2 trivial-glue ceiling and disclosed
rather than folded into a worker's commit.

## Adversarial review

Cold Claude gate-reviewer, read-only, did not write any reviewed code (implementer = codex). Artifact
persisted verbatim at `REVIEW-ADVERSARIAL.md` — not streamed, not summarised.

**Verdict: REPROVA**, on the close-out half, not the code. It passed every code attack (integrity,
real consumption, test honesty, ownership, gate justification, reachability, slop) and blocked on:

| Blocker | Status |
|---|---|
| No `EVIDENCE.md`; I7/I8/I9 are declaration-shaped and nothing was declared | **RESOLVED** — this file. The reviewer was right: I had verified F-02 and never written it down, and unwritten = didn't happen. |
| `apps/web/vitest.chip.config.ts` on disk; could not tell if tracked | **RESOLVED** — deleted; proven never tracked in any commit on any branch, never in the index. |
| L2 live drive not run | **RESOLVED by the hub**, `7099d9f` — see the ladder and "What the L2 drive PAID". Not resolved by this chip, and not resolvable by it. |

One code ressalva accepted and fixed (`protocol` guard, `aceff011`). Two declined as correct-as-is:
the redundant `enabled` + `importId` guard (cheap, and the route param genuinely cannot be empty),
and `routes/importacoes.tsx` as a pass-through (it matches `routes/integracoes.tsx` exactly —
convention, not invention).

## Standing REPORTs

1. **`pages/vinculos/useErpImports.ts` should leave `/vinculos`.** `pages/importacoes/` now imports it
   cross-directory (`ImportacaoSection.tsx:6`), leaving the promoted screen coupled to the directory
   it was promoted out of — an unfinished promotion. Not fixable here: I9 forbids touching
   `pages/vinculos/` beyond the two named files, and a collision with CHIP-VINC-NEUTRO would cost more
   than the coupling. Hub ruling R4 accepted this and took the move onto its own board, **after
   VINC-NEUTRO merges**.
2. **DISCHARGED — the CODPROD suspicion was raised and ruled out.** The drive did show
   `vinculados = 0`, and the hub suspected the leading-zero undercount exactly as this REPORT asked,
   then refuted it with two independent measurements (`L2-hub-live-drive.md` §"O `vinculados = 0` é
   VERDADEIRO"): (a) `count(*) filter (where codprod ~ '^0')` on the import = **0**, so the defect
   CHIP-ANCHORS-3 fixes has nothing to bite on here; (b) the ranges are disjoint — the import's
   CODPROD are `1001..1055` while the 29 resolved links' `internal_product_id` are `15956..42194`.
   The fixture is the M-01 synthetic catalogue, not Sankhya data. The zero is TRUE, and it is not
   this chip's FE nor ANCHORS-3's bug.
3. **`enfileirados` is a queue depth at an instant**, not an import-history total — it falls as the
   queue drains. `Fila lida em` renders beside it so a smaller number on a later visit reads as
   drainage rather than data loss.
4. **Contract drift** in `milestone.md` / `validation-contract.md` — see I7.

## Field findings (candidates for ratification under core §0)

- **FINDING-1 — a mid-install `node_modules` fabricates errors that impersonate real defects.**
  Observed during this chip's own `npm ci`: `TS1005: '}' expected` INSIDE
  `node_modules/@types/node/zlib.d.ts`; `Cannot find module 'node:url'` in `apps/web/vite.config.ts`;
  and an `aria-query` directory holding `lib/` + `LICENSE` but no `package.json`, which took down all
  65 vitest files at setup with what reads like a dependency-version conflict. `node_modules/.bin`
  appearing is NOT an install-complete signal — it lands early, and I wrongly used it as one, drawing
  three wrong conclusions before catching it. **Two identical consecutive tsc runs did not catch it
  either.** The only honest signal is the install process exiting. Written into
  `dispatches/bindings-import-chain-v2.md` so dispatched workers cannot repeat it.
- **FINDING-2 — sandbox vitest blindness has teeth.** `codex exec --sandbox workspace-write` on
  Windows cannot run vitest (esbuild `Access is denied`), so all three workers committed on typecheck
  alone. Dispatch #2 broke two suites it could not see and reported "typecheck clean, committed". The
  chip's post-dispatch vitest re-run caught 2 broken files plus 1 genuine timing defect. That re-run
  is the verification of record, not a formality.
- **FINDING-3 — `New-DispatchPrompt.ps1` does not validate the role string.** A typo
  (`ImplementStandard`) assembles cleanly and fails later at `Invoke-CodexDispatch.ps1` with
  `ROLE-UNKNOWN`. Cheap fix: validate against `roles.psd1` at assembly time.
- **FINDING-4 — Windows PowerShell 5.1 mangles the UTF-8 `HarnessDispatch.psm1`** (em-dash inside a
  string → `Token 're-vendor' inesperado`). `pwsh` 7 required, which the profile already binds —
  worth stating as hard, not preferred.
- **FINDING-5 — WITHDRAWN.** The vitest jest-dom fs-guard failure was a consequence of the junction
  technique, not a real environment defect. It did not recur against the stock config after `npm ci`.
  Per hub ruling R2 it dies with the technique and does NOT enter the profile. Recorded here so the
  withdrawal is on the record rather than the finding quietly disappearing.
- **FINDING-6 — codex binary in use is 0.144.4**, while profile F-ENV-9 recommends 0.145.0-alpha.18
  (not installed). All three dispatches completed exit 0. Noted, not a blocker.

## What the L2 drive PAID

The waiver's named risk is settled and should not be carried forward as open. `GET
/erp/imports/{id}/chain` had never been driven live by anyone; the risk was a composition-root wiring
defect of the M-02 catalog-503 class, invisible to all 521 tests because **every test in this diff
mocks `useClient` directly**. The hub drove it: `200 OK` on the wire, counters matching direct SQL.
**The wiring stands up.** With it, the runtime halves of I1, I2, I4, I6 and I7 are no longer open.

Also settled by the drive, and worth recording because it was on this list as unproven: light/dark
rendering. Measured by computed value rather than screenshot — `light` bg `rgb(251,250,247)` / text
`rgb(37,41,31)`, `dark` bg `rgb(22,24,20)` / text `rgb(233,232,226)`, and **0** literal colour classes
under `main` in both, by a `-(slate|blue|emerald|red|gray|zinc|…)-\d` sweep. Computed value is the
better instrument here anyway: a screenshot shows a colour, not whether it came from a token.

## What is STILL not proven

Stated plainly, because a close that hides this is worth less than one that doesn't.

- **`vinculados` was never exercised non-zero.** The drive rendered a true `0`, cross-checked against
  the DB — so the counter is proven honest, but the path where it carries a real linked count is
  untraversed. Proving it needs an import whose CODPROD match real links, i.e. client data. The hub
  declares this limitation in its own report rather than letting the drive read as fuller coverage
  than it was; I cite the gap rather than inherit an inflated claim.
- **The absent-field half of I5 cannot fire against a conforming server** (see the I5 scope note). It
  is unit-tested and will stay that way until a real version drift occurs — which is the point of the
  guard, not a hole in it.
- The `—` path for `protocol` and `queue_read_at` likewise: guarded and unit-tested, never observed
  live, because the server always sends them.

## Close

**P6 GATE PENDING — hub's side. This chip claims no verdict.**

**Correction of form, on the hub's ruling and it is right.** An earlier revision of this section led
with `AGREEMENT — P6 discharged` while its own next sentence said the `P6-DUAL-GATE:` line belongs to
the hub. Those cannot both be true. `AGREEMENT` is a verdict word — it means two independent sides
converged — and this chip had ONE side, a cold Claude gate-reviewer, which returned **REPROVA**. A
header that promises more than the body delivers is the exact shape this mission's failures keep
taking: CHIP-ANCHORS-2 sent `AGREEMENT` on three Claude passes and zero GPT, and the real gate later
refuted BOTH sides and produced CORR-1..CORR-6 including a blocker no earlier pass had seen;
CHIP-VINC-NEUTRO was returned the same day for the mirror-image defect, two passes both
`gpt-5.6-sol` medium with no Opus side. R-26: verbatim is a claim about FORM, and that header claimed
a verdict that does not exist. Substance below is unchanged — only the claim about it is corrected.

What this chip actually delivers: **ladder complete, gate not run.** L0 `ran` (15 errors, none mine),
L1 `ran` (65 files / 521 tests green), L2 `ran` **by the hub** (`7099d9f`, re-confirmed against this
chip's final HEAD — see below). Adversarial review dispatched, persisted verbatim, verdict REPROVA on
the close-out; all three of its blockers now closed, two by this chip and the third by the hub's
drive. That is a single cold pass, not a dual gate, and it does not become one by being good.

Two things stay open by their nature and are named above rather than buried: `vinculados` non-zero,
and the absent-field half of I5. The dual gate — Opus full + GPT-5.6 Sol medium, independent — is the
hub's to run and to write, and the side that implemented does not review: all four code commits came
from codex `gpt-5.6-luna/high`, so the Opus side must be a cold session, never this one.

### L2 re-confirmed against the final HEAD

The hub's first drive ran against the tree at `67e4a3d`/`aceff011` while this branch was still
moving; the CLOSED event was sent at `42d7c2d1`. The hub re-drove with the worktree confirmed at
`42d7c2d11d7bbd5d1cbceeb45380e5e980220d80` and got byte-identical numbers (`#001-E · 55 · 0 · 55`,
`Fila lida em: 28/07/2026, 13:42`). So `aceff011`'s `protocol` guard regressed nothing the first
drive had measured, and `42d7c2d1` is evidence-only. **L2 holds for the HEAD this chip actually
closed** — which is the claim that matters, and it was not the claim the first report could support.
Recorded because it is the live-drive form of `base_sha` is a FLOOR: a drive has a floor too, and the
hub disclosed that its own was unanchored.

**The `protocol` guard cannot be exercised live, confirmed by the hub against the schema:** the column
is `NOT NULL` with `CHECK (protocol ~ '^#[0-9]{3,}-E$')`, so a conforming server cannot return `200`
without one. The unit test is the only possible proof and the hub accepted it as such. Same shape as
the I5 scope note — declared, not inflated.

### Findings, as ratified by the hub

**F-1 RATIFIED into profile §3** (mid-install `node_modules` fabricates errors that impersonate real
defects) — the hub's note is worth keeping attached: it is the sibling of the 15-error trap, since in
both cases the observable is STABLE and still does not discriminate. **F-2 RATIFIED** (sandbox vitest
blindness; a codex slice reporting "typecheck clean, committed" over a red suite is blind, not
dishonest, and the chip's re-run is the verification of record). **F-3/F-4 ACCEPTED as upstream** —
they belong to `mnfs-harness`, not this repo's profile. **F-5 withdrawal ACCEPTED.** **F-6 noted, not
blocking.**

REPORT 4 (contract drift) was **corrected by the hub at `79b5967`** across all five sites, and it was
worse than reported: `milestone.md:137` did not merely cite the phantom path, it instructed a future
chip to refactor `useActiveErpSource` off `localStorage` — work that had already landed as
`activeSource.ts`. That line would have ordered a rebuild of existing code, which is the failure this
chip's pack had to pre-empt by hand for F-02.

Two items from the hub's report that concern method rather than this chip's code, recorded so they
survive the transcript: an **unreproduced** anomaly where a radio click checked the DOM input without
firing a `PUT` (one occurrence, no repro — real as a class, uncontrolled-input drift, but not scope
and not a defect of this diff), and the hub's own disclosed method error where a coordinate click hit
a card underneath and WROTE the active source, briefly manufacturing a "screen disagrees with server"
that only the DB's `set_at` unmasked. Both are the hub's to hold; neither changes a criterion here.
They belong in the record because a drive that reports only its successes is worth less than one that
reports how it nearly fooled itself.
