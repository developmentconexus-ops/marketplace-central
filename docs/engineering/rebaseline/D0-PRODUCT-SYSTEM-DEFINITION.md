# D0 — Product / System Definition

> **Status:** WORKING — D0 OPEN, FINAL CORRECTIONS APPLIED; PENDING FINAL COLD REVIEW  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Authority:** operator-approved D0 decisions recorded here; later D-stages own technical realization  
> **Last updated:** 2026-08-15

## 1. Purpose and boundary

D0 defines what Marketplace Central (MPC) is, who it serves, which outcomes belong to Product 1.0, and where the product boundary ends.

D0 does **not** decide target contexts/modules, canonical database identities, schema, API shape, frontend topology, event/outbox mechanics, job/runtime topology or provider-specific transport. Those belong to D1–D7.

Existing code, schemas, OpenAPI, tests, runtime and historical ADRs are current-state evidence unless the active rebaseline explicitly marks a constraint binding.

---

## 2. D0.1 — Product mission

**Accepted by operator.**

Marketplace Central is the internal **Marketplace Operations Control Plane** of the organization.

It combines authoritative internal/business facts with marketplace observations so operators can understand operational reality, detect divergence, make grounded decisions, execute controlled actions in participating systems, and verify/reconcile the result.

```text
observe → understand → reconcile → decide/policy → execute → verify → audit/reconcile
```

MPC is neither a marketplace dashboard nor an ERP replacement. External systems remain authoritative for facts/processes inherently theirs. MPC owns the cross-system marketplace operating semantics: intent, policy application, workflow, correlation, controlled execution, operational state, divergence, audit and reconciliation.

---

## 3. D0.2 / D0.3 — Product 1.0 scope and capability boundary

**Accepted by operator.**

Product 1.0 is **Marketplace Operations + Commercial Intelligence**.

Its normal lifecycle is:

```text
internal product
  → channel readiness / linkage
  → listing + marketplace availability control
  → competitive market observation
  → pricing / simulation / expected profitability
  → controlled decision / policy
  → marketplace action
  → sale / marketplace order
  → Business Order Intent / business-system materialization
  → fulfillment readiness / physical execution
  → Invoicing Intent / fiscal materialization
  → provider-requirement closure / provider-accepted dispatch readiness
  → packing / dispatch
  → shipment / delivery lifecycle
  → essential cancellation / return / refund lifecycle when applicable
  → marketplace settlement
  → optional bank cash-receipt evidence when an accepted source exists
  → realized economics / reconciliation / learning
```

Product 1.0 capabilities are:

1. **Product & Channel Readiness** — determine which internal products can operate in a marketplace, their linkage/readiness and missing/conflicting conditions.
2. **Marketplace Listing Operations** — create, inspect and control listings/material channel state and verify convergence.
3. **Marketplace Availability Control** — derive and maintain sellable marketplace availability from eligible authoritative inventory facts/rules plus applicable MPC-owned allocation policy, with automatic normal-path synchronization and explicit uncertainty/failure.
4. **Competitive Intelligence** — observe comparable market offers/prices, expose competitive position/change and represent insufficient evidence honestly.
5. **Pricing & Profitability Intelligence** — combine market evidence and explainable internal economics under an explicit Cost Basis to simulate price, expected margin/profitability and trade-offs.
6. **Decision & Policy Control** — translate evidence/recommendations into permitted, approval-required or prohibited actions under governing rules/policies.
7. **Order & Invoicing Business-System Operations** — express `Business Order Intent`, cause/verify the corresponding native business order, and later express readiness-gated `Invoicing Intent` and verify the authoritative fiscal/documentary result without importing ERP-native operation types into MPC semantics.
8. **Marketplace Fulfillment / Dispatch** — progress marketplace sales through eligible physical fulfillment execution, provider-required prerequisites/data/artifacts/readiness, conference, invoicing trigger when applicable, packing and verified dispatch handoff without becoming a company-wide WMS.
9. **Shipment / Delivery Observation & Exceptions** — observe shipment after dispatch until a relevant terminal outcome and surface material delivery exceptions without becoming a TMS/carrier.
10. **Essential Post-Sale Operations** — coordinate cancellation, return and refund consequences for MPC-controlled marketplace sales without becoming general CRM/SAC or company-wide reverse logistics.
11. **Reconciliation & Exception Operations** — turn missing evidence, ambiguous outcomes and cross-system divergence into explicit work instead of plausible defaults or silent success.
12. **Economic Evidence & Realized Profitability** — preserve an explainable economic evidence chain from simulation through order economics and marketplace/payment settlement; when an accepted authoritative bank-side source exists, extend the chain to bank cash-receipt evidence without turning bank integration into a Product 1.0 launch gate.

Mercado Livre and Sankhya are the first concrete systems used to prove relevant capabilities. Provider/ERP-specific mechanics belong to D4.

### 3.1 Action authority model

**Accepted by operator.**

Product 1.0 supports:

- **human-controlled execution by default:** evidence/analysis → decision/approval → MPC executes → verifies/reconciles;
- **policy-driven automatic execution where explicitly authorized;**
- **human review for uncertainty, low confidence, policy conflict, failed/non-convergent actions or exceptions.**

Routine sufficiently-known policy-valid marketplace availability synchronization is automatic and does not require per-change human approval. Fully autonomous commercial repricing is not a launch gate.

### 3.2 Policy/rule provenance and operability

**Accepted by operator.**

A governing rule used by MPC may be **MPC-owned**, **externally governed** or **derived**. MPC preserves provenance and must not silently convert an external rule into editable MPC-owned policy.

Where a policy is legitimately MPC-owned, Product 1.0 must allow organization-operable configuration without code edits. The later policy model must support deterministic default/inherited policy plus explicit more-specific overrides where proven business scopes justify them; effective policy and provenance must be explainable. D0 does not create a generic speculative rules engine.

This applies to availability-allocation policy and internal operating-time targets. An organization may deliberately choose a stricter internal target than a provider deadline, but internal policy never rewrites or relaxes an externally governed obligation.

