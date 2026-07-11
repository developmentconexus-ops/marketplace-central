# Repository-Native Development Harness — Pragmatic V1

```yaml
status: approved_final_simple_v1
owner: Portfolio Orchestrator
approved_by: operator
approved_on: 2026-07-11
supersedes: cold-clone and clean-machine requirements from the 2026-07-10 design
research: .mnfs/MIS-001-mercado-livre-operating-cockpit/research/R-08-codex-development-harness.md
```

## Final Authority — Simple Session Protocol

This section supersedes any later text that implies F-09, a synthetic eval
corpus, a long state machine, per-feature final review, or M-08 as a permanent
product blocker.

The active topology is:

```text
Portfolio -> visible Milestone -> bounded Feature Plan + Execution workers
Portfolio <- milestone checkpoints <- feature commit/evidence handoffs
Milestone complete -> one fixed-SHA review -> proportional QA
```

Portfolio and Milestone communicate through native Codex task messages using
the packet and checkpoint schemas in `.agents/skills/mpc-goal-harness/SKILL.md`.
Each dispatch passes file paths, knowledge route IDs, constraints, base SHA,
proof targets, and callback task ID. It never copies broad documents or relies
on transcript history. The context file is compiled after Feature Plan and is
the read authority for Feature Execution.

F-10 provides the optional current-checkout impact gate. F-05 provides context,
knowledge routes, checkpoints, handoff validation, and optional seam/worktree
coordination. F-09 remains preserved WIP and is not part of V1 acceptance.
The next operational proof is creation of a fresh Portfolio task that resumes
Marketplace Central product development.

## Objective

Turn a Codex goal into long-running, restartable, high-quality engineering work
inside Marketplace Central. The harness must give each agent the smallest
sufficient repository context, enforce project architecture and safety, route
work to the right task or subagent, validate in proportion to risk, and persist
enough state that another session can resume without hidden chat history.

The harness accelerates product development. It is not a clean-machine
simulator, a virtual machine, a second CI system, or a generic multi-agent
platform.

## Product Outcome

An operator can give the portfolio task a goal. The task reconciles it with
MNFS, compiles an exact context pack, opens or resumes the appropriate milestone
task, delegates bounded feature work, monitors and steers it, runs the required
gates, requests real QA only when the feature needs it, and returns accepted
commits plus durable evidence. Subsequent sessions resume from repository
artifacts rather than rediscovering the codebase.

## Principles

1. **Context is compiled, not accumulated.** Objective, interfaces, invariants,
   paths, commands, and stop conditions are selected from canonical sources.
2. **One coherent outcome per visible task.** Portfolio and milestone tasks keep
   decisions; feature exploration, implementation, tests, and review use bounded
   subagents.
3. **The repository is memory.** MNFS owns execution state, Git owns code state,
   governance owns machine facts, and validation artifacts own evidence.
4. **Global maximum before implementation.** Planning identifies the owning
   domain, interfaces, consumers, source of truth, wrong legacy surface, and
   smallest durable abstraction before code changes begin.
5. **Risk determines ceremony.** Low-risk work flows quickly; money, tenant,
   credentials, external systems, provider writes, and irreversible actions get
   stronger gates and real evidence.
6. **Parallelize information, serialize collisions.** Read-heavy investigation,
   tests, and reviews may run in parallel. A shared seam has one writer.
7. **Native Codex first.** Use tasks, subagents, skills, goals, worktrees, browser,
   and thread steering already supplied by Codex. Do not build an app-server
   client for V1.
8. **No ceremonial proof.** A gate must catch a real failure mode or it does not
   block product development.

## Explicit Non-Goals

- No cold clone of the current repository.
- No dependency-cache purge or repeated Go/npm/Docker provisioning per feature.
- No local simulation of a new machine or CI runner.
- No AutoGen, vector database, knowledge graph, or custom agent runtime.
- No task per micro-step and no recursive agent swarm.
- No parallel writers in one checkout or shared seam.
- No fake, mock, compile-only, or preflight evidence promoted to live.
- No automatic provider write without the existing actor, policy, linkage,
  idempotency, timestamp, and audit requirements.

## Four-Plane Architecture

### 1. Knowledge plane

Canonical inputs remain:

- `AGENTS.md` and rare nested `AGENTS.md` files for durable operating rules;
- `ARCHITECTURE.md` and ADRs for rationale and frozen decisions;
- `contracts/governance/` for modules, invariants, runtime keys, lanes, shared
  seams, and knowledge routes;
