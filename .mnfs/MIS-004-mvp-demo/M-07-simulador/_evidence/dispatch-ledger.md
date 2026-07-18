# M-07-simulador — Dispatch Ledger

Chip: CHIP-M07 · branch `chip/m07-simulador` · base `8b6c4b3093f9465cd3b91209b054af4fa702171a`
Hub: local_7c41fdd8-dd6f-487c-a7ea-1b511e214fa9 · Contingency lane D-23 (codex DEAD → Claude-only)

Row written AT DISPATCH TIME (role · model · effort · prompt-pack path), completed with verdict/output artifact.

| # | Phase | Role | Model | Effort | Prompt-pack | Output/verdict artifact | Status |
|---|-------|------|-------|--------|-------------|------------------------|--------|
| D0 | boot | codebase investigator (base map) | sonnet | default | (inline) | agent a200b36ff4995ed09 report (in-transcript) | done |
| D1 | P2 | batch planner (cold Opus) | opus | default | `_evidence/packs/P2-planner.md` | `_evidence/M07-plan.md` (14 slices) | done · accepted-with-corrections (SECTION 0: C-A no 4th DIFAL table→columns on pricing_difal_rates, 0055–0058+0059 reserve, fixture 51→55; C-B npm not pnpm) |
| D2 | P3 | F01-S1 implementer (decimal helper) | sonnet | default | `_evidence/packs/IMPL-BINDINGS.md` + M07-plan F01-S1 | `pricing/domain/decimal.go`+`decimal_test.go` GREEN (5 tests), carrier commit `bb0a276f` (see F-boot-2) | done · orch-reviewed accepted (formal P4 batches w/ S2) |
| D3 | P3 | F01-S2 implementer (CalcProfile+DIFAL seed) | sonnet | default | `_evidence/packs/IMPL-BINDINGS.md` (rev: orch sole committer) + M07-plan F01-S2 | (pending) | dispatched |

## Findings log (for hub ratification, core §0)
- F-boot-1: chip-m07 worktree was NOT provisioned by launch flow (siblings m05/m08 had theirs). Session home `sleepy-wing-6d7500` is not a git worktree — resolves to the shared main checkout. Boot `git checkout -b` therefore transiently moved the shared main HEAD (reflog HEAD@{2}); hub recovered via D-45 (main@896c8659, clean). Self-resolved: created `.claude/worktrees/chip-m07-simulador` off 8b6c4b30 matching sibling convention. Recommend launch flow pre-provision chip worktrees. NON-BLOCKING — batch into first substantive event.
- F-boot-2 (SHARED-INDEX COMMIT RACE): the chip worktree has ONE git index shared by the orchestrator process and every implementer subagent (subagents run git via Bash cwd inside the same worktree). During F01-S1: the S1 implementer `git add`ed its 2 write_set files; the orchestrator concurrently `git add`ed `M07-plan.md` and committed — the commit swept the implementer's already-staged files into the orchestrator's commit `bb0a276f` ("docs(m07): fold HUB v2 SDK-PRICING-M07…"). Net: decimal.go/decimal_test.go ARE correct + GREEN + in HEAD, but landed under a docs-labelled commit, not their own message/trailer. No history rewritten (prohibited + shared-worktree unsafe). **Mitigation adopted (this chip, effective S2):** ORCHESTRATOR IS SOLE COMMITTER — implementers write+test+report only, never touch the index; orchestrator serializes one clean commit per reviewed slice. Recommend hub ratify sole-committer for shared-index worktrees (or give each subagent its own worktree/index). NON-BLOCKING — batch into first substantive event (COMMITTED ic04-ports).
