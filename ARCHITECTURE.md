# Marketplace Central Architecture

## Status

This architecture is approved as the official foundation of the repository. Frozen decisions below must not be rediscussed without an explicit ADR.

Accepted architecture decisions are indexed in [`docs/architecture/decisions/`](docs/architecture/decisions/README.md).

## Reference baseline

This repository mirrors the structural discipline of MetalShopping Final and adopts the operating rigor proven in MetalDocs:

- GitHub: https://github.com/leandrotcawork/MetalShopping_Final
- The engineering bar, module pattern, and platform conventions are inherited from MetalShopping
- Runtime truth, contract truth, wiki truth, and verification truth must stay aligned
- Architecture contradictions stop local feature work until classified and resolved

## North Star

Marketplace Central is an internal Mercado Livre operations and intelligence cockpit:

- Independent monorepo, architecturally compatible with future MetalShopping integration
- Server-first: all business logic lives in Go
- Modular monolith core with explicit module boundaries
- PostgreSQL as the only canonical state
- Thin web client consuming SDK-generated methods
- Oracle-backed ERP data is the internal source of truth for products, stock, price, cost, taxes, and sales history
- Mercado Livre is the first operational marketplace control plane for listings, stock, price, orders, fees, and questions
- MPC handles intelligence and guardrails: stock reconciliation, price simulation, order profitability, alerts, and action audit

## Frozen decisions

1. Independent monorepo — merge into MetalShopping is a future module migration, not a rewrite
2. Go `apps/server_core` is the canonical backend
3. `apps/web` is a thin React client — no business logic
4. PostgreSQL is the only canonical persistence — no SQLite, no localStorage as source of truth
5. Single-tenant operation, but every business table carries `tenant_id` (tenant-ready by design)
6. Stable API routes without `/v1` prefix in URLs — versioning is in the OpenAPI document
7. External marketplace integrations enter only through the `connectors` module via port interfaces
8. Mercado Livre is the first operational control plane; VTEX is legacy and must not receive new feature work
9. Scheduler-based polling is acceptable initially; webhook/notification support may be added where Mercado Livre provides reliable event topics
10. Frontend consumes only `packages/sdk-runtime` — never calls backend directly
11. Global-maximum design beats local patches: do not extend legacy VTEX abstractions to solve Mercado Livre problems
12. Internal ERP reads happen through MPC-owned ports implemented by Oracle adapters inside `apps/server_core`; no legacy `MS_DATABASE_URL` read path remains in the target architecture

## Layout

```
apps/
  server_core/          # Go backend — canonical business logic
    cmd/
      server/           # HTTP server entrypoint
      migrate/          # Migration runner
    internal/
      composition/      # Module registration and dependency injection
      modules/
        catalog/        # Product entities for pricing simulation
        marketplaces/   # Marketplace accounts and pricing policies
        pricing/        # Price simulation engine
        messaging/      # [planned] Centralized customer messages from all marketplaces
        orders/         # [implemented; M-06 milestone blocked] Order tracking with SLA monitoring
        alerts/         # [planned] SLA guardrails and notifications
        connectors/     # Marketplace API adapters (Mercado Livre first; legacy adapters inventoried before deletion)
        integrations/   # Integration lifecycle (install/auth/credential/fee-sync operations)
      platform/
        config/         # Environment configuration
        httpx/          # HTTP helpers (JSON writer, router)
        logging/        # Structured logger
        pgdb/           # PostgreSQL pool, tenant helpers
    migrations/         # Sequential SQL migrations
    tests/
      unit/             # Unit tests with stub repositories

  web/                  # Thin React client (Vite + React Router)
    src/
      app/              # Route definitions
      main.tsx          # Entry point

contracts/
  api/                  # OpenAPI spec — source of truth for HTTP behavior
  events/               # [reserved] Async event contracts
  governance/           # [reserved] Runtime governance schemas

packages/
  sdk-runtime/          # TypeScript client for web-to-core communication
  ui/                   # Shared UI primitives (Button, SurfaceCard, etc.)
  feature-marketplaces/ # Marketplace settings page
  feature-simulator/    # Pricing simulator page

docs/
  marketplaces/         # Per-marketplace API reference docs (ML, Magalu, Amazon, etc.)
  IMPLEMENTATION_PLAN.md
```

