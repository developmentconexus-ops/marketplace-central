# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **W1 + W2 + W3 CANONICAL; W4 ACCEPTED IN-STAGE; Whole-W4 independent review + GPT final adjudication CONVERGED / RESTRUCTURE W4-LOCAL; operator final Whole-W4 ratification = NEXT**  
> **Implementation:** **BLOCKED until D9 is accepted**  
> **Last updated:** 2026-08-19

## 1. Authority path

A fresh session reads, in order:

1. `AGENTS.md`
2. this router
3. `docs/engineering/standards/root-cause-global-maximum-method.md`
4. `ARCHITECTURE.md`
5. `docs/engineering/rebaseline/DECISION-RECONCILIATION-BASELINE.md`
6. `docs/architecture/decisions/README.md`
7. `docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md`
8. `docs/engineering/rebaseline/D1-DOMAINS-BOUNDARIES.md`
9. `docs/engineering/rebaseline/D2-IDENTITY-TENANT-DATA-OWNERSHIP.md`
10. `docs/engineering/rebaseline/D3-COMMUNICATION-EVENTS.md`
11. `docs/engineering/rebaseline/D4-EXTERNAL-INTEGRATIONS.md`
12. `docs/engineering/rebaseline/D4-R1-PUBLICATION-INPUT.md`
13. `docs/engineering/rebaseline/D5-API.md`
14. `docs/engineering/rebaseline/D5-B2-PRODUCT-OPERATION-SURFACE.md`
15. `docs/engineering/rebaseline/D5-B2-OPERATION-ADMISSION-MATRIX.md`
16. `docs/engineering/rebaseline/D5-B2-WIRE-CONTRACT.md` — canonical W1
17. `docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md` — canonical W2
18. `docs/engineering/rebaseline/D5-B2-W3-COLLECTION-GRAMMAR.md` — canonical W3
19. `docs/engineering/rebaseline/D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT.md` — accepted W4 authority/design home
20. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
21. code/OpenAPI/schemas/tests/runtime only as evidence when needed

This file alone owns **program status, allowed/blocked work and exact next action**. Detailed semantics remain in the accepted artifacts above.

`D5-B2-W4-WHOLE-COHERENCE-REVIEW-CANDIDATE.md` and `AI-DIALOG.md` are **NON-AUTHORITATIVE review evidence**. The converged Whole-W4 corrections below do not amend D2/W4 until final operator ratification and canonical filing.

`docs/engineering/rebaseline/cockpit.html` remains a non-authoritative visual projection and is synchronized only after canonical status changes.

Legacy/current code, OpenAPI, IdP roles/scopes, middleware, provider roles/scopes and frontend guards remain evidence only, never Product access authority.

---

## 2. Program state

```text
D0 — Product / System Definition                         CLOSED / ACCEPTED
D1 — Domains / Boundaries                                CLOSED / ACCEPTED
D2 — Identity / Tenant / Data Ownership                  CLOSED / ACCEPTED
D3 — Communication / Events                              CLOSED / ACCEPTED
D4 — External Integrations                               CLOSED / ACCEPTED AS A WHOLE
D4-R1 — Publication Input & Listing Authoring            ACCEPTED / CANONICAL
Decision Reconciliation                                  ACCEPTED / CANONICAL
D5 — API                                                  OPEN / ACTIVE
  B1 Semantic API Model                                  ACCEPTED / CANONICAL
  B2 Product Operation / Resource Surface                 OPEN / ACTIVE
    B2-A Client/Auth                                     ACCEPTED IN-STAGE
    Operation Admission Matrix                           ACCEPTED / RATIFIED
    Whole-Matrix Global Coherence                        ACCEPTED / RATIFIED
    Wire Contract
      W1 Resource / Path / HTTP Grammar                  ACCEPTED / CANONICAL
      W2 Request / Response Schema Grammar               ACCEPTED / CANONICAL
      W3 Collections / Query / Cursor Grammar            ACCEPTED / CANONICAL
      W4 Permission / Client-Class Enforcement           ACCEPTED IN-STAGE
        Whole-W4 lead review                             COMPLETE / RESTRUCTURE W4-LOCAL
        Fable Whole-W4 independent review                COMPLETE
        GPT final adjudication                           CONVERGED
        operator final Whole-W4 ratification             NEXT
      Technical non-Product ingress classification       BLOCKED BY W4
      Final Problem/media consistency                    BLOCKED BY SEQUENCE
      Single OpenAPI wire authority/tooling decision     BLOCKED BY SEQUENCE
D6 — Frontend                                             BLOCKED BY D5
D7 — Runtime / Jobs / Transactions                        BLOCKED
D8 — Golden Flows                                         BLOCKED
D9 — Adversarial Architecture Review                      BLOCKED
Implementation                                            BLOCKED UNTIL D9
```

