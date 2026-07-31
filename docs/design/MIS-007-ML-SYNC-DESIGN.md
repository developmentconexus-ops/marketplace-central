# MIS-007 — ML-Sync: Design Ratificado (brainstorm 2026-07-31)

Spec de design validada em sessão de brainstorming com o operador. Insumo direto do
`/mission-planning` da MIS-007. Toda decisão abaixo foi aprovada explicitamente; fontes
citadas por doc oficial (context7 `/websites/developers_mercadolivre_br` + web) ou medição
live própria (D-120/D-121).

Relacionado: [SYSTEM-BLUEPRINT.md](SYSTEM-BLUEPRINT.md), [STORAGE-SCHEMA.md](STORAGE-SCHEMA.md),
[evidence/INTEGRATION-FINDINGS-D120.md](evidence/INTEGRATION-FINDINGS-D120.md),
`.mnfs/MIS-006-integracao-fundacao/mission.md` (fundação fechada).

---

## 1. Objetivo e critério de sucesso

**Operação diária confiável.** O vendedor abre /anuncios e /pedidos de manhã e tudo está lá,
atualizado sozinho, rápido (<2s), com margem real por pedido. Zero botão "refresh" obrigatório.

Contexto: MIS-006 fechou a fundação (products_mirror, adapters ERP, vínculo, sync_state,
erro unificado). O lado ERP existe e está provado. O lado ML ainda não foi ingerido — esta
missão traz anúncios e pedidos para dentro, do jeito nativo, sem repetir a implementação
desconexa antiga (telas chamando API externa no read, consumidores lendo fontes divergentes).

Construção base→topo, passo a passo: esta missão é SÓ LEITURA do ML. Escrita (estoque, preço,
publicação, optin de catálogo) é missão futura — mas o desenho nasce write-ready.

## 2. Escopo / Não-escopo

**Dentro (duas ondas + chip inicial):**
- Onda 0: orders-decoupling — matar os 4 sítios que chamam ML no read.
- Onda 1: sync de anúncios (backfill + contínuo + colunas E3 + estoque ML + comissão/frete do
  próprio anúncio + divergência de estoque + re-vínculo).
- Onda 2: sync de pedidos (backfill 12 meses + incremental 5min + order_shipments +
  decomposição persistida + Fila/SLA + auditoria de tarifa) + webhook ingest.
- Pré-requisitos de adapter (fecham buracos F-ADAPTER-1): backoff exponencial + jitter +
  honrar Retry-After no 429; regra DTO `Raw json.RawMessage`.
- Fonte ativa ERP: sankhya (F-LINK-1 resolvido — Sankhya 31/31 sob predicado vendável).

**Fora (com destino nomeado):**
- Mercado/concorrência, catalog offers, sellers, oportunidades, `ml_tariffs` sweep → MIS-008.
- Simulação de produto NÃO vinculado (categoria por preditor + frete manual) → MIS-008,
  par natural de Oportunidades.
- Tabela de-para categoria nossa↔ML + preditor de categoria → missão de publicação (write).
  Nesta missão categoria vem DE GRAÇA no ingest do anúncio existente (`category_id`).
- Writes ML de qualquer tipo (estoque, preço, catalog optin) → missão futura; write vivo
  exige autorização explícita do operador.
- Onboarding saga completa (tela de progresso) → depois; `sync_state` já grava progresso.

## 3. Arquitetura — núcleo nativo × adapter (aprovada)

Sistema é NATIVO: núcleo tem os conceitos; adapters só conversam com terceiros. Núcleo nunca
importa tipo do ML; adapter nunca vaza payload cru para cima.

**Conceitos do núcleo (agnósticos de provider — valem para Shopee amanhã):**
- **Publicação** (listing): presença de um produto num canal. Já existe; ganha colunas E3.
  Chave: `tenant + installation + provider_listing_id` — provider é dado, não tipo.