## Module responsibilities

### `catalog` (implemented — foundation)

Product entities used by pricing, product linking, and future Mercado Livre listing operations.

Scope: product CRUD, SKU/EAN management, cost tracking.

### `marketplaces` (implemented — foundation)

Marketplace account configuration and pricing policies (commission, fixed fees, freight, SLA thresholds).

Scope: account registration, policy management, connection status tracking.

### `pricing` (implemented — foundation)

Price simulation engine. Calculates margin, commission impact, freight cost, and viability per product per marketplace.

Scope: simulation execution, snapshot persistence, manual price overrides, margin alerts.

### `product_links` (active validated foundation — Mercado Livre first)

Maps internal products/SKUs to Mercado Livre listing and variation identifiers.

Scope: link creation, confidence/state tracking, duplicate detection, audit trail.

### `internal_read` (active foundation — Oracle first)

Owns MPC's read contracts and Oracle adapters for internal ERP facts.

Scope: product, stock, price, cost, tax, and sales-history reads through MPC-owned ports and typed domain models.

Rules:
- Application/domain code depends on MPC-owned read contracts, never on Oracle SQL or driver types.
- Oracle mapping and query semantics live only in `adapters/oracle`.
- Missing or ambiguous source facts surface as explicit quality states, never silent zero/default values.
- Read access is global-maximum and contract-first: no ad hoc SQL from downstream modules.

### `inventory` (active validated foundation — Mercado Livre first)

Compares internal ERP stock with Mercado Livre announced stock and proposes or applies safe stock actions.

Scope: stock snapshots, safety buffers, divergence detection, manual approval, action audit.

Reads from: internal stock views via ports.
Writes to: Mercado Livre only through connector capabilities after policy checks.

### `orders` (implemented — M-06 milestone blocked)

Order monitoring and reconciliation for Mercado Livre. Tracks order lifecycle, items, fees, shipping, cancellation reasons, and internal product links.

Scope: order polling/notifications, status tracking, cancellation analysis, dispatch guardrails.

Reads from: Mercado Livre APIs via `connectors` adapters and internal product/cost providers.

### `profitability` (implemented — M-06 milestone blocked)

Calculates per-order and per-item contribution using Mercado Livre revenue/fees/freight and internal ERP cost/tax inputs.

Scope: margin snapshots, manual cost adjustments, data quality flags, profitability alerts.

### `messaging` (planned)

Centralizes Mercado Livre questions/messages first. Multi-marketplace inbox is deferred until the Mercado Livre workflow is stable.

Scope: message polling from marketplace APIs, unified thread view, reply dispatch, response time tracking.

Read from: marketplace APIs via `connectors` adapters.
Write to: marketplace APIs (reply only) via `connectors` adapters.

### `alerts` (planned — phase 2)

SLA guardrails and notification engine. Monitors messaging response times, order dispatch deadlines, and pricing thresholds.

Scope: deadline calculation, alert generation, notification dispatch (initially in-app, future email/webhook).

Reads from: `messaging` and `orders` modules for SLA data.

### `connectors` (implemented foundation)

Marketplace API adapters. Mercado Livre is the first target for live operations. Legacy VTEX surfaces are not part of the target architecture.
Other marketplaces (Magalu, Amazon, Shopee, etc.) are future, capability-driven additions only after Mercado Livre operations are reliable.

Scope: authentication management, API request/response mapping, rate limiting, error handling.

Pattern: one adapter package per marketplace under `connectors/adapters/`. The module owns the port interfaces; adapters implement them.

