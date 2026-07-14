# F-02-observability-baseline

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-01
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-2
lifecycle_scope: feature
```

## Mission

MIS-002 (`../../mission.md`), IC-01, risks R1/R8/R10.

## Milestone

M-01 foundation-observability. Depends on F-01 (instrument the committed surface).

## Brief

Make Oracle reads observable and capture the three baseline facts that gate M-02: (a) per-port-method latency logging + slow-query flag (>1s) + periodic `db.Stats()` logging; (b) via the governed live lane, record active TGFPRO count, median RTT, and the EXPLAIN PLAN verdict for the IC-01 catalog page query.

## Inputs

- `Database` interface already exposes `Stats() sql.DBStats` (unused today).
- Decorator seam: wrap `ports.Reader`/`ports.SankhyaLinkageReader` implementations at composition root with a timing decorator (slog).
- Live lane: `scripts/run-live-oracle-docker.ps1` + `reader_live_test.go` pattern (`//go:build cgo`, `MPC_ORACLE_LIVE_TEST=1`); extend with a baseline subtest — read-only, sanitized output like `_fixed-sha-oracle-evidence.md`.
- Candidate page SQL shape: IC-01 entity (product + stock SUM subquery + price + cost, keyset ORDER BY internal_product_id, FETCH FIRST 50).

## Expected Output

- Timing decorator: `method=<port method> duration_ms=<n> slow_query=<bool>` per call; `pool_stats open=<n> in_use=<n> wait_count=<n>` every 30s (env-tunable).
- Baseline subtest emitting machine-readable lines (e.g. `MPC_BASELINE_TGFPRO_ACTIVE_COUNT=<n>`, `MPC_BASELINE_RTT_MS=<n>`, `MPC_BASELINE_PAGE_PLAN=<INDEX|FULL_SCAN>`).
- One commit.

## Constraints

- No SQL text, bind values, credentials, or raw driver errors in logs (extend `safeOracleCause` discipline).
- No metrics library/Prometheus — slog structured lines only (lean; internal tool).
- Live lane stays read-only; no schema objects created (EXPLAIN PLAN via `DBMS_XPLAN.DISPLAY` after `EXPLAIN PLAN FOR` is acceptable — plan table is session-scoped; if plan table unavailable, record verdict from `V$SQL_PLAN` alternative or mark blocked with the exact Oracle error).

## Negative Scenarios

- While the decorator wraps a reader, when a call fails, the system shall log `method`, `duration_ms`, and the oracle numeric code only — never the raw driver message.
- While the live lane lacks credentials, when the baseline subtest runs, the system shall skip with the existing `MPC_ORACLE_LIVE_TEST` gate message, not fail the suite.

## Validation Expectations

- Log transcript showing latency line per call, forced slow call with `slow_query=true`, and a `pool_stats` line.
- Baseline lines captured into `M-01 validation-result.md` with the three concrete values.
- `GOCACHE=.gocache go test ./...` green (unit path uses fake queryer; live subtest gated).

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` created during feature execution.

## Handoff

- Current status: briefed
- Next owner: Feature Implementer (mpc-implementer, gpt-5.6-luna high)
- Next action: `spec.md` → `plan.md` → implement → evidence
- Required files/evidence: `validation.md`; baseline values surfaced to milestone `validation-result.md`
- Blockers or open decisions: None
