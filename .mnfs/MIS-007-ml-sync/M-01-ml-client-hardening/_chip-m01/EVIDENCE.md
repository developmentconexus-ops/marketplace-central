# M-01-ml-client-hardening — Evidence Pack

```yaml
milestone: M-01-ml-client-hardening
mission: MIS-007-ml-sync
base_sha: e6dcf75e8017a0629419505ca9ef30f6c0f189ce
orchestrator: Opus milestone session (dispatch-only, zero production code written by orchestrator)
review_model: lean (operator override) — 1 adversarial sonnet reviewer per feature + this pack; no dual gate
status_as_of: 2026-08-01, uncommitted working tree on top of base_sha
```

Two features dispatched in parallel (disjoint write-sets at dispatch time): F-01
resilience-decorator, F-02 items-multiget-raw-dto. F-01's first pass narrowly wired retry
into one call site (`ReadFeeQuote`) instead of the milestone's actual named choke point
(`doRawWithHeaders`), leaving F-02's multiget/prices readers — the mission's stated reason
this milestone exists — unprotected. Orchestrator caught this via independent code reading
(not from either agent's self-report), directed a correction that widened retry to the real
choke point via an HTTP-method branch (GET retries, non-GET never does), and classified the
resulting file-touch-restriction conflict as a formal amendment to `validation-contract.md`
C6 rather than silently exceeding the declared write-set. A second, smaller defect (per-item
Raw cap documented/spec'd as 1MB but implemented as 256KiB, with a stale comment and a
tautological test) was found by F-02's adversarial review and corrected the same way — code
fix dispatched, spec amended in `feature.md` with rationale, not silently left inconsistent.

## Criterion: M01-C1 — Retry-After honrado com tempo NOMEADO
Status: **PASS**
- Command: `cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/connectors/adapters/mercado_livre/... -run TestResilienceDecoratorHonorsRetryAfterHeader -v -count=1`
- Actual (orchestrator's own run, not the implementer's claim): `--- PASS: TestResilienceDecoratorHonorsRetryAfterHeader (2.39s)`, log line `measured total elapsed = 2.3922754s, inter-call gap = 2.3917459s` — real httptest server returns `429` + `Retry-After: 2` then `200`; assertion is on `time.Now()`-measured elapsed ≥ 2s, not "eventually succeeds."
- Independently reconfirmed by the F-01 adversarial reviewer agent (separate run): `elapsed=2.1704736s, inter-call gap=2.1677997s` (jitter varies run-to-run by design, always ≥ 2s floor).
- Artifact: `resilience_decorator_test.go:61-114`, `TestResilienceDecoratorHonorsRetryAfterHeader`.

## Criterion: M01-C2 — Token-bucket por installation sob concorrência
Status: **PASS**
- Command: `... -run TestResilienceDecoratorTokenBucketThrottlesConcurrentRequests -v -count=1`
- Actual: `--- PASS (0.40s)`, log: `gap[1]=99.6126ms, gap[2]=99.8742ms, gap[3]=100.1954ms, gap[4]=100.1531ms, total elapsed for 5 concurrent requests at 600/min = 400.3972ms (min expected ~400ms)` — assertion is on OBSERVED timestamps from 5 real goroutines, not on config values.
- Configurability confirmed structurally (not hardcoded, fact #11 stays `assumed`): `CapabilityAdapterConfig.RateLimitPerMinute` (`capability_adapter.go:41`) wired into `newResilienceDecorator` (`capability_adapter.go:97`); default (900/min) applied only when `<=0` (`resilience_decorator.go:149-151`).
- Singleton wiring confirmed real (not per-request): `internal/composition/root.go:371` calls `NewCapabilityAdapter` exactly once at process composition root — the bucket map is genuinely process-lifetime shared.
- Gap (non-blocking, registered as debt below): composition root does not currently pass a non-default `RateLimitPerMinute` from env/config — mechanism is configurable, no operator-facing knob wired yet.
- Artifact: `resilience_decorator_test.go:121-189`; `resilience_decorator.go:86-132` (`tokenBucket.Take` — reservation under mutex, sleep outside lock).

## Criterion: M01-C3 — Budget esgotado → erro tipado nomeando o contexto
Status: **PASS**
- Command: `... -run TestResilienceDecoratorBudgetExhaustedReturnsTypedError -v -count=1`
- Actual: `--- PASS (0.13s)`; `MaxRetryAttempts=3`; `domain.ErrorCodeOf(err) == ErrCodeProviderRateLimited`; message contains `"3 attempts"`; `errors.As` recovers `*RateLimitExhaustedError{Attempts:3}`, `LastRetryAfter > 0`; exactly 3 real provider calls observed on the fake server.
- Artifact: `resilience_decorator.go:57-64` (`RateLimitExhaustedError`), `:255-286` (`doRetryable` exhaustion path); `resilience_decorator_test.go:204-252`.

## Criterion: M01-C4 — Multiget 20/batch com Raw populado
Status: **PASS**
- Command: `... -run TestGetItemsMultigetPartitionsInBatchesOf20 -v -count=1`
- Actual: 45 ids → exactly 3 HTTP calls, batch sizes `[20, 20, 5]` observed on the fake transport, DTO order preserved against the global input order (not per-batch). `TestGetItemsMultigetSuccessfulDTOsHaveNonEmptyRaw` confirms non-empty `Raw` per item, byte-identical to that item's own sub-object (not the whole batch array, not empty, not aliased).
- Per-item error isolation tested with a genuinely MIXED batch (not the weaker all-fail/all-succeed shape): `TestGetItemsMultigetPerItemErrorDoesNotFailBatch` — 45 ids, item index 17 fails with 404, other 44 decode cleanly, batch call itself returns no top-level error.
- Boundary cases verified: 0 ids → short-circuit, zero calls (`TestGetItemsMultigetEmptyIDsIsNoOp`); exactly-20/exactly-40 verified by partition-loop logic read (`start += 20` loop exits cleanly at boundary).
- Truncation: implemented per-item cap is **256KiB, not 1MB** as originally stated in `feature.md`/`mission.md` ADR-03 — found by F-02 adversarial review, classified and resolved (see Amendment below), not silently shipped wrong.
- Artifact: `items_multiget_reader.go:169-343`; `items_multiget_reader_test.go`.

## Criterion: M01-C5 — Write marcado no-retry (opt-out provado)
Status: **PASS**
- Command: `... -run TestResilienceDecoratorStockWriteDoesNotRetry -v -count=1`
- Actual: 1 PUT call, `elapsed=67ms` against a `RetryBaseDelay=500ms` config specifically chosen so an accidental retry would be caught by elapsed time — stayed fast, confirming no retry occurred.
- Structural guarantee (not test-luck): `doRawWithHeaders` (`capability_adapter.go:745-748`) routes any `method != http.MethodGet` to `doRawWithHeadersNoRetry` BEFORE the resilience decorator's retry loop is ever reached. All three write call sites converge here: stock write (`capability_adapter.go:466`, direct), price write (`price_writer.go:61`, via `doRawWithIdempotency`), listing write (`listing_writer.go:44`, same path) — all PUT.
- Regression check: `price_writer_test.go`'s pre-existing rate-limit subtest (`calls==1`, unmodified, zero diff) and `listing_writer_test.go`'s 429 case (unmodified, zero diff) both still PASS — the no-retry guarantee holds without needing those tests touched, because it's enforced by the method branch, not by each test's own assumptions.
- Artifact: `capability_adapter.go:745-767`.

## Criterion: M01-C6 — Lanes verdes + freeze do adapter (write-set discipline)
Status: **PASS**
- Command 1 (unit + package suite, fresh): `cd apps/server_core && GOCACHE=$(pwd)/.gocache go build ./... && GOCACHE=$(pwd)/.gocache go vet ./...` → both clean, zero output.
- Command 2: `GOCACHE=$(pwd)/.gocache GOMODCACHE=$(pwd)/.gomodcache go test ./internal/modules/connectors/... -count=1` → `ok` for `melhorenvio` (3.21s), `mercado_livre` (12.10s), `application` (1.93s), `domain` (1.56s), `transport` (2.12s); `events`/`ports`/`readmodel` have no test files (expected).
- Command 3 (hermetic integration lane, per contract's evidence requirement): `npm run harness:integration` (repo root) → `target=ephemeral-postgres`, `migrations=embedded` (72 applied), `status=passed`, `run_id=56b877e5001f457dadd0638014660ba2`, `run_dir=scripts/.runs/56b877e5001f457dadd0638014660ba2`.
- Command 4: `git diff --stat HEAD -- internal/modules/connectors/adapters/mercado_livre/` →
  ```
  capability_adapter.go              | 80 ++++++++++++++++++++--
  items_multiget_reader.go           |  2 +-
  items_multiget_reader_test.go      | 33 +++++++--
  pricing_reader_test.go             | 45 +++++++++++-
  4 files changed, 145 insertions(+), 15 deletions(-)
  ```
  plus untracked `resilience_decorator.go`, `resilience_decorator_test.go` (new files).
- Write-set matches C6's amended file list exactly (see amendment in `validation-contract.md` C6): among pre-existing files, only `capability_adapter.go` + `pricing_reader_test.go` + `items_multiget_reader_test.go` touched; everything else is new files. Zero migrations, zero UI, zero schema changes.
- `git diff HEAD -- internal/modules/connectors/adapters/mercado_livre/capability_adapter.go` at F-02's own commit (`e6dcf75e`) confirmed empty — F-02's Constraints ("zero edits to capability_adapter.go") honored at its own commit boundary, before F-01's later correction touched it.

## Amendments filed (contract conflicts classified, not silently bypassed)

1. **`validation-contract.md` C6** (2026-08-01): Required Outcome names `doRawWithHeaders` as
   THE choke point; C6's literal file list originally allowed only `capability_adapter.go`
   among existing files. Making the choke point actually universal (the Required Outcome's own
   demand) required correcting two pre-existing tests that locked in the WRONG (pre-fix)
   single-attempt behavior. Resolved by amending C6's evidence text to name the two additional
   files, with full rationale — not by touching more files than declared without saying so.
2. **`F-02-items-multiget-raw-dto/feature.md`** (2026-08-01): spec said "1MB" for the per-item
   Raw cap; implementation is 256KiB by deliberate design (the outer whole-response 1MB
   `LimitReader` at `capability_adapter.go:807` is shared across up to 20 items per batch, so a
   true 1MB per-item cap is unreachable/untestable — a >1MB single item would already blow the
   outer cap before this reader sees a parseable array). Resolved by amending the feature spec
   to state 256KiB explicitly with rationale, fixing the code comment that incorrectly said
   "(1MB)", and pinning the constant's value in the test (`if itemMultigetRawCap != 256*1024`)
   so future drift fails loudly instead of silently changing what "oversize" means.

## Registered debts (non-blocking, not contract violations)

- **Test-suite latency footprint**: because the C1/C5 fix is method-based (every GET retries),
  five pre-existing 429 tests across the package that were never touched (`prices_reader`,
  `shipping_reader` ×2, `catalog_offers`, `catalog_identity`) now each burn the full default
  retry budget (~4-8s) before resolving, instead of near-instant single-attempt failure — none
  assert on call count or elapsed time so none broke, but package test wall-clock grew by
  roughly 25s. Production implication: any GET reader hitting persistent 429 now takes several
  seconds to fail where it used to fail fast; sequential-GET readers (catalog identity) pay
  this per call. Not a blocking failure per any stated criterion. Owner for a future decision
  (tighten those five tests' budgets for CI speed, or accept as the correct production
  trade-off): whoever picks up test-suite speed work next; flagged here so it isn't lost.
- **No operator-facing rate-limit config knob**: `RateLimitPerMinute` is real and wired, but
  `internal/composition/root.go` doesn't yet source it from env/config — production runs on the
  900/min default. Fine for this milestone (fact #11 stays `assumed`), but M-04/M-06 (the
  actual backfill/multiget consumers) may want to tune it once real 429 telemetry exists.

## Build/vet final state (orchestrator's own independent verification, not implementer-reported)

```
$ GOCACHE=.gocache go build ./...          → clean
$ GOCACHE=.gocache go vet ./...            → clean
$ go test ./internal/modules/connectors/... -count=1
ok  .../connectors/adapters/melhorenvio   3.210s
ok  .../connectors/adapters/mercado_livre 12.104s
ok  .../connectors/application            1.932s
ok  .../connectors/domain                 1.556s
ok  .../connectors/transport              2.124s
$ npm run harness:integration
status=passed  run_id=56b877e5001f457dadd0638014660ba2
```

## Outstanding before merge (not yet done)

- **M01-U1/M01-U2 (user-drive, browser QA of /anuncios, /pedidos, /precos, /integracoes)**:
  deferred per this contract's own line 19-20 ("comportamento live é certificado
  transitivamente pelos live-drives de M-04/M-06") for the *live provider seam*, but the
  **pre-merge baseline capture** (Evidence Requirements, B-6: a real operation traversing the
  adapter, captured and PII/auth-header-scrubbed BEFORE merge, "irrecuperável depois") has NOT
  been captured — this session cannot boot the dev stack (hub-owned seam per
  HARNESS-PROFILE.md §6). This needs a `REQUEST dev-stack-capture` to the hub before merge, or
  the hub performing the capture directly.
- Milestone-diff cold review (this pack functions as that reviewer's evidence base) — commit
  not yet made; changes are still uncommitted working tree on top of `e6dcf75e`.
