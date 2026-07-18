# M-04 vinculos-import-ui — Dispatch Ledger

Chip: CHIP-M04 · milestone branch `chip/m04-vinculos-import-ui` · base SHA `28b8447c82fc05f626d3e42404fa3267cc82e953`
Milestone worktree (all commits/dispatches/lanes/evidence): `C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-m04-vinculos-import-ui` (hub-prepared; D-correction 2026-07-18). Orchestrator session cwd is scaffold `hungry-bell-ddf2a3` — every command cd's into the milestone worktree (PowerShell resets cwd between calls).
Codex binary: `C:\Users\leandro.theodoro\.codex\plugins\.plugin-appserver\codex.exe` (0.145.0-alpha.18, per profile §3 F-ENV-9)
Scratchpad (active session): `C:\Users\LEANDR~1.THE\AppData\Local\Temp\claude\C--Users-leandro-theodoro-Documents-marketplace-central--claude-worktrees-hungry-bell-ddf2a3\6631ff9c-7132-4829-b19e-f9d73eab1cbe\scratchpad`
(prior session 2ec05736 died mid-D01 with no output — re-dispatched here.)

| # | role | model/effort | path | prompt | output | status |
|---|---|---|---|---|---|---|
| D01 | investigator (baseline read: product_links, OpenAPI, SDK, FE seams) | gpt-5.6-luna / medium | OS-process (re-dispatch; prior companion attempt died with session) | `scratchpad/prompt-d01-investigator.md` | `evidence/d01-investigator-baseline.md` | COMPLETED exit 0 2026-07-18 |
| D02 | investigator (gap-closer: SDK gen-vs-handauthored + new-domain-client recipe, matcher/anchor data feasibility, FE tab/modal/toast pattern + web-query invalidation) | gpt-5.6-luna / medium | OS-process | `scratchpad/prompt-d02-investigator.md` | — | **FAILED exit 1 (bvhk0fs5c): codex USAGE LIMIT exhausted until Jul-25 — post-demo. F-ENV finding.** |
| D02b | investigator (same prompt as D02) — Claude fallback (codex quota down) | Claude sonnet subagent | Agent tool SYNC | `scratchpad/prompt-d02-investigator.md` | `evidence/d02-investigator-gaps.md` | pending BLOCKED ruling on Claude-fallback lane; read-only so proceeding |
| P2 | batch planner (plan-only; slice cards S1-S9, write-sets, contract-satisfiability IC-01-A2/IC-05, verification map C01-C05, DAG, migration ledger, S1 fixtures) — Claude-only lane D-23 | **cold Opus** subagent (Plan agent, read-only) | Agent tool SYNC | inline plan-only prompt (inputs: planning-readiness-M04.md + d01/d02 + ic01-amended-A2.md + milestone/feature/validation-contract) | `PLAN-M04.md` | COMPLETED 2026-07-18 · open_questions EMPTY = dispatch-ready. Plan agent read-only ⇒ orchestrator persisted verbatim. Key reconciliations: SDK-01 inline index.ts (productLinks.ts retired), routes `/product-links/link-resolutions/*` (live prefix), migrations 0065/0066/0067, undo single+batch, single-worktree ⇒ serialize implementers. |

## Claude-only contingency lane — RATIFIED by hub (D-23, 2026-07-18)

Hub ruling on F-ENV-M04-1. Profile §12 tail "Contingency lane — codex quota outage" (commits 60da30b4 + 4d1a7cfb on main; do NOT pull — lane text is operative binding). Expiry = codex return 2026-07-25 OR MIS-004 demo close, whichever first.

| Role | Contingency binding |
|---|---|
| P2 planner | cold Opus SUBAGENT (fresh context, plan-only; sonnet fallback if Opus limits bite). Plan outputs unchanged (write-sets, contract-satisfiability vs IC-01/IC-05, verification map). |
| Implementer | sonnet subagent (anti-slop contract + slice cards unchanged). |
| Slice reviewer | independent sonnet; implementer ≠ reviewer. |
| P6 dual gate | cold Opus + INDEPENDENT second sonnet, separate dispatch, ADVERSARIAL-REFUTE prompt (replaces GPT side; refutation framing MANDATORY — compensates lost cross-vendor diversity). Agreement/reconciliation rules unchanged. |
| Investigator | sonnet (haiku ok trivial greps). |

RESOLVED by hub 2026-07-18 (main commit 0937621b) — all 4 open questions closed, P2 UNBLOCKED:
- **Finding A → GRANTED (D-24, lock SDK-01):** temporary ADDITIVE-ONLY lock on `packages/sdk-runtime/src/index.ts`, scoped to product-links block. New batch-preview/batch/undo methods + enriched `ProductLinkCandidateItem` (confidence/band/reasons) inline adjacent to existing 7-method block (:1508-1545, types :850-949). ZERO edits/reformat of other domains. Lock releases at CLOSED; CLOSED payload MUST call out the index.ts diff. `productLinks.ts` new-file grant RETIRED. M-02 F-04 may also touch index.ts additively (market block) — disjoint regions, hub resolves textual merge conflict.
- **Finding B → RATIFIED (D-25, IC-01 Amendment A2):** written into `.mnfs/MIS-004-mvp-demo/research/identity-matching-interface-contract.md` @ 0937621b = verdade única (S1 fixtures read the FILE, not summaries). Anchors cross-side = seller_sku(exact→codprod, ACCEPT-grade ONLY in concordant pair) + ean(corroboration, unproved flag). seller_sku ALONE = MEDIA/REVIEW ("EAN ausente ⇒ max REVIEW" binding). Auto-ACCEPT/ALTA = seller_sku+ean concordant same codprod, no hard-neg. Title-only=BAIXA. SKU/EAN conflict or title hard-neg = cap BAIXA+AGAINST. marca/refforn = UNAVAILABLE in reasons[] (ADR-17). Thresholds 85/50 unchanged. F-01 authors ≥8 own fixtures incl Doka hard-negative. M-02 propagation hub-side (no M-02 change).
- Finding C: page-local tabs/modal/toast + page-local query keys — sound, no ruling needed.

## F-ENV-M04-1 (2026-07-18) — codex quota exhausted pre-demo
codex-cli usage limit hit after D01; resets Jul-25 (post-demo 2026-07-20). Blocks all codex roles (planner Sol-medium, implementers Luna-high/Sol-low, dual-gate GPT side). Mission-wide (all wave-B chips). BLOCKED sent to hub for ruling on Claude-only execution lane. Sanctioned already: sonnet fallback IMPLEMENTER (§1). Read-only investigators fall back to Claude sonnet without ruling.
