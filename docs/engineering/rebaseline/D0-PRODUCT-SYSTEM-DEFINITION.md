# D0 — Product / System Definition

> **Status:** WORKING — D0 OPEN, NOT YET ACCEPTED AS A WHOLE  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Authority:** operator-approved D0 decisions recorded here; later D-stages own technical realization  
> **Last updated:** 2026-08-14

## 1. Purpose

D0 defines what Marketplace Central is, who it serves, which problems and outcomes belong to Product 1.0, and where the product boundary ends.

D0 does **not** decide target contexts/modules, identity model, database schema, API shape, frontend topology, runtime/jobs, or provider-specific transport. Those belong to D1–D7.

Current code, schemas, APIs and historical ADRs are evidence only unless the active rebaseline marks a constraint binding.

## 2. Accepted D0 decisions

### D0.1 — Product mission: Marketplace Operations Control Plane

**Accepted by operator.**

Marketplace Central is the internal **Marketplace Operations Control Plane** of the company.

It combines internal business facts with marketplace observations so operators can understand the real operational state, detect divergences, make grounded decisions, execute controlled actions in participating systems, and verify/reconcile the result afterward.

Its fundamental operating loop is:

```text
observe → understand → reconcile → decide/policy → execute → verify → audit/reconcile
```

Marketplace Central is not merely an intelligence dashboard and is not a replacement ERP or marketplace platform.

It can coordinate and execute operations on both sides of the boundary. Examples include creating/changing marketplace listings and prices, and causing marketplace-originated orders to be created and invoiced in Sankhya through an accepted integration path.

External systems remain authorities for facts that inherently belong to them; the exact data/identity/write ownership matrix is deferred to D2/D4. MPC owns the cross-system operational control semantics: intent, policy application, workflow, correlation, controlled execution, operational state, divergence, audit and reconciliation.

### D0.2 — Product 1.0 scope: Operations + Commercial Intelligence (A+)

**Accepted by operator.**

Product 1.0 must close a meaningful operational loop **and** remove major manual commercial-analysis work. Competitive intelligence, pricing and profitability are therefore core Product 1.0 capabilities, not optional future dashboard features.

Product 1.0 must support the combined loop:

```text
internal product
  → channel readiness / linkage
  → marketplace listing operations
  → competitive market observation
  → pricing & expected profitability
  → controlled decision / policy
  → marketplace action
  → sale / marketplace order
  → ERP operation / invoicing
  → fulfillment / dispatch
  → delivery/post-sale lifecycle when applicable
  → realized economics
  → reconciliation / exception handling
```

The provider-independent product requirement is to observe the relevant market for a product and compare MPC's offer against semantically comparable external offers. Mercado Livre is the first provider used to prove that capability; provider-specific mechanisms belong to D4.

### D0.3 — Product 1.0 capability boundary

**Accepted by operator.**

Product 1.0 includes these product capabilities:

1. **Product & Channel Readiness** — understand which internal products can operate in a marketplace, their channel linkage/readiness, and missing/conflicting conditions.
2. **Marketplace Listing Operations** — create, inspect and control listings and material operational state such as price/availability, then verify whether requested changes converged.
3. **Competitive Intelligence** — observe comparable market offers/prices, expose competitive position and meaningful changes, and represent insufficient comparison evidence honestly.
4. **Pricing & Profitability Intelligence** — combine relevant internal economics and market observations to calculate price scenarios, expected margin/profitability and decision-relevant trade-offs.
5. **Decision & Policy Control** — translate observations/recommendations into permitted, approval-required or prohibited actions according to governing company rules/policies.
6. **Order-to-ERP Operations** — receive/understand marketplace orders and coordinate the corresponding business operations in Sankhya, including order creation and invoicing where authorized.
7. **Marketplace Fulfillment / Dispatch** — progress marketplace orders through physical separation, conference, invoicing trigger, packing and dispatch handoff without becoming a company-wide WMS.
8. **Shipment / Delivery Observation & Exceptions** — continue observing the marketplace shipment after dispatch handoff until a relevant terminal outcome, surfacing delays, delivery failures, returns or other delivery exceptions without becoming a company-wide TMS or carrier platform.
9. **Essential Post-Sale Operations** — control the operational response to marketplace cancellations, returns and refunds when they affect an MPC-controlled sale, coordinating consequences across the marketplace/ERP/fulfillment/economic loop without becoming a general CRM/SAC platform.
10. **Reconciliation & Exception Operations** — identify cases where an expected cross-system result cannot be proven or systems diverge, and make them actionable instead of silently assuming success.
11. **Realized Profitability** — determine realized sale economics, including material delivery/post-sale reversals or adjustments, and compare expected versus realized results so material variance can be explained.

