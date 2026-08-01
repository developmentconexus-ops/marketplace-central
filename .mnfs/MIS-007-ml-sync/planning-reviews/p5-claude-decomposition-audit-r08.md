# P5 Decomposition Audit — round 08 — VERBATIM

```yaml
type: planning-review-verbatim
phase: P5
round: 08
created: 2026-08-01
auditor: cold Claude Opus crew (task af1ee8372fd436a93; operator-ratified waiver — Sol P5
  retroactive mandatory before status: planned)
input_manifest: p5-input-r08.sha256 (digest e74150e7d2ee1c482f77f45a8b44461e70f05931ee2c9e408b3d3cac15e38679;
  auditor recomputed self-digest MATCH + 46/46 OK)
transport: task output-file 0 bytes (fragile transport recurrence); recovered from session
  transcript JSONL, longest task-notification candidate (25,056 chars envelope → 24,223
  chars result), persisted SAME turn the notification arrived
verdict: PASS (zero blocking; 5 advisory F-r08-1..5)
---
```

# P5 Decomposition Audit — round 08 (MIS-007-ml-sync)

**Auditor:** cold Claude Opus crew, independent
**Date:** 2026-08-01
**Input manifest:** `planning-reviews/p5-input-r08.sha256` — self-digest recomputed = `e74150e7d2ee1c482f77f45a8b44461e70f05931ee2c9e408b3d3cac15e38679` **MATCH**; `sha256sum -c` → **46/46 OK**, zero mismatches. All 46 manifested files read; `planning-reviews/*` read only as PART A context.
**Prerequisite verification base:** repo working tree at `C:\Users\leandro.theodoro\Documents\marketplace-central`, branch `main`.

---

## PART A — closure verification of r07

### F-r07-1 (advisory) — missing M-09⤳M-02 soft edge — **CONFIRMED CLOSED**

DAG block, `mission.md:273-274`:
> `M-09 (⤳ M-04/M-06: entidades acendem sozinhas; ⤳ M-02: last_incremental_at real — pré-fix NULL uniforme é honesto e o gate passa; edge dura → M-08 acima)`

Edge table row, `mission.md:293`:
> `| M-09⤳M-02 | `last_incremental_at` REAL exige fix `incremental` (M-02 F-03, IC-05) — pré-fix, NULL uniforme é honesto e o gate do M-09 PASSA; lane A paraleliza sem ordem de close (F-r07-1) |`

Both loci carry IC-05's own qualifier. Closed without contradiction.

### F-r07-2 (advisory) — matrix false universal on root.go — **CONFIRMED CLOSED at both declared loci**

M-03 root.go cell, `mission.md:301`:
> `edita região orders existente `:576-601` in-place (troca readers A/B `:591-592` por readers de banco, deleções inclusas — mesma classe de exceção do M-07; F-r07-2), hub arbitra`

M-07 root.go cell, `mission.md:305`:
> `edita região pricing existente `:828-858` + remoção de imports tarifflive/tariffcomposite (`root.go:99,101` — F-r04-5); hub arbitra (região-edit, uma das DUAS exceções — a outra é M-03 região orders; F-r07-2)`

The matrix the hub adjudicates root.go collisions with no longer carries the false universal. **Two same-class residuals survive in milestone bodies** — the fold touched only the matrix. Recorded as F-r08-1 and F-r08-2 (advisory, same severity r07 assigned the parent finding).

### F-r07-3 (blocking) — `{null,0,0}` rendered as a CONFIGURATION verdict — **CONFIRMED CLOSED**

IC-05, `research/sync-health-interface-contract.md:106-111`:
> `M-09 pode fechar antes do M-08 — FE discrimina pelo `last_notification_at === null` e renderiza o FATO observado: "nenhuma notificação recebida". Rótulo de VEREDITO de configuração ("webhook não configurado") é PROIBIDO: pós-M-08, instalação CONFIGURADA com inbox vazio (registro feito e nenhum evento ainda; janela quieta; worker travado) produz o MESMO `{null,0,0}` — o payload não distingue configuração de silêncio, e a tela não pode afirmar o que o payload não carrega (honest-unknown; auditoria P5 r07 F-r07-3).`

