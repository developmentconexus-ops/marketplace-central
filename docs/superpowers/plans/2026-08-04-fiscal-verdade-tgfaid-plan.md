# Verdade fiscal ex-ante — a fonte é o ERP, não uma tabela nossa

**Data:** 2026-08-04 · **Missão:** MIS-008 · **Metodologia:** `/mc-planning` (fases 0–6 percorridas)
**Evidência bruta:** `.mnfs/MIS-008/evidence/fiscal-2026-08-04/` (rodadas B, C, F, I, J, K, L, M, N, O, P, Q
contra o Oracle vivo `METALPRD`, somente leitura, via lane Docker)

> Este documento é **relatório + plano**. A parte 1 responde o que o operador pediu: como o sistema
> vai funcionar e por que o erro para de acontecer. As partes 2–7 são o plano no formato da casa.

---

## Parte 1 — Relatório ao operador

### 1.1 A regra que você me deu está confirmada, ao centavo, na sua própria nota

Você disse: *"a base é menos icms e difal se tiver"* — porque *"ambos são ICMS, o que muda é que o
DIFAL é ICMS recolhido para o lugar de destino"*.

Medi na nota **212253 / NUNOTA 895507** (BA, 2026-07-21, produto 41912). O Sankhya guarda a base já
reduzida na coluna `BASERED` do `TGFDIN` — não é cálculo meu, é o número que o ERP gravou
(`roundK.txt`, bloco K6):

| CODIMP | imposto | BASE | **BASERED** | alíquota | VALOR |
|---|---|---|---|---|---|
| 1 | ICMS | 299,90 | 299,90 | 7 | **20,99** |
| — | DIFAL (mesma linha) | 299,90 | — | 10 (= 17 − 7) | **29,99** |
| 6 | PIS | 299,90 | **248,92** | 1,65 | **4,11** |
| 7 | COFINS | 299,90 | **248,92** | 7,60 | **18,92** |
| 12/14 | IBS/CBS | 299,90 | 225,89 | 0,1 / 0,9 | 0,23 / 2,03 |

`299,90 − 20,99 − 29,99 = 248,92`. **Exatamente a sua regra.** E o encadeamento continua: a base de
IBS/CBS é `248,92 − 4,11 − 18,92 = 225,89`.

Isso não é uma nota só. Em 2026, contra o próprio `BASERED`, **PIS bate em 7.617/7.617 itens e COFINS
em 7.627/7.627** (`roundL.txt`, L3). E a fórmula completa, com todos os abatimentos:

```
base do item          = VLRTOT − VLRDESC                     (6.218/6.220 itens de 2026; 40/40 no TOP 306)
ICMS próprio          = base × a_inter                        a_inter ∈ {4 importado, 7, 12}
DIFAL                 = base × (interna_destino − a_inter)    DIFERENÇA SIMPLES — 316/316 itens
FCP                   = base × perc_fcp                       só RJ, 2%
base PIS/COFINS       = base − ICMS − DIFAL − FCP − ST_retido
PIS 1,65% + COFINS 7,60% sobre essa base
base IBS/CBS          = base PIS/COFINS − PIS − COFINS
```

**A forma da sua regra já está no nosso código** (`pricing/domain/icms.go:177-178`): a base do
PIS/COFINS é `P × (1 − a) − S`, que é literalmente `P − ICMS − DIFAL − S`. Não é aí que erramos.

### 1.2 Onde erramos — dois defeitos, e o segundo é maior do que eu pensava

**Defeito A — gross-up.** `icms.go:154-158` calcula a alíquota total como
`a = a_int × (1 − a_inter) / (1 − a_int)` (imposto "por dentro"). O ERP usa **diferença simples**:
`DIFAL = base × (ALIQINTDEST − ALIQUOTA)`, com `TIPCALCDIFAL = 0` e `BASEDIFAL = BASE` em
**316/316** itens medidos. Foi uma decisão registrada (D-41, "DIFAL pela lei, não pelo ERP") com
consequência aceita de 2% a 10% de divergência. A medição de hoje mostra que a consequência é maior
e cai toda para o lado errado: **nós superestimamos o imposto, logo subestimamos a margem**.

**Defeito B — a tabela de alíquota interna é estruturalmente incapaz de acertar.** Nós semeamos
27 linhas à mão em `icms_aliquota_interna` (`migrations/0094_icms_matrix.sql:85`), uma por UF. Medi
de onde o ERP tira esse número: **`METALPRD.TGFAID`**, chaveada por **(CODUF, CODPROD)** — alíquota
interna **por produto e por UF**, não por UF.

E ela explica tudo:

