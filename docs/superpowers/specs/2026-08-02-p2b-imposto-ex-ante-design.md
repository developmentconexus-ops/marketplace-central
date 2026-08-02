# P2.b — Imposto ex-ante do produto (escopo A)

Data: 2026-08-02 · Missão MIS-007 · Fatia posterior ao P2 · Status: **desenho aprovado, com as emendas abaixo**

---

## ⚠ EMENDA 1 (2026-08-02, posterior ao commit desta spec) — a refutação da §1 era inválida

**Leia isto antes de qualquer coisa nas §1, §4 e §9.** Três afirmações desta spec são falsas e estão corrigidas aqui. O texto original fica abaixo como registro do que se acreditava, **não como instrução**.

### O que estava errado

A §1 e a §9 tratam `299,90 × 13,5% = 40,49` como prova de que a `pricing_difal_rates` errava +35% contra a nota fiscal de 29,99. A conta que produz 40,49 não é a da tabela morta — é a da **matriz atual do ERP**:

```
299,90 × 20,5% − 299,90 × 7% = 61,4795 − 20,993 = 40,4865 ≈ 40,49
```

`ALIQINTDEST` da BA no grupo 122 é **20,5**, que é também a alíquota legal da Bahia desde fev/2024. **Quem estava defasada era a nota 895507**, que usou 17,0. A nota-irmã do dia anterior (895436), mesmo grupo, mesma UF, mesma TOP, mesmo comprador PF, usou **20,5** e reconcilia ao centavo:

```
895436:  189,90 × 20,5% − 189,90 × 7% = 25,64   ✓ (VLRICMS 13,29 ✓)
895507:  299,90 × 17,0% − 299,90 × 7% = 29,99   ✓ (VLRICMS 20,99 ✓)
```

A refutação foi feita contra a única nota errada do par. Não foi erro de método — a nota certa está a um dia de distância e não havia como saber que existia.

### Por que a nota estava errada

`TGFHICM` registra que a linha `UFORIG=13, UFDEST=5, TIPRESTRICAO2='I'/122` foi editada **6 vezes na tarde de 2026-07-20**, entre 14:02:32 e 14:42:15, por `LEANDROTH` e `SUP`, alternando `CODTRIB` entre 60 e 0 e preenchendo `ALIQINTDEST=20,5` e `ALIQSUBTRIB=18`. São as 6 alterações mais recentes da base inteira antes de 31/07. A BA foi a 20,5% em fev/2024 e a matriz só acompanhou em jul/2026 — **29 meses de atraso**. Atraso de cadastro fiscal é crônico nesta base, não excepcional.

Descartadas por medição: redução de base (`BASEDIFAL` = 100% do `VLRNOTA` nas duas notas, `REDBASE` vazio) e erro de leitura da fórmula (ela reconcilia ao centavo nos dois documentos).

### O que isso muda no desenho

**A pergunta não é "lei contra ERP". É matriz do ERP contra documentos do ERP**, que discordam entre si — o pedido 895422 diz `CODTRIB=60, 17,0`; a NF-e do mesmo pedido diz `CODTRIB=0, 20,5`.

**Decisão: a fonte é a MATRIZ.** Satisfaz a diretriz de espelhar o ERP (é dado do ERP, não tabela legal externa), é determinística e auditável, e é a única que existe **ex-ante** — quando precificamos, o documento ainda não foi emitido. Isto não é escolha de gosto: espelhar o documento é logicamente impossível nesta fatia.

**A `pricing_difal_rates` continua morrendo**, mas por outro motivo: não porque errava, e sim porque era uma segunda cópia mantida à mão de um dado que o ERP já tem. O veredito "erra +35%" está **retirado**.

### Números corrigidos do pedido BA de R$ 299,90 (grupo 122, `CODTRIB=0`, `ALIQUOTA=7`, `ALIQINTDEST=20,5`, `FCP=0`)

| componente | valor |
|---|---|
| ICMS | 20,99 |
| DIFAL | **40,49** |
| FCP | 0 |
| PIS/COFINS | 25,85 |
| `imposto` (ICMS+FCP+PIS/COFINS) | 46,84 |
| **margem** | **−6,10 (−2,03%)** |

