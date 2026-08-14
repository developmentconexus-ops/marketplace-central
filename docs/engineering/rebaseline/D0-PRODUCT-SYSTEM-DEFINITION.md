# D0 — Product / System Definition

> **Status:** WORKING — D0 OPEN, NOT YET ACCEPTED AS A WHOLE  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Authority:** operator-approved D0 decisions recorded here; later D-stages own technical realization  
> **Last updated:** 2026-08-14

## 1. Purpose

D0 defines what Marketplace Central is, who it serves, which problems and outcomes belong to Product 1.0, and where the product boundary ends.

D0 does **not** decide target contexts/modules, identity model, database schema, API shape, frontend topology, runtime/jobs, or provider-specific transport. Those belong to D1–D7.

Current code, schemas, APIs and historical ADRs are evidence only unless the active rebaseline marks a constraint binding.

## 2. D0.1 — Product mission

**Accepted by operator.**

Marketplace Central is the internal **Marketplace Operations Control Plane** of the company.

It combines internal business facts with marketplace observations so operators can understand the real operational state, detect divergences, make grounded decisions, execute controlled actions in participating systems, and verify/reconcile the result afterward.

Its fundamental loop is:

```text
observe → understand → reconcile → decide/policy → execute → verify → audit/reconcile
```

Marketplace Central is not merely an intelligence dashboard and is not a replacement ERP or marketplace platform. It can coordinate and execute operations on both sides of the boundary, including marketplace listing/price actions and marketplace-originated order/invoicing workflows in Sankhya through an accepted integration path.

External systems remain authorities for facts that inherently belong to them. MPC owns the cross-system operational control semantics: intent, policy application, workflow, correlation, controlled execution, operational state, divergence, audit and reconciliation.

## 3. D0.2 / D0.3 — Product 1.0 scope and capability boundary

**Accepted by operator.**

Product 1.0 is **Marketplace Operations + Commercial Intelligence (A+)**. It must close the operating loop and remove major manual commercial-analysis work.

The Product 1.0 loop is:

```text
internal product
  → channel readiness / linkage
  → listing + marketplace availability control
  → competitive market observation
  → pricing & expected profitability
  → controlled decision / policy
  → marketplace action
  → sale / marketplace order
  → ERP operation / invoicing
  → fulfillment / dispatch
  → shipment / delivery lifecycle
  → essential cancellation / return / refund lifecycle when applicable
  → realized economics
  → reconciliation / exception handling
```

The accepted Product 1.0 capabilities are:

1. **Product & Channel Readiness** — determine which internal products can operate in a marketplace, their linkage/readiness, and missing/conflicting conditions.
2. **Marketplace Listing Operations** — create, inspect and control listings and material channel state, then verify whether requested changes converged.
3. **Marketplace Availability Control** — keep published marketplace availability coherently aligned with sellable availability derived from authoritative stock facts/rules, automatically on the normal path and with explicit exceptions when evidence or convergence is uncertain.
4. **Competitive Intelligence** — observe comparable market offers/prices, expose competitive position and meaningful changes, and represent insufficient comparison evidence honestly.
5. **Pricing & Profitability Intelligence** — combine relevant internal economics and market observations to calculate price scenarios, expected margin/profitability and decision-relevant trade-offs.
6. **Decision & Policy Control** — translate observations/recommendations into permitted, approval-required or prohibited actions according to governing company rules/policies.
7. **Order-to-ERP Operations** — receive/understand marketplace orders and coordinate corresponding operations in Sankhya, including order creation and invoicing where authorized.
8. **Marketplace Fulfillment / Dispatch** — progress marketplace orders through physical separation, conference, invoicing trigger, packing and dispatch handoff without becoming a company-wide WMS.
9. **Shipment / Delivery Observation & Exceptions** — continue observing shipment after dispatch handoff until a relevant terminal outcome, surfacing delays, delivery failures, returns or other material exceptions without becoming a TMS/carrier platform.
10. **Essential Post-Sale Operations** — control the operational response to marketplace cancellations, returns and refunds when they affect an MPC-controlled sale, coordinating consequences across marketplace/ERP/fulfillment/economic workflows without becoming general CRM/SAC.
11. **Reconciliation & Exception Operations** — identify cases where an expected cross-system result cannot be proven or systems diverge, making them explicit work instead of silently assuming success.
12. **Realized Profitability** — determine realized sale economics, including material delivery/post-sale reversals or adjustments, and compare expected versus realized results.

The provider-independent requirement is to express these capabilities in business terms. Mercado Livre is the first provider used to prove them; provider-specific mechanics belong to D4.

