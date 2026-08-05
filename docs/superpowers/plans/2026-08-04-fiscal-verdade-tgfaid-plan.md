# Verdade fiscal ex-ante — método do ERP, dado nosso, divergência registrada

**Data:** 2026-08-04 · **Missão:** MIS-008 · **Metodologia:** `/mc-planning` (fases 0–6)
**Árvore medida:** `main` @`14360d34`
**Evidência bruta:** `.mnfs/MIS-008/evidence/fiscal-2026-08-04/` (rodadas B, C, F, I–V, contra o
Oracle vivo `METALPRD`, somente leitura, via lane Docker)
**Regras da simulação:** [`docs/superpowers/specs/2026-08-04-regras-fiscais-simulacao.md`](../specs/2026-08-04-regras-fiscais-simulacao.md)

> **Revisão 4.** A revisão 1 elegeu o `TGFAID` como fonte de cálculo — invertido, a rodada R refutou.
> A revisão 2 corrigiu o desenho (**método do ERP, dado nosso, divergência registrada**). A revisão 3
> fechou as **regras completas** (rodadas S–V + pesquisa). Esta revisão é a que o operador pediu
> depois de ratificar as recomendações: **varredura módulo a módulo**, com padrão de código, padrão
> de erro, padrão de acesso a dados, interface entre módulos e seam de frontend medidos um por um,
> antes de escrever as tasks. A varredura achou **nove coisas que as revisões anteriores não
> sabiam** (§1.9) — três delas mudam o tamanho do trabalho, e duas o **encolhem**.

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
(`pricing/domain/icms.go:175-177`: `P × (1 − a) − S`). O errado é só o `a`.

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
`migrations/0094_icms_matrix.sql:85-113` — **23 das 27 linhas** trazem:

```
fonte = 'legislação estadual vigente, sem alteração recente conhecida'
vigente_desde = '2000-01-01'
```

Afirmação sem fonte com cara de fonte. Só **4** linhas citam lei de verdade: AL (Lei 9.776/2025,
`:87`), DF (Lei 7.326/2023, `:92`), PR (Lei 21.850/2023, `:101`), RJ (Lei 10.253/2023 + LC 210/2023,
`:104`).

BA 20,5, MA 23 e PE 20,5 são três dessas 23 — e ganharam confirmação independente pelo próprio ERP
corrigido (§1.3).

**Ratificado pelo operador em 2026-08-04:** os 27 valores vieram da tabela que **ele forneceu e
validou**. Ou seja, o defeito nunca foi "número sem origem" — foi **origem real trocada por uma
citação legislativa falsa**. A correção não é conseguir 20 leis antes de calcular; é o campo `fonte`
**dizer a verdade sobre de onde o número veio**:

| tier | UFs | `fonte` | `lei` |
|---|---|---|---|
| lei citada | AL, DF, PR, RJ | a lei | preenchida |
| lei + confirmação independente do ERP corrigido | BA, MA, PE | a lei + a nota de confirmação (R6) | preenchida |
| **validado pelo operador** | as outras 20 | `'tabela validada pelo operador em 2026-08-04'` | **NULL** |

`lei` **NULL** é informação honesta ("não temos a citação"), não um desconhecido virando default — o
número existe, é validado, e calcula. O que **some** é a frase que fingia ser legislação. Quando o
contador trouxer o número da lei de uma UF, ela sobe de tier com um `UPDATE`, sem migração de desenho.
A tela nomeia o tier, para a pendência nunca alegar mais procedência do que tem.

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

Consequência para o plano: a curadoria por UF (Task B1) é humana. Nenhuma fonte automática a dispensa.

### 1.6.2 As regras completas, definidas antes de planejar

Estão em [`specs/2026-08-04-regras-fiscais-simulacao.md`](../specs/2026-08-04-regras-fiscais-simulacao.md),
medidas nas rodadas S–V. Resumo do que muda o plano:

- **Uma fórmula só:** `base PIS/COFINS = base − ICMS − DIFAL − ST_anterior`, fechando em
  **6.258 de 6.274 itens de 2026 (99,7%)** (`roundV.txt` V1). É a forma que `icms.go:175-177` já tem.
- **O FCP NÃO sai da base.** Seis dos 16 que não fecham são RJ com resíduo = exatamente −FCP. Regra
  medida, não suposição — e o nosso código precisa fazer igual. **Isto tem consequência de código
  maior do que parecia; ver §1.9 F-3.**
- **Dentro de MG, 88% é CST 60**: ICMS já retido por ST, **zero ICMS na venda**; só o ST anterior sai
  da base. **Fora de MG, 98% é CST 00**: ICMS + DIFAL, e 422 de 765 **também** com ST anterior.
- **Consequência de desenho:** a simulação precisa do **CST** como entrada, não só da UF. Fonte: a
  matriz que já espelhamos (`icms_matrix_mirror.codtrib`).
- **Crédito/ressarcimento de ST não existe na nota:** `VLRSUBST`, `VLRREPREDST` e `VLRICMSUFDEST` são
  **zero em 100%** das notas (`roundT.txt` T1). O ressarcimento é apuração **mensal** do SPED. Ex-ante
  só cabe faixa, nunca valor. Lacuna estrutural, registrada.
- **IBS/CBS já estão nas notas** desde 2025-12-10 (0,1 / 0 / 0,9), ~1% do preço e subindo pela
  transição. Entram na simulação com alíquota em configuração.
- **TOP 313 (entrega futura e-commerce) carrega imposto** — 28/43 itens com ICMS valorado
  (`roundT.txt` T9). Já está no escopo do sync da matriz
  (`internal_read/adapters/oracle/icms_matrix.go:33 icmsMatrixTOPs = {306, 313}`).

### 1.7 Como o sistema vai funcionar

**Calcular como o ERP calcula, com o dado que a lei manda, e registrar a diferença.**

| grandeza | fonte | papel |
|---|---|---|
| método (diferença simples, partilha 100%) | medido no ERP, 783/783 | **como** se calcula |
| `a_inter` (interestadual) | regra da origem MG: 4% importado (`ORIGPROD ∈ {1,2,3,8}`), 12% SP/RJ/PR/SC/RS, 7% resto | já existe e está certo (`icms.go:186 aInterFor`, 184/184) |
| interna do destino — **devido** | **nossa `icms_aliquota_interna`, com lei e data** | **o dado do cálculo** |
| interna do destino — **previsto** | **`TGFAID` espelhado (uf, codprod)** | **só para o gate**, nunca para a margem |
| FCP | `icms_aliquota_interna.fcp_embutido` (já semeado: AL 1,0 · RJ 2,0 · SE 1,0) | componente de custo que **não** abate a base |
| ST | `icms_matrix_mirror` + `products_mirror.st_retido_entrada` | componente |

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

- **Método e dado deixam de ser a mesma decisão.** Separados, cada um tem dono: método = medido no
  ERP; dado = a lei, citada.
- **A pendência é acionável.** Não é "algo está errado no fiscal" — é `CODPROD 41912, BA, cadastro
  17, devido 20,5, Lei X, R$ 10,50 nesta nota`.
- **A varredura é do catálogo inteiro, não da nota.** 10.049 produtos da BA estão com o valor velho
  agora. O gate mostra isso antes de virar 10.049 notas erradas.
- **Divergência nunca vira número.** ADR-17: célula ausente ou ambígua sai como desconhecido nomeado.

**Não para, e é honesto dizer:**

- **20 das 27 linhas legais ainda não têm fonte.** Task B1 existe por isso, e ela precisa de ti ou do
  contador — eu não invento citação de lei.
- **O `TGFAID` não tem trilha de alteração** — 5 colunas, sem data (`roundR.txt` R4). O passado
  anterior ao primeiro sync é perdido (mesma classe da D-28).
- **Ex-ante ≠ ex-post no ST.** Estrutural, vira dívida, não se apaga.
- **`TIPCALCDIFAL = 0` também é parâmetro de cadastro.** Se o contador disser que para consumidor
  final não-contribuinte cabe base dupla (LC 190/2022), o método muda — e aí entra pelo gate, não por
  hardcode. **Pergunta para o fiscal, não para mim.**
- **`red_base` nunca foi avaliado.** `icms_matrix_mirror.red_base` é não-nulo em **63%** das linhas
  (`migrations/0094_icms_matrix.sql:47`) e **nenhum cálculo o lê**. Na nota 895507 o ERP não reduziu
  (`BASERED = BASE = 299,90`), mas isso é uma nota. **Vira medição própria (Task D3), não implementação
  às cegas.**
- **ES** segue não-calculável (D-54). **MG interna** por produto segue aberta (D-55).

### 1.9 O que a varredura módulo a módulo achou (novo nesta revisão)

O operador mandou medir tudo antes de planejar. Nove achados; três mudam o tamanho do trabalho.

