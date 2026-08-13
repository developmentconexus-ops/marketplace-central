# Marketplace Central — Architecture Rebaseline

> **Role:** canonical current-program entrypoint after `AGENTS.md`  
> **Program:** Technical Architecture Deep Dive  
> **Status:** D0 — Current State & Authority Baseline — IN PROGRESS  
> **Implementation:** BLOCKED until D9 is accepted  
> **Evidence baseline:** `main@de1dc88bcef5a6ed5515378e7c646682c0bc15d2`  
> **Working branch:** `docs/architecture-rebaseline`  
> **Tracking PR:** #40 (draft while design is under review)  
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

Git history is the archive. Deleted documents may be recovered only when a current stage explicitly needs historical evidence.

---

## 1. Product objective

Marketplace Central is an internal operations and intelligence system for marketplace commerce, initially Mercado Livre, backed by real Sankhya/Oracle operational facts.

The target system must support a trustworthy loop from internal product/stock/cost/tax truth through marketplace observations, linking, pricing/action decisions, safe external writes, orders, reconciliation and realized profitability.

The rebaseline exists because the current repository grew through multiple architectural eras and now contains duplicated authorities across modules, API shapes, database state, integration boundaries and frontend surfaces.

The objective is not to preserve that history. It is to establish one coherent target before the next large implementation wave.

---

## 2. Binding decisions already accepted

These decisions govern the deep dive unless a new material constraint causes an explicit ADR amendment.

### RB-01 — Plan the whole technical system before implementation

Architecture folders alone are insufficient. Context ownership, identity, data, communication, events, adapters, API, frontend, runtime and golden flows must be decided deeply enough that an implementation task does not need to invent architectural semantics.

### RB-02 — No product implementation before D9

The flow is:

```text
Architecture Rebaseline
  ↓
Technical System Design
  ↓
Context-by-context Design
  ↓
Data / API / Events / Adapters / Frontend Design
  ↓
Golden Flow Simulation
  ↓
Contradiction / Global-Maximum Review
  ↓
Implementation DAG
  ↓
Implementation Plan
  ↓
Codex / implementation
```

### RB-03 — No `old/`/archive source tree

Git history is the historical archive. Legacy source and documentation are deleted after any still-valid property is absorbed into the current authority.

### RB-04 — Hard cutover is allowed

There are no production users requiring backward compatibility. Routes, schemas, IDs, package APIs and frontend navigation may break when the target design requires it.

Compatibility is opt-in and needs a measured reason.

### RB-05 — Main must remain epistemically trustworthy

Hard cutover does not mean permanently red or ambiguous `main`.

Preferred landing unit:

```text
new authority
  → all consumers/writers cut over
  → old authority unreachable
  → old code/schema/route deleted
  → proof green
  → one authority on main
```

### RB-06 — Do not copy current architecture into the target by default

Current code, schema, OpenAPI, packages and old ADRs are evidence of the present system and migration cost. They are not sufficient arguments that the target should have the same boundaries.

Use the inversion test in the Global-Maximum Method.

### RB-07 — No broker without need

A Postgres transactional outbox plus bounded consumers is the default async candidate. Kafka/NATS/RabbitMQ require a concrete need that Postgres cannot adequately satisfy.

### RB-08 — Clean DB baseline is an allowed candidate

Before deciding, classify every current state class as re-derivable, re-authorizable, configuration, human decision or non-rederivable audit/history.

### RB-09 — Old 13-context protocol is a hypothesis

The previous proposal (`account`, `catalog`, `listings`, `linking`, `inventory`, `costing`, `tax`, `pricing`, `orders`, `reconciliation`, `profitability`, `intelligence`, `changecontrol`) contains useful measured ideas but is not accepted wholesale. D1 re-adjudicates the context set.

### RB-10 — Issues/old docs do not define architecture

Issues are evidence and traceability. Historical plans/specs/handoffs do not become a roadmap merely because they remain in Git history.

---

## 3. Program gates

