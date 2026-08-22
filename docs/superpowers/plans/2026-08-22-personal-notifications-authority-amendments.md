# Personal Notifications Authority Amendments Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Admit Marketplace Central Personal Notifications as a Launch-V1 Product responsibility, carry that decision coherently through D0→D8 authority, add the topbar bell + organization-scoped Inbox to frontend authority, and preserve the existing modular-monolith/PostgreSQL/River architecture without Product implementation before D9.

**Architecture:** Personal Notifications is a small Organization-scoped, human-Principal-targeted supporting semantic owner. Source owners commit their own facts and atomically enqueue one durable River reaction; Personal Notifications creates idempotent Notification state in its own transaction. The browser reads canonical Notification state through the Product API; the topbar bell is an Inbox entry point, while realtime remains an optional disposable wake-up optimization rather than truth.

**Tech Stack:** Documentation authority under `docs/engineering/rebaseline`, OpenAPI 3.1.2, Node verification scripts, React + TypeScript frontend authority, TanStack Query/Router, Go modular monolith target, PostgreSQL, River durable work, optional same-origin SSE wake-up seam only if later justified.

**Spec:** `docs/engineering/rebaseline/D6-R2-NOTIF-01-NOTIFICATION-ARCHITECTURE-DESIGN.md`

## Global Constraints

- Product implementation remains **BLOCKED UNTIL accepted D9** throughout this plan.
- Repository current authority > this plan > chat/history.
- Keep PR #61 Draft/unmerged unless the operator explicitly authorizes merge.
- One coherent gate lands before the next; every task ends with exact-head verification.
- Personal Notifications owns awareness state only; it is not Work, Audit, Governance authorization, source truth, acknowledgement, or a cross-owner workflow owner.
- Notification is Organization-scoped and targets an exact human Principal `H`; A/S do not gain a human Inbox.
- Initial required trigger is only: Operational Work becomes explicitly assigned or reassigned to an exact human Principal.
- Reading/archiving Notification never mutates Work or any source owner.
- Notification never grants source access; navigation re-authorizes the source at click time.
- No `seen`, numeric unread aggregate, mark-all-read, bulk archive, preferences, subscriptions, digests, e-mail, push, generic template engine, generic entity graph, external broker, generic EventStore, or cross-Organization Inbox baseline.
- No source-owner + Notification cross-owner business transaction. Source commit + River `InsertTx` durable handoff is atomic; Notification persistence is consumer-owner-local.
- No Kafka, RabbitMQ, NATS, Redis queue, or second generic outbox.
- Realtime is not correctness-critical. Persistent PostgreSQL Notification state + Product API read is the baseline.
- B00 global IA remains LOCKED. Only the bounded topbar utility slot may reopen for the bell.
- The current 99-operation / 30-Permission census is not protected by preference once Notification is accepted scope; any increase must be the smallest derivable D5 result.
- Preserve exact gate markers in `docs/roadmap.md`; do not reformat them casually.
- Keep bootstrap authority pack below 20,480 bytes.

---

## File Structure Map

The plan deliberately avoids rewriting accepted historical stage documents. It uses one NOTIF-01 cross-stage authority amendment plus bounded owning-stage wire/runtime/proof amendments, matching the repository's existing repair pattern.

**Create:**

- `docs/engineering/rebaseline/D6-R2-NOTIF-01-AUTHORITY-AMENDMENT.md` — accepted D0/D1/D2/D3 semantic amendments and trigger census; this is the upstream semantic authority for the feature.
- `docs/engineering/rebaseline/D5-R3-PERSONAL-NOTIFICATIONS.md` — exact Product operation, Permission, wire, error, concurrency and client-class amendment.
- `contracts/api/product/paths-notifications.yaml` — Notification Product paths only.
- `scripts/verify-product-notifications.mjs` — executable Notification wire/negative-control proof.
- `docs/engineering/rebaseline/D6-R3-PERSONAL-NOTIFICATIONS-FRONTEND.md` — bell, Inbox, route/state/action/owner/Permission frontend authority.
- `qualification/d6-r2-wireframes/b02-notifications.html` — executable low-fidelity Inbox prototype after the B00 utility-slot gate.
- `docs/engineering/rebaseline/D7-R2-PERSONAL-NOTIFICATIONS-RUNTIME.md` — target runtime realization over PostgreSQL + River and optional realtime seam disposition.
- `docs/engineering/rebaseline/D8-R3-PERSONAL-NOTIFICATIONS-PROOF.md` — smallest composed Notification falsifier.

