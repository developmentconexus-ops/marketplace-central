# Marketplace Central — Architecture Rebaseline

> **Role:** canonical current-program entrypoint after `AGENTS.md`  
> **Program:** Technical Architecture Deep Dive  
> **Status:** D0 — Current State & Authority Baseline — IN PROGRESS  
> **Implementation:** **BLOCKED until D9 is accepted**  
> **Evidence baseline:** `main@de1dc88bcef5a6ed5515378e7c646682c0bc15d2`  
> **Last updated:** 2026-08-13

## Why this file exists

This is the one place a fresh session uses to answer:

- where the program is;
- what has been accepted;
- what remains open;
- what is forbidden right now;
- which document to read next;
- the exact next action.

There is deliberately no parallel roadmap, wiki progress page, dated handoff tree or active legacy implementation plan.

Git history is the archive. Deleted documents are recovered only when a current stage explicitly needs historical evidence.

---

## 1. Product objective

Marketplace Central is an internal operations and intelligence system for marketplace commerce, initially Mercado Livre, backed by real Sankhya/Oracle operational facts.

The target system must support a trustworthy loop from internal product/stock/cost/tax truth through marketplace observations, linking, pricing/action decisions, safe external writes, orders, reconciliation and realized profitability.

The rebaseline exists because the repository grew through multiple architectural eras and accumulated competing authorities across modules, API shapes, database state, integration boundaries, frontend surfaces and documentation.

The objective is not to preserve that history. It is to establish one coherent technical system before the next implementation wave.

---

## 2. Binding rebaseline decisions

### RB-01 — Plan the technical system before product implementation

Architecture folders alone are insufficient. Context ownership, identity, persistence, communication, events, external adapters, API, frontend, runtime and golden flows must be decided deeply enough that an implementation task does not invent architectural semantics.

### RB-02 — No product implementation before D9

The governing flow is:

```text
D0 current state / authority
  ↓
D1 context adjudication
  ↓
D2 identity / data ownership
  ↓
D3 communication / events
  ↓
D4 external integrations
  ↓
D5 HTTP API
  ↓
D6 frontend
  ↓
D7 runtime / scheduler / transactions / outbox
  ↓
D8 golden-flow simulation
  ↓
D9 adversarial global-maximum review
  ↓
implementation DAG
  ↓
implementation plan
  ↓
implementation
```

### RB-03 — No `old/` or archive source tree

Git history is the historical archive. Legacy source and documentation are removed from the active tree after any still-valid property/evidence is absorbed into a current authority.

### RB-04 — Hard cutover is allowed

There are no production users requiring backward compatibility. Routes, schemas, IDs, package APIs and frontend navigation may break when the accepted target design requires it.

Compatibility is opt-in and requires a measured consumer/constraint.

### RB-05 — `main` must remain epistemically trustworthy

Hard cutover does not mean an indefinitely red or ambiguous `main`.

Preferred landing unit:

```text
new authority
  → all consumers/writers cut over
  → old authority unreachable
  → old code/schema/route removed
  → proof green
  → one authority on main
```

### RB-06 — Current structure is evidence, not target proof

Current code, schema, OpenAPI, packages, gates and old ADRs explain the present system and migration cost. They are not sufficient arguments that the target should preserve their boundaries.

Use the inversion test from the Root-Cause / Global-Maximum method.

### RB-07 — No broker without a measured need

A PostgreSQL transactional outbox plus bounded consumers is the default async candidate. Kafka, NATS, RabbitMQ or another broker requires a concrete need that Postgres cannot adequately satisfy.

### RB-08 — A clean database baseline is an allowed candidate

Before deciding, D2 classifies every current state class as re-derivable, re-authorizable, configuration, human decision or non-rederivable audit/history.

### RB-09 — The old 13-context proposal is a hypothesis

The previous proposal (`account`, `catalog`, `listings`, `linking`, `inventory`, `costing`, `tax`, `pricing`, `orders`, `reconciliation`, `profitability`, `intelligence`, `changecontrol`) contains useful measured ideas but is not accepted wholesale. D1 re-adjudicates the context set from domain semantics and ownership.

### RB-10 — Issues and historical documents do not define target architecture

Issues are evidence/traceability. Historical plans, specs and handoffs do not become a roadmap merely because Git retains them.

---

## 3. Program gates

