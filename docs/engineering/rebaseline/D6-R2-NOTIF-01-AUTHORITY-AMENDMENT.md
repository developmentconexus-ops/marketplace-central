# NOTIF-01 — Personal Notifications Authority Amendment

> **Status:** OPEN — D0 PRODUCT-SCOPE CANDIDATE / OPERATOR ADJUDICATION REQUIRED
> **Trigger:** operator-approved [Personal Notification Architecture Design](D6-R2-NOTIF-01-NOTIFICATION-ARCHITECTURE-DESIGN.md)
> **Execution plan:** [NOTIF-01 Authority Amendment Plan](D6-R2-NOTIF-01-AUTHORITY-AMENDMENT-PLAN.md)
> **Current accepted authority:** D0–D8 + D5-R2/D8-R2 remain authoritative except where a later NOTIF-01 section is explicitly operator-ratified
> **Current Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. Amendment law

NOTIF-01 is a bounded frontend-discovered Product gap. It does not rewrite accepted stage history and it does not grant later-stage decisions early.

Authority progresses only in this order:

```text
D0 scope
→ D1 owner/boundaries
→ D2 identity/data ownership
→ D3 communication/events
→ D5 Product wire
→ D6 frontend
→ D7 runtime realization
→ D8 composed proof
```

Only the **D0 section below is candidate in this gate**. D1+ remain blocked until explicit operator approval of D0.

---

# D0 bounded amendment — Personal Notification Inbox

## 2. Product 1.0 scope addition

**Candidate for operator ratification.**

Product 1.0 adds **Personal Notification Inbox** as a bounded supporting capability of the Marketplace Operations Control Plane.

Purpose:

> allow an exact human Principal to discover personally relevant committed MPC facts, retain personal awareness state, and navigate back toward the current source authority without turning Notification into the source fact, workflow, authorization, audit record or acknowledgement of that fact.

This capability is additive to Marketplace Operations + Commercial Intelligence. It does not alter the accepted marketplace operating lifecycle, source-system authority model, Work lifecycle, Governance model, fulfillment model or commercial-intelligence authority.

## 3. Product-level authority

MPC **OWNS** Personal Notification awareness state for the recipient, bounded to:

```text
personal Notification existence
recipient
unread / read state
archived / active-Inbox state
source correlation sufficient to explain/navigate the awareness item
```

At D0 level, Notification is explicitly distinct from:

```text
originating business truth
Operational Work obligation/resolution
Governance authorization/decision
Audit/history authority
source-object access
acknowledgement that the originating fact was read, accepted or resolved
```

Reading or archiving a Notification changes only Notification awareness state. It never resolves Work and never mutates or acknowledges the source condition by authority.

## 4. Launch-V1 human need and first required trigger

The first proved Launch-V1 trigger is intentionally narrow:

```text
Operational Work becomes explicitly assigned or reassigned
→ to an exact human Principal
→ that human must be able to discover the resulting personal Notification
```

The originating Work remains the actionable obligation and source of its own responsibility/assignment/escalation/resolution semantics.

D0 does **not** authorize generic notifications for every sale, listing, shipment, provider event, CRUD change, Work role/group, approval or asynchronous outcome. Further trigger classes require later bounded evidence and owning-authority review.

## 5. Recipient, Organization and access posture

Personal Notification Inbox is:

- scoped to the current `Organization`;
- targeted to an exact **human Principal**;
- not a cross-Organization/global personal feed;
- not a human-style Inbox for automation/system Principals;
- never a grant of access to the originating source object.

When a human follows a Notification to its source, current source authorization remains authoritative. Losing source access does not retroactively erase the recipient's Notification awareness history.

Exact identity representation, persistence ownership, RLS/isolation mechanics and source-reference shape belong to D2/D7 and are **not selected by this D0 amendment**.

## 6. Minimal Product outcome / no count requirement

Launch V1 requires durable distinction between unread/read and active/archived Notification state.

D0 does **not** require an exact unread aggregate count. A later frontend may surface only the truthful existence of one-or-more unread Notifications. A numeric badge/count requires separate authoritative aggregate evidence if later proven necessary.

D0 does not select route shape, topbar structure, API operations, Permission names, event mechanics, River jobs, SSE, `LISTEN/NOTIFY`, tables or schemas. Those remain later-stage authority.

## 7. Explicit Launch-V1 non-goals

NOTIF-01 does **not** make any of the following Product 1.0 requirements:

```text
seen state
numeric unread aggregate
mark-all-read
bulk archive
notification preferences
per-kind user subscriptions
digests
e-mail notifications
mobile/web push
generic notification template platform
generic EventStore/pub-sub platform
external broker
cross-Organization Inbox
A/S human-style Inbox
generic entity/reference graph
Notification-triggered source mutation
```

These remain deferred unless later material evidence explicitly reopens the relevant authority.

## 8. D0 negative controls

The target must fail review if it attempts to:

1. reuse **Operational Work** itself as the personal Inbox;
2. make Notification read/archive resolve or mutate Work/source truth;
3. make Notification an Audit or acknowledgement authority;
4. infer source access from possession of a Notification;
5. create a cross-Organization Inbox;
6. expose a human-style Inbox to A/S Principals by symmetry;
7. make e-mail/push/preferences/subscriptions a Launch-V1 gate;
8. create a generic Notification/subscription/event platform before a proved consumer requires it;
9. preserve the current 99-operation/30-Permission census by forcing Notification into unrelated existing Product capabilities.

## 9. D0 coherence result

The amendment is coherent with accepted D0 because:

- MPC already requires human review and explicit Work for uncertainty/failure/exception conditions;
- accepted D0 explicitly distinguishes escalation from notification, so Work cannot truthfully absorb personal awareness semantics;
- Notification adds awareness only and does not alter the existing operating lifecycle or external-system authority;
- YAGNI keeps the addition narrower than a messaging/subscription platform.

No D1 owner, D2 identity/schema, D3 event, D5 operation/Permission, D6 UI realization or D7 mechanism becomes accepted merely because this D0 candidate exists.

## 10. Gate

```text
D0 NOTIF-01 candidate     READY FOR OPERATOR REVIEW
D1 amendment              BLOCKED
D2 amendment              BLOCKED
D3 amendment              BLOCKED
D5 Product/OAD change     BLOCKED
D6 bell/Inbox authority   BLOCKED
D7 runtime amendment      BLOCKED
D8 proof amendment        BLOCKED
Product implementation    BLOCKED UNTIL D9
```

**Exact next action:** operator adjudicates this D0 Product-scope amendment only. If approved, mark the D0 section accepted and then open the bounded D1 Personal Notifications owner/boundary gate. Do not open D1 before that decision.
