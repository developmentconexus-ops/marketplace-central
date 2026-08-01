# P5 Decomposition Audit — round 05 (VERBATIM)

```yaml
type: planning-review-audit
phase: P5
round: 05
created: 2026-07-31
auditor: cold Claude Opus crew (task aad8fe745b084ab00; operator-ratified waiver — Sol P5
  retroactive mandatory before status: planned)
input_manifest: p5-input-r05.sha256 (digest eefe38a5b903af0c54baaaac58436a32b68eca033bdaefdcfee2d19e5e77c22e)
provenance: task output-file returned 0 bytes (fifth instance of the transport fragility);
  text below persisted VERBATIM from the task-notification (recovered from the session
  transcript JSONL) in the same turn it arrived; transport HTML entities decoded.
verdict: NEEDS-REVISION (F-r04-1..6 all CONFIRMED CLOSED; 2 new blocking F-r05-1..2,
  4 advisory F-r05-3..6)
```

---

# P5 Decomposition Audit — MIS-007-ml-sync — Round 05

## 1. Header

### Manifest verification

- **Manifest self-digest**: `eefe38a5b903af0c54baaaac58436a32b68eca033bdaefdcfee2d19e5e77c22e` — **MATCHES** expected.
- **`sha256sum -c planning-reviews/p5-input-r05.sha256`** (Git Bash, from mission root): **46 OK / 0 FAILED** (46 lines).
- Manifest integrity confirmed; audit proceeded.

### Files read (46/46 manifested, no sampling)

Mission-root-relative:

1. `M-01-ml-client-hardening/F-01-resilience-decorator/feature.md`
2. `M-01-ml-client-hardening/F-02-items-multiget-raw-dto/feature.md`
3. `M-01-ml-client-hardening/milestone.md`
4. `M-02-sync-core-seam/F-01-core-ddl/feature.md`
5. `M-02-sync-core-seam/F-02-fee-divergence-ports/feature.md`
6. `M-02-sync-core-seam/F-03-scheduler-incremental-cursor/feature.md`
7. `M-02-sync-core-seam/F-04-read-guard-allowlist/feature.md`
8. `M-02-sync-core-seam/milestone.md`
9. `M-03-orders-shipment-persist/F-01-ml-ingest-readers/feature.md`
10. `M-03-orders-shipment-persist/F-02-ingest-order-v1/feature.md`
11. `M-03-orders-shipment-persist/F-03-read-path-switch/feature.md`
12. `M-03-orders-shipment-persist/milestone.md`
13. `M-04-listings-backfill-ingest/F-01-listings-ddl/feature.md`
14. `M-04-listings-backfill-ingest/F-02-mass-closure-replacement/feature.md`
15. `M-04-listings-backfill-ingest/F-03-backfill-cursor-ingest/feature.md`
16. `M-04-listings-backfill-ingest/F-04-scheduler-refresh-wiring/feature.md`
17. `M-04-listings-backfill-ingest/milestone.md`
18. `M-05-listings-fees-divergence/F-01-camada2-fee-ingest/feature.md`
19. `M-05-listings-fees-divergence/F-02-stock-divergence/feature.md`
20. `M-05-listings-fees-divergence/F-03-anuncios-fe-contract/feature.md`
21. `M-05-listings-fees-divergence/milestone.md`
22. `M-06-orders-backfill-decomposition/F-01-backfill-incremental/feature.md`
23. `M-06-orders-backfill-decomposition/F-02-decomposition-camada3/feature.md`
24. `M-06-orders-backfill-decomposition/F-03-audit-fe-pedidos/feature.md`
25. `M-06-orders-backfill-decomposition/milestone.md`
26. `M-07-pricing-fee-read/F-01-fee-read-resolver/feature.md`
27. `M-07-pricing-fee-read/F-02-precos-provenance-fe/feature.md`
28. `M-07-pricing-fee-read/milestone.md`
29. `M-08-webhook-ingest/F-01-inbox-endpoint/feature.md`
30. `M-08-webhook-ingest/F-02-worker-callback/feature.md`
31. `M-08-webhook-ingest/milestone.md`
32. `M-09-sync-observability/F-01-sync-health-endpoint/feature.md`
33. `M-09-sync-observability/F-02-integracoes-health-section/feature.md`
34. `M-09-sync-observability/milestone.md`
35. `mission.md`
36. `research/channel-fees-interface-contract.md`
37. `research/codebase-ingest-side.md`
38. `research/codebase-read-side.md`
39. `research/divergences-interface-contract.md`
40. `research/external-ml-api-facts.md`
41. `research/listings-sync-interface-contract.md`
42. `research/orders-persistence-interface-contract.md`
43. `research/p5-prerequisites.md`
44. `research/sync-health-interface-contract.md`
45. `research/sync-ingest-ports-interface-contract.md`
46. `research/webhook-inbox-interface-contract.md`

