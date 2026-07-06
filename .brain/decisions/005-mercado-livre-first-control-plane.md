# ADR-005: Mercado Livre First Control Plane

**Date:** 2026-07-06
**Status:** accepted

## Context
Marketplace Central was originally planned around VTEX as the operational commerce hub.
The business no longer uses VTEX as the ecommerce platform and the urgent operating need is Mercado Livre stock, price, order, and margin control using Sankhya/MetalShopping data.

## Decision
Marketplace Central will remove VTEX from the target architecture and treat Mercado Livre as the first operational marketplace control plane.
Sankhya/MetalShopping remains the source of truth for internal product, stock, price, cost, tax, and sales data.

## Rationale
Continuing VTEX work would optimize a dead path and preserve a local maximum.
Mercado Livre solves the current revenue and operational risk directly: overselling, manual stock updates, unclear fees, and unknown per-order margin.
Other marketplaces stay visible only as future provider catalog entries until Mercado Livre operations are reliable.

## Consequences
VTEX docs, roadmap items, connector work, and environment keys are legacy and must not receive new feature work.
New work should start from Mercado Livre product links, stock reconciliation, order ingestion, and profitability.
Any remaining VTEX code must be inventoried before deletion to avoid accidental contract or migration breakage.

## Alternatives Considered
- Keep VTEX as primary: rejected because the business no longer uses it.
- Stay multi-marketplace-first: rejected because it delays the immediate Mercado Livre operational fix.
- Patch existing VTEX connector into a generic connector: rejected because it would preserve wrong abstractions.
