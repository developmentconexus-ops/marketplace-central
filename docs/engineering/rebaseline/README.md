# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **B2-A + Operation Matrix Blocks 1–5 ACCEPTED IN-STAGE; Whole-Matrix review candidate PREPARED; independent Fable review = NEXT**  
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
15. `docs/engineering/rebaseline/D5-B2-OPERATION-ADMISSION-MATRIX.md`
16. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
17. code, OpenAPI, schemas, tests and runtime only as current-state evidence when needed

This file alone owns **where the program is and what happens next**. `ARCHITECTURE.md` owns stable cross-stage constraints; the Decision Reconciliation Baseline routes current decision generations; the ADR registry owns ADR status; accepted D-stage/B2 artifacts own detailed semantics in their scope.

The current review candidate `D5-B2-WHOLE-MATRIX-REVIEW-CANDIDATE.md` and `AI-DIALOG.md` are **not authority** and are deliberately excluded from the authority path. They are review input only.

`D5-API.md` remains D5-B1 authority. Its old next-action wording is a pre-B2-opening snapshot. Never reconstruct target authority from memory, chat, Git history, retired ADRs, `AI-DIALOG.md`, review candidates or current code/OpenAPI shape.

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
  ↓
D4 — External Integrations — CLOSED / ACCEPTED AS A WHOLE
  └─ D4-R1 Publication Input & Listing Authoring — ACCEPTED / CANONICAL
       └─ R1-G1 ML initial publication × Availability — PASS-B
  ↓
Decision Reconciliation Baseline — ACCEPTED / CANONICAL
  ↓
D5 — API — OPEN / ACTIVE
  ├─ B1 Semantic API Model & Contract Laws — ACCEPTED / CANONICAL
  └─ B2 Product Operation / Resource Surface — OPEN / ACTIVE
       ├─ B2-A Client & Authentication Admission Model — ACCEPTED IN-STAGE
       ├─ Matrix Block 1 — Identity/Access + Portfolio + Readiness — ACCEPTED IN-STAGE
       ├─ Matrix Block 2 — Offering + Price + Availability — ACCEPTED IN-STAGE
       ├─ Matrix Block 3 — Market Intelligence + Commercial Economics — ACCEPTED IN-STAGE
       ├─ Matrix Block 4 — Governance + Sales + Materialization — ACCEPTED IN-STAGE
       ├─ Matrix Block 5 — Fulfillment + Post-Sale + Work + P compositions — ACCEPTED IN-STAGE
       └─ Whole-Matrix Global Coherence
            ├─ lead review — RESTRUCTURE NOW / B2-local corrections identified
            ├─ review candidate — PREPARED / NON-AUTHORITATIVE
            └─ Fable independent review — NEXT
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

## 3. Accepted B2 baseline before whole-matrix review

### B2-A — Client/Auth

- Product API authentication uses one standards-based OIDC/OAuth boundary.
- Humans use Authorization Code + PKCE semantics; confidential machine clients use Client Credentials/service-account semantics.
- MPC remains authority for Principal, Organization Membership, AccessRole/Permission/RoleAssignment and all business decisions.
- Tokens are audience-bound to MPC API; no global/static MPC Product API key or IdP-role business authority.
- Keycloak is the first implementation/proof candidate; D7 owns provider/deployment/realm/secrets/token-lifetime realization.

### Matrix Blocks 1–5

- **Block 1:** minimal D2 access context/role assignment, Portfolio Installation lifecycle/configuration, and marketplace-context Readiness Product discovery/requirements/correspondence; no PIM/IAM/integration platform.
- **Block 2:** Listing actual state is Offering Q; `ListingIntent` is create/edit authoring/tracking; `PriceIntent` is separate; Availability owns Sellable Availability; no giant Listing CRUD, direct price/stock set or generic async Operation.
- **Block 3:** Market Intelligence exposes competitive interpretation; Economics owns stateless scenario evaluation plus durable material L0/L1/L2 lineage; no Recommendation/Simulation authority, generic ledger or price actuation in Economics.
- **Block 4:** Governance decisions/delegations remain authorization-only; Sales is externally originated/read-centric; Materialization creates BusinessOrder/Invoicing intents from accepted owner reactions; no direct Sankhya/order/invoice/retry/workflow API.
- **Block 5:** Fulfillment exposes physical checkpoints/nodes/artifacts, Shipment remains external read observation, Post-Sale uses canonical scoped Resolution, Work owns responsibility/lifecycle without source truth, and one provisional Sale operational P was admitted for whole-matrix challenge.

These block decisions remain accepted in-stage until the whole-matrix review package is adjudicated and operator-ratified.

## 4. Current non-authoritative whole-matrix review package

