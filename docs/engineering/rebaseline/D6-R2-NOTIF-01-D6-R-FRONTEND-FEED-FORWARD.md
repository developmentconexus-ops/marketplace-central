# NOTIF-01 D6-R — Frontend Feed-Forward

> **Status:** OPERATOR-APPROVED DESIGN + WRITTEN SPEC / P8 B00-R2 RENDERED CANDIDATE / VISUAL ADJUDICATION REQUIRED
> **Operator approval:** 2026-08-23 — approved the Global-Maximum structural direction and subsequently approved this written spec before B00-R2 artifact authoring
> **Parent authority:** [D5-R4 Canonical Product OAD Wire Proof](D6-R2-NOTIF-01-D5-R4-OAD-WIRE-PROOF.md) + [Notification Architecture Design](D6-R2-NOTIF-01-NOTIFICATION-ARCHITECTURE-DESIGN.md) + [Frontend Product Experience Planning Method v2.1](../../development/frontend-product-experience-planning-method.md)
> **Existing frontend authority:** [D6-R2 P8 Block Ledger](D6-R2-P8-BLOCK-LEDGER.md) — B00/B01 operator-`LOCKED`; B00-R2 rendered candidate; B10 suspended
> **Product wire:** 104 Product operations · 31 ordinary Permissions · Principal kinds H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Purpose and boundary

D6-R feeds the accepted Personal Notifications Product authority forward into frontend experience authority without reopening Product semantics, global IA or runtime mechanics.

It decides the candidate human structure for:

```text
B00-R2 bounded topbar utility-slot reopen
→ bell / unread-presence awareness
→ bounded recent-Inbox preview
→ full personal Inbox utility route
→ Configurações > Notificações routing administration
```

It does **not**:

- change the global sidebar groups or Organization workspace model;
- add Product operations or Permissions;
- select React implementation details beyond accepted D6 architecture;
- select D7 PostgreSQL/River/SSE mechanics;
- add visual-design authority such as palette, typography, iconography, radius, shadow or branding polish;
- set any P8 block to `LOCKED` before rendered HTML inspection and explicit operator adjudication.

## 2. P0–P3 bounded authority result

### 2.1 Human jobs

**Personal awareness user**

```text
When a material Product occurrence happens while I may be elsewhere in MPC,
I need to notice it, understand enough context to triage it and continue to the current source truth,
so that important developments do not depend on me already being inside the source workspace.
```

**Notification routing administrator**

```text
When my Organization needs specific humans to receive specific operational awareness families,
I need to see the fixed Product-defined routing slots and configure exact human recipients,
so that responsibility is explicit without turning ordinary Permission into recipient selection.
```

### 2.2 Accepted wire coverage

| Human need | Owner | Read/write authority | Frontend home |
| --- | --- | --- | --- |
| unread presence / recent awareness / full Inbox | PersonalNotifications | `ListMyNotifications` | topbar bell + preview + full Inbox |
| read/unread + active/archived awareness | PersonalNotifications | `UpdateMyNotificationAwarenessState` | Inbox item controls |
| current routing state | PersonalNotifications | `ListNotificationRoutes` | Configurações > Notificações |
| human recipient discovery | IdentityAccess | `ListNotificationRouteRecipientCandidates` | routing editor selector |
| configure/unconfigure routing | PersonalNotifications | `SetNotificationRoute` | routing editor |

No orphan NOTIF-01 human operation remains.

## 3. IA decision

### 3.1 Global IA remains locked

The existing B00 sidebar remains unchanged. Notifications does not become a new global navigation mass.

```text
VISÃO GERAL
OFERTA
OPERAÇÃO
ESTRATÉGIA E INTELIGÊNCIA
CONTROLE
CONFIGURAÇÕES
```

`Configurações > Notificações` is a local Settings destination under the already-locked `Configurações` mass.

The full personal Inbox is an Organization-scoped **utility route reached primarily from the topbar bell**, not a sidebar destination.

### 3.2 Candidate route identities

Exact frontend spelling remains subject to P8/P9 route binding, but the approved semantic homes are:

```text
personal Inbox utility route
  e.g. /org/:organizationId/notificacoes

notification routing settings
  e.g. /org/:organizationId/configuracoes/notificacoes
```

Neither route creates Product API authority.

## 4. Alternatives challenged

### A — topbar bell + anchored preview + full Inbox utility route + inline Settings editor

**SELECTED — Global Maximum.**

It preserves awareness as a transversal utility while keeping durable triage/history in a route and routing configuration inside Settings. It fits current authority without adding API, Permission or parallel frontend truth.

### B — bell opens a drawer containing the entire Inbox

**REJECTED.**

A material filtered/paginated awareness surface should not depend on overlay lifetime. It weakens durable navigation, mobile recovery, URL/deep-link clarity and accessibility for the full Inbox job.

### C — Notifications becomes a new sidebar mass and routing gets one route per NotificationKind

**REJECTED.**

It overstates awareness as a primary product domain in the mental model and fragments ten small fixed routing slots into unnecessary navigation/route mass.

## 5. B00-R2 — bounded topbar utility-slot reopen

**Rendered P8 candidate:** [`qualification/d6-r2-wireframes/b00-r2-notifications.html`](../../../qualification/d6-r2-wireframes/b00-r2-notifications.html) — executable low-fidelity HTML; **not operator-`LOCKED` yet**.

The B00 physical/context shell remains binding:

```text
desktop sidebar ≈264 px
page header ≈64 px
Organization = only global workspace
page-local Marketplace Installation host
mobile drawer + stacked local context
```

Only one bounded utility slot is added to the page header.

Desktop structural candidate:

```text
┌──────────────────────────────────────────────────────────────┐
│ page title       page-local context / Installation      bell │
└──────────────────────────────────────────────────────────────┘
```

Mobile structural candidate:

```text
┌─────────────────────────────────┐
│ menu   page title          bell │
│        local context            │
└─────────────────────────────────┘
```

The bell:

- is Organization-scoped;
- is available to any authenticated human currently eligible/member in the active Organization because self-Inbox has no ordinary `notifications.read` Permission;
- does **not** depend on `notifications.manage`;
- closes its preview and invalidates/refetches incompatible Notification state on Organization switch;
- never introduces a second workspace/context selector.

## 6. Bell state and preview

### 6.1 Unread-presence state

The bell uses the accepted bounded probe:

```text
ListMyNotifications(
  archive_state = active,
  read_state = unread,
  limit = 1
)
```

Structural states are distinct:

```text
successful non-empty  → unread-known-present → dot/presence indicator
successful empty      → unread-known-empty   → no dot
request unavailable   → knowledge unavailable → visible unavailable treatment, never empty
```

No numeric unread count is inferred.

### 6.2 Preview pattern

Bell activation opens a **bounded anchored popover** on desktop/tablet and an equivalent bounded overlay treatment on narrow screens without turning the full Inbox into a drawer-only workflow.

Preview hierarchy:

```text
Notificações
→ small recent active-awareness list
→ each item: unread/read distinction + subject label + kind/outcome context + time
→ source continuation when currently authorized
→ "Ver todas" → full Inbox utility route
```

The preview is intentionally bounded; it is not an infinite-scroll Inbox and has no independent search/filter platform.

Opening a source **does not implicitly mark the Notification read**. Navigation and awareness mutation remain separate semantics; source navigation re-enters current source authorization.

## 7. Full personal Inbox

### 7.1 Representation

The Inbox uses a **structured vertical list**, not a table.

Reason: the primary job is heterogeneous awareness triage and continuation, not cross-row comparison of homogeneous fields. Notification kinds have different human meanings and two families have typed outcome atoms.

### 7.2 Primary controls

The Inbox exposes only admitted controls:

```text
archive lens:  Ativas | Arquivadas
read lens:     Todas | Não lidas | Lidas
kind filter:   bounded accepted NotificationKind selection
pagination:    limit/cursor behavior through normal UI continuation
```

No baseline:

```text
text search
unread count
total_count
mark-all-read
bulk archive
severity/priority
saved views
source-family filter DSL
```

### 7.3 Item hierarchy

Each list item may show only Notification-owned projection plus presentation-safe human language:

```text
read/unread distinction
subject_display_label
NotificationKind in user language
F02 offering_async_result_outcome when applicable
F14 authorization_decision_outcome when applicable
source_occurred_at
created_at when delayed materialization distinction matters
active/archived state
```

Material controls:

```text
Abrir origem
Marcar como lida / Marcar como não lida
Arquivar / Restaurar
```

`Abrir origem` never mutates source or Notification automatically. Read/archive controls invoke only `UpdateMyNotificationAwarenessState` with stale-write protection.

### 7.4 Material failure states

The Inbox must distinguish:

```text
loading / revalidating
known populated
known empty
request unavailable
stale awareness write
source continuation denied/not-found under current authorization
```

Unavailable never renders as known-empty.

## 8. Configurações > Notificações

### 8.1 Access and placement

The routing-settings destination is visible/usable only when current access includes `notifications.manage`.

Hidden navigation reduces noise only; direct/server authorization remains authoritative.

### 8.2 Collection structure

Because the Product exposes exactly ten fixed ORG_ROUTED kinds, the leading structure is a **single fixed route list** rather than one page per kind.

Each row shows:

```text
human NotificationKind label
current state: Configurado | Sem configuração de destinatários
current recipient presentation when resolvable from current candidate identities
Editar action
```

`UNCONFIGURED` is presented as **Sem configuração de destinatários**, not a generic disabled switch.

### 8.3 Editing pattern

Editing expands the affected row inline.

```text
route row
→ recipient selector
→ current/draft selected recipients
→ Cancelar
→ Salvar
→ Remover configuração when currently configured
```

No mandatory modal/drawer is introduced. The stale/error state remains adjacent to the route being changed.

`Remover configuração` writes the accepted `UNCONFIGURED` desired state; it is not `CONFIGURED([])` and does not delete temporal route history.

## 9. Recipient discovery and save-time eligibility

`ListNotificationRouteRecipientCandidates` returns only current eligible human:

```text
principal_id
display_name
```

It deliberately does not expose role/Permission/IdP/access-administration state and does not pre-authorize the route write.

Therefore the UX must tolerate:

```text
candidate listed
→ candidate becomes ineligible, leaves Organization or lacks kind-specific source-read eligibility
→ SetNotificationRoute rejects current submitted recipient
```

The editor preserves the local draft, keeps the affected route open and presents the rejected recipient selection as needing correction. The frontend must not infer or disclose the underlying Permission name merely to explain the failure.

The accepted validation-error grammar (`errors[].location + detail`) is sufficient for the server to identify a rejected submitted selection without adding a new Product operation.

## 10. Historical configured recipient no longer in current candidates

A current `NotificationRouteView` may contain a `recipient_principal_id` that is absent from the current candidate collection because eligibility continuity was broken.

The frontend must not:

```text
show the opaque PrincipalID as user-facing identity
invent a display name
call access.read as a hidden prerequisite
silently reactivate an old binding
```

Candidate structural treatment:

```text
Destinatário não elegível · configuração anterior
```

The row remains visibly inconsistent/resolvable rather than pretending the old binding is current. A new explicit routing decision is required to establish valid current recipients.

## 11. P5 delta

The approved D6-R design adds frontend homes only; Product surface remains unchanged.

```text
G00-E  topbar Notification utility slot
U01    bounded recent-Inbox preview popover
R128   full personal Inbox utility route
R129   Configurações > Notificações routing page
```

Consequences:

```text
new global sidebar destinations   0
new Product operations            0
new ordinary Permissions          0
new content routes/pages          2
new material popover surfaces     1
B00 reopen                         utility slot only
```

These identifiers are D6-R planning labels, not Product identities.

## 12. P7 data-feasibility line

### Bell / preview / Inbox

```text
fields/summaries         PRESENT-IN-AUTHORITY
identity source          PRESENT-IN-AUTHORITY
read/archive write       PRESENT-IN-AUTHORITY
bounded kind filtering   PRESENT-IN-AUTHORITY
unread presence          PRESENT-IN-AUTHORITY via limit=1 probe
numeric count            NOT ADMITTED / NOT REQUIRED
```

