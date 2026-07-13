# F-07 Order-Specific Oracle Tax Provenance

```yaml
id: F-07
type: feature
status: briefed
owner: Feature Implementer
parent: M-06
created: 2026-07-12
updated: 2026-07-12
validation_level: QA-1
lifecycle_scope: feature
```

## Objective

Make Oracle tax inputs truthful at one resolved Mercado Livre order-item scope.
The current product-plus-date query may aggregate unrelated Oracle sale rows;
it must not be presented as tax provenance for a specific marketplace sale.

## Acceptance Criteria

1. The MPC-owned tax read contract names the exact source identity needed to
   select one Oracle sale line and returns that identity in provenance.
2. The Oracle adapter uses exact, read-only source predicates. It never
product/date-aggregates unrelated rows as an order-item fact.

## Brief

Replace product-and-date Oracle tax aggregation with an exact, read-only sale
line contract. When a verified Oracle document and line identity is absent,
keep tax missing and profitability incomplete.

## Expected Output

- Exact Oracle document/line tax selection with returned provenance.
- Explicit missing tax when verified source identity is unavailable.
- Deterministic internal-read and profitability boundary evidence.
3. When no verified Oracle sale-line identity exists, ICMS/IPI/PIS/COFINS stay
   nil/missing and profitability remains incomplete; no matching heuristic,
   estimate, or zero default is introduced.
4. Deterministic tests cover exact selection, ambiguous/missing provenance,
   partial tax, and the resolved-order profitability composition boundary.
5. The Feature records the exact operational gap to real resolved-order margin
   evidence. Live Oracle reads are allowed only through the registered runner;
   no Oracle/provider write and no secret/PII capture is allowed.

## Constraints

- Preserve domain/application/ports/adapters boundaries.
- Do not change authentication or manual-adjustment surfaces.
- Do not reissue Candidate A approval or mutate its link/audit evidence.
- CUSSEMICM is per-unit and is quantity-extended once; Mercado Livre sale fee
  and Oracle tax components are line amounts and are not quantity-multiplied.
- Stop on a contract/ownership conflict or if real source linkage requires an
  owner-approved business mapping not present in repository truth.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: write `spec.md` and `plan.md`, compile/validate scoped context,
  implement the smallest safe contract correction, run focused proof, write
  `validation.md`, and create one intentional commit.
