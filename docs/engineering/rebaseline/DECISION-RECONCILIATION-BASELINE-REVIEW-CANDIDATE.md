# Decision Reconciliation Baseline — CONVERGED REVIEW CANDIDATE

> **Status:** NON-AUTHORITATIVE — Fable reviewed / GPT adjudicated / awaiting explicit operator ratification  
> **Purpose:** reconcile the accepted D0→D4 + D4-R1 architecture and accepted D5-B1 contract laws before D5-B2 turns them into concrete API/code-facing contracts.  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Review evidence:** `AI-DIALOG.md`, Fable Decision Reconciliation Baseline Independent Review, 2026-08-18  
> **Important:** this file is a proposed **decision-generation/routing index**, not a second semantic architecture and not an ADR-status authority.

---

# 1. Why this baseline exists

D0→D4 were intentionally iterative. Later decisions refined, split or superseded earlier/legacy choices. D4-R1 was discovered only when D5-B2 attempted to derive the concrete Product API surface. D5-B1 is already accepted; D5-B2 will bind architecture to operations, schemas, Permissions and implementation-facing contracts.

The repository therefore needs one compact reconciliation point answering:

> **Which decision generation is current, where is its semantic authority, what older idea must not be implemented, what remains genuinely open/deferred, and which invariants technical stages may not silently re-decide?**

This is an authority-routing problem, not a new product-design problem.

The baseline exists to prevent four recurring defect classes:

1. resurrecting an older decision because its file still exists;
2. implementing both an old and a new model “to be safe”;
3. treating an Unknown/Deferred item as permission to invent a default;
4. letting implementation convenience silently reopen architecture.

---

# 2. Scope-based authority model

If accepted, this baseline owns only **decision-generation reconciliation/routing**:

- which generation of a decision is current;
- which accepted artifact owns the detailed meaning;
- which historical idea routes to which current semantic home;
- which questions remain intentionally open and which later stage owns them.

It does **not** own detailed product/domain semantics already defined by D0–D4/D4-R1/D5-B1.

It does **not** own legacy ADR file status. `docs/architecture/decisions/README.md` remains the **sole ADR-status/disposition authority**.

Scope ownership:

1. `docs/engineering/rebaseline/README.md` — current program status and exact next action;
2. `ARCHITECTURE.md` — stable cross-stage constraints;
3. this baseline — current decision-generation routing map;
4. ADR registry — sole status/disposition index for ADR files;
5. accepted D-stage artifacts — detailed semantic authority in their stage scope;
6. Evidence Register / code / schemas / OpenAPI / tests / Git history — supporting evidence unless an active authority explicitly carries their meaning.

If this baseline ever disagrees with an accepted semantic home, **the baseline is defective and must shrink/correct; it does not override the semantic authority**.

If a future target ADR changes a decision, ADR status changes in the ADR registry and the semantic home changes in the responsible authority; this baseline is updated only if its routing map is materially affected.

---

# 3. Current decision generation

This section is intentionally index-sized. Detailed rules remain in the named authority.

## 3.1 Product/system boundary — CURRENT

**Semantic home:** D0 + stable `ARCHITECTURE.md` constraints.

Current generation:

- MPC is the internal **Marketplace Operations Control Plane + Commercial Intelligence** product;
- external systems retain authority for facts/processes inherently theirs;
- MPC owns the cross-system operating semantics required to observe → understand → reconcile → decide/policy → execute → verify → audit/reconcile;
- Mercado Livre is first marketplace proof; Sankhya is first business-system proof;
- Organization remains explicit;
- MPC is not an ERP replacement, Product/PIM/MDM master, generic marketplace/integration platform, universal workflow engine or unrestricted autonomous-agent platform;
- enterprise-generic repositioning requires D0/D1 reopen, never implementation generalization.

Historical ideas such as dashboard-only product, ERP-shaped product, generic provider hub and MPC Product/PIM do not belong to the current generation.

---

## 3.2 Business authority boundaries — CURRENT

**Semantic home:** D1.

The current semantic authorities are:

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

Load-bearing boundary rules include:

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

Legacy package/module names do not define current boundaries.

