# P7 QA-Validator — Execution-Grounded Re-run + API Live-Drive (round-1)

Milestone: **M-01-listings-read-spine** (MIS-003)
Role: execution-grounded corroboration pass (read-only; NEVER fixes defects; only re-runs + records).
Contract: `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-01-listings-read-spine/validation-contract.md` (M01-C01..C10).
Date: 2026-07-15. Worktree SHA head: `b1179767`. Env: go1.26.4, windows/amd64.

## Lanes executed (race-proof ephemeral Postgres, Docker-assigned loopback port)

- **Lane A** (primary) — `docker version`=29.6.1; ephemeral `postgres:16-bookworm` on host port 65251; `CREATE DATABASE` ok after 2 tries; `go run ./cmd/testdb migrate` → applied 36; `GOCACHE=<worktree>\.gocache`.
  `go test -tags=integration -v -count=1 -run "TestListingsRefresh|TestListingsRead|TestOperationRunRepositoryBeginExclusive" ./tests/integration` → **PASS** `ok ... tests/integration 14.426s`. **TEST_EXIT=0**.
  (Note: `TestOperationRunRepositoryBeginExclusive*` lives in the `integrations/adapters/postgres` package, NOT `./tests/integration`, so the given regex did not execute it here — covered by Lane B.)
- **Lane B** (supplemental, C02 atomicity) — ephemeral `postgres:16-bookworm` on host port 52095; migrate applied 36.
  `go test -tags=integration -v -count=1 -run "TestOperationRunRepositoryBeginExclusive" ./internal/modules/integrations/adapters/postgres` → **PASS** `ok ... 4.903s`. **TEST_EXIT=0**.

Full stdout: `_gate-evidence/round-1/rerun-lane.log` (+ `rerun-lane-c02.log`).

Real composition per feature validation.md: real `httpx.NewRouteClassMux` + production `ReadHandler`/`ReadService` + `RefreshHandler`/`RefreshService` + real listings-owned integrations gateway over the published boundary + real integrations Postgres installation/operation-run repos + real listings Postgres repo over harness migrations (0036). Sole external substitutions are deterministic in-process fakes at the provider boundary (`connectors/ports.ListingReader`) and the Oracle cost boundary (`ports.CostReader`) — explicit fake evidence, not live-provider evidence.

## API live-drive summary

- **Surface**: API/service. This IS the live API drive for an API/service milestone — the M-01 integration suite drives real HTTP end-to-end through the production composition over httptest against real ephemeral Postgres.
- **Tool**: Go `-tags=integration` httptest over the real production composition.
- **Outcome**: **validated** for C01..C09 (real request→response→persisted-effect assertions all green against a fresh integrated DB). C10 live-provider surface: **could-not-drive** (deferred, see below).

## Per-criterion re-run table