### Routing Settings

```text
10 fixed route slots           PRESENT-IN-AUTHORITY
current configured IDs         PRESENT-IN-AUTHORITY
human labels for current candidates PRESENT-IN-AUTHORITY
configure/unconfigure          PRESENT-IN-AUTHORITY
stale-write precondition       PRESENT-IN-AUTHORITY
kind-specific eligibility write validation PRESENT-IN-AUTHORITY
permission disclosure          NOT ADMITTED / NOT REQUIRED
```

No backend falsifier remains at this stage.

## 13. Assumption register

| Assumption | Evidence level | Influences | P12 probe | Status |
| --- | --- | --- | --- | --- |
| Organizations have a practically enumerable number of current human routing candidates using cursor pagination without a search operation | inferred / current Product scale not yet measured here | R129 recipient selector | walkthrough with realistic member counts; verify candidate discovery remains usable without server-side search | OPEN |
| A bounded recent preview is enough for quick awareness while full triage belongs to R128 | strong cross-product reference + operator-approved direction | U01 | walkthrough unread-heavy scenario and mobile behavior | OPEN |

If candidate-directory scale is falsified, reopen the smallest IdentityAccess read contract. Do not add search/filter DSL preventively.

## 14. Accessibility and responsive structural laws

P8/P11 must preserve:

- bell reachable by keyboard with an accessible name and non-color-only unread/unavailable state;
- preview focus enters predictably and returns to the bell when dismissed;
- Escape dismisses the preview where appropriate;
- preview and Inbox items have meaningful heading/reading order;
- no essential item action is hover-only;
- inline Settings expansion preserves focus context and exposes validation adjacent to the affected recipient/route;
- narrow widths may stack item metadata/actions but cannot hide the only source continuation or awareness mutation path;
- no drag-only recipient selection;
- Organization switch closes incompatible preview/editor state before new Organization truth is presented.

## 15. Negative controls

D6-R fails if frontend planning introduces:

```text
Notifications as a new global sidebar mass
numeric unread count inferred from pagination
notification source data as client-owned truth
open-source = implicit mark-read mutation
notifications.manage as prerequisite for self Inbox
access.read as hidden prerequisite for routing settings
role/Permission disclosure in recipient selector
client-side route authorization from candidate presence
CONFIGURED([]) or generic disabled switch
mark-all-read / bulk archive
text-search/filter platform by symmetry
mandatory drawer as the only full-Inbox realization
modal-only route editor that hides stale/error recovery
silent resurrection of historically ineligible bindings
new Product API merely to simplify one screen
```

## 16. P8 block order after written-spec approval

Rendered work is sequential and operator-reviewed:

```text
B00-R2  topbar bell + preview integration       RENDERED CANDIDATE
  ↓ operator visual adjudication / LOCK or revision
B11     full personal Inbox                    NOT RENDERED
  ↓ operator visual adjudication / LOCK or revision
B12     Configurações > Notificações           NOT RENDERED
  ↓ operator visual adjudication / LOCK or revision
P9 exact Screen Contracts / bidirectional wire trace
```

P8 artifacts are executable low-fidelity HTML, not images or visual-design comps.

No later block becomes baseline before the current material block is operator-`LOCKED`, unless the operator explicitly authorizes parallel candidate work.

## 17. Gate

```text
D5-R4 canonical OAD                      PROVED / CANONICAL 104/31
D6-R structural direction                OPERATOR-APPROVED
D6-R written spec                         OPERATOR-APPROVED
P8 B00-R2                                 RENDERED CANDIDATE / NOT LOCKED
B11 / B12                                 NOT RENDERED
D7-R                                      BLOCKED BY D6-R
D8-R                                      BLOCKED BY D7-R
Product implementation                    BLOCKED UNTIL D9
```

**Exact next action:** operator visually adjudicates only the rendered B00-R2 topbar bell + bounded recent-Inbox preview. Do not render B11/B12 as baseline, begin D7-R/D8-R, resume B10, merge PR #61 or implement Product code before this visual gate closes, unless the operator explicitly authorizes parallel candidate work.