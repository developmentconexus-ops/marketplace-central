# CHIP-SIM-PARAMS — Evidence Pack

**Chip:** CHIP-SIM-PARAMS (MIS-004 / M-07 simulador). FE-only `apps/web/src/pages/precos/**`.
**Goal:** full `/precos` design parity vs `design/handoff/Simulador.dc.html`.
**Commit:** `56e4f0b7` on branch `chip/sim-params` (worktree `.claude/worktrees/chip-sim-params`). NOT pushed.
**Constraints honored:** no backend / no new endpoints / no SDK changes; read-only demo (D-57);
honest "—" for unknowns (ADR-17); `comissao_pct` never sent (resolver chain intact); tokens only
from `apps/web/src/index.css`; no server boot / no docker / no dep install.

## Status: CLOSED — dual gate in AGREEMENT

Hub ruling (local_d614b8a7): gate 2 runs as a Claude adversarial refuter (HARNESS-CORE §1
sanctioned fallback; codex/GPT-5.6 Sol exhausted til 2026-07-25, post-demo — standing block,
not a chip hold). Priced-demo REQUEST resolved as a stack-pointing gap (data EXISTS in shared
dev DB; worktree stack resolved an isolated empty DB) — final priced live-drive is hub-owned
against the shared stack, no chip action.

## Deliverables (CHIP.md §b) — implemented

| Deliverable | Element / testid | State |
|---|---|---|
| ParamsDrawer — recalcula-ao-vivo badge | `params-badge-live` | ✅ |
| ParamsDrawer — MODALIDADES ML read-only (active=comissão %, others "—") | `params-modalidades-ml` | ✅ |
| ParamsDrawer — REGRAS DO ML verbatim, não-editável | `params-regras-ml` | ✅ |
| ParamsDrawer — Restaurar padrão (client-side, no save) | `params-reset` | ✅ |
| ParamsDrawer — DIFAL in-drawer context + "ver tabela completa" link → DifalDrawer | `params-difal-context` | ✅ |
| Header — "Produtos: N ▾" pill + dropdown select | `produtos-pill` / `produtos-dropdown` / `produtos-option-<id>` | ✅ |
| Header — CEP chip (honest `CEP → {difal_destino_uf ?? "—"}`) | `cep-chip` | ✅ |
| Header — resumo params echo + Parâmetros trigger; standalone DIFAL button REMOVED | `params-trigger` | ✅ |
| Sim panel — ONE unified 380px/12px card (header strip + Preço/Margem grid) | `<aside>` 380px | ✅ |
| Sim panel — quick chips (mediana / menor conc. / cobrir Z%) | `quick-chip-mediana` / `-menor` / `-cobrir` | ✅ |
| Sim panel — MODALIDADE 3-col segmented | 3-col grid | ✅ |
| Sim panel — decomposição / solver / frete-nota / comparação / aplicar / cenários regions | `region-*` / `frete-nota` | ✅ |

## Gate verdicts

### Gate A — cold Opus P6 reviewer (harness:gate-reviewer, read-only)
**Verdict: PASS-WITH-NITS.** All 6 rubric criteria PASS with file:line evidence — scope discipline,
no-hardcode, ADR-17 honesty, deliverable coverage, backward-compat, test integrity.
Nits (non-blocking): QuickChips reuses MarketComparison's deduped `["market","aggregates",id]` query;
no golden EXEMPLO-IO fixture.

### Gate B — Claude adversarial refuter (harness skeptic, read-only) — hub-sanctioned fallback
**Verdict: NO-REFUTATION.** Prompted to refute the diff; read all 4 changed source files + tests
+ MarketComparison context. Refuted against all 5 invariants — none breached (file:line table):
1 scope `precos/**` PASS; 2 no `comissao_pct` sent PASS (SolverPanel.tsx:56-63, PricingPage.tsx:180-188,
tests assert `.not.toHaveProperty`); 3 ADR-17 "—" PASS (QuickChips gates on status==="OK";
MODALIDADE/CEP honest "—"); 4 SolverPanel backward-compat PASS (props optional, uncontrolled path
intact); 5 no new network/endpoint/SDK PASS (reuses MarketComparison seam).
**Secondary flag actioned:** inline `bg-[rgba(22,24,20,0.35)]` scrim (ParamsDrawer) → FIXED to
`bg-ink/35` token (ratified app convention `bg-ink/50`, BatchPreviewModal.tsx:61); precos suite
re-run 88 pass / 1 baseline red, no regression.

