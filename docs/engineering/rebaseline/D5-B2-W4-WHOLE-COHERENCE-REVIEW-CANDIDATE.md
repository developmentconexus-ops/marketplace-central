# D5-B2 — Whole-W4 Permission / Client-Class Coherence Review Candidate

> **Status:** NON-AUTHORITATIVE LEAD REVIEW CANDIDATE  
> **Review subject:** accepted-in-stage W4 as one Product access-enforcement system  
> **Authority:** none — review evidence only until operator ratification and canonical filing  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Reviewed:** 2026-08-19

## 1. Review outcome

**Method outcome:** `RESTRUCTURE NOW — W4-LOCAL`.

The accepted W4 direction is globally sound. No current finding requires a D0→W3 semantic parent reopen. Two small but material W4-local gaps should be corrected before W4 can be accepted as a whole.

## 2. Full-sweep results that passed

The lead review re-derived the admitted operation inventory and confronted it against W4:

```text
admitted Product operations                           95
W4 operation enforcement rows                        95
unmapped admitted operations                           0
additional Product operations introduced               0
stored ordinary Permissions                            29
Permission hierarchy/wildcard relations                 0
```

The following survived adversarial challenge:

- B2 `both` refinement: ordinary Q/read operations → `H|A|S`; consequential business-authoring/coordination C → `H|A`; side-effect-free `EvaluatePriceScenario` → `H|A|S`;
- `system` is not a substitute for `automation` on business-authoring operations;
- all currently human-only operations remain justified by accepted operation admission; no proven current machine consumer is blocked;
- `RecordSeparation` and `RecordPacking` remain human-only baseline;
- `RecordPhysicalConference` and `RecordDispatchHandoff` remain human OR explicitly qualified `system`;
- 29 Permissions remain flat/exact; AccessRole may bundle but no `manage ⇒ read`, `execute ⇒ read`, prefix or wildcard semantics emerge;
- `listing.manage != price.manage` remains essential;
- `fulfillment.read != fulfillment.execute != fulfillment.manage` remains essential;
- Economics and Governance read/manage/decision splits remain proportionate;
- `offering.read` remains sufficient for Listing/ListingIntent/PriceIntent reads; no `price.read` by symmetry;
- `sales.manage` and `post_sale.manage` are linguistically broad but create no wildcard/future-operation authority and do not justify churn-only renaming;
- `fulfillment.execute` protecting FulfillmentArtifact reads remains proportionate because those reads are operational/PII-sensitive; Principal-class restrictions still prevent automation/system from establishing human-only physical facts;
- `GetCurrentAccessContext` remains the only authenticated/no-Organization discovery exception;
- Organization no-Membership → privacy-preserving 404, versus current-Membership but missing Permission/client-class admission → 403, remains coherent with W1/W2 cross-Organization privacy;
- business `approval-required`, rejection/prohibition, unavailable/unknown and Governance outcomes remain separate from ordinary access failure;
- current MPC Membership/RoleAssignment/Permission remains access authority; IdP roles/scopes/provider roles/frontend guards remain non-authoritative;
- monotonic AccessRole/AuthorizationDelegation revocations retain ordinary caller admission and only relax stale-target concurrency semantics;
- Structural Inversion against current Keycloak roles/scopes, middleware and frontend visibility: **PASS**.

---

## 3. W4-G1 — current Principal access eligibility / binding failure grammar — REVISE

### Evidence / gap

B2-A already requires:

- fail-closed resolution to one MPC Principal;
- revocation/disablement to stop future access without rewriting history.

W4 currently says:

```text
authenticate
→ resolve exactly one Principal
→ Membership
```

but does not explicitly require that the resolved Principal is **currently eligible for Product access**. A disabled/revoked Principal could therefore appear to pass W4 if its historical Membership/RoleAssignments still exist.

The failure grammar also does not distinguish:

1. a valid-looking external token that cannot resolve uniquely to any accepted MPC Principal; versus
2. a successfully resolved MPC Principal whose current identity/access lifecycle disables further Product access.

### Root cause

Authentication binding success and current MPC Principal access eligibility were collapsed into one informal “resolve Principal” step.

### Corrected invariant / direction

The W4 enforcement sequence becomes proportionately:

```text
1. authenticate / accept external credential
2. resolve exactly one MPC Principal
3. require current MPC Principal access eligibility under D2 identity/access state
4. identify Product operation / apply access-context exception
5. resolve current Organization Membership when Organization-scoped
...
```

Failure classes:

- missing/invalid/untrusted/wrong-audience credential, or credential that cannot resolve **uniquely** to an accepted MPC Principal → `401 authentication-required`;
- successfully resolved MPC Principal that is currently disabled/not access-eligible → `403 access-denied`;
- access-eligible Principal with no current Membership in the path Organization → privacy-preserving `404 resource-not-found`;
- then ordinary Permission/client-class rules continue as already accepted.

W4 does **not** invent a new Principal lifecycle enum/resource. D2 identity/access remains owner of Principal lifecycle; W4 only consumes current eligibility.

An impossible duplicate external binding that resolves one external identity to multiple Principals is a violated D2 invariant/server integrity fault, not a new public access-state taxonomy.

### Alternatives

1. treat successful binding as perpetual Principal eligibility — rejected: defeats accepted disablement semantics;
2. encode enable/disable only in IdP/token — rejected: duplicates MPC Principal lifecycle and cannot cover MPC-owned disablement;
3. **consume current D2 Principal access eligibility after binding — selected Global Maximum.**

