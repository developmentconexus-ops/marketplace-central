# P5 Decomposition Audit — round 04 (VERBATIM)

```yaml
type: planning-review-audit
phase: P5
round: 04
created: 2026-07-31
auditor: cold Claude Opus crew (task a4e749acef0c99f7b; operator-ratified waiver — Sol P5
  retroactive mandatory before status: planned)
input_manifest: p5-input-r04.sha256 (digest b425888df8a8f2705fe198eafe35a266b21850c96fadaaa96f0198b859adf55a)
provenance: task output-file returned 0 bytes (fourth instance of the transport fragility);
  text below persisted VERBATIM from the task-notification (recovered from the session
  transcript JSONL) in the same turn it arrived; transport HTML entities decoded.
verdict: NEEDS-REVISION (P-1..P-5 all CONFIRMED CLOSED; 1 new blocking F-r04-1,
  5 advisory F-r04-2..6)
```

---

# P5 Decomposition Audit — MIS-007-ml-sync — Round 04 (cold, independent)

## 1. Header

### Manifest verification

| Step | Result |
| --- | --- |
| sha256 of `planning-reviews/p5-input-r04.sha256` | `b425888df8a8f2705fe198eafe35a266b21850c96fadaaa96f0198b859adf55a` — **MATCHES** expected self-digest |
| `sha256sum -c planning-reviews/p5-input-r04.sha256` (Git Bash, mission root) | **46 OK, 0 FAILED** |

Manifest integrity confirmed. Audit proceeded.

### Files read (all 46 manifested, no sampling)

`mission.md` · `research/channel-fees-interface-contract.md` (IC-01) · `research/divergences-interface-contract.md` (IC-02) · `research/orders-persistence-interface-contract.md` (IC-03) · `research/webhook-inbox-interface-contract.md` (IC-04) · `research/sync-health-interface-contract.md` (IC-05) · `research/sync-ingest-ports-interface-contract.md` (IC-06) · `research/listings-sync-interface-contract.md` (IC-07) · `research/p5-prerequisites.md` · `research/codebase-ingest-side.md` · `research/codebase-read-side.md` · `research/external-ml-api-facts.md` · `M-01-ml-client-hardening/{milestone.md, F-01-resilience-decorator/feature.md, F-02-items-multiget-raw-dto/feature.md}` · `M-02-sync-core-seam/{milestone.md, F-01-core-ddl, F-02-fee-divergence-ports, F-03-scheduler-incremental-cursor, F-04-read-guard-allowlist}/feature.md` · `M-03-orders-shipment-persist/{milestone.md, F-01-ml-ingest-readers, F-02-ingest-order-v1, F-03-read-path-switch}/feature.md` · `M-04-listings-backfill-ingest/{milestone.md, F-01-listings-ddl, F-02-mass-closure-replacement, F-03-backfill-cursor-ingest, F-04-scheduler-refresh-wiring}/feature.md` · `M-05-listings-fees-divergence/{milestone.md, F-01-camada2-fee-ingest, F-02-stock-divergence, F-03-anuncios-fe-contract}/feature.md` · `M-06-orders-backfill-decomposition/{milestone.md, F-01-backfill-incremental, F-02-decomposition-camada3, F-03-audit-fe-pedidos}/feature.md` · `M-07-pricing-fee-read/{milestone.md, F-01-fee-read-resolver, F-02-precos-provenance-fe}/feature.md` · `M-08-webhook-ingest/{milestone.md, F-01-inbox-endpoint, F-02-worker-callback}/feature.md` · `M-09-sync-observability/{milestone.md, F-01-sync-health-endpoint, F-02-integracoes-health-section}/feature.md`.

Prior-round artifacts read as **context only, not audited**: `planning-reviews/p5-claude-decomposition-audit-r02.md`, `p5-reconciliation-r02.md`, `p5-claude-decomposition-audit-r03.md`, `p5-reconciliation-r03.md` (plus `p5-claude-decomposition-audit-r01.md` and `p5-passes-r01.md` consulted solely to confirm no finding below reopens a prior closure).

