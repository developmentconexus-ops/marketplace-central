# Repository-Native Agent Harness Design

```yaml
status: approved_transition_design
owner: Portfolio Orchestrator
approved_by: operator
approved_on: 2026-07-10
transitions_from: docs/superpowers/specs/2026-07-10-codex-portfolio-orchestration-design.md
supersedes_on: M-08 knowledge-authority cutover acceptance
```

## Objective

Build a repository-native development control plane that turns a Codex `/goal`
objective into bounded, restartable, evidence-backed engineering work. The
harness must maximize autonomy, implementation quality, speed, and token
efficiency without sacrificing Marketplace Central architecture, safety, or
live-evidence honesty.

The harness is part of the repository's engineering platform. It is not a set
of prompts and it is not a second application architecture.

## Design Principles

1. **One fact, one owner.** Status, contracts, runtime configuration, and
   architecture rules must not have competing writable sources.
2. **Context is compiled, not remembered.** Workers receive only goal-relevant
   facts selected from canonical sources and current code.
3. **Deterministic controls precede model judgment.** Scripts and tests enforce
   objective rules; agents handle ambiguity, design, and implementation.
4. **Risk determines ceremony.** Review depth, model choice, serial execution,
   and live evidence scale with blast radius.
5. **Outcome beats transcript.** Commits, diffs, tests, target-labelled
   evidence, and runtime state prove completion. Conversation history does not.
6. **Global maximum with YAGNI.** Remove wrong legacy surfaces and redundant
   abstractions, but do not build a generic orchestration product before this
   repository proves the need.

## Explicit Non-Goals

- Do not add AutoGen or another multi-agent runtime.
- Do not replace Codex's native task, subagent, worktree, goal, skill, or hook
  capabilities.
- Do not make natural-language specs the only truth when code, OpenAPI, or
  runtime evidence owns the behavior.
- Do not create a knowledge graph, vector database, or RAG service.
- Do not permit parallel writes to shared seams.
- Do not preserve `.brain` as a compatibility layer.
- Do not claim mocks, fakes, compile-only checks, or preflight output as live
  integration evidence.

## Source-of-Truth Topology

### Durable layers

| Layer | Canonical owner | Purpose |
| --- | --- | --- |
| Repository policy | `AGENTS.md` plus nested `AGENTS.md` files | Stable rules, truth order, commands, safety |
| Architecture | `ARCHITECTURE.md`, `contracts/governance/modules.json`, `contracts/governance/invariants.json`, and ADRs under `docs/architecture/decisions/` | Narrative rationale plus machine-owned module boundaries and invariants |
| HTTP contract | `contracts/api/marketplace-central.openapi.yaml` | API behavior; SDK changes remain atomic with it |
| Governance registry | `contracts/governance/*.json` | Machine-readable modules, invariants, runtime configuration, lanes, shared seams |
| Human operating knowledge | `wiki/` | Explanations and runbooks rendered from or linked to canonical contracts |
| Execution ledger | `.mnfs/` | Mission, milestone, feature status, blockers, accepted SHA, and evidence paths |
| Executable truth | code, tests, harness output, live validation artifacts | Actual behavior and completion proof |

### `.brain` retirement

`.brain` is removed after a migration audit. Current, unique architecture
decisions move to `docs/architecture/decisions/`; current execution state moves
to `.mnfs`; current human guidance moves to the wiki. Historical plans remain
available through Git history and are not copied forward merely for
compatibility.

`AGENTS.md`, handoffs, plans, and runbooks must stop instructing workers to
read or update `.brain`. Historical documents may mention it only when clearly
labelled as superseded history.

## Governance Registry

Use JSON because PowerShell can read it natively, JSON Schema can validate it,
and no new runtime dependency is required. Knowledge-authority cutover is
atomic: an ADR and `AGENTS.md` truth order assign each registry its facts;
duplicated status, configuration lists, and invariant definitions are removed
from narrative documents in the same slice. Semantic drift checks then enforce
the declared one-way ownership instead of mediating two writable truths.

### `contracts/governance/modules.json`

Owns each active module's owner, lifecycle status, allowed dependencies,
ports, source systems, owned state, external side effects, contract paths, and
validation commands. `ARCHITECTURE.md` keeps rationale and topology and links
to this registry instead of repeating volatile status. It distinguishes:

- `owned_state`: PostgreSQL state owned by Marketplace Central;
- `external_fact_source`: Oracle or marketplace facts read through ports;
- `external_side_effect`: provider or customer-facing writes.

### `contracts/governance/runtime-config.json`

Owns canonical environment keys, legacy aliases, owner, sensitivity class,
allowed execution lanes, approved readers, and removal status. Runtime code
must converge on typed configuration owners; new direct environment reads
outside approved configuration boundaries fail governance validation. Wiki
runbooks link to this contract and do not maintain a second key list.

### `contracts/governance/execution-lanes.json`

Defines unit, integration, live, browser, and provider-write commands,
permissions, required inputs, network policy, database policy, evidence type,
and side-effect gates.

Unit execution uses a newly constructed child-process environment with a small
safe allowlist. It never relies on a growing denylist of known external keys.

### `contracts/governance/invariants.json`

Owns stable machine-verifiable invariant IDs with scope, severity, rationale,
and deterministic verifier. `AGENTS.md` retains short human operating rules and
links to these IDs rather than duplicating their full definitions. Initial
invariants include module boundaries,
tenant scope, OpenAPI/SDK atomicity, explicit unknown data quality, provider
write safety, evidence honesty, and `pgxpool.Pool` ownership.

### `contracts/governance/shared-seams.json`

Defines paths that require a single active writer: OpenAPI/SDK, migration
sequence, composition root, dependency locks, architecture/ADRs, and provider
capability contracts.

Every registry file has a JSON Schema under
`contracts/governance/schemas/`. Schema validation is necessary but not
sufficient; semantic checkers compare registry declarations with the current
repository.

## Context Compiler

The context compiler is a deterministic PowerShell module invoked through the
versioned root harness. It accepts an MNFS feature path and base SHA, reads the
governance registry, inspects only relevant repository paths, and emits an
ignored run artifact.

### Input

- active goal-derived MNFS mission, milestone, or feature;
- base commit SHA;
- requested changed paths or owning module;
- risk classification;
- optional explicit external targets for validation.

### Output contract

`scripts/.runs/<run-id>/context-pack.json` contains:

- objective and observable done state;
- base SHA and working environment;
- owning module and risk level;
- ordered required reads;
- applicable invariant IDs;
- allowed and forbidden paths;
- shared-seam leases;
- allowed and forbidden side effects;
- acceptance criterion to test/evidence mapping;
- exact validation commands and target labels;
- model/reasoning recommendation;
- stop conditions, retry budget, and handoff schema;
- source paths and content hashes used to compile the pack.

The pack targets at most 2,000 model-input tokens before a worker opens its
required sources. Large documents are linked, not copied. A stale source hash,
changed base SHA, or missing criterion mapping invalidates the pack before
implementation.

## `/goal` Lifecycle

Codex owns the active thread goal. Repository scripts must not depend on an
undocumented local goal store.

The repo skill bridges the surfaces:

1. Read the active goal already supplied to the Codex task.
2. Classify it as portfolio, milestone, feature, or bounded task.
3. Reconcile the objective into the appropriate `.mnfs` artifact.
4. Run eligibility and risk classification.
5. Invoke the context compiler for the active feature or task.
6. Dispatch bounded workers using the compiled pack.
7. Persist status and evidence in `.mnfs`, not in conversation memory.
8. Audit the original goal requirement-by-requirement before completion.

The skill orchestrates; manifests and scripts own deterministic behavior.

## Orchestration Model

Use a risk-adaptive hybrid rather than a permanent three-level hierarchy.

- The portfolio task is the persistent control plane.
- A milestone is always a durable MNFS boundary, but needs a resident Codex
  task only for complex, long-running, or independently owned execution.
- Feature workers receive fresh bounded context and no portfolio transcript.
- Read-heavy research, test analysis, and review may run in parallel.
- Write-heavy work runs in parallel only from the same accepted base SHA, in
  separate worktrees, with disjoint paths and no shared seam.

### State machine

