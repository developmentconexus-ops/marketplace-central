# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D4 — EXTERNAL INTEGRATIONS — OPEN / ACTIVE; D4-B1 ACCEPTED / CANONICAL; D4-B2 ACCEPTED / CANONICAL; D4-B3 NEXT**  
> **Implementation:** BLOCKED until D9 is accepted  
> **Last updated:** 2026-08-17

## 1. Authority path

A fresh session reads, in order:

1. `AGENTS.md`
2. this file
3. `docs/engineering/standards/root-cause-global-maximum-method.md`
4. `ARCHITECTURE.md`
5. `docs/architecture/decisions/README.md`
6. `docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md`
7. `docs/engineering/rebaseline/D1-DOMAINS-BOUNDARIES.md`
8. `docs/engineering/rebaseline/D2-IDENTITY-TENANT-DATA-OWNERSHIP.md`
9. `docs/engineering/rebaseline/D3-COMMUNICATION-EVENTS.md`
10. `docs/engineering/rebaseline/D4-EXTERNAL-INTEGRATIONS.md`
11. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
12. code, OpenAPI, schemas, tests and runtime only as current-state evidence when needed

This file alone answers **where the program is and what happens next**. Stable architecture belongs in `ARCHITECTURE.md`; accepted/current stage semantics belong in D-stage artifacts; Git history is the archive.

Do not reconstruct target authority from memory, legacy package shape, historical plans, `AI-DIALOG.md`, review candidates or stale docs.

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
  └─ Final Global Coherence + YAGNI / Overengineering / Future-Cost — COMPLETED / PASS
  ↓
D4 — External Integrations — OPEN / ACTIVE
  ├─ B1 External Contract Grounding — ACCEPTED / CANONICAL
  ├─ B2 Mercado Livre Operational Contract — ACCEPTED / CANONICAL
  │    └─ Installation Evidence Gate — CLOSED / PASS
  ├─ B3 Sankhya Business-System Contract — NEXT / NOT YET OPENED
  ├─ B4 Market / Economics / Settlement Contract — NOT YET OPENED
  └─ Final D4 Global Coherence + YAGNI / Overengineering / Future-Cost — NOT STARTED
  ↓
D5 — API
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

## 3. Accepted baseline

### D0 — CLOSED

Product 1.0 is **Marketplace Operations + Commercial Intelligence**. MPC is the marketplace operations control plane: external systems retain authority for facts/processes inherently theirs while MPC owns the cross-system marketplace operating semantics needed to observe, understand, decide, execute, verify and reconcile.

D0 authority and non-goals are defined only in `D0-PRODUCT-SYSTEM-DEFINITION.md`.

### D1 — CLOSED

`D1-DOMAINS-BOUNDARIES.md` is the accepted D1 authority. It defines 12 business boundaries, explicit ownership/non-ownership, cross-cutting non-domain treatment, legal semantic authority edges, forbidden boundary violations, legacy semantic disposition, later-stage defers and reopen triggers.

The 12 D1 business boundaries do **not** imply 12 services, databases, processes or deployments.

### D2 — CLOSED

`D2-IDENTITY-TENANT-DATA-OWNERSHIP.md` is the accepted D2 authority. Its accepted direction includes:

- canonical identity follows semantic authority;
- Organization is the canonical tenant/isolation root;
- Marketplace Installation and SourceInstance qualify external namespaces without collapsing credentials/transport into identity;
- source-qualified external Product, Listing/Variation, Sale/Order, Shipment and financial identities;
- stable domain-local identity for material MPC Business Intents and accepted Work/Governance/Economics identities;
- explicit Principal/identity/access state with a strict business-authority fence;
- one D1 semantic write authority per canonical business meaning;
- exact value/knowledge/time semantics and honest unknowns;
- clean target persistence baseline and legacy-ADR rehoming gates.

### D3 — CLOSED

`D3-COMMUNICATION-EVENTS.md` is the accepted D3 authority. Its accepted direction includes:

- semantic hybrid **Q/C/E/P** communication;
- current truth from owners; material historical occurrence from the smallest sufficient durable authority;
- recoverable consequential propagation;
- explicit Organization scope;
- semantic idempotency under duplicate delivery;
- no global arrival-order authority;
- known/known-empty/unknown/unavailable query semantics;
- accepted/rejected/pending/ambiguous capability semantics where applicable;
- no blind replay of ambiguous external effects;
- projections as rebuildable read state, never write authority;
- shared execution-safety mechanisms verify owner-issued proofs without owning business disposition/authorization;
- no generic Event Bus/Command Bus/Workflow engine/event sourcing/exactly-once/global-ordering/runtime topology choice.