M-09/F-02 Brief discriminator, `:34-36`:
> `discriminador de estado inicial = `last_notification_at === null` → subseção mostra o FATO: "nenhuma notificação recebida" (NUNCA veredito de configuração — payload `{null,0,0}` é idêntico p/ não-configurado e configurado-silencioso; IC-05 pina, F-r07-3)`

Third EARS, `:44-47`:
> `While webhook block está no estado canônico inicial (last_notification_at null), when card renderiza, the subseção shall dizer "nenhuma notificação recebida" — o fato, nunca o veredito "não configurado" (F-r07-3); não esconder, não inventar — pending/dropped 0 do estado inicial NÃO renderizam como atividade.`

New Inputs/Outputs, `:76-77`: `bloco webhook — estado inicial `{null,0,0}` → FATO "nenhuma notificação recebida" (nunca veredito de configuração — F-r07-3)`.

**Live-artifact sweep** (case-insensitive `não configurado|nao configurado`, mission root, `planning-reviews/` excluded): exactly **one** hit — `M-09-sync-observability/F-02-.../feature.md:46`, inside the EARS prohibition clause. The configuration-verdict label survives nowhere as a rendered string. The activity-fact label is byte-consistent across IC-05, Brief, EARS and Inputs/Outputs.

### F-r07-4 (advisory) — keep-absent range bisects the statement — **CONFIRMED CLOSED, verified against the repo**

`research/codebase-ingest-side.md:96-97`:
> `Keep-absent (ADR-04): `writer.go:104-112` (doc comment em `:97-103` — range medido, P5 r07 F-r07-4)`

Repo check, `apps/server_core/internal/modules/internal_read/adapters/mirror/writer.go`: doc comment occupies `:97-103`; `const keepAbsentSQL = ` at `:104`; statement terminates `AND codigo_produto <> ALL($2::text[])`` at `:112`. `upsertSQL` is `:74-95`. Corrected range is **byte-exact**.

### F-r07-5 (advisory) — `845-850` under-range on a DELETION instruction — **CONFIRMED CLOSED, verified against the repo**

`M-07-pricing-fee-read/milestone.md:33`: `Simulação hoje faz GETs vivos no ML por cotação (root.go:845-851)`
`M-07-pricing-fee-read/milestone.md:70-72`:
> `a deleção muda a fiação root.go:845-851 inteira (o `}` de fechamento em `:851` incluso; `var tariffResolver` em `:844` SOBREVIVE — F-r07-5, coerente com F-01)`

Repo check, `apps/server_core/internal/composition/root.go`: `:844 var tariffResolver pricingports.TariffResolver = calcTariffResolver`; `:845 if feeReader, ferr := marketplaceCapabilities.FeeQuoteReader(...)`; `:851 }`; `:852 batchOrch.WithTariffResolver(tariffResolver)`. Both the deleted block and the surviving declaration are exact. Zero live occurrences of `845-850` remain.

### F-r07-6 (advisory) — M-09/F-02 missing `## Inputs/Outputs` — **CONFIRMED CLOSED**

`M-09-.../F-02-integracoes-health-section/feature.md:69-78` carries the section, self-labelled `(Seção adicionada — auditoria P5 r07 F-r07-6.)`, binding the input to IC-05 `§entities[]` + the webhook block and enumerating four render states as output. Section sweep confirms **4/4 FE briefs** now carry both `## Inputs/Outputs` and `## Interaction Model`.

### Observation (stale blocker, M-02/F-01 Handoff) — **CONFIRMED DISCHARGED**

`M-02-sync-core-seam/F-01-core-ddl/feature.md:89-90`:
> `Blockers or open decisions: none — p5-prerequisites §2 COMPLETO no repo (`buyer_fiscal_reader.go` + DTO enumerados; blocker descarregado, observação r07).`

### Reopened-closure sweep (r01–r06) — **CLEAN, zero reopened**

