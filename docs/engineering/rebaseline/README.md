# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status / router after `AGENTS.md`  
> **Current phase:** D0 — PRODUCT / SYSTEM DEFINITION — CLOSED / ACCEPTED AS A WHOLE; D1 — DOMAINS / BOUNDARIES is the exact next stage  
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
D0 — Product / System Definition — CLOSED / ACCEPTED
  ↓
D1 — Domains / Boundaries — NEXT, NOT YET OPENED
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

Product implementation remains blocked until D9 is accepted. D0 acceptance authorizes opening D1; it does not pre-authorize any particular D1 domain/context layout.

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

Use `docs/engineering/standards/root-cause-global-maximum-method.md`. Classify questions as MUST DECIDE NOW, SHOULD DECIDE NOW or CAN DEFER SAFELY. Later stages may explicitly reopen an earlier decision when new evidence creates a material contradiction; silent contradiction is prohibited.

Global Maximum is **not** permission for infinite review or maximum abstraction. YAGNI removes speculative capability/accidental complexity, not correctness or required invariants.

## 5. D0–D9 core questions

| Stage | Core question |
|---|---|
| **D0** | What exactly are we building, for whom, what belongs inside/outside, and what is Product 1.0? — **ANSWERED / CLOSED** |
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
6. `docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md` — accepted D0 product/system authority;
7. the active D1 artifact once D1 is opened;
8. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md` — supporting evidence only;
9. code, OpenAPI, schemas, tests and runtime — current-state evidence.

Historical plans/specs/handoffs/wikis do not become target authority merely because they remain in Git history.

## 7. D0 final disposition — ACCEPTED AS A WHOLE

The operator accepted D0 after:

- iterative D0.1–D0.7n adjudication;
- broad current provider/competitor research;
- an explicit YAGNI/overengineering review;
- a bounded subtractive closure review;
- an independent adversarial Fable review;
- adjudication of Fable findings under the repository method rather than by deference;
- a final cold review of the corrected authority path.

No material Product 1.0 semantic remains for an implementer to invent at D0 altitude.

### 7.1 Accepted product boundary

Product 1.0 is **Marketplace Operations + Commercial Intelligence**.

MPC owns the cross-system marketplace operating model/control semantics while external systems retain authority for facts/processes inherently theirs.

The accepted lifecycle covers:

- product/channel readiness and listing operation;
- availability from explicit eligible inventory + rules/policy;
- competitive/pricing/profitability intelligence;
- controlled decision/policy and external action;
- Business Order Intent and business-system materialization;
- Invoicing Intent and readiness-gated fiscal materialization;
- fulfillment/provider-requirement closure/dispatch;
- shipment/delivery terminal observation;
- essential cancellation/return/refund consequences;
- marketplace/payment settlement and realized economics;
- reconciliation/exception work, deadlines/internal targets and operational ownership.

### 7.2 Final independent-review amendments

The final Fable review produced small boundary clarifications, not an architectural rebaseline:

- **Composite offers/kits:** real component-dependent semantics may not be silently flattened; provider-native composite support is not a Product 1.0 launch gate unless the selected first Mercado Livre flow requires it.
- **Organization posture:** first launch proof uses one Organization; Multi-Organization/SaaS operation is not a launch gate; Organization remains a real identity and tenant-ready isolation remains binding.
- **Marketplace Installation health/reputation:** provider-authoritative health/reputation may feed portfolio attention as observed/derived evidence; reputation management/SAC remains outside.
- **Economic Evidence Chain:** core Product 1.0 economic lineage closes through L2 Marketplace/Payment Settlement. L3 Bank Cash Receipt remains a semantically distinct optional extension and bank integration is not a launch gate.
- **Self-containment:** `(A+)` was removed; buyer Q&A/chat scope was made unambiguous; authority-map combined notation was defined.

### 7.3 Cross-cutting truth/action-safety invariants

Accepted D0 truth/action-safety properties include:

- unknown is not plausible zero/default;
- freshness is use-sensitive;
- fresh is not complete;
- not observed is not does-not-exist;
- provider success/submission is not automatically convergence;
- ambiguous external writes are not blindly retried;
- historical material decisions remain explainable from decision-time evidence/policy/authority;
- multi-target actions preserve intended blast radius and partial/ambiguous outcomes;
- approval is not permanent permission; materially governing conditions remain valid enough at execution time;
- external obligations and MPC-owned internal targets are distinct;
- actionable work has role ownership; assignment is not authorization; material closure requires evidence.

These are **cross-cutting invariants, not automatic D1 contexts/modules/frameworks**. D1–D7 must implement the smallest mechanism inside real owning boundaries.

## 8. Explicit safe defers — not D0 blockers

The following remain intentionally owned by later stages:

- **D1:** final domains/contexts/boundaries and legacy-module disposition;
- **D2:** canonical identities, tenant/isolation model, persistence ownership and exact value/time/evidence representations;
- **D3:** event/synchronous/projection matrix, webhook/event semantics and outbox decisions;
- **D4:** exact Mercado Livre/Sankhya/payment/bank contracts, provider capability profiles, native fields/statuses/artifacts, composite/provider-mode support and source-specific completeness/freshness evidence;
- **D5:** HTTP/API shape, bulk contracts, version/precondition/idempotency contract surface;
- **D6:** portfolio/attention/work inbox, approvals, countdowns, notifications and other UI topology;
- **D7:** schedulers, workers, queues, polling, retries, locks/versioning, transactions, compensation, timing and deployment topology;
- **D8:** end-to-end proof of accepted Product 1.0 flows;
- **D9:** final adversarial architecture review.

Product 1.0 also intentionally does not require paid ads/media, buyer Q&A/chat, general CRM/SAC, campaign authoring, company-wide WMS/TMS/reverse logistics, company-wide accounting/treasury/bank reconciliation, broad demand forecasting/purchasing, Multi-Organization/SaaS operation, every future marketplace, a universal ERP framework, a marketplace-hub compatibility layer or unrestricted autonomous AI.

## 9. Exact next action

**Open D1 — Domains / Boundaries with the operator.**

D1 must start from the accepted D0 semantics and independently determine the smallest real business domains/boundaries. It must **not** infer a standalone context/module for every named D0 invariant.

In particular, do not begin D2–D9 target design prematurely and do not start product implementation. Existing code/module/context shape remains evidence, not target authority.

## 10. Fresh-session success test

A fresh session should conclude correctly that:

- cleanup is DONE;
- D0 is **CLOSED / ACCEPTED AS A WHOLE**;
- Product 1.0 is Marketplace Operations + Commercial Intelligence;
- the final Fable review was adjudicated and resulted in scope clarifications/subtractions, not new frameworks;
- one Organization is enough for the launch proof, but Organization remains a real identity;
- composite semantics are never silently flattened, but composite provider support is not a launch gate by default;
- Marketplace Installation health/reputation can feed portfolio attention without becoming reputation management;
- core economic lineage closes through L2 settlement; L3 bank cash evidence is optional;
- canonical MPC semantics precede ERP/provider mapping;
- provider capabilities/authorities are context-sensitive and claimed MPC-controlled paths close provider requirements;
- D0 truth/action-safety properties are cross-cutting invariants, not automatic future contexts;
- marketplace hubs remain benchmark/competitive evidence, not target runtime dependencies;
- implementation remains blocked until D9;
- current code/docs remain evidence, not target authority;
- the exact next work is **D1 — Domains / Boundaries**.

If it cannot, the current authority path is incomplete or contradictory.
