# M-08 Execution Guide — Pragmatic Harness V1

## Required Reads

F-10 is the transition: until F-05 changes the repository bootstrap, it obeys
the current root `AGENTS.md` start list and records that larger read set without
claiming the V1 context target.

F-05 must atomically shorten the root bootstrap and this guide without losing
policy. The accepted V1 order is:

1. short root `AGENTS.md` plus the active execution pointer;
2. validated current context pack;
3. only the architecture, contract, interface, consumer, test, and QA selectors
   named by that pack;
4. a targeted discovery read only when a named route gap blocks planning.

The context metric counts every harness-requested initial read, including the
bootstrap and selectors; a full document cannot sit outside the measurement.
Prior feature plans/logs are historical evidence and load only through an
explicit selector.

## Orchestration Contract

- Portfolio task: persistent goal owner and integration authority.
- Milestone task: one visible task with exclusive checkout lease. Request
  `gpt-5.6-terra`/high where the task surface supports it and record the actual
  capability; bounded built-in subagents use the host-supported default rather
  than assuming an unavailable model alias.
- Depth: one. Do not let workers recursively fan out.
- Context: compile/validate the active pack, then read only declared selectors
  and sources. Broad discovery requires a concrete missing-knowledge blocker.
- Writers: one per checkout. Parallel writers require separate worktrees,
  accepted same base, disjoint paths, and no shared seam.
- Worktrees: Git coordination only. Reuse normal dependencies/caches. Namespace
  ports/database/Compose only when concurrent runtimes would collide.
- Reviews: L0 self, L1 one combined final, L2 contract gate plus one final,
  L3 parallel SPEC/SAFETY and QUALITY once on one fixed commit.
- Corrections: one consolidated batch and one focused revalidation. A new
  material finding replans or blocks.
- Evidence: concise, target-labelled, redacted, and linked by relative path.
- Never reset, revert, stash, clean, overwrite, delete unknown state, use WSL,
  expose secrets/PII, or claim fake evidence as live.
- No cold clone, cold provisioning, cache purge, repeated dependency install,
  custom app-server client, or M-09+ product work in this milestone.

## Feature Order

Accepted: F-01, F-03, F-06, F-07, F-08.
Superseded evidence: F-02, F-04.

1. **F-10 pragmatic harness cutover** — retire cold command/lane/tests and add
   context-driven impact evidence in the normal checkout.
2. **F-05 goal and orchestration control plane** — knowledge routes,
   `State.psm1`, leases/resume, repo skill, native milestone task dispatch,
   bounded built-in subagents, optional worktree, and compact handoff.
3. **F-09 eval and dogfood** — deterministic failure corpus, fresh task,
   steering, bounded implementation/review, selected QA, artifact-only resume,
   timing/context metrics, and milestone gate.

## Gate Selection

- Always: context current, allowed paths, shared seams, governance, changed-path inventory.
- Go change: focused package tests with repository `GOCACHE` convention and `gofmt`.
- Frontend change: impacted workspace test/build and SDK-only data access check.
- API change: OpenAPI and SDK together plus contract tests.
- Migration/composition: single-writer seam plus integration proof.
- DB behavior: F-03 ephemeral PostgreSQL.
- Oracle/provider/browser: explicit real target only when acceptance depends on it.
- Provider write: separate command and complete safety/audit gates.

## Output Contract

```text
status:
commit:
changed_paths:
context_pack:
tests_and_targets:
review:
blockers:
next:
```

No raw logs or transcript replay.
