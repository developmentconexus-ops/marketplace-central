# M-03 Internal Read Design

**Date:** 2026-07-07
**Mission:** `MIS-001`
**Milestone:** `M-03-mnos-sankhya-read-contract`
**Status:** Drafted after design approval

## Goal

Introduce an MPC-owned internal read seam for Sankhya/MNOS-derived product, stock, price, cost, tax, and sales facts so future `product_links`, `inventory`, and `profitability` modules can consume stable contracts without ad hoc SQL, Sankhya writes, or ERP mirroring.

## Problem

MPC has mission-level contract research for MNOS/Sankhya semantics, but no production code surface that owns those semantics. Planned modules need a shared seam for:

- product candidate lookup for linking;
- sellable stock calculation with explicit scope and source metadata;
- cost/tax lookup with missing-value quality flags;
- future price and sales reads.

Without an MPC-owned seam, later milestones would either duplicate contract logic, leak infrastructure details into business modules, or extend unrelated modules such as `catalog` or `connectors`.

## Constraints

- Preserve MNOS semantics from `.mnfs/MIS-001-mercado-livre-operating-cockpit/research/mnos-sankhya-read-interface-contract.md`.
- Default sellable stock remains `SUM(ESTOQUE - RESERVADO)` with `CODEMP IN (1,2)` and `CODLOCAL=10101`.
- `CODLOCAL=10108` showroom stock must not contribute to default sellable stock.
- Initial margin cost basis is `CUSSEMICM`, never `CUSVARIAVEL`.
- Missing product, stock, cost, tax, or ambiguous data must become explicit quality flags, never zero defaults.
- Sankhya access remains read-only.
- MPC must not mirror whole ERP tables; only MPC-owned state or auditable snapshots are allowed.
- Secrets and DSN values must not appear in logs, errors, or validation artifacts.

## Approaches Considered

### Option 1: New shared `internal_read` module

Create `apps/server_core/internal/modules/internal_read/` with `domain`, `application`, `ports`, and adapters. Future business modules consume it through focused ports or application services.

Pros:
- ownership is explicit;
- keeps internal source contracts separate from marketplace adapters;
- scales to `inventory`, `product_links`, and `profitability` without overloading `catalog`.

Cons:
- introduces a new module before direct consumers exist.

### Option 2: Extend `catalog`

Put stock/cost/tax/product read semantics into `catalog`.

Pros:
- fewer new directories.

Cons:
- wrong ownership boundary;
- pulls stock/cost/tax/sales into a product module;
- makes later consumers depend on unrelated concepts.

### Option 3: Place the seam in `connectors`

Treat Sankhya like an external provider adapter.

Pros:
- centralizes integration-style code.

Cons:
- contradicts architecture because `connectors` is for marketplace providers;
- blurs internal truth with external marketplace APIs.

## Chosen Design

Adopt **Option 1**.

Create a new shared module:

- `apps/server_core/internal/modules/internal_read/domain/`
- `apps/server_core/internal/modules/internal_read/application/`
- `apps/server_core/internal/modules/internal_read/ports/`
- `apps/server_core/internal/modules/internal_read/adapters/fake/`
- `apps/server_core/internal/modules/internal_read/adapters/oracle/`

This module owns the MPC contract for internal reads. It does not expose public HTTP routes in M-03. It provides application-facing read APIs and fake/test adapters first, while the Oracle adapter remains read-only and safe by construction.

## Module Responsibilities

### Domain

Own stable value objects and enums for:

- internal product candidate;
- sellable stock result and source scope;
- current price result;
- cost-as-of result;
- tax input result;
- sales history result;
- quality flags such as `missing_product`, `missing_stock`, `missing_cost`, `missing_tax`, `ambiguous_product`, `stale_source`.

### Ports

Own application-facing read interfaces for operations aligned with IC-002:

- `FindProductsForLinking`
- `GetSellableStock`
- `GetCurrentPrice`
- `GetCostAsOf`
- `GetSalesHistory`
- `GetTaxInputs`

### Application

Provide a thin service that delegates to the configured reader port, preserves typed errors, and centralizes any invariant checks needed by future consumers. This layer should stay orchestration-only, not SQL-aware.

### Adapters

- `fake/`: deterministic in-memory/test adapter used by unit tests and future module tests.
- `oracle/`: read-only adapter shape and bootstrap guardrails for future live Sankhya access. In M-03 it can remain minimal, but must preserve no-write and no-secret-leak guarantees.

