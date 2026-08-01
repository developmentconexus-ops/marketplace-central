# Resultado Esperado — /anuncios e /pedidos (spec de campo)

```yaml
id: MIS-007-EXPECTED-RESULT
type: interface-contract
status: draft
owner: hub
parent: MIS-007
created: 2026-08-01
supersedes_scope_of: [M-05, M-06, M-07, M-08]
```

## 0. Por que este documento existe

O operador parou a missão com um diagnóstico: *"nossos contratos não estavam bem
definidos e foi sendo implementado, implementado"*. A medição confirma. A doença tem
uma forma exata e é reproduzível:

**O contrato foi escrito de um jeito que não consegue reprovar nada.**
`ListingReadModel` declara `price`, `published_quantity`, `quality_score`, `sales_30d`
como `required` **e** `nullable`. `null` é contrato-legal. Logo nenhum teste de
contrato, nenhum `tsc`, nenhum gate consegue distinguir três coisas completamente
diferentes:

1. o valor é honestamente desconhecido (ADR-17, correto);
2. o valor nunca foi implementado;
3. o valor existia e uma milestone o apagou (regressão).

Foi por esse buraco que `price_amount` sumiu de 34/34 anúncios sem nenhum gate ver.

Este documento fecha o buraco: define, **campo a campo**, o que o Mercado Livre
entrega, o que persistimos, o que a tela mostra, e — o que faltava — **qual das três
classes acima cada campo pertence**, com o critério que prova.

Ordem de leitura: §1 regras → §2 o que o ML dá de fato → §3 /anuncios → §4 /pedidos
→ §5 registro de defeitos → §6 mudanças de contrato → §7 impacto em M-05..M-08.

---

## 1. Regras (vinculantes a partir daqui)

### R-1 — Toda coluna persistida tem exatamente um dos quatro selos

| selo | significado | como o contrato expressa |
|---|---|---|
| `SEMPRE` | valor sempre existe; ausência é defeito | OpenAPI `required`, **não** `nullable` |
| `CONDICIONAL(<condição>)` | ausência é legítima só sob condição nomeada | `required` + `nullable` + **descrição obriga a nomear a condição** |
| `SEM-FONTE` | o provider não fornece; nunca vai preencher | campo **removido** do contrato, ou marcado `deprecated` com data |
| `PLANEJADO(<milestone>)` | ainda não implementado; milestone nomeada | `nullable` + descrição cita a milestone dona |

Um campo `nullable` sem uma dessas quatro classificações escrita na descrição é
contrato inválido. Isso mata a forma "ADR-17 como desculpa universal" que hoje
aparece em quase todo campo de `OrderRead`.

### R-2 — ADR-17 é para lacuna medida, não para lacuna não implementada

Texto canônico (não existia como documento; passa a existir aqui):

> **ADR-17.** Fato operacional desconhecido nunca vira zero, string vazia, `false`,
> data mínima ou "ok". Desconhecido se representa como ausência explícita (`null`) e a
> tela renderiza o glifo de desconhecido, **nomeando a razão quando ela for conhecida**.
> ADR-17 protege o *dado*; **não** autoriza declarar um campo como desconhecido porque
> ninguém o implementou. Campo não implementado é `PLANEJADO`, com milestone dona, não
> "honest-unknown".

Corolário medido: gravar `''` (string vazia) onde o provider não mandou nada **viola
ADR-17** — é fabricação de valor. Ver DEF-09.

### R-3 — Todo `EXCLUDED.x` em upsert é uma arma carregada

Uma coluna que o produtor não seta é sobrescrita com NULL a cada re-sync. Isso não é
teórico: é a causa raiz de DEF-01. Regra: **em upsert de espelho, coluna cujo produtor
pode legitimamente não ter valor usa `COALESCE(EXCLUDED.x, tabela.x)`; coluna que o
produtor sempre tem usa `EXCLUDED.x` cru.** A escolha entre os dois é decisão de
contrato e vai escrita ao lado da coluna, não improvisada por quem escreve o SQL.

Estado atual: `listings` usa `EXCLUDED.x` cru em **todas** as colunas (`repository.go:436-446`).
`orders_marketplace_orders` é misto e o misto está certo por acaso, não por decisão
registrada (`order_repo.go:668-687` usa COALESCE no bloco fiscal/buyer, cru no resto).
Itens de pedido (`order_repo.go:856-878`) são cru em tudo — mesma arma carregada do
DEF-01, ainda não disparada.

### R-4 — Prova de preenchimento é por contagem, não por presença

`fill(coluna) = não-nulo / total` medido no banco real, mais **um caso negativo**: um
registro em que a coluna deve legitimamente ser nula. Sem o caso negativo, `34/34` não
distingue "preenchido" de "preenchido com lixo constante". Sem a contagem, uma asserção
de presença não pega valor errado.

