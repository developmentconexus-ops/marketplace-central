# F-02-link-candidate-engine

```yaml
id: F-02
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

Generate product link candidates from listing snapshots using exact EAN and `seller_sku` first, then title heuristic candidates that require approval.

## Inputs

- Listing snapshots.
- IC-002 FindProductsForLinking.

## Expected Output

- Candidate states: `manual`, `exact_sku`, `exact_ean`, `title_match`, `unresolved`, `conflict`.
- Conflict when more than one exact product/listing mapping is plausible.

## Constraints

- Do not auto-resolve ambiguous title matches.
- Do not treat missing EAN/SKU as failure; produce unresolved or title candidates.

## Negative Scenarios

- Multiple exact EAN results create `conflict`.
- No candidates creates `unresolved`.

## Validation Expectations

- Unit tests cover exact EAN, exact SKU, title fallback, no match, and conflict.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: Create `spec.md`.
- Required files/evidence: F-02/validation.md.
- Blockers or open decisions: None.