**Modify:**

- `docs/engineering/rebaseline/D6-R2-NOTIF-01-NOTIFICATION-ARCHITECTURE-DESIGN.md` — already operator-approved; only add links/status as gates close.
- `docs/index.md` — route the active Notification authority pack without bloating bootstrap.
- `docs/roadmap.md` — sole mutable status/next-action authority.
- `contracts/api/product/openapi.yaml` — route the accepted Notification Product paths.
- `contracts/api/product/components.yaml` — canonical Notification schemas only; no parallel DTO model.
- `package.json` — append the Notification proof to the full gate only after the wire is selected.
- `docs/engineering/rebaseline/D6-R2-P8-BLOCK-LEDGER.md` — bounded B00 utility-slot reopen, B02 status, and resume gates.
- `qualification/d6-r2-wireframes/b00-app-shell.html` — add only the approved bell utility slot after D5/D6 authority exists.
- PR #61 body — keep exact current checkpoint and verification synchronized.

**Do not create:**

- Product runtime source code;
- migrations;
- React components;
- Go Notification packages;
- River worker code;
- SSE handlers;
- broker/outbox infrastructure.

Those belong to post-D9 implementation, not this authority plan.

---

### Task 1: Ratify D0 Product Scope Amendment

**Files:**
- Create: `docs/engineering/rebaseline/D6-R2-NOTIF-01-AUTHORITY-AMENDMENT.md`
- Modify: `docs/roadmap.md`

**Interfaces:**
- Consumes: operator-approved Notification design at `D6-R2-NOTIF-01-NOTIFICATION-ARCHITECTURE-DESIGN.md`.
- Produces: an explicit D0 amendment stating that Personal Notification Inbox is Product 1.0 scope and defining its product-level authority/non-goals without choosing owner identity, wire, or runtime.

- [ ] **Step 1: Revalidate repository state before any write**

Check PR #61 is open, Draft, unmerged, branch is `stage/d6-r2-frontend-realization`, and current HEAD/CI are green. Stop if repository authority changed materially.

- [ ] **Step 2: Write the D0 section of the authority amendment**

Create the amendment document with a D0 section that records exactly:

```text
Product 1.0 adds Personal Notification Inbox.
Purpose: personal awareness of committed MPC facts for an exact human Principal.
Authority class: OWN for Notification read/archive state only.
Notification != Work != Audit != authorization != source truth != acknowledgement.
No source access is granted by Notification.
Initial launch trigger: Work assigned/reassigned to exact human Principal.
No e-mail/push/preferences/subscriptions/generic notification platform launch gate.
No exact unread count requirement; bell presence indicator is sufficient baseline.
```

Do not name HTTP operations, Permission strings, tables, events, River jobs, or SSE in the D0 section.

- [ ] **Step 3: Add D0 negative controls**

The amendment must explicitly reject:

```text
using Operational Work as the Inbox
using Notification read/archive as Work/source resolution
cross-Organization Inbox
A/S human Inbox
email/push as Launch-V1 requirement
generic notification/subscription platform
```

- [ ] **Step 4: Update roadmap to `NOTIF-01 D0 amendment candidate`**

Keep D0 globally `ACCEPTED / CLOSED`; route the bounded amendment as candidate authority and set the exact next action to operator adjudication of D0 only. Do not open D1 yet.