| medição | resultado |
|---|---|
| TGFAID reproduz a `ALIQINTDEST` cobrada nas notas (TOP 305/306 desde 2025) | **784 de 784. Zero divergências, zero linhas faltando** (`roundM.txt`, M3/M4) |
| tamanho | 1.051.700 linhas · 40.450 produtos · 26 UFs (`roundM.txt`, M2) |
| cobertura do nosso catálogo vendável (10.051 produtos `USOPROD='R'`, `ATIVO='S'`) | **10.051 de 10.051, com as 26 UFs cada** — mín. = máx. = 26 (`roundN.txt`, N1) |
| a BA da sua nota | produto 41912 → **17**; produto 39563 → **20,5**. Mesma UF, mesmo grupo fiscal, dois dias seguidos (`roundL.txt` L4, `roundI.txt` I3) |

Ou seja: o que eu tinha classificado como "3 UFs divergentes" e "a BA varia por produto" era o mesmo
fato o tempo todo — **a alíquota interna do destino é um dado por produto, e o ERP tem esse dado
completo para todo o nosso catálogo.**

Comparando nossa tabela semeada com o que o ERP vai cobrar (moda do TGFAID no catálogo vendável,
`roundQ.txt` Q1):

| UF | nossa tabela | TGFAID | Δ | | UF | nossa tabela | TGFAID | Δ |
|---|---|---|---|---|---|---|---|---|
| AC | 19,0 | 17 | **+2,0** | | PB | 20,0 | 17 | **+3,0** |
| AL | 21,5 | 17 | **+4,5** | | PE | 20,5 | 17 | **+3,5** |
| AM | 20,0 | 17 | **+3,0** | | PI | 22,5 | 17 | **+5,5** |
| AP | 18,0 | 17 | **+1,0** | | PR | 19,5 | 18 | **+1,5** |
| BA | 20,5 | 17 | **+3,5** | | RJ | 22,0 | 19 (+2 FCP) | **+1,0** |
| CE | 20,0 | 17 | **+3,0** | | RN | 20,0 | 17 | **+3,0** |
| DF | 20,0 | 18 | **+2,0** | | RO | 19,5 | 17 | **+2,5** |
| ES | 17,0 | 17 | 0 | | RR | 20,0 | 17 | **+3,0** |
| GO | 19,0 | 19 | 0 | | RS | 17,0 | 17 | 0 |
| MA | 23,0 | 17 | **+6,0** | | SC | 17,0 | 17 | 0 |
| MS | 17,0 | 17 | 0 | | SE | 20,0 | 17 | **+3,0** |
| MT | 17,0 | 17 | 0 | | SP | 18,0 | 18 | 0 |
| PA | 19,0 | 17 | **+2,0** | | TO | 20,0 | 17 | **+3,0** |

**19 das 26 UFs divergem.** MG não está no TGFAID (é a UF de origem; a interna de MG é caso à parte).

### 1.3 Quanto custa, na sua nota

Preço 299,90, MG → BA, produto 41912, origem nacional:

| cenário | ICMS próprio | DIFAL | PIS/COFINS | **total de tributos** | Δ |
|---|---|---|---|---|---|
| **ERP (o que foi realmente cobrado)** | 20,99 | 29,99 | 23,03 | **74,01** | — |
| só o gross-up errado (com interna 17 correta) | 20,99 | 36,14 | 22,46 | 79,59 | +5,58 |
| só a tabela errada (20,5, com diferença simples) | 20,99 | 40,49 | 22,05 | 83,53 | +9,52 |
| **nosso código hoje (os dois juntos)** | 20,99 | 50,93 | 21,09 | **93,01** | **+19,00** |

**19,00 em 299,90 = 6,34% do preço, inteiramente contra a margem.** Os dois erros se compõem: juntos
custam mais que a soma deles.

Em GO — 725 dos 784 itens medidos, o destino real do negócio — nossa tabela **acerta** (19) e sobra
só o gross-up: ERP `7% + 12% + 9,25%×(1−0,19) = 26,49%` do preço contra os nossos `29,05%`.
**2,56 pontos de preço** em toda venda interestadual para GO.

Numa operação onde *"margem marketplace é muito apertada"*, 2,5 a 6 pontos de preço é a diferença
entre recusar um pedido bom e aceitar um ruim.

### 1.4 Como o sistema vai funcionar

Cinco fontes, cada uma com dono único e nenhuma inventada:

