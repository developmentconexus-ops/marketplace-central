---
name: mpc-goal-harness
description: Run Marketplace Central development through a visible Portfolio session, one Milestone session, bounded Feature plan/execution workers, compact context files, checkpoints, final review, and proportional QA.
---

# Marketplace Central Development Harness

The harness is a session protocol. Native Codex tasks carry conversation;
MNFS, Git, context files, and validation artifacts carry durable truth.

## Session topology

```text
Portfolio hub
  -> one clean native `mpc-milestone` subagent
       -> native bounded Feature subagents (plan + execution)
       -> native final fixed-SHA review subagent
       -> native proportional QA subagent
  <- `needs_input` or compact terminal handoff only
```

One active milestone outcome owns exactly one visible native Milestone
subagent. Spawn the project agent type `mpc-milestone` with clean context and
pass only the Portfolio packet. Portfolio does not open another Milestone for
investigation, review, QA, or a follow-up question inside that outcome. The
Milestone dispatches its bounded Feature, review, and QA children and remains
the sole outcome owner until its automatic terminal return or explicit
replacement.

Project config limits the native tree to Portfolio depth 0, Milestone depth 1,
and Feature/review/QA depth 2. Milestone children never delegate. Create a
standalone task only when native controls are unavailable or the user requires
independent ownership; label it `fresh-session fallback`. Never make that
fallback the normal topology.

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

The parent channel is outcome-based, not a progress feed:

- The Milestone reports only `needs_input` for a genuine human/authority
  decision and `terminal` when no autonomous work remains. Native completion
  delivery alone does not guarantee that a dormant Portfolio task wakes, so
  the Portfolio must not rely on delivery as its wake-up mechanism.
- Feature acceptance, SHA freeze, review, and QA progress stay in the
  Milestone plan and durable evidence artifacts. They are not injected into
  Portfolio context.
- Feature, review, and QA tasks report only to Milestone. Portfolio reads no
  child transcript or raw log and sends no messages to child tasks.
- Portfolio resumes the same Milestone after a human answer and otherwise
  steers it only for priority or an explicit stop. It does not poll during
  normal progress or create a second Milestone to ask a question.
- Whenever a dispatched Milestone may outlast the current Portfolio turn,
  Portfolio must create one native Codex heartbeat on itself before ending
  that turn. The heartbeat inspects only the single active Milestone's native
  summary/status; it never reads child logs, transcripts, Feature tasks, review
  tasks, or QA tasks. If that Milestone is completed or its required callback
  is overdue, the heartbeat wakes Portfolio and requests exactly one compact
  `needs_input` or `terminal` result from the same Milestone. Otherwise it does
  nothing and waits for its next native invocation.
- Portfolio deletes the heartbeat immediately after consuming the Milestone's
  `needs_input` or `terminal` result, or when that Milestone is stopped or
  replaced. At most one heartbeat exists for the single active Milestone; a
  replacement dispatch creates its own heartbeat only if it may outlast the
  current turn. Never build a custom agent, hook, cron loop, app server,
  scheduler, or status service.

## Visible Milestone plan

The Milestone calls native `update_plan` when accepting the Portfolio packet
and after every state transition. Its plan is the visible progress mirror for
internal Feature, review, and QA work, not durable orchestration truth. The
Portfolio plan tracks milestone outcomes only. Resume from packets, MNFS, Git,
validation artifacts, and context files.

## Portfolio -> Milestone packet

Start a visible task with these fields:

```text
agent_type: mpc-milestone
role: Milestone Orchestrator
parent_task_id: <Portfolio task/thread ID>
objective: <one milestone outcome>
base_sha: <accepted 40-char SHA>
mission_file: <mission.md>
milestone_file: <milestone.md>
execution_guide: <execution-guide.md>
knowledge_routes: <route IDs>
constraints: <paths, seams, side effects, stop conditions>
qa_contract: <validation-contract.md>
next: reconcile repository truth and proceed until needs_input or terminal
```

Spawn the named custom agent with no inherited conversation turns. The
Milestone validates the packet, updates its visible plan, and proceeds without
a routine acknowledgment message. Before every child dispatch, run the
read-only preflight below. Portfolio may steer, resume, or interrupt the
Milestone; it never edits Milestone-owned paths concurrently.

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

Every `needs_input` or terminal message between Portfolio and Milestone, and
every Feature return inside the Milestone, uses:

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

Keep parent messages below 2,000 characters and reference evidence paths rather
than copying output. Resume from checkpoints, MNFS, Git, validation artifacts,
and context files. Native task IDs are correlation metadata. If native agent
controls are unavailable, use a labelled **fresh-session fallback** and pass
the same packet; do not claim native parity.

## Boundary

The single project custom agent `.codex/agents/mpc-milestone.toml` and depth
limit in `.codex/config.toml` are the only persistent orchestration runtime
configuration. Tasks, subagents, read/steer/interrupt, and worktrees remain
operator-observed native capabilities. Do not add more custom agents, hooks,
app servers, VMs, custom schedulers, synthetic eval products, cold clones, or a
second CI without a separately reviewed failure that requires them.