- [ ] **Step 5: Verify exact-head repository gate**

Run through CI-equivalent repository full gate:

```bash
npm run gate:full
```

Expected: `gate: PASS`, existing Product OAD still 99/30/H-A-S, Notification wire unchanged.

- [ ] **Step 6: Operator gate**

Stop. Present only the D0 amendment for operator approval. Do not proceed to D1 until explicit approval.

- [ ] **Step 7: Commit status after approval**

Commit the approved D0 section/status with a message equivalent to:

```bash
git commit -m "docs(notif): ratify D0 notification scope amendment"
```

Run `npm run gate:full` again on the final D0 HEAD.

---

### Task 2: Ratify D1 Supporting Semantic Owner and Edges

**Files:**
- Modify: `docs/engineering/rebaseline/D6-R2-NOTIF-01-AUTHORITY-AMENDMENT.md`
- Modify: `docs/roadmap.md`

**Interfaces:**
- Consumes: accepted D0 Notification scope.
- Produces: `Personal Notifications` as the only new supporting semantic owner and a minimal D1 edge set.

- [ ] **Step 1: Add the D1 owner definition**

Record:

```text
Personal Notifications owns:
- Notification lifecycle as personal awareness state
- exact recipient Principal
- read/unread state
- archived/not-archived state
- Notification-local source correlation needed for navigation/deduplication

Personal Notifications does not own:
- source business meaning
- Work responsibility/assignment/escalation/resolution
- Governance authorization
- source access
- Audit
- delivery channels
```

- [ ] **Step 2: Admit only the proved semantic edge**

The baseline D1 edge set is exactly:

```text
Operational Work -> Personal Notifications
```

Meaning: Work may communicate a committed assignment/reassignment occurrence that makes the Work personally relevant to an exact human Principal. Personal Notifications decides Notification-local idempotency/state; it does not reinterpret Work truth.

Do not add generic `AnyDomain -> Notifications` or Notifications→all-domains edges.

- [ ] **Step 3: Record future-edge reopen law**

Future Authorization/action-outcome/etc. triggers require the trigger census to prove an exact recipient, independent personal-awareness value, and stable occurrence discriminator. Until then, they are not D1 edges.

- [ ] **Step 4: Add D1 forbidden-boundary controls**

Reject:

```text
Personal Notifications becoming a workflow/event hub
producer writing Notification private state
Notification mutating source owner
Notification becoming generic user-task domain
Notification becoming a platform-wide polymorphic entity graph
```

- [ ] **Step 5: Update roadmap to D1 candidate and verify**

Run:

```bash
npm run gate:full
```

Expected: existing Product OAD unchanged; no D2/D3/D5 semantics claimed yet.

- [ ] **Step 6: Operator gate**

Stop for explicit D1 approval.

- [ ] **Step 7: Commit approved D1 gate and reverify**

Use a commit message equivalent to:

```bash
git commit -m "docs(notif): ratify D1 notification owner boundary"
```

Run `npm run gate:full` on the final D1 HEAD.

---

### Task 3: Ratify D2 Notification Identity, Ownership, Recipient, and Source Reference

**Files:**
- Modify: `docs/engineering/rebaseline/D6-R2-NOTIF-01-AUTHORITY-AMENDMENT.md`
- Modify: `docs/roadmap.md`

**Interfaces:**
- Consumes: accepted `Personal Notifications` owner and `Operational Work -> Personal Notifications` edge.
- Produces: canonical Notification identity/state model and Organization/Principal isolation semantics.

- [ ] **Step 1: Define canonical Notification identity**

Record one MPC-owned opaque `NotificationID`, Organization-owned and non-reusable. No ID encodes Principal, Work, kind, or timestamp.

- [ ] **Step 2: Define the minimal D2 state**

Baseline state is exactly:

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

For the first accepted trigger, `source_work_ref` is a typed Work reference, not a generic entity reference. Do not create a union until a second trigger owner is actually accepted.