Automation/human approval is a cross-cutting control-plane property, not a separate business capability.

### D0.3a — Action authority model

**Accepted by operator.**

Product 1.0 supports:

- **human-controlled execution by default:** recommendation/analysis → operator approval → MPC executes → verifies/reconciles;
- **policy-driven automatic execution when explicitly authorized:** an accepted policy may approve a bounded action automatically;
- **human review for exceptions, uncertainty, low confidence or policy violations.**

A fully autonomous repricing system is not required as a launch gate. The product must, however, have the semantics necessary to support explicit bounded automation without bypassing policy, audit or reconciliation.

### D0.3b — Policy/rule provenance is explicit

**Accepted by operator.**

Marketplace Central must not assume that every business rule or commercial policy is authored inside MPC.

A governing rule used by MPC may be:

- **MPC-owned** — intentionally defined and governed inside Marketplace Central;
- **externally governed** — sourced from Sankhya/another ERP or another authoritative system and consumed by MPC;
- **derived** — mechanically computed from authoritative facts/rules without becoming an independent source of truth.

The product-level invariant is that MPC must know enough provenance to distinguish these classes and must not silently turn an externally governed rule into an editable MPC-owned copy.

A Commercial / Marketplace Manager may define or change **MPC-owned** commercial policies within their authority, but cannot override an externally authoritative ERP/system rule merely because MPC consumes it. If business requirements need that external rule changed, the change must occur through the system/authority that owns it or through an explicitly accepted cross-system workflow.

The exact ownership matrix, synchronization mechanism, conflict semantics, freshness rules and provider contracts are deferred to D2/D4.

## 3. Product 1.0 non-goals currently safe to defer

The following are not required to define Product 1.0 as complete unless later D0 evidence proves otherwise:

- paid ads/media management;
- automated buyer Q&A or buyer chat;
- general CRM/SAC and broad customer-service automation;
- broad reputation optimization;
- company-wide demand forecasting or purchasing;
- company-wide logistics/WMS/TMS replacement;
- marketplaces beyond what is needed to prove the first Mercado Livre operating loop;
- unrestricted AI/autonomous commercial decisions without explicit policy;
- a generic integration framework for many speculative providers.

These may become future MPC capabilities; deferral is not a permanent exclusion. Essential cancellation/return/refund operations tied to an MPC-controlled marketplace sale and shipment/delivery observation for that sale are **not** deferred by this list; D0.7 includes them in Product 1.0.

## 4. Operator-provided Sankhya evidence / constraint

This is **D0/D4 evidence**, not a provider-transport decision made in D0.

The operator reports an already-proven Sankhya application integration in another app using application credentials (`client id`, `client secret`, `X-Token`; no secret values are recorded here).

Operational writes such as order creation and invoicing can be performed through Sankhya's application API. The operator prefers this path for writes rather than direct Oracle database writes because it is a safer system-owned operation boundary.

The environment also has DB Explorer access for database inspection.

Consequences for later design:

- D0 may require MPC to create/invoice orders in Sankhya without assuming direct database writes;
- D4 must evaluate and ratify the exact Sankhya read/write capability contracts and transport boundaries;
- D2/D4 must account for business rules/policies whose authority remains in Sankhya/another ERP rather than duplicating them as MPC-owned configuration;
- existing binding Oracle-read constraints are not silently reopened here; the target read path and any need for direct Oracle access are adjudicated only by the stage that owns that decision.

## 5. Stable constraints carried into D0

D0 remains constrained by current accepted repository authority, including:

- Mercado Livre first;
- Sankhya/Oracle external to MPC;
- Go backend as canonical business execution;
- React frontend is not a second business authority;
- PostgreSQL stores MPC-owned canonical state;
- unknown/absent facts do not become plausible known defaults;
- external writes require explicit authority/policy, duplicate protection, auditability and reconciliation;
- ambiguous external-write outcomes are not blindly retried;
- provider PII is minimized;
- provider-specific protocol details remain behind provider boundaries.

