# F-02-catalog-product-360

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-13
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-2
lifecycle_scope: feature
```

## Mission

MIS-001 MVP operator journey.

## Milestone

M-13 Integrated Operator Workspaces.

## Brief

Create `/products` and `/products/:productId` as the canonical catalog and Product
360, joining source facts and links without moving domain calculations into React.

## Inputs

- Passed M-09 product contract, existing catalog/classification/pricing/listing/stock/
  sales/profitability SDK reads, IC-003 ProductRef and quality states.

## Inputs/Outputs

- Input: positive productId and optional installation/filter context.
- Output sections: Overview, Listings, Stock, Price and Margin, Sales, Data and Evidence.
- Lists use IC-003 ordering and deep-link targets.

## Interaction Model

- Product list search/filter opens a stable detail route.
- Product detail sections preserve product header and installation context.
- Listing and sale rows navigate to their canonical detail routes; back restores product.

## State Model

Each section independently renders loading, error/retry, empty, and the server-owned
IC-003 quality. Stale retains its last value/time; unknown renders null + reason;
conflict keeps evidence and blocks dependent actions. React never infers these states.

## Negative Scenarios

- Nonpositive/noninteger ID: `invalid_identity`.
- Product absent: `not_found` with return to Product list.
- One product has conflicting listing links: render conflict and block simulation link.

## Expected Output

One Product 360 makes Sankhya identity, listings, safe stock, sales, margin, and data
quality understandable without duplicating source truth.

## Constraints

- No client-side stock/margin/tax/fee calculations.
- Classifications become product filters/sections; no destructive removal of admin data.
- Owned paths: feature-products, feature-local detail/table components, SDK
  consumption, and this root. `packages/ui` is read-only for this feature.
- Forbidden paths: provider adapters, server business rules, Operations and Sales ownership.

## Criteria IDs

- M-13-C02 Product/listing/sale cross-links.
- M-13-C05 State and responsive consistency.
- M-13-C06 SDK-only thin client.

## Validation Expectations

- Product 1001 route survives reload and displays separate CODPROD/EAN/reference.
- Known zero and unknown null render differently with source and observed time.
- Listing and sale links preserve installation/product context.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: Briefed.
- Next owner: Feature Implementer after F-01 shell acceptance.
- Next action: Create spec/plan and implement Product 360 using accepted SDK reads.
- Required files/evidence: M-09 `validation-result.md`, IC-003, and this feature's `validation.md`.
- Blockers or open decisions: Missing server-owned fact requires IC/OpenAPI/SDK amendment.
