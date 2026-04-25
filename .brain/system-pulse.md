# System Pulse - Marketplace Central
> Last updated: 2026-04-25 | Session: #18

## Project Identity

**Name:** Marketplace Central (MPC)  
**Purpose:** Intelligence and control surface for marketplace operations (pricing, integrations, messaging, orders, SLA).  
**Target:** Brazilian marketplace sellers using VTEX and marketplace channels.  
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

- ADR-001: MPC reads products directly from MetalShopping Postgres
- ADR-002: MPC own tables live in `mpc` schema on MetalShopping cluster
- ADR-003: Split integrations delivery into operational specs after foundation (OAuth -> Fee Sync -> UX)
- ADR-004: Integration catalog plugin framework (provider/auth/sync self-registration)

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
| `connectors` | partial | Provider-specific integration surfaces (legacy/transitional + VTEX/ME flows) |

---

## Established Patterns

- Module structure: `domain/ application/ ports/ adapters/ transport/ events/ readmodel/`
- Structured errors: `MODULE_ENTITY_REASON`
- Handler logs: `action`, `result`, `duration_ms`
- Tenant-safe access: all business queries scoped by `tenant_id`
- Frontend must use `sdk-runtime` (no direct backend fetch in features)
- Money: `float64` in domain, `numeric(14,2)` in Postgres

---

## Current Phase

**Phase 4 - VTEX connector** (in progress)

Recently completed phase:
1. `T-031` Complete provider catalog definitions
2. `T-032` Redesign marketplace catalog and integration UX
3. `T-033` Align OpenAPI + `sdk-runtime` provider metadata contract
4. `T-034` Close catalog foundation with backend/frontend verification and evidence ledger

Next in sequence:
1. Start VTEX connector infrastructure (`4.1`)
2. Continue catalog enhancement for VTEX operational flows (`4.2`)
3. Expose VTEX connector UX in frontend (`4.3`)

---

## Recent Changes

- 2026-04-25: Completed `T-031` through `T-034` with commits `b922940`, `1540a8a`, and `d915806`; added six-provider catalog baseline and provider-first marketplace UX
- 2026-04-25: Added rollout evidence ledger at `docs/superpowers/evidence/2026-04-25-marketplace-catalog-ux-data-foundation.md` with passing backend/frontend/build verification
- 2026-04-25: Restructured `phase-3b` into plugin-framework-driven tasks `T-031` through `T-034`
- 2026-04-25: Replaced manual integration wiring with registration-backed provider definitions, auth factories, and fee syncer factories; `root.go` now consumes registries
- 2026-04-25: Completed roadmap task `T-030` and added ADR-004 to document the integration catalog plugin framework

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
| `apps/server_core/internal/modules/integrations/adapters/providers/registry.go` | Integration plugin registration and adapter/syncer assembly |
| `packages/sdk-runtime/src/index.ts` | Typed client methods |
| `.brain/decisions/004-integration-catalog-plugin-framework.md` | Integration plugin framework ADR |

---

## Known Risks

- `.brain/` remains gitignored by default (project memory can diverge across machines)
- `.agents/skills/nexus-*` are currently untracked local project tooling
- Migration runner `cmd/migrate/main.go` still needs production-hardening workflow
- OAuth callback success evidence for external provider consent remains dependent on sandbox availability (tracked as `DONE_WITH_CONCERNS` in `T-029`)
- Windows environments may require local absolute `GOCACHE` for stable test runs
- Frontend tests include non-blocking React `act(...)` warnings for marketplace loading-state test setup
