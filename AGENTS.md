# Marketplace Central — Harness Bootstrap

Development runs under the hub-and-chips harness: **`docs/superpowers/HARNESS.md` is BINDING** —
read it before any milestone/feature work. The hub session (boots via the `harness-hub` skill)
authors milestone chips (spawn_task, operator launches on Opus, worktree isolation), owns
merging/deploy/shared seams, and adjudicates parallelism via the collision matrix. Milestone
sessions orchestrate; feature plans come from GPT-5.6 Sol medium, implementation from GPT-5.6
Luna high (standard) / Sol low (complex) via `/codex:rescue`, bulk reads from Luna-medium
investigators. Milestone close = dual gate (full Opus + GPT-5.6 Sol medium, agreement required)
then fresh browser QA. Only QA passes a milestone. Chips talk to the hub only via events
(`CLOSED`/`BLOCKED`/`ESCALATION`/`REQUEST`/`SPLIT-REQUEST`).

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
