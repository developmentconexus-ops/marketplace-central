# Verdade fiscal ex-ante — método do ERP, dado nosso, divergência registrada

**Data:** 2026-08-04 · **Missão:** MIS-008 · **Metodologia:** `/mc-planning` (fases 0–6)
**Evidência bruta:** `.mnfs/MIS-008/evidence/fiscal-2026-08-04/` (rodadas B, C, F, I–V, contra o
Oracle vivo `METALPRD`, somente leitura, via lane Docker)
**Regras da simulação:** [`docs/superpowers/specs/2026-08-04-regras-fiscais-simulacao.md`](../specs/2026-08-04-regras-fiscais-simulacao.md)

> **Revisão 3.** A revisão 1 elegeu o `TGFAID` como fonte de cálculo e rebaixou nossa tabela legal a
> mero alerta — estava invertido; a rodada R refutou. A revisão 2 corrigiu o desenho (**método do ERP,
> dado nosso, divergência registrada**). Esta revisão fecha as **regras completas antes de planejar**,
> a pedido do operador: rodadas S–V mediram a matriz de regras (spec acima) e uma pesquisa fechou as
> duas disputas do §1.6 **a favor da nossa tabela**.

---

## Parte 1 — Relatório ao operador

### 1.1 O que estava em disputa

Duas coisas diferentes que eu tinha colapsado numa só:

| | quem está certo | evidência |
|---|---|---|
| **método** de cálculo (diferença simples vs gross-up) | **o ERP** | 316/316 itens com `TIPCALCDIFAL = 0`, `BASEDIFAL = BASE`, `ALIQPARADIFAL = ALIQINTDEST − ALIQUOTA`, `PERCPARTDIFAL = 100` em 783/783 (`roundR.txt` R5) |
| **dado** (qual é a alíquota interna do destino) | **nossa tabela** | rodada R, abaixo |

O nosso código erra o **método**. O ERP erra o **dado**. Nenhum dos dois, sozinho, produz o número
certo.

### 1.2 A regra da base, confirmada ao centavo

Nota **212253 / NUNOTA 895507**, BA, 2026-07-21, produto 41912, base 299,90. O Sankhya grava a base
já reduzida em `TGFDIN.BASERED` (`roundK.txt`, K6):

| CODIMP | imposto | BASE | **BASERED** | alíquota | VALOR |
|---|---|---|---|---|---|
| 1 | ICMS | 299,90 | 299,90 | 7 | **20,99** |
| — | DIFAL (mesma linha) | 299,90 | — | 10 (= 17 − 7) | **29,99** |
| 6 | PIS | 299,90 | **248,92** | 1,65 | **4,11** |
| 7 | COFINS | 299,90 | **248,92** | 7,60 | **18,92** |
| 12/14 | IBS/CBS | 299,90 | 225,89 | 0,1 / 0,9 | 0,23 / 2,03 |

`299,90 − 20,99 − 29,99 = 248,92`. **Base menos ICMS e menos DIFAL — exatamente a tua regra.** E
encadeia: base de IBS/CBS = `248,92 − 4,11 − 18,92 = 225,89`.

Não é uma nota só: contra `BASERED`, em 2026, PIS bate **7.617/7.617** e COFINS **7.627/7.627**
(`roundL.txt` L3). A **forma** dessa regra já está no nosso código
(`pricing/domain/icms.go:177`: `P × (1 − a) − S`). O errado é só o `a`.

### 1.3 A prova de que o cadastro do ERP é que está velho

Eu suspeitava que a variação por produto no `TGFAID` pudesse ser redução legítima (produto com base
reduzida por lei). Medi. **Não é.**

**`PERCREDBASEDEST > 0` em ZERO produtos, nas 26 UFs** (`roundR.txt` R3, coluna `N_COM_RED`). Não
existe redução por produto no catálogo. Logo a variação de alíquota dentro da mesma UF não tem
justificativa fiscal — é cadastro em estados diferentes de atualização.

E a distribuição prova isso de forma gritante. Bahia, catálogo vendável (`roundR.txt` R1):

| ALIQINTDEST | produtos |
|---|---|
| **17** | **10.049** |
| 18 | 1 |
| **20,5** | **1** |

Não é "a BA varia por produto". É **10.049 linhas velhas e 2 corrigidas à mão.** Mesmo grupo fiscal
(122) nos dois produtos que eu tinha comparado, sem redução em nenhum (`roundR.txt` R2).

O mesmo padrão em todo o resto (`roundR.txt` R3): **1 alíquota distinta** em 21 das 26 UFs, e 2 ou 3
exatamente onde alguém corrigiu alguma coisa — BA, DF, MA, PE, RJ.

**E o valor corrigido bate com a nossa tabela.** Onde o ERP foi consertado, ele foi consertado
**para o nosso número**:

| UF | valor velho no ERP | valor corrigido no ERP | nossa tabela |
|---|---|---|---|
| BA | 17 | **20,5** | **20,5** ✓ |
| MA | 17 | **23** | **23** ✓ |
| PE | 17 | **20,5** | **20,5** ✓ |

Confirmação independente: nossa tabela não é chute. Quando o pessoal do fiscal descobre o problema
numa nota e vai arrumar o cadastro, o número que eles põem é o nosso.

**O caso mais eloquente é a tua própria nota** (`roundR.txt` R6, alíquotas efetivamente cobradas):

```
BA | 20,5 | 1 nota | 2026-07-20
BA | 17   | 1 nota | 2026-07-21   <-- a tua, NUNOTA 895507
```

Dois dias seguidos, mesmo estado, alíquotas diferentes. No dia 20 saiu certa; no dia 21 saiu com 17
porque **aquele produto** ainda não tinha sido corrigido. É a correção manual, produto a produto,
acontecendo hoje sem sistema nenhum — exatamente o buraco que tu queres fechar.

