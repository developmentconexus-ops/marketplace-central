# D6-R2 — AuthorizationRequest Fable Closure Operator Ratification

> **Status:** OPERATOR-RATIFIED / ACCEPTED / GLOBAL-MAXIMUM CLOSED
> **Operator adjudication:** 2026-08-24
> **Adjudication evidence:** [AuthorizationRequest Independent Fable Review Adjudication](D6-R2-AUTHORIZATION-REQUEST-FABLE-ADJUDICATION.md)
> **Reviewed candidate:** `5ea52e223b2c4f077a15872300a54c20e89f4c6e`
> **Final bounded-fix checkpoint before operator ruling:** `84dca874eeb1ef272568b1c30af6c2ada2260a03`
> **Canonical Product:** 106 Product operations · 31 ordinary Permissions · H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Operator decision

The operator explicitly responded **“Aprovado”** after the independent Fable review, GPT adjudication, bounded F-3/F-4 Product-wire corrections and executable GREEN proof were presented.

This ratifies exactly:

- F-3: remove `requester_or_initiator_principal_id` and `predecessor_authorization_request_id` from the public H-only actionable AuthorizationRequest projections because no current Product consumer can legitimately use or dereference them;
- D2-R6 internal Governance requester/initiator correlation and bounded predecessor lineage remain canonical and unchanged;
- F-4: the known-no-effect material-validity 503 uses the exact Product Problem type `https://conexus.fun/marketplace-central/problems/product/authorization-validity-unavailable`;
- any 503 without that exact parseable semantic discriminator is ambiguous potentially accepted and cannot become a blind new decision attempt;
- F-1/F-2 proof fixes are accepted;
- F-5 and the remaining Fable runtime list are accepted as D7-R obligations, not D0-D6 defects;
- the independent review/adjudication is accepted and closed.

## 2. Global-Maximum closure

The challenged alternatives remain rejected under present Product requirements:

```text
A  query-only target/revision projection        REJECTED — preserves P9-F1 root cause
B  pending AuthorizationDecision                REJECTED — identity/lifecycle + Permission collision
C  AuthorizationRequest + immutable Decision    ACCEPTED — GLOBAL MAXIMUM
D  generic approval/workflow engine              REJECTED — no present consumer / YAGNI
```

No material contradiction survives against D1-R2, D2-R6, D3-R3, D5-R6, locked B110 or final P9. No hidden 107th Product operation or new ordinary Permission is required.

Therefore the AuthorizationRequest redesign is **GLOBAL-MAXIMUM CLOSED** for the current Product authority.

## 3. Proof accepted by this ratification

The bounded fixes were already GREEN before this operator ruling:

```text
RED   7a974f98baa16b2b34e2e492740e8f3754c0b0f2 / CI #603 FAILURE
GREEN 00cba52fd6956169c25864c3a0d9dfc37dd578b0 / CI #604 SUCCESS

final adjudication checkpoint
84dca874eeb1ef272568b1c30af6c2ada2260a03
CI #606 SUCCESS
pr-title #677/#678 SUCCESS
```

Accepted proof properties include:

```text
Decision target/review pairing falsifier         PASS
actionable public projection                     MINIMAL
semantic 503 discriminator                       SPECIFIC_PROBLEM_TYPE
P9 operation-home ↔ canonical OAD cross-check    10/10
Fable supplemental negative controls             4/4
historical Product proof                         95/29 + 99/30 PASS
current Product proof                            106/31 PASS
```

## 4. D7-R feed-forward

D7-R must consume the eight runtime obligations already recorded in the adjudication, including the critical replay ordering law:

```text
same Idempotency-Key + same semantic command + committed intake
→ recover original committed Decision/result
→ before rejecting on the now-stale If-Match
```

This record does not select database schema, transaction primitive, queue, scheduler, worker, router, framework or deployment mechanics.

## 5. Scope of this approval

This approval closes only the AuthorizationRequest/Fable redesign gate.

It does **not**:

- close D6-R2 as a whole;
- resume B10;
- complete D7-R or D8-R;
- authorize merge of PR #61;
- authorize Product implementation;
- bypass D9.

Exact continuation is the bounded NOTIF-01 D7-R runtime/jobs/transactions repair against accepted D7 authority and the recorded obligations. D8-R remains blocked by D7-R; B10 remains suspended.