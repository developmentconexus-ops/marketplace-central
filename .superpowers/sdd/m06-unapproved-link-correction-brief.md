# M-06 cold-gate correction: unapproved candidate leakage

## Context

The full live gate found a persisted order item for Mercado Livre listing
`MLB4834373620` with `link_quality=unresolved` while
`internal_product_id=20303`. No product link has been approved. The current
`orders/adapters/productlinks.LinkReader` copies a candidate's internal product
ID even though candidate states map to unresolved/conflict. Profitability then
checks only whether the ID is nil before querying Oracle. This allowed an
unapproved candidate to produce known Oracle cost/tax inputs.

This contradicts F-02 spec: if an order item does not have a resolved product
link, cost/tax remain missing and quality reflects the link truth. It also
invalidates prior live cost/tax evidence as proof of a resolved link.

## Required correction

Use strict TDD. Preserve the heavily dirty worktree and do not edit unrelated
files, reset, revert, stage, or commit.

1. Add a focused adapter test proving link candidates in exact-EAN, conflict,
   and unresolved states may preserve their quality but never expose an
   `InternalProductID` to orders. A persisted workflow link may expose an ID
   only when its state is `resolved`; rejected/conflict/unresolved links must
   not expose an ID either.
2. Run the focused test and record genuine RED caused by the current leakage.
3. Correct the source mapping in `orders/adapters/productlinks` so only an
   explicitly resolved persisted product link can provide the internal product
   ID.
4. Add a profitability application regression proving a non-resolved
   `OrderItemFact` with a nonnil internal ID still does not call internal_read;
   cost and all tax inputs remain missing with the exact link-quality truth.
5. Run that focused test and preserve genuine RED before production changes to
   profitability.
6. Add defense-in-depth in profitability: Oracle cost/tax reads require both
   `LinkQuality == resolved` and a nonnil product ID. All other link qualities
   follow the existing link-blocked path. Do not convert unknown values to
   zero or realized values.
7. Run focused GREEN and impacted suites, `gofmt`, and boundary `rg` proving
   orders imports remain only under `profitability/adapters/orders`.
8. Do not approve or mutate any real product link. Do not claim the prior live
   Oracle values prove a resolved-link flow.

## Expected paths

- `apps/server_core/internal/modules/orders/adapters/productlinks/link_reader.go`
- `apps/server_core/internal/modules/orders/adapters/productlinks/link_reader_test.go`
- `apps/server_core/internal/modules/profitability/application/service.go`
- `apps/server_core/internal/modules/profitability/application/service_test.go`

## Report

Write `.superpowers/sdd/m06-unapproved-link-correction-report.md` containing:
root-cause trace, RED outputs, changed paths, GREEN outputs, gofmt/boundary
evidence, and any concerns. Return only status plus a one-line summary.
