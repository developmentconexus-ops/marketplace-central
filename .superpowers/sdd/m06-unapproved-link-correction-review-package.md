# Review package: M-06 unapproved candidate leakage correction

## Scope

Review only the correction described by:

- Brief: `.superpowers/sdd/m06-unapproved-link-correction-brief.md`
- Report: `.superpowers/sdd/m06-unapproved-link-correction-report.md`

## Changed paths

- `apps/server_core/internal/modules/orders/adapters/productlinks/link_reader.go`
- `apps/server_core/internal/modules/orders/adapters/productlinks/link_reader_test.go`
- `apps/server_core/internal/modules/profitability/application/service.go`
- `apps/server_core/internal/modules/profitability/application/service_test.go`

These module trees pre-existed as untracked content in a heavily dirty shared
worktree, so Git cannot produce a trustworthy before/after diff without
capturing unrelated user work. Read the four scoped current files and compare
their behavior to the brief and the RED/GREEN transcript in the report.

## Binding constraints

- An unapproved candidate is not a resolved product link.
- Only an explicitly resolved persisted product link may expose an internal
  product ID to orders.
- Profitability must require both resolved quality and nonnil product ID before
  reading Oracle/internal facts.
- Rejected, conflict, unresolved, and missing links keep cost/tax nil and retain
  their exact quality truth.
- No unknown-to-zero or unknown-to-realized conversion.
- Profitability may import orders internals only in
  `profitability/adapters/orders`.
- No real product link approval belongs to this correction.

## Reviewer output contract

Return two explicit verdicts:

1. `SPEC: APPROVED` or `SPEC: REJECTED` with exact findings.
2. `QUALITY: APPROVED` or `QUALITY: REJECTED` with actionable findings and
   severity.

Do not edit files.
