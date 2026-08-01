# P5 Claude Decomposition Audit — round 07 (VERBATIM)

```yaml
type: planning-review-verbatim
phase: P5
round: 07
task_id: a09131a36e7a7dbb2
auditor: cold Claude Opus crew (general-purpose, read-only mandate; operator-ratified
  waiver -- Sol P5 retroactive mandatory before status: planned)
input_manifest: p5-input-r07.sha256 (digest 60f898ef70b6f932fca22b49e3b3d8570e4e0d28e9776798a25b3c6741b17ea3)
persisted: 2026-07-31, same turn the task notification arrived
provenance: recovered VERBATIM from session transcript JSONL (longest task-notification
  candidate, html-unescaped, <result> block) -- SEVENTH round using this path
verdict: NEEDS-REVISION (PART A: all F-r06-1..6 CONFIRMED CLOSED, zero reopened;
  1 blocking F-r07-3 + 5 advisory F-r07-1/2/4/5/6 + 1 stale-blocker observation)
```

---

# P5 Decomposition Audit — MIS-007-ml-sync — Round 07

```yaml
type: planning-review-audit
phase: P5
round: 07
mission: MIS-007-ml-sync
date: 2026-07-31
auditor: cold independent decomposition auditor (read-only)
mission_root: C:\Users\leandro.theodoro\Documents\marketplace-central\.mnfs\MIS-007-ml-sync
repo_root: C:\Users\leandro.theodoro\Documents\marketplace-central
input_manifest: planning-reviews/p5-input-r07.sha256
manifest_self_digest_expected: 60f898ef70b6f932fca22b49e3b3d8570e4e0d28e9776798a25b3c6741b17ea3
manifest_self_digest_recomputed: 60f898ef70b6f932fca22b49e3b3d8570e4e0d28e9776798a25b3c6741b17ea3  # MATCH
manifest_entries_verified: 46/46 OK (sha256sum -c, zero failures)
scope_audited: exactly the 46 manifested files; planning-reviews/* read only as PART A context
```

---

## PART A — Closure verification of r06 findings

### F-r06-1 (blocking, r06) — layer-3 detail mandate unscoped by fee_kind — **CONFIRMED CLOSED**

Fix present at all three declared loci, and the three loci agree.

`research/channel-fees-interface-contract.md:142-147` (Canonical Examples, rejection):

> Rejeição canônica: upsert de camada 3 **fee_kind=`commission`** com `detail` sem `sale_fee_unit`/`quantity` → writer recusa (...) row camada 3 fee_kind=`freight` do shipment (produtor M-06 F-02: subject_type=`order`, origem=`api_shipment`) NÃO tem decomposição sale_fee_unit/quantity, `detail` NULL é aceito (auditoria P5 r06 F-r06-1).

`research/channel-fees-interface-contract.md:106-111` (Persistence Expectations) now defines both siblings once — camada 3 comissão `value = sale_fee_unit × quantity` with `detail` OBRIGATÓRIO, and:

> Camada 3 frete (produtor M-06 F-02): `value` = custo seller do shipment; `subject_type=order`, `origem=api_shipment`; `detail` NULL permitido — NÃO existe decomposição sale_fee_unit/quantity p/ frete (F-r06-1).

`M-02-sync-core-seam/F-02-fee-divergence-ports/feature.md:31-34` (guard clause) and `:75-76` (Negative Scenario) both scope the guard to commission and make the freight-accepted side a test obligation:

> UpsertFee camada 3 fee_kind=`commission` sem detail obrigatório → erro nomeado (IC-01); camada 3 fee_kind=`freight` sem detail → ACEITO (F-r06-1 — o teste cobre os DOIS lados).

Cross-consistency re-check (IC-01 Persistence × IC-01 Canonical Examples × M-02/F-02 × M-06/F-02 × M-07/F-01): the producer brief `M-06-orders-backfill-decomposition/F-02-decomposition-camada3/feature.md:29-31` emits the freight row as `fee_kind=freight, subject_type=order, subject_id=<provider_order_id>, origem=api_shipment` — byte-consistent with the two IC-01 loci and with the M-02 writer guard. No new contradiction introduced. M-06 Done Means (`frete_seller`) and IC-03's canonical `−22.90` freight term are reachable.

### F-r06-2 (blocking, r06) — `ResolveListingFees` tuple lacked `detail` — **CONFIRMED CLOSED**

`research/channel-fees-interface-contract.md:71-77` (Required Outputs):

