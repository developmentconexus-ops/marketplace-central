# M-01 foundation-observability — Milestone Checkpoint

```yaml
id: M-01
mission: MIS-002-oracle-read-rearchitecture
type: milestone-checkpoint
verdict: PASS
frozen_sha: 4fde7c0e
base_sha: a1d4aedd
closed: 2026-07-14
```

## Verdict

**PASS** — QA Validator (mnfs-workflow:qa-validator) issued binary PASS against `validation-contract.md` C01..C04 on frozen SHA `4fde7c0e`, corroborated independently (re-ran targeted tests, grepped adapter boundary, validated recorded baseline). Independent verifier (mpc-verifier, gpt-5.6-luna, high) had also returned PASS with no defects.

## Integrated commits (base a1d4aedd)

| SHA | Scope |
| --- | --- |
| e0b1a8d7 | F-01 adapter hardening — godror native pool default 12 (env `MPC_ORACLE_POOL_MAX_SESSIONS`), http.Server timeouts, route-class deadlines (interactive 15s / batch 120s → 504 `{"error":"deadline_exceeded"}`), config validation (unknown-never-zero) |
| f51cabca | F-02 observability — timing decorator (`method`/`duration_ms`/`slow_query`), periodic `pool_stats`, gated live baseline probe `TestOracleLiveBaseline` |
| 4fde7c0e | correction — `/admin/fee-schedules/seed` classified batch; pool stats moved behind adapter boundary via neutral `ports.PoolStats` (observability no longer imports `database/sql`) |

## Criteria

- **C01** adapter refactor committed + tests green + redaction (safeOracleCause/wrapOracleError) preserved — PASS
- **C02** pool 12 + http timeouts + route deadlines 15s/120s + seed batch + unknown-never-zero — PASS
- **C03** latency/slow_query/pool_stats logging, no SQL/creds leak, boundary respected (no `database/sql` in observability) — PASS
- **C04** live baseline captured + recorded with M-02 gating verdict — PASS

## C04 baseline (GATES M-02 SQL)

Captured via governed read-only lane `scripts/run-live-oracle-docker.ps1 -EmitBaseline` (whitelist-sanitized). Frozen SHA 4fde7c0e.

| Fact | Value | Consequence for M-02 |
| --- | --- | --- |
| Active product count (`TGFPRO ATIVO='S'`) | **10520** | R10 does NOT fire — below assumed 30k–100k; keyset scales trivially |
| Median Oracle RTT | **26 ms** | R8 confirmed — WAN-class; N+1 costly per round trip, single-query-per-page + batch ports critical, TTLs lean high |
| Catalog page query plan | **FULL_SCAN** | R1 confirmed — **M-02 plan verdict = `fallback: JOIN base TGF* tables`** (NOT index-supported keyset); build page by joining base TGF* tables or complete DBA index review |

## Governance-infra note (uncommitted, outside feature commits)

To capture C04, two shared governance files were extended (operator-authorized 2026-07-14, currently **uncommitted** in working tree):
- `scripts/run-live-oracle-docker.ps1` — added opt-in `-EmitBaseline` switch + whitelist-only sanitizer `Get-LiveOracleBaselineLines` (surfaces only the 3 sanitized baseline markers; default C05 path + output suppression unchanged; redaction discipline preserved).
- `docker/live-oracle/profile.json` — `go_test_regex` set to `^TestOracleLiveBaseline$`.

Decision on committing these shared files deferred to user (governance infra, not milestone feature scope).

## Next (M-02)

M-02-catalog-batch-cutover proceeds with the recorded plan verdict: build the one-query catalog page against base `TGF*` tables (fallback), not the view-based keyset. RTT 26ms → prioritize single-query-per-page + batch ports, raise cache TTLs.
