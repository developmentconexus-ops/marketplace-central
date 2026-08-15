# D0 — Product / System Definition

> **Status:** WORKING — D0 OPEN, NOT YET ACCEPTED AS A WHOLE  
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

Product 1.0 is **Marketplace Operations + Commercial Intelligence (A+)**.

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
  → cash-receipt evidence when available
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
12. **Economic Evidence & Realized Profitability** — preserve an explainable economic evidence chain from simulation through order economics, marketplace settlement and cash-receipt evidence where available; reconcile stage-to-stage variance, calculate attributable realized profitability and create calibration/exception work when evidence proves a model or operational discrepancy.

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

---

## 4. Product 1.0 non-goals safe to defer

Unless later D0 evidence reopens them, Product 1.0 does not require:

- paid ads/media management;
- automated buyer Q&A/chat;
- general CRM/SAC or broad complaint/reputation automation;
- company-wide demand forecasting/purchasing;
- company-wide WMS/TMS/reverse-logistics replacement;
- company-wide accounting, treasury, accounts-payable/receivable or bank reconciliation;
- marketplaces beyond what is needed to prove the first Mercado Livre operating loop;
- unrestricted AI/autonomous commercial action without explicit policy;
- a generic integration framework for speculative providers;
- a universal ERP model;
- runtime dependency on marketplace hubs such as ANYMARKET or Magis5 merely to reach marketplaces MPC is intended to control directly.

Shipment/delivery observation, essential cancellation/return/refund consequences, provider-requirement closure for claimed normal fulfillment paths, materially time-bound obligations, marketplace economic settlement and cash evidence needed to close the marketplace economic chain are **not** deferred merely because their source systems are external.

---

## 5. External-system evidence / target direction

### 5.1 Sankhya evidence

The operator reports an already-proven Sankhya application integration in another application using application credentials (`client id`, `client secret`, `X-Token`; no secret values are recorded here). Operational writes such as order creation and invoicing can be performed through Sankhya's application API; system-owned APIs are preferred for writes rather than direct Oracle writes. DB Explorer/database inspection is also available.

This is D0/D4 evidence, not target transport or canonical domain design. D4 independently ratifies exact read/write capability contracts. `CODEMP`, `CODLOC`, TOPs, native cost variants and other Sankhya structures are integration evidence, **not automatically canonical MPC concepts**.

### 5.2 Marketplace/hub landscape evidence

Current research covers Mercado Livre, Amazon, Magalu, Casas Bahia, Leroy/Mirakl, Shopee, MadeiraMadeira, ANYMARKET, Magis5 and Bling. Detailed observations remain in `EVIDENCE-REGISTER.md`.

**Accepted target direction:** MPC is itself the marketplace operations/control-plane product. For marketplaces it supports, the target direction is an MPC-owned provider integration boundary rather than routing marketplace operation through another marketplace hub by default.

ANYMARKET, Magis5 and similar marketplace hubs are benchmark/competitive evidence, not Product 1.0 runtime dependencies. Bling may later be integrated as an ERP/business system when used in that role. Shared provider technology such as Mirakl may justify technical reuse without changing the business marketplace identity.

---

## 6. Stable constraints carried into D0

- Mercado Livre first.
- Sankhya/Oracle external to MPC.
- Go backend is canonical business execution.
- React frontend is not a second business authority.
- PostgreSQL stores MPC-owned canonical state.
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

### 8.1 Product-level authority map

