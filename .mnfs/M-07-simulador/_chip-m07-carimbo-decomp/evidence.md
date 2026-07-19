# CHIP-M07-CARIMBO-DECOMP — evidence pack

**Mission:** MIS-004-mvp-demo · **Milestone:** M-07 simulador · **Chip:** carimbo-decomp
**Branch:** `chip/m07-carimbo-decomp` · **Base:** `3f8b39048bb5e33bdee019c1b429382c1ad1c35b`
**Head:** `d36cb87a0f30911f76f4eb3b2d6b36fa38c40a89`
**Commit range:** `3f8b3904..d36cb87a` (2 commits: feat + dual-gate fix)

## Scope delivered
FE-only. Surfaced the tarifa provenance carimbo (fonte + degrau + ESTIMATIVA pill +
FreshnessIndicator) on the DecompositionPanel Comissão and Frete rows, reusing SolverPanel's
proven badge via a new shared module. Honest-absence (ADR-17) preserved: NO-DATA keeps "—",
carimbo suppressed; custo/imposto/DIFAL/taxa/tarifa-full rows stay value-only.

### Field finding (reported to hub)
Dispatch context assumed the SDK type lags (`tarifa` block not in `PricingDecomposeResponse`)
and mandated a local widening cast as SolverPanel does. **At base SHA the SDK ALREADY carries
it**: `PricingTarifa` / `PricingTarifaComponent` / `PricingTarifaFrete` +
`PricingDecomposeResponse.tarifa?` and `PricingSolveResponse.tarifa?`
(`packages/sdk-runtime/src/index.ts:1507-1546`). So NO local SDK cast was needed on the
decompose path — `PricingPage` passes `decomposeQuery.data.tarifa` straight through into a
tolerant local `TariffBlock` prop (SDK types assign to it cleanly). SolverPanel's pre-existing
`SolveResult` cast was left untouched (out of scope; harmless).

## Changed files (write-set reconciliation — `git diff --name-only 3f8b3904..d36cb87a`)
| Path | Declared owned seam | Status |
|---|---|---|
| apps/web/src/pages/precos/tariffBadge.tsx (NEW) | yes (optional new module) | declared |
| apps/web/src/pages/precos/DecompositionPanel.tsx | yes | declared |
| apps/web/src/pages/precos/DecompositionPanel.test.tsx | yes | declared |
| apps/web/src/pages/precos/PricingPage.tsx | yes | declared |
| apps/web/src/pages/precos/PricingPage.test.tsx | yes | declared |
| apps/web/src/pages/precos/SolverPanel.tsx | yes (badge extraction) | declared |

All within owned seam `apps/web/src/pages/precos/**`. Zero undeclared paths. Chip-local
`apps/web/vitest.config.chip.ts` was created for the verify lane and DELETED before both
commits — not in the diff.

## Dispatch ledger
| Role | Model | Dispatch | Output |
|---|---|---|---|
| Implementer (fallback — codex quota dead til 2026-07-25) | claude sonnet (Agent) | TDD impl-pack, full spec inline | 6 files, vitest 75→ green |
| Dual gate A (cold) | **opus** (Agent, clean context) | REVIEW-STANDARD, git read-only, diff at chip.diff | VERDICT: PASS (1 suggestion, 2 nits) |
| Dual gate B (adversarial) | **sonnet** (Agent, clean context) | refute-the-diff, git read-only | round 1: FAIL (1 important) → delta re-review after fix |
| Dual gate delta re-review | opus + sonnet (resumed) | delta.diff (4195e63d..d36cb87a) | see reconciliation |

Codex roles (planner/implementer/gate GPT side) NOT dispatched — quota exhausted til
2026-07-25 (memory codex-quota-exhausted-jul25); Claude-only dual gate per chip dispatch.

## Verification (P5 — all `ran`, worktree with node_modules junction to main)
- **vitest** (`npx vitest run --config vitest.config.chip.ts`, 9 precos suites): **76 passed**
  (baseline 70 → +6: 5 carimbo tests + 1 ADR-17 regression test). `ran`.
- **vite build** (`npx vite build`): **✓ built**. `ran`.
- **tsc** (`npx tsc --noEmit`): 463 total errors, **ALL** the known env baseline
  (`toBeInTheDocument does not exist on type 'Assertion'` — jest-dom matcher types unresolved
  via the node_modules junction; present repo-wide on unchanged files, e.g. VinculosPage,
  and on unchanged SolverPanel.test). **ZERO non-baseline errors on the changed surface.**
  Canonical tsc baseline is the hub env; what this chip guarantees is zero NEW type errors on
  its files. One real type error found during verify (`sem_dados` on `PricingTarifaComponent`
  in a test fixture) was fixed before commit. `ran`.

## Dual-gate reconciliation
Round 1: **Opus PASS**, **Sonnet FAIL** — disagreement on SEVERITY of the same fact: the
DecompositionPanel frete carimbo (source `tarifa.frete`) was decoupled from the frete value's
own null-check (`decomposition.frete`), so a backend contradiction could render a carimbo
beside "—" (ADR-17 honest-absence). Opus rated it non-blocking (unreachable: comissão is
non-null by type; frete backend-contract-guaranteed); Sonnet rated it blocking.

Resolution (adopted the stricter reviewer): ADR-17 honest-absence is a §7 non-negotiable, so
the invariant was FE-ENFORCED rather than backend-trusted — frete carimbo now gates on
`decomposition.frete !== null` (commit d36cb87a) + regression test. Comissão intentionally
left unguarded (non-null by SDK type → a guard would be dead code / anti-slop). Delta
re-reviewed by both gates.

**Final dual-gate verdict:** **PASS + PASS (agreement).** Delta re-review on delta.diff:
- Opus (cold): PASS — frete gate mirrors `Value`'s null check; comissão exemption type-sound; regression test real, not theater.
- Sonnet (adversarial): PASS — finding fully resolved; carimbo can no longer render beside "—" regardless of tarifa.frete; test fails if guard removed.

## Defers / notes
- SolverPanel's `SolveResult = PricingSolveResponse & { tarifa?: TariffBlock }` cast is now
  redundant (SDK carries `tarifa`) but left as-is — out of scope, zero behavior impact.
- Opus nit: the "no carimbo on custo/imposto" test title slightly overclaims (asserts the two
  carimbos + custo/imposto values, not the absence of a custo carimbo). Left as-is (no testid
  on those rows to assert against; code is correct).
- Browser P7 QA is the hub's post-merge responsibility (chip ran no server / docker / .env).
