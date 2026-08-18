# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D4 — EXTERNAL INTEGRATIONS — OPEN / ACTIVE; D4-B1 ACCEPTED / CANONICAL; D4-B2 ACCEPTED / CANONICAL; D4-B3 ACCEPTED / CANONICAL; D4-B4 ACCEPTED / CANONICAL; FINAL D4 GLOBAL COHERENCE NEXT**  
> **Implementation:** BLOCKED until D9 is accepted  
> **Last updated:** 2026-08-18

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
  ├─ B3 Sankhya Business-System Contract — ACCEPTED / CANONICAL
  ├─ B4 Market / Economics / Settlement Contract — ACCEPTED / CANONICAL
  │    ├─ M1 Market Evidence lane — CLOSED / PASS
  │    ├─ E1 Expected / Order Economic Evidence — CLOSED / PASS
  │    └─ S1 Realized / Release Evidence — CLOSED / PASS
  └─ Final D4 Global Coherence + YAGNI / Overengineering / Future-Cost — NEXT / NOT YET COMPLETED
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

### D4 — OPEN / B1+B2+B3+B4 ACCEPTED

`D4-EXTERNAL-INTEGRATIONS.md` is the active D4 authority.

D4-B1 accepted:

- concrete provider/business-system adapters implement consumer-owned semantic ports; no Integration business domain or universal provider entity graph;
- Marketplace Installation and Sankhya SourceInstance bind external namespaces explicitly and fail closed where authoritative markers expose mismatch;
- credentials/auth are protocol/runtime secrets, not business identity;
- notification is a trigger/pointer; current external meaning comes from authoritative reread where material;
- point/enumeration/delta/notification coverage claims are operation-scoped and fail-honest;
- Integration Support, Provider Effective Capability/Requirement and Effective Business Capability remain distinct authorities;
- later external effects distinguish acceptance/ambiguity from convergence and name authoritative reconciliation surfaces;
- **Sankhya API Gateway is the target transport; Direct Oracle/database is outside target architecture and is never an implicit fallback.**

D4-B2 accepted:

- Mercado Livre provider topology stays inside the adapter; Item, User Product, Family, Catalog Product, provider stock locations, Claim and Return do not become MPC business ontology by normalization;
- listing creation/observation is current-provider/mode aware while Offering retains Listing/Price intent and convergence;
- Product↔channel provider identifiers are evidence for Readiness, never automatic correspondence authority;
- price/availability writes preserve provider-effective prerequisites, shared-effect scope and authoritative reread/convergence;
- seller-managed does not automatically mean API-writable and provider-managed Full is not silently treated as MPC-controlled;
- seller Order search completion does not prove cancellation-inclusive Sales coverage;
- Order and Shipment remain separate provider resources;
- fulfillment responsibility, fiscal/label/readiness requirements and external SLA remain provider-context-sensitive evidence translated for Fulfillment closure;
- essential Claim/Return/reverse-shipment support stays bounded;
- no universal `OperatingMode`/provider graph/framework is introduced.

The B2 Installation Evidence Gate is **CLOSED / PASS**. The selected current first Mercado Livre proof context remains time-bound: User Product; non-multi-origin Item availability; direct-price candidate; seller-operated `me2 / xd_drop_off`; live Shipment SLA; real Claim/Return + reverse Shipment. B2 residual R1 — first controlled Price/Availability write + reread — remains D8. B2 residual R2 — live fiscal/label progression — uses the accepted B3 materialization contract and remains D8.

D4-B3 accepted:

- the sanctioned Sankhya Gateway/API surface is sufficient for the currently claimed Product 1.0 business-system contract under explicit SourceInstance capability fences;
- Product remains externally identified by `SourceInstance + native Product key`; company/location/control/cost/fiscal provider dimensions remain external evidence, not MPC business ontology;
- `CRUDServiceProvider.loadRecords` is admitted only through a bounded root-entity/read fence; arbitrary SQL/subqueries/cross-table criteria/Oracle escape hatches are not admitted;
- Business-System Party Resolution and Delivery Destination Realization are distinct bounded Materialization prerequisites; neither creates an MPC Customer/Party/Address master;
- ambiguous native party correspondence fails closed; a human-adjudicated resolution may be retained as bounded durable Materialization correspondence when provider reread cannot reconstruct the decision;
- concurrent/repeated first-time materializations for the same unresolved fiscal identity must not independently create duplicate native parties; D7 chooses the mechanism;
- transaction shipping evidence never silently authorizes customer-master overwrite or creation of another customer merely for another address; unsupported destination realization becomes explicit Work / `external-required`;
- the current Sankhya contact-based alternate-destination path is a conditioned concrete candidate whose first consequential proof remains D8;
- Business Order Intent maps to a bounded explicit Sankhya target binding; provider-native TOP/NUNOTA/status/choreography never becomes MPC semantics;
- Invoicing Intent remains readiness-gated and correlates to a distinct native fiscal result; first irreversible selected-lane `313→306` effect remains D8;
- Expected Tax uses `POST /v1/fiscal/impostos/calculo` through the stable model binding proven on 2026-08-18; MPC does not duplicate the Sankhya tax engine;
- negotiation type is fiscally material for Expected Tax; current proven binding preserves type `27`;
- provider-returned tax component/value/provenance is preserved honestly;
- consequential Sankhya writes preserve explicit Organization + SourceInstance, owning intent/correlation, acceptance/ambiguity, authoritative reread and no-blind-retry semantics;
- no generic ERP/workflow/customer framework and no Direct Oracle fallback is introduced.

D4-B4 accepted:

- **Semantic Core + Provider-Enriched Evidence** is the target: MPC does not collapse marketplaces to a lowest common denominator and does not mirror arbitrary provider payloads into business ontology;
- materially useful provider-specific evidence may be preserved when it serves a named Product 1.0 consumer/correctness property; unsupported equivalents on another provider remain honestly unsupported/not-applicable/unavailable/unknown rather than fabricated;
- Mercado Livre enriched Market Evidence may include `price_to_win`, catalog offer/winner evidence, buyer-facing shipping/free-shipping state, shipping tags and boosts/reasons while Market Intelligence retains comparability/competitive interpretation and Offering retains Price Intent;
- a real case proved price-only market comparison materially misleading: own `69.90 + 44.94` buyer-facing shipping vs winner `79.90 + 0`, with `price_to_win=26.75` and shipping/boost evidence;
- expected selling fee, expected seller shipping, Order transaction fee and realized seller Shipment cost remain distinct evidence classes;
- current `listing_prices` behavior requires explicit qualification/fail-open fences: HTTP 200 does not prove fields were consumed, listing type/category/response shape must be validated, provider currency is preserved and silent-ignore behavior is falsified proportionately;
- transaction fee granularity/decomposition is source-specific; live multi-unit evidence re-proved Order `sale_fee` per unit while Payment exposes multiple directionally distinct charge rows;
- Payment approval, money release, refund/reversal, withdrawal/payout and Bank Cash Receipt remain distinct rungs;
- the same bound Mercado Livre Installation token can read the selected Payment API; **no separate Mercado Pago credential is required for the selected lane**;
- `money_release_date` alone does not prove release, `fee_details` is not complete fee evidence, and `net_received_amount` is not post-refund realized authority;
- real refund-after-release evidence is appended rather than rewriting earlier release and can feed Commercial Economics plus Post-Sale Resolution without authority transfer;
- broader account-movement population remains a bounded safe defer until a real unanchored-movement/period-completeness consumer appears;
- report generation is not admitted as “read support” by convenience;
- no generic financial ledger, universal fee model, generic CollectorPort or unadjudicated scraping path is introduced.

No D0/D1/D2/D3 or D4-B1/B2/B3 reopen is required by B4.

## 4. Engineering method and repo lifecycle

Engineering reasoning follows the **DevelopmentConexus Engineering Method** identified in `AGENTS.md`; the local file in this authority path is only the consumed context copy.