Isso substitui o padrão que falhou no M-04: `sum(jsonb_array_length(raw->'variations')) = 0`
sobre uma coluna `raw` que é NULL em 34/34 — a soma dá zero sem discriminar nada. Ver DEF-02.

### R-5 — Nenhuma milestone nova sem esta spec como âncora

Os briefs de M-05..M-08 foram escritos contra o contrato frouxo. Precisam ser
re-ancorados nesta tabela antes de despacho. §7 diz o que muda em cada um.

---

## 2. O que o Mercado Livre entrega de fato

Fonte: respostas reais capturadas em `docs/design/evidence/ml-api/`. Não é documentação
decorada — é o corpo que a conta do operador devolveu.

### 2.1 `GET /items?ids=…` (multiget) — base de /anuncios

Chaves do item, verificadas presentes: `id, title, status, sub_status, price, base_price,
original_price, currency_id, listing_type_id, available_quantity, initial_quantity,
sold_quantity, category_id, condition, permalink, thumbnail, date_created, last_updated,
tags, catalog_product_id, catalog_listing, seller_custom_field, attributes[],
variations[], shipping{mode, logistic_type, free_shipping, local_pick_up, store_pick_up,
methods, dimensions, tags}, health, seller_address, official_store_id, warranty,
sale_terms, pictures, domain_id, family_id, user_product_id, inventory_id, …`

Dois fatos decisivos:

- **`price`, `currency_id`, `listing_type_id`, `initial_quantity` existem.** Não há
  limitação de API nenhuma. O que falta é nosso DTO declarar.
- **`health` vem `null` em 100% dos itens capturados.** Nossa coluna `quality_score`
  promete um número que a API não entrega para esta conta.

### 2.2 Endpoints que **não** existem para o que prometemos

| o que o contrato promete | realidade medida |
|---|---|
| `sales_30d` (vendas 30 dias) | nenhum endpoint capturado fornece. `sold_quantity` é acumulado desde sempre, não janela. |
| `quality_score` | `health` = `null` sempre. |
| multiget de shipments | `GET /shipments?ids=…` → **404**. Só existe leitura 1-a-1. |
| `/orders/{id}/discounts` | **404** `discount_not_found` na conta. |
| item de terceiro | **403** sempre (`access_denied` ou PolicyAgent). Preço de concorrente só via `/products/{id}/items`. |

### 2.3 Pedidos e envios — onde cada valor mora

| valor | endpoint | caminho |
|---|---|---|
| comissão | `GET /orders/{id}` | `order_items[].sale_fee` — **por unidade**, multiplicar por `quantity` |
| valores do pedido | `GET /orders/{id}` | `total_amount`, `paid_amount`, `shipping_cost`, `currency_id` |
| pagamentos | `GET /orders/{id}` | `payments[]` (inclui `marketplace_fee`, `taxes_amount`, `coupon_amount`) |
| status do envio | `GET /shipments/{id}` | `status`, `substatus` |
| modalidade logística | `GET /shipments/{id}` | **`logistic.type`** (aninhado — não é chave plana) |
| rastreio | `GET /shipments/{id}` | `tracking_number`, `tracking_method` |
| SLA — prazo | `GET /shipments/{id}` | `lead_time.estimated_delivery_limit.date` — é daqui que `sla_limit_at` vem hoje |
| SLA — situação | `GET /shipments/{id}/sla` | `expected_date`, `status` (vocabulário próprio: `on_time`…) — **endpoint nunca chamado hoje** |
| custo de frete | `GET /shipments/{id}/costs` | `gross_amount`, `receiver.cost`, `senders[].cost` |
| tarifa por preço | `GET /sites/MLB/listing_prices?price=X` | `sale_fee_amount`, `sale_fee_details{percentage_fee, fixed_fee, gross_amount}`, `listing_fee_amount` |
| cursor incremental | `GET /orders/search?date_last_updated.from=…` | `last_updated` do pedido |

O vocabulário de `status` do envio e o de `/sla` **são diferentes**. Não podem cair na
mesma coluna. Hoje `order_shipments` tem `status` e `sla_status` separados — a separação
está certa; o preenchimento não (DEF-08).

---

## 3. /anuncios — tabela de verdade

Estado medido: 34 anúncios, `listings` no dev stack, pós-sweep do M-04.