`source_occurrence_discriminator` must identify the particular accepted Work assignment/reassignment occurrence sufficiently for semantic deduplication. The D3 task will bind it to the committed producer occurrence; D2 only requires stability and non-guessing.

- [ ] **Step 3: Define recipient law**

`recipient_principal_id` must resolve to an exact current/historical MPC Principal identity whose kind is human for Inbox behavior. Notification retention does not rewrite when the Principal later loses membership/access.

- [ ] **Step 4: Define state semantics**

Use:

```text
read_at = null       -> unread
read_at != null      -> read
archived_at = null   -> active Inbox
archived_at != null  -> archived
```

No `seen`, delivered, dismissed, acknowledged, resolved, or severity field.

- [ ] **Step 5: Define isolation/access invariants**

Notification rows are Organization-owned and recipient-owned for Product reads/mutations. Cross-Organization references are forbidden. Notification source reference does not grant source access.

- [ ] **Step 6: Update roadmap, gate, and operator adjudication**

Run:

```bash
npm run gate:full
```

Stop for D2 approval before D3.

- [ ] **Step 7: Commit approved D2 gate and reverify**

Use a commit message equivalent to:

```bash
git commit -m "docs(notif): ratify D2 notification identity model"
```

Run the full gate again.

---

### Task 4: Ratify D3 Trigger Census and Recoverable Communication

**Files:**
- Modify: `docs/engineering/rebaseline/D6-R2-NOTIF-01-AUTHORITY-AMENDMENT.md`
- Modify: `docs/roadmap.md`

**Interfaces:**
- Consumes: D1 Work→Notifications edge and D2 Notification occurrence model.
- Produces: exact event-worthiness, duplicate/order/recovery semantics for the initial trigger; no transport selection beyond accepted D3 semantics.

- [ ] **Step 1: Run the Notification Trigger Census**

Evaluate candidate triggers against the five accepted criteria in the design spec. Baseline expected disposition:

```text
Work assignment/reassignment to exact human Principal -> ADMIT
Authorization Decision -> DEFER pending exact-recipient/useful-awareness proof
async action outcome -> DEFER pending exact-recipient/useful-awareness proof
shipment/listing/sale CRUD/state changes -> REJECT as generic notification fan-out
all Work openings to role/group -> REJECT without exact recipient/subscription semantics
```

Document evidence/rationale for each disposition.

- [ ] **Step 2: Define the committed producer occurrence**

D3 must define a producer-owned Work assignment/reassignment occurrence whose stable discriminator is recoverable from durable Work meaning/history and can be used by Notifications for duplicate suppression. Do not invent a universal EventID; use the smallest Work-owned occurrence discriminator consistent with existing D3 rules.

- [ ] **Step 3: Classify communication**

Record the initial edge as:

```text
Operational Work commits assignment/reassignment occurrence
  -> E to wake Personal Notifications
  -> no Q required for Notification creation if the occurrence already carries exact recipient + stable Work reference + occurrence discriminator
  -> source navigation later reads Work through ordinary Product authority
```

If current Work truth must be revalidated to avoid notifying a superseded recipient, specify the exact Q/revalidation condition; do not make every historical occurrence depend on current mutable assignment if that would erase legitimate awareness history.

- [ ] **Step 4: Define delivery/recovery semantics**

Binding semantics:

```text
at-least-once / repeat-safe
no exactly-once claim
late/duplicate delivery cannot duplicate Notification business state
loss is detectable/recoverable
transport/job state is not Notification history
```

- [ ] **Step 5: Explicitly reject broker/event-platform expansion**

No Kafka/RabbitMQ/NATS, no generic EventStore, no universal event envelope, no event-per-CRUD, no source-owner direct Notification table writes.

- [ ] **Step 6: Update roadmap, run gate, operator adjudication**

Run:

```bash
npm run gate:full
```