O pedido **dá prejuízo**. A §1 dizia "margem cai de 69,23 (23%) para ~4,4 (1,5%)" — o número certo é **negativo**. A tela mostra 69,23 hoje.

> ⚠ **Armadilha para quem for escrever teste:** a comissão deste pedido também é **40,49**. Dois componentes diferentes com o mesmo valor. Uma asserção que confunda os dois passa por acidente.

---

## ⚠ EMENDA 2 — DIFAL não é calculável em PR, RS e SC

Medido na matriz do recorte real (8 grupos × 10 UFs que os pedidos reais tocaram), contando só células com `CODTRIB=0` (DIFAL devido):

- **21 células calculáveis** — BA, ES, MA, SP, RJ, PE
- **20 células mudas** — `CODTRIB=0` sem `ALIQINTDEST`, concentradas em **PR, RS e SC**
- **6 sem linha** — grupo 311 em BA, MA, PR, RJ, RS, SC

O ERP **cobra** DIFAL do PR (11 notas com 18,0) tirando o número de um lugar que não está em `TGFICM`, `TSIUFS` nem `TGFPAR`, e que não foi localizado.

**Decisão do operador:** PR/RS/SC saem como **desconhecido explícito com motivo legível**, e o par (produto, UF) vira **pendência apontando a lacuna na matriz do ERP**. Nenhuma tabela fiscal é mantida por nós. O conserto pertence ao Sankhya. Custo aceito: 7 notas em 12 meses (de 74) não precificáveis até alguém preencher a matriz.

Isto **reforça** a §2 (nada de tabela legal própria) em vez de contradizê-la.

## ⚠ EMENDA 3 — separar "se incide" de "quanto incide"

`TGFICM.CODTRIB` responde **se** o DIFAL incide (60 = não, 0 = sim) com cobertura total no recorte real — só o grupo 311 tem 6 buracos. `ALIQINTDEST` responde **quanto** em cerca de metade das UFs. As duas perguntas têm confiabilidade diferente e o código deve tratá-las como fatos separados, nunca como uma leitura só.

---

*(texto original abaixo — histórico, com as correções acima prevalecendo)*

## 1. O problema, medido

`/pedidos` mostra margem para 33 de 38 pedidos. O número está errado, e erra **para cima**.

Rastreamento de um pedido real de R$ 299,90 para a Bahia:

| componente | origem | veredito |
|---|---|---|
| comissão 40,49 | ML `sale_fee` × quantidade | correto |
| frete 23,65 | ML `/shipments/{id}/costs` → `senders[].cost` | correto |
| custo 154,53 | Oracle `TGFCUS.CUSSEMICM`, `DTATUAL <= data do pedido` | correto |
| **imposto 12,00** | 4% de um perfil que **não existe** | **fabricado** |
| **DIFAL 0,00** | zero explícito de um switch nunca ligado | **fabricado** |

`pricing_calc_profiles` tem **0 linhas**. `calc_repository.go:49` responde `pgx.ErrNoRows` com
`NewDefaultCalcProfile()` = SIMPLES 4%, DIFAL desligado. O operador ratificou **Regime Normal**.
O imposto exibido é a alíquota de um regime que ninguém configurou.

Nota fiscal real de mesmo valor e mesmo destino (NUNOTA 895507, BA, base 299,90) cobrou
**ICMS 20,99 + DIFAL 29,99**. Somando PIS/COFINS de saída (~25,83), o imposto verdadeiro é
**~76,8** contra os **12,00** exibidos. A margem cai de 69,23 (23%) para **~4,4 (1,5%)**.

Esta é a classe de defeito do ADR-17 com disfarce novo: no P1 o desconhecido virava `NULL` e
aparecia como `—`; aqui vira **default plausível**. É pior, porque em branco o operador
desconfia e R$ 12,00 ele acredita.

## 2. Escopo

**Dentro:** imposto ex-ante por produto × UF de destino, para venda Mercado Livre a
consumidor final pessoa física, calculado antes da venda, para piso de preço e simulação.

**Fora, com motivo:**

