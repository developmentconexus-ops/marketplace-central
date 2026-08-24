# NOTIF-01 D2-R3 — Explicit Route Unconfigure Cutover

> **Status:** OPEN — TARGETED D2 REOPEN / OPERATOR ADJUDICATION REQUIRED
> **Trigger:** [D5-F1 Operation Surface & Route Reversibility Finding](D6-R2-NOTIF-01-D5-F1-OPERATION-SURFACE-ROUTE-REVERSIBILITY.md)
> **Accepted base:** D2-R + D2-R2 remain ACCEPTED / OPERATOR-RATIFIED; all unaffected clauses remain unchanged
> **Scope law:** repair only the inability to intentionally return one ORG_ROUTED kind from configured routing to unconfigured routing while preserving temporal replay/recovery correctness
> **Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Why D2 reopens narrowly

Accepted D2-R2 deliberately admitted:

```text
before first route revision
→ UNCONFIGURED

configured route revision
→ one-or-more exact human recipient bindings
```

and deferred an explicit later route-removal/disable meaning until a real Product consumer proved it.

D5-R3 now proves that consumer: Organization notification routing is a mutable Settings capability. A configuration that can be created but never intentionally removed is not a complete Product configuration lifecycle.

The repair does not create `configured empty`, per-user opt-out, subscriptions or a notification preference engine.

---

## 2. Route semantic identity remains unchanged — CANDIDATE

The natural route key remains:

```text
NotificationRouteKey = (
  organization_id,
  notification_kind
)
```

No RoutingConfigID is introduced.

Route history remains immutable owner-local revision lineage under that key.

---

## 3. Route revision state becomes an explicit closed union — CANDIDATE

Each committed route revision is exactly one of:

```text
CONFIGURED
  recipient_bindings = one-or-more

UNCONFIGURED
  recipient_bindings = none / not applicable
```

This is a route-state distinction, not `configured with zero recipients`.

Semantic meaning:

- `CONFIGURED` — future source occurrences whose source-owner commit cutover falls under this revision resolve against its explicit recipient bindings;
- `UNCONFIGURED` — future source occurrences whose source-owner commit cutover falls under this revision have no ORG_ROUTED personal recipients for this kind.

The initial pre-history state before any route revision is semantically `UNCONFIGURED` as already accepted. D2-R3 merely allows an explicit later committed revision to restore that state.

---

## 4. Temporal cutover law — CANDIDATE

D2-R2 route-time selection remains unchanged:

```text
source occurrence O
→ choose route revision current at O.source_committed_at
```

Therefore:

```text
t1 R1 = CONFIGURED [A]
t2 occurrence O1 commits
t3 R2 = UNCONFIGURED
t4 delayed O1 is processed
→ O1 still resolves under R1 → A
```

while:

```text
t1 R1 = CONFIGURED [A]
t2 R2 = UNCONFIGURED
t3 occurrence O2 commits
t4 O2 is processed
→ O2 resolves under R2 → no personal recipients
```

and:

```text
t1 R2 = UNCONFIGURED
t2 occurrence O2 commits
t3 R3 = CONFIGURED [B]
t4 delayed/recovered O2 is processed
→ no backfill; O2 remains no-recipient-at-commit
```

A later reconfiguration never retroactively changes an occurrence committed under an unconfigured revision.

---

## 5. Reconfiguration after unconfigure — CANDIDATE

A later explicit routing decision may create a new `CONFIGURED` revision:

```text
R1 CONFIGURED [A@epoch-1]
R2 UNCONFIGURED
R3 CONFIGURED [B@epoch-7]
```

R3 captures current recipient eligibility-continuity epochs exactly as D2-R2 already requires.

No old binding is resurrected by convenience. Reconfiguration is a new explicit routing decision.

---

## 6. Historical retention and current state — CANDIDATE

Current route state is the latest committed route revision under the natural key.

Prior revisions remain durable only to the extent required for:

- delayed occurrence materialization;
- recovery/replay;
- explaining which routing decision governed an already-created Notification when materially required.

D2-R3 does not select physical retention windows/table layout. D7 must preserve enough history to satisfy accepted recoverability.

An explicit `UNCONFIGURED` revision is not treated as deletion of earlier route history.

---

## 7. Non-goals / negative controls — CANDIDATE

This repair does **not** admit:

- `CONFIGURED` with an empty recipient set;
- per-user mute/opt-out preferences;
- user-authored subscriptions;
- scheduling or quiet hours;
- email/push/digest channel preferences;
- role/group/Permission-derived routing;
- default-admin fallback;
- custom NotificationKinds;
- arbitrary routing predicates/DSL;
- deleting temporal route history;
- rewriting historical Notification recipients after unconfigure;
- backfilling occurrences committed during an unconfigured interval;
- D5 HTTP spelling, D7 storage or D6 settings UI implementation.

---

## 8. D5 feed-forward consequence — CANDIDATE

If this repair is accepted, D5-R3 may represent one route slot as a desired state with exactly two semantic outcomes:

```text
configured(one-or-more exact human Principals)
unconfigured
```

A single owner-local Product operation may update between those states using stale-write protection. Exact HTTP method/body/precondition encoding remains D5 wire work after the operation table is ratified.

---

## 9. Gate

```text
D2-R / D2-R2                     ACCEPTED / OPERATOR-RATIFIED
D3-R                              ACCEPTED / OPERATOR-RATIFIED
D5-F1 reversibility finding       PROVED CONSUMER
D2-R3 route-unconfigure cutover   READY FOR OPERATOR REVIEW
D5-R3 operation-table ratification BLOCKED BY D2-R3
canonical Product OAD             UNCHANGED
D6 / D7 / D8                      BLOCKED for NOTIF-01
Product implementation            BLOCKED UNTIL D9
```

**Exact next action:** operator adjudicates only this bounded D2-R3 repair. If approved, return immediately to the already-derived four-operation D5-R3 table and adjudicate it before any OpenAPI modification.