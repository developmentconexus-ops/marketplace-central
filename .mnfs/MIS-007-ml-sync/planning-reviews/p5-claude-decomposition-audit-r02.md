<!-- provenance: Agent task a2e8de99471e5fb45 ("P5 decomposition re-audit r02", cold Opus,
     read-only). Task output-file was 0 bytes; result arrived in the task-notification and was
     persisted VERBATIM in the same turn it arrived (rule from p5-reconciliation-r01.md).
     Transport HTML entities (&lt; &gt; &amp;) decoded to their literal characters. -->

# P5 Decomposition Audit — round 02 (MIS-007-ml-sync)

Cold read-only auditor, operator-ratified waiver substituting the GPT-5.6 Sol P5 touchpoint. Sol retroactive pass remains MANDATORY before `status: planned` regardless of this verdict. Nothing was edited, written, or created.

---

## 1. Header

**Manifest self-digest verified: YES.**
`sha256sum planning-reviews/p5-input-r02.sha256` → `fc0522e1665949cf573b72e5ee855105ceef5b14d3f0564cde00dabd3ebb1c6c *planning-reviews/p5-input-r02.sha256` — byte-identical to the expected digest.

**`sha256sum -c planning-reviews/p5-input-r02.sha256` (run from mission root, Git Bash): 46 OK / 0 FAILED.**

**Files read: ALL 46 (no sampling).**

M-01: `F-01-resilience-decorator/feature.md`, `F-02-items-multiget-raw-dto/feature.md`, `milestone.md` · M-02: `F-01-core-ddl/feature.md`, `F-02-fee-divergence-ports/feature.md`, `F-03-scheduler-incremental-cursor/feature.md`, `F-04-read-guard-allowlist/feature.md`, `milestone.md` · M-03: `F-01-ml-ingest-readers/feature.md`, `F-02-ingest-order-v1/feature.md`, `F-03-read-path-switch/feature.md`, `milestone.md` · M-04: `F-01-listings-ddl/feature.md`, `F-02-mass-closure-replacement/feature.md`, `F-03-backfill-cursor-ingest/feature.md`, `F-04-scheduler-refresh-wiring/feature.md`, `milestone.md` · M-05: `F-01-camada2-fee-ingest/feature.md`, `F-02-stock-divergence/feature.md`, `F-03-anuncios-fe-contract/feature.md`, `milestone.md` · M-06: `F-01-backfill-incremental/feature.md`, `F-02-decomposition-camada3/feature.md`, `F-03-audit-fe-pedidos/feature.md`, `milestone.md` · M-07: `F-01-fee-read-resolver/feature.md`, `F-02-precos-provenance-fe/feature.md`, `milestone.md` · M-08: `F-01-inbox-endpoint/feature.md`, `F-02-worker-callback/feature.md`, `milestone.md` · M-09: `F-01-sync-health-endpoint/feature.md`, `F-02-integracoes-health-section/feature.md`, `milestone.md` · `mission.md` · research: `channel-fees-interface-contract.md`, `codebase-ingest-side.md`, `codebase-read-side.md`, `divergences-interface-contract.md`, `external-ml-api-facts.md`, `listings-sync-interface-contract.md`, `orders-persistence-interface-contract.md`, `p5-prerequisites.md`, `sync-health-interface-contract.md`, `sync-ingest-ports-interface-contract.md`, `webhook-inbox-interface-contract.md`.

Prior-round context read but **not** audited: `planning-reviews/p5-claude-decomposition-audit-r01.md`, `planning-reviews/p5-reconciliation-r01.md`, `planning-reviews/p5-passes-r01.md`, `planning-reviews/p3-*`.