Spot-checks all intact at their loci: F-r06-1 commission-scoped rejection `IC-01:142-147` + freight bullet `:106-111`; F-r06-2 `detail` in the `ResolveListingFees` tuple `IC-01:71-77` with the "seria insatisfazível contra o port" rationale preserved; F-r06-5 `root.go:845-851` in `M-07/F-01:35-36`; F-r06-6 `## Inputs/Outputs` present in `M-04/F-02` and `M-08/F-02:70-81`; F-r05-1 all-rows scan + ERP sentinel `M-09/F-01:36-41`; F-r05-2 three-term formula `M-06/F-03:26-33` + Constraints `:63-66`; F-r05-3 widened listings grant `mission.md:319-323`; F-r04-2 reference-injection `M-08/F-02:59-61`; F-r04-4 installations-port-not-new-store `M-08/F-02:44-49`; F-r04-5 hub-resolved import block `mission.md:311-313` + `M-07/F-01:92-93`; F-r04-6 write-DAG enumeration `mission.md:330-337`; F-r03-P2 ADR-14 ≤1-COMMIT `mission.md:260-261` + `:314-316`; F-r02-N2/N3/N4/N8; GREATEST `last_success_at` `IC-05:95-101` → `M-09/F-01:60-65` → `M-09/F-02:104`; F-8 origem vocabulary `IC-01:65-67`.

---

## PART B — fresh full audit

### Check 1 — DAG edge completeness and justification — **1 advisory**

Every hard edge in the table names a forcing artifact; every milestone's own `## Dependencies` reconciles with the table (M-01/M-02 roots; M-03←M-01,M-02; M-04←M-01,M-02; M-05←M-04,M-02; M-06←M-03,M-02 + soft M-05; M-07←M-02 + soft M-05; M-08←M-06,M-09; M-09 hard-free). Hard graph is acyclic; the M-08→M-09 hard edge plus M-09's soft "M-08 (bloco webhook real)" is a legitimate soft/hard pair already ordered by lanes (A before D) — no missing-edge finding. One ASCII/table disagreement recorded below.

**F-r08-5 · advisory · DAG edge completeness and justification**
Excerpt, `mission.md:268-270`:
> ```
> M-01 ──┬─→ M-03 ──┬─→ M-06 ──→ M-08
> M-02 ──┤          │             ↑
>        ├─→ M-04 ──┴─→ M-05     M-09 (porta WebhookStatsReader IC-05)
> ```
Locus: `mission.md:268` col 18 (`┬`), `:269` col 18 (`│`), `:270` col 18 (`┴`).
Offending token: the vertical connector `┬`/`│`/`┴` at column 18.
Defect: the glyphs form a connector between the M-03 line and the `M-04 ──…─→ M-05` line, drawing a cross-lane edge that the authoritative edge table (`:277-293`, "edge = artefato que força") declares in neither direction — there is no `M-05→M-03` row and no `M-06→M-04` row, and `M-05/milestone.md` Dependencies lists only `M-04`, `M-02`, `(M-01 transitiva via M-04)`. The diagram is ambiguous between two readings, both unbacked.
Yes-if: the vertical is removed so M-03's branch terminates at M-06 and M-05 is fed only by M-04 (matching the table and `M-05/milestone.md`), **or** an artifact is named that forces the drawn edge and a table row is added for it.

### Check 2 — canonical six-axis disjointness, seam locks, migration ranges — **2 advisory**

Migration ranges verified disjoint and against the repo: tip is `0085_erp_import_products_sellable_fields.sql`; M-02 `0086-0089`, M-04 `0090-0092`, M-08 `0093`, M-06 `0094-0095` (reserve, gap permitted) — each milestone's `Migration block:` line matches its matrix cell exactly; M-01/M-03/M-05/M-07/M-09 declare none. Package/table/OpenAPI/FE axes are disjoint or carry a recorded additive lock (`orders/**` lane inheritance `:317-318`; listings additive grant `:319-323`; `sync/application/` additive lock `:324-327`; import block hub-resolved `:311-313`). Two milestone bodies contradict the corrected matrix.

