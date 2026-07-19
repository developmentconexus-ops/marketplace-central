# CHIP-M08-LIST-TIMEOUT — Dispatch Ledger

- **Mission:** MIS-004-mvp-demo · **Milestone:** M-08-pedidos
- **Origin:** REOPEN of CHIP-M08-BUYER (FINDING-M08-LIST-TIMEOUT, demo-critical T-1)
- **Hub:** local_7c41fdd8-dd6f-487c-a7ea-1b511e214fa9
- **Branch:** chip/m08-list-timeout (forked from main tip)
- **BaseSha:** `5ab5714cf2b114ebb156896c213739e7e43662a0`

## Commit

| SHA | Summary |
|-----|---------|
| 73a86c4e | fix(orders): buyer fiscal is drawer-only — list no longer times out (Enrich drops fiscal; EnrichOne detail path; handleGet→EnrichOne) |

## Gate results

- Backend: build OK · vet OK · test all `ok` (orders application/transport/adapters) · gofmt clean (3 files).
- Governance lane: `status=passed` (only pre-existing baseline exceptions).
- P6 dual gate (re-gate of delta): cold Opus PASS + adversarial sonnet PASS → agreement.

## Constraints honored

- READ-only, zero ML writes. No migration. Contract intact (comprador_fiscal DTO unchanged).
- Degrade semantics unchanged (honest absence / warn-once).
- Forked from current main tip (5ab5714c); touched files were byte-identical to the merged chip tip, so no divergence.
- No push. `git branch -d` only. No reset/revert/stash/clean/WSL. All writes to worktree path.

## Out of scope (documented)

- List per-order shipment lookup still sequential — separate optimization, not required to clear the regression.

## Hub owns

- Merge to main.
- P7 re-drive of /pedidos list (no 504) + drawer fiscal/shipping fields = M-08 re-close gate.