| campo na tela | rótulo | payload | coluna | fill | fonte ML | selo correto | veredito |
|---|---|---|---|---|---|---|---|
| MLB | `MLB` | `provider_listing_id` | idem | 34/34 | `id` | SEMPRE | **OK** |
| título | `TÍTULO` | `title` | idem | 34/34 | `title` | SEMPRE | **OK** |
| preço | `PREÇO` | `price.{amount,currency}` | `price_amount`/`price_currency` | **0/34** | `price`/`currency_id` | SEMPRE | **DEF-01 REGRESSÃO** |
| estoque | `EST.` | `published_quantity` | idem | **0/34** | `initial_quantity` | SEMPRE | **DEF-01 REGRESSÃO** |
| modalidade | (detalhe) | `listing_type.{code,label}` | `listing_type_code` | **0/34** | `listing_type_id` | SEMPRE | **DEF-01 REGRESSÃO** |
| qualidade | `QUAL.` | `quality_score` | idem | 0/34 | `health` = null | **SEM-FONTE** | **DEF-03 remover do contrato** |
| vendas 30d | (não renderizado) | `sales_30d` | idem | 0/34 | nenhuma | **SEM-FONTE** | **DEF-04 campo morto** |
| status | (abas/pills) | `status` | idem | 34/34 | `status` verbatim | SEMPRE | **OK, mas ver DEF-05** |
| sync | `SYNC` | `sync_state` | idem | 34/34 | computado | SEMPRE | **OK** |
| erro de sync | `PENDÊNCIA` | `sync_error` | idem | 0/34 | computado | CONDICIONAL(sem erro) | **OK** |
| vínculo | `PRODUTO` | `link.{state,product_id,seller_sku}` | tabela de vínculo | — | interno | SEMPRE | **OK** |
| custo | (detalhe) | `cost` | ERP | — | ERP, não ML | CONDICIONAL(sem custo no ERP) | **OK** |
| margem | `Margem est.` | `below_margin_worst_case` | derivado | — | derivado | CONDICIONAL(sem custo) | **OK** |
| estoque ERP (grupo) | `ERP est. —` | **nenhum** | — | — | — | — | **DEF-06 literal hardcoded** |
| sinal de mercado | `Vs. mercado` | `market_signal.*` | outra tabela | — | `/products/{id}/items` | CONDICIONAL(sem vínculo/sem evidência) | **OK** |
| tarifa | **não existe** | — | `channel_fees` c2 | — | `/sites/MLB/listing_prices` | PLANEJADO(M-05) | pendente |
| divergência estoque | **não existe** | — | `divergences` | — | ML × ERP | PLANEJADO(M-05) | pendente |

Resumo: das 4 colunas visíveis que o operador viu em branco, **3 são regressão de uma
causa só** e **1 não tem fonte no ML**.

---

## 4. /pedidos — tabela de verdade

Estado medido: 38 pedidos (31 `paid`, 7 `cancelled`), 38 envios.

### 4.1 Cabeçalho e identidade

| campo | payload | coluna | fill | selo | veredito |
|---|---|---|---|---|---|
| nº do pedido | `provider_order_id` | idem | 38/38 | SEMPRE | OK |
| data | `provider_created_at` | idem | 38/38 | SEMPRE | OK |
| status do provider | `status` | `provider_status` | 38/38 | SEMPRE | OK |
| detalhe do status | `provider_status_detail` | idem | 38/38 não-nulo · **0/38 não-vazio** | CONDICIONAL | **DEF-09 `''` no lugar de desconhecido** |
| cancelamento | — | `cancellation_detail` | 38/38 não-nulo · **0/38 não-vazio** | CONDICIONAL(não cancelado) | **DEF-09** — e 7 pedidos ESTÃO cancelados |
| bucket | `bucket` | idem | 38/38 | SEMPRE | OK (derivado) |
| faturado | `faturado_at` | idem | 1/38 | CONDICIONAL(não faturado) | OK — fato nosso, não do ML |
| cursor incremental | — | `date_last_updated_ml` | **0/38** | SEMPRE | **DEF-07 — sync incremental não tem chave** |

### 4.2 Envio

| campo | coluna | fill | fonte ML | selo | veredito |
|---|---|---|---|---|---|
| status do envio | `status` | 38/38 | `/shipments/{id}.status` | SEMPRE | OK |
| substatus | `substatus` | **1/38** | `.substatus` | CONDICIONAL | OK — parseado (`:31,:144`); 1/38 é o dado, não o código |
| modalidade logística | `logistic_type` | **0/38** | `.logistic.type` | SEMPRE | **DEF-08 — chave nunca declarada no struct** |
| método de rastreio | `tracking_method` | **0/38** | `.tracking_method` | CONDICIONAL | **DEF-08 — idem** |
| código de rastreio | `tracking_number` | 35/38 | `.tracking_number` | CONDICIONAL(não despachado) | OK |
| SLA — prazo | `sla_limit_at` | 38/38 | `/sla.expected_date` | SEMPRE | OK |
| SLA — situação | `sla_status` | **0/38** | `/sla.status` | SEMPRE | **DEF-08 — o endpoint `/sla` nunca é chamado** |
| custo bruto/vendedor | `cost_gross`/`cost_seller` | 38/38 | `/costs` | SEMPRE | OK |
| destino cidade/UF | `dest_city`/`dest_state` | 38/38 | `.destination` | SEMPRE | OK |