| # | achado medido | `file:line` | efeito no plano |
|---|---|---|---|
| **F-1** | `DecomposeInput` **já tem o campo `ICMSCell`**, com o caminho D-41 já implementado; o que falta é o `CalcService` preenchê-lo | `pricing/domain/decompose.go:35-75` (campo), `:190-222` (ramo D-41), `application/calc_service.go:271-279` (nunca seta) | **ENCOLHE.** A antiga Task 6 era "levar o motor ao simulador"; é composição + contrato, **zero** trabalho de domínio |
| **F-2** | O **solver** (`margem-alvo → preço`) não tem caminho D-41 nenhum: monta `DecomposeInput` sem `ICMSCell` e calcula o teto por `100 − comissão − aliquota − difal` | `pricing/domain/solve.go:289-296 margemDecompose`, `:102-110 ceilingPct`, `:112-114 difalApplied` | **CRESCE.** Segundo consumidor do imposto fabricado, invisível nas revisões 1–3. Sem ele, `/precos` mostra decomposição certa e preço-alvo errado na mesma tela |
| **F-3** | O FCP está **embutido** na alíquota legal (RJ 22,0 com `fcp_embutido` 2,0; AL 21,5/1,0; SE 20,0/1,0), e `AliquotaInternaFor` **não lê `fcp_embutido`** — logo a base do PIS/COFINS é abatida também pelo FCP, contra a regra medida | seed `0094_icms_matrix.sql:87,104,111`; leitor `pricing/adapters/postgres/matrix_reader.go:53-67` | **CRESCE (pouco).** `a` precisa virar **dois**: `a_custo` (com FCP) e `a_base` (sem). Sem isso, RJ/AL/SE erram a base do PIS/COFINS em `FCP × 9,25%` |
| **F-4** | Operação **interna** (MG→MG) sai com DIFAL ≠ 0: o código deriva `ICMS = P×a_inter` e `DIFAL = P×a − ICMS` mesmo quando destino = origem, produzindo ICMS 7% + DIFAL 11% em vez de ICMS 18% + DIFAL 0 | `pricing/domain/icms.go:141`, `:163-165` | **Defeito novo, mesma fatia da correção do `a`.** A soma está certa, a **quebra** está errada — e é a quebra que aparece na tela |
| **F-5** | `icms_matrix_mirror.aliquota`, `.red_base` e `.perc_fcp` são **escritos e nunca lidos**: `MatrixCell` publica só `{Found, CodTrib, Ambiguo}` | `pricing/ports/tax_matrix.go:12-16`, `adapters/postgres/matrix_reader.go:34-39` | Confirma que a alíquota do cálculo é **só** a nossa tabela legal (desenho correto), e transforma `red_base` em medição própria (Task D3) |
| **F-6** | Existem **dois** conjuntos de UF interestadual-12 em conflito: um **exclui** MG, o outro **inclui** | `pricing/domain/icms.go:12-14 aInterUF12` vs `domain/difal.go:14-16 interestadual12` | O colapso da `pricing_difal_rates` mata o segundo. Um conjunto, com a citação da Res. Senado 22/1989 |
| **F-7** | Matar a `pricing_difal_rates` atravessa **cinco** camadas, não uma: tabela + 3 métodos de repositório + porta congelada + 2 operações OpenAPI + 5 schemas + 3 tipos e 2 métodos de SDK + drawer e mutação no FE + **1 exceção de governança** | `migrations/0056`, `adapters/postgres/calc_repository.go:100-176`, `ports/calc_ports.go:39`, `ports/calc_ports_contract_test.go:55-120`, `openapi:2626-2660,4566-4633`, `sdk-runtime/src/index.ts:1524-1552,2380-2382`, `web/src/pages/precos/DifalDrawer.tsx`, `PricingPage.tsx:129-156`, `contracts/governance/invariants.json:50-56` | Fatia própria (C), com substituto no lugar — não um `DROP` que apaga o instrumento manual do operador |
| **F-8** | Nem `/pricing/decompose` nem `/pricing/solve` carregam UF de destino: os dois usam `PricingCalcInput`, e o destino vem de `profile.DifalDestinoUF` (uma UF fixa por tenant) | `transport/calc_handler.go:287-296 calcInputDTO`, `application/calc_service.go:216-225`, `:443-456 applyDifal`, `sdk-runtime/src/index.ts:1571-1580` | Simular "o mesmo produto para BA e para GO" é **impossível hoje**. Um campo em `PricingCalcInput` serve às duas operações |
| **F-9** | `PricingDecomposition` (DTO, OpenAPI e SDK) não publica `icms_saida`, `pis_cofins` nem `restituicao_st` — que o domínio **já calcula** — enquanto o schema de pedidos **já os publica** com a redação ADR-17 | `transport/calc_handler.go:261-273`, `openapi:4699 PricingDecomposition` vs `openapi:5860-5934` (pedidos), `sdk-runtime/src/index.ts:1582-1594` | **Molde pronto.** Copiar a forma do schema de pedidos, não inventar |

**Leitura conjunta:** o trabalho de domínio é **menor** do que a revisão 3 supunha (F-1), e o de
contrato/FE é **maior** (F-2, F-7, F-8, F-9). É exatamente o tipo de troca que só aparece medindo.

### 1.10 Varredura de máximo global e de código morto (2026-08-05)

Pedido do operador antes de executar: *a arquitetura está certa, ou vai precisar de refactor? o que dá
para **apagar**, não só somar?* Medido no tree e no Postgres de dev.

**O que está saudável — e é a maior parte:**

| verificação | resultado |
|---|---|
| legado VTEX em código vivo | **zero**. Só migração histórica (`0005`), a de remoção (`0081`) e docs. O reset da MIS-001 realmente terminou |
| tabelas órfãs (migração sem nenhuma referência em Go) | **1 em 62** — `pricing_manual_overrides` |
| tabelas cujo único chamador é `_test.go` | **nenhuma** — a classe do P2.b (`icms_matrix_mirror` com 0 linhas) não se repetiu |
| portas suspeitas de estarem mortas | `ports/degrau3.go` (102 linhas) **está viva**: wired em `root.go:966-968` e implementada em `adapters/tarifflive/resolver.go:35`. Falso positivo meu, verificado e descartado |
| anatomia de módulo | consistente nos 20 módulos; o molde de adapter/porta é o mesmo em todo lado |

**Veredito: não há refactor arquitetural pendente.** O problema não é a arquitetura — é **duplicação
concentrada na área fiscal**, que é exatamente o que este plano desfaz.

**O que a varredura achou de novo, e que muda tasks:**

| # | achado | evidência | efeito |
|---|---|---|---|
| **F-10** | **Terceira cópia da tabela de 27 UFs**, hardcoded em Go — `domain/difal_seed.go`, 64 linhas — com **zero consumidor de produção** (só `_test.go`). Já era código morto antes deste plano, e é constante de negócio em código, que o gate de auto-revisão proíbe | `difal_seed.go:12-45`; `grep DifalSeed` fora de testes = nada | C3 passa a nomeá-lo. **−64 linhas** que ninguém tinha visto |
| **F-11** | As duas tabelas de alíquota divergem **exatamente pelo FCP**, em três UFs e só nelas: AL 21,5 vs 20,5 (fcp 1,0) · RJ 22,0 vs 20,0 (fcp 2,0) · SE 20,0 vs 19,0 (fcp 1,0) | `SELECT … FROM icms_aliquota_interna l JOIN pricing_difal_rates d ON d.uf = l.uf WHERE l.aliquota <> d.interna_pct` | **Prova aritmética de F-3**: `pricing_difal_rates.interna_pct` **é** o `a_base` (sem FCP) e `icms_aliquota_interna.aliquota` é o `a_custo` (com FCP). `a_base = aliquota − fcp_embutido` **reproduz a tabela que morre linha a linha** — dropar não perde informação nenhuma. Vira caso de teste em A1 |
| **F-12** | `pricing_manual_overrides` (`0004:11`) é a única tabela órfã do banco: 0 linhas, **0 referências em Go**, e a wiki já a documenta como morta (`wiki/modules/pricing.md:214`) | grep de 26 nomes de tabela sobre `internal/`: 25 aparecem, essa não | Task nova (C4). É o **quarto** mecanismo de override do módulo |
| **F-13** | **42 de 104 métodos do `sdk-runtime` não têm nenhum chamador** em `apps/web/src`, `packages/feature-*` ou `packages/web-query` — 40% da superfície de contrato | lista dos 42 medida em 2026-08-05 | **Fora do escopo deste plano.** Vira dívida com o instrumento certo registrado, não task |

**Contagem que fecha o argumento da decisão 5:** o módulo `pricing` tem **quatro** mecanismos de
override — `pricing_manual_overrides` (morto), `pricing_difal_rates.override_*` (0 usos),
o estado inline do FE (`wiki/modules/pricing.md:262`) e o `icms_aliquota_interna_override` que eu ia
criar. Nenhum deles em uso. Não era uma escolha de modelagem: era um **padrão de acumulação**.

**Armadilha de medição desta varredura, registrada para não se repetir:** `pg_stat_user_tables.
n_live_tup` reportou `listings = 0`; a contagem real é **34**. É estimativa de autovacuum, não fato.
Toda contagem aqui é `count(*)` via `query_to_xml`. E a primeira tentativa de achar operação de
contrato órfã comparou `operationId` do OpenAPI com o FE — mas o SDK é escrito à mão e os nomes
divergem (`listMarketplaceOrders` no OpenAPI, `listOrders` no SDK), então aquela lista era ruído. O
instrumento correto é **método do SDK → chamador**, que é o que produziu os 42.

---

## Parte 2 — Medição (as 8 perguntas, com `file:line`)

**1. O que já existe que faz parte disto?**
`pricing/domain/icms.go:101 TaxesForItem` monta a cadeia inteira (a_inter → a → ICMS → DIFAL → base
PIS/COFINS → ST) e **já tem consumidor de produção** em pedidos.
`pricing/domain/decompose.go:35-75` **já aceita `ICMSCell`** e `:190-222` já implementa o ramo D-41
(F-1). `pricing/domain/difal.go:59 DifalForUF` já implementa a diferença simples correta sobre a
tabela errada. `internal_read/adapters/oracle/icms_matrix.go:52` é o molde exato de extração Oracle →
espelho versionado, e `internal_read/application/icms_matrix_job.go:55` + `composition/root.go:723-747`
são o molde exato de job + relógio (`syncapp.Scheduler`, escopo ERP, 24 h, `RunOnce` no boot).

**2. Onde o defeito realmente mora?**
Três sítios no mesmo arquivo, não um:
- `pricing/domain/icms.go:154-158` — o `Quo(num, den)` do gross-up (o defeito principal).
- `pricing/domain/icms.go:163-165` — a quebra ICMS/DIFAL aplicada também à operação interna (F-4).
- `pricing/domain/icms.go:144-145` — MG a 18% **literal** enquanto `icms_aliquota_interna` tem
  `('MG', 18.0)` (`0094:98`). Constante de negócio hardcoded.

Mais dois fora dele: `pricing/domain/solve.go:289-296` (o solver nunca vê `ICMSCell`, F-2) e
`pricing/application/calc_service.go:271-279` (o simulador nunca preenche `ICMSCell`, F-1).

**3. Quem mais consome esse caminho?**
`TaxesForItem` — **um** chamador de produção: `orders/adapters/pricingtax/reader.go:83`. `ICMSCell` —
**um** sítio de construção de produção: `reader.go:201`. Porta: `orders/ports/tax_reader.go:43-45`,
chamada de `orders/application/enrich_service.go:424`, consumida em `:390-400 resolveProfitability`.
`Decompose` — chamado por `calc_service.go` (decompose e cenários) e por `solve.go:290`.
`DifalForUF` — quatro sítios: `ports/calc_ports.go:39`, `application/calc_service.go:453`,
`transport/calc_handler.go:156` e os testes congelados de `ports/calc_ports_contract_test.go`.

**4. O que o contrato já diz?**
`TaxMatrixReader` (`pricing/ports/tax_matrix.go:25-37`): `CellFor(ctx, tenantID, ufOrigem, ufDestino,
grupoICMS) (MatrixCell, error)` e `AliquotaInternaFor(ctx, uf) (*string, error)` — **esta segunda não
recebe `tenantID`**, porque `icms_aliquota_interna` é tabela **global** (a lei não é por tenant;
`0094:26-35` não tem coluna `tenant_id`). É a assinatura que muda em F-3 (precisa devolver também o
`fcp_embutido`).
HTTP: `/pricing/decompose` (`openapi:2718`, `operationId: pricingDecompose`) e `/pricing/solve`
(`openapi:2747`, `pricingSolveTarget`) recebem o **mesmo** `PricingCalcInput` — que não tem UF (F-8).
`/pricing/difal` (`:2626 listPricingDifal`) e `/pricing/difal/{uf}` (`:2637
putPricingDifalOverride`) são as duas operações que morrem com a tabela (F-7).
`PricingDecomposition` (`openapi:4699`) exige `preco, comissao, taxa_fixa, imposto,
componentes_desconhecidos` e **para em `margem_pct`** — sem os componentes que o domínio calcula (F-9);
o schema de pedidos (`openapi:5860-5934`) já os publica e é o molde.

