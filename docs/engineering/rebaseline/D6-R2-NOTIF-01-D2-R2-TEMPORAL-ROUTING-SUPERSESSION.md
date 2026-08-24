# NOTIF-01 D2-R2 — Temporal Routing & Late Supersession Repair

> **Status:** OPEN — TARGETED D2 REOPEN / OPERATOR ADJUDICATION REQUIRED
> **Trigger:** [D3-F1 Temporal Routing & Reordered Suppression Finding](D6-R2-NOTIF-01-D3-F1-TEMPORAL-ORDERING-FINDING.md)
> **Accepted base:** [D2-R Ratification](D6-R2-NOTIF-01-D2-R-RATIFICATION.md) + [D2-R Identity & Data Ownership](D6-R2-NOTIF-01-D2-R-IDENTITY-DATA-OWNERSHIP.md)
> **Scope law:** every unaffected D2-R clause remains ACCEPTED; this document repairs only route-time lineage, eligibility continuity and the two already-approved suppression rules under arbitrary delivery order
> **Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Why the accepted D2-R reopens narrowly

D3 must prove correctness when committed events are delayed, duplicated, reordered or replayed. Three counterexamples cannot be solved by queue-order assumptions without violating accepted D2 meaning:

1. a delayed ORG_ROUTED occurrence can cross a later route edit;
2. a generic Notification can materialize before the preferred Work/Post-Sale awareness, or vice versa;
3. a route entry containing only PrincipalID can silently resume after revoke→re-enable without an explicit routing decision.

The repair below adds only the durable lineage necessary to make those cases deterministic. It does not reopen the fourteen family meanings, producer edges, typed source union, Notification identity, read/archive semantics, or Product scope.

---

## 2. Source commit time becomes explicit — CANDIDATE

D2-R already separates `source_occurred_at` from Notification `created_at`. D3 proves that ORG_ROUTED cutover additionally needs the time at which the **MPC source owner committed the accepted attention occurrence**.

Canonical Notification state therefore gains:

```text
source_committed_at
```

with semantics:

```text
source_occurred_at
= when the source/business occurrence is semantically effective/occurred

source_committed_at
= when the MPC source owner durably committed the admitted attention occurrence

created_at
= when Personal Notifications durably materialized this personal Notification
```

These times may differ and must not be rewritten into one another.

Routing applicability is based on the source-owner **commit cutover**, not provider webhook time and not transport delivery time.

Exact database/clock carrier remains D7. D7 must realize an unambiguous committed-before/after cutover between a route revision and a source occurrence; an implementation that leaves a same-boundary race semantically ambiguous is non-conformant.

---

## 3. ORG_ROUTED history is versioned/effective state — CANDIDATE

The semantic route identity remains natural and unchanged:

```text
NotificationRouteKey = (organization_id, notification_kind)
```

No RoutingConfigID is introduced.

However, replacing prior recipients cannot erase the route version needed to recover an older delayed source occurrence. Each committed route update therefore produces an immutable historical **route revision** under that same semantic key:

```text
organization_id
notification_kind
route_revision
committed_at
recipient_bindings  // one or more
```

`route_revision` is owner-local ordering/concurrency lineage for one `NotificationRouteKey`; it is not a global sequence or business entity identity.

Current route is the latest committed revision. Prior revisions remain sufficient to resolve still-recoverable source occurrences whose `source_committed_at` fell under that revision.

Baseline still has no proved “configured intentionally to nobody” state. Before the first route revision, the kind is `UNCONFIGURED`. Every admitted revision has one-or-more configured recipient bindings.

If a future Product operation needs explicit route removal/disable distinct from absence, D2/D5 reopen then rather than encoding configured-empty now.

### 3.1 Temporal route-selection law — CANDIDATE

For one ORG_ROUTED source occurrence O:

```text
R = route revision that was current when O was committed by its source owner
```

Personal Notifications must resolve O using R even if:

- delivery occurs after later route revisions;
- the event is replayed;
- a recovery sweep discovers O later.

Therefore:

```text
t1 route = [A]
t2 source occurrence O commits
t3 route = [B]
t4 O is delivered/recovered
→ O uses the t2 route lineage, not current [B]
```

If no route had yet been configured at O's commit cutover, later route creation does not backfill O.

This does not require global event ordering. It requires only a deterministic cutover relation between the source occurrence commit and the route history relevant to its `(Organization, NotificationKind)`.

---

## 4. Route recipient binding includes eligibility continuity — CANDIDATE

