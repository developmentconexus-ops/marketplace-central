# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status / router after `AGENTS.md`  
> **Current phase:** D0 — PRODUCT / SYSTEM DEFINITION — OPEN; bounded closure review completed, awaiting explicit whole-stage acceptance  
> **Implementation:** BLOCKED until D9 is accepted  
> **Evidence baseline:** `main@de1dc88bcef5a6ed5515378e7c646682c0bc15d2`  
> **Last updated:** 2026-08-15

## 1. Why this file exists

This is the one place a fresh session uses to determine where the program is, what is authoritative now, what is prohibited now, what the exact next action is, and when the current phase is finished.

There is deliberately no parallel roadmap, wiki progress page, permanent session-handoff tree or active legacy implementation plan. Git history is the archive.

## 2. Documentary / governance cleanup — DONE

PR **#41** removed/retargeted retired authority, stopped tooling from recreating it, made the authority chain self-contained and closed `npm run gate:full` green without weakening controls/ratchets.

That cleanup did **not** decide target domains, identity/tenant model, database schema, API/frontend/runtime topology, provider/ERP adapter design, auth/permissions or transaction/event/outbox architecture. Existing code, schema, OpenAPI, tests and runtime remain current-state evidence, not target authority.

## 3. Current design program

```text
DOCUMENTARY / GOVERNANCE AUTHORITY CLEANUP — DONE
  ↓
D0 — Product / System Definition — OPEN / READY FOR OPERATOR ACCEPTANCE
  ↓
D1 — Domains / Boundaries
  ↓
D2 — Identity / Tenant / Data Ownership
  ↓
D3 — Communication / Events
  ↓
D4 — External Integrations
  ↓
D5 — API
  ↓
D6 — Frontend
  ↓
D7 — Runtime / Jobs / Transactions
  ↓
D8 — Golden Flows
  ↓
D9 — Adversarial Architecture Review
  ↓
Implementation DAG / Plan
  ↓
Implementation
```

Product implementation remains blocked until D9 is accepted. D1 is not authorized until D0 receives explicit whole-stage operator acceptance.

## 4. D-stage decision method

Each material decision follows:

```text
needed evidence
  → alternatives
  → mature patterns / current external facts
  → trade-offs
  → recommendation
  → operator discussion
  → explicit decision
  → recorded contract
  → implications for later stages
```

Use `docs/engineering/standards/root-cause-global-maximum-method.md`. Classify questions as MUST DECIDE NOW, SHOULD DECIDE NOW or CAN DEFER SAFELY. Later stages may explicitly reopen earlier decisions when new evidence creates a material contradiction; silent contradiction is prohibited.

Global Maximum is **not** permission for infinite review or maximum abstraction. D0 now has an explicit stop rule: no new microtopic is opened unless a new real/reachable finding changes Product 1.0 meaning/authority/boundary, would otherwise force an implementer to invent business semantics, and does not clearly belong to D1–D7.

## 5. D0–D9 core questions

| Stage | Core question |
|---|---|
| **D0** | What exactly are we building, for whom, what belongs inside/outside, and what is Product 1.0? |
| **D1** | Which capabilities/domains exist and who owns each responsibility/state? |
| **D2** | Who/what are canonical identities and which authority owns each class of data? |
| **D3** | How do components coordinate without duplicate authority? |
| **D4** | How do marketplaces, Sankhya/Oracle, payment/bank and future external systems enter the product? |
| **D5** | What contracts expose accepted capabilities/semantics? |
| **D6** | How does the frontend represent workflows without duplicating business authority? |
| **D7** | How are execution, scheduling, concurrency, retries, transactions and recovery handled? |
| **D8** | Do important end-to-end flows remain coherent through success/failure/retry/reconciliation? |
| **D9** | Where can the accepted design contradict itself, overbuild, under-specify or fail under real constraints? |

## 6. Documentation authority

Read active authority in this order:

1. `AGENTS.md`;
2. **this file**;
3. `docs/engineering/standards/root-cause-global-maximum-method.md`;
4. `ARCHITECTURE.md`;
5. `docs/architecture/decisions/README.md`;
6. `docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md`;
7. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md` — supporting evidence only;
8. code, OpenAPI, schemas, tests and runtime — current-state evidence.

Historical plans/specs/handoffs/wikis do not become target authority merely because they remain in Git history.

## 7. Current supporting evidence

`EVIDENCE-REGISTER.md` contains provider/competitor evidence for Mercado Livre, Amazon, Magalu, Casas Bahia, Leroy/Mirakl, Shopee, MadeiraMadeira, ANYMARKET, Magis5 and Bling, plus Sankhya and repository evidence. It also records current evidence for asynchronous/partial provider outcomes, price approval/safety controls and duplicate/out-of-order event delivery.

These sources inform decisions but do not create target abstractions by themselves. Additional measurement is performed only when a D-stage decision requires it.

## 8. D0 bounded closure review — COMPLETED, no product blocker found

The closure review was adversarial and subtractive rather than another open-ended discovery round.

### 8.1 Lifecycle / product-boundary coverage

No material Product 1.0 lifecycle responsibility remains ownerless or semantically undefined at D0 altitude:

- product/channel readiness and listing operation;
- availability from explicit eligible inventory + governing facts/rules/policy;
- competitive/pricing/profitability intelligence;
- controlled decision/policy and external action;
- Business Order Intent and business-system materialization;
- Invoicing Intent and readiness-gated fiscal materialization;
- fulfillment/provider-requirement closure/dispatch;
- shipment/delivery terminal observation;
- essential cancellation/return/refund consequences;
- marketplace settlement, cash evidence where available and realized economics;
- reconciliation/exception work, deadlines/internal targets and operational ownership.

### 8.2 Stable-authority contradiction check

The accepted D0 direction is consistent with `ARCHITECTURE.md` and the binding ADR registry constraints: Mercado Livre first, Sankhya/Oracle external, provider protocol behind adapter boundaries, honest unknown/partial observations, controlled auditable/reconcilable writes, no blind retry, provider PII minimization, PostgreSQL for MPC-owned canonical state, Go backend business authority and React not a second policy authority.

One **documentary** contradiction was found during closure: `AGENTS.md` still stated that D0 had not opened. It was corrected on this branch. No product-semantic contradiction was found from that issue.

### 8.3 Explicit safe defers — not D0 blockers

The following are intentionally left to their owning stages rather than pulled into D0:

- **D1:** final contexts/boundaries and legacy-module disposition;
- **D2:** canonical identities, tenant/isolation model, persistence ownership and exact value/time/evidence representations;
- **D3:** event/synchronous/projection matrix, webhook/event semantics and outbox decisions;
- **D4:** exact Mercado Livre/Sankhya/payment/bank contracts, provider capability profiles, native fields/statuses/artifacts and source-specific completeness/freshness evidence;
- **D5:** HTTP/API shape, bulk contracts, version/precondition/idempotency contract surface;
- **D6:** portfolio/attention/work inbox, approvals, countdowns, notifications and other UI topology;
- **D7:** schedulers, workers, queues, polling, retries, locks/versioning, transactions, compensation, timing and deployment topology;
- **D8:** end-to-end proof of the accepted Product 1.0 flows;
- **D9:** final adversarial architecture review.

### 8.4 Explicit Product 1.0 non-goals / safe scope subtraction

Current Product 1.0 does not need campaign/discount-campaign authoring as a separate control surface, paid ads/media, buyer chat/Q&A, general CRM/SAC, company-wide WMS/TMS/reverse logistics, company-wide accounting/treasury, broad demand forecasting/purchasing, every future marketplace, a universal ERP framework, a generic marketplace-hub compatibility layer or unrestricted autonomous AI decisions.

Promotion/discount effects that materially affect price/order/economics remain observable/attributable facts even though campaign authoring itself is deferred.

### 8.5 Anti-overengineering result

D0.7j–D0.7n are **cross-cutting truth/action-safety properties**, not an instruction to create one context/module/framework per property. In particular, D1 must not infer standalone `Freshness`, `Coverage`, `DecisionLineage`, `BulkEngine` or `ApprovalRevalidation` domains merely because D0 names those invariants. Later stages should implement the smallest mechanisms inside the real owning capabilities/boundaries.

This is a binding interpretation of the D0/YAGNI intent unless later material evidence proves an independent business domain actually exists.

### 8.6 Closure conclusion

**No material D0 blocker remains after the bounded review.** No additional D0.7 microtopic is currently justified under the stop rule.

D0 is therefore **ready for explicit operator acceptance as a whole**, but remains OPEN until that acceptance is given.

## 9. Exact next action

**Operator decision: accept D0 — Product / System Definition as a whole, or identify a concrete material contradiction/blocker.**

Accepted decisions currently include D0.1–D0.6 and D0.7a–D0.7n. Important consequences include:

- MPC is Marketplace Operations + Commercial Intelligence (A+) and owns the marketplace operating model, not external native facts;
- canonical MPC semantics precede ERP/provider mapping;
- `Selling Entity`, `Inventory Source / Inventory Scope`, availability-allocation policy, `Fulfillment Node / Fulfillment Scope`, `Cost Observation / Cost Basis`, `Business Order Intent` and `Invoicing Intent` are distinct accepted semantics;
- Economic Evidence Chain = `Simulation → Order Economics → Marketplace Settlement → Cash Receipt evidence where available`;
- provider capabilities/authorities are context-sensitive and claimed MPC-controlled paths close provider requirements;
- marketplace hubs such as ANYMARKET/Magis5 are benchmark/competitive evidence, not target runtime dependencies;
- external obligations and MPC-owned internal targets remain distinct;
- actionable work has role ownership, assignment is not authorization, escalation is not notification and material work closes through evidence;
- evidence is freshness- and coverage-aware;
- historical material actions remain explainable from decision-time lineage;
- multi-target actions expose intended blast radius and partial/ambiguous results;
- approval is not permanent permission: consequential execution revalidates only materially governing conditions, without treating irrelevant drift as automatic invalidation.

Nothing else is authorized. In particular: **do not start D1 until the operator explicitly accepts D0 as a whole**, do not start product implementation, and do not reopen documentary cleanup.

## 10. Fresh-session success test

A fresh session should conclude correctly that:

- cleanup is DONE;
- D0 is OPEN but its bounded closure review is COMPLETE and found no material product blocker;
- D0.1–D0.6 and D0.7a–D0.7n are operator-approved individually;
- D0 as a whole still awaits explicit operator acceptance;
- implementation remains blocked until D9;
- no new D0 microtopic should be invented absent a new material finding satisfying the stop rule;
- D0 truth/action-safety invariants are cross-cutting properties, not automatic future contexts/modules;
- current code/docs remain evidence, not target authority;
- the exact next action is **operator acceptance of D0 as a whole or identification of a concrete material blocker**.

If it cannot, the current authority path is incomplete or contradictory.
