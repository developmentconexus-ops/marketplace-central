# Marketplace Central Architecture

> **Status:** stable product-level constraints during Architecture Rebaseline  
> **Detailed target architecture:** intentionally under D0–D9 design; see `docs/README.md`  
> **Last updated:** 2026-08-18

## Purpose

This file contains only architecture constraints that remain stable while the detailed system is re-adjudicated.

It is deliberately **not** a catalog of current modules, tables, routes, frontend packages or processes. Those are current-state evidence in D0 and target-design questions in D1–D9.

A new structural decision becomes durable only after the relevant D-stage is accepted and, when appropriate, an ADR records it.

## Product North Star

Marketplace Central is an internal operations and intelligence system for marketplace commerce, initially Mercado Livre, backed by real Sankhya business-system operational facts.

It must support trustworthy flows for:

- internal product identity and source observations;
- marketplace listings/variations and channel observations;
- product↔listing linkage with explicit evidence;
- stock/cost/tax/price semantics;
- safe external actions with preview/policy/audit/reconciliation;
- marketplace orders and realized profitability;
- operational views whose freshness/completeness is honest.

## Stable platform constraints

1. **Independent monorepo.** Marketplace Central remains its own repository/application. Future integration into a broader product does not justify coupling current domain boundaries to another repository.
2. **Go backend is canonical business execution.** `apps/server_core` is the server-side application. Business policy is not duplicated in React.
3. **React frontend is a client, not a second domain authority.** Server state is managed through TanStack Query (ADR-021) unless a material finding explicitly reopens that decision.
4. **PostgreSQL stores MPC-owned canonical state.** External systems are sources/dependencies, not alternate writable application stores.
5. **Sankhya is external to MPC and is integrated through its sanctioned API Gateway.** Business-system access stays behind MPC-owned consumer ports/adapters. Direct Oracle/database access is explicitly outside the target architecture and is not a fallback path. If the Gateway/API surface cannot satisfy a materially required Product 1.0 claim, the architecture stops and re-adjudicates the requirement/transport explicitly; it does not fall back to Oracle by convenience. Legacy Oracle/godror code remains current-state/history evidence only.
6. **Mercado Livre first.** Other marketplace providers are deferred until the Mercado Livre operating loop is coherent and the adapter protocol is proven (ADR-005).
7. **Marketplace provider boundary.** Provider integrations enter through vendor adapters and implement ports owned by consuming business contexts (ADR-033). Provider wire DTO/protocol knowledge stays inside the vendor boundary; exact target package layout remains later realization detail unless a stage explicitly freezes it.
8. **Honest absence.** Unknown facts do not become plausible zero/default values. `internal/kernel/fact` is an accepted primitive for uncertainty where semantically appropriate (ADR-034); D2 decides its correct scope rather than forcing it onto every value.
9. **Exactness where the domain requires it.** Money/tax/cost/pricing values must not lose correctness through floating-point convenience. D2 owns the exact shared/domain representation.
10. **Tenant-ready data isolation is a real invariant.** The exact tenant runtime/RLS model is under D7, but tenant isolation may not depend solely on developers remembering predicates.
11. **External writes are controlled.** A provider write has explicit authority/policy, duplicate protection, auditability and reconciliation. An ambiguous outcome is not blindly retried (ADR-029).
12. **Provider PII is minimized.** Raw external PII is not retained merely because a payload contains it (ADR-025).
13. **Partial observations are honest.** Absence from a partial provider/source pull does not prove closure/deletion (ADR-027).
14. **No compatibility tax without a consumer.** There are no production users requiring current route/schema/package compatibility; hard cutover is allowed under ADR-035.
15. **Git history is history.** Active source/document trees do not keep `old/` copies or parallel legacy roadmaps.
16. **Semantic Product API.** MPC clients use a semantic/domain-oriented Product API, not provider/integration ontology. Organization-owned Product API operations are path-scoped under `/organizations/{organization_id}/...`; provider/business-system protocol ingress remains a separate D4 boundary. OpenAPI is the single machine-readable Product API wire authority; supported client contracts derive/conform to it and server behavior conforms to the same contract. A hand-written second wire authority is not target architecture.
17. **Publication authoring is not Product mastery.** External Product identity/truth remains source-qualified; Product & Channel Readiness owns publication requirements/correspondence/source-level readiness; Marketplace Offering Operations owns `ListingIntent` as the one create/edit authoring identity and owns draft dispatchability. Baseline listing-value resolution is `FOLLOW_SOURCE | EXPLICIT_OVERRIDE` only. Provider/source acquisition remains D4 mechanism feeding consumer-owned ports; no MPC PIM, `PublicationPreparation`, generic `SourceProductObservation`, rule/mapping engine, connector platform or AI-specific authority path is target architecture. A provider request may jointly realize multiple owner-issued meanings through technical execution machinery without merging business ownership.

