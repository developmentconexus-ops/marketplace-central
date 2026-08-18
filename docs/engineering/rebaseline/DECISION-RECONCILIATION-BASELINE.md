# Decision Reconciliation Baseline

> **Status:** ACCEPTED / CANONICAL  
> **Purpose:** always-read decision-generation/routing index for the accepted D0→D4 + D4-R1 architecture and accepted D5-B1 contract laws before D5-B2 binds them to concrete API/code-facing contracts.  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Accepted:** 2026-08-18  
> **Operator ratification:** 2026-08-18

---

# 1. Role

This baseline answers one question:

> **Which decision generation is current, where is its semantic authority, what older idea must not be implemented, what remains genuinely open/deferred, and which invariants technical stages may not silently re-decide?**

It is an authority-routing artifact, not a second semantic architecture.

It exists because D0→D4 were intentionally iterative and D4-R1 was discovered only while D5-B2 was being derived. Without reconciliation, a later implementation/review could resurrect a stale ADR, stale candidate, legacy module shape or earlier formulation simply because it still existed somewhere in the repository.

---

# 2. Scope-based authority model

This baseline owns only **decision-generation reconciliation/routing**:

- which generation of a decision is current;
- which accepted artifact owns the detailed meaning;
- which older architectural idea routes to which current semantic home;
- which questions remain intentionally open and which later stage owns them.

It does **not** own detailed product/domain semantics. It does **not** own ADR file status.

Scope ownership:

1. `docs/engineering/rebaseline/README.md` — current program status and exact next action;
2. `ARCHITECTURE.md` — stable cross-stage constraints;
3. this baseline — current decision-generation routing map;
4. `docs/architecture/decisions/README.md` — sole ADR file status/disposition authority;
5. accepted D-stage artifacts — detailed semantic authority in their stage scope;
6. Evidence Register, code, schemas, OpenAPI, tests and Git history — supporting evidence unless active authority explicitly carries their meaning.

If this baseline disagrees with an accepted semantic home, **this baseline is defective and must shrink/correct; it never overrides the semantic authority**.

Absence from this baseline is never permission. `ARCHITECTURE.md` and accepted D-stage artifacts remain the complete constraint set.

---

# 3. Current decision generation

## 3.1 Product / system boundary

**Semantic home:** D0 + `ARCHITECTURE.md`.

Current generation:

- MPC is the internal **Marketplace Operations Control Plane + Commercial Intelligence** product;
- external systems retain authority for facts/processes inherently theirs;
- MPC owns cross-system marketplace operating semantics needed to observe → understand → reconcile → decide/policy → execute → verify → audit/reconcile;
- Mercado Livre is first marketplace proof; Sankhya is first business-system proof;
- Organization remains explicit;
- MPC is not an ERP replacement, Product/PIM/MDM master, generic marketplace/integration platform, universal workflow engine or unrestricted autonomous-agent platform;
- enterprise-generic repositioning requires D0/D1 reopen, never implementation generalization.

Do not resurrect dashboard-only, ERP-shaped, generic-provider-hub or MPC-Product-master interpretations.

## 3.2 Business authority boundaries

**Semantic home:** D1.

Current semantic authorities:

1. Marketplace Portfolio
2. Product & Channel Readiness
3. Marketplace Offering Operations
4. Availability Control
5. Market Intelligence
6. Commercial Economics
7. Controlled Action Governance
8. Marketplace Sales
9. Business-System Materialization
10. Fulfillment Lifecycle
11. Post-Sale Resolution
12. Operational Work

Plus the accepted non-domain D2 identity/access substrate.

These are semantic authorities, **not** a required number of services, databases or processes.

Load-bearing laws:

- one business meaning → one semantic authority;
- mechanism ≠ authority;
- Product master remains external;
- no generic Integration/Evidence/Mutation/Workflow/SLA/Policy business domain;
- Governance authorization ≠ ordinary Permission ≠ domain disposition;
- Work lifecycle ≠ originating business truth;
- Offering Listing/Price meaning ≠ Sellable Availability;
- Economics interpretation ≠ marketplace write;
- Materialization ≠ physical Fulfillment;
- provider resource grouping never merges MPC business authorities.

