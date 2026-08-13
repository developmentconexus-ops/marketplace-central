# Marketplace Central — Technical Architecture Rebaseline Program

> **Status:** APPROVED PROGRAM  
> **Date:** 2026-08-13  
> **Implementation:** prohibited until D9 acceptance  
> **Progress authority:** `docs/architecture/rebaseline/README.md`

## 1. Purpose

Marketplace Central grew from a legacy foundation into a mixed architecture. New contexts and external adapters are already appearing beside 21 legacy modules, a large historical schema/API surface, duplicated documentation, manual contract surfaces and verification instruments that were originally anchored to the old tree.

The rebaseline exists to avoid a second generation of local decisions being made during implementation.

The target is not “clean folders”. Before implementation, the system must have an explicit answer for every material boundary:

- who owns each business fact;
- who may write it;
- how it is identified;
- where it is persisted;
- which communication is synchronous vs asynchronous;
- which events exist and why;
- which external adapter owns provider/ERP protocol knowledge;
- which API operation exposes each use case;
- which frontend surface consumes it;
- which runtime process executes it;
- which legacy surface becomes deletable;
- which proof demonstrates that one authority replaced another.

If an implementation worker still needs to invent one of those answers while writing feature code, the design program is incomplete.

## 2. Approved rebaseline decisions

### RB-01 — Rebaseline before further major product expansion

Do not add the next context merely because an old plan says it is next. Architecture mapping/adjudication precedes implementation.

### RB-02 — No source-code cemetery

Do not move legacy code into `old/`, `legacy/` or equivalent. Git history is the archive. A source cemetery remains searchable, importable, compilable and copyable and therefore preserves a second architecture.

### RB-03 — Hard cutover is allowed

There are no production users requiring compatibility with current routes, package APIs, schema shapes or frontend aliases. Compatibility exists only when a measured external constraint requires it.

### RB-04 — Hard cutover still requires a trustworthy merge state

A branch may be temporarily broken while assembling an atomic cutover. The merged state must be mechanically understandable and verifiable.

Preferred shape:

```text
new authority
  -> migrate its consumers
  -> migrate persistence/contracts
  -> prove zero legal old consumers/writers
  -> delete old authority
  -> merge one coherent authority
```

Not:

```text
new authority + compatibility layer + deprecated authority + unspecified future cleanup
```

### RB-05 — Previous architecture documents are evidence, not automatic target authority

The old 13-context protocol, ADRs, wiki, plans, audits and domain drafts contain useful measurements and design ideas. Their structural conclusions must be re-adjudicated. Current-tree shape and migration cost are valid sequencing evidence, not proof that a target boundary is right.

### RB-06 — Issues/findings are evidence grouped by root cause

Do not let issue numbering become the architecture. Several issues caused by one authority defect belong to one root-cause program.

### RB-07 — API/frontend compatibility may reset

There are no external API clients whose backward compatibility currently constrains the redesign. Legacy aliases and hand-maintained compatibility may be removed.

### RB-08 — No broker without a measured need

PostgreSQL transactional outbox + consumers is the default candidate. Kafka, NATS, RabbitMQ or another broker requires a concrete isolation/throughput/delivery need that Postgres cannot satisfy adequately.

### RB-09 — Clean DB baseline is a legal candidate

A full schema reset may be preferable to carrying a migration chain for authorities we intentionally delete, but D2 must first classify every existing state as re-derivable, re-authorizable/configurable, human-authored or audit-critical.

### RB-10 — Technical detail precedes implementation planning

Architecture folder shape alone is insufficient. D1–D8 must define operational details and D9 must review them as one system before any implementation plan is written.

## 3. Current-state baseline to interrogate

At `main@de1dc88bcef5a6ed5515378e7c646682c0bc15d2`:

- `internal/modules/` contains 21 module directories;
- `internal/contexts/` contains `catalog` and `listings`;
- `internal/adapters/` already separates `erp` and `marketplace` boundaries;
- `cmd/` contains 7 executable entrypoints;
- frontend feature ownership is split between packages and app-local routes/pages;
- legacy route redirects are still present;
- OpenAPI, SDK/runtime and routing surfaces grew historically with manual authorities;
- the migration chain contains both legacy and new-context state;
- governance checks still intentionally describe parts of the current legacy state and are migrated only with the corresponding architecture cutover.

The current implementation is an input to the program, not the target specification.

## 4. Provisional top-level dependency direction

