# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **B2-A + Operation Matrix Blocks 1–5 ACCEPTED IN-STAGE; Whole-Matrix Global Coherence review = NEXT**  
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

`D5-API.md` remains D5-B1 authority. Its old next-action wording is a pre-B2-opening snapshot. Current B2 status/next action is defined only here.

Never reconstruct target authority from memory, chat, Git history, retired ADRs, `AI-DIALOG.md`, stale candidates or current code/OpenAPI shape.

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
       └─ Whole-Matrix Global Coherence — NEXT / UNDER REVIEW
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

## 3. Accepted B2 routing summary

### B2-A — Client/Auth

- Product API authentication uses one standards-based OIDC/OAuth boundary.
- Humans use Authorization Code + PKCE semantics; confidential machine clients use Client Credentials/service-account semantics.
- MPC remains authority for Principal, Organization Membership, AccessRole/Permission/RoleAssignment and every business decision.
- Tokens are audience-bound to MPC API; client/request never supplies effective Principal or business approval.
- No global/static MPC Product API key or duplicate IdP-role business authority.
- Keycloak remains first implementation/proof candidate; D7 owns concrete provider/deployment/realm/secrets/token-lifetime realization.

### Matrix Block 1

Admitted only consumer-proven D2 access context/role-assignment administration, Portfolio Marketplace Installation lifecycle/configuration + Selling Entity discovery, and Readiness marketplace-context Product discovery/readiness/publication requirements/correspondence.

Not admitted: Product/PIM CRUD, generic IAM platform, Organization SaaS provisioning, generic Integration API, provider sync/refresh commands or speculative bulk.

### Matrix Block 2

- Listing actual state is Offering Q; no direct provider-shaped Listing mutation.
- `ListingIntent` is the single create/edit authoring/tracking identity with draft concurrency and controlled submit.
- `PriceIntent` is separate from Economics and uses durable intent semantics.
- Sellable Availability is Availability-owned; no public `SetAvailableQuantity` or baseline public AvailabilityIntent authoring.
- Provider may jointly serialize Offering + Availability inputs without ownership merge/cross-owner atomicity.
- No generic `LongRunningOperation`, giant Listing CRUD, direct Price set or generic mutation/action surface.

### Matrix Block 3

- Market Intelligence exposes competitive position/comparable evidence, not generic MarketObservation CRUD/collector commands.
- `EvaluatePriceScenario` is stateless/side-effect-free; simulations/recommendations do not gain durable IDs by default.
- Expected/Sale Economics preserve honest L0/L1/L2 lineage, coverage and R1/R2 reconciliation without one mutable profitability row.
- Commercial policy remains Economics-owned; Governance does not acquire business thresholds.
- Economic Attribution is persistent Economics state; explicit ambiguous resolution is human baseline.
- No generic financial ledger, universal Reconciliation resource, public `ReconcileNow`, bank/R3 API or price actuation inside Economics.

### Matrix Block 4

- Governance exposes Authorization Decisions and bounded Delegation/Grant administration; approval never mutates Intent or executes effects.
- Marketplace Sales is externally originated and read-centric; Product clients do not create/update provider sales.
- explicit human Selling Entity attribution resolution is admitted only for genuine Sales ambiguity.
- BusinessOrderIntent is Materialization-owned and normally created from committed Sale meaning, not client commands.
- InvoicingIntent is Materialization-owned and normally created/advanced from Fulfillment physical-readiness checkpoints, not direct invoice commands.
- Party Resolution exposes bounded human resolution without Customer/CRM mastery.
- Destination Realization Q is admitted; write resolution remains conditioned on D8 controlled proof.
- no direct Sankhya TOP/NUNOTA/order/invoice/retry/workflow API and no blind replay after possible native acceptance.

### Matrix Block 5

- Fulfillment exposes explicit physical checkpoints and Fulfillment Node configuration without generic status/WMS/TMS authority.
- physical facts may be established only by a client/Principal legitimately capable of establishing them; a generic automation token cannot fabricate physical conference.
- Shipment remains source-qualified external identity and read-only observation in the selected lane.
- Post-Sale exposes canonical scoped Resolution read/create while provider Claim/Return/refund action vocabulary stays outside Product API until concrete actions are proven.
- Work owns responsibility/assignment/hold/escalation/resolution submission; arbitrary user-created Task/Case and direct close are not baseline.
- `GetSaleOperationalView` is the only baseline cross-owner P composition and remains read-only, component-permission-bound and non-authoritative for writes/concurrency.

