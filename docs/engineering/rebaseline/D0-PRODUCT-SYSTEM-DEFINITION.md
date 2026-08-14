# D0 — Product / System Definition

> **Status:** WORKING — D0 OPEN, NOT YET ACCEPTED AS A WHOLE  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Authority:** operator-approved D0 decisions recorded here; later D-stages own technical realization  
> **Last updated:** 2026-08-14

## 1. Purpose and boundary

D0 defines what Marketplace Central is, who it serves, which problems and outcomes belong to Product 1.0, and where the product boundary ends.

D0 does **not** decide target contexts/modules, canonical identity keys, database schema, API shape, frontend topology, runtime/jobs, event/outbox mechanics or provider-specific transport. Those belong to D1–D7.

Current code, schemas, APIs and historical ADRs are evidence only unless the active rebaseline marks a constraint binding.

---

## 2. D0.1 — Product mission

**Accepted by operator.**

Marketplace Central is the internal **Marketplace Operations Control Plane** of the company.

It combines internal business facts with marketplace observations so operators can understand the real operational state, detect divergences, make grounded decisions, execute controlled actions in participating systems, and verify/reconcile the result afterward.

Its fundamental loop is:

```text
observe → understand → reconcile → decide/policy → execute → verify → audit/reconcile
```

Marketplace Central is not merely an intelligence dashboard and is not a replacement ERP or marketplace platform. It can coordinate and execute operations on both sides of the boundary, including marketplace listing/price actions and marketplace-originated order/invoicing workflows through an accepted ERP integration path.

External systems remain authorities for facts that inherently belong to them. MPC owns the cross-system operational control semantics: intent, policy application, workflow, correlation, controlled execution, operational state, divergence, audit and reconciliation.

---

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
3. **Marketplace Availability Control** — keep published marketplace availability coherently aligned with sellable availability derived from authoritative stock facts/rules and applicable MPC-owned allocation policy, automatically on the normal path and with explicit exceptions when evidence or convergence is uncertain.
4. **Competitive Intelligence** — observe comparable market offers/prices, expose competitive position and meaningful changes, and represent insufficient comparison evidence honestly.
5. **Pricing & Profitability Intelligence** — combine relevant internal economics under an explicit, explainable Cost Basis with market observations to calculate price scenarios, expected margin/profitability and decision-relevant trade-offs.
6. **Decision & Policy Control** — translate observations/recommendations into permitted, approval-required or prohibited actions according to governing company rules/policies.
7. **Order-to-ERP Operations** — receive/understand marketplace orders and coordinate corresponding operations in the participating ERP, including order creation and invoicing where authorized.
8. **Marketplace Fulfillment / Dispatch** — progress marketplace orders through eligible fulfillment execution, including physical separation, conference, invoicing trigger when applicable, packing and dispatch handoff, without becoming a company-wide WMS.
9. **Shipment / Delivery Observation & Exceptions** — continue observing shipment after dispatch handoff until a relevant terminal outcome, surfacing delays, delivery failures, returns or other material exceptions without becoming a TMS/carrier platform.
10. **Essential Post-Sale Operations** — control the operational response to marketplace cancellations, returns and refunds when they affect an MPC-controlled sale, coordinating consequences across marketplace/ERP/fulfillment/economic workflows without becoming general CRM/SAC.
11. **Reconciliation & Exception Operations** — identify cases where an expected cross-system result cannot be proven or systems diverge, making them explicit work instead of silently assuming success.
12. **Realized Profitability** — determine realized sale economics from attributable authoritative/derived economic facts, including the relevant Cost Basis and material delivery/post-sale reversals or adjustments, and compare expected versus realized results.

The provider-independent requirement is to express these capabilities in business terms. Mercado Livre and Sankhya are the first concrete systems used to prove relevant capabilities; provider-specific mechanics belong to D4.

### D0.3a — Action authority model

**Accepted by operator.**

Product 1.0 supports:

- **human-controlled execution by default:** recommendation/analysis → operator approval → MPC executes → verifies/reconciles;
- **policy-driven automatic execution when explicitly authorized:** an accepted policy may approve a bounded action automatically;
- **human review for exceptions, uncertainty, low confidence or policy violations.**

A fully autonomous repricing engine is not required as a launch gate. The product must nevertheless support explicit bounded automation without bypassing policy, audit or reconciliation.

Marketplace availability synchronization is different from discretionary commercial repricing: routine stock/availability changes are expected to execute automatically when governing facts/rules are sufficiently known and the action is valid. Human attention is concentrated on uncertainty, policy conflicts and failed/non-convergent updates.

### D0.3b — Policy/rule provenance is explicit

**Accepted by operator.**

Marketplace Central must not assume every business rule or commercial policy is authored inside MPC.

A governing rule used by MPC may be:

- **MPC-owned** — intentionally defined and governed inside Marketplace Central;
- **externally governed** — sourced from an ERP or another authoritative system and consumed by MPC;
- **derived** — mechanically computed from authoritative facts/rules without becoming an independent source of truth.

