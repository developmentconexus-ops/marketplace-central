# NOTIF-01 D3-R — Communication & Propagation Ratification

> **Status:** ACCEPTED / OPERATOR-RATIFIED
> **Accepted artifact:** [D3-R Communication & Propagation Contract](D6-R2-NOTIF-01-D3-R-COMMUNICATION-PROPAGATION.md), blob `2a7e28edddd551de4093c5e3e69f548a10ca393e`
> **Accepted parents:** D2-R + D2-R2 remain ACCEPTED / OPERATOR-RATIFIED
> **Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## Ratified result

The operator ratified the bounded NOTIF-01 D3 communication contract:

1. all fourteen accepted source-owner reactions are committed-fact `E` boundaries;
2. source business transitions never depend on Personal Notifications success;
3. `OFFERING_ASYNC_ACTION_RESULT`, `NEW_MARKETPLACE_SALE`, and `AUTHORIZATION_DECISION_RESULT` are immutable-occurrence awareness and may materialize from the committed occurrence without source-currentness Q;
4. the other eleven current-attention families require source-owner `Q` currentness reconciliation before a new current personal-awareness item is materialized after delay/recovery;
5. DIRECT_SOURCE and OWNER_DERIVED carry exact producer-owned human recipients; ORG_ROUTED carries no recipient and resolves the historical route revision current at source-owner commit cutover;
6. current identity/access eligibility remains separately authoritative, including D2-R2 eligibility-continuity epochs for ORG_ROUTED recipient bindings;
7. propagation is recoverable, duplicate-safe, replay-safe, and globally unordered; no exactly-once or queue-order correctness assumption is admitted;
8. the two approved precedence rules converge identically under generic-first or preferred-first delivery through the bounded D2-R2 replacement basis/supersession model;
9. realtime browser wake-up is optional D6/D7 convenience and never the persistence/correctness authority.

No provider webhook is a Product event, no generic event bus/business event envelope is introduced, and no D5/D6/D7 mechanism is accepted through this ratification.

## Gate

```text
D0-R                           ACCEPTED
D1-R                           ACCEPTED
D2-R                           ACCEPTED
D2-R2                          ACCEPTED
D3-R                           ACCEPTED / OPERATOR-RATIFIED
D5-R3 Product operation surface OPEN / NEXT
D6 / D7 / D8                   BLOCKED for NOTIF-01
Product implementation         BLOCKED UNTIL D9
```

**Exact next action:** open only the bounded D5-R3 Product operation-admission gate. Freeze the operation/Permission/client-class table before editing the canonical Product OpenAPI.