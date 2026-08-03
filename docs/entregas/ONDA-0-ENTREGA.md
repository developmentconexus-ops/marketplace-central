# Onda 0 — Ficha de Entrega

`status: FASE 1 (cadeia no código) COMPLETA — fase 2 (banco real + live drive) pendente`
`método: docs/METODO-DE-REVISAO.md`
`base medida: main @e66ce013 (merge de f00-scheduler-pedidos + worktree-p2-dinheiro-real-pedidos + p2b-imposto-ex-ante + fa3-idade-honesta)`

> **Leia isto primeiro.** Toda linha "Veredito" desta fase é o que o **código** permite.
> Distinguir VIVO de MUDO exige contar linhas na base real — é a fase 2. Até lá, um campo
> com cadeia completa aparece como `VIVO?`, não como VIVO.

---

## Tela `/pedidos`

Rota FE: `apps/web/src/pages/pedidos/PedidosPage.tsx` ·
operações: `listOrders` (`GET /orders`), `getOrder` (`GET /orders/{provider_order_id}`),
`getOrdersSummary`, `markOrderFaturado`.

Fato de arquitetura que vale para a tela inteira: **lista e detalhe passam pelo MESMO
enriquecimento**. `handleReadList` (`transport/http_handler.go:215-222`) e `handleGet`
(`:261`) chamam ambos `mapEnrichedOrder`. Não existe campo que só o drawer saiba calcular —
o que aparece "—" na lista aparece "—" no drawer, pela mesma razão.

### 1.1 Colunas da lista (`PedidosTable.tsx:95-146`)

| Campo na tela | Componente | Chega por | Coluna/origem | Se desconhecido | Veredito |
|---|---|---|---|---|---|
| Pedido | `:102` | `provider_order_id` | `orders_marketplace_orders.provider_order_id` ← ML `GET /orders/search` | nunca (chave) | VIVO? |
| Data | `:107` | `provider_created_at` | idem, `date_created` do ML | `UnknownValue` | VIVO? |
| Comprador | `:112` | `buyer.display` | colunas `buyer_*` (migração 0089) | `UnknownValue` | VIVO? |
| Itens | `:117` | `items[].title` | `order_items` ← payload ML | `UnknownValue` | VIVO? |
| Valor | `:122` | `total` | `orders_marketplace_orders.total` | `UnknownValue` | VIVO? |
| Retorno | `:48-62` | `retorno_liquido` + `margem_pct` | **derivado**, `domain.BuildProfitability` | `UnknownValue` + hint | **ver 1.3** |
| SLA | `:26-37` | `sla.due` / `sla.atrasado` | `order_shipments` (0088) | `UnknownValue` | VIVO? |
| DIFAL | `:66-78` | `difal.amount` / `paid` / `due_date` | motor D-41 (`adapters/pricingtax`) | `UnknownValue` + hint | **ver 1.3** |
| Ação | `:80-93` | `bucket` | derivado `DeriveOrderBucket` | sempre derivável | **botão inerte** |

**A coluna Ação é decorativa.** `renderAcao` (`:83-91`) renderiza o botão com `disabled`
fixo e `title="disponível em breve"`. O rótulo muda com o bucket, o clique não faz nada.
Veredito: **ÓRFÃO deliberado** — precisa estar dito na tela, não só no código.

### 1.2 Drawer de detalhe (`PedidoDrawer.tsx`)

Blocos: Itens · Decomposição + DIFAL · Comprador fiscal (ERP) · Linha do tempo · Fatos.