| fora | motivo |
|---|---|
| imposto ex-post (o que a nota cobrou) | subsistema B, fatia própria |
| identidade pedido ↔ nota | subsistema C, fatia própria |
| ressarcimento de ST no custo | decisão do operador 2026-08-02: margem real é da parte pós-pedido. E o erro é **pessimista** — deixa o piso conservador. Dívida **D-18** |
| ERP contra a lei (alíquota devida) | exige tabela legal mantida; é o `pricing_difal_rates` de volta com boa intenção. Feature separada, com dono |
| `GetICMSCeiling` do `/anuncios` | outro conceito (teto), outro consumidor. Suspeito, mas não é desta costura. Dívida **D-19** |
| perfil fiscal como tela de configuração | não há o que configurar: a fonte é o ERP |

## 3. Fonte da verdade

**`METALPRD.TGFICM`, consultada por `TGFPRO.GRUPOICMS`.**

A decisão fiscal já foi tomada pelo Sankhya. Nós **lemos**, não derivamos. Isto não é
preferência de estilo — é o que separa esta fatia do `pricing_difal_rates` que ela mata.

Chave de consulta:

```
interestadual:  UFORIG = 13 (MG)
                UFDEST = <destino>
                TIPRESTRICAO2 = 'I'
                CODRESTRICAO2 = TGFPRO.GRUPOICMS
intra-MG:       posições invertidas — UFDEST = 13, TIPRESTRICAO = 'I', CODRESTRICAO = GRUPOICMS
```

Colunas lidas: `CODTRIB`, `ALIQUOTA`, `ALIQINTDEST`, `PERCICMSFCP`.

### Por que a linha genérica basta

`TIPRESTRICAO` não é uma dimensão — é **tag polimórfica de tipo**, e a letra diz para qual
tabela `CODRESTRICAO` aponta (`O`=TOP, `I`=GRUPOICMS, `G`=grupo de produto, `P`=produto,
`N`=genérica). As "duplicatas" por (grupo, UFDEST) não são conflito: são uma linha genérica
`N/0` mais exceções por TOP.

A precedência entre linhas concorrentes **não foi decodificada** e não é decodificável com o
dado disponível. Isso não nos alcança: **as TOPs 306 e 313 (e-commerce) não têm nenhuma linha
em `TGFICM`**, logo toda venda nossa cai na genérica `N/0`, que é unívoca.

Consequência de escopo: se algum dia uma TOP de e-commerce ganhar linha própria em `TGFICM`,
esta premissa quebra silenciosamente. O detector da §6 cobre isso.

### Por que `ALIQINTDEST` e não `ALIQUFDEST`

As duas existem e divergem (PR 19 × 17, MA 22 × 23). As notas reais reconciliam com
**`ALIQINTDEST`**. `ALIQUFDEST` é a coluna que o leitor de teto do `/anuncios` usa — motivo
da dívida D-19.

## 4. O cálculo

Comprador é **sempre** pessoa física não-contribuinte (venda Mercado Livre). Isso está
ratificado nos dois lados:

- **lei:** `LC 87/96 art. 4º §2º` (redação da LC 190/2022) — destinatário não-contribuinte,
  o **remetente** é o contribuinte do DIFAL; destinatário contribuinte, o destinatário apura.
- **dados:** na TOP 306, PF sem IE (`CLASSIFICMS='C'`) → CODTRIB 0, cobra ICMS + DIFAL, em
  26 de 29 linhas (as 3 exceções são intra-MG). PJ com IE (`'X'`) → CODTRIB 60/10, ICMS zero.

Para não-contribuinte a base do DIFAL é **única**. Base dupla só existe para contribuinte —
não nos alcança.

Por unidade, sobre o valor `V`:

```
CODTRIB = 60   →  ICMS = 0 · DIFAL = 0 · FCP = 0        (cadeia encerrada por ST)
CODTRIB = 0    →  ICMS  = V × ALIQUOTA                   (7% N/NE/CO+ES · 12% S/SE)
                  DIFAL = V × ALIQINTDEST − ICMS
                  FCP   = V × PERCICMSFCP                (só RJ, 2%)
CODTRIB = 10   →  DESCONHECIDO                           (ver §5)
CODTRIB outro  →  DESCONHECIDO

PIS/COFINS     =  V × 9,25% sobre base reduzida (BASERED ≈ 93,1% de BASE)
                  ≈ 8,62% do valor cheio — aproximação conservadora documentada
IPI            =  0, com guard nos 6 produtos com TEMIPIVENDA='S'
```

