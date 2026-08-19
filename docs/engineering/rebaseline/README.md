# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **Operation Matrix + Whole-Matrix RATIFIED; Wire W1 + W2 + W3 CANONICAL; W4 Permission → Operation / Client-Class Enforcement ACCEPTED IN-STAGE; Whole-W4 adversarial coherence = NEXT**  
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
19. `docs/engineering/rebaseline/D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT.md` — accepted W4 design home
20. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
21. code/OpenAPI/schemas/tests/runtime only as current-state evidence when needed

This file alone owns **program status, allowed/blocked work and exact next action**. Detailed semantics remain in the accepted artifacts above.

`AI-DIALOG.md` is protocol-only with no active review cycle. Former W2/W3 staging/review artifacts are absent from the active tree; Git history is the archive.

`docs/engineering/rebaseline/cockpit.html` is a **non-authoritative visual projection**. It never participates in the authority path and is synchronized only after canonical status changes.

Legacy/current code, routes, OpenAPI, SDK, IdP roles/scopes, middleware and frontend route/button visibility remain evidence only, never target access authority.

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
      W3 Collections / Pagination / Filter / Search /
         Cursor Grammar                                  ACCEPTED / CANONICAL
      W4 Permission → Operation / Client-Class
         Enforcement                                     ACCEPTED IN-STAGE
        Whole-W4 adversarial coherence                   NEXT
      Technical non-Product ingress classification       BLOCKED BY W4
      Final Problem/media consistency                    BLOCKED BY SEQUENCE
      Single OpenAPI wire authority/tooling decision     BLOCKED BY SEQUENCE
