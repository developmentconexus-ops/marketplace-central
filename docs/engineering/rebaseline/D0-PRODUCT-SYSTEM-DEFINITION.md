# D0 — Product / System Definition

> **Status:** WORKING — D0 OPEN, NOT YET ACCEPTED AS A WHOLE  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Authority:** operator-approved D0 decisions recorded here; later D-stages own technical realization  
> **Last updated:** 2026-08-14

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

A governing rule used by MPC may be:

- **MPC-owned** — intentionally defined/governed inside MPC;
- **externally governed** — authoritative in an ERP/provider/other system and consumed by MPC;
- **derived** — mechanically concluded from authoritative facts/rules without becoming an independent source of truth.

MPC must preserve provenance and must not silently convert an external rule into an editable MPC-owned copy.

Where a policy is legitimately MPC-owned, Product 1.0 must allow organization-operable configuration without code edits. The policy model must support deterministic inherited/default policy plus explicit more-specific overrides where later-proven business scopes justify them; effective policy and provenance must be explainable. D0 does not mandate every candidate scope or create a generic speculative rules engine.

This same rule applies to internal operational-time targets: an organization may deliberately choose a stricter target than a provider deadline, or define an internal target where no external deadline exists. Internal targets remain MPC-owned policy and never rewrite or relax an externally governed obligation.

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

Essential shipment/delivery observation, cancellation/return/refund consequences, provider-requirement closure for claimed normal fulfillment paths, materially time-bound operational obligations, marketplace economic settlement and the cash evidence needed to close the marketplace economic chain are **not** deferred merely because their source systems are external.

---

## 5. External-system evidence / target direction

### 5.1 Sankhya evidence

This is **D0/D4 evidence**, not target transport or canonical MPC domain design.

The operator reports an already-proven Sankhya application integration in another application using application credentials (`client id`, `client secret`, `X-Token`; no secret values are recorded here). Operational writes such as order creation and invoicing can be performed through Sankhya's application API; the operator prefers system-owned APIs for writes rather than direct Oracle writes. DB Explorer/database inspection is also available.

Consequences:

- Product 1.0 may require order/invoicing operations without assuming direct database writes;
- D4 must independently ratify exact read/write capability contracts and transports;
- externally governed Sankhya rules remain external authority;
- existing binding Oracle-read constraints are not silently reopened;
- `CODEMP`, `CODLOC`, TOPs, native cost variants and other Sankhya structures are integration evidence, **not automatically canonical MPC concepts**.

### 5.2 Marketplace/hub landscape evidence

Current provider and competitor research covers Mercado Livre, Amazon, Magalu, Casas Bahia, Leroy/Mirakl, Shopee, MadeiraMadeira, ANYMARKET, Magis5 and Bling. The detailed observations live in `EVIDENCE-REGISTER.md`; this artifact records only accepted product implications.

**Accepted target direction:** MPC is itself the marketplace operations/control-plane product. For marketplaces Product 1.0 or later releases choose to support, the target direction is an MPC-owned provider integration boundary rather than routing marketplace operation through another marketplace hub by default.

ANYMARKET, Magis5 and similar marketplace hubs are benchmark/competitive evidence, not Product 1.0 runtime dependencies. Bling may only become a future business-system/ERP integration when used in that role; its marketplace-hub role is not a required dependency.

This does not preclude future evidence from justifying an intermediary integration, but no intermediary abstraction or compatibility tax is authorized in the current target design.

### 5.3 Time-bound obligation evidence

Current provider research confirms that marketplace operations can expose authoritative deadlines, not merely descriptive timestamps. For example, Mercado Livre exposes a maximum shipment-dispatch date/time whose breach can affect reputation/visibility. External deadlines therefore have operational consequences and must remain distinguishable from internal organization targets.

Mature reliability practice also supports keeping an internal target stricter than an external commitment to create operating safety margin. D0 adopts the business-semantic principle without importing SRE-specific implementation machinery.

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

**Invariant:** the accepted internal normal path does not intentionally invoice before physical confirmation that the correct items are available/separated. Physical inconsistency blocks normal invoicing and becomes explicit exception work.

Externally operated Fulfillment Nodes may prove readiness differently in later stages; Product 1.0 need not implement every fulfillment mode in its first deployment.

