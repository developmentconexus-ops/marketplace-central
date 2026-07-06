# F-02-vtex-active-surface-removal

```yaml
id: F-02
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

Remove active VTEX runtime and UI surfaces identified by F-01 while preserving generic marketplace/integrations foundations.

## Inputs

- F-01 inventory.
- OpenAPI, SDK runtime, server composition/router, frontend router/nav.

## Expected Output

- No active VTEX routes, SDK methods/types, adapter registration, frontend pages/nav, or tests remain.
- Historical docs only remain if explicitly marked legacy.

## Constraints

- Do not delete Mercado Livre or generic provider catalog behavior.
- Do not touch unrelated provider definitions except where tests need updated fixtures.

## Inputs/Outputs

- Input: F-01 inventory.
- Output: code/doc diffs plus validation commands proving active VTEX removal.

## Negative Scenarios

- If OpenAPI still exposes `/connectors/vtex/*`, validation fails M-01-C02.
- If frontend still links to VTEX publisher, validation fails M-01-C01.

## Validation Expectations

- OpenAPI/SDK/frontend/router tests pass.
- `rg` shows no active VTEX path outside legacy-marked docs.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: Create `spec.md`.
- Required files/evidence: F-02/validation.md.
- Blockers or open decisions: None.