## 6. D0.4 — Actors / operational users

**Accepted by operator.**

D0 distinguishes human actors by the operational responsibility they carry. "Operator" is an umbrella description, not a single undifferentiated persona.

The accepted actor classes are:

1. **Marketplace Operations Operator** — day-to-day channel operations around products, listings, prices, marketplace state, orders, divergences and operational exceptions within permitted authority.
2. **Fulfillment / Dispatch Operator** — physical fulfillment execution for marketplace orders: work queue, separation/conference, invoicing trigger when valid, packing, dispatch handoff and exception reporting.
3. **Commercial / Marketplace Manager** — commercial authority for MPC-owned marketplace policies, pricing/margin boundaries, higher-impact approvals, exception decisions and bounded automation policies.
4. **Owner / Administrator / Policy Approver** — system/organization governance: integrations, users, exceptional/high-impact authority and controls that must remain above ordinary commercial operation.

This actor model does **not** define JWT roles, permission tables or technical authorization implementation. Those are later-stage concerns.

### D0.4a — Fulfillment / Dispatch Operator authority

**Accepted by operator.**

The fulfillment workflow follows the business sequence:

```text
marketplace order / ERP order readiness
  → fulfillment work queue
  → physical separation
  → physical verification / conference
  → if valid: operator triggers invoicing through MPC
  → MPC causes the corresponding Sankhya invoicing operation through the later-accepted integration boundary
  → invoicing result is verified
  → packing
  → dispatch / carrier handoff
  → completion or exception
```

The target property is that **an order is not intentionally invoiced through the normal fulfillment path before the operator has physically confirmed that the correct items are available and separated**.

If the operator finds missing stock, wrong item, damaged material, quantity divergence or another physical inconsistency, normal invoicing is blocked and the order becomes an operational exception requiring resolution.

D0 does not decide permission representation, workflow persistence or Sankhya transport.

### D0.4b — Marketplace Operations Operator authority

**Accepted by operator.**

The Marketplace Operations Operator owns routine marketplace control within governing policy. At product/workflow level this actor may:

- prepare/correct product-channel readiness and linkage where evidence is sufficient;
- create and publish a new listing when the product is fully ready and all governing conditions are satisfied;
- inspect/edit listings and permitted operational state;
- analyze competitive position, price scenarios and expected profitability;
- execute price changes and other bounded marketplace actions inside policy;
- pause/reactivate or correct marketplace operational state when policy-compliant;
- investigate and progress ordinary divergences/exceptions inside operational authority.

The actor does not own the policies/rules themselves and cannot redefine a boundary merely to make an action permissible.

```text
fully ready + inside governing policy
  → Marketplace Operations Operator may decide/execute

outside policy / insufficient evidence / higher-impact exception
  → Commercial / Marketplace Manager review or other explicit escalation
```

Initial listing creation is not approval-gated merely because it is the first publication.

### D0.4c — Commercial / Marketplace Manager authority

**Accepted by operator, with policy-provenance constraint.**

The Commercial / Marketplace Manager is the ordinary commercial authority for the marketplace operation, but only over policies and decisions that are actually within the manager's authority.

This actor may, at product/workflow level:

- define/change **MPC-owned** margin floors, price boundaries, approval thresholds and other marketplace commercial policies;
- approve/reject commercial exceptions escalated by Marketplace Operations Operators;
- authorize higher-impact price/listing actions when policy requires review;
- define which bounded action classes may execute automatically and the commercial constraints around that automation;
- review commercial/competitive/profitability intelligence at managerial scope and convert accepted strategy into MPC-owned operating policy;
- suspend or narrow an MPC-owned automation/policy when commercial risk requires intervention.

This actor is **not automatically the authority over rules sourced from Sankhya/another ERP or another governing system**. Externally governed rules remain externally governed unless a later accepted ownership decision explicitly transfers authority.

The manager also does not administer integration credentials, security controls, users, structural tenant/organization settings or bypass audit/reconciliation controls merely because they affect commercial operations.

Ordinary commercial policy should not require Owner approval by default. Escalation to Owner/Admin is reserved for organization/system governance or explicitly high-impact exceptional authority.

### D0.4d — Owner / Administrator / Policy Approver authority

**Accepted by operator.**

