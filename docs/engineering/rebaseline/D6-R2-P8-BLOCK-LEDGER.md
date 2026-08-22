# D6-R2 P8 — Structural Wireframe Block Ledger

> **Status:** OPEN / ACTIVE — B00 rendered as HTML `CANDIDATE`; no block is operator-`LOCKED`
> **Parent:** [D6-R2 Complete Frontend Realization Closure](D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md)
> **P5 input:** [Complete Screen / Material-Surface Inventory](D6-R2-P5-SCREEN-SURFACE-INVENTORY.md)
> **Method:** [Frontend Product Experience Planning Method v2.1](../../development/frontend-product-experience-planning-method.md)
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. P8 artifact law for Marketplace Central

P8 wireframes are **executable low-fidelity HTML structural prototypes**, not visual-design comps and not static images.

They decide only material UX structure such as:

```text
shell and region placement
relative size / density class
navigation and context placement
reading / action order
material state placement
responsive transformation
interaction needed to prove the block
```

They deliberately do **not** decide final palette, typography, iconography, radius, shadows, illustration, branding polish or final component styling. Those belong to later visual-design handoff and may not silently alter locked UX structure.

A static styled image generated during the B00 attempt is **NON-AUTHORITATIVE / DISCARDED** and is not P8 evidence. The only current B00 rendered candidate is the HTML artifact below.

## 2. B00 — App Shell + global IA

**Status:** `CANDIDATE` — awaiting operator review of the rendered/interactable HTML.

**Artifact:** [`qualification/d6-r2-wireframes/b00-app-shell.html`](../../../qualification/d6-r2-wireframes/b00-app-shell.html)

**P6:** NOT TRIGGERED.  
**P7:** NOT TRIGGERED.

### 2.1 Structural candidate

Desktop baseline:

```text
prototype-only evidence strip
└─ product frame
   ├─ persistent left sidebar ≈ 264 px
   │  ├─ Marketplace Central text identity
   │  ├─ Organization selector — global workspace
   │  └─ accepted grouped primary navigation
   └─ main
      ├─ page header ≈ 64 px
      │  ├─ current page title
      │  └─ page-local Marketplace Installation host when admitted/required
      └─ page-owned content host with ≈ 24 px outer padding
```

Responsive candidate:

```text
desktop  persistent sidebar
→ tablet narrower/collapsible navigation
→ mobile drawer + stacked page-local context
```

No responsive transition may remove Organization/Installation qualification or change navigation meaning.

### 2.2 Accepted IA rendered

```text
VISÃO GERAL
  Visão geral

OPERAÇÕES
  Preparação
  Publicações
  Disponibilidade
  Vendas
  Expedição
  Pós-venda

ESTRATÉGIA E INTELIGÊNCIA
  Performance
    Resumo
    Publicações
    Mídia
  Mercado
  Economia

CONTROLE
  Trabalho
  Aprovações

CONFIGURAÇÕES
  Configurações
```

No `Dashboard`, `Strategy`, `Analytics`, `Materialização`, `ERP`, generic `Integrações`, global search or permanent global right rail is introduced.

### 2.3 Executable proof behavior

The HTML candidate intentionally proves these B00 laws:

- navigation changes only page identity / contextual host; it does not create Product state;
- Organization switching clears Marketplace Installation context;
- organization-wide destinations do not show an ambient marketplace account;
- exact-required destinations show an explicit missing-account state until a Marketplace Installation is selected;
- all-or-exact destinations may expose `all` only when the page mode admits it;
- no first/default Marketplace Installation is silently selected;
- prototype-only scenarios exercise no-access and stale/deep-linked Organization blocking;
- mobile viewport transforms navigation into a drawer without changing IA order or context semantics.

All page-owned content is represented only by a neutral placeholder. B00 therefore does not decide Overview cards, table layouts, KPI hierarchy, form patterns or any later block structure.

## 3. Operator adjudication gate

B00 may become `LOCKED` only after the operator opens/interacts with the HTML artifact and explicitly adjudicates the rendered structure.

Possible dispositions:

```text
LOCKED       explicit operator approval of rendered B00
CANDIDATE    changes requested or review pending
FINDING      material contradiction discovered
```

No dependent P8 block becomes baseline before B00 is `LOCKED`, unless the operator explicitly authorizes parallel progression under Method v2.1.