### Repo anchors verified directly this round (check 5 discipline — plan claims not trusted)

`apps/server_core/internal/modules/sync/composition/products_job.go:22-26` · `internal/modules/sync/application/scheduler.go:46,160` · `internal/modules/sync/adapters/postgres/sync_state_repo.go:59-81` (`RecordSuccess`, incl. `timestampArg(!incremental,…)/timestampArg(incremental,…)` at `:79`) · `apps/server_core/migrations/0075_sync_sync_state.sql` (full) · `migrations/` highest = `0085_erp_import_products_sellable_fields.sql` · `migrations/listings_test.go:20-30` · `internal/platform/httpx/route_deadline.go:20-30` · `internal/composition/root.go:99,101,259-272,366-382,826-860` · `internal/modules/listings/adapters/postgres/repository.go:386-398` (MASS-CLOSURE) · `internal/modules/internal_read/adapters/mirror/writer.go:70-108` (upsert/keep-absent idiom) · `internal/modules/orders/domain/order_bucket.go:19,48` + `order_bucket_test.go:8` · `internal/modules/orders/transport/http_handler.go:248,496,549` · `internal/modules/orders/adapters/postgres/order_repo.go:337-342,378,467,602` · `internal/modules/orders/application/import_service.go:46` · `internal/modules/orders/application/enrich_service.go:188-196` · `internal/modules/orders/ports/buyer_fiscal_reader.go:19-21` · `internal/modules/pricing/ports/tariff_resolver.go:30-32` · `internal/modules/pricing/adapters/{tarifflive,tariffcomposite}/` (only consumers = `root.go:99,101,849,850`) · `internal/modules/integrations/adapters/postgres/{credential_repo.go:21-147, installation_repo.go:26-150}` · `internal/modules/integrations/application/auth_flow_service.go:691` · `internal/modules/integrations/adapters/mercadolivre/auth_adapter.go:181-192,250-261` · `migrations/0016_integrations_foundation.sql:157` · `packages/sdk-runtime/src/index.ts:2113,2130,2144,2288,2290,2359,2361,2446` · tree-wide grep for the pre-amendment ADR-14 strong form.

---

## 2. Verification of r03 findings (P-1 … P-5)

### P-1 (BLOCKING, r03) — `phase` mandate vs. the real products job — **CLOSED, genuinely**

- `research/sync-ingest-ports-interface-contract.md` §Compatibility Rules now narrows the mandate to jobs created by this mission and pins tolerant parse (`phase` absent ⇒ treated as non-incremental, never an error).
- `M-02-sync-core-seam/F-03-scheduler-incremental-cursor/feature.md` carries the narrowing consistently across Brief items 1-2, Expected Output, Negative Scenarios and Validation Expectations; `M-02-sync-core-seam/milestone.md` Done Means matches.
- Against ratified scope: `mission.md` ADR-07 covers only the new jobs — the narrowed form is inside the ratified scope, not a widening.
- Against code: `products_job.go:22-26` has `Source/Processed/CompletedAt` and **no** `phase` field; `scheduler.go:46` `JobFunc` is unchanged; `scheduler.go:160` still hardcodes `incremental=false`. The products regression is now falsifiable against the real job (absent `phase` ⇒ `incremental=false`, which is exactly what the existing cursor produces).
- r02 N-1's closure clauses survive: no brief asks for a `JobFunc` signature change, and `products_job.go` remains listed as untouched. **No earlier closure reopened.**

### P-2 (BLOCKING, r03) — ADR-14 strong form / client literal / R-7 trigger — **CLOSED, verified exhaustively**