### 3.3 Organization posture for Product 1.0

**Accepted in final D0 closure review.**

The first Product 1.0 proof requires one real `Organization`. Multi-Organization/SaaS operation, tenant provisioning, cross-organization administration and SaaS billing are **not Product 1.0 launch gates**.

`Organization` remains an explicit non-collapsed business identity and must not be hard-coded away as a global singleton merely because the first proof uses one organization. The stable tenant-ready isolation invariant remains binding; D2/D7 decide the smallest identity/isolation mechanism that satisfies it.

### 3.4 Composite offers / kits

**Accepted in final D0 closure review.**

A business-system item that already behaves as an independently authoritative sellable SKU may operate as an ordinary Product 1.0 product even when its commercial description says kit/combo/conjunto.

A provider-native or MPC-recognized **composite offer** whose sellable availability or business-order materialization materially depends on component composition is **not a Product 1.0 launch requirement unless the selected first Mercado Livre operating flow requires it**. If such an offer is encountered, MPC must not silently flatten real composition semantics into an ordinary single-product assumption: the path must be explicitly supported with composition preserved or explicitly unsupported/external-required.

D1/D2/D4 decide the minimum composition model and provider/business-system mapping only if the accepted operating flow requires it.

### 3.5 Marketplace Installation operational health / reputation

**Accepted in final D0 closure review.**

Provider-authoritative Marketplace Installation health/reputation evidence, where exposed and materially relevant to marketplace operation, participates in MPC portfolio attention as **OBSERVE / DERIVE** evidence.

Product 1.0 may surface health/reputation degradation and correlate it with operational facts. Reputation optimization, complaint management, buyer service/conversation management and automated reputation-management workflows remain outside Product 1.0.

---

## 4. Product 1.0 non-goals safe to defer

Unless later material evidence reopens them, Product 1.0 does not require:

- paid ads/media management;
- buyer Q&A/chat or buyer-conversation management;
- general CRM/SAC, complaint handling or reputation-management automation;
- campaign/discount-campaign authoring as a separate control surface;
- company-wide demand forecasting/purchasing;
- company-wide WMS/TMS/reverse-logistics replacement;
- company-wide accounting, treasury, accounts-payable/receivable or bank reconciliation;
- bank/Open Finance/statement integration as a Product 1.0 launch gate;
- Multi-Organization/SaaS operation, tenant provisioning or SaaS billing as a launch gate;
- provider-native/composite kit semantics unless the selected first Mercado Livre operating flow requires them;
- marketplaces beyond what is needed to prove the first Mercado Livre operating loop;
- unrestricted AI/autonomous commercial action without explicit policy;
- a generic integration framework for speculative providers;
- a universal ERP model;
- runtime dependency on marketplace hubs such as ANYMARKET or Magis5 merely to reach marketplaces MPC is intended to control directly.

Observed promotion/discount effects that materially alter price, order economics or realized economics remain attributable evidence even though campaign authoring itself is deferred.

Shipment/delivery observation, essential cancellation/return/refund consequences, provider-requirement closure for claimed normal fulfillment paths, materially time-bound obligations, marketplace/payment settlement and the evidence needed to explain realized marketplace economics are **not** deferred merely because their source systems are external.

---

## 5. External-system evidence / target direction

### 5.1 Sankhya evidence

The operator reports an already-proven Sankhya application integration in another application using application credentials (`client id`, `client secret`, `X-Token`; no secret values are recorded here). Operational writes such as order creation and invoicing can be performed through Sankhya's application API; system-owned APIs are preferred for writes rather than direct Oracle writes. DB Explorer/database inspection is also available.

This is D0/D4 evidence, not target transport or canonical domain design. D4 independently ratifies exact read/write capability contracts. `CODEMP`, `CODLOC`, TOPs, native cost variants and other Sankhya structures are integration evidence, **not automatically canonical MPC concepts**.

### 5.2 Marketplace / provider / competitor evidence

Current research covers Mercado Livre, Amazon, Magalu, Casas Bahia, Leroy/Mirakl, Shopee, MadeiraMadeira, ANYMARKET, Magis5 and Bling. Detailed observations remain in `EVIDENCE-REGISTER.md`.

**Accepted target direction:** MPC is itself the marketplace operations/control-plane product. For marketplaces it supports, the target direction is an MPC-owned direct provider integration boundary rather than routing marketplace operation through another marketplace hub by default.

ANYMARKET, Magis5 and similar marketplace hubs are benchmark/competitive evidence, not Product 1.0 runtime dependencies. Bling may later be integrated as an ERP/business system when used in that role. Shared provider technology such as Mirakl may justify technical reuse without changing the business marketplace identity.

Current evidence also supports these D0 closure conclusions:

- provider/fulfillment mode can materially change capability and responsibility;
- bulk/provider writes can be asynchronous, partial or later-discovered rather than one atomic success;
- provider-side approval/safety controls can keep a proposed change pending or block it;
- notifications/webhooks can be duplicated, lost or out of order and are not automatically complete/current truth;
- Mercado Livre exposes provider-native composite/kit behavior and seller-reputation evidence, confirming those are real provider concerns even though D0 deliberately keeps their Product 1.0 scope minimal.

---

## 6. Stable constraints carried into D0

- Mercado Livre first.
- Sankhya/Oracle external to MPC.
- Go backend is canonical business execution.
- React frontend is not a second business authority.
- PostgreSQL stores MPC-owned canonical state.
- tenant-ready data isolation remains a real invariant even though multi-Organization SaaS operation is not a Product 1.0 launch gate.
- unknown/absent facts do not become plausible known defaults.
- external writes require explicit authority/policy, duplicate protection, auditability and reconciliation.
- ambiguous external-write outcomes are not blindly retried.
- provider PII is minimized.
- provider-specific protocol details remain behind provider boundaries.

---

## 7. D0.4 — Actors / operational users

**Accepted by operator.**

### Marketplace Operations Operator

