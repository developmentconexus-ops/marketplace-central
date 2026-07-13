# F-15 Partial Tax Snapshot Quality Correction

```yaml
id: F-15
type: correction-feature
status: briefed
owner: Feature Implementer
parent: M-06
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-1
lifecycle_scope: feature
correction_attempt: 2
max_correction_attempts: 2
```

## Brief

Retain known partial tax amounts while making every non-complete tax input
quality keep item and order snapshots incomplete and contribution/margin nil.

## Objective

Correct the independent-review-proven F-14 aggregation defect: a known tax
amount whose input quality is `partial` or otherwise incomplete must retain the
known amount but keep item and order snapshots incomplete, with no realized
contribution or margin.

## Required behavior

- During snapshot input application, every tax component whose amount is nil
  **or whose quality is not `complete`** must set the missing-tax/incomplete
  accumulator condition. Continue retaining known non-nil tax amounts.
- Item snapshot rollup must remain incomplete when all four tax components are
  known but partial because TOP 306 lineage is partial.
- Order snapshot rollup must inherit the same incomplete/missing-tax condition
  from any item; it must not compute contribution or margin.
- A complete TOP 306 lineage with an incomplete exact-line source quality must
  behave the same way.
- Fully complete exact-line tax inputs must retain current complete-margin
  behavior. Unknown values remain nil and no amount becomes zero by default.
- Add an import-to-snapshot regression that proves the exact counterexample,
  plus complete-lineage/incomplete-source and complete-control cases.
- Correct F-14 evidence overclaim append-only and append attempt 2 to the
  authoritative correction ledger. C03 remains explicitly deferred/failing.

## Scope

Only profitability snapshot aggregation/tests and M-06 correction/evidence
artifacts. Do not change lineage resolution, Oracle reads, APIs/SDK, migrations,
UI, runtime config, manual adjustments, authentication, or authorization.

## Expected Output

- Snapshot aggregation records missing tax for every nil amount or
  non-complete tax quality while retaining all known amounts.
- Import-to-snapshot regressions prove partial-lineage/all-known components,
  complete-lineage/incomplete-source, and fully complete control behavior.
- F-14 evidence is corrected truthfully and correction attempt 2 is appended
  without weakening the explicitly deferred/failing C03 constraint.
- Focused and repository-wide fake Go proof, scoped diff evidence, and one
  intentional Feature commit establish the bounded correction.

## Proof and handoff

Run focused profitability tests, repository-wide Go tests with
`GOCACHE=.gocache`, and `git diff --check`; write `spec.md`, `plan.md`, scoped
`context.json`, and `validation.md`; create exactly one intentional commit and
return it for a replacement fixed-SHA review.
