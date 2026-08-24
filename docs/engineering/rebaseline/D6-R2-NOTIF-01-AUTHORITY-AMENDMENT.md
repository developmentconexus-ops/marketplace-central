# NOTIF-01 — Personal Notifications Authority Amendment

> **Status:** OPEN — D0 ACCEPTED / D1 ACCEPTED / D2 IDENTITY-DATA-OWNERSHIP CANDIDATE / OPERATOR ADJUDICATION REQUIRED
> **Trigger:** operator-approved [Personal Notification Architecture Design](D6-R2-NOTIF-01-NOTIFICATION-ARCHITECTURE-DESIGN.md)
> **Execution plan:** [NOTIF-01 Authority Amendment Plan](D6-R2-NOTIF-01-AUTHORITY-AMENDMENT-PLAN.md)
> **Current Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. Amendment law

NOTIF-01 is a bounded frontend-discovered Product gap. It does not rewrite accepted stage history or grant later-stage decisions early.

```text
D0 scope — ACCEPTED
→ D1 owner/boundaries — ACCEPTED
→ D2 identity/data ownership — CANDIDATE
→ D3 communication/events — BLOCKED
→ D5 Product wire — BLOCKED
→ D6 frontend — BLOCKED
→ D7 runtime realization — BLOCKED
→ D8 composed proof — BLOCKED
```

---

# D0 bounded amendment — Personal Notification Inbox — ACCEPTED

Product 1.0 adds **Personal Notification Inbox** as a bounded supporting capability of the Marketplace Operations Control Plane.

Purpose:

> allow an exact human Principal to discover personally relevant committed MPC facts, retain personal awareness state, and navigate back toward current source authority without Notification becoming source fact, workflow, authorization, audit record or acknowledgement.

MPC owns Notification awareness state only: Notification existence, recipient, unread/read, active/archived, and bounded source correlation. Reading or archiving changes Notification only; it never resolves Work or mutates/acknowledges source truth.

Initial Launch-V1 trigger:

```text
Operational Work becomes explicitly assigned or reassigned
→ exact human Principal
→ that human must be able to discover the resulting personal Notification
```

The Inbox is Organization-scoped, human-Principal-targeted, not cross-Organization, not a human-style A/S Inbox, and never a source-access grant. D0 requires no exact unread aggregate count.

Accepted D0 non-goals remain: `seen`, numeric unread aggregate, mark-all-read, bulk archive, preferences, subscriptions, digests, e-mail/push, generic templates, generic EventStore/pub-sub, external broker, cross-Organization Inbox, A/S human Inbox, generic entity graph, and Notification-triggered source mutation.

---

# D1 bounded amendment — Personal Notifications semantic owner — ACCEPTED

**Personal Notifications** is a small supporting semantic owner inside the modular-monolith direction. It owns only the personal-awareness lifecycle, exact recipient reference, unread/read, active/archived, and Notification-local source correlation required for explanation/navigation/deduplication.

It does not own originating business meaning, Work responsibility/assignment/escalation/closure, Governance authorization, source access, Audit/history authority, delivery channels, generic subscriptions/preferences/templates, or cross-owner workflow progression.

The only admitted baseline semantic edge is:

```text
Operational Work → Personal Notifications
```

Work remains authority for assignment/reassignment and the actionable obligation. Personal Notifications never changes Work state, and producers never write Personal Notifications private state.

D1 admits no `AnyDomain → Personal Notifications`, Notifications→all-domains, event-per-CRUD fan-out, role/group broadcast, generic subscription routing, workflow/event hub or generic task domain. Future producer edges require a bounded D1 reopen proving exact recipient, independent awareness value, stable occurrence identity and no source-owner distortion.

D1 deliberately does not choose Q/C/E/P, event names, identity/schema, API, Permission, River or realtime.

---

# D2 bounded amendment — Notification identity and data ownership

## 2. Canonical Notification identity — CANDIDATE

`Notification` is MPC-owned canonical durable state under **Personal Notifications** and receives one stable opaque **NotificationID**.

Binding identity laws:

- NotificationID is non-business-semantic and non-reusable;
- it does not encode Organization, Principal, Work, kind, timestamp or lifecycle state;
- Notification identity does not change when read/archive state changes;
- Notification identity does not change when the source Work later changes;
- a reassignment never retargets or rewrites an earlier Notification: it creates a distinct source occurrence and, when D3 later admits delivery, a distinct Notification for the new recipient.

Exact UUID/ULID/database encoding remains later realization.

## 3. Organization ownership and recipient identity — CANDIDATE

Every Notification belongs to exactly one canonical `Organization` isolation scope and references exactly one canonical MPC `Principal` as recipient.

For the Product 1.0 Inbox behavior admitted by D0:

```text
recipient Principal kind = human
```

Recipient identity is the canonical Principal identity, never e-mail, username, display name, OIDC subject or role name by convenience.

The recipient attached to one Notification is historical and stable. Later Work reassignment, Membership/RoleAssignment change, Principal disablement or source-access revocation does not rewrite the Notification into another recipient or another Organization.

