# M-08-repository-integrity-harness

```yaml
id: M-08
type: milestone
status: in_progress
owner: Mission Strategist
parent: MIS-001
created: 2026-07-10
updated: 2026-07-10
validation_level: QA-3
lifecycle_scope: milestone
```

## Mission

MIS-001 Mercado Livre Operating Cockpit.

## Outcome

Marketplace Central has a reconciled baseline and a repository-native
development control plane that compiles bounded context from canonical truth,
routes work by risk, isolates deterministic and live execution, and produces
evidence-backed `/goal` completion. A clean validated SHA can seed independent
Codex worktrees without hidden chat history or mutable runtime collisions.

## Why This Milestone Exists

F-01 recovered the former dirty checkout into intentional commits. F-02 proved
the command taxonomy but also proved that hand-maintained environment denylists
and post-build requirement discovery are not durable. Execution truth remains
duplicated across `.brain`, `.mnfs`, architecture, and wiki. Parallel work is
unsafe until knowledge authority, context compilation, runtime isolation, and
shared-seam ownership are enforced globally.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | Baseline recovery and truth reconciliation | Classify every dirty path, supersede stale artifacts, and create intentional commits without reset, revert, stash, or loss. |
| F-02 | Hermetic execution lanes baseline | Preserve the blocked denylist-based spike as evidence; F-08 replaces its incomplete isolation architecture. |
| F-03 | Ephemeral PostgreSQL and canonical migrations | Give each integration run an isolated database, canonical migrations/seeds, cleanup, and guards against dev/live DB mutation. |
| F-04 | Deterministic cold gate and evidence manifest | Aggregate Go/JS/build/boundary checks and record target-aware redacted evidence by run ID. |
| F-05 | Goal, worktree, and Codex lifecycle | Prove supported project agents/hooks, parameterize runtime identity, and bridge `/goal` through a repo skill without hidden state. |
| F-06 | Knowledge authority cutover | Rehome current ADRs, establish one truth owner per fact, update active guidance, and remove `.brain` without a compatibility layer. |
| F-07 | Governance registry and context compiler | Define validated machine contracts and compile hash-bound, criterion-complete context packs from MNFS and repository truth. |
| F-08 | Hermetic child environment and harness modules | Replace F-02 denylist isolation with fresh allowlisted child environments and a typed modular runner while preserving command compatibility. |
| F-09 | Harness eval and dogfood closure | Run deterministic regression cases, fresh-task/worktree proof, cold gate, and milestone evidence before portfolio reuse. |

## Dependencies

- Existing dirty checkout on `main`; its contents are evidence, not disposable
  state.
- Passed M-01 through M-05 artifacts and blocked M-06 evidence.
- No parallel implementation worktree may start before F-05 proves accepted
  isolation; M-08 shared seams remain single-writer.

## Risks

- Accidental loss or mixing of user-owned changes during baseline recovery.
- A harness-only patch could hide real inventory or contract failures.
- Live credentials could leak into deterministic tests or evidence.
- Fixed Compose names/ports could make later worktrees interfere.
- A registry could become a second writable truth unless authority cutover is
  atomic.
- Codex agent/hook/model assumptions may differ on the installed host; scripts
  remain mandatory fallback until a capability spike passes.

## Done Means

- Every pre-M-08 dirty path is assigned to an intentional commit or explicitly
  retained with owner/reason; no reset/revert/stash is used.
- `git status --short` is empty at the accepted baseline SHA.
- Unit tests run with no `.env`, network, PostgreSQL, Oracle, or provider access.
- PostgreSQL integration tests use a generated `mpc_test_*` database and leave
  the development database unchanged.
- Live Oracle/provider commands are explicit opt-ins and emit target-labelled,
  secret-redacted evidence.
- Cold gate uses the same commands intended for CI and produces a run manifest.
- A native Codex worktree can start from the accepted SHA with an independent
  project/port/database namespace and a passing baseline.
- `.brain` is absent; current ADRs live under `docs/architecture/decisions/`,
  execution status lives only in `.mnfs`, and active guidance names the new
  truth order.
- A valid MNFS feature produces a source-hashed context pack with risk,
  criteria, paths, side effects, commands, evidence types, and stop conditions;
  stale or incomplete packs fail closed.
- Allowed-path and shared-seam conflicts fail before writes; accepted outcomes
  retain base/commit SHA and target-labelled evidence.
- A fresh Codex task can use the repo skill to reconcile `/goal`, dispatch a
  bounded worker, and resume from files without portfolio transcript.
- Deterministic harness eval cases pass; unsupported Codex surfaces have an
  explicit script-based fallback rather than an unproved control.

## Handoff

- Current status: In progress; F-01, F-03, F-06, F-07, and F-08 are accepted;
  F-02 remains the superseded blocked baseline.
- Next owner: Milestone Orchestrator (`gpt-5.6-terra`, reasoning `high`).
- Next action: Dispatch F-04 deterministic cold gate from the accepted F-03
  PostgreSQL and F-08 hermetic execution seams.
- Required files/evidence: F-*/validation.md and M-08/validation-result.md.
- Blockers or open decisions: No parallel implementation writer until F-05
  proves worktree isolation; F-02 cannot be accepted by another denylist patch.

## Correction Handoff

- QA failure summary: Not applicable during planning.
- Correction scope: Not applicable.
- Attempts used/remaining: 0/2.
- Next artifact: `validation-result.md` after execution.
- Revalidation evidence required: Criteria M-08-C01 through M-08-C12.