---

## 3.3 Identity / tenant / data ownership — CURRENT

**Semantic home:** D2.

Current generation:

- Organization is tenant/isolation root;
- MPC-owned canonical IDs are stable opaque identities;
- Marketplace Installation and SourceInstance qualify external namespaces without becoming credentials/protocol;
- Selling Entity, Inventory Source and Fulfillment Node retain their accepted distinct meanings;
- Principal is accountable actor identity, with human/automation/system distinctions;
- Product remains `SourceInstance + native Product key`; no MPC Product mirror/master;
- provider Listing/Variation, Sale/Order, Shipment and native financial movements remain source-qualified external identities;
- material Business Intents remain owner-local identities; no universal Action/Command/Intent owner;
- Membership / AccessRole / Permission / RoleAssignment stay in the D2 ordinary-access substrate; consequential authorization stays Governance-owned;
- Organization-owned durable state/evidence is explicitly Organization-scoped;
- unknown/empty/absent/not-applicable distinctions and material provenance/time remain honest;
- historical snapshots explain past decisions without becoming current authority;
- automation recurrence never silently reverses standing human decisions in the same semantic scope;
- pre-rebaseline persistence imposes no compatibility/migration model on the target.

Bare native identity, separate canonical Tenant, Product mirror and universal entity/evidence graphs are not current target structures.

---

## 3.4 Communication / failure / recovery — CURRENT

**Semantic home:** D3.

Current grammar:

- **Q** — current owner meaning;
- **C** — request owner capability/work;
- **E** — committed producer occurrence with an independent consumer reaction;
- **P** — read-only composition across authorities.

Current safety laws include:

- communication never transfers authority;
- Organization scope remains explicit and durable communication/recovery state can recover that scope without ambient guessing;
- duplicate delivery is expected; semantic idempotency prevents duplicate business effects;
- no global ordering or exactly-once authority;
- known/known-empty/unknown/unavailable and material partial/freshness/provenance remain honest;
- `accepted != completed != externally applied != converged`;
- accepted/rejected/pending/ambiguous are preserved where materially reachable;
- no blind replay after possible external acceptance;
- projections never become write/current-truth/concurrency authority;
- cross-owner workflows are correlated/convergent; no cross-owner atomicity is invented;
- **recoverable consequential propagation:** a required reaction cannot be lost into a silent permanent stall;
- **evidence-edge occurrence recoverability:** material historical occurrences remain recoverable from the smallest sufficient durable authority; latest state cannot erase them;
- cutover cannot silently discard pending reactions whose accepted Product 1.0 progression still depends on them;
- a source-committed material actionable condition ends represented in Work or explicitly reconciled; ownerless actionable work is not an accepted terminal state;
- reusable external-effect safety machinery verifies owner-issued proofs/correlation without owning business validity or authorization.

Transport/runtime mechanism remains D7.

---

## 3.5 External integration boundary — CURRENT

**Semantic home:** D4-B1/B2/B3/B4 + D4-R1.

Current generation:

- **consumer owns meaning; adapter owns protocol**;
- D4 owns acquisition/protocol/capability/coverage/effect translation, not a business domain or generic evidence store;
- one external acquisition may feed multiple consumer-owned semantic ports without one consumer/D4 owning the provider resource wholesale;
- credentials/auth are mechanism, not identity/business truth;
- notifications/callbacks are acquisition evidence/pointers; authoritative reread establishes material current provider meaning;
- coverage is operation-scoped;
- Integration Support, Provider Effective Capability and Effective Business Capability remain distinct;
- consequential external effects preserve owner intent/correlation, prerequisites, acceptance/ambiguity, blast radius and authoritative reread/convergence;
- provider 2xx never proves business convergence by itself;
- provider richness is preserved for named correctness/consumer needs without lowest-common-denominator flattening or raw DTO mirroring;
- Sankhya target transport is the sanctioned Gateway; Direct Oracle/database is not an admitted fallback;
- provider/business-system topology remains adapter-local realization; detailed current Mercado Livre/Sankhya surfaces and time-bound proof-lane facts remain in D4, not in this index;
- source admissibility is explicit: absence of an admitted market-data source never authorizes fabricated evidence or an unadjudicated scraping source;
- externally governed requirements/policy evidence never silently becomes editable MPC policy; internal rules/targets may not relax external obligations.

