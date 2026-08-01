# P5 Reconciliation — round 01 (MIS-007-ml-sync)

```yaml
type: planning-review-reconciliation
phase: P5
round: 01
created: 2026-07-31
auditor_artifact: p5-claude-decomposition-audit-r01.md
auditor: cold Claude Opus crew (operator-ratified waiver; Sol P5 retroactive mandatory before status: planned)
input_manifest: p5-input-r01.sha256 (digest 30d48ee15cacfe2650a025870c0afa596dcc5ac9a2793c9eaaef07fd48dd5d51)
verdict_received: NEEDS-REVISION (10 blocking F-1..F-10, 9 advisory A-1..A-9)
disposition: ALL findings ACCEPTED as valid; fold applied in full; re-audit on r02 manifest required
```

## Process defect (recorded)

The auditor's output returned in-context (Agent task `a8580ea15a44ddbba`) and was NOT
persisted at dispatch time. Recovered VERBATIM from the session transcript after context
compaction and persisted to `p5-claude-decomposition-audit-r01.md`. Rule going forward:
persist reviewer stdout to its artifact in the same turn the result arrives, before any
fold work. (Same class as "unwritten = didn't happen" — the audit nearly vanished.)

## Per-finding disposition (all ACCEPTED, no downgrades)

### Blocking

- **F-1 (nonexistent ownership paths — 9 loci)** — FIXED. Canonical paths everywhere:
  `internal/modules/connectors/adapters/mercado_livre/` (M-01 milestone + F-01 + F-02),
  `internal/modules/channelfees|divergences`, `modules/sync/application/scheduler.go`
  (M-02 milestone + F-03), `internal/modules/listings/**` (M-04 milestone:50, F-02:65,
  F-03:75), mission matrix M-01 cell (scoped-exclusivity wording, folds A-4).
- **F-2 (observed_at vs coletado_em)** — FIXED. All channel_fees provenance references
  renamed `coletado_em`: M-05 F-01 (brief, EARS, negative, validation), M-05 F-03 (DTO,
  EARS, milestone ×2), M-07 F-01 (brief, EARS, negative), M-07 F-02 (brief, EARS, shape),
  M-07 milestone (×2). Residual `observed_at` hits verified to be IC-02 divergences
  vocabulary (`expected_observed_at`/`observed_observed_at`) — distinct contract, correct.
- **F-3 (camada 2 percent vs amount)** — FIXED with IC-01 canonical winning: camada 2 =
  `value_type=amount` + `detail {percentage_fee, fixed_fee, financing_add_on_fee,
  price_used, listing_type_id}`. M-05 F-01 brief/EARS/validation, M-05 F-03 DTO
  (`tarifa {amount, detail, origem, coletado_em}`) + FE column, M-05 milestone. M-06 F-03
  audit formula re-pinned: esperado_unit = detail.percentage_fee × unit_price/100 +
  detail.fixed_fee; esperado_total = × quantity; NEVER the camada-2 `amount` (anchored to
  price_used, not order price).
- **F-4 (false sum-of-parts invariant)** — FIXED. Invariant now
  `liquido = receita_bruta − comissao_total − frete_seller − custo_produto` (canonical
  239.70−48.66−22.90−95.10=73.04) in M-06 F-02 EARS + Constraints + milestone Done Means;
  any absent part ⇒ liquido/margem ABSENT + incompleto[], never partial sum.
- **F-5 (installation_id param IC-05 forbids)** — FIXED. M-09 F-01: param dropped at all
  3 loci; tenant from ctx, tenant scan; negative scenario re-written (tenant without
  installation/sync → null-honest, not 404).
- **F-6 (webhook block shape stated 3 incompatible ways)** — FIXED. IC-05:74 de-hedged to
  the exact canonical initial state `{"last_notification_at":null,"pending":0,
  "dropped_24h":0}` (timestamp null = never observed; counters 0 = true empty count).
  M-09 F-01 brief/EARS/validation, M-09 F-02 (FE discriminator = last_notification_at
  null), M-09 milestone quote it verbatim.