Legacy module/package names do not define current boundaries.

## 3.3 Identity / tenant / data ownership

**Semantic home:** D2.

Current generation:

- Organization is tenant/isolation root;
- MPC-owned canonical IDs are stable opaque identities;
- Marketplace Installation and SourceInstance qualify external namespaces without becoming credentials/protocol;
- Selling Entity, Inventory Source and Fulfillment Node retain distinct accepted meanings;
- Principal is accountable actor identity with human/automation/system distinctions;
- Product remains `SourceInstance + native Product key`; no MPC Product mirror/master;
- provider Listing/Variation, Sale/Order, Shipment and native financial movements remain source-qualified external identities;
- material Business Intents remain owner-local identities; no universal Action/Command/Intent owner;
- Membership / AccessRole / Permission / RoleAssignment stay in the D2 ordinary-access substrate; consequential authorization stays Governance-owned;
- Organization-owned durable state/evidence is explicitly Organization-scoped;
- unknown/empty/absent/not-applicable distinctions and material provenance/time remain honest;
- historical snapshots explain past decisions without becoming current authority;
- automation recurrence never silently reverses a standing human decision in the same semantic scope;
- pre-rebaseline persistence imposes no compatibility/migration model on the target.

Bare native identity, separate canonical Tenant, Product mirror and universal entity/evidence graphs are not current target structures.

## 3.4 Communication / failure / recovery

**Semantic home:** D3.

Current grammar:

- **Q** — current owner meaning;
- **C** — request owner capability/work;
- **E** — committed producer occurrence with an independent consumer reaction;
- **P** — read-only composition across authorities.

Current safety laws include:

- communication never transfers authority;
- Organization scope remains explicit and durable communication/recovery state can recover scope without ambient guessing;
- duplicate delivery is expected; semantic idempotency prevents duplicate business effect;
- no global ordering or exactly-once authority;
- known/known-empty/unknown/unavailable plus material partial/freshness/provenance remain honest;
- `accepted != completed != externally applied != converged`;
- accepted/rejected/pending/ambiguous are preserved where materially reachable;
- no blind replay after possible external acceptance;
- projections never become write/current-truth/concurrency authority;
- cross-owner workflows are correlated/convergent; no cross-owner atomicity is invented;
- **recoverable consequential propagation:** a required reaction cannot be lost into a silent permanent stall;
- **evidence-edge occurrence recoverability:** material historical occurrences remain recoverable from the smallest sufficient durable authority; latest state cannot erase them;
- cutover cannot silently discard pending reactions required for accepted Product 1.0 progression;
- a source-committed material actionable condition ends represented in Work or explicitly reconciled;
- reusable external-effect safety machinery verifies owner-issued proofs/correlation without owning business validity or authorization.

Transport/runtime realization remains D7.

## 3.5 External integration boundary

**Semantic home:** D4-B1/B2/B3/B4 + D4-R1.

Governing rule: **consumer owns meaning; adapter owns protocol**.

Current generation:

- D4 owns acquisition/protocol/capability/coverage/effect translation, not a business domain or generic evidence store;
- one provider acquisition may feed several consumer-owned semantic ports without D4 or one consumer owning the provider payload wholesale;
- credentials/auth are mechanism, not identity/business truth;
- notification/callback is acquisition evidence/pointer; authoritative reread establishes material current provider meaning;
- point/enumeration/delta/notification coverage is operation-scoped;
- Integration Support / Provider Effective Capability / Effective Business Capability remain distinct;
- external effects preserve owner intent/correlation, prerequisites, acceptance/ambiguity and authoritative reread/convergence;
- provider 2xx never proves business convergence by itself;
- provider-specific richness is kept for named consumers/correctness needs without lowest-common-denominator flattening or raw DTO mirroring;
- Sankhya target transport is the sanctioned API Gateway; Direct Oracle/database is not target fallback;
- provider-native Mercado Livre and Sankhya topology remains adapter-local realization, not MPC ontology;
- source admissibility is explicit: missing evidence never authorizes fabrication or an unadjudicated scraping/source path;
- provider PII is minimized.