| Concern | MPC authority |
|---|---|
| Internal master product / ERP-native fiscal and cost facts | **OBSERVE / CONSUME** |
| Native stock/reservation/availability facts in ERP/WMS/3PL/provider | **OBSERVE / CONSUME** |
| Product ↔ marketplace linkage and readiness conclusion | **OWN / DERIVE** |
| Marketplace Installation membership/configuration | **OWN**; provider account identity/state remains provider-owned |
| Selling Entity semantic and transaction attribution | **OWN**; corresponding external legal/company records may remain externally authoritative |
| Inventory Source identity and Inventory Scope eligibility | **OWN** for MPC semantics/configuration; native stock source/facts remain external |
| Availability Allocation Policy | **OWN** when MPC-governed |
| Sellable Availability | **OWN / DERIVE** from eligible authoritative facts/rules/policy |
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
| Bank cash-receipt evidence | bank/accepted bank-data source authoritative; MPC **OBSERVES / RECONCILES** |
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
8. MPC canonical semantics come from marketplace-operating needs, not ERP/provider ontology.
9. ERP/provider integration is semantic translation, not field renaming; unsupported/ambiguous mapping is explicit.
10. Selling Entity, Inventory Source, Fulfillment Node and Cost Basis remain distinct unless explicit business evidence relates them.
11. Inventory Scope is explicit eligibility, not “all stock we can find”.
12. Availability Allocation Policy changes allocation, not authoritative stock truth.
13. MPC-owned policy defaults/overrides resolve deterministically with visible provenance.
14. Cost Observation is not a bare number; Cost Basis is not an ERP cost-type alias.
15. Missing/ambiguous cost does not silently fall back; current cost is not historical transaction cost by default.
16. Cost does not absorb all sale economics; fees, freight, taxes, discounts, subsidies and reversals remain separately attributable when material.
17. Business Order Intent is not an ERP order type/TOP; marketplace order, MPC intent and ERP-native order remain distinct authorities.
18. Missing/ambiguous order mapping is explicit; ambiguous order-write outcome is not blind retry.
19. No separate `Order Execution Scope` is introduced without independent business semantics.
20. Invoicing Intent is not an ERP/native fiscal operation; order existence does not equal invoicing readiness.
21. Unknown invoicing prerequisite does not become ready; ambiguous invoicing outcome is not blind retry.
22. No separate `Fiscal Scope` / `Invoicing Scope` is introduced without independent business evidence.
23. Post-sale fiscal consequences remain controlled/orchestrated without importing native fiscal taxonomy into D0.
24. Economic evidence stages do not overwrite each other.
25. Simulation variance is not automatically a simulator defect; calibration requires materially comparable context.
26. Marketplace settlement is not bank cash receipt.
27. Unattributable financial movement is not invented attribution.
28. Payout/cash matching is not assumed 1:1 with an order.
29. Simulator calibration is evidence-driven.
30. Provider capability/authority is context-sensitive; marketplace brand alone does not determine responsibility.
31. A claimed MPC-controlled normal fulfillment path must close provider requirements.
32. Native provider states/artifacts remain provider-native rather than universal MPC domain types by analogy.
33. Unsupported/external-required provider work is explicit.
34. Marketplace hubs are not Product 1.0 runtime dependencies by default.
35. Shared provider technology does not change business-provider identity.
36. External obligation and internal target are different authorities; internal policy may tighten but never relax the external obligation.
37. Relative time policy requires an explicit trustworthy time anchor.
38. Internal-target breach and external-obligation breach are distinct states.
39. Material time drives operational attention.
40. Actionable work is not operationally ownerless.
41. Assignment is not authorization.
42. Escalation is not notification; it means greater/different responsibility, attention or authority is required.
43. Material work does not close by arbitrary dismissal; accepted resolution evidence is required.
44. Freshness is use-sensitive; D0 does not define a universal age threshold.
45. Stale evidence is not false evidence, but cannot masquerade as current truth.
46. Unknown freshness is not current.
47. Acquisition failure affects operational confidence and remains distinguishable from the last successful observation.
48. Unsafe action may be blocked/degraded by insufficient freshness when current evidence is materially required.
49. **Fresh is not complete.** A newly observed subset does not prove the relevant source population was fully observed.
50. **Not observed is not does-not-exist.** Absence from a partial/unknown-coverage observation cannot become a terminal or negative fact merely because a record was not returned.
51. **Coverage is scoped.** Any claim of sufficient completeness must be relative to an explicit population/range/scope/obligation rather than a global `providerComplete=true` assertion.
52. **Individual observed facts do not prove population completeness.** A fact can be valid while portfolio-level coverage remains partial/unknown.
53. **Portfolio health and reconciliation closure require sufficient coverage for the claimed universe.** Partial/unknown coverage cannot truthfully produce `no issues`, `nothing missing` or `reconciled` merely because observed subsets agree.
54. **Webhook/callback receipt is not completeness proof.** Notification delivery may contribute evidence but cannot alone prove every relevant change/record was observed unless the provider contract independently guarantees that property and D4 verifies it.
55. **Current state is not historical decision basis.** Material historical explanation must use decision-time evidence/provenance rather than mutable current values.
56. **Current policy is not necessarily the policy that governed a historical action.** Historical explanation preserves the governing policy/rule meaning/provenance effective for that decision.
57. **Action audit is not decision explanation.** Knowing who/what invoked an action is insufficient without enough evidence, authority/approval and uncertainty context to explain why it was permitted/recommended/executed.
58. **Re-running a current calculation is not reproducing a historical calculation.** Simulator, cost, fee, rule or policy changes must not rewrite the meaning of a prior decision.
59. **Material external action remains correlatable to intent + authority + decision-time evidence + result.** The explanation must survive later mutation of current state.
60. D0 does not require event sourcing, universal snapshots or copying every source fact. D2/D3/D7 later choose the smallest retention/reference/snapshot mechanism that satisfies accepted historical explainability.
61. **Multi-target action is not one opaque boolean result.** A material action affecting multiple targets preserves an intended target scope and sufficiently granular outcomes to distinguish confirmed, rejected, ambiguous and not-executed members where those distinctions are material.
62. **Intended target scope is not reconstructed from mutable current filters.** Historical blast radius must remain explainable from the target universe intended/authorized for that action, even if current eligibility/filter membership changes later.
63. **Batch-level authorization does not erase target-level authority/readiness differences.** A collective decision cannot silently authorize a member that independently requires different authority, approval, evidence or readiness.
64. **Partial success is neither total success nor total failure.** Confirmed effects remain confirmed while rejected/ambiguous/not-executed members remain independently visible/reconcilable.
65. **An ambiguous member does not make whole-batch retry safe.** MPC must not blindly retry already-confirmed members merely because other members failed or became uncertain.
66. **MPC does not invent cross-provider atomicity.** Unless an accepted provider contract proves atomic behavior, external multi-target effects may be partial and Product 1.0 must remain correct under partial outcomes.
67. **Automation at scale does not hide blast radius.** Policy-driven mass effects remain observable/reconcilable even when no per-target human approval is required.
68. Provider/ERP/payment/bank transport, timer, scheduler, polling, webhook, TTL/cache, freshness threshold, pagination, checkpoint, backfill, coverage-proof, decision-record, batch API, queue, concurrency, compensation and retry-orchestration mechanics remain later-stage responsibilities.

