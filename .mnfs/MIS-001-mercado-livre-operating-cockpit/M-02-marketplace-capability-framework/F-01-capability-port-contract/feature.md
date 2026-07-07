# F-01-capability-port-contract

```yaml
id: F-01
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

Define business-facing capability ports matching IC-001 for listing read, stock read/write, and order read. Keep provider ids as strings and provider payloads outside business modules.

## Inputs

- IC-001.
- Existing `integrations` and `connectors` packages.

## Expected Output

- Go interfaces and domain DTOs for capability use.
- Unit tests showing a business service can use a fake capability implementation.

## Constraints

- Do not implement Stock Seguro policy here.
- Do not add provider-specific endpoint structs to domain/application.

## Inputs/Outputs

- Input: provider code, installation/account refs, listing/order ids.
- Output: normalized snapshots and structured errors from IC-001.

## Negative Scenarios

- Unsupported capability returns one defined structured error.
- Missing credential/install returns one defined structured error.

## Validation Expectations

- Go tests compile with fake capability implementations.
- `rg` confirms no Mercado Livre HTTP adapter import in business application/domain layers.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: accepted.
- Next owner: Milestone Orchestrator.
- Next action: Dispatch F-02 using the accepted capability contracts from F-01.
- Required files/evidence: `spec.md`, `plan.md`, `validation.md`, and connector capability changed paths.
- Blockers or open decisions: None.
