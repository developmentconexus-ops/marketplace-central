# F-07 order-specific tax provenance specification

```yaml
id: F-07
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: F-07
created: 2026-07-12
updated: 2026-07-12
validation_level: QA-1
lifecycle_scope: feature
```

## Feature ID

F-07-order-specific-tax-provenance

## Problem

The Oracle tax read currently selects by product, incidence, and sale date and
sums every matching row. That result can combine unrelated Oracle sale lines
and is therefore not truthful provenance for one Mercado Livre order item.

## Requirements

- The MPC-owned read contract accepts an exact Oracle sale-line identity
  (`NUNOTA`, `SEQUENCIA`) and returns that identity with the tax result.
- The Oracle adapter selects with exact `NUNOTA` and `SEQUENCIA` predicates and
  retains the product predicate as a consistency check; it performs no
  product/date aggregation.
- An absent or invalid sale-line identity does not invoke Oracle and returns
  nil ICMS/IPI/PIS/COFINS plus `missing_tax`.
- Profitability may request tax only when its resolved order-item fact carries
  a verified Oracle sale-line identity. Existing resolved product linkage alone
  is insufficient and therefore remains incomplete for tax.
- Known partial components remain known, while missing components keep the
  profitability item and order incomplete.

## Acceptance Evidence

- Focused internal-read adapter tests inspect exact selection and missing
  provenance behavior.
- Focused profitability tests cover resolved order items with and without an
  exact Oracle sale-line identity and partial tax composition.

## Non-Goals

- Discovering or approving a Mercado Livre-to-Oracle business mapping.
- Adding Oracle/provider writes, heuristics, date matching, or zero defaults.
- Changing public APIs, SDKs, protected research, Candidate A evidence, or
  manual adjustments.

## Design

Introduce a small `TaxSourceIdentity` value containing Oracle document and line
numbers. Carry it from `OrderItemFact` through profitability's internal-read
port into `ports.TaxInput`. The Oracle adapter queries `VW_IMPOSTO_ITEM` for
exact document, sequence, product, and incidence identity. `TaxInputs` exposes
the same identity as provenance. The fake and profitability boundaries return
explicit missing tax without calling Oracle when the identity is unavailable.

This does not assert that MPC currently knows the real Oracle identity. It
creates the truthful contract and preserves the operational gap until an
owner-approved upstream linkage supplies it.

## Edge Cases

- Zero or negative document/line numbers are unverified and yield missing tax.
- Exact identity with no rows yields all nil values and `missing_tax`.
- A subset of tax components may be present; absent components remain nil.
- A resolved internal product without Oracle sale-line identity still yields
  missing tax and no Oracle tax query.

## Acceptance Criteria

### F07-AC01 — Exact source identity and selection

- Traces to milestone criterion ID: `M-06-C02`.
- Proven by: focused internal-read domain, fake, and Oracle adapter tests.

### F07-AC02 — Unknown provenance remains incomplete

- Traces to milestone criterion ID: `M-06-C02`.
- Proven by: focused profitability service tests.

### F07-AC03 — Partial tax preserves known values without completed margin

- Traces to milestone criterion ID: `M-06-C02`.
- Proven by: focused profitability snapshot tests.

## Handoff

- Current status: `spec_ready`
- Next owner: Feature Implementer
- Next action: Follow `plan.md`, compile/validate scoped context, and implement.
- Required files/evidence: feature brief, plan, focused Go test evidence.
- Blockers or open decisions: Real Oracle sale-line linkage remains an upstream
  operational input; the scoped implementation must represent its absence.