Current Product access remains separately authoritative: retained Notification history does not allow a disabled/non-member Principal to bypass current access gates.

## 4. Minimal canonical Notification state — CANDIDATE

The smallest durable state is:

```text
notification_id
organization_id
recipient_principal_id
kind
source_work_ref
source_occurrence_discriminator
created_at
read_at?
archived_at?
revision
```

Semantics:

```text
read_at = null       → unread
read_at != null      → read
archived_at = null   → active Inbox
archived_at != null  → archived
```

`created_at` is the Notification occurrence creation time owned by Personal Notifications; it is not silently substituted for the Work occurrence time.

`revision` is owner-local concurrency lineage for later safe Product mutation. D2 requires stable stale-write distinguishability but does not choose ETag/header/wire/database representation; D5/D7 own those mechanics.

No canonical `seen`, delivered, dismissed, acknowledged, resolved, severity, priority or generic status field is admitted.

## 5. Typed Work source reference — CANDIDATE

Because the only accepted D1 producer is Operational Work, the baseline source correlation is an explicit **typed Work reference**, not a generic entity reference or polymorphic platform graph.

```text
source_work_ref
```

references the accepted canonical Work identity in the same Organization scope. It does not transfer Work authority to Personal Notifications and does not grant the recipient source access.

Do **not** introduce a `source_ref` union yet. A broader closed union becomes admissible only after a second producer owner is separately accepted by D1.

Personal Notifications may retain only the source correlation needed to explain/deduplicate/navigate the awareness item. Current Work business state must not be copied into Notification as a second current authority merely to avoid a source read.

## 6. Source occurrence discriminator and duplicate identity — CANDIDATE

A Notification generated from Work assignment/reassignment must correlate to one particular committed Work-owned assignment occurrence, not merely to current `(work_id, recipient)` state.

Therefore `source_occurrence_discriminator` must be stable enough that:

```text
same Work occurrence delivered twice
→ same semantic Notification

distinct reassignment occurrence
→ distinct semantic Notification

a Work assigned to A, later B, later A again
→ the two assignments to A remain distinct occurrences
```

D2 does **not** invent a universal EventID. D3 must bind this discriminator to the smallest durable Work-owned occurrence semantics that satisfy recoverability and deduplication.

## 7. Write authority and history — CANDIDATE

Personal Notifications is the only write authority for Notification lifecycle state.

- Operational Work owns the source assignment/reassignment fact.
- Identity/access owns canonical Principal and Organization access state.
- Personal Notifications stores references; it does not mutate either authority.
- Notification read/archive history is not evidence that Work/source truth was read, accepted, acknowledged or resolved.
- Historical Notification meaning is not rewritten merely because current source state/access changed.

Deletion/retention policy is not selected by this amendment; target behavior must not silently erase material Notification history merely as a side effect of Work reassignment or access revocation.

## 8. D2 isolation and reference laws — CANDIDATE

- cross-Organization Notification references are invalid;
- `source_work_ref` must stay within the Notification Organization scope;
- Notification ownership must not be inferred from Work, recipient e-mail, Installation or process-global context;
- recipient scoping is a Product-access requirement, but exact PostgreSQL/RLS enforcement remains D7;
- typed references follow accepted D2 law; no universal `{entity_type, entity_id}` graph is introduced.

## 9. D2 negative controls — CANDIDATE

The target must fail review if it:

1. derives NotificationID from Work/Principal/time or recycles it;
2. uses e-mail/username/role as recipient identity instead of Principal;
3. changes an existing Notification recipient when Work is reassigned;
4. deduplicates only by `(work_id, recipient)` and therefore collapses two distinct assignment occurrences;
5. creates a generic source-entity graph before another producer is accepted;
6. copies current Work state into Notification as second authority;
7. treats Notification possession as source access;
8. deletes/rewrites history merely because Membership, access or current Work assignment changed;
9. adds `seen`, delivered, acknowledged, severity or generic workflow status by symmetry;
10. chooses table/RLS/API/event/River/SSE mechanics inside D2.

## 10. D2 coherence result — CANDIDATE

This model follows accepted D2 laws: canonical identity follows the new D1 owner; Organization remains the isolation root; Principal remains the canonical actor identity; typed references avoid a universal entity graph; historical meaning is stable; and one meaning has one write authority.

No D3 communication form, D5 wire/Permission, D6 UI realization or D7 persistence/runtime mechanism becomes accepted through this D2 candidate.

## 11. Gate

```text
D0 NOTIF-01 amendment     ACCEPTED
D1 owner/boundary         ACCEPTED
D2 identity/data ownership READY FOR OPERATOR REVIEW
D3 amendment              BLOCKED
D5 Product/OAD change     BLOCKED
D6 bell/Inbox authority   BLOCKED
D7 runtime amendment      BLOCKED
D8 proof amendment        BLOCKED
Product implementation    BLOCKED UNTIL D9
```

**Exact next action:** operator adjudicates this D2 Notification identity/data-ownership candidate only. If approved, mark D2 accepted and open only the bounded D3 trigger/communication gate.