# NOTIF-01 D1-R — Family Semantic Closure

> **Status:** FAMILY CONTRACTS OPERATOR-APPROVED / D1-R FINAL ADJUDICATION NEXT
> **Approved artifact:** [Notification Family Semantic Contracts](D6-R2-NOTIF-01-NOTIFICATION-FAMILY-SEMANTIC-CONTRACTS.md), blob `13618d9a6d0c8583e0c669b9aee6dd2f84014995`
> **Parent candidate:** [D1-R Producer Edges & Notification Routing Boundary Correction](D6-R2-NOTIF-01-D1-R-PRODUCER-ROUTING-BOUNDARY-CORRECTION.md)
> **Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Ratified semantic contract set

The operator approved all fourteen Launch-V1 Notification Family Semantic Contracts, including each family's:

```text
human meaning
attention value
source owner
birth transition
not-this boundary
audience / consumer
deep-link home
repeat law
suppression / overlap law
conceptual user-language explanation
```

The approved families remain:

```text
MARKETPLACE_INSTALLATION_ATTENTION
OFFERING_ASYNC_ACTION_RESULT
AVAILABILITY_ATTENTION
ECONOMIC_RECONCILIATION_ATTENTION
NEW_MARKETPLACE_SALE
SALE_ATTENTION
MATERIALIZATION_ATTENTION
FULFILLMENT_ACTIONABLE
FULFILLMENT_ATTENTION
SHIPMENT_EXCEPTION
POST_SALE_ATTENTION
WORK_ASSIGNMENT
AUTHORIZATION_ACTION_REQUIRED
AUTHORIZATION_DECISION_RESULT
```

Fourteen is a derived result, not a protected count. Before D2-R, a family must split if a subtype proves materially different source authority, human job, audience strategy, deep-link continuation, repeat law or interaction. Generic `reason:string` or hidden routing logic may not conceal that difference.

## 2. Binding D1-R clarification A — Materialization audience

`MATERIALIZATION_ATTENTION` has one stable `ORG_ROUTED` semantic audience:

```text
Organization
+ MATERIALIZATION_ATTENTION
→ configured business-system-materialization / backoffice operations humans
```

Routing does not vary by hidden materialization subtype.

If downstream physical meaning changes, Fulfillment owns and emits its own admitted awareness occurrence:

```text
FULFILLMENT_ACTIONABLE
or
FULFILLMENT_ATTENTION
```

This prevents a subtype-dependent routing DSL from appearing inside one NotificationKind.

## 3. Binding D1-R clarification B — Sale → Post-Sale awareness precedence

Marketplace Sales and Post-Sale Resolution may both own valid facts for the same external development. Source truth is never suppressed.

For personal awareness only, when all of the following are proven:

```text
same bounded underlying consequence correlation
+ recipient P would receive SALE_ATTENTION
+ recipient P would receive POST_SALE_ATTENTION
+ PostSaleResolution is already the richer owned continuation
```

Personal Notifications may:

```text
suppress SALE_ATTENTION for P
→ deliver POST_SALE_ATTENTION for P
```

The rule is per-recipient and bounded to this proved overlap. It does not create a generic cross-kind deduplication engine and does not change Sale or Post-Sale state.

The previously approved source-alert → `WORK_ASSIGNMENT` replacement remains a separate bounded suppression case.

## 4. Final D1-R candidate after semantic closure

Read the parent D1-R with this closure record. Its final candidate meaning is:

```text
10 explicit source-owner → Personal Notifications edges
→ 14 operator-approved semantic awareness families
→ one Personal Notifications supporting semantic owner
→ DIRECT_SOURCE recipient meaning stays with producer
→ OWNER_DERIVED responsibility/authority stays with producer
→ ORG_ROUTED exact-human configuration belongs to Personal Notifications
→ bounded Work-replacement suppression
→ bounded Sale→Post-Sale precedence
```

Personal Notifications owns awareness vocabulary/lifecycle, bounded A3 routing configuration, historical recipient awareness and the two proved per-recipient suppression semantics. It does not own source business truth, Work assignment, Governance authority, ordinary Permissions/access, provider protocol, event delivery or workflow progression.

## 5. Negative controls

D1-R must still fail if it introduces:

- `AnyDomain → Personal Notifications`;
- provider webhook/topic → Product Notification directly;
- Permission holder = notification recipient by default;
- implicit all-member/admin broadcast;
- subtype-dependent routing hidden inside one kind;
- generic cross-kind dedup/routing/subscription DSL;
- Notification as source truth/workflow authority;
- e-mail/push/digest/broker/runtime mechanics in D1;
- event names, database schema, OAD or frontend implementation before their owning gates.

## 6. Gate

```text
H3 / Trigger + Audience Census       OPERATOR-APPROVED
D0-R Product scope                   ACCEPTED
14 family semantic contracts         OPERATOR-APPROVED
D1-R semantic clarifications A/B     OPERATOR-APPROVED
D1-R final boundary authority        READY FOR OPERATOR REVIEW
D2-R                                 BLOCKED
D3 / D5 / D6 / D7 / D8             BLOCKED for NOTIF-01
Product implementation               BLOCKED UNTIL D9
```

**Exact next action:** operator adjudicates the final D1-R boundary authority as the parent D1-R plus this semantic closure. Do not begin D2-R before explicit D1-R ratification.
