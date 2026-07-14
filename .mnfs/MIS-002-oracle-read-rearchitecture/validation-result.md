# MIS-002 oracle-read-rearchitecture — Mission Validation Result

```yaml
id: MIS-002-validation-result
type: mission-validation-result
parent: MIS-002
contract: .mnfs/MIS-002-oracle-read-rearchitecture/validation-contract.md
validation_level: QA-2
validated_sha: 314b1ef34bdecb91b5c5695fc7f63068b2149c5a
branch: main (all milestones M-01..M-05 merged)
owner: QA Validator
created: 2026-07-14
verdict: PASS
```

## Summary

- Status: Complete
- Validation verdict: **PASS**
- Contract checked: MIS-02-C01..C05 (all Required, all QA-2)
- All Evidence commands were re-executed live against merged main at
  `314b1ef3`. Every command was green; no blocking-failure line of any
  criterion was observed. Live-Oracle facts are cited only from governed-lane
  outputs already recorded in M-01/M-02 artifacts (per contract: mock evidence
  is never presented as live-Oracle proof; no live Oracle was run in this
  session).
- Go commands ran in `apps/server_core` with
  `GOCACHE="C:/Users/leandro.theodoro/Documents/marketplace-central/.gocache"`
  (absolute; relative `.gocache` fails on this machine).

## Milestone rollups (all PASS)

| Milestone | Rollup | Verdict | Frozen SHA |
| --- | --- | --- | --- |
| M-01 foundation-observability | `M-01-foundation-observability/validation-result.md` | PASS | 4fde7c0e |
| M-02 catalog-batch-cutover | `M-02-catalog-batch-cutover/validation-result.md` | PASS | 648adbdf |
| M-03 batch-inventory-profitability-sankhya | `M-03-batch-inventory-profitability-sankhya/validation-result.md` | PASS | 37930585 |
| M-04 server-cache | `M-04-server-cache/validation-result.md` | PASS | 56d5be9e |
| M-05 web-tanstack | `M-05-web-tanstack/validation-result.md` | PASS (dual review) | c2aea877 |

## Criteria

### MIS-02-C01 — Catalog page is one Oracle query: PASS

Command run (from `apps/server_core`, absolute GOCACHE):

```text
go test ./internal/modules/catalog/... ./internal/modules/internal_read/... -run 'PageQueryCount|CatalogPage' -v
```

Actual: no test named `PageQueryCount` exists; the contract's indicative
pattern is satisfied by the real fake-queryer count test
`TestCatalogPageUsesOneQueryForEveryListSize` (recorded substitution). Observed
verbatim:

```text
=== RUN   TestCatalogPageUsesOneQueryForEveryListSize
=== RUN   TestCatalogPageUsesOneQueryForEveryListSize/limit-1
=== RUN   TestCatalogPageUsesOneQueryForEveryListSize/limit-50
=== RUN   TestCatalogPageUsesOneQueryForEveryListSize/limit-100
--- PASS: TestCatalogPageUsesOneQueryForEveryListSize (0.00s)
--- PASS: TestCatalogPageCursorChainIsGaplessAndNonOverlapping (0.00s)
--- PASS: TestCatalogPageRoutesFollowThreePageCursorChain (0.01s)
ok  	marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle	4.651s
ok  	marketplace-central/apps/server_core/internal/modules/catalog/transport	5.068s
```

The fake queryer records exactly 1 Oracle query per `ListCatalogProductFacts`
page call at limits 1, 50, and 100. The composed httptest
(`TestComposedCatalogHTTPNoCacheBypassesAndRepopulates`, see C05) additionally
proves 1 Oracle call for two warm HTTP page requests.

Baseline (1+3N eliminated): M-01 rollup C04 captured the governed live-Oracle
baseline via `scripts/run-live-oracle-docker.ps1 -EmitBaseline` at frozen SHA
4fde7c0e — "RTT 26 ms is WAN-class … sequential N+1 (today 1+3N) is expensive
per round trip". M-02 rollup C01/C05 records the cutover with governed-lane
evidence: "exactly one query" per page against live Oracle (see
`M-02-catalog-batch-cutover/validation-result.md`, PASS at 648adbdf). The
1+3N pattern no longer exists in the page path: the transport's page reader is
the single-query `catalog_page` adapter (optionally wrapped by cache), proven
by the fake-queryer count above.

Blocking failure observed: **No**.

### MIS-02-C02 — Unknown facts never become zero: PASS

Command run:

```text
go test ./internal/modules/internal_read/... -run 'Nullable|Quality' -v
```

