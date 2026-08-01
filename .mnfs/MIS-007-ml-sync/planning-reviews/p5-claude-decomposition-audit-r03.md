# P5 Decomposition Audit — round 03 (VERBATIM)

```yaml
type: planning-review-audit
phase: P5
round: 03
created: 2026-07-31
auditor: cold Claude Opus crew (task a37f6975f3368ccfc; operator-ratified waiver — Sol P5
  retroactive mandatory before status: planned)
input_manifest: p5-input-r03.sha256 (digest 94faf38b37e907073434b423fd331faa6828b001ef6a3bcc86c8b7ef0d1e7407)
provenance: task output-file returned 0 bytes (third instance of the transport fragility);
  text below recovered VERBATIM from the task-notification in the session transcript JSONL
  (line 1164) and persisted post-compaction; transport HTML entities (&lt; &gt; &amp;) decoded.
verdict: NEEDS-REVISION
```

---

Verification complete. Deliverable follows.

---

# P5 Decomposition Audit — round 03 (MIS-007-ml-sync)

Cold read-only auditor, operator-ratified waiver substituting the GPT-5.6 Sol P5 touchpoint. Sol retroactive pass remains MANDATORY before `status: planned` regardless of this verdict. **Nothing was edited, written, or created.**

---

## 1. Header

**Manifest self-digest verified: YES.**
`sha256sum planning-reviews/p5-input-r03.sha256` → `94faf38b37e907073434b423fd331faa6828b001ef6a3bcc86c8b7ef0d1e7407` — byte-identical to the expected digest.

**`sha256sum -c planning-reviews/p5-input-r03.sha256` (run from mission root, Git Bash): 46 OK / 0 FAILED.**

**Files read: ALL 46 (no sampling).** Mission-root-relative:

`mission.md` · `research/channel-fees-interface-contract.md` · `research/codebase-ingest-side.md` · `research/codebase-read-side.md` · `research/divergences-interface-contract.md` · `research/external-ml-api-facts.md` · `research/listings-sync-interface-contract.md` · `research/orders-persistence-interface-contract.md` · `research/p5-prerequisites.md` · `research/sync-health-interface-contract.md` · `research/sync-ingest-ports-interface-contract.md` · `research/webhook-inbox-interface-contract.md` · `M-01-ml-client-hardening/milestone.md` · `M-01-ml-client-hardening/F-01-resilience-decorator/feature.md` · `M-01-ml-client-hardening/F-02-items-multiget-raw-dto/feature.md` · `M-02-sync-core-seam/milestone.md` · `M-02-sync-core-seam/F-01-core-ddl/feature.md` · `M-02-sync-core-seam/F-02-fee-divergence-ports/feature.md` · `M-02-sync-core-seam/F-03-scheduler-incremental-cursor/feature.md` · `M-02-sync-core-seam/F-04-read-guard-allowlist/feature.md` · `M-03-orders-shipment-persist/milestone.md` · `M-03-orders-shipment-persist/F-01-ml-ingest-readers/feature.md` · `M-03-orders-shipment-persist/F-02-ingest-order-v1/feature.md` · `M-03-orders-shipment-persist/F-03-read-path-switch/feature.md` · `M-04-listings-backfill-ingest/milestone.md` · `M-04-listings-backfill-ingest/F-01-listings-ddl/feature.md` · `M-04-listings-backfill-ingest/F-02-mass-closure-replacement/feature.md` · `M-04-listings-backfill-ingest/F-03-backfill-cursor-ingest/feature.md` · `M-04-listings-backfill-ingest/F-04-scheduler-refresh-wiring/feature.md` · `M-05-listings-fees-divergence/milestone.md` · `M-05-listings-fees-divergence/F-01-camada2-fee-ingest/feature.md` · `M-05-listings-fees-divergence/F-02-stock-divergence/feature.md` · `M-05-listings-fees-divergence/F-03-anuncios-fe-contract/feature.md` · `M-06-orders-backfill-decomposition/milestone.md` · `M-06-orders-backfill-decomposition/F-01-backfill-incremental/feature.md` · `M-06-orders-backfill-decomposition/F-02-decomposition-camada3/feature.md` · `M-06-orders-backfill-decomposition/F-03-audit-fe-pedidos/feature.md` · `M-07-pricing-fee-read/milestone.md` · `M-07-pricing-fee-read/F-01-fee-read-resolver/feature.md` · `M-07-pricing-fee-read/F-02-precos-provenance-fe/feature.md` · `M-08-webhook-ingest/milestone.md` · `M-08-webhook-ingest/F-01-inbox-endpoint/feature.md` · `M-08-webhook-ingest/F-02-worker-callback/feature.md` · `M-09-sync-observability/milestone.md` · `M-09-sync-observability/F-01-sync-health-endpoint/feature.md` · `M-09-sync-observability/F-02-integracoes-health-section/feature.md`

