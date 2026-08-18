# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **B2-A Client & Authentication Admission Model accepted in-stage; Operation Admission Matrix = NEXT**  
> **Decision Reconciliation:** **ACCEPTED / CANONICAL**  
> **Implementation:** BLOCKED until D9 is accepted  
> **Last updated:** 2026-08-18

## 1. Authority path

A fresh session reads, in order:

1. `AGENTS.md`
2. this file
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
15. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
16. code, OpenAPI, schemas, tests and runtime only as current-state evidence when needed

This file alone owns **where the program is and what happens next**. `ARCHITECTURE.md` owns stable constraints; the Decision Reconciliation Baseline routes current decision generations; the ADR registry owns ADR status; D-stage artifacts own detailed semantics.

`D5-API.md` remains the accepted D5-B1 authority. Its former “next action” wording is a pre-B2-opening snapshot; current D5-B2 status/next action is defined here and detailed in the active B2 artifact.

Do not reconstruct target authority from memory, Git history, retired ADRs, `AI-DIALOG.md`, review candidates or current code shape.

## 2. Program state

```text
Documentary / governance cleanup — DONE
  ↓
D0 — Product / System Definition — CLOSED / ACCEPTED
  ↓
D1 — Domains / Boundaries — CLOSED / ACCEPTED
  ↓
D2 — Identity / Tenant / Data Ownership — CLOSED / ACCEPTED
  ↓
D3 — Communication / Events — CLOSED / ACCEPTED
  ├─ B1 Communication Topology & Edge Matrix — ACCEPTED / CANONICAL
  ├─ B2 Communication Contract & Failure Semantics — ACCEPTED / CANONICAL
  └─ Final Global Coherence — PASS
  ↓
D4 — External Integrations — CLOSED / ACCEPTED AS A WHOLE
  ├─ B1 External Contract Grounding — ACCEPTED / CANONICAL
  ├─ B2 Mercado Livre Operational Contract — ACCEPTED / CANONICAL
  ├─ B3 Sankhya Business-System Contract — ACCEPTED / CANONICAL
  ├─ B4 Market / Economics / Settlement — ACCEPTED / CANONICAL
  ├─ Original Global Coherence — PASS
  └─ R1 Publication Input & Listing Authoring — ACCEPTED / CANONICAL
       └─ R1-G1 ML initial publication × Availability — PASS-B
  ↓
Decision Reconciliation Baseline — ACCEPTED / CANONICAL
  ├─ D0→D4/D4-R1 + D5-B1 decision set — RECONCILED / COHERENT
  ├─ legacy ADR active tree — reduced to real D7/Fact/transition residues
  └─ stale pre-R1 D5-B2 candidate — RETIRED TO GIT HISTORY
  ↓
D5 — API — OPEN / ACTIVE
  ├─ B1 Semantic API Model & Contract Laws — ACCEPTED / CANONICAL
  └─ B2 Product Operation / Resource Surface — OPEN / ACTIVE
       ├─ B2-A Client & Authentication Admission Model — ACCEPTED IN-STAGE
       └─ Operation Admission Matrix — NEXT / UNDER DERIVATION
  ↓
D6 — Frontend
  ↓
D7 — Runtime / Jobs / Transactions
  ↓
D8 — Golden Flows
  ↓
D9 — Adversarial Architecture Review
  ↓
Implementation DAG / Plan
  ↓
Implementation
```

Product implementation remains blocked until D9 is accepted.

## 3. Accepted baseline — routing summary

Use `DECISION-RECONCILIATION-BASELINE.md` to discover the current decision generation, then read the detailed semantic home.

Load-bearing current conclusions:

- **D0:** MPC is Marketplace Operations Control Plane + Commercial Intelligence; external systems retain their own truth; no ERP/PIM/generic integration/workflow product.
- **D1:** 12 semantic business authorities remain accepted; they do not imply 12 services/processes/databases.
- **D2:** Organization is isolation root; Product remains source-qualified external identity; domain intents are owner-local; ordinary access is distinct from Governance/business disposition; Principal includes human/automation/system and interactive AuthN is external via OIDC.
- **D3:** Q/C/E/P hybrid, recoverable consequential propagation, durable evidence-occurrence recovery, no global exactly-once/order, no blind replay of ambiguous effects, no cross-owner atomicity.
- **D4:** consumer owns meaning; adapter owns protocol; sanctioned Sankhya Gateway only; provider richness without DTO mirroring; external-effect reread/convergence; explicit source admissibility and PII minimization.
- **D4-R1:** Readiness owns requirements/source-level readiness; Offering owns `ListingIntent` create/edit draft + dispatchability; baseline values are `FOLLOW_SOURCE | EXPLICIT_OVERRIDE`; no PublicationPreparation/PIM/source-observation/rules/AI framework; provider may jointly realize owner-issued meanings without ownership merge.
- **D5-B1:** semantic Product API, Organization path scope, source-qualified wire identity, honest knowledge/freshness, fail-closed consequential idempotency, RFC 9457 problems, one OpenAPI wire authority, hard cutover, operation-local bulk only.
- **D5-B2-A:** Product API clients use one standards-based OIDC/OAuth authentication boundary; humans use Authorization Code + PKCE semantics, confidential machine clients use Client Credentials/service-account semantics, MPC remains Principal/Membership/Permission authority, tokens are audience-bound to MPC API, no global/static MPC Product API key or duplicate IdP-role business authority is baseline. Keycloak remains the first implementation/proof candidate; D7 owns concrete provider/deployment/realm realization.

