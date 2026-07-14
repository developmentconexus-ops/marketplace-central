# Milestone Validation Contract

```yaml
id: M-02
type: milestone-validation-contract
status: planned
owner: Mission Strategist
parent: MIS-002
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-2
lifecycle_scope: milestone
```

## Milestone ID

M-02

## QA Level

QA-2

## Required Outcome

Catalog page = 1 Oracle query, served through the IC-01 envelope, contracts synchronized.

## Criteria

## Criterion: One Oracle query per page
ID: M-02-C01
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: `GOCACHE=.gocache go test ./internal/modules/internal_read/... ./internal/modules/catalog/... -run 'CatalogPage' -v` (counting fake queryer)
- Expected: query count == 1 for page sizes 1, 50, 100; count independent of item count (no N+1)
- Actual:
- Artifact: `M-02-catalog-batch-cutover/F-01-catalog-page-port/validation.md`
Blocking failure: query count grows with N
Blocking failure observed: No
Owner: QA Validator

## Criterion: Envelope and pagination conform to IC-01
ID: M-02-C02
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: httptest transcript: page 1 (no cursor) → follow `next_cursor` → last page
- Expected: bodies contain exactly `items`, `next_cursor`, `page_size`, `as_of`; `next_cursor` null on last page; item JSON field names/types match IC-01 example (decimal strings, RFC3339 `as_of`); pages non-overlapping and gapless across the cursor chain
- Actual:
- Artifact: `M-02-catalog-batch-cutover/F-02-catalog-api-envelope/validation.md`
Blocking failure: envelope shape or cursor chain deviates from IC-01
Blocking failure observed: No
Owner: QA Validator

## Criterion: Error matrix enforced
ID: M-02-C03
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: httptest: garbage cursor; `limit=0`; `limit=101`; reader returning wrapped Oracle error
- Expected: 400 `invalid_cursor`; 400 `invalid_limit` (both out-of-range cases, allowed range 1..100 stated in body); 503 `source_unavailable` with no driver detail in body
- Actual:
- Artifact: F-02 `validation.md`
Blocking failure: any wrong status/code or leaked driver text
Blocking failure observed: No
Owner: QA Validator

## Criterion: Unknown facts stay null with quality flags
ID: M-02-C04
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: `GOCACHE=.gocache go test ... -run 'Nullable|Quality|Ambiguous' -v` with fixture rows missing stock/price/cost and one product with duplicate active price rows
- Expected: JSON `null` for the missing fact + corresponding `quality` entry (`missing_stock`/`missing_price`/`missing_cost`); duplicate active price rows → `current_price.amount=null` + `ambiguous_price` flag, page still 200; never `0`
- Actual:
- Artifact: F-01 `validation.md`
Blocking failure: zero substitution anywhere in the new port
Blocking failure observed: No
Owner: QA Validator

## Criterion: OpenAPI and sdk-runtime synchronized
ID: M-02-C05
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: `git show --stat <F-02 commit>` + `npm run build` in `packages/sdk-runtime`
- Expected: same commit touches OpenAPI spec + sdk-runtime client (new list method typed with envelope); sdk builds green; old unpaginated catalog listing route removed or marked deprecated per IC-01
- Actual:
- Artifact: F-02 `validation.md`
Blocking failure: route shipped without spec+SDK in same change
Blocking failure observed: No
Owner: QA Validator

## Criterion: Search route bounded with null cursor
ID: M-02-C06
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: httptest: `GET /catalog/products/search?q=PARAFUSO&limit=50`; `limit=51`; fake queryer count
- Expected: 200 envelope with `next_cursor` null and ≤50 items sorted by `internal_product_id` asc from exactly 1 Oracle query; `limit=51` → 400 `invalid_limit` (search range 1..50 per IC-01)
- Actual:
- Artifact: `M-02-catalog-batch-cutover/F-02-catalog-api-envelope/validation.md`
Blocking failure: search paginating (non-null cursor) or unbounded query
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

Command transcripts, httptest bodies verbatim, commit stats. Live-Oracle re-check optional here (M-01 lane owns live evidence).

## Blocking Failures

Per criterion.

## Retry Policy

- correction_attempts: 0
- max_correction_attempts: 2
- last_validation_result: none

## Handoff

- Current status: planned
- Next owner: Milestone Orchestrator
- Next action: execute F-01 → F-02, then QA
- Required files/evidence: feature validation.md files, `validation-result.md`
- Blockers or open decisions: M-01-C04 gate
