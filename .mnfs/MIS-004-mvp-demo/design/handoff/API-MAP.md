# Mapa de viabilidade — UI × API Mercado Livre (+ Mercado Pago + ERP)

Análise feita em jul/2026 contra a documentação oficial (developers.mercadolivre.com.br e mercadopago.com.br/developers). Cada bloco de cada tela está classificado:

- ✅ **API direta** — endpoint oficial cobre o dado/ação
- ⚠️ **Derivado** — a API dá a matéria-prima; o valor exibido é calculado/agregado por nós (ou exige armazenar snapshots)
- ❌ **Não vem do ML** — vem do ERP ou é regra/feature nossa (banco próprio)

Regras gerais da plataforma: OAuth 2.0 authorization code + refresh; rate limit ~1500 req/min por seller; webhooks (topics: `orders_v2`, `shipments`, `claims` [marketplace claims], `items`, `payments`, `item competition`, `fbm stock operations`) — usar webhooks como fonte de eventos e GET como reconciliação.

---

## Pedidos.dc.html

| Bloco | Fonte | Status |
|---|---|---|
| Lista/fila de pedidos, status, valores, itens | `GET /orders/search?seller=$ID` + webhook `orders_v2` | ✅ |
| SLA de envio (data-limite, atrasado) | `GET /shipments/$ID` (`handling_time`, `estimated_delivery_*`) e `GET /shipments/$ID/delays` | ✅ |
| Etiqueta | `GET /shipment_labels?shipment_ids=` | ✅ |
| Rastreio | `GET /shipments/$ID` (`tracking_number`, `status`, histórico) + webhook `shipments` | ✅ |
| NF-e (emitir/consultar) | API de invoicing MLB (`/users/$ID/invoices/...`); emissão pode ser do ERP — decidir na implementação | ✅/ERP |
| Comprador (dados fiscais) | `buyer.billing_info.id` no order → `GET /orders/billing-info/$SITE/$BILLING_INFO_ID` (fluxo novo; legado deprecado) | ✅ |
| Timeline ML → interno | composição: eventos de order + shipment + nossos (faturamento ERP) | ⚠️ |
| Retorno real por pedido (decomposição) | tarifas: `GET /orders/$ID` (`sale_fee` por item) + `GET /shipments/$ID/costs` (frete do vendedor c/ descontos); custo: ERP | ⚠️ |
| Coluna/bloco DIFAL, agendamento, lembrete, "marcar pago" | ❌ ML não calcula DIFAL. UF destino vem do shipment (`receiver_address.state`) ✅; alíquota vem da nossa tabela (Configurações); agendamento/lembrete/pago = banco próprio | ❌ (feature nossa) |
| Kanban "status vem do ML" | mesmos orders/shipments; correto não permitir arrastar | ✅ |

## Pedidos — abas Cancelados / Devoluções

| Bloco | Fonte | Status |
|---|---|---|
| Reclamações/devoluções, motivo | `GET /post-purchase/v1/claims/search` (filtros stage/status/reason_id) + detalhe `/claims/$ID`; devolução: `GET /marketplace/v2/claims/$CLAIM_ID/returns` (tipos claim/dispute/automatic) | ✅ |
| Impacto na reputação (chip "rep ↓") | `GET /marketplace/v2/claims/$ID/affects-reputation` (inclui `has_incentive`: responder em 48h evita impacto) | ✅ |
| Status do reembolso | payments do order (webhook `payments`) + status da claim | ✅ |
| Reversa (retorno ao estoque) | return c/ tracking próprio; conferência: `POST /post-purchase/v1/returns/$RETURN_ID/return-review` (OK ou fail c/ evidência via `/claims/$ID/returns/attachments`) | ✅ |
| Ação "Contestar" | `open_dispute` / mensagens `POST /marketplace/v2/claims/$ID/actions/send-message`; reembolso parcial: `/claims/$ID/partial-refund/available-offers` | ✅ |
| Motivo de cancelamento simples (sem claim) | `GET /orders/$ID` (`cancel_detail`) | ✅ |
| Reintegração ao estoque físico | ERP | ❌ (ERP) |

## Anuncios.dc.html