### Repo anchors independently verified (Read/Grep, outside the manifest)

All paths under `C:\Users\leandro.theodoro\Documents\marketplace-central`:

| Anchor claimed by plan | Verified |
| --- | --- |
| `apps/server_core/internal/modules/sync/adapters/postgres/sync_state_repo.go:62,74-79` COALESCE-per-type write | EXACT |
| `apps/server_core/internal/composition/root.go:99` `pricingtariffcomposite` import, `:101` `pricingtarifflive` import | EXACT |
| `root.go:856` `pricingHandler = pricingHandler.WithCalc(calcSvc)` (value-builder idiom, cited as `:854-857`) | EXACT (line 856, inside cited range) |
| `root.go:837` `pricingtariffdefaults.NewResolver`; `:849-850` live/composite; `:852` `batchOrch.WithTariffResolver`; `:853-855` calcSvc chain; `:858` `pricingHandler.Register` | EXACT |
| `root.go:259-272` `registerBatchRoutes` list | EXACT |
| `root.go:370-378` `AccessTokenResolver` wiring | EXACT |
| `root.go:591` `ordersShipmentReader` (sítio A) | EXACT |
| `internal/modules/integrations/adapters/postgres/installation_repo.go:81` `ListInstallations` selecting `external_account_id` | EXACT |
| `internal/modules/integrations/adapters/mercadolivre/auth_adapter.go:192` and `:261` `ProviderAccountID: normalizeAnyString(payload.UserID)` | EXACT (both) |
| `internal/modules/integrations/application/auth_flow_service.go:691` `ExternalAccountID: firstNonEmpty(...)` | EXACT |
| `migrations/0016_integrations_foundation.sql:33` `external_account_id` column; `0017:30-32` partial unique index | EXISTS |
| `internal/modules/orders/domain/order_bucket.go:48` `func DeriveOrderBucket(...)` | EXACT |
| `internal/modules/pricing/ports/tariff_resolver.go:14-24` `TariffRequest`, `:30-32` `TariffResolver` | EXACT |
| `internal/modules/pricing/adapters/postgres/calc_repository.go:239` `GetTariffDefaults` (materialize-on-read) | EXACT |
| `internal/modules/pricing/transport/calc_handler.go:461` `handlePutTariffDefaults` | EXACT |
| `internal/modules/internal_read/adapters/mirror/writer.go:74-105` upsert + keep-absent idiom | EXACT |
| `internal/modules/connectors/adapters/mercado_livre/capability_adapter.go:712` `doRawWithHeaders` (ADR-02 seam) | EXACT |
| `internal/modules/connectors/adapters/mercado_livre/buyer_fiscal_reader.go` | EXISTS |
| `internal/modules/sync/application/scheduler.go:160` `RecordSuccess(..., false)` hardcoded | EXACT |
| `internal/modules/sync/composition/scheduler.go:11` `const InstallationScopeERP = "erp"`; `products_job.go:50` uses it | EXACT |
| `internal/modules/listings/domain/filter.go:9` `FilterKeys`; `listings/transport/query.go:31` `filter.` prefix | EXACT |
| `internal/modules/integrations/transport/run_read_handler.go:34` `GET /sync/runs` | EXACT |
| `internal/modules/market/ports/product_identity_reader.go:14` `LinkedListings` | EXISTS |
| `internal/modules/product_links/application/import_service.go:84` `AbsorbProviderSnapshots`; `listings/adapters/connectors/source.go:18,44` | EXACT |
| `packages/sdk-runtime/src/index.ts:2144` `listSyncRuns`, `:2288` `listOrders`, `:2290` `getOrder` | EXACT |
| `apps/web/src/pages/integracoes/IntegracoesPage.tsx:558-574` mount region (file is 574 lines) | EXACT |
| `apps/web/src/pages/pedidos/PedidoDrawer.tsx:355-368` `CompradorFiscalSection` | EXACT |
| Migration tip = `0085_erp_import_products_sellable_fields.sql` (72 `.sql` files) — 0086+ free | CONFIRMED |
| Migration runner `internal/platform/migrate/runner.go:33-84` — per-filename `schema_migrations` tracking, lexicographic apply | CONFIRMED (out-of-lane-order arrival of 0093 after 0094-0095 is safe; no finding) |
| `sync/transport/` package does not exist (M-09 creates it); `sync/application/scheduler.go` exists (M-02 owner) | CONFIRMED |