### 4.3 Financeiro (drawer — hoje quase tudo `—`)

| campo | payload | selo correto | dono |
|---|---|---|---|
| comissão | `decomposicao.comissao` | PLANEJADO | **M-06** — fonte existe (`sale_fee` × qty) |
| taxa fixa | `decomposicao.taxa_fixa` | PLANEJADO | **M-06** — fonte: `sale_fee_details.fixed_fee` |
| frete | `decomposicao.frete` | PLANEJADO | **M-06** — fonte já persistida (`cost_seller`) |
| imposto | `decomposicao.imposto` | PLANEJADO | **M-06** — fonte interna, não ML |
| DIFAL | `decomposicao.difal` / `difal.*` | PLANEJADO | **M-06** — fonte interna |
| tarifa Full | `decomposicao.tarifa_full` | PLANEJADO | **M-06** |
| custo do produto | `decomposicao.custo` | CONDICIONAL(sem vínculo/sem custo ERP) | M-06 |
| margem / retorno | `margem_valor`/`margem_pct`/`retorno_liquido` | PLANEJADO | **M-06** — colunas `net_amount`/`margin_pct`/`decomposition` 0/38 **por desenho** |

Aqui o comportamento atual está **certo**: `—` com a razão nomeada. O erro é o contrato
chamar isso de honest-unknown em vez de `PLANEJADO(M-06)`.

### 4.4 Comprador e fiscal

| campo | selo | veredito |
|---|---|---|
| apelido do comprador | SEMPRE | OK (38/38) |
| nome/documento/endereço fiscal | CONDICIONAL(**só na rota de detalhe**) | **OK mas mal declarado** — a lista nunca resolve (decisão de performance, FINDING-M08-LIST-TIMEOUT). O contrato precisa dizer isso; hoje parece lacuna de dado. |

### 4.5 KPIs da tela

| KPI | origem | veredito |
|---|---|---|
| Novos / A faturar / A enviar / Enviados | contados **no cliente sobre a página carregada** | **DEF-10** — `GET /orders/summary` e o tipo `OrderSummary` existem e a tela não usa. Correto por acaso com 38 pedidos; mente na primeira página cheia. |
| DIFAL a pagar | **constante `UNKNOWN_KPI_VALUE` no código** | **DEF-11** — nunca lê payload nenhum |

---

## 5. Registro de defeitos