Detailed provider realization is deliberately not duplicated here.

---

## 3.6 Publication / listing authoring seam — CURRENT

**Semantic home:** D4-R1 + parent D1/D2/D3 authority.

Current generation:

- external Product remains external/source-qualified;
- Readiness owns publication requirements, correspondence, source candidates and **source-level readiness**;
- Offering owns `ListingIntent` as the one create/edit authoring identity and owns **draft dispatchability** using Readiness meaning via the accepted Readiness→Offering Q edge;
- there is no separate `PublicationPreparation` aggregate;
- baseline resolution is `FOLLOW_SOURCE | EXPLICIT_OVERRIDE` only;
- no generic DERIVED/rules/mapping DSL;
- source acquisition is D4 mechanism feeding consumer-owned ports; there is no generic `SourceProductObservation` business owner;
- embedded adapters need no self-HTTP; a future physically external connector gets a bounded ingress contract only when a concrete consumer exists;
- human and automation authoring use the Semantic Product API; automation cannot impersonate source truth or silently reverse a human override;
- listing-specific media does not create Product-media master authority;
- provider requirement/schema evolution remains source-qualified evidence/history, not universal ProductAttribute ontology;
- one provider request may jointly realize multiple owner-issued meanings without merging ownership;
- initial Mercado Livre publication × Availability is **PASS-B** under D4-R1: Offering does not own quantity; Availability issues its own meaning/input; technical execution may serialize both; convergence remains owner-specific;
- create/edit may be multi-step/partial/asynchronous; there is no `createListing = success` shortcut.

PublicationPreparation, SourceProductObservation owner, generic listing rule engine, Product media master and AI-specific authority path are not current target structures.

---

## 3.7 API laws already frozen before B2 — CURRENT

**Semantic home:** D5-B1 + `ARCHITECTURE.md`.

B2 must preserve:

- semantic/domain-oriented Product API distinct from protocol ingress;
- Organization-owned Product API paths under `/organizations/{organization_id}/...`;
- same-Organization validation of secondary references;
- Principal/access context distinct from client-authored business fields;
- Permission/invoke access distinct from business disposition/Governance;
- source-qualified external identities on the wire;
- honest knowledge/freshness/provenance;
- business outcomes distinct from API transport problems;
- fail-closed idempotency by default for consequential intake unless operation-local structural/natural idempotency is proven;
- concurrency protection only where stale client state is materially unsafe;
- RFC 9457 Problem Details for API-level failures;
- one machine-readable Product API wire authority: OpenAPI;
- derived/conformant clients and conformant server behavior with negative proof that conformance controls fire;
- hard cutover/no compatibility tax absent an entitled consumer;
- bulk only when a real operation/workflow proves member-level need.

Legacy routes, provider resource vocabulary, generic mutation endpoints and manual OpenAPI+manual SDK dual authority do not define the target.

---

# 4. Historical idea → current semantic home

This table is **routing only**. It does not assign ADR file status; ADR status belongs exclusively to the ADR registry.

