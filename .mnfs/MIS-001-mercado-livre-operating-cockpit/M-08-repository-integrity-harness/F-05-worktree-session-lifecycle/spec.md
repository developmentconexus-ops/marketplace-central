# F-05 — Goal Orchestration Control Plane Spec

```yaml
id: F-05
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: F-05
created: 2026-07-11
updated: 2026-07-11
validation_level: QA-3
lifecycle_scope: feature
```

## Feature ID

F-05-worktree-session-lifecycle

## Problem

The current harness compiles a generic feature pack and runs registered impact
commands, but it has no canonical knowledge-route selector, durable ignored
state, lease, checkpoint, or artifact-only resume contract. Its root bootstrap
also asks every task to read broad documents rather than a compact pack.

## Requirements

- Add schema-validated knowledge routes that select exact repository sources for
  a declared module/change concern and fail closed on stale or unknown routes.
- Compile context sources as ordered selectors, measure the complete bootstrap
  plus selector set, and record a source-by-source `overflow_reason` only for
  L2/L3 packs above the target.
- Add `State.psm1` for ignored run state, exclusive checkout/seam leases,
  canonical checkpoint IDs, safe recovery, resume, and compact handoff checks.
- Add a repository skill that routes a goal through MNFS reconciliation,
  context compilation, risk-selected registered gates, bounded native
  subagents, and artifact-only resume. Native task IDs are optional metadata.
- Keep only a compact bootstrap in root `AGENTS.md` and the active execution
  guide; put operational detail in the skill and the context pack.
- Select deterministic review/QA policies for L0–L3 without executing arbitrary
  command text or treating fake evidence as real QA.

## Non-Goals

- No custom app server, custom agent runtime, hook, automation, Docker, Oracle,
  provider, browser, product, or M-09+ work.
- No worktree creation, clone, cache purge, dependency installation, or VM-like
  isolation. Worktrees remain optional Git coordination.
- No milestone pass or external/live validation claim.

## Design

`knowledge-routes.json` is the compact canonical index. `Context.psm1` resolves
only declared route IDs into ordered selectors and adds the pack bootstrap to
the same measured source list. `State.psm1` writes ignored JSON records under
`scripts/.runs/state`, obtains leases atomically with create-new files, and
never removes or overwrites a stale lease during recovery. Checkpoints name the
feature, base/current commits, paths, evidence, blocker, and next action; a
native task ID may only be correlation metadata.

The skill documents native Codex control surfaces as operator-observed
capabilities, not a persisted task API. The risk router returns registered
command IDs and named review/real-QA requirements; the impact gate remains the
only executor of registered argv.

## Edge Cases

- An unknown route, missing source, mutated hash, or unlisted selector blocks
  dispatch rather than broadening the read set.
- A second writer for the same checkout or seam is rejected before it writes.
- A stale lease is preserved for explicit recovery disposition; recovery never
  resets, deletes, or restores files.
- L0/L1 packs over 2,000 tokens fail; L2/L3 packs need a reason on every excess
  source.
- Recursive subagent dispatch and raw-log handoff are rejected by the skill
  contract.

## Acceptance Criteria

### F05-AC01 — Compact route-selected context

- Traces to milestone criterion ID: `M-08-C09`.
- Proven by: `pwsh scripts/tests/harness-orchestration.tests.ps1` and context
  compiler tests.

### F05-AC02 — Writer ownership and recovery

- Traces to milestone criterion ID: `M-08-C10`.
- Proven by: `pwsh scripts/tests/harness-orchestration.tests.ps1`.

### F05-AC03 — Native capability boundary

- Traces to milestone criterion ID: `M-08-C11`.
- Proven by: repository-skill assertions in
  `pwsh scripts/tests/harness-orchestration.tests.ps1`.

### F05-AC04 — Global-maximum dispatch packet

- Traces to milestone criterion ID: `M-08-C13`.
- Proven by: route and risk fixtures in
  `pwsh scripts/tests/harness-orchestration.tests.ps1`.

### F05-AC05 — Artifact-only checkpoint resume

- Traces to milestone criterion ID: `M-08-C14`.
- Proven by: checkpoint/resume fixtures in
  `pwsh scripts/tests/harness-orchestration.tests.ps1`.

### F05-AC06 — Proportional registered gates

- Traces to milestone criterion ID: `M-08-C15`.
- Proven by: risk-router fixtures and `pwsh scripts/tests/harness-impact.tests.ps1`.

### F05-AC07 — Native worktree coordination

- Traces to milestone criterion ID: `M-08-C06`.
- Proven by: `pwsh scripts/tests/harness-orchestration.tests.ps1` creates a
  detached native worktree at the named SHA, verifies its separate task state,
  and confirms the primary checkout status is unchanged.

## Handoff

- Current status: `spec_ready`
- Next owner: Feature Implementer
- Next action: Write the scoped plan and implement the control-plane seams.
- Required files/evidence: feature brief, milestone contract, compact context,
  state/lease fixtures, and registered-gate evidence.
- Blockers or open decisions: None.
