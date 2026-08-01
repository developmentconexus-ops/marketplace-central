# Candidatos a ADR — decisões do operador em 2026-08-01

```yaml
id: MIS-007-ADR-CANDIDATES
type: decision-draft
status: draft
owner: hub
parent: MIS-007
ratified_by: operador, 2026-08-01
```

Seis regras decididas pelo operador que **não existem em nenhum ADR hoje**. Ficam aqui como
rascunho até virarem ADR numerado. Nenhuma delas é nova invenção: cada uma nasceu de um
defeito medido nesta missão.

---

## ADR-C1 — Capability ausente é um terceiro estado

**Contexto.** ADR-17 distingue *tem valor* de *desconhecido*, e proíbe desconhecido virar
zero. Com mais de um provider aparece um terceiro caso que ADR-17 não cobre: o provider
**não tem o conceito**. Saúde de anúncio existe no Mercado Livre (`/item/{id}/performance`);
outro marketplace pode não ter nada equivalente.

**Decisão.** A porta declara capability. Três estados, distintos em tipo:

| estado | significado | representação | tela |
|---|---|---|---|
| suportado, com valor | provider deu o dado | valor | mostra |
| suportado, desconhecido | provider suporta, não sabe agora | `null` (ADR-17) | `—` **com o motivo** |
| **não suportado** | provider não tem o conceito | capability ausente | **coluna some** |

**Por quê.** Colapsar não-suportado em `null` faz a tela mostrar `—`, e o operador lê "o
provider ainda não me respondeu". É mentira diferente do zero, mesmo estrago: leva a esperar
por um dado que nunca vem.

**Consequência.** Toda porta opcional declara capability em compile-time (mesma técnica dos
asserts que pegaram o catalog-503 no M-02). Tela consulta capability antes de renderizar
coluna. Contrato: campo não suportado é **ausente** do payload, não nulo.

---

## ADR-C2 — Writer declara as colunas que possui

**Contexto medido.** `UpsertPulledRows`
(`internal/modules/listings/adapters/postgres/repository.go:436-446`) faz
`ON CONFLICT DO UPDATE SET x = EXCLUDED.x` cru em toda coluna. O mapper de produção não
preenche `price_amount`, `price_currency`, `listing_type_code`, `published_quantity` ⇒ **cada
re-sync grava NULL por cima de valor que já esteve correto**. Foi assim que os preços da tela
morreram. `order_items` (`order_repo.go:856-878`) tem a mesma forma — arma carregada, ainda
não disparada.

**Alternativas descartadas.**
- *`COALESCE` em tudo*: nunca apaga, mas nunca apaga quando **deve**. Anúncio que perde o SKU
  no provider fica com o SKU velho para sempre. Troca perda de dado por **dado zumbi**, que é
  pior de auditar porque parece certo.
- *Regra por coluna documentada ao lado dela*: funciona, mas é disciplina humana. Coluna nova
  entra errada por esquecimento e nada no tipo impede.

**Decisão.** Cada escritor declara explicitamente o conjunto de colunas do qual é dono. O
`UPDATE` toca **somente** esse conjunto. Sem `COALESCE`, sem adivinhação.

**Consequência.**
- Ninguém apaga campo de outro escritor — o backfill do ML não encosta em coluna do espelho
  Sankhya nem do xlsx.
- O dono **ainda consegue apagar de verdade**: se o ML remover o valor, o dono escreve `NULL`
  e isso é a verdade, não perda.
- Casa com ADR-04 (um dono por seam), agora no grão de coluna, não de tabela.
- Escritor sem conjunto declarado não compila / não passa no gate.

---

## ADR-C3 — O selo governa a forma do contrato

**Contexto medido.** `ListingReadModel` no OpenAPI declara `price`, `published_quantity`,
`quality_score`, `sales_30d` como **`required` E `nullable`**. Campo obrigatório que pode ser
nulo é uma forma que **não recusa payload nenhum**. Não existe resposta que reprove. Foi por
essa porta que 34 preços nulos chegaram à tela sem nenhum gate acusar.

**Decisão.** Todo campo carrega um selo, e o selo determina a forma no OpenAPI:

| selo | significado | forma no contrato |
|---|---|---|
| `SEMPRE` | provider sempre manda | `required`, **não** `nullable` — nulo **reprova** |
| `CONDICIONAL(x)` | só quando `x` | `nullable`, com `x` nomeado na descrição |
| `SEM-FONTE` | API não fornece | campo **não existe** |
| `PLANEJADO(M-xx)` | há fonte, ninguém implementou | campo **não existe ainda**, com milestone dona |
| `NÃO-SUPORTADO` | provider não tem o conceito | ausente (ADR-C1) |

**Consequência de ordem, aceita conscientemente.** `price` só pode virar não-nulável **depois**
que o backfill preencher. Isso não é obstáculo — é a força motriz: o contrato só aperta quando
o dado existe, e a partir daí regressão reprova sozinha.

