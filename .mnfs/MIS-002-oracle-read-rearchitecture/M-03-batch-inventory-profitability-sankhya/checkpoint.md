# M-03 batch-inventory-profitability-sankhya — Milestone Checkpoint

Mission: MIS-002-oracle-read-rearchitecture
Milestone: M-03-batch-inventory-profitability-sankhya
Role: Milestone Orchestrator (visible session, parallel-lane plan)
Hub: `local_00112a95`
Date: 2026-07-14

## Verdict

**QA PASS** — independent mpc-verifier (fixed_sha_review THEN proportional_qa).

- Fixed-SHA review: PASS
- C01–C06: PASS
- Full `go test ./...`: PASS
- Blockers: none

## SHAs

- Accepted base SHA: `a714acf7aa87f4e93d1c44102cc51504a532c7dc`
- Main integrated (M-02 catalog batch): `f55bc5ba` (confirmed ancestor of milestone tip)
- **QA-verified code SHA (reviewed by verifier): `379305857b8e00fd288ea9a9429fafeedc82ac2a`**
- Branch: `claude/unruffled-boyd-48984d` (worktree `.claude/worktrees/sleepy-wing-6d7500`)
- Milestone tip after this checkpoint = one docs-only commit on top of `379305857` (checkpoint + validation-result evidence; no code change).

## What shipped

Parallel-lane execution, integration order A→B→C, one intentional commit per feature, then merge onto main.

- **F-01 stock-batch** (`856c89a8`): `oraclebatch` pkg (`Chunks`, `Semaphore`), stock batch port `GetStockFactsByIDs`, IN-list chunk 500, stock-risk consumption. Unknown facts → nil + `missing_stock`, never zero.
- **F-02 cost-tax-batch** (`3be5bed6`): cost/tax batch ports (ceil(N/500)), profitability consumption, SalesHistory peek@5001 → `truncated=true` (never silent), `ImportMarginInputs.Limit>200 → 422 limit_exceeded` before any Oracle call.
- **F-03 sankhya-linkage-lean** (`a6ca7737`): startup-once ValidateConfiguration, per-candidate N+1 → one IN-list lines query, raw causes through `safeOracleCause`, linkage never cached, canonical-identity preserved (commit 97fd4b58).

## Merge/integration commits

- F-01 merge `7e5e9081`, F-02 merge `56dce72a`, F-03 merge `01c4e31f`, onto-main merge `d4017db7`.
- composition/root.go three-way conflict resolved: single `strings` import; one `oracleBatchPermits` (F-01 robust version with slog.Warn); one shared `oracleBatchSemaphore := oraclebatch.NewSemaphore(oracleBatchPermits(os.Getenv))` injected into both batch stock reader and profitability batch reader; F-03 startup ValidateConfiguration block. Kept M-02 catalog page reader wiring intact (declared conflict point).

## C03 correction (post-first-review)

First verifier review at `d4017db7` returned **FAIL on C03**: F-03 fixed the Oracle adapter but the **orders module** request path still called `ValidateConfiguration` per request (`orders/application/assisted_sankhya_linkage_service.go:140,175,198` + orders wrapper `adapters/internalread/sankhya_linkage_reader.go:39`), issuing Oracle Ping + metadata queries every call. Startup validation alone did not stop them.

Scoped correction (commit `37930585`) — orders sankhya-linkage vertical (same feature vertical, not a foreign milestone seam; orders is the live request path wired via `orderstransport.NewSankhyaLinkageHandler` in root.go):
- Removed the three request-path `ValidateConfiguration` calls.
- Orders wrapper `ValidateConfiguration` now reports the in-memory configuration guard only (no source delegation → no Oracle Ping/metadata).
- Typed states preserved: nil/misconfigured reader → `ConfigurationInvalid`; runtime unavailability surfaces through the actual candidate/descendant reads (`FindCandidates` → `lineageErrorState` → `Unavailable`).
- Tests updated: assert zero request-path validation; inject unavailability through the candidate read.
- Post-fix: only remaining `ValidateConfiguration` non-test call site is `root.go:404` (startup) — matches C03 "ValidateConfiguration only at startup".

Re-review at `379305857` → **PASS** (all criteria + full suite).

## Contract note (non-blocking)

M-03-C03's registered Evidence.Command names `./internal/modules/product_links/... -run 'Linkage'`, which matches **no tests** — the linkage reader/evidence lives in `internal_read/adapters/oracle` (`TestSankhyaLinkageCandidateReadsChunkLinesWithoutRequestValidation`). Substantive C03 evidence was taken from the oracle + orders commands. Recommend correcting the contract command to the actual test location in a follow-up.

## Evidence

- `.mnfs/MIS-002-oracle-read-rearchitecture/M-03-batch-inventory-profitability-sankhya/validation-result.md` (verifier rollup at `379305857`)

## Terminal action

Per hub directive (operator: hub owns main merges): NO self-merge to main. Terminal callback sent to hub `local_00112a95` with verdict + this checkpoint path + branch + QA-verified SHA. Hub merges to main and runs its ladder.
