# D6-R2 — AuthorizationRequest Independent Fable Review Adjudication

> **Status:** GPT-ADJUDICATED CANDIDATE / OPERATOR RATIFICATION REQUIRED
> **Reviewed candidate:** `5ea52e223b2c4f077a15872300a54c20e89f4c6e`
> **Review branch:** `review/d6r2-authorization-request-fable`
> **Fable verdict:** `ACCEPT WITH BOUNDED FIXES`
> **Canonical census:** 106 Product operations · 31 ordinary Permissions · H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Review isolation and architectural verdict

The independent review was run from the exact P9 candidate. Revalidation after the Fable response established:

```text
review vs candidate   2 commits ahead / 0 behind
merge base            exact candidate 5ea52e22...
changed paths          docs/work/current/ai-dialog.md only
```

The review therefore remained isolated evidence and did not mutate candidate authority.

Fable independently challenged four credible structures:

```text
A  query-only target/revision projection
B  pending AuthorizationDecision
C  AuthorizationRequest + immutable AuthorizationDecision
D  generic approval/workflow engine
```

Its conclusion is that **C remains the Global-Maximum architecture under present Product requirements**. No material contradiction survives against D1-R2, D2-R6 or D3-R3; no generic workflow/request platform and no hidden 107th Product operation is justified.

GPT adjudication agrees with that architectural conclusion. The original P9-F1 root cause is structurally resolved rather than relocated.

## 2. Finding adjudication

### F-1 — Decision target/review-basis pairing proof blind spot

**Disposition: ACCEPTED — repository proof fix.**

The OAD already structurally pairs the four `AuthorizationDecision.target` kinds to their matching immutable `review_basis`, but the prior D5-R6 verifier did not directly falsify that invariant.

Bounded fix: supplemental bundle-level proof now asserts all four pairings and contains a mutation negative control that fails if the pairing is stripped.

No Product semantics change.

### F-2 — P9 operation homes not cross-checked against canonical OAD

**Disposition: ACCEPTED — repository proof fix.**

The prior P9 verifier proved the ten `P9_OP_HOME` markers inside P9 but did not prove those names still existed in the canonical bundled OAD.

Bounded fix: supplemental proof extracts all ten P9 homes and requires all ten exact `operationId`s in the bundled OAD, with a fake-operation negative control.

No Product semantics change.

### F-3 — no-consumer opaque correlation fields on actionable views

**Disposition: ACCEPTED AS BOUNDED D5 SHRINK — OPERATOR RATIFICATION REQUIRED.**

D5-R5 exposed these optional fields on the H-only actionable projection:

```text
requester_or_initiator_principal_id
predecessor_authorization_request_id
```

Neither B110 nor P9 has a legitimate current human consumer for them. The first is an opaque identity without a sanctioned presentation need; the second cannot be dereferenced through the actionable-only read after its predecessor becomes terminal. Adding a history endpoint/filter merely to justify those fields would violate YAGNI.

Candidate correction removes both fields from `ActionableAuthorizationRequestListItem` and all four actionable detail variants.

**D2-R6 is unchanged:** Governance still retains requester/initiator correlation and bounded predecessor lineage internally where they carry identity/history meaning. This supersedes only the two public actionable-projection fields from D5-R5 §4; it does not erase the canonical internal lineage.

Census remains 106/31.

### F-4 — semantic 503 indistinguishable from ambiguous infrastructure 503

**Disposition: ACCEPTED WITH STRONGER WIRE CORRECTION — OPERATOR RATIFICATION REQUIRED.**

The previous AuthorizationRequest overlay used `about:blank` for the semantic 503 meaning:

```text
material authorization-basis validity unavailable
→ no Decision recorded
→ request remains PENDING
```

That is weaker than MPC's existing convention for Product-semantic Problems and cannot safely distinguish a known-no-effect response from a proxy/bodyless/other 503 whose acceptance may be ambiguous.

