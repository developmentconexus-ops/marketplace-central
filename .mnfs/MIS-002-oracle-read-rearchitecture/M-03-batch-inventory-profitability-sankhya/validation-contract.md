# Milestone Validation Contract

```yaml
id: M-03
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

M-03

## QA Level

QA-2

## Required Outcome

All remaining N+1 Oracle consumers batched; linkage reader lean and redacted.

## Criteria

## Criterion: Stock batch chunks correctly
ID: M-03-C01
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: `GOCACHE=.gocache go test ./internal/modules/inventory/... ./internal/modules/internal_read/... -run 'StockBatch' -v`
- Expected: counting fake queryer shows ceil(N/500) queries for N ∈ {1, 500, 501, 1200}; results keyed by id; ids absent in Oracle → nil fact + `missing_stock` flag (not zero)
- Actual:
- Artifact: `M-03-batch-inventory-profitability-sankhya/F-01-stock-batch/validation.md`
Blocking failure: per-id queries remain or zero substitution
Blocking failure observed: No
Owner: QA Validator

## Criterion: Cost/tax batch adopted by profitability
ID: M-03-C02
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: `GOCACHE=.gocache go test ./internal/modules/profitability/... -run 'Batch|QueryCount' -v`
- Expected: profitability computation over N products issues ceil(N/500) cost + ceil(N/500) tax queries; sales query capped at 5000 rows per IC-01 with `truncated=true` marker when cap hit
- Actual:
- Artifact: `M-03-batch-inventory-profitability-sankhya/F-02-cost-tax-batch/validation.md`
Blocking failure: N+1 remains in profitability path
Blocking failure observed: No
Owner: QA Validator

## Criterion: Linkage reader lean
ID: M-03-C03
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: `GOCACHE=.gocache go test ./internal/modules/internal_read/adapters/oracle/... -run 'TestSankhyaLinkage' -v` (linkage tests live in the oracle adapter package, not product_links)
- Expected: candidate scoring for k candidates = 1 candidates query + 1 IN-list lines query (fake queryer count); no Ping/metadata validation queries in the call path (ValidateConfiguration only at startup)
- Actual:
- Artifact: `M-03-batch-inventory-profitability-sankhya/F-03-sankhya-linkage-lean/validation.md`
Blocking failure: per-candidate queries or per-call validation remain
Blocking failure observed: No
Owner: QA Validator

## Criterion: Linkage redaction gap closed
ID: M-03-C04
Level: Milestone
Type: Security
Required: Yes
Status: Pending
Evidence:
- Command: `GOCACHE=.gocache go test ./... -run 'Redact|SafeCause' -v` including new linkage-reader case
- Expected: forced driver error with DSN-like text in message → returned error and logs contain oracle numeric code only; test asserts absence of username/host/service substrings
- Actual:
- Artifact: F-03 `validation.md`
Blocking failure: raw driver text reachable from linkage path
Blocking failure observed: No
Owner: QA Validator

## Criterion: Batch semaphore bounds concurrency
ID: M-03-C05
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: unit test launching 8 concurrent batch-port calls against instrumented fake
- Expected: max in-flight observed == 4 (env-tunable `MPC_ORACLE_BATCH_PERMITS`, default 4); interactive-class calls unaffected by exhausted permits
- Actual:
- Artifact: F-01 `validation.md` (helper is shared; asserted once)
Blocking failure: unbounded batch concurrency or interactive path blocked by semaphore
Blocking failure observed: No
Owner: QA Validator

## Criterion: ImportMarginInputs ceiling enforced
ID: M-03-C06
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: httptest/unit: import request with `Limit=200` and `Limit=201`
- Expected: `Limit=200` accepted; `Limit=201` → 422 `{"error":"limit_exceeded"}` with the 200 cap named in body and zero Oracle calls (fake queryer count 0)
- Actual:
- Artifact: `M-03-batch-inventory-profitability-sankhya/F-02-cost-tax-batch/validation.md`
Blocking failure: oversized import accepted or silently clamped
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

Command transcripts + fake-queryer counts; redaction assertions verbatim.

## Blocking Failures

Per criterion.

## Retry Policy

- correction_attempts: 0
- max_correction_attempts: 2
- last_validation_result: none

## Handoff

- Current status: planned
- Next owner: Milestone Orchestrator
- Next action: execute F-01 → F-02 → F-03, then QA
- Required files/evidence: feature validation.md files, `validation-result.md`
- Blockers or open decisions: None