### Commercial / Marketplace Manager

Ordinary authority over legitimate **MPC-owned** commercial policies: margin floors, price boundaries, approval thresholds, availability-allocation policy, bounded automation, internal operating targets and commercial exceptions. This actor is not automatically authority over externally governed ERP/provider rules or integration/security administration. Exact delegation over each policy class belongs later.

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
| Effective provider capability/requirement interpretation for a claimed workflow | **OWN / DERIVE** from provider-native mode/state/contract evidence; provider remains authority for native capability/state/artifact semantics |
| Provider-required native fulfillment prerequisites, acknowledgements, shipment readiness and artifacts | provider authoritative where native; MPC **OBSERVES / ORCHESTRATES / RECONCILES** closure inside claimed normal paths |
| Externally governed operational deadline/window | provider/other external authority; MPC **OBSERVES** and preserves source/time evidence |
| MPC-owned Internal Operational Target | **OWN** as organization policy; it may be stricter than an external obligation but cannot rewrite/relax it |
| Operational urgency, remaining time, safety margin and breach/attention interpretation | **OWN / DERIVE** from authoritative obligation/target evidence and current workflow state |
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
| Economic Evidence Chain, variance classification, simulator-calibration cases | **OWN / DERIVE** |
| Realized profitability interpretation | **OWN / DERIVE** from attributable economic evidence |
| MPC audit/reconciliation/exception state | **OWN** |

### 8.2 Boundary invariants

