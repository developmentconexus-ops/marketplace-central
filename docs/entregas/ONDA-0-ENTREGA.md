# Onda 0 — Ficha de Entrega

`status: FASE 1 (cadeia no código) EM CURSO — fase 2 (banco real + live drive) pendente`
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

*Cadeia ainda não percorrida — próximo passo da fase 1.*

---

## Fase 2 — o que falta medir

1. Contagem no banco real dos 39 pedidos: quantos têm UF de destino, quantos têm **todos**
   os itens com vínculo interno e snapshot fiscal, quantos produzem margem não-nula.
   Separa VIVO de MUDO em toda a seção 1.
2. Linha `orders` no `/sync/health` visível na tela.
3. Live drive das quatro telas, com container id e `Created` no veredito.
4. Lane de governança por diff de conjunto contra `main`, medida fora da árvore montada
   (decide o F-4).
