# NOTIF-01 D2-R6 — Authorization Request Ratification

> **Status:** OPERATOR-RATIFIED / ACCEPTED
> **Ratified candidate:** [D2-R6 Authorization Request Identity & Lifecycle](D6-R2-NOTIF-01-D2-R6-AUTHORIZATION-REQUEST-IDENTITY-LIFECYCLE.md)
> **Ratified candidate blob:** `41b06b60dcfd3c3bb391d4cb05de6102bf5b870c`
> **Parent authority:** D2 Identity / Tenant / Data Ownership
> **Boundary revalidation:** [D1-R2](D6-R2-NOTIF-01-D1-R2-AUTHORIZATION-REQUEST-BOUNDARY-REVALIDATION.md) — CURRENT STRUCTURE CONFIRMED
> **Canonical Product wire:** unchanged — 104 Product operations · 31 ordinary Permissions · H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## Ratified contract

The operator ratifies the exact D2-R6 candidate with these binding outcomes:

1. `AuthorizationRequest` is a canonical MPC-owned identity under **Controlled Action Governance**.
2. `AuthorizationRequest` is distinct from the action-owning Business Intent and from terminal `AuthorizationDecision`.
3. One material authorization episode has one stable opaque `AuthorizationRequestID`; a later reauthorization episode receives a new ID even when the governed target and its coarse resource revision are unchanged.
4. Baseline lifecycle is exactly `PENDING | DECIDED | INVALIDATED`.
5. One request has `0..1` terminal `AuthorizationDecision`; a decided or invalidated request is never reopened.
6. `INVALIDATED` is not a fabricated human rejection. It represents a request that can no longer truthfully be decided under its preserved authorization basis.
7. Each request retains one immutable, bounded, typed authorization review/basis snapshot sufficient for truthful human review and later explanation. It is not source-owner current truth, a general read mirror, arbitrary JSON, generic metadata or a capability token.
8. Current eligible decision Principals remain dynamic Governance-owned meaning. They are not immutable request recipients and may change while request identity remains stable.
9. Request concurrency and business validity are separate meanings. The request owns an owner-local revision for stale/concurrent decision protection.
10. A whole-target resource ETag/revision is **not** the authorization-validity oracle. Irrelevant target drift must not force reapproval; material governing drift must not be ignored merely because a coarse target revision stayed unchanged.
11. Before a decision is committed, Governance must revalidate current decision eligibility plus material validity of the preserved authorization episode through the accepted action-owner boundary when source/business meaning is required.
12. Material reauthorization may preserve one bounded `predecessor_authorization_request_id`; this is not a generic relation graph.
13. Retry/recovery of the same semantic authorization episode resolves to the same request; genuinely new authorization requires a new episode.
14. `AuthorizationDecision` remains immutable historical decision occurrence and belongs to exactly one request.
15. Candidate Notification feed-forward is `AUTHORIZATION_ACTION_REQUIRED → AuthorizationRequestRef`; `AUTHORIZATION_DECISION_RESULT` is not changed by D2-R6 alone.
16. No workflow engine, BPMN, quorum/voting, approval-stage DSL, generic case management, generic entity graph, request comments/chat or speculative expiry model is admitted.

## Authority result

D1 remains unchanged: Controlled Action Governance was already the correct semantic owner. The repair is identity/lifecycle only; action disposition, Business Intent and execution-time validity remain action-owner authority.

The Product OAD remains **104/31**. D2-R6 does not select HTTP paths, Product operation count, storage schema, runtime topology or frontend mechanics.

## Required continuation

```text
D2-R6 ACCEPTED
→ D3-R3 communication/recovery feed-forward
→ D5 bounded Product-wire repair + executable proof
→ D6 P9 Screen Contracts / bidirectional trace closure
→ independent Fable review + adjudication of the entire AuthorizationRequest redesign
→ only then may the redesign be declared Global-Maximum closed and D7-R open
```

The independent Fable review is a **required closure gate**, not optional commentary. The review pack must include at minimum D1-R2, D2-R6, D3-R3, the resulting D5/OAD repair, final P9 trace, and Notification/Operational-Work consequences. Findings are evidence and must be adjudicated against current authority before closure.

## Exact next action

Derive and adjudicate only D3-R3. Do not modify D5/OpenAPI, resume P9/B10, begin D7-R/D8-R or implement Product code first.