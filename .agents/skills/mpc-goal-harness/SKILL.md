---
name: mpc-goal-harness
description: Run Marketplace Central through a visible Portfolio hub, one manually started Milestone session, bounded generic workers, compact context, terminal callbacks, final review, and proportional QA.
---

# Marketplace Central Development Harness

The harness is a session protocol. Native Codex tasks carry conversation;
MNFS, Git, context files, and validation artifacts carry durable truth.

## Validated runtime boundary

The current Codex app can send a message from one visible task to another and
wake the destination. A manual-session probe on 2026-07-13 also established
that native child dispatch exposes task_name, message, and context forking, but
no agent_type, model, or reasoning selector.

Therefore:

- never claim that a child is mpc-implementer or mpc-verifier;
- never use task_name as evidence of a custom agent type;
- never claim a child model or effort that the runtime did not report;
- dispatch only generic direct children with complete role packets;
- keep child count bounded because per-child cost cannot be pinned;
- stop and revalidate this boundary before restoring custom-agent config.

## Session topology

~~~text
Portfolio hub (visible, normally dormant)
  -> prepares one copyable Milestone prompt and stops
  -> user manually starts one clean visible Milestone task at Terra/medium
       -> generic direct Feature Implementer child
       -> generic direct fixed-SHA Reviewer child
       -> generic direct proportional QA child
  <- Milestone persists terminal checkpoint, then calls send_message_to_thread
~~~

One active milestone outcome owns exactly one standalone visible Milestone
task. Portfolio never creates that task. It prepares the exact prompt, shows it
to the user, and stops. The user creates a clean project task manually, selects
gpt-5.6-terra with medium, and pastes the prompt. The prompt points to this
skill and assigns the Milestone Orchestrator role; the root is generic.

The Milestone is the root of its native child tree. Project configuration
limits that tree to depth 1 and three concurrent threads. Children never
delegate. Human questions stay in the visible Milestone. Portfolio stays
dormant while it works and never edits Milestone-owned paths.

## Dispatch and cost policy

- Visible Milestone: user selects gpt-5.6-terra with medium.
- Direct children: runtime-managed model and reasoning; record
  runtime_managed, not Luna, Terra, or an inferred value.
- Default budget: one coherent Feature Implementer run, one final fixed-SHA
  review, and one proportional QA run.
- Discovery stays inside the Implementer unless a separate read-only child
  materially reduces context or latency.
- A correction run is allowed only for a named failed criterion and inside the
  validation contract retry cap.
- Any extra implementation split, discovery child, or repeated broad
  verification requires a compact human cost decision in the Milestone.
- If exact child model, effort, or custom-agent identity becomes required,
  stop with a runtime capability conflict.

Portfolio owns priority, dependencies, milestone start/stop, and the final
mission view. Milestone owns feature order, child dispatch, integration,
review, QA, and Portfolio communication. An Implementer plans and executes one
coherent vertical unit; small plan slices remain inside that run.

## Event-driven control plane

- needs_input: ask the human directly in the visible Milestone. Do not relay
  the question through Portfolio.
- terminal: persist and validate one compact checkpoint, then explicitly call
  send_message_to_thread with portfolio_task_id.
- Native completion is not a callback. A final Milestone response without the
  cross-task message does not wake Portfolio.
- If cross-task messaging is unavailable or fails, retain the checkpoint and
  return a pasteable callback so the user can resume Portfolio manually.
- Portfolio validates the referenced checkpoint before acting and never
  reconstructs truth from the Milestone transcript.
- Feature, freeze, review, and QA progress stay in the Milestone plan and
  durable evidence. Only the terminal outcome reaches Portfolio.
- No heartbeat, callback guard, cron, hook, polling loop, app server, or custom
  scheduler is part of the normal harness.

## Visible plans

Milestone calls native update_plan when accepting the packet and after each
state transition. Its plan mirrors Feature, review, and QA work; it is not
durable truth. Portfolio plan tracks milestone outcomes only. Resume from
packets, MNFS, Git, validation artifacts, context files, and checkpoints.

## Portfolio to manually created Milestone

Portfolio resolves its exact native task ID through app task controls. It must
not guess. It emits one copyable prompt with:

~~~text
session_type: manually_created_standalone_visible_root
role: Milestone Orchestrator
model_to_select: gpt-5.6-terra
reasoning_effort_to_select: medium
harness_skill: .agents/skills/mpc-goal-harness/SKILL.md
portfolio_task_id: <exact Portfolio task/thread ID>
objective: <one milestone outcome>
base_sha: <accepted 40-char SHA>
mission_file: <mission.md>
milestone_file: <milestone.md>
execution_guide: <execution-guide.md>
knowledge_routes: <route IDs>
constraints: <paths, seams, side effects, stop conditions>
qa_contract: <validation-contract.md>
next: reconcile repository truth and proceed until needs_input or terminal
~~~

The prompt also says:

1. user, not Portfolio, creates the new task;
2. user selects Terra/medium before sending it;
3. generic root reads AGENTS.md, this skill, packet files, and only named
   knowledge selectors;
4. terminal requires checkpoint validation followed by
   send_message_to_thread(portfolio_task_id, compact_payload);
5. Milestone must not use create_thread or emulate a custom root.

After emitting the prompt, Portfolio stops. It does not poll or create the
Milestone.

