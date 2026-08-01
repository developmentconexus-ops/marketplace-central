# Interface Contract — listings E3 + listing_variations

```yaml
id: IC-07
type: interface-contract
status: planned
owner: Mission Strategist
parent: MIS-007
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: support
```

## Boundary

Shape persistido de anúncios entre: M-04 (writer único — backfill/scheduler/refresh),
M-05 (consome p/ camada 2 + divergência de estoque), FE /anuncios (colunas ratificadas
design §8), matcher de re-vínculo (AbsorbProviderSnapshots). DDL aditiva no M-04 (range
próprio) — NÃO no M-02, porque `listings` já existe e o único writer é do M-04.

## Why This Contract Exists

M-04 e M-05 são milestones distintos sobre a mesma tabela; e o design §7 criou ambiguidade
que 2 workers resolveriam diferente: E3 lista `commission_amount/pct, free_shipping_cost`
em `listings` E define `channel_fees` — duas casas pro mesmo fato.

**Decisão de planning (registrada aqui, autoridade: "decisão fina no planning" — design §7
item 6):** fee mora SÓ em `channel_fees` (IC-01, ADR-09 — proveniência obrigatória);
`listings` NÃO ganha colunas de comissão/frete; o DTO de /anuncios COMPÕE listing +
channel_fees camada 2 (com coletado_em, que a tela ratificada exibe). Dual-write = seam de
drift, morto aqui.

## Resources Or Entities

- `listings` — colunas aditivas E3 (abaixo); PK NÃO muda (sentinela `'-'` mantida — ADR-13).
- `listing_variations` — NOVA child table (M-04).
- Writer único: `ApplyCompletedPull` re-semantizado (ADR-06) atrás de `IngestListing`
  (IC-06).

## Operations

| Operation | Trigger | Input | Output | Notes |
| --- | --- | --- | --- | --- |
| IngestListing | backfill/scheduler diário/refresh em lote | provider_listing_id (hidratação multiget 20) | upsert listing + variations + camada 2 (M-05) + divergência estoque (M-05) | writer ÚNICO; status = verdade do provider SOMENTE |
| ListListings (existente) | GET /listings | filtros atuais + `filter.divergentes` (IC-02; idioma `filter.` — `query.go:31`) | DTO atual + campos E3 + fees compostos + divergences | ordenação existente preservada; paginação existente |

## Fields

### Required Inputs — colunas aditivas `listings` (todas NULL = honesto-desconhecido)

Do design §7 E3, verbatim: `sold_quantity int`, `category_id text`, `condition text`,
`permalink text`, `thumbnail text`, `date_created_ml timestamptz`, `tags text[]`,
`catalog_product_id text`, `shipping_mode text`, `free_shipping boolean`,
`logistic_type text`, `available_quantity int` (estoque ML, grão listing quando sem
variação).

EXCLUÍDOS do E3 (decisão acima): `commission_amount/pct`, `free_shipping_cost` → IC-01.

Suporte ao ciclo de vida (ADR-06): `last_seen_at timestamptz`, `absent_since timestamptz`
NULL; `raw jsonb` NULL (payload `/items` — permitido, ADR-03) + marcador de truncamento.

### Required Inputs — `listing_variations` (nova)

- PK `(tenant_id, installation_id, provider, provider_listing_id, variation_id)` — tuple
  verbatim do ADR-13 (P7 r02 ★2-A: `installation_id` é OBRIGATÓRIO no tuple porque o PK
  real do parent `listings` é `(tenant_id, installation_id, provider_listing_id,
  variation_id)` — `0036_listings.sql:2-31`, `research/codebase-ingest-side.md:65` — e a
  child sem `installation_id` não endereça rows do parent por instalação; `provider`
  mantido conforme A-13 ratificada, `p3-reconciliation-r01.md:37`).
- `price numeric(14,2)` NULL; `available_quantity int` NULL; `sold_quantity int` NULL;
  `seller_sku text` NULL; `attributes jsonb` NULL; `last_seen_at`; `absent_since` NULL;
  created/updated.
- SEM raw próprio (raw do listing pai cobre).

### Required Outputs

DTO /anuncios: colunas ratificadas design §8 — vendidos, data criação ML, categoria, tags,
catálogo, estoque ML × ERP com badge ⚠, comissão+frete por anúncio COM data de coleta
(composto de IC-01), filtro "divergentes".

