# CHIP-SIM-PARAMS — Validation Result

**Scope:** FE-only (`apps/web/src/pages/precos/**`). Full design parity of `/precos` vs
`design/handoff/Simulador.dc.html`. No backend / no new endpoints / no SDK changes.
Read-only demo (D-57). Honest "—" for unknowns (ADR-17). No hardcoded % (commission
stays on the backend resolver chain; `comissao_pct` never sent).

## Deliverables checklist

### Parâmetros drawer (ParamsDrawer.tsx) — DONE
- [x] Fixed overlay 400px + backdrop (`params-drawer`, `params-backdrop`)
- [x] "recalcula ao vivo" header badge (`params-badge-live`)
- [x] IMPOSTO regime segmented (Simples / L. Presumido) + alíquota (0–35 validated)
- [x] LIMIARES dots (verde/amber) + helper copy
- [x] MODALIDADES ML read-only box (`params-modalidades-ml`) — active tier shows the
      resolved comissão %, others honest "—" (never fabricated)
- [x] REGRAS DO ML — NÃO EDITÁVEIS box verbatim (`params-regras-ml`)
- [x] Tarifa Full (R$) input
- [x] DIFAL toggle pill + destino context + 27-UF select + "ver tabela completa por UF →"
      link → opens DifalDrawer (`params-difal-context`)
- [x] Footer: Restaurar padrão (client-side reset, no save) + Concluído (`params-reset`)

### Header (PricingPage.tsx) — DONE
- [x] Title "Preços & Simulador" (22px bold)
- [x] "Produtos: N ▾" pill + dropdown listing loaded catalog (`produtos-pill`,
      `produtos-dropdown`, `produtos-option-<id>`); click selects + closes
- [x] Destino chip (`cep-chip`) — honest `CEP → {difal_destino_uf ?? "—"}`
- [x] resumo params echo (regime + alíquota + DIFAL flag), right-aligned
- [x] "⚙ Parâmetros de cálculo" button (`params-trigger`)
- [x] Standalone "DIFAL por UF" button REMOVED — reached only via the drawer link
      (design parity)

### Sim panel — ONE unified 380px card (PricingPage.tsx) — DONE
- [x] Header strip "Simular · {name}" + sku (mono, right)
- [x] Preço de venda + Margem desejada 2-col grid (margin lifted from SolverPanel via
      controlled `target`/`onTargetChange`/`hideInput`)
- [x] Quick chips: mediana / menor conc. (from ML market, omitted when no confident
      price — ADR-17) / cobrir X% (seeds solver target) — `QuickChips.tsx`
- [x] MODALIDADE 3-col segmented; active tier shows margem %, others "—"
- [x] Per-component decomposition mono card (`region-decomposicao`, DecompositionPanel)
- [x] Margem-alvo → preço solver (`region-solver`, SolverPanel hideInput)
- [x] Frete rule note (`frete-nota`) — ≥R$79 free-shipping vs <R$79 fixed fee
- [x] Market comparison (`region-comparacao`), Apply action (`region-aplicar`),
      Scenarios (`region-cenarios`)

## Live DOM parity (Tailwind-computed, dev @ :5175)

Measured on a live render (ERP product used only to exercise layout; see gap below):

| Property                 | Prototype            | App (computed)          | ✓ |
|--------------------------|----------------------|-------------------------|---|
| Sim card width           | `width:380px`        | `380px`                 | ✓ |
| Sim card radius          | `border-radius:12px` | `12px` (rounded-card)   | ✓ |
| Preço/Margem grid        | `1fr 1fr`            | `168px 168px`           | ✓ |
| MODALIDADE grid          | `1fr 1fr 1fr`        | `111px 111px 111px`     | ✓ |
| Body row gap             | `gap:14px`           | `14px`, flex-wrap       | ✓ |
| Matrix column            | `flex:1;min-width`   | `flex-grow:1; min 480px`| ✓ |
| Header                   | full-width, top      | full-width above body   | ✓ |
| Params drawer            | `400px` fixed        | `400px` fixed           | ✓ |

Header live-verified: `produtos-pill`, `cep-chip`, resumo, `params-trigger` present;
no `difal-trigger` (removed). All sim-panel regions present in the single card.