MPC must preserve enough provenance to distinguish these classes and must not silently turn an externally governed rule into an editable MPC-owned copy.

For policy classes that are legitimately MPC-owned, Product 1.0 must support organization-operable configuration without requiring code edits. The policy model must be capable of inherited/default policy plus explicit more-specific overrides where a later accepted business scope requires them. The effective policy and its provenance must be explainable; precedence must not depend on hidden evaluation order.

D0 does **not** enumerate every future policy or mandate every conceivable scope. Candidate scopes such as organization/default, marketplace installation, product group/category, product and offer are carried forward for later adjudication and are only implemented when a real business need justifies them.

The exact ownership matrix, policy schemas/scope hierarchy, synchronization mechanism, conflict semantics, freshness rules and provider contracts are deferred to later stages.

---

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

---

## 5. Operator-provided Sankhya evidence / constraint

This is **D0/D4 evidence**, not a provider-transport decision and not a canonical MPC domain model.

The operator reports an already-proven Sankhya application integration in another app using application credentials (`client id`, `client secret`, `X-Token`; no secret values are recorded here).

Operational writes such as order creation and invoicing can be performed through Sankhya's application API. The operator prefers this path for writes rather than direct Oracle database writes because it is a safer system-owned operation boundary. The environment also has DB Explorer access for database inspection.

Consequences for later design:

- Product 1.0 may require MPC to create/invoice orders in the participating ERP without assuming direct database writes;
- D4 must evaluate and ratify exact Sankhya read/write capability contracts and transport boundaries for the first ERP integration;
- later stages must account for rules/policies whose authority remains in Sankhya/another ERP rather than duplicating them as MPC-owned configuration;
- existing binding Oracle-read constraints are not silently reopened here;
- Sankhya-specific constructs such as `CODEMP`, `CODLOC`, cost fields/types or other native identifiers are **integration evidence only** unless a later MPC-domain decision independently proves that an equivalent business concept belongs in the canonical product model.

---

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

---

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

Owns physical fulfillment execution for marketplace orders when the selected Fulfillment Node is operated by the organization and the work is assigned to this actor: work queue, separation/conference, invoicing trigger when valid, packing, dispatch handoff and exception reporting.

Accepted normal sequence for the internal-fulfillment path:

```text
marketplace order / ERP order readiness
  → eligible / selected Fulfillment Node
  → fulfillment queue
  → physical separation
  → physical conference
  → if valid: operator triggers invoicing through MPC
  → MPC causes ERP invoicing through the later-accepted integration boundary
  → invoicing result is verified
  → packing
  → dispatch handoff
  → completion or exception
```

**Invariant:** on the accepted internal-fulfillment normal path, an order is not intentionally invoiced before the operator has physically confirmed that the correct items are available and separated. Missing stock, wrong item, damaged material, quantity divergence or another physical inconsistency blocks normal invoicing and becomes an operational exception.

A Fulfillment Node may later represent an externally operated fulfillment capability where these physical steps are performed by a provider rather than this human actor; D0.7e.3 defines the node semantics without requiring all such modes in the first deployment.

### Commercial / Marketplace Manager

Is the ordinary commercial authority for the marketplace operation, but only over policies/decisions actually within that actor’s authority. This actor may:

- define/change **MPC-owned** margin floors, price boundaries, approval thresholds and marketplace commercial policies;
- define/change MPC-owned marketplace availability-allocation policies within delegated authority, including policies that intentionally expose only a bounded share of otherwise eligible stock;
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

