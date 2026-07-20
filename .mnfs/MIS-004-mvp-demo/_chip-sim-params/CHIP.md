# CHIP-SIM-PARAMS — Simulador: Parâmetros de cálculo + painel a paridade de design

```yaml
id: CHIP-SIM-PARAMS
type: chip-context-pack
parent: MIS-004 / M-07 (simulador)
owner_on_launch: milestone session (Opus, worktree isolation)
hub_session_ref: title-match "Dispatch Hub" (capture real local_ id on first event)
lane: FE-only (apps/web)
parallel_safe_with: CHIP-ANUN-IA (disjoint file trees; shared read-only: apps/web/src/index.css tokens)
```

## (a) Why / context

C10 design-fidelity visual re-pass (2026-07-20, screenshots both sides, prototype served on
`http://localhost:7395/Simulador.dc.html`) found: Simulador **core is faithful** (matrix +
bidirectional panel + decomposition + market-evidence). The gap is the **Parâmetros de cálculo
drawer** — it carries the skeleton but omits the design's read-only ML-rules context and the
"restore defaults" affordance. Operator flagged this directly: "como a tela comporta os
parâmetros de cálculo… temos que replicar o design". Evidence: `validation-result.md §2`.

Design truth = `.mnfs/MIS-004-mvp-demo/design/handoff/Simulador.dc.html` (the `.dc.html`, NOT the
distilled reference). Fidelity HIGH: tokens/copy/layout final. Data divergence ≠ finding.

## (b) Deliverable (scope) — files owned

Enrich EXISTING components (no greenfield):

1. **`apps/web/src/pages/precos/ParamsDrawer.tsx`** (comp @46, footer @143–155):
   - Add **MODALIDADES ML — read-only box**: "Comissão Clássico {x}% · Premium {y}% — vem do ML"
     + "Tarifa Full por unidade R$ {z}". Rates are display-only (source below). Non-editable.
   - Add **"REGRAS DO ML — NÃO EDITÁVEIS"** gray box, static copy verbatim from prototype:
     "Taxa fixa R$ 6,50 em vendas abaixo de R$ 79 · frete grátis obrigatório (vendedor paga) a
     partir de R$ 79 · frete calculado por peso e CEP · custo do produto vem do ERP".
   - Add **"Restaurar padrão"** button in footer (left of Cancelar/Salvar) — resets the drawer
     form fields to their default values (client-side; no persistence call beyond existing Salvar).
   - Add **"recalcula ao vivo"** header badge next to the drawer title.
   - Add **DIFAL in-drawer context line**: current UF destino + rate derived from the analysis CEP
     (e.g. "SP — do CEP da análise · 4%"), with a link/button that opens the existing
     `DifalDrawer` for the full UF table. (Keep DifalDrawer as the full table — richer than the
     prototype mini-table; do NOT delete it.)
2. **`apps/web/src/pages/precos/DecompositionPanel.tsx`** and/or **`SolverPanel.tsx`**:
   - Add sim-panel **quick chips**: "mediana R$X", "menor conc. R$Y", "cobrir Z%" — one-tap seeds
     that populate Preço de venda from the already-loaded market comparison (median / lowest valid
     offer / cover-margin). Read-only compute from existing data; no new fetch.
3. **`apps/web/src/pages/precos/PricingPage.tsx`** (@60, header near matrix):
   - Add **CEP chip** ("CEP {origem} → {destino}") in the header, from existing analysis CEP.
   - Add **"Produtos: N ▾"** filter pill (count of matrix rows; opens existing product filter if
     present, else a simple count display — do NOT invent a new filter backend).

Commission % source (display-only, already in FE):
- `apps/web/src/pages/precos/tariffBadge.tsx:10` `TariffComponent.valor` (comissão %)
- `apps/web/src/pages/precos/DecompositionPanel.tsx:69–73` comissão line + tariff carimbo (fonte/degrau)
- `apps/web/src/pages/precos/SolverPanel.tsx:137–142` solver tarifa block

## (c) Out of scope / DO NOT

- **No backend / no new endpoints / no OpenAPI/SDK changes.** Everything is FE display of data the
  page already loads. If a value is genuinely not in FE state, show honest "—" (ADR-17) — never
  fabricate; if you believe backend data is required, send `REQUEST` to hub, do not add a fetch.
- **No ML writes, no mutation.** Read-only demo (D-57). "Salvar parâmetros" keeps its existing
  behavior only; "Restaurar padrão" is client-side form reset.
- **Do not touch** the matrix columns, the solver math, DifalDrawer's UF table logic, or any file
  outside `apps/web/src/pages/precos/**` (+ read-only `index.css` tokens).
- **Do not** edit `.env*`, run a server / bind :8080, load docker, or install deps. Need the dev
  stack? `REQUEST` to hub. Tokens come from `apps/web/src/index.css:13–36` — reuse, never hardcode
  hex, never Inter/Roboto, no gradients/emojis (HANDOFF non-regression rules).

## (d) Validation contract (only QA passes)

- **Visual parity**: render app `/precos` vs `http://localhost:7395/Simulador.dc.html` (hub serves
  the static prototype). Open Params drawer on both; screenshot side-by-side. Assert each (b) item
  present and shaped per prototype. Light AND dark (`data-theme`).
- Tokens: computed-style of drawer bg/accent/ink = paper+green vars; fonts Instrument Sans +
  IBM Plex Mono (numbers). No new hardcoded colors.
- `apps/web` **tsc clean** (cite main baseline = 10 type-only, do not regress) + **vitest green**
  for any new/changed component (junction node_modules per FE-chip vitest rule; delete temp vitest
  config pre-commit).
- Evidence pack in `.mnfs/MIS-004-mvp-demo/_chip-sim-params/`: screenshots both sides + tsc/vitest
  logs + a checklist mapping each (b) deliverable to a rendered element. Unwritten = didn't happen.

## (e) Dual gate + events

- Close = P6 dual-gate (cold Opus reviewer + GPT-5.6 Sol medium, agreement) then fresh browser QA
  live-drive of `/precos`. Only QA passes.
- Talk to hub via events only: `CLOSED` (with evidence) / `BLOCKED` / `ESCALATION` / `REQUEST` /
  `SPLIT-REQUEST` / `COMMITTED`. Anything needing the dev stack, a dep, or backend data = `REQUEST`.

## (f) Superseded-protocol denylist (skill pin)

Obey `docs/HARNESS-PROFILE.md §10` superseded-protocol denylist. Do NOT invoke deleted MNFS
execution skills (feature-execution & siblings, removed @6b29412; stale codex cache 0.1.0). Binding
doctrine = `HARNESS-CORE.md` + `docs/HARNESS-PROFILE.md` + mission `## Parallel Execution Plan`.
Never push. `git branch -d` never `-D`. Never reset/revert/stash/clean unknown state (uncommitted
FIX-4 batch-commission work in tree is NOT yours — leave it).