The Owner / Administrator / Policy Approver governs concerns that sit above routine marketplace commercial operation, including at product/workflow level:

- integration/connection administration and organization-level system configuration;
- user/access governance at an organizational level (technical realization deferred);
- exceptional/high-impact authority not delegated to the Commercial / Marketplace Manager;
- governance boundaries around who may define policies or authorize automation;
- emergency suspension/containment of risky external actions or automation when organizational control requires it;
- resolution/escalation when a policy or integration authority conflict cannot be settled inside normal marketplace operations.

The Owner/Admin is **not** intended to approve routine commercial changes merely to create a longer approval chain. Commercial decisions inside delegated authority terminate at the Commercial / Marketplace Manager.

No actor may use its authority to disable mandatory audit/reconciliation/safety invariants or silently convert externally governed business rules into local editable copies.

## 7. D0.5 — System boundary / authority classes

**Accepted by operator.**

Marketplace Central uses three product-level authority classes:

- **OWN** — MPC is the business authority for the concern/state;
- **ORCHESTRATE** — another system remains authoritative for the underlying business fact/process, while MPC owns the cross-system operational intent/control/workflow around its participation in marketplace operations;
- **OBSERVE / DERIVE** — MPC consumes authoritative external facts and may derive decision-support information without becoming source of truth for the underlying fact.

The governing principle is:

> **MPC owns the marketplace operating model. External systems own the facts and processes that inherently belong to them. MPC orchestrates their participation in the marketplace operating loop.**

This classification is conceptual. D0 does not decide tables, context ownership, sync topology, ports, APIs or persistence mechanisms.

### D0.5a — Product 1.0 authority map

| Concern | Product-level MPC authority |
|---|---|
| Internal master product facts, ERP base cost and ERP-governed fiscal facts | **OBSERVE / DERIVE** |
| Business rules/policies whose authority remains in Sankhya/another ERP/system | **OBSERVE / CONSUME** |
| Product ↔ marketplace linkage/evidence maintained by MPC | **OWN** |
| Marketplace readiness and readiness assessment | **OWN / DERIVE** |
| Cross-system operational intent, workflow, correlation, divergence and exception state | **OWN** |
| Competitive intelligence and comparable-market interpretation | **OWN / DERIVE** |
| Pricing scenarios and expected profitability | **OWN / DERIVE** |
| MPC-owned marketplace commercial policies | **OWN** |
| Actual listing/channel state that inherently exists in the marketplace provider | provider is authoritative; MPC **ORCHESTRATES** and observes |
| MPC intent to create/change/pause/reactivate a marketplace listing or price | **OWN / ORCHESTRATE** |
| Marketplace-originated order facts | marketplace/provider is authoritative for originating channel facts; MPC **ORCHESTRATES** |
| ERP order/invoice/accounting facts produced in Sankhya | Sankhya is authoritative; MPC **ORCHESTRATES** the marketplace workflow that causes/observes them |
| Marketplace-order fulfillment workflow: queue, conference, invoicing trigger, packing/dispatch progression, exceptions | **OWN**, while underlying ERP/carrier facts retain their external authority |
| Shipment/delivery state after dispatch handoff | marketplace/carrier/provider remains authoritative for native logistics facts; MPC **OBSERVES / ORCHESTRATES** delivery exceptions and sale-lifecycle closure |
| Marketplace cancellation/return/refund channel facts | marketplace/provider remains authoritative; MPC **ORCHESTRATES** the cross-system response and observes/reconciles consequences |
| ERP reversal/credit/fiscal/accounting facts caused by post-sale operations | Sankhya/ERP remains authoritative; MPC **ORCHESTRATES** the marketplace post-sale workflow that requires/observes them |
| Realized profitability interpretation | **OWN / DERIVE** from authoritative realized facts, including material delivery/post-sale adjustments |
| Audit/reconciliation records for MPC-controlled operations | **OWN** |

### D0.5b — Boundary invariants

Product 1.0 must preserve these invariants regardless of later implementation shape:

1. Observing an external fact does not transfer ownership of that fact to MPC.
2. Orchestrating an external process does not make MPC the source of truth for the external system's resulting native record.
3. An MPC-derived commercial conclusion must preserve the provenance/freshness needed to understand which authoritative facts/rules produced it.
4. Externally governed rules are not silently copied into mutable MPC-owned policy.
5. MPC-owned workflow state is not replaced by guessing from one provider response; ambiguous or divergent cross-system outcomes remain explicit until reconciled.
6. Provider/ERP-specific mechanisms remain implementation concerns for later stages; the Product 1.0 boundary is stated in business terms.

