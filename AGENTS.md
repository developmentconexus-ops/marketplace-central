# Marketplace Central — Harness Bootstrap

Repository truth is ordered: `ARCHITECTURE.md`/ADRs, OpenAPI plus SDK,
`contracts/governance/`, wiki, `.mnfs/`, then tests/builds/commits. Stop and
classify architecture, contract, runtime, ownership, or verification conflicts.

Portfolio prepares one copyable `/goal` handoff, then the user starts one clean
visible standalone Milestone session manually with the desired root
model/effort. The Milestone role comes from
`.agents/skills/mpc-goal-harness/` plus that handoff. Milestone dispatches the
project custom agent `mpc-implementer`, pinned to `gpt-5.6-luna` with high
reasoning; milestone-end review is the dual gate — Codex review (`mpc-verifier`
config: `gpt-5.6-sol`, medium) plus an independent Claude review — per the
mission execution plan. Milestone requests those agents directly
by their configured names and gives each one the bounded role packet. Portfolio passes
mission/milestone paths, the accepted base SHA, knowledge routes, constraints,
QA contract, and its task ID for the terminal callback. Milestone owns feature
order and starts bounded workers. A Feature plans and executes in one session
by default; after `plan.md`, compile/validate its context file and read only its
selectors. Use the harness skill for exact packets and handoffs. Do not replay
transcripts or load broad documents without a named route gap.

Keep domain/application/ports/adapters/transport boundaries; tenant queries
scope `tenant_id`; provider payloads remain at adapters; unknown operational
facts never become zero/default. API changes update OpenAPI and `sdk-runtime`
together. Provider writes need resolved linkage, explicit policy/source time,
duplicate protection, and audit. Mocks prove contract behavior, never live
integration.

One writer owns a checkout/shared seam. Do not reset, revert, stash, clean,
delete unknown state, use WSL, expose secrets/PII, cold-clone, purge caches, or
install dependencies as a feature ritual. Use `GOCACHE=.gocache` for Go tests.
Each Feature returns one intentional commit and impacted evidence to Milestone.
After all Features integrate, Milestone requests one fixed-SHA review and
proportional QA. Ordinary ambiguities are asked directly in the visible
Milestone. Cross-task messages go to Portfolio only when the user explicitly
requests escalation or on terminal. Only QA may pass a milestone. On terminal,
Milestone persists a compact checkpoint first, then explicitly sends its path
and verdict to the Portfolio task with native cross-task messaging; final
response alone is not a callback.
