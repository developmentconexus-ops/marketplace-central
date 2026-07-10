# F-03-data-quality-rules

```yaml
id: F-03
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-03
created: 2026-07-06
updated: 2026-07-07
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-001 Mercado Livre Operating Cockpit.

## Milestone

M-03 Oracle internal-read contract.

## Brief

Preserve and harden reusable data-quality rules so Oracle-backed internal facts degrade visibly and truthfully across `product_links`, `inventory`, `orders`, and `profitability`.

## Inputs

- Rewritten F-01 contract.
- Rewritten F-02 adapter behavior.
- Planned downstream blocked/incomplete states.

## Expected Output

- Shared value objects/enums for missing, ambiguous, stale, unsupported, and source-unavailable states.
- Tests prove missing values never become business-safe defaults by accident.
- Module docs and validation language reflect real evidence boundaries.

## Constraints

- Do not put UI copy or workflow policy here.
- Do not hide unsupported Oracle query shapes behind fake success values.

## Validation Expectations

- Unit tests cover quality-state stability and nil-preserving semantics.
- Validation artifacts explicitly separate fake-seam proof from live Oracle proof.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: Rewrite `spec.md` with the Oracle-first quality-state model.
- Required files/evidence: F-03/validation.md.
- Blockers or open decisions: exact stale/freshness thresholds may depend on module-specific consumers.
