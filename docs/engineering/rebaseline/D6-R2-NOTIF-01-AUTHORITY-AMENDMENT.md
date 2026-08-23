# NOTIF-01 — Personal Notifications Authority Amendment

> **Status:** OPEN — D0 ACCEPTED / D1 OWNER-BOUNDARY CANDIDATE / OPERATOR ADJUDICATION REQUIRED
> **Trigger:** operator-approved [Personal Notification Architecture Design](D6-R2-NOTIF-01-NOTIFICATION-ARCHITECTURE-DESIGN.md)
> **Execution plan:** [NOTIF-01 Authority Amendment Plan](D6-R2-NOTIF-01-AUTHORITY-AMENDMENT-PLAN.md)
> **Current accepted authority:** D0 NOTIF-01 amendment is accepted; D1+ remain unchanged except where a later NOTIF-01 section is explicitly operator-ratified
> **Current Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. Amendment law

NOTIF-01 is a bounded frontend-discovered Product gap. It does not rewrite accepted stage history and it does not grant later-stage decisions early.

```text
D0 scope — ACCEPTED
→ D1 owner/boundaries — CANDIDATE
→ D2 identity/data ownership — BLOCKED
→ D3 communication/events — BLOCKED
→ D5 Product wire — BLOCKED
→ D6 frontend — BLOCKED
→ D7 runtime realization — BLOCKED
→ D8 composed proof — BLOCKED
```

---

# D0 bounded amendment — Personal Notification Inbox

## 2. Product 1.0 scope addition — ACCEPTED

Product 1.0 adds **Personal Notification Inbox** as a bounded supporting capability of the Marketplace Operations Control Plane.

Purpose:

> allow an exact human Principal to discover personally relevant committed MPC facts, retain personal awareness state, and navigate back toward the current source authority without turning Notification into the source fact, workflow, authorization, audit record or acknowledgement of that fact.

This capability is additive to Marketplace Operations + Commercial Intelligence. It does not alter the accepted marketplace operating lifecycle, source-system authority model, Work lifecycle, Governance model, fulfillment model or commercial-intelligence authority.

## 3. Product-level authority — ACCEPTED

MPC **OWNS** Personal Notification awareness state for the recipient, bounded to:

```text
personal Notification existence
recipient
unread / read state
archived / active-Inbox state
source correlation sufficient to explain/navigate the awareness item
```

Notification is distinct from originating business truth, Operational Work obligation/resolution, Governance authorization/decision, Audit/history authority, source-object access, and acknowledgement that the originating fact was read, accepted or resolved.

Reading or archiving a Notification changes only Notification awareness state. It never resolves Work and never mutates or acknowledges the source condition by authority.

## 4. Launch-V1 human need and first required trigger — ACCEPTED

```text
Operational Work becomes explicitly assigned or reassigned
→ to an exact human Principal
→ that human must be able to discover the resulting personal Notification
```

The originating Work remains the actionable obligation and source of its own responsibility/assignment/escalation/resolution semantics.

D0 does **not** authorize generic notifications for every sale, listing, shipment, provider event, CRUD change, Work role/group, approval or asynchronous outcome. Further trigger classes require later bounded evidence and owning-authority review.

## 5. Recipient, Organization and access posture — ACCEPTED

Personal Notification Inbox is Organization-scoped, targeted to an exact human Principal, not a cross-Organization/global feed, not a human-style Inbox for A/S Principals, and never a grant of source access.

When a human follows a Notification to its source, current source authorization remains authoritative. Losing source access does not retroactively erase the recipient's Notification awareness history.

Exact identity representation, persistence ownership, RLS/isolation mechanics and source-reference shape belong to D2/D7.

## 6. Minimal Product outcome / no count requirement — ACCEPTED

Launch V1 requires durable distinction between unread/read and active/archived Notification state.

D0 does **not** require an exact unread aggregate count. A later frontend may surface only the truthful existence of one-or-more unread Notifications. A numeric badge/count requires separate authoritative aggregate evidence if later proven necessary.

D0 does not select route shape, topbar structure, API operations, Permission names, event mechanics, River jobs, SSE, `LISTEN/NOTIFY`, tables or schemas.

## 7. Explicit Launch-V1 non-goals — ACCEPTED

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

## 8. D0 negative controls — ACCEPTED

The target must fail review if it attempts to:

1. reuse Operational Work itself as the personal Inbox;
2. make Notification read/archive resolve or mutate Work/source truth;
3. make Notification an Audit or acknowledgement authority;
4. infer source access from possession of a Notification;
5. create a cross-Organization Inbox;
6. expose a human-style Inbox to A/S Principals by symmetry;
7. make e-mail/push/preferences/subscriptions a Launch-V1 gate;
8. create a generic Notification/subscription/event platform before a proved consumer requires it;
9. preserve the 99-operation/30-Permission census by forcing Notification into unrelated capabilities.

---

# D1 bounded amendment — Personal Notifications semantic owner

## 9. Supporting semantic owner — CANDIDATE

Add **Personal Notifications** as a small supporting semantic owner. This is a semantic boundary inside the accepted modular-monolith direction; it does **not** imply a service, database, process or deployment boundary.

Personal Notifications owns only:

```text
Notification personal-awareness lifecycle
exact recipient Principal reference
unread / read state
active-Inbox / archived state
Notification-local source correlation needed for explanation/navigation/deduplication
```

It explicitly does not own:

```text
originating business meaning or source resolution
Work responsibility / assignment / escalation / closure
Governance authorization or decision semantics
source-object access
Audit/history authority
delivery channels such as e-mail/push
generic subscriptions/preferences/templates
cross-owner workflow progression
```

The boundary exists because personal awareness has an independent lifecycle (`unread/read`, `active/archived`) and cannot truthfully be absorbed by Operational Work without making Work responsible for non-actionable awareness.

## 10. Minimal semantic edge — CANDIDATE

The only admitted baseline business edge is:

```text
Operational Work → Personal Notifications
```

Meaning:

> when Operational Work commits an assignment or reassignment occurrence to an exact human Principal, that committed Work-owned fact may make the Work personally relevant to that recipient; Personal Notifications owns whether/how that occurrence is represented as Notification awareness state.

Ownership does not move:

- Work remains authority for the assignment/reassignment and the actionable obligation;
- Personal Notifications never changes Work state;
- Notification read/archive never closes or acknowledges Work;
- the producer never writes Personal Notifications private state.

D1 does not choose whether the edge is Q/C/E/P, an event name, job payload, transaction mechanism or HTTP call. Those belong to D3/D7/D5.

## 11. No generic fan-out owner — CANDIDATE

D1 does **not** admit:

```text
AnyDomain → Personal Notifications
Notifications → all domains
event-per-CRUD notification fan-out
role/group broadcast semantics
generic subscription-driven routing
```

Future trigger owners (for example Governance or asynchronous action outcomes) require a bounded D1 reopen proving all of:

1. a committed producer-owned fact;
2. an exact human recipient without a generic subscription engine;
3. independent personal-awareness value rather than duplicate noise;
4. a stable occurrence discriminator;
5. no distortion of the originating owner's responsibility.

## 12. D1 forbidden boundary violations — CANDIDATE

The target must fail review if it makes Personal Notifications:

1. a workflow/event hub;
2. a generic task/ticket owner;
3. a source-business-truth or source-resolution authority;
4. a cross-owner mutable entity shared with producers;
5. a platform-wide polymorphic entity/reference graph;
6. an access-grant mechanism;
7. a delivery-channel platform by symmetry;
8. a reason to merge Work and Notification lifecycles.

## 13. D1 coherence result — CANDIDATE

The D1 addition is the smallest boundary coherent with accepted authority:

- Notification has independent recipient-facing state and lifecycle;
- Operational Work cannot absorb that lifecycle without violating its actionable-work charter;
- only one proved producer edge is admitted now;
- mechanism remains separate from authority;
- the owner remains supporting and narrow rather than becoming a generic communications platform.

No D2 identity/schema, D3 communication form, D5 operation/Permission, D6 UI realization or D7 mechanism becomes accepted through this D1 candidate.

## 14. Gate

```text
D0 NOTIF-01 amendment     ACCEPTED
D1 owner/boundary         READY FOR OPERATOR REVIEW
D2 amendment              BLOCKED
D3 amendment              BLOCKED
D5 Product/OAD change     BLOCKED
D6 bell/Inbox authority   BLOCKED
D7 runtime amendment      BLOCKED
D8 proof amendment        BLOCKED
Product implementation    BLOCKED UNTIL D9
```

**Exact next action:** operator adjudicates this D1 supporting semantic owner and the single `Operational Work → Personal Notifications` edge. If approved, mark D1 accepted and open only the bounded D2 identity/data-ownership gate.