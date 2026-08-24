# D6-R2 P8 — Structural Wireframe Block Ledger

> **Status:** OPEN / ACTIVE — B00 App Shell + global IA operator-`LOCKED`; B01 operator-`LOCKED`; B00-R2 Notification utility operator-`LOCKED`; B11 Personal Inbox operator-`LOCKED`; B12 Notification Routing Settings **RENDERED CANDIDATE / VISUAL ADJUDICATION REQUIRED**; OP-READ-01 RESOLVED; B10 SUSPENDED
> **Parent:** [D6-R2 Complete Frontend Realization Closure](D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md)
> **P4-R1 reopen:** [Global IA / Operational Mass Reopen](D6-R2-P4-R1-GLOBAL-IA-OPERATIONAL-MASS-REOPEN.md)
> **D5-R2 repair:** [Operational Read Projection Repair](D5-R2-OPERATIONAL-READ-PROJECTION-REPAIR.md)
> **D8-R2 revalidation:** [GF-02 Operational Read Revalidation](D8-R2-OPERATIONAL-READ-REVALIDATION.md)
> **NOTIF-01 D6-R:** [Frontend Feed-Forward](D6-R2-NOTIF-01-D6-R-FRONTEND-FEED-FORWARD.md)
> **P5 input:** [Complete Screen / Material-Surface Inventory](D6-R2-P5-SCREEN-SURFACE-INVENTORY.md)
> **Method:** [Frontend Product Experience Planning Method v2.1](../../development/frontend-product-experience-planning-method.md)
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. P8 artifact law

P8 artifacts are **executable low-fidelity HTML structural prototypes**, not visual-design comps and not static images.

They decide:

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

**Operator adjudication:** `LOCKED` on 2026-08-22.  
**Artifact:** [`qualification/d6-r2-wireframes/b00-app-shell.html`](../../../qualification/d6-r2-wireframes/b00-app-shell.html)

Locked physical/context laws:

```text
desktop persistent sidebar ≈264 px
Organization = only global workspace
page header ≈64 px with page-local Installation host
page-owned content ≈24 px outer padding
tablet collapsible navigation
mobile drawer + stacked local context
```

Also locked:

- Organization switching clears Marketplace Installation context;
- organization-wide routes have no ambient marketplace account;
- exact-required routes block until one exact Installation is selected;
- no hidden/default Installation;
- no-access/stale Organization blocks explicitly;
- responsive transformation never changes context meaning.

### 2.1 Global IA — LOCKED

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

Notifications does not become another global sidebar mass.

### 2.2 B00-R2 — Notification utility slot — LOCKED

**Artifact:** [`qualification/d6-r2-wireframes/b00-r2-notifications.html`](../../../qualification/d6-r2-wireframes/b00-r2-notifications.html).  
**Structural verifier:** `scripts/verify-d6-r-b00-r2-wireframe.mjs`.  
**Operator visual adjudication:** **APPROVED 2026-08-23**.  
**Operator `LOCK`:** **YES**.

Locked structure:

```text
G00-E topbar Notification utility slot
→ Organization-scoped bell
→ unread-known-present / unread-known-empty / knowledge-unavailable
→ bounded U01 recent-Inbox preview
→ explicit source continuation separate from awareness mutation
→ B11 "Ver todas" continuation
```

Locked laws:

- no global sidebar destination;
- self-Inbox bell does not depend on `notifications.manage`;
- unread presence is a dot/presence state, never an inferred number;
- unavailable is not known-empty;
- `Abrir origem` does not implicitly mark read;
- awareness mutation changes only Notification state;
- Organization switch closes/invalidate incompatible Notification context;
- desktop and mobile preserve the B00 context model.

The reviewed HTML remains a `CANDIDATE` snapshot; the operator `LOCK` lives in this ledger so the reviewed evidence is not rewritten after adjudication.

### 2.3 B11 — Full Personal Inbox — LOCKED

**Artifact:** [`qualification/d6-r2-wireframes/b11-notifications-inbox.html`](../../../qualification/d6-r2-wireframes/b11-notifications-inbox.html).  
**Structural verifier:** `scripts/verify-d6-r-b11-inbox-wireframe.mjs`.  
**First complete GREEN candidate:** `90409e6bc7adb711c63a0ef7ed8f0b52ced761c8`.  
**Operator visual adjudication:** **APPROVED 2026-08-23**.  
**Operator `LOCK`:** **YES — R128 full personal Inbox structure**.

Locked structure:

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

Locked laws:

- structured list, not table/grid;
- all 14 accepted NotificationKinds are filterable;
- F02/F14 typed outcomes are present when applicable;
- no text search, unread count, `total_count`, mark-all-read, bulk archive, severity/priority or saved views;
- source navigation and awareness mutation remain separate;
- known-empty != unavailable;
- stale awareness write and current-source denied states are recoverable;
- Organization switch invalidates preview/cursor/transient Inbox state;
- mobile stacks actions without hiding material controls.

Executable proof:

```text
d6_r_b11_status=CANDIDATE
d6_r_b11_representation=STRUCTURED_LIST
d6_r_b11_notification_kinds=14/14
d6_r_b11_totals_bulk_search=FORBIDDEN
d6_r_b11_b12=NOT_RENDERED
d6_r_b11_wireframe=PASS
```

The proof reports the reviewed HTML snapshot as `CANDIDATE`; the operator `LOCK` is authority in this ledger.

### 2.4 B12 — Notification Routing Settings — RENDERED CANDIDATE / NOT LOCKED

