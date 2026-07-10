# Feature Spec

```yaml
id: F-03
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: F-03
created: 2026-07-07
updated: 2026-07-07
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-03-data-quality-rules

## Problem

The Oracle-first read boundary must make data quality explicit enough that downstream modules never confuse missing, ambiguous, stale, unsupported, or unavailable facts with valid operational data.

## Requirements

- Keep shared quality flags stable and intentionally named.
- Prove missing product, stock, cost, tax, and source freshness remain explicit quality states.
- Prove missing numerics remain `nil`, not `0`.
- Support unsupported-query and source-unavailable states explicitly where the Oracle boundary cannot answer safely.
- Keep validation wording honest about which states are proven by fake/unit tests versus real Oracle evidence.

## Non-Goals

- Writing UI text.
- Deciding business approval policy for each downstream module.
- Replacing live Oracle validation with fake-only evidence.

## Acceptance Criteria

- Quality-state tests prove stability and nil-preserving behavior for required internal-read outputs.
- Fake seam and Oracle adapter tests both honor the same quality-state contract.
- Validation artifacts explicitly avoid claiming real integration proof from mock/fake-only evidence.

## Handoff

- Current status: `spec_ready`
- Next owner: Feature Implementer
- Next action: write the execution plan and harden the quality-state implementation/tests
- Required files/evidence: feature brief, spec, milestone contract, validation expectations
- Blockers or open decisions: freshness thresholds may stay configurable if downstream modules need different tolerances