| id | defeito | classe | causa raiz (file:line) | dono |
|---|---|---|---|---|
| **DEF-01** | `price_amount`, `price_currency`, `published_quantity`, `listing_type_code` NULL em 34/34 | **REGRESSÃO (M-04)** | `items_multiget_reader.go:136-161` não declara `price`/`currency_id`/`listing_type_id`/`initial_quantity` → `multiget_mapper.go:71-148` não seta → `repository.go:436-446` grava NULL por cima | chip novo |
| **DEF-02** | M04-C3 aprovado com prova vazia | **PASS VAZIO** | `listings.raw` 0/34; a soma sobre NULL dá 0 | hub — reverter veredito |
| **DEF-03** | `quality_score` prometido sem fonte | **SEM-FONTE** | `health` = null em 100% | contrato |
| **DEF-04** | `sales_30d` no contrato e no SDK, nenhuma tela renderiza, nenhum produtor | **CAMPO MORTO** | — | contrato |
| **DEF-05** | dump do ML diz `active=9 / paused=18`; banco diz `active=0 / paused=27 / under_review=7` | **PROVÁVEL DUMP VELHO** — não bloqueante | Evidência contra "estamos mangling": `status` é gravado verbatim sem remap (`multiget_mapper.go:106`) e o banco contém `under_review`, valor que **só pode ter vindo do ML**. Se houvesse normalização quebrada, `under_review` não sobreviveria. Dump F1 é de D-120; sweep é de 2026-08-01. Confirmar com uma contagem viva antes de fechar. | medição |
| **DEF-06** | `ERP est. —` literal fixo no cabeçalho de grupo | **NUNCA IMPLEMENTADO** | `AnunciosTable.tsx:223-266`; `ListingGroup` não tem o campo | M-05 |
| **DEF-07** | `date_last_updated_ml` 0/38 | **NUNCA IMPLEMENTADO** | é a chave do incremental de pedidos | **M-06 — bloqueante** |
| **DEF-08** | `logistic_type`, `tracking_method`, `sla_status` 0/38 | **NUNCA IMPLEMENTADO (confirmado)** | `shipment_ingest_reader.go:21-41` — `mlIngestShipmentResponse` não declara `logistic`/`logistic_type`/`tracking_method`; `mapShipmentDetail:140-177` nunca atribui os três; `GET /shipments/{id}/sla` **nunca é chamado** (omissão deliberada documentada em `:53-58`). Plumbing a jusante está correto (`*string` + coluna nullable). | chip |
| **DEF-09** | `provider_status_detail` e `cancellation_detail` gravados como `''` | **VIOLAÇÃO ADR-17 (confirmado)** | `order_detail.go:19,38` — os campos são `string`, não `*string`, contra a própria convenção declarada no cabeçalho do arquivo (`:13-14`). Chave ML ausente → `""`; NULL é **inalcançável** para essas colunas. | chip |
| **DEF-16** | `cancellation_detail` copia a MESMA fonte de `provider_status_detail` | **LÓGICA ERRADA** | `order_ingest_reader.go:180` e `:186` atribuem os dois a partir de `order.StatusDetail`. Motivo de cancelamento nunca terá valor próprio — e há 7 pedidos cancelados. | chip |
| **DEF-10** | KPIs de /pedidos contados no cliente sobre página paginada | **LÓGICA ERRADA** | `PedidosPage.tsx:41-64,176-177`; `OrderSummary` existe e não é usado | chip |
| **DEF-11** | KPI "DIFAL a pagar" é constante de desconhecido | **NUNCA IMPLEMENTADO** | `PedidosPage.tsx:292-297` | M-06 |
| **DEF-12** | `listings.raw` 0/34 — payload do ML não é guardado | **LACUNA DE EVIDÊNCIA** | impede auditar espelho contra origem | chip |
| **DEF-13** | upsert de itens de pedido sem COALESCE | **ARMA CARREGADA** | `order_repo.go:856-878` — mesma forma do DEF-01, ainda não disparada | R-3 |
| **DEF-14** | contrato `required`+`nullable` sem classe | **CONTRATO INVÁLIDO** | OpenAPI listings + orders | §6 |
| **DEF-15** | ADR-17 aplicado repo inteiro sem texto canônico | **DOUTRINA FANTASMA** | só uma linha em `HARNESS-PROFILE.md:246` | R-2 resolve |

---

## 6. Mudanças de contrato exigidas

1. **Remover** `sales_30d` de `ListingReadModel`, do SDK e da coluna (DEF-04). Campo sem
   fonte, sem produtor e sem leitor não é contrato: é dívida disfarçada.
2. **`quality_score` → `deprecated`** com a razão escrita (`health` nulo na API). Se um dia
   houver fonte, volta com produtor junto — nunca antes.
3. **`price`, `published_quantity`, `listing_type`, `status` → `required`, NÃO `nullable`**.
   Depois do DEF-01 corrigido, `null` nesses campos passa a ser reprovação automática de
   contrato, e a regressão que aconteceu fica impossível de repetir em silêncio.
4. **Toda propriedade `nullable` restante ganha classe na descrição** (`CONDICIONAL(razão)`
   ou `PLANEJADO(M-0X)`). Sem classe = contrato inválido (R-1).
5. **`OrderRead.comprador_fiscal`** ganha na descrição: resolvido **só** na rota de detalhe,
   por decisão de performance — não é lacuna de dado.
6. Renomear `ListOrdersResponse` **ou** `OrderPage` para casarem (drift só de nome, mas é
   drift).

---

## 7. Impacto em M-05 … M-08

| milestone | como estava | o que muda |
|---|---|---|
| **M-05** listings-fees-divergence | assume `listings` íntegro | **bloqueado por DEF-01.** Divergência de estoque compara ML × ERP; hoje o lado ML (`published_quantity`) é NULL em 34/34 — a comparação seria contra nada. Absorve DEF-06. |
| **M-06** orders-backfill-decomposition | backfill 12m + incremental 5min + decomposição | **bloqueado por DEF-07**: incremental usa `date_last_updated_ml`, 0/38. Absorve DEF-11. Escopo confirmado correto: as fontes de comissão/frete/taxa fixa existem e estão mapeadas em §2.3. |
| **M-07** pricing-fee-read | lê tarifa do ledger | inalterado — não depende de nada acima. Pode ir na frente. |
| **M-08** webhook-ingest | pedido novo em segundos | inalterado, mas **depende do M-06** (chama o `IngestOrder`). Continua por último. |

