# MIS-006 Research — Contract Inventory, Decisions, Scenario IO, Flaw Scope

Sources read in full: INTEGRATION-DATA-CONTRACT.md, SCENARIO-WALKTHROUGH.md,
ML-API-QUERY-CATALOG.md, SYSTEM-BLUEPRINT.md, memory ml-api-live-verdicts-d120,
memory lenient-xlsx-import-two-source. Base sha 138aac3d.

---

## (1) Contract inventory E1–E7

| # | Name | Purpose | Current shape (summary) | Mission concern coverage | Verdict |
|---|------|---------|--------------------------|---------------------------|---------|
| E1 | Conta ML (Installation) | OAuth identity + token lifecycle | seller_id/nickname/site/email; token+expiry; connection state machine; onboarding saga (connect→backfill→link→collect) | none of the 4 (product_mirror/adapter/vinculo/sync_state) directly — onboarding saga *sequences* E4/E6 but those are out of MIS-006 scope | **UNTOUCHED** by this mission (belongs to ML-sync/onboarding missions) |
| E2 | Produto interno (catálogo) | Canonical product model both adapters converge to | 10 required fields: codigo_produto, custo, estoque, local, marca, grupo, preco, ean, referencia, descricao. Sankhya column hypotheses `[TESTAR-SKW]` unconfirmed; xlsx PT aliases to be guaranteed in lenient parser. Only #1 blocks a row (proposed); others = NULL honest-unknown w/ completeness badge | **THIS IS the products_mirror contract** — direct hit on source-adapter + products_mirror concern | **NEEDS EXTENSION**: current mirror model only has 1,2,3,8,9,10 (per SYSTEM-BLUEPRINT §2); fields 4 (local), 5 (marca), 6 (grupo), 7 (preco) require schema migration. This is THE core contract for MIS-006. |
| E3 | Anúncio (listing ML) | Full listing data from ML | mlb_id/title/status/sub_status/price/qty/category/condition/photos/dates/sku/ean/variations/shipping/catalog_product_id/health/tags — ~20 fields, many `✗ ADICIONAR`. Derived: comissao, frete_gratis_custo, visitas, preco_venda_atual | vínculo chain touches this (match target for E4), but backfill/sync itself is E3 territory = ML-sync mission | **UNTOUCHED core** (listings sync is non-scope); only relevant as the READ side of the EAN-match in E4 |
| E4 | Vínculo (anúncio↔produto) | Link internal product to ML listing | Candidates + human approval model; EAN+SKU=ALTA; contradiction beats EAN. Planned changes: auto-trigger snapshot+candidate-gen after backfill/import; EAN-exato-único→auto-approve `[DECISÃO D4 RATIFICADO]`; no-EAN listing→REVIEW w/ visible reason | **DIRECT hit** — vinculo chain concern, core to MIS-006 (auto-link audit) | **NEEDS EXTENSION**: policy amendment (auto-approve unique exact EAN) + trigger wiring + audit trail are new mission work. Idempotency (A8: unique key `(internal_product_id, provider_listing_id)`, respect manual override) also new. |
| E5 | Pedido | Order full contract (items, buyer, fiscal, financial breakdown, logistics/SLA) | Very large contract — identity/state, line items, buyer+NF, financial decomposition (valor_sobrado calc), logistics/SLA. Backfill 12mo (D5 ratified), incremental via date_last_updated cursor | **OUT of scope** — explicitly ML-sync/orders mission (non-scope list: "Sync de pedidos + decomposição + Fila/SLA — missão ML-sync") | **UNTOUCHED** |
| E6 | Concorrência (market) | Competitor offer data + aggregates | Per-offer: item_id, price, seller_id, seller identity, sold_quantity (proved impossible via 403), shipping, listing_type, official_store. Aggregates: quantos_concorrentes, min/median/max, buy_box_winner, price_to_win, freshness | **OUT of scope** — "Mercado/concorrência + tarifas live — missão mercado" (non-scope list) | **UNTOUCHED**, but E6's consumption of products_mirror (Oportunidades join) is the reason MIS-006 exists — mirror is upstream dependency, not itself E6 work |
| E7 | Produto de catálogo ML + tarifas | Catalog identity resolution + commission/fee tables | catalog_product_id resolution, product name/attrs/category, buy_box status; tariffs by category/listing_type/price-range (ml_tariffs) | F3.7 discovery route (`descobrir_produto_catalogo(ean)` = F3.1 in query catalog) is **explicitly in MIS-006 scope** ("Rota de descoberta F3.7... entra na missão"); tariff sync itself is out (mercado mission) | **PARTIAL EXTENSION**: only the EAN→catalog_product_id discovery lookup (F3.1/F3.2) is in scope, as the mechanism that breaks chicken-egg for unlinked products. Full E7 (tariffs, buy_box) stays untouched. |

