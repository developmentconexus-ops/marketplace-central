---
name: mpc-goal-harness
description: Run Marketplace Central development through a visible Portfolio session, one Milestone session, bounded Feature plan/execution workers, compact context files, checkpoints, final review, and proportional QA.
---

# Marketplace Central Development Harness

The harness is a session protocol. Native Codex tasks carry conversation;
MNFS, Git, context files, and validation artifacts carry durable truth.

## Session topology

```text
Portfolio hub (visible, normally dormant)
  -> one standalone visible `mpc-milestone` task
       -> `mpc-implementer` subagent for a coherent vertical unit
       -> `mpc-verifier` subagent for final fixed-SHA review
       -> `mpc-verifier` subagent for proportional QA
  <- user resumes Portfolio from a compact terminal checkpoint
```

One active milestone outcome owns exactly one standalone visible Milestone
task. Start the project agent type `mpc-milestone` with clean context and pass
only the Portfolio packet. Human questions stay in that visible Milestone;
Portfolio does not relay them. The Milestone dispatches bounded children and
remains the sole outcome owner through terminal evidence.

The Milestone task is the root of its own native tree. Project config limits
that tree to depth 1: its implementation, review, and QA children never
delegate. Portfolio is a roadmap hub, not the parent runtime of a nested agent
tree, and stays dormant while the Milestone works.

## Dispatch model policy

Choose an explicit model and reasoning effort when creating every visible
Milestone or bounded Feature task; never inherit an expensive user default.

- Milestone orchestration: `gpt-5.6-terra` with `medium`.
- Coherent vertical implementation unit: `gpt-5.6-luna` with `high`.
- Simple read-only discovery: `gpt-5.6-luna` with `medium`.
- Fixed-SHA review: `gpt-5.6-luna` with `high`.
- Proportional deterministic QA: `gpt-5.6-luna` with `medium` when the caller
  can select it; the registered verifier remains Luna/high when native custom
  agent selection cannot vary effort per run.
- Escalate the Milestone to Terra/high only for a recorded cross-cutting,
  irreversible, security-sensitive, or repeatedly unresolved decision. Never
  escalate silently.

Pass the chosen values to the native visible-task creation call and include
them in the packet checkpoint. A task that needs a higher tier returns a
compact escalation request; it does not silently consume it.

- Portfolio owns product priority, dependencies, milestone start/stop, and the
  final mission view. It does not ingest feature logs.
- Milestone owns its checkout, feature order, integration, and communication
  with Portfolio. Only Milestone dispatches Feature, review, and QA subagents.
- An Implementer plans and executes one coherent vertical unit in the same
  session by default. Small plan slices and adjacent microfeatures stay inside
  that session. It never spawns another writer.
- Review and QA run after all Feature commits are integrated. An early review
  is allowed only for an irreversible or cross-cutting decision.

## Event-driven control plane

The parent channel is outcome-based, not a progress feed:

- The Milestone asks `needs_input` directly in its own visible task for a
  genuine human/authority decision and records `terminal` when no autonomous
  work remains. Native completion delivery does not reliably wake a different
  dormant task, so the protocol does not pretend that Portfolio receives an
  automatic callback.
- Feature acceptance, SHA freeze, review, and QA progress stay in the
  Milestone plan and durable evidence artifacts. They are not injected into
  Portfolio context.
- Feature, review, and QA tasks report only to Milestone. Portfolio reads no
  child transcript or raw log and sends no messages to child tasks.
- The user answers human questions in the same Milestone task. Portfolio steers
  only priority, pause, replacement, or explicit stop; it does not poll normal
  progress or create another Milestone for a follow-up question.
- No heartbeat, callback guard, cron, hook, polling loop, app server, or custom
  scheduler is part of the normal harness. When the Milestone is terminal, its
  durable checkpoint lets the user explicitly resume Portfolio without replay.

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
portfolio_task_id: <Portfolio task/thread ID for correlation only>
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

Start a standalone visible task with the named custom agent and no inherited
conversation turns. The
Milestone validates the packet, updates its visible plan, and proceeds without
a routine acknowledgment message. Before every child dispatch, run the
read-only preflight below. Portfolio may steer, resume, or interrupt the
Milestone; it never edits Milestone-owned paths concurrently.

## Milestone -> Implementer packet

```text
agent_type: mpc-implementer
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

### Implementer Plan

Read the brief and only the supplied knowledge. Write `spec.md` and `plan.md`.
The plan fixes owner, ports/interfaces, consumers, legacy decision, explicit
unknown states, paths, commands, and proof. Stop only when a cross-worker or
irreversible decision is unresolved.

### Implementer Execution

After `plan.md`, compile and validate the context file. Pass its path, not its
contents, when a fresh worker/session is used. Read only selectors, implement,
run impacted commands, write `validation.md`, create one intentional commit,
and return the compact handoff to Milestone.

## Context and token constraints

- Default child budget per Milestone is one Implementer run, one final review,
  and one proportional QA run. A correction run is allowed only after a named
  failed criterion and within the validation contract's retry cap. Any extra
  implementation split or repeated broad verification requires a compact human
  cost checkpoint before dispatch.
- Implementers run targeted proof. The integrated broad ladder runs once in
  proportional QA; do not repeat the same full suite in every small slice.
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

Every persisted `needs_input` or terminal handoff between Portfolio and
Milestone, and every Implementer return inside the Milestone, uses:

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
The legacy `heartbeat` event remains schema-valid only for historical artifact
compatibility. New dispatches must not create it.

Keep parent messages below 2,000 characters and reference evidence paths rather
than copying output. Resume from checkpoints, MNFS, Git, validation artifacts,
and context files. Native task IDs are correlation metadata. If native agent
controls are unavailable, stop and tell the user which visible-task capability
is missing; do not silently replace the Milestone with a nested Portfolio child.

## Boundary

The three project agents (`mpc-milestone`, `mpc-implementer`, and
`mpc-verifier`) plus the depth/thread limits in `.codex/config.toml` are the
only persistent orchestration runtime configuration. Tasks, subagents,
read/steer/interrupt, and worktrees remain operator-observed native
capabilities. Do not add hooks, app servers, VMs, custom schedulers, synthetic
eval products, cold clones, or a second CI without a separately reviewed
failure that requires them.
