# MIS-004 — Mission Validation Result

```yaml
id: MIS-004-VALIDATION-RESULT
type: validation-result
owner: Dispatch Hub
parent: MIS-004
```

## §screen-inventory (criterion C10) — design-fidelity live re-pass

```yaml
run: 2026-07-20
base: main tip @97ef7b09 (dev stack :5174/:8080, docker)
method: app route vs handoff/.dc.html prototype, side-by-side; structure + copy + tokens
evidence_mode: computed-style + a11y tree + rendered-DOM text (rasterizer broken F-ENV-10, no pixel-diff)
scope_this_run: operator-scoped to demo-critical screens — Anúncios, Simulador, Pedidos
              (Dashboard measured then deferred; Produto Detalhe + Vínculos NOT re-measured this run)
caveat: working tree carries uncommitted FIX-4 batch-commission changes (backend + openapi + sdk only).
        apps/web untouched → FE served at :5174 = clean main FE. Structural parity unaffected
        (DESIGN-REFERENCE: data divergence ≠ finding).
```

### Verdict summary

| Screen | Prototype | Core structural parity | Tokens (light+dark) | Residual gaps |
|---|---|---|---|---|
| **Anúncios** `/anuncios` | `Anuncios.dc.html` | ✗ **DIVERGES** (IA differs — CORRECTED 2026-07-20, visual re-pass) | ✓ paper+green, Instrument Sans/IBM Plex Mono | header IA + missing actions + drawer |
| **Simulador** `/precos` | `Simulador.dc.html` | ✓ **MATCH** (post gap#1 D-107/D-112) | ✓ | gap#6 (undeclared) + 3 minor |
| **Pedidos** `/pedidos` | `Pedidos.dc.html` | ✓ **MATCH** (KPIs/views/fila) | ✓ | fila-row density + filters + deliberate read-only |
| Dashboard `/` | `Dashboard.dc.html` | ✗ **DIVERGES** (MPC rebuild) | ✓ | operator-DEFERRED this run |

**Net (CORRECTED 2026-07-20):** no token/theme regressions on any screen. BUT **Anúncios
diverges at the information-architecture level** (header IA, Resumo panel, table-first vs
dashboard-first) — first pass wrongly ✓-marked it from column parity read off source, never
rendered; operator caught it, visual re-pass confirms ✗. Simulador + Pedidos rows below were
ALSO recorded from source read, not screenshots — **treat as unverified until visual re-pass**
(⚠ flagged). Deliberate exclusions (D-57 read-only, ADR-17 honest-nulls) still hold. Anúncios
IA divergence is demo-path visible → operator scope call.

---

### §1 Anúncios `/anuncios` vs `Anuncios.dc.html` — CORRECTED (visual re-pass 2026-07-20)

**Method note:** first pass recorded ✓ MATCH from column-count parity via templated source
read — NEVER rendered prototype side-by-side. Operator rejected ("não tem nada a ver").
Re-run with real screenshots (prototype served over `http://localhost:7395/Anuncios.dc.html`,
app at :5174). Verdict corrected to **✗ DIVERGES** — divergence is information-architecture, not cosmetic.

**✓ present / match (component level)**
- Table columns (MLB·TÍTULO·PRODUTO·PREÇO·EST.·SYNC·QUAL.·PENDÊNCIA) map to prototype grid.
- Chip % vs mercado in PREÇO; bulk action bar = superset vs prototype 3.
- Tokens paper+green light+dark; Instrument Sans + IBM Plex Mono. No theme regression.

**✗ DIVERGES — screen composition / information-architecture**
- **Header:** prototype = compact title + "1.284 · exceções:" + **4 inline clickable exception
  chips** + **[Agrupar ▾][Exportar][+ Criar anúncio]**. App = "**Workspace operacional**" banner
  + generic subtitle, **no inline chips, no Exportar, no Criar** in header.
- **Counters:** app buries exception counts in an **8-cell "Resumo" stat panel** (dashboard
  pattern) instead of the header exception-chips. Not the design's layout.
- **Search:** app leads with a **full-width "Buscar anúncios" box**; prototype has none up top.
- **Body:** prototype is **data-table-first** (table renders immediately). App pushes the table
  far down behind Resumo + search + a large empty gap.
- **Drawer:** prototype ships a **rich 300px sticky detail drawer** (foto, red error box, stat
  grid Preço/Est/Margem/Qualidade, actions, LINHA DO TEMPO) open by default. App drawer not
  equivalent / not shown.
- Net: prototype = dense operational table workspace; app = generic dashboard framing. **Different IA.**

**Deliberate subset (not the cause of divergence):** "+ Criar anúncio" / mutating actions
disabled = read-only demo D-57 + zero-ML-writes. But the header IA, Resumo panel, big search,
and missing drawer are NOT covered by any ruling → **undeclared, operator triage.**

---

### §2 Simulador `/precos` vs `Simulador.dc.html`

**✓ present / match**
- **7-col matrix EXACT**: SKU·DESCRIÇÃO·CUSTO·NOSSO PREÇO·PREÇO MERCADO·MARGEM·VEREDICTO (34 rows).
- Honest ADR-17: CUSTO real per row; **3 OK verdicts lit** (90008 R$229.2 / 90009 R$749.88 / 90010 R$1681.11), rest SEM_EVIDENCIA / MERCADO_INSUFICIENTE; NOSSO PREÇO / MARGEM "—"; "novo" tags; zero fabricated R$0.
- Bidirectional panel: Preço de venda · Clássico/Premium/Full · Margem alvo · Calcular preço · Aplicar preço.
- ⚙ Parâmetros de cálculo drawer, live-recalc: regime SIMPLES/PRESUMIDO · Alíquota · Limiar verde · Limiar amarela · Tarifa Full · Habilitar DIFAL · Salvar.
- Separate "DIFAL por UF" control.

**≠ / ✗ residual**
- **gap#6 (UNDECLARED — carries over from SCREEN-INVENTORY):**
  - No **"Restaurar padrão"** in the Params drawer (only Cancelar / Salvar).
  - No read-only **"MODALIDADES ML — Comissão Clássico/Premium (vem do ML), REGRAS DO ML — NÃO EDITÁVEIS"** box.
- **CEP chip** ("CEP 88301 → 01310") absent.
- "categoria ▾" filter absent.
- "NA ANÁLISE / SUGESTÕES DO RADAR" toggle absent.

---

### §3 Pedidos `/pedidos` vs `Pedidos.dc.html`

**✓ present / match**
- 5 KPIs: NOVOS · A FATURAR · A ENVIAR · ENVIADOS · DIFAL A PAGAR (honest 0/0/0/24/—, "decomposição pendente").
- 3 views Fila / Lista / Kanban (proper a11y tablist).
- Fila "ordenada por urgência (atrasados primeiro, depois SLA de envio)"; explicit "DIFAL indisponível nesta fila".
- Row-click opens order drawer (deep-link `?order=`); fiscal drawer detail validated at M-08 P7 (D-90).
- Tokens paper+green; IBM Plex Mono valores.

**≠ / ✗ residual**
- **Fila rows are thin** vs prototype 9/10-col: no COMPRADOR, no ITENS count, no SLA-envio chip, no DIFAL chip per row (app row = date · canal · "· item · R$" · — · sem ação).
- Filters **período** + **logística** absent from header.
- Per-row action "sem ação"; drawer actions disabled — deliberate read-only (D-57).
- DIFAL sum "—" — hub C2 decompositor unwired, ADR-17 honest (deliberate).

---

### §4 Dashboard `/` vs `Dashboard.dc.html` — DEFERRED (operator, this run)

Measured, then de-scoped. Recorded for the record:
- Tokens/fonts ✓ (light+dark).
- **Structural divergence:** app = MPC-rebuilt catalog-health cards (Vendas hoje/7d, Anúncios ativos/erro/margem, Sem vínculo/GTIN); prototype = fulfilment cards (Para embalar/enviar/entregar/faturar) + Saúde-dos-anúncios donut + Mais vendidos + RESUMO (GMV) + CANAIS (4 marketplaces).
- Divergence is a design-intent choice (mission accepted-assumption "Dashboard recriado MPC"; prototype data is fake multi-channel, MVP is ML-only) — **undeclared** in any milestone doc. Not ruled. Out of this run's scope.

---

### Operator triage — undeclared gaps not yet ruled accept-for-demo

Ranked by demo-visibility:

1. **Anúncios IA divergence** (visual-confirmed): header exception-chips + Exportar + Criar
   absent; counts pushed into a Resumo dashboard panel; full-width search added; table buried
   below a gap; rich 300px drawer missing. NOT cosmetic — whole-screen shape differs from design.
   Biggest demo-path visual gap on a wave-B screen. Operator scope call before demo.
2. **Simulador gap#6** (Params drawer): no "Restaurar padrão" + no read-only ML-commission box. Small FE add.
3. **Pedidos fila-row density**: buyer/SLA/DIFAL-chip columns thinner than prototype. FE add, data mostly available.
4. **Simulador CEP chip / categoria / radar toggle** absent. Secondary.
5. **Dashboard** taxonomy divergence — separate scope decision (deferred).

⚠ **Simulador + Pedidos not yet visually re-verified** — their §2/§3 findings came from source
read, same weak method that mis-scored Anúncios. Re-pass with screenshots before trusting.

Deliberate-and-ruled (not findings): read-only actions across all screens (D-57), ADR-17 honest-nulls (DIFAL sum, margem categórica), zero-ML-writes exclusions (Criar anúncio).
