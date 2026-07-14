# Milestone Validation Contract

```yaml
id: M-01
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

M-01

## QA Level

QA-2

## Required Outcome

Hardened, observable, committed Oracle foundation + recorded baseline (COUNT, RTT, EXPLAIN PLAN verdict) gating M-02.

## Criteria

## Criterion: Adapter refactor committed and validated
ID: M-01-C01
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: `git log --oneline -3` + `GOCACHE=.gocache go test ./internal/modules/internal_read/... -v`
- Expected: one intentional commit containing the adapter refactor (config/pool/timeouts/Database interface); all adapter tests pass including the 3 new config tests (bounded defaults, governed overrides, unsafe-value rejection)
- Actual:
- Artifact: `M-01-foundation-observability/F-01-adapter-hardening-pool-deadlines/validation.md`
Blocking failure: refactor still uncommitted or tests red
Blocking failure observed: No
Owner: QA Validator

## Criterion: Pool 12 and route-class deadlines active
ID: M-01-C02
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: unit test reading effective config (no env → MaxSessions 12) + httptest against an interactive route with a handler stalled 16s
- Expected: default `PoolMaxSessions=12`; stalled interactive request returns 504 `deadline_exceeded` at ~15s; batch-class route with same stall at 16s returns 200 (120s budget)
- Actual:
- Artifact: F-01 `validation.md`
Blocking failure: no deadline fires, or single deadline applied to batch routes
Blocking failure observed: No
Owner: QA Validator

## Criterion: Latency, slow-query and pool-stats logging
ID: M-01-C03
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: run server with fake reader, issue request, capture slog output
- Expected: one log line per port call with `method=<name> duration_ms=<n>`; a forced >1s call emits `slow_query=true`; a `pool_stats` line with `open/in_use/wait_count` appears within the configured interval; no SQL text or credentials in any line
- Actual:
- Artifact: `M-01-foundation-observability/F-02-observability-baseline/validation.md`
Blocking failure: missing latency lines or leaked SQL/credentials
Blocking failure observed: No
Owner: QA Validator

## Criterion: Baseline evidence recorded (gates M-02)
ID: M-01-C04
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: governed lane `scripts/run-live-oracle-docker.ps1` extended run (read-only): `SELECT COUNT(*) FROM METALPRD.TGFPRO WHERE ATIVO='S'`; RTT sample (repeated 1-row query timing); `EXPLAIN PLAN` for the IC-01 catalog page query
- Expected: three concrete numbers/verdicts written to `validation-result.md`: active product count; median RTT ms; plan verdict = `index-supported keyset` OR `fallback: JOIN base TGF* tables` (decision recorded for M-02)
- Actual:
- Artifact: `M-01-foundation-observability/validation-result.md`
Blocking failure: M-02 dispatched without these three facts
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

Command transcripts + log excerpts; live numbers only from the governed lane.

## Blocking Failures

Per criterion above.

## Retry Policy

- correction_attempts: 0
- max_correction_attempts: 2
- last_validation_result: none

## Handoff

- Current status: planned
- Next owner: Milestone Orchestrator
- Next action: execute features, then QA validation
- Required files/evidence: feature validation.md files, `validation-result.md`
- Blockers or open decisions: None
