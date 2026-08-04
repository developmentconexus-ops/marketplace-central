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

## Fase 2 — medição no banco real

Ambiente medido: `marketplace-central-postgres-1`, banco `marketplace_central`, em
2026-08-03. Toda contagem abaixo é `select` executado, não estimativa.

### 5.0 O binário que estava rodando era mais velho que o merge

Primeira medição da fase, antes de qualquer tela: a imagem `marketplace-central-backend`
foi construída em `2026-08-03T20:00:19Z`; os merges terminaram às `22:29Z`
(`e66ce013`). **A stack de pé estava 2h29 atrás do código.** A prova independente está no
banco: `schema_migrations` tem 82 linhas contra 83 arquivos no repo, e a única faltante é
exatamente `0093_icms_matrix.sql` — a migração que veio no merge de `p2b-imposto-ex-ante`.

Sem esse passo, o live drive teria medido a onda anterior e chamado de onda 0. A stack foi
reconstruída antes de qualquer captura de tela.

### 5.1 Insumos de `/pedidos` (39 pedidos reais)

| Insumo | Medido | Consequência |
|---|---|---|
| pedidos | 39 | — |
| itens | 39 (1 por pedido) | — |
| itens com `internal_product_id` | **34 / 39** | 5 pedidos sem vínculo → todo componente fiscal nil |
| itens com `unit_price` | 39 / 39 | ok |
| itens com `sale_fee_amount` | 39 / 39 | Comissão resolve para todos |
| remessas com `cost_seller` | 39 / 39 | Frete resolve para todos |
| remessas com `dest_state` | 39 / 39, em 9 UFs (SP 15, RJ 11, MG 3, PR 3, BA 2, ES 2, RS 1, PE 1, MA 1) | `destinoUF` nunca é vazio |
| `icms_matrix_mirror` | **tabela inexistente** (migração não aplicada; após o rebuild, existente e vazia) | ver F-8 |
| `products_mirror` | 10.632 linhas; colunas fiscais só passam a existir após a migração | ver F-8 |
| `orders_marketplace_orders.decomposition` / `net_amount` / `margin_pct` | 0 não-nulos em 39 | ver F-9 |