### D0.3a — Action authority model

**Accepted by operator.**

Product 1.0 supports:

- **human-controlled execution by default:** recommendation/analysis → operator approval → MPC executes → verifies/reconciles;
- **policy-driven automatic execution when explicitly authorized:** an accepted policy may approve a bounded action automatically;
- **human review for exceptions, uncertainty, low confidence or policy violations.**

A fully autonomous repricing engine is not required as a launch gate. The product must nevertheless support explicit bounded automation without bypassing policy, audit or reconciliation.

Marketplace availability synchronization is different from discretionary commercial repricing: routine stock/availability changes are expected to execute automatically when the governing facts/rules are sufficiently known and the action is valid. Human attention is concentrated on uncertainty, policy conflicts and failed/non-convergent updates.

### D0.3b — Policy/rule provenance is explicit

**Accepted by operator.**

Marketplace Central must not assume every business rule or commercial policy is authored inside MPC.

A governing rule used by MPC may be:

- **MPC-owned** — intentionally defined and governed inside Marketplace Central;
- **externally governed** — sourced from Sankhya/another ERP or another authoritative system and consumed by MPC;
- **derived** — mechanically computed from authoritative facts/rules without becoming an independent source of truth.

MPC must preserve enough provenance to distinguish these classes and must not silently turn an externally governed rule into an editable MPC-owned copy.

The exact ownership matrix, synchronization mechanism, conflict semantics, freshness rules and provider contracts are deferred to D2/D4.

## 4. Product 1.0 non-goals currently safe to defer

The following are not required to call Product 1.0 complete unless later D0 evidence reopens them:

- paid ads/media management;
- automated buyer Q&A or buyer chat;
- general CRM/SAC and broad customer-service automation;
- broad reputation optimization;
- company-wide demand forecasting or purchasing;
- company-wide logistics/WMS/TMS replacement;
- marketplaces beyond what is needed to prove the first Mercado Livre operating loop;
- unrestricted AI/autonomous commercial decisions without explicit policy;
- a generic integration framework for many speculative providers.

Essential cancellation/return/refund operations and shipment/delivery observation tied to an MPC-controlled marketplace sale are **not** deferred.

## 5. Operator-provided Sankhya evidence / constraint

This is **D0/D4 evidence**, not a provider-transport decision made in D0.

The operator reports an already-proven Sankhya application integration in another app using application credentials (`client id`, `client secret`, `X-Token`; no secret values are recorded here).

Operational writes such as order creation and invoicing can be performed through Sankhya's application API. The operator prefers this path for writes rather than direct Oracle database writes because it is a safer system-owned operation boundary. The environment also has DB Explorer access for database inspection.

Consequences for later design:

- D0 may require MPC to create/invoice orders in Sankhya without assuming direct database writes;
- D4 must evaluate and ratify exact Sankhya read/write capability contracts and transport boundaries;
- D2/D4 must account for rules/policies whose authority remains in Sankhya/another ERP rather than duplicating them as MPC-owned configuration;
- existing binding Oracle-read constraints are not silently reopened here.

## 6. Stable constraints carried into D0

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

## 7. D0.4 — Actors / operational users

**Accepted by operator.**

“Operator” is an umbrella description, not one undifferentiated persona. Product 1.0 has four accepted actor classes.

### Marketplace Operations Operator

Owns routine day-to-day marketplace control inside governing policy. This actor may:

- prepare/correct product-channel readiness and linkage where evidence is sufficient;
- create/publish a new listing when the product is fully ready and governing conditions are satisfied;
- inspect/edit listings and permitted operational state;
- analyze competitive position, price scenarios and expected profitability;
- execute bounded price/listing actions inside policy;
- investigate and progress ordinary divergences/exceptions inside operational authority;
- inspect and resolve marketplace-availability exceptions, while routine policy-valid availability synchronization remains automatic rather than requiring approval per stock change.

The actor cannot redefine policies/rules merely to make an action permissible.

```text
fully ready + inside governing policy
  → Marketplace Operations Operator may decide/execute

outside policy / insufficient evidence / higher-impact exception
  → Commercial / Marketplace Manager review or other explicit escalation
```

### Fulfillment / Dispatch Operator

Owns physical fulfillment execution for marketplace orders: work queue, separation/conference, invoicing trigger when valid, packing, dispatch handoff and exception reporting.

Accepted normal sequence:

