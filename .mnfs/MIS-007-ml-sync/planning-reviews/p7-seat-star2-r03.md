# P7 seat único — ★2 Consistency — round r03 (verbatim)

```yaml
round: r03
seat: 1/1 (rodada enxuta dirigida pelo operador — só o critério que falhou no r02)
scope: "★2 Consistency (re-review escopada)"
model: claude sonnet (cold mission-reviewer)
manifest: NENHUM (desvio registrado — shell wedge B-8; input = árvore em disco,
  imutável durante a rodada, lista prescrita no prompt)
transport: Agent tool síncrono (result in-context, completo)
verdict: "★2 PASS"
persisted: 2026-08-01
operator_ruling: "Deixe assim e mande sonnet subagents conferir e somente um revisar
  gastou muitos creditos" (2026-08-01)
```

Resultado verbatim do assento (abaixo, sem edição):

---

Consistent everywhere: `(tenant_id, provider, provider_shipment_id)`. No divergence found. Based on this thorough sweep of the mission tree, both previously failed loci are cleanly repaired and no new divergence was introduced or discovered elsewhere in the tree.

## ★2 Consistency — Round r03 (scoped re-review)

**Verdict: ★2 PASS**

### Findings

**★2-A — `listing_variations` PK tuple.** Now uniform (5-column, `installation_id` included) across all binding loci:
- `mission.md:229-231` (ADR-13): `(tenant_id, installation_id, provider, provider_listing_id, variation_id)`
- `research/listings-sync-interface-contract.md:65-70` (IC-07): same tuple, with reasoning tying it to parent `listings` PK (`0036_listings.sql:2-31`)
- `M-04-listings-backfill-ingest/F-01-listings-ddl/feature.md:31-32`: same tuple, cites "ADR-13/IC-07"

Rationale: grepped `listing_variations|variation_id` across the whole tree; the only remaining 4-column statements are inside excluded/historical files (`planning-reviews/p7-seat5-doublepass-r02.md`, `p7-claude-readiness-r02.md`) that record the pre-repair defect for the record — not binding artifacts.

**★2-B — Webhook route slug vs `provider_code`.** Now uniform (`mercado_livre`) at both loci named in the r02 finding:
- `research/webhook-inbox-interface-contract.md:31-34` (IC-04): "`provider` = slug `mercado_livre`, mesmo vocabulário provider_code das 4 superfícies ... `mercadolivre` é nome de PACKAGE Go, nunca slug de rota/coluna"
- `M-08-webhook-ingest/F-01-inbox-endpoint/feature.md:39-40`: "slug provider_code das 4 superfícies (`mercado_livre` ... não confundir com package Go `mercadolivre`)"

Cross-checked `sync-health-interface-contract.md` (IC-05) — it does not restate the slug token at all (inbox aggregates are read by field name only), so no third locus to diverge. `M-08-webhook-ingest/milestone.md`, `M-08-webhook-ingest/validation-contract.md`, `M-08-webhook-ingest/F-02-worker-callback/feature.md`, `M-09-sync-observability/F-01-sync-health-endpoint/feature.md`, `M-09-sync-observability/F-02-integracoes-health-section/feature.md`, `M-09-sync-observability/validation-contract.md` contain no slug restatement that could drift.

