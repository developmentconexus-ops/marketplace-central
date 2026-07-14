# F-01-stock-batch

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-03
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-2
lifecycle_scope: feature
```

## Mission

MIS-002. IC-01 batch rules (chunk 500, semaphore, route classes).

## Milestone

M-03. First feature — also establishes the shared batch helper (chunker + semaphore) reused by F-02.

## Brief

Add batch port `GetStockFactsByIDs(ctx, ids []int64) (map[int64]*StockFact, error)` with IN-list chunking (≤500 per query, ORA-01795 guard) and the batch Oracle semaphore (default 4 permits, `MPC_ORACLE_BATCH_PERMITS`); convert `inventory/application/stock_risk_service.go:29-84` from per-product calls to one batch call.

## Inputs

- Current N+1: `apps/server_core/internal/modules/inventory/application/stock_risk_service.go:29-84`.
- Existing stock SQL in `internal_read/adapters/oracle/reader.go` (TGFEST aggregation) — reuse column mapping.
- M-02 F-01 adapter shape (page port) as the pattern reference.
- IC-01: chunking rule, semaphore spec, batch route class.

## Expected Output

- Port + adapter + shared `oraclebatch` helper (chunk + semaphore, adapter-internal).
- Stock-risk service consumes batch port; missing ids → nil + `missing_stock` flag.
- Semaphore concurrency test (M-03-C05 evidence lives here).
- One intentional commit.

## Constraints

- Semaphore applies to batch-class calls only; interactive ports untouched.
- Unknown stock ≠ 0 — nil + flag (mission MIS-02-C02).
- SQL stays in adapter; map keyed results in adapter, domain types out.
- `GOCACHE=.gocache` for tests.

## Inputs/Outputs

Input ids deduplicated by adapter; output map contains ONLY ids found; caller derives missing set. Chunks executed sequentially within one call (permits bound cross-call concurrency, not intra-call).

## Negative Scenarios

- While ids has 501 entries, when the port runs, the system shall issue exactly 2 queries.
- While ids is empty, when the port runs, the system shall return an empty map with zero Oracle calls.
- While one chunk fails, when the port runs, the system shall return the wrapped error and no partial map.

## Validation Expectations

- Fake-queryer counts for N ∈ {1,500,501,1200}.
- Concurrency test: 8 parallel calls → max 4 in flight.
- Stock-risk service test proves single batch call per run.
- `GOCACHE=.gocache go test ./...` green.

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` at execution time.

## Handoff

- Current status: briefed
- Next owner: Feature Implementer (mpc-implementer, gpt-5.6-luna high)
- Next action: spec → plan → implement → evidence
- Required files/evidence: `validation.md`
- Blockers or open decisions: None
