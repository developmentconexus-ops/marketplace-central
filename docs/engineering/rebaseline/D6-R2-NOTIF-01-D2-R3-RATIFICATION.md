# NOTIF-01 D2-R3 — Route Unconfigure Ratification

> **Status:** ACCEPTED / OPERATOR-RATIFIED
> **Accepted artifact:** [D2-R3 Explicit Route Unconfigure Cutover](D6-R2-NOTIF-01-D2-R3-ROUTE-UNCONFIGURE-CUTOVER.md), blob `1907b2b7801a06a56a81ccfd1e6b4e611fcccbb7`
> **Accepted base:** D2-R + D2-R2 remain ACCEPTED / OPERATOR-RATIFIED
> **Triggering consumer:** [D5-F1 Operation Surface & Route Reversibility Finding](D6-R2-NOTIF-01-D5-F1-OPERATION-SURFACE-ROUTE-REVERSIBILITY.md)
> **Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## Ratified result

The operator ratified only the bounded route-lifecycle repair required by the proved Settings consumer:

1. `NotificationRouteKey = (Organization, NotificationKind)` remains the natural semantic key; no RoutingConfigID is introduced.
2. Route revisions form a closed state union: `CONFIGURED(one-or-more exact human recipient bindings)` or `UNCONFIGURED`.
3. `UNCONFIGURED` is not `CONFIGURED([])` and is not a personal mute/subscription state.
4. Temporal selection remains source-owner-commit based: delayed/replayed occurrences use the route revision current at `source_committed_at`.
5. An explicit later `UNCONFIGURED` revision affects only future source occurrences under that cutover; older configured occurrences retain their historical route and occurrences committed while unconfigured are never backfilled by later reconfiguration.
6. Reconfiguration after unconfigure is a new explicit routing decision and captures new current eligibility-continuity epochs; old bindings never resurrect.
7. Historical route lineage remains recoverable; D7 chooses physical retention/storage.

No per-user opt-out, quiet hours, email/push preference, role/Permission routing, custom NotificationKind, routing DSL, configured-empty state, history deletion, OAD spelling or runtime/frontend mechanism is admitted.

## Gate

```text
D0-R                           ACCEPTED
D1-R                           ACCEPTED
D2-R                           ACCEPTED
D2-R2                          ACCEPTED
D2-R3                          ACCEPTED / OPERATOR-RATIFIED
D3-R                           ACCEPTED
D5-R3 operation admission      OPEN / NEXT
canonical Product OAD          UNCHANGED
D6 / D7 / D8                   BLOCKED for NOTIF-01
Product implementation         BLOCKED UNTIL D9
```

**Exact next action:** return immediately to the bounded D5-R3 four-operation admission table and adjudicate it before any canonical Product OpenAPI modification.