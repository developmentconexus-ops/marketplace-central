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
       -> native bounded Feature subagents (plan + execution)
       -> native final fixed-SHA review subagent
       -> native proportional QA subagent
  <- compact checkpoints and final handoff
```

One active milestone outcome owns exactly one visible Milestone task. Portfolio
does not open another Milestone task for investigation, review, QA, or a
follow-up question inside that same outcome. The existing Milestone dispatches
its bounded Feature, review, and QA tasks and remains the sole callback owner
until it returns a terminal handoff or Portfolio explicitly replaces it.

Use native Milestone-owned subagents for Feature, review, and QA work. Create a
standalone child task only when native child controls are unavailable or when
the user needs independently owned follow-up work; label that packet
`fresh-session fallback`. Never make the fallback the normal topology.

## Dispatch model policy

Choose an explicit model and reasoning effort when creating every visible
Milestone or bounded Feature task; never inherit an expensive user default.

- Portfolio and implementation Milestone: `gpt-5.6-terra` with `high`.
- Feature implementation: `gpt-5.6-terra` with `high`.
- Read-only discovery, review, or QA: `gpt-5.6-luna` with `medium`.
- Use `gpt-5.6-sol` or `ultra` only when the Portfolio records a concrete
  cross-cutting or irreversible decision that the lower tier cannot resolve.

Pass the chosen values to the native visible-task creation call and include
them in the packet checkpoint. A task that needs a higher tier returns a
compact escalation request; it does not silently consume it.

- Portfolio owns product priority, dependencies, milestone start/stop, and the
  final mission view. It does not ingest feature logs.
- Milestone owns its checkout, feature order, integration, and communication
  with Portfolio. Only Milestone dispatches Feature, review, and QA subagents.
- A Feature worker plans and executes one feature in the same session by
  default. It never spawns another writer.
- Review and QA run after all Feature commits are integrated. An early review
  is allowed only for an irreversible or cross-cutting decision.

## Event-driven control plane

Callbacks are push-based, not user-driven polling:

- Milestone sends its checkpoint to `portfolio_task_id` using the native
  thread-message capability before it waits, ends a turn, or starts a new
  bounded child task. It never waits for the user to request a status update.
- Required callback moments: packet accepted; accepted Feature; external or
  owner blocker; fixed SHA frozen; review verdict; QA verdict; terminal
  handoff. Each callback uses the standard checkpoint fields below.
- Feature, review, and QA tasks report only to Milestone. Portfolio reads no
  child transcript or raw log and sends no messages to child tasks.
- Portfolio consumes pushed checkpoints and only steers the one Milestone for
  priority, a true blocker, or an explicit stop. It does not poll with
  `read_thread` during normal progress and does not create a second Milestone
  merely to ask a question within the active outcome.
- For work expected to outlast a normal turn, Portfolio may create one native
  Codex heartbeat on itself as a fallback. The heartbeat checks only whether a
  required callback is overdue, then requests one compact status from the
  Milestone. Never build a custom agent, hook, cron loop, or status service.

## Visible Milestone plan

Call native `update_plan` when accepting the Portfolio packet and after every
state transition: dispatch, Feature acceptance, blocker, SHA freeze, review,
QA, and handoff. The plan is a visible progress mirror for the task UI, not
durable orchestration truth; resume from packets, MNFS files, Git, and context.

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
updates the visible plan, and sends another after each accepted Feature, on a
blocker, and after QA. Before dispatching, run the read-only preflight below.
Portfolio may steer or interrupt Milestone; it never edits Milestone-owned
paths concurrently.

## Milestone -> Feature packet

```text
role: Feature Implementer
milestone_task_id: <native task/thread ID used for the return>
feature_id: <feature/work-item identity>
feature_file: <repository-relative feature.md>
context_files: [<repository-relative mission/milestone/guide and knowledge files>]
knowledge_routes: [<route IDs>]
base_sha: <accepted milestone SHA>
allowed_paths: [<exact repository-relative paths>]
forbidden_paths: [<exact repository-relative paths>]
shared_seams: [<exclusive repository-relative seams>]
side_effects: {allowed: [<effects>], forbidden: [<effects>]}
proof: {command_ids: [<registered IDs>], evidence_targets: [<paths>]}
stop_conditions: [<architecture, contract, ownership, runtime, or QA conflicts>]
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
schema_version:
milestone_id:
source_task_id:
parent_task_id:
event_id:
sequence:
event_type:
dispatch_id:
work_item_id:
base_sha:
commit_or_frozen_sha:
emitted_at:
callback_due_at:
status:
checkpoint_id:
commit:
changed_paths:
evidence:
review:
blockers:
next:
```

From the repository root, validate checkpoint JSON with:

```text
python .agents/skills/mpc-goal-harness/scripts/validate_checkpoint.py --checkpoint <checkpoint.json>
```

The schema is `.agents/skills/mpc-goal-harness/scripts/checkpoint_schema.json`;
a prior-state JSON input makes duplicate event IDs and non-monotonic sequence
invalid. Do not manufacture unknown SHA, commit, review, or operational facts:
use explicit `null` where the schema allows it.

From the repository root, run before a child dispatch:

```text
python .agents/skills/mpc-goal-harness/scripts/dispatch_preflight.py --packet <packet.json> --accepted-base-sha <Milestone-accepted-SHA> --current-writers <snapshot.json>
```

The SHA and `{"authoritative": true, "writers": []}`
snapshot are read-only Milestone inputs, not packet claims or a persistent
ledger. Preflight fails closed for a wrong/nonexistent accepted SHA, missing
packet field/file or feature/work-item identity, duplicate/completed/stale work
marker, absent authoritative one-writer state, writer-path overlap, or callback
that is not the parent task. It does not create, mutate, or resume a task. A
`heartbeat` checkpoint is valid only after `callback_due_at`; it is an overdue
fallback, never normal progress control.

Resume from the checkpoint, MNFS files, Git, and context file. Native task IDs
are correlation metadata. If task controls are unavailable, use a labelled
**fresh-session fallback** and pass the same packet; do not claim native parity.

## Boundary

Tasks, subagents, read/steer/interrupt, and worktrees are **operator-observed**
native capabilities. Do not build custom agents, hooks, app servers, VMs,
automation, synthetic eval products, cold clones, or a second CI for this flow.