Storing only `principal_id` cannot distinguish uninterrupted eligibility from revoke→re-enable.

The D2 identity/access substrate therefore supplies one opaque **OrganizationPrincipalEligibilityEpoch** for the current uninterrupted interval in which Principal P is both:

```text
currently Product-access eligible
AND
currently a member of Organization O
```

The epoch changes whenever either continuity is broken and later re-established. It is:

- not Principal identity;
- not Membership identity;
- not an AccessRole or Permission snapshot;
- not an OIDC/session/token value;
- not business responsibility;
- an opaque D2 identity/access continuity fact used only where a consumer must distinguish “same access episode” from “later re-established access”.

An ORG_ROUTED historical recipient binding is therefore:

```text
recipient_principal_id
eligibility_epoch
```

At route-edit time, Personal Notifications captures the current epoch through the D2 identity/access boundary.

At source-occurrence materialization/recovery time, the recipient is eligible only when a current Q proves:

```text
same Organization
+ human Principal
+ currently eligible/member
+ current eligibility_epoch == bound eligibility_epoch
```

Consequences:

```text
route configured for P
→ P revoked
→ P later re-enabled/re-added
→ old route binding does NOT silently reactivate
```

An explicit routing edit is required to bind the new eligibility episode.

Ordinary role/Permission changes do not by themselves create a new eligibility epoch unless they break the accepted Product-access/Membership continuity. Exact D5 Permission compatibility may still add an independent eligibility check later; Permission never selects recipients.

---

## 5. Bounded late-supersession lifecycle — CANDIDATE

The approved suppression rules require convergence even when the generic awareness item already became visible before the preferred item arrives.

A Notification therefore gains the following **optional bounded owner-local awareness state**:

```text
superseded_at?
superseded_by_notification_id?
supersession_reason?
```

All three are absent together or present together.

Allowed reasons are closed:

```text
work_assignment_replacement
post_sale_precedence
```

No other reason exists by symmetry.

Semantic laws:

- supersession is Personal Notifications awareness state only;
- it never changes/deletes source-domain truth;
- it never changes Work/Post-Sale/Governance state;
- it is not user archive/read state;
- `read_at` and `archived_at` remain unchanged/orthogonal;
- a superseded Notification is no longer a current personal-awareness item for active Inbox/unread-presence semantics;
- historical identity/correlation remains stable so replay remains idempotent;
- `superseded_by_notification_id` must reference a Notification in the same Organization and for the same historical recipient;
- a Notification may itself later be superseded only through another explicitly allowed bounded rule; cycles are invalid.

An already-visible item may therefore converge honestly instead of being deleted, auto-archived or silently rewritten.

---

## 6. Preferred-awareness reverse correlation — CANDIDATE

Late supersession must also work when the **preferred** Notification arrives first. Personal Notifications therefore preserves reverse replacement basis only on the two preferred kinds that need it.

### 6.1 `WORK_ASSIGNMENT` replacement basis

When a Work assignment is known to arise from an exact admitted source attention occurrence, the `WORK_ASSIGNMENT` Notification may carry one bounded `work_replacement_basis`:

```text
replaced_notification_kind
replaced_source_ref
replaced_source_occurrence_key
```

The `replaced_notification_kind + replaced_source_ref` pair must satisfy the already-accepted D2-R kind/source compatibility matrix and must represent the exact source occurrence that causally produced the Work obligation/assignment lineage.

It is not inferred from title, timestamp, target label or arbitrary shared entity identity.

If the generic source Notification arrives later, Personal Notifications can prove it is replaced for this recipient without queue-order assumptions.

### 6.2 `POST_SALE_ATTENTION` replacement basis

When a Post-Sale Resolution is the richer continuation of the exact Sale attention consequence, its `POST_SALE_ATTENTION` Notification may carry:

```text
sale_attention_replacement_basis:
  marketplace_sale_ref
  sale_attention_source_occurrence_key
```

This basis is fixed to `SALE_ATTENTION`; it is not a generic cross-kind relation.

If the delayed `SALE_ATTENTION` arrives later for the same recipient, the richer Post-Sale item remains preferred.

### 6.3 Preferred-first vs generic-first convergence — CANDIDATE

**Generic first:**

```text
generic Notification exists
→ preferred Notification arrives with exact replacement basis
→ generic becomes superseded by preferred
```

**Preferred first:**

