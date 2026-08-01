# P7 r01 — Claude crew seat 5/5 verbatim (independent adversarial double-pass ★2+★7)

```yaml
round: r01
seat: 5/5
scope: "★2 Consistency + ★7 Security Posture (independent adversarial double-pass)"
model: opus (cold, read-only)
manifest: planning-reviews/p7-input-r01.sha256
manifest_digest: 3429ceb0a3c91e09c2eef18af526dd14bb44a46ddc9c94303bbc146238cd9f3e
task_id: a4591b9f4bdf708e4
transport: task-notification result (in-context), persisted same turn
verdict: "★2 FAIL (2.1 divergences[] DTO fork IC-02 vs M-05/F-03; 2.2 ADR-09 total sem enforcement em /pedidos); ★7 PASS (neste assento; fold une com FAIL do seat 4); 5 advisories"
persisted: 2026-08-01
```

Verbatim result below (unedited):

---

# P7 Readiness Review — MIS-007-ml-sync — Claude crew seat 5/5 — round r01

**Seat**: Claude crew 5 of 5 — INDEPENDENT ADVERSARIAL DOUBLE-PASS
**Scope**: ★2 Consistency + ★7 Security Posture ONLY. Seven-★ verdict NOT computed (per dispatch).
**Manifest**: `C:\Users\leandro.theodoro\Documents\marketplace-central\.mnfs\MIS-007-ml-sync\planning-reviews\p7-input-r01.sha256` — 65 entries; digest line read verbatim at `:66`:
`# top-level digest (sha256 of the 65 sorted entries above): 3429ceb0a3c91e09c2eef18af526dd14bb44a46ddc9c94303bbc146238cd9f3e` — **matches** the dispatch-provided expected digest.

---

## ★2 Consistency — **FAIL**

Two cross-file shape/vocabulary divergences survive a full grep of enums, ports, routes, migration blocks and line anchors. Everything else on the hunt list checked clean (see Verified-clean below).

### Defect 2.1 — `divergences[]` DTO field set contradicts IC-02

**Cited excerpt (contract):** `research/divergences-interface-contract.md:63`
> `DTOs de listings/orders ganham `divergences: [{kind, detected_at}]` (array vazio = sem divergência aberta) — flag persistida, zero cálculo no read.`

**Cited excerpt (consumer):** `M-05-listings-fees-divergence/F-03-anuncios-fe-contract/feature.md:30`
> ``divergences[]` (IC-02 §DTO — kind/expected/observed/timestamps das rows abertas); param`

**Defect locus:** `M-05-listings-fees-divergence/F-03-anuncios-fe-contract/feature.md:30`
**Offending token:** `kind/expected/observed/timestamps`, attributed to "IC-02 §DTO". IC-02 §DTO declares exactly two keys: `kind`, `detected_at`. `expected`, `observed` and plural `timestamps` are the *table* columns (IC-02 Database Shape), not the DTO. The feature also renders them: `:38` "ListingDetailPanel mostra os dois lados + timestamps da divergência." — two-sided display is impossible from the 2-key DTO. The third consumer, `M-06-orders-backfill-decomposition/F-03-audit-fe-pedidos/feature.md:35`, names the same `divergences[]` with no field set at all, and `M-05/validation-contract.md:123` gates the OpenAPI/SDK pair on "divergences[] + param" without pinning the shape — so no artifact resolves the fork. Two milestones in the SAME lane C ship the same DTO field under one name with two different field sets, and the serialized contract commit (ADR-14) will land whichever writer runs first.

**Yes-if:** amend `research/divergences-interface-contract.md:63` to declare the single DTO field set the FE actually needs (`{kind, expected, observed, detected_at}` — the only variant that satisfies the already-approved `:38` "os dois lados" requirement), and make `M-05-listings-fees-divergence/F-03/feature.md:30` cite that exact key list; OR strike "expected/observed/timestamps" from `:30` and delete the "dois lados" render at `:38`. No new scope either way — one of the two existing statements must go.

### Defect 2.2 — fee-provenance mandate stated as TOTAL, unenforced on /pedidos

**Cited excerpt (parent, total claim):** `mission.md:190-191` (ADR-09)
> `Todo fee em tela carrega (camada, origem, coletado_em); número sem proveniência = milestone reprovado.`

**Cited excerpt (IC-01 restating it):** `research/channel-fees-interface-contract.md:72-73`
> `proveniência SEMPRE junto do número; consumidor que exibe número sem proveniência reprova milestone (ADR-09).`

