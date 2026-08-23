# NOTIF-01 D5-F1 — Operation Surface & Route Reversibility Finding

> **Status:** D5-R3 ANALYSIS FINDING / TARGETED D2 REOPEN REQUIRED BEFORE OPERATION-TABLE RATIFICATION
> **Accepted inputs:** D0-R + D1-R + D2-R + D2-R2 + D3-R, all operator-ratified where applicable
> **Parent D5 authority:** [D5 API](D5-API.md) + [D5-B2 Operation Admission Matrix](D5-B2-OPERATION-ADMISSION-MATRIX.md)
> **Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. D5 admission law applied to Personal Notifications

D5 requires the smallest external operation surface justified by real Product 1.0 consumers, with one semantic owner, explicit Organization scope, client class, access authority and complete mutation safety tuple before path/schema spelling.

Applying that law to the accepted bell/Inbox + Organization routing consumers yields **four** public Product operations. No OAD edit is authorized yet.

## 2. Candidate operation table

| Candidate operation | Consumer | Client | Class | Ordinary access | Admission direction |
| --- | --- | --- | --- | --- | --- |
| `ListMyNotifications` | topbar bell preview, unread-presence probe, full personal Inbox | human only | Q | exact authenticated human recipient + current Organization Membership/Product eligibility; no separate ordinary Permission | **ADMIT candidate** |
| `UpdateMyNotificationAwarenessState` | mark read/unread and archive/unarchive one personal Notification | human only | C/update | same self-only recipient rule | **ADMIT candidate** |
| `ListNotificationRoutes` | `Configurações > Notificações` current routing state | human only | Q | new `notifications.manage` | **ADMIT candidate** |
| `SetNotificationRoute` | configure or intentionally return one ORG_ROUTED kind to unconfigured state | human only | C/update | new `notifications.manage` | **ADMIT candidate, but exposes D2 gap** |

Candidate Product surface consequence if later ratified exactly as above:

```text
99 + 4 = 103 Product operations
30 + 1 = 31 ordinary Permissions
Principal kinds remain H/A/S
```

Counts are consequences, not goals.

## 3. Why personal Inbox read does not gain `notifications.read`

The Inbox is exact-recipient, self-only awareness state. A second ordinary Permission axis would allow this contradiction:

```text
accepted Work/routing semantics target human P
+ P remains an eligible Organization member
+ P lacks notifications.read
→ P cannot discover the personal awareness explicitly addressed to P
```

No Product 1.0 consumer requires selectively granting one member access to their own Inbox while denying another member access to theirs. Therefore baseline self-Inbox access is a bounded self-only capability tied to the authenticated human recipient plus current Organization membership/Product eligibility.

This does **not** grant source access. Navigation still re-enters current source authorization.

Cross-Principal Inbox read/search remains rejected.

## 4. `ListMyNotifications` consumer contract — candidate

One collection operation serves bell + full Inbox rather than creating separate preview/count endpoints.

Admit only consumer-proven controls:

```text
archive_state: active | archived
read_state: all | unread | read
limit
cursor
```

Baseline default is current, non-superseded, active awareness ordered by Notification `created_at` descending with a stable tie-breaker. `source_occurred_at` remains separately exposed so delayed materialization does not falsify business chronology.

The bell's unread dot does not require an aggregate/count operation:

```text
ListMyNotifications(
  archive_state=active,
  read_state=unread,
  limit=1
)
```

successful empty collection → no current unread personal item
successful non-empty collection → unread presence

Reject baseline:

- exact unread numeric aggregate;
- total count;
- mark-all-read;
- bulk archive;
- text search;
- arbitrary source/kind filter DSL;
- admin-wide/cross-Principal notification search;
- separate bell-preview endpoint;
- public superseded-history filter merely for debugging.

Superseded Notifications remain durable owner history/correctness state but are not current personal-awareness collection members.

## 5. `UpdateMyNotificationAwarenessState` safety tuple — candidate