## Architecture Rebaseline authority

ADR-035 governs detailed target-design authority during D0–D9.

The following are **not currently frozen target decisions** even if old code/ADRs once specified them:

- exact context/module set beyond accepted D1 semantic authorities;
- exact database schemas/tables/FKs and physical persistence realization;
- exact scheduler/process/worker topology;
- exact HTTP path namespace and operation set beyond D5-B1's Organization path-scoping law;
- generated server/client technology choice;
- frontend feature/package topology;
- exact transaction/outbox implementation;
- legacy `connectors`, `integrations`, `marketplaces`, `mutations`, `sync`, `internal_read`, `dashboard` or other module structures;
- old manual SDK or proxy-table synchronization mechanisms.

The legacy direct-Oracle/godror Sankhya shape is no longer merely “unfrozen”: **D4-B1 explicitly supersedes it as target architecture.**

The legacy manual OpenAPI+SDK same-commit shape is likewise no longer open target structure: **D5-B1 supersedes ADR-016 and requires one machine-readable Product API wire authority.**

A future session must not infer target architecture from existing artifacts before their responsible D-stage accepts the relevant meaning.

## Target reasoning shape

The rebaseline is testing — not blindly accepting — a top-level shape with:

```text
apps/server_core/internal/
  contexts/       business authorities
  adapters/       external-system translations
  kernel/         tiny shared value semantics only
  platform/       technical runtime mechanisms without business policy
  composition/    final assembly only
  views/          rebuildable read projections when justified
```

The exact context set and allowed semantic edges are D1 outcomes. Runtime/package realization remains subject to later stages.

## Communication principles

D3 is accepted and canonical:

- **Q:** current producer-owned meaning required for the consumer's current decision;
- **C:** caller asks the owner to accept/perform owner-owned work;
- **E:** already-committed producer-owned fact with a real independent consumer reaction;
- **P:** multiple authorities composed for reading without becoming write authority.

Communication may duplicate, arrive late/out of order, fail or replay without changing business truth. Cross-context SQL/private implementation access is not unnamed communication. D7 chooses transport/runtime realization without moving authority.

## External-integration principles

D4-B1, D4-B2, D4-B3, D4-B4 and the later accepted **D4-R1 Publication Input & Listing Authoring amendment** are canonical. External integrations obey these constraints:

- consumer context owns semantic meaning/port;
- adapter owns provider/business-system protocol, DTOs, auth and pagination;
- Marketplace Installation/SourceInstance qualify the correct external namespace without becoming credentials or provider IDs;
- namespace mismatch fails closed where authoritative source/provider markers allow it to be detected;
- provider notification/callback is acquisition evidence; authoritative reread establishes current provider state where material;
- point/enumeration/delta/notification coverage is operation-scoped; incomplete/unavailable never becomes plausible absence;
- Integration Support and Provider Effective Capability/Requirement do not become Effective Business Capability;
- external-effect contracts distinguish acceptance/ambiguity from convergence and name an authoritative reread/reconciliation surface;
- Mercado Livre Item/User Product/Catalog/stock/Order/Shipment/Claim topology stays provider-local and does not become MPC business ontology by normalization;
- source Product acquisition may feed multiple consumer-owned semantic ports but does not create a generic source-Product business owner/store; embedded adapters need not loop through a public HTTP API merely for symmetry;
- publication requirements remain provider-authoritative evidence translated for Readiness; they do not create a universal ProductAttribute/ProviderRequirement business framework;
- Listing creation/editing uses one Offering-owned `ListingIntent` authoring model; Readiness owns source-level requirement/correspondence meaning and Offering owns draft dispatchability from that meaning;
- `FOLLOW_SOURCE` preserves external source authority/provenance; `EXPLICIT_OVERRIDE` is MPC-authored listing intent state with Principal attribution; recurring automation never silently reverses a standing human override;
- listing media may originate from source evidence or MPC authoring for the listing context without creating an MPC Product-media master;
- a provider call that physically requires several owner-owned values may be one **joint technical realization** of those owner meanings; the execution mechanism correlates owner-issued inputs but cannot create a hidden `Availability → Offering`/`Fulfillment → Offering` business edge, recalculate another owner's answer or merge convergence semantics;
- current Mercado Livre initial User Product publication × Availability is closed as **PASS-B**: active creation may require an Availability-issued quantity in the same provider call; Offering never owns that quantity and Availability evaluates its own convergence after authoritative reread;
- publication create/edit may be multi-step, partial and asynchronous; early `2xx/201/202` never proves whole-listing convergence, and shared-User-Product blast radius remains explicit;
- provider stock/price/fulfillment capability is context-sensitive; seller-managed does not automatically mean API-writable, provider-managed Full is not silently treated as MPC-controlled, and shared User Product effects cannot silently widen intended/authorized scope;
- seller Order-search completion does not prove cancellation-inclusive Sales coverage when provider documentation/runtime behavior does not establish a reliable complete universe;
- the first current Mercado Livre proof lane selected by D4-B2 is deliberately narrow and time-bound; future provider configuration changes re-evaluate capability rather than creating speculative universal mode support;
- Sankhya uses the sanctioned API Gateway target path; Direct Oracle is not an admitted fallback;
- Sankhya Product/company/location/inventory/control/cost/tax/party/document concepts remain provider-local evidence/realization and never become MPC business ontology by normalization;
- Business-System Party Resolution and Delivery Destination Realization are distinct bounded Materialization prerequisites; neither creates an MPC Customer/Party/Address master;
- transaction-specific delivery evidence never silently authorizes customer-master overwrite or another customer record merely to represent another destination; unsupported destination realization remains explicit Work / `external-required`;
- Business Order Intent and Invoicing Intent remain MPC-owned semantics while TOP, NUNOTA, provider statuses and native choreography remain adapter-local;
- Expected Tax is delegated to the sanctioned Sankhya fiscal engine under an explicit, revalidatable SourceInstance binding; MPC does not reimplement the provider tax engine and does not turn unproven/absent tax components into plausible zeros;
- provider-native negotiation/configuration may be fiscally material and therefore participates in concrete binding validation without becoming MPC commercial ontology;
- consequential business-system effects preserve source-qualified correlation, authoritative reread and no-blind-retry semantics;
- **provider-rich, semantics-first:** MPC does not discard materially useful provider evidence merely because another marketplace lacks an equivalent capability; shared semantics are normalized only where meanings genuinely align, while provider-distinct evidence remains source-qualified/optional and never becomes universal MPC ontology merely because one provider exposes it;
- capability richness does not authorize payload mirroring: provider evidence is retained only for a named consumer/correctness need or materially required non-reobservable evidence, with PII minimization intact;
- Market Evidence such as provider competition status, `price_to_win`, winner/offer shipping, free-shipping tags or boosts may enrich Market Intelligence when exposed; they never become Price Intent or automatic price recommendation;
- expected sale fee, expected seller shipping, Order transaction fee, billed charge/rebate, Payment approval, release/account impact, refund/reversal, withdrawal/payout and Bank Cash Receipt remain distinct evidence/meaning rungs;
- source-specific fee/financial decomposition and granularity are preserved rather than forced into one generic `Fee` model;
- current provider read surfaces that may silently ignore/fallback on request qualifiers require fail-honest validation/falsification proportional to the economic claim; transport 200 alone is not semantic sufficiency;
- a separate payment-provider credential is not invented when the bound Marketplace Installation credential already proves the selected sanctioned Payment read path; capability remains Installation/context-sensitive and revalidated on material change;
- absence of an admitted external market-data source never authorizes fabricated evidence or an unadjudicated scraping path by convenience; a materially new market-data source requires explicit source, legality/trust, coverage and provenance adjudication before its evidence can support MPC claims;
- current unstable provider/reference behavior must be verified against current official/real behavior for the concrete decision that depends on it;
- live integration claims require real-dependency evidence, not only mocks;
- no speculative universal provider/integration/ERP/workflow/customer/financial framework is introduced.