### O que NÃO subtrair de novo

`TGFCUS.CUSSEMICM` não é "custo sem ICMS". É o resultado de `METAL_CUSTO` e já está líquido
de **ST, IPI e crédito de PIS/COFINS**; no ramo com ST o ICMS próprio **fica dentro**.
Subtrair qualquer um deles outra vez é o erro simétrico do que temos hoje, e igualmente caro.

O que falta na venda é só o que está acima. Em particular o **débito** de PIS/COFINS da saída:
o crédito é da entrada e já está no custo; o débito é da receita e hoje não existe em lugar
nenhum. Somar os dois **não é dupla contagem** — é o não-cumulativo funcionando.

### Aritmética verificada

`DIFAL = BASE × ALIQINTDEST − VLRICMS` reconcilia exato nas notas testadas pelo especialista.
Caso BA: 299,90 × 17% = 50,98 − 20,99 = **29,99**, idêntico à nota.

A nossa `pricing_difal_rates` daria 299,90 × 13,5% = **40,49** — erro de +35%.

## 5. Estados de conhecimento

**Selo por componente, nunca por pedido.** Um pedido pode ter ICMS conhecido e DIFAL
desconhecido, e a tela precisa dizer exatamente isso.

Isto não exige tipo novo: `domain.OrderDecomposition` já tem `ComponentesDesconhecidos
[]string`, e `BuildProfitability` já anula a margem quando qualquer componente é nil. O
produtor existente é estendido; nenhum segundo produtor é criado.

Estados dos 19 pares (produto × UF) que os pedidos reais já exerceram:

| estado | pares | o que a tela mostra |
|---|---|---|
| completo | 14 | todos os componentes com valor, margem calculada |
| DIFAL desconhecido | 3 — 39563/PR, 39587/PR, 39587/RS | ICMS com valor; DIFAL marcado; margem vazia |
| tudo desconhecido | 2 — 15956/RJ, 42194/SP | componentes fiscais marcados; margem vazia |

**Nunca zero e nunca pior caso.** Zero é a mentira otimista que temos hoje. Pior caso foi
medido: superestima o ICMS em **+17,5%** nos nossos pares (real 0,31% do faturamento contra
17,8% no pior caso, 57×) e faria o operador recusar venda lucrativa.

Cada componente desconhecido carrega **motivo legível**, não só o estado:

```
DIFAL: desconhecido — ALIQINTDEST vazio na matriz do ERP para o grupo 122 em PR
ICMS:  desconhecido — grupo 311 sem linha para RJ na matriz do ERP
ICMS:  desconhecido — CODTRIB=10 não caracterizado para comprador pessoa física
```

## 6. Detector de conflitos do ERP

Requisito do operador, 2026-08-02: *"pegar erros e conflitos, ERP não é 100% verdade"*.

Escopo deliberadamente restrito a **inconsistências internas** — detectáveis cruzando o ERP
consigo mesmo, sem tabela legal nenhuma:

| conflito | detecção | caso real |
|---|---|---|
| matriz sem linha | consulta por (grupo, UFDEST) não retorna linha | 15956/RJ, grupo 311 sem RJ |
| coluna vazia | `ALIQINTDEST` nulo com `CODTRIB=0` | 39563/PR, 39587/PR, 39587/RS |
| previsto ≠ realizado | CODTRIB da matriz ≠ CODTRIB da nota emitida do mesmo par | 42194/SP — matriz 0, nota 10 |
| alíquota mudou sem log | mesmo par, notas com `ALIQINTDEST` diferente, sem linha em `TGFHICM` entre as datas | RJ 18,0 em jun/2026 → 19,0 em jul/2026 |
| TOP ganhou linha | `TGFICM` passa a ter linha para TOP 306 ou 313 | quebra a premissa da §3 |

Cada conflito vira **pendência marcada no produto** e **alerta na venda**. Não bloqueia nada:
informa. O quarto caso é o que pega o `19,0` do RJ, e pega **sem saber a lei** — basta ver que
a matriz mudou de comportamento sem registrar a mudança.

