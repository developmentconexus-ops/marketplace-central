# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **B2-A + Operation Matrix Blocks 1–3 ACCEPTED IN-STAGE; Block 4 Governance + Sales + Materialization = NEXT**  
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
       └─ Matrix Block 4 — Governance + Sales + Materialization — NEXT / UNDER DERIVATION
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
- Humans use Authorization Code + PKCE semantics.
- Confidential machine clients use Client Credentials/service-account semantics.
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
- Inventory Source and Availability policy/configuration remain Availability-owned.
- Provider may jointly serialize Offering + Availability inputs without ownership merge or cross-owner atomicity.
- No generic `LongRunningOperation`, giant Listing CRUD, direct Price set or generic mutation/action surface.

### Matrix Block 3

- Market Intelligence exposes competitive position/comparable evidence, not generic MarketObservation CRUD or collector commands.
- `EvaluatePriceScenario` is a stateless, side-effect-free Commercial Economics capability; simulations/recommendations do not gain durable IDs by default.
- Expected Economics and Sale Economics preserve honest L0/L1/L2 lineage, coverage and reconciliation without one mutable profitability row.
- Commercial policy remains Economics-owned; Governance does not acquire business thresholds.
- Economic Attribution is persistent Economics state where meaning/correlation is exact/partial/ambiguous/unresolved; explicit ambiguous resolution is human baseline.
- No generic financial ledger, universal Reconciliation resource, public `ReconcileNow`, bank/R3 API or price actuation inside Economics.

## 4. What is prohibited now

While D5-B2 is OPEN / ACTIVE:

- do not begin D6–D9 target design or implementation;
- do not silently alter accepted D0–D4/D4-R1/D5-B1 or ratified B2 in-stage decisions;
- do not derive operations from legacy routes/current OpenAPI/provider endpoints;
- do not recreate Product/PIM master, generic Integration/Mutation/Workflow/Rules/AI authority, generic finance ledger or market collector platform;
- do not merge Offering, Availability, Economics, Governance, Sales, Materialization or Fulfillment because a provider/ERP workflow combines fields or calls;
- do not create global/shared MPC API keys or treat Keycloak roles/Organizations as MPC business authority;
- do not expose provider OAuth, source sync, readiness/market refresh, reconcile-now or runtime commands merely for operational convenience;
- do not weaken Organization scope, source-qualified identity, honest knowledge/freshness, Permission/Governance separation, idempotency, concurrency, ambiguity, recovery, multi-target scope or convergence laws;
- do not add compatibility/versioning or generic bulk without a real entitled consumer/workflow;
- do not treat retained legacy ADRs, `AI-DIALOG.md`, chat summaries or Git history as current target authority.

## 5. Exact next action

**Derive D5-B2 Operation Admission Matrix Block 4 — Controlled Action Governance + Marketplace Sales + Business-System Materialization.**

The block must establish the smallest Product API surface for:

1. **Controlled Action Governance**
   - Authorization Decision read/create semantics for real approval consumers;
   - authorization grant/delegation administration only where Product 1.0 actually needs it;
   - exact intended/authorized scope preservation;
   - ordinary Permission distinct from domain disposition/Governance;
   - approval never mutating Intent, executing effects or waiving execution-time revalidation;
   - no generic approval/workflow/policy engine.

2. **Marketplace Sales**
   - source-qualified marketplace Sale listing/get semantics and honest acquisition/freshness;
   - transaction-specific Selling Entity attribution;
   - no synthetic MPC sale/order alias merely for normalization;
   - no downstream Materialization/Fulfillment/Economics/Post-Sale ownership leakage;
   - provider Order/Pack/Shipment topology remains D4 evidence, not Product ontology.

3. **Business-System Materialization**
   - Business Order Intent + native-order convergence;
   - Invoicing Intent + authoritative fiscal/document convergence;
   - decide whether these intents are directly client-created or initiated only by accepted upstream owner conditions/workflows;
   - bounded Party Resolution and Destination Realization operations where a real human/client decision is required;
   - no Customer/Address master CRUD;
   - no direct `/sankhya/orders`, `/sankhya/invoices`, TOP/NUNOTA/status/choreography Product API;
   - consequential create idempotency and no-blind-retry after possible native acceptance;
   - explicit ambiguous/duplicate/native-correlation failures and Work rather than hidden retry.

For every candidate, apply the complete B2 admission tuple: consumer, client class, owner, Q/C/P, Organization, Permission, subject identity, knowledge/outcome, idempotency, concurrency/preconditions, provider enrichment, collection semantics and bulk.

Do not spell final paths/schemas yet.

If a required operation cannot fit accepted ownership/communication/external-contract meaning, reopen only the implicated parent decision.

Implementation remains blocked until D9.

## 6. Fresh-session success test

A fresh session must conclude unambiguously:

- D0→D4 + D4-R1 accepted/canonical;
- Decision Reconciliation accepted/canonical;
- D5-B1 accepted/canonical;
- D5-B2 OPEN / ACTIVE;
- B2-A and Matrix Blocks 1–3 accepted in-stage;
- Block 4 Governance + Sales + Materialization is the exact next action;
- Product API auth remains OIDC/OAuth with MPC-owned access/business authority;
- no stale pre-R1 B2 candidate is active;
- implementation remains blocked until D9.

If not, the active authority tree is inconsistent.