---

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
| Internal master product facts, ERP-native cost/fiscal facts | **OBSERVE / CONSUME** as external facts; MPC may derive normalized economic conclusions without taking ownership of the source facts |
| ERP/WMS/3PL/provider-native physical/on-hand/reserved/other stock facts | **OBSERVE / CONSUME** |
| Business rules/policies whose authority remains in an ERP/another system | **OBSERVE / CONSUME** |
| Product ↔ marketplace linkage/evidence maintained by MPC | **OWN** |
| Marketplace readiness/readiness assessment | **OWN / DERIVE** |
| MPC Inventory Source identity/membership | **OWN** as MPC inventory semantics/configuration; native source identities and stock facts retain their external authority |
| Inventory Scope / source eligibility for an offer or policy context | **OWN** as MPC operating policy/configuration |
| MPC-owned availability-allocation policy | **OWN** |
| Sellable marketplace availability derived from eligible sources, authoritative stock/rules and applicable MPC-owned allocation policy | **OWN / DERIVE** as an MPC operating conclusion; underlying stock facts/rules retain their source authority |
| MPC Fulfillment Node identity/membership | **OWN** as MPC fulfillment-operating semantics/configuration; native facility/provider identities and physical execution facts retain their relevant external authority |
| Fulfillment Scope / node eligibility for an order, offer or policy context | **OWN** as MPC operating policy/configuration |
| Fulfillment-node selection/routing intent once governing policy decides it | **OWN / ORCHESTRATE**; exact routing algorithm is deferred |
| Cost Observation normalized for MPC economic use | **OBSERVE / DERIVE** from authoritative economic facts with meaning, time/context and provenance preserved |
| Cost Basis / cost-selection semantics | **OWN** when intentionally governed by MPC; **OBSERVE / CONSUME** when externally governed; the effective basis and provenance must remain explicit |
| Intent to publish/update marketplace quantity/availability | **OWN / ORCHESTRATE** |
| Actual marketplace quantity/availability state | provider authoritative; MPC **OBSERVES / ORCHESTRATES** convergence |
| Cross-system intent, workflow, correlation, divergence and exception state | **OWN** |
| Competitive intelligence / comparable-market interpretation | **OWN / DERIVE** |
| Pricing scenarios / expected profitability | **OWN / DERIVE** |
| MPC-owned marketplace commercial policies | **OWN** |
| Marketplace account/installation membership under an MPC organization | **OWN** as MPC organization/control-plane configuration; provider account identity itself remains provider-owned |
| MPC Selling Entity identity/membership and transaction attribution | **OWN** as MPC business/control-plane semantics; an ERP or legal registry may remain authoritative for corresponding external legal/company identities |
| Actual listing/channel state inherent in marketplace | provider authoritative; MPC **OBSERVES / ORCHESTRATES** |
| Marketplace listing/price mutation intent | **OWN / ORCHESTRATE** |
| Marketplace-originated order facts | marketplace/provider authoritative; MPC **ORCHESTRATES** |
| ERP-native order/invoice/accounting facts | participating ERP authoritative; MPC **ORCHESTRATES** marketplace workflow around them |
| Marketplace-order fulfillment workflow | **OWN**, while native ERP/WMS/provider/carrier execution facts retain their external authority where applicable |
| Shipment/delivery state after dispatch | marketplace/carrier/provider authoritative; MPC **OBSERVES / ORCHESTRATES** exceptions/lifecycle closure |
| Marketplace cancellation/return/refund facts | marketplace/provider authoritative; MPC **ORCHESTRATES** cross-system response |
| ERP-native reversal/credit/fiscal/accounting facts | participating ERP authoritative; MPC **ORCHESTRATES** post-sale workflow |
| Realized profitability interpretation | **OWN / DERIVE** from attributable authoritative/derived economic facts, including the applicable Cost Basis and material delivery/post-sale adjustments |
| Audit/reconciliation records for MPC-controlled operations | **OWN** |

### D0.5b — Boundary invariants

1. Observing an external fact does not transfer ownership of that fact to MPC.
2. Orchestrating an external process does not make MPC source of truth for the external system’s native record.
3. An MPC-derived conclusion preserves provenance/freshness needed to understand which authoritative facts/rules produced it.
4. Externally governed rules are not silently copied into mutable MPC-owned policy.
5. MPC-owned workflow state is not replaced by guessing from one provider response; ambiguous/divergent outcomes remain explicit until reconciled.
6. **Unknown availability is not zero.** If MPC cannot determine sellable availability with sufficient confidence, it must not invent a plausible quantity merely to continue synchronization; uncertainty becomes explicit operational state.
7. A routine policy-valid availability change does not require human approval merely because stock changed; failed/non-convergent/uncertain updates become exception work.
8. An organization and a marketplace seller account/installation are not the same identity merely because the first deployment uses one account.
9. Any listing/order/action/policy whose meaning depends on a marketplace installation must be attributable to the correct installation; no implicit “current seller account” may become product semantics.
10. **MPC canonical business semantics are defined from marketplace-operating needs, not copied from an ERP ontology.** Native ERP identifiers, tables, fields and taxonomies remain at the integration boundary unless the MPC domain independently requires an equivalent concept.
11. ERP integration is semantic translation, not merely field renaming: an adapter maps authoritative facts/capabilities/commands/results between the ERP’s native model and MPC semantics and must represent unsupported or uncertain mappings explicitly rather than inventing equivalence.
12. **Selling Entity is a canonical MPC business concept, not an alias for ERP company/branch.** When a marketplace operation materially depends on which business/legal/fiscal entity is acting, that attribution must be explicit and preserved.
13. Selling Entity must not be silently collapsed with Inventory Source, Fulfillment Node, cost scope or another business dimension merely because a particular ERP or first deployment represents them together. Any relationship between those dimensions must be explicit and justified by business semantics.
14. **Inventory Source is a canonical MPC inventory-origin concept, not an alias for ERP warehouse/location/company.** One MPC Inventory Source may require one or many native external structures to represent it, and one native external structure may contribute to multiple MPC-relevant semantics when justified.
15. **Inventory Scope is explicit eligibility, not “all stock we can find”.** Only Inventory Sources allowed by the offer’s governing scope/policy may contribute to Sellable Availability; stock outside that scope must not leak into marketplace availability merely because it exists.
16. Inventory Source must not be silently collapsed with Selling Entity or Fulfillment Node. Their relationships are explicit business rules/configuration unless a later accepted decision proves stronger identity.
17. **Availability Allocation Policy changes allocation, not stock truth.** An MPC-owned policy may intentionally expose less than full eligible availability — for example a percentage, reserve or cap — without rewriting authoritative inventory facts.
18. MPC-owned policy defaults/overrides must resolve deterministically and retain visible provenance. Conflicting or ambiguous effective policy becomes explicit configuration/exception state rather than hidden ordering behavior.
19. **Fulfillment Node is a canonical MPC execution concept, not an alias for Inventory Source, ERP warehouse/location or physical address.** It represents the business-recognized fulfillment execution point/capability responsible for marketplace-order physical work; it may be internally operated or externally operated by a provider.
20. **Fulfillment Scope is explicit eligibility.** Only Fulfillment Nodes allowed by the governing order/offer/policy context may be selected for the operation. No implicit “default warehouse” may become product semantics when more than one node is material.
21. A Fulfillment Node may use inventory associated with one or more Inventory Sources, but that relationship is explicit. Inventory promise and physical execution are different semantics even when the first deployment uses the same facility for both.
22. **Cost Observation is not a bare number.** Any cost used materially by MPC must preserve enough economic meaning, amount/currency, applicable temporal/business context and provenance to explain what the number represents and where it came from.
23. **Cost Basis is not an ERP cost-type alias.** It expresses which cost semantic is economically appropriate for a particular analysis/decision/transaction. An adapter must not choose a convenient native cost field/type silently when the required basis is unsupported or ambiguous.
24. **Missing or ambiguous cost does not silently fall back.** Unsupported basis, missing historical evidence or conflicting mappings become explicit uncertainty/configuration/exception state unless a later accepted policy explicitly authorizes a specific fallback semantic.
25. **Current cost is not historical transaction cost by default.** Expected and realized profitability must preserve the temporal/economic evidence required by their respective use cases; recomputing a past sale using a convenient current cost without explicit semantics is prohibited.
26. **Cost is not the entire sale economy.** Marketplace fees, freight effects, taxes, discounts, subsidies, reversals and other material economic components remain separately attributable facts/conclusions rather than being hidden inside an opaque “cost” number.
27. Provider/ERP-specific mechanisms remain implementation concerns for later stages.

