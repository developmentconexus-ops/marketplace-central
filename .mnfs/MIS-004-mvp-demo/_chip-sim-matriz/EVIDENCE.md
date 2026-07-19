# CHIP-SIM-MATRIZ — Evidence Pack

**Branch:** `chip/sim-matriz` · **Base:** `783cbc0d` · **Head:** `aeee4b53`
**Scope:** `apps/web/src/pages/precos/**` ONLY · **Mission:** MIS-004-mvp-demo, design-fidelity gap #1 (criterion C10)

## Deliverables

The design's **product matrix** is now the `/precos` (Simulador) main surface.

- **`PricingMatrix.tsx`** (new) — matrix with the 7 design columns
  `SKU · DESCRIÇÃO · CUSTO · NOSSO PREÇO · PREÇO MERCADO · MARGEM · VEREDICTO`,
  one row per catalog product. Row-click hands selection to the page.
- **`PricingPage.tsx`** — restructured: matrix in the `flex-1` main column; the
  **existing** DecompositionPanel + SolverPanel (+ market / apply / scenarios regions)
  moved into a conditional 380px `<aside>` opened by row selection. Parâmetros + DIFAL
  drawers preserved at page root. No panel rebuilt.
- **`MarketComparison.tsx`** — one-line change only: `apiBaseUrl()` exported so the
  matrix reuses it (no duplicate helper).

### Data lanes (reuse only — no new endpoints; OpenAPI/SDK FROZEN)
- **CUSTO / NOSSO PREÇO** — straight off `CatalogProductFact` (no per-row fetch).
- **PREÇO MERCADO / VEREDICTO** — ONE `listMarketAggregates(codprods[])` call, indexed by
  `product_id`. VEREDICTO maps the aggregate **status only** → `OK / SEM_EVIDENCIA /
  MERCADO_INSUFICIENTE`. The categorical MARGIN verdict is M-07-owned and NOT invented.
- **MARGEM** — per-row existing `pricingDecompose`, `comissao_pct` OMITTED so the backend
  resolver chain (COTACAO → PADRAO) runs, never a MANUAL override.
- **"novo" tag** — per-row `listListingsByProduct`, shown ONLY on a query RESOLVED to zero
  listings (lane disabled when no installation is known).

### ADR-17 (unknown ≠ zero/default) — honest across every state
Missing custo → `—`. No market evidence → `—` + honest verdict. Missing/blocked margin →
`UnknownValue` (`margem: M-07`). **Loading** market/decompose → honest unknown (`carregando…`),
never a verdict. **Errored** market/decompose → honest unknown (`falha ao carregar` /
`falha ao calcular`), never SEM_EVIDENCIA nor the M-07 placeholder. Rank `19º/23` is absent
from `MarketPriceIntelAggregate` → correctly omitted (reported as FINDING below).

## Gates

**Web vitest** (chip junction to main `node_modules` + throwaway `vitest.config.chip.ts`,
deleted pre-commit): full `src/pages/precos/**` suite — **78/78 pass**, 9/9 suites.
`PricingMatrix.test.tsx` = 10 tests incl. EXEMPLO-IO golden (a)/(b)/(c), single-market-call,
rank-absent, comissao_pct-absent, INSUFFICIENT_MARKET, market-error honesty, decompose-error
honesty, and both "novo" cases.

**L0 tsc** (`apps/web`, `tsc --noEmit`): **10 errors = exact documented main baseline**,
**ZERO in `src/pages/precos`**. (See memory `web-tsc-lane-cross-branch-resolution`.)

**P6 dual gate — Claude side (cold Opus + adversarial refuter, agreement required):**

*Round 1* (base diff): Cold Opus `gate-reviewer/opus` **PASS**; adversarial `gate-reviewer/sonnet`
**FAIL** — flagged D1 (ADR-17: loading/error states collapsed into confirmed SEM_EVIDENCIA / the
M-07 placeholder) and D3 (INSUFFICIENT_MARKET branch untested). D2 (missing markers) = this
paperwork, not code. **No agreement → fix, do not merge.**

Fixes: commit `02ed4a14` gates VEREDICTO/PREÇO MERCADO on `marketQuery` resolution and the margin
cell on decompose success/error; adds INSUFFICIENT_MARKET + market-error + decompose-error tests.
Commit `aeee4b53` hardens those tests to positively assert the error-branch hint title.

*Round 2* (hardened diff): Cold Opus `gate-reviewer/opus` **VERDICT PASS**, blockers none.
Adversarial `gate-reviewer/sonnet` **VERDICT PASS** — conceded the functional fix on all attack
angles; flagged two non-blocking test-rigor gaps (assert positive error-hint title), applied in
`aeee4b53`.

**P6-DUAL-GATE: AGREEMENT** — both reviewers PASS on the hardened head `aeee4b53`.

**LIVE-VERIFIED: pending** — hub-owned P7 browser QA. FE-only render slice reusing FROZEN
endpoints; zero provider writes, no live integration in this diff. Chip does NOT self-drive P7.

## FINDING (for hub ratification)
`MarketPriceIntelAggregate` (SDK `packages/sdk-runtime/src/market.ts`) carries **no rank/position
field** — `position`/`rank` exist only on the unrelated `MarketPriceIntelSignal`. The design's
`19º/23` rank in PREÇO MERCADO is therefore not renderable from current data and is omitted, per
the dispatch's "rank ONLY if available in data" rule. If the demo needs rank, it requires a
backend/SDK change (out of this chip's frozen-contract scope).

## Hard-rule compliance
Zero ML writes. No server booted / `:8080` / `.env` / docker touched. No push / reset / revert /
stash / clean / `-D` / WSL / cold-clone / cache-purge. No deps installed. Scope confined to
`apps/web/src/pages/precos/**` (verified via `git diff --stat` on the range). Throwaway
`vitest.config.chip.ts` deleted pre-commit; `node_modules` junction not committed.
