# Marketplace Central — Architecture Rebaseline Design

**Date:** 2026-08-13  
**Status:** WRITTEN DESIGN — operator review required before implementation planning  
**Baseline:** `main@de1dc88bcef5a6ed5515378e7c646682c0bc15d2`  
**Scope:** architecture, ownership, contracts, communication, persistence, API, frontend, external boundaries, legacy disposition and proof strategy. No product implementation is authorized by this document alone.

---

## 0. Why this rebaseline exists

Marketplace Central is in the middle of a structural migration. The repository already contains two competing shapes:

- the legacy modular-monolith tree under `apps/server_core/internal/modules/`;
- the emerging context/adapters/kernel shape under `apps/server_core/internal/contexts/`, `internal/adapters/` and `internal/kernel/`.

The repository is not in production use. There is therefore no product requirement to preserve legacy routes, legacy module APIs, legacy database shapes, compatibility adapters or historical migrations merely to avoid breaking current users.

The goal is not to rewrite for aesthetic reasons. The goal is to remove the class of defects caused by deciding ownership and integration semantics locally during implementation.

The implementation program that follows this design must be able to answer, before writing feature code:

1. Which context owns this business fact?
2. Which component is the single writer/authority?
3. Which database object stores it, if any?
4. Is communication synchronous, asynchronous or a projection concern?
5. Which contract is public and which port is required?
6. Which external adapter translates provider/ERP semantics?
7. Which event exists, who produces it and who consumes it?
8. Which API operation exposes the behavior?
9. Which frontend surface consumes it?
10. Which legacy surface becomes deletable when the new authority lands?
11. What proof demonstrates that the old and new authorities are not simultaneously live?

If a feature PR still needs to invent one of these answers while implementing the feature, the rebaseline is incomplete.

---

## 1. Binding engineering method

### 1.1 Decision

Marketplace Central will adopt the Root-Cause / Global-Maximum Engineering Method proven in MetalDocs, adapted to this repository.

The project-local canonical method will live at:

`docs/engineering/standards/root-cause-global-maximum-method.md`

`AGENTS.md` will be a short binding bridge to that document, not a second copy of the method.

### 1.2 Binding principle

> Always simplify the code, never simplify correctness. Find the root cause, determine whether the proposed solution is only a local maximum, and prefer the global structure that makes the defect class unrepresentable or mechanically impossible at the strongest reasonable boundary.

The method must distinguish essential complexity from accidental complexity, define a local maximum versus a global maximum, require explicit authority/boundary ownership, define proof before implementation, and impose a bounded stop rule so that “global maximum” cannot become endless review.

### 1.3 Enforcement preference

For each material invariant, prefer the strongest reasonable layer:

1. structure/API makes invalid state unrepresentable;
2. type system makes it invalid;
3. database/schema constraint makes it invalid;
4. runtime boundary fails closed;
5. tests prove reachable behavior;
6. static guard/lint detects the violation;
7. documentation/convention.

A lower layer is legitimate when the stronger layer cannot express the complete property without disproportionate cost or loss of required flexibility. The reason must be explicit.

### 1.4 Guard policy

A new custom guard is not the default response to architectural drift. Before adding one, ask whether the duplicated authority, legal import path, legal API shape or legal state can be removed instead.

A guard is legitimate when it protects a distinct material property that cannot yet be expressed more strongly. A transitional guard must name its successor and deletion condition.

---

## 2. Decisions already approved by the operator on 2026-08-13

The following are decisions of this rebaseline, not merely hypotheses from older documents.

### D-RB-01 — Rebaseline before the next major feature implementation

The immediate program is architecture mapping and adjudication. New functional contexts are not added merely because they are next in the old sequence.

### D-RB-02 — No `old/` source-code cemetery

Legacy code will not be moved into an `old/`, `legacy/` or equivalent importable tree as a preservation strategy.

Git history is the archive.

Reason: an `old/` tree remains searchable, importable, compilable and copyable by humans and agents, creating a second architecture inside the active repository.

### D-RB-03 — Hard cutover is allowed; compatibility is not a default constraint

Because the platform has no production users depending on current behavior, routes, package APIs, IDs, schema shape and frontend navigation may break when the target design requires it.

Compatibility layers need a positive reason. “The old code existed” is not a reason.

### D-RB-04 — Hard cutover does not mean permanently broken `main`

A feature may be broken inside an implementation branch while a coherent cutover is being assembled. The merge target must remain epistemically trustworthy: its build/gates must distinguish intentional absence from accidental breakage.

The preferred cutover unit is:

```text
new authority lands
  -> consumers migrate
  -> persistence/contract cutover completes
  -> old authority has zero legal consumers
  -> old authority is deleted
  -> main is green and has one authority
```

Not:

```text
new authority + compatibility adapter + deprecated old authority + future cleanup promise
```

unless a measured constraint forces a transitional period.