Stop for D3 approval. D5 remains blocked until D0–D3 are all accepted.

- [ ] **Step 7: Commit approved D3 gate and reverify**

Use a commit message equivalent to:

```bash
git commit -m "docs(notif): ratify D3 notification propagation"
```

Run full gate again.

---

### Task 5: Derive and Prove the D5 Product Wire

**Files:**
- Create: `docs/engineering/rebaseline/D5-R3-PERSONAL-NOTIFICATIONS.md`
- Create: `contracts/api/product/paths-notifications.yaml`
- Modify: `contracts/api/product/openapi.yaml`
- Modify: `contracts/api/product/components.yaml`
- Create: `scripts/verify-product-notifications.mjs`
- Modify: `package.json`
- Modify: `docs/index.md`
- Modify: `docs/roadmap.md`

**Interfaces:**
- Consumes: accepted D0–D3 NOTIF-01 authority.
- Produces: one canonical Product wire for own-Inbox reads/mutations, exact ordinary Permission mapping, exact operation census, generated-projection proof, and hard negative controls.

- [ ] **Step 1: Read only the D5 operation/schema/access subpack needed for the wire decision**

Use current D5 API + Operation Admission Matrix + W1/W2/W3/W4 as needed. Do not import frontend preferences into D5.

- [ ] **Step 2: Compare the two bounded wire shapes**

Evaluate only:

```text
A. explicit transition operations for read/unread/archive/unarchive
B. one bounded Notification-state PATCH/update operation plus list
```

Select the smaller D5-conformant wire that preserves precondition/concurrency semantics, auditability, discoverability, client-class/Permission clarity, and does not create a generic patch/action API. Record the rejected alternative and why.

Do not admit public Notification creation, delete, bulk mutation, preferences, mark-all-read, exact-count aggregate, admin search, or machine Inbox.

- [ ] **Step 3: Freeze exact operation/Permission table in D5-R3 before editing OAD**

The D5-R3 document must contain a closed table with for every admitted Notification operation:

```text
operationId
method + path
Q/C classification
semantic owner = PersonalNotifications
required ordinary Permission
principal kinds
Organization/recipient scope
request body/precondition semantics
success and material error outcomes
```

After this table is operator-approved, its exact operation IDs and Permission strings become the expected values in the verifier. No tolerant or inferred census is allowed after approval.

- [ ] **Step 4: Operator gate on D5 semantic wire before OAD edits**

Stop and obtain explicit approval of the operation/Permission table.

- [ ] **Step 5: Write the failing Notification wire verifier**

Create `scripts/verify-product-notifications.mjs` so it bundles the canonical OAD and asserts at least:

```javascript
const notificationOps = operations(document)
  .filter(({ operation }) => operation['x-mpc-semantic-owner'] === 'PersonalNotifications');
assert(notificationOps.length === EXPECTED_NOTIFICATION_OPERATION_IDS.length,
  'Notification operation census mismatch');
sameSet(notificationOps.map(({ operation }) => operation.operationId),
  EXPECTED_NOTIFICATION_OPERATION_IDS,
  'Notification operationIds');
for (const { operation } of notificationOps) {
  assert((operation['x-mpc-principal-kinds'] ?? []).every((k) => k === 'H'),
    `${operation.operationId} must be human-only`);
}
assert(!operations(document).some(({ operation }) => operation.operationId === 'CreateNotification'),
  'public CreateNotification is forbidden');
```

Also assert exact Permission strings from the approved table, no generic `entity_type/entity_id`, no `seen`, no numeric unread aggregate operation, and no public recipient selector.

- [ ] **Step 6: Run verifier before OAD change and confirm RED**

Run:

```bash
node scripts/verify-product-notifications.mjs
```

Expected: FAIL because the approved Notification Product wire is not yet present.

- [ ] **Step 7: Implement the canonical OAD minimally**