| Gate | Question | Required output | Status |
|---|---|---|---|
| **D0** | What actually exists and which authorities are live? | current-state topology, import/runtime/API/DB/event/frontend census, document authority, legacy inventory, contradictions | **IN PROGRESS** |
| **D1** | Which business contexts should exist? | context admission/rejection, responsibilities, commands, queries, contracts, ports, legacy destination | PENDING |
| **D2** | What are the identities and data authorities? | canonical IDs, external refs, table/schema ownership, writers/readers, temporal/knowledge semantics, RLS, reset decision | PENDING |
| **D3** | How do internal components communicate? | sync-call matrix, event catalog, projection rules, transaction/outbox ownership, ordering/idempotency | PENDING |
| **D4** | How do external systems integrate? | Mercado Livre and Sankhya capability catalogs, auth/account ownership, pagination, rate limit, retry, freshness, write/readback semantics | PENDING |
| **D5** | What is the external HTTP contract? | operation inventory keep/redesign/delete, target OpenAPI, error model, pagination/idempotency/concurrency, generation/runtime-validation decision | PENDING |
| **D6** | How does the frontend map to product capabilities? | route/screen/API map, feature topology, query/cache/invalidation/error/loading/empty states | PENDING |
| **D7** | How does the system run? | process topology, scheduler, leases, workers, transaction spine, outbox dispatcher, readiness/shutdown/observability | PENDING |
| **D8** | Does the design survive real end-to-end scenarios? | fully traced golden flows with ownership, transaction, failure, retry, freshness and proof at each hop | PENDING |
| **D9** | Is this the global maximum for current constraints? | adversarial contradiction review, local-vs-global findings, YAGNI pass, implementation dependency DAG/proof matrix, material blockers = 0 | PENDING |

No gate is complete because its document exists. It completes only after its target questions are answered, contradictions are dispositioned and the operator accepts the written result.

---

## 4. Required technical maps before implementation

The deep dive must produce these as accepted design outputs, not implementation-time guesses.

### 4.1 Context map

For every candidate context:

- business responsibility/language;
- state/lifecycle authority;
- commands and queries;
- published contracts;
- required ports;
- events produced/consumed;
- external adapters used;
- tables/state owned;
- API/frontend surfaces;
- legacy replaced;
- deletion gate.

### 4.2 Data ownership map

For every persistent object/table:

- context owner;
- source-of-truth classification;
- primary/unique identity;
- tenant boundary;
- one writer authority;
- legal readers;
- temporal semantics;
- unknown/estimated/not-applicable semantics;
- exact monetary representation where applicable;
- RLS/constraints;
- retention/rebuildability;
- idempotency/concurrency key where relevant;
- legacy disposition.

### 4.3 Identity map

At minimum adjudicate:

- tenant/human identity;
- marketplace channel account / installation identity;
- ERP connection / instance identity;
- product / source-product identity;
- listing / variation internal and external identity;
- order / external-order identity;
- linking/decision/evidence identity.

External identifiers are not assumed to be canonical internal identities.

### 4.4 Internal communication map

Every cross-context relationship is deliberately one of:

- synchronous capability/query because the result is required now;
- asynchronous event because another component reacts to a committed fact;
- projection/read model because several authorities are combined for reading.

Cross-context SQL or accidental imports are not acceptable unclassified communication modes.

### 4.5 Event catalog

For each admitted event define:

- producer/authority;
- event type/version and entity/aggregate identity;
- tenant/correlation/causation;
- transaction/outbox boundary;
- payload contract;
- consumers;
- ordering key when needed;
- dedupe/idempotency;
- retry/dead-letter/escalation behavior where needed;
- schema/version evolution.

Events are not used merely to avoid a normal synchronous call.

### 4.6 External capability map

For Mercado Livre and Sankhya/Oracle, record per operation/capability:

- external operation/query;
- read/write;
- owning context port;
- auth/account scope;
- pagination/cursor;
- rate limits;
- freshness/completeness semantics;
- retry rules;
- timeout/ambiguous-outcome semantics;
- idempotency;
- raw payload/evidence retention;
- error translation;
- read-after-write/reconciliation.

### 4.7 HTTP/API map

For every current operation:

`current operation → owner → keep / redesign / delete`

For every target operation define:

- operationId;
- method/path;
- owning context/use case;
- auth/tenant;
- request/response;
- problem/error model;
- pagination/filter/sort;
- idempotency/concurrency;
- cache/freshness semantics;
- frontend/external consumer.

D5 must choose and prove the actual OpenAPI generation/runtime-validation approach against this repository.

### 4.8 Frontend map

