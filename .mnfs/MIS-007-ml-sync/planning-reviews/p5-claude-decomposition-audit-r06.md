# P5 Claude Decomposition Audit — round 06 (VERBATIM)

```yaml
type: planning-review-verbatim
phase: P5
round: 06
task_id: aa9d98bf0fcdd8c7d
auditor: cold Claude Opus crew (general-purpose, read-only mandate; operator-ratified
  waiver -- Sol P5 retroactive mandatory before status: planned)
input_manifest: p5-input-r06.sha256 (digest 5fcb02fbc59f470a6663bc47c6d12c2704a228ab061c6da979845f1697535bd6)
persisted: 2026-07-31, same turn the task notification arrived
provenance: task output-file transport; recovered VERBATIM from session transcript JSONL
  (task-notification <result> block, html-unescaped) -- SIXTH round using this path
verdict: NEEDS-REVISION (PART A: all F-r05-1..6 CONFIRMED CLOSED, zero reopened;
  2 blocking F-r06-1/F-r06-2 + 4 advisory F-r06-3..6)
```

---

# P5 Decomposition Audit — MIS-007-ml-sync — Round 06

**Auditor:** cold independent Claude Opus crew (read-only; zero files written)
**Date:** 2026-07-31
**Input manifest:** `planning-reviews/p5-input-r06.sha256` — read first; self-digest recomputed as `5fcb02fbc59f470a6663bc47c6d12c2704a228ab061c6da979845f1697535bd6` (**matches** the value supplied in the dispatch). Full `sha256sum -c` over all 46 entries: **46 OK, 0 mismatched**. The audit below is provably against the frozen post-fold content.
**Scope:** exactly the 46 manifested files. `planning-reviews/*` read only as PART A context (`p5-reconciliation-r05.md`).
**Repo tip for prerequisite verification:** `C:\Users\leandro.theodoro\Documents\marketplace-central`, branch `main` @ `dd89d4b3`.

---

# PART A — Closure verification of F-r05-1 .. F-r05-6

## F-r05-1 (blocking) — M-09/F-01 scan predicate could never return products — **CONFIRMED CLOSED**

Locus: `M-09-sync-observability/F-01-sync-health-endpoint/feature.md:36-41` (Brief) and `:87` (Negative Scenarios).

Current text, verbatim:

> `SEM param `installation_id` — tenant vem do ctx (IC-05 pina); a leitura varre TODAS as rows de `sync_state` do tenant, independente de `installation_id` — INCLUI o sentinela de escopo ERP `installation_id = "erp"` (`sync/composition/scheduler.go:11`), onde vive a row de products; scan restrito a instalações ML devolveria ZERO products e tornaria o Done Means inatingível (auditoria P5 r05 F-r05-1).`

Negative scenario now reads:

> `- Tenant sem NENHUMA row em `sync_state` → entities todas null-honestas + webhook`

Repo check: `apps/server_core/internal/modules/sync/composition/scheduler.go:11` — `const InstallationScopeERP = "erp"`; consumed at `products_job.go:50` — `scheduler := syncapp.NewScheduler(repo, InstallationScopeERP, interval, time.Now)`. **Exact.** The fix names the real sentinel and the emptiness case is now defined as absence of ROW, not absence of ML installation. Consistent with IC-05:94-95 (entity-agnostic) and M-09 Done Means. No new contradiction introduced.

## F-r05-2 (blocking) — audit formula dropped `financing_add_on_fee` — **CONFIRMED CLOSED**

Locus 1: `M-06-orders-backfill-decomposition/F-03-audit-fe-pedidos/feature.md:27-30` (Brief):

> `esperado_unit = detail.percentage_fee × unit_price/100 + detail.fixed_fee + detail.financing_add_on_fee; esperado_total = esperado_unit × quantity — TODA componente de fee do detail canônico IC-01 entra (dropar financing_add_on_fee abriria divergência falsa PERMANENTE em anúncio com parcelamento — auditoria P5 r05 F-r05-2)`

Locus 2: same file `:63-67` (Constraints, the "NUNCA do amount" pin) carries the identical three-term formula. **Both loci fixed** — no half-fix.

Locus 3: `research/channel-fees-interface-contract.md:117-125` — canonical camada-2 example now `"value":15.99` with the 5-key detail and the derivation note:

> `(Exemplo aritmeticamente consistente com a fórmula ratificada da auditoria 3→2 — `value = percentage_fee × price_used/100 + fixed_fee + financing_add_on_fee` = 12.5% × 79.90 + 6.00 + 0 = 15.99; auditoria P5 r05 F-r05-2.)`

Arithmetic re-verified by hand: 12.5% × 79.90 = 9.9875; + 6.00 + 0 = 15.9875 → 15.99 at `numeric(14,4)`/2-dp presentation. Camada-3 example unchanged and still consistent: 8.11 × 3 = 24.33. **Closed.**

