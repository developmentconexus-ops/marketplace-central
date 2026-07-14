# M-03-batch-inventory-profitability-sankhya

```yaml
id: M-03
type: milestone
status: planned
owner: Mission Strategist
parent: MIS-002
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-2
lifecycle_scope: milestone
```

## Mission

MIS-002 (`../mission.md`), IC-01 batch rules (`../research/catalog-read-interface-contract.md`).

## Outcome

Remaining N+1 consumers converted to batch ports: stock-risk (inventory) and profitability read stock/cost/tax facts via chunked batch queries (IN-list ≤500 per IC-01); Sankhya linkage reader loses its per-call ValidateConfiguration ping and per-candidate line N+1, and its raw-error redaction gap is closed. Batch routes run under the 120s batch route class with the Oracle batch semaphore (4 permits).

## Why This Milestone Exists

Catalog (M-02) is the hot path, but `profitability/application/service.go:239-287`, `inventory/application/stock_risk_service.go:29-84` and `sankhya_linkage_reader.go` repeat the same chatty pattern; leaving them keeps pool pressure and slow batch jobs. The linkage reader also leaks raw driver causes (`sankhya_linkage_reader.go:367-369`) — security fix owned here.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | stock-batch | batch port `GetStockFactsByIDs(ctx, ids)` (chunk 500); stock-risk service consumes it — 1..k queries per run instead of N |
| F-02 | cost-tax-batch | batch ports for cost + tax facts by id set; profitability service consumes them; sales-cap 5000 rule per IC-01 |
| F-03 | sankhya-linkage-lean | remove per-call ValidateConfiguration (startup-once), collapse per-candidate line queries to one IN-list query, route raw causes through safeOracleCause; linkage NEVER cached (mission MIS-02-C05) |

## Dependencies

M-01 passed (deadlines/observability in place). No dependency on M-02 code (runs in parallel with M-02 per `../execution-plan.md`; the batch pattern reference is IC-01, not M-02's source).

All three features run IN PARALLEL. The shared chunker API is contract-fixed here so no feature waits for another:

```go
// package internal_read/adapters/oracle/oraclebatch — owned by F-01
// Chunks splits ids preserving order; max is always 500 (IC-01).
func Chunks(ids []int64, max int) [][]int64
```

F-02 and F-03 code against this signature from minute one; the package lands with F-01 and integrates at orchestrator accept time.

## Feature Parallelization

| Lane | Content | Starts |
| --- | --- | --- |
| A | F-01: `oraclebatch` pkg + stock batch port + stock-risk consumption | immediately |
| B | F-02: cost/tax batch ports + profitability consumption + ImportMarginInputs ceiling | immediately (chunker API fixed above) |
| C | F-03: Sankhya linkage lean (startup-once validation, IN-list lines, redaction) | immediately (uses same fixed API) |

Seam ownership (one writer per file): F-01 owns `oraclebatch` + stock port file; F-02 owns cost/tax port files; F-03 owns `product_links` module only. Integration order at accept: A → B → C (B and C rebase on A's chunker before their accept).

## Risks

R3 (semaphore starvation — permits=4, batch class only), R4 (ORA-01795 — chunking rule 500), R6 (redaction regression — F-03 explicit tests).

## Done Means

Stock-risk run over N products issues ceil(N/500) stock queries (fake queryer proof); profitability batch same pattern for cost+tax; linkage candidate scoring issues 1 line query total; no ValidateConfiguration in request path; redaction test proves no raw cause; `GOCACHE=.gocache go test ./...` green.

## Handoff

- Current status: planned
- Next owner: Milestone Orchestrator
- Next action: dispatch F-01, F-02, F-03 simultaneously (three parallel mpc-implementer workers, gpt-5.6-luna high); integrate in order A→B→C
- Required files/evidence: feature `validation.md` files; `validation-result.md`
- Blockers or open decisions: None

## Correction Handoff

Not applicable during initial planning.