Owns routine marketplace control inside policy: readiness/linkage preparation, listing creation/editing where ready/allowed, competitive/pricing analysis, bounded price/listing actions and operational exception work. This actor does not redefine governing policy merely to make an action permissible.

### Fulfillment / Dispatch Operator

For an internally operated selected Fulfillment Node, owns the physical work queue, separation/conference, invoicing trigger when readiness is proven, provider-required fulfillment/handoff steps that belong to the accepted workflow, packing, dispatch handoff and exception reporting.

Accepted internal normal path:

```text
marketplace sale / business-order readiness
  → eligible / selected Fulfillment Node
  → separation
  → physical conference
  → if valid + invoicing prerequisites proven: trigger invoicing through MPC
  → authoritative invoicing result verified/reconciled
  → satisfy/orchestrate provider-required dispatch prerequisites/artifacts/readiness
  → packing / dispatch handoff
```

The internal normal path does not intentionally invoice before physical confirmation that the correct items are available/separated. Physical inconsistency blocks normal invoicing and becomes explicit exception work.

Externally operated Fulfillment Nodes may prove readiness differently in later stages; Product 1.0 need not implement every fulfillment mode in its first proof.

### Commercial / Marketplace Manager

Ordinary authority over legitimate **MPC-owned** commercial policies: margin floors, price boundaries, approval thresholds, availability-allocation policy, bounded automation, internal operating targets and commercial exceptions. This actor is not automatically authority over externally governed ERP/provider rules or integration/security administration.

### Owner / Administrator / Policy Approver

Governs organization-level access/integrations, exceptional authority above delegated commercial management, governance boundaries for automation/policy and emergency containment. It is not a routine commercial approval bottleneck.

No actor may configure away mandatory audit/reconciliation/safety invariants, externally governed obligations or silently convert external rules into local mutable copies.

---

## 8. D0.5 — System boundary / authority classes

**Accepted by operator.**

- **OWN** — MPC is business authority for the concern/state.
- **ORCHESTRATE** — another system owns the native fact/process; MPC owns cross-system intent/control/workflow around it.
- **OBSERVE / DERIVE** — MPC consumes external authoritative facts and may derive decision-support/operational conclusions without becoming source of truth for the base fact.

> **MPC owns the marketplace operating model. External systems own the facts and processes that inherently belong to them. MPC orchestrates their participation in the marketplace operating loop.**

### 8.0 Authority-map notation

The three authority classes above remain the only product-level authority classes.

Combined words in the table describe **different aspects of one concern**, not new authority classes:

- `consume`, `record` and `observe` describe acquisition/preservation of externally authoritative evidence;
- `derive` describes an MPC conclusion produced from evidence;
- `orchestrate` / `reconcile` describe MPC-owned cross-system intent/control around externally authoritative facts/results;
- `own` describes MPC business authority for the named concern.

### 8.1 Product-level authority map

| Concern | MPC authority |
|---|---|
| Internal master product / ERP-native fiscal and cost facts | **OBSERVE / CONSUME** |
| Native stock/reservation/availability facts in ERP/WMS/3PL/provider | **OBSERVE / CONSUME** |
| Product ↔ marketplace linkage and readiness conclusion | **OWN / DERIVE** |
| Marketplace Installation membership/configuration | **OWN**; provider account identity/state remains provider-owned |
| Marketplace Installation health/reputation native evidence | provider authoritative; MPC **OBSERVES / DERIVES** portfolio attention |
| Selling Entity semantic and transaction attribution | **OWN**; corresponding external legal/company records may remain externally authoritative |
| Inventory Source identity and Inventory Scope eligibility | **OWN** for MPC semantics/configuration; native stock source/facts remain external |
| Availability Allocation Policy | **OWN** when MPC-governed |
| Sellable Availability | **OWN / DERIVE** from eligible authoritative facts/rules/policy |
| Provider-native composite offer/component facts | provider/business-system authority where native; MPC **OBSERVES / ORCHESTRATES** only when the accepted flow supports composition |
| Provider actual listing/price/availability state | provider authoritative; MPC **OBSERVES / ORCHESTRATES** |
| Listing/price/availability mutation intent | **OWN / ORCHESTRATE** |
| Fulfillment Node / Fulfillment Scope / routing intent | **OWN / ORCHESTRATE**; native execution facts stay with their authority |
| Effective provider capability/requirement interpretation | **OWN / DERIVE** from provider-native evidence; provider remains authority for native state/artifact semantics |
| Provider-required fulfillment prerequisites/readiness/artifacts | provider authoritative where native; MPC **OBSERVES / ORCHESTRATES / RECONCILES** closure inside claimed normal paths |
| External operational deadline/window | external authority; MPC **OBSERVES** source/time evidence |
| Internal Operational Target | **OWN** as MPC organization policy |
| Operational urgency/safety-margin/breach interpretation | **OWN / DERIVE** |
| Actionable-work ownership/assignment/escalation/resolution state | **OWN** as MPC operational-control semantics; actor/policy authorization remains separate |
| Observation/acquisition time and provenance for external facts | source/provider owns native fact time where supplied; MPC **RECORDS** observation provenance and **DERIVES** freshness-for-use |
| Observation coverage/completeness conclusion | **OWN / DERIVE** from source enumeration/observation evidence; underlying facts remain externally authoritative |
| Material decision/action lineage and historical explanation | **OWN** as MPC audit/control semantics, built from decision-time evidence/provenance, governing policy/rules, authority/approval context, uncertainty and correlated action/result evidence |
| Multi-target intended action scope and aggregate/granular execution outcome | **OWN** as MPC control/reconciliation semantics; each provider-native result remains authoritative in its source system |
| Execution-time decision/approval validity conclusion | **OWN / DERIVE** from the authorized intent plus current material authority, policy, evidence-sufficiency, readiness and mandatory safety conditions; underlying external facts remain with their authorities |
| Cost Observation | **OBSERVE / DERIVE** from attributable authoritative economics |
| Cost Basis | **OWN** when MPC-governed; **OBSERVE / CONSUME** when externally governed |
| Competitive/pricing/expected-profitability interpretation | **OWN / DERIVE** |
| Marketplace order source facts | marketplace/provider authoritative |
| Business Order Intent + correlation/materialization workflow | **OWN**; ERP-native order remains ERP-authoritative |
| Invoicing readiness + Invoicing Intent + correlation workflow | **OWN / DERIVE** for readiness, **OWN** for intent/workflow; native fiscal result remains ERP/fiscal-system authoritative |
| Shipment/delivery native facts | provider/carrier authoritative; MPC **OBSERVES / ORCHESTRATES** |
| Cancellation/return/refund native facts | provider authoritative; MPC **ORCHESTRATES** cross-system consequences |
| Native ERP reversal/credit/fiscal/accounting facts | ERP/fiscal system authoritative; MPC **ORCHESTRATES** |
| Simulation / Expected Economics | **OWN / DERIVE** from explicit inputs, Cost Basis, market/provider facts and policies |
| Order Economics interpretation | **OWN / DERIVE** from authoritative sale/order/payment/shipping/economic facts |
| Marketplace/payment-account settlement movements | marketplace/payment provider authoritative; MPC **OBSERVES / DERIVES / RECONCILES** |
| Bank cash-receipt evidence | bank/accepted bank-data source authoritative when available; MPC **OBSERVES / RECONCILES**; not a Product 1.0 launch gate |
| Economic Evidence Chain / variance classification / calibration cases | **OWN / DERIVE** |
| Realized profitability interpretation | **OWN / DERIVE** from attributable economic evidence |
| MPC audit/reconciliation/exception state | **OWN** |

