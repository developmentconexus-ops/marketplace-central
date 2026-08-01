# Resultado esperado de /anuncios e /pedidos — design

```yaml
data: 2026-08-01
status: aguardando revisão do operador
missão: MIS-007-ml-sync
escopo ratificado: A+B (buracos + margem/competição), regime tributário Normal
fontes: research/ml-api-capability-map.md · research/EXPECTED-RESULT-SPEC.md · research/live-probe-results.md · research/ADR-CANDIDATES.md
```

## 1. O problema real

`/anuncios` mostra preço errado. `/anuncios` e `/pedidos` têm vários campos em branco.

O diagnóstico do operador estava certo: *"nossos contratos não estavam bem definidos e foi
sendo implementado, implementado"*. A investigação confirmou o mecanismo exato, e ele é pior
do que "falta implementar".

Os quatro campos nulos em 34/34 anúncios (`price`, `currency`, `listing_type`,
`published_quantity`) **não morreram na tela nem no banco**. Morreram no primeiro milímetro:
o DTO do adapter **não declarava a chave** que o Mercado Livre manda
(`items_multiget_reader.go:136-161`). `encoding/json` ignora chave desconhecida por padrão —
sem erro, sem log, sem sintoma. O mapper não mapeia o que não existe no struct. O upsert
grava `NULL`. O contrato aceita `NULL` porque declara o campo `required` **e** `nullable`.
A tela mostra `—`.

**Nenhuma camada mentiu. Cada uma passou adiante fielmente um nada que nasceu na primeira.**

Mesma classe de defeito, quatro ocorrências, dois adapters. Não é bug: é buraco de método.

E há um agravante medido: `UpsertPulledRows`
(`listings/adapters/postgres/repository.go:436-446`) faz `SET x = EXCLUDED.x` cru em toda
coluna. Como o mapper de produção não preenche esses campos, **cada re-sync grava `NULL` por
cima de valor que já esteve correto**. Os preços não estão faltando — foram **apagados**.
`order_items` (`order_repo.go:856-878`) tem a mesma forma: arma carregada, ainda não disparada.

## 2. Princípio que governa o desenho

Todo campo carrega um **selo**, e o selo é uma afirmação verificável:

| selo | significado |
|---|---|
| `SEMPRE` | o provider sempre manda |
| `CONDICIONAL(x)` | só quando `x` |
| `SEM-FONTE` | a API não fornece |
| `PLANEJADO(M-xx)` | há fonte, ninguém implementou — **com milestone dona nomeada** |
| `NÃO-SUPORTADO` | este provider não tem o conceito |

A regra anti-disfarce importa mais que a tabela: **`SEM-FONTE` é para o que a API não dá.
Campo não implementado é `PLANEJADO`.** Foi marcando não-implementado como
"honest-unknown" que `sales_30d` e `quality_score` sumiram do radar por uma missão inteira.

E o selo não é comentário — ele **governa a forma do contrato** (ADR-C3), então uma
classificação errada quebra o build em vez de virar um traço na tela.

## 3. Arquitetura — o terceiro no jogo

O sistema fala com marketplaces por **adapter**. O domínio não conhece Mercado Livre
(ADR-C4). Vocabulário de provider morre no adapter, e isso é verificável por `grep`, não por
julgamento: `sale_fee`, `price_to_win`, `sub_status`, `marketplace_fee`, `logistic_type` fora
de `adapters/mercado_livre/` reprovam no gate.

Com múltiplos providers aparece um estado que ADR-17 não cobria (ADR-C1):

| estado | representação | tela |
|---|---|---|
| suportado, com valor | valor | mostra |
| suportado, desconhecido | `null` | `—` **com o motivo** |
| **não suportado pelo provider** | capability ausente | **coluna some** |

Colapsar "não-suportado" em `null` faz o operador ler "ainda não me responderam" e esperar
por um dado que nunca vem. É mentira diferente do zero, com o mesmo estrago.

### Conceitos canônicos e a tradução do ML

| conceito canônico | Mercado Livre traduz de |
|---|---|
| saúde do anúncio | `GET /item/{id}/performance` |
| motivo de bloqueio | `sub_status[]` + `GET /moderations/infractions/{user_id}` |
| competição de oferta | `GET /items/{id}/price_to_win?version=v2` |
| comissão da linha | `order_items[].sale_fee × quantity` |
| custo de envio do vendedor | `GET /shipments/{id}/costs → senders[].cost` |
| motivo de cancelamento | `cancel_detail{group, code, description, requested_by, date}` |