### Parent reopen

None. W4-local enforcement completion of already-accepted D2/B2-A semantics.

---

## 4. W4-G2 — physical qualification ordering and non-self-assertion — REVISE

### Evidence / gap

W4 correctly separates `system` Principal kind from operation-specific physical-establisher qualification, but the current enforcement order groups:

```text
enforce allowed Principal kind / special qualification
→ enforce exact ordinary Permission
```

The physical qualification is not merely Principal kind. It is a sensitive Fulfillment-specific epistemic predicate and may depend on server-known Principal/source provisioning. Evaluating it before ordinary Permission is unnecessary and risks treating qualification as a parallel ordinary-access authority.

The accepted W2 physical evidence rule also says server attributes effective Principal/source/time and clients cannot author trusted-physical-evidence claims. W4 should make the same fence explicit for system qualification.

### Corrected invariant / direction

Split broad Principal-kind admission from additional physical qualification:

```text
current Membership
→ allowed Principal kind
→ exact ordinary Permission
→ if operation admits qualified system: server-resolved operation-specific physical qualification
→ same-Organization/resource resolution
→ W1/W2/domain/Governance gates
```

Rules:

- `H` on an admitted physical operation uses the existing accountable-human baseline after `fulfillment.execute`;
- `S` is admitted for PhysicalConference/DispatchHandoff only after `fulfillment.execute` **and** server-resolved qualification for that exact checkpoint capability;
- `A` remains ineligible for physical checkpoint establishment;
- caller request body, token role/scope or arbitrary claim cannot self-assert `physical_system`, `trusted`, `qualified`, evidence-source authority or equivalent;
- successful system-established physical checkpoint preserves the actual effective Principal/source/time provenance already required by W2;
- qualification remains Fulfillment-operation-specific; do not promote it to Permission, AccessRole, Principal kind or generic `system_capabilities[]` graph.

### Alternatives

1. physical qualification as another Permission — rejected: conflates ordinary access with epistemic authority;
2. physical qualification as a fourth Principal kind — rejected: duplicates D2 taxonomy and overgeneralizes one bounded need;
3. evaluate qualification before Permission — works functionally but weakens authority separation and may perform sensitive qualification resolution for callers that lack the operation Permission;
4. **kind → exact Permission → operation-specific server qualification — selected Global Maximum.**

### Parent reopen

None. W4-local ordering/fence correction.

---

## 5. Strong challenges that survived unchanged

### 5.1 H/A/S on `both` reads

**KEEP.**

B2-A already admits human and non-human Product clients; W4's `H|A|S` on ordinary `both` Qs grants no actual data access without current Membership + explicitly assigned exact Permission. Narrowing every read to `H|A` would add operation-specific exclusions without a correctness property while blocking bounded system consumers such as fulfillment workstations/label systems from reused reads.

This does not make `system` business automation because consequential business-authoring C operations remain H/A or narrower.

### 5.2 `fulfillment.execute` for artifact reads and physical capabilities

**KEEP.**

The shared Permission is deliberate least privilege from the ratified operation matrix. Client-class/physical qualification prevents the read capability from silently authorizing system/automation establishment of physical facts. A separate `fulfillment.artifact.read` would add Permission surface without a proven consumer needing that split.

### 5.3 404 versus 403 privacy split

**KEEP.**

No-Membership path-Organization access uses 404 to avoid confirming another Organization. Once current Membership exists, exact Permission or client-class failure may honestly use 403. Secondary-reference privacy remains W1/W2 fail-closed.

### 5.4 Flat Permissions

**KEEP.**

Flat exact Permissions make role bundles explicit and keep Product operation admission reviewable. Implicit inheritance/wildcards would create a second permission language and future-operation authority by naming convention.

### 5.5 `EvaluatePriceScenario` under `economics.read`

**KEEP.**

The operation is side-effect-free/stateless and was explicitly ratified as non-consequential. A separate execute Permission would be permission-by-HTTP-verb rather than business meaning.

---

## 6. Global Maximum after W4-G1/G2

```text
external AuthN
→ unique MPC Principal binding
→ current Principal access eligibility
→ current Organization Membership
→ allowed Principal kind
→ flat exact Permission
→ bounded operation-specific physical qualification only where required
→ resource/reference privacy
→ W1/W2 safety
→ owner business disposition
→ Governance when applicable
→ owner intake/effect
```

while preserving:

```text
no IdP-role Product authority
no OAuth-scope Permission duplication
no Permission hierarchy/wildcards
no generic IAM/ReBAC DSL
no fourth physical-system Principal kind
no generic machine-capability graph
no Governance widening of client class
no system-as-business-automation shortcut
```

## 7. Reopen classification

- D0 product boundary: **NO REOPEN**
- D1 semantic owners: **NO REOPEN**
- D2 identity/access: **NO REOPEN**; W4 consumes existing Principal lifecycle/eligibility
- D3 communication: **NO REOPEN**
- D4/D4-R1: **NO REOPEN**
- D5-B1/B2 operation inventory/W1/W2/W3: **NO REOPEN**
- W4: **targeted local corrections G1/G2 required before final acceptance**

## 8. Operator decision requested

Ratify/revise W4-G1 and W4-G2 as the Whole-W4 lead correction direction.

Until operator decision:

- accepted W4 artifact remains current in-stage authority;
- this candidate is evidence only;
- do not silently apply G1/G2 to W4;
- do not begin technical ingress classification or later Wire obligations.
