# HUB RULINGS — CHIP-IMPORT-CHAIN (received 2026-07-28, hub session local_99feb041)

Filed verbatim in substance because they change what this chip's EVIDENCE must contain.

## R1 — env: profile §3 beats the pack. Run `npm ci`. (CORRECTION OF PACK, hub's error)

`docs/HARNESS-PROFILE.md` §3 is `ratified` (2026-07-16, REQUEST from CHIP-SAT) and literal: fresh
chip worktrees have no `node_modules`; run `npm ci` at the worktree root before web `tsc`/`vitest`
lanes. Lockfile-faithful install = env prep, NOT a dep change — no REQUEST needed. Never reuse
another checkout's `node_modules`: npm workspace symlinks would resolve workspace packages to the
OTHER tree's sources, silently validating the wrong code.

The pack's junction + chip-vitest-config instruction is the PRE-ratification technique, carried into
the pack from memory. The pack is wrong; the profile governs.

**Why my mitigation was rejected — and the hub is right.** I argued the `resolve.alias` pinned
`@marketplace-central/*` to this worktree's `packages/`. The hub's objection is not that the alias
failed; it is that **the evidence I cited does not discriminate**. The 15 tsc errors with the
composition QueueRow×2 + VinculoDrawer×1 + 12 baseline also exist in `main` (they landed with the
CHIP-ANCHORS-2 merge, `dbdcdfb1`). Reading the WRONG tree produces the identical count and the
identical composition. The observable passes in both worlds, so it proves neither. Same lesson as
`npx --no-install tsc` passing vacuously.

Action: junctions unlinked (`rmdir` — link only; main checkout verified intact at 99 / 2 entries),
`npm ci` at the worktree root, baselines redone against the REAL `apps/web/vitest.config.ts`,
`apps/web/vitest.chip.config.ts` deleted. If the two measurements agree, the agreement is itself the
evidence the alias was sound — which is what I could not produce before.

## R2 — the vitest fs-guard finding: HOLD ratification

`Failed to load url .../@testing-library/jest-dom/dist/vitest.mjs` with the file present, because the
junction's realpath resolves outside the vite root. Diagnosis accepted as good, but it is a
CONSEQUENCE of the junction. Re-measure post-`npm ci`: if it disappears, the finding dies with the
technique and does NOT enter the profile; if it survives, it is a real finding and the hub takes it
to ratification under core §0. No dispatch spent on it now.

## R3 — `ImportacaoSection` rendered in TWO places: resolution APPROVED

The pack cited `VinculosPage.tsx:159` and missed `IntegracoesPage.tsx:449` (+ import at `:6`) — hub
scan error. Approved as proposed: **removed from `/vinculos`, KEPT on `/integracoes`, `/importacoes`
owns the chain.** The history on the upload screen is the receipt for the upload the operator just
performed; removing it would be the silent regression I2 exists to prevent. `pages/integracoes/` is
already mine by the matrix — no new grant.

**I2 is amended:** EVIDENCE must declare BOTH sites and the destination of each. One screen vanishing
without a record is a regression; two is worse.

## R4 — `useErpImports.ts` stays put: APPROVED, debt is the hub's

The module is already imported from `pages/integracoes/` (`IntegracoesPage.tsx:7`,
`useErpImportUpload.ts:4`), and I9 forbids touching `pages/vinculos/` beyond the two named files.
`../vinculos/useErpImports` is ugly and correct — a collision with CHIP-VINC-NEUTRO would cost more
than the coupling. REPORT accepted; **the hub moves the module after VINC-NEUTRO merges**, not this
chip, not now.

## R5 — route outside the gate: APPROVED

`/importacoes` outside `InstallationGatedRoutes`. The justification stands because it was not
invented: the comment at `AppRouter.tsx:62-65` already states that catalog/stock/import screens read
the ERP mirror, which exists before any marketplace. The screen sat behind the gate by accident of
living in `/vinculos`.

**EVIDENCE must record that the gate was WRONG for this screen, quoting that comment** — that is what
makes the placement a declared correction instead of a silent change.

## On the self-disclosure

Writing `ImportChainPanel.tsx` + `useErpImportChain.ts` inline exceeded the core §4.2 glue ceiling;
deleted, never committed, reported before being asked. Recorded as the correct handling. Work goes
through workers with failing-test-first.

L2 remains the hub's; REQUEST when green.