### 1.4 Quanto isso custou na tua nota — e para os dois lados

Preço 299,90, MG → BA, produto nacional, alíquota interna correta **20,5**:

| cenário | ICMS | DIFAL | base PIS/COF | PIS/COFINS | **total** | erro |
|---|---|---|---|---|---|---|
| **correto** (método ERP + dado nosso) | 20,99 | **40,49** | 238,42 | 22,05 | **83,53** | — |
| **ERP real** (método certo, dado 17 velho) | 20,99 | 29,99 | 248,92 | 23,03 | **74,01** | **−9,52 recolhido a MENOS** |
| **nosso código** (dado 20,5 certo, gross-up) | 20,99 | 50,93 | 227,98 | 21,09 | **93,01** | **+9,48 superestimado** |

Os dois erram quase a mesma coisa, em direções opostas, e **nenhum dos dois é a verdade**:

- **O ERP subrecolheu 9,52 nessa nota.** Isso não é margem a mais — é passivo fiscal.
- **Nós superestimamos 9,48**, o que empurra a margem para baixo e faz recusar venda boa.

Em **GO** — 725 dos 784 itens medidos, o destino real do negócio — o cadastro do ERP está **certo**
(19, bate com a nossa tabela). Lá só sobra o nosso gross-up: ERP 26,49% do preço contra os nossos
29,05%. **2,56 pontos de preço** em toda venda para GO.

### 1.5 O buraco do nosso lado: a tabela "legal" não tem lei

Se o desenho é "julgar o dado do ERP pela lei", a tabela que julga precisa ser lei. Ela não é ainda.
`migrations/0094_icms_matrix.sql:85` — **23 das 27 linhas** trazem:

```
fonte = 'legislação estadual vigente, sem alteração recente conhecida'
vigente_desde = '2000-01-01'
```

Afirmação sem fonte com cara de fonte. Só **4** linhas citam lei de verdade: AL (Lei 9.776/2025),
DF (Lei 7.326/2023), PR (Lei 21.850/2023), RJ (Lei 10.253/2023 + LC 210/2023).

BA 20,5, MA 23 e PE 20,5 são três dessas 23 — e ganharam confirmação independente pelo próprio ERP
corrigido (§1.3). As outras 20 continuam sem nada. **Isso vira a task de maior valor do plano**, não
uma nota de rodapé: sem fonte, a pendência que a tela mostrar não sustenta uma conversa com o fiscal.

### 1.6 As duas UFs em disputa — resolvidas pela pesquisa, a nosso favor

A revisão 2 dizia "não dá para decidir isso medindo". Estava errado: dava para **pesquisar**. Feito
em 2026-08-04:

| UF | ERP corrigido para | nossa tabela | veredito | fonte |
|---|---|---|---|---|
| **DF** | 18 | **20,0** | **nós** | Lei 7.326/2023, vigente 21/01/2024. O 18 é o valor **pré-2024** |
| **RJ** | 19 (+2 FCP) | **20,0 + 2 FCP = 22** | **nós** | LC 210/2023, vigente 20/03/2024 — página oficial SEFAZ-RJ, atualizada 23/06/2026. O RJ era 18+2 antes disso; **"19" não corresponde a nenhum valor histórico** |

Ou seja: as duas também eram cadastro parcialmente atualizado, não disputa legítima. **Nossa tabela
ganha as quatro** (BA, MA, DF, RJ).

**Ressalva honesta:** BA, MA e DF vieram de fonte secundária (legisweb); só o RJ veio do portal da
própria SEFAZ. Antes de gravar `lei` no banco, o número precisa ser conferido no diário oficial ou no
site da SEFAZ. **Pesquisa não é procedência.**

### 1.6.1 Não existe fonte única — e isso é estrutural

O operador pediu para procurar uma base aberta/online com todas as alíquotas atualizadas. **Não
existe**, e a razão é jurídica: alíquota interna é prerrogativa de cada estado; o CONFAZ só
centraliza o que exige convênio. Medido:

| candidato | veredito |
|---|---|
| `github.com/brasil-js/icms` | **morto** — 5 commits, todos de 2016 |
| IBPT / De Olho no Imposto | resolve outro problema — carga **total aproximada por NCM** (Lei 12.741/2012), não isola ICMS |
| BrasilAPI | sem endpoint fiscal |
| agregadores comerciais | matriz completa, **sem lei e sem data de atualização** — checagem cruzada, nunca fonte |
| CONFAZ | repositório de atos, não tabela |
| **SEFAZ estaduais** | **a fonte** — uma a uma; RJ mantém página oficial de verdade |

Uma exceção boa: a **interestadual** (4/7/12) tem base federal única e estável — **Resolução do
Senado nº 22/1989, alterada pela nº 13/2012**. Essa parte da conta pode ficar em código com a
citação no comentário. É a única.

Consequência para o plano: a Task 3 continua sendo curadoria humana por UF. Nenhuma fonte automática
a dispensa.

### 1.6.2 As regras completas, definidas antes de planejar

O operador pediu para fechar **todas** as regras da simulação antes de terminar o plano. Estão em
[`specs/2026-08-04-regras-fiscais-simulacao.md`](../specs/2026-08-04-regras-fiscais-simulacao.md),
medidas nas rodadas S–V. Resumo do que muda o plano:

- **Uma fórmula só:** `base PIS/COFINS = base − ICMS − DIFAL − ST_anterior`, fechando em
  **6.258 de 6.274 itens de 2026 (99,7%)** (`roundV.txt` V1). É a forma que `icms.go:177` já tem.
- **O FCP NÃO sai da base.** Seis dos 16 que não fecham são RJ com resíduo = exatamente −FCP. Regra
  medida, não suposição — e o nosso código precisa fazer igual.
