# D6-R2 P8 — Structural Wireframe Block Ledger

> **Status:** OPEN / ACTIVE — B00 App Shell + corrected global IA operator-`LOCKED`; B01 content operator-`LOCKED`; B00-R2 Notification utility operator-`LOCKED`; B11 Personal Inbox **RENDERED CANDIDATE / VISUAL ADJUDICATION REQUIRED**; OP-READ-01 RESOLVED; B10 SUSPENDED
> **Parent:** [D6-R2 Complete Frontend Realization Closure](D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md)
> **P4-R1 reopen:** [Global IA / Operational Mass Reopen](D6-R2-P4-R1-GLOBAL-IA-OPERATIONAL-MASS-REOPEN.md)
> **D5-R2 repair:** [Operational Read Projection Repair](D5-R2-OPERATIONAL-READ-PROJECTION-REPAIR.md)
> **D8-R2 revalidation:** [GF-02 Operational Read Revalidation](D8-R2-OPERATIONAL-READ-REVALIDATION.md)
> **NOTIF-01 D6-R:** [Frontend Feed-Forward](D6-R2-NOTIF-01-D6-R-FRONTEND-FEED-FORWARD.md)
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

### 2.2 B00-R2 — Notification utility-slot reopen — LOCKED

**Approved design:** [NOTIF-01 D6-R Frontend Feed-Forward](D6-R2-NOTIF-01-D6-R-FRONTEND-FEED-FORWARD.md).  
**Rendered artifact:** [`qualification/d6-r2-wireframes/b00-r2-notifications.html`](../../../qualification/d6-r2-wireframes/b00-r2-notifications.html).  
**Structural verifier:** `scripts/verify-d6-r-b00-r2-wireframe.mjs` — first GREEN artifact commit `c725b89a9f652cc3628eb158684b067ae80696d3` passed the repository full gate.  
**Operator visual adjudication:** **APPROVED 2026-08-23**.  
**Operator `LOCK`:** **YES — bounded topbar bell + U01 preview structure**.

The locked B00-R2 structure preserves the already-locked B00 sidebar, Organization workspace and page-local Installation model and adds only:

```text
G00-E topbar Notification utility slot
→ Organization-scoped bell
→ unread-known-present / unread-known-empty / knowledge-unavailable distinction
→ bounded U01 recent-Inbox preview
→ explicit source continuation separate from awareness mutation
→ B11 "Ver todas" continuation
```

Locked structural laws:

- no new global sidebar destination;
- self-Inbox bell does not depend on `notifications.manage`;
- unread presence is a dot/presence state, never a numeric count inferred from pagination;
- request unavailable is visibly distinct from known-empty;
- `Abrir origem` does not implicitly mark read;
- explicit `Marcar como lida` changes only Notification awareness state;
- Organization switch closes incompatible preview state before new Organization awareness is presented;
- desktop preserves title + page-local context + bell; mobile preserves menu + title + bell with local context below;
- B11 full Inbox and B12 routing Settings are separate blocks and were not locked by the B00-R2 adjudication.

No final visual-design decision is implied by the grayscale HTML treatment.

### 2.3 B11 — Full Personal Inbox — RENDERED CANDIDATE / NOT LOCKED

**Approved structural direction:** [NOTIF-01 D6-R Frontend Feed-Forward](D6-R2-NOTIF-01-D6-R-FRONTEND-FEED-FORWARD.md).  
**Rendered artifact:** [`qualification/d6-r2-wireframes/b11-notifications-inbox.html`](../../../qualification/d6-r2-wireframes/b11-notifications-inbox.html).  
**Structural verifier:** `scripts/verify-d6-r-b11-inbox-wireframe.mjs`.  
**First complete GREEN candidate:** `90409e6bc7adb711c63a0ef7ed8f0b52ced761c8` — repository full gate PASS.  
**Operator visual adjudication:** **REQUIRED**.  
**Operator `LOCK`:** **NO — not yet adjudicated from rendered HTML**.

The rendered candidate inherits B00/B00-R2 and realizes only R128 full personal Inbox:

```text
Organization-scoped utility route outside sidebar
→ structured vertical awareness list
→ Ativas | Arquivadas
→ Todas | Não lidas | Lidas
→ closed NotificationKind filter
→ cursor continuation
→ explicit source continuation
→ explicit read/unread + archive/restore awareness mutation
```

Candidate structural laws shown in the HTML:

- structured list, not table/grid;
- all fourteen accepted NotificationKinds are filterable; no source-family/filter DSL;
- F02 and F14 typed result examples are statically inspectable and visibly rendered;
- no text search, unread count, `total_count`, mark-all-read, bulk archive, severity/priority or saved views;
- `Abrir origem` remains separate from awareness mutation and simulates current source re-authorization;
- read/unread and archive/restore mutate only Notification awareness state and preserve stale-write recovery;
- known-empty is distinct from request-unavailable;
- stale awareness write and source-access-denied remain explicit recoverable states;
- Organization switch invalidates preview, cursor and transient Inbox state before presenting the new Organization;
- narrow widths stack item actions without hiding source continuation or awareness controls;
- B12 `Configurações > Notificações` is **not rendered**.

Executable proof reports:

```text
d6_r_b11_status=CANDIDATE
d6_r_b11_representation=STRUCTURED_LIST
d6_r_b11_notification_kinds=14/14
d6_r_b11_totals_bulk_search=FORBIDDEN
d6_r_b11_b12=NOT_RENDERED
d6_r_b11_wireframe=PASS
```

No final visual-design decision is implied by the grayscale HTML treatment.

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

B01 inherits the locked B00 global navigation. Its content/state lock is not reopened by IA-01, B00-R2 or B11.

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

Its global placement is now locked as `OFERTA > Preparação`. B10 remains suspended while the accepted NOTIF-01 dependency chain completes D6-R → D7-R → D8-R; the immediate frontend gate is operator adjudication of B11.

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

[D5-R2](D5-R2-OPERATIONAL-READ-PROJECTION-REPAIR.md) enriches only existing owner-local operational list projections/filters and preserved the then-current 99-operation / 30-Permission Product surface without changing owners or write surface.

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

## 8. Operator adjudication records

### 8.1 B00-R1

```text
rendered artifact reviewed: YES
operator disposition:       LOCKED
material changes requested: NONE
physical/context shell:     LOCKED / unchanged
corrected global IA:        LOCKED
visual-design decisions:    NONE
```

### 8.2 B00-R2

```text
rendered artifact reviewed: YES
operator disposition:       LOCKED
material changes requested: NONE
bounded utility slot:       LOCKED
bell + U01 preview:         LOCKED
B11 / B12:                  NOT LOCKED BY THIS GATE
visual-design decisions:    NONE
```

### 8.3 B11

```text
rendered artifact reviewed: NO / PENDING OPERATOR
operator disposition:       PENDING
structural verifier:        PASS
B12:                        NOT RENDERED
visual-design decisions:    NONE
```

A later material user-model finding may reopen only the smallest affected authority; preference alone does not.

## 9. Exact next action

Operator visually adjudicates **only B11 — full personal Inbox** from the rendered executable low-fidelity HTML. Valid dispositions are revision/finding or explicit operator `LOCK` of the B11 structure.

Do **not** render B12 as baseline, begin D7-R/D8-R, resume B10, merge PR #61 or implement Product code before this B11 visual gate closes, unless the operator explicitly authorizes parallel candidate work under the frontend method.