### D4 — OPEN / B1+B2 ACCEPTED

`D4-EXTERNAL-INTEGRATIONS.md` is the active D4 authority.

D4-B1 accepted:

- concrete provider/business-system adapters implement consumer-owned semantic ports; no Integration business domain or universal provider entity graph;
- Mercado Livre Installation binds fail-closed to the authoritative external seller namespace, including acquisition-time attribution when the provider exposes a namespace marker;
- Sankhya SourceInstance remains stable across sanctioned credential/protocol mechanics while the authoritative namespace is unchanged;
- credentials/auth are protocol/runtime secrets, not business identity;
- provider notification is a trigger/pointer and current provider meaning comes from authoritative reread where material;
- point/enumeration/delta/notification coverage claims are operation-scoped and fail-honest;
- Integration Support, Provider Effective Capability/Requirement and Effective Business Capability remain distinct authorities;
- later external effects must define acceptance/ambiguity and an authoritative reconciliation surface; HTTP/provider transport success is not convergence;
- **Sankhya API Gateway is the target transport for MPC↔Sankhya integration. Direct Oracle/database access is explicitly outside the target architecture and is not a fallback path;**
- B3 must prove Gateway/API correctness, coverage and operational viability; if the sanctioned surface cannot satisfy a required Product 1.0 claim, B3 stops and returns to explicit operator/architecture adjudication rather than enabling database access implicitly;
- ADR-004 D4 plugin-framework meaning is superseded; ADR-010 polling-only D4 meaning is superseded while D7 runtime residue remains; ADR-006/007 are historical for target architecture after their transport-independent lessons were rehomed.

D4-B2 accepted:

- Mercado Livre provider topology stays inside the adapter; Item, User Product, Family, Catalog Product, provider stock locations, Claim and Return do not become MPC business ontology by normalization;
- listing creation/observation is mode-aware and current-provider-authoritative while Offering retains Listing/Price intent/convergence;
- Product↔channel provider identifiers are evidence for Readiness, never automatic correspondence authority;
- price observation authority is separated from current write mechanism; Price Automation is a Provider Effective restriction and provider 2xx never substitutes for authoritative convergence proof;
- User Product shared-field effects, including applicable availability writes, must not silently widen intended/authorized scope;
- stock writability depends on concrete provider resource plus site, seller configuration and current resource/listing context; seller-managed does not automatically mean API-writable, and provider-managed Full stock is not seller/MPC-writable by convenience;
- stale provider stock version conflict is a definitive rejected precondition followed by reread/redecision, not ambiguous blind retry;
- Order and Shipment remain separate provider resources; seller Order search alone does **not** prove cancellation-inclusive Sales coverage because current official documentation and real Installation behavior conflict on canceled-Order inclusion;
- fulfillment responsibility, provider fiscal/label/readiness requirements and external SLA/deadline remain provider-context-sensitive evidence translated for Fulfillment closure;
- essential Claim/Return/reverse-shipment capability is included without expanding into CRM/SAC; refund/payment/settlement movement authority remains B4;
- Full/provider-operated observation alone does not prove MPC-controlled Availability convergence or internally operated Fulfillment Node execution;
- no universal `OperatingMode`/provider graph/framework and no support for every documented ML mode is introduced.

The B2 Installation Evidence Gate closed **PASS** using a fresh read-only real-dependency probe. The selected current first Mercado Livre proof context is:

- current seller/listing model: User Product; complete seller scan observed 34/34 current Items as UP, zero legacy variations;
- current Item↔UP relation: 1:1 across the observed population at measurement time;
- current Availability lane: non-multi-origin Item-path `available_quantity`, with shared-UP blast-radius revalidation retained;
- current Price candidates: not blocked by Price Automation at measurement time;
- current Sale/Fulfillment lane: seller-operated `me2 / xd_drop_off`; no Full/multi-origin/flex lane observed;
- Shipment SLA proven live; real post-sale Claim/Return + reverse Shipment evidence observed.

These are time-bound Installation facts, not permanent provider promises. B2 residual R1 — first controlled ML Price/Availability write + authoritative reread/convergence — belongs D8. B2 residual R2 — live selected-lane fiscal/label progression — is constrained by B3 materialization semantics and later proven in D8.