`docs/engineering/rebaseline/D5-B2-WHOLE-MATRIX-REVIEW-CANDIDATE.md` records the lead's operator-approved review direction for independent challenge.

Lead disposition: **RESTRUCTURE NOW — B2-local only; no parent-stage reopen currently justified.**

Proposed corrections under review:

1. **ADD** ListingIntent-scoped authored-media intake under Offering; no ProductAsset/media master.
2. **ADD** Fulfillment-owned internal operating-target Q/C; external provider deadline remains distinct.
3. **DEFER** generic `SubmitWorkResolution`; use source-owner-specific resolution capabilities unless a concrete bounded evidence-submission need is later proven.
4. **DEFER** `GetSaleOperationalView` P until D6 proves repeated consumer need/benefit.

Proposed hardenings under review:

5. `ResolveBusinessSystemPartyResolution` requires client idempotency by default.
6. `GetCurrentAccessContext` is a bounded platform-scoped D2 discovery Q; Organization-owned business routes remain explicit Organization-path scoped.
7. Authorization Delegation update/revoke gains stale-state concurrency/precondition protection where material.

These corrections are not canonical merely because they appear here; Fable review + GPT adjudication + operator ratification precede consolidation into the active matrix.

## 5. What is prohibited now

While the independent Whole-Matrix review is open:

- do not begin resource/path/schema/OpenAPI crystallization yet;
- do not begin D6–D9 design or implementation;
- do not mutate accepted parent D0–D4/D4-R1/D5-B1 semantics by review convenience;
- do not treat the candidate, Fable output or `AI-DIALOG.md` as authority;
- do not derive operations from legacy routes/current OpenAPI/provider endpoint shape;
- do not recreate Product/PIM, generic Integration/Mutation/Action/Operation/Workflow/Rules/AI authority, generic finance ledger, Task/Case engine or market collector platform;
- do not weaken Organization scope, source-qualified identity, honest knowledge/freshness, Permission/Governance separation, idempotency, concurrency, ambiguity, recovery, multi-target scope or convergence laws;
- do not create direct client commands for owner reactions already owned by D3 flows;
- do not allow a generic machine token to fabricate physical facts;
- do not add compatibility/versioning or bulk without a real consumer/workflow.

## 6. Exact next action

**Run one independent Fable review of the coherent D5-B2 Whole-Matrix package before any wire-contract design.**

Follow the canonical **Standard Fable review workflow** in `developmentconexus-ops/conexus-methodology/README.md`.

Fable must:

1. independently reconstruct this repository's authority from `AGENTS.md` + this router;
2. read accepted B2-A + Blocks 1–5 and the non-authoritative `D5-B2-WHOLE-MATRIX-REVIEW-CANDIDATE.md`;
3. apply the DevelopmentConexus Method and search for a materially better Global Maximum, not agreement;
4. challenge duplicate/missing authority, Product 1.0 reachability, client classes, Permissions, Q/C/P, idempotency, concurrency, owner-trigger vs client-trigger, Organization/source identity, provider richness, YAGNI/future cost, Structural Inversion and missing operations;
5. specifically attack the four proposed corrections and three hardenings rather than assuming them correct;
6. append **material findings only** to the active `AI-DIALOG.md` cycle with `APPROVE / REVISE / REJECT` and handoff to GPT;
7. modify no other repository file unless separately authorized by the operator.

After Fable, GPT independently adjudicates every material finding. Round 2 occurs only if a material contradiction survives. The converged package then requires operator ratification before corrections are consolidated into the active matrix and the disposable candidate is removed.

If the Whole-Matrix package is then accepted, the next B2 sub-batch is **Wire Contract / Resource-Path-Schema Grammar**: concrete path/resource hierarchy, standard vs owner-specific HTTP methods, request/response families, Problem Details, idempotency/preconditions, pagination/filter/search and OpenAPI spelling — still without D6/D7 implementation choices.

Implementation remains blocked until D9.

## 7. Fresh-session success test

A fresh session must conclude unambiguously:

- D0→D4 + D4-R1 accepted/canonical;
- Decision Reconciliation accepted/canonical;
- D5-B1 accepted/canonical;
- D5-B2 OPEN / ACTIVE;
- B2-A and Matrix Blocks 1–5 accepted in-stage;
- the Whole-Matrix lead review found B2-local corrections but no parent-stage reopen;
- `D5-B2-WHOLE-MATRIX-REVIEW-CANDIDATE.md` is non-authoritative and prepared for Fable;
- `AI-DIALOG.md` has the D5-B2 Whole-Matrix review cycle open;
- independent Fable review is the exact next action;
- wire-contract design is blocked until that review is adjudicated/ratified;
- implementation remains blocked until D9.

If not, the active authority tree is inconsistent.