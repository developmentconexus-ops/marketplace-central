# F-01-canonical-product-identity

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-09
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-1
lifecycle_scope: feature
```

## Mission

MIS-001 MVP operator journey.

## Milestone

M-09 Canonical Product Foundation.

## Brief

Define one product identity contract in catalog/domain/OpenAPI/SDK where positive
`internal_product_id` equals Sankhya CODPROD and all other identifiers remain separate.

## Inputs

- `internal_read/domain.ProductCandidate` and source metadata.
- Current catalog Product, product-link DTOs, OpenAPI, and SDK types.
- IC-003 ProductRef and SourceFact.

## Inputs/Outputs

- Input: positive CODPROD and nullable EAN/reference/brand/group/source facts.
- Output: product DTO with `internal_product_id`, separate identifier fields, nullable
  numeric source facts, `source`, `quality`, `observed_at`, and nullable
  `quality_reason` required whenever quality is `unknown|stale`.
- Existing public contract changes update OpenAPI and SDK in one commit.

## Expected Output

All product consumers share one typed positive integer identity and can distinguish
unknown from known zero without provider/manufacturer identifiers leaking into it.

## Negative Scenarios

- CODPROD <= 0: reject as `invalid_identity`; emit no product.
- EAN/reference/seller SKU only: do not synthesize internal_product_id.
- Missing numeric fact: output null plus `unknown` and nonblank `quality_reason`.
- Conflicting identities: return `identity_conflict`; do not choose one.

## Constraints

- Preserve tenant scope and domain/provider boundaries.
- No database cutover, UI redesign, provider write, or authentication work.
- Owned paths: catalog/internal_read domain/ports, OpenAPI, SDK, and this feature root.
- Forbidden paths: web workspaces, provider adapters outside read mapping, M-06 evidence.

## Criteria IDs

- M-09-C01 Canonical identity.
- M-09-C02 Honest nullable facts.

## Validation Expectations

- Contract tests serialize CODPROD 1001 separately from EAN/reference/seller SKU.
- Tests distinguish missing value from known zero.
- OpenAPI and SDK field/type parity is exact.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: Briefed.
- Next owner: Feature Implementer.
- Next action: Create `spec.md`, then `plan.md`, compile context, and implement one commit.
- Required files/evidence: this brief, IC-003, M-09 contract, and this feature's `validation.md`.
- Blockers or open decisions: Stop on any unproved mapping from legacy string IDs to CODPROD.