### 8.2 Boundary invariants

1. Observing an external fact does not transfer ownership of it to MPC.
2. Orchestrating an external process does not make MPC source of truth for its native record.
3. Derived conclusions preserve enough provenance/freshness/coverage/context to explain their evidence.
4. Externally governed rules are not silently copied into mutable MPC policy.
5. Ambiguous/divergent outcomes remain explicit until reconciled.
6. **Unknown availability is not zero** or another plausible quantity.
7. Organization, Marketplace Installation and Selling Entity do not collapse merely because a first deployment uses one of each.
8. The first Product 1.0 proof using one Organization does not make Organization a hard-coded global singleton; Multi-Organization/SaaS operation is nevertheless not a Product 1.0 launch gate.
9. MPC canonical semantics come from marketplace-operating needs, not ERP/provider ontology.
10. ERP/provider integration is semantic translation, not field renaming; unsupported/ambiguous mapping is explicit.
11. Selling Entity, Inventory Source, Fulfillment Node and Cost Basis remain distinct unless explicit business evidence relates them.
12. Inventory Scope is explicit eligibility, not “all stock we can find”.
13. Availability Allocation Policy changes allocation, not authoritative stock truth.
14. MPC-owned policy defaults/overrides resolve deterministically with visible provenance.
15. Provider-native/composite composition that materially affects availability or order materialization is not silently flattened into single-product semantics.
16. Cost Observation is not a bare number; Cost Basis is not an ERP cost-type alias.
17. Missing/ambiguous cost does not silently fall back; current cost is not historical transaction cost by default.
18. Cost does not absorb all sale economics; fees, freight, taxes, discounts, subsidies and reversals remain separately attributable when material.
19. Business Order Intent is not an ERP order type/TOP; marketplace order, MPC intent and ERP-native order remain distinct authorities.
20. Missing/ambiguous order mapping is explicit; ambiguous order-write outcome is not blind retry.
21. No separate `Order Execution Scope` is introduced without independent business semantics.
22. Invoicing Intent is not an ERP/native fiscal operation; order existence does not equal invoicing readiness.
23. Unknown invoicing prerequisite does not become ready; ambiguous invoicing outcome is not blind retry.
24. No separate `Fiscal Scope` / `Invoicing Scope` is introduced without independent business evidence.
25. Post-sale fiscal consequences remain controlled/orchestrated without importing native fiscal taxonomy into D0.
26. Economic evidence stages do not overwrite each other.
27. Simulation variance is not automatically a simulator defect; calibration requires materially comparable context.
28. Marketplace settlement is not bank cash receipt.
29. L3 bank cash-receipt evidence is an optional extension when an accepted authoritative source exists; bank integration is not a Product 1.0 launch gate.
30. Unattributable financial movement is not invented attribution.
31. Payout/cash matching is not assumed 1:1 with an order.
32. Simulator calibration is evidence-driven.
33. Provider capability/authority is context-sensitive; marketplace brand alone does not determine responsibility.
34. A claimed MPC-controlled normal fulfillment path must close provider requirements.
35. Native provider states/artifacts remain provider-native rather than universal MPC domain types by analogy.
36. Unsupported/external-required provider work is explicit.
37. Marketplace hubs are not Product 1.0 runtime dependencies by default.
38. Shared provider technology does not change business-provider identity.
39. Marketplace Installation health/reputation may affect portfolio attention without turning MPC into a reputation-management/SAC product.
40. External obligation and internal target are different authorities; internal policy may tighten but never relax the external obligation.
41. Relative time policy requires an explicit trustworthy time anchor.
42. Internal-target breach and external-obligation breach are distinct states.
43. Material time drives operational attention.
44. Actionable work is not operationally ownerless.
45. Assignment is not authorization.
46. Escalation is not notification; it means greater/different responsibility, attention or authority is required.
47. Material work does not close by arbitrary dismissal; accepted resolution evidence is required.
48. Freshness is use-sensitive; D0 does not define a universal age threshold.
49. Stale evidence is not false evidence, but cannot masquerade as current truth.
50. Unknown freshness is not current.
51. Acquisition failure affects operational confidence and remains distinguishable from the last successful observation.
52. Unsafe action may be blocked/degraded by insufficient freshness when current evidence is materially required.
53. **Fresh is not complete.** A newly observed subset does not prove the relevant source population was fully observed.
54. **Not observed is not does-not-exist.** Absence from a partial/unknown-coverage observation cannot become a terminal or negative fact merely because a record was not returned.
55. **Coverage is scoped.** Any claim of sufficient completeness must be relative to an explicit population/range/scope/obligation rather than a global `providerComplete=true` assertion.
56. **Individual observed facts do not prove population completeness.** A fact can be valid while portfolio-level coverage remains partial/unknown.
57. **Portfolio health and reconciliation closure require sufficient coverage for the claimed universe.** Partial/unknown coverage cannot truthfully produce `no issues`, `nothing missing` or `reconciled` merely because observed subsets agree.
58. **Webhook/callback receipt is not completeness proof.** Notification delivery may contribute evidence but cannot alone prove every relevant change/record was observed unless the provider contract independently guarantees that property and D4 verifies it.
59. **Current state is not historical decision basis.** Material historical explanation must use decision-time evidence/provenance rather than mutable current values.
60. **Current policy is not necessarily the policy that governed a historical action.** Historical explanation preserves the governing policy/rule meaning/provenance effective for that decision.
61. **Action audit is not decision explanation.** Knowing who/what invoked an action is insufficient without enough evidence, authority/approval and uncertainty context to explain why it was permitted/recommended/executed.
62. **Re-running a current calculation is not reproducing a historical calculation.** Simulator, cost, fee, rule or policy changes must not rewrite the meaning of a prior decision.
63. **Material external action remains correlatable to intent + authority + decision-time evidence + result.** The explanation must survive later mutation of current state.
64. **Multi-target action is not one opaque boolean result.** A material action affecting multiple targets preserves an intended target scope and sufficiently granular outcomes to distinguish confirmed, rejected, ambiguous and not-executed members where those distinctions are material.
65. **Intended target scope is not reconstructed from mutable current filters.** Historical blast radius must remain explainable from the target universe intended/authorized for that action, even if current eligibility/filter membership changes later.
66. **Batch-level authorization does not erase target-level authority/readiness differences.** A collective decision cannot silently authorize a member that independently requires different authority, approval, evidence or readiness.
67. **Partial success is neither total success nor total failure.** Confirmed effects remain confirmed while rejected/ambiguous/not-executed members remain independently visible/reconcilable.
68. **An ambiguous member does not make whole-batch retry safe.** MPC must not blindly retry already-confirmed members merely because other members failed or became uncertain.
69. **MPC does not invent cross-provider atomicity.** Unless an accepted provider contract proves atomic behavior, external multi-target effects may be partial and Product 1.0 must remain correct under partial outcomes.
70. **Automation at scale does not hide blast radius.** Policy-driven mass effects remain observable/reconcilable even when no per-target human approval is required.
71. **Approval is not permanent permission.** A material decision/approval authorizes execution under its governing decision context; it does not waive later-invalidated authority, policy, readiness, evidence-sufficiency or mandatory safety conditions.
72. **Any state change is not automatic invalidation.** Execution-time revalidation is materiality-based; irrelevant drift must not create unnecessary reapproval or block routine work.
73. **Material governing drift may invalidate prior authorization.** If the facts/rules/authority/readiness on which permission materially depended no longer hold, MPC must block, recompute or require new authorization rather than execute from stale permission.
74. **Approval cannot waive mandatory safety invariants.** An earlier human or automated decision does not make an ambiguous write safe to retry, stale/insufficient evidence current, or a newly prohibited action permissible.
75. **Decision-time validity is not execution-time validity.** Consequential execution must still be valid enough when it occurs, whether the original decision came from a human, policy or automation.
76. **Multi-target approval does not force one execution-validity result for every member.** Material revalidation may preserve execution for still-valid targets while blocking/recomputing/reapproving others.
77. D0 does not require event sourcing, universal snapshots or copying every source fact. D2/D3/D7 later choose the smallest retention/reference/snapshot mechanism that satisfies accepted historical explainability.
78. D0 does not choose locks, optimistic-concurrency/version tokens, approval expiry durations, compare-and-swap, distributed locking or reapproval UI. D2/D5/D6/D7 choose the smallest execution-time validity mechanism needed by each action.
79. Provider/ERP/payment/bank transport, timer, scheduler, polling, webhook, TTL/cache, freshness threshold, pagination, checkpoint, backfill, coverage-proof, decision-record, batch API, queue, concurrency, compensation and retry-orchestration mechanics remain later-stage responsibilities.

