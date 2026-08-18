# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — NEXT / NOT YET OPENED; D4 CLOSED / ACCEPTED AS A WHOLE**  
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
D4 — External Integrations — CLOSED / ACCEPTED AS A WHOLE
  ├─ B1 External Contract Grounding — ACCEPTED / CANONICAL
  ├─ B2 Mercado Livre Operational Contract — ACCEPTED / CANONICAL
  │    └─ Installation Evidence Gate — CLOSED / PASS
  ├─ B3 Sankhya Business-System Contract — ACCEPTED / CANONICAL
  ├─ B4 Market / Economics / Settlement Contract — ACCEPTED / CANONICAL
  │    ├─ M1 Market Evidence lane — CLOSED / PASS
  │    ├─ E1 Expected / Order Economic Evidence — CLOSED / PASS
  │    └─ S1 Realized / Release Evidence — CLOSED / PASS
  └─ Final Global Coherence + YAGNI / Overengineering / Future-Cost — COMPLETED / PASS
  ↓
D5 — API — NEXT / NOT YET OPENED
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

`D1-DOMAINS-BOUNDARIES.md` is the accepted D1 authority. It defines 12 business boundaries, explicit ownership/non-ownership, legal semantic edges, cross-cutting non-domain treatment, forbidden boundary violations, legacy semantic disposition and reopen triggers.

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
- current truth from owners and material historical occurrence from the smallest sufficient durable authority;
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

### D4 — CLOSED / ACCEPTED AS A WHOLE

`D4-EXTERNAL-INTEGRATIONS.md` is the accepted D4 authority.

D4-B1 accepted:

- concrete provider/business-system adapters implement consumer-owned semantic ports; no Integration business domain or universal provider entity graph;
- Marketplace Installation and Sankhya SourceInstance bind external namespaces explicitly and fail closed where authoritative markers expose mismatch;
- credentials/auth are protocol/runtime secrets, not business identity;
- notification is a trigger/pointer; authoritative reread establishes current external meaning where material;
- point/enumeration/delta/notification coverage claims are operation-scoped and fail-honest;
- Integration Support, Provider Effective Capability/Requirement and Effective Business Capability remain distinct authorities;
- external effects distinguish acceptance/ambiguity from convergence and name authoritative reconciliation surfaces;
- **Sankhya API Gateway is the target transport; Direct Oracle/database is outside target architecture and is never an implicit fallback.**

D4-B2 accepted:

- Mercado Livre provider topology stays inside the adapter; Item, User Product, Family, Catalog Product, provider stock locations, Claim and Return do not become MPC business ontology by normalization;
- Offering retains Listing/Price Intent and convergence; Availability retains Sellable Availability authority;
- Product↔channel provider identifiers are evidence for Readiness, never automatic correspondence authority;
- price/availability writes preserve provider-effective prerequisites, blast radius and authoritative reread/convergence;
- seller-managed does not automatically mean API-writable and provider-managed Full is not silently treated as MPC-controlled;
- seller Order search completion does not prove cancellation-inclusive Sales coverage;
- Order and Shipment remain separate provider resources;
- fulfillment responsibility, provider requirements and SLA remain context-sensitive external evidence translated for Fulfillment closure;
- essential Claim/Return/reverse-shipment support stays bounded;
- no universal OperatingMode/provider graph/framework is introduced.

The B2 Installation Evidence Gate is **CLOSED / PASS**. The selected first Mercado Livre proof context is time-bound and remains subject to capability revalidation. First controlled Price/Availability write+reread and selected fiscal/label progression remain D8 proofs.

D4-B3 accepted:

