# NOTIF-01 D5-R3 — Final Product Operation Admission Table

> **Status:** OPEN — FINAL GLOBAL-MAXIMUM CANDIDATE / OPERATOR ADJUDICATION REQUIRED
> **Accepted parents:** D0-R + D1-R + D2-R/R2/R3/R4 + [D2-R5 Ratification](D6-R2-NOTIF-01-D2-R5-RATIFICATION.md) + D3-R + D3-R1 + [D3-R2](D6-R2-NOTIF-01-D3-R2-TYPED-RESULT-CONTINUATION-FEED-FORWARD.md)
> **Global-Maximum inputs:** [D5-F3](D6-R2-NOTIF-01-D5-F3-GLOBAL-MAXIMUM-OPERATION-PERMISSION-REVIEW.md) surviving findings + [D5-F4 Recipient Discovery](D6-R2-NOTIF-01-D5-F4-RECIPIENT-DISCOVERY-GLOBAL-MAXIMUM.md) ACCEPTED / OPERATOR-RATIFIED
> **Parent D5 authority:** [D5 API](D5-API.md) + [D5-B2 Operation Admission Matrix](D5-B2-OPERATION-ADMISSION-MATRIX.md) + accepted W1/W2/W3/W4 laws
> **Scope:** freeze the complete known Launch-V1 Product operation / client-class / Permission / projection / safety surface before any canonical OpenAPI modification
> **Current Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Governing admission result

The complete known NOTIF-01 consumer set requires exactly five new Product operations and one new ordinary Permission:

```text
Personal Inbox
  ListMyNotifications
  UpdateMyNotificationAwarenessState

Organization routing settings
  ListNotificationRoutes
  ListNotificationRouteRecipientCandidates
  SetNotificationRoute
```

This result is not census-driven. Four operations were explicitly challenged and rejected because they leave the routing administrator without a least-privilege human recipient-discovery contract. More operations remain admissible if later evidence proves a distinct Product consumer.

If this table is operator-ratified exactly, the later Product wire consequence is:

```text
99 + 5 = 104 Product operations
30 + 1 = 31 ordinary Permissions
Principal kinds remain H/A/S
```

Counts are consequences only. **This document does not edit the canonical OAD.**

---

## 2. Final candidate operation matrix

| operationId | Semantic owner | Real Product consumer | Client | D3 class | Ordinary access | Admission |
| --- | --- | --- | --- | --- | --- | --- |
| `ListMyNotifications` | Personal Notifications | topbar bell preview, unread-presence probe, full personal Inbox | H only | Q | exact authenticated human recipient + current Organization Membership/Product eligibility; no separate ordinary Permission | **ADMIT** |
| `UpdateMyNotificationAwarenessState` | Personal Notifications | mark one own Notification read/unread and active/archived | H only | C/update | same exact self-recipient + current Organization eligibility rule | **ADMIT** |
| `ListNotificationRoutes` | Personal Notifications | `Configurações > Notificações` current Organization routing state | H only | Q | `notifications.manage` | **ADMIT** |
| `ListNotificationRouteRecipientCandidates` | IdentityAccess | least-privilege human selector for routing administration | H only | Q | `notifications.manage` | **ADMIT** |
| `SetNotificationRoute` | Personal Notifications | configure or intentionally unconfigure one Product-defined ORG_ROUTED kind | H only | C/update | `notifications.manage` | **ADMIT** |

All five operations are scoped to the exact Organization under accepted D5 organization-path laws. Exact path/method/schema spelling remains later wire work.

A/S do not receive a human Inbox, human recipient selector or notification-settings UI capability by symmetry. Internal source-owner → Personal Notifications reactions remain D3/D7 internal behavior, not public Product operations.

---

## 3. Self-Inbox access law

Baseline self Inbox does **not** introduce `notifications.read`.

The Product has already selected and materialized an exact human recipient. No proved Launch-V1 consumer requires a second ordinary Permission axis that selectively denies one eligible member access to awareness already addressed to that same member.

Both Inbox operations therefore require:

```text
authenticated Principal kind = H
+ exact path Organization
+ current Membership/Product-access eligibility in that Organization
+ Notification.recipient_principal_id == authenticated Principal
```

This capability is self-only. It never authorizes another Principal's Inbox.

Possession/readability of a Notification does not grant source access. Navigation through typed `source_ref` re-enters the source owner's current Product authorization.

