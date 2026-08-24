# NOTIF-01 D5-R4 — Canonical Product OAD Wire Proof

> **Status:** PROVED / CANONICAL WIRE AMENDMENT
> **Ratified parent:** [D5-R3 Final Product Operation Admission Ratification](D6-R2-NOTIF-01-D5-R3-RATIFICATION.md)
> **Wire authority:** `contracts/api/product/openapi.yaml` + `contracts/api/product/paths-notifications.yaml`
> **Exact proved candidate:** `2517727d438d6ddb2229782d86385505312a119c`
> **CI evidence:** `ci #526` SUCCESS + `pr-title #583` SUCCESS on the exact candidate
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Proven wire result

The canonical Product OAD now realizes exactly the five ratified NOTIF-01 operations:

```text
ListMyNotifications
UpdateMyNotificationAwarenessState
ListNotificationRoutes
ListNotificationRouteRecipientCandidates
SetNotificationRoute
```

and the one ratified ordinary Permission:

```text
notifications.manage
```

Current canonical wire consequence:

```text
Product operations:      104 / 104
ordinary Permissions:     31 / 31
Principal kinds:           H / A / S only
OpenAPI:                   3.1.2
stable origin:             https://conexus.fun
active runtime baseline:   NONE
```

Counts are consequences of accepted authority, not protected design targets.

## 2. TDD evidence

The NOTIF-01 proof was introduced before the wire. Exact RED candidate `da87ca6c007bf3f4760a482a635935cea2b61099` passed the existing historical Product/API, authentication, Performance and OP-READ controls and then failed specifically with:

```text
operation not found: ListMyNotifications
```

The wire then proceeded through falsifier-driven corrections rather than guard weakening.

The final exact candidate `2517727d438d6ddb2229782d86385505312a119c` proves:

```text
product_oad_baseline_non_regression                 PASS  95/29
product_oad_historical_99_non_regression            PASS  99/30
product_oad_current_generated_projection_semantics  PASS
product_oad_operations                              104/104
product_oad_auth_negative_controls                  5/5
performance_evidence_knowledge_proof                PASS
operational_read_contract_proof                     PASS
notification_oad_operations                         104/104
notification_oad_permissions                        31/31
notification_oad_negative_controls                  8/8
notification_oad                                    PASS
repository full gate                                PASS
legacy runtime population                           0
```

The 99-operation verifier remains preserved as a historical non-regression fixture. The current proof composes 104 over that accepted lineage rather than rewriting history to pretend the older surface was always 104.

## 3. NOTIF-01 wire properties proved

The focused proof falsifies at least:

1. all five operation IDs, exact path/method shapes and H-only client admission;
2. `PersonalNotifications` ownership for Inbox/routes and `IdentityAccess` ownership for recipient discovery;
3. self-Inbox exact-recipient + current-Organization-Membership law with no ordinary `notifications.read`;
4. `notifications.manage` for route reads, recipient discovery and route writes without implying `access.read`;
5. exact 31-Permission effective access vocabulary;
6. bounded archive/read/`NotificationKind`/limit/cursor Inbox controls;
7. recipient-candidate projection restricted to `principal_id + display_name`;
8. required stale-write preconditions and absence of client idempotency keys for owner-local desired-state updates;
9. exact fourteen `NotificationKind` branches and ten ORG_ROUTED kinds;
10. required F02/F14 typed immutable outcome atoms and rejection of generic result/status/reason/payload/template fields;
11. F14 `AuthorizationTargetRef` continuation with no retained `AuthorizationDecisionRef` navigation variant;
12. exact ORG_ROUTED `NotificationKind -> source-read Permission` eligibility mapping;
13. candidate discovery never becoming route-write authorization by contract.

Generated TypeScript and Go projections are deterministic and compile/test under the existing Product OAD proof profile.

## 4. OP-READ proof boundary correction

The accepted D5-R2 OP-READ verifier previously asserted the then-current global census `99/30` because D5-R2 itself was required not to expand Product surface. D5-R3 later legitimately adds NOTIF-01 operations.

The current OP-READ proof therefore owns only its actual invariant set:

```text
owner-local operational projections
accepted operational filters
no synthetic operational workflow/stage/priority/dashboard truth
```

Global Product census proof is delegated to the canonical Product OAD verifier, which now proves both the historical 99/30 fixture and current 104/31 surface. The OP-READ protection was not weakened; a superseded cross-stage census assumption was removed from the wrong owner.

## 5. Negative controls retained

The proved wire still rejects the ratified NOTIF-01 anti-requirements, including:

```text
notifications.read
cross-Principal/admin Inbox
unread count / total_count
mark-all-read / bulk archive
custom NotificationKind/template CRUD
email/push/digest preference platform
routing/filter DSL
public SSE/event-stream Product operation
generic result/status/reason/payload/template fields
notifications.manage implying access.read
recipient discovery exposing role/Permission/IdP state
Personal Notifications duplicating the member directory
candidate discovery authorizing route writes
F14 navigation through AuthorizationDecisionRef
```

No Product runtime, router, database schema, River job, SSE endpoint or frontend implementation is selected by this proof.

## 6. Gate

```text
D5-R3 final operation admission    ACCEPTED / OPERATOR-RATIFIED
D5-R4 canonical OAD wire           PROVED / CANONICAL
Product OAD                        104 operations / 31 Permissions / H-A-S
active runtime baseline            NONE
D6-R NOTIF frontend feed-forward   NEXT
D7-R / D8-R                        BLOCKED until prior NOTIF gate lands
Product implementation             BLOCKED UNTIL D9
```

**Exact next action:** execute only the bounded NOTIF-01 D6-R frontend feed-forward: bell + personal Inbox + notification routing settings plus the accepted bounded B00 topbar utility-slot reopen, deriving exact frontend interaction/surface authority from the proved 104/31 OAD. Do not begin D7 runtime realization, D8 composed proof or Product implementation first.