---

## 9. D0.6 — Product 1.0 completion / user-observable outcomes

**Accepted by operator.**

Product 1.0 is complete only when MPC is demonstrably usable as the normal operational control plane, not merely when individual capabilities exist in isolation.

The acceptance bar is user-observable:

1. **Attention is portfolio-driven, not manual-search driven.** Operators can see what is healthy, changed, divergent, blocked, approval-required or otherwise actionable without inspecting products/external systems one by one.
2. **An eligible internal product can reach a verified marketplace state through MPC.** The normal path covers readiness/linkage, commercial analysis, creation/publication and observation of the real channel state.
3. **Inventory eligibility and allocation are explicit and marketplace availability remains operationally coherent without per-change manual work.** MPC knows which Inventory Sources may contribute to an offer, applies authoritative rules and applicable MPC-owned allocation policy, derives Sellable Availability, updates the marketplace and verifies convergence. Uncertainty, missing mapping, policy ambiguity or failure becomes explicit work rather than guessed quantity, unintended stock aggregation or silent drift.
4. **Competitive/pricing intelligence replaces major manual analysis with explainable economics.** MPC exposes grounded comparable-market position, relevant internal economics, price scenarios and expected profitability using an explicit Cost Basis and attributable Cost Observations/economic facts; insufficient or ambiguous cost evidence is visible rather than silently substituted.
5. **Decision closes into controlled action.** Authorized human or bounded policy-driven decisions can become external actions with policy enforcement, auditability, verification and reconciliation.
6. **Material transaction context is explicit.** When a marketplace operation depends on which business/legal/fiscal entity is acting, MPC can attribute the workflow to the correct Selling Entity rather than relying on an implicit ERP/default-company assumption.
7. **Fulfillment responsibility is explicit.** When fulfillment execution location/capability is material, MPC knows which Fulfillment Nodes are eligible and which node is responsible for the operation rather than relying on an implicit ERP warehouse/default shipping point.
8. **A marketplace sale traverses the normal operating loop through MPC.** Order recognition, ERP operation, eligible/selected fulfillment, conference, invoicing trigger/verification when applicable, packing and dispatch do not require hidden manual system hopping as a normal step.
9. **Delivery remains visible through terminal outcome.** After dispatch, MPC continues observing until delivered, returned, cancelled or equivalent terminal state; material delivery exceptions become explicit work.
10. **Essential post-sale changes remain inside the controlled lifecycle.** Cancellation, return/refund effects can be progressed through necessary cross-system response/reconciliation without dropping normal operation back to manual system hopping.
11. **Failures become explicit work.** Missing evidence, ambiguous external results, integration failures and physical/order/availability/fulfillment/delivery/post-sale divergences surface with what is known/unknown and what requires action.
12. **The economic loop closes.** MPC compares expected versus realized profitability using attributable economic facts and explicit cost semantics appropriate to each temporal/transaction context, including material delivery/cancellation/return/refund effects, so the result can be explained rather than merely recomputed from current values.
13. **Organizational governance is operable without code edits.** Actors can exercise legitimate MPC-owned authorities while externally governed rules remain externally governed and mandatory safety/audit invariants cannot be configured away.