Zero anchor named by the plan was found false or missing.

---

## 2. Verification of r04 findings

### F-r04-1 (BLOCKING in r04) — **CLOSED, fix real and complete**

- `research/sync-health-interface-contract.md:94-100` now defines `last_success_at` as `GREATEST(last_full_sync_at, last_incremental_at)`, "NULL só quando ambos NULL", cites the writer at `sync_state_repo.go:62,74-79`, and mandates the negative fixture ("entidade com full velho + incremental recente ⇒ `last_success_at` = o incremental, badge fresco").
- `M-09-sync-observability/F-01-sync-health-endpoint/feature.md:54-61` (Inputs) carries GREATEST verbatim with the explicit "NUNCA só o full, que congela em entidade incremental-only". The `= last_full_sync_at` equation is **absent** from the whole brief (grepped).
- `M-09/F-01:98-100` (Validation): "Fixture negativa F-r04-1: entidade com `last_full_sync_at` velho + `last_incremental_at` recente → `last_success_at` = o incremental no JSON; implementação que iguala ao full REPROVA."
- `M-09/F-02:89-91`: FE-badge negative fixture present.
- **Writer semantics independently confirmed**: `sync_state_repo.go` `ON CONFLICT ... DO UPDATE SET last_full_sync_at = COALESCE(EXCLUDED.last_full_sync_at, sync_state.last_full_sync_at), last_incremental_at = COALESCE(EXCLUDED.last_incremental_at, sync_state.last_incremental_at)` with `timestampArg(!incremental, at), timestampArg(incremental, at)` — only the column matching the run type is written. The r04 premise (full freezes) is a repo fact, and GREATEST is the only honest definition. Fix arm holds.

### F-r04-2 — **CLOSED**

- `IC-05:47-55` pins reference/pointer injection, names the forbidden idiom verbatim ("Builder por VALOR que retorna cópia (idioma `root.go:854-857` `pricingHandler = pricingHandler.WithCalc(...)`) é PROIBIDO aqui"), and requires "Prova obrigatória na ROTA, não na porta".
- Idiom independently confirmed at `root.go:856` — inside the cited range; the citation is accurate.
- Propagated: `M-09/F-01:101-105` (fake injected via setter asserted on the **mounted route**, not the isolated service), `M-09-sync-observability/milestone.md:77-79` (Done Means: seam proved on the registered route), `M-08/F-02:56-61` (Expected Output: "injeção por referência/ponteiro, nunca builder-cópia").

### F-r04-3 — **CLOSED**

- `IC-03:61-66` states the bucket enum comes from the **existing** `domain.DeriveOrderBucket` at `orders/domain/order_bucket.go:48` "assinatura e truth table INALTERADAS; já vive no núcleo", that M-03 moves only the CALL SITE, and explicitly deletes the prior false prose ("a prosa anterior 'transport de market' / 'função MOVE p/ o núcleo' era FALSA contra o repo — deletada").
- Independently confirmed: `order_bucket.go:48 func DeriveOrderBucket(providerStatus, shipmentStatus string, tags []string, faturado bool) OrderBucket` — already in `orders/domain`. False prose is gone (grepped for "transport de market" across the tree: zero hits). Deletion, not honest-unknown — correct per R-25.

### F-r04-4 — **CLOSED, and every named anchor is real**