**Chip de correção antes da onda C.** Todos os defeitos abaixo têm a mesma assinatura —
**o DTO do adapter não declara a chave que o ML manda** — e por isso vão num chip só, com
um guard único que os impede de voltar:

| defeito | superfície |
|---|---|
| DEF-01 | `items_multiget_reader` — `price`, `currency_id`, `listing_type_id`, `initial_quantity` |
| DEF-08 | `shipment_ingest_reader` — `logistic.type`, `tracking_method`, + chamar `/sla` |
| DEF-09 | `order_detail` — `string` → `*string` nos dois campos |
| DEF-16 | `order_ingest_reader:186` — fonte própria para cancelamento |
| DEF-12 | persistir `listings.raw` |
| DEF-13 | COALESCE no upsert de itens (R-3) |

DEF-10 (KPI de /pedidos) é FE isolado e pode ir junto. **DEF-05 é medição, não código:
resolver antes de escrever o chip** — se estivermos gravando `paused` sobre `active`, o
escopo do chip muda.

---

## 8. Como cada campo se prova (mata o pass vazio)

Para qualquer campo desta spec, a evidência aceitável é **uma** das duas:

- **Fill medido com caso negativo** (R-4): `não-nulo/total` no banco real **mais** um
  registro em que a ausência é legítima e o campo está de fato nulo. Sem o caso negativo,
  não é prova.
- **Must-fail nomeado**: injetar o defeito, mostrar o teste vermelho **citando o nome do
  teste**, corrigir, mostrar verde. `status=passed` sem `failure_token=test=` é
  byte-idêntico a "nenhum teste rodou".

Prova proibida: asserção de presença de chave, contagem sobre coluna que pode ser NULL, e
qualquer agregação cujo zero seja indistinguível de instrumento morto — foi exatamente
assim que o DEF-02 passou.

---

## §9 — Emendas pós-estudo da API (2026-08-01)

Fonte: `research/ml-api-capability-map.md`. Duas lanes de pesquisa (catálogo/anúncio e
vendas/envio/financeiro) contra doc oficial via context7 + os dumps reais. Onde este §9
contradiz §5, **§9 ganha** — §5 foi escrito antes do estudo.

| DEF | veredito anterior | veredito corrigido |
|---|---|---|
| DEF-03 `quality_score` | SEM-FONTE | **ERRADO. Endpoint errado.** Fonte real = `GET /item/{id}/performance` (`score`, `level`, `buckets[].variables[].rules[]` com ação corretiva). `health` no item vem `null` porque está sendo **descontinuado** pelo ML, não porque não existe qualidade. |
| DEF-04 `sales_30d` | campo morto | Confirmado **sem endpoint**. `sold_quantity` é vitalício. Só derivando de `/orders/search` por item+data ⇒ é campo **DERIVADO**, e o contrato tem que dizer isso. |
| DEF-05 status 34 vs 27 | discrepância aberta | **RESOLVIDO pelos próprios dumps**: `all=34, active=9, paused=18, closed=0`. Os 7 restantes são `under_review`. Conta mudou entre a captura (julho) e agora. Não é bug de gravação. |
| DEF-16 `cancellation_detail` | duplica `status_detail` | Confirmado defeito, e **existe fonte real distinta**: objeto `cancel_detail{group, code, description, requested_by, date}` no pedido. Ler `status_detail` ali é perda de dado, não só duplicação. |

### Achados novos que mudam milestone

**N-1 (M-05, M-07) — `listing_prices` sem `logistic_type` calcula tarifa errada.**
A doc chama o parâmetro de **crucial** para o `fixed_fee` correto. Medido: na categoria
testada o `fixed_fee` cai de 6,24 para 0 entre R$ 12,50 e R$ 13,00. Simulador que não passa
`logistic_type` mente. Vira critério de aceite do M-05/M-07.

**N-2 (margem, todos) — `payments[].marketplace_fee` NÃO é a comissão.**
No pedido real medido: `marketplace_fee = 0` com `sale_fee = 120,43` no mesmo pedido.
Comissão confiável = **`order_items[].sale_fee × quantity`** (`sale_fee` é **por unidade** —
provado por aritmética: 120,43 / 729,90 = 16,5 %; se fosse total da linha daria 8,25 %, fora
de qualquer faixa do ML). A doc é silente: todos os exemplos usam `quantity=1`.

**N-3 (margem, M-07) — `lead_time.cost` é custo do COMPRADOR.**
Custo do vendedor = `/shipments/{id}/costs` → `senders[].cost`, já líquido de desconto. No
envio real: comprador 0, vendedor **19,85**. Ler `lead_time.cost=0` como "frete de graça"
erra a margem em R$ 19,85 por pedido.

