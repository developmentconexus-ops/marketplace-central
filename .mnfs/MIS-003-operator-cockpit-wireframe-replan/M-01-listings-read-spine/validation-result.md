# Milestone Validation Result — M-01-listings-read-spine

```yaml
milestone: M-01-listings-read-spine
mission: MIS-003-operator-cockpit-wireframe-replan
round: 1
review_sha: e2cde3648a6bdd534afc0ae076a08a93d06d7c7a
diff_base: d0d30d68
validation_level: QA-2
status: Blocked
verdict: Blocked
state_transition: NOT APPLIED (no --apply; operator-deferred C10)
```

## Summary

M-01 (IC-02 listings read-spine) — refresh ingestion + five read endpoints + OpenAPI/SDK — passes every code-verifiable acceptance criterion (M01-C01..C09) under a fresh execution-grounded re-run and a five-member cold independent reviewer crew. The sole open item is **M01-C10 (live-provider-read)**, a `Required: Yes` criterion that needs a real connected Mercado Livre installation + OAuth credentials; it was **operator-deferred** this session and honestly recorded `could-not-drive` — no fabricated evidence. Per the fold rule (runnable milestone, no live-driven evidence) the folded verdict is **Blocked**, not Fail: there is no code defect to correct, only missing live evidence.

## Findings

- **P5 verification ladder** (HEAD 54dfdf7): build 0, vet 0, governance-validate PASS, full-tree `go test ./...` 0 (after fixing two M-01-introduced sibling-package regressions, commit 54dfdf7), integration owned-packages green, SDK tsc 0 + 43 vitest, governance-drift RED classified non-blocking (branch-topology + validate-parity, HUB RULING B).
- **P6 dual gate** (SHA e2cde36): full Opus → MILESTONE PASS; GPT-5.6 Sol medium → MILESTONE FAIL on C07 only. Adjudicated PASS under harness truth order — C07 divergence was stale `.mnfs` wording superseded by operator-ratified D-22 + binding OpenAPI (`below_margin_worst_case`). C01–C09 PASS, C10 NOT-VERIFIED/deferred.
- **P7 qa-validator re-run + API live-drive** (fresh ephemeral Postgres, this session): `TEST_EXIT=0` on both lanes; C01–C09 all `reproduced`; C10 `could-not-drive`.
- **P7 cold crew** (5 subagents): ★2,★3,★4,★5,★7 PASS (incl. adversarial ★2 + ★4); ★1 FAIL solely on C10-deferred; ★6 N/A.

## Structural precondition result

`bash scripts/status-integrity.sh .mnfs/MIS-003-operator-cockpit-wireframe-replan` (Git Bash, read-only) →
`STATUS-INTEGRITY OK (6 milestone(s) checked)`, exit 0. No claim/proof, dangling-milestone, or evidence-existence violation. Precondition for a Pass satisfied (not sufficient alone). Post-write attestation not required (no `status: passed` written).

## Contract checked

`validation-contract.md` (QA-2), criteria M01-C01..C10.

## Per-criterion result

| ID | Criterion | Verdict | Evidence |
|----|-----------|---------|----------|
| M01-C01 | Refresh ingestion upserts and closes | **PASS** | `TestListingsRefreshSeedsIC02RowsAndClosesMissing` PASS (fresh lane `rerun-lane.log:75-86`); composite tenant-scoped PK, close-only-removed. |
| M01-C02 | Concurrent refresh guarded | **PASS** | `TestListingsRefreshRejectsConcurrentRunWithActiveID` PASS + `TestOperationRunRepositoryBeginExclusiveIsAtomic` PASS (`rerun-lane-c02.log:68-72`); atomic advisory-lock, 409 with active run id. |
| M01-C03 | Unmappable→unknown/NULL honesty | **PASS** | asserted in C01 test: `sales_30d IS NULL`, `status='unknown'`, `listing_type` NULL, no 0/default (`listings_refresh_test.go:177-183`; `mapper.go:110-119`). |
| M01-C04 | List endpoint contract | **PASS** | `TestListingsReadContractEndToEnd/{small_page_cursor_walk,all_filter_keys,q_*}` PASS; keyset title→provider→variation, filter/q, nullable→JSON null. |
| M01-C05 | By-product grouping (null-last) | **PASS** | `.../by_product_cursor_walk_tie_order_and_null_last` PASS. |
| M01-C06 | Error matrix (status + code) | **PASS** | `.../error_matrix` PASS — 400 installation_required / invalid_filter / invalid_cursor, 404 listing_not_found, 409 refresh_in_progress. |
| M01-C07 | below_margin unknown honesty | **PASS** | `.../null_cost_honesty_known_margin_and_summary` PASS; `belowMargin` returns nil (never false), field `*bool`→JSON null (`read_service.go:430-432`, `read_model.go:136`). Field name `below_margin_worst_case` per ratified D-22 + binding OpenAPI (contract literal `below_margin` is stale — see reconciliation). |
| M01-C08 | OpenAPI/SDK same-commit | **PASS** | commits `77845a59` + `1f0bbc66` each pair openapi.yaml + sdk-runtime same commit (`git show --stat`); P5 recorded `GOV_API_SDK_SPLIT` green (worktree re-run environmentally `GOV_SCHEMA_INVALID`, disclosed). |
| M01-C09 | List perf (Q1) | **PASS** | `TestListingsReadPerformance2000` PASS — fresh p95 **3.1384 ms** (ceiling 500), Index Only Scan `idx_listings_f02_title_key` no Seq Scan, summary aggregate query count=1 (`rerun-lane.log:49-55`). |
| M01-C10 | Live read ingestion (live-provider lane) | **COULD-NOT-DRIVE (deferred)** | Requires real ML connected installation + OAuth; not provisioned this session; operator-deferred (`qa-validator-report.md:39-43`). No fabricated pass. **BLOCKING for close until executed or Required-downgraded.** |

