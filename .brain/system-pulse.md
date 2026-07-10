# System Pulse - Marketplace Central
> Last updated: 2026-07-09 | Session: #21

## Project Identity

**Name:** Marketplace Central (MPC)  
**Purpose:** Internal Mercado Livre operations and intelligence cockpit (stock, pricing, integrations, orders, profitability, alerts).
**Target:** Mercado Livre operations with Oracle-backed ERP data as internal truth.
**Future:** Designed to merge into MetalShopping as a module.

---

## Technology Stack

### Backend
- Language: Go 1.25.x
- Database: PostgreSQL via `pgx/v5` (`pgxpool.Pool` only)
- Pattern: Modular monolith, ports and adapters

### Frontend
- React + Vite
- SDK: `packages/sdk-runtime`
- Testing: Vitest

### Monorepo
- Go workspace: `go.work` (`apps/server_core`)
- Node workspace: npm workspaces (`apps/web`, `packages/*`)

---

## Active Architectural Decisions

- ADR-001: superseded; direct MetalShopping Postgres internal reads are no longer target architecture
- ADR-002: superseded; `MS_DATABASE_URL` same-cluster assumption is no longer target architecture
- ADR-003: Split integrations delivery into operational specs after foundation (OAuth -> Fee Sync -> UX)
- ADR-004: Integration catalog plugin framework (provider/auth/sync self-registration)
- ADR-005: Mercado Livre first control plane; VTEX is legacy and receives no new work
- ADR-006: MPC reads internal ERP facts through Oracle adapters inside `apps/server_core`

---

## Architecture Overview

```
apps/server_core/
  cmd/server/
  cmd/migrate/
  internal/
    composition/
    modules/
      catalog/
      classifications/
      marketplaces/
      integrations/
      pricing/
      connectors/
      product_links/   # active: internal product <-> Mercado Livre listing links with candidate + audit workflow
      internal_read/   # active: Oracle ERP read contracts + adapters
      inventory/       # planned: internal ERP stock vs ML announced stock reconciliation
      orders/          # planned: Mercado Livre order ingestion and monitoring
      profitability/   # planned: per-order/per-item margin and data quality
    platform/
      config/
      httpx/
      logging/
      pgdb/
  migrations/  # 0001-0018

apps/web/
packages/sdk-runtime/
contracts/api/
```

---

## Module Responsibilities

| Module | Status | Scope |
|---|---|---|
| `catalog` | active | Products, taxonomy, enrichments |
| `classifications` | active | Product classification management |
| `marketplaces` | active consumer | Accounts/policies and bridge to integrations |
| `integrations` | foundation complete | Provider catalog, installations, credentials, auth sessions, capability states, operation runs, base APIs/SDK |
| `pricing` | active | Simulation engine and batch orchestration |
| `connectors` | active foundation | Provider-specific integration surfaces; Mercado Livre first, Melhor Envio auth active, VTEX removed from runtime |
| `product_links` | active foundation | Internal product/SKU to Mercado Livre listing/variation mapping with live snapshot import, candidate generation, and audited operator resolution workflow |
| `internal_read` | active reset target | MPC-owned Oracle read contracts and adapters for internal facts |
| `inventory` | active foundation | Stock policy, risk classification, manual action audit, API/SDK, and Stock Seguro dashboard validated for Mercado Livre |
| `orders` | planned | Mercado Livre order ingestion, shipment/status tracking, cancellation reasons |
| `profitability` | planned | Per-order/per-item margin from ML revenue/fees/freight and internal ERP cost/tax inputs |

---

## Established Patterns

- Module structure: `domain/ application/ ports/ adapters/ transport/ events/ readmodel/`
- Structured errors: `MODULE_ENTITY_REASON`
- Handler logs: `action`, `result`, `duration_ms`
- Tenant-safe access: all business queries scoped by `tenant_id`
- Frontend must use `sdk-runtime` (no direct backend fetch in features)
- Money: `float64` in domain, `numeric(14,2)` in Postgres
- Runtime/contract/wiki/execution truth must be classified when they disagree
- Evidence before closure; no task is done without verification and bounded defers
- Global-maximum design beats extending legacy paths

---

## Current Phase

**Phase 6 - Orders + Margin ML** (execution)

Recently completed phase:
1. `T-031` Complete provider catalog definitions
2. `T-032` Redesign marketplace catalog and integration UX
3. `T-033` Align OpenAPI + `sdk-runtime` provider metadata contract
4. `T-034` Close catalog foundation with backend/frontend verification and evidence ledger

Next in sequence:
1. Start M-06/F-01 order ingestion on top of passed M-05 Stock Seguro and M-04 product-link truth
2. Define margin input quality states for fee, freight, cost, tax, and manual adjustments without zero-default fallbacks
3. Research Mercado Livre orders/shipments/payment semantics against official docs before profitability calculations
4. Keep provider stock writes explicitly manual/audited and outside automated flows

---

## Recent Changes