| grandeza | fonte | estado |
|---|---|---|
| `a_inter` (interestadual) | regra da origem MG: 4% importado (`ORIGPROD ∈ {1,2,3,8}`), 12% para SP/RJ/PR/SC/RS, 7% no resto | **já existe e está certo** — `icms.go:186`, backtest 184/184 interestadual |
| interna do destino | **`TGFAID(CODUF, CODPROD)` espelhada** | **a construir** — hoje é tabela semeada à mão |
| FCP | `TGFICM.PERCICMSFCP`, **já espelhado** | RJ 142/142 células com FCP, todas as outras UFs zero — bate com as notas |
| ST / `CODTRIB` | `icms_matrix_mirror` | já existe, 2.498 células vigentes, 2 ambíguas (D-54) |
| ST retido / restituição | `TGFEFDVMRSTDIA` espelhada | já existe; `TGFC185F` é a verdade ex-post (ver §1.5) |

Com isso o fluxo fica:

1. **Simulação de preço.** Você escolhe o produto e a UF de destino. O sistema resolve
   `a_inter` + interna(produto, UF) + FCP + ST e devolve **o mesmo número que a nota vai cobrar**.
   Margem verdadeira, não estimativa.
2. **Pedido entra.** A UF do comprador vira a célula fiscal do pedido; o custo do pedido usa a mesma
   conta. Nada de alíquota de regime genérica.
3. **Antes de faturar.** O sistema compara **previsto (TGFAID espelhado)** × **devido (tabela legal
   nossa)**. Divergiu → pendência nomeada na tela, com produto, UF, as duas alíquotas e a diferença
   em reais. É aqui que *"garantir que o ERP não vai gerar coisa errada"* acontece: você corrige o
   cadastro no Sankhya **antes** da nota sair.
4. **Depois de faturar.** Reconciliação contra o `TGFDIN` real da nota. Divergência vira registro,
   não some.

### 1.5 Por que o erro para de acontecer — e onde ele **não** para

Para de acontecer porque:

- **A fonte deixa de ser uma alegação nossa e passa a ser o cadastro que o ERP obedece.** 784/784 não
  é "quase sempre"; é a definição operacional de "certo". Uma tabela semeada à mão envelhece em
  silêncio (D-45: nenhum órgão publica tabela consultável); um espelho envelhece **com data e com
  contador de falha no `sync_state`**, igual aos outros seis espelhos.
- **A dimensão passa a ser a certa.** Nenhuma quantidade de manutenção manual faz 27 linhas
  representarem 261.326 pares (produto, UF). A BA da sua nota já provava isso.
- **O gross-up morre e vira uma conta só, com dois nomes na tela.** ICMS e DIFAL passam a ser duas
  linhas derivadas do mesmo `a`, nunca dois cálculos.
- **Divergência vira pendência, não vira número.** ADR-17: célula ausente, ambígua, ou cadastro
  contraditório sai como desconhecido nomeado. Nunca como zero, nunca como aproximação.

**Não para de acontecer, e é honesto dizer:**

- **Ex-ante ≠ ex-post no ST.** Antes da nota, o ST retido só é conhecível pela média diária
  (`TGFEFDVMRSTDIA`). Na hora de faturar, o ERP prefere o valor real da nota (`TGFC185F`). São
  grandezas diferentes por natureza — a simulação carrega faixa, não certeza. Isso vira dívida
  registrada, não é apagável.
- **Se o cadastro do Sankhya estiver errado, nós vamos acertar o número errado.** É exatamente por
  isso que o gate pré-nota compara com a tabela legal em vez de jogá-la fora. As 19 UFs divergentes
  da §1.2 são **duas hipóteses ainda não separadas**: ou nossa semeadura está errada, ou o cadastro
  do ERP está defasado e a empresa está recolhendo DIFAL a menos. **Só você pode decidir qual
  investigar primeiro** — é decisão fiscal, não técnica (§5, decisão 2).
- **De onde sai o FCP do RJ eu medi; qual regra o ERP usa para aplicá-lo, não.** O valor bate
  (`TGFICM.PERCICMSFCP`, RJ 2%, já espelhado), mas `TGFAID.PERCICMSFCP` está nulo nesses produtos e
  o `TGFICM_INT_TRIB` está vazio (`roundP.txt`). Registro como desconhecido nomeado, com o
  observável que o resolveria: uma nota nova para AL, PB ou SE (as outras UFs com FCP geral).
- **ES continua não-calculável** enquanto D-54 viver: cadastro contraditório de verdade
  (`I/199+S/-1` diz CODTRIB 10 / 12%, `N/0+I/199` diz CODTRIB 0 / 7%), pesando 683 + 2.592 produtos.
- **MG interna** tem alíquotas por produto que a matriz não resolve (D-55, 57 de 705 itens). O TGFAID
  não cobre MG. Fica em aberto.

---

## Parte 2 — Medição (as 8 perguntas, com `file:line`)

