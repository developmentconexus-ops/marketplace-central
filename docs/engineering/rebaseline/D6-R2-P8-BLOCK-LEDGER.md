# D6-R2 P8 — Structural Wireframe Block Ledger

> **Status:** OPEN / ACTIVE — B00 + B01 operator-`LOCKED`; B10 P6 DERIVED / P7 NOT TRIGGERED; B10 HTML `CANDIDATE` rendered
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
**P6/P7:** NOT TRIGGERED.

Locked baseline:

```text
desktop persistent sidebar ≈264 px
+ Organization as only global workspace
+ page header ≈64 px with page-local Installation host
+ page-owned content ≈24 px outer padding
→ tablet collapsible navigation
→ mobile drawer + stacked local context
```

Locked IA remains:

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

Locked behavior:

- Organization switching clears Marketplace Installation context;
- organization-wide routes have no ambient marketplace account;
- exact-required routes block until one exact Installation is selected;
- all-or-exact routes expose `all` only where admitted;
- no hidden/default Installation;
- no-access/stale Organization blocks explicitly;
- responsive transformation never changes IA/context meaning;
- no Dashboard/Strategy/Analytics/Materialização/ERP/global search/permanent right rail.

The lock is structural only; it chooses no visual-design system.

## 3. B01 — Overview

**Status:** `LOCKED` — operator approved after executable HTML review on 2026-08-22.  
**Homes:** R01.  
**Artifact:** [`qualification/d6-r2-wireframes/b01-overview.html`](../../../qualification/d6-r2-wireframes/b01-overview.html)  
**P6/P7:** NOT TRIGGERED.

### 3.1 A04 — operator-adjudicated contextual priority

```text
Work known + actionable
→ attention expands and leads

Work known-empty
→ attention collapses
→ never implies healthy operation

Work unknown/unavailable
→ uncertainty remains visible
→ never infers zero/empty/healthy

all cases
→ marketplace/account + Performance + Economics orientation remains visible
```

Locked B01 negative controls:

```text
“Tudo certo” from empty Work
resolve/dismiss/close Work actions
frontend-generated collector freshness
health score
/dashboard authority
cross-provider aggregate KPI
hidden global totals
write authority acquired by composition
```

## 4. B10 — Preparation

**Homes:** R10.  
**Status:** `CANDIDATE` — executable HTML rendered; awaiting operator review.  
**Artifact:** [`qualification/d6-r2-wireframes/b10-preparation.html`](../../../qualification/d6-r2-wireframes/b10-preparation.html)  
**P6:** [DERIVED](D6-R2-P6-B10-PREPARATION-REFERENCE-STUDY.md).  
**P7:** NOT TRIGGERED.

P6 reference evidence converged on **search/list triage → selected exact-subject detail**. Persistent split-view and full material inline expansion are rejected as baseline.

### 4.1 Candidate structural modes

```text
SEARCH MODE
├─ exact Marketplace Installation in B00 local context
├─ search input
├─ optional explicit SourceInstance narrowing
└─ structured results
   ├─ source-product presentation
   ├─ SourceInstance + native key
   ├─ compact readiness state
   ├─ requirements/correspondence signal when known
   └─ open preparation

SELECTED-SUBJECT MODE
├─ back to search/results; search context is navigation state
├─ exact SourceInstance + native key + Installation header
├─ readiness / knowledge state
├─ publication requirements
├─ correspondence state
│  └─ resolve / clear only when Product admits it
├─ explicit re-read after consequential correspondence effect
└─ continue to ListingIntent only when current state permits
```

The selected subject remains inside the accepted `Preparação` workspace. P9 later binds the exact URL/search carriers; B10 does not invent Product identity.

### 4.2 Executable scenarios

The HTML candidate exposes four structural scenarios:

```text
READY
→ requirements known/attended
→ correspondence resolved
→ ListingIntent continuation available

MISSING_REQUIREMENTS
→ readiness partial/known
→ missing requirements remain explicit
→ continuation blocked

CORRESPONDENCE_NEEDED
→ no candidate auto-selected
→ admitted resolve/clear only
→ after consequential effect, continuation remains blocked
→ explicit re-read/revalidation required before assuming convergence

KNOWLEDGE_UNAVAILABLE
→ unknown/unavailable remains explicit
→ requirements do not become an empty list
→ correspondence is not inferred
→ continuation blocked
```

The results also include two different SourceInstances to visibly prove that omission of the source filter is multi-source search, not a hidden source default.

### 4.3 Responsive law

B10 inherits B00. On narrow widths, the result collection becomes stacked rows and selected-subject regions become one column. The interface may change arrangement but never drops Installation, SourceInstance or native-key qualification.

### 4.4 B10 negative controls

```text
hidden/default SourceInstance
first search result silently selected
Marketplace Installation used as ambient authority
source product edited as MPC master data
provider write exposed from Readiness
frontend-invented readiness percentage/health score
known-empty collapsed into unknown/unavailable/unsupported
correspondence represented as generic product editing
blind success after correspondence mutation
bulk preparation/actions invented without authority/evidence
permanent split-view promoted without evidence
```

## 5. B10 operator adjudication gate

B10 may become `LOCKED` only after the operator opens/interacts with the HTML candidate and explicitly adjudicates its structure.

```text
LOCKED       rendered structure approved
CANDIDATE    changes requested / review pending
FINDING      material contradiction discovered
```

Do not render B20 Publications core as baseline before B10 is `LOCKED` unless the operator explicitly authorizes parallel progression.