---

## 9. D0.6 — Product 1.0 completion / user-observable outcomes

**Accepted by operator; refined by later D0.7 decisions.**

Product 1.0 is complete only when MPC is demonstrably usable as the normal marketplace operations control plane.

1. **Attention is portfolio-driven.** Healthy/changed/divergent/blocked/approval-required/actionable state is visible without external manual product-by-product checking.
2. **An eligible internal product reaches verified marketplace state through MPC.** Readiness/linkage → commercial analysis → creation/publication → real channel observation.
3. **Inventory eligibility/allocation is explicit and availability converges.** Eligible Inventory Sources + authoritative facts/rules + applicable MPC policy derive Sellable Availability; uncertainty/failure is explicit.
4. **Competitive/pricing intelligence replaces major manual analysis with explainable economics.** Cost Basis and evidence are visible; insufficient evidence is honest.
5. **Decision closes into controlled action.** Human or bounded policy decisions become auditable external actions with verification/reconciliation.
6. **Material business context is explicit.** Selling Entity, installation and other proven semantics are not hidden defaults.
7. **Fulfillment responsibility is explicit.** Eligible/selected Fulfillment Node and physical responsibility are known when material.
8. **A marketplace sale becomes the correct business-system order through MPC.** Business Order Intent is mapped, executed, correlated and verified without leaking native order types into canonical semantics.
9. **Invoicing is readiness-gated and verifiable.** Invoicing Intent is emitted only after sufficient evidence; native result is observed/reconciled.
10. **Provider requirements close inside claimed normal fulfillment paths.** MPC determines effective provider requirements/capabilities, surfaces/orchestrates required data/artifacts/acknowledgements and proves provider-specific readiness without hidden routine provider-UI work.
11. **The normal sale lifecycle runs through MPC** for responsibilities it claims to control.
12. **Material time-bound obligations drive attention.** External deadlines and organization-owned internal targets remain distinguishable and relative targets have known anchors.
13. **Actionable work has an owner and can progress/escalate.** Assignment is supported when useful and never bypasses authorization.
14. **Work closes through evidence.** Material unresolved divergence cannot disappear through arbitrary dismissal.
15. **Material evidence is freshness-aware.** Stale/unknown-freshness evidence and degraded acquisition affect confidence/attention and can block unsafe normal action.
16. **Portfolio/absence/reconciliation claims are coverage-aware.** MPC distinguishes sufficiently-covered, partial and unknown coverage for the relevant population; `no issues`, terminal absence and reconciliation closure require enough coverage evidence for the claim being made.
17. **Material historical decisions/actions remain explainable.** MPC can explain why a consequential recommendation/approval/action was permitted or executed using decision-time evidence, governing policy/rule provenance, authority/approval and known uncertainty rather than reconstructing history from mutable current state.
18. **Material multi-target actions expose blast radius and partial outcomes.** MPC preserves the intended target scope and can distinguish confirmed, rejected, ambiguous and not-executed effects sufficiently to reconcile partial execution without pretending the whole operation was atomically successful/failed.
19. **Shipment remains visible through terminal outcome.** Material delivery exceptions become explicit work.
20. **Essential post-sale changes remain controlled.** Cancellation/return/refund and required fiscal/economic consequences remain orchestrated/reconciled.
21. **Failures become explicit work.** Missing evidence, ambiguity, unsupported capability, deadline uncertainty/breach, stale evidence, partial/unknown coverage, partial action outcome and divergence identify what is known/unknown, who is responsible and what action/authority is required.
22. **The economic evidence chain closes as far as authoritative evidence is available.** Expected Economics → Order Economics → Marketplace Settlement → Cash Receipt evidence remain distinct.
23. **Simulator drift becomes actionable learning** only after materially comparable historical evidence isolates model/provider-rule drift.
24. **Realized profitability is explainable** from attributable economic evidence rather than convenient current values.
25. **Organizational governance is operable without code edits** while external authority and mandatory safety/audit invariants remain protected.

