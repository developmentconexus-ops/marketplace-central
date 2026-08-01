# P3 counter-proposal — Opus frio (waiver Claude-crew, seat substituto do Sol) — r01

Proveniência: subagente Opus independente, cego ao candidato Claude
(`p3-claude-candidate-r01.md` negado), mesmo manifest `p3-input-r01.sha256`
(drift disclosed: mission.md mudou pós-freeze por churn de status + waiver — mismatch
esperado e isolado, verificado pelo próprio seat). Output VERBATIM abaixo, colado como
primeiro ato após a notificação (regra profile §11).

---

## Manifest verification

Recomputed SHA-256 over the five manifest entries at `.mnfs\MIS-007-ml-sync\planning-reviews\p3-input-r01.sha256`:

| File | Manifest digest | Recomputed | Verdict |
|---|---|---|---|
| `.mnfs/MIS-007-ml-sync/mission.md` | `55bdfc37…6f54e` | `f6b9a1b9…02dee` | **MISMATCH** (expected — drift disclosed in the brief) |
| `.mnfs/MIS-007-ml-sync/research/codebase-ingest-side.md` | `1e9710a0…c1861678` | identical | MATCH |
| `.mnfs/MIS-007-ml-sync/research/codebase-read-side.md` | `cde669cb…95efbef` | identical | MATCH |
| `.mnfs/MIS-007-ml-sync/research/external-ml-api-facts.md` | `1633ab8b…b06c38d2` | identical | MATCH |
| `docs/design/MIS-007-ML-SYNC-DESIGN.md` | `d39ddba9…38a353` | identical | MATCH |

Top-digest `50c997c3…7876281e` therefore does not recompute; the divergence is isolated to `mission.md`. The mission.md I read carries `status: draft`, `planning_phase: scope`, gate P1 answered 2026-07-31 (4 rows incl. the "Apply" write-mode waiver), and the P2-corrected fee outcome — consistent with the disclosed churn. **I treated the manifest as the file list and mission.md's current content as authoritative for P1 answers.** No other file was opened; the forbidden `planning-reviews/` siblings were not read.

One correctness note carried forward: the ratified design at `docs/design/MIS-007-ML-SYNC-DESIGN.md:86` still cites `fee_sync.go:29` as a live seed. Research refutes it (`codebase-read-side.md:151-155`). mission.md already absorbed the correction; the **design doc has not been amended**. Any milestone brief that copies §4's "o seed estático 16%/22% (`fee_sync.go:29`) morre" would be shipping a false claim about the repo.

---

## Architecture spine

Fourteen ADR-lite entries. Each is a choice two downstream workers could make incompatibly and only discover at merge or at live-drive.