**Regra anti-disfarce (herda da spec §1).** `SEM-FONTE` é para o que a API não dá. Campo não
implementado é `PLANEJADO`, **com milestone dona nomeada**. Marcar não-implementado como
`SEM-FONTE` é como `sales_30d` e `quality_score` sumiram do radar por uma missão inteira.

---

## ADR-C4 — Nome de provider morre no adapter

**Contexto.** Ordem direta do operador: o sistema fala com terceiros por adapter, e o domínio
não conhece Mercado Livre.

**Decisão.** Nenhum tipo de domínio, porta, tabela ou campo de contrato cita vocabulário de
provider. Conceitos são canônicos; a tradução mora no adapter.

| conceito canônico | Mercado Livre traduz de |
|---|---|
| saúde do anúncio | `/item/{id}/performance` |
| motivo de bloqueio | `sub_status[]` + `/moderations/infractions` |
| competição de oferta | `price_to_win` + `boosts[]` |
| comissão da linha | `order_items[].sale_fee × quantity` |
| custo de envio do vendedor | `/shipments/{id}/costs → senders[].cost` |
| motivo de cancelamento | `cancel_detail{group,code,description,requested_by}` |

**Consequência — verificável por string.** `sale_fee`, `price_to_win`, `sub_status`,
`marketplace_fee`, `logistic_type` e afins fora de `adapters/mercado_livre/` **reprovam no
gate**. É grep, não julgamento.

---

## ADR-C5 — DIFAL é cálculo nosso, não campo do provider

**Contexto.** O KPI "DIFAL A PAGAR" em `/pedidos` é constante chumbada que nunca leu payload
nenhum. Decisão do operador: tem que ser calculado de verdade — venda de MG para GO gera
DIFAL, e o sistema tem que saber disso.

**Decisão.** DIFAL é motor de regras do **nosso domínio**. O provider fornece apenas insumos
(UF de destino, valores, vínculo do produto). Regime tributário ratificado pelo operador:
**Regime Normal** (Lucro Real/Presumido).

**Forma.** Entrada = fato canônico de venda. Saída =
`{devido, base, valor, fcp, motivo_se_não_devido, parâmetros_usados}`.

**Os parâmetros usados são gravados junto com o valor.** Alíquota muda por decreto; sem isso
ninguém reproduz em auditoria o número calculado há seis meses.

**Tabelas são dado versionado por vigência, nunca constante em código**: alíquota
interestadual por par origem/destino, alíquota interna por UF × NCM, FCP por UF × NCM.

**ADR-17 aplicado.** DIFAL exige NCM (vem do ERP, precisa de vínculo) e UF de destino. Sem os
dois o valor é **desconhecido, nunca zero**. O KPI declara cobertura:
`R$ X sobre 22 de 38 pedidos — 16 sem base fiscal`. Somatório que finge cobrir tudo é a
mentira mais cara desta tela.

**Decisões de contador, não de engenharia — nenhuma vira código sem ratificação:**
1. Base legal da UF de destino: `billing_info` ou endereço de entrega.
2. Frete entra ou não na base de cálculo.
3. Existe UF em que o marketplace tem responsabilidade de retenção.

Os três entram no motor como **parâmetro de configuração**, não como suposição do
desenvolvedor.

---

## ADR-C6 — Payload cru persistido e reconciliado contra o DTO

**Contexto medido.** Quatro campos vieram `NULL` em 34/34 anúncios porque o DTO do adapter
**não declarava a chave** que o ML manda (`items_multiget_reader.go:136-161` não tem `price`,
`currency_id`, `listing_type_id`, `initial_quantity`). `encoding/json` ignora chave
desconhecida por padrão: **sem erro, sem log, sem sintoma**. Mesma classe, 4 ocorrências,
2 adapters.

**Alternativas descartadas.**
- *`DisallowUnknownFields`*: falha alto, mas o provider adiciona campo sem avisar — a ingestão
  de produção quebraria por um campo que não nos interessa. Frágil contra terceiro.
- *Fixture golden por endpoint*: barato e determinístico, mas só pega o que estava no fixture
  no dia da captura. Cego a campo novo.

**Decisão.** Persistir o payload cru por recurso e rodar reconciliação: chave presente no raw
e **ausente** do DTO vira alerta. Pega campo novo do provider sem ninguém reler documentação.

**Consequência.** `listings.raw` está `NULL` hoje (DEF-12) — deixa de ser dívida e vira
**pré-requisito**.

**Condição de segurança, inegociável.** O payload cru carrega PII: nome, documento e endereço
do comprador, e o endereço completo do vendedor. **Scrub no momento da captura, nunca depois.**
Raw sujo persistido é vazamento permanente; "limpo depois" não existe. `billing_info` cru
nunca persiste.