- **Custo de canal**: comissão + frete + taxas, em 3 camadas de verdade (ver §4). Cada valor
  persiste com camada + proveniência + data. Núcleo só conhece "camada mais forte ganha".
- **Divergência**: conceito único (estoque ERP×ML, tarifa realizada×estimada; extensível a
  categoria descontinuada). Mesmo shape, mesma UI de aviso (badge).
- **Envio** (order_shipments): SLA, custo real, endereço. Nativo; ML só alimenta.
- **Inbox de notificações**: entrada de webhook de QUALQUER provider (ver §5).
- **sync_state**: já existe, cadence-agnostic; ganha entidades `listings` e `orders`.

**Só o adapter ML sabe:** que comissão vem de `listing_prices`, frete de `shipping_options`,
verdade do pedido de `order_items[].sale_fee`; paginação scan/multiget; headers `X-Costs-New`
/ `x-format-new`; backoff/429; DTOs com `Raw json.RawMessage`.

**Teste da fronteira (critério de aceitação arquitetural):** adicionar um provider novo =
escrever 1 adapter, zero mudança em tabela, tela ou serviço de núcleo. Item da missão que
violar isso = design errado.

## 4. Ciclo de vida de comissão/frete — 3 camadas (aprovado, validado por benchmark)

| Camada | Quando existe | Fonte (adapter ML) | Uso |
|---|---|---|---|
| 1. Cotação (quote) | produto sem anúncio | `GET /sites/MLB/listing_prices?price&category_id&listing_type_id` (+`billable_weight`,`logistic_type`,`shipping_mode`) → `sale_fee_amount`, `percentage_fee`, `fixed_fee`, `financing_add_on_fee` | simulação pré-anúncio (MIS-008) |
| 2. Anúncio (listing) | produto anunciado | comissão: `listing_prices` reconsultado com categoria+preço atuais do item; frete: `/items/{id}/shipping_options` (campo `free_options`) | simulação do já-anunciado; estimativa de margem |
| 3. Pedido (realized) | venda real | `order_items[].sale_fee` (POR UNIDADE — medição live T2, doc não declara) + `/shipments/{id}/costs` com header `X-Costs-New` | verdade absoluta; margem REAL |

Regras:
- Camada mais forte disponível ganha.
- **Auditoria 3→2**: pedido chega com sale_fee real ≠ comissão estimada do anúncio →
  divergência gravada + aviso. Nunca sobrescrever silenciosamente — divergência é informação.
- Todo valor com proveniência (camada, origem, coletado_em). Número sem origem não existe.
  O seed estático 16%/22% (`fee_sync.go:29`) morre.
- **Frete NÃO tem camada 1**: não existe endpoint de cotação por dimensão/peso sem `item_id`
  (pesquisa confirmou ausência na doc). Simulação pré-anúncio mostra frete
  honesto-desconhecido ou campo manual — nunca número inventado (ADR-17).

Benchmark que sustenta o modelo: Bling (líder) usa o valor real cobrado importado do pedido
(nossa camada 3); AnyMarket usa regra percentual manual (modelo que já abandonamos); nenhum
hub confirmado dá aviso proativo de mudança de tarifa — a auditoria 3→2 é diferencial.

## 5. Webhook + reconciliação (aprovado, com fatos verificados)

Padrão: **webhook é o mecanismo primário** (latência de segundos), **scheduler é a
reconciliação** (completude garantida). Nunca um sem o outro.

Fatos da doc oficial (context7):
- Callback = URL pública, POST; payload traz só `resource` + `topic` + `user_id` + `attempts`
  — o dado NÃO vem; buscar o resource depois. Timestamps UTC.
- ML tenta entregar **8 vezes em ~1 hora** antes de considerar perdida; perdidas ficam em
  `GET /missed_feeds` por **2 dias**. Buraco >2 dias = só reconciliação salva.