**LIVE-VERIFIED** (dev @ :5175, computed-style DOM inspection — screenshots blocked by
F-ENV-10 rasterizer): header row (pill/CEP/resumo/params), Produtos dropdown open+list,
Params drawer all sections (400px fixed, badge/modalidades-ml/regras-ml/reset/difal-link),
sim card 380px @12px radius with 2-col (168×168) + 3-col (111×3) grids and every region
(quick-chips, decomposicao, solver, frete-nota, comparacao, aplicar, cenarios). Two-column
body confirmed `flex-wrap gap:14px`, matrix `flex-1 min-w-480`, aside `380 shrink-0`;
side-by-side not force-visualizable because the headless viewport reports width 0
(F-ENV-10) — the CSS mirrors the prototype's own responsive rule 1:1. No new provider
integration (reuses the IC-03 MarketComparison seam), so no `LIVE-WAIVED-BY-OPERATOR`
provider marker applies.

## P6 dual gate

- **Gate 1 (cold Opus, harness:gate-reviewer, read-only):** PASS-WITH-NITS — all 6 rubric
  criteria PASS with file:line; both nits non-blocking (shared TanStack-deduped queryKey;
  EXEMPLO-IO golden not mandated for this chip — hub ACK).
- **Gate 2 (Claude adversarial refuter, harness skeptic, read-only — hub-sanctioned fallback
  per HARNESS-CORE §1, codex/GPT-5.6 Sol exhausted til 2026-07-25):** NO-REFUTATION.
  Refuted against all 5 invariants (scope `precos/**`; no `comissao_pct` sent; ADR-17 "—";
  SolverPanel backward-compat; no new network/endpoint/SDK) — none breached, file:line table.
  Secondary flag: inline `bg-[rgba(22,24,20,0.35)]` scrim in ParamsDrawer → **FIXED** to the
  ratified `bg-ink/35` token (app convention = `bg-ink/50`, BatchPreviewModal.tsx:61); re-ran
  precos suite (88 pass / 1 baseline red), no regression.

**P6-DUAL-GATE: AGREEMENT**

## Tests

`npx vitest run --config vitest.chip.config.ts` → **97 passed, 1 failed (11 files)**.
The single failure is `PricingMatrix.test.tsx > (a) priced row` — a **pre-existing
baseline red**: `git diff --name-only HEAD` shows neither `PricingMatrix.tsx` nor its
test among the changed files. Not a regression from this chip.

New/updated coverage:
- `ParamsDrawer.test.tsx` (12) — parity surface + Concluído + restore-defaults
- `QuickChips.test.tsx` (4) — mediana/menor seeds, cobrir→target, ADR-17 omission
- `SolverPanel.test.tsx` — unchanged, still green (controlled props are backward-compat)
- `PricingPage.test.tsx` — produtos picker toggle+list, DIFAL via drawer link, regions

## tsc

`npx tsc -p tsconfig.json --noEmit` → **0 errors under `src/pages/precos/**`**.
15 remaining errors are all baseline (AnunciosPage/anunciosQueries/ProdutoPage/
mutations/ListingsRefresh), none introduced by this chip.

## Known gap (backend data, not FE)

The worktree dev stack's catalog has only ERP-bulk facts (ids < 90000) — the priced/
listed demo block (90000+, e.g. Eletrodo 6013) is not in this DB. `/catalog/products`
at the priced cursor returns `items:[]`. Live sim-panel render for the parity
measurements above used the ERP block (temporary, reverted — cursor restored to
`btoa("90000")`). Decomposition/market values live-showed the honest "—"/placeholder
because ERP facts carry no price. The priced-path rendering (mono decomposition with
tarifa carimbos, margin chips) is covered deterministically by the vitest fixtures
(`PricingPage.test` 90001 fixture asserts `decomp-tarifa-comissao` = "Cotação").
→ REQUEST to hub: provision the priced demo block in the dev DB for final browser QA.

## Files changed

- `apps/web/src/pages/precos/ParamsDrawer.tsx` + `.test.tsx`
- `apps/web/src/pages/precos/PricingPage.tsx` + `.test.tsx`
- `apps/web/src/pages/precos/SolverPanel.tsx` (controlled target props, backward-compat)
- `apps/web/src/pages/precos/QuickChips.tsx` + `.test.tsx` (new)

`vitest.chip.config.ts` is a temp runner (junctioned node_modules + fs.allow) and is
deleted before commit — never staged.