Completion statement:

> **A company can take its internal products, determine marketplace readiness, publish and operate offers, derive marketplace availability from explicitly eligible inventory sources, governing facts/rules and MPC-owned allocation policies, preserve the correct business entity context, route physical fulfillment through an eligible recognized fulfillment node, analyze price/profitability with explainable cost semantics and attributable economics, make and execute decisions under policy, receive sales, progress them through ERP and physical fulfillment, follow delivery to a terminal outcome, handle essential cancellation/return/refund consequences, surface/reconcile exceptions, and understand the realized economic result — using Marketplace Central as the normal marketplace operations control plane.**

### D0.6a — Normal-path rule

The normal operational path must be executable through MPC for responsibilities Product 1.0 claims to control.

Direct use of Mercado Livre, Sankhya or another participating external system remains legitimate when:

- the responsibility inherently belongs to that external system and is intentionally outside MPC scope;
- investigation/support requires native-system inspection;
- an exceptional recovery path explicitly requires it.

Direct external-system hopping must **not** be a hidden required step in an otherwise claimed MPC normal workflow.

Detailed end-to-end proof scenarios remain a D8 responsibility; D0 defines the outcomes D8 must eventually prove.

---

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

Product 1.0 is responsible for **maintaining marketplace availability coherently with Sellable Availability derived from authoritative stock facts, governing rules and applicable MPC-owned allocation policy**.

The normal path is automatic:

```text
authoritative stock facts / reservations / governing rules
  → Inventory Scope / eligible Inventory Sources
  → applicable MPC-owned allocation policy
  → derive Sellable Availability
  → policy-valid MPC intent
  → update marketplace availability
  → observe / verify convergence
  → success or explicit exception
```

MPC does not become owner of physical stock merely because it controls marketplace availability. Stock, reservation and other source facts remain owned by their authoritative system; some rules may be externally governed, some MPC-owned, and some values derived.

Routine, sufficiently-known and policy-valid stock/availability changes do not require human approval per change. Human intervention is for uncertainty, conflict, policy violation, failed update or non-convergence.

**Unknown availability must not silently become zero or another plausible quantity.** If governing facts/rules are insufficient, MPC must represent uncertainty explicitly rather than inventing availability.

Exact stock authority, reservation semantics, allocation arithmetic/rounding, buffers, source reads, event/sync strategy and provider update mechanisms belong to later stages.

### D0.7d — Marketplace account / installation multiplicity

**Accepted by operator.**

A Marketplace Central organization may operate **one or more marketplace seller accounts/installations** under the same control plane.

Product-level cardinality is therefore:

```text
Organization 1 → N Marketplace Installations
```

The first real deployment may validate Product 1.0 using only one Mercado Livre seller account. That does **not** authorize collapsing organization identity into marketplace-account identity or hardcoding single-account semantics into the target model.

Product-level requirements:

- listings, orders, external actions, reconciliation and other channel facts/workflows whose meaning depends on the seller account must be attributable to the correct marketplace installation;
- policies may later be scoped at organization level, installation level or another explicitly accepted scope; D0 does not invent the exact hierarchy;
- consolidated organization-level views may combine multiple installations without erasing their source/account provenance;
- installation/account membership is an MPC organization/control-plane concern, while the provider remains authoritative for its native seller-account identity and state;
- exact identity keys, credential storage, connection lifecycle, routing and isolation mechanisms belong to D2/D4.

This decision establishes the correct cardinality without making multi-account feature breadth a launch blocker: Product 1.0 may initially operate one real account while retaining a model in which adding another installation does not require redefining what an organization, listing or order means.

### D0.7e — ERP-agnostic business semantics before ERP mapping

**Accepted by operator.**

Marketplace Central must define the business concepts it needs from **marketplace-operating semantics first**, independently of how Sankhya or any future ERP happens to represent them.

The governing direction is:

```text
business / MPC need
  → canonical MPC semantic concept or capability
  → integration contract
  → ERP-specific semantic mapping
      ├── Sankhya representation
      └── another ERP representation
```

The reverse direction is prohibited as a design shortcut:

```text
ERP has field/entity X
  → therefore MPC must have canonical concept X
```

ERP integration is therefore a **semantic translation boundary**, not a database-field mirroring exercise. The adapter may need to combine several ERP-native structures to satisfy one MPC concept, or map one ERP structure into several MPC-relevant facts. If an ERP cannot provide a required MPC semantic reliably, the integration must expose unsupported/incomplete/unknown evidence rather than inventing an equivalence.