| Campo | Componente | Origem real | Veredito |
|---|---|---|---|
| custo unit. / `≈` / snapshot ERP | `:123-129` | `internalread` → Sankhya `TGFCUS` | VIVO?, com **D-22** (CODEMP fixo em 1) |
| CODPROD | `:130` | `items[].internal_product_id` ← vínculo | VIVO? |
| Comissão | `:224` | `sumSaleFee(order.Items)` (`enrich_service.go:393`) | VIVO? |
| Taxa fixa | `:225` | **nada** — `order_decomposition.go:132` empilha `taxa_fixa` em desconhecidos **sem condição** | **MUDO estrutural** |
| Frete | `:226` | `senderFreight(shipment)` (`:395`) | VIVO? |
| Imposto | `:227` | **nada** — `BuildProfitability` não popula mais (`order_decomposition.go:14-19`) | **FANTASMA na tela** |
| ICMS saída | `:228` | motor D-41 por item (`pricingtax/reader.go:58`) | VIVO? |
| DIFAL | `:229` | idem | VIVO? |
| PIS/COFINS | `:230` | idem | VIVO? |
| Restituição ST | `:231-235` | idem, é **crédito** (soma) | VIVO? |
| Tarifa Full | `:236` | **nada** — `order_decomposition.go:145`, mesmo caso da Taxa fixa | **MUDO estrutural** |
| Custo | `:237` | `sumItemCosts` (`:394`) | VIVO? |
| Margem valor / Margem % / Retorno líquido | `:239-241` | derivados | **ver 1.3** |
| DIFAL · Rota | `:248` | `Difal.UFRoute` — nunca preenchido sem decomposer (`enrich_service.go:380-383`) | **MUDO estrutural** |
| DIFAL · Vencimento / Pago | `:249-254` | idem | **MUDO estrutural** |
| Nota fiscal | `:340-346` | `nf_state` — é marcador de vínculo, **não é número de NF** | VIVO?, rótulo ambíguo |
| Rastreio / Modalidade / Transportadora / Destino / Destinatário / Frete real | `:347-408` | `order_shipments` (0088) | VIVO? |
| Motivo do cancelamento | `:385-389` | `cancellation_detail` (`0093_orders_status_details_nullable`) | VIVO? |
| Comprador fiscal (nome/doc/endereço) | `:442-455` | `buyer_fiscal_reader` sobre colunas 0089 | VIVO? |
| Linha do tempo | `:270-286` | só timestamps presentes; evento sem data é **omitido**, não inventado | VIVO? |
| "Foi faturado" | `:493-502` | `markOrderFaturado` → grava `faturado_at` na NOSSA base | **VIVO (única ação que age)** |
| Faturar via ERP / Etiqueta / Marcar enviado / DIFAL agendar / Devolução… | `:503-513` | nada — `disabled` fixo | **ÓRFÃO deliberado (ruling D-57)** |

### 1.3 A regra que decide Margem, Retorno e DIFAL

`domain.BuildProfitability` (`order_decomposition.go:154`) só produz margem quando
**todos** estes sete são conhecidos:

`Total`, `Comissao`, `Frete`, `ICMSSaida`, `Difal`, `PisCofins`, `Custo`, `RestituicaoST`

E os quatro tributários só existem se, para **todo** item do pedido (regra tudo-ou-nada,
`pricingtax/reader.go:46-49`):

1. o pedido tem **UF de destino** conhecida (`reader.go:59` — sem ela, tudo nulo);
2. o item tem **vínculo interno** (`InternalProductID`) e preço unitário (`:145`);
3. o produto tem **snapshot fiscal** em `products_mirror` (`:154`);
4. existe **célula da matriz ICMS** para (MG, destino, grupo_icms) e alíquota interna
   cadastrada para a UF de destino.

**Consequência operacional:** um único item sem vínculo em um pedido de dez zera a
margem do pedido inteiro — de propósito (ADR-17: margem parcial seria mentira). Quantos
dos pedidos reais passam nos quatro filtros é **a medição principal da fase 2**.

### 1.4 Achados desta tela

**F-1 — a tela explica o "—" com o motivo errado.** `PedidosTable.tsx:54` e todo o bloco
de hints do drawer (`PedidoDrawer.tsx:196-197`, `:224-241`) dizem ao operador
*"decomposição de custos ainda não disponível (hub C2)"*. Isso era verdade quando nenhum
produtor existia. Hoje `enrich_service.go:391-400` popula Comissão, Custo e Frete, e o
motivo real de um "—" é outro — e está **no próprio payload**, em
`decomposicao.componentes_desconhecidos`, que a lista nem lê. O operador recebe uma
explicação falsa quando a explicação verdadeira já viajou junto com o dado.