**1. O que já existe que faz parte disto?**
`pricing/domain/icms.go:101 TaxesForItem` já monta a cadeia inteira (a_inter → a → ICMS → DIFAL →
base PIS/COFINS → ST). `pricing/domain/difal.go:DifalForUF` já implementa **a diferença simples**
(`computeEfetivoPct = max(interna − interestadual, 0)`) sobre `pricing_difal_rates` (27 linhas, com
override do operador). `internal_read/adapters/oracle/icms_matrix.go:53` já é o molde exato de
extração TGFICM → espelho versionado. `internal_read/adapters/oracle/sync.go:~276
sankhyaFiscalSTSQL` já traz `S` e `R` do `TGFEFDVMRSTDIA` e **bate com o ERP** — nada a mudar lá.

**2. Onde o defeito realmente mora?**
Dois sítios, ambos em `pricing/domain/icms.go`: (a) linha **154-158**, o `Quo(num, den)` do gross-up;
(b) a dependência de `ICMSCell.AliquotaInterna` (`icms.go:67`), alimentada por
`pricing/adapters/postgres/matrix_reader.go:53 AliquotaInternaFor(ctx, uf)` — assinatura **só com
UF**, que é a dimensão errada. O defeito visível (margem errada na tela) não mora na tela.

**3. Quem mais consome esse caminho?**
`TaxesForItem` tem **um** chamador de produção: `orders/adapters/pricingtax/reader.go:83`.
`ICMSCell` é construído em **um** sítio de produção: `reader.go:201`. Porta:
`orders/ports/tax_reader.go:44 TaxesForItems`, chamada de
`orders/application/enrich_service.go:424`. `AliquotaInternaFor` é declarada em
`pricing/ports/tax_matrix.go:36`, implementada em `matrix_reader.go:53`, mais um fake em
`orders/adapters/pricingtax/reader_test.go:34`. **O simulador de preço não passa por nada disso.**

**4. O que o contrato já diz?**
`AliquotaInternaFor(ctx, uf string) (*string, error)` — `pricing/ports/tax_matrix.go:36`. Trocar a
dimensão para (uf, codprod) **muda a porta**, e a porta tem teste de integração próprio
(`matrix_reader_integration_test.go:124` e `:159`) que precisa ser reescrito, não apagado.

**5. Qual é o estado vivo real?**
`icms_aliquota_interna`: 27 linhas vigentes (medido em Postgres hoje; colunas reais são
`uf, aliquota, fcp_embutido, fonte, lei, vigente_desde, vigente_ate` — **não** `interna_pct`, que eu
errei ao recordar e o `psql` reprovou). `icms_matrix_mirror`: 2.498 vigentes não-ambíguas + 2
ambíguas; RJ com `perc_fcp > 0` em 142/142 células, todas as demais UFs zero. `pricing_difal_rates`:
27. `pricing_calc_profiles`: **0 linhas**. `products_mirror`: 10.638. `TGFAID` no Oracle:
1.051.700 linhas / 40.450 produtos / 26 UFs.

**6. O que prova que está quebrado hoje e o que provará consertado?**
Quebrado: para o produto 41912 e UF `BA`, o nosso caminho devolve DIFAL 50,93 e o `TGFDIN` da nota
895507 registra **29,99** (`roundK.txt` K6). Consertado: o mesmo par devolve 29,99, e um backtest
sobre os 784 itens interestaduais de 2025-2026 bate `ALIQINTDEST` em 784/784. **Contagem no banco e
comparação contra nota emitida — nunca lane verde.**

**7. Orçamento de custo?**
Zero chamada a provider. Um `SELECT` a mais no Oracle por sync (o TGFAID restrito aos `CODPROD` que
existem em `products_mirror` — ~261 mil linhas, contra as ~10 mil da matriz atual). Nenhuma consulta
extra por render de tela: a leitura é por (uf, codprod) em índice.

**8. O que quebra em silêncio às 3 da manhã?**
Se o sync do TGFAID falhar, a alíquota interna congela na última versão vigente e **a tela continua
verde** — mesma classe do §5.1 do design da MIS-008 (falha de token invisível). Por isso o espelho
entra no `sync_state` com `entity = 'icms_aliquota_interna'` e a idade aparece na tela junto do
componente, igual ao que a Onda 0 fez com `/mercado`. Sem isso a fatia está incompleta.

---

## Parte 3 — O que já existe (varredura anti-redundância)

Rodado: `grep -rn "aliquota_interna\|AliquotaInterna\|ICMSCell\|DifalForUF" apps/server_core`.

