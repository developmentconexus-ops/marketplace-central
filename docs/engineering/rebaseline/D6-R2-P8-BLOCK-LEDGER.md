# D6-R2 P8 — Structural Wireframe Block Ledger

> **Status:** OPEN / ACTIVE — B00 physical shell + B01 content operator-`LOCKED`; B00 global IA REOPENED by IA-01; OP-READ-01 RESOLVED; corrected B00 render NEXT; B10 SUSPENDED
> **Parent:** [D6-R2 Complete Frontend Realization Closure](D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md)
> **P4-R1 reopen:** [Global IA / Operational Mass Reopen](D6-R2-P4-R1-GLOBAL-IA-OPERATIONAL-MASS-REOPEN.md)
> **D5-R2 repair:** [Operational Read Projection Repair](D5-R2-OPERATIONAL-READ-PROJECTION-REPAIR.md)
> **D8-R2 revalidation:** [GF-02 Operational Read Revalidation](D8-R2-OPERATIONAL-READ-REVALIDATION.md)
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

They do **not** decide final palette, typography, iconography, radius, shadows, illustration, branding polish or final component styling.

## 2. B00 — App Shell + global IA — PARTIAL REOPEN / RERENDER NEXT

**Physical/context shell:** `LOCKED`.  
**Global IA grouping/labels:** `REOPENED` by operator-approved finding IA-01.  
**Artifact to correct:** [`qualification/d6-r2-wireframes/b00-app-shell.html`](../../../qualification/d6-r2-wireframes/b00-app-shell.html)

The following B00 laws remain locked:

```text
desktop persistent sidebar ≈264 px
Organization is the only global workspace
page header ≈64 px with page-local Installation host
page-owned content ≈24 px outer padding
→ tablet collapsible navigation
→ mobile drawer + stacked local context
```

Also preserved:

- Organization switching clears Marketplace Installation context;
- organization-wide routes have no ambient marketplace account;
- exact-required routes block until one exact Installation is selected;
- no hidden/default Installation;
- no-access/stale Organization blocks explicitly;
- responsive transformation never changes context meaning.

### 2.1 IA-01 grouping — operator-approved direction

```text
VISÃO GERAL
  Visão geral

OFERTA
  Preparação
  Anúncios
  Preços
  Disponibilidade

OPERAÇÃO
  Visão operacional
  Vendas
  Expedição
  Pós-venda

ESTRATÉGIA E INTELIGÊNCIA
  Performance
    Resumo
    Anúncios
    Mídia
  Mercado
  Economia

CONTROLE
  Trabalho
  Aprovações

CONFIGURAÇÕES
  Configurações
```

This grouping is operator-adjudicated direction but is not yet a restored B00 global-frame `LOCK`: it must be visibly re-rendered and separately reviewed.

Existing technical route identities remain where possible; user-facing `Anúncios` may continue to use `/publicacoes/*`, and `Preços` may continue to use `/publicacoes/precos`. The only candidate new route is the organization-wide operational landing `/org/:organizationId/operacao`.

## 3. B01 — Overview

**Content/state hierarchy:** `LOCKED` — operator approved after executable HTML review on 2026-08-22.  
**Artifact:** [`qualification/d6-r2-wireframes/b01-overview.html`](../../../qualification/d6-r2-wireframes/b01-overview.html)

A04 remains locked as contextual priority:

```text
Work known + actionable -> attention expands and leads
Work known-empty        -> attention collapses; never implies health
Work unknown/unavailable -> uncertainty remains visible
all cases               -> marketplace/account + Performance + Economics orientation remains visible
```

B01 will later inherit the corrected global navigation. Its content/state lock is not reopened by IA-01.

## 4. B10 — Preparation — SUSPENDED

**Internal pattern:** P6 `DERIVED`; P7 NOT TRIGGERED.  
**P6 evidence:** [B10 Preparation Reference Study](D6-R2-P6-B10-PREPARATION-REFERENCE-STUDY.md).  
**Rendered artifact:** [`qualification/d6-r2-wireframes/b10-preparation.html`](../../../qualification/d6-r2-wireframes/b10-preparation.html)  
**Current disposition:** `SUSPENDED CANDIDATE` — do not operator-lock until corrected global IA is restored.

The B10 internal pattern remains valid evidence:

```text
search/list triage
→ selected exact-subject detail
→ readiness / requirements / correspondence
→ explicit reread after consequential correspondence effect
→ continuation only from admitted current state
```

Its previous placement under `OPERAÇÕES` is superseded by `OFERTA > Preparação`.

## 5. Operational landing — CANDIDATE CONCEPT ONLY

Operator approved **cockpit hybrid oriented to action**, not one global Kanban:

```text
1. PRECISA DE ATENÇÃO
   explicit exception / Work / ambiguity / post-sale attention

2. TRABALHO OPERACIONAL NORMAL
   normal actionable execution

3. ACOMPANHAMENTO
   in-progress / dispatched / delivered without immediate action

4. SPECIALIST ENTRY POINTS
   Vendas · Expedição · Pós-venda · Trabalho
```

A global `Nova -> Faturar -> Separar -> Conferir -> Embalar -> Enviar` Kanban is rejected because no cross-owner Product lifecycle owns those columns. A Kanban-like view may later be studied inside Fulfillment/Expedição only.

## 6. OP-READ-01 — RESOLVED

[D5-R2](D5-R2-OPERATIONAL-READ-PROJECTION-REPAIR.md) enriches only existing owner-local operational list projections/filters and preserves 99 operations, 30 Permissions, H/A/S, owners and write surface.

[D8-R2](D8-R2-OPERATIONAL-READ-REVALIDATION.md) confirms GF-02 remains coherent: choreography, effect semantics, physical authority and source-domain boundaries are unchanged.

Binding negative controls remain:

```text
no /operational-dashboard Product endpoint
no OperationalWorkflow owner
no cross-owner synthetic lifecycle
no operational_stage / next_action / priority / total_count
no global totals inferred from one paginated page
no N+1 detail fan-out as production baseline
```

The operational landing may later compose the repaired owner-local reads, but it acquires no business/write authority.

## 7. Global-Maximum / YAGNI law

For every remaining block:

```text
human job + accepted authority
→ smallest coherent UX structure
→ exact owner/wire trace
→ bounded upstream repair only on material falsifier
→ no screen-shaped API or parallel frontend truth
→ no speculative platform capability
```

## 8. Exact next action

Re-render **B00 App Shell + global IA** only, preserving the locked physical/context/responsive shell and applying the operator-approved IA-01 navigation grouping. The corrected HTML remains `CANDIDATE` until the operator opens/interacts with it and explicitly restores the global-frame `LOCK`.

Do not render the Visão operacional or resume B10/B20 before corrected B00 is operator-`LOCKED`.
