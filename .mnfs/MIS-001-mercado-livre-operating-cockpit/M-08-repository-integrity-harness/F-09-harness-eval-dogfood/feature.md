# F-09-harness-eval-dogfood

```yaml
id: F-09
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

Prove the completed development control plane against deterministic regressions
and a real fresh-task dogfood flow, including orchestration, steering, bounded
implementation/review, selected QA, and artifact-only resume.

## Inputs

- Accepted F-10 impact gate and F-05 orchestration control plane.
- Accepted F-03/F-06/F-07/F-08 foundations.
- Real repository failure cases and target-labelled evidence rules.
- One pinned dogfood slice: after the eval runner exists, a fresh bounded worker
  adds the negative `out-of-scope-write` fixture and its manifest row under
  F-09-owned test/eval paths. Deterministic tests—not self-review—own its verdict.

## Expected Output

- Versioned eval cases with deterministic expected verdicts and reason codes.
- Result manifest with case ID, duration, total harness-requested initial-context
  estimate, required-read count,
  route-miss/unrelated-read count, target, exit code, and artifact path.
- One visible milestone task compiles context, dispatches bounded work, receives
  checkpoints, steers if needed, validates, reviews, and accepts or blocks.
- A fresh continuation task resumes from repository artifacts without the
  portfolio transcript.
- Evidence rollup consumed by independent milestone QA, based on current
  checkout gates and required real targets, never cold clone or clean-cache proof.

## Constraints

- Deterministic graders own objective outcomes; a model cannot override failure.
- Dogfood may not change product behavior or execute a provider write.
- No cold clone, dependency reprovisioning, cache reset, or full-suite-by-default.
- Comparative token/lead-time improvement is prospective; record available
  metrics honestly and do not invent unavailable model telemetry.
- Native task controls are operator-observed dogfood. `codex exec` may prove a
  fresh structured session but is never claimed as app-visible steering parity.

## Negative Scenarios

- Unknown-to-zero, OpenAPI without SDK, forbidden module import, mock-as-live,
  unsafe provider write, stale context, out-of-scope path, or competing seam is accepted: fail.
- L1 work receives L3 ceremony or L3 work avoids real QA: fail risk-router case.
- Fresh task loads unrelated routes or leaves an unexplained route miss: fail efficiency case.
- Interrupted work cannot reconstruct next action from files: fail resume case.
- Parallel writers share checkout/seam or a task spawns recursively: fail orchestration case.
- Any active cold command/lane/criterion remains: fail cutover case.

## Validation Expectations

- Every case emits its pinned verdict, stable reason, duration, target, and relative artifact path.
- Context fixtures count bootstrap plus selectors, meet the 2,000-token target
  or carry justified L2/L3 overflow, and contain no unrelated module source.
- Fresh dogfood records a checkpoint ID, optional task correlation ID, subagent
  activity, compact handoff, commit, changed paths, selected commands, review
  verdict, and evidence targets.
- A second session states the exact next or completed state using only named artifacts.
- No secret/PII is persisted and no fake result satisfies a real-target criterion.

## Dogfood Slice

- Objective: add one deterministic negative fixture proving an attempted changed
  path outside the pack is rejected before acceptance.
- Allowed ownership: `contracts/governance/harness-evals.json`,
  `scripts/tests/harness-eval.tests.ps1`, and
  `scripts/tests/fixtures/harness-eval/out-of-scope-write/**` only.
- Observable completion: the new case emits the pinned rejection reason and all
  pre-existing eval cases retain their verdicts.
- Independence: the task proves F-05 orchestration and resume behavior; the eval
  runner's deterministic grader proves the fixture result, so the worker does
  not grade itself.

## Criterion Mapping

| Criterion | Ownership | Minimum proof |
| --- | --- | --- |
| M-08-C04 | Integration | Dogfood evidence labels fake versus any selected real QA and contains no secret/PII. |
| M-08-C14 | Integration | Visible task, bounded worker, checkpoint, steer path, and fresh continuation are observed. |
| M-08-C16 | Primary | All pinned eval cases plus the exact dogfood slice and efficiency measurements pass. |
| M-08-C17 | Integration | Eval confirms no active cold command/lane/criterion remains. |

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution.

## Handoff

- Current status: Briefed after pragmatic replan.
- Next owner: Feature Implementer after F-05 acceptance.
- Next action: Create spec/plan for deterministic corpus and one dogfood flow.
- Required files/evidence: eval manifest/results, fresh-task dispatch, checkpoint,
  resume proof, impact-gate evidence, and final milestone rollup.
- Blockers or open decisions: None; the dogfood slice is pinned above.
