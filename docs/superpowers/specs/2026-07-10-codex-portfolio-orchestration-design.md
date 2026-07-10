# Codex Portfolio Orchestration Design

## Objective

Use one persistent Codex thread as Portfolio Orchestrator and one fresh
`gpt-5.6-terra` thread per milestone. Preserve MNFS as execution truth while
keeping prompts, agent returns, and durable handoffs token-efficient.

## Operating Model

- The portfolio thread owns dependency order, dispatch, status, review,
  integration order, and mission-level gates. It does not implement in files
  concurrently owned by a milestone thread.
- Each milestone thread receives one bounded context pack, owns one checkout
  or worktree, and delegates worker-sized tasks plus independent SPEC/QUALITY
  reviews.
- M-08 runs serially in the current checkout because the uncommitted M-03 to
  M-06 state is absent from HEAD. Later independent milestones use native
  Codex worktrees only after M-08 proves a clean baseline.
- Shared seams are serialized: OpenAPI plus SDK, migration numbering,
  `composition/root.go`, root dependency manifests, architecture/ADRs, and
  provider capability contracts.

## Context Pack Contract

Target: at most 2,000 input tokens before required source reads.

Every dispatch names:

1. concrete objective and observable done state;
2. exact required-read paths in order;
3. frozen decisions and interface contracts;
4. allowed and forbidden paths;
5. side effects allowed and forbidden;
6. exact validation commands and evidence target type;
7. commit and handoff contract;
8. known blockers and assumptions.

Workers do not inherit portfolio history. Subagents use `fork_turns="none"`
with self-contained prompts. Broad repository search requires a concrete
blocker; required reads and targeted `rg` come first.

## Token and Output Policy

- Milestone budget guide: 60% implementation, 20% validation, 10% independent
  review, 10% handoff.
- Cavecrew handles bounded locate/edit/review tasks. Cross-cutting feature work
  stays with the milestone orchestrator or a dedicated feature worker.
- Agent status and handoff use compressed output:

```text
status:
done:
evidence:
blockers:
next:
```

- Raw logs stay in ignored validation artifacts. Conversation output carries
  only decisive lines, exit code, target type, and artifact path.
- Specs, contracts, security warnings, irreversible actions, and code remain
  uncompressed when compression could change meaning.

## Evidence Honesty

Every result labels its target as `fake`, `ephemeral-postgres`, `live-oracle`,
`live-provider`, or `browser`. Mocks never prove live integration. Secrets and
buyer PII never enter prompts, logs, screenshots committed to the repository,
or validation artifacts.

## Dispatch Lifecycle

1. Portfolio checks dependency DAG, clean base SHA, owned paths, and baseline
   gate.
2. Portfolio creates a project thread using the requested model and local or
   worktree environment.
3. Milestone produces feature specs/plans just in time, executes TDD, and
   obtains independent reviews.
4. Milestone returns a compact handoff with commits, changed paths, commands,
   evidence types, blockers, and next exact action.
5. Portfolio reads the handoff, reruns integration-owned gates, and either
   accepts, returns a scoped correction, or blocks.
6. Merge order follows shared-seam ownership, not completion time.

## Parallelism Rules

Parallel implementation is allowed only when branches start from the same
validated base and owned paths do not overlap. Planning/research may run in
parallel earlier. No parallel writes to OpenAPI/SDK, migrations,
`composition/root.go`, package locks, architecture/ADRs, or the same module.

## Portfolio Headline Order

- M-08 Repository Integrity and Deterministic Harness — blocking, serial.
- M-09 Canonical Product Identity and Oracle Cutover.
- M-10 Provider Runtime Consolidation.
- M-11 Durable External Write Execution.
- M-12 Listing Identity/SKU Synchronization.
- Resume blocked M-06 with a paid resolved-link live scenario and full cold
  gate.
- Execute existing M-07 Commercial Intelligence only after M-06 passes.

The present SaaS scope is tenant-safe and provider-neutral internal operation.
Billing, tenant onboarding, and multi-company authentication are deferred by
YAGNI; no design may make them harder to add later.

