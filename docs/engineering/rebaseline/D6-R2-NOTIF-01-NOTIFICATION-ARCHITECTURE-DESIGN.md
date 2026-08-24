# D6-R2 NOTIF-01 — Personal Notification Architecture Design

> **Status:** OPERATOR-APPROVED DESIGN / AUTHORITY-AMENDMENT PLAN NEXT
> **Verification:** repository full gate passed on the written design candidate before operator spec adjudication
> **Trigger:** operator-requested Notification assessment during D6-R2 frontend realization
> **Comparative evidence:** MetalDocs Notification direction; not Marketplace Central authority
> **Current Product authority:** D0–D8 + D5-R2/D8-R2 remain accepted until targeted amendments are separately approved
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. Problem

Marketplace Central already distinguishes source truth, Operational Work, Governance approvals and read-only attention, but has no Product-owned personal Notification/Inbox state. D0 explicitly distinguishes escalation from notification; Work therefore cannot be repurposed as a personal awareness inbox without distorting its accepted responsibility.

The proven human need is narrower than workflow:

> a human Principal must be able to discover personally relevant committed MPC facts, retain unread/read/archive state, and navigate back to the current source authority without Notification becoming source truth, authorization, Work, Audit or acknowledgement of the originating fact.

This design treats Notification as a new Product responsibility candidate and intentionally pauses dependent frontend progression until the upstream authority is resolved.

## 2. Selected architecture

Admit **Personal Notifications** as a small supporting semantic owner inside the existing modular monolith.

It owns only personal Inbox state for an Organization-scoped human recipient.

It does **not** own:

- the originating business fact;
- Work responsibility, assignment, escalation or resolution;
- Governance authorization;
- Audit/history truth;
- source-object access;
- external e-mail/push delivery;
- generic subscriptions/rules;
- event-bus infrastructure;
- a cross-owner workflow lifecycle.

Notification is durable Product state, not an ephemeral frontend projection.

## 3. Minimal Product state

A Notification candidate carries only the smallest semantics needed for a personal Inbox:

```text
Notification
  notification_id          MPC-owned opaque identity
  organization_id          canonical Organization scope
  recipient_principal_id   exact human Principal recipient
  kind                     closed Product-defined notification kind
  source_ref               bounded typed reference to the originating owner subject/occurrence
  created_at
  read_at?                  null = unread; timestamp = read
  archived_at?              null = active Inbox; timestamp = archived
```

No generic `{entity_type, entity_id}` platform graph is admitted. `source_ref` must be a closed typed union derived only from proved trigger consumers.

No `seen` state is admitted now. `seen` may reopen only if a distinct Product need is proven between “rendered somewhere” and “read/opened”.

No Notification content copy may become a substitute for current source authority. A bounded title/summary may exist for Inbox usability, but sensitive/current detail must be re-read from the source owner when opened.

## 4. Initial trigger law

Do not create event-per-CRUD notifications.

The first mandatory trigger is the strongest currently evidenced personal-awareness occurrence:

```text
Operational Work becomes explicitly assigned or reassigned to a human Principal
  -> Personal Notification for that exact new Principal
```

The Notification says that Work became personally relevant; Work remains the obligation authority.

Reading or archiving the Notification does not alter Work.

Further triggers such as delayed Authorization Decisions, asynchronous action outcomes, material ambiguity or other personal outcomes require a **Notification Trigger Census** and separate proof that:

1. the producer owns a committed fact;
2. an exact human recipient can be derived without a generic subscription engine;
3. personal awareness is useful and not merely duplicate noise;
4. the source fact is not itself an actionable Work obligation that should remain represented only as Work;
5. duplicates can be suppressed by stable semantic occurrence identity.

## 5. Cross-owner delivery semantics

Marketplace Central must **not** copy MetalDocs' proposed cross-owner shared business transaction if Personal Notifications is a separate semantic owner.

Accepted MPC D7 law remains:

```text
BEGIN source-owner transaction
  commit source-owned fact
  River InsertTx durable Notification reaction
COMMIT

River worker
  enter exact Organization scope
  re-enter Personal Notifications owner
  apply semantic idempotency
  create Notification in Notifications-owned transaction
  complete durable job
```

The correctness guarantee is:

> a committed source occurrence that requires a Notification cannot silently lose its required Notification reaction.

It is **not** an exactly-once execution claim and does not require source-owner + Notification state in one cross-owner transaction.

Duplicate/redelivered River work must not create duplicate Notification business state.

## 6. Access and security

Notification never grants access to its source.

Opening a Notification performs a current authorized Product read of the referenced source object/subject. If current access was revoked or the source is no longer visible, the Notification remains the recipient's Inbox state but the source deep link fails closed under current authority.

Notifications are Organization-scoped and recipient-scoped. Cross-Organization Inbox access is forbidden.

The browser experience is for human Principal `H`. A/S machine Principals do not gain a human Inbox merely because the backend supports Notification reactions.

Exact ordinary Permission wording and operation mapping are deferred to the bounded D5 amendment; they must not be forced into an unrelated existing Permission merely to preserve the current 30-Permission census.

## 7. Product API shape — semantic requirement only

D5 will derive exact operationIds, paths, methods, error semantics and Permission mapping only after the targeted D0/D1/D2/D3 amendments are accepted.

The minimum required capability classes are:

```text
Q  list own Notifications in current Organization
C  mutate own read/unread state
C  archive/unarchive own Notification when admitted
```

Creation is an internal accepted reaction, not a public browser `CreateNotification` command.

No baseline requirement for:

- admin-wide Notification search;
- arbitrary recipient creation;
- delete;
- notification preference center;
- bulk mutation;
- exact unread aggregate count;
- e-mail/push subscriptions;
- generic notification templates.

