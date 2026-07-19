# CHIP-SPIKE-T3 — dispatch ledger

Codex dead until 2026-07-25 (quota) → Claude-only lane. All roles below ran as
Agent-tool subagents or in-session Claude work. No codex dispatches.

| # | Role | Model | Purpose | Outcome |
|---|---|---|---|---|
| 1 | recon + build | in-session (Opus 4.8) | port + adapter + orchestration + transport + composition, TDD per slice | all green L0/L1 |
| 2 | governance lane | pwsh `scripts/harness.ps1 -Command governance` | drift check, clean detached worktree, 40-hex BaseSha `72be3aa9…` | `status=passed` |
| 3 | P6 gate A — cold review | Opus (fresh subagent, no context) | holistic correctness/ADR-17/anti-slop/boundary/fixture/error rubric | PASS-WITH-NITS |
| 4 | P6 gate B — adversarial | Sonnet (subagent) | try-to-break: flag drop, nil-panic, fixture lies, boundary leak, error swallow, build | SHIP |

**Dual-gate agreement:** both independently cleared the same six categories; no blocker
from either. Only convergent nit with a code change (fetched_at) was applied (`64adc8cc`).

**Clean detached governance worktree:** `.../Temp/claude/govwt-spike-t3` (created via
`git worktree add --detach 72be3aa9…`, removed after run) — avoids the `.claude/worktrees`
scanner false-fail and the dirty `go.work.sum` in the chip worktree.

**Gotchas hit (memory-confirmed):**
- PowerShell/Bash cwd resets between calls → single-command `cd` + absolute paths.
- Bootstrap `go mod` hit main checkout first → redone with worktree-absolute path.
- `go.work.sum` churn → staged feature files explicitly, did not commit the sum.