Time-bound provider facts and concrete protocol surfaces remain in D4; this baseline does not restate them.

## 3.6 Publication / listing authoring seam

**Semantic home:** D4-R1 + D1/D2/D3 parents.

Current generation:

- external Product remains external/source-qualified;
- Readiness owns publication requirements, Product↔channel correspondence, source candidates and **source-level readiness**;
- Offering owns `ListingIntent` as the one **create/edit authoring/draft identity** and owns **draft dispatchability** from current Readiness meaning;
- no separate `PublicationPreparation` aggregate;
- listing-value resolution is **`FOLLOW_SOURCE | EXPLICIT_OVERRIDE`** at baseline only;
- no generic DERIVED/rules/mapping DSL;
- source acquisition remains D4 mechanism/evidence feeding consumer-owned ports; no generic `SourceProductObservation` business owner;
- embedded source adapters do not require self-HTTP; external connector ingress gets a wire contract only when a real connector creates a consumer;
- human and automation authoring use the same Semantic Product API authority; automation cannot impersonate source truth or silently reverse human overrides;
- media may be source-qualified or ListingIntent-specific MPC authoring without creating Product-media master;
- provider requirement/schema churn remains source-qualified D4 evidence and historical intent context, not universal ProductAttribute ontology;
- provider execution may jointly realize multiple owner-issued meanings in one physical request without merging business ownership;
- Mercado Livre initial publication × Availability is **PASS-B**: Offering never owns quantity; Availability issues its own meaning/input; technical execution may serialize both and each owner independently evaluates convergence;
- publication create/edit may be multi-step, partial and asynchronous; no `createListing = success` simplification.

Do not resurrect PublicationPreparation, SourceProductObservation owner/service, generic source-ingestion API, ProductAsset master, mapping/rules engine, AI-specific API or separate create/edit architectures.

## 3.7 API laws already accepted before B2

**Semantic home:** D5-B1 + `ARCHITECTURE.md`.

B2 may not re-decide:

- semantic/domain-oriented Product API distinct from protocol ingress;
- Organization-owned Product API under `/organizations/{organization_id}/...`;
- same-Organization resolution for secondary references;
- Principal/auth context is not a client-authored business field;
- Permission/invoke access remains distinct from domain disposition/Governance;
- source-qualified external identity on the wire; no bare native correlation key;
- honest knowledge/freshness/provenance on reads;
- accepted/rejected/pending/ambiguous business outcomes distinct from transport problems;
- fail-closed idempotency key by default for consequential intake unless an operation proves structural/natural idempotency;
- optimistic concurrency only where stale client state is materially unsafe;
- RFC 9457 Problem Details for API-level failures;
- one machine-readable Product API wire authority: OpenAPI;
- supported clients derive/conform and server behavior conforms; conformance controls must be shown to fire;
- hard cutover/no compatibility tax absent an entitled consumer;
- bulk only when a real operation/workflow proves member-level semantics.

---

# 4. Historical idea → current semantic home

This table is semantic routing only. **ADR status is owned solely by the ADR registry.**

