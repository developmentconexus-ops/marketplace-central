# NOTIF-01 D5-F2 — Inbox Presentation Context Finding

> **Status:** D5-R3 ANALYSIS FINDING / TARGETED D2 REOPEN REQUIRED BEFORE OPERATION-TABLE RATIFICATION
> **Accepted inputs:** D0-R + D1-R + D2-R + D2-R2 + D2-R3 + D3-R
> **Existing design evidence:** [Notification Architecture Design](D6-R2-NOTIF-01-NOTIFICATION-ARCHITECTURE-DESIGN.md) already requires a bounded recent-Inbox preview and permits a bounded title/summary for usability
> **Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Consumer that exposes the gap

The already-approved topbar bell and full personal Inbox need to render several Notifications in one collection read. A representation containing only:

```text
kind
source_ref
source_occurrence_key
timestamps
```

is semantically correct but insufficiently human-operable for multiple same-kind items.

Examples:

```text
WORK_ASSIGNMENT
WORK_ASSIGNMENT
WORK_ASSIGNMENT
```

cannot be meaningfully distinguished by a person if the only differentiator is an opaque WorkID. Similar problems exist for materialization, fulfillment, approvals and other source refs whose canonical identity is intentionally opaque.

## 2. Why per-item source rereads are not the baseline answer

Making the frontend or Product read path reread every referenced source solely to render each Inbox row creates avoidable coupling:

```text
ListMyNotifications
→ N source refs
→ N owner reads
```

That would:

- create frontend/server N+1 behavior for an ordinary Inbox list;
- make basic Inbox comprehensibility depend on simultaneous availability of several source owners;
- degrade historical Notification readability after source access loss/deletion even though the Notification remains valid personal history;
- encourage a screen-shaped aggregation surface or generic payload later to repair the UX.

Source reread remains mandatory when the user opens/navigates to the source because current authorization and current truth belong there. It is not required merely to distinguish rows in the personal Inbox.

## 3. Existing D2 stop trigger is now met

Accepted D2-R intentionally forbids free-form payload/template/source-state copies and says that if D5/D6 later proves an immutable presentation snapshot is required, D2 must reopen for explicit typed fields.

That consumer is now proved.

The required state is much smaller than a generic payload: one bounded, immutable, notification-safe **subject display label** sufficient to distinguish the awareness item from neighboring items. Final sentence/title copy remains derived from `NotificationKind` by the client presentation layer.

## 4. Required semantic shape

Candidate minimum:

```text
notification_subject_display_label
```

Properties required before D5 can ratify the operation table:

- non-authoritative human presentation only;
- immutable after Notification creation;
- source-owner-derived for the exact admitted occurrence;
- safe to retain even if later source access is revoked;
- never used for identity, routing, deduplication, authorization, navigation or source mutation;
- no buyer name/address/payment/fiscal detail, credentials/secrets, raw provider body or arbitrary source payload;
- no generic metadata/template-variable bag;
- final localized Notification title remains derived from `NotificationKind`, not stored as source truth.

Examples only:

```text
NEW_MARKETPLACE_SALE        → "Mercado Livre · Pedido #200391"
FULFILLMENT_ATTENTION       → "Venda #200388"
WORK_ASSIGNMENT             → "Revisar materialização da venda #200340"
AUTHORIZATION_ACTION_REQUIRED → "Alteração de preço · Anúncio MLB-..."
```

The examples are presentation guidance, not frozen copy.

## 5. Why one subject label is enough now

No second summary/detail field is proved necessary for baseline correctness. Deadline, status, buyer, amount and other current/sensitive detail remains source-owned and is read only where the later D5/D6 interaction genuinely requires it.

A future consumer that proves a second immutable presentation atom must reopen the smallest D2/D3 contract rather than expanding this field into JSON.

## 6. Gate

```text
D2-R3 route reversibility               ACCEPTED / OPERATOR-RATIFIED
D3-R communication                      ACCEPTED
D5-R3 four-operation direction          DERIVED / NOT YET RATIFIED
D5-F2 human-operable Inbox finding      PROVED
D2-R4 presentation snapshot             REQUIRED / NEXT
canonical Product OAD                   UNCHANGED
D6 / D7 / D8                            BLOCKED for NOTIF-01
Product implementation                  BLOCKED UNTIL D9
```

**Exact next action:** adjudicate only the bounded D2-R4 Notification presentation-snapshot repair. If accepted, revalidate its feed-forward effect on D3 and return to the four-operation D5-R3 table before any OpenAPI edit.