**Approved structural direction:** [NOTIF-01 D6-R Frontend Feed-Forward](D6-R2-NOTIF-01-D6-R-FRONTEND-FEED-FORWARD.md).  
**Rendered artifact:** [`qualification/d6-r2-wireframes/b12-notification-routing-settings.html`](../../../qualification/d6-r2-wireframes/b12-notification-routing-settings.html).  
**Structural verifier:** `scripts/verify-d6-r-b12-routing-settings-wireframe.mjs`.  
**First complete GREEN candidate:** `ae18ba0a9aa0faa4ea866e035c312e08b1e9c111` — repository full gate PASS.  
**Operator visual adjudication:** **REQUIRED**.  
**Operator `LOCK`:** **NO**.

Rendered candidate structure:

```text
Configurações > Notificações
→ notifications.manage-conditioned local Settings destination
→ exactly 10 fixed ORG_ROUTED route rows
→ Configurado | Sem configuração de destinatários
→ inline row editor
→ IdentityAccess candidate directory: principal_id + display_name only
→ cursor continuation for candidates
→ SetNotificationRoute with If-Match
→ explicit save / cancel / remove configuration
```

Candidate laws shown in the HTML:

- Notifications remains inside the existing `Configurações` mass; no new global sidebar mass;
- `ListNotificationRoutes` owns current route presentation through PersonalNotifications;
- `ListNotificationRouteRecipientCandidates` remains an IdentityAccess purpose-bound discovery projection;
- candidate discovery is usability, **not** route-write authorization;
- selector exposes only `principal_id + display_name`; no `access.read`, roles, Permissions or IdP state;
- the Product-defined set is exactly 10 ORG_ROUTED kinds;
- `UNCONFIGURED` is shown as **Sem configuração de destinatários**;
- `CONFIGURED([])` and generic disabled-switch semantics are forbidden;
- editor expands inline in the affected row; no mandatory modal/drawer baseline;
- `SetNotificationRoute` carries `If-Match` and save-time server revalidation;
- save-time recipient rejection preserves the draft and identifies the selection needing correction without Permission disclosure;
- a historical recipient absent from current candidates is rendered only as **Destinatário não elegível · configuração anterior**, never as an opaque ID or invented identity;
- candidate continuation is cursor-based; no candidate-search API/platform is introduced preventively;
- unavailable route truth is distinct from ten unconfigured routes;
- Organization switch invalidates editor, candidate pages and route state before presenting the new Organization;
- mobile stacks route rows/editor/actions without changing authority meaning.

Executable proof:

```text
d6_r_b12_status=CANDIDATE
d6_r_b12_route_slots=10/10
d6_r_b12_editor=INLINE_ROW
d6_r_b12_permission=notifications.manage
d6_r_b12_access_read=FORBIDDEN
d6_r_b12_candidate_disclosure=MINIMAL
d6_r_b12_configured_empty=FORBIDDEN
d6_r_b12_wireframe=PASS
```

No final visual-design decision is implied by the grayscale HTML treatment.

## 3. B01 — Overview — LOCKED

**Artifact:** [`qualification/d6-r2-wireframes/b01-overview.html`](../../../qualification/d6-r2-wireframes/b01-overview.html).  
**Operator adjudication:** `LOCKED` on 2026-08-22.

Locked contextual priority:

```text
Work known + actionable  -> attention expands and leads
Work known-empty         -> attention collapses; never implies health
Work unknown/unavailable -> uncertainty remains visible
all cases                -> marketplace/account + Performance + Economics orientation remains visible
```

## 4. B10 — Preparation — SUSPENDED CANDIDATE

**P6:** `DERIVED`; P7 NOT TRIGGERED.  
**Artifact:** [`qualification/d6-r2-wireframes/b10-preparation.html`](../../../qualification/d6-r2-wireframes/b10-preparation.html).

Its internal search/list → exact-subject detail → readiness/correspondence → explicit reread pattern remains valid evidence, and its global placement remains `OFERTA > Preparação`. It stays suspended while NOTIF-01 closes D6-R → D7-R → D8-R.

## 5. Operational landing — CANDIDATE CONCEPT ONLY

Operator-approved concept remains the action-oriented hybrid cockpit:

```text
1. PRECISA DE ATENÇÃO
2. TRABALHO OPERACIONAL NORMAL
3. ACOMPANHAMENTO
4. SPECIALIST ENTRY POINTS
```

One global cross-owner Kanban remains rejected.

## 6. OP-READ-01 — RESOLVED

[D5-R2](D5-R2-OPERATIONAL-READ-PROJECTION-REPAIR.md) and [D8-R2](D8-R2-OPERATIONAL-READ-REVALIDATION.md) remain binding. Negative controls remain:

```text
no /operational-dashboard Product endpoint
no OperationalWorkflow owner
no cross-owner synthetic lifecycle
no operational_stage / next_action / priority / total_count
no global totals inferred from one paginated page
no N+1 detail fan-out as production baseline
```

## 7. Global-Maximum / YAGNI law

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
visual-design decisions:    NONE
```

### 8.3 B11

```text
rendered artifact reviewed: YES
operator disposition:       LOCKED
material changes requested: NONE
R128 full Inbox structure:  LOCKED
visual-design decisions:    NONE
```

### 8.4 B12

```text
rendered artifact reviewed: NO / PENDING OPERATOR
operator disposition:       PENDING
structural verifier:        PASS
D7-R:                       BLOCKED
visual-design decisions:    NONE
```

A later material user-model finding may reopen only the smallest affected authority; preference alone does not.

## 9. Exact next action

Operator visually adjudicates **only B12 — `Configurações > Notificações` routing Settings** from the rendered executable low-fidelity HTML. Valid dispositions are revision/finding or explicit operator `LOCK`.

Do **not** begin D7-R/D8-R, resume B10, merge PR #61 or implement Product code before this B12 visual gate closes, unless the operator explicitly authorizes parallel candidate work under the frontend method.