Rejected baseline:

```text
notifications.read
admin/cross-Principal Inbox
Organization-wide shared Inbox
A/S human-style Inbox
```

---

## 4. `ListMyNotifications` — final candidate contract

### 4.1 Consumer-proven collection controls

Admit only:

```text
archive_state: active | archived
read_state: all | unread | read
notification_kind: one-or-more accepted NotificationKind values
limit
cursor
```

`notification_kind` is a closed bounded filter over accepted Product-defined kinds. A frontend grouping such as `Expedição` may compile to a closed set of kinds; no independent `source_family` identity or generic filter DSL is created.

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
- successful non-empty collection → unread presence;
- request failure/unavailability → knowledge unavailable, never empty.

### 4.3 Complete bounded Inbox item projection

The collection representation must expose the smallest Notification-owned state required by the accepted bell/full-Inbox consumer without source-per-row N+1 reads:

```text
notification_id
kind
subject_display_label
source_ref
source_occurred_at
created_at
read state
archive state
revision / accepted concurrency carrier

offering_async_result_outcome?    // required iff OFFERING_ASYNC_ACTION_RESULT
  converged | rejected | ambiguous | divergent

authorization_decision_outcome?   // required iff AUTHORIZATION_DECISION_RESULT
  authorize | reject
```

Laws:

- `subject_display_label` is the accepted immutable D2-R4 presentation snapshot; never source truth or navigation authority.
- The two result atoms are immutable, source-owned in meaning and kind-constrained under accepted D2-R5.
- No generic `result`, `status`, `reason`, `summary`, metadata, payload or template-variable field is admitted.
- `AUTHORIZATION_DECISION_RESULT.source_ref` is `AuthorizationTargetRef`, not `AuthorizationDecisionRef`; Governance retains decision-occurrence authority through the exact source occurrence key.
- `source_occurred_at` and `created_at` remain distinct so delayed materialization does not falsify chronology.

### 4.4 Rejected collection surface

```text
exact unread numeric aggregate
total_count
text search
source-family business identity
arbitrary source/kind filter DSL
separate bell-preview endpoint
GetNotification detail operation by symmetry
admin-wide/cross-Principal search
public superseded-history/debug collection
```

No `GetNotification` is admitted because the list item is the complete bounded Notification-owned Inbox representation and current source detail belongs to the source owner.

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

A stale client must not overwrite a later read/archive mutation or a later supersession revision. A superseded Notification cannot be resurrected by a stale awareness update.

One desired-state update operation is admitted instead of action-shaped `markRead`, `markUnread`, `archive`, `unarchive` endpoints because these are ordinary state axes of one owner resource.

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

`notifications.manage` is the one new ordinary Permission.

It means:

> manage the Organization's Personal Notifications routing configuration for the fixed Product-defined ORG_ROUTED NotificationKinds and discover the minimum human recipient identities required to perform that exact job.

It does **not** mean:

- receive those Notifications;
- read another user's Inbox;
- imply `access.read`;
- manage Membership/AccessRoles;
- expose role/Permission/IdP administration state;
- define custom NotificationKinds;
- create a generic subscription/rules system.

Permission remains access authority only. It never selects a recipient by itself.

---

## 7. `ListNotificationRoutes` — final candidate contract

The operation lists the fixed Product-defined baseline ORG_ROUTED slots and their **current** desired state for the Organization.

Baseline slot set is exactly:

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

Each current slot is semantically exactly one of:

```text
UNCONFIGURED
```

or:

```text
CONFIGURED
+ one-or-more exact human recipient_principal_ids
```

The representation includes owner-local current revision/precondition lineage sufficient for safe `SetNotificationRoute`; exact HTTP carrier remains later wire work.

Historical route revisions are not a public baseline collection; they remain D2/D7 recoverability state.

Rejected:

```text
custom NotificationKind catalog
route history browser
role/group recipients
Permission-derived recipient selection
rule predicates
quiet hours
channel preferences
per-user subscription settings
```

---

## 8. `ListNotificationRouteRecipientCandidates` — accepted Global-Maximum admission

This Q is owned by **IdentityAccess**, not Personal Notifications.

Its semantic meaning is:

> list current human Principals who are currently Product-access eligible and current members of the exact Organization, solely as human-readable references for Notification route recipient selection.

Minimum projection:

```text
principal_id
display_name
```

Collection mechanics:

```text
limit
cursor
```

The operation is:

```text
H only
Organization-scoped
notifications.manage
```

It does **not** expose by baseline:

```text
principal_kind        // all returned candidates are already H by contract
role_keys
Permissions
OIDC identity
email / username
Membership internals
eligibility_epoch
access administration state
```

It is not a generic Principal search API and does not replace `ListOrganizationMembers`.

The candidate list is **usability/discovery only**. Observing a candidate never authorizes a later route write.

---

## 9. `SetNotificationRoute` — complete safety tuple

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

### 9.1 Authoritative recipient validation

For every submitted recipient, the server revalidates at write time:

```text
same exact Organization
+ human Principal
+ current Organization Membership
+ current Product-access eligibility
+ current source-read eligibility required by the selected NotificationKind
```

Required source-read eligibility floor:

```text
MARKETPLACE_INSTALLATION_ATTENTION      → portfolio.read
AVAILABILITY_ATTENTION                  → availability.read
ECONOMIC_RECONCILIATION_ATTENTION       → economics.read
NEW_MARKETPLACE_SALE                    → sales.read
SALE_ATTENTION                          → sales.read
MATERIALIZATION_ATTENTION               → materialization.read
FULFILLMENT_ACTIONABLE                  → fulfillment.read
FULFILLMENT_ATTENTION                   → fulfillment.read
SHIPMENT_EXCEPTION                      → fulfillment.read
POST_SALE_ATTENTION                     → post_sale.read
```

These Permissions are eligibility checks only. They never select recipients and are not exposed by the candidate-discovery operation.

Only after current validation succeeds does Personal Notifications capture the accepted D2-R2 eligibility-continuity epoch server-side.

A candidate shown earlier and revoked before save must be rejected safely.

### 9.2 Mutation safety tuple

```text
consequence class:
  non-consequential owner-local Product configuration mutation

idempotency:
  structural desired-state semantics
  no client idempotency key

concurrency / precondition:
  REQUIRED
```

A stale settings page cannot silently overwrite a newer routing decision. Repeating the same desired state under the valid current precondition is semantically idempotent.

Temporal routing remains D2-R2/D2-R3 authority: a route revision affects future source occurrences by source-owner commit cutover and never retroactively retargets/backfills older occurrences.

---

## 10. Product surface consequence — candidate, not yet wire authority

If and only if this final table is operator-ratified:

```text
Product operations:      99 → 104
ordinary Permissions:    30 → 31
new Permission:          notifications.manage
Principal kinds:         H/A/S unchanged
```

The canonical Product OAD remains **99/30** until the separately executed wire-authoring/proof step.

No count is protected by preference.

---

## 11. Global negative controls

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
role/Permission/all-admin recipient selection
routing/filter DSL
generic event-stream/SSE Product operation
Notification action that mutates/acknowledges source truth
generic result/status/reason/payload/template fields
source-per-row N+1 reread as required Inbox rendering baseline
notifications.manage implying access.read
recipient candidate projection exposing role_keys/Permissions/IdP data
Personal Notifications owning a duplicate member directory
candidate discovery being treated as route-write authorization
AuthorizationDecisionRef retained in Notification source union without a consumer
governance.read granted to requester merely for F14 continuation
```

SSE/realtime may later be a disposable D6/D7 wake-up seam; it is not a Product operation admitted here.

---

## 12. Gate

```text
D0-R / D1-R                         ACCEPTED
D2-R / R2 / R3 / R4 / R5           ACCEPTED / OPERATOR-RATIFIED
D3-R                                ACCEPTED / OPERATOR-RATIFIED
D3-R1 / D3-R2                       PASS / NO TOPOLOGY REOPEN
D5-F3 surviving findings            INCORPORATED
D5-F4 recipient-discovery decision  ACCEPTED / OPERATOR-RATIFIED
D5-R3 final five-operation table    READY FOR OPERATOR REVIEW
canonical Product OAD               UNCHANGED — 99/30
D6 / D7 / D8                        BLOCKED for NOTIF-01
Product implementation              BLOCKED UNTIL D9
```

**Exact next action:** operator adjudicates only this final D5-R3 table. If approved, author and executably prove the canonical Product OpenAPI amendment for exactly these five operations + `notifications.manage`, yielding the derived 104/31 wire surface. Do not add operations, schemas or Permissions beyond this ratified contract by symmetry or convenience.