> `ResolveListingFees` retorna por fee_kind: `{value, value_type, currency, layer, detail, origem, coletado_em}` (...) `detail` = o jsonb VERBATIM da row resolvida (camada 2 = tupla canônica de 5 chaves; NULL quando a row não tem detail) — declarado UMA vez aqui, escopado por camada, nunca re-inventado por consumidor; sem ele o braço camada-2 do M-07 (`detail.percentage_fee/fixed_fee`) seria insatisfazível contra o port (auditoria P5 r06 F-r06-2).

`M-02-sync-core-seam/F-02-fee-divergence-ports/feature.md:40-42` restates the full tuple in EARS ("the reader shall retornar {value, value_type, currency, layer:2, detail, origem, coletado_em} — detail = jsonb VERBATIM da row (IC-01 §Required Outputs; F-r06-2)"). `M-07-pricing-fee-read/F-01-fee-read-resolver/feature.md:45-46` is unchanged and now satisfiable through the port alone. Semantics stated once, layer-scoped; no consumer re-invents them. No collision with the F-r06-1 scoping: `detail` NULL is a legal tuple value, and only the camada-3 commission writer path mandates content.

### F-r06-3 (advisory, r06) — missing M-06⤳M-05 soft edge — **CONFIRMED CLOSED**

`mission.md:272` (DAG block): `M-06 (⤳ M-05: auditoria 3→2 precisa de camada 2 populada — fica MUDA sem ela, não quebra)`.
`mission.md:288` (edge-justification table): `| M-06⤳M-05 | auditoria 3→2 precisa de camada 2 populada — sem ela auditoria fica MUDA, não quebra (já declarado em M-06/milestone; mesma classe do M-07⤳M-05 — F-r06-3) |`. Both present, with the non-blocking qualifier matching M-07⤳M-05.

### F-r06-4 (advisory, r06) — residual 3-key detail enumerations — **CONFIRMED CLOSED**

`M-05-listings-fees-divergence/F-01-camada2-fee-ingest/feature.md` (Validation fixture) and `M-05-listings-fees-divergence/F-03-anuncios-fe-contract/feature.md` (DTO shape) both now carry the full 5-key canonical tuple `{percentage_fee, fixed_fee, financing_add_on_fee, price_used, listing_type_id}`. No 3-key enumeration survives anywhere in the manifested set (swept).

### F-r06-5 (advisory, r06) — drifted anchor ranges — **CLOSED AT ALL SIX ENUMERATED LOCI (measured), TWO RESIDUALS OF THE SAME CLASS REMAIN** (see F-r07-2, F-r07-3)

Every corrected anchor re-measured against the repo working tree, not accepted from prose:

| Claim | Locus | Measured | Verdict |
| --- | --- | --- | --- |
| `root.go:845-851` = `if feeReader…` block; `:844` survives | M-07/F-01:36-38 | `:844` `var tariffResolver pricingports.TariffResolver = calcTariffResolver`; `:845` `if feeReader, ferr := marketplaceCapabilities.FeeQuoteReader(...)`; `:851` `}` | correct |
| IC-05 builder idiom `root.go:856` | IC-05:49-51 | `:856` `pricingHandler = pricingHandler.WithCalc(calcSvc)` | correct |
| `writer.go:74-95` (`upsertSQL`) + `:104-112` (`keepAbsentSQL`) | M-02/F-02:52-53, M-04/F-02, IC-03:99, IC-07:97, mission.md:96, mission.md:151-153 | `:74` `const upsertSQL = \`` … `:95` `updated_at = now()\``; `:97-103` doc comment; `:104` `const keepAbsentSQL = \`` … `:112` | correct at all six |
| `listings_test.go:25` = required-substring list; `:101` = regex body extractor | mission.md:215, codebase-ingest-side.md | `:24` `for _, required := range []string{`, `:25` `"CREATE TABLE IF NOT EXISTS listings",`; `:99` `func createTableBody(...)`, `:101` `re := regexp.MustCompile(...)` | correct |

Sampled beyond the r06 list (all correct): `capability_adapter.go:681-688` providerDiag; `capability_adapter.go:26` AccessTokenResolver; `root.go:370-378` credential wiring; `root.go:590-592`/`:591-592` live readers A/B; `root.go:661-678` scheduler region and `:668-678` products-scheduler idiom; `order_bucket.go:48` `DeriveOrderBucket`; `repository.go:390-394` MASS-CLOSURE UPDATE; `installation_repo.go:81` `ListInstallations`; `sync_state_repo.go:35` `Read`; `composition/scheduler.go:11` `const InstallationScopeERP = "erp"`; `run_read_handler.go:34` `/sync/runs` mount; `route_deadline.go:23-28` class consts; `query.go:31` `filter.` prefix; `index.ts:2144` `listSyncRuns` method idiom; `IntegracoesPage.tsx:508` `listIntegrationInstallations`.