### D-RB-05 — Existing protocol is evidence/hypothesis, not automatic authority

`docs/superpowers/specs/2026-08-06-protocolo-de-codigo-design.md` contains valuable measured work and several already-proven mechanisms. Its 13-context target is not accepted wholesale merely because it is documented.

Each structural conclusion is re-adjudicated against domain semantics, failure modes and the root-cause/global-maximum method.

### D-RB-06 — Issues are evidence; architecture is grouped by root cause and target property

Existing issues remain useful for traceability. They do not independently define the architecture program.

If five issues are symptoms of one authority problem, the implementation program is organized around the authority problem and the issues are acceptance evidence.

### D-RB-07 — API and frontend compatibility may be reset

There is no external-client compatibility requirement today. Legacy redirects, hand-maintained route aliases, stale packages and historical request/response shapes may be removed rather than preserved.

### D-RB-08 — No broker without a measured need

Async communication may use a transactional Postgres outbox and in-process/job consumers initially. Kafka, NATS, RabbitMQ or another broker requires a concrete scale, isolation or delivery requirement that Postgres cannot satisfy adequately.

### D-RB-09 — Database reset is an allowed candidate, not yet an automatic action

If all business state is re-derivable from Sankhya, marketplaces or re-authorizable configuration, the new architecture may establish a clean schema baseline instead of preserving a long migration chain built around authorities being removed.

Before a reset, all existing state must be classified by recoverability. Non-rederivable human decisions/audit material cannot be discarded silently.

---

## 3. Current-state evidence at the baseline

### 3.1 Backend topology

`apps/server_core/internal/modules/` contains 21 legacy module directories:

1. `catalog`
2. `channelfees`
3. `classifications`
4. `connectors`
5. `dashboard`
6. `divergences`
7. `erp_import`
8. `integrations`
9. `internal_read`
10. `inventory`
11. `listings`
12. `market`
13. `marketplaces`
14. `mutations`
15. `orders`
16. `pricing`
17. `product_links`
18. `profitability`
19. `sourcekind`
20. `sync`
21. `tenant_config`

`apps/server_core/internal/contexts/` contains two contexts in the new structural mold:

- `catalog`
- `listings`

`apps/server_core/internal/adapters/` already has the new top-level external-boundary split:

- `erp/`
- `marketplace/`

This means the repository is not choosing between two hypothetical architectures. It is already running a partial migration and therefore needs one explicit target before a third or fourth context increases the mixed-state surface.

### 3.2 Runtime entrypoints

`apps/server_core/cmd/` contains seven entrypoint directories:

- `catalogingest`
- `listingsingest`
- `listingsreprocess`
- `migrate`
- `mlprobe`
- `server`
- `testdb`

The rebaseline must classify each final runtime as one of:

- long-running serving process;
- scheduled/background worker capability;
- operator/admin command;
- migration/test tooling;
- temporary probe with explicit deletion condition.

A business capability must not require a permanent dedicated executable merely because the first vertical slice was easiest to drive that way.

### 3.3 Frontend topology

`packages/` currently includes:

- `feature-classifications`
- `feature-inventory`
- `feature-products`
- `feature-simulator`
- `sdk-runtime`
- `ui`
- `web-query`

while additional product surfaces live directly under `apps/web/src/routes` and `apps/web/src/pages`.

`AppRouter.tsx` still contains legacy redirects:

- `/products` -> `/catalogo`
- `/product-links` -> `/vinculos`
- `/inventory/stock-seguro` -> `/estoque`
- `/orders` -> `/pedidos`
- `/integrations` -> `/integracoes`
- `/simulator` -> `/precos`

No production user requires those aliases. Their default target disposition is DELETE unless an actual external dependency is measured.

### 3.4 Contract topology

`contracts/api/marketplace-central.openapi.yaml` is OpenAPI 3.1 and already exposes a broad API surface including mutation/change-control operations and many legacy domain operations.

The historical audit found contract truth duplicated across OpenAPI, runtime handlers/types, SDK and routing surfaces. The proposed protocol already points toward generated Go/TypeScript shapes from OpenAPI, but generator selection and dialect normalization remain decisions that must be proved, not assumed.

### 3.5 Persistence topology

The current migration chain contains legacy schemas/tables for marketplace definitions, pricing, connectors, classifications, simulator and other pre-rebaseline concepts. New contexts were added late in the chain.

The fact that a migration exists is historical evidence, not proof that its table belongs in the target model.

---

## 4. Target top-level shape — ratified direction

The exact set of contexts remains subject to adjudication, but the top-level dependency shape is ratified:

```text
apps/server_core/internal/
  contexts/       # business authorities / bounded contexts
  adapters/       # external-system translation and transport
  kernel/         # very small cross-context value semantics
  platform/       # technical runtime capabilities with no business ownership
  composition/    # final assembly / wiring only
  views/          # disposable read projections, if required
```

