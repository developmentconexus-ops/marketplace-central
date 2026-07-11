---
name: mpc-goal-harness
description: Run Marketplace Central development through a visible Portfolio session, one Milestone session, bounded Feature plan/execution workers, compact context files, checkpoints, final review, and proportional QA.
---

# Marketplace Central Development Harness

The harness is a session protocol. Native Codex tasks carry conversation;
MNFS, Git, context files, and validation artifacts carry durable truth.

## Session topology

```text
Portfolio session
  -> one visible Milestone session
       -> bounded Feature workers (plan + execution)
       -> one final fixed-SHA review
       -> proportional QA
  <- compact checkpoints and final handoff
```

- Portfolio owns product priority, dependencies, milestone start/stop, and the
  final mission view. It does not ingest feature logs.
- Milestone owns its checkout, feature order, integration, and communication
  with Portfolio. Only Milestone dispatches Feature workers.
- A Feature worker plans and executes one feature in the same session by
  default. It never spawns another writer.
- Review and QA run after all Feature commits are integrated. An early review
  is allowed only for an irreversible or cross-cutting decision.

## Portfolio -> Milestone packet

Start a visible task with these fields:

```text
role: Milestone Orchestrator
portfolio_task_id: <native task/thread ID used for checkpoint replies>
objective: <one milestone outcome>
base_sha: <accepted 40-char SHA>
mission_file: <mission.md>
milestone_file: <milestone.md>
execution_guide: <execution-guide.md>
knowledge_routes: <route IDs>
constraints: <paths, seams, side effects, stop conditions>
qa_contract: <validation-contract.md>
next: reconcile features and send the first checkpoint
```

Milestone acknowledges the packet, sends a checkpoint to `portfolio_task_id`,
and sends another after each accepted Feature, on a blocker, and after QA.
Portfolio may steer or interrupt Milestone; it never edits Milestone-owned
paths concurrently.

## Milestone -> Feature packet

```text
role: Feature Implementer
milestone_task_id: <native task/thread ID used for the return>
feature_file: <feature.md>
context_files: <mission/milestone/guide plus required knowledge paths>
knowledge_routes: <route IDs>
base_sha: <accepted milestone SHA>
allowed_paths: <exact paths>
forbidden_paths: <exact paths>
shared_seams: <exclusive seams>
side_effects: <allowed and forbidden>
proof: <registered command IDs and evidence targets>
stop_conditions: <architecture, contract, ownership, runtime, or QA conflicts>
```

### Feature Plan

Read the brief and only the supplied knowledge. Write `spec.md` and `plan.md`.
The plan fixes owner, ports/interfaces, consumers, legacy decision, explicit
unknown states, paths, commands, and proof. Stop only when a cross-worker or
irreversible decision is unresolved.

### Feature Execution

After `plan.md`, compile and validate the context file. Pass its path, not its
contents, when a fresh worker/session is used. Read only selectors, implement,
run impacted commands, write `validation.md`, create one intentional commit,
and return the compact handoff to Milestone.

## Context and token constraints

- Prompts pass file paths and selectors, not copied document bodies.
- Initial reads are bootstrap + packet files + selected knowledge only.
- No transcript replay, repository-wide scan, raw logs, unrelated milestone
  history, repeated reads, or full implementation tree without a named gap.
- A route gap permits one targeted search. Stable discoveries update the
  canonical knowledge route in the accepted slice.
- The 2,000-token estimate is a dispatch budget: L0/L1 must fit; necessary
  L2/L3 overflow names the reason per source.
- Checkpoints contain decisions and evidence paths, never logs.

## Ownership and validation

One writer owns a checkout/shared seam. Worktrees coordinate disjoint writers;
they are not VMs. Use registered command IDs only. Fake evidence proves only
deterministic behavior. Oracle, provider, database, browser, and provider-write
targets remain distinct and run only when the milestone contract requires them.

Milestone integrates Feature commits, freezes one SHA, then requests one review
and proportional QA. Only QA writes/passes `validation-result.md`.

## Checkpoint and handoff

Every message between Portfolio and Milestone, and every Feature return, uses:

```text
status:
checkpoint_id:
commit:
changed_paths:
evidence:
review:
blockers:
next:
```

Resume from the checkpoint, MNFS files, Git, and context file. Native task IDs
are correlation metadata. If task controls are unavailable, use a labelled
**fresh-session fallback** and pass the same packet; do not claim native parity.

## Boundary

Tasks, subagents, read/steer/interrupt, and worktrees are **operator-observed**
native capabilities. Do not build custom agents, hooks, app servers, VMs,
automation, synthetic eval products, cold clones, or a second CI for this flow.
