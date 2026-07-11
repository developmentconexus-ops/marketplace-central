# Marketplace Central — Harness Bootstrap

Repository truth is ordered: `ARCHITECTURE.md`/ADRs, OpenAPI plus SDK,
`contracts/governance/`, wiki, `.mnfs/`, then tests/builds/commits. Stop and
classify architecture, contract, runtime, ownership, or verification conflicts.

For an MNFS feature, read the active milestone execution guide, compile and
validate its current context pack, then read only its selectors. Use
`.agents/skills/mpc-goal-harness/` for progressive orchestration detail; do not
replay history or load broad documents without a named route gap.

Keep domain/application/ports/adapters/transport boundaries; tenant queries
scope `tenant_id`; provider payloads remain at adapters; unknown operational
facts never become zero/default. API changes update OpenAPI and `sdk-runtime`
together. Provider writes need resolved linkage, explicit policy/source time,
duplicate protection, and audit. Mocks prove contract behavior, never live
integration.

One writer owns a checkout/shared seam. Do not reset, revert, stash, clean,
delete unknown state, use WSL, expose secrets/PII, cold-clone, purge caches, or
install dependencies as a feature ritual. Use `GOCACHE=.gocache` for Go tests.
Each completed task has one intentional commit and impacted evidence; only QA
may pass a milestone.
