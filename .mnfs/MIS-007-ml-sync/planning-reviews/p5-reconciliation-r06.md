# P5 Reconciliation — round 06 (MIS-007-ml-sync)

```yaml
type: planning-review-reconciliation
phase: P5
round: 06
created: 2026-07-31
auditor_artifact: p5-claude-decomposition-audit-r06.md
auditor: cold Claude Opus crew (task aa9d98bf0fcdd8c7d; operator-ratified waiver — Sol P5
  retroactive mandatory before status: planned)
input_manifest: p5-input-r06.sha256 (digest 5fcb02fbc59f470a6663bc47c6d12c2704a228ab061c6da979845f1697535bd6;
  auditor recomputed 46/46 OK before reading)
verdict_received: NEEDS-REVISION (PART A: all F-r05-1..6 CONFIRMED CLOSED, zero r01–r04
  closures reopened; 2 blocking F-r06-1/F-r06-2 + 4 advisory F-r06-3..6)
disposition: ALL findings ACCEPTED as valid; fold applied in full (2 blocking + 4
  advisories + same-class residual sweep); re-audit r07 required
```

## Process note

Verbatim recovered from the session transcript JSONL (SIXTH round on this recovery path)
and persisted in the SAME turn the notification arrived. First extraction attempt picked a
truncated later-line candidate (4 chars); corrected by selecting the LONGEST
task-notification string across all matching JSONL lines — 34,337 chars persisted. Rule
honored with no gap.

## Per-finding disposition (all ACCEPTED, no downgrades)

### Blocking

- **F-r06-1 (layer-3 detail mandate unscoped by fee_kind — M-02 writer would refuse M-06's
  freight row)** — FIXED via the auditor's FIRST arm (scope the guard to commission; the
  freight row IS intended layer 3 — realized cost from the actual shipment):
  - IC-01 §Canonical Examples rejection line now reads "camada 3 **fee_kind=`commission`**
    com `detail` sem `sale_fee_unit`/`quantity`" and names the freight row (M-06 F-02,
    subject_type=order, origem=api_shipment) as detail-NULL-accepted.
  - IC-01 §Persistence Expectations gains the sibling bullet defining camada-3 frete once
    (value = custo seller do shipment; detail NULL permitido — no sale_fee_unit/quantity
    decomposition exists for freight).
  - `M-02/F-02` guard clause and Negative Scenario both scoped to fee_kind=commission,
    with the freight-accepted side made an explicit test obligation ("o teste cobre os
    DOIS lados").
  The two ratified briefs (M-02 writer, M-06 producer) are now simultaneously
  implementable; M-06 Done Means (frete_seller in the decomposition) and IC-03's canonical
  example (− 22.90 freight term) are reachable.
- **F-r06-2 (ResolveListingFees tuple lacked `detail`; M-07 EARS unsatisfiable against the
  port)** — FIXED via the auditor's FIRST arm (smaller edit, preserves M-07 F-02
  provenance surface): IC-01 §Required Outputs tuple is now
  `{value, value_type, currency, layer, detail, origem, coletado_em}` with detail defined
  ONCE, scoped by layer (jsonb VERBATIM of the resolved row; camada 2 = canonical 5-key
  tuple; NULL when the row has none) — per the auditor's interaction note with F-r06-1,
  the semantics are stated in one place, never re-invented per consumer. `M-02/F-02` EARS
  restates the full tuple. `M-07/F-01:43-44` (detail.percentage_fee/fixed_fee da camada 2)
  is now satisfiable through the port alone — unchanged, as it was the correct promise.

### Advisory (all applied)

- **F-r06-3 (missing M-06⤳M-05 soft edge at mission level)** — edge added to the DAG
  block and as a row in the edge-justification table with the forcing rationale already
  declared in `M-06/milestone.md:50-51` (auditoria 3→2 precisa de camada 2 populada — MUDA
  sem ela, não quebra) and the same non-blocking qualifier as M-07⤳M-05. Hub now sees the
  soft ordering when adjudicating lane-C close order.
- **F-r06-4 (F-r05-5 residual class: two more 3-key detail enumerations)** — both fixed
  with the full 5-key tuple: `M-05/F-01` Validation fixture (subset assertion would let a
  writer omitting financing_add_on_fee pass green — lesson chip-import-chain recorded
  in-line) and `M-05/F-03` DTO shape (the FE contract locus — silently dropping
  financing_add_on_fee from /listings was the more consequential of the two).
- **F-r06-5 (four anchor line-ranges drifted/mischaracterized)** — all corrected to
  measured values: M-07/F-01 root.go range `844-850` → `845-851` with `:844`
  (`var tariffResolver`) explicitly marked SOBREVIVE (the range instructs a deletion —
  off-by-one both directions was the dangerous case); IC-05 builder-idiom `854-857` →
  `:856` (the actual `WithCalc` line); writer.go idiom `74-105` → dual-range `74-95`
  (`upsertSQL`) + `104-112` (`keepAbsentSQL`) at M-02/F-02 and M-04/F-02 — residual sweep
  found and fixed the SAME stale range in four more manifested loci the auditor did not
  enumerate (IC-03:99, IC-07:97, mission.md:95, mission.md:151); `listings_test.go:25`
  re-labelled as the required-substring list with the regex idiom re-anchored at `:101`
  (`createTableBody`) in codebase-ingest-side.md and mission.md ADR-12. Note: r01's A-2
  yes-if had set `74-105`; r06 measured that the range bisects `keepAbsentSQL` — the
  dual-range citation supersedes it (recorded here; no r01 closure semantics reopened, the
  idiom claim itself was always true).
- **F-r06-6 (two data-shape briefs missing ## Inputs/Outputs)** — sections added citing
  the binding sources already named in each brief body: `M-08/F-02` (IC-04 §Enums status
  machine + IC-04 §Fields; outputs = IngestOrder call + the REAL WebhookStatsReader impl
  with IC-05 §InboxHealth shape and reference-semantics injection) and `M-04/F-02` (IC-06
  run-complete rule + IC-07 E3 columns; outputs = last_seen_at advance, absent_since only
  on complete runs, status verbatim from provider).

## Effect on prior artifacts

- `p5-input-r06.sha256` INVALIDATED by the fold edits. New manifest: `p5-input-r07.sha256`.
- No r01–r05 closure reopened (auditor Part A clean across all seven spot-checks; this
  fold touched prior-closure loci only additively — the one intersection, r01 A-2's
  `74-105` range, is superseded by measurement and recorded above).
- IC-01 gained: detail in the resolution tuple (single-source, layer-scoped), the
  commission scoping of the layer-3 detail mandate, and the camada-3 freight persistence
  bullet. Mission DAG gained the M-06⤳M-05 soft edge. Two briefs gained Inputs/Outputs.

## Next

1. Freeze `p5-input-r07.sha256` over the same 46-file set (post-fold content).
2. Re-dispatch cold decomposition auditor (r07) on the frozen manifest; persist verbatim
   IMMEDIATELY on return (transcript-JSONL recovery path; pick the LONGEST candidate).
3. Advance to P6 only on r07 PASS. Sol P5 retroactive touchpoint remains MANDATORY before
   `status: planned` (≥ 2026-08-05).
