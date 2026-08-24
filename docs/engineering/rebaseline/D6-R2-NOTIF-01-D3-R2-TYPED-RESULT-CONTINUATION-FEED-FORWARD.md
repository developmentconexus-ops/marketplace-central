# NOTIF-01 D3-R2 — Typed Result & Requester Continuation Feed-Forward Revalidation

> **Status:** PASS / NO COMMUNICATION-TOPOLOGY REOPEN REQUIRED
> **Accepted parents:** [D3-R Communication & Propagation Ratification](D6-R2-NOTIF-01-D3-R-RATIFICATION.md) + [D2-R5 Ratification](D6-R2-NOTIF-01-D2-R5-RATIFICATION.md)
> **Scope:** revalidate only F02/F14 immutable result propagation and F14 requester continuation identity inside the already-accepted source-owner committed-fact `E` model
> **Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Question

D2-R5 adds only:

```text
OFFERING_ASYNC_ACTION_RESULT
  + offering_async_result_outcome

AUTHORIZATION_DECISION_RESULT
  + authorization_decision_outcome
  + source_ref = AuthorizationTargetRef
```

D3-R already permits source-owned `E` contracts to carry the smallest stable immutable occurrence facts required by Personal Notifications, while forbidding mutable aggregate mirrors, free-form payloads, generic envelopes and source-authority transfer.

The bounded question is:

> Can these exact typed atoms and the already-accepted target reference cross the existing F02/F14 committed-fact `E` boundaries while preserving owner authority, replay/recovery, recipient derivation and immutable-occurrence semantics?

**Result: YES / PASS.**

---

## 2. F02 — Offering async result

The existing source-owned contract remains:

```text
Marketplace Offering Operations
  OfferingAsyncActionResultCommitted E
        ↓
Personal Notifications
  OFFERING_ASYNC_ACTION_RESULT
```

Its minimum occurrence contract now additionally carries exactly:

```text
offering_async_result_outcome
  = converged | rejected | ambiguous | divergent
```

Laws:

- Offering owns the result meaning at source commit.
- The value is immutable for the exact `source_occurrence_key`.
- A materially different later result in the accepted uncertainty lineage is a distinct committed occurrence, not mutation of the earlier result atom.
- The result remains historical/presentation meaning; it does not authorize or mutate Offering state.
- F02 remains an immutable-occurrence family and still requires no source-currentness Q merely to materialize the committed historical result.

No `result:string`, reason blob, provider payload or source reread is introduced.

---

## 3. F14 — Governance requester result

The existing source owner and communication form remain:

```text
Controlled Action Governance
  AuthorizationDecisionResultCommitted E
        ↓
Personal Notifications
  AUTHORIZATION_DECISION_RESULT
```

The occurrence contract now carries:

```text
source_ref = AuthorizationTargetRef

authorization_decision_outcome
  = authorize | reject
```

while:

```text
source_occurrence_key
  = stable Governance-owned discriminator for the exact committed decision occurrence
```

This preserves the authority split:

- Governance owns the decision occurrence and immutable `authorize | reject` meaning.
- The action-owning target remains the requester continuation subject.
- Carrying `AuthorizationTargetRef` does not transfer Governance decision authority to the target owner or Personal Notifications.
- The requester is not granted `governance.read` merely because Governance produced the result.
- Navigation from the Notification re-enters the target owner's current Product authorization.

F14 remains an immutable-occurrence family and still does not require source-currentness Q merely to preserve the committed decision result.

---

## 4. Revised bounded occurrence shapes

The accepted common NOTIF-01 occurrence atoms remain unchanged:

```text
organization_id
source_ref
source_occurrence_key
source_occurred_at
source_committed_at
subject_display_label
```

Only the exact families add their constrained atom:

```text
OfferingAsyncActionResultCommitted
  + offering_async_result_outcome

AuthorizationDecisionResultCommitted
  + authorization_decision_outcome
```

`AuthorizationDecisionResultCommitted.source_ref` is now `AuthorizationTargetRef` under D2-R5.

No other family gains either result field by symmetry.

---

## 5. Replay, recovery and ordering

Accepted D3-R behavior is unchanged:

- duplicate/redelivery of one source occurrence converges semantically and preserves the same typed result atom;
- recovery must reconstruct the exact accepted immutable result atom together with the occurrence rather than recomputing it from later mutable source state;
- no global event ordering, EventID, SagaID or exactly-once delivery claim is introduced;
- DIRECT_SOURCE recipient derivation for F02/F14 remains the exact producer-owned human recipient;
- current Organization Membership/Product-access eligibility remains independently authoritative before materialization;
- the F14 target reference never substitutes for the Governance-owned occurrence key;
- D2-R2/D2-R3 routing and supersession semantics are unaffected because F02/F14 are DIRECT_SOURCE immutable-occurrence families.

---

## 6. Negative controls

D3-R2 fails if realization introduces any of:

```text
result: string
reason: string
status: string
metadata: {}
payload: {}
template_variables: {}
source mutable snapshot
current-source reread to reconstruct historical result
AuthorizationDecisionRef as F14 navigation workaround
governance.read grant merely for Notification continuation
result-based routing or authorization
generic event envelope as Product authority
new communication form or source-owner edge
```

It also fails if replay of one `source_occurrence_key` can carry a different immutable result value.

---

## 7. Coherence result

D2-R5 is fully expressible inside the accepted D3 topology:

```text
Offering commits exact async result
  + typed immutable outcome
        ↓ E
Personal Notifications

Governance commits exact decision result
  + AuthorizationTargetRef
  + typed immutable decision outcome
        ↓ E
Personal Notifications
```

No new D1 edge, source owner, Q/C/P form, broker, event platform, public event API, runtime mechanism or generic payload is required.

Therefore **D3-R remains ACCEPTED and D3-R2 is PASS / NO TOPOLOGY REOPEN**.

## 8. Gate

```text
D2-R5 typed result/continuation   ACCEPTED / OPERATOR-RATIFIED
D3-R2 bounded feed-forward        PASS / NO TOPOLOGY REOPEN
D5-R3 final five-operation table  OPEN / NEXT
canonical Product OAD             UNCHANGED — 99/30
D6 / D7 / D8                      BLOCKED for NOTIF-01
Product implementation            BLOCKED UNTIL D9
```

**Exact next action:** rebuild and operator-adjudicate the final D5-R3 NOTIF-01 operation-admission table using D5-F3 surviving findings, accepted D5-F4 recipient discovery, accepted D2-R5 and D3-R2 PASS. Do not edit the canonical Product OpenAPI first.