The exact context set is D1, but the intended structural families are approved:

```text
apps/server_core/internal/
  contexts/       business authorities / bounded contexts
  adapters/       external system translation and protocol
  kernel/         tiny shared value semantics only when truly universal
  platform/       technical mechanisms, no business policy
  composition/    assembly/wiring only
  views/          rebuildable cross-context projections when required
```

### Context mold candidate

```text
contexts/<name>/
  contracts/      what the context intentionally publishes
  port/           capabilities the context requires
  internal/
    domain/
    application/
    postgres/
  module.go       public facade that assembles private internals
```

Go `internal/` is preferred when it can make an illegal dependency impossible instead of merely linting it.

### Adapter direction

```text
context owns port
      ↑
adapter implements it
      ↓
external provider / ERP
```

Vendor wire DTOs, OAuth payloads, provider URL/error shapes, Oracle driver types and raw SQL semantics are not domain contracts.

## 5. D0 — Authority and current-state rebaseline

### Goal

Make a fresh session unable to accidentally follow an obsolete roadmap or doctrine.

### Required outputs

- one read order and truth hierarchy;
- one Root-Cause / Global-Maximum method;
- one status/progress document;
- legacy decision disposition;
- evidence register carrying useful facts without carrying stale conclusions;
- deletion of competing wiki/mission/plan/handoff authorities;
- explicit next gate.

### D0 does NOT decide

- final context count;
- final DB schema/reset;
- event catalog;
- API shape/generator;
- frontend package topology;
- process topology.

Those decisions are deliberately deferred to the technical gates below.

## 6. D1 — Context Adjudication

### Question

**Which business contexts actually exist, and which current modules are mechanisms/projections/adapters rather than contexts?**

### For every candidate context

Document:

- business capability/question owned;
- state with independent lifecycle;
- invariants;
- commands;
- queries;
- published contracts;
- required ports;
- external sources that inform it without owning its semantics;
- potential events;
- frontend/read surfaces;
- current tables/modules/routes that appear to belong to it;
- nearest neighbor and why merge/split is correct;
- legacy disposition/delete gate.

### Admission test

A context is not admitted because a folder/name already exists. It must own business meaning/state/lifecycle that would become ambiguous or coupled if placed elsewhere.

Candidates from the previous proposal (`account`, `catalog`, `listings`, `linking`, `inventory`, `costing`, `tax`, `pricing`, `orders`, `reconciliation`, `profitability`, `intelligence`, `changecontrol`) are hypotheses, not a required count.

Special scrutiny:

- whether `costing` is distinct from catalog/finance facts;
- whether `profitability` is an authority or a reproducible projection/calculation record;
- whether `intelligence` is a context or read model;
- whether `account` incorrectly combines tenant/human identity with external installations/credentials;
- whether `dashboard`, `internal_read`, `sync`, `mutations`, `integrations`, `marketplaces`, `market`, `sourcekind` are mechanisms or misplaced responsibilities.

### D1 output

A context map plus an explicit `KEEP / MOVE / MERGE / DISSOLVE / REPLACE-DELETE` disposition for every current module.

## 7. D2 — Identity and Data Ownership

### Question

**What identities and persisted facts exist, and who is the single authority/writer for each?**

### Identity model

Adjudicate at minimum:

- TenantID / UserID if humans/auth are in scope;
- ChannelAccountID vs provider seller/account reference;
- ERPConnectionID / SourceInstanceRef;
- ProductID vs CODPROD/EAN/REFFORN/source product keys;
- ListingID vs external item ID;
- VariationID vs external variation ID;
- OrderID vs external order ID;
- link/decision/change request identities.

Binding question: **internal identity is not automatically an external identifier.** If CODPROD is made canonical internally, D2 must prove why that remains valid with multiple ERP/source instances rather than inherit it from legacy code.

### Data map

For every target table/state record:

- context owner;
- truth type: authoritative / derived / projection / raw observation / audit / configuration;
- primary/unique keys;
- tenant isolation;
- FKs/constraints and rationale;
- one write authority;
- legal readers;
- temporal semantics (`source_updated_at`, `observed_at`, effective period where relevant);
- unknown/estimated/not-applicable semantics;
- exact monetary representation;
- idempotency/concurrency key;
- retention/rebuildability;
- event emitted, if any.

### DB reset decision

Inventory current data into:

1. re-derivable from Sankhya;
2. re-derivable from marketplace/provider;
3. configuration/re-authorizable credentials;
4. human-authored decisions;
5. audit/evidence requiring preservation.

Only then choose forward migration vs clean baseline.

## 8. D3 — Internal Communication and Events

### Question

**How do contexts/components communicate without hidden database coupling or event cargo culting?**

Every cross-context need must be exactly one of:

1. **synchronous request/query** — caller needs the answer to complete the current use case;
2. **asynchronous event** — a committed fact occurred and another capability may react later;
3. **projection/view** — user-facing read combines multiple authorities but owns no business command/state authority.

Create a matrix:

| Producer/owner | Consumer | Information/use | Sync/event/projection | Contract | Why |
|---|---|---|---|---|---|

### Event catalog

For every admitted event define:

- stable event type/version;
- producer (one authority);
- aggregate/entity key where applicable;
- tenant/correlation/causation fields;
- transaction/outbox boundary;
- payload contract;
- ordering key if needed;
- consumers;
- consumer idempotency/dedup;
- retry/dead-letter/escalation behavior;
- version compatibility policy.

No event exists merely to avoid a normal synchronous call.

## 9. D4 — External Integration Contracts

### Question

**Exactly how do Mercado Livre and Sankhya/Oracle capabilities enter and leave the system?**

Static provider docs from the legacy repository are not authority. D4 re-verifies material provider behavior against current official documentation and live/read-only evidence when required.

### Mercado Livre capability matrix

Map at minimum:

- authorization/install/account discovery;
- token refresh/revocation/health;
- listings/items/variations;
- stock;
- price;
- fees/commission/shipping;
- orders/payments/shipments/cancellations;
- questions/messages if still product scope;
- notifications/webhooks;
- reads vs writes;
- pagination/cursor semantics;
- rate limits/backoff;
- retryability;
- idempotency;
- timeout-after-send ambiguity;
- raw payload/evidence retention;
- read-after-write/convergence verification;
- owning context port for each capability.

### Sankhya/Oracle matrix

Map at minimum:

- canonical product/source identity;
- stock by company/location;
- costs and method/as-of semantics;
- fiscal/tax source facts;
- sales/order/document history used by MPC;
- pagination/read consistency;
- query timeout/pool behavior;
- source timestamps/effective periods;
- mapping of Oracle types/errors;
- which context owns each consumer-facing port.

A context never imports Oracle or provider wire types.

## 10. D5 — HTTP API Contract

### Question

**What is the external application contract after legacy compatibility is removed?**

Inventory every current `operationId` and classify:

`KEEP / REDESIGN / RENAME / DELETE / INTERNAL-NOT-HTTP`.

For every target operation define:

- context/use-case owner;
- HTTP method/path;
- identity/auth/tenant semantics;
- request/response types;
- stable problem/error semantics;
- pagination/sort/filter semantics;
- idempotency/concurrency requirements;
- cache/freshness semantics;
- frontend consumer.

D5 must decide the actual OpenAPI toolchain and prove generation/validation against this repo. The target is one contract authority, not another hand-sync guard.

Also adjudicate whether a broad edge namespace such as `/api/* -> backend` removes duplicate Caddy/Vite route tables without harming OAuth/callback needs.

## 11. D6 — Frontend Contract

### Question

**How does a thin frontend represent product behavior without becoming a second domain authority?**

Map each screen/route to API/use-case/context.

For every surface define:

- route and navigation ownership;
- data query/mutation;
- cache/query-key ownership;
- loading/empty/stale/error states;
- unknown-data rendering;
- installation/account prerequisites;
- optimistic update policy, if any;
- error containment boundary;
- whether aggregation comes from a server projection or simple UI composition.

Adjudicate whether per-feature npm packages provide a real boundary/reuse benefit or whether `apps/web/src/features/*` is simpler for one application. Monorepo does not imply package-per-screen.

Legacy redirects have no default preservation requirement.

## 12. D7 — Runtime / Scheduler / Transactions / Outbox

### Question

**What processes run, what do they own, and how are lifecycle/concurrency/failure boundaries enforced?**

Decide based on failure/resource/isolation needs, not current `cmd/` count:

- HTTP serving process;
- scheduler;
- outbox dispatcher/consumers;
- long-running workers if justified;
- operator/admin commands;
- migrations;
- probes/test binaries and deletion conditions.

Define:

- transaction ownership;
- raw transaction APIs allowed only where appropriate;
- outbox commit semantics;
- leases/fencing/concurrency;
- cursor advancement atomicity;
- shutdown/readiness;
- retries/backoff;
- dead/outcome-unknown handling;
- runtime DB identity/RLS expectations;
- deployment topology.

Do not split into multiple processes solely because enterprise systems often do so.

## 13. D8 — Golden Flow Simulation

### Question

**Can the entire design explain real workflows without inventing ownership or transport mid-flow?**

Simulate at least:

### Flow A — ERP product to canonical operational state

```text
Sankhya/Oracle
 -> ERP adapter
 -> catalog/source ingest
 -> canonical/source identity resolution
 -> persistence/raw evidence
 -> events/projections/consumers
```

### Flow B — Marketplace listing observation and linking

```text
channel account/credential
 -> Mercado Livre adapter
 -> listings observation
 -> listing/variation identity
 -> linking decision/candidate
 -> inventory/pricing consumers
```

### Flow C — Safe external price/stock change

```text
user/use case
 -> recommendation/desired state
 -> preview/change control policy
 -> authorization/approval if required
 -> provider write
 -> success | failure | outcome unknown
 -> read-back
 -> verified | diverged
 -> audit/event
```

### Flow D — Marketplace order to realized profitability/reconciliation

```text
order observation
 -> item identity/link
 -> fees/shipping/payment facts
 -> ERP cost/tax/document facts
 -> reconciliation
 -> profitability/calculation record or projection
 -> UI/API
```

For every hop D8 names IDs, owner, contract, transaction, table/state, event, idempotency, retry and failure behavior. Any “we decide this during implementation” gap reopens the owning D-stage.

## 14. D9 — Adversarial Global-Maximum Review

### Question

**Do D1–D8 form one coherent system, or did each gate optimize locally?**

Review across axes:

- duplicated authorities;
- context cycles/tight synchronous coupling;
- cross-context DB ownership leaks;
- generic adapters/interfaces that encode the first vendor as “universal”;
- event misuse;
- duplicated identity vocabularies;
- API/frontend contract duplication;
- runtime mechanisms duplicated by contexts;
- invalid states still representable;
- transitional components without exit;
- unnecessary enterprise infrastructure;
- missing proof/golden-flow gaps.

D9 produces:

1. accepted target architecture;
2. explicit residual risks/deferred scope;
3. legacy deletion/cutover map;
4. implementation dependency DAG;
5. proof matrix for later implementation.

Only operator acceptance of D9 changes implementation authorization from BLOCKED to AUTHORIZED FOR PLANNING.

## 15. What implementation workers may decide later

After D9, workers may make local choices that do not change architecture, such as:

- private function decomposition;
- local variable names;
- equivalent SQL expression inside an owned repository;
- test helper structure;
- formatting and small internal refactors.

They may **not** silently decide:

- context ownership;
- new tables/authorities;
- sync vs async communication;
- event semantics;
- external credential ownership;
- retry/idempotency policy;
- API behavior not in the accepted contract;
- frontend business policy;
- new platform machinery;
- compatibility layers;
- deviation from the D1–D9 model.

A material ambiguity stops the slice and reopens the relevant design gate instead of being patched locally.

## 16. Documentation discipline during D0–D9

- `docs/architecture/rebaseline/README.md` is the only progress/status authority.
- This file is the program/gate authority.
- `ARCHITECTURE.md` is the concise architectural constitution.
- The root-cause/global-maximum standard is the decision method.
- Each D-stage gets at most one canonical design result plus tightly scoped evidence attachments when necessary.
- External API behavior is reverified from current official sources during D4; stale copied provider manuals are not retained as permanent authority.
- Git history preserves deleted plans/audits/ADRs/handoffs.
- No `wiki/`, `.mnfs/`, historical “implementation plan”, parallel roadmap or session handoff tree is created.

## 17. Definition of success

The rebaseline succeeds when a fresh implementation worker can receive a slice and answer, without architectural invention:

```text
what owns this behavior?
what owns this state?
what ID is canonical?
what table/constraint represents it?
who may write/read it?
is the interaction sync/event/projection?
which external adapter/port is involved?
which HTTP operation exposes it?
which frontend surface consumes it?
which transaction/idempotency/retry semantics apply?
what legacy is deleted?
what proof demonstrates completion?
```

That is the point at which implementation planning becomes cheaper than continued design.