**Gaps needing a NEW contract (E8+):**
- **sync_state** — no existing E-contract models this. SYSTEM-BLUEPRINT §4 and mission.md require `sync_state` (cursor per account+entity, last_error, retries) as foundation-skeleton scope. Needs its own E8 entry: fields = tenant, entity_type, cursor, last_run_at, last_error, status.
- **Active-source config (per-tenant)** — contract section 1b describes it prose-only ("Fonte ativa = configuração de tenant em banco"), never tabled as an entity. Needs E9: tenant_id, active_source (sankhya|xlsx|catalogo_cliente), set_at, set_by.
- **Auto-link audit trail** — E4 mentions "audit" repeatedly (auto_exact_ean audit note in SCENARIO-WALKTHROUGH Fase 2) but no field-level schema given anywhere. Needs E10 (or E4 extension): link_id, rule_matched (exact_ean/manual/etc), actor (system|operator), created_at, superseded_by (for A8 idempotency/override tracking).
- **erp_import_protocols / erp_import_products** — exist in code already (0046 migration, lenient-xlsx memory) but are not formalized as an E-contract in INTEGRATION-DATA-CONTRACT.md; SYSTEM-BLUEPRINT §2 table lists them informally. Worth promoting to a formal E-number if the mission writes a validation contract citing them.

---

## (2) Decisions needing resolution

### D3 — Campos fiscais (fiscal fields)
Status in doc: **Pendente**, recommendation given: "honesto-desconhecido + badge de completude."
Evidence: E2 doesn't list a dedicated fiscal field (no NCM/CFOP/CST in the 10-field table) —
fiscal fields belong mostly to E5 Pedido (doc fiscal CPF/CNPJ, nome/razão social, IE — via
`billing-info`, proven live T4 in ML-api-verdicts memory: `x-version:2` OK, IE presence for PJ
still unsampled). For products_mirror (E2) itself, NO fiscal field is named as required —
the "10 campos obrigatórios RATIFICADO D-120" table is closed and does not include NCM/tax
class. **Recommendation**: D3 as scoped to E2/products_mirror is effectively **N/A / already
resolved as "not required"** — the 10-field table is ratified closed. If D3 is about
E5 fiscal capture (NF), that's out of MIS-006 scope (Pedido = ML-sync mission). Flag this
ambiguity for planning: D3 may be a stale decision-id inherited from before the E2 scope was
ratified to exactly 10 fields.