Add only the approved paths/schemas and update canonical `openapi.yaml` routing. Use `components.yaml` as the shared schema authority; do not create handwritten frontend DTOs or a second schema home.

Schemas must preserve:

```text
NotificationID
Organization scope
recipient Principal reference where response semantics require it
closed notification kind
Work source reference + source occurrence discriminator
read_at?
archived_at?
revision/concurrency carrier as required by D5
```

- [ ] **Step 8: Add verifier to full gate and run GREEN**

Run:

```bash
node scripts/verify-product-notifications.mjs
npm run gate:full
```

Expected: Notification proof PASS; all existing Product/auth/performance/operational-read proofs remain PASS; generated TypeScript/Go projections remain deterministic/valid.

- [ ] **Step 9: Record the new exact Product census**

Update roadmap and D5-R3 with the actual post-Notification operation and ordinary-Permission counts produced by the accepted wire. Do not preserve 99/30 artificially and do not estimate counts from memory.

- [ ] **Step 10: Commit and exact-head verify**

Commit in logical slices, ending with a final D5-R3/OAD proof commit. After the last write, run full CI and record exact HEAD + checks in PR #61.

---

### Task 6: Reopen D6 Only for Bell Utility Slot + Personal Inbox

**Files:**
- Create: `docs/engineering/rebaseline/D6-R3-PERSONAL-NOTIFICATIONS-FRONTEND.md`
- Modify: `docs/engineering/rebaseline/D6-R2-P8-BLOCK-LEDGER.md`
- Modify: `qualification/d6-r2-wireframes/b00-app-shell.html`
- Create: `qualification/d6-r2-wireframes/b02-notifications.html`
- Modify: `docs/roadmap.md`

**Interfaces:**
- Consumes: canonical Notification OAD/Permission surface from D5-R3.
- Produces: exact bell/Inbox interaction authority and two operator-reviewed structural prototypes without changing the locked global sidebar IA.

- [ ] **Step 1: Map Notification user needs to exact Product operations**

D6-R3 must bind:

```text
bell visibility -> current access context / exact Notification Permission
unread dot -> authoritative existence query/read semantics from D5; never page-count inference
bell preview -> bounded recent own-Inbox read
full Inbox -> organization-scoped Notification collection
read/unread -> exact D5 operation
archive/unarchive -> exact D5 operation
source click -> source route + current source authorization
```

- [ ] **Step 2: Reopen B00 only for topbar utility slot**

Preserve locked sidebar groups, Organization selector, Installation semantics, dimensions, responsive laws, and access states. Add only a bell utility slot to the header.

The structural candidate must prove:

```text
no unread -> bell without dot
known unread exists -> bell + dot
Notification knowledge unavailable -> no false known-empty/known-zero signal
Organization switch -> preview closes and query scope changes
```

- [ ] **Step 3: Render corrected B00 HTML candidate**

Modify only `b00-app-shell.html`. Do not add Inbox content to the shell itself.

- [ ] **Step 4: Operator B00 utility-slot gate**

Stop. Only the operator can restore B00 `LOCKED` after viewing the bell placement/responsive behavior.

- [ ] **Step 5: Define B02 Personal Notifications block**

Add B02 to the P8 ledger. The Inbox candidate must include:

```text
active Inbox: unread + read
archived view
knowledge unavailable/request failure
recent preview from bell
source navigation affordance
read/unread mutation
archive/unarchive mutation
no source mutation
```

No numeric unread count, bulk actions, preferences, or “mark all read”.

- [ ] **Step 6: Render `b02-notifications.html`**

Build executable low-fidelity HTML only. It decides shell-relative placement, preview/full-Inbox structure, action order, states, and responsive behavior; it does not decide palette/branding/final visual system.

- [ ] **Step 7: Operator B02 gate**

Stop for explicit candidate changes or `LOCKED` decision.

- [ ] **Step 8: Verify and commit D6 authority**

Run:

```bash
npm run gate:full
```

