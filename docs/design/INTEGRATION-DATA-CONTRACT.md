# Contrato de Dados da Integração — DRAFT v0 (D-120)

Status: **RASCUNHO para definição conjunta**. Este documento define, entidade por entidade:
(a) quais campos nosso sistema PRECISA, (b) de onde cada campo vem (fonte da verdade),
(c) COMO é alimentado (mecanismo + cadência). Nada implementa nada — é o contrato que a
implementação vai obedecer. Itens `[DECISÃO]` precisam de ratificação do operador.
Itens `[TESTAR]` precisam de prova live contra a API ML antes de virar contrato.

Princípio-mestre (já ratificado): **Postgres local = fonte de verdade para a UI.**
ML e ERP alimentam o banco; telas NUNCA chamam ML no caminho de leitura.
Enrichment acontece no ingest, não no read.

---

## 1. Mapa de fontes

| Fonte | O que fornece | Mecanismo de entrada |
|-------|--------------|----------------------|
| ERP Sankhya (API) | produtos, custo, estoque, fiscal | sync sob demanda + agendado `[DECISÃO D1]` |
| Excel (.xlsx) | mesmos dados de produto, p/ cliente sem API | upload manual em /integracoes (protocolo de import) |
| ML OAuth | identidade da conta (seller_id, nickname, site) | fluxo connect |
| ML REST (pull) | anúncios, pedidos, shipments, catálogo, tarifas | backfill batch + reconciliação agendada |
| ML Webhooks (push) | eventos de mudança (orders_v2, items, ...) | callback público `[DECISÃO D2]` |
| Operador | vínculos aprovados, faturado, parâmetros de preço | UI |

---

## 1b. Módulo canônico de Produto + adapters (RATIFICADO D-120)

Não existe "priorizar xlsx vs Sankhya". Existe UM módulo de Produto nosso (modelo canônico,
seção E2) e N adapters que **convergem** dados de fontes diferentes para esse modelo:

```
xlsx (planilha)  ──parse──►  snapshot persistido ──┐
Sankhya (DB/API) ──query──►  read-through + cache ──┤──► ProductSourceAdapter (port)
[futuro: outro ERP] ────────────────────────────────┘         │
                                                    Modelo canônico de Produto (E2)
                                                              │
                              ┌───────────────┬───────────────┼───────────────┐
                           pricing         vínculos (E4)   mercado (E6)    pedidos (E5)
```

**Estratégia de persistência por TIPO de fonte** (ratificado):
- **Fonte-arquivo (xlsx)**: TEM que persistir — o arquivo é efêmero. Snapshot versionado por
  protocolo de import (modelo atual, manter).
- **Fonte-banco (Sankhya DB/API)**: puxa direto (read-through + cache TTL, modelo atual do
  adapter Oracle) — a fonte é viva, não duplicamos. `[DECISÃO D7]` avaliar snapshot leve
  opcional (histórico de custo/estoque p/ análise + funcionamento offline) — não bloqueia a base.
- Ambos entregam o MESMO contrato pelo mesmo port. Consumidor (pricing, vínculo, mercado)
  não sabe nem pode saber de onde veio.

Fonte ativa = **configuração de tenant em banco** (não env var de boot — `MC_ERP_SOURCE` morre).

**Tela de Importações (observabilidade — ratificado D4):** toda alimentação (upload xlsx OU
sync Sankhya) gera protocolo visível: quando, fonte, quantos produtos, campos presentes/ausentes
(completude: custo? estoque? EAN?), erros, e **o que aconteceu depois** — quantos vínculos
gerados/auto-aprovados, quantos produtos entraram na coleta de mercado. Import não é fim,
é começo de cadeia — e a cadeia tem que ser auditável na UI.

---

## 1c. Por que produto do xlsx nunca aparecia em Oportunidades (cadeia quebrada hoje)

Sua dor exata, rastreada. Para um produto aparecer em Oportunidades ele precisa percorrer:

```
import xlsx → dataset ATIVO → /catalog/products → identidade ML (EAN→catalog_product_id)
   → COLETA de ofertas → aggregate com ≥5 sellers "OK" → aparece
```