For every screen/route:

`screen → query/mutation → API operation → owning context/projection`

Also decide feature/package topology, query keys, invalidation, loading, empty, stale, error and permission/account states.

The frontend must not become a second business-policy authority.

### 4.9 Runtime map

Classify each executable/runtime capability as serving process, worker/background capability, operator command, migration/test tooling or temporary probe with deletion gate.

Decide scheduler, leases/fencing, process isolation, readiness, shutdown, transaction ownership, outbox dispatch and observability from measured needs rather than enterprise convention.

### 4.10 Golden flows

Trace real flows through every boundary and failure mode. Minimum families:

1. Sankhya product → catalog authority → listing/linking/inventory/pricing decision;
2. Mercado Livre listing observation → mapping/linking → stock/price interpretation;
3. safe external price/stock write → ambiguous outcome → readback → verified/diverged;
4. Mercado Livre order → product/cost/tax/fees → reconciliation → realized profitability.

If a flow reaches “decide during implementation”, the design is not complete.

---

## 5. Stage working protocol

For each D-stage:

1. Read this file and the current stage document.
2. Measure the current system; distinguish known fact from hypothesis.
3. Use current official/external technical sources when a decision depends on provider/library behavior.
4. Group findings by root cause/target property, not historical issue number.
5. Propose credible alternatives and identify local-maximum risks.
6. Record the chosen authority/boundary and proof.
7. Run a contradiction/self-review against all earlier accepted stages.
8. Present the stage for operator acceptance.
9. Update this status table and exact next action in the same documentation change.
10. Only then open the next D-stage.

A later stage may reopen an earlier one only on a material contradiction or changed constraint, and this status table must show it explicitly.

---

## 6. Documentation authority

Read active docs in this order:

1. `AGENTS.md` — routing/process/current prohibitions.
2. **This file** — sole current program status/progress/next action.
3. `docs/engineering/standards/root-cause-global-maximum-method.md` — decision method.
4. `ARCHITECTURE.md` — stable product-level constraints.
5. `docs/architecture/decisions/README.md` — ADR current status; check registry before treating an old ADR as target authority.
6. Current D-stage document(s).
7. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md` — supporting evidence only.
8. `contracts/api/`, `contracts/governance/`, code, migrations and tests — current-state/runtime evidence, not automatic target architecture.
9. Explicitly labeled supporting operational references.

There is intentionally no active `wiki/`, `.mnfs/`, `.harness/`, `.superpowers/`, `docs/superpowers/`, legacy `IMPLEMENTATION_PLAN.md` or session-handoff tree after the D0 documentation hygiene change.

---

## 7. D0 facts already established

At `main@de1dc88b...`:

- `internal/modules/` contains **21** legacy module directories;
- `internal/contexts/` contains **2** new contexts: `catalog`, `listings`;
- `internal/adapters/` already separates `erp/` and `marketplace/`;
- `cmd/` has **7** entrypoint directories;
- frontend organization is mixed between `packages/feature-*` and `apps/web/src/routes/pages`;
- legacy redirects remain in `AppRouter.tsx`;
- OpenAPI is large and contract authority has historically been manually duplicated;
- migrations span several architectural eras;
- current governance contains controls/ratchets scoped to `internal/modules`, so those controls are transitional current-state evidence rather than proof of the target boundary.

Detailed established facts and remaining census work live in `D0-current-state-and-authority.md`.

---

## 8. Exact next action

**Complete D0 — Current State & Authority Baseline. Do not start D1 yet.**

Finish the current-state census across:

1. package/import graph, fan-in/out and SCCs;
2. runtime composition/reachability and entrypoint classification;
3. table/schema ownership and every productive writer/reader;
4. API/route/handler/SDK/frontend-consumer topology;
5. current event/outbox/job/scheduler topology;
6. Mercado Livre and Sankhya adapter reachability/protocol surfaces;
7. frontend feature/route/query topology;
8. recoverability classification of current database state;
9. surviving references to retired documentary authorities;
10. root-cause grouping of measured contradictions.

Then present D0 for operator acceptance.

Only after **D0 ACCEPTED** does D1 Context Adjudication begin.

---

## 9. Session restart rule

A fresh session should be able to say, after reading only `AGENTS.md`, this file and the active D-stage document:

- current phase;
- accepted decisions;
- prohibited work;
- current evidence baseline;
- exact unfinished measurements;
- exact next gate.

If it cannot, the documentation topology has regressed and must be fixed before continuing product design.