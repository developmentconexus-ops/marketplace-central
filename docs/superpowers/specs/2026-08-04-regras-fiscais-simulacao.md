# Regras fiscais da simulação — matriz completa, medida

**Data:** 2026-08-04 · **Missão:** MIS-008
**Evidência:** `.mnfs/MIS-008/evidence/fiscal-2026-08-04/` rodadas S, T, U, V (Oracle vivo
`METALPRD`, somente leitura). Todo número aqui tem contagem e origem. O que não foi medido está
marcado **NÃO MEDIDO**.

Este documento define **todas as regras antes de planejar**, a pedido do operador. O plano
(`2026-08-04-fiscal-verdade-tgfaid-plan.md`) implementa o que está aqui.

---

## 1. A fórmula é UMA só

```
base do item          = VLRTOT − VLRDESC
ICMS próprio          = base × a_inter          (0 quando CST 60)
DIFAL                 = base × (a_interna_dest − a_inter)   (0 quando intra-MG ou CST 60)
FCP                   = base × perc_fcp         (só RJ, 2%)
ST anterior           = valor já pago a montante, embutido no item

base PIS/COFINS       = base − ICMS − DIFAL − ST_anterior
PIS 1,65% + COFINS 7,60% sobre essa base
base IBS/CBS          = base PIS/COFINS − PIS − COFINS
IBS 0,1% + CBS 0,9%
```

**Medido (`roundV.txt` V1, itens de 2026):**

| destino | fecha (≤ R$ 0,02) | não fecha |
|---|---|---|
| MG (interna) | **5.995** | 3 |
| fora de MG | **263** | 13 |
| **total** | **6.258 (99,7%)** | 16 |

É o que o operador sempre disse: **tudo que é ICMS sai da base do PIS/COFINS** — próprio, DIFAL,
e o ST que já foi pago antes. Um só conceito, três nomes.

**Exceção medida — o FCP NÃO sai da base.** Seis dos 16 que não fecham são RJ, e em todos o resíduo
é **exatamente −FCP** (−48,59; −37,19; −16,20; −6,20; −6,00; −3,40 — `roundV.txt` V2). O ERP abate
ICMS + DIFAL e deixa o FCP dentro da base. Nossa fórmula tem que fazer igual.

Os 10 restantes: 2 são teto (ST anterior maior que a própria base — MG 883660, STANT 5.429,94 sobre
VLRTOT 5.624,63 com desconto de 1.687,39), 1 é ST anterior não abatido (GO 896879), e 7 têm resíduo
positivo sem causa identificada. **NÃO MEDIDO** — 7 itens em 6.274. Entram como pendência nomeada,
nunca como ajuste.

---

## 2. Os dois cenários, medidos

O que muda entre "dentro de MG" e "fora de MG" **não é a fórmula — é o CST do produto**, e o CST
já está na matriz que espelhamos (`icms_matrix_mirror.codtrib`, 2.498 células vigentes).

`roundV.txt` V3, itens desde 2025:

| destino | CST | itens | com ICMS | com ST anterior |
|---|---|---|---|---|
| **MG (interna)** | **60** | **16.215 (88%)** | **0** | 11.033 |
| MG (interna) | 00 | 2.178 | 2.178 | 29 |
| MG (interna) | 40 | 3 | 0 | 0 |
| **fora de MG** | **00** | **765 (98%)** | **764** | 422 |
| fora de MG | 60 | 17 | 0 | 9 |
| fora de MG | 10 | 10 | 0 | 6 |

**Leitura:**

- **Dentro de MG, 88% é CST 60** — ICMS já recolhido por substituição tributária lá atrás.
  **Zero ICMS na venda.** O único ICMS na conta é o **ST anterior**, que sai da base do PIS/COFINS.
  A intuição do operador ("dentro de MG só ST") está certa e agora tem número.
- **Fora de MG, 98% é CST 00** — tributado integralmente: ICMS próprio + DIFAL. Metade desses
  (422/765) **também** carrega ST anterior, e aí os três termos saem da base juntos.
  Confirmado item a item: nota 896886/1 (GO) → ICMS 10,02 + DIFAL 17,17 + ST 11,93 = 39,12 =
  abatido, **resíduo zero** (`roundU.txt` U1).
- **CST 60 fora de MG existe** (17 itens) — produto ainda sob ST no destino. Não é impossível; é raro.
  Não pode virar `else`.

**Consequência de desenho:** a simulação **precisa do CST como entrada**, não só da UF. Mesmo produto,
mesmo preço, dentro ou fora de MG = imposto completamente diferente. A fonte do CST é a matriz que já
temos.

---

## 3. Crédito / ressarcimento de ST — a resposta honesta

**Não aparece na nota.** Medido (`roundT.txt` T1, itens desde 2025):

| coluna | dentro de MG | fora de MG |
|---|---|---|
| `VLRSUBSTANT` (ST anterior, já pago) | 11.062 / 18.396 | 437 / 792 |
| `VLRSUBST` (ST a recolher nesta nota) | **0** | **0** |
| `VLRREPREDST` (ressarcimento) | **0** | **0** |
| `VLRICMSUFDEST` | **0** | **0** |

