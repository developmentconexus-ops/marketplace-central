# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D1 — DOMAINS / BOUNDARIES — OPEN / PENDING INDEPENDENT REVIEW**  
> **Implementation:** BLOCKED until D9 is accepted  
> **Last updated:** 2026-08-16

## 1. Authority path

A fresh session reads, in order:

1. `AGENTS.md`
2. this file
3. `docs/engineering/standards/root-cause-global-maximum-method.md`
4. `ARCHITECTURE.md`
5. `docs/architecture/decisions/README.md`
6. `docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md`
7. `docs/engineering/rebaseline/D1-DOMAINS-BOUNDARIES.md`
8. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
9. code, OpenAPI, schemas, tests and runtime only as current-state evidence when needed

This file alone answers **where the program is and what happens next**. Stable architecture belongs in `ARCHITECTURE.md`; accepted stage semantics belong in D-stage artifacts; Git history is the archive.

Do not reconstruct target authority from memory, legacy package shape, historical plans or stale docs.

## 2. Program state

```text
Documentary / governance cleanup — DONE
  ↓
D0 — Product / System Definition — CLOSED / ACCEPTED
  ↓
D1 — Domains / Boundaries — OPEN / PENDING INDEPENDENT REVIEW
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

## 3. Accepted baseline

### D0 — CLOSED

Product 1.0 is **Marketplace Operations + Commercial Intelligence**. MPC is the marketplace operations control plane: external systems retain authority for facts/processes inherently theirs while MPC owns the cross-system marketplace operating semantics needed to observe, understand, decide, execute, verify and reconcile.

D0 authority and non-goals are defined only in `D0-PRODUCT-SYSTEM-DEFINITION.md`.

### D1 — decisions complete; closure review pending

`D1-DOMAINS-BOUNDARIES.md` consolidates the operator-approved D1 decisions:

- 12 business boundaries with explicit ownership/non-ownership;
- `Mechanism ≠ Authority` as the cross-cutting boundary rule;
- policy/evidence/time/audit/execution-safety/integration capability treatment;
- semantic authority edges and forbidden boundary violations;
- semantic disposition of all legacy modules/current contexts;
- explicit D2–D7 defers and reopen triggers.

D1 is **not closed yet**. The current artifact is the candidate contract for independent review.

## 4. Engineering method

All material decisions follow:

`docs/engineering/standards/root-cause-global-maximum-method.md`

Do not replace that method with local shorthand. In particular:

- current implementation is evidence, not destiny;
- unknown is not a plausible default;
- solve root cause/defect class, not one symptom;
- preserve essential complexity and remove accidental complexity;
- YAGNI removes speculative capability, not correctness or justified future seams;
- define proof before implementation;
- accepted decisions reopen only on material new evidence.

## 5. What is prohibited now

Until D1 closes:

- do not begin D2–D9 target design prematurely;
- do not implement product architecture/features;
- do not choose schema/IDs/persistence ownership;
- do not choose events/outbox/sync-vs-async communication;
- do not choose provider/ERP transport contracts;
- do not choose HTTP/frontend/runtime topology;
- do not preserve or reject a legacy module merely because it exists;
- do not turn cross-cutting invariants into new contexts/frameworks without independent business authority;
- do not silently contradict `D1-DOMAINS-BOUNDARIES.md`.

Existing code/module/context shape remains current-state evidence only.

## 6. Exact next action

Complete the D1 closure pipeline in this order:

1. **Independent adversarial Fable review** of the complete `D1-DOMAINS-BOUNDARIES.md` contract.
2. Adjudicate every material finding using the repository method rather than by deference.
3. Run the final internal **Global Coherence + YAGNI / Overengineering / Future-Cost** review across corrected D1.
4. Obtain explicit operator approval of D1 as a whole.
5. Mark D1 **CLOSED / ACCEPTED** and make D2 the exact next stage.

Do **not** skip the independent review gate or perform the final global closure review before its findings are adjudicated.

## 7. Fresh-session success test

A fresh session should conclude that:

- D0 is closed/accepted;
- D1 decisions are consolidated but D1 remains open pending independent review;
- `D1-DOMAINS-BOUNDARIES.md` is the candidate D1 contract;
- 12 business boundaries do not imply 12 services/databases/processes;
- shared mechanisms do not acquire business authority by reuse;
- current modules/contexts are evidence, not target authority;
- D2–D9 and implementation remain blocked;
- the exact next action is the **independent Fable review of D1**.

If it cannot, the authority path is incomplete or contradictory.