- the sanctioned Sankhya Gateway/API surface is sufficient for the currently claimed Product 1.0 business-system contract under explicit SourceInstance capability fences;
- Product remains `SourceInstance + native Product key`; company/location/control/cost/fiscal provider dimensions remain external evidence, not MPC ontology;
- `CRUDServiceProvider.loadRecords` is admitted only through a bounded root-entity/read fence; arbitrary SQL/Oracle escape hatches are not admitted;
- Party Resolution and Destination Realization are distinct bounded Materialization prerequisites and create no Customer/Party/Address master authority;
- ambiguous correspondence fails closed; safe human-adjudicated correspondence may be retained where reread cannot reconstruct the decision;
- destination evidence never silently authorizes customer-master overwrite or duplicate customer creation;
- Business Order Intent and Invoicing Intent remain MPC semantics while TOP/NUNOTA/status/choreography stay provider-local;
- Expected Tax is delegated to the sanctioned Sankhya fiscal engine under the proven stable binding; MPC does not duplicate the tax engine;
- consequential writes preserve Organization+SourceInstance, owning intent/correlation, acceptance/ambiguity, authoritative reread and no-blind-retry semantics;
- no generic ERP/workflow/customer framework and no Direct Oracle fallback is introduced.

D4-B4 accepted:

- **Semantic Core + Provider-Enriched Evidence** is the target: MPC does not collapse marketplaces to a lowest common denominator and does not mirror arbitrary provider payloads into business ontology;
- provider-specific evidence is retained when it serves a named Product 1.0 consumer/correctness property; unsupported equivalents on another provider remain honestly unsupported/not-applicable/unavailable/unknown;
- Mercado Livre enriched Market Evidence may include `price_to_win`, catalog offer/winner evidence, buyer-facing shipping/free-shipping state, shipping tags and boosts/reasons while Market Intelligence owns competitive interpretation and Offering owns Price Intent;
- a real price/shipping/winner case proved that price-only competitive comparison is materially misleading;
- expected selling fee, expected seller shipping, Order transaction fee and realized seller Shipment cost remain distinct evidence classes;
- `listing_prices` requires explicit qualification/fail-open fences and proportional falsification; HTTP 200 does not prove a submitted field was consumed;
- fee granularity/decomposition is source-specific;
- Payment approval, money release, refund/reversal, withdrawal/payout and Bank Cash Receipt remain distinct rungs;
- the same bound Mercado Livre Installation token can read the selected Payment API; **no separate Mercado Pago credential is required for the selected lane**;
- `money_release_date` alone does not prove release, `fee_details` is not complete fee evidence, and `net_received_amount` is not post-refund realized authority;
- real refund-after-release evidence is appended rather than rewriting earlier release and can feed Commercial Economics plus Post-Sale Resolution without authority transfer;
- broader account-movement population remains a bounded safe defer until a real consumer appears;
- report generation is not admitted as read support by convenience;
- no generic financial ledger, universal fee model, generic CollectorPort or unadjudicated scraping path is introduced.

Final D4 Global Coherence accepted:

- **CURRENT STRUCTURE CONFIRMED / PASS**;
- no duplicate or missing business authority;
- B2/B3/B4 specialize B1 without weakening its namespace/coverage/effect-safety rules;
- Provider Richness preserves essential provider capability without provider overfit;
- economic rungs remain distinct without a finance ledger;
- Sankhya and marketplace/provider bindings remain replaceable realization rather than core ontology;
- D7/D8 defers are safe and trigger-bounded;
- YAGNI, future-cost, later-stage leakage and legacy-ADR coherence all pass;
- no additional D4 batch and no D0–D4 reopen is required.

Two final coherence fences are binding:

1. **D4 evidence contract is not D4 evidence authority/store.** D4 preserves enough evidence across the boundary; persistent MPC semantic ownership follows the D1/D2 owner. Technical caches/raw acquisition artifacts, if later justified, are D7 mechanism and never canonical business truth.
2. **Provider resource ownership does not move wholesale to one consumer.** One provider acquisition may translate into multiple consumer-owned semantic views/ports. No consumer owns the provider payload as a whole, and no generic provider-resource/raw-payload entity bypasses D1/D3 authority edges.

## 4. Engineering method and repo lifecycle

Engineering reasoning follows the **DevelopmentConexus Engineering Method** identified in `AGENTS.md`; the local file in this authority path is only the consumed context copy.