Elos quebrados (evidência no POSTMORTEM):
1. **Dataset ativo**: `catalogo_cliente` nunca vira dataset ativo (I2 — `WithActiveSource`
   morto em HEAD); e se `MC_ERP_SOURCE` de boot ≠ xlsx, o backend inteiro lê Oracle e ignora
   seu upload (I1). Verificar qual dos dois te pegou.
2. **Coleta nunca disparada**: import NÃO dispara coleta de mercado. Coleta = 1 clique manual
   por produto dentro da página do produto (S5/M-tabela). Importar 2012 produtos = zero coletas.
3. **Corte silencioso**: mesmo coletado, aggregate exige ≥5 sellers distintos para "OK" (M3);
   1–4 concorrentes = produto some sem explicação.
4. **Vínculo MLB**: para OPORTUNIDADES o vínculo não é o elo (identidade é por EAN direto no
   catálogo ML); vínculo é o elo para REPRECIFICAÇÃO/sinal por anúncio — e a geração de
   candidatos também não é disparada por nada (S4).

O contrato desta base conserta os 4 elos: source por tenant + trigger de cadeia pós-import
(gerar vínculos + enfileirar coleta) + threshold visível + tela de Importações mostrando a cadeia.

---

## 2. Entidades e campos

### E1 — Conta ML (Installation)
| Campo | Fonte | Alimentação |
|-------|-------|-------------|
| seller_id, nickname, site_id, email | ML `/users/me` | no connect (OAuth callback) |
| access/refresh token, expiry | ML OAuth | connect + refresh proativo; **refresh-retry no request path quando expirado** (hoje falha seco — I8) |
| status da conexão | nosso | máquina de estados existente (manter) |

**Onboarding saga (novo, obrigatório):** connect → dispara automaticamente: backfill de
anúncios (E3) → snapshot+geração de candidatos de vínculo (E4) → coleta de mercado em lote
p/ vinculados (E6). Com tela de progresso. Hoje: connect não dispara nada (S5).

### E2 — Produto interno (catálogo) — 10 campos obrigatórios (RATIFICADO D-120)

Contrato mínimo que TODA fonte (Sankhya e xlsx) tem que fornecer. Coluna Sankhya =
hipótese a confirmar com a sessão especialista Oracle/Sankhya antes de ratificar
(`[TESTAR-SKW]`); alias xlsx = a garantir no parser lenient.

| # | Campo canônico | Sankhya (Oracle) | xlsx (aliases PT) | Consumidor | Nota |
|---|----------------|------------------|-------------------|------------|------|
| 1 | codigo_produto | `TGFPRO.CODPROD` | "Código"/"Cód. Produto" | todos | chave interna |
| 2 | custo | `TGFCUS` (qual custo: gerencial? reposição? média?) `[TESTAR-SKW]` | "Custo" | pricing, decomposição, margem de pedido | |
| 3 | estoque | `TGFEST.ESTOQUE` (− reservado?) `[TESTAR-SKW]` | "Estoque"/"Saldo" | pricing, anúncios (sync de qty futuro) | por local |
| 4 | local (depósito) | `TGFEST.CODLOCAL` → `TGFLOC.DESCRLOCAL` `[TESTAR-SKW]` | "Local"/"Depósito" | estoque multi-local | novo — não existe no modelo atual |
| 5 | marca | `TGFPRO.MARCA` `[TESTAR-SKW]` | "Marca" | filtros, matcher (contradição lexical), anúncio novo | novo |
| 6 | grupo | `TGFPRO.CODGRUPOPROD` → `TGFGRU.DESCRGRUPOPROD` `[TESTAR-SKW]` | "Grupo"/"Categoria" | filtros, mapeamento p/ categoria ML | novo. NÃO usar como chave de tarifa ML (lição FIX-4: comissão é por categoria ML, não taxonomia ERP) |
| 7 | preco (venda ERP) | tabela de preço vigente `TGFTAB`/`TGFEXC` (qual NUTAB?) `[TESTAR-SKW]` | "Preço"/"Preço de Venda" | referência no simulador, comparação c/ preço do anúncio | novo |
| 8 | ean | `TGFBAR.CODBARRA` ou campo em TGFPRO `[TESTAR-SKW]` | "EAN"/"GTIN"/"Cód. Barras" | vínculo (ALTA), identidade catálogo ML | |
| 9 | referencia | `TGFPRO.REFERENCIA` | "Referência"/"Ref." | matcher SKU, exibição | |
| 10 | descricao | `TGFPRO.DESCRPROD` | "Descrição"/"Produto" | exibição, matcher título | |