| Gate | Question | Required output | Status |
|---|---|---|---|
| **D0** | What actually exists and which authorities are live? | current-state topology, import/runtime/API/DB/event/frontend census, document authority, legacy inventory, contradictions | **IN PROGRESS** |
| **D1** | Which business contexts should exist? | context admission/rejection, responsibilities, commands, queries, contracts, ports, legacy destination | PENDING |
| **D2** | What are the identities and data authorities? | canonical IDs, external refs, table/schema ownership, writer/readers, temporal/knowledge semantics, RLS, reset decision | PENDING |
| **D3** | How do internal components communicate? | sync-call matrix, event catalog, projection rules, transaction/outbox ownership, ordering/idempotency | PENDING |
| **D4** | How do external systems integrate? | Mercado Livre and Sankhya capability catalogs, auth/account ownership, pagination, rate limit, retry, freshness, write/readback semantics | PENDING |
| **D5** | What is the external HTTP contract? | operation inventory keep/redesign/delete, target OpenAPI, error model, pagination/idempotency/concurrency, generator choice | PENDING |
| **D6** | How does the frontend map to product capabilities? | route/screen map, API consumer map, feature topology, query/cache/error/loading/empty states | PENDING |
| **D7** | How does the system run? | process topology, scheduler, leases, workers, transaction spine, outbox dispatcher, readiness/shutdown/observability | PENDING |
| **D8** | Does the design survive real end-to-end scenarios? | fully traced golden flows with ownership, transaction, failure, retry, freshness and proof at each hop | PENDING |
| **D9** | Is this the global maximum for current constraints? | adversarial contradiction review, local-vs-global findings, YAGNI pass, proof plan, unresolved blockers = 0 | PENDING |

No gate is complete because its document exists. It completes only after its target questions are answered, contradictions are dispositioned and the operator accepts the result.

---

## 4. Required maps before implementation

The deep dive must produce these maps as real design outputs, not implementation-time guesses.

### 4.1 Context map

For every candidate context:

- business responsibility and language;
- state/lifecycle authority;
- commands and queries;
- published contracts;
- required ports;
- events produced/consumed;
- external adapters used;
- tables owned;
- API/frontend surfaces;
- legacy replaced;
- deletion gate.

### 4.2 Data ownership map

For every persistent object/table:

- context owner;
- source-of-truth classification;
- primary/unique identity;
- tenant boundary;
- writer authority;
- legal readers;
- temporal semantics;
- knowledge/unknown semantics;
- RLS/constraints;
- retention/rebuildability;
- idempotency key where relevant;
- legacy disposition.

### 4.3 Identity map

At minimum adjudicate:

- Tenant/User identity;
- marketplace channel account/install identity;
- ERP connection/instance identity;
- Product/SourceProduct identity;
- Listing/Variation external identity;
- Order/external order identity;
- linking identity/evidence.

External identifiers are not assumed to be canonical internal identities.

### 4.4 Internal communication map

Every cross-context relationship is deliberately one of:

- synchronous query/capability because the result is required now;
- asynchronous event because another component reacts to a fact that already happened;
- projection/read model because several authorities are combined for reading.

No cross-context SQL or accidental import is an unclassified communication mode.

### 4.5 Event catalog

For each event:

- producer/authority;
- event/version/aggregate identity;
- tenant/correlation/causation;
- transaction/outbox boundary;
- consumers;
- ordering key;
- dedupe/idempotency;
- retry/dead-letter handling if needed;
- schema evolution.

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
- raw payload retention policy;
- error translation;
- read-after-write/reconciliation.

### 4.7 HTTP/API map

For every current operation:

`current operation → owner → keep / redesign / delete`

For every target operation:

- operationId;
- method/path;
- owning context/use case;
- auth/tenant;
- request/response;
- problem/error model;
- pagination/filter/sort;
- idempotency/concurrency;
- cache semantics;
- frontend/external consumer.

### 4.8 Frontend map

For every screen/route:

`screen → query/mutation → API operation → owning context/projection`

Also decide feature/package topology, query keys, invalidation, loading, empty, stale, error and permission/account states.

### 4.9 Runtime map

