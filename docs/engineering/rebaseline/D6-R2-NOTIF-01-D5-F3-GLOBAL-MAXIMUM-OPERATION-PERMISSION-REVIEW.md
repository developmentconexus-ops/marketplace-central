# NOTIF-01 D5-F3 — Global Maximum Operation / Permission Review

> **Status:** D5-R3 GLOBAL-MAXIMUM REVIEW / TARGETED D2-R5 REQUIRED BEFORE FINAL TABLE RATIFICATION
> **Accepted inputs:** D0-R + D1-R + D2-R/R2/R3/R4 + D3-R + D3-R1
> **Reviewed candidate:** [D5-R3 Product Operation Admission Table](D6-R2-NOTIF-01-D5-R3-OPERATION-ADMISSION-TABLE.md)
> **Operator direction:** operation/Permission counts are consequences only; Global Maximum completeness outranks preserving 103/31 or any prior census
> **Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Review question

The operator explicitly rejected any interpretation that the current operation/Permission census is a constraint. The review therefore asks:

> Against every accepted NOTIF-01 human job, semantic family, access law and frontend consumer, what is the complete smallest Product contract — even if that requires more operations, Permissions or typed state than the current candidate?

Result:

- **four public Product operations remain sufficient**;
- **one new ordinary Permission remains sufficient**;
- the prior candidate omitted one bounded collection filter already proved by the accepted census;
- the prior D2 presentation state is insufficient to express two result-bearing Notification families without rereading source state;
- `AUTHORIZATION_DECISION_RESULT` has the wrong continuation reference for its exact requester consumer;
- no additional operation/Permission is justified merely to repair those state/continuation gaps.

Counts remain consequences, not design objectives.

---

## 2. Operation count review — four remain sufficient

All accepted Product consumers are covered by:

```text
ListMyNotifications
UpdateMyNotificationAwarenessState
ListNotificationRoutes
SetNotificationRoute
```

The review explicitly challenged and still rejects, for lack of a distinct consumer:

```text
GetNotification
GetUnreadCount
GetNotificationPreview
CreateNotification
DeleteNotification
MarkAllRead
BulkArchive
ListNotificationKinds
GetNotificationRoute
ListNotificationRouteHistory
cross-Principal/admin Inbox operations
public SSE/event-stream operation
```

`ListMyNotifications` is the bounded collection read; the list item is the complete Notification-owned Inbox representation. Source detail remains source-owner authority.

---

## 3. Collection completeness correction — bounded `NotificationKind` filter

The operator-approved Trigger/Audience Census already proved future Inbox **filters/grouping by Product-defined kind/source family** as a real frontend consumer. The frozen D5-R3 table admitted only archive/read filters and therefore under-covered accepted evidence.

`ListMyNotifications` must additionally admit a bounded closed kind filter:

```text
notification_kind = one or more accepted NotificationKind values
```

Exact repeated-query/array wire grammar remains later D5 wire spelling.

A frontend source-family preset may compile to the corresponding closed set of Product kinds. No independent `source_family` business identity is created because family grouping is a presentation mapping over the already-stable Product-defined kinds.

Examples:

```text
Expedição
→ FULFILLMENT_ACTIONABLE
 + FULFILLMENT_ATTENTION
 + SHIPMENT_EXCEPTION

Governança
→ AUTHORIZATION_ACTION_REQUIRED
 + AUTHORIZATION_DECISION_RESULT
```

This is **not** an arbitrary filter DSL. Text search, generic source predicates, custom expressions and aggregation/count facets remain unproved.

No new Product operation or Permission is required for this correction.

---

## 4. Permission review — one new Permission remains semantically correct

The review does **not** minimize Permissions by count. It challenges whether the accepted jobs need additional independent access capabilities.

### 4.1 Self Inbox

No `notifications.read` is justified.

Self Inbox is exact-recipient awareness already addressed to the authenticated human. A second Permission would create an unsupported product state where an eligible human is the canonical recipient but cannot read their own awareness.

Cross-Principal Inbox remains a different, unproved job and is not smuggled into self access.

### 4.2 Routing administration

`notifications.manage` remains a real distinct ordinary Permission because choosing Organization awareness recipients is neither `access.manage` nor any source-domain manage capability.

No separate `notifications.routes.read` is justified because no read-only routing-auditor consumer is currently accepted. `ListNotificationRoutes` exists for the same settings administration job as `SetNotificationRoute`.

No per-kind manage Permission is invented without evidence of delegated notification-administration boundaries. If later Product evidence requires Sales-routing administrators to be isolated from Fulfillment-routing administrators, D5 reopens and may split access then.

### 4.3 Source-access eligibility uses existing Permissions, not new notification Permissions

Notification routing never grants source access. For ORG_ROUTED kinds, configuration/materialization must fail closed for a recipient who cannot currently read the source continuation class.

The bounded eligibility floor uses existing Permissions:

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

These Permissions are **eligibility checks only**. They never select recipients. The Organization still explicitly configures exact human Principals.