**5. Qual é o estado vivo real?**
`icms_aliquota_interna`: 27 linhas vigentes; **23 com `fonte` placeholder e `vigente_desde =
2000-01-01`**; `fcp_embutido` ≠ 0 em três (AL 1,0 · RJ 2,0 · SE 1,0).
`icms_matrix_mirror`: 2.498 vigentes + 2 ambíguas; `red_base` ≠ 0 em ~63% das linhas de origem;
alimentado a cada 24 h pelo job de `root.go:728`.
`pricing_difal_rates`: 27 (seed `0057`). `pricing_calc_profiles`: **0 linhas** — logo todo simulador
cai em `domain/calcprofile.go:19 defaultSimplesAliquotaPct = "4"`.
`products_mirror`: 10.638, com as 5 colunas fiscais de `0094:1-19`.
`METALPRD.TGFAID`: 1.051.700 linhas / 40.450 produtos / 26 UFs, **5 colunas, sem trilha**; no catálogo
vendável, 10.051 produtos × 26 UFs, `PERCREDBASEDEST > 0` em **zero**.
Migrações: **83** arquivos (`platform/migrate/runner_test.go:25` e `:64`, o mesmo número em dois
sítios); maior prefixo usado **0096**; **0095 não existe** e **0093 está duplicado**
(`0093_orders_status_details_nullable.sql` e `0093_sync_state_market_queue_entity_split.sql`) sem
entrada em `invariants.json` — **próximo prefixo livre = 0097**.

**6. O que prova que está quebrado hoje e o que provará consertado?**
Quebrado, quatro observáveis distintos:
- produto 41912, `BA`, 299,90 → nosso caminho devolve DIFAL **50,93**; devido **40,49**; `TGFDIN` da
  nota 895507 registra **29,99**;
- mesmo produto, destino `MG` → devolve ICMS 7% + DIFAL 11% em vez de ICMS 18% + DIFAL 0 (F-4);
- destino `RJ` → base do PIS/COFINS abatida em 22% em vez de 20% (F-3);
- em `/precos`, a linha "(−) Imposto" (`DecompositionPanel.tsx:85`) mostra **4% do preço** vindo de
  `defaultSimplesAliquotaPct`, e o solver usa o mesmo 4% no teto.

Consertado: os quatro invertem, **e** o gate lista 41912/BA como pendência com Δ 10,50.
**Aceite = contagem no banco e confronto com nota emitida, nunca lane verde.**

**7. Orçamento de custo?**
Zero chamada a provider (nada aqui toca Mercado Livre; o balde de 900/min de
`connectors/adapters/mercado_livre/resilience_decorator.go` fica intacto). Um `SELECT` a mais no
Oracle a cada 24 h para o espelho de previsão (`TGFAID` restrito aos `CODPROD` de `products_mirror`
≈ 261 mil linhas, contra 1,05 milhão da tabela cheia) — mesmo relógio e mesmo escopo do job da matriz.
Por render: o cálculo lê 27 linhas da tabela legal (uma consulta por UF, já em cache de request no
adapter de pedidos, `reader.go` resolve **uma vez por pedido**, não por item); o gate lê por índice
`(tenant_id, uf_destino, codprod)`.

**8. O que quebra em silêncio às 3 da manhã?**
Se o sync do `TGFAID` falhar, o gate para de acusar divergência e **a tela fica verde** — mesmo modo
de falha do §5.1 do design da MIS-008. Mitigação medida: o job entra no `syncapp.Scheduler` com
entidade própria em `sync/domain/sync_state.go` (o molde é `EntityICMSMatrix` em `:29-34`, com o
`Valid()` em `:45` e o teste de valor de fio em `sync_state_test.go:29-34`), de modo que a falha vira
`sync_state.last_error` e aparece no `SyncHealthCard` de `/integracoes`. O cursor carrega o que o
ciclo **fez**, não que ele ocorreu — molde em `internal_read/application/icms_matrix_job.go:35-40`
(`Cells/Ambiguos/Written/CompletedAt`), porque um ticker que roda com `Cells = 0` tem que ser
distinguível de um que sincronizou. **A idade do espelho aparece na mesma linha da pendência.**

---

## Parte 3 — O que já existe (varredura anti-redundância)

`grep -rn "aliquota_interna|AliquotaInterna|ICMSCell|DifalForUF|TaxesForItem" apps packages contracts`.

| existe hoje | `file:line` | por que não serve como está |
|---|---|---|
| `TaxesForItem` | `pricing/domain/icms.go:101` | Estrutura certa, três defeitos pontuais (§2.2). **Editar**, nunca duplicar. |
| `DecomposeInput.ICMSCell` + ramo D-41 | `pricing/domain/decompose.go:35-75`, `:190-222` | **Já serve.** Falta só o chamador preencher (F-1). Criar um segundo caminho seria a redundância. |
| `DifalForUF` — diferença simples | `pricing/domain/difal.go:59` | Fórmula certa sobre a tabela errada. **Segunda** implementação de DIFAL do repo. Colapsar as duas é o trabalho; criar uma terceira seria a redundância. |
| `icms_aliquota_interna` (27 linhas, global) | `migrations/0094_icms_matrix.sql:26-39` | **Continua sendo a fonte do cálculo** — a dimensão (UF) está certa (R3). Falta procedência em 23 linhas e leitura do `fcp_embutido`. |
| `ICMSMatrixReader` + writer + job + relógio | `internal_read/adapters/oracle/icms_matrix.go`, `adapters/mirror/icms_matrix_writer.go`, `application/icms_matrix_job.go`, `composition/root.go:723-747` | **Molde completo** do espelho de previsão: extração crua → resolução pura → escrita versionada → job com cursor → scheduler com `RunOnce`. Copiar a forma inteira, não só o SQL. |
| `syncapp.Scheduler` | `sync/application/scheduler.go` | O job novo registra nele. Segundo ticker seria redundância (erro já cometido em F-A3 Slice B). |
| `MatrixReader` (padrão de acesso a dados) | `pricing/adapters/postgres/matrix_reader.go:12,31-67` | **Molde** de adapter Postgres: `var _ ports.X = (*T)(nil)`, `*pgxpool.Pool`, `QueryRow`, `::text` em numérico, `pgx.ErrNoRows` → estado honesto. |
| schema de pedidos com componentes fiscais | `contracts/api/marketplace-central.openapi.yaml:5860-5934` | **Molde** do que `PricingDecomposition` precisa publicar (F-9), redação ADR-17 inclusa. |
| `DifalDrawer` (instrumento manual do operador) | `apps/web/src/pages/precos/DifalDrawer.tsx` | Morre com a tabela, **mas o instrumento não pode sumir** — vira drawer de alíquota interna com `lei`/`vigência`/`fonte`, que é estritamente mais útil. |

**Artefatos novos, com justificativa medida:**

| novo | por que nenhum existente serve |
|---|---|
| tabela `erp_aliquota_interna_prevista` | Nenhuma tabela tem a chave `(uf, codprod)`: `icms_matrix_mirror` é `(uf, grupo)` e `icms_aliquota_interna` é `(uf)`. A divergência de cadastro **é** por produto (BA: 10.049 numa alíquota, 2 noutras). |
| ~~tabela `icms_aliquota_interna_override`~~ **CORTADA** | Medido em 2026-08-04: `pricing_difal_rates` tem **27 linhas e 0 overrides**. O caminho existe inteiro (OpenAPI, SDK, drawer, colunas, `CHECK`) e **nunca foi usado**. Além disso, o tier `operador-validado` da tabela legal (B1) **já é** o override, versionado e com procedência — a tabela nova seria o segundo mecanismo para a mesma coisa. Ver decisão 5. |
| campo `uf_destino` em `PricingCalcInput` | Não existe campo de destino em `calcInputDTO` (`calc_handler.go:287-296`); o destino atual é uma UF fixa por tenant no perfil (F-8), o que torna a simulação por UF impossível. |

---

## Parte 4 — Máximo local vs global

1. **Quantas cópias do conceito existem?** DIFAL: **duas** (`icms.go` gross-up, `difal.go` diferença
   simples) sobre **duas** tabelas. Alíquota interna: **duas** (`icms_aliquota_interna`,
   `pricing_difal_rates`). Conjunto interestadual-12: **dois**, em conflito (F-6). Caminho de imposto
   no simulador: **dois** (decompose e solve), e só um deles vai ser consertado se ninguém olhar (F-2).
   ≥2 em quatro conceitos ⇒ o conserto pertence onde convergem.
2. **A causa está uma camada abaixo do sintoma?** Sim. A margem errada na tela é o `a` errado no
   domínio; o "4% de imposto" na tela é uma constante em `calcprofile.go:19` com a tabela de perfis
   vazia. **Não** é a porta `TaxMatrixReader` — R3 refutou a mudança de dimensão.
3. **Já existe seam?** Sim, quatro: `syncapp.Scheduler`, o quarteto de espelho versionado do
   `icms_matrix_mirror`, `TaxesForItem`/`Decompose`, e o padrão de adapter Postgres do `MatrixReader`.
   Todos reusados; nenhum paralelo criado.
4. **Estou estendendo legado?** Não — estou **removendo**: `pricing_difal_rates` e a porta congelada
   IC-04 saem (fatia C), com substituto equivalente no lugar.

**O máximo local seria consertar as 20 linhas da tabela** — e continuaria sem detectar que 10.049
produtos da BA vão sair com 17 no cadastro do ERP, sem consertar o solver, e sem levar nada disso à
tela. O global é: um `a`, uma tabela de alíquota, um caminho de imposto para decompose **e** solve,
gate por produto, e o contrato publicando o que o domínio calcula.

**Raio declarado, arquivo a arquivo:** `pricing/domain/{icms,solve,difal}.go`,
`pricing/ports/{tax_matrix,calc_ports,calc_repository}.go`,
`pricing/adapters/postgres/{matrix_reader,calc_repository}.go`,
`pricing/application/calc_service.go`, `pricing/transport/calc_handler.go`,
`internal_read/{domain,adapters/oracle,adapters/mirror,application}` (arquivos novos),
`sync/domain/sync_state.go`, `composition/root.go`, **3 migrações**, OpenAPI + SDK (mesmo commit),
`apps/web/src/pages/precos/{PricingPage,DecompositionPanel,DifalDrawer}.tsx` + testes ao lado,
`contracts/governance/invariants.json`.

---

## Parte 5 — Decisões (ratificadas pelo operador em 2026-08-04)

> *"Pode seguir com suas recomendação, conferi e concordo com suas taxas e pricing difal rates some."*