- 2026-04-27: Added Amazon draft-mode sandbox fallback in auth adapter (allow `amzn1.sp.solution.*` only when `MPC_PROVIDER_AMAZON_AUTH_VERSION=beta`) to unblock authorize flow before production app registration
- 2026-04-27: Improved integration error surfacing: backend now preserves wrapped configuration context and frontend shows actionable marketplace action errors
- 2026-04-27: Aligned provider panel action gating with backend status rules (fee sync disabled unless installation is `connected` or `degraded`) and validated via browser checks on `/marketplaces`
- 2026-04-25: Completed `T-031` through `T-034` with commits `b922940`, `1540a8a`, and `d915806`; added six-provider catalog baseline and provider-first marketplace UX
- 2026-04-25: Added rollout evidence ledger at `docs/superpowers/evidence/2026-04-25-marketplace-catalog-ux-data-foundation.md` with passing backend/frontend/build verification
- 2026-07-06: Accepted ADR-005 to remove VTEX from the target architecture and pivot Marketplace Central to a Mercado Livre first internal cockpit backed by internal ERP data
- 2026-07-06: Completed M-01/F-01 VTEX surface inventory and M-01/F-02 active-surface removal; `/connectors/vtex/*`, VTEX SDK/runtime/frontend surfaces, and VTEX route wiring were removed from active code
- 2026-07-06: M-01 validation passed; MIS-001 and M-02 moved into execution with an M-02 execution guide dispatching F-01 capability port contract first
- 2026-07-08: Completed the marketplace integration foundation slice with connection snapshot projection, runtime capability honesty, credential-safe provider execution, and API/SDK/UI alignment for integrations
- 2026-07-08: Real Mercado Livre platform reads/probes are now validated for `account_probe`, `listing_read`, `order_read`, `fee_quote_read`, and `stock_read`; validation remains read-only and does not cover provider writes
- 2026-07-08: Completed M-04/F-01 listing snapshot import with a dedicated `product_links` module slice, manual import API/SDK contract, tenant-scoped Postgres upsert, and live Mercado Livre validation proving idempotent re-import
- 2026-07-08: Implemented M-04/F-02 candidate generation with exact-first rules, persisted candidate storage, and API/SDK contract; positive Oracle-backed live validation is currently blocked because `10.55.10.101:1521` is unreachable from this session, but degraded-runtime `503 PRODUCT_LINKS_INTERNAL_READ_UNAVAILABLE` behavior was validated
- 2026-07-09: Completed M-04/F-03 link resolution workflow with dedicated `feature-product-links` UI, workflow API/SDK contract, persisted `product_links` + `product_link_audit_entries`, and live approve/reject/manual-resolve validation against the Mercado Livre Docker stack
- 2026-07-09: M-04 Product Links milestone passed; the next execution target is M-05 Stock Seguro, where downstream stock actions must consume the explicit `resolved` / `rejected` / `conflict` / `unresolved` link states as safety gates
- 2026-07-09: M-05 Stock Seguro milestone passed; real QA found and fixed Vite proxy shadowing, inventory JSON contract drift, and blocker-summary miscount before closeout
- 2026-07-09: Browser and real API validation now prove Stock Seguro reads real Mercado Livre listing snapshots, shows blocked states with policy/source evidence, and exposes inventory-owned risk/action contracts through API and SDK
- 2026-07-07: Superseded the `MS_DATABASE_URL` / MetalShopping direct-read premise; Oracle becomes the internal truth boundary through MPC-owned adapters inside `apps/server_core`
- 2026-07-07: M-03 passed with real Oracle validation on Windows; `internal_read` now owns live-validated product, stock, price, cost, sales, and tax reads via `godror` + Oracle Instant Client + MPC-owned Oracle adapters

---

## Key File Locations

| File | Purpose |
|---|---|
| `AGENTS.md` | Engineering rules and guardrails |
| `ARCHITECTURE.md` | Frozen architecture decisions |
| `IMPLEMENTATION_PLAN.md` | Top-level phased execution plan |
| `contracts/api/marketplace-central.openapi.yaml` | API source of truth |
| `apps/server_core/internal/composition/root.go` | DI and module wiring |
| `apps/server_core/internal/modules/integrations/` | Integrations platform module |
| `apps/server_core/internal/modules/integrations/adapters/amazon/auth_adapter.go` | Amazon LWA/SP-API auth URL + token exchange behavior |
| `apps/server_core/internal/modules/integrations/adapters/providers/registry.go` | Integration plugin registration and adapter/syncer assembly |
| `packages/sdk-runtime/src/index.ts` | Typed client methods |
| `wiki/modules/integrations.md` | Amazon auth identifier mapping and runtime troubleshooting notes |
| `.brain/decisions/004-integration-catalog-plugin-framework.md` | Integration plugin framework ADR |
| `.brain/decisions/005-mercado-livre-first-control-plane.md` | Mercado Livre first pivot and VTEX retirement ADR |

---

## Known Risks

- `.brain/` remains gitignored by default (project memory can diverge across machines)
- `.agents/skills/nexus-*` are currently untracked local project tooling
- Migration runner `cmd/migrate/main.go` still needs production-hardening workflow
- OAuth callback success evidence for external provider consent remains dependent on sandbox availability (tracked as `DONE_WITH_CONCERNS` in `T-029`)
- Windows environments may require local absolute `GOCACHE` for stable test runs
- Frontend tests include non-blocking React `act(...)` warnings for marketplace loading-state test setup
- Amazon production consent still blocked until real `amzn1.sellerapps.app.*` application ID is available; current `beta` fallback is temporary and should not be kept for production release
- Legacy VTEX residue is now limited to forward-only migrations and historical docs that still need explicit closeout classification in M-01/F-03
- Old `MS_DATABASE_URL` and same-cluster assumptions still exist in historical docs/artifacts and must not be reused as execution truth
- Mercado Livre stock writes can create oversell/undersell risk unless guarded by product links, safety buffers, idempotency, and audit