**F-2 — a linha "Imposto" não pode ter valor.** `PedidoDrawer.tsx:227` renderiza uma
linha cujo produtor foi aposentado (`order_decomposition.go:14-19`). Ela vai mostrar "—"
para sempre. É ruído que compete visualmente com ICMS saída/PIS-COFINS, que são os campos
que substituíram justamente ela.

**F-3 — "Taxa fixa" e "Tarifa Full" também nunca terão valor nesta onda**, por construção
(`order_decomposition.go:132` e `:145`, `append` incondicional). Diferente de F-2, aqui a
ausência está documentada e é honesta — mas a tela não distingue "ainda não" de "não por
esta via", e as duas somem no mesmo "—".

**F-4 — `mustRat` dá `panic` em adaptador de produção** (`pricingtax/reader.go:238-244`).
Candidato a `GOV_PRODUCTION_PANIC`; a lane de governança da fase 2 decide.

---

## Tela `/integracoes`

### 2.1 Cartão "Contas conectadas" (`ConnectionHealthCard.tsx`)

| Campo | Componente | Cadeia | Veredito |
|---|---|---|---|
| Nome da conta | `:85` | `installation.display_name` | VIVO? |
| Selo de estado | `:86-91` | `connection.state` ← `ProjectConnectionSnapshot` (`domain/connection_snapshot.go:86-166`) ← `integration_installations.status` | VIVO? |
| Botão **Reautorizar** | `:93-102` | `startIntegrationReauthorization` → `auth_handler.go:180` | **VIVO — age de verdade** |
| Motivo (erro cru do provider) | `:114-118` | `connection.reauth_reason` ← `cause.Error()` (`auth_flow_service.go:619`) | VIVO? |

**A cadeia da falha de token está fechada.** `degradeAfterRefreshFailure`
(`auth_flow_service.go:582-622`) classifica o erro: terminal (refresh token revogado) →
`requires_reauth` + health `critical`; transitório persistente acima da política →
`degraded` + "Repetindo automaticamente"; soluço isolado → **não muda a tela de propósito**
(`:602-607`), para não treinar o operador a ignorar o aviso. Era exatamente o buraco que a
MIS-008 abriu como "falha de refresh é invisível em todas as camadas". Está fechado no
código; a fase 2 confirma no navegador.

### 2.2 Cartão "Saúde do sync" (`SyncHealthCard.tsx`)

| Campo | Componente | Cadeia | Veredito |
|---|---|---|---|
| Linha por entidade | `:44-76` | `GET /sync/health` ← `sync_state` por `(tenant_id, installation_id, entity)` | VIVO? |
| Selo ok / N falhas / nunca | `:13-42` | `consecutive_failures` e `last_success_at` — **estado, nunca corte de tempo** | VIVO? |
| Idade | `:56-64` | `formatRelativeAge(last_success_at)`; sem sucesso = "nunca" | VIVO? |
| Fase | `:54` | `phase` (ADR-07: `backfill`/`incremental`/`sweep`) | VIVO? |
| Notificações (webhook) | `:78-114` | `webhook.last_notification_at` / `pending` / `dropped_24h` | VIVO? |

Dois acertos que valem registro: o selo vermelho **ganha do verde** mesmo com sucesso
recente (`:14-15`) — falhar agora é falhar; e o estado indeterminado do react-query
(`status:"pending"` com `fetchStatus` "paused"/"idle") renderiza "Estado desconhecido" em
vez de cartão em branco (`:182-193`), defeito que já apareceu ao vivo duas vezes.

**Pendência de fase 2:** a linha **Orders** só existe depois do primeiro ciclo do scheduler
do F-00. Às 20:50 do dia 03/08 o `sync_state` tinha `listings`, `market`, `market_queue`,
`products` — sem `orders`. Às 22:13 o F-00 já tinha `orders` com `fails=0`. Confirmar na
tela.