No D0/D1/D2/D3 reopen is currently required.

## 4. Engineering method and repo lifecycle

Engineering reasoning follows the **DevelopmentConexus Engineering Method** identified in `AGENTS.md`; the local file in this authority path is only the consumed context copy.

This router defines the Marketplace Central D0–D9 status/lifecycle and allowed work. It is repo-specific specialization, not a second organizational engineering method. Conflicts inside the organizational method's scope must be surfaced, never silently reinterpreted.

## 5. What is prohibited now

While D4-B3 is the active next decision batch:

- do not begin B4 or D5–D9 target design prematurely;
- do not implement product architecture/features;
- do not silently alter accepted D0–D4-B2 authority;
- do not let Sankhya API payloads, current Oracle schema, legacy adapters, DTOs or historical ADRs become target business authority by inheritance;
- do not create semantic dependencies outside D1 or bypass D2/D3 semantics through integration code;
- do not introduce Direct Oracle/database access for Sankhya as an implementation shortcut, read optimization or Gateway fallback;
- do not infer a Gateway fact/command is sufficient merely because an endpoint exists; exact semantics, namespace, coverage and operational viability must be proven for the Product 1.0 claim;
- do not choose HTTP/frontend/runtime topology prematurely;
- do not treat `AI-DIALOG.md`, review candidates or reviewer/chat summaries as target authority.

Existing code/module/context/schema shape remains current-state evidence only.

## 6. Exact next action

**Open D4-B3 — Sankhya Business-System Contract from canonical D4-B1+B2.**

B3 must independently prove the sanctioned Sankhya API Gateway surface required by Product 1.0, including:

- SourceInstance/authoritative namespace qualification and any exposed source identity markers;
- authoritative Product/native Product key and provider evidence needed by Readiness;
- inventory facts with company/location/as-of semantics needed by Availability;
- cost/tax evidence and source qualifiers needed by Commercial Economics;
- Business Order Intent materialization into the correct native Sankhya order and authoritative result correlation/reread;
- Invoicing Intent materialization into the correct native fiscal/document result and authoritative result correlation/reread;
- pagination/window/change semantics and honest coverage;
- request/response size, rate limits, timeout/blocking behavior and operational viability where material;
- explicit unsupported/unknown facts or commands rather than Oracle shortcuts.

B3 must preserve the D4-B1 rule: if a materially required Product 1.0 fact/command cannot be satisfied correctly and operationally through the sanctioned Gateway/API surface, **STOP / SPLIT PREREQUISITE** and return to explicit operator/architecture adjudication. Direct Oracle/database is not admitted as fallback.

Do not begin B4 until B3 is accepted under the current router. Implementation remains blocked until D9.

## 7. Fresh-session success test

A fresh session should conclude that:

- D0, D1, D2 and D3 are **CLOSED / ACCEPTED**;
- D4 is **OPEN / ACTIVE**;
- D4-B1 is **ACCEPTED / CANONICAL**;
- D4-B2 is **ACCEPTED / CANONICAL** and its Installation Evidence Gate is **CLOSED / PASS**;
- D4-B3 is **NEXT / NOT YET OPENED**;
- B4 is **NOT YET OPENED**;
- consumer-owned ports and provider-local protocol remain distinct;
- Marketplace Installation/SourceInstance external namespace binding is explicit and fail-closed;
- Mercado Livre Item/UP/stock/Order/Shipment/Claim topology remains provider-local evidence rather than MPC business ontology;
- seller Order search completion alone does not prove cancellation-inclusive coverage;
- selected current Mercado Livre proof lane is User Product + non-multi-origin Item availability + direct-price candidate + seller-operated `xd_drop_off`, with effect convergence/fiscal-label progression reserved for D8/B3+D8 proof;
- **Sankhya API Gateway is the target transport and Direct Oracle is not an admitted fallback;**
- a Gateway/API capability gap causes STOP / explicit re-adjudication rather than database access by convenience;
- current modules/contexts/schema/legacy ADRs remain evidence, not target authority by inheritance;
- implementation remains blocked until D9;
- the exact next action is D4-B3 independent design/review from canonical D4-B1+B2 plus current official Sankhya evidence and real sanctioned-API measurements where needed.

If it cannot, the authority path is incomplete or contradictory.