D6 — Frontend                                             BLOCKED BY D5
D7 — Runtime / Jobs / Transactions                        BLOCKED
D8 — Golden Flows                                         BLOCKED
D9 — Adversarial Architecture Review                      BLOCKED
Implementation                                            BLOCKED UNTIL D9
```

Whole-W2 and Whole-W3 adversarial review cycles are complete and incorporated into canonical Wire artifacts. W4 is the only active Wire design surface.

---

## 3. Load-bearing authority W4 may not weaken

### 3.1 Identity / ordinary access

- external OIDC/OAuth authenticates; MPC owns Principal, Membership, RoleAssignment, AccessRole and Permission;
- D2 Principal kinds remain `human | automation | system`; no fourth `physical_system`/agent/service-account kind;
- AccessRole is a product-defined bundle; Permission is the stable ordinary-access capability consumed by backend/business entry checks;
- ordinary Permission does **not** prove business disposition, automation eligibility, Governance authorization, execution validity or epistemic ability;
- Organization-owned non-human actions remain explicitly Organization-scoped;
- `GET /access-context` is the one platform-scoped self-only Product discovery Q.

### 3.2 W1/W2/W3

- Organization is explicit path scope and never self-authorization;
- one Product Problem Details catalog exists in W2;
- business rejection/approval-required/unknown/external-required remain Product semantics rather than generic access failure;
- source-qualified identities and same-Organization privacy remain fail-closed;
- W3 collection/query semantics are independent from access enforcement;
- implementation/provider/DB/Keycloak mechanisms do not become Product authority.

### 3.3 Accepted W4 core

- exactly **95/95** admitted Product operations have one W4 access/client-class row; no new Product operation was added;
- exactly **29 stored ordinary Permissions** remain in the baseline; `authenticated` for `/access-context` is a special operation condition, not a Permission;
- Permissions are **flat and exact**: no `*.manage ⇒ *.read`, no `*.execute ⇒ *.read`, no wildcard/prefix implication;
- AccessRoles may bundle exact Permissions but do not change Permission semantics;
- B2 coarse `both` is crystallized: ordinary Q/read → `H|A|S`; side-effect-free `EvaluatePriceScenario` → `H|A|S`; consequential/business-authoring C admitted to both → `H|A` unless explicitly narrower/different;
- `system` is not business automation by default; AI/agent/repricer/business automation uses `automation` Principal semantics;
- client-class admission is separate from Permission and from Governance/business disposition;
- Governance never widens a human-only Product operation to automation/system;
- `fulfillment.execute` does not itself prove physical fact authority;
- `RecordSeparation` + `RecordPacking` are human-only baseline;
- `RecordPhysicalConference` + `RecordDispatchHandoff` allow human OR an explicitly qualified `system` Principal/source;
- the physical-system qualification is operation/Fulfillment-specific, not a Permission, Principal kind or generic machine-capability graph;
- IdP roles/scopes/provider roles/frontend guards never independently grant Product access;
- current Membership/RoleAssignment/Permission state governs future access; revocation must not be defeated by stale token-carried Product role authority;
- no-Membership Organization access fails closed without proving another Organization exists;
- missing Permission or disallowed Principal kind/required physical qualification is ordinary access denial, not owner business rejection;
- fail-safe monotonic target revocation semantics remain separate from caller admission.

---

## 4. Prohibited now

While Whole-W4 coherence is next:

- do not begin technical ingress classification, final OpenAPI/tooling closure, D6–D9 target design or implementation;
- do not reopen accepted D0→D5-B1/B2/W1/W2/W3 for naming/style or implementation convenience;
- do not implement Keycloak realm/client/role mapping, token claims, middleware, RLS/ACL, cache or persistence mechanics;
- do not derive Product Permission from IdP roles/scopes, provider roles/scopes, frontend guards or current middleware;
- do not add Permission hierarchy/wildcards/prefix matching or generic IAM/ReBAC/policy DSL;
- do not make `system` a generic business-automation escape hatch;
- do not create `physical_system` as a fourth Principal kind or generic `system_capabilities[]` merely to satisfy checkpoints;
- do not collapse ordinary access, Governance, business disposition and physical epistemic qualification;
- do not add Permission mappings for deferred/rejected operations by symmetry;
- do not choose exact OpenAPI minor/generator or W3 numeric limit defaults yet.

---

## 5. Exact next action

**Run the Whole-W4 adversarial coherence sweep over all 95 admitted Product operations, all 29 ordinary Permissions and every allowed Principal/client-class rule.**

The review must challenge at least:

1. operation coverage is exactly 95/95 with zero duplicates/unmapped/new operations;
2. every operation has exactly one ordinary Permission/special access condition and one allowed Principal set;
3. no coarse `both` ambiguity remains after H/A/S crystallization;
4. every human-only operation has a real reason and does not block a currently accepted machine consumer;
5. every H/A operation genuinely represents business automation rather than a system fact/source;
6. generic `system` cannot accidentally author Listing/Price/Work/Post-Sale business state;
7. the two qualified-system physical checkpoints have enough epistemic fencing without a generic capability framework;
8. Separation/Packing human-only versus PhysicalConference/DispatchHandoff qualified-system split remains coherent;
9. `fulfillment.execute` for artifact reads remains proportionate to PII/operational sensitivity;
10. Permission vocabulary/splits preserve least privilege without hidden implication — especially `listing.manage != price.manage`, Fulfillment read/execute/manage and Economics/Governance splits;
11. broad-looking names such as `sales.manage` / `post_sale.manage` do not create wildcard authority or require premature rename;
12. `access-context`, no-Membership privacy, secondary-reference privacy and current revocation behavior are coherent;
13. auth/client-class/access failure remains distinct from domain business outcomes and Governance;
14. side-effect-free `EvaluatePriceScenario` under `economics.read` is honest and does not need a fake execute Permission;
15. monotonic AccessRole/AuthorizationDelegation revocations retain caller-access enforcement while avoiding stale-target authority survival;
16. Structural Inversion against current Keycloak roles/scopes, middleware and frontend visibility still passes;
17. no D0→W3 parent reopen is actually required.

If no material contradiction survives, present Whole-W4 to the operator for final ratification before making W4 canonical and moving to technical non-Product ingress classification.

If a material contradiction survives, return only the smallest implicated W4 row/Permission/client-class or parent operation admission to the Decision Loop. Do not proceed by convenience.

Implementation remains blocked until D9.

---

## 6. Fresh-session success test

A fresh session must conclude:

- D0→D4/D4-R1 + D5-B1 are accepted/canonical;
- D5-B2 Operation Matrix + Whole-Matrix are ratified;
- W1, W2 and W3 are canonical;
- `D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT.md` is the accepted W4 design home;
- W4 maps **95/95** admitted Product operations using **29 flat exact Permissions** plus the special authenticated `/access-context` condition;
- H/A/S, business automation and physical-system epistemic qualification remain distinct;
- **Whole-W4 adversarial coherence is the exact next action**;
- technical ingress/OpenAPI/D6–D9 remain blocked by sequence;
- implementation remains blocked until D9.

If not, the active authority tree is inconsistent.