### Fato de custo que governa a cadência

Não existe multiget de shipments (404 medido). Visitas, performance, `price_to_win` e
infrações por item são **1 chamada por item**. O custo de sincronização é dominado por
sub-recurso 1-a-1, **não** pela listagem.

⇒ Cadência **por conceito**, não uma só: preço/estoque em ciclo curto; saúde e visitas em
ciclo longo. Tratar tudo no mesmo ciclo transforma 34 anúncios em 100+ chamadas por rodada.

## 4. /anuncios — resultado esperado

**Identidade e estado.** id, título, permalink, thumb, `status` verbatim, categoria e domínio.
Já funciona. Falta: **motivo do bloqueio** — `sub_status[]` (medido: `[]` quando ativo, ou
seja, só carrega informação quando há problema) somado às infrações.

**Preço e estoque** — o buraco principal. `price`, `currency`, `listing_type`,
`available_quantity`, `initial_quantity` são todos `SEMPRE`; quatro deles estão `NULL` em
34/34 hoje.

Duas armadilhas de semântica:
- **`sold_quantity` é vitalício, não janela.** Usar como "vendas do mês" é erro de leitura.
- **`initial_quantity` ≠ `available + sold`** — reposição incrementa disponível sem tocar o
  inicial. São dois números distintos, não redundância.

`sales_30d` é **derivado** de `/orders/search` por item e data — a API não tem endpoint de
vendas por janela. O contrato precisa dizer que é derivado.

Estoque ERP é `CONDICIONAL(vinculado)` e hoje é a string chumbada `"ERP est. —"`
(`AnunciosTable.tsx:223-266`).

**Saúde do anúncio** — fonte confirmada com payload real: `score`, `level`, `level_wording`,
e `buckets[].variables[].rules[]` com `status`, `progress` (0–1) e `wordings.link` apontando
a ação de correção. Não é uma nota: é a lista acionável de o-que-corrigir. Infrações são
**1 chamada para a conta inteira** (barata); performance e visitas são 1-a-1 (ciclo longo).

**Competição** — `PLANEJADO`, com bloqueio externo nomeado. `price_to_win` responde 200, mas
devolve `item_not_opted_in` mesmo com o anúncio ativo: falta inscrição em catálogo do lado do
ML. Provado por controle — mesmo resultado com o item pausado e ativo, o que elimina o status
como causa. Destrava com opt-in do operador, não com código.

**Tarifa.** Não existe "tarifa deste anúncio publicado" — só a função paramétrica
`listing_prices`. Tarifa por anúncio é composição nossa.

**Fora de escopo, declarado:** estoque Full (`inventory_id` null em toda a amostra) — fica
`CONDICIONAL(inventory_id≠null)`, não implementado.

## 5. /pedidos — resultado esperado

**O pedido.** `status` e `status_detail` verbatim. Duas correções:
- **`cancellation_detail` está vazio sempre.** Lê `status_detail`, que é `null` em 7/7 dos
  pedidos cancelados. A fonte real é o objeto `cancel_detail`, presente em 7/7, com
  `requested_by` ∈ {buyer, meli, seller}. Não é dado duplicado — é dado zero.
- `string` em vez de `*string` no tipo de domínio torna `NULL` **inalcançável**: desconhecido
  vira `''`. Fere ADR-17 no tipo, antes de qualquer código rodar.

**Dinheiro real** — o núcleo do bloco B, e onde estão as duas armadilhas caras:

| conceito | fórmula | armadilha medida |
|---|---|---|
| receita bruta da linha | `(unit_price + discounts.full) × quantity` | — |
| **comissão** | `sale_fee × quantity` | `marketplace_fee` veio **0** no mesmo pedido em que `sale_fee` era 120,43 |
| **frete do vendedor** | `/shipments/{id}/costs → senders[].cost` | `lead_time.cost` é o custo do **comprador**; 5/5 com `senders.cost` de 23,65 a 138,95 e `lead_time.cost` = 0 |
| subsídio do ML | `costs.senders[].discounts[]` | explica por que o custo caiu |
| custo do produto | espelho Sankhya | `CONDICIONAL(vinculado)` |
| margem | receita − comissão − frete vendedor − custo produto | qualquer cálculo por `marketplace_fee` está inflado |

`sale_fee` é **por unidade** — provado por aritmética no payload real (120,43 / 729,90 =
16,5 %; se fosse total da linha daria 8,25 %, fora de qualquer faixa de comissão do ML). A
documentação é incapaz de desambiguar: todos os exemplos dela usam `quantity=1`.