**N-4 (M-06) — backfill de 12 meses é risco aberto.**
`order.date_last_updated.from` em `/orders/search` está **confirmado** como o mecanismo
incremental correto. Mas **scan/scroll não é documentado para `/orders/search`** (só para
`/items/search` e `/questions/search`) e o teto de offset é `NÃO VERIFICADO`. Acima de ~1000
pedidos o backfill pode não ter caminho. **Medir antes de planejar o M-06.**

**N-5 (M-08) — política de webhook do brief está errada.**
Real: **8 tentativas em janela de 1 h** (brief diz 5), receptor precisa responder **200 em
≤ 500 ms**, `missed_feeds` guarda só **2 dias**, e falha repetida **desativa o tópico** para
a aplicação. Desativação silenciosa de tópico é perda de sincronismo sem alarme — precisa de
detecção própria.

**N-6 (novo, sem dono) — vocabulário é ABERTO em toda parte.**
`status_detail` de pedido, `tags`/`static_tags`, `payments[].status`, substatus de envio
(100+, sem tabela canônica na doc; o nosso `invoice_pending` não consta de nenhuma lista
oficial), `logistic.mode/type/direction`. Confirma ADR-06/IC-07: **verbatim sempre, enum
nunca**. Qualquer CHECK ou `switch` exaustivo sobre esses campos quebra em produção.

**N-7 (novo, sem dono) — domínios inteiros que o produto não tem.**
- Competição de catálogo: `GET /items/{id}/price_to_win?version=v2` com `status`
  winning/losing, `price_to_win`, `winner`, e **`boosts[]`** (frete grátis, full,
  parcelamento) — a resposta para "por que perco com preço igual".
- Pendência/moderação: `GET /moderations/infractions/{user_id}` → `reason` + `remedy`.
  É o que deveria alimentar a coluna PENDÊNCIA que hoje está vazia.
- Pós-venda: `GET /post-purchase/v1/claims/search` — reclamação, devolução, mediação,
  `affects_reputation`, `available_actions`. Nossa tela **não tem o conceito**.
- Linha do tempo de envio: `/shipments/{id}/history`. Nota fiscal do envio:
  `/shipments/{id}/invoice_data`.

### Custo de sincronização — fato de desenho

Não existe multiget de shipments (404 medido). Visitas são 1 item por chamada. Performance,
price_to_win e infrações são 1 item por chamada. **O custo é dominado por sub-recurso 1-a-1,
não pela listagem.** Qualquer plano que trate "buscar pedidos" como uma chamada erra por uma
ordem de grandeza.

---

## §10 — Correções por MEDIÇÃO ao vivo (2026-08-01)

Fonte: `research/live-probe-results.md`, sonda `cmd/mlprobe -round3` contra a conta real.
**§10 ganha de §9 e de §5.** §9 veio de documentação; §10 veio do payload. Onde a doc
contradiz a medição, a medição manda.

### Refutado

**N-1 CAI (parcialmente).** `logistic_type` **não alterou nada** em `listing_prices`:
12 combinações (`drop_off` · `xd_drop_off` · `fulfillment` · ausente) na categoria MLB270310
devolveram `sale_fee_details` **idêntico**. A doc chama o parâmetro de crucial; a conta diz
que não é — ao menos nessa categoria. **Não vira critério de aceite do M-05/M-07** sem
medição em segunda categoria. Passar o parâmetro segue correto por higiene; *depender* dele
não se justifica com a evidência atual.

### Resolvido

**N-4 FECHADO — M-06 tem caminho.** `/orders/search`: teto de offset = **10000**
(`limit.maximum_exceeded` acima disso), e **`search_type=scan` é aceito** — HTTP 200 com
`scroll_id` real — apesar de a doc só documentar scan para `/items/search` e
`/questions/search`. Backfill de 12 meses acima de 10k pedidos é viável por scan.
Risco de milestone fechado.

**DEF-16 FECHADO, e é pior que o diagnóstico.** Nos 7 pedidos cancelados: `cancel_detail`
presente em **7/7** com `requested_by` ∈ {buyer, meli, seller}; e `status_detail` **null em
todos**. Não é dado duplicado — é **dado zero**. `cancellation_detail` grava string vazia
sempre (o tipo `string`, não `*string`, transforma o null em `''`). Fonte correta =
`cancel_detail.description` + `.requested_by` + `.group` + `.code`.

**N-3 CONFIRMADO 5/5.** `senders[0].cost` entre **23,65 e 138,95**; `lead_time.cost` = **0**
em todos os cinco. Margem que lê `lead_time` perde o frete inteiro.

