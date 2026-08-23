# NOTIF-01 D5-R3 — Product Operation Admission Table

> **Status:** OPEN — FROZEN CANDIDATE / OPERATOR ADJUDICATION REQUIRED
> **Accepted parents:** D0-R + D1-R + D2-R + D2-R2 + D2-R3 + D2-R4 + D3-R; [D3-R1 Presentation Feed-Forward](D6-R2-NOTIF-01-D3-R1-PRESENTATION-FEED-FORWARD-REVALIDATION.md) = PASS
> **Parent D5 authority:** [D5 API](D5-API.md) + [D5-B2 Operation Admission Matrix](D5-B2-OPERATION-ADMISSION-MATRIX.md)
> **Scope:** freeze the smallest Product operation / client-class / Permission / safety surface before any canonical OpenAPI modification
> **Current Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Governing admission result

The accepted NOTIF-01 consumers require exactly four new Product operations and one new ordinary Permission.

```text
Personal Inbox
  ListMyNotifications
  UpdateMyNotificationAwarenessState

Organization routing settings
  ListNotificationRoutes
  SetNotificationRoute
```

No public Notification creation operation exists. No separate bell/count endpoint exists. No machine-style Inbox exists.

If this table is operator-ratified exactly as written, the later Product wire consequence is:

```text
99 + 4 = 103 Product operations
30 + 1 = 31 ordinary Permissions
Principal kinds remain H/A/S
```

Counts are consequences, never protected targets. **This document does not edit the OAD.**

---

## 2. Frozen candidate operation matrix

| operationId | Real Product consumer | Client | D3 class | Ordinary access | Admission |
| --- | --- | --- | --- | --- | --- |
| `ListMyNotifications` | topbar bell preview, unread-presence probe, full personal Inbox | H only | Q | exact authenticated human recipient + current Organization Membership/Product eligibility; no separate ordinary Permission | **ADMIT** |
| `UpdateMyNotificationAwarenessState` | mark one own Notification read/unread and active/archived | H only | C/update | same exact self-recipient + current Organization eligibility rule | **ADMIT** |
| `ListNotificationRoutes` | `Configurações > Notificações` routing administration | H only | Q | `notifications.manage` | **ADMIT** |
| `SetNotificationRoute` | configure or intentionally unconfigure one Product-defined ORG_ROUTED kind | H only | C/update | `notifications.manage` | **ADMIT** |

All four operations are Organization-scoped under the already-accepted D5 organization-path law. Exact path/method/schema spelling is later D5 wire work.

A/S do not receive a human Inbox or notification-settings UI capability by symmetry. Internal source-owner → Personal Notifications reactions remain D3/D7 internal behavior, not Product operations.

---

## 3. Self-Inbox access law

Baseline self Inbox does **not** introduce `notifications.read`.

The Product has already selected and materialized an exact human recipient. No proved Product 1.0 consumer requires a second ordinary Permission axis that selectively denies one eligible member access to awareness already addressed to that same member.

Therefore both Inbox operations require:

```text
authenticated Principal kind = H
+ exact path Organization
+ current Membership/Product-access eligibility in that Organization
+ Notification.recipient_principal_id == authenticated Principal
```

This capability is self-only. It does not authorize reading another Principal's Inbox.

Possession/readability of a Notification does not grant source access. Navigation through typed `source_ref` re-enters the source owner's current Product authorization.

Rejected baseline:

```text
notifications.read
admin/cross-Principal Inbox
search another Principal's Notifications
Organization-wide shared Inbox
A/S human-style Inbox
```

---

## 4. `ListMyNotifications` — frozen admission contract

### 4.1 Consumer-proven collection controls

Admit only:

```text
archive_state: active | archived
read_state: all | unread | read
limit
cursor
```

Default collection meaning:

```text
current non-superseded active awareness
ordered by Notification.created_at DESC
+ stable deterministic tie-breaker
```