A empresa **nunca recolhe ST na venda** (é revendedora — o ST foi pago na compra) e **nunca registra
ressarcimento na nota**.

O ressarcimento existe, mas na **apuração mensal do SPED**: `roundT.txt` T6 encontrou
`TGFC185F`, `TGFEFDCC185`, `TGFEFDFC185`, `TGFEFDVMRST`, `TGFEFDVMRSTDIA`. O registro C185 do EFD é
exatamente ressarcimento / restituição / complementação de ST.

`TGFEFDVMRSTDIA` (`roundT.txt` T5) é a média móvel diária, 4,78 milhões de linhas por tipo, de
2020-01-01 a 2026-07-01, com quatro `TIPIMPOSTO`: **I** (1.917.494 com valor), **S** (1.917.046),
**B** (1.898.773) e **F** (**zero com valor, sempre** — `TIPMEDIA` D e I). O significado exato de
cada letra é **NÃO MEDIDO**; o adapter atual (`sync.go`, `sankhyaFiscalSTSQL`) já lê `S` e `R`
e bate com o ERP.

**Para a simulação isso significa:**

- O ST anterior é **custo dentro do produto**, e sai da base do PIS/COFINS. Isso é ex-ante, sabemos.
- O ressarcimento é **crédito mensal**, apurado depois, fora da nota. **Não é conhecível ex-ante por
  item.** A simulação pode mostrá-lo como faixa a partir da média móvel, **nunca como valor certo**.
- Isso é a lacuna estrutural ex-ante/ex-post já registrada. **Não se apaga inventando um número.**

---

## 4. PIS/COFINS — praticamente constante, com 12 exceções

`roundS.txt` S5, desde 2025:

| CODIMP | alíquota | itens |
|---|---|---|
| 6 (PIS) | 1,65 | 23.211 |
| 6 (PIS) | **0** | **12** |
| 7 (COFINS) | 7,60 | 23.223 |

As 12 com PIS zero (`roundT.txt` T8) são todas TOP 306, em MG, PR, RJ e SP, e **têm COFINS 7,60
normal**. PIS zero com COFINS cheio não é monofásico (monofásico zera os dois). **NÃO MEDIDO** por
quê. Classe de pendência, não regra.

**Regra:** 1,65% e 7,60% vêm de configuração, não de literal no código. Item fora disso → pendência.

---

## 5. IBS / CBS — já estão nas notas

`roundS.txt` S6:

| CODIMP | imposto | alíquota | itens | primeira | última |
|---|---|---|---|---|---|
| 12 | IBS UF | 0,1% | 8.496 | 2025-12-10 | 2026-08-04 |
| 13 | IBS municipal | **0** | 8.496 | 2025-12-10 | 2026-08-04 |
| 14 | CBS | 0,9% | 8.496 | 2025-12-10 | 2026-08-04 |

Entraram em 2025-12-10 e estão em toda nota desde então. Base = base PIS/COFINS − PIS − COFINS
(confirmado na nota 895507: 248,92 − 4,11 − 18,92 = 225,89).

**1% do preço hoje, e a alíquota sobe pela transição da reforma.** A simulação inclui, com as
alíquotas em configuração — nunca literal, justamente porque vão mudar.

`CODIMP 9` aparece em 19.000+ linhas e tem **valor zero em 100% delas** (`roundS.txt` S1). Placeholder.
Ignorar, mas não excluir do espelho — se um dia tiver valor, é mudança de regime.

---

## 6. Operações no escopo

`roundS.txt` S10 e `roundT.txt` T9, desde 2025:

| TOP | descrição | notas | tem imposto? |
|---|---|---|---|
| **305** | NFE ENTREGA PEDIDO (NOVA) | 7.151 | sim |
| **313** | PEDIDO ENTREGA FUTURA ECOMMERCE | 41 | **sim** — CODIMP 1 com valor em 28/43 |
| **306** | NFE ENTREGA E-COMMERCE | 39 | sim |

**313 não é "só pedido"** — carrega imposto. Eu tinha assumido que era só o pedido. Entra no escopo.

---

## 7. Alíquota interestadual — origem MG

Regra da origem, já implementada e correta (`pricing/domain/icms.go:186 aInterFor`, backtest 184/184):

| condição | a_inter |
|---|---|
| `ORIGPROD ∈ {1,2,3,8}` (importado) | **4%** |
| destino ∈ {SP, RJ, PR, SC, RS} | **12%** |
| demais | **7%** |

**Base legal única e federal:** Resolução do Senado nº 22/1989, alterada pela nº 13/2012. É a única
peça da conta com fonte nacional estável — pode ficar em código **com a citação no comentário**.

---

## 8. Alíquota interna do destino — a peça sem fonte única

