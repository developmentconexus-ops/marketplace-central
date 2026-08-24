# NOTIF-01 D3-F1 — Temporal Routing & Reordered Suppression Finding

> **Status:** D3 ANALYSIS FINDING / TARGETED D2 REOPEN REQUIRED
> **Accepted inputs:** [D2-R Ratification](D6-R2-NOTIF-01-D2-R-RATIFICATION.md) + [D2-R Identity & Data Ownership](D6-R2-NOTIF-01-D2-R-IDENTITY-DATA-OWNERSHIP.md) + accepted [D3 Communication / Events](D3-COMMUNICATION-EVENTS.md)
> **Scope:** prove Notification communication correctness under delayed, duplicate, out-of-order and replayed propagation before selecting D7 mechanics
> **Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Why D3 stops

Accepted D3 requires consequential `E` propagation to remain correct under loss, delay, duplication, out-of-order delivery and replay. Arrival order cannot become business authority. D2-R intentionally left late suppression ordering for D3 proof.

Three concrete counterexamples show that the current D2-R state is insufficient even though its semantic direction remains correct. D3 therefore stops and reopens only the responsible D2 clauses. No D1 meaning or Product scope changes.

## 2. Falsifier A — delayed ORG_ROUTED occurrence crosses a route change

Accepted D2-R says routing changes affect future source occurrences only.

Counterexample:

```text
t1 route NEW_MARKETPLACE_SALE = [A]
t2 Sales commits occurrence O
t3 route changes to [B]
t4 delayed O reaches Personal Notifications
```

If Personal Notifications resolves against **current** routing at t4, B receives an occurrence that happened while A was the configured audience and A silently loses it. If it uses no durable route history, replay cannot reconstruct the correct recipient set.

Therefore current mutable `(Organization, NotificationKind) -> recipients` alone is insufficient. Routing requires versioned/effective history and the source occurrence requires a durable producer-commit cutover anchor distinct from delivery time.

## 3. Falsifier B — preferred awareness arrives after/before the generic item

Two accepted precedence rules must survive arbitrary delivery order:

```text
source alert -> WORK_ASSIGNMENT
SALE_ATTENTION -> POST_SALE_ATTENTION
```

### Generic first

```text
t1 source alert materializes for P
t2 correlated Work/Post-Sale awareness arrives
```

Without a bounded owner-local supersession state, the already-visible generic Notification cannot converge to the approved single preferred awareness item without abusing archive/delete/read state.

### Preferred first

```text
t1 WORK_ASSIGNMENT or POST_SALE_ATTENTION materializes for P
t2 delayed generic source occurrence arrives
```

Without durable reverse correlation on the preferred awareness item, Personal Notifications cannot prove that the late generic item is the exact replaced candidate. Waiting an arbitrary time window is not correctness and queue ordering is not authority.

Therefore D2 needs only the specific late-supersession state/correlation required by the two already-approved precedence rules; a generic cross-kind relation engine remains forbidden.

## 4. Falsifier C — revoked routing recipient silently resurrects

Accepted D2-R says access/Membership revocation must not leave latent routing responsibility that silently resumes after re-eligibility.

Counterexample:

```text
t1 route K includes Principal P
t2 P loses Product eligibility or Organization Membership
t3 no K occurrence happens while P is ineligible
t4 P later becomes eligible/member again
t5 K occurrence happens
```

If the route stores only `principal_id` and resolution checks only current eligibility, P silently resumes receiving K at t5 even though no explicit routing decision re-added P.

A best-effort cleanup event is insufficient because missed cleanup would reintroduce the bug. The route recipient binding therefore needs a stable identity/access **eligibility-continuity anchor** captured at configuration time and checked by Q at resolution time. A break/re-establishment of eligibility yields a new anchor; the old route binding fails closed until an explicit routing edit binds the new episode.

This remains identity/access meaning, not Permission-derived routing.

## 5. Communication direction that remains valid

The D1/D2 source semantics are not falsified:

- all fourteen source-owner -> Personal Notifications reactions remain legitimate committed-fact **E** candidates;
- DIRECT_SOURCE carries exact source-owned human recipient;
- OWNER_DERIVED carries exact producer-resolved human recipients;
- ORG_ROUTED leaves recipient selection to Personal Notifications;
- current identity/access eligibility is queried through the accepted D2 substrate **Q** boundary;
- Notification creation/read/archive never mutates source truth.

D3 final classification is blocked only until the three state-lineage gaps above are repaired.

## 6. Required targeted D2 repair

The smallest repair must add only:

1. versioned/effective ORG_ROUTED history, without a synthetic routing business identity;
2. a producer commit/cutover time or equivalent durable anchor distinct from source business time and delivery time;
3. recipient eligibility-continuity binding for ORG_ROUTED entries;
4. bounded late-supersession state for an already-materialized replaced Notification;
5. bounded reverse replacement correlation on `WORK_ASSIGNMENT` and `POST_SALE_ATTENTION` so preferred-first delivery is also deterministic.

It must not add:

```text
global queue ordering
universal EventID
generic relation/causal graph
generic notification status machine
routing DSL
generic cross-kind dedup engine
broker/outbox/runtime selection
```

## 7. Gate

```text
D2-R base identity/data ownership       ACCEPTED / OPERATOR-RATIFIED
D3-F1 communication analysis            FALSIFIER FOUND
D2-R2 targeted temporal repair          REQUIRED / NEXT
D3 final communication gate             BLOCKED BY D2-R2
D5 / D6 / D7 / D8                       BLOCKED for NOTIF-01
Product implementation                   BLOCKED UNTIL D9
```

**Exact next action:** adjudicate the bounded D2-R2 temporal-routing / late-supersession repair derived from this finding. Do not choose D3 event names or D7 transport mechanics first.
