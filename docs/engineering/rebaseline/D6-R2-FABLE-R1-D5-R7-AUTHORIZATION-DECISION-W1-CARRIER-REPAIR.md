# D6-R2 Fable R-1 / D5-R7 — AuthorizationDecision W1 Carrier Repair

> **Status:** OPERATOR-RATIFIED / ACCEPTED
> **Operator adjudication:** 2026-08-24 — Fable R-1 accepted; typed-request-ETag alternative selected
> **Trigger:** independent methodology + whole-repository Fable review `review/d6r2-methodology-whole-repo-fable @ 975a9176cb4019960166e436d63e33349819c046`
> **Canonical grammar owner:** [D5-B2 W1 Wire Contract](D5-B2-WIRE-CONTRACT.md)
> **Prior AuthorizationRequest wire proof:** [D5-R6](D6-R2-NOTIF-01-D5-R6-AUTHORIZATION-REQUEST-OAD-WIRE-PROOF.md)
> **Affected composed authorities:** [final P9](D6-R2-P9-AUTHORIZATION-REQUEST-BIDIRECTIONAL-SCREEN-CONTRACTS.md) · [D7-R](D6-R2-NOTIF-01-D7-R-AUTHORIZATION-REQUEST-RUNTIME-REPAIR.md) · [D8-R](D6-R2-NOTIF-01-D8-R-AUTHORIZATION-REQUEST-GOLDEN-FLOW-REVALIDATION.md)
> **Canonical Product:** 106 Product operations · 31 ordinary Permissions · H/A/S only
> **Implementation:** BLOCKED UNTIL accepted D9

```text
D5R7_W1_CARRIER:TYPED_REQUEST_ETAG
D5R7_FAILURES:MISSING_INVALID_422_STALE_409
D5R7_REPLAY_ORDER:IDEMPOTENCY_BEFORE_REVISION_PRECONDITION
D5R7_B110:REVALIDATED_STRUCTURE_UNAFFECTED
D5R7_SUPERSEDES:D5R6_P9_D7R_D8R_CARRIER_ONLY
D5R7_NEW_PRODUCT_OPERATIONS:0
D5R7_NEW_PERMISSIONS:0
D5R7_NEW_BUSINESS_MEANING:0
```

## 1. Finding and adjudication

The accepted W1 grammar and the later AuthorizationRequest wire had one material carrier contradiction.

W1 §14.2 says an owner custom-method request target such as:

```text
/resource/{id}:verb
```

is not the literal protected base-resource request target and therefore does **not** carry the base resource validator in HTTP `If-Match` unless separately ratified alias semantics exist. MPC created no such alias. The same strong opaque owner validator must instead travel as typed technical request data; missing/invalid typed ETag uses `422`, and stale typed ETag uses `409 resource-revision-conflict`.

D5-R6 later admitted:

```text
POST /organizations/{organization_id}/authorization-requests/{authorization_request_id}:decide
If-Match: <AuthorizationRequest ETag>
Idempotency-Key: ...
{ outcome }
```

with `428/412`. Final P9, D7-R and D8-R then composed that carrier. The semantic model remained coherent, but two ratified wire authorities could not both be implemented literally.

The operator accepts the Fable finding and selects the smallest repair: **preserve W1 universally and correct only the later `:decide` carrier**. No W1 alias exception is created.

## 2. Canonical current `CreateAuthorizationDecision` wire

The operation identity, URI, semantic owner, Permission and Principal-kind admission are unchanged:

```text
POST /organizations/{organization_id}/authorization-requests/{authorization_request_id}:decide
owner       ControlledActionGovernance
class       C
Permission  governance.decide
Principal   H only
header      Idempotency-Key
body        { etag, outcome }
```

`etag` is the exact strong opaque AuthorizationRequest owner-local validator returned by the actionable request detail. It remains concurrency proof only; it is not governed-target validity authority.

Canonical request body:

```json
{
  "etag": "\"opaque-authorization-request-validator\"",
  "outcome": "authorize"
}
```

The frontend still does **not** supply:

```text
governed target
target ETag
current material-validity claim
current eligibility claim
requester identity
```

## 3. Failure grammar

The current custom-method grammar is:

```text
missing/invalid typed etag      → 422 validation-error
typed etag stale                → 409 resource-revision-conflict
current state/intake conflict   → 409 accepted conflict grammar
idempotency material mismatch   → 422 accepted validation grammar
exact semantic validity unknown → exact typed 503 authorization-validity-unavailable
```

`412` and `428` are no longer admitted on `CreateAuthorizationDecision`; they remain valid W1 responses for standard methods whose literal HTTP request target is the protected resource.

