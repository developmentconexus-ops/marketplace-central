# F-02-provider-capability-registration

```yaml
id: F-02
type: feature-brief
status: accepted
owner: Milestone Orchestrator
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

- Current status: accepted.
- Next owner: Milestone Orchestrator.
- Next action: Dispatch F-03 using the accepted business capability vocabulary from F-01 and F-02.
- Required files/evidence: `spec.md`, `plan.md`, `validation.md`, and the updated contract/provider registration paths.
- Blockers or open decisions: None.
