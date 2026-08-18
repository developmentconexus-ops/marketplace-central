# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 NEXT / NOT YET OPENED**  
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
11. `docs/engineering/rebaseline/D4-R1-PUBLICATION-INPUT.md`
12. `docs/engineering/rebaseline/D5-API.md`
13. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
14. code, OpenAPI, schemas, tests and runtime only as current-state evidence when needed

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
  ├─ Original D4 Global Coherence + YAGNI / Overengineering / Future-Cost — COMPLETED / PASS
  └─ R1 Publication Input & Listing Authoring targeted amendment — ACCEPTED / CANONICAL
       ├─ Pre-D5-B2 whole-product Global Coherence — COMPLETED / PASS after adjudicated corrections
       └─ R1-G1 Mercado Livre initial publication × Availability — CLOSED / PASS-B
  ↓
D5 — API — OPEN / ACTIVE
  ├─ B1 Semantic API Model & Contract Laws — ACCEPTED / CANONICAL
  └─ B2 Product Operation / Resource Surface — NEXT / NOT YET OPENED
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
- freshness-for-use through owner-controlled provenance/time where material;
- accepted/rejected/pending/ambiguous capability semantics where applicable;
- no blind replay of ambiguous external effects;
- projections as rebuildable read state, never write authority;
- shared execution-safety mechanisms verify owner-issued proofs without owning business disposition/authorization;
- no generic Event Bus/Command Bus/Workflow engine/event sourcing/exactly-once/global-ordering/runtime topology choice.

### D4 — CLOSED / ACCEPTED AS A WHOLE

`D4-EXTERNAL-INTEGRATIONS.md` remains the accepted whole-stage D4 authority for B1/B2/B3/B4 and the original closure review. `D4-R1-PUBLICATION-INPUT.md` is a later accepted targeted amendment discovered during D5-B2 preparation; it extends D4 without creating a separate product architecture.

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

The B2 Installation Evidence Gate is **CLOSED / PASS**. The selected first Mercado Livre proof context is time-bound and remains subject to capability revalidation. First controlled marketplace writes remain D8 proofs.

D4-B3 accepted:

- the sanctioned Sankhya Gateway/API surface is sufficient for the currently claimed Product 1.0 business-system contract under explicit SourceInstance capability fences;
- Product remains `SourceInstance + native Product key`; company/location/control/cost/fiscal provider dimensions remain external evidence, not MPC ontology;
- bounded sanctioned reads never become arbitrary SQL/Oracle escape hatches;
- Party Resolution and Destination Realization remain distinct Materialization prerequisites and create no Customer/Party/Address master authority;
- Business Order Intent and Invoicing Intent remain MPC semantics while TOP/NUNOTA/status/choreography stay provider-local;
- Expected Tax is delegated to the sanctioned Sankhya fiscal engine under the proven stable binding; MPC does not duplicate the tax engine;
- consequential writes preserve Organization+SourceInstance, owning intent/correlation, acceptance/ambiguity, authoritative reread and no-blind-retry semantics;
- no generic ERP/workflow/customer framework and no Direct Oracle fallback is introduced.

D4-B4 accepted:

- **Semantic Core + Provider-Enriched Evidence** is the target: MPC does not collapse marketplaces to a lowest common denominator and does not mirror arbitrary provider payloads into business ontology;
- provider-specific evidence is retained when it serves a named Product 1.0 consumer/correctness property; unsupported equivalents on another provider remain honestly unsupported/not-applicable/unavailable/unknown;
- Market Intelligence retains competitive interpretation; Commercial Economics retains expected/order/realized economic meaning; Offering retains Price Intent;
- expected fee, expected seller shipping, Order fee, billed charges/rebates, Payment approval/release/refund and bank cash evidence remain distinct evidence rungs;
- broader account-movement population and R3 bank-side evidence remain bounded defers until a real consumer appears;
- report generation is not admitted as read support by convenience;
- no generic financial ledger, universal fee model, generic CollectorPort or unadjudicated scraping path is introduced.

Original final D4 Global Coherence accepted:

- no duplicate or missing business authority;
- Provider Richness preserves essential provider capability without provider overfit;
- Sankhya and marketplace/provider bindings remain replaceable realization rather than core ontology;
- D7/D8 defers are safe and trigger-bounded;
- YAGNI, future-cost, later-stage leakage and legacy-ADR coherence all pass.

Two original coherence fences remain binding:

1. **D4 evidence contract is not D4 evidence authority/store.** Persistent MPC semantic ownership follows the D1/D2 owner; technical caches/raw acquisition artifacts never become canonical business truth merely by persistence.
2. **Provider resource ownership does not move wholesale to one consumer.** One provider acquisition may translate into multiple consumer-owned semantic views/ports; no generic provider-resource/raw-payload entity bypasses D1/D3 authority edges.

#### D4-R1 — Publication Input & Listing Authoring — ACCEPTED / CANONICAL

The targeted amendment closes the publication-input seam without changing the 12 D1 boundaries:

- Product master remains external/source-qualified; no MPC Product/PIM master appears;
- **Readiness** owns publication requirements, correspondence, source candidates and source-level readiness;
- **Offering** owns `ListingIntent` as the single create/edit authoring/draft identity and owns draft dispatchability from current Readiness meaning;
- Listing values use only **`FOLLOW_SOURCE` or `EXPLICIT_OVERRIDE`** at baseline; no generic DERIVED/rule/mapping engine exists;
- humans and automation Principals author through the normal Semantic Product API; automation recurrence never silently reverses standing human overrides;
- source acquisition remains D4 evidence/mechanism feeding consumer-owned ports; no generic `SourceProductObservation` business owner exists;
- embedded source adapters do not require self-HTTP; external connector ingress remains a prepared seam until a real connector creates a wire-contract consumer;
- media may be source-qualified or MPC-authored for one ListingIntent without creating a Product-media master;
- provider requirement/schema churn remains source-qualified D4 evidence and historical ListingIntent context, not a universal ProductAttribute ontology;
- a provider request may **jointly realize multiple owner-issued meanings** through D4/D7 execution mechanics without creating a new semantic edge or merging business authority;
- `R1-G1 Mercado Livre initial publication × Availability = CLOSED / PASS-B`: Offering never owns quantity; Availability issues its own meaning/input; a technical execution mechanism may serialize both into one provider request and each owner later evaluates its own convergence;
- publication create/edit is allowed to be multi-step, partial and asynchronous; early provider `2xx/201/202` never implies whole-operation convergence;
- no D0/D1/D2/D3/D5-B1 reopen was required.

The first controlled real Mercado Livre creation remains a D8 proof, including authoritative reread and shared-User-Product blast-radius verification.

### D5 — OPEN / ACTIVE

`D5-API.md` is the current D5 authority.

D5-B1 **Semantic API Model & Contract Laws** is ACCEPTED / CANONICAL with `RESTRUCTURE NOW` relative to the current API shape.

B1 establishes:

- a semantic/domain-oriented MPC Product API, distinct from provider/business-system protocol ingress;
- HTTP/REST resource semantics where truthful plus explicit owner-specific operations where CRUD would lie;
- no generic Mutation/Action/Command/Operation business API owner;
- Organization-owned Product API operations under `/organizations/{organization_id}/...`;
- fail-closed same-Organization resolution for secondary body/query references;
- ordinary Principal access distinct from business disposition and Governance authorization;
- Q/C/P wire semantics preserving knowledge, freshness/provenance and effect meaning;
- source-qualified external identity at the wire boundary; no bare provider/native correlation key;
- accepted/rejected/pending/ambiguous outcomes where materially reachable;
- mandatory fail-closed idempotency key by default for consequential intake, with only explicit operation-local structural-idempotency exemptions;
- optimistic concurrency only where stale client state is materially unsafe;
- RFC 9457 Problem Details for API-level failures, distinct from valid domain outcomes;
- provider-rich, source-qualified enrichment without provider DTO/payload ontology;
- OpenAPI as the single machine-readable Product API wire authority, with supported client contracts derived/conformant and server behavior mechanically conformant;
- conformance controls must be shown to fire through negative drift fixtures;
- hard cutover with no compatibility/versioning tax absent a real consumer;
- bulk admitted only per real operation/workflow with member-level correctness;
- ADR-016 historical.

D5-B2 must incorporate D4-R1 explicitly. In particular, it must not introduce a Product/PIM API, separate PublicationPreparation resource, SourceProductObservation business API, generic mapping/rules API, AI-specific authoring surface or `createListing = success` semantics.

## 4. Engineering method and repo lifecycle

Engineering reasoning follows the **DevelopmentConexus Engineering Method** identified in `AGENTS.md`; the local file in this authority path is only the consumed context copy.

This router defines the Marketplace Central D0–D9 status/lifecycle and allowed work. It is repo-specific specialization, not a second organizational engineering method. Conflicts inside the organizational method's scope must be surfaced, never silently reinterpreted.

## 5. What is prohibited now

While **D5 — API** is OPEN / ACTIVE and B2 is next:

- do not begin D6–D9 target design before D5 is accepted as a whole;
- do not implement product architecture/features; implementation remains blocked until D9;
- do not silently alter accepted D0–D4, D4-R1 or D5-B1 authority;
- do not treat current OpenAPI/routes/SDK/controller/package shape as D5 target authority by inheritance;
- do not expose provider DTO/resource vocabulary as MPC API semantics when a D1-owned semantic contract exists;
- do not revive generic `/mutations`, `/commands`, provider-resource or integration-platform business surfaces;
- do not introduce an MPC Product/PIM master, `PublicationPreparation` aggregate, generic `SourceProductObservation` authority, generic provider-field bag or generic listing transformation/rules engine;
- do not create a separate AI/agent authority path; automation uses D2 Principal + normal Product API semantics;
- do not make Offering own Availability/Fulfillment meaning merely because provider protocol combines fields in one physical request;
- do not turn D5-B2 into a generic CRUD/API framework exercise disconnected from Product 1.0 consumers;
- do not make frontend or runtime topology decisions as side effects of API operation design;
- do not use projections/caches as consequential write or concurrency authority;
- do not weaken Organization path scope, same-Organization secondary-reference checks, source-qualified external identity, honest knowledge/freshness, authorization, idempotency, precondition, ambiguity or external-convergence semantics already accepted;
- do not add compatibility/versioning machinery without a real entitled consumer;
- do not treat `AI-DIALOG.md`, deleted review candidates or chat/reviewer summaries as target authority.

