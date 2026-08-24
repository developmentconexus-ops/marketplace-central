# NOTIF-01 D5-R6 — AuthorizationRequest OAD Wire Proof

> **Status:** PROVED / CANONICAL
> **Ratified Product surface:** [D5-R5 AuthorizationRequest Product Surface](D6-R2-NOTIF-01-D5-R5-AUTHORIZATION-REQUEST-PRODUCT-SURFACE.md)
> **Identity authority:** [D2-R6 Ratification](D6-R2-NOTIF-01-D2-R6-RATIFICATION.md)
> **Communication authority:** [D3-R3 Ratification](D6-R2-NOTIF-01-D3-R3-RATIFICATION.md)
> **Canonical wire:** `contracts/api/product/openapi.yaml`
> **Result:** 106 Product operations · 31 ordinary Permissions · H/A/S only
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Proof claim

D5-R6 proves that the operator-ratified `AuthorizationRequest` Product surface can be expressed by the canonical Product OAD without widening ordinary access, exposing owner-to-owner workflow capabilities, coupling authorization validity to a whole-target ETag, or regressing the accepted historical Product generations.

The final Product delta is exactly:

```text
+ ListMyActionableAuthorizationRequests
+ GetMyActionableAuthorizationRequest
~ CreateAuthorizationDecision  // same operation identity, request-anchored wire

104 → 106 Product operations
31  → 31 ordinary Permissions
H/A/S unchanged
```

The count is a consequence of the admitted semantics, not an architecture target.

## 2. Canonical actionable-read wire

### `ListMyActionableAuthorizationRequests`

```text
GET /organizations/{organization_id}/authorization-requests
owner       ControlledActionGovernance
class       Q
Permission  governance.decide
Principal   H only
query       limit, cursor only
```

It returns only currently `PENDING` requests the exact current human caller is currently eligible to decide in that Organization. It does not grant Governance history access and does not expose search, generic lifecycle filtering, role/Permission filters, assignee filters, totals or saved views.

### `GetMyActionableAuthorizationRequest`

```text
GET /organizations/{organization_id}/authorization-requests/{authorization_request_id}
owner       ControlledActionGovernance
class       Q
Permission  governance.decide
Principal   H only
```

Success is current-actionability scoped. The response strong ETag represents the `AuthorizationRequest` owner-local concurrency revision. It is not the governed target resource ETag and is not the material authorization-validity oracle.

Possession of `AuthorizationRequestID` or a Notification remains non-authoritative.

## 3. Canonical decision wire

`CreateAuthorizationDecision` is reanchored to the canonical request:

```text
POST /organizations/{organization_id}/authorization-requests/{authorization_request_id}:decide
owner       ControlledActionGovernance
class       C
Permission  governance.decide
Principal   H only
headers     If-Match + Idempotency-Key
body        { outcome: authorize | reject }
```

The two request-trust carriers protect different properties:

```text
If-Match        → stale/concurrent decision protection on AuthorizationRequest
Idempotency-Key → ambiguous acceptance / safe retry of the same human decision command
```

The client no longer supplies a governed target or target ETag. Before commit, Governance performs the accepted D3-R3 current request, current exact-human eligibility, IdentityAccess and material-basis validity revalidation.

Wire outcomes preserve fail-honest distinctions:

- `412` / `428` — request precondition stale/missing;
- `409` — current request/state/idempotent intake conflict;
- `422` — validation / idempotency-key material mismatch;
- `503` — current material authorization-basis validity cannot be sufficiently established now; no Decision is fabricated and the request remains pending.

## 4. Typed authorization review basis

The current OAD contains exactly four closed typed review-basis families:

```text
listing_intent
price_intent
business_order_intent
invoicing_intent
```

The actionable detail is a closed four-variant `oneOf` pairing the governed target kind with the same review-basis kind structurally. A target/review mismatch cannot be represented merely by convention.

### ListingIntent basis

Preserves the authorization-purpose ListingIntent snapshot:

```text
listing_intent_id
source_product
target
desired
requirements_revision?  // observed provenance when present
```

### PriceIntent basis

Preserves:

```text
price_intent_id
target
desired_price
current_price_observation?
expected_economics?
minimum_contribution_margin_rate?
```

Optional current-price/economic/policy values are historical authorization evidence only when material to the accepted episode; they do not become current Offering/Economics authority.

### BusinessOrderIntent basis

Preserves:

```text
business_order_intent_id
sale_snapshot
target_source_instance_id
party_resolution
destination_realization
```

### InvoicingIntent basis

Preserves:

```text
invoicing_intent_id
sale_snapshot
business_order_snapshot
fulfillment_execution_id?
```

No review basis admits generic `payload`, `metadata`, arbitrary fields or an entity graph.

## 5. AuthorizationDecision history

`AuthorizationDecision` remains immutable Governance history and now explicitly carries:

```text
authorization_decision_id
authorization_request_id
typed governed target ref without target ETag
immutable typed review_basis
outcome
deciding Principal
decided_at
```

The decision representation structurally preserves the same target/review-basis pairing used by the request. Historical Governance reads remain `governance.read` and H/A/S; actionable request reads/write remain `governance.decide` and H-only. No Permission implication is introduced.

## 6. Notification feed-forward