**F-r08-1 · advisory · six-axis disjointness / seam locks**
Excerpt, `M-07-pricing-fee-read/milestone.md:55-57`:
> `região pricing do root.go `:828-858` (ÚNICA exceção à regra linha-nova — edita bloco existente; hub arbitra o merge, R-7)`
Locus: `M-07-pricing-fee-read/milestone.md:56`.
Offending token: `ÚNICA`.
Defect: surviving instance of the exact false universal F-r07-2 corrected. `mission.md:305` now states `uma das DUAS exceções — a outra é M-03 região orders`; `mission.md:301` states M-03 edits its orders region in-place. The r07 residual sweep was case-sensitive on `única exceção` and could not see this uppercase form.
Yes-if: the token is restated as one of two region-edit exceptions naming M-03's orders region `:576-601`, matching `mission.md:301`/`:305`.

**F-r08-2 · advisory · six-axis disjointness / seam locks**
Excerpt, `M-03-orders-shipment-persist/milestone.md:55-57`:
> `1 linha ancorada root.go (região orders `:576-601` — edita SÓ dentro dela), entradas A/B da allowlist (remoção)`
Locus: `M-03-orders-shipment-persist/milestone.md:56`.
Offending token: `1 linha ancorada`.
Defect: same false universal as F-r08-1, opposite artifact. `M-03/F-03:27-28` swaps **two** reader lines (`root.go:591-592`), rewrites `:597-599`, and deletes the live call sites; `mission.md:301` now says `edita região orders existente :576-601 in-place (troca readers A/B :591-592 por readers de banco, deleções inclusas)`. A milestone orchestrator reading only the milestone body would grant M-03 a single-line root.go footprint it does not have.
Yes-if: the Ownership cell restates the real shape verbatim from `mission.md:301`.

### Check 3 — feature-level write-DAG — **PASS**

Every cross-feature write-set overlap carries a serial edge or a recorded grant, enumerated at `mission.md:330-337`: allowlist file (owner M-02/F-04) written by M-03/F-03 and M-07/F-01, serialized by lane B → lane C; `listings/**` extension by M-05/F-01–F-03 gated on M-04 close under the widened additive lock; `orders/**` by M-06 gated on M-03 close (lane inheritance, never simultaneous); root.go region edits arbitrated by the hub; import block ownerless and hub-resolved. The M-08 ↔ M-09 root.go seam is disjoint by construction — `M-08/F-02:57-59` writes `WithWebhookStatsReader(...)` **from M-08's own anchored region**, with `código/construção do M-09 NÃO é editado`. No residual overlap without an edge or grant.

### Check 4 — contract satisfiability — **PASS**

All IC canonical arithmetic re-derived and correct: IC-01 camada 2 `12.5% × 79.90 + 6.00 = 15.9875 → 15.99` (`:125-132`); IC-01 camada 3 `8.11 × 3 = 24.33` (`:134-139`); IC-03 `239.70 − 48.66 − 22.90 − 95.10 = 73.04`, `73.04/239.70 = 30.47%`; IC-01 tolerance `R$0.01` (`:186`) consistent with the M-06/F-03 divergence formula. Brief promises are satisfiable against the ratified set: M-07/F-01's camada-2 arm reads `detail.percentage_fee/fixed_fee`, which IC-01 `:71-77` guarantees in the `ResolveListingFees` tuple; M-06/F-03's three-term expected-fee formula consumes exactly the keys IC-01's camada-2 5-key `detail` tuple provides; M-07/F-02's DTO vocabulary `api_listing_prices | config` is a per-surface narrowing of IC-01 `:67`'s four-value mission vocabulary (order/shipment origins cannot arise in a pricing simulation) — narrowing, not contradiction; M-08/F-02's `WebhookStatsReader` injection is satisfiable only under IC-05 `§Seam`'s reference/pointer rule, which the brief cites and re-asserts.

### Check 5 — prerequisite existence — **1 advisory**