### A-01 — ML resilience is a decorator on the single choke point, not per-call-site retry
- **Decision.** Backoff (exponential + jitter), `Retry-After` honoring on 429, and a per-installation token-bucket limiter are implemented **once**, wrapping `doRawWithHeaders` (`capability_adapter.go:712-731`) — the one function every ML request already passes through. Budget is shared across goroutines (one limiter per `ProviderAccountRef`), not per-caller. 429 stops returning `ErrCodeProviderRateLimited` on first contact (`:654-655`, `:462-465`, `:578-583`); it returns it only after the retry budget is spent, and the exhausted error names attempts + last `Retry-After`.
- **Prevents.** Two milestones (listings backfill, orders backfill) each writing their own retry loop with different semantics; and the real failure mode — independent goroutine pools each staying "under" ~1500 req/min while their *sum* is not (`external-ml-api-facts.md:21`, fact #11 is `assumed`, so the limiter must be configurable, not a compiled-in constant).
- **Must preserve.** `AccessTokenResolver` injection (`capability_adapter.go:26`, wired `root.go:370-378`) and the existing OAuth refresh backoff (`refresh_policy.go:18-27`) are a *different* mechanism — do not merge them. Non-429 error mapping and `providerDiag` 512-rune clipping (`:681-688`) stay.
- **Trade-off.** A choke-point decorator cannot express per-endpoint policy (e.g. "never retry a write"). Accepted: this mission is read-only; writes (`PUT /items/{id}`, `:444`) get an explicit no-retry opt-out flag now so the future write mission does not have to unpick it.
- **Validation impact.** Must-fail: a fake transport returning `429 + Retry-After: 2` must make the test **name** the elapsed wait; a test asserting only "eventually succeeds" is a vacuous pass. Rate-limiter proven by a fixture issuing N concurrent calls and asserting the observed request timestamps, not by reading the config.

### A-02 — New ML endpoints go in new files in the `mercadolivre` package, never appended to `capability_adapter.go`
- **Decision.** `capability_adapter.go` is frozen to the resilience change (A-01) plus multiget. Listings-side endpoints (`/items?ids=`, `/items/{id}/shipping_options`) and orders-side endpoints (`/shipments/{id}/sla`, `/costs`, `billing_info`, `/packs/{id}`) land in **separate files**, following the precedent already in the package: `shipping_reader.go`, `buyer_fiscal_reader.go`, `catalog_match_reader.go`.
- **Prevents.** The single hottest merge collision in the mission. Every milestone needs an ML endpoint; a 1100-line god-file with two parallel workers appending methods produces conflicts on every merge, and conflicts in a file that also contains A-01's retry logic are the dangerous kind.
- **Must preserve.** One `ProviderCapabilitySet` declaration (`:79-92`) — capability registration remains centralized even though implementations spread.
- **Trade-off.** More files, and the capability struct grows collaborator fields. Cheaper than serialized dispatch.
- **Validation impact.** Gate check is mechanical: `git diff --stat` on `capability_adapter.go` outside M-01 should be ~zero lines.

### A-03 — Raw payload is carried in the DTO, persisted **selectively**, and never persisted for PII-bearing resources
- **Decision.** DTO rule `Raw json.RawMessage` applies to every ML DTO (design §6). Persistence is narrower: raw is stored in an additive `raw jsonb` column only for `listings`, `orders`, `order_shipments`. **`billing_info` raw is NOT persisted** — it carries buyer document + fiscal address (`buyer_fiscal_reader.go:59-94`; fact #16). The existing `raw_provider_ref` (path-only, `capability_adapter.go:362`) is kept and extended, not replaced.
- **Prevents.** A worker reading "regra DTO Raw json.RawMessage persistido" (design §6, line 150) literally and dumping fiscal payloads into Postgres — against AGENTS.md's PII bar, and compounding the already-open PII-scrub pendency on `docs/design/evidence/ml-api/` (`external-ml-api-facts.md:7`).
- **Must preserve.** `io.LimitReader(resp.Body, 1<<20)` (`:747`) — keep a cap; raw must be bounded, and oversize raw stores a truncation marker rather than silently a prefix (a truncated JSON blob that still parses as text is a lie).
- **Trade-off.** Reprocessing a billing_info bug requires a refetch. Acceptable — refetch is authenticated and idempotent.
- **Validation impact.** A test asserts no persisted column ever contains the buyer document string, driven from a fixture that *does* contain one.

### A-04 — Ingest is resource-addressed and single-path; enumeration is a separate, replaceable concern
- **Decision.** For each entity there is exactly one writer with signature `Ingest{Order,Listing}(ctx, tenant, installation, providerResourceID) error`, idempotent by upsert on the natural key. **Every** producer — backfill, scheduler reconciliation, webhook worker, manual refresh button — calls that one function. Enumerators (scan/`scroll_id`, `orders/search` with `date_last_updated` + `sort=date_desc`, inbox drain) only produce IDs.
- **Prevents.** Exactly the defect the mission exists to erase: today `orders` has two writers with different semantics (`import_service.go:46-82` persists; `enrich_service.go:150-162` computes at read and throws away). A wave-shaped split naturally regrows this — the webhook milestone would otherwise write its own upsert.
- **Must preserve.** The canonical upsert idiom: `ON CONFLICT (natural key) DO UPDATE` + keep-absent, never DELETE (`internal_read/adapters/mirror/writer.go:74-105`); `orders_marketplace_order_items` conflict target is `(…, mpc_line_id)` and `mpc_line_id` is trigger-immutable (`0033:20-34`).
- **Trade-off.** Per-resource ingest forgoes batch-SQL efficiency. Mitigated by multiget hydration (20/call) feeding N ingest calls inside one transaction boundary per batch.
- **Validation impact.** Q3's criterion is discharged only against this seam: re-deliver the same notification → row count unchanged and `updated_at` moves; kill backfill mid-run → resume produces no duplicates. Both must be run against the *same* function the webhook uses, or they prove nothing about the webhook.

### A-05 — Read-path decoupling is a **shrinking allowlist guard**, and no live source is deleted before its persisted replacement lands in the same merge
- **Decision.** This is the counter-proposal's load-bearing rule. Two parts:
  1. An architectural test lands in the **first** schedulable milestone: it enumerates every read-time site that may touch the ML client, as an explicit allowlist of 4 entries (`GET /orders`, `GET /orders/{id}`, `POST /pricing/decompose`, `POST /pricing/solve` — `codebase-read-side.md:20`). Any *new* read-time ML dependency fails immediately. Each milestone that lands a persisted replacement **deletes its own entry** in the same commit. The guard reaches zero at mission close; Q1's "nenhum sítio de transport importa client ML" is then a measured fact, not an end-of-mission hope.
  2. **A milestone may not remove a live read source unless its persisted source ships in the same merge.**
- **Prevents.** The regression window a literal Onda-0-first split creates. Deleting `enrich_service.go:150-162` before `order_shipments` exists blanks, on the operator's daily screen: the ratified `SLA` column (`PedidosTable.tsx:95-146`), the Fila/Kanban `bucket` (`http_handler.go:549`, `DeriveOrderBucket` consumes *live* shipment status), `frete_real`, `destinatario`, `destino_uf/cep`, `rastreio` (`:557-567`), and in the drawer `comprador_fiscal` (`:462-467`). That is a multi-milestone outage of the mission's own success criterion ("tudo está lá"), caused by the plan, not by a bug.
- **Must preserve.** The end state is non-negotiable and ratified: zero ML calls in read, `<2s`. This ADR changes only *when* each of the 4 dies, never *whether*.
- **Trade-off.** The two ML-in-read families die at different times instead of in one "Onda 0" chip, so the mission is never in a clean "phase 0 complete" state. Bought with: no user-visible regression at any merge point. Asymmetry is real and evidence-based — sites **C/D** (pricing) *can* die early, because their fallback already exists and is honest (`pricing_tariff_defaults` degrau-4 with provenance `config`, per mission.md:36-37); sites **A/B** (orders) cannot, because they have no fallback at all.
- **Validation impact.** Every milestone that deletes an allowlist entry must show before/after on the *same* browser screen: the field still renders the same value, now sourced from Postgres. Deleting the entry without the screenshot pair is not acceptance.

### A-06 — MASS-CLOSURE dies by separating "absent from this pull" from "closed at ML"
- **Decision.** `ApplyCompletedPull` (`listings/adapters/postgres/repository.go:383-465`) stops issuing the unfiltered `UPDATE … SET status='closed'` (`:390-394`). Replacement adopts the ERP mirror contract verbatim in shape: per-row upsert keyed by the listing natural key, plus a run-scoped keep-absent marker (`absent_in_last_snapshot` / `stale_since` analog, `writer.go:97-105`, ADR-04). `status` becomes provider truth only — it is set from `/items` payload, never inferred from absence. Marking-absent is applied **only after a run is declared complete**, and a run truncated by 429/deadline/cancellation is never complete.
- **Prevents.** The single hardest incompatibility in the mission: *resumable backfill and MASS-CLOSURE cannot coexist*. Design §6 mandates "retomável por cursor"; today a partial pull closes the entire catalog. A listings worker that ships the cursor before the closure fix ships a catalog-wiper on the first 429 — which A-01 makes *survivable* but does not make *impossible*.
- **Must preserve.** `ApplyCompletedPull` remains the **sole** production writer of `listings`. The replacement does not add a second writer; it changes the semantics of the existing one.
- **Trade-off.** A listing genuinely closed at ML but never returned by scan lingers as `absent` rather than `closed` until a full sweep confirms. That is the honest state and matches ADR-04; the alternative fabricates provider truth.
- **Validation impact.** Must-fail with the failure named: seed N listings, run a backfill that aborts after page 1, assert **zero** rows flipped to `closed`. Fixture must exceed one page (R-3; `listings_read_test.go` cursor precedent).

### A-07 — Backfill cursor contract: `nil` means "erase", so a finished backfill returns a terminal cursor
- **Decision.** Both new jobs return a non-nil cursor whose JSON carries an explicit phase (`{"phase":"backfill","scroll_id":…}` → `{"phase":"incremental","watermark":…}`). Returning `nil` on completion is forbidden.
- **Prevents.** `scheduler.go:42-45` — a returned nil cursor **deletes** the persisted cursor. Two workers each writing "backfill done → return nil, we're caught up" produce an infinite backfill loop that re-scans the whole catalog every tick, silently, forever. This is invisible in unit tests and only shows as rate-limit burn in the live-drive.
- **Must preserve.** `JobFunc` signature (`scheduler.go:46`); `RegisterJob` one-job-per-entity fail-closed (`:85-101`); `Read`/`RecordSuccess`/`RecordFailure` seam (`:22-34`) and the atomic `consecutive_failures+1` (`sync_state_repo.go:95-102`). The `entity` enum already contains `listings` and `orders` (`0075:12-14`) — **no migration is needed** to register them; a milestone proposing one is wrong.
- **Trade-off.** Phase lives in the opaque cursor rather than a typed column. Keeps 0075 untouched, which the transversal "additive only" rule prefers.
- **Validation impact.** Kill-and-resume test asserts the persisted `cursor` JSONB is non-null after a successful terminal run. Byte-exact string comparison on round-tripped JSONB is banned.

### A-08 — Two scheduler instances, not one loop that learns cadence
- **Decision.** Cadences differ (orders 5 min, listings daily). Instantiate **two** `Scheduler`s with different intervals, following the existing composition pattern (`synccomposition.NewProductsScheduler`, `root.go:672-677`). The `schedule jsonb` column (0075) stays unread this mission; making the loop cadence-driven is explicitly out of scope. Separately and **once**, in the foundation milestone, fix `scheduler.go:160`'s hardcoded `incremental=false` so `last_incremental_at` becomes real.
- **Prevents.** Two workers independently "improving" the shared scheduler loop — one teaching it to read `schedule`, one adding a per-entity interval map — colliding in `sync/application/scheduler.go`, a file neither milestone owns.
- **Must preserve.** Cadence-agnostic `Start` (`:105-119`), per-entity failure isolation (`:124-139`), and "cursor read error skips the cycle, never fabricates nil" (`:141-161`).
- **Trade-off.** Three goroutine tickers at boot instead of one. Negligible; `root.go:661-663` already runs three.
- **Validation impact.** `last_incremental_at` populated is a precondition for Q4's /integracoes health criterion — if the `incremental=false` fix is not in the foundation milestone, the observability milestone will render a permanently-empty field and QA will (correctly) fail it.

### A-09 — Fee layers: exactly what dies is enumerated, and the already-dead is not re-killed
- **Decision.** Given P2, the ledger is:
  - **Already dead — no work, and no milestone may claim it.** The flat 16%/22% ML seeder (`registry/mercado_livre.go:46-48` no-op; `registry_test.go:90-103`; rows dropped by `0081:15-16`). `TestBuildFeeSyncersSeedsNothing` asserts `len(BuildFeeSyncers()) == 0`. `fee_sync.go:29` **does not exist**. The 22% figure never existed in production — only in `fee_schedule_service_test.go:178-179`.
  - **Dies in this mission.** (a) The live degrau-3 at read time: `pricingtarifflive.NewResolver` wiring (`root.go:845-851`) + `tarifflive/resolver.go:43-69`, replaced by a `channel_fees` lookup. (b) `auth_adapter.go:47-48` `"baseline_commission_percent": 0.16` — a number with no provenance, no layer, no collection date; it either dies or is re-expressed as a `channel_fees` row with `origem='config'`.
  - **Survives, re-labelled.** `pricing_tariff_defaults` 13/16 (`0068:13-14`, materialized `calc_repository.go:240-246`, editable via `PutTariffDefaults`) becomes the explicit lowest-strength fallback carrying provenance `config`.
  - **Neither adopted nor extended.** `FeeSyncScheduler` (15 min, `root.go:663`) and `RegisterFeeSyncerFactory` — machinery with no production caller that returns `INTEGRATIONS_FEE_SYNC_UNSUPPORTED` (`marketplace_executor.go:62-66`). `channel_fees` is fed by the listings/orders ingest, **not** by the marketplace fee-schedule syncer. Naming this prevents a worker from "finishing" the dead path.
- **Prevents.** A milestone brief asserting a total claim ("the static seed dies") that is already true → an unfalsifiable acceptance criterion, and a worker burning a milestone deleting code that is not there.
- **Must preserve.** `order_items[].sale_fee` is **per unit** (fact #1, live-measured, undocumented) — layer 3 arithmetic is `sale_fee × quantity`. `FeeSource: "api_sync"` (`registry/mercado_livre.go:20`) stays honest. `MissingSaleFee` flagging in profitability (`service.go:1014-1015`) must keep firing rather than being satisfied by a fabricated `channel_fees` row.
- **Trade-off.** `channel_fees` ships with layer-1 columns unpopulated (MIS-008). Accepted per mission.md:99-101 — an unconsumed column forces no rework.
- **Validation impact.** Every fee shown on a screen must render its provenance triple (layer, origin, collected-at). A displayed number with no provenance is a milestone failure regardless of correctness.

### A-10 — `divergences` is one-open-row-per-(entity, kind), upserted, with **both** sides' observation timestamps
- **Decision.** Per P1 (dedicated table). Schema fixed centrally, before any producer: `(tenant_id, provider, entity_type, entity_id, kind)` as the natural key with at most one unresolved row; columns `expected_value`, `observed_value`, `expected_source`, `observed_source`, `expected_observed_at`, `observed_observed_at`, `detected_at`, `resolved_at NULL`. Detection is **upsert** (open or update the open row); convergence sets `resolved_at` (mission.md:104-106). `kind ∈ {estoque, tarifa}` with a CHECK, extensible additively.
- **Prevents.** The textbook naive-split break: the stock producer (listings ingest, daily) reasonably writes append-only event rows; the tariff producer (orders ingest, every 5 min) reasonably writes one-open-row. Both are defensible; together they make the "divergentes" filter and the badge count meaningless, and a 5-minute append-only writer generates unbounded rows. One shape must be decided before either producer exists.
- **Must preserve.** R-5's mitigation is a **schema** requirement, not a code requirement: both timestamps are NOT NULL, because "ERP sellable cut and ML `available_quantity` were observed at different times" is the false-positive mechanism and it is unprovable without them.
- **Trade-off.** History of a flapping divergence is lost (only current + last resolution). Recoverable additively later via an events sibling; not worth the row volume now.
- **Validation impact.** Design §9's two-direction proof binds here: create divergence → badge appears; converge → `resolved_at` set and badge disappears **on the same screen in the same drive**. A test that only asserts appearance is a half-proof.

### A-11 — Webhook: the payload is a pointer, never data; the endpoint always answers 200
- **Decision.** Per P1 (untrusted hint + IP log-only + `orders_v2` only). `POST /webhooks/{provider}` is thin transport: bounded body read, insert into `notifications_inbox`, respond 200 — always, including unknown topic, unparseable body, and unknown `user_id`. Only `resource` and `topic` are used, and only as a pointer to an **authenticated refetch** through A-04's ingester. `user_id` is mapped to an installation via our own credential store; a `user_id` we do not own is recorded and dropped. Source IP is persisted and compared to the official list (fact #9) **log-only, never a reject** — the list is mutable. Dedupe on `(provider, topic, resource, notification_id)`. Route registered method-aware (`"POST /webhooks/{provider}"`, the `market/transport/http_handler.go:56` idiom) as **interactive class — explicitly not in `registerBatchRoutes`** (`root.go:259-272`).
- **Prevents.** (a) A forged notification injecting values — nothing from the body reaches a domain table. (b) Non-200 responses triggering ML's 8-retries/~1h storm (fact #7) against our own endpoint. (c) The IP allowlist quietly becoming an enforcement rule and going deaf when ML rotates IPs. (d) The endpoint landing in the batch class, where a 120s deadline invites synchronous processing inside the request.
- **Must preserve.** This is the system's **first unauthenticated public write surface**. Body size cap; no error-envelope detail leakage (moot — it returns 200); the `apierror` envelope contract (`apierror.go:25`) still applies to the few non-200 shapes if any exist. Scheduler reconciliation stays mandatory this mission (design §5) — missed_feeds' 2-day window (fact #7) is not a safety net for a >2-day gap.
- **Trade-off.** Always-200 means a broken inbox insert is invisible to ML. Compensated by /integracoes surfacing "last notification received" (Q4) plus the 5-minute reconciliation.
- **Validation impact.** Q2 is discharged by a **forged POST** carrying a plausible-but-wrong `resource` — assert zero domain-table writes and a persisted inbox row with an off-allowlist IP logged. Live-drive: mutate a real listing/order → notification → inbox → ingest → screen (design §9).

### A-12 — Migration numbers are pre-allocated by the hub in disjoint ranges; never picked by scanning
- **Decision.** Next free is **0086** (`codebase-ingest-side.md:47-48`). Each milestone brief carries an explicit reserved range (e.g. `0086-0087`, `0088-0089`, …). Workers use their assigned numbers and no others. All migrations additive; no destructive `ALTER` on `products_mirror` or `listings`; each new migration ships a regex assert test in the existing style (`migrations/listings_test.go:25`, `product_link_decisions_test.go:23`).
- **Prevents.** Two parallel worktrees both taking 0086 — a conflict invisible until merge, and one that a rushed resolution "fixes" by renaming a file whose name is the primary key of `schema_migrations` (`runner.go:33-40`). A renamed already-applied migration re-runs.
- **Must preserve.** `cmd/server` does not migrate on boot (`main.go:15-55`); migration is `cmd/migrate` only. `0021` is already duplicated — do not "clean it up" in this mission.
- **Trade-off.** Reserved-but-unused numbers leave gaps. The tree already has gaps (0002, 0040-0042, 0048-0049, 0054, 0059-0064, 0077); harmless.
- **Validation impact.** Merge-time check: no two files with the same numeric prefix; every applied filename immutable.

### A-13 — `listing_variations` is additive; `listings`' PK does not move this mission
- **Decision.** Formalize `listing_variations` as a child table keyed `(tenant_id, installation_id, provider, provider_listing_id, variation_id)`, populated by the same sole writer (A-06). `listings` **keeps** its current PK including `variation_id` with the `'-'` sentinel (`0036:2-31`) for this mission; the child table is the new home for per-variation facts, and the read model prefers it when present.
- **Prevents.** A PK change on `listings` is a destructive change radiating into: the sole writer, the read repo backing /anuncios, `product_link_listing_snapshots` (PK carries `provider_variation_id`, `0022:1-17`), and `product_links` (PK carries `provider_variation_id`, `0025:1-15`). Design §7 says "listing_variations formalizada" — that is satisfiable additively, and the additive reading is the one compatible with the transversal "migrações aditivas, nunca ALTER destrutivo".
- **Must preserve.** Identity anchors already resolved outside `listings` must not be orphaned. Critically: `Source.ReadPage` feeds `AbsorbProviderSnapshots` (`listings/adapters/connectors/source.go:17-19, 54-89`), which resolves the EAN/SKU chicken-and-egg **before** variations are flattened. Design §6's "coletar ids primeiro, hidratar depois" restructures exactly this path — the re-vínculo step (design §6, Onda 1 item 5) silently loses its input if the observer is not re-fed from the new hydration stage. This is a hidden coupling between the backfill rewrite and product-link quality.
- **Trade-off.** Two places carry variation identity during this mission. The consolidating PK move becomes a named future item rather than a mid-mission destructive migration.
- **Validation impact.** After the backfill rewrite, assert `product_link_listing_snapshots` row count and anchor content are non-regressive versus the pre-change pull — a must-fail that names "snapshot observer starved".

### A-14 — Composition root and the OpenAPI/SDK pair are hub-owned serialized seams
- **Decision.** (a) Each milestone's wiring enters `root.go` as **one** call to a module-local composition constructor (the `synccomposition.NewProductsScheduler` precedent, `root.go:672-677`) at a named anchor region; the hub arbitrates `root.go` and is the resolver of record for any conflict there. (b) `contracts/api/marketplace-central.openapi.yaml` + `packages/sdk-runtime/src/index.ts` land in the **same commit** (`HARNESS-PROFILE.md:248`) under the hub-held contract lock (`:193`, `:202`); **at most one milestone with FE-facing route/schema changes is in flight at a time.**
- **Prevents.** `root.go` is touched by literally every milestone (`:591-599` orders enrich, `:701-714` listings, `:728-736` refresh, `:845-852` pricing, `:672-677` scheduler) — it is the guaranteed collision. And the SDK is **hand-written**, not generated (`codebase-read-side.md:105`): the client is one object literal (`index.ts:2113-2330+`), so two milestones adding methods conflict in one region, while YAML path additions in disjoint alphabetical regions usually would not. The generated-SDK intuition is wrong here and would produce a bad parallelism call.
- **Must preserve.** Error-code unions are sourced verbatim from the YAML `code` enum (`index.ts:1727-1796`); `hasCode<C>()` validates the `details` fields the spec guarantees (`:1897-1909`, `REQUIRED_DETAIL_FIELDS` `:1867-1873`). New codes enter the union and the spec together or `hasCode` lies. Byte-identical key locks between handler DTO, OpenAPI schema, and SDK interface (`orders/transport/http_handler.go:290-296`) hold for every extended payload.
- **Trade-off.** Contract-touching milestones serialize even when their code is disjoint. This is the constraint that shapes the wave plan below more than any code dependency.
- **Validation impact.** Gate check per merge: does the diff change the YAML without the SDK, or the SDK without the YAML? Either alone is a blocker.

---

## Milestone split

Nine milestones. The cut is **not** wave-per-milestone: the design's three waves are re-cut along *seam ownership* (one owner per shared file) and *replacement-before-deletion* (A-05), which yields three genuinely parallel lanes where a wave split yields one serial chain with a regression window in the middle.

| id | headline | why here |
|---|---|---|
| **M-01** `ml-client-hardening` | Backoff+jitter+Retry-After+shared limiter+raw-DTO at `doRawWithHeaders`, plus multiget `/items?ids=`. Zero schema, zero UI. | Design §6 makes it the prerequisite of both waves. Every backfill multiplies request volume against an adapter with literally zero retry; landing it after any backfill means the backfill is written against a client whose failure semantics then change under it. |
| **M-02** `sync-core-seam` | The cross-cutting contracts: `channel_fees`, `divergences`, `order_shipments` migrations (0086-0088); the fee-layer-resolution and divergence-recording ports; the scheduler `incremental` fix; the shrinking read-guard allowlist test. One owner, no provider code. | This is the counter-proposal's main structural move. Extracting the shapes two ingest milestones would otherwise each invent (A-10, A-09, A-04) is what makes M-03 ∥ M-04 and M-05 ∥ M-06 ∥ M-07 safe. Independent of M-01 — can run beside it. |
| **M-03** `orders-shipment-persist` | `order_shipments` populated (shipment + `/sla` + `/costs` + `billing_info`) via the resource-addressed ingester on the existing import path; **sites A and B die**; bucket/SLA/frete/rastreio/comprador-fiscal now read from Postgres. | The inversion of "Onda 0 first". Persisting the enrichment *is* the decoupling of /pedidos; doing it in one milestone means the ML-in-read deletion never outruns its replacement (A-05). |
| **M-04** `listings-backfill-ingest` | Scan-ids-then-multiget backfill with resumable cursor, MASS-CLOSURE replacement, E3 columns, `listing_variations` additive, daily scheduler, `available_quantity`, bulk manual refresh. | Owns the `listings` writer end-to-end. Fully disjoint from M-03's module — the two are the primary parallel pair. |
| **M-05** `listings-fees-divergence` | Layer-2 fee write (`listing_prices` requeried with the ingested `category_id`/price + `shipping_options`), stock divergence detection at ingest, `/anuncios` columns + ⚠ badge + "divergentes" filter. | Split out of M-04 so the ingest milestone stays QA-able and so the /anuncios FE surface has its own browser-driven contract. Consumes M-04's output; touches no orders code. |
| **M-06** `orders-backfill-decomposition` | 12-month backfill (`date_last_updated` + `sort=date_desc`), 5-minute incremental, **persisted** margin decomposition with frozen cost, layer-3 fee write, tariff audit 3→2 divergence, indexed bucket. | Extends the ingester M-03 established rather than competing with it for the same files. |
| **M-07** `pricing-fee-read` | Pricing resolver reads `channel_fees` layers with provenance; degrau-3-live dies; `pricing_tariff_defaults` relabelled as `config` fallback; **sites C and D die**. | Pulled out of the "kill all 4 sites" chip because the pricing resolver chain (`root.go:845-851`, `tarifflive/resolver.go`, `pricingtariffcomposite`) is a distinct seam shared with `POST /pricing/simulations/batch` (`root.go:852`) — one owner, and not the orders owner. |
| **M-08** `webhook-ingest` | `notifications_inbox` + `POST /webhooks/{provider}` + in-process worker + callback registration on the ngrok fixed domain, `orders_v2` only. | Last on the orders chain because the worker's only job is to call M-06's ingester (A-04). Building it earlier means building a worker with nothing to call. |
| **M-09** `sync-observability` | `/integracoes` sync-health-per-entity + webhook status section. | The one genuinely free-floating lane: `sync_state` (0075) and `GET /sync/runs` already exist and `listSyncRuns` has **zero FE consumers** (`codebase-read-side.md:135`). Buildable from day one against the `products` entity; the ML entities light up as later milestones register jobs. |

**Dependency edges — each named with the artifact that forces it**

| edge | forcing artifact |
|---|---|
| M-03 → M-01 | `order_shipments` ingest issues 3-4 parallel GETs per order (`shipping_reader.go:149-180`, no multiget — fact #2) against an adapter where 429 is terminal on first contact (`capability_adapter.go:654-655`). |
| M-04 → M-01 | Backfill hydration is the multiget consumer; today hydration is N+1 (`capability_adapter.go:246-255, 292-301`) and `Ingestion.Pull` accumulates the whole catalog in memory (`ingestion.go:34-76`, 10 000-page cap). |
| M-03 → M-02 | `order_shipments` DDL + the `Divergence`/`FeeLayer` ports M-03's ingester writes through. |
| M-04 → M-02 | Same DDL/ports; plus the scheduler `incremental=false` fix (`scheduler.go:160`) that M-04's daily job needs to report truthfully. |
| M-05 → M-04 | Layer-2 requery needs the ingested `category_id` + current price (fact #4: `listing_prices` needs `logistic_type`/`shipping_mode`/`billable_weight`); stock divergence needs `available_quantity`, which M-04 introduces. |
| M-06 → M-03 | M-03 owns `IngestOrder(providerOrderID)` (A-04); M-06 extends it rather than forking a second writer — the exact defect `import_service.go` + `enrich_service.go` represent today. |
| M-06 → M-02 | `channel_fees` layer-3 rows + `divergences` rows for the 3→2 audit. |
| M-07 → M-02 | The fee-layer resolution port; without it the pricing resolver has nothing to read after degrau-3 is unwired at `root.go:845-851`. |
| M-07 ⤳ M-05 | **Data-quality edge, not a compile edge.** M-07 is correct without M-05 (falls to `config` provenance); it is *useful* after it. The hub can schedule M-07 early if the lane is free — the honest fallback is a ratified outcome, not a regression. |
| M-08 → M-06 | The worker's sole action is calling the ingester (A-04, design §5 "duas portas do mesmo caminho"). |
| M-09 → (none, compile) ⤳ M-04/M-06 | `sync_state` 0075 and `GET /sync/runs` already exist; entity rows for `listings`/`orders` only appear once those jobs register (`root.go:672-677` currently registers `products` alone). |
| **All → A-14** | `root.go` and the OpenAPI/SDK pair are hub-serialized; this is a *scheduling* edge that overrides code-level disjointness for any two milestones both changing FE-facing contracts. |

**Suggested lanes.** Wave A: M-01 ∥ M-02 ∥ M-09. Wave B: M-03 ∥ M-04. Wave C: M-05 ∥ M-06 ∥ M-07. Wave D: M-08. Three parallel workers sustained for most of the mission, versus a strictly serial chain under a wave-per-milestone cut.

---

## Top risks

**R-A — Premature read-decoupling blanks ratified /pedidos columns for multiple milestones.** Highest-likelihood plan-induced defect. Design §6 places orders-decoupling first as a "chip pequeno"; the four sites are the *sole* producers of `SLA`, `bucket`, `frete_real`, `destinatario`, `destino_uf/cep`, `rastreio` and `comprador_fiscal` (`http_handler.go:549-567`, `:462-467`; column set `PedidosTable.tsx:95-146`). *Mitigation:* A-05 — replacement-before-deletion, enforced by the shrinking allowlist guard, with a per-milestone before/after screenshot pair on the same field. *Trigger:* any milestone diff that removes an ML read call without adding a persisted source in the same merge.

**R-B — Resumable backfill lands on top of MASS-CLOSURE and wipes the catalog.** `ApplyCompletedPull` closes everything unconditionally (`repository.go:390-394`); design mandates a resumable cursor. A partial run — and A-01 makes partial runs *more* frequent by retrying rather than failing fast — closes the whole catalog, and the badge/divergence work downstream then computes against a dead catalog. *Mitigation:* A-06 sequenced strictly before the cursor inside M-04, proven by an abort-after-page-1 must-fail asserting zero `status='closed'` flips, on a >1-page fixture (R-3). *Trigger:* the cursor commit landing before the closure-semantics commit.

**R-C — Parallel-lane merge collisions in seams no milestone owns.** Three concrete instances, each invisible until merge: migration numbers (next free 0086, `0021` already duplicated); `root.go` (touched by all nine); and the **hand-written** SDK client object literal (`index.ts:2113-2330+`) paired with the OpenAPI same-commit rule. The SDK being hand-written rather than generated is the non-obvious one and inverts the usual parallelism intuition. *Mitigation:* A-12 (hub pre-allocates disjoint ranges in the brief), A-14 (module-local composition constructors + hub-held contract lock, at most one FE-facing contract change in flight). *Trigger:* two dispatched briefs without disjoint migration ranges, or two in-flight milestones both listing OpenAPI in their file scope.

**R-D — PII escapes via the new raw-persistence rule and the new public endpoint.** The mission simultaneously introduces "persist `Raw json.RawMessage`" and an unauthenticated `POST /webhooks/{provider}`, while `billing_info` carries buyer document and fiscal address (`buyer_fiscal_reader.go:59-94`) and an unresolved PII-scrub pendency already exists on `docs/design/evidence/ml-api/` (`external-ml-api-facts.md:7`). A worker implementing design §6 literally persists fiscal payloads. *Mitigation:* A-03 (raw persisted selectively, `billing_info` excluded, bounded with an explicit truncation marker) + A-11 (webhook body never reaches a domain table); a test asserting no persisted column contains the buyer document, driven from a fixture that contains one. *Trigger:* any migration adding a `raw jsonb` column to a buyer/fiscal-bearing table.

**R-E — Divergence becomes two mutually-incompatible semantics in one table, and false-positives make the badge noise.** Two producers on different cadences (listings daily, orders every 5 min) will pick different write models — append-events versus one-open-row — both defensible in isolation; the "divergentes" filter and badge count then mean nothing, and the 5-minute producer grows unbounded rows. Layered on top is mission R-5: an ERP sellable cut and an ML `available_quantity` observed hours apart look like a divergence. *Mitigation:* A-10 — one shape decided in M-02 before either producer exists, upsert-keyed, with **both** observation timestamps NOT NULL so staleness is distinguishable from disagreement; two-direction proof (create → badge; converge → `resolved_at` + badge gone) in the same drive. *Trigger:* the first badge appearing on a listing synced minutes earlier.

*Sharpened, lower-tier:* the `nil`-cursor-erases trap (`scheduler.go:42-45`) producing a silent infinite re-backfill (A-07); the hydration restructure starving `AbsorbProviderSnapshots` and degrading re-vínculo (A-13); and a milestone brief inheriting design §4's obsolete `fee_sync.go:29` claim and shipping an unfalsifiable "the seed dies" criterion (A-09).

---

## Checked paths

Read in full, in the mandated order:

1. `C:\Users\leandro.theodoro\Documents\marketplace-central\docs\design\MIS-007-ML-SYNC-DESIGN.md`
2. `C:\Users\leandro.theodoro\Documents\marketplace-central\.mnfs\MIS-007-ml-sync\mission.md`
3. `C:\Users\leandro.theodoro\Documents\marketplace-central\.mnfs\MIS-007-ml-sync\research\codebase-ingest-side.md`
4. `C:\Users\leandro.theodoro\Documents\marketplace-central\.mnfs\MIS-007-ml-sync\research\codebase-read-side.md`
5. `C:\Users\leandro.theodoro\Documents\marketplace-central\.mnfs\MIS-007-ml-sync\research\external-ml-api-facts.md`

Read for manifest verification only:

6. `C:\Users\leandro.theodoro\Documents\marketplace-central\.mnfs\MIS-007-ml-sync\planning-reviews\p3-input-r01.sha256`

No other file was opened. No repository browsing beyond the above; every `path:line` cited in this document is quoted from the research notes, not from independent inspection. Nothing was written or mutated; the only command executed was a read-only `sha256sum` over the six paths above.