- `M-08/F-02:25-28` (Brief) and `:44-50` (Inputs) re-source the mapping to `integration_installations.external_account_id`, read "pelo repo/serviço de installations EXISTENTE (`installation_repo.go:81` ListInstallations já seleciona `external_account_id`) atrás de PORTA do package webhook — NUNCA SQL cross-module direto; sem store novo". `AccessTokenResolver` is retained only as token source ("não é invertível p/ lookup").
- `IC-04:41` DrainInbox row fixed: "user_id→installation via `integration_installations.external_account_id` lido por porta do package (F-r04-4)".
- Anchors verified EXACT: `installation_repo.go:81`, `auth_adapter.go:192`, `auth_adapter.go:261`, `auth_flow_service.go:691`, column at `migrations/0016:33`. The cross-module boundary rule (AGENTS.md ports/adapters) is respected.

### F-r04-5 — **CLOSED**

- `mission.md:303-306`: "O BLOCO DE IMPORTS do root.go não pertence a nenhum milestone — hub-resolved: cada milestone adiciona/remove SÓ os imports que a própria região ancorada exige; conflito mecânico de import = merge do hub (P5 r04 F-r04-5)."
- Matrix cell `mission.md:298` names `root.go:99,101`; `M-07/F-01:89-92` Owned paths name "remoção dos imports `root.go:99,101` (tarifflive/tariffcomposite — sem eles o arquivo não compila)".
- Line numbers verified EXACT against `root.go`.

### F-r04-6 — **CLOSED**

- `mission.md:320-327` enumerates the residual overlaps: "arquivo da allowlist do guard (dono M-02 F-04) → escrito por M-03 F-03 (remoção A/B) e M-07 F-01 (remoção C/D), serializados pelo edge lane B → lane C; bloco de imports do root.go → hub-resolved, sem dono". Universal claim is scoped, not absolute: "Nenhum write-set overlap **restante** sem edge nomeado ou resolução do hub registrada".
- The lane B → lane C edge is real in the DAG (`mission.md:265-268`), so the serialization claim is discharged.

**No earlier r01/r02/r03 closure was reopened.** Spot-checked survivors: N-2 (`baseline_commission_percent` intocada — `IC-01:106-109`, `M-07/F-01:38-40`), N-1 (`products_job.go` intocado — `M-02/F-03` Expected Output), N-3 (materialize-on-read totality — verified at `calc_repository.go:239`), N-8 (`sync/application/` additive lock — `mission.md:314-317`), P-2 (≤1 FE contract commit — `mission.md:307-309`), P-5 (M-05 root.go cell `—` — `mission.md:296`, `M-05/F-01` Expected Output), F-8 (`api_shipping_options` reserve — `IC-01:67-69`).

---

## 3. New findings

### F-r05-1 — BLOCKING — checks 4, 5, 6

- **Locus**: `.mnfs/MIS-007-ml-sync/M-09-sync-observability/F-01-sync-health-endpoint/feature.md:36-37` (also `:83`).
- **Offending token (verbatim)**: "SEM param `installation_id` — tenant vem do ctx (IC-05 pina); a leitura **varre as instalações ML do tenant**." And `:83`: "Tenant sem instalação ML ou sem sync ainda → entities todas null-honestas".
- **Why it is a defect**: the scan predicate is scoped to ML installations, but `products` — the one live entity the milestone is required to show — is registered in `sync_state` under a non-ML installation id. Verified in repo: `apps/server_core/internal/modules/sync/composition/scheduler.go:11` `const InstallationScopeERP = "erp"` ("the installation_id used for ERP-sourced entities (products via xlsx/sankhya) that are **not tied to a marketplace installation**"), consumed at `products_job.go:50 syncapp.NewScheduler(repo, InstallationScopeERP, interval, time.Now)`. An ML-installation scan therefore returns **zero** products rows. This contradicts three ratified statements simultaneously:
  - `research/sync-health-interface-contract.md:91` — "uma row por entidade **REGISTRADA em `sync_state`** (products incluído — o endpoint é entidade-agnóstico)" (IC-05 is declared "binding integral" by `M-09/milestone.md:17`);
  - the brief's own EARS at `M-09/F-01:44-45` — "While sync_state tem row de products, when GET roda, the entity products shall carregar timestamps reais";
  - `M-09/milestone.md:74-75` Done Means — "`GET /sync/health` responde payload IC-05 com **products (entidade viva de MIS-006) real**".
  
  It is also the only live-data proof the milestone has before M-04/M-06 land, so the brief as written makes its own Done Means unreachable. This is a plan-time narrowing (check 6) resting on an unverified assumption about where products' sync_state rows live (check 5), producing a contract the IC set cannot satisfy (check 4).