Examples such as Sankhya `CODEMP`, `CODLOC`, native cost variants, stock structures and fiscal constructs are evidence later stages may inspect; they are not the vocabulary from which D0/D1 define the MPC domain.

This principle does **not** require building a speculative universal ERP abstraction. MPC should define only the business semantics actually required by accepted marketplace workflows.

#### D0.7e.1 — Selling Entity

**Accepted by operator.**

`Selling Entity` is the first ERP-independent canonical business dimension accepted for MPC.

A **Selling Entity** is the business entity to which a marketplace commercial operation belongs when that distinction is material to who is acting as seller and to the transaction's business, legal or fiscal consequences.

Product-level cardinality permits:

```text
Organization 1 → N Selling Entities
```

The purpose of the concept is to answer, in MPC language:

> **Which business/legal/fiscal entity of this organization is acting for this marketplace transaction?**

Product-level requirements:

- an order, invoice-triggering workflow, reversal or other marketplace operation whose business/fiscal meaning depends on the acting entity must be attributable to the correct Selling Entity;
- organization identity is not Selling Entity identity merely because a first deployment has one legal/company entity;
- Selling Entity is not defined as an ERP company/branch identifier and does not inherit an ERP ontology;
- a marketplace installation is not inherently bound one-to-one to a Selling Entity; any routing relationship between them is explicit policy/configuration, not structural identity;
- Selling Entity does **not** automatically determine inventory source, fulfillment node, costing scope or another operational dimension. Those dimensions remain separate until their own business semantics are decided;
- a participating ERP adapter must translate the MPC Selling Entity semantic to the ERP-native structures needed to execute/read the operation. If that mapping is missing or ambiguous, the adapter/workflow must expose unsupported/unknown/configuration-required state rather than selecting an arbitrary default entity.

The first deployment may use one Selling Entity without hardcoding single-entity semantics into the product.

#### D0.7e.2 — Inventory Source and Inventory Scope

**Accepted by operator.**

`Inventory Source` and `Inventory Scope` are canonical MPC inventory semantics defined independently of any ERP/WMS/provider storage model.

An **Inventory Source** is a business-recognized source or pool of inventory whose authoritative stock facts may be eligible to contribute to the sellable availability of a marketplace offer.

An **Inventory Scope** is the explicit eligibility definition — directly or through governing policy — that determines which Inventory Sources may contribute to a particular offer/operation and which must not.

The purpose is to answer, in MPC language:

> **Which authoritative inventory sources are allowed to contribute to what we promise as sellable for this marketplace offer?**

Product-level requirements:

- an Inventory Source is not defined as an ERP warehouse/location/company, physical address, marketplace-fulfillment warehouse or any other provider-native structure;
- authoritative stock facts for an Inventory Source may come from an ERP, WMS, marketplace-fulfillment provider, 3PL or another accepted authority;
- an ERP/provider mapping is semantic rather than cardinality-preserving: a single MPC Inventory Source may map to one or many native structures, and multiple native structures may need to be combined to represent one source; no `1:1` mapping is assumed;
- Inventory Scope determines eligibility. Stock from a source outside the governing scope must not contribute to Sellable Availability merely because that stock exists or is technically reachable;
- Sellable Availability is a derived MPC operating conclusion from eligible Inventory Sources, their authoritative stock/reservation/availability facts and the applicable external/MPC-owned rules, buffers or policies;
- Inventory Source is not Selling Entity and is not Fulfillment Node. A business may relate these concepts explicitly, but D0 does not collapse them into one identity;
- if source mapping, source eligibility or source facts are incomplete/ambiguous, MPC must surface uncertainty/configuration-required state instead of silently aggregating stock from an arbitrary/default location;
- the first deployment may use one Inventory Source and a trivial Inventory Scope without hardcoding single-source semantics into Product 1.0.

Exact source identities, native ERP/WMS mappings, stock ledger semantics, reservations, aggregation formulas, source precedence, buffers and synchronization/event mechanics belong to D2/D4/D7.

##### D0.7e.2a — Availability Allocation Policy

**Accepted as a D0 product requirement from operator-provided use case; detailed policy catalog deferred.**

MPC must support configurable **MPC-owned availability-allocation policies** that intentionally expose less than the full otherwise-eligible availability when the organization chooses to reserve capacity/stock for other channels or operational needs.

Representative use case:

```text
eligible availability = 100
MPC marketplace allocation policy = 70%
→ marketplace allocation is derived from 70% of eligible availability
```

The example establishes the policy class, **not** the final arithmetic contract. Rounding, minimum/maximum caps, fixed reserves, safety stock, sequencing with other rules and other allocation forms remain later design decisions.

Product-level requirements:

- allocation policy modifies the MPC marketplace-allocation conclusion; it never rewrites the authoritative stock fact;
- MPC-owned allocation policy is organization-operable without code changes;
- the eventual policy model must support inherited/default configuration and explicit more-specific overrides where justified by real business scope;
- the effective policy, its source/provenance and why a more-specific override won must be observable;
- D0 does not freeze every supported scope. Candidate scopes include organization/default, marketplace installation, product group/category, product and offer; later stages determine the minimal justified set rather than implementing all speculatively;
- conflicts or ambiguity in effective policy become explicit configuration/exception state rather than being resolved by hidden ordering.