Candidate correction defines the exact semantic Problem type:

```text
https://conexus.fun/marketplace-central/problems/product/authorization-validity-unavailable
```

Only a valid response carrying that exact Product Problem type is **known no-effect**. After reread/revalidation and a new explicit human confirmation, a new semantic attempt may use a new Idempotency-Key.

Any 503 that does not carry the exact typed semantic Problem — including bodyless, unparsable or infrastructure 503 — is **ambiguous potentially accepted**. It must not become a blind new attempt; same-semantic-attempt recovery preserves the same Idempotency-Key until authoritative convergence is established.

Census remains 106/31.

### F-5 — idempotent replay must precede If-Match evaluation

**Disposition: ACCEPTED AS D7-R OBLIGATION — not a D0-D6 defect.**

If a decision commits and its `201` is lost, the request revision has already advanced. Replaying the exact same semantic attempt with the same Idempotency-Key must resolve the committed intake before evaluating the now-stale If-Match; otherwise a precondition-first runtime returns 412 and defeats the admitted ambiguous-acceptance recovery contract.

D7-R must prove:

```text
same key + same semantic command + committed intake
→ return/recover original 201 Decision
→ before current-resource If-Match rejection
```

A genuinely different attempt does not receive this replay treatment.

## 3. RED → GREEN evidence

Supplemental verifier:

`script/verify-authorization-request-fable-fixes.mjs` (repository path: `scripts/verify-authorization-request-fable-fixes.mjs`)

RED:

```text
HEAD 7a974f98baa16b2b34e2e492740e8f3754c0b0f2
CI #603 FAILURE
all previous 95/29 + 99/30 + 106/31 + P8 + P9 proofs PASS first
failure = actionable AuthorizationRequest wire leaked requester_or_initiator_principal_id
```

GREEN candidate:

```text
HEAD 00cba52fd6956169c25864c3a0d9dfc37dd578b0
CI #604 SUCCESS
pr-title #675 SUCCESS

authorization_request_fable_decision_pairing=PASS
authorization_request_fable_actionable_projection=MINIMAL
authorization_request_fable_503_discriminator=SPECIFIC_PROBLEM_TYPE
authorization_request_fable_p9_wire_crosscheck=10/10
authorization_request_fable_negative_controls=4/4
authorization_request_fable_bounded_fixes=PASS
```

Historical 95/29 + 99/30 and current 106/31 remain green.

## 4. D7-R feed-forward obligations

The independent review identified runtime obligations, not missing Product semantics. D7-R must explicitly consume and prove at least:

1. committed same-key replay resolution before If-Match evaluation;
2. typed semantic validity-unavailable 503 vs ambiguous/non-semantic 503 handling;
3. current eligible-decider derivation freshness without delayed events becoming sole revocation authority;
4. material-validity Q timeout/unavailability realization matching the typed 503 contract;
5. zero-decider Work materialization and reconciliation;
6. F13/F14 delivery plus currentness revalidation;
7. invalidation reconciliation preventing silent stale-PENDING requests;
8. Idempotency-Key storage, scope and retention sufficient for admitted replay semantics.

These obligations do not select a server framework, queue, database schema or transport mechanism here.

## 5. Exact remaining gate

F-1/F-2 are proof-only fixes and require no architecture reopen. F-5 is a D7-R feed-forward obligation.

F-3 and F-4 modify bounded parts of the operator-ratified D5-R5 Product wire. Therefore this GPT adjudication **does not self-ratify them**.

Exact remaining human decision:

```text
operator ratifies F-3 bounded public-projection shrink
+ operator ratifies F-4 typed semantic-503 contract
+ operator accepts the independent-review adjudication
→ AuthorizationRequest redesign may be declared GLOBAL-MAXIMUM CLOSED
→ D7-R becomes NEXT, not silently started
```

Until that explicit operator decision: D7-R/D8-R, B10 resumption, merge and Product implementation remain blocked.