**Explicitamente fora:** comparar o ERP com a alíquota legal devida. Exige tabela de alíquotas
mantida por alguém, o que é `pricing_difal_rates` reencarnado. Se for feito, é feature própria
com dono próprio.

## 7. Tempo: congelar, e invalidar por log

`TGFICM` **não é versionada na prática**: 4.601 linhas com `DTINICIOVIGENCIA`/`DTFIMVIGENCIA`
nulos em 100%, `SEQUENCIA=1` em todas, zero chaves multi-versão. É sobrescrita in-place.
Alíquota histórica **não é reconstruível** a partir dela.

Existe log parcial: `TGFHICM`, 99 linhas, com `DHALTER`. As alterações estão acelerando —
3 em 2024, 34 em 2025, **62 até 01/ago/2026**.

Duas consequências obrigatórias:

1. **O imposto congela no momento da venda.** Recalcular pela matriz atual reescreve margem
   histórica sozinha quando um estado muda alíquota. Isto é a mesma regra que o custo já
   segue (`DTATUAL <= data do pedido`) — consistência, não invenção.
2. **Cache invalidado por `MAX(TGFHICM.DHALTER)`, não por TTL.** Um TTL de 24h serve uma
   alíquota velha por até 24h sem saber; o log muda no instante da alteração.

Limite honesto: o log é **incompleto** — o `19,0` de julho apareceu sem alteração logada.
Invalidação por log é melhor que TTL, não é garantia. Registrado como **D-20**.

## 8. Onde encaixa no código

Máximo global aplicado como regra: **estender o produtor existente, nunca criar um segundo
produtor do mesmo número.** Foi exatamente o erro que derrubou a primeira versão do plano P2.

| arquivo | ação |
|---|---|
| `internal_read/adapters/oracle/` | **criar** leitor fiscal por (GRUPOICMS, UFDEST) → CODTRIB + alíquotas. Leitor novo, tabela compartilhada com `icms_ceiling.go`, pergunta diferente |
| `internal_read/ports/` | **estender** com a porta do leitor fiscal |
| `orders/adapters/pricingtax/reader.go` | **reescrever** para consumir o leitor fiscal. Deixa de consumir `ProfileSource` |
| `orders/domain/order_decomposition.go` | **estender** — componentes de ICMS, DIFAL, FCP, PIS/COFINS com selo |
| `orders/application/enrich_service.go` | `resolveTaxes` passa a devolver os componentes selados |
| `pricing/adapters/postgres/calc_repository.go` | **não tocar** |
| `pricing_difal_rates` (tabela + seed 0057) | **morre** — nenhum consumidor após a reescrita acima |

Sobre matar a tabela: a fatia remove o **consumo** (o `RateForUF` sai do caminho de `/pedidos`)
e **não** faz `DROP TABLE`. A tabela pertence ao módulo `pricing`, mesma costura de D-21, e
derrubá-la por migração daqui violaria ADR-04. A migração de remoção é uma fatia do dono de
`pricing`, com a verificação de que `/precos` também deixou de ler. Até lá a tabela fica órfã e
o teste da §9 (`o valor não é 40,49`) é o que impede alguém de religá-la por engano.

### A costura

`calc_repository.go` é do módulo `pricing`; P2.b é do M-06, dono de `orders/**`. Decisão do
operador 2026-08-02: **orders para de consumir o perfil**, e `pricing` não é tocado. ADR-04
respeitado — um escritor por costura, e escrevemos só na nossa.

Consequência aceita: `NewDefaultCalcProfile()` (SIMPLES 4%) sobrevive em `pricing`. Se `/precos`
o consome, aquela tela continua fabricando. Dívida **D-21**.

