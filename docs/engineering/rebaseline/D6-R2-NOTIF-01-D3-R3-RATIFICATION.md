# NOTIF-01 D3-R3 — Authorization Request Communication Ratification

> **Status:** OPERATOR-RATIFIED / ACCEPTED
> **Ratified candidate:** [D3-R3 Authorization Request Communication & Recovery](D6-R2-NOTIF-01-D3-R3-AUTHORIZATION-REQUEST-COMMUNICATION-RECOVERY.md)
> **Accepted identity authority:** [D2-R6 Ratification](D6-R2-NOTIF-01-D2-R6-RATIFICATION.md)
> **Boundary authority:** [D1-R2](D6-R2-NOTIF-01-D1-R2-AUTHORIZATION-REQUEST-BOUNDARY-REVALIDATION.md) — CURRENT STRUCTURE CONFIRMED
> **Canonical Product wire:** unchanged — 104 Product operations · 31 ordinary Permissions · H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## Ratified communication contract

1. A new authorization episode is requested from the action owner to Governance by **C** and carries a stable owner-local authorization episode anchor. Ambiguous retry/recovery of the same episode resolves to the same `AuthorizationRequest`; a genuinely new reauthorization episode uses a new anchor and new `AuthorizationRequestID`.
2. Current decision eligibility is dynamic Governance-owned truth. It is resolved from Governance authority/delegation plus current identity/access truth and is revalidated on every decision attempt.
3. `AUTHORIZATION_ACTION_REQUIRED` is driven by a Governance-owned current human decision-eligibility episode. Personal Notifications receives only Organization, `AuthorizationRequestRef`, exact human Principal, occurrence discriminator and source time.
4. A delayed/replayed F13 cannot create new actionable awareness after the request is terminal or the human is no longer eligible. Existing historical Notifications are not rewritten.
5. A `PENDING` request with a **known-empty** current eligible-human set is immediately a material blocking Governance condition. There is no fallback approver, default admin, role inference, notification routing fallback or unowned grace period. Governance feeds the existing Operational Work boundary; Work owns the obligation lifecycle but never authorization.
6. A human decision binds to `AuthorizationRequestID` plus request-local concurrency revision/precondition. Concurrent deciders cannot create multiple terminal Decisions.
7. Before committing a Decision, Governance revalidates request state/revision, exact current human eligibility and the material validity of the preserved authorization episode. Current source/business validity is queried from the action owner when required.
8. Action-owner validity is semantically `VALID | INVALID | UNKNOWN_OR_UNAVAILABLE`. `UNKNOWN_OR_UNAVAILABLE` never defaults to valid.
9. `INVALID` terminates a still-current request as `INVALIDATED`; it never fabricates a human `reject`.
10. A committed `AuthorizationDecision` is a Governance-owned occurrence and propagates by recoverable **E** to the action owner. `authorize` supplies authorization evidence only; it never executes the governed action or waives execution-time validity.
11. The same committed Decision may independently feed F14 to the exact human requester when accepted lineage permits. Duplicate Decision propagation cannot duplicate the semantic result Notification.
12. Action-owner material invalidation may awaken Governance by **E**, followed by Q/revalidation where currentness matters. Governance alone writes request lifecycle.
13. Request invalidation propagates recoverably to the action owner when owner-local progression depends on it. No F14 decision-result Notification exists without an actual Decision.
14. Reauthorization is action-owner meaning and always creates a new episode/request; terminal historical requests and Decisions are never reopened or rewritten.
15. Event transport, queues, workers, persistence and transactions remain D7 mechanisms. Event delivery is never the sole authority for current eligibility, current validity or recovery.

## Authority fences

D3-R3 does not move Business Intent, business disposition, source current truth, execution-time validity, ordinary access, Work assignment, Notification state or provider protocol into Governance.

No workflow engine, generic command bus, generic case platform, universal event identity, exactly-once claim, distributed transaction or fallback approver is admitted.

## Product consequence

No D5 consequence is preselected by this ratification. The canonical Product OAD remains **104/31**. D5 must independently derive the smallest sustainable Product surface from the accepted D1-R2 + D2-R6 + D3-R3 semantics.

## Required Global-Maximum closure review

The previously ratified requirement remains binding:

```text
D3-R3 ACCEPTED
→ D5 bounded wire derivation + executable OAD proof
→ D6 P9 final Screen Contracts / bidirectional trace
→ independent Fable review of the complete AuthorizationRequest redesign
→ adjudicate every material finding
→ only then may the redesign be declared Global-Maximum closed and D7-R open
```

The Fable pack must include D1-R2, D2-R6, D3-R3, the resulting D5/OAD repair, final P9 trace, Notification consequences and the zero-decider Operational Work path.

## Exact next action

Derive the D5 Product surface in design first. Do not edit the canonical OAD until the operator approves the D5 surface. Do not resume P9/B10, begin D7-R/D8-R, merge PR #61 or implement Product code first.
