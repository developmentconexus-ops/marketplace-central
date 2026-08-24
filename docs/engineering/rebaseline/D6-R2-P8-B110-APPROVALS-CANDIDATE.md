# D6-R2 P8 — B110 Approvals AuthorizationRequest Candidate

> **Status:** RENDERED CANDIDATE / VISUAL ADJUDICATION REQUIRED / NOT LOCKED
> **P5 input:** [B110 AuthorizationRequest Targeted Supersession](D6-R2-P5-B110-AUTHORIZATION-REQUEST-SUPERSESSION.md)
> **Canonical wire:** [D5-R6 AuthorizationRequest OAD Proof](D6-R2-NOTIF-01-D5-R6-AUTHORIZATION-REQUEST-OAD-WIRE-PROOF.md)
> **Artifact:** [`qualification/d6-r2-wireframes/b110-approvals.html`](../../../qualification/d6-r2-wireframes/b110-approvals.html)
> **Verifier:** `scripts/verify-d6-r-b110-approvals-wireframe.mjs`
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Trigger and boundary

D2-R6/D5-R6 materially changed the Approvals human model from target/revision-shaped decision handling to a canonical pre-decision `AuthorizationRequest` plus separate immutable `AuthorizationDecision` history. Method v2.1 therefore requires a targeted P8 reopen before final P9.

This block does **not** reopen the global sidebar, B00 shell, Notifications blocks, Delegations Settings, Product OAD or final visual design.

## 2. Rendered candidate structure

```text
CONTROLE > Aprovações

local lens: Para decidir
→ structured list of exact-human actionable AuthorizationRequests
→ request detail route
→ one of four typed review-basis regions
→ Rejeitar | Aprovar
→ inline confirmation while evidence remains visible
→ explicit stale / validity-unavailable / success states

local lens: Histórico
→ immutable AuthorizationDecision list
→ decided_from / decided_before
→ Decision detail
→ no decision controls
```

The two lenses remain independently conditioned:

```text
governance.decide → Para decidir
governance.read   → Histórico
```

Neither Permission implies the other.

## 3. Locked authority preserved by the candidate

- one global `Aprovações` destination only;
- no global `Histórico de aprovações` destination;
- Organization remains the only global workspace;
- server authorization remains authoritative despite hidden local lenses;
- actionable queue is a structured recognition/review list, not a comparison table;
- cursor continuation only; no total, search, approver filter, lifecycle platform or bulk decisions;
- F13 deep-links with `AuthorizationRequestRef`, but Notification remains awareness, never a capability token;
- source continuation reauthorizes current owner truth and does not gain source-read from Governance.

## 4. Typed review-basis proof

The rendered source statically demonstrates all four D5-R6 variants:

```text
listing_intent
price_intent
business_order_intent
invoicing_intent
```

No generic review payload/basis is admitted. The structure uses one consistent decision grammar while the evidence region remains type-specific.

## 5. Consequential decision law

`CreateAuthorizationDecision` is represented as:

```text
If-Match        → AuthorizationRequest concurrency
Idempotency-Key → same-command ambiguous retry protection
body            → outcome only
```

No governed target ETag is client-authored. The candidate uses inline confirmation rather than a mandatory modal so the human keeps the exact review context visible while confirming.

Recovery remains fail-honest:

- stale/412 → reread current Governance truth; no blind resubmission and no history disclosure without `governance.read`;
- material-validity unavailable/503 → no Decision recorded, request remains pending, reread/revalidation is safe, consequential write auto-retry is forbidden;
- success → Decision recorded; source open is a separate reauthorized continuation.

## 6. Deterministic visual scenarios

The HTML exposes these operator-testable states:

```text
actionable-populated
actionable-empty
request-price
request-listing
request-business-order
request-invoicing
inline-confirmation
stale-decision
validity-unavailable
decision-recorded
history-list
history-detail
decide-only
read-only
f13-deep-link
organization-switch
```

Mobile stacks approval actions without changing authority meaning.

## 7. Executable proof

RED first:

```text
CI #596
all prior authority/proofs PASS
B110 rendered artifact missing
```

GREEN on first rendered artifact:

```text
HEAD db0c87bc110bf0a05185b23c0c3841e38c6f8579
CI #597 SUCCESS
pr-title #666 SUCCESS

d6_r_b110_status=CANDIDATE
d6_r_b110_lenses=ACTIONABLE+HISTORY
d6_r_b110_permissions=governance.decide|governance.read_INDEPENDENT
d6_r_b110_review_basis=4/4
d6_r_b110_confirmation=INLINE
d6_r_b110_stale_and_503=EXPLICIT
d6_r_b110_notification=AWAWARENESS_NOT_CAPABILITY
d6_r_b110_wireframe=PASS
```

Historical 95/29 + 99/30 and current 106/31 OAD proofs also pass in the same full gate.

## 8. Operator gate

`B110` remains **NOT LOCKED**. The exact next action is operator visual walkthrough of the rendered HTML. Final P9 cannot begin until explicit operator `LOCK`; Fable review remains after final P9 and before Global-Maximum closure / D7-R.