Existing OpenAPI/code/schema remains current-state evidence only.

## 6. Exact next action

**Open D5-B2 — Product Operation / Resource Surface from accepted D0–D4 + D4-R1 + D5-B1 authority.**

Do not begin from the legacy route list or from the existing pre-R1 B2 review candidate by inheritance. Re-derive the smallest coherent Product 1.0 external operation surface from accepted owners and real product consumers.

For each candidate Product API operation, establish proportionately:

- real Product 1.0 consumer/use;
- one accepted semantic owner / accepted D2 substrate authority;
- Q / C / P interaction class;
- Organization path scope and same-Organization secondary-reference rule;
- ordinary Permission requirement;
- canonical or source-qualified subject identity;
- read knowledge/freshness/provenance semantics where applicable;
- consequential Intent/outcome/idempotency/precondition/concurrency semantics where applicable;
- projection/read-only status for composed views;
- provider-enriched fields only where a named consumer/correctness need exists;
- pagination/filter/sort/cursor only where the real consumer requires them;
- bulk only where a real Product 1.0 workflow requires member-level bulk semantics.

For publication/listing operations specifically, B2 must preserve:

- Readiness requirements/source-level readiness versus Offering draft dispatchability;
- one Offering-owned `ListingIntent` create/edit authoring model;
- `FOLLOW_SOURCE | EXPLICIT_OVERRIDE` baseline resolution only;
- normal automation-Principal authoring and human-override supersession safety;
- no current generic external source-ingestion Product API; external connector ingress is admitted only when a real connector requires a wire contract;
- joint technical realization of owner-issued Listing/Availability/Fulfillment meanings without owner merger;
- multi-step/partial/asynchronous provider outcomes and owner-specific convergence.

B2 must not choose D6 screens or D7 runtime/generator framework as side effects. If an operation cannot fit accepted D1/D2/D3/D4/D4-R1/D5-B1 meaning without distortion, stop and reopen only the implicated parent decision.

Implementation remains blocked until D9.

## 7. Fresh-session success test

A fresh session should conclude that:

- D0, D1, D2 and D3 are **CLOSED / ACCEPTED**;
- D4 is **CLOSED / ACCEPTED AS A WHOLE** with B1/B2/B3/B4 accepted plus the later **D4-R1 Publication Input & Listing Authoring targeted amendment ACCEPTED / CANONICAL**;
- D4-R1 whole-product coherence found one coherent system and no missing/duplicate D1 boundary;
- `R1-G1 Mercado Livre initial publication × Availability = CLOSED / PASS-B` without D1/D3 reopen;
- Product master remains external; Readiness owns requirements/source-level readiness; Offering owns the one create/edit ListingIntent draft and draft dispatchability;
- publication values use `FOLLOW_SOURCE | EXPLICIT_OVERRIDE` at baseline; no PIM/rule engine/PublicationPreparation/SourceProductObservation authority exists;
- source acquisition feeds consumer-owned ports; embedded adapters need no self-HTTP and future external ingress remains consumer-triggered rather than speculative;
- provider requests may jointly realize multiple owner-issued meanings without business-authority merger;
- publication effects may be multi-step/partial/asynchronous and require owner-specific authoritative reread/convergence;
- consumer-owned semantics and provider-local protocol remain distinct;
- provider-rich evidence is retained when materially useful without becoming universal ontology or payload mirroring;
- Sankhya API Gateway remains the target transport and Direct Oracle is not an admitted fallback;
- **D5 — API is OPEN / ACTIVE**;
- **D5-B1 Semantic API Model & Contract Laws is ACCEPTED / CANONICAL**;
- the Product API is semantic/domain-oriented and Organization-scoped;
- source-qualified external identity, honest knowledge/freshness, effect ambiguity and consequential idempotency survive the wire;
- OpenAPI is the single machine-readable Product API wire authority; manual SDK duplication is not target architecture;
- ADR-016 is historical;
- hard cutover remains allowed because no production compatibility consumer exists;
- **D5-B2 Product Operation / Resource Surface is NEXT / NOT YET OPENED**;
- the existing pre-R1 D5-B2 review candidate is non-authoritative and must be revised/re-derived before any B2 ratification;
- implementation remains blocked until D9.

If it cannot, the authority path is incomplete or contradictory.