## Boundaries

- Future `inventory`, `product_links`, and `profitability` modules must not import Oracle adapter code directly.
- `catalog` remains owner of current product catalog behavior; it is not promoted into a generic Sankhya read module.
- `connectors` remains owner of marketplace provider capabilities only.
- `tenant_id` scope in MPC must stay distinct from Sankhya `CODEMP` and `CODLOCAL` filters.
- Provider IDs remain strings; Sankhya IDs remain numeric and are not mixed.

## Feature Breakdown

### F-01: Sankhya read contract import

Deliver:

- MPC-owned domain types and quality flags under `internal_read/domain`.
- Port definitions under `internal_read/ports`.
- Feature artifacts citing exact MNOS evidence.

Out of scope:

- live Oracle implementation depth;
- public HTTP routes;
- inventory/product-link business workflows.

### F-02: Read adapter seam

Deliver:

- application service under `internal_read/application`;
- fake adapter for unit tests;
- Oracle adapter seam with read-only and secret-safe shape;
- composition wiring only if needed for non-HTTP consumers.

Out of scope:

- production write paths;
- full live Oracle validation;
- business-module integration beyond seam readiness.

### F-03: Data quality rules

Deliver:

- reusable quality flag semantics in domain;
- tests proving missing values stay flagged, not converted to zero;
- blocked/incomplete examples for stock and margin inputs.

Out of scope:

- UI copy;
- product exclusion policy;
- downstream consumer behavior beyond shared contract semantics.

## Data Flow

1. Future module asks `internal_read` application service for product/stock/cost/tax facts.
2. Application service delegates to configured reader port.
3. Adapter returns typed domain results with source metadata and quality flags.
4. Future module decides how to block, degrade, or continue based on those flags.

No public route is added in this milestone. No direct Sankhya calls from React. No business logic in adapters.

## Error Handling

- Source unavailability maps to structured errors such as `SANKHYA_READ_UNAVAILABLE`.
- Product ambiguity maps to `SANKHYA_PRODUCT_AMBIGUOUS`.
- Missing cost and tax return successful typed results with `missing_cost` or `missing_tax` quality flags instead of zero values.
- Configuration/bootstrap errors may name missing environment keys but must never include secret values.

## Testing Strategy

M-03 evidence should be unit-test-first and fake-adapter-driven.

Required proof:

- sellable stock uses `SUM(ESTOQUE - RESERVADO)` semantics;
- default stock scope uses `CODEMP IN (1,2)` and `CODLOCAL=10101`;
- `CODLOCAL=10108` contributes zero under default policy;
- missing product/cost/tax become typed quality flags;
- missing numeric values are never converted to `0`;
- read seam errors do not leak DSN or password values;
- no write path exists in the Oracle adapter seam.

## File Structure

Expected primary paths:

- `apps/server_core/internal/modules/internal_read/domain/*.go`
- `apps/server_core/internal/modules/internal_read/application/*.go`
- `apps/server_core/internal/modules/internal_read/ports/*.go`
- `apps/server_core/internal/modules/internal_read/adapters/fake/*.go`
- `apps/server_core/internal/modules/internal_read/adapters/oracle/*.go`
- `apps/server_core/internal/composition/root.go`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-01-sankhya-read-contract-import/{spec.md,plan.md,validation.md}`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-02-read-adapter-seam/{spec.md,plan.md,validation.md}`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-03-data-quality-rules/{spec.md,plan.md,validation.md}`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/validation-result.md`

## Risks

- Drift from MNOS semantics if fake fixtures invent unsupported meanings.
- Over-designing the module before real consumers arrive.
- Secret leakage in bootstrap or adapter errors.
- Accidental ownership creep into `catalog` or `connectors`.

## Mitigations

- Cite MNOS source files directly in feature specs and validation.
- Keep M-03 focused on contract/seam/quality only.
- Make fake adapter the primary test harness; keep Oracle adapter minimal.
- Add explicit tests for showroom exclusion, `CUSSEMICM`, and missing-value flags.

## Success Criteria

M-03 is ready to pass when:

- all three features have `spec.md`, `plan.md`, implementation evidence, and `validation.md`;
- code owns explicit internal read contracts and a reusable fake seam;
- unit tests prove M-03-C01 and M-03-C02;
- no Sankhya write path exists;
- no secret leak is observed;
- milestone validation writes `validation-result.md` with verdict grounded in current evidence.