1. Observing an external fact does not transfer ownership of it to MPC.
2. Orchestrating an external process does not make MPC source of truth for its native record.
3. Derived conclusions preserve enough provenance/freshness/context to explain their evidence.
4. Externally governed rules are not silently copied into mutable MPC policy.
5. Ambiguous/divergent outcomes remain explicit until reconciled.
6. **Unknown availability is not zero** or another plausible quantity.
7. Organization, Marketplace Installation and Selling Entity are not the same identity merely because a first deployment uses one of each.
8. MPC canonical semantics come from marketplace-operating needs, not ERP/provider ontology.
9. ERP/provider integration is semantic translation, not field renaming; unsupported/ambiguous mapping is explicit.
10. Selling Entity, Inventory Source, Fulfillment Node and Cost Basis remain distinct semantics unless explicit business evidence relates them.
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
24. **Economic evidence stages do not overwrite each other.** Simulation, order economics, marketplace settlement and cash evidence remain distinguishable/provenanced layers.
25. **Simulation variance is not automatically a simulator defect.** A calibration defect requires comparable material context or an explanation of changed inputs/rules/facts.
26. **Marketplace settlement is not bank cash receipt.** Provider/payment-account movement and bank-side credit are different evidence layers.
27. **Unattributable financial movement is not invented attribution.** Unknown/aggregate adjustments remain unresolved rather than being spread across sales merely to close arithmetic.
28. **Payout/cash matching is not assumed 1:1 with an order.** Aggregated many-sale/many-movement payouts must preserve lineage and uncertainty.
29. **Simulator calibration is evidence-driven.** Confirmed model/rule drift creates explicit calibration work; explainable context variance does not masquerade as a model bug.
30. **Provider capability/authority is context-sensitive.** Marketplace brand alone does not prove what the provider/seller/MPC can or must do for an offer/order; effective capability may depend on installation, native operating/fulfillment mode and accepted provider contract evidence.
31. **A claimed MPC-controlled normal fulfillment path must close provider requirements.** Required provider data handoffs, acknowledgements, readiness and artifacts must be surfaced, satisfied or orchestrated through MPC; a hidden routine provider-UI step means that path is incomplete.
32. **Native provider states/artifacts remain native.** Labels, provider fiscal documents, shipment states and analogous provider artifacts do not become universal MPC domain types merely because several providers expose analogous constructs.
33. **Unsupported/external-required provider work is explicit.** If a required step cannot be supported through the accepted integration, the workflow may be explicitly limited/exceptional but must not masquerade as fully MPC-executable.
34. **Marketplace hubs are not Product 1.0 runtime dependencies by default.** ANYMARKET, Magis5 and similar systems are competitive/architectural evidence; no intermediary compatibility layer is introduced without new business evidence.
35. **Shared provider technology does not change business-provider identity.** A future shared protocol/family implementation such as Mirakl may reuse technical infrastructure, but the marketplace/business provider remains the marketplace being operated.
36. **External obligation and internal target are different authorities.** An organization may choose an earlier/stricter Internal Operational Target, but that target cannot rewrite, waive or relax a provider/external obligation.
37. **Relative time policy requires an explicit time anchor.** A rule such as “act within one hour of receipt” is only meaningful when the triggering event/time and provenance are known; unknown anchor/deadline is not equivalent to no deadline.
38. **Internal-target breach and external-obligation breach are distinct states.** Missing an internal safety target must not be reported as provider breach; meeting an internal target does not by itself prove an external obligation was satisfied.
39. **Time drives attention, not hidden background arithmetic.** Material remaining time, safety margin, risk and breach state must be available to operational prioritization rather than existing only as raw timestamps.
40. Provider/ERP/payment/bank transport, timer and scheduler mechanics remain later-stage responsibilities.

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
10. **Provider requirements close inside claimed normal fulfillment paths.** MPC determines the effective provider requirements/capabilities for the supported operating context, surfaces or orchestrates required data/artifacts/acknowledgements and proves provider-specific readiness needed to reach dispatch handoff without hidden routine provider-UI work.
11. **The normal sale lifecycle runs through MPC.** Order materialization → fulfillment/conference → invoicing → provider-requirement closure → packing/dispatch without hidden routine system hopping for responsibilities the product claims to control.
12. **Material time-bound obligations drive attention.** Provider/external deadlines and organization-owned internal targets remain distinguishable, relative targets have a known anchor, and operational attention can prioritize approaching/breached obligations rather than requiring operators to discover deadlines manually.
13. **Shipment remains visible through terminal outcome.** Material delivery exceptions become explicit work.
14. **Essential post-sale changes remain controlled.** Cancellation/return/refund and required fiscal/economic consequences remain orchestrated/reconciled.
15. **Failures become explicit work.** Missing evidence, ambiguity, unsupported capability, deadline uncertainty/breach and divergence identify what is known/unknown and what action is required.
16. **The economic evidence chain closes as far as authoritative evidence is available.** MPC preserves and compares Expected Economics → Order Economics → Marketplace Settlement → Cash Receipt evidence, without collapsing distinct stages into one opaque number.
17. **Simulator drift becomes actionable learning.** Materially comparable simulation/order evidence can identify model/provider-rule drift and create calibration work; explainable input/context changes are classified separately.
18. **Realized profitability is explainable.** The system can attribute material cost/fees/freight/tax/post-sale/settlement effects and distinguish unresolved economics instead of recomputing history from convenient current values.
19. **Organizational governance is operable without code edits.** Legitimate MPC-owned authorities, including internal operating-target policy, are configurable while external authority and mandatory safety/audit invariants remain protected.

Completion statement:

> **A company can take internal products, determine marketplace readiness, publish and operate offers, maintain availability from explicit inventory scope/policy, monitor competition and profitability, make and execute controlled decisions, materialize marketplace sales in business/fiscal systems, fulfill them through explicit eligible nodes, close provider-required prerequisites/artifacts/readiness for the supported operating path, prioritize and meet time-bound marketplace/business obligations under both external deadlines and organization-owned internal targets, dispatch and observe delivery/post-sale consequences, preserve economic evidence from simulation through order/settlement/cash where available, reconcile divergence, learn when its own simulator is wrong, and explain realized economic outcome — using Marketplace Central as the normal marketplace operations control plane.**

### 9.1 Normal-path rule

The normal operational path must be executable through MPC for responsibilities Product 1.0 claims to control.

Direct use of a marketplace, ERP or another external system remains legitimate for intentionally external responsibilities, investigation/support and explicit exceptional recovery. It must **not** be a hidden routine step in a workflow MPC claims to control.

If a provider-required operation cannot be performed through an accepted Product 1.0 integration, that limitation must be explicit; the affected path cannot be presented or proven as a fully MPC-controlled normal path.

D8 later proves detailed golden flows; D0 defines what those flows must prove.

---

## 10. D0.7 — Product completeness / contradiction review

