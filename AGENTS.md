# Marketplace Central — Harness Bootstrap

Repository truth is ordered: `ARCHITECTURE.md`/ADRs, OpenAPI plus SDK,
`contracts/governance/`, wiki, `.mnfs/`, then tests/builds/commits. Stop and
classify architecture, contract, runtime, ownership, or verification conflicts.

Portfolio starts one visible Milestone session and passes mission/milestone
paths, the accepted base SHA, knowledge routes, constraints, QA contract, and
its task ID for checkpoints. Milestone owns feature order and starts bounded
Feature workers. A Feature plans and executes in one session by default; after
`plan.md`, compile/validate its context file and read only its selectors. Use
`.agents/skills/mpc-goal-harness/` for the exact session packets and handoffs.
Do not replay transcripts or load broad documents without a named route gap.

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
Milestone sends compact checkpoints to Portfolio, then requests one fixed-SHA
review and proportional QA after all Features integrate. Only QA may pass a
milestone.
