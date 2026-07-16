# Marketplace Central — Harness Bootstrap

Development runs under the hub-and-chips harness. **BINDING doctrine = `docs/HARNESS-CORE.md`
(method) + `docs/HARNESS-PROFILE.md` (this repo's bindings) + the active mission's
`## Parallel Execution Plan` (`.mnfs/MIS-*/mission.md`)** — read core §4–§5 and the profile
before any milestone/feature work (`docs/HARNESS.md` is a pointer). The hub session (boots via
the `harness-hub` skill) authors milestone chips (spawn_task, operator launches on Opus,
worktree isolation), owns merging/deploy/shared seams, and adjudicates parallelism via the
collision matrix. Milestone sessions orchestrate; feature plans come from GPT-5.6 Sol medium,
implementation from GPT-5.6 Luna high (standard) / Sol low (complex) via `/codex:rescue`
(Claude sonnet = sanctioned fallback per core §1), bulk reads from Luna-medium investigators.
Milestone close = dual gate (cold Opus + GPT-5.6 Sol medium, agreement required) then fresh
QA live-drive. Only QA passes a milestone. Chips talk to the hub only via events
(`CLOSED`/`BLOCKED`/`ESCALATION`/`REQUEST`/`SPLIT-REQUEST`/`COMMITTED`/`ACK`).

Repository truth is ordered: `ARCHITECTURE.md`/ADRs, OpenAPI plus SDK,
`contracts/governance/`, wiki, `.mnfs/`, then tests/builds/commits. Stop and
classify architecture, contract, runtime, ownership, or verification conflicts.

Keep domain/application/ports/adapters/transport boundaries; tenant queries
scope `tenant_id`; provider payloads remain at adapters; unknown operational
facts never become zero/default. API changes update OpenAPI and `sdk-runtime`
together. Provider writes need resolved linkage, explicit policy/source time,
duplicate protection, and audit. Mocks prove contract behavior, never live
integration. Live ML writes require explicit operator authorization.

One writer owns a checkout/shared seam. Do not reset, revert, stash, clean,
delete unknown state, use WSL, expose secrets/PII, cold-clone, purge caches, or
install dependencies as a ritual (dep change = `REQUEST` to the hub). Use
`GOCACHE=.gocache` for Go tests. Never push without explicit operator
permission. Evidence lives in `.mnfs` artifacts — unwritten = didn't happen.
