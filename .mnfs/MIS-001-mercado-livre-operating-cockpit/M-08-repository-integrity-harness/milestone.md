# M-08-repository-integrity-harness

```yaml
id: M-08
type: milestone
status: in_progress
owner: Mission Strategist
parent: MIS-001
created: 2026-07-10
updated: 2026-07-11
validation_level: QA-3
lifecycle_scope: milestone
```

## Mission

MIS-001 Mercado Livre Operating Cockpit.

## Outcome

Marketplace Central has a repository-native development control plane that
turns a goal into bounded, restartable, architecture-aware implementation. It
routes exact context, milestone tasks, subagents, leases, tests, reviews, and
real QA without relying on hidden transcript history or VM-like local
reproducibility gates.

## Why This Milestone Exists

Long-running Codex work previously depended on broad rediscovery, duplicated
knowledge, unbounded review loops, and conversation memory. Accepted M-08 work
has already established a clean baseline, canonical knowledge, governance,
hash-current context packs, safe child processes, and ephemeral PostgreSQL.

The F-04 cold-clone experiment demonstrated that simulating a clean machine
inside the Windows checkout costs more than it protects and blocks the product
on host-specific Git behavior. The operator explicitly superseded that model.
The remaining milestone builds the actual development workflow: context,
orchestration, state, risk-selected gates, QA, and dogfood.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | Baseline recovery and truth reconciliation | Accepted clean, intentional baseline without losing user work. |
| F-02 | Hermetic execution lanes baseline | Superseded denylist spike retained as evidence; F-08 owns the accepted process model. |
| F-03 | Ephemeral PostgreSQL and canonical migrations | Accepted disposable integration database, canonical migrations, cleanup, and dev invariance. |
| F-04 | Cold gate experiment | Superseded blocked experiment; evidence is retained and no acceptance depends on it. |
| F-06 | Knowledge authority cutover | Accepted one-owner truth topology and `.brain` retirement. |
| F-07 | Governance registry and context compiler | Accepted schemas, drift checks, hash-current context packs, paths, seams, and risk metadata. |
| F-08 | Safe child runtime and harness modules | Accepted allowlisted environments, subprocess execution, redaction, and stable command surface. |
| F-10 | Pragmatic harness cutover | Remove cold-only runtime and replace it with a current-checkout, task-declared impact/evidence gate. |
| F-05 | Goal and orchestration control plane | Add knowledge routes, state/leases/resume, goal skill, native task workflow, bounded built-in subagents, and optional worktree dispatch. |
| F-09 | Harness eval and dogfood closure | Prove context efficiency, deterministic guards, fresh-task orchestration, steering, review, QA routing, and artifact-only resume. |

## Dependencies

- F-01, F-03, F-06, F-07, and F-08 are accepted foundations.
- F-02 and F-04 are historical superseded experiments and do not gate V1.
- F-10 precedes F-05 so no new control plane depends on cold-only vocabulary.
- F-05 precedes F-09.
- Parallel product milestones require accepted F-05 lease/worktree behavior;
  read-only planning and research may already run in parallel.

## Risks

- Overloading `AGENTS.md` or context packs could recreate context bloat.
- Automatically opening too many tasks/subagents could increase tokens and
  coordination cost instead of reducing them.
- Narrative knowledge routes could drift unless paths and selectors are
  mechanically validated.
- Hooks or experimental app-server methods could become accidental hard
  dependencies.
- Parallel writers could collide on OpenAPI/SDK, migrations, composition,
  dependency locks, ADRs, or provider capability contracts.
- A convenient fake/preflight could be misreported as live QA.

## Done Means

- A goal is reconciled to one MNFS mission, milestone, and next eligible feature.
- The total harness-requested initial read set—bootstrap plus pack selectors—
  targets at most 2,000 estimated tokens and names exact interfaces, paths,
  seams, criteria, commands, risk, and stop conditions. Necessary L2/L3 overflow
  is source-justified; stale packs fail.
- A visible milestone task can dispatch bounded depth-one subagents, receive
  compact results, steer or interrupt them, and resume from repository state.
- One writer owns each checkout and shared seam; competing or out-of-scope
  writes fail before acceptance.
- L0-L3 risk selects proportional planning, testing, review, and QA without a
  universal full/cold gate.
- Unit execution does not inherit credentials; PostgreSQL integration is
  ephemeral; Oracle/provider/browser/provider-write evidence remains explicit.
- The old cold command, cold-provision lane, cold-only tests, and cold acceptance
  criteria are absent from the active harness.
- Deterministic evals reject known architecture, context, evidence, and safety
  failures.
- A fresh task completes one dogfood flow and another session resumes it using
  MNFS, Git, context pack, and validation artifacts only.
- Architecture, governance, wiki, MNFS, code, and evidence agree at the accepted SHA.

## Handoff

- Current status: Replanned in progress; accepted foundations preserved, cold model superseded.
- Next owner: Milestone Orchestrator (request `gpt-5.6-terra`/high when the
  visible-task surface supports it; record the actual capability).
- Next action: Dispatch F-10 pragmatic harness cutover as the sole writer.
- Required files/evidence: F-10, F-05, and F-09 execution artifacts plus `validation-result.md`.
- Blockers or open decisions: None; do not reopen cold-clone diagnosis.

## Correction Handoff

- QA failure summary: Not applicable; this is an operator-approved scope replan.
- Correction scope: Not applicable until the next formal milestone gate.
- Attempts used/remaining: 0/2 after replan.
- Next artifact: F-10 `spec.md`, `plan.md`, and `validation.md`.
- Revalidation evidence required: M-08 required criteria excluding superseded C05/C12.
