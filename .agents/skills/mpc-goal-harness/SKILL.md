---
name: mpc-goal-harness
description: SUPERSEDED 2026-07-15 — do NOT invoke. The binding harness is docs/superpowers/HARNESS.md (hub-and-chips); workers use the harness-worker skill. This file is retained only for the dispatch_preflight.py concept pending migration.
---

> **SUPERSEDED (2026-07-15).** This /goal + Portfolio-Hub protocol was replaced by the
> hub-and-chips harness: `docs/superpowers/HARNESS.md` (BINDING) + `.agents/skills/harness-worker`.
> Do not follow anything below. Retained temporarily because
> `scripts/dispatch_preflight.py` contains a fail-closed dispatch-packet validation concept
> queued for migration into the new harness (hub queue item).

# Marketplace Central Goal Harness (RETIRED)

This harness is a lightweight session protocol. `/goal` keeps a Milestone
running; this skill defines ownership, delegation, validation, and handoff.
MNFS, Git, and named evidence remain durable truth.

## Topology

~~~text
Portfolio Hub (visible, normally dormant)
  -> emits one copyable /goal handoff and stops
  -> user starts one clean visible Milestone task
       -> mpc-implementer: Luna/high, one coherent writer at a time
       -> mpc-verifier: Luna/high, fixed-SHA review then proportional QA
  <- explicit user-requested escalation or terminal handoff
~~~

The Hub owns roadmap priority, dependencies, and milestone selection. The
Milestone owns one milestone outcome, feature order, integration, review, QA,
and communication. The Hub does not create or poll the Milestone and does not
edit Milestone-owned paths while it runs.

The user chooses the Milestone root model and reasoning in the app. The worker
profiles are project-scoped custom agents:

- `.codex/agents/mpc-implementer.toml`: `gpt-5.6-luna`, `high`;
- `.codex/agents/mpc-verifier.toml`: `gpt-5.6-luna`, `high`.

Codex identifies a custom agent by its `name`. The Milestone requests
`mpc-implementer` or `mpc-verifier` directly by name and supplies the bounded
packet for that run. Their TOML profiles fix Luna/high; no separate capability
probe or custom-agent preflight belongs in milestone execution.

Project config keeps `max_depth = 1`; children never delegate. Keep at most one
writer active in a checkout or shared seam.

## Hub to Milestone

The Hub resolves its exact task ID, prepares one handoff under 4,000
characters, shows it to the user, and stops. Use this shape:

~~~text
/goal Execute <milestone_id> as the Marketplace Central Milestone Orchestrator.
Read AGENTS.md, use $mpc-goal-harness, and follow only the named files and
knowledge routes below. Continue until a blocking human decision or terminal.

hub_task_id: <exact Hub task/thread ID>
milestone_id: <ID>
objective: <one outcome and definition of done>
accepted_base_sha: <40-char SHA>
mission_file: <repository-relative path>
milestone_file: <repository-relative path>
execution_guide: <repository-relative path>
knowledge_routes: [<route IDs>]
constraints: <allowed paths, shared seams, forbidden effects, stop conditions>
qa_contract: <repository-relative path>
~~~

Long details belong in the named repository files, not in the goal. The user
creates the task manually, selects the desired Milestone root model/effort,
and pastes the handoff. The Hub never calls `create_thread` for this flow.

## Milestone loop

1. Read `AGENTS.md`, this skill, the handoff files, and only named knowledge
   selectors. Reconcile HEAD, accepted base SHA, dirty state, completed work,
   and writer ownership once.
2. If a product, authority, contract, or irreversible choice is genuinely
   ambiguous, ask the user directly in the visible Milestone and pause. Do not
   route ordinary questions through the Hub.
   Formatting-only normalization or bounded parser compatibility that preserves
   the same contract meaning is an autonomous remediation, not a user-input
   stop; record the deterministic proof with the implementing change.
3. Maintain a visible plan at major gates only: accepted, coherent unit done,
   SHA frozen, review, QA, terminal.
4. Dispatch `mpc-implementer` for one coherent vertical unit. Do not create an
   agent per function or plan step. The Milestone root does not edit while the
   child writer is active.
5. Accept only intentional commits with targeted proof and named evidence.
   Repeat implementation only for another coherent unit or a named failed
   criterion inside the milestone retry cap.