Completion statement:

> **A company can take internal products, determine marketplace readiness, publish and operate offers, maintain availability from explicit inventory scope/policy, monitor competition and profitability, make controlled decisions from sufficiently current and sufficiently covered evidence, retain enough decision-time lineage to explain material historical actions, execute at scale without hiding blast radius or partial outcomes, materialize sales in business/fiscal systems, fulfill them through eligible nodes, close provider requirements, prioritize time-bound obligations, route actionable work to explicit owners, close work through evidence, dispatch and observe delivery/post-sale consequences, preserve economic evidence through settlement/cash where available, reconcile divergence without mistaking partial observation or partial execution for closure, learn when its own simulator is wrong, and explain realized economic outcome — using Marketplace Central as the normal marketplace operations control plane.**

### 9.1 Normal-path rule

The normal operational path must be executable through MPC for responsibilities Product 1.0 claims to control. Direct use of an external system remains legitimate only for intentionally external responsibilities, investigation/support and explicit exceptional recovery; hidden routine system hopping means the claimed workflow is incomplete.

If a provider-required operation cannot be performed through an accepted Product 1.0 integration, that limitation must be explicit; the affected path cannot be presented/proven as fully MPC-controlled.

D8 later proves detailed golden flows; D0 defines what those flows must prove.

