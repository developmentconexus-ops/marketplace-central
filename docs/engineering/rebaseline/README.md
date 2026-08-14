# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status / router after `AGENTS.md`  
> **Current phase:** DOCUMENTARY / GOVERNANCE AUTHORITY CLEANUP — IN PROGRESS  
> **Implementation:** BLOCKED until D9 is accepted  
> **Evidence baseline:** `main@de1dc88bcef5a6ed5515378e7c646682c0bc15d2`  
> **Last updated:** 2026-08-14

## 1. Why this file exists

This is the one place a fresh session uses to determine:

- where the program is;
- what is authoritative now;
- what is prohibited now;
- what the exact next action is;
- when the current phase is finished.

There is deliberately no parallel roadmap, wiki progress page, permanent session-handoff tree or active legacy implementation plan. Git history is the archive.

## 2. Current phase: documentary / governance authority cleanup

The current PR is **#41**. Its purpose is not to design the target software. Its purpose is to remove competing legacy authority so the future design process starts from one unambiguous control plane.

### In scope now

- remove or retarget references to retired documents and authority trees;
- remove stale milestone/feature/wave ownership (`M-xx`, `F-xx`, old missions/plans) from active authority;
- align governance registries, gates, workflows and scripts where they consume retired documentary authority;
- prevent auxiliary tools from recreating retired documentation trees;
- keep `AGENTS.md`, this file, `ARCHITECTURE.md`, the ADR registry and `contracts/governance/` mutually coherent;
- verify the cleanup without weakening gates or raising ratchet baselines merely to make them pass;
- prove that a fresh session finds one authority path and one exact next action.

### Explicitly out of scope now

Do **not** redesign, refactor, migrate, choose or delete legacy product/runtime code merely because it looks old.

In particular, the cleanup does not decide:

- `modules` versus `contexts`;
- domain boundaries;
- identity or tenant model;
- database target schema;
- API target contract;
- frontend target topology;
- runtime/process/job topology;
- Mercado Livre/Sankhya target adapter design;
- auth/permissions;
- transaction/event/outbox architecture.

Existing code, schema, OpenAPI, tests and runtime are **evidence about the present system**, not authority for the target system.

A product/runtime finding discovered during cleanup is recorded as evidence and adjudicated only when the corresponding D-stage needs it.

### Narrow exception: documentary consumers inside tooling

A tool may be changed during cleanup only when the change is necessary to stop it from consuming or recreating retired documentary authority.

Example: `apps/server_core/cmd/mlprobe` currently references retired `docs/design/...` material and writes evidence into `docs/design/evidence/ml-api`. Retargeting that documentary output/reference is cleanup; redesigning the probe's marketplace behavior is not.

## 3. Cleanup completion criteria

The documentary cleanup is DONE only when all of the following hold:

1. no retired document competes as architecture/program authority;
2. no active governance registry points to deleted authority as current authority;
3. gates/workflows/scripts do not depend conceptually on retired documentary authority;
4. auxiliary tools no longer recreate retired documentation trees;
5. current governance is self-contained;
6. verification is green without weakening controls or inflating ratchets;
7. no material dead reference remains in the active authority path;
8. a fresh session can identify one authority path and one exact next action without chat history.

When these conditions pass, the cleanup stops. We do not extend it into a general codebase audit.

## 4. What happens after cleanup

After cleanup, the target-design program begins with **D0 — Product / System Definition**.

The governing sequence is:

```text
DOCUMENTARY / GOVERNANCE AUTHORITY CLEANUP
  ↓
D0 — Product / System Definition
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

Each D-stage is a decision process, not a pretext for auditing every legacy file first.

For each material decision:

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

Use `docs/engineering/standards/root-cause-global-maximum-method.md` for non-trivial decisions.

The current repository is inspected **on demand** according to the decision being made. Existing work is useful evidence, but it does not get a vote merely because it already exists.

For each question distinguish:

- **MUST DECIDE NOW** — implementation would otherwise invent semantics;
- **SHOULD DECIDE NOW** — materially affects downstream architecture;
- **CAN DEFER SAFELY** — operational/configuration detail that can be decided later without creating architectural ambiguity.

A later D-stage may explicitly reopen an earlier decision when new evidence creates a material contradiction. Silent contradiction is not allowed.

## 6. D0–D9 questions

| Stage | Core question |
|---|---|
| **D0 — Product / System Definition** | What exactly are we building, for whom, what problem does it solve, what belongs inside/outside, and what is Product 1.0? |
| **D1 — Domains / Boundaries** | Which capabilities/domains exist and who owns each responsibility/state? |
| **D2 — Identity / Tenant / Data Ownership** | Who/what are the canonical identities and which authority owns each class of data? |
| **D3 — Communication / Events** | How do components coordinate state and communicate without duplicate authority? |
| **D4 — External Integrations** | How do Mercado Livre, Sankhya/Oracle and future external systems enter the product? |
| **D5 — API** | What contracts expose the accepted capabilities and semantics? |
| **D6 — Frontend** | How does the application represent workflows and consume capabilities without duplicating business authority? |
| **D7 — Runtime / Jobs / Transactions** | How are execution, scheduling, concurrency, retries, transactions and failure recovery handled? |
| **D8 — Golden Flows** | Do the important end-to-end flows remain coherent through success, partial failure, retry and reconciliation? |
| **D9 — Adversarial Architecture Review** | Where can the accepted design contradict itself, overbuild, under-specify or fail under real constraints? |

Each stage consults legacy code/runtime/schema only to answer the questions actually in front of it.

## 7. Legacy disposition policy

Legacy product/runtime units are not classified for deletion during the current documentary cleanup.

During the relevant D-stage they may later be classified as, for example:

- KEEP;
- KEEP AS REFERENCE;
- REFACTOR;
- MIGRATE / MOVE;
- REPLACE;
- DELETE.

No classification is granted solely by age, directory name, or current reachability discovered incidentally during documentary cleanup.

## 8. Documentation authority

Read active authority in this order:

1. `AGENTS.md` — bootstrap, process, prohibitions;
2. **this file** — sole current status and exact next action;
3. `docs/engineering/standards/root-cause-global-maximum-method.md` — decision method;
4. `ARCHITECTURE.md` — stable product/platform constraints that have actually survived rebaseline authority review;
5. `docs/architecture/decisions/README.md` — ADR status registry;
6. accepted/current D-stage artifact(s), once design begins;
7. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md` — supporting evidence only;
8. code, OpenAPI, schemas, tests and runtime — current-state evidence.

Historical plans/specs/handoffs/wikis do not become target authority because they remain available in Git history.

## 9. Current supporting evidence

`docs/engineering/rebaseline/EVIDENCE-REGISTER.md` contains already-collected facts worth carrying forward. They are supporting evidence, not prerequisites that must all be exhaustively reproduced before D0 starts.

Any additional codebase measurement is performed when a D-stage decision requires it.

## 10. Exact next action

**Finish PR #41 documentary / governance authority cleanup. Do not start D0 design yet.**

Specifically:

1. reconcile remaining active consumers of retired documentary authority;
2. retarget tools that recreate retired documentation paths (including `mlprobe`);
3. remove/retarget stale governance references and ownership language;
4. run the repository's existing verification gates without weakening them;
5. perform the fresh-session authority test;
6. delete the temporary session handoff once canonical authority alone is sufficient;
7. mark documentary cleanup DONE.

Then stop cleaning and open **D0 — Product / System Definition** with the operator.

## 11. Fresh-session success test

After cleanup, a fresh session should be able to read `AGENTS.md` and this file and state correctly:

- we are not implementing product features yet;
- the cleanup has either a precise remaining action or is DONE;
- historical code/docs are evidence, not target authority;
- the next design stage is D0 Product / System Definition;
- implementation begins only after the design program is accepted.

If it cannot, authority cleanup is not finished.