```text
marketplace order / ERP order readiness
  → fulfillment queue
  → physical separation
  → physical conference
  → if valid: operator triggers invoicing through MPC
  → MPC causes Sankhya invoicing through the later-accepted integration boundary
  → invoicing result is verified
  → packing
  → dispatch handoff
  → completion or exception
```

**Invariant:** an order is not intentionally invoiced through the normal fulfillment path before the operator has physically confirmed that the correct items are available and separated. Missing stock, wrong item, damaged material, quantity divergence or another physical inconsistency blocks normal invoicing and becomes an operational exception.

### Commercial / Marketplace Manager

Is the ordinary commercial authority for the marketplace operation, but only over policies/decisions actually within that actor’s authority. This actor may:

- define/change **MPC-owned** margin floors, price boundaries, approval thresholds and marketplace commercial policies;
- approve/reject commercial exceptions;
- authorize higher-impact actions when policy requires review;
- define bounded automation classes and their commercial constraints;
- review competitive/profitability intelligence and convert accepted strategy into MPC-owned operating policy;
- suspend/narrow MPC-owned automation when commercial risk requires intervention.

This actor is **not automatically the authority over externally governed rules** and does not administer integration credentials, security controls or organizational access governance.

### Owner / Administrator / Policy Approver

Governs concerns above routine marketplace commercial operation, including:

- integration/connection administration and organization-level configuration;
- user/access governance at organizational level;
- exceptional/high-impact authority not delegated to the Commercial / Marketplace Manager;
- governance boundaries around who may define policies or authorize automation;
- emergency suspension/containment of risky external actions/automation;
- resolution/escalation when a policy or integration authority conflict cannot be settled in normal operations.

The Owner/Admin is not intended to approve routine commercial changes merely to create a longer approval chain.

No actor may disable mandatory audit/reconciliation/safety invariants or silently convert externally governed rules into local editable copies.

## 8. D0.5 — System boundary / authority classes

**Accepted by operator.**

Marketplace Central uses three product-level authority classes:

- **OWN** — MPC is the business authority for the concern/state;
- **ORCHESTRATE** — another system remains authoritative for the underlying business fact/process, while MPC owns the cross-system operational intent/control/workflow around it;
- **OBSERVE / DERIVE** — MPC consumes authoritative external facts and may derive decision-support information without becoming source of truth for the underlying fact.

Governing principle:

> **MPC owns the marketplace operating model. External systems own the facts and processes that inherently belong to them. MPC orchestrates their participation in the marketplace operating loop.**

### D0.5a — Product 1.0 authority map

| Concern | Product-level MPC authority |
|---|---|
| Internal master product facts, ERP base cost/fiscal facts | **OBSERVE / DERIVE** |
| ERP physical/on-hand/reserved/other stock facts | **OBSERVE / CONSUME** |
| Business rules/policies whose authority remains in Sankhya/another system | **OBSERVE / CONSUME** |
| Product ↔ marketplace linkage/evidence maintained by MPC | **OWN** |
| Marketplace readiness/readiness assessment | **OWN / DERIVE** |
| Sellable marketplace availability derived from authoritative stock/rules | **OWN / DERIVE** as an MPC operating conclusion; underlying stock facts/rules retain their source authority |
| Intent to publish/update marketplace quantity/availability | **OWN / ORCHESTRATE** |
| Actual marketplace quantity/availability state | provider authoritative; MPC **OBSERVES / ORCHESTRATES** convergence |
| Cross-system intent, workflow, correlation, divergence and exception state | **OWN** |
| Competitive intelligence / comparable-market interpretation | **OWN / DERIVE** |
| Pricing scenarios / expected profitability | **OWN / DERIVE** |
| MPC-owned marketplace commercial policies | **OWN** |
| Actual listing/channel state inherent in marketplace | provider authoritative; MPC **OBSERVES / ORCHESTRATES** |
| Marketplace listing/price mutation intent | **OWN / ORCHESTRATE** |
| Marketplace-originated order facts | marketplace/provider authoritative; MPC **ORCHESTRATES** |
| ERP order/invoice/accounting facts | Sankhya authoritative; MPC **ORCHESTRATES** marketplace workflow around them |
| Marketplace-order fulfillment workflow | **OWN**, while underlying ERP/carrier facts retain external authority |
| Shipment/delivery state after dispatch | marketplace/carrier/provider authoritative; MPC **OBSERVES / ORCHESTRATES** exceptions/lifecycle closure |
| Marketplace cancellation/return/refund facts | marketplace/provider authoritative; MPC **ORCHESTRATES** cross-system response |
| ERP reversal/credit/fiscal/accounting facts | ERP authoritative; MPC **ORCHESTRATES** post-sale workflow |
| Realized profitability interpretation | **OWN / DERIVE** from authoritative realized facts, including material delivery/post-sale adjustments |
| Audit/reconciliation records for MPC-controlled operations | **OWN** |

