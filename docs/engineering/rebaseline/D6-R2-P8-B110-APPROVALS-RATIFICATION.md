# D6-R2 P8 — B110 Approvals Operator Ratification

> **Status:** OPERATOR-RATIFIED / ACCEPTED / LOCKED
> **Operator adjudication:** 2026-08-24
> **Reviewed candidate:** [B110 Approvals Candidate](D6-R2-P8-B110-APPROVALS-CANDIDATE.md)
> **P5 input:** [B110 AuthorizationRequest Targeted Supersession](D6-R2-P5-B110-AUTHORIZATION-REQUEST-SUPERSESSION.md)
> **Canonical wire:** [D5-R6 AuthorizationRequest OAD Proof](D6-R2-NOTIF-01-D5-R6-AUTHORIZATION-REQUEST-OAD-WIRE-PROOF.md)
> **Reviewed artifact:** [`qualification/d6-r2-wireframes/b110-approvals.html`](../../../qualification/d6-r2-wireframes/b110-approvals.html)
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Operator decision

The operator visually adjudicated the exact rendered B110 candidate and responded **“Aprovado”** after interactive review. Under Frontend Method v2.1, this is the explicit operator `LOCK` for B110 only.

The reviewed HTML remains an immutable `data-p8-status="candidate"` evidence snapshot. It must not be edited to self-claim operator authority; this ratification record owns the `LOCK` disposition.

## 2. Locked B110 contract

B110 locks the existing `CONTROLE > Aprovações` destination with two local, independently authorized lenses:

```text
Para decidir → governance.decide
Histórico    → governance.read
```

Neither Permission implies the other.

Locked laws:

- one global `Aprovações` destination only; no second global history destination;
- `Para decidir` is the exact-human actionable `AuthorizationRequest` queue using a structured list and cursor continuation only;
- request detail is a real route and renders exactly one of the four typed review-basis families: `listing_intent`, `price_intent`, `business_order_intent`, `invoicing_intent`;
- `CreateAuthorizationDecision` uses inline confirmation with evidence still visible, request-local `If-Match`, `Idempotency-Key`, and outcome-only client body;
- consequential decision auto-retry is forbidden;
- stale/412 rereads Governance truth and never silently overwrites or leaks history without `governance.read`;
- validity unavailable/503 records no Decision and leaves the request pending;
- successful Decision does not execute the source action; the action owner retains execution-time revalidation/authority;
- `Abrir origem` reauthorizes the current source owner and Governance never grants source-read by implication;
- `Histórico` is immutable `AuthorizationDecision` truth with admitted date filters and no decision controls;
- F13 deep-links by `AuthorizationRequestRef`, but Notification remains awareness, never a capability token;
- Organization switch invalidates request/history transient context;
- mobile stacks decision actions without changing authority semantics;
- no approval search, totals, bulk decision model, approver filter, generic review payload or workflow/lifecycle platform is admitted.

## 3. Executable evidence

TDD sequence:

```text
RED  CI #596
all prior authority/OAD/P8 proofs PASS
B110 rendered artifact missing

GREEN HEAD db0c87bc110bf0a05185b23c0c3841e38c6f8579
CI #597 SUCCESS
pr-title #666 SUCCESS
```

Final reviewed candidate checkpoint before operator lock:

```text
HEAD 3e9c24b35f3449558512f22d8e1e5a931620c97b
CI #598 SUCCESS
pr-title #667/#668 SUCCESS
```

Structural proof:

```text
d6_r_b110_lenses=ACTIONABLE+HISTORY
d6_r_b110_permissions=governance.decide|governance.read_INDEPENDENT
d6_r_b110_review_basis=4/4
d6_r_b110_confirmation=INLINE
d6_r_b110_stale_and_503=EXPLICIT
d6_r_b110_notification=AWAWARENESS_NOT_CAPABILITY
d6_r_b110_wireframe=PASS
```

Historical 95/29 + 99/30 and current 106/31 Product OAD proofs remain green in the same gate.

## 4. Supersession and next gate

This record supersedes only the prior B110 `RENDERED CANDIDATE / NOT LOCKED` disposition. It does not imply D6-R2 closure, D7-R authorization, D8-R authorization, merge authorization, B10 resumption or Product implementation.

Exact next work is the final P9 Screen Contract + bidirectional backend trace over the locked affected frontend surfaces and canonical 106/31 Product authority. After P9, the complete AuthorizationRequest package must receive independent Fable review and finding adjudication before Global-Maximum closure or D7-R may open.