| Historical/earlier idea | Current semantic route |
|---|---|
| marketplace dashboard as product | D0 control-plane product loop |
| MPC Product mirror/master | D2 external Product identity + D1 Readiness/Offering consumers |
| separate canonical Tenant | D2 Organization root |
| generic Integration business domain | Portfolio business meaning + D4 protocol + D7 mechanics |
| provider plugin/self-registration framework | concrete D4 adapters; shared mechanics only when later proven |
| Direct Oracle/godror Sankhya target | D4 sanctioned Gateway-only target transport |
| `SELLER_SKU == CODPROD` | D2 source-qualified identity + D4 evidence + Readiness correspondence |
| CODPROD+EAN universal auto-link rule | Readiness-owned corroboration policy under D2 safety law |
| generic Mutation business owner | owner-local intents + Governance + D3/D7 execution safety |
| generic divergence ledger as truth owner | source-domain correctness + Operational Work lifecycle |
| sync/polling phase as product semantics | owner freshness/coverage + D4 acquisition + D7 runtime |
| provider DTO/resource model in core | consumer-owned semantics + adapter-local protocol |
| lowest-common-denominator marketplace model | semantic core + provider-enriched evidence |
| generic market CollectorPort framework | explicit source admissibility + later D7 mechanics only if proven |
| generic Customer/Address master for Sankhya | bounded Materialization Party/Destination realization |
| MPC tax engine | sanctioned Sankhya fiscal engine + Economics interpretation |
| global Fee/Payment/Settlement MPC entity | source-qualified financial evidence + Economics attribution/reconciliation |
| PublicationPreparation aggregate | Readiness Q + Offering ListingIntent draft |
| SourceProductObservation business service | D4 acquisition feeding consumer-owned ports |
| generic listing transformation/rules engine | `FOLLOW_SOURCE | EXPLICIT_OVERRIDE`; later targeted reopen only on real repetition |
| dedicated AI business/API authority | D2 automation Principal + ordinary Product API/Governance |
| Offering-owned listing quantity | Availability-owned meaning; joint technical realization only |
| manual OpenAPI + manual SDK wire authority | one OpenAPI wire authority + derived/conformant clients/server |
| compatibility/versioning for current legacy API | hard cutover while no entitled consumer exists |

---

# 5. What remains genuinely open / deferred

This index does not convert later-stage work into architecture Unknowns, and it does not allow later stages to reopen already-settled meaning.

## D5-B2 current work

- exact Product 1.0 operation/resource inventory;
- exact request/response schemas;
- exact Permission→operation map;
- final paths/nouns inside B1 laws;
- pagination/filter/sort/cursor only for proven consumers;
- operation-local bulk only where proven;
- concrete OpenAPI/server/client generation and conformance tooling.

## D6 deferred

- screens/navigation/editor topology;
- projection/view composition required by UX;
- frontend feature/package structure.

## D7 deferred

- process/server/worker/scheduler topology;
- transaction/outbox/queue/cursor/lease/lock mechanisms;
- token refresh/cache/secret realization;
- RLS/runtime Organization-isolation realization;
- idempotency persistence/TTL/locking;
- retry/backoff/rate-control mechanisms;
- media blob/cache/CDN realization;
- production deployment topology;
- retained runtime ADR residues named by the ADR registry.

## D8 proof obligations / deferred concrete capabilities

- first controlled Mercado Livre create/Price/Availability effects + authoritative reread/convergence;
- shared-User-Product blast-radius real-write proof;
- selected Sankhya irreversible fiscal progression;
- controlled alternate destination/contact realization before claiming it;
- first consequential native-party write if the golden flow reaches it;
- unproven fiscal/provider branches only when a selected golden flow materially depends on them.

## D9

- final adversarial architecture/system review;
- implementation remains blocked until D9 acceptance.

## Bounded Unknown/Deferred items that remain honest

Detailed homes remain D4/D4-R1. Current classes include:

- unselected Mercado Livre modes/configurations;
- paused/zero-quantity representation-first creation on the selected Installation;
- broader payment/account movement universe absent a consumer;
- R3 bank-side reconciliation until an accepted bank source exists;
- controlled-product marketplace path until selection/interchangeability is proven;
- post-invoice fiscal return / selected reversal branches until required;
- fiscal components/branches not yet proven for a claimed flow;
- reusable source→listing transformations, Product-media library or external connector wire contract until real repeated consumers appear.

Unknown/Deferred is **never** permission to invent a plausible default.

---

# 6. Implementation reconciliation guard

Technical stages choose mechanisms inside accepted meaning. The following list highlights high-risk invariants; **it is not exhaustive**. `ARCHITECTURE.md` and accepted D-stage artifacts remain the complete constraint set. **Absence from this index is never permission.**

Technical work MUST NOT silently re-decide by convenience:

1. Product ownership — external/source-qualified, not MPC master.
2. Organization root/isolation semantics.
3. the D1 business authorities and accepted semantic edges.
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

If a technical design appears to require violating one of these, that is evidence for a targeted architecture reopen — not permission for a workaround.