Installation-/SourceInstance-specific proof-lane details, B4 Payment/fee observations, D4-R1 provider-documentation evidence and later D7/D8 proof obligations remain in canonical D4 artifacts rather than becoming stable global platform constants here.

## API and frontend

D5-B1 is accepted and canonical.

The Product API obeys these stable constraints:

- client-facing vocabulary follows MPC semantic owners, not Mercado Livre/Sankhya/provider DTO/resource nouns;
- provider/business-system protocol ingress is a separate D4 boundary and is not part of the normal Product SDK;
- Organization-owned Product API operations are scoped under `/organizations/{organization_id}/...`;
- secondary Organization-owned references in body/query must resolve inside that path Organization;
- provider/native identifiers are source-qualified through Marketplace Installation / SourceInstance or an unambiguous operation scope; bare external IDs are never Product API correlation keys;
- Q results preserve known/known-empty/unknown/unavailable/partial and owner-controlled freshness/provenance where material;
- consequential C outcomes preserve accepted/rejected/pending/ambiguous and never collapse acceptance into completion/application/convergence;
- consequential intake requires a fail-closed idempotency key by default, with only explicit operation-local structural-idempotency exemptions;
- RFC 9457 Problem Details represents API-level failures; valid business/domain outcomes do not become transport problems by convenience;
- provider-rich evidence may appear only as source-qualified, owner-bounded enrichment for a named Product 1.0 need; raw provider payload mirroring and lowest-common-denominator flattening are both rejected;
- OpenAPI is the single machine-readable Product API wire authority; supported SDK/client contracts derive/conform to it and server behavior conforms to the same contract;
- contract/conformance controls must be shown to fire through negative fixtures;
- no current API compatibility/versioning machinery is required because there is no entitled production client;
- bulk is operation-local and admitted only for a real consumer/workflow with member-level correctness.

The exact operation inventory, final path nouns beyond Organization scoping, request/response schemas, Permission mapping, pagination/filter/sort/bulk decisions and concrete generator/server technology remain D5/later realization questions.

D6 maps every target screen to explicit accepted API/query/capability ownership and decides frontend feature/package topology. D6 may not invent a second client-side business authority.

## Runtime and persistence

D7 decides serving/worker/scheduler/outbox/transaction/cursor/secret/deployment topology for the admitted target transports. Do not preserve a dedicated executable, poller or database driver because it exists today.

D2 classifies persistent state and the clean target baseline. Historical migrations do not automatically define the target model.

## Proof bar

A structural rule should, where reasonable, fail at the strongest available boundary:

- illegal private import → compile failure;
- invalid value combination → type/constructor/schema failure;
- foreign/unowned write → structurally unavailable or mechanically blocked;
- contract drift → generation/validation red;
- missing required consequential idempotency key → API rejection before durable intake;
- same key reused for materially different semantic request → explicit API problem;
- cross-Organization secondary resource reference → fail closed before semantic resolution/effect;
- bare external identity collision across Installation/SourceInstance → schema/SDK negative fixture red;
- materially stale known response without required owner provenance → contract review/fixture red;
- RLS/isolation bypass → boot/integration failure;
- custom guard → negative fixture proves it fires;
- external namespace mismatch → fail closed before attribution/effect where the source exposes authoritative qualification;
- partial acquisition → cannot pass as complete in contract/integration proof;
- attempted Direct Oracle wiring for Sankhya target integration → architecture/governance failure, not accepted fallback;
- listing publication missing a required owner-issued Availability/other-owner input → fail closed before provider dispatch rather than copy/default the value into Offering;
- automation recurrence attempting to silently replace a standing human listing override → reject or require explicit supersession semantics with preserved lineage;
- provider-rich evidence missing on another provider → honest unsupported/not-applicable/unavailable, never suppression of a richer supported provider or fabricated equivalence;
- economic/provider request field silently ignored/fallbacked → cannot pass as sufficiently qualified evidence merely because transport returned 2xx.

A green artifact that did not execute the relevant subject is no proof.

## Current stage

Read `docs/README.md` for the sole current status and exact next action.