**DIFAL — base de UF resolvida na prática.** UF do envio == UF do `billing_info` em **38/38**
pedidos. A ambiguidade documental persiste, mas não é bloqueante operacionalmente. Bônus:
condição de contribuinte de ICMS é detectável em `billing_info.taxes.taxpayer_type` +
`inscriptions.state_registration` — insumo direto do motor do ADR-C5.

**PENDÊNCIA tem fonte confirmada.** Os 7 anúncios `under_review` aparecem entre as **55
infrações** retornadas por `/moderations/infractions/{user_id}`, com `reason` e `remedy`
reais. A coluna vazia da tela tem de onde ser preenchida.

**Multiget: o limite real é 20.** Não 50. A doc se contradizia entre duas páginas; medido,
com 400 claro acima disso. O servidor deduplica ids **antes** de contar.

**`/shipments/{id}/sla` existe** e devolve `{expected_date, status, service, last_updated}`
(amostra de 1). Hoje nunca é chamado; o `SLALimitAt` atual vem de `lead_time`, que é outra
grandeza.

### Medido contra estado transitório — REFAZER

**M3 e M5 da primeira rodada não valem.** A conta estava em **modo férias** (anúncios
desativados pelo operador), o que produziu zero itens `active` — e daí
`/item/{id}/performance` → 400 *"Only status active is supported"* e `price_to_win` →
`item_not_opted_in`. Anúncios reativados; medição refeita. **Não tirar conclusão de produto
da primeira rodada.**

Lição de método, e não é a primeira vez nesta missão: *estado transitório da conta
disfarçado de capacidade ausente da API*. Antes de declarar "a API não dá", confirmar o
estado da conta no momento da medição.

### Dívida de segurança aberta por esta rodada

`sanitize()` do `cmd/mlprobe` não cobre as chaves `identification` e `state_tax_id` ⇒
`docs/design/evidence/ml-api/T0-identity.json` guardava **CNPJ e inscrição estadual do
vendedor em claro**. Arquivo untracked, nunca commitado. A própria lane também gravou IE de
comprador na primeira passada, detectou e corrigiu redigindo **na captura**. Correção do
`piiKeys` + remoção do arquivo despachadas. O diretório `docs/design/evidence/ml-api/`
continua pendente de scrub completo antes de qualquer compartilhamento.

### §10.1 — M3/M5 remedidos com anúncios ativos

Conta reativada pelo operador. Confirmado **ao vivo** (via `GET /items?ids=`, nunca pelo
Postgres, que estava velho): **9 de 34** listings `active`; `sub_status[]` = `[]` em 9/9.

**M3 — MEDIDO. Qualidade tem fonte e tem dado.**
`GET /item/MLB4834219830/performance` → 200:
`score=49` · `level="bad"` · `level_wording="Básica"` · `calculated_at` · `buckets[]` (2:
`USER_PRODUCT` e `ITEM`) → `variables[].rules[]` com `status` (COMPLETED/PENDING), `progress`
(0–1, fracionário), `mode` (`OPPORTUNITY` em toda a amostra) e `wordings{label, link, title}`.
O `link` aponta a ação de correção. Amostra **n=1** — é o único ativo não-catálogo da conta
hoje. Reconfirmar quando houver mais.

**M5 — MEDIDO, e o veredito é de produto, não de API.**
8 dos 10 listings com `catalog_product_id` confirmados `active`. `price_to_win?version=v2`
nos 8: **idêntico à rodada com a conta em férias** — `status=not_listed`,
`reason=item_not_opted_in`, `boosts=null`.

Reativar o anúncio **não basta**: existe uma etapa separada de opt-in/inscrição em catálogo
do lado do ML, independente do status do anúncio. Provado por controle — mesmo resultado com
o item pausado e com o item ativo, o que elimina o status como causa.

`boosts[]` (`boosted` vs `opportunity`) segue **NÃO MEDIDO** e seguirá até algum item atingir
`status=listed`.

### §10.2 — Escopo do bloco B, corrigido pela medição

| metade | observável hoje | bloqueio |
|---|---|---|
| margem (comissão + frete do vendedor + custo ERP) | **sim** — 38 pedidos reais | nenhum |
| qualidade (`/performance` + infrações) | **sim** — payload real medido | nenhum |
| competição (`price_to_win` + `boosts[]`) | **não** | opt-in de catálogo no ML — ação do operador, externa ao código |

Margem e qualidade entram no escopo com drive vivo e QA capaz de reprovar. Competição é
`PLANEJADO(M-xx)`, **com o bloqueio nomeado**: não é falta de fonte (a fonte existe e
responde 200), é falta de inscrição da conta. Destrava remedindo M5 depois do opt-in.