Superseded Notifications remain durable correctness/history state but are not public current-Inbox collection members. No public superseded-history filter is admitted.

### 4.2 Bell unread-presence proof

No unread-count operation is required.

```text
ListMyNotifications(
  archive_state = active,
  read_state = unread,
  limit = 1
)
```

- successful empty collection → no current unread personal awareness;
- successful non-empty collection → unread presence.

A request failure/unavailability never means empty.

### 4.3 Minimum human-usable item projection

The collection representation must expose proportionately enough accepted Notification-owned state for the bell and full Inbox without source-per-row N+1 reads:

```text
notification_id
kind
subject_display_label
source_ref
source_occurred_at
created_at
read state
archive state
revision / later D5 concurrency carrier where applicable
```

`subject_display_label` is the accepted immutable D2-R4 presentation snapshot; it is not source truth or navigation authority.

`source_occurred_at` and `created_at` remain distinct so delayed materialization does not falsify chronology.

Exact schema names/encoding are wire work after this table is ratified.

### 4.4 Rejected collection surface

```text
exact unread numeric aggregate
total_count
text search
arbitrary kind/source filter DSL
separate bell-preview endpoint
GetNotification detail operation by symmetry
admin-wide/cross-Principal search
public superseded-history/debug collection
```

No `GetNotification` is admitted because the list item itself is the complete bounded Inbox representation and source detail belongs to the source owner.

---

## 5. `UpdateMyNotificationAwarenessState` — complete safety tuple

The operation changes only Personal Notifications-owned awareness state for one exact own Notification.

Desired semantic axes:

```text
read | unread
active | archived
```

They remain orthogonal. The operation does not acknowledge, resolve, authorize, dismiss or mutate the source.

Safety tuple:

```text
consequence class:
  non-consequential owner-local awareness mutation

idempotency:
  structural desired-state semantics
  no client idempotency key

concurrency / precondition:
  REQUIRED
```

Reason for precondition: a stale client must not overwrite a later read/archive mutation or a later D2-R2 supersession revision. A superseded Notification is no longer a current Inbox item; a stale mutation racing with supersession must fail/reconcile rather than resurrect visibility.

One desired-state update operation is admitted instead of four action-shaped endpoints (`markRead`, `markUnread`, `archive`, `unarchive`) because D2 already defines these as ordinary state axes of the same owner resource.

Rejected:

```text
delete Notification
mark-all-read
bulk mutation
source acknowledgement
source mutation/action
```

---

## 6. Routing administration and `notifications.manage`

`notifications.manage` is the one new ordinary Permission candidate.

It means:

> manage the Organization's Personal Notifications routing configuration for the fixed Product-defined ORG_ROUTED NotificationKinds.

It does **not** mean:

- receive those Notifications;
- read another user's Inbox;
- manage Membership/AccessRoles;
- define custom kinds;
- create a generic subscription/rules system.

Permission remains access authority only. It never selects a recipient by itself.

The existing AccessRole/Permission model can carry this additional ordinary Permission without changing identity or role semantics; exact role composition/default assignment is not invented by this D5 operation table.

---

## 7. `ListNotificationRoutes` — frozen admission contract

The operation lists the fixed Product-defined baseline ORG_ROUTED slots and their **current** desired state for the Organization.

Baseline slot set is exactly the ten accepted ORG_ROUTED kinds:

```text
MARKETPLACE_INSTALLATION_ATTENTION
AVAILABILITY_ATTENTION
ECONOMIC_RECONCILIATION_ATTENTION
NEW_MARKETPLACE_SALE
SALE_ATTENTION
MATERIALIZATION_ATTENTION
FULFILLMENT_ACTIONABLE
FULFILLMENT_ATTENTION
SHIPMENT_EXCEPTION
POST_SALE_ATTENTION
```

Each current slot is semantically:

```text
UNCONFIGURED
```

or:

```text
CONFIGURED
+ one-or-more exact human recipient_principal_ids
```

The representation includes owner-local current revision/precondition lineage sufficient for safe `SetNotificationRoute`; exact HTTP carrier remains later wire work.

Recipient human display names are not duplicated into Personal Notifications configuration by convenience. The Settings client composes recipient IDs with the already-admitted Organization member directory (`ListOrganizationMembers`) when human labels are needed.

Historical route revisions are not a public baseline collection; they remain D2/D7 recoverability state.

Rejected:

```text
custom NotificationKind catalog
route history browser
role/group recipients
Permission-derived recipients
rule predicates
quiet hours
channel preferences
per-user subscription settings
```

---

## 8. `SetNotificationRoute` — complete safety tuple

The operation mutates exactly one natural route slot:

```text
(Organization, NotificationKind)
```

Only the ten accepted ORG_ROUTED kinds are valid targets.

Desired state is exactly one of:

```text
CONFIGURED
  recipient_principal_ids = one-or-more exact human Principals

UNCONFIGURED
```

`UNCONFIGURED` is the accepted D2-R3 state; it is not `CONFIGURED([])`.

For `CONFIGURED`, Personal Notifications validates each recipient as an exact human Principal currently belonging/eligible in the same Organization and captures the accepted D2-R2 eligibility-continuity epoch. Permission/role never selects the recipients automatically.

Safety tuple:

```text
consequence class:
  non-consequential owner-local Product configuration mutation

idempotency:
  structural desired-state semantics
  no client idempotency key

concurrency / precondition:
  REQUIRED
```

A stale settings page cannot silently overwrite a newer routing decision. A successful repeated write of the same desired state under the valid current precondition is semantically idempotent.

Temporal routing remains D2-R2/D2-R3 authority: the new route revision affects future source occurrences by source-owner commit cutover and never retroactively retargets/backfills older occurrences.

---

## 9. Product surface consequence — candidate, not yet wire authority

If and only if this table is operator-ratified:

```text
Product operations:      99 → 103
ordinary Permissions:    30 → 31
new Permission:          notifications.manage
Principal kinds:         H/A/S unchanged
```

The canonical Product OAD remains **99/30** until the next separately executed D5 wire-authoring/proof step.

No count is preserved or expanded by preference.

---

## 10. Global negative controls

This D5-R3 table fails review if it admits any of the following by symmetry/convenience:

```text
CreateNotification Product operation
GetNotification detail operation without a distinct consumer
DeleteNotification
notifications.read
cross-Principal/admin Inbox
A/S human Inbox
GetUnreadCount / total_count
mark-all-read / bulk archive
public superseded-history/debug API
custom NotificationKind/template CRUD
email / push / digest settings
per-user subscription or mute engine
role/Permission/all-admin routing
routing/filter DSL
generic event-stream/SSE Product operation
Notification action that mutates/acknowledges source truth
generic payload/template fields
source-per-row N+1 reread as required Inbox rendering baseline
```

SSE/realtime may later be a disposable D6/D7 wake-up seam; it is not a Product operation admitted here.

---

## 11. Gate

```text
D0-R / D1-R                         ACCEPTED
D2-R / R2 / R3 / R4                ACCEPTED / OPERATOR-RATIFIED
D3-R                                ACCEPTED / OPERATOR-RATIFIED
D3-R1 presentation feed-forward    PASS
D5-R3 four-operation table         READY FOR OPERATOR REVIEW
canonical Product OAD              UNCHANGED — 99/30
D6 / D7 / D8                       BLOCKED for NOTIF-01
Product implementation             BLOCKED UNTIL D9
```

**Exact next action:** operator adjudicates only this frozen D5-R3 operation-admission table. If approved, author/prove the canonical Product OpenAPI amendment for exactly these four operations + `notifications.manage`; do not add routes/schemas beyond the ratified table and accepted D2/D3 feed-forward semantics.