---

## Tela `/mercado`

| Campo | Componente | Cadeia | Veredito |
|---|---|---|---|
| Contagem por aba | `MercadoPage.tsx:219` | walk do catálogo + `/market` | VIVO? |
| Idade por linha | `RepricingTable.tsx:101`, `OportunidadesTable.tsx:95` | `market_signal.evidence.fetched_at` | VIVO? |
| **"fonte: busca pública ML · coletado X"** | `MercadoPage.tsx:284` | `X` = `summary.as_of` do **resumo de anúncios** (`:132`, `:141`, `:192`) | **F-5 — ver abaixo** |
| Margem mín: 12% | `:279-281` | constante no FE | ÓRFÃO deliberado (não é configurável) |
| Categoria ▾ | `:271-278` | `disabled` fixo | **ÓRFÃO deliberado** |

**F-5 — o cabeçalho atribui à busca pública ML uma data que é de outra coisa.** O texto
promete "coletado" para a "busca pública ML", mas o valor vem de
`getListingsSummary().as_of` — a idade da **sincronia de anúncios**, não a da evidência de
mercado. A data é real; a pergunta que ela responde é outra. Por linha a idade está certa
(`evidence.fetched_at`); o agregado do cabeçalho é que está trocado. Esta é a classe exata
que o método chama de MENTIRA por atribuição: nenhum teste reprova, a tela fica verde, e o
operador decide reprecificar olhando uma idade que não é a da evidência.

---

## Tela `/anuncios`

### 4.1 Colunas da linha (`AnunciosTable.tsx:180-210`)

| Campo na tela | Componente:linha | Coluna Postgres / origem | Origem externa | Veredito |
|---|---|---|---|---|
| `provider_listing_id` | `:185` | `listings.provider_listing_id` | ML `/items` multiget | VIVO |
| id da variação | `:186-188` | `listings.variation_id` | ML `/items` → `variations[].id` | VIVO |
| Título (botão) | `:189-197` | `listings.title` | ML `/items` → `title` | VIVO |
| Produto (vínculo) | `renderProductCell:65-78` | `product_links.state` / `product_id` | interno (M-05), não é ML | VIVO |
| Preço + chip de delta | `renderPriceCell:90-132` | `listings.price_amount` + `market_signal` | ML `/items` + evidência de mercado | VIVO (condicional — ver 4.2) |
| Qtd. publicada | `:200` | `listings.published_quantity` | ML `/items` → **`initial_quantity`** | VIVO |
| Pill `sync_state` | `:201` | `listings.sync_state` | — (escrito pelo mapper) | **MUDO por construção — F-7** |
| `pending_issue.message_pt` | `:202-208` | derivado (`read_service.go:769-783`) | — | VIVO parcial — F-7 mata 2 dos 4 ramos |
| Cabeçalho: `listing_count` | `renderGroupHeader:217` | `len(group.listings)` (`read_service.go:330`, `:370`) | — | VIVO |
| Cabeçalho: pill de erros | `renderGroupHeader` (`errorCount`) | `sync_state === "error"` | — | **FANTASMA — F-7** |
| Estoque do ERP | *removido* | sem produtor (ADR-C1) | — | removido de propósito, documentado no componente |

Uma precisão que evita um falso achado: o pill de erros do cabeçalho é calculado no FE
sobre `group.listings`, enquanto o número do cabeçalho vem de `listing_count`. Parece
assimétrico, mas não é — `read_service.go:330` e `:370` atribuem
`ListingCount = len(...)` da própria fatia carregada. Contagem e pill enxergam sempre o
mesmo conjunto. **Não é defeito.**

### 4.2 A regra que decide o chip de preço

`DeriveSignalStatus` (`listings/domain/signal.go:65-76`) é pura e tem quatro saídas:

```
sem product_id                      -> SEM_VINCULO        (célula vira link "sem vínculo")
com vínculo, sem sinal              -> NO_PRICE_EVIDENCE  (preço puro, sem chip)
sinal com fetched_at > 1h           -> STALE
caso contrário                      -> OK
```