| # | decisão | estado |
|---|---|---|
| 1 | **Método = diferença simples** (`a = a_interna`), como o ERP | **RATIFICADA.** Registrado que o método também é parâmetro de cadastro: se o contador disser que cabe base dupla (LC 190/2022), entra pelo gate, nunca por hardcode |
| 2 | **Alíquotas conferidas** — BA 20,5 · MA 23 · PE 20,5 · DF 20 · RJ 20+2 | **RATIFICADA.** E a procedência das 27 linhas também: a tabela foi **fornecida e validada pelo operador**. Logo B1 **não bloqueia cálculo** — troca a citação falsa pela origem real, com `lei` NULL onde não temos o número. Subir de tier (validado → lei citada) é `UPDATE`, tarefa contínua do contador |
| 3 | **DF e RJ resolvidas a nosso favor** | **RATIFICADA.** Deixam de ser "disputa" e viram linha normal da tabela, com lei citada |
| 4 | **`pricing_difal_rates` morre** | **RATIFICADA.** Fatia C, com substituto — ver decisão 5 |

**Decisão 5 — NÃO existe tabela de override. Uma tabela só (REVISADA por medição).**

A revisão anterior defendia duas tabelas com quatro razões. O operador perguntou se isso não era
máximo local. Era — e a resposta certa estava uma pergunta acima, numa medição que eu não tinha feito.
Feita agora, no Postgres de dev:

```
SELECT count(*) AS total, count(override_interna_pct) AS com_override FROM pricing_difal_rates;
 total | com_override
-------+--------------
    27 |            0
```

**O override tem zero uso.** O caminho existe inteiro — `putPricingDifalOverride` no OpenAPI, método
no SDK, `input` por UF no `DifalDrawer`, colunas e `CHECK` no banco — e **nunca foi usado uma vez**.
Construir `icms_aliquota_interna_override` seria **portar uma feature sem usuário**, e a discussão
"uma tabela ou duas" seria uma discussão sobre um artefato que não devia existir.

**E há uma redundância pior por baixo:** depois do B1, a tabela legal tem o tier
`operador-validado` — 20 das 27 linhas são, literalmente, *o número que o operador pôs*. **Isso já é
o override**, permanente, versionado e visível na tela. Uma tabela de override em cima de um valor
já marcado como operador-validado é o **segundo mecanismo para a mesma coisa** — exatamente a classe
que a varredura anti-redundância (fase 2) existe para pegar, e que eu deixei passar porque estava
copiando a *forma* da `pricing_difal_rates` em vez de perguntar se a forma era necessária.

**Desenho final: uma tabela.** `icms_aliquota_interna`, global, bitemporal, com `procedencia`,
`fonte`, `lei` e **`actor`**. Corrigir uma alíquota fecha a linha vigente e abre outra — o mesmo
padrão de escrita versionada do resto do repo, com histórico de auditoria de graça. O que era
"override" vira **caminho de edição da própria tabela**, com quem editou e quando.

Ganhos sobre o desenho de duas tabelas, todos concretos: −1 tabela, −1 migração, sem `LEFT JOIN` de
resolução, sem regra de precedência para escrever e testar, e o override passa a **ter histórico** —
que a tabela separada, do jeito que eu tinha desenhado (PK `(tenant_id, uf)`, último vence), **não
teria**. Num sistema fiscal, "quem mudou a alíquota, para quanto e quando" é material de auditoria.
Meu desenho anterior perdia isso; este ganha.

**O que se perde, dito na cara:** a tabela é global, então editar uma alíquota edita para todos os
tenants. Hoje é um tenant, e **não existe alíquota interna de ICMS legítima por tenant** — a lei é a
mesma para quem vende de MG para a BA. O controle certo é **quem pode editar**, não uma cópia por
tenant. Vira dívida nomeada (D2): quando multi-tenant for real, a edição da tabela legal exige papel
de administrador.

**Ressalva honesta da medição:** `0` é o estado **agora**. Não existe histórico que prove que nunca
houve um override — a tabela não guarda. O que se afirma é o medido: nenhum override vivo.

**Decisão 6 — o drawer (RATIFICADA, ajustada à decisão 5).** `DifalDrawer` mostra hoje
interna/interestadual/efetivo por UF, com um `input` de override que ninguém usou. O substituto
mostra **alíquota interna, FCP embutido, procedência (tier), lei, vigência, fonte, quem editou e
quando**, e o campo de edição escreve **na tabela legal**, subindo o tier para `operador-validado`
com `actor` registrado. Mesma função manual, um mecanismo só, com procedência.

---

## Parte 6 — Os padrões da casa que o executor tem que seguir

Medidos nesta árvore. Um plano que entrega código fora destes padrões custa uma rodada inteira.

### 6.1 Anatomia de módulo

`apps/server_core/internal/modules/<nome>/` com cinco camadas: `domain/` (puro, sem I/O),
`ports/` (as interfaces que o módulo **possui**), `application/` (serviços sobre portas),
`adapters/<tecnologia>/`, `transport/` (HTTP). `contracts/governance/modules.json` é a verdade
checada por máquina.

Dependências permitidas, medidas (`modules.json:19-22`):

```
pricing       -> catalog, internal_read, marketplaces
orders        -> connectors, integrations, internal_read, pricing, product_links, sync
internal_read -> catalog, inventory, sync
```

Consequências para este plano:
- `orders/adapters/pricingtax` importar `pricing/domain` é **permitido** (`orders -> pricing`).
- `pricing` importar `internal_read` é **permitido**, mas **não é preciso**: o espelho de previsão é
  lido pelo módulo que faz o gate, não pelo cálculo.
- **`pricing -> tenant_config` NÃO está na lista** e já é importado hoje
  (`pricing/adapters/postgres/product_fiscal_reader.go:10`), sem exceção registrada. É violação
  **pré-existente**; o checador compara contra a tip da `main`, então o plano **não pode adicionar
  nenhuma ocorrência nova** dessa aresta. Se uma task precisar de `tenant_config` num arquivo novo de
  `pricing`, **para e escala** — não inventa exceção.

### 6.2 Padrão de acesso a dados (Postgres)

Molde: `pricing/adapters/postgres/matrix_reader.go`.

```go
var _ ports.TaxMatrixReader = (*MatrixReader)(nil)   // :12  asserção em tempo de compilação

type MatrixReader struct{ pool *pgxpool.Pool }        // :19  pgxpool, NUNCA database/sql

err := r.pool.QueryRow(ctx, `
    SELECT aliquota::text                             // :56  numérico SEMPRE via ::text
    FROM icms_aliquota_interna
    WHERE uf = $1 AND vigente_ate IS NULL
`, uf).Scan(&aliquota)
if err == pgx.ErrNoRows { return nil, nil }           // :60  ausência = estado honesto, nunca zero
```

Regras que caem daí:
- **`pgxpool`, nunca `database/sql`** — invariante `postgres-driver` / `GOV_POSTGRES_DRIVER`.
- **Numérico sai por `::text`** para nunca passar por `float64`. Dinheiro e alíquota são string
  decimal no domínio (`domain.Money`, `ParseRat`, `FormatRatHalfUp`).
- **`pgx.ErrNoRows` vira estado nomeado** (`Found: false` / `nil`), nunca `0`.
- **Toda consulta com `tenant_id` no predicado** — exceto `icms_aliquota_interna`, que é global por
  desenho (a lei não é por tenant); essa exceção fica dita no comentário da consulta.
- **Escrita versionada** segue `adapters/mirror/icms_matrix_writer.go`: transação, ler as abertas,
  igual → só `synced_at`; diferente → fecha e insere; sumiu → fecha e **não reabre**; lista vazia →
  `ErrEmptyICMSMatrix`, falha fechada (`:23-28`).

### 6.3 Padrão de erro

- **Domínio:** funções puras não devolvem erro; entrada inválida é validada **antes**, na
  `application`. `mustRat` entra em `panic` de propósito e por isso `calc_service` valida decimal
  antes de chamar (`ErrInvalidPrice`/`ErrInvalidRate`, `calc_service.go:16-24`).
- **Application:** sentinelas exportadas, `errors.New("CÓDIGO_EM_CAIXA_ALTA")`, comparadas com
  `errors.Is`.
- **Transport:** `apierror.Write(w, status, code, message, details map[string]any)`
  (`platform/apierror/apierror.go:25`) — envelope `error.code` / `error.message` / `error.details`.
  5xx **nunca** vaza mensagem interna: `calc_handler.go:63-70` troca por `"internal error"`. Sucesso
  sai por `httpx.WriteJSON`.
- **`panic` em código de produção é bloqueado** (`GOV_PRODUCTION_PANIC`, `exception_mode:
  exact-occurrence`). Há **duas** exceções vivas em `pricing`: `decompose.go mustRat` (1 ocorrência,
  `invariants.json:41-48`) e `difal.go computeEfetivoPct` (2 ocorrências, `:49-56`). **Ao apagar
  `computeEfetivoPct`, a exceção correspondente sai no mesmo commit** — exceção que aponta para
  símbolo inexistente é registro podre.

### 6.4 Interface entre módulos

A porta é declarada **no consumidor**, não no produtor, quando o consumidor é quem define a forma:
- `orders/ports/tax_reader.go:43` declara `TaxReader`; quem satisfaz é
  `orders/adapters/pricingtax/reader.go:45`, que por dentro fala com `pricing/domain`.
- `internal_read/application/icms_matrix_job.go:18-27` declara `icmsMatrixResolver` e
  `icmsMatrixApplier` **minúsculos, locais**, exatamente para a camada `application` nunca importar
  `adapters/` (invariante `application-import` / `GOV_APPLICATION_IMPORT`).

Porta opcional é `nil`-segura: `enrich_service.go:409` devolve tudo desconhecido quando o leitor é
`nil`, em vez de entrar em pânico. Todo leitor novo segue esse idioma.

Agregação **tudo-ou-nada** é contrato, não detalhe: um item sem vínculo leva **o pedido inteiro** a
desconhecido (`orders/ports/tax_reader.go:6-10`), porque soma parcial é número inventado.

### 6.5 Seam de contrato (OpenAPI ↔ SDK ↔ handler)

- Fonte: `contracts/api/marketplace-central.openapi.yaml`. Cliente: `packages/sdk-runtime/src/`,
  **escrito à mão**, não gerado.
- **Os dois no mesmo commit** — invariante `api-sdk-atomicity` / `GOV_API_SDK_SPLIT`,
  `exception_mode: none`. Uma task, um commit, sempre.
- Paridade de `pricing` é verificada **só do lado do SDK** (`packages/sdk-runtime/src/index.test.ts`,
  que lê o YAML por caminho relativo ao `cwd`): **não existe** `openapi_contract_test.go` em
  `pricing/transport` — os arquivos são `calc_handler.go`, `calc_handler_test.go`, `http_handler.go`,
  `http_handler_error_test.go`. Se a fatia B criar paridade do lado Go, o idioma é o de
  `market/transport/openapi_contract_test.go` (fatiar string do YAML; não há parser de YAML no repo).
- **O frontend nunca faz `fetch`** — invariante `frontend-fetch` / `GOV_FRONTEND_FETCH`. Tudo pelo
  cliente do SDK, via `useClient()`.

### 6.6 Frontend

- Páginas em `apps/web/src/pages/<área>/`; teste **ao lado** do componente (`X.tsx` / `X.test.tsx`).
- Valor desconhecido renderiza `<UnknownValue hint=… />` de `@marketplace-central/ui`
  (`DecompositionPanel.tsx`), **nunca** `0` nem `—` fabricado.