This router defines the Marketplace Central D0–D9 status/lifecycle and allowed work. It is repo-specific specialization, not a second organizational engineering method. Conflicts inside the organizational method's scope must be surfaced, never silently reinterpreted.

## 5. What is prohibited now

While **D5 — API** is the next unopened stage:

- do not begin D6–D9 target design before D5 is accepted;
- do not implement product architecture/features; implementation remains blocked until D9;
- do not silently alter accepted D0–D4 authority;
- do not treat current OpenAPI/routes/SDK/controller/package shape as D5 target authority by inheritance;
- do not expose provider DTO/resource vocabulary as MPC API semantics when a D1-owned semantic contract exists;
- do not turn D5 into a generic CRUD/API framework exercise disconnected from Product 1.0 consumers;
- do not make frontend or runtime topology decisions as side effects of API design;
- do not use projections/caches as consequential write authority;
- do not weaken known/unknown/unavailable/partial, authorization, idempotency, precondition, ambiguity or external-convergence semantics already accepted in D0–D4;
- do not treat `AI-DIALOG.md`, deleted review candidates or chat/reviewer summaries as target authority.

Existing OpenAPI/code/schema remains current-state evidence only.

## 6. Exact next action

**Open D5 — API from accepted D0–D4 authority.**

D5 must decide the smallest coherent API contract by which clients interact with the accepted business authorities, proportionately including:

- public query/capability/mutation operations mapped to D1 owners and D3 Q/C/P semantics;
- stable resource/operation naming in MPC business language rather than provider DTO vocabulary;
- Organization and Principal/access scoping at the API boundary;
- known / known-empty / unknown / unavailable / partial representation where material;
- pagination/filter/sort/cursor contracts only where real consumers need them;
- command/precondition/idempotency/concurrency/duplicate semantics for consequential actions;
- accepted/rejected/pending/ambiguous outcomes where the owner/external effect can reach them;
- bulk/member-level partial outcome semantics only where a real Product 1.0 workflow requires bulk;
- validation/error/problem semantics that preserve domain/provider distinction without leaking raw provider errors as business truth;
- OpenAPI/generation/SDK authority and drift-prevention rules;
- ordinary access Permission→API-operation enforcement while keeping business disposition and consequential authorization with their accepted owners;
- API cuts required by provider-rich Market Intelligence and economic views without forcing lowest-common-denominator provider fields into universal contracts.

D5 must **not** choose frontend component topology (D6), worker/queue/outbox/retry/transaction/deployment topology (D7), golden-flow execution/proof choreography (D8), or product implementation.

Implementation remains blocked until D9.

## 7. Fresh-session success test

A fresh session should conclude that:

- D0, D1, D2, D3 and D4 are **CLOSED / ACCEPTED**;
- D4-B1/B2/B3/B4 are **ACCEPTED / CANONICAL**;
- B2 Installation gate and B4 M1/E1/S1 are **CLOSED / PASS**;
- Final D4 Global Coherence is **COMPLETED / PASS** with no earlier-stage reopen;
- consumer-owned semantics and provider-local protocol remain distinct;
- provider-rich evidence is retained when materially useful without becoming universal ontology or payload mirroring;
- D4 does not own a generic persistent evidence store;
- one provider resource may feed multiple consumer-owned semantic views without transferring whole-resource authority;
- Sankhya API Gateway remains the target transport and Direct Oracle is not an admitted fallback;
- Market Intelligence owns competitive interpretation, Commercial Economics owns L0/L1/L2 attribution/reconciliation, Offering owns Price Intent, Post-Sale owns refund-consequence closure, and Materialization/Fulfillment retain their accepted meanings;
- D7/D8 proof obligations do not reopen D4 merely because their future concrete effect has not yet executed;
- **D5 — API is NEXT / NOT YET OPENED**;
- implementation remains blocked until D9;
- the exact next action is to open D5 from D0–D4 authority, not to implement or jump to frontend/runtime design.

If it cannot, the authority path is incomplete or contradictory.
