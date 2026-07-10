# M-08-repository-integrity-harness

```yaml
id: M-08
type: milestone
status: planned
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

Marketplace Central has a reconciled, intentionally committed baseline and a
deterministic development harness that isolates unit, ephemeral PostgreSQL,
live integration, browser, and cold-gate execution. A clean validated SHA can
seed independent Codex worktrees without losing current M-03 through M-06
work.

## Why This Milestone Exists

The current checkout contains 88 tracked modifications and 215 untracked
paths. Test execution inherits live `.env` values and can mutate the persistent
development database. New worktrees from HEAD would omit current product
reality. Parallel milestone execution is unsafe until these boundaries are
fixed globally.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | Baseline recovery and truth reconciliation | Classify every dirty path, supersede stale artifacts, and create intentional commits without reset, revert, stash, or loss. |
| F-02 | Hermetic execution lanes | Separate unit, integration, live, browser, and provider-write commands with environment allowlists. |
| F-03 | Ephemeral PostgreSQL and canonical migrations | Give each integration run an isolated database, canonical migrations/seeds, cleanup, and guards against dev/live DB mutation. |
| F-04 | Deterministic cold gate and evidence manifest | Aggregate Go/JS/build/boundary checks and record target-aware redacted evidence by run ID. |
| F-05 | Worktree and session lifecycle | Parameterize stacks, define context packs/handoffs, update runbooks, and create the reusable local Codex skill. |

## Dependencies

- Existing dirty checkout on `main`; its contents are evidence, not disposable
  state.
- Passed M-01 through M-05 artifacts and blocked M-06 evidence.
- No parallel implementation worktree may start before M-08 passes.

## Risks

- Accidental loss or mixing of user-owned changes during baseline recovery.
- A harness-only patch could hide real inventory or contract failures.
- Live credentials could leak into deterministic tests or evidence.
- Fixed Compose names/ports could make later worktrees interfere.

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

## Handoff

- Current status: Planned and approved for dispatch.
- Next owner: Milestone Orchestrator (`gpt-5.6-terra`).
- Next action: Read `execution-guide.md`, create F-01 spec/plan, and execute
  baseline recovery serially in the current checkout.
- Required files/evidence: F-*/validation.md and M-08/validation-result.md.
- Blockers or open decisions: No worktree until F-01 accepts a clean baseline.

## Correction Handoff

- QA failure summary: Not applicable during planning.
- Correction scope: Not applicable.
- Attempts used/remaining: 0/2.
- Next artifact: `validation-result.md` after execution.
- Revalidation evidence required: Criteria M-08-C01 through M-08-C07.