- Estado de carga/erro: `LoadingState` / `ErrorState` de `@marketplace-central/ui`
  (`DifalDrawer.tsx:48-51`).
- Ancoragem de teste por `data-testid` (`decomposition-panel`, `blocking-sem-custo`,
  `difal-off-warning`, `difal-drawer`, `region-solver`).
- Decimal pt-BR entra e sai por `./ptbrDecimal` (`ptBrRateToDot`).

### 6.7 Migrações e governança

- Prefixo numérico único (`GOV_MIGRATION_PREFIX`). **Próximo livre: `0097`** (maior existente `0096`;
  `0095` não existe; `0093` está duplicado sem exceção registrada — **não repetir o padrão**).
- **Toda migração nova bumpa a contagem canônica em DOIS sítios**:
  `internal/platform/migrate/runner_test.go:25` e `:64` (hoje `83` nos dois). Esquecer isso deixa a
  lane de unidade vermelha por motivo alheio à feature.
- Módulo novo → entrada em `modules.json`. **Este plano não cria módulo.**
- Exceção de governança que deixa de valer sai no **mesmo commit** que a torna obsoleta (§6.3).

### 6.8 Lanes (copiar e colar, com o diretório)

```bash
cd apps/server_core && GOCACHE="$PWD/.gocache" go test ./internal/modules/pricing/...
```

```bash
cd apps/server_core && GOCACHE="$PWD/.gocache" go test -tags=integration ./tests/integration/...
```

```bash
cd apps/web && npx --no-install vitest run
```

```bash
npm run harness:governance -- -BaseSha 14360d34c87a53393e57bf61ad1e597dd55edd18
```

Fatos de lane que o plano respeita: o `cd` **faz parte** da lane de FE; a lane hermética se
autodescobre por `//go:build integration` **nas cinco primeiras linhas** do arquivo; a lane de
governança quer SHA de **40 dígitos** e worktree limpa; **o chip nunca sobe servidor, porta ou
`.env`** — a stack de dev é seam do hub.

---

## Parte 7 — Tasks

Quatro fatias. Cada task commita. Vermelho antes de verde, com o texto da falha esperada escrito
**antes** de rodar.

Ordem obrigatória entre fatias: **A → B → C → D**. Dentro da fatia, a ordem das tasks.

---

### Fatia A — o motor (domínio puro, zero mudança de contrato)

#### A1 — Goldens do domínio: método do ERP com dado devido

**Arquivo:** `apps/server_core/internal/modules/pricing/domain/icms_erp_golden_test.go` (novo)

Casos transcritos de nota emitida, com o **dado devido**, não o cobrado:

| # | caso | entrada | esperado |
|---|---|---|---|
| 1 | nota 895507 (BA), devido | P=299,90 · a_inter=7 · interna=**20,5** · FCP=0 · S=0 | ICMS 20,99 · DIFAL **40,49** · base PC 238,42 · PIS/COFINS 22,05 |
| 2 | nota 895507 (BA), reproduzir o ERP | mesma coisa com interna=**17** | ICMS 20,99 · DIFAL 29,99 · base PC **248,92** · PIS/COFINS 23,03 — **bate com `TGFDIN` linha a linha** |
| 3 | GO típico | P=299,90 · a_inter=7 · interna=19 | DIFAL 35,99 · base PC 242,92 |
| 4 | **RJ — FCP não abate** | P=299,90 · a_inter=12 · interna=**22,0** · fcp_embutido=**2,0** | ICMS 35,99 · DIFAL 30,00 · **FCP 6,00** · base PC **239,91** (`= P × (1 − 0,20)`, **não** `× (1 − 0,22)`) |
| 5 | **intra-MG CST 00** | P=136,66 · destino=**MG** · interna=18 | ICMS **24,60** · DIFAL **0,00** · base PC 112,06 (`roundJ.txt` I4: 882236/2) |
| 6 | **intra-MG CST 60** (o caso de 88%) | P=110,00 · destino=MG · CST=60 · ST_ant=8,13 | ICMS **0** · DIFAL **0** · base PC **101,87** (`roundT.txt` T10: 897156/1) |
| 7 | **interestadual com ST anterior** | P=143,10 · ICMS=10,02 · DIFAL=17,17 · ST_ant=11,93 | base PC **103,98**, resíduo **zero** (`roundU.txt` U1: 896886/1) |
| 8 | **RJ real, FCP fora do abatimento** | P=1.896,00 − 36,60 · ICMS=223,13 · DIFAL=111,56 · FCP=37,19 | abatido = **334,69** (só ICMS+DIFAL); incluir o FCP erra por 37,19 (`roundV.txt` V2: 892328/1) |
| 9 | UF sem linha na tabela legal | interna=nil | `Unknown` contém `icms_saida`, `difal`, `pis_cofins` — **nenhum número** |
| 10 | célula ambígua | `Ambiguo=true` | idem — nunca escolhe candidata |

O caso 2 fecha o argumento: **com o dado do ERP, o nosso motor tem que reproduzir o `TGFDIN` do ERP ao
centavo.** Se reproduz, o método está provado idêntico, e toda diferença restante é dado.

**Controles negativos nomeados** (rodar contra o código de hoje, `main` @`14360d34`):

| caso | falha esperada hoje | o que ela prova |
|---|---|---|
| 1 | `DIFAL = 50.93, want 40.49` | o gross-up |
| 4 | `base PC = 233,92, want 239,91` | o FCP abatendo indevidamente (F-3) |
| 5 | `ICMS = 9.57, want 24.60` **e** `DIFAL = 15.03, want 0.00` | a quebra errada na operação interna (F-4) |

Qualquer outra mensagem ⇒ **o teste está errado, não o código**. Parar e reescrever o teste.

- [ ] Escrever os dez goldens, cada valor com a conta no comentário.
- [ ] Rodar a lane de unidade e **colar as três falhas exatas no commit**.

#### A2 — Matar o gross-up, separar FCP, consertar a operação interna

**Arquivos:** `pricing/domain/icms.go`, `pricing/ports/tax_matrix.go`,
`pricing/adapters/postgres/matrix_reader.go`, `orders/adapters/pricingtax/reader.go`

Três correções e uma extensão de porta, num commit só porque um golden (caso 4) só fica verde com as
três.

**(a) O `a` deixa de ser gross-up** (`icms.go:154-158`):

```go
// a é a carga interna do destino. O ERP calcula o DIFAL por DIFERENÇA
// SIMPLES: TGFDIN.ALIQPARADIFAL = ALIQINTDEST − ALIQUOTA, com
// TIPCALCDIFAL = 0, BASEDIFAL = BASE e PERCPARTDIFAL = 100 em 783/783
// itens de notas TOP 305/306 desde 2025. D-41 ("DIFAL pela lei") era sobre
// o DADO, não sobre o método. Ver
// .mnfs/MIS-008/evidence/fiscal-2026-08-04/roundR.txt (R5) e o plano
// docs/superpowers/plans/2026-08-04-fiscal-verdade-tgfaid-plan.md §1.1.
aCusto = new(big.Rat).Quo(mustRat(*cell.AliquotaInterna), cem)
```

**(b) O FCP sai do abatimento, mas fica no custo** (novo, F-3). Duas alíquotas derivadas de uma fonte:

```go
// aCusto inclui o FCP (a tabela legal embute o FCP na alíquota headline
// onde ele é geral — D-43; RJ 22,0 = 20 ICMS + 2 FECP, migrations/0094:104).
// aBase EXCLUI o FCP porque, medido, o ERP abate ICMS + DIFAL da base do
// PIS/COFINS e DEIXA O FCP DENTRO: 6 itens do RJ com resíduo exatamente
// igual a −FCP (roundV.txt V2). Uma fonte, duas alíquotas, nunca dois
// cálculos independentes.
aBase = new(big.Rat).Sub(aCusto, fcp)
```

`base = P × (1 − aBase) − S`. O componente `FCP = P × fcp` passa a existir em `TaxComponents`.

**(c) Operação interna não tem DIFAL** (`icms.go:163-165`, F-4):

```go
// Destino = origem é operação INTERNA: o ICMS de saída é a alíquota
// interna cheia e o DIFAL é 0 EXPLÍCITO (zero legítimo, não desconhecido).
// A derivação interestadual (ICMS = P×a_inter, DIFAL = P×a − ICMS) só vale
// quando há duas UFs; aplicada a MG→MG ela devolvia ICMS 7% + DIFAL 11%,
// soma certa e quebra errada — e é a quebra que aparece na tela.
if cell.UFDestino == ufOrigemMG {
    icmsOper = new(big.Rat).Mul(p, aCusto)
    difal    = new(big.Rat)
} else { … }
```

**(d) MG deixa de ser literal** (`icms.go:144-145`): `18/100` sai; MG lê `icms_aliquota_interna` como
qualquer outra UF (a linha existe, `0094:98`). Constante de negócio em código é gate de auto-revisão.

**(e) A porta passa a devolver o FCP.** `AliquotaInternaFor` hoje devolve `*string`
(`ports/tax_matrix.go:36`); passa a devolver a alíquota **e** o `fcp_embutido`, os dois `nil` juntos
quando a UF não tem linha vigente:

```go
// AliquotaInterna é o par (alíquota headline, FCP embutido) de
// icms_aliquota_interna para uma UF. Os dois vêm da MESMA linha vigente:
// nil significa "UF sem linha" (D-37) e nunca "FCP zero".
type AliquotaInterna struct {
    AliquotaPct    string
    FcpEmbutidoPct string
}
AliquotaInternaFor(ctx context.Context, uf string) (*AliquotaInterna, error)
```

O adapter passa a ler `SELECT aliquota::text, fcp_embutido::text` (`matrix_reader.go:55-59`) e o
comentário `:51-52` ("fcp_embutido is not read here — D-43") **é reescrito com a medição que o
refuta**, não apagado.

- [ ] Aplicar (a)–(e).
- [ ] A1 passa. **Nomear no commit quais goldens viraram verdes.**
- [ ] `decompose_icms_golden_test.go` e `icms_test.go` quebram — os esperados foram escritos contra o
      gross-up. **Recalcular e reescrever, nunca apagar**; cada valor novo com a conta no comentário.
- [ ] `orders/adapters/pricingtax/reader.go` acompanha a nova assinatura (é o único chamador de
      produção, `reader.go:83`/`:201`); a resolução **uma vez por pedido** não muda.

#### A3 — O solver entra no mesmo caminho

**Arquivos:** `pricing/domain/solve.go`, `pricing/domain/solve_golden_test.go`

**Medido (F-2):** `margemDecompose` (`solve.go:289-296`) monta `DecomposeInput` **sem `ICMSCell`**, e
`ceilingPct` (`:102-110`) calcula o teto por `100 − comissão − aliquota − difal`. Sem esta task,
`/precos` mostra decomposição correta e preço-alvo com 4% fabricado **na mesma tela**.

