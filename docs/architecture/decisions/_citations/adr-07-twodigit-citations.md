# ADR-07 (two-digit citations) — what they assert
**Harvested:** 2026-08-05 · **Two-digit citations:** 42 (live code: 6, .mnfs only: 36) · **Three-digit `ADR-007` citations:** 11
**Existing document `007-godror-oci-oracle-runtime.md`:** "Godror with OCI is the canonical Oracle runtime" — DOES NOT MATCH: this is the confirmed collision named in the task brief. The document is about the Oracle Go driver (godror/ODPI-C/Instant Client); every two-digit `ADR-07` citation in code and `.mnfs` is about either the sync scheduler's phase vocabulary or the MIS-004 frontend shell — neither mentions Oracle drivers.

## Assertion A1 — Sync scheduler phase vocabulary is exactly `backfill | incremental | sweep`; only `phase="incremental"` ever resolves `incremental=true` (unrecognized/absent phase — including `"sweep"` — tolerantly resolves `false`, never a hard error); a run's cursor must be non-nil with an explicit phase, and only jobs required to distinguish backfill/incremental/sweep carry the field (the legacy `products` cursor stays untouched) — **amended 2026-08-01** to close a `sweep→true` bug (mission: MIS-007-ml-sync)
- Citations: 32 (live code: 6)
- Verbatim: "ADR-07 (amended 2026-08-01) ratifies the phase vocabulary as backfill | incremental | sweep. Only \"incremental\" resolves"
- Anchors:
  - `apps/server_core/internal/modules/sync/application/scheduler.go:165`
  - `apps/server_core/internal/modules/sync/application/scheduler_test.go:270`
  - `apps/server_core/internal/modules/sync/application/scheduler_test.go:300`
  - `apps/server_core/internal/modules/listings/composition/scheduler_sync_state_integration_test.go:85`
  - `apps/server_core/internal/modules/listings/application/backfill.go:94`
  - `apps/server_core/internal/modules/orders/application/orders_job.go:13`
  - `.mnfs/MIS-007-ml-sync/mission.md:183`
  - `.mnfs/MIS-007-ml-sync/M-02-sync-core-seam/_chip-m02/EVIDENCE.md:82`
  - `.mnfs/MIS-007-ml-sync/planning-reviews/p5-claude-decomposition-audit-r03.md:105`
  - `.mnfs/MIS-007-ml-sync/planning-reviews/p5-reconciliation-r03.md:30`
  - `docs/entregas/ONDA-0-ENTREGA.md:142`
  - `docs/superpowers/plans/2026-08-03-f00-scheduler-pedidos-plan.md:118`

Note: `p5-claude-decomposition-audit-r03.md:105` shows this rule was under active dispute mid-planning — a widened IC-06 draft made `phase` mandatory on every cursor, contradicting the narrow ratified `mission.md:168-170` scope (only new M-04/M-06 jobs) and the code's own `ProductsCursor` (no `phase` field); fixed back to the narrow ratified form in `p5-reconciliation-r03.md:30-34` (P-1).

## Assertion A2 — Frontend retheme-first: the initial FE milestone rebuilds the app shell (paper+green tokens, `data-theme` light/dark, PT-BR copy, canonical nav pills; Mercado/Repasses shown "em breve"; Vínculos kept out of the global nav); all new screens are born inside this shell (mission: MIS-004-mvp-demo)
- Citations: 8 (live code: 0)
- Verbatim: "ADR-07 FE retheme-first: M-03 entrega tokens papel+verde, data-theme light/dark, fontes, pills canônicas"
- Anchors:
  - `.mnfs/MIS-004-mvp-demo/mission.md:91`
  - `.mnfs/MIS-004-mvp-demo/mission.md:98`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-claude-candidate-r01.md:54`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-claude-candidate-r01.md:99`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-sol-counterproposal-r01.md:65`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-reconciliation-r01.md:15`
  - `.mnfs/MIS-004-mvp-demo/M-03-shell-retheme/milestone.md:25`

## Contradictions
- Two unrelated decisions share the token `ADR-07` (MIS-007 scheduler phase vocabulary vs. MIS-004 FE shell retheme). Neither matches the three-digit `007-godror-oci-oracle-runtime.md` document, confirming the collision the task brief flagged as already known.
- Internal MIS-007 planning churn: the P3 candidate `ADR-07` meant "typed fields + raw payload persisted as `json.RawMessage` on ingest entities" (`.mnfs/MIS-007-ml-sync/planning-reviews/p3-claude-candidate-r01.md:73`, reconciled toward shape `A-03` in `p3-reconciliation-r01.md:20`) — a third, distinct meaning that never survived to the ratified `mission.md`, which instead assigns `ADR-07` to the scheduler cursor/phase rule (`mission.md:183`).
- Internal MIS-004 planning churn: `p3-reconciliation-r01.md:14` shows a second reviewer's candidate numbering `ADR-07` for "snapshot validity follows ADR-17" (durability), a fourth distinct reading, before reconciliation settled on the retheme meaning for the ratified number.

## Amendments
- MIS-007's A1 carries an explicit code-level amendment dated 2026-08-01: `scheduler.go:165` — "ADR-07 (amended 2026-08-01) ratifies the phase vocabulary as backfill | incremental | sweep. Only \"incremental\" resolves" — fixing a prior `inferIncremental` bug where `case "incremental", "sweep": return true` treated `"sweep"` as incremental; both `scheduler_test.go` table-driven cases were renamed to `"...falls to tolerant default (ADR-07 has no sweep phase)"` and now assert `false`. Evidence trail: `.mnfs/MIS-007-ml-sync/M-02-sync-core-seam/_chip-m02/EVIDENCE.md:82-102`.
