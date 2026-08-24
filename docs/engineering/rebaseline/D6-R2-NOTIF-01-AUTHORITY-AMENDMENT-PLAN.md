# NOTIF-01 — Personal Notifications Authority Amendment Plan

> **Execution:** Use `superpowers:subagent-driven-development` (recommended) or `superpowers:executing-plans` task-by-task. This plan changes architecture/contracts/proofs only; Product implementation remains blocked until accepted D9.

**Goal:** Carry the operator-approved Personal Notifications design coherently through D0→D8, add the bounded topbar bell + organization-scoped Inbox to frontend authority, and preserve the existing modular-monolith/PostgreSQL/River architecture.

**Architecture:** `PersonalNotifications` is a small supporting semantic owner for Organization-scoped awareness state targeted to an exact human Principal. Source owners commit only their own facts and atomically enqueue a River reaction; the Notification owner writes its own durable state. The browser reads Notification truth through the Product API. Realtime, if later admitted, is only a disposable wake-up/refetch optimization.

**Spec:** `D6-R2-NOTIF-01-NOTIFICATION-ARCHITECTURE-DESIGN.md`

## Global constraints

- Repository current authority > this plan > conversation history.
- PR #61 remains Draft/unmerged unless explicitly authorized.
- One coherent gate lands before the next; every stage ends with exact-head verification and operator adjudication.
- No Product/runtime implementation before D9.
- Notification != Work != Audit != authorization != source truth != acknowledgement.
- Initial mandatory trigger only: Operational Work becomes explicitly assigned or reassigned to an exact human Principal.
- Read/unread/archive state affects Notification only.
- Notification never grants source access; source reads re-authorize current access.
- No `seen`, numeric unread aggregate, mark-all-read, bulk archive, preferences, subscriptions, digest, e-mail, push, generic template engine, generic entity graph, external broker, generic EventStore, or cross-Organization Inbox baseline.
- No source-owner + Notification cross-owner business transaction. Atomicity is source fact + River `InsertTx`; Notification persistence is consumer-owner-local.
- No Kafka/RabbitMQ/NATS/Redis queue/second generic outbox.
- Persistent PostgreSQL Notification state + Product API read is the correctness baseline.
- B00 global IA stays LOCKED; only the topbar utility slot may reopen for the bell.
- Do not preserve the current 99 operations / 30 Permissions by preference if accepted Notification capability requires a minimal increase.
- Preserve exact `docs/roadmap.md` gate markers and keep bootstrap authority pack < 20,480 bytes.

## Planned files

Create as gates mature:

```text
docs/engineering/rebaseline/D6-R2-NOTIF-01-AUTHORITY-AMENDMENT.md
docs/engineering/rebaseline/D5-R3-PERSONAL-NOTIFICATIONS.md
contracts/api/product/paths-notifications.yaml
scripts/verify-product-notifications.mjs
docs/engineering/rebaseline/D6-R3-PERSONAL-NOTIFICATIONS-FRONTEND.md
qualification/d6-r2-wireframes/b02-notifications.html
docs/engineering/rebaseline/D7-R2-PERSONAL-NOTIFICATIONS-RUNTIME.md
docs/engineering/rebaseline/D8-R3-PERSONAL-NOTIFICATIONS-PROOF.md
```

Modify as required:

```text
docs/index.md
docs/roadmap.md
contracts/api/product/openapi.yaml
contracts/api/product/components.yaml
package.json
docs/engineering/rebaseline/D6-R2-P8-BLOCK-LEDGER.md
qualification/d6-r2-wireframes/b00-app-shell.html
PR #61 body
```

Do **not** create migrations, Go Notification packages, React runtime components, River worker code, SSE handlers, or broker/outbox infrastructure in this plan.

---

## Task 1 — D0 Product-scope amendment

**Files:** create/update `D6-R2-NOTIF-01-AUTHORITY-AMENDMENT.md`, update `docs/roadmap.md`.

**Produces:** Product-level admission only; no owner/wire/runtime selection.

- [ ] Revalidate PR #61, branch, HEAD, main, and CI before writes.
- [ ] Create the D0 section recording exactly:

```text
Personal Notification Inbox is Launch-V1 Product scope.
Purpose = personal awareness of committed MPC facts for an exact human Principal.
MPC owns Notification read/archive state only.
Notification is not Work/Audit/authorization/source truth/acknowledgement.
Notification grants no source access.
Initial trigger = Work assignment/reassignment to exact human Principal.
No e-mail/push/preferences/subscriptions/platform-notification requirement.
No exact unread count requirement; bell presence indicator is enough baseline.
```

- [ ] Add D0 negative controls: no Work-as-Inbox, no source resolution by Notification mutation, no cross-Organization Inbox, no A/S human Inbox, no generic notification platform.
- [ ] Update roadmap to `NOTIF-01 D0 CANDIDATE`; keep global D0 `ACCEPTED / CLOSED` and do not open D1.
- [ ] Run `npm run gate:full`; expected OAD remains 99/30/H-A-S because D5 is untouched.
- [ ] Stop for operator D0 approval.
- [ ] After approval, mark D0 amendment accepted, commit, and rerun full gate on exact HEAD.

---

## Task 2 — D1 supporting semantic owner + edge

**Files:** update `D6-R2-NOTIF-01-AUTHORITY-AMENDMENT.md`, roadmap.

**Produces:** one new semantic owner and the smallest edge set.

- [ ] Define `Personal Notifications` ownership:

```text
owns: Notification lifecycle, exact recipient, read/unread, archive state,
      Notification-local source correlation/deduplication
not owns: source meaning, Work lifecycle, Governance, Audit, source access,
          external delivery channels
```

- [ ] Admit only:

```text
Operational Work -> Personal Notifications
```

for a committed assignment/reassignment occurrence targeted to an exact human Principal.

- [ ] Explicitly reject generic `AnyDomain -> Notifications`, producer writes to Notification private state, Notification-triggered source mutation, workflow/event-hub authority, and a universal entity graph.
- [ ] Record that future Authorization/action-outcome triggers require a new trigger-census proof before a new D1 edge exists.
- [ ] Run full gate, stop for D1 approval, then commit accepted D1 amendment and reverify.

---

## Task 3 — D2 identity/state/isolation model

**Files:** update `D6-R2-NOTIF-01-AUTHORITY-AMENDMENT.md`, roadmap.

**Produces:** canonical Notification identity/state grammar.

- [ ] Define one opaque, non-reusable MPC `NotificationID` scoped to Organization.
- [ ] Freeze the initial semantic state:

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

- [ ] Keep `source_work_ref` typed to Work for V1. Do not create a generic union until a second trigger owner is actually accepted.
- [ ] Require `source_occurrence_discriminator` to identify the exact assignment/reassignment occurrence used for semantic deduplication; D3 binds its producer meaning.
- [ ] `read_at=null` = unread; non-null = read. `archived_at=null` = active; non-null = archived. No `seen`, delivered, dismissed, acknowledged, resolved, severity, or priority.
- [ ] Recipient Inbox behavior is human-Principal only; historical Notification remains even if later access/membership changes.
- [ ] Cross-Organization state/reference access is forbidden; source reference never grants access.
- [ ] Run full gate, stop for D2 approval, commit accepted D2 amendment, reverify.

---

## Task 4 — D3 trigger census + recoverable propagation

**Files:** update `D6-R2-NOTIF-01-AUTHORITY-AMENDMENT.md`, roadmap.

**Produces:** exact event-worthiness and recovery semantics without transport-platform expansion.

- [ ] Run the trigger census with expected baseline disposition:

```text
Work assignment/reassignment to exact human Principal -> ADMIT
Authorization Decision                              -> DEFER
async action outcome                                -> DEFER
sale/listing/shipment CRUD/state changes             -> REJECT generic fan-out
all Work opened to role/group                        -> REJECT without exact recipient semantics
```

- [ ] Define the Work-owned committed assignment/reassignment occurrence and the smallest stable occurrence discriminator. Do not introduce a universal EventID.
- [ ] Classify the admitted edge as `E` after Work commit. The occurrence must carry exact Organization, Work reference, exact new human recipient, and stable discriminator sufficient for Notification creation.
- [ ] If D3 determines a current Work Q is needed to avoid notifying a superseded recipient, record the exact revalidation condition; do not erase legitimate historical occurrence semantics.
- [ ] Bind `at-least-once / repeat-safe`, no exactly-once claim, duplicate-safe Notification creation, detectable/recoverable loss, and transport/job state not being Notification history.
- [ ] Reject broker, generic EventStore, universal envelope, event-per-CRUD, and direct producer writes to Notification state.
- [ ] Run full gate, stop for D3 approval, commit and reverify. D5 stays blocked until D0–D3 are accepted.

