# NOTIF-01 D2-R2 — Temporal Routing & Late Supersession Ratification

> **Status:** ACCEPTED / OPERATOR-RATIFIED
> **Accepted artifact:** [D2-R2 Temporal Routing & Late Supersession Repair](D6-R2-NOTIF-01-D2-R2-TEMPORAL-ROUTING-SUPERSESSION.md), blob `5e6f9c936d3b4df9de5ed8ad0d88fbf68eab316d`
> **Triggering proof:** [D3-F1 Temporal Routing & Reordered Suppression Finding](D6-R2-NOTIF-01-D3-F1-TEMPORAL-ORDERING-FINDING.md)
> **Accepted base:** D2-R remains ACCEPTED / OPERATOR-RATIFIED; unaffected D2-R clauses remain unchanged
> **Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## Ratified repair

The operator ratified only the bounded temporal-lineage repair proven necessary by D3-F1:

1. `source_committed_at` is distinct from source occurrence time and Notification creation time;
2. ORG_ROUTED state retains versioned/effective history under natural key `(Organization, NotificationKind)` so delayed/replayed occurrences resolve against the route current at source-owner commit cutover;
3. ORG_ROUTED recipient bindings include an opaque Organization+Principal eligibility-continuity epoch so revoke → re-enable cannot silently resurrect prior routing responsibility;
4. already-materialized generic awareness may be superseded only by the two approved bounded precedence rules;
5. `WORK_ASSIGNMENT` and `POST_SALE_ATTENTION` may preserve only their approved typed reverse replacement basis so preferred-first and generic-first arrival converge identically.

The repair does **not** admit global event ordering/EventID, configured-empty routing, Permission-derived recipients, generic Notification status/precedence/relation graphs, routing/subscription DSL, broker/runtime mechanics, OAD changes or frontend implementation.

## Gate

```text
D0-R                           ACCEPTED
D1-R                           ACCEPTED
D2-R                           ACCEPTED
D2-R2                          ACCEPTED / OPERATOR-RATIFIED
D3 communication/propagation   OPEN / NEXT
D5 / D6 / D7 / D8             BLOCKED for NOTIF-01
Product implementation         BLOCKED UNTIL D9
```

**Exact next action:** resume only the bounded NOTIF-01 D3 communication/propagation gate against the accepted D2-R + D2-R2 model.