Repo anchors verified directly (not via the plan's claims): `apps/server_core/internal/modules/` module list; highest migration `0085_erp_import_products_sellable_fields.sql`; `internal/composition/root.go:259-272` `registerBatchRoutes`; `root.go:826-858` pricing region; `internal/modules/sync/application/scheduler.go:46,160`; `sync/composition/products_job.go:90`; `pricing/ports/tariff_resolver.go:14-24,30-32`; `pricing/adapters/tariffdefaults/resolver.go`; `pricing/adapters/postgres/calc_repository.go:239-269`; `integrations/adapters/mercadolivre/auth_adapter.go:34-54`; `integrations/transport/run_read_handler.go:30-34`; `listings/transport/http_handler.go:28-30`; `listings/domain/filter.go:9`; `listings/transport/query.go:31`; `apps/web/src/pages/integracoes/IntegracoesPage.tsx:508,558-574`; `packages/sdk-runtime/src/index.ts:122,2144,2359,2361,2374`; `contracts/api/marketplace-central.openapi.yaml:4941`; `wiki/framework/provider-metadata-contract.md:41`.

---

## 2. r01 blocking findings F-1 … F-10

**F-1 — nonexistent ownership paths (9 loci). CLOSED.**
Loci checked: `M-01/milestone.md:47-53`, `M-01/F-01:75`, `M-02/milestone.md:47-52`, `M-04/milestone.md:50`, `mission.md:279`. `M-02/milestone.md:47-50`: *"packages novos de núcleo p/ fees/divergences sob `apps/server_core/internal/modules/` (ex.: `internal/modules/channelfees/`, `internal/modules/divergences/` — nomes finais no spec, layering AGENTS obriga `modules/`, DONO é este milestone)"*; `:51` `apps/server_core/internal/modules/sync/application/scheduler.go`. Verified against the tree: `internal/modules/` exists with `connectors`, `listings`, `orders`, `pricing`, `sync`, `integrations`; `sync/application/scheduler.go` exists. New packages are named as created by their owner, which is legitimate. (One residual ellipsis locus → NEW ADVISORY N-7.)

**F-2 — `observed_at` vs `coletado_em`. CLOSED.**
`research/channel-fees-interface-contract.md:63` `coletado_em timestamptz NOT NULL`. Every M-05/M-07 provenance locus now uses it: `M-05/F-01:27-34`, `M-05/F-03:26`, `M-05/milestone.md:77`, `M-07/milestone.md:26` (*"DTOs de /pricing carregam PROVENIÊNCIA (origem + coletado_em)"*), `:80`, `M-07/F-01:31` (*"origem IC-01 + coletado_em"*), `:43`, `:80`, `M-07/F-02:28,33,37,65`. Zero `observed_at` outside IC-02's distinct `expected_observed_at`/`observed_observed_at` vocabulary.

**F-3 — camada 2 percent vs amount. CLOSED.**
`IC-01:111-118` canonical `"value_type":"amount","value":16.45` + `detail{percentage_fee, fixed_fee, financing_add_on_fee, price_used, listing_type_id}`. Propagated at `M-05/F-01:27-34`, `M-05/F-03:26` (`tarifa {amount, detail, origem, coletado_em}`), and the audit formula is re-anchored at `M-06/F-03:27-28`: *"esperado_unit = detail.percentage_fee × unit_price/100 + detail.fixed_fee; esperado_total = esperado_unit × quantity"* — never the camada-2 `amount`.

**F-4 — false sum-of-parts invariant. CLOSED.**
`M-06/F-02:69-72`: *"Soma-das-partes é invariante testada: liquido = receita_bruta − comissao_total − frete_seller − custo_produto (IC-03 canonical; tolerância centavo por arredondamento — regra pinada na spec, half-even)"*; matches `IC-03:113-116` `239.70 − 48.66 − 22.90 − 95.10 = 73.04`. `M-06/milestone.md:80-81` carries the same form; absent part ⇒ liquido/margem ABSENT + `incompleto[]`.

**F-5 — `installation_id` param IC-05 forbids. CLOSED.**
`M-09/F-01:37-38`: *"SEM param `installation_id` — tenant vem do ctx (IC-05 pina); a leitura varre as instalações ML do tenant."* `:78` Inputs/Outputs: *"Input: NENHUM param — tenant do ctx (IC-05)."* `:81` negative scenario: *"Tenant sem instalação ML ou sem sync ainda → entities todas null-honestas + webhook canônico inicial (não 404 — IC-05 pina)."* Matches `IC-05:60` *"Nenhum parâmetro além do tenant do contexto."*

**F-6 — webhook block stated 3 incompatible ways. CLOSED.**
`IC-05:88-92` de-hedged to the exact canonical state. Quoted verbatim at `M-09/F-01:32` and `:47` (`{"last_notification_at":null,"pending":0,"dropped_24h":0}`), `M-09/milestone.md:25-27`, and the FE discriminator at `M-09/F-02:33-35`: *"discriminador de estado inicial = `last_notification_at === null` → subseção mostra 'webhook não configurado'"*, with `F-02:45` *"pending/dropped 0 do estado inicial NÃO renderizam como atividade"*.

**F-7 — missing M-08→M-09 edge + cross-region wiring swap. CLOSED.**
DAG `mission.md:254` `M-09 (porta WebhookStatsReader IC-05)` with the `↑` edge at `:253`; edge table row `mission.md:272`: *"M-08→M-09 | porta WebhookStatsReader (IC-05) + injeção `WithWebhookStatsReader` da região ancorada do M-08"*. Binding seam section added at `IC-05:36-48` including *"Código e linha de construção do M-09 NUNCA editados pelo M-08"* and *"Porta opcional ⇒ compile-time assert obrigatório (lição catalog-503…)"*. Producer side `M-09/F-01:33-38` + `M-09/milestone.md:52-55`; consumer side `M-08/F-02:50-54` (*"injeção via `WithWebhookStatsReader(...)` chamado DA REGIÃO ANCORADA DO M-08 no root.go"*); `M-08/milestone.md:43-44` Dependencies now lists M-09.

**F-8 — camada-2 freight promised with no producer. CLOSED.**
`mission.md:238` M-05 headline: *"frete = honesto-desconhecido, sem shipping_options nesta missão"*. `IC-01:65-67` marks `api_shipping_options` an additive-future reserve with no producer. `M-07/F-02:26-28`: *"vocabulário IC-01 com produtor NESTA missão: api_listing_prices | config — `api_shipping_options` fica no CHECK do schema como reserva aditiva futura, SEM produtor aqui"*.

**F-9 — M-05 fallback path writes in an ML adapter it does not own. CLOSED.**
`M-01/F-02:54-59` now carries the named deliverable: *"**Fonte de fee camada 2 (decisão de posse, auditoria P5 F-9)**: a execução VERIFICA se o multiget expõe `sale_price` com `?context=channel_marketplace` … Se NÃO expõe, este feature TAMBÉM entrega reader dedicado de prices (`GET /items/{id}/prices`, arquivo novo no mesmo package) … p/ que M-05 F-01 consuma UMA fonte pronta e nunca escreva no package do adapter (posse M-01)."* `M-05/F-01` consumes one ready source; its open-decision line is gone.

**F-10 — config fallback with two owners. CLOSED (ownership axis).**
`M-02/F-02:26-31` narrows `ChannelFeeReader` to ledger-only, confirmed by its negative scenario `:71-72`: *"Resolve p/ listing sem NENHUMA row de ledger → resultado 'desconhecido' tipado, nunca zero (fallback config é do consumidor M-07, fora deste port)."* `M-07/F-01:29-30` owns the config step: *"fallback `pricing_tariff_defaults` (resolver existente `pricingtariffdefaults.NewResolver` root.go:837 REUSADO como base da cadeia — degrau `config`)"* — verified verbatim at `root.go:837`. **Caveat:** the reconciliation's added claim *"M-07 F-01 EARS-3 (honest-absent) is now reachable"* (`p5-reconciliation-r01.md:78-79`) is **REFUTED against the code** — carried as NEW BLOCKING **N-3** below.

---

## 3. r01 advisories A-1 … A-9

- **A-1** (migration sub-range drift) — **CLOSED at both cited loci**: `mission.md:235` and `IC-01:133` now read `0086-0089`. Two *uncited sibling* loci still drift → NEW ADVISORY **N-6**.
- **A-2** (writer.go anchor) — **CLOSED**: `mission.md:95` and `M-02/F-02` Inputs both cite `writer.go:74-105`; matches `codebase-ingest-side.md:95-96`.
- **A-3** (missing `M-05→M-02` row; M-08→M-06 artifact name) — **CLOSED**: `mission.md:266` *"M-05→M-02 | ports ChannelFeeWriter (camada 2) + DivergenceRecorder (estoque) — IC-01/IC-02"*; `:271` *"worker chama o IngestOrder ESTENDIDO do M-06 (decomposição+camada 3) — não o v1 do M-03"*.
- **A-4** (M-01 matrix cell over-claims) — **CLOSED**: `mission.md:279` *"exclusividade sobre os ARQUIVOS EXISTENTES (único a editar capability_adapter.go); arquivos novos no dir permitidos downstream (M-03 F-01 readers)"*.
- **A-5** (missing Interaction Model) — **CLOSED**: `## Interaction Model` present at `M-05/F-03:71-80`, `M-06/F-03:70-79`, `M-07/F-02:68-74`.
- **A-6** (missing Inputs/Outputs on asymmetric briefs) — **CLOSED**: `M-04/F-01:52-56` and `M-03/F-03:60-65`.
- **A-7** (unratified 7-day staleness rule) — **CLOSED**: `M-07/F-02:33` *"SEM regra de staleness (limiar de stale = decisão futura, fora desta missão)"*, `:58-59` *"⚠ SÓ para origem=config — nenhuma heurística de staleness nesta missão"*. Repo-wide grep over the mission tree finds no 7-day threshold anywhere.
- **A-8** (stale migration count) — **CLOSED**: `research/codebase-read-side.md:162` reads `0085`; verified highest on disk is `0085_erp_import_products_sellable_fields.sql`.
- **A-9** (freight row subject) — **CLOSED**: `M-06/F-02:29-31` pins `subject_type=order`, `subject_id=<provider_order_id>`.

---

## 4. NEW findings (fresh full pass)

### N-1 — BLOCKING — check 4 (contract satisfiability) + check 2 (ownership)
**Locus:** `M-02-sync-core-seam/F-03-scheduler-incremental-cursor/feature.md:49-51` (Expected Output).
**Offending token:** `+ tipo de retorno do job carregando o tipo do run — mudança aditiva compatível com products`.
The same brief's Inputs at `:43` states *"IC-06 (binding: JobFunc signature intocada, …)"*, and `research/sync-ingest-ports-interface-contract.md:151` Must Preserve states *"`JobFunc` signature (scheduler.go:46)"*. Verified in code: `apps/server_core/internal/modules/sync/application/scheduler.go:46` is `type JobFunc func(ctx context.Context, cursor json.RawMessage) (json.RawMessage, error)`. A Go function type admits **no** additive return change; and `internal/modules/sync/composition/products_job.go:90` returns `syncapp.JobFunc`, so any change forces an edit there — a path the same brief lists under Forbidden paths at `:74` (*"jobs concretos"*). The brief therefore mandates a change that is simultaneously forbidden by its own ownership block, by IC-06, and by its own Inputs line. (Secondary defect: the phrase is ambiguous between "job's return type" and "run type carried inside the returned cursor" — a Must-Preserve signature must not be left ambiguous in a brief.)
**Yes-if:** Expected Output pins that the run type is derived from the mandatory `cursor.phase` (already binding at `F-03:30-31` and `IC-06:136-137`, and already the source used by `IC-05:87`), leaving `JobFunc` byte-identical — **or** IC-06 Must Preserve is amended and `sync/composition/products_job.go` is added to M-02 F-03's owned paths as a named updated call site.

### N-2 — BLOCKING — check 5 (prerequisite existence) + check 4 (contract satisfiability) + check 2 (ownership)
**Locus:** `M-07-pricing-fee-read/F-01-fee-read-resolver/feature.md:37-38`; `M-07-pricing-fee-read/milestone.md:32-33`, `:53`, `:76`; `research/channel-fees-interface-contract.md:106`.
**Offending tokens:** `0.16 de auth_adapter.go:47-48 deletado (miss ali vira erro tipado ou honest-absent — spec pina contra o call site)`; `um miss cai em 0.16 hardcoded silencioso`; `fallback silencioso 0.16`; `auth_adapter.go:47-48 (0.16) MORRE no M-07`.
Verified in code: `internal/modules/integrations/adapters/mercadolivre/auth_adapter.go:42-48` is `domain.ProviderDefinition.Metadata` registered at package `init()` — `"baseline_commission_percent": 0.16` is a **UI/catalog metadata key**, not a pricing resolver fallback. There is **no call site**: a repo-wide search finds exactly one Go occurrence (that literal) and no reader by name. The actual miss path is degrau-4 `pricingtariffdefaults` (13.00/16.00), per `root.go:844-851` and `pricing/adapters/tariffcomposite` — so *"um miss cai em 0.16"* is false and *"spec pina contra o call site"* is unfalsifiable. Worse, the key is governed by published contracts M-07 does not own: `wiki/framework/provider-metadata-contract.md:41` lists it `number` **required = yes**; `contracts/api/marketplace-central.openapi.yaml:4941` names it a stable metadata key; `packages/sdk-runtime/src/index.ts:122` declares `baseline_commission_percent?: number`. `mission.md:285` limits M-07's OpenAPI/SDK axis to *"/pricing DTOs proveniência (par)"*, so the deletion is unowned on the contract axis and assigned to no milestone. Note `mission.md:180` preserved the honest disjunction — *"morre ou vira row `origem='config'`"* — which IC-01:106 collapsed to deletion without addressing the published contract.
**Yes-if:** M-07 F-01 takes the second ratified arm of `mission.md:180` (re-express the baseline as a `channel_fees` row with `origem='config'`, leaving the metadata key intact), **or** the deletion arm names the contract fallout as owned deliverables — `wiki/framework/provider-metadata-contract.md:41` (required → optional), `contracts/api/marketplace-central.openapi.yaml:4941`, `packages/sdk-runtime/src/index.ts:122` — and the false "fallback silencioso / call site / um miss cai em 0.16" language is deleted (R-25: falsehood is deleted, not hedged).

### N-3 — BLOCKING — check 4 (unreachable acceptance branch)
**Locus:** `M-07-pricing-fee-read/F-01-fee-read-resolver/feature.md:47-48` (EARS-3).
**Offending token:** `While nenhuma fonte resolve (nem defaults), when Resolve roda, the resultado shall ser erro/honest-absent`.
Verified in code: `pricing/adapters/tariffdefaults/resolver.go` `Resolve` is **total** for commission — it always returns `domain.ComponentResolution{Valor: &comissaoPct, Fonte: domain.FontePadrao, Degrau: 4}`, with `error` only on store failure; and its store call `calc_repository.go:239-269` `GetTariffDefaults` **materializes** the row (`INSERT … ON CONFLICT DO NOTHING`, DB column defaults 13.00/16.00) before selecting it. The "nem defaults" state is therefore unreachable through the very chain F-01 reuses (`:29-30` *"resolver existente `pricingtariffdefaults.NewResolver` … REUSADO como base da cadeia"*), and F-01 forbids the only edits that could make it reachable — `:46` *"13.00/16.00 materialize-on-read `calc_repository.go:239-269` INTOCADO"*, `:65-66`, `:89` Forbidden paths *"calc_repository defaults"*. This also refutes `p5-reconciliation-r01.md:78-79`.
**Yes-if:** EARS-3 is re-scoped to the branch that actually exists — degrau-4 **store error** ⇒ typed error, never a constant — or it is deleted and the F-01 truth table at `:94` is reduced from 4 fixtures to the 3 reachable ones (camada2 hit, camada1 hit, config), with `M-07/milestone.md:79` *"sem ledger → `config` com origem visível"* left as the total statement.

### N-4 — BLOCKING — check 6 (verbatim propagation of a ratified ADR value) + check 4
**Locus:** `research/webhook-inbox-interface-contract.md:88-90` against `mission.md:198`; consequence at `M-08-webhook-ingest/milestone.md:77` and `M-08-webhook-ingest/F-02-worker-callback/feature.md:81`.
**Offending tokens:** IC-04 `UNIQUE (provider, notification_id) WHERE notification_id IS NOT NULL` … `Sem notification_id: sem dedupe de transport`, versus ratified ADR-11 at `mission.md:198` `Dedupe (provider, topic, resource, notification_id)`.
Two ratified layers state different dedupe keys and the plan records no amendment. The narrower key is not merely a restatement: IC-04:53 itself declares `notification_id text` NULL (*"campo `_id` do payload quando presente"*), and `research/external-ml-api-facts.md:16` (fact #6, the verified payload inventory) lists only `resource`+`topic`+`user_id`+`attempts`+`sent` — `_id` is not among the verified fields. In the no-`_id` branch the partial index is inert, so ML's 8 retries/1h (fact #7, `external-ml-api-facts.md:17`) produce 8 inbox rows. Against that, `M-08/milestone.md:77` asserts the **total** claim *"Replay da mesma notificação → zero duplicata"* and `M-08/F-02:81` *"Replay → zero duplicata (contagem)"* — a total acceptance criterion over a partial guard whose gap the contract itself names (the R-24 / "guard parcial sob frase total" class).
**Yes-if:** IC-04:88 restores the ADR-11 tuple with a NULL-safe expression (e.g. `UNIQUE (provider, topic, resource, COALESCE(notification_id,''))`) and `M-08/F-01:30-32` quotes it — **or** `mission.md:198` is amended with the narrowing recorded, and the Done-Means at `M-08/milestone.md:77` + `M-08/F-02:81` are narrowed verbatim to *"replay com `_id` presente"*, with the no-`_id` path explicitly documented as covered only by IngestOrder idempotence.

### N-5 — ADVISORY — check 6
**Locus:** `M-01-ml-client-hardening/F-02-items-multiget-raw-dto/feature.md:25` and `:43`.
**Offending token:** `fato #3` / `Fato #3 (multiget 20/call) em research/external-ml-api-facts.md`.
The source is named explicitly, so the citation is checkable and wrong: `external-ml-api-facts.md:13` fact **#3** is *"Cotação de frete por dimensão/peso sem `item_id` NÃO existe"*; multiget-20 is fact **#13** (`:23`, *"Items scan: 1000/batch com `scroll_id`; multiget 20"*). Root cause: two independent 1..N fact namespaces are both cited as "fato #N" — `external-ml-api-facts.md` (1..20, used by M-01 and M-08 F-01:50) and `p5-prerequisites.md` (1..12, used by M-05/M-06/M-07/M-09).
**Yes-if:** both loci read `fato #13`, and every "fato #N" citation names its file the way `M-01/F-02:43` already does.

### N-6 — ADVISORY — check 6 (A-1 residual at uncited sibling loci)
**Locus:** `research/channel-fees-interface-contract.md:30` and `research/divergences-interface-contract.md:29`.
**Offending token:** `range 0086-0088`.
The ratified allocation is `0086-0089` (`mission.md:235`, `mission.md:280`, `M-02/milestone.md:53-54`, `M-02/F-01:25`, and IC-01's own `:133`). The r01 fold corrected only the two loci A-1 named; the same sentence pattern survives one line into each contract's Resources section, so IC-01 now contradicts itself between `:30` and `:133`.
**Yes-if:** both read `0086-0089`.

### N-7 — ADVISORY — check 2 (ownership loci must name paths)
**Locus:** `M-02-sync-core-seam/F-02-fee-divergence-ports/feature.md:76`.
**Offending token:** `apps/server_core/internal/...` fees/divergences.
An ellipsis is not a path; the canonical form is pinned one file away at `M-02/milestone.md:47-50` (`internal/modules/channelfees/`, `internal/modules/divergences/`, with *"layering AGENTS obriga `modules/`"*). This is the exact string shape r01 F-1 was raised against.
**Yes-if:** `:76` quotes `apps/server_core/internal/modules/channelfees/` and `apps/server_core/internal/modules/divergences/`.

### N-8 — ADVISORY — check 2 / check 3 (matrix does not record a same-package split)
**Locus:** `mission.md:280` (M-02 cell) and `mission.md:287` (M-09 cell).
**Offending tokens:** M-02 `scheduler.go (fix pontual)`; M-09 `endpoint /sync/health + FE seção`.
Both milestones run in lane A and both add files to the **same Go package** `apps/server_core/internal/modules/sync/application/`: M-02 F-03 owns `scheduler.go` plus *"helper de contrato de cursor (arquivo novo no package sync)"* (`M-02/F-03:72-73`), while M-09 F-01 owns *"`sync/application/` (read service + porta)"* (`M-09/F-01:87`). File-level disjointness is asserted in prose on both sides (`M-09/milestone.md:57` *"`scheduler.go` INTOCADO (M-02 é o dono do fix)"*; `M-09/F-01:89` Forbidden paths), so this is not a live collision — but the matrix's Go-packages axis shows two owners of one package with no recorded additive lock, unlike the listings case which the mission does record at `mission.md:297-298`.
**Yes-if:** a matrix rule row records the split (M-02 owns `sync/application/scheduler.go` + cursor-contract helper; M-09 additive-only for `sync/transport/**` and `sync/application/health_*`), in the same form as the M-04/M-05 listings additive lock.

### N-9 — ADVISORY — check 6 (wire name vs repo idiom)
**Locus:** `research/divergences-interface-contract.md:134`; echoed at `research/listings-sync-interface-contract.md:46` and `:93`.
**Offending token:** ``/listings` ganha query param `divergentes=true``.
Verified in code: `listings/transport/query.go:31` only interprets keys prefixed `filter.`; `listings/domain/filter.go:9` holds the bare key list. A bare `?divergentes=true` is therefore **silently ignored** — the filter is a no-op with no 400 to reveal it. The correct wire form appears only in a feature brief (`M-05/F-03:30` *"novo `filter.divergentes=true` (entra em `domain.FilterKeys`, `filter.go:9`)"*, and `M-05/milestone.md:80`), while IC-02 is the binding artifact a spec is told not to re-decide.
**Yes-if:** IC-02:134 (and IC-07:46) read `filter.divergentes=true`, or name the `filter.` prefix idiom with the `query.go:31` anchor.

### N-10 — ADVISORY — check 4 (acceptance criterion with no instrument)
**Locus:** `M-07-pricing-fee-read/F-01-fee-read-resolver/feature.md:96`; `M-07-pricing-fee-read/milestone.md:76-77`.
**Offending token:** `Must-fail 0.16: reintroduzir constante → guard nomeia` / `0.16 morto em auth_adapter (guard de allowlist -2 entradas C/D + must-fail)`.
The named instrument is the ADR-05 read-guard, defined at `M-02/F-04:25-26` as *"teste arquitetural que enumera TODO sítio read-time que pode tocar o client ML, como allowlist explícita de exatamente 4 entradas"* — its domain is ML read-time call sites (sites C/D are `root.go:847` and `:848` per `codebase-read-side.md:48`), not numeric literals. No instrument in the plan can fail on a re-introduced `0.16`, so the must-fail cannot be discharged. The two criteria are also conflated into one line at `milestone.md:76-77`.
**Yes-if:** the constant guard is named as its own instrument (e.g. a source-assert test in the pricing package that fails on a commission literal), separated from the allowlist −2 criterion — or the must-fail is dropped and only the allowlist shrink is claimed.

---

## 5. VERDICT

**NEEDS-REVISION.**

All ten r01 blocking findings (F-1 … F-10) are genuinely closed at their loci, and all nine advisories (A-1 … A-9) are applied. The fold is honest — but it introduced/left four blocking defects, three of them in loci the fold did not touch, and one (N-3) that refutes a claim the reconciliation itself made.

Blocking, must be resolved before P6:
- **N-1** — `M-02/F-03:49-51`: job return-type change contradicts IC-06 Must Preserve, its own Inputs, and its own Forbidden paths (`products_job.go:90`).
- **N-2** — `M-07/F-01:37-38` + `M-07/milestone.md:32-33,53,76` + `IC-01:106`: `auth_adapter.go:47` `baseline_commission_percent` is published-contract metadata with no pricing call site; deletion is unowned on the OpenAPI/SDK/wiki axis and rests on a false premise.
- **N-3** — `M-07/F-01:47-48`: EARS-3 unreachable — degrau-4 is total over a materialize-on-read store the feature forbids touching.
- **N-4** — `IC-04:88-90` vs `mission.md:198`: dedupe key silently narrowed from the ratified ADR-11 tuple; total "zero duplicata" Done Means (`M-08/milestone.md:77`, `M-08/F-02:81`) unsatisfiable in the no-`_id` branch the contract itself names.

Advisory (fold recommended, not gating): **N-5** … **N-10**.

Re-audit on an r03 manifest is required after the fold. The Sol P5 retroactive touchpoint remains mandatory before `status: planned` (≥ 2026-08-05) regardless of this verdict.