- **Dentro de MG, 88% é CST 60**: ICMS já retido por ST, **zero ICMS na venda**; só o ST anterior sai
  da base. **Fora de MG, 98% é CST 00**: ICMS + DIFAL, e 422 de 765 **também** com ST anterior.
  A intuição do operador confirmada com número (`roundV.txt` V3).
- **Consequência de desenho:** a simulação precisa do **CST** como entrada, não só da UF. Fonte: a
  matriz que já espelhamos (`icms_matrix_mirror.codtrib`).
- **Crédito/ressarcimento de ST não existe na nota:** `VLRSUBST`, `VLRREPREDST` e `VLRICMSUFDEST` são
  **zero em 100%** das notas (`roundT.txt` T1). O ressarcimento é apuração **mensal** do SPED
  (`TGFC185F`/`TGFEFDCC185`/`TGFEFDFC185` sobre a média móvel `TGFEFDVMRSTDIA`). Ex-ante só cabe
  faixa, nunca valor. Lacuna estrutural, registrada.
- **IBS/CBS já estão nas notas** desde 2025-12-10 (0,1 / 0 / 0,9), ~1% do preço e subindo pela
  transição. Entram na simulação com alíquota em configuração.
- **TOP 313 (entrega futura e-commerce) carrega imposto** — 28/43 itens com ICMS valorado
  (`roundT.txt` T9). Eu tinha assumido que era só pedido. **Entra no escopo.**

### 1.7 Como o sistema vai funcionar

**Calcular como o ERP calcula, com o dado que a lei manda, e registrar a diferença.**

| grandeza | fonte | papel |
|---|---|---|
| método (diferença simples, partilha 100%) | medido no ERP, 783/783 | **como** se calcula |
| `a_inter` (interestadual) | regra da origem MG: 4% importado (`ORIGPROD ∈ {1,2,3,8}`), 12% SP/RJ/PR/SC/RS, 7% resto | já existe e está certo (`icms.go:186`, 184/184) |
| interna do destino — **devido** | **nossa `icms_aliquota_interna`, com lei e data** | **o dado do cálculo** |
| interna do destino — **previsto** | **`TGFAID` espelhado (uf, codprod)** | **só para o gate**, nunca para a margem |
| FCP | `TGFICM.PERCICMSFCP`, já espelhado (RJ 142/142, resto 0) | componente |
| ST | `icms_matrix_mirror` + `TGFEFDVMRSTDIA`, já existem | componente |

Fluxo:

1. **Simulação / pedido.** Calcula com o método do ERP e a alíquota **devida**. É a margem
   verdadeira — a que vale para decidir preço, porque é o que a empresa deve pagar.
2. **Antes de faturar.** Compara `TGFAID(produto, UF)` com o devido. Divergiu → pendência nomeada com
   **CODPROD, UF, valor no cadastro, valor devido, diferença em reais e a lei**. É a linha exata que
   alguém vai corrigir no Sankhya. Hoje isso acontece à mão, uma nota por vez, depois do erro.
3. **Depois de faturar.** Reconcilia contra o `TGFDIN` real. Se saiu com o valor velho, fica
   registrado quanto se recolheu a menos.

### 1.8 Por que o erro para — e onde não para

Para porque:

- **Método e dado deixam de ser a mesma decisão.** Hoje "DIFAL pela lei" virou gross-up no código,
  que não era o que tu pediste. Separados, cada um tem dono: método = medido no ERP; dado = a lei,
  citada.
- **A pendência é acionável.** Não é "algo está errado no fiscal" — é `CODPROD 41912, BA, cadastro
  17, devido 20,5, Lei X, R$ 10,50 nesta nota`. Sem isso a tela vira ruído.
- **A varredura é do catálogo inteiro, não da nota.** 10.049 produtos da BA estão com o valor velho
  agora. O gate mostra isso antes de virar 10.049 notas erradas — hoje só aparece uma por vez,
  depois de emitida.
- **Divergência nunca vira número.** ADR-17: célula ausente ou ambígua sai como desconhecido nomeado.

**Não para, e é honesto dizer:**

- **20 das 27 linhas legais ainda não têm fonte.** Enquanto não tiverem, o gate compara um número
  nosso com um número deles. Task 3 existe por isso, e ela precisa de ti ou do contador — eu não
  invento citação de lei.
- **DF e RJ estão em disputa real** (§1.6), não em "ERP velho".
- **O `TGFAID` não tem trilha de alteração** — 5 colunas, sem data (`roundR.txt` R4). Só dá para
  detectar mudança comparando espelho contra espelho ao longo do tempo. É o que o espelho versionado
  faz, mas o passado anterior ao primeiro sync é perdido (mesma classe da D-28).
- **Ex-ante ≠ ex-post no ST.** Antes da nota só existe a média diária (`TGFEFDVMRSTDIA`); o ERP usa o
  real (`TGFC185F`) na emissão. Estrutural, vira dívida, não se apaga.
- **`TIPCALCDIFAL = 0` também é parâmetro de cadastro.** O ERP tem campo para outros métodos; está
  configurado em base única. Se o contador disser que para consumidor final não-contribuinte cabe
  base dupla (LC 190/2022), o método muda — e aí entra pelo gate, não por hardcode. **Pergunta para o
  fiscal, não para mim.**
- **ES** segue não-calculável (D-54, cadastro contraditório de verdade). **MG interna** por produto
  segue aberta (D-55); o `TGFAID` não cobre MG (é a origem).

---

## Parte 2 — Medição (as 8 perguntas, com `file:line`)