---

## Task 5 — D5 exact Product wire + executable proof

**Files:** create `D5-R3-PERSONAL-NOTIFICATIONS.md`, `paths-notifications.yaml`, `verify-product-notifications.mjs`; modify OAD/components/package/index/roadmap.

**Produces:** exact operationIds, Permission mapping, final Product census, and generated-wire proof.

- [ ] Read only the D5 operation/schema/collection/access subpacks required by this wire.
- [ ] Compare exactly two wire shapes:

```text
A. explicit state-transition operations for read/unread/archive/unarchive + list
B. one bounded Notification-state update/PATCH operation + list
```

Choose the smaller D5-conformant form that preserves concurrency/preconditions, auditability, client-class/Permission clarity, and avoids a generic patch/action surface.

- [ ] Freeze an operator-reviewable D5-R3 table **before** editing OAD. Every admitted operation row must contain:

```text
operationId
method/path
Q or C
semantic owner = PersonalNotifications
ordinary Permission
principal kinds
Organization/recipient scope
precondition/concurrency semantics
success + material errors
```

- [ ] Explicitly reject public Notification creation, delete, recipient selection, admin-wide search, bulk, mark-all-read, preferences, exact unread-count aggregate, and machine Inbox.
- [ ] Stop for operator approval of the D5 semantic wire table.
- [ ] After approval, write `scripts/verify-product-notifications.mjs` first and run it RED. The verifier must use the exact approved operationId/Permission sets and assert:

```text
all PersonalNotifications Product operations match approved set
all browser Inbox operations are H-only
no public CreateNotification
no generic entity_type/entity_id
no seen field
no exact unread-count operation
no caller-selectable recipient
```

- [ ] Add canonical paths/schemas to `paths-notifications.yaml`, `components.yaml`, and `openapi.yaml`; no parallel DTO/schema home.
- [ ] Add the verifier to `gate` / `gate:full` only after the canonical wire exists.
- [ ] Run:

```bash
node scripts/verify-product-notifications.mjs
npm run gate:full
```

Expected: Notification proof PASS plus all existing OAD/auth/Performance/operational-read/generator proofs PASS.

- [ ] Record the **actual** new Product operation and Permission counts from the accepted wire. Never preserve 99/30 by preference.
- [ ] Final D5-R3 exact-head CI must be green before D6 opens.

---

## Task 6 — D6 bell utility slot + Inbox frontend authority

**Files:** create `D6-R3-PERSONAL-NOTIFICATIONS-FRONTEND.md`, update P8 ledger/B00, create B02 Inbox HTML, update roadmap.

**Produces:** exact frontend interaction authority only; no React implementation.

- [ ] Bind bell visibility, unread-presence query, preview, full Inbox, read/unread, archive/unarchive, and source navigation to the exact D5 operations/Permission.
- [ ] Reopen B00 **only** for the topbar utility slot. Preserve locked sidebar groups, Organization/Installation semantics, dimensions, responsive laws, and access behavior.
- [ ] B00 candidate states must prove:

```text
known no unread -> bell, no dot
known unread exists -> bell + dot
Notification knowledge unavailable -> no false empty/zero inference
Organization switch -> close preview + rescope query
```

- [ ] Render updated executable `b00-app-shell.html`; stop for operator utility-slot `LOCKED` adjudication.
- [ ] Add P8 block `B02 — Personal Notifications` only after B00 utility slot is locked.
- [ ] Render `b02-notifications.html` with active Inbox (unread/read), archived view, unavailable state, bell preview, source navigation, read/unread and archive/unarchive. No count, bulk, preferences, mark-all-read, or source mutation.
- [ ] Stop for operator B02 lock.
- [ ] Run full gate and commit D6-R3 only after required operator locks.

---

## Task 7 — D7 runtime amendment over existing PostgreSQL + River

**Files:** create `D7-R2-PERSONAL-NOTIFICATIONS-RUNTIME.md`, update roadmap.

**Produces:** target runtime mechanics/proof contract; still no Product runtime code.

- [ ] Bind Work owner commit to existing D7-C pattern:

```text
BEGIN Work owner transaction
  commit assignment/reassignment occurrence
  River InsertTx durable PersonalNotifications reaction
COMMIT
```

No direct Notification write and no second outbox.

- [ ] Consumer worker: enter exact Organization, re-enter PersonalNotifications boundary, apply semantic duplicate predicate over accepted occurrence discriminator, create Notification in owner-local transaction, transactionally complete River work when useful.
- [ ] Apply existing D7-B Organization/RLS/composite-FK/runtime-role laws to Notification persistence.
- [ ] Make no exactly-once claim; River uniqueness is optimization, semantic idempotency is correctness.
- [ ] Dispose realtime based on D6 evidence:

```text
DEFER, if ordinary query/focus refetch is sufficient for Launch V1
or
OPTIONAL WAKE-UP, if proven useful: same-origin SSE only invalidates/refetches Inbox truth
```

If optional wake-up is admitted, PostgreSQL LISTEN/NOTIFY may be best-effort IPC only **after** canonical Notification persistence; it is not truth/history/count and cannot be required for correctness. No WebSocket/broker baseline.

- [ ] D7 falsifiers must include atomic source-state+River handoff, rollback leaving no job, duplicate/redelivery safety, cross-Organization job tampering, lost optional wake-up, and River state not becoming Notification state.
- [ ] Run full gate, stop for operator D7 approval, commit and reverify.

---

## Task 8 — D8 bounded Notification falsifier

**Files:** create `D8-R3-PERSONAL-NOTIFICATIONS-PROOF.md`, update roadmap.

**Produces:** one composed Notification falsifier, not an exhaustive catalog.

- [ ] Define `GF-N01`:

```text
human assigns/reassigns Work
-> Work commits occurrence
-> durable River reaction exists
-> one semantic Notification materializes
-> Inbox returns unread
-> read mutates Notification only
-> archive mutates Notification only
-> source click re-authorizes Work access
```

- [ ] Adversarial branches: duplicate/redelivery, consumer crash before/after Notification commit, source access revoked, cross-Organization read/mutation, Inbox read unavailable, optional realtime signal loss, machine Principal human-Inbox attempt.
- [ ] Do not reopen GF-01/GF-02/GF-03/SR-01 unless a concrete shared invariant was actually changed.
- [ ] Run full gate, stop for operator D8-R3 approval, commit and reverify.

---

## Task 9 — reconcile D6-R2 and resume frontend planning

**Files:** update D6-R2 closure, P5 inventory, P8 ledger, roadmap, PR #61 body.

**Produces:** current screen/flow/operation inventory and the next legitimate D6-R2 gate.

- [ ] Add the bounded Personal Notification need/flow and re-run only delta coverage against the final Product operation census.
- [ ] Add only accepted Notification utility/full-Inbox surfaces to P5; preserve all previously locked B00/B01 decisions except the bounded bell delta already ratified.
- [ ] Restore B10 from “suspended for NOTIF-01” to its previously valid `OFERTA > Preparação` candidate status; do not silently lock it.
- [ ] Set roadmap exact next action to the next unresolved block in the existing P8 sequence.
- [ ] Run final `npm run gate:full` and record exact HEAD, CI, Product operation count, Permission count, H/A/S profile, generated-projection proof, Notification proof, operational-read proof, bootstrap bytes, durable-doc reachability, and legacy runtime population.
- [ ] Synchronize PR #61; keep Draft/unmerged.

---

## Self-review

**Spec coverage:** supporting owner, Work assignment/reassignment trigger, durable state, Organization/H recipient isolation, read/archive independence, recoverable River propagation, no broker, topbar bell, full Inbox, no numeric count, re-authorization on click, optional realtime only, D8 proof, and D6-R2 resume are all assigned to explicit tasks.

**No placeholders:** Deferred items are closed later-stage authority choices with defined candidate sets and operator gates; there is no `TBD/TODO`, generic “add error handling”, or speculative platform work.

**Naming:** narrative = `Personal Notifications`; Product semantic-owner metadata = `PersonalNotifications`; initial fields = `notification_id`, `organization_id`, `recipient_principal_id`, `kind`, `source_work_ref`, `source_occurrence_discriminator`, `created_at`, `read_at`, `archived_at`, `revision`. Purely mechanical D5 naming changes require explicit D5-R3 rationale.