- `mission.md:257-258`, `:305-307`, `:359` all read the amended form ("≤1 **COMMIT** de contrato FE em voo; código paraleliza"). ADR-14 body (~:220-228) matches.
- IC restatements amended in step: `orders-persistence-interface-contract.md:152-154`, `listings-sync-interface-contract.md:143-145`, `sync-health-interface-contract.md:146-148`.
- **Tree-wide grep over the manifested set returns zero surviving instances of the strong form** ("≤1 milestone … em voo"). No locus escaped.
- R-7 trigger now reads "2 COMMITS de contrato simultâneos ou não arbitrados pelo hub" (`mission.md:359`) — it no longer fires on the ratified lane C (M-05∥M-06∥M-07), which the matrix rule at `:305-307` explicitly permits.

### P-3 (ADVISORY, r03) — Q3 "zero duplicata" scope — **CLOSED**

`mission.md:330` now reads "zero duplicata de **DOMÍNIO** (IngestOrder idempotente; dedupe de inbox só com `_id` — ADR-11 emendado; ver M-08 Done Means)". Consistent with `M-08-webhook-ingest/milestone.md:77-80` and `F-02` Validation Expectations (assertion on the *effect*, not on inbox row count).

### P-4 (ADVISORY, r03) — SDK client-literal anchor — **CLOSED, verified against the file**

`packages/sdk-runtime/src/index.ts`: `return {` at `:2113`, closing `};` at `:2446`. Every method anchor cited in briefs is accurate (`listListings` `:2130`, `listSyncRuns` `:2144`, `listOrders` `:2288`, `getOrder` `:2290`, `runPricingSimulation` `:2359`, `runBatchSimulation` `:2361`).

### P-5 (ADVISORY, r03) — M-05 wiring vs. `root.go` matrix cell — **CLOSED**

`M-05/F-01` Expected Output and `M-05/F-02` Expected Output both state the wiring lives inside the listings composition package under M-04's recorded additive lock and that "`root.go` NÃO é tocado pelo M-05 (matriz: célula root.go = `—`; auditoria P5 r03 P-5)". `mission.md:296` matrix row M-05, `root.go wiring` cell = `—`. Consistent.

**PART A result: all five r03 findings are closed at their loci; the chosen arms' rationale holds against ratified scope and against code; no r01/r02 closure was reopened by the fixes.**

---

## 3. NEW findings (fresh audit, checks 1-8 over the entire manifested tree)

### F-r04-1 — BLOCKING — checks 4, 6 — `last_success_at := last_full_sync_at` makes the health screen lie for every incremental-only entity

**Locus:** `.mnfs/MIS-007-ml-sync/M-09-sync-observability/F-01-sync-health-endpoint/feature.md:57-60` (Inputs); consumed by `M-09-sync-observability/F-02-integracoes-health-section/feature.md:29-32`; underspecified at `research/sync-health-interface-contract.md:64-94` (Required Outputs).

**Offending token, verbatim** (`M-09/F-01/feature.md`, Inputs):

> `last_error jsonb, consecutive_failures int NOT NULL DEFAULT 0 — todo campo de entities[] tem coluna real; `phase` deriva do cursor jsonb; `last_success_at` = last_full_sync_at).`

**Why it is a defect.** `last_success_at` is *not* a column; it is a derived field, and IC-05 never defines it — `research/sync-health-interface-contract.md:83-94` documents `phase`, the NULL semantics and the `webhook` block, but says nothing about which success `last_success_at` denotes. The brief closes that gap with `last_success_at = last_full_sync_at`. Measured against the real writer, that mapping is wrong the moment M-02 F-03 lands: `sync_state_repo.go:62,79` passes `timestampArg(!incremental, at)` to `last_full_sync_at` and `timestampArg(incremental, at)` to `last_incremental_at`, and the upsert at `:74-75` uses `COALESCE(EXCLUDED.…, sync_state.…)`, so an incremental run **leaves `last_full_sync_at` frozen**. Today the bug at `scheduler.go:160` (`incremental=false` hardcoded) masks this — every run writes `last_full_sync_at`. After the ADR-08 fix that M-02 F-03 owns and that IC-05:94 itself declares a precondition, `orders` (backfill once, then incremental every 5 min per M-06) renders "última sincronização: há 7 dias" while syncing healthily every 5 minutes, and `listings` behaves the same after its first full sweep. `M-09/F-02` feeds exactly this field into the relative timestamp and the freshness badge ("verde=sucesso recente … cinza=nunca"), so the screen shows a stale/never reading for a healthy entity. That is the failure `M-09-sync-observability/milestone.md:66-67` names as disqualifying ("Tela de saúde que mente é pior que nenhuma") and it defeats Q4, whose validation criterion is a browser-driven milestone gate (`mission.md:331`). It is a plan-time defect, not an implementation escalation: the mapping is written into the brief and the contract leaves no other authority to appeal to.

