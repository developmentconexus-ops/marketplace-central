# NOTIF-01 D3-R — Communication & Propagation Contract

> **Status:** OPEN — D3-R CANDIDATE / OPERATOR ADJUDICATION REQUIRED
> **Accepted parents:** [D3 Communication / Events](D3-COMMUNICATION-EVENTS.md) + [D2-R Ratification](D6-R2-NOTIF-01-D2-R-RATIFICATION.md) + [D2-R2 Ratification](D6-R2-NOTIF-01-D2-R2-RATIFICATION.md)
> **Semantic families:** [Notification Family Semantic Contracts](D6-R2-NOTIF-01-NOTIFICATION-FAMILY-SEMANTIC-CONTRACTS.md) — OPERATOR-APPROVED
> **Scope:** communication semantics only; no Product API, PostgreSQL/River topology, broker, SSE or frontend implementation
> **Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Purpose

D3-R answers one bounded question:

> How do the ten accepted source owners cause the fourteen accepted personal-awareness families to converge correctly under delay, loss, duplication, out-of-order delivery and replay without moving source authority into Personal Notifications?

D3-R does not change D0/D1 family meaning, D2 identity, routing ownership or Notification lifecycle. It selects only semantic communication forms and failure/recovery contracts.

---

## 2. Governing result — all fourteen source reactions are committed-fact `E`

Every admitted source-owner → Personal Notifications edge is a committed-fact **E** boundary.

Why:

```text
source owner commits meaning it owns
→ that fact is valid without Personal Notifications
→ Personal Notifications reacts independently by materializing personal awareness
```

The source transition never waits for Notification success and never becomes invalid because awareness is delayed.

There is no baseline source-domain `C CreateNotification`, no producer write into Personal Notifications private state and no generic `AnyDomainNotificationCandidate` business contract.

Personal Notifications may use **Q** after `E` when current owner truth is required to prevent stale awareness. `P` remains a later read/UX composition option; the canonical Inbox is Personal Notifications state, not a projection over source domains.

### 2.1 Source-owned semantic `E` contracts

The semantic contracts are source-owned facts. Names below are authority labels, not a physical message-schema/queue naming mandate:

| Source owner | Source-owned committed-fact `E` | Personal Notifications mapping |
| --- | --- | --- |
| Marketplace Portfolio | `MarketplaceInstallationAttentionOccurred` | `MARKETPLACE_INSTALLATION_ATTENTION` |
| Marketplace Offering Operations | `OfferingAsyncActionResultCommitted` | `OFFERING_ASYNC_ACTION_RESULT` |
| Availability Control | `AvailabilityAttentionOccurred` | `AVAILABILITY_ATTENTION` |
| Commercial Economics | `EconomicReconciliationAttentionOccurred` | `ECONOMIC_RECONCILIATION_ATTENTION` |
| Marketplace Sales | `MarketplaceSaleFirstConfirmed` | `NEW_MARKETPLACE_SALE` |
| Marketplace Sales | `SaleAttentionOccurred` | `SALE_ATTENTION` |
| Business-System Materialization | `MaterializationAttentionOccurred` | `MATERIALIZATION_ATTENTION` |
| Fulfillment Lifecycle | `FulfillmentBecameActionable` | `FULFILLMENT_ACTIONABLE` |
| Fulfillment Lifecycle | `FulfillmentAttentionOccurred` | `FULFILLMENT_ATTENTION` |
| Fulfillment Lifecycle | `ShipmentExceptionOccurred` | `SHIPMENT_EXCEPTION` |
| Post-Sale Resolution | `PostSaleAttentionOccurred` | `POST_SALE_ATTENTION` |
| Operational Work | `WorkAssignedToHuman` | `WORK_ASSIGNMENT` |
| Controlled Action Governance | `AuthorizationActionRequiredForHuman` | `AUTHORIZATION_ACTION_REQUIRED` |
| Controlled Action Governance | `AuthorizationDecisionResultCommitted` | `AUTHORIZATION_DECISION_RESULT` |

`NotificationKind` remains Personal Notifications vocabulary. Producers own the source-fact contracts; Personal Notifications owns the fixed mapping from each admitted source contract to its approved `NotificationKind`.

No provider webhook/topic is one of these `E` contracts. Provider evidence must first become committed owner meaning.

---

## 3. Minimum immutable occurrence contract

Every admitted `E` preserves only the stable immutable facts Personal Notifications needs for correctness:

```text
organization_id
source_ref
source_occurrence_key
source_occurred_at
source_committed_at
```

Audience fields depend on the accepted strategy:

```text
DIRECT_SOURCE
  + recipient_principal_id

OWNER_DERIVED
  + recipient_principal_id
  // one semantic eligibility occurrence per exact human

ORG_ROUTED
  + no recipient field
  // Personal Notifications resolves historical route lineage
```

Only preferred-awareness contracts may additionally carry the bounded D2-R2 reverse replacement basis:

```text
WorkAssignedToHuman
  + work_replacement_basis?

PostSaleAttentionOccurred
  + sale_attention_replacement_basis?
```

No common free-form payload, arbitrary metadata, provider DTO, source mutable snapshot, generic correlation list or universal EventID is admitted.

### 3.1 One owner-local occurrence key, not one global event identity

`source_occurrence_key` answers only:

> Is this the same source-owned material occurrence delivered again, or a distinct occurrence?

It is scoped by explicit Organization + producer/source reference. It does not establish total order and is not a universal `EventID`/SagaID.

---

## 4. Awareness currentness — `E` is not automatically permission to show stale UI

D3's accepted progression law applies to personal awareness too: a late committed fact may wake the consumer, while current source truth is revalidated when currentness is material to the human job.

D3-R therefore classifies the fourteen families into two bounded profiles.

### 4.1 Immutable-occurrence awareness — materialize from `E` without source-currentness `Q`

These occurrences remain meaningful even if source state later progresses:

```text
OFFERING_ASYNC_ACTION_RESULT
NEW_MARKETPLACE_SALE
AUTHORIZATION_DECISION_RESULT
```

The underlying result/first-confirmation/decision occurrence is historical durable meaning. Personal Notifications still performs current recipient identity/access eligibility checks, but does not suppress the item merely because the source later changed.

### 4.2 Current-attention awareness — `E → Q` before a new current Notification is materialized

These families describe an attention condition or responsibility that must still be relevant when recovered/delivered late:

```text
MARKETPLACE_INSTALLATION_ATTENTION
AVAILABILITY_ATTENTION
ECONOMIC_RECONCILIATION_ATTENTION
SALE_ATTENTION
MATERIALIZATION_ATTENTION
FULFILLMENT_ACTIONABLE
FULFILLMENT_ATTENTION
SHIPMENT_EXCEPTION
POST_SALE_ATTENTION
WORK_ASSIGNMENT
AUTHORIZATION_ACTION_REQUIRED
```

For these, Personal Notifications uses the source owner's public semantic boundary to revalidate the **exact source occurrence**, not merely the subject's latest generic status.

The owner-specific Q must distinguish semantically:

```text
STILL_RELEVANT
  exact occurrence remains a current personal-attention reason

NO_LONGER_RELEVANT
  source progressed/resolved/superseded such that showing this awareness now is stale

UNKNOWN_OR_UNAVAILABLE
  owner cannot currently prove either result
```

Consequences:

- `STILL_RELEVANT` → recipient resolution/materialization may continue;
- `NO_LONGER_RELEVANT` → reaction is explicitly reconciled as no current Notification; it is not a silent propagation loss;
- `UNKNOWN_OR_UNAVAILABLE` → never becomes false/empty/resolved; recovery remains pending/retryable.

No universal cross-domain `GetAttentionStatus` API is created. Each producer owns the smallest Q semantics needed for its own occurrence class.

### 4.3 Examples

```text
FulfillmentBecameActionable O
→ delivered after separation already started
→ owner Q proves O is no longer the current actionable-start handoff
→ no stale “ready to start” Notification
```

```text
WorkAssignedToHuman O1 (A assigned)
→ Work later B, then A again under O3
→ delayed O1 arrives
→ Work Q proves current assignment occurrence is O3, not O1
→ O1 does not materialize as current awareness
```

```text
AuthorizationActionRequiredForHuman O
→ authority/delegation changed before delivery
→ Governance Q proves that exact Principal eligibility episode is no longer current
→ no stale approval request Notification
```

---

## 5. Audience resolution contract

Recipient selection is evaluated only after any required currentness Q succeeds.

### 5.1 DIRECT_SOURCE

The producer includes the exact historical human Principal already owned by source lineage:

```text
OFFERING_ASYNC_ACTION_RESULT
WORK_ASSIGNMENT
AUTHORIZATION_DECISION_RESULT
```

Personal Notifications does not replace the recipient with route configuration.