The operation updates only the authenticated recipient's exact Notification and supports the two orthogonal desired-state axes:

```text
read | unread
active | archived
```

It does not acknowledge/resolve/dismiss the source.

Safety tuple:

```text
consequence class:
  non-consequential owner-local awareness mutation

idempotency:
  structural desired-state semantics; no client idempotency key

concurrency / precondition:
  required — stale mutation cannot silently overwrite a later read/archive/supersession revision
```

One owner-local update operation is preferred over four action-shaped endpoints because read/unread and archive/unarchive are ordinary mutable state of the same Notification resource and D2 explicitly makes those axes orthogonal.

No delete operation is admitted.

## 6. Routing operations and Permission — candidate

`notifications.manage` is a new ordinary Product capability because managing who receives Personal Notifications is semantically distinct from managing AccessRoles/Membership. It does not transfer routing ownership to the identity/access substrate and does not make Permission itself a recipient selector.

`ListNotificationRoutes` returns the bounded Product-defined ORG_ROUTED slots and current state only. There is no generic kind registry or subscription catalog.

Frontend recipient selection composes with the already-accepted `ListOrganizationMembers`; Personal Notifications does not duplicate the Organization member directory merely for convenience.

`SetNotificationRoute` is one desired-state operation per `(Organization, NotificationKind)` slot. It must support both legitimate admin outcomes:

```text
CONFIGURED
→ one-or-more exact human Principal recipients

UNCONFIGURED
→ no personal routing for future occurrences of this kind
```

Safety tuple:

```text
consequence class:
  non-consequential owner-local Product configuration mutation

idempotency:
  structural desired-state semantics; no client idempotency key

concurrency / precondition:
  required — stale route state cannot overwrite a newer routing decision
```

No custom rules, groups, roles, expressions, per-user subscription preferences or notification delivery channels are admitted.

## 7. Targeted D2 falsifier — routing must be reversible

Accepted D2-R/D2-R2 model supports:

```text
before any route revision → UNCONFIGURED
configured route revision → one-or-more recipient bindings
```

but deliberately did not admit an explicit later transition back to `UNCONFIGURED`.

The now-proven Product settings consumer makes irreversible configuration invalid:

```text
admin configures NEW_MARKETPLACE_SALE → [A]
admin later decides this Organization should have no personal bell routing for new sales
```

Forcing a permanent recipient after first configuration, inventing a dummy recipient or treating stale route history as deleted would violate Product behavior and D2-R2 temporal recovery semantics.

Therefore D5-R3 stops before operation-table ratification and reopens only D2 routing history to admit an explicit **unconfigured route revision/cutover**.

This is not `configured empty` and not a personal opt-out system. It is Organization-level absence of routing for one Product-defined kind.

## 8. Negative controls

D5-R3 must continue to reject:

- `CreateNotification` public operation;
- delete Notification;
- cross-user/admin Inbox;
- `notifications.read` without a separately proven selective-access consumer;
- A/S human-style Inbox;
- unread count endpoint;
- mark-all-read/bulk mutations;
- custom NotificationKind/template CRUD;
- e-mail/push/digest preferences;
- per-user subscription engine;
- route by role/Permission/all-admin fallback;
- generic routing/filter DSL;
- notification event-stream/SSE Product operation merely because D3 uses E internally;
- notification actions that mutate/acknowledge source business state.

## 9. Gate

```text
D3-R communication/propagation       ACCEPTED / OPERATOR-RATIFIED
D5-R3 four-operation direction       DERIVED / NOT YET RATIFIED
D5-F1 route reversibility falsifier  FOUND
D2-R3 route-unconfigure repair       REQUIRED / NEXT
canonical Product OAD                UNCHANGED
D6 / D7 / D8                         BLOCKED for NOTIF-01
Product implementation               BLOCKED UNTIL D9
```

**Exact next action:** adjudicate only the targeted D2-R3 route-unconfigure temporal-history repair. Then return to this four-operation D5-R3 table for operator ratification before any OpenAPI edit.