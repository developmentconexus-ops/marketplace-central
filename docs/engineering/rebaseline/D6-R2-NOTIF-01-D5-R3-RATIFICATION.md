# NOTIF-01 D5-R3 — Final Product Operation Admission Ratification

> **Status:** ACCEPTED / OPERATOR-RATIFIED 2026-08-23
> **Accepted artifact:** [Final D5-R3 Product Operation Admission Table](D6-R2-NOTIF-01-D5-R3-OPERATION-ADMISSION-TABLE.md), blob `983beec633dbe250c0fcfd4d1e502cd67f71f06c`
> **Accepted parents:** D0-R + D1-R + D2-R/R2/R3/R4/R5 + D3-R/R1/R2
> **Global-Maximum inputs:** D5-F3 surviving findings + [D5-F4 Recipient Discovery](D6-R2-NOTIF-01-D5-F4-RECIPIENT-DISCOVERY-GLOBAL-MAXIMUM.md)
> **Current Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## Ratified result

The operator ratified the complete known Launch-V1 NOTIF-01 Product operation/access surface exactly as frozen by the accepted artifact:

```text
ListMyNotifications
UpdateMyNotificationAwarenessState
ListNotificationRoutes
ListNotificationRouteRecipientCandidates
SetNotificationRoute
```

Binding consequences:

1. the first two operations are H-only self-recipient capabilities with current Organization Membership/Product eligibility and no separate ordinary `notifications.read` Permission;
2. `ListNotificationRoutes`, `ListNotificationRouteRecipientCandidates` and `SetNotificationRoute` are H-only and require the one new ordinary Permission `notifications.manage`;
3. `ListNotificationRouteRecipientCandidates` is owned by IdentityAccess, is purpose-bounded to routing recipient discovery, and exposes only current eligible human `principal_id` + `display_name` with bounded collection pagination; it does not imply `access.read`, expose role/Permission/IdP state, duplicate the member directory into Personal Notifications or authorize a route write;
4. `ListMyNotifications` admits only the ratified archive/read state controls, one-or-more closed `NotificationKind` filtering and ordinary limit/cursor collection mechanics; no count/search/filter DSL or source-family Product identity is admitted;
5. the Inbox item carries the accepted D2-R4 presentation atom and D2-R5 kind-constrained immutable F02/F14 outcomes; `AUTHORIZATION_DECISION_RESULT` continues through `AuthorizationTargetRef`, not `AuthorizationDecisionRef`;
6. `UpdateMyNotificationAwarenessState` and `SetNotificationRoute` are desired-state owner-local updates with structural idempotency and required stale-write preconditions; neither uses a client idempotency key;
7. every configured route recipient is revalidated server-side at write time for exact Organization, human Principal, current Membership/Product eligibility and the ratified NotificationKind→source-read Permission floor before the D2-R2 eligibility epoch is captured;
8. routing Permission/access never selects recipients, recipient discovery never authorizes a write, and Personal Notifications never becomes an identity/member directory;
9. the global negative controls in the accepted table remain binding, including no public create/delete/count/bulk/admin Inbox, no generic payload/template/result/status fields, no routing DSL and no SSE Product operation.

If and only if the canonical OAD is separately authored and executably proved from this ratified contract, the derived Product wire consequence is:

```text
99 + 5 = 104 Product operations
30 + 1 = 31 ordinary Permissions
Principal kinds remain H/A/S
```

Counts are consequences, not design targets.

## Gate

```text
D2-R / R2 / R3 / R4 / R5           ACCEPTED / OPERATOR-RATIFIED
D3-R                                 ACCEPTED / OPERATOR-RATIFIED
D3-R1 / D3-R2                        PASS
D5-F4 recipient discovery            ACCEPTED / OPERATOR-RATIFIED
D5-R3 final operation admission      ACCEPTED / OPERATOR-RATIFIED
canonical Product OAD                UNCHANGED — 99/30
D6 / D7 / D8 for NOTIF-01            BLOCKED until wire authoring/proof
Product implementation               BLOCKED UNTIL D9
```

**Exact next action:** author and executably prove the canonical Product OpenAPI amendment for exactly the five ratified operations plus `notifications.manage`, deriving 104 Product operations / 31 ordinary Permissions. Do not add Product operations, Permissions or generic Notification machinery by symmetry or convenience.