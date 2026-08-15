# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status / router after `AGENTS.md`  
> **Current phase:** D0 — PRODUCT / SYSTEM DEFINITION — OPEN; working, not yet accepted as a whole.  
> **Implementation:** BLOCKED until D9 is accepted  
> **Evidence baseline:** `main@de1dc88bcef5a6ed5515378e7c646682c0bc15d2`  
> **Last updated:** 2026-08-15

## 1. Why this file exists

This is the one place a fresh session uses to determine where the program is, what is authoritative now, what is prohibited now, what the exact next action is, and when the current phase is finished.

There is deliberately no parallel roadmap, wiki progress page, permanent session-handoff tree or active legacy implementation plan. Git history is the archive.

## 2. Closed phase: documentary / governance authority cleanup

The cleanup landed in PR **#41**. Its purpose was not to design the target software. Its purpose was to remove competing legacy authority so the design process starts from one unambiguous control plane. It is closed.

### Explicitly out of scope, then and now

Do **not** redesign, refactor, migrate, choose or delete legacy product/runtime code merely because it looks old.

The cleanup did not decide target domains, identity/tenant model, database schema, API/frontend/runtime topology, provider/ERP adapter design, auth/permissions or transaction/event/outbox architecture. Existing code, schema, OpenAPI, tests and runtime are **evidence about the present system**, not target authority.

## 3. Cleanup completion record

The documentary cleanup is DONE: retired authority trees were removed/retargeted, governance no longer points to them as current authority, auxiliary tools no longer recreate them, the authority chain is self-contained, `npm run gate:full` closed green without weakening controls/ratchets, and a fresh session can identify one authority path and one exact next action.

The cleanup stops here. It is not extended into a general codebase audit.

## 4. Current design program

D0 — Product / System Definition is open. Its working artifact is `docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md`.

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

## 5. D-stage decision method

Each material decision follows:

```text
needed evidence
  → alternatives
  → relevant mature patterns / external facts
  → trade-offs
  → recommendation
  → operator discussion
  → explicit decision
  → recorded contract / artifact
  → implications for later stages
```

Use `docs/engineering/standards/root-cause-global-maximum-method.md` for non-trivial decisions. Classify questions as MUST DECIDE NOW, SHOULD DECIDE NOW or CAN DEFER SAFELY. Current code/runtime/schema is inspected on demand. Later stages may explicitly reopen an earlier decision when new evidence creates a material contradiction; silent contradiction is not allowed.

## 6. D0–D9 questions

| Stage | Core question |
|---|---|
| **D0** | What exactly are we building, for whom, what belongs inside/outside, and what is Product 1.0? |
| **D1** | Which capabilities/domains exist and who owns each responsibility/state? |
| **D2** | Who/what are the canonical identities and which authority owns each class of data? |
| **D3** | How do components coordinate without duplicate authority? |
| **D4** | How do marketplaces, Sankhya/Oracle, payment/bank and future external systems enter the product? |
| **D5** | What contracts expose accepted capabilities/semantics? |
| **D6** | How does the frontend represent workflows without duplicating business authority? |
| **D7** | How are execution, scheduling, concurrency, retries, transactions and recovery handled? |
| **D8** | Do important end-to-end flows remain coherent through success/failure/retry/reconciliation? |
| **D9** | Where can the accepted design contradict itself, overbuild, under-specify or fail under real constraints? |

## 7. Legacy disposition policy

Legacy product/runtime units are not classified for deletion merely because D0 is open. During the relevant D-stage they may later be classified KEEP, KEEP AS REFERENCE, REFACTOR, MIGRATE/MOVE, REPLACE or DELETE. No classification is granted solely by age, directory name or incidental reachability.

## 8. Documentation authority

Read active authority in this order:

1. `AGENTS.md` — bootstrap, process, prohibitions;
2. **this file** — sole current status and exact next action;
3. `docs/engineering/standards/root-cause-global-maximum-method.md` — decision method;
4. `ARCHITECTURE.md` — stable constraints that survived authority review;
5. `docs/architecture/decisions/README.md` — ADR status registry;
6. `docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md` — active D0 artifact;
7. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md` — supporting evidence only;
8. code, OpenAPI, schemas, tests and runtime — current-state evidence.

Historical plans/specs/handoffs/wikis do not become target authority merely because they remain in Git history.

## 9. Current supporting evidence

`docs/engineering/rebaseline/EVIDENCE-REGISTER.md` contains already-collected facts worth carrying forward. Provider/competitor evidence includes Mercado Livre, Amazon, Magalu, Casas Bahia, Leroy/Mirakl, Shopee, MadeiraMadeira, ANYMARKET, Magis5 and Bling. These are evidence, not target provider abstractions.

Current time-obligation evidence includes provider-authoritative dispatch deadlines and the mature operating pattern of using stricter internal targets as safety margin without rewriting the external obligation.

Additional codebase/provider measurement is performed only when a D-stage decision requires it.

## 10. Exact next action

**Continue D0 with the operator: D0.7 Product completeness / contradiction review — operational evidence freshness / staleness.**

Accepted D0.1–D0.6 plus D0.7a–D0.7i are recorded in `docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md`.

Important accepted decisions now include:

- ERP/provider semantics are translated into MPC business semantics rather than dictating them;
- `Selling Entity`, `Inventory Source / Inventory Scope`, MPC-owned availability-allocation policy, `Fulfillment Node / Fulfillment Scope`, `Cost Observation / Cost Basis`, `Business Order Intent` and `Invoicing Intent` are distinct accepted semantics;
- the Economic Evidence Chain is `Simulation → Order Economics → Marketplace Settlement → Cash Receipt evidence where available`, with evidence-driven simulator calibration and no fabricated payout/order attribution;
- provider capabilities/authorities are context-sensitive and any claimed MPC-controlled path must close provider-required prerequisites/data/artifacts/readiness;
- ANYMARKET, Magis5 and similar hubs are benchmark/competitive evidence, not Product 1.0 runtime dependencies; target direction is MPC-owned direct provider boundaries;
- provider/external time obligations remain external authority, organization-owned Internal Operational Targets remain MPC policy, relative targets require an explicit time anchor, and internal targets may tighten but never relax external obligations;
- actionable operational work has a durable owning role/responsibility; individual assignment is optional/distinct, assignment never grants authority, automation failure may create human-owned work, escalation means increased/different responsibility/attention/authority, and material closure requires accepted resolution evidence rather than arbitrary dismissal.

The next contradiction is freshness. MPC can still produce a confident but wrong decision if the underlying marketplace/ERP/provider evidence is stale while presented as current. D0 must decide whether material evidence freshness/observation age must explicitly influence readiness, decision safety and portfolio attention rather than remain a hidden integration/cache detail.

Polling intervals, TTLs, cache design, source-specific cadence, webhook/polling mechanics and scheduler implementation belong to D2/D3/D4/D7.

Nothing else is authorized. In particular: do not start product implementation, do not begin D1–D9 before D0 is accepted, and do not reopen documentary cleanup.

D0 remains open until its product/system definition is complete, internally coherent, adversarially reviewed with the operator and explicitly accepted as a whole.

## 11. Fresh-session success test

A fresh session should conclude correctly that:

- documentary/governance cleanup is DONE;
- D0 is OPEN and not accepted as a whole;
- D0.1–D0.6 and D0.7a–D0.7i are operator-approved;
- Product 1.0 is Marketplace Operations + Commercial Intelligence (A+), not an ERP/marketplace/accounting/task-management replacement;
- marketplace availability uses explicit eligible inventory + rules + MPC policy; known policy-valid sync is automatic and unknown is not zero;
- organization, Marketplace Installation and Selling Entity identities do not collapse;
- canonical MPC semantics precede ERP/provider mapping;
- native ERP company/location/order/invoicing/cost constructs stay behind semantic integration boundaries unless independent business semantics justify otherwise;
- ambiguous order/invoicing writes are not blindly retried;
- the Economic Evidence Chain preserves separate simulation/order/settlement/cash evidence and payouts are not assumed 1:1 with orders;
- provider capabilities/authorities are context-sensitive and claimed MPC-controlled paths close provider requirements without hidden routine provider-UI work;
- native provider artifacts remain provider-native semantics;
- marketplace hubs remain benchmark/competitive evidence, not target runtime dependencies;
- external time obligations and MPC-owned Internal Operational Targets are separate authorities; relative targets require explicit anchors and internal policy cannot relax external obligations;
- time participates in attention, safety-margin and breach semantics;
- actionable work has explicit role responsibility; assignment is distinct from authorization; escalation changes required attention/responsibility/authority; material work closes through resolution evidence;
- implementation remains blocked until D9;
- the exact next work is **D0.7 Product completeness review — operational evidence freshness / staleness**.

If it cannot, the current authority path is incomplete or contradictory.