- **Yes-if**: `M-09/F-01:36-37` replaces "varre as instalações ML do tenant" with a scan of **all `sync_state` rows for the tenant regardless of `installation_id`** — explicitly naming the ERP scope sentinel `sync/composition/scheduler.go:11` (`installation_id = "erp"`) as in-scope — and `:83` is re-worded so the empty case is "tenant sem NENHUMA row em sync_state", matching IC-05:91's entity-agnostic requirement.

### F-r05-2 — BLOCKING — checks 4, 6

- **Locus**: `.mnfs/MIS-007-ml-sync/M-06-orders-backfill-decomposition/F-03-audit-fe-pedidos/feature.md:27-28` (restated as a Constraint at `:61`).
- **Offending token (verbatim)**: "esperado_unit = detail.percentage_fee × unit_price/100 + detail.fixed_fee; esperado_total = esperado_unit × quantity — |delta| > R$0,01 (IC-01 tolerância) → divergência kind=`tarifa`".
- **Why it is a defect**: the camada-2 `detail` ratified by IC-01 has **five** keys, one of which is a fee component the formula silently drops. `research/channel-fees-interface-contract.md:119-120` canonical example: `"detail":{"percentage_fee":12.5,"fixed_fee":6.0,"financing_add_on_fee":0,"price_used":79.90,"listing_type_id":"gold_special"}`. The producer writes it: `M-05-listings-fees-divergence/F-01-camada2-fee-ingest/feature.md:28` — "`detail` {percentage_fee, fixed_fee, **financing_add_on_fee**, price_used, listing_type_id} (canonical IC-01 — camada 2 NÃO é percent seco)". `research/divergences-interface-contract.md:92` pins the semantics as "expected = camada 2 (IC-01) **valorada no preço da venda**" — i.e. the whole camada-2 fee revalued, not a two-term subset. With any listing where `financing_add_on_fee > 0`, the pinned formula understates `esperado` by that amount on every line and opens a permanent, never-resolving `tarifa` divergence — the exact false alarm the same brief promises to avoid ("auditoria muda, nunca falso alarme", `:26`). The formula is pinned as non-negotiable by the brief's own Constraint at `:60-62` ("pinar a fórmula") and divergence semantics are Must-Not-Decide-In-Feature-Execution (`IC-02`), so execution cannot honestly repair it — this is a plan-time defect, not an implementation escalation. (Corroborating smell: IC-01's own example does not reconcile — `12.5% × 79.90 + 6.00 = 15.99`, but the row's `value` is `16.45`, a 0.46 gap the two-term formula cannot explain.)
- **Yes-if**: `M-06/F-03:27-28` and `:61` state the expected-unit formula over **every** fee component IC-01 declares in the camada-2 `detail` — minimally `esperado_unit = detail.percentage_fee × unit_price/100 + detail.fixed_fee + detail.financing_add_on_fee` — **or** the brief pins, verbatim and with the ML-API evidence line that supports it, why `financing_add_on_fee` is excluded from the seller-borne commission; and IC-01's canonical example at `:118-120` is made arithmetically self-consistent with whichever formula is ratified.

### F-r05-3 — ADVISORY — checks 2, 3

- **Locus**: `.mnfs/MIS-007-ml-sync/M-05-listings-fees-divergence/F-03-anuncios-fe-contract/feature.md:90-93` vs `mission.md:312-313` and `mission.md:296`.
- **Offending token (verbatim)**: F-03 Owned paths — "`listings/transport/`, `listings/application/read_*`, `listings/adapters/postgres/repository.go` (query read aditiva)"; the registered grant reads "lock aditivo registrado: **arquivos de application de listings**, additive-only".
- **Why it is a defect**: `M-04-listings-backfill-ingest/F-02-mass-closure-replacement/feature.md:65-66` owns `listings/adapters/postgres/repository.go` outright, and `mission.md:295` gives M-04 the whole `listings/**` tree. The recorded additive-lock is worded strictly narrower (application files only) than the write-set three M-05 features actually claim (transport + adapters/postgres). The M-05 Go-packages matrix cell (`mission.md:296`) likewise names only "ingest ext camada2/divergência (dentro de listings app…), FE /anuncios". Time-safety holds — the serial edge M-05→M-04 (`mission.md:278`) and "PÓS-close" make simultaneity impossible — so this is a recorded-grant/write-set mismatch, not a live collision; hence advisory. But a P6 gate reading the grant literally would flag M-05 F-03's adapter and transport edits as out-of-grant.
- **Yes-if**: `mission.md:312-313` widens the registered grant to name the actual surfaces ("arquivos de application, transport e `adapters/postgres/repository.go` de listings, additive-only, pós-close do M-04"), and the M-05 Go-packages cell at `mission.md:296` is aligned to the same enumeration.