Ausência de campo obrigatório: linha importa COM pendência visível no protocolo (badge de
completude por campo), nunca zero/default silencioso (ADR-17). `[DECISÃO D3-refinada]`: algum
dos 10 deve BLOQUEAR a linha (ex.: sem codigo_produto = rejeita)? Proposta: só #1 bloqueia.

Alimentação: import por protocolo (snapshot versionado p/ xlsx; read-through p/ Sankhya —
seção 1b), source ativa por tenant. Campos 4–7 exigem migração do modelo canônico atual
(hoje só tem 1,2,3,8,9,10 + enrichment manual de dimensões).

### E3 — Anúncio (listing ML) — "todas as informações" (RATIFICADO D-120)

Fonte primária: `GET /items?ids=...` (multiget, lote 20) / `GET /items/{id}`. Contrato
completo — cada campo com json-path; tudo persiste na tabela `listings` (expandida):

| Campo canônico | json-path ML | Status hoje | Nota |
|----------------|--------------|-------------|------|
| mlb_id | `id` | ✓ temos | chave (+ `variation_id`) |
| titulo | `title` | ✓ | |
| status | `status` | ✓ | active/paused/closed/under_review |
| sub_status | `sub_status[]` | ✗ ADICIONAR | ex.: deleted, out_of_stock, freeze |
| preco | `price`, `currency_id` | ✓ | |
| preco_original | `original_price` / `base_price` | ✗ | promoções `[TESTAR]` semântica dos dois |
| qtd_disponivel | `available_quantity` | ✓ | |
| qtd_vendida | `sold_quantity` | ✗ | vendas acumuladas do anúncio |
| qtd_inicial | `initial_quantity` | ✗ | |
| tipo_anuncio | `listing_type_id` | ✓ | gold_special (Clássico) / gold_pro (Premium) |
| categoria_ml | `category_id` | ✗ | **chave da comissão** (FIX-4) |
| condicao | `condition` | ✗ | new/used |
| permalink | `permalink` | ✗ | abrir no ML |
| foto | `thumbnail`, `pictures[]` | ✗ | thumbnail na lista; pictures no detalhe |
| video | `video_id` | ✗ | opcional |
| criado_em | `date_created` | ✗ | **habilita filtro/ordenação por data** (L1) |
| atualizado_em | `last_updated` | ✓ | |
| inicio/fim | `start_time`, `stop_time` | ✗ | opcional |
| sku_vendedor | `seller_custom_field` / attr `SELLER_SKU` | ✓ | matcher |
| ean | attributes `GTIN` | ✓ | matcher + identidade catálogo |
| marca/modelo | attributes `BRAND`, `MODEL` | ✗ | matcher lexical |
| variações | `variations[]` (id, price, available_quantity, sold_quantity, attribute_combinations, seller_custom_field, picture_ids) | parcial | hoje só id `[TESTAR]` multiget devolve variations completas? |
| envio | `shipping.mode`, `.free_shipping`, `.logistic_type`, `.tags` | ✗ | full/flex/me2; frete grátis |
| catalogo | `catalog_product_id`, `catalog_listing` | ✗ | **ponte direta p/ E6** sem re-busca por EAN |
| saude | `health` | ✗ | qualidade do anúncio |
| tags | `tags[]` | ✗ | ex.: good_quality_thumbnail, catalog_boost |
| loja_oficial | `official_store_id` | ✗ | |
| garantia | `warranty` | ✗ | opcional |