**Defect locus:** `research/orders-persistence-interface-contract.md:104-108` — the canonical `decomposicao` key list is `receita_bruta`, `comissao_total`, `frete_seller`, `custo_produto`, `custo_fonte`, `custo_congelado_em`, `liquido`, `margem_pct`, `computado_em`, `incompleto`.
**Offending token/value:** `comissao_total` and `frete_seller` — fee amounts with **no `origem` and no `coletado_em`** anywhere in the object, then rendered on screen by `M-06-orders-backfill-decomposition/F-03-audit-fe-pedidos/feature.md:36-39` (PedidoDrawer "Decomposição financeira": receita, comissão, frete, custo, líquido, margem). A grep of the whole `M-06-*` tree for `origem|proveni` returns only two hits, both in the ingest feature (`F-02-decomposition-camada3/feature.md:29` `origem=api_order`, `:31` `origem=api_shipment`) — zero hits in the FE feature and zero in `M-06/validation-contract.md`. The M-06 gate at `M-06/validation-contract.md:143-150` drives `assert text "Comissão"` with Expected "drawer decompõe receita/comissão/frete/custo/líquido" — a fee number on screen with no provenance assertion. The mission gate does not cover it either: `validation-contract.md:159-161` (MIS07-C7) Command is scoped to `/precos` and `/anuncios` only, while its Blocking failure at `:166` is unscoped — "valor de fee sem origem".

**Yes-if:** either (a) add `origem` + `coletado_em` to the `comissao_total`/`frete_seller` terms in the IC-03 `decomposicao` shape at `research/orders-persistence-interface-contract.md:104-108` and add the provenance assertion to the `M-06/validation-contract.md:143-150` drive (values already exist in the layer-3 ledger written by `M-06/F-02:29,31` — no new collection scope), or (b) scope ADR-09 at `mission.md:190` to the surfaces it actually gates (`/precos`, `/anuncios`) and scope MIS07-C7's blocking failure at `validation-contract.md:166` to match. A total claim that one approved consumer provably violates is the exact "prosa falsa em contrato" class.

### Verified clean (adversarial pass, no defect)
- **Migrations**: M-02 `0086-0089`, M-04 `0090-0092`, M-08 `0093`, M-06 `0094-0095` — disjoint, gapless, and matching the `mission.md` matrix. Repo max is `0085_erp_import_products_sellable_fields.sql`; `apps\server_core\internal\platform\migrate\runner.go:20-53` keys `schema_migrations` by filename and applies all pending, so M-08's `0093` landing after M-06's `0094` in lane order is safe (repo already carries numeric gaps and a duplicated `0021`).
- **Ports**: `IngestOrder`, `IngestListing`, `WebhookStatsReader`, `WithWebhookStatsReader`, `WithTariffResolver` — identical spelling and signature across IC-05/IC-06 and every milestone/feature that names them.
- **Routes**: `/sync/health`, `/webhooks/{provider}`, `/orders/import` — consistent; `/orders/import` confirmed in `root.go:259-272` `registerBatchRoutes`.
- **Inbox status enum**: `received|processing|done|malformed|unmapped|dropped` — the CHECK at `research/webhook-inbox-interface-contract.md:59-61` covers every terminal named by `M-08/validation-contract.md:69-71`; `topic_ignored` is used as a *reason* alongside `done`, not as a 7th status. No divergence.
- **Fee vocabulary**: `api_listing_prices | api_order | api_shipment | config` — identical in IC-01:66-67, MIS07-C7:161-162, M-05-VC, M-06 F-02.
- **List ordering declared**: IC-02 `detected_at desc` (`:38`); IC-05 entities "por nome asc". No undeclared collection ordering found.
- **Repo line anchors** spot-verified and ACCURATE: `root.go:590-592` (ordersCost/Shipment/BuyerFiscal readers), `root.go:844-852` (`tariffResolver` … `batchOrch.WithTariffResolver`), `enrich_service.go:194-198` (`EnrichOne`), `route_deadline.go:23-28` (15s/120s classes).
- **DAG/lane assignments**: A(M-01∥M-02∥M-09) → B(M-03∥M-04) → C(M-05∥M-06∥M-07) → D(M-08) agree between `mission.md` and every milestone Handoff (`M-02:182`, `M-05:152`, `M-06:174`, `M-08:170`, `M-09:182`).
- **Self-refuted hypothesis (reported for honesty, not a defect)**: IC-07:155 "classe batch p/ pull" appeared false because `/listings/refresh` is absent from `registerBatchRoutes`; `apps\server_core\internal\modules\listings\transport\http_handler.go:28` registers `RegisterRouteClass("/listings/refresh", httpx.BatchRouteClass)` per-module. Claim is TRUE.