Detailed rules, Unknowns, proof obligations and reopen triggers remain in the named D-stage homes.

## 4. What is prohibited now

While **D5-B2 is OPEN / ACTIVE**:

- do not begin D6–D9 target design before D5 is accepted as a whole;
- do not implement product features; implementation remains blocked until D9;
- do not silently alter accepted D0–D4/D4-R1/D5-B1 or accepted in-stage B2-A meaning;
- do not derive B2 from the legacy route list or the retired pre-R1 candidate;
- do not preserve current OpenAPI/routes/SDK/controller/package shape by inheritance;
- do not recreate Product/PIM master, PublicationPreparation, SourceProductObservation owner, generic Mutation/Workflow/Integration platform, generic listing rule engine or AI-specific authority path;
- do not move Availability/Fulfillment meaning into Offering because a provider combines fields in one request;
- do not create a global/shared MPC API key, browser client secret, client-supplied Principal, or treat Keycloak/IdP roles/Organizations as MPC business authority by name similarity;
- do not choose Keycloak deployment/realm/HA/token-lifetime/secret-storage mechanics inside B2; those remain D7 realization;
- do not weaken Organization scope, source-qualified identity, honest knowledge/freshness, Permission/Governance separation, idempotency, precondition, ambiguity, recovery or convergence laws;
- do not add compatibility/versioning machinery without a real entitled consumer;
- do not treat retained legacy ADRs as target authority beyond the exact residue named by the ADR registry;
- do not treat `AI-DIALOG.md`, Git-history ADRs or review artifacts as target authority.

## 5. Exact next action

**Derive the D5-B2 Operation Admission Matrix from Product 1.0 actors/consumers and accepted semantic owners.**

The matrix is derived owner by owner, not route by route. For every candidate Product API operation establish:

- real Product 1.0 actor/client consumer and concrete use;
- allowed client class: human, machine/automation/system, or both;
- exactly one accepted semantic owner or D2 substrate authority;
- Q / C / P interaction class;
- explicit Organization path scope and same-Organization secondary-reference rule;
- ordinary Permission requirement distinct from business disposition/Governance;
- canonical/source-qualified subject identity;
- honest knowledge/freshness/provenance for reads;
- consequential Intent/outcome/idempotency/precondition/concurrency laws where applicable;
- owner-specific convergence and multi-step/partial outcome semantics where effects are involved;
- projection/read-only status where composition is justified;
- provider-enriched fields only for named consumer/correctness needs;
- pagination/filter/sort/cursor only when a real consumer requires them;
- bulk only for a real workflow with member-level correctness;
- D4-R1 publication authoring through `ListingIntent`, never Product/PIM or provider-field-bag semantics.

Admission predicate:

> **A Product API operation exists only when a real Product 1.0 client needs to read an accepted owner's meaning (Q), ask one accepted owner to perform/accept owner-owned work (C), or consume a justified read-only composition (P). Symmetry, current code, provider endpoints and internal implementation convenience are insufficient.**

Derivation order:

1. D2 identity/access client needs;
2. Marketplace Portfolio;
3. Product & Channel Readiness;
4. Marketplace Offering Operations;
5. Availability Control;
6. Market Intelligence;
7. Commercial Economics;
8. Controlled Action Governance;
9. Marketplace Sales;
10. Business-System Materialization;
11. Fulfillment Lifecycle;
12. Post-Sale Resolution;
13. Operational Work;
14. justified read-only P compositions.

Do **not** spell final paths/schemas until the admission inventory is coherent enough that naming cannot hide duplicate/missing authority.

If an operation cannot fit accepted ownership/identity/communication/external-contract meaning without distortion, stop and reopen only the implicated parent decision.

Implementation remains blocked until D9.

## 6. Fresh-session success test

A fresh session must conclude:

- D0, D1, D2, D3, D4 and D4-R1 are accepted/canonical;
- Decision Reconciliation Baseline is accepted/canonical and is routing, not a second semantic architecture;
- ADR registry contains only real unresolved/transition legacy residues plus future target ADRs;
- Sankhya Gateway-only target transport and no Direct Oracle fallback are unambiguous;
- D5-B1 is accepted/canonical;
- **D5-B2 is OPEN / ACTIVE**;
- **B2-A Client & Authentication Admission Model is accepted in-stage**;
- Product API authentication is OIDC/OAuth standards-based while MPC retains Principal/Membership/Permission/business authority;
- Keycloak is first implementation/proof candidate but D7 owns concrete provider/deployment/realm realization;
- **Operation Admission Matrix is NEXT / UNDER DERIVATION**;
- the stale pre-R1 B2 candidate is not in the active tree;
- `AI-DIALOG.md` contains only the reusable review protocol, not historical review authority;
- implementation remains blocked until D9;
- exact next action is to derive admitted Product 1.0 operations from real consumers and accepted semantic owners.

If it cannot, the active authority tree is inconsistent.