| existe hoje | `file:line` | por que não serve como está |
|---|---|---|
| `DifalForUF` — diferença simples correta | `pricing/domain/difal.go` | Serve a fórmula, **não** a fonte: lê `pricing_difal_rates`, outra tabela, também por UF só. É a **segunda** implementação de DIFAL do repo, e a que está certa na conta. Colapsar as duas é o trabalho, não criar uma terceira. |
| `TaxesForItem` | `pricing/domain/icms.go:101` | Estrutura certa, um parâmetro errado (`a`). **Editar**, nunca duplicar. |
| `ICMSMatrixReader` | `internal_read/adapters/oracle/icms_matrix.go` | É o **molde** do adapter novo: extração crua + resolução pura + espelho versionado. Copiar a forma, não o conteúdo. |
| `sankhyaFiscalSTSQL` | `internal_read/adapters/oracle/sync.go` | Já correto e medido contra o ERP. **Nenhuma mudança.** |
| `icms_aliquota_interna` (27 linhas) | `migrations/0094_icms_matrix.sql:26` | **Não morre.** Vira a tabela *legal* do gate pré-nota (previsto × devido). Deixa de ser a fonte do cálculo. |
| `syncapp.Scheduler` | seam de sync existente | O job novo registra nele. Segundo ticker seria a redundância. |

**Nenhum artefato novo sem citação:** o único artefato realmente novo é a tabela espelho
`icms_aliquota_interna_erp` (nome a fixar na task) + seu adapter. Justificativa: nenhuma tabela
existente tem a chave (uf, codprod) — `icms_matrix_mirror` é (uf, grupo) e `icms_aliquota_interna` é
(uf).

---

## Parte 4 — Máximo local vs global

1. **Quantas cópias do conceito existem?** DIFAL: **duas** (`icms.go` gross-up, `difal.go` diferença
   simples), sobre **duas** tabelas. Alíquota interna: **duas** (`icms_aliquota_interna`,
   `pricing_difal_rates`). ≥2 ⇒ o conserto pertence onde elas convergem.
2. **A causa está uma camada abaixo do sintoma?** Sim. A margem errada na tela é a alíquota errada no
   domínio, que é a **dimensão** errada na porta. Consertar as 19 linhas da tabela semeada seria o
   máximo local — e continuaria errado na BA, no RJ e em todo produto com exceção.
3. **Já existe seam?** Sim, três: `syncapp.Scheduler`, o padrão de espelho versionado do
   `icms_matrix_mirror`, e `TaxesForItem`. Todos reusados.
4. **Estou estendendo legado?** Não. `pricing_difal_rates` é IC-04 e **não cresce**: ou vira o
   override do gate legal, ou é inventariado e removido em fatia própria. Decisão do operador (§5).

**O conserto global é maior e eu o tomo assim mesmo**, com raio declarado: 1 porta, 1 adapter,
1 domínio, 1 migração, 1 job de sync, 2 testes de integração reescritos. Nenhuma tela muda de forma
nesta fatia — muda o número que ela mostra.

---

## Parte 5 — Decisões que só o operador toma

Nenhuma task abaixo começa sem estas três respostas.

**Decisão 1 — o gross-up (D-41) cede?**
Você decidiu em 2026-08-02 calcular DIFAL "pela lei", por dentro, aceitando divergir do Sankhya.
A medição mostra o preço: +5,58 em 299,90 só por isso, sempre contra a margem, e a divergência não é
teórica — o ERP é quem paga. **Recomendo diferença simples** (o que o ERP faz, 316/316), com a
fórmula legal preservada no gate como "devido" se você quiser vigiar a diferença.

**Decisão 2 — as 19 UFs divergentes: quem está errado?**
Ou a semeadura legal está errada, ou o cadastro do ERP está defasado (e a empresa recolhe DIFAL a
menos, exposta a autuação). Não dá para decidir isso medindo — nenhuma nota foi emitida para 17
dessas UFs. **Recomendo:** espelhar o TGFAID para o cálculo (é o que vai acontecer de fato) **e**
subir o gate legal marcando as 19 como divergência conhecida, para você levar ao contador.

**Decisão 3 — `pricing_difal_rates` (IC-04) morre nesta fatia ou vira dívida?**
Ela tem 27 linhas, override do operador e a fórmula certa. Se o gate legal absorver o override, ela
é redundância pura. **Recomendo:** absorver o override no gate e registrar a remoção como fatia
própria — nunca "depois", que nesta casa nunca aconteceu.

---

## Parte 6 — Tasks

Ordem obrigatória. Cada task commita. Vermelho antes de verde, com o texto da falha esperada escrito
**antes** de rodar.

Lanes (medidas, com diretório):

```bash
cd apps/server_core && GOCACHE="$PWD/.gocache" go test ./internal/modules/pricing/...
```

```bash
cd apps/server_core && GOCACHE="$PWD/.gocache" go test -tags=integration ./tests/integration/...
```

---

