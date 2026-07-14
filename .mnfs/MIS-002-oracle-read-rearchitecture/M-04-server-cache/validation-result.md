# Milestone Validation Result — M-04 Server Cache

```yaml
id: M-04-validation-result
type: milestone-validation-result
parent: M-04
mission: MIS-002-oracle-read-rearchitecture
validated_sha: 56d5be9ecb02c361b6df889d1baf49ec3b1c48c2
branch: claude/focused-borg-9c9811
validation_level: QA-2
owner: QA Validator
created: 2026-07-14
updated: 2026-07-14
verdict: PASS
supersedes: PASS at df2cac6a31f68c8980ccfa47da51cb40dfbcacee (incomplete; missed 3 evidence-level defects later found by Codex gpt-5.6-sol)
```

## Summary

- Status: Complete
- Validation verdict: **PASS**
- Contract checked: `.mnfs/MIS-002-oracle-read-rearchitecture/M-04-server-cache/validation-contract.md` (M-04-C01..C05)
- Scope: milestone M-04-server-cache at frozen SHA 56d5be9e

All five criteria were re-derived independently from the contract's literal
Evidence lines. Every test cited was mutation-tested by me at this SHA: I broke
the production behavior it guards, confirmed the test FAILED, then restored.
Nine mutations were run; all nine were caught. The three defects that
invalidated the prior PASS are genuinely closed, and each was re-verified by
breaking production behavior rather than by reading the evidence document.

Production `.go` diff vs HEAD is empty at hand-off (`git diff --stat HEAD --
'*.go'` → no output), confirming this round changed tests, evidence, and the
contract only. No git reset/revert/stash/clean was used; every mutation was
reverted by editing the file back.

## Findings

### M-04-C01 — TTL hit/expiry semantics per data class — PASS

- Evidence run: `go test ./internal/modules/internal_read/adapters/cache -run 'FreshnessCache' -count=1 -v` → 11/11 PASS (fake clock).
- Contract specifics verified in production `cache.go:28-31,54-58`: catalog 300s (`5*time.Minute`), stock 45s, price/cost 120s (`2*time.Minute`), env-tunable via `MPC_CACHE_TTL_CATALOG|_STOCK|_PRICECOST`.
- Mutation: removed the TTL expiry check in `lookup` (`c.clock.Now().Sub(item.created) >= ttl`) → `TestFreshnessCacheTTLPerClass` FAILED (`cache_test.go:201: catalog expiry: calls=1 ... third=12:00:00`). Restored; passes.
- Contract blocking failure (stale serve past TTL / wrong as_of propagation): guarded, not observed.

### M-04-C02 — Singleflight collapses concurrent misses — PASS (with residual risk)

- Contract line 51 requires 20 goroutines / same key / cold cache / blocking fake reader; exactly 1 call; all 20 same page + same `as_of`. `TestFreshnessCacheSingleflight` (cache_test.go:256-286) does literally that and passes.
- Mutation: defeated collapsing by making each `DoChan` group key unique.
  - `TestFreshnessCacheErrorNotCached` FAILED **20/20** (`-count=20`) — it has a readiness barrier (all 20 waiters started + 50ms settle) and pins total downstream calls to exactly 2.
  - `TestFreshnessCacheSingleflight` FAILED **0/20** — it only waits for `blockStart` (first entry) before releasing, so the other 19 can hit the now-warm cache and still observe 1 call.
- The contract's blocking failure (">1 call for identical concurrent key") is deterministically guarded by `ErrorNotCached`. validation.md discloses this weakness rather than overstating it. PASS; the weak happy-path detector is recorded as a non-blocking residual risk.

### M-04-C03 — no-cache bypass and linkage exclusion — PASS

