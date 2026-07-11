# ADR-006: MPC owns Oracle internal-read adapters inside `apps/server_core`

**Date:** 2026-07-07
**Status:** accepted

## Context

Marketplace Central needs reliable internal product, stock, price, cost, tax, and sales facts for Mercado Livre operations.
The older `MS_DATABASE_URL` / MetalShopping Postgres read path is no longer considered a correct source model for the target cockpit.
The project already uses a modular-monolith Go core with explicit ports and adapters, and the internal-read boundary must now match that architecture instead of preserving a legacy shortcut.

## Decision

Marketplace Central will read internal ERP facts directly from Oracle through MPC-owned ports and adapters inside `apps/server_core`.

The target structure is:

- module-owned read contracts in MPC;
- Oracle adapter/query/mapping ownership inside MPC;
- downstream modules consume typed domain models only;
- PostgreSQL remains MPC-owned persistence only.

No `MS_DATABASE_URL` direct-read path remains in the target architecture.

## Rationale

- Keeps the business boundary inside the canonical backend where the domain already lives.
- Avoids creating a networked internal service before there is multi-consumer or independent-operations pressure.
- Preserves a clean extraction path later because the core depends on ports, not Oracle details.
- Eliminates the false abstraction of reading from the wrong legacy store just because it is easier to wire.
- Supports global-maximum structure: no ad hoc SQL in business modules, no hardcoded source semantics spread through the app.

## Consequences

- MPC must own Oracle query, mapping, secret-handling, and unavailability semantics explicitly.
- `internal_read` becomes a first-class module/boundary for `product_links`, `inventory`, `orders`, and `profitability`.
- Real validation for internal-read milestones must use Oracle-backed evidence, not only fake seams.
- Historical MNOS mapping knowledge can be used as reference evidence, but not as outsourced execution truth.

## Alternatives Considered

- Keep direct MetalShopping/Postgres reads: rejected because the source model is considered wrong for the target system.
- Introduce a separate internal service now: rejected because it adds operational and network complexity without a proven shared-platform need yet.