---

## 9. D0.6 — Product 1.0 completion / user-observable outcomes

**Accepted by operator; refined by later D0.7 decisions and final closure review.**

Product 1.0 is complete only when MPC is demonstrably usable as the normal marketplace operations control plane.

1. **Attention is portfolio-driven.** Healthy/changed/divergent/blocked/approval-required/actionable state is visible without external manual product-by-product checking; materially relevant Marketplace Installation health/reputation can participate in that attention.
2. **An eligible internal product reaches verified marketplace state through MPC.** Readiness/linkage → commercial analysis → creation/publication → real channel observation.
3. **Inventory eligibility/allocation is explicit and availability converges.** Eligible Inventory Sources + authoritative facts/rules + applicable MPC policy derive Sellable Availability; uncertainty/failure is explicit.
4. **Composite semantics are honest when encountered.** Product 1.0 does not require provider-native composite offers for launch, but composition that materially affects availability/order materialization is not silently flattened.
5. **Competitive/pricing intelligence replaces major manual analysis with explainable economics.** Cost Basis and evidence are visible; insufficient evidence is honest.
6. **Decision closes into controlled action.** Human or bounded policy decisions become auditable external actions with verification/reconciliation.
7. **Material business context is explicit.** Selling Entity, installation and other proven semantics are not hidden defaults.
8. **Fulfillment responsibility is explicit.** Eligible/selected Fulfillment Node and physical responsibility are known when material.
9. **A marketplace sale becomes the correct business-system order through MPC.** Business Order Intent is mapped, executed, correlated and verified without leaking native order types into canonical semantics.
10. **Invoicing is readiness-gated and verifiable.** Invoicing Intent is emitted only after sufficient evidence; native result is observed/reconciled.
11. **Provider requirements close inside claimed normal fulfillment paths.** MPC determines effective provider requirements/capabilities, surfaces/orchestrates required data/artifacts/acknowledgements and proves provider-specific readiness without hidden routine provider-UI work.
12. **The normal sale lifecycle runs through MPC** for responsibilities it claims to control.
13. **Material time-bound obligations drive attention.** External deadlines and organization-owned internal targets remain distinguishable and relative targets have known anchors.
14. **Actionable work has an owner and can progress/escalate.** Assignment is supported when useful and never bypasses authorization.
15. **Work closes through evidence.** Material unresolved divergence cannot disappear through arbitrary dismissal.
16. **Material evidence is freshness-aware.** Stale/unknown-freshness evidence and degraded acquisition affect confidence/attention and can block unsafe normal action.
17. **Portfolio/absence/reconciliation claims are coverage-aware.** MPC distinguishes sufficiently-covered, partial and unknown coverage for the relevant population; `no issues`, terminal absence and reconciliation closure require enough coverage evidence for the claim being made.
18. **Material historical decisions/actions remain explainable.** MPC can explain why a consequential recommendation/approval/action was permitted or executed using decision-time evidence, governing policy/rule provenance, authority/approval and known uncertainty rather than reconstructing history from mutable current state.
19. **Material multi-target actions expose blast radius and partial outcomes.** MPC preserves the intended target scope and can distinguish confirmed, rejected, ambiguous and not-executed effects sufficiently to reconcile partial execution without pretending the whole operation was atomically successful/failed.
20. **Material authorization remains valid at execution time.** A prior human/policy/automation decision does not execute blindly after material governing context changes; irrelevant drift does not force needless reapproval, while material invalidation blocks/recomputes/requires new authorization as appropriate.
21. **Shipment remains visible through terminal outcome.** Material delivery exceptions become explicit work.
22. **Essential post-sale changes remain controlled.** Cancellation/return/refund and required fiscal/economic consequences remain orchestrated/reconciled.
23. **Failures become explicit work.** Missing evidence, ambiguity, unsupported capability, deadline uncertainty/breach, stale evidence, partial/unknown coverage, partial action outcome and invalidated authorization identify what is known/unknown, who is responsible and what action/authority is required.
24. **Core marketplace economic lineage closes through marketplace/payment settlement.** Expected Economics → Order Economics → Marketplace/Payment Settlement remain distinct and reconciled as far as authoritative provider/payment evidence permits.
25. **Bank cash evidence is an optional extension, not launch gate.** When an accepted authoritative bank-side source exists, MPC may extend the chain through Cash Receipt evidence without redefining L0–L2.
26. **Simulator drift becomes actionable learning** only after materially comparable historical evidence isolates model/provider-rule drift.
27. **Realized profitability is explainable** from attributable economic evidence rather than convenient current values.
28. **Organizational governance is operable without code edits** while external authority and mandatory safety/audit invariants remain protected.