This router defines the Marketplace Central D0–D9 status/lifecycle and allowed work. It is repo-specific specialization, not a second organizational engineering method. Conflicts inside the organizational method's scope must be surfaced, never silently reinterpreted.

## 5. What is prohibited now

While the **Final D4 Global Coherence + YAGNI / Overengineering / Future-Cost review** is the active next action:

- do not begin D5–D9 target design prematurely;
- do not implement product architecture/features;
- do not silently alter accepted D0–D4-B4 authority;
- do not create a generic provider/ERP/financial framework from repeated vocabulary alone;
- do not collapse provider-rich evidence to a lowest common denominator merely for cross-provider symmetry;
- do not mirror provider payloads or PII merely because an API exposes them;
- do not move Market Intelligence, Commercial Economics, Offering, Post-Sale, Fulfillment or Materialization meaning into D4 adapters;
- do not convert bounded D7/D8 defers into reasons to reopen B1–B4 without new material evidence;
- do not introduce Direct Oracle/database access for Sankhya;
- do not choose HTTP/frontend/runtime topology as a side effect of D4 closure;
- do not treat `AI-DIALOG.md`, deleted review candidates or reviewer/chat summaries as target authority.

Existing code/module/context/schema shape remains current-state evidence only.

## 6. Exact next action

**Run the Final D4 Global Coherence + YAGNI / Overengineering / Future-Cost review across canonical D4-B1+B2+B3+B4.**

The review must challenge proportionately:

- duplicate or missing authority across D1 consumers and D4 evidence contracts;
- contradictions between B1 generic boundary and B2/B3/B4 specializations;
- provider-overfit vs lowest-common-denominator pressure;
- whether Provider Richness preserves useful evidence without becoming payload/ontology mirroring;
- whether expected/order/billed/released/refunded evidence remains distinct without creating a generic finance ledger;
- whether Sankhya/Marketplace provider-specific bindings remain replaceable realization rather than core semantics;
- whether repeated transport/coverage/falsification/correlation mechanisms should remain local/shared mechanism without business authority;
- whether D7/D8 defers are safe and owned by the correct later stage;
- YAGNI / overengineering / foreseeable retrofit traps, including second-provider and second-business-system replacement tests.

If no material contradiction survives, close D4 and route the next session to **D5 — API**. Do not begin product implementation; implementation remains blocked until D9.

## 7. Fresh-session success test

A fresh session should conclude that:

- D0, D1, D2 and D3 are **CLOSED / ACCEPTED**;
- D4 is **OPEN / ACTIVE** pending only its final Global Coherence review;
- D4-B1 is **ACCEPTED / CANONICAL**;
- D4-B2 is **ACCEPTED / CANONICAL** and its Installation Evidence Gate is **CLOSED / PASS**;
- D4-B3 is **ACCEPTED / CANONICAL**;
- D4-B4 is **ACCEPTED / CANONICAL** and M1/E1/S1 are **CLOSED / PASS**;
- consumer-owned semantic ports and provider-local protocol remain distinct;
- provider-rich evidence is preserved when materially useful without becoming universal MPC ontology or payload mirroring;
- unsupported provider capabilities remain honest rather than suppressing richer capabilities available elsewhere;
- Market Intelligence owns competitive interpretation; Commercial Economics owns L0/L1/L2 attribution/reconciliation; Offering owns Price Intent; Post-Sale owns refund-consequence closure;
- the selected ML Payment path uses the bound Installation credential without requiring a separate Mercado Pago credential for the proven lane;
- Sankhya API Gateway remains the target transport and Direct Oracle is not an admitted fallback;
- Party Resolution and Destination Realization remain distinct bounded Materialization prerequisites;
- D7/D8 proof obligations do not reopen B1–B4 merely because the concrete future effect has not yet executed;
- implementation remains blocked until D9;
- the exact next action is the final D4 Global Coherence review, not D5 and not implementation.

If it cannot, the authority path is incomplete or contradictory.
