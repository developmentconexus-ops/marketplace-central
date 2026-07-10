# Feature Spec

```yaml
id: F-01
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: F-01
created: 2026-07-07
updated: 2026-07-07
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-01-oracle-read-contract-redesign

## Reference Evidence

- ADR-006 in `.brain/decisions/006-oracle-internal-read-owned-by-mpc.md`
- Legacy Oracle mapping material from `C:\Users\leandro.theodoro\Documents\MNOS\...` as reference-only evidence
- Current `internal_read` module shapes in `apps/server_core/internal/modules/internal_read`

## Problem

Marketplace Central already has an `internal_read` seam, but it was briefed against a superseded architecture. The project now needs an MPC-owned Oracle read contract that expresses internal facts in business-safe types and policies without relying on `MS_DATABASE_URL`, imported legacy assumptions, or scattered hardcoded semantics.

## Requirements

- Requirement: define MPC-owned domain and port contracts for product lookup, sellable stock, current price, cost-as-of, tax inputs, and sales history.
  - Acceptance evidence: `internal_read/domain` and `internal_read/ports` compile and are consumed without Oracle-specific imports outside the adapter boundary.
- Requirement: separate source semantics from operator policy.
  - Acceptance evidence: stock scope, company/location filters, cost basis selection, and freshness expectations are modeled explicitly as policy/value objects or typed config inputs, not hidden constants inside downstream modules.
- Requirement: preserve quality-state semantics for missing, ambiguous, stale, or unsupported facts.
  - Acceptance evidence: domain tests prove missing facts remain explicit quality states and nullable fields remain `nil`.
- Requirement: identify which Oracle-backed semantics are provisional reference evidence versus which become MPC-owned contract truth.
  - Acceptance evidence: spec and contract docs explicitly label reference evidence and MPC-owned decisions.
- Requirement: keep the contract read-only and MPC-owned.
  - Acceptance evidence: no ERP write path, no mirror tables, and no downstream SQL reach-through are introduced by this feature.

## Non-Goals

- Implement live Oracle queries.
- Finalize every SQL object name or performance optimization.
- Introduce HTTP routes, scheduler jobs, or module business workflows.
- Preserve the old `SANKHYA_*` error taxonomy as-is if a better MPC-owned naming scheme is needed.

## Design

This feature redefines `internal_read` as a first-class MPC boundary. The domain layer owns the canonical read models and quality flags. The ports layer owns the read operations and typed policy inputs. The contract explicitly distinguishes:

- source facts: product identity, stock rows, price rows, cost rows, tax rows, sales rows, source timestamps;
- MPC policy: default sellable-stock scope, supported stock modes, preferred cost basis, freshness thresholds, ambiguity behavior;
- consumer-facing outcomes: resolved values, quality flags, unsupported-query errors, and source metadata.

The contract must support future changes in Oracle object names or query strategy without forcing `inventory`, `product_links`, `orders`, or `profitability` to change their imports or business logic.

## Edge Cases

- A product may exist but not be uniquely linkable by EAN/SKU/title.
- Stock semantics may require multiple companies/locations and explicit exclusions.
- Cost and tax may be partially available for a product/date combination.
- Some inputs may need product-level queries while others need group-level or date-window queries.
- Oracle evidence may reveal unsupported query shapes; these must become explicit contract errors, not silent best-effort behavior.

## Acceptance Criteria

- Criterion: MPC owns the read contract independently of legacy Postgres/MNOS runtime assumptions.
  - Traces to milestone criterion ID: M-03-C01
  - Proven by: focused `go test` on `internal_read/domain` and `internal_read/ports`, plus import-boundary inspection.
- Criterion: contract semantics distinguish explicit policy, source facts, and quality states.
  - Traces to milestone criterion ID: M-03-C02
  - Proven by: contract tests covering nullable fields, quality flags, and typed policy inputs.

## Handoff

- Current status: `spec_ready`
- Next owner: Feature Implementer
- Next action: write the execution plan and update the domain/ports contract accordingly
- Required files/evidence: feature brief, spec, milestone contract, validation expectations
- Blockers or open decisions: exact default stock and cost policy values must be confirmed against real Oracle evidence during execution