Completion statement:

> **A company can take internal products, determine marketplace readiness, publish and operate offers, maintain availability from explicit inventory scope/policy, observe materially relevant installation health, monitor competition and profitability, make controlled decisions from sufficiently current and sufficiently covered evidence, retain enough decision-time lineage to explain material historical actions, revalidate material authorization before consequential execution, execute at scale without hiding blast radius or partial outcomes, materialize sales in business/fiscal systems, fulfill them through eligible nodes, close provider requirements, prioritize time-bound obligations, route actionable work to explicit owners, close work through evidence, dispatch and observe delivery/post-sale consequences, preserve economic evidence through marketplace/payment settlement and optionally bank cash evidence when an accepted source exists, reconcile divergence without mistaking partial observation or partial execution for closure, learn when its own simulator is wrong, and explain realized economic outcome — using Marketplace Central as the normal marketplace operations control plane.**

### 9.1 Normal-path rule

The normal operational path must be executable through MPC for responsibilities Product 1.0 claims to control. Direct use of an external system remains legitimate only for intentionally external responsibilities, investigation/support and explicit exceptional recovery; hidden routine system hopping means the claimed workflow is incomplete.

If a provider-required operation cannot be performed through an accepted Product 1.0 integration, that limitation must be explicit; the affected path cannot be presented/proven as fully MPC-controlled.

D8 later proves detailed golden flows; D0 defines what those flows must prove.

---

## 10. D0.7 — Product completeness / contradiction review

D0.7 adversarially tested Product 1.0 for missing lifecycle responsibility, contradiction and accidental scope gap.

### D0.7a — Essential post-sale lifecycle

**Accepted.** Cancellations, returns and refunds tied to MPC-controlled sales remain inside the controlled lifecycle, including required cross-system/fiscal/economic consequences. Buyer Q&A/chat and buyer-conversation management, general CRM/SAC, complaint/reputation-management automation and company-wide reverse logistics remain outside.

### D0.7b — Shipment / delivery lifecycle

**Accepted.** MPC observes relevant shipment/delivery state after dispatch through terminal outcome and turns material delay/failure/return-to-sender states into explicit work without becoming TMS/carrier.

### D0.7c — Stock / marketplace availability control

**Accepted.** Sellable Availability is derived from authoritative stock/reservation facts, Inventory Scope, governing rules and applicable MPC-owned allocation policy; routine known policy-valid synchronization is automatic and verified. Unknown availability does not become zero/plausible quantity.

### D0.7d — Marketplace installation multiplicity

**Accepted and clarified at closure.** `Organization 1 → N Marketplace Installations`. The first Product 1.0 proof may use one Organization and one Mercado Livre account, but organization identity never collapses into seller-account identity and Organization is not hard-coded away as a global singleton. Multi-Organization/SaaS operation is not a Product 1.0 launch gate.

### D0.7e — ERP-agnostic business semantics before ERP mapping

**Accepted.** MPC defines business semantics from marketplace-operating needs first; adapters translate them to/from Sankhya or another ERP. Native ERP constructs do not become canonical MPC concepts by existence alone.