Not a reopened closure: no audit round (r01/r02/r03) adjudicated `last_success_at`. The only prior mention is the strategist's own unaudited pass log, `planning-reviews/p5-passes-r01.md:52-53` ("`last_success_at` := last_full_sync_at. SATISFIABLE"), which is precisely the claim refuted here by the writer's column-selection semantics.

**Yes-if.** IC-05 §Required Outputs defines `last_success_at` as the most recent success of any kind — `GREATEST(last_full_sync_at, last_incremental_at)`, NULL only when both are NULL — and `M-09/F-01` Inputs drops the `last_success_at = last_full_sync_at` equation in favour of that definition; **or** IC-05 renames the field to `last_full_sync_at`, exposes both columns verbatim, and pins that the FE freshness badge and relative timestamp read the later of the two, with a named negative fixture (entity whose only recent success is incremental ⇒ must render fresh/green, never "há N dias"/cinza).

---

### F-r04-2 — ADVISORY — checks 3, 5 — IC-05's `WithWebhookStatsReader` seam leaves an arm that silently drops the injection

**Locus:** `research/sync-health-interface-contract.md:44-48`; restated at `M-09-sync-observability/F-01-sync-health-endpoint/feature.md:33-36`, `M-09-sync-observability/milestone.md:51-55`, `M-08-webhook-ingest/F-02-worker-callback/feature.md:49-54`, `M-08-webhook-ingest/milestone.md:54-57`.

**Offending token, verbatim** (`sync-health-interface-contract.md:44`):

> `Troca: health service expõe `WithWebhookStatsReader(...)` (setter/builder).`

**Why it is a defect.** The contract admits both arms, but only one of them satisfies the other two ratified clauses — "M-08 chama o setter DA PRÓPRIA região ancorada" and "Código e linha de construção do M-09 NUNCA editados" — because M-09 F-01 owns construction *and* mount in a single anchored line (`M-09/F-01` Expected Output: "impl default … + setter … + par OpenAPI+SDK + linha root.go"), and M-08's region necessarily comes later in `root.go`. Under the builder arm (value receiver returning a copy, the idiom already in the file at `root.go:854-857`, `pricingHandler = pricingHandler.WithCalc(calcSvc)`), M-08's later call mutates a copy that nothing serves; the endpoint keeps returning the default `{null,0,0}` forever. The named mitigation is blind to it: the compile-time assert (IC-05:48) proves the port is *implemented*, not that the *registered handler* reads the injected impl, and M-09's own proof obligation ("teste da porta … fake injetado via setter", `M-09/milestone.md:78`) passes under both arms because it exercises the service in isolation, never the route. The only discriminating observable in the plan is `M-08/milestone.md:83` ("attempts ≥5 → dropped, visível na saúde"), at the very end of lane D — the catalog-503 class the contract itself cites, moved to the latest possible point.

**Yes-if.** IC-05 §Seam pins the mutation semantics as reference/pointer — the injected reader must be observable through the handler already registered by M-09 — and adds the proof obligation that the discriminating assertion is made on the **route** (`GET /sync/health` returns inbox-derived stats after injection), not on the port; **or**, if the builder arm is kept, IC-05 pins that M-09's anchored region constructs but does not mount, and names the mission-level wiring point where the mount happens after M-08's region.

