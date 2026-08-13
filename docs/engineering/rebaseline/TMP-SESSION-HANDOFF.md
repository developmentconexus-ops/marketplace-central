# TEMP — Marketplace Central Rebaseline Session Handoff

> **Temporary continuity document. Not a design authority.**
>
> This file exists only so a fresh session can continue the current rebaseline without depending on chat memory. If anything here conflicts with `AGENTS.md` or `docs/engineering/rebaseline/README.md`, those canonical files win.
>
> Delete this file once the documentation cleanup and D0 bootstrap are stable enough that a fresh session can continue from the canonical authorities alone.

## 1. Where we are now

Marketplace Central is **not in implementation**. We are in an **Architecture Rebaseline / Technical Architecture Deep Dive**.

The current clean review surface is **PR #41**, branch `docs/architecture-rebaseline-clean`.

PR #40 is superseded and must not be merged. It temporarily created a second documentation authority and is historical evidence only.

The active sequence is:

```text
D0 current state / authority
  ↓
D1 context adjudication
  ↓
D2 identity / data ownership
  ↓
D3 internal communication / events
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

**No product implementation is authorized before D9 is accepted.**

## 2. Why we are doing this

The repository grew through several architectural eras and accumulated legacy in multiple dimensions:

- business modules and new contexts coexist;
- contracts/API shapes have historically been duplicated manually;
- database state and migrations reflect several old models;
- frontend organization is mixed;
- external integrations and account/credential responsibilities are not yet fully adjudicated;
- documentation, roadmaps, handoffs and old workflow trees became conflicting authorities;
- some current verification rules still describe the old `internal/modules` world.

The goal is not to make the current shape cleaner. The goal is to determine the **best coherent target system**, map the entire technical design deeply enough that implementation does not invent architecture on the fly, and then hard-cut from legacy to the accepted target.

## 3. Binding engineering method

Use the root-cause / global-maximum method for every non-trivial decision.

Core rule:

> Simplify code, never simplify correctness. Find the root cause, identify the target property, test whether the proposed fix is only a local maximum, and prefer the structure that removes duplicate authorities or makes the defect class impossible at the strongest reasonable boundary.

YAGNI removes speculative capability, duplicate machinery and unnecessary abstraction. It does **not** remove required invariants, fail-closed behavior, exact semantics or proof.

Before choosing a design, name:

1. observed symptom/evidence;
2. root cause;
3. target property/invariant;
4. authority/owner;
5. boundary where the property should be enforced;
6. local-maximum alternatives;
7. global-maximum candidate;
8. proof that distinguishes broken from fixed.

Do not add guards or synchronizers when the underlying duplicate authority can be removed structurally.

## 4. Legacy policy already decided

### No `old/` source tree

Do **not** move legacy source into an `old/` folder.

Git history is the archive. An `old/` tree would remain importable/searchable/copyable and would recreate a second architecture inside the repository.

### Hard cutover is allowed

There are no production users requiring backward compatibility. It is acceptable later to:

- break routes;
- replace IDs;
- replace schemas;
- reset/rebuild the database if D2 proves the state is safely re-derivable/re-authorizable;
- delete legacy modules;
- delete old frontend routes/packages;
- remove obsolete compatibility paths.

Compatibility is opt-in and must have a measured consumer/constraint.

### Main should still remain trustworthy

Hard cutover does not mean leaving `main` indefinitely broken.

Preferred landing unit:

```text
new authority
  → writers/readers/consumers cut over
  → old authority unreachable
  → old code/schema/route deleted
  → proof green
  → one authority on main
```

## 5. Documentation cleanup comes before deeper design work

The repository has accumulated ambiguous and conflicting documentation. The current session deliberately stops before trying to clean everything because large cleanup attempts have already shown that they can create another competing authority.

The **next session must first finish the documentation-authority cleanup around PR #41**, not start product design or implementation.

Classify every surviving active document as one of:

- **CANONICAL KEEP** — current authority for a defined topic;
- **SUPPORTING EVIDENCE** — useful, explicitly non-authoritative;
- **ABSORB THEN DELETE** — contains a valid property/decision that must be moved to the current authority before deletion;
- **DELETE** — stale, duplicate, unrelated or conflicting with the new program;
- **HISTORICAL ONLY** — remove from active tree; use Git history if needed later.

Target documentation topology:

1. `AGENTS.md` — fresh-session bootstrap, routing, prohibitions;
2. `docs/engineering/rebaseline/README.md` — **only** current program status, accepted decisions and exact next action;
3. `docs/engineering/standards/root-cause-global-maximum-method.md` — binding engineering method;
4. `ARCHITECTURE.md` — stable product-level constraints only;
5. `docs/architecture/decisions/README.md` + accepted ADRs — durable design decisions;
6. current D-stage document(s) — active design work;
7. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md` — supporting evidence only;
8. code/OpenAPI/schema/tests/runtime — evidence of what exists today, not automatic proof of target design.