---

## 10. D0.7 — Product completeness / contradiction review

D0.7 adversarially tests Product 1.0 for missing lifecycle responsibility, contradiction and accidental scope gap.

### D0.7a — Essential post-sale lifecycle

**Accepted.** Cancellations, returns and refunds tied to MPC-controlled sales remain inside the controlled lifecycle, including required cross-system/fiscal/economic consequences. General CRM/SAC, buyer messaging, reputation automation and company-wide reverse logistics remain outside.

### D0.7b — Shipment / delivery lifecycle

**Accepted.** MPC observes relevant shipment/delivery state after dispatch through terminal outcome and turns material delay/failure/return-to-sender states into explicit work without becoming TMS/carrier.

### D0.7c — Stock / marketplace availability control

**Accepted.** Sellable Availability is derived from authoritative stock/reservation facts, Inventory Scope, governing rules and applicable MPC-owned allocation policy; routine known policy-valid synchronization is automatic and verified. Unknown availability does not become zero/plausible quantity.

### D0.7d — Marketplace installation multiplicity

**Accepted.** `Organization 1 → N Marketplace Installations`. First proof may use one Mercado Livre account, but organization identity never collapses into seller-account identity.

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

**Accepted.** Product 1.0 preserves:

```text
L0 Simulation / Expected Economics
  → R1 simulation reconciliation
L1 Order Economics
  → R2 settlement reconciliation
L2 Marketplace/payment-account Settlement
  → R3 cash reconciliation
L3 Bank Cash Receipt evidence, when available
```

Stages do not overwrite each other. Simulation/order variance is classified before becoming a calibration defect. Settlement is distinct from bank receipt. Payouts may aggregate many orders/movements and unattributable movement is never fabricated. This closes marketplace economic lineage without turning MPC into company-wide finance/accounting/treasury.

### D0.7g — Context-sensitive provider capabilities / Provider Requirement Closure

**Accepted.** Marketplace brand alone does not determine capability/authority; effective responsibility can vary by installation/offer/order/native operating mode. Every flow claimed as MPC-controlled must surface/satisfy/orchestrate required provider prerequisites/data/artifacts/readiness or explicitly mark the path unsupported/external-required. Native provider artifacts/states remain provider-native.

MPC targets direct provider boundaries. ANYMARKET, Magis5 and similar hubs are benchmark/competitive evidence, not Product 1.0 runtime dependencies. Shared provider technology such as Mirakl may justify implementation reuse without changing business marketplace identity.

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

**Accepted by operator.**

For material operator- or policy-driven actions that affect multiple marketplace targets, Product 1.0 preserves an explicit **intended target scope / blast radius** and sufficiently granular execution outcome semantics so differences and partial effects remain observable/reconcilable rather than being collapsed into one opaque batch status.

Product-level requirements:

- the action's intended target universe is attributable to the decision/authorization that created it; later changes in filters, eligibility, policy membership or current state must not silently redefine what that historical action intended to affect;
- collective/batch authorization does not erase material target-level differences in authority, policy, evidence sufficiency, readiness or provider capability;
- the execution/reconciliation view can distinguish, where material, targets/effects that are **confirmed**, **rejected/failed**, **ambiguous/unknown outcome** and **not executed/not attempted**;
- confirmed target effects remain confirmed even when other members fail; partial success is neither total success nor total failure;
- ambiguous or failed members do not make blind retry of the whole intended set safe because already-confirmed targets may otherwise be duplicated/reapplied;
- MPC does not promise cross-target/provider all-or-nothing atomicity unless a later verified provider contract actually guarantees it;
- policy-driven automation that affects many targets is subject to the same blast-radius/outcome observability and reconciliation property even when routine policy-valid execution does not require per-target human approval;
- D0 does not mandate a `BatchJob`/`BulkOperation` entity, per-target database row, saga, compensation engine, chunking, worker pool, progress UI or a specific retry algorithm.