Os 5 pedidos sem vínculo batem exatamente com a previsão do plano do P2.b ("5 pedidos sem
vínculo de produto", Task 7). A previsão estava certa; o que não aconteceu foi a Task 7.

### 5.2 `/anuncios` — F-7 confirmado no dado

```
select sync_state, count(*) from listings group by 1;
 sync_state | count
------------+-------
 synced     |    34
```

Uma única linha. Zero variância na coluna que a tela transforma em pill, em contador de
resumo e em dois filtros de exceção. O F-7 sai de "por construção" para **medido**.

### 5.3 `/mercado` — F-6 dimensionado

`market_aggregates` tem **4 linhas**. O teto aritmético do F-6 (100 renovações/hora contra
validade de 1 h) portanto **não está apertando hoje** — com 4 produtos, o ciclo de 30 min
renova tudo com folga. O achado continua válido como propriedade do sistema, não como
sintoma atual: ele passa a morder quando o operador coletar mais de ~100 produtos, o que o
uso normal de "Atualizar agora" produz sozinho. Registrar agora é mais barato que descobrir
depois olhando uma tela cheia de "desatualizado" sem explicação.

### 5.4 Achados novos da fase 2

**F-8 — a matriz de ICMS tem leitor, escritor, testes e nenhum chamador.**
`internal_read/adapters/oracle/icms_matrix.go` (`ICMSMatrixReader`) e
`internal_read/adapters/mirror/icms_matrix_writer.go` (`ICMSMatrixWriter`) existem, estão
cobertos por teste de integração, e são construídos em **exatamente dois lugares no repo**,
ambos arquivos `_test.go`. Nenhuma composição, nenhum job de sync, nenhum endpoint, nenhum
binário em `cmd/` os invoca. O plano do P2.b previa isso: rodar o sync da matriz era o
primeiro passo da **Task 7**, que nunca foi executada.

Consequência em cadeia, e é a mais cara desta onda:

```
icms_matrix_mirror vazia
  -> MatrixCellFor não acha célula para nenhum (UF, grupo)
  -> pricing/domain devolve componente desconhecido por item
  -> reader.go all-or-nothing: ICMS saída, DIFAL, PIS/COFINS e restituição ST = nil p/ o pedido inteiro
  -> BuildProfitability não fecha os sete componentes
  -> Margem e Retorno líquido = nil em 39 de 39 pedidos
```

Ou seja: **a funcionalidade central do P2.b está construída e desligada.** O comportamento
é honesto (a tela mostra "—", não um número inventado — ADR-17 preservado), mas a promessa
da onda não está entregue. E a diferença entre "construído" e "entregue" é invisível a
todas as lanes: `go build` passa, 120 pacotes de teste passam, o front compila, e nenhuma
tela fica vermelha.

Isto é a mesma classe que a MIS-006 já pagou uma vez ("scheduler construído mas nunca
ligado", `composition/root.go:730-737` documenta a lição no próprio código). Reincidência de
classe: cabe **stop-the-line** — o conserto não é só ligar esta matriz, é ter um observável
que reprove porta de produção sem chamador de produção.

**F-9 — três colunas de pedido sem leitor e sem escritor.**
`orders_marketplace_orders.decomposition`, `net_amount` e `margin_pct` (migração 0089) têm
0 não-nulos em 39 linhas. Uma varredura por escritores e leitores em Go encontra **zero** de
cada: `orders/domain/order.go:80` diz em comentário que os três são "escopo do M-06,
deliberadamente ausentes" — e o M-06 fechou sem preenchê-los. O caminho vivo calcula a
decomposição em tempo de leitura (`enrich_service.go:390`), que é a decisão certa. As três
colunas são resíduo de um desenho abandonado. Não mentem para o operador (nada as lê), mas
são exatamente a redundância que a revisão de código deve remover: ou nascem um consumidor,
ou saem do esquema.

---

## §6 — As cinco buscas de máximo local (método §4)

Todas com número medido, não impressão.

### 6.1 Fórmula sem consumidor (busca 2)

Varredura de todo `apps/server_core/internal/modules/*/domain/`, contando usos de cada
função exportada **fora do arquivo que a declara** (testes incluídos como consumidor
legítimo):

```
funções exportadas em */domain/: 98
sem nenhum uso fora do arquivo de declaração: 3
  - IsAutomatic            product_links/domain/product_link_decision.go
  - IsSankhyaLinkageError  internal_read/domain/sankhya_linkage.go
  - ValidateRowLenient     erp_import/domain/validation.go
```

**3 em 98.** A suspeita de "fórmula morta em massa" no domínio está refutada por medição.
As três são dívida pequena e nomeada, não um sintoma estrutural.

### 6.2 Operação de contrato sem consumidor (busca 3)

Três contagens, porque a pergunta tem três lados:

| Medida | Número |
|---|---|
| `operationId` no OpenAPI | 111 |
| símbolos expostos pelo `sdk-runtime` | 114 |
| `operationId` sem símbolo homônimo no SDK (nem a string aparece lá) | **13** |
| `operationId` sem nenhum consumidor de produção no monorepo FE | **50 / 111** |

O segundo número é o achado; o terceiro é contexto (metade das operações não tem tela, o
que para uma plataforma em construção é esperado, não defeito).

**F-10 — o contrato e o SDK divergem em nome, e o gate mede a coisa errada.**
Dos 13, seis são renomeações silenciosas — o contrato publica um nome, o cliente expõe
outro:

| OpenAPI | `sdk-runtime` |
|---|---|
| `getOrdersSummary` | `getOrderSummary` |
| `getMarketplaceOrder` | `getOrder` |
| `upsertProductEnrichment` | `updateProductEnrichment` |
| `runPricingBatchSimulation` | `runBatchSimulation` |
| `createPricingSimulation` | `runPricingSimulation` |
| `getMelhorEnvioOAuthStatus` | `getMelhorEnvioStatus` |

Os outros sete não têm par nenhum. Caso verificado ponta a ponta:
`pricingGetTariffDefaults` / `pricingPutTariffDefaults` existem no contrato
(`marketplace-central.openapi.yaml:2773`, `:2795`) **e** no servidor
(`pricing/transport/calc_handler.go:37-38`), e não existem no SDK nem no front. Backend
serve, contrato publica, ninguém consome.

O que torna isto uma classe, e não seis erros de digitação: o invariante que deveria pegar
mede outra coisa. `GOV_API_SDK_SPLIT` (`scripts/harness/Policy.psm1:452-454`) verifica
apenas se `contracts/api/…yaml` e `packages/sdk-runtime/` **mudaram no mesmo diff** —
atomicidade de commit, nunca correspondência de nomes. Renomear no SDK sem tocar no
contrato passa; publicar operação no contrato sem gerar método passa. O gate está verde
porque pergunta "mexeram nos dois?", não "os dois dizem a mesma coisa?".

### 6.3 Motor duplicado, campo sem produtor, abstração vazia (buscas 1, 4 e 5)

- **Busca 1 (motor duplicado):** nenhum encontrado nesta onda. O risco histórico
  (segundo motor de preço ao lado de `pricing/domain`) foi eliminado no re-escopo do
  P2.b — os 17 tasks viraram 7 exatamente por isso.
- **Busca 4 (campo sem produtor):** é o veredito ÓRFÃO/FANTASMA da §1.2, já contabilizado
  nas seções por tela — F-2, F-3, F-7 e F-9.
- **Busca 5 (abstração que não abstrai):** as portas desta onda protegem o domínio de
  dependências que mudam de verdade (Oracle, ML, Postgres) e têm implementação real mais
  implementação de teste. Não é cerimônia. Uma exceção a observar, não a cortar: o par
  `ICMSMatrixReader`/`ICMSMatrixWriter` do F-8 tem porta, implementação e teste — e nenhum
  chamador de produção. Não é abstração vazia; é abstração **desligada**.

---

## Fase 2 — o que ainda falta

1. Live drive das quatro telas com a stack reconstruída, com container id e `Created` no
   veredito (a reconstrução está em curso).
2. Linha `orders` no `/sync/health` visível na tela.
3. Lane de governança por diff de conjunto contra `main`, medida fora da árvore montada
   (decide o F-4). Medido até aqui: `governance-validate` passa, mas **não é ele que checa
   `panic`** — o check de `production-panic` vive em `Test-GovernanceDrift`
   (`Policy.psm1:360-373`), que exige `-BaseSha`. Não há exceção registrada para
   `orders/adapters/pricingtax/reader.go` em `invariants.json` (as quatro existentes são
   para outros arquivos), então a expectativa é reprovação.