**Derivados por anúncio (enrichment no ingest, persistidos):**
| Campo | Fonte | Nota |
|-------|-------|------|
| comissao_valor + % | `GET /sites/MLB/listing_prices?price=&listing_type_id=&category_id=` | `sale_fee_amount`/`sale_fee_details` `[TESTAR]` shape por categoria |
| frete_gratis_custo | `GET /users/{seller}/shipping_options/free?item_id=` | `[TESTAR]` se algum campo do item já traz; senão fica no ingest, nunca no read |
| visitas | `GET /items/{id}/visits` ou `/items/visits?ids=` | `[TESTAR]` existe multiget de visits? cadência diária |
| preco_venda_atual | `GET /items/{id}/sale_price?context=channel_marketplace` | armadilha já conhecida (400 sem context) |

**Alimentação:**
1. *Backfill* (onboarding + botão Atualizar): `/users/{id}/items/search` com
   `search_type=scan` (catálogos >1000) → ids → **multiget `GET /items?ids=` em lotes de 20**
   com goroutines bounded (ex.: 5 workers) + backoff em 429. 2006 anúncios ≈ 101 calls
   (vs 2027 hoje). `[TESTAR]` multiget: shape da resposta, campos por item, limite de ids.
2. *Incremental*: webhook topic `items` `[DECISÃO D2]` → refetch pontual do item alterado.
3. *Reconciliação*: job agendado (ex.: 1x/dia) re-varre ids e marca closed **somente com
   varredura completa provada** (fix do L2 — truncation nunca pode fechar anúncio).

### E4 — Vínculo (anúncio ↔ produto)
Modelo e policy atuais mantidos (candidatos + aprovação humana; EAN+SKU=ALTA;
contradição vence EAN). Mudanças:
| Mudança | Motivo |
|---------|--------|
| Trigger automático de snapshot+geração após backfill e após import de catálogo | hoje endpoints órfãos (S4) |
| Policy amendment: EAN exato único → ALTA/auto-aprovar `[DECISÃO D4]` | operador já aprovou 31 em massa no D-95 |
| Anúncio sem EAN → fica REVIEW com motivo visível na UI | honestidade |

### E5 — Pedido — contrato completo (RATIFICADO D-120)

Requisito do operador: quantidade, produto vendido, anúncio vendido, comissão, frete,
comprador + dados de Nota Fiscal, valor total, **valor sobrado** (líquido). Mapeamento:

**Identidade e estado:**
| Campo | json-path / fonte | Status hoje |
|-------|-------------------|-------------|
| pedido_id | `id` | ✓ |
| pack_id | `pack_id` | ✗ ADICIONAR — N pedidos podem = 1 envio `[TESTAR]` |
| status + detalhe | `status`, `status_detail`, `tags[]` | ✓ (tags não expostas no DTO — expor) |
| canal | `context.channel` | ✗ — marketplace vs mshops (afeta tarifa) |
| datas | `date_created`, `date_closed`, `last_updated` | ✓ persistidas, **expor na lista/Fila** (P4) |

**Itens vendidos (por linha):**
| Campo | json-path / fonte | Status hoje |
|-------|-------------------|-------------|
| anuncio_vendido | `order_items[].item.id` + `.variation_id` + `.title` | ✓ |
| quantidade | `order_items[].quantity` | ✓ |
| valor_unitario | `order_items[].unit_price` + `currency_id` | ✓ (sem currency — corrigir P6) |
| produto_vendido (nosso) | via vínculo E4 (item.id → codprod) | ✓ funciona (link_reader) — manter |
| comissao_valor | `order_items[].sale_fee` | ✓ capturado `[TESTAR]` semântica: por unidade ou total da linha? multiplicar por qty? |
| custo_produto | E2.custo via vínculo, congelado na data da venda | parcial (CostReader as-of) |

**Comprador + Nota Fiscal:**
| Campo | json-path / fonte | Status hoje |
|-------|-------------------|-------------|
| comprador | `buyer.id`, `buyer.nickname` | ✓ |
| doc fiscal (CPF/CNPJ) | `GET /orders/{id}/billing-info` → doc_type, doc_number | ✓ (2-step) |
| nome/razão social | billing-info → name/business_name `[TESTAR]` shape PF vs PJ | parcial |
| endereço de entrega | shipment → `receiver_address` (rua, número, cidade, UF, CEP) | parcial (só UF/CEP) — completar p/ NF |
| inscrição estadual / IE | billing-info `[TESTAR]` vem p/ PJ? | ✗ |