**1. O que já existe que faz parte disto?**
`pricing/domain/icms.go:101 TaxesForItem` monta a cadeia inteira (a_inter → a → ICMS → DIFAL → base
PIS/COFINS → ST). `pricing/domain/difal.go:DifalForUF` **já implementa a diferença simples correta**
(`computeEfetivoPct = max(interna − interestadual, 0)`) sobre `pricing_difal_rates`, com override do
operador. `internal_read/adapters/oracle/icms_matrix.go:53` é o molde exato de extração Oracle →
espelho versionado. `internal_read/adapters/oracle/sync.go` (`sankhyaFiscalSTSQL`) já traz o ST e
bate com o ERP — nada a mudar lá.

**2. Onde o defeito realmente mora?**
`pricing/domain/icms.go:154-158`: o `Quo(num, den)` do gross-up. **Um sítio.** A revisão 1 apontava
um segundo defeito na dimensão da porta; a rodada R desmente — com `N_COM_RED = 0` em 26/26 UFs, não
existe variação legítima por produto, então **a dimensão (UF) da tabela legal está certa** e a porta
não muda.

**3. Quem mais consome esse caminho?**
`TaxesForItem` tem **um** chamador de produção: `orders/adapters/pricingtax/reader.go:83`.
`ICMSCell` é construído em **um** sítio de produção: `reader.go:201`. Porta:
`orders/ports/tax_reader.go:44 TaxesForItems`, chamada de `orders/application/enrich_service.go:424`.
`AliquotaInternaFor` é declarada em `pricing/ports/tax_matrix.go:36` e implementada em
`pricing/adapters/postgres/matrix_reader.go:53`. **O simulador de preço não passa por nada disso.**

**4. O que o contrato já diz?**
`AliquotaInternaFor(ctx context.Context, uf string) (*string, error)` —
`pricing/ports/tax_matrix.go:36`. **Fica como está** (resposta 2). A porta **nova** é a do gate, que é
per-produto por natureza: o cadastro errado do ERP é por produto. Nenhuma operação OpenAPI existente
muda nesta fatia.

**5. Qual é o estado vivo real?**
`icms_aliquota_interna`: 27 linhas vigentes, colunas `uf, aliquota, fcp_embutido, fonte, lei,
vigente_desde, vigente_ate`; **23 com `fonte` placeholder e `vigente_desde = 2000-01-01`**.
`icms_matrix_mirror`: 2.498 vigentes + 2 ambíguas; RJ `perc_fcp > 0` em 142/142, resto 0.
`pricing_difal_rates`: 27. `pricing_calc_profiles`: **0 linhas**. `products_mirror`: 10.638.
`METALPRD.TGFAID`: 1.051.700 linhas / 40.450 produtos / 26 UFs, **5 colunas, sem trilha de
alteração**; no catálogo vendável, 10.051 produtos × 26 UFs, `PERCREDBASEDEST > 0` em **zero**.

**6. O que prova que está quebrado hoje e o que provará consertado?**
Quebrado: produto 41912, UF `BA`, preço 299,90 → nosso caminho devolve DIFAL **50,93**; o devido é
**40,49**; o `TGFDIN` da nota 895507 registra **29,99**. Os três números são diferentes e nenhum
software mostra isso hoje. Consertado: o cálculo devolve 40,49 **e** o gate lista o produto 41912/BA
como pendência de cadastro com Δ 10,50. **Aceite = contagem no banco e confronto com nota emitida,
nunca lane verde.**

**7. Orçamento de custo?**
Zero chamada a provider. Um `SELECT` a mais no Oracle por sync do espelho de previsão (`TGFAID`
restrito aos `CODPROD` de `products_mirror` ≈ 261 mil linhas, contra 1,05 milhão da tabela cheia).
Nenhuma consulta extra por render: o cálculo lê 27 linhas; o gate lê por índice (uf, codprod).

**8. O que quebra em silêncio às 3 da manhã?**
Se o sync do `TGFAID` falhar, o gate para de acusar divergência e **a tela fica verde** — o mesmo
modo de falha do §5.1 do design da MIS-008. Por isso o espelho entra no `sync_state` com entidade
própria e a idade aparece junto da pendência. Um gate silencioso é pior que gate nenhum: ele afirma
"nada divergente" sem ter olhado.

---

## Parte 3 — O que já existe (varredura anti-redundância)

`grep -rn "aliquota_interna\|AliquotaInterna\|ICMSCell\|DifalForUF" apps/server_core`.

| existe hoje | `file:line` | por que não serve como está |
|---|---|---|
| `DifalForUF` — diferença simples correta | `pricing/domain/difal.go` | Fórmula certa sobre a tabela errada (`pricing_difal_rates`, paralela à `icms_aliquota_interna`). **Segunda** implementação de DIFAL do repo. Colapsar as duas é o trabalho; criar uma terceira seria a redundância. |
| `TaxesForItem` | `pricing/domain/icms.go:101` | Estrutura certa, um parâmetro errado. **Editar**, nunca duplicar. |
| `icms_aliquota_interna` (27 linhas) | `migrations/0094_icms_matrix.sql:26` | **Continua sendo a fonte do cálculo** — a dimensão (UF) está certa (R3). Falta procedência em 23 linhas (Task 3). |
| `ICMSMatrixReader` | `internal_read/adapters/oracle/icms_matrix.go` | **Molde** do adapter de previsão: extração crua + espelho versionado + entidade no `sync_state`. Copiar a forma. |
| `sankhyaFiscalSTSQL` | `internal_read/adapters/oracle/sync.go` | Já correto e medido. **Nenhuma mudança.** |
| `syncapp.Scheduler` | seam de sync existente | O job novo registra nele. Segundo ticker seria redundância (erro já cometido em F-A3 Slice B). |