- OpenAPI plus `packages/sdk-runtime` for HTTP contract truth;
- `.mnfs/` for mission, milestone, feature, status, and evidence ownership;
- code and tests for executable truth.

The context compiler resolves a feature's owning module and change concerns,
then emits selectors rather than copying entire documents. A selector names a
path, business reason, and either a Markdown heading, symbol, interface file, or
full-contract read. Source hashes and base SHA make stale packs fail closed.

The total harness-requested initial read set—short bootstrap plus pack
selectors—targets at most 2,000 estimated tokens. F-05 shortens the root
bootstrap and routes detailed workflow through the repo skill. An L2/L3
overflow is allowed only when the pack records why each additional source is required.
Required code reads are ordered and bounded. Broad search is allowed when a
named uncertainty or missing route blocks planning; the discovered stable route
is then added to the registry instead of rediscovered forever.

### 2. Control plane

- One persistent portfolio task owns the active goal, dependency order,
  milestone dispatch, integration order, and final mission audit.
- One visible Codex task owns each active long-running milestone.
- A milestone task delegates bounded work to depth-one subagents and sends
  checkpoints back to the portfolio task.
- Thread tools start, read, steer, interrupt, title, pin, and archive tasks when
  available. Repository scripts never treat thread history as canonical state.
- A repo skill under `.agents/skills/mpc-goal-harness/` implements the workflow
  through progressive disclosure.
- V1 uses Codex's built-in depth-one subagents with bounded prompts. Project
  custom agents are deferred until dogfood proves a repeated role-specific gap.

The app-server API is documented but the local CLI marks it experimental. V1
does not depend on it. `codex exec --json --output-schema` is the deterministic
fallback for a fresh bounded session or future CI.

### 3. Execution plane

Agents edit the normal checkout or a native Codex worktree. Worktrees exist for
parallel Git ownership, not machine isolation. They reuse the normal developer
toolchain and caches. Only genuinely mutable runtime identities—such as an
active Compose project, port, or database—need distinct namespaces when two
runtimes actually run concurrently.

The stable harness entrypoint routes commands but does not contain business
logic. Focused modules own policy, context, state/leases, execution, environment,
PostgreSQL lifecycle, and evidence.

### 4. Validation and QA plane

Every feature declares the proof required by each acceptance criterion. The
risk router selects the minimum sufficient gate:

| Risk | Typical work | Execution and review |
| --- | --- | --- |
| L0 | docs, generated or mechanical change | deterministic check + self-review |
| L1 | one module, no shared seam or live target | TDD + impacted tests + one combined final review |
| L2 | cross-module, API/SDK, migration, composition | contract gate before build + one fixed-commit final review |
| L3 | money, tenant, credentials, real DB/provider, external write | serial owner + independent SPEC/SAFETY and QUALITY once + real QA where applicable |

The impact gate runs in the current checkout with existing dependencies. It
executes governance, formatting, boundary checks, and task-declared tests/builds.
It does not install dependencies, pull images, or invent a full-suite requirement
for every feature.

Target labels stay explicit: `fake`, `ephemeral-postgres`, `live-oracle`,
`live-provider`, `browser`, and `provider-write`.

## Long-Running State Machine

```text
GOAL_CAPTURED
  -> MISSION_RECONCILED
  -> MILESTONE_SELECTED
  -> FEATURE_ELIGIBLE
  -> CONTEXT_COMPILED
  -> GLOBAL_MAXIMUM_PLAN_GATE
  -> LEASED
  -> IMPLEMENTING
  -> IMPACT_VALIDATED
  -> REVIEWED
       -> ACCEPTED
       -> ONE_CORRECTION_BATCH -> FOCUSED_REVALIDATION -> ACCEPTED | REPLAN_OR_BLOCK
  -> INTEGRATED
  -> GOAL_AUDITED
  -> DONE
```

Each transition is reconstructible from a harness checkpoint ID, base SHA,
feature path, allowed paths, shared seams, commands, evidence paths, commit SHA,
and next action. Native task/thread IDs are optional correlation metadata, not
canonical state. A stale lease never deletes, resets, restores, or overwrites work.

## Global-Maximum Plan Gate

Before implementation, the planner must answer:

- Which module owns the business rule and state?
- Which ports/interfaces isolate external systems and other modules?
- Which existing consumers and tests constrain the change?
- Which canonical source owns each fact?
- Is a legacy path being extended when it should be removed or bypassed?
- What provider-neutral abstraction is justified by a current consumer?
- What smaller abstraction would be a patch or duplicate an existing policy?
- Which unknown states must remain explicit instead of becoming zero/default?
- Which contract, migration, SDK, UI, runtime, and QA surfaces are actually
  affected?

L2/L3 work receives a bounded architecture investigation before the plan is
frozen. Model review cannot override deterministic boundary failures.

## Token-Efficiency Contract

- Harness-requested initial-context target, including bootstrap and selectors:
  at most 2,000 estimated tokens; a necessary L2/L3 overflow names each
  additional source and its `overflow_reason`.
- Main task retains decisions and summaries, not raw logs or broad scans.
- Subagent prompts include exact objective, read set, allowed/forbidden paths,
  output schema, and stop conditions.
- Subagent returns are concise and file-referenced.
- Only the active feature is loaded; prior history is linked by artifact path.
- New stable knowledge discovered during work updates a canonical route or
  runbook in the same accepted slice when ownership is clear.
- Handoffs contain commit, changed paths, commands, evidence target, blockers,
  and next action—never a transcript replay.

Prospective metrics are pack size, required-read count, route misses, unrelated
reads, elapsed time, first-pass review result, correction batches, and
target-labelled evidence. Actual model-token telemetry is recorded only when
the active Codex surface exposes it reliably.

## Parallelism and Leases

Read-only investigations, independent tests, and independent final reviews can
run in parallel. Write work can run in parallel only when it starts from the
same accepted base, uses separate worktrees, has disjoint allowed paths, and
owns no common shared seam. OpenAPI/SDK, migrations, composition root,
dependency locks, architecture/ADRs, and provider capability contracts stay
single-writer.

One concurrency slot remains available for orchestration, review, or recovery.
Depth stays one unless the operator explicitly changes it.

## Codex Surface Roles

- `AGENTS.md`: short durable rules and truth routing.
- Repo skill: reusable goal-to-MNFS orchestration method.
- Custom agents: optional follow-up only after a repeated dogfood gap.
- Hooks: optional fast advisory/preflight checks only after a deterministic
  script exists; versioned scripts own mandatory enforcement.
- Tasks/threads: visible milestone control and steering.
- Worktrees/handoff: parallel Git ownership and foreground/background movement.
- Browser: UI QA against the running product.
- `codex exec`: fresh-session and CI-style structured fallback.
- Automations: heartbeat or maintenance only after the workflow is stable; no
  acceptance decision or provider write.

## V1 Implementation Sequence

1. Preserve F-04 evidence and mark its cold design superseded.
2. Remove the cold command, cold-provision lane, snapshot/provisioning code, and
   cold-only tests while retaining reusable redaction/outcome primitives.
3. Add an impact gate whose inventory comes from the active context pack.
4. Add knowledge routes and context selectors for module interfaces, contracts,
   consumers, and validation commands.
5. Add ignored state, leases, resume/recovery, compact handoffs, and the repo
   goal skill using built-in subagents.
6. Add custom agents, nested guidance, hooks, or automation only if dogfood
   produces a concrete repeated gap and an owner-approved follow-up.
7. Dogfood one milestone task and bounded feature flow, including steering,
   review, PostgreSQL or browser/live QA only when its contract requires it.
8. Run deterministic evals built from real repository failures and close M-08.

## Completion Contract

M-08 is complete when:

- a goal routes to the correct MNFS milestone and feature;
- a hash-current pack names exact context, paths, seams, commands, risk, and
  stop conditions; the total harness-requested bootstrap-plus-selector set
  meets the 2,000-token target or carries justified L2/L3 overflow;
- a visible milestone task can dispatch, monitor, steer, and resume bounded
  subagents from repository artifacts;
- a competing seam or out-of-scope write fails before acceptance;
- the impact gate selects and runs only required deterministic checks;
- database, Oracle, provider, browser, and provider-write evidence remain
  distinct and real targets are exercised only when required;
- a fresh Portfolio task resumes product development from mission, Git,
  knowledge routes, and the milestone handoff without transcript replay;
- no active command or criterion depends on cold clone, clean caches, or local
  clean-machine simulation;
- architecture, governance, wiki, MNFS, code, and accepted evidence agree.
