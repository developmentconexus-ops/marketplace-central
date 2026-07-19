# MIS-004 — App-Wide Screen Inventory (design-fidelity gate)

```yaml
id: MIS-004-SCREEN-INVENTORY
type: qa-inventory
owner: Dispatch Hub (Wave C)
parent: MIS-004
created: 2026-07-19
base_sha: 390d79ab
design_truth:
  - .mnfs/MIS-004-mvp-demo/research/design-screens-2026-07-17.md   # R-02 semantics
  - .mnfs/MIS-004-mvp-demo/design/DESIGN-REFERENCE.md              # tokens/binding
evidence_mode: computed-style + a11y tree (pixel-diff impossible — rasterizer broken F-ENV-10)
```

## Why this exists

Operator concern (2026-07-19): the mission-close QA gate (C01) is a **journey rehearsal** — it
walks the 6-step demo path and checks each stop reaches its milestone-defined state. It does NOT
sweep every screen detail-by-detail, and it says nothing about screens **off** the demo path or
screens that are **missing/placeholder**. This inventory is the app-wide map the gate was missing.
It is registered as mission criterion **C10** (below), run live at mission close.

**Method honesty:** no pixel screenshots on this machine (F-ENV-10). Fidelity = computed-style
values + accessibility-tree structure, read per element. This catches *presence/absence/wrong-value*
and *wrong-token*; it does NOT catch sub-pixel layout drift. Operator accepted this limit.

## Status legend

✓ present · ✗ absent · ≠ differs · ⏳ in-flight (chip not merged — not measured) · ▢ placeholder/stub

## Screen map (route × page × design counterpart)

| Route | Impl today | On demo journey | Design counterpart (R-02) | Structural status |
|---|---|---|---|---|
| `/` (index) | DashboardPage | yes | Dashboard | ⏳ **M-09 in flight** — not measured |
| `/catalogo/produtos/:id` | ProdutoRoute | yes | Produto Detalhe | ⏳ **M-06 in flight** — not measured |
| `/anuncios` | AnunciosPage | yes | Anúncios | measured ↓ §1 |
| `/precos` | precos/PricingPage | yes | Simulador | measured ↓ §2 |
| `/pedidos` | pedidos/** | yes | Pedidos | measured ↓ §3 |
| `/vinculos` | vinculos/** | yes | Vínculos | measured ↓ §4 |
| `/catalogo` | CatalogPageWrapper | partial | ~Produto (list) | real screen |
| `/estoque` | StockSeguroPageWrapper | no | aba Estoque (MIS-005) | real screen, off-MVP |
| `/classifications` | ClassificationsPageWrapper | no | — (no MVP counterpart) | real screen, off-theme |
| `/protocolos/:id` | ProtocoloPage | indirect | — (governance surface) | real screen, off-theme |
| `/integracoes` | WorkspacePlaceholder | no | Config→Integrações (MIS-005) | ▢ placeholder |
| `/marketplaces` | WorkspacePlaceholder | no | "Mercado" = em breve | ▢ placeholder |
| nav Mercado / Repasses | disabled spans "em breve" | no | design chose "em breve" | ▢ honest stub |
| `/products`,`/orders`,`/simulator`,… | LegacyRedirect | n/a | n/a | redirect |

## Field-level findings (measured screens, base 390d79ab)

Full element tables live in the investigator report (session transcript). Rollup ✓/✗/≠ per screen:

| Screen | ✓ | ✗ (absent) | ≠ (differs) |
|---|---|---|---|
| **Anúncios** | ~11 | 4 — Exportar, +Criar anúncio, estoque `0⚠`, drawer Vincular/Editar | 6 — chips no Resumo (não header) & rótulos, "Agrupar" é checkbox, bulk 6-set vs 3, grupo-row raso (sem chevron/meta ERP/pill erro), drawer ações disabled, contagem no Resumo |
| **Simulador** | ~9 | 7 — **matriz-tabela**, **coluna VEREDICTO**, chip CEP, atalhos mediana/menor/cobrir-verde, "Criar anúncio", nota frete≥R$79, Restaurar padrão, box comissões read-only | 3 — picker (aside simples s/ busca), DIFAL em drawer separado (não no Params), bidirecional = 2 painéis |
| **Pedidos** | ~11 | 2 — filtros período/logística, banner DIFAL condicional | 6 — DIFAL soma "—" (hub C2), fila sem tier DIFAL, decomposição toda honest-null, NF sem nº/deep-link ERP, ações disabled, status chip s/ "entregue" |
| **Vínculos** | ~4 | 3 — header importação/Nova importação, fallback "Criar produto", negativo Criar/Ignorar | 7 — KPIs≠tabs & rótulos, tabela 6-col product-cêntrica vs 9-col anúncio-cêntrica, GTIN dobrado, "Resolvidos" = stub TODO, **tema slate/blue/emerald off vs paper+green**, ações |

## Gaps that touch the demo path — for operator triage

Ranked by demo-visibility (Anúncios→Simulador→Pedidos→Vínculos):

1. **Simulador sem matriz-tabela nem coluna VEREDICTO** (`pages/precos/PricingPage.tsx:216-322`) —
   layout single-product em vez da matriz SKU·…·VEREDICTO do design. Maior divergência estrutural.
2. **Vínculos off-theme + tabela reorientada** (`pages/vinculos/**`) — página inteira em
   slate/blue/emerald vs tokens paper+green ratificados; 6-col vs 9-col; "Resolvidos" = stub TODO.
3. **Anúncios linhas-grupo rasas** (`AnunciosTable.tsx:198-208`) — sem chevron/meta "ERP est·N
   anúncios"/pill ok-erro; "Agrupar por produto" é caminho de demo provável.
4. **Ações inertes em todos os drawers** (Anúncios/Pedidos/Fila/Lista/Kanban) — botões disabled "em
   breve"; demo não aciona Faturar/Etiqueta/Corrigir/Pausar/Vincular. **Deliberado** por ruling D-57
   (read-only) + zero-writes ML — provável aceitável, confirmar.
5. **Pedidos decomposição/DIFAL/retorno honest-null** — "—" em KPIs e drawer. **Deliberado** por
   ADR-17 + decompositor hub C2 não fiado — correto, mas "sem números" visual em vários pontos.
6. **Simulador Params** sem "Restaurar padrão" nem box comissões read-only; DIFAL em drawer à parte.

### Deliberate vs undeclared

- Itens #4 e #5 têm **ruling documentada** (D-57, ADR-17) → deliberate-exclusion, aceitáveis p/ demo salvo objeção.
- Itens #1, #2, #3, #6 **não têm declaração de omissão** em nenhum `feature.md`/`milestone.md`
  (investigator não achou; M-05 F-02 só declara mudanças aditivas sobre W1). → **decisão de escopo
  do operador/Strategist antes da demo.**

## Live-pass protocol (runs at mission close, criterion C10)

After M-06 + M-09 merge (so ⏳ rows become measurable), on the clean docker dev stack:
1. Re-run this field inventory including Dashboard + Produto Detalhe.
2. Per demo-path screen: read computed-style of theme-bearing elements (bg, font-family, accent) →
   assert paper+green + Instrument Sans/IBM Plex Mono (catches Vínculos off-theme regression class).
3. Read a11y tree per screen → assert design elements present (columns/chips/drawer sections/KPIs).
4. Record present/absent/differs per element in `validation-result.md §screen-inventory`.
5. Any demo-path screen with an **undeclared** missing/differs element that the operator has NOT
   ruled accept-for-demo ⇒ feeds mission FAIL triage (not auto-fail — operator scope call).