## 4. What is prohibited now

While D5-B2 is OPEN / ACTIVE:

- do not begin D6–D9 target design or implementation;
- do not silently alter accepted D0–D4/D4-R1/D5-B1 or ratified B2 in-stage decisions;
- do not derive operations from legacy routes/current OpenAPI/provider endpoints;
- do not recreate Product/PIM master, generic Integration/Mutation/Workflow/Rules/AI authority, generic finance ledger, Task/Case engine or market collector platform;
- do not merge Offering, Availability, Economics, Governance, Sales, Materialization, Fulfillment, Post-Sale or Work because provider/ERP workflow combines fields or calls;
- do not create global/shared MPC API keys or treat Keycloak roles/Organizations as MPC business authority;
- do not expose provider OAuth, source sync, readiness/market refresh, reconcile-now, ERP retry or runtime commands merely for convenience;
- do not create direct client commands for owner reactions already defined by D3 (for example normal-path BusinessOrderIntent/InvoicingIntent creation);
- do not allow a machine client to establish a human/physical fact solely because it holds ordinary API Permission;
- do not allow a projection/P endpoint to bypass component permissions or become write/concurrency/authorization/retry authority;
- do not weaken Organization scope, source-qualified identity, honest knowledge/freshness, Permission/Governance separation, idempotency, concurrency, ambiguity, recovery, multi-target scope or convergence laws;
- do not add compatibility/versioning or generic bulk without a real entitled consumer/workflow;
- do not treat retained legacy ADRs, `AI-DIALOG.md`, chat summaries or Git history as current target authority.

## 5. Exact next action

**Run D5-B2 Whole-Matrix Global Coherence + YAGNI / Permissions / Client-Class / Missing-Operation review across B2-A + Blocks 1–5 as one system.**

The review must challenge, proportionately:

1. duplicate or missing business/API authority;
2. Product 1.0 outcome reachability from legitimate Product clients or accepted owner-triggered reactions;
3. human vs machine/system client-class correctness, especially physical facts and standing human decisions;
4. Permission coherence and least privilege without one-per-endpoint fragmentation;
5. Q/C/P classification honesty;
6. consequential idempotency defaults and every claimed structural exemption;
7. concurrency/lost-update coverage independent of idempotency;
8. owner-triggered versus client-triggered lifecycle decisions;
9. Organization/source-qualified identity and same-Organization secondary-reference safety;
10. unknown/unavailable/partial/stale/provider-enriched read semantics;
11. generic-abstraction pressure: Product/PIM, Integration, Mutation/Action/Operation, Workflow, Rule, Finance, Task/Case, Provider graph;
12. pagination/filter/search/bulk admission only for real collection/workflow consumers;
13. P projection permission/freshness/concurrency fences;
14. future-cost/YAGNI seams for second marketplace/business system, stronger automation, additional fulfillment/post-sale modes and eventual public/external clients;
15. Structural Inversion against legacy routes/OpenAPI/controllers;
16. explicit missing-operation challenge.

Allowed review outcomes:

- `RESTRUCTURE NOW`
- `CURRENT STRUCTURE CONFIRMED`
- `STOP / SPLIT PREREQUISITE`

Do not spell final Product API paths/schemas until this Whole-Matrix review is operator-ratified.

If the matrix passes, the following B2 sub-batch will define resource/path grammar, standard versus owner-specific HTTP methods, request/response schema families, status/problem outcomes, `Idempotency-Key` and precondition placement, pagination/filter/search grammar and OpenAPI operation spelling — still without D6/D7 implementation choices.

Implementation remains blocked until D9.

## 6. Fresh-session success test

A fresh session must conclude unambiguously:

- D0→D4 + D4-R1 accepted/canonical;
- Decision Reconciliation accepted/canonical;
- D5-B1 accepted/canonical;
- D5-B2 OPEN / ACTIVE;
- B2-A and Matrix Blocks 1–5 accepted in-stage;
- Whole-Matrix Global Coherence review is the exact next action;
- Product API auth remains OIDC/OAuth with MPC-owned access/business authority;
- BusinessOrder/Invoicing normal-path intent creation remains owner-triggered, not client-commanded;
- physical Fulfillment facts cannot be fabricated by generic automation authority;
- Work/projections never replace originating owner truth;
- no stale pre-R1 B2 candidate is active;
- implementation remains blocked until D9.

If not, the active authority tree is inconsistent.