### D0.5b — Boundary invariants

1. Observing an external fact does not transfer ownership of that fact to MPC.
2. Orchestrating an external process does not make MPC source of truth for the external system’s native record.
3. An MPC-derived conclusion preserves provenance/freshness needed to understand which authoritative facts/rules produced it.
4. Externally governed rules are not silently copied into mutable MPC-owned policy.
5. MPC-owned workflow state is not replaced by guessing from one provider response; ambiguous/divergent outcomes remain explicit until reconciled.
6. **Unknown availability is not zero.** If MPC cannot determine sellable availability with sufficient confidence, it must not invent a plausible quantity merely to continue synchronization; uncertainty becomes explicit operational state.
7. A routine policy-valid availability change does not require human approval merely because stock changed; failed/non-convergent/uncertain updates become exception work.
8. Provider/ERP-specific mechanisms remain implementation concerns for later stages.

## 9. D0.6 — Product 1.0 completion / user-observable outcomes

**Accepted by operator.**

Product 1.0 is complete only when MPC is demonstrably usable as the normal operational control plane, not merely when individual capabilities exist in isolation.

The acceptance bar is user-observable:

1. **Attention is portfolio-driven, not manual-search driven.** Operators can see what is healthy, changed, divergent, blocked, approval-required or otherwise actionable without inspecting products/external systems one by one.
2. **An eligible internal product can reach a verified marketplace state through MPC.** The normal path covers readiness/linkage, commercial analysis, creation/publication and observation of the real channel state.
3. **Marketplace availability remains operationally coherent without per-change manual work.** When governing stock/rules change and sellable availability is known, MPC updates marketplace availability on the normal path and verifies convergence; uncertainty or failure becomes explicit work rather than guessed quantity or silent drift.
4. **Competitive/pricing intelligence replaces major manual analysis.** MPC exposes grounded comparable-market position, relevant internal economics, price scenarios, expected profitability and insufficient-evidence cases at portfolio and product level.
5. **Decision closes into controlled action.** Authorized human or bounded policy-driven decisions can become external actions with policy enforcement, auditability, verification and reconciliation.
6. **A marketplace sale traverses the normal operating loop through MPC.** Order recognition, ERP operation, fulfillment queue, conference, invoicing trigger/verification, packing and dispatch do not require hidden manual system hopping as a normal step.
7. **Delivery remains visible through terminal outcome.** After dispatch, MPC continues observing until delivered, returned, cancelled or equivalent terminal state; material delivery exceptions become explicit work.
8. **Essential post-sale changes remain inside the controlled lifecycle.** Cancellation, return/refund effects can be progressed through necessary cross-system response/reconciliation without dropping normal operation back to manual system hopping.
9. **Failures become explicit work.** Missing evidence, ambiguous external results, integration failures and physical/order/availability/delivery/post-sale divergences surface with what is known/unknown and what requires action.
10. **The economic loop closes.** MPC compares expected versus realized profitability using attributable authoritative facts, including material delivery/cancellation/return/refund effects.
11. **Organizational governance is operable without code edits.** Actors can exercise legitimate MPC-owned authorities while externally governed rules remain externally governed and mandatory safety/audit invariants cannot be configured away.

Completion statement:

> **A company can take its internal products, determine marketplace readiness, publish and operate offers, keep marketplace availability coherent with governing stock/rules, monitor market position and profitability at portfolio scale, make and execute decisions under policy, receive sales, progress them through ERP and physical fulfillment, follow delivery to a terminal outcome, handle essential cancellation/return/refund consequences, surface/reconcile exceptions, and understand the realized economic result — using Marketplace Central as the normal marketplace operations control plane.**

### D0.6a — Normal-path rule

The normal operational path must be executable through MPC for responsibilities Product 1.0 claims to control.

Direct use of Mercado Livre, Sankhya or another participating external system remains legitimate when:

- the responsibility inherently belongs to that external system and is intentionally outside MPC scope;
- investigation/support requires native-system inspection;
- an exceptional recovery path explicitly requires it.

Direct external-system hopping must **not** be a hidden required step in an otherwise claimed MPC normal workflow.

Detailed end-to-end proof scenarios remain a D8 responsibility; D0 defines the outcomes D8 must eventually prove.

## 10. D0.7 — Product completeness / contradiction review