**Margem é `CONDICIONAL(vínculo ERP)`.** Sem vínculo não há custo de produto ⇒ não se mostra
margem parcial, mostra-se "sem custo vinculado". Margem pela metade é pior que margem nenhuma.

**Envio.** `logistic.type` e `tracking_method` nunca foram gravados — o DTO não declara as
chaves. `/shipments/{id}/sla` **nunca é chamado**; o `SLALimitAt` de hoje vem de `lead_time`,
que é outra grandeza. Faltam ainda `/history` (linha do tempo) e `/invoice_data` (nota
fiscal — hoje `nf_state` é derivado de vínculo, não de nota).

Substatus de envio: **100+, sem tabela canônica na documentação**, e o `invoice_pending` da
nossa própria conta não aparece em nenhuma lista oficial. ⇒ verbatim, sem enum, sem `CHECK`.
Confirma ADR-06/IC-07.

**KPIs.** Hoje Novos/Faturar/Enviar/Enviados são contados **no cliente sobre a página
paginada** (`PedidosPage.tsx:41-64`) — mudam ao virar a página. Alvo: agregação server-side
sobre o conjunto todo.

## 6. DIFAL

O KPI "DIFAL A PAGAR" é constante chumbada que nunca leu payload nenhum. Decisão do operador:
tem que ser calculado de verdade — venda de MG para GO gera DIFAL.

**DIFAL não é campo do provider. É cálculo nosso** (ADR-C5). O ML fornece apenas insumos.

Regime ratificado: **Normal** (Lucro Real/Presumido) — o que mantém o DIFAL devido na venda
interestadual a consumidor final não contribuinte.

```
base_difal = (valor_operação − ICMS_origem) / (1 − alíq_interna_destino)
DIFAL      = base_difal × alíq_interna_destino − valor_operação × alíq_interestadual
FCP        = base_difal × alíq_fcp_destino     (recolhimento separado)
```

Entrada: fato canônico de venda. Saída:
`{devido, base, valor, fcp, motivo_se_não_devido, parâmetros_usados}`.

**Os parâmetros usados são gravados junto com o valor.** Alíquota muda por decreto; sem isso
ninguém reproduz em auditoria o número calculado há seis meses.

Tabelas são **dado versionado por vigência, nunca constante em código**: alíquota
interestadual por par origem/destino, alíquota interna por UF × NCM, FCP por UF × NCM.

**Insumos e origem.** UF de origem e regime vêm da nossa configuração. UF de destino vem do
ML — e a medição resolveu a dúvida na prática: UF do envio == UF do `billing_info` em
**38/38**. Condição de contribuinte é detectável em `billing_info.taxes.taxpayer_type` +
`inscriptions.state_registration`. NCM vem do Sankhya, pelo vínculo.

**ADR-17 aplicado.** Sem NCM ou sem UF de destino, o DIFAL é **desconhecido, nunca zero**. O
KPI declara cobertura: `R$ X sobre 22 de 38 pedidos — 16 sem base fiscal`. Somatório que
finge cobrir tudo é a mentira mais cara desta tela.

**Decisões de contador — nenhuma vira código sem ratificação:** base legal da UF de destino;
frete entra ou não na base; existe UF em que o marketplace tem responsabilidade de retenção.
Os três entram como **parâmetro de configuração**, não como suposição do desenvolvedor.

## 7. Os cinco saltos e o mecanismo de cada um

A cadeia ML → DTO → domínio → banco → API → tela falhou em silêncio em cada elo. Cada salto
ganha um mecanismo que **falha alto**:

| salto | falha medida | mecanismo |
|---|---|---|
| ML → DTO | chave não declarada; `encoding/json` ignora sem erro | **raw persistido + reconciliação** (ADR-C6): chave no raw e ausente do DTO vira alerta. Pega campo novo do provider sem reler documentação |
| DTO → domínio | campo existe nos dois lados, ninguém ligou (`logistic_type`, `tracking_method`) | tipo canônico, sem vocabulário de provider (ADR-C4), verificável por `grep` |
| domínio → banco | `string` em vez de `*string` (NULL inalcançável); `EXCLUDED.x` cru apagando valor alheio | **writer declara as colunas que possui** (ADR-C2): o `UPDATE` toca só o conjunto declarado |
| banco → API | `required` + `nullable` = forma que não recusa payload nenhum | **selo governa o shape** (ADR-C3): `SEMPRE` → não-nulável, nulo reprova |
| API → tela | string chumbada, constante, KPI sobre página | célula nunca contém literal; capability ausente esconde coluna; KPI server-side com cobertura |

