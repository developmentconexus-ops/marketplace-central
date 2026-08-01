# P5 Reconciliation — round 04 (MIS-007-ml-sync)

```yaml
type: planning-review-reconciliation
phase: P5
round: 04
created: 2026-07-31
auditor_artifact: p5-claude-decomposition-audit-r04.md
auditor: cold Claude Opus crew (task a4e749acef0c99f7b; operator-ratified waiver — Sol P5
  retroactive mandatory before status: planned)
input_manifest: p5-input-r04.sha256 (digest b425888df8a8f2705fe198eafe35a266b21850c96fadaaa96f0198b859adf55a)
verdict_received: NEEDS-REVISION (r03 P-1..P-5 all CONFIRMED CLOSED, no prior closure
  reopened; 1 new blocking F-r04-1, 5 advisory F-r04-2..F-r04-6)
disposition: ALL findings ACCEPTED as valid; fold applied in full (blocking + all 5
  advisories); re-audit r05 required
```

## Process note

Task output-file came back 0 bytes — FOURTH instance of the transport fragility. Verbatim
persisted in the SAME turn the notification arrived (recovered from the session transcript
JSONL, entities decoded, provenance header recorded). Rule honored with no gap this round.

## Per-finding disposition (all ACCEPTED, no downgrades)

### Blocking

- **F-r04-1 (`last_success_at = last_full_sync_at` lies for incremental-only entities)** —
  FIXED via the auditor's FIRST arm: IC-05 §Required Outputs now DEFINES the field —
  `last_success_at = GREATEST(last_full_sync_at, last_incremental_at)`, NULL only when
  both NULL — with the writer semantics recorded (`sync_state_repo.go:62,74-79` COALESCE
  freezes `last_full_sync_at` under incremental runs once M-02 F-03 lands) and the
  mandatory negative fixture named in the contract. M-09 F-01 Inputs drops the
  `= last_full_sync_at` equation in favour of the GREATEST definition; F-01 Validation
  gains the negative fixture (old full + recent incremental ⇒ JSON carries the
  incremental; equating to full REPROVA); M-09 F-02 Validation gains the FE side (same
  entity ⇒ badge VERDE/fresh, never "há N dias"/cinza). The field now has a single
  binding authority (IC-05) and the failure `M-09/milestone.md` forbids ("tela de saúde
  que mente") is instrumented at endpoint AND render level.

### Advisory (all applied)

- **F-r04-2 (setter/builder arm silently drops M-08's injection)** — FIRST arm: IC-05
  §seam pins injection as REFERENCE/pointer semantics — the injected reader must be
  observable through the handler ALREADY REGISTERED by M-09; value-receiver builder
  returning a copy (idiom `root.go:854-857`) is named PROIBIDO at this seam; the proof
  obligation moves to the ROUTE (`GET /sync/health` post-injection returns inbox-derived
  stats). Propagated: M-09 F-01 port test asserts via the mounted handler, M-09 milestone
  Done-Means, M-08 F-02 Expected Output.
- **F-r04-3 (IC-03 false locus + phantom move of `DeriveOrderBucket`)** — false prose
  DELETED per doctrine: IC-03 now states the function already lives in the núcleo
  (`orders/domain/order_bucket.go:48`, signature and truth table unchanged) and M-03
  moves only the CALL SITE from read to ingest; deletion of "transport de market" /
  "função MOVE p/ o núcleo" recorded in-line with the refutation source (§1 of
  p5-prerequisites). Briefs were already correct — no brief edit needed.
- **F-r04-4 (`user_id`→installation sourced from a non-invertible resolver)** — M-08 F-02
  Brief + Inputs re-sourced: map = `integration_installations.external_account_id`
  (persisted at `auth_adapter.go:192,261` → `auth_flow_service.go:691`), read through the
  EXISTING installations repo/service (`installation_repo.go:81`) behind a PORT owned by
  the webhook package — never direct cross-module SQL, no new store; `AccessTokenResolver`
  retained ONLY as the token source for the authenticated refetch (explicitly noted as
  non-invertible). Residual sweep caught the same false source in IC-04's DrainInbox row
  ("credential store PRÓPRIO") — also fixed.
- **F-r04-5 (root.go import block unowned; M-07 range excludes :99,101)** — mission matrix
  transversal rule now states the import block belongs to NO milestone (hub-resolved;
  each milestone adds/removes only the imports its own anchored region requires); M-07's
  matrix cell and F-01 Owned paths name the `root.go:99,101` import removals
  (tarifflive/tariffcomposite) alongside the `:828-858` region.
- **F-r04-6 (write-DAG universal omits two real overlaps)** — the enumeration now names
  the guard-allowlist file (owner M-02 F-04; writers M-03 F-03 A/B and M-07 F-01 C/D,
  serialized by the lane B → lane C edge) and the root.go import block (hub-resolved per
  F-r04-5); the universal is scoped to "sem edge nomeado ou resolução do hub registrada".

## Effect on prior artifacts

- `p5-input-r04.sha256` INVALIDATED by the fold edits. New manifest: `p5-input-r05.sha256`.
- No r01/r02/r03 closure reopened (auditor confirmed Part A clean; this fold touched none
  of those loci except additively).
- IC-05 gained two normative pins (last_success_at definition; seam mutation semantics);
  IC-03 and IC-04 each lost one false sentence.

## Next

1. Freeze `p5-input-r05.sha256` over the same 46-file set (post-fold content).
2. Re-dispatch cold decomposition auditor (r05) on the frozen manifest; persist verbatim
   IMMEDIATELY on return.
3. Advance to P6 only on r05 PASS. Sol P5 retroactive touchpoint remains MANDATORY before
   `status: planned` (≥ 2026-08-05).