### P6-DUAL-GATE: AGREEMENT (cold Opus PASS + refuter NO-REFUTATION)

## Live marker

**LIVE-VERIFIED** (dev @ :5175, computed-style DOM inspection — screenshots blocked by F-ENV-10
rasterizer): header row (pill/CEP/resumo/params), Produtos dropdown open+list, Params drawer all
sections (400px fixed), sim card 380px @12px radius with 2-col (168×168) + 3-col (111×3) grids and
every region. Body two-column `flex-wrap gap:14px`, matrix `flex-1 min-w-480`, aside `380 shrink-0`.
Side-by-side not force-visualizable (headless viewport width 0 under F-ENV-10); CSS mirrors the
prototype responsive rule 1:1.

**Priced-path live gap:** dev DB has only ERP facts (<90000); `/catalog/products@btoa("90000")`=[]
so the priced mono-decomposition can't render live. Covered deterministically by `PricingPage.test`
90001 fixture (`decomp-tarifa-comissao`="Cotação"). → REQUEST: provision priced demo block for final QA.

**Priced-path gap CLOSED by hub (stack-pointing, not data — REQUEST resolved).** The 90000+ block
lives in the SHARED dev stack, not the chip's isolated worktree DB. Hub drove it against the shared
stack after syncing the 4 chip source files (ParamsDrawer/PricingPage/QuickChips/SolverPanel) onto
the mounted tree (HMR).

LIVE-VERIFIED: http://localhost:5174/precos — selected priced product 90008 (PAPELEIRA DECA FLEX,
custo ERP R$118.99, real competitor data `ml_catalog_offers` coletado 2026-07-19T21:09:59Z: mediana
R$229.2 / menor R$179.9). Entered venda R$229,20 → mono-decomposition rendered live: (−)Comissão 30.94
**carimbada "Cotação · degrau 3 · Atualização dos dados"** (real tariff-table tier badge), (−)Imposto
9.17, (−)DIFAL 0.00, (−)Tarifa Full 0.00, (−)Custo ERP 118.99, honest "Margem sem DIFAL — não use para
decisão de venda interestadual." Decomposição + tarifa carimbo + real custo + honest margin all present
against live shared-stack data. Priced live-drive PASS (hub-owned, ledger D-118).

## Deterministic gate

- `npx vitest run --config vitest.chip.config.ts` → **97 passed, 1 failed (11 files)**. The single
  failure is `PricingMatrix.test.tsx > (a) priced row` — pre-existing baseline red; `git diff
  --name-only HEAD` shows neither `PricingMatrix.tsx` nor its test in the changeset. Not a regression.
- `npx tsc -p tsconfig.json --noEmit` → **0 errors under `src/pages/precos/**`** (15 baseline errors
  elsewhere, none introduced here).

## Files changed

- `apps/web/src/pages/precos/ParamsDrawer.tsx` + `.test.tsx`
- `apps/web/src/pages/precos/PricingPage.tsx` + `.test.tsx`
- `apps/web/src/pages/precos/SolverPanel.tsx` (controlled target props, backward-compat)
- `apps/web/src/pages/precos/QuickChips.tsx` + `.test.tsx` (new)

`vitest.chip.config.ts` = temp runner (junctioned node_modules + fs.allow), deleted before commit,
never staged.

## Hub-owned follow-ups (not chip blockers)

1. **Merge** `chip/sim-params` → main (not pushed by chip; hub owns).
2. **Priced live-drive** — mono-decomposition + tarifa carimbos + margin chips against the shared
   dev stack at http://localhost:5174/precos (data exists there; hub-owned QA).

See `validation-result.md` (companion). Full deliverable→element checklist and live DOM parity table there.