**Financeiro (a "conta" do pedido — decomposição persistida no ingest):**
| Campo | Fonte | Nota |
|-------|-------|------|
| valor_total_venda | `total_amount` / `payments[].transaction_amount` | `[TESTAR]` qual é canônico com desconto/cupom |
| valor_pago | `paid_amount` / `payments[].total_paid_amount` | inclui frete pago pelo comprador |
| comissao_total | Σ sale_fee | |
| frete_vendedor | `/shipments/{id}/costs` (header `x-format-new: true`) → custo do seller (desconto de frete grátis) | `[TESTAR]` shape: senders[].cost vs receiver.cost |
| impostos_retidos | `taxes`, retenções `[TESTAR]` DIFAL/imposto na API? | senão: calculado (motor fiscal nosso) |
| cupom | `coupon.amount` | ✗ |
| **valor_sobrado (líquido)** | CALCULADO: total − comissão − frete_vendedor − impostos − custo_produto | = decomposição; motor hoje wired nil (S4) — vira entrega da base |

**Logística / SLA (Fila de verdade):**
| Campo | Fonte | Nota |
|-------|-------|------|
| shipment_id | `shipping.id` | ✓ |
| status/substatus envio | `/shipments/{id}` → status, substatus | ✓ capturado no read — mover p/ ingest |
| **SLA despacho** | shipment → `estimated_handling_limit` `[TESTAR]` | prazo p/ despachar — coluna da Fila |
| SLA entrega | `lead_time.estimated_delivery_limit` | ✓ |
| atrasado | `delayed`/tags | ✓ |
| tipo logístico | shipment `logistic_type` (full/flex/me2) | ✗ |
| rastreio | tracking_number + URL transportadora | ✓ |
| faturado_at | NOSSO (operador marca; migração 0074 em main) | our-DB-only, nunca escreve ML |

Buckets da Fila derivados no ingest e persistidos; Fila FILTRA pendentes (fix P3).

**Alimentação:**
1. *Backfill* (onboarding): `/orders/search?seller=X&sort=date_desc` paginado até horizonte
   `[DECISÃO D5: quantos meses?]`. `[TESTAR]` param de sort aceito + paginação máxima.
2. *Incremental*: webhook `orders_v2` `[DECISÃO D2]`; fallback scheduler (ex.: 5min) com
   `order.date_last_updated.from=<último sync>` — cursor persistido. Isto mata o "pedidos de
   hoje não aparecem" (P1/P2) mesmo sem webhook.
3. *Enrichment no ingest*: shipment (3 calls) + decomposição calculados quando o pedido
   entra/atualiza, persistidos. Read = SQL puro. Mata os ~10s (S3).
4. *Nossos campos* (our-DB-only, nunca escreve ML): `faturado_at` (existente em main),
   bucket derivado, vínculo de item.

### E6 — Concorrência (RATIFICADO D-120: preço, quantos, QUEM são, quantos venderam)

**Por oferta concorrente (persistir cada uma):**
| Campo | json-path / fonte | Status hoje |
|-------|-------------------|-------------|
| item_id do concorrente | `/products/{id}/items` → `results[].item_id` | ✗ **dropado hoje** (M6) — sem ele nada é rastreável |
| preco | `results[].price` + currency | ✓ |
| seller_id | `results[].seller_id` | ✓ |
| **quem é** (nome) | `GET /users/{seller_id}` (público) → nickname, reputação (`seller_reputation.level_id`, `power_seller_status`) | ✗ ADICIONAR `[TESTAR]` shape público + cachear (seller muda pouco) |
| **quantos venderam** | `sold_quantity` do anúncio concorrente | `[TESTAR]` vem no `/products/{id}/items`? senão: multiget `/items?ids=` dos concorrentes (lote 20) |
| frete grátis | `results[].shipping.free_shipping` / mode / logistic_type | parcial |
| tipo anúncio | `results[].listing_type_id` | ✗ |
| loja oficial | `results[].official_store_id` | ✗ — buy-box correlaciona |
| condição | `results[].condition` | ✓ |