The current 99-operation / 30-Permission census is **not** protected by preference if Notifications is accepted Product scope. Any census increase must be the smallest derived consequence of the new capability.

## 8. Frontend experience

### 8.1 Topbar bell — operator-approved requirement

The B00 physical shell remains locked except for a **bounded topbar utility-slot reopen** required by this new Product responsibility.

For an authenticated human in a current Organization:

```text
page title / page-local context                 [bell]
                                                   │
                                                   └─ Personal Inbox entry point
```

The bell is Organization-scoped, not cross-Organization aggregation.

Baseline behavior:

- bell is shown only when current access exposes the Notification Inbox capability;
- no unread number is inferred from a paginated collection;
- baseline attention indicator is a simple dot/presence state: “there is at least one unread Notification”;
- clicking the bell opens a bounded recent-Inbox preview;
- preview offers `Ver todas` to an organization-scoped full Inbox route, e.g. `/org/:organizationId/notificacoes`;
- the full Inbox is a utility surface reached from the topbar, not a sixth global sidebar mass;
- Organization switch invalidates/refetches Notification state and closes incompatible preview state;
- marking read/archive updates Notification only;
- source navigation re-checks current source authorization.

A numeric unread badge requires a future authoritative aggregate read if proved useful.

### 8.2 Frontend state law

TanStack Query owns Notification server state. Bell-open/popover state is ephemeral UI state. Notification read/archive state is never mirrored as frontend-owned durable truth.

The Inbox must distinguish at least:

```text
unread
read
archived
knowledge unavailable / request failed
```

An unavailable Inbox read must never be displayed as “no notifications”.

## 9. Realtime posture

Realtime is a delivery optimization, not Notification truth.

Correctness baseline:

```text
PostgreSQL Notification state
  -> Product API read
  -> TanStack Query Inbox
```

Candidate low-cost realtime seam to validate later:

```text
Notification committed
  -> best-effort lightweight wake-up
  -> SSE same-origin
  -> browser invalidates/refetches Notification query
```

PostgreSQL `LISTEN/NOTIFY` is a candidate wake-up mechanism only. It must not be the sole record, history or exact unread count source.

A lost wake-up must not lose a Notification. Browser focus/refetch or another later read recovers the persistent state.

Do not make a mandatory `NOTIFY` call capable of rolling back canonical Notification persistence merely to achieve realtime UX. If LISTEN/NOTIFY is selected, the implementation design must preserve durable Notification commit independently from optional wake-up failure.

No Kafka, RabbitMQ, NATS, Redis queue or second generic outbox is admitted. Existing PostgreSQL + River remains the durable reaction substrate unless a measured later reopen proves it insufficient.

## 10. Error and failure semantics

Required behavior includes:

- source commit succeeds + process dies before consumer execution -> River handoff eventually recreates the required reaction;
- reaction delivery duplicates -> one semantic Notification;
- Notification creation fails transiently -> durable reaction retries safely;
- source access is revoked after Notification creation -> Notification remains, source read fails closed;
- Notification query fails -> UI shows unavailable, not empty;
- realtime signal is lost -> later Product read recovers state;
- mark-read/archive concurrency is handled by the accepted Product precondition/concurrency grammar derived in D5;
- Notification mutation never mutates the source object.

## 11. YAGNI exclusions

Launch-V1 candidate explicitly excludes until separate evidence:

```text
seen state
numeric unread aggregate
mark-all-read
bulk archive
notification preferences
per-kind user subscriptions
digests
email
mobile/web push
generic template engine
generic EventStore/pub-sub platform
external broker
cross-Organization inbox
A/S human-like inbox
generic entity-reference graph
Notification-triggered source mutation
```

## 12. Targeted reopen sequence

With this written design operator-approved, authority changes proceed in dependency order:

```text
D0-R?  add Personal Notification Inbox to Product 1.0 launch scope
  ↓
D1-R?  add supporting Personal Notifications semantic owner + only proved edges
  ↓
D2-R?  Notification identity/ownership/recipient/source-ref/history semantics
  ↓
D3-R?  trigger/event census + recoverable propagation contracts
  ↓
D5-R?  exact Product operations / Permissions / OAD / generated proof
  ↓
D6-R?  bell + Inbox interaction authority; bounded B00 topbar utility reopen
  ↓
D7-R?  realize on existing PostgreSQL + River; validate optional realtime seam
  ↓
D8-R?  add smallest composed Notification falsifier
  ↓
resume D6-R2 dependent frontend planning
```

D4 is not expected to reopen because Personal Notifications has no provider transport responsibility.

No Product code is authorized by this sequence.

## 13. Smallest falsifiable proof contract

At minimum the later composed proof must falsify:

```text
1. Assign/reassign Work to human Principal -> required Notification eventually exists for the new Principal.
2. Duplicate River reaction -> no duplicate semantic Notification.
3. Read Notification -> Work/source state unchanged.
4. Archive Notification -> Work/source state unchanged.
5. Notification -> source access is re-authorized at click time.
6. Cross-Organization Notification read is blocked.
7. Lost realtime wake-up -> persistent Notification is still recovered by read.
8. Inbox read failure -> UI never presents known-empty.
9. Machine Principal does not receive human Inbox behavior.
10. No external broker is required for correctness.
```

## 14. Decision summary

Selected direction:

> **Marketplace Central Personal Notifications is a small supporting semantic owner for Organization-scoped, human-Principal-targeted awareness state. Required Notifications are produced through recoverable owner-local event/reaction semantics over existing PostgreSQL + River. Notification is not Work, Audit, authorization, source truth or acknowledgement. The React shell gains a bounded topbar bell entry point and organization-scoped Inbox; realtime may later use a disposable SSE wake-up seam but persistent Product state remains the sole truth.**