- Topic recomendado para vendas: `orders_v2`. Existem `items`, `shipments`, `messages`.
- Filtro opcional por IP de origem: 54.88.218.97, 18.215.140.160, 18.213.114.129, 18.206.34.84.
- Limite de resposta "200 em 500ms": NÃO confirmado na doc indexada (exemplo oficial mostra
  498ms aceito). Design independe: responder 200 imediato, processar depois.

Componentes:
- `POST /webhooks/{provider}` — transport fino: valida, grava em `notifications_inbox`,
  responde 200 imediato (milissegundos).
- Worker consome inbox → chama o MESMO ingest idempotente do scheduler (upsert por resource
  id). Webhook e scheduler são duas portas do mesmo caminho.
- Inbox = auditoria + dedupe (por `_id`/attempts) + retry.
- Registro do callback no app ML: domínio ngrok FIXO do operador (meses estável, tratado como
  URL de produção; troca futura de domínio = 1 edição no cadastro do app, zero código).
- Scheduler continua obrigatório nesta missão (pedidos 5min, anúncios diário); relaxar
  cadência só depois que webhook provar saúde, em missão futura.

## 6. Fluxo por onda (aprovado)

**Onda 0 — orders-decoupling (chip pequeno, primeiro):** os 4 sítios que chamam ML no read
morrem. Tela lê Postgres, ponto.

**Onda 1 — Anúncios:**
1. Backfill: scan ids (1000/batch, scroll_id; coletar ids primeiro, hidratar depois) →
   multiget 20 → upsert `listings` E3 + `listing_variations`. Retomável por cursor.
2. Scheduler diário + refresh manual em lote na tela.
3. Ingest puxa junto: comissão camada 2, frete, estoque ML (`available_quantity`).
4. Divergência de estoque calculada no INGEST (não no read): disponível ERP
   (estoque−reservado, corte vendável) via vínculo × estoque ML → flag persistida → badge.
5. Re-vínculo pós-backfill: matcher roda contra anúncios completos e frescos.

**Onda 2 — Pedidos:**
1. Backfill 12 meses, cursor `date_last_updated` + `sort=date_desc` (bug `date_asc` provado
   live T7).
2. Scheduler incremental 5min (reconciliação).
3. Ingest por pedido (goroutines com budget): `/shipments/{id}` + `/sla` + `/costs` +
   `billing_info` → `order_shipments`. Multiget de shipments NÃO existe (T5) — GETs paralelos.
4. Decomposição PERSISTIDA no ingest: total − comissão(sale_fee×qty, por unidade) − frete
   real − custo(mirror via vínculo) = líquido + margem%. Custo CONGELADO no pedido.
5. Auditoria camada 3→2 → divergência de tarifa → aviso.
6. Fila/SLA: bucket derivado no ingest, indexado.
7. Webhook ingest: inbox + worker + registro callback (topic `orders_v2` primeiro; `items`
   se sobrar folga).

**Pré-requisito das duas ondas (antes de qualquer backfill):** backoff exponencial + jitter +
honrar Retry-After no 429 (doc oficial recomenda; ~1500 req/min por seller). Regra DTO:
tipado nos campos consumidos + `Raw json.RawMessage` persistido.

## 7. Modelo de dados (aprovado)

1. `listings` ESTENDIDA (E3): sold_quantity, category_id, condition, permalink, thumbnail,
   date_created_ml, tags, catalog_product_id, shipping_mode/free, logistic_type,
   commission_amount/pct, free_shipping_cost, available_quantity (estoque ML) +
   `listing_variations` formalizada.
2. `channel_fees` NOVA — 3 camadas com proveniência (camada, origem, coletado_em).
3. `order_shipments` NOVA — SLA, custos reais, endereço NF, tracking.
4. `orders` estendida — pack_id, decomposição JSONB + líquido + margem%, bucket indexado.
5. `notifications_inbox` NOVA — provider, topic, resource, dedupe, received_at, processed_at.
6. `divergences` NOVA (ou flag por entidade — decisão fina no planning): tipo
   (estoque|tarifa), entidade, esperado × observado, detectado_em, resolvido_em.
