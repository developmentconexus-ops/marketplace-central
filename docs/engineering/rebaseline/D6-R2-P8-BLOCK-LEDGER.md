# D6-R2 P8 — Structural Wireframe Block Ledger

> **Status:** OPEN / ACTIVE — B00 App Shell + corrected global IA operator-`LOCKED`; B01 content operator-`LOCKED`; OP-READ-01 RESOLVED; B10 SUSPENDED pending current upstream review
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

## 2. B00 — App Shell + global IA — LOCKED

**Operator adjudication:** `LOCKED` on 2026-08-22 after review of the corrected executable HTML.  
**Artifact:** [`qualification/d6-r2-wireframes/b00-app-shell.html`](../../../qualification/d6-r2-wireframes/b00-app-shell.html)

The following B00 laws are locked:

```text
desktop persistent sidebar ≈264 px
Organization is the only global workspace
page header ≈64 px with page-local Installation host
page-owned content ≈24 px outer padding
→ tablet collapsible navigation
→ mobile drawer + stacked local context
```

Also locked:

- Organization switching clears Marketplace Installation context;
- organization-wide routes have no ambient marketplace account;
- exact-required routes block until one exact Installation is selected;
- no hidden/default Installation;
- no-access/stale Organization blocks explicitly;
- responsive transformation never changes context meaning.

### 2.1 IA-01 corrected grouping — LOCKED

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

Existing technical route identities remain where possible; user-facing `Anúncios` may continue to use `/publicacoes/*`, and `Preços` may continue to use `/publicacoes/precos`. The organization-wide operational landing remains the one new frontend route identity candidate, e.g. `/org/:organizationId/operacao`; exact route spelling is frontend navigation authority, not Product API authority.

B00 does **not** decide the internal content of Visão operacional. It locks only the destination's place in global IA and its organization-wide context.

## 3. B01 — Overview — LOCKED

**Content/state hierarchy:** `LOCKED` — operator approved after executable HTML review on 2026-08-22.  
**Artifact:** [`qualification/d6-r2-wireframes/b01-overview.html`](../../../qualification/d6-r2-wireframes/b01-overview.html)

A04 remains locked as contextual priority:

```text
Work known + actionable  -> attention expands and leads
Work known-empty         -> attention collapses; never implies health
Work unknown/unavailable -> uncertainty remains visible
all cases                -> marketplace/account + Performance + Economics orientation remains visible
```

B01 inherits the locked B00 global navigation. Its content/state lock is not reopened by IA-01.

## 4. B10 — Preparation — SUSPENDED CANDIDATE

**Internal pattern:** P6 `DERIVED`; P7 NOT TRIGGERED.  
**P6 evidence:** [B10 Preparation Reference Study](D6-R2-P6-B10-PREPARATION-REFERENCE-STUDY.md).  
**Rendered artifact:** [`qualification/d6-r2-wireframes/b10-preparation.html`](../../../qualification/d6-r2-wireframes/b10-preparation.html)

The B10 internal pattern remains valid evidence:

```text
search/list triage
→ selected exact-subject detail
→ readiness / requirements / correspondence
→ explicit reread after consequential correspondence effect
→ continuation only from admitted current state
```

Its global placement is now locked as `OFERTA > Preparação`. B10 remains suspended only because the operator requested an upstream Notification architecture evaluation before frontend progression resumes.

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

The operational landing may later compose repaired owner-local reads, but it acquires no business/write authority.

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

## 8. B00-R1 operator adjudication record

```text
rendered artifact reviewed: YES
operator disposition:       LOCKED
material changes requested: NONE
physical/context shell:     LOCKED / unchanged
corrected global IA:        LOCKED
visual-design decisions:    NONE
```

This closes IA-01. A later material user-model finding may reopen only the smallest affected authority; preference alone does not.

## 9. Exact next action

Pause dependent frontend progression while the operator-requested **Notification architecture assessment** is performed against current D0–D8 + D5-R2/D8-R2 authority. Do not introduce Notification Product semantics, operations, owner, event contract, runtime mechanism or frontend surface until that assessment is operator-adjudicated.