---

## ★7 Security Posture — **PASS**

Rubric bar: Security present in `## Quality Attributes` with ≥1 mitigation AND ≥1 Security-typed validation criterion. Both hold, and each probed surface has an owner.

**Cited excerpt (mitigation):** `mission.md:344` (Q2)
> `webhook = hint não-confiável; dado real sempre refetch na API autenticada; dedupe no inbox; IP origem gravado + comparado à allowlist oficial (log-only)`

**Cited excerpt (Security-typed criteria):** `M-08-webhook-ingest/validation-contract.md:38-42` (`ID: M08-C1` / `Type: Security`), Expected at `:46-50`:
> `200 + row no inbox + ZERO escrita em tabela de domínio (o worker REFETCH na API autenticada e o recurso não existe → status terminal, sem dado)`

Independent probe of the five surfaces named in the dispatch:
- **Forgery** — covered semantically, not by signature: `research/webhook-inbox-interface-contract.md:109-111` makes authenticated refetch the rejection mechanism; gated by M08-C1 (Security) and mission `validation-contract.md:116-132` MIS07-C5 (Security).
- **Replay** — `research/webhook-inbox-interface-contract.md:88` `UNIQUE (provider, notification_id) WHERE notification_id IS NOT NULL`; `M-08/validation-contract.md:79-93` (M08-C3) proves the honest case: with `_id` → transport dedupe; **without** `_id` → extra inbox rows are declared harmless and the real proof is domain idempotence of `IngestOrder`. Correct placement of the guarantee.
- **PII / raw** — ADR-03 raw-selective/PII-never enforced by TWO Security-typed criteria: `M-02/validation-contract.md:145-158` (M02-C7, `Type: Security`, "nenhuma coluna `raw`/jsonb de billing_info em tabela com dado fiscal/comprador") and MIS07-C5's grep `0086*..0095*`. Webhook body is pointer-never-data (ADR-11), asserted at M08-C1.
- **Tenant scoping** — single tenant/installation is an explicitly recorded accepted assumption (`mission.md:101-111`), and IC-06 (`research/sync-ingest-ports-interface-contract.md:35-37`) still routes tenant through repo scoping rather than dropping it. Not a silent omission.
- **Live ML write authorization** — the one write to the ML account (callback registration) is gated by a dedicated Security criterion, `M-08/validation-contract.md:132-146` (M08-C6, `Type: Security`): "autorização EXPLÍCITA do operador, com data, ANTES do registro do callback no app ML; sem ela, M08-C5 fica `could-not-run` e o milestone NÃO passa por stub" — satisfies the AGENTS.md live-ML-writes rule without a stub escape hatch.

Public-surface abuse limits are addressed at the envelope level (always-200, 64KB cap, interactive 15s class — `M-08/validation-contract.md:58-71`). Gaps below are advisory, not ★7 failures.

---

## Advisory findings (do not change any verdict)

1. **No flood/rate-limit criterion on the first unauthenticated surface.** `research/webhook-inbox-interface-contract.md:17` — "público, NÃO autenticado — 1ª superfície assim no sistema". Body cap and deadline exist; sustained-volume behavior (inbox growth, worker saturation) has no criterion. `dropped_24h` in IC-05 gives partial visibility only.
2. **`{provider}` path slug is unconstrained.** `POST /webhooks/{provider}` with `UNIQUE (provider, notification_id)`: an arbitrary slug value forks the dedupe key. Pinning `provider` to the known set at the route is a one-line contract tightening.
3. **`source_ip` / raw body retention unaddressed.** `research/webhook-inbox-interface-contract.md:57` stores `source_ip inet` NOT NULL; no retention/purge statement anywhere in the tree.
4. **`processing` rows are invisible in health.** IC-05's webhook block is `{last_notification_at, pending, dropped_24h}`; a worker wedged mid-row (status `processing`) counts in neither `pending` nor `dropped_24h`. M08-C4 covers shutdown but not a stuck row.
5. **M09-C6 is a deferred criterion inside a milestone contract** (`M-09/validation-contract.md:141-156`) — self-declared as not blocking M-09's close, only the mission's. Clear as written; flagged only because a deferred criterion in a milestone contract is easy to lose at close time.