**Adversarial sweep beyond the two named loci (no new divergence found):**
- `listing_variations` grain/format used identically in `M-05-listings-fees-divergence/F-02-stock-divergence/feature.md:28,68` (variation grain matches IC-07 §Persistence Expectations).
- `order_shipments` PK `(tenant_id, provider, provider_shipment_id)` identical across `M-02-sync-core-seam/F-01-core-ddl/feature.md:30`, `M-02-sync-core-seam/validation-contract.md:49`, `M-03-orders-shipment-persist/F-02-ingest-order-v1/feature.md:31`, `M-03-orders-shipment-persist/validation-contract.md:62-63`.
- `channel_fees` vocab (`origem` CHECK `api_listing_prices|api_shipping_options|api_order|api_shipment|config`, `fee_kind`, `subject_id` formats) identical between `research/channel-fees-interface-contract.md` and its producers `M-05-listings-fees-divergence/F-01-camada2-fee-ingest`/`M-06-orders-backfill-decomposition/F-02-decomposition-camada3/feature.md:26-30`.
- `divergences` vocab (`kind ∈ {estoque, tarifa}`, `entity_type ∈ {listing, order_line}`, entity_id formats) identical between `research/divergences-interface-contract.md` and `M-05-listings-fees-divergence/F-02-stock-divergence/feature.md:30-31` and `M-06-orders-backfill-decomposition/F-03-audit-fe-pedidos/feature.md:31-32,75`.
- IC-05 `GetSyncHealth` payload field names (`entity, last_success_at, last_incremental_at, consecutive_failures, phase, last_error`; `webhook: {last_notification_at, pending, dropped_24h}`) reproduced verbatim in `M-09-sync-observability/F-01-sync-health-endpoint/feature.md:29-30` and `F-02-integracoes-health-section/feature.md:71-73`.
- Error Matrices checked for completeness against each IC's declared Operations/Error Cases (IC-01..IC-07) — each feature-returned status/case (409 refresh_in_progress, 200-always webhook, 500 internal-read) has a corresponding Error Matrix row; no undeclared feature-returned case found.
- List/collection orderings declared: `listing_variations`/`listings` sort (`p5-prerequisites.md:95`, unchanged), IC-02 `ListOpenByEntity` (`detected_at desc`), IC-05 `entities` (name asc) — all named, none silently unordered.

### Coverage

**Read line-by-line:**
- `mission.md`, `validation-contract.md` (mission root)
- `research/listings-sync-interface-contract.md`, `research/webhook-inbox-interface-contract.md`, `research/sync-health-interface-contract.md`, `research/channel-fees-interface-contract.md`, `research/divergences-interface-contract.md`, `research/orders-persistence-interface-contract.md`, `research/sync-ingest-ports-interface-contract.md`
- `M-04-listings-backfill-ingest/milestone.md`, `M-04-listings-backfill-ingest/F-01-listings-ddl/feature.md`, `M-04-listings-backfill-ingest/F-04-scheduler-refresh-wiring/feature.md`
- `M-08-webhook-ingest/milestone.md`, `M-08-webhook-ingest/validation-contract.md`, `M-08-webhook-ingest/F-01-inbox-endpoint/feature.md`, `M-08-webhook-ingest/F-02-worker-callback/feature.md`
- `M-09-sync-observability/F-01-sync-health-endpoint/feature.md`, `M-09-sync-observability/F-02-integracoes-health-section/feature.md`, `M-09-sync-observability/validation-contract.md`
- `M-02-sync-core-seam/validation-contract.md`
- `M-01-ml-client-hardening/validation-contract.md`
- `M-07-pricing-fee-read/validation-contract.md`
- `M-05-listings-fees-divergence/F-02-stock-divergence/feature.md`
- `M-06-orders-backfill-decomposition/F-02-decomposition-camada3/feature.md`
- `planning-reviews/p3-reconciliation-r01.md` (A-13 excerpt only, offset 30-44)

**Grep-swept (pattern matches inspected, full files not opened unless flagged above):**
- `mercado_livre|mercadolivre|mercado-livre` across entire mission tree
- `listing_variations|variation_id` across entire mission tree
- `kind.*tarifa|tarifa.*kind|AuditRealizedVsEstimated|entity_type.*order_line` across entire mission tree
- `provider_shipment_id|order_shipments` across `M-02-sync-core-seam/` and `M-03-orders-shipment-persist/`
- `mission\.md$` glob (confirmed single canonical mission.md at root; all `.sha256`/other hits are unrelated freeze manifests)

**Not opened (excluded per dispatch instructions):** all `planning-reviews/p7-*` review outputs, all `planning-reviews/p5-claude-decomposition-audit-*`, `p3-claude-candidate-r01.md`, `p3-opus-counterproposal-r01.md`, `sol-unavailable-p3-r01.md`, `p5-passes-r01.md`, all `*.sha256`, scratch files — grep hits inside disregarded.

**Not read this round (outside touched-file/adversarial-sweep scope, no grep hit suggesting drift):** `M-01/F-01`, `M-01/F-02` feature bodies; `M-02/F-01..F-03` feature bodies (grep-swept); `M-03/F-03` body (grep-swept); `M-05/F-01`, `F-03`, milestone, VC bodies (grep-swept); `M-06/F-01`, `F-03`, milestone, VC bodies (grep-swept); `M-07/F-01`, `F-02` bodies; `M-09/milestone.md` body.

No overall seven-★ verdict computed — per dispatch instructions this scoped ★2-only finding is folded by the dispatching session.
