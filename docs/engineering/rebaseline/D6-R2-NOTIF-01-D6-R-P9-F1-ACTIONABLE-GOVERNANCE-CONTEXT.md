# NOTIF-01 D6-R P9-F1 — Actionable Governance Context Finding

> **Status:** OPEN — P9 MATERIAL FALSIFIER / OPERATOR ADJUDICATION REQUIRED
> **Trigger:** P9 exact Screen Contract + bidirectional wire trace after P8 operator `LOCK`
> **Parent frontend authority:** [D6-R Frontend Feed-Forward](D6-R2-NOTIF-01-D6-R-FRONTEND-FEED-FORWARD.md) + [P8 Ratification](D6-R2-NOTIF-01-D6-R-P8-RATIFICATION.md)
> **Product wire:** canonical 104 Product operations · 31 ordinary Permissions · H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Finding

P9 cannot truthfully close `AUTHORIZATION_ACTION_REQUIRED` or the already-admitted `Aprovações` consumer with the current Product wire.

Accepted semantics say:

```text
AUTHORIZATION_ACTION_REQUIRED
→ OWNER_DERIVED exact currently valid decision Principal
→ source_ref = AuthorizationTargetRef
```

Current Notification `AuthorizationTargetRef` carries only the closed target identity:

```text
listing_intent_id
price_intent_id
business_order_intent_id
invoicing_intent_id
```

Current Governance write requires:

```text
CreateAuthorizationDecision
→ governance.decide
→ H only
→ request target = AuthorizationTarget
→ AuthorizationTarget requires exact target ID + current strong ETag
```

The current read surface exposes only completed `AuthorizationDecision` records through `ListAuthorizationDecisions` / `GetAuthorizationDecision`, both under `governance.read`.

There is no Product read for **the current pending authorization context that the current human is actually eligible to decide**.

This also falsifies the pre-existing D6 interaction-map claim:

```text
S110 Aprovações
S111 Decisão / contexto de aprovação
→ target/revision-explicit
→ Get/CreateAuthorizationDecision
```

`GetAuthorizationDecision` requires an `AuthorizationDecisionID`, but under accepted D2 meaning that ID does not exist before the consequential decision occurrence is created.

## 2. Root cause

The missing object is not a Notification payload and not a target-domain permission.

> **Governance owns current decision eligibility and the pending governed-action meaning, but the Product API exposes the decision command without exposing a purpose-bounded current actionable-decision projection from which a human can review the decision context and obtain the exact target revision the command must bind.**

Therefore:

```text
F13 Notification
→ tells P that Governance needs P
→ P may have governance.decide
→ governance.decide does not imply governance.read
→ governance.decide does not imply offering.read/materialization.read
→ Notification does not contain target ETag/current decision context
→ P cannot reliably construct or review CreateAuthorizationDecision
```

W4 explicitly makes all those gates independent. Adding a hidden implication would violate accepted authority.

## 3. Target invariant

A currently eligible human decision Principal must be able to complete the human job without receiving broad unrelated source access:

```text
current Governance eligibility
→ sanctioned current actionable authorization projection
→ enough current decision-relevant presentation to make the decision
→ exact current governed target revision
→ CreateAuthorizationDecision
→ current Governance eligibility + target revision revalidated
```

The projection must never become target business authority or a generic cross-domain payload.

## 4. Credible alternatives

### A — require every approver to also hold target-owner read Permissions

Examples:

```text
governance.decide + offering.read
governance.decide + materialization.read
```

**REJECTED — Local Maximum.**

It couples independent capabilities, broadens source disclosure beyond the exact governed action, contradicts flat exact Permission semantics and makes Governance operability depend on role-bundle accidents.

### B — add ETag and decision details to `AUTHORIZATION_ACTION_REQUIRED` Notification

**REJECTED.**

A Notification is historical awareness, not current decision authority. The snapshot can become stale, would duplicate target/current Governance state and still would not repair the normal `/aprovacoes` queue consumer.

### C — remove target ETag from `CreateAuthorizationDecision` and let the server decide against whatever is current

**REJECTED.**

This weakens an accepted concurrency/correctness property. A human could approve a materially changed target they never reviewed. Server freshness alone is not human review of the exact revision.

### D — add one approval-context read independently inside each target owner

**REJECTED.**

It duplicates Governance eligibility semantics across Offering and Materialization, creates multiple authorities for “can this human decide now?” and forces the frontend to orchestrate owner-specific authorization rules.

### E — one purpose-bounded actionable Governance projection

**SELECTED — GLOBAL MAXIMUM CANDIDATE.**

Add one Human-only Product read, provisionally:

```text
ListMyActionableAuthorizations
```

Semantics:

```text
operation class: read-only projection (P)
frontend clients: H only
ordinary access: governance.decide
scope: exact Organization
projection authority: ControlledActionGovernance for decision eligibility
component truth: governed target remains owned by its source owner
```

