# M-08 Execution Guide

## Required Reads

Read in order; do not perform broad discovery first.

1. `AGENTS.md`
2. `ARCHITECTURE.md`
3. `wiki/README.md`
4. `.mnfs/MIS-001-mercado-livre-operating-cockpit/mission.md`
5. `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/milestone.md`
6. `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-contract.md`
7. `docs/superpowers/specs/2026-07-10-repository-native-agent-harness-design.md`

Load one feature brief only when that feature starts.

## Orchestration Contract

- Model: `gpt-5.6-terra`; reasoning: `high`.
- Environment: current local checkout, serial. Do not create a worktree before
  M-08-C01 passes.
- Preserve all dirty state. No reset, revert, stash, clean, checkout-overwrite,
  or deletion of unknown files.
- Do not implement M-09 through M-12 or change provider listing data.
- Use Cavecrew for bounded investigation/review. Cross-cutting changes remain
  orchestrator-owned.
- Each feature creates `spec.md`, `plan.md`, and `validation.md` and follows
  TDD. Review depth follows risk: L0 mechanical, L1 one combined final review,
  L2 contract gate plus final review, L3 independent final spec/safety and
  quality reviews on one fixed commit.
- One correction batch and one focused revalidation are allowed. A new material
  finding returns to replan/block; it does not open a micro-review loop.
- One intentional commit per accepted task. Never stage unrelated paths.
- Keep status and handoff compressed; keep contracts/security warnings clear.

## Feature Order

F-01 is accepted. F-02 is a blocked baseline and receives no further denylist
patch. Remaining order:

1. F-06 cuts knowledge authority over and removes `.brain`.
2. F-07 defines governance contracts and context compilation.
3. F-08 replaces F-02 isolation with fresh child environments and modular
   execution.
4. F-03 completes ephemeral PostgreSQL on the accepted lane contract.
5. F-04 adds evidence/state and the deterministic cold gate.
6. F-05 proves Codex capabilities, goal/skill flow, leases, worktree, and
   runtime namespace isolation.
7. F-09 runs regression corpus, dogfood, fresh-task/worktree proof, and M08
   closure evidence.

## Output Contract

```text
status:
done:
evidence:
blockers:
next:
```

No raw log dumps. Provide command, target type, exit code, and artifact path.
