# EVIDENCE — CHIP-T1-FE (solver-labels)

- Mission: MIS-004-mvp-demo · Milestone: M-07-simulador · DEMO-CRÍTICO (demo 2026-07-20)
- Branch: `chip/t1-fe-solver-labels`
- Base: `main` @ `18fbd91a` (ledger D-78 — DESIGN-TARIFAS-ML ratified; wave dispatched)
- Binding design: `DESIGN-TARIFAS-ML.md` §4.4 Layer 3 + §6.1/6.2 (Simulador surface only)
- Scope (exclusive): `apps/web/src/pages/precos/**` — files touched: `SolverPanel.tsx`, `SolverPanel.test.tsx`

## Boot integrity note (FINDING)

Dispatched into worktree `elastic-vaughan-cdd03d` @ `3d158885` — **wrong base**: it lacked
`.mnfs/MIS-004-mvp-demo`, `apps/web/src/pages/precos/`, and the chip branch was **not
pre-provisioned**. Correct base = `main` @ `18fbd91a` (carries MIS-004 + precos + design doc).
Created `chip/t1-fe-solver-labels` from `18fbd91a` via `git worktree add`. → Field finding
for hub: chip worktrees are not always pre-provisioned; dispatch base SHA should be pinned to
the 40-hex, not "current main HEAD" (ambiguous when the launching session sits on a stale worktree).

## Deliverables (design §4.4 Layer 3 + §6.1/6.2)

1. Result branches keyed on backend `code` + `desconhecidos`:
   - `SEM_FRETE` / `desconhecidos:["frete"]` → actionable guidance (exact copy: "Sem dados de
     frete para este produto. Cadastre dimensões (peso, altura, largura, comprimento) OU vincule
     um anúncio ML para cotar o frete."). Never "inatingível", never a blank ceiling.
   - `SEM_CUSTO` / `desconhecidos:["custo"]` / `blocking_state==="SEM_CUSTO"` → existing blocking banner.
   - `UNREACHABLE_TARGET` (or legacy no-code) → "melhor margem possível: X%" ONLY when
     `ceiling_pct` non-empty; blank ceiling never renders as a lone "%".
   - Residual not-reached / no legitimate ceiling → honest generic `solver-incomplete` banner.
2. Price result renders `tarifa` badges (fonte/degrau/data via `FreshnessIndicator`) + `ESTIMATIVA`
   pill (degrau 4 / `estimativa`). NO-DATA (`sem_dados`, nil, or blank value) → `UnknownValue` "—",
   carimbos suppressed. Never R$0, never misleading green (ADR-17).
3. Reads `result.desconhecidos` + new `tarifa` block with tolerant parsing — missing/partial
   tarifa renders without badges, no crash. `tarifa` typed locally (`SolveResult` widen) since the
   SDK `PricingSolveResponse` does not carry the block yet — **sdk-runtime NOT touched** (owned by
   backend chip CHIP-T1; reconcile on their COMMITTED schema).

## Gates (all green)

| Gate | Command | Result |
|---|---|---|
| Unit (target) | `vitest run SolverPanel.test.tsx` | **16 passed** (7 baseline + 9 new) |
| Unit (full web, regression) | `vitest run` (apps/web) | **233 passed** (31 files) |
| Typecheck delta | `tsc --noEmit -p apps/web/tsconfig.json` | **0 real delta on my files** (only repo-wide jest-dom matcher tooling gap, present on baseline `main` too; production `SolverPanel.tsx` = 0 errors) |
| Build | `vite build` (apps/web) | **exit 0** |

Baseline (main) tsc error count 352 → worktree 366; the +14 are all pre-existing-category
jest-dom `Assertion` matcher errors from added test assertions, not new type defects.

## P6 dual gate

- Cold **Opus** reviewer: **PASS** (minors/nits only; ADR-17 honesty, precedence, tolerant parsing,
  copy match, scope all confirmed clean).
- Adversarial **sonnet** reviewer: initial **CHANGES-REQUIRED**, 4 findings — all in scope, all valid:
  - F1 (BLOCKER) empty-string `valor` bypassed NO-DATA → blank "R$ " + false ESTIMATIVA → FIXED (trim/empty guard + test).
  - F2 (MAJOR) `unreachable` ignored `code` → DADOS_INCOMPLETOS+ceiling faked "X%" → FIXED (gate to `UNREACHABLE_TARGET`/no-code + test).
  - F3 (MAJOR) banners ignored `reached` → stale code hid real price → FIXED (all banners gated `!reached` + test).
  - F4 (MAJOR) no reset on product/modalidade change → stale result under new product → FIXED (`useEffect` reset + test).
  - Nit: carimbos beside "—" → FIXED (suppressed on NO-DATA).
- Re-review: confirmation requested with updated diff (agreement closure).

## Test infra (chip worktree has no node_modules)

Junctioned `node_modules` → main checkout; chip-local `apps/web/vitest.config.chip.ts`
(`server.fs.allow` + absolute setupFiles) to run under the junction. **Chip config DELETED before
final commit** — not shipped. Verification here = vitest + tsc + build only; browser/visual QA vs
DESIGN-REFERENCE is hub-owned P7.
