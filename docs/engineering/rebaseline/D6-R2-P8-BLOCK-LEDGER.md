# D6-R2 P8 — Structural Wireframe Block Ledger

> **Status:** OPEN / ACTIVE — B00 + B01 operator-`LOCKED`; B10 bounded P6 is NEXT
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
```

This `LOCKED` status is operator authority only; CI verifies repository coherence but does not grant the lock.

## 4. A04 — Overview priority — OPERATOR-ADJUDICATED

The original static choices A/B/C/D/E were challenged before B01 rendering. The operator approved the following **contextual-priority** law instead of any always-dominant signal:

```text
if Work is known + actionable
→ operational-attention region expands and leads the page

if Work is known-empty
→ operational-attention region collapses to a short explicit known-empty state
→ it MUST NOT say or imply “everything is healthy”

if Work knowledge is unavailable/unknown
→ the region remains visible as unavailable/unknown
→ it MUST NOT infer zero, empty or healthy

in all three cases
→ a compact marketplace/account + Performance + Economics orientation remains visible
```

This is **not** a new health model. The Overview remains a read-only composition of admitted owner reads.

### 4.1 Review findings incorporated

1. A pure Work-first page was rejected because an empty Work result could become a false green state if the underlying detection/evidence path is incomplete.
2. Slow commercial degradation may never produce discrete Work, so the bounded Performance/Economics orientation is a safety net, not decoration.
3. No alert-triage lifecycle is invented. Current Work authority does **not** admit generic `resolve`, `dismiss` or `close`; B01 may navigate to Work only.
4. No “verified X minutes ago” or collector-health freshness is fabricated. Such freshness may be displayed only when the owning Product authority actually exposes it.
5. No health score, hidden total, cross-provider KPI or synthetic dashboard authority is admitted.

A04 is therefore **resolved for B01 hierarchy by operator adjudication**. It does not create a Product capability and does not reopen D0–D8.

## 5. B01 — Overview

**Status:** `LOCKED` — operator approved after opening/interacting with the executable HTML on 2026-08-22.  
**Homes:** R01.  
**Artifact:** [`qualification/d6-r2-wireframes/b01-overview.html`](../../../qualification/d6-r2-wireframes/b01-overview.html)  
**P6:** NOT TRIGGERED.  
**P7:** NOT TRIGGERED.

### 5.1 Locked structure

B01 inherits the entire locked B00 shell unchanged. Its page-owned content is:

```text
Visão geral
├─ contextual operational-attention region
│  ├─ ACTIONABLE            -> expanded preview + navigation to Trabalho
│  ├─ KNOWN_EMPTY           -> collapsed explicit known-empty state
│  └─ KNOWLEDGE_UNAVAILABLE -> visible unknown/unavailable state
└─ balanced orientation — always present
   ├─ Marketplace / account posture             [MarketplacePortfolio]
   ├─ per-Installation Performance entry/result [MarketplacePerformanceIntelligence]
   └─ bounded economic orientation               [CommercialEconomics]
```

The Work preview is bounded evidence, not a global inferred count. It exposes no screen-shaped mutation and cannot resolve underlying source truth.

### 5.2 Independent-owner law

Each Overview region remains permission/evidence conditioned independently. Missing one region does not authorize the host page to infer that owner's truth.

```text
not permitted != known empty
unknown/unavailable != zero
Work empty != healthy operation
Performance absent != zero performance
Economics absent != zero/healthy economics
```

### 5.3 Responsive law

The B00 responsive shell remains locked. B01 content may stack at narrower widths, but reading priority stays:

```text
contextual attention state
→ marketplace/account orientation
→ Performance orientation
→ Economics orientation
```

When Work is `KNOWN_EMPTY`, its collapsed band does not consume dominant vertical space. When Work is `KNOWLEDGE_UNAVAILABLE`, its uncertainty remains visible rather than disappearing.

### 5.4 Locked negative controls

```text
“Tudo certo” inferred from empty Work
resolve/dismiss/close Work actions
frontend-generated collector freshness
health score
/dashboard authority
cross-provider aggregate KPI
hidden global totals from one paginated list
owner write authority acquired by composition
```

## 6. B01 operator adjudication record

```text
rendered artifact reviewed: YES
operator disposition:          LOCKED
material changes requested:    NONE
visual-design decisions locked: NONE
```

The operator lock covers B01 structure and state hierarchy only. It does not choose final metrics, typography, colors, visual component treatment or a new Product capability.

## 7. B10 — Preparation — P6 NEXT

**Homes:** R10.  
**Status:** `NOT RENDERED` — bounded P6 reference study must complete first.  
**P6:** TRIGGERED.  
**P7:** CONDITIONAL after P6.

B10 must support one complete human job without hidden source authority:

```text
exact Organization + Marketplace Installation
→ search admitted source products
→ select exact SourceInstance + native product key
→ inspect readiness + publication requirements
→ resolve/clear correspondence when needed and authorized
→ re-read readiness
→ continue to ListingIntent only when the current state permits
```

The real structural uncertainty is whether this search-first, multi-source, readiness-and-resolution task is best represented as:

```text
split-view / master-detail
progressive detail inside one workspace
list/search → dedicated detail state
```

P6 must study mature products solving analogous catalog-diagnostics/readiness tasks and extract task-pattern evidence rather than visual fashion. If two or more materially credible structures remain after the bounded study, P7 opens; otherwise record why one pattern is sufficient.

Do **not** render B10 HTML until P6 has a durable result and any required P7 adjudication is complete.