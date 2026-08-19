# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **W1 + W2 + W3 CANONICAL; W4 ACCEPTED IN-STAGE; lead Whole-W4 review COMPLETE / RESTRUCTURE W4-LOCAL; operator ratification of W4-G1/G2 = NEXT**  
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

`D5-B2-W4-WHOLE-COHERENCE-REVIEW-CANDIDATE.md` is **NON-AUTHORITATIVE lead review evidence** and is deliberately outside the authority path. W4-G1/G2 do not modify W4 or parent authority until operator ratification and later canonical filing.

`AI-DIALOG.md` is protocol-only with no active review cycle. `docs/engineering/rebaseline/cockpit.html` is a non-authoritative visual projection and is synchronized only after canonical status changes.

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
        Whole-W4 lead coherence                          COMPLETE / RESTRUCTURE W4-LOCAL
        W4-G1/G2 operator direction                      NEXT
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

## 3. Accepted W4 authority that remains current during review

- W4 maps **95/95** admitted Product operations; zero operation added/unmapped;
- baseline has **29 flat exact stored Permissions**; `/access-context` uses authenticated special condition, not a Permission;
- D2 Principal kinds remain `H human | A automation | S system`;
- ordinary `both` Q/read → H/A/S; consequential/business C admitted to both → H/A; side-effect-free `EvaluatePriceScenario` → H/A/S;
- systems do not become business automation; AI/agent/repricer/business automation uses automation Principal semantics;
- Permission is flat: no manage→read, execute→read, wildcard or prefix implication; AccessRoles may bundle exact Permissions;
- client-class admission, ordinary Permission, owner business disposition, Governance and physical epistemic authority remain separate;
- `RecordSeparation` and `RecordPacking` are H-only baseline;
- `RecordPhysicalConference` and `RecordDispatchHandoff` admit H or explicitly qualified S;
- physical qualification is not Permission, AccessRole, fourth Principal kind or generic machine-capability graph;
- IdP roles/scopes/provider roles/frontend guards never grant Product access;
- no-Membership Organization access is privacy-preserving 404; current Membership but missing Permission/class admission is 403;
- monotonic AccessRole/AuthorizationDelegation revocation changes target concurrency semantics only, never caller-access requirements.

---

## 4. Lead Whole-W4 findings — NON-AUTHORITATIVE UNTIL OPERATOR DECISION

Lead review re-derived all 95 operations/29 Permissions and found **two W4-local gaps; no parent-stage reopen**:

### W4-G1 — current Principal eligibility / binding failure grammar — REVISE

- after external AuthN, credential must resolve uniquely to one MPC Principal;
- then W4 must explicitly consume **current D2 Principal access eligibility** before Membership;
- missing/invalid/untrusted/wrong-audience or not-uniquely-resolvable credential → `401 authentication-required`;
- resolved but currently disabled/not access-eligible Principal → `403 access-denied`;
- eligible Principal with no current path-Organization Membership → privacy-preserving `404 resource-not-found`;
- no new Principal lifecycle enum/resource is introduced; D2 remains owner.

### W4-G2 — physical qualification ordering / non-self-assertion — REVISE

Refine the access sequence:

```text
current Membership
→ allowed Principal kind
→ exact ordinary Permission
→ if applicable: server-resolved operation-specific physical qualification
→ resource/reference privacy
→ W1/W2/domain/Governance
```

- client body/token/role/scope cannot self-assert trusted/physical/qualified authority;
- H retains the accepted accountable-human physical baseline with `fulfillment.execute`;
- S requires `fulfillment.execute` + server-resolved qualification for PhysicalConference/DispatchHandoff;
- successful system-established checkpoint preserves actual Principal/source/time provenance under W2;
- qualification remains Fulfillment-operation-specific and does not become Permission/Principal kind/generic capability graph.

The following survived lead attack unchanged: H/A/S on accepted `both` reads; `fulfillment.execute` for artifact reads; 404-vs-403 privacy split; flat Permissions; `EvaluatePriceScenario` under `economics.read`; all currently human-only operation admissions.

---

## 5. Prohibited now

Until operator ratifies/revises W4-G1/G2:

- do not silently apply G1/G2 to accepted W4;
- do not begin technical ingress classification, final OpenAPI/tooling, D6–D9 or implementation;
- do not treat the Whole-W4 candidate/cockpit/current code/middleware as authority;
- do not add Principal kinds, Permission hierarchy/wildcards, IAM/ReBAC DSL or generic system capability graph;
- do not infer Product access from Keycloak/OAuth/provider roles/scopes or frontend visibility;
- do not give deferred/rejected operations Permission rows by symmetry;
- do not choose D7 Keycloak/claim/cache/RLS/middleware realization yet.

---

## 6. Exact next action

**Operator reviews and ratifies/revises W4-G1 and W4-G2 as the lead Whole-W4 correction direction.**

If ratified, determine whether proportional independent Fable challenge adds material value before final W4 canonical consolidation. W4 is a high-impact access/security contract, so independent challenge is expected unless evidence makes it redundant. Reviewer output remains evidence only.

After final Whole-W4 convergence + operator final ratification:

1. consolidate ratified corrections into the single W4 authority;
2. remove Whole-W4 review candidate; Git history remains archive;
3. reset any review channel used;
4. update cockpit as non-authoritative projection;
5. advance to technical non-Product ingress classification.

Implementation remains blocked until D9.

---

## 7. Fresh-session success test

A fresh session must conclude:

- W1/W2/W3 are canonical;
- W4 accepted authority maps 95/95 operations with 29 flat exact Permissions;
- Whole-W4 lead review found only W4-G1/G2 and `RESTRUCTURE NOW — W4-LOCAL`;
- Whole-W4 candidate is non-authoritative evidence;
- **operator ratification/revision of W4-G1/G2 is the exact next action**;
- no parent reopen is currently proven;
- later Wire obligations/D6–D9/implementation remain blocked.

If not, the active authority tree is inconsistent.
