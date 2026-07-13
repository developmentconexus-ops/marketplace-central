# M-09 Inventory Clock QA-Unblock Validation

- work item: `M-09-QA-CORR-01`
- dispatch base: `ee71d7ab2a66f1a55bcc5dc2e9928c34dab78eb2`
- scope: test-only runtime-relative observation time in
  `TestStockRiskServiceClassifiesOversellAndFilters`
- change: replaced the fixed `2026-07-09T12:00:00Z` observation with
  `time.Now().UTC()`; both evidence sources still share the same observation.
- unchanged: production clock and freshness policy, oversell classification,
  quantities, filter input, and assertions.
- side effects: no production, dependency, database, network, provider, or
  Oracle changes or writes.

## Proof

All commands ran from `apps/server_core` with
`GOCACHE=C:\Users\leandro.theodoro\Documents\marketplace-central\.gocache` and
uncached test execution via `-count=1`.

1. `go test ./internal/modules/inventory/application -run '^TestStockRiskServiceClassifiesOversellAndFilters$' -count=1`
   - PASS: `ok marketplace-central/apps/server_core/internal/modules/inventory/application 1.221s`
2. `go test ./internal/modules/inventory/application -count=1`
   - PASS: `ok marketplace-central/apps/server_core/internal/modules/inventory/application 1.082s`
3. `go test ./... -count=1`
   - PASS, exit 0; inventory application passed in `2.293s` and the complete
     `server_core` package set reported no failures.

## Result

Correction proof is PASS. The work item is ready for the Milestone to freeze
the containing commit and request fixed-SHA review. This implementer result does
not itself pass M-09 proportional QA.
