# P7 r01 — Claude crew seat 2/5 verbatim (★2 Consistency + ★3 Seam Ownership)

```yaml
round: r01
seat: 2/5
scope: "★2 Consistency + ★3 Seam Ownership"
model: opus (cold, read-only)
manifest: planning-reviews/p7-input-r01.sha256
manifest_digest: 3429ceb0a3c91e09c2eef18af526dd14bb44a46ddc9c94303bbc146238cd9f3e
task_id: a860342ef076c8dbd
transport: task-notification result (in-context), persisted same turn
verdict: "★2 FAIL (divergences[] DTO fork + citação §DTO inexistente); ★3 FAIL (write-set M-05 fora do grant: listings/composition + listings/ports); 9 advisories A-1..A-9"
persisted: 2026-08-01
```

Verbatim result below (unedited):

---

# P7 Mission-Readiness Review — MIS-007-ml-sync — Claude crew seat 2 of 5

**Seat:** 2 of 5 (cold, independent) · **Round:** r01
**Criterion scope:** ★2 Consistency + ★3 Seam Ownership ONLY. No seven-★ verdict computed.
**Mission root:** `C:\Users\leandro.theodoro\Documents\marketplace-central\.mnfs\MIS-007-ml-sync`
**Manifest digest line read verbatim** (`planning-reviews/p7-input-r01.sha256:66`):

```
# top-level digest (sha256 of the 65 sorted entries above): 3429ceb0a3c91e09c2eef18af526dd14bb44a46ddc9c94303bbc146238cd9f3e
```

Matches the expected digest supplied in the dispatch.

---

## ★2 Consistency — **FAIL**

**Passing sweeps** (recorded so the fold can see the procedure ran, not just the failure): migration ranges disjoint and gap-legal (M-02 `0086-0089`, M-04 `0090-0092`, M-08 `0093`, M-06 `0094-0095` reserve; repo tip is 0085 — verified `research/codebase-ingest-side.md:47`); route prefixes coherent (`POST /webhooks/{provider}`, `GET /sync/health`, `/listings`+`filter.divergentes`, `/orders`, `/pricing`); port names identical across IC and features (`ChannelFeeWriter`/`ChannelFeeReader`, `DivergenceRecorder`/`DivergenceReader`, `OrderIngestor`/`ListingIngestor`, `WebhookStatsReader`+`WithWebhookStatsReader`); enum vocabularies agree (`phase` backfill|incremental|sweep, inbox status, `kind` estoque|tarifa, `entity_type`, `subject_type`, `fee_kind`, `value_type`, `layer`, `origem`); every list operation declares an order (IC-02 `detected_at` desc, IC-05 entities name asc, IC-04 `received_at` asc, IC-07/IC-03 "ordenação existente preservada"); Error Matrices cover the feature-returned cases for IC-04, IC-05, IC-07.

**The divergence — DTO shape of `divergences[]`.**

Parent contract, `research/divergences-interface-contract.md:63` (verbatim):

> DTOs de listings/orders ganham `divergences: [{kind, detected_at}]` (array vazio = sem divergência aberta) — flag persistida, zero cálculo no read.

Feature, `M-05-listings-fees-divergence/F-03-anuncios-fe-contract/feature.md:30` (verbatim):

> `divergences[]` (IC-02 §DTO — kind/expected/observed/timestamps das rows abertas)

and `:38`:

> ListingDetailPanel mostra os dois lados + timestamps da divergência.

and `:70`:

> DTO canonical: IC-02 §DTO + IC-01 §Canonical Examples camada 2

- **Rubric criterion:** ★2 — "grep every enum value, error code, data shape/field name … across ALL files against the interface contract"; a feature-returned shape that the parent contract does not declare is a divergence.
- **Defect locus:** `M-05-listings-fees-divergence/F-03-anuncios-fe-contract/feature.md:30` (repeated `:70`; consumed at `:38`).
- **Offending tokens:** `IC-02 §DTO` and `kind/expected/observed/timestamps das rows abertas`.
- **Why load-bearing, not cosmetic:** F-03 `:79-81` pins *"dados de tarifa/divergência vêm na MESMA query de lista (sem query paralela); panel usa o item da lista + fetch de detalhe existente"*. Under IC-02's two-field array the detail panel at `:38` has no source for "os dois lados + timestamps" — the requirement is unsatisfiable against the parent contract. M05-C5 (`M-05/validation-contract.md:123`) then commits *"OpenAPI (DTO tarifa + divergences[] + param) e `sdk-runtime` no MESMO commit"*, so the unresolved shape is what gets frozen into the hand-written SDK under the ADR-14 contract lock.
- Secondary: **IC-02 has no section named `§DTO`.** Its sections are Boundary / … / Required Inputs / Required Outputs / Enums And Statuses / Error Cases / Error Matrix / Persistence Expectations / Canonical Examples / Database Shape / Seed Data / Timestamp And ID Semantics / Compatibility Rules. The pinned shape lives in **§Required Outputs**. The citation points at nothing.

