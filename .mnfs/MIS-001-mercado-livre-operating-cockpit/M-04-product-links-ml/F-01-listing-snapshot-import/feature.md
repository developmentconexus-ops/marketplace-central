# F-01-listing-snapshot-import

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-04
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
```

## Brief

Fetch existing Mercado Livre listings through the capability adapter and persist normalized listing snapshots for link candidate generation.

## Inputs

- IC-001 ListListings/ReadListing.
- Connected Mercado Livre installation.

## Expected Output

- Listing snapshots include provider item id, optional variation id, title, status, seller SKU, EAN, available quantity, fetched timestamp.

## Constraints

- Do not publish new listings.
- Do not mutate provider state.

## Validation Expectations

- Tests cover listing with no variation and listing with variation.
- Re-import is idempotent by provider listing/variation id.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: Create `spec.md`.
- Required files/evidence: F-01/validation.md.
- Blockers or open decisions: None.
