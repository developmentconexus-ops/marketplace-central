# F-02 Slice 7 validation

## Scope and composition

Integration tests close M01-C04..C09 for the read-only listings spine. They drive real HTTP requests through `httpx.NewRouteClassMux`, the production `ReadHandler.Register`, `NewReadService`, the real Postgres listings repository over harness migrations, and the real integrations Postgres repository/service/adapter for installation existence. The sole external-data substitution is a deterministic in-process `ports.CostReader` fake at the Oracle boundary; it is explicit fake evidence, not live-Oracle evidence. Because `root.go` deliberately composes an unavailable marketplace-policy source, the tests use the same explicit pricing-policy contract as the listings read-service unit tests.

## Contract evidence

`TestListingsReadContractEndToEnd` covers the six-row real-cursor page walk; exact envelope keys and present-null fields; every filter key (including `exception=sync_error`); title/provider-ID/seller-SKU search; by-product counts, states, tie ordering and null-last grouping; variation detail identity; last-ten timeline ordering; malformed/absent 404s; the status plus `error.code` matrix; ADR-17 null-cost honesty; deterministic known below-margin behavior and summary counters; and cross-tenant isolation across list, q, filter, get, group, summary, and cursor continuation.

## Performance evidence

`TestListingsReadPerformance2000` seeds 2,000 tenant/installation-scoped listing, product-link, and snapshot rows, warms the endpoint, records 100 sequential `limit=50` calls, calculates nearest-rank p95, captures the keyset `EXPLAIN`, and traces the summary endpoint to prove exactly one listings conditional-aggregate query.

- p95: **3.2563 ms** (nearest-rank over 100 sequential `GET /listings?limit=50` calls, 5 warmups) — well under the 500 ms ceiling.
- keyset plan/index line: `Index Only Scan using idx_listings_f02_title_key on listings l` with `Index Cond: ((tenant_id = …) AND (installation_id = …) AND (ROW(title, provider_listing_id, variation_id) > ROW('Title 1000','MLBPERF1000','-')))`, wrapped by `Limit (rows=51)`. No `Seq Scan`. Planner stats primed with `ANALYZE listings` after the 2,000-row bulk load so the plan is deterministic.
- summary listings aggregate query count: **1** (exactly one tenant-scoped `count(*) FILTER` conditional-aggregate query, traced via pgx `QueryTracer` — D-20 one-query rule; the Oracle cost batch is a separate bounded port call, not a per-row query).
- environment/SHA and 100 samples: go=go1.26.4, os=windows/amd64, SHA=(local worktree); 100 samples all in the 1.5–4.1 ms band (min ~1.55 ms, max ~4.03 ms).
- exact lane commands: ephemeral `postgres:16-bookworm` → `CREATE DATABASE mpc_test_<32hex>` → `go run ./cmd/testdb migrate` (applied 36) → `go test -tags=integration -run TestListingsRead -v -count=1 ./tests/integration`.

## Lane result (2026-07-15, milestone-owner run)

`go test -tags=integration -run TestListingsRead -v` → **PASS** (`ok marketplace-central/apps/server_core/tests/integration 5.994s`):
- `TestListingsReadPerformance2000` PASS (0.77s) — p95 3.2563 ms, keyset Index Only Scan, summary query count 1.
- `TestListingsReadContractEndToEnd` PASS (2.06s) — all 8 subtests green: cursor walk + JSON null contract, all filter keys, q (title/provider-id/SKU), by-product tie/null-last, detail timeline + 404, error matrix (status + `error.code`), null-cost honesty + known below-margin + summary counters, tenant isolation across all read paths + cursor.

Real ephemeral Postgres, real migrations (0036), real registrar/handler/repository; the only substitution is the deterministic `ports.CostReader` fake at the Oracle boundary (explicit fake, not live-Oracle evidence). Completes **M01-C04..C09**.

## Slice 8 validation — optional margin-reader degradation

Status: **quick_validation_passed**

Test-first evidence:

- Added failing ceiling-outage and cost-outage cases for `List`, `ByProduct`, and `Get`; dependent `below_margin` pass-through cases; and installation/policy hard-error guards.
- Before the service change, `go test ./internal/modules/listings/application` exited **1** with the expected optional-reader errors (`read ICMS ceiling` / `read listing costs`) and dependent-filter failures.

Implementation evidence:

- Changed only `apps/server_core/internal/modules/listings/application/read_service.go` and `apps/server_core/internal/modules/listings/application/read_service_test.go`.
- Ceiling or cost outage now returns rows/groups/detail with nullable cost and margin facts; ceiling outage skips cost lookup; cost-only detail outage preserves known ICMS rows with null margin results; dependent filters restart from the original cursor and pass candidates through. Installation, policy, repository, and timeline errors remain hard.

Commands and results from `apps/server_core` (absolute cache: `C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\m01-listings\apps\server_core\.gocache`):

- `GOCACHE=$(pwd)/.gocache go test ./internal/modules/listings/...` — **exit 0**, all listings packages passed.
- `GOCACHE=$(pwd)/.gocache go build ./...` — initial **exit 1** from sandbox Git safe-directory VCS stamping; rerun with an in-process `safe.directory=*` Git environment override — **exit 0**. No repository or Git configuration was changed.
- `GOCACHE=$(pwd)/.gocache go vet ./internal/modules/listings/...` — **exit 0**.

Evidence type: `ran` for every command above. No transport, OpenAPI, SDK, migration, adapter, or other-module changes.

### Slice 8 wiring addendum — NO-STUB (composition root)