## Enums And Statuses

- `status` de listing: verbatim do provider (`active`,`paused`,`closed`,... sem CHECK —
  vocabulário do ML).
- Ausência ≠ closed: `absent_since` marcado só após run COMPLETO (IC-06); `status` NUNCA
  inferido de ausência.

## Error Cases

Rotas existentes; envelope vigente; nenhum code novo além do filtro `divergentes`
(parâmetro inválido cai no code de validação existente).

## Error Matrix

| Case | Status | Code | Notes |
| --- | --- | --- | --- |
| `filter.divergentes` não-booleano | 400 | code de validação vigente | mesmo envelope das outras rotas |
| `POST /listings/refresh` com run em curso | 409 | `refresh_in_progress` (existente) | comportamento PRESERVADO, medido `codebase-read-side.md:87`; 202 quando livre (row adicionada por referência — P7 r02 A-12) |

## Persistence Expectations

- Upsert por row (idioma writer.go:74-95 `upsertSQL` + keep-absent `:104-112` — F-r06-5);
  MASS-CLOSURE (`repository.go:390-394`) MORTO —
  substituição DENTRO do M-04, ANTES do cursor retomável (R-B: cursor antes da semântica =
  catalog-wiper).
- `available_quantity`: grão variação quando variações existem (soma NÃO materializada no
  pai — leitor soma); grão listing caso contrário. Divergência de estoque usa o grão do
  vínculo (IC-02).
- Re-vínculo pós-backfill: hidratação nova CONTINUA alimentando `AbsorbProviderSnapshots`
  (âncoras EAN/SKU) — acoplamento nomeado ADR-13; snapshots/âncoras não-regressivos vs
  pull pré-mudança.

## Canonical Examples

Upsert de listing com variação (dois níveis, mesma transação):

```json
{"listing":{"provider_listing_id":"MLB123","status":"active","category_id":"MLB1055",
 "sold_quantity":37,"available_quantity":null,"tags":["good_quality_thumbnail"],
 "logistic_type":"fulfillment","free_shipping":true},
 "variations":[{"variation_id":"181","price":79.9,"available_quantity":5,"seller_sku":"SKU-A"}]}
```

Rejeição canônica: run incompleto tenta marcar absent → recusa (must-fail IC-06).

## Database Shape

- Migrações aditivas no range do M-04; teste regex por migração.
- `listing_variations` PK acima; índice `(tenant_id, provider_listing_id)`.
- PKs de `product_link_listing_snapshots` (0022) e `product_links` (0025) intocadas.

## Seed Data

Nenhum. Fixtures: catálogo >1 página de scan (R-3); listing com e sem variação; fixture
de abort pós-página-1 (ADR-06 must-fail).

## Timestamp And ID Semantics

- `date_created_ml` = do provider; `last_seen_at` = clock nosso no upsert;
  `absent_since` = clock nosso na marcação pós-run-completo.

## Compatibility Rules

- Colunas E3 futuras = aditivas.
- Consolidação de PK (tirar sentinela `'-'`) = missão futura nomeada, NUNCA nesta (ADR-13).

## Route Namespace

Nenhuma rota nova; DTO/param novos em /listings = par OpenAPI/SDK no milestone FE
respectivo (M-05; ADR-14 emendado P5 r03 P-2 — commit de contrato serializado pelo hub,
≤1 COMMIT FE em voo).

## Transport And Integration

Refresh manual em lote: rota batch existente (202 async) — chip nunca sobe servidor;
live-drive de priced/refresh é do hub.

## Must Preserve

- Writer único; PK com sentinela; snapshots observer alimentado; paginação/ordenação
  existentes de /listings; classe batch p/ pull.

## Must Not Decide In Feature Execution

- Exclusão das colunas de fee do E3 (decidido AQUI); grão de estoque; shape de variations;
  semântica absent; local da DDL (M-04, não M-02).

## Validation Impact

- Must-fail ADR-06: abort pós-página-1 → zero flips `closed`, fixture >1 página.
- Must-fail ADR-13: "snapshot observer starved" — âncoras não-regressivas.
- M0X-U*: /anuncios mostra colunas E3 reais + estoque ML×ERP + badge + filtro, dirigido em
  browser; comissão exibida vem de channel_fees com coletado_em visível (compostura provada
  na tela, não no JSON).
