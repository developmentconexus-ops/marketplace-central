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

**Next exact decision with the operator:** identify the real actors who operate or govern Marketplace Central and the decisions each actor must be able to make.

This is a product/workflow decision only. Do not define JWT roles, permission tables or auth implementation in D0.

Candidate actors are evidence to discuss, not accepted yet, for example:

- marketplace operator;
- commercial/marketplace manager;
- owner/administrator or policy approver.

D0.4 should establish whether these actor classes are sufficient and where approval/accountability differs between them.

## 7. Resume contract for a fresh session

A fresh session must read `AGENTS.md`, `docs/engineering/rebaseline/README.md`, the canonical engineering method, `ARCHITECTURE.md`, the ADR registry, and then this D0 artifact.

It should conclude:

- D0 is open and not yet accepted as a whole;
- D0.1–D0.3a above are operator-approved decisions;
- no D1+ target architecture may be invented yet;
- Sankhya API availability for writes is evidence to carry into D4, not a D0 transport decision;
- the exact next decision is **D0.4 — Actors / operational users**.