### 4.1 `contexts/`

A context owns business meaning and invariants. It publishes only intentional contracts and declares the ports it requires.

Preferred mold:

```text
contexts/<name>/
  contracts/              # what the context intentionally publishes
  port/                   # capabilities the context requires from outside
  internal/
    domain/
    application/
    postgres/
  module.go               # public facade; assembles its own internals
```

Go `internal/` is preferred over a custom guard whenever it can make cross-boundary imports impossible.

### 4.2 `adapters/`

Adapters own wire/system-specific semantics and implement ports declared by consuming contexts.

Target family:

```text
adapters/
  marketplace/<vendor>/
    <vendor>.go
    internal/api/
    <capability implementations...>

  erp/sankhyaoracle/
    internal/oracle/
    <context-facing capability implementations...>

  spreadsheet/<import type>/
```

No context should need to import Mercado Livre DTOs, Oracle driver types, raw SQL helpers or OAuth wire semantics.

ADR-033 already ratifies the marketplace direction: external marketplace integrations enter through `adapters/marketplace/<vendor>`, implementing ports owned by consumers; the legacy `connectors` module receives no new marketplace code.

### 4.3 `kernel/`

Kernel admission must be difficult.

A type belongs here only when its meaning and invariants are genuinely identical across multiple contexts and it has no vendor/context ownership.

Candidate members from the existing protocol include tenant/channel identity, exact decimal/money, knowledge/fact semantics, provenance and effective periods. Each candidate still requires evidence that centralization reduces duplicated authority rather than creating a generic domain bucket.

### 4.4 `platform/`

Platform owns technical mechanisms without domain policy, for example:

- PostgreSQL connection/transaction infrastructure;
- HTTP boundary primitives;
- logging/telemetry;
- scheduling/leases/backoff mechanics;
- outbox dispatch mechanics;
- config loading;
- migrations.

Platform must not own business enums such as “sync type = products/orders/market.” Business-specific cursor/freshness semantics belong to the consuming context.

### 4.5 `composition/`

Composition wires public facades and external adapter bundles. It is not a place to reimplement provider adapters or domain orchestration.

A recurring smell is composition needing to name private implementation types. If that happens, first test whether the public facade is wrong before adding an import exception.

### 4.6 `views/`

A view/projection exists only when a read concern legitimately combines several authorities for a user-facing query/dashboard.

A view is:

- rebuildable/disposable;
- not authoritative;
- incapable of accepting business commands;
- explicit about `as_of`/freshness/completeness when material;
- fed from versioned contracts/events or explicit read contracts.

A “dashboard module” with repository access to every context is not the target.

---

## 5. Context adjudication — the 13-context proposal is a candidate set

The existing protocol proposes:

`account`, `catalog`, `listings`, `linking`, `inventory`, `costing`, `tax`, `pricing`, `orders`, `reconciliation`, `profitability`, `intelligence`, `changecontrol`.

The rebaseline will not accept a context because its name appears in that list. Each candidate must pass the following test.

### 5.1 Context admission test

For each candidate, answer:

1. What business question does this context own that no neighbor owns?
2. What state can change independently here?
3. What invariant/lifecycle does it enforce?
4. What commands does it accept?
5. What facts/contracts does it publish?
6. Which source systems inform it but do not own its semantics?
7. If merged with its nearest neighbor, which named failure mode appears?
8. If split from its nearest neighbor, which concrete coupling is reduced?
9. Can it be independently tested from contracts without reading its internals?
10. Does it deserve authoritative persistence, or is it only a projection/use case?

A context with no independent lifecycle/authority is presumed to be a feature, projection or component rather than a bounded context.

### 5.2 Candidate context matrix

| Candidate | Provisional authority | Status before adjudication | Main question |
|---|---|---|---|
| `account` | tenants, users, channel/ERP installations, credential refs | **ADJUDICATE** | Are human IAM and external installation identity one domain or two? |
| `catalog` | canonical product identity, variants, source observations/imports | **PROVISIONAL / already implemented in new mold** | Does costing belong outside catalog; are canonical IDs fully source-opaque? |
| `listings` | observed marketplace listing/variation state | **PROVISIONAL / already implemented in new mold** | Which published facts are stable across marketplace vendors? |
| `linking` | product<->external-variation match candidates, decisions, evidence | **ADJUDICATE** | Is matching a distinct decision lifecycle or catalog/listings application logic? |
| `inventory` | stock observations, safety policy, sellable/desired availability | **ADJUDICATE** | Separate observation from desired-state/write control cleanly. |
| `costing` | effective-dated costs by method/company/location | **ADJUDICATE** | Is cost lifecycle independent enough from catalog to justify a context? |
| `tax` | fiscal classification/rules/determinations with effective time | **ADJUDICATE** | Which values are copied source facts vs MPC determinations? |
| `pricing` | scenarios, policy resolution, expected fees/freight, recommendations | **ADJUDICATE** | Keep pricing predictive; do not absorb realized profitability. |
| `orders` | observed marketplace orders/items/payments/shipments/lifecycle | **ADJUDICATE** | Which financial fields are observed facts versus derived profitability? |
| `reconciliation` | pairing of ERP documents/items with marketplace orders/items | **ADJUDICATE** | Is reconciliation independently stateful/historical enough for own authority? |
| `profitability` | realized calculation with exact versioned inputs/policy | **ADJUDICATE** | If it is merely a dashboard formula, make it a projection; if historical/recomputable, context may be justified. |
| `intelligence` | competition observations/matching confidence/signals | **ADJUDICATE** | If it only aggregates existing facts, it is a view, not an authority. |
| `changecontrol` | preview/approval/execution/retry/verification/audit of external writes | **ADJUDICATE** | Ensure it owns write protocol, not every domain's desired-state policy. |