The defect that invalidated the prior PASS (two half-proofs that never met) is
closed. `TestComposedCatalogHTTPNoCacheBypassesAndRepopulates`
(`tests/unit/cache_composed_test.go:102`) is genuinely composed: real
`catalogtransport.Handler` → real `cacheadapter.NewCatalogPageReader` → fake
Oracle, driven through `httptest` with a real `Cache-Control: no-cache` header
and asserting on the decoded HTTP response body. It crosses the seam.

Step 3 genuinely proves repopulation with fresh data, and would **not** pass if
bypass returned stale data:

- Mutation A (bypass ignored): `bypass = true` disabled → FAILED `cache_composed_test.go:119: bypass Oracle calls=1, want 2`.
- Mutation B (bypass fetches fresh but skips repopulation): wrapped `storeIfCurrent` in `if !bypass` → composed test FAILED `cache_composed_test.go:130: repopulated as_of=16:00:00, want bypass as_of 16:00:01` — it caught step 3 serving the **stale original**. Under the identical defect the old cache-only `TestFreshnessCacheBypassAndLinkageExclusion` stayed **green**, independently confirming the retraction that `cache_test.go:585` is insensitive and superseded.
- Stale-bypass is also excluded by line 121 (`bypassed.AsOf.After(warm.AsOf)`, strict).
- Mutation C (linkage exclusion — the provider-write-safety blocking failure): added a `Reader.FindProductsForLinking` override routing through the catalog cache → FAILED `cache_test.go:594: linkage was cached: calls=1`. Exclusion is structural (embedded port, only the two catalog methods overridden); the 3-identical-call assertion is sensitive.
- All restored. Contract blocking failures (linkage served from cache; bypass ignored): guarded, not observed.

### M-04-C04 — Bounded memory and observability — PASS

- `defaultMaxEntries = 100_000` matches the mission's ≤100k-products assumption; env-tunable via `MPC_CACHE_MAX_ENTRIES`.
- Mutation A (unbounded growth): disabled the LRU eviction loop → `TestFreshnessCacheLRUAndLogs` FAILED (`cache_test.go:618: cache exceeded LRU cap: size=3`).
- Mutation B (raw key leak): added `"key", "raw-product-id-12345"` to `cacheLog` → FAILED (`cache_test.go:638: log record 1 has unexpected attribute "key"`). The test parses JSON and enforces an exact attribute allowlist `{time, level, msg, cache, key_class}`, so any id/arg/value/digest leak fails.
- Both restored. Contract blocking failure (unbounded growth or raw key values logged): guarded, not observed.

### M-04-C05 — Evict-on-mutation invalidates matching fact class — PASS

The defect that invalidated the prior PASS (entirely unguarded success path, no
composed flow) is closed. Production `InvalidateClass("catalog")` sits at
`assisted_sankhya_linkage_service.go:256`, **after** a successful
`AppendConfirmation`.

- **Confirm is genuinely successful**, not erroring early and passing for the wrong reason: `cache_composed_test.go:223` fails the test on any Confirm error (`t.Fatalf("successful Confirm: %v", err)`), and the flow reaches real persistence via `composedLinkageRepository.AppendConfirmation`.
- **The post-confirm refresh proves EVICTION, not TTL expiry**: the fake clock advances **1s** (line 216) against a **5m** catalog TTL (line 96) — 1s ≪ 5m, so the entry is still live. Corroborated empirically: under Mutation A the post-confirm log emitted `cache=hit key_class=catalog`, proving the entry was still within TTL and that the unmutated second Oracle call can only be eviction.
- **pricecost preservation proves class-scoped invalidation** and is mutation-sensitive: Mutation C (made `InvalidateClass` wipe every entry regardless of class) → FAILED `cache_composed_test.go:231: pricecost Oracle calls=2, want warm entry preserved at 1`.
- Mutation A (delete success-path `InvalidateClass("catalog")`): **both** guards fired — composed FAILED `cache_composed_test.go:236: post-confirm catalog Oracle calls=1, want 2`, and direct FAILED `assisted_sankhya_linkage_service_test.go:430: invalidations=[], want exactly [catalog]`. Deleting the production call no longer leaves the suite green.
- Failed/rolled-back write → no eviction: retained at the application layer (`TestEvictOnMutationFailedWriteDoesNotInvalidate`, `TestApplyManualStockActionDoesNotInvalidateFailedOrRejectedWrites`, `TestAssistedSankhyaConfirmDoesNotInvalidateFailedPersistence`, `TestImportMarginInputsDoesNotInvalidateFailedPersistence`).
- All restored. Contract blocking failure (stale L2 entry served after a successful mutation of that fact class): guarded, not observed.

