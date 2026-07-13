# F-02-oracle-catalog-cutover

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-09
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-2
lifecycle_scope: feature
```

## Mission

MIS-001 MVP operator journey.

## Milestone

M-09 Canonical Product Foundation.

## Brief

Make catalog product reads use MPC-owned internal_read/Oracle semantics and remove
the active MetalShopping/MSDB catalog source without copying ERP state wholesale.

## Inputs

- Accepted F-01 identity contract.
- M-03 internal-read ports/adapters and real Oracle read lane.
- Current catalog repository/service/transport behavior.

## Inputs/Outputs

- Input: tenant-scoped list/get/search using CODPROD, name, EAN, or reference.
- Output: canonical product DTOs ordered name asc then CODPROD asc, with source quality.
- Search never treats provider seller SKU as canonical internal SKU.

## Expected Output

Catalog list/get/search operate without MS_DATABASE_URL and return honest Oracle-backed
product facts suitable for Product 360.

## Negative Scenarios

- Oracle unavailable: `source_unavailable`; no fallback to stale MSDB.
- Duplicate/ambiguous search facts: explicit `identity_conflict`.
- Missing cost/stock/price: null/unknown + nonblank reason, not zero.
- Query would require Oracle write: stop; internal_read is read-only.

## Constraints

- No wholesale Oracle mirror or new canonical product table.
- Provider payloads stay outside catalog.
- Owned paths: catalog application/ports/adapters/transport, internal_read read seam,
  composition required for cutover, and this feature root.
- Forbidden paths: web route redesign, production writes/auth, M-06 QA artifacts.

## Criteria IDs

- M-09-C02 Honest nullable facts.
- M-09-C03 Legacy runtime removed.
- M-09-C05 Real Oracle product read.

## Validation Expectations

- Deterministic tests prove list/get/search order, identity, unavailable source, and null facts.
- Real read-only probe returns a positive CODPROD and source metadata.
- Composition starts without MS_DATABASE_URL.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: Briefed.
- Next owner: Feature Implementer.
- Next action: Create execution artifacts and implement after accepted F-01 commit.
- Required files/evidence: F-01 `validation.md`, M-03 contract, and this feature's `validation.md`.
- Blockers or open decisions: Stop if current Oracle contract lacks a required product fact.