`SetNotificationRoute(CONFIGURED)` must reject a recipient who lacks the current required source-read eligibility for that kind rather than create a latent unusable binding. D3 materialization/recovery revalidates current eligibility; a later Permission revocation temporarily blocks new awareness while the route decision remains. Permission regrant may resume the still-explicit route because D2-R2 deliberately treats ordinary Permission changes separately from Membership/access-continuity epochs.

DIRECT_SOURCE/OWNER_DERIVED kinds similarly preserve current access independently from recipient derivation; their exact continuation access is owner/target-specific rather than a new notification Permission.

---

## 5. D2 presentation falsifier A — result-bearing kinds cannot express the approved human meaning

Accepted semantic contracts require:

```text
OFFERING_ASYNC_ACTION_RESULT
→ human needs the materially committed result:
   converged | rejected | ambiguous | divergent

AUTHORIZATION_DECISION_RESULT
→ requester needs the committed disposition:
   authorize | reject
```

Current Notification-owned presentation state has only:

```text
kind
subject_display_label
```

That can identify the subject but cannot truthfully distinguish the accepted result variants.

Examples of the defect:

```text
kind = OFFERING_ASYNC_ACTION_RESULT
subject_display_label = "Alteração de preço · Anúncio X"
```

cannot tell the client whether copy should communicate convergence, rejection, ambiguity or divergence.

Likewise:

```text
kind = AUTHORIZATION_DECISION_RESULT
subject_display_label = "Alteração de preço · Anúncio X"
```

cannot distinguish “Sua solicitação foi aprovada” from “Sua solicitação foi rejeitada”, despite that distinction being explicit in the operator-approved family contract.

Forcing one source reread merely to know the Notification's own immutable result defeats the bounded Inbox presentation model and makes result copy unavailable when the source read is inaccessible/unavailable.

The smallest repair is **typed closed result state**, not `summary`, `reason`, `metadata` or a payload template engine.

---

## 6. D2 continuation falsifier B — Governance result points to the wrong consumer continuation

Accepted `AUTHORIZATION_DECISION_RESULT` audience is the exact requester/initiator, not the Governance administrator/approver.

Current D2-R compatibility uses:

```text
AUTHORIZATION_DECISION_RESULT
→ AuthorizationDecisionRef
```

But canonical D5 authority requires `governance.read` for `GetAuthorizationDecision`, while the requester may legitimately have only the action-owning target permissions. `governance.read` must not be granted broadly merely to make a Notification link work.

The governed target is already a closed accepted identity:

```text
AuthorizationTargetRef =
    ListingIntentRef
  | PriceIntentRef
  | BusinessOrderIntentRef
  | InvoicingIntentRef
```

For the requester-result Notification, the correct continuation is therefore the exact governed target. Governance decision identity remains the committed occurrence discriminator/history under Governance, but it need not be the Notification navigation/source subject.

Smallest repair:

```text
AUTHORIZATION_DECISION_RESULT.source_ref
→ AuthorizationTargetRef

source_occurrence_key
→ stable Governance-owned decision occurrence discriminator
  (the canonical AuthorizationDecision identity may supply this meaning)
```

This lets the requester continue in the action-owning workspace under current target authorization without acquiring Governance read authority.

No new Product operation or ordinary Permission is required.

---

## 7. Targeted D2-R5 repair required

Before the D5-R3 operation table can be finally ratified, D2 must reopen only to:

1. add one typed immutable Offering async-result outcome for `OFFERING_ASYNC_ACTION_RESULT`;
2. add one typed immutable Governance decision disposition for `AUTHORIZATION_DECISION_RESULT`;
3. change the latter kind's source continuation from `AuthorizationDecisionRef` to `AuthorizationTargetRef` while retaining Governance-owned occurrence discrimination;
4. remove `AuthorizationDecisionRef` from the Notification source union if no admitted kind still consumes it;
5. keep all generic payload/reason/template machinery rejected.

D3 then needs only bounded feed-forward revalidation of the two typed immutable result atoms and changed F14 target ref.

---

## 8. Global-Maximum conclusion

The complete known Launch-V1 NOTIF-01 Product surface is **not** being held to four operations or one Permission by preference.

After challenging every accepted consumer:

```text
public operations required      4
new ordinary Permissions        1
bounded kind filter             required
new typed result atoms          2
F14 continuation correction     required
```

More operations/Permissions remain fully admissible if a concrete downstream consumer or safety property later proves them. They are not added now because no accepted job requires them, not because of a census target.

## 9. Gate

```text
D2-R/R2/R3/R4                    ACCEPTED
D3-R/R1                          ACCEPTED / PASS
D5-F3 Global Maximum review      COMPLETE
D2-R5 typed result/continuation  REQUIRED / NEXT
D5-R3 final table ratification   BLOCKED BY D2-R5
canonical Product OAD            UNCHANGED — 99/30
D6 / D7 / D8                     BLOCKED for NOTIF-01
Product implementation           BLOCKED UNTIL D9
```

**Exact next action:** adjudicate only the bounded D2-R5 typed-result / Governance-requester-continuation repair. Do not preserve 103/31 by preference and do not edit the canonical Product OpenAPI first.