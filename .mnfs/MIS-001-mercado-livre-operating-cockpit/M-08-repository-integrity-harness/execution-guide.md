# M-08 Execution Guide

## Required Reads

Read in order; do not perform broad discovery first.

1. `AGENTS.md`
2. `ARCHITECTURE.md`
3. `.brain/system-pulse.md`
4. `.brain/roadmap.json`
5. `.mnfs/MIS-001-mercado-livre-operating-cockpit/mission.md`
6. `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/milestone.md`
7. `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-contract.md`
8. `docs/superpowers/specs/2026-07-10-codex-portfolio-orchestration-design.md`

Load one feature brief only when that feature starts.

## Orchestration Contract

- Model: `gpt-5.6-terra`; reasoning: `ultra`.
- Environment: current local checkout, serial. Do not create a worktree before
  M-08-C01 passes.
- Preserve all dirty state. No reset, revert, stash, clean, checkout-overwrite,
  or deletion of unknown files.
- Do not implement M-09 through M-12 or change provider listing data.
- Use Cavecrew for bounded investigation/review. Cross-cutting changes remain
  orchestrator-owned.
- Each feature creates `spec.md`, `plan.md`, and `validation.md`, follows TDD,
  and gets independent SPEC then QUALITY review.
- One intentional commit per accepted task. Never stage unrelated paths.
- Keep status and handoff compressed; keep contracts/security warnings clear.

## Feature Order

F-01 is serial and blocking. After its accepted clean baseline:

- F-02 starts before F-03 because environment lanes define DB ownership.
- F-04 depends on F-02 and F-03.
- F-05 may draft documentation after F-02 freezes command names, but its final
  worktree proof depends on F-04.

## Output Contract

```text
status:
done:
evidence:
blockers:
next:
```

No raw log dumps. Provide command, target type, exit code, and artifact path.