**Único artefato novo:** a tabela de previsão do ERP + seu adapter. Justificativa medida: nenhuma
tabela existente tem a chave (uf, codprod) — `icms_matrix_mirror` é (uf, grupo) e
`icms_aliquota_interna` é (uf), e a divergência de cadastro **é** por produto (BA: 10.049 numa
alíquota, 2 noutras).

---

## Parte 4 — Máximo local vs global

1. **Quantas cópias do conceito existem?** DIFAL: **duas** (`icms.go` gross-up, `difal.go` diferença
   simples) sobre **duas** tabelas. Alíquota interna: **duas** (`icms_aliquota_interna`,
   `pricing_difal_rates`). ≥2 ⇒ o conserto pertence onde convergem.
2. **A causa está uma camada abaixo do sintoma?** Sim, mas menos fundo do que a revisão 1 dizia. A
   margem errada na tela é o `a` errado no domínio. **Não** é a porta — R3 refutou isso.
3. **Já existe seam?** Sim: `syncapp.Scheduler`, o padrão de espelho versionado do
   `icms_matrix_mirror`, e `TaxesForItem`. Todos reusados.
4. **Estou estendendo legado?** Não. `pricing_difal_rates` (IC-04) não cresce: absorve-se o override
   no gate ou remove-se em fatia própria (decisão 4).

**O máximo local seria consertar as 19 linhas da tabela** — e continuaria sem detectar que 10.049
produtos da BA vão sair com 17 no cadastro do ERP. O global é: método certo + dado citado + gate por
produto. Raio declarado: 1 domínio, 1 migração de dados, 1 adapter, 1 job, 1 tela de pendência, 2
testes de integração reescritos. **A porta pública não muda.**

---

## Parte 5 — Decisões que só o operador toma

Nenhuma task começa sem estas.

**Decisão 1 — método: diferença simples, confirmado?**
Medido: `TIPCALCDIFAL = 0`, `PERCPARTDIFAL = 100`, 783/783. Tu disseste que o método do ERP está
certo. **Recomendo adotar** e registrar que o método também é parâmetro de cadastro — se o contador
disser que cabe base dupla (LC 190/2022) para consumidor final, entra pelo gate, não por hardcode.

**Decisão 2 — as 20 UFs sem fonte: quem cita a lei?**
Pesquisado (§1.6.1): **não existe fonte consolidada**. O caminho é percorrer as 26 SEFAZ, uma a uma,
anexando lei + vigência. Não é trabalho de agente sozinho — o número da lei precisa de conferência no
diário oficial. **Recomendo:** eu faço o levantamento por UF e entrego a lista para tu ou o contador
ratificar antes de virar `UPDATE`. **Sem ratificação a Task 3 não fecha.**

**Decisão 3 — DF (18 vs 20) e RJ (19 vs 20+2): quem está certo?**
**Resolvida a nosso favor** (§1.6): DF 20 (Lei 7.326/2023) e RJ 20+2 (LC 210/2023, fonte oficial
SEFAZ-RJ). Os valores do ERP são pré-2024. Falta só tu ratificares — e conferir BA/MA/DF, que vieram
de fonte secundária.

**Decisão 4 — `pricing_difal_rates` (IC-04): o que é, e o que faço com ela.**

*O que é:* uma **segunda tabela de alíquotas**, com 27 linhas (uma por UF), criada antes da
`icms_aliquota_interna`. Guarda a interna do destino e um campo de **override do operador** — um
valor que tu podes fixar à mão para uma UF, sobrepondo o calculado. Quem lê é
`pricing/domain/difal.go DifalForUF`, e **essa função já faz a conta certa** (diferença simples,
`max(interna − interestadual, 0)`) — é o único lugar do repo que sempre esteve certo.

*O problema:* são duas tabelas dizendo a mesma coisa. Duas alíquotas internas para a mesma UF, sem
regra de qual vence. O código de pedidos lê uma; o de DIFAL lê a outra. Divergirem é questão de
tempo, e quando divergirem ninguém vai saber qual estava certa.

*Recomendo:* **absorver o override na tabela legal e apagar a `pricing_difal_rates`, em fatia
própria, logo depois desta.** Razões: (a) o override é útil e não pode sumir — vira coluna em
`icms_aliquota_interna`, com quem alterou e quando; (b) a fórmula de `DifalForUF` fica, é a certa;
(c) apagar na mesma fatia mistura conserto de cálculo com remoção de tabela, e se algo quebrar não
se sabe qual dos dois foi. Fatia própria com aceite por `count(*)`. **Nesta casa "depois" nunca
aconteceu**, então entra como task numerada, não como dívida.

---

## Parte 6 — Tasks

Ordem obrigatória. Cada task commita. Vermelho antes de verde, com o texto da falha esperada escrito
**antes** de rodar.

```bash
cd apps/server_core && GOCACHE="$PWD/.gocache" go test ./internal/modules/pricing/...
```

```bash
cd apps/server_core && GOCACHE="$PWD/.gocache" go test -tags=integration ./tests/integration/...
```

---

### Task 1 — Golden do domínio: método do ERP com dado devido

**Arquivo:** `apps/server_core/internal/modules/pricing/domain/icms_erp_golden_test.go` (novo)

Casos transcritos de nota emitida, com o **dado devido**, não o cobrado:

| caso | entrada | esperado |
|---|---|---|
| nota 895507 (BA), devido | P=299,90, a_inter=7, interna=**20,5**, FCP=0, S=0 | ICMS 20,99 · DIFAL **40,49** · base PC 238,42 · PIS/COFINS 22,05 |
| nota 895507 (BA), reproduzir o ERP | mesma coisa com interna=17 | ICMS 20,99 · DIFAL 29,99 · base PC **248,92** · PIS/COFINS 23,03 — **bate com `TGFDIN` linha a linha** |
| GO típico | P=299,90, a_inter=7, interna=19 | DIFAL 35,99 · base PC 242,92 |
| RJ com FCP | P=299,90, a_inter=12, interna=19, FCP=2 | DIFAL 20,99 · FCP 6,00 (`roundI.txt` I2: 893649 → `VLRDIFALDEST` 20,99) |
| intra-MG CST 00 | P=136,66, interna=18 | ICMS 24,59 · DIFAL **0** · base PC 112,07 (`roundJ.txt` I4: 882236/2) |
| **intra-MG CST 60** (o caso de 88%) | P=110,00, CST=60, ST_ant=8,13 | ICMS **0** · DIFAL **0** · base PC **101,87** (`roundT.txt` T10: 897156/1) |
| **interestadual com ST anterior** | P=143,10, ICMS=10,02, DIFAL=17,17, ST_ant=11,93 | base PC **103,98**, resíduo **zero** (`roundU.txt` U1: 896886/1) |
| **RJ — FCP NÃO abate a base** | P=1.896,00 − 36,60, ICMS=223,13, DIFAL=111,56, FCP=37,19 | abatido = **334,69** (só ICMS+DIFAL); incluir o FCP erra por 37,19 (`roundV.txt` V2: 892328/1) |

O segundo caso fecha o argumento: **com o dado do ERP, o nosso motor tem que reproduzir o `TGFDIN` do
ERP ao centavo.** Se reproduz, o método está provado idêntico, e toda diferença restante é dado.

Os três últimos vêm das rodadas S–V e cobrem os dois cenários que o operador nomeou (dentro de MG só
ST; fora de MG ICMS + DIFAL) mais a exceção do FCP. **CST é entrada do cálculo** — sem ele o mesmo
produto no mesmo preço dá imposto errado por 88% dos itens.

**Controle negativo nomeado:** rodar contra o código de hoje. O caso BA-devido deve falhar com
`DIFAL = 50.93, want 40.49`. Outra mensagem ⇒ o teste está errado, não o código.

- [ ] Escrever os oito goldens.
- [ ] Rodar a lane de unidade e **colar a falha exata** no commit.

---

### Task 2 — Matar o gross-up

**Arquivo:** `apps/server_core/internal/modules/pricing/domain/icms.go:144-178`

```go
// a é a carga total do destino: própria + DIFAL. O ERP calcula o DIFAL por
// DIFERENÇA SIMPLES: TGFDIN.ALIQPARADIFAL = ALIQINTDEST − ALIQUOTA, com
// TIPCALCDIFAL = 0, BASEDIFAL = BASE e PERCPARTDIFAL = 100 em 783/783 itens
// de notas TOP 305/306 desde 2025.
//
// D-41 ("DIFAL pela lei") era sobre o DADO, não sobre o método: o método do
// ERP está certo, o cadastro dele é que está velho. Ver
// .mnfs/MIS-008/evidence/fiscal-2026-08-04/roundR.txt (R5) e o plano
// docs/superpowers/plans/2026-08-04-fiscal-verdade-tgfaid-plan.md §1.1.
a = aInt
```

`ICMS_oper = P × a_inter` e `DIFAL = P × a − ICMS_oper` seguem derivados do mesmo `a` — uma fonte,
duas linhas. A base do PIS/COFINS (`icms.go:175-178`) **não muda**: `P × (1 − a) − S` já é
`P − ICMS − DIFAL − S`.

- [ ] Aplicar.
- [ ] Task 1 passa. Nomear no commit qual golden virou verde.
- [ ] `decompose_icms_golden_test.go` e `icms_test.go` quebram — os esperados foram escritos contra o
      gross-up. **Recalcular e reescrever, nunca apagar**; cada valor novo com a conta no comentário.

---

### Task 3 — Procedência da tabela legal (a task de maior valor)

**Bloqueada pela decisão 2.** Sem fonte, o gate não sustenta conversa com o fiscal.

**Migração:** `apps/server_core/migrations/00NN_icms_interna_procedencia.sql` (prefixo medido contra o
maior existente no momento de escrever + o fixture de contagem em
`internal/platform/migrate/runner_test.go`).

- [ ] Para cada uma das 20 UFs sem fonte: `UPDATE` com `fonte`, `lei` e `vigente_desde` reais, **ou**
      marcar explicitamente como não verificada. Nunca deixar o placeholder passando por citação.
- [ ] `CHECK` que impeça o texto placeholder de voltar.
- [ ] BA 20,5 / MA 23 / PE 20,5 ganham a nota de confirmação independente pelo ERP corrigido
      (`roundR.txt` R6) além da lei.
- [ ] DF e RJ ficam marcadas **em disputa** (decisão 3) — e uma UF em disputa **não calcula
      silenciosamente**: sai como pendência (ADR-17).
- [ ] `TestICMSAliquotaInternaSeedHasAll27UFsExactlyOnce`
      (`migrations/icms_matrix_shape_test.go:105`) ganha asserção de procedência. **Restado, não
      apagado.**

---

### Task 4 — Espelho da previsão do ERP

**Migração:** `apps/server_core/migrations/00NN_erp_aliquota_interna_prevista.sql`

```sql
CREATE TABLE IF NOT EXISTS erp_aliquota_interna_prevista (
    tenant_id     text        NOT NULL,
    uf_destino    text        NOT NULL,
    codprod       integer     NOT NULL,
    aliquota      numeric(6,3) NOT NULL,
    perc_fcp      numeric(6,3),
    red_base_dest numeric(6,3),
    vigente_desde timestamptz NOT NULL,
    vigente_ate   timestamptz,
    synced_at     timestamptz NOT NULL,
    PRIMARY KEY (tenant_id, uf_destino, codprod, vigente_desde)
);
CREATE UNIQUE INDEX IF NOT EXISTS erp_aliq_interna_prevista_vigente
    ON erp_aliquota_interna_prevista (tenant_id, uf_destino, codprod)
    WHERE vigente_ate IS NULL;
```

