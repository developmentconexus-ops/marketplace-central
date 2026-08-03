# Repo map — measured anchors

Every anchor here was read from the tree. Line numbers rot; **paths, symbol names and file
counts are the durable part** — re-measure a line before you cite it in a plan, and name the
tree you measured it in.

## Contents

- [Layout](#layout)
- [Go module anatomy](#go-module-anatomy)
- [Module registry (measured)](#module-registry-measured)
- [The contract seam](#the-contract-seam)
- [Frontend](#frontend)
- [Platform packages](#platform-packages)
- [Sync and scheduling seam](#sync-and-scheduling-seam)
- [Provider seam and rate budget](#provider-seam-and-rate-budget)
- [Governance registries](#governance-registries)
- [Verification lanes](#verification-lanes)
- [Where documents live](#where-documents-live)

## Layout

```
apps/server_core/     Go backend — all business logic
apps/web/             thin React client (Vite + React Router)
apps/landing/         marketing surface, not part of the cockpit
packages/sdk-runtime/ the ONLY way the web talks to the backend
packages/ui/          shared primitives
packages/web-query/   shared data-query layer + shared formatters
packages/feature-*/   workspace screens (classifications, connectors, inventory,
                      orders, products, simulator)
contracts/api/        OpenAPI — source of truth for HTTP behaviour
contracts/governance/ machine-checked registries (see below)
```

Some `packages/feature-*` surfaces are historical and superseded by pages under
`apps/web/src/pages/`. Before planning FE work in a `feature-*` package, grep for its
component from `apps/web/src` and confirm it is actually mounted — an unmounted package is
dead code and editing it produces a plan that changes nothing on screen.

## Go module anatomy

Every module under `apps/server_core/internal/modules/<name>/` uses the same five layers:

| Layer | Holds | Never holds |
|---|---|---|
| `domain/` | entities, value objects, pure rules | I/O, SQL, provider types |
| `ports/` | the interfaces the module owns | implementations |
| `application/` | services orchestrating ports | SQL, HTTP, provider payloads |
| `adapters/<tech>/` | postgres / oracle / provider implementations | business rules |
| `transport/` | HTTP handlers, query parsing, OpenAPI contract tests | business rules |

Read `ports/` first when planning against a module you don't know. `orders/ports/` alone
carries 15 interfaces — the module's real contract surface is there, not in `application/`.

Adapters are named by technology, and a module may have several:
`orders/adapters/{postgres,integrations,internalread,pricingtax,productlinks}/`.

## Module registry (measured)

20 modules exist under `internal/modules/`: `catalog`, `channelfees`, `classifications`,
`connectors`, `dashboard`, `divergences`, `erp_import`, `integrations`, `internal_read`,
`inventory`, `listings`, `market`, `marketplaces`, `mutations`, `orders`, `pricing`,
`product_links`, `profitability`, `sourcekind`, `sync`, `tenant_config`.

`contracts/governance/modules.json` is the machine-checked truth for each one: its root,
`composition_required`, its `openapi_prefixes`, and its allowed `dependencies`. **A new
module without an entry fails `GOV_MODULE_COVERAGE`.** A new cross-module import that is not
in that module's `dependencies` fails `GOV_MODULE_DEPENDENCY`; importing another module's
`adapters/` (rather than its published interfaces) fails `GOV_MODULE_LAYER`.

The file also carries `temporary_exceptions`, each with a `removal_owner`. If your plan
touches a path named in an exception, either remove the exception or say why it survives.

## The contract seam

- Spec: `contracts/api/marketplace-central.openapi.yaml`
- Client: `packages/sdk-runtime/src/index.ts` (plus domain files, e.g. `market.ts`,
  `dashboard.ts`, `erpImport.ts`, `activeSource.ts`)

The SDK is **hand-maintained**, not generated. That is why the atomicity invariant
(`api-sdk-atomicity` → `GOV_API_SDK_SPLIT`) exists: the two files must move in the same
commit, or the client and the spec disagree with nothing to catch it.

Types are spread across the SDK's files — `MarketPriceIntelAggregate` lives in
`packages/sdk-runtime/src/market.ts`, not `index.ts`. Grep the whole `src/` directory for a
type name; do not assume `index.ts`.

Parity is enforced two ways, both worth extending in a plan that adds routes:

1. **Transport-side contract tests**, e.g.
   `internal/modules/market/transport/openapi_contract_test.go` — reads the YAML by relative
   path and string-slices each path block, asserting `operationId`, request/response `$ref`s
   and status codes. No YAML parser or spectral in the repo; follow the string-slicing idiom.
2. **SDK-side contract tests** under `packages/sdk-runtime/src/*.test.ts`, which read the
   same YAML via `resolve(process.cwd(), "../../contracts/api/marketplace-central.openapi.yaml")`.
   That relative path is why the FE test lane's working directory is part of the lane.

Frontend must never `fetch` directly — invariant `frontend-fetch` (`GOV_FRONTEND_FETCH`)
scans `apps/web/src` and `packages`.

## Frontend

- Pages: `apps/web/src/pages/<area>/` — `dashboard`, `importacoes`, `integracoes`, `mercado`,
  `mutations`, `pedidos`, `precos`, `produto`, `vinculos`, plus top-level listing screens.
- Tests sit beside the component (`AnunciosTable.tsx` / `AnunciosTable.test.tsx`).
- Shared formatting lives in `packages/web-query/src/index.ts`. A formatting defect there is
  a defect on **every** screen that imports it — check the call-site count before deciding
  where the fix belongs.
- Shared components live in `packages/ui/src/`. Duplicate components across `ui` and
  `web-query` have existed; grep by component name across both before adding one.

## Platform packages

`apps/server_core/internal/platform/`: `apierror`, `archguard`, `config`, `httpx`, `logging`,
`migrate`, `msdb`, `pgdb`, `rawkeys`.

- `apierror/` owns the error envelope (`error.code` / `error.message` / `error.details`) that
  the SDK's typed error surface mirrors. New error codes are a contract change.
- `pgdb/` owns pool creation and tenant helpers.
- `migrate/runner_test.go` **hardcodes the canonical migration count** (79 at time of
  writing). Any plan adding a migration bumps that fixture in the same task, or the unit lane
  goes red for a reason unrelated to the feature.

## Sync and scheduling seam

`internal/modules/sync/` is the scheduling seam, and it is a seam — not a utility.

- `sync/domain/entity.go` declares the valid entities: `products`, `listings`, `orders`,
  `market`, `tariffs`. An entity being declared does **not** mean a job is registered for it.
- `sync/application/scheduler.go` provides `JobFunc`, `RegisterJob`, `RunOnce`, `Start`,
  and the `RecordSuccess` / `RecordFailure` cursor writes.
- `sync_state` (migration `0075_sync_sync_state.sql`) is keyed
  `(tenant_id, installation_id, entity)` and carries `cursor`, `schedule`,
  `last_full_sync_at`, `last_incremental_at`, `last_error`, `consecutive_failures`.
  **There is no `last_success_at` column** — a plan that names one costs a round.

Registering a job on this scheduler inherits `sync_state` reconciliation, per-entity failure
isolation, cursor persistence, and the `RecordFailure` → sync-health → operator-visible card
path. Building a second ticker beside it throws all four away, and the fourth is the one that
turns a broken job into a green screen. **Prefer registration; a parallel scheduler needs a
written justification.**

Known live gap: `Scheduler.Start` is a plain ticker with no boot run and no persisted
due-time, so a long interval (24h) starves after every restart. Harmless at short intervals.
Registered as a harness debt — cite it, don't rediscover it.

## Provider seam and rate budget

Mercado Livre lives at
`internal/modules/connectors/adapters/mercado_livre/` (note the underscore in the directory).

`resilience_decorator.go` holds `defaultRateLimitPerMinute = 900` and a shared token bucket;
exhaustion surfaces as `domain.CapabilityError{Code: ErrCodeProviderRateLimited}`.

The bucket is **shared across all callers**. Any plan that adds a provider-calling job states
its budget explicitly: calls per item × items per cycle × cycles per hour, and what it starves
if it saturates. A blind sweep of the vendable catalogue (≈2 900 products, up to 6 calls each)
is ≈17 500 calls — it saturates the bucket for ~20 minutes and starves listings and orders
sync. That arithmetic belongs in the plan, not in the incident afterwards.

Oracle/ERP reads go through `internal_read` ports with adapters under `adapters/oracle`.
Oracle is **read-only**; no plan writes to it.

## Governance registries

`contracts/governance/`:

| File | What it constrains |
|---|---|
| `modules.json` | module roots, composition requirement, OpenAPI prefixes, allowed dependencies |
| `invariants.json` | the 10 machine checks and their reason codes |
| `shared-seams.json` | exclusive paths that only one writer may hold at a time |
| `execution-lanes.json`, `runtime-config.json`, `knowledge-routes.json`, `harness-evals.json` | lane and routing config |

The invariants, by reason code: `GOV_MODULE_COVERAGE`, `GOV_MODULE_DEPENDENCY`,
`GOV_MODULE_LAYER`, `GOV_COMPOSITION_MISSING`, `GOV_APPLICATION_IMPORT`,
`GOV_POSTGRES_DRIVER`, `GOV_PRODUCTION_PANIC`, `GOV_API_SDK_SPLIT`, `GOV_FRONTEND_FETCH`,
`GOV_MIGRATION_PREFIX`.

Shared seams with exclusive paths — a plan that touches one must say who owns it for the
duration: the api-sdk pair, `apps/server_core/migrations`, `composition/root.go`,
the dependency graph (`go.mod`, `package.json`, `package-lock.json`),
`docs/architecture/decisions`, harness control files, and the provider-capability contract.

`composition/root.go` is ~940 lines and is an exclusive seam. Two tracks editing it in
parallel is a guaranteed merge conflict; plan the wiring task as one owner.

## Verification lanes

From the repo root, unless stated otherwise:

```bash
npm run harness:unit
```

```bash
npm run harness:integration
```

```bash
npm run harness:governance -- -BaseSha <full-40-hex-sha>
```

```bash
npm run docker:dev
```

Lane facts a plan must respect:

- **Go builds/tests bind an absolute `GOCACHE`**, run from `apps/server_core`. A relative
  `.gocache` breaks when the working directory shifts mid-pipeline.
- **The FE vitest lane is `cd apps/web && npx --no-install vitest run`.** The `cd` is part of
  the lane: from the repo root the same command globs a different project set and turns
  contract tests red for reasons unrelated to the code.
- **The governance lane runs from a clean worktree and needs the full 40-hex base SHA.** A
  short SHA yields `GOV_SEMANTIC_DRIFT id=base-sha-invalid`.
- **The integration lane self-provisions Postgres** (fresh `mpc_test_<32hex>` per run) and
  self-discovers packages by `//go:build integration`. A new integration package joins by
  existing; the build tag must be in the first lines of the file.
- **The dev stack is a hub-owned seam**: `docker compose` only, server `:8080`, web `:5174`.
  A plan never instructs an executor to boot its own server, bind those ports, or load `.env`
  into session environment variables.
- **Never dump an environment.** Diagnose one variable by name (`printenv THE_VAR`), never
  `docker inspect` or bare `printenv`.

Known pre-existing reds are allowlisted in HARNESS-PROFILE §2 — cite the entry rather than
re-proving non-linkage. `TestListingsReadContractEndToEnd` is the live one.

## Where documents live

| Document | Path |
|---|---|
| Architecture + frozen decisions | `ARCHITECTURE.md` |
| ADRs | `docs/architecture/decisions/` (004, 005, 006, 007, 008) |
| Repo harness bindings | `docs/HARNESS-PROFILE.md` |
| Designs / specs | `docs/superpowers/specs/YYYY-MM-DD-<topic>-design.md` |
| Plans | `docs/superpowers/plans/YYYY-MM-DD-<feature>-plan.md` |
| Missions and evidence | `.mnfs/MIS-*/` |
| Harness debts | `.mnfs/HARNESS-DEBTS.md` |
| Marketplace API reference | `docs/marketplaces/` |

Debts follow the format **measured fact → cost paid → candidate fix**. Nothing goes in as an
opinion without the case that paid for it.
