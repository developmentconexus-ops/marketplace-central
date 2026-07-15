# M-01 Dispatch Ledger

Chip session: M-01-listings-read-spine. Branch `mis-003/m-01-listings-read-spine` off base SHA `d0d30d68` (rebased from 1eb99600 after hub committed the harness amendment; docs-only fast-forward).
Isolation: dedicated git worktree `.claude/worktrees/m01-listings` (session re-rooted). Established after the shared-tree hub commit clobbered HEAD onto `main` — chip now isolated, hub owns main.
Harness amendment (operator-ratified 2026-07-15): planning is batch up-front; all codex dispatches via `/codex:rescue --model <m> --effort <e> --wait`; raw `codex exec` only for the precondition probe.

## Precondition
- Codex probe: `codex 0.144.4` present; `CODEX_OK` (gpt-5.6-sol) + `LUNA_OK` (gpt-5.6-luna) returned. Sandbox healthy.

## P2 PLAN

| # | Feature | Role | Model / effort | Path | Result |
|---|---|---|---|---|---|
| D1 | F-01 | Feature planner | gpt-5.6-sol / medium (raw `codex exec` — pre-amendment) | evidence/planner-sol-medium.log → plan.md | DELIVERED — 5 slices, real file:line anchors, 15 open risks; surfaced 2 blocking contradictions |

### D1 findings triage
- **BLOCKING (hub):** migration block 0033–0035 already consumed at base; connectors capability seam lacks price/currency/modality (IC-02-required) → cross-module edit outside listings lane. → BLOCKED event sent to hub.
- **Milestone-session ratified** (documented, non-blocking): price_currency nullable; 409 active-run id at `error.details.operation_run_id`; duplicate-canonical-key → fail-honest before persistence; no cross-module FK (installation existence enforced at app layer); modality allowlist authored here; known-but-not-connected installation → new error case (scoped connected-only for M-01 tests).
- **ESCALATION (out-of-scope):** ADR-12/ADR-17 have no formal record under docs/architecture/decisions (behavior unambiguous; architecture owner repairs).
- **Deferred:** stale queued/running run crash-recovery (manual cleanup for M-01).

| D2 | F-02 | Feature planner | gpt-5.6-sol / medium (streamed to live log; two stuck restarts, see note) | F-02/plan.md | DELIVERED — 32.8KB/219 lines. 7 slices, 7 shared-seam reqs for F-01, real file:line anchors, 5 blockers surfaced (R1–R5) |

### D2 findings triage (post-plan)
- **OPERATOR-RATIFIED:** R2 below_margin formula → contract-literal (DECISIONS D-16).
- **Milestone-ratified:** R3 exception precedence/group severity (D-17); R5 keyset stability (D-18); R4 timeline event source `listing_sync_events` folded into F-01 (D-19).
- **ESCALATION → HUB:** R1 Oracle-cost vs C09 single-query summary (D-20). Blocks only Slice 5 + below_margin counter; Slices 1–4/6/7 proceed.
- **Live-visibility infra:** codex now streamed to a tailable log + local dashboard (http://127.0.0.1:7391) rendering codex events as chat cards. Native task panel does NOT stream codex (wrapper hides child proc → shows "Parado"); log+dashboard is the real live view. Two planner restarts: first raw `codex exec` hung 3.7h (no reasoning output); relaunch via /codex:rescue got interrupted; final clean raw-exec-with-live-log run delivered plan.md at 11:59:50.

### Process deviation (recorded, one-time)
Harness §4.1 (amended) requires ONE batched planner pass per milestone covering ALL features + shared seam. M-01 ran TWO passes (D1 F-01, D2 F-02). Cause: D1 (F-01 planner) was dispatched BEFORE the batch-once amendment was ratified (old one-feature-at-a-time model). Post-amendment the correct move was a single combined pass; instead F-02 was planned separately. Sunk cost — F-01 already planned; a combined re-plan would waste D1. Coherence recovered by feeding F-01/plan.md + DECISIONS.md into D2 as inputs (F-02 designs against F-01's ports; F-02's Shared-Seam-Requirements folded back into F-01 Slice 3). M-02..M-06: single batched planning pass, no repeat.

## Hub rulings applied
- Migration block 0036–0037 granted → 0036_listings.sql (see DECISIONS D-1).
- Connectors capability seam contract-lock (option a), additive-only + CONNECTORS_PROVIDER_AUTH (see DECISIONS D-3). Lock ends at CLOSED; connectors diff called out in CLOSED payload.

## Base / isolation
- Rebased 1eb99600 → d0d30d68 (hub harness amendment, docs-only ff).
- Isolated worktree .claude/worktrees/m01-listings after shared-tree clobber onto main.

## P3 IMPLEMENT
- Starts at F-01 Slice 1 once D2 (F-02 plan) lands and Shared-Seam-Requirements are folded into F-01 repository design. Slices: Luna-high (standard) / Sol-low (complex), failing-test-first, commit per green, independent sonnet review per slice before next.

| # | Feature/Slice | Role | Model / effort | Log | Result |
|---|---|---|---|---|---|
| I1 | F-01 Slice 1 (schema+domain+modules.json) | Implementer | gpt-5.6-luna / high (direct `codex exec`, live log) | scratchpad/agent__f01-slice1.log | IN PROGRESS — dispatched. Scope: 0036_listings.sql (listings + listing_sync_events per D-19), domain listing.go, migration+domain RED-first tests, modules.json register. TDD, GOCACHE absolute, governance green no prefix-exception. |

Note: D-19 folded — `listing_sync_events` table created in the SAME 0036 migration in Slice 1 (schema only; write logic stays Slice 3).

### Live visibility (multi-agent)
Dashboard rebuilt as multi-agent (http://127.0.0.1:7391): sidebar lists every `scratchpad/agent__<id>.log`, operator picks which worker's chat stream to view. Each dispatch writes its own `agent__<id>.log` + `agent__<id>.done` sentinel on exit; server `/agents` reports state live|idle|done. Current agents: f01-slice1 (live), f02-planner (done).

### Below_margin / ICMS (parallel, non-blocking F-01)
ICMS consultation SENT to DB specialist `local_ec787804` (3 questions per hub route-a ruling — see cost-read-verification.md). below_margin (F-02 Slice 2 enrichment + Slice 5 counter) parked pending reply. F-01 Slice 1 has no cost dependency → proceeds.