---

# 7. Legacy ADR reconciliation / active-tree cleanup

## 7.1 Single ADR-status authority

The ADR registry is the **only active authority for ADR file status/disposition**.

This baseline may say where an old architectural idea routes semantically (§4) and may record the cleanup set approved by the operator, but it does not create a second `historical/superseded/reopened` vocabulary for individual ADRs.

After canonicalization, the registry should become a compact transition index containing:

- retained unresolved legacy ADRs only;
- one statement that all other pre-rebaseline ADRs have retired from the active tree and remain in Git history;
- future target ADR-036+ entries when genuinely warranted.

## 7.2 Corrected retained legacy set

Independent review falsified the proposed ADR-003 retention: its OAuth → fee-sync → UX sequencing has no unique current D9 consumer after D4-B1/D5 rehoming.

**Proposed retained set after operator ratification:**

| ADR | Why it remains active | Owning later stage / retirement condition |
|---|---|---|
| 008 | deployment/publisher/host evidence still needed for target production topology | D7; retire when deployment topology is adjudicated |
| 010 | acquisition cadence/polling/runtime evidence remains useful while D7 mechanics are open | D7; retire when acquisition runtime is adjudicated |
| 017 | legacy Fact/domain-judgment evidence still participates in D2's explicit rehoming gate | retire only when replacement target Fact ADR rehomes still-valid clauses |
| 018 | D7 execution-safety/runtime residue remains intentionally open | D7; retire after runtime mechanism adjudication |
| 026 | D7 scheduler/runtime residue remains intentionally open | D7; retire after scheduler/runtime adjudication |
| 030 | worker/scheduler/installation-topology evidence remains open | D7; retire after process/scheduler topology adjudication |
| 034 | current Fact implementation/evidence anchor explicitly retained by D2 | retire with 017 when target Fact ADR lands |
| 035 | governs the D0–D9 rebaseline authority transition itself | retire only when D0–D9 program closes |

**ADR-003 is proposed to retire to Git history** in the canonicalization commit; the ADR registry row must be adjudicated in that same atomic change.

## 7.3 ADR-035 stale-snapshot fence

ADR-035 remains a transition authority, but its embedded “still-binding” / “reopened” tables are a **2026-08-14 historical snapshot**, not current disposition authority.

Canonicalization must add a dated amendment/fence to ADR-035 stating:

> Current ADR dispositions are owned by the ADR registry and accepted D-stage artifacts. ADR-035's embedded decision tables do not override later adjudication; notably D4-B1 superseded Direct Oracle/godror target transport and ADR-006/007 are no longer current target constraints.

Its governing rebaseline role remains unchanged.

## 7.4 Retire-to-Git-history set

Subject to the operator's final ratification, active-tree retirement is safe for:

`001, 002, 003, 004, 005, 006, 007, 009, 011, 012, 013, 014, 015, 016, 019, 020, 021, 022, 023, 024, 025, 027, 028, 029, 031, 032, 033`.

For carried-constraint ADRs in that set, current meaning has already been rehomed into `ARCHITECTURE.md` and/or accepted D-stage authority. Git history remains the archive and ADR numbers are never reused.

## 7.5 Citation-harvest cleanup

`docs/architecture/decisions/_citations/` is legacy citation archaeology, not target authority.

During canonicalization:

- retire citation-harvest files whose citing ADRs all retire;
- retain only citation files still directly required by the retained ADR set (including Fact and remaining D7 residues);
- retire `RENUMBERING-REGISTRY.md` unless a retained ADR still materially depends on it; target ADR numbering remains ADR-036+ through the active registry, not through citation archaeology.

No `old/` or archive directory is created; Git history is sufficient.

---

# 8. Disposable review / stale candidate cleanup

After this reconciliation is ratified and canonically filed:

1. delete `docs/engineering/rebaseline/D5-B2-REVIEW-CANDIDATE.md` — it predates D4-R1 and must not seed the new B2;
2. update every router/read-path/fresh-session reference in the **same atomic change** so no deleted candidate is still named as live input;
3. replace `AI-DIALOG.md` with its protocol/header only, preserving the review channel for D5-B2→D9 while Git history retains prior rounds;
4. future review candidates are deleted when their accepted meaning is canonically filed.

