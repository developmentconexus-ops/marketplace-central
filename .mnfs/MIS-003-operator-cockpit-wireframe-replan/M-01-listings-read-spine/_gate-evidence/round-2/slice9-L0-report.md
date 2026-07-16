# Slice 9 — L0/L1 deterministic pre-pass report (§7)

```yaml
base_sha: d50636e7 (slice-8 evidence commit; code tip a595f36c)
candidate: working-tree diff (uncommitted) — 9 impl files + new migration 0037
gocache: absolute (worktree apps/server_core/.gocache)
buildvcs: false (sidesteps sandbox Git safe.directory VCS-stamp; no repo config changed)
run_by: milestone-owner (chip), independent of the worker
```

## L0 (precedes review dispatch)

- `go build ./...` → **exit 0**
- `go vet ./...` → **exit 0** (whole repo)

## L1 unit (runs ∥ review per §15)

- `go test -count=1 ./internal/modules/listings/... ./internal/composition/... ./internal/platform/migrate/... ./migrations/...` → **exit 0**
  - connectors, integrations, application, domain, ports, transport, composition, platform/migrate, migrations all `ok`.
  - `internal/platform/migrate` green = the two `runner_test.go` 36→37 migration-count guards pass with 0037 present.
- SDK (`packages/sdk-runtime`): `tsc --noEmit` → **exit 0**; `vitest run` → **43/43 PASS** (worker's earlier vitest exit 1 was a transient worktree resolver/access-denied env glitch; re-run clean by milestone-owner).

## L1 integration lane (over 0037) — GREEN

Ephemeral `postgres:16-bookworm`, CREATE DATABASE retry (ok after 12 tries, first-boot 3D000 race defeated). From `apps/server_core`, absolute GOCACHE:

- `go run ./cmd/testdb migrate` #1 → **applied 37 migration(s)** (0037 auto-discovered via `//go:embed *.sql`), exit 0.
- `go run ./cmd/testdb migrate` #2 → **applied 0 migration(s)** (idempotent), exit 0.
- Live constraint verify (`pg_get_constraintdef` for `listings_status_check`):
  `CHECK ((status = ANY (ARRAY['active','paused','closed','unknown','under_review','inactive','payment_required','not_yet_active'])))` — grown set applied at RUNTIME, exact 8 values, name preserved.
- `go test -tags=integration -run TestListingsRead -v -count=1 ./tests/integration` → **PASS**, exit 0 (27.612s):
  - `TestListingsReadPerformance2000` PASS (3.99s) — nearest-rank p95 = 3.17ms, index-only scan, summary conditional-aggregate query count=1.
  - `TestListingsReadContractEndToEnd` PASS (20.85s) — all 8 subtests: small_page cursor walk + JSON null contract, all_filter_keys, q/title/provider_id/seller_sku, by_product cursor tie-order + null-last, detail variation timeline + not_found, error_matrix, null_cost_honesty + known margin + summary, tenant_isolation all read paths + cursors.

This closes the schema/runtime proof the unit lane (source-text-only assertions) cannot: the CHECK DROP+ADD applies cleanly against real Postgres and accepts the grown status set.

## Test-first (RED→GREEN) — worker evidence

RED proofs captured before impl: compile-fail (new consts referenced), mapper table fail (4 statuses → unknown), migration-count fail (36≠37). All GREEN post-impl. Full worker evidence in F-02 validation.md "Slice 9 validation".

## Scope

Exactly 9 impl files (domain listing.go/_test, connectors mapper.go/_test, transport query_test, platform/migrate runner_test, contracts openapi, sdk index.ts) + NEW migrations/0037_listings_status.sql. No root.go, no query.go, no 0036 edit, no F-01 beyond mapper/domain, no main-tree file. `docker/dev/*.sh` dirty = pre-existing hub CRLF env-prep (excluded from commit).
