# NOTIF-01 D5-R5 — AuthorizationRequest Product Surface Ratification

> **Status:** OPERATOR-RATIFIED / ACCEPTED SURFACE — D5-R6 OAD WIRE PROOF REQUIRED
> **Accepted parents:** [D2-R6](D6-R2-NOTIF-01-D2-R6-RATIFICATION.md) + [D3-R3](D6-R2-NOTIF-01-D3-R3-RATIFICATION.md)
> **Trigger:** D6 P9 actionable-approval falsifier and operator-approved Global-Maximum restructure
> **Canonical Product OAD before proof:** 104 Product operations · 31 ordinary Permissions · H/A/S
> **Expected consequence if proof passes:** 106 Product operations · 31 ordinary Permissions · H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Ratified Product-surface decision

The operator ratifies the smallest Product surface that makes the accepted `AuthorizationRequest` identity usable without exposing D3 owner-to-owner communication as a public workflow API.

Exactly two new Product operations are admitted:

```text
ListMyActionableAuthorizationRequests
GetMyActionableAuthorizationRequest
```

The existing `CreateAuthorizationDecision` remains one Product operation but is reanchored from a caller-supplied governed target to the canonical `AuthorizationRequest` resource.

No new ordinary Permission is admitted.

## 2. Exact operation contract

### 2.1 `ListMyActionableAuthorizationRequests`

```text
GET /organizations/{organization_id}/authorization-requests
class       Q
owner       ControlledActionGovernance
Permission  governance.decide
Principal   H only
```

Meaning:

> Return only `PENDING` AuthorizationRequests that the exact current human Principal is currently eligible to decide in the exact path Organization.

This is a purpose-bounded self-actionability read, not Governance history/read access.

Baseline query controls are exactly:

```text
limit
cursor
```

No search, generic state filter, assignee filter, target-owner filter, Permission/role filter, total count or saved view is admitted.

### 2.2 `GetMyActionableAuthorizationRequest`

```text
GET /organizations/{organization_id}/authorization-requests/{authorization_request_id}
class       Q
owner       ControlledActionGovernance
Permission  governance.decide
Principal   H only
```

The operation succeeds only while the exact request is currently actionable by the exact human caller. Possession of an `AuthorizationRequestID` or Notification never grants access.

Success returns a purpose-bounded actionable request view plus a strong ETag representing the **AuthorizationRequest owner-local revision**. The ETag does not represent the whole target resource and is not the material-authorization-validity oracle.

### 2.3 `CreateAuthorizationDecision`

Canonical wire target becomes:

```text
POST /organizations/{organization_id}/authorization-requests/{authorization_request_id}:decide
class       C
owner       ControlledActionGovernance
Permission  governance.decide
Principal   H only
```

Required request trust/concurrency carriers are both:

```text
If-Match        = current AuthorizationRequest validator
Idempotency-Key = stable client intake key
```

They solve different correctness problems:

```text
If-Match        → concurrent/stale decision protection
Idempotency-Key → ambiguous acceptance / exact retry recovery
```

The JSON body contains only:

```json
{"outcome":"authorize | reject"}
```

The human client does not supply target identity, target ETag, effective Principal, authority/delegation facts, `authorized=true`, source current state or material-validity claims.

Before decision commit, Governance applies the accepted D3-R3 sequence: request still PENDING, request revision current, exact human still eligible, current identity/access gates valid, and material authorization basis VALID through the action-owner boundary when current source meaning is required.

## 3. Why two reads are required

### One-read collection-as-detail alternative — rejected

A collection-only API could be filtered by request ID, but that makes a collection endpoint a pseudo-detail resource, weakens deep-link semantics and leaves no natural resource representation whose ETag is the request concurrency validator.

### Full AuthorizationRequest CRUD/history surface — rejected

Public create/invalidate/reauthorize/general history endpoints would expose D3 owner-to-owner communication without a Product consumer and would move Governance toward a generic workflow/case API.

### Selected — two reads + reanchored existing decision write

This is the smallest surface that simultaneously supports:

```text
Aprovações queue
F13 Notification deep-link
exact request detail
request-local ETag concurrency
human decision
least privilege without source-owner read Permissions
```

## 4. Actionable read representation

The actionable detail is a **purpose-bounded view**, not a generic `AuthorizationRequest` history resource.

Common fields are:

```text
authorization_request_id
typed governed target ref without target ETag
immutable subject_display_label
immutable typed review_basis
requester_or_initiator_principal_id?  // historical/correlation only
predecessor_authorization_request_id? // bounded reauthorization lineage only
created_at
```

`PENDING` and current caller eligibility are admission predicates for these operations, so they are not re-exposed as client-authored or independently writable state.

