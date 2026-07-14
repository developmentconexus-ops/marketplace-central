# M-09 C05 Governed Oracle Evidence (closeout re-run)

- frozen_sha: `2eabecbc806d652543b60583e036d2e3e686e10e`
- source: `oracle/sankhya`
- observed_at: `2026-07-14T22:00:25.7498931+00:00`
- read_only: `true`
- positive_codprod_observed: `true`
- command: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/run-live-oracle-docker.ps1`
- status: `passed`
- phase: `complete`
- exit_code: `0`

Run through the default governed lane after the per-mode selector split
(commit f099a52c) and the governance envelope classification (commit
2eabecbc): the default C05 lane executes only
`^TestOracleLiveSmoke$/^product_lookup$` (SELECT-only) and the runner enforced
exit 0 plus the `MPC_C05_POSITIVE_CODPROD_OBSERVED=true` marker. Prior interim
passes at `e75a3b88` (combined selector, dual-review-rejected) and `f099a52c`
are superseded by this record.

Companion `-EmitBaseline` mode proof at `f099a52c`
(2026-07-14T21:42:50Z window; docs/governance-only diff to this SHA):
`status=passed`, `exit_code=0`, `MPC_BASELINE_TGFPRO_ACTIVE_COUNT=10523`,
`MPC_BASELINE_RTT_MS=27`, `MPC_BASELINE_PAGE_PLAN=FULL_SCAN`.

No sensitive, source-row, or personal data is included; container output
remained suppressed. Supersedes the blocked C05 attempt recorded in
`checkpoint.md` (deadlock at `a1d4aedd`) and refreshes
`F-02-oracle-catalog-cutover/_fixed-sha-oracle-evidence.md` (prior pass at
`230dc783`).
