# ADR-004: Integration Catalog Plugin Framework

**Date:** 2026-04-25
**Status:** accepted

## Context
Marketplace definitions already used plugin-style registration, but integrations still required central edits in provider registry and composition root. This made each new provider harder to add and increased wiring churn.

## Decision
Integration providers self-register catalog definitions, auth adapter factories, and optional fee syncers. Composition root consumes registries instead of owning provider-specific construction.

## Rationale
This matches the product model of installable integrations while keeping provider-specific configuration close to provider code. Registration reduces repeated central edits without introducing runtime plugin loading complexity.

## Consequences
Future providers are added through provider-owned packages plus side-effect imports where Go requires registration. Duplicate provider codes are treated as startup configuration errors.

## Alternatives Considered
- Keep root wiring: rejected because it scales poorly with each provider.
- Dynamic runtime plugins: rejected as unnecessary complexity for the current modular monolith.
