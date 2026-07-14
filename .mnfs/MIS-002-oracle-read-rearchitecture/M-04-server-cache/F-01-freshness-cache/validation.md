# Feature Validation

```yaml
id: F-01
type: feature-validation
status: correction_validation_passed
owner: Correction Worker
parent: M-04
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-2
lifecycle_scope: feature
```

## Summary

The fixed-SHA review failures assigned to this correction worker are fixed and
targeted validation passed. This is not a milestone or formal QA verdict.

## Correction Scope and Result

All changes remain on top of `05c2c012`; no reset, revert, stash, clean, merge,
port change, OpenAPI change, sdk-runtime change, or dependency change was made.

1. P1-A stale-after-mutation race: added a per-class generation fence. The
   generation is captured before downstream I/O, invalidation increments it,
   and the compare plus store occurs under one mutex acquisition. A fenced
   load still returns its value to its paying caller but cannot populate cache.
2. P1-B no-cache bypass: MaxAge=0 calls use a separate singleflight namespace,
   so they never receive a normal caller's cached double-check result and still
   repopulate the regular cache key.
3. P1-C cache aliasing: catalog pages and cost, tax, and stock maps now deep
   copy every pointer, slice, and nested `SourceMetadata.ObservedAt` pointer on
   both store and return. Nil maps, nil slices, nil pointers, and nil map values
   remain nil.
4. P2-D waiter cancellation: singleflight waits use `DoChan` and select on the
   waiter context, while live waiters retain one downstream call.
5. P2-E validation evidence: added failed or rolled-back persistence tests for
   inventory, assisted Sankhya linkage, and profitability. Inventory also covers
   provider rejection with a persisted non-Applied state. Claims below name only
   tests that exist and ran.

## Changed Paths

- `apps/server_core/internal/modules/internal_read/adapters/cache/cache.go`
- `apps/server_core/internal/modules/internal_read/adapters/cache/cache_test.go`
- `apps/server_core/internal/modules/inventory/application/stock_action_service_test.go`
- `apps/server_core/internal/modules/orders/application/assisted_sankhya_linkage_service_test.go`
- `apps/server_core/internal/modules/profitability/application/service_test.go`
- `.mnfs/MIS-002-oracle-read-rearchitecture/M-04-server-cache/F-01-freshness-cache/validation.md`

## Commands Run

Every Go test command below used:

```powershell
$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\focused-borg-9c9811\.gocache'
```

- `go test ./internal/modules/internal_read/adapters/cache/ -count=1 -v`
  - Status: Pass
  - Actual excerpts:

    ```text
    --- PASS: TestFreshnessCacheFencesInFlightLoadAfterInvalidation
    --- PASS: TestFreshnessCacheBypassUsesSeparateSingleflightNamespace
    --- PASS: TestFreshnessCacheWaiterContextCancellation
    --- PASS: TestFreshnessCacheDeepCopiesCachedFacts
    --- PASS: TestFreshnessCacheSingleflight
    PASS
    ok   marketplace-central/apps/server_core/internal/modules/internal_read/adapters/cache 1.303s
    ```

- `go test ./... -run FreshnessCache -count=1 -v`
  - Status: Pass
  - Actual excerpt:

    ```text
    --- PASS: TestFreshnessCacheDeepCopiesCachedFacts
    --- PASS: TestFreshnessCacheErrorNotCached
    --- PASS: TestFreshnessCacheBypassAndLinkageExclusion
    --- PASS: TestFreshnessCacheLRUAndLogs
    PASS
    ok   marketplace-central/apps/server_core/internal/modules/internal_read/adapters/cache 5.609s
    ```

- `go test ./... -run EvictOnMutation -count=1 -v`
  - Status: Pass
  - Actual excerpt:

    ```text
    --- PASS: TestEvictOnMutationFailedWriteDoesNotInvalidate
    --- PASS: TestEvictOnMutation
    PASS
    ok   marketplace-central/apps/server_core/internal/modules/internal_read/adapters/cache 6.190s
    ```

- `go test ./internal/modules/inventory/application -run TestApplyManualStockActionDoesNotInvalidateFailedOrRejectedWrites -count=1 -v`
  - Status: Pass
  - Actual excerpt:

    ```text
    --- PASS: TestApplyManualStockActionDoesNotInvalidateFailedOrRejectedWrites/persistence_failure
    --- PASS: TestApplyManualStockActionDoesNotInvalidateFailedOrRejectedWrites/provider_rejection
    PASS
    ```

- `go test ./internal/modules/orders/application -run TestAssistedSankhyaConfirmDoesNotInvalidateFailedPersistence -count=1 -v`
  - Status: Pass
  - Actual excerpt:

    ```text
    --- PASS: TestAssistedSankhyaConfirmDoesNotInvalidateFailedPersistence
    PASS
    ```

- `go test ./internal/modules/profitability/application -run TestImportMarginInputsDoesNotInvalidateFailedPersistence -count=1 -v`
  - Status: Pass
  - Actual excerpt:

    ```text
    --- PASS: TestImportMarginInputsDoesNotInvalidateFailedPersistence
    PASS
    ```

- `go test ./... -count=1`
  - Status: Pass
  - Actual excerpts:

    ```text
    ok   marketplace-central/apps/server_core/internal/modules/internal_read/adapters/cache 3.011s
    ok   marketplace-central/apps/server_core/internal/modules/inventory/application 4.479s
    ok   marketplace-central/apps/server_core/internal/modules/orders/application 3.483s
    ok   marketplace-central/apps/server_core/internal/modules/profitability/application 2.474s
    ok   marketplace-central/apps/server_core/tests/integration 3.698s
    ok   marketplace-central/apps/server_core/tests/unit 3.984s
    ```

- `git diff --check`
  - Status: Pass; no whitespace errors. Git reported only normal LF/CRLF
    normalization warnings for the changed files.

## Evidence Mapped to M-04 Criteria

| Criterion | Evidence actually run | Result |
| --- | --- | --- |
| M-04-C01 | `TestFreshnessCacheTTLPerClass`: fake-clock catalog, inventory, and pricecost TTLs plus environment configuration | Pass |
| M-04-C02 | `TestFreshnessCacheSingleflight`: 20 live waiters and one downstream call; `TestFreshnessCacheErrorNotCached`: failed loads retry; `TestFreshnessCacheWaiterContextCancellation`: canceled waiter returns without waiting for leader | Pass |
| M-04-C03 | `TestFreshnessCacheBypassUsesSeparateSingleflightNamespace`: concurrent normal leader and MaxAge=0 caller make two downstream calls; `TestFreshnessCacheBypassAndLinkageExclusion`: bypass repopulates and linkage remains uncached | Pass |
| M-04-C04 | `TestFreshnessCacheLRUAndLogs`: cap, recency, and redacted structured logs; `TestFreshnessCacheDeepCopiesCachedFacts`: returned mutable fields cannot alias cache state | Pass |
| M-04-C05 | `TestFreshnessCacheFencesInFlightLoadAfterInvalidation`: stale in-flight catalog load is not stored; `TestEvictOnMutation`: successful class eviction and unrelated retention; catalog, inventory, linkage, and profitability failed-write tests: no invalidation after failed persistence; inventory provider rejection: no invalidation for non-Applied state | Pass |

## Risks and Remaining Blockers

- Cache invalidation remains intentionally per-process; no cross-process bus was
  added in this correction.
- No remaining blockers were found in the required targeted or full-suite runs.
- Fixed-SHA re-review and milestone QA remain the responsibility of the
  Milestone Orchestrator; this document does not declare either one passed.

## Handoff

- Status: `correction_validation_passed`
- Original blocking failures resolved: Yes
- Next owner: Milestone Orchestrator / QA Validator
- Handoff reason: all five assigned corrections are implemented with targeted
  evidence and the full `server_core` suite is green.