```text
preferred Notification exists with exact replacement basis
→ generic source occurrence arrives later
→ Personal Notifications detects the exact replacement basis
→ generic candidate is suppressed before becoming a current awareness item
```

If a suppressed candidate had never materialized before the preferred item, D2-R2 does not require minting a useless visible Notification solely for symmetry. The preferred Notification's durable replacement basis remains sufficient for replay/dedup of that late candidate.

If the generic Notification had already materialized, its stable identity remains and it is marked superseded rather than deleted.

---

## 7. Bounded precedence can chain but never generalize — CANDIDATE

The two accepted rules may legitimately compose when each causal relation is independently proven, for example:

```text
SALE_ATTENTION
→ POST_SALE_ATTENTION
→ WORK_ASSIGNMENT
```

This does not create a generic precedence engine. Each edge must independently satisfy one of the two closed rules and exact typed occurrence correlation.

No arbitrary `priority`, `precedence`, `related_entities[]`, `replaces_kind:string`, causal graph or user-authored dedup rule is admitted.

`WORK_ASSIGNMENT` is never superseded merely because a generic source alert exists. `POST_SALE_ATTENTION` supersedes `SALE_ATTENTION` only for the same bounded consequence/recipient. No cycle is valid.

---

## 8. Revised minimal canonical state — CANDIDATE

The accepted D2-R Notification state becomes, with additions marked by meaning rather than physical schema:

```text
notification_id
organization_id
recipient_principal_id
kind
source_ref
source_occurrence_key
source_occurred_at
source_committed_at        // added by D2-R2
created_at
read_at?
archived_at?
superseded_at?            // bounded late-replacement only
superseded_by_notification_id?
supersession_reason?
work_replacement_basis?    // WORK_ASSIGNMENT only
sale_attention_replacement_basis? // POST_SALE_ATTENTION only
revision
```

No free-form payload, metadata, generic status, generic replacement relation or delivery-state machine is introduced.

Route history is separate Personal Notifications canonical configuration/history under `NotificationRouteKey`; it is not embedded into every Notification. Historical Notification recipient remains the final personal identity; D5 may decide whether exposing route revision lineage to clients has any consumer.

---

## 9. D3 proof obligations after this repair — CANDIDATE

If D2-R2 is accepted, D3 must still prove:

1. source-owner -> Personal Notifications propagation form for all fourteen families;
2. minimal immutable occurrence contract;
3. recoverable propagation without exactly-once claims;
4. DIRECT_SOURCE and OWNER_DERIVED recipient handling;
5. ORG_ROUTED historical route-revision selection at source commit cutover;
6. current eligibility-epoch check at materialization/recovery;
7. semantic dedup under replay;
8. generic-first and preferred-first convergence for both bounded suppression rules;
9. no global order dependency;
10. no source business transition depends on Notification success.

D7 then chooses concrete PostgreSQL/River/outbox/wake-up mechanics that realize the accepted D3 semantics.

---

## 10. Negative controls — CANDIDATE

This targeted repair fails if it:

1. changes any of the fourteen Notification family meanings;
2. adds a new producer edge;
3. makes delivery order the routing/suppression authority;
4. resolves delayed ORG_ROUTED occurrences against current route instead of occurrence-time route lineage;
5. uses provider/source event time as the route cutover when MPC owner commit occurred later;
6. stores only PrincipalID and permits revoke→re-enable to silently resume old routing responsibility;
7. uses Permission/role as the eligibility continuity anchor;
8. auto-archives/marks-read/deletes an already-materialized generic Notification to fake supersession;
9. introduces a generic Notification status machine or generic replacement relation;
10. infers replacement from text/timestamp/entity similarity;
11. creates global event ordering/EventID/causal graph;
12. selects broker, River job, DB table/RLS, HTTP endpoint or frontend behavior in D2.

---

## 11. Gate

```text
D2-R base identity/data ownership       ACCEPTED / OPERATOR-RATIFIED
D3-F1 temporal/reorder analysis         TARGETED FALSIFIER
D2-R2 temporal routing/supersession     READY FOR OPERATOR REVIEW
D3 final communication/propagation      BLOCKED BY D2-R2
D5 Product/OAD                          BLOCKED
D6 bell/Inbox/settings                  BLOCKED
D7 runtime                              BLOCKED
D8 proof                                BLOCKED
Product implementation                  BLOCKED UNTIL D9
```

**Exact next action:** operator adjudicates only this D2-R2 repair. If accepted, return directly to D3 and prove the fourteen-family propagation contract. Do not open D5/D6/D7 first.