Current connector baseline: integrations framework, Mercado Livre OAuth, seeded fee baseline, and live listing/order read capabilities used by validated product-link, inventory, and M-06 evidence. No live provider stock mutation is claimed; stock writes remain policy-gated and were validated without executing a real mutation.

### `integrations` (implemented foundation)

Integration lifecycle module for provider definitions, installation records, auth flows,
credential lifecycle, and fee-sync operation tracking.

Scope: provider registry, installation draft/connection states, OAuth/API-key auth flow,
credential rotation, operation runs, capability state transitions, and scheduled refresh/cleanup jobs.

## Internal ERP Access

Marketplace Central reads internal operational facts directly from Oracle through adapters inside `apps/server_core`.

Design rules:

- Oracle is an external dependency behind module-owned ports, not a second business store.
- `internal_read` owns the read contracts used by `product_links`, `inventory`, `orders`, and `profitability`.
- SQL/query knowledge stays inside Oracle adapters and helper packages owned by that boundary.
- PostgreSQL stores only MPC-owned operational state, audit, projections, and snapshots.
- Removing or changing an Oracle query shape must not force business-module rewrites; only adapter implementations should move.

## Platform packages

Located in `apps/server_core/internal/platform/`:

- `config/` — environment variable loading
- `httpx/` — JSON response writer, router factory, middleware
- `logging/` — structured logger
- `pgdb/` — PostgreSQL pool creation, tenant context helpers

These are shared infrastructure — not business logic. They mirror MetalShopping's `internal/platform/` structure.

## Communication flow

### Current (foundation)

```
web → sdk-runtime → server_core HTTP handlers → application services → postgres
```

### Target (phase 2+)

```
web → sdk-runtime → server_core HTTP handlers → application services → postgres
                                                       ↓
                                              scheduler jobs / webhooks
                                                       ↓
                                              connectors adapters → Mercado Livre APIs
                                                       ↓
                                              inventory/orders/profitability/messaging modules → postgres
                                                       ↓
                                              alerts module → notifications
```

### Rules

- Web client never calls marketplace APIs directly
- Connectors never own business state — they fetch and deliver to domain modules
- Scheduler runs polling jobs at configured intervals; Mercado Livre notifications may later reduce polling where reliable
- Synchronous HTTP requests from web never depend on connector availability

## Database

- Engine: PostgreSQL for MPC-owned state
- All business tables carry `tenant_id` as part of the primary key or with NOT NULL constraint
- Migrations are sequential files in `apps/server_core/migrations/`
- Naming: `NNNN_description.sql`
- No down migrations — forward-only

## External Data Sources

- Oracle ERP: source of truth for internal product, stock, price, cost, tax, and sales inputs, consumed through `internal_read` ports/adapters
- Mercado Livre APIs: source of truth for marketplace listing, order, fee, and question state, consumed through `connectors` adapters

Legacy note:
- `MS_DATABASE_URL` and direct MetalShopping/Postgres internal-read assumptions are no longer part of the target architecture.

## Future MetalShopping integration

When the time comes to merge MPC into MetalShopping:

1. Move `apps/server_core/internal/modules/*` into MetalShopping's module directory
2. Register modules in MetalShopping's composition root
3. Migrate database tables (add to MetalShopping's migration sequence)
4. Move `packages/feature-*` into MetalShopping's frontend packages
5. Point SDK methods to MetalShopping's API routes

The merge should be a module migration, not a rewrite. This is why structure compatibility matters now.

## Related documents

- `AGENTS.md` — daily operational rules
- `docs/architecture/decisions/README.md` — accepted architecture decisions
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/mission.md` — current execution truth
- `IMPLEMENTATION_PLAN.md` — historical reconciliation only
- `contracts/api/marketplace-central.openapi.yaml` — API source of truth
- `docs/marketplaces/*.md` — per-marketplace API reference


