# M-08-repository-integrity-harness

```yaml
id: M-08
type: milestone
status: passed
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
| F-10 | Pragmatic harness cutover | Accepted current-checkout, task-declared impact/evidence gate with cold execution removed. |
| F-05 | Simple session orchestration | Accepted knowledge/context/state foundations; final protocol defines Portfolio ↔ Milestone checkpoints and Milestone → Feature plan/execution. |
| F-09 | Synthetic eval WIP | Rejected from V1. Preserved as WIP evidence only; it is not active or required for completion. |

## Dependencies

- F-01, F-03, F-06, F-07, and F-08 are accepted foundations.
- F-02 and F-04 are historical superseded experiments and do not gate V1.
- F-10 precedes F-05 so no active workflow depends on cold-only vocabulary.
- F-09 is outside the active dependency graph.
- Product work may resume after final review/QA of the simple F-10/F-05 flow.

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
- Portfolio and Milestone exchange the defined packets and compact checkpoints;
  Milestone dispatches Feature Plan + Execution and receives one commit/evidence
  handoff without transcript replay.
- A fresh Portfolio task can resume product development from mission, Git,
  knowledge routes, and the milestone handoff.
- Architecture, governance, wiki, MNFS, code, and evidence agree at the accepted SHA.

## Handoff

- Current status: Passed on fixed commit `0adae4d8203718a3a6a0058314b2a3d61b363bea`; F-10/F-05 active and F-09 rejected from V1.
- Next owner: fresh Portfolio session.
- Next action: Resume M-06 at its paid resolved-link Oracle + Mercado Livre evidence gap.
- Required files/evidence: `validation-result.md` and the simple session packets in the repo skill.
- Blockers or open decisions: None; harness work is closed.

## Correction Handoff

- QA failure summary: Not applicable; this is an operator-approved scope replan.
- Correction scope: Not applicable until the next formal milestone gate.
- Attempts used/remaining: 0/2 after replan.
- Next artifact: F-10 `spec.md`, `plan.md`, and `validation.md`.
- Revalidation evidence required: M-08 required criteria excluding superseded C05/C12.
