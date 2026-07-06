# F-02-provider-capability-registration

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-02
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-001 Mercado Livre Operating Cockpit.

## Milestone

M-02 Marketplace capability framework.

## Brief

Align provider definitions and capability health with business capability names so the UI and services can know whether listing read, stock read/write, and order read are supported for Mercado Livre.

## Inputs

- Existing provider registry.
- IC-001 capability statuses.

## Expected Output

- Mercado Livre provider definition exposes required capability states.
- Future providers can remain unavailable/blocked without implementing operations.

## Constraints

- Do not imply future providers are operational.
- Do not remove integration auth lifecycle.

## Validation Expectations

- Provider registry tests assert Mercado Livre required capability support.
- Provider catalog UI can still represent unavailable/blocked providers.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: Create `spec.md`.
- Required files/evidence: F-02/validation.md.
- Blockers or open decisions: None.