6. Freeze the integrated SHA. Dispatch `mpc-verifier` in `fixed_sha_review`
   mode with no allowed write paths.
7. After review passes, dispatch `mpc-verifier` again in `proportional_qa`
   mode. QA runs only registered commands and may write only named QA evidence
   and `validation-result.md`. Only QA may pass the milestone.
8. Persist the compact terminal checkpoint, send the terminal handoff to the
   Hub, and end the goal. Do not restart workers or callbacks after terminal.

No discovery child exists by default. Use one only when a bounded read-only
investigation materially reduces context or latency. Parallel writers require
separate worktrees and disjoint seams; otherwise work sequentially.

## Milestone to Implementer

Dispatch the named `mpc-implementer` custom agent with no inherited
conversation turns and a compact packet:

~~~text
agent_type: mpc-implementer
role: Feature Implementer
milestone_task_id: <Milestone task/thread ID>
work_item: <coherent unit ID and feature file>
base_sha: <accepted SHA>
context_files: [<mission, milestone, guide, selected knowledge files>]
knowledge_routes: [<route IDs>]
allowed_paths: [<exact paths>]
forbidden_paths: [<exact paths>]
shared_seams: [<exclusive seams>]
side_effects: {allowed: [<effects>], forbidden: [<effects>]}
proof: {command_ids: [<registered IDs>], evidence_targets: [<paths>]}
stop_conditions: [<architecture, contract, ownership, runtime, QA conflicts>]
~~~

The Implementer plans and executes the coherent unit, stays within the packet,
runs targeted proof, writes `validation.md`, creates one intentional commit,
and returns changed paths, commands, evidence, blockers, and next. It never
delegates or contacts the Hub.

## Milestone to Verifier

Dispatch the named `mpc-verifier` custom agent after integration:

~~~text
agent_type: mpc-verifier
role: Milestone Verifier
mode: fixed_sha_review | proportional_qa
frozen_sha: <40-char integrated SHA>
contract_files: [<milestone and validation contract>]
evidence_files: [<named evidence>]
allowed_write_paths: [] | [<QA evidence paths>]
registered_commands: [<IDs allowed in this mode>]
stop_conditions: [<verification conflicts>]
~~~

Review is read-only and reports Pass/Fail with actionable findings. QA runs
only the registered proportional ladder and writes only explicitly allowed
evidence. The verifier never delegates or contacts the Hub.

## Communication

There are exactly three communication events:

- `needs_input_local`: ask the user directly in the Milestone and pause. No Hub
  message.
- `escalation_requested`: only after the user explicitly asks to escalate.
  Persist a compact checkpoint and send the Hub the decision needed, blockers,
  evidence paths, and next action. This is not a progress update or polling
  channel.
- `terminal_handoff`: always after passed, failed, or externally blocked.
  Persist the compact checkpoint first, then send its path and verdict to the
  exact `hub_task_id` with `send_message_to_thread`.

Terminal payload:

~~~text
event_type: terminal_handoff
milestone_id: <ID>
source_task_id: <Milestone task/thread ID>
checkpoint: <repository-relative path>
status: <passed | failed | externally_blocked>
sha: <frozen SHA or null>
evidence: [<paths only>]
blockers: [<compact items>]
next: <exact Hub action>
~~~

If cross-task messaging fails, retain the checkpoint and return the same
payload for the user to paste into the Hub. A final response alone is not a
callback. No heartbeat, polling loop, hook, scheduler, or app server belongs
in the normal harness.

## Context and safety

- Pass paths and selectors, not copied document bodies or transcripts.
- Initial reads are bootstrap, handoff files, and named knowledge only.
- A route gap permits one targeted search; stable findings update the
  canonical route.
- Never reset, revert, stash, clean, cold-clone, expose secrets/PII, or invent
  evidence.
- Preserve domain/application/ports/adapters/transport boundaries, tenant
  scope, explicit unknown states, OpenAPI plus SDK lockstep, and provider-write
  safeguards from `AGENTS.md`.
- Mocks prove deterministic contracts, never live integration.

## Legacy artifacts

The scripts and schema under `scripts/` remain only for validating historical
version-1 checkpoints and dispatch packets. New Milestones do not create
heartbeats, sequence ledgers, writer snapshots, or per-transition checkpoint
events.