| ID | Criterion | Command / test | Recorded (F-0x validation.md) | Observed (this re-run) | Verdict |
|----|-----------|----------------|-------------------------------|------------------------|---------|
| C01 | Refresh upserts & closes | `TestListingsRefreshSeedsIC02RowsAndClosesMissing` | PASS (1.69s) — 6 rows on composite PK, `variation_id='-'` default, nullable facts NULL, unrecognized status→`'unknown'`, 2nd refresh closes only removed row | PASS (0.29s) | **reproduced** |
| C02 | Concurrent refresh guarded | `TestListingsRefreshRejectsConcurrentRunWithActiveID` + `TestOperationRunRepositoryBeginExclusiveIsAtomic` | PASS — 202 id A / 409 `refresh_in_progress` w/ active id at `error.details.operation_run_id`; atomicity: exactly one queued run | RejectsConcurrent PASS (0.24s, Lane A); BeginExclusiveIsAtomic PASS (1.91s, Lane B) | **reproduced** |
| C03 | Unmappable → unknown/NULL honesty | asserted inside `TestListingsRefreshSeedsIC02RowsAndClosesMissing` (`sales_30d IS NULL` all 6; unrecognized status persisted `'unknown'`; `listing_type` NULL — no guessed enum, no 0 substitution) | PASS | PASS (0.29s) | **reproduced** |
| C04 | List endpoint contract | `TestListingsReadContractEndToEnd/{small_page_cursor_walk_and_JSON_null_contract, all_filter_keys, q_title_provider_id_and_seller_sku}` | PASS — 6-row cursor walk, title ASC then id ASC, `exception=sync_error` filter, q search, present-null JSON | all 3 subtests PASS | **reproduced** |
| C05 | By-product grouping (null-last) | `TestListingsReadContractEndToEnd/by_product_cursor_walk_tie_order_and_null_last` | PASS — unlinked grouped under synthetic null group ordered LAST | PASS | **reproduced** |
| C06 | Error matrix (status + code) | `TestListingsReadContractEndToEnd/error_matrix` | PASS — installation_required/invalid_filter/invalid_cursor/listing_not_found/refresh_in_progress on status AND `error.code` | PASS | **reproduced** |
| C07 | below-margin unknown honesty (D-22) | `TestListingsReadContractEndToEnd/null_cost_honesty_known_margin_and_summary` | PASS — `cost:null`, honest null (not false), summary counter excludes it | PASS | **reproduced** |
| C08 | OpenAPI/SDK same-commit | `git show --stat` of endpoint commits + governance | endpoint commits pair OpenAPI+SDK; `GOV_API_SDK_SPLIT` green at P5 | git: `77845a59` (POST /refresh) touches `openapi.yaml`+`sdk-runtime` same commit; `1f0bbc66` (4 GET routes) touches both same commit; `cb1c17e1` (summary) registers NO route ("Route NOT registered — Slice 6 owns exposure") — no path added w/o pairing. governance-validate re-run in THIS worktree fails `GOV_SCHEMA_INVALID` on ALL 6 governance JSONs (environmental schema-tooling issue, NOT `GOV_API_SDK_SPLIT`; matches known worktree-checkout governance false-fail) | **reproduced** (substance via git; named governance lane green not re-confirmable here — environmental GOV_SCHEMA_INVALID, P5-green record stands) |
| C09 | List performance (p95<500ms; single-aggregate summary) | `TestListingsReadPerformance2000` | PASS — p95 3.2563ms; keyset Index Only Scan, no Seq Scan; summary aggregate query count=1 | PASS (3.28s) — **p95=3.1384ms** (100 samples 1.03–3.98ms); `Index Only Scan using idx_listings_f02_title_key`, no Seq Scan; **summary conditional-aggregate query count=1** | **reproduced** |
| C10 | Live read ingestion (live-provider lane) | `POST /listings/refresh` vs connected real ML installation | Pending | **could-not-drive** | **could-not-drive** |

## C10 — could-not-drive record

Requires a real connected Mercado Livre installation + valid OAuth credentials, which are NOT provisioned this session. Operator explicitly deferred C10. Reason: "live ML creds + connected installation not provisioned this session — operator-deferred to a dedicated post-P6 live-provider lane before P8 close." Not attempted (hard constraint: do NOT handle/request real provider credentials). No fabricated pass.

## Mismatches / defects

None. Zero mismatches across C01..C09; every re-runnable + high-risk criterion reproduced against a fresh integrated milestone DB. TEST_EXIT=0 on both lanes.

Caveat (not a milestone-code defect): the `harness.ps1 governance-validate` lane in this worktree fails at `GOV_SCHEMA_INVALID` for all six `contracts/governance/*.json` files — a governance schema-tooling/environment condition (consistent with the known worktree-checkout governance false-fail), independent of the C08 `GOV_API_SDK_SPLIT` rule, which did not fire. C08 substance corroborated directly via `git show --stat`.

## Evidence paths

- `_gate-evidence/round-1/rerun-lane.log` — Lane A full stdout (docker run, migrate, `go test -v`, TEST_EXIT) + appended Lane B.
- `_gate-evidence/round-1/rerun-lane-c02.log` — Lane B (BeginExclusive atomicity) stdout.
- `_gate-evidence/round-1/qa-validator-report.md` — this report.