- **F-7 (missing M-08→M-09 edge + cross-region wiring swap)** — FIXED. Hard edge added to
  mission DAG diagram + edge table (`porta WebhookStatsReader (IC-05) + injeção
  WithWebhookStatsReader da região ancorada do M-08`). Indirection specified: M-09 F-01
  delivers default impl + `WithWebhookStatsReader(...)` setter + compile-time assert
  (catalog-503 lesson); M-08 F-02/milestone call it from M-08's OWN anchored region —
  M-09's construction line never edited. IC-05 gained a binding seam section (port
  signature + default + swap rule) — also closes the previously dangling "assinatura
  IC-05" reference in the briefs. M-08 milestone Dependencies now lists M-09.
- **F-8 (camada-2 freight promised, no producer)** — FIXED via drop-and-reserve:
  mission.md M-05 headline drops `+ shipping_options` (frete = honest-unknown); M-07 F-02
  vocabulary = `api_listing_prices | config` for this mission; IC-01 notes
  `api_shipping_options` in the CHECK is an additive-future reserve with NO producer this
  mission.
- **F-9 (M-05 fallback path writes in ML adapter it doesn't own)** — FIXED via the
  audit's second yes-if arm: M-01 F-02 gains the named deliverable (verify multiget
  exposes sale_price with `?context=channel_marketplace`; if not, ALSO deliver dedicated
  `GET /items/{id}/prices` reader as a new file in the ML package). M-05 F-01 now consumes
  ONE ready source; its open-decision line closed.
- **F-10 (config fallback with two owners)** — FIXED per recommended arm: config step
  owned by M-07 F-01 only (composes `pricingtariffdefaults.NewResolver` as chain base).
  M-02 F-02 `ChannelFeeReader` narrowed to LEDGER-ONLY (camada 2→1→typed-absent, never
  config) in brief/EARS/Inputs/negative; M-02 milestone Done Means updated; IC-01
  §Enums/resolution amended with the owner split. M-07 F-01 EARS-3 (honest-absent) is now
  reachable.

### Advisory (all applied)

- **A-1** — mission.md:235 + IC-01 Database Shape → 0086-0089.
- **A-2** — writer.go anchor → `:74-105` in mission.md (M-02 F-02 Inputs already cited
  74-105 post-F-10 rewrite).
- **A-3** — `M-05→M-02` edge row added; M-08→M-06 forcing artifact renamed to "IngestOrder
  ESTENDIDO do M-06 (decomposição+camada 3) — não o v1 do M-03".
- **A-4** — folded into F-1 (M-01 matrix cell: exclusivity over EXISTING files; new files
  permitted downstream, M-03 F-01 named).
- **A-5** — `## Interaction Model` added to M-05 F-03, M-06 F-03, M-07 F-02 (ownership,
  refetch=same-payload, stale policy, no new routes).
- **A-6** — `## Inputs/Outputs` added to the two asymmetric briefs: M-04 F-01 (DDL
  verbatim IC-07) and M-03 F-03 (byte-identical comprador_fiscal + additive IC-03 shape).
- **A-7** — 7-day staleness rule REMOVED (unratified product surface); ⚠ limited to
  origem=config; coletado_em displayed without judgment; threshold = future decision.
- **A-8** — codebase-read-side.md:162 corrected 0081→0085.
- **A-9** — M-06 F-02 freight row pinned subject_type=`order`,
  subject_id=`<provider_order_id>` (frete é do pedido, não da linha).

## Effect on prior artifacts

- `p5-passes-r01.md` stands as historical record; its Pass 1 claim "Sem edit cruzado de
  arquivo" (M-08↔M-09) was REFUTED by F-7 and is superseded by the setter design; its
  Pass 3 clearance of M-05 F-01 path B was superseded by F-9.
- `p5-input-r01.sha256` INVALIDATED by the fold edits. New manifest: `p5-input-r02.sha256`.

## Next

1. Freeze `p5-input-r02.sha256` over the same 46-file set (post-fold content).
2. Re-dispatch cold decomposition auditor (r02) on the frozen manifest; persist verbatim
   IMMEDIATELY on return.
3. Advance to P6 only on r02 PASS. Sol P5 retroactive touchpoint remains MANDATORY before
   `status: planned` (≥ 2026-08-05).