**Yes-if (either arm closes it, both inside approved scope):**
1. Amend `research/divergences-interface-contract.md` §Required Outputs to declare the open-row projection the panel needs, using its own already-ratified column vocabulary (`kind`, `expected_value`, `observed_value`, `expected_source`, `observed_source`, `expected_observed_at`, `observed_observed_at`, `detected_at` — all NOT NULL at `:51-59`), **or**
2. Narrow `F-03:30,38,70` to `{kind, detected_at}` and name the declared surface the two sides come from.

In both arms, replace the `IC-02 §DTO` citation with `IC-02 §Required Outputs` at `F-03:30` and `F-03:70`.

---

## ★3 Seam Ownership — **FAIL**

**Passing sweeps:** every declared-parallel set is disjoint on the six matrix axes. Lane A (M-01∥M-02∥M-09): the one shared package `sync/application/` carries a named additive lock at `mission.md:324-327` (M-02 owns `scheduler.go`, M-09 additive-only `sync/application/health_*` + `sync/transport/**`) and M-09/F-01:94-97 matches it exactly. Lane B (M-03∥M-04): `orders/**` vs `listings/**` disjoint; M-03's new ML readers inside the M-01 adapter dir are explicitly permitted by `mission.md:299`; both root.go regions are hub-serialized. Lane C (M-05∥M-06∥M-07): shared `channel_fees`/`divergences` writes partitioned by named seam locks (layer 2 + kind=estoque vs layer 3 + kind=tarifa, `M-05/milestone.md:59-61`); FE routes /anuncios vs /pedidos vs /precos disjoint; contract commits serialized ≤1-in-flight by ADR-14. Every migration-adding milestone carries an explicit pre-allocated block. M-08→M-09 seam owned by IC-05 with the injection-by-reference rule pinned. Repo anchor spot-checked and true: `apps/server_core/internal/modules/sync/composition/scheduler.go:11` → `const InstallationScopeERP = "erp"`.

**The hole — M-05's declared write surface is outside its registered grant.**

Mission grant, `mission.md:319-323` (verbatim, exhaustive enumeration):

> **Posse de ingest de listings**: M-05 estende superfícies do M-04 PÓS-close (lock aditivo registrado: `listings/application/**`, `listings/transport/**` e `listings/adapters/postgres/repository.go`, additive-only — o write-set REAL do M-05 F-03 inclui transport e o repository de leitura, não só application; grant alinhado ao write-DAG na auditoria P5 r05 F-r05-3).

Matrix cell repeats the same three paths, `mission.md:303`: *"(listings application/transport/`adapters/postgres/repository.go` via lock aditivo do M-04 PÓS-close — F-r05-3)"*.

Features declare writes outside all three. `M-05-listings-fees-divergence/F-01-camada2-fee-ingest/feature.md:55-58` (verbatim):

> A fiação (construção do writer e injeção no pipeline) vive DENTRO do package de composition de listings, sob o lock aditivo registrado do M-04 — `root.go` NÃO é tocado pelo M-05 (matriz da missão: célula root.go = `—`; auditoria P5 r03 P-5).

Identical claim at `M-05-listings-fees-divergence/F-02-stock-divergence/feature.md:50-52`, plus a second unlisted path at `F-02:75`:

> Owned paths: avaliador novo + porta de leitura mirror em listings; fiação pós-sweep.

- **Rubric criterion:** ★3 — enumerate cross-worker seams; each must be owned by an IC or ADR; every predicted shared-seam addition must be NAMED as a seam lock.
- **Defect loci:** `M-05/F-01/feature.md:55-58`; `M-05/F-02/feature.md:50-52` and `:75`; the grant that fails to cover them, `mission.md:319-323` and the M-05 cell `mission.md:303`.
- **Offending tokens:** `package de composition de listings` and `porta de leitura mirror em listings` — neither `listings/composition/**` nor `listings/ports/**` appears in the registered lock.
- **Why the claim of cover is false, verified against the repo (factual-claim check only):** composition is a **sibling** package, not a subtree of `application/` — `apps/server_core/internal/modules/sync/composition/` and `.../product_links/composition/` exist, and `apps/server_core/internal/modules/listings/` today contains only `adapters/`, `application/`, `domain/`, `ports/`, `transport/`. The listings composition package does not exist yet; **M-04 creates and owns it** (`M-04/F-04/feature.md:63`: *"Owned paths: composição do módulo listings (arquivo novo), linha ancorada no root.go"*; matrix `mission.md:302` gives M-04 all of `listings/**`). So M-05 writes into an M-04-owned path with no grant naming it.
- **Provenance of the hole, so the planner doesn't re-close it the same way:** `planning-reviews/p5-reconciliation-r03.md:73-78` asserted the cover — *"the composition package is already the additive-lock surface"* — but never added the path to the grant; the grant was later rewritten at r05 (F-r05-3) to match the real write-set and **still** enumerates only the three paths. An assertion in a reconciliation is not a grant line.
- **Concurrency axis is clean:** M-04 is closed before lane C opens, so this is an ownership-grant defect, not a live collision. It still fails ★3 because the seam addition is not NAMED, and three artifacts state the lock three non-identical ways (see advisory A-1).