---

### F-r04-3 — ADVISORY — checks 5, 6 — IC-03 states a false locus for `DeriveOrderBucket` and mandates a move that cannot happen

**Locus:** `research/orders-persistence-interface-contract.md:61-63`.

**Offending token, verbatim:**

> `Valores = enum HOJE produzido por `DeriveOrderBucket` (transport de market) — função MOVE para o núcleo INALTERADA no M-03; derivação no INGEST (nunca no read). P5 verifica o enum por prerequisite-existence.`

**Why it is a defect.** Two of the three assertions are false against the repo. The function lives in the núcleo already — `apps/server_core/internal/modules/orders/domain/order_bucket.go:48`, `func DeriveOrderBucket(providerStatus, shipmentStatus string, tags []string, faturado bool) OrderBucket` — so "MOVE para o núcleo" names work that does not exist. And it is not in "transport de market": a tree-wide grep finds it only in the `orders` module (call sites `orders/transport/http_handler.go:549`, `orders/adapters/postgres/order_repo.go:378`; truth table `orders/domain/order_bucket_test.go:8`); the `market` module has no reference at all. The IC deferred the check to P5 ("P5 verifica o enum por prerequisite-existence"); P5 performed it — `research/p5-prerequisites.md` §1 records the domain location and signature verbatim — and the refuted prose was not corrected. The substantive decision behind the sentence (move the *derivation call site* from read to ingest) is correct and is propagated correctly to the briefs (`M-03/milestone.md`, `M-03/F-02/feature.md:31-33`: "MESMA função (`order_bucket.go:48`)", truth table "INTOCADA"), which is why this is advisory rather than blocking — but IC-03 is declared binding by `M-03-orders-shipment-persist/milestone.md:17`, a "move" of `order_bucket.go` falls outside every M-03 feature's Owned paths (F-02 owns `orders/application|ports|adapters`; F-03 owns `orders/application|adapters|transport`; `orders/domain/` is in neither), and the mission's own doctrine is that false prose in a contract is deleted, not annotated.

**Yes-if.** IC-03 §Required Inputs — extensão `orders` reads that `bucket` carries the enum produced by the existing `domain.DeriveOrderBucket` (`orders/domain/order_bucket.go:48`, signature and truth table unchanged), and that M-03 moves only the **call site** from the read path to the ingest path — deleting "(transport de market)" and "função MOVE para o núcleo".

---

### F-r04-4 — ADVISORY — check 5 — M-08's `user_id`→installation map is sourced from a resolver that runs the other way

**Locus:** `M-08-webhook-ingest/F-02-worker-callback/feature.md:44-46` (Inputs); Brief restatement at `:26-27`.

**Offending token, verbatim:**

> `credential store existente (AccessTokenResolver wiring root.go:370-378 — mapa user_id→installation a partir das credenciais JÁ armazenadas; sem store novo).`

**Why it is a defect.** `root.go:370-382` is a closure that takes a fully-populated `connectorsdomain.ProviderAccountRef` (tenant + installation + provider + account) and returns an access token. It consumes `ProviderAccountID`; it does not index by it, and it cannot be inverted. The credential store offers no lookup by account either — `credential_repo.go` exposes `NextCredentialVersion`, `GetActiveCredential(installationID)`, `SaveCredentialVersion`, `DeactivateCredential`, `DeactivateAllForInstallation`, all keyed by installation. The capability M-08 needs does exist, but somewhere else: the ML `user_id` is persisted as the installation's `external_account_id` (`auth_adapter.go:192,261` → `auth_flow_service.go:691` `firstNonEmpty(payload.ProviderAccountID, inst.ExternalAccountID)` → `installation_repo.go:118-145` `ApplyConnectionSnapshot`), readable via `InstallationRepository.ListInstallations` (`installation_repo.go:81`, which already selects `external_account_id`). Naming the wrong source in a feature whose Owned paths are limited to the new webhook package invites the implementer to invent the lookup — and the cheapest invention is direct SQL against `integration_installations` from the webhook package, the cross-module read the mission forbids elsewhere (`M-05/F-02` Expected Output: "leitura de mirror por porta (NUNCA SQL cross-module direto)"). r03 recorded this prerequisite as grounded on the strength of `auth_adapter.go:192,261`; that grounds the *datum* being persisted, not the existence of a query in the required direction.

