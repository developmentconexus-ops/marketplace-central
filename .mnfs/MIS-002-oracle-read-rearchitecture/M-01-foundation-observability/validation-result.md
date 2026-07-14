# M-01 foundation-observability — Milestone Validation Result

```yaml
id: M-01
type: milestone-validation-result
parent: MIS-002
frozen_sha: 4fde7c0e
base_sha: a1d4aedd
updated: 2026-07-14
status: qa-pass
qa_verdict: PASS
qa_validator: mnfs-workflow:qa-validator
```

## Integrated commits

| SHA | Scope |
| --- | --- |
| e0b1a8d7 | F-01 adapter hardening: godror pool default 12, http.Server timeouts, route-class deadlines (interactive 15s / batch 120s → 504 `deadline_exceeded`), config validation (unknown-never-zero) |
| f51cabca | F-02 observability: timing decorator (`method`/`duration_ms`/`slow_query`), periodic `pool_stats`, gated live baseline probe (`TestOracleLiveBaseline`) |
| 4fde7c0e | correction: `/admin/fee-schedules/seed` classified batch (operator-authorized); pool stats moved behind adapter boundary (neutral `ports.PoolStats`, observability no longer imports `database/sql`) |

Two user-authored tooling commits (`909f61e6`, `777785e1`) interleave the history; not milestone scope.

## Criteria status (pre-QA, orchestrator + independent verifier)

- **C01** adapter refactor committed + tests green — PASS (fixed-SHA review 4fde7c0e).
- **C02** pool default 12 + route-class deadlines (15s/120s, 504 body); seed batch — PASS.
- **C03** latency / slow-query / pool-stats logging; no SQL/creds leaked; redaction preserved — PASS.
- **C04** baseline evidence (below) — CAPTURED.

Independent verifier (mpc-verifier, gpt-5.6-luna, fixed_sha_review on 4fde7c0e): VERDICT PASS, no defects.

## C04 — Baseline evidence (gates M-02)

Captured via the governed read-only lane `scripts/run-live-oracle-docker.ps1 -EmitBaseline` (dockerized godror + Instant Client, creds from governed `.env`, output whitelist-sanitized to three integer/enum markers). Frozen SHA 4fde7c0e.

| Fact | Value |
| --- | --- |
| Active product count (`SELECT COUNT(*) FROM METALPRD.TGFPRO WHERE ATIVO='S'`) | **10520** |
| Median Oracle RTT | **26 ms** |
| Catalog page query plan verdict | **FULL_SCAN** |

### M-02 decision (recorded, per C04)

- **Plan verdict = `fallback: JOIN base TGF* tables`.** The IC-01 catalog page query (product + stock SUM subquery + price + cost, keyset ORDER BY `internal_product_id`, FETCH FIRST 50) resolves to a FULL_SCAN, not an index-supported keyset. M-02 MUST build the one-query page by joining the base `TGF*` tables directly (or complete a DBA index review) rather than the naive view-based keyset. Confirms risk **R1**.
- **RTT 26 ms** is WAN-class (LAN ≈ 1–5 ms). Confirms risk **R8**: sequential N+1 (today 1+3N) is expensive per round trip; the single-query-per-page and batch ports are even more critical, and L2/TanStack TTLs should lean toward the higher end. Not a blocker for M-01.
- **Count 10,520** active products — below the assumed 30k–100k band. Risk **R10** does NOT fire; keyset pagination scales trivially at this size.

## Governed-lane note

The lane script was extended (opt-in `-EmitBaseline`) to surface only the three whitelisted sanitized baseline markers; the default C05 path and full output suppression are unchanged. Operator-authorized 2026-07-14.

## Handoff

- Next owner: QA Validator (proportional QA against `validation-contract.md` M-01-C01..C04).
- Then: milestone checkpoint + terminal handoff to Portfolio.