- [ ] `SolveInput` ganha `ICMSCell *ICMSCell`, repassado em `margemDecompose`.
- [ ] `ceilingPct` passa a derivar a carga fiscal da célula quando ela existe; com célula
      desconhecida, o solver **bloqueia** (`Desconhecidos`), nunca supõe alíquota.
- [ ] **Vermelho primeiro:** golden novo "alvo 20% para BA" que hoje devolve preço calculado com 4% e
      passa a devolver o preço calculado com 20,5.
- [ ] `difalApplied()` (`:112-114`) some junto com o caminho legado na fatia C — aqui só ganha o ramo
      novo, para a fatia A não depender de C.

---

### Fatia B — o dado (procedência, gate e contrato)

#### B1 — Procedência da tabela legal

**Decisão 2 ratificada:** os 27 valores são a tabela que o operador forneceu e validou. Esta task
**não bloqueia cálculo nenhum** — ela troca a citação legislativa falsa pela origem real e torna o
tier legível na tela (§1.5).

**Migração:** `apps/server_core/migrations/0097_icms_interna_procedencia.sql`

```sql
ALTER TABLE icms_aliquota_interna
  ADD COLUMN IF NOT EXISTS procedencia text NOT NULL DEFAULT 'operador-validado'
    CHECK (procedencia IN ('lei-citada', 'lei-citada-confirmada-erp', 'operador-validado')),
  ADD COLUMN IF NOT EXISTS actor text;   -- quem editou; NULL nas linhas que vieram do seed
```

`actor` entra aqui, e não numa migração da fatia C, porque é a **mesma** tabela: decisão 5 mata a
tabela de override e o que era "override" vira edição versionada desta linha (C1).

- [ ] As 20 linhas com o placeholder (`0094:86,88,89,90,91,93,94,95,96,97,98,99,100,102,103,105,106,
      107,108,109,110,111,112`): `fonte` = `'tabela validada pelo operador em 2026-08-04'`, `lei` =
      **NULL**, `procedencia` = `'operador-validado'`.
- [ ] AL/DF/PR/RJ → `'lei-citada'`. BA/MA/PE → `'lei-citada-confirmada-erp'`, com a nota da rodada R6
      (o ERP corrigido à mão convergiu para o nosso número) somada à lei.
- [ ] `CHECK` que impeça o texto `'legislação estadual vigente, sem alteração recente conhecida'` de
      voltar. Placeholder que finge ser lei é a falsidade que esta task existe para apagar.
- [ ] **`vigente_desde` das 20 NÃO é atualizada.** Ela faz parte da PK `(uf, vigente_desde)`
      (`0094:38`), então mexer nela é `DELETE`+`INSERT`, não `UPDATE` — e não temos data legislativa
      medida para pôr no lugar. Fica `2000-01-01`, e é `procedencia` que diz que essa data não é uma
      vigência apurada. **A tela nunca renderiza "vigente desde 2000" para tier
      `operador-validado`** — renderiza "vigência não informada".
- [ ] O tier viaja até a tela: o drawer (C2) e a pendência do gate (B3) nomeiam a procedência, para a
      pendência nunca alegar mais do que tem.
- [ ] Subir de tier é `UPDATE` de `lei` + `procedencia`, sem migração de desenho — é a tarefa contínua
      do contador, fora do caminho crítico deste plano.
- [ ] `TestICMSAliquotaInternaSeedHasAll27UFsExactlyOnce`
      (`migrations/icms_matrix_shape_test.go`) ganha asserção de procedência. **Restado, não apagado.**
- [ ] **Vermelho primeiro:** teste que assere `procedencia` das 27 linhas e **zero** ocorrências do
      texto placeholder — hoje falha com 23 ocorrências.
- [ ] **Bumpar `runner_test.go:25` e `:64` de 83 para 84.**

#### B2 — Espelho da previsão do ERP

**Migração:** `apps/server_core/migrations/0098_erp_aliquota_interna_prevista.sql`

```sql
CREATE TABLE IF NOT EXISTS erp_aliquota_interna_prevista (
    tenant_id     text         NOT NULL,
    uf_destino    text         NOT NULL,
    codprod       integer      NOT NULL,
    aliquota      numeric(6,3) NOT NULL,
    perc_fcp      numeric(6,3),
    red_base_dest numeric(6,3),
    vigente_desde timestamptz  NOT NULL,
    vigente_ate   timestamptz,
    synced_at     timestamptz  NOT NULL,
    PRIMARY KEY (tenant_id, uf_destino, codprod, vigente_desde)
);
CREATE UNIQUE INDEX IF NOT EXISTS erp_aliq_interna_prevista_vigente
    ON erp_aliquota_interna_prevista (tenant_id, uf_destino, codprod)
    WHERE vigente_ate IS NULL;
```

O nome carrega o papel: **previsão do que o ERP vai fazer**, nunca fonte de cálculo. Quem ler isto
daqui a seis meses não pode confundir de novo — foi o erro da revisão 1 deste plano.

**Quarteto no molde do `icms_matrix` (§3), arquivo por arquivo:**

| novo | molde |
|---|---|
| `internal_read/adapters/oracle/erp_aliquota_interna.go` | `oracle/icms_matrix.go:52-121` |
| `internal_read/adapters/mirror/erp_aliquota_interna_writer.go` | `mirror/icms_matrix_writer.go:82-187` |
| `internal_read/application/erp_aliquota_interna_job.go` | `application/icms_matrix_job.go:18-92` |
| entidade `EntityERPAliquotaPrevista` | `sync/domain/sync_state.go:29-34` + `Valid()` `:45` + teste de valor de fio `sync_state_test.go:29-34` |
| registro do relógio | `composition/root.go:723-747` (scheduler próprio, escopo ERP, `RunOnce` antes do `Start`) |

```sql
SELECT u.UF, a.CODPROD, a.ALIQINTDEST, a.PERCICMSFCP, a.PERCREDBASEDEST
FROM METALPRD.TGFAID a
JOIN METALPRD.TSIUFS u ON u.CODUF = a.CODUF AND u.CODPAIS = 55
WHERE a.CODPROD IN (<CODPROD presentes em products_mirror>)
```

≈261 mil linhas em vez de 1,05 milhão. `TSIUFS` com `CODPAIS = 55` é obrigatório (o `CODUF` do
Sankhya **não** é IBGE). `PERCREDBASEDEST` é lido e **medido zero em 100% do catálogo** (R3) — se
algum dia deixar de ser zero, é sinal de mudança de regime e o gate tem que acusar, não ignorar.

- [ ] Escritor **falha fechada** com lista vazia (molde `ErrEmptyICMSMatrix`,
      `icms_matrix_writer.go:23-28`): aplicar vazio fecharia toda linha aberta do tenant.
- [ ] Cursor registra `Linhas/Divergentes/Written/CompletedAt` — ticker com `Linhas=0` tem que ser
      distinguível de sync real (molde `icms_matrix_job.go:35-40`).
- [ ] Teste de integração hermético (`//go:build integration` nas **cinco primeiras linhas**,
      `tests/integration/`): duas versões (uma fechada, uma vigente) → leitor devolve a vigente.
      **Controle negativo:** `(uf, produto)` sem linha devolve `nil`, nunca default.
- [ ] Rodar e ver falhar **antes** do adapter existir.
- [ ] **Bumpar `runner_test.go` de 84 para 85.**

#### B3 — Gate previsto × devido

- [ ] Para cada (produto vendável, UF), comparar `erp_aliquota_interna_prevista` com
      `icms_aliquota_interna` (única fonte da alíquota devida, §Decisão 5).
- [ ] Pendência traz **CODPROD, UF, cadastro, devido, Δ em reais no preço corrente, e a lei**. Sem
      esses cinco campos ela não é acionável e a task não fecha.
- [ ] Agregar por UF **e** listar produto — 10.049 pendências na BA precisam virar uma linha
      navegável, não 10.049 alertas.
- [ ] **Falha ou atraso do sync do espelho aparece na MESMA linha.** Gate silencioso afirma "nada
      divergente" sem ter olhado — pior que gate nenhum. Caminho falha → pixel: job devolve erro →
      `syncapp.Scheduler.RecordFailure` → `sync_state.last_error` → `SyncHealthCard` em
      `/integracoes` → e a idade do espelho ao lado da contagem de pendências.
- [ ] **Aceite por observável:** a tela lista BA com 10.049 produtos divergentes (cadastro 17, devido
      20,5) e nomeia o produto 41912. **`count(*)` no banco, nunca lane verde.**

#### B4 — O contrato publica o que o domínio calcula (OpenAPI + SDK + handler, **um commit**)

**Medido (F-8, F-9):** `PricingCalcInput` não tem UF e serve `/pricing/decompose` **e**
`/pricing/solve`; `PricingDecomposition` para em `margem_pct`. O molde do que falta já existe no
schema de pedidos (`openapi:5860-5934`).

- [ ] `PricingCalcInput` ganha `uf_destino` (nullable). Ausente ⇒ cai no `profile.difal_destino_uf`
      atual, preservando o comportamento de hoje; presente ⇒ vence.
- [ ] `PricingDecomposition` ganha `icms_saida`, `difal`, `fcp`, `pis_cofins`, `restituicao_st`, todos
      nullable, com a **mesma redação ADR-17** do schema de pedidos. `imposto` fica marcado como legado
      em retirada — como já está feito em `openapi:5860-5934`.
- [ ] `calc_handler.go:261-273 decompositionDTO` publica os campos novos;
      `calc_handler.go:287-296 calcInputDTO` lê `uf_destino`.
- [ ] `calc_service.go:271-279` **preenche `DecomposeInput.ICMSCell`** — que já existe
      (`decompose.go:35-75`, F-1) — resolvendo a célula pela UF do request. Sem UF resolvível:
      componente **desconhecido nomeado**, nunca 4%.
- [ ] `packages/sdk-runtime/src/index.ts:1571-1594` acompanha, **no mesmo commit**
      (`GOV_API_SDK_SPLIT`, `exception_mode: none`).
- [ ] `packages/sdk-runtime/src/index.test.ts` ganha o caso de paridade dos campos novos.
- [ ] **Vermelho primeiro:** teste de handler que pede `/pricing/decompose` com `uf_destino: "BA"` e
      espera `icms_saida = "20.99"` — hoje falha com `campo ausente`.

#### B5 — A tela mostra a verdade

**Arquivos:** `apps/web/src/pages/precos/{PricingPage,DecompositionPanel}.tsx` + testes ao lado

- [ ] `PricingPage.tsx:220-232` passa a mandar `uf_destino` no `decomposeInput` (seletor de UF ao lado
      da modalidade; default = `profile.difal_destino_uf`).
- [ ] `DecompositionPanel.tsx:85` — a linha `(−) Imposto` dá lugar a **`(−) ICMS`, `(−) DIFAL`,
      `(−) FCP`, `(−) PIS/COFINS`, `(+) Restituição ST`**, cada uma com `<Value>` que renderiza
      `<UnknownValue>` quando nula. **Nenhum 0 fabricado.**
