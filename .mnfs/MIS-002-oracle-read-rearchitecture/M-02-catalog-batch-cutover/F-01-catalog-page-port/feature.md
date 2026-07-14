# F-01-catalog-page-port

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-02
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-2
lifecycle_scope: feature
```

## Mission

MIS-002. Contract of record: IC-01 (`../../research/catalog-read-interface-contract.md`) — do not restate values, reference them.

## Milestone

M-02 catalog-batch-cutover. Requires M-01 passed; consume the M-01-C04 plan verdict before writing SQL.

## Brief

Replace the catalog listing's entity-composition reads (1+3N) with one use-case-shaped port `ListCatalogProductFacts(ctx, cursor Cursor, limit int) (CatalogFactPage, error)` backed by a single set-based Oracle query: product base + stock aggregate subquery + price + cost joined, keyset-paginated by internal product id (`FETCH FIRST :limit+1 ROWS ONLY` peek for next_cursor). Plus the search variant `SearchCatalogProductFacts(ctx, q string, limit int)` — same single-query shape, text predicate, bounded `FETCH FIRST` (limit 1..50), NextCursor always nil, sort `internal_product_id` asc (IC-01 search row).

## Inputs

- Current hot path: `apps/server_core/internal/modules/catalog/adapters/internalread/reader.go:67-111` (3 sequential calls/product).
- Existing entity ports in `internal_read/ports` and Oracle SQL in `internal_read/adapters/oracle/reader.go` (source of table/view names, column mappings, redaction helpers).
- IC-01 sections: entity shape (nullable facts + quality flags, decimal strings), cursor semantics (base64 of last CODPROD), limit ceiling, chunking rule.
- M-01 `validation-result.md` plan verdict: if `FULL_SCAN`, use the recorded fallback query shape.

## Expected Output

- New port + domain page types in `internal_read` (port layer), Oracle adapter implementation, fake-queryer unit tests (query count, cursor math, null mapping).
- Catalog module consumes the new port for listing; old per-product composition removed from the listing flow (other callers untouched).
- One intentional commit.

## Constraints

- SQL/driver types stay inside `internal_read/adapters/oracle`.
- Unknown facts → nil pointer + quality flag; NEVER zero (mission MIS-02-C02).
- Wrap errors via existing `wrapOracleError`/`safeOracleCause`; no raw driver text.
- Do not modify transport/OpenAPI (F-02 owns that).
- `GOCACHE=.gocache` for tests.

## Inputs/Outputs

Port signature above; `CatalogFactPage{Items []CatalogProductFact, NextCursor *Cursor, AsOf time.Time}`. Field-level shape fixed in IC-01. Cursor invalid → typed domain error the transport maps to 400 `invalid_cursor`.

## Negative Scenarios

- While a cursor decodes to a non-numeric or negative id, when the port is called, the system shall return the typed invalid-cursor error without touching Oracle.
- While a product row lacks a cost row, when the page is built, the system shall set `cost=nil` and append `missing_cost` to quality flags.
- While a product has multiple active price rows, when the page is built, the system shall set `current_price.amount=null` and append `ambiguous_price` to quality flags without failing the page (IC-01 error matrix, per-item).
- While Oracle is unreachable, when the query runs, the system shall return the wrapped source-unavailable error with redacted cause.

## Validation Expectations

- Counting fake queryer: 1 query for pages of 1/50/100 items.
- Cursor chain test: 3 pages non-overlapping, gapless, last page NextCursor nil.
- Null-mapping test per fact type.
- `GOCACHE=.gocache go test ./...` green.

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` at execution time.

## Handoff

- Current status: briefed
- Next owner: Feature Implementer (mpc-implementer, gpt-5.6-luna high)
- Next action: spec → plan → implement → evidence
- Required files/evidence: `validation.md` with transcripts
- Blockers or open decisions: M-01-C04 plan verdict must be read first
