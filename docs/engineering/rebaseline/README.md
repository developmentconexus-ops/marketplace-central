# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D4 — EXTERNAL INTEGRATIONS — OPEN / ACTIVE; D4-B1 ACCEPTED / CANONICAL; D4-B2 OPEN / CANONICAL CONTRACT CORE; INSTALLATION EVIDENCE GATE OPEN**  
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
  ├─ B2 Mercado Livre Operational Contract — OPEN / CANONICAL CONTRACT CORE
  │    └─ Installation Evidence Gate — OPEN / BLOCKS B2 WHOLE-BATCH CLOSURE
  ├─ B3 Sankhya Business-System Contract — NOT YET OPENED
  ├─ B4 Market / Economics / Settlement Contract — NOT YET OPENED
  └─ Final Global Coherence + YAGNI / Overengineering / Future-Cost — NOT STARTED
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

### D4 — OPEN / B1 ACCEPTED / B2 CORE CANONICAL

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
- the former Oracle/godror target existed historically because the project did not yet have a known/usable sanctioned Sankhya API path; that historical reason no longer constrains the target;
- B3 must prove Gateway/API correctness, coverage and operational viability; if the sanctioned surface cannot satisfy a required Product 1.0 claim, B3 stops and returns to explicit operator/architecture adjudication rather than enabling database access implicitly;
- ADR-004 D4 plugin-framework meaning is superseded; ADR-010 polling-only D4 meaning is superseded while D7 runtime residue remains; ADR-006/007 are historical for target architecture after their transport-independent lessons were rehomed.

D4-B2 canonical contract core now establishes:

- Mercado Livre provider topology stays inside the adapter; Item, User Product, Family, Catalog Product, provider stock locations, Claim and Return do not become MPC business ontology by normalization;
- listing creation/observation is mode-aware and current-provider-authoritative while Offering retains Listing/Price intent/convergence;
- Product↔channel provider identifiers are evidence for Readiness, never automatic correspondence authority;
- price observation authority is separated from current write mechanism; Price Automation is a Provider Effective restriction and provider 2xx never substitutes for authoritative convergence proof;
- User Product shared-field effects, including applicable availability writes, must not silently widen intended/authorized scope;
- stock writability depends on concrete provider resource **plus site, seller configuration and current resource/listing context**; seller-managed does not automatically mean API-writable, and provider-managed Full stock is not seller/MPC-writable by convenience;
- stale provider stock version conflict is a definitive rejected precondition followed by reread/redecision, not ambiguous blind retry;
- Order and Shipment remain separate provider resources; seller Order-search coverage explicitly excludes canceled Orders under current behavior and therefore cannot claim complete Sales/cancellation coverage;
- fulfillment responsibility, provider fiscal/label/readiness requirements and external SLA/deadline remain provider-context-sensitive evidence translated for Fulfillment closure;
- essential Claim/Return/reverse-shipment capability is included without expanding into CRM/SAC; refund/payment/settlement movement authority remains B4;
- Full/provider-operated observation alone does not prove MPC-controlled Availability convergence or internally operated Fulfillment Node execution;
- no universal `OperatingMode`/provider graph/framework and no support for every documented ML mode is introduced.

The **Installation Evidence Gate remains OPEN**. Public/provider documentation proves conditional contracts but does not prove which current modes/lane capabilities the real Metal Nobre Installation actually exposes. B2 therefore is not accepted as a whole yet.

No D0/D1/D2/D3 reopen is currently required.

## 4. Engineering method and repo lifecycle

Engineering reasoning follows the **DevelopmentConexus Engineering Method** identified in `AGENTS.md`; the local file in this authority path is only the consumed context copy.

This router defines the Marketplace Central D0–D9 status/lifecycle and allowed work. It is repo-specific specialization, not a second organizational engineering method. Conflicts inside the organizational method's scope must be surfaced, never silently reinterpreted.

## 5. What is prohibited now

While the D4-B2 Installation Evidence Gate is the active next work:

- do not begin B3/B4 or D5–D9 target design prematurely under the current router;
- do not implement product architecture/features;
- do not silently alter accepted D0–D4-B2 canonical contract authority;
- do not let provider APIs, ERP schema, current adapters, legacy DTOs or historical ADRs become target business authority by inheritance;
- do not create semantic dependencies outside D1 or bypass D2/D3 semantics through integration code;
- do not promote a provider callback/notification/2xx into domain truth or convergence without the accepted D4 evidence/reread/reconciliation contract;
- do not infer a seller-writable ML stock path from typology ownership alone; site/seller/resource enablement must be established;
- do not claim seller Order search covers canceled Orders or use search completion as complete cancellation coverage;
- do not infer current seller tags, User Product migration state, Price Automation, stock mode or fulfillment lane from historical probes;
- do not introduce Direct Oracle/database access for Sankhya as an implementation shortcut or API fallback;
- do not choose HTTP/frontend/runtime topology prematurely;
- do not treat `AI-DIALOG.md`, deleted review candidates or reviewer/chat summaries as target authority.

Existing code/module/context/schema shape remains current-state evidence only.

## 6. Exact next action

**Execute/admit D4-B2 §4.8 — the read-only Mercado Livre Installation Evidence Gate.**

The probe must use current real provider evidence without writes, secret leakage or buyer-PII retention and establish only the minimum facts needed to select the smallest Product 1.0 proof lane set:

- seller tags relevant to publication/stock model (`user_product_seller`, `warehouse_management`, `multiwarehouse` where applicable);
- selected listing topology: legacy/User Product, Item↔UP/Catalog relations and real variation/composite presence where material;
- Price Automation state for candidate price-write proof listings;
- stock typology **and concrete write surface actually enabled for this MLB seller/site/listing**, including `/items.available_quantity` applicability and UP/shared effect scope;
- recent selected Order/Shipment fulfillment/logistics/fiscal/label/SLA contexts;
- material current seller/listing moderation/restriction evidence;
- Claim/Return evidence if available, without interpreting sample absence as provider nonexistence.

Gate outcome must identify the smallest real lane set that can truthfully prove the accepted D0 capabilities. If no real seller-writable availability lane or required seller-operated fulfillment lane exists for a claimed Product 1.0 proof, surface the targeted product/proof conflict rather than fabricating support.

The canonical B2 mode-conditional contract core does not need to be re-reviewed merely because the probe selects one mode over another. A probe result that contradicts a contract assumption reopens only the implicated B2/upstream decision.

Under the current router, do not open B3/B4 until this B2 gate is resolved or the operator explicitly changes sequencing.

Implementation remains blocked until D9.

## 7. Fresh-session success test

A fresh session should conclude that:

- D0, D1, D2 and D3 are **CLOSED / ACCEPTED**;
- D4 is **OPEN / ACTIVE**;
- D4-B1 is **ACCEPTED / CANONICAL**;
- D4-B2 has a **CANONICAL CONTRACT CORE but remains OPEN because the Installation Evidence Gate is unresolved**;
- B3/B4 are **NOT YET OPENED** under the current router;
- consumer-owned ports and provider-local protocol remain distinct;
- Marketplace Installation/SourceInstance external namespace binding is explicit and fail-closed;
- Mercado Livre Item/UP/stock/Order/Shipment/Claim topology remains provider-local evidence rather than MPC business ontology;
- seller-managed stock is not assumed API-writable without site/seller/resource enablement evidence;
- UP shared effects cannot silently widen Listing/Availability intent scope;
- seller Order-search completion does not prove canceled-Order coverage;
- Full/provider-operated lanes do not prove MPC-controlled availability or internal physical fulfillment responsibilities they do not own;
- **Sankhya API Gateway is the target transport and Direct Oracle is not an admitted fallback;**
- current modules/contexts/schema/legacy ADRs remain evidence, not target authority by inheritance;
- implementation remains blocked until D9;
- the exact next action is the read-only real Mercado Livre Installation Evidence Gate.

If it cannot, the authority path is incomplete or contradictory.