```text
BACKLOG
  -> ELIGIBILITY_GATE
  -> RISK_CLASSIFIED
  -> CONTRACT_FROZEN
  -> CONTEXT_COMPILED
  -> READY_SERIAL | READY_PARALLEL
  -> IMPLEMENTING
  -> QUICK_VALIDATED
  -> REVIEW_GATE
       -> ACCEPTED
       -> CORRECTION_BATCH
            -> FOCUSED_REVALIDATION
                 -> ACCEPTED
                 -> REPLAN_OR_BLOCK
  -> SEAM_INTEGRATION_QUEUE
  -> INTEGRATED
  -> COLD_GATE
  -> DONE
```

No new micro-correction loop follows `FOCUSED_REVALIDATION`. A new material
finding means contract/context failure and returns to `REPLAN_OR_BLOCK`.

## Risk and Model Routing

| Level | Typical scope | Execution | Review | Default model |
| --- | --- | --- | --- | --- |
| L0 | docs, generated files, mechanical check | parallel when paths are disjoint | deterministic checks; sampled owner review | Luna/medium |
| L1 | local behavior inside one module | bounded worker | one combined final review | Luna/high |
| L2 | cross-module, API/SDK, migration, composition | one writer per seam | contract gate plus final owner review | Terra/high or Luna/high when fully bounded |
| L3 | credentials, real DB/provider, money, tenant, external write, recovery | serial | independent safety/spec review plus live evidence where applicable | Terra/high |

The portfolio orchestrator remains Terra/high. Ultra is never automatic. Sol
or higher reasoning is an explicit escalation for unresolved architecture, not
a default worker setting. These selections remain advisory until the Codex
capability spike proves the installed host accepts the project agent schema,
model aliases, reasoning values, and sandbox declarations.

With four concurrency slots, one remains available for review or recovery.
The control plane must not fill every slot with writers.

## Review Policy

- L0: no mandatory independent reviewer.
- L1: one combined correctness/maintainability reviewer after quick validation.
- L2: contract completeness checked before implementation; one final review on
  a fixed commit.
- L3: independent spec/safety and quality reviews may run in parallel, once,
  against the same fixed commit.
- Findings are consolidated into `spec_gap`, `implementation_defect`, or
  `evidence_gap` before returning to the builder.
- One correction batch and one focused revalidation are allowed.
- A documentation-only correction does not reopen code quality review.
- An evidence gap does not authorize production-code churn unless it proves a
  behavior defect.

## Work and Shared-Seam Leases

Every dispatched writer records:

- owner task/thread;
- base SHA;
- branch/worktree;
- allowed paths;
- shared seams;
- dependent tasks;
- lease status and expiry/recovery action.

Leases live in ignored run state during execution and their accepted outcome is
recorded in the MNFS feature validation artifact. A stale lease never permits
automatic deletion or reset of a worktree.

## Harness Components

Keep `scripts/harness.ps1` as the stable entrypoint and split behavior into
focused modules under `scripts/harness/`:

- `Policy.psm1`: load and validate governance contracts;
- `Context.psm1`: compile and verify context packs;
- `Environment.psm1`: build lane-specific child environments;
- `Execution.psm1`: invoke lanes and capture subprocess results;
- `Evidence.psm1`: write redacted trace and outcome manifests;
- `State.psm1`: eligibility, leases, resume, and recovery state.

Root npm aliases remain thin delegates to the PowerShell entrypoint. Business
rules remain in Go modules, never in the harness.

## Codex Surfaces

### Nested repository instructions

Add short `AGENTS.md` files only where rules differ materially:

- `apps/server_core/AGENTS.md`;
- `apps/web/AGENTS.md`;
- `contracts/AGENTS.md`;
- module-local instructions only for high-risk ownership boundaries.

### Custom agents

Project-scoped `.codex/agents/*.toml` definitions provide narrow roles:

- explorer: read-only, Luna/medium;
- feature-worker: workspace-write, Luna/high;
- cross-module-worker: workspace-write, Terra/high;
- reviewer: read-only, Luna/high;
- live-validator: least privilege, Terra/high, no provider write by default.

Depth remains one. Workers do not recursively fan out.

### Hooks

Project hooks enforce fast deterministic checks such as context-pack presence,
allowed-path boundaries, secret scanning, and stop-time evidence completeness.
Hooks complement rather than replace sandboxing and final validation. A fresh
task capability spike must prove hook discovery, supported lifecycle events,
Windows command resolution, exit semantics, and JSON output. Until that proof
passes, mandatory controls remain in the versioned harness and final gate.

### Repo skill