### Evidence accuracy (F-01 validation.md) — accurate, retractions present

- C02: states `TestFreshnessCacheSingleflight` is a weak guard (~15/100, no waiter-parking barrier) with `TestFreshnessCacheErrorNotCached` as the real guard (20/20). **Confirmed accurate** — my independent measurement was 0/20 vs 20/20, consistent with (and if anything more conservative than) the disclosure.
- C03: states `cache_test.go:585` bypass-repopulation is insensitive and superseded by the composed response-body assertion. **Confirmed accurate** by Mutation B above.
- Both retractions are present (Sixth correction round §3, and the C02/C03 rows of the criteria table). No overstatement found; the document correctly labels itself feature evidence, not a QA verdict.

### Contract Retry Policy — truthful

`correction_attempts: 6`, `max_correction_attempts: 2`, with an explicit
user-authorized overrun disclosure for rounds 3-6. `last_validation_result`
records the df2cac6a PASS **and** the Codex gpt-5.6-sol FAIL at that same SHA,
and states "validation at the new SHA is pending" — it does **not** claim a
result that has not happened. The prior finding is closed.

## Evidence

Artifact paths (absolute):

- `C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\focused-borg-9c9811\.mnfs\MIS-002-oracle-read-rearchitecture\M-04-server-cache\validation-contract.md`
- `C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\focused-borg-9c9811\.mnfs\MIS-002-oracle-read-rearchitecture\M-04-server-cache\F-01-freshness-cache\validation.md`
- `C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\focused-borg-9c9811\apps\server_core\tests\unit\cache_composed_test.go`
- `C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\focused-borg-9c9811\apps\server_core\internal\modules\internal_read\adapters\cache\cache.go`
- `C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\focused-borg-9c9811\apps\server_core\internal\modules\internal_read\adapters\cache\cache_test.go`
- `C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\focused-borg-9c9811\apps\server_core\internal\modules\orders\application\assisted_sankhya_linkage_service.go`
- `C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\focused-borg-9c9811\apps\server_core\internal\modules\orders\application\assisted_sankhya_linkage_service_test.go`

Evidence/commands (all with `GOCACHE=C:\...\focused-borg-9c9811\.gocache`, all `-count=1` unless noted):

- `go test ./... -count=1` → all packages ok, zero FAIL lines.
- `go vet ./...` → clean, no diagnostics.
- `go test ./internal/modules/internal_read/adapters/cache -run 'FreshnessCache' -count=1 -v` → 11/11 PASS.
- `go test ./tests/unit ./internal/modules/orders/application -run 'TestComposedCatalogHTTPNoCacheBypassesAndRepopulates|TestComposedLinkageConfirmEvictsCatalogButPreservesPriceCost|TestAssistedSankhyaConfirmInvalidatesCatalogAfterSuccessfulPersistence' -count=1` → PASS.
- `go test ... -run 'TestFreshnessCacheErrorNotCached' -count=20` and `... -run 'TestFreshnessCacheSingleflight$' -count=20` under the singleflight mutation → 20 fails vs 0 fails.

Mutation ledger (9 mutations, 9 caught, all restored):

| # | Criterion | Production mutation | Guard result |
| --- | --- | --- | --- |
| 1 | C01 | TTL expiry check removed in `lookup` | TTLPerClass FAIL |
| 2 | C02 | `DoChan` group key made unique (no collapsing) | ErrorNotCached FAIL 20/20; Singleflight 0/20 |
| 3 | C03 | bypass detection disabled | Composed FAIL (calls=1, want 2) |
| 4 | C03 | bypass skips `storeIfCurrent` (no repopulation) | Composed FAIL (stale as_of); old cache-only test stayed green |
| 5 | C03 | linkage routed through catalog cache | BypassAndLinkageExclusion FAIL |
| 6 | C04 | LRU eviction loop disabled | LRUAndLogs FAIL (size=3) |
| 7 | C04 | raw `key` attribute added to cache log | LRUAndLogs FAIL (unexpected attribute) |
| 8 | C05 | success-path `InvalidateClass("catalog")` deleted | Composed FAIL + direct FAIL (`invalidations=[]`) |
| 9 | C05 | `InvalidateClass` wipes all classes | Composed FAIL (pricecost calls=2) |

Working-tree integrity: `git diff --stat HEAD -- '*.go'` → empty at hand-off;
`git rev-parse HEAD` → 56d5be9e (frozen). The ' M' entries reported by
`git status` on `cache.go` / `assisted_sankhya_linkage_service.go` /
`stock_action_service.go` are the known Windows CRLF stat artifact — verified:
content diff vs HEAD is zero lines. Not a finding.

Unrun checks: `-race` not run (`CGO_ENABLED=0`, no C compiler). Assessed: does
**not** block any criterion. No M-04 criterion's Evidence line or blocking
failure requires race detection; C02 is satisfied by call-count/`as_of`
assertions under a blocking fake reader. No race coverage is claimed by the
evidence.

Blocking failures: **None observed.**

## Risks

Non-blocking residual risks (recorded, not blocking advancement):

1. `TestFreshnessCacheSingleflight` is a weak regression detector (0/20 under my mutation) because it releases the blocking reader after only the first waiter enters, with no all-waiters-parked barrier. C02's blocking failure is deterministically covered by `TestFreshnessCacheErrorNotCached`, so this is a future-regression sensitivity gap, not a current defect. Recommend adding a readiness barrier to the happy-path test in later work.
2. `cache_test.go:585` remains an insensitive assertion in the tree. It is superseded by the composed C03 test and is now correctly labeled in the evidence; harmless but dead weight.
3. Race detection unavailable in this environment (no C compiler). Concurrency is exercised only via deterministic fake-clock/blocking-reader tests. External environment limitation; recommend race coverage wherever CI provides cgo.
4. Cache invalidation is intentionally per-process; no cross-process bus. Known mission scope, out of M-04.
5. Correction rounds 3-6 exceeded `max_correction_attempts: 2`. Disclosed and user-authorized in the visible Milestone session, and those rounds were test/evidence-only. Process observation for Mission Strategist, not a technical blocker.

## Recommendation

**PASS M-04-server-cache at 56d5be9ecb02c361b6df889d1baf49ec3b1c48c2.**

All five required criteria have sufficient, mutation-verified evidence; no
scope-relevant blocking failure was observed; the production `.go` diff vs the
prior reviewed SHA is empty, so this round changed only tests, evidence, and the
contract. The three findings that invalidated the prior PASS are each closed and
each was re-verified by breaking the production behavior rather than by reading
the claim. No correction scope is recommended; the residual risks above are
recorded as accepted, none requiring rework before advancement.

## Next Handoff

- Next owner: Milestone Orchestrator
- Handoff reason: milestone QA verdict issued; M-04 is cleared to advance.
- Required next inputs: none from QA. Orchestrator may proceed to milestone
  closeout / next milestone. Residual risks 1-3 should be carried forward as
  mission-level notes rather than M-04 rework.
- Verdict is final for this SHA. Any further commit to M-04 invalidates this
  result and requires re-validation at the new SHA.