The detail schema is a closed `oneOf` of four target-specific variants so target kind and review-basis kind cannot mismatch structurally.

## 5. Four typed review-basis families

The request snapshot must be sufficient for an approver who legitimately has `governance.decide` but does not necessarily have `offering.read` or `materialization.read`.

### 5.1 ListingIntent

Snapshot includes the exact ListingIntent identity, source Product reference, listing target, declarative `desired` authoring state and optional `requirements_revision` observed for that authorization episode.

It does not become current Offering/Readiness truth.

### 5.2 PriceIntent

Snapshot includes PriceIntent identity, price target and desired price. It may additionally preserve typed price/economic evidence when that evidence was material to the action owner's `approval-required` disposition:

```text
current price observation (price + observed_at)
expected economics snapshot
minimum contribution-margin policy value
```

Those snapshots remain historical evidence/provenance, not current Economics/Offering authority.

### 5.3 BusinessOrderIntent

Snapshot includes BusinessOrderIntent identity, immutable Sale snapshot, target SourceInstance, PartyResolution and DestinationRealization relevant to the materialization decision.

### 5.4 InvoicingIntent

Snapshot includes InvoicingIntent identity, immutable Sale snapshot, the relevant BusinessOrderIntent snapshot and optional FulfillmentExecution correlation where present.

No generic payload, metadata bag, arbitrary entity map or source DTO is admitted in any review basis.

## 6. AuthorizationDecision feed-forward

`AuthorizationDecision` remains immutable Governance history and now belongs explicitly to one `authorization_request_id`.

Its Product representation preserves the immutable target/review basis needed to explain what was decided so a `governance.read` consumer does not need source-owner read Permission merely to interpret Governance history.

Decision list/history access remains separate:

```text
ListAuthorizationDecisions / GetAuthorizationDecision
→ governance.read
→ H/A/S
```

Actionability remains:

```text
ListMyActionableAuthorizationRequests
GetMyActionableAuthorizationRequest
CreateAuthorizationDecision
→ governance.decide
→ H only
```

No Permission implication is created.

## 7. Notification and Work feed-forward

### F13

`AUTHORIZATION_ACTION_REQUIRED` changes source identity from `AuthorizationTargetRef` to `AuthorizationRequestRef`.

F13 still grants no request/source capability by possession; the actionable GET rechecks current Governance eligibility.

### F14

`AUTHORIZATION_DECISION_RESULT` remains target-oriented under the accepted requester continuation semantics.

### Zero eligible deciders

The existing Work Product operations remain sufficient. The Work origin closed union gains an `authorization_request` variant anchored by `AuthorizationRequestID`; no special Work operation is added.

## 8. Explicit non-operations

The Product surface does not admit:

```text
CreateAuthorizationRequest
InvalidateAuthorizationRequest
ReauthorizeAuthorizationRequest
ListAllAuthorizationRequests
ResolveAuthorizationRequestRecipients
CreateNoApproverWork
```

The corresponding owner-to-owner semantics remain D3/internal owner boundaries unless a future named Product consumer proves otherwise.

## 9. Expected census

If D5-R6 executable wire proof passes:

```text
Product operations: 104 + 2 = 106
ordinary Permissions: 31
Principal kinds: H / A / S only
```

`106` is a consequence of the admitted surface, not a target metric.

## 10. Required proof

D5-R6 must falsifiably prove at minimum:

- exactly 106 Product operations and 31 ordinary Permissions;
- both new reads are H-only + `governance.decide` + ControlledActionGovernance;
- list query surface is only `limit,cursor`;
- actionable detail returns request-local ETag;
- target refs contain no target ETag;
- detail target kind and review-basis kind are structurally paired;
- exactly four closed typed review-basis families and no generic payload/metadata bag;
- `CreateAuthorizationDecision` uses request path + `If-Match` + `Idempotency-Key` and body only `outcome`;
- Decision history carries `authorization_request_id` + immutable review basis;
- F13 points to `AuthorizationRequestRef`; F14 remains target-oriented;
- Work origin admits `authorization_request` without a new Work operation;
- no public create/invalidate/reauthorize AuthorizationRequest operation exists;
- generated OpenAPI projections remain deterministic/valid under the repository full gate.

## 11. Closure sequence

```text
D5-R5 surface ACCEPTED
→ D5-R6 OAD wire + executable proof
→ D6 P9 final Screen Contracts / bidirectional trace
→ independent Fable review + finding adjudication
→ only then Global-Maximum redesign closure / D7-R eligibility
```

**Exact next action:** implement the D5-R6 proof test first, observe RED against canonical 104/31, then implement the smallest OAD correction to GREEN.