O pacote `pricingtax` mantém o nome mas perde o significado do comentário de cabeçalho ("adapta
o perfil de cálculo do pricing"). Renomear é ruído de diff numa fatia que já mexe no arquivo
inteiro; o comentário é reescrito, o pacote fica. Se incomodar depois, é rename mecânico.

## 9. Testes

Regra: **teste que não falha antes não é teste.** Cada teste abaixo tem um controle negativo
nomeado que precisa ser visto vermelho antes da implementação.

| teste | controle negativo |
|---|---|
| BA 299,90 → ICMS 20,99 + DIFAL 29,99 | com o código atual dá 12,00 + 0. Reproduzir esse vermelho é o gate |
| `pricing_difal_rates` daria 40,49 | asserção explícita de que o valor **não** é 40,49 — pega regressão para a tabela morta |
| CODTRIB 60 → ICMS, DIFAL e FCP todos zero | injetar CODTRIB 0 e ver a asserção falhar |
| `ALIQINTDEST` nulo → DIFAL desconhecido, ICMS **com valor** | preencher a coluna e ver o selo sumir. Prova que o selo é por componente |
| margem nula quando qualquer componente é desconhecido | preencher todos e ver a margem aparecer |
| desconhecido nunca serializa como `0` | asserção sobre o JSON, não sobre o struct |
| RJ → FCP 2%; SP/BA/PR/RS/ES → FCP ausente | inverter e ver falhar |
| PIS/COFINS de saída entra na conta | remover e ver a margem subir ~8,6% |
| detector: par sem linha na matriz vira pendência | dar linha ao par e ver a pendência sumir |
| detector: `ALIQINTDEST` divergente entre notas sem `TGFHICM` vira pendência | alinhar as notas e ver sumir |

Verificação final em browser real, dirigida: abrir `/pedidos`, achar o pedido BA de 299,90,
confirmar que o imposto **não** é 12,00, que o DIFAL **não** é 0, e que um dos 5 pares selados
mostra motivo legível em vez de número.

## 10. Riscos e limites declarados

| # | risco | tratamento |
|---|---|---|
| R1 | precedência entre linhas concorrentes de `TGFICM` não decodificada | não nos alcança hoje (TOPs 306/313 sem linha); detector §6 avisa se mudar |
| R2 | FCP fora do RJ é *provável, não provado* — BA/PR/RS/ES têm 20–29 linhas em 6 anos. Só SP tem n robusto (206) | aceito. Fechar exige conferir NCM contra lista estadual — trabalho jurídico |
| R3 | log `TGFHICM` incompleto | D-20 |
| R4 | `CODEMP=1` fixo no leitor de custo, vendável é `CODEMP(1,2)` | D-17, pré-existente |
| R5 | 8,62% de PIS/COFINS sobre valor cheio é aproximação da base reduzida | documentado no código, não no comentário — constante nomeada com a medição que a originou |
| R6 | divergência entre matriz do ERP e alíquota legal (RJ) | fora de escopo por decisão; **sinalizado ao operador como pendência fiscal da empresa**, separada do projeto |

## 11. Dívidas criadas

| id | dívida |
|---|---|
| D-18 | ST recuperável dentro de `CUSSEMICM` — custo superestimado, margem pessimista. Empresa ressarce via CFOP 1603 mensalmente |
| D-19 | `GetICMSCeiling` do `/anuncios` usa `ALIQUFDEST` com `MAX()` sobre linhas de especificidade diferente |
| D-20 | invalidação de cache por `TGFHICM` é melhor que TTL, não é garantia |
| D-21 | `NewDefaultCalcProfile()` SIMPLES 4% sobrevive em `pricing` |

## 12. Procedência

Toda afirmação factual sobre o ERP nesta spec veio de medição em `METALPRD` (read-only,
2026-08-02) pela sessão especialista Sankhya, registrada em
`Documents/MNOS/docs/product/sankhya-fiscal-venda.md` e
`Documents/MNOS/docs/product/sankhya-custo-formula.md`.

Afirmações jurídicas com fonte primária citada: `LC 87/96 art. 4º §2º` (LC 190/2022),
`Convênio ICMS 142/2018`, `Convênio ICMS 235/2021`, `LC 214/2025`.

Duas hipóteses minhas foram **refutadas** por medição no curso deste desenho e estão
registradas para que não voltem: (a) `CLASSIFICMS` como primeira dimensão de restrição do
`TGFICM` — refutada por tipo (VARCHAR × NUMBER); (b) o `19,0` do RJ como versão antiga
sobrescrita — refutada pela cronologia (o 19 vem **depois** do 18).