| Bloco | Fonte | Status |
|---|---|---|
| Lista de anúncios, preço, estoque, modalidade | `GET /users/$ID/items/search` + `GET /items?ids=` (multiget) + webhook `items` | ✅ |
| Alterar preço/estoque (fila de sync) | `PUT /items/$ID` — a *fila com preview e protocolo* é arquitetura nossa | ✅ (ação) / ❌ (fila) |
| Comissão real por anúncio | `GET /items/$ID/sale_price` + `GET /sites/MLB/listing_prices?price=` | ✅ |
| Exceções/qualidade | Listing quality API (`/item/$ID/performance` etc.) | ✅ |
| vs. mercado (chip %) | derivado da busca (ver Mercado) | ⚠️ |

## Produto Detalhe.dc.html

| Bloco | Fonte | Status |
|---|---|---|
| Cadastro (GTIN, NCM, custo, dimensões) | ERP | ❌ (ERP) |
| Aba Estoque: físico/reservado | ERP (reservado cruza c/ orders abertos ✅) | ❌/⚠️ |
| Aba Estoque: Full separado | `GET /inventories/$INVENTORY_ID/stock/fulfillment` + operações `GET /stock/fulfillment/operations/search` + webhook `fbm stock operations`. Atenção: inbound ao Full só pelo Seller Center (API só consulta) | ✅ (consulta) |
| Aba Estoque: movimentações | ERP + fulfillment operations | ⚠️ |
| Alerta de ruptura (dias de cobertura) | cálculo nosso (vendas 30d ÷ estoque) | ⚠️ |
| Aba Pedidos: lista + mini-métricas | agregação de `/orders/search` por item | ⚠️ |
| Aba Histórico: preço nosso vs. mercado 90d | ❌ ML não dá série histórica de preços — armazenar snapshot diário da busca no nosso banco desde o dia 1 | ❌ (snapshot próprio) |
| Aba Histórico: auditoria (quem mudou o quê) | log próprio | ❌ (feature nossa) |
| Visitas do anúncio | `GET /items/$ID/visits` (API oficial de visits) | ✅ |

## Mercado.dc.html (radar)

| Bloco | Fonte | Status |
|---|---|---|
| Preços de concorrentes, mediana, posição na busca | `GET /sites/MLB/search?q=` (exige token desde 2023, dados públicos). Mediana/posição = calculadas por nós sobre o resultado | ⚠️ |
| Sugestão de preço oficial | `GET /marketplace/benchmarks/user/$USER_ID/items` (price references: suggested_price, lowest_price, percent_difference, costs) — usar como 2ª fonte do radar | ✅ |
| Catálogo: ganhar o destaque | `GET /products/$ID/items_to_win` (price to win) + webhook `item competition` | ✅ |
| Tendências/demanda (Oportunidades) | `GET /trends/MLB` + busca por categoria; cruzamento c/ catálogo ERP é nosso | ⚠️ |
| Monitorados: termos | busca salva + snapshot diário (nosso scheduler) | ⚠️ |
| Monitorados: vendedores | busca por seller (`/sites/MLB/search?seller_id=`) + snapshot | ⚠️ |
| Monitorados: anúncios específicos | `GET /items/$ID` (público) + snapshot p/ variação | ⚠️ |
| "Variação 7d" e alertas | exigem os snapshots — sem histórico retroativo | ⚠️ |

## Simulador.dc.html

| Bloco | Fonte | Status |
|---|---|---|
| Comissão por modalidade + taxa fixa < R$79 | `GET /sites/MLB/listing_prices?price=X` (retorna sale_fee por listing_type, inclui taxa fixa) — **não hardcodar 12%/17%/6,50: buscar da API por categoria+preço** | ✅ |
| Frete (≥ R$79 vendedor paga) | `GET /users/$ID/shipping_options/free` e `/items/$ID/shipping_options` (por CEP) | ✅ |
| Tarifa Full por unidade | ⚠️ custo real de fulfillment aparece no billing; por item/antecipado = estimativa nossa (param em Configurações, como está na UI) | ⚠️ |
| Custo do produto | ERP | ❌ (ERP) |
| Imposto (Simples/alíquota) | parâmetro nosso | ❌ (config nossa) |
| Preço↔margem bidirecional, cenários, veredicto | lógica nossa | ❌ (feature nossa) |
| Aplicar no anúncio | `PUT /items/$ID {"price": X}` via fila de sync | ✅ |
| Preço de mercado no painel | ver Mercado (busca + benchmarks) | ⚠️ |