D0.7 adversarially tests Product 1.0 for missing lifecycle responsibility, contradiction and accidental scope gap. Each finding is classified as MUST DECIDE NOW, SHOULD DECIDE NOW or CAN DEFER SAFELY.

### D0.7a — Essential post-sale lifecycle

**Accepted.** Cancellations, returns and refunds tied to MPC-controlled sales remain inside the controlled lifecycle, including required cross-system/fiscal/economic consequences. General CRM/SAC, buyer messaging, reputation automation and company-wide reverse logistics remain outside.

### D0.7b — Shipment / delivery lifecycle

**Accepted.** MPC observes relevant shipment/delivery state after dispatch through terminal outcome and turns material delay/failure/return-to-sender states into explicit work without becoming TMS/carrier.

### D0.7c — Stock / marketplace availability control

**Accepted.** Sellable Availability is derived from authoritative stock/reservation facts, Inventory Scope, governing rules and applicable MPC-owned allocation policy; routine known policy-valid synchronization is automatic and verified. Unknown availability does not become zero/plausible quantity.

### D0.7d — Marketplace installation multiplicity

**Accepted.** `Organization 1 → N Marketplace Installations`. First proof may use one Mercado Livre account, but organization identity must never collapse into seller-account identity. Channel-dependent facts/workflows retain installation attribution.

### D0.7e — ERP-agnostic business semantics before ERP mapping

**Accepted.** MPC defines canonical semantics from marketplace business needs first; adapters translate them to/from Sankhya or another ERP. Native ERP fields/entities do not create canonical MPC concepts by existence alone, and unsupported semantic mappings remain explicit.

#### D0.7e.1 — Selling Entity

**Accepted.** `Organization 1 → N Selling Entities`. Selling Entity answers which business/legal/fiscal entity is acting when material. It is not an ERP company identifier, marketplace installation, Inventory Source, Fulfillment Node or Cost Basis.

#### D0.7e.2 — Inventory Source / Inventory Scope

**Accepted.** Inventory Source represents a business-recognized inventory source/pool; Inventory Scope governs which sources may contribute to an offer. Native ERP/WMS/provider mapping need not be 1:1. Stock outside scope does not leak into availability.

##### D0.7e.2a — Availability Allocation Policy

**Accepted product requirement; detailed catalog deferred.** MPC may intentionally expose less than full eligible stock, e.g. `eligible=100`, allocation policy `70%`, derived marketplace allocation based on that policy. The policy changes allocation, not authoritative stock. Defaults/overrides must be deterministic/explainable; exact policy types, arithmetic, rounding and final scope hierarchy belong later.

#### D0.7e.3 — Fulfillment Node / Fulfillment Scope

**Accepted.** Fulfillment Node is a recognized physical-fulfillment execution point/capability, internal or external; Fulfillment Scope governs eligibility. Inventory promise and fulfillment execution are distinct semantics even when one facility performs both. First deployment may use one internal node.

#### D0.7e.4 — Cost Observation / Cost Basis

**Accepted.** Cost Observation preserves economic meaning, amount/currency, time/business context and provenance. Cost Basis expresses which cost semantic is appropriate for a decision/transaction and is not an ERP cost-type alias. Unsupported/ambiguous cost remains explicit; historical economics does not silently substitute current cost.

#### D0.7e.5 — Business Order Intent

**Accepted.** Business Order Intent is MPC-owned intent/correlation for materializing a marketplace sale in the participating business system. ERP TOP/document/operation codes stay in semantic mapping. Missing mapping is explicit; ambiguous write is not blind retry. No separate `Order Execution Scope` is justified yet.

#### D0.7e.6 — Invoicing Intent

**Accepted.** Invoicing Intent is MPC-owned readiness-gated fiscal/documentary materialization intent. Order existence is not readiness. Native fiscal documents remain externally authoritative; an API success response alone does not prove materialization. Ambiguous invoicing is not blind retry. No separate `Fiscal/Invoicing Scope` is justified yet. Material post-sale fiscal consequences remain orchestrated/reconciled.

### D0.7f — Economic Evidence Chain / settlement / cash / simulator calibration

**Accepted by operator after adversarial review and current-source research.**

Product 1.0 preserves four distinct progressive evidence layers:

```text
L0 — Simulation / Expected Economics
       ↓ R1: simulation reconciliation
L1 — Order Economics
       ↓ R2: settlement reconciliation
L2 — Marketplace / payment-account Settlement
       ↓ R3: cash reconciliation
L3 — Bank Cash Receipt evidence, when available
```

These layers form an **Economic Evidence Chain**; later evidence refines what is known but does not overwrite earlier evidence/history.

Simulation/order variance is classified before becoming a simulator defect. Explained input/context variance, facts knowable only at/after sale, provider-rule drift, confirmed model drift and insufficient evidence remain distinct causes. Only materially comparable residual discrepancy becomes a `Simulator Calibration Case`.

Order Economics is reconciled with authoritative marketplace/payment settlement evidence. Explained fees/refunds/chargebacks/adjustments remain attributable; unexplained or unassignable differences become explicit financial exceptions.

Marketplace payout/withdrawal evidence is distinct from bank-side cash receipt. Payouts may aggregate many orders/movements; MPC does not invent `Order 1 → Bank Transaction 1` or proportionally fabricate attribution merely to close totals.

This closes the **marketplace operation's own economic lineage** without expanding MPC into company-wide accounting, treasury or bank reconciliation. Exact Mercado Pago/report APIs, SFTP, Open Finance/bank integration, statement formats, matching algorithms/tolerances/windows and persistence belong later.

### D0.7g — Context-sensitive provider capabilities and Provider Requirement Closure

**Accepted by operator after broad current-provider/competitor research.**

Marketplace fulfillment/dispatch responsibility is not determined by marketplace brand alone. The effective operating contract can vary by marketplace installation, offer/order context and provider-native fulfillment/logistics mode. Product 1.0 must therefore reason from the **effective provider capabilities, authorities and requirements for the supported operating flow**, not from a universal assumption that every marketplace exposes the same workflow.

For every marketplace operating flow Product 1.0 claims as an **MPC-controlled normal path**, MPC must determine enough provider context to know seller/provider/MPC responsibility; surface, satisfy or orchestrate required prerequisites/data handoffs/artifacts; observe provider-specific readiness/acknowledgement; and turn missing capability/rejected handoff/non-convergence into explicit work.

Native labels, provider fiscal documents, shipment/fulfillment statuses, manifests, QR codes or analogous artifacts/states remain provider-native authority. D0 does not create universal provider artifact/status types merely because several providers expose analogous constructs.

A provider-required manual UI step may remain explicitly unsupported/external-required or exceptional, but cannot remain a hidden routine step in a flow claimed as fully executable through MPC.

Product 1.0 only needs end-to-end closure for the Mercado Livre operating mode(s) selected for the first real proof; it does not implement every future marketplace mode merely to preserve architectural correctness.

#### Direct marketplace target direction / competitor-hub boundary

MPC is itself the marketplace operations/control-plane product. ANYMARKET, Magis5 and similar marketplace hubs are **benchmark/competitive evidence, not target runtime dependencies**. The accepted direction is MPC-owned direct provider integration boundaries for marketplaces it supports. No generic intermediary compatibility layer is introduced without future independent business evidence.

Shared provider technology such as Mirakl may later justify technical reuse without changing the business marketplace identity or pretending marketplace contracts are identical.

### D0.7h — Time-bound operational obligations and organization-owned internal targets

**Accepted by operator, including the organization-policy refinement.**

For responsibilities Product 1.0 claims to control, materially time-bound marketplace/business obligations must participate explicitly in MPC operational semantics rather than remaining timestamps operators must discover manually in provider UIs.

MPC distinguishes at least two authority classes:

1. **External Operational Obligation** — a provider/contract/external-business deadline or permitted action window. Its deadline/window and source remain externally governed facts.
2. **Internal Operational Target** — an MPC-owned organization policy expressing when the organization wants the operation completed or responded to. It may deliberately be stricter/earlier than the external obligation to create safety margin, and it may exist even when no external deadline exists.

Example:

```text
provider obligation: dispatch within 2 days
organization target: dispatch within 1 day
```

The internal target does not replace the provider deadline. MPC preserves both and derives the operational safety margin/attention state from them.

Time policy may also be **relative to an explicit operational anchor**, for example:

```text
trigger/anchor: event received at 10:00
organization policy: act within 1 hour
internal target: 11:00
```

The exact event classes in scope are determined by accepted Product 1.0 workflows. The example does **not** reopen automated buyer chat/Q&A, which remains a non-goal unless separately re-adjudicated.

Product-level requirements:

- provider/external obligations and organization targets retain separate authority/provenance;
- an Internal Operational Target may tighten an external obligation but cannot waive, overwrite or make a later external deadline acceptable;
- internal policy may use organization defaults and later-proven more-specific scopes under the deterministic inheritance/override principles already accepted for MPC-owned policy;
- a relative target requires an explicit, sufficiently trustworthy time anchor. Unknown trigger/deadline is explicit uncertainty, not `no deadline`;
- MPC may derive remaining time, safety margin, urgency, risk, internal-target breach and external-obligation breach for portfolio attention;
- missing an internal target and breaching an external obligation are different operational facts and must not be conflated;
- exact provider deadline fields, timezones, calendars, cutoffs, timers, schedulers, notifications, escalation thresholds and UI countdown mechanics belong to D4/D6/D7.

This is a cross-cutting control-plane property, not a separate business capability or a generic SLA engine.

### Next D0.7 question

The next material product-completeness question is **operational work ownership / assignment / escalation**.

D0 already requires failures, approvals, ambiguities, calibration cases and approaching/breached obligations to become explicit work. A control plane can still fail if actionable work is visible but nobody is explicitly responsible for progressing it. D0 must therefore decide whether Product 1.0 needs explicit operational ownership/assignment and escalation semantics for actionable cases, while leaving queues, notification channels, routing algorithms and UI mechanics to later stages.

Other remaining D0.7 findings continue to be classified adversarially; D0 closes only when no material Product 1.0 semantic is left for implementation to invent.

---

## 11. Resume contract for a fresh session

A fresh session must read `AGENTS.md`, `docs/engineering/rebaseline/README.md`, `docs/engineering/standards/root-cause-global-maximum-method.md`, `ARCHITECTURE.md`, the ADR registry and then this artifact.

It should conclude:

- D0 is OPEN and not yet accepted as a whole; implementation remains blocked until D9;
- D0.1–D0.6 and D0.7a–D0.7h are operator-approved;
- Product 1.0 is Marketplace Operations + Commercial Intelligence (A+), not an ERP/marketplace/accounting replacement;
- `Organization 1 → N Marketplace Installations` and `Organization 1 → N Selling Entities`; identities do not collapse;
- canonical MPC semantics are defined before ERP/provider mapping;
- Selling Entity, Inventory Source/Scope, Fulfillment Node/Scope, Cost Observation/Basis, Business Order Intent and Invoicing Intent are accepted semantics;
- availability is derived from eligible authoritative inventory + rules + policy; routine known policy-valid sync is automatic; unknown is not zero;
- MPC-owned allocation policy may intentionally expose less stock (e.g. 70%) with deterministic/explainable future scope/override semantics;
- Business Order Intent and Invoicing Intent remain separate from native ERP TOP/document operations; ambiguous writes are not blindly retried;
- the Economic Evidence Chain is `Simulation → Order Economics → Marketplace Settlement → Cash Receipt evidence where available`;
- simulation/order variance is classified before creating a simulator-calibration defect; payout/cash matching is not assumed 1:1;
- provider capability/authority depends on effective operating context; claimed MPC-controlled paths must close provider-required prerequisites/data/artifacts/readiness;
- native provider artifacts/states remain provider-owned semantics rather than universal MPC types;
- ANYMARKET, Magis5 and similar hubs are benchmark/competitive evidence, not Product 1.0 runtime dependencies; direct MPC-owned provider boundaries are the current target direction;
- materially time-bound obligations are explicit: provider/external obligation and MPC-owned organization target remain separate, relative targets require explicit anchors, and time participates in portfolio attention/breach semantics;
- an internal operating target may be stricter than a provider deadline but cannot relax or overwrite that external obligation;
- current code/docs remain evidence, not target authority;
- no D1+ target architecture may be invented yet;
- the exact next work is **D0.7 Product completeness review — operational work ownership / assignment / escalation**.