Tracked review-candidate inventory at independent review time contained only the reconciliation candidate and stale D5-B2 candidate.

---

# 9. Fresh-session target read shape

After acceptance/canonicalization the intended read order is:

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

This requires coordinated updates to `AGENTS.md` and the router when the baseline becomes canonical.

The baseline is **always read** because it is the cheap routing map; keeping provider/runtime detail out of it is what makes that sustainable.

A fresh reviewer must not need retired ADRs, `AI-DIALOG.md`, stale candidates or legacy implementation to discover target architecture.

---

# 10. Proof strategy

Canonicalization is valid only if these checks pass:

1. every CURRENT route in this baseline resolves to accepted authority;
2. no accepted CURRENT decision is retired by cleanup;
3. no retired ADR contains unique still-needed target meaning absent from active authority;
4. every retained ADR names a concrete later-stage consumer and retirement condition;
5. no ADR number has independent status vocabulary in more than one active authority;
6. no retained active file can reasonably resurrect Direct Oracle/godror as current target transport;
7. every D0 cross-cutting safety class and `ARCHITECTURE.md` stable constraint maps either to a highlighted guard/law here or remains explicitly complete in its semantic home — absence from the baseline never weakens it;
8. bounded Unknown/Deferred items remain open in their detailed home;
9. the stale D5-B2 candidate can disappear without losing target authority;
10. `AI-DIALOG` history can disappear from the active file without breaking the future review channel;
11. after cleanup, a fresh reviewer can answer product, owners, identity, communication, external boundary, publication ownership, action safety, open stages and exact next action without archaeology;
12. Structural Inversion: deleting/inverting current implementation/OpenAPI does not change the reconciled architecture.

---

# 11. Reopen triggers

Reopen this reconciliation only when its routing map materially changes because:

- an accepted D-stage decision is amended/reopened;
- a new target ADR changes current decision routing;
- a later real proof invalidates an architectural assumption;
- a Product requirement changes D0/D1 ownership;
- a retained legacy ADR is adjudicated and leaves the active tree;
- D0–D9 closes and ADR-035/transition machinery retires.

Do not reopen for implementation naming, package layout, framework preference or rediscovery of a retired Git-history decision.

---

# 12. Converged outcome

**Proposed outcome:** `CURRENT D0→D4/D4-R1 + D5-B1 DECISION SET RECONCILED / COHERENT`.

Independent review conclusion incorporated here:

- no contradiction among D0→D4/D4-R1 + D5-B1;
- no missing business authority/boundary;
- no superseded decision presented as current;
- no current decision proposed for accidental retirement;
- no D6/D7 mechanism frozen by this index;
- no material architecture prerequisite blocks D5-B2;
- baseline is safe as an always-read routing authority only if it remains compact and scope-limited;
- corrected legacy KEEP set is `{008, 010, 017, 018, 026, 030, 034, 035}`;
- ADR-003 retires;
- stale D5-B2 candidate retires;
- AI-DIALOG resets to protocol header after this review;
- Git history remains the sole archive of retired decision generations and review dialogue.

If the operator explicitly ratifies this converged package, canonicalization must be **one coherent cleanup/authority change**:

1. replace this candidate with `DECISION-RECONCILIATION-BASELINE.md` as ACCEPTED/CANONICAL;
2. add it to `AGENTS.md` + router read order;
3. shrink the ADR registry and make it the sole ADR-status authority;
4. add the ADR-035 stale-snapshot fence;
5. retire all non-KEEP legacy ADR files and unneeded citation-harvest files;
6. delete stale `D5-B2-REVIEW-CANDIDATE.md` and update router references atomically;
7. reset `AI-DIALOG.md` to its protocol/header;
8. preserve D5-B2 as **NEXT / NOT YET OPENED** until this canonicalization verifies cleanly;
9. after verification, D5-B2 may be opened from current authority, not from any prior candidate;
10. implementation remains blocked until D9.

**This file remains non-authoritative until explicit operator ratification.**