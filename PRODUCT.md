# Marketplace Central — Product Design Context

> **Role:** derived design-context aid for frontend tooling. This file is not Product, architecture, program-status or implementation authority. Current authority remains in [`docs/roadmap.md`](docs/roadmap.md) and the owners routed by [`docs/index.md`](docs/index.md).

## Register

product

## Users

- Marketplace Operations Operators preparing products, controlling listings and resolving operational exceptions inside accepted policy.
- Commercial / Marketplace Managers making explainable commercial decisions and governing legitimate MPC-owned policy.
- Fulfillment / Dispatch Operators progressing eligible physical work without losing provider or business-system truth.
- Owners / Administrators / Policy Approvers governing exceptional authority, access and integration boundaries.

Their work is operational, information-dense and consequential. They need exact Organization, marketplace, source and subject context; honest knowledge states; clear next actions; and recovery after uncertain or conflicting outcomes.

## Product Purpose

Marketplace Central is the internal Marketplace Operations Control Plane + Commercial Intelligence product. It combines authoritative internal/business facts with marketplace observations so operators can understand reality, detect divergence, decide within policy, execute controlled actions and verify or reconcile outcomes.

Success means the operator can move through `observe → understand → reconcile → decide/policy → execute → verify → audit/reconcile` without the frontend inventing business truth or hiding uncertainty.

## Brand Personality

Grounded, controlled, explainable.

The interface should communicate operational confidence through explicit scope, provenance, knowledge state and consequences—not through decorative certainty or oversimplification.

## Anti-references

- generic marketplace dashboards that optimize for attractive metrics instead of operational decisions;
- provider-admin clones that expose transport/provider vocabulary as the Product mental model;
- backend-shaped CRUD/navigation or endpoint inventories presented as UX;
- opaque automation, AI recommendations or readiness scores that conceal evidence and authority;
- decorative SaaS surfaces that reduce information density, hide decision-critical facts or use color as the only meaning carrier.

## Design Principles

1. Show exact operating context before action: Organization is global; marketplace account and source-qualified subject are explicit where material.
2. Preserve truth distinctions: known-empty, unknown, unavailable, partial, unsupported, conflicting, pending, ambiguous and converged must not collapse.
3. Put the human job before backend topology while keeping every material read, write and identity traceable to accepted authority.
4. Make consequential actions explicit, reversible or reconcilable as authority permits; never imply success before authoritative reread.
5. Use progressive disclosure for technical evidence without hiding facts required for the current decision.

## Accessibility & Inclusion

No formal WCAG conformance level is currently accepted as Product authority. Current structural requirements are keyboard-operable material controls, logical focus order, semantic controls and labels, non-color-only meaning, screen-reader-plausible structure, and responsive transformations that preserve scope and authority meaning. Motion must remain non-essential and honor reduced-motion preference when introduced.

## Canonical Sources

- [`docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md`](docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md)
- [`docs/engineering/rebaseline/D6-FRONTEND.md`](docs/engineering/rebaseline/D6-FRONTEND.md)
- [`docs/development/frontend-product-experience-planning-method.md`](docs/development/frontend-product-experience-planning-method.md)
- [`docs/roadmap.md`](docs/roadmap.md)