---

## 3. Accepted W4 authority that remains current until final ratification

- W4 maps **95/95** admitted Product operations; zero additions/unmapped operations;
- there are **29 flat exact stored Permissions**; `/access-context` uses an authenticated special condition, not a stored Permission;
- D2 Principal kinds remain `H human | A automation | S system`;
- ordinary `both` Q/read → H/A/S; consequential/business C admitted to both → H/A; side-effect-free `EvaluatePriceScenario` → H/A/S;
- `system` never substitutes for business `automation`;
- Permission is flat/exact: no manage→read, execute→read, prefix or wildcard implication; AccessRoles may bundle exact Permissions;
- Principal kind, Membership, Permission, client-class admission, owner business disposition, Governance and physical epistemic authority remain separate gates;
- `RecordSeparation` / `RecordPacking` are H-only; `RecordPhysicalConference` / `RecordDispatchHandoff` are H or explicitly qualified S;
- physical qualification is not Permission, AccessRole, fourth Principal kind or generic machine-capability graph;
- IdP/OAuth/provider roles/scopes/frontend guards never independently grant Product access;
- no-Membership path Organization is privacy-preserving 404; current Membership with missing Permission/client-class admission is 403;
- monotonic AccessRole/AuthorizationDelegation revocation changes target-concurrency semantics only, never caller admission.

---

## 4. Converged Whole-W4 corrections — REVIEW RESULT, NOT YET CANONICAL

Whole-W4 independent review re-derived all 95 operations and 29 Permissions, confirmed the accepted matrix, and found five additive gaps beyond lead G1/G2. GPT adjudication accepted them with no Round 2.

### W4-C1 — Principal access eligibility anchored in D2

Final filing proposes one bounded D2 confirmation:

> **Current Principal access eligibility is Principal-scoped revocable D2 identity/access state. Disabling/revoking that eligibility blocks future Product access, including `/access-context`, without deleting Organization Membership/RoleAssignment or rewriting historical actor attribution. Exact lifecycle/representation mechanics remain later realization.**

This confirms D2 §6.1 Principal lifecycle ownership plus B2-A's already-accepted disablement requirement; it creates no new domain, Permission or IAM framework.

Failure grammar:

- missing/invalid/untrusted/wrong-audience credential → `401 authentication-required`;
- accepted credential resolving to zero MPC Principals → `401 authentication-required`;
- duplicate binding / more than one Principal resolution → fail closed as existing W2 `500 internal-error`, never select one;
- exactly one Principal resolved but currently ineligible → `403 access-denied`;
- eligible Principal without current path-Organization Membership → privacy-preserving `404 resource-not-found`.

### W4-C2 — current-authority property for every mutable/revocable access gate

Every mutable/revocable access fact is evaluated against its **current authoritative state** at the Product boundary. External credential validity/binding follows the accepted authentication authority. MPC-owned Principal eligibility, Membership, RoleAssignment/Permission, Principal kind and operation-specific physical qualification follow current MPC authority.

Stale token claims, cached snapshots or retired provisioning records never remain access authority after revocation. D7 may cache only while preserving this semantic property.

### W4-C3 — `/access-context` eligibility

`GET /access-context` requires:

```text
valid accepted authentication
+ exactly one MPC Principal binding
+ current Principal access eligibility
```

It waives only Organization Membership and Product Permission because it discovers them. An eligible Principal with zero memberships receives a successful empty Organization set; an ineligible Principal receives 403 and no membership/role/permission disclosure.

### W4-C4 — physical qualification ordering / non-self-assertion

For system-admitted physical checkpoints:

```text
current Membership
→ allowed Principal kind
→ exact ordinary Permission
→ current server-resolved operation-specific physical qualification
→ resource/reference privacy
→ W1/W2/domain/Governance gates
```