| Older architectural idea | Current semantic home / instruction |
|---|---|
| Marketplace dashboard as product | D0 control-plane product loop; do not implement dashboard-only architecture |
| Product mirror / MPC Product master | D2 external Product identity + D1 Readiness consumers; do not create MPC Product master |
| separate canonical Tenant | D2 Organization root |
| generic Integration business domain | Portfolio business meaning + D4 protocol + D7 mechanics |
| provider plugin/self-registration business framework | concrete D4 adapters; shared technical mechanism only when later proven |
| Direct Oracle/godror Sankhya target | D4 sanctioned Gateway only |
| `SELLER_SKU == CODPROD` | D2/D4 evidence + Readiness correspondence; not identity law |
| CODPROD+EAN unattended link formula | D2 corroboration safety + Readiness policy |
| generic Mutation business owner | domain-local intents + Governance + D3/D7 execution safety |
| generic divergence ledger as truth owner | source-domain correctness + Work lifecycle |
| `sync`/polling phase as product semantics | domain freshness/coverage + D4/D7 mechanics |
| provider DTO/resource model in core | consumer-owned semantic ports + adapter-local protocol |
| lowest-common-denominator marketplace model | semantic core + provider-enriched evidence |
| generic CollectorPort/market-source framework | explicit source admissibility + later D7 mechanism if proven |
| generic Customer/Address master for Sankhya | bounded Materialization Party/Destination realization |
| MPC tax engine | sanctioned Sankhya fiscal engine + Economics interpretation |
| global Fee/Payment/Settlement MPC entity | source-qualified external movements + Economics attribution/reconciliation |
| PublicationPreparation aggregate | Offering ListingIntent draft + Readiness Q |
| SourceProductObservation business service | D4 acquisition feeding consumer-owned ports |
| generic listing transformation/rules engine | `FOLLOW_SOURCE | EXPLICIT_OVERRIDE`; targeted future reopen only on repeated need |
| dedicated AI business/API authority | D2 automation Principal + ordinary Product API/Governance |
| Product API listing quantity owned by Offering | Availability-owned meaning; joint technical realization only |
| manual OpenAPI + manual SDK authority | OpenAPI one wire authority + conformant/derived clients/server |
| compatibility/versioning for legacy API | hard cutover until a real entitled consumer appears |

---

# 5. What remains genuinely open

## D5-B2

- exact Product 1.0 operation/resource inventory;
- request/response schemas;
- Permission→operation mapping;
- concrete paths/nouns under B1 laws;
- pagination/filter/sort/cursor only where a real consumer proves need;
- operation-local bulk only where proven;
- concrete OpenAPI/server/client generation and conformance tooling.

## D6

- screens/navigation/editor topology;
- projection/view composition required by UX;
- frontend feature/package structure.

## D7

- process/server/worker/scheduler topology;
- transaction/outbox/queue/cursor/lease/lock mechanisms;
- token refresh/cache/secret realization;
- RLS/runtime Organization-isolation realization;
- idempotency persistence/TTL/locking;
- retry/backoff/rate-control mechanisms;
- media blob/cache/CDN realization;
- production deployment topology;
- retained runtime ADR residues named by the ADR registry.

## D8

- first controlled Mercado Livre create/Price/Availability effects + authoritative reread/convergence;
- shared-User-Product blast-radius real-write proof;
- selected irreversible Sankhya fiscal progression;
- controlled alternate destination/contact realization before claiming it;
- first consequential native-party write if the golden flow reaches it;
- unproven provider/fiscal branches only when a selected golden flow materially depends on them.

## D9

- final adversarial architecture/system review;
- implementation remains blocked until D9 acceptance.

## Bounded Unknown / Deferred classes

Detailed homes remain D4/D4-R1. Current classes include unselected Mercado Livre modes/configurations, paused/zero-quantity representation-first creation, broader payment/account movement universe, R3 bank-side reconciliation, controlled-product marketplace paths, post-invoice return/reversal branches, unproven fiscal components and future reusable source→listing transformations/media library/external connector wire contract.

Unknown/Deferred is **never** permission to invent a plausible default.

---

# 6. Implementation reconciliation guard

The following list highlights high-risk invariants; **it is not exhaustive**. `ARCHITECTURE.md` and accepted D-stage artifacts remain the complete constraint set. Absence from this index is never permission.

Technical stages MUST NOT silently re-decide by convenience:

1. Product ownership — external/source-qualified, not MPC master.
2. Organization root/isolation semantics.
3. D1 business authorities and accepted semantic edges.
4. Q/C/E/P meaning, recoverable consequential propagation and evidence-occurrence recovery.
5. domain-local Intent ownership versus generic mutation/workflow authority.
6. Governance versus Permission versus domain disposition.
7. Work lifecycle versus originating business truth; material actionable conditions cannot become ownerless silent state.
8. external/provider identity versus MPC canonical identity.
9. consumer-owned port / adapter-owned protocol boundary.
10. sanctioned Sankhya Gateway-only target transport.
11. no blind retry after ambiguous possible external acceptance.
12. provider richness without provider DTO mirroring.
13. Provider PII minimization.
14. explicit source admissibility; missing data never authorizes fabricated/unadjudicated source acquisition.
15. externally governed obligation/policy provenance; MPC cannot silently relax external obligations.
16. Offering / Readiness / Availability publication ownership split.
17. ListingIntent as the one create/edit authoring identity.
18. `FOLLOW_SOURCE | EXPLICIT_OVERRIDE` baseline and standing-human-decision safety.
19. no hidden PIM/source-observation/rule/connector/AI business framework.
20. joint technical realization may compose owner-issued meanings but may not create ownership transfer or hidden semantic edges.
21. owner-specific convergence after provider effects.
22. execution-time validity/authorization: an earlier approval does not authorize blind execution after material governing drift and cannot waive mandatory safety invariants.
23. material multi-target separation: intended scope, authorized scope and attempted/outcome scope remain distinct; member-level confirmed/rejected/ambiguous/not-executed outcomes survive.
24. cutover/recovery cannot silently discard pending reactions required for accepted lifecycle progression.
25. semantic Product API with explicit Organization scope and source-qualified identity.
26. OpenAPI as the one machine-readable Product API wire authority; no second manual SDK/wire authority.
27. hard cutover while no entitled production compatibility consumer exists.

A technical design that appears to require violating one of these is evidence for a targeted architecture reopen — not permission for a workaround.

---

# 7. ADR / active-tree reconciliation

The ADR registry is the **only active authority for ADR file status/disposition**.

The 2026-08-18 reconciliation retired all fully rehomed pre-rebaseline ADRs from the active tree. The retained legacy residue set is owned and enumerated only by the registry. Git history is the archive; ADR numbers are never reused.

`ADR-035` remains transition authority through D0–D9, but its embedded 2026-08-14 still-binding/reopened tables are historical snapshot evidence only; current ADR disposition is owned by the registry and later accepted D-stage authority.

The citation-harvest tree is retained only for citation files still directly referenced by retained legacy residues. No `old/` or in-repo archive tree exists.

---

# 8. Fresh-session read shape

Always read:

```text
AGENTS.md
→ rebaseline README/router
→ DevelopmentConexus Method
→ ARCHITECTURE.md
→ Decision Reconciliation Baseline
→ ADR registry (only unresolved technical/transition residues)
→ accepted/current D-stage artifact(s) needed for the work
→ Evidence Register
→ implementation evidence only when necessary
```

A fresh reviewer must not need `AI-DIALOG.md`, stale candidates, retired ADRs or legacy implementation to discover current target architecture.

---

# 9. Reopen triggers

Reconcile this baseline only when its routing map materially changes because:

- an accepted D-stage decision is amended/reopened;
- a new target ADR changes current decision routing;
- later real proof invalidates an architectural assumption;
- a Product requirement changes D0/D1 ownership;
- a retained legacy ADR is adjudicated and leaves the active tree;
- D0–D9 closes and ADR-035/transition machinery retires.

Do not reopen for implementation naming, package layout, framework preference or rediscovery of a retired Git-history decision.

---

# 10. Reconciliation verdict

**CURRENT D0→D4/D4-R1 + D5-B1 DECISION SET RECONCILED / COHERENT.**

Independent review and GPT adjudication found:

- no contradiction among accepted D0→D4/D4-R1 + D5-B1;
- no missing business authority/boundary;
- no D6/D7 mechanism frozen by this index;
- no material architecture prerequisite blocks D5-B2;
- Structural Inversion passes: deleting/inverting current implementation/OpenAPI does not change the reconciled architecture.

**Exact program status/next action remains owned only by the router.**