### Task 1 — Golden do domínio: a conta do ERP, escrita antes do código

**Arquivo:** `apps/server_core/internal/modules/pricing/domain/icms_erp_golden_test.go` (novo)

Quatro casos **transcritos de nota emitida**, não inventados:

| caso | entrada | esperado |
|---|---|---|
| nota 895507 (BA) | P=299,90, a_inter=7, interna=17, FCP=0, S=0 | ICMS 20,99 · DIFAL 29,99 · base PC 248,92 · PIS/COFINS 23,03 |
| GO típico | P=299,90, a_inter=7, interna=19 | DIFAL 35,99 · base PC 242,92 |
| RJ com FCP | P=299,90, a_inter=12, interna=19, FCP=2 | DIFAL 20,99 · FCP 6,00 (`roundI.txt` I2: 893649 → VLRDIFALDEST 20,99) |
| intra-MG | P=136,66, interna=18 | ICMS 24,59 · DIFAL **0** · base PC 112,07 (`roundJ.txt` I4: 882236/2) |

**Controle negativo nomeado:** rodar contra o código de hoje. O caso BA deve falhar com
`DIFAL = 50.93, want 29.99`. Se falhar com outra mensagem, o teste está errado, não o código.

- [ ] Escrever os quatro goldens.
- [ ] Rodar a lane de unidade e **colar a falha exata** no commit.

---

### Task 2 — Matar o gross-up

**Arquivo:** `apps/server_core/internal/modules/pricing/domain/icms.go:144-178`

Trocar o ramo não-MG por diferença simples:

```go
// a é a carga total do destino: própria + DIFAL. O ERP calcula o DIFAL por
// DIFERENÇA SIMPLES (TGFDIN.ALIQPARADIFAL = ALIQINTDEST − ALIQUOTA, com
// TIPCALCDIFAL = 0 e BASEDIFAL = BASE em 316/316 itens medidos em 2026).
// Gross-up por dentro (D-41) foi medido em +5,58 numa nota de 299,90, sempre
// contra a margem. Evidência: .mnfs/MIS-008/evidence/fiscal-2026-08-04/roundC.txt (B6).
a = aInt
```

`ICMS_oper = P × a_inter` e `DIFAL = P × a − ICMS_oper` continuam derivados do mesmo `a` — uma fonte,
duas linhas. A base do PIS/COFINS (`icms.go:175-178`) **não muda**: `P × (1 − a) − S` já é
`P − ICMS − DIFAL − S`.

- [ ] Aplicar.
- [ ] Task 1 passa. Comentar no commit qual golden virou verde.
- [ ] `decompose_icms_golden_test.go` e `icms_test.go` vão quebrar: os valores esperados foram
      escritos contra o gross-up. **Recalcular e reescrever, nunca apagar** — cada valor novo com a
      conta ao lado no comentário.

---

### Task 3 — Espelho do `TGFAID` (a fatia grande)

**Migração:** `apps/server_core/migrations/0095_icms_aliquota_interna_erp.sql`
(prefixo a confirmar contra o maior existente **no momento de escrever** + o fixture de contagem
`internal/platform/migrate/runner_test.go`).

```sql
CREATE TABLE IF NOT EXISTS icms_aliquota_interna_erp (
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
CREATE UNIQUE INDEX IF NOT EXISTS icms_aliquota_interna_erp_vigente
    ON icms_aliquota_interna_erp (tenant_id, uf_destino, codprod)
    WHERE vigente_ate IS NULL;
```

**Adapter:** `internal_read/adapters/oracle/icms_interna.go`, no molde de `icms_matrix.go`:

```sql
SELECT u.UF, a.CODPROD, a.ALIQINTDEST, a.PERCICMSFCP, a.PERCREDBASEDEST
FROM METALPRD.TGFAID a
JOIN METALPRD.TSIUFS u ON u.CODUF = a.CODUF AND u.CODPAIS = 55
WHERE a.CODPROD IN (<os CODPROD que existem no products_mirror>)
```

Restrição pelo catálogo espelhado: 10.051 × 26 ≈ 261 mil linhas, não 1,05 milhão. `PERCICMSFCP` é
lido mas **medido nulo** em todos os produtos do catálogo — entra como coluna, nunca como fonte do
FCP (que sai de `icms_matrix_mirror.perc_fcp`, RJ 142/142).

**Job:** registrar em `syncapp.Scheduler` com `entity = 'icms_aliquota_interna_erp'` no `sync_state`.

- [ ] Teste de integração hermético: insere duas versões (uma fechada, uma vigente) e prova que o
      reader devolve a vigente. **Controle negativo:** UF/produto sem linha devolve `nil`, nunca
      um default.
- [ ] Rodar e ver falhar antes do adapter existir.