**Yes-if (either arm, inside approved scope):**
1. Extend the registered lock at `mission.md:319-323` **and** the M-05 cell at `mission.md:303` to `listings/composition/**` and `listings/ports/**` (additive-only), and restate `M-05/milestone.md:53-57` with the identical enumeration; **or**
2. Relocate the M-05 wiring and the mirror-read port into `listings/application/**`, which the grant already covers, and drop the phrase *"package de composition de listings"* from `F-01:55-58` and `F-02:50-52`.

---

## Advisory findings (never flip a verdict)

- **A-1** — `M-05/milestone.md:53-55` states the same lock a third way: *"arquivos de application/adapters de listings, additive-only"* — drops `transport` (which F-03:91 relies on) and widens `adapters` beyond `repository.go`. Three renditions of one lock across `mission.md:303`, `mission.md:319-323`, and the milestone.
- **A-2** — `M-05/F-03/feature.md:91-94` owned paths omit `AnunciosPage.tsx`, though the `filter.divergentes` control and the `listListings` call site live there (`codebase-read-side.md:86` → `AnunciosPage.tsx:163`). No rival owner (the M-05 cell grants /anuncios), so no collision.
- **A-3** — `IC-07:65` keys `listing_variations` `(tenant_id, provider, provider_listing_id, variation_id)` while the parent `listings` PK is `(tenant_id, installation_id, provider_listing_id, variation_id)`, which `IC-07:36` pins unchanged — no FK expressible for a table `IC-07:110-117` requires written in the same transaction. Same class in `IC-02` (`entity_id` carries no installation) and `IC-03` (`order_shipments`). Practically collision-free (MLB ids are account-scoped), hence advisory.
- **A-4** — `M-03/F-02/feature.md:25` says *"assinatura IC-06 verbatim"* while `:68` declares a tenant-less signature. IC-06 explicitly permits tenant-by-repo-scoping, so this is wording, not divergence.
- **A-5** — `M-06/F-03/feature.md:73` cites *"IC-03 §DTO"*; IC-03 has no such section (it is §Required Outputs). Same broken-section idiom as the ★2 locus — worth fixing in the same pass.
- **A-6** — `research/channel-fees-interface-contract.md` Error Matrix carries only the placeholder row `| — | — | — | sem endpoint próprio; …` while M-07 F-01 introduces a typed defaults/config failure surfaced through `/pricing`.
- **A-7** — Migration ordering: M-08 owns `0093` (lane D) and M-06 reserves `0094-0095` (lane C). `mission.md:328` permits the reserve to stay empty but does not address the filename-ordered runner (`migrate/runner.go` applies missing files in filename order) if the reserve is used before 0093 merges.
- **A-8** — `M-05/F-01:55-58` keeps the M-05 root.go cell at `—`; whether that holds depends on whether M-04's composition constructor can build the `ChannelFeeWriter` from the pool it already receives (idiom `synccomposition.NewProductsScheduler(pool, …)`) or must grow an argument at the root.go call site. Unprovable either way from the artifacts.
- **A-9** — Positive check, no action: the repo anchor `sync/composition/scheduler.go:11` cited at `M-09/F-01:38` is exact.

---

## Manifested paths actually read

All paths below are relative to `C:\Users\leandro.theodoro\Documents\marketplace-central\.mnfs\MIS-007-ml-sync\`.

**Read in full (56 of 65):**
`planning-reviews/p7-input-r01.sha256` (manifest itself) · `mission.md` · `validation-contract.md`
`research/channel-fees-interface-contract.md` · `research/divergences-interface-contract.md` · `research/orders-persistence-interface-contract.md` · `research/webhook-inbox-interface-contract.md` · `research/sync-health-interface-contract.md` · `research/sync-ingest-ports-interface-contract.md` · `research/listings-sync-interface-contract.md` · `research/external-ml-api-facts.md` · `research/codebase-ingest-side.md` · `research/codebase-read-side.md` · `research/p5-prerequisites.md`
`M-01…M-09/milestone.md` (9) · `M-01…M-09/validation-contract.md` (9) · all 22 `M-*/F-*/feature.md`

**Read partially:** `planning-reviews/p5-reconciliation-r03.md` (`:60-97`, for the P-5 fold that asserts the composition cover).

**Not read line-by-line — swept by targeted grep only** (`composition|lock aditivo|additive lock`, `divergences|§DTO`, `observed_at`): `planning-reviews/p3-reconciliation-r01.md`, `planning-reviews/p5-reconciliation-r01.md`, `-r02.md`, `-r04.md`, `-r05.md`, `-r06.md`, `-r07.md`, `-r08.md`. Neither FAIL above is closed by any hit in those files; disclosed so the fold can weigh it.

**Repo files opened only to verify a factual claim an artifact makes about the repo** (never to invent requirements), absolute paths:
- `C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\internal\modules\sync\composition\scheduler.go` (line 11)
- glob `C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\internal\modules\*\composition\*.go`
- glob `C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\internal\modules\listings\**\*.go`