It Qs the D2 identity/access substrate for current Organization membership/Product-access eligibility before materialization. Definitive ineligibility means no current Inbox item; identity/access unavailability remains retryable/unknown.

`WORK_ASSIGNMENT` additionally requires the source-currentness Q from §4.2 so a delayed old assignment episode does not notify merely because the same Principal happens to be assigned again later.

### 5.2 OWNER_DERIVED

`AUTHORIZATION_ACTION_REQUIRED` is emitted as one exact Principal eligibility occurrence per human resolved by Governance authority semantics.

The event carries that Principal. Personal Notifications never recomputes approvers from Permission/role.

Before materialization:

1. Governance Q confirms that exact decision-eligibility occurrence is still current;
2. identity/access Q confirms current Organization/Product eligibility.

Either definitive failure prevents stale personal awareness; unavailable/unknown remains retryable.

### 5.3 ORG_ROUTED

The producer carries no recipients and never queries Personal Notifications routing during its business commit.

Personal Notifications resolves:

```text
route revision R
= revision current at source-owner commit cutover
```

using accepted D2-R2 route history and `source_committed_at` semantics.

Rules:

- later route edits never retarget an older delayed occurrence;
- no route configured at source commit → zero personal recipients for that occurrence; later configuration does not backfill it;
- each historical recipient binding carries the accepted eligibility-continuity epoch;
- current identity/access Q must confirm same Organization, human Principal, current eligibility/Membership and the same bound epoch;
- epoch mismatch/ineligibility fails closed for that recipient and later re-enable does not resurrect the old binding;
- identity/access unavailable/unknown is retryable, not interpreted as no recipient.

Permission may later be an additional D5 eligibility check; it never selects the route.

---

## 6. Recoverable propagation

For an admitted source occurrence whose semantic reaction remains required by §§4–5:

> loss is detectable and recoverable; producer commit is never rolled back because Notification propagation failed.

D3-R makes no exactly-once claim and selects no queue/outbox/broker.

Semantic recovery must be able to reach one of these terminal conclusions per candidate recipient:

```text
MATERIALIZED
  one semantic Notification exists

SUPERSEDED_OR_REPLACED
  accepted preferred-awareness rule already represents the same personal job

NO_LONGER_RELEVANT
  source owner proves current-attention occurrence is stale

NO_RECIPIENT_AT_COMMIT
  ORG_ROUTED had no route at source commit cutover

RECIPIENT_INELIGIBLE
  exact recipient/binding is definitively not eligible under accepted D2 rules
```

`UNKNOWN_OR_UNAVAILABLE` is not terminal and remains recoverable.

D7 may realize detection/recovery with River, owner-local durable reaction records, reconciliation scans or another smaller mechanism. The transport log never becomes source business authority.

### 6.1 Immutable-occurrence recovery

For `OFFERING_ASYNC_ACTION_RESULT`, `NEW_MARKETPLACE_SALE` and `AUTHORIZATION_DECISION_RESULT`, sufficient durable owner state/history must allow the occurrence to be recovered without inventing a universal event store.

### 6.2 Current-attention recovery

For current-attention families, recovery may reconcile the committed occurrence against current producer meaning. If the producer proves it is no longer relevant, creating a stale Notification is not required.

This is explicit reconciliation, not silent event loss.

---

## 7. Duplicate, replay and anti-regression

Delivery/recovery may happen multiple times.

Personal Notifications semantic uniqueness remains the accepted D2 key:

```text
Organization
+ recipient Principal
+ NotificationKind
+ typed source_ref
+ source_occurrence_key
```

Therefore duplicate/replay of one source occurrence for one recipient converges on one Notification or one terminal suppressed/reconciled outcome.

Replay never:

- creates new Notification history for the same occurrence;
- re-runs current policy and rewrites the historical source occurrence;
- changes `source_occurred_at`/`source_committed_at` to delivery time;
- rebinds an ORG_ROUTED occurrence to a newer route revision;
- resurrects a superseded generic Notification;
- repeats an external effect.

No global queue order or global sequence is business truth.

---

## 8. Bounded replacement/supersession under arbitrary order

D2-R2 is now the authority for the only two replacement families. D3-R defines how propagation converges.

### 8.1 Generic-first

```text
generic source Notification G already materialized for P
→ preferred occurrence arrives
→ preferred passes its own currentness/recipient checks
→ preferred Notification F materializes idempotently
→ Personal Notifications supersedes G by F using the closed reason
```