## 8. D0.6 — Product 1.0 completion / user-observable outcomes

**Accepted by operator.**

Product 1.0 is complete only when Marketplace Central is demonstrably usable as the normal operational control plane, not merely when individual capabilities exist in isolation.

The acceptance bar is user-observable:

1. **Attention is portfolio-driven, not manual-search driven.** Operators can see what is healthy, changed, divergent, blocked, approval-required or otherwise actionable without manually inspecting products and external systems one by one.
2. **An eligible internal product can reach a verified marketplace state through MPC.** The normal path covers readiness/linkage, commercial analysis, creation/publication and subsequent observation of the real channel state.
3. **Competitive/pricing intelligence replaces major manual analysis.** MPC can expose grounded comparable-market position, relevant internal economics, price scenarios, expected profitability and insufficient-evidence cases at portfolio and individual-product level.
4. **Decision closes into controlled action.** Authorized human or bounded policy-driven decisions can become external actions with policy enforcement, auditability, verification and reconciliation rather than ending as recommendations that require ordinary manual execution elsewhere.
5. **A marketplace sale can traverse the normal operating loop through MPC.** The normal path covers order recognition, corresponding ERP operation, fulfillment queue, physical separation/conference, invoicing trigger and verification, packing and dispatch handoff without hidden manual system hopping as a required normal step.
6. **Delivery remains visible through terminal outcome.** After dispatch handoff, MPC continues observing the relevant shipment/delivery lifecycle until delivered, returned, cancelled or another equivalent terminal state; delays/failures/returns become explicit operational work rather than forcing routine native-system monitoring.
7. **Essential post-sale changes remain inside the controlled lifecycle.** Cancellation, return or refund events that materially affect an MPC-controlled marketplace sale can be understood and progressed through the necessary cross-system response/reconciliation without dropping the normal operation back into manual system hopping.
8. **Failures become explicit work.** Missing evidence, ambiguous external results, integration failures and physical/order/delivery/post-sale divergences are surfaced with what is known/unknown and what requires action instead of silently becoming plausible success/defaults.
9. **The economic loop closes.** MPC can compare expected versus realized profitability using attributable authoritative facts, including material delivery/cancellation/return/refund effects, so a reversed or materially adjusted sale does not remain represented using stale pre-adjustment economics.
10. **Organizational governance is operable without code edits.** Actors can exercise the legitimate MPC-owned authorities assigned in D0.4, while externally governed rules remain governed externally and mandatory safety/audit invariants are not configurable away.

The Product 1.0 completion statement is:

> **A company can take its internal products, determine marketplace readiness, publish and operate offers, monitor market position and profitability at portfolio scale, make and execute decisions under policy, receive sales, progress them through ERP and physical fulfillment, follow delivery to a terminal outcome, handle essential cancellation/return/refund consequences, surface/reconcile exceptions, and understand the realized economic result — using Marketplace Central as the normal marketplace operations control plane.**

### D0.6a — Normal-path rule

The normal operational path must be executable through MPC for responsibilities that Product 1.0 claims to control.

Direct use of Mercado Livre, Sankhya or another participating external system remains legitimate when:

- the responsibility inherently belongs to that external system and is intentionally outside MPC scope;
- investigation/support requires native-system inspection;
- an exceptional recovery path explicitly requires it.

Direct external-system hopping must **not** be a hidden required step in an otherwise claimed MPC normal workflow. A Product 1.0 flow with routine manual gaps between MPC, marketplace and ERP is incomplete even if each individual integration exists.

Detailed end-to-end proof scenarios remain a D8 responsibility; D0.6 defines the outcomes that D8 must eventually prove.

## 9. D0.7 — Product completeness / contradiction review

**OPEN — exact next D0 work.**

Before D0 can be accepted as a whole, review the accepted Product 1.0 definition adversarially for missing business lifecycle responsibilities, contradictions and accidental scope gaps.

### D0.7a — Essential post-sale lifecycle

**Accepted by operator.**

Product 1.0 includes the essential operational lifecycle for **cancellations, returns and refunds** when they affect an MPC-controlled marketplace sale.

