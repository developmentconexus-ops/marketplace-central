# Milestone Validation Contract

```yaml
id: M-01
type: milestone-validation-contract
status: planned
owner: Mission Strategist
parent: MIS-003
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-2
lifecycle_scope: milestone
```

## Milestone ID

M-01-listings-read-spine

## QA Level

QA-2 — integration lane (ephemeral-postgres + stubbed capability) for contracts; one live-provider-read criterion (read-only).

## Required Outcome

IC-02 fully implemented: `listings` table populated by refresh ingestion; five endpoints serving canonical shapes; OpenAPI + sdk-runtime updated same commit.

## Criteria

## Criterion: Refresh ingestion upserts and closes
ID: M01-C01
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: integration test: seed stub capability with 6 listings, refresh, remove one from stub, refresh again
- Expected: first refresh → 6 rows keyed `(tenant_id, installation_id, provider_listing_id, variation_id)`; second refresh → removed row has `status='closed'`, other 5 unchanged
- Actual:
- Artifact: `F-01-listings-module-ingestion/validation.md` test transcript
Blocking failure: missing row, wrong key collision, or absent-row not closed
Blocking failure observed: No
Owner: QA Validator

## Criterion: Concurrent refresh guarded
ID: M01-C02
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: integration test: two `POST /listings/refresh` for same installation, second while first in-flight
- Expected: first → 202 `{"operation_run_id": "<id>"}`; second → 409 body `{"error":{"code":"refresh_in_progress"}}` including the active run id
- Actual:
- Artifact: `F-01.../validation.md`
Blocking failure: duplicate concurrent runs or second request 5xx
Blocking failure observed: No
Owner: QA Validator

## Criterion: Unmappable provider values become unknown
ID: M01-C03
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: unit test feeding adapter an unrecognized modality string and null sales metric
- Expected: row persists `listing_type_code` SQL NULL (IC-02: `listing_type` nullable, no `unknown` code — `unknown` is a `status` value only) and `sales_30d` SQL NULL; no default/zero substitution
- Actual:
- Artifact: `F-01.../validation.md`
Blocking failure: guessed enum value or 0 in place of NULL
Blocking failure observed: No
Owner: QA Validator

## Criterion: List endpoint contract
ID: M01-C04
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: integration tests over IC-02 seed (MLBTEST0001..0006): full page walk; `filter.exception=sync_error`; `q=` search
- Expected: cursor walk yields all 6 exactly once, sorted title ASC then listing_id ASC; exception filter returns only the seeded sync-error row; response items match IC-02 Canonical Examples field-for-field with nullable facts as JSON null
- Actual:
- Artifact: `F-02-listings-read-api/validation.md`
Blocking failure: pagination duplicate/skip, wrong sort, or shape divergence from IC-02
Blocking failure observed: No
Owner: QA Validator

## Criterion: By-product grouping
ID: M01-C05
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: integration test `GET /listings/by-product` over seed containing linked + unlinked listings
- Expected: groups cursor over products; unlinked listings grouped under synthetic null-product group ordered LAST
- Actual:
- Artifact: `F-02.../validation.md`
Blocking failure: unlinked listings dropped or null group not last
Blocking failure observed: No
Owner: QA Validator

## Criterion: Error matrix complete
ID: M01-C06
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: table test covering IC-02 error matrix
- Expected: missing installation → 400 `installation_required`; `filter.bogus=x` → 400 `invalid_filter`; garbage cursor → 400 `invalid_cursor`; unknown composite id → 404 `listing_not_found`; concurrent refresh → 409 `refresh_in_progress` — each asserted on status AND `error.code`
- Actual:
- Artifact: `F-02.../validation.md`
Blocking failure: any row returning different status/code
Blocking failure observed: No
Owner: QA Validator

## Criterion: below_margin unknown honesty
ID: M01-C07
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: integration test on seeded listing with NULL cost
- Expected: item JSON contains `"cost": null` and `"below_margin": null` (not `false`); summary counter `below_margin` excludes it
- Actual:
- Artifact: `F-02.../validation.md`
Blocking failure: `below_margin: false` or 0-cost computation for null cost
Blocking failure observed: No
Owner: QA Validator

## Criterion: OpenAPI/SDK same-commit
ID: M01-C08
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: `git show --stat` of each feature commit touching endpoints
- Expected: every commit adding/changing a `/listings*` path touches `contracts/api/marketplace-central.openapi.yaml` AND `packages/sdk-runtime` in the same commit; governance `GOV_API_SDK_SPLIT` lane green
- Actual:
- Artifact: milestone `validation-result.md` diff excerpts
Blocking failure: endpoint commit without paired OpenAPI+SDK change
Blocking failure observed: No
Owner: QA Validator

## Criterion: List performance target (Q1)
ID: M01-C09
Level: Milestone
Type: Performance
Required: Yes
Status: Pending
Evidence:
- Command: seed 2000 listings rows (script in evidence), run 100 sequential `GET /listings?installation_id=X&limit=50` against dev-local server, record latencies
- Expected: p95 < 500ms; summary endpoint executes as single SQL query (EXPLAIN/query-log proof)
- Actual:
- Artifact: milestone `validation-result.md` latency table + query log
Blocking failure: p95 ≥ 500ms at 2k rows
Blocking failure observed: No
Owner: QA Validator

## Criterion: Live read ingestion (live-provider-read lane)
ID: M01-C10
Level: Milestone
Type: QA
Required: Yes
Status: Pending
Evidence:
- Command: `POST /listings/refresh` against the connected real installation (read-only lane, `scripts/run-live-oracle-docker.ps1` environment conventions)
- Expected: operation run completes `succeeded`; `SELECT count(*) FROM listings` > 0 with tenant_id scoped rows; sample row shows real provider_listing_id `MLB…` and canonical enums (no `unknown` flood: <20% unknown status)
- Actual:
- Artifact: milestone `validation-result.md` live-lane section (run id, counts, sanitized sample row — no tokens)
Blocking failure: run fails, zero rows, or >20% unmapped statuses (adapter mapping gap)
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

Feature quick proofs in `F-0*/validation.md`; QA rollup `validation-result.md` here with fixed SHA + dual-gate records.

## Blocking Failures

Any criterion's blocking failure; tenant_id missing from any new query (Q2) blocks regardless.

## Retry Policy

- correction_attempts: 0
- max_correction_attempts: 2
- last_validation_result: none

## Handoff

- Current status: planned.
- Next owner: QA Validator (after F-01/F-02 accepted by Milestone Orchestrator).
- Next action: execute criteria at fixed SHA.
- Required files/evidence: as above.
- Blockers or open decisions: none.