This is intentionally narrower than defining a generic rules engine. Product 1.0 needs configurable policy for accepted marketplace-operating semantics; later stages design the smallest policy model that satisfies those concrete needs.

#### D0.7e.3 — Fulfillment Node and Fulfillment Scope

**Accepted by operator.**

`Fulfillment Node` and `Fulfillment Scope` are canonical MPC fulfillment semantics defined independently of any ERP/WMS/provider location model.

A **Fulfillment Node** is a business-recognized fulfillment execution point or capability responsible for performing marketplace-order physical work such as separation, conference, packing and dispatch handoff, according to the capabilities and operating model of that node.

A Fulfillment Node is not necessarily only a physical address or internally operated warehouse. It may represent, where a real accepted workflow requires it, an internal CD/store/dispatch operation, a 3PL fulfillment operation, marketplace-managed fulfillment or another recognized execution capability. This semantic breadth does **not** require Product 1.0 to implement every fulfillment model in its first deployment.

A **Fulfillment Scope** is the explicit eligibility definition — directly or through governing policy — that determines which Fulfillment Nodes may execute a given marketplace order/offer/workflow and which may not.

The purpose is to answer, in MPC language:

> **Which recognized fulfillment execution points/capabilities are allowed to perform this marketplace order, and which node is responsible once one is selected?**

Product-level requirements:

- Fulfillment Node is not defined as an ERP warehouse/location/company, Inventory Source, physical address or provider-native fulfillment identifier;
- node identity and eligibility belong to MPC operating semantics, while native facility/provider identities and physical execution facts remain with their appropriate authoritative systems;
- one Fulfillment Node may map to one or many native ERP/WMS/provider structures, and an external native structure may participate in more than one MPC semantic when justified; no `1:1` mapping is assumed;
- Fulfillment Scope determines eligibility; an ineligible node must not be selected merely because it is technically reachable or contains stock;
- selecting/routing among eligible nodes is a later policy/decision problem; D0 establishes that eligibility and the responsible node must be explicit when materially relevant, not that a particular routing algorithm exists;
- a Fulfillment Node may be related to one or more Inventory Sources, but that relationship is explicit. Inventory Source answers **what inventory may be promised**; Fulfillment Node answers **where/by whom the physical fulfillment work is executed**;
- Selling Entity, Inventory Source and Fulfillment Node remain separate semantics unless an explicit business rule relates them;
- if node mapping, eligibility or responsibility is ambiguous, MPC surfaces configuration-required/exception state rather than silently choosing a default warehouse/location;
- the first deployment may use one internally operated Fulfillment Node and a trivial Fulfillment Scope without hardcoding single-node or internal-only semantics into the canonical model.

Exact node identities, addresses, native ERP/WMS/provider mappings, node capability taxonomy, routing policy/algorithm, capacity/SLA logic and execution contracts belong to later stages.

#### D0.7e.4 — Cost Observation and Cost Basis

**Accepted by operator.**

`Cost Observation` and `Cost Basis` are canonical MPC economic semantics defined independently of any ERP-native cost taxonomy.

A **Cost Observation** is an attributable economic observation or derived conclusion about cost that preserves enough meaning to state what the value represents, its amount/currency, relevant temporal/business context and provenance.

A **Cost Basis** is the explicit economic semantic or governing selection policy that determines which cost meaning/evidence is appropriate for a particular pricing analysis, expected-profitability decision, transaction interpretation or realized-profitability calculation.

The purpose is to answer, in MPC language:

> **Which cost is economically appropriate for this decision/transaction, what does it mean, when/contextually where is it valid, and what evidence produced it?**

Product-level requirements:

- a cost used materially by MPC is not represented as an unexplained number; its semantic meaning and provenance must remain inspectable;
- Cost Basis is not defined as an ERP cost field/type/name. A participating adapter translates the required MPC cost semantic to authoritative ERP/other-source evidence when that mapping is supportable;
- one Cost Basis may require combining/transforming multiple native facts, and one native cost construct may satisfy different MPC uses only where its semantic meaning actually matches; no `1:1` mapping is assumed;
- if a requested Cost Basis cannot be satisfied reliably, MPC represents unsupported/missing/ambiguous cost evidence rather than choosing another native cost silently;
- expected profitability uses cost evidence appropriate to the decision context/time, while realized profitability uses evidence appropriate to the realized transaction/economic outcome. They are not required to be numerically identical and their variance may itself be decision-relevant;
- historical economic interpretation must not silently substitute a current product cost for the cost semantic relevant to the historical transaction;
- Cost Observation/Cost Basis do not absorb all economics into one opaque value. Marketplace fees, freight effects, taxes, discounts, subsidies, refunds/reversals and other material components remain separately attributable where they affect profitability;
- Cost Basis may be MPC-owned policy, externally governed rule or derived selection depending on the business authority. Its provenance/authority class must remain explicit;
- where Cost Basis is legitimately MPC-owned and configurable, later stages may apply the same inherited/default plus explicit override principles accepted for MPC policy, but D0 does not freeze a cost-basis catalog or every override scope;
- the first deployment may use a single accepted Cost Basis if that satisfies the real business need, without hardcoding a particular Sankhya cost variant as the universal MPC meaning of cost.