---

### Task 4 — A porta muda de dimensão (contrato)

**Arquivos:** `pricing/ports/tax_matrix.go:36`, `pricing/adapters/postgres/matrix_reader.go:53`,
`orders/adapters/pricingtax/reader.go:201`, `orders/adapters/pricingtax/reader_test.go:34`.

```go
// AliquotaInternaFor devolve a alíquota interna do destino PARA O PRODUTO,
// espelhada de METALPRD.TGFAID (chave real: CODUF + CODPROD). Medido:
// reproduz 784/784 das ALIQINTDEST cobradas em notas de 2025-2026.
// nil = sem linha vigente → pendência nomeada (ADR-17), nunca aproximação.
AliquotaInternaFor(ctx context.Context, uf string, codprod int) (*string, error)
```

- [ ] Reescrever `matrix_reader_integration_test.go:124` e `:159` contra a tabela nova —
      **restaurados, não apagados** (o `:159` é o teste de nil, e ele fica).
- [ ] `reader.go:201` passa o `codprod` que já tem em mãos.
- [ ] Backtest ao vivo: script que compara a alíquota resolvida contra `TGFDIN.ALIQINTDEST` nos 784
      itens. **Aceite = 784/784.** Menos que isso, a task não fecha.

---

### Task 5 — O simulador entra no mesmo seam

**Medido:** `ICMSCell` é construído em **um** sítio de produção — `reader.go:201`, do módulo de
pedidos. O simulador (`pricing/transport/calc_handler.go`) nunca recebe célula; cai em
`decompose.go:136 pctOfPrice(in.AliquotaPct, preco)`, e `AliquotaPct` vem do perfil de regime, cujo
default é `calcprofile.go:19 defaultSimplesAliquotaPct = "4"` — com `pricing_calc_profiles` em
**0 linhas**. **Uma empresa de regime normal simulando com SIMPLES 4%.**

Este é o gate de sítio de composição da §4.3 do handoff: sem esta task, as 1–4 entregam um motor
correto sem consumidor na tela que você mais usa.

- [ ] `DecomposeInput` do simulador recebe `ICMSCell` resolvida por (uf, codprod).
- [ ] Sem UF escolhida: componente **desconhecido nomeado**, não 4%.
- [ ] **Aceite por observável:** simular o produto 41912 para BA em `/precos` e ler na tela
      ICMS 20,99 / DIFAL 29,99 sobre 299,90.
- [ ] **Precondição do live drive:** container do backend construído de um commit **≥** o desta
      fatia. Comparar `CreatedAt` do container com o commit — binário velho faz o live drive mentir.

---

### Task 6 — Gate pré-nota: previsto × devido

- [ ] `icms_aliquota_interna` (27 linhas legais) **permanece** e passa a ser o lado "devido".
- [ ] Divergência entre espelho e legal vira pendência com produto, UF, as duas alíquotas e a
      diferença em reais.
- [ ] As 19 UFs da §1.2 aparecem imediatamente. **Isso é o gate funcionando, não um bug.**
- [ ] Falha do sync do espelho vira idade visível na linha — sem isso a fatia não fecha (§2 pergunta 8).

---

### Task 7 — Dívidas: fechar e abrir, com a medição de cada uma

**Fechadas por esta fatia:**
`D-41` (gross-up cede — se decisão 1 for a recomendada) · `D-45` (tabela sem fonte automática: a
fonte existe e é o TGFAID) · `D-46` (a_inter por ORIGPROD confirmado, 784/784) · fila item 4
(espelhar TGFAID) · fila item 5 (as 27 linhas legais viram gate).

**Abertas / reafirmadas:**

| id | fato medido |
|---|---|
| `D-56` | reescrita: o ERP usa diferença simples (316/316). Deixa de ser "divergência intencional" e vira "defeito quantificado em 6,34% do preço numa nota real". |
| **novo** | ex-ante ≠ ex-post no ST: `TGFEFDVMRSTDIA` (média diária) contra `TGFC185F` (real da nota). Estrutural, não apagável. |
| **novo** | fonte da **regra** de FCP desconhecida. Valor bate por `TGFICM.PERCICMSFCP` (RJ 2%); `TGFAID.PERCICMSFCP` nulo; `TGFICM_INT_TRIB` vazio. Observável que resolve: nota nova para AL, PB ou SE. |
| **novo** | as 19 UFs onde a semeadura legal diverge do cadastro do ERP — risco fiscal a levar ao contador. |
| `D-54` | ES ambíguo de verdade; venda para ES sai não-calculável. |
| `D-55` | MG interna por produto; TGFAID não cobre MG. |
| `D-17` | `CODEMP = 1` fixo — inalterado por esta fatia. |
| `D-28` | histórico da matriz começa no primeiro sync; `TGFHICM` guarda o passado. Vale igual para o espelho novo. |