### 5.3 Specific challenges that must be resolved

#### `account`

Do not create a god-context simply because several concepts contain “account.” Human principals/session/authorization and a marketplace installation/seller identity have different lifecycles. They may still belong together, but the reason must be domain evidence rather than naming convenience.

#### `costing`

This is the most explicitly questionable split in the older protocol. Effective-dated cost method/company/location behavior supports separation, but if catalog is the only legal consumer and lifecycle is not independently meaningful, a separate context may add accidental ceremony.

#### `profitability`

A versioned realized calculation that stores exact inputs, method/policy version and explanation can be authoritative. A dashboard computation over orders/cost/tax is a projection. These are different products; the design must pick one.

#### `intelligence`

Do not promote analytics to a bounded context unless it owns genuinely new observations/decisions. Aggregation alone is not authority.

---

## 6. Provisional legacy disposition map

This table is not permission to delete. It ensures every legacy module has a destination question before implementation begins.

| Legacy module | Provisional destination | Disposition class | Deletion evidence required |
|---|---|---|---|
| `catalog` | `contexts/catalog` | REPLACE | all authoritative reads/writes use context; zero legal imports of old module; old tables classified/migrated/reset |
| `listings` | `contexts/listings` | REPLACE | all ingest/read consumers use context; provider DTOs stay in adapters; old storage no longer authoritative |
| `channelfees` | likely `pricing` and/or observed channel/account facts | DISSOLVE / ADJUDICATE | fee authority decided; no duplicated fee formula/schedule source |
| `classifications` | likely `tax` | DISSOLVE / ADJUDICATE | fiscal classification ownership and effective-time semantics decided |
| `connectors` | `adapters/marketplace/*` | DELETE AFTER CUTOVER | all provider capabilities supplied by new adapters; zero imports/writes through connectors |
| `dashboard` | `views/` or direct context queries | DISSOLVE | each dashboard datum has named authority; projection rebuildable and non-authoritative |
| `divergences` | likely `reconciliation`, `inventory`, `linking` or view-specific logic | DISSOLVE / ADJUDICATE | every divergence type has an owning domain decision, not a generic divergence bucket |
| `erp_import` | `adapters/erp/sankhyaoracle` + consuming context ingest | DISSOLVE | raw ERP access confined to adapter; import state/cursors owned by consumer context |
| `integrations` | account/installations or dedicated integration capability after adjudication | REPLACE / ADJUDICATE | provider installation/auth lifecycle has one authority |
| `internal_read` | split among context ports + ERP adapter | DISSOLVE | no generic read-domain authority; no downstream Oracle knowledge |
| `inventory` | `contexts/inventory` | REPLACE | observed/desired stock semantics owned once; external writes go through approved write boundary |
| `market` | likely `intelligence` observations or a projection | DISSOLVE / ADJUDICATE | no generic “market” bucket; each observation has explicit semantics |
| `marketplaces` | installation/account + pricing policy | DISSOLVE | account identity and pricing policy authorities are separated/ratified |
| `mutations` | likely `changecontrol` | REPLACE / ADJUDICATE | protocol state machine and external-write authority ratified |
| `orders` | `contexts/orders` | REPLACE | provider observations mapped once; financial derivations owned elsewhere |
| `pricing` | `contexts/pricing` | REPLACE | one exact money/decimal vocabulary; no duplicate simulation engine |
| `product_links` | likely `contexts/linking` | REPLACE / ADJUDICATE | linking decision/evidence lifecycle ratified |
| `profitability` | `contexts/profitability` or projection | ADJUDICATE | authoritative versus derived role settled before migration |
| `sourcekind` | likely `kernel/provenance` and per-context source semantics | DISSOLVE | no generic defaulting source enum remains |
| `sync` | `platform/scheduler` mechanics + per-context cursors/policy | DISSOLVE | platform contains no domain sync enum; each context owns freshness/cursor semantics |
| `tenant_config` | account/tenant domain or platform config depending on field semantics | DISSOLVE / ADJUDICATE | each config key has one owner and channel (deployment vs invocation vs business configuration) |