**Por que "writer declara colunas" e não `COALESCE` em tudo:** `COALESCE` nunca apaga — nem
quando deve. Anúncio que perde o SKU no provider ficaria com o SKU velho para sempre. Isso
troca perda de dado por **dado zumbi**, que é pior de auditar porque parece certo.

**Segurança do raw, inegociável.** O payload cru carrega PII: nome, documento e endereço do
comprador, e o endereço completo do vendedor. **Scrub no momento da captura, nunca depois.**
Raw sujo persistido é vazamento permanente; "limpo depois" não existe. `billing_info` cru
nunca persiste.

## 8. Impacto nos milestones

**M-06 (incremental/backfill) — risco fechado.** `order.date_last_updated.from` é o mecanismo
correto, confirmado. E a medição resolveu a dúvida que a documentação deixava aberta:
`/orders/search` tem teto de offset em **10000**, mas **aceita `search_type=scan`** (HTTP 200
com `scroll_id` real), apesar de a doc só documentar scan para `/items/search` e
`/questions/search`. Backfill de 12 meses acima de 10k pedidos é viável.

**M-05 / M-07 (simulação de tarifa) — alerta rebaixado.** A documentação chama `logistic_type`
de crucial para o `fixed_fee`. A medição refutou: 12 combinações na categoria MLB270310
devolveram `sale_fee_details` idêntico. Passar o parâmetro segue correto por higiene;
**depender** dele não se justifica com a evidência atual. Não vira critério de aceite sem
medição em segunda categoria.

**M-08 (webhooks) — o brief está errado.** A política real é **8 tentativas em janela de 1 h**
(o brief diz 5), o receptor precisa responder **200 em ≤ 500 ms**, `missed_feeds` guarda só
**2 dias**, e falha repetida **desativa o tópico** para a aplicação. Desativação silenciosa é
perda de sincronismo sem alarme — precisa de detecção própria.

**Pré-requisito promovido.** `listings.raw` está `NULL` hoje. Era dívida; com ADR-C6 vira
pré-requisito da reconciliação.

## 9. Escopo — o que fica de fora, e por quê

| item | razão |
|---|---|
| pós-venda / claims (reclamação, devolução, mediação) | domínio inteiro novo; estoura a MIS-007. Vira missão própria — e reclamação aberta é a coisa mais urgente de um pedido, então não some do radar |
| competição de catálogo (`price_to_win`, `boosts[]`) | `PLANEJADO`, bloqueio externo nomeado: falta opt-in de catálogo no ML. Destrava com ação do operador, não com código |
| estoque Full | conta não usa; `inventory_id` null em toda a amostra |
| visitas por anúncio | fonte confirmada, valor de produto não priorizado nesta missão |

## 10. Lacunas que continuam abertas

Nenhuma foi preenchida com suposição:

1. `boosts[]` — estrutura nunca observada populada (depende do opt-in).
2. `/performance` — amostra n=1; reconfirmar com mais anúncios ativos.
3. `logistic_type` em `listing_prices` — medido em uma única categoria.
4. Rota canônica de `billing_info` (a nossa responde 200; a documentada é outra).
5. Prefixo correto dos relatórios de faturamento (a própria doc grafa de dois jeitos).
6. Vocabulário de `payments[].status` no domínio ML.
7. Limiar real de 429 e header `Retry-After` (100 chamadas seguidas, zero 429).
8. API de devoluções (Returns) — existe, não detalhada.
9. Três decisões fiscais pendentes de ratificação do contador (§6).
10. Escopo OAuth e retenção/LGPD exigidos por `billing_info`.

## 11. Método — a lição que custou caro duas vezes

**Estado transitório da conta se disfarça de capacidade ausente da API.**

A primeira rodada de medição concluiu que `/performance` e `price_to_win` não serviam nesta
conta. Estava errado: a conta estava em **modo férias**, com os anúncios desativados. Depois
da reativação, `/performance` passou a responder 200 com dado real.

E o mesmo controle serviu para o caso oposto: `price_to_win` devolveu `item_not_opted_in`
**com o item pausado e com o item ativo**, o que elimina o status como causa e prova que
falta um opt-in separado. Mesma técnica, dois vereditos opostos, ambos sustentados.

Antes de declarar "a API não dá", confirmar o estado da conta no momento da medição.