### F-r06-6 (advisory, r06) — missing `## Inputs/Outputs` — **CONFIRMED CLOSED**

`M-08-webhook-ingest/F-02-worker-callback/feature.md:73-80` now carries the section (inbox FIFO input with IC-04 §Enums status machine verbatim; Output 1 = `IngestOrder` call; Output 2 = the real `WebhookStatsReader` with IC-05 §InboxHealth shape and reference-semantics injection), self-labelled "(Seção adicionada — auditoria P5 r06 F-r06-6.)". `M-04-listings-backfill-ingest/F-02-mass-closure-replacement/feature.md` likewise. Full sweep of all 22 briefs confirms the section is present in every API/data-shape brief; the four without it are non-API/non-data (M-01/F-01 decorator, M-02/F-04 guard test, M-04/F-04 wiring) plus the FE brief M-09/F-02 (see F-r07-5).

### Reopened-closure sweep (r01–r05) — **CLEAN**

| Prior closure | Spot-check locus | Result |
| --- | --- | --- |
| F-r05-1 all-rows `sync_state` scan + ERP sentinel | `M-09/F-01:35-40` — "a leitura varre TODAS as rows de `sync_state` do tenant, independente de `installation_id` — INCLUI o sentinela de escopo ERP `installation_id = "erp"` (`sync/composition/scheduler.go:11`)…scan restrito a instalações ML devolveria ZERO products"; Negative Scenario `:83-85` | intact (sentinel verified in repo at `scheduler.go:11`) |
| F-r05-2 three-term formula | `M-06/F-03` Brief + Constraints ("percentage_fee × unit_price/100 + fixed_fee + financing_add_on_fee … dropar financing_add_on_fee abriria divergência falsa PERMANENTE"); IC-01:124-132 example `15.99` (12.5% × 79.90 + 6.00 + 0 = 15.99) | intact |
| ADR-14 ≤1-COMMIT FE contract in flight | `mission.md` lanes note + matrix cross-rule; IC-03:156-158; IC-07:144-146 | intact |
| GREATEST `last_success_at` | IC-05:95-101 with mandatory negative fixture; restated at M-09/F-01 Inputs and M-09/F-02 Validation | intact |
| `external_account_id` webhook mapping | `M-08/F-02:44-50` — port over existing installations repo, no new store, `AccessTokenResolver` explicitly non-invertible | intact |
| write-DAG enumeration + widened listings grant | `mission.md:317-321` (`listings/application/**`, `listings/transport/**`, `listings/adapters/postgres/repository.go`) and `:328-335` | intact |
| phase-tolerant-parse scoping | `M-09/F-01` Constraints ("cursor de formato desconhecido → phase null (nunca erro)") + IC-05:102 | intact |
| ADR-06 absent ≠ closed | IC-07:81-82; IC-06 run-complete rule; M-04/F-02 | intact |

No r01–r05 closure reopened.

---

## PART B — Fresh full audit (8 checks over the 46 manifested files)

### Check 1 — DAG edge completeness and justification — **1 advisory finding**

Every hard edge in `mission.md:274-292` is justified by a named forcing artifact, and every cross-milestone consumption I could trace has a matching edge: M-03/M-04→M-01 (adapter), M-03/M-04/M-05/M-06/M-07→M-02 (DDL + ports), M-05→M-04, M-06→M-03 (writer heredity), M-08→M-06, M-08→M-09 (`WebhookStatsReader` port), plus soft M-07⤳M-05, M-06⤳M-05, M-09⤳M-04/M-06. Lanes A/B/C/D are consistent with the edge set.

**F-r07-1 — advisory — DAG edge completeness**
Cited excerpt (`research/sync-health-interface-contract.md:108-113`):
> `last_incremental_at` REAL exige o fix `incremental` (ADR-08, M-02 F-03) — pré-condição nomeada do critério de MISSÃO/live: (...) a dependência M-02 F-03 permanece SOFT no M-09 (auditoria P5 r05 F-r05-4).