A module may only disappear when its authority is accounted for, not simply because the directory is undesirable.

---

## 7. Internal communication protocol

The architecture recognizes three different communication needs. They must not be conflated.

### 7.1 Synchronous contract/port call

Use when the caller cannot complete the current use case without an immediate answer.

Examples of legitimate shape:

```text
pricing use case
  -> pricing-owned port
  -> costing/catalog/inventory facade or adapter implementation
  -> typed answer
```

Rules:

- no cross-context SQL as a shortcut;
- no import of another context's `internal/` package;
- the consumer owns the port it needs;
- the producer owns the contract it publishes;
- a synchronous dependency must be explicit in the context dependency map.

Do not publish an event and synchronously wait for it merely to avoid a direct dependency.

### 7.2 Asynchronous domain/integration event

Use when a fact has already happened and downstream work can react without participating in the producer's transaction decision.

Baseline event envelope candidate:

```text
event_id
schema/event_type
version
tenant_id
occurred_at
aggregate/entity reference when applicable
correlation_id
causation_id when applicable
payload
```

Each event catalog row must name:

- semantic name and version;
- producer (exactly one authority);
- triggering state transition;
- payload contract;
- consumers;
- idempotency/deduplication key;
- ordering requirement, if any;
- retry/dead-letter/poison behavior;
- retention/replay expectations;
- whether it is domain-internal, cross-context integration, or projection feed.

### 7.3 Projection/view feed

Use for dashboards/search/read models that combine multiple authorities.

A projection consumer must not become a hidden command path or business-rule owner.

### 7.4 Transactional publication direction

The preferred initial mechanism is a Postgres transactional outbox written in the same transaction as the authoritative state transition, then dispatched by platform/runtime mechanics.

This direction remains subject to proof against actual throughput and delivery requirements, but external brokers are not introduced without a measured reason.

---

## 8. External integration map — required before implementation

A statement such as “Mercado Livre adapter exists” is insufficient. Every external capability must be catalogued at operation level.

Required matrix columns:

| Field | Meaning |
|---|---|
| Provider/system | Mercado Livre, Sankhya Oracle, spreadsheet, etc. |
| Capability | listings feed, orders read, price update, catalog import, etc. |
| External operation | endpoint/query/procedure identifier |
| Direction | read / write / auth / callback |
| Consuming context | owner of the port |
| Adapter package | exact implementation boundary |
| Credential principal | seller/install/user/application identity used |
| Pagination/cursor | provider semantics and replay constraints |
| Rate limit | provider policy and coordination scope |
| Freshness | source timestamps versus observed time |
| Retry | safe, unsafe, conditional |
| Idempotency | provider key / local protocol |
| Timeout ambiguity | can the remote action succeed after local timeout? |
| Raw retention | exact bytes/value/hash requirement |
| Error translation | stable domain/integration errors versus raw provider error |
| Verification/read-after-write | how convergence is proven for writes |

### 8.1 Provider knowledge boundary

Mercado Livre-specific concepts such as API DTO names, `scroll_id`, OAuth payload shapes, `SELLER_SKU` wire representation, provider error bodies and provider URL structure belong under the Mercado Livre adapter tree.

Contexts may know a runtime channel code/account reference when business semantics require channel identity; they do not know the provider's HTTP representation.

### 8.2 ERP boundary

Oracle SQL, `godror`, Sankhya table names and query-specific null/default behavior belong to the Sankhya Oracle adapter.

Context contracts should expose business facts with explicit knowledge/freshness semantics, not Oracle rows.

---

## 9. Persistence map — required before schema work

Every target table/projection must appear in a database ownership catalog before its migration lands.

Required columns:

| Field | Required question |
|---|---|
| schema.table | physical object |
| context owner | business authority |
| semantic truth | what fact/state it means |
| writer | exactly which repository/use case writes |
| readers | direct owner reads, contract projection, admin only |
| tenant strategy | tenant key/RLS/role expectations |
| primary identity | canonical/local/source-derived? |
| uniqueness | business uniqueness and tenant scope |
| temporal semantics | observed/effective/source-updated/created |
| exactness | decimal/money/quantity representation |
| raw evidence | whether source payload/hash is retained |
| retention | permanent, bounded, rebuildable |
| cross-context references | semantic reference and integrity mechanism |
| legacy predecessor | table(s) it replaces |
| recovery source | marketplace/ERP/human decision/not rederivable |

### 9.1 Strong candidates already supported by current evidence

- one authoritative owner per table/state family;
- one intended writer authority;
- tenant isolation must not rely on developers remembering a `WHERE tenant_id = ?` predicate;
- serving application role must not defeat RLS by ownership/superuser/BYPASSRLS privileges;
- financial values must not use `float64` as authoritative storage;
- raw source observations that justify a mapped fact should preserve enough evidence for replay/audit when replay is a product requirement.