O nome carrega o papel: **previsão do que o ERP vai fazer**, nunca fonte de cálculo. Quem ler isto
daqui a seis meses não pode confundir de novo — foi o erro da revisão 1 deste plano.

**Adapter:** `internal_read/adapters/oracle/erp_aliquota_interna.go`, no molde de `icms_matrix.go`:

```sql
SELECT u.UF, a.CODPROD, a.ALIQINTDEST, a.PERCICMSFCP, a.PERCREDBASEDEST
FROM METALPRD.TGFAID a
JOIN METALPRD.TSIUFS u ON u.CODUF = a.CODUF AND u.CODPAIS = 55
WHERE a.CODPROD IN (<CODPROD presentes em products_mirror>)
```

≈261 mil linhas em vez de 1,05 milhão. `PERCREDBASEDEST` é lido e **medido zero em 100% do catálogo**
(R3) — se algum dia deixar de ser zero, é sinal de mudança de regime e o gate tem que acusar, não
ignorar.

**Job:** registrar em `syncapp.Scheduler` com entidade própria no `sync_state`.

- [ ] Teste de integração hermético: duas versões (uma fechada, uma vigente) → reader devolve a
      vigente. **Controle negativo:** (uf, produto) sem linha devolve `nil`, nunca default.
- [ ] Rodar e ver falhar antes do adapter existir.

---

### Task 5 — Gate previsto × devido

- [ ] Para cada (produto vendável, UF), comparar `erp_aliquota_interna_prevista` com
      `icms_aliquota_interna`.
- [ ] Pendência traz **CODPROD, UF, cadastro, devido, Δ em reais no preço corrente, e a lei**. Sem
      esses cinco campos ela não é acionável e a task não fecha.
- [ ] Agregar por UF **e** listar produto — 10.049 pendências na BA precisam virar uma linha
      navegável, não 10.049 alertas.
- [ ] **Aceite por observável:** a tela lista BA com 10.049 produtos divergentes (cadastro 17,
      devido 20,5) e nomeia o produto 41912.
- [ ] Falha ou atraso do sync do espelho aparece **na mesma linha**. Gate silencioso afirma "nada
      divergente" sem ter olhado — pior que gate nenhum.

---

### Task 6 — O simulador entra no mesmo seam

**Medido:** `ICMSCell` é construído em **um** sítio de produção — `reader.go:201`, do módulo de
pedidos. O simulador (`pricing/transport/calc_handler.go`) nunca recebe célula; cai em
`decompose.go:136 pctOfPrice(in.AliquotaPct, preco)`, e `AliquotaPct` vem do perfil de regime, cujo
default é `calcprofile.go:19 defaultSimplesAliquotaPct = "4"` — com `pricing_calc_profiles` em
**0 linhas**. **Empresa de regime normal simulando com SIMPLES 4%.**

Gate de sítio de composição: sem esta task, as tasks 1–5 entregam motor correto sem consumidor na
tela que mais se usa.

- [ ] `DecomposeInput` do simulador recebe `ICMSCell` resolvida pela alíquota **devida**.
- [ ] Sem UF escolhida: componente **desconhecido nomeado**, não 4%.
- [ ] **Aceite por observável:** simular o produto 41912 para BA em `/precos` e ler na tela
      ICMS 20,99 / DIFAL 40,49 sobre 299,90.
- [ ] **Precondição do live drive:** container do backend construído de um commit **≥** o desta
      fatia. Comparar `CreatedAt` do container com o commit — binário velho faz o live drive mentir.

---

### Task 7 — Dívidas

**Fechadas:** `D-41` (era sobre o dado, não o método — reescrita, não removida) · `D-45` (a fonte da
previsão existe: `TGFAID`; a fonte legal continua sendo pesquisa humana) · `D-46` (a_inter por
`ORIGPROD` confirmado, 784/784) · fila item 4 · fila item 5.

**Abertas / novas:**

| id | fato medido |
|---|---|
| `D-56` reescrita | nosso gross-up diverge do método do ERP; +9,48 numa nota real de 299,90. |
| **nova** | 20 das 27 linhas legais sem procedência (`0094_icms_matrix.sql:85`). |
| **nova** | DF (18 vs 20) e RJ (19 vs 20+2) — corrigidas no ERP para valor ≠ do nosso. Disputa real. |
| **nova** | ERP subrecolhe DIFAL onde o cadastro está velho: −9,52 na nota 895507; 10.049 produtos da BA na mesma situação hoje. **Exposição fiscal, não margem.** |
| **nova** | `TGFAID` sem trilha de alteração (5 colunas, R4) — passado anterior ao 1º sync é perdido. Mesma classe da `D-28`. |
| **nova** | ex-ante ≠ ex-post no ST: `TGFEFDVMRSTDIA` (média diária) vs `TGFC185F` (real da nota). Estrutural. |
| **nova** | `TIPCALCDIFAL` é parâmetro de cadastro; base dupla (LC 190/2022) não foi avaliada. Pergunta ao contador. |
| **nova** | 7 itens de 6.274 com resíduo positivo na fórmula única, sem causa identificada (`roundV.txt` V2: PE 880915, SP 882116, SC 884972…). Observável que resolve: dump de todos os componentes de uma nota inteira. |
| **nova** | 12 itens com PIS alíquota **0** e COFINS 7,60 normal (`roundT.txt` T8) — não é monofásico (zeraria os dois). Provável erro de cadastro do produto. |
| **nova** | significado de `TIPIMPOSTO` I/S/B/F em `TGFEFDVMRSTDIA` não medido; `F` nunca tem valor em 4,9 milhões de linhas. |
| `D-54` | ES ambíguo de verdade; venda para ES sai não-calculável. |
| `D-55` | MG interna por produto; `TGFAID` não cobre MG. |
| `D-17` | `CODEMP = 1` fixo — inalterado por esta fatia. |