## F-r05-3 (advisory) — listings additive-lock grant narrower than M-05's real write-set — **CONFIRMED CLOSED**

Locus 1: `mission.md:312-316` (transversal rules):

> `- **Posse de ingest de listings**: M-05 estende superfícies do M-04 PÓS-close (lock aditivo registrado: `listings/application/**`, `listings/transport/**` e `listings/adapters/postgres/repository.go`, additive-only — o write-set REAL do M-05 F-03 inclui transport e o repository de leitura, não só application; grant alinhado ao write-DAG na auditoria P5 r05 F-r05-3).`

Locus 2: `mission.md:296` (M-05 Go-packages matrix cell) carries the same enumeration. Cross-checked against `M-05/F-03` Owned paths — the grant is now a superset. **Closed.**

## F-r05-4 (advisory) — IC-05 "sem ela QA reprova" vs M-09 SOFT dependency — **CONFIRMED CLOSED**

Locus: `research/sync-health-interface-contract.md:108-111`. The failure condition is now scoped to the mission/live criterion, with the pre-fix uniform-NULL state explicitly declared honest and gate-passing. `M-02 F-03` remains declared SOFT in M-09's dependency block (unchanged, already correct). No contradiction remains between IC-05 and the M-09 gate.

## F-r05-5 (advisory) — M-05/F-01 EARS enumerated 3 of 5 detail keys — **CLOSED AT THE TWO NAMED LOCI; RESIDUAL OF THE SAME CLASS REMAINS**

Fixed loci, verbatim:

- `M-05-listings-fees-divergence/F-01-camada2-fee-ingest/feature.md:37-39` — `the sistema shall upsert row camada 2 com value amount + detail canônico IC-01 COMPLETO (5 chaves: percentage_fee, fixed_fee, financing_add_on_fee, price_used, listing_type_id — F-r05-5) + coletado_em.`
- `M-05-listings-fees-divergence/milestone.md:76-78` — `row camada 2 com amount + detail canônico IC-01 COMPLETO (5 chaves: percentage_fee, fixed_fee, financing_add_on_fee, price_used, listing_type_id — F-r05-5)`

Residuals still enumerating only 3 keys — see **F-r06-4** below. The r05 fold's own "residual sweep" (reconciliation `:62-65`) missed two further instances of the identical class. Closure of the finding as scoped is confirmed; the class is not fully swept.

## F-r05-6 (advisory) — mission Handoff stale at P4 — **CONFIRMED CLOSED**

Locus: `mission.md:371-385`. Now reads `planning_phase = decompose (P5)` with the audit loop and rounds r02–r05 named; waiver row preserved verbatim (`:379-381` — `Sol P3/P5/P7 RETROATIVOS a partir de 2026-08-05, todos obrigatórios ANTES de `status: planned``); `planning-reviews/p5-*` present in Required artifact paths (`:382-383`). **Closed.**

## Reopened-closure sweep (r01–r04)

All seven named spot-checks located and observed intact after the r05 fold:

| prior closure | locus | state |
| --- | --- | --- |
| ADR-14 ≤1-COMMIT amendment (r03 P-2) | `mission.md:219-227`, restated `:257-258`, `:307-309` | intact, three loci agree |
| GREATEST `last_success_at` (r04 F-r04-1) | `research/sync-health-interface-contract.md:94-100`; M-09/F-01 restatement | intact; repo re-verified: `sync_state_repo.go:74-79` COALESCE-per-column confirms `last_full_sync_at` freezes on `incremental=true`, so GREATEST is the only honest definition |
| WebhookStatsReader reference semantics (r04 F-r04-2) | IC-05 §seam; `M-08/F-02:56-61` (`injeção por referência/ponteiro, nunca builder-cópia`) | intact; repo re-verified: builder-by-value idiom is real at `root.go:856` |
| `external_account_id` webhook mapping (r04 F-r04-4) | IC-04; `M-08/F-02:44-49` (port-of-package, no cross-module SQL) | intact; repo re-verified: `migrations/0016_integrations_foundation.sql:33` defines the column; `installation_repo.go:81` already selects it |
| root.go import-block ownership (r04 F-r04-5) | `mission.md:303-306`; `M-07/F-01:90-91` | intact; repo re-verified: `root.go:99,101` are the tarifflive/tariffcomposite imports |
| write-DAG guard-allowlist enumeration (r04 F-r04-6) | `mission.md:325-330` | intact; A/B removal (M-03 F-03) and C/D removal (M-07 F-01) both named and serialized by the lane B→C edge |
| phase-tolerant-parse scoping (r03 P-1) | `research/sync-ingest-ports-interface-contract.md:136-139`; `M-02/F-03` | intact; repo re-verified: `products_job.go:22-26` `ProductsCursor{Source,Processed,CompletedAt}` — **no `phase` field**, exactly as IC-06:138 asserts |

**No r01–r04 closure reopened.**

---

# PART B — Fresh full audit (8 checks over all 46 manifested files)

## Check 1 — DAG edge completeness and justification

Lanes `A M-01∥M-02∥M-09 → B M-03∥M-04 → C M-05∥M-06∥M-07 → D M-08` are consistent with the 13-row edge-justification table (`mission.md:272-286`) and with every milestone's own `## Dependencies` block, with one asymmetry: **F-r06-3** below. No false edges found: each hard edge names a concrete forcing artifact (ports, migrations, cursor contract, single-writer inheritance), and each soft edge (`⤳`) is labelled data-quality rather than compile.

### F-r06-3 — advisory — DAG edge completeness

**Cited excerpt** (`M-06-orders-backfill-decomposition/milestone.md:50-51`, verbatim):

> `- M-05 SOFT (auditoria 3→2 precisa de camada 2 populada — sem ela auditoria fica muda, não quebra; lane C permite qualquer ordem de close).`

corroborated at `M-06-orders-backfill-decomposition/F-03-audit-fe-pedidos/feature.md:51` — `M-05 F-01 (camada 2 — dependência SOFT: sem ela auditor muda)`.

**Defect locus:** `mission.md:265-270` (Dependency DAG block), `mission.md:272-286` (edge-justification table), and `mission.md:251` (Milestone Strategy row for M-06, `depende de: M-03, M-02`).
**Offending token:** the absent `M-06⤳M-05` edge.
**Defect:** the mission spine draws and justifies every other soft edge — `M-07⤳M-05` (`:268`, `:283`) and `M-09⤳M-04/M-06` (`:269`, `:286`) — but omits `M-06⤳M-05`, which two downstream artifacts declare. The hub adjudicates lane-C close order from the mission DAG; a soft edge invisible at that level can be scheduled as if it did not exist.
**Yes-if:** add `M-06⤳M-05` to the DAG block and one row to the edge table with the forcing artifact already written in `M-06/milestone.md:50-51` (`auditoria 3→2 precisa de camada 2 populada`) and the same non-blocking qualifier used for `M-07⤳M-05`. Advisory only — the milestone-level declaration means no execution correctness depends on it, and a fixture may plant camada-2 rows directly since M-02 owns the DDL.

## Check 2 — Canonical six-axis disjointness, seam locks, migration ranges

**Clean.** Each of the six axes was cross-read for all 9 milestones (`mission.md:290-300`) against each milestone's `## Ownership & Concurrency` block.

- **Go packages/files:** `orders/**` held serially (M-03 lane B → M-06 lane C, `:310-311`); `listings/**` held by M-04 with M-05's additive lock enumerated (`:312-316`); `sync/application/` split by file with M-09 additive-only on new `health_*` files and `scheduler.go` reserved to M-02 F-03 (`:317-320`); `capability_adapter.go` exclusive to M-01 with new files in the dir permitted downstream (`:292`).
- **Migrations:** M-02 `0086-0089` · M-04 `0090-0092` · M-08 `0093` · M-06 `0094-0095` (reserve, may stay empty). **Disjoint, contiguous, no overlap.** Each range is restated identically in three loci (mission matrix, milestone `Migration block`, feature brief). Repo tip verified: highest existing migration is `0085_erp_import_products_sellable_fields.sql`; **no file 0086 or higher exists** — the allocation starts clean.
- **DB write-sets:** `channel_fees` partitioned by layer (M-05 layer 2 / M-06 layer 3), `divergences` partitioned by kind (M-05 `estoque` / M-06 `tarifa`), both partitions inside the same natural key `(…, layer, fee_kind)` / `(…, kind)` so upserts cannot collide. M-07 declared read-only over `channel_fees`.
- **OpenAPI/SDK:** the seam is correctly identified as the hand-written client literal `index.ts:2113-2446` (repo-verified: `return {` at 2113, closing `};` at 2446) with ≤1 contract COMMIT in flight, hub-serialized inside lane C.
- **FE routes:** `/anuncios` (M-05), `/pedidos` (M-06), `/precos` (M-07), `/integracoes` (M-09) — disjoint.
- **root.go:** one anchored constructor line per milestone; M-07's region-edit `:828-858` declared the single exception, hub-arbitrated; the import block declared ownerless and hub-resolved.

## Check 3 — Feature-level write-DAG

**Clean.** `mission.md:323-330` enumerates every cross-milestone write-set collision with its serializing edge or recorded grant: M-05.F-01/F-02 after M-04 close under the additive lock; M-06.* after M-03 close by lane inheritance; M-07's root.go pricing region hub-arbitrated; the guard-allowlist file (owner M-02 F-04) written by M-03 F-03 and M-07 F-01, serialized by the lane B→C edge; the root.go import block hub-resolved. Feature-level `Owned paths` / `Forbidden paths` blocks were read for all 25 briefs and no residual overlap without an edge or grant was found. The F-r05-3 widening removed the one previously-real gap (M-05 F-03's transport + repository writes).

## Check 4 — Contract satisfiability

**Two blocking defects.** Arithmetic and predicate consistency across IC examples was otherwise verified by hand and is sound: IC-01 camada 2 (15.99), IC-01 camada 3 (8.11 × 3 = 24.33), IC-03 decomposition (239.70 − 48.66 − 22.90 − 95.10 = 73.04, `margem_pct` 30.47), tolerances (`estoque` 0, `tarifa` R$0.01) restated identically in IC-01, IC-02, M-05 and M-06.

### F-r06-1 — **blocking** — contract satisfiability

**Cited excerpt A** (`M-06-orders-backfill-decomposition/F-02-decomposition-camada3/feature.md:25-31`, verbatim):

> `grava (a) rows camada 3 do fee ledger — uma por linha de pedido, subject_type=`order_line`, subject_id=`<provider_order_id>:<provider_item_id>` (formato IC-01), value = TOTAL da linha (sale_fee_unit × quantity — sale_fee é POR UNIDADE, fato T2), detail obrigatório `{"sale_fee_unit":..,"quantity":..}`, origem=`api_order`; frete seller do shipment (M-03) → row fee_kind=`freight`, subject_type=`order`, subject_id=`<provider_order_id>` (frete é do pedido, não da linha — IC-01), origem=`api_shipment`;`

**Cited excerpt B** (`M-02-sync-core-seam/F-02-fee-divergence-ports/feature.md:31`, verbatim):

> `camada 3 NUNCA na resolução; recusa de camada 3 sem `detail.sale_fee_unit`/`quantity`.`

reinforced at the same file `:70` — `- UpsertFee camada 3 sem detail obrigatório → erro nomeado (IC-01).` — and at `research/channel-fees-interface-contract.md:135`:

> `Rejeição canônica: upsert de camada 3 com `detail` sem `sale_fee_unit`/`quantity` → writer recusa (constraint de aplicação; teste nomeia).`

**Defect locus:** the guard clause in `M-02/F-02:31` and `:70`, and IC-01 `:135`, versus the producer in `M-06/F-02:29-31`.
**Offending token:** the unqualified `camada 3` in the rejection rule (no `fee_kind` scope).
**Defect:** the freight row that M-06 F-02 promises to write is grammatically and semantically **inside** the "(a) rows camada 3" clause — it is a layer-3 row (`origem=api_shipment`, realized from the actual shipment) with `fee_kind=freight`, `subject_type=order`, and **no `detail`** (a shipment freight cost has no `sale_fee_unit`/`quantity` decomposition). IC-01 `§Persistence Expectations:102-104` correctly scopes the detail mandate — `Camada 3 **comissão**: … `detail` OBRIGATÓRIO com `{"sale_fee_unit": x, "quantity": n}`` — but IC-01 `§Canonical Examples:135` and both M-02 F-02 loci drop the `comissão` scope. The M-02 writer, built to the brief as written, **refuses the M-06 freight row**. M-06's own Done Means (`milestone.md:80-82`, `frete_seller` inside the canonical decomposition) and IC-03's canonical example (`− 22.90` freight term) cannot be satisfied. This is a plan-time defect: two ratified briefs cannot both be implemented.
**Yes-if:** scope the guard to commission at all three loci — `recusa de camada 3 **fee_kind=commission** sem detail.sale_fee_unit/quantity` in `M-02/F-02:31` and `:70`, and the same qualifier in IC-01 `:135` (aligning the Canonical Examples wording to the already-correct `§Persistence Expectations:102`) — **or**, if a layer-3 freight row is not intended at all, move the freight row to the layer IC-01 §Enums assigns to freight resolution and re-state its layer explicitly in `M-06/F-02:29-31`. Either arm is inside approved scope; no new product scope is required.

### F-r06-2 — **blocking** — contract satisfiability

**Cited excerpt A** (`research/channel-fees-interface-contract.md:71-73`, `## Required Outputs`, verbatim):

> ``ResolveListingFees` retorna por fee_kind: `{value, value_type, currency, layer, origem, coletado_em}` — proveniência SEMPRE junto do número; consumidor que exibe número sem proveniência reprova milestone (ADR-09).`

restated without `detail` at `M-02-sync-core-seam/F-02-fee-divergence-ports/feature.md:37-38`:

> `- While row camada 2 existe p/ o listing, when ResolveListingFees, the reader shall retornar {value, layer:2, origem, coletado_em} — nunca o config.`

**Cited excerpt B** (`M-07-pricing-fee-read/F-01-fee-read-resolver/feature.md:43-44`, verbatim):

> `- While anúncio vinculado tem camada 2, when Resolve roda, the comissão shall vir do ledger (detail.percentage_fee/fixed_fee da camada 2) com origem api_listing_prices + coletado_em.`

**Defect locus:** IC-01 `§Required Outputs:71-73` (and its restatement at `M-02/F-02:37-38`) versus `M-07/F-01:43-44`.
**Offending token:** `detail.percentage_fee/fixed_fee` in an EARS clause whose only data source is `ChannelFeeReader.ResolveListingFees`.
**Defect:** M-07 F-01 is defined (`:27-28` — `resolve comissão via ChannelFeeReader (IC-01) em cascata camada 2 … → camada 1 … → fallback`) as consuming the port exclusively, and its `Forbidden paths` (`:93`) bar it from `channel_fees` writes and from `listings`, so it has no second read path. The port's contracted return tuple contains no `detail`, so the M-07 EARS clause is unsatisfiable against the ratified IC-01 output shape. This is not symmetrical with M-05 F-03, which composes `detail` in the listings **read service** directly from `channel_fees` (`M-05/F-03:26`) rather than through the port — that path is fine and is precisely what makes the omission invisible at first read.
**Yes-if:** add `detail` to the resolution tuple in IC-01 `§Required Outputs:71-73` (`{value, value_type, currency, layer, detail, origem, coletado_em}`) and mirror it in `M-02/F-02:37-38` — **or** re-word `M-07/F-01:43-44` to promise only what the tuple carries (`value` + `layer:2` + `origem` + `coletado_em`), dropping the `detail.percentage_fee/fixed_fee` parenthetical. The first arm is the smaller edit and preserves the M-07 F-02 provenance surface; either is inside approved scope. Note this arm interacts with F-r06-1: whichever arm is chosen, the `detail` semantics must be stated once and scoped by layer, not restated per-consumer.

## Check 5 — Prerequisite existence

**Substantially clean.** Every assumed-existing symbol, wiring, table, and endpoint named in the briefs was either verified against the repo or explicitly assigned to a named creating upstream feature. Verified this round (repo-root, file:line):

- `sync/composition/scheduler.go:11` `InstallationScopeERP = "erp"`; consumed `products_job.go:50` — **exact**.
- `sync/application/scheduler.go:46` `JobFunc`; `:85` `RegisterJob`; `:42-45` nil-cursor-erases doc; `:160` `incremental` hardcoded `false` — **all exact**.
- `sync/adapters/postgres/sync_state_repo.go:62` `RecordSuccess`; `:74-79` COALESCE-per-column upsert; `:95-102` `RecordFailure` — **exact**.
- `migrations/0075_sync_sync_state.sql:23-28` columns exact; `entity` has **no DB CHECK** (comment `:12-14` declares it a semantic enum validated in Go) — the Go enum `sync/domain/sync_state.go:18-22` does contain `listings` and `orders`, so IC-06:104's operative conclusion (**no migration needed to register jobs**) holds.
- Highest migration `0085`; nothing ≥ 0086 — **exact**.
- `listings/transport/query.go:31` `filter.` prefix; `listings/domain/filter.go:9` `FilterKeys` — **exact**.
- `listings/adapters/postgres/repository.go:390-394` the unconditional `UPDATE listings SET status='closed'` MASS-CLOSURE statement — **exact**.
- `root.go:259-272` `registerBatchRoutes`; `:672-677` products scheduler bootstrap; `:99,101` tarifflive/tariffcomposite imports; `:837` `pricingtariffdefaults.NewResolver`; `:852` batchOrch; `:853-855` calcSvc — **all exact**.
- `orders/transport/http_handler.go:549` `DeriveOrderBucket` call site; `:85` `/orders/import`; `:89` `/orders/{provider_order_id}` — **exact**.
- `profitability/application/service.go:1014-1015` `MissingSaleFee` flag — **exact**.
- `integrations/adapters/mercadolivre/auth_adapter.go:47` `"baseline_commission_percent": 0.16` inside the provider `Metadata` map — **exact**; the N-2 "no call site in pricing" claim holds.
- `migrations/0016_integrations_foundation.sql:33` `external_account_id` — **exact**.
- `packages/sdk-runtime/src/index.ts:2113-2446` client object literal — **exact**.
- `capability_adapter.go`: 429 handling at `:654-655`, `:462-463`, `:578-580` (single attempt, no retry); `:444` `PUT /items/{id}`; `:79-92` `ProviderCapabilitySet`; `:681-688` `providerDiag` 512-rune clip — **exact**.
- `pricing/adapters/tarifflive/resolver.go:35` `NewResolver`; `:43` `ResolveCommission` — **exact**.
- `apps/web/src/pages/integracoes/IntegracoesPage.tsx:558-574` mount order — **exact**.
- `integrations/domain/refresh_policy.go` exists at the path the research note cites.
- **Assigned-to-creating-feature (correctly absent from the repo):** `GET /sync/health` (M-09 F-01), `notifications_inbox` (M-08 F-01, 0093), `channel_fees` / `divergences` / `order_shipments` (M-02 F-01, 0086-0089), `listing_variations` (M-04 F-01, 0090-0092). Zero occurrences under `apps/server_core` — consistent with the plan; no false "already exists" assumption anywhere.

Residual: four anchor line-ranges have drifted or are mischaracterized. Raised as one consolidated advisory.

### F-r06-5 — advisory — prerequisite anchor drift

**Cited excerpts and their measured contradictions:**

1. `M-07-pricing-fee-read/F-01-fee-read-resolver/feature.md:35` — `fiação root.go:844-850 vira `tariffResolver := feeledger.NewResolver(...)``. Measured: the `if feeReader, ferr := marketplaceCapabilities.FeeQuoteReader(...)` block spans **845-851**; line 844 is the `var tariffResolver` declaration that must **survive**, and 851 is the closing brace the brief's range omits. Off-by-one in both directions on a range the brief instructs a chip to delete.
2. `research/sync-health-interface-contract.md` — cites `root.go:854-857` for the builder-by-value idiom. Measured: `pricingHandler = pricingHandler.WithCalc(calcSvc)` is at **856**; 854-855 belong to the `calcSvc` builder chain and 857 is `}`.
3. `M-02-sync-core-seam/F-02-fee-divergence-ports/feature.md:48` and `M-04-listings-backfill-ingest/F-02-mass-closure-replacement/feature.md:28` — both cite `internal_read/adapters/mirror/writer.go:74-105` as the binding upsert idiom. Measured: `upsertSQL` is **74-95**; the keep-absent constant (`keepAbsentSQL`, ADR-04's actual "never DELETE" idiom) is a **separate constant at 104-112**. The cited range stops at 105, cutting the keep-absent statement in half and omitting the part both briefs actually depend on.
4. `research/codebase-ingest-side.md:50` — `Há testes de migração no próprio diretório (`migrations/listings_test.go:25` …) que fazem assert por regex no SQL`, echoed as a pattern instruction in `mission.md:212` (`teste regex por migração (estilo `listings_test.go:25`)`). Measured: line 25 is `"CREATE TABLE IF NOT EXISTS listings",` — a plain string literal in a `required` slice. The `regexp.MustCompile` logic is at **:101** inside `createTableBody`. A chip copying "the style at `:25`" copies substring matching, not regex.

**Defect locus:** the four line-ranges above.
**Offending token:** the numeric ranges `844-850`, `854-857`, `74-105`, and `:25`.
**Yes-if:** correct each range to the measured value (`845-851`; `856`; `74-95` for the upsert plus `104-112` for keep-absent, cited as two ranges; `:101` for the regex idiom with `:25` re-labelled as the required-substring list). Advisory — none of these change a decision, but the r05 reconciliation recorded "28 anchor classes verified EXACT", and these four were not among them. Lesson `chip-anchors-3-f5-unidade` applies: a cited range that does not contain the statement it names is prose with the form of a fact.

## Check 6 — Propagation of approved values

Numbers, enums, key names, and route names were traced from each IC into every restating brief. Correct throughout: the `phase` vocabulary `{backfill, incremental, sweep}` (IC-06:87 → M-02 F-03, M-04 F-03, M-06 F-01, M-09 F-01); the `notifications_inbox` status enum and the `attempts ≥ 5 → dropped` rule (IC-04:59-60, `:69-70` → `M-08/F-02:36-38`, exact match); the `subject_id` composite formats with `:` separator (IC-01:51-56 → M-06 F-02, IC-02); the `origem` vocabulary with `api_shipping_options` correctly held as an additive reserve with no producer this mission (IC-01:65-67); the tolerances; the IC-03 decomposition key list (10 keys, restated complete at `M-06/F-02:33-34`); the migration ranges; the `versao:1` marker; `filter.divergentes=true`.

### F-r06-4 — advisory — propagation (F-r05-5 residual class)

**Cited excerpt A** (`M-05-listings-fees-divergence/F-01-camada2-fee-ingest/feature.md:90-91`, `## Validation Expectations`, verbatim):

> `- Fixture: ingest de N anúncios → N rows camada 2 verificadas por SELECT (amount + detail exatos do fixture: percentage_fee, fixed_fee, price_used).`

**Cited excerpt B** (`M-05-listings-fees-divergence/F-03-anuncios-fe-contract/feature.md:26-27`, verbatim):

> ``tarifa` (objeto: `amount` string decimal, `detail` {percentage_fee, fixed_fee, price_used} quando presente, `origem`, `coletado_em` — COMPOSTO de channel_fees camada 2 no read service, IC-07 …; shape = IC-01 canonical camada 2)`

**Defect locus:** `M-05/F-01:91` and `M-05/F-03:26`.
**Offending token:** the 3-key enumeration `percentage_fee, fixed_fee, price_used`.
**Defect:** IC-01's canonical camada-2 `detail` has **five** keys (`:119-120`), and the r05 fold established the full tuple at `M-05/F-01:37-39` and `M-05/milestone.md:76-78`. These two loci were missed. `M-05/F-03:26` is the more consequential of the two: it is the **DTO shape** for `/listings`, so a 3-key object shipped there silently drops `financing_add_on_fee` from the FE contract — the exact component whose omission F-r05-2 identified as generating a permanent false `tarifa` divergence. `M-05/F-01:91` weakens the fixture assertion to a subset, so a writer that omits `financing_add_on_fee` passes green (lesson `chip-import-chain`: an observable that passes in both worlds is not evidence).
**Yes-if:** replace both enumerations with the 5-key tuple already written at `M-05/F-01:37-39`, or replace them with a citation to IC-01 §Canonical Examples camada 2 as binding rather than re-enumerating. Advisory — IC-01 remains binding by reference in both briefs (`M-05/F-01:71`, `M-05/F-03:26` `shape = IC-01 canonical camada 2`), so the contradiction is resolvable at execution time without a decision.

## Check 7 — Required brief sections

Mechanical sweep of all 25 feature briefs: `## Negative Scenarios` **25/25**, `## Validation Expectations` **25/25**, `## Ownership` **25/25**. All four UI features carry `## Interaction Model` (`M-05/F-03`, `M-06/F-03`, `M-07/F-02`, `M-09/F-02`). Validation expectations name inspectable proof throughout — SELECT-verified fixtures, named must-fail tests, golden DTO baselines, >1-page fixtures (R-3 / lesson CHIP-MERCADO), and hub live-drive criteria. Six briefs lack `## Inputs/Outputs`; four of those are non-shape features (`M-01/F-01` decorator, `M-02/F-04` guard test, `M-04/F-04` wiring, `M-09/F-02` UI with Interaction Model instead) where the section would be empty.

### F-r06-6 — advisory — required sections

**Cited excerpt:** `M-08-webhook-ingest/F-02-worker-callback/feature.md` and `M-04-listings-backfill-ingest/F-02-mass-closure-replacement/feature.md` — neither file contains a `## Inputs/Outputs` heading.

**Defect locus:** section list of the two briefs.
**Offending token:** the missing `## Inputs/Outputs` heading.
**Defect:** both are data-shape features. `M-08/F-02` drives the `notifications_inbox` status machine and, in the same brief (`:55-61`), delivers the **real implementation of `WebhookStatsReader`** — an IC-05 port whose output shape (`last_notification_at`/`pending`/`dropped_24h`) reaches `GET /sync/health`. `M-04/F-02` changes the write semantics of `ApplyCompletedPull` including a possible signature change (`:48-50` — `assinatura pode ganhar run metadata (run_started_at, completo:bool)`). Both are drift-risk surfaces the section exists to pin.
**Yes-if:** add `## Inputs/Outputs` to each citing the binding source already named in the brief body — IC-04 `§status flow` + IC-05 `§InboxHealth` for `M-08/F-02`; IC-06 run-complete rule + IC-07 E3 columns for `M-04/F-02`. Advisory — both briefs defer to a binding IC by name in their `## Inputs` block, so no shape is left undecided.

## Check 8 — No implementation planning, no new product scope

**Clean.** All 25 briefs stop at behavioral contract. Where a brief touches implementation, it is pinning a **seam constraint** rather than a plan: "não criar segundo writer" (`M-04/F-02:48`), "`ports.TariffResolver` INTOCADO" (`M-07/F-01:66`), "additive-only sobre M-04" (`M-05/F-01:62-64`), "spec pina valor" for the worker polling interval (`M-08/F-02:66-67`) — the last of which explicitly defers the number to execution, which is the correct direction. Every brief carries `## Execution Artifact Rules` deferring spec/plan/validation to the implementer. Scope is bounded by the ratified design: no brief introduces a screen, entity, provider, or capability outside the mission `## Scope`; `api_shipping_options` is held as an explicitly producer-less reserve (IC-01:65-67); camada 1 is documented as dormant until MIS-008 (`M-07/F-01:72-73`); LISTEN/NOTIFY, fee history, and margin sorting are each named and deferred.

---

# Verdict

**NEEDS-REVISION**

Two blocking findings (**F-r06-1**, **F-r06-2**), both check-4 contract satisfiability, both located at the `channel_fees` `detail` seam. Their yes-if conditions are compatible with each other and with all approved scope — F-r06-1 scopes the layer-3 detail mandate by `fee_kind`, F-r06-2 adds `detail` to (or removes it from) the resolution output tuple; a single coordinated edit to IC-01 plus its two restating briefs discharges both. No mutually incompatible conditions exist, so this is not BLOCKED. PART A: all six r05 findings CONFIRMED CLOSED (F-r05-5 closed as scoped, with a same-class residual raised as advisory F-r06-4); no r01–r04 closure reopened. Four advisory findings (F-r06-3 .. F-r06-6) do not affect the verdict.

---

# Files read

**Manifest and reconciliation (PART A context only):**
- `planning-reviews/p5-input-r06.sha256`
- `planning-reviews/p5-reconciliation-r05.md`

**All 46 manifested files:**
- `mission.md`
- `M-01-ml-client-hardening/milestone.md`, `M-01-ml-client-hardening/F-01-resilience-decorator/feature.md`, `M-01-ml-client-hardening/F-02-items-multiget-raw-dto/feature.md`
- `M-02-sync-core-seam/milestone.md`, `M-02-sync-core-seam/F-01-core-ddl/feature.md`, `M-02-sync-core-seam/F-02-fee-divergence-ports/feature.md`, `M-02-sync-core-seam/F-03-scheduler-incremental-cursor/feature.md`, `M-02-sync-core-seam/F-04-read-guard-allowlist/feature.md`
- `M-03-orders-shipment-persist/milestone.md`, `M-03-orders-shipment-persist/F-01-ml-ingest-readers/feature.md`, `M-03-orders-shipment-persist/F-02-ingest-order-v1/feature.md`, `M-03-orders-shipment-persist/F-03-read-path-switch/feature.md`
- `M-04-listings-backfill-ingest/milestone.md`, `M-04-listings-backfill-ingest/F-01-listings-ddl/feature.md`, `M-04-listings-backfill-ingest/F-02-mass-closure-replacement/feature.md`, `M-04-listings-backfill-ingest/F-03-backfill-cursor-ingest/feature.md`, `M-04-listings-backfill-ingest/F-04-scheduler-refresh-wiring/feature.md`
- `M-05-listings-fees-divergence/milestone.md`, `M-05-listings-fees-divergence/F-01-camada2-fee-ingest/feature.md`, `M-05-listings-fees-divergence/F-02-stock-divergence/feature.md`, `M-05-listings-fees-divergence/F-03-anuncios-fe-contract/feature.md`
- `M-06-orders-backfill-decomposition/milestone.md`, `M-06-orders-backfill-decomposition/F-01-backfill-incremental/feature.md`, `M-06-orders-backfill-decomposition/F-02-decomposition-camada3/feature.md`, `M-06-orders-backfill-decomposition/F-03-audit-fe-pedidos/feature.md`
- `M-07-pricing-fee-read/milestone.md`, `M-07-pricing-fee-read/F-01-fee-read-resolver/feature.md`, `M-07-pricing-fee-read/F-02-precos-provenance-fe/feature.md`
- `M-08-webhook-ingest/milestone.md`, `M-08-webhook-ingest/F-01-inbox-endpoint/feature.md`, `M-08-webhook-ingest/F-02-worker-callback/feature.md`
- `M-09-sync-observability/milestone.md`, `M-09-sync-observability/F-01-sync-health-endpoint/feature.md`, `M-09-sync-observability/F-02-integracoes-health-section/feature.md`
- `research/channel-fees-interface-contract.md`, `research/divergences-interface-contract.md`, `research/orders-persistence-interface-contract.md`, `research/webhook-inbox-interface-contract.md`, `research/sync-health-interface-contract.md`, `research/sync-ingest-ports-interface-contract.md`, `research/listings-sync-interface-contract.md`, `research/p5-prerequisites.md`, `research/external-ml-api-facts.md`, `research/codebase-read-side.md`, `research/codebase-ingest-side.md`

**Repo files inspected for check 5 (paths relative to repo root):**
- `apps/server_core/internal/composition/root.go`
- `apps/server_core/internal/modules/sync/composition/scheduler.go`, `.../sync/composition/products_job.go`
- `apps/server_core/internal/modules/sync/application/scheduler.go`, `.../sync/adapters/postgres/sync_state_repo.go`, `.../sync/domain/sync_state.go`
- `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go`
- `apps/server_core/internal/modules/listings/adapters/postgres/repository.go`, `.../listings/adapters/connectors/source.go`, `.../listings/transport/query.go`, `.../listings/domain/filter.go`
- `apps/server_core/internal/modules/orders/transport/http_handler.go`
- `apps/server_core/internal/modules/pricing/adapters/tarifflive/resolver.go`, `.../pricing/adapters/postgres/calc_repository.go`
- `apps/server_core/internal/modules/profitability/application/service.go`
- `apps/server_core/internal/modules/integrations/adapters/mercadolivre/auth_adapter.go`, `.../integrations/domain/refresh_policy.go`, `.../integrations/adapters/postgres/installation_repo.go`
- `apps/server_core/internal/modules/internal_read/adapters/mirror/writer.go`
- `apps/server_core/migrations/0016_integrations_foundation.sql`, `.../0075_sync_sync_state.sql`, `.../listings_test.go`, and the full `apps/server_core/migrations/` listing
- `apps/web/src/pages/integracoes/IntegracoesPage.tsx`
- `packages/sdk-runtime/src/index.ts`
