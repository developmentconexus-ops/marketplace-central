# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **B2-A accepted; Matrix Blocks 1 and 2 accepted; Block 3 Market Intelligence + Commercial Economics = NEXT**  
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

This file alone owns **where the program is and what happens next**. `ARCHITECTURE.md` owns stable constraints; the Decision Reconciliation Baseline routes current decision generations; the ADR registry owns ADR status; accepted D-stage/B2 artifacts own detailed semantics.

`D5-API.md` remains the accepted D5-B1 authority. Its former next-action wording is a pre-B2-opening snapshot. Current B2 state and next action live here.

Do not reconstruct target authority from memory, Git history, retired ADRs, stale candidates, `AI-DIALOG.md` or current code shape.

## 2. Program state

```text
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
  ↓
Decision Reconciliation Baseline — ACCEPTED / CANONICAL
  ↓
D5 — API — OPEN / ACTIVE
  ├─ B1 Semantic API Model & Contract Laws — ACCEPTED / CANONICAL
  └─ B2 Product Operation / Resource Surface — OPEN / ACTIVE
       ├─ B2-A Client & Authentication Admission Model — ACCEPTED IN-STAGE
       ├─ Matrix Block 1 — Identity/Access + Portfolio + Readiness — ACCEPTED IN-STAGE
       ├─ Matrix Block 2 — Offering + Price + Availability — ACCEPTED IN-STAGE
       └─ Matrix Block 3 — Market Intelligence + Commercial Economics — NEXT / UNDER DERIVATION
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

## 3. Accepted current B2 decisions

### B2-A — Client / authentication

- one standards-based OIDC/OAuth boundary for Product API clients;
- humans use Authorization Code + PKCE semantics;
- confidential machine clients use Client Credentials/service-account semantics;
- MPC remains authority for Principal, Membership, AccessRole, Permission and all business/Governance meaning;
- tokens are audience-bound to MPC API and another sibling application's token is not valid merely because the same IdP issued it;
- no global/static MPC API key or browser secret baseline;
- Keycloak is first implementation/proof candidate; D7 owns concrete provider/deployment/realm/HA/secrets/token-lifetime realization.

### Matrix Block 1 — Identity/Access + Portfolio + Readiness

Admitted:

- current access context + minimal Membership/RoleAssignment administration;
- Marketplace Installation lifecycle/configuration + Selling Entity discovery;
- marketplace-context source Product search;
- Product↔channel readiness;
- publication requirements;
- explicit correspondence lifecycle.

Rejected/deferred:

- IAM/custom-role platform;
- SaaS Organization provisioning;
- Product/PIM CRUD;
- generic Integration/provider catalog;
- provider OAuth as Product API;
- Product/source sync/refresh commands;
- speculative bulk.

### Matrix Block 2 — Offering + Price + Availability

Admitted:

- provider-actual Marketplace Listing Q surface;
- `ListingIntent` create/edit draft lifecycle and Q tracking;
- declarative draft update + discard + submit with concurrency where stale overwrite is material;
- `PriceIntent` Q + explicit create for exact desired price actuation;
- Sellable Availability Q + convergence;
- Inventory Source lifecycle/configuration;
- Availability allocation/scope-policy Q + human configuration.

Binding fences:

- no giant Listing CRUD owning content + price + stock + fulfillment;
- no direct provider Listing create/update/set-price/set-stock Product operations;
- ListingIntent is the one create/edit authoring identity;
- PriceIntent stays separate from ListingIntent and from Commercial Economics calculation;
- no public PriceDraft baseline;
- no public AvailabilityIntent creation merely because the internal identity exists;
- no public sync/refresh mechanism;
- no generic `LongRunningOperation`; owner-local Intents are the Product tracking resources;
- R1-G1 remains owner-preserving: Offering and Availability may be jointly serialized by D4/D7 without ownership merge;
- provider early success never equals whole-operation convergence;
- intended, authorized and actual attempted/provider-affected scopes remain distinct when blast radius is material.

Current Permission floor across accepted blocks:

- `access.read`, `access.manage`
- `portfolio.read`, `portfolio.manage`
- `readiness.read`, `readiness.manage`
- `offering.read`, `listing.manage`, `price.manage`
- `availability.read`, `availability.manage`

These are ordinary-access capabilities only; they do not grant business disposition or Governance authority.

## 4. What is prohibited now

While D5-B2 is OPEN / ACTIVE:

- do not begin D6–D9 target design before D5 is accepted as a whole;
- do not implement product features; implementation remains blocked until D9;
- do not silently alter accepted D0–D4/D4-R1/D5-B1 or accepted in-stage B2 meaning;
- do not derive routes from legacy OpenAPI/controllers/packages or the retired pre-R1 B2 candidate;
- do not recreate Product/PIM, PublicationPreparation, SourceProductObservation, generic Integration/Mutation/Workflow/Rule/AI authority;
- do not merge Listing representation, Price, Availability, Economics or Fulfillment because provider protocol places them together;
- do not create generic async `Operation` business identity when the owner-local Intent already tracks the work;
- do not choose Keycloak deployment/runtime mechanics in B2;
- do not weaken Organization scope, source-qualified identity, honest knowledge/freshness, Permission/Governance separation, idempotency, concurrency, ambiguity, recovery or convergence;
- do not add compatibility/versioning machinery without a real entitled consumer;
- do not treat retained legacy ADRs/review dialogue/current implementation as target authority.

## 5. Exact next action

**Derive D5-B2 Operation Admission Matrix Block 3 — Market Intelligence + Commercial Economics.**

The block must establish the smallest Product API surface that preserves these accepted authorities:

### Market Intelligence

- comparable-market observations and source-qualified evidence;
- comparability interpretation;
- competitive position/change;
- market-evidence sufficiency/insufficiency;
- no generic provider-market payload mirror;
- no public collector/scraper/source-refresh command unless a real Product client requires an owner-level capability rather than D4/D7 mechanism.

### Commercial Economics

- Cost Basis and exact-money economic interpretation;
- price/profitability simulation and trade-offs;
- expected economics (L0);
- order economics (L1);
- realized/settlement economics (L2);
- variance/calibration/reconciliation;
- modeled/observed/realized provenance remains distinct;
- Economics never writes marketplace price; actual price actuation stays Offering `PriceIntent`.

The block must decide proportionately:

1. whether price/economics simulation is stateless, durable, or split according to real consumer need — no persistent `Simulation` entity merely for history;
2. which market evidence/competitive-position reads need pagination/filtering and which do not;
3. how insufficient/stale/partial evidence blocks false precision;
4. whether any recommendation resource is independently justified or is simply an economic conclusion returned by simulation/current analysis;
5. how expected/order/realized economics remain separate without one mutable profitability row overwriting evidence stages;
6. the minimum economic reconciliation read/capability surface without a generic finance ledger;
7. whether any Economics policy configuration belongs in Product API now and which concrete actor needs it;
8. no campaign/discount authoring merely because observed promotions affect price/economics.

For every candidate operation record the full B2 admission tuple: consumer, client class, owner, Q/C/P, Organization, Permission, identity, knowledge/freshness, consequence/idempotency/concurrency, provider enrichment, collection semantics, bulk, ADMIT/REJECT/DEFER.

Do **not** spell final HTTP paths/schemas until Block 3 operation inventory is coherent.

If a required operation cannot fit accepted authority without distortion, reopen only the implicated parent decision.

Implementation remains blocked until D9.

## 6. Fresh-session success test

A fresh session must conclude:

- D0–D4 and D4-R1 are accepted/canonical;
- Decision Reconciliation is accepted/canonical;
- D5-B1 is accepted/canonical;
- D5-B2 is OPEN / ACTIVE;
- B2-A is accepted in-stage;
- Matrix Blocks 1 and 2 are accepted in-stage;
- Block 2 preserves ListingIntent/PriceIntent/Availability authority split and rejects giant Listing CRUD/direct stock-price mutation/generic Operation;
- Block 3 Market Intelligence + Commercial Economics is NEXT;
- implementation remains blocked until D9;
- exact next action is to derive market/economic Product operations, not paths or code.

If it cannot, the active authority tree is inconsistent.