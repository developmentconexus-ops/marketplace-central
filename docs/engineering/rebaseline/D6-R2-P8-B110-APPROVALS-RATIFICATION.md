# D6-R2 P8 — B110 Approvals Operator Ratification

> **Status:** OPERATOR-RATIFIED / ACCEPTED / LOCKED / CURRENT EVIDENCE  
> **Operator adjudication:** 2026-08-24  
> **Artifact:** `qualification/d6-r2-wireframes/b110-approvals.html`  
> **Current Product:** 106 operations / 31 ordinary Permissions / H-A-S

## 1. Operator LOCK

The operator visually approved the rendered B110 structure. The immutable HTML may still self-label `candidate`; this record owns the LOCK.

B110 preserves one `CONTROLE > Aprovações` destination with two independently authorized lenses:

```text
Para decidir → governance.decide
Histórico    → governance.read
```

Neither Permission implies the other.

## 2. Locked human structure

### Para decidir

- exact-human actionable `AuthorizationRequest` queue;
- structured list + cursor only;
- request detail is a real route;
- exactly one typed review-basis family: `listing_intent | price_intent | business_order_intent | invoicing_intent`;
- Approve/Reject uses inline confirmation with evidence still visible;
- consequential decision auto-retry is forbidden;
- current source continuation reauthorizes the source owner;
- successful AuthorizationDecision does **not** execute the target action;
- F13 deep-link by `AuthorizationRequestRef` is awareness only, not capability.

### Histórico

- immutable `AuthorizationDecision` history;
- independent `governance.read`;
- admitted date filters only;
- no Approve/Reject controls.

No approval search, total count, bulk decision, approver filter, generic review payload or workflow/case platform is admitted.

## 3. Current decision carrier after bounded W1 repair

The later operator-ratified W1 correction changed **only the transport carrier**, not B110 structure:

```text
POST /organizations/{organization_id}/authorization-requests/{authorization_request_id}:decide
Idempotency-Key: required
body:
  etag: current AuthorizationRequest StrongETag
  outcome: authorize | reject
```

Current failure grammar:

```text
missing/invalid body.etag → 422
stale Request revision    → 409 resource-revision-conflict
current state conflict    → 409
known validity unavailable→ typed 503, known no Decision recorded
```

There is no `If-Match`, 412 or 428 on this custom `:decide` operation.

This change was explicitly revalidated as **STRUCTURE UNAFFECTED / P8 reopen NO**. The same human stale/recovery state remains visible; only the underlying technical carrier/status spelling changed.

## 4. Scope/authority laws

- `governance.decide` does not grant Governance history or source-owner read;
- `governance.read` does not grant decision authority;
- Notification does not supply Request ETag, decision eligibility or source access;
- Work may represent a zero-decider operational blocking condition but never becomes approver;
- Organization switch invalidates transient Request/history context;
- mobile stacking does not change decision/source authority.

Exact current wire/state binding lives in the current P9 contract. Mutable stage/next action belongs only to `docs/roadmap.md`.