- [ ] **Vermelho primeiro:** o teste ao lado assere o **valor** (`20,99`), não a presença do rótulo —
      asserção de presença passa com qualquer número.
- [ ] **Aceite por observável (live drive):** simular o produto 41912 para **BA** em `/precos` e ler na
      tela ICMS **20,99** / DIFAL **40,49** sobre 299,90; trocar para **MG** e ler ICMS **24,60** /
      DIFAL **0,00**; trocar para **RJ** e ver a linha de **FCP** aparecer.
- [ ] **Precondição do live drive:** container do backend construído de um commit **≥** o desta fatia.
      Comparar `CreatedAt` do container com o commit — **binário velho faz o live drive mentir**.

---

### Fatia C — colapsar `pricing_difal_rates` (fatia própria, decisão 4)

Fatia separada de propósito: misturar conserto de cálculo com remoção de tabela torna impossível saber
qual dos dois quebrou. **Substituto primeiro, remoção depois** — nesta ordem, dentro da fatia.

#### C1 — Caminho de edição da tabela legal (substitui o override, decisão 5)

**Sem migração** — as colunas `procedencia` e `actor` já entraram no `0097` (B1). Esta task é escrita
versionada, não schema.

- [ ] **Precondição medida, reconferir antes de começar:**
      `SELECT count(override_interna_pct) FROM pricing_difal_rates` deve continuar **0**. Se alguém
      tiver usado o override entre o plano e a execução, **para** — a premissa da decisão 5 caiu e a
      task muda. Medido em 2026-08-04: `27 total / 0 com override`.
- [ ] `UpdateAliquotaInterna(ctx, uf, aliquota, actor, now)` na `application`: **fecha** a linha
      vigente (`vigente_ate = now`) e **insere** a nova com `procedencia = 'operador-validado'`,
      `fonte = 'editado por <actor> em <data>'`, `lei = NULL`, `actor`. Mesmo padrão de escrita
      versionada de `adapters/mirror/icms_matrix_writer.go:110-160` — histórico de auditoria de graça,
      que a tabela de override descartada não teria.
- [ ] `AliquotaInternaFor` continua lendo **uma** linha (`vigente_ate IS NULL`) e passa a devolver
      `Procedencia` e `Actor` junto — um número sem procedência é o problema que este plano existe para
      resolver.
- [ ] **Vermelho primeiro:** teste de integração que edita a alíquota da BA de 20,5 para 21, lê a
      vigente (21, `operador-validado`, actor) **e** conta **2** linhas para a BA — a fechada e a
      vigente. Um `UPDATE` no lugar do fecha-e-insere passa na primeira asserção e **falha na
      contagem**; é essa segunda que prova o histórico.
- [ ] Registrar em `.mnfs/HARNESS-DEBTS.md`: a tabela legal é **global**, então a edição vale para
      todos os tenants. Correto enquanto a alíquota interna for objetiva (é a lei, não preferência) e
      houver um tenant. Multi-tenant real ⇒ a edição exige papel de administrador, **não** uma cópia
      da lei por tenant.

#### C2 — Trocar as operações HTTP (OpenAPI + SDK + handler + FE, **um commit**)

- [ ] `/pricing/difal` e `/pricing/difal/{uf}` (`openapi:2626-2660`) dão lugar a
      `/pricing/aliquotas-internas` e `/pricing/aliquotas-internas/{uf}`, devolvendo **alíquota, FCP
      embutido, procedência (tier), lei, vigência, fonte, quem editou e quando** (decisão 6). O campo
      de edição escreve na **tabela legal** via C1 — não existe tabela de override.
- [ ] Os 5 schemas `PricingDifal*` (`openapi:4566-4633`) e os 3 tipos + 2 métodos do SDK
      (`index.ts:1524-1552`, `:2380-2382`) saem no **mesmo commit**.
- [ ] `DifalDrawer.tsx` vira `AliquotaInternaDrawer.tsx`; `DifalDrawer.test.tsx` é **reescrito**, não
      apagado — cada asserção que existia ganha equivalente na tabela nova.
- [ ] `PricingPage.tsx:129-133` e `:149-156` acompanham.
- [ ] O disclaimer `"seed padrão 2026 — não é orientação fiscal"` (`difal.go:6`,
      `DifalDrawer.tsx:8`) dá lugar ao que a tabela nova sustenta **por linha**: a lei e a data onde
      existem; `'validado pelo operador em 2026-08-04'` + "vigência não informada" no tier
      `operador-validado` (B1). Disclaimer global vira procedência por UF.

#### C3 — Remover as tabelas e a fórmula duplicada

**Migração:** `apps/server_core/migrations/0099_drop_dead_pricing_tables.sql` — **as duas quedas na
mesma migração**, porque são o mesmo fato (tabela de pricing sem consumidor) e assim a contagem de
migrações fecha em 86, não 87.

- [ ] `DROP TABLE pricing_difal_rates`.
- [ ] `DROP TABLE pricing_manual_overrides` (F-12). **Zero código a remover** — a tabela nunca teve
      referência em Go; existe só como `CREATE` no `0004:11` desde a fundação. A wiki já a descreve
      como não sendo fonte de verdade (`wiki/modules/pricing.md:214`) e o parágrafo `:262` — que
      explica que os overrides inline vivem no estado do request do FE — **fica**, porque continua
      verdadeiro; o que sai é a menção à tabela. Precondição a reconferir: `count(*) = 0` (medido
      2026-08-05).
- [ ] `pricing/domain/difal_seed.go` **inteiro** (64 linhas) — terceira cópia da tabela de 27 UFs,
      hardcoded em Go, **sem nenhum consumidor de produção** (F-10). Sai junto com `difal.go` porque
      depende de `InterestadualPct` e `computeEfetivoPct`. É também a constante de negócio em código
      que o gate de auto-revisão proíbe (§9). O comentário do `0057_pricing_difal_seed.sql` que diz
      espelhar `domain.buildDifalSeed` morre com a tabela que ele semeava.
- [ ] `pricing/domain/difal.go`: `DifalRate`, `DifalOverride`, `DifalForUF`, `DifalForUFResult`,
      `computeEfetivoPct`, `interestadual12`, `InterestadualPct` saem. **`aInterUF12`
      (`icms.go:12-14`) fica como conjunto único**, com o comentário reescrito para citar a
      **Resolução do Senado 22/1989, alterada pela 13/2012** — a única parte da conta com base federal
      citável (§1.6.1). O conflito MG entre os dois conjuntos (F-6) morre aqui.
- [ ] `ports/calc_ports.go:39 DifalForUF` e `ports/calc_repository.go` (`RateForUF`,
      `ListDifalRates`, `UpsertDifalOverride`) saem; `adapters/postgres/calc_repository.go:100-176`
      idem; `application/calc_service.go:443-456 applyDifal` idem;
      `domain/solve.go:112-114 difalApplied` idem.
- [ ] `ports/calc_ports_contract_test.go:55-120` congela a porta IC-04 que está sendo retirada. **É
      teste alheio: ele é REESCRITO para congelar a porta nova, nunca apagado** — a propriedade que ele
      protege (override respeitado, tenant escopado, 2 casas decimais preservadas) continua valendo e
      continua testada.
- [ ] **Remover a exceção `production-panic-pricing-difal-efetivo`** de
      `contracts/governance/invariants.json:49-56` **no mesmo commit** — ela aponta para
      `computeEfetivoPct`, que deixa de existir (§6.3).
- [ ] **Bumpar `runner_test.go` de 85 para 86.**
- [ ] **Aceite por observável:** `count(*)` — **zero** tabelas com alíquota interna além de
      **`icms_aliquota_interna`** — uma só; e o mesmo produto/UF devolve o mesmo DIFAL antes e
      depois da fatia, medido nos dois lados.
- [ ] **Aceite de não-perda de informação (F-11), medido antes de dropar:** para as **27** UFs,
      `icms_aliquota_interna.aliquota − fcp_embutido` reproduz `pricing_difal_rates.interna_pct`
      **linha a linha**. Já medido em 2026-08-05: divergem só AL (Δ1,0), RJ (Δ2,0) e SE (Δ1,0), e o Δ
      **é exatamente o FCP** nas três. Reconferir na execução; se alguma linha divergir por outro
      valor, **para** — a tabela que morre carregaria informação que a que fica não tem.

---

### Fatia D — dívidas e medições que ficam abertas

#### D1 — Fechar dívidas

**Fechadas:** `D-41` (era sobre o dado, não o método — **reescrita**, não removida) · `D-45` (a fonte
da previsão existe: `TGFAID`; a fonte legal continua sendo curadoria humana) · `D-46` (a_inter por
`ORIGPROD` confirmado, 784/784) · `D-36`/`D-37` no que toca alíquota ausente (agora é pendência
nomeada) · fila itens 4 e 5.

#### D2 — Registrar as novas em `.mnfs/HARNESS-DEBTS.md` e no ledger da MIS-008

| id | fato medido |
|---|---|
| `D-56` reescrita | nosso gross-up diverge do método do ERP; +9,48 numa nota real de 299,90 |
| **nova** | 20 das 27 linhas legais sem procedência (`0094:85-113`) |
| **nova** | ERP subrecolhe DIFAL onde o cadastro está velho: −9,52 na nota 895507; 10.049 produtos da BA na mesma situação hoje. **Exposição fiscal, não margem** |
| **nova** | `TGFAID` sem trilha de alteração (5 colunas, R4) — passado anterior ao 1º sync é perdido. Mesma classe da `D-28` |
| **nova** | ex-ante ≠ ex-post no ST: média diária (`TGFEFDVMRSTDIA`) vs real da nota (`TGFC185F`). Estrutural |
| **nova** | `TIPCALCDIFAL` é parâmetro de cadastro; base dupla (LC 190/2022) não foi avaliada. Pergunta ao contador |
| **nova** | 7 itens de 6.274 com resíduo positivo na fórmula única, sem causa identificada (`roundV.txt` V2: PE 880915, SP 882116, SC 884972…). Observável que resolve: dump de todos os componentes de uma nota inteira |
| **nova** | 12 itens com PIS alíquota **0** e COFINS 7,60 normal (`roundT.txt` T8) — não é monofásico (zeraria os dois). Provável erro de cadastro do produto |
| **nova** | significado de `TIPIMPOSTO` I/S/B/F em `TGFEFDVMRSTDIA` não medido; `F` nunca tem valor em 4,9 milhões de linhas |
| **nova (harness)** | `0093` duplicado em `apps/server_core/migrations` **sem** entrada em `invariants.json`, ao contrário do `0021` que tem. `0095` ausente |
| **nova (harness)** | `pricing -> tenant_config` (`product_fiscal_reader.go:10`) não está em `modules.json:20` nem nas exceções — aresta viva não registrada |
| **nova** | IBS/CBS já valorados nas notas desde 2025-12-10 e **ausentes de todo o nosso cálculo** (`TaxComponents`, `Decomposition`, DTO, OpenAPI, SDK). ~1% do preço hoje, subindo pela transição |
| **nova (F-13)** | **42 dos 104 métodos do `sdk-runtime` não têm chamador** em `apps/web/src`, `packages/feature-*` nem `packages/web-query` — 40% da superfície de contrato. **Fora do escopo desta missão**; é varredura própria. Instrumento que produziu o número (o único válido — comparar `operationId` do OpenAPI com o FE dá ruído, porque o SDK é escrito à mão e renomeia): para cada `async <nome>(` exportado de `packages/sdk-runtime/src/index.ts`, procurar o nome nos três pacotes acima. Antes de apagar qualquer um, separar **"nunca teve tela"** de **"tela foi removida"** — os dois casos se parecem e só o segundo é lixo puro |
| **nova (F-12)** | `pricing_manual_overrides` (`0004:11`) é a **única** tabela órfã do schema. Fechada por C3 nesta missão; a dívida fica registrada porque o padrão — quatro mecanismos de override no mesmo módulo, nenhum em uso — é o que precisa não se repetir |
| `D-54` | ES ambíguo de verdade; venda para ES sai não-calculável |
| `D-55` | MG interna por produto; `TGFAID` não cobre MG (é a origem) |
| `D-17` | `CODEMP = 1` fixo — inalterado por esta fatia |