**Yes-if.** `M-08/F-02` Inputs names `integration_installations.external_account_id` (populated at `auth_flow_service.go:691`) as the mapping source, read through the existing installation service / `ListInstallations` (`installation_repo.go:81`) behind a port owned by the webhook package — no direct cross-module SQL, no new store — and drops the `AccessTokenResolver` citation, which stays only as the token source for the authenticated refetch.

---

### F-r04-5 — ADVISORY — check 2 — the `root.go` axis allocates constructor lines but never the import block, and M-07's declared range excludes edits it must make

**Locus:** `mission.md:298` (matrix row M-07, `root.go wiring` cell), `mission.md:303-304` (transversal rule), `M-07-pricing-fee-read/F-01-fee-read-resolver/feature.md:90`.

**Offending token, verbatim** (`mission.md:298`):

> `edita região pricing existente (única exceção — hub arbitra)`

and (`M-07/F-01/feature.md:90`) `região pricing root.go `:828-858``.

**Why it is a defect.** M-07 deletes `pricing/adapters/tarifflive/` and collapses `tariffcomposite/`; those packages have exactly four references in the tree, and two of them — `root.go:99` and `root.go:101` — are **import lines outside `:828-858`**. The declared range is therefore a total claim that the milestone's own work falsifies; the file will not compile if M-07 stops at `:858`. The same gap is structural rather than M-07-specific: every milestone whose cell reads "1 linha ancorada" (M-03, M-04, M-06, M-08, M-09) must also add an import, and the import block is the one region of `root.go` that no cell owns and no rule partitions — exactly the seam R-7 (`mission.md:359`) is written about. The transversal rule "hub = resolver of record" arguably absorbs the mechanical conflicts, but it is stated for constructor regions, not for imports, and M-04∥M-03 (lane B) and M-06∥M-07 (lane C) will hit it concurrently.

**Yes-if.** `mission.md:303-304` states that the `root.go` **import block is hub-resolved and owned by no milestone** (each milestone adds/removes only the imports its own anchored region requires), and M-07's matrix cell plus `M-07/F-01` Ownership name `root.go:99,101` (import removals for `tarifflive`/`tariffcomposite`) alongside the `:828-858` region.

---

### F-r04-6 — ADVISORY — checks 2, 3 — the write-DAG's universal "no residual overlap" claim omits two real overlaps

**Locus:** `mission.md:320-322`.

**Offending token, verbatim:**

> `Nenhum write-set overlap restante sem edge nomeado.`

