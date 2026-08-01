# P5 Reconciliation — round 02 (MIS-007-ml-sync)

```yaml
type: planning-review-reconciliation
phase: P5
round: 02
created: 2026-07-31
auditor_artifact: p5-claude-decomposition-audit-r02.md
auditor: cold Claude Opus crew (task a2e8de99471e5fb45; operator-ratified waiver — Sol P5
  retroactive mandatory before status: planned)
input_manifest: p5-input-r02.sha256 (digest fc0522e1665949cf573b72e5ee855105ceef5b14d3f0564cde00dabd3ebb1c6c)
verdict_received: NEEDS-REVISION (F-1..F-10 + A-1..A-9 all CONFIRMED CLOSED; 4 new blocking
  N-1..N-4, 6 new advisory N-5..N-10)
disposition: ALL findings ACCEPTED as valid; fold applied in full; re-audit r03 required
```

## Process note

Task output-file came back 0 bytes; the verdict arrived only in the task-notification and
was persisted VERBATIM in the same turn (rule from r01 honored — no recovery needed this
time, but the 0-byte output-file is a second instance of the same transport fragility).

## Per-finding disposition (all ACCEPTED, no downgrades)

### Blocking

- **N-1 (job return-type change vs IC-06 Must Preserve)** — FIXED via the auditor's FIRST
  arm: M-02 F-03 Brief + Expected Output now pin that the run type is DERIVED from the
  mandatory `cursor.phase` of the returned cursor (IC-06:136-137, same source IC-05 uses);
  `JobFunc` BYTE-IDENTICAL, no return type changes, no concrete job edited
  (`products_job.go` untouched). Ambiguity removed at both loci.
- **N-2 (0.16 deletion rests on false premise + unowned published contracts)** — FIXED via
  the FIRST arm, ADAPTED (recorded deviation): `baseline_commission_percent: 0.16`
  (`auth_adapter.go:42-48`) is provider-catalog METADATA (published contract wiki/OpenAPI/
  SDK) with NO pricing call site — it stays INTACT; all false "fallback silencioso /
  call site / miss cai em 0.16" language DELETED per R-25 (falsehood is deleted, not
  hedged). Loci fixed: M-07 F-01 Brief + Ownership; M-07 milestone Mission/Outcome/Why/
  Ownership/Risks/Done-Means; IC-01:106 + Must-Not-Decide:177; mission.md ADR-09 (:180
  amendment recorded — prior disjunction "morre ou vira row" dropped as false-premised);
  mission matrix M-07 cell; M-01 milestone cross-reference (:47-50).
  **Adaptation rationale (recorded, not a downgrade):** the arm's literal "re-express as a
  `channel_fees` row origem='config'" would contradict the ratified F-10 owner-split
  (config step = COMPOSITION at the consumer M-07 F-01 via `pricingtariffdefaults`, never a
  ledger row — IC-01 §resolution). The baseline's semantics (16%) already live in the
  config degrau (`pricing_tariff_defaults` 16.00). Both arms' shared core — metadata key
  intact, no unowned contract deletion, false language gone — is satisfied.
- **N-3 (EARS-3 unreachable)** — FIXED via the FIRST arm: EARS-3 re-scoped to the branch
  that exists — degrau-4 STORE ERROR ⇒ typed error, never a constant (`Resolve` of
  tariffdefaults is total; materialize-on-read guarantees the row). Truth table stays 4
  fixtures with the 4th = "store de defaults falhando → erro TIPADO". This also formally
  registers the refutation of `p5-reconciliation-r01.md:78-79` ("EARS-3 is now reachable")
  — that claim was WRONG; r01 reconciliation stands as historical record, superseded here.
- **N-4 (dedupe key silently narrowed from ADR-11)** — FIXED via the SECOND arm (amendment
  recorded): mission.md ADR-11 now records the narrowing to
  `UNIQUE (provider, notification_id) WHERE notification_id IS NOT NULL` with rationale
  (`_id` not in the verified payload — fato #6 external-ml-api-facts; full tuple w/
  COALESCE would permanently blackhole legitimate re-notification of the same resource);
  IC-04:88-92 cross-references the amendment; M-08 milestone Done-Means + F-02 validation
  narrowed verbatim — replay WITH `_id` → zero new row; replay WITHOUT `_id` → extra rows,
  zero-duplicate proven at the DOMAIN level (IngestOrder idempotence, assertion on effect
  not on inbox count).
  **Arm-1 rejection rationale (recorded):** the NULL-safe full-tuple index would dedupe
  FOREVER two legitimately distinct notifications for the same (provider, topic, resource)
  once the first reaches a terminal status — a missed-update defect worse than the
  duplicate rows it prevents.

### Advisory (all applied)

- **N-5** — `fato #3` → `fato #13` at M-01 F-02:25 + :44; full sweep of every "fato #N"
  citation in manifested files now names its source file (9 additional loci:
  M-01 milestone + F-01 → `external-ml-api-facts.md`; M-08 F-01 → `external-ml-api-facts.md`;
  M-07 milestone + F-01 + F-02, M-06 F-03, M-05 F-03, M-09 F-01 + F-02 →
  `p5-prerequisites.md`).
- **N-6** — `0086-0088` → `0086-0089` at IC-01:30 and IC-02:29 (uncited siblings of A-1).
- **N-7** — M-02 F-02 Ownership ellipsis replaced with the canonical
  `internal/modules/channelfees/` + `internal/modules/divergences/` paths.
- **N-8** — mission matrix gained the `sync/application/` package-split rule row (M-02 F-03
  owns `scheduler.go` + cursor-contract helper; M-09 additive-only `health_*` +
  `sync/transport/**`), same idiom as the listings additive lock.
- **N-9** — IC-02:134 rewritten to `filter.divergentes=true` naming the `filter.` prefix
  idiom with the `query.go:31` + `filter.go:9` anchors and the silent-ignore hazard;
  IC-07 Operations row + Error Matrix aligned; M-05 milestone:26 residual bare
  `divergentes=true` also fixed.
- **N-10** — the 0.16 must-fail died with N-2 (nothing to kill); allowlist −2 (C/D) is now
  its OWN criterion with its OWN instrument (the M-02 F-04 guard test; must-fail =
  reintroduce an ML read-time call in pricing → allowlist names the site) at M-07 F-01
  Validation + milestone Done-Means — no criterion left without an instrument.

## Effect on prior artifacts

- `p5-reconciliation-r01.md:78-79` claim "EARS-3 is now reachable" REFUTED by N-3;
  superseded by this round (r01 file untouched — historical record).
- `p5-input-r02.sha256` INVALIDATED by the fold edits. New manifest: `p5-input-r03.sha256`.
- mission.md ADR-09 and ADR-11 now carry recorded amendments (N-2, N-4) — briefs cite them.

## Next

1. Freeze `p5-input-r03.sha256` over the same 46-file set (post-fold content).
2. Re-dispatch cold decomposition auditor (r03) on the frozen manifest; persist verbatim
   IMMEDIATELY on return.
3. Advance to P6 only on r03 PASS. Sol P5 retroactive touchpoint remains MANDATORY before
   `status: planned` (≥ 2026-08-05).