### 9.2 Decisions deliberately NOT ratified yet

#### One Postgres schema per context

The existing protocol proposes it. The rebaseline will test whether schema boundaries materially strengthen ownership/permissions or merely add naming ceremony for a single modular monolith.

#### No foreign keys across context schemas

This is not accepted as dogma. A cross-context FK can create coupling, but removing every FK can also discard a stronger database invariant and replace it with convention/reconciliation.

Decision rule:

- if lifecycles are truly independent and only semantic references cross the boundary, no FK may be correct;
- if synchronous referential integrity is a real invariant, a FK may be the strongest correct mechanism;
- if many cross-context FKs are required, first question whether the context boundary is wrong.

#### Database reset

Before reset, classify every current state item:

1. rederivable from Sankhya;
2. rederivable from marketplace;
3. deployment/configuration;
4. credential/re-authorizable;
5. human-authored decision/configuration;
6. historical/audit evidence that cannot be recreated.

Only categories whose loss is explicitly acceptable can be destroyed.

If categories 5/6 are empty or disposable, a new clean baseline migration becomes a preferred candidate over preserving historical tables for compatibility.

---

## 10. Identity, exactness, knowledge and time — kernel/domain adjudication

The old protocol correctly identified several cross-cutting semantic defects, but the rebaseline will preserve properties rather than universalize mechanisms blindly.

### 10.1 Canonical identity

Canonical MPC identity must not accidentally equal one source system's identifier when the product is expected to integrate multiple sources.

Source keys belong in explicit source-reference structures. The final identity strategy for each aggregate must state:

- canonical ID type/generator;
- tenant scope;
- source references;
- merge/supersession semantics if identities are reconciled;
- external provider IDs as typed source references, not naked strings where confusion is material.

### 10.2 Exact money/decimal

Authoritative monetary arithmetic/storage cannot rely on binary floating point.

A shared exact value type is a strong candidate for kernel admission because money semantics recur across pricing, orders, costing, taxes and profitability. The exact type/API still needs contract/DB/JSON compatibility proof.

### 10.3 Knowledge / `Fact[T]`

The property “unknown is never silently converted into zero/default” is ratified.

The mechanism “every business value in every contract must be `Fact[T]`” is not ratified.

Use knowledge-state wrapping for observations/derived facts where absence, uncertainty, estimation or not-applicable are real domain states.

Do not wrap deterministic configuration merely to make all fields visually uniform.

Examples:

- observed marketplace price -> knowledge/fact semantics are plausible;
- observed ERP stock -> knowledge/fact semantics are plausible;
- estimated freight -> knowledge/fact semantics are plausible;
- user-configured page size -> ordinary integer;
- explicitly configured target margin -> exact decimal/config value unless the business model itself makes it uncertain.

### 10.4 Time

The design must distinguish at least:

- when the source says a fact changed/effectively applies;
- when MPC observed it;
- legal/business effective period where relevant;
- local record creation/update time.

One timestamp must not silently stand for another.

---

## 11. External write/change-control protocol

No live marketplace write should be treated as a simple HTTP command with boolean success.

The candidate `changecontrol` authority must be adjudicated around these properties:

1. desired state is distinct from last requested state;
2. last observed provider state is distinct from both;
3. preview/approval policy is explicit for material writes;
4. a timeout after request transmission may produce `outcome_unknown`, not automatic failure;
5. retries depend on idempotency/provider semantics;
6. convergence is proven by later observation when the provider response does not itself prove final state;
7. audit records actor, intent, source inputs, policy decision, remote attempt and verification.

The context that owns the business decision keeps its desired-state policy. The change-control mechanism owns the safe execution protocol. It must not become a god-context containing pricing, inventory and order business rules.

---

## 12. API rebaseline

### 12.1 Contract authority

Target direction:

```text
OpenAPI source
  -> generated/validated Go server contract
  -> generated TypeScript client/types
```

Handwritten business implementation remains handwritten; transport shapes should not have multiple authorities.

Before generator adoption:

- normalize/validate OpenAPI dialect usage;
- evaluate generator support against actual 3.1 constructs;
- prove request validation/error mapping behavior;
- prove generated client behavior needed by the frontend;
- pin tool versions and add regeneration drift proof.

### 12.2 Operation inventory

Every operation in the target API must be represented in a matrix:

| operationId | method/path | owning context/view | use case | request | response | errors | auth/tenant | idempotency | frontend consumer | disposition |
|---|---|---|---|---|---|---|---|---|---|---|

Each current operation is classified as:

- KEEP;
- RENAME/RESHAPE;
- DELETE;
- INTERNAL/ADMIN only;
- REPLACE by projection/query.

No route remains solely because an old frontend references it; the frontend and API are rebaselined together.