7. `sync_state` — entidades novas `listings`/`orders`, sem mudança estrutural.

Transversais: NULL = honesto-desconhecido (nunca 0 fake); custo congelado no pedido; nada de
delete físico; tudo tenant-scoped; migrações aditivas, nunca ALTER destrutivo no mirror.

## 8. Telas (aprovado — nenhuma tela nova)

- **/anuncios**: colunas reais (vendidos, data criação ML, categoria, tags, catálogo),
  estoque ML × ERP com badge ⚠ + filtro "divergentes", comissão+frete por anúncio com data
  de coleta, refresh manual em lote.
- **/pedidos**: <2s, zero call ML no read; margem real persistida; custo congelado; coluna
  SLA; Fila com bucket. Pedido novo em segundos (webhook) ou ≤5min (reconciliação).
- **/precos**: produto vinculado simula com camada 2/3 real; proveniência visível;
  divergência de tarifa ⚠.
- **/integracoes**: saúde do sync por entidade (sync_state) + status webhook (última
  notificação recebida).

## 9. Validação (aprovado — padrão MIS-006)

- Contrato de validação por milestone com critérios dirigidos em browser real (mandato M0X-U*).
- Fixtures multi-página (paginação só se prova com >1 página).
- Live-drive final: conta ML real, backfill completo medido, pedido real decomposto, webhook
  disparado por evento real (mexer em anúncio → notificação → inbox → ingest → tela).
- Must-fail: injetar defeito, teste NOMEIA a falha.
- Divergência provada nas 2 direções: criar → badge aparece; resolver → some.

## 10. Registro de fatos que bindam o planning

- `order_items[].sale_fee` é POR UNIDADE — medição live própria (T2); doc não declara.
- Multiget de shipments NÃO existe (T5) — GETs individuais paralelos.
- Cotação de frete por dimensão/peso sem item NÃO existe na doc — frete pré-anúncio =
  honesto-desconhecido.
- `listing_prices` exige `logistic_type`/`shipping_mode`/`billable_weight` para o valor
  bater com o cobrado (doc de comissão-por-vender).
- Comissão dentro de `/items/{id}/sale_price`: não confirmado — usar `listing_prices`.
- missed_feeds: 2 dias; 8 tentativas/~1h; payload sem dado; `orders_v2` recomendado.
- Rate limit: ~1500 req/min por seller (fonte parcial); doc oficial recomenda backoff
  exponencial + jitter.
- `POST /items/catalog_listings` (optin catálogo de anúncio existente) existe — é WRITE,
  fora desta missão.
- Predição de categoria: `GET /sites/MLB/domain_discovery/search?q=` (limit 1–8) — MIS-008+.
- MIS-008 = mercado/concorrência (herda F3.7, T13-T16, MC-11, ml_tariffs, simulação de
  não-vinculado, oportunidades).

## 11. Decisões da sessão (trilha)

1. Alvo = operação diária confiável (A).
2. Anúncios + pedidos na mesma missão, ondas duras, corte YAGNI se estourar (C).
3. Missão só leitura; write depois; desenho write-ready com avisos de divergência.
4. Divergência de estoque = badge/coluna em /anuncios (A), sem central de avisos.
5. Não-escopo confirmado + exceção: comissão/frete do PRÓPRIO anúncio entra (insumo da margem).
6. Modelo 3 camadas ratificado após benchmark (Bling valida camada 3; auditoria 3→2 é
   diferencial; AnyMarket-style regra manual rejeitada).
7. Categoria: só ingest do `category_id` real (A); de-para e preditor ficam para missão de
   publicação; simulação de não-vinculado → MIS-008 (A).
8. Webhook primário + scheduler reconciliação; ngrok domínio fixo = URL de produção por ora.