---

## Manifested paths actually read

`mission.md`; `validation-contract.md`; `planning-reviews/p7-input-r01.sha256` (the manifest itself); `research/channel-fees-interface-contract.md`; `research/divergences-interface-contract.md`; `research/orders-persistence-interface-contract.md`; `research/webhook-inbox-interface-contract.md`; `research/sync-health-interface-contract.md`; `research/sync-ingest-ports-interface-contract.md`; `research/listings-sync-interface-contract.md`; `research/external-ml-api-facts.md`; `research/codebase-read-side.md`; `M-01-ml-client-hardening/milestone.md`, `/F-01-resilience-decorator/feature.md`, `/F-02-items-multiget-raw-dto/feature.md`; `M-02-sync-core-seam/milestone.md`, `/F-01-core-ddl/feature.md`, `/F-02-fee-divergence-ports/feature.md`, `/F-03-scheduler-incremental-cursor/feature.md`, `/F-04-read-guard-allowlist/feature.md`, `/validation-contract.md`; `M-03-orders-shipment-persist/milestone.md`, `/F-01-ml-ingest-readers/feature.md`, `/F-02-ingest-order-v1/feature.md`, `/F-03-read-path-switch/feature.md`; `M-04-listings-backfill-ingest/milestone.md`, `/F-01-listings-ddl/feature.md`, `/F-02-mass-closure-replacement/feature.md`, `/F-03-backfill-cursor-ingest/feature.md`, `/F-04-scheduler-refresh-wiring/feature.md`; `M-05-listings-fees-divergence/milestone.md`, `/F-01-camada2-fee-ingest/feature.md`, `/F-02-stock-divergence/feature.md`, `/F-03-anuncios-fe-contract/feature.md`, `/validation-contract.md`; `M-06-orders-backfill-decomposition/milestone.md`, `/F-01-backfill-incremental/feature.md`, `/F-02-decomposition-camada3/feature.md`, `/F-03-audit-fe-pedidos/feature.md`, `/validation-contract.md`; `M-07-pricing-fee-read/milestone.md`, `/F-01-fee-read-resolver/feature.md`, `/F-02-precos-provenance-fe/feature.md`; `M-08-webhook-ingest/milestone.md`, `/F-01-inbox-endpoint/feature.md`, `/F-02-worker-callback/feature.md`, `/validation-contract.md`; `M-09-sync-observability/milestone.md`, `/F-01-sync-health-endpoint/feature.md`, `/validation-contract.md`.

**Unread-path sweep (manifested, NOT read — declared, not assumed):** `research/codebase-ingest-side.md`; `research/p5-prerequisites.md`; `planning-reviews/p3-reconciliation-r01.md`; `planning-reviews/p5-reconciliation-r01.md` … `r08.md` (8 files); `M-01-ml-client-hardening/validation-contract.md`; `M-03-orders-shipment-persist/validation-contract.md`; `M-04-listings-backfill-ingest/validation-contract.md`; `M-07-pricing-fee-read/validation-contract.md`; `M-09-sync-observability/F-02-integracoes-health-section/feature.md`. Neither ★2 defect can be repaired by these: 2.1 is a fork between two contract-and-feature statements that a reconciliation log cannot overwrite, and 2.2's locus is IC-03's key list plus the M-06 gate, both read in full. `M-07/validation-contract.md` being unread is the one residual risk for 2.2 — it could add /precos provenance, but /precos is already covered by MIS07-C7 and the defect is on /pedidos (M-06), which was read.

**Repo files read solely to verify factual claims (absolute paths):** `C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\internal\composition\root.go`; `...\internal\modules\orders\application\enrich_service.go`; `...\internal\platform\httpx\route_deadline.go`; `...\internal\modules\sync\application\scheduler.go`; `...\internal\modules\listings\transport\query.go`; `...\internal\modules\listings\transport\http_handler.go`; `...\internal\modules\listings\domain\filter.go`; `...\internal\modules\orders\domain\order_bucket.go`; `...\internal\modules\connectors\adapters\mercado_livre\capability_adapter.go`; `...\internal\platform\migrate\runner.go`; `...\apps\server_core\migrations\` (listing).