#### D0.7e.1 — Selling Entity

**Accepted.** `Organization 1 → N Selling Entities`. Selling Entity answers which business/legal/fiscal entity is acting when material and remains distinct from Marketplace Installation, Inventory Source, Fulfillment Node and Cost Basis.

#### D0.7e.2 — Inventory Source / Inventory Scope

**Accepted.** Inventory Source represents a business-recognized inventory source/pool; Inventory Scope governs which sources may contribute to an offer. Native mapping need not be 1:1. Stock outside scope does not leak into availability.

##### D0.7e.2a — Availability Allocation Policy

**Accepted product requirement.** MPC-owned allocation policy may intentionally expose less than full eligible availability (for example 70%). Policy changes allocation, not stock truth. Exact catalog/arithmetic/scope hierarchy belong later.

#### D0.7e.3 — Fulfillment Node / Fulfillment Scope

**Accepted.** Fulfillment Node is a recognized physical-fulfillment execution point/capability; Fulfillment Scope governs eligibility. Inventory promise and fulfillment execution remain distinct semantics.

#### D0.7e.4 — Cost Observation / Cost Basis

**Accepted.** Cost Observation preserves meaning/value/time/business context/provenance. Cost Basis expresses which cost semantic is economically appropriate and is not an ERP cost-type alias. Unsupported/ambiguous cost remains explicit; historical economics does not silently use current cost.

#### D0.7e.5 — Business Order Intent

**Accepted.** Business Order Intent is MPC-owned intent/correlation for materializing a marketplace sale in a participating business system. ERP TOP/document/operation codes stay in integration mapping. Missing mapping is explicit; ambiguous writes are not blindly retried. No separate `Order Execution Scope` is justified yet.

#### D0.7e.6 — Invoicing Intent

**Accepted.** Invoicing Intent is readiness-gated MPC fiscal/documentary materialization intent. Order existence is not readiness; native fiscal result remains external authority; ambiguous invoicing is not blind retry. No separate `Fiscal/Invoicing Scope` is justified yet.

### D0.7f — Economic Evidence Chain

**Accepted and simplified at closure.**

Product 1.0 preserves distinct evidence layers:

```text
L0 Simulation / Expected Economics
  → R1 simulation reconciliation
L1 Order Economics
  → R2 settlement reconciliation
L2 Marketplace/payment-account Settlement

optional extension when an accepted source exists:
  → R3 cash reconciliation
L3 Bank Cash Receipt evidence
```

L0–L2 are the Product 1.0 core marketplace economic lineage. L3/R3 preserve the semantic distinction between provider/payment settlement and bank receipt, but **bank-side integration is not a Product 1.0 launch gate**.

Stages do not overwrite each other. Simulation/order variance is classified before becoming a calibration defect. Settlement is distinct from bank receipt. Payouts may aggregate many orders/movements and unattributable movement is never fabricated. This closes marketplace economic lineage without turning MPC into company-wide finance/accounting/treasury.

### D0.7g — Context-sensitive provider capabilities / Provider Requirement Closure

**Accepted.** Marketplace brand alone does not determine capability/authority; effective responsibility can vary by installation/offer/order/native operating mode. Every flow claimed as MPC-controlled must surface/satisfy/orchestrate required provider prerequisites/data/artifacts/readiness or explicitly mark the path unsupported/external-required. Native provider artifacts/states remain provider-native.

MPC targets direct provider boundaries. ANYMARKET, Magis5 and similar hubs are benchmark/competitive evidence, not Product 1.0 runtime dependencies. Shared provider technology such as Mirakl may justify implementation reuse without changing business marketplace identity.

Provider-native/composite offers follow the same rule: composition is not a launch requirement unless the selected first flow needs it, but a real composition dependency cannot be hidden behind single-product semantics.

### D0.7h — Time-bound operational obligations / internal targets

**Accepted.** External Operational Obligation remains provider/contract authority. Internal Operational Target is MPC-owned organization policy and may intentionally be stricter. Relative targets require explicit trustworthy anchors. Internal breach and external breach are distinct. Time participates in portfolio attention; exact clocks/timers/schedulers belong later.

### D0.7i — Operational work ownership / assignment / escalation / closure

**Accepted.** Material actionable work has a durable owning role/responsibility. Individual assignment is optional/distinct; assignment never grants authority. Automation failure may create human-owned work. Escalation means greater/different responsibility, attention or authority—not merely notification. Material work closes only through accepted resolution evidence, not arbitrary dismissal. MPC is not a generic ticket/task manager.

### D0.7j — Operational evidence freshness / staleness

**Accepted.** Material evidence used for decisions/readiness/automation/external action must carry enough observation/acquisition provenance to judge freshness for that use. Freshness is use-sensitive. Stale last-known values remain last-known rather than false/zero, but cannot masquerade as current. Unknown freshness and failed refresh attempts are explicit. Insufficient freshness may degrade confidence, trigger work or block unsafe normal action. Exact TTL/polling/cache mechanics belong later.

### D0.7k — Observation coverage / completeness uncertainty

**Accepted.** Product 1.0 distinguishes observed facts from evidence that the relevant source population was observed completely enough for a conclusion. Coverage is scoped to an explicit universe. Fresh is not complete; not-observed is not does-not-exist; individual facts do not prove population completeness; portfolio health, terminal absence and reconciliation closure require sufficient coverage; webhook/callback receipt is not automatic completeness proof. Exact pagination/backfill/checkpoint mechanics belong later.

### D0.7l — Decision/action evidence lineage and historical reproducibility

**Accepted.** For material recommendations, approvals, decisions and externally consequential actions, Product 1.0 preserves enough decision-time lineage to explain why the action was recommended, permitted and/or executed without reconstructing the past from mutable current state. Material intent, governing evidence/provenance, effective policy/rules, actor/automation authority/approval, freshness/coverage/uncertainty and correlated result remain explainable. D0 does not mandate event sourcing, universal snapshots or a canonical `DecisionRecord`; D1/D2/D3/D7 choose the smallest sufficient mechanism.