Commit D6-R3 + B00/B02 ledger changes only after operator locks the required artifacts.

---

### Task 7: Ratify D7 Runtime Realization Without New Infrastructure

**Files:**
- Create: `docs/engineering/rebaseline/D7-R2-PERSONAL-NOTIFICATIONS-RUNTIME.md`
- Modify: `docs/roadmap.md`

**Interfaces:**
- Consumes: D3 propagation contract + D5 Notification semantics + accepted D7-A/B/C.
- Produces: target runtime mechanics/proof contract only; no runtime code.

- [ ] **Step 1: Bind source commit to existing River `InsertTx` law**

Specify:

```text
Work owner transaction:
  commit assignment/reassignment occurrence
  InsertTx one durable PersonalNotifications reaction
COMMIT
```

No direct Notification-row write from Work and no second outbox.

- [ ] **Step 2: Bind consumer transaction/idempotency**

The consumer worker enters exact Organization scope, resolves exact human recipient, applies a uniqueness constraint/semantic duplicate predicate over the accepted source occurrence identity, writes Notification state in a PersonalNotifications-owned transaction, and completes River work transactionally where useful.

No exactly-once claim.

- [ ] **Step 3: Bind persistence/isolation expectations**

Notification persistence is Organization-owned and covered by the accepted D7-B RLS/composite-FK pattern. Runtime role remains non-owner/NOBYPASSRLS. Cross-Organization consumer tampering must fail.

- [ ] **Step 4: Dispose realtime explicitly**

Because realtime is not required for Notification correctness, choose one of two outcomes based on D6 evidence:

```text
DEFER: ordinary Query refetch/focus is sufficient for Launch-V1 baseline
or
ADMIT OPTIONAL WAKE-UP: same-origin SSE whose payload only invalidates/refetches Inbox truth
```

If optional wake-up is admitted, PostgreSQL LISTEN/NOTIFY may be used only as best-effort IPC after canonical Notification persistence; it must not be the sole record, count, history, or a mandatory operation capable of making Notification correctness depend on wake-up delivery.

No WebSocket/broker baseline.

- [ ] **Step 5: Define D7 falsifiers**

At minimum:

```text
source commit + missing River handoff impossible under one transaction
source rollback leaves no runnable reaction
consumer crash/redelivery does not duplicate Notification
cross-Organization job tampering blocked
lost optional realtime wake-up does not lose Notification
River state never becomes Notification business state
```

- [ ] **Step 6: Operator gate + full repository gate**

Run `npm run gate:full`, stop for operator approval, then commit D7-R2 and reverify exact HEAD.

---

### Task 8: Add the Smallest D8 Notification Falsifier

**Files:**
- Create: `docs/engineering/rebaseline/D8-R3-PERSONAL-NOTIFICATIONS-PROOF.md`
- Modify: `docs/roadmap.md`

**Interfaces:**
- Consumes: accepted D0–D7 Notification amendments.
- Produces: one composed Notification golden-path/failure proof contract that protects the whole feature without creating an exhaustive test catalog.

- [ ] **Step 1: Define GF-N01**

Compose:

```text
human assigns/reassigns Work
-> Work commits assignment occurrence
-> durable River reaction exists
-> Personal Notification materializes exactly once semantically
-> Inbox read exposes unread Notification
-> mark read changes Notification only
-> archive changes Notification only
-> source click re-authorizes current Work access
```

- [ ] **Step 2: Define adversarial branches**

Include:

```text
duplicate/redelivered River reaction
consumer crash before/after Notification commit
source access revoked after Notification creation
cross-Organization read/mutation attempt
Inbox read unavailable
optional realtime signal loss
machine Principal attempting human Inbox use
```

- [ ] **Step 3: Preserve D8 scope discipline**

Do not reopen GF-01/GF-02/GF-03/SR-01 except where a concrete shared invariant is materially touched. Notification proof is additive and bounded.

