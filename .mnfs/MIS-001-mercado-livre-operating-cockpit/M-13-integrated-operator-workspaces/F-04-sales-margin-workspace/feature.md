# F-04-sales-margin-workspace

```yaml
id: F-04
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-13
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-3
lifecycle_scope: feature
```

## Mission

MIS-001 MVP operator journey.

## Milestone

M-13 Integrated Operator Workspaces.

## Brief

Evolve Orders into `/sales` and routable sale detail, connecting order/item identity,
listing, canonical product, assisted Sankhya linkage, margin inputs, and profit quality.

## Inputs

- Existing Orders/Profitability UI and SDK reads/import/calculate operations.
- F-11 assisted Sankhya-linkage SDK operations.
- IC-003 SaleRef, SourceFact, deep links, and actor label.

## Inputs/Outputs

- Input: installation, quality/filter/page, provider order ID.
- Output: sale/order facts, item identities, Product/Listing links, linkage evidence,
  inputs, adjustments, snapshots, null/quality/source/time, and contribution/margin.

## Interaction Model

- Sale list opens stable detail; detail links to Product and Listing.
- Import/recalculate actions refetch and retain selected sale.
- Assisted linkage confirmation is local evidence; actor copy is unverified.

## State Model

Profit quality retains existing domain values; source quality uses IC-003.
Not calculated, incomplete, negative, not realized, conflict, and complete remain
visually distinct. No state is inferred from amount alone.

## Negative Scenarios

- Order absent for installation: `not_found`.
- Linkage ambiguous: `identity_conflict`, retain candidate evidence.
- Missing tax/cost/fee/freight: null plus existing missing flag; no realized margin.
- Manual adjustment lacks required local label/reason: `invalid_request`; no record.

## Expected Output

The operator traces a sale from Mercado Livre through listing/product/Sankhya facts to
an explainable margin without authentication being a prerequisite for read analysis.

## Constraints

- Preserve current profitability formulas and M-06 evidence; this feature reorganizes UX.
- Manual records show `operator_supplied_unverified`; no production-auth claim.
- Owned paths: feature-orders, route wrappers, relevant SDK consumption, this feature root.
- `packages/ui` is read-only; sales-specific components stay in `feature-orders`.
- Forbidden paths: profitability formula changes, provider writes, auth/RBAC.

## Criteria IDs

- M-13-C02 Product/listing/sale cross-links.
- M-13-C05 State and responsive consistency.
- M-13-C06 SDK-only thin client.

## Validation Expectations

- Deep link reloads ORDER-MVP-1 and preserves installation.
- Unknown tax is null and prevents complete margin.
- Product and Listing links use canonical identities, and conflict stays explicit.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: Briefed.
- Next owner: Feature Implementer after Product/Listing routes stabilize.
- Next action: Create spec/plan and adapt the existing Orders workspace.
- Required files/evidence: M-06 feature `validation.md` files, IC-003, and this feature's `validation.md`.
- Blockers or open decisions: A formula defect is separate correction authority, not UX scope.
