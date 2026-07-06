# F-03-architecture-truth-alignment

```yaml
id: F-03
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

Align architecture, wiki, roadmap-facing docs, and quality guidance so Mercado Livre-first hub architecture is the active truth.

## Inputs

- F-01 and F-02 evidence.
- `ARCHITECTURE.md`, `wiki/`, `.brain/`, OpenAPI.

## Expected Output

- Docs name Mercado Livre as first operational adapter and preserve future marketplace capability architecture.
- VTEX retained only as historical legacy, if at all.

## Constraints

- Do not rewrite historical evidence files unless they are active planning truth.
- Do not invent implemented modules that remain planned.

## Validation Expectations

- Docs reference capability-based architecture and Stock Seguro priority.
- No active planning doc instructs future VTEX work.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: Create `spec.md`.
- Required files/evidence: F-03/validation.md.
- Blockers or open decisions: None.