### D6 — Cadência de coleta (collection cadence)
Status: **Pendente**, recommendation: "diária, vinculados + watchlist."
Evidence: SYSTEM-BLUEPRINT §4 scheduler diagram shows `mercado: diário (D6 pendente)` alongside
pedidos 5min / anúncios diário / tarifas semanal — so the *pattern* (daily cadence, scheduled
via sync_state cursor) is already assumed elsewhere in the design even though D6 itself is
unratified. Collection (E6) itself is non-scope for MIS-006, BUT the `sync_state` skeleton
MIS-006 must build needs to accommodate a cadence field/enum generically (not hardcode any one
entity's schedule) since D6 isn't ratified yet. **Recommendation**: sync_state schema should be
cadence-agnostic (store interval or cron per entity_type, not hardcode "daily"); actual D6
ratification (daily vs other) can happen in the mercado mission without touching sync_state's
shape.

### F3.7 viability — `descobrir_produto_catalogo(ean)`
**UNPROVEN.** From ML-API-QUERY-CATALOG.md F3.1 (`identidade_catalogo(ean)`) is marked DOC only
(not LIVE-verified) — confirmed only against official docs, never tested live. The
"VEREDICTOS LIVE" section (T1-T12) in the same file does NOT include a test of F3.1/EAN-search;
the live round covered T1-T12 (items, orders, shipments, competitors, tariffs) but skipped
catalog-by-EAN entirely. SCENARIO-WALKTHROUGH.md ADENDO 1 explicitly flags: "F3.7 NOVA — rota a
provar T13" and Adendo A1: "F3.7 é hipótese até provar... Rodada live obrigatória antes de
ratificar." A separate live round (**T13–T16**, distinct from the already-executed T0–T12) is
defined in SCENARIO-WALKTHROUGH.md §"Rodada LIVE 2":
- **T13**: `GET /products/search?site_id=MLB&status=active&product_identifier={EAN}` against 10
  real EANs from the client's mirror (mix with/without our own listing) → does F3.7 exist? what
  shape? EAN hit-rate?
- **T14**: for each catalog_product found, `GET /products/{id}` + `GET /products/{id}/items` →
  is demand data sufficient for an Oportunidade row?
- **T15**: public search fallback `GET /sites/MLB/search?q=` for non-catalogable products —
  check whether it is PolicyAgent-blocked like `/items` (T8/F2 precedent: 403 on third-party
  reads is a proven pattern for this account/policy).
- **T16**: full simulation for product B: T14 aggregates + `ml_tariffs` sweep + shipping (F4.2)
  → is potential margin end-to-end computable?

Read-only GET endpoints only; EANs to pull from `erp_import_products` of protocol **#004-E**
(prospect catalog, 2012 products). Tool: extend `apps/server_core/cmd/mlprobe` (round 3) —
confirmed present at `apps/server_core/cmd/mlprobe/{main.go,followup.go}` (untracked, per repo
git status). Evidence output target: `docs/design/evidence/ml-api/`.
**Recommendation**: MIS-006 planning must run T13–T16 via mlprobe before ratifying F3.7 as part
of the mission's discovery-route deliverable; mission.md's own Outcome section already
conditions this ("UNPROVEN, exige rodada live T13-T16 no planning").

### Stale-policy (F-XLSX-1)
**Confirmed operator-decided = keep-absent + stale flag, no physical delete.**
Direct citations:
- INTEGRATION-DATA-CONTRACT.md §1b: "Fonte-arquivo (xlsx): TEM que persistir... Snapshot
  versionado por protocolo de import."
- SCENARIO-WALKTHROUGH.md FASE 5: "xlsx novo protocolo #006-E sem produto B → mirror rebuild
  marca B ausente (política: soft-flag `absent_in_last_snapshot`, nunca delete silencioso)."
- mission.md line 45: "upsert-merge keep-absent (F-XLSX-1), NULL=honesto-desconhecido" — cited
  directly as ratified mission scope, sourced from `docs/design/evidence/INTEGRATION-FINDINGS-D120.md`.
This is the strongest-evidenced of the four items — treat as **already resolved**, not a live
decision; MIS-006 implementation must implement `absent_in_last_snapshot` as a soft-flag column
on `products_mirror`, never a DELETE on rebuild.

---

## (3) Scenario IO pairs (exemplo-IO markers)

### Product A — Sankhya, already listed
- Input: `codigo_produto=90008`, descricao "Torneira Cozinha Mesa Bica Móvel",
  `ean=7891234500017`, custo=R$82.00, estoque=120 (local CD-01), preco(tabela ERP)=R$169.00,
  marca=Ferragens, grupo=Torneiras. Source: Sankhya Oracle (`TGFPRO⋈TGFEST⋈TGFCUS⋈TGFBAR`).
- Expected mirror row: `products_mirror(tenant, 90008) → {custo 82.00, estoque 120, ean
  ...017, ...}`, `sync_state(products)` cursor+timestamp updated.
- Expected link outcome: EAN `...017` matches exactly 1 existing listing `MLB4735304125` →
  **AUTO-APPROVED** link (D4 policy) + audit note `'auto_exact_ean prot#005'` + market
  collection auto-enqueued via `catalog_product_id` of our own listing (route F3.3, already
  proven live).
- Downstream render (Radar tab): NOSSO PREÇO R$169.00 (from listing, not custo — fix M2), best
  competitor R$154.90 (seller MAPRON, 3_yellow), position 3rd of 4, collected 2h ago.
- Pedido flow: sold → decomposition `R$169 − R$27.89 commission(16.5% unit) − R$19.85 shipping
  − R$82 cost = R$39.26 net (23%)`.

### Product B — xlsx, unlisted → Oportunidade
- Input: `codigo_produto=74606`, descricao "Misturador Monocamando Lavatório Docol",
  `ean=7897586000745`, custo=R$410.00, estoque=35 (local LOJA), preco(tabela ERP)=R$899.00,
  marca=Docol, grupo=Misturadores. Source: xlsx upload → protocol `#005-E` (parsed via lenient
  parser, PT aliases/preamble/multi-sheet union) → `erp_import_products` immutable history →
  mirror rebuild.
- Expected mirror row: `products_mirror(tenant, 74606) → {custo 410.00, ...}`.
- Expected link outcome: **NO listing** of ours carries EAN `...745` → **NOT an error** — this
  IS the Oportunidades input. Enters `sync_state(catalog_discovery)` queue for F3.7 EAN-lookup
  (`GET /products/search?product_identifier=7897586000745` → hypothetical
  `catalog_product_id=MLB22624877`, or empty = non-catalogable, see Adendo A3).
- Discovery result (if F3.7 proven): `GET /products/{id}/items` → 5 offers R$780–950 → written
  to `market_aggregates` for codprod 74606 even though we don't sell it.
- Oportunidades tab render: Produto 74606, Estoque/Custo 35un/R$410, Mercado "5 vendedores,
  R$780–950", Margem potencial "sell at R$849 → ~R$214 net", Veredicto "You don't sell — create
  listing?", Volume proxy = leader's 4,400 historical transactions (sellers_cache, not real
  sold_quantity — T8/F2 limit).