Every load-bearing anchor relied on by a brief was verified against the working tree, all exact unless noted:
`pricing/ports/tariff_resolver.go` — `TariffRequest` `:14-24`, `TariffResolver` `:30-32` ✔ | `root.go:837` `pricingtariffdefaults.NewResolver` ✔ | `root.go:844/845/851/852/853-855` ✔ | `root.go:99,101` tarifflive/tariffcomposite imports ✔ | `root.go:259-272` `registerBatchRoutes` (8 patterns) ✔ | `root.go:589-601` orders wiring, `:591` shipment reader, `:592` buyer-fiscal reader, `:597-599` `NewEnrichServiceWithReaders` ✔ | `root.go:370-378` `AccessTokenResolver` ✔ | `capability_adapter.go:712` `doRawWithHeaders` ✔ | `installation_repo.go:81` `func (r *InstallationRepository) ListInstallations` with `external_account_id` in the SELECT ✔ | `auth_adapter.go:192` and `:261` `ProviderAccountID: normalizeAnyString(payload.UserID)` ✔ | `auth_flow_service.go:691` `ExternalAccountID: firstNonEmpty(...)` ✔ | `auth_adapter.go:47` `baseline_commission_percent: 0.16` (inside cited `:42-48`) ✔ | `calc_handler.go:461` `handlePutTariffDefaults`, `:38` `PUT /pricing/tariff-defaults` ✔ | `calc_repository.go:239` `GetTariffDefaults` materialize-on-read ✔ | `route_deadline.go:23-28` const block, 15s/120s ✔ | `order_bucket.go:48` `DeriveOrderBucket` ✔ | `writer.go:74-95` / `:97-103` / `:104-112` ✔ | `IntegracoesPage.tsx:508` `listIntegrationInstallations`, `:558-574` section order verbatim ✔ | SDK `index.ts:2144` `listSyncRuns`, `:2359` `runPricingSimulation`, `:2361` `runBatchSimulation`, `:2374` `pricingDecompose` ✔ | migrations tip `0085` ✔.

**F-r08-3 · advisory · prerequisite existence**
Excerpt, `M-03-orders-shipment-persist/milestone.md:31-32`:
> `Hoje `GET /orders/{id}` faz 3-4 GETs vivos no ML por pedido (`enrich_service.go:192`, readers em root.go:590-592)`
Locus: `M-03-orders-shipment-persist/milestone.md:32`.
Offending tokens: `enrich_service.go:192` and `root.go:590-592`.
Defect, two parts. (a) `enrich_service.go:192` lands inside the doc comment; `func (s EnrichService) EnrichOne(` is at `:194` and the body ends `:198` — `research/codebase-read-side.md:41` already carries the correct `:194-198`, so the milestone contradicts its own source note. (b) `root.go:590` is `ordersCostReader := newOrdersCostReaderAdapter(internalReadSvc, internalReadAvailable)` — the **internal** cost reader, not an ML live-read site; the ML readers are `:591` (shipment) and `:592` (buyer fiscal), which is what `mission.md:301` and `M-03/F-03:27-28` both name. Same class as F-r07-5: a range on a kill-list instruction that swallows a line which must survive.
Yes-if: restated as `enrich_service.go:194-198` and `root.go:591-592`, matching `research/codebase-read-side.md:25,41` and `mission.md:301`.

### Check 6 — propagation of approved ADR/IC values into briefs — **1 advisory**

Numbers, enums, key names and routes propagate correctly: migration numbers (all 9 milestones); `orders_v2` topic and the four official IP addresses (log-only) from `external-ml-api-facts.md` #8/#9 into IC-04 and M-08; 64KB cap and always-200 into `M-08/F-01`; status enum `received/processing/done/malformed/unmapped/dropped` and the ≥5 threshold into `M-08/F-02:72-74`; `{backfill, incremental, sweep}` phase vocabulary; `15s`/`120s` deadline classes; `/webhooks/{provider}` path with **no** SDK method; `GREATEST(full, incremental)`; origem enum; the `0.16` metadata carve-out.

**F-r08-4 · advisory · propagation of approved ADR/IC values into briefs**
Excerpt, `M-08-webhook-ingest/F-02-worker-callback/feature.md:78-79`:
> `Output 2: impl REAL de `WebhookStatsReader` — shape de saída BINDING em IC-05 §InboxHealth (`last_notification_at`, `pending`, `dropped_24h`), consumido por `GET /sync/health``

and `M-09-sync-observability/F-02-integracoes-health-section/feature.md:73`:
> `§InboxHealth: last_notification_at, pending, dropped_24h) — a spec não re-decide shape.`

