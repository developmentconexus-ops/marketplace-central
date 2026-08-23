# NOTIF-01 D1-R2 — Authorization Request Boundary Revalidation

> **Status:** PASS — CURRENT STRUCTURE CONFIRMED
> **Trigger:** operator-ratified supersession of the P9-F1 query-only remedy
> **Parent authority:** [D1 Domains / Boundaries](D1-DOMAINS-BOUNDARIES.md)
> **Upstream product authority:** [D0 Product / System Definition](D0-PRODUCT-SYSTEM-DEFINITION.md)
> **Downstream targeted review:** [D2-R6 Authorization Request Identity & Lifecycle](D6-R2-NOTIF-01-D2-R6-AUTHORIZATION-REQUEST-IDENTITY-LIFECYCLE.md)
> **Canonical Product wire:** unchanged — 104/31
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Question

Does the newly proved pre-decision authorization identity require a new D1 business domain, a moved authority boundary or a new semantic edge?

## 2. Existing D1 authority

`Controlled Action Governance` already owns:

```text
authorization delegation/grant semantics
authorization decision/context
authorized target-scope snapshot
authority context
decision state
provenance/correlation
```

It explicitly does not own:

```text
business thresholds/disposition
business intent
intended target selection
readiness/economic/operational validity
execution
provider protocol
```

D1 also already admits the bidirectional semantic edge:

```text
action-owning domain
⇄ Controlled Action Governance
```

The action owner supplies the domain-owned Business Intent, intended scope, effective action disposition and authorization-relevant context. Governance applies authorization-specific authority/delegation semantics and returns authorization context/decision. The action owner remains responsible for execution-time validity and execution.

## 3. New evidence classification

`AuthorizationRequest` means:

> one concrete pre-decision authorization episode in which Governance has accepted authorization-specific context for a governed target and a consequential authorization decision has not yet occurred.

That meaning is wholly authorization-specific.

It does not create a new independent business responsibility. It does not own the target-domain intent, policy or execution. It is the missing pre-decision lifecycle state inside the already accepted Governance responsibility.

## 4. Alternatives

### A — new Approval/Workflow domain

**Rejected.** The meaning already belongs to Governance. A new domain would duplicate authority and invite generic workflow semantics that Product 1.0 does not require.

### B — push pending approval state back into every action-owning domain

**Rejected.** Action owners own `approval-required` disposition and Business Intent, but current authorization eligibility/context/decision lifecycle belongs to Governance. Duplicating the request lifecycle per owner would create multiple authorities for the same authorization meaning.

### C — keep Governance owner and add the missing canonical pre-decision identity

**Selected.** This preserves one meaning/one authority and uses the existing D1 edge.

## 5. Decision

**CURRENT STRUCTURE CONFIRMED.**

No D1 boundary or semantic edge changes.

```text
Marketplace Offering Operations / Business-System Materialization
→ own Business Intent + disposition + execution-time validity

Controlled Action Governance
→ owns AuthorizationRequest + authority/delegation evaluation + AuthorizationDecision/context

Personal Notifications
→ owns awareness only

Operational Work
→ owns responsibility/work lifecycle only if a separate actionable-work obligation exists
```

## 6. Boundary laws for D2-R6

The D2 repair must preserve:

1. `AuthorizationRequest` is an MPC-owned Governance identity, not a generic platform request.
2. It references a closed governed target; it never becomes a universal entity graph.
3. Its retained review/authorization context is purpose-bounded and cannot become source-owner current truth.
4. Action disposition remains action-owner authority.
5. Decision authority remains Governance authority.
6. Execution-time validity remains action-owner authority after any authorization decision.
7. No request lifecycle state may execute the governed action directly.
8. No ordinary Permission implication is introduced by the identity.

## 7. Reopen triggers

Reopen D1 only if later derivation proves one of these:

- `AuthorizationRequest` needs to own business policy/disposition rather than authorization semantics;
- it becomes a generic workflow/case platform for unrelated business responsibilities;
- it requires a semantic dependency not already present in the accepted action-owner ⇄ Governance edge;
- another domain must independently write the same authorization-request meaning.

None is evidenced now.

## 8. Gate

```text
D1 owner/boundary              CURRENT STRUCTURE CONFIRMED
new business domain            NO
new D1 semantic edge           NO
D2 identity/lineage reopen      REQUIRED / NEXT
D3/D5/OAD                      UNCHANGED / BLOCKED BY D2-R6
```

**Exact next action:** derive and adjudicate only D2-R6 canonical `AuthorizationRequest` identity/lifecycle.