The current Notification union still has exactly 14 kinds.

F13 now uses:

```text
AUTHORIZATION_ACTION_REQUIRED
→ AuthorizationRequestRef
```

This enables exact request continuation while preserving the rule that Notification possession grants no capability; the actionable Governance GET rechecks current eligibility.

F14 remains target-oriented:

```text
AUTHORIZATION_DECISION_RESULT
→ AuthorizationTargetRef
```

so the requester can continue to the original work without requiring Governance history access.

The existing NOTIF laws remain intact: self Inbox has no `notifications.read`, route management uses only `notifications.manage`, recipient discovery remains minimal, and read/archive state remains Notification-owned.

## 7. Zero-decider Work feed-forward

No new Work operation was added.

The existing Work Product surface is reused and the closed `WorkOrigin` union gains:

```text
kind: authorization_request
authorization_request_id
```

`ListWork.origin_kind` admits `authorization_request` accordingly.

This is the Product representation of the D3-R3 law:

```text
PENDING request + known-empty current eligible-human set
→ explicit Governance actionable condition
→ Operational Work responsibility lifecycle
```

Work never becomes authorization authority and no fallback approver is fabricated.

## 8. Explicitly absent public workflow surface

The proved OAD does not admit:

```text
CreateAuthorizationRequest
InvalidateAuthorizationRequest
ReauthorizeAuthorizationRequest
ListAllAuthorizationRequests
ResolveAuthorizationRequestRecipients
CreateNoApproverWork
```

Action-owner ↔ Governance request intake, invalidation and recovery remain owner-to-owner D3 semantics until a real Product consumer proves otherwise.

## 9. RED → GREEN evidence

### RED

At `735f36f014fa15d64469bc3f6f0c6b68be21ca07`, the new verifier was installed before the wire change.

The historical/current preconditions passed through the accepted 104-operation generation and the new verifier failed exactly on:

```text
AuthorizationRequest D5-R6 Product operation count must be 106, found 104
```

This proves the new control could fire before implementation.

### GREEN iterations

The first wire candidate exposed two proof-infrastructure issues before semantic closure:

1. historical 99-generation proof roots copied the current OpenAPI but did not rewind the new overlay refs; the harness was corrected to rewind the D5-R6 delta while preserving the historical 95/29 and 99/30 contracts;
2. Redocly detected a local component-name collision with the prior Notification-module `AuthorizationTargetRef`; the new overlay-local schema was renamed while the canonical root alias remains `AuthorizationTargetRef`.

Neither correction weakened the D5-R6 verifier.

The next current-state failure was the older NOTIF verifier's literal 104 census. Only that superseded current census was changed to 106; its 8 negative controls were preserved unchanged.

## 10. Final executable proof

At semantic GREEN HEAD `5c1133fc034ec940e599f8a5e9f2495bd6f5a117`, CI #590 completed successfully.

Observed proof:

```text
product_oad_baseline_non_regression=PASS
product_oad_operations=95/95
product_oad_permissions=29/29

product_oad_operations=99/99
product_oad_permissions=30/30
product_oad_historical_99_non_regression=PASS

product_oad_current_generated_projection_semantics=PASS
product_oad_operations=106/106
product_oad_auth_negative_controls=5/5

authorization_request_oad_operations=106/106
authorization_request_oad_permissions=31/31
authorization_request_oad_review_basis=4/4
authorization_request_oad_negative_controls=10/10
authorization_request_oad=PASS

notification_oad_operations=106/106
notification_oad_permissions=31/31
notification_oad_negative_controls=8/8
notification_oad=PASS

operational_read_contract_proof=PASS
B00-R2/B11/B12 structural proofs=PASS
repository full gate=PASS
```

Generated TypeScript and Go projections are deterministic and compile/test under the existing proof harness.

## 11. Canonical result

```text
Product operations       106
ordinary Permissions     31
Principal kinds          H / A / S
stable origin            https://conexus.fun
active runtime baseline  NONE
implementation           BLOCKED UNTIL accepted D9
```

D5-R5 + D5-R6 now supersede the prior 104-operation current wire only for the bounded AuthorizationRequest repair. Historical 95/29, 99/30 and the prior 104/31 generation remain proof evidence, not current Product wire authority.

## 12. Required continuation

```text
D5-R6 PROVED / CANONICAL
→ D6 P9 final Screen Contracts + bidirectional wire trace
→ independent Fable review of complete AuthorizationRequest redesign
→ adjudicate every material finding
→ only then may redesign be declared Global-Maximum closed / D7-R become eligible
```

P9 must prove at minimum:

- `/aprovações` actionable queue ↔ `ListMyActionableAuthorizationRequests`;
- exact request detail/F13 continuation ↔ `GetMyActionableAuthorizationRequest`;
- decision ↔ request ETag + Idempotency-Key + `CreateAuthorizationDecision`;
- stale / no-longer-actionable / current-validity-unavailable UX;
- no source-owner read Permission workaround;
- Governance history remains separate from actionable self reads;
- zero-decider request is represented through Work, not the user's actionable approvals queue.

**Exact next action:** resume only D6 P9 final Screen Contracts / bidirectional trace against canonical 106/31. Keep B10, D7-R, D8-R and Product implementation blocked.