Pesquisa (2026-08-04, um subagente, resultado integral em §9): **não existe fonte oficial única,
consolidada e mantida.** É estrutural — alíquota interna é prerrogativa de cada estado, não é matéria
de convênio CONFAZ. Cada UF publica no seu próprio site, ou não publica.

Descartados por medição:

- **`github.com/brasil-js/icms`** — 5 commits, **todos de 2016**. Morto. Não usar.
- **IBPT / "De Olho no Imposto"** — resolve outro problema: carga tributária **aproximada por NCM**
  (Lei 12.741/2012), somando federal+estadual+municipal. Não isola a alíquota interna de ICMS.
- **Agregadores comerciais** (brasilnfe.com.br e similares) — têm a matriz completa, mas **sem data de
  atualização e sem citar lei**. Servem como checagem cruzada, não como fonte.
- **BrasilAPI** — não tem endpoint fiscal.
- **CONFAZ** — é repositório de atos legais, não tabela consolidada.

Serve:

- **SEFAZ estaduais**, uma a uma. O RJ mantém página oficial atualizada (23/06/2026) citando
  LC 210/2023 — [portal.fazenda.rj.gov.br/pagamentos/aliquotas-internas](https://portal.fazenda.rj.gov.br/pagamentos/aliquotas-internas/).

**Decisão que isso força:** a tabela legal continua sendo curadoria humana, mas **com citação
obrigatória** — hoje 23 das 27 linhas não têm (`migrations/0094_icms_matrix.sql:85`). O trabalho é
percorrer as 26 SEFAZ e anexar lei + vigência. Não tem atalho, e nenhuma fonte automática dispensa
isso.

### 8.1 As duas disputas, resolvidas pela pesquisa

| UF | ERP | nossa tabela | veredito da pesquisa | fonte |
|---|---|---|---|---|
| **DF** | 18 | **20** | **20% correto**; 18% é o valor pré-2024 | Lei 7.326/2023, vigente 21/01/2024 |
| **RJ** | 19 (+2 FCP) | **20 (+2 FCP) = 22** | **20+2 correto**; "19" não corresponde a nenhum valor histórico (RJ era 18+2 antes de mar/2024) | LC 210/2023, vigente 20/03/2024 — **fonte oficial SEFAZ-RJ** |
| BA | 17 (velho) | **20,5** | confirmado | Lei 14.629/2023, vigente 07/02/2024 |
| MA | 17 (velho) | **23** | confirmado | Lei 12.426/2024, vigente 23/02/2025 |

**Nossa tabela ganha as quatro.** O que eu tinha classificado como "disputa que só o contador
resolve" era, nas duas, cadastro do ERP parcialmente atualizado.

**Ressalva honesta:** BA, MA e DF vieram de fonte secundária (legisweb); só o RJ veio do site da
própria SEFAZ. Antes de gravar `lei` no banco, o número da lei precisa ser conferido no diário
oficial ou no site da SEFAZ do estado. **Pesquisa não é procedência.**

---

## 9. O que a simulação vai mostrar

O operador pediu: *"é bom considerar o pior caso, mas é bom poder ver os outros"*.

Por UF de destino, com o mesmo produto e preço:

| linha | conteúdo |
|---|---|
| **componentes** | ICMS próprio, DIFAL, FCP, ST anterior, PIS, COFINS, IBS, CBS — cada um com sua alíquota e sua base |
| **cenário devido** | com a alíquota da tabela legal. **É o número da margem.** |
| **cenário ERP** | com a alíquota do cadastro espelhado. É o que a nota vai sair se ninguém corrigir |
| **divergência** | a diferença entre os dois, em reais e em pontos de preço, com a lei citada |
| **pior caso** | a UF de maior carga entre as vendáveis |
| **desconhecidos** | UF sem procedência, CST ambíguo (ES), ressarcimento de ST — **nomeados, nunca zerados** |

Regra transversal (ADR-17): componente que não se sabe **não vira zero e não vira média**. Sai
nomeado, e a margem que depende dele sai marcada como incompleta.

---

## 10. Lacunas assumidas

| # | lacuna | o observável que resolve |
|---|---|---|
| 1 | 20 das 27 UFs sem citação de lei | percorrer as SEFAZ; sem atalho automático (§8) |
| 2 | ressarcimento de ST não é conhecível ex-ante por item | é apuração mensal do SPED (§3); faixa, nunca valor |
| 3 | 7 itens de 6.274 com resíduo positivo sem causa | dump dos componentes de uma nota inteira |
| 4 | 12 itens com PIS 0 e COFINS 7,6 | cadastro do produto — provável erro, classe de pendência |
| 5 | significado de `TIPIMPOSTO` I/S/B/F em `TGFEFDVMRSTDIA` | documentação Sankhya ou dump correlacionado com C185 |
| 6 | ES: cadastro contraditório (D-54) | venda para ES sai não-calculável |
| 7 | MG interna por produto (D-55); `TGFAID` não cobre MG | — |
| 8 | `TIPCALCDIFAL` = 0 é parâmetro; base dupla (LC 190/2022) não avaliada | pergunta ao contador |