`.agents/skills/mpc-goal-harness/` contains one focused workflow skill with
progressive references and scripts. It activates for Marketplace Central goal,
mission, milestone, feature, orchestration, and harness requests. It does not
duplicate architecture or registry content.

## Evidence Model

Separate trace from outcome:

- trace: tool calls, commands, exit codes, timings, and decisions, stored as
  ignored redacted JSONL;
- outcome: base/commit SHA, changed paths, tests, target-labelled evidence,
  remaining risks, and acceptance, persisted in MNFS validation artifacts.

Supported target labels remain explicit: `fake`, `ephemeral-postgres`,
`live-oracle`, `live-provider`, and `browser`. Provider-write evidence adds
actor, idempotency record, resolved link, policy, source timestamp, before/after
values, and provider response without secret or buyer PII.

## Harness Evaluation

The harness ships with an eval corpus built from real repository failures.
Each case runs in an isolated fixture and has deterministic graders where
possible.

Initial suites:

1. instruction and context selection;
2. architecture/module boundary safety;
3. runtime environment and side-effect isolation;
4. contract atomicity and tenant safety;
5. validation/evidence honesty;
6. worktree/shared-seam coordination;
7. recovery and resume;
8. token and latency efficiency.

Initial regression cases include missed external environment aliases,
OpenAPI-without-SDK changes, forbidden cross-module imports, unknown-to-zero,
mock-as-live claims, unsafe provider writes, migration conflicts, stale context
packs, writes outside allowed paths, and competing shared-seam writers.

M08 completion requires the deterministic eval corpus to pass and emit timing,
command, and target-classification metrics. Comparative efficiency is a
prospective operational KPI, not an M08 pass condition, because F01/F02 did not
capture a reproducible model-input-token baseline.

After M08, measure at least six representative features:

- at least 20% lower median lead time and model-input tokens versus the M08
  F01/F02 baseline;
- no more than one correction batch per accepted feature;
- higher first-pass acceptance without weaker cold-gate results;
- zero unapproved provider writes, worktree loss, secret/PII leakage, or
  evidence-class inflation;
- no status divergence between `.mnfs`, generated views, and accepted commits.

## Migration Strategy

1. Replan M08 before implementation: add feature briefs and validation criteria
   mapping every new completion requirement; record F02 as blocked baseline
   superseded by the allowlist architecture rather than spending another
   denylist correction.
2. Audit current truth, create the new ADR location, migrate only current unique
   decisions, atomically update `AGENTS.md`, `ARCHITECTURE.md`, wiki, MNFS, and
   the active execution guide, then delete `.brain`. This accepts the
   knowledge-authority cutover and activates this design as successor.
3. Add governance schemas and registries under the newly accepted truth order;
   remove duplicated facts from narrative owners in the same slice.
4. Add semantic drift checks and context compilation.
5. Refactor the unit lane to a child-process allowlist and modularize the
   harness entrypoint.
6. Add MNFS state/lease/evidence integration.
7. Run a Codex capability spike, then add only supported repo-scoped agents,
   hooks, nested instructions, and goal skill. Unsupported model routing stays
   advisory; unsupported enforcement remains in scripts.
8. Complete ephemeral PostgreSQL, cold-gate, and runtime namespace work; then
   run the eval corpus and dogfood the harness on remaining M08 work.
9. Run the full cold gate in a clean worktree and record M08 milestone QA.

Each step is independently testable and ends in an intentional commit. A later
step may not claim success from an earlier step's fake or preflight evidence.

## Completion Contract

The harness is complete only when:

- `/goal` can be reconciled into MNFS and produce a valid context pack;
- risk classification selects bounded execution/review policy;
- unit, integration, live, browser, and provider-write lanes enforce their
  declared boundaries;
- registry/schema/code drift checks pass;
- allowed-path and shared-seam conflicts fail closed;
- trace and outcome evidence are redacted and correctly labelled;
- supported custom agents, hooks, and repo skill operate from a fresh Codex
  task, while unsupported surfaces have an explicit script-based fallback;
- eval corpus passes with recorded deterministic timing and outcome metrics;
- a clean worktree reproduces the cold gate;
- `.brain` and active references to it are removed;
- architecture, wiki, MNFS, OpenAPI/SDK, code, and runtime evidence agree.