There must not be multiple active roadmaps, session-handoff trees, wiki authorities, legacy implementation plans or duplicate architecture directories.

## 6. What is already measured about the codebase

At the current baseline used by PR #41:

- `apps/server_core/internal/modules/` has **21 legacy module directories**;
- `apps/server_core/internal/contexts/` has **2 new contexts**: `catalog` and `listings`;
- adapters are already separated under `internal/adapters/erp` and `internal/adapters/marketplace`;
- `cmd/` has **7 entrypoint directories**;
- frontend organization is mixed between `packages/feature-*` and route/page code inside the web app;
- legacy frontend redirects still exist;
- OpenAPI is large and contract authority has historically been manually duplicated;
- migrations span multiple architecture eras;
- some governance/ratchets are still scoped around `internal/modules` and are transitional evidence, not target architecture.

These facts are starting points only. D0 must freshly reproduce the complete census.

## 7. D0 must finish before D1

Do not adjudicate the final context set yet.

D0 must produce a measured current-state map covering:

1. package/import graph, fan-in/fan-out, SCCs and cross-boundary imports;
2. runtime composition and reachability of every entrypoint/job/command;
3. table/schema ownership, every productive writer and legal/actual readers;
4. API → route → handler → use case → SDK → frontend consumer topology;
5. event/outbox/job/scheduler topology;
6. Mercado Livre adapter reachability and protocol surfaces;
7. Sankhya/Oracle adapter reachability and protocol surfaces;
8. frontend route/feature/query/mutation topology;
9. database recoverability classification;
10. remaining references to retired documentary authorities;
11. measured contradictions grouped by root cause.

Only after D0 is presented and accepted does D1 begin.

## 8. Deep technical design required before implementation

Architecture folders are not enough. Before implementation planning, the program must decide the following in detail.

### D1 — Context map

For every candidate context define:

- business responsibility/language;
- state and lifecycle authority;
- commands/queries;
- published contracts;
- required ports;
- events produced/consumed;
- external adapters used;
- persistent state owned;
- API/frontend surfaces;
- legacy replaced;
- deletion gate.

The old proposed context list is a **hypothesis**, not a conclusion. `costing`, `profitability`, `intelligence`, `account` and other boundaries must be re-adjudicated from domain semantics.

### D2 — Identity + data ownership

Define canonical internal IDs versus external IDs for tenant, user, channel account/installation, ERP connection, product/source-product, listing/variation, order/external-order, link/decision/evidence, etc.

For every table/object define owner, writer, readers, tenant boundary, keys, constraints, RLS, temporal semantics, exact-money semantics, unknown/estimated/not-applicable semantics, retention/rebuildability, idempotency/concurrency and legacy disposition.

Also decide whether a clean DB reset/baseline is safer than carrying historical migrations.

### D3 — Internal communication + events

Every cross-context dependency must be deliberately one of:

- synchronous capability/query because the result is needed now;
- asynchronous event because another component reacts to a committed fact;
- projection/read model because several authorities are combined for reading.

No accidental cross-context SQL/import is an accepted communication mode.

For every event define producer, type/version, identity, tenant, correlation/causation, transaction/outbox boundary, payload, consumers, ordering, dedupe/idempotency, retry/escalation and evolution.

Default async candidate is PostgreSQL transactional outbox + bounded consumers. Do not add Kafka/NATS/RabbitMQ without a measured need.

### D4 — External integrations

Map Mercado Livre and Sankhya/Oracle operation by operation:

- external operation/query;
- read/write;
- owning context port;
- auth/account scope;
- pagination/cursor;
- rate limits;
- freshness/completeness;
- retry rules;
- timeout/ambiguous outcome;
- idempotency;
- raw payload/evidence retention;
- error translation;
- read-after-write/reconciliation.

Provider DTOs/protocol knowledge stay inside adapters. Contexts know ports/contracts, not Mercado Livre/Oracle internals.

### D5 — HTTP API

Create a complete current operation inventory:

`current operation → owner → KEEP / REDESIGN / DELETE`

For every target operation define operationId, method/path, owner/use case, auth/tenant, request/response, problem model, pagination/filter/sort, idempotency/concurrency, cache/freshness and frontend consumer.

OpenAPI should be evaluated as the single HTTP contract authority feeding generated Go/TypeScript surfaces; do not preserve manual duplication merely because it exists today.

### D6 — Frontend

Map every screen/route:

`screen → query/mutation → API operation → owning context/projection`

Decide feature topology, package boundaries, query keys, invalidation, loading/empty/stale/error/permission/account states.

Frontend must not become a second business-policy authority.

### D7 — Runtime

Decide process topology from measured needs: HTTP server, scheduler, workers, outbox dispatcher, leases/fencing, readiness, shutdown, transactions and observability.

Do not create extra services/processes merely because an enterprise system might have them.

### D8 — Golden flows

At minimum simulate end-to-end:

1. Sankhya product → canonical catalog → listing/linking/inventory/pricing decision;
2. Mercado Livre listing → ingest → linking → stock/price interpretation;
3. safe external price/stock write → timeout/unknown outcome → readback → verified/diverged;
4. Mercado Livre order → items/fees/freight → product/cost/tax → reconciliation → realized profitability.

Every hop must identify owner, ID, transaction, table/state, event/call, retry, failure behavior, freshness and proof.

If any flow reaches “decide during implementation”, design is incomplete.

### D9 — Adversarial global-maximum review

Before implementation planning:

- challenge every context boundary;
- challenge every authority;
- look for duplicate state/contract/vocabulary;
- distinguish essential versus accidental complexity;
- apply YAGNI;
- verify legacy deletion gates;
- verify no design only works because current code already looks that way;
- material contradictions = 0.

Only after D9:

```text
implementation DAG
→ implementation plan
→ Codex/implementation
```

## 9. Temporary decisions about the legacy source tree

Do not delete legacy source merely because it is under `internal/modules`.

During D0/D1 classify each legacy unit as:

- KEEP as target authority;
- REFACTOR IN PLACE;
- MOVE;
- MERGE INTO CONTEXT;
- REPLACE THEN DELETE;
- DELETE NOW because it is provably unreachable/unneeded.

Once a target owner and cutover plan are accepted, legacy deletion is preferred over compatibility wrappers.

## 10. What the next session should do

Start with:

1. read `AGENTS.md`;
2. read `docs/engineering/rebaseline/README.md`;
3. read this temporary handoff;
4. inspect PR #41, **not** superseded PR #40;
5. finish/document the documentation-authority cleanup without creating another authority tree;
6. verify a fresh session can locate one roadmap/status/next action;
7. continue D0 measurements listed in the canonical README;
8. update canonical D0/current-state evidence, not this temporary file, with actual measurements;
9. present D0 for operator acceptance;
10. only then begin D1.

### Explicitly prohibited now

- product implementation;
- implementation plan for the product rebaseline;
- creating new business contexts because the old 13-context proposal says so;
- moving source into `old/`;
- deleting reachable legacy source before its authority/consumers are mapped;
- preserving compatibility without a measured consumer;
- creating another roadmap/wiki/handoff authority;
- treating current code/OpenAPI/schema as automatic target design;
- allowing Codex/another agent to decide architecture during implementation.

## 11. Deletion condition for this TMP file

Delete `TMP-SESSION-HANDOFF.md` when all of the following are true:

1. PR #41 (or its accepted successor) has established a single documentation authority topology;
2. `AGENTS.md` points unambiguously to the canonical rebaseline README;
3. the canonical README contains current stage, accepted decisions, prohibitions and exact next action;
4. the active D-stage document contains all unfinished work needed for continuation;
5. a fresh session can continue without this temporary file.

Until then this file is only a continuity aid. It must never become a roadmap or design authority.
