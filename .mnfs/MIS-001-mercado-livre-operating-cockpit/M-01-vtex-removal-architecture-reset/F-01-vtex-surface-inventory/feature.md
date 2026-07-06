# F-01-vtex-surface-inventory

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-01
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-001 Mercado Livre Operating Cockpit.

## Milestone

M-01 VTEX removal and architecture reset.

## Brief

Produce a precise inventory of VTEX references by category before deletion: backend routes/adapters/ports/tests, OpenAPI schemas/routes, SDK types/methods/tests, frontend pages/nav/tests, docs/wiki/env/migrations.

## Inputs

- ADR-005.
- Current repository tree.
- M-01 validation contract.

## Expected Output

- `spec.md` classifies every VTEX reference as remove, legacy-doc-retain, or migration-risk.
- `validation.md` includes exact `rg` commands and categorized counts.

## Constraints

- Do not delete files in this feature.
- Do not classify Mercado Livre or generic marketplace references as VTEX.

## Validation Expectations

- `validation.md` names at least backend, contracts, SDK, frontend, docs, tests, env, and migrations categories.
- Every active VTEX route or SDK method has a removal owner in F-02.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: Create `spec.md`.
- Required files/evidence: F-01/validation.md.
- Blockers or open decisions: None.
