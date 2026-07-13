# F-03-catalog-compatibility-cutover

```yaml
id: F-03
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

Remove remaining active MSDB composition/config and reconcile MPC-owned enrichment,
classification, pricing, and product-link references to canonical CODPROD without guessing.

## Inputs

- Accepted F-01/F-02 contracts and commits.
- Existing enrichment/classification/product-link foreign references.
- Governance inventory of MSDB exceptions.

## Inputs/Outputs

- Input: existing MPC records and deterministic legacy-to-CODPROD evidence.
- Output: preserved records keyed by canonical product identity or explicit unmapped/conflict records.
- Server runtime no longer reads MS_DATABASE_URL/MS_TENANT_ID.

## Expected Output

All MVP product consumers use CODPROD; valid enrichments survive, ambiguous records remain
unlinked, and active runtime/config/governance no longer depends on MSDB.

## Negative Scenarios

- More than one CODPROD candidate: retain unlinked and report `identity_conflict`.
- No deterministic candidate: retain evidence and report `not_found`; do not drop silently.
- Migration or composition requires rewriting historical M-06 evidence: stop.

## Constraints

- Forward-only migration; no reset/revert/stash/clean.
- OpenAPI/SDK atomicity for public shape changes.
- Owned paths: catalog-owned migrations/repos, product consumer adapters, composition,
  platform/msdb removal, governance registry exception, and this feature root.
- Forbidden paths: M-13 UI redesign and provider writes.

## Criteria IDs

- M-09-C03 Legacy runtime removed.
- M-09-C04 Deterministic compatibility.

## Validation Expectations

- Fixtures prove mapped, unmapped, and ambiguous records.
- Exact repository scan finds no active platform/msdb/config dependency.
- Catalog, classification, pricing, product-link, OpenAPI, and SDK suites exit 0.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: Briefed.
- Next owner: Feature Implementer.
- Next action: Create execution artifacts after F-02 acceptance.
- Required files/evidence: accepted F-01/F-02 `validation.md` files and this feature's `validation.md`.
- Blockers or open decisions: Any ambiguous production mapping requires operator evidence.