The collection lists only authorization contexts the **current human is currently eligible to decide**.

It supports both consumers:

```text
/org/:organizationId/aprovacoes
→ current actionable queue

AUTHORIZATION_ACTION_REQUIRED source_ref
→ same operation with exact closed target filter
→ zero-or-one current actionable context for that human
```

No `AuthorizationRequestID` is invented. Current actionable identity remains the closed governed target identity/revision already required by Governance.

## 5. Candidate Product contract constraints

The exact wire is D5 repair work if this direction is ratified, but the admission contract is already bounded.

### Inputs

```text
Organization path scope
limit / cursor
optional exact closed target filter:
  target_kind = listing_intent | price_intent | business_order_intent | invoicing_intent
  target_id   = exact MPC target ID
```

`target_kind + target_id` is a closed Governance lookup grammar, not a universal entity graph.

### Minimum result semantics

Each item must carry:

```text
current AuthorizationTarget including the current strong ETag
current Governance actionability for the current human
bounded current decision-relevant presentation
```

The presentation must be a **closed typed union by governed target kind** or an equivalently explicit bounded schema. It must contain enough current source-owner-derived information for a human to decide the exact revision without dereferencing the target through an unrelated broad read Permission.

It must not expose:

```text
raw provider/business-system payload
generic metadata/payload
arbitrary target JSON
all source-owner fields
roles / Permission internals
a generic entity graph
```

The source owner remains authority for target meaning and revision. Governance owns only decision eligibility/actionability and the governed-decision contract.

### Access and failure

```text
not current decision Principal
→ absent from collection / exact filtered lookup yields no actionable item

current eligibility revoked while screen open
→ reread/write fails closed

target revision changed
→ old actionable context cannot authorize decision on the new revision

CreateAuthorizationDecision
→ still revalidates current Governance eligibility + exact target revision
```

`governance.decide` does not become `governance.read` and does not grant target-owner read operations.

## 6. Consequence if ratified

The smallest admitted Product-surface consequence is:

```text
104 + 1 = 105 Product operations
31 ordinary Permissions unchanged
Principal kinds H/A/S unchanged
```

The number 105 is a consequence, not a target.

The existing operations remain:

```text
ListAuthorizationDecisions / GetAuthorizationDecision
→ completed-decision read/history under governance.read

CreateAuthorizationDecision
→ consequential decision under governance.decide

ListMyActionableAuthorizations
→ current actionable human decision context under governance.decide
```

## 7. P9 feed-forward after repair

If the owner-level repair is accepted and proved, P9 resumes with:

```text
F13 AUTHORIZATION_ACTION_REQUIRED
→ navigate by AuthorizationTargetRef, never a not-yet-existing AuthorizationDecisionID
→ resolve exact current actionable context via ListMyActionableAuthorizations exact target filter
→ render current bounded decision context
→ CreateAuthorizationDecision using the returned current AuthorizationTarget + ETag

Aprovações queue
→ ListMyActionableAuthorizations
→ completed/history view may independently use ListAuthorizationDecisions where governance.read exists
```

F14 `AUTHORIZATION_DECISION_RESULT` remains unchanged: requester continuation uses `AuthorizationTargetRef`, not Governance history access.

## 8. Falsifiers / negative controls

This repair fails if it:

- makes `governance.decide` imply `governance.read`, `offering.read`, `materialization.read` or target write authority;
- stores target current truth inside Notifications;
- removes target revision binding from the decision command;
- invents `AuthorizationRequestID` merely for routing symmetry;
- creates one approval read API per target owner;
- introduces generic `payload`, `metadata`, `entity_type/entity_id` or arbitrary JSON;
- exposes non-actionable authorizations to the current human through this purpose-bound operation;
- lets candidate/actionable presence alone authorize a later write without fresh server revalidation;
- rewrites F14 requester continuation back to `AuthorizationDecisionID`.

## 9. Decision

**Recommendation:** `RESTRUCTURE NOW` through the one purpose-bounded Governance actionable projection before final P9 Screen Contract freeze.

No P8 wireframe needs to be reopened. The finding is a backend/Product operability gap exposed by P9 traceability, not a visual-layout defect.

## 10. Gate

```text
P8 B00-R2 / B11 / B12             OPERATOR-LOCKED
P9 Screen Contracts               OPEN / BLOCKED BY P9-F1
P9-F1 Governance projection       OPERATOR ADJUDICATION REQUIRED
canonical Product OAD             104/31 UNCHANGED UNTIL RATIFICATION
D7-R / D8-R                       BLOCKED
Product implementation            BLOCKED UNTIL D9
```

**Exact next action:** operator adjudicates only the P9-F1 Global-Maximum direction. Do not edit the canonical OAD, freeze P9, begin D7-R/D8-R or implement Product code before this gate closes.