- Query shape for the tab (already given as canonical join):
  `SELECT pm.*, ma.* FROM products_mirror pm JOIN market_aggregates ma ON ma.codprod=pm.codigo_produto LEFT JOIN product_links pl ON pl.internal_product_id=pm.codigo_produto WHERE pl.link_id IS NULL AND pm.estoque_total > 0 ORDER BY ma.margem_potencial DESC`
- **Real-world caveat (A6):** actual prospect protocol #004-E (2012 products) has custo/estoque
  NULL (not R$410 like the illustrative example) — margin ranking must fall back to
  demand-based ordering with label "margem indisponível — custo não informado" when custo is
  NULL, per ADR-17 (never fabricate custo=0).

---

## (4) Flaws F1–F15 relevance to MIS-006

| Flaw | One-line | In/Out scope |
|------|----------|---------------|
| F1 | Oportunidades 100% client-side, no backend | **OUT** — Oportunidades/mercado query backend is the "mercado" mission (E6), not MIS-006; but MIS-006's products_mirror is the prerequisite table F1's eventual fix will JOIN against |
| F2 | "Monitorado" scope doesn't exist | **OUT** — mercado mission (scope/coverage of collection targets) |
| F3 | "Não vendemos" exclusion absent (no LEFT JOIN product_links) | **OUT** structurally (query lives in mercado/Oportunidades backend) but **IN** dependency — product_links (E4) is MIS-006 scope, so this flaw's fix becomes possible only after MIS-006 |
| F4 | Collection is 1-codprod-manual, no scheduler/batch | **OUT** — E6 collection scheduling is mercado mission |
| F5 | Auto-approve EAN link absent, ACCEPT is label-only | **IN** — directly D4 policy, E4 vínculo chain, core MIS-006 deliverable |
| F6 | MASS-CLOSURE on partial listing pull | **OUT** — E3 listings sync, ML-sync mission |
| F7 | date_created of listing never captured | **OUT** — E3 |
| F8 | Sequential refresh, no multiget | **OUT** — E3/F1.2 |
| F9 | Order import page-1/limit-20, no date/sort | **OUT** — E5 |
| F10 | Enrich N+1 on read (~10s) | **OUT** — E5 order enrichment |
| F11 | Decomposer wired nil | **OUT** — E5 financial decomposition |
| F12 | Competitor item_id + winner identity dropped | **OUT** — E6 |
| F13 | 4 divergent winner structs | **OUT** — E6/E7 |
| F14 | `market_*.product_id` has no FK to mirror | **IN (partial)** — the FK target (`products_mirror`) is MIS-006's deliverable; adding the FK itself may be MIS-006 scope (cheap, since mirror is being created) even though `market_*` tables belong to mercado mission — flag for planning as a possible low-cost inclusion |
| F15 | Orphaned listing on account switch | **OUT** — E1/E3 (installation lifecycle) |
| (unlabeled) | sold_quantity of competitor impossible (403) | **OUT** — E6, already resolved as honest-unknown, no action needed |

**Chicken-egg (the mission's namesake defect)**, per mission.md: `Collect(codprod)` resolves
market via `LinkedListings(codprod)` (`collection_pipeline_service.go:269`) → product without a
link has no path to market. This isn't numbered F1-F15 but is the root cause MIS-006 exists to
fix — directly addressed by F3.7 discovery route + products_mirror + auto-link chain (E2+E4+E7-partial).

---

## Compressed summary (returned to caller)
