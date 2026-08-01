# P7 r01 — Claude crew seat 1/5 verbatim (★1 Completeness + ★5 Traceability)

```yaml
round: r01
seat: 1/5
scope: "★1 Completeness + ★5 Traceability"
model: opus (cold, read-only)
manifest: planning-reviews/p7-input-r01.sha256
manifest_digest: 3429ceb0a3c91e09c2eef18af526dd14bb44a46ddc9c94303bbc146238cd9f3e
task_id: ad7d57ab2fc9eec9e
transport: task-notification result (in-context), persisted same turn
verdict: "★1 FAIL (mission.md:28 'comissão/frete camada 2' — frete sem produtor/critério/Non-Scope); ★5 PASS; 7 advisories A1–A7"
persisted: 2026-08-01
```

Verbatim result below (unedited):

---

# P7 Readiness Review — MIS-007-ml-sync — round r01 — Claude crew seat 1/5

**Scope:** ★1 Completeness + ★5 Traceability only. No seven-★ verdict computed (dispatching session folds the crew).
**Mission root:** `C:\Users\leandro.theodoro\Documents\marketplace-central\.mnfs\MIS-007-ml-sync`
**Manifest digest line read verbatim** (`planning-reviews/p7-input-r01.sha256:66`):

```
# top-level digest (sha256 of the 65 sorted entries above): 3429ceb0a3c91e09c2eef18af526dd14bb44a46ddc9c94303bbc146238cd9f3e
```

Matches the expected digest in the dispatch. Honest limitation: I read the digest line and the 65 entry lines; I did **not** recompute per-file SHA-256 (no hashing performed — read-only text tools).

---

## ★1 Completeness — **FAIL**

**Procedure run:** enumerated every requirement/criteria ID in the parent brief (`mission.md`): Outcome bullets 1–5, Domain Scope items (2/3), (5), (6), (7), (8), (9), (10), Non-Scope, ADR-01..ADR-14, Q1/Q2/Q3/Q4/Q6 (+Q5/Q7 declined), R-1..R-7, and MIS07-C1..C9; then mapped each to ≥1 downstream artifact (milestone / feature / validation criterion). Then grepped the to-be-built scope for `TBD|TODO|FIXME|???|XXX`.

**Placeholder grep — clean.** Zero `TBD`, zero English-placeholder `TODO`. The four uppercase `TODO` hits are Portuguese `TODO/TODOS` = "every/all" (`research/divergences-interface-contract.md:37` "TODO ingest da entidade"; `M-09-sync-observability/milestone.md:66` "TODO campo sem observação"; `M-02.../F-04:25`; `M-08.../validation-contract.md:67`), not placeholders. `XXX` hit is `XXXtransport.NewHandler` in a repo-fact citation (`research/codebase-read-side.md:125`).

**Coverage — 24 of 25 enumerated requirements covered.** Spot-proof of the density: 25 features exist and 25 are listed in the nine `## Features` tables (1:1, no orphan directory, no unlisted feature); Q1→MIS07-C2+M03-C*, Q2→MIS07-C5+M08-C1+M02 PII criterion, Q3→MIS07-C3+M04-C1+M06-C2+M08-C3, Q4→MIS07-C8+M09-C1..C6, Q6→MIS07-C6, R-1→M-04 evidence+M-01 bucket, R-3→M04-C6+M05-C3, R-4→M06-C3+M06-U1, R-5→M05-C2, R-6→MIS07-C5+M02 PII+M-03 F-01 `billing_info` assert, Domain Scope (9)→M-09, (10)→M-08.

### Defect

**Defect locus:** `mission.md:28-29`
**Offending token:** `comissão/frete camada 2` (the `frete` half)

Verbatim excerpt, `mission.md:28-29`:

> - `listings` completa (E3 + estoque ML + comissão/frete camada 2) sincronizada por
>   backfill retomável + scheduler diário + refresh manual.

No downstream artifact produces listing-level freight at layer 2. The binding contract declares the opposite, verbatim `research/channel-fees-interface-contract.md:65`:

> - `api_shipping_options` no CHECK é RESERVA aditiva: NENHUM produtor nesta missão (frete de

and `research/channel-fees-interface-contract.md:87`:

> honest-absent do M-07 fica inalcançável. Frete: camada 2 → **honesto-desconhecido** (frete

The mission's own Milestone Strategy row agrees, `mission.md:253`:

> | M-05 `listings-fees-divergence` | Camada 2 (`listing_prices` com category_id/price ingeridos; frete = honesto-desconhecido, sem shipping_options nesta missão), …

The owning milestone's Required Outcome never mentions freight (`M-05-listings-fees-divergence/validation-contract.md:30-34` — camada 2 = `percentage_fee, fixed_fee, financing_add_on_fee, price_used, listing_type_id`, per M05-C1:46-48). Freight exists only at **layer 3, on orders** (`M-06-orders-backfill-decomposition/validation-contract.md:85-86` "frete camada 3 = custo seller do shipment"), which is a different entity from the `listings` promise at `mission.md:28`.

Crucially, the element was never relocated: `## Non-Scope` (`mission.md:64-73`) lists mercado/catalog offers/sellers/`ml_tariffs` sweep, de-para categoria, writes ML, onboarding saga, webhook topic `items`, and cadence relaxation — **freight-on-listings is absent from that list.** So the mission Outcome states a delivered capability that no milestone owns, no validation criterion asserts, and no Non-Scope line disclaims. A hub orchestrator reading only `mission.md` Outcome would dispatch expecting listing freight and be told at M-05 close that it is honest-unknown by design.

**Yes-if (minimal, grounded in approved scope — no new scope):** either (a) edit `mission.md:28` to read `E3 + estoque ML + comissão camada 2` and add one Non-Scope line naming freight-at-layer-2 as declined with its recorded reason and destination, verbatim-consistent with `mission.md:253` and `research/channel-fees-interface-contract.md:65,87` (`api_shipping_options` = reserva aditiva, sem produtor nesta missão); or (b) if listing freight is genuinely in scope, name the owning milestone and add a producer + validation criterion. (a) matches every other artifact in the tree and is a one-line documentation repair.

---

## ★5 Traceability — **PASS**

**Procedure run:** backward (every artifact → parent requirement), forward (every requirement → verification), decomposition (no orphan goals), plus assumption traceability against `Clarified Decisions` → `Accepted assumptions:`.

**Backward — no orphans.** All 9 milestones appear in the mission's Milestone Strategy table (`mission.md:249-257`). All 25 feature briefs are listed in their milestone `## Features` table (grep of `^| F-0` over `M-*/milestone.md` returns exactly 25 rows: M-01/2, M-02/4, M-03/3, M-04/4, M-05/3, M-06/3, M-07/2, M-08/2, M-09/2), and the manifest carries exactly those 25 `feature.md` paths — no unlisted artifact, no dangling table row. Every feature brief names both parents, e.g. `M-08-webhook-ingest/F-02-worker-callback/feature.md:15-21`:

> ## Mission
>
> MIS-007 ml-sync.
>
> ## Milestone
>
> M-08 webhook-ingest.

**Forward — every requirement reaches a verification.** Each milestone carries `M0X-C*` plus operator-mandated `M0X-U*` browser-driven criteria, and the mission contract closes the loop at `validation-contract.md:188-196` (MIS07-C9 = "leitura dos 9 `M-*/validation-result.md`" → "9/9 verdict PASS por QA Validator"). Cross-milestone seams are traced to a criterion on both sides, e.g. the `WebhookStatsReader` port is published by M-09 (`M-09-sync-observability/validation-contract.md:76-88`, M09-C3 "default → bloco webhook BYTE-IGUAL ao canônico IC-05 … fake injetado via `WithWebhookStatsReader`") and its real implementation is discharged by M-08 (`M-08-webhook-ingest/F-02-worker-callback/feature.md:78-83`), with the dependency edge declared in the mission edge table (`mission.md:291` "| M-08→M-09 | porta WebhookStatsReader (IC-05) + injeção `WithWebhookStatsReader` da região ancorada do M-08 |").

**Decomposition — no orphan goals.** Every Domain Scope item (`mission.md:52-62`) lands on a named milestone; every ADR-01..14 is claimed by ≥1 milestone ownership cell (`mission.md:297-307`).

**Assumption traceability — recorded, not orphan.** `mission.md:101-112`:

> - Accepted assumptions:
>   - `channel_fees` nasce com schema das 3 camadas; esta missão popula só camadas 2 e 3
>     (camada 1 = MIS-008). Reversível: coluna extra sem consumidor não força retrabalho.
>   - Worker do webhook roda in-process (goroutine no server_core); …
>   - Divergência auto-resolve no ingest: …
>   - Installation/tenant única existente (conta ML do operador). …
> - Owner decisions still open: None — gate P1 respondido 2026-07-31.

I red-teamed for invented cross-worker decisions adopted without an operator answer and absent from the record. Every candidate I found is written down with provenance: the dormant camada-1 arm (`M-07-pricing-fee-read/F-01-fee-read-resolver/feature.md`, "ninguém ESCREVE camada 1 … documentado como adormecido") traces to the first Accepted assumption; the three spine amendments carry their round of origin (`mission.md:182` "EMENDA (auditoria P5 r02 N-2)", `:203` "EMENDA P5 r02 N-4", `:227` "EMENDA (auditoria P5 r03 P-2)"), and the r02 amendment explicitly reverses the P3-ratified `baseline_commission_percent` disposition recorded at `planning-reviews/p3-reconciliation-r01.md:17`, naming the falsified premise. No orphan assumption found. (The seven IC-level "decisões finas" recorded outside the `Accepted assumptions:` heading are logged as advisory A2 below — recorded ⇒ not orphan ⇒ does not fail this criterion.)

---

## Advisory (should-meet) — never flips PASS/FAIL

| # | Locus | Finding |
|---|-------|---------|
| A1 | `mission.md:369` | R-2 mitigation cites `missed_feeds 2 dias`, but `mission.md:239` declares "missed_feeds não-consumido" e `research/webhook-inbox-interface-contract.md:96` reads "missed_feeds (`GET /missed_feeds`): NÃO consumido nesta missão". One mitigation leg is fictional; the other two (scheduler 5min, /integracoes última notificação) are real and discharged by M08-C5/M08-U3 and M-09. **Belongs to ★2 Consistency — flagged cross-scope for the ★2 seat.** |
| A2 | `mission.md:235-240` | Seven binding cross-worker "decisões finas" (fee only in `channel_fees`; `/sync/health` new instead of reusing `/sync/runs`; DDL placement M-02/M-04/M-08; webhook in OpenAPI without SDK method; missed_feeds unconsumed; R$0.01 tolerance; pinned composite id formats) are recorded under `## Architecture Spine`, not under `Clarified Decisions` → `Accepted assumptions:` where ★5's procedure names them. Recorded and traceable, so not orphan — but a reader following the rubric's heading will not find them. |
| A3 | `mission.md:343` | Q1's Validation criterion column says "critério de milestone: read medido <2s com dados reais". Only M-03 carries a milestone-level <2s (`M-03.../validation-contract.md:128` orders detail, `:139` drawer). `/anuncios` and the `/pedidos` list have no milestone <2s criterion — only mission-level MIS07-C2 (`validation-contract.md:60-73`). Coverage holds via MIS07-C2; the Q1 row's phrasing overstates where the criterion lives. |
| A4 | `mission.md:379-394` | Handoff cites required artifacts absent from the frozen P7 manifest: `p5-claude-decomposition-audit-r08.md`, `p5-input-r09.sha256`, `p3-claude-candidate-r01.md`, `p3-opus-counterproposal-r01.md`. P7 reviewers cannot verify those verdict claims from the frozen 65-entry input. (The audit/candidate files do exist in the tree — confirmed incidentally via grep hits — but are outside the manifest.) |
| A5 | `M-09-sync-observability/F-02-integracoes-health-section/feature.md:65` | "Refetch: polling leve (30s) OU manual — spec decide contra idioma web-query da página" — a live choice deferred to spec rather than decided. Bounded ("sem websocket"), single-owner, so not a completeness gap; ★7 Unambiguity territory. |
| A6 | `M-09-sync-observability/validation-contract.md:148-149` | M09-C6 is explicitly deferred: "critério DIFERIDO — não bloqueia close do M-09; bloqueia close da MISSÃO via MIS07-C8". Traceable and honest, but M-09 can close with a Required criterion unverified; MIS07-C8 evidence should name M09-C6 explicitly so the deferral cannot be lost. |
| A7 | `mission.md:374` | R-7 (merge collision on ownerless seams) is the only risk with no validation criterion anywhere in the 10 contracts — mitigated purely by hub process (ADR-12 ranges, ADR-14 serialization). Weaker forward-trace than R-1..R-6, which all reach a `M0X-C*`. |

---

## Manifested paths actually read

**Read in full (58 of 65):**
`mission.md` · `validation-contract.md` · `research/channel-fees-interface-contract.md` · `research/divergences-interface-contract.md` · `research/orders-persistence-interface-contract.md` · `research/webhook-inbox-interface-contract.md` · `research/sync-health-interface-contract.md` · `research/sync-ingest-ports-interface-contract.md` · `research/listings-sync-interface-contract.md` · `research/p5-prerequisites.md` · `research/external-ml-api-facts.md` · `research/codebase-ingest-side.md` · `research/codebase-read-side.md` · `planning-reviews/p3-reconciliation-r01.md` · `planning-reviews/p5-reconciliation-r08.md` · `M-01-ml-client-hardening/{milestone.md, validation-contract.md, F-01-resilience-decorator/feature.md, F-02-items-multiget-raw-dto/feature.md}` · `M-02-sync-core-seam/{milestone.md, validation-contract.md, F-01-core-ddl, F-02-fee-divergence-ports, F-03-scheduler-incremental-cursor, F-04-read-guard-allowlist}/feature.md` · `M-03-orders-shipment-persist/{milestone.md, validation-contract.md, F-01-ml-ingest-readers, F-02-ingest-order-v1, F-03-read-path-switch}/feature.md` · `M-04-listings-backfill-ingest/{milestone.md, validation-contract.md, F-01-listings-ddl, F-02-mass-closure-replacement, F-03-backfill-cursor-ingest, F-04-scheduler-refresh-wiring}/feature.md` · `M-05-listings-fees-divergence/{milestone.md, validation-contract.md, F-01-camada2-fee-ingest, F-02-stock-divergence, F-03-anuncios-fe-contract}/feature.md` · `M-06-orders-backfill-decomposition/{milestone.md, validation-contract.md, F-01-backfill-incremental, F-02-decomposition-camada3, F-03-audit-fe-pedidos}/feature.md` · `M-07-pricing-fee-read/{milestone.md, validation-contract.md, F-01-fee-read-resolver, F-02-precos-provenance-fe}/feature.md` · `M-08-webhook-ingest/{milestone.md, validation-contract.md, F-01-inbox-endpoint, F-02-worker-callback}/feature.md` · `M-09-sync-observability/{milestone.md, validation-contract.md, F-01-sync-health-endpoint, F-02-integracoes-health-section}/feature.md` · `planning-reviews/p7-input-r01.sha256`

**Unread-path sweep — 7 manifested files opened only by targeted grep, not end-to-end:**
`planning-reviews/p5-reconciliation-r01.md` … `p5-reconciliation-r07.md`. Justification: for ★1/★5 these are historical fold records of superseded artifact states, not the current requirement set. I queried each for `verdict_received|disposition:|REJECT|not folded|deferred` and every one returns the same disposition shape — e.g. `p5-reconciliation-r01.md:12` "ALL findings ACCEPTED as valid; fold applied in full", `p5-reconciliation-r07.md:16` "ALL findings ACCEPTED as valid; fold applied in full (1 blocking + 5 …" — with zero rejected and zero unfolded findings. `p5-reconciliation-r08.md:13-14,84` closes the loop ("PART A: F-r07-1..6 + observation ALL CONFIRMED CLOSED; reopened-closure sweep r01–r06 CLEAN, zero reopened"; "No r01–r07 closure reopened"). No open finding in r01–r07 could therefore leave an uncovered requirement or an orphan in the current tree. Repo files: none read this round — no planning artifact claim required repo verification for ★1/★5 beyond what the artifacts already cite with file:line.