Defect locus: `mission.md:274-292` (DAG block + edge-justification table) — no `M-09⤳M-02` row.
Offending token: the absent edge. IC-05 names a soft dependency of M-09 on M-02 F-03 that exists nowhere at mission level, while the exactly analogous M-07⤳M-05 and M-06⤳M-05 soft edges are both declared there (the latter added by F-r06-3 for precisely this reason). M-09 and M-02 are the two milestones the mission runs *in parallel in lane A*, so the hub adjudicating lane-A close order is the reader who most needs the edge and is the one who cannot see it.
Yes-if: `mission.md` gains a `M-09⤳M-02` row carrying IC-05's own qualifier (pre-fix uniform NULL is honest, M-09's gate passes without it), matching the form of the M-06⤳M-05 row.

### Check 2 — Canonical six-axis disjointness, seam locks, migration ranges — **1 advisory finding**

Migration ranges: repo tip measured at `apps/server_core/migrations/0085_erp_import_products_sellable_fields.sql`. Allocations M-02 `0086-0089`, M-04 `0090-0092`, M-08 `0093`, M-06 `0094-0095` (reserve, gap permitted) are disjoint, contiguous above tip, and each milestone body agrees with its matrix cell (M-08/F-01 `0093`; M-04 `0090-0092`; M-03, M-05, M-07, M-09 declare no migration and take none). No milestone ALTERs another's table: `listings` DDL is deliberately M-04 (IC-07 Boundary), `orders`/`order_shipments`/`channel_fees`/`divergences` DDL is M-02, `notifications_inbox` is M-08.

Shared-table writes in the same lane are partitioned and declared: lane C has M-05 writing `channel_fees` camada 2 + `divergences` kind `estoque`, M-06 writing camada 3 + kind `tarifa` — disjoint row partitions through a port merged in lane A, with IC-02's one-open-row-per-(entity,kind) partial unique keeping them non-interfering. Additive locks are registered for both contested packages (`listings/**` M-04→M-05; `sync/application/**` M-02→M-09, M-09 additive-only with `scheduler.go` untouched).

**F-r07-2 — advisory — six-axis ownership (root.go axis)**
Cited excerpt (`mission.md:262`, M-07 row, root.go cell):
> edita região pricing existente `:828-858` + remoção de imports tarifflive/tariffcomposite (`root.go:99,101` — F-r04-5); **única exceção**, hub arbitra

Cited counter-excerpt (`M-03-orders-shipment-persist/F-03-read-path-switch/feature.md:27-28`):
> (root.go:597-599) troca `ordersShipmentReaderAdapter`/`ordersBuyerFiscalReaderAdapter` (root.go:591-592) por readers de banco; sites vivos A (shipment) e B (fiscal) DELETADOS

Defect locus: `mission.md:258` (M-03 row, root.go cell = `1 linha ancorada`) contradicted by `M-03/milestone.md:56` ("região orders `:576-601` — edita SÓ dentro dela") and M-03/F-03 above.
Offending token: `única exceção` (and the understated `1 linha ancorada` for M-03).
Measured: `root.go:591` `ordersShipmentReader := newOrdersShipmentReaderAdapter(...)`, `:592` `ordersBuyerFiscalReader := newOrdersBuyerFiscalReaderAdapter(...)`, `:597-599` the `ordersEnrichSvc` chain — all pre-existing lines M-03 F-03 rewrites/deletes. M-03 is therefore the *same class* of exception as M-07: an in-place edit of an existing root.go region with deletions, not an appended constructor line. The matrix is the artifact the hub adjudicates root.go collisions with, and it carries a false universal. No unresolvable collision follows (lane B → lane C serializes M-03 and M-07; root.go is hub-resolver-of-record either way), hence advisory.
Yes-if: the M-03 root.go cell states the region (`:576-601`, in-place edit incl. deletion of the A/B reader lines) and the M-07 cell's `única exceção` is either removed or restated as "duas exceções (M-03 região orders, M-07 região pricing)".

### Check 3 — Feature-level write-DAG — **PASS**

Every write-set overlap I could construct has a named serial edge or a recorded hub resolution: allowlist file (owner M-02 F-04) written by M-03 F-03 and M-07 F-01, serialized lane B → lane C; `orders/**` M-03 then M-06 by lane heredity; `listings/**` M-04 then M-05 under the widened additive lock; `sync/application/` M-02 owner + M-09 additive-only; root.go import block hub-resolved with no owner; root.go anchored regions one per milestone (M-03 `:576-601`, M-06 `:661-678`, M-07 `:828-858`, M-08/M-09 own lines). `mission.md:328-335` enumerates these and closes with the total claim "Nenhum write-set overlap restante sem edge nomeado ou resolução do hub registrada" — which held against every overlap I could construct. M-05 explicitly disclaims root.go at both F-01:57 and F-02:51-52, matching its matrix cell `—`.

### Check 4 — Contract satisfiability — **1 blocking finding**

Arithmetic re-derived in every IC canonical example: IC-01 camada 2 `12.5% × 79.90 + 6.00 + 0 = 15.99` ✓; IC-01 camada 3 `8.11 × 3 = 24.33` ✓; IC-03 `239.70 − 48.66 − 22.90 − 95.10 = 73.04` ✓ and `73.04 / 239.70 = 30.47%` ✓; IC-03 partial example correctly omits `liquido`/`margem_pct` with `incompleto:["custo_produto"]`. IC-02/IC-04/IC-05/IC-06/IC-07 examples carry no arithmetic. `sale_fee` per-unit semantics propagate identically through IC-01, IC-03, M-06/F-02 and M-06/F-03.

**F-r07-3 — BLOCKING — contract satisfiability / honest-unknown**
Cited excerpt (`research/sync-health-interface-contract.md:103-107`):
> `webhook`: estado canônico inicial EXATO quando inbox ainda não existe: `{"last_notification_at":null,"pending":0,"dropped_24h":0}` (**timestamp null = nunca observado**; contadores ZERO = contagem vazia é fato conhecido). M-09 pode fechar antes do M-08 — FE discrimina pelo `last_notification_at === null` e renderiza **"webhook não configurado"**.

Propagated verbatim to `M-09-sync-observability/F-02-integracoes-health-section/feature.md:35` and to its third EARS at `:42-45`:
> While webhook block está no estado canônico inicial (last_notification_at null), when card renderiza, the subseção shall dizer "webhook não configurado" (não esconder, **não inventar** (...))

Defect locus: `research/sync-health-interface-contract.md:106-107` (source of the mapping) and `M-09/F-02:35` + `:42-45` (consumer).
Offending token: the string `"webhook não configurado"` bound to the discriminator `last_notification_at === null`.
Defect: the contract defines the null in the same sentence as an **activity** fact ("nunca observado") and then instructs the screen to render it as a **configuration** verdict. The two are the same null. There is a hard edge `M-08→M-09`, and M-08/F-02 replaces the default reader with a real one reading the inbox — after which a *configured* installation with an empty inbox (registration done and no event yet; ML quiet window; worker wedged; inbox INSERT failing) produces byte-identical `{null,0,0}` and the card asserts the webhook is not configured. No field in the ratified IC-05 payload distinguishes the two states, so the brief's promise is a claim the contract cannot ground — and it inverts the diagnosis in precisely the failure state Q4 ("/integracoes mostra saúde por entidade + última notificação recebida", `mission.md` Quality Attributes) exists to surface. This is the honest-unknown class the mission enforces everywhere else ("Campos de entities sem observação = null JSON (nunca 0, nunca string vazia)", M-09/F-01), applied inconsistently to the webhook block.
Yes-if: the rendered string states the observed fact rather than a configuration verdict (e.g. an activity-phrased label at IC-05:106-107, M-09/F-02:35 and the third EARS at `:42-45`), so that no reachable state makes the card assert something the payload does not carry. (A configuration signal in IC-05 §InboxHealth produced by M-08 would also close it, but that is a payload addition — the label restatement closes it inside ratified scope.)

### Check 5 — Prerequisite existence — **PASS**

Every symbol/anchor a brief assumes was verified in the repo working tree (file:line measured, not quoted from prose): `DeriveOrderBucket` `orders/domain/order_bucket.go:48` (signature intact — IC-03:63-66's correction of the earlier false prose holds); `buyer_fiscal_reader.go` present in both connectors adapter and orders ports, with the full DTO enumerated in `research/p5-prerequisites.md` §2; `SyncStateRepository.Read` `sync_state_repo.go:35`; `InstallationScopeERP = "erp"` `sync/composition/scheduler.go:11`; `RecordSuccess(..., false)` hardcode `sync/application/scheduler.go:160`; `ListInstallations` `installation_repo.go:81`; `AccessTokenResolver` wiring `root.go:370-378`; `providerDiag` `capability_adapter.go:681-688`; MASS-CLOSURE `UPDATE listings SET status='closed'` `listings/adapters/postgres/repository.go:390-394`; `filter.` prefix parse `listings/transport/query.go:31`; route classes `route_deadline.go:23-28`; `/sync/runs` mount `run_read_handler.go:34`; SDK method idiom `packages/sdk-runtime/src/index.ts:2144`; `IntegracoesPage.tsx:508` fetch idiom and the section order (ActiveSourceCard → SellableAssortmentCard → UploadCard → ProviderConnectCard → ImportacaoSection) measured at `:567-571` inside the cited `:558-574` function body. Everything not yet existing is assigned to a named creating feature (`notifications_inbox`→M-08 F-01; `order_shipments`/`channel_fees`/`divergences`→M-02 F-01; `listing_variations`→M-04 F-01; `WebhookStatsReader` port→M-09 F-01, real impl→M-08 F-02).

### Check 6 — Propagation of approved ADR/IC values into briefs — **2 advisory findings**

Values propagate correctly across the set: migration numbers, `fee_kind`/`origem`/`subject_type` vocabularies, the camada-2 5-key detail tuple, the three-term formula, IC-04 status machine (`received`→`processing`→`done` | re-`received` attempts++ | ≥5 `dropped`; terminals `malformed`/`unmapped`), the 64KB LimitReader and always-200 rule, `dropped_24h`/`pending`/`last_notification_at` key names, routes `POST /webhooks/{provider}` and `GET /sync/health` with SDK `getSyncHealth` and the deliberate SDK-method-absence for the webhook path. M-07/F-02 correctly scopes the provenance vocabulary to `api_listing_prices | config` with `api_shipping_options` reserve-only.

**F-r07-4 — advisory — anchor accuracy (F-r06-5 residual class)**
Cited excerpt (`research/codebase-ingest-side.md:96`):
> Keep-absent (ADR-04): `writer.go:97-105` — `keepAbsentSQL` marca `absent_in_last_snapshot=true` + `stale_since=COALESCE(...)` SÓ nas rows da própria source; nunca DELETE.

Defect locus: `research/codebase-ingest-side.md:96`. Offending token: `97-105`.
Measured: `writer.go:97-103` is the doc comment, `:104` opens `const keepAbsentSQL`, `:105` is `UPDATE products_mirror`, `:112` closes it — the cited range bisects the statement, which is exactly the defect r06 measured. Line `:95` of the same file already carries the corrected `writer.go:74-95` for `upsertSQL`, so the file is internally inconsistent. The r06 residual sweep enumerated IC-03:99, IC-07:97, mission.md:95(→96), mission.md:151 and missed this manifested locus, which is the *source map* the other loci cite.
Yes-if: `:96` reads `writer.go:104-112`, matching the six already-corrected loci.

**F-r07-5 — advisory — anchor accuracy on a deletion instruction (F-r06-5 residual class)**
Cited excerpts (`M-07-pricing-fee-read/milestone.md:33` and `:70`):
> Simulação hoje faz GETs vivos no ML por cotação (root.go:845-850) e um miss cai no…
> …a deleção muda a fiação root.go:845-850 inteira; o `if feeReader…` guard morre junto

Defect locus: `M-07-pricing-fee-read/milestone.md:33` and `:70`. Offending token: `845-850`.
Measured: the `if feeReader, ferr := ...` block runs `:845` through its closing `}` at `:851`. `M-07/F-01:36-38` already carries the corrected `845-851` with `:844` marked SOBREVIVE — so the milestone and its own feature brief disagree, and the stale one is attached to the word `deleção`: an under-range on a deletion instruction leaves a dangling `}`, which is the dangerous direction r06 named.
Yes-if: both milestone occurrences read `845-851`, consistent with M-07/F-01.

### Check 7 — Required brief sections — **1 advisory finding**

Swept all 22 feature briefs: `## Negative Scenarios`, `## Validation Expectations` and `## Ownership` present in 22/22; `## Interaction Model` present in all four UI briefs (M-05/F-03, M-06/F-03, M-07/F-02, M-09/F-02); `## Inputs/Outputs` present in 18/22, absent only from M-01/F-01 (internal decorator), M-02/F-04 (guard test), M-04/F-04 (wiring) — none of which are API/data-shape features — and M-09/F-02. Validation expectations are inspectable rather than assertion-of-presence throughout (golden JSON with exact nulls; the mandatory F-r04-1 negative fixture stated on both the BE and FE side; must-fail fixtures for ADR-06 abort-post-page-1 and ADR-13 starved-observer; route-class proof with the negative control kept outside `registerBatchRoutes`; replay asserted on the *domain effect*, not on inbox row counts).

**F-r07-6 — advisory — required brief sections**
Cited excerpt (`M-09-sync-observability/F-02-integracoes-health-section/feature.md`, section list): brief carries `## Inputs`, `## Expected Output`, `## Constraints`, `## Interaction Model`, `## Negative Scenarios` — no `## Inputs/Outputs`.
Defect locus: `M-09/F-02` section set. Offending token: the absent `## Inputs/Outputs` heading.
Defect: the other three FE briefs of this mission (M-05/F-03, M-06/F-03, M-07/F-02) all carry both `## Inputs/Outputs` and `## Interaction Model`, and r06's F-r06-6 established the section as the binding-source citation for data shapes. M-09/F-02 consumes a payload whose shape is binding (IC-05 entities[] + InboxHealth) and renders discriminated states from it, so the omission is a uniformity gap on the one brief where the payload-to-label mapping is contested (see F-r07-3).
Yes-if: M-09/F-02 gains `## Inputs/Outputs` citing IC-05 §entities + §InboxHealth as the binding input shape and the four render states as output — or the mission records FE briefs as exempt, which would then require removing the section from the other three.

### Check 8 — No implementation planning, no new product scope — **PASS**

No brief plans implementation: sizing/interval/idiom decisions are consistently deferred ("spec pina valor" for the inbox poll interval; "spec segue o idioma da página" for FE data fetching; "spec decide contra idioma web-query"), and every brief closes with `## Execution Artifact Rules` assigning spec/plan/validation to execution. No new product scope appears: the webhook callback registration remains gated on explicit operator authorization (M-08/F-02 Handoff), LISTEN/NOTIFY is explicitly named as a future additive, PK-sentinel consolidation is deferred to a named future mission (IC-07:140), fee columns stay excluded from the E3 set by IC-07's recorded decision, and `api_shipping_options` stays reserve-only.

---

## Verdict

**NEEDS-REVISION**

One blocking finding (F-r07-3). All six r06 findings are CONFIRMED CLOSED and no r01–r05 closure was reopened; the remaining five findings are advisory and do not affect the verdict.

| ID | Severity | Check | Locus |
| --- | --- | --- | --- |
| F-r07-1 | advisory | DAG edge completeness | `mission.md:274-292` (missing `M-09⤳M-02`) |
| F-r07-2 | advisory | six-axis ownership (root.go) | `mission.md:258`, `mission.md:262` |
| F-r07-3 | **blocking** | contract satisfiability / honest-unknown | `research/sync-health-interface-contract.md:106-107`; `M-09/F-02:35,42-45` |
| F-r07-4 | advisory | anchor accuracy | `research/codebase-ingest-side.md:96` |
| F-r07-5 | advisory | anchor accuracy (deletion range) | `M-07-pricing-fee-read/milestone.md:33,70` |
| F-r07-6 | advisory | required brief sections | `M-09/F-02` |

Also observed, not raised as a finding but worth the fold's attention: `M-02-sync-core-seam/F-01-core-ddl/feature.md` Handoff still declares "Blockers or open decisions: campos fiscais dependem de p5-prerequisites §2 (investigador em curso no planning — estará no repo antes do dispatch)", while `research/p5-prerequisites.md` §2 is manifested and complete (adapter entry `buyer_fiscal_reader.go:59`, two-step flow `:71-107`, decode structs `:18-52`, domain DTO `buyer_fiscal.go:16-22`). The blocker is discharged; the declaration is stale.

---

## Files read

Mission root `C:\Users\leandro.theodoro\Documents\marketplace-central\.mnfs\MIS-007-ml-sync`; repo root `C:\Users\leandro.theodoro\Documents\marketplace-central`. Paths below are relative to those roots, as required by the output contract.

**Manifest + PART A context (mission root)**
- `planning-reviews/p5-input-r07.sha256`
- `planning-reviews/p5-reconciliation-r06.md`

**Audited set — all 46 manifested files (mission root)**
- `mission.md`
- `research/channel-fees-interface-contract.md` (IC-01)
- `research/divergences-interface-contract.md` (IC-02)
- `research/orders-persistence-interface-contract.md` (IC-03)
- `research/webhook-inbox-interface-contract.md` (IC-04)
- `research/sync-health-interface-contract.md` (IC-05)
- `research/sync-ingest-ports-interface-contract.md` (IC-06)
- `research/listings-sync-interface-contract.md` (IC-07)
- `research/p5-prerequisites.md`
- `research/codebase-ingest-side.md`
- `research/codebase-read-side.md`
- `research/external-ml-api-facts.md`
- `M-01-ml-client-hardening/milestone.md`, `M-01-ml-client-hardening/F-01-resilience-decorator/feature.md`, `M-01-ml-client-hardening/F-02-items-multiget-raw-dto/feature.md`
- `M-02-sync-core-seam/milestone.md`, `M-02-sync-core-seam/F-01-core-ddl/feature.md`, `M-02-sync-core-seam/F-02-fee-divergence-ports/feature.md`, `M-02-sync-core-seam/F-03-scheduler-incremental-cursor/feature.md`, `M-02-sync-core-seam/F-04-read-guard-allowlist/feature.md`
- `M-03-orders-shipment-persist/milestone.md`, `M-03-orders-shipment-persist/F-01-ml-ingest-readers/feature.md`, `M-03-orders-shipment-persist/F-02-ingest-order-v1/feature.md`, `M-03-orders-shipment-persist/F-03-read-path-switch/feature.md`
- `M-04-listings-backfill-ingest/milestone.md`, `M-04-listings-backfill-ingest/F-01-listings-ddl/feature.md`, `M-04-listings-backfill-ingest/F-02-mass-closure-replacement/feature.md`, `M-04-listings-backfill-ingest/F-03-backfill-cursor-ingest/feature.md`, `M-04-listings-backfill-ingest/F-04-scheduler-refresh-wiring/feature.md`
- `M-05-listings-fees-divergence/milestone.md`, `M-05-listings-fees-divergence/F-01-camada2-fee-ingest/feature.md`, `M-05-listings-fees-divergence/F-02-stock-divergence/feature.md`, `M-05-listings-fees-divergence/F-03-anuncios-fe-contract/feature.md`
- `M-06-orders-backfill-decomposition/milestone.md`, `M-06-orders-backfill-decomposition/F-01-backfill-incremental/feature.md`, `M-06-orders-backfill-decomposition/F-02-decomposition-camada3/feature.md`, `M-06-orders-backfill-decomposition/F-03-audit-fe-pedidos/feature.md`
- `M-07-pricing-fee-read/milestone.md`, `M-07-pricing-fee-read/F-01-fee-read-resolver/feature.md`, `M-07-pricing-fee-read/F-02-precos-provenance-fe/feature.md`
- `M-08-webhook-ingest/milestone.md`, `M-08-webhook-ingest/F-01-inbox-endpoint/feature.md`, `M-08-webhook-ingest/F-02-worker-callback/feature.md`
- `M-09-sync-observability/milestone.md`, `M-09-sync-observability/F-01-sync-health-endpoint/feature.md`, `M-09-sync-observability/F-02-integracoes-health-section/feature.md`

**Repo files read for prerequisite/anchor verification (repo root)**
- `apps/server_core/internal/composition/root.go`
- `apps/server_core/internal/modules/internal_read/adapters/mirror/writer.go`
- `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go`
- `apps/server_core/internal/modules/connectors/adapters/mercado_livre/buyer_fiscal_reader.go` (existence)
- `apps/server_core/internal/modules/orders/domain/order_bucket.go`
- `apps/server_core/internal/modules/orders/application/enrich_service.go`
- `apps/server_core/internal/modules/orders/ports/buyer_fiscal_reader.go` (existence)
- `apps/server_core/internal/modules/listings/adapters/postgres/repository.go`
- `apps/server_core/internal/modules/listings/transport/query.go`
- `apps/server_core/internal/modules/sync/adapters/postgres/sync_state_repo.go`
- `apps/server_core/internal/modules/sync/application/scheduler.go`
- `apps/server_core/internal/modules/sync/composition/scheduler.go`
- `apps/server_core/internal/modules/integrations/adapters/postgres/installation_repo.go`
- `apps/server_core/internal/modules/integrations/transport/run_read_handler.go`
- `apps/server_core/internal/platform/httpx/route_deadline.go`
- `apps/server_core/migrations/listings_test.go`
- `apps/server_core/migrations/` (directory listing — tip = `0085_erp_import_products_sellable_fields.sql`)
- `packages/sdk-runtime/src/index.ts`
- `apps/web/src/pages/integracoes/IntegracoesPage.tsx`

No file was created, edited, or written during this audit.
