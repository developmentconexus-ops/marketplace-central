# Feature Spec

```yaml
id: F-02
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: F-02
created: 2026-07-08
updated: 2026-07-08
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-02-link-candidate-engine

## Problem

Imported listing snapshots still are not operationally useful until MPC can turn them into explicit, reviewable product-link candidates with exact-first confidence semantics and unresolved/conflict safety states.

## Requirements

- Read persisted `product_link_listing_snapshots` as the source of listing truth for candidate generation.
- Reuse the MPC-owned `internal_read` boundary for product lookup; do not query Oracle or provider adapters from `product_links` directly.
- Generate and persist candidate states: `manual`, `exact_sku`, `exact_ean`, `title_match`, `unresolved`, `conflict`.
- Use `seller_sku` and `ean` as exact signals before any title fallback.
- If more than one exact mapping is plausible, persist `conflict` instead of guessing.
- If exact `seller_sku` and exact `ean` disagree, persist `conflict`.
- Title fallback may produce reviewable `title_match` candidates, but it must not auto-resolve ambiguity.
- No exact/title hit must persist `unresolved`.
- Expose a manual trigger so the generation can be validated live before any background scheduler exists.

## Non-Goals

- Final operator approval/rejection workflow.
- Automatic resolution from title heuristics.
- Any Mercado Livre write action.

## Acceptance Criteria

- Focused tests prove exact SKU, exact EAN, exact conflict, title fallback, and unresolved behavior.
- Generated candidates are persisted and replace stale candidates for the same listing identity instead of duplicating indefinitely.
- OpenAPI and `sdk-runtime` expose the generation surface.
- Live validation proves candidate generation against real Mercado Livre snapshots and real Oracle reads.

## Handoff

- Current status: `spec_ready`
- Next owner: Feature Implementer
- Next action: implement generation, persistence, and live validation
- Required files/evidence: plan, validation, API contract, persisted candidate evidence
- Blockers or open decisions: none