## Vinculos e Importacao.dc.html

| Bloco | Fonte | Status |
|---|---|---|
| SKU do anúncio | `seller_sku` / `seller_custom_field` no item | ✅ |
| GTIN/EAN do anúncio | atributo `GTIN` em `attributes` do item | ✅ |
| Match SKU→GTIN→título + score de confiança | algoritmo nosso (os dados dos dois lados existem: item ML + produto ERP) | ⚠️ (lógica nossa) |
| "Vincular não altera o ML" | correto — vínculo é tabela nossa | ❌ (feature nossa) |
| Protocolo de lote | nosso | ❌ (feature nossa) |

## Repasses.dc.html

| Bloco | Fonte | Status |
|---|---|---|
| Extrato de liberações (conciliação) | **Mercado Pago Releases Report API**: `POST /v1/account/release_report` (gerar por período), config em `/v1/account/release_report/config`, agendável. Inclui bloqueios/desbloqueios, disputas, reembolsos, saques | ✅ |
| Tarifas ML por período | Billing API: `GET /billing/integration/periods/...` (+ `/payment/details`) | ✅ |
| Previsão de recebimento (calendário) | `money_release_date` em cada payment (webhook `payments`) — somar por data = nossa agregação | ⚠️ |
| Vendas − tarifas − fretes − retenções por repasse | derivado do release report (linha a linha, `SOURCE_ID` ↔ payment ↔ order via `EXTERNAL_REFERENCE`) | ⚠️ |
| Retenções (reclamações/reversas) | release report marca bloqueios por claim; cruzar c/ claims API | ✅ |
| Chip "bateu/divergente" c/ ERP + ação Conciliar | comparação nossa (extrato ML × lançamento ERP) | ❌ (feature nossa) |
| Nota: report é assíncrono (gera em minutos, CSV) — arquitetura deve gerar/ingerir agendado, não on-demand na tela | | ⚠️ |

## Configuracoes.dc.html

| Bloco | Fonte | Status |
|---|---|---|
| Tabela DIFAL por UF (padrão + exceções) | 100% nossa — ML não fornece dados fiscais estaduais. "Padrão 2026" deve ser mantido por nós (seed no banco) | ❌ (feature nossa) |
| Parâmetros de cálculo globais | nossos | ❌ (feature nossa) |
| Integrações (status ML/ERP) | token OAuth válido + health-check; ERP conforme fornecedor | ⚠️ |
| Notificações | nossas (alimentadas por webhooks ✅) | ❌/✅ |

## Dashboard.dc.html
Agregações dos mesmos endpoints acima (orders, items, claims, billing) — nada exclusivo. ⚠️ tudo derivado.

---

## Resumo executivo p/ implementação

**Direto da API (baixo risco):** pedidos, envios, etiquetas, claims/devoluções/reputação, anúncios + preço/estoque, comissões e frete reais (listing_prices/shipping_options), estoque Full, billing + release report (repasses), visits, benchmarks de preço, webhooks.

**Derivado (médio — exige agregação/scheduler):** mediana/posição de mercado, monitorados + variação (snapshots diários próprios, sem retroativo), previsão de recebimentos, retorno real por pedido, métricas por produto.

**Nosso banco/ERP (a API não ajuda):** DIFAL (tabela + agendamento + lembretes), custo/estoque físico/NCM (ERP), vínculos anúncio↔produto, fila de sync com protocolo, cenários do simulador, auditoria, conciliação ERP, histórico de preços.

**Pegadinhas confirmadas na doc:** billing-info legado deprecado (usar fluxo novo); inbound Full não tem API (só consulta); release report é assíncrono/CSV; `/sites/MLB/search` exige token; comissões variam por categoria — nunca hardcodar; responder claim em 48h pode evitar impacto de reputação (`has_incentive`).