---

### Task 8 — Colapsar `pricing_difal_rates` (fatia própria, **logo após** esta)

**Depende da decisão 4.** Não entra na mesma fatia das tasks 1–6 de propósito: misturar conserto de
cálculo com remoção de tabela torna impossível saber qual dos dois quebrou.

- [ ] Migrar o override do operador de `pricing_difal_rates` para colunas em `icms_aliquota_interna`
      (valor, quem alterou, quando). **O override não pode sumir** — é o instrumento manual que existe.
- [ ] `DifalForUF` (`pricing/domain/difal.go`) **fica** — a fórmula é a certa. Só troca a fonte de
      leitura.
- [ ] `DROP TABLE pricing_difal_rates` em migração própria.
- [ ] **Aceite por observável:** `count(*)` — zero tabelas com alíquota interna além de
      `icms_aliquota_interna`; e o mesmo produto/UF devolve o mesmo DIFAL antes e depois.

---

## Parte 7 — O roadmap, revisto

**Estado medido:** Onda 0 fechada como código (`F-A1`/`F-A2` @`be8fc56c`, `F-A3` @`afb6b54a`);
`F-00` bloqueada por D-16; P2.b na `main` @`65aacc7`. **Onda 1 não começou.**

O design da MIS-008 (§7) trava `F-06` (consolidar margem) e `F-09` (API entrega a linha pronta) em
`D-17`(fiscal). A medição mostra que o travamento não era `CODEMP`: era método errado no nosso código
e dado velho no do ERP, e **nenhum dos dois estava visível em lugar nenhum**.

**Proposta: Onda 0.5 — "verdade fiscal" — entre a Onda 0 e a Onda 1**, com as tasks da Parte 6.

- **`F-06` consolida três motores de margem em um.** Consolidar sobre entrada errada produz **um**
  número errado no lugar de três, e apaga a evidência de que estava errado.
- **`F-09` publica a linha pronta.** Publicar SIMPLES 4% fabricado num contrato é caro de desfazer.
- **Disjunção com a Onda 1 por arquivo:** esta onda toca `pricing/domain`, `pricing/adapters`,
  `internal_read/adapters/oracle`, duas migrações e uma tela de pendência. A Onda 1 toca
  `orders/application/service.go`, `batch_orchestrator.go`, `pedidosFormatters.ts`,
  `mercadoFormatters.ts`, OpenAPI + SDK. **Interseção esperada vazia — mas isso é alegação até rodar
  `git worktree list` e `git diff --name-only main...<branch>` por branch em voo.** Medir antes de
  despachar em paralelo.
- **Ondas 2 e 3 inalteradas.** `F-01`–`F-04`, `F-10`, `F-A4`, `F-A5` são independentes do fiscal.

**O que repensar, não só evoluir:** o design tratava o fiscal como bloqueio de outras fatias. Ele é a
**linha de chegada** de metade delas — margem verdadeira é o produto. E a rodada R mostrou uma coisa
que o design não previa: **existe hoje um processo manual, produto a produto, corrigindo o cadastro
do ERP depois que a nota sai errada.** A Onda 0.5 não é só conserto de cálculo; é a primeira vez que
esse processo ganha instrumento. Sugiro promovê-la a onda própria com aceite em nota emitida.

---

## Parte 8 — Gates de auto-revisão (fase 6)

- [x] Todo `file:line` citado foi aberto nesta sessão.
- [x] Nenhum artefato novo sem citação do que existe e por que não serve (Parte 3).
- [x] Local vs global respondido por escrito (Parte 4) — **e a revisão 1 foi refutada por medição,
      não por opinião** (R3: `N_COM_RED = 0` mata a mudança de dimensão da porta).
- [x] Nenhuma constante de negócio nova hardcoded. Alíquota devida vem da tabela com lei; `a_inter` é
      regra da origem com D-46 medido; o método é medido, não literal.
- [x] Nenhum desconhecido vira `0`/default: UF sem procedência ou em disputa sai como pendência.
- [x] Nenhum mock em seam de integração; aceite é contagem no banco e nota emitida.
- [x] Toda asserção do plano reprova no código de hoje (Task 1 traz a mensagem exata).
- [x] Teste alheio alterado é **reescrito**, nunca apagado (Tasks 2 e 3).
- [x] OpenAPI + SDK: **esta fatia não muda contrato HTTP.** Se a Task 6 expuser UF no request do
      simulador, OpenAPI e `sdk-runtime` entram **no mesmo commit**.
- [x] `tenant_id` presente na tabela nova e no índice vigente.
- [x] Falha do sync novo é visível na tela junto da pendência (Task 5), com o caminho falha → pixel.
- [x] Sem chamada a provider; orçamento na Parte 2 pergunta 7.
- [x] Lanes copiam-e-colam, com diretório.
- [x] Governança: duas migrações novas → prefixos únicos **e** o fixture de contagem em
      `internal/platform/migrate/runner_test.go`. Nenhum módulo novo.
- [x] Nada aqui dá push, reset, revert, stash, clean, instala dependência, despeja ambiente, escreve
      no Oracle, nem faz escrita viva em provider.

**Pendências honestas deste plano:** os prefixos de migração e os números de linha de `icms.go` foram
medidos na árvore de hoje (`main` @`90c9396e`) — confirmar ambos na hora de escrever o código. E as
Tasks 3 e 5 dependem das decisões 2 e 3, que são do operador e do contador, não minhas.
