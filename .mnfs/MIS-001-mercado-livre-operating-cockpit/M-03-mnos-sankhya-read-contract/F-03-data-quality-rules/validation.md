# Feature Validation

```yaml
id: F-03
type: feature-validation
status: passed
owner: Feature Implementer
parent: F-03
created: 2026-07-07
updated: 2026-07-07
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-03-data-quality-rules

## Summary

The Oracle-first rewrite preserves and extends the quality-state guardrails. Missing values stay nil, ambiguity stays explicit, and local proof does not overclaim live Oracle integration.

## Current Validation State

- Result: Passed for local contract/adapter behavior
- Result owner: Feature Implementer
- Decision date: 2026-07-07
- Final feature state for handoff: ready_for_downstream_consumers

## Evidence

- Command: `cd apps/server_core; $env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/modules/internal_read/...`
  - Result: Pass
- Command: `rg -n "INSERT|UPDATE|DELETE|MERGE" apps/server_core/internal/modules/internal_read/adapters/oracle`
  - Result: Pass (no matches)

## Observed

- `missing_product`, `ambiguous_product`, `missing_stock`, `missing_price`, `missing_cost`, `missing_tax`, and `stale_source` remain stable quality flags.
- Fake adapter behavior now mirrors the same nil-preserving contract as the Oracle adapter surface.
- Validation wording now distinguishes local fake/unit proof from blocked live Oracle proof.

## Scope Declaration

- contract_validated: Yes
- integration_validated: No
- blocked_for_real_validation: F-03 itself is locally proven; milestone live Oracle proof still belongs to M-03/F-02

## Handoff

- Current status: `passed`
- Next owner: Feature Implementer
- Next action: keep these states unchanged while exercising real Oracle reads
- Required files/evidence: milestone live-validation evidence
- Blockers or open decisions: none for local contract behavior
