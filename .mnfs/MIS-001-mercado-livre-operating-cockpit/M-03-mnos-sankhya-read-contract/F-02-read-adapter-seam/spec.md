# Feature Spec

```yaml
id: F-02
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: F-02
created: 2026-07-07
updated: 2026-07-07
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-02-oracle-adapter-implementation

## Problem

The project needs a real Oracle-backed implementation of `ports.Reader`. A guarded seam that always returns unavailable is no longer enough; M-03 now requires actual adapter ownership, query structure, mapper discipline, and live validation capability.

## Requirements

- Implement Oracle adapter configuration and connection bootstrap inside `internal_read/adapters/oracle`.
- Isolate query text/builders, row mapping, and adapter error translation inside the Oracle boundary.
- Keep an application-facing service over `ports.Reader` so downstream modules stay adapter-agnostic.
- Keep a deterministic fake adapter for unit and downstream tests.
- Keep secrets out of config/bootstrap/runtime errors and artifacts.

## Non-Goals

- Shipping stock, product-link, profitability, or order workflows in this feature.
- Adding HTTP routes.
- Solving every future performance optimization before proving correctness.

## Design

This feature turns `internal_read/adapters/oracle` into the real implementation boundary. The adapter owns:

- config loading and validation;
- driver/connection lifecycle;
- query organization by read surface;
- row-to-domain mapping;
- MPC-owned error translation;
- source timestamp capture.

The application service continues to expose the same `ports.Reader` surface. Downstream business modules remain unable to import Oracle driver packages, SQL helpers, or config structs directly.

## Edge Cases

- Oracle connectivity may be valid but a specific query shape may still be unsupported.
- Price/cost/tax rows may exist partially or with date-window ambiguity.
- Adapter mapping must preserve `nil`/quality-state semantics even when Oracle rows are sparse.
- Runtime evidence must avoid using synthetic success data when the milestone claims real Oracle behavior.

## Acceptance Criteria

- Criterion: Oracle adapter exists as the real `ports.Reader` implementation with isolated query/mapping ownership.
  - Traces to milestone criterion ID: M-03-C03
  - Proven by: focused package tests plus adapter-boundary inspection.
- Criterion: config/bootstrap/runtime errors are secret-safe and read-only discipline is preserved.
  - Traces to milestone criterion ID: M-03-C04
  - Proven by: config tests, grep checks, and validation artifacts.

## Handoff

- Current status: `spec_ready`
- Next owner: Feature Implementer
- Next action: write the execution plan and implement the Oracle adapter boundary
- Required files/evidence: feature brief, spec, milestone contract, validation expectations
- Blockers or open decisions: exact Oracle runtime stack and credential source integration must be finalized during implementation
