# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status / router after `AGENTS.md`  
> **Current phase:** D0 — PRODUCT / SYSTEM DEFINITION — OPEN; working, not yet accepted as a whole.  
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

## 2. Closed phase: documentary / governance authority cleanup

The cleanup landed in PR **#41**. Its purpose was not to design the target software. Its purpose was to remove competing legacy authority so the design process starts from one unambiguous control plane. It is closed; §3 records how each completion criterion was discharged.

The scope below is kept because it defines what the cleanup was allowed to touch, and therefore what it did **not** decide.

### In scope for that cleanup

- remove or retarget references to retired documents and authority trees;
- remove stale milestone/feature/wave ownership (`M-xx`, `F-xx`, old missions/plans) from active authority;
- align governance registries, gates, workflows and scripts where they consume retired documentary authority;
- prevent auxiliary tools from recreating retired documentation trees;
- keep `AGENTS.md`, this file, `ARCHITECTURE.md`, the ADR registry and `contracts/governance/` mutually coherent;
- verify the cleanup without weakening gates or raising ratchet baselines merely to make them pass;
- prove that a fresh session finds one authority path and one exact next action.

### Explicitly out of scope, then and now

Do **not** redesign, refactor, migrate, choose or delete legacy product/runtime code merely because it looks old.

In particular, the cleanup did not decide:

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

A tool could be changed during cleanup only when the change was necessary to stop it from consuming or recreating retired documentary authority.

Worked example, now discharged: `apps/server_core/cmd/mlprobe` referenced retired `docs/design/...` material and wrote evidence into `docs/design/evidence/ml-api`. It now writes to `/workspace/output/ml-api` and cites nothing under `docs/design/`. Retargeting that documentary output/reference was cleanup; redesigning the probe's marketplace behavior would not have been.

## 3. Cleanup completion criteria, and how each was discharged

The documentary cleanup was DONE only when all of the following held. Each row records the criterion and the evidence that closed it, so a later session can re-check the claim instead of trusting it.

| # | Criterion | Discharged by |
|---|---|---|
| 1 | no retired document competes as architecture/program authority | the retired trees (`.mnfs/`, `docs/superpowers/`, `docs/design/`, `docs/HARNESS-PROFILE.md`, `docs/engineering/repo-audit-2026-08-07/`) are deleted from the repository; `docs/README.md` carries the removal record. A checkout that predates the removal can still hold gitignored leftovers under those paths on disk — they are untracked local residue, never authority, and `git ls-files` is the check that settles it |
| 2 | no active governance registry points to deleted authority as current authority | governance exceptions carry `re_adjudicate_in: D<N>` instead of milestone/feature `removal_owner`; `scripts/tests/governance-contracts.tests.ps1` asserts that field and fails if either side is absent |
| 3 | gates/workflows/scripts do not depend conceptually on retired documentary authority | every `HARNESS-PROFILE` / `GATE-TOPOLOGY` / `docs/superpowers` citation in `scripts/`, `.github/`, `eslint.config.mjs`, `.golangci.yml`, `vitest.config.ts`, `contracts/gate/baselines.json` and `deploy/` was replaced by the rule itself or retargeted to the owning D-stage |
| 4 | auxiliary tools no longer recreate retired documentation trees | `cmd/mlprobe` writes to `/workspace/output/ml-api`; `scripts/harness/pack-measure.sh`, whose only subject was measuring an evidence pack inside `.mnfs/`, is deleted and had no invoker |
| 5 | current governance is self-contained | the authority chain in §8 resolves entirely to files in this tree |
| 6 | verification is green without weakening controls or inflating ratchets | `npm run gate:full` — 17 lanes, 0 failed. No baseline in `contracts/gate/baselines.json` was raised, no lane disabled, no test skipped, no scanner exemption added |
| 7 | no material dead reference remains in the active authority path | swept; what remains are quoted-inline provenance citations in historical ADRs and `_citations/`, which the ADR registry explicitly retains, and product-code comments recording where a fact came from |
| 8 | a fresh session can identify one authority path and one exact next action without chat history | §11, run cold against `AGENTS.md` alone |

The cleanup stops here. It is not extended into a general codebase audit.

## 4. Current design program

After cleanup, the target-design program begins with **D0 — Product / System Definition**. D0 is now open. Its current working artifact is `docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md`.

The governing sequence is:

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

Legacy product/runtime units are not classified for deletion merely because D0 is open.

During the relevant D-stage they may later be classified as, for example:

- KEEP;
- KEEP AS REFERENCE;
- REFACTOR;
- MIGRATE / MOVE;
- REPLACE;
- DELETE.

No classification is granted solely by age, directory name, or current reachability discovered incidentally.

## 8. Documentation authority

Read active authority in this order:

1. `AGENTS.md` — bootstrap, process, prohibitions;
2. **this file** — sole current status and exact next action;
3. `docs/engineering/standards/root-cause-global-maximum-method.md` — decision method;
4. `ARCHITECTURE.md` — stable product/platform constraints that have actually survived rebaseline authority review;
5. `docs/architecture/decisions/README.md` — ADR status registry;
6. `docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md` — current D0 artifact while D0 is active;
7. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md` — supporting evidence only;
8. code, OpenAPI, schemas, tests and runtime — current-state evidence.

Historical plans/specs/handoffs/wikis do not become target authority because they remain available in Git history.

## 9. Current supporting evidence

`docs/engineering/rebaseline/EVIDENCE-REGISTER.md` contains already-collected facts worth carrying forward. They are supporting evidence, not prerequisites that must all be exhaustively reproduced before D0 continues.

Any additional codebase measurement is performed when a D-stage decision requires it.

## 10. Exact next action

**Continue D0 with the operator: D0.7e — define ERP-independent Fiscal / Invoicing Semantics needed by MPC.**

Accepted D0.1–D0.6 decisions plus D0.7a essential post-sale, D0.7b shipment/delivery lifecycle, D0.7c automatic marketplace availability control, D0.7d `Organization 1 → N Marketplace Installations`, D0.7e's ERP-agnostic semantic-translation principle, D0.7e.1 `Selling Entity`, D0.7e.2 `Inventory Source / Inventory Scope`, D0.7e.2a availability-allocation policy requirement, D0.7e.3 `Fulfillment Node / Fulfillment Scope`, D0.7e.4 `Cost Observation / Cost Basis`, and D0.7e.5 `Business Order Intent` are recorded in `docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md`.

`Business Order Intent` is the MPC-owned semantic intent to materialize a marketplace sale in the participating business/ERP system while preserving material business context and source provenance. ERP-native order types/TOPs/document codes remain integration/native semantics; the ERP owns the native order record, while MPC owns the intent, correlation and materialization workflow. Missing/ambiguous mapping becomes explicit configuration/exception state, and ambiguous write outcomes are not blindly retried.

D0.7e.5 intentionally **does not introduce a separate canonical `Order Execution Scope`**. Material order execution context is carried by Business Order Intent plus explicit governing policies/integration mappings; a new scope concept would require independent business evidence later.

The availability-policy requirement remains intentionally scoped at D0 to the existence of configurable MPC-owned allocation policy. Exact policy catalog, scope hierarchy, override precedence details and arithmetic remain explicitly deferred for later adjudication rather than forgotten.

D0.7e continues to prohibit designing the MPC canonical domain by copying Sankhya-native constructs such as `CODEMP`, `CODLOC`, TOPs, cost variants or other ERP-specific structures. First define the business semantic MPC needs; later D2/D4 map Sankhya or another ERP into that semantic.

Nothing else is authorized. In particular: do not start product implementation, do not begin D1–D9 before D0 is accepted, and do not reopen the documentary cleanup.

D0 remains open until its product/system definition is complete, internally coherent, adversarially reviewed with the operator, and explicitly accepted as a whole.

## 11. Fresh-session success test

A fresh session should be able to read `AGENTS.md`, this file and the current D0 artifact and state correctly:

- documentary/governance cleanup is DONE;
- D0 Product / System Definition is OPEN and not yet accepted as a whole;
- D0.1–D0.6 and D0.7a–D0.7e.5 recorded in the D0 artifact are operator-approved decisions/principles/requirements;
- essential cancellation/return/refund operations remain inside the controlled sale lifecycle without expanding Product 1.0 into general CRM/SAC;
- shipment/delivery remains visible through a terminal outcome without turning MPC into a TMS;
- marketplace availability is automatically maintained from governing authoritative stock/rules/policies when sufficiently known; uncertainty/failure becomes explicit work and MPC does not become physical-stock authority;
- one MPC organization may control one or more marketplace installations, even if the first deployment uses one Mercado Livre seller account; organization identity is not marketplace-account identity;
- canonical MPC business semantics come from marketplace-operating needs, not from Sankhya/another ERP ontology;
- ERP integration is semantic translation; unsupported/incomplete mappings become explicit rather than guessed equivalence;
- `Selling Entity` is the canonical MPC concept for the acting business/legal/fiscal entity when material, independent from ERP-native company identifiers and from inventory/fulfillment/cost dimensions unless explicit business rules relate them;
- `Inventory Source` and `Inventory Scope` are canonical MPC inventory semantics; stock outside the governing scope does not contribute to Sellable Availability merely because it exists;
- MPC-owned availability-allocation policy may intentionally reserve/expose only part of eligible stock, including percentage-style use cases such as `70%`; exact policy catalog/scopes/arithmetic remain for later adjudication;
- `Fulfillment Node` and `Fulfillment Scope` are canonical MPC fulfillment semantics; node eligibility/responsibility is explicit and is not inferred from an ERP warehouse or Inventory Source;
- `Cost Observation` and `Cost Basis` are canonical MPC economic semantics; cost meaning/time-context/provenance must be explicit, ERP cost variants are integration evidence only, and unsupported/ambiguous cost does not silently fall back;
- `Business Order Intent` is canonical MPC intent/correlation semantics; ERP-native order types/TOPs stay behind the semantic integration boundary, and no separate `Order Execution Scope` is currently justified;
- ambiguous order materialization cannot be blindly retried or silently treated as failure/success;
- historical/realized economics must not silently use current cost as a substitute, and cost remains distinct from marketplace fees, freight, taxes and other economic components;
- MPC owns the marketplace operating model while external systems retain authority for facts/processes that inherently belong to them;
- rules/policies may be MPC-owned, externally governed or derived; MPC must preserve that provenance rather than silently duplicating authority;
- Product 1.0 requires its claimed normal operational path to be executable through MPC rather than relying on hidden routine manual system hopping;
- historical code/docs are evidence, not target authority;
- implementation remains blocked until D9 is accepted;
- the exact next work is **D0.7e — define ERP-independent Fiscal / Invoicing Semantics needed by MPC**.

If it cannot, the current authority path is incomplete or contradictory.