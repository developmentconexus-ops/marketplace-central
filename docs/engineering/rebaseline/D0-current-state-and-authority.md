# D0 — Current State & Authority Baseline

> **Stage:** D0  
> **Status:** IN PROGRESS  
> **Evidence baseline:** `main@de1dc88bcef5a6ed5515378e7c646682c0bc15d2`  
> **Purpose:** describe what exists now without converting current implementation into target-architecture authority.

## 1. D0 target property

Before choosing contexts or target contracts, Marketplace Central must have one falsifiable current-state model of:

- code/package dependencies;
- runtime entrypoints and reachability;
- persistence ownership and writers;
- HTTP/API/SDK/frontend consumers;
- async/event/job paths;
- external integration boundaries;
- documentation/decision authority;
- recoverability of current persistent state.

D0 is complete when later design stages can reason from this census without rediscovering the repository or trusting stale historical plans.

## 2. Established repository facts

### 2.1 Mixed backend topology

The baseline has **21** directories under `apps/server_core/internal/modules/`:

`catalog`, `channelfees`, `classifications`, `connectors`, `dashboard`, `divergences`, `erp_import`, `integrations`, `internal_read`, `inventory`, `listings`, `market`, `marketplaces`, `mutations`, `orders`, `pricing`, `product_links`, `profitability`, `sourcekind`, `sync`, `tenant_config`.

The emerging target mold has **2** contexts under `apps/server_core/internal/contexts/`:

- `catalog`
- `listings`

`apps/server_core/internal/adapters/` already contains:

- `erp/`
- `marketplace/`

This is a live mixed architecture, not a hypothetical rewrite.

### 2.2 Runtime entrypoints

`apps/server_core/cmd/` contains seven entrypoint directories:

- `catalogingest`
- `listingsingest`
- `listingsreprocess`
- `migrate`
- `mlprobe`
- `server`
- `testdb`

D7 will decide the target process topology. D0 must first determine which are production-reachable, operator-only, test/migration tooling or temporary vertical-slice probes.

### 2.3 Frontend topology

Current reusable packages include:

- `feature-classifications`
- `feature-inventory`
- `feature-products`
- `feature-simulator`
- `sdk-runtime`
- `ui`
- `web-query`

Other product surfaces live directly in `apps/web/src/routes`, `apps/web/src/pages` and `apps/web/src/app`.

Current `AppRouter.tsx` still carries legacy redirects:

- `/products` → `/catalogo`
- `/product-links` → `/vinculos`
- `/inventory/stock-seguro` → `/estoque`
- `/orders` → `/pedidos`
- `/integrations` → `/integracoes`
- `/simulator` → `/precos`

There is no production compatibility requirement that justifies preserving them by default.

### 2.4 Contract topology

`contracts/api/marketplace-central.openapi.yaml` is the current HTTP contract artifact and contains a broad API surface from several architectural eras.

Historically, OpenAPI, Go handler/request shapes, TypeScript SDK methods and proxy/route tables have been manually synchronized. D5 must decide the target generation/validation mechanism; D0 must first enumerate the actual current mapping and drift.

### 2.5 Governance topology

`scripts/gate.ps1` is the active local/CI gate implementation. Current governance includes transitional ratchets whose scopes still name `internal/modules` and an API/SDK atomicity rule that assumes the current hand-authored SDK model.

Therefore:

- the gate is real and load-bearing today;
- its current rule vocabulary is not automatic proof of the target design;
- D0/D1/D5 must classify which controls survive, which are retargeted and which become removable because a stronger boundary replaces them.

### 2.6 Persistence topology

The migration chain spans multiple product/architecture eras: marketplace definitions/accounts, pricing/simulator, connectors/integrations, classifications, legacy modules and newer context state.

Migration existence proves historical runtime requirements, not target ownership.

### 2.7 External boundaries

Accepted current direction:

- Mercado Livre is the first operational marketplace;
- external marketplace integration enters through `internal/adapters/marketplace/<vendor>` and consumer-owned ports;
- Oracle/Sankhya access is MPC-owned and `godror`/OCI is the current driver/runtime path;
- `docs/operations/live-oracle-docker.md` remains a current operational reference for real read-only Oracle validation until D4/D7 replace or confirm it.

## 3. Documentation authority census

The baseline accumulated multiple competing planning surfaces:

- `IMPLEMENTATION_PLAN.md` — explicitly “Active but legacy”;
- root `EVIDENCE.md` — snapshot of a July worktree/mission;
- `.mnfs/` mission/debt/history tree;
- `docs/superpowers/{plans,specs,handoffs,evidence,runbooks}` spanning April–August designs;
- `docs/engineering/repo-audit-2026-08-07/` discovery/synthesis artifacts;
- `docs/design/` system blueprints, storage schemas, handoffs and evidence from older designs;
- `docs/research/` dated product/provider investigations and implementation handoffs;
- `docs/entregas/`, `docs/audit/` historical delivery/postmortem artifacts;
- `wiki/` with a second architecture/framework/module/quality/vision tree;
- provider docs for deferred marketplaces;
- `docs/HARNESS-PROFILE.md`, `docs/METODO-DE-REVISAO.md`, `docs/REVIEW-LEARNINGS.md` tied to retired wave/chip/hub processes.

The cleanup policy is:

### CANONICAL / survives

- `AGENTS.md`
- `docs/engineering/rebaseline/README.md`
- `docs/engineering/standards/root-cause-global-maximum-method.md`
- `ARCHITECTURE.md`
- `docs/architecture/decisions/README.md` and ADR records, with ADR-035 governing current/reopened status
- current D-stage design/evidence
- current machine-readable contracts/governance/runtime code

### SUPPORTING / explicitly non-authoritative

- `docs/engineering/defect-class-catalog.md` — historical defect-class evidence, not target architecture
- `docs/operations/live-oracle-docker.md` — current live-validation procedure, not target runtime design

### DELETE FROM LIVE TREE AFTER ABSORPTION

- `.mnfs/`
- `wiki/`
- `IMPLEMENTATION_PLAN.md`
- `EVIDENCE.md`
- `docs/superpowers/`
- `docs/engineering/repo-audit-2026-08-07/`
- `docs/design/`
- `docs/research/`
- `docs/entregas/`
- `docs/audit/`
- stale/deferred provider references under `docs/marketplaces/`
- loose superseded architecture/design/reference files outside the ADR registry
- retired harness/wave review doctrine files
- production deployment runbook tied to ADR-008 while D7 is explicitly re-adjudicating runtime/deploy topology

Git preserves all deleted content at the baseline SHA.

## 4. ADR authority state

ADR-035 establishes two classes during the rebaseline:

### Still binding as constraints unless a material finding reopens them

- ADR-005 — Mercado Livre first
- ADR-006 — MPC-owned Oracle reads
- ADR-007 — godror/OCI current Oracle runtime
- ADR-009 — fee provenance
- ADR-013 — webhook payload is pointer, not trusted domain data
- ADR-021 — TanStack Query owns frontend server state
- ADR-025 — raw provider PII is not retained
- ADR-027 — absence from partial pull is not closure
- ADR-029 — provider writes are not blindly retried
- ADR-033 — marketplace adapters + consumer-owned ports
- ADR-034 — `Fact` primitive supersedes the old unknown-is-zero prose mechanism; D2 still decides its proper scope

### Reopened for target design by ADR-035

All architecture tied to old module boundaries, hand-written SDK authority, specific scheduler/process topology, old mutation protocol, old market/divergence modules, old identity/linking assumptions and old mirror semantics is historical until its D-stage re-adjudicates it. See the ADR registry for exact mapping.

## 5. Root-cause families already visible

These are hypotheses for D0 grouping, not target solutions:

1. **Architecture authority split:** `modules` and `contexts` are both live.
2. **Contract authority split:** OpenAPI/runtime/SDK/routes have been manually synchronized.
3. **Identity/tenant authority split:** current runtime historically allowed boot/default tenant semantics while new commands moved toward explicit invocation identity.
4. **External-account/credential ownership:** provider installation/account/credential knowledge is spread across generic integrations and command/adapters.
5. **Data ownership split:** legacy tables/mirrors and new context tables coexist; exact writers/readers still need full census.
6. **Boundary instrumentation tied to legacy topology:** governance ratchets still measure `internal/modules`.
7. **Async/process semantics evolved incrementally:** `sync`, scheduler decisions, commands and newer ingest flows need one reachability map.
8. **Frontend/route authority split:** feature packages, app routes, hand-authored SDK and legacy redirects coexist.
9. **Evidence/document authority split:** multiple historical roadmaps/specs could independently direct a new session.

D0 must prove/refine these families with actual graph/runtime/data measurements before D1 uses them.

## 6. Remaining D0 measurements

The following are intentionally **not yet claimed as complete facts**.

### D0-M1 — Package/import graph

Measure all first-party Go packages and internal imports:

- SCCs/cycles;
- context → module and module → context edges;
- cross-module/context contract edges;
- adapter/platform/composition inversions;
- fan-in/fan-out;
- packages reachable only from tests/tooling.

### D0-M2 — Runtime composition/reachability

For each `cmd` and server route/job:

- entrypoint;
- composed services/modules/contexts;
- long-running vs operator/test/probe;
- direct external dependencies;
- lifecycle/shutdown;
- duplicated wiring.

### D0-M3 — Database ownership/write graph

Inventory every table/view/schema and every productive SQL site:

- all writers;
- all readers;
- foreign-context/module SQL;
- RLS/role behavior;
- duplicated state;
- immutable raw observations;
- tables whose state is fully re-derivable.

### D0-M4 — API/route/SDK/frontend topology

For each OpenAPI operation:

- handler/router reachability;
- request/response authority;
- SDK method;
- frontend/external consumer;
- proxy/edge routing dependency;
- orphan/legacy/duplicate surface.

### D0-M5 — Event/job/async topology

Inventory:

- domain/integration events;
- outbox tables/writers/consumers;
- pollers/schedulers/tickers;
- retries/backoff;
- cursor/checkpoint state;
- single-writer assumptions;
- async paths with no durable trace.

### D0-M6 — External-adapter topology

Map current Mercado Livre and Sankhya/Oracle code to consuming business capabilities and identify provider/vendor knowledge outside adapter boundaries.

### D0-M7 — Frontend topology

Map routes/screens/packages/query hooks/SDK methods/cache keys and business logic leakage.

### D0-M8 — Database recoverability

Classify current data by:

- re-derivable from Sankhya;
- re-derivable from Mercado Livre;
- configuration;
- re-authorizable credential/account state;
- human decision;
- non-rederivable audit/history.

This measurement decides whether a clean database baseline is the global maximum.

### D0-M9 — Active stale references

After documentary deletion, scan active source/workflows/contracts for references to removed authority paths. Replace instructions that remain operationally meaningful; historical comments may stay only if clearly historical and non-normative.

## 7. D0 exit criteria

D0 may be proposed for acceptance only when:

- M1–M9 are measured;
- counts/universes are explicit;
- current-state diagrams/matrices agree with runtime/code;
- document authority has one entry path;
- old documents cannot direct current work;
- contradictions are grouped by root cause;
- no target context/data/API choice has been smuggled in merely because current code has that shape;
- exact D1 questions are derived from D0 evidence.

## 8. Exact current next action

Continue **D0-M1 through D0-M9**, beginning with package/import graph and database ownership/write graph because they constrain the context-adjudication questions D1 must answer.

Do not start D1 implementation or create a third context.