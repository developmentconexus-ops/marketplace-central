# F-09 Harness Eval and Dogfood — Specification

```yaml
id: F-09
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: M-08
created: 2026-07-11
updated: 2026-07-11
validation_level: QA-3
lifecycle_scope: feature
```

## Feature ID

F-09-harness-eval-dogfood

## Problem

The accepted harness has deterministic safety and orchestration primitives, but
it needs a versioned regression corpus and one independently graded fresh-task
slice to prove that those primitives remain effective without portfolio
transcript history or a cold-machine workflow.

## Requirements

- Version a closed eval corpus with a pinned rejected verdict and stable reason
  for unknown-to-zero, OpenAPI without SDK, forbidden import, mock-as-live,
  unsafe provider write, stale context, out-of-scope path, competing seam,
  risk misroute, unrelated read, resume failure, recursive fan-out, and active
  cold surface.
- Emit a deterministic result manifest per case with case ID, duration, total
  initial-context estimate, required-read count, route-miss/unrelated-read
  counts, target, exit, and relative artifact path.
- Keep the grader in repository code and reject malformed or unknown case
  inputs; the manifest never supplies executable text.
- Run a fresh bounded dogfood worker only for the pinned out-of-scope fixture
  and manifest row. The deterministic grader, not that worker, determines the
  result.
- Record a canonical checkpoint and a fresh artifact-only resume observation.
  Task identifiers remain optional correlation metadata.
- Use only registered current-checkout impact commands and label all executed
  evidence as fake/contract. This feature has no real target requirement.

## Non-Goals

- No product, API, SDK, Docker, network, Oracle, browser, provider, or
  provider-write work.
- No clone, cache purge, dependency installation, custom app server, custom
  agents, hooks, automation, or milestone verdict.
- No claim that fixture evidence proves live integration.

## Design

`contracts/governance/harness-evals.json` is a declarative corpus of closed
case facts and expected outcomes. `scripts/harness/Evals.psm1` grades those
facts with explicit invariant-to-reason mappings and writes a redacted result
manifest below the ignored `scripts/.runs` directory. The harness impact
registry gets one `harness-evals` command ID that invokes the deterministic
test script through structured argv. The test suite asserts every pinned
verdict/reason and checks the result shape and efficiency metrics.

The fresh worker owns exactly the manifest row, test coverage, and fixture
under the brief's dogfood paths; it cannot grade itself. The visible
orchestrator records the dispatch, checkpoint, compact result, and a fresh
continuation that reads the feature artifacts, pack, Git, and checkpoint only.

## Edge Cases

- An unknown case ID, missing required result field, unknown target, or
  non-relative artifact path is rejected deterministically.
- Initial context includes bootstrap and selectors. A value above 2,000 needs
  risk L2/L3 and a per-source overflow reason; unrelated reads must remain 0.
- A route miss must be resolved or explicitly blocked, never silently ignored.
- A task/thread ID may be absent without changing the canonical checkpoint.

## Acceptance Criteria

### F09-AC01 — Deterministic safety corpus

- Traces to milestone criterion ID: `M-08-C16`.
- Proven by: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/harness-eval.tests.ps1`.

### F09-AC02 — Target-labelled dogfood and resume

- Traces to milestone criterion ID: `M-08-C04`.
- Proven by: the eval test and `validation.md` fresh-worker/checkpoint/resume
  evidence; all executed cases report `fake` contract evidence.

### F09-AC03 — Bounded orchestration continuity

- Traces to milestone criterion ID: `M-08-C14`.
- Proven by: canonical checkpoint/resume fixture plus the fresh worker's
  compact handoff recorded in `validation.md`.

### F09-AC04 — Registered proportional gate and pragmatic cutover

- Traces to milestone criterion ID: `M-08-C15`.
- Proven by: current-checkout impact run for the validated F-09 pack and the
  registered risk gate assertions.

### F09-AC05 — Pragmatic cutover remains active

- Traces to milestone criterion ID: `M-08-C17`.
- Proven by: the active-cold-surface eval case and harness alias boundary test.

## Handoff

- Current status: `spec_ready`
- Next owner: Feature Implementer
- Next action: Write the machine work contract and RED evaluator tests.
- Required files/evidence: feature brief, current context pack, corpus, test
  results, dogfood checkpoint, and impact outcome.
- Blockers or open decisions: None.
