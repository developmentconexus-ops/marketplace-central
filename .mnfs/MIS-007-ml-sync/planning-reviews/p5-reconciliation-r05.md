# P5 Reconciliation — round 05 (MIS-007-ml-sync)

```yaml
type: planning-review-reconciliation
phase: P5
round: 05
created: 2026-07-31
auditor_artifact: p5-claude-decomposition-audit-r05.md
auditor: cold Claude Opus crew (task aad8fe745b084ab00; operator-ratified waiver — Sol P5
  retroactive mandatory before status: planned)
input_manifest: p5-input-r05.sha256 (digest eefe38a5b903af0c54baaaac58436a32b68eca033bdaefdcfee2d19e5e77c22e)
verdict_received: NEEDS-REVISION (r04 F-r04-1..6 all CONFIRMED CLOSED, no prior closure
  reopened, 28 anchor classes verified EXACT against repo; 2 new blocking F-r05-1/F-r05-2,
  4 advisory F-r05-3..F-r05-6)
disposition: ALL findings ACCEPTED as valid; fold applied in full (2 blocking + 4
  advisories); re-audit r06 required
```

## Process note

Task output-file came back 0 bytes — FIFTH instance of the transport fragility. Verbatim
recovered from the session transcript JSONL (html.unescape + `<result>` extraction) and
persisted with provenance header in the SAME turn the notification arrived. Rule honored
with no gap this round.

## Per-finding disposition (all ACCEPTED, no downgrades)

### Blocking

- **F-r05-1 (M-09 scan predicate cannot return products)** — FIXED per the yes-if verbatim:
  `M-09/F-01` Brief now reads "a leitura varre TODAS as rows de `sync_state` do tenant,
  independente de `installation_id` — INCLUI o sentinela de escopo ERP
  `installation_id = "erp"` (`sync/composition/scheduler.go:11`), onde vive a row de
  products"; the negative scenario at the old `:83` is re-worded to "Tenant sem NENHUMA
  row em `sync_state`" with the empty case explicitly defined as absence of ROW, not
  absence of ML installation. This restores consistency with IC-05:91 (entity-agnostic),
  F-01's own EARS (products entity carries real timestamps), and M-09 Done Means
  ("products real"). Repo fact confirmed: products register under the ERP sentinel
  (`products_job.go:50` via `InstallationScopeERP` at `scheduler.go:11`), so an ML-only
  scan returns zero products — the Done Means would be unattainable.
- **F-r05-2 (audit formula drops `financing_add_on_fee`)** — FIXED via the auditor's
  chosen arm: `M-06/F-03` formula now reads `esperado_unit = detail.percentage_fee ×
  unit_price/100 + detail.fixed_fee + detail.financing_add_on_fee` at BOTH loci (Brief
  and the Constraints "NUNCA do amount" pin), with the rationale recorded (dropping the
  component opens a PERMANENT false `tarifa` divergence on any listing with installment
  financing). IC-01's canonical camada-2 example made arithmetically self-consistent:
  `"value":16.45` → `15.99` (= 12.5% × 79.90 + 6.00 + 0), with the derivation note added
  under the example.

### Advisory (all applied)

- **F-r05-3 (listings additive-lock grant narrower than M-05's real write-set)** — mission
  matrix transversal rule widened: grant = `listings/application/**`,
  `listings/transport/**` and `listings/adapters/postgres/repository.go`, additive-only
  pós-close do M-04; M-05's Go-packages matrix cell aligned to the same enumeration.
  Write-DAG now matches the F-03 Owned paths.
- **F-r05-4 (IC-05 "sem ela QA reprova" contradicts M-09 SOFT dependency)** — arm 1: the
  failure condition is scoped to the MISSION/live criterion — it fires only after
  M-04/M-06 register phase-bearing jobs and the field stays NULL. Pre-fix, uniform NULL
  is honest per IC-05 §NULLs and the M-09 gate PASSES with it; M-02 F-03 remains SOFT in
  M-09's dependency block (unchanged — already correct).
- **F-r05-5 (M-05/F-01 EARS enumerates 3 of 5 detail keys)** — EARS now names the full
  5-key canonical tuple (percentage_fee, fixed_fee, financing_add_on_fee, price_used,
  listing_type_id). Residual sweep caught the SAME 3-key enumeration in M-05
  `milestone.md` Done Means — also fixed (same finding class, recorded here).
- **F-r05-6 (mission Handoff stale at P4)** — Handoff updated to actual state:
  planning_phase = decompose (P5), P4 closed, audit loop in progress (r02–r05 folded),
  next action = close the decomposition audit then P6/P7. Waiver rows preserved VERBATIM;
  `planning-reviews/p5-*` added to required artifact paths.

## Effect on prior artifacts

- `p5-input-r05.sha256` INVALIDATED by the fold edits. New manifest: `p5-input-r06.sha256`.
- No r01–r04 closure reopened (auditor Part A clean; fold touched none of those loci
  except additively).
- IC-01 gained a self-consistent canonical example + derivation note; IC-05 gained the
  live-criterion scoping of the incremental precondition; mission.md gained the widened
  grant and a truthful Handoff.

## Next

1. Freeze `p5-input-r06.sha256` over the same 46-file set (post-fold content).
2. Re-dispatch cold decomposition auditor (r06) on the frozen manifest; persist verbatim
   IMMEDIATELY on return (assume 0-byte output-file; use the transcript-JSONL recovery
   path).
3. Advance to P6 only on r06 PASS. Sol P5 retroactive touchpoint remains MANDATORY before
   `status: planned` (≥ 2026-08-05).
