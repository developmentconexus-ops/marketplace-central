# NOTIF-01 D2-R5 — Typed Result & Requester Continuation Ratification

> **Status:** ACCEPTED / OPERATOR-RATIFIED 2026-08-23
> **Accepted artifact:** [D2-R5 Typed Result & Requester Continuation Repair](D6-R2-NOTIF-01-D2-R5-TYPED-RESULT-CONTINUATION.md), blob `fceecde8a207f8874038d65e4a292566f88be468`
> **Accepted base:** D2-R + D2-R2 + D2-R3 + D2-R4 remain ACCEPTED / OPERATOR-RATIFIED
> **Trigger:** [D5-F3 Global Maximum Operation / Permission Review](D6-R2-NOTIF-01-D5-F3-GLOBAL-MAXIMUM-OPERATION-PERMISSION-REVIEW.md)
> **Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## Ratified result

The operator ratified exactly the bounded D2-R5 repair required by the two proved result-bearing Notification families and the requester continuation falsifier:

1. `OFFERING_ASYNC_ACTION_RESULT` requires immutable `offering_async_result_outcome` with the closed values `converged | rejected | ambiguous | divergent`.
2. `AUTHORIZATION_DECISION_RESULT` requires immutable `authorization_decision_outcome` with the closed Governance values `authorize | reject`.
3. Both result atoms are kind-constrained presentation/history meaning only; they are never generic `result`, `status`, `reason`, metadata, routing, authorization, deduplication or source-mutation authority.
4. The source owner owns each result meaning; Personal Notifications owns only the immutable retained value after materialization. Replay/redelivery of the same `source_occurrence_key` must preserve the same result value.
5. `AUTHORIZATION_DECISION_RESULT.source_ref` changes from `AuthorizationDecisionRef` to the already-accepted closed `AuthorizationTargetRef` so the exact requester can continue in the action-owning target workspace without being granted `governance.read` merely for Notification navigation.
6. Governance remains the owner of the committed Authorization Decision occurrence. `source_occurrence_key` retains Governance-owned exact occurrence discrimination even though the Notification continuation points to the target.
7. `AuthorizationDecisionRef` leaves the closed Notification `source_ref` union because no admitted Launch-V1 NotificationKind still consumes it; the canonical Governance identity itself remains unchanged.
8. An immutable `authorize` result never claims that the target action later executed or converged. Current target navigation re-enters the target owner's current Product authorization and may fail closed after later access revocation.

No new Product operation, ordinary Permission, communication form, runtime mechanism, generic payload/template facility or frontend behavior is accepted by D2-R5.

## Gate

```text
D2-R / R2 / R3 / R4              ACCEPTED / OPERATOR-RATIFIED
D2-R5 typed result/continuation   ACCEPTED / OPERATOR-RATIFIED
D3-R2 bounded feed-forward        REVALIDATION NEXT
D5-R3 final five-operation table  BLOCKED until D3-R2 PASS
canonical Product OAD             UNCHANGED — 99/30
D6 / D7 / D8                      BLOCKED for NOTIF-01
Product implementation            BLOCKED UNTIL D9
```

**Exact next action:** revalidate only the accepted D3-R source-owner committed-fact `E` contracts for the two typed immutable result atoms and the F14 `AuthorizationTargetRef` continuation. Do not edit the Product OpenAPI first.