Classify each executable/runtime capability as serving process, worker/background capability, operator command, migration/test tooling or temporary probe with deletion gate.

Decide scheduler, leases, process isolation, readiness, shutdown, transactions, outbox dispatch and observability only from measured needs.

### 4.10 Golden flows

Trace real flows through every boundary and failure mode. Minimum families:

1. Sankhya product → catalog authority → listing/linking/inventory/pricing decision;
2. Mercado Livre listing observation → mapping/linking → stock/price interpretation;
3. safe external price/stock write → ambiguous outcome → readback → verified/diverged;
4. Mercado Livre order → product/cost/tax/fees → reconciliation → realized profitability.

If a flow reaches “we decide this during implementation,” the design is not complete.

---

## 5. Stage working protocol

For each D-stage:

1. Read this file and the current stage document.
2. Measure the current system; distinguish known fact from hypothesis.
3. Use outside/official technical sources when the decision depends on current provider/library behavior.
4. Group findings by root cause/target property, not old issue number.
5. Propose alternatives and identify local maximum risks.
6. Record the chosen authority/boundary and proof.
7. Run a contradiction/self-review against all earlier accepted stages.
8. Present the stage for operator acceptance.
9. Update this status table and exact next action in the same documentation change.
10. Only then open the next D-stage.

A later stage may reopen an earlier one only on a material contradiction or changed constraint, and the status table must show that explicitly.

---

## 6. Current documentation authority

Read current docs as:

1. `AGENTS.md` — routing, process, current prohibitions.
2. **This file** — sole current program status/progress.
3. `docs/engineering/standards/root-cause-global-maximum-method.md` — decision method.
4. `ARCHITECTURE.md` — stable product-level constraints only.
5. `docs/architecture/decisions/README.md` — ADR registry; check status before reading an ADR as authority.
6. Current D-stage document(s).
7. `contracts/api/`, `contracts/governance/` and code — current-state evidence and current runtime contracts, not automatic target architecture.
8. Explicitly labeled supporting references.

There is intentionally no live `wiki/`, `.mnfs/`, `docs/superpowers/`, legacy `IMPLEMENTATION_PLAN.md` or root `EVIDENCE.md` after the documentation hygiene change. Those remain recoverable in Git history.

---

## 7. Current D0 facts already established

At `main@de1dc88b...`:

- `internal/modules/` contains **21** legacy module directories;
- `internal/contexts/` contains **2** new contexts: `catalog`, `listings`;
- `internal/adapters/` already separates `erp/` and `marketplace/`;
- `cmd/` has **7** entrypoint directories;
- frontend organization is mixed between `packages/feature-*` and `apps/web/src/routes/pages`;
- legacy redirects remain in `AppRouter.tsx`;
- OpenAPI is large and contract authority has historically been manually duplicated;
- migrations span several architectural eras;
- current governance still contains controls/ratchets scoped to `internal/modules`, so those controls are transitional evidence rather than proof of the target boundary.

Detailed established facts and remaining D0 census live in `D0-current-state-and-authority.md`.

---

## 8. Exact next action

**Complete D0 — Current State & Authority Baseline.**

The next session must not start D1 yet. It must finish the current-state census across:

1. package/import graph and SCCs;
2. runtime composition/reachability;
3. table/schema ownership and all writers/readers;
4. API/route/handler/SDK/frontend-consumer topology;
5. current event/outbox/job/scheduler topology;
6. Mercado Livre and Sankhya adapter reachability;
7. frontend feature/route/query topology;
8. recoverability classification of current database state;
9. surviving references to deleted documentary authorities;
10. root-cause grouping of the measured contradictions.

Then present D0 for operator acceptance.

Only after **D0 ACCEPTED** does D1 Context Adjudication begin.

---

## 9. Session restart rule

A fresh session should be able to say, after reading only `AGENTS.md`, this file and the current D-stage document:

- current phase;
- accepted decisions;
- prohibited work;
- current evidence baseline;
- exact unfinished measurements;
- exact next gate.

If it cannot, the documentation structure has regressed and must be fixed before continuing product design.