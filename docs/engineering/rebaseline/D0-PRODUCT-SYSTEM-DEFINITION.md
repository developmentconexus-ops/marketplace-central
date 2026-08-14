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

External systems remain authorities for facts that inherently belong to them; the exact data/identity/write ownership matrix is deferred to D2/D4. MPC owns the cross-system operational control semantics: intent, policy, workflow, correlation, controlled execution, operational state, divergence, audit and reconciliation.

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
5. **Decision & Policy Control** — translate observations/recommendations into permitted, approval-required or prohibited actions according to company policy.
6. **Order-to-ERP Operations** — receive/understand marketplace orders and coordinate the corresponding business operations in Sankhya, including order creation and invoicing where authorized.
7. **Reconciliation & Exception Operations** — identify cases where an expected cross-system result cannot be proven or systems diverge, and make them actionable instead of silently assuming success.
8. **Realized Profitability** — determine realized sale economics and compare expected versus realized results so material variance can be explained.

Automation/human approval is a cross-cutting control-plane property, not a ninth product capability.

### D0.3a — Action authority model

**Accepted by operator.**

Product 1.0 supports:

- **human-controlled execution by default:** recommendation/analysis → operator approval → MPC executes → verifies/reconciles;
- **policy-driven automatic execution when explicitly authorized:** an accepted policy may approve a bounded action automatically;
- **human review for exceptions, uncertainty, low confidence or policy violations.**

A fully autonomous repricing system is not required as a launch gate. The product must, however, have the semantics necessary to support explicit bounded automation without bypassing policy, audit or reconciliation.

## 3. Product 1.0 non-goals currently safe to defer

The following are not required to define Product 1.0 as complete unless later D0 evidence proves otherwise:

- paid ads/media management;
- automated buyer Q&A or buyer chat;
- general CRM;
- broad reputation optimization;
- company-wide demand forecasting or purchasing;
- company-wide logistics/WMS replacement;
- marketplaces beyond what is needed to prove the first Mercado Livre operating loop;
- unrestricted AI/autonomous commercial decisions without explicit policy;
- a generic integration framework for many speculative providers.

These may become future MPC capabilities; deferral is not a permanent exclusion.

## 4. Operator-provided Sankhya evidence / constraint

This is **D0/D4 evidence**, not a provider-transport decision made in D0.

The operator reports an already-proven Sankhya application integration in another app using application credentials (`client id`, `client secret`, `X-Token`; no secret values are recorded here).

Operational writes such as order creation and invoicing can be performed through Sankhya's application API. The operator prefers this path for writes rather than direct Oracle database writes because it is a safer system-owned operation boundary.

The environment also has DB Explorer access for database inspection.

Consequences for later design:

- D0 may require MPC to create/invoice orders in Sankhya without assuming direct database writes.
- D4 must evaluate and ratify the exact Sankhya read/write capability contracts and transport boundaries.
- Existing binding Oracle-read constraints are not silently reopened here; the target read path and any need for direct Oracle access are adjudicated only by the stage that owns that decision.

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

## 6. Current open D0 decision

### D0.4 — Actors / operational users

**Partially accepted by operator; responsibility/authority boundaries remain open.**

D0 distinguishes human actors by the operational responsibility they carry. "Operator" is an umbrella description, not a single undifferentiated persona.

The actor classes accepted so far are:

1. **Marketplace Operations Operator** — performs day-to-day channel operations around products, listings, prices, marketplace state, orders, divergences and operational exceptions within permitted authority.
2. **Fulfillment / Dispatch Operator** — owns the physical fulfillment execution for marketplace orders: identifying work to prepare, separating and physically checking the correct items, progressing valid orders into invoicing, packing, preparing the shipment/dispatch handoff, reporting completion/problems and progressing the order through the physical dispatch workflow exposed by MPC.
3. **Commercial / Marketplace Manager** — governs or approves commercial decisions such as price/margin boundaries, higher-impact changes, exception decisions and bounded automation policies.
4. **Owner / Administrator / Policy Approver** — governs system-level configuration, integrations, users and exceptional/high-impact policy authority where that responsibility is not delegated elsewhere.

