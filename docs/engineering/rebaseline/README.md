# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **B2-A + Operation Matrix Blocks 1–4 ACCEPTED IN-STAGE; Block 5 Fulfillment + Post-Sale + Work + P compositions = NEXT**  
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
       └─ Matrix Block 5 — Fulfillment + Post-Sale + Work + P compositions — NEXT / UNDER DERIVATION
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

## 4. What is prohibited now

While D5-B2 is OPEN / ACTIVE:

- do not begin D6–D9 target design or implementation;
- do not silently alter accepted D0–D4/D4-R1/D5-B1 or ratified B2 in-stage decisions;
- do not derive operations from legacy routes/current OpenAPI/provider endpoints;
- do not recreate Product/PIM master, generic Integration/Mutation/Workflow/Rules/AI authority, generic finance ledger or market collector platform;
- do not merge Offering, Availability, Economics, Governance, Sales, Materialization, Fulfillment, Post-Sale or Work because provider/ERP workflow combines fields or calls;
- do not create global/shared MPC API keys or treat Keycloak roles/Organizations as MPC business authority;
- do not expose provider OAuth, source sync, readiness/market refresh, reconcile-now, ERP retry or runtime commands merely for convenience;
- do not create direct client commands for owner reactions already defined by D3 (for example normal-path BusinessOrderIntent/InvoicingIntent creation);
- do not weaken Organization scope, source-qualified identity, honest knowledge/freshness, Permission/Governance separation, idempotency, concurrency, ambiguity, recovery, multi-target scope or convergence laws;
- do not add compatibility/versioning or generic bulk without a real entitled consumer/workflow;
- do not treat retained legacy ADRs, `AI-DIALOG.md`, chat summaries or Git history as current target authority.

## 5. Exact next action

**Derive D5-B2 Operation Admission Matrix Block 5 — Fulfillment Lifecycle + Post-Sale Resolution + Operational Work + justified read-only P compositions.**

The block must establish the smallest Product API surface for:

1. **Fulfillment Lifecycle**
   - Fulfillment Node read/configuration only where Product clients need it;
   - physical separation/conference/packing/dispatch checkpoints as owner-specific client capabilities rather than generic workflow status;
   - provider-requirement closure/readiness under Fulfillment authority, with provider artifacts/evidence remaining D4 source truth;
   - source-qualified Shipment/delivery observation through relevant terminal outcomes;
   - physical conference/readiness as the legitimate owner signal feeding Materialization, never a direct invoice command;
   - no company-wide WMS/TMS, generic shipment mutation or provider protocol mirror.

2. **Post-Sale Resolution**
   - scoped Post-Sale Resolution read/create/advance semantics only for real cancellation/return/refund consequence workflows;
   - 0..N resolutions per Sale and line/quantity scope where material;
   - provider Claim/Return/refund/reverse-shipment resources remain external/source-qualified evidence, not MPC ontology;
   - no CRM/SAC/general reverse-logistics platform and no direct provider refund/cancel API by protocol vocabulary;
   - closure only when applicable consequences are sufficiently evidenced, never from one provider terminal flag alone.

3. **Operational Work**
   - Work list/get/assignment/responsibility/lifecycle operations for real operator queues;
   - Work closure never mutates originating source truth;
   - resolution evidence returns to source owner for semantic closure when required;
   - source-domain independent resolution can reconcile Work through its own events;
   - no generic Task/Case/Workflow engine beyond accepted Work meaning.

4. **Read-only P compositions**
   - admit only cross-owner read compositions with a real operator/manager consumer, such as operational attention/cockpit summaries;
   - projection remains read-only and may carry component freshness/partiality;
   - no P endpoint may become write/concurrency/authorization/retry authority;
   - D6 screen topology does not become API authority merely because a dashboard needs data.

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
- B2-A and Matrix Blocks 1–4 accepted in-stage;
- Block 5 Fulfillment + Post-Sale + Work + P compositions is the exact next action;
- Product API auth remains OIDC/OAuth with MPC-owned access/business authority;
- BusinessOrder/Invoicing normal-path intent creation remains owner-triggered, not client-commanded;
- no stale pre-R1 B2 candidate is active;
- implementation remains blocked until D9.

If not, the active authority tree is inconsistent.