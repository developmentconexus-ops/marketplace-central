# AGENTS — Marketplace Central

## On every session start

1. Read `ARCHITECTURE.md` — canonical architecture contract and frozen decisions
2. Read `wiki/README.md` — index for module docs, architecture notes, and operational references
3. Read `.brain/system-pulse.md` and `.brain/roadmap.json` — current project state and execution plan
4. Read `IMPLEMENTATION_PLAN.md` only when validating historical phase scope or reconciling older roadmap notes
5. After any correction: document the lesson in commit message or PR description

## Source of truth

- `ARCHITECTURE.md` is the stable architecture source of truth
- `.brain/` is Nexus working memory: current status, roadmap, session continuity, and ADR index
- `wiki/README.md` is the human knowledge index for modules, operations, and implementation context
- `contracts/api/marketplace-central.openapi.yaml` is the HTTP API contract
- `IMPLEMENTATION_PLAN.md` is legacy/top-level phase history, not the primary day-to-day plan

When architecture changes, update `ARCHITECTURE.md`, the relevant wiki page, and Nexus/ADR records in the same task.

## Behavioral guardrails (LLM)

These guidelines reduce common coding mistakes. They favor caution over speed; for trivial tasks, use judgment.

### 1) Think before coding

- State assumptions explicitly before implementation
- If uncertain, ask instead of guessing
- If multiple interpretations exist, present options explicitly
- Call out simpler approaches when they exist
- If scope or intent is unclear, pause and clarify

### 2) Simplicity first

- Implement the minimum code that solves the requested problem
- Do not add features, configurability, or abstractions that were not requested
- Avoid defensive handling for impossible scenarios
- If a solution can be much smaller, rewrite it smaller
- Sanity check: if a senior reviewer would call it overcomplicated, simplify

### 3) Surgical changes

- Touch only what is needed for the request
- Do not refactor or reformat adjacent unrelated code
- Match existing local style and patterns
- If unrelated dead code is found, mention it but do not remove it unless asked
- Remove only unused code/imports created by your own change
- Every changed line must trace directly to the user request

### 4) Goal-driven execution

- Turn requests into verifiable outcomes
- For bug fixes: reproduce first (test or clear check), then fix, then verify
- For validations: write failing checks/tests first when practical, then make them pass
- For refactors: verify behavior parity before and after
- For multi-step work, define concise step → verification pairs before coding

These guidelines are working when diffs are smaller, rewrites are fewer, and clarifications happen before implementation.

## Engineering bar

Every decision passes this filter:
*"Would a Stripe or Google senior engineer approve this in code review?"*

- Names are self-documenting — no comment needed to understand them
- Errors carry structured codes: `MODULE_ENTITY_REASON` (e.g. `PRICING_SIMULATION_INVALID`)
- Every handler logs `action`, `result`, `duration_ms`
- Every write is idempotent and retry-safe
- Tenant-owned business tables carry `tenant_id`; global/system reference tables must be explicitly documented

## Absolute rules — violation = stop and fix immediately

### Go

- Every tenant-owned Postgres query must scope by `tenant_id`
- Global/system reference queries are allowed only when the table is intentionally shared, documented, and protected by schema or service rules
- Every HTTP handler validates request method and returns structured JSON errors
- Every new module is registered in `composition/root.go` with dependency injection
- Integration providers self-register definitions/auth/sync factories; composition imports provider packages only to activate registration
- Transport layer never contains business logic — delegate to application service
- Application service never imports `net/http` or database packages — use ports
- Domain entities are pure Go structs with no external dependencies
- Adapters implement port interfaces — one adapter per external dependency
- `pgxpool.Pool` is the only database access mechanism — no raw `sql.DB`
- No `panic()` in production code — return errors
- All monetary values use `float64` in domain, `numeric(14,2)` in Postgres until an ADR changes this convention
- Use local Go cache for test/build commands: `GOCACHE=.gocache`

### Frontend

- Data only via `sdk-runtime` methods — no direct `fetch()` to backend
- No business logic in React components — pricing, margin, commission calculations belong in Go
- No local persistence (localStorage, SQLite) as source of truth — Postgres is canonical
- Loading + error + empty state on every data-fetching component
- Feature packages (`packages/feature-*`) own page-level UI
- Shared primitives live in `packages/ui`

### Process

- No code task marked done without appropriate verification, passing impacted tests/builds, and a commit
- Docs-only tasks require at least a diff review/proofread and a commit
- One commit per completed task — no uncommitted work at session end
- Legacy files from the old Next.js monolith must not be reintroduced
- Every new endpoint must exist in `contracts/api/marketplace-central.openapi.yaml`
- Every new migration file is sequential: `NNNN_description.sql`
- `packages/generated/` never edited manually (when SDK generation is added)

## Module structure

Every module in `apps/server_core/internal/modules/*` must follow:

```
domain/        — entities and value objects (pure Go, no imports)
application/   — use cases and service layer (imports domain + ports)
ports/         — interfaces for external dependencies
adapters/      — implementations of ports (postgres, http clients, etc.)
transport/     — HTTP handlers (imports application + platform/httpx)
events/        — event types for async communication (future)
readmodel/     — query-optimized views (future)
```

## Marketplace and integration plugin patterns

There are two related layers:

- `marketplaces`: business configuration, account/policy links, marketplace definitions, and fee schedules
- `integrations`: technical operations, provider catalog, installations, credentials, auth state, capability state, and operation runs

Integration providers use self-registration:

```go
func init() {
    providers.RegisterDefinition(definition)
    providers.RegisterAuthFactory(providerCode, factory)
    providers.RegisterFeeSyncerFactory(providerCode, factory)
}
```

Adding a provider should normally mean adding a provider package, plus one composition side-effect import if the package is not already imported.

Runtime marketplace operations live behind connector/capability ports. A marketplace connector may implement only the capabilities it supports:

```go
type MarketplaceConnector interface {
    FetchMessages(ctx context.Context) ([]Message, error)
    FetchOrders(ctx context.Context) ([]Order, error)
    ReplyToMessage(ctx context.Context, messageID string, body string) error
}
```

One adapter per marketplace/provider (`vtex`, `mercado_livre`, `magalu`, etc.). The owning module defines the port; adapters live under `adapters/` and must not own tenant business state directly.

## Commit format

`<type>(<scope>): <what>` — feat | fix | docs | chore | refactor | test

Examples:
- `feat(pricing): add margin threshold alerts`
- `fix(connectors): handle VTEX token refresh on 401`
- `docs(architecture): freeze messaging module scope`

## Skill map

| Task | Reference |
|---|---|
| Any Go implementation | This file + `ARCHITECTURE.md` + relevant wiki page |
| Code/module knowledge | `wiki/README.md` (index to architecture, modules, operations) |
| Database changes | `apps/server_core/migrations/` |
| API contract changes | `contracts/api/marketplace-central.openapi.yaml` |
| Frontend feature | `packages/feature-*/` + `packages/sdk-runtime/` |
| Current planning/status | `.brain/roadmap.json` + `.brain/system-pulse.md` |
| Historical phase context | `IMPLEMENTATION_PLAN.md` |

## Integration with MetalShopping

This repository is designed as a future module of MetalShopping. Key compatibility rules:

- Module structure mirrors `MetalShopping_Final/apps/server_core/internal/modules/*`
- Platform packages mirror `MetalShopping_Final/apps/server_core/internal/platform/*`
- Database schema uses prefix `mpc_` or dedicated tables (no collision with MS tables)
- When the merge happens, modules move to MetalShopping's monorepo with minimal rewrite
- The same Postgres cluster can be shared (different schema or same schema with table prefixes)