### 12.3 Edge routing simplification candidate

A strong global-max candidate is a coarse edge split such as:

```text
/api/* -> backend
/*     -> SPA
```

with explicit health/operational endpoints as needed, rather than hand-maintaining business route lists in Caddy, Vite and Go.

This is a candidate because the historical audit measured several manually synchronized route authorities. It must be checked against current OAuth callback/static-routing requirements before ratification.

The repository keeps the existing frozen decision of no `/v1` prefix unless a new ADR changes it; `/api` would be an edge namespace, not semantic API versioning.

---

## 13. Frontend rebaseline

The web client remains thin: business rules live in Go/context contracts.

### 13.1 Legacy navigation

The current legacy redirects have no preservation requirement and are default DELETE candidates.

### 13.2 Feature packaging

The current frontend mixes npm packages (`feature-*`) with app-local routes/pages.

The target is not automatically “one package per screen.” Before preserving the workspace-per-feature pattern, measure whether those packages have independent reuse/build/test/ownership value.

Simpler candidate:

```text
apps/web/src/
  app/
  features/
    integrations/
    catalog/
    listings/
    linking/
    inventory/
    pricing/
    orders/
    reconciliation/
    profitability/
  shared/

packages/
  api-client/ or sdk-runtime/
  ui/ only where genuinely shared
```

Package boundaries are retained only where they enforce useful dependency/reuse properties.

### 13.3 Data access

Frontend code consumes the generated/authoritative client. No feature authors ad hoc `fetch` URLs to bypass the contract.

TanStack Query/query helpers may remain if they provide useful cache/lifecycle semantics, but `web-query` must not become a second hand-authored API contract.

### 13.4 Failure containment

Route/app failure boundaries are a product behavior concern, not merely a lint concern. The rebaseline will preserve the prior finding that a single render failure should not necessarily destroy the entire SPA, but exact boundaries follow final screen topology.

---

## 14. Runtime and scheduler map

The final runtime topology must be designed around capabilities, not historical binaries.

### 14.1 Serving runtime

`cmd/server` is the likely HTTP composition root.

### 14.2 Background work

Polling/scheduled ingest may be hosted by one worker process, the serving process, or separate operational binaries depending on isolation/resource/failure requirements.

The architecture does not ratify one executable per business feed.

### 14.3 Scheduler ownership

Platform scheduler owns technical mechanics only:

- registration;
- tenant/account partitioning;
- lease/fencing token if concurrency requires it;
- rate-limit coordination primitive;
- backoff;
- telemetry;
- lifecycle/shutdown.

Each context owns:

- cursor shape;
- overlap window;
- deduplication key;
- freshness meaning;
- full-scan policy;
- event/state semantics.

### 14.4 Temporary commands

`mlprobe`, dedicated reprocessors and other operator probes must be classified as permanent admin tooling or given a deletion condition. “Useful during development” is not enough for a permanent runtime surface.

---

## 15. Golden product flow

Architecture is accepted only if it can support an end-to-end real flow, not merely isolated green packages.

A provisional Marketplace Central golden flow is:

```text
BOOT
 -> operator/user identity resolved
 -> tenant/installations available
 -> connect Mercado Livre account
 -> connect/read Sankhya source
 -> ingest canonical products
 -> ingest marketplace listings
 -> link product <-> listing variation with evidence
 -> observe stock/cost/tax/fee facts
 -> produce pricing recommendation/preview
 -> optionally approve/execute controlled write
 -> ingest orders
 -> reconcile marketplace order/item with ERP document/item
 -> compute realized profitability from versioned inputs
 -> show current state/projections in frontend
 -> retain sufficient evidence to explain/replay material decisions
```

Each stage will receive:

- owning context/view;
- API operation(s);
- tables/state;
- event/port dependencies;
- external operation;
- frontend screen;
- failure modes;
- proof fixture/live evidence;
- legacy components it invalidates.

A root cause that makes one of these stages untrustworthy is prioritized before lower-value cleanup.

---

## 16. Proof model

### 16.1 Architecture proof

For every structural rule, state the mechanism and the counterexample that must fail.

Examples:

- cross-context private import -> compiler must reject;
- vendor DTO outside vendor tree -> compiler/static detector must reject;
- generated contract drift -> regeneration check must produce diff/red;
- app DB role bypasses RLS -> boot/integration proof must fail;
- legacy writer still mutates replaced table -> ownership scan/integration proof must fail.

### 16.2 Positive behavior proof

Architecture-only tests do not prove the product works. The final program requires real Postgres integration and bounded live-provider/Oracle evidence for capabilities that depend on external semantics, with writes requiring explicit authorization.

Mocks remain valid for deterministic contract behavior but do not certify live integration.

### 16.3 Anti-vacuity

A green lane must prove the intended test/check actually ran. Existing gate principles remain binding.

### 16.4 Stop rule

A design decision stops being re-litigated when:

- root cause and target property are explicit;
- authority/boundary is coherent;
- chosen mechanism has a named failure mode it prevents;
- relevant counterexample/positive proof exists or has a precise proof plan;
- no material contradiction remains;
- remaining concerns are speculative or non-material.

Reopen only on new material evidence or changed constraints.

---

## 17. Rebaseline artifact set

This program should end with a small canonical set, not another pile of overlapping plans.

### Canonical

1. `ARCHITECTURE.md` — product-level architecture/north star and frozen decisions.
2. `docs/architecture/decisions/*` — material architecture decisions/amendments.
3. `docs/engineering/standards/root-cause-global-maximum-method.md` — engineering decision method.
4. `AGENTS.md` — operational bridge/router to authorities and verification.
5. `contracts/api/marketplace-central.openapi.yaml` — external HTTP contract source.
6. `contracts/governance/*` — mechanical rules/registries that remain necessary.

### Rebaseline working/decision artifact

This document remains the design record until its material decisions are either:

- incorporated into `ARCHITECTURE.md`;
- ratified by an ADR;
- implemented as a structural/mechanical rule;
- explicitly rejected with evidence.

Old plans/audits remain history and evidence. They are not silently rewritten to look as if they had predicted the final design.

---

## 18. Sequencing constraints after this written design is approved

This is not yet the implementation plan. It defines only dependency order.

```text
R0  Ratify engineering method + AGENTS bridge
 |
R1  Complete current-state ownership inventory
 |
R2  Adjudicate target contexts + kernel + runtime
 |
 +----> external capability matrix
 +----> DB ownership/recoverability matrix
 +----> API operation disposition matrix
 +----> frontend surface disposition matrix
 +----> event/port dependency matrix
 |
R3  Ratify architectural contradictions via ADR/ARCHITECTURE updates
 |
R4  Decide clean DB baseline versus migration preservation
 |
R5  Select one narrow real vertical slice that crosses the target boundaries
 |
R6  Prove/refute target mold with that slice
 |
R7  Migrate/delete legacy authority by authority
 |
R8  Remove transitional instruments and stale historical active guidance
```

The vertical slice must be narrow in breadth and deep in correctness. It is a falsifier for the architecture, not an excuse to implement all 13 candidate contexts before learning.

---

## 19. What this design deliberately does not decide yet

The following are open adjudications with explicit evidence requirements, not placeholders:

1. **Exact context count.** Resolve with the context admission test and current domain flows.
2. **Whether IAM/human identity belongs with marketplace/ERP installation identity.** Resolve by lifecycle/principal/authorization analysis.
3. **Whether `costing` is separate from `catalog`.** Resolve by independent lifecycle/consumers/invariants.
4. **Whether `profitability` is authority or projection.** Resolve by requirements for historical recomputation/versioned inputs/explanation.
5. **Whether `intelligence` is authority or projection.** Resolve by existence of independently owned observations/decisions.
6. **One DB schema per context.** Resolve by permissions/ownership benefits versus ceremony.
7. **Cross-context FK policy.** Resolve per invariant; no blanket ban or blanket allowance.
8. **Clean database reset.** Resolve by state recoverability inventory.
9. **Exact `Fact[T]` scope.** Preserve unknown-not-zero property; apply wrapper only to semantically uncertain facts.
10. **Frontend package topology.** Resolve by actual reuse/enforcement needs.
11. **OpenAPI generator choice and 3.1 normalization.** Resolve by spike against actual contract and required generated runtime/client behavior.
12. **`/api/*` edge namespace.** Resolve against OAuth callbacks, current deployment/proxy routing and stable-route frozen decision.
13. **Background-process topology.** Resolve by failure isolation, scheduling, resource and operational needs rather than current `cmd/` count.

No implementation task may silently settle one of these. It must either reference the adjudication result or stop and produce the missing decision.

---

## 20. Spec self-review

### Placeholder scan

This spec intentionally contains no `TBD` or `TODO`. Unsettled decisions are listed in §19 with an evidence-based adjudication rule.

### Internal consistency

- Hard cutover is allowed, but deletion requires authority accounting and a green/verifiable merge state.
- Existing accepted ADRs remain authoritative until explicitly amended.
- The older 13-context protocol is treated as evidence, not overwritten or silently promoted.
- The method prefers stronger enforcement but does not ban legitimate lower-level guards.
- The design preserves unknown-not-zero without requiring a universal `Fact[T]` wrapper.
- The design permits DB reset without assuming all current state is disposable.

### Scope check

This is intentionally broader than one feature because its job is to define the architecture program and prevent later feature PRs from making incompatible local decisions. It does **not** itself authorize the multi-context migration.

### Review gate

The next action after operator approval of this written spec is to write the implementation plan for the **rebaseline program**, starting with R0/R1. Product-feature implementation remains blocked until the relevant rebaseline decisions are ratified.