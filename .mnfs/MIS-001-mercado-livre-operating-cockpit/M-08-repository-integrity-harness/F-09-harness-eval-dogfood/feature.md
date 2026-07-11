# F-09-harness-eval-dogfood

```yaml
id: F-09
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-08
created: 2026-07-10
updated: 2026-07-10
validation_level: QA-3
lifecycle_scope: feature
```

## Mission

MIS-001 Mercado Livre Operating Cockpit.

## Milestone

M-08 Repository Integrity and Deterministic Harness.

## Brief

Prove the completed control plane against deterministic regressions and one
fresh-task/worktree dogfood flow before milestone closure.

## Inputs

- Accepted F-03 through F-08 contracts, real past failure cases, clean accepted
  SHA, versioned cold gate, Codex capability results, and ignored run storage.

## Expected Output

- Versioned eval manifest and isolated positive/negative fixtures.
- Result manifest with case ID, duration, target type, exit code, and artifacts.
- Fresh `/goal` task compiles context, dispatches bounded work, validates,
  hands off, and resumes from repository artifacts.
- Final clean-worktree cold-gate and M08 QA evidence.

## Constraints

- Deterministic graders own objective verdicts; model rubric cannot override a
  failed safety/control case.
- Fake fixtures never prove Oracle/provider/browser integration.
- Comparative token/lead-time improvement is recorded prospectively and is not
  an M08 pass condition.

## Negative Scenarios

- Unknown-to-zero, OpenAPI-without-SDK, forbidden module import, mock-as-live,
  unsafe provider write, stale context, out-of-scope path, or competing seam is
  accepted: fail suite.
- Fresh task requires hidden portfolio transcript: fail dogfood.
- Worktree changes primary checkout or shares mutable runtime identity: fail.

## Validation Expectations

- Every declared case produces its pinned verdict and target classification.
- Two cold-gate runs from one clean SHA have identical command inventory and
  exit classification.
- Fresh task/worktree proof records supported Codex controls or explicit
  deterministic fallbacks and leaves primary checkout unchanged.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution.

## Handoff

- Current status: Briefed.
- Next owner: Feature Implementer after F-03 through F-08.
- Next action: Create `spec.md` and `plan.md` for eval manifest, fixtures, and
  fresh-task dogfood.
- Required files/evidence: eval results, two cold manifests, fresh-task and
  worktree proof, final milestone rollup.
- Blockers or open decisions: All implementation features accepted.