Prior-round context read but **not** audited: `planning-reviews/p5-claude-decomposition-audit-r02.md`, `planning-reviews/p5-reconciliation-r02.md`.

Repo anchors verified directly this round (not via the plan's claims): `apps/server_core/internal/modules/sync/application/scheduler.go:36-46,155-161`; `internal/modules/sync/composition/products_job.go:19-25,90,98-133`; `internal/modules/integrations/adapters/mercadolivre/auth_adapter.go:34-54,181-192,250-261`; `internal/modules/pricing/adapters/tariffdefaults/resolver.go` (full); `internal/composition/root.go:366-382,576-601,826-858`; `internal/modules/listings/transport/query.go:28-40`; `internal/modules/listings/domain/filter.go:9`; `internal/modules/orders/domain/order_bucket.go:48`; `internal/modules/sync/adapters/postgres/sync_state_repo.go:35`; `apps/server_core/migrations/` (highest = `0085_erp_import_products_sellable_fields.sql`); `apps/web/src/pages/integracoes/IntegracoesPage.tsx:558-574`; `packages/sdk-runtime/src/index.ts:122,2113,2144,2446`; `contracts/api/marketplace-central.openapi.yaml:4941`; `wiki/framework/provider-metadata-contract.md:41`.

---

## 2. r02 findings N-1 … N-10

**N-1 — job return-type change vs IC-06 Must Preserve. CLOSED.**
`M-02/F-03:49-54` Expected Output now reads: *"o scheduler DERIVA o tipo do run do `cursor.phase` obrigatório no JSON que o job retorna (IC-06:136-137 — campo já binding; mesma fonte que IC-05 usa). `JobFunc` fica BYTE-IDÊNTICA (`scheduler.go:46` — Must Preserve IC-06); NENHUM tipo de retorno muda; NENHUM job concreto editado (`products_job.go` intocado — auditoria P5 r02 N-1)."* Verified in code: `scheduler.go:46` is `type JobFunc func(ctx context.Context, cursor json.RawMessage) (json.RawMessage, error)` and `products_job.go:90` returns `syncapp.JobFunc` — both left untouched by the new text. The ambiguity r02 flagged is gone. **The mechanism the arm moved the requirement onto is itself defective — NEW BLOCKING P-1 below; that is a new defect, not a failure to close N-1.**

**N-2 — `0.16` deletion on a false premise + unowned published contracts. CLOSED. Adaptation judged SOUND.**
`mission.md:180-183`: *"EMENDA (auditoria P5 r02 N-2): `baseline_commission_percent: 0.16` (`auth_adapter.go:42-48`) é METADATA de catálogo do provider — contrato publicado (wiki required, OpenAPI stable key, SDK typed), SEM call site em pricing; fica INTOCADA (a disjunção anterior 'morre ou vira row `origem='config'`' caiu — premissa de fallback silencioso era falsa)."* Propagated verbatim at `M-07/milestone.md:27-29`, `M-07/F-01:38-40`, `IC-01:106-109`, `IC-01:177` (Must-Not-Decide), and `M-01/milestone.md:47-50` (*"`internal/modules/integrations/adapters/mercadolivre/` — OAuth/metadata de catálogo, NINGUÉM edita nesta missão — P5 r02 N-2"*). Every false token r02 quoted (`fallback silencioso 0.16`, `um miss cai em 0.16 hardcoded silencioso`, `spec pina contra o call site`, `auth_adapter.go:47-48 (0.16) MORRE no M-07`) is gone from the tree.

**Adversarial judgment on the ADAPTED arm (no `channel_fees` config row created).** The adaptation genuinely satisfies the finding; it does not leave a gap.
- The finding's two defects were (a) a false premise and (b) an unowned deletion on the OpenAPI/SDK/wiki axis. Under the adaptation there is **no deletion at all**, so (b) cannot exist: verified `wiki/framework/provider-metadata-contract.md:41` still `| baseline_commission_percent | number | yes |`, `contracts/api/marketplace-central.openapi.yaml:4941` still names it a stable key, `packages/sdk-runtime/src/index.ts:122` still `baseline_commission_percent?: number;` — none is claimed by any milestone. (a) is deleted, not hedged (R-25 honored).
- The literal arm-1 text ("re-express as a `channel_fees` row `origem='config'`") would have contradicted ratified F-10, which `IC-01:79-85` pins: the reader resolves *"SÓ o ledger — camada 2 → 1 → **ausente TIPADO**"* and *"o degrau `config` é COMPOSIÇÃO do consumidor pricing (M-07 F-01)"*. Writing a config row into the ledger would have reopened the exact two-owner defect r01 F-10 closed. The rejection rationale is correct, and the recorded deviation is in `mission.md` ADR-09, not only in the reconciliation.
- Verified in code that the premise for leaving it intact holds: `auth_adapter.go:42-48` is `domain.ProviderDefinition.Metadata` set inside `func init()`; a repo-wide search returns exactly one Go occurrence (`auth_adapter.go:47`) and zero readers by name. No pricing call site exists.
- Non-gating observation (not a finding): the metadata key carries `0.16` while the `config` degrau carries `16.00` — two units for the same 16%. The plan touches neither and makes no claim about the unit, so no plan defect.

**N-3 — EARS-3 unreachable. CLOSED.**
`M-07/F-01:48-50`: *"While o store de defaults FALHA (única lacuna real do degrau-4 — `Resolve` do tariffdefaults é total: materialize-on-read garante a row; auditoria P5 r02 N-3), when Resolve roda, the resultado shall ser erro TIPADO — NUNCA uma constante embutida."* Verified in code: `tariffdefaults/resolver.go` `Resolve` returns `domain.TariffResolution{}, fmt.Errorf("tariffdefaults: get tariff defaults: %w", err)` **only** on the store call, and otherwise always populates `comissao.Valor = &comissaoPct` — total for commission, error only on store failure. The re-scoped branch is exactly the one that exists. Truth table `M-07/F-01:96-97` matches: *"4 fixtures (camada2 hit, camada1 hit, config, store de defaults falhando → erro TIPADO)"*.

**N-4 — dedupe key silently narrowed from ADR-11. CLOSED. Adaptation judged SOUND, with one residual advisory.**
Amendment recorded at the ratified layer — `mission.md:201-206`: *"Dedupe de transporte (EMENDA P5 r02 N-4, estreitado da tupla original): `UNIQUE (provider, notification_id) WHERE notification_id IS NOT NULL` — `_id` NÃO está no payload verificado (fato #6, `external-ml-api-facts.md`) e a tupla cheia com COALESCE bloquearia p/ SEMPRE re-notificação legítima do mesmo resource; sem `_id` não há dedupe de transporte e a idempotência REAL é do IngestOrder (ADR-04) + reconciliação 5min."*

**IC-04 cross-reference present** (`research/webhook-inbox-interface-contract.md`, Persistence Expectations): *"(Chave estreitada da tupla original de ADR-11 — emenda registrada em `mission.md` ADR-11, P5 r02 N-4: tupla cheia com COALESCE bloquearia re-notificação legítima do mesmo resource p/ sempre.)"*

**M-08 Done-Means narrowed verbatim** (`M-08/milestone.md:77-80`): *"Replay com `_id` presente → zero row nova (upsert attempts_provider). Sem `_id` NÃO há dedupe de transporte (IC-04/ADR-11 emendado — P5 r02 N-4): replay vira rows extras inofensivas; a prova de zero-duplicata é de DOMÍNIO — IngestOrder idempotente, zero efeito duplicado."*
**F-02 validation asserts on the effect, not the count** (`M-08/F-02:81-83`): *"replay SEM `_id` → rows extras processadas com ZERO efeito duplicado de domínio (IngestOrder idempotente — asserção no EFEITO, não na contagem de inbox; IC-04/N-4)."*
**M-08 F-01 consistency verified** (`M-08/F-01:29-31`): *"Dedupe UNIQUE(provider,notification_id) parcial → duplicata faz upsert de attempts_provider, não row nova."* — no total "zero duplicata" claim anywhere in F-01. The "guard parcial sob frase total" class is discharged at every M-08 locus.

Amendment completeness is genuine: rationale recorded at the ADR, cross-referenced at IC-04, and both consumer loci narrowed. One locus the fold did not reach → **ADVISORY P-3** below (`mission.md:325`).

**N-5 — `fato #3` miscitation / two fact namespaces. CLOSED.**
`M-01/F-02:25` now *"batches de 20 — fato #13"*; `:44` *"Fato #13 (multiget 20/call) em `research/external-ml-api-facts.md`"*. Full sweep verified: every `fato #N` in the manifested tree names its source file — `M-01/F-01:48`, `M-01/milestone.md:63`, `M-08/F-01:50` → `external-ml-api-facts.md`; `M-05/F-03:49`, `M-06/F-03:49`, `M-07/F-01:54`, `M-07/milestone.md:69`, `M-09/F-02:49` → `p5-prerequisites.md`; `mission.md:203` → `external-ml-api-facts.md`. Zero unqualified citations remain.

**N-6 — `0086-0088` residual at uncited siblings. CLOSED.**
`research/channel-fees-interface-contract.md:30` *"Tabela `channel_fees` (nova, M-02, range 0086-0089)"*; `research/divergences-interface-contract.md:29` *"Tabela `divergences` (nova, M-02, range 0086-0089)"*. IC-01 self-consistent with its own `:136`. Repo-wide check: no `0086-0088` survives anywhere in the tree.

**N-7 — ellipsis path in an ownership block. CLOSED.**
`M-02/F-02:76-78`: *"Owned paths: packages novos (`apps/server_core/internal/modules/channelfees/` e `apps/server_core/internal/modules/divergences/` — nomes finais no spec, layering AGENTS obriga `modules/`) + seus testes."* Canonical form now matches `M-02/milestone.md:47-50`.

**N-8 — same-package split not recorded in the matrix. CLOSED.**
`mission.md:307-310`: *"**Posse do package `sync/application/`** (lane A, M-02 ∥ M-09 — lock aditivo registrado, P5 r02 N-8): M-02 F-03 é dono de `scheduler.go` + helper de contrato de cursor (arquivo novo); M-09 é ADDITIVE-ONLY — arquivos novos `sync/application/health_*` + `sync/transport/**`; `scheduler.go` intocado pelo M-09."* Same idiom as the listings additive lock at `:305-306`.

**N-9 — bare `divergentes=true` wire name. CLOSED.**
`research/divergences-interface-contract.md:134-137`: *"`/listings` ganha o filtro `filter.divergentes=true` — idioma do … (`listings/domain/filter.go:9`); um `?divergentes=true` sem prefixo seria silenciosamente [ignorado]"*. `IC-07:46` Operations row: *"`filter.divergentes` (IC-02; idioma `filter.` — `query.go:31`)"*; `IC-07:93` error row `filter.divergentes` não-booleano → 400. `M-05/milestone.md:26` residual fixed. Verified in code: `listings/transport/query.go:31` `if strings.HasPrefix(key, "filter.")`; `listings/domain/filter.go:9` `FilterKeys` bare list — the idiom is stated correctly.

**N-10 — must-fail with no instrument. CLOSED.**
The `0.16` must-fail died with the N-2 adaptation (no locus in the tree names it). The allowlist criterion now stands alone with its own instrument: `M-07/F-01:98-99` *"Allowlist (M-02 F-04) -2 entradas C/D: teste do guard atualizado; must-fail = reintroduzir chamada ML read-time em pricing → allowlist reprova nomeando o sítio"*; identical at `M-07/milestone.md:80-81`. Instrument domain (ML read-time call sites, `M-02/F-04:25-26`) now matches the must-fail's subject.

**All ten r02 findings CLOSED at their loci.**

---

## 3. NEW findings (fresh full pass, 8 checks)

### P-1 — BLOCKING — check 4 (contract satisfiability) + check 5 (prerequisite existence) + check 6 (ratified ADR value widened)

**Loci:** `M-02-sync-core-seam/F-03-scheduler-incremental-cursor/feature.md:26-28`, `:49-51`, `:71`; `research/sync-ingest-ports-interface-contract.md` Compatibility Rules; against `mission.md:168-170` (ADR-07) and `M-02-sync-core-seam/milestone.md:76`.

**Offending tokens.** `M-02/F-03:26-28`: *"`incremental` deixa de ser hardcoded false (`scheduler.go:160`): passa a refletir o tipo do run DERIVADO do `cursor.phase` **obrigatório** do cursor que o job retorna"*. `M-02/F-03:71` Negative Scenarios: *"JSON de cursor sem `phase` → contrato falha nomeando o campo."* IC-06 Compatibility Rules: *"Campos extras no cursor são livres POR entidade (opaco), MAS `phase` é obrigatório e estável."*

**Contradiction, three ways.**
1. **Against the ratified ADR.** `mission.md:168-170` ADR-07 scopes the requirement narrowly: *"`scheduler.go:42-45` APAGA cursor nil ⇒ **jobs novos** retornam cursor não-nil com phase explícita"*. IC-06 and F-03 widened a *jobs-novos* rule into a total one over every cursor. That is P4/P5 widening a ratified spine value — check 6.
2. **Against the code.** Verified: `sync/composition/products_job.go:22-25` is `type ProductsCursor struct { Source string \`json:"source"\`; Processed int \`json:"processed"\`; CompletedAt time.Time \`json:"completed_at"\` }` — **no `phase` field**, and `:123-127` marshals exactly that struct as the returned cursor. The only live job in the repo violates the total statement today.
3. **Against F-03's own ownership block.** `M-02/F-03:52-53` Expected Output: *"NENHUM job concreto editado (`products_job.go` intocado — auditoria P5 r02 N-1)"*; `:77` Forbidden paths: *"jobs concretos"*. So the brief mandates a mandatory field on a cursor it forbids itself from changing, on the one job that exists.

**Consequence the plan does not resolve.** With `phase` absent from products' cursor the scheduler cannot derive the run type for products, so `M-02/milestone.md:76` Done-Means *"`incremental` reflete tipo do run; job products inalterado (regressão verde)"* is unfalsifiable on the only real job — F-03's own validation (`:82`) proves the fix with a *"job fake"* only. Lane A also carries the opposite rule for the same field: `M-09/F-01:71-72` Constraints pins *"parse tolerante: cursor de formato desconhecido → phase null (nunca erro)"*, and `M-09/milestone.md:33-35` *"campos que dependem do fix incremental (M-02 F-03) ficam honestos até lá"*. Two lane-A artifacts, same field, opposite totality — mandatory-and-failing vs tolerant-and-null.

**Yes-if:** either (a) IC-06's Compatibility Rule and `M-02/F-03:26-28,71` are narrowed to ADR-07's ratified wording — `phase` obrigatório nos cursores dos **jobs novos** (M-04/M-06), scheduler com parse tolerante (`phase` ausente ⇒ `incremental=false`, comportamento de hoje) matching `M-09/F-01:71-72` verbatim, and the contract test's domain is stated as the new jobs; **or** (b) `apps/server_core/internal/modules/sync/composition/products_job.go` is added to M-02 F-03's owned paths as a named updated call site (one additive field on `ProductsCursor`), with the products regression asserting the pre-existing `source`/`processed`/`completed_at` fields unchanged, and `M-02/F-03:52-53` + `:77` amended to stop forbidding it.

### P-2 — BLOCKING — check 4 (contract satisfiability) + check 1 (DAG/lane structure)

**Loci:** `mission.md:222` and `mission.md:253` (ADR-14, ratified) against `mission.md:252`, `mission.md:300-302`, `M-05/milestone.md:61-62`, `M-06/milestone.md:61-62`, `M-07/milestone.md:61-63`; risk row `mission.md:354`.

**Offending tokens.** `mission.md:222`: *"SDK é ESCRITO À MÃO (client = objeto literal `index.ts:2113-2330`) ⇒ **no máximo 1 milestone com mudança de contrato FE em voo**."* `mission.md:253`: *"Edge transversal ADR-14: ≤1 milestone com contrato FE em voo **(sobrepõe disjunção de código)**."* Versus `mission.md:252`: *"**C** M-05∥M-06∥M-07"* and `mission.md:300-302`: *"≤1 milestone FE-contract em voo — na lane C, M-05/M-06/M-07 têm TODOS par FE ⇒ hub serializa os **merges** de contrato dentro da lane (**código pode paralelizar**; o COMMIT de contrato não)."*

**Why this is unsatisfiable as written.** `:253` states the ADR-14 edge **overrides code disjointness** — i.e. it is a constraint on being in flight, strictly stronger than file-level disjointness. `:300-302` then resolves the identical situation the opposite way, permitting three FE-contract milestones in flight simultaneously and downgrading the rule to merge ordering. The matrix confirms all three carry the axis: `mission.md:291` M-05 `/listings DTOs+param (par)`, `:292` M-06 `/orders DTOs margem (par)`, `:293` M-07 `/pricing DTOs proveniência (par)`. The mission's own risk register makes the collision explicit: `mission.md:354` R-7 (likelihood **H**) names as its trigger *"2 milestones em voo listando OpenAPI"* — lane C as planned trips its own trigger by construction, permanently, not as an exception. The ICs restate only the strong form (`IC-03:153` *"ADR-14: ≤1 milestone FE-contract em voo"*; `IC-07:144` *"ADR-14 ≤1 em voo"*; `IC-05:147` *"M-09 conta como milestone FE-contract em voo (ADR-14: ≤1)"*), so the relaxation exists at exactly one locus and contradicts every other statement of the rule. A milestone orchestrator cannot dispatch lane C without knowing which of the two governs, and the answer changes the lane structure.

**Yes-if:** either (a) `mission.md:222` and `:253` are amended to the rule actually being planned — *"≤1 COMMIT de contrato por vez; código paralelo permitido"* — with the hand-written client literal named as a hub-arbitrated seam in the same form as M-07's pricing region of `root.go` (`mission.md:293`, *"edita região pricing existente (única exceção — hub arbitra)"*), and R-7's trigger at `:354` re-worded so it does not fire on the ratified plan; **or** (b) lane C is re-drawn with the three FE-contract milestones serialized, and `:252` / the three `Runs in parallel with` lines are corrected to match.

### P-3 — ADVISORY — check 6 (N-4 amendment not propagated to the mission-level criterion)

**Locus:** `mission.md:325` (Quality Attributes, row Q3).
**Offending token:** *"reprocessar notificação → zero duplicata"*.
The N-4 fold narrowed this claim at `M-08/milestone.md:77-80` and `M-08/F-02:81-83` to *replay com `_id`* for row count and to a domain-effect assertion for the no-`_id` branch, but the mission-level validation criterion is still stated unqualified. P6 authors validation contracts from this table; unqualified, it reads as a row-count criterion that the ratified partial key cannot satisfy in the no-`_id` branch — the same total-phrase-over-partial-guard shape N-4 was raised against.
**Yes-if:** `:325` reads *"reprocessar notificação → zero duplicata de DOMÍNIO (IngestOrder idempotente; dedupe de inbox só com `_id` — ADR-11 emendado)"*, or cross-references `M-08/milestone.md:77-80`.

### P-4 — ADVISORY — check 6 (stale anchor inside a ratified ADR)

**Locus:** `mission.md:222` (ADR-14).
**Offending token:** `` `index.ts:2113-2330` ``.
Verified in code: the client object literal opens at `packages/sdk-runtime/src/index.ts:2113` (`return {`) and closes at `:2446` (`};`); line 2330 falls mid-literal, inside `createProfitabilityManualAdjustment`. `research/p5-prerequisites.md` already carries the correct `index.ts:2113-2446`, so the ADR contradicts the prerequisites file it was built from. Same class as r01 A-2/A-8.
**Yes-if:** `:222` reads `index.ts:2113-2446`.

### P-5 — ADVISORY — check 2 (matrix cell asserts an axis the features contradict)

**Locus:** `mission.md:291` (M-05 row, `root.go wiring` column) against `M-05/F-01:53-54` and `M-05/F-02:47-50`.
**Offending token:** the M-05 `root.go wiring` cell is `—`, while `M-05/F-01:53-54` Expected Output promises *"Passo de fee no ingest … + **fiação do ChannelFeeWriter no pipeline**"* and `M-05/F-02:47-50` promises *"Avaliador + **fiação pós-sweep** + leitura de mirror por porta"*. Both require constructing M-02's postgres impls (`channelfees`, `divergences`) and injecting them into the listings ingest composition. If that construction lands in `root.go` the cell is false; if it lands inside `listings/composition` under M-04's additive lock, no brief says so. Every other milestone that wires something carries an explicit `1 linha ancorada` cell (`:289`, `:290`, `:292`, `:294`, `:295`). Not a live collision — M-04 is closed before lane C opens — but the axis is unallocated where R-7 (likelihood H) is the named risk.
**Yes-if:** the M-05 cell reads `1 linha ancorada (writers de fee/divergência)`, **or** `M-05/F-01` and `M-05/F-02` pin that the wiring lives inside the listings composition package under the recorded additive lock and never touches `root.go`.

**Checks passing with no finding:** check 3 (feature write-DAG — every cross-milestone overlap carries a named serial edge or a recorded additive lock at `mission.md:303-311,315-317`); check 5 for the remaining assumed symbols (`DeriveOrderBucket` `order_bucket.go:48`, `sync_state_repo.go:35` `Read`, `root.go:370-378` `AccessTokenResolver`/`credentialResolver` — and the M-08 `user_id`→installation map is grounded: `auth_adapter.go:192,261` persist ML `user_id` as `ProviderAccountID`; `root.go:591-592` shipment/buyer-fiscal readers; migrations highest `0085` with M-02/M-04/M-08/M-06 ranges `0086-0089`/`0090-0092`/`0093`/`0094-0095` disjoint and all above it); check 7 (all 25 briefs carry Negative Scenarios, Ownership and Validation Expectations; all 4 FE briefs carry `## Interaction Model` — `M-05/F-03`, `M-06/F-03`, `M-07/F-02`, `M-09/F-02`; the 6 briefs without `## Inputs/Outputs` are internal/test-only surfaces with no asymmetric I/O shape); check 8 (no implementation planning beyond seam/wiring pins; no product scope beyond the ratified design).

---

## 4. VERDICT

**NEEDS-REVISION.**

All ten r02 findings (N-1 … N-10) are genuinely closed at their loci, and both recorded arm-adaptations hold up under adversarial reading: **N-2's** adaptation removes the deletion entirely, so the published wiki/OpenAPI/SDK contracts stay intact and unowned-axis exposure is zero, while F-10's owner-split is preserved — no gap; **N-4's** amendment is recorded at the ADR with rationale, cross-referenced at IC-04, and narrowed verbatim at both M-08 consumer loci with the no-`_id` proof correctly relocated to the domain effect — complete but for one un-updated mission-level criterion (advisory).

Blocking, must be resolved before P6:
- **P-1** — `M-02/F-03:26-28,49-53,71,77` + IC-06 Compatibility Rules vs ratified ADR-07 `mission.md:168-170` and the code: `phase` is made mandatory over all cursors while `products_job.go:22-25` `ProductsCursor` has no such field and the brief forbids editing it; `M-09/F-01:71-72` pins the opposite (tolerant, phase null) for the same field in the same lane.
- **P-2** — `mission.md:222,253` vs `:252,300-302` and the three lane-C milestone files: ADR-14 ratifies ≤1 FE-contract milestone in flight *overriding code disjointness*, and the Parallel Execution Plan schedules three, tripping the mission's own R-7 trigger (`:354`) by construction.

Advisory (fold recommended, not gating): **P-3**, **P-4**, **P-5**.

A re-audit on an r04 manifest is required after the fold. The Sol P5 retroactive touchpoint remains MANDATORY before `status: planned` (≥ 2026-08-05) regardless of this verdict.