O horizonte de 1 h é `signalStaleTTL` (`signal.go:26`), espelhando o `snapshotTTL` do
módulo market (`collection_pipeline_service.go:19`). Isso liga o chip desta tela a uma
aritmética que vive em outro arquivo — ver F-6.

### 4.3 Achados desta tela

**F-6 — o teto de renovação é menor que a validade da evidência.** Três constantes,
escritas em três arquivos, formam um sistema:

| Constante | Valor | Onde |
|---|---|---|
| validade do snapshot / horizonte STALE | 1 h | `collection_pipeline_service.go:19`, `signal.go:26` |
| cadência do scheduler de mercado | 30 min | `composition/root.go:740` |
| teto de produtos por ciclo | 50 | `composition/root.go:745` |

Renovação máxima sustentada = 50 produtos a cada 30 min = **100 produtos/hora**. A
evidência vence em 1 h. Logo, se o operador tiver mais de ~100 produtos com evidência
coletada, uma parte deles fica **permanentemente STALE**, e nenhuma tela explica por quê —
o chip diz "desatualizado" e some a razão. O catálogo vendável ratificado tem 2.923
produtos; a coleta é manual por aba em `/mercado` ("Atualizar agora"), então o conjunto
coletado cresce conforme o operador usa a ferramenta, e o teto é atingido pelo uso normal,
não por um caso extremo. Também vale registrar que `Start` só dispara no primeiro tick
(`sync/application/scheduler.go:109`): depois de um restart, a primeira renovação de
mercado é 30 min depois. Quantos produtos têm agregado hoje é medição da fase 2.

**F-7 — `sync_state` de anúncio nunca sai de `synced`.** O único produtor de
`sync_state` no caminho de sincronia é `multiget_mapper.go:114`, que escreve a constante
`ListingSyncStateSynced`. Uma varredura por escritores de `listings` encontra exatamente
dois em produção: o upsert (`postgres/repository.go:448`, que recebe o valor do mapper) e
`UPDATE listings SET absent_since` (`:543`). Nenhum escreve `'error'`, `'stale'` ou um
`sync_error` não-nulo — só os testes escrevem. Consequência em cascata, toda invisível às
lanes porque os testes fabricam o estado que a produção não produz:

- o pill de erros do cabeçalho é sempre `0`;
- `pending_issue` perde 2 dos seus 4 ramos (`sync_error`, `stale`) —
  `read_service.go:770-775` é código inalcançável em produção;
- os contadores `sync_error` e `stale` de `/listings/summary`
  (`postgres/repository.go:320-321`) são sempre `0`;
- os filtros de exceção `sync_error` e `stale` (`repository.go:112`, `:114`) sempre
  retornam vazio.

Isto não é um bug de uma linha: é uma **promessa de superfície sem produtor**. Ou o sync de
anúncios passa a marcar falha por item, ou os quatro consumidores acima são retirados da
tela. A escolha é do operador; o que não pode continuar é a tela afirmar "0 erros" quando o
que ela sabe dizer é "não sei medir erro".

---

## Fase 2 — o que falta medir

1. Contagem no banco real dos 39 pedidos: quantos têm UF de destino, quantos têm **todos**
   os itens com vínculo interno e snapshot fiscal, quantos produzem margem não-nula.
   Separa VIVO de MUDO em toda a seção 1.
2. Linha `orders` no `/sync/health` visível na tela.
3. `select sync_state, count(*) from listings group by 1` — confirma o F-7 no dado real
   (esperado: uma única linha, `synced`).
4. Contagem de produtos com agregado de mercado e distribuição de idade — dimensiona o F-6
   (acima de ~100, a cauda permanentemente STALE existe hoje, não no futuro).
5. Live drive das quatro telas, com container id e `Created` no veredito.
6. Lane de governança por diff de conjunto contra `main`, medida fora da árvore montada
   (decide o F-4).