The exact typed semantic 503 contract is unchanged: only the specific `authorization-validity-unavailable` Product Problem is known-no-effect; other uncertain/non-semantic 503 remains ambiguous potentially accepted.

## 4. Idempotency / D7-R carrier-neutral composition

The accepted D7-R namespace remains unchanged:

```text
Organization
+ effective Principal
+ CreateAuthorizationDecision operation identity
+ Idempotency-Key digest
```

Within that namespace the semantic fingerprint is now mechanically:

```text
authorization_request_id
+ typed request etag
+ outcome
```

instead of request + `If-Match` + outcome.

The binding ordering law is carrier-neutral:

```text
exact committed idempotent replay
→ BEFORE current revision-precondition evaluation
```

Therefore a lost `201` followed by the exact same Principal/key/request/etag/outcome replays the original committed Decision even though that supplied request revision is now historically stale. A genuinely new attempt evaluates the current request and typed revision proof normally.

No current eligibility, material-validity, F13/F14, zero-decider Work, invalidation-recovery or PITR law changes.

## 5. P9 / B110 bounded revalidation

The operator-LOCKED B110 structural experience remains valid:

```text
queue/detail/history hierarchy          UNAFFECTED
inline approve/reject confirmation      UNAFFECTED
typed review basis                      UNAFFECTED
Idempotency-Key attempt lifecycle       UNAFFECTED
stale-decision recovery state           UNAFFECTED
503 validity-unavailable state          UNAFFECTED
source continuation/current reauth      UNAFFECTED
```

Only the technical write annotation changes:

```text
old historical P8/P9 evidence: If-Match + outcome
current authority:             body.etag + outcome
```

and stale request revision is now represented through accepted `409` revision-conflict grammar instead of `412`.

The reviewed P8 HTML remains immutable historical evidence; it is not rewritten merely to alter a transport carrier that does not change the operator interaction, hierarchy, state placement or responsive behavior.

```text
B110 LOCK disposition = REVALIDATE
result                = STRUCTURE UNAFFECTED
P8 reopen             = NO
```

### 5.1 LOCK impact sweep for this material upstream correction

| LOCKED block | Disposition | Reason |
| --- | --- | --- |
| B00 App Shell / IA | UNAFFECTED | no route/context/shell meaning changed |
| B01 Overview | UNAFFECTED | no operation/state used by Overview changed |
| B00-R2 Notification utility | UNAFFECTED | Notification awareness wire unchanged |
| B11 Personal Inbox | UNAFFECTED | Personal Notification read/write wire unchanged |
| B12 Notification Routing Settings | UNAFFECTED | route-management wire unchanged |
| B110 Approvals | REVALIDATE → STRUCTURE UNAFFECTED | only `CreateAuthorizationDecision` revision carrier/status code changed; human interaction is unchanged |

## 6. D8-R bounded disposition

No Golden Flow set or business choreography changes.

GF-01/GF-02 still exercise the same four AuthorizationRequest review-basis kinds. The only mechanical replacement is:

```text
same Principal/key exact retry
→ replay before revision-precondition evaluation
```

rather than the historical carrier-specific phrase “before stale If-Match evaluation”. SR-01 continuity/PITR law is unchanged. GF-03 remains not materially affected. P1–P6 live probes are not reopened.

## 7. Exact supersession boundary

This document supersedes **only carrier-specific current meaning** in the following accepted historical artifacts:

- D5-R6 references to `If-Match` + `412/428` for `CreateAuthorizationDecision`;
- final P9 references/trace to AuthorizationRequest `If-Match` and `412` stale decision for that operation;
- D7-R references to `REQUEST_IFMATCH_OUTCOME` and `IDEMPOTENCY_BEFORE_IF_MATCH` for that operation;
- D8-R references to exact replay before stale `If-Match` evaluation.

Everything else in those artifacts remains accepted and is consumed unchanged.

Historical proof text is intentionally not rewritten. Current consumers must compose those artifacts with D5-R7.

## 8. Accepted result

```text
AuthorizationRequest model              UNCHANGED
AuthorizationDecision model             UNCHANGED
CreateAuthorizationDecision operation   UNCHANGED
Product operations                      106
ordinary Permissions                    31
Principal kinds                         H / A / S
request revision authority              AuthorizationRequest strong opaque ETag
custom-method revision carrier          typed body.etag
missing/invalid revision proof          422
stale revision proof                    409 resource-revision-conflict
Idempotency-Key                         REQUIRED
exact committed replay ordering         BEFORE REVISION PRECONDITION
B110 structure                          REVALIDATED / UNAFFECTED
Golden Flow set                         UNCHANGED
new Product/runtime architecture        NONE
```

The R-1 material contradiction is closed by this bounded D5 wire repair. Product implementation remains blocked until accepted D9.