**OPEN — exact next D0 work.**

Before D0 can be accepted as a whole, review the accepted Product 1.0 definition adversarially for missing business lifecycle responsibilities, contradictions and accidental scope gaps.

### D0.7a — Essential post-sale lifecycle

**Accepted by operator.**

Product 1.0 includes essential operational handling of **cancellations, returns and refunds** when they affect an MPC-controlled marketplace sale. MPC coordinates the cross-system response, surfaces unresolved outcomes as explicit exceptions and reconciles operational/economic consequences. Marketplace/provider remains authoritative for native marketplace facts; ERP remains authoritative for native fiscal/accounting/reversal records.

This does not make Product 1.0 general CRM/SAC, buyer messaging, complaint/mediation automation, reputation platform or company-wide reverse logistics.

### D0.7b — Shipment / delivery lifecycle after dispatch handoff

**Accepted by operator.**

Product 1.0 continues observing shipment/delivery after physical dispatch handoff until a relevant terminal outcome such as delivered, returned, cancelled or equivalent provider state. Delays, failed attempts, returns-to-sender and material delivery exceptions become MPC operational work.

This does not make Product 1.0 a carrier, route planner, fleet system, company-wide transport manager or TMS.

### D0.7c — Stock / marketplace availability control

**Accepted by operator.**

Product 1.0 is responsible for **maintaining marketplace availability coherently with sellable availability derived from the authoritative stock facts and governing rules applicable to that offer**.

The normal path is automatic:

```text
authoritative stock facts / reservations / governing rules
  → derive sellable marketplace availability
  → policy-valid MPC intent
  → update marketplace availability
  → observe / verify convergence
  → success or explicit exception
```

MPC does not become the owner of physical stock merely because it controls marketplace availability. Stock, reservation and other source facts remain owned by their authoritative system; some availability rules may be ERP-governed, some MPC-owned, and some availability values derived.

Routine, sufficiently-known and policy-valid stock/availability changes do not require human approval per change. Human intervention is for uncertainty, conflict, policy violation, failed update or non-convergence.

**Unknown availability must not silently become zero or another plausible quantity.** If the governing facts/rules are insufficient, MPC must represent the uncertainty explicitly and follow the later-accepted safety/reconciliation policy rather than inventing availability.

Exact stock authority, reservation semantics, buffers, source reads, event/sync strategy and provider update mechanisms belong to D2–D4/D7.

### Next D0.7 question

The review continues. The next material product-boundary question is **marketplace account / installation multiplicity inside one organization**.

D0 must decide whether Product 1.0 is conceptually limited to one Mercado Livre seller account/installation per organization or whether one organization may operate multiple seller accounts/installations under the same MPC control plane. This affects downstream identity, policy scope, linkage, order attribution, integration ownership and data isolation, so it must not be invented locally in D2/D4.

This is a product cardinality decision only. Exact identity keys, tenant schema, credentials, connection models and routing belong to D2/D4.

Other findings discovered during D0.7 are classified as MUST DECIDE NOW, SHOULD DECIDE NOW or CAN DEFER SAFELY. D0 closes only when no material Product 1.0 semantic is being left for implementation to invent.

## 11. Resume contract for a fresh session

A fresh session must read `AGENTS.md`, `docs/engineering/rebaseline/README.md`, the canonical engineering method, `ARCHITECTURE.md`, the ADR registry, and then this D0 artifact.

It should conclude:

- D0 is open and not yet accepted as a whole;
- D0.1–D0.6 are operator-approved product decisions;
- D0.7a is accepted: essential cancellation/return/refund operations remain inside the controlled sale lifecycle without expanding Product 1.0 into general CRM/SAC;
- D0.7b is accepted: shipment/delivery remains visible through terminal outcome and material delivery exceptions become MPC operational work without turning MPC into a TMS;
- D0.7c is accepted: marketplace availability is maintained automatically from governing authoritative stock/rules when sufficiently known, while uncertainty/failure becomes explicit work and MPC does not become physical-stock authority;
- MPC uses **OWN / ORCHESTRATE / OBSERVE-DERIVE** and owns the marketplace operating model without taking ownership of external facts merely because it consumes or causes them;
- Sankhya API availability for writes is evidence to carry into D4, not a D0 transport decision;
- business rules/policies may be MPC-owned, externally governed or derived, and later D2/D4 must preserve that provenance;
- no D1+ target architecture may be invented yet;
- the exact next work is **D0.7 Product completeness review — marketplace account / installation multiplicity inside one organization**.