Read/archive state on G is not rewritten to fake replacement. Source truth is untouched.

### 8.2 Preferred-first

```text
preferred Notification F exists with exact reverse replacement basis
→ delayed generic source occurrence arrives
→ Personal Notifications proves exact kind/ref/occurrence match
→ generic candidate does not become a current awareness item
```

If a historical G had previously existed it remains superseded; if G had never been materialized, D3 does not require minting a useless visible record solely for symmetry.

### 8.3 Closed rules only

Allowed replacement edges remain exactly:

```text
source-family alert → WORK_ASSIGNMENT
SALE_ATTENTION → POST_SALE_ATTENTION
```

They may chain when independently proven:

```text
SALE_ATTENTION
→ POST_SALE_ATTENTION
→ WORK_ASSIGNMENT
```

No arbitrary priority/precedence table, generic `replaces_kind`, related-entity graph or user-authored dedup rule is admitted.

A preferred item that fails its own currentness/eligibility checks does not supersede a still-valid generic item merely because a correlation exists.

---

## 9. Route and identity changes are not Notification events by default

Routing edits are Personal Notifications-owned configuration state, not source-domain events.

Identity/access revocation correctness remains Q-authoritative. A Membership/eligibility change event may later optimize cleanup, but missed cleanup cannot be correctness-critical because historical route bindings carry eligibility epochs and materialization Q fails closed.

No `RoleChanged`/`PermissionChanged` event becomes routing authority.

---

## 10. No realtime correctness dependency

D3-R does not require browser push/SSE/WebSocket.

Correctness is:

```text
source owner committed occurrence
→ recoverable Personal Notifications reaction
→ durable Notification state
→ Product read later can recover the truth
```

A later D6/D7 same-origin SSE or other wake-up may invalidate/refetch browser queries, but loss of that wake-up must never lose a Notification or define unread truth.

---

## 11. D3-R negative controls

D3-R fails review if it:

1. uses `C CreateNotification` or producer writes instead of committed-fact `E`;
2. introduces one generic `AnyDomainNotificationCandidate` business event;
3. maps provider webhooks/topics directly to Product Notifications;
4. carries mutable whole-source payloads merely to avoid Q;
5. creates universal EventID/SagaID/global ordering;
6. materializes current-attention Notifications from delayed `E` without owner-currentness revalidation;
7. treats Q unavailable/unknown as false/empty/no-longer-relevant;
8. resolves ORG_ROUTED against current route rather than source-commit route history;
9. uses Permission/role/all-admin as recipient selection;
10. lets revoke→re-enable reactivate an old route binding;
11. depends on queue ordering or waiting windows for the two replacement rules;
12. turns supersession into archive/read/delete/source mutation;
13. introduces a generic cross-kind precedence/dedup engine;
14. makes source business commit depend on Notification success;
15. claims exactly-once delivery;
16. selects River/PostgreSQL schema/outbox/broker/SSE/API/frontend mechanics inside D3.

---

## 12. Coherence result — CANDIDATE

The bounded communication model is:

```text
14 source-owned committed-fact E contracts
        ↓
3 immutable-occurrence awareness families
11 current-attention E→Q families
        ↓
audience resolution
  DIRECT_SOURCE
  OWNER_DERIVED
  ORG_ROUTED-at-source-commit-route-revision
        ↓
current D2 identity/access Q
        ↓
semantic idempotency + recoverable propagation
        ↓
2 bounded reorder-safe replacement rules
        ↓
Personal Notifications durable awareness state
```

The model preserves source authority, requires no global event platform, creates no source-to-Notification transaction dependency, and remains realizable by later D7 mechanisms without selecting them here.

---

## 13. Gate

```text
D0-R Product scope                    ACCEPTED
D1-R boundary authority               ACCEPTED
D2-R identity/data ownership          ACCEPTED
D2-R2 temporal repair                 ACCEPTED / OPERATOR-RATIFIED
D3-R communication/propagation        READY FOR OPERATOR REVIEW
D5 Product/OAD                        BLOCKED
D6 bell/Inbox/settings                BLOCKED
D7 runtime                            BLOCKED
D8 proof                              BLOCKED
Product implementation                BLOCKED UNTIL D9
```

**Exact next action:** operator adjudicates only this D3-R communication/propagation candidate. If approved, close D3-R and open only D5-R3 Product/OAD authority. Do not edit the OAD, alter B00, select River/runtime mechanics or resume B10/B20 before explicit D3-R ratification.