Locus: `M-08/F-02/feature.md:78` and `M-09/F-02/feature.md:73`.
Offending token: `IC-05 §InboxHealth`.
Defect: IC-05 contains no section by that name. Its headings are `Boundary`, `Why This Contract Exists`, `Resources Or Entities`, `Seam M-09 ↔ M-08`, `Operations`, `Fields` (`Required Inputs` / `Required Outputs`), `Enums And Statuses`, `Error Cases`, `Error Matrix`, `Persistence Expectations`, `Canonical Examples`, `Database Shape`, `Seed Data`, `Timestamp And ID Semantics`, `Compatibility Rules`, `Route Namespace`, `Transport And Integration`, `Must Preserve`, `Must Not Decide In Feature Execution`, `Validation Impact`. The webhook block — including the `{null,0,0}` canonical state and the F-r07-3 prohibition — lives in **IC-05 §Fields → §Required Outputs** at `:103-111`. `InboxHealth` is the name of IC-04's Operations row (`research/webhook-inbox-interface-contract.md:42`), which itself defers with `ver IC-05`. Both citations were introduced by the F-r06-6 and F-r07-6 folds, so a section declared BINDING points at a handle that resolves nowhere. The three field names cited are correct, so the shape is still recoverable by reading IC-05 in full — hence advisory, not blocking.
Yes-if: both briefs cite the real heading (`IC-05 §Required Outputs — bloco webhook`), or IC-05 gains a heading with that name (contract edit, hub-serialized under ADR-14).

### Check 7 — required brief sections — **PASS**

Sweep over all 25 feature briefs: `## Negative Scenarios` **25/25**; `## Validation Expectations` **25/25**; `## Interaction Model` present in exactly the 4 UI briefs (`M-05/F-03`, `M-06/F-03`, `M-07/F-02`, `M-09/F-02`) and in no non-UI brief; `## Inputs/Outputs` present in 22/25. The three without it are correctly exempt as neither API nor data features and say so in their own text: `M-01/F-01` (resilience decorator over an unchanged `doRawWithHeaders` signature — no wire or DB shape), `M-02/F-04` (architectural guard test; allowlist is test data, `allowlist como dado no próprio teste`), `M-04/F-04` (composition/wiring only — `Sem contrato FE novo (rota de refresh já existe)`). Validation expectations are inspectable-proof throughout — must-fail guards (`M-02/F-04`, `M-07/F-01`), truth-table fixtures with named errors (`M-07/F-01:100-101`), effect-not-count assertions (`M-08/F-02:101-103`), negative controls outside `registerBatchRoutes` (`M-08/F-01:76`), and the F-r04-1 incremental-only badge fixture (`M-09/F-02:102-104`).

### Check 8 — no implementation planning, no new product scope — **PASS**

Briefs stay at contract/ownership level; implementation is consistently deferred (`Execução cria spec/plan/validation`; `técnica exata = spec`; `spec pina valor`). Scope discipline holds at the known pressure points: no staleness threshold in M-07/F-02 (`limiar não-ratificado não entra`); camada 1 implemented but documented dormant with no producer; `api_shipping_options` reserved in the CHECK with no producer; `/market/collections` explicitly deferred to MIS-008; no `LISTEN/NOTIFY`; the ML callback registration is gated on explicit operator authorization at `M-08/milestone.md:68-70` and `M-08/F-02:31-32`, consistent with the live-ML-writes rule.

---

**Verdict: PASS**

Zero blocking findings. Five advisory findings (F-r08-1 … F-r08-5) do not change the verdict. Note that F-r08-1 and F-r08-2 are same-class residuals of F-r07-2 that the r07 residual sweep could not see because it was case-sensitive; a case-insensitive re-sweep is the cheap guard against this class recurring.

---

## Files read