- [ ] **Step 4: Operator gate and repository verification**

Run `npm run gate:full`; stop for operator approval; commit and reverify.

---

### Task 9: Reconcile D6-R2 Inventory and Resume Frontend Planning

**Files:**
- Modify: `docs/engineering/rebaseline/D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md`
- Modify: `docs/engineering/rebaseline/D6-R2-P5-SCREEN-SURFACE-INVENTORY.md`
- Modify: `docs/engineering/rebaseline/D6-R2-P8-BLOCK-LEDGER.md`
- Modify: `docs/roadmap.md`
- Modify: PR #61 body

**Interfaces:**
- Consumes: all accepted NOTIF-01 amendments and locked B00/B02 frontend artifacts.
- Produces: a reconciled D6-R2 screen/need/operation inventory and exact next frontend block with no stale 99-operation assumptions.

- [ ] **Step 1: Re-run P1/P2/P3 coverage only for the bounded delta**

Add the Personal Notification user need/flow and prove the final Product operation families are all findable. Do not rerun unrelated accepted actor/flow design from preference.

- [ ] **Step 2: Update P5 inventory**

Add only the new utility/full-Inbox surface(s) actually accepted. Preserve all previously locked B00/B01 and existing route homes unless Notification authority materially requires a bounded delta.

- [ ] **Step 3: Reconcile B10 suspension**

Once NOTIF-01 is fully accepted and B00 utility slot/B02 are locked, return B10 to its previously valid `OFERTA > Preparação` candidate state. Do not silently mark B10 locked; its prior operator review gate remains separate.

- [ ] **Step 4: Update roadmap exact next action**

Set the next action to the next unresolved D6-R2 block according to the existing P8 sequence and operator locks. Do not skip directly to implementation-readiness or D9.

- [ ] **Step 5: Final full-gate verification**

Run:

```bash
npm run gate:full
```

Record exact final HEAD, CI, Product operation count, ordinary Permission count, H/A/S profile, generated projection proof, Notification proof, operational read proof, bootstrap bytes, durable-doc reachability, and legacy runtime population.

- [ ] **Step 6: Synchronize PR #61 body**

Keep it Draft/unmerged and describe NOTIF-01 as accepted authority only if every preceding operator gate is actually closed.

---

## Plan Self-Review

### Spec coverage

- Supporting semantic owner: Tasks 1–4.
- Work assignment/reassignment initial trigger: Tasks 1–4 and 8.
- Notification durable state/read/archive: Tasks 2–5.
- Organization/human Principal isolation: Tasks 3, 5, 7, 8.
- No Work/source mutation: Tasks 1, 2, 5, 6, 8.
- Recoverable River propagation: Tasks 4, 7, 8.
- No external broker/outbox: Tasks 4 and 7.
- Topbar bell + unread presence dot: Task 6.
- Full Inbox utility route: Task 6.
- No numeric unread aggregate: Tasks 5 and 6.
- Current source authorization on click: Tasks 3, 5, 6, 8.
- Realtime as optional non-truth seam: Tasks 6–8.
- No Product implementation before D9: global constraint and every runtime/frontend task.
- D8 falsifier: Task 8.
- D6-R2 reconciliation/resume: Task 9.

### Placeholder scan

No `TBD`, `TODO`, “implement later”, generic error-handling placeholder, or speculative provider/runtime dependency is admitted. Decisions intentionally deferred by the approved spec are represented as explicit later authority gates with a closed set of candidate dispositions, not placeholders.

### Type/name consistency

Use the semantic owner name **PersonalNotifications** in Product wire metadata and narrative label **Personal Notifications** in documents/UI language. Use `notification_id`, `organization_id`, `recipient_principal_id`, `source_work_ref`, `source_occurrence_discriminator`, `read_at`, `archived_at`, and `revision` consistently unless an accepted D5 schema grammar requires a purely mechanical naming adjustment documented in D5-R3.