### F-r05-4 — ADVISORY — check 4

- **Locus**: `research/sync-health-interface-contract.md:107-108` vs `M-09-sync-observability/milestone.md:34-35` and `:46-47`.
- **Offending token (verbatim)**: IC-05 — "`last_incremental_at` REAL exige o fix `incremental=false` (ADR-08, M-02) — pré-condição nomeada; **sem ela QA reprova**." M-09 milestone Dependencies — "Nenhuma dura (sync_state 0075 existe). **SOFT**: M-02 F-03 (last_incremental_at vira significativo)"; and `:34-35` — "campos que dependem do fix incremental (M-02 F-03) ficam **honestos** até lá (IC-05 §NULLs)".
- **Why it is a defect**: IC-05 is declared "binding integral" for M-09 (`M-09/milestone.md:17`), so the two artifacts hand a QA validator opposite instructions for the same lane-A milestone: one says the milestone fails without M-02 F-03, the other says lane A has no hard dependency and the field is honestly null until then. Repo evidence favors the milestone's arm — `sync/application/scheduler.go:160` hardcodes `RecordSuccess(..., false)`, so pre-fix `last_incremental_at` is uniformly NULL, which is honest-unknown, not a lie — and M-09's mandatory negative fixture (`M-09/F-01:98-100`) is fixture-seeded and therefore satisfiable without M-02. Advisory rather than blocking: no build-order hazard, but the gate criterion is ambiguous and could reject a correct M-09.
- **Yes-if**: `IC-05:107-108` is scoped so the failure condition attaches to the mission-level/live criterion for an entity that actually runs incremental (i.e. after M-04/M-06 register phase-bearing jobs), and explicitly names M-09's pre-fix contract (uniform NULL = honest, per IC-05 §NULLs) as passing — **or** `M-09/milestone.md:46-47` promotes M-02 F-03 to a hard dependency and the DAG at `mission.md:265-269,272-286` gains the corresponding M-09→M-02 edge.

### F-r05-5 — ADVISORY — checks 6, 7

- **Locus**: `.mnfs/MIS-007-ml-sync/M-05-listings-fees-divergence/F-01-camada2-fee-ingest/feature.md:38`.
- **Offending token (verbatim)**: EARS — "the sistema shall upsert row camada 2 com value amount + detail (**percentage_fee/fixed_fee/price_used**) + coletado_em".
- **Why it is a defect**: the same brief's own Brief line `:28` and IC-01's canonical `detail` (`:119-120`) require five keys; the EARS enumerates three, dropping `financing_add_on_fee` and `listing_type_id`. EARS is the executable acceptance sentence — an implementer discharging it literally writes a `detail` that is contract-incomplete, and the missing `financing_add_on_fee` is precisely the key F-r05-2 turns on. Narrowing of a ratified value inside the propagation chain.
- **Yes-if**: `M-05/F-01:38` enumerates the full IC-01 `detail` tuple (`percentage_fee, fixed_fee, financing_add_on_fee, price_used, listing_type_id`) or references it as "detail canonical IC-01 (5 chaves)" without an abbreviated inline list.

### F-r05-6 — ADVISORY — check 7

- **Locus**: `.mnfs/MIS-007-ml-sync/mission.md:368-373` (Handoff) vs `mission.md:6,13` (front matter).
- **Offending token (verbatim)**: "Current status: draft — P3 fechado … **planning_phase = architecture (P4)**." / "Next action: **autorar interface contracts compartilhados (P4)**; depois corpos de milestone + briefs + Parallel Execution Plan (P5)."
- **Why it is a defect**: front matter says `planning_phase: decompose` (`:13`), and all seven interface contracts plus nine milestone bodies and twenty-three feature briefs exist and are manifested. The Handoff block — the artifact a downstream P6/P7 reader consults for "what is next" — points at work already completed, and is the only remaining place in the mission file that still describes the plan as pre-P4. The waiver row at `:374-376` (Sol P3/P5/P7 retroactive from 2026-08-05, mandatory before `status: planned`) is current and must survive the edit.
- **Yes-if**: `mission.md:368-373` is updated to the actual state (`P5 fechado, planning_phase = decompose`; next action = P6 dual gate over the decomposition), preserving the waiver row at `:374-376` verbatim.

### Checks that produced no finding

- **Check 1 (DAG)**: every edge in `mission.md:272-286` is justified by a named artifact and each is corroborated by the feature briefs; no false edge found. Soft edges are marked `⤳` and distinguished. The one candidate missing edge (M-09→M-02) is adjudicated as F-r05-4 rather than a check-1 finding, since the milestone's own arm is repo-correct.
- **Check 2 (six-axis / migrations)**: ownership cells are disjoint on all six axes given the lane order; migration ranges 0086-0089 / 0090-0092 / 0093 / 0094-0095 are disjoint and start above the verified tip 0085. The out-of-lane-order allocation (M-08 lane D gets 0093 while M-06 lane C gets 0094-0095) was checked against the actual runner — `internal/platform/migrate/runner.go:68-82` skips already-applied filenames and applies the rest in lexicographic order, tracking per filename in `schema_migrations`, so late arrival of a lower number is safe and the two tables are independent. **No finding.**
- **Check 3 (feature write-DAG)**: the guard-allowlist overlap and root.go import block are both named and resolved (F-r04-6, verified); the only residual mismatch is F-r05-3.
- **Check 5 (prerequisite existence)**: all 28 anchor classes named by the plan were verified directly in the repo (table in §1). Zero false anchors. The single unverified assumption found was the products/installation-scope premise, reported as F-r05-1.
- **Check 7 (section completeness)**: all 23 feature briefs carry Brief, Inputs, Expected Output, Constraints, Negative Scenarios, Ownership, Validation Expectations (mechanically checked). All four FE-bearing briefs (M-05 F-03, M-06 F-03, M-07 F-02, M-09 F-02) carry an Interaction Model. Validation expectations name minimum inspectable proof throughout (must-fail on guards, negative fixtures, route-level seam assertions, golden DTOs, before/after on the same screen).
- **Check 8 (no implementation planning / no new scope)**: briefs pin contracts, anchors, and formulas but defer technique to `spec.md` ("técnica exata = spec", `M-02/F-04` Inputs). No new product scope was introduced beyond the ratified design; out-of-scope items are explicitly recorded as such (`IC-04:96-97` missed_feeds; `IC-01:67-69` `api_shipping_options` reserve; `M-07/F-01:72-73` camada-1 dormant arm).

---

## 4. Verdict

**NEEDS-REVISION.**

Blocking:

- **F-r05-1** — `M-09/F-01:36-37,83`: the `/sync/health` scan predicate is scoped to ML installations and therefore cannot return `products` (registered under `installation_id = "erp"`, verified at `sync/composition/scheduler.go:11` and `products_job.go:50`), contradicting IC-05:91, the brief's own EARS at `:44-45`, and M-09's Done Means at `milestone.md:74-75`.
- **F-r05-2** — `M-06/F-03:27-28,61`: the pinned `tarifa` expected-value formula drops `financing_add_on_fee`, a camada-2 `detail` component ratified in `IC-01:119` and written by `M-05/F-01:28`, guaranteeing systematic false divergences under IC-02:92 semantics that execution is forbidden to re-decide.

Advisory (does not gate): F-r05-3 (listings additive-lock grant narrower than the claimed write-set), F-r05-4 (IC-05 "QA reprova" vs M-09 SOFT dependency), F-r05-5 (M-05/F-01 EARS enumerates 3 of 5 `detail` keys), F-r05-6 (mission.md Handoff stale at P4).

All six r04 findings are verified closed at their loci, with the chosen yes-if arms' rationale independently confirmed against the repo; no r01/r02/r03 closure was reopened.

Waiver note: this audit substitutes the GPT-5.6 Sol P5 touchpoint under the operator-ratified waiver. Sol's retroactive P5 pass remains mandatory before `status: planned`, regardless of this verdict.