**Agregados por produto (computados no ingest da coleta):**
| Campo | Regra |
|-------|-------|
| quantos_concorrentes | sellers distintos, own-seller excluído (manter; corrigir no-op silencioso M5) |
| min/mediana/max | ofertas válidas (new, BRL, ACCEPT) |
| buy_box_winner | `/products/{id}` — **um único struct DTO** (fix M6) |
| price_to_win | `/items/{id}/price_to_win?version=v2` por anúncio nosso |
| freshness | timestamp DA COLETA renderizado na UI + signal_status (fix M4) |

Alimentação: job de coleta em LOTE com fila + rate-limit (mata o 1-clique-por-produto),
disparado por: onboarding, pós-import de catálogo, agendado `[DECISÃO D6]`, botão manual.
Threshold ≥N sellers p/ "OK" vira visível/configurável, não corte silencioso (M3).

### E7 — Produto de catálogo ML + tarifas ("tudo do produto, comissão, frete, código MLB")

**Produto de catálogo ML (por produto nosso com identidade resolvida):**
| Campo | json-path / fonte | Nota |
|-------|-------------------|------|
| codigo_mlb (catálogo) | `catalog_product_id` (via anúncio E3, ou `GET /products/search?site_id=MLB&product_identifier=EAN`) | pai vs filho: ofertas só existem em LEAF (`children_ids` fan-out — regra já conhecida) |
| nome, atributos, fotos | `GET /products/{id}` → name, attributes[], pictures | exibição + conferência de match |
| categoria | `GET /products/{id}` → domain/category `[TESTAR]` path exato | chave de comissão |
| status buy-box | `buy_box_winner` | struct único (M6) |

**Tarifas (por categoria + tipo de anúncio + faixa de preço):**
| Campo | Fonte | Alimentação |
|-------|-------|-------------|
| comissao % e valor | `GET /sites/MLB/listing_prices?price=&category_id=&listing_type_id=` → `sale_fee_amount`, `sale_fee_details` | sync agendado REAL (matar seed hardcoded 16%/22% — `fee_sync.go`) + tabela histórica VENDA\|COTACAO (design pendente — memória `ml-tariff-design-pending`) |
| taxa fixa (item barato) | mesmo endpoint `[TESTAR]` — regra de custo fixo p/ preço < R$79 vem em `sale_fee_details`? R$79 hardcode BANIDO | |
| custo frete grátis | `GET /users/{seller}/shipping_options/free?item_id=` | no ingest do anúncio (E3 derivados) |

---

## 3. Decisões (status D-120)

| ID | Decisão | Status |
|----|---------|--------|
| D1 | Arquitetura de fonte de produto | **RATIFICADO**: módulo canônico + adapters convergentes (seção 1b). xlsx persiste snapshot; Sankhya read-through. Ambos nesta fase, adaptados corretamente — não é escolha de prioridade. |
| D2 | Mecanismo incremental | **RATIFICADO**: scheduler-first (pedidos ~5min via `date_last_updated`, anúncios/mercado diário). Webhook quando houver URL pública estável. |
| D3 | Campos fiscais obrigatórios? | Pendente. Recomendação: honesto-desconhecido + badge de completude. |
| D4 | Auto-aprovar vínculo EAN-exato-único | **RATIFICADO**: auto-aprovar com audit trail + tela de Importações mostrando toda a cadeia (o que entrou, o que vinculou, o que coletou). |
| D5 | Horizonte backfill de pedidos | **RATIFICADO**: 12 meses. |
| D6 | Cadência coleta de mercado | Pendente. Recomendação: diária, vinculados + watchlist. |
| D7 | Snapshot leve opcional p/ fonte-banco (histórico custo/estoque, offline) | Pendente — não bloqueia a base. |

## 4. Próximos passos da fase Base
1. Ratificar decisões D1–D6.
2. Desenhar fluxogramas: onboarding saga, sync de anúncios, sync de pedidos, import de catálogo.
3. Rodada `[TESTAR]`: provar live cada endpoint/campo marcado (multiget, scan, sort de orders, SLA no shipment, pack_id).
4. Só então: plano de implementação da base (sync engine + onboarding).
