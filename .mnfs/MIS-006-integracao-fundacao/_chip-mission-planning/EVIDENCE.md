# EVIDENCE — MIS-006 mission-planning close

```yaml
chip: mission-planning (MIS-006-integracao-fundacao)
kind: planning (P0→P7) — NOT a code chip; no implementation, no diff, no server
base_sha: 138aac3dff20438d8ddc509daf4171d82e5e45f6
worktree: .claude/worktrees/epic-lehmann-4ffbad
close_condition: deliverables 1-5 written AND P7 readiness gate = Ready
status: CLOSED — planning artifacts written, gate Ready
```

## Nature of this pack (honest scope)

This is a **mission-planning** close, not a feature/chip implementation. The operator's hard
constraint was **do NOT implement — deliver PLANO + INVENTÁRIO DE REFACTOR only**. Therefore:

- No source diff, no code, no migration applied, no server booted, no browser QA. Those are
  **execution-time** concerns owned by the milestones this plan defines.
- The gate is the mission-planning **P7 dual-model readiness gate**, not a chip **P6 dual gate**.
  The `P6-DUAL-GATE: AGREEMENT` marker below is recorded in its P7 analogue form (dual **crew**
  agreement on a frozen content-addressed manifest), with the Sol-model substitution disclosed as
  an operator-authorized waiver (codex quota-wall).

## Deliverables (5/5 written)

| # | Deliverable | Artifact(s) |
|---|-------------|-------------|
| 1 | Inventário de refactor legacy/local-maximum (file:line, refactor/remove/keep + porquê) | `research/refactor-inventory-backend.md`, findings in `docs/design/evidence/INTEGRATION-FINDINGS-D120.md` (F-XLSX-1 / F-SDK-1 / F-ADAPTER-1) |
| 2 | Decomposição milestones + DAG + Ownership & Concurrency + Parallel Execution Plan | `mission.md` §Milestone Strategy + §Parallel Execution Plan; `architecture-map.md`; M-01..M-07 `milestone.md` §Ownership & Concurrency |
| 3 | ICs / contratos (E1-E7 estendidos: E2.1/Eport.1/E4.1/E7.1/E8/E9/E10 + Error Matrix) | `interface-contracts-mis006.md` |
| 4 | Decisões D3 / D6 / F3.7-viability / stale-policy | `mission.md` §Decisions Resolution (D1-D7, F3.7 PENDENTE owner-authority, stale=keep-absent) |
| 5 | Bindings integração real (Sankhya Oracle, ML F3.7 EAN-discovery) | `mission.md` §Real-integration bindings; REQUESTs enumerated below (owner-gated, not yet sent) |

## Gate verdicts

**P7-DUAL-GATE: AGREEMENT** (P7 analogue of P6 dual gate — dual independent cold crew on one
frozen manifest; joint verdict computed, not chosen)

```
final_round:            03 (cap 3 respected)
final_manifest:         planning-reviews/p7-input-r03.sha256
final_manifest_digest:  76522374df7f5cafbf8c6d720ee33fdc7a0f0b2fb155922bab141b3e2ccd0906
claude_cold_crew:       Ready   (folded r01-r03; p7-claude-readiness-r0{1,3}.md)
sol_high (gpt-5.6):     unavailable → rebound to independent Claude cold crew
                        (OPERATOR WAIVER: codex quota-wall, reset 2026-07-25, pre-authorized
                         in dispatch brief; sol-unavailable-p7-r03.md)
joint_verdict:          Ready
readiness-review.md:    written (full round ledger + repair disposition)
```

Round ledger (all blocking findings closed in-artifact, zero downgrade/vote-away):

- **R01** (a9ba45f1) — Needs revision: 5 blocking (★3 M-07 sync_state key drift; ★2 matrix M-07
  migration cell contradiction; ★1 no IC Error Matrix; ★2 undeclared list ordering; ★3 M-07
  ungranted migration on M-02-owned products_mirror) → all fixed.
- **R02** (52ac5829) — Needs revision: 1 new adversarial (★2 catalogPage repointed to
  products_mirror with no ORDER BY → unstable pagination) → fixed (`ORDER BY codigo_produto ASC`
  + M03-C11b). Six criteria held PASS.
- **R03** (76522374) — **Ready**: focused ★2 re-gate (operator directive: "roda só review do 2,
  não full crew toda vez") = **★2 PASS**; six held-PASS unmoved; 2 non-blocking advisories applied.
  Computed seven-★ fold = all PASS.

## Live-proof marker / operator waiver

**No live code marker** — nothing was implemented, so there is no running surface to browser-QA.
This is correct for a planning mission and is the honest state, not a skipped verification.

The mission's **live-proof seams are DEFERRED to execution** and enumerated as owner-gated REQUESTs
(not sent — awaiting operator OK, per side-effectful-action discipline):

1. **Sankhya Oracle db-consult** (hub relays to MNOS): map `TGFPRO/TGFEST/TGFCUS/TGFBAR/TGFTAB/TGFGRU`
   before M-04. — decides M-04 adapter field bindings.
2. **ML F3.7 live T13-T16** (mlprobe, #004-E EANs, active credential NEVER exposed, read-only):
   validates EAN-discovery viability — resolves the F3.7 `PENDENTE` decision.

Sol true-model confirming pass owed post-2026-07-25 (non-blocking) if execution has not begun.

## Manifest of written planning artifacts (this worktree)

Verified against on-disk (Glob, this worktree's repo-root `.mnfs/`):

```
.mnfs/MIS-006-integracao-fundacao/
  mission.md                         (status: planned)
  validation-contract.md
  architecture-map.md
  interface-contracts-mis006.md
  readiness-review.md
  research/{refactor-inventory.md, refactor-inventory-backend.md,
           refactor-inventory-frontend.md, contracts-decisions-scenario.md}
  M-01-sync-state-scheduler/{milestone.md, validation-contract.md}
  M-02-mirror-port-active-source/{milestone.md, validation-contract.md}
  M-03-xlsx-adapter/{milestone.md, validation-contract.md}
  M-04-sankhya-adapter/{milestone.md, validation-contract.md}
  M-05-auto-vinculo/{milestone.md, validation-contract.md}
  M-06-telas-sdk/{milestone.md, validation-contract.md}
  M-07-f37-discovery/{milestone.md, validation-contract.md}   (conditional, gated live T13-T16)
  planning-reviews/{p7-input-r01.sha256, p7-input-r02.sha256, p7-input-r03.sha256,
                    p7-claude-readiness-r01.md, p7-claude-readiness-r02.md,
                    p7-claude-readiness-r03.md, sol-unavailable-p7-r03.md}
```

**Honest gaps in the artifact set:**
- No `p3-*` / `p5-*` Sol reconciliation artifacts on disk. The P3 (co-planner) and P5
  (decomposition audit) Sol touchpoints were NOT run as separate persisted artifacts — the mission
  was authored under the codex quota-wall from the start (same wall as P7 Sol). This is a
  dual-model-completeness gap, disclosed here rather than papered over.
- `p7-claude-readiness-r02.md` was filed at close (round ran live; per-round fold persisted
  retroactively — see its header note).