### D0.7m — Multi-target action scope / blast radius / partial-result semantics

**Accepted.** For material operator- or policy-driven actions affecting multiple targets, Product 1.0 preserves the intended target scope/blast radius and sufficiently granular outcomes. Target-level authority/readiness differences remain meaningful; confirmed/rejected/ambiguous/not-executed effects do not collapse into one opaque boolean; provider non-atomicity is not hidden; whole-batch blind retry cannot erase already-confirmed effects. D5/D6/D7 choose the smallest bulk-operation mechanism.

### D0.7n — Decision / approval validity at execution time

**Accepted by operator after explicit YAGNI/overengineering review.**

A material human, policy or automated decision/approval authorizes execution only while the **material governing context** on which that authorization depends remains valid enough for the intended action.

Product-level requirements:

- an approval is authorization under a decision context, not permanent permission independent of later facts;
- before consequential execution, MPC must be able to detect material invalidation of governing authority, policy, readiness, evidence sufficiency/freshness/coverage, provider capability or mandatory safety conditions when those conditions materially matter to the action;
- irrelevant state drift does not automatically invalidate a decision or require human reapproval; revalidation is proportional to the conditions that actually governed the decision;
- when material governing context changed, MPC may block execution, recompute the decision, obtain fresher/sufficient evidence or require a new authorization rather than executing from stale permission;
- an earlier approval cannot waive mandatory safety invariants such as ambiguous-write protection, externally governed prohibition or required evidence/readiness;
- the same principle applies to human approvals and policy/automation decisions;
- for multi-target actions, materially changed context may invalidate only some targets; the decision model must not force one validity result over the entire batch merely for convenience;
- D0 does not choose locking, optimistic concurrency/version tokens, approval-expiry duration, compare-and-swap, distributed locks or reapproval UI.

This is a product safety/authority property. D2/D5/D6/D7 later choose the smallest execution-time validity/revalidation mechanism justified by each action.

### Final closure amendments from independent adversarial review

An independent adversarial review was adjudicated under the repository method rather than accepted by deference. The operator approved these final amendments:

1. **Composite offers:** real component-dependent semantics are explicit when encountered, but provider-native composite offers are not a launch requirement by default.
2. **Organization posture:** the first proof uses one Organization; Multi-Organization/SaaS operation is not a launch gate; Organization remains a real identity and is not collapsed into a singleton.
3. **Marketplace Installation health/reputation:** provider-native health/reputation may feed portfolio attention as observed/derived evidence; reputation-management/SAC remains out.
4. **Economic chain:** L3 bank cash evidence remains semantically distinct but is an optional extension, not a launch gate.
5. **Document self-containment:** the orphan `(A+)` label was removed, buyer Q&A/chat scope was made unambiguous, and authority-map combined notation was defined.

These amendments do not create new Product 1.0 capabilities or new generic frameworks. They close scope ambiguity and, where applicable, **subtract launch scope**.

### D0.7 stop rule — no open-ended microtopic expansion

D0.7n remains the last planned microtopic. A new D0 question may only be opened if a fresh finding passes all of these tests:

1. there is a real/reachable scenario or material evidence;
2. the answer changes Product 1.0 meaning, responsibility, authority or boundary;
3. if left undecided, an implementer would have to invent business semantics;
4. the question does not clearly belong to D1–D7 implementation/architecture mechanics.

If those tests are not satisfied, the issue is deferred to the owning later stage or removed as accidental complexity. This applies the repository's YAGNI and bounded-review rules to the architecture process itself.

### Anti-overengineering interpretation for later stages

D0 truth/action-safety properties such as freshness, coverage, decision lineage, blast radius and execution-time validity are **cross-cutting invariants**, not instructions to create standalone contexts/modules/frameworks named after those properties.

D1–D7 must choose the smallest mechanism inside the real owning business/technical boundaries. A standalone domain/framework is justified only by new material business evidence establishing an independent responsibility.

---

## 11. Final D0 closure / resume contract

The bounded closure review plus independent adversarial review found no remaining product-semantic blocker after the amendments above.

A fresh session must read `AGENTS.md`, `docs/engineering/rebaseline/README.md`, `docs/engineering/standards/root-cause-global-maximum-method.md`, `ARCHITECTURE.md`, the ADR registry and then this artifact.

Before D0 is finally marked closed, the cold review must confirm:

- the authority documents agree on D0 status and exact next action;
- Product 1.0 remains Marketplace Operations + Commercial Intelligence, not an ERP/marketplace/accounting/task-management/reputation-management replacement;
- the first launch proof requires one Organization but does not erase Organization as an identity;
- composite offers with real composition dependencies are never silently flattened, while composite support is not a launch gate by default;
- Marketplace Installation health/reputation can feed portfolio attention without expanding into SAC/reputation management;
- core economic lineage closes through L2 settlement; L3 bank receipt is optional;
- canonical MPC semantics precede ERP/provider mapping;
- availability derives from eligible authoritative inventory + rules + policy; unknown is not zero;
- provider capability/authority is context-sensitive and claimed MPC-controlled paths close provider requirements;
- marketplace hubs remain benchmark/competitive evidence, not target runtime dependencies;
- external time obligations and MPC-owned Internal Operational Targets remain distinct;
- actionable work has explicit role ownership; assignment is distinct from authorization; material work closes through evidence;
- material evidence is freshness-aware and coverage-aware;
- historical material actions preserve decision-time lineage;
- multi-target actions preserve intended blast radius and granular partial/ambiguous outcomes;
- material decisions/approvals are revalidated only against materially governing conditions at execution time;
- D0.7 remains bounded by the stop rule;
- current code/docs remain evidence, not target authority;
- no D1+ target architecture is invented inside D0.

If those checks pass, D0 is ready to be recorded **ACCEPTED AS A WHOLE / CLOSED**, and the exact next stage is **D1 — Domains / Boundaries**. No product implementation is authorized before D9.