**Why it is a defect.** The enumeration immediately above it lists three cross-milestone collisions (M-05→post-M-04, M-06→post-M-03, M-07's `root.go` pricing region) and then closes with a universal. At least two further overlaps exist and are not enumerated: (a) the **guard allowlist file** created by M-02 F-04 (`M-02/F-04/feature.md:25-26,49` — "allowlist explícita de exatamente 4 entradas", "allowlist como dado no próprio teste") is written by M-03 F-03 ("allowlist (remoção A/B)", `M-03/F-03/feature.md:76`) and by M-07 F-01 ("allowlist (C/D)", `M-07/F-01/feature.md:90`), i.e. two milestones write a third's file, and neither appears in the M-03 or M-07 `Go packages/files` matrix cells (`mission.md:294,298`); (b) the `root.go` import block per F-r04-5. Both overlaps are in fact serialized — lane B precedes lane C, so M-03's removal always lands before M-07's — so the plan is sound and this is advisory; what is false is the totality of the claim, which is the artefact a P6 gate will read as the disjointness proof.

**Yes-if.** `mission.md:320-322` enumerates the guard-allowlist file (owner M-02 F-04; writers M-03 F-03 and M-07 F-01, serialized by the lane B → lane C edge) and the `root.go` import block, or replaces the universal with a scoped statement naming the axes it actually covers.

---

### Checks passing with no finding

- **Check 1 (DAG completeness/justification).** Every edge in `mission.md:272-286` is backed by a named forcing artefact; each is real against the briefs. No false edge: `M-08→M-09` is genuinely hard (M-08 F-02 cannot call a setter that does not exist) and is honestly labelled "na prática sempre satisfeita" since M-09 is lane A; `M-07⤳M-05` and `M-09⤳M-04/M-06` are correctly typed as data-quality, not compile. No missing edge found: M-05's dependency on M-01 F-02's fee reader is transitive through M-04→M-01 and is declared as such (`M-05/milestone.md`, "(M-01 transitiva via M-04.)"); M-05 F-02's read of `products_mirror` targets a MIS-006 asset that already exists.
- **Check 2 (six-axis disjointness), other than F-r04-5/6.** Migrations are disjoint and all above the repo's highest (`0085`): M-02 `0086-0089`, M-04 `0090-0092`, M-08 `0093`, M-06 `0094-0095` (reserve, gap ratified at `mission.md:316`); out-of-order application is safe per the migration runner's set-difference semantics (retired in r01). Shared *tables* are partitioned explicitly and the partitions hold: `channel_fees` M-05 (camada 2) vs M-06 (camada 3) vs M-07 (read-only), `divergences` M-05 (`kind=estoque`) vs M-06 (tarifa) — disjoint natural keys, and the partial unique `WHERE resolved_at IS NULL` (0087) separates rows by kind. `orders/**` is held by M-03 in lane B and inherited by M-06 in lane C, never simultaneously (`mission.md:308-309`). `sync/application/` carries the registered additive lock (M-02 F-03 owns `scheduler.go`; M-09 is `health_*` + `sync/transport/**` only, `mission.md:312-315`). FE surfaces are disjoint (/pedidos M-03→M-06, /anuncios M-05, /precos M-07, IntegracoesPage M-09).
- **Check 3 (feature write-DAG), other than F-r04-6's enumeration.** Every cross-milestone write-set overlap I could construct carries either a serial lane edge or the recorded additive-lock grant; within milestones the internal DAGs are stated and consistent (notably `M-05` F-01∥F-02→F-03 with disjoint write sets `channel_fees` vs `divergences`, and `M-08` F-01→F-02).
- **Check 4 (contract satisfiability), other than F-r04-1.** The camada-2 fee source is a closed loop, not a dangling promise: `M-01/F-02` Expected Output commits to verifying `?context=channel_marketplace` on the multiget and, failing that, to shipping the dedicated `GET /items/{id}/prices` reader in the same package, so `M-05/F-01`'s "fonte única pronta / Blockers: none" is satisfiable and no ownership boundary is crossed. The allowlist arithmetic closes exactly (4 entries: M-03 removes A/B, M-07 removes C/D). IC-01/IC-03 keep camada 3 out of M-03 and inside M-06 consistently.
- **Check 5 (prerequisite existence), other than F-r04-3/4.** Every anchor I re-checked is accurate at the cited line, including ones no prior round verified: `route_deadline.go:23-28` (the `RouteClass`/deadline const block), `root.go:259-272` (`registerBatchRoutes`, the negative control for the deadline-class proof), `repository.go:390-394` (the MASS-CLOSURE `UPDATE listings SET status='closed'`), `writer.go:74-105` (upsert + keep-absent idiom), `listings_test.go:25` (migration regex idiom), `order_repo.go:467-514` / `:603`, `import_service.go:46`, `buyer_fiscal_reader.go:19-21`, `enrich_service.go:192`, `tariff_resolver.go:30-32`, `root.go:844-850` (M-07's rewire target is exactly the `var tariffResolver` + `if feeReader` block). `0075_sync_sync_state.sql` really does carry `cursor jsonb`, `last_full_sync_at`, `last_incremental_at`, `last_error jsonb`, `consecutive_failures integer NOT NULL DEFAULT 0`, and has no `phase` column — consistent with IC-05 deriving `phase` from the cursor.
- **Check 6 (ADR/IC propagation), other than F-r04-1/3.** No widening or narrowing found in the briefs relative to the ratified ADR/IC set; the ADR-11 amendment (dedupe only with `_id`), ADR-13 (PK sentinel untouched), ADR-03 (no `raw` for `billing_info`) and ADR-05 (zero ML calls on the read path) propagate verbatim wherever restated.
- **Check 7 (artifact-class completeness).** All 25 feature briefs carry Brief, Inputs, Expected Output, Constraints, Negative Scenarios, Ownership, Validation Expectations and Handoff; all four FE briefs carry `## Interaction Model` (`M-05/F-03`, `M-06/F-03`, `M-07/F-02`, `M-09/F-02`); the briefs without `## Inputs/Outputs` are internal or test-only surfaces with no asymmetric I/O shape. Validation Expectations name inspectable minimums throughout (SELECT-verified row counts and column values, golden JSON with exact nulls, two-direction divergence fixtures, `>1 page` fixtures per R-3, must-fail tokens naming the failure, deadline-class proof with the negative control outside `registerBatchRoutes`).
- **Check 8 (no implementation planning, no new scope).** Briefs stop at seam and wiring level; open engineering values are explicitly deferred to spec (`M-08/F-02` inbox poll interval, `M-09/F-02` refetch cadence, `M-09/F-01` repo read shape). No product scope beyond the ratified design: the ngrok callback runbook, the `orders_v2`-only topic, the `/sync/health` new-endpoint decision and the fee-column exclusion from E3 are all traceable to ratified ADRs/ICs, and the operator-authorization gate for the ML callback registration is stated in both `M-08/milestone.md:68-70` and `F-02`.

---

## 4. VERDICT

**NEEDS-REVISION.**

Blocking findings (1):

- **F-r04-1** — `M-09/F-01/feature.md:57-60` pins `last_success_at = last_full_sync_at`, which `sync_state_repo.go:62,74-79` freezes for any entity running incrementally once M-02 F-03 lands; `M-09/F-02` renders that field as the freshness badge, so /integracoes reports a healthy 5-minute sync as stale or never — the exact failure `M-09/milestone.md:66-67` forbids and the criterion Q4 gates in a browser. IC-05 never defines the field, so there is no higher authority to appeal to at execution time. Yes-if: define `last_success_at` in IC-05 as the most recent success of any kind (`GREATEST(last_full_sync_at, last_incremental_at)`, NULL only when both NULL) and drop the equation from the brief, with a negative fixture proving an incremental-only entity renders fresh.

Advisory findings (5): F-r04-2 (setter/builder arm silently drops the M-08 injection), F-r04-3 (IC-03's false `DeriveOrderBucket` locus and phantom move), F-r04-4 (M-08's `user_id`→installation map cited from a one-directional resolver), F-r04-5 (`root.go` import block unowned; M-07's declared range excludes `:99,101`), F-r04-6 (write-DAG universal omits two real, though serialized, overlaps). Advisories do not gate.

Part A is clean: all five r03 findings are genuinely closed, the chosen arms hold against ratified scope and against code, and no earlier closure was reopened.

Standing gate note, unchanged by this verdict: the operator waiver (`planning-reviews/sol-unavailable-p3-r01.md`) makes the GPT-5.6 Sol P5 touchpoint retroactively mandatory from 2026-08-05, before `status: planned`.