This is a product safety/reconciliation property. D5/D6/D7 later choose the smallest bulk-action contract, execution and presentation mechanics that preserve scope, authority and partial outcomes.

### Next D0.7 question

The next material product-completeness question is **decision/approval validity at execution time**.

D0 now preserves decision-time lineage and multi-target blast radius, but a consequential action may execute after the recommendation/approval that authorized it. Between approval and execution, authoritative facts, provider capability/readiness, evidence freshness/coverage, MPC-owned policy or an external obligation may materially change.

D0 must decide whether Product 1.0 treats an earlier approval/decision as permanently valid, or requires material execution-time preconditions/authority to remain valid enough for the action actually being executed. The likely target property is that approval is authorization under a decision context—not an eternal waiver of later-invalidated safety/policy conditions. This is a product safety/authority question, not a locking, optimistic-concurrency, version token or reapproval UI decision. Exact revalidation/versioning mechanics belong to D2/D5/D6/D7.

Other remaining D0.7 findings continue to be classified adversarially; D0 closes only when no material Product 1.0 semantic is left for implementation to invent.

---

## 11. Resume contract for a fresh session

A fresh session must read `AGENTS.md`, `docs/engineering/rebaseline/README.md`, `docs/engineering/standards/root-cause-global-maximum-method.md`, `ARCHITECTURE.md`, the ADR registry and then this artifact.

It should conclude:

- D0 is OPEN and not yet accepted as a whole; implementation remains blocked until D9;
- D0.1–D0.6 and D0.7a–D0.7m are operator-approved;
- Product 1.0 is Marketplace Operations + Commercial Intelligence (A+), not an ERP/marketplace/accounting/task-management replacement;
- organization, Marketplace Installation and Selling Entity identities do not collapse;
- canonical MPC semantics precede ERP/provider mapping;
- Selling Entity, Inventory Source/Scope, Fulfillment Node/Scope, Cost Observation/Basis, Business Order Intent and Invoicing Intent are accepted semantics;
- availability derives from eligible authoritative inventory + rules + policy; routine known policy-valid synchronization is automatic; unknown is not zero;
- the Economic Evidence Chain preserves separate simulation/order/settlement/cash evidence and does not fabricate payout/order attribution;
- provider capability/authority is context-sensitive and claimed MPC-controlled paths close provider requirements without hidden routine provider-UI work;
- marketplace hubs remain benchmark/competitive evidence, not target runtime dependencies;
- external time obligations and MPC-owned Internal Operational Targets are distinct; relative targets require explicit anchors and internal policy cannot relax external obligations;
- actionable work has explicit role ownership; assignment is distinct from authorization; escalation changes required attention/responsibility/authority; material work closes through evidence;
- material evidence is freshness-aware and coverage-aware: stale/unknown-freshness evidence cannot masquerade as current truth, fresh is not complete, and portfolio/absence/reconciliation claims require sufficient coverage;
- material historical recommendations/approvals/actions preserve enough decision-time lineage to explain their evidence, policy/rule provenance, authority/approval and uncertainty without relying on mutable current state;
- material multi-target actions preserve their intended blast radius and sufficiently granular partial/ambiguous outcomes; batch-level success/failure or whole-batch blind retry cannot erase already-confirmed or unresolved members;
- current code/docs remain evidence, not target authority;
- no D1+ target architecture may be invented yet;
- the exact next work is **D0.7 Product completeness review — decision/approval validity at execution time**.