A request body, client-declared station/device/evidence source, token claim, IdP/OAuth/provider role/scope or frontend state cannot self-assert trusted/qualified physical authority. Successful system-established checkpoints preserve server-attributed effective Principal/source/time.

### W4-C5 — mutation-response disclosure is operation-scoped

A mutation/capability Permission authorizes the W1/W2 response representation for that **same operation and exact operation subject**. This disclosure does not grant the corresponding Get/List/Search operation, another subject/resource or Permission inheritance.

Do not require a second read Permission merely to receive the operation's normal response, and do not vary normal response shape based on possession of a read Permission.

Flat-Permission negative controls are therefore about **read-operation admission**, e.g. `portfolio.manage` does not grant `portfolio.read` operations and `fulfillment.manage` does not grant `fulfillment.read`/`fulfillment.execute` operations.

### W4-C6 — sensitive/special read cases are explicitly reviewed, no new Permission now

- FulfillmentArtifact List/Get remains `fulfillment.execute`; W4 explicitly resolves the ratified `human/both read` ambiguity to H/A/S for the read operation, without conferring checkpoint authority. A later proven consumer needing artifact read without execute access reopens only this split.
- `ListAuthorizationDelegations` remains H + `governance.manage`; standing delegation topology is authorization-management-sensitive and no proven read-only auditor consumer requires a split. A real auditor use reopens only this read boundary.
- `GetBusinessSystemPartyResolution` and `GetDestinationRealization` remain under `materialization.read`; that Permission is explicitly PII-bearing and must be assigned accordingly. A real consumer needing Materialization tracking without party/destination detail reopens the smallest read/Permission boundary.

No speculative `artifact.read`, `delegation.read` or `materialization.pii.read` Permission is introduced. Stored Permission count remains 29.

### Whole-W4 dispositions preserved

- 95/95 operations mapped; zero operations added by symmetry;
- 29/29 Permissions used; no orphan Permission;
- three Principal kinds only: H/A/S;
- flat exact Permissions confirmed as Global Maximum;
- `EvaluatePriceScenario` remains H/A/S under `economics.read`;
- physical checkpoint matrix unchanged except current qualification/revocation hardening;
- `fulfillment.execute` artifact boundary retained;
- 404/403 Organization privacy split retained;
- Governance/business/ordinary-access separation retained;
- monotonic revocations retain full caller admission;
- Structural Inversion against current IdP/middleware/frontend/provider access shape: PASS;
- no D0/D1/D3/D4/D4-R1/D5-B1/W1/W2/W3/Operation-Matrix reopen;
- D2 requires one-line authority confirmation only, not semantic reopen;
- Round 2 not required.

---

## 5. Prohibited now

Until final operator Whole-W4 ratification:

- do not amend D2 or W4 with the converged package;
- do not begin technical ingress classification, final OpenAPI/tooling, D6–D9 or implementation;
- do not treat Whole-W4 candidate, `AI-DIALOG.md`, cockpit, current code or middleware as authority;
- do not add Principal kinds, Permission hierarchies/wildcards, IAM/ReBAC DSL, generic machine-capability graph or speculative sensitive-read Permissions;
- do not infer Product access from Keycloak/OAuth/provider roles/scopes or frontend visibility;
- do not choose D7 claim/cache/RLS/middleware/provisioning realization yet.

---

## 6. Exact next action

**Operator performs final Whole-W4 ratification of the converged package in §4, including the one-line D2 Principal-access-eligibility confirmation.**

If ratified:

1. add the bounded D2 confirmation;
2. consolidate all ratified corrections into the single W4 authority;
3. remove the Whole-W4 review candidate; Git history remains archive;
4. reset `AI-DIALOG.md` to protocol-only;
5. update cockpit as non-authoritative projection;
6. advance router to **technical non-Product ingress classification**.

Implementation remains blocked until D9.

---

## 7. Fresh-session success test

A fresh session must conclude:

- W1/W2/W3 are canonical;
- W4 remains accepted in-stage and maps 95/95 operations with 29 flat exact Permissions;
- Whole-W4 lead + independent Fable + GPT adjudication are complete and converged;
- converged package includes G1/G2 plus F-W4-1…F-W4-5 and a one-line D2 Principal-access-eligibility confirmation;
- **operator final Whole-W4 ratification is the exact next action**;
- technical ingress/later Wire obligations/D6–D9/implementation remain blocked.

If not, the active authority tree is inconsistent.