Actual (all PASS, verbatim test names):

```text
--- PASS: TestFakeReaderMissingStockStaysNilWithQualityFlag (0.00s)
--- PASS: TestCatalogPageMapsNullableFactsAndAmbiguousPrice (0.00s)
--- PASS: TestSankhyaLinkageReaderPreservesOneToManyNullableDescendants (0.00s)
--- PASS: TestRequiredQualityFlagsRemainExplicit (0.00s)
--- PASS: TestMissingCostStaysNilWithQualityFlag (0.00s)
--- PASS: TestQualityFlagsAreStable (0.00s)
ok  	marketplace-central/apps/server_core/internal/modules/internal_read/domain	2.411s
ok  	marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle	3.011s
ok  	marketplace-central/apps/server_core/internal/modules/internal_read/adapters/fake	2.378s
```

Missing stock/price/cost stay `nil` (JSON `null`) with explicit quality flags
(`missing_stock`, `missing_cost`, ambiguous-price handling); domain contract
tests pin the flag set. No new batch/page port emits 0 for an unknown fact —
corroborated at milestone level by M-03 rollup C01/C02 ("missing facts remain
nil with `missing_stock`… no zero substitution").

Blocking failure observed: **No**.

### MIS-02-C03 — Layer boundaries and build health: PASS

Commands run:

1. `go build ./...` in `apps/server_core` → `BUILD_OK` (exit 0).
2. `go test ./...` in `apps/server_core` (absolute GOCACHE) → exit 0, every
   package `ok` (full transcript tail shows all modules incl.
   `internal_read/*`, `catalog/*`, `orders/*`, `profitability/*`,
   `tests/unit`, `tests/integration` — zero FAIL lines).
3. `npm run build` in `apps/web` → green: `✓ 1832 modules transformed … ✓
   built in 5.19s` (only benign "use client" bundler warnings from
   node_modules).
4. Boundary grep (ripgrep + `grep -rln`) for `"database/sql"` and `godror`
   across all `apps/server_core/**/*.go`.

Actual grep result — driver imports are confined to exactly one package,
`internal/modules/internal_read/adapters/oracle`:

- `godror`: only `internal_read/adapters/oracle/open_cgo.go`.
- `database/sql`: only `internal_read/adapters/oracle/*.go`
  (`database.go`, `reader.go`, `batch_reader.go`, `stock_batch_reader.go`,
  `catalog_page.go`, `sankhya_linkage_reader.go`, `open_cgo.go`, and their
  `_test.go` files).
- Other oracle-named packages checked: the only oracle adapter directories in
  the tree are `internal_read/adapters/oracle` and its subpackage
  `internal_read/adapters/oracle/oraclebatch`; `oraclebatch` imports neither
  `godror` nor `database/sql`. `orders/adapters/internalread` (Sankhya linkage
  consumer) and `internal_read/observability` use only neutral ports (M-01
  correction 4fde7c0e moved pool stats behind `ports.PoolStats`). No driver
  import exists outside the adapter package.

Blocking failure observed: **No**.

### MIS-02-C04 — No credential or raw driver detail leaks: PASS

Command run:

```text
go test ./internal/modules/internal_read/... -run 'SafeCause|Redact|Leak' -v
```

Actual (all PASS):

```text
--- PASS: TestLoadConfigFromEnvDoesNotLeakSecretValues (0.00s)
--- PASS: TestSankhyaLinkageOracleErrorsRedactDSNCauseAndLogs (0.01s)
ok  	marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle	3.099s
```

Log-statement review of new code (read-only inspection at 314b1ef3):

- Cache layer (`internal_read/adapters/cache/cache.go:255-256`): the only log
  statement is `slog.Info("internal_read.cache", "cache", status,
  "key_class", class)` — status/class fields only, no keys, values,
  credentials, or driver text. `TestFreshnessCacheLRUAndLogs` PASS pins this.
- M-01 instrumentation (`internal_read/observability/timing.go`,
  `pool_stats.go`): `oracle_read` lines carry method/duration attrs with the
  Oracle numeric code extracted via `safeOracleCodePattern`; `pool_stats`
  lines carry only `open`/`in_use`/`wait_count` integers. No
  username/password/connect-string/raw driver text in any statement.
- Sankhya raw-cause gap: fixed in M-03 (rollup C04 — all new Oracle error
  paths use `wrapOracleError`/`sankhyaOracleError` with `safeOracleCause`);
  `TestSankhyaLinkageOracleErrorsRedactDSNCauseAndLogs` re-run green here.
- 503 mapping: `TestCatalogPageRoutesMapSourceAndDeadlineErrors` (run under
  C01) passes with `source_unavailable` / `deadline_exceeded` subtests.

Blocking failure observed: **No**.

### MIS-02-C05 — as_of and force-refresh work end to end: PASS

Command run:

```text
go test ./tests/unit/... -run 'Composed' -v
go test ./internal/modules/internal_read/adapters/cache/... ./internal/modules/internal_read/adapters/oracle/... -run 'FreshnessCache|SankhyaLinkageReaderHasNoFreshnessPolicyOrCachePath' -v
```

Actual (all PASS):

```text
--- PASS: TestComposedCatalogHTTPNoCacheBypassesAndRepopulates (0.01s)
--- PASS: TestComposedLinkageConfirmEvictsCatalogButPreservesPriceCost (0.00s)
ok  	marketplace-central/apps/server_core/tests/unit	4.769s
--- PASS: TestFreshnessCacheTTLPerClass ... TestFreshnessCacheBypassAndLinkageExclusion, TestFreshnessCacheLRUAndLogs (11/11)
--- PASS: TestSankhyaLinkageReaderHasNoFreshnessPolicyOrCachePath (0.00s)
ok  	marketplace-central/apps/server_core/internal/modules/internal_read/adapters/cache	1.399s
```

The composed httptest (`tests/unit/cache_composed_test.go`, from M-04, re-run
now against merged main) exercises the exact contract transcript over real
`catalog/transport` HTTP handlers wired through the real cache adapter:

1. GET `/catalog/products` (miss) → Oracle calls = 1.
2. GET again (hit) → Oracle calls still 1 and `as_of` **equal** to the first
   response's `as_of` (asserted `warm.AsOf.Equal(warmAgain.AsOf)`).
3. Clock advanced, GET with `Cache-Control: no-cache` → Oracle calls = 2 and
   `as_of` **strictly newer** (asserted `bypassed.AsOf.After(warm.AsOf)`),
   then the bypass result repopulates the cache (4th GET stays at 2 calls).

Linkage never cached — structural exclusion confirmed by inspection of
`internal_read/adapters/cache/cache.go`: the cache package defines decorators
only for `CatalogPageReader` and `BatchReader` with classes
`catalog`/`pricecost`/(stock); it contains no linkage type, class, or code
path at all (zero matches for `linkage` in the file).
`TestSankhyaLinkageReaderHasNoFreshnessPolicyOrCachePath` and
`TestFreshnessCacheBypassAndLinkageExclusion` pass, and
`TestComposedLinkageConfirmEvictsCatalogButPreservesPriceCost` proves confirm
invalidates catalog (fresh Oracle read, newer `as_of`) while linkage data is
read straight through. Cross-checked with `M-04-server-cache/validation-result.md`
(PASS at 56d5be9e). Client end (`as_of` + manual refresh with no-cache) covered
by M-05 rollup C03 (PASS, dual review at c2aea877).

Blocking failure observed: **No**.

## Evidence notes and substitutions

- Test-name substitutions recorded: `PageQueryCount` →
  `TestCatalogPageUsesOneQueryForEveryListSize` (C01); batch-side query-count
  behavior additionally exists as
  `TestImportMarginInputsBatchQueryCountN1200` (profitability), green in the
  full suite.
- Live-Oracle evidence is cited from governed lane outputs in M-01/M-02
  rollups only; mocks/fakes here prove contract behavior, never live
  integration.
- Contract file left untouched (status/Actual recorded here per mission
  instruction).

## Residual risks (carried from milestones, accepted, non-blocking)

1. `-race` unavailable: no cgo toolchain on the validation machine, so the
   concurrency suite (singleflight, fences) runs without the race detector.
2. Cache is per-process with no invalidation bus: multi-instance deployments
   can serve divergent `as_of` until TTL; acceptable for current single-node
   topology.
3. Singleflight regression coverage is behavioral but weak (timing-based
   joins); a refactor could regress coalescing without failing a test loudly.
4. `apps/web` TypeScript typecheck carries a pre-existing TS2688 baseline
   (missing type-definition reference) accepted at M-05; `npm run build`
   (vite) is green and is the contract's build gate.

## Blockers

None.

## Handoff

- Next owner: Mission Strategist (mission closeout).
- Verdict: **PASS** — all five criteria pass on live re-run at merged main
  `314b1ef3`; no blocking-failure line observed; residual risks above are
  documented and accepted at milestone level.