## Milestone to generic Feature Implementer

Run read-only dispatch preflight, then spawn one direct generic child with no
inherited conversation turns. task_name is a correlation label, not an agent
type. The message contains:

~~~text
runtime_agent_type: generic
role: Feature Implementer
runtime_model_policy: runtime_managed
milestone_task_id: <native Milestone task/thread ID>
feature_id: <feature/work-item identity>
feature_file: <repository-relative feature.md>
context_files: [<mission/milestone/guide and selected knowledge files>]
knowledge_routes: [<route IDs>]
base_sha: <accepted milestone SHA>
allowed_paths: [<exact repository-relative paths>]
forbidden_paths: [<exact repository-relative paths>]
shared_seams: [<exclusive repository-relative seams>]
side_effects: {allowed: [<effects>], forbidden: [<effects>]}
proof: {command_ids: [<registered IDs>], evidence_targets: [<paths>]}
stop_conditions: [<architecture, contract, ownership, runtime, QA conflicts>]
~~~

The child reads only supplied files and selectors; writes spec.md and plan.md
when no accepted contract exists; compiles and validates context after
planning; implements inside allowed paths; runs targeted proof; writes
validation.md; creates one intentional commit; and returns a compact handoff.
It never delegates or sends messages to Portfolio.

## Milestone to generic Verifier

Review and QA are separate generic direct-child runs after accepted Feature
commits integrate:

~~~text
runtime_agent_type: generic
role: Milestone Verifier
mode: fixed_sha_review | proportional_qa
runtime_model_policy: runtime_managed
milestone_task_id: <native Milestone task/thread ID>
frozen_sha: <40-char integrated SHA>
contract_files: [<milestone and validation contract paths>]
evidence_files: [<named validation evidence paths>]
allowed_write_paths: [] | [<QA validation/evidence paths>]
registered_commands: [<command IDs allowed for this mode>]
stop_conditions: [<verification conflicts>]
~~~

In fixed_sha_review the child is read-only and returns Pass/Fail with
actionable findings. In proportional_qa it runs only registered commands and
writes only explicitly allowed validation/evidence paths. Only QA may pass the
milestone. Neither verifier delegates.

## Context and ownership constraints

- Prompts pass paths and selectors, never copied document bodies.
- Initial reads are bootstrap, packet files, and selected knowledge only.
- No transcript replay, repository-wide scan, raw logs, unrelated milestone
  history, or full tree without a named route gap.
- A route gap permits one targeted search. Stable discoveries update the
  canonical knowledge route in the accepted slice.
- The 2,000-token estimate is a dispatch budget: L0/L1 must fit; necessary
  L2/L3 overflow names the reason.
- Checkpoints contain decisions and evidence paths, never logs.
- One writer owns a checkout and shared seam. Never edit concurrently with a
  child writer.
- Worktrees coordinate truly disjoint writers; they are not VMs.
- Mocks prove deterministic contract behavior, never live integration.

Milestone integrates Feature commits, freezes one SHA, then requests one review
and one proportional QA. Only QA writes or passes validation-result.md.

## Checkpoint and terminal callback

Every persisted needs_input or terminal handoff, and every Implementer return,
uses this schema. For terminal, persist and validate it before messaging.

~~~text
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
~~~

Validate from repository root:

~~~powershell
python .agents/skills/mpc-goal-harness/scripts/validate_checkpoint.py --checkpoint <checkpoint.json>
~~~

For terminal only, send a payload under 2,000 characters to portfolio_task_id
using native send_message_to_thread. Do not override Hub model or effort:

~~~text
event_type: terminal_handoff
milestone_id: <id>
source_task_id: <Milestone task/thread ID>
checkpoint: <repository-relative checkpoint path>
status: <passed | failed | externally_blocked>
commit_or_frozen_sha: <SHA or null>
evidence: [<paths only>]
blockers: [<compact items>]
next: <exact Portfolio action>
~~~

Receiving Portfolio validates the checkpoint, reconciles repository truth,
updates its plan, then closes the outcome or prepares the next manual
Milestone prompt. It does not read child transcripts.

The checkpoint schema is
.agents/skills/mpc-goal-harness/scripts/checkpoint_schema.json. Unknown facts
remain explicit null where allowed. Legacy heartbeat remains schema-valid only
for historical compatibility; new dispatches must not create it.

## Child dispatch preflight

Before a Feature child:

~~~powershell
python .agents/skills/mpc-goal-harness/scripts/dispatch_preflight.py --packet <packet.json> --accepted-base-sha <Milestone-accepted-SHA> --current-writers <snapshot.json>
~~~

The SHA and {"authoritative": true, "writers": []} snapshot are read-only
Milestone inputs. Preflight fails closed for wrong/nonexistent SHA, missing
packet field/file or identity, stale/duplicate markers, absent authoritative
one-writer state, writer overlap, or callback target not equal to parent. It
does not create, mutate, or resume tasks.

## Boundary

Root and child roles are defined by this skill plus packets.
.codex/config.toml contains only multi-agent enablement and bounded
depth/thread limits; there are no project custom-agent registrations. Tasks,
generic subagents, cross-task messaging, read/steer/interrupt, and worktrees
remain operator-observed native capabilities. Do not add hooks, app servers,
VMs, custom schedulers, synthetic eval products, cold clones, or a second CI
without a separately reviewed failure that requires them.