#### D3 — Medição própria: `red_base` (não implementar às cegas)

**Fato:** `icms_matrix_mirror.red_base` é não-nulo em **63%** das linhas (`0094:47`) e **nenhum
cálculo o lê** (F-5). Na nota 895507 o ERP **não** reduziu (`BASERED = BASE = 299,90`) — mas isso é
uma nota, e uma nota já fabricou veredito errado nesta missão.

- [ ] Rodada nova contra o Oracle vivo: para os itens de 2026 cuja célula `(uf_destino, grupo)` tem
      `REDBASE ≠ 0`, comparar `TGFDIN.BASERED` do ICMS com `BASE`. Se forem iguais em ~100%, `red_base`
      é campo de cadastro sem efeito nesta operação e vira comentário. Se divergirem, é componente do
      cálculo e vira task própria.
- [ ] **Não escrever código de `red_base` antes desta medição.** Implementar 63% das células com uma
      regra suposta é exatamente a classe de erro que esta missão está desfazendo.

---

### Balanço de linhas — honesto, não otimista

O pedido era *diminuir linha, não só somar*. Onde o plano diminui, medido com `wc -l`:

| sai | linhas | tipo |
|---|---|---|
| `pricing/domain/difal.go` | 90 | **apagado** |
| `pricing/domain/difal_seed.go` | 64 | **apagado** (F-10, já era morto) |
| `pricing/domain/difal_test.go` | 201 | **apagado** com o que ele testava |
| `adapters/postgres/calc_repository.go:100-176` | ~77 | apagado |
| `application/calc_service.go:443-456` + `ports/calc_ports.go:39` + `domain/solve.go:112-114` | ~25 | apagado |
| blocos DIFAL do OpenAPI (`:2626-2660`, `:4566-4633`) e do SDK (`:1524-1552`, `:2380-2382`) | ~110 | apagado |
| `DifalDrawer.tsx` + `.test.tsx` | 170 | **reescrito, não apagado** — não conta como redução |
| `ports/calc_ports_contract_test.go:55-120` | 66 | **reescrito, não apagado** |

**~570 linhas saem de vez** do módulo `pricing`, mais duas tabelas do schema. O que entra em `pricing`
é pequeno: a correção da F-3 troca uma divisão por uma atribuição, e a F-4 substitui um split errado.

**A parte que aumenta, dita sem maquiar:** a Fatia B constrói um pipeline de sync novo em
`internal_read` (`TGFAID` → `erp_aliquota_interna_prevista` → gate de pendência). Isso é **capacidade
que não existe hoje**, não duplicação — sem ela o operador continua descobrindo a divergência de
cadastro depois da nota emitida, que é o motivo da missão. Somando tudo, a contagem total de linhas do
repo provavelmente **sobe um pouco**.

O número que importa não é esse. É este: **cópias do conceito "alíquota interna por UF" vão de 3 para
1** (`pricing_difal_rates` + `difal_seed.go` + `icms_aliquota_interna` → só a última), e **mecanismos
de override vão de 4 para 0**, substituídos por escrita versionada na própria tabela legal. Duplicação
conceitual é o que fez esta missão existir; contagem bruta de linhas é o sintoma, não a doença.

---

## Parte 8 — O roadmap, revisto

**Estado medido:** Onda 0 fechada como código (`F-A1`/`F-A2` @`be8fc56c`, `F-A3` @`afb6b54a`);
`F-00` bloqueada por D-16; P2.b na `main` @`65aacc7`; matriz de ICMS com relógio de verdade
(`root.go:723-747`). **Onda 1 não começou.**

O design da MIS-008 (§7) trava `F-06` (consolidar margem) e `F-09` (API entrega a linha pronta) em
`D-17`(fiscal). A medição mostra que o travamento não era `CODEMP`: era **método errado no nosso
código, dado velho no do ERP, e nenhum dos dois visível em lugar nenhum**.

**Onda 0.5 — "verdade fiscal" — entre a Onda 0 e a Onda 1**, com as fatias A–D.

- **`F-06` consolida três motores de margem em um.** Consolidar sobre entrada errada produz **um**
  número errado no lugar de três, e apaga a evidência de que estava errado.
- **`F-09` publica a linha pronta.** Publicar SIMPLES 4% fabricado num contrato é caro de desfazer.
- **Disjunção com a Onda 1 por arquivo:** esta onda toca `pricing/*`, `internal_read/*`,
  `sync/domain/sync_state.go`, `composition/root.go`, 3 migrações, OpenAPI + SDK e
  `apps/web/src/pages/precos/*`. A Onda 1 toca `orders/application/service.go`,
  `batch_orchestrator.go`, `pedidosFormatters.ts`, `mercadoFormatters.ts`, **OpenAPI + SDK** e
  **`composition/root.go`**. **Interseção NÃO é vazia**: os dois seams exclusivos (o par api-sdk e o
  `root.go` de ~940 linhas) são disputados. **Serializar por seam ou dar o dono a uma das duas —
  medindo com `git worktree list` e `git diff --name-only main...<branch>` antes de despachar em
  paralelo.** A revisão 3 dizia "interseção esperada vazia"; era alegação, e a medição a desmente.
- **Ondas 2 e 3 inalteradas.** `F-01`–`F-04`, `F-10`, `F-A4`, `F-A5` são independentes do fiscal.

**O que repensar, não só evoluir:** o design tratava o fiscal como bloqueio de outras fatias. Ele é a
**linha de chegada** de metade delas — margem verdadeira é o produto. E a rodada R mostrou o que o
design não previa: **existe hoje um processo manual, produto a produto, corrigindo o cadastro do ERP
depois que a nota sai errada.** A Onda 0.5 é a primeira vez que esse processo ganha instrumento.

---

## Parte 9 — Gates de auto-revisão (fase 6)

- [x] Todo `file:line` citado foi **aberto nesta sessão**, na árvore `main` @`14360d34`.
- [x] Nenhum artefato novo sem citação do que existe e por que não serve (Parte 3), incluindo os três
      artefatos novos com justificativa medida.
- [x] Local vs global respondido por escrito (Parte 4) — e a revisão 1 foi refutada **por medição**
      (R3: `N_COM_RED = 0` mata a mudança de dimensão da porta), não por opinião.
- [x] **Nenhuma constante de negócio hardcoded sobrevive:** o `18` de MG sai (A2-d); a alíquota vem da
      tabela com lei; o `a_inter` é regra da origem com a Res. Senado citada; o método é medido.
      O `4%` de `calcprofile.go:19` deixa de alcançar o caminho fiscal (B4).
- [x] Nenhum desconhecido vira `0`/default: UF **sem linha na tabela**, célula ausente, célula ambígua
      e destino não resolvido saem como pendência nomeada (A1 casos 9–10, B4, B5). UF **com** linha
      calcula, e a tela diz de que tier veio o número — procedência fraca é informação publicada,
      nunca ausência disfarçada.
- [x] Nenhum mock em seam de integração; aceite é `count(*)` no banco, nota emitida e live drive.
- [x] Toda asserção do plano **reprova no código de hoje** — A1 traz as três mensagens exatas; A3, B4
      e B5 trazem o vermelho nomeado.
- [x] Teste alheio alterado é **reescrito**, nunca apagado: `decompose_icms_golden_test.go` e
      `icms_test.go` (A2), `icms_matrix_shape_test.go` (B1), `calc_ports_contract_test.go` e
      `DifalDrawer.test.tsx` (C).
- [x] **OpenAPI + SDK numa task, num commit** — B4 e C2, os dois explicitamente.
- [x] `tenant_id` no predicado de toda consulta nova; a única exceção (`icms_aliquota_interna`,
      global) está **dita e justificada** (§6.2).
- [x] Falha do sync novo é visível na tela, com o caminho **falha → pixel** nomeado (B3).
- [x] Sem chamada a provider; o balde de 900/min fica intacto; orçamento na Parte 2 pergunta 7.
- [x] Lanes copiam-e-colam, com diretório (§6.8), e o `BaseSha` de 40 dígitos está escrito.
- [x] Governança planejada: **3 migrações → 3 prefixos únicos a partir de `0097`**, e o fixture de
      contagem em `runner_test.go:25` **e** `:64` bumpado **em cada uma** (83 → 86). Nenhum módulo
      novo. A exceção `production-panic-pricing-difal-efetivo` sai no commit que a torna obsoleta.
      Nenhuma ocorrência nova da aresta `pricing -> tenant_config`.
- [x] Nada aqui dá push, reset, revert, stash, clean, instala dependência, despeja ambiente, escreve
      no Oracle, sobe servidor/porta/`.env`, nem faz escrita viva em provider.

**Pendências honestas deste plano:**

1. **20 UFs calculam com procedência `operador-validado`, não com lei citada.** Ratificado, e visível
   na tela por tier (B1). O que continua aberto é a curadoria incremental do contador — que **não
   bloqueia nada** e não tem data.
2. **Decisão 5 foi REVISADA por medição depois de ratificada.** O operador perguntou se duas tabelas
   não eram máximo local. Eram: `pricing_difal_rates` tem **0 overrides em 27 linhas**, então a tabela
   de override some e o desenho fica com **uma** tabela — menos uma migração, menos um join, e com
   histórico de auditoria que o desenho de duas não teria. **A pergunta do operador encolheu o plano;
   as minhas quatro razões defendiam o artefato errado.** Precondição a reconferir na execução: o
   contador de overrides continua `0`.
3. **`red_base` (D3) é medição, não implementação** — deliberadamente.
4. **A interseção de seam com a Onda 1 não é vazia** (Parte 8). Isso precisa de adjudicação do hub
   antes de qualquer despacho paralelo.
5. Prefixos de migração e âncoras de linha foram medidos hoje; **reconferir na hora de escrever o
   código** — números de linha apodrecem entre branches.