**Manifested (46/46), relative to `.mnfs/MIS-007-ml-sync/`**
`mission.md`
`M-01-ml-client-hardening/milestone.md`, `M-01-ml-client-hardening/F-01-resilience-decorator/feature.md`, `M-01-ml-client-hardening/F-02-items-multiget-raw-dto/feature.md`
`M-02-sync-core-seam/milestone.md`, `M-02-sync-core-seam/F-01-core-ddl/feature.md`, `M-02-sync-core-seam/F-02-fee-divergence-ports/feature.md`, `M-02-sync-core-seam/F-03-scheduler-incremental-cursor/feature.md`, `M-02-sync-core-seam/F-04-read-guard-allowlist/feature.md`
`M-03-orders-shipment-persist/milestone.md`, `M-03-orders-shipment-persist/F-01-ml-ingest-readers/feature.md`, `M-03-orders-shipment-persist/F-02-ingest-order-v1/feature.md`, `M-03-orders-shipment-persist/F-03-read-path-switch/feature.md`
`M-04-listings-backfill-ingest/milestone.md`, `M-04-listings-backfill-ingest/F-01-listings-ddl/feature.md`, `M-04-listings-backfill-ingest/F-02-mass-closure-replacement/feature.md`, `M-04-listings-backfill-ingest/F-03-backfill-cursor-ingest/feature.md`, `M-04-listings-backfill-ingest/F-04-scheduler-refresh-wiring/feature.md`
`M-05-listings-fees-divergence/milestone.md`, `M-05-listings-fees-divergence/F-01-camada2-fee-ingest/feature.md`, `M-05-listings-fees-divergence/F-02-stock-divergence/feature.md`, `M-05-listings-fees-divergence/F-03-anuncios-fe-contract/feature.md`
`M-06-orders-backfill-decomposition/milestone.md`, `M-06-orders-backfill-decomposition/F-01-backfill-incremental/feature.md`, `M-06-orders-backfill-decomposition/F-02-decomposition-camada3/feature.md`, `M-06-orders-backfill-decomposition/F-03-audit-fe-pedidos/feature.md`
`M-07-pricing-fee-read/milestone.md`, `M-07-pricing-fee-read/F-01-fee-read-resolver/feature.md`, `M-07-pricing-fee-read/F-02-precos-provenance-fe/feature.md`
`M-08-webhook-ingest/milestone.md`, `M-08-webhook-ingest/F-01-inbox-endpoint/feature.md`, `M-08-webhook-ingest/F-02-worker-callback/feature.md`
`M-09-sync-observability/milestone.md`, `M-09-sync-observability/F-01-sync-health-endpoint/feature.md`, `M-09-sync-observability/F-02-integracoes-health-section/feature.md`
`research/channel-fees-interface-contract.md`, `research/codebase-ingest-side.md`, `research/codebase-read-side.md`, `research/divergences-interface-contract.md`, `research/external-ml-api-facts.md`, `research/listings-sync-interface-contract.md`, `research/orders-persistence-interface-contract.md`, `research/p5-prerequisites.md`, `research/sync-health-interface-contract.md`, `research/sync-ingest-ports-interface-contract.md`, `research/webhook-inbox-interface-contract.md`

**PART A context only (not audited)**
`planning-reviews/p5-input-r08.sha256`, `planning-reviews/p5-claude-decomposition-audit-r07.md`, `planning-reviews/p5-reconciliation-r07.md`

**Repo, for check-5 prerequisite verification (relative to repo root)**
`apps/server_core/internal/composition/root.go`
`apps/server_core/internal/modules/pricing/ports/tariff_resolver.go`
`apps/server_core/internal/modules/pricing/transport/calc_handler.go`
`apps/server_core/internal/modules/pricing/adapters/postgres/calc_repository.go`
`apps/server_core/internal/modules/orders/application/enrich_service.go`
`apps/server_core/internal/modules/orders/domain/order_bucket.go`
`apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go`
`apps/server_core/internal/modules/integrations/adapters/mercadolivre/auth_adapter.go`
`apps/server_core/internal/modules/integrations/application/auth_flow_service.go`
`apps/server_core/internal/modules/integrations/adapters/postgres/installation_repo.go`
`apps/server_core/internal/modules/internal_read/adapters/mirror/writer.go`
`apps/server_core/internal/platform/httpx/route_deadline.go`
`apps/server_core/migrations/` (listing)
`apps/web/src/pages/integracoes/IntegracoesPage.tsx`
`packages/sdk-runtime/src/index.ts`