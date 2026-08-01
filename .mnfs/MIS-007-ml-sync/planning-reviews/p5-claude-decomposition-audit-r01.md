<!-- Verbatim output of the cold P5 decomposition auditor (Agent task a8580ea15a44ddbba),
recovered from the session transcript after context compaction. The in-context return was
not persisted at dispatch time (process defect, recorded in p5-reconciliation-r01.md). -->

# P5 Cold Decomposition Audit — MIS-007-ml-sync — round 01

**Auditor:** cold Claude Opus crew (substituting GPT-5.6 Sol P5 touchpoint under the operator-ratified waiver; Sol retroactive before `status: planned`).
**Manifest:** `planning-reviews/p5-input-r01.sha256` — digest `30d48ee15cacfe2650a025870c0afa596dcc5ac9a2793c9eaaef07fd48dd5d51` (file's own SHA-256 matches the stated digest). `sha256sum -c` over the manifest: **46 OK, 0 failures**. Complete manifest inspected — no sampling.
**Base:** repo tip `dd89d4b3`, read-only. Zero writes to any manifest file or repo file.

## Files read (46/46 manifest entries + 1 non-manifest input)

`mission.md`
`M-01-ml-client-hardening/milestone.md`, `.../F-01-resilience-decorator/feature.md`, `.../F-02-items-multiget-raw-dto/feature.md`
`M-02-sync-core-seam/milestone.md`, `.../F-01-core-ddl/feature.md`, `.../F-02-fee-divergence-ports/feature.md`, `.../F-03-scheduler-incremental-cursor/feature.md`, `.../F-04-read-guard-allowlist/feature.md`
`M-03-orders-shipment-persist/milestone.md`, `.../F-01-ml-ingest-readers/feature.md`, `.../F-02-ingest-order-v1/feature.md`, `.../F-03-read-path-switch/feature.md`
`M-04-listings-backfill-ingest/milestone.md`, `.../F-01-listings-ddl/feature.md`, `.../F-02-mass-closure-replacement/feature.md`, `.../F-03-backfill-cursor-ingest/feature.md`, `.../F-04-scheduler-refresh-wiring/feature.md`
`M-05-listings-fees-divergence/milestone.md`, `.../F-01-camada2-fee-ingest/feature.md`, `.../F-02-stock-divergence/feature.md`, `.../F-03-anuncios-fe-contract/feature.md`
`M-06-orders-backfill-decomposition/milestone.md`, `.../F-01-backfill-incremental/feature.md`, `.../F-02-decomposition-camada3/feature.md`, `.../F-03-audit-fe-pedidos/feature.md`
`M-07-pricing-fee-read/milestone.md`, `.../F-01-fee-read-resolver/feature.md`, `.../F-02-precos-provenance-fe/feature.md`
`M-08-webhook-ingest/milestone.md`, `.../F-01-inbox-endpoint/feature.md`, `.../F-02-worker-callback/feature.md`
`M-09-sync-observability/milestone.md`, `.../F-01-sync-health-endpoint/feature.md`, `.../F-02-integracoes-health-section/feature.md`
`research/channel-fees-interface-contract.md`, `research/divergences-interface-contract.md`, `research/orders-persistence-interface-contract.md`, `research/webhook-inbox-interface-contract.md`, `research/sync-health-interface-contract.md`, `research/sync-ingest-ports-interface-contract.md`, `research/listings-sync-interface-contract.md`, `research/codebase-ingest-side.md`, `research/codebase-read-side.md`, `research/external-ml-api-facts.md`, `research/p5-prerequisites.md`
Non-manifest (audited per reading order step 6, not accepted): `planning-reviews/p5-passes-r01.md`

**Repo anchors spot-checked (18, requirement ≥8):** `capability_adapter.go:79-92`, `:444`, `:462`, `:578-579`, `:654-655`, `:681`, `:712`; `root.go:259-272`, `:661-678`, `:835-858`; `orders/domain/order_bucket.go:48`; `listings/domain/filter.go:9`; `listings/ports/cursor.go:60`; `listings/adapters/connectors/source.go:54-90`; `orders/application/enrich_service.go:~192`; `integrations/adapters/mercadolivre/auth_adapter.go:47`; `sync/adapters/postgres/sync_state_repo.go:35`; `pricing/ports/tariff_resolver.go:14-32`; `migrations/0075_sync_sync_state.sql`; `platform/migrate/runner.go`; `migrations/` highest = `0085`. All cited **line numbers** are accurate. Migration out-of-order hazard investigated and **retired**: `runner.go` loads every applied filename into a set and applies anything absent (`sort.Strings` over pending only, `schema_migrations(filename TEXT PRIMARY KEY)`), so M-08's `0093` landing after M-06's `0094-0095` is safe — not a finding.

---

## FINDINGS

### F-1 — BLOCKING — check 5 (prerequisite existence) + check 2 (canonical disjointness)
**Nine ownership locks address Go packages that do not exist in the tree.**
`M-01-ml-client-hardening/F-01-resilience-decorator/feature.md:75`:
&gt; `- Owned paths: \`apps/server_core/internal/integrations/adapters/connectors/mercadolivre/\``

Same string at `M-01/F-02/feature.md:76` and `M-01/milestone.md:47` (`/**`). Also: `M-02/F-03/feature.md:72` `apps/server_core/internal/sync/application/scheduler.go`; `M-02/milestone.md:48` `internal/channelfees/`, `internal/divergences/`; `M-04/F-02/feature.md:65` `apps/server_core/internal/listings/adapters/postgres/repository.go`; `M-04/F-03/feature.md:75` `apps/server_core/internal/listings/application/**`; `M-04/milestone.md:50` `apps/server_core/internal/listings/**`; `mission.md:277` ownership cell `connectors/mercadolivre/**`.
**Offending token:** the segment `internal/&lt;module&gt;` (missing `modules/`), and `integrations/adapters/connectors/mercadolivre` as a directory.
`apps/server_core/internal/` contains exactly `{composition, modules, platform, testsupport}`. Real paths: `internal/modules/connectors/adapters/mercado_livre/` (dir `mercado_livre`, package `mercadolivre` — `research/codebase-ingest-side.md:7` states this correctly), `internal/modules/sync/application/scheduler.go`, `internal/modules/listings/**`. This is not cosmetic: the nearest **real** directory to M-01's exclusive-surface string is `internal/modules/integrations/adapters/mercadolivre/` — the OAuth `auth_adapter.go` package, whose `:47-48` `baseline_commission_percent: 0.16` is explicitly **M-07's** to delete (`mission.md:283`, `M-07/milestone.md:26`). Read literally, M-01's blanket lock claims a package M-07 owns; read charitably, M-01 has no enforceable lock at all. Ownership locks are the only instrument the collision matrix has, and they currently point at nothing.
**Yes-if:** all nine loci are rewritten to the canonical `apps/server_core/internal/modules/...` paths (M-01 → `internal/modules/connectors/adapters/mercado_livre/`, explicitly NOT `internal/modules/integrations/adapters/mercadolivre/`), and the two new core packages of `M-02/milestone.md:48` are placed under `internal/modules/&lt;module&gt;/` per the layering doctrine in AGENTS.md rather than at `internal/&lt;name&gt;/`.

### F-2 — BLOCKING — check 6 (verbatim propagation)
**Six briefs, including two published-DTO surfaces, name a `channel_fees` provenance column that IC-01 does not define.**
`research/channel-fees-interface-contract.md:63`:
&gt; `` `coletado_em timestamptz` NOT NULL (hora do NOSSO fetch) ``

and `:68-69` — `ResolveListingFees` returns `{value, value_type, currency, layer, origem, coletado_em}`. Consumers instead write `observed_at`: `M-05/F-01:28` (`` `observed_at`=hora da chamada ``), `:37`, `:41`, `:71`, `:88`; `M-05/F-03:26` and `:39` (**FE/SDK DTO** — `{percent, origem, observed_at}`); `M-05/milestone.md:68`, `:76`; `M-07/F-01:31`, `:43`, `:80`; `M-07/F-02:27`, `:31`, `:62` (**FE/SDK DTO** — `{"origem":"api_listing_prices","observed_at":"..."}`); `M-07/milestone.md:27`, `:80`.
**Offending token:** `observed_at` applied to `channel_fees`. IC-02 legitimately uses `*_observed_at` for `divergences`, which is exactly how the name leaked; but for the fee ledger the column is `coletado_em`, and two milestones would ship the wrong name into OpenAPI + the hand-written SDK under the same-commit contract lock.
**Yes-if:** every `channel_fees` reference is renamed to `coletado_em` (or IC-01 is amended and re-ratified to rename the column, and both FE DTOs quote whichever name wins verbatim).

### F-3 — BLOCKING — check 6 (verbatim propagation)
**M-05 pins layer-2 as `percent`; IC-01's canonical layer-2 example pins `amount`.**
`M-05/F-01/feature.md:27`:
&gt; `subject_id=provider_listing_id, layer=2, fee_kind=\`commission\`, value_type=\`percent\`,`

`research/channel-fees-interface-contract.md:106-110` (Camada 2 comissão, **produtor M-05**):
&gt; `"value_type":"amount","value":16.45,"currency":"BRL","origem":"api_listing_prices", "detail":{"percentage_fee":12.5,"fixed_fee":6.0,...,"price_used":79.90,...}`

**Offending token:** `value_type=\`percent\`` at M-05 F-01:27 (propagated to `M-05/F-03:26` as a `percent` string in the FE DTO). The CHECK at IC-01:57 admits both, so this will not fail at the DDL — it fails at the reader: M-07's cascade (`M-07/F-01:23-31`) consumes the same rows, and the M-02 writer/reader contract tests will encode whichever shape the canonical example states. `amount` also carries the fixed-fee component (`fixed_fee":6.0`) that a bare percentage silently drops.
**Yes-if:** one representation is pinned for `(layer=2, fee_kind=commission)` across IC-01's canonical example, M-05 F-01, M-05 F-03 and M-07 F-01/F-02 — and if `percent` wins, IC-01's example and the `detail` contract are amended so the fixed-fee component is not lost.

### F-4 — BLOCKING — check 6 (verbatim propagation) + check 4 (contract satisfiability)
**M-06's sum-of-parts invariant is arithmetically false against IC-03's canonical decomposition.**
`M-06/F-02/feature.md:66-67`:
&gt; `- Soma-das-partes é invariante testada: liquido = receita_bruta − comissao_total − frete_seller (tolerância centavo por arredondamento ...)`

Same formula in the EARS at `:40-41` ("`decomposição com liquido = receita − comissao − frete`") and echoed at `M-06/milestone.md:80`.
`research/orders-persistence-interface-contract.md:113-115`:
&gt; `"receita_bruta":239.70,"comissao_total":48.66,"frete_seller":22.90, "custo_produto":95.10,... "liquido":73.04`

239.70 − 48.66 − 22.90 = **168.14**, not 73.04. IC-03 subtracts `custo_produto`; IC-03:105 further binds "`liquido`/`margem_pct` só computados com TODOS os insumos presentes".
**Offending token:** the invariant string `liquido = receita_bruta − comissao_total − frete_seller`. This is the *tested* invariant — a green test would certify a wrong number on `/pedidos`, which is precisely the failure class R-4 exists to catch.
**Yes-if:** the invariant reads `liquido = receita_bruta − comissao_total − frete_seller − custo_produto`, with the `incompleto[]` branch (custo absent ⇒ no `liquido`, no `margem_pct`) kept as the only exception, matching IC-03:105 and the canonical example at `:113-123`.

### F-5 — BLOCKING — check 6 (verbatim propagation)
**`GET /sync/health` requires a parameter IC-05 forbids, and cites IC-05 as the authority for it.**
`M-09/F-01/feature.md:32`:
&gt; `Param \`installation_id\` obrigatório (idioma de /sync/runs).`

reinforced at `:69` ("installation_id ausente = 400 com code") and `:73`:
&gt; `- installation_id inexistente → entities todas null-honestas (não 404 ...; IC-05 pina).`

`research/sync-health-interface-contract.md:44-46`:
&gt; `### Required Inputs` / `Nenhum parâmetro além do tenant do contexto.`

and `:40` — `GetSyncHealth ... | tenant (ctx) |`; `:110` — "Consultas: `sync_state` scan por tenant".
**Offending token:** `Param \`installation_id\` obrigatório`, plus the false attribution `IC-05 pina` at `:73` (IC-05 pins the opposite). The FE consumer `M-09/F-02` renders one section for the tenant; a required per-installation param changes the OpenAPI path, the SDK signature and the query shape — all inside the hub-serialized contract seam.
**Yes-if:** M-09 F-01 drops the parameter and reads by tenant scan per IC-05:44-46 and :110 (the repo already scopes tenant in the constructor — `sync_state_repo.go:30`, `SyncStateRepository{pool, tenantID}`), or IC-05 is amended to make `installation_id` required and the canonical example + `M-09/F-02` are updated with it.

### F-6 — BLOCKING — check 6 (verbatim propagation)
**The webhook block's pre-M-08 shape is stated three incompatible ways; the golden fixture would encode the wrong one.**
`M-09/F-01/feature.md:33-35`:
&gt; `Bloco webhook via porta \`WebhookStatsReader\` ... com implementação NULA (retorna null-stats)` ... `Campos sem observação = null JSON (nunca 0, nunca string vazia).`

and `:88` — "Teste da porta: impl nula → bloco webhook null"; `:87` — "JSON golden com nulls exatos".
`research/sync-health-interface-contract.md:105` (canonical honest initial state):
&gt; `"webhook":{"last_notification_at":null,"pending":0,"dropped_24h":0}`

with `:74` hedging both ways — "`webhook`: NULL/zeros quando inbox ainda não existe".
**Offending token:** `null-stats` / "nunca 0" against the canonical `"pending":0,"dropped_24h":0`. IC-05 also binds "NULL viaja como null literal, nunca omitido (shape estável p/ SDK tipado)" — a whole-null block and a zero-filled block are different SDK types, and M-08 F-02 later swaps the implementation under that same type.
**Yes-if:** IC-05:74 is de-hedged to a single shape (the canonical example's — `last_notification_at: null`, counters `0`, since 0 pending rows in a non-existent inbox is a true count, not a fabricated unknown), and `M-09/F-01:34`, `:87-88` quote it verbatim.

### F-7 — BLOCKING — check 1 (edge completeness) + check 3 (write-DAG)
**The M-08↔M-09 port dependency has no edge in the DAG, and the wiring swap crosses into M-09's anchored root.go region. The author's Pass 1 clears this on a claim that is false.**
`mission.md:259-271` edge table lists 11 edges; none is M-08→M-09. The DAG block at `:252-257` shows `M-09` free-floating with only the soft `⤳ M-04/M-06`.
`M-08/F-02/feature.md:51-53`:
&gt; `implementação REAL de \`WebhookStatsReader\` (porta IC-05 publicada pelo M-09 com impl nula) ... troca da fiação na região ancorada do M-08 (a impl nula do M-09 morre aqui).`

`planning-reviews/p5-passes-r01.md:23-25`:
&gt; `M-08 ↔ M-09 → \`WebhookStatsReader\`: M-09 publica porta + impl nula (F-01); M-08 fornece impl real + troca fiação NA REGIÃO DELE (F-02). ... Sem edit cruzado de arquivo.`

**Offending tokens:** the absent edge row, and `Sem edit cruzado de arquivo`. Both milestones edit `root.go` (matrix `mission.md:284-285`: M-08 "1 linha ancorada", M-09 "1 linha ancorada"). The constructor that *receives* `WebhookStatsReader` is M-09's health-service constructor on M-09's anchored line (`M-09/F-01:35` — "Registro: `httpx.InteractiveRouteClass` + 1 linha ancorada root.go"). Replacing the null implementation means editing the argument list on **M-09's** line; no indirection (registry, setter, late-bind) is specified in either brief. The refuted claim is the sole basis on which Pass 1 declared the seam clean.
**Yes-if:** a hard edge `M-08→M-09 | porta WebhookStatsReader (IC-05) + linha de fiação` is added to `mission.md:259-271`, AND either M-09 F-01 specifies the indirection that lets M-08 supply the real implementation without touching M-09's line, or the matrix registers root.go's health-wiring line as an additive lock transferred to M-08 with the hub as resolver of record (R-7).

### F-8 — BLOCKING — check 4 (contract satisfiability)
**Camada-2 freight is promised by the milestone strategy and exposed in a published FE vocabulary, but every producer brief explicitly excludes it.**
`mission.md:238` (M-05 headline):
&gt; `Camada 2 (\`listing_prices\` com category_id/price ingeridos + shipping_options), divergência de estoque no ingest, ...`

`M-05/F-01/feature.md:31-32` (the only camada-2 producer):
&gt; `Frete: NÃO observado aqui (honest-unknown IC-01 — sem chute de R$79).`

`research/channel-fees-interface-contract.md:62` keeps `api_shipping_options` in the `origem` CHECK, and `M-07/F-02/feature.md:26-27` ships it in the FE vocabulary:
&gt; `ADITIVO por componente de tarifa \`origem\` (vocabulário IC-01: api_listing_prices | api_shipping_options | config)`

**Offending token:** `+ shipping_options` at mission.md:238 against `Frete: NÃO observado aqui`. No brief in M-01..M-09 emits an `origem='api_shipping_options'` row. The FE would carry a rendering branch no producer can reach — the "unproved branch" class flagged in the MIS-006 postmortem.
**Yes-if:** either the mission M-05 row drops `+ shipping_options` and M-07 F-02's vocabulary drops `api_shipping_options` for this mission (IC-01's CHECK may keep it as an additive-future value with a one-line note), or a named feature is added that produces it — with an owner, a range, and a spot in the lane plan.

### F-9 — BLOCKING — check 2 (disjointness) + check 3 (write-DAG)
**M-05 F-01's fallback path writes inside the ML adapter package, which M-05 does not own and which no edge or lock grants it.**
`M-05/F-01/feature.md:29-31`:
&gt; `se o payload de multiget do M-01 F-02 já carrega, ZERO GET extra; senão, GET /items/{id}/prices em lote respeitando token-bucket.`

M-05's ownership cell (`mission.md:281`) is `ingest ext camada2/divergência (dentro de listings app via lock aditivo do M-04 PÓS-close), FE /anuncios` — no ML adapter surface. A new ML endpoint is, per ADR-02/A-02 and `mission.md:277`, a new file in the ML package owned by M-01 (lane A, long closed by the time lane C runs). Pass 3 (`p5-passes-r01.md:92-96`) rules this "não é ★5: nenhum caminho depende de símbolo inexistente" — true about symbols, silent about ownership: path B requires a producer M-05 cannot write and M-01 has not been asked to build.
**Offending token:** `senão, GET /items/{id}/prices em lote`.
**Yes-if:** the multiget field question is resolved during M-01 F-02 (the milestone that already reads `/items`) and its outcome recorded, so M-05 F-01 has exactly one path; or path B is assigned to M-01 F-02 as a named deliverable (new file in `internal/modules/connectors/adapters/mercado_livre/`) with the M-05→M-01 edge added.

### F-10 — BLOCKING — check 4 (contract satisfiability)
**The config fallback has two owners, which makes M-07 F-01's honest-absent branch unreachable, and puts a pricing-owned table behind an M-02 core port.**
`M-02/F-02/feature.md:26-28`:
&gt; `\`ChannelFeeReader.ResolveListingFees\` (IC-01): ... resolução comissão camada 2 → 1 → fallback \`pricing_tariff_defaults\` com proveniência \`config\``

with Inputs at `:45-46` citing `pricing_tariff_defaults` leitura existente (`calc_repository.go:240-246`).
`M-07/F-01/feature.md:26-30`:
&gt; `resolve comissão via ChannelFeeReader (IC-01) em cascata camada 2 ... → camada 1 ... → fallback \`pricing_tariff_defaults\` (resolver existente \`pricingtariffdefaults.NewResolver\` root.go:837 REUSADO como base da cadeia — degrau \`config\`)`

and its third EARS at `:46-48`:
&gt; `- While nenhuma fonte resolve (nem defaults), when Resolve roda, the resultado shall ser erro/honest-absent — NUNCA 0.16 nem qualquer constante.`

**Offending token:** `fallback \`pricing_tariff_defaults\`` appearing in both briefs. The defaults are materialize-on-read (`calc_repository.go:240-246`, `INSERT ... ON CONFLICT DO NOTHING`, 13.00/16.00 — verified) so the M-02 reader always answers; the M-07 EARS-3 branch can then never fire, and a test for it is unfalsifiable — the criterion-vacuous-against-the-type failure class. Secondarily, M-02's ownership cell (`mission.md:278`) carries no pricing surface, yet its core fee package would read a pricing table and M-07 owns `pricing/**`.
**Yes-if:** the config step is assigned to exactly one owner — recommended M-07 F-01, which already composes `pricingtariffdefaults.NewResolver` as the chain base — and `M-02/F-02:26-28` is narrowed to "camada 2 → camada 1 → honest-absent, nunca config" (which its own first EARS at `:34-35` already implies: "nunca o config"), with the `:36-37` EARS and Inputs `:45-46` updated accordingly, so M-07 F-01 EARS-3 becomes reachable and testable.

---

### ADVISORY

**A-1 — check 6.** Migration sub-range drift: `mission.md:235` says M-02 owns "DDL 0086-0088" and `research/channel-fees-interface-contract.md:126` says "migração no range do M-02 (0086-0088)", while `mission.md:278`, `M-02/milestone.md:51-52` and `M-02/F-01:25` all say **0086-0089** (0089 = additive `orders` columns). Allocation is otherwise disjoint and correct (0086-0089 / 0090-0092 / 0093 / 0094-0095; highest existing 0085 verified). *Yes-if:* mission.md:235 and IC-01:126 read 0086-0089.

**A-2 — check 5.** `mission.md` and `M-02/F-02:45-46` cite the upsert idiom as `internal_read/adapters/mirror/writer.go:74-95`; the keep-absent half (`keepAbsentSQL`, the ADR-04 "never DELETE" clause) begins at `:97`. Verified range is `74-105` (`research/codebase-ingest-side.md:95-96` gets this right). *Yes-if:* the anchor reads `74-105`.

**A-3 — check 1.** Edge-table asymmetry: `M-06→M-02` is listed even though it is transitively implied via M-03, but `M-05→M-02` is absent although M-05 writes `channel_fees` and `divergences` (M-02's tables and ports). Separately, `mission.md:270` names the M-08→M-06 forcing artifact as "worker chama IngestOrder", and `M-08/F-02:28` attributes `IngestOrder` to "(M-06 ...)", while IC-06:17-18 and Pass 3 assign creation to **M-03 F-02** (M-06 only extends it) — the artifact that actually forces M-08 after M-06 is the incremental/decomposition extension, unnamed. *Yes-if:* the M-05→M-02 row is added and the M-08→M-06 artifact is renamed to the extension it depends on.

**A-4 — check 2.** `mission.md:277` gives M-01 `connectors/mercadolivre/**` with the parenthetical "(ÚNICO a editar capability_adapter.go)", while `M-03/F-01:26,76` places new files in the same package. The lanes serialize (M-03→M-01), so there is no simultaneity, but the blanket `**` contradicts its own parenthetical. *Yes-if:* the cell reads `connectors/adapters/mercado_livre/` with exclusivity scoped to existing files, new files explicitly permitted to downstream milestones.

**A-5 — check 7.** `## Interaction Model` exists only in `M-09/F-02`. Absent from the three other FE workflow features where state drift is plausible: `M-05/F-03` (filter `divergentes=true` interacting with the existing cursor/sort and the badge), `M-06/F-03` (drawer + Fila/Kanban margin surfaces), `M-07/F-02` (⚠ chip states across single vs batch simulation). *Yes-if:* each gains the section, or states why no client state is introduced.

**A-6 — check 7.** `## Inputs/Outputs` absent from `M-01/F-01`, `M-02/F-04`, `M-03/F-03`, `M-04/F-01`, `M-04/F-02`, `M-04/F-04`, `M-08/F-02`. Most are internal seams where its absence is defensible, but two are asymmetric with their siblings: `M-04/F-01` is a DDL brief without it while `M-02/F-01` (also DDL) has it; `M-03/F-03` switches what `GET /orders` serves. All 25 briefs do have `## Negative Scenarios`. *Yes-if:* the two asymmetric briefs gain the section.

**A-7 — check 8.** `M-07/F-02:31-32` introduces a product rule not present in any ratified ADR or IC: `observed_at stale (&gt; 7 dias — limiar pinado aqui, planning decision)`. P5 may pin thresholds, but a user-visible staleness warning is new product surface relative to the ratified set. *Yes-if:* the threshold is recorded in IC-01 or the mission's clarified decisions, or the ⚠ is limited to `origem=config` for this mission.

**A-8 — check 5.** `research/codebase-read-side.md:162` states "migrações vão até 0081"; the tree's highest is **0085** (`0085_erp_import_products_sellable_fields.sql`, verified) — as `research/codebase-ingest-side.md:47` and `p5-prerequisites.md` item 11 both correctly report. Harmless to the allocation (0086+ is free either way) but it is a false fact in a manifest input. *Yes-if:* corrected to 0085.

**A-9 — check 6.** `M-06/F-02:29-30` specifies the layer-3 freight row as `fee_kind=\`freight\`, origem=\`api_shipment\`` without pinning `subject_type`; IC-01:50 admits `order` and `order_line`, and IC-01:52-55 pins a different `subject_id` format for each. *Yes-if:* the brief states `subject_type=order` with `subject_id=&lt;provider_order_id&gt;` (the natural grain for a per-shipment cost).

---

## Audit of `planning-reviews/p5-passes-r01.md` (author's self-declared 3/3 PASS)

- **Pass 1 item 4 is refuted** — see F-7. "Sem edit cruzado de arquivo" is false; both milestones edit `root.go` and M-08 must edit M-09's anchored line absent an indirection no brief specifies. Pass 1's PASS verdict rests on it.
- **Pass 2 item 2 rubber-stamps a false invariant** — it declares the sum-of-parts/`incompleto[]` interaction SATISFIABLE without checking the formula itself, which is arithmetically wrong against IC-03's own canonical example (F-4).
- **Pass 3 verified symbols but never verified owned paths** — every "fato #N" row is accurate against the code, yet nine ownership locks address directories that do not exist (F-1). Prerequisite-existence for a decomposition audit must cover the lock targets, not only the referenced symbols.
- **Pass 3's "Assunção ABERTA" disposition is under-scoped** — the two-path camera-2 fee source is not a symbol-existence problem, it is an ownership problem (F-9).
- Correctly handled and confirmed by this audit: the IC-06 in-round amendment (ports in the owning module) is present in `research/sync-ingest-ports-interface-contract.md:29-36`; the 0075 semantic-enum claim is exact (verified in the migration and in `sync/domain`); the migration ranges are disjoint; the "no migration needed to register listings/orders" claim is correct.

## Facts re-verified against the binding context (no contradiction found)

`sale_fee` **per unit** — `mission.md:184`, IC-01:117-118 (`sale_fee_unit 8.11 × qty 3 = 24.33`), `M-06/F-02:29`, `external-ml-api-facts.md:11`: consistent everywhere. `incremental=false` hardcoded at `scheduler.go:160` — IC-06:106-108 and M-02 F-03 name it as M-02's single fix: consistent. Nil-cursor erasure at `:42-45` — IC-06:79-81 forbids nil returns and requires a valid terminal cursor: consistent. `ApplyCompletedPull` mass-closure at `repository.go:390-394` — M-04 F-02 targets exactly that statement: consistent. Semantic entity enum without DB CHECK — IC-06:104-105 and Pass 3: consistent.

---

## VERDICT

**NEEDS-REVISION** — 10 BLOCKING findings (F-1 … F-10), 9 ADVISORY.

Blocking set: F-1 nonexistent ownership paths; F-2 `observed_at` vs `coletado_em`; F-3 layer-2 `percent` vs `amount`; F-4 false `liquido` invariant; F-5 `/sync/health` `installation_id` vs IC-05 (with a false IC citation); F-6 webhook block shape; F-7 missing M-08→M-09 edge + crossed root.go region (refutes Pass 1); F-8 unproduced `api_shipping_options`; F-9 M-05 fee-source path B in an unowned package; F-10 doubly-owned config fallback making M-07's honest-absent branch unreachable.

Input set was complete and verifiable (46/46), so this verdict is not sampled. Re-audit required after revision; the Sol P5 retroactive touchpoint remains mandatory before `status: planned`.