## Crew composition and per-★ fold result

| ★ | Folded | Reviewers |
|---|--------|-----------|
| ★1 coverage | FAIL → Blocked | A (+ union with qa-validator C10 could-not-drive) |
| ★2 evidence honesty | PASS | B + adversarial |
| ★3 verifiability | PASS | B |
| ★4 integration/composition | PASS | C + adversarial |
| ★5 traceability | PASS | A |
| ★6 correction integrity | N/A | not dispatched (round 1, no corrections) |
| ★7 security/tenancy | PASS | C |

Folded review persisted at `milestone-review.md` (round 1). No sub-reviewer FAIL downgraded.

## Re-run corroboration sample

| Criterion | Command | Recorded | Observed | Result |
|-----------|---------|----------|----------|--------|
| C01–C06 | `go test -tags=integration -run "TestListingsRefresh\|TestListingsRead" ./tests/integration` | PASS | PASS | reproduced |
| C02 (atomicity) | `go test -tags=integration -run TestOperationRunRepositoryBeginExclusive ./internal/modules/integrations/adapters/postgres` | PASS | PASS | reproduced |
| C09 | `go test -tags=integration -run TestListingsReadPerformance2000 ./tests/integration` | p95 3.2563 ms | p95 3.1384 ms | reproduced (both < 500 ms) |
| C08 | `git show --stat 77845a59 1f0bbc66` | paired OpenAPI+SDK | paired OpenAPI+SDK | reproduced (substance) |
| C10 | live-provider `POST /listings/refresh` vs real installation | — | not provisioned | could-not-reproduce (operator-deferred) |

Lane: race-proof Docker-assigned-port ephemeral Postgres (`docker run -p 127.0.0.1::5432` + `docker port` discovery + retry-CREATE-DATABASE), 36 migrations, `TEST_EXIT=0`.

## Live runtime validation

- Surface: **API/service** (no apps/web in M-01; frontend is M-02).
- Tool: Go integration httptest over the REAL production composition (real route-class mux + ReadHandler/ReadService + RefreshService/gateway + Postgres repositories over real migrations), request→response→persisted-effect against a fresh integrated DB.
- Outcome: **validated** for C01–C09 (each endpoint driven live against real state). **could-not-drive** for C10 (live-provider surface — no ML creds, operator-deferred).

## Artifact paths

- `milestone-review.md` (folded crew, round 1)
- `_gate-evidence/round-1/rerun-lane.log`
- `_gate-evidence/round-1/rerun-lane-c02.log`
- `_gate-evidence/round-1/qa-validator-report.md`
- `F-01-listings-module-ingestion/validation.md`, `F-02-listings-read-api/validation.md`
- `DISPATCH-LEDGER.md` (P5 ladder + P6 dual gate)

## Blocking failures with defect loci

1. **M01-C10 (Required) has no satisfying evidence** — `validation-contract.md:164-177`; recorded `could-not-drive` at `_gate-evidence/round-1/qa-validator-report.md:39`. NOT a code defect — missing live-provider environment (operator-deferred). This is the sole reason the verdict is Blocked rather than Pass.

No code/composition/security defect found by any reviewer.

## Recommended correction scope

None in code. To reach Pass, EITHER:
- (a) execute the C10 live-provider-read lane against the connected real ML installation (record run id, `SELECT count(*)` tenant-scoped, sanitized sample row, <20% unknown-status) — operator provisions creds/env; OR
- (b) operator/hub ratifies a `Required: Yes`→deferred amendment of C10 in the milestone contract (routes C10 to a dedicated post-merge live lane), then re-gate on C01–C09.

Doc reconciliation (non-blocking): update `validation-contract.md:127` C07 wording `below_margin`→`below_margin_worst_case` to match ratified D-22 + binding OpenAPI.

## Required next inputs

- Live ML connected installation + OAuth credentials (operator-supplied; never entered by the milestone session), OR an operator/hub ruling deferring C10.

## Next handoff

Milestone Orchestrator holds at **Blocked (pending C10)**. No state transition applied (no `--apply`). C10 + doc reconciliation are the only items between here and P8 close. Re-gate is a fresh cold crew + re-run + live pass over the full milestone once C10 evidence exists or C10 is contract-downgraded.

---
**STATUS NOT FLIPPED — re-run with --apply to record the milestone status transition.**