The product-level responsibility is to:

- understand the authoritative marketplace post-sale event/state;
- stop or redirect still-avoidable downstream work when appropriate;
- coordinate the required cross-system operational response with ERP/fulfillment/marketplace participants;
- surface unresolved or contradictory outcomes as explicit exceptions;
- reconcile the resulting operational and economic consequences;
- ensure realized profitability reflects material reversals/adjustments rather than freezing the economics at the pre-return/pre-refund state.

This does **not** make Product 1.0 a general CRM/SAC, buyer-messaging system, complaint/mediation automation suite, reputation platform or company-wide reverse-logistics platform. Those remain deferrable unless later evidence reopens them.

Marketplace/provider remains authoritative for native marketplace cancellation/return/refund facts; ERP remains authoritative for its native fiscal/accounting/reversal records; MPC owns/orchestrates the cross-system control workflow and its exception/reconciliation state.

Exact statuses, return logistics, fiscal reversal mechanics, refund contracts and provider-specific behavior belong to later D-stages.

### D0.7b — Shipment / delivery lifecycle after dispatch handoff

**Accepted by operator.**

Product 1.0 continues observing the shipment/delivery lifecycle of an MPC-controlled marketplace sale after physical dispatch handoff until the order reaches a relevant terminal outcome such as delivered, returned, cancelled or an equivalent terminal provider state.

The product-level responsibility is to:

- observe the authoritative shipment/delivery state supplied by the marketplace/carrier/provider;
- surface delays, failed delivery attempts, returns-to-sender and other material delivery exceptions as operational work;
- correlate delivery exceptions with the MPC-controlled sale workflow;
- route a returned/cancelled/post-sale consequence into the appropriate post-sale/reconciliation path;
- avoid forcing routine operator monitoring in the native marketplace/carrier system merely to know whether an MPC-controlled sale completed.

This does **not** make Product 1.0 a carrier, route planner, fleet system, company-wide transport manager or TMS. Native carrier/provider logistics facts remain externally authoritative; MPC observes/orchestrates only what is needed to close the marketplace sale lifecycle and its exceptions.

Exact tracking mechanisms, carrier/provider contracts and status vocabularies belong to later D-stages.

### Next D0.7 question

The review continues. The next material boundary is **stock / marketplace availability control**. Product 1.0 already claims listing availability operations and observes ERP facts, but D0 has not yet stated whether maintaining marketplace availability coherently with the governing internal stock/availability rules is a normal MPC responsibility or merely a manual operator action.

This must be decided at product level before D0 closes. Exact stock authority, reservation semantics, buffers, source reads, event/sync strategy and provider update mechanisms belong to D2–D4/D7.

Other findings discovered during D0.7 are classified as MUST DECIDE NOW, SHOULD DECIDE NOW or CAN DEFER SAFELY. D0 closes only when no material Product 1.0 semantic is being left for implementation to invent.

## 10. Resume contract for a fresh session

A fresh session must read `AGENTS.md`, `docs/engineering/rebaseline/README.md`, the canonical engineering method, `ARCHITECTURE.md`, the ADR registry, and then this D0 artifact.

It should conclude:

- D0 is open and not yet accepted as a whole;
- D0.1–D0.3b are operator-approved product decisions;
- D0.4 is accepted with four actors and explicit responsibility/authority boundaries;
- D0.5 is accepted: MPC uses **OWN / ORCHESTRATE / OBSERVE-DERIVE** and owns the marketplace operating model without taking ownership of external facts merely because it consumes or causes them;
- D0.6 is accepted: Product 1.0 has a user-observable completion bar and the normal path for MPC-owned/orchestrated responsibilities must be executable through MPC rather than requiring routine manual system hopping;
- D0.7a is accepted: essential cancellation/return/refund operations remain inside the MPC-controlled sale lifecycle without expanding Product 1.0 into general CRM/SAC;
- D0.7b is accepted: shipment/delivery state remains visible through a terminal outcome and material delivery exceptions become MPC operational work without turning MPC into a TMS;
- Sankhya API availability for writes is evidence to carry into D4, not a D0 transport decision;
- business rules/policies may be MPC-owned, externally governed or derived, and later D2/D4 must preserve that provenance;
- no D1+ target architecture may be invented yet;
- the exact next work is **D0.7 Product completeness review — stock / marketplace availability control**.
