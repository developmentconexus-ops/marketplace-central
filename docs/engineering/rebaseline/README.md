# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status / router after `AGENTS.md`  
> **Current phase:** D0 — PRODUCT / SYSTEM DEFINITION — OPEN; working, not yet accepted as a whole.  
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
D0 — Product / System Definition — OPEN
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

Product implementation remains blocked until D9 is accepted.

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

`EVIDENCE-REGISTER.md` contains provider/competitor evidence for Mercado Livre, Amazon, Magalu, Casas Bahia, Leroy/Mirakl, Shopee, MadeiraMadeira, ANYMARKET, Magis5 and Bling, plus Sankhya and repository evidence. These sources inform decisions but do not create target abstractions by themselves.

Additional measurement is performed only when a D-stage decision requires it.

## 8. Exact next action

**Continue D0 with the operator: D0.7 Product completeness / contradiction review — decision/action evidence lineage and reproducibility.**

Accepted D0.1–D0.6 plus D0.7a–D0.7k are recorded in `D0-PRODUCT-SYSTEM-DEFINITION.md`.

Important accepted decisions include:

- MPC is Marketplace Operations + Commercial Intelligence (A+) and owns the marketplace operating model, not external native facts;
- canonical MPC semantics precede ERP/provider mapping;
- `Selling Entity`, `Inventory Source / Inventory Scope`, MPC-owned availability-allocation policy, `Fulfillment Node / Fulfillment Scope`, `Cost Observation / Cost Basis`, `Business Order Intent` and `Invoicing Intent` are distinct accepted semantics;
- availability uses explicit eligible authoritative inventory + rules + MPC policy; routine known policy-valid synchronization is automatic; unknown is not zero;
- the Economic Evidence Chain is `Simulation → Order Economics → Marketplace Settlement → Cash Receipt evidence where available`; stages remain distinct, simulator calibration is evidence-driven and payout/order attribution is never fabricated;
- provider capabilities/authorities are context-sensitive and every flow claimed as MPC-controlled must close provider-required prerequisites/data/artifacts/readiness;
- ANYMARKET, Magis5 and similar hubs are benchmark/competitive evidence, not Product 1.0 runtime dependencies; direct MPC-owned provider boundaries are the target direction;
- External Operational Obligations and MPC-owned Internal Operational Targets are distinct; relative targets require explicit anchors and internal policy may tighten but never relax external obligations;
- actionable work has durable role ownership; individual assignment is distinct from authorization; escalation changes responsibility/attention/authority; material work closes through accepted resolution evidence;
- material evidence is freshness-aware and use-sensitive; stale/unknown-freshness evidence and failed acquisition cannot masquerade as current truth;
- observation coverage is explicit and scoped: **fresh is not complete**, **not observed is not does-not-exist**, callback receipt is not automatic completeness proof, and portfolio health/absence/reconciliation closure require sufficient coverage for the claimed universe.

The next contradiction is historical explainability. Even if evidence was fresh and sufficiently covered at decision time, current source state and MPC-owned policy may later change. D0 must decide whether material recommendations/approvals/actions retain enough decision-time facts, policy/rule provenance, authority/approval and uncertainty context to explain **why that historical action was permitted/recommended/executed**, instead of reconstructing history from current values.

This is a product auditability/explainability question. Snapshot/reference/hash/event-log/retention implementation belongs to D2/D3/D7.

Nothing else is authorized: do not start product implementation, do not begin D1–D9 before D0 is accepted, and do not reopen documentary cleanup.

D0 remains open until its product/system definition is complete, internally coherent, adversarially reviewed and explicitly accepted as a whole.

## 9. Fresh-session success test

A fresh session should conclude correctly that:

- cleanup is DONE;
- D0 is OPEN and not accepted as a whole;
- D0.1–D0.6 and D0.7a–D0.7k are operator-approved;
- implementation remains blocked until D9;
- Product 1.0 is Marketplace Operations + Commercial Intelligence (A+), not an ERP/marketplace/accounting/task-management replacement;
- canonical MPC semantics precede ERP/provider mapping and native provider/ERP constructs do not become canonical by existence alone;
- ambiguous external writes are not blindly retried;
- provider requirements, time obligations, work ownership, freshness and coverage all participate in truthful operational control;
- marketplace hubs remain benchmark/competitive evidence, not target runtime dependencies;
- current code/docs remain evidence, not target authority;
- the exact next work is **D0.7 Product completeness review — decision/action evidence lineage and reproducibility**.

If it cannot, the current authority path is incomplete or contradictory.