Exact native cost mappings, source fields/APIs, cost taxonomy, calculation methods, effective-time lookup, currency handling, allowed fallback semantics, policy scopes and profitability formulas belong to later stages.

#### Next D0.7e decision

Define the next ERP-independent business dimension: **Order Execution Scope**.

The question is not “which Sankhya TOP/company/fields create the order?”. The question is:

> **What must MPC represent so a marketplace sale can be translated into the correct business-system order operation, with the right execution context and governing rules, without making an ERP-native order type/configuration part of the MPC canonical domain?**

D0 must decide whether a canonical order-execution scope/intent semantic is required and what business distinctions it must preserve. Exact ERP order types, operation codes, document models, field mappings, API calls and routing mechanics belong to later stages.

A later ERP-independent dimension still to test is fiscal/invoicing scope. D0 must not assume order execution, Selling Entity, Inventory Source, Fulfillment Node or invoicing semantics are identical merely because the first ERP may bind them together.

Other findings discovered during D0.7 are classified as MUST DECIDE NOW, SHOULD DECIDE NOW or CAN DEFER SAFELY. D0 closes only when no material Product 1.0 semantic is being left for implementation to invent.

---

## 11. Resume contract for a fresh session

A fresh session must read `AGENTS.md`, `docs/engineering/rebaseline/README.md`, the canonical engineering method, `ARCHITECTURE.md`, the ADR registry, and then this D0 artifact.

It should conclude:

- D0 is open and not yet accepted as a whole;
- D0.1–D0.6 are operator-approved product decisions;
- D0.7a is accepted: essential cancellation/return/refund operations remain inside the controlled sale lifecycle without expanding Product 1.0 into general CRM/SAC;
- D0.7b is accepted: shipment/delivery remains visible through terminal outcome and material delivery exceptions become MPC operational work without turning MPC into a TMS;
- D0.7c is accepted: marketplace availability is maintained automatically from governing authoritative stock/rules/policies when sufficiently known, while uncertainty/failure becomes explicit work and MPC does not become physical-stock authority;
- D0.7d is accepted: one MPC organization may own/control one or more marketplace installations, even if Product 1.0 is initially proven with one Mercado Livre account;
- D0.7e is accepted: canonical MPC business semantics are defined from marketplace-operating needs before ERP-specific mapping; ERP integration semantically translates rather than dictating the MPC domain;
- D0.7e.1 is accepted: `Selling Entity` is a canonical MPC business dimension for the acting business/legal/fiscal entity, independent from ERP company identifiers, marketplace installation, Inventory Source, Fulfillment Node and Cost Basis unless explicit business rules relate them;
- D0.7e.2 is accepted: `Inventory Source` identifies business-recognized inventory sources, `Inventory Scope` explicitly governs which sources may contribute to an offer, and Sellable Availability is derived from eligible sources plus authoritative facts/rules rather than copied from an external stock field;
- D0.7e.2a records the Product 1.0 requirement for configurable MPC-owned availability-allocation policy, including percentage-style allocation such as `70%`, while the exact policy catalog/scopes/arithmetic are intentionally deferred;
- D0.7e.3 is accepted: `Fulfillment Node` identifies a recognized physical-fulfillment execution point/capability, `Fulfillment Scope` governs which nodes are eligible, and fulfillment semantics remain distinct from inventory promise, Selling Entity and provider-native locations;
- D0.7e.4 is accepted: `Cost Observation` preserves cost meaning/value/time-context/provenance, `Cost Basis` explicitly selects the economically appropriate cost semantic, and ERP-native cost types cannot become silent MPC defaults/fallbacks;
- historical/realized economics must not silently use current cost as a substitute, and cost remains distinct from other economic components such as marketplace fees, freight and taxes;
- organization identity must not be collapsed into marketplace seller-account or Selling Entity identity;
- Sankhya-native concepts such as `CODEMP`, `CODLOC`, TOPs and cost variants are integration evidence, not automatically canonical MPC concepts;
- MPC uses **OWN / ORCHESTRATE / OBSERVE-DERIVE** and owns the marketplace operating model without taking ownership of external facts merely because it consumes or causes them;
- Sankhya API availability for writes is evidence to carry into D4, not a D0 transport decision;
- business rules/policies may be MPC-owned, externally governed or derived, and later stages must preserve that provenance;
- no D1+ target architecture may be invented yet;
- the exact next work is **D0.7e — define ERP-independent Order Execution Scope semantics needed by MPC**.