---

## Parte 7 — O roadmap, revisto

**Estado medido:** Onda 0 fechada como código (`F-A1`/`F-A2` @`be8fc56c`, `F-A3` @`afb6b54a`);
`F-00` bloqueada por D-16; P2.b na `main` @`65aacc7`. **Onda 1 não começou.**

O design da MIS-008 (§7) trazia Onda 1 = remoção, com `F-06` (consolidar margem) e `F-09` (API entrega
a linha pronta) **explicitamente bloqueadas por `D-17`(fiscal)**. Esta medição mostra que o bloqueio
era mais fundo do que o registro dizia: não era só `CODEMP` fixo, era **a alíquota estar na dimensão
errada**.

**Proposta: uma Onda 0.5 — "verdade fiscal" — entre a Onda 0 e a Onda 1**, com as tasks da Parte 6.

Por quê, e não simplesmente empurrar para dentro da Onda 1:

- **`F-06` consolida três motores de margem em um.** Consolidar enquanto a entrada está errada
  produz **um** número errado no lugar de três — e some com a evidência de que estava errado. A
  ordem "remoção antes de conserto" do design vale para código que vai morrer; aqui o que está
  errado é o dado que sobrevive à consolidação.
- **`F-09` entrega a linha pronta para `/precos`.** Entregar a linha com SIMPLES 4% fabricado é
  publicar a mentira num contrato, e contrato publicado é caro de mudar.
- **A Onda 0.5 é disjunta da Onda 1 por arquivo.** Ela toca `pricing/domain`, `pricing/ports`,
  `pricing/adapters`, `internal_read/adapters/oracle` e uma migração. A Onda 1 toca
  `orders/application/service.go`, `batch_orchestrator.go`, `pedidosFormatters.ts`,
  `mercadoFormatters.ts`, OpenAPI + SDK. **Interseção vazia** — mas isso é alegação até rodar
  `git worktree list` e `git diff --name-only main...<branch>` para cada branch em voo. Fazer antes
  de despachar em paralelo.
- **Ondas 2 e 3 não mudam.** `F-01`, `F-02`, `F-03`, `F-04`, `F-10`, `F-A4`, `F-A5` são
  independentes do fiscal.

**O que eu recomendo repensar, não só evoluir:** o design tratava o fiscal como *bloqueio* de outras
fatias. Ele é, na verdade, a **linha de chegada** de metade delas — a margem verdadeira é o produto.
Sugiro promover a verdade fiscal de "dívida que bloqueia" para **onda própria com critério de aceite
em nota emitida** (784/784), e manter a Onda 1 como está, atrás dela.

---

## Parte 8 — Gates de auto-revisão (fase 6)

- [x] Todo `file:line` citado foi aberto nesta sessão.
- [x] Nenhum artefato novo sem citação do que existe e por que não serve (Parte 3).
- [x] Local vs global respondido por escrito (Parte 4).
- [x] Nenhuma constante de negócio nova hardcoded — alíquota vem do espelho; `a_inter` continua regra
      da origem, com D-46 medido.
- [x] Nenhum desconhecido vira `0`/default: sem linha vigente → `nil` + motivo (ADR-17).
- [x] Nenhum mock em seam de integração; o aceite é contagem no banco e nota emitida.
- [x] Toda asserção do plano reprova no código de hoje (Task 1 traz a mensagem exata).
- [x] Teste alheio alterado é **reescrito**, nunca apagado (Task 2 e Task 4).
- [ ] OpenAPI + SDK: **esta fatia não muda contrato HTTP.** Se a Task 5 expuser UF no request do
      simulador, OpenAPI e `sdk-runtime` entram **no mesmo commit**.
- [x] `tenant_id` presente na tabela nova e no índice vigente.
- [x] Falha do sync novo é visível na tela (Task 6), com o caminho falha → pixel nomeado.
- [x] Sem chamada a provider; orçamento declarado na Parte 2 pergunta 7.
- [x] Lanes copiam-e-colam, com diretório.
- [x] Governança: migração nova → prefixo único **e** o fixture de contagem em
      `internal/platform/migrate/runner_test.go`. Módulo novo: nenhum.
- [x] Nada aqui dá push, reset, revert, stash, clean, instala dependência, despeja ambiente, escreve
      no Oracle, nem faz escrita viva em provider.

**Pendência honesta deste plano:** o prefixo `0095` e os números de linha de `icms.go` foram medidos
na árvore de hoje (`main` @`b86e912c`). Confirmar ambos no momento de escrever o código — números de
linha apodrecem entre branches.
