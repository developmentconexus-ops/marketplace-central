# D6-R2 P8 — Structural Wireframe Block Ledger

> **Status:** OPEN / ACTIVE — B00 operator-`LOCKED`; B01 is NEXT
> **Parent:** [D6-R2 Complete Frontend Realization Closure](D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md)
> **P5 input:** [Complete Screen / Material-Surface Inventory](D6-R2-P5-SCREEN-SURFACE-INVENTORY.md)
> **Method:** [Frontend Product Experience Planning Method v2.1](../../development/frontend-product-experience-planning-method.md)
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. P8 artifact law for Marketplace Central

P8 wireframes are **executable low-fidelity HTML structural prototypes**, not visual-design comps and not static images.

They decide only material UX structure:

```text
shell / region placement
relative size / density class
navigation / context placement
reading / action order
material state placement
responsive transformation
interaction needed to prove the block
```

They do **not** decide final palette, typography, iconography, radius, shadows, illustration, branding polish or final component styling. Those belong to later visual-design handoff and may not silently alter locked UX structure.

The static styled image generated during the first B00 attempt is **NON-AUTHORITATIVE / DISCARDED** and is not P8 evidence.

## 2. B00 — App Shell + global IA

**Status:** `LOCKED` — operator approved after opening/interacting with the executable HTML on 2026-08-22.

**Artifact:** [`qualification/d6-r2-wireframes/b00-app-shell.html`](../../../qualification/d6-r2-wireframes/b00-app-shell.html)

**P6:** NOT TRIGGERED.  
**P7:** NOT TRIGGERED.

### 2.1 Locked structural baseline

Desktop:

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

Responsive:

```text
desktop  persistent sidebar
→ tablet narrower/collapsible navigation
→ mobile drawer + stacked page-local context
```

No responsive transition may remove Organization/Installation qualification or change navigation meaning.

### 2.2 Locked IA

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

No `Dashboard`, `Strategy`, `Analytics`, `Materialização`, `ERP`, generic `Integrações`, global search or permanent global right rail is admitted by B00.

### 2.3 Locked behavior

- navigation changes page/context only; it does not create Product state;
- Organization switching clears Marketplace Installation context;
- organization-wide destinations do not show an ambient marketplace account;
- exact-required destinations show an explicit missing-account state until an Installation is selected;
- all-or-exact destinations expose `all` only when the page mode admits it;
- no first/default Marketplace Installation is silently selected;
- no-access and stale/deep-linked Organization states block explicitly;
- mobile navigation becomes a drawer without changing IA order or context semantics.

All page-owned content remains outside B00. The lock therefore does **not** decide Overview cards, table layouts, KPI hierarchy, form patterns or any later block structure.

## 3. B00 operator adjudication record

```text
rendered artifact reviewed: YES
operator disposition:          LOCKED
material changes requested:    NONE
visual-design decisions locked: NONE
next dependent block allowed:  B01 Overview
```

This `LOCKED` status is operator authority only; it is not inferred from CI or assistant judgment.

## 4. B01 — Overview — NEXT

**Homes:** R01.  
**P6:** NOT TRIGGERED.  
**P7:** NOT TRIGGERED.

P5 already limits Overview to a small read-only composition that may expose, when separately permitted:

```text
Marketplace Installation posture
bounded economic summary
bounded Work preview
per-Installation Performance entry point
```

A04 remains OPEN: accepted authority does not yet prove which permitted signal deserves first visual priority. External reference research cannot decide organization-specific priority, so B01 must begin with the smallest operator walkthrough needed to resolve or explicitly carry that priority assumption.

### Exact B01 pre-render question

Before rendering the HTML candidate, establish what the user should notice first on opening `Visão geral`:

```text
A  operational attention / Work requiring action
B  marketplace/account operating posture
C  commercial/economic posture
D  Performance entry / recent channel result
E  no dominant signal; intentionally balanced orientation
```

The answer controls hierarchy only. It does not create new Product reads, hidden counts, cross-provider aggregation or `/dashboard` authority.

No B01 HTML is baseline until this A04 input is recorded and the resulting rendered candidate is separately operator-reviewed.