# F-02 Specification - Margin Input Model

```yaml
id: M-06-F-02
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: M-06
created: 2026-07-09
updated: 2026-07-09
```

## Goal

Create the first `profitability` slice that models margin inputs honestly, with explicit source, timestamp, currency, amount, and quality state for each input component before any final contribution calculation is attempted.

## Architectural Position

- `orders` owns persisted marketplace order snapshots.
- `internal_read` owns Oracle-backed cost, tax, and internal fact contracts.
- `profitability` owns margin input composition, quality truth, and manual adjustment audit.
- `F-02` must not calculate final margin or contribution totals; that belongs to `F-03`.

## Scope

This feature includes:

- a new `profitability` module
- a normalized margin-input model at order and order-item scope
- component types for revenue, sale fee, cost, tax, freight, commission, and manual adjustments
- explicit quality states when inputs are missing, unresolved, stale, manual, or complete
- manual adjustment persistence with audit fields
- a backend/API/SDK surface to inspect margin inputs and add manual adjustments

This feature does not include:

- final order profitability calculation
- profitability UI
- automatic freight inference from unsupported provider fields

## Input Truth

### Revenue And Sale Fee

Read from persisted `orders` snapshots:

- item revenue from `quantity * unit_price`
- sale fee from item `sale_fee_amount` when available
- if item sale fee is absent, the model may fallback to order-level fee only as an explicit degraded quality state, never as silent certainty

### Cost Basis

Read from `internal_read` cost contract:

- `CostBasisCUSSEMICM` is the only allowed initial basis
- `CUSVARIAVEL` is forbidden in this slice
- missing cost stays `nil` with explicit quality state

### Tax Inputs

Read from `internal_read` tax contract:

- ICMS, IPI, PIS, COFINS fields stay explicit
- missing tax stays `nil` with explicit quality state

### Freight And Commission

For `F-02`, freight and extra commission inputs are modeled as:

- absent/unknown by default
- manually adjustable at order or item scope
- never defaulted to zero

## Domain Shape

### Component Kinds

Required component kinds:

- `revenue`
- `sale_fee`
- `cost`
- `tax_icms`
- `tax_ipi`
- `tax_pis`
- `tax_cofins`
- `freight`
- `commission`
- `manual_adjustment`

### Quality States

Required quality states:

- `complete`
- `missing`
- `unresolved_link`
- `rejected_link`
- `conflict_link`
- `stale`
- `manual`
- `partial`

These states are profitability-owned and should summarize operator-facing input confidence, even when sourced from `internal_read` quality flags.

### Source Metadata

Every component must carry:

- `kind`
- `scope` (`order` or `item`)
- `amount`
- `currency`
- `source_system`
- `source_reference`
- `observed_at`
- `quality_state`
- `quality_reason`

### Manual Adjustment Model

Manual adjustments must support:

- scope: order-level or item-level
- category: freight, commission, cost, generic adjustment
- signed amount
- currency
- reason/note
- operator metadata
- created timestamp

Adjustments are append-only audit events for this slice. Updates should happen by inserting a new adjustment row rather than mutating historical entries.

## Resolution Rules

### Link Dependency

If an order item does not have a resolved product link:

- profitability input assembly still succeeds
- cost/tax components remain missing
- quality state reflects missing/unresolved/rejected/conflict link truth

### Unknown Value Policy

If fee, freight, cost, or tax is unknown:

- keep the component amount `nil`
- emit explicit non-complete quality
- never normalize to `0`

## Persistence Design

Add first-class profitability tables:

- `profitability_margin_inputs`
- `profitability_manual_adjustments`

The first table persists the latest assembled input snapshot per order component. The second persists append-only manual adjustments with actor/note audit.

## API Surface

Recommended initial surface:

- `POST /profitability/margin-inputs/import`
- `GET /profitability/margin-inputs?installation_id=...`
- `POST /profitability/manual-adjustments`
- `GET /profitability/manual-adjustments?installation_id=...`

## Validation Requirements

Tests must cover:

- complete input shape
- missing cost
- missing freight
- manual freight
- manual commission
- adjustment audit persistence

Runtime evidence should prove:

- real order snapshot can be transformed into profitability inputs
- missing link remains explicit quality, not blocker
- manual adjustment persists with audit fields

## Non-Negotiables

- `CUSSEMICM` only for initial cost basis
- no zero-default fallback for unknown values
- no direct SQL into `orders` or `product_links` internals from other modules; use adapters/ports
- no UI-first shortcut that skips backend truth
