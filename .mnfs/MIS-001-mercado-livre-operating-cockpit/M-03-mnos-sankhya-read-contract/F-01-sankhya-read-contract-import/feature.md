# F-01-oracle-read-contract-redesign

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-03
created: 2026-07-06
updated: 2026-07-07
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-001 Mercado Livre Operating Cockpit.

## Milestone

M-03 Oracle internal-read contract.

## Brief

Redesign the internal-read contract so MPC owns the domain types, read policies, and input/output seams for Oracle-backed product, stock, price, cost, tax, and sales facts.

## Inputs

- ADR-006.
- Existing `internal_read` module surface.
- Legacy MNOS Oracle mapping knowledge as reference evidence only.

## Expected Output

- MPC-owned contract doc and typed interfaces define required fields, policy objects, and quality states.
- The contract names what is configurable policy versus what is source semantics.
- No downstream module needs Oracle table, view, driver, or SQL knowledge.

## Constraints

- Do not preserve `MS_DATABASE_URL` assumptions.
- Do not copy legacy SQL or mapping blindly without re-owning the contract.
- Do not hardcode business semantics in downstream modules.

## Validation Expectations

- Contract identifies the load-bearing Oracle-backed facts required by `product_links`, `inventory`, `orders`, and `profitability`.
- Missing/ambiguous/stale source facts are typed quality states.
- The feature stays at contract/policy ownership level; adapter execution belongs to F-02.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: Rewrite `spec.md` for the Oracle-first contract.
- Required files/evidence: F-01/validation.md.
- Blockers or open decisions: exact Oracle object names and query shapes must be confirmed during execution.
