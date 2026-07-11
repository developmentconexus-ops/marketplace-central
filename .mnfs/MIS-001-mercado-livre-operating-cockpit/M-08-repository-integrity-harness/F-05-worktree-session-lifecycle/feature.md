# F-05-goal-orchestration-control-plane

```yaml
id: F-05
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-08
created: 2026-07-10
updated: 2026-07-11
validation_level: QA-3
lifecycle_scope: feature
```

## Mission

MIS-001 Mercado Livre Operating Cockpit.

## Milestone

M-08 Repository Integrity Harness.

## Brief

Connect a Codex goal to canonical repository knowledge, one visible milestone
task, bounded subagents, leases, risk-selected gates, compact handoffs, and
artifact-only resume. Use native Codex surfaces instead of implementing a custom
agent runtime.

## Inputs

- Accepted F-06/F-07/F-08 knowledge, context, environment, and execution seams.
- F-10 impact-gate command and evidence shape.
- Official Codex task, built-in subagent, skill, worktree/handoff, browser, and
  `codex exec` capabilities recorded in R-08.
- Current module, invariant, runtime, lane, and shared-seam registries.

## Expected Output

- Machine-validated knowledge routes from module/change concern to exact
  architecture sections, interface files/symbols, consumers, tests, and QA commands.
- Root `AGENTS.md` and the active execution guide use a short bootstrap and
  route detailed workflow/knowledge through the repo skill and current pack;
  no durable policy is lost or duplicated.
- Context packs expose ordered selectors; the measured bootstrap-plus-selector
  set targets at most 2,000 initial tokens and requires a source-by-source
  `overflow_reason` for L2/L3 overflow.
- `State.psm1` owns ignored run state, optional task correlation, leases, checkpoints,
  resume, recovery, and compact handoff validation.
- `.agents/skills/mpc-goal-harness/` performs goal classification, MNFS
  reconciliation, context compilation, milestone dispatch/steering, gate
  selection, acceptance audit, and resume through progressive references.
- Built-in Codex subagents receive bounded role prompts and depth remains one.
- Native worktree dispatch is optional and only selected for safe parallel
  writers; normal local caches/dependencies are reused.

## Inputs/Outputs

- Input: active goal, mission/milestone/feature path, accepted base SHA,
  requested paths or module, optional real QA target.
- Dispatch packet: objective, observable done, exact read selectors,
  allowed/forbidden paths, seams, risk, gates, commands, stop conditions, output schema.
- Checkpoint: checkpoint ID, optional native-task correlation ID, feature state,
  base/current commit, changed paths, completed evidence, blocker, and next action.
- Handoff: accepted commit, changed paths, criterion-to-evidence mapping,
  remaining risks, and next owner.

## Constraints

- Repository artifacts are authoritative; task transcript and native goal are control surfaces only.
- No custom app-server client or dependency on experimental APIs in V1.
- Hooks may warn or preflight but may not be the only enforcement.
- No architecture prose duplicated into the repo skill or agent files.
- Custom agents, nested `AGENTS.md`, hooks, and scheduled automation are
  deferred unless F-09 dogfood proves a concrete repeated gap.
- No task per micro-step, recursive fan-out, or parallel writers on one seam.
- No cold clone, cache reset, repeated dependency installation, or VM-like isolation.
- Worktree ignored-file copying is explicit; secrets never enter tracked files or prompts.

## Negative Scenarios

- Stale/missing context source: dispatch fails with a context reason code.
- Out-of-scope path or competing seam: writer is not started or accepted.
- Missing owner/interface route for L2/L3 work: plan blocks for targeted investigation.
- A mandatory full-document read bypasses the pack or its context measurement:
  context validation fails.
- Task disappears or session compacts: resume uses checkpoint and repository artifacts.
- Unsupported task-control surface: orchestration stops or uses a labelled
  fresh-session fallback; it never claims app-visible steering parity.
- A subagent attempts recursive fan-out or returns raw logs: orchestration rejects the handoff.
- Fake evidence is offered for a real QA criterion: feature remains unaccepted.

## Validation Expectations

- A fresh task receives a schema-valid pack whose total harness-requested
  bootstrap-plus-selector set meets the 2,000-token target or carries justified
  L2/L3 overflow and names its first action from declared routes.
- Two simulated writers requesting one seam yield one lease and one stable rejection.
- After an interrupted worker, a fresh session reconstructs the exact next
  action from MNFS, Git, context pack, and checkpoint without portfolio transcript.
- Operator-observed dogfood shows a milestone task dispatching at least one
  bounded subagent and reading, steering, and consuming its compact result
  through supported current-host tools; the repository persists only checkpoint
  state and optional correlation IDs.
- L1 and L3 fixtures select different review/QA policies with exact expected commands.

## Criterion Mapping

| Criterion | Ownership | Minimum proof |
| --- | --- | --- |
| M-08-C06 | Primary | Native worktree writer leaves primary checkout unchanged. |
| M-08-C07 | Primary | Fresh dispatch packet contains every required field and exact read selectors. |
| M-08-C09 | Primary | Knowledge routes compile current packs and stale/unrelated routes fail. |
| M-08-C10 | Primary | Conflicting seam and out-of-scope writer fixtures reject before acceptance. |
| M-08-C11 | Primary | Repo skill uses supported native goal/task/subagent controls or labels a fallback. |
| M-08-C13 | Primary | Goal fixtures select owner, interfaces, consumers, legacy decision, risk, and proof. |
| M-08-C14 | Primary | Checkpoint/resume reproduces the exact next action without transcript. |
| M-08-C15 | Primary | L0-L3 fixtures select proportional deterministic, review, and real-QA gates. |

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution.

## Handoff

- Current status: Briefed after pragmatic replan.
- Next owner: Feature Implementer after F-10 acceptance.
- Next action: Create spec/plan from the native control-plane contract.
- Required files/evidence: capability matrix, knowledge routes, state/lease tests,
  repo-skill validation, operator-observed native-task proof, and resume proof.
- Blockers or open decisions: None; app-server remains future optional scope.