Per hub NO-STUB doctrine (docs/HARNESS.md §4/§5, ratified 2026-07-15) + hub in-worktree verification,
the corrective slice ALSO removes the composition-root stubs (see
`../HUB-EVENT-ADJUDICATION-cost-policy-wiring.md` FINAL and `slice8-corrective-brief.md` FINAL ADDENDUM):

- `apps/server_core/internal/composition/root.go`: listings cost reader `NewBatchReader(nil,…)` → `NewBatchReader(oracleDB,…)`
  (real Oracle handle, identical to profitability `:474`); policy reader stub `unavailablePolicyService`
  → real Postgres-backed `marketSvc`, keeping the `unavailableListingPolicyReader` degrade wrapper; dead
  `unavailablePolicyService` type + unused `marketplacesdomain` import removed.
- Zero-value `oracleDB` safety (Oracle env absent): `BatchReader.ensureBatchAvailable`
  (`internal_read/adapters/oracle/batch_reader.go:204-212`) returns `ReadErrorSourceUnavailable` before any
  DB call → clean degrade via fix (a), no panic. Verified by inspection + the cost-outage unit cases.
- Verification (chip, absolute GOCACHE, `apps/server_core`): `go build ./...` **exit 0**;
  `go vet ./internal/composition/... ./internal/modules/listings/...` **exit 0**;
  `go test -count=1 ./internal/modules/listings/...` **exit 0** (fresh). Independent cold review (cavecrew): CLEAN.
- Wiring runtime proof deferred to the C10 live re-drive (composition wiring is out-of-band for unit tests).
  No reachable stub remains; NO-STUB satisfied by removal, not deferral (deferral record marked SUPERSEDED).

## Slice 9 validation

Status: **implementation green; SDK Vitest environment-blocked**

Test-first RED proof (before implementation):

- `[ran]` From `apps/server_core`: `GOCACHE=C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\m01-listings\apps\server_core\.gocache go test -count=1 ./internal/modules/listings/...` — **exit 1**. Before domain constants existed, the new domain/mapper test references failed with four `undefined: ListingStatus...` errors; transport also rejected `under_review`, `inactive`, `payment_required`, and `not_yet_active` as invalid status values.
- `[ran]` After Step 2 and before mapper implementation: `go test -count=1 ./internal/modules/listings/domain ./internal/modules/listings/transport ./internal/modules/listings/adapters/connectors` — **exit 1**; domain and transport passed, while mapper reported `status "under_review" = "unknown", want "under_review"`.
- `[ran]` After adding `0037_listings_status.sql` and before count updates: `go test -count=1 ./internal/platform/migrate/...` — **exit 1**; both inventory tests observed 37 migrations while expecting 36.

Files changed for Slice 9:

- `apps/server_core/internal/modules/listings/domain/listing.go`
- `apps/server_core/internal/modules/listings/domain/listing_test.go`
- `apps/server_core/internal/modules/listings/adapters/connectors/mapper.go`
- `apps/server_core/internal/modules/listings/adapters/connectors/mapper_test.go`
- `apps/server_core/internal/modules/listings/transport/query_test.go`
- `apps/server_core/migrations/0037_listings_status.sql` (new)
- `apps/server_core/internal/platform/migrate/runner_test.go`
- `contracts/api/marketplace-central.openapi.yaml`
- `packages/sdk-runtime/src/index.ts`

Validation commands (all marked `ran`):

- `[ran]` `go build ./...` — **exit 1** on the first run due Git VCS safe-directory stamping; rerun with process-only `GIT_CONFIG_*` safe-directory environment override — **exit 0**. No repository config changed.
- `[ran]` `go vet ./...` — **exit 0**, with the same process-only safe-directory environment override.
- `[ran]` `go test -count=1 ./internal/modules/listings/... ./internal/composition/... ./internal/platform/migrate/... ./migrations/...` — **exit 0**.
- `[ran]` From `packages/sdk-runtime`: `npx --no-install tsc --noEmit` — **exit 0**.
- `[ran]` From `packages/sdk-runtime`: `npx --no-install vitest run --config vitest.config.ts` — **exit 1** before test execution; resolver reported access denied reading an ancestor and could not resolve the worktree `vitest.config.ts`.
- `[ran]` From `packages/sdk-runtime`: `npm test` — **exit 1** with the same pre-test Vitest resolver/access-denied failure.
- `[ran]` From `packages/sdk-runtime`: `npx --no-install vitest run --root . --config .\vitest.config.ts` — **exit 1** with the same resolver/access-denied failure; explicit root did not change the environment result.

Integration lane: not run, per the operator instruction that the milestone owner drives ephemeral Postgres and no Docker/server boot is allowed in this worker session.

Locus differences: `docs/HARNESS.md` was not present under this worktree; it was not read from the off-limits main tree. SDK Vitest remains environment/tooling-blocked, with no assertion failure observed. No other plan locus differed.

### Slice 9 — L1 integration lane over 0037 (GREEN)

- Ephemeral postgres:16; migrate #1 = applied 37, #2 = applied 0 (idempotent).
- `listings_status_check` live constraint = `active,paused,closed,unknown,under_review,inactive,payment_required,not_yet_active` (8 values, name preserved).
- `go test -tags=integration -run TestListingsRead` PASS exit 0 (27.6s): ContractEndToEnd (8 subtests) + Performance2000 (p95 3.17ms).
- §14 cold sonnet review: no blocking; single important (lane-evidence gap) CLOSED by this run. Effective APPROVE. See `_gate-evidence/round-2/slice9-review.md`.