This actor model does **not** define JWT roles, permission tables or technical authorization implementation. Those are later-stage concerns.

The Fulfillment / Dispatch Operator makes an important product-boundary distinction: MPC Product 1.0 includes the marketplace-order fulfillment workflow needed to move a marketplace sale from ERP/order readiness through physical separation, invoicing, packing and dispatch handoff. This does **not** make MPC a company-wide WMS or logistics platform.

#### D0.4a — Fulfillment / Dispatch Operator authority

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

The Fulfillment / Dispatch Operator is therefore not limited to packing already-invoiced orders. Within the accepted fulfillment workflow, this actor has business authority to confirm that the physical order is correct and then request/trigger invoicing.

The target property is that **an order is not intentionally invoiced through the normal fulfillment path before the operator has physically confirmed that the correct items are available and separated**.

If the operator finds missing stock, wrong item, damaged material, quantity divergence or another physical inconsistency, the normal invoicing transition is blocked and the order becomes an operational exception requiring resolution rather than being invoiced optimistically.

D0 does not decide how this authority is represented in permissions, how the workflow state is persisted, or how the Sankhya call is implemented. Those belong to later stages.

#### D0.4b — Marketplace Operations Operator authority

**Accepted by operator.**

The Marketplace Operations Operator owns routine day-to-day marketplace control within company policy. At product/workflow level this actor may:

- prepare and correct product/channel readiness information and operational linkage where the underlying evidence is sufficient;
- create and publish a new marketplace listing when the product is fully ready and the action satisfies all current policies and required conditions;
- inspect and edit listings and their permitted operational state;
- analyze competitive position, price scenarios and expected profitability;
- execute price changes and other bounded marketplace actions when they remain inside accepted commercial policies;
- pause/reactivate or correct marketplace operational state when the action is policy-compliant;
- investigate and progress ordinary marketplace divergences/exceptions that are within the actor's operational authority.

The actor does **not** own the commercial policies themselves. It cannot redefine margin floors, price boundaries, approval thresholds or automation authority merely to make an action permissible.

The product decision rule is:

```text
fully ready + inside accepted policy
  → Marketplace Operations Operator may decide/execute

outside policy / insufficient evidence / higher-impact exception
  → Commercial / Marketplace Manager review or other explicit escalation
```

Initial listing creation is therefore **not approval-gated merely because it is the first publication**. If the product is fully ready and all governing policies are satisfied, the Marketplace Operations Operator can create/publish it directly through MPC. The exact readiness criteria, policy representation and authorization mechanism are later-stage decisions.

#### Next exact D0.4 decision

Define the responsibility/authority boundary of the **Commercial / Marketplace Manager**, especially what this actor may approve, which commercial policies it may define/change, and what remains reserved for Owner / Administrator / Policy Approver.

Do this at product/workflow level only; do not design permissions/auth implementation yet.

## 7. Resume contract for a fresh session

A fresh session must read `AGENTS.md`, `docs/engineering/rebaseline/README.md`, the canonical engineering method, `ARCHITECTURE.md`, the ADR registry, and then this D0 artifact.

It should conclude:

- D0 is open and not yet accepted as a whole;
- D0.1–D0.3a above are operator-approved decisions;
- D0.4 has four accepted actor classes: Marketplace Operations Operator, Fulfillment / Dispatch Operator, Commercial / Marketplace Manager, and Owner / Administrator / Policy Approver;
- D0.4a is accepted: the fulfillment operator physically separates/checks the order and, when valid, triggers invoicing through MPC before packing/dispatch; physical inconsistencies block normal invoicing and become exceptions;
- D0.4b is accepted: the marketplace operator controls routine marketplace operations within policy and may create/publish a new listing without separate commercial approval when the product is fully ready and policy-compliant;
- D0.4 remains open for the responsibility/approval/accountability boundaries of Commercial / Marketplace Manager and Owner / Administrator / Policy Approver;
- no D1+ target architecture may be invented yet;
- Sankhya API availability for writes is evidence to carry into D4, not a D0 transport decision;
- the exact next work is the **Commercial / Marketplace Manager** authority boundary inside D0.4.
