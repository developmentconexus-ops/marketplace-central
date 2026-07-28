# EVIDENCE — CHIP-IMPORT-CHAIN (MIS-006 · M-06 · wave 2)

Branch `chip/import-chain`. Base floor `5441fe18f64171ef61cb03b51b5bf66e2922e4eb`. HEAD `aceff011`.
Scope: `apps/web` only — screens `/importacoes` (new) and `/integracoes` (verification only).

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
| L2 live drive | `could-not-run (hub-owned)` | REQUEST sent to hub `local_99feb041`. Chips never boot a server or bind `:8080`. **This is the open item — see "What is NOT proven".** |
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

## Criteria

### I1 — `/importacoes` route exists; gate placement decided AND justified — PASS (code) / L2 open (behavior)

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

Behavioral half — the screen mounting with no installation — needs the live drive; no unit test
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

Component and test moved by `git mv` (history preserved). **All 14 base assertions survive
byte-identical; exactly one added.** Verified by string:
`diff <(git show 5441fe18:…/vinculos/ImportacaoSection.test.tsx | grep "expect(") <(git show HEAD:…/importacoes/ImportacaoSection.test.tsx | grep "expect(")`
→ single line added, the `Ver cadeia` href assertion. No assertion weakened.

### I3 — `VinculosPage.tsx` restricted to two lines — PASS

`git diff 5441fe18f64171ef61cb03b51b5bf66e2922e4eb HEAD -- apps/web/src/pages/vinculos/VinculosPage.tsx`
→ exactly two hunks, both pure deletions: the `import { ImportacaoSection }` line and the
`<ImportacaoSection />` line. No third hunk, no rename, no reflow, no comment. No collision with
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

**Scope note, carried up rather than buried:** the OpenAPI declares all five fields `required` and
non-nullable (`contracts/api/marketplace-central.openapi.yaml:8078`) and the server always emits
integers. So against a CONFORMING server the `—` path cannot fire — the guard defends against
version drift between server and client, which is precisely when honesty is worth the most, but it
means **I5 is untested against production reality by construction.** A real `0` reaching this screen
means the BACKEND decided zero; the FE cannot distinguish that from "unknown" and must not try.

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

Runtime behavior of the toggle (persists across reload, invalidates globally) is L2 and belongs to
the hub's drive.

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
| L2 live drive not run | **OPEN** — hub's, by design. |

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
2. **If the live drive shows a suspiciously low `vinculados`**, suspect the CODPROD leading-zero
   undercount that CHIP-ANCHORS-3 is fixing in parallel. Record the observed number; do not attribute
   it to this chip's FE.
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

## What is NOT proven

Stated plainly, because a close that hides this is worth less than one that doesn't.

- **The chain endpoint has never been driven live by anyone.** Everything above proves the FE calls
  the SDK method and renders honestly whatever comes back. What no unit test here can prove is that
  `GET /erp/imports/{id}/chain` is actually WIRED in the composition root — **every test in this diff
  mocks `useClient` directly**, so a lost-decorator / composition-root defect of the M-02 catalog-503
  class would pass all 521 tests and surface only on this screen. That is the named risk the
  operator's 2026-07-28 waiver was granted against, and it comes due at L2.
- The screens mounting with no marketplace installation connected (I1's behavioral half).
- The active-source toggle persisting across reload and invalidating globally (I7's runtime half).
- Whether the `—` path can ever fire against the real server (see I5 scope note).
- Light/dark rendering of the new screens.

## Close

**AGREEMENT — P6 discharged.** Ladder: L0 `ran` (15, none mine), L1 `ran` (65 files / 521 tests
green), L2 `could-not-run (hub-owned)`. Adversarial review dispatched, persisted, and its two
chip-side blockers resolved; its third is L2. The `P6-DUAL-